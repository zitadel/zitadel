import { accessSync, constants, readFileSync } from "node:fs";
import { request as httpRequest } from "node:http";
import { createServer, request as httpsRequest } from "node:https";
import { createRequire } from "node:module";

/**
 * Parses a Go-style duration string into milliseconds.
 *
 * This uses the same format as Go's time.ParseDuration (e.g. "30s", "5m",
 * "1m30s") to stay consistent with the ZITADEL backend, which uses
 * time.Duration for its own ZITADEL_DATABASE_POSTGRES_AWAITINITIALCONN
 * setting. Operators configuring both services use the same syntax.
 *
 * Valid units are "ms", "s", "m", and "h".
 * Returns null if the string is empty, zero, or contains invalid syntax.
 */
function parseDuration(raw) {
  const trimmed = raw.trim();
  if (!trimmed) return null;

  let totalMs = 0;
  let remaining = trimmed;
  const unitMap = { ms: 1, s: 1_000, m: 60_000, h: 3_600_000 };

  while (remaining.length > 0) {
    const match = remaining.match(/^(\d+)(ms|s|m|h)/);
    if (!match) return null;
    totalMs += parseInt(match[1], 10) * unitMap[match[2]];
    remaining = remaining.slice(match[0].length);
  }

  return totalMs > 0 ? totalMs : null;
}

/**
 * Polls an HTTP(S) URL until it returns HTTP 200 or the deadline is reached.
 */
async function pollReady(url, timeoutMs) {
  const parsedUrl = new URL(url);
  const doRequest = parsedUrl.protocol === "https:" ? httpsRequest : httpRequest;
  const deadline = Date.now() + timeoutMs;

  while (Date.now() < deadline) {
    const ok = await new Promise((resolve) => {
      const req = doRequest(parsedUrl, { method: "GET", timeout: 5_000, rejectUnauthorized: false }, (res) => {
        res.resume();
        resolve(res.statusCode === 200);
      });
      req.on("error", () => resolve(false));
      req.on("timeout", () => {
        req.destroy();
        resolve(false);
      });
      req.end();
    });

    if (ok) return true;

    await new Promise((r) => setTimeout(r, 1_000));
  }

  return false;
}

/**
 * Waits for the ZITADEL API to become ready before starting the server.
 *
 * When ZITADEL_API_AWAITINITIALCONN is set (e.g. "60s", "5m", "120s"), this
 * function blocks until GET {ZITADEL_API_URL}/debug/ready returns HTTP 200
 * or the timeout expires. This replaces the wait4x init container that was
 * previously used in the Helm chart.
 *
 * If the timeout expires without a successful response, the process exits
 * with code 1.
 */
async function awaitApiReady() {
  const raw = process.env.ZITADEL_API_AWAITINITIALCONN;
  if (!raw) return;

  const apiUrl = process.env.ZITADEL_API_URL;
  if (!apiUrl) {
    console.error("ZITADEL_API_AWAITINITIALCONN is set but ZITADEL_API_URL is not configured.");
    process.exit(1);
  }

  const timeoutMs = parseDuration(raw);
  if (timeoutMs === null) {
    console.error(
      `Invalid ZITADEL_API_AWAITINITIALCONN value "${raw}". Expected a Go-style duration like "30s", "5m", or "1m30s".`,
    );
    process.exit(1);
  }

  const readyUrl = new URL("/debug/ready", apiUrl).toString();
  console.log(`Waiting up to ${raw} for API to be ready at ${readyUrl} ...`);

  const ok = await pollReady(readyUrl, timeoutMs);
  if (ok) {
    console.log("API is ready.");
  } else {
    console.error(`Timed out after ${raw} waiting for API readiness at ${readyUrl}`);
    process.exit(1);
  }
}

await awaitApiReady();

if (process.env.ZITADEL_TLS_ENABLED !== "true") {
  // The build script must rename the Next.js-generated server.js to next-server.js.
  // Validate that the file exists before importing to provide a clearer error if the build step fails.
  try {
    accessSync(new URL("./next-server.js", import.meta.url), constants.R_OK);
  } catch (err) {
    console.error(
      'Unable to find or read "./next-server.js". Ensure the build script has renamed Next.js\'s server.js to next-server.js.',
      err,
    );
    process.exit(1);
  }
  await import("./next-server.js");
} else {
  const certPath = process.env.ZITADEL_TLS_CERTPATH;
  const keyPath = process.env.ZITADEL_TLS_KEYPATH;

  if (!certPath || !keyPath) {
    console.error(
      "ZITADEL_TLS_ENABLED is set to true but ZITADEL_TLS_CERTPATH or ZITADEL_TLS_KEYPATH are missing. Both must be set when TLS is enabled.",
    );
    process.exit(1);
  }

  try {
    accessSync(certPath, constants.R_OK);
  } catch (err) {
    console.error(`Certificate file is not readable: ${certPath}`, err);
    process.exit(1);
  }

  try {
    accessSync(keyPath, constants.R_OK);
  } catch (err) {
    console.error(`Key file is not readable: ${keyPath}`, err);
    process.exit(1);
  }

  const cert = readFileSync(certPath);
  const key = readFileSync(keyPath);

  // Requires Node.js >= 20.11.0 (or >= 21.2.0) for import.meta.dirname support.
  const appDir = import.meta.dirname;

  const requiredServerFilesPath = new URL(".next/required-server-files.json", import.meta.url);

  let requiredServerFiles;
  try {
    accessSync(requiredServerFilesPath, constants.R_OK);
    const requiredServerFilesContent = readFileSync(requiredServerFilesPath, "utf8");
    requiredServerFiles = JSON.parse(requiredServerFilesContent);
  } catch (err) {
    console.error(
      "Failed to load .next/required-server-files.json. The Next.js build may be incomplete or the file may contain invalid JSON.",
      err,
    );
    process.exit(1);
  }
  process.env.__NEXT_PRIVATE_STANDALONE_CONFIG = JSON.stringify(requiredServerFiles.config);

  process.env.NODE_ENV = "production";
  process.chdir(appDir);

  const currentPort = parseInt(process.env.PORT, 10) || 3000;
  const hostname = process.env.HOSTNAME || "0.0.0.0";

  // KEEP_ALIVE_TIMEOUT is expected to be in milliseconds.
  // Reject values that are negative or unreasonably large to avoid resource issues.
  const MAX_KEEP_ALIVE_TIMEOUT_MS = 10 * 60 * 1000; // 10 minutes

  let keepAliveTimeout = parseInt(process.env.KEEP_ALIVE_TIMEOUT, 10);
  if (
    Number.isNaN(keepAliveTimeout) ||
    !Number.isFinite(keepAliveTimeout) ||
    keepAliveTimeout < 0 ||
    keepAliveTimeout > MAX_KEEP_ALIVE_TIMEOUT_MS
  ) {
    keepAliveTimeout = undefined;
  }

  const require = createRequire(import.meta.url);
  require("next");
  // NOTE: This imports an internal Next.js API that is not part of the public, semver-stable surface.
  // - Verified with Next.js version 15.5.10 (see package.json).
  // - This may break on Next.js upgrades; if Next.js is updated, re-check that
  //   `next/dist/server/lib/start-server` still exists and that `getRequestHandlers` behaves as expected.
  // - This is currently required to support HTTPS with self-signed certificates in production mode,
  //   as there is no equivalent public API in Next.js at this time.
  const { getRequestHandlers } = require("next/dist/server/lib/start-server");

  const httpsServer = createServer({ key, cert });

  if (keepAliveTimeout !== undefined) {
    httpsServer.keepAliveTimeout = keepAliveTimeout;
  }

  const { requestHandler, upgradeHandler } = await getRequestHandlers({
    dir: appDir,
    port: currentPort,
    isDev: false,
    server: httpsServer,
    hostname,
    keepAliveTimeout,
  });

  httpsServer.on("request", async (req, res) => {
    try {
      await requestHandler(req, res);
    } catch (err) {
      res.statusCode = 500;
      res.end("Internal Server Error");
      console.error(`Failed to handle request for ${req.url}`);
      console.error(err);
    }
  });

  httpsServer.on("upgrade", async (req, socket, head) => {
    try {
      await upgradeHandler(req, socket, head);
    } catch (err) {
      socket.destroy();
      console.error(`Failed to handle upgrade for ${req.url}`);
      console.error(err);
    }
  });

  httpsServer.on("error", (err) => {
    console.error("Failed to start HTTPS server");
    console.error(err);
    process.exit(1);
  });

  let cleanupStarted = false;
  const cleanup = () => {
    if (cleanupStarted) return;
    cleanupStarted = true;
    httpsServer.close((err) => {
      if (err) {
        console.error(err);
        process.exit(1);
      } else {
        process.exit(0);
      }
    });
  };

  if (!process.env.NEXT_MANUAL_SIG_HANDLE) {
    process.on("SIGINT", cleanup);
    process.on("SIGTERM", cleanup);
  }

  httpsServer.listen(currentPort, hostname, () => {
    const formattedHost = hostname.includes(":") ? `[${hostname}]` : hostname;
    const readyUrl = new URL(`https://${formattedHost}:${currentPort}`);
    console.log(`> Ready on ${readyUrl}`);
  });
}

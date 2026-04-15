/**
 * Mock Zitadel server for testing ZITADEL_API_AWAITINITIALCONN.
 *
 * Returns 503 on /debug/ready until READY_AFTER_MS milliseconds have elapsed,
 * then switches to 200. All other endpoints return minimal ConnectRPC-compatible
 * responses so the login app can start normally once the readiness gate passes.
 */
const http = require("http");

const PORT = process.env.PORT || 8080;
const READY_AFTER_MS = parseInt(process.env.READY_AFTER_MS || "5000", 10);

const startTime = Date.now();

const sendResponse = (res, statusCode, body) => {
  res.writeHead(statusCode, { "Content-Type": "application/json" });
  res.end(JSON.stringify(body));
};

const server = http.createServer((req, res) => {
  const { method, url } = req;
  console.log(`[${new Date().toISOString()}] ${method} ${url}`);

  if (url === "/debug/ready") {
    const elapsed = Date.now() - startTime;
    if (elapsed < READY_AFTER_MS) {
      console.log(`  -> 503 (not ready, ${elapsed}ms / ${READY_AFTER_MS}ms)`);
      sendResponse(res, 503, { status: "not ready" });
    } else {
      console.log(`  -> 200 (ready after ${elapsed}ms)`);
      sendResponse(res, 200, { status: "ok" });
    }
    return;
  }

  if (url.includes("GetLoginSettings") || url.includes("SettingsService")) {
    sendResponse(res, 200, { settings: {} });
  } else if (url.includes("GetBrandingSettings")) {
    sendResponse(res, 200, { settings: {} });
  } else if (url.includes("SessionService") || url.includes("Session")) {
    sendResponse(res, 200, { session: {} });
  } else if (url.includes("UserService") || url.includes("User")) {
    sendResponse(res, 200, { user: {} });
  } else if (url.includes("OrganizationService") || url.includes("Organization")) {
    sendResponse(res, 200, { result: [] });
  } else if (url.includes("IdentityProviderService") || url.includes("IdentityProvider")) {
    sendResponse(res, 200, { identityProviders: [] });
  } else {
    sendResponse(res, 200, {});
  }
});

server.on("error", (err) => console.error("Server error:", err));

server.listen(PORT, "0.0.0.0", () => {
  console.log(`Mock Zitadel server listening on port ${PORT}`);
  console.log(`Will become ready after ${READY_AFTER_MS}ms`);
});

process.on("SIGTERM", () => {
  console.log("Shutting down mock server...");
  server.close(() => {
    console.log("Mock server stopped");
    process.exit(0);
  });
});

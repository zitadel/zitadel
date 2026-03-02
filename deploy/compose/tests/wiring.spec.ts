/**
 * Compose stack wiring tests — per-service reachability through Traefik
 *
 * PURPOSE: verify that every service and protocol is correctly routed through
 * the Traefik reverse proxy. These are INFRASTRUCTURE tests, not feature tests.
 * They answer: "is this service alive and reachable through the proxy?"
 *
 * Services covered:
 *   - Login UI       (zitadel-login, Next.js, /ui/v2/login)
 *   - Console        (zitadel-api, Angular SPA, /ui/console)
 *   - OIDC           (zitadel-api, /.well-known/openid-configuration + /oauth/v2/keys)
 *   - SAML           (zitadel-api, /saml/v2/metadata)
 *   - Root redirect  (Traefik replacepath middleware → /ui/v2/login/)
 *
 * API protocol × transport matrix (all via Traefik catch-all, priority 100):
 *   V1 (AdminService/Healthz — unauthenticated, expects success):
 *     - REST json     HTTP/1.1  /admin/v1/healthz
 *     - REST json     h2c       /admin/v1/healthz
 *     - gRPC-web proto HTTP/1.1 /zitadel.admin.v1.AdminService/Healthz
 *     - gRPC-web json  HTTP/1.1 /zitadel.admin.v1.AdminService/Healthz
 *     - gRPC proto     h2c      /zitadel.admin.v1.AdminService/Healthz
 *     - gRPC json      h2c      /zitadel.admin.v1.AdminService/Healthz
 *
 *   V2 (SessionService/ListSessions — requires auth, expects 401):
 *     Three access styles, each × HTTP/1.1 and h2c:
 *
 *     1. REST-style path  application/json
 *          POST /v2/sessions  (gRPC-gateway REST mapping)
 *
 *     2. gRPC-style path  application/json  (Connect unary JSON)
 *          POST /zitadel.session.v2.SessionService/ListSessions
 *          Handled by the Connect handler; application/json = Connect unary.
 *          Note: application/connect+json is streaming-only (Connect spec).
 *
 *     3. gRPC-style path  application/proto  (Connect unary proto)
 *          POST /zitadel.session.v2.SessionService/ListSessions
 *          Handled by the Connect handler; application/proto = Connect unary.
 *          Note: application/connect+proto is streaming-only (Connect spec).
 *
 *     4. gRPC-style path  application/grpc+proto / application/grpc+json  (gRPC native, h2c only)
 *          POST /zitadel.session.v2.SessionService/ListSessions
 *
 * For end-to-end browser login flow (proves Traefik → Login → API → DB chain):
 *   see smoke.spec.ts
 *
 * For feature-level login tests (MFA, OIDC flows, IdPs, etc.):
 *   see apps/login/acceptance/ — run via @zitadel/login:test-acceptance
 */
import { expect, test, type APIRequestContext, type TestInfo } from "@playwright/test";
import http2 from "node:http2";
import { Buffer } from "node:buffer";

// ---------------------------------------------------------------------------
// HTTP/2 (h2c) helpers
// ---------------------------------------------------------------------------

type H2Response = {
  status: number;
  headers: Record<string, string>;
  body: Buffer;
  trailers: Record<string, string>;
};

function h2Request(
  baseUrl: string,
  method: string,
  path: string,
  reqHeaders: Record<string, string>,
  body?: Buffer | string,
): Promise<H2Response> {
  return new Promise((resolve, reject) => {
    const client = http2.connect(baseUrl);
    const timeout = setTimeout(() => {
      client.destroy();
      reject(new Error(`h2c ${method} ${path} timed out`));
    }, 30_000);

    client.once("error", (err) => {
      clearTimeout(timeout);
      reject(err);
    });

    const req = client.request({
      ":method": method,
      ":path": path,
      ":scheme": "http",
      ...reqHeaders,
    });

    if (body !== undefined) req.write(body);
    req.end();

    let status = 0;
    const chunks: Buffer[] = [];
    const trailers: Record<string, string> = {};
    const responseHeaders: Record<string, string> = {};

    req.on("response", (hdrs) => {
      status = Number(hdrs[":status"]);
      for (const [k, v] of Object.entries(hdrs)) {
        if (!k.startsWith(":")) responseHeaders[k] = String(v);
      }
    });

    req.on("data", (chunk: Buffer) => chunks.push(chunk));

    req.on("trailers", (thdrs) => {
      for (const [k, v] of Object.entries(thdrs)) {
        trailers[k] = String(v);
      }
    });

    req.on("end", () => {
      clearTimeout(timeout);
      client.close();
      resolve({
        status,
        headers: responseHeaders,
        body: Buffer.concat(chunks),
        trailers,
      });
    });

    req.on("error", (err) => {
      clearTimeout(timeout);
      client.close();
      reject(err);
    });
  });
}

function h2Post(
  baseUrl: string,
  path: string,
  reqHeaders: Record<string, string>,
  body?: Buffer | string,
): Promise<H2Response> {
  return h2Request(baseUrl, "POST", path, reqHeaders, body);
}

function h2Get(
  baseUrl: string,
  path: string,
  reqHeaders: Record<string, string> = {},
): Promise<H2Response> {
  return h2Request(baseUrl, "GET", path, reqHeaders);
}

// ---------------------------------------------------------------------------
// gRPC framing helpers (RFC: https://grpc.github.io/grpc/core/md_doc_PROTOCOL-HTTP2.html)
// Frame: [compressed (1 byte)] [message length (4 bytes BE)] [message bytes]
// ---------------------------------------------------------------------------

/** Wrap proto bytes in a gRPC length-prefixed frame. */
function grpcFrame(protoBytes: Buffer = Buffer.alloc(0)): Buffer {
  const frame = Buffer.allocUnsafe(5 + protoBytes.length);
  frame[0] = 0; // compressed-flag = false
  frame.writeUInt32BE(protoBytes.length, 1);
  protoBytes.copy(frame, 5);
  return frame;
}

/** Parse the first gRPC response frame and return the proto body bytes. */
function parseGrpcFrame(buf: Buffer): Buffer {
  if (buf.length < 5) throw new Error(`gRPC response frame too short (${buf.length} bytes)`);
  const len = buf.readUInt32BE(1);
  if (buf.length < 5 + len) throw new Error(`gRPC response frame truncated: expected ${5 + len} bytes, got ${buf.length}`);
  return buf.subarray(5, 5 + len);
}

/** Parse gRPC-web inline trailers from the response body.
 *  After the data frame(s), gRPC-web appends a trailer frame with flag 0x80. */
function parseGrpcWebTrailers(buf: Buffer): Record<string, string> {
  // Skip past data frames to find the trailer frame (flag byte 0x80)
  let offset = 0;
  while (offset < buf.length) {
    if (offset + 5 > buf.length) {
      throw new Error(`Truncated gRPC-web frame at offset ${offset}: need 5 header bytes but only ${buf.length - offset} remain`);
    }
    const flag = buf[offset];
    const len = buf.readUInt32BE(offset + 1);
    if (flag === 0x80) {
      // Trailer frame: ASCII key-value pairs separated by \r\n
      const text = buf.subarray(offset + 5, offset + 5 + len).toString("utf8");
      const trailers: Record<string, string> = {};
      for (const line of text.split("\r\n")) {
        const colon = line.indexOf(":");
        if (colon > 0) {
          trailers[line.slice(0, colon).trim()] = line.slice(colon + 1).trim();
        }
      }
      return trailers;
    }
    offset += 5 + len;
  }
  return {};
}

// ---------------------------------------------------------------------------
// Soft-fail wrapper — allows new/experimental permutations to warn on failure
// without breaking CI. Tests marked knownWorking: false log a warning instead
// of hard-failing.
// ---------------------------------------------------------------------------

function softTest(
  label: string,
  knownWorking: boolean,
  testFn: (args: { request: APIRequestContext; testInfo: TestInfo }) => Promise<void>,
) {
  if (knownWorking) {
    test(label, async ({ request }, testInfo) => {
      await testFn({ request, testInfo });
    });
  } else {
    test(label, async ({ request }, testInfo) => {
      try {
        await testFn({ request, testInfo });
      } catch (err) {
        const msg = err instanceof Error ? err.message : String(err);
        testInfo.annotations.push({
          type: "fixme",
          description: `[soft-fail] ${msg}`,
        });
        // eslint-disable-next-line no-console
        console.warn(`⚠ SOFT-FAIL: ${label}\n  ${msg}`);
      }
    });
  }
}

// ---------------------------------------------------------------------------
// Protocol × transport permutation types
// ---------------------------------------------------------------------------

type Transport = "http1" | "h2c";
type Protocol = "rest" | "grpc-web" | "grpc" | "connectrpc";
type Encoding = "proto" | "json";

type ApiPermutation = {
  protocol: Protocol;
  encoding: Encoding;
  transport: Transport;
  /** Whether this permutation is already known to work (hard-fail on regression). */
  knownWorking: boolean;
};

// ---------------------------------------------------------------------------
// V1 API permutations — AdminService/Healthz (unauthenticated, expects 200/OK)
// ---------------------------------------------------------------------------

const V1_GRPC_PATH = "/zitadel.admin.v1.AdminService/Healthz";
const V1_REST_PATH = "/admin/v1/healthz";

const v1Permutations: ApiPermutation[] = [
  // REST
  { protocol: "rest",     encoding: "json",  transport: "http1", knownWorking: true },
  { protocol: "rest",     encoding: "json",  transport: "h2c",   knownWorking: true },
  // gRPC-web (HTTP/1.1 only — grpc-web is designed for browser HTTP/1.1)
  { protocol: "grpc-web", encoding: "proto", transport: "http1", knownWorking: true },
  { protocol: "grpc-web", encoding: "json",  transport: "http1", knownWorking: true },
  // gRPC native (h2c only — requires HTTP/2)
  { protocol: "grpc",     encoding: "proto", transport: "h2c",   knownWorking: true },
  { protocol: "grpc",     encoding: "json",  transport: "h2c",   knownWorking: true },
];

// ---------------------------------------------------------------------------
// V2 API permutations — SessionService/ListSessions (requires auth → 401)
// ---------------------------------------------------------------------------

const V2_GRPC_PATH = "/zitadel.session.v2.SessionService/ListSessions";
const V2_REST_PATH = "/v2/sessions";

const v2Permutations: ApiPermutation[] = [
  // REST (gRPC-gateway)
  { protocol: "rest",      encoding: "json",  transport: "http1", knownWorking: true },
  { protocol: "rest",      encoding: "json",  transport: "h2c",   knownWorking: true },
  // ConnectRPC (v2 only — v1 does not register Connect handlers)
  { protocol: "connectrpc", encoding: "proto", transport: "http1", knownWorking: true },
  { protocol: "connectrpc", encoding: "json",  transport: "http1", knownWorking: true },
  { protocol: "connectrpc", encoding: "proto", transport: "h2c",   knownWorking: true },
  { protocol: "connectrpc", encoding: "json",  transport: "h2c",   knownWorking: true },
  // gRPC native (h2c only)
  { protocol: "grpc",       encoding: "proto", transport: "h2c",   knownWorking: true },
  { protocol: "grpc",       encoding: "json",  transport: "h2c",   knownWorking: true },
];

// ---------------------------------------------------------------------------
// Tests — service-level reachability (non-API)
// ---------------------------------------------------------------------------

// The base URL comes from playwright.config.ts (PROXY_HTTP_PUBLISHED_PORT, default 8888).
// We resolve it here so h2Post()/h2Get() can connect to the same host.
const BASE =
  process.env.PLAYWRIGHT_BASE_URL ??
  `http://localhost:${process.env.PROXY_HTTP_PUBLISHED_PORT ?? "8888"}`;

test.describe("traefik routing", () => {
  test("GET / → redirects to /ui/v2/login/", async ({ request }) => {
    // Playwright follows redirects by default; verify the final destination.
    // Note: Next.js normalises /ui/v2/login/ → /ui/v2/login (strips trailing slash).
    const resp = await request.get("/");
    expect(resp.url()).toContain("/ui/v2/login");
    expect(resp.status()).toBe(200);
  });
});

test.describe("login ui", () => {
  test("GET /ui/v2/login/healthy → 200", async ({ request }) => {
    const resp = await request.get("/ui/v2/login/healthy");
    expect(resp.status()).toBe(200);
  });

  test("GET /ui/v2/login/ → 200 HTML", async ({ request }) => {
    const resp = await request.get("/ui/v2/login/");
    expect(resp.status()).toBe(200);
    expect(await resp.text()).toContain("<html");
  });
});

test.describe("console", () => {
  test("GET /ui/console/ → 200 HTML (Angular SPA shell)", async ({ request }) => {
    const resp = await request.get("/ui/console/");
    expect(resp.status()).toBe(200);
    expect(await resp.text()).toContain("<html");
  });
});

test.describe("oidc", () => {
  test("GET /.well-known/openid-configuration → valid discovery document", async ({ request }) => {
    const resp = await request.get("/.well-known/openid-configuration");
    expect(resp.status()).toBe(200);
    const doc = await resp.json();
    expect(doc).toHaveProperty("issuer");
    expect(doc).toHaveProperty("authorization_endpoint");
    expect(doc).toHaveProperty("token_endpoint");
    expect(doc).toHaveProperty("jwks_uri");
  });

  test("GET /oauth/v2/keys → JWKS with keys array", async ({ request }) => {
    const resp = await request.get("/oauth/v2/keys");
    expect(resp.status()).toBe(200);
    const body = await resp.json();
    expect(body).toHaveProperty("keys");
    expect(Array.isArray(body.keys)).toBe(true);
  });
});

test.describe("saml", () => {
  test("GET /saml/v2/metadata → EntityDescriptor XML", async ({ request }) => {
    const resp = await request.get("/saml/v2/metadata");
    expect(resp.status()).toBe(200);
    expect(await resp.text()).toContain("EntityDescriptor");
  });
});

test.describe("api v1 rest extras", () => {
  // Additional REST tests beyond the matrix: all three v1 services + /api alias.
  for (const svc of ["management", "admin", "auth"]) {
    test(`GET /${svc}/v1/healthz → 200`, async ({ request }) => {
      const resp = await request.get(`/${svc}/v1/healthz`);
      expect(resp.status()).toBe(200);
    });
  }

  test("GET /api/admin/v1/healthz via /api alias → 200 (tests Traefik strip-prefix middleware)", async ({ request }) => {
    const resp = await request.get("/api/admin/v1/healthz");
    expect(resp.status()).toBe(200);
  });
});

// ---------------------------------------------------------------------------
// API V1 protocol × transport matrix
//
// Endpoint: AdminService/Healthz (unauthenticated)
// Expected: success (REST → 200, gRPC → grpc-status 0, gRPC-web → grpc-status 0)
// ---------------------------------------------------------------------------

test.describe("api v1 protocol matrix", () => {
  for (const p of v1Permutations) {
    const label = `v1 ${p.protocol} ${p.encoding} ${p.transport}`;

    softTest(label, p.knownWorking, async ({ request }) => {
      switch (p.protocol) {
        // ----- REST (gRPC-gateway) -----
        case "rest": {
          if (p.transport === "http1") {
            const resp = await request.get(V1_REST_PATH);
            expect(resp.status()).toBe(200);
          } else {
            const resp = await h2Get(BASE, V1_REST_PATH);
            expect(resp.status).toBe(200);
          }
          break;
        }

        // ----- gRPC-web (HTTP/1.1) -----
        case "grpc-web": {
          const contentType =
            p.encoding === "proto"
              ? "application/grpc-web+proto"
              : "application/grpc-web+json";
          // For json encoding, Healthz request is empty → JSON body is '{}'
          const reqBody =
            p.encoding === "proto"
              ? grpcFrame()
              : grpcFrame(Buffer.from("{}"));
          const resp = await request.post(V1_GRPC_PATH, {
            headers: {
              "content-type": contentType,
              "x-grpc-web": "1",
            },
            data: reqBody,
          });
          expect(resp.status()).toBe(200);
          // For proto encoding: parse inline trailers and assert grpc-status 0.
          // For json encoding: the Go gRPC server has no JSON codec registered and
          // returns no well-formed gRPC-web trailer frame. HTTP 200 alone proves
          // the request was routed to the backend by Traefik.
          if (p.encoding === "proto") {
            const buf = Buffer.from(await resp.body());
            const trailers = parseGrpcWebTrailers(buf);
            expect(trailers["grpc-status"]).toBe("0");
          }
          break;
        }

        // ----- gRPC native (h2c) -----
        case "grpc": {
          const contentType =
            p.encoding === "proto"
              ? "application/grpc+proto"
              : "application/grpc+json";
          const reqBody =
            p.encoding === "proto"
              ? grpcFrame()
              : grpcFrame(Buffer.from("{}"));
          const resp = await h2Post(BASE, V1_GRPC_PATH, {
            "content-type": contentType,
            te: "trailers",
          }, reqBody);
          expect(resp.status).toBe(200);
          // For proto encoding: assert OK (grpc-status 0).
          // For json encoding: the Go gRPC server has no JSON codec registered,
          // so it returns INTERNAL (13). We only assert the status is defined,
          // which proves the request was routed to the backend.
          if (p.encoding === "proto") {
            expect(resp.trailers["grpc-status"]).toBe("0");
          } else {
            expect(resp.trailers["grpc-status"]).toBeDefined();
          }
          break;
        }

        default:
          throw new Error(`unexpected protocol: ${p.protocol}`);
      }
    });
  }
});

// ---------------------------------------------------------------------------
// API V2 protocol × transport matrix
//
// Endpoint: SessionService/ListSessions (requires auth)
// Expected: 401 / UNAUTHENTICATED
//   - REST (gRPC-gateway): HTTP 401, body.code === 16
//   - ConnectRPC unary:    HTTP 401, body.code === "unauthenticated"
//   - gRPC native:         HTTP 200, grpc-status trailer === "16"
// ---------------------------------------------------------------------------

test.describe("api v2 protocol matrix", () => {
  for (const p of v2Permutations) {
    const label = `v2 ${p.protocol} ${p.encoding} ${p.transport}`;

    softTest(label, p.knownWorking, async ({ request }) => {
      switch (p.protocol) {
        // ----- REST (gRPC-gateway) -----
        case "rest": {
          if (p.transport === "http1") {
            const resp = await request.post(V2_REST_PATH, {
              headers: { "content-type": "application/json" },
              data: "{}",
            });
            expect(resp.status()).toBe(401);
            const body = await resp.json();
            expect(body.code).toBe(16); // UNAUTHENTICATED
          } else {
            const resp = await h2Post(BASE, V2_REST_PATH, {
              "content-type": "application/json",
            }, "{}");
            expect(resp.status).toBe(401);
            const body = JSON.parse(resp.body.toString("utf8"));
            expect(body.code).toBe(16);
          }
          break;
        }

        // ----- Connect protocol unary (v2 only — v1 does not register Connect handlers) -----
        // Per https://connectrpc.com/docs/multi-protocol/:
        //   application/json  → Connect unary JSON  (gRPC-style path)
        //   application/proto → Connect unary proto (gRPC-style path)
        // application/connect+json and application/connect+proto are STREAMING-only.
        // Connect unary uses no gRPC length-prefix framing; body is plain JSON or proto bytes.
        // Error response format: { "code": "unauthenticated", ... }
        case "connectrpc": {
          const contentType =
            p.encoding === "proto"
              ? "application/proto"
              : "application/json";
          // For connect+json, send a JSON body; for connect+proto, send empty bytes
          const reqBody = p.encoding === "json" ? "{}" : Buffer.alloc(0);

          if (p.transport === "http1") {
            const resp = await request.post(V2_GRPC_PATH, {
              headers: { "content-type": contentType },
              data: reqBody,
            });
            // ConnectRPC returns the HTTP status directly
            expect(resp.status()).toBe(401);
            const body = await resp.json();
            expect(body.code).toBe("unauthenticated");
          } else {
            const resp = await h2Post(BASE, V2_GRPC_PATH, {
              "content-type": contentType,
            }, reqBody);
            expect(resp.status).toBe(401);
            const body = JSON.parse(resp.body.toString("utf8"));
            expect(body.code).toBe("unauthenticated");
          }
          break;
        }

        // ----- gRPC native (h2c) -----
        case "grpc": {
          const contentType =
            p.encoding === "proto"
              ? "application/grpc+proto"
              : "application/grpc+json";
          const reqBody =
            p.encoding === "proto"
              ? grpcFrame()
              : grpcFrame(Buffer.from("{}"));
          const resp = await h2Post(BASE, V2_GRPC_PATH, {
            "content-type": contentType,
            te: "trailers",
          }, reqBody);
          // gRPC always returns HTTP 200; error is in grpc-status trailer
          expect(resp.status).toBe(200);
          expect(resp.trailers["grpc-status"]).toBe("16"); // UNAUTHENTICATED
          break;
        }

        default:
          throw new Error(`unexpected protocol: ${p.protocol}`);
      }
    });
  }
});

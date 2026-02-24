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
 *   - API v1 REST    (gRPC-Gateway, /management/v1 + /admin/v1 + /auth/v1)
 *   - gRPC v1 h2c              (Traefik catch-all rule, priority 100, grpc.health.v1.Health)
 *   - gRPC-web                  (Traefik catch-all rule, priority 100, grpc.health.v1.Health)
 *   - API v2 REST HTTP/1.1      (Traefik catch-all rule, priority 100, gRPC-gateway JSON at /v2/sessions)
 *   - API v2 REST HTTP/2 (h2c)  (Traefik catch-all rule, priority 100, gRPC-gateway JSON at /v2/sessions via h2c)
 *   - Root redirect            (Traefik replacepath middleware → /ui/v2/login/)
 *
 * For end-to-end browser login flow (proves Traefik → Login → API → DB chain):
 *   see smoke.spec.ts
 *
 * For feature-level login tests (MFA, OIDC flows, IdPs, etc.):
 *   see apps/login/acceptance/ — run via @zitadel/compose:test-login-acceptance
 */
import { expect, test } from "@playwright/test";
import http2 from "node:http2";
import { Buffer } from "node:buffer";

// ---------------------------------------------------------------------------
// HTTP/2 (h2c) helper — used for the gRPC h2c test
// ---------------------------------------------------------------------------

type H2Response = {
  status: number;
  headers: Record<string, string>;
  body: Buffer;
  trailers: Record<string, string>;
};

function h2Post(
  baseUrl: string,
  path: string,
  reqHeaders: Record<string, string>,
  body?: Buffer | string,
): Promise<H2Response> {
  return new Promise((resolve, reject) => {
    const client = http2.connect(baseUrl);
    const timeout = setTimeout(() => {
      client.destroy();
      reject(new Error(`h2c request to ${path} timed out`));
    }, 30_000);

    client.once("error", (err) => {
      clearTimeout(timeout);
      reject(err);
    });

    const req = client.request({
      ":method": "POST",
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

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// The base URL comes from playwright.config.ts (PROXY_HTTP_PUBLISHED_PORT, default 8080).
// We resolve it here so h2Post() can connect to the same host.
const BASE =
  process.env.PLAYWRIGHT_BASE_URL ??
  `http://localhost:${process.env.PROXY_HTTP_PUBLISHED_PORT ?? "8080"}`;

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

test.describe("api v1 rest (grpc-gateway)", () => {
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

test.describe("grpc v1 (h2c)", () => {
  // gRPC-web uses HTTP/1.1 with application/grpc-web+proto.
  // A 5-byte length-prefix frame is the same as gRPC, but trailers are
  // inlined as a DATA frame (flag byte 0x80). The health check requires no auth.
  test("grpc.health.v1.Health/Check via gRPC-web → SERVING", async ({ request }) => {
    const reqFrame = grpcFrame();
    const resp = await request.post("/grpc.health.v1.Health/Check", {
      headers: {
        "content-type": "application/grpc-web+proto",
        "x-grpc-web": "1",
      },
      data: reqFrame,
    });
    expect(resp.status()).toBe(200);
    const buf = Buffer.from(await resp.body());
    // First frame (0x00) = data
    const proto = parseGrpcFrame(buf);
    expect(proto[0]).toBe(0x08); // field 1, wire type 0
    expect(proto[1]).toBe(0x01); // ServingStatus.SERVING = 1
  });
});

test.describe("api v2 rest (grpc-gateway) http/1.1", () => {
  // ZITADEL's v2 APIs are served via the gRPC-gateway at /v2/... REST paths.
  // Connect-RPC wire format is not accessible through Traefik (the gRPC-gateway
  // catch-all intercepts /zitadel.*.*/Method paths).
  // An unauthenticated POST → 401 proves Traefik routes to zitadel-api and
  // the gRPC-gateway auth middleware is active.
  test("POST /v2/sessions → 401 unauthenticated", async ({ request }) => {
    const resp = await request.post("/v2/sessions", {
      headers: { "content-type": "application/json" },
      data: "{}",
    });
    expect(resp.status()).toBe(401);
    const body = await resp.json();
    // gRPC-gateway encodes UNAUTHENTICATED as numeric code 16
    expect(body.code).toBe(16);
  });
});

test.describe("api v2 rest (grpc-gateway) http/2 (h2c)", () => {
  // Same path, same routing — but over h2c (HTTP/2 cleartext).
  // Traefik accepts h2c and forwards to zitadel-api over h2c.
  // Verifies the full h2c path works for plain-HTTP v2 REST traffic.
  test("POST /v2/sessions over h2c → 401 unauthenticated", async () => {
    const resp = await h2Post(
      BASE,
      "/v2/sessions",
      { "content-type": "application/json" },
      "{}",
    );
    expect(resp.status).toBe(401);
    const body = JSON.parse(resp.body.toString("utf8"));
    // gRPC-gateway encodes UNAUTHENTICATED as numeric code 16
    expect(body.code).toBe(16);
  });
});

test.describe("grpc-web", () => {
  // Traefik routes via catch-all (priority 100) → zitadel-api h2c.
  // No dedicated gRPC router needed: the h2c backend handles gRPC natively.
  // grpc.health.v1.Health/Check requires no authentication.
  test("grpc.health.v1.Health/Check → SERVING", async () => {
    // HealthCheckRequest { service: "" } = empty proto → 0 message bytes
    const reqFrame = grpcFrame();

    const resp = await h2Post(
      BASE,
      "/grpc.health.v1.Health/Check",
      {
        "content-type": "application/grpc+proto",
        te: "trailers",
      },
      reqFrame,
    );

    // gRPC always uses HTTP 200; actual status is in the grpc-status trailer
    expect(resp.status).toBe(200);
    expect(resp.trailers["grpc-status"]).toBe("0"); // 0 = OK

    // HealthCheckResponse { status: SERVING (1) }
    // Proto encoding: field 1, varint 1 → [0x08, 0x01]
    const proto = parseGrpcFrame(resp.body);
    expect(proto.length).toBe(2);
    expect(proto[0]).toBe(0x08); // field 1, wire type 0 (varint)
    expect(proto[1]).toBe(0x01); // ServingStatus.SERVING = 1
  });
});

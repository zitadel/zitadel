// @vitest-environment node
import { NextRequest } from "next/server";
import { afterEach, describe, expect, test, vi } from "vitest";

vi.mock("./lib/service-url", () => ({
  getServiceConfig: () => ({ serviceConfig: { baseUrl: "https://idp.example.com" } }),
}));
vi.mock("./lib/server/security-settings", () => ({
  getIframeOrigins: vi.fn().mockResolvedValue(null),
}));

import { getIframeOrigins } from "./lib/server/security-settings";
import { proxy } from "./proxy";

afterEach(() => {
  vi.unstubAllEnvs();
  vi.mocked(getIframeOrigins).mockResolvedValue(null);
});

describe("strict Login CSP", () => {
  test.each(["/_not-found", "/_global-error"])("does not expose a prerendered internal document: %s", async (path) => {
    vi.stubEnv("CSP_NONCE_ENABLED", "true");
    const response = await proxy(new NextRequest(`https://login.example.com${path}`));
    expect(response.status).toBe(404);
    expect(await response.text()).toBe("");
    expect(response.headers.get("Cache-Control")).toBe("private, no-store");
  });

  test("overwrites caller headers with fresh 128-bit server nonces and one uncached rendering policy", async () => {
    vi.stubEnv("CSP_NONCE_ENABLED", "true");
    vi.stubEnv("CSP_FETCH_ENABLED", "false");
    const request = () =>
      new NextRequest("https://login.example.com/loginname", {
        headers: {
          "x-zitadel-csp-nonce": "caller-selected",
          "Content-Security-Policy": "script-src 'nonce-caller-selected'",
          "Content-Security-Policy-Report-Only": "script-src 'nonce-caller-selected'",
        },
      });
    const first = await proxy(request());
    const second = await proxy(request());
    const nonce = first.headers.get("x-middleware-request-x-zitadel-csp-nonce")!;
    expect(Buffer.from(nonce, "base64")).toHaveLength(16);
    expect(nonce).not.toBe("caller-selected");
    expect(second.headers.get("x-middleware-request-x-zitadel-csp-nonce")).not.toBe(nonce);
    const policy = first.headers.get("Content-Security-Policy")!;
    expect(policy).toContain(`script-src 'nonce-${nonce}' 'strict-dynamic'`);
    expect(first.headers.get("x-middleware-request-content-security-policy")).toBe(policy);
    expect(first.headers.get("x-middleware-override-headers")).not.toContain("content-security-policy-report-only");
    expect(first.headers.get("Cache-Control")).toBe("private, no-store");
  });

  test("strict policy also survives settings lookup failure", async () => {
    vi.stubEnv("CSP_NONCE_ENABLED", "true");
    vi.stubEnv("CSP_FETCH_ENABLED", "true");
    vi.mocked(getIframeOrigins).mockRejectedValue(new Error("fixture unavailable"));
    const response = await proxy(new NextRequest("https://login.example.com/password"));
    expect(response.headers.get("Content-Security-Policy")).toContain("'strict-dynamic'");
    expect(response.headers.get("X-Frame-Options")).toBe("deny");
  });

  test("legacy mode stays off by default and strips a spoofed theme nonce", async () => {
    vi.stubEnv("CSP_NONCE_ENABLED", "false");
    vi.stubEnv("CSP_FETCH_ENABLED", "false");
    const response = await proxy(
      new NextRequest("https://login.example.com/loginname", {
        headers: { "x-zitadel-csp-nonce": "caller-selected" },
      }),
    );
    expect(response.headers.get("Content-Security-Policy")).toContain("script-src 'self' 'unsafe-inline' 'unsafe-eval'");
    expect(response.headers.get("x-middleware-request-x-zitadel-csp-nonce")).toBeNull();
  });
});

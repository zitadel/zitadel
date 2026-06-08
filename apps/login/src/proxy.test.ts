import { describe, expect, test } from "vitest";

import { isProxyPath, normalizeProxyPathname, PROXY_MATCHER } from "./proxy";

describe("proxy path normalization", () => {
  test("normalizes callback paths under the login base path", () => {
    expect(normalizeProxyPathname("/ui/v2/login/idps/callback", "/ui/v2/login")).toBe("/idps/callback");
    expect(normalizeProxyPathname("/ui/v2/login/oidc/v1/userinfo", "/ui/v2/login")).toBe("/oidc/v1/userinfo");
  });

  test("leaves non-prefixed paths unchanged", () => {
    expect(normalizeProxyPathname("/idps/callback", "/ui/v2/login")).toBe("/idps/callback");
    expect(normalizeProxyPathname("/loginname", "/ui/v2/login")).toBe("/loginname");
  });

  test("matches proxy callback paths with and without the base path", () => {
    expect(isProxyPath("/idps/callback", "/ui/v2/login")).toBe(true);
    expect(isProxyPath("/ui/v2/login/idps/callback", "/ui/v2/login")).toBe(true);
    expect(isProxyPath("/ui/v2/login/oidc/v1/userinfo", "/ui/v2/login")).toBe(true);
    expect(isProxyPath("/ui/v2/login/loginname", "/ui/v2/login")).toBe(false);
  });

  test("declares the bare callback route in the middleware matcher", () => {
    expect(PROXY_MATCHER).toContain("/idps/callback");
  });
});

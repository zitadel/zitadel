import { describe, expect, test } from "vitest";
import { buildCSP } from "./csp";

describe("buildCSP", () => {
  test("nonce policy admits only trusted scripts and preserves supplied styles and iframe settings", () => {
    const nonce = "AAECAwQFBgcICQoLDA0ODw==";
    const csp = buildCSP({ nonce, iframeOrigins: ["https://portal.example.com"] });
    const scripts = csp.split("; ").find((directive) => directive.startsWith("script-src "));
    expect(scripts).toBe(`script-src 'nonce-${nonce}' 'strict-dynamic'`);
    expect(csp).toContain("script-src-attr 'none'");
    expect(csp).toContain("base-uri 'none'");
    expect(csp).toContain("style-src 'self' 'unsafe-inline'");
    expect(csp).toContain("frame-ancestors https://portal.example.com");
  });

  test.each(["caller-nonce", "'unsafe-inline'", "AAECAwQFBgcICQoLDA0ODw==; script-src *"])(
    "rejects an invalid nonce: %s",
    (nonce) => expect(() => buildCSP({ nonce })).toThrow("Invalid CSP nonce"),
  );

  test("returns all base directives with safe defaults", () => {
    const csp = buildCSP();

    expect(csp).toContain("default-src 'self'");
    expect(csp).toContain("script-src 'self' 'unsafe-inline' 'unsafe-eval'");
    expect(csp).toContain("connect-src 'self'");
    expect(csp).toContain("style-src 'self' 'unsafe-inline'");
    expect(csp).toContain("font-src 'self'");
    expect(csp).toContain("img-src 'self'");
    expect(csp).toContain("object-src 'none'");
    expect(csp).toContain("frame-ancestors 'none'");
  });

  test("adds serviceUrl to img-src and font-src", () => {
    const csp = buildCSP({ serviceUrl: "https://my-instance.zitadel.cloud" });

    expect(csp).toContain("img-src 'self' https://my-instance.zitadel.cloud");
    expect(csp).toContain("font-src 'self' https://my-instance.zitadel.cloud");
  });

  test("keeps frame-ancestors as 'none' when iframeOrigins is empty", () => {
    const csp = buildCSP({ iframeOrigins: [] });

    expect(csp).toContain("frame-ancestors 'none'");
  });

  test("overrides frame-ancestors when iframeOrigins are provided", () => {
    const csp = buildCSP({
      iframeOrigins: ["https://app.example.com", "https://other.example.com"],
    });

    expect(csp).toContain("frame-ancestors https://app.example.com https://other.example.com");
    expect(csp).not.toContain("frame-ancestors 'none'");
  });

  test("combines serviceUrl and iframeOrigins", () => {
    const csp = buildCSP({
      serviceUrl: "https://zitadel.mycompany.com",
      iframeOrigins: ["https://portal.mycompany.com"],
    });

    expect(csp).toContain("img-src 'self' https://zitadel.mycompany.com");
    expect(csp).toContain("font-src 'self' https://zitadel.mycompany.com");
    expect(csp).toContain("frame-ancestors https://portal.mycompany.com");
    expect(csp).not.toContain("frame-ancestors 'none'");
  });

  test("base directives are preserved when options are set", () => {
    const csp = buildCSP({
      serviceUrl: "https://example.com",
      iframeOrigins: ["https://embed.example.com"],
    });

    expect(csp).toContain("default-src 'self'");
    expect(csp).toContain("script-src 'self' 'unsafe-inline' 'unsafe-eval'");
    expect(csp).toContain("connect-src 'self'");
    expect(csp).toContain("style-src 'self' 'unsafe-inline'");
    expect(csp).toContain("object-src 'none'");
  });
});

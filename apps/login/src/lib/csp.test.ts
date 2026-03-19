import { describe, expect, test } from "vitest";
import { buildCSP, resolveImageSources } from "./csp";

describe("buildCSP", () => {
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

  test("adds serviceUrl to img-src only", () => {
    const csp = buildCSP({ serviceUrl: "https://my-instance.zitadel.cloud" });

    expect(csp).toContain("img-src 'self' https://my-instance.zitadel.cloud");
    expect(csp).toMatch(/font-src 'self'(?:;| |$)/);
  });

  test("adds image sources from custom hosts and normalizes to origins", () => {
    const csp = buildCSP({
      imageSources: ["idp.zitadel.com", "https://idp.zitadel.com/assets/v1/logo", "localhost:8080"],
    });

    expect(csp).toContain("img-src 'self' https://idp.zitadel.com http://localhost:8080");
  });

  test("adds both host:port and host-only origins when service URL has a custom port", () => {
    const csp = buildCSP({
      serviceUrl: "http://zitadel.myorg.local:8080/ui/v2/login",
    });

    expect(csp).toContain("img-src 'self' http://zitadel.myorg.local:8080 http://zitadel.myorg.local");
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

  test("resolveImageSources keeps custom hosts and serviceUrl", () => {
    const sources = resolveImageSources({
      serviceUrl: "http://zitadel.myorg.local:8080",
      publicHost: "login.myorg.local:3021",
      customRequestHeaders: "x-zitadel-public-host:zitadel.myorg.local,x-forwarded-proto:http",
    });

    expect(sources).toEqual([
      "zitadel.myorg.local",
      "login.myorg.local:3021",
      "http://zitadel.myorg.local:8080",
    ]);
  });
});

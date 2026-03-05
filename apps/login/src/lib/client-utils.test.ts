import { describe, expect, test, vi, beforeEach } from "vitest";
import { isSafeRedirectUri, sanitizeRedirectUri } from "./client-utils";

describe("isSafeRedirectUri", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("should reject empty uri", async () => {
    expect(await isSafeRedirectUri("")).toBe(false);
    expect(await isSafeRedirectUri(undefined as unknown as string)).toBe(false);
  });

  test("should allow relative paths", async () => {
    expect(await isSafeRedirectUri("/dashboard")).toBe(true);
    expect(await isSafeRedirectUri("/ui/console")).toBe(true);
    expect(await isSafeRedirectUri("/")).toBe(true);
  });

  test("should reject protocol-relative paths", async () => {
    expect(await isSafeRedirectUri("//evil.com")).toBe(false);
    expect(await isSafeRedirectUri("\\\\evil.com")).toBe(false);
  });

  test("should allow absolute URLs matching the current host", async () => {
    expect(await isSafeRedirectUri("https://my-zitadel.com/console")).toBe(true);
  });

  test("should allow absolute URLs to external domains", async () => {
    expect(await isSafeRedirectUri("https://evil.com")).toBe(true);
    expect(await isSafeRedirectUri("https://settings.com/dashboard")).toBe(true);
  });

  test("should reject non-http(s) protocols", async () => {
    expect(await isSafeRedirectUri("javascript:alert(1)")).toBe(false);
    expect(await isSafeRedirectUri("data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==")).toBe(false);
    expect(await isSafeRedirectUri("file:///etc/passwd")).toBe(false);
  });

  test("should gracefully reject invalid, unparsable URLs that aren't relative", async () => {
    expect(await isSafeRedirectUri("not-a-valid-url")).toBe(false);
    expect(await isSafeRedirectUri("http://%")).toBe(false);
  });
});

describe("sanitizeRedirectUri", () => {
  test("should return undefined for empty or falsy uri", () => {
    expect(sanitizeRedirectUri("")).toBeUndefined();
    expect(sanitizeRedirectUri(undefined as unknown as string)).toBeUndefined();
  });

  test("should reconstruct relative paths", () => {
    expect(sanitizeRedirectUri("/dashboard")).toBe("/dashboard");
    expect(sanitizeRedirectUri("/ui/console")).toBe("/ui/console");
    expect(sanitizeRedirectUri("/")).toBe("/");
    expect(sanitizeRedirectUri("/signedin?loginName=foo")).toBe("/signedin?loginName=foo");
    expect(sanitizeRedirectUri("/path?a=1#hash")).toBe("/path?a=1#hash");
  });

  test("should reject protocol-relative paths", () => {
    expect(sanitizeRedirectUri("//evil.com")).toBeUndefined();
  });

  test("should reject absolute URLs without allowedOrigin on server", () => {
    expect(sanitizeRedirectUri("https://example.com")).toBeUndefined();
    expect(sanitizeRedirectUri("https://example.com/path")).toBeUndefined();
  });

  test("should allow same-origin absolute URLs with allowedOrigin", () => {
    expect(sanitizeRedirectUri("https://my-host.com/dashboard", "https://my-host.com")).toBe("https://my-host.com/dashboard");
    expect(sanitizeRedirectUri("https://my-host.com/", "https://my-host.com")).toBe("https://my-host.com/");
    expect(sanitizeRedirectUri("https://my-host.com/path?q=1#frag", "https://my-host.com")).toBe("https://my-host.com/path?q=1#frag");
  });

  test("should reject cross-origin absolute URLs even with allowedOrigin", () => {
    expect(sanitizeRedirectUri("https://evil.com", "https://my-host.com")).toBeUndefined();
    expect(sanitizeRedirectUri("https://evil.com/path", "https://my-host.com")).toBeUndefined();
  });

  test("should reject URLs with embedded credentials (userinfo)", () => {
    expect(sanitizeRedirectUri("https://user:pass@example.com/path", "https://example.com")).toBeUndefined();
    expect(sanitizeRedirectUri("https://user@example.com", "https://example.com")).toBeUndefined();
    expect(sanitizeRedirectUri("https://user:pass@example.com", undefined, true)).toBeUndefined();
  });

  test("should allow any http(s) URL when trustOrigin is true", () => {
    expect(sanitizeRedirectUri("https://external.com/callback", undefined, true)).toBe("https://external.com/callback");
    expect(sanitizeRedirectUri("https://other.com/path?q=1", undefined, true)).toBe("https://other.com/path?q=1");
    expect(sanitizeRedirectUri("http://dev.local:8080/", undefined, true)).toBe("http://dev.local:8080/");
  });

  test("should reject non-http protocols even when trustOrigin is true", () => {
    expect(sanitizeRedirectUri("javascript:alert(1)", undefined, true)).toBeUndefined();
    expect(sanitizeRedirectUri("data:text/html,test", undefined, true)).toBeUndefined();
  });

  test("should return undefined for non-http(s) protocols", () => {
    expect(sanitizeRedirectUri("javascript:alert(1)")).toBeUndefined();
    expect(sanitizeRedirectUri("data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==")).toBeUndefined();
    expect(sanitizeRedirectUri("file:///etc/passwd")).toBeUndefined();
  });

  test("should return undefined for invalid, unparsable URLs", () => {
    expect(sanitizeRedirectUri("not-a-valid-url")).toBeUndefined();
    expect(sanitizeRedirectUri("http://%")).toBeUndefined();
  });
});

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

  test("should return relative paths as-is", () => {
    expect(sanitizeRedirectUri("/dashboard")).toBe("/dashboard");
    expect(sanitizeRedirectUri("/ui/console")).toBe("/ui/console");
    expect(sanitizeRedirectUri("/")).toBe("/");
    expect(sanitizeRedirectUri("/signedin?loginName=foo")).toBe("/signedin?loginName=foo");
  });

  test("should reject protocol-relative paths", () => {
    expect(sanitizeRedirectUri("//evil.com")).toBeUndefined();
  });

  test("should return reconstructed href for safe absolute URLs", () => {
    expect(sanitizeRedirectUri("https://my-zitadel.com/console")).toBe("https://my-zitadel.com/console");
    expect(sanitizeRedirectUri("http://localhost:8080/dashboard")).toBe("http://localhost:8080/dashboard");
    expect(sanitizeRedirectUri("https://example.com/path?q=1#frag")).toBe("https://example.com/path?q=1#frag");
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

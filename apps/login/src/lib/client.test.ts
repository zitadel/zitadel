import { describe, expect, test, vi, beforeEach } from "vitest";
import { getNextUrl, isSafeRedirectUri } from "./client";

// Mock next/headers
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

// Mock host helper
vi.mock("./server/host", () => ({
  getPublicHostWithProtocol: vi.fn(),
}));

describe("getNextUrl", () => {
  const command = { loginName: "test-user" };

  beforeEach(() => {
    vi.clearAllMocks();
    delete (process.env as any).DEFAULT_REDIRECT_URI;
    delete (process.env as any).NEXT_PUBLIC_BASE_PATH;
  });

  test("should use DEFAULT_REDIRECT_URI if set", async () => {
    process.env.DEFAULT_REDIRECT_URI = "https://env-override.com";
    const result = await getNextUrl(command);
    expect(result).toBe("https://env-override.com");
  });

  test("should use host-based redirect if DEFAULT_REDIRECT_URI is set to a path (starting with '/')", async () => {
    const { headers } = await import("next/headers");
    const { getPublicHostWithProtocol } = await import("./server/host");

    process.env.DEFAULT_REDIRECT_URI = "/dashboard";
    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getPublicHostWithProtocol).mockReturnValue("https://my-host.com");
    process.env.NEXT_PUBLIC_BASE_PATH = "/ui/v2/login";

    const result = await getNextUrl(command);
    expect(result).toBe("https://my-host.com/dashboard");
  });

  test("should use defaultRedirectUri if env is NOT set", async () => {
    const result = await getNextUrl(command, "https://settings.com");
    expect(result).toBe("https://settings.com");
  });

  test("should fallback to relative signedin page if everything else fails (the new default)", async () => {
    const { headers } = await import("next/headers");
    vi.mocked(headers).mockRejectedValue(new Error("No headers"));

    const result = await getNextUrl(command);
    expect(result).toBe("/signedin?loginName=test-user");
  });
});

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

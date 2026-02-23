import { describe, expect, test, vi, beforeEach } from "vitest";
import { getConfiguredRedirectUri, getNextUrl } from "./client";

// Mock next/headers
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

// Mock host helper
vi.mock("./server/host", () => ({
  getPublicHostWithProtocol: vi.fn(),
}));

describe("getConfiguredRedirectUri", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    delete (process.env as any).DEFAULT_REDIRECT_URI;
  });

  test("should return DEFAULT_REDIRECT_URI if set to absolute URL", async () => {
    process.env.DEFAULT_REDIRECT_URI = "https://env-override.com";
    const result = await getConfiguredRedirectUri();
    expect(result).toBe("https://env-override.com");
  });

  test("should return host-based redirect if DEFAULT_REDIRECT_URI is a path", async () => {
    const { headers } = await import("next/headers");
    const { getPublicHostWithProtocol } = await import("./server/host");

    process.env.DEFAULT_REDIRECT_URI = "/dashboard";
    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getPublicHostWithProtocol).mockReturnValue("https://my-host.com");

    const result = await getConfiguredRedirectUri();
    expect(result).toBe("https://my-host.com/dashboard");
  });

  test("should return defaultRedirectUri from settings if env is not set", async () => {
    const result = await getConfiguredRedirectUri("https://settings.com");
    expect(result).toBe("https://settings.com");
  });

  test("should return null if DEFAULT_REDIRECT_URI is a path but host cannot be determined", async () => {
    const { headers } = await import("next/headers");
    process.env.DEFAULT_REDIRECT_URI = "/dashboard";
    vi.mocked(headers).mockRejectedValue(new Error("No headers"));

    const result = await getConfiguredRedirectUri();
    expect(result).toBeNull();
  });

  test("should return null if neither env nor settings provide a redirect", async () => {
    const { headers } = await import("next/headers");
    vi.mocked(headers).mockRejectedValue(new Error("No headers"));

    const result = await getConfiguredRedirectUri();
    expect(result).toBeNull();
  });

  test("should return null if defaultRedirectUri is unsafe", async () => {
    const result = await getConfiguredRedirectUri("javascript:alert(1)");
    expect(result).toBeNull();
  });
});

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

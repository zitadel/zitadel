import { describe, expect, test, vi, beforeEach } from "vitest";
import { getNextUrl } from "./client";

// Mock next/headers
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

// Mock host helper
vi.mock("./server/host", () => ({
  getPublicHostWithProtocol: vi.fn(),
}));

// Mock auth-flow to prevent transitive winston import via logger.ts
vi.mock("./server/auth-flow", () => ({
  completeAuthFlow: vi.fn(),
}));

describe("getNextUrl", () => {
  const command = { loginName: "test-user" };

  beforeEach(() => {
    vi.clearAllMocks();
    delete (process.env as any).DEFAULT_REDIRECT_URI;
    delete (process.env as any).NEXT_PUBLIC_BASE_PATH;
  });

  test("should use same-origin DEFAULT_REDIRECT_URI", async () => {
    const { headers } = await import("next/headers");
    const { getPublicHostWithProtocol } = await import("./server/host");

    process.env.DEFAULT_REDIRECT_URI = "https://my-host.com/callback";
    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getPublicHostWithProtocol).mockReturnValue("https://my-host.com");

    const result = await getNextUrl(command);
    expect(result).toBe("https://my-host.com/callback");
  });

  test("should use host-based redirect if DEFAULT_REDIRECT_URI is set to a path (starting with '/')", async () => {
    const { headers } = await import("next/headers");
    const { getPublicHostWithProtocol } = await import("./server/host");

    process.env.DEFAULT_REDIRECT_URI = "/dashboard";
    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getPublicHostWithProtocol).mockReturnValue("https://my-host.com");

    const result = await getNextUrl(command);
    expect(result).toBe("https://my-host.com/dashboard");
  });

  test("should use same-origin defaultRedirectUri", async () => {
    const { headers } = await import("next/headers");
    const { getPublicHostWithProtocol } = await import("./server/host");

    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getPublicHostWithProtocol).mockReturnValue("https://my-host.com");

    const result = await getNextUrl(command, "https://my-host.com/console");
    expect(result).toBe("https://my-host.com/console");
  });

  test("should use relative defaultRedirectUri", async () => {
    const { headers } = await import("next/headers");
    vi.mocked(headers).mockRejectedValue(new Error("No headers"));

    const result = await getNextUrl(command, "/console");
    expect(result).toBe("/console");
  });

  test("should allow cross-origin DEFAULT_REDIRECT_URI (admin-controlled env var)", async () => {
    const { headers } = await import("next/headers");
    const { getPublicHostWithProtocol } = await import("./server/host");

    process.env.DEFAULT_REDIRECT_URI = "https://external.com/callback";
    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getPublicHostWithProtocol).mockReturnValue("https://my-host.com");

    const result = await getNextUrl(command);
    expect(result).toBe("https://external.com/callback");
  });

  test("should reject cross-origin defaultRedirectUri and fall back", async () => {
    const { headers } = await import("next/headers");
    const { getPublicHostWithProtocol } = await import("./server/host");

    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getPublicHostWithProtocol).mockReturnValue("https://my-host.com");

    const result = await getNextUrl(command, "https://external.com/path");
    expect(result).toBe("/signedin?loginName=test-user");
  });

  test("should fallback to relative signedin page if everything else fails", async () => {
    const { headers } = await import("next/headers");
    vi.mocked(headers).mockRejectedValue(new Error("No headers"));

    const result = await getNextUrl(command);
    expect(result).toBe("/signedin?loginName=test-user");
  });

  test("should reject unsafe DEFAULT_REDIRECT_URI with javascript: protocol", async () => {
    process.env.DEFAULT_REDIRECT_URI = "javascript:alert(1)";
    const result = await getNextUrl(command);
    expect(result).toBe("/signedin?loginName=test-user");
  });

  test("should reject unsafe defaultRedirectUri with data: protocol", async () => {
    const result = await getNextUrl(command, "data:text/html,<script>alert(1)</script>");
    expect(result).toBe("/signedin?loginName=test-user");
  });
});



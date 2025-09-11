import { describe, expect, test, vi, beforeEach, afterEach } from "vitest";
import { getOriginalHost, getOriginalHostWithProtocol } from "./host";

// Mock the Next.js headers function
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

describe("Host utility functions", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("getOriginalHost", () => {
    test("should return x-forwarded-host when available", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-forwarded-host") return "zitadel.com";
          if (key === "x-original-host") return "backup.com";
          if (key === "host") return "internal.vercel.app";
          return null;
        }),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHost();
      expect(result).toBe("zitadel.com");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-forwarded-host");
    });

    test("should fall back to x-original-host when x-forwarded-host is not available", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-forwarded-host") return null;
          if (key === "x-original-host") return "original.com";
          if (key === "host") return "internal.vercel.app";
          return null;
        }),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHost();
      expect(result).toBe("original.com");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-forwarded-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-original-host");
    });

    test("should fall back to host when forwarded headers are not available", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-forwarded-host") return null;
          if (key === "x-original-host") return null;
          if (key === "host") return "fallback.com";
          return null;
        }),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHost();
      expect(result).toBe("fallback.com");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-forwarded-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-original-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("host");
    });

    test("should throw error when no host is found", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => null),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      await expect(getOriginalHost()).rejects.toThrow("No host found in headers");
    });

    test("should throw error when host is empty string", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => ""),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      await expect(getOriginalHost()).rejects.toThrow("No host found in headers");
    });

    test("should throw error when host is not a string", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => 123),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      await expect(getOriginalHost()).rejects.toThrow("No host found in headers");
    });
  });

  describe("getOriginalHostWithProtocol", () => {
    test("should return https for production domain", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => "zitadel.com"),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("https://zitadel.com");
    });

    test("should return http for localhost", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => "localhost:3000"),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("http://localhost:3000");
    });

    test("should return http for localhost without port", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => "localhost"),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("http://localhost");
    });

    test("should return https for custom domain", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => "auth.company.com"),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("https://auth.company.com");
    });
  });

  describe("Real-world scenarios", () => {
    test("should handle Vercel rewrite scenario", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn((key: string) => {
          // Simulate Vercel rewrite: zitadel.com/login -> login-zitadel-qa.vercel.app
          if (key === "x-forwarded-host") return "zitadel.com";
          if (key === "host") return "login-zitadel-qa.vercel.app";
          return null;
        }),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("https://zitadel.com");
    });

    test("should handle CloudFlare proxy scenario", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-forwarded-host") return "auth.company.com";
          if (key === "x-original-host") return null;
          if (key === "host") return "cloudflare-worker.workers.dev";
          return null;
        }),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHost();
      expect(result).toBe("auth.company.com");
    });

    test("should handle development environment", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "host") return "localhost:3000";
          return null;
        }),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("http://localhost:3000");
    });

    test("should handle staging environment with subdomain", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-forwarded-host") return "staging-auth.company.com";
          if (key === "host") return "staging-internal.vercel.app";
          return null;
        }),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("https://staging-auth.company.com");
    });
  });

  describe("Edge cases", () => {
    test("should handle IPv4 addresses", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => "192.168.1.100:3000"),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("https://192.168.1.100:3000");
    });

    test("should handle IPv6 addresses", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => "[::1]:3000"),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("https://[::1]:3000");
    });

    test("should handle hosts with ports", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => "zitadel.com:8080"),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("https://zitadel.com:8080");
    });

    test("should handle localhost with different ports", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn(() => "localhost:8080"),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHostWithProtocol();
      expect(result).toBe("http://localhost:8080");
    });

    test("should handle priority order correctly", async () => {
      const { headers } = await import("next/headers");
      const mockHeaders = {
        get: vi.fn((key: string) => {
          // All headers are present, should return x-forwarded-host (highest priority)
          if (key === "x-forwarded-host") return "priority1.com";
          if (key === "x-original-host") return "priority2.com";
          if (key === "host") return "priority3.com";
          return null;
        }),
      };

      vi.mocked(headers).mockResolvedValue(mockHeaders as any);

      const result = await getOriginalHost();
      expect(result).toBe("priority1.com");
      // Should only call x-forwarded-host since it's available
      expect(mockHeaders.get).toHaveBeenCalledWith("x-forwarded-host");
      expect(mockHeaders.get).toHaveBeenCalledTimes(1);
    });
  });
});

import { describe, expect, test, vi, beforeEach, afterEach } from "vitest";
import { getInstanceHost, getPublicHostWithProtocol, getPublicHost } from "./host";

describe("Host utility functions", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("getInstanceHost", () => {
    test("should use x-zitadel-instance-host when available", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-instance-host") return "instance.zitadel.cloud";
          if (key === "x-zitadel-forward-host") return "forward.zitadel.cloud";
          return null;
        }),
      } as any;

      const result = getInstanceHost(mockHeaders);
      expect(result).toBe("instance.zitadel.cloud");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-instance-host");
    });

    test("should use x-zitadel-forward-host when x-zitadel-instance-host is not available", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-instance-host") return null;
          if (key === "x-zitadel-forward-host") return "forward.zitadel.cloud";
          return null;
        }),
      } as any;

      const result = getInstanceHost(mockHeaders);
      expect(result).toBe("forward.zitadel.cloud");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-instance-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-forward-host");
    });

    test("should return null when neither x-zitadel-instance-host nor x-zitadel-forward-host are available", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-instance-host") return null;
          if (key === "x-zitadel-forward-host") return null;
          if (key === "x-forwarded-host") return "accounts.mycompany.com";
          if (key === "host") return "internal.server";
          return null;
        }),
      } as any;

      const result = getInstanceHost(mockHeaders);
      expect(result).toBeNull();
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-instance-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-forward-host");
    });
  });

  describe("getPublicHostWithProtocol", () => {
    test("should return https for production domain", () => {
      const mockHeaders = {
        get: vi.fn(() => "zitadel.com"),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("https://zitadel.com");
    });

    test("should return http for localhost", () => {
      const mockHeaders = {
        get: vi.fn(() => "localhost:3000"),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("http://localhost:3000");
    });

    test("should return http for localhost without port", () => {
      const mockHeaders = {
        get: vi.fn(() => "localhost"),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("http://localhost");
    });

    test("should return https for custom domain", () => {
      const mockHeaders = {
        get: vi.fn(() => "auth.company.com"),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("https://auth.company.com");
    });
  });

  describe("Real-world scenarios", () => {
    test("should handle Vercel rewrite scenario", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          // Simulate Vercel rewrite: zitadel.com/login -> login-zitadel-qa.vercel.app
          if (key === "x-forwarded-host") return "zitadel.com";
          if (key === "host") return "login-zitadel-qa.vercel.app";
          return null;
        }),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("https://zitadel.com");
    });

    test("should handle CloudFlare proxy scenario", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-forwarded-host") return "auth.company.com";
          if (key === "x-original-host") return null;
          if (key === "host") return "cloudflare-worker.workers.dev";
          return null;
        }),
      } as any;

      const result = getPublicHost(mockHeaders);
      expect(result).toBe("auth.company.com");
    });

    test("should handle development environment", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "host") return "localhost:3000";
          return null;
        }),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("http://localhost:3000");
    });

    test("should handle staging environment with subdomain", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-forwarded-host") return "staging-auth.company.com";
          if (key === "host") return "staging-internal.vercel.app";
          return null;
        }),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("https://staging-auth.company.com");
    });

    test("should prioritize x-zitadel-instance-host in multi-tenant scenario", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-instance-host") return "customer.zitadel.cloud";
          if (key === "x-forwarded-host") return "accounts.company.com";
          if (key === "host") return "internal.vercel.app";
          return null;
        }),
      } as any;

      const result = getInstanceHost(mockHeaders);
      expect(result).toBe("customer.zitadel.cloud");
    });
  });

  describe("Edge cases", () => {
    test("should handle IPv4 addresses", () => {
      const mockHeaders = {
        get: vi.fn(() => "192.168.1.100:3000"),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("https://192.168.1.100:3000");
    });

    test("should handle IPv6 addresses", () => {
      const mockHeaders = {
        get: vi.fn(() => "[::1]:3000"),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("https://[::1]:3000");
    });

    test("should handle hosts with ports", () => {
      const mockHeaders = {
        get: vi.fn(() => "zitadel.com:8080"),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("https://zitadel.com:8080");
    });

    test("should handle localhost with different ports", () => {
      const mockHeaders = {
        get: vi.fn(() => "localhost:8080"),
      } as any;

      const result = getPublicHostWithProtocol(mockHeaders);
      expect(result).toBe("http://localhost:8080");
    });
  });

  describe("getPublicHost", () => {
    test("should use x-zitadel-public-host when available", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-public-host") return "public.zitadel.cloud";
          if (key === "x-forwarded-host") return "accounts.company.com";
          return null;
        }),
      } as any;

      const result = getPublicHost(mockHeaders);
      expect(result).toBe("public.zitadel.cloud");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-public-host");
    });

    test("should use x-zitadel-forward-host when x-zitadel-public-host is not available", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-public-host") return null;
          if (key === "x-zitadel-forward-host") return "forward.zitadel.cloud";
          if (key === "x-forwarded-host") return "accounts.company.com";
          return null;
        }),
      } as any;

      const result = getPublicHost(mockHeaders);
      expect(result).toBe("forward.zitadel.cloud");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-public-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-forward-host");
    });

    test("should use x-forwarded-host when neither x-zitadel-public-host nor x-zitadel-forward-host is available", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-public-host") return null;
          if (key === "x-zitadel-forward-host") return null;
          if (key === "x-forwarded-host") return "accounts.company.com";
          if (key === "host") return "internal.server";
          return null;
        }),
      } as any;

      const result = getPublicHost(mockHeaders);
      expect(result).toBe("accounts.company.com");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-public-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-zitadel-forward-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-forwarded-host");
    });

    test("should fall back to host when x-forwarded-host is not available", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-public-host") return null;
          if (key === "x-zitadel-forward-host") return null;
          if (key === "x-forwarded-host") return null;
          if (key === "host") return "localhost:3000";
          return null;
        }),
      } as any;

      const result = getPublicHost(mockHeaders);
      expect(result).toBe("localhost:3000");
      expect(mockHeaders.get).toHaveBeenCalledWith("x-forwarded-host");
      expect(mockHeaders.get).toHaveBeenCalledWith("host");
    });

    test("should throw error when no host is found", () => {
      const mockHeaders = {
        get: vi.fn(() => null),
      } as any;

      expect(() => getPublicHost(mockHeaders)).toThrow("No host found in headers");
    });

    test("should differ from getInstanceHost when x-zitadel-instance-host is present", () => {
      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-instance-host") return "instance.zitadel.cloud";
          if (key === "x-forwarded-host") return "accounts.company.com";
          if (key === "host") return "internal.server";
          return null;
        }),
      } as any;

      const instanceHost = getInstanceHost(mockHeaders);
      const publicHost = getPublicHost(mockHeaders);

      expect(instanceHost).toBe("instance.zitadel.cloud");
      expect(publicHost).toBe("accounts.company.com");
      expect(instanceHost).not.toBe(publicHost);
    });
  });
});

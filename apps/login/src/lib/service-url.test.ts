import { describe, expect, test, beforeEach, afterEach, vi } from "vitest";
import { getServiceConfig, constructUrl } from "./service-url";
import { NextRequest } from "next/server";

describe("Service URL utilities", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  describe("getServiceConfig", () => {
    test("should throw when ZITADEL_API_URL is not set", () => {
      process.env.ZITADEL_API_URL = undefined as any;

      const mockHeaders = {
        get: vi.fn(() => null),
      } as any;

      expect(() => getServiceConfig(mockHeaders)).toThrow("ZITADEL_API_URL is not set");
    });

    test("should return only baseUrl when x-zitadel-forward-host is not present (self-hosted)", () => {
      process.env.ZITADEL_API_URL = "https://zitadel.mycompany.com";

      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-forward-host") return null;
          if (key === "host") return "mycompany.com";
          return null;
        }),
      } as any;

      const result = getServiceConfig(mockHeaders);

      expect(result.serviceConfig.baseUrl).toBe("https://zitadel.mycompany.com");
      expect(result.serviceConfig.instanceHost).toBeUndefined();
      expect(result.serviceConfig.publicHost).toBe("mycompany.com");
    });

    test("should use x-zitadel-forward-host when present (multi-tenant)", () => {
      process.env.ZITADEL_API_URL = "https://api.zitadel.cloud";

      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-forward-host") return "customer.zitadel.cloud";
          if (key === "host") return "customer.zitadel.cloud";
          return null;
        }),
      } as any;

      const result = getServiceConfig(mockHeaders);

      expect(result.serviceConfig.baseUrl).toBe("https://api.zitadel.cloud");
      expect(result.serviceConfig.instanceHost).toBe("customer.zitadel.cloud");
      expect(result.serviceConfig.publicHost).toBe("customer.zitadel.cloud");
    });

    test("should strip protocol from instanceHost and publicHost", () => {
      process.env.ZITADEL_API_URL = "https://api.zitadel.cloud";

      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-forward-host") return "https://customer.zitadel.cloud";
          if (key === "host") return "customer.zitadel.cloud";
          return null;
        }),
      } as any;

      const result = getServiceConfig(mockHeaders);

      expect(result.serviceConfig.instanceHost).toBe("customer.zitadel.cloud");
      expect(result.serviceConfig.publicHost).toBe("customer.zitadel.cloud");
    });

    test("should throw when host header is missing", () => {
      process.env.ZITADEL_API_URL = "https://api.zitadel.cloud";

      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-forward-host") return null;
          if (key === "host") return null;
          return null;
        }),
      } as any;

      expect(() => getServiceConfig(mockHeaders)).toThrow("No host found in headers");
    });

    test("should handle host with port number", () => {
      process.env.ZITADEL_API_URL = "https://api.zitadel.cloud";

      const mockHeaders = {
        get: vi.fn((key: string) => {
          if (key === "x-zitadel-forward-host") return "customer.zitadel.cloud:443";
          if (key === "host") return "customer.zitadel.cloud:443";
          return null;
        }),
      } as any;

      const result = getServiceConfig(mockHeaders);

      expect(result.serviceConfig.publicHost).toBe("customer.zitadel.cloud:443");
    });
  });

  describe("constructUrl", () => {
    test("should construct URL with x-zitadel-forward-host when present", () => {
      process.env.NEXT_PUBLIC_BASE_PATH = "";
      const mockRequest = {
        headers: {
          get: vi.fn((key: string) => {
            if (key === "x-zitadel-forward-host") return "customer.zitadel.cloud";
            if (key === "host") return "customer.zitadel.cloud";
            return null;
          }),
        },
        nextUrl: {
          protocol: "https:",
        },
      } as any;

      const result = constructUrl(mockRequest as NextRequest, "/test");

      expect(result.hostname).toBe("customer.zitadel.cloud");
      expect(result.pathname).toBe("/test");
      expect(result.protocol).toBe("https:");
    });

    test("should fall back to x-forwarded-host when x-zitadel-forward-host is not present", () => {
      process.env.NEXT_PUBLIC_BASE_PATH = "";
      const mockRequest = {
        headers: {
          get: vi.fn((key: string) => {
            if (key === "x-zitadel-forward-host") return null;
            if (key === "x-forwarded-host") return "mycompany.com";
            return null;
          }),
        },
        nextUrl: {
          protocol: "https:",
        },
      } as any;

      const result = constructUrl(mockRequest as NextRequest, "/oauth/authorize");

      expect(result.hostname).toBe("mycompany.com");
      expect(result.pathname).toBe("/oauth/authorize");
    });

    test("should fall back to host header when no forwarded headers present", () => {
      const mockRequest = {
        headers: {
          get: vi.fn((key: string) => {
            if (key === "x-zitadel-forward-host") return null;
            if (key === "x-forwarded-host") return null;
            if (key === "host") return "localhost:3000";
            return null;
          }),
        },
        nextUrl: {
          protocol: "http:",
        },
      } as any;

      const result = constructUrl(mockRequest as NextRequest, "/test");

      expect(result.hostname).toBe("localhost");
      expect(result.port).toBe("3000");
    });

    test("should use protocol from nextUrl.protocol (not from headers)", () => {
      const mockRequest = {
        headers: {
          get: vi.fn((key: string) => {
            // Even if x-forwarded-proto is present, it should be ignored
            if (key === "x-forwarded-proto") return "http";
            if (key === "host") return "example.com";
            return null;
          }),
        },
        nextUrl: {
          protocol: "https:", // This should be used
        },
      } as any;

      const result = constructUrl(mockRequest as NextRequest, "/test");

      // Should use https: from nextUrl.protocol, not http from header
      expect(result.protocol).toBe("https:");
    });

    test("should include base path when NEXT_PUBLIC_BASE_PATH is set", () => {
      process.env.NEXT_PUBLIC_BASE_PATH = "/login";

      const mockRequest = {
        headers: {
          get: vi.fn((key: string) => {
            if (key === "host") return "example.com";
            return null;
          }),
        },
        nextUrl: {
          protocol: "https:",
        },
      } as any;

      const result = constructUrl(mockRequest as NextRequest, "/test");

      expect(result.pathname).toBe("/login/test");
    });
  });
});

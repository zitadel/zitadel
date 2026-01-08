import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import { v4 as uuidv4 } from "uuid";
import {
  getFingerprintId,
  setFingerprintIdCookie,
  getFingerprintIdCookie,
  getOrSetFingerprintId,
  getUserAgent,
} from "./fingerprint";

// Mock dependencies
vi.mock("uuid");
vi.mock("next/headers", () => ({
  cookies: vi.fn(),
  headers: vi.fn(),
}));
vi.mock("next/server", () => ({
  userAgent: vi.fn(),
}));
vi.mock("@zitadel/client", () => ({
  create: vi.fn(),
}));

import { cookies, headers } from "next/headers";
import { userAgent } from "next/server";
import { create } from "@zitadel/client";

describe("fingerprint", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("getFingerprintId", () => {
    it("should generate a UUID v4", async () => {
      const mockUuid = "550e8400-e29b-41d4-a716-446655440000";
      vi.mocked(uuidv4).mockReturnValue(mockUuid as any);

      const result = await getFingerprintId();

      expect(result).toBe(mockUuid);
      expect(uuidv4).toHaveBeenCalledOnce();
    });

    it("should generate unique IDs on multiple calls", async () => {
      const uuid1 = "550e8400-e29b-41d4-a716-446655440000";
      const uuid2 = "6ba7b810-9dad-11d1-80b4-00c04fd430c8";

      vi.mocked(uuidv4)
        .mockReturnValueOnce(uuid1 as any)
        .mockReturnValueOnce(uuid2 as any);

      const result1 = await getFingerprintId();
      const result2 = await getFingerprintId();

      expect(result1).toBe(uuid1);
      expect(result2).toBe(uuid2);
      expect(result1).not.toBe(result2);
    });
  });

  describe("setFingerprintIdCookie", () => {
    it("should set cookie with correct name and value", async () => {
      const mockSet = vi.fn();
      const mockCookies = { set: mockSet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const fingerprintId = "test-fingerprint-id";
      await setFingerprintIdCookie(fingerprintId);

      expect(mockSet).toHaveBeenCalledWith({
        name: "fingerprintId",
        value: fingerprintId,
        httpOnly: true,
        path: "/",
        maxAge: 31536000,
      });
    });

    it("should set httpOnly flag to true", async () => {
      const mockSet = vi.fn();
      const mockCookies = { set: mockSet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      await setFingerprintIdCookie("test-id");

      const callArgs = mockSet.mock.calls[0][0];
      expect(callArgs.httpOnly).toBe(true);
    });

    it("should set cookie path to root", async () => {
      const mockSet = vi.fn();
      const mockCookies = { set: mockSet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      await setFingerprintIdCookie("test-id");

      const callArgs = mockSet.mock.calls[0][0];
      expect(callArgs.path).toBe("/");
    });

    it("should set maxAge to 1 year (31536000 seconds)", async () => {
      const mockSet = vi.fn();
      const mockCookies = { set: mockSet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      await setFingerprintIdCookie("test-id");

      const callArgs = mockSet.mock.calls[0][0];
      expect(callArgs.maxAge).toBe(31536000);
      expect(callArgs.maxAge).toBe(365 * 24 * 60 * 60);
    });

    it("should handle special characters in fingerprint ID", async () => {
      const mockSet = vi.fn();
      const mockCookies = { set: mockSet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const specialId = "test-id-with-!@#$%^&*()";
      await setFingerprintIdCookie(specialId);

      const callArgs = mockSet.mock.calls[0][0];
      expect(callArgs.value).toBe(specialId);
    });

    it("should handle empty fingerprint ID", async () => {
      const mockSet = vi.fn();
      const mockCookies = { set: mockSet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      await setFingerprintIdCookie("");

      expect(mockSet).toHaveBeenCalledWith(
        expect.objectContaining({
          value: "",
        }),
      );
    });

    it("should handle very long fingerprint IDs", async () => {
      const mockSet = vi.fn();
      const mockCookies = { set: mockSet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const longId = "a".repeat(1000);
      await setFingerprintIdCookie(longId);

      const callArgs = mockSet.mock.calls[0][0];
      expect(callArgs.value).toBe(longId);
    });
  });

  describe("getFingerprintIdCookie", () => {
    it("should retrieve fingerprint cookie", async () => {
      const mockGet = vi.fn().mockReturnValue({
        name: "fingerprintId",
        value: "test-fingerprint",
      });
      const mockCookies = { get: mockGet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const result = await getFingerprintIdCookie();

      expect(mockGet).toHaveBeenCalledWith("fingerprintId");
      expect(result).toEqual({
        name: "fingerprintId",
        value: "test-fingerprint",
      });
    });

    it("should return undefined if cookie doesn't exist", async () => {
      const mockGet = vi.fn().mockReturnValue(undefined);
      const mockCookies = { get: mockGet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const result = await getFingerprintIdCookie();

      expect(result).toBeUndefined();
    });

    it("should handle null cookie value", async () => {
      const mockGet = vi.fn().mockReturnValue(null);
      const mockCookies = { get: mockGet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const result = await getFingerprintIdCookie();

      expect(result).toBeNull();
    });
  });

  describe("getOrSetFingerprintId", () => {
    it("should return existing cookie value if present", async () => {
      const existingId = "existing-fingerprint-id";
      const mockGet = vi.fn().mockReturnValue({
        name: "fingerprintId",
        value: existingId,
      });
      const mockCookies = { get: mockGet, set: vi.fn() };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const result = await getOrSetFingerprintId();

      expect(result).toBe(existingId);
      expect(mockCookies.set).not.toHaveBeenCalled();
    });

    it("should generate and set new ID if cookie doesn't exist", async () => {
      const newUuid = "new-generated-uuid";
      const mockGet = vi.fn().mockReturnValue(undefined);
      const mockSet = vi.fn();
      const mockCookies = { get: mockGet, set: mockSet };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);
      vi.mocked(uuidv4).mockReturnValue(newUuid as any);

      const result = await getOrSetFingerprintId();

      expect(result).toBe(newUuid);
      expect(mockSet).toHaveBeenCalledWith({
        name: "fingerprintId",
        value: newUuid,
        httpOnly: true,
        path: "/",
        maxAge: 31536000,
      });
    });

    it("should not regenerate ID on subsequent calls with existing cookie", async () => {
      const existingId = "persistent-id";
      const mockGet = vi.fn().mockReturnValue({
        name: "fingerprintId",
        value: existingId,
      });
      const mockCookies = { get: mockGet, set: vi.fn() };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const result1 = await getOrSetFingerprintId();
      const result2 = await getOrSetFingerprintId();

      expect(result1).toBe(existingId);
      expect(result2).toBe(existingId);
      expect(uuidv4).not.toHaveBeenCalled();
    });
  });

  describe("getUserAgent", () => {
    it("should create user agent with all components", async () => {
      const mockFingerprintId = "test-fingerprint";
      const mockGet = vi.fn().mockReturnValue({
        name: "fingerprintId",
        value: mockFingerprintId,
      });
      const mockCookies = { get: mockGet, set: vi.fn() };
      vi.mocked(cookies).mockResolvedValue(mockCookies as any);

      const mockHeaders = new Map([
        ["user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"],
        ["x-forwarded-for", "192.168.1.1"],
      ]);
      vi.mocked(headers).mockResolvedValue({
        get: (key: string) => mockHeaders.get(key) ?? null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: { type: "desktop", vendor: "Dell", model: "XPS" },
        engine: { name: "Blink", version: "120.0" },
        os: { name: "Windows", version: "10" },
        browser: { name: "Chrome", version: "120.0" },
      } as any);

      const mockUserAgent = {
        ip: "192.168.1.1",
        header: {
          "user-agent": {
            values: ["Mozilla/5.0 (Windows NT 10.0; Win64; x64)"],
          },
        },
        description: "Chrome, 120.0, desktop, Dell, XPS, , Blink, 120.0, , Windows, 10, ",
        fingerprintId: mockFingerprintId,
      };
      vi.mocked(create).mockReturnValue(mockUserAgent as any);

      const result = await getUserAgent();

      expect(result).toEqual(mockUserAgent);
      expect(create).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({
          fingerprintId: mockFingerprintId,
          ip: "192.168.1.1",
        }),
      );
    });

    it("should use x-forwarded-for header for IP if available", async () => {
      const mockGet = vi.fn().mockReturnValue({
        value: "test-id",
      });
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);

      const mockHeaders = new Map([
        ["x-forwarded-for", "203.0.113.1"],
        ["remoteAddress", "192.168.1.1"],
      ]);
      vi.mocked(headers).mockResolvedValue({
        get: (key: string) => mockHeaders.get(key) ?? null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: {},
        engine: {},
        os: {},
        browser: {},
      } as any);

      await getUserAgent();

      expect(create).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({
          ip: "203.0.113.1",
        }),
      );
    });

    it("should fallback to remoteAddress if x-forwarded-for is not available", async () => {
      const mockGet = vi.fn().mockReturnValue({
        value: "test-id",
      });
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);

      const mockHeaders = new Map([["remoteAddress", "192.168.1.100"]]);
      vi.mocked(headers).mockResolvedValue({
        get: (key: string) => mockHeaders.get(key) ?? null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: {},
        engine: {},
        os: {},
        browser: {},
      } as any);

      await getUserAgent();

      expect(create).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({
          ip: "192.168.1.100",
        }),
      );
    });

    it("should use empty string for IP if no headers available", async () => {
      const mockGet = vi.fn().mockReturnValue({
        value: "test-id",
      });
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);

      const mockHeaders = new Map();
      vi.mocked(headers).mockResolvedValue({
        get: (key: string) => mockHeaders.get(key) ?? null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: {},
        engine: {},
        os: {},
        browser: {},
      } as any);

      await getUserAgent();

      expect(create).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({
          ip: "",
        }),
      );
    });

    it("should parse user agent header values", async () => {
      const mockGet = vi.fn().mockReturnValue({
        value: "test-id",
      });
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);

      const userAgentString = "Mozilla/5.0,AppleWebKit/537.36,Chrome/120.0";
      const mockHeaders = new Map([["user-agent", userAgentString]]);
      vi.mocked(headers).mockResolvedValue({
        get: (key: string) => mockHeaders.get(key) ?? null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: {},
        engine: {},
        os: {},
        browser: {},
      } as any);

      await getUserAgent();

      expect(create).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({
          header: {
            "user-agent": {
              values: ["Mozilla/5.0", "AppleWebKit/537.36", "Chrome/120.0"],
            },
          },
        }),
      );
    });

    it("should build description from device components", async () => {
      const mockGet = vi.fn().mockReturnValue({
        value: "test-id",
      });
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);

      vi.mocked(headers).mockResolvedValue({
        get: () => null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: { type: "mobile", vendor: "Apple", model: "iPhone 15" },
        engine: { name: "WebKit", version: "605.1.15" },
        os: { name: "iOS", version: "17.0" },
        browser: { name: "Safari", version: "17.0" },
      } as any);

      await getUserAgent();

      const createCall = vi.mocked(create).mock.calls[0]?.[1];
      expect(createCall).toBeDefined();
      expect(createCall?.description).toContain("Safari");
      expect(createCall?.description).toContain("17.0");
      expect(createCall?.description).toContain("mobile");
      expect(createCall?.description).toContain("Apple");
      expect(createCall?.description).toContain("iPhone 15");
      expect(createCall?.description).toContain("WebKit");
      expect(createCall?.description).toContain("iOS");
    });

    it("should handle missing device components gracefully", async () => {
      const mockGet = vi.fn().mockReturnValue({
        value: "test-id",
      });
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);

      vi.mocked(headers).mockResolvedValue({
        get: () => null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: {},
        engine: {},
        os: {},
        browser: {},
      } as any);

      await getUserAgent();

      expect(create).toHaveBeenCalled();
      const createCall = vi.mocked(create).mock.calls[0]?.[1];
      expect(createCall?.description).toBeDefined();
    });

    it("should handle partial device information", async () => {
      const mockGet = vi.fn().mockReturnValue({
        value: "test-id",
      });
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);

      vi.mocked(headers).mockResolvedValue({
        get: () => null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: { vendor: "Samsung" },
        engine: { name: "Blink" },
        os: { name: "Android" },
        browser: { name: "Chrome" },
      } as any);

      await getUserAgent();

      const createCall = vi.mocked(create).mock.calls[0]?.[1];
      expect(createCall).toBeDefined();
      expect(createCall?.description).toContain("Samsung");
      expect(createCall?.description).toContain("Blink");
      expect(createCall?.description).toContain("Android");
      expect(createCall?.description).toContain("Chrome");
    });

    it("should generate new fingerprint ID if none exists", async () => {
      const newUuid = "newly-generated-fingerprint";
      const mockGet = vi.fn().mockReturnValue(undefined);
      const mockSet = vi.fn();
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: mockSet } as any);
      vi.mocked(uuidv4).mockReturnValue(newUuid as any);

      vi.mocked(headers).mockResolvedValue({
        get: () => null,
      } as any);

      vi.mocked(userAgent).mockReturnValue({
        device: {},
        engine: {},
        os: {},
        browser: {},
      } as any);

      await getUserAgent();

      expect(create).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({
          fingerprintId: newUuid,
        }),
      );
    });
  });
});

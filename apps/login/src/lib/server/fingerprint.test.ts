import {
  addIdpFlowProof,
  getBindingIdCookie,
  getOrSetBindingId,
  getOrSetFingerprintId,
  getUserAgent,
} from "@/lib/server/fingerprint";
import { isIdpFlowProofValid } from "@/lib/server/idp-flow-proof";
import { v4 as uuidv4 } from "uuid";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { create } from "@zitadel/client";
import { cookies, headers } from "next/headers";
import { userAgent } from "next/server";
import { getSecuritySettings } from "./zitadel";

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
vi.mock("./service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://api.example.com" } })),
}));
vi.mock("./zitadel", () => ({
  getSecuritySettings: vi.fn(),
}));

describe("fingerprint", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // default: login is not embedded in an iframe
    vi.mocked(getSecuritySettings).mockResolvedValue(undefined);
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllEnvs();
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
        sameSite: "lax",
        secure: false,
      });
    });

    it("should use sameSite none with secure when iframe embedding is enabled", async () => {
      const mockSet = vi.fn();
      vi.mocked(cookies).mockResolvedValue({ get: vi.fn().mockReturnValue(undefined), set: mockSet } as any);
      vi.mocked(getSecuritySettings).mockResolvedValue({ embeddedIframe: { enabled: true } } as any);

      await getOrSetFingerprintId();

      const callArgs = mockSet.mock.calls[0][0];
      expect(callArgs.sameSite).toBe("none");
      expect(callArgs.secure).toBe(true);
    });

    it("should use sameSite lax when iframe embedding is disabled", async () => {
      const mockSet = vi.fn();
      vi.mocked(cookies).mockResolvedValue({ get: vi.fn().mockReturnValue(undefined), set: mockSet } as any);
      vi.mocked(getSecuritySettings).mockResolvedValue({ embeddedIframe: { enabled: false } } as any);

      await getOrSetFingerprintId();

      expect(mockSet.mock.calls[0][0].sameSite).toBe("lax");
    });

    it("should fall back to sameSite lax when the security settings can't be loaded", async () => {
      const mockSet = vi.fn();
      vi.mocked(cookies).mockResolvedValue({ get: vi.fn().mockReturnValue(undefined), set: mockSet } as any);
      vi.mocked(getSecuritySettings).mockRejectedValue(new Error("unavailable"));

      // must not throw, the fingerprint is still needed by every flow
      await getOrSetFingerprintId();

      expect(mockSet.mock.calls[0][0].sameSite).toBe("lax");
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

  describe("getOrSetBindingId", () => {
    it("should return the existing binding id cookie", async () => {
      const mockGet = vi.fn((name: string) => (name === "bindingId" ? { name, value: "binding-1" } : undefined));
      const mockSet = vi.fn();
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: mockSet } as any);

      expect(await getOrSetBindingId()).toBe("binding-1");
      expect(await getBindingIdCookie()).toEqual({ name: "bindingId", value: "binding-1" });
      expect(mockSet).not.toHaveBeenCalled();
    });

    it("should set a new random binding id with the iframe-aware attributes", async () => {
      const mockSet = vi.fn();
      vi.mocked(cookies).mockResolvedValue({ get: vi.fn().mockReturnValue(undefined), set: mockSet } as any);
      vi.mocked(uuidv4).mockReturnValue("new-binding-id" as any);
      vi.mocked(getSecuritySettings).mockResolvedValue({ embeddedIframe: { enabled: true } } as any);

      const result = await getOrSetBindingId();

      expect(result).toBe("new-binding-id");
      expect(mockSet).toHaveBeenCalledWith({
        name: "bindingId",
        value: "new-binding-id",
        httpOnly: true,
        path: "/",
        maxAge: 31536000,
        sameSite: "none",
        secure: true,
      });
    });
  });

  describe("addIdpFlowProof", () => {
    it("should set a flow proof bound to this browser and the session to link", async () => {
      vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");
      const mockGet = vi.fn((name: string) => (name === "bindingId" ? { name, value: "binding-1" } : undefined));
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);

      const params = new URLSearchParams();
      await addIdpFlowProof(params, "session-1");
      const signInParams = new URLSearchParams();
      await addIdpFlowProof(signInParams);

      expect(isIdpFlowProofValid(params.get("flowFingerprint") ?? undefined, "binding-1", "session-1")).toBe(true);
      expect(isIdpFlowProofValid(signInParams.get("flowFingerprint") ?? undefined, "binding-1")).toBe(true);
    });
  });

  describe("getUserAgent", () => {
    it("should send the fingerprintId but never the binding id to the API", async () => {
      const mockGet = vi.fn((name: string) => ({ name, value: name === "bindingId" ? "secret-binding" : "public-fp" }));
      vi.mocked(cookies).mockResolvedValue({ get: mockGet, set: vi.fn() } as any);
      vi.mocked(headers).mockResolvedValue({ get: () => null } as any);
      vi.mocked(userAgent).mockReturnValue({ device: {}, engine: {}, os: {}, browser: {} } as any);

      await getUserAgent();

      const sent = vi.mocked(create).mock.calls[0][1];
      expect(sent).toEqual(expect.objectContaining({ fingerprintId: "public-fp" }));
      expect(JSON.stringify(sent)).not.toContain("secret-binding");
      expect(mockGet).not.toHaveBeenCalledWith("bindingId");
    });

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

    describe("client ip", () => {
      // the ip getUserAgent sends for the given request headers
      async function ipFor(requestHeaders: Record<string, string>) {
        vi.mocked(cookies).mockResolvedValue({
          get: vi.fn().mockReturnValue({ value: "test-id" }),
          set: vi.fn(),
        } as any);
        vi.mocked(headers).mockResolvedValue({ get: (key: string) => requestHeaders[key] ?? null } as any);
        vi.mocked(userAgent).mockReturnValue({ device: {}, engine: {}, os: {}, browser: {} } as any);

        await getUserAgent();

        return vi.mocked(create).mock.calls.at(-1)?.[1]?.ip;
      }

      it.each([
        ["a single entry written by the trusted proxy", "198.51.100.7", "198.51.100.7"],
        ["the rightmost entry of a list", "203.0.113.5, 198.51.100.7", "198.51.100.7"],
        ["a list the client prepended entries to", "9.9.9.9, 1.2.3.4, 198.51.100.7", "198.51.100.7"],
        ["an entry with surrounding whitespace", "   198.51.100.7   ", "198.51.100.7"],
        ["an IPv4 entry with a port", "198.51.100.7:4711", "198.51.100.7"],
        ["a bare IPv6 entry", "2001:db8::1", "2001:db8::1"],
        ["a bracketed IPv6 entry", "[2001:db8::1]", "2001:db8::1"],
        ["a bracketed IPv6 entry with a port", "[2001:db8::1]:4711", "2001:db8::1"],
        ["nothing, for an obfuscated last entry", "198.51.100.7, unknown", ""],
        ["nothing, for a last entry that is no address", "not-an-ip", ""],
      ])("should send %s", async (_name, forwardedFor, expected) => {
        expect(await ipFor({ "x-forwarded-for": forwardedFor })).toBe(expected);
      });

      it("should skip the configured number of trusted proxy hops", async () => {
        vi.stubEnv("TRUSTED_PROXY_HOPS", "1");

        expect(await ipFor({ "x-forwarded-for": "9.9.9.9, 198.51.100.7, 203.0.113.5" })).toBe("198.51.100.7");
      });

      it("should send no ip when the list is shorter than the trusted hops", async () => {
        vi.stubEnv("TRUSTED_PROXY_HOPS", "2");

        expect(await ipFor({ "x-forwarded-for": "198.51.100.7" })).toBe("");
      });

      it("should fall back to the default hops for an invalid TRUSTED_PROXY_HOPS", async () => {
        vi.stubEnv("TRUSTED_PROXY_HOPS", "not-a-number");

        expect(await ipFor({ "x-forwarded-for": "9.9.9.9, 198.51.100.7" })).toBe("198.51.100.7");
      });

      it("should prefer the forwarded header over the direct connection address", async () => {
        expect(await ipFor({ "x-forwarded-for": "198.51.100.7", remoteAddress: "192.0.2.1" })).toBe("198.51.100.7");
      });

      it("should use the direct connection address when no proxy set the header", async () => {
        expect(await ipFor({ remoteAddress: "198.51.100.7:4711" })).toBe("198.51.100.7");
        expect(await ipFor({ "x-forwarded-for": "", remoteAddress: "198.51.100.7" })).toBe("198.51.100.7");
      });

      it("should send no ip when neither header is set", async () => {
        expect(await ipFor({})).toBe("");
      });
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

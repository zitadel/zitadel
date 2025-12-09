/**
 * Unit tests for U2F server actions.
 *
 * Tests core business logic for:
 * - U2F device registration
 * - U2F device verification
 * - Token name generation from user agent
 * - Session validation
 * - Error handling
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { addU2F, verifyU2F } from "./u2f";
import * as zitadelModule from "@/lib/zitadel";
import * as cookiesModule from "../cookies";
import type { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import type { GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";

vi.mock("@/lib/zitadel");
vi.mock("../cookies");
vi.mock("./host");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
}));
vi.mock("next/server", () => ({
  userAgent: vi.fn(() => ({
    browser: { name: "Firefox" },
    device: { vendor: "Dell", model: "" },
    os: { name: "Windows" },
  })),
}));
vi.mock("@/lib/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://zitadel-test.zitadel.cloud" } })),
}));

describe("U2F server actions", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockSessionId = "session123";
  const mockUserId = "user123";
  const mockU2FId = "u2f123";

  beforeEach(async () => {
    vi.clearAllMocks();

    const { getPublicHost } = await import("./host");
    vi.mocked(getPublicHost).mockReturnValue("zitadel.com:443");
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("addU2F", () => {
    it("should register U2F device successfully", async () => {
      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        loginName: "test@example.com",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      const mockSession: GetSessionResponse = {
        session: {
          id: mockSessionId,
          factors: {
            user: {
              id: mockUserId,
              loginName: "test@example.com",
            },
          },
        } as Session,
      } as GetSessionResponse;

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSession);
      vi.mocked(zitadelModule.registerU2F).mockResolvedValue({
        u2fId: mockU2FId,
      } as any);

      const result = await addU2F({ sessionId: mockSessionId });

      expect(zitadelModule.registerU2F).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        userId: mockUserId,
        domain: "zitadel.com",
      });
      expect(result).toEqual({ u2fId: mockU2FId });
    });

    it("should extract hostname from host with port", async () => {
      const { getPublicHost } = await import("./host");
      vi.mocked(getPublicHost).mockReturnValue("localhost:3000");

      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          factors: { user: { id: mockUserId } },
        } as Session,
      } as GetSessionResponse);
      vi.mocked(zitadelModule.registerU2F).mockResolvedValue({ u2fId: mockU2FId } as any);

      await addU2F({ sessionId: mockSessionId });

      expect(zitadelModule.registerU2F).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        userId: mockUserId,
        domain: "localhost",
      });
    });

    it("should return error when session cookie not found", async () => {
      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(null as any);

      const result = await addU2F({ sessionId: mockSessionId });

      expect(result).toEqual({ error: "Could not get session" });
      expect(zitadelModule.registerU2F).not.toHaveBeenCalled();
    });

    it("should return error when session has no user", async () => {
      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          id: mockSessionId,
          factors: {},
        } as Session,
      } as GetSessionResponse);

      const result = await addU2F({ sessionId: mockSessionId });

      expect(result).toEqual({ error: "Could not get session" });
    });

    it("should throw error when hostname cannot be extracted", async () => {
      const { getPublicHost } = await import("./host");
      vi.mocked(getPublicHost).mockReturnValue("");

      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          factors: { user: { id: mockUserId } },
        } as Session,
      } as GetSessionResponse);

      await expect(addU2F({ sessionId: mockSessionId })).rejects.toThrow("Could not get hostname");
    });
  });

  describe("verifyU2F", () => {
    const mockPublicKeyCredential = {
      id: "credential123",
      type: "public-key",
      rawId: new ArrayBuffer(8),
    };

    it("should verify U2F device successfully", async () => {
      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      const mockSession: GetSessionResponse = {
        session: {
          factors: {
            user: { id: mockUserId },
          },
        } as Session,
      } as GetSessionResponse;

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSession);
      vi.mocked(zitadelModule.verifyU2FRegistration).mockResolvedValue({
        details: { verified: true },
      } as any);

      await verifyU2F({
        u2fId: mockU2FId,
        publicKeyCredential: mockPublicKeyCredential,
        sessionId: mockSessionId,
      });

      expect(zitadelModule.verifyU2FRegistration).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        request: expect.objectContaining({
          u2fId: mockU2FId,
          publicKeyCredential: mockPublicKeyCredential,
          userId: mockUserId,
        }),
      });
    });

    it("should generate token name from user agent when not provided", async () => {
      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          factors: { user: { id: mockUserId } },
        } as Session,
      } as GetSessionResponse);
      vi.mocked(zitadelModule.verifyU2FRegistration).mockResolvedValue({} as any);

      await verifyU2F({
        u2fId: mockU2FId,
        publicKeyCredential: mockPublicKeyCredential,
        sessionId: mockSessionId,
      });

      expect(zitadelModule.verifyU2FRegistration).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        request: expect.objectContaining({
          tokenName: "Dell , Windows, Firefox",
        }),
      });
    });

    it("should use provided passkey name", async () => {
      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          factors: { user: { id: mockUserId } },
        } as Session,
      } as GetSessionResponse);
      vi.mocked(zitadelModule.verifyU2FRegistration).mockResolvedValue({} as any);

      await verifyU2F({
        u2fId: mockU2FId,
        passkeyName: "My YubiKey",
        publicKeyCredential: mockPublicKeyCredential,
        sessionId: mockSessionId,
      });

      expect(zitadelModule.verifyU2FRegistration).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        request: expect.objectContaining({
          tokenName: "My YubiKey",
        }),
      });
    });

    it("should return error when session has no user", async () => {
      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          factors: {},
        } as Session,
      } as GetSessionResponse);

      const result = await verifyU2F({
        u2fId: mockU2FId,
        publicKeyCredential: mockPublicKeyCredential,
        sessionId: mockSessionId,
      });

      expect(result).toEqual({ error: "Could not get session" });
      expect(zitadelModule.verifyU2FRegistration).not.toHaveBeenCalled();
    });

    it("should handle user agent with no device info", async () => {
      const { userAgent } = await import("next/server");
      vi.mocked(userAgent).mockReturnValue({
        browser: { name: "Safari" },
        device: { vendor: "", model: "" },
        os: { name: "macOS" },
      } as any);

      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          factors: { user: { id: mockUserId } },
        } as Session,
      } as GetSessionResponse);
      vi.mocked(zitadelModule.verifyU2FRegistration).mockResolvedValue({} as any);

      await verifyU2F({
        u2fId: mockU2FId,
        publicKeyCredential: mockPublicKeyCredential,
        sessionId: mockSessionId,
      });

      expect(zitadelModule.verifyU2FRegistration).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        request: expect.objectContaining({
          tokenName: " macOS, Safari",
        }),
      });
    });
  });
});

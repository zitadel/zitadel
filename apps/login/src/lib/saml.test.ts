import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import { getSAMLFormUID, setSAMLFormCookie, getSAMLFormCookie, loginWithSAMLAndSession } from "./saml";

// Mock dependencies
vi.mock("uuid");
vi.mock("next/headers", () => ({
  cookies: vi.fn(),
}));
vi.mock("./session", () => ({
  isSessionValid: vi.fn(),
}));
vi.mock("@/lib/server/loginname", () => ({
  sendLoginname: vi.fn(),
}));
vi.mock("@/lib/zitadel", () => ({
  createResponse: vi.fn(),
  getLoginSettings: vi.fn(),
}));
vi.mock("@zitadel/client", () => ({
  create: vi.fn((schema, data) => data),
}));

import { v4 as uuidv4 } from "uuid";
import { cookies } from "next/headers";
import { isSessionValid } from "./session";
import { sendLoginname } from "@/lib/server/loginname";
import {createResponse, getLoginSettings, ServiceConfig} from "@/lib/zitadel";

describe("saml", () => {
  let mockCookies: any;

  beforeEach(() => {
    vi.clearAllMocks();
    mockCookies = {
      get: vi.fn(),
      set: vi.fn(),
    };
    vi.mocked(cookies).mockResolvedValue(mockCookies);

    // Suppress console logs during tests
    vi.spyOn(console, "log").mockImplementation(() => {});
    vi.spyOn(console, "warn").mockImplementation(() => {});
    vi.spyOn(console, "error").mockImplementation(() => {});
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("getSAMLFormUID", () => {
    it("should generate a UUID v4", async () => {
      const mockUuid = "saml-uid-123-456";
      vi.mocked(uuidv4).mockReturnValue(mockUuid as any);

      const result = await getSAMLFormUID();

      expect(result).toBe(mockUuid);
      expect(uuidv4).toHaveBeenCalledOnce();
    });

    it("should generate unique UIDs on multiple calls", async () => {
      const uuid1 = "saml-uid-1";
      const uuid2 = "saml-uid-2";

      vi.mocked(uuidv4)
        .mockReturnValueOnce(uuid1 as any)
        .mockReturnValueOnce(uuid2 as any);

      const result1 = await getSAMLFormUID();
      const result2 = await getSAMLFormUID();

      expect(result1).toBe(uuid1);
      expect(result2).toBe(uuid2);
      expect(result1).not.toBe(result2);
    });
  });

  describe("setSAMLFormCookie", () => {
    const mockUid = "test-saml-uid";

    beforeEach(() => {
      vi.mocked(uuidv4).mockReturnValue(mockUid as any);
    });

    it("should set SAML form cookie with correct parameters", async () => {
      const samlValue = "SAMLRequest=encodedvalue";

      const uid = await setSAMLFormCookie(samlValue);

      expect(uid).toBe(mockUid);
      expect(mockCookies.set).toHaveBeenCalledWith({
        name: mockUid,
        value: samlValue,
        httpOnly: true,
        secure: expect.any(Boolean),
        sameSite: "lax",
        path: "/",
        maxAge: 5 * 60,
      });
    });

    it("should set httpOnly to true", async () => {
      await setSAMLFormCookie("test-value");

      const callArgs = mockCookies.set.mock.calls[0][0];
      expect(callArgs.httpOnly).toBe(true);
    });

    it("should set sameSite to 'lax'", async () => {
      await setSAMLFormCookie("test-value");

      const callArgs = mockCookies.set.mock.calls[0][0];
      expect(callArgs.sameSite).toBe("lax");
    });

    it("should set path to root", async () => {
      await setSAMLFormCookie("test-value");

      const callArgs = mockCookies.set.mock.calls[0][0];
      expect(callArgs.path).toBe("/");
    });

    it("should set maxAge to 5 minutes (300 seconds)", async () => {
      await setSAMLFormCookie("test-value");

      const callArgs = mockCookies.set.mock.calls[0][0];
      expect(callArgs.maxAge).toBe(300);
      expect(callArgs.maxAge).toBe(5 * 60);
    });

    it("should warn when value exceeds 4000 characters", async () => {
      const largeSamlValue = "x".repeat(5000);

      await setSAMLFormCookie(largeSamlValue);

      expect(console.warn).toHaveBeenCalledWith(expect.stringContaining("SAML form cookie value is large"));
      expect(console.warn).toHaveBeenCalledWith(expect.stringContaining("5000 characters"));
    });

    it("should handle empty SAML value", async () => {
      const uid = await setSAMLFormCookie("");

      expect(uid).toBe(mockUid);
      expect(mockCookies.set).toHaveBeenCalledWith(
        expect.objectContaining({
          value: "",
        }),
      );
    });

    it("should handle very long SAML values", async () => {
      const longValue = "x".repeat(10000);

      const uid = await setSAMLFormCookie(longValue);

      expect(uid).toBe(mockUid);
      expect(console.warn).toHaveBeenCalled();
    });

    it("should throw error when cookie setting fails", async () => {
      const error = new Error("Cookie set failed");
      mockCookies.set.mockRejectedValue(error);

      await expect(setSAMLFormCookie("test-value")).rejects.toThrow("Failed to set SAML form cookie");
    });

    it("should log error with details when cookie setting fails", async () => {
      const error = new Error("Storage quota exceeded");
      mockCookies.set.mockRejectedValue(error);

      await expect(setSAMLFormCookie("test-value")).rejects.toThrow();

      expect(console.error).toHaveBeenCalledWith(
        expect.stringContaining("Failed to set SAML form cookie"),
        expect.objectContaining({
          error,
          uid: mockUid,
        }),
      );
    });
  });

  describe("getSAMLFormCookie", () => {
    it("should retrieve SAML form cookie by UID", async () => {
      const uid = "test-saml-uid";
      const cookieValue = "SAMLRequest=encodedvalue";

      mockCookies.get.mockReturnValue({
        name: uid,
        value: cookieValue,
      });

      const result = await getSAMLFormCookie(uid);

      expect(result).toBe(cookieValue);
      expect(mockCookies.get).toHaveBeenCalledWith(uid);
    });

    it("should return null if cookie not found", async () => {
      mockCookies.get.mockReturnValue(undefined);

      const result = await getSAMLFormCookie("non-existent-uid");

      expect(result).toBeNull();
      expect(console.warn).toHaveBeenCalledWith(expect.stringContaining("SAML form cookie not found"));
    });

    it("should return null if cookie has empty value", async () => {
      mockCookies.get.mockReturnValue({
        name: "test-uid",
        value: "",
      });

      const result = await getSAMLFormCookie("test-uid");

      expect(result).toBeNull();
      expect(console.warn).toHaveBeenCalledWith(expect.stringContaining("empty value"));
    });

    it("should handle errors gracefully", async () => {
      const error = new Error("Cookie read failed");
      mockCookies.get.mockImplementation(() => {
        throw error;
      });

      const result = await getSAMLFormCookie("test-uid");

      expect(result).toBeNull();
      expect(console.error).toHaveBeenCalledWith(expect.stringContaining("Error retrieving SAML form cookie"), error);
    });

    it("should handle null cookie value", async () => {
      mockCookies.get.mockReturnValue({
        name: "test-uid",
        value: null,
      });

      const result = await getSAMLFormCookie("test-uid");

      expect(result).toBeNull();
    });

    it("should retrieve large SAML values", async () => {
      const largeValue = "x".repeat(5000);
      mockCookies.get.mockReturnValue({
        name: "test-uid",
        value: largeValue,
      });

      const result = await getSAMLFormCookie("test-uid");

      expect(result).toBe(largeValue);
    });
  });

  describe("loginWithSAMLAndSession", () => {
    const mockSession: any = {
      id: "session-123",
      factors: {
        user: {
          loginName: "user@example.com",
          organizationId: "org-123",
        },
      },
    };

    const mockCookie: any = {
      id: "session-123",
      token: "session-token-123",
      loginName: "user@example.com",
    };

    const baseParams = {
      serviceConfig: {baseUrl: "https://example.com"} as ServiceConfig,
      samlRequest: "saml-request-id",
      sessionId: "session-123",
      sessions: [mockSession],
      sessionCookies: [mockCookie],
    };

    it("should return redirect URL for valid session", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockResolvedValue({
        url: "https://sp.example.com/acs",
      } as any);

      const result = await loginWithSAMLAndSession(baseParams);

      expect(result).toEqual({
        redirect: "https://sp.example.com/acs",
      });
    });

    it("should validate session using isSessionValid", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockResolvedValue({
        url: "https://sp.example.com/acs",
      } as any);

      await loginWithSAMLAndSession(baseParams);

      expect(isSessionValid).toHaveBeenCalledWith({
        serviceConfig: baseParams.serviceConfig,
        session: mockSession,
      });
    });

    it("should call createResponse with session info", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockResolvedValue({
        url: "https://sp.example.com/acs",
      } as any);

      await loginWithSAMLAndSession(baseParams);

      expect(createResponse).toHaveBeenCalledWith({
        serviceConfig: baseParams.serviceConfig,
        req: expect.objectContaining({
          samlRequestId: "saml-request-id",
          responseKind: {
            case: "session",
            value: {
              sessionId: "session-123",
              sessionToken: "session-token-123",
            },
          },
        }),
      });
    });

    it("should redirect to re-authentication if session is invalid", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(false);
      vi.mocked(sendLoginname).mockResolvedValue({
        redirect: "/loginname",
      });

      const result = await loginWithSAMLAndSession(baseParams);

      expect(sendLoginname).toHaveBeenCalledWith({
        loginName: "user@example.com",
        organization: "org-123",
        requestId: "saml_saml-request-id",
      });
      expect(result).toEqual({
        redirect: "/loginname",
      });
    });

    it("should return error if session not found", async () => {
      const params = {
        ...baseParams,
        sessionId: "non-existent-session",
      };

      const result = await loginWithSAMLAndSession(params);

      expect(result).toEqual({
        error: "Session not found or invalid",
      });
    });

    it("should return error if session cookie not found", async () => {
      const params = {
        ...baseParams,
        sessionCookies: [],
      };

      vi.mocked(isSessionValid).mockResolvedValue(true);

      const result = await loginWithSAMLAndSession(params);

      expect(result).toEqual({
        error: "Session not found or invalid",
      });
    });

    it("should return error if createResponse returns no URL", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockResolvedValue({} as any);

      const result = await loginWithSAMLAndSession(baseParams);

      expect(result).toEqual({
        error: "An error occurred!",
      });
    });

    it("should handle SAML request already used (error code 9)", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockRejectedValue({ code: 9 });
      vi.mocked(getLoginSettings).mockResolvedValue({
        defaultRedirectUri: "https://example.com/default",
      } as any);

      const result = await loginWithSAMLAndSession(baseParams);

      expect(result).toEqual({
        redirect: "https://example.com/default",
      });
    });

    it("should fallback to /signedin if no defaultRedirectUri (error code 9)", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockRejectedValue({ code: 9 });
      vi.mocked(getLoginSettings).mockResolvedValue({} as any);

      const result = await loginWithSAMLAndSession(baseParams);

      expect(result).toEqual({
        redirect: expect.stringContaining("/signedin"),
      });
    });

    it("should include loginName and organization in /signedin redirect", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockRejectedValue({ code: 9 });
      vi.mocked(getLoginSettings).mockResolvedValue({} as any);

      const result = await loginWithSAMLAndSession(baseParams);

      expect(result).toEqual({
        redirect: expect.stringContaining("loginName=user%40example.com"),
      });
      expect(result).toEqual({
        redirect: expect.stringContaining("organization=org-123"),
      });
    });

    it("should return error for unknown errors from createResponse", async () => {
      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockRejectedValue(new Error("Network error"));

      const result = await loginWithSAMLAndSession(baseParams);

      expect(result).toEqual({
        error: "Unknown error occurred",
      });
    });

    it("should handle session without user factors", async () => {
      const sessionWithoutUser = {
        id: "session-123",
        factors: {},
      } as any;

      const params = {
        ...baseParams,
        sessions: [sessionWithoutUser],
      };

      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockResolvedValue({
        url: "https://sp.example.com/acs",
      } as any);

      const result = await loginWithSAMLAndSession(params);

      expect(result).toEqual({
        redirect: "https://sp.example.com/acs",
      });
    });

    it("should handle multiple sessions with correct selection", async () => {
      const otherSession: any = {
        id: "session-456",
        factors: {
          user: {
            loginName: "other@example.com",
          },
        },
      };

      const params = {
        ...baseParams,
        sessions: [otherSession, mockSession],
      };

      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockResolvedValue({
        url: "https://sp.example.com/acs",
      } as any);

      await loginWithSAMLAndSession(params);

      expect(isSessionValid).toHaveBeenCalledWith({
        serviceConfig: baseParams.serviceConfig,
        session: mockSession,
      });
    });

    it("should match session cookie by ID", async () => {
      const otherCookie: any = {
        id: "session-456",
        token: "other-token",
        loginName: "other@example.com",
      };

      const params = {
        ...baseParams,
        sessionCookies: [otherCookie, mockCookie],
      };

      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockResolvedValue({
        url: "https://sp.example.com/acs",
      } as any);

      await loginWithSAMLAndSession(params);

      expect(createResponse).toHaveBeenCalledWith({
        serviceConfig: baseParams.serviceConfig,
        req: expect.objectContaining({
          responseKind: {
            case: "session",
            value: {
              sessionId: "session-123",
              sessionToken: "session-token-123",
            },
          },
        }),
      });
    });

    it("should handle error code 9 without loginName", async () => {
      const sessionWithoutLoginName = {
        id: "session-123",
        factors: {
          user: {
            organizationId: "org-123",
          },
        },
      } as any;

      const params = {
        ...baseParams,
        sessions: [sessionWithoutLoginName],
      };

      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockRejectedValue({ code: 9 });
      vi.mocked(getLoginSettings).mockResolvedValue({} as any);

      const result = await loginWithSAMLAndSession(params);

      expect(result).toEqual({
        redirect: expect.stringContaining("/signedin"),
      });
      expect(result).toEqual({
        redirect: expect.not.stringContaining("loginName="),
      });
    });

    it("should handle error code 9 without organization", async () => {
      const sessionWithoutOrg = {
        id: "session-123",
        factors: {
          user: {
            loginName: "user@example.com",
          },
        },
      } as any;

      const params = {
        ...baseParams,
        sessions: [sessionWithoutOrg],
      };

      vi.mocked(isSessionValid).mockResolvedValue(true);
      vi.mocked(createResponse).mockRejectedValue({ code: 9 });
      vi.mocked(getLoginSettings).mockResolvedValue({} as any);

      const result = await loginWithSAMLAndSession(params);

      expect(result).toEqual({
        redirect: expect.stringContaining("/signedin"),
      });
      expect(result).toEqual({
        redirect: expect.not.stringContaining("organization="),
      });
    });
  });
});

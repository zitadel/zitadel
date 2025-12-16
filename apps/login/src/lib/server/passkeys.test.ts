import { describe, expect, test, vi, beforeEach } from "vitest";
import { sendPasskey } from "./passkeys";

// Mock all the dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: vi.fn(),
  Duration: vi.fn(),
  Timestamp: vi.fn(),
  timestampDate: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  getLoginSettings: vi.fn(),
  getUserByID: vi.fn(),
}));

vi.mock("./cookie", () => ({
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("../cookies", () => ({
  getSessionCookieById: vi.fn(),
  getSessionCookieByLoginName: vi.fn(),
  getMostRecentSessionCookie: vi.fn(),
}));

vi.mock("../verify-helper", () => ({
  checkEmailVerification: vi.fn(),
}));

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

// Mock translations - returns the key itself for testing
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

describe("sendPasskey", () => {
  let mockHeaders: any;
  let mockGetServiceUrlFromHeaders: any;
  let mockGetLoginSettings: any;
  let mockGetUserByID: any;
  let mockSetSessionAndUpdateCookie: any;
  let mockGetSessionCookieById: any;
  let mockGetSessionCookieByLoginName: any;
  let mockGetMostRecentSessionCookie: any;
  let mockCheckEmailVerification: any;
  let mockCompleteFlowOrGetUrl: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    // Import mocked modules
    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getLoginSettings, getUserByID } = await import("../zitadel");
    const { setSessionAndUpdateCookie } = await import("./cookie");
    const { getSessionCookieById, getSessionCookieByLoginName, getMostRecentSessionCookie } = await import("../cookies");
    const { checkEmailVerification } = await import("../verify-helper");
    const { completeFlowOrGetUrl } = await import("../client");

    // Setup mocks
    mockHeaders = vi.mocked(headers);
    mockGetServiceUrlFromHeaders = vi.mocked(getServiceConfig);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockGetUserByID = vi.mocked(getUserByID);
    mockSetSessionAndUpdateCookie = vi.mocked(setSessionAndUpdateCookie);
    mockGetSessionCookieById = vi.mocked(getSessionCookieById);
    mockGetSessionCookieByLoginName = vi.mocked(getSessionCookieByLoginName);
    mockGetMostRecentSessionCookie = vi.mocked(getMostRecentSessionCookie);
    mockCheckEmailVerification = vi.mocked(checkEmailVerification);
    mockCompleteFlowOrGetUrl = vi.mocked(completeFlowOrGetUrl);

    // Default mock implementations
    mockHeaders.mockResolvedValue(new Headers());
    mockGetServiceUrlFromHeaders.mockReturnValue({
      serviceUrl: "https://example.com",
    });
    mockGetLoginSettings.mockResolvedValue({
      multiFactorCheckLifetime: {
        seconds: BigInt(300),
        nanos: 0,
      },
    });
  });

  describe("Session Cookie Retrieval", () => {
    test("should return error when session cookie not found by sessionId", async () => {
      mockGetSessionCookieById.mockResolvedValue(null);

      const result = await sendPasskey({
        sessionId: "test-session-id",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        error: "verify.errors.couldNotFindSession",
      });
      expect(mockGetSessionCookieById).toHaveBeenCalledWith({
        sessionId: "test-session-id",
      });
    });

    test("should return error when session cookie is not found by loginName", async () => {
      mockGetSessionCookieByLoginName.mockResolvedValue(null);

      const result = await sendPasskey({
        loginName: "test@example.com",
        organization: "org-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        error: "verify.errors.couldNotFindSession",
      });
      expect(mockGetSessionCookieByLoginName).toHaveBeenCalledWith({
        loginName: "test@example.com",
        organization: "org-123",
      });
    });

    test("should return error when no session cookie found (most recent fallback)", async () => {
      mockGetMostRecentSessionCookie.mockResolvedValue(null);

      const result = await sendPasskey({
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        error: "verify.errors.couldNotFindSession",
      });
      expect(mockGetMostRecentSessionCookie).toHaveBeenCalled();
    });
  });

  describe("Session Update Failures", () => {
    beforeEach(() => {
      mockGetSessionCookieById.mockResolvedValue({
        id: "session-123",
        token: "session-token",
        loginName: "test@example.com",
      });
    });

    test("should return error when setSessionAndUpdateCookie fails", async () => {
      mockSetSessionAndUpdateCookie.mockResolvedValue({
        error: "Failed to update session",
      });

      const result = await sendPasskey({
        sessionId: "session-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        error: "verify.errors.couldNotUpdateSession",
      });
    });

    test("should return error when setSessionAndUpdateCookie returns undefined", async () => {
      mockSetSessionAndUpdateCookie.mockResolvedValue(undefined as any);

      const result = await sendPasskey({
        sessionId: "session-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        error: "verify.errors.couldNotUpdateSession",
      });
    });
  });

  describe("User Validation", () => {
    beforeEach(() => {
      mockGetSessionCookieById.mockResolvedValue({
        id: "session-123",
        token: "session-token",
        loginName: "test@example.com",
      });
      mockSetSessionAndUpdateCookie.mockResolvedValue({
        sessionId: "session-123",
        sessionToken: "new-token",
        factors: {
          user: {
            id: "user-123",
          },
        },
      });
    });

    test("should return error when getUserByID fails", async () => {
      mockGetUserByID.mockRejectedValue(new Error("User not found"));

      const result = await sendPasskey({
        sessionId: "session-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        error: "verify.errors.couldNotGetUser",
      });
    });
  });

  describe("Successful Passkey Verification", () => {
    beforeEach(() => {
      mockGetSessionCookieById.mockResolvedValue({
        id: "session-123",
        token: "session-token",
        loginName: "test@example.com",
      });
      mockSetSessionAndUpdateCookie.mockResolvedValue({
        id: "session-123",
        factors: {
          user: {
            id: "user-123",
            loginName: "test@example.com",
          },
        },
      });
      mockGetUserByID.mockResolvedValue({
        user: {
          id: "user-123",
          type: {
            case: "human",
            value: {
              email: {
                isVerified: true,
              },
            },
          },
        },
      });
      mockCheckEmailVerification.mockResolvedValue(true);
    });

    test("should redirect on successful verification without requestId", async () => {
      mockCompleteFlowOrGetUrl.mockResolvedValue({ redirect: "/dashboard" });

      const result = await sendPasskey({
        sessionId: "session-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        redirect: "/dashboard",
      });
    });

    test("should redirect on successful verification with requestId", async () => {
      mockCompleteFlowOrGetUrl.mockResolvedValue({ redirect: "/auth/callback" });

      const result = await sendPasskey({
        sessionId: "session-123",
        requestId: "request-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        redirect: "/auth/callback",
      });
    });

    test("should redirect for email verification when required", async () => {
      mockCheckEmailVerification.mockReturnValue({ redirect: "/verify" });

      const result = await sendPasskey({
        sessionId: "session-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toHaveProperty("redirect");
      if ("redirect" in result) {
        expect(result.redirect).toContain("/verify");
      }
    });
  });

  describe("Fallback Error Handling - Critical Fix", () => {
    beforeEach(() => {
      mockGetSessionCookieById.mockResolvedValue({
        id: "session-123",
        token: "session-token",
        loginName: "test@example.com",
      });
      mockSetSessionAndUpdateCookie.mockResolvedValue({
        id: "session-123",
        factors: {
          user: {
            id: "user-123",
            loginName: "test@example.com",
          },
        },
      });
      mockGetUserByID.mockResolvedValue({
        user: {
          id: "user-123",
          type: {
            case: "human",
            value: {
              email: {
                isVerified: true,
              },
            },
          },
        },
      });
      mockCheckEmailVerification.mockResolvedValue(true);
    });

    test("should return fallback error when completeFlowOrGetUrl returns undefined", async () => {
      mockCompleteFlowOrGetUrl.mockResolvedValue(undefined);

      const result = await sendPasskey({
        sessionId: "session-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        error: "verify.errors.couldNotDetermineRedirect",
      });
    });

    test("should return fallback error when completeFlowOrGetUrl returns empty string", async () => {
      mockCompleteFlowOrGetUrl.mockResolvedValue("");

      const result = await sendPasskey({
        sessionId: "session-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(result).toEqual({
        error: "verify.errors.couldNotDetermineRedirect",
      });
    });
  });

  describe("Custom Lifetime Handling", () => {
    beforeEach(() => {
      mockGetSessionCookieById.mockResolvedValue({
        id: "session-123",
        token: "session-token",
        loginName: "test@example.com",
      });
      mockSetSessionAndUpdateCookie.mockResolvedValue({
        id: "session-123",
        factors: {
          user: {
            id: "user-123",
            loginName: "test@example.com",
          },
        },
      });
      mockGetUserByID.mockResolvedValue({
        user: {
          id: "user-123",
          type: {
            case: "human",
            value: {
              email: {
                isVerified: true,
              },
            },
          },
        },
      });
      mockCheckEmailVerification.mockResolvedValue(true);
      mockCompleteFlowOrGetUrl.mockResolvedValue({ redirect: "/dashboard" });
    });

    test("should use custom lifetime when provided", async () => {
      await sendPasskey({
        sessionId: "session-123",
        lifetime: { seconds: BigInt(600), nanos: 0 } as any,
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({
          lifetime: expect.objectContaining({
            seconds: BigInt(600),
          }),
        }),
      );
    });

    test("should use default lifetime from login settings when not provided", async () => {
      await sendPasskey({
        sessionId: "session-123",
        checks: { webAuthN: { credentialAssertionData: {} } } as any,
      });

      expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({
          lifetime: expect.objectContaining({
            seconds: BigInt(300),
          }),
        }),
      );
    });
  });
});

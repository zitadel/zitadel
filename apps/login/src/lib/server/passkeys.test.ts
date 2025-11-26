/**
 * Unit tests for passkeys server actions.
 *
 * Tests core business logic for:
 * - Passkey registration with session validation
 * - Passkey registration via code/link
 * - Session validity checking
 * - User verification requirements
 * - Passkey verification and naming
 * - Authentication flow completion
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { registerPasskeyLink, verifyPasskeyRegistration, sendPasskey } from "./passkeys";
import * as zitadelModule from "@/lib/zitadel";
import * as cookiesModule from "../cookies";
import * as cookieModule from "./cookie";
import * as verifyHelperModule from "../verify-helper";
import * as clientModule from "../client";
import { create, Duration, Timestamp } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import type { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import type { GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";

vi.mock("@/lib/zitadel");
vi.mock("../cookies");
vi.mock("./cookie");
vi.mock("../verify-helper");
vi.mock("../client");
vi.mock("./host");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
}));
vi.mock("next/server", () => ({
  userAgent: vi.fn(() => ({
    browser: { name: "Chrome" },
    device: { vendor: "Apple", model: "iPhone" },
    os: { name: "iOS" },
  })),
}));
vi.mock("@/lib/service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(() => ({ serviceUrl: "https://zitadel-test.zitadel.cloud" })),
}));

describe("Passkeys server actions", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockUserId = "user123";
  const mockSessionId = "session123";
  const mockLoginName = "test@example.com";
  const mockOrganization = "org123";

  const mockFutureDate = {
    seconds: BigInt(Math.floor(Date.now() / 1000) + 3600),
    nanos: 0,
  } as Timestamp;

  const mockPastDate = {
    seconds: BigInt(Math.floor(Date.now() / 1000) - 3600),
    nanos: 0,
  } as Timestamp;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { getOriginalHost } = await import("./host");
    vi.mocked(getOriginalHost).mockResolvedValue("localhost:3000");
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("registerPasskeyLink", () => {
    it("should return error when neither sessionId nor userId provided", async () => {
      const result = await registerPasskeyLink({});

      expect(result).toEqual({ error: "Either sessionId or userId must be provided" });
    });

    it("should throw error with valid session when code not provided", async () => {
      const mockSession: GetSessionResponse = {
        session: {
          id: mockSessionId,
          expirationDate: mockFutureDate,
          factors: {
            user: { id: mockUserId },
            password: { verifiedAt: mockFutureDate },
          },
        } as Session,
      } as GetSessionResponse;

      const mockSessionCookie = {
        id: mockSessionId,
        token: "token123",
        loginName: mockLoginName,
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSession);

      await expect(registerPasskeyLink({ sessionId: mockSessionId })).rejects.toThrow("Missing code in response");
    });

    it("should require user verification when session is expired and user has auth methods", async () => {
      const mockSession: GetSessionResponse = {
        session: {
          id: mockSessionId,
          expirationDate: mockPastDate,
          factors: {
            user: { id: mockUserId },
          },
        } as Session,
      } as GetSessionResponse;

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue({
        id: mockSessionId,
        token: "token123",
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSession);
      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [1],
      } as any);

      const result = await registerPasskeyLink({ sessionId: mockSessionId });

      expect(result).toEqual({
        error: "You have to authenticate or have a valid User Verification Check",
      });
    });

    it("should check user verification when session is expired and no auth methods exist", async () => {
      const mockSession: GetSessionResponse = {
        session: {
          id: mockSessionId,
          expirationDate: mockPastDate,
          factors: {
            user: { id: mockUserId },
          },
        } as Session,
      } as GetSessionResponse;

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue({
        id: mockSessionId,
        token: "token123",
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSession);
      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);
      vi.mocked(verifyHelperModule.checkUserVerification).mockResolvedValue(false);

      const result = await registerPasskeyLink({ sessionId: mockSessionId });

      expect(result).toEqual({ error: "User Verification Check has to be done" });
      expect(verifyHelperModule.checkUserVerification).toHaveBeenCalledWith(mockUserId);
    });

    it("should create registration code when user verification passes but no code provided", async () => {
      const mockSession: GetSessionResponse = {
        session: {
          id: mockSessionId,
          expirationDate: mockPastDate,
          factors: {
            user: { id: mockUserId },
          },
        } as Session,
      } as GetSessionResponse;

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue({
        id: mockSessionId,
        token: "token123",
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSession);
      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [],
      } as any);
      vi.mocked(verifyHelperModule.checkUserVerification).mockResolvedValue(true);
      vi.mocked(zitadelModule.createPasskeyRegistrationLink).mockResolvedValue({
        code: { id: "code456", code: "XYZ789" },
      } as any);
      vi.mocked(zitadelModule.registerPasskey).mockResolvedValue({
        passkeyId: "passkey456",
      } as any);

      await registerPasskeyLink({ sessionId: mockSessionId });

      expect(zitadelModule.createPasskeyRegistrationLink).toHaveBeenCalled();
      expect(zitadelModule.registerPasskey).toHaveBeenCalledWith(
        expect.objectContaining({
          code: { id: "code456", code: "XYZ789" },
        }),
      );
    });

    it("should handle userId + code flow and create session", async () => {
      const mockUser = {
        userId: mockUserId,
        preferredLoginName: mockLoginName,
      };

      const mockCreatedSession: Session = {
        id: "newsession123",
        factors: {
          user: { id: mockUserId, loginName: mockLoginName },
        },
      } as Session;

      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: mockUser as any,
      } as any);
      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue(mockCreatedSession);
      vi.mocked(zitadelModule.registerPasskey).mockResolvedValue({
        passkeyId: "passkey789",
      } as any);

      const result = await registerPasskeyLink({
        userId: mockUserId,
        code: "DEF456",
        codeId: "code789",
      });

      expect(zitadelModule.getUserByID).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
      });
      expect(cookieModule.createSessionAndUpdateCookie).toHaveBeenCalledWith({
        checks: expect.objectContaining({
          user: expect.objectContaining({
            search: {
              case: "loginName",
              value: mockLoginName,
            },
          }),
        }),
        requestId: undefined,
      });
      expect(zitadelModule.registerPasskey).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
        code: { id: "code789", code: "DEF456" },
        domain: "localhost",
      });
      expect(result).toEqual({ passkeyId: "passkey789" });
    });

    it("should return error when user not found in userId flow", async () => {
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: null,
      } as any);

      const result = await registerPasskeyLink({
        userId: mockUserId,
        code: "CODE123",
        codeId: "codeId123",
      });

      expect(result).toEqual({ error: "User not found" });
    });

    it("should return error when session has no user", async () => {
      const mockSession: GetSessionResponse = {
        session: {
          id: mockSessionId,
          factors: {},
        } as Session,
      } as GetSessionResponse;

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue({
        id: mockSessionId,
        token: "token123",
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSession);

      const result = await registerPasskeyLink({ sessionId: mockSessionId });

      expect(result).toEqual({ error: "Could not determine user from session" });
    });
  });

  describe("verifyPasskeyRegistration", () => {
    it("should verify passkey with sessionId", async () => {
      const mockPublicKeyCredential = { id: "cred123", type: "public-key" };
      const mockSession: GetSessionResponse = {
        session: {
          factors: {
            user: { id: mockUserId },
          },
        } as Session,
      } as GetSessionResponse;

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue({
        id: mockSessionId,
        token: "token123",
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSession);
      vi.mocked(zitadelModule.verifyPasskeyRegistration).mockResolvedValue({} as any);

      await verifyPasskeyRegistration({
        passkeyId: "passkey123",
        publicKeyCredential: mockPublicKeyCredential,
        sessionId: mockSessionId,
      });

      expect(zitadelModule.verifyPasskeyRegistration).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        request: expect.objectContaining({
          passkeyId: "passkey123",
          publicKeyCredential: mockPublicKeyCredential,
          userId: mockUserId,
        }),
      });
    });

    it("should verify passkey with userId", async () => {
      const mockPublicKeyCredential = { id: "cred456", type: "public-key" };

      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: { userId: mockUserId } as any,
      } as any);
      vi.mocked(zitadelModule.verifyPasskeyRegistration).mockResolvedValue({} as any);

      await verifyPasskeyRegistration({
        passkeyId: "passkey456",
        publicKeyCredential: mockPublicKeyCredential,
        userId: mockUserId,
      });

      expect(zitadelModule.getUserByID).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
      });
      expect(zitadelModule.verifyPasskeyRegistration).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        request: expect.objectContaining({
          userId: mockUserId,
        }),
      });
    });

    it("should generate passkey name from user agent when not provided", async () => {
      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue({
        id: mockSessionId,
        token: "token123",
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          factors: { user: { id: mockUserId } },
        } as Session,
      } as GetSessionResponse);
      vi.mocked(zitadelModule.verifyPasskeyRegistration).mockResolvedValue({} as any);

      await verifyPasskeyRegistration({
        passkeyId: "passkey789",
        publicKeyCredential: {},
        sessionId: mockSessionId,
      });

      expect(zitadelModule.verifyPasskeyRegistration).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        request: expect.objectContaining({
          passkeyName: "Apple iPhone, iOS, Chrome",
        }),
      });
    });

    it("should use provided passkey name", async () => {
      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue({
        id: mockSessionId,
        token: "token123",
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue({
        session: {
          factors: { user: { id: mockUserId } },
        } as Session,
      } as GetSessionResponse);
      vi.mocked(zitadelModule.verifyPasskeyRegistration).mockResolvedValue({} as any);

      await verifyPasskeyRegistration({
        passkeyId: "passkey999",
        passkeyName: "My Custom Key",
        publicKeyCredential: {},
        sessionId: mockSessionId,
      });

      expect(zitadelModule.verifyPasskeyRegistration).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        request: expect.objectContaining({
          passkeyName: "My Custom Key",
        }),
      });
    });

    it("should throw error when neither sessionId nor userId provided", async () => {
      await expect(
        verifyPasskeyRegistration({
          passkeyId: "passkey000",
          publicKeyCredential: {},
        }),
      ).rejects.toThrow("Either sessionId or userId must be provided");
    });

    it("should throw error when user not found in userId flow", async () => {
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: null,
      } as any);

      await expect(
        verifyPasskeyRegistration({
          passkeyId: "passkey111",
          publicKeyCredential: {},
          userId: mockUserId,
        }),
      ).rejects.toThrow("User not found");
    });
  });

  describe("sendPasskey", () => {
    const mockSessionCookie = {
      id: mockSessionId,
      token: "token123",
      loginName: mockLoginName,
      organization: mockOrganization,
      creationTs: "1234567890",
      expirationTs: "1234567890",
      changeTs: "1234567890",
    };

    it("should return error when session not found", async () => {
      vi.mocked(cookiesModule.getMostRecentSessionCookie).mockResolvedValue(null as any);

      const result = await sendPasskey({});

      expect(result).toEqual({ error: "Could not find session" });
    });

    it("should update session and complete flow with requestId", async () => {
      const mockUpdatedSession: Session = {
        id: mockSessionId,
        factors: {
          user: {
            id: mockUserId,
            loginName: mockLoginName,
            organizationId: mockOrganization,
          },
        },
      } as Session;

      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        multiFactorCheckLifetime: { seconds: BigInt(900), nanos: 0 } as Duration,
      } as any);
      vi.mocked(cookieModule.setSessionAndUpdateCookie).mockResolvedValue(mockUpdatedSession as any);
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          type: {
            case: "human",
            value: { email: { email: "test@example.com", isVerified: true } },
          },
        } as any,
      } as any);
      vi.mocked(verifyHelperModule.checkEmailVerification).mockReturnValue(undefined);
      vi.mocked(clientModule.completeFlowOrGetUrl).mockResolvedValue({
        redirect: "https://app.example.com/callback",
      });

      const result = await sendPasskey({
        sessionId: mockSessionId,
        requestId: "oidc_request123",
        organization: mockOrganization,
        checks: create(ChecksSchema, { webAuthN: {} }),
      });

      expect(result).toEqual({ redirect: "https://app.example.com/callback" });
      expect(clientModule.completeFlowOrGetUrl).toHaveBeenCalledWith(
        {
          sessionId: mockSessionId,
          requestId: "oidc_request123",
          organization: mockOrganization,
        },
        undefined,
      );
    });

    it("should use default lifetime when not configured", async () => {
      const mockUpdatedSession: Session = {
        id: mockSessionId,
        factors: {
          user: {
            id: mockUserId,
            loginName: mockLoginName,
          },
        },
      } as Session;

      vi.mocked(cookiesModule.getMostRecentSessionCookie).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);
      vi.mocked(cookieModule.setSessionAndUpdateCookie).mockResolvedValue(mockUpdatedSession as any);
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          type: { case: "human", value: {} },
        } as any,
      } as any);
      vi.mocked(verifyHelperModule.checkEmailVerification).mockReturnValue(undefined);
      vi.mocked(clientModule.completeFlowOrGetUrl).mockResolvedValue({
        redirect: "https://app.example.com",
      });

      await sendPasskey({ checks: create(ChecksSchema, {}) });

      expect(cookieModule.setSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({
          lifetime: {
            seconds: BigInt(60 * 60 * 24),
            nanos: 0,
          },
        }),
      );
    });

    it("should redirect for email verification when required", async () => {
      const mockUpdatedSession: Session = {
        id: mockSessionId,
        factors: {
          user: { id: mockUserId },
        },
      } as Session;

      vi.mocked(cookiesModule.getMostRecentSessionCookie).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);
      vi.mocked(cookieModule.setSessionAndUpdateCookie).mockResolvedValue(mockUpdatedSession as any);
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          type: { case: "human", value: { email: { isVerified: false } } },
        } as any,
      } as any);
      vi.mocked(verifyHelperModule.checkEmailVerification).mockReturnValue({
        redirect: "/verify-email",
      });

      const result = await sendPasskey({});

      expect(result).toEqual({ redirect: "/verify-email" });
    });
  });
});

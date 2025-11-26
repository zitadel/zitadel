/**
 * Unit tests for password server actions.
 *
 * Tests core business logic for:
 * - Password reset functionality
 * - Password authentication and validation
 * - Lockout handling
 * - Failed attempt tracking
 * - Session creation after password verification
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { resetPassword, sendPassword } from "./password";
import * as zitadelModule from "@/lib/zitadel";
import * as cookiesModule from "../cookies";
import * as cookieModule from "./cookie";
import * as clientModule from "../client";
import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import type { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

vi.mock("@/lib/zitadel");
vi.mock("../cookies");
vi.mock("./cookie");
vi.mock("../client");
vi.mock("./host");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
}));
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() =>
    Promise.resolve((key: string, values?: Record<string, any>) => {
      const translations: Record<string, string> = {
        "errors.couldNotSendResetLink": "Could not send reset link",
        "errors.couldNotCreateSessionForUser": "Could not create session for user",
        "errors.couldNotVerifyPassword": "Could not verify password",
        "errors.failedToAuthenticate":
          "Failed to authenticate. {failedAttempts} of {maxPasswordAttempts} attempts used. {lockoutMessage}",
        "errors.failedToAuthenticateNoLimit": "Failed to authenticate. {failedAttempts} attempts used.",
        "errors.accountLockedContactAdmin": "Account locked, please contact admin",
      };
      let translation = translations[key] || key;

      if (values) {
        Object.keys(values).forEach((k) => {
          translation = translation.replace(`{${k}}`, String(values[k]));
        });
      }

      return translation;
    }),
  ),
}));
vi.mock("@/lib/service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(() => ({ serviceUrl: "https://zitadel-test.zitadel.cloud" })),
}));

describe("Password server actions", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockLoginName = "test@example.com";
  const mockUserId = "user123";
  const mockOrganization = "org123";
  const mockPassword = "SecurePassword123!";

  beforeEach(async () => {
    vi.clearAllMocks();

    const { getOriginalHostWithProtocol } = await import("./host");
    vi.mocked(getOriginalHostWithProtocol).mockResolvedValue("https://localhost:3000");
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("resetPassword", () => {
    it("should send password reset link for valid user", async () => {
      const mockUser = {
        userId: mockUserId,
        userName: mockLoginName,
      };

      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [mockUser],
        details: { totalResult: BigInt(1) },
      } as any);

      vi.mocked(zitadelModule.passwordReset).mockResolvedValue({} as any);

      const result = await resetPassword({
        loginName: mockLoginName,
        organization: mockOrganization,
      });

      expect(result).not.toHaveProperty("error");
      expect(zitadelModule.passwordReset).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
        urlTemplate: expect.stringContaining("code={{.Code}}"),
      });
    });

    it("should include requestId in reset URL when provided", async () => {
      const mockUser = {
        userId: mockUserId,
        userName: mockLoginName,
      };

      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [mockUser],
        details: { totalResult: BigInt(1) },
      } as any);

      vi.mocked(zitadelModule.passwordReset).mockResolvedValue({} as any);

      await resetPassword({
        loginName: mockLoginName,
        requestId: "oidc_request123",
      });

      expect(zitadelModule.passwordReset).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
        urlTemplate: expect.stringContaining("requestId=oidc_request123"),
      });
    });

    it("should return error when user not found", async () => {
      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [],
        details: { totalResult: BigInt(0) },
      } as any);

      const result = await resetPassword({
        loginName: "nonexistent@example.com",
      });

      expect(result).toEqual({ error: "Could not send reset link" });
      expect(zitadelModule.passwordReset).not.toHaveBeenCalled();
    });

    it("should return error when multiple users found", async () => {
      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [{}, {}],
        details: { totalResult: BigInt(2) },
      } as any);

      const result = await resetPassword({
        loginName: mockLoginName,
      });

      expect(result).toEqual({ error: "Could not send reset link" });
    });
  });

  describe("sendPassword", () => {
    const mockChecks = create(ChecksSchema, {
      password: { password: mockPassword },
    });

    const mockSession: Session = {
      id: "session123",
      factors: {
        user: {
          id: mockUserId,
          loginName: mockLoginName,
          organizationId: mockOrganization,
        },
      },
    } as Session;

    it("should authenticate user with valid password and existing session", async () => {
      const mockSessionCookie = {
        id: "session123",
        token: "token123",
        loginName: mockLoginName,
        creationTs: "1234567890",
        expirationTs: "1234567890",
        changeTs: "1234567890",
      };

      vi.mocked(cookiesModule.getSessionCookieByLoginName).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        passwordCheckLifetime: { seconds: BigInt(3600), nanos: 0 },
      } as any);
      vi.mocked(cookieModule.setSessionAndUpdateCookie).mockResolvedValue(mockSession as any);
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          userId: mockUserId,
          userName: mockLoginName,
          state: 2,
          type: {
            case: "human",
            value: {
              profile: { displayName: "Test User" },
              email: { email: mockLoginName },
            },
          },
        },
      } as any);
      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [1],
      } as any);
      vi.mocked(clientModule.completeFlowOrGetUrl).mockResolvedValue({
        redirect: "/accounts",
      } as any);

      const result = await sendPassword({
        loginName: mockLoginName,
        organization: mockOrganization,
        checks: mockChecks,
      });

      expect(result).toHaveProperty("redirect");
      expect(cookieModule.setSessionAndUpdateCookie).toHaveBeenCalled();
    });

    it("should create new session when no session cookie exists", async () => {
      const mockUser = {
        userId: mockUserId,
        userName: mockLoginName,
        preferredLoginName: mockLoginName,
      };

      vi.mocked(cookiesModule.getSessionCookieByLoginName).mockRejectedValue(new Error("No session"));
      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [mockUser as any],
        details: { totalResult: BigInt(1) },
      } as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        passwordCheckLifetime: { seconds: BigInt(3600), nanos: 0 },
      } as any);
      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue(mockSession);

      await sendPassword({
        loginName: mockLoginName,
        checks: mockChecks,
      });

      expect(cookieModule.createSessionAndUpdateCookie).toHaveBeenCalledWith({
        checks: expect.objectContaining({
          user: expect.objectContaining({ search: { case: "userId", value: mockUserId } }),
          password: expect.objectContaining({ password: mockPassword }),
        }),
        requestId: undefined,
        lifetime: { seconds: BigInt(3600), nanos: 0 },
      });
    });

    it("should handle failed password attempts with lockout", async () => {
      const error = {
        failedAttempts: 5,
      };

      vi.mocked(cookiesModule.getSessionCookieByLoginName).mockRejectedValue(new Error("No session"));
      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [{ userId: mockUserId }],
        details: { totalResult: BigInt(1) },
      } as any);
      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockRejectedValue(error);
      vi.mocked(zitadelModule.getLockoutSettings).mockResolvedValue({
        maxPasswordAttempts: BigInt(5),
      } as any);

      const result = await sendPassword({
        loginName: mockLoginName,
        checks: mockChecks,
      });

      expect(result).toMatchObject({
        error: expect.stringContaining("Failed to authenticate"),
      });
    });

    it("should handle failed attempts without lockout limit", async () => {
      const error = {
        failedAttempts: 3,
      };

      vi.mocked(cookiesModule.getSessionCookieByLoginName).mockRejectedValue(new Error("No session"));
      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [{ userId: mockUserId }],
        details: { totalResult: BigInt(1) },
      } as any);
      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockRejectedValue(error);
      vi.mocked(zitadelModule.getLockoutSettings).mockResolvedValue({
        maxPasswordAttempts: BigInt(0),
      } as any);

      const result = await sendPassword({
        loginName: mockLoginName,
        checks: mockChecks,
      });

      expect(result).toMatchObject({
        error: "Failed to authenticate. 3 attempts used.",
      });
    });

    it("should return error when user not found (security)", async () => {
      vi.mocked(cookiesModule.getSessionCookieByLoginName).mockRejectedValue(new Error("No session"));
      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [],
        details: { totalResult: BigInt(0) },
      } as any);

      const result = await sendPassword({
        loginName: "nonexistent@example.com",
        checks: mockChecks,
      });

      expect(result).toEqual({ error: "Could not verify password" });
    });

    it("should pass requestId through to session creation", async () => {
      const mockUser = {
        userId: mockUserId,
        userName: mockLoginName,
      };

      vi.mocked(cookiesModule.getSessionCookieByLoginName).mockRejectedValue(new Error("No session"));
      vi.mocked(zitadelModule.listUsers).mockResolvedValue({
        result: [mockUser as any],
        details: { totalResult: BigInt(1) },
      } as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);
      vi.mocked(cookieModule.createSessionAndUpdateCookie).mockResolvedValue(mockSession);

      await sendPassword({
        loginName: mockLoginName,
        checks: mockChecks,
        requestId: "oidc_request789",
      });

      expect(cookieModule.createSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({
          requestId: "oidc_request789",
        }),
      );
    });
  });
});

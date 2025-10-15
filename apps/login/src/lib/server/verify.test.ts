/**
 * Unit tests for the verify server actions.
 *
 * These tests replace the integration tests from verify.cy.ts and invite.cy.ts which tested:
 * - Email verification error handling
 * - Invite code verification with error handling
 * - Redirect to authenticator setup after successful invite verification
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { sendVerification } from "./verify";
import * as zitadelModule from "../zitadel";
import * as sessionModule from "../session";
import type { VerifyEmailResponse, VerifyInviteCodeResponse } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import type {
  GetUserByIDResponse,
  ListAuthenticationMethodTypesResponse,
} from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import type { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import type { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

// Mock all dependencies
vi.mock("../zitadel");
vi.mock("../session");
vi.mock("../cookies", () => ({
  getSessionCookieByLoginName: vi.fn(() => Promise.resolve(undefined)),
}));
vi.mock("../fingerprint");
vi.mock("../client");
vi.mock("./cookie");
vi.mock("./host");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
  cookies: vi.fn(() =>
    Promise.resolve({
      set: vi.fn(),
      get: vi.fn(),
    }),
  ),
}));
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() =>
    Promise.resolve((key: string) => {
      const translations: Record<string, string> = {
        "errors.couldNotVerifyEmail": "Could not verify email",
        "errors.couldNotVerifyInvite": "Could not verify invite",
        "errors.couldNotVerify": "Could not verify",
        "errors.couldNotLoadUser": "Could not load user",
      };
      return translations[key] || key;
    }),
  ),
}));
vi.mock("../service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(() => ({ serviceUrl: "https://zitadel-test.zitadel.cloud" })),
}));

describe("sendVerification server action", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockUserId = "221394658884845598";
  const mockCode = "abc123";
  const mockOrganization = "256088834543534543";
  const mockRequestId = "req123";

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("email verification", () => {
    it("should return error when email verification fails", async () => {
      // Mock verification failure
      vi.mocked(zitadelModule.verifyEmail).mockRejectedValue(new Error("Verification failed"));

      const result = await sendVerification({
        userId: mockUserId,
        code: mockCode,
        isInvite: false,
        organization: mockOrganization,
      });

      expect(result).toHaveProperty("error");
      if ("error" in result) {
        expect(result.error).toBe("Could not verify email");
      }
    });

    it("should return error when email verification returns error", async () => {
      // Mock verification returning error (simulating error case)
      vi.mocked(zitadelModule.verifyEmail).mockResolvedValue({
        error: "Invalid code",
      } as any);

      const result = await sendVerification({
        userId: mockUserId,
        code: mockCode,
        isInvite: false,
      });

      expect(result).toHaveProperty("error");
    });

    it("should return error when user cannot be loaded after verification", async () => {
      // Mock successful verification
      vi.mocked(zitadelModule.verifyEmail).mockResolvedValue({
        details: {},
      } as VerifyEmailResponse);

      // Mock user lookup failure
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: undefined,
      } as GetUserByIDResponse);

      const result = await sendVerification({
        userId: mockUserId,
        code: mockCode,
        isInvite: false,
      });

      expect(result).toHaveProperty("error");
      if ("error" in result) {
        expect(result.error).toBe("Could not load user");
      }
    });
  });

  describe("invite verification", () => {
    it("should return error when invite verification fails", async () => {
      // Mock verification failure
      vi.mocked(zitadelModule.verifyInviteCode).mockRejectedValue(new Error("Verification failed"));

      const result = await sendVerification({
        userId: mockUserId,
        code: mockCode,
        isInvite: true,
        organization: mockOrganization,
      });

      expect(result).toHaveProperty("error");
      if ("error" in result) {
        expect(result.error).toBe("Could not verify invite");
      }
    });

    it("should call verifyInviteCode with correct parameters", async () => {
      // Mock successful invite verification
      vi.mocked(zitadelModule.verifyInviteCode).mockResolvedValue({
        details: {},
      } as VerifyInviteCodeResponse);

      // Mock user
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          userId: mockUserId,
          state: 1,
          username: "john@example.com",
          preferredLoginName: "john@example.com",
          type: {
            case: "human",
            value: {
              email: {
                email: "john@example.com",
                isVerified: true,
              },
            },
          },
        },
      } as GetUserByIDResponse);

      await sendVerification({
        userId: mockUserId,
        code: mockCode,
        isInvite: true,
        organization: mockOrganization,
        requestId: mockRequestId,
      });

      // Verify that verifyInviteCode was called with correct parameters
      expect(vi.mocked(zitadelModule.verifyInviteCode)).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
        verificationCode: mockCode,
      });
    });
  });

  describe("successful verification flow", () => {
    it("should complete verification flow for existing user with authentication methods", async () => {
      // Mock successful verification
      vi.mocked(zitadelModule.verifyEmail).mockResolvedValue({
        details: {},
      } as any);

      // Mock user with authentication methods
      vi.mocked(zitadelModule.getUserByID).mockResolvedValue({
        user: {
          userId: mockUserId,
          state: 1,
          username: "john@example.com",
          type: {
            case: "human",
            value: {
              email: {
                email: "john@example.com",
                isVerified: true,
              },
            },
          },
        },
      } as GetUserByIDResponse);

      // Mock authentication methods - user has methods
      vi.mocked(zitadelModule.listAuthenticationMethodTypes).mockResolvedValue({
        authMethodTypes: [1], // PASSWORD
      } as ListAuthenticationMethodTypesResponse);

      // Mock login settings
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        allowUsernamePassword: true,
      } as LoginSettings);

      // Mock session - return undefined as it will be created later in the flow
      vi.mocked(sessionModule.loadMostRecentSession).mockResolvedValue(undefined as unknown as Session);

      const { completeFlowOrGetUrl } = await import("../client");
      vi.mocked(completeFlowOrGetUrl).mockResolvedValue({
        redirect: "/dashboard",
      } as { redirect: string });

      const result = await sendVerification({
        userId: mockUserId,
        code: mockCode,
        isInvite: false,
        loginName: "john@example.com",
        organization: mockOrganization,
        requestId: mockRequestId,
      });

      expect(result).toHaveProperty("redirect");
    });
  });
});

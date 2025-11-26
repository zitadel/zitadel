/**
 * Unit tests for session server actions.
 *
 * Tests core business logic for:
 * - MFA skip functionality
 * - Session continuation
 * - Session updates with checks and challenges
 * - Lifetime management
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { skipMFAAndContinueWithNextUrl, continueWithSession, updateSession } from "./session";
import * as zitadelModule from "@/lib/zitadel";
import * as cookiesModule from "../cookies";
import * as clientModule from "../client";
import { create, Duration } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import type { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

// Mock all dependencies
vi.mock("@/lib/zitadel");
vi.mock("../cookies");
vi.mock("../client");
vi.mock("./host");
vi.mock("./cookie");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
}));
vi.mock("@/lib/service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(() => ({ serviceUrl: "https://zitadel-test.zitadel.cloud" })),
}));

describe("Session server actions", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockUserId = "user123";
  const mockLoginName = "test@example.com";
  const mockOrganization = "org123";
  const mockSessionId = "session123";
  const mockRequestId = "oidc_request123";

  beforeEach(async () => {
    vi.clearAllMocks();

    const { getOriginalHost } = await import("./host");
    vi.mocked(getOriginalHost).mockResolvedValue("localhost:3000");
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("skipMFAAndContinueWithNextUrl", () => {
    it("should skip MFA and continue with sessionId and requestId", async () => {
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        defaultRedirectUri: "https://app.example.com",
      } as any);
      vi.mocked(zitadelModule.humanMFAInitSkipped).mockResolvedValue({} as any);
      vi.mocked(clientModule.completeFlowOrGetUrl).mockResolvedValue({
        redirect: "https://app.example.com/callback",
      });

      const result = await skipMFAAndContinueWithNextUrl({
        userId: mockUserId,
        sessionId: mockSessionId,
        requestId: mockRequestId,
        organization: mockOrganization,
      });

      expect(result).toEqual({ redirect: "https://app.example.com/callback" });
      expect(zitadelModule.humanMFAInitSkipped).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
      });
      expect(clientModule.completeFlowOrGetUrl).toHaveBeenCalledWith(
        {
          sessionId: mockSessionId,
          requestId: mockRequestId,
          organization: mockOrganization,
        },
        "https://app.example.com",
      );
    });

    it("should skip MFA and continue with loginName only", async () => {
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        defaultRedirectUri: "https://app.example.com",
      } as any);
      vi.mocked(zitadelModule.humanMFAInitSkipped).mockResolvedValue({} as any);
      vi.mocked(clientModule.completeFlowOrGetUrl).mockResolvedValue({
        redirect: "https://app.example.com/callback",
      });

      const result = await skipMFAAndContinueWithNextUrl({
        userId: mockUserId,
        loginName: mockLoginName,
        organization: mockOrganization,
      });

      expect(result).toEqual({ redirect: "https://app.example.com/callback" });
      expect(clientModule.completeFlowOrGetUrl).toHaveBeenCalledWith(
        {
          loginName: mockLoginName,
          organization: mockOrganization,
        },
        "https://app.example.com",
      );
    });

    it("should return error when neither sessionId nor loginName provided", async () => {
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);
      vi.mocked(zitadelModule.humanMFAInitSkipped).mockResolvedValue({} as any);

      const result = await skipMFAAndContinueWithNextUrl({
        userId: mockUserId,
      });

      expect(result).toEqual({ error: "Could not skip MFA and continue" });
    });

    it("should handle MFA skip errors", async () => {
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);
      vi.mocked(zitadelModule.humanMFAInitSkipped).mockRejectedValue(new Error("MFA skip failed"));

      await expect(
        skipMFAAndContinueWithNextUrl({
          userId: mockUserId,
          sessionId: mockSessionId,
          requestId: mockRequestId,
        }),
      ).rejects.toThrow("MFA skip failed");
    });
  });

  describe("continueWithSession", () => {
    const mockSession: Session = {
      id: mockSessionId,
      factors: {
        user: {
          id: mockUserId,
          loginName: mockLoginName,
          organizationId: mockOrganization,
        },
      },
    } as Session;

    it("should continue with session using sessionId and requestId", async () => {
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        defaultRedirectUri: "https://app.example.com",
      } as any);
      vi.mocked(clientModule.completeFlowOrGetUrl).mockResolvedValue({
        redirect: "https://app.example.com/callback",
      });

      const result = await continueWithSession({
        ...mockSession,
        requestId: mockRequestId,
      });

      expect(result).toEqual({ redirect: "https://app.example.com/callback" });
      expect(clientModule.completeFlowOrGetUrl).toHaveBeenCalledWith(
        {
          sessionId: mockSessionId,
          requestId: mockRequestId,
          organization: mockOrganization,
        },
        "https://app.example.com",
      );
    });

    it("should continue with session using loginName only", async () => {
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        defaultRedirectUri: "https://app.example.com",
      } as any);
      vi.mocked(clientModule.completeFlowOrGetUrl).mockResolvedValue({
        redirect: "https://app.example.com/callback",
      });

      const result = await continueWithSession(mockSession);

      expect(result).toEqual({ redirect: "https://app.example.com/callback" });
      expect(clientModule.completeFlowOrGetUrl).toHaveBeenCalledWith(
        {
          loginName: mockLoginName,
          organization: mockOrganization,
        },
        "https://app.example.com",
      );
    });

    it("should handle missing user factors", async () => {
      const invalidSession = {
        id: mockSessionId,
        factors: {},
      } as Session;

      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);

      const result = await continueWithSession(invalidSession);

      expect(result).toBeUndefined();
    });
  });

  describe("updateSession", () => {
    const mockSessionCookie = {
      id: mockSessionId,
      token: "token123",
      loginName: mockLoginName,
      organization: mockOrganization,
      creationTs: "1234567890",
      expirationTs: "1234567890",
      changeTs: "1234567890",
    };

    const mockChecks = create(ChecksSchema, {
      password: { password: "password123" },
    });

    it("should update session by sessionId", async () => {
      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);

      const { setSessionAndUpdateCookie } = await import("./cookie");
      vi.mocked(setSessionAndUpdateCookie).mockResolvedValue({
        id: mockSessionId,
      } as any);

      await updateSession({
        sessionId: mockSessionId,
        checks: mockChecks,
      });

      expect(cookiesModule.getSessionCookieById).toHaveBeenCalledWith({ sessionId: mockSessionId });
      expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({
          recentCookie: mockSessionCookie,
          checks: mockChecks,
        }),
      );
    });

    it("should update session by loginName", async () => {
      vi.mocked(cookiesModule.getSessionCookieByLoginName).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);

      const { setSessionAndUpdateCookie } = await import("./cookie");
      vi.mocked(setSessionAndUpdateCookie).mockResolvedValue({
        id: mockSessionId,
      } as any);

      await updateSession({
        loginName: mockLoginName,
        organization: mockOrganization,
        checks: mockChecks,
      });

      expect(cookiesModule.getSessionCookieByLoginName).toHaveBeenCalledWith({
        loginName: mockLoginName,
        organization: mockOrganization,
      });
    });

    it("should use most recent session when no identifiers provided", async () => {
      vi.mocked(cookiesModule.getMostRecentSessionCookie).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);

      const { setSessionAndUpdateCookie } = await import("./cookie");
      vi.mocked(setSessionAndUpdateCookie).mockResolvedValue({
        id: mockSessionId,
      } as any);

      await updateSession({
        checks: mockChecks,
      });

      expect(cookiesModule.getMostRecentSessionCookie).toHaveBeenCalled();
    });

    it("should return error when session cookie not found", async () => {
      vi.mocked(cookiesModule.getSessionCookieById).mockResolvedValue(null as any);

      const result = await updateSession({
        sessionId: mockSessionId,
        checks: mockChecks,
      });

      expect(result).toEqual({ error: "Could not find session" });
    });

    it("should set domain for webAuthN challenges", async () => {
      const mockChallenges: any = {
        webAuthN: {
          publicKeyCredentialCreationOptions: {},
        },
      };

      vi.mocked(cookiesModule.getMostRecentSessionCookie).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({} as any);

      const { setSessionAndUpdateCookie } = await import("./cookie");
      vi.mocked(setSessionAndUpdateCookie).mockResolvedValue({
        id: mockSessionId,
      } as any);

      await updateSession({
        challenges: mockChallenges,
        checks: mockChecks,
      });

      expect(mockChallenges.webAuthN.domain).toBe("localhost");
    });

    it("should use appropriate lifetime for different check types", async () => {
      vi.mocked(cookiesModule.getMostRecentSessionCookie).mockResolvedValue(mockSessionCookie as any);
      vi.mocked(zitadelModule.getLoginSettings).mockResolvedValue({
        multiFactorCheckLifetime: { seconds: BigInt(900), nanos: 0 } as Duration,
        secondFactorCheckLifetime: { seconds: BigInt(600), nanos: 0 } as Duration,
      } as any);

      const { setSessionAndUpdateCookie } = await import("./cookie");
      vi.mocked(setSessionAndUpdateCookie).mockResolvedValue({
        id: mockSessionId,
      } as any);

      // WebAuthN check should use multiFactorCheckLifetime
      const webAuthNChecks = create(ChecksSchema, {
        webAuthN: {},
      });

      await updateSession({
        checks: webAuthNChecks,
      });

      expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(
        expect.objectContaining({
          lifetime: { seconds: BigInt(900), nanos: 0 },
        }),
      );
    });
  });
});

/**
 * Unit tests for cookie server actions.
 *
 * Tests core business logic for:
 * - Session creation and cookie management
 * - IDP intent session creation
 * - Session updates with challenges
 * - Password attempt handling
 * - Session validation and error handling
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { createSessionAndUpdateCookie, createSessionForIdpAndUpdateCookie, setSessionAndUpdateCookie } from "./cookie";
import * as cookiesModule from "@/lib/cookies";
import * as zitadelModule from "@/lib/zitadel";
import { create, Duration, Timestamp } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import type { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import type { GetSessionResponse } from "@zitadel/proto/zitadel/session/v2/session_service_pb";

vi.mock("@/lib/cookies");
vi.mock("@/lib/zitadel");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
}));
vi.mock("@/lib/service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(() => ({ serviceUrl: "https://zitadel-test.zitadel.cloud" })),
}));

describe("Cookie server actions", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockSessionId = "session123";
  const mockSessionToken = "token123";
  const mockLoginName = "test@example.com";
  const mockUserId = "user123";
  const mockOrganization = "org123";

  const mockCreationDate = {
    seconds: BigInt(1234567890),
    nanos: 0,
  } as Timestamp;

  const mockSession: Session = {
    id: mockSessionId,
    creationDate: mockCreationDate,
    changeDate: mockCreationDate,
    expirationDate: mockCreationDate,
    factors: {
      user: {
        id: mockUserId,
        loginName: mockLoginName,
        organizationId: mockOrganization,
      },
    },
  } as Session;

  const mockSessionResponse: GetSessionResponse = {
    session: mockSession,
  } as GetSessionResponse;

  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(zitadelModule.getSecuritySettings).mockResolvedValue({
      embeddedIframe: { enabled: false },
    } as any);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("createSessionAndUpdateCookie", () => {
    const mockChecks = create(ChecksSchema, {
      user: { search: { case: "loginName", value: mockLoginName } },
      password: { password: "password123" },
    });

    const mockLifetime = {
      seconds: BigInt(3600),
      nanos: 0,
    } as Duration;

    it("should create session and update cookie successfully", async () => {
      vi.mocked(zitadelModule.createSessionFromChecks).mockResolvedValue({
        sessionId: mockSessionId,
        sessionToken: mockSessionToken,
      } as any);

      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSessionResponse);
      vi.mocked(cookiesModule.addSessionToCookie).mockResolvedValue(undefined);

      const result = await createSessionAndUpdateCookie({
        checks: mockChecks,
        requestId: "oidc_request123",
        lifetime: mockLifetime,
      });

      expect(result).toEqual(mockSession);
      expect(zitadelModule.createSessionFromChecks).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        checks: mockChecks,
        lifetime: mockLifetime,
      });
      expect(zitadelModule.getSession).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        sessionId: mockSessionId,
        sessionToken: mockSessionToken,
      });
      expect(cookiesModule.addSessionToCookie).toHaveBeenCalledWith({
        session: expect.objectContaining({
          id: mockSessionId,
          token: mockSessionToken,
          loginName: mockLoginName,
          organization: mockOrganization,
          requestId: "oidc_request123",
        }),
        iFrameEnabled: false,
      });
    });

    it("should use default lifetime when not provided", async () => {
      vi.mocked(zitadelModule.createSessionFromChecks).mockResolvedValue({
        sessionId: mockSessionId,
        sessionToken: mockSessionToken,
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSessionResponse);
      vi.mocked(cookiesModule.addSessionToCookie).mockResolvedValue(undefined);

      await createSessionAndUpdateCookie({
        checks: mockChecks,
        requestId: undefined,
      });

      expect(zitadelModule.createSessionFromChecks).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        checks: mockChecks,
        lifetime: {
          seconds: BigInt(24 * 60 * 60),
          nanos: 0,
        },
      });
    });

    it("should enable iframe when security settings allow", async () => {
      vi.mocked(zitadelModule.getSecuritySettings).mockResolvedValue({
        embeddedIframe: { enabled: true },
      } as any);
      vi.mocked(zitadelModule.createSessionFromChecks).mockResolvedValue({
        sessionId: mockSessionId,
        sessionToken: mockSessionToken,
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSessionResponse);
      vi.mocked(cookiesModule.addSessionToCookie).mockResolvedValue(undefined);

      await createSessionAndUpdateCookie({
        checks: mockChecks,
        requestId: undefined,
        lifetime: mockLifetime,
      });

      expect(cookiesModule.addSessionToCookie).toHaveBeenCalledWith(
        expect.objectContaining({
          iFrameEnabled: true,
        }),
      );
    });

    it("should throw error when session creation fails", async () => {
      vi.mocked(zitadelModule.createSessionFromChecks).mockResolvedValue(null as any);

      await expect(
        createSessionAndUpdateCookie({
          checks: mockChecks,
          requestId: undefined,
          lifetime: mockLifetime,
        }),
      ).rejects.toThrow("Could not create session");
    });

    it("should throw error when session has no loginName", async () => {
      vi.mocked(zitadelModule.createSessionFromChecks).mockResolvedValue({
        sessionId: mockSessionId,
        sessionToken: mockSessionToken,
      } as any);

      const invalidSession = {
        session: {
          id: mockSessionId,
          factors: {},
        },
      } as GetSessionResponse;

      vi.mocked(zitadelModule.getSession).mockResolvedValue(invalidSession);

      await expect(
        createSessionAndUpdateCookie({
          checks: mockChecks,
          requestId: undefined,
          lifetime: mockLifetime,
        }),
      ).rejects.toThrow("could not get session or session does not have loginName");
    });
  });

  describe("createSessionForIdpAndUpdateCookie", () => {
    const mockIdpIntent = {
      idpIntentId: "intent123",
      idpIntentToken: "intentToken123",
    };

    it("should create IDP session and update cookie successfully", async () => {
      vi.mocked(zitadelModule.createSessionForUserIdAndIdpIntent).mockResolvedValue({
        sessionId: mockSessionId,
        sessionToken: mockSessionToken,
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSessionResponse);
      vi.mocked(cookiesModule.addSessionToCookie).mockResolvedValue(undefined);

      const result = await createSessionForIdpAndUpdateCookie({
        userId: mockUserId,
        idpIntent: mockIdpIntent,
        requestId: "oidc_request123",
        lifetime: { seconds: BigInt(3600), nanos: 0 } as Duration,
      });

      expect(result).toEqual(mockSession);
      expect(zitadelModule.createSessionForUserIdAndIdpIntent).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
        idpIntent: mockIdpIntent,
        lifetime: { seconds: BigInt(3600), nanos: 0 },
      });
    });

    it("should use default lifetime when not provided for IDP", async () => {
      vi.mocked(zitadelModule.createSessionForUserIdAndIdpIntent).mockResolvedValue({
        sessionId: mockSessionId,
        sessionToken: mockSessionToken,
      } as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSessionResponse);
      vi.mocked(cookiesModule.addSessionToCookie).mockResolvedValue(undefined);

      await createSessionForIdpAndUpdateCookie({
        userId: mockUserId,
        idpIntent: mockIdpIntent,
        requestId: undefined,
      });

      expect(zitadelModule.createSessionForUserIdAndIdpIntent).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        userId: mockUserId,
        idpIntent: mockIdpIntent,
        lifetime: {
          seconds: BigInt(24 * 60 * 60),
          nanos: 0,
        },
      });
    });

    it("should throw error when IDP session creation fails", async () => {
      vi.mocked(zitadelModule.createSessionForUserIdAndIdpIntent).mockResolvedValue(null as any);

      await expect(
        createSessionForIdpAndUpdateCookie({
          userId: mockUserId,
          idpIntent: mockIdpIntent,
          requestId: undefined,
        }),
      ).rejects.toThrow("Could not create session");
    });

    it("should handle password attempt errors from IDP session", async () => {
      const error = {
        failedAttempts: 3,
      };

      vi.mocked(zitadelModule.createSessionForUserIdAndIdpIntent).mockRejectedValue(error);

      await expect(
        createSessionForIdpAndUpdateCookie({
          userId: mockUserId,
          idpIntent: mockIdpIntent,
          requestId: undefined,
        }),
      ).rejects.toMatchObject({
        error: "Failed to authenticate: You had 3 password attempts.",
        failedAttempts: 3,
      });
    });
  });

  describe("setSessionAndUpdateCookie", () => {
    const mockRecentCookie = {
      id: mockSessionId,
      token: mockSessionToken,
      loginName: mockLoginName,
      organization: mockOrganization,
      creationTs: "1234567890",
      expirationTs: "1234567890",
      changeTs: "1234567890",
    };

    const mockLifetime = {
      seconds: BigInt(3600),
      nanos: 0,
    } as Duration;

    it("should update session and cookie successfully", async () => {
      const mockUpdatedSession = {
        sessionToken: "newToken456",
        details: {
          changeDate: mockCreationDate,
        },
        challenges: undefined,
      };

      vi.mocked(zitadelModule.setSession).mockResolvedValue(mockUpdatedSession as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSessionResponse);
      vi.mocked(cookiesModule.updateSessionCookie).mockResolvedValue(undefined);

      const result = await setSessionAndUpdateCookie({
        recentCookie: mockRecentCookie,
        lifetime: mockLifetime,
      });

      expect(result).toMatchObject({
        ...mockSession,
        challenges: undefined,
      });
      expect(zitadelModule.setSession).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        sessionId: mockSessionId,
        sessionToken: mockSessionToken,
        challenges: undefined,
        checks: undefined,
        lifetime: mockLifetime,
      });
      expect(cookiesModule.updateSessionCookie).toHaveBeenCalledWith({
        id: mockSessionId,
        session: expect.objectContaining({
          id: mockSessionId,
          token: "newToken456",
          loginName: mockLoginName,
        }),
        iFrameEnabled: false,
      });
    });

    it("should include requestId when provided", async () => {
      const mockUpdatedSession = {
        sessionToken: "newToken456",
        details: { changeDate: mockCreationDate },
        challenges: undefined,
      };

      vi.mocked(zitadelModule.setSession).mockResolvedValue(mockUpdatedSession as any);
      vi.mocked(zitadelModule.getSession).mockResolvedValue(mockSessionResponse);
      vi.mocked(cookiesModule.updateSessionCookie).mockResolvedValue(undefined);

      await setSessionAndUpdateCookie({
        recentCookie: mockRecentCookie,
        requestId: "oidc_request999",
        lifetime: mockLifetime,
      });

      expect(cookiesModule.updateSessionCookie).toHaveBeenCalledWith({
        id: mockSessionId,
        session: expect.objectContaining({
          requestId: "oidc_request999",
        }),
        iFrameEnabled: false,
      });
    });

    it("should throw error when session update fails", async () => {
      vi.mocked(zitadelModule.setSession).mockResolvedValue(null as any);

      await expect(
        setSessionAndUpdateCookie({
          recentCookie: mockRecentCookie,
          lifetime: mockLifetime,
        }),
      ).rejects.toThrow();
    });

    it("should handle password attempt errors", async () => {
      const error = {
        findDetails: vi.fn(() => [{ failedAttempts: 5 }]),
      };

      vi.mocked(zitadelModule.setSession).mockRejectedValue(error);

      await expect(
        setSessionAndUpdateCookie({
          recentCookie: mockRecentCookie,
          lifetime: mockLifetime,
        }),
      ).rejects.toMatchObject({
        error: "Failed to authenticate: You had 5 password attempts.",
        failedAttempts: 5,
      });
    });
  });
});

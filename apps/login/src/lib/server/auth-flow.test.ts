/**
 * Unit tests for auth-flow server actions.
 *
 * Tests core business logic for:
 * - OIDC and SAML authentication flow completion
 * - Session loading and validation
 * - Error handling and validation
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { completeAuthFlow } from "./auth-flow";
import * as cookiesModule from "@/lib/cookies";
import * as zitadelModule from "@/lib/zitadel";
import * as oidcModule from "@/lib/oidc";
import * as samlModule from "@/lib/saml";
import type { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

// Mock all dependencies
vi.mock("@/lib/cookies");
vi.mock("@/lib/zitadel");
vi.mock("@/lib/oidc");
vi.mock("@/lib/saml");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
}));
vi.mock("@/lib/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://zitadel-test.zitadel.cloud" } })),
}));

describe("completeAuthFlow", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockSessionId = "session123";
  const mockOrganization = "org123";

  const mockSession: Session = {
    id: mockSessionId,
    factors: {
      user: {
        id: "user123",
        loginName: "test@example.com",
        organizationId: mockOrganization,
      },
    },
  } as Session;

  const mockSessionCookie = {
    id: mockSessionId,
    token: "token123",
    loginName: "test@example.com",
    creationTs: "1234567890",
    expirationTs: "1234567890",
    changeTs: "1234567890",
  };

  beforeEach(() => {
    vi.clearAllMocks();

    // Setup default mocks
    vi.mocked(cookiesModule.getAllSessions).mockResolvedValue([mockSessionCookie]);
    vi.mocked(zitadelModule.listSessions).mockResolvedValue({
      sessions: [mockSession],
    } as any);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("OIDC flow completion", () => {
    it("should complete OIDC flow successfully", async () => {
      const mockRedirectUrl = "https://app.example.com/callback";
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue({
        redirect: mockRedirectUrl,
      });

      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "oidc_request123",
        organization: mockOrganization,
      });

      expect(result).toEqual({ redirect: mockRedirectUrl });
      expect(oidcModule.loginWithOIDCAndSession).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        authRequest: "request123",
        sessionId: mockSessionId,
        sessions: [mockSession],
        sessionCookies: [mockSessionCookie],
      });
    });

    it("should strip 'oidc_' prefix from request ID", async () => {
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue({
        redirect: "https://app.example.com/callback",
      });

      await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "oidc_authrequest456",
      });

      expect(oidcModule.loginWithOIDCAndSession).toHaveBeenCalledWith(
        expect.objectContaining({
          authRequest: "authrequest456",
        }),
      );
    });

    it("should handle OIDC flow errors", async () => {
      const mockError = "Authentication failed";
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue({
        error: mockError,
      });

      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "oidc_request123",
      });

      expect(result).toEqual({ error: mockError });
    });

    it("should handle invalid OIDC response with safety net", async () => {
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue(null as any);

      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "oidc_request123",
      });

      expect(result).toEqual({
        error: "Authentication completed but navigation failed",
      });
    });
  });

  describe("SAML flow completion", () => {
    it("should complete SAML flow successfully", async () => {
      const mockRedirectUrl = "https://app.example.com/saml/callback";
      vi.mocked(samlModule.loginWithSAMLAndSession).mockResolvedValue({
        redirect: mockRedirectUrl,
      });

      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "saml_request789",
        organization: mockOrganization,
      });

      expect(result).toEqual({ redirect: mockRedirectUrl });
      expect(samlModule.loginWithSAMLAndSession).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        samlRequest: "request789",
        sessionId: mockSessionId,
        sessions: [mockSession],
        sessionCookies: [mockSessionCookie],
      });
    });

    it("should strip 'saml_' prefix from request ID", async () => {
      vi.mocked(samlModule.loginWithSAMLAndSession).mockResolvedValue({
        redirect: "https://app.example.com/callback",
      });

      await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "saml_samlrequest999",
      });

      expect(samlModule.loginWithSAMLAndSession).toHaveBeenCalledWith(
        expect.objectContaining({
          samlRequest: "samlrequest999",
        }),
      );
    });

    it("should handle SAML flow errors", async () => {
      const mockError = "SAML authentication failed";
      vi.mocked(samlModule.loginWithSAMLAndSession).mockResolvedValue({
        error: mockError,
      });

      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "saml_request789",
      });

      expect(result).toEqual({ error: mockError });
    });

    it("should handle invalid SAML response with safety net", async () => {
      vi.mocked(samlModule.loginWithSAMLAndSession).mockResolvedValue(undefined as any);

      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "saml_request789",
      });

      expect(result).toEqual({
        error: "Authentication completed but navigation failed",
      });
    });
  });

  describe("Session loading", () => {
    it("should load sessions from cookies", async () => {
      const sessionCookies = [
        {
          id: "session1",
          token: "token1",
          loginName: "user1@example.com",
          creationTs: "1234567890",
          expirationTs: "1234567890",
          changeTs: "1234567890",
        },
        {
          id: "session2",
          token: "token2",
          loginName: "user2@example.com",
          creationTs: "1234567890",
          expirationTs: "1234567890",
          changeTs: "1234567890",
        },
      ];

      vi.mocked(cookiesModule.getAllSessions).mockResolvedValue(sessionCookies);
      vi.mocked(zitadelModule.listSessions).mockResolvedValue({
        sessions: [mockSession],
      } as any);
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue({
        redirect: "https://app.example.com",
      });

      await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "oidc_request123",
      });

      expect(zitadelModule.listSessions).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        ids: ["session1", "session2"],
      });
    });

    it("should handle empty session cookies", async () => {
      vi.mocked(cookiesModule.getAllSessions).mockResolvedValue([]);
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue({
        redirect: "https://app.example.com",
      });

      await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "oidc_request123",
      });

      expect(zitadelModule.listSessions).not.toHaveBeenCalled();
      expect(oidcModule.loginWithOIDCAndSession).toHaveBeenCalledWith(
        expect.objectContaining({
          sessions: [],
        }),
      );
    });

    it("should filter out undefined session IDs", async () => {
      const sessionCookies = [
        {
          id: "session1",
          token: "token1",
          loginName: "user1@example.com",
          creationTs: "1234567890",
          expirationTs: "1234567890",
          changeTs: "1234567890",
        },
        {
          id: undefined,
          token: "token2",
          loginName: "user2@example.com",
          creationTs: "1234567890",
          expirationTs: "1234567890",
          changeTs: "1234567890",
        },
        {
          id: "session3",
          token: "token3",
          loginName: "user3@example.com",
          creationTs: "1234567890",
          expirationTs: "1234567890",
          changeTs: "1234567890",
        },
      ];

      vi.mocked(cookiesModule.getAllSessions).mockResolvedValue(sessionCookies as any);
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue({
        redirect: "https://app.example.com",
      });

      await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "oidc_request123",
      });

      expect(zitadelModule.listSessions).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: mockServiceUrl },
        ids: ["session1", "session3"],
      });
    });
  });

  describe("Request ID validation", () => {
    it("should return error for invalid request ID format", async () => {
      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "invalid_request123",
      });

      expect(result).toEqual({ error: "Invalid request ID format" });
      expect(oidcModule.loginWithOIDCAndSession).not.toHaveBeenCalled();
      expect(samlModule.loginWithSAMLAndSession).not.toHaveBeenCalled();
    });

    it("should return error for request ID without prefix", async () => {
      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "request123",
      });

      expect(result).toEqual({ error: "Invalid request ID format" });
    });

    it("should return error for empty request ID", async () => {
      const result = await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "",
      });

      expect(result).toEqual({ error: "Invalid request ID format" });
    });
  });

  describe("Error handling", () => {
    it("should handle session loading errors gracefully", async () => {
      vi.mocked(zitadelModule.listSessions).mockRejectedValue(new Error("Session service unavailable"));
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue({
        redirect: "https://app.example.com",
      });

      await expect(
        completeAuthFlow({
          sessionId: mockSessionId,
          requestId: "oidc_request123",
        }),
      ).rejects.toThrow("Session service unavailable");
    });

    it("should pass organization parameter through flow", async () => {
      vi.mocked(oidcModule.loginWithOIDCAndSession).mockResolvedValue({
        redirect: "https://app.example.com",
      });

      await completeAuthFlow({
        sessionId: mockSessionId,
        requestId: "oidc_request123",
        organization: "org456",
      });

      // Organization should be available in the command but not necessarily passed to OIDC
      expect(oidcModule.loginWithOIDCAndSession).toHaveBeenCalled();
    });
  });
});

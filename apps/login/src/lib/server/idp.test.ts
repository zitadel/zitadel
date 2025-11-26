/**
 * Unit tests for the IDP server actions.
 *
 * These tests replace the integration tests from register-idp.cy.ts which tested:
 * - IDP redirect URL generation
 * - Starting identity provider flow
 * - Redirecting users to the correct IDP authentication URL
 */

import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { redirectToIdp } from "./idp";
import * as zitadelModule from "../zitadel";

// Mock all dependencies
vi.mock("../zitadel");
vi.mock("./cookie");
vi.mock("./host");
vi.mock("../client");
vi.mock("next/headers", () => ({
  headers: vi.fn(() => Promise.resolve(new Map())),
}));
vi.mock("next/navigation", () => ({
  redirect: vi.fn((url: string) => {
    // Throw error to simulate redirect behavior
    throw new Error(`REDIRECT: ${url}`);
  }),
}));
vi.mock("../service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(() => ({ serviceUrl: "https://zitadel-test.zitadel.cloud" })),
}));

describe("redirectToIdp server action", () => {
  const mockServiceUrl = "https://zitadel-test.zitadel.cloud";
  const mockIdpId = "idp-123";
  const mockIdpUrl = "https://accounts.google.com/oauth2/auth?client_id=...";
  const mockOrganization = "256088834543534543";
  const mockRequestId = "req123";

  beforeEach(async () => {
    vi.clearAllMocks();

    const { getOriginalHost } = await import("./host");
    vi.mocked(getOriginalHost).mockResolvedValue("localhost:3000");
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("IDP flow initiation", () => {
    it("should start IDP flow and redirect to IDP URL", async () => {
      // Mock IDP flow start
      vi.mocked(zitadelModule.startIdentityProviderFlow).mockResolvedValue(mockIdpUrl);

      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "google");
      formData.set("organization", mockOrganization);
      formData.set("requestId", mockRequestId);

      try {
        await redirectToIdp(undefined, formData);
        // Should not reach here
        expect(true).toBe(false);
      } catch (error: any) {
        // Should redirect
        expect(error.message).toContain("REDIRECT:");
        expect(error.message).toContain(mockIdpUrl);
      }

      expect(vi.mocked(zitadelModule.startIdentityProviderFlow)).toHaveBeenCalledWith({
        serviceUrl: mockServiceUrl,
        idpId: mockIdpId,
        urls: expect.objectContaining({
          successUrl: expect.stringContaining("/idp/google/success"),
          failureUrl: expect.stringContaining("/idp/google/failure"),
        }),
      });
    });

    it("should include query parameters in success and failure URLs", async () => {
      vi.mocked(zitadelModule.startIdentityProviderFlow).mockResolvedValue(mockIdpUrl);

      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "google");
      formData.set("organization", mockOrganization);
      formData.set("requestId", mockRequestId);

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Expected redirect
        expect(error.message).toBeDefined();
      }

      const callArgs = vi.mocked(zitadelModule.startIdentityProviderFlow).mock.calls[0][0];
      expect(callArgs.urls.successUrl).toContain(`organization=${mockOrganization}`);
      expect(callArgs.urls.successUrl).toContain(`requestId=${mockRequestId}`);
      expect(callArgs.urls.failureUrl).toContain(`organization=${mockOrganization}`);
      expect(callArgs.urls.failureUrl).toContain(`requestId=${mockRequestId}`);
    });

    it("should handle LDAP provider differently by redirecting to LDAP page", async () => {
      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "ldap");
      formData.set("organization", mockOrganization);

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        expect(error.message).toContain("REDIRECT:");
        expect(error.message).toContain("/idp/ldap");
        expect(error.message).toContain(`idpId=${mockIdpId}`);
      }

      // Should not call startIdentityProviderFlow for LDAP
      expect(vi.mocked(zitadelModule.startIdentityProviderFlow)).not.toHaveBeenCalled();
    });

    it("should handle linkOnly parameter", async () => {
      vi.mocked(zitadelModule.startIdentityProviderFlow).mockResolvedValue(mockIdpUrl);

      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "github");
      formData.set("linkOnly", "true");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Expected redirect
        expect(error.message).toBeDefined();
      }

      const callArgs = vi.mocked(zitadelModule.startIdentityProviderFlow).mock.calls[0][0];
      expect(callArgs.urls.successUrl).toContain("link=true");
    });
  });

  describe("error handling", () => {
    it("should return error when IDP flow cannot be started", async () => {
      // Mock IDP flow start failure
      vi.mocked(zitadelModule.startIdentityProviderFlow).mockResolvedValue(null);

      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "google");

      const result = await redirectToIdp(undefined, formData);

      expect(result).toHaveProperty("error");
      expect(result?.error).toBe("Unexpected response from IDP flow");
    });

    it("should return error when IDP flow returns unexpected response", async () => {
      // Mock unexpected response (empty string instead of URL)
      vi.mocked(zitadelModule.startIdentityProviderFlow).mockResolvedValue("");

      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "google");

      const result = await redirectToIdp(undefined, formData);

      expect(result).toHaveProperty("error");
      expect(result?.error).toBe("Unexpected response from IDP flow");
    });
  });

  describe("different IDP providers", () => {
    it("should handle Google IDP", async () => {
      vi.mocked(zitadelModule.startIdentityProviderFlow).mockResolvedValue("https://accounts.google.com/oauth2/auth");

      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "google");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        expect(error.message).toContain("accounts.google.com");
      }
    });

    it("should handle GitHub IDP", async () => {
      vi.mocked(zitadelModule.startIdentityProviderFlow).mockResolvedValue("https://github.com/login/oauth/authorize");

      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "github");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        expect(error.message).toContain("github.com");
      }
    });

    it("should handle GitLab IDP", async () => {
      vi.mocked(zitadelModule.startIdentityProviderFlow).mockResolvedValue("https://gitlab.com/oauth/authorize");

      const formData = new FormData();
      formData.set("id", mockIdpId);
      formData.set("provider", "gitlab");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        expect(error.message).toContain("gitlab.com");
      }
    });
  });
});

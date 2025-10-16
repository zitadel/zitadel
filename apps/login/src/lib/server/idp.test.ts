import { describe, expect, test, vi, beforeEach, afterEach } from "vitest";
import { redirectToIdp } from "./idp";

// Mock all the dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("next/navigation", () => ({
  redirect: vi.fn((url: string) => {
    throw new Error(`REDIRECT: ${url}`);
  }),
}));

vi.mock("../service-url", () => ({
  getServiceUrlFromHeaders: vi.fn(),
}));

vi.mock("./host", () => ({
  getOriginalHost: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  startIdentityProviderFlow: vi.fn(),
}));

describe("redirectToIdp", () => {
  let mockHeaders: any;
  let mockGetServiceUrlFromHeaders: any;
  let mockGetOriginalHost: any;
  let mockStartIdentityProviderFlow: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    // Import mocked modules
    const { headers } = await import("next/headers");
    const { getServiceUrlFromHeaders } = await import("../service-url");
    const { getOriginalHost } = await import("./host");
    const { startIdentityProviderFlow } = await import("../zitadel");

    // Setup mocks
    mockHeaders = vi.mocked(headers);
    mockGetServiceUrlFromHeaders = vi.mocked(getServiceUrlFromHeaders);
    mockGetOriginalHost = vi.mocked(getOriginalHost);
    mockStartIdentityProviderFlow = vi.mocked(startIdentityProviderFlow);

    // Default mock implementations
    mockHeaders.mockResolvedValue({} as any);
    mockGetServiceUrlFromHeaders.mockReturnValue({
      serviceUrl: "https://api.example.com",
    });
    mockGetOriginalHost.mockResolvedValue("example.com");
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("postErrorRedirectUrl parameter handling", () => {
    test("should include postErrorRedirectUrl in success and failure URLs when provided", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");
      formData.append("requestId", "req123");
      formData.append("organization", "org123");
      formData.append("postErrorRedirectUrl", "https://app.example.com/error");

      mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT:");
      }

      expect(mockStartIdentityProviderFlow).toHaveBeenCalledWith({
        serviceUrl: "https://api.example.com",
        idpId: "idp123",
        urls: {
          successUrl: expect.stringContaining("postErrorRedirectUrl=https%3A%2F%2Fapp.example.com%2Ferror"),
          failureUrl: expect.stringContaining("postErrorRedirectUrl=https%3A%2F%2Fapp.example.com%2Ferror"),
        },
      });

      // Verify both URLs contain all expected parameters
      const callArgs = mockStartIdentityProviderFlow.mock.calls[0][0];
      const successUrl = callArgs.urls.successUrl;
      const failureUrl = callArgs.urls.failureUrl;

      expect(successUrl).toContain("requestId=req123");
      expect(successUrl).toContain("organization=org123");
      expect(successUrl).toContain("postErrorRedirectUrl=https%3A%2F%2Fapp.example.com%2Ferror");

      expect(failureUrl).toContain("requestId=req123");
      expect(failureUrl).toContain("organization=org123");
      expect(failureUrl).toContain("postErrorRedirectUrl=https%3A%2F%2Fapp.example.com%2Ferror");
    });

    test("should not include postErrorRedirectUrl in URLs when not provided", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");
      formData.append("requestId", "req123");
      formData.append("organization", "org123");

      mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT:");
      }

      expect(mockStartIdentityProviderFlow).toHaveBeenCalledWith({
        serviceUrl: "https://api.example.com",
        idpId: "idp123",
        urls: {
          successUrl: expect.not.stringContaining("postErrorRedirectUrl"),
          failureUrl: expect.not.stringContaining("postErrorRedirectUrl"),
        },
      });

      const callArgs = mockStartIdentityProviderFlow.mock.calls[0][0];
      expect(callArgs.urls.successUrl).not.toContain("postErrorRedirectUrl");
      expect(callArgs.urls.failureUrl).not.toContain("postErrorRedirectUrl");
    });

    test("should not include postErrorRedirectUrl when it is an empty string", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");
      formData.append("postErrorRedirectUrl", "");

      mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT:");
      }

      const callArgs = mockStartIdentityProviderFlow.mock.calls[0][0];
      expect(callArgs.urls.successUrl).not.toContain("postErrorRedirectUrl");
      expect(callArgs.urls.failureUrl).not.toContain("postErrorRedirectUrl");
    });

    test("should include postErrorRedirectUrl in LDAP redirect URL", async () => {
      const formData = new FormData();
      formData.append("id", "ldap123");
      formData.append("provider", "ldap");
      formData.append("requestId", "req123");
      formData.append("organization", "org123");
      formData.append("postErrorRedirectUrl", "/custom-error");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT: /idp/ldap?");
        expect(error.message).toContain("requestId=req123");
        expect(error.message).toContain("organization=org123");
        expect(error.message).toContain("postErrorRedirectUrl=%2Fcustom-error");
      }
    });

    test("should handle postErrorRedirectUrl with special characters", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");
      formData.append("postErrorRedirectUrl", "https://app.example.com/error?code=123&message=test");

      mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT:");
      }

      const callArgs = mockStartIdentityProviderFlow.mock.calls[0][0];
      const successUrl = new URL(callArgs.urls.successUrl);
      const failureUrl = new URL(callArgs.urls.failureUrl);

      // Verify the URL is properly encoded
      expect(successUrl.searchParams.get("postErrorRedirectUrl")).toBe(
        "https://app.example.com/error?code=123&message=test",
      );
      expect(failureUrl.searchParams.get("postErrorRedirectUrl")).toBe(
        "https://app.example.com/error?code=123&message=test",
      );
    });

    test("should preserve postErrorRedirectUrl alongside linkOnly parameter", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");
      formData.append("linkOnly", "true");
      formData.append("postErrorRedirectUrl", "/custom-error");

      mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT:");
      }

      const callArgs = mockStartIdentityProviderFlow.mock.calls[0][0];
      const successUrl = callArgs.urls.successUrl;
      const failureUrl = callArgs.urls.failureUrl;

      // Both parameters should be present
      expect(successUrl).toContain("link=true");
      expect(successUrl).toContain("postErrorRedirectUrl=%2Fcustom-error");
      expect(failureUrl).toContain("link=true");
      expect(failureUrl).toContain("postErrorRedirectUrl=%2Fcustom-error");
    });

    test("should handle relative postErrorRedirectUrl paths", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "github");
      formData.append("postErrorRedirectUrl", "/loginname");

      mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT:");
      }

      const callArgs = mockStartIdentityProviderFlow.mock.calls[0][0];
      expect(callArgs.urls.successUrl).toContain("postErrorRedirectUrl=%2Floginname");
      expect(callArgs.urls.failureUrl).toContain("postErrorRedirectUrl=%2Floginname");
    });
  });

  describe("General redirect behavior", () => {
    test("should return error when IDP flow returns null", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");

      mockStartIdentityProviderFlow.mockResolvedValue(null);

      const result = await redirectToIdp(undefined, formData);

      expect(result).toEqual({ error: "Unexpected response from IDP flow" });
    });

    test("should redirect when IDP flow returns a valid URL", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");

      mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

      try {
        await redirectToIdp(undefined, formData);
        // Should not reach here
        expect(true).toBe(false);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toBe("REDIRECT: https://idp.example.com/auth");
      }
    });
  });
});

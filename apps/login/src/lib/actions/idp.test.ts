import { createNewSessionForLDAP, redirectToIdp } from "@/lib/actions/idp";
import { createIdpFlowProof, isIdpFlowProofValid } from "@/lib/server/idp-flow-proof";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";

// Mock all the dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("next/navigation", () => ({
  redirect: vi.fn((url: string) => {
    throw new Error(`REDIRECT: ${url}`);
  }),
}));

vi.mock("@/lib/server/cookies", () => ({
  getSessionCookieById: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("@/lib/server/host", () => ({
  getInstanceHost: vi.fn(),
  getPublicHost: vi.fn(),
}));

vi.mock("@/lib/server/zitadel", () => ({
  startIdentityProviderFlow: vi.fn(),
  startLDAPIdentityProviderFlow: vi.fn(),
}));

vi.mock("@/lib/server/fingerprint", () => ({
  getBindingIdCookie: vi.fn(),
  addIdpFlowProof: vi.fn(),
}));

// the real helper reads the binding id cookie, here the browser binding id is "fp123"
async function mockAddIdpFlowProof() {
  const { addIdpFlowProof } = await import("@/lib/server/fingerprint");
  vi.mocked(addIdpFlowProof).mockImplementation(async (params, linkToSessionId) => {
    params.set("flowFingerprint", createIdpFlowProof("fp123", linkToSessionId));
  });
}

describe("redirectToIdp", () => {
  let mockHeaders: any;
  let mockGetServiceUrlFromHeaders: any;
  let mockGetInstanceHost: any;
  let mockGetPublicHost: any;
  let mockStartIdentityProviderFlow: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");

    // Import mocked modules
    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("@/lib/server/service-url");
    const { getInstanceHost, getPublicHost } = await import("@/lib/server/host");
    const { startIdentityProviderFlow } = await import("@/lib/server/zitadel");
    const { getBindingIdCookie } = await import("@/lib/server/fingerprint");
    const { getSessionCookieById } = await import("@/lib/server/cookies");

    // Setup mocks
    mockHeaders = vi.mocked(headers);
    mockGetServiceUrlFromHeaders = vi.mocked(getServiceConfig);
    mockGetInstanceHost = vi.mocked(getInstanceHost);
    mockGetPublicHost = vi.mocked(getPublicHost);
    mockStartIdentityProviderFlow = vi.mocked(startIdentityProviderFlow);
    const mockGetBindingIdCookie = vi.mocked(getBindingIdCookie);
    const mockGetSessionCookieById = vi.mocked(getSessionCookieById);

    // Default mock implementations
    mockGetBindingIdCookie.mockResolvedValue({ name: "bindingId", value: "fp123" });
    await mockAddIdpFlowProof();
    mockGetSessionCookieById.mockResolvedValue({} as any);
    mockHeaders.mockResolvedValue({} as any);
    mockGetServiceUrlFromHeaders.mockReturnValue({
      serviceConfig: { baseUrl: "https://api.example.com" },
    });
    mockGetInstanceHost.mockReturnValue("example.com");
    mockGetPublicHost.mockReturnValue("example.com");
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllEnvs();
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
        serviceConfig: { baseUrl: "https://api.example.com" },
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

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT:");
      }

      expect(mockStartIdentityProviderFlow).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
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

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

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

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

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
      formData.append("sessionId", "session123");
      formData.append("postErrorRedirectUrl", "/custom-error");

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

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
      expect(successUrl).toContain("linkToSessionId=session123");
      expect(successUrl).toContain("postErrorRedirectUrl=%2Fcustom-error");
      expect(failureUrl).toContain("linkToSessionId=session123");
      expect(failureUrl).toContain("postErrorRedirectUrl=%2Fcustom-error");
    });

    test("should handle relative postErrorRedirectUrl paths", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "github");
      formData.append("postErrorRedirectUrl", "/loginname");

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

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

  describe("loginHint parameter handling", () => {
    test("should forward the loginHint to the IDP flow when provided", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");
      formData.append("loginHint", "user@example.com");

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

      try {
        await redirectToIdp(undefined, formData);
      } catch (error: any) {
        // Redirect throws in tests
        expect(error.message).toContain("REDIRECT:");
      }

      expect(mockStartIdentityProviderFlow).toHaveBeenCalledWith(
        expect.objectContaining({
          idpId: "idp123",
          urls: expect.objectContaining({ loginHint: "user@example.com" }),
        }),
      );
    });

    test("should not set a loginHint when it is missing or empty", async () => {
      for (const loginHint of [undefined, ""]) {
        mockStartIdentityProviderFlow.mockReset();
        mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

        const formData = new FormData();
        formData.append("id", "idp123");
        formData.append("provider", "google");
        if (loginHint !== undefined) {
          formData.append("loginHint", loginHint);
        }

        try {
          await redirectToIdp(undefined, formData);
        } catch (error: any) {
          // Redirect throws in tests
          expect(error.message).toContain("REDIRECT:");
        }

        const callArgs = mockStartIdentityProviderFlow.mock.calls[0][0];
        expect(callArgs.urls.loginHint).toBeUndefined();
      }
    });
  });

  describe("General redirect behavior", () => {
    test("should return error when IDP flow returns null", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");

      mockStartIdentityProviderFlow.mockResolvedValue(null);

      const result = await redirectToIdp(undefined, formData);

      expect(result).toEqual({ error: "Could not start IDP flow" });
    });

    test("should redirect when IDP flow returns a valid URL", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

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

  describe("flowFingerprint (browser binding)", () => {
    test("should add a flow proof bound to this browser's binding id to the success URL", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

      await expect(redirectToIdp(undefined, formData)).rejects.toThrow("REDIRECT:");

      const successUrl = mockStartIdentityProviderFlow.mock.calls[0][0].urls.successUrl;
      const flowFingerprint = new URL(successUrl).searchParams.get("flowFingerprint") ?? undefined;
      expect(isIdpFlowProofValid(flowFingerprint, "fp123")).toBe(true);
      expect(isIdpFlowProofValid(flowFingerprint, "other-fp")).toBe(false);
    });

    test("should bind the flow proof to the session when linking", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");
      formData.append("sessionId", "session123");

      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

      await expect(redirectToIdp(undefined, formData)).rejects.toThrow("REDIRECT:");

      const successUrl = new URL(mockStartIdentityProviderFlow.mock.calls[0][0].urls.successUrl);
      const flowFingerprint = successUrl.searchParams.get("flowFingerprint") ?? undefined;
      expect(isIdpFlowProofValid(flowFingerprint, "fp123", "session123")).toBe(true);
      expect(isIdpFlowProofValid(flowFingerprint, "fp123", "other-session")).toBe(false);
      expect(isIdpFlowProofValid(flowFingerprint, "fp123")).toBe(false);
    });

    test("should reject linking to a session this browser does not own", async () => {
      const { getSessionCookieById } = await import("@/lib/server/cookies");
      vi.mocked(getSessionCookieById).mockResolvedValue(undefined);
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "google");
      formData.append("sessionId", "victim-session");

      const result = await redirectToIdp(undefined, formData);

      expect(result).toEqual({ error: "Invalid session for linking" });
      expect(mockStartIdentityProviderFlow).not.toHaveBeenCalled();
    });

    test("should not add a flow proof to the LDAP credentials page redirect", async () => {
      const formData = new FormData();
      formData.append("id", "idp123");
      formData.append("provider", "ldap");

      const error = await redirectToIdp(undefined, formData).catch((e: Error) => e);

      expect((error as Error).message).toContain("REDIRECT: /idp/ldap?");
      expect((error as Error).message).not.toContain("flowFingerprint");
    });
  });
});

describe("createNewSessionForLDAP", () => {
  beforeEach(async () => {
    vi.clearAllMocks();
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("@/lib/server/service-url");
    const { startLDAPIdentityProviderFlow } = await import("@/lib/server/zitadel");
    const { getSessionCookieById } = await import("@/lib/server/cookies");

    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } } as any);
    await mockAddIdpFlowProof();
    vi.mocked(getSessionCookieById).mockResolvedValue({} as any);
    vi.mocked(startLDAPIdentityProviderFlow).mockResolvedValue({
      nextStep: { case: "idpIntent", value: { userId: "user123", idpIntentId: "intent123", idpIntentToken: "token123" } },
    } as any);
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllEnvs();
  });

  test("should add a flow proof bound to this browser's binding id to the process URL", async () => {
    const result = await createNewSessionForLDAP({
      username: "user",
      password: "secret",
      idpId: "idp123",
      link: false,
      requestId: "req123",
    });

    expect("redirect" in result && result.redirect).toContain("/idp/ldap/process?");
    const params = new URLSearchParams((result as { redirect: string }).redirect.split("?")[1]);
    expect(params.get("id")).toBe("intent123");
    expect(params.get("requestId")).toBe("req123");
    expect(isIdpFlowProofValid(params.get("flowFingerprint") ?? undefined, "fp123")).toBe(true);
    expect(isIdpFlowProofValid(params.get("flowFingerprint") ?? undefined, "other-fp")).toBe(false);
  });

  test("should bind the flow proof to the session to link", async () => {
    const result = await createNewSessionForLDAP({
      username: "user",
      password: "secret",
      idpId: "idp123",
      link: true,
      linkToSessionId: "session123",
    });

    const params = new URLSearchParams((result as { redirect: string }).redirect.split("?")[1]);
    expect(params.get("linkToSessionId")).toBe("session123");
    expect(isIdpFlowProofValid(params.get("flowFingerprint") ?? undefined, "fp123", "session123")).toBe(true);
    expect(isIdpFlowProofValid(params.get("flowFingerprint") ?? undefined, "fp123")).toBe(false);
  });

  test("should reject linking to a session this browser does not own", async () => {
    const { getSessionCookieById } = await import("@/lib/server/cookies");
    const { startLDAPIdentityProviderFlow } = await import("@/lib/server/zitadel");
    vi.mocked(getSessionCookieById).mockResolvedValue(undefined);

    const result = await createNewSessionForLDAP({
      username: "user",
      password: "secret",
      idpId: "idp123",
      link: true,
      linkToSessionId: "victim-session",
    });

    expect(result).toEqual({ error: "Invalid session for linking" });
    expect(startLDAPIdentityProviderFlow).not.toHaveBeenCalled();
  });
});

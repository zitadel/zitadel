import { describe, expect, test, vi, beforeEach, afterEach } from "vitest";
import { sendLoginname } from "./loginname";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getIDPByID } from "../zitadel";

// Mock all the dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("../idp", () => ({
  idpTypeToIdentityProviderType: vi.fn(),
  idpTypeToSlug: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  getActiveIdentityProviders: vi.fn(),
  getIDPByID: vi.fn(),
  getLoginSettings: vi.fn(),
  getOrgsByDomain: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  listIDPLinks: vi.fn(),
  searchUsers: vi.fn(),
  startIdentityProviderFlow: vi.fn(),
}));

vi.mock("./cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("./host", () => ({
  getInstanceHost: vi.fn(),
  getPublicHost: vi.fn(),
}));

// this returns the key itself that can be checked not the translated value
vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

describe("sendLoginname", () => {
  // Mock modules
  let mockHeaders: any;
  let mockCreate: any;
  let mockGetServiceUrlFromHeaders: any;
  let mockGetLoginSettings: any;
  let mockSearchUsers: any;
  let mockCreateSessionAndUpdateCookie: any;
  let mockListAuthenticationMethodTypes: any;
  let mockListIDPLinks: any;
  let mockGetInstanceHost: any;
  let mockGetPublicHost: any;
  let mockStartIdentityProviderFlow: any;
  let mockGetActiveIdentityProviders: any;
  let mockGetIDPByID: any;
  let mockIdpTypeToSlug: any;
  let mockGetOrgsByDomain: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    // Import mocked modules
    const { headers } = await import("next/headers");
    const { create } = await import("@zitadel/client");
    const { getServiceConfig } = await import("../service-url");
    const {
      getLoginSettings,
      searchUsers,
      listAuthenticationMethodTypes,
      listIDPLinks,
      startIdentityProviderFlow,
      getActiveIdentityProviders,
      getOrgsByDomain,
    } = await import("../zitadel");
    const { createSessionAndUpdateCookie } = await import("./cookie");
    const { getInstanceHost, getPublicHost } = await import("./host");
    const { idpTypeToSlug } = await import("../idp");

    // Setup mocks
    mockHeaders = vi.mocked(headers);
    mockCreate = vi.mocked(create);
    mockGetServiceUrlFromHeaders = vi.mocked(getServiceConfig);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockSearchUsers = vi.mocked(searchUsers);
    mockCreateSessionAndUpdateCookie = vi.mocked(createSessionAndUpdateCookie);
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockListIDPLinks = vi.mocked(listIDPLinks);
    mockGetInstanceHost = vi.mocked(getInstanceHost);
    mockGetPublicHost = vi.mocked(getPublicHost);
    mockStartIdentityProviderFlow = vi.mocked(startIdentityProviderFlow);
    mockGetActiveIdentityProviders = vi.mocked(getActiveIdentityProviders);
    mockGetIDPByID = vi.mocked(getIDPByID);
    mockIdpTypeToSlug = vi.mocked(idpTypeToSlug);
    mockGetOrgsByDomain = vi.mocked(getOrgsByDomain);

    // Default mock implementations
    mockHeaders.mockResolvedValue({} as any);
    mockGetServiceUrlFromHeaders.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
    mockGetInstanceHost.mockReturnValue("example.com");
    mockGetPublicHost.mockReturnValue("example.com");
    mockIdpTypeToSlug.mockReturnValue("google");
    mockGetIDPByID.mockResolvedValue({
      id: "idp123",
      name: "Google",
      type: "GOOGLE",
    });
    // Default: org discovery returns empty result
    mockGetOrgsByDomain.mockResolvedValue({ result: [] });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("Error cases", () => {
    test("should return error when login settings cannot be retrieved", async () => {
      mockGetLoginSettings.mockResolvedValue(null);

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "errors.couldNotGetLoginSettings" });
    });

    test("should return error when user search fails", async () => {
      mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: true });
      mockSearchUsers.mockResolvedValue({ error: "Search failed" });

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "Search failed" });
    });

    test("should return error when search result has no result field", async () => {
      mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: true });
      mockSearchUsers.mockResolvedValue({});

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "errors.couldNotSearchUsers" });
    });

    test("should return error when more than one user found", async () => {
      mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: true });
      mockSearchUsers.mockResolvedValue({
        result: [
          { userId: "user1", preferredLoginName: "user1@example.com" },
          { userId: "user2", preferredLoginName: "user2@example.com" },
        ],
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "errors.moreThanOneUserFound" });
    });
  });

  describe("Single user found - authentication method handling", () => {
    const mockUser = {
      userId: "user123",
      preferredLoginName: "user@example.com",
      details: { resourceOwner: "org123" },
      type: { case: "human", value: { email: { email: "user@example.com" } } },
      state: UserState.ACTIVE,
    };

    const mockSession = {
      factors: {
        user: {
          id: "user123",
          loginName: "user@example.com",
          organizationId: "org123",
        },
      },
    };

    beforeEach(() => {
      mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: true });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
    });

    test("should redirect to verify when user has no authentication methods", async () => {
      mockListAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [] });

      const result = await sendLoginname({
        loginName: "user@example.com",
        requestId: "req123",
      });

      expect(result).toHaveProperty("redirect");
      expect((result as any).redirect).toMatch(/^\/verify\?/);
      expect((result as any).redirect).toContain("loginName=user%40example.com");
      expect((result as any).redirect).toContain("send=true");
      expect((result as any).redirect).toContain("invite=true");
      expect((result as any).redirect).toContain("requestId=req123");
    });

    describe("Single authentication method", () => {
      test("should redirect to password when user has only password method and it's allowed", async () => {
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
          requestId: "req123",
        });

        expect(result).toHaveProperty("redirect");
        expect((result as any).redirect).toMatch(/^\/password\?/);
        expect((result as any).redirect).toContain("loginName=user%40example.com");
        expect((result as any).redirect).toContain("requestId=req123");
      });

      test("should attempt IDP redirect when password is not allowed but user has IDP links", async () => {
        mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });
        mockListIDPLinks.mockResolvedValue({
          result: [{ idpId: "idp123" }],
        });
        mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({ redirect: "https://idp.example.com/auth" });
        expect(mockListIDPLinks).toHaveBeenCalledWith({
          serviceConfig: { baseUrl: "https://api.example.com" },
          userId: "user123",
        });
      });

      test("should return error when password not allowed and no IDP links available", async () => {
        mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });
        mockListIDPLinks.mockResolvedValue({ result: [] });
        mockGetActiveIdentityProviders.mockResolvedValue({ identityProviders: [] });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({
          error: "errors.localAuthenticationNotAllowed",
        });
      });

      test("should redirect to organization IDP when password not allowed, no user IDP links, but organization has active IDP", async () => {
        mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });
        mockListIDPLinks.mockResolvedValue({ result: [] });
        mockGetActiveIdentityProviders.mockResolvedValue({
          identityProviders: [{ id: "org-idp-123", type: 0 }],
        });
        mockIdpTypeToSlug.mockReturnValue("google");
        mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://org-idp.example.com/auth" });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({ redirect: "https://org-idp.example.com/auth" });
        expect(mockGetActiveIdentityProviders).toHaveBeenCalledWith({
          serviceConfig: { baseUrl: "https://api.example.com" },
          orgId: "org123", // User's organization from resourceOwner
        });
      });

      test("should redirect to passkey when user has only passkey method and it's allowed", async () => {
        mockGetLoginSettings.mockResolvedValue({ passkeysType: PasskeysType.ALLOWED, allowLocalAuthentication: true });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSKEY],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
          requestId: "req123",
        });

        expect(result).toHaveProperty("redirect");
        expect((result as any).redirect).toMatch(/^\/passkey\?/);
        expect((result as any).redirect).toContain("loginName=user%40example.com");
        expect((result as any).redirect).toContain("requestId=req123");
      });

      test("should return error when passkeys are not allowed", async () => {
        mockGetLoginSettings.mockResolvedValue({ passkeysType: PasskeysType.NOT_ALLOWED, allowLocalAuthentication: true });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSKEY],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({
          error: "errors.passkeysNotAllowed",
        });
      });

      test("should return error when passkeys are allowed but allowLocalAuthentication is false", async () => {
        mockGetLoginSettings.mockResolvedValue({
          passkeysType: PasskeysType.ALLOWED,
          allowLocalAuthentication: false,
        });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSKEY],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({
          error: "errors.passkeysNotAllowed",
        });
      });

      test("should redirect to IDP when user has only IDP method", async () => {
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.IDP],
        });
        mockListIDPLinks.mockResolvedValue({
          result: [{ idpId: "idp123" }],
        });
        mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({ redirect: "https://idp.example.com/auth" });
      });

      test("should NOT create session when ignoreUnknownUsernames is true", async () => {
        mockGetLoginSettings.mockResolvedValue({
          allowLocalAuthentication: true,
          ignoreUnknownUsernames: true,
        });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
          requestId: "req123",
        });

        expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();

        expect(result).toHaveProperty("redirect");
        expect((result as any).redirect).toMatch(/^\/password\?/);
        expect((result as any).redirect).toContain("loginName=user%40example.com");
      });
    });

    describe("Multiple authentication methods", () => {
      test("should prefer passkey when multiple methods available", async () => {
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.PASSKEY],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toHaveProperty("redirect");
        expect((result as any).redirect).toMatch(/^\/passkey\?/);
        expect((result as any).redirect).toContain("altPassword=true"); // password is allowed
      });

      test("should return error when allowLocalAuthentication is false (disabling both password and passkey)", async () => {
        mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.PASSKEY],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({
          error: "errors.localAuthenticationNotAllowed",
        });
      });

      test("should redirect to IDP when no passkey but IDP available", async () => {
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.IDP],
        });
        mockListIDPLinks.mockResolvedValue({
          result: [{ idpId: "idp123" }],
        });
        mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({ redirect: "https://idp.example.com/auth" });
      });

      test("should redirect to password when no passkey or IDP, only password available and allowed", async () => {
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toBeDefined();
        expect(result?.redirect).toMatch(/^\/password\?/);
      });

      test("should return error when password is only method in multi-method scenario but not allowed", async () => {
        mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });
        mockListIDPLinks.mockResolvedValue({ result: [] });
        mockGetActiveIdentityProviders.mockResolvedValue({ identityProviders: [] });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({
          error: "errors.localAuthenticationNotAllowed",
        });
      });
    });
  });

  describe("User not found scenarios", () => {
    beforeEach(() => {
      mockSearchUsers.mockResolvedValue({ result: [] });
    });

    test("should redirect to single IDP when register allowed but password not allowed", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowRegister: true,
        allowLocalAuthentication: false,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({
        identityProviders: [{ id: "idp123", type: "OIDC" }],
      });
      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth" });

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ redirect: "https://idp.example.com/auth" });
    });

    test("should redirect to register when both register and password allowed", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowRegister: true,
        allowLocalAuthentication: true,
        ignoreUnknownUsernames: false,
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        organization: "org123",
        requestId: "req123",
      });

      expect(result).toBeDefined();
      expect(result?.redirect).toMatch(/^\/register\?/);
      expect(result?.redirect).toContain("organization=org123");
      expect(result?.redirect).toContain("requestId=req123");
      expect(result?.redirect).toContain("email=user%40example.com");
    });

    test("should redirect to password when ignoreUnknownUsernames is true", async () => {
      mockGetLoginSettings.mockResolvedValue({
        ignoreUnknownUsernames: true,
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        requestId: "req123",
        organization: "org123",
        ignoreUnknownUsernames: true,
      });

      expect(result).toBeDefined();
      expect(result?.redirect).toMatch(/^\/password\?/);
      expect(result?.redirect).toContain("loginName=user%40example.com");
      expect(result?.redirect).toContain("requestId=req123");
      expect(result?.redirect).toContain("organization=org123");
    });

    test("should return error when user not found and no registration allowed", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowRegister: false,
        allowLocalAuthentication: true,
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "errors.userNotFound" });
    });

    test("should discover organization from domain suffix when user not found without org context", async () => {
      // Mock login settings for instance level (no org context)
      mockGetLoginSettings
        .mockResolvedValueOnce({
          allowRegister: true,
          allowLocalAuthentication: true,
          ignoreUnknownUsernames: false,
        })
        // Mock login settings for discovered org - must include all necessary flags
        .mockResolvedValueOnce({
          allowDomainDiscovery: true,
          allowRegister: true,
          allowLocalAuthentication: true,
          ignoreUnknownUsernames: false,
        });

      // Mock org discovery to return one org with matching domain
      mockGetOrgsByDomain.mockResolvedValue({
        result: [{ id: "discovered-org-123", name: "Example Org" }],
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        requestId: "req123",
        // No organization parameter - this is the key test scenario
      });

      expect(result).toBeDefined();
      expect(result?.redirect).toMatch(/^\/register\?/);
      expect(result?.redirect).toContain("organization=discovered-org-123");
      expect(result?.redirect).toContain("requestId=req123");
      expect(result?.redirect).toContain("email=user%40example.com");

      // Verify org discovery was called with correct domain
      expect(mockGetOrgsByDomain).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        domain: "example.com",
      });
    });

    test("should redirect to IDP with discovered org when user not found and only IDP allowed", async () => {
      // Mock login settings for instance level (no org context)
      mockGetLoginSettings
        .mockResolvedValueOnce({
          allowRegister: true,
          allowLocalAuthentication: false,
        })
        // Mock login settings for discovered org - must include all necessary flags
        .mockResolvedValueOnce({
          allowDomainDiscovery: true,
          allowRegister: true,
          allowLocalAuthentication: false,
        });

      // Mock org discovery to return one org with matching domain
      mockGetOrgsByDomain.mockResolvedValue({
        result: [{ id: "discovered-org-456", name: "Example Org" }],
      });

      mockGetActiveIdentityProviders.mockResolvedValue({
        identityProviders: [{ id: "idp123", type: "OIDC" }],
      });
      mockStartIdentityProviderFlow.mockResolvedValue({ url: "https://idp.example.com/auth?org=discovered-org-456" });

      const result = await sendLoginname({
        loginName: "user@company.com",
        requestId: "req123",
        // No organization parameter
      });

      expect(result).toEqual({ redirect: "https://idp.example.com/auth?org=discovered-org-456" });

      // Verify org discovery was called
      expect(mockGetOrgsByDomain).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        domain: "company.com",
      });

      // Verify IDP redirect was called with discovered org
      expect(mockGetActiveIdentityProviders).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        orgId: "discovered-org-456",
      });
    });

    test("should not discover org if domain discovery is disabled", async () => {
      mockGetLoginSettings
        .mockResolvedValueOnce({
          allowRegister: true,
          allowLocalAuthentication: true,
          ignoreUnknownUsernames: false,
        })
        // Mock login settings for org with domain discovery disabled
        .mockResolvedValueOnce({
          allowDomainDiscovery: false,
        });

      mockGetOrgsByDomain.mockResolvedValue({
        result: [{ id: "10987654321", name: "Example Org" }],
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        // No organization parameter
      });

      // Should return error since discovery is disabled and no org context
      expect(result).toEqual({ error: "errors.userNotFound" });
    });

    test("should not discover org if multiple orgs match the domain", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowRegister: true,
        allowLocalAuthentication: true,
        ignoreUnknownUsernames: false,
      });

      // Mock org discovery to return multiple orgs
      mockGetOrgsByDomain.mockResolvedValue({
        result: [
          { id: "12345678910", name: "Example Org 1" },
          { id: "10987654321", name: "Example Org 2" },
        ],
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        // No organization parameter
      });

      // Should return error since multiple orgs match
      expect(result).toEqual({ error: "errors.userNotFound" });
    });

    test("should use provided organization instead of discovering when org context exists", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowRegister: true,
        allowLocalAuthentication: true,
        ignoreUnknownUsernames: false,
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        organization: "123456",
        requestId: "req123",
      });

      expect(result).toBeDefined();
      expect(result?.redirect).toMatch(/^\/register\?/);
      expect(result?.redirect).toContain("organization=123456");

      // Verify org discovery was NOT called since org was provided
      expect(mockGetOrgsByDomain).not.toHaveBeenCalled();
    });

    test("should redirect to password when ignoreUnknownUsernames is true, allowRegister is true, but allowLocalAuthentication is false (User not found)", async () => {
      mockSearchUsers.mockResolvedValue({ result: [] });
      mockGetLoginSettings.mockResolvedValue({
        allowRegister: true,
        allowLocalAuthentication: false,
        ignoreUnknownUsernames: true,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({ identityProviders: [] });

      const result = await sendLoginname({
        loginName: "user@example.com",
        ignoreUnknownUsernames: true,
      });

      expect(result).not.toEqual({ error: "errors.userNotFound" });
      expect(result).toHaveProperty("redirect");
      expect((result as any).redirect).toMatch(/^\/password\?/);
    });

    test("should redirect to password when ignoreUnknownUsernames is true, user found but rejected due to disableLoginWithEmail", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" }, phone: { phone: "123456" } } },
        state: UserState.ACTIVE,
      };

      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockGetLoginSettings.mockResolvedValue({
        disableLoginWithEmail: true,
        ignoreUnknownUsernames: true,
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        ignoreUnknownUsernames: true,
      });

      expect(result).not.toEqual({ error: "errors.userNotFound" });
      expect(result).toHaveProperty("redirect");
      expect((result as any).redirect).toMatch(/^\/password\?/);
    });

    test("should allow login with email when disableLoginWithPhone is true and preferred login name differs from email", async () => {
      // Regression test for: https://github.com/zitadel/zitadel/issues/11518
      // When preferredLoginName (e.g. username@org-domain) differs from the user's email, logging
      // in with the email while disableLoginWithPhone=true must NOT be blocked.
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@orgdomain.com", // org-domain scoped login name
        details: { resourceOwner: "org123" },
        type: {
          case: "human",
          value: { email: { email: "user@test.com" }, phone: { phone: "+1234567890" } },
        },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@orgdomain.com",
            organizationId: "org123",
          },
        },
      };

      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockGetLoginSettings.mockResolvedValue({
        disableLoginWithPhone: true,
        allowLocalAuthentication: true,
      });
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD],
      });

      const result = await sendLoginname({
        loginName: "user@test.com", // logging in with email, not phone
      });

      // Must NOT return "User not found" — email-based login must be allowed when only phone is disabled
      expect(result).not.toEqual({ error: "errors.userNotFound" });
      expect(mockCreateSessionAndUpdateCookie).toHaveBeenCalled();
    });

    test("should block login with phone number when disableLoginWithPhone is true", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@orgdomain.com",
        details: { resourceOwner: "org123" },
        type: {
          case: "human",
          value: { email: { email: "user@example.com" }, phone: { phone: "+1234567890" } },
        },
        state: UserState.ACTIVE,
      };

      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockGetLoginSettings.mockResolvedValue({
        disableLoginWithPhone: true,
      });

      const result = await sendLoginname({
        loginName: "+1234567890", // logging in with phone — must be blocked
      });

      expect(result).toEqual({ error: "errors.userNotFound" });
      expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test("should redirect to password when ignoreUnknownUsernames is true and more than one user found", async () => {
      mockGetLoginSettings.mockResolvedValue({
        ignoreUnknownUsernames: true,
      });
      mockSearchUsers.mockResolvedValue({
        result: [
          { userId: "user1", preferredLoginName: "user1@example.com" },
          { userId: "user2", preferredLoginName: "user2@example.com" },
        ],
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        ignoreUnknownUsernames: true,
      });

      expect(result).not.toEqual({ error: "errors.moreThanOneUserFound" });
      expect(result).toHaveProperty("redirect");
      expect((result as any).redirect).toMatch(/^\/password\?/);
    });

    test("should return generic error when user not active and ignoreUnknownUsernames is true", async () => {
      mockSearchUsers.mockResolvedValue({
        result: [
          {
            userId: "user1",
            state: UserState.ACTIVE,
            preferredLoginName: "user1",
            type: { case: "human", value: { email: { isVerified: true } } },
          },
        ],
      });
      mockGetLoginSettings.mockResolvedValue({
        ignoreUnknownUsernames: true,
        allowLocalAuthentication: true,
      });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD],
      });
      // Mock createSessionAndUpdateCookie to fail with user not active error
      mockCreateSessionAndUpdateCookie.mockRejectedValue({
        rawMessage: "Errors.User.NotActive (SESSION-Gj4ko)",
      });

      const result = await sendLoginname({ loginName: "user1", ignoreUnknownUsernames: true });

      expect(result).toEqual({ redirect: "/password?loginName=user1" });
      // With ignoreUnknownUsernames: true, we skip session creation, so this mock is NOT called
      expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test("should NOT create session and return redirect when ignoreUnknownUsernames is true and user is valid", async () => {
      mockSearchUsers.mockResolvedValue({
        result: [
          {
            userId: "user1",
            state: UserState.ACTIVE,
            preferredLoginName: "user1",
            type: { case: "human", value: { email: { isVerified: true } } },
          },
        ],
      });
      mockGetLoginSettings.mockResolvedValue({
        ignoreUnknownUsernames: true,
        allowLocalAuthentication: true,
      });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD],
      });

      const result = await sendLoginname({ loginName: "user1", ignoreUnknownUsernames: true });

      expect(result).toEqual({ redirect: "/password?loginName=user1" });
      expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test("should redirect to password when ignoreUnknownUsernames is true and password not allowed", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "org123",
          },
        },
      };

      mockGetLoginSettings.mockResolvedValue({
        allowLocalAuthentication: false,
        ignoreUnknownUsernames: true,
      });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD],
      });
      mockListIDPLinks.mockResolvedValue({ result: [] });
      mockGetActiveIdentityProviders.mockResolvedValue({ identityProviders: [] });

      const result = await sendLoginname({
        loginName: "user@example.com",
        ignoreUnknownUsernames: true,
      });

      expect(result).not.toEqual({ error: "errors.localAuthenticationNotAllowed" });
      expect(result).toHaveProperty("redirect");
      expect((result as any).redirect).toMatch(/^\/password\?/);
    });

    test("should redirect to password when ignoreUnknownUsernames is true and passkeys not allowed", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "org123",
          },
        },
      };

      mockGetLoginSettings.mockResolvedValue({
        passkeysType: PasskeysType.NOT_ALLOWED,
        ignoreUnknownUsernames: true,
      });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSKEY],
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        ignoreUnknownUsernames: true,
      });

      expect(result).not.toEqual({ error: "errors.passkeysNotAllowed" });
      expect(result).toHaveProperty("redirect");
      expect((result as any).redirect).toMatch(/^\/password\?/);
    });
  });

  describe("Edge cases", () => {
    test("should handle initial user state", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.INITIAL,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "org123",
          },
        },
      };

      mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: true });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "errors.initialUserNotSupported" });
    });

    test("should handle organization parameter in all redirects", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "org123",
          },
        },
      };

      mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: true });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD],
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
        organization: "custom-org",
        requestId: "req123",
      });

      expect(result).toBeDefined();
      expect(result?.redirect).toContain("organization=custom-org");
      expect(result?.redirect).toContain("requestId=req123");
    });

    test("should redirect to password with INPUT loginName when ignoreUnknownUsernames is true, even if user preferredLoginName is different", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "org123",
          },
        },
      };

      mockGetLoginSettings.mockResolvedValue({
        allowLocalAuthentication: true,
        ignoreUnknownUsernames: true,
      });
      // Mock search result returns a user with resolved/different loginName
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD],
      });

      // INPUT login name is just "user"
      const result = await sendLoginname({
        loginName: "user",
        ignoreUnknownUsernames: true,
      });

      expect(result).toHaveProperty("redirect");
      // Expect redirect to contain input "user" not resolved "user@example.com"
      expect((result as any).redirect).toContain("loginName=user");
      expect((result as any).redirect).not.toContain("loginName=user%40example.com");
    });

    test("should redirect to passkey with INPUT loginName when ignoreUnknownUsernames is true", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "org123",
          },
        },
      };

      mockGetLoginSettings.mockResolvedValue({
        passkeysType: PasskeysType.ALLOWED,
        allowLocalAuthentication: true,
        ignoreUnknownUsernames: true,
      });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSKEY],
      });

      const result = await sendLoginname({
        loginName: "user",
        ignoreUnknownUsernames: true,
      });

      expect(result).toHaveProperty("redirect");
      expect((result as any).redirect).toMatch(/^\/passkey\?/);
      expect((result as any).redirect).toContain("loginName=user");
      expect((result as any).redirect).not.toContain("loginName=user%40example.com");
    });

    test("should redirect to passkey with INPUT loginName when ignoreUnknownUsernames is true (multi-method)", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "org123",
          },
        },
      };

      mockGetLoginSettings.mockResolvedValue({
        passkeysType: PasskeysType.ALLOWED,
        allowLocalAuthentication: true,
        ignoreUnknownUsernames: true,
      });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSKEY, AuthenticationMethodType.PASSWORD],
      });

      const result = await sendLoginname({
        loginName: "user",
        ignoreUnknownUsernames: true,
      });

      expect(result).toHaveProperty("redirect");
      expect((result as any).redirect).toMatch(/^\/passkey\?/);
      expect((result as any).redirect).toContain("loginName=user");
      expect((result as any).redirect).not.toContain("loginName=user%40example.com");
    });
    test("should use CONTEXT settings to HIDE username even if USER settings would show it", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "user-org" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "user-org",
          },
        },
      };

      // Mock implementation to return different settings based on organization
      mockGetLoginSettings.mockImplementation(async (args: any) => {
        if (args.organization === "context-org") {
          return { allowLocalAuthentication: true, ignoreUnknownUsernames: true };
        }
        if (args.organization === "user-org") {
          return { allowLocalAuthentication: true, ignoreUnknownUsernames: false };
        }
        return {};
      });

      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD],
      });

      const result = await sendLoginname({
        loginName: "input-name",
        organization: "context-org", // Context has ignore=true
        ignoreUnknownUsernames: true,
      });

      expect(result).toHaveProperty("redirect");
      // Should result in input-name because context says HIDE
      expect((result as any).redirect).toContain("loginName=input-name");
      expect((result as any).redirect).not.toContain("loginName=user%40example.com");
    });

    test("should use CONTEXT settings to SHOW username even if USER settings would hide it", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "user-org" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      const mockSession = {
        factors: {
          user: {
            id: "user123",
            loginName: "user@example.com",
            organizationId: "user-org",
          },
        },
      };

      mockGetLoginSettings.mockImplementation(async (args: any) => {
        if (args.organization === "context-org") {
          return { allowLocalAuthentication: true, ignoreUnknownUsernames: false };
        }
        if (args.organization === "user-org") {
          return { allowLocalAuthentication: true, ignoreUnknownUsernames: true };
        }
        return {};
      });

      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ session: mockSession, sessionCookie: {} });
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD],
      });

      const result = await sendLoginname({
        loginName: "input-name",
        organization: "context-org", // Context has ignore=false
        ignoreUnknownUsernames: false,
      });

      expect(result).toHaveProperty("redirect");
      // Should result in resolved name because context says SHOW
      expect((result as any).redirect).toContain("loginName=user%40example.com");
    });
  });
});

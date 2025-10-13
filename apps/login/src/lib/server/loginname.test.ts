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
  getServiceUrlFromHeaders: vi.fn(),
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
  getOriginalHost: vi.fn(),
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
  let mockGetOriginalHost: any;
  let mockStartIdentityProviderFlow: any;
  let mockGetActiveIdentityProviders: any;
  let mockGetIDPByID: any;
  let mockIdpTypeToSlug: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    // Import mocked modules
    const { headers } = await import("next/headers");
    const { create } = await import("@zitadel/client");
    const { getServiceUrlFromHeaders } = await import("../service-url");
    const {
      getLoginSettings,
      searchUsers,
      listAuthenticationMethodTypes,
      listIDPLinks,
      startIdentityProviderFlow,
      getActiveIdentityProviders,
    } = await import("../zitadel");
    const { createSessionAndUpdateCookie } = await import("./cookie");
    const { getOriginalHost } = await import("./host");
    const { idpTypeToSlug } = await import("../idp");

    // Setup mocks
    mockHeaders = vi.mocked(headers);
    mockCreate = vi.mocked(create);
    mockGetServiceUrlFromHeaders = vi.mocked(getServiceUrlFromHeaders);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockSearchUsers = vi.mocked(searchUsers);
    mockCreateSessionAndUpdateCookie = vi.mocked(createSessionAndUpdateCookie);
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockListIDPLinks = vi.mocked(listIDPLinks);
    mockGetOriginalHost = vi.mocked(getOriginalHost);
    mockStartIdentityProviderFlow = vi.mocked(startIdentityProviderFlow);
    mockGetActiveIdentityProviders = vi.mocked(getActiveIdentityProviders);
    mockGetIDPByID = vi.mocked(getIDPByID);
    mockIdpTypeToSlug = vi.mocked(idpTypeToSlug);

    // Default mock implementations
    mockHeaders.mockResolvedValue({} as any);
    mockGetServiceUrlFromHeaders.mockReturnValue({ serviceUrl: "https://api.example.com" });
    mockGetOriginalHost.mockResolvedValue("example.com");
    mockIdpTypeToSlug.mockReturnValue("google");
    mockGetIDPByID.mockResolvedValue({
      id: "idp123",
      name: "Google",
      type: "GOOGLE",
    });
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
      mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: true });
      mockSearchUsers.mockResolvedValue({ error: "Search failed" });

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "Search failed" });
    });

    test("should return error when search result has no result field", async () => {
      mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: true });
      mockSearchUsers.mockResolvedValue({});

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "errors.couldNotSearchUsers" });
    });

    test("should return error when more than one user found", async () => {
      mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: true });
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
      mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: true });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue(mockSession);
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
        mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });
        mockListIDPLinks.mockResolvedValue({
          result: [{ idpId: "idp123" }],
        });
        mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({ redirect: "https://idp.example.com/auth" });
        expect(mockListIDPLinks).toHaveBeenCalledWith({
          serviceUrl: "https://api.example.com",
          userId: "user123",
        });
      });

      test("should return error when password not allowed and no IDP links available", async () => {
        mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });
        mockListIDPLinks.mockResolvedValue({ result: [] });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({
          error: "errors.usernamePasswordNotAllowed",
        });
      });

      test("should redirect to passkey when user has only passkey method and it's allowed", async () => {
        mockGetLoginSettings.mockResolvedValue({ passkeysType: PasskeysType.ALLOWED });
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
        mockGetLoginSettings.mockResolvedValue({ passkeysType: PasskeysType.NOT_ALLOWED });
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
        mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({ redirect: "https://idp.example.com/auth" });
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

      test("should not show password alternative when password is not allowed", async () => {
        mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.PASSKEY],
        });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toBeDefined();
        expect(result?.redirect).toMatch(/^\/passkey\?/);
        expect(result?.redirect).toContain("altPassword=false"); // password is not allowed
      });

      test("should redirect to IDP when no passkey but IDP available", async () => {
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.IDP],
        });
        mockListIDPLinks.mockResolvedValue({
          result: [{ idpId: "idp123" }],
        });
        mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

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
        mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: false });
        mockListAuthenticationMethodTypes.mockResolvedValue({
          authMethodTypes: [AuthenticationMethodType.PASSWORD],
        });
        mockListIDPLinks.mockResolvedValue({ result: [] });

        const result = await sendLoginname({
          loginName: "user@example.com",
        });

        expect(result).toEqual({
          error: "errors.usernamePasswordNotAllowed",
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
        allowUsernamePassword: false,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({
        identityProviders: [{ id: "idp123", type: "OIDC" }],
      });
      mockStartIdentityProviderFlow.mockResolvedValue("https://idp.example.com/auth");

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ redirect: "https://idp.example.com/auth" });
    });

    test("should redirect to register when both register and password allowed", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowRegister: true,
        allowUsernamePassword: true,
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
        allowUsernamePassword: true,
      });

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "errors.userNotFound" });
    });
  });

  describe("Edge cases", () => {
    test("should handle session creation failure", async () => {
      const mockUser = {
        userId: "user123",
        preferredLoginName: "user@example.com",
        details: { resourceOwner: "org123" },
        type: { case: "human", value: { email: { email: "user@example.com" } } },
        state: UserState.ACTIVE,
      };

      mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: true });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue({ factors: {} }); // No user in session

      const result = await sendLoginname({
        loginName: "user@example.com",
      });

      expect(result).toEqual({ error: "errors.couldNotCreateSession" });
    });

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

      mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: true });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue(mockSession);

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

      mockGetLoginSettings.mockResolvedValue({ allowUsernamePassword: true });
      mockSearchUsers.mockResolvedValue({ result: [mockUser] });
      mockCreate.mockReturnValue({});
      mockCreateSessionAndUpdateCookie.mockResolvedValue(mockSession);
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
  });
});

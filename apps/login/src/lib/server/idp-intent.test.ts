import { describe, expect, test, vi, beforeEach, afterEach } from "vitest";
import { processIDPCallback } from "./idp-intent";
import { AutoLinkingOption } from "@zitadel/proto/zitadel/idp/v2/idp_pb";
import crypto from "crypto";

// Mock all the dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: vi.fn((schema: any, data: any) => data),
  ConnectError: class extends Error {
    code: any;
    constructor(message: string, code: any) {
      super(message);
      this.code = code;
    }
  },
  Code: {
    AlreadyExists: 6,
  },
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  retrieveIDPIntent: vi.fn(),
  getIDPByID: vi.fn(),
  updateHuman: vi.fn(),
  addIDPLink: vi.fn(),
  listUsers: vi.fn(),
  addHuman: vi.fn(),
  getLoginSettings: vi.fn(),
  getOrgsByDomain: vi.fn(),
  getActiveIdentityProviders: vi.fn(),
  getUserByID: vi.fn(),
  getDefaultOrg: vi.fn(),
  getSession: vi.fn(),
}));

vi.mock("./idp", () => ({
  createNewSessionFromIdpIntent: vi.fn(),
}));

vi.mock("@/lib/cookies", () => ({
  getSessionCookieById: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("../fingerprint", () => ({
  getFingerprintIdCookie: vi.fn(),
}));

describe("processIDPCallback", () => {
  // Mock modules
  let mockHeaders: any;
  let mockGetServiceUrlFromHeaders: any;
  let mockRetrieveIDPIntent: any;
  let mockGetIDPByID: any;
  let mockUpdateHuman: any;
  let mockAddIDPLink: any;
  let mockListUsers: any;
  let mockAddHuman: any;
  let mockGetLoginSettings: any;
  let mockGetOrgsByDomain: any;
  let mockGetActiveIdentityProviders: any;
  let mockGetUserByID: any;
  let mockGetDefaultOrg: any;
  let mockGetSession: any;
  let mockCreateNewSessionFromIdpIntent: any;
  let mockGetSessionCookieById: any;
  let mockGetFingerprintIdCookie: any;

  const defaultParams = {
    provider: "google",
    id: "intent123",
    token: "token123",
    requestId: "req123",
    organization: "org123",
  };

  const defaultIntent = {
    idpInformation: {
      idpId: "idp123",
      userId: "user123",
      userName: "testuser",
    },
    userId: "user123",
    addHumanUser: {
      username: "testuser",
      profile: {
        givenName: "Test",
        familyName: "User",
        displayName: "Test User",
      },
      email: {
        email: "test@example.com",
      },
    },
    updateHumanUser: {
      username: "testuser",
      profile: {
        givenName: "Test",
        familyName: "User 1",
        displayName: "Test User 1",
      },
      email: {
        email: "test@example.com",
      },
    },
  };

  const defaultIdp = {
    id: "idp123",
    config: {
      options: {
        isAutoUpdate: false,
        isLinkingAllowed: false,
        isCreationAllowed: false,
        isAutoCreation: false,
        autoLinking: undefined,
      },
    },
  };

  beforeEach(async () => {
    vi.clearAllMocks();

    // Import mocked modules
    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const {
      retrieveIDPIntent,
      getIDPByID,
      updateHuman,
      addIDPLink,
      listUsers,
      addHuman,
      getLoginSettings,
      getOrgsByDomain,
      getActiveIdentityProviders,
      getUserByID,
      getDefaultOrg,
      getSession,
    } = await import("../zitadel");
    const { createNewSessionFromIdpIntent } = await import("./idp");
    const { getSessionCookieById } = await import("@/lib/cookies");
    const { getFingerprintIdCookie } = await import("../fingerprint");

    // Setup mocks
    mockHeaders = vi.mocked(headers);
    mockGetServiceUrlFromHeaders = vi.mocked(getServiceConfig);
    mockRetrieveIDPIntent = vi.mocked(retrieveIDPIntent);
    mockGetIDPByID = vi.mocked(getIDPByID);
    mockUpdateHuman = vi.mocked(updateHuman);
    mockAddIDPLink = vi.mocked(addIDPLink);
    mockListUsers = vi.mocked(listUsers);
    mockAddHuman = vi.mocked(addHuman);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockGetOrgsByDomain = vi.mocked(getOrgsByDomain);
    mockGetActiveIdentityProviders = vi.mocked(getActiveIdentityProviders);
    mockGetUserByID = vi.mocked(getUserByID);
    mockGetDefaultOrg = vi.mocked(getDefaultOrg);
    mockGetSession = vi.mocked(getSession);

    mockCreateNewSessionFromIdpIntent = vi.mocked(createNewSessionFromIdpIntent);
    mockGetSessionCookieById = vi.mocked(getSessionCookieById);
    mockGetFingerprintIdCookie = vi.mocked(getFingerprintIdCookie);

    // Default mock implementations
    mockHeaders.mockResolvedValue({} as any);
    mockGetServiceUrlFromHeaders.mockReturnValue({
      serviceConfig: { baseUrl: "https://api.example.com" },
    });
    mockRetrieveIDPIntent.mockResolvedValue(defaultIntent);
    mockGetIDPByID.mockResolvedValue(defaultIdp);
    mockCreateNewSessionFromIdpIntent.mockResolvedValue({
      redirect: "https://app.example.com/success",
    });

    // Default mocks for validation functions
    mockGetUserByID.mockResolvedValue({
      userId: "user123",
      details: {
        resourceOwner: "org123",
      },
    });
    mockGetLoginSettings.mockResolvedValue({
      allowExternalIdp: true,
    });
    mockGetActiveIdentityProviders.mockResolvedValue({
      identityProviders: [{ id: "idp123", name: "Test IDP" }],
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("Parameter validation", () => {
    test("should return error redirect when provider is missing", async () => {
      const result = await processIDPCallback({
        provider: "",
        id: "intent123",
        token: "token123",
      });

      expect(result.redirect).toContain("/idp//failure");
      expect(mockRetrieveIDPIntent).not.toHaveBeenCalled();
    });

    test("should return error redirect when id is missing", async () => {
      const result = await processIDPCallback({
        provider: "google",
        id: "",
        token: "token123",
      });

      expect(result.redirect).toContain("/idp/google/failure");
      expect(mockRetrieveIDPIntent).not.toHaveBeenCalled();
    });

    test("should return error redirect when token is missing", async () => {
      const result = await processIDPCallback({
        provider: "google",
        id: "intent123",
        token: "",
      });

      expect(result.redirect).toContain("/idp/google/failure");
      expect(mockRetrieveIDPIntent).not.toHaveBeenCalled();
    });

    test("should preserve requestId and organization in error redirect", async () => {
      const result = await processIDPCallback({
        provider: "google",
        id: "",
        token: "token123",
        requestId: "req123",
        organization: "org123",
      });

      expect(result.redirect).toContain("requestId=req123");
      expect(result.redirect).toContain("organization=org123");
    });
  });

  describe("Intent retrieval errors", () => {
    test("should return error redirect when IDP information is missing", async () => {
      mockRetrieveIDPIntent.mockResolvedValue({
        idpInformation: undefined,
        userId: "user123",
      });

      const result = await processIDPCallback(defaultParams);

      expect(result.redirect).toContain("/idp/google/failure");
      expect(result.redirect).toContain("error=missing_idp_info");
    });

    test("should return error when IDP is not found", async () => {
      mockGetIDPByID.mockResolvedValue(null);

      const result = await processIDPCallback(defaultParams);

      expect(result.error).toBe("errors.idpNotFound");
    });

    test("should handle retrieval errors gracefully", async () => {
      mockRetrieveIDPIntent.mockRejectedValue(new Error("Network error"));

      const result = await processIDPCallback(defaultParams);

      expect(result.redirect).toContain("/idp/google/failure");
      expect(result.redirect).toContain("error=Network+error");
    });
  });

  describe("CASE 1: User exists and should sign in", () => {
    test("should create session for existing user without auto-update", async () => {
      const result = await processIDPCallback(defaultParams);

      expect(mockRetrieveIDPIntent).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        id: "intent123",
        token: "token123",
      });
      expect(mockCreateNewSessionFromIdpIntent).toHaveBeenCalledWith({
        userId: "user123",
        idpIntent: {
          idpIntentId: "intent123",
          idpIntentToken: "token123",
        },
        requestId: "req123",
        organization: "org123",
      });
      expect(result.redirect).toBe("https://app.example.com/success");
    });

    test("should auto-update user profile when enabled", async () => {
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            isAutoUpdate: true,
          },
        },
      });

      await processIDPCallback(defaultParams);

      expect(mockUpdateHuman).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        request: expect.objectContaining({
          userId: "user123",
          profile: defaultIntent.updateHumanUser.profile,
          email: defaultIntent.updateHumanUser.email,
        }),
      });
    });

    test("should continue session creation even if auto-update fails", async () => {
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            isAutoUpdate: true,
          },
        },
      });
      mockUpdateHuman.mockRejectedValue(new Error("Update failed"));

      const result = await processIDPCallback(defaultParams);

      expect(mockCreateNewSessionFromIdpIntent).toHaveBeenCalled();
      expect(result.redirect).toBe("https://app.example.com/success");
    });

    test("should return error when session creation fails", async () => {
      mockCreateNewSessionFromIdpIntent.mockResolvedValue({
        error: "Session creation error",
      });

      const result = await processIDPCallback(defaultParams);

      expect(result.error).toBe("Session creation error");
    });

    test("should return error when session creation returns neither redirect nor error", async () => {
      mockCreateNewSessionFromIdpIntent.mockResolvedValue({});

      const result = await processIDPCallback(defaultParams);

      expect(result.error).toBe("errors.sessionCreationFailed");
    });
  });

  describe("CASE 2: Link IDP to existing user", () => {
    const sessionId = "session123";
    const fingerprintValue = "fp123";
    const linkFingerprint = crypto
      .createHash("sha256")
      .update(sessionId + fingerprintValue)
      .digest("hex");

    const linkParams = {
      ...defaultParams,
      sessionId,
      linkFingerprint,
    };

    beforeEach(() => {
      mockGetFingerprintIdCookie.mockResolvedValue({ name: "fingerprintId", value: fingerprintValue });
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
      });
      mockGetSessionCookieById.mockResolvedValue({ id: sessionId, token: "token123" });
      mockGetSession.mockResolvedValue({
        session: {
          factors: {
            user: {
              id: "user123",
            },
          },
        },
      });
    });

    test("should link IDP and create session when linking is allowed", async () => {
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            isLinkingAllowed: true,
          },
        },
      });

      const result = await processIDPCallback(linkParams);

      expect(mockGetSessionCookieById).toHaveBeenCalledWith({ sessionId: "session123" });

      expect(mockAddIDPLink).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        idp: {
          id: "idp123",
          userId: "user123",
          userName: "testuser",
        },
        userId: "user123",
      });
      expect(mockCreateNewSessionFromIdpIntent).toHaveBeenCalled();
      expect(result.redirect).toBe("https://app.example.com/success");
    });

    test("should return error redirect when linking is not allowed", async () => {
      const result = await processIDPCallback(linkParams);

      expect(result.redirect).toContain("/idp/google/linking-failed");
      expect(result.redirect).toContain("error=linking_not_allowed");
      expect(mockAddIDPLink).not.toHaveBeenCalled();
    });

    test("should return error redirect when linking fails", async () => {
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            isLinkingAllowed: true,
          },
        },
      });
      mockAddIDPLink.mockRejectedValue(new Error("Linking failed"));

      const result = await processIDPCallback(linkParams);

      expect(result.redirect).toContain("/idp/google/linking-failed");
    });

    test("should return error when session creation fails after linking", async () => {
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            isLinkingAllowed: true,
          },
        },
      });
      mockCreateNewSessionFromIdpIntent.mockResolvedValue({
        error: "Session error",
      });

      const result = await processIDPCallback(linkParams);

      expect(result.error).toBe("Session error");
    });
  });

  describe("CASE 3: Auto-linking by email", () => {
    beforeEach(() => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined, // No existing userId
      });
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            autoLinking: AutoLinkingOption.EMAIL,
            isLinkingAllowed: true,
          },
        },
      });
    });

    test("should auto-link user by email and create session", async () => {
      const foundUser = {
        userId: "found123",
        details: {
          resourceOwner: "org123",
        },
      };
      mockListUsers.mockResolvedValue({
        result: [foundUser],
      });

      const result = await processIDPCallback(defaultParams);

      expect(mockListUsers).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        email: "test@example.com",
        organizationId: "org123",
      });
      expect(mockAddIDPLink).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        idp: {
          id: "idp123",
          userId: "user123",
          userName: "testuser",
        },
        userId: "found123",
      });
      expect(mockCreateNewSessionFromIdpIntent).toHaveBeenCalledWith({
        userId: "found123",
        idpIntent: {
          idpIntentId: "intent123",
          idpIntentToken: "token123",
        },
        requestId: "req123",
        organization: "org123",
      });
      expect(result.redirect).toBe("https://app.example.com/success");
    });

    test("should continue to next case when no user found by email", async () => {
      mockListUsers.mockResolvedValue({
        result: [],
      });

      const result = await processIDPCallback(defaultParams);

      expect(mockAddIDPLink).not.toHaveBeenCalled();
      expect(result.redirect).toContain("/idp/google/account-not-found");
    });

    test("should return error redirect when auto-linking fails", async () => {
      mockListUsers.mockResolvedValue({
        result: [{ userId: "found123" }],
      });
      mockAddIDPLink.mockRejectedValue(new Error("Linking failed"));

      const result = await processIDPCallback(defaultParams);

      expect(result.redirect).toContain("/idp/google/linking-failed");
    });
  });

  describe("CASE 3: Auto-linking by username", () => {
    beforeEach(() => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
      });
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            autoLinking: AutoLinkingOption.USERNAME,
            isLinkingAllowed: true,
          },
        },
      });
    });

    test("should auto-link user by username", async () => {
      mockListUsers.mockResolvedValue({
        result: [
          {
            userId: "found123",
            details: {
              resourceOwner: "org123",
            },
          },
        ],
      });

      const result = await processIDPCallback(defaultParams);

      expect(mockListUsers).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userName: "testuser",
        organizationId: "org123",
      });
      expect(mockAddIDPLink).toHaveBeenCalled();
      expect(result.redirect).toBe("https://app.example.com/success");
    });
  });

  describe("CASE 4: Auto-creation of user", () => {
    beforeEach(() => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
      });
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            isAutoCreation: true,
          },
        },
      });
    });

    test("should auto-create user and create session", async () => {
      mockAddHuman.mockResolvedValue({
        userId: "newuser123",
      });

      const result = await processIDPCallback(defaultParams);

      expect(mockAddHuman).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        request: expect.objectContaining({
          username: "testuser",
          profile: defaultIntent.addHumanUser.profile,
          email: defaultIntent.addHumanUser.email,
          organization: expect.objectContaining({
            org: { case: "orgId", value: "org123" },
          }),
        }),
      });
      expect(mockCreateNewSessionFromIdpIntent).toHaveBeenCalledWith({
        userId: "newuser123",
        idpIntent: {
          idpIntentId: "intent123",
          idpIntentToken: "token123",
        },
        requestId: "req123",
        organization: "org123",
      });
      expect(result.redirect).toBe("https://app.example.com/success");
    });

    test("should resolve organization from username domain", async () => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
        addHumanUser: {
          ...defaultIntent.addHumanUser,
          username: "user@example.com",
        },
      });
      mockGetOrgsByDomain.mockResolvedValue({
        result: [{ id: "org-from-domain" }],
      });
      mockGetLoginSettings.mockResolvedValue({
        allowDomainDiscovery: true,
      });
      mockAddHuman.mockResolvedValue({ userId: "newuser123" });

      await processIDPCallback({
        ...defaultParams,
        organization: undefined,
      });

      expect(mockGetOrgsByDomain).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        domain: "example.com",
      });
      expect(mockAddHuman).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        request: expect.objectContaining({
          organization: expect.objectContaining({
            org: { case: "orgId", value: "org-from-domain" },
          }),
        }),
      });
    });

    test("should fallback to default organization when not resolved", async () => {
      mockGetDefaultOrg.mockResolvedValue({ id: "default-org" });
      mockAddHuman.mockResolvedValue({ userId: "newuser123" });

      await processIDPCallback({
        ...defaultParams,
        organization: undefined,
      });

      expect(mockGetDefaultOrg).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
      });
      expect(mockAddHuman).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        request: expect.objectContaining({
          organization: expect.objectContaining({
            org: { case: "orgId", value: "default-org" },
          }),
        }),
      });
    });

    test("should return error when no organization context can be determined", async () => {
      mockGetDefaultOrg.mockResolvedValue(null);

      const result = await processIDPCallback({
        ...defaultParams,
        organization: undefined,
      });

      expect(mockGetDefaultOrg).toHaveBeenCalled();
      expect(mockAddHuman).not.toHaveBeenCalled();
      expect(result.redirect).toContain("/idp/google/failure");
      expect(result.redirect).toContain("error=no_organization_context");
    });

    test("should return error redirect when user creation fails", async () => {
      mockAddHuman.mockRejectedValue(new Error("Creation failed"));

      const result = await processIDPCallback(defaultParams);

      expect(result.redirect).toContain("/idp/google/failure");
      expect(result.redirect).toContain("error=user_creation_failed");
    });

    test("should return error when session creation fails after user creation", async () => {
      mockAddHuman.mockResolvedValue({ userId: "newuser123" });
      mockCreateNewSessionFromIdpIntent.mockResolvedValue({
        error: "Session error",
      });

      const result = await processIDPCallback(defaultParams);

      expect(result.error).toBe("Session error");
    });
  });

  describe("CASE 5: Manual user creation allowed", () => {
    beforeEach(() => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
      });
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            isCreationAllowed: true,
          },
        },
      });
    });

    test("should redirect to complete registration page with user data", async () => {
      const result = await processIDPCallback(defaultParams);

      expect(result.redirect).toContain("/idp/google/complete-registration");
      expect(result.redirect).toContain("id=intent123");
      expect(result.redirect).toContain("token=token123");
      expect(result.redirect).toContain("requestId=req123");
      expect(result.redirect).toContain("organization=org123");
      expect(result.redirect).toContain("idpId=idp123");
      expect(result.redirect).toContain("idpUserId=user123");
      expect(result.redirect).toContain("idpUserName=testuser");
      expect(result.redirect).toContain("givenName=Test");
      expect(result.redirect).toContain("familyName=User");
      expect(result.redirect).toContain("email=test%40example.com");
    });

    test("should fallback to default organization for registration", async () => {
      mockGetDefaultOrg.mockResolvedValue({ id: "default-org" });

      const result = await processIDPCallback({
        ...defaultParams,
        organization: undefined,
      });

      expect(mockGetDefaultOrg).toHaveBeenCalled();
      expect(result.redirect).toContain("/idp/google/complete-registration");
      expect(result.redirect).toContain("organization=default-org");
    });

    test("should redirect to registration failed when no organization context available", async () => {
      mockGetDefaultOrg.mockResolvedValue(null);

      const result = await processIDPCallback({
        ...defaultParams,
        organization: undefined,
      });

      expect(mockGetDefaultOrg).toHaveBeenCalled();
      expect(result.redirect).toContain("/idp/google/registration-failed");
      expect(result.redirect).toContain("id=intent123");
    });

    test("should resolve organization from domain for registration", async () => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
        addHumanUser: {
          ...defaultIntent.addHumanUser,
          username: "user@example.com",
        },
      });
      mockGetOrgsByDomain.mockResolvedValue({
        result: [{ id: "org-from-domain" }],
      });
      mockGetLoginSettings.mockResolvedValue({
        allowDomainDiscovery: true,
      });

      const result = await processIDPCallback({
        ...defaultParams,
        organization: undefined,
      });

      expect(result.redirect).toContain("organization=org-from-domain");
    });
  });

  describe("CASE 6: No user found and creation not allowed", () => {
    beforeEach(() => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
      });
    });

    test("should redirect to account not found page", async () => {
      const result = await processIDPCallback(defaultParams);

      expect(result.redirect).toContain("/idp/google/account-not-found");
      expect(result.redirect).toContain("id=intent123");
      expect(result.redirect).toContain("requestId=req123");
      expect(result.redirect).toContain("organization=org123");
    });
  });

  describe("Priority of cases", () => {
    test("should prioritize existing user sign-in over auto-linking", async () => {
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            autoLinking: AutoLinkingOption.EMAIL,
          },
        },
      });

      await processIDPCallback(defaultParams);

      // Should not search for users when userId already exists
      expect(mockListUsers).not.toHaveBeenCalled();
      expect(mockCreateNewSessionFromIdpIntent).toHaveBeenCalledWith(
        expect.objectContaining({
          userId: "user123",
        }),
      );
    });

    test("should prioritize auto-linking over auto-creation", async () => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
      });
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            autoLinking: AutoLinkingOption.EMAIL,
            isLinkingAllowed: true,
            isAutoCreation: true,
          },
        },
      });
      mockListUsers.mockResolvedValue({
        result: [
          {
            userId: "found123",
            details: {
              resourceOwner: "org123",
            },
          },
        ],
      });

      await processIDPCallback(defaultParams);

      // Should link, not create
      expect(mockAddIDPLink).toHaveBeenCalled();
      expect(mockAddHuman).not.toHaveBeenCalled();
    });

    test("should prioritize auto-creation over manual creation", async () => {
      mockRetrieveIDPIntent.mockResolvedValue({
        ...defaultIntent,
        userId: undefined,
      });
      mockGetIDPByID.mockResolvedValue({
        ...defaultIdp,
        config: {
          options: {
            ...defaultIdp.config.options,
            isAutoCreation: true,
            isCreationAllowed: true,
          },
        },
      });
      mockAddHuman.mockResolvedValue({ userId: "newuser123" });

      const result = await processIDPCallback(defaultParams);

      // Should auto-create, not redirect to manual form
      expect(mockAddHuman).toHaveBeenCalled();
      expect(result.redirect).toBe("https://app.example.com/success");
      expect(result.redirect).not.toContain("complete-registration");
    });
  });

  describe("postErrorRedirectUrl handling", () => {
    test("should preserve postErrorRedirectUrl in all redirects", async () => {
      const paramsWithError = {
        ...defaultParams,
        postErrorRedirectUrl: "https://app.example.com/error",
      };

      mockRetrieveIDPIntent.mockRejectedValue(new Error("Test error"));

      const result = await processIDPCallback(paramsWithError);

      expect(result.redirect).toContain("postErrorRedirectUrl=https%3A%2F%2Fapp.example.com%2Ferror");
    });
  });
});

describe("validateIDPLinkingPermissions", () => {
  let mockGetLoginSettings: any;
  let mockGetActiveIdentityProviders: any;
  let validateIDPLinkingPermissions: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    // Import mocked modules
    const { getLoginSettings, getActiveIdentityProviders } = await import("../zitadel");
    const { validateIDPLinkingPermissions: validate } = await import("./idp-intent");

    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockGetActiveIdentityProviders = vi.mocked(getActiveIdentityProviders);
    validateIDPLinkingPermissions = validate;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("Organization login settings validation", () => {
    test("should return false when allowExternalIdp is false", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowExternalIdp: false,
      });

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(false);
      expect(mockGetLoginSettings).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        organization: "org123",
      });
      expect(mockGetActiveIdentityProviders).not.toHaveBeenCalled();
    });

    test("should return false when login settings are undefined", async () => {
      mockGetLoginSettings.mockResolvedValue(undefined);

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(false);
    });

    test("should return false when allowExternalIdp is missing", async () => {
      mockGetLoginSettings.mockResolvedValue({});

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(false);
    });
  });

  describe("Active IDP validation", () => {
    test("should return false when IDP is not in active list", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowExternalIdp: true,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({
        identityProviders: [
          { id: "idp456", name: "Other IDP" },
          { id: "idp789", name: "Another IDP" },
        ],
      });

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(false);
      expect(mockGetActiveIdentityProviders).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        orgId: "org123",
        linking_allowed: true,
      });
    });

    test("should return false when identityProviders list is empty", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowExternalIdp: true,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({
        identityProviders: [],
      });

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(false);
    });

    test("should return false when identityProviders is undefined", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowExternalIdp: true,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({});

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(false);
    });
  });

  describe("Successful validation", () => {
    test("should return true when all validations pass", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowExternalIdp: true,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({
        identityProviders: [
          { id: "idp123", name: "Target IDP" },
          { id: "idp456", name: "Other IDP" },
        ],
      });

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(true);
      expect(mockGetLoginSettings).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        organization: "org123",
      });
      expect(mockGetActiveIdentityProviders).toHaveBeenCalledWith({
        serviceConfig: { baseUrl: "https://api.example.com" },
        orgId: "org123",
        linking_allowed: true,
      });
    });

    test("should find IDP in middle of list", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowExternalIdp: true,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({
        identityProviders: [
          { id: "idp111", name: "IDP 1" },
          { id: "idp222", name: "IDP 2" },
          { id: "idp123", name: "Target IDP" },
          { id: "idp333", name: "IDP 3" },
        ],
      });

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(true);
    });
  });

  describe("Edge cases", () => {
    test("should handle getLoginSettings throwing an error", async () => {
      mockGetLoginSettings.mockRejectedValue(new Error("Network error"));

      await expect(
        validateIDPLinkingPermissions({
          serviceConfig: { baseUrl: "https://api.example.com" },
          userOrganizationId: "org123",
          idpId: "idp123",
        }),
      ).rejects.toThrow("Network error");
    });

    test("should handle getActiveIdentityProviders throwing an error", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowExternalIdp: true,
      });
      mockGetActiveIdentityProviders.mockRejectedValue(new Error("Network error"));

      await expect(
        validateIDPLinkingPermissions({
          serviceConfig: { baseUrl: "https://api.example.com" },
          userOrganizationId: "org123",
          idpId: "idp123",
        }),
      ).rejects.toThrow("Network error");
    });

    test("should perform case-sensitive IDP ID comparison", async () => {
      mockGetLoginSettings.mockResolvedValue({
        allowExternalIdp: true,
      });
      mockGetActiveIdentityProviders.mockResolvedValue({
        identityProviders: [{ id: "IDP123", name: "Target IDP" }],
      });

      const result = await validateIDPLinkingPermissions({
        serviceConfig: { baseUrl: "https://api.example.com" },
        userOrganizationId: "org123",
        idpId: "idp123",
      });

      expect(result).toBe(false);
    });
  });
});

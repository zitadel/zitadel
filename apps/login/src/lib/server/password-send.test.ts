import { create } from "@zitadel/client";
import { ChecksSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { UserState } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { sendPassword } from "./password";

// Mock dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: vi.fn((schema, data) => data),
  ConnectError: class extends Error {
    code: number;
    constructor(msg: string, code: number) {
      super(msg);
      this.code = code;
    }
  },
  timestampDate: (ts: any) => new Date(ts.seconds * 1000),
}));

vi.mock("@zitadel/client/v2", () => ({
  createUserServiceClient: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("../zitadel", () => ({
  getLoginSettings: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  getSession: vi.fn(),
  setPassword: vi.fn(),
  createServerTransport: vi.fn(),
  listUsers: vi.fn(),
  getLockoutSettings: vi.fn(),
  getPasswordExpirySettings: vi.fn(),
  getUserByID: vi.fn(),
  searchUsers: vi.fn(),
}));

vi.mock("../cookies", () => ({
  getSessionCookieByLoginName: vi.fn(),
  getSessionCookieById: vi.fn(),
}));

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("@/lib/server/cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("../verify-helper", () => ({
  checkPasswordChangeRequired: vi.fn(),
  checkEmailVerification: vi.fn(),
  checkMFAFactors: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("@/lib/metrics", () => ({
  recordAuthAttempt: vi.fn(),
  recordAuthFailure: vi.fn(),
  recordAuthSuccess: vi.fn(),
}));

describe("sendPassword", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetSessionCookieByLoginName: any;
  let mockSearchUsers: any;
  let mockGetLoginSettings: any;
  let mockCreateSessionAndUpdateCookie: any;
  let mockCompleteFlowOrGetUrl: any;
  let mockListAuthenticationMethodTypes: any;
  let mockCheckMFAFactors: any;
  let mockSetSessionAndUpdateCookie: any;
  let mockGetUserByID: any;
  let mockRecordAuthFailure: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getSessionCookieByLoginName } = await import("../cookies");
    const { getLoginSettings, listAuthenticationMethodTypes, searchUsers, getUserByID } = await import("../zitadel");
    const { createSessionAndUpdateCookie, setSessionAndUpdateCookie } = await import("@/lib/server/cookie");
    const { completeFlowOrGetUrl } = await import("../client");
    const { checkMFAFactors } = await import("../verify-helper");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetSessionCookieByLoginName = vi.mocked(getSessionCookieByLoginName);
    mockSearchUsers = vi.mocked(searchUsers);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockCreateSessionAndUpdateCookie = vi.mocked(createSessionAndUpdateCookie);
    mockCompleteFlowOrGetUrl = vi.mocked(completeFlowOrGetUrl);
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockCheckMFAFactors = vi.mocked(checkMFAFactors);
    mockSetSessionAndUpdateCookie = vi.mocked(setSessionAndUpdateCookie);
    mockGetUserByID = vi.mocked(getUserByID);
    mockRecordAuthFailure = vi.mocked((await import("@/lib/metrics")).recordAuthFailure);

    mockHeaders.mockResolvedValue({});
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
  });

  test("should create session and verify password when no session exists", async () => {
    // 1. No existing session
    mockGetSessionCookieByLoginName.mockResolvedValue(null);

    // 2. User exists
    mockSearchUsers.mockResolvedValue({
      result: [{ userId: "user123", type: { case: "human", value: {} }, state: UserState.ACTIVE }],
    });

    // 3. Login settings
    mockGetLoginSettings.mockResolvedValue({
      passwordCheckLifetime: { seconds: BigInt(86400) },
      allowLocalAuthentication: true,
    });

    // 4. Create session success
    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "new-session-id",
        factors: {
          user: { id: "user123", loginName: "testuser", organizationId: "org123" },
          password: { verifiedAt: { seconds: 100 } },
        },
      },
      sessionCookie: {
        id: "new-session-id",
        token: "token123",
        organization: "org123",
      },
    });

    // 5. Auth methods (no MFA)
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD],
    });

    // 6. MFA check passes
    mockCheckMFAFactors.mockResolvedValue(null);

    // 7. Redirect
    mockCompleteFlowOrGetUrl.mockResolvedValue({ redirect: "/next-page" });

    const result = await sendPassword({
      loginName: "testuser",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    // Verification
    expect(mockSearchUsers).toHaveBeenCalledWith(expect.objectContaining({ searchValue: "testuser" }));
    expect(mockCreateSessionAndUpdateCookie).toHaveBeenCalled();
    expect(mockCompleteFlowOrGetUrl).toHaveBeenCalled();
    expect(result).toEqual({ redirect: "/next-page" });
  });

  const existingCookie = { id: "session-id", token: "token123", loginName: "testuser", organization: "orgB" };

  test("refuses password when existing session cookie belongs to an org with local auth disabled", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(existingCookie);
    mockGetLoginSettings.mockResolvedValue({
      passwordCheckLifetime: { seconds: BigInt(86400) },
      allowLocalAuthentication: false,
    });

    const result = await sendPassword({
      loginName: "testuser",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    expect(result).toEqual({ error: "errors.localAuthenticationNotAllowed" });
    expect(mockSetSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("updates existing session when the cookie org allows local auth", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(existingCookie);
    mockGetLoginSettings.mockResolvedValue({
      passwordCheckLifetime: { seconds: BigInt(86400) },
      allowLocalAuthentication: true,
    });
    mockSetSessionAndUpdateCookie.mockResolvedValue({
      id: "session-id",
      factors: {
        user: { id: "user123", loginName: "testuser", organizationId: "orgB" },
        password: { verifiedAt: { seconds: 100 } },
      },
    });
    mockGetUserByID.mockResolvedValue({
      user: { userId: "user123", type: { case: "human", value: {} }, state: UserState.ACTIVE },
    });
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD],
    });
    mockCheckMFAFactors.mockResolvedValue(null);
    mockCompleteFlowOrGetUrl.mockResolvedValue({ redirect: "/next-page" });

    const result = await sendPassword({
      loginName: "testuser",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    expect(mockSetSessionAndUpdateCookie).toHaveBeenCalled();
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(result).toEqual({ redirect: "/next-page" });
  });

  function mockNoSessionUserInOrgB(contextSettings: Record<string, unknown>, userOrgSettings: unknown) {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockSearchUsers.mockResolvedValue({
      result: [
        {
          userId: "user123",
          preferredLoginName: "testuser",
          type: { case: "human", value: {} },
          state: UserState.ACTIVE,
          details: { resourceOwner: "orgB" },
        },
      ],
    });
    mockGetLoginSettings.mockImplementation(async ({ organization }: { organization?: string }) =>
      organization === "orgB" ? userOrgSettings : contextSettings,
    );
  }

  test("refuses password when no session and the user's org has local auth disabled (context org allows)", async () => {
    mockNoSessionUserInOrgB({ allowLocalAuthentication: true }, { allowLocalAuthentication: false });

    const result = await sendPassword({
      loginName: "testuser",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    expect(result).toEqual({ error: "errors.localAuthenticationNotAllowed" });
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("masks the refusal when context org ignores unknown usernames", async () => {
    mockNoSessionUserInOrgB(
      { allowLocalAuthentication: true, ignoreUnknownUsernames: true },
      { allowLocalAuthentication: false },
    );

    const result = await sendPassword({
      loginName: "testuser",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    expect(result).toEqual({ error: "errors.failedToAuthenticateNoLimit" });
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("fails closed when the user's org login settings cannot be loaded", async () => {
    mockNoSessionUserInOrgB({ allowLocalAuthentication: true }, undefined);

    const result = await sendPassword({
      loginName: "testuser",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    expect(result).toEqual({ error: "errors.couldNotVerifyPassword" });
    expect(mockRecordAuthFailure).toHaveBeenCalledWith("password", "settings_unavailable", undefined);
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("masks the failure to load the user's org settings when context org ignores unknown usernames", async () => {
    mockNoSessionUserInOrgB({ allowLocalAuthentication: true, ignoreUnknownUsernames: true }, undefined);

    const result = await sendPassword({
      loginName: "testuser",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    expect(result).toEqual({ error: "errors.failedToAuthenticateNoLimit" });
    expect(mockCreateSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("records the failure when the context org has local auth disabled", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockGetLoginSettings.mockResolvedValue({ allowLocalAuthentication: false });

    const result = await sendPassword({
      loginName: "testuser",
      organization: "orgA",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    expect(result).toEqual({ error: "errors.localAuthenticationNotAllowed" });
    expect(mockRecordAuthFailure).toHaveBeenCalledWith("password", "local_auth_not_allowed", "orgA");
    expect(mockSearchUsers).not.toHaveBeenCalled();
  });

  test("uses the user's org login settings for the later checks", async () => {
    const userOrgSettings = { allowLocalAuthentication: true };
    mockNoSessionUserInOrgB({ allowLocalAuthentication: true }, userOrgSettings);
    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: {
        id: "new-session-id",
        factors: {
          user: { id: "user123", loginName: "testuser", organizationId: "orgB" },
          password: { verifiedAt: { seconds: 100 } },
        },
      },
      sessionCookie: { id: "new-session-id", token: "token123", organization: "orgB" },
    });
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD],
    });
    mockCheckMFAFactors.mockResolvedValue(null);
    mockCompleteFlowOrGetUrl.mockResolvedValue({ redirect: "/next-page" });

    const result = await sendPassword({
      loginName: "testuser",
      organization: "orgA",
      checks: create(ChecksSchema, { password: { password: "password123" } }),
    });

    expect(result).toEqual({ redirect: "/next-page" });
    expect(mockGetLoginSettings).toHaveBeenCalledTimes(2);
    expect(mockCheckMFAFactors).toHaveBeenCalledWith(
      expect.anything(),
      expect.anything(),
      userOrgSettings,
      expect.anything(),
      "orgA",
      undefined,
    );
  });
});

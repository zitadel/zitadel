import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { changePassword, checkSessionAndSetPassword, resetPassword, sendPassword } from "./password";

// Mock dependencies
vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("@zitadel/client", () => ({
  create: vi.fn(),
  ConnectError: class extends Error {
    code: number;
    constructor(msg: string, code: number) {
      super(msg);
      this.code = code;
    }
  },
  timestampDate: (ts: any) => new Date(ts.seconds * 1000),
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
  getLockoutSettings: vi.fn(),
  passwordReset: vi.fn(),
  getUserByID: vi.fn(),
  setUserPassword: vi.fn(),
  getPasswordExpirySettings: vi.fn(),
  searchUsers: vi.fn(),
}));

vi.mock("./cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("../cookies", () => ({
  getSessionCookieById: vi.fn(),
  getSessionCookieByLoginName: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("../verify-helper", () => ({
  checkEmailVerification: vi.fn(),
  checkMFAFactors: vi.fn(),
  checkPasswordChangeRequired: vi.fn(),
  checkUserVerification: vi.fn(),
}));

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

describe("checkSessionAndSetPassword", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetSessionCookieById: any;
  let mockGetSession: any;
  let mockListAuthenticationMethodTypes: any;
  let mockGetLoginSettings: any;
  let mockGetLockoutSettings: any;
  let mockSetPassword: any;
  let mockSetSessionAndUpdateCookie: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getSessionCookieById } = await import("../cookies");
    const { getSession, listAuthenticationMethodTypes, getLoginSettings, getLockoutSettings, setPassword } =
      await import("../zitadel");
    const { setSessionAndUpdateCookie } = await import("./cookie");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetSessionCookieById = vi.mocked(getSessionCookieById);
    mockGetSession = vi.mocked(getSession);
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockGetLockoutSettings = vi.mocked(getLockoutSettings);
    mockSetPassword = vi.mocked(setPassword);
    mockSetSessionAndUpdateCookie = vi.mocked(setSessionAndUpdateCookie);

    mockHeaders.mockResolvedValue({});
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });

    // Default session setup
    mockGetSessionCookieById.mockResolvedValue({ id: "session123", token: "token123", organization: "org123" });
    mockGetSession.mockResolvedValue({
      session: {
        factors: {
          user: { id: "user123", organizationId: "org123" },
          password: { verifiedAt: { seconds: Math.floor(Date.now() / 1000) } },
        },
      },
    });

    // Default: Only password method
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD],
    });

    // Default: login settings with password check lifetime
    mockGetLoginSettings.mockResolvedValue({
      forceMfa: false,
      passwordCheckLifetime: { seconds: BigInt(60 * 60 * 24), nanos: 0 },
    });

    // Default: current password verification succeeds
    mockSetSessionAndUpdateCookie.mockResolvedValue({});

    mockSetPassword.mockResolvedValue({});
  });

  test("should verify current password and set new password", async () => {
    await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "oldpassword",
      password: "newpassword",
    });

    expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({
        recentCookie: expect.objectContaining({ id: "session123" }),
      }),
    );
    expect(mockSetPassword).toHaveBeenCalled();
  });

  test("should return error when current password is incorrect", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue(new Error("invalid password"));

    const result = await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "wrongpassword",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "change.errors.currentPasswordInvalid" });
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should return generic error when current password is incorrect and ignoreUnknownUsernames is true", async () => {
    mockGetLoginSettings.mockResolvedValue({
      forceMfa: false,
      passwordCheckLifetime: { seconds: BigInt(60 * 60 * 24), nanos: 0 },
      ignoreUnknownUsernames: true,
    });
    mockSetSessionAndUpdateCookie.mockRejectedValue(new Error("invalid password"));

    const result = await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "wrongpassword",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "change.errors.couldNotVerifyPassword" });
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should pass requestId from session cookie to setSessionAndUpdateCookie", async () => {
    mockGetSessionCookieById.mockResolvedValue({
      id: "session123",
      token: "token123",
      organization: "org123",
      requestId: "oidc-request-456",
    });

    await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "oldpassword",
      password: "newpassword",
    });

    expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({
        requestId: "oidc-request-456",
      }),
    );
  });

  test("should return lockout error with attempt counts when failedAttempts error is thrown", async () => {
    mockSetSessionAndUpdateCookie.mockRejectedValue({ failedAttempts: 3 });
    mockGetLockoutSettings.mockResolvedValue({ maxPasswordAttempts: BigInt(5) });

    const result = await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "wrongpassword",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "errors.failedToAuthenticate" });
    expect(mockGetLockoutSettings).toHaveBeenCalled();
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should return generic error for failedAttempts when ignoreUnknownUsernames is true", async () => {
    mockGetLoginSettings.mockResolvedValue({
      forceMfa: false,
      passwordCheckLifetime: { seconds: BigInt(60 * 60 * 24), nanos: 0 },
      ignoreUnknownUsernames: true,
    });
    mockSetSessionAndUpdateCookie.mockRejectedValue({ failedAttempts: 3 });

    const result = await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "wrongpassword",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "errors.failedToAuthenticateNoLimit" });
    expect(mockGetLockoutSettings).not.toHaveBeenCalled();
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should return error when session cookie not found", async () => {
    mockGetSessionCookieById.mockResolvedValue(null);

    const result = await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "oldpassword",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "Could not load session cookie" });
    expect(mockSetSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should return error when session has no user", async () => {
    mockGetSession.mockResolvedValue({
      session: { factors: {} },
    });

    const result = await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "oldpassword",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "errors.couldNotLoadSession" });
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should use default lifetime when login settings have no passwordCheckLifetime", async () => {
    mockGetLoginSettings.mockResolvedValue({ forceMfa: false });

    await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "oldpassword",
      password: "newpassword",
    });

    expect(mockSetSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({
        lifetime: expect.objectContaining({ seconds: BigInt(60 * 60 * 24) }),
      }),
    );
    expect(mockSetPassword).toHaveBeenCalled();
  });

  test("should handle setPassword failure with failed precondition", async () => {
    mockSetPassword.mockRejectedValue({ code: 9, message: "User is not yet initialized" });

    const result = await checkSessionAndSetPassword({
      sessionId: "session123",
      currentPassword: "oldpassword",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "errors.failedPrecondition" });
  });
});

describe("sendPassword", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetSessionCookieByLoginName: any;
  let mockSearchUsers: any;
  let mockGetLoginSettings: any;
  let mockCreateSessionAndUpdateCookie: any;
  let mockSetSessionAndUpdateCookie: any;
  let mockGetLockoutSettings: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getSessionCookieByLoginName } = await import("../cookies");
    const { getLoginSettings, getLockoutSettings, searchUsers } = await import("../zitadel");
    const { createSessionAndUpdateCookie, setSessionAndUpdateCookie } = await import("./cookie");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetSessionCookieByLoginName = vi.mocked(getSessionCookieByLoginName);
    mockSearchUsers = vi.mocked(searchUsers);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockCreateSessionAndUpdateCookie = vi.mocked(createSessionAndUpdateCookie);
    mockSetSessionAndUpdateCookie = vi.mocked(setSessionAndUpdateCookie);
    mockCreateSessionAndUpdateCookie = vi.mocked(createSessionAndUpdateCookie);
    mockSetSessionAndUpdateCookie = vi.mocked(setSessionAndUpdateCookie);
    mockGetLockoutSettings = vi.mocked(getLockoutSettings);

    const { completeFlowOrGetUrl } = await import("../client");
    vi.mocked(completeFlowOrGetUrl).mockResolvedValue({ redirect: "https://example.com" });

    // eslint-disable-next-line
    const verifyHelper = await import("../verify-helper");
    vi.mocked(verifyHelper.checkPasswordChangeRequired).mockReturnValue(undefined);
    vi.mocked(verifyHelper.checkEmailVerification).mockReturnValue(undefined);
    vi.mocked(verifyHelper.checkMFAFactors).mockResolvedValue(undefined);

    mockHeaders.mockResolvedValue({});
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
  });

  test("should return generic error when user not found and ignoreUnknownUsernames is true", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockSearchUsers.mockResolvedValue({ result: [] });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });

    const result = await sendPassword({
      loginName: "unknown@example.com",
      checks: { password: { password: "password" } } as any,
    });

    expect(result).toEqual({ error: "errors.failedToAuthenticateNoLimit" });
  });

  test("should return specific error when user not found and ignoreUnknownUsernames is false", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockSearchUsers.mockResolvedValue({ result: [] });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: false });

    const result = await sendPassword({
      loginName: "unknown@example.com",
      checks: { password: { password: "password" } } as any,
    });

    expect(result).toEqual({ error: "errors.couldNotVerifyPassword" });
  });

  test("should return generic error when password verification fails and ignoreUnknownUsernames is true", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockSearchUsers.mockResolvedValue({
      result: [{ userId: "user123", type: { case: "human", value: {} }, state: 1 }],
    });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });
    mockCreateSessionAndUpdateCookie.mockRejectedValue({ failedAttempts: 1 });

    const result = await sendPassword({
      loginName: "user@example.com",
      checks: { password: { password: "wrong" } } as any,
    });

    expect(result).toEqual({ error: "errors.failedToAuthenticateNoLimit" });
  });

  test("should return generic error when session creation fails with unknown error and ignoreUnknownUsernames is true", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockSearchUsers.mockResolvedValue({
      result: [{ userId: "user123", type: { case: "human", value: {} }, state: 1 }],
    });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });
    // Simulate an error that is NOT a failed attempt error (e.g. database error)
    mockCreateSessionAndUpdateCookie.mockRejectedValue(new Error("Some internal error"));

    const result = await sendPassword({
      loginName: "user@example.com",
      checks: { password: { password: "correct" } } as any,
    });

    expect(result).toEqual({ error: "errors.failedToAuthenticateNoLimit" });
  });

  test("should return specific error with lockout info when password verification fails and ignoreUnknownUsernames is false", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockSearchUsers.mockResolvedValue({
      result: [{ userId: "user123", type: { case: "human", value: {} }, state: 1 }],
    });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: false });
    mockCreateSessionAndUpdateCookie.mockRejectedValue({ failedAttempts: 1 });
    mockGetLockoutSettings.mockResolvedValue({ maxPasswordAttempts: BigInt(5) });

    const result = await sendPassword({
      loginName: "user@example.com",
      checks: { password: { password: "wrong" } } as any,
    });

    expect(result).toEqual({
      error: "errors.failedToAuthenticate",
    });
  });

  test("should recreate session when session verification fails with keys/session terminated error and ignoreUnknownUsernames is true", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue({
      id: "session123",
      token: "token123",
      organization: "org123",
    });

    const terminatedError = { message: "session already terminated" };
    mockSetSessionAndUpdateCookie.mockRejectedValue(terminatedError);

    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });

    mockSearchUsers.mockResolvedValue({
      result: [
        {
          userId: "user123",
          type: {
            case: "human",
            value: {
              email: { email: "user@example.com", isVerified: true },
              phone: { phone: "+1234567890", isVerified: true },
            },
          },
          state: 1, // Active
          preferredLoginName: "user@example.com",
        },
      ],
    });

    mockCreateSessionAndUpdateCookie.mockResolvedValue({
      session: { factors: { user: { id: "user123", loginName: "user@example.com" } } },
      sessionCookie: { id: "newSession", token: "newToken" },
    });

    // Execute
    await sendPassword({
      loginName: "user@example.com",
      checks: { password: { password: "password" } } as any,
    });

    expect(mockSetSessionAndUpdateCookie).toHaveBeenCalled();
    expect(mockCreateSessionAndUpdateCookie).toHaveBeenCalled();
  });
});
describe("resetPassword", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockSearchUsers: any;
  let mockGetLoginSettings: any;
  let mockPasswordReset: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getLoginSettings, passwordReset, searchUsers } = await import("../zitadel");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockSearchUsers = vi.mocked(searchUsers);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockPasswordReset = vi.mocked(passwordReset);

    mockHeaders.mockResolvedValue({ get: vi.fn(() => "example.com") });
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
  });

  test("should return generic success when user not found and ignoreUnknownUsernames is true", async () => {
    mockSearchUsers.mockResolvedValue({ result: [] });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });

    const result = await resetPassword({
      loginName: "unknown@example.com",
    });

    expect(result).toEqual({});
    expect(mockPasswordReset).not.toHaveBeenCalled();
  });

  test("should return specific error when user not found and ignoreUnknownUsernames is false", async () => {
    mockSearchUsers.mockResolvedValue({ result: [] });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: false });

    const result = await resetPassword({
      loginName: "unknown@example.com",
    });

    expect(result).toEqual({ error: "errors.couldNotSendResetLink" });
    expect(mockPasswordReset).not.toHaveBeenCalled();
  });
});

describe("changePassword", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetUserByID: any;
  let mockGetLoginSettings: any;
  let mockSetUserPassword: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getUserByID, getLoginSettings, setUserPassword } = await import("../zitadel");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetUserByID = vi.mocked(getUserByID);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockSetUserPassword = vi.mocked(setUserPassword);

    mockHeaders.mockResolvedValue({ get: vi.fn(() => "example.com") });
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
  });

  test("should return generic error when user not found and ignoreUnknownUsernames is true", async () => {
    mockGetUserByID.mockResolvedValue({}); // User not found
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });

    const result = await changePassword({
      userId: "unknown",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "set.errors.couldNotSetPassword" });
    expect(mockSetUserPassword).not.toHaveBeenCalled();
  });

  test("should return specific error when user not found and ignoreUnknownUsernames is false", async () => {
    mockGetUserByID.mockResolvedValue({}); // User not found
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: false });

    const result = await changePassword({
      userId: "unknown",
      password: "newpassword",
    });

    expect(result).toEqual({ error: "errors.couldNotSendResetLink" });
    expect(mockSetUserPassword).not.toHaveBeenCalled();
  });
});

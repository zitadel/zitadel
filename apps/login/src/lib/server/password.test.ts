import { describe, expect, test, vi, beforeEach } from "vitest";
import { checkSessionAndSetPassword, sendPassword, resetPassword, changePassword } from "./password";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";

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
  passwordReset: vi.fn(),
  getUserByID: vi.fn(),
  setUserPassword: vi.fn(),
  getPasswordExpirySettings: vi.fn(),
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
  let mockSetPassword: any; // Service account
  let mockCreateUserServiceClient: any; // User session
  let mockSetPasswordUser: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getSessionCookieById } = await import("../cookies");
    const { getSession, listAuthenticationMethodTypes, getLoginSettings, setPassword } = await import("../zitadel");
    const { createUserServiceClient } = await import("@zitadel/client/v2");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetSessionCookieById = vi.mocked(getSessionCookieById);
    mockGetSession = vi.mocked(getSession);
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockSetPassword = vi.mocked(setPassword);
    mockCreateUserServiceClient = vi.mocked(createUserServiceClient);

    mockHeaders.mockResolvedValue({});
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });

    // Default session setup
    mockGetSessionCookieById.mockResolvedValue({ id: "session123", token: "token123" });
    mockGetSession.mockResolvedValue({
      session: {
        factors: {
          user: { id: "user123", organizationId: "org123" },
          password: { verifiedAt: { seconds: Math.floor(Date.now() / 1000) } }, // Password verified recently
        },
      },
    });

    // Default: Only password method
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD],
    });

    // Default: No forced MFA
    mockGetLoginSettings.mockResolvedValue({ forceMfa: false });

    // Mock user service client
    mockSetPasswordUser = vi.fn().mockResolvedValue({});
    mockCreateUserServiceClient.mockReturnValue({
      setPassword: mockSetPasswordUser,
    });

    mockSetPassword.mockResolvedValue({});
  });

  test("should use user session when no MFA is configured", async () => {
    await checkSessionAndSetPassword({ sessionId: "session123", password: "newpassword" });

    expect(mockCreateUserServiceClient).toHaveBeenCalled();
    expect(mockSetPasswordUser).toHaveBeenCalled();
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should use service account when MFA is configured but NOT verified in session", async () => {
    // User has TOTP configured
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.TOTP],
    });

    // Session only has password verified (no OTP)
    const now = Math.floor(Date.now() / 1000);
    mockGetSession.mockResolvedValue({
      session: {
        factors: {
          user: { id: "user123", organizationId: "org123" },
          password: { verifiedAt: { seconds: now - 60 } }, // Verified 1 minute ago
          // otp missing
        },
      },
    });

    await checkSessionAndSetPassword({ sessionId: "session123", password: "newpassword" });

    // EXPECTATION: Should use service account (mockSetPassword)
    // CURRENTLY: Will fail this test and use user session
    expect(mockSetPassword).toHaveBeenCalled();
    expect(mockCreateUserServiceClient).not.toHaveBeenCalled();
  });

  test("should use user session when MFA is configured AND verified in session", async () => {
    // User has TOTP configured
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.TOTP],
    });

    // Session has both password and OTP verified
    mockGetSession.mockResolvedValue({
      session: {
        factors: {
          user: { id: "user123", organizationId: "org123" },
          password: { verifiedAt: { seconds: Math.floor(Date.now() / 1000) } },
          totp: { verifiedAt: { seconds: Math.floor(Date.now() / 1000) } }, // Verified
        },
      },
    });

    await checkSessionAndSetPassword({ sessionId: "session123", password: "newpassword" });

    expect(mockCreateUserServiceClient).toHaveBeenCalled();
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should fail when MFA is configured but not verified, and password verification is too old", async () => {
    // User has TOTP configured
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.TOTP],
    });

    // Session has password verified 10 minutes ago (600 seconds)
    const now = Math.floor(Date.now() / 1000);
    mockGetSession.mockResolvedValue({
      session: {
        factors: {
          user: { id: "user123", organizationId: "org123" },
          password: { verifiedAt: { seconds: now - 600 } },
          // otp missing
        },
      },
    });

    const result = await checkSessionAndSetPassword({ sessionId: "session123", password: "newpassword" });

    expect(result).toEqual({ error: "errors.passwordVerificationTooOld" });
    expect(mockSetPassword).not.toHaveBeenCalled();
  });

  test("should use service account when MFA is configured but not verified, and password verification is recent", async () => {
    // User has TOTP configured
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.TOTP],
    });

    // Session has password verified 1 minute ago (60 seconds)
    const now = Math.floor(Date.now() / 1000);
    mockGetSession.mockResolvedValue({
      session: {
        factors: {
          user: { id: "user123", organizationId: "org123" },
          password: { verifiedAt: { seconds: now - 60 } },
          // otp missing
        },
      },
    });

    await checkSessionAndSetPassword({ sessionId: "session123", password: "newpassword" });

    expect(mockSetPassword).toHaveBeenCalled();
  });
});

describe("sendPassword", () => {
  let mockHeaders: any;
  let mockGetServiceConfig: any;
  let mockGetSessionCookieByLoginName: any;
  let mockListUsers: any;
  let mockGetLoginSettings: any;
  let mockCreateSessionAndUpdateCookie: any;
  let mockSetSessionAndUpdateCookie: any;
  let mockGetLockoutSettings: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { getSessionCookieByLoginName } = await import("../cookies");
    const { listUsers, getLoginSettings, getLockoutSettings } = await import("../zitadel");
    const { createSessionAndUpdateCookie, setSessionAndUpdateCookie } = await import("./cookie");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockGetSessionCookieByLoginName = vi.mocked(getSessionCookieByLoginName);
    mockListUsers = vi.mocked(listUsers);
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
    mockListUsers.mockResolvedValue({ details: { totalResult: BigInt(0) }, result: [] });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });

    const result = await sendPassword({
      loginName: "unknown@example.com",
      checks: { password: { password: "password" } } as any,
    });

    expect(result).toEqual({ error: "errors.failedToAuthenticateNoLimit" });
  });

  test("should return specific error when user not found and ignoreUnknownUsernames is false", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockListUsers.mockResolvedValue({ details: { totalResult: BigInt(0) }, result: [] });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: false });

    const result = await sendPassword({
      loginName: "unknown@example.com",
      checks: { password: { password: "password" } } as any,
    });

    expect(result).toEqual({ error: "errors.couldNotVerifyPassword" });
  });

  test("should return generic error when password verification fails and ignoreUnknownUsernames is true", async () => {
    mockGetSessionCookieByLoginName.mockResolvedValue(null);
    mockListUsers.mockResolvedValue({
      details: { totalResult: BigInt(1) },
      result: [{ userId: "user123" }],
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
    mockListUsers.mockResolvedValue({
      details: { totalResult: BigInt(1) },
      result: [{ userId: "user123" }],
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
    mockListUsers.mockResolvedValue({
      details: { totalResult: BigInt(1) },
      result: [{ userId: "user123" }],
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

    mockListUsers.mockResolvedValue({
      details: { totalResult: BigInt(1) },
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
  let mockListUsers: any;
  let mockGetLoginSettings: any;
  let mockPasswordReset: any;

  beforeEach(async () => {
    vi.clearAllMocks();

    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("../service-url");
    const { listUsers, getLoginSettings, passwordReset } = await import("../zitadel");

    mockHeaders = vi.mocked(headers);
    mockGetServiceConfig = vi.mocked(getServiceConfig);
    mockListUsers = vi.mocked(listUsers);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    mockPasswordReset = vi.mocked(passwordReset);

    mockHeaders.mockResolvedValue({ get: vi.fn(() => "example.com") });
    mockGetServiceConfig.mockReturnValue({ serviceConfig: { baseUrl: "https://api.example.com" } });
  });

  test("should return generic success when user not found and ignoreUnknownUsernames is true", async () => {
    mockListUsers.mockResolvedValue({ details: { totalResult: BigInt(0) }, result: [] });
    mockGetLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });

    const result = await resetPassword({
      loginName: "unknown@example.com",
    });

    expect(result).toEqual({});
    expect(mockPasswordReset).not.toHaveBeenCalled();
  });

  test("should return specific error when user not found and ignoreUnknownUsernames is false", async () => {
    mockListUsers.mockResolvedValue({ details: { totalResult: BigInt(0) }, result: [] });
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

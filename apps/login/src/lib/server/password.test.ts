import { describe, expect, test, vi, beforeEach } from "vitest";
import { checkSessionAndSetPassword } from "./password";
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
}));

vi.mock("../cookies", () => ({
  getSessionCookieById: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
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
          password: { verifiedAt: { seconds: 100 } }, // Password verified
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
          password: { verifiedAt: { seconds: 100 } },
          totp: { verifiedAt: { seconds: 100 } }, // Verified
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

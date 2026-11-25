import { ClassifiedConnectError } from "@/lib/server/grpc/interceptors/error-classification";
import { Code, ConnectError } from "@connectrpc/connect";
import { PasskeysType, SecondFactorType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { getEnrollmentAuthorizationError, getSetupAuthorizationError } from "./enrollment-guard";

vi.mock("@zitadel/client", () => ({
  // Convert a {seconds} Timestamp-like object into a real Date for expiry checks.
  timestampDate: vi.fn((ts: { seconds: bigint | number }) => new Date(Number(ts.seconds) * 1000)),
}));

vi.mock("@/lib/server/zitadel", () => ({
  listAuthenticationMethodTypes: vi.fn(),
  getLoginSettings: vi.fn(),
}));

vi.mock("@/lib/server/verify-helper", () => ({
  checkUserVerification: vi.fn(),
}));

const serviceConfig = { baseUrl: "https://example.com" } as any;

const apiError = (code: Code, message: string) => new ClassifiedConnectError(new ConnectError(message, code));

// A timestamp far in the future / past, expressed the way the proto client shapes them.
const future = { seconds: BigInt(Math.floor(Date.now() / 1000) + 3600), nanos: 0 };
const past = { seconds: BigInt(Math.floor(Date.now() / 1000) - 3600), nanos: 0 };

describe("getEnrollmentAuthorizationError", () => {
  let mockListAuthenticationMethodTypes: any;
  let mockCheckUserVerification: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const { listAuthenticationMethodTypes } = await import("@/lib/server/zitadel");
    const { checkUserVerification } = await import("@/lib/server/verify-helper");
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockCheckUserVerification = vi.mocked(checkUserVerification);
  });

  test("allows enrollment when the session has a verified password and the user has no second factor", async () => {
    // e.g. the forced first-time MFA setup right after the password step
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD],
    });

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: { factors: { user: { id: "user-1" }, password: { verifiedAt: future } } } as any,
    });

    expect(result).toBeNull();
    // an authenticated session must not trigger the user-verification fallback
    expect(mockCheckUserVerification).not.toHaveBeenCalled();
  });

  test.each([
    ["TOTP", AuthenticationMethodType.TOTP],
    ["OTP SMS", AuthenticationMethodType.OTP_SMS],
    ["OTP email", AuthenticationMethodType.OTP_EMAIL],
    ["U2F", AuthenticationMethodType.U2F],
  ])(
    "rejects a password-only session when the user has %s configured",
    async (_name: string, method: AuthenticationMethodType) => {
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD, method],
      });

      const result = await getEnrollmentAuthorizationError({
        serviceConfig,
        userId: "user-1",
        session: { factors: { user: { id: "user-1" }, password: { verifiedAt: future } } } as any,
      });

      expect(result).toBe("You have to verify your second factor first");
      expect(mockCheckUserVerification).not.toHaveBeenCalled();
    },
  );

  test("rejects an IDP-only session when the user has a second factor configured", async () => {
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.IDP, AuthenticationMethodType.TOTP],
    });

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: { factors: { user: { id: "user-1" }, intent: { verifiedAt: future } } } as any,
    });

    expect(result).toBe("You have to verify your second factor first");
  });

  test.each([
    ["TOTP", AuthenticationMethodType.TOTP, { totp: { verifiedAt: future } }],
    ["OTP SMS", AuthenticationMethodType.OTP_SMS, { otpSms: { verifiedAt: future } }],
    ["OTP email", AuthenticationMethodType.OTP_EMAIL, { otpEmail: { verifiedAt: future } }],
    ["U2F", AuthenticationMethodType.U2F, { webAuthN: { verifiedAt: future, userVerified: false } }],
  ])(
    "allows enrollment when the session verified the configured %s",
    async (_name: string, method: AuthenticationMethodType, factor: object) => {
      mockListAuthenticationMethodTypes.mockResolvedValue({
        authMethodTypes: [AuthenticationMethodType.PASSWORD, method],
      });

      const result = await getEnrollmentAuthorizationError({
        serviceConfig,
        userId: "user-1",
        session: { factors: { user: { id: "user-1" }, password: { verifiedAt: future }, ...factor } } as any,
      });

      expect(result).toBeNull();
    },
  );

  test("rejects when the session verified a second factor the user does not have configured", async () => {
    // a verified OTP email check must not stand in for the user's configured TOTP
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.TOTP],
    });

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: {
        factors: { user: { id: "user-1" }, password: { verifiedAt: future }, otpEmail: { verifiedAt: future } },
      } as any,
    });

    expect(result).toBe("You have to verify your second factor first");
  });

  test("allows enrollment when the session has a user-verified passkey factor", async () => {
    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: { factors: { user: { id: "user-1" }, webAuthN: { verifiedAt: future, userVerified: true } } } as any,
    });

    expect(result).toBeNull();
    // a user-verified passkey is multi-factor on its own, no method lookup needed
    expect(mockListAuthenticationMethodTypes).not.toHaveBeenCalled();
  });

  test("rejects a presence-only passkey factor (userVerified false, user has auth methods)", async () => {
    // A WebAuthn assertion without user-verification is a second-factor (U2F) check, not
    // primary authentication — it must not authorize attaching a new authenticator.
    mockListAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [1] });

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: { factors: { user: { id: "user-1" }, webAuthN: { verifiedAt: future, userVerified: false } } } as any,
    });

    expect(result).toBe("You have to authenticate or have a valid User Verification Check");
  });

  test("allows enrollment when the session has a verified IDP intent factor", async () => {
    mockListAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [AuthenticationMethodType.IDP] });

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: { factors: { user: { id: "user-1" }, intent: { verifiedAt: future } } } as any,
    });

    expect(result).toBeNull();
  });

  test("rejects an expired session even if a primary factor was verified (user has auth methods)", async () => {
    mockListAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [1] });

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: {
        factors: { user: { id: "user-1" }, password: { verifiedAt: past } },
        expirationDate: past,
      } as any,
    });

    expect(result).toBe("You have to authenticate or have a valid User Verification Check");
  });

  test("rejects an identify-only session when the user already has auth methods", async () => {
    mockListAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [1] });

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "victim-1",
      session: { factors: { user: { id: "victim-1" } } } as any,
    });

    expect(result).toBe("You have to authenticate or have a valid User Verification Check");
    expect(mockCheckUserVerification).not.toHaveBeenCalled();
  });

  test("rejects an identify-only session with no auth methods when user verification is missing", async () => {
    mockListAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [] });
    mockCheckUserVerification.mockResolvedValue(false);

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "new-user",
      session: { factors: { user: { id: "new-user" } } } as any,
    });

    expect(result).toBe("User Verification Check has to be done");
    expect(mockCheckUserVerification).toHaveBeenCalledWith("new-user");
  });

  test("allows onboarding: no auth methods and a valid prior user-verification check", async () => {
    mockListAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [] });
    mockCheckUserVerification.mockResolvedValue(true);

    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "new-user",
      session: { factors: { user: { id: "new-user" } } } as any,
    });

    expect(result).toBeNull();
  });

  describe("when the authentication methods cannot be listed", () => {
    const sessions = [
      ["an identify-only session", { factors: { user: { id: "user-1" } } }],
      ["an authenticated session", { factors: { user: { id: "user-1" }, password: { verifiedAt: future } } }],
    ] as const;

    test.each(sessions)("fails closed for %s when ZITADEL rejects the call", async (_, session) => {
      // e.g. the user was deleted mid-flow
      mockListAuthenticationMethodTypes.mockRejectedValue(apiError(Code.NotFound, "User not found (QUERY-Dfbg2)"));
      // would allow the onboarding path, if the guard failed open
      mockCheckUserVerification.mockResolvedValue(true);

      const result = await getEnrollmentAuthorizationError({ serviceConfig, userId: "user-1", session: session as any });

      expect(result).toBe("You have to authenticate or have a valid User Verification Check");
      expect(mockCheckUserVerification).not.toHaveBeenCalled();
    });

    test.each([
      ["a server fault", apiError(Code.Internal, "Internal error")],
      ["the API rejecting the login's own credentials", apiError(Code.Unauthenticated, "Token invalid")],
    ])("rethrows %s", async (_, error) => {
      mockListAuthenticationMethodTypes.mockRejectedValue(error);

      for (const [, session] of sessions) {
        await expect(
          getEnrollmentAuthorizationError({ serviceConfig, userId: "user-1", session: session as any }),
        ).rejects.toBe(error);
      }
    });
  });
});

describe("getSetupAuthorizationError", () => {
  let mockListAuthenticationMethodTypes: any;
  let mockGetLoginSettings: any;

  const passwordSession = {
    factors: { user: { id: "user-1", organizationId: "org-1" }, password: { verifiedAt: future } },
  } as any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const { listAuthenticationMethodTypes, getLoginSettings } = await import("@/lib/server/zitadel");
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockGetLoginSettings = vi.mocked(getLoginSettings);
    // password only, no second factor: the enrollment guard allows it
    mockListAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [AuthenticationMethodType.PASSWORD] });
  });

  test("returns the enrollment guard error without loading the login settings", async () => {
    mockListAuthenticationMethodTypes.mockResolvedValue({
      authMethodTypes: [AuthenticationMethodType.PASSWORD, AuthenticationMethodType.TOTP],
    });

    const result = await getSetupAuthorizationError({
      serviceConfig,
      session: passwordSession,
      userId: "user-1",
      method: "passkey",
    });

    expect(result).toBe("You have to verify your second factor first");
    expect(mockGetLoginSettings).not.toHaveBeenCalled();
  });

  test("allows a passkey when passkeys are allowed, using the settings of the user's organization", async () => {
    mockGetLoginSettings.mockResolvedValue({ passkeysType: PasskeysType.ALLOWED, secondFactors: [] });

    const result = await getSetupAuthorizationError({
      serviceConfig,
      session: passwordSession,
      userId: "user-1",
      method: "passkey",
    });

    expect(result).toBeNull();
    expect(mockGetLoginSettings).toHaveBeenCalledWith({ serviceConfig, organization: "org-1" });
  });

  test("rejects a passkey when passkeys are not allowed", async () => {
    mockGetLoginSettings.mockResolvedValue({ passkeysType: PasskeysType.NOT_ALLOWED, secondFactors: [] });

    const result = await getSetupAuthorizationError({
      serviceConfig,
      session: passwordSession,
      userId: "user-1",
      method: "passkey",
    });

    expect(result).toBe("Passkeys are not allowed");
  });

  test.each([
    ["totp", SecondFactorType.TOTP],
    ["u2f", SecondFactorType.U2F],
    ["otp_sms", SecondFactorType.OTP_SMS],
    ["otp_email", SecondFactorType.OTP_EMAIL],
  ] as const)("allows %s when the login settings allow it", async (method, type) => {
    mockGetLoginSettings.mockResolvedValue({ passkeysType: PasskeysType.NOT_ALLOWED, secondFactors: [type] });

    const result = await getSetupAuthorizationError({ serviceConfig, session: passwordSession, userId: "user-1", method });

    expect(result).toBeNull();
  });

  test.each([["totp"], ["u2f"], ["otp_sms"], ["otp_email"]] as const)(
    "rejects %s when the login settings don't allow it",
    async (method) => {
      // every other second factor is allowed
      const others = [
        SecondFactorType.TOTP,
        SecondFactorType.U2F,
        SecondFactorType.OTP_SMS,
        SecondFactorType.OTP_EMAIL,
      ].filter(
        (t) =>
          t !==
          {
            totp: SecondFactorType.TOTP,
            u2f: SecondFactorType.U2F,
            otp_sms: SecondFactorType.OTP_SMS,
            otp_email: SecondFactorType.OTP_EMAIL,
          }[method],
      );
      mockGetLoginSettings.mockResolvedValue({ passkeysType: PasskeysType.ALLOWED, secondFactors: others });

      const result = await getSetupAuthorizationError({ serviceConfig, session: passwordSession, userId: "user-1", method });

      expect(result).toBe("This second factor is not allowed");
    },
  );

  test("fails closed when the login settings can't be loaded", async () => {
    mockGetLoginSettings.mockResolvedValue(undefined);

    const result = await getSetupAuthorizationError({
      serviceConfig,
      session: passwordSession,
      userId: "user-1",
      method: "totp",
    });

    expect(result).toBe("Could not load login settings");
  });
});

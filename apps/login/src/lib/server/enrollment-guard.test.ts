import { beforeEach, describe, expect, test, vi } from "vitest";
import { getEnrollmentAuthorizationError } from "./enrollment-guard";

vi.mock("@zitadel/client", () => ({
  // Convert a {seconds} Timestamp-like object into a real Date for expiry checks.
  timestampDate: vi.fn((ts: { seconds: bigint | number }) => new Date(Number(ts.seconds) * 1000)),
}));

vi.mock("@/lib/zitadel", () => ({
  listAuthenticationMethodTypes: vi.fn(),
}));

vi.mock("@/lib/verify-helper", () => ({
  checkUserVerification: vi.fn(),
}));

const serviceConfig = { baseUrl: "https://example.com" } as any;

// A timestamp far in the future / past, expressed the way the proto client shapes them.
const future = { seconds: BigInt(Math.floor(Date.now() / 1000) + 3600), nanos: 0 };
const past = { seconds: BigInt(Math.floor(Date.now() / 1000) - 3600), nanos: 0 };

describe("getEnrollmentAuthorizationError", () => {
  let mockListAuthenticationMethodTypes: any;
  let mockCheckUserVerification: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const { listAuthenticationMethodTypes } = await import("@/lib/zitadel");
    const { checkUserVerification } = await import("@/lib/verify-helper");
    mockListAuthenticationMethodTypes = vi.mocked(listAuthenticationMethodTypes);
    mockCheckUserVerification = vi.mocked(checkUserVerification);
  });

  test("allows enrollment when the session has a verified password factor", async () => {
    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: { factors: { user: { id: "user-1" }, password: { verifiedAt: future } } } as any,
    });

    expect(result).toBeNull();
    // an authenticated session must not trigger the fallback checks
    expect(mockListAuthenticationMethodTypes).not.toHaveBeenCalled();
    expect(mockCheckUserVerification).not.toHaveBeenCalled();
  });

  test("allows enrollment when the session has a user-verified passkey factor", async () => {
    const result = await getEnrollmentAuthorizationError({
      serviceConfig,
      userId: "user-1",
      session: { factors: { user: { id: "user-1" }, webAuthN: { verifiedAt: future, userVerified: true } } } as any,
    });

    expect(result).toBeNull();
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
});

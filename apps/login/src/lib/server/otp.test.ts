import { Code, ConnectError } from "@connectrpc/connect";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { ClassifiedConnectError } from "../grpc/interceptors/error-classification";
import { addOTP } from "./otp";

vi.mock("@/lib/zitadel", () => ({
  getSession: vi.fn(),
  addOTPEmail: vi.fn(),
  addOTPSMS: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("../cookies", () => ({
  getSessionCookieById: vi.fn(),
}));

vi.mock("./enrollment-guard", () => ({
  getEnrollmentAuthorizationError: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

const session = {
  session: { id: "session-1", factors: { user: { id: "user-1", loginName: "user@example.com" } } },
};

describe("addOTP", () => {
  let getSessionCookieById: any;
  let getSession: any;
  let addOTPEmail: any;
  let addOTPSMS: any;
  let getEnrollmentAuthorizationError: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    getSessionCookieById = vi.mocked((await import("../cookies")).getSessionCookieById);
    getSession = vi.mocked((await import("@/lib/zitadel")).getSession);
    addOTPEmail = vi.mocked((await import("@/lib/zitadel")).addOTPEmail);
    addOTPSMS = vi.mocked((await import("@/lib/zitadel")).addOTPSMS);
    getEnrollmentAuthorizationError = vi.mocked((await import("./enrollment-guard")).getEnrollmentAuthorizationError);

    getSessionCookieById.mockResolvedValue({ id: "session-1", token: "token-1" });
    getSession.mockResolvedValue(session);
    getEnrollmentAuthorizationError.mockResolvedValue(null);
  });

  test("returns an error without a session cookie", async () => {
    getSessionCookieById.mockResolvedValue(undefined);

    const result = await addOTP({ method: "email", sessionId: "session-1" });

    expect(result).toEqual({ error: "errors.couldNotGetSession" });
    expect(addOTPEmail).not.toHaveBeenCalled();
  });

  test("returns an error when the session has no user", async () => {
    getSession.mockResolvedValue({ session: { id: "session-1", factors: {} } });

    const result = await addOTP({ method: "email", sessionId: "session-1" });

    expect(result).toEqual({ error: "errors.couldNotGetSession" });
    expect(addOTPEmail).not.toHaveBeenCalled();
  });

  test("rejects an unauthorized (identify-only) session without adding the factor", async () => {
    getEnrollmentAuthorizationError.mockResolvedValue("You have to authenticate or have a valid User Verification Check");

    const result = await addOTP({ method: "sms", sessionId: "session-1" });

    expect(result).toEqual({ error: "You have to authenticate or have a valid User Verification Check" });
    expect(getEnrollmentAuthorizationError).toHaveBeenCalledWith({
      serviceConfig: { baseUrl: "https://example.com" },
      session: session.session,
      userId: "user-1",
    });
    expect(addOTPSMS).not.toHaveBeenCalled();
  });

  test("adds OTP via email for the user of the session", async () => {
    addOTPEmail.mockResolvedValue({});

    const result = await addOTP({ method: "email", sessionId: "session-1" });

    expect(result).toEqual({ success: true });
    expect(addOTPEmail).toHaveBeenCalledWith({ serviceConfig: { baseUrl: "https://example.com" }, userId: "user-1" });
    expect(addOTPSMS).not.toHaveBeenCalled();
  });

  test("adds OTP via SMS for the user of the session", async () => {
    addOTPSMS.mockResolvedValue({});

    const result = await addOTP({ method: "sms", sessionId: "session-1" });

    expect(result).toEqual({ success: true });
    expect(addOTPSMS).toHaveBeenCalledWith({ serviceConfig: { baseUrl: "https://example.com" }, userId: "user-1" });
    expect(addOTPEmail).not.toHaveBeenCalled();
  });

  test.each([
    ["email", "addOTPEmail"],
    ["sms", "addOTPSMS"],
  ] as const)("treats an already added OTP via %s as success", async (method, fn) => {
    const add = fn === "addOTPEmail" ? addOTPEmail : addOTPSMS;
    add.mockRejectedValue(new ClassifiedConnectError(new ConnectError("already ready", Code.AlreadyExists)));

    const result = await addOTP({ method, sessionId: "session-1" });

    expect(result).toEqual({ success: true });
  });

  test.each([
    ["email", "addOTPEmail", "errors.emailNotVerified"],
    ["sms", "addOTPSMS", "errors.phoneNotVerified"],
  ] as const)("asks to verify first when OTP via %s is not ready", async (method, fn, expected) => {
    const add = fn === "addOTPEmail" ? addOTPEmail : addOTPSMS;
    add.mockRejectedValue(new ClassifiedConnectError(new ConnectError("not ready", Code.FailedPrecondition)));

    const result = await addOTP({ method, sessionId: "session-1" });

    expect(result).toEqual({ error: expected });
  });

  test.each([
    ["email", "addOTPEmail", "errors.couldNotAddEmail"],
    ["sms", "addOTPSMS", "errors.couldNotAddSms"],
  ] as const)("returns an error when the backend rejects adding OTP via %s", async (method, fn, expected) => {
    const add = fn === "addOTPEmail" ? addOTPEmail : addOTPSMS;
    add.mockRejectedValue(new ClassifiedConnectError(new ConnectError("denied", Code.PermissionDenied)));

    const result = await addOTP({ method, sessionId: "session-1" });

    expect(result).toEqual({ error: expected });
  });

  test("rejects an unknown method", async () => {
    const result = await addOTP({ method: "time-based" as any, sessionId: "session-1" });

    expect(result).toEqual({ error: "errors.invalidMethod" });
    expect(addOTPEmail).not.toHaveBeenCalled();
    expect(addOTPSMS).not.toHaveBeenCalled();
  });
});

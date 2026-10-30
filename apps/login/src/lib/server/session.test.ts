import { Code } from "@connectrpc/connect";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { clearSession, updateOrCreateSession } from "./session";

vi.mock("@zitadel/client", async () => {
  const { Code } = await import("@connectrpc/connect");
  return { Code, create: vi.fn((_schema: unknown, data: unknown) => data), Duration: {} };
});

vi.mock("@/lib/zitadel", () => ({
  deleteSession: vi.fn(),
  getLoginSettings: vi.fn(),
  getSecuritySettings: vi.fn(),
  humanMFAInitSkipped: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  listUsers: vi.fn(),
}));

vi.mock("@/lib/server/cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("../cookies", () => ({
  getMostRecentSessionCookie: vi.fn(),
  getSessionCookieById: vi.fn(),
  getSessionCookieByLoginName: vi.fn(),
  removeSessionFromCookie: vi.fn(),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("../session", () => ({
  isSessionValid: vi.fn(),
}));

vi.mock("./host", () => ({
  getPublicHost: vi.fn(() => "test.com"),
}));

vi.mock("./loginname", () => ({
  sendLoginname: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

vi.mock("@/lib/grpc/interceptors/error-classification", () => ({
  isClassifiedError: (error: unknown) => error !== null && typeof error === "object" && "code" in error,
}));

const cookie = { id: "session-1", token: "token-1", loginName: "user@example.com" };

describe("clearSession", () => {
  let deleteSession: any;
  let getSecuritySettings: any;
  let getSessionCookieById: any;
  let removeSessionFromCookie: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/zitadel");
    const cookies = await import("../cookies");
    deleteSession = vi.mocked(zitadel.deleteSession);
    getSecuritySettings = vi.mocked(zitadel.getSecuritySettings);
    getSessionCookieById = vi.mocked(cookies.getSessionCookieById);
    removeSessionFromCookie = vi.mocked(cookies.removeSessionFromCookie);

    getSessionCookieById.mockResolvedValue(cookie);
    getSecuritySettings.mockResolvedValue({ embeddedIframe: { enabled: true } });
  });

  test("deletes the session and prunes the cookie entry on success", async () => {
    deleteSession.mockResolvedValue({ details: {} });

    const res = await clearSession({ sessionId: "session-1" });

    expect(res).toBeUndefined();
    expect(deleteSession).toHaveBeenCalledWith(expect.objectContaining({ sessionId: "session-1", sessionToken: "token-1" }));
    expect(removeSessionFromCookie).toHaveBeenCalledWith({ session: cookie, iFrameEnabled: true });
  });

  test("prunes the cookie entry when the server rejects the cookie token (session gone or token stale)", async () => {
    deleteSession.mockRejectedValue({ code: Code.PermissionDenied });

    const res = await clearSession({ sessionId: "session-1" });

    expect(res).toBeUndefined();
    expect(removeSessionFromCookie).toHaveBeenCalledWith({ session: cookie, iFrameEnabled: true });
  });

  test("keeps the cookie entry and reports an error on any other failure", async () => {
    deleteSession.mockRejectedValue({ code: Code.Unavailable });

    const res = await clearSession({ sessionId: "session-1" });

    expect(res).toEqual({ error: "couldNotClearSession" });
    expect(removeSessionFromCookie).not.toHaveBeenCalled();
  });

  test("attempts the server-side delete even when security settings cannot be loaded", async () => {
    deleteSession.mockResolvedValue({ details: {} });
    getSecuritySettings.mockRejectedValue(new Error("settings unavailable"));

    await expect(clearSession({ sessionId: "session-1" })).rejects.toThrow("settings unavailable");

    expect(deleteSession).toHaveBeenCalledTimes(1);
    expect(removeSessionFromCookie).not.toHaveBeenCalled();
  });

  test("does nothing when the session id is not in the cookie", async () => {
    getSessionCookieById.mockResolvedValue(undefined);

    const res = await clearSession({ sessionId: "unknown" });

    expect(res).toBeUndefined();
    expect(deleteSession).not.toHaveBeenCalled();
  });
});

describe("updateOrCreateSession OTP SMS resend cooldown", () => {
  let getSessionCookieById: any;
  let setSessionAndUpdateCookie: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    vi.useFakeTimers();
    vi.setSystemTime(1_000_000);

    const zitadel = await import("@/lib/zitadel");
    const cookies = await import("../cookies");
    const serverCookie = await import("@/lib/server/cookie");
    getSessionCookieById = vi.mocked(cookies.getSessionCookieById);
    setSessionAndUpdateCookie = vi.mocked(serverCookie.setSessionAndUpdateCookie);

    vi.mocked(zitadel.getLoginSettings).mockResolvedValue({ secondFactorCheckLifetime: { seconds: BigInt(60) } } as any);
    setSessionAndUpdateCookie.mockResolvedValue({ id: "session-1", factors: {} });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  test("refuses another OTP SMS within the cooldown without calling the API", async () => {
    getSessionCookieById.mockResolvedValue({ ...cookie, otpSmsSentTs: `${1_000_000 - 10_000}` });

    const res = await updateOrCreateSession({ sessionId: "session-1", challenges: { otpSms: {} } as any });

    expect(res).toEqual({ error: "errors.resendCooldown", retryAfterSeconds: 20 });
    expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("sends the OTP SMS once the cooldown has passed and returns the next cooldown", async () => {
    getSessionCookieById.mockResolvedValue({ ...cookie, otpSmsSentTs: `${1_000_000 - 30_000}` });

    const res = await updateOrCreateSession({ sessionId: "session-1", challenges: { otpSms: {} } as any });

    expect(setSessionAndUpdateCookie).toHaveBeenCalledTimes(1);
    expect(res).toEqual(expect.objectContaining({ sessionId: "session-1", retryAfterSeconds: 30 }));
  });

  test("requests the OTP SMS only once when it has to create the session first", async () => {
    const cookies = await import("../cookies");
    const serverCookie = await import("@/lib/server/cookie");
    vi.mocked(cookies.getSessionCookieByLoginName).mockResolvedValue(undefined);
    vi.mocked(serverCookie.createSessionAndUpdateCookie).mockResolvedValue({ sessionCookie: cookie } as any);

    await updateOrCreateSession({ loginName: "user@example.com", challenges: { otpSms: {} } as any });

    expect(vi.mocked(serverCookie.createSessionAndUpdateCookie).mock.calls[0][0]).not.toHaveProperty("challenges");
    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(expect.objectContaining({ challenges: { otpSms: {} } }));
  });

  test("does not apply the cooldown to code checks", async () => {
    getSessionCookieById.mockResolvedValue({ ...cookie, otpSmsSentTs: `${1_000_000}` });

    const res = await updateOrCreateSession({ sessionId: "session-1", checks: { otpSms: { code: "123456" } } as any });

    expect(setSessionAndUpdateCookie).toHaveBeenCalledTimes(1);
    expect(res).toEqual(expect.objectContaining({ sessionId: "session-1", retryAfterSeconds: undefined }));
  });
});

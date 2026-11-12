import { Code } from "@connectrpc/connect";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { clearSession, updateOrCreateSession } from "./session";

const mockLogger = vi.hoisted(() => ({
  debug: vi.fn(),
  info: vi.fn(),
  warn: vi.fn(),
  error: vi.fn(),
}));

vi.mock("@/lib/logger", () => ({
  createLogger: () => mockLogger,
}));

vi.mock("@zitadel/client", async () => {
  const { Code } = await import("@connectrpc/connect");
  return {
    Code,
    create: vi.fn((_schema: unknown, data: unknown) => data),
    Duration: {},
    timestampDate: (timestamp: { seconds: bigint }) => new Date(Number(timestamp.seconds) * 1000),
  };
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

vi.mock("../session", async () => {
  const actual = await vi.importActual<typeof import("../session")>("../session");
  return {
    isSessionValid: vi.fn(),
    reportMissingSessionCookie: vi.fn(),
    resolveSession: vi.fn(),
    hasVerifiedPrimaryFactor: actual.hasVerifiedPrimaryFactor,
    SESSION_RESOLUTION_MESSAGE_KEYS: actual.SESSION_RESOLUTION_MESSAGE_KEYS,
  };
});

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

describe("updateOrCreateSession", () => {
  const sessionCookie = { id: "session-1", token: "token-1", loginName: "user@example.com", organization: "org-1" };
  const inFuture = { seconds: BigInt(Math.floor(Date.now() / 1000) + 3600) };
  const inPast = { seconds: BigInt(Math.floor(Date.now() / 1000) - 3600) };
  const passwordVerified = {
    id: "session-1",
    expirationDate: inFuture,
    factors: { user: { id: "user-1" }, password: { verifiedAt: inPast } },
  };
  const userOnly = { id: "session-1", expirationDate: inFuture, factors: { user: { id: "user-1" } } };

  const totpCheck = { totp: { code: "123456" } } as any;
  const smsChallenge = { otpSms: {} } as any;
  const passkeyChallenge = { webAuthN: { domain: "test.com" } } as any;

  let getSessionCookieByLoginName: any;
  let getSessionCookieById: any;
  let getMostRecentSessionCookie: any;
  let setSessionAndUpdateCookie: any;
  let createSessionAndUpdateCookie: any;
  let listUsers: any;
  let reportMissingSessionCookie: any;
  let resolveSession: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/zitadel");
    const cookies = await import("../cookies");
    const serverCookie = await import("@/lib/server/cookie");
    const sessionLib = await import("../session");

    getSessionCookieByLoginName = vi.mocked(cookies.getSessionCookieByLoginName);
    getSessionCookieById = vi.mocked(cookies.getSessionCookieById);
    getMostRecentSessionCookie = vi.mocked(cookies.getMostRecentSessionCookie);
    setSessionAndUpdateCookie = vi.mocked(serverCookie.setSessionAndUpdateCookie);
    createSessionAndUpdateCookie = vi.mocked(serverCookie.createSessionAndUpdateCookie);
    listUsers = vi.mocked(zitadel.listUsers);
    reportMissingSessionCookie = vi.mocked(sessionLib.reportMissingSessionCookie);
    resolveSession = vi.mocked(sessionLib.resolveSession);

    vi.mocked(zitadel.getLoginSettings).mockResolvedValue({} as any);
    getSessionCookieByLoginName.mockResolvedValue(sessionCookie);
    listUsers.mockResolvedValue({ details: { totalResult: BigInt(1) }, result: [{ userId: "user-1" }] });
    createSessionAndUpdateCookie.mockResolvedValue({
      session: { id: "session-2", factors: { user: { id: "user-1" } } },
      sessionCookie: { ...sessionCookie, id: "session-2" },
      challenges: undefined,
    });
  });

  describe("no cookie entry for the session", () => {
    beforeEach(() => {
      getSessionCookieByLoginName.mockResolvedValue(undefined);
      getSessionCookieById.mockResolvedValue(undefined);
      getMostRecentSessionCookie.mockResolvedValue(undefined);
      reportMissingSessionCookie.mockResolvedValue("cookie_missing");
    });

    test.each([
      ["a TOTP check", { checks: totpCheck }],
      ["an SMS challenge", { challenges: smsChallenge }],
    ])("does not create a session for %s and reports the session as lost", async (_name, step) => {
      const res = await updateOrCreateSession({ loginName: "user@example.com", organization: "org-1", ...step });

      expect(res).toEqual({ error: "sessionResolution.cookieMissing", sessionLost: true });
      expect(reportMissingSessionCookie).toHaveBeenCalledWith("update_session", {
        loginName: "user@example.com",
        organization: "org-1",
        sessionId: undefined,
      });
      expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
      expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test("reports the session as lost when neither a login name nor a cookie entry identifies it", async () => {
      reportMissingSessionCookie.mockResolvedValue("session_entry_not_found");

      const res = await updateOrCreateSession({ sessionId: "session-1", challenges: passkeyChallenge });

      expect(res).toEqual({ error: "sessionResolution.entryNotFound", sessionLost: true });
      expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test("still creates the session that starts a passkey login", async () => {
      setSessionAndUpdateCookie.mockResolvedValue({ id: "session-2", factors: { user: { id: "user-1" } } });

      const res = await updateOrCreateSession({ loginName: "user@example.com", challenges: passkeyChallenge });

      expect(createSessionAndUpdateCookie).toHaveBeenCalledTimes(1);
      expect(reportMissingSessionCookie).not.toHaveBeenCalled();
      expect(res).toMatchObject({ sessionId: "session-2" });
    });
  });

  describe("the session update fails", () => {
    const updateError = { code: Code.InvalidArgument };

    beforeEach(() => {
      setSessionAndUpdateCookie.mockRejectedValue(updateError);
    });

    test("rethrows in a second factor step while the session is intact (e.g. a wrong code)", async () => {
      resolveSession.mockResolvedValue({ session: passwordVerified });

      await expect(updateOrCreateSession({ loginName: "user@example.com", checks: totpCheck })).rejects.toBe(updateError);

      expect(resolveSession).toHaveBeenCalledWith({ baseUrl: "https://example.com" }, sessionCookie, "update_session");
      expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test.each([
      ["session_token_rejected", "sessionResolution.cookieInvalid"],
      ["server_session_not_found", "sessionResolution.sessionNotFound"],
      ["session_lookup_failed", "sessionResolution.lookupFailed"],
    ])("returns %s in a second factor step instead of creating a session", async (reason, message) => {
      resolveSession.mockResolvedValue({ reason });

      const res = await updateOrCreateSession({ loginName: "user@example.com", challenges: smsChallenge });

      expect(res).toEqual({ error: message, sessionLost: true });
      expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
      expect(mockLogger.warn).toHaveBeenCalledWith(
        expect.stringContaining("LOGIN_SESSION_UPDATE_FAILED"),
        expect.objectContaining({
          code: "LOGIN_SESSION_UPDATE_FAILED",
          reason,
          sessionId: "session-1",
          secondFactorStep: true,
          grpcCode: Code.InvalidArgument,
        }),
      );
    });

    test("reports an expired session in a second factor step as ended", async () => {
      resolveSession.mockResolvedValue({ session: { ...passwordVerified, expirationDate: inPast } });

      const res = await updateOrCreateSession({ loginName: "user@example.com", checks: totpCheck });

      expect(res).toEqual({ error: "sessionResolution.sessionNotFound", sessionLost: true });
      expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
      expect(mockLogger.warn).toHaveBeenCalledWith(
        expect.stringContaining("LOGIN_SESSION_UPDATE_FAILED"),
        expect.objectContaining({ expired: true }),
      );
    });

    test("does not replace a valid session with a verified first factor outside a second factor step", async () => {
      resolveSession.mockResolvedValue({ session: passwordVerified });

      await expect(updateOrCreateSession({ loginName: "user@example.com", challenges: passkeyChallenge })).rejects.toBe(
        updateError,
      );

      expect(listUsers).not.toHaveBeenCalled();
      expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test.each([
      ["only holds the user check", { session: userOnly }],
      ["can no longer be resolved", { reason: "server_session_not_found" }],
      // a returning user whose cookie still holds an old session must be able to start a passkey login
      ["has expired", { session: { ...passwordVerified, expirationDate: inPast } }],
    ])("creates a new session when the existing one %s and no second factor is asked for", async (_name, existing) => {
      resolveSession.mockResolvedValue(existing);

      const res = await updateOrCreateSession({ loginName: "user@example.com", challenges: passkeyChallenge });

      expect(createSessionAndUpdateCookie).toHaveBeenCalledTimes(1);
      expect(res).toMatchObject({ sessionId: "session-2" });
    });
  });

  test("updates the session of the cookie entry", async () => {
    setSessionAndUpdateCookie.mockResolvedValue(passwordVerified);

    const res = await updateOrCreateSession({ loginName: "user@example.com", organization: "org-1", checks: totpCheck });

    expect(getSessionCookieByLoginName).toHaveBeenCalledWith({ loginName: "user@example.com", organization: "org-1" });
    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(expect.objectContaining({ recentCookie: sessionCookie }));
    expect(resolveSession).not.toHaveBeenCalled();
    expect(res).toMatchObject({ sessionId: "session-1" });
  });
});

import { getLoginSettings, getSession, listUsers, searchUsers } from "@/lib/zitadel";
import { Code } from "@connectrpc/connect";
import { create } from "@zitadel/client";
import { ChallengesSchema, RequestChallengesSchema } from "@zitadel/proto/zitadel/session/v2/challenge_pb";
import { SessionSchema } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { ChecksSchema, GetSessionResponseSchema } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { LoginSettingsSchema } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { ListUsersResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { getMostRecentSessionCookie, getSessionCookieById, getSessionCookieByLoginName } from "../cookies";
import { createSessionAndUpdateCookie, setSessionAndUpdateCookie } from "./cookie";
import { clearSession, updateOrCreateSession } from "./session";

vi.mock("@zitadel/client", async () => {
  const { Code } = await import("@connectrpc/connect");
  return {
    Code,
    create: vi.fn((_schema: unknown, data: unknown) => data),
    Duration: {},
    timestampMs: (value: { seconds: bigint; nanos: number }) => Number(value.seconds) * 1000 + value.nanos / 1e6,
  };
});

vi.mock("@/lib/zitadel", () => ({
  deleteSession: vi.fn(),
  getLoginSettings: vi.fn(),
  getSecuritySettings: vi.fn(),
  getSession: vi.fn(),
  humanMFAInitSkipped: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  listUsers: vi.fn(),
  searchUsers: vi.fn(),
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

describe("passkey sessions", () => {
  const email = "user@example.com";
  const canonicalLoginName = `${email}@freightcheck.test`;
  const organization = "freightcheck";
  const sessionCookie = {
    ...cookie,
    loginName: canonicalLoginName,
    organization,
    creationTs: "",
    expirationTs: "",
    changeTs: "",
  };
  const session = create(SessionSchema, {
    id: cookie.id,
    factors: { user: { id: "fc-user", loginName: canonicalLoginName, organizationId: organization } },
    expirationDate: { seconds: BigInt(4102444800), nanos: 0 },
  });
  const users = create(ListUsersResponseSchema, {
    details: { totalResult: BigInt(1) },
    result: [
      {
        userId: "fc-user",
        preferredLoginName: canonicalLoginName,
        details: { resourceOwner: organization },
        type: { case: "human", value: { email: { email } } },
      },
    ],
  });
  const challenge = create(RequestChallengesSchema, { webAuthN: { domain: "test.com" } });
  const returnedChallenge = create(ChallengesSchema, {
    webAuthN: { publicKeyCredentialRequestOptions: { publicKey: { challenge: "expected-challenge" } } },
  });
  const assertion = create(ChecksSchema, { webAuthN: { credentialAssertionData: { id: "credential" } } });
  const failure = { error: "couldNotFindSession" };

  beforeEach(() => {
    vi.resetAllMocks();
    vi.mocked(getLoginSettings).mockResolvedValue(create(LoginSettingsSchema, { ignoreUnknownUsernames: true }));
    vi.mocked(searchUsers).mockResolvedValue(users);
    vi.mocked(listUsers).mockResolvedValue(users);
    vi.mocked(getSession).mockResolvedValue(create(GetSessionResponseSchema, { session }));
    vi.mocked(createSessionAndUpdateCookie).mockResolvedValue({ session, sessionCookie, challenges: returnedChallenge });
    vi.mocked(setSessionAndUpdateCookie).mockResolvedValue({ ...session, challenges: returnedChallenge });
  });

  test("creates a fresh scoped session by user ID and returns its challenge without replacing it", async () => {
    const result = await updateOrCreateSession({ loginName: email, organization, challenges: challenge, requestId: "oidc" });
    expect(result).toMatchObject({ sessionId: cookie.id, factors: session.factors, challenges: returnedChallenge });
    expect(searchUsers).toHaveBeenCalledWith(
      expect.objectContaining({
        searchValue: email,
        organizationId: organization,
        loginSettings: expect.objectContaining({ ignoreUnknownUsernames: true }),
      }),
    );
    expect(getSessionCookieByLoginName).toHaveBeenCalledWith({ loginName: canonicalLoginName, organization });
    expect(createSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({
        checks: { user: { search: { case: "userId", value: "fc-user" } } },
        challenges: challenge,
        requestId: "oidc",
      }),
    );
    expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(listUsers).not.toHaveBeenCalled();
  });

  test("reuses the canonical remembered session for a new challenge", async () => {
    vi.mocked(getSessionCookieByLoginName).mockResolvedValue(sessionCookie);
    await updateOrCreateSession({ loginName: email, organization, challenges: challenge });
    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({ recentCookie: sessionCookie, challenges: challenge }),
    );
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("verifies only the explicitly supplied challenge session", async () => {
    vi.mocked(getSessionCookieById).mockResolvedValue(sessionCookie);
    const result = await updateOrCreateSession({ loginName: email, organization, sessionId: cookie.id, checks: assertion });
    expect(result).toMatchObject({ sessionId: cookie.id });
    expect(getSession).toHaveBeenCalledWith(expect.objectContaining({ sessionId: cookie.id, sessionToken: cookie.token }));
    expect(setSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({ recentCookie: sessionCookie, checks: assertion }),
    );
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(getSessionCookieByLoginName).not.toHaveBeenCalled();
  });

  test("rejects an assertion without a session ID even with a remembered session", async () => {
    vi.mocked(getMostRecentSessionCookie).mockResolvedValue(sessionCookie);
    vi.mocked(getSessionCookieByLoginName).mockResolvedValue(sessionCookie);
    expect(await updateOrCreateSession({ loginName: email, organization, checks: assertion })).toEqual(failure);
    expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test.each(["challenge", "assertion"])("rejects a missing cookie for an explicit %s session", async (operation) => {
    expect(
      await updateOrCreateSession({
        loginName: email,
        organization,
        sessionId: "missing",
        ...(operation === "challenge" ? { challenges: challenge } : { checks: assertion }),
      }),
    ).toEqual(failure);
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test.each([{ organization: "foreign" }, { loginName: "someone-else@freightcheck.test" }])(
    "rejects a mismatched signed session cookie: %j",
    async (override) => {
      vi.mocked(getSessionCookieById).mockResolvedValue({ ...sessionCookie, ...override });
      expect(
        await updateOrCreateSession({ loginName: email, organization, sessionId: cookie.id, checks: assertion }),
      ).toEqual(failure);
      expect(getSession).not.toHaveBeenCalled();
      expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
      expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
    },
  );

  test.each([
    { id: "other-session" },
    { factors: { user: { ...session.factors!.user!, id: "other-user" } } },
    { factors: { user: { ...session.factors!.user!, organizationId: "foreign" } } },
    { factors: undefined },
  ])("rejects a missing or mismatched server-side identity: %j", async (override) => {
    vi.mocked(getSessionCookieById).mockResolvedValue(sessionCookie);
    vi.mocked(getSession).mockResolvedValue(
      create(GetSessionResponseSchema, {
        session: create(SessionSchema, {
          id: override.id ?? session.id,
          factors: "factors" in override ? create(SessionSchema, { factors: override.factors }).factors : session.factors,
        }),
      }),
    );
    expect(await updateOrCreateSession({ loginName: email, organization, sessionId: cookie.id, checks: assertion })).toEqual(
      failure,
    );
    expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("recovers an expired remembered session with the same scoped user and a new challenge", async () => {
    vi.mocked(getSessionCookieById).mockResolvedValue(sessionCookie);
    vi.mocked(getSession).mockResolvedValue(
      create(GetSessionResponseSchema, {
        session: create(SessionSchema, {
          id: session.id,
          factors: session.factors,
          expirationDate: { seconds: BigInt(1), nanos: 0 },
        }),
      }),
    );
    vi.mocked(createSessionAndUpdateCookie).mockResolvedValue({
      session: { ...session, id: "replacement" },
      sessionCookie: { ...sessionCookie, id: "replacement" },
      challenges: returnedChallenge,
    });
    const result = await updateOrCreateSession({
      loginName: email,
      organization,
      sessionId: cookie.id,
      challenges: challenge,
    });
    expect(result).toMatchObject({ sessionId: "replacement" });
    expect(createSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({ checks: { user: { search: { case: "userId", value: "fc-user" } } } }),
    );
    expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("never recreates an expired session during assertion verification", async () => {
    vi.mocked(getSessionCookieById).mockResolvedValue(sessionCookie);
    vi.mocked(getSession).mockResolvedValue(
      create(GetSessionResponseSchema, {
        session: create(SessionSchema, {
          id: session.id,
          factors: session.factors,
          expirationDate: { seconds: BigInt(1), nanos: 0 },
        }),
      }),
    );
    expect(await updateOrCreateSession({ loginName: email, organization, sessionId: cookie.id, checks: assertion })).toEqual(
      failure,
    );
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
    expect(setSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("rejects an update that returns a different session for the same user", async () => {
    vi.mocked(getSessionCookieById).mockResolvedValue(sessionCookie);
    vi.mocked(setSessionAndUpdateCookie).mockResolvedValue({
      ...session,
      id: "other-session",
      challenges: returnedChallenge,
    });
    expect(await updateOrCreateSession({ loginName: email, organization, sessionId: cookie.id, checks: assertion })).toEqual(
      failure,
    );
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test.each([Code.PermissionDenied, Code.NotFound, Code.Unavailable])(
    "does not replace an inaccessible session (%s)",
    async (code) => {
      vi.mocked(getSessionCookieById).mockResolvedValue(sessionCookie);
      vi.mocked(getSession).mockRejectedValue({ code });
      expect(
        await updateOrCreateSession({ loginName: email, organization, sessionId: cookie.id, challenges: challenge }),
      ).toEqual(failure);
      expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
    },
  );

  test.each([
    { result: [] },
    { error: "multiple users" },
    { result: [users.result[0], users.result[0]] },
    { result: [{ ...users.result[0], details: { resourceOwner: "foreign" } }] },
  ])("masks unknown, ambiguous, and foreign-only results: %j", async (result) => {
    vi.mocked(searchUsers).mockResolvedValue(result as never);
    expect(await updateOrCreateSession({ loginName: email, organization, challenges: challenge })).toEqual(failure);
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("passes email and phone restrictions to the shared resolver", async () => {
    vi.mocked(getLoginSettings).mockResolvedValue(
      create(LoginSettingsSchema, {
        disableLoginWithEmail: true,
        disableLoginWithPhone: true,
        ignoreUnknownUsernames: true,
      }),
    );
    vi.mocked(searchUsers).mockResolvedValue({ result: [] });
    expect(await updateOrCreateSession({ loginName: email, organization, challenges: challenge })).toEqual(failure);
    expect(searchUsers).toHaveBeenCalledWith(
      expect.objectContaining({
        loginSettings: expect.objectContaining({ disableLoginWithEmail: true, disableLoginWithPhone: true }),
      }),
    );
    expect(listUsers).not.toHaveBeenCalled();
  });

  test("rejects an email alias when email login is disabled even if the resolver returns a user", async () => {
    vi.mocked(getLoginSettings).mockResolvedValue(create(LoginSettingsSchema, { disableLoginWithEmail: true }));
    expect(await updateOrCreateSession({ loginName: email, organization, challenges: challenge })).toEqual(failure);
    expect(createSessionAndUpdateCookie).not.toHaveBeenCalled();
  });

  test("allows a canonical login name when email and phone login are disabled", async () => {
    vi.mocked(getLoginSettings).mockResolvedValue(
      create(LoginSettingsSchema, { disableLoginWithEmail: true, disableLoginWithPhone: true }),
    );
    expect(
      await updateOrCreateSession({ loginName: canonicalLoginName, organization, challenges: challenge }),
    ).toMatchObject({ sessionId: cookie.id });
  });

  test("allows a phone alias only when phone login is enabled", async () => {
    const phone = "+15555550123";
    vi.mocked(searchUsers).mockResolvedValue(
      create(ListUsersResponseSchema, {
        result: [
          {
            userId: "fc-user",
            preferredLoginName: canonicalLoginName,
            details: { resourceOwner: organization },
            type: { case: "human", value: { phone: { phone } } },
          },
        ],
      }),
    );
    expect(await updateOrCreateSession({ loginName: phone, organization, challenges: challenge })).toMatchObject({
      sessionId: cookie.id,
    });
    vi.mocked(getLoginSettings).mockResolvedValue(create(LoginSettingsSchema, { disableLoginWithPhone: true }));
    expect(await updateOrCreateSession({ loginName: phone, organization, challenges: challenge })).toEqual(failure);
  });

  test("keeps unscoped login-name lookup from introducing global email discovery", async () => {
    await updateOrCreateSession({ loginName: canonicalLoginName, challenges: challenge });
    expect(searchUsers).not.toHaveBeenCalled();
    expect(listUsers).toHaveBeenCalledWith(expect.objectContaining({ loginName: canonicalLoginName }));
  });

  test("keeps a session-only MFA request scoped to the signed cookie organization", async () => {
    vi.mocked(getSessionCookieById).mockResolvedValue(sessionCookie);
    const result = await updateOrCreateSession({ sessionId: cookie.id, checks: assertion });
    expect(result).toMatchObject({ sessionId: cookie.id });
    expect(searchUsers).toHaveBeenCalledWith(
      expect.objectContaining({ organizationId: organization, searchValue: canonicalLoginName }),
    );
  });
});

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

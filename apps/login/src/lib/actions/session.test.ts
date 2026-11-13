import {
  clearSession,
  continueWithSession,
  skipMFAAndContinueWithNextUrl,
  updateOrCreateSession,
} from "@/lib/actions/session";
import { Code } from "@connectrpc/connect";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";

vi.mock("@zitadel/client", async () => {
  const { Code } = await import("@connectrpc/connect");
  return { Code, create: vi.fn((_schema: unknown, data: unknown) => data), Duration: {} };
});

vi.mock("@/lib/server/zitadel", () => ({
  deleteSession: vi.fn(),
  getLoginSettings: vi.fn(),
  getSecuritySettings: vi.fn(),
  humanMFAInitSkipped: vi.fn(),
  listUsers: vi.fn(),
  searchUsers: vi.fn(),
}));

vi.mock("@/lib/server/session-cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  setSessionAndUpdateCookie: vi.fn(),
}));

vi.mock("@/lib/server/cookies", () => ({
  getMostRecentSessionCookie: vi.fn(),
  getSessionCookieById: vi.fn(),
  getSessionCookieByLoginName: vi.fn(),
  removeSessionFromCookie: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("@/lib/actions/flow", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

vi.mock("@/lib/server/session", () => ({
  isSessionValid: vi.fn(),
  loadSession: vi.fn(),
  hasVerifiedPrimaryFactor: vi.fn((session: any) => ({ valid: !!session?.factors?.password?.verifiedAt })),
}));

vi.mock("@/lib/server/mfa-helper", () => ({
  shouldEnforceMFA: vi.fn(() => false),
}));

vi.mock("@/lib/server/host", () => ({
  getPublicHost: vi.fn(() => "test.com"),
  getPublicHostWithProtocol: vi.fn(() => "https://test.com"),
}));

vi.mock("@/lib/actions/loginname", () => ({
  sendLoginname: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(async () => (key: string) => key),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

vi.mock("@/lib/server/grpc/interceptors/error-classification", () => ({
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
    const zitadel = await import("@/lib/server/zitadel");
    const cookies = await import("@/lib/server/cookies");
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

// Server actions are public endpoints: these tests send what a direct POST could send, including
// fields that are not part of the command types.
describe("server actions resolve the session from the caller's cookie", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
  });

  const authenticatedSession = {
    id: "session-1",
    factors: {
      user: { id: "user-1", loginName: "user@example.com", organizationId: "org-1" },
      password: { verifiedAt: { seconds: BigInt(1) } },
    },
  };

  let zitadel: any;
  let cookies: any;
  let sessionLib: any;
  let mfaHelper: any;
  let cookieHelpers: any;
  let client: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    zitadel = vi.mocked(await import("@/lib/server/zitadel"));
    cookies = vi.mocked(await import("@/lib/server/cookies"));
    sessionLib = vi.mocked(await import("@/lib/server/session"));
    mfaHelper = vi.mocked(await import("@/lib/server/mfa-helper"));
    cookieHelpers = vi.mocked(await import("@/lib/server/session-cookie"));
    client = vi.mocked(await import("@/lib/actions/flow"));

    cookies.getSessionCookieById.mockResolvedValue({ ...cookie, organization: "org-1" });
    sessionLib.loadSession.mockResolvedValue(authenticatedSession);
    zitadel.getLoginSettings.mockResolvedValue({});
    client.completeFlowOrGetUrl.mockResolvedValue({ redirect: "/signedin" });
  });

  describe("skipMFAAndContinueWithNextUrl", () => {
    test("marks the user of the caller's session, not a user id from the request", async () => {
      const res = await skipMFAAndContinueWithNextUrl({ sessionId: "session-1", userId: "victim-user" } as any);

      expect(sessionLib.loadSession).toHaveBeenCalledWith(expect.anything(), expect.objectContaining({ token: "token-1" }));
      expect(zitadel.humanMFAInitSkipped).toHaveBeenCalledWith(expect.objectContaining({ userId: "user-1" }));
      expect(res).toEqual({ redirect: "/signedin" });
    });

    test("rejects a caller without a session cookie", async () => {
      cookies.getSessionCookieById.mockResolvedValue(undefined);

      const res = await skipMFAAndContinueWithNextUrl({ sessionId: "victim-session", userId: "victim-user" } as any);

      expect(res).toEqual({ error: expect.any(String) });
      expect(zitadel.humanMFAInitSkipped).not.toHaveBeenCalled();
    });

    test("rejects a session without a verified primary factor", async () => {
      sessionLib.loadSession.mockResolvedValue({
        ...authenticatedSession,
        factors: { user: authenticatedSession.factors.user },
      });

      const res = await skipMFAAndContinueWithNextUrl({ sessionId: "session-1" });

      expect(res).toEqual({ error: expect.any(String) });
      expect(zitadel.humanMFAInitSkipped).not.toHaveBeenCalled();
    });

    test("rejects when the login policy enforces MFA", async () => {
      mfaHelper.shouldEnforceMFA.mockReturnValue(true);

      const res = await skipMFAAndContinueWithNextUrl({ sessionId: "session-1" });

      expect(res).toEqual({ error: expect.any(String) });
      expect(zitadel.humanMFAInitSkipped).not.toHaveBeenCalled();
    });
  });

  describe("continueWithSession", () => {
    test("validates the session loaded from ZITADEL, not a session object from the request", async () => {
      sessionLib.isSessionValid.mockResolvedValue(true);

      await continueWithSession({
        sessionId: "session-1",
        factors: { user: { id: "victim-user", loginName: "victim@example.com" } },
      } as any);

      expect(sessionLib.isSessionValid).toHaveBeenCalledWith(expect.objectContaining({ session: authenticatedSession }));
      expect(client.completeFlowOrGetUrl).toHaveBeenCalledWith(
        expect.objectContaining({ loginName: "user@example.com" }),
        undefined,
      );
    });

    test("rejects a session id that is not in the caller's cookie", async () => {
      cookies.getSessionCookieById.mockResolvedValue(undefined);

      const res = await continueWithSession({ sessionId: "victim-session" });

      expect(res).toEqual({ error: "couldNotContinueSession" });
      expect(sessionLib.isSessionValid).not.toHaveBeenCalled();
    });
  });

  describe("updateOrCreateSession", () => {
    test("takes the lifetime from the user's login policy and drops primary checks", async () => {
      zitadel.getLoginSettings.mockResolvedValue({ secondFactorCheckLifetime: { seconds: BigInt(300), nanos: 0 } });
      cookieHelpers.setSessionAndUpdateCookie.mockResolvedValue(authenticatedSession);

      await updateOrCreateSession({
        sessionId: "session-1",
        organization: "org-with-a-longer-lifetime",
        checks: { password: { password: "guess" }, otpEmail: { code: "123456" } },
        lifetime: { seconds: BigInt(10 * 365 * 24 * 60 * 60), nanos: 0 },
      } as any);

      expect(zitadel.getLoginSettings).toHaveBeenCalledWith(expect.objectContaining({ organization: "org-1" }));

      const [command] = cookieHelpers.setSessionAndUpdateCookie.mock.calls[0];
      expect(command.lifetime).toEqual({ seconds: BigInt(300), nanos: 0 });
      expect(command.checks).toEqual(expect.objectContaining({ otpEmail: { code: "123456" } }));
      expect(command.checks.password).toBeUndefined();
    });

    test("gives a TOTP check the second factor lifetime", async () => {
      zitadel.getLoginSettings.mockResolvedValue({
        secondFactorCheckLifetime: { seconds: BigInt(300), nanos: 0 },
        multiFactorCheckLifetime: { seconds: BigInt(600), nanos: 0 },
      });
      cookieHelpers.setSessionAndUpdateCookie.mockResolvedValue(authenticatedSession);

      await updateOrCreateSession({ sessionId: "session-1", checks: { totp: { code: "123456" } } } as any);

      const [command] = cookieHelpers.setSessionAndUpdateCookie.mock.calls[0];
      expect(command.lifetime).toEqual({ seconds: BigInt(300), nanos: 0 });
    });

    test("builds the OTP email link and the WebAuthN domain from the request host", async () => {
      vi.stubEnv("NEXT_PUBLIC_BASE_PATH", "/ui/v2/login");
      cookieHelpers.setSessionAndUpdateCookie.mockResolvedValue(authenticatedSession);

      await updateOrCreateSession({
        sessionId: "session-1",
        requestId: "oidc_123&evil=1",
        challenges: {
          otpEmail: { deliveryType: { case: "sendCode", value: { urlTemplate: "https://evil.example/?code={{.Code}}" } } },
          webAuthN: { domain: "evil.example", userVerificationRequirement: 1 },
        },
      } as any);

      const [command] = cookieHelpers.setSessionAndUpdateCookie.mock.calls[0];
      expect(command.challenges.otpEmail.deliveryType).toEqual({
        case: "sendCode",
        value: {
          urlTemplate:
            "https://test.com/ui/v2/login/otp/email?code={{.Code}}&userId={{.UserID}}&sessionId={{.SessionID}}&requestId=oidc_123%26evil%3D1",
        },
      });
      expect(command.challenges.webAuthN).toEqual({ domain: "test.com", userVerificationRequirement: 1 });
    });

    test("binds the challenges of a new session to the request host too", async () => {
      vi.stubEnv("NEXT_PUBLIC_BASE_PATH", "/ui/v2/login");
      cookies.getSessionCookieByLoginName.mockResolvedValue(undefined);
      cookieHelpers.createSessionAndUpdateCookie.mockResolvedValue({ sessionCookie: cookie });
      cookieHelpers.setSessionAndUpdateCookie.mockResolvedValue(authenticatedSession);

      await updateOrCreateSession({
        loginName: "victim@example.com",
        challenges: {
          otpEmail: { deliveryType: { case: "sendCode", value: { urlTemplate: "https://evil.example/?code={{.Code}}" } } },
        },
      } as any);

      const [createCommand] = cookieHelpers.createSessionAndUpdateCookie.mock.calls[0];
      expect(createCommand.challenges.otpEmail.deliveryType.value.urlTemplate).toMatch(
        /^https:\/\/test\.com\/ui\/v2\/login\/otp\/email\?/,
      );
    });
  });
});

describe("updateOrCreateSession without a session cookie", () => {
  const ORG = "org-1";
  const passkeyChallenge = { webAuthN: { domain: "", userVerificationRequirement: 1 } } as any;
  const user = {
    userId: "user-1",
    preferredLoginName: "jane@example.com@org.example.com",
    details: { resourceOwner: ORG },
    type: { case: "human", value: { email: { email: "jane@example.com" }, phone: { phone: "" } } },
  };

  let getLoginSettings: any;
  let searchUsers: any;
  let createSessionAndUpdateCookie: any;
  let setSessionAndUpdateCookie: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/server/zitadel");
    const cookies = await import("@/lib/server/cookies");
    const cookie = await import("@/lib/server/session-cookie");
    getLoginSettings = vi.mocked(zitadel.getLoginSettings);
    searchUsers = vi.mocked(zitadel.searchUsers);
    createSessionAndUpdateCookie = vi.mocked(cookie.createSessionAndUpdateCookie);
    setSessionAndUpdateCookie = vi.mocked(cookie.setSessionAndUpdateCookie);

    vi.mocked(cookies.getSessionCookieByLoginName).mockResolvedValue(undefined);
    getLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true });
    createSessionAndUpdateCookie.mockResolvedValue({ sessionCookie: { id: "session-1", token: "token-1" } });
    setSessionAndUpdateCookie.mockResolvedValue({
      id: "session-1",
      factors: { user: { id: "user-1" } },
      challenges: { webAuthN: { publicKeyCredentialRequestOptions: {} } },
    });
  });

  test("starts the session by user id when the typed email belongs to exactly one user of the org", async () => {
    // Given the loginname page passed on the raw email, as it does while ignoreUnknownUsernames applies
    searchUsers.mockResolvedValue({ result: [user] });

    // When the passkey page requests its challenge
    const res = await updateOrCreateSession({
      loginName: "Jane@Example.com",
      organization: ORG,
      challenges: passkeyChallenge,
    });

    // Then the session is created for the resolved user, not by the raw input as a login name
    expect(searchUsers).toHaveBeenCalledWith(
      expect.objectContaining({ searchValue: "Jane@Example.com", organizationId: ORG }),
    );
    expect(createSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({ checks: { user: { search: { case: "userId", value: "user-1" } } } }),
    );
    expect(res).toEqual(expect.objectContaining({ sessionId: "session-1", challenges: expect.anything() }));
  });

  test("keeps the login name search when no user matches, so unknown users fail as before", async () => {
    // Given nobody matches the input
    searchUsers.mockResolvedValue({ result: [] });
    createSessionAndUpdateCookie.mockRejectedValue({ code: Code.NotFound, isUserError: true });

    // When the passkey page requests its challenge
    const res = await updateOrCreateSession({
      loginName: "nobody@example.com",
      organization: ORG,
      challenges: passkeyChallenge,
    });

    // Then the session is tried by login name and the generic error comes back
    expect(createSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({ checks: { user: { search: { case: "loginName", value: "nobody@example.com" } } } }),
    );
    expect(res).toEqual({ error: "couldNotFindSession" });
  });

  test("does not resolve an email when the user's org disallows login with email", async () => {
    // Given the user's org disables email login
    searchUsers.mockResolvedValue({ result: [user] });
    getLoginSettings.mockResolvedValue({ ignoreUnknownUsernames: true, disableLoginWithEmail: true });
    createSessionAndUpdateCookie.mockRejectedValue({ code: Code.NotFound, isUserError: true });

    // When the passkey page requests its challenge with the email
    const res = await updateOrCreateSession({
      loginName: "jane@example.com",
      organization: ORG,
      challenges: passkeyChallenge,
    });

    // Then the email is not used to find the user
    expect(createSessionAndUpdateCookie).toHaveBeenCalledWith(
      expect.objectContaining({ checks: { user: { search: { case: "loginName", value: "jane@example.com" } } } }),
    );
    expect(res).toEqual({ error: "couldNotFindSession" });
  });
});

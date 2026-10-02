import { Code, ConnectError } from "@zitadel/client";
import crypto from "crypto";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { ClassifiedConnectError } from "../grpc/interceptors/error-classification";
import {
  REGISTRATION_MARKER_LIFETIME_SECONDS,
  createRegistrationMarker,
  verifyRegistrationMarker,
} from "../registration-marker";
import { signSession } from "../session-cookie-signature";
import { registerUser } from "./register";

// In-memory cookie jar shared by the next/headers mock and the fingerprint mock.
const cookieJar = new Map<string, string>();
// Attributes of the last `set` per cookie name, and the names of all cookies set, in order.
const cookieAttributes = new Map<string, Record<string, unknown>>();
const cookieSets: string[] = [];

vi.mock("next/headers", () => ({
  headers: vi.fn(async () => ({})),
  cookies: vi.fn(async () => ({
    get: (name: string) => (cookieJar.has(name) ? { name, value: cookieJar.get(name) } : undefined),
    set: (cookie: { name: string; value: string }) => {
      cookieJar.set(cookie.name, cookie.value);
      cookieAttributes.set(cookie.name, cookie);
      cookieSets.push(cookie.name);
    },
    delete: (name: string) => {
      cookieJar.delete(name);
    },
  })),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("../service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://api.example.com" } })),
}));

vi.mock("@/lib/zitadel", () => ({
  addHumanUser: vi.fn(),
  addIDPLink: vi.fn(),
  getLoginSettings: vi.fn(),
  getUserByID: vi.fn(),
  listAuthenticationMethodTypes: vi.fn(),
  listUsers: vi.fn(),
  setUserPassword: vi.fn(),
}));

vi.mock("@/lib/server/cookie", () => ({
  createSessionAndUpdateCookie: vi.fn(),
  createSessionForIdpAndUpdateCookie: vi.fn(),
}));

vi.mock("../fingerprint", () => ({
  getOrSetFingerprintId: vi.fn(async () => cookieJar.get("fingerprintId") ?? "fingerprint-1"),
  getFingerprintIdCookie: vi.fn(async () =>
    cookieJar.has("fingerprintId") ? { name: "fingerprintId", value: cookieJar.get("fingerprintId") } : undefined,
  ),
}));

vi.mock("../verify-helper", async (importOriginal) => {
  const actual = await importOriginal<typeof import("../verify-helper")>();
  return {
    ...actual,
    checkEmailVerification: vi.fn(),
    checkMFAFactors: vi.fn(),
  };
});

vi.mock("../client", () => ({
  completeFlowOrGetUrl: vi.fn(),
}));

const ORG = "org123";
const EMAIL = "ada@example.com";
const FINGERPRINT = "fingerprint-1";
const EXISTING_USER_ID = "existing-user";
const MARKER = "zitadel-registration";
const SECRET = "register-test-signing-secret-0123456789abcdef";
const OTHER_SECRET = "a-different-signing-secret-0123456789abcdef";

const SECRET_ENV_VARS = [
  "ZITADEL_SESSION_COOKIE_SECRET",
  "ZITADEL_SERVICE_USER_TOKEN",
  "SYSTEM_USER_PRIVATE_KEY",
  "SYSTEM_USER_PRIVATE_KEY_FILE",
  "ZITADEL_LOGINCLIENT_KEYFILE",
] as const;

function unsetAllSecrets() {
  for (const name of SECRET_ENV_VARS) {
    vi.stubEnv(name, undefined);
  }
}

function markerFor(userId: string, organization = ORG, nowMs = Date.now()) {
  const marker = createRegistrationMarker(userId, organization, nowMs);
  if (!marker) {
    throw new Error("test setup: no signing secret configured");
  }
  return marker;
}

function expectValidMarkerFor(userId: string) {
  expect(verifyRegistrationMarker(cookieJar.get(MARKER), { userId, organization: ORG })).toBe(true);
}

function verificationCookieFor(userId: string, fingerprint = FINGERPRINT) {
  return crypto.createHash("sha256").update(`${userId}:${fingerprint}`).digest("hex");
}

function alreadyExists() {
  return new ClassifiedConnectError(new ConnectError("Errors.User.AlreadyExists", Code.AlreadyExists));
}

function humanUser(userId: string, isVerified = false) {
  return {
    userId,
    type: { case: "human", value: { email: { email: EMAIL, isVerified } } },
  };
}

function sessionFor(userId: string) {
  return {
    session: {
      id: "session-1",
      factors: { user: { id: userId, loginName: EMAIL, organizationId: ORG } },
    },
  };
}

describe("registerUser", () => {
  let zitadel: any;
  let cookie: any;
  let client: any;
  let verifyHelper: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    cookieJar.clear();
    cookieAttributes.clear();
    cookieSets.length = 0;
    cookieJar.set("fingerprintId", FINGERPRINT);
    unsetAllSecrets();
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", SECRET);

    zitadel = vi.mocked(await import("@/lib/zitadel"));
    cookie = vi.mocked(await import("@/lib/server/cookie"));
    client = vi.mocked(await import("../client"));
    verifyHelper = vi.mocked(await import("../verify-helper"));

    zitadel.getLoginSettings.mockResolvedValue({
      allowRegister: true,
      allowLocalAuthentication: true,
      passwordCheckLifetime: { seconds: BigInt(3600) },
      defaultRedirectUri: "",
    });
    zitadel.getUserByID.mockImplementation(async ({ userId }: { userId: string }) => ({ user: humanUser(userId) }));
    zitadel.listAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [] });
    zitadel.listUsers.mockResolvedValue({ details: { totalResult: BigInt(1) }, result: [humanUser(EXISTING_USER_ID)] });
    zitadel.setUserPassword.mockResolvedValue({ details: {} });
    cookie.createSessionAndUpdateCookie.mockImplementation(async ({ checks }: any) => sessionFor(checks.user.search.value));
    verifyHelper.checkEmailVerification.mockResolvedValue(undefined);
    client.completeFlowOrGetUrl.mockResolvedValue({ redirect: "/signedin" });
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllEnvs();
  });

  const passkeyCommand = { email: EMAIL, firstName: "Ada", lastName: "Lovelace", organization: ORG, requestId: "req-1" };
  const passwordCommand = { ...passkeyCommand, password: "Secr3t!Pass" };

  describe("fresh registration", () => {
    test("passkey: creates the user, sets the verification cookie and redirects to /passkey/set", async () => {
      zitadel.addHumanUser.mockResolvedValue({ userId: "new-user" });

      const result = await registerUser(passkeyCommand);

      expect(zitadel.addHumanUser).toHaveBeenCalledWith(expect.objectContaining({ email: EMAIL, password: undefined }));
      expect(result).toEqual({
        redirect: `/passkey/set?loginName=${encodeURIComponent(EMAIL)}&organization=${ORG}&requestId=req-1`,
      });
      expect(cookieJar.get("verificationCheck")).toBe(verificationCookieFor("new-user"));
      expect(zitadel.listUsers).not.toHaveBeenCalled();
    });

    test("passkey: issues a signed registration marker for the new user, valid for 15 minutes", async () => {
      zitadel.addHumanUser.mockResolvedValue({ userId: "new-user" });

      await registerUser(passkeyCommand);

      expectValidMarkerFor("new-user");
      expect(verifyRegistrationMarker(cookieJar.get(MARKER), { userId: "new-user", organization: "other-org" })).toBe(false);
      expect(cookieAttributes.get(MARKER)).toEqual(
        expect.objectContaining({
          httpOnly: true,
          path: "/",
          sameSite: "lax",
          maxAge: REGISTRATION_MARKER_LIFETIME_SECONDS,
        }),
      );
    });

    test("passkey without a signing secret: no marker is issued", async () => {
      unsetAllSecrets();
      zitadel.addHumanUser.mockResolvedValue({ userId: "new-user" });

      const result = await registerUser(passkeyCommand);

      expect(result).toEqual({
        redirect: `/passkey/set?loginName=${encodeURIComponent(EMAIL)}&organization=${ORG}&requestId=req-1`,
      });
      expect(cookieJar.has(MARKER)).toBe(false);
    });

    test("password: creates the user with the password and completes the flow", async () => {
      zitadel.addHumanUser.mockResolvedValue({ userId: "new-user" });

      const result = await registerUser(passwordCommand);

      expect(zitadel.addHumanUser).toHaveBeenCalledWith(expect.objectContaining({ password: "Secr3t!Pass" }));
      const { checks } = cookie.createSessionAndUpdateCookie.mock.calls[0][0];
      expect(checks.user.search.value).toBe("new-user");
      expect(checks.password.password).toBe("Secr3t!Pass");
      expect(zitadel.setUserPassword).not.toHaveBeenCalled();
      expect(client.completeFlowOrGetUrl).toHaveBeenCalled();
      expect(result).toEqual({ redirect: "/signedin" });
      expect(cookieJar.has(MARKER)).toBe(false);
    });
  });

  describe("resuming an unfinished registration with a valid registration marker", () => {
    beforeEach(() => {
      zitadel.addHumanUser.mockRejectedValue(alreadyExists());
      cookieJar.set(MARKER, markerFor(EXISTING_USER_ID));
    });

    test("password: sets the password on the existing user and signs in with it", async () => {
      const result = await registerUser(passwordCommand);

      expect(zitadel.listUsers).toHaveBeenCalledWith(expect.objectContaining({ userName: EMAIL, organizationId: ORG }));
      expect(zitadel.setUserPassword).toHaveBeenCalledWith(
        expect.objectContaining({ userId: EXISTING_USER_ID, password: "Secr3t!Pass" }),
      );
      expect(zitadel.setUserPassword.mock.calls[0][0].code).toBeUndefined();
      const { checks, lifetime } = cookie.createSessionAndUpdateCookie.mock.calls[0][0];
      expect(checks.user.search.value).toBe(EXISTING_USER_ID);
      expect(checks.password.password).toBe("Secr3t!Pass");
      expect(lifetime).toEqual({ seconds: BigInt(3600) });
      expect(verifyHelper.checkEmailVerification).toHaveBeenCalled();
      expect(result).toEqual({ redirect: "/signedin" });
      expect(cookieJar.has(MARKER)).toBe(false);
    });

    test("passkey: continues to /passkey/set for the existing user without setting a password", async () => {
      const result = await registerUser(passkeyCommand);

      expect(zitadel.setUserPassword).not.toHaveBeenCalled();
      const { checks } = cookie.createSessionAndUpdateCookie.mock.calls[0][0];
      expect(checks.user.search.value).toBe(EXISTING_USER_ID);
      expect(checks.password).toBeUndefined();
      expect(result).toEqual({
        redirect: `/passkey/set?loginName=${encodeURIComponent(EMAIL)}&organization=${ORG}&requestId=req-1`,
      });
      expect(cookieJar.get("verificationCheck")).toBe(verificationCookieFor(EXISTING_USER_ID));
    });

    test("passkey: a resume keeps the original marker (no new cookie for it, original expiry)", async () => {
      const original = markerFor(EXISTING_USER_ID, ORG, Date.now() - 10 * 60 * 1000);
      cookieJar.set(MARKER, original);

      const result = await registerUser(passkeyCommand);

      expect(result).toEqual({
        redirect: `/passkey/set?loginName=${encodeURIComponent(EMAIL)}&organization=${ORG}&requestId=req-1`,
      });
      expect(cookieSets).not.toContain(MARKER);
      expect(cookieJar.get(MARKER)).toBe(original);
    });

    test("passkey: resuming does not extend the window (14 minutes ok, 16 minutes refused)", async () => {
      vi.useFakeTimers({ toFake: ["Date"] });
      const createdAt = new Date("2026-10-02T10:00:00Z").getTime();
      vi.setSystemTime(createdAt);
      zitadel.addHumanUser.mockResolvedValueOnce({ userId: EXISTING_USER_ID });
      cookieJar.delete(MARKER);

      await registerUser(passkeyCommand);
      const issued = cookieJar.get(MARKER);
      expect(issued).toBeDefined();

      vi.setSystemTime(createdAt + 14 * 60 * 1000);
      expect(await registerUser(passkeyCommand)).toEqual({
        redirect: `/passkey/set?loginName=${encodeURIComponent(EMAIL)}&organization=${ORG}&requestId=req-1`,
      });
      expect(cookieJar.get(MARKER)).toBe(issued);

      cookie.createSessionAndUpdateCookie.mockClear();
      vi.setSystemTime(createdAt + 16 * 60 * 1000);
      expect(await registerUser(passkeyCommand)).toEqual({ error: "errors.couldNotCreateUser" });
      expect(cookie.createSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test("password: setUserPassword answering with an error ends in couldNotCreateUser without a session", async () => {
      zitadel.setUserPassword.mockResolvedValue({ error: "Failed to set password" });

      const result = await registerUser(passwordCommand);

      expect(result).toEqual({ error: "errors.couldNotCreateUser" });
      expect(cookie.createSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test("password: setUserPassword throwing ends in couldNotCreateUser without a session", async () => {
      zitadel.setUserPassword.mockRejectedValue(new Error("boom"));

      const result = await registerUser(passwordCommand);

      expect(result).toEqual({ error: "errors.couldNotCreateUser" });
      expect(cookie.createSessionAndUpdateCookie).not.toHaveBeenCalled();
    });

    test("password while local authentication is not allowed: returns localAuthenticationNotAllowed, nothing written", async () => {
      zitadel.getLoginSettings.mockResolvedValue({ allowRegister: true, allowLocalAuthentication: false });

      const result = await registerUser(passwordCommand);

      expect(result).toEqual({ error: "errors.localAuthenticationNotAllowed" });
      expect(zitadel.addHumanUser).not.toHaveBeenCalled();
      expect(zitadel.setUserPassword).not.toHaveBeenCalled();
      expect(cookie.createSessionAndUpdateCookie).not.toHaveBeenCalled();
    });
  });

  describe("refusing to resume", () => {
    let validMarker: string;

    beforeEach(() => {
      zitadel.addHumanUser.mockRejectedValue(alreadyExists());
      validMarker = markerFor(EXISTING_USER_ID);
      cookieJar.set(MARKER, validMarker);
    });

    function expectNothingWritten(result: unknown) {
      expect(result).toEqual({ error: "errors.couldNotCreateUser" });
      expect(zitadel.setUserPassword).not.toHaveBeenCalled();
      expect(cookie.createSessionAndUpdateCookie).not.toHaveBeenCalled();
      expect(cookieJar.has("verificationCheck")).toBe(false);
    }

    function expectRefused(result: unknown) {
      expectNothingWritten(result);
      expect(cookieJar.get(MARKER)).toBe(validMarker);
    }

    function expectNoLookup() {
      expect(zitadel.listUsers).not.toHaveBeenCalled();
      expect(zitadel.listAuthenticationMethodTypes).not.toHaveBeenCalled();
    }

    test.each([
      ["password", passwordCommand],
      ["passkey", passkeyCommand],
    ])("%s: no registration marker", async (_name, command) => {
      cookieJar.delete(MARKER);

      expectNothingWritten(await registerUser(command));
      expectNoLookup();
      expect(cookieJar.has(MARKER)).toBe(false);
    });

    test("malformed registration marker: no lookup", async () => {
      cookieJar.set(MARKER, "not-a-marker");

      expectNothingWritten(await registerUser(passwordCommand));
      expectNoLookup();
    });

    test.each([
      ["password", passwordCommand],
      ["passkey", passkeyCommand],
    ])("%s: marker for another user id", async (_name, command) => {
      cookieJar.set(MARKER, markerFor("someone-else"));

      expectNothingWritten(await registerUser(command));
      expect(zitadel.listAuthenticationMethodTypes).not.toHaveBeenCalled();
    });

    test("marker for another organization", async () => {
      cookieJar.set(MARKER, markerFor(EXISTING_USER_ID, "other-org"));

      expectNothingWritten(await registerUser(passwordCommand));
      expectNoLookup();
    });

    test("expired marker", async () => {
      cookieJar.set(
        MARKER,
        markerFor(EXISTING_USER_ID, ORG, Date.now() - (REGISTRATION_MARKER_LIFETIME_SECONDS + 1) * 1000),
      );

      expectNothingWritten(await registerUser(passwordCommand));
      expectNoLookup();
    });

    test("marker with a tampered signature", async () => {
      const lastChar = validMarker.slice(-1);
      cookieJar.set(MARKER, validMarker.slice(0, -1) + (lastChar === "A" ? "B" : "A"));

      expectNothingWritten(await registerUser(passwordCommand));
      expectNoLookup();
    });

    test("marker with a changed expiry", async () => {
      const [userId, expiresAt, sig] = validMarker.split(".");
      cookieJar.set(MARKER, `${userId}.${Number(expiresAt) + 3600}.${sig}`);

      expectNothingWritten(await registerUser(passwordCommand));
      expectNoLookup();
    });

    test("marker signed with a different secret", async () => {
      vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", OTHER_SECRET);
      cookieJar.set(MARKER, markerFor(EXISTING_USER_ID));
      vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", SECRET);

      expectNothingWritten(await registerUser(passwordCommand));
      expectNoLookup();
    });

    test.each([
      ["password", passwordCommand],
      ["passkey", passkeyCommand],
    ])("%s: no signing secret configured", async (_name, command) => {
      unsetAllSecrets();

      expectRefused(await registerUser(command));
      expectNoLookup();
    });

    test("a session cookie for the user (user check only), without a marker", async () => {
      cookieJar.delete(MARKER);
      cookieJar.set(
        "sessions",
        JSON.stringify([
          signSession({
            id: "session-from-loginname",
            token: "token",
            loginName: EMAIL,
            organization: ORG,
            creationTs: `${Date.now()}`,
            expirationTs: `${Date.now() + 60_000}`,
            changeTs: `${Date.now()}`,
          }),
        ]),
      );

      expectNothingWritten(await registerUser(passwordCommand));
    });

    test.each([
      ["password", passwordCommand],
      ["passkey", passkeyCommand],
    ])("%s: the user already has an authentication method", async (_name, command) => {
      zitadel.listAuthenticationMethodTypes.mockResolvedValue({ authMethodTypes: [1] });

      expectRefused(await registerUser(command));
    });

    test("the authentication method lookup fails", async () => {
      zitadel.listAuthenticationMethodTypes.mockRejectedValue(new Error("unavailable"));

      expectRefused(await registerUser(passwordCommand));
    });

    test("the e-mail address is verified", async () => {
      zitadel.listUsers.mockResolvedValue({ details: {}, result: [humanUser(EXISTING_USER_ID, true)] });

      expectRefused(await registerUser(passwordCommand));
    });

    test("the user is not a human user", async () => {
      zitadel.listUsers.mockResolvedValue({
        details: {},
        result: [{ userId: EXISTING_USER_ID, type: { case: "machine", value: {} } }],
      });

      expectRefused(await registerUser(passwordCommand));
    });

    test("no user is found", async () => {
      zitadel.listUsers.mockResolvedValue({ details: {}, result: [] });

      expectRefused(await registerUser(passwordCommand));
    });

    test("two users are found", async () => {
      zitadel.listUsers.mockResolvedValue({
        details: {},
        result: [humanUser(EXISTING_USER_ID), humanUser("other-user")],
      });

      expectRefused(await registerUser(passwordCommand));
    });

    test("the user lookup fails", async () => {
      zitadel.listUsers.mockRejectedValue(new Error("unavailable"));

      expectRefused(await registerUser(passwordCommand));
    });

    test.each([
      ["a plain error", () => new Error("network")],
      ["a different gRPC error", () => new ClassifiedConnectError(new ConnectError("invalid", Code.InvalidArgument))],
    ])("addHumanUser fails with %s: no lookup at all", async (_name, makeError) => {
      zitadel.addHumanUser.mockRejectedValue(makeError());

      expectRefused(await registerUser(passwordCommand));
      expect(zitadel.listUsers).not.toHaveBeenCalled();
      expect(zitadel.listAuthenticationMethodTypes).not.toHaveBeenCalled();
    });
  });
});

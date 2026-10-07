import { createHmac, hkdfSync } from "crypto";
import { mkdtempSync, writeFileSync } from "fs";
import { tmpdir } from "os";
import { join } from "path";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import {
  MIN_SESSION_COOKIE_SECRET_LENGTH,
  getSessionCookieSecretConfigError,
  getSessionCookieSecretStartupNotice,
  getSessionCookieSecrets,
  hasSessionCookieSecret,
  hmacWithDerivedKey,
  isUsingCredentialFallback,
  parseAndVerifySessions,
  signSession,
  signaturesEqual,
  stripSessionSignature,
  verifySession,
} from "./session-cookie-signature";

const session = {
  id: "session-1",
  token: "token-1",
  loginName: "user@example.com",
  organization: "org-1",
  creationTs: "1700000000000",
  expirationTs: "1800000000000",
  changeTs: "1700000000000",
};

const OLD_SECRET = "old-secret-".padEnd(MIN_SESSION_COOKIE_SECRET_LENGTH, "o");
const NEW_SECRET = "new-secret-".padEnd(MIN_SESSION_COOKIE_SECRET_LENGTH, "n");

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

describe("session-cookie-signature", () => {
  beforeEach(() => {
    unsetAllSecrets();
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  test("signs and verifies a session", () => {
    const signed = signSession(session);

    expect(signed.sig).toEqual(expect.any(String));
    expect(verifySession(signed)).toBe(true);
    expect(stripSessionSignature(signed)).toEqual(session);
  });

  test("signature is HMAC-SHA256 over the canonical entry with the HKDF-derived session cookie key", () => {
    const key = Buffer.from(
      hkdfSync("sha256", "test-session-cookie-secret-at-least-32-chars", "", "zitadel-login-session-cookie-v1", 32),
    );
    const canonical = JSON.stringify([
      "v1",
      Object.fromEntries(Object.entries(session).sort(([a], [b]) => (a < b ? -1 : 1))),
    ]);

    expect(signSession(session).sig).toBe(createHmac("sha256", key).update(canonical).digest("base64url"));
  });

  test("hmacWithDerivedKey separates keys by info", () => {
    const secret = "test-session-cookie-secret-at-least-32-chars";

    expect(hmacWithDerivedKey(secret, "info-a", "payload")).toBe(hmacWithDerivedKey(secret, "info-a", "payload"));
    expect(hmacWithDerivedKey(secret, "info-a", "payload")).not.toBe(hmacWithDerivedKey(secret, "info-b", "payload"));
    expect(hmacWithDerivedKey(secret, "info-a", "payload")).toBe(
      createHmac("sha256", Buffer.from(hkdfSync("sha256", secret, "", "info-a", 32)))
        .update("payload")
        .digest("base64url"),
    );
  });

  test("signaturesEqual compares exactly", () => {
    expect(signaturesEqual("abc", "abc")).toBe(true);
    expect(signaturesEqual("abc", "abd")).toBe(false);
    expect(signaturesEqual("abc", "abcd")).toBe(false);
  });

  test("rejects a missing signature", () => {
    expect(verifySession(session)).toBe(false);
  });

  test("rejects entries with an empty id or token even if correctly signed", () => {
    expect(() => signSession({ ...session, token: "" })).toThrow(/non-empty id and token/);
    expect(() => signSession({ ...session, id: "" })).toThrow(/non-empty id and token/);

    // a signature over an empty token must never verify, regardless of how it was produced
    const signed = signSession(session);
    expect(verifySession({ ...signed, token: "" })).toBe(false);
  });

  test("does not confuse field boundaries in the signed payload", () => {
    const signed = signSession({ ...session, id: "ab", token: "c" });

    expect(verifySession({ ...signed, id: "a", token: "bc" })).toBe(false);
  });

  test("rejects a forged session id with an otherwise valid signature", () => {
    const signed = signSession(session);

    expect(verifySession({ ...signed, id: "guessed-victim-session" })).toBe(false);
  });

  test("covers all fields of the entry, not only id and token", () => {
    const signed = signSession(session);

    expect(verifySession({ ...signed, loginName: "other@example.com" })).toBe(false);
    expect(verifySession({ ...signed, organization: "other-org" })).toBe(false);
    expect(verifySession({ ...signed, changeTs: "9999999999999" })).toBe(false);
    expect(verifySession({ ...signed, expirationTs: "9999999999999" })).toBe(false);
    expect(verifySession({ ...signed, requestId: "injected-request" })).toBe(false);
  });

  test("is independent of key order and of undefined optional fields", () => {
    const signed = signSession({ ...session, organization: undefined, requestId: undefined });

    // simulate the JSON round-trip through the cookie: undefined keys disappear, key order may change
    const roundTripped = JSON.parse(JSON.stringify(signed));
    const reordered = Object.fromEntries(Object.entries(roundTripped).reverse());

    expect(roundTripped).not.toHaveProperty("organization");
    expect(verifySession(roundTripped)).toBe(true);
    expect(verifySession(reordered)).toBe(true);
  });

  test("rejects a guessed id with an empty token and no signature", () => {
    expect(
      verifySession({
        ...session,
        id: "guessed-victim-session",
        token: "",
      }),
    ).toBe(false);
  });

  test("rejects a signature created with a different secret", () => {
    const signed = signSession(session);
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "other-secret-".padEnd(MIN_SESSION_COOKIE_SECRET_LENGTH, "x"));

    expect(verifySession(signed)).toBe(false);
  });

  test("drops unsigned and tampered entries when parsing a cookie", () => {
    const valid = signSession(session);
    const forged = { ...session, id: "stolen-or-guessed-id", token: "", sig: valid.sig };

    expect(parseAndVerifySessions(JSON.stringify([forged, valid, session]))).toEqual([session]);
  });

  test("returns an empty list for malformed cookie values", () => {
    expect(parseAndVerifySessions("not-json")).toEqual([]);
    expect(parseAndVerifySessions(JSON.stringify({ id: "session-1" }))).toEqual([]);
    expect(parseAndVerifySessions(undefined)).toEqual([]);
  });

  test("supports secret rotation: signs with the first secret, verifies with all", () => {
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", OLD_SECRET);
    const signedWithOld = signSession(session);

    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", `${NEW_SECRET}, ${OLD_SECRET}`);
    expect(verifySession(signedWithOld)).toBe(true);
    const signedWithNew = signSession(session);
    expect(signedWithNew.sig).not.toEqual(signedWithOld.sig);

    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", NEW_SECRET);
    expect(verifySession(signedWithNew)).toBe(true);
    expect(verifySession(signedWithOld)).toBe(false);
  });

  test("falls back to ZITADEL_SERVICE_USER_TOKEN when no dedicated secret is set", () => {
    unsetAllSecrets();
    vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");

    expect(hasSessionCookieSecret()).toBe(true);
    expect(verifySession(signSession(session))).toBe(true);
  });

  test("falls back to SYSTEM_USER_PRIVATE_KEY when no other secret is set", () => {
    unsetAllSecrets();
    vi.stubEnv("SYSTEM_USER_PRIVATE_KEY", "system-user-private-key");

    expect(verifySession(signSession(session))).toBe(true);
  });

  test("falls back to SYSTEM_USER_PRIVATE_KEY_FILE and ZITADEL_LOGINCLIENT_KEYFILE", () => {
    const dir = mkdtempSync(join(tmpdir(), "session-cookie-signature-"));
    const systemKeyFile = join(dir, "system.key");
    const loginClientKeyFile = join(dir, "login-client.key");
    writeFileSync(systemKeyFile, "system-user-private-key-from-file\n");
    writeFileSync(loginClientKeyFile, "login-client-key-from-file\n");

    unsetAllSecrets();
    vi.stubEnv("SYSTEM_USER_PRIVATE_KEY_FILE", systemKeyFile);
    expect(getSessionCookieSecrets()).toEqual(["system-user-private-key-from-file"]);
    expect(verifySession(signSession(session))).toBe(true);

    unsetAllSecrets();
    vi.stubEnv("ZITADEL_LOGINCLIENT_KEYFILE", loginClientKeyFile);
    expect(getSessionCookieSecrets()).toEqual(["login-client-key-from-file"]);
    expect(verifySession(signSession(session))).toBe(true);
  });

  test("ignores an unreadable key file", () => {
    unsetAllSecrets();
    vi.stubEnv("ZITADEL_LOGINCLIENT_KEYFILE", "/does/not/exist.key");

    expect(hasSessionCookieSecret()).toBe(false);
    expect(() => signSession(session)).toThrow(/signing secret is not configured/);
  });

  test("keeps cookies valid while migrating between credential types", () => {
    unsetAllSecrets();
    vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");
    const signedWithPat = signSession(session);

    // both credentials configured during the migration
    vi.stubEnv("SYSTEM_USER_PRIVATE_KEY", "system-user-private-key");
    expect(verifySession(signedWithPat)).toBe(true);
  });

  test("a dedicated secret takes precedence over API credentials", () => {
    vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");

    expect(getSessionCookieSecrets()[0]).toBe("test-session-cookie-secret-at-least-32-chars");
  });

  test("does not verify when no signing secret is configured", () => {
    const signed = signSession(session);
    unsetAllSecrets();

    expect(hasSessionCookieSecret()).toBe(false);
    expect(verifySession(signed)).toBe(false);
    expect(() => signSession(session)).toThrow(/signing secret is not configured/);
  });

  test("rejects a dedicated secret shorter than the minimum length instead of falling back", () => {
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "too-short");
    vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");

    expect(getSessionCookieSecretConfigError()).toMatch(/at least 32 characters/);
    expect(getSessionCookieSecrets()).toEqual([]);
    expect(hasSessionCookieSecret()).toBe(false);
    expect(() => signSession(session)).toThrow(/at least 32 characters/);
  });

  test("rejects the rotation list if any value is shorter than the minimum length", () => {
    const signed = signSession(session);
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", `${NEW_SECRET},short`);

    expect(getSessionCookieSecretConfigError()).toMatch(/1 value\(s\) are shorter/);
    expect(verifySession(signed)).toBe(false);
  });

  test("accepts a dedicated secret of exactly the minimum length", () => {
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "x".repeat(MIN_SESSION_COOKIE_SECRET_LENGTH));

    expect(getSessionCookieSecretConfigError()).toBeUndefined();
    expect(verifySession(signSession(session))).toBe(true);
  });

  test("does not apply the minimum length to API credentials", () => {
    unsetAllSecrets();
    vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "short-token");

    expect(getSessionCookieSecretConfigError()).toBeUndefined();
    expect(verifySession(signSession(session))).toBe(true);
  });

  describe("startup notice", () => {
    test("is empty when a dedicated secret is configured", () => {
      vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");

      expect(isUsingCredentialFallback()).toBe(false);
      expect(getSessionCookieSecretStartupNotice()).toBeUndefined();
    });

    test("warns that the credential fallback is deprecated", () => {
      unsetAllSecrets();
      vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");

      expect(isUsingCredentialFallback()).toBe(true);
      const notice = getSessionCookieSecretStartupNotice();
      expect(notice?.level).toBe("warn");
      expect(notice?.message).toMatch(/Rotating that credential signs all users out/);
      expect(notice?.message).toMatch(/deprecated/);
    });

    test("reports an error for a too short dedicated secret", () => {
      vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "too-short");
      vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");

      expect(getSessionCookieSecretStartupNotice()).toEqual({
        level: "error",
        message: expect.stringMatching(/at least 32 characters/),
      });
    });

    test("reports an error when no secret is available at all", () => {
      unsetAllSecrets();

      expect(getSessionCookieSecretStartupNotice()?.level).toBe("error");
    });
  });
});

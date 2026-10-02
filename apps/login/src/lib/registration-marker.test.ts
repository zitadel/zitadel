import { createHmac, hkdfSync } from "crypto";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import {
  REGISTRATION_MARKER_LIFETIME_SECONDS,
  createRegistrationMarker,
  readRegistrationMarker,
  verifyRegistrationMarker,
} from "./registration-marker";

const SECRET = "registration-marker-test-secret-0123456789";
const OTHER_SECRET = "another-registration-marker-secret-9876543210";
const USER_ID = "312345678901234567";
const ORG = "298765432109876543";
const NOW_MS = 1_800_000_000_000;

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

function hmac(secret: string, info: string, payload: string) {
  const key = Buffer.from(hkdfSync("sha256", secret, "", info, 32));
  return createHmac("sha256", key).update(payload).digest("base64url");
}

describe("registration-marker", () => {
  beforeEach(() => {
    unsetAllSecrets();
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", SECRET);
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  test("creates a marker <userId>.<expiresAt>.<signature> valid for 15 minutes", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);
    const expiresAt = NOW_MS / 1000 + REGISTRATION_MARKER_LIFETIME_SECONDS;

    expect(REGISTRATION_MARKER_LIFETIME_SECONDS).toBe(15 * 60);
    expect(marker).toBe(
      `${USER_ID}.${expiresAt}.${hmac(SECRET, "zitadel-login-registration-v1", `registration:${USER_ID}:${expiresAt}:${ORG}`)}`,
    );
    expect(verifyRegistrationMarker(marker, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(true);
  });

  test("is bound to the user id", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);

    expect(verifyRegistrationMarker(marker, { userId: "other-user", organization: ORG }, NOW_MS)).toBe(false);
  });

  test("is bound to the organization", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);

    expect(verifyRegistrationMarker(marker, { userId: USER_ID, organization: "other-org" }, NOW_MS)).toBe(false);
  });

  test("expires after 15 minutes", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);
    const lifetimeMs = REGISTRATION_MARKER_LIFETIME_SECONDS * 1000;

    expect(verifyRegistrationMarker(marker, { userId: USER_ID, organization: ORG }, NOW_MS + lifetimeMs - 1000)).toBe(true);
    expect(verifyRegistrationMarker(marker, { userId: USER_ID, organization: ORG }, NOW_MS + lifetimeMs)).toBe(false);
  });

  test("rejects a changed expiry", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS)!;
    const [userId, expiresAt, sig] = marker.split(".");

    expect(
      verifyRegistrationMarker(
        `${userId}.${Number(expiresAt) + 3600}.${sig}`,
        { userId: USER_ID, organization: ORG },
        NOW_MS,
      ),
    ).toBe(false);
  });

  test("rejects a tampered signature", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS)!;
    const lastChar = marker.slice(-1);
    const tampered = marker.slice(0, -1) + (lastChar === "A" ? "B" : "A");

    expect(verifyRegistrationMarker(tampered, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(false);
  });

  test("rejects a marker signed with a different secret", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", OTHER_SECRET);

    expect(verifyRegistrationMarker(marker, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(false);
  });

  test("accepts a marker signed with any configured secret (rotation)", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", `${OTHER_SECRET},${SECRET}`);

    expect(verifyRegistrationMarker(marker, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(true);
  });

  test("uses a key separate from the session cookie signature key", () => {
    const expiresAt = NOW_MS / 1000 + REGISTRATION_MARKER_LIFETIME_SECONDS;
    const payload = `registration:${USER_ID}:${expiresAt}:${ORG}`;
    const withSessionKey = `${USER_ID}.${expiresAt}.${hmac(SECRET, "zitadel-login-session-cookie-v1", payload)}`;
    const withoutDerivation = `${USER_ID}.${expiresAt}.${createHmac("sha256", SECRET).update(payload).digest("base64url")}`;

    expect(verifyRegistrationMarker(withSessionKey, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(false);
    expect(verifyRegistrationMarker(withoutDerivation, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(false);
  });

  test("works with the API credential fallback like the session cookie signature", () => {
    unsetAllSecrets();
    vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");

    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);

    expect(marker).toBeDefined();
    expect(verifyRegistrationMarker(marker, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(true);
  });

  test("without a signing secret no marker is created and none verifies", () => {
    const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);
    unsetAllSecrets();

    expect(createRegistrationMarker(USER_ID, ORG, NOW_MS)).toBeUndefined();
    expect(verifyRegistrationMarker(marker, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(false);
  });

  test("with an invalid dedicated secret no marker is created (no fallback to the credential)", () => {
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "too-short");
    vi.stubEnv("ZITADEL_SERVICE_USER_TOKEN", "service-user-token");

    expect(createRegistrationMarker(USER_ID, ORG, NOW_MS)).toBeUndefined();
  });

  test("does not create a marker for an empty or ambiguous user id", () => {
    expect(createRegistrationMarker("", ORG, NOW_MS)).toBeUndefined();
    expect(createRegistrationMarker("a.b", ORG, NOW_MS)).toBeUndefined();
    expect(createRegistrationMarker("a:b", ORG, NOW_MS)).toBeUndefined();
  });

  test.each([
    ["undefined", undefined],
    ["empty", ""],
    ["two parts", `${USER_ID}.123`],
    ["four parts", `${USER_ID}.123.sig.extra`],
    ["non-numeric expiry", `${USER_ID}.abc.sig`],
    ["empty signature", `${USER_ID}.${NOW_MS / 1000 + 60}.`],
  ])("rejects a malformed value (%s)", (_name, value) => {
    expect(verifyRegistrationMarker(value, { userId: USER_ID, organization: ORG }, NOW_MS)).toBe(false);
    expect(readRegistrationMarker(value, ORG, NOW_MS)).toBeUndefined();
  });

  describe("readRegistrationMarker", () => {
    test("returns the user id of a valid marker for the organization", () => {
      const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);

      expect(readRegistrationMarker(marker, ORG, NOW_MS)).toBe(USER_ID);
    });

    test("returns undefined for another organization, an expired marker, a tampered signature or another secret", () => {
      const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS)!;

      expect(readRegistrationMarker(marker, "other-org", NOW_MS)).toBeUndefined();
      expect(readRegistrationMarker(marker, ORG, NOW_MS + REGISTRATION_MARKER_LIFETIME_SECONDS * 1000)).toBeUndefined();
      expect(readRegistrationMarker(marker.slice(0, -1) + (marker.endsWith("A") ? "B" : "A"), ORG, NOW_MS)).toBeUndefined();

      vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", OTHER_SECRET);
      expect(readRegistrationMarker(marker, ORG, NOW_MS)).toBeUndefined();
    });

    test("returns undefined without a signing secret", () => {
      const marker = createRegistrationMarker(USER_ID, ORG, NOW_MS);
      unsetAllSecrets();

      expect(readRegistrationMarker(marker, ORG, NOW_MS)).toBeUndefined();
    });
  });
});

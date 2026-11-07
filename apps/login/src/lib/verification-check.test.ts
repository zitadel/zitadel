import { createHash } from "crypto";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const mockCookies = { get: vi.fn(), set: vi.fn() };

vi.mock("next/headers", () => ({
  cookies: vi.fn(async () => mockCookies),
}));

vi.mock("./fingerprint", () => ({
  getFingerprintIdCookie: vi.fn(),
  getOrSetFingerprintId: vi.fn(),
}));

import { getFingerprintIdCookie, getOrSetFingerprintId } from "./fingerprint";
import { MIN_SESSION_COOKIE_SECRET_LENGTH, signSession } from "./session-cookie-signature";
import {
  createVerificationCheck,
  hasValidUserVerificationCheckCookie,
  isVerificationCheckValid,
  setUserVerificationCheckCookie,
  VERIFICATION_CHECK_MAX_AGE_SECONDS,
} from "./verification-check";

const ENV_KEYS = [
  "ZITADEL_SESSION_COOKIE_SECRET",
  "ZITADEL_SERVICE_USER_TOKEN",
  "SYSTEM_USER_PRIVATE_KEY",
  "SYSTEM_USER_PRIVATE_KEY_FILE",
  "ZITADEL_LOGINCLIENT_KEYFILE",
] as const;

describe("verification-check", () => {
  const originalEnv: Record<string, string | undefined> = {};

  beforeEach(() => {
    vi.clearAllMocks();
    for (const key of ENV_KEYS) {
      originalEnv[key] = process.env[key];
      delete process.env[key];
    }
    process.env.ZITADEL_SERVICE_USER_TOKEN = "service-user-token";
  });

  afterEach(() => {
    for (const key of ENV_KEYS) {
      if (originalEnv[key] === undefined) {
        delete process.env[key];
      } else {
        process.env[key] = originalEnv[key];
      }
    }
  });

  it("accepts a check issued for the same user and fingerprint", () => {
    const value = createVerificationCheck("user-1", "fp-1");
    expect(isVerificationCheckValid(value, "user-1", "fp-1")).toBe(true);
  });

  it("rejects a check for another user or another fingerprint", () => {
    const value = createVerificationCheck("user-1", "fp-1");
    expect(isVerificationCheckValid(value, "user-2", "fp-1")).toBe(false);
    expect(isVerificationCheckValid(value, "user-1", "fp-2")).toBe(false);
  });

  it("rejects the legacy, client computable sha256(userId:fingerprintId) value", () => {
    const forged = createHash("sha256").update("victim:attacker-fp").digest("hex");
    expect(isVerificationCheckValid(forged, "victim", "attacker-fp")).toBe(false);
  });

  it("rejects a value signed with a secret the server does not know", () => {
    const value = createVerificationCheck("user-1", "fp-1");
    process.env.ZITADEL_SERVICE_USER_TOKEN = "another-token";
    expect(isVerificationCheckValid(value, "user-1", "fp-1")).toBe(false);
  });

  it("enforces the expiry server side and rejects a tampered expiry", () => {
    const now = Date.now();
    const value = createVerificationCheck("user-1", "fp-1", now);

    expect(isVerificationCheckValid(value, "user-1", "fp-1", now + (VERIFICATION_CHECK_MAX_AGE_SECONDS - 1) * 1000)).toBe(
      true,
    );
    expect(isVerificationCheckValid(value, "user-1", "fp-1", now + (VERIFICATION_CHECK_MAX_AGE_SECONDS + 1) * 1000)).toBe(
      false,
    );

    const [version, expiresAt, signature] = value.split(".");
    const extended = `${version}.${Number(expiresAt) + 3600}.${signature}`;
    expect(isVerificationCheckValid(extended, "user-1", "fp-1", now)).toBe(false);
  });

  it("rejects malformed values", () => {
    for (const value of [undefined, "", "v1", "v1..", "v2.123.abc", "v1.abc.def", "a.b.c.d"]) {
      expect(isVerificationCheckValid(value, "user-1", "fp-1")).toBe(false);
    }
  });

  it("does not let field boundaries be shifted between userId and fingerprintId", () => {
    const value = createVerificationCheck("user:1", "fp");
    expect(isVerificationCheckValid(value, "user", "1:fp")).toBe(false);
  });

  it("prefers the dedicated secret and supports rotation", () => {
    const oldSecret = "old-secret-".padEnd(MIN_SESSION_COOKIE_SECRET_LENGTH, "o");
    const newSecret = "new-secret-".padEnd(MIN_SESSION_COOKIE_SECRET_LENGTH, "n");
    process.env.ZITADEL_SESSION_COOKIE_SECRET = oldSecret;
    const value = createVerificationCheck("user-1", "fp-1");

    process.env.ZITADEL_SESSION_COOKIE_SECRET = `${newSecret}, ${oldSecret}`;
    expect(isVerificationCheckValid(value, "user-1", "fp-1")).toBe(true);
  });

  it("fails closed without any secret", () => {
    delete (process.env as Record<string, string | undefined>).ZITADEL_SERVICE_USER_TOKEN;
    expect(() => createVerificationCheck("user-1", "fp-1")).toThrow();
    expect(isVerificationCheckValid("v1.9999999999.abc", "user-1", "fp-1")).toBe(false);
  });

  it("uses a key that is independent of the session cookie signature", () => {
    const session = signSession({ id: "user-1", token: "fp-1" });
    const [, , signature] = createVerificationCheck("user-1", "fp-1").split(".");
    expect(signature).not.toBe(session.sig);
  });

  it("round trips through the cookie helpers", async () => {
    vi.mocked(getOrSetFingerprintId).mockResolvedValue("fp-1");
    await setUserVerificationCheckCookie("user-1");

    expect(mockCookies.set).toHaveBeenCalledWith(
      expect.objectContaining({ name: "verificationCheck", httpOnly: true, maxAge: VERIFICATION_CHECK_MAX_AGE_SECONDS }),
    );

    const { value } = mockCookies.set.mock.calls[0][0];
    mockCookies.get.mockReturnValue({ value });
    vi.mocked(getFingerprintIdCookie).mockResolvedValue({ name: "fingerprintId", value: "fp-1" } as any);

    expect(await hasValidUserVerificationCheckCookie("user-1")).toBe(true);
    expect(await hasValidUserVerificationCheckCookie("user-2")).toBe(false);
  });
});

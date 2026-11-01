import { createHash } from "crypto";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import {
  checkDeviceCodeCheckCookie,
  createDeviceCodeCheck,
  DEVICE_CODE_CHECK_MAX_AGE_SECONDS,
  isDeviceCodeCheckValid,
  setDeviceCodeCheckCookie,
} from "./device-code-check";
import { createVerificationCheck } from "./verification-check";

const { jar } = vi.hoisted(() => ({ jar: new Map<string, string>() }));

vi.mock("next/headers", () => ({
  cookies: vi.fn(async () => ({
    get: (name: string) => (jar.has(name) ? { name, value: jar.get(name) } : undefined),
    set: ({ name, value }: { name: string; value: string }) => jar.set(name, value),
    delete: (name: string) => jar.delete(name),
  })),
}));

describe("deviceCodeCheck cookie", () => {
  beforeEach(() => {
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");
    jar.clear();
    jar.set("fingerprintId", "fingerprint-1");
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  test("passes for the user code and browser that set it", async () => {
    await setDeviceCodeCheckCookie("ABCD-1234");

    expect(await checkDeviceCodeCheckCookie("ABCD-1234")).toBe(true);
    expect(await checkDeviceCodeCheckCookie("WXYZ-9876")).toBe(false);
  });

  test("fails when the fingerprint changed", async () => {
    await setDeviceCodeCheckCookie("ABCD-1234");
    jar.set("fingerprintId", "fingerprint-2");

    expect(await checkDeviceCodeCheckCookie("ABCD-1234")).toBe(false);
  });

  test("fails without a fingerprint cookie", async () => {
    await setDeviceCodeCheckCookie("ABCD-1234");
    jar.delete("fingerprintId");

    expect(await checkDeviceCodeCheckCookie("ABCD-1234")).toBe(false);
  });

  test("rejects the previous unsigned sha256(userCode:fingerprintId) value", async () => {
    jar.set("deviceCodeCheck", createHash("sha256").update("ABCD-1234:fingerprint-1").digest("hex"));

    expect(await checkDeviceCodeCheckCookie("ABCD-1234")).toBe(false);
  });

  test("rejects a value signed with a secret the server does not know", () => {
    const value = createDeviceCodeCheck("ABCD-1234", "fingerprint-1");
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "another-secret-".padEnd(32, "x"));

    expect(isDeviceCodeCheckValid(value, "ABCD-1234", "fingerprint-1")).toBe(false);
  });

  test("enforces the expiry server side", () => {
    const now = Date.now();
    const value = createDeviceCodeCheck("ABCD-1234", "fingerprint-1", now);

    expect(
      isDeviceCodeCheckValid(value, "ABCD-1234", "fingerprint-1", now + (DEVICE_CODE_CHECK_MAX_AGE_SECONDS - 1) * 1000),
    ).toBe(true);
    expect(
      isDeviceCodeCheckValid(value, "ABCD-1234", "fingerprint-1", now + (DEVICE_CODE_CHECK_MAX_AGE_SECONDS + 1) * 1000),
    ).toBe(false);
  });

  test("is not interchangeable with a verification check for the same values", () => {
    const now = Date.now();
    // both proofs live 5 minutes, so the same fields and time give the same payload; only the purpose key differs
    const verificationCheck = createVerificationCheck("ABCD-1234", "fingerprint-1", now);

    expect(isDeviceCodeCheckValid(verificationCheck, "ABCD-1234", "fingerprint-1", now)).toBe(false);
  });
});

import {
  checkDeviceCodeCheckCookie,
  createDeviceCodeCheck,
  DEVICE_CODE_CHECK_MAX_AGE_SECONDS,
  isDeviceCodeCheckValid,
  setDeviceCodeCheckCookie,
} from "@/lib/server/device-code-check";
import { createVerificationCheck } from "@/lib/server/verification-check";
import { getSecuritySettings } from "@/lib/server/zitadel";
import { createHash } from "crypto";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";

const { jar, setCookie } = vi.hoisted(() => {
  const jar = new Map<string, string>();
  return { jar, setCookie: vi.fn(({ name, value }: { name: string; value: string }) => jar.set(name, value)) };
});

vi.mock("next/headers", () => ({
  cookies: vi.fn(async () => ({
    get: (name: string) => (jar.has(name) ? { name, value: jar.get(name) } : undefined),
    set: setCookie,
    delete: (name: string) => jar.delete(name),
  })),
  headers: vi.fn(async () => new Headers()),
}));
vi.mock("./service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://api.example.com" } })),
}));
vi.mock("./zitadel", () => ({
  getSecuritySettings: vi.fn(),
}));

describe("deviceCodeCheck cookie", () => {
  beforeEach(() => {
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "test-session-cookie-secret-at-least-32-chars");
    jar.clear();
    jar.set("bindingId", "binding-id-1");
  });

  afterEach(() => {
    vi.unstubAllEnvs();
  });

  test("passes for the user code and browser that set it", async () => {
    await setDeviceCodeCheckCookie("ABCD-1234");

    expect(await checkDeviceCodeCheckCookie("ABCD-1234")).toBe(true);
    expect(await checkDeviceCodeCheckCookie("WXYZ-9876")).toBe(false);
  });

  test("sets the cookie with SameSite and Secure so it also works in an iframe", async () => {
    vi.mocked(getSecuritySettings).mockResolvedValueOnce({ embeddedIframe: { enabled: true } } as any);

    await setDeviceCodeCheckCookie("ABCD-1234");

    expect(setCookie).toHaveBeenCalledWith(
      expect.objectContaining({
        name: "deviceCodeCheck",
        httpOnly: true,
        maxAge: DEVICE_CODE_CHECK_MAX_AGE_SECONDS,
        sameSite: "none",
        secure: true,
      }),
    );
  });

  test("fails when the binding id changed", async () => {
    await setDeviceCodeCheckCookie("ABCD-1234");
    jar.set("bindingId", "binding-id-2");

    expect(await checkDeviceCodeCheckCookie("ABCD-1234")).toBe(false);
  });

  test("fails without a binding id cookie", async () => {
    await setDeviceCodeCheckCookie("ABCD-1234");
    jar.delete("bindingId");

    expect(await checkDeviceCodeCheckCookie("ABCD-1234")).toBe(false);
  });

  test("rejects the previous unsigned sha256(userCode:fingerprintId) value", async () => {
    jar.set("deviceCodeCheck", createHash("sha256").update("ABCD-1234:binding-id-1").digest("hex"));

    expect(await checkDeviceCodeCheckCookie("ABCD-1234")).toBe(false);
  });

  test("rejects a value signed with a secret the server does not know", () => {
    const value = createDeviceCodeCheck("ABCD-1234", "binding-id-1");
    vi.stubEnv("ZITADEL_SESSION_COOKIE_SECRET", "another-secret-".padEnd(32, "x"));

    expect(isDeviceCodeCheckValid(value, "ABCD-1234", "binding-id-1")).toBe(false);
  });

  test("enforces the expiry server side", () => {
    const now = Date.now();
    const value = createDeviceCodeCheck("ABCD-1234", "binding-id-1", now);

    expect(
      isDeviceCodeCheckValid(value, "ABCD-1234", "binding-id-1", now + (DEVICE_CODE_CHECK_MAX_AGE_SECONDS - 1) * 1000),
    ).toBe(true);
    expect(
      isDeviceCodeCheckValid(value, "ABCD-1234", "binding-id-1", now + (DEVICE_CODE_CHECK_MAX_AGE_SECONDS + 1) * 1000),
    ).toBe(false);
  });

  test("is not interchangeable with a verification check for the same values", () => {
    const now = Date.now();
    // both proofs live 5 minutes, so the same fields and time give the same payload; only the purpose key differs
    const verificationCheck = createVerificationCheck("ABCD-1234", "binding-id-1", now);

    expect(isDeviceCodeCheckValid(verificationCheck, "ABCD-1234", "binding-id-1", now)).toBe(false);
  });
});

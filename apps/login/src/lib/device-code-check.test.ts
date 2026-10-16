import { beforeEach, describe, expect, test, vi } from "vitest";
import { checkDeviceCodeCheckCookie, setDeviceCodeCheckCookie } from "./device-code-check";

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
    jar.clear();
    jar.set("fingerprintId", "fingerprint-1");
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
});

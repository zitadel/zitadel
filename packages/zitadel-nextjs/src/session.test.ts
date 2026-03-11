import { beforeEach, describe, expect, test, vi } from "vitest";

const cookiesMock = vi.fn();
const decryptMock = vi.fn();

vi.mock("next/headers", () => ({
  cookies: cookiesMock,
}));

vi.mock("./crypto.js", () => ({
  decrypt: decryptMock,
}));

describe("getSession", () => {
  beforeEach(() => {
    vi.resetModules();
    vi.clearAllMocks();
    delete process.env.ZITADEL_COOKIE_SECRET;
  });

  test("returns null when cookie secret is too short", async () => {
    process.env.ZITADEL_COOKIE_SECRET = "short";
    const mod = await import("./session.js");

    await expect(mod.getSession()).resolves.toBeNull();
    expect(cookiesMock).not.toHaveBeenCalled();
  });

  test("returns session for valid epoch-seconds expiry", async () => {
    process.env.ZITADEL_COOKIE_SECRET = "abcdefghijklmnopqrstuvwxyz123456";
    cookiesMock.mockResolvedValueOnce({
      get: vi.fn(() => ({ value: "encrypted" })),
    });
    decryptMock.mockResolvedValueOnce(
      JSON.stringify({
        accessToken: "access-token",
        expiresAt: Math.floor(Date.now() / 1000) + 3600,
      }),
    );
    const mod = await import("./session.js");

    await expect(mod.getSession()).resolves.toMatchObject({
      accessToken: "access-token",
    });
  });

  test("returns null for expired epoch-seconds expiry", async () => {
    process.env.ZITADEL_COOKIE_SECRET = "abcdefghijklmnopqrstuvwxyz123456";
    cookiesMock.mockResolvedValueOnce({
      get: vi.fn(() => ({ value: "encrypted" })),
    });
    decryptMock.mockResolvedValueOnce(
      JSON.stringify({
        accessToken: "access-token",
        expiresAt: Math.floor(Date.now() / 1000) - 60,
      }),
    );
    const mod = await import("./session.js");

    await expect(mod.getSession()).resolves.toBeNull();
  });
});

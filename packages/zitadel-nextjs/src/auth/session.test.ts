import { beforeEach, describe, expect, test, vi } from "vitest";

const createSessionRpc = vi.fn();
const setSessionRpc = vi.fn();
const getSessionRpc = vi.fn();
const deleteSessionRpc = vi.fn();
const createCallbackRpc = vi.fn();
const getOIDCSessionMock = vi.fn();

vi.mock("@zitadel/zitadel-js/auth/bearer-token", () => ({
  createBearerTokenTransport: vi.fn(() => ({})),
}));

vi.mock("@zitadel/zitadel-js/api/v2", () => ({
  createSessionServiceClient: vi.fn(() => ({
    createSession: createSessionRpc,
    setSession: setSessionRpc,
    getSession: getSessionRpc,
    deleteSession: deleteSessionRpc,
  })),
  createOIDCServiceClient: vi.fn(() => ({
    createCallback: createCallbackRpc,
  })),
}));

vi.mock("../session.js", () => ({
  getSession: getOIDCSessionMock,
}));

describe("auth/session", () => {
  beforeEach(() => {
    vi.resetModules();
    vi.clearAllMocks();
    delete process.env.ZITADEL_API_URL;
    delete process.env.ZITADEL_SERVICE_USER_TOKEN;
    delete process.env.ZITADEL_COOKIE_SECRET;
  });

  test("throws when api URL is missing", async () => {
    const mod = await import("./session.js");

    await expect(
      mod.createSession({ accessToken: "token" }),
    ).rejects.toThrow(
      "apiUrl option or ZITADEL_API_URL environment variable is required",
    );
  });

  test("throws when no access token source is available", async () => {
    const mod = await import("./session.js");
    getOIDCSessionMock.mockResolvedValueOnce(null);

    await expect(
      mod.createSession({ apiUrl: "https://api.example.com" }),
    ).rejects.toThrow(
      "accessToken option, ZITADEL_SERVICE_USER_TOKEN, or an active OIDC session is required",
    );
  });

  test("falls back to OIDC session access token", async () => {
    const mod = await import("./session.js");
    getOIDCSessionMock.mockResolvedValueOnce({
      accessToken: "oidc-token",
      expiresAt: Math.floor(Date.now() / 1000) + 3600,
    });
    createSessionRpc.mockResolvedValueOnce({
      sessionId: "1",
      sessionToken: "2",
    });

    await mod.createSession({ apiUrl: "https://api.example.com" });

    expect(createSessionRpc).toHaveBeenCalled();
  });

  test("normalizes oidc_ auth request IDs when creating callbacks", async () => {
    const mod = await import("./session.js");
    process.env.ZITADEL_API_URL = "https://api.example.com";
    process.env.ZITADEL_SERVICE_USER_TOKEN = "service-token";
    createCallbackRpc.mockResolvedValueOnce({
      callbackUrl: "https://app.example.com/callback",
    });

    const callbackUrl = await mod.createCallback({
      authRequestId: "oidc_12345",
      sessionId: "sid",
      sessionToken: "stok",
    });

    expect(createCallbackRpc).toHaveBeenCalledWith(
      expect.objectContaining({
        authRequestId: "12345",
      }),
    );
    expect(callbackUrl).toBe("https://app.example.com/callback");
  });
});

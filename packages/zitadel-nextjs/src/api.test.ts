import { beforeEach, describe, expect, test, vi } from "vitest";

const createGrpcTransportMock = vi.fn(() => ({}));
const createAuthorizationBearerInterceptorMock = vi.fn(
  () => (next: unknown) => next,
);
const newSystemTokenMock = vi.fn();
const getSessionMock = vi.fn();

vi.mock("@zitadel/zitadel-js", () => ({
  createGrpcTransport: createGrpcTransportMock,
  createAuthorizationBearerInterceptor: createAuthorizationBearerInterceptorMock,
}));

vi.mock("@zitadel/zitadel-js/node", () => ({
  newSystemToken: newSystemTokenMock,
}));

vi.mock("@zitadel/zitadel-js/v2", () => ({
  createUserServiceClient: vi.fn(() => ({})),
  createSettingsServiceClient: vi.fn(() => ({})),
  createSessionServiceClient: vi.fn(() => ({})),
  createOIDCServiceClient: vi.fn(() => ({})),
  createSAMLServiceClient: vi.fn(() => ({})),
  createOrganizationServiceClient: vi.fn(() => ({})),
  createFeatureServiceClient: vi.fn(() => ({})),
  createIdpServiceClient: vi.fn(() => ({})),
  createActionServiceClient: vi.fn(() => ({})),
}));

vi.mock("./session.js", () => ({
  getSession: getSessionMock,
}));

describe("createZitadelApiClient", () => {
  beforeEach(() => {
    vi.resetModules();
    vi.clearAllMocks();

    delete process.env.ZITADEL_API_URL;
    delete process.env.ZITADEL_SERVICE_USER_TOKEN;
    delete process.env.ZITADEL_SERVICE_USER_KEY_ID;
    delete process.env.ZITADEL_SERVICE_USER_ID;
    delete process.env.ZITADEL_SERVICE_USER_PRIVATE_KEY;
  });

  test("uses explicit accessToken first", async () => {
    const mod = await import("./api.js");

    await mod.createZitadelApiClient({
      apiUrl: "https://api.example.com",
      accessToken: "explicit-token",
    });

    expect(createAuthorizationBearerInterceptorMock).toHaveBeenCalledWith(
      "explicit-token",
    );
    expect(getSessionMock).not.toHaveBeenCalled();
    expect(newSystemTokenMock).not.toHaveBeenCalled();
  });

  test("falls back to service user token", async () => {
    process.env.ZITADEL_SERVICE_USER_TOKEN = "service-token";
    const mod = await import("./api.js");

    await mod.createZitadelApiClient({
      apiUrl: "https://api.example.com",
    });

    expect(createAuthorizationBearerInterceptorMock).toHaveBeenCalledWith(
      "service-token",
    );
    expect(getSessionMock).not.toHaveBeenCalled();
  });

  test("falls back to private key JWT credentials", async () => {
    process.env.ZITADEL_SERVICE_USER_KEY_ID = "kid";
    process.env.ZITADEL_SERVICE_USER_ID = "service-user-id";
    process.env.ZITADEL_SERVICE_USER_PRIVATE_KEY = "private-key";
    newSystemTokenMock.mockResolvedValueOnce("jwt-token");
    const mod = await import("./api.js");

    await mod.createZitadelApiClient({
      apiUrl: "https://api.example.com",
      serviceUserTokenExpiresInSeconds: 120,
    });

    expect(newSystemTokenMock).toHaveBeenCalledWith({
      keyId: "kid",
      key: "private-key",
      issuer: "service-user-id",
      audience: "https://api.example.com",
      expiresInSeconds: 120,
    });
    expect(createAuthorizationBearerInterceptorMock).toHaveBeenCalledWith(
      "jwt-token",
    );
    expect(getSessionMock).not.toHaveBeenCalled();
  });

  test("falls back to OIDC session token when service credentials are missing", async () => {
    getSessionMock.mockResolvedValueOnce({
      accessToken: "session-token",
      expiresAt: Math.floor(Date.now() / 1000) + 3600,
    });
    const mod = await import("./api.js");

    await mod.createZitadelApiClient({
      apiUrl: "https://api.example.com",
    });

    expect(createAuthorizationBearerInterceptorMock).toHaveBeenCalledWith(
      "session-token",
    );
  });

  test("throws on partial private key JWT configuration", async () => {
    process.env.ZITADEL_SERVICE_USER_KEY_ID = "kid-only";
    const mod = await import("./api.js");

    await expect(
      mod.createZitadelApiClient({
        apiUrl: "https://api.example.com",
      }),
    ).rejects.toThrow("Incomplete private key JWT configuration");
  });
});

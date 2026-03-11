import { beforeEach, describe, expect, test, vi } from "vitest";
import {
  createOIDCAuthorizationUrl,
  createOIDCEndSessionUrl,
  discoverOIDCAuthorizationServer,
  exchangeOIDCAuthorizationCode,
  refreshOIDCTokens,
} from "./oidc.js";

const {
  discoveryRequestMock,
  processDiscoveryResponseMock,
  validateAuthResponseMock,
  authorizationCodeGrantRequestMock,
  processAuthorizationCodeResponseMock,
  refreshTokenGrantRequestMock,
  processRefreshTokenResponseMock,
  noneMock,
} = vi.hoisted(() => ({
  discoveryRequestMock: vi.fn(),
  processDiscoveryResponseMock: vi.fn(),
  validateAuthResponseMock: vi.fn(),
  authorizationCodeGrantRequestMock: vi.fn(),
  processAuthorizationCodeResponseMock: vi.fn(),
  refreshTokenGrantRequestMock: vi.fn(),
  processRefreshTokenResponseMock: vi.fn(),
  noneMock: vi.fn(() => "none-auth"),
}));

vi.mock("oauth4webapi", () => ({
  discoveryRequest: discoveryRequestMock,
  processDiscoveryResponse: processDiscoveryResponseMock,
  validateAuthResponse: validateAuthResponseMock,
  authorizationCodeGrantRequest: authorizationCodeGrantRequestMock,
  processAuthorizationCodeResponse: processAuthorizationCodeResponseMock,
  refreshTokenGrantRequest: refreshTokenGrantRequestMock,
  processRefreshTokenResponse: processRefreshTokenResponseMock,
  None: noneMock,
}));

const issuerMetadata = {
  authorization_endpoint: "https://issuer.example.com/oauth/v2/authorize",
  end_session_endpoint: "https://issuer.example.com/oidc/v1/end_session",
};

describe("oidc helpers", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("discovers OIDC authorization server metadata", async () => {
    discoveryRequestMock.mockResolvedValueOnce({ ok: true });
    processDiscoveryResponseMock.mockResolvedValueOnce(issuerMetadata);

    const discovered = await discoverOIDCAuthorizationServer(
      "https://issuer.example.com",
    );

    expect(discoveryRequestMock).toHaveBeenCalledWith(
      new URL("https://issuer.example.com"),
      { algorithm: "oidc" },
    );
    expect(discovered).toBe(issuerMetadata);
  });

  test("creates authorization URL with PKCE/state/prompt", () => {
    const url = createOIDCAuthorizationUrl({
      authorizationServer: issuerMetadata,
      clientId: "client-id",
      redirectUri: "http://localhost:3000/api/auth/callback",
      scopes: ["openid", "profile"],
      state: "state-123",
      codeChallenge: "challenge-456",
      prompt: "login",
    });

    expect(url.toString()).toContain("response_type=code");
    expect(url.searchParams.get("client_id")).toBe("client-id");
    expect(url.searchParams.get("redirect_uri")).toBe(
      "http://localhost:3000/api/auth/callback",
    );
    expect(url.searchParams.get("scope")).toBe("openid profile");
    expect(url.searchParams.get("state")).toBe("state-123");
    expect(url.searchParams.get("code_challenge")).toBe("challenge-456");
    expect(url.searchParams.get("code_challenge_method")).toBe("S256");
    expect(url.searchParams.get("prompt")).toBe("login");
  });

  test("throws when authorization endpoint is missing", () => {
    expect(() =>
      createOIDCAuthorizationUrl({
        authorizationServer: {},
        clientId: "client-id",
        redirectUri: "http://localhost:3000/callback",
        scopes: ["openid"],
        state: "state-123",
        codeChallenge: "challenge-456",
      }),
    ).toThrow("authorization_endpoint");
  });

  test("exchanges authorization code and maps token response", async () => {
    validateAuthResponseMock.mockReturnValueOnce(
      new URLSearchParams([["code", "abc"]]),
    );
    authorizationCodeGrantRequestMock.mockResolvedValueOnce({ ok: true });
    processAuthorizationCodeResponseMock.mockResolvedValueOnce({
      access_token: "access-token",
      id_token: "id-token",
      refresh_token: "refresh-token",
      expires_in: 3600,
      scope: "openid profile",
      token_type: "bearer",
    });

    const result = await exchangeOIDCAuthorizationCode({
      authorizationServer: issuerMetadata,
      clientId: "client-id",
      callbackRequestUrl:
        "http://localhost:3000/api/auth/callback?code=abc&state=state-123",
      callbackUrl: "http://localhost:3000/api/auth/callback",
      expectedState: "state-123",
      codeVerifier: "verifier-123",
    });

    expect(validateAuthResponseMock).toHaveBeenCalled();
    expect(authorizationCodeGrantRequestMock).toHaveBeenCalled();
    expect(result).toMatchObject({
      accessToken: "access-token",
      idToken: "id-token",
      refreshToken: "refresh-token",
      expiresIn: 3600,
      tokenType: "bearer",
    });
  });

  test("refreshes tokens and maps response", async () => {
    refreshTokenGrantRequestMock.mockResolvedValueOnce({ ok: true });
    processRefreshTokenResponseMock.mockResolvedValueOnce({
      access_token: "new-access-token",
      refresh_token: "new-refresh-token",
      expires_in: 1800,
      token_type: "bearer",
    });

    const result = await refreshOIDCTokens({
      authorizationServer: issuerMetadata,
      clientId: "client-id",
      refreshToken: "refresh-token",
    });

    expect(refreshTokenGrantRequestMock).toHaveBeenCalled();
    expect(result).toMatchObject({
      accessToken: "new-access-token",
      refreshToken: "new-refresh-token",
      expiresIn: 1800,
      tokenType: "bearer",
    });
  });

  test("creates end-session URL when endpoint exists", () => {
    const url = createOIDCEndSessionUrl({
      authorizationServer: issuerMetadata,
      idTokenHint: "id-token",
      postLogoutRedirectUri: "http://localhost:3000/auth",
    });

    expect(url?.toString()).toContain("id_token_hint=id-token");
    expect(url?.searchParams.get("post_logout_redirect_uri")).toBe(
      "http://localhost:3000/auth",
    );
  });

  test("returns null when end-session endpoint is missing", () => {
    expect(
      createOIDCEndSessionUrl({
        authorizationServer: {},
      }),
    ).toBeNull();
  });
});

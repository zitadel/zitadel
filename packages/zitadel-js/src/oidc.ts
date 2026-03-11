import * as oauth from "oauth4webapi";

export type OIDCAuthorizationServer = oauth.AuthorizationServer;

export interface OIDCAuthenticationResult {
  accessToken: string;
  idToken?: string;
  refreshToken?: string;
  expiresIn?: number;
  scope?: string;
  tokenType: string;
}

export async function discoverOIDCAuthorizationServer(
  issuerUrl: string,
): Promise<OIDCAuthorizationServer> {
  const issuer = new URL(issuerUrl);
  const discoveryResponse = await oauth.discoveryRequest(issuer, {
    algorithm: "oidc",
  });
  return oauth.processDiscoveryResponse(issuer, discoveryResponse);
}

export function createOIDCAuthorizationUrl(options: {
  authorizationServer: OIDCAuthorizationServer;
  clientId: string;
  redirectUri: string;
  scopes: string[];
  state: string;
  codeChallenge: string;
  prompt?: string;
}): URL {
  const authorizationEndpoint = options.authorizationServer.authorization_endpoint;
  if (!authorizationEndpoint) {
    throw new Error(
      "OIDC issuer metadata is missing authorization_endpoint",
    );
  }

  const authUrl = new URL(authorizationEndpoint);
  authUrl.searchParams.set("response_type", "code");
  authUrl.searchParams.set("client_id", options.clientId);
  authUrl.searchParams.set("redirect_uri", options.redirectUri);
  authUrl.searchParams.set("scope", options.scopes.join(" "));
  authUrl.searchParams.set("state", options.state);
  authUrl.searchParams.set("code_challenge", options.codeChallenge);
  authUrl.searchParams.set("code_challenge_method", "S256");
  if (options.prompt) {
    authUrl.searchParams.set("prompt", options.prompt);
  }

  return authUrl;
}

function mapTokenResult(
  tokenResult: oauth.TokenEndpointResponse,
): OIDCAuthenticationResult {
  return {
    accessToken: tokenResult.access_token,
    idToken: tokenResult.id_token,
    refreshToken: tokenResult.refresh_token,
    expiresIn: tokenResult.expires_in,
    scope: tokenResult.scope,
    tokenType: tokenResult.token_type,
  };
}

export async function exchangeOIDCAuthorizationCode(options: {
  authorizationServer: OIDCAuthorizationServer;
  clientId: string;
  callbackRequestUrl: URL | string;
  callbackUrl: string;
  expectedState: string;
  codeVerifier: string;
}): Promise<OIDCAuthenticationResult> {
  const client: oauth.Client = { client_id: options.clientId };
  const callbackRequestUrl =
    typeof options.callbackRequestUrl === "string"
      ? new URL(options.callbackRequestUrl)
      : options.callbackRequestUrl;

  const params = oauth.validateAuthResponse(
    options.authorizationServer,
    client,
    callbackRequestUrl,
    options.expectedState,
  );

  const tokenResponse = await oauth.authorizationCodeGrantRequest(
    options.authorizationServer,
    client,
    oauth.None(),
    params,
    options.callbackUrl,
    options.codeVerifier,
  );

  const tokenResult = await oauth.processAuthorizationCodeResponse(
    options.authorizationServer,
    client,
    tokenResponse,
  );

  return mapTokenResult(tokenResult);
}

export async function refreshOIDCTokens(options: {
  authorizationServer: OIDCAuthorizationServer;
  clientId: string;
  refreshToken: string;
}): Promise<OIDCAuthenticationResult> {
  const client: oauth.Client = { client_id: options.clientId };
  const tokenResponse = await oauth.refreshTokenGrantRequest(
    options.authorizationServer,
    client,
    oauth.None(),
    options.refreshToken,
  );
  const tokenResult = await oauth.processRefreshTokenResponse(
    options.authorizationServer,
    client,
    tokenResponse,
  );

  return mapTokenResult(tokenResult);
}

export function createOIDCEndSessionUrl(options: {
  authorizationServer: OIDCAuthorizationServer;
  idTokenHint?: string;
  postLogoutRedirectUri?: string;
}): URL | null {
  const endSessionEndpoint = options.authorizationServer.end_session_endpoint;
  if (!endSessionEndpoint) {
    return null;
  }

  const logoutUrl = new URL(endSessionEndpoint);
  if (options.idTokenHint) {
    logoutUrl.searchParams.set("id_token_hint", options.idTokenHint);
  }
  if (options.postLogoutRedirectUri) {
    logoutUrl.searchParams.set(
      "post_logout_redirect_uri",
      options.postLogoutRedirectUri,
    );
  }

  return logoutUrl;
}

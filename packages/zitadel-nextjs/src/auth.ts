import * as oauth from "oauth4webapi";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { generatePKCE, generateState } from "@zitadel/zitadel-js";
import { encrypt, decrypt } from "./crypto.js";

// ---------------------------------------------------------------------------
// Configuration types
// ---------------------------------------------------------------------------

export interface AuthOptions {
  /** The ZITADEL issuer URL, e.g. `https://my.zitadel.cloud`. */
  issuerUrl: string;
  /** The OIDC client ID registered in ZITADEL. */
  clientId: string;
  /** The callback URL registered in ZITADEL (must match exactly). */
  callbackUrl: string;
  /** Scopes to request. Defaults to `["openid", "profile", "email"]`. */
  scopes?: string[];
  /**
   * Secret used to encrypt session cookies. Must be at least 32 characters.
   * Falls back to the `ZITADEL_COOKIE_SECRET` environment variable.
   */
  cookieSecret?: string;
  /** Additional prompt parameter, e.g. `"login"` to force re-authentication. */
  prompt?: string;
}

export interface CallbackResult {
  /** The OAuth 2.0 access token. */
  accessToken: string;
  /** The OIDC ID token (if returned). */
  idToken?: string;
  /** Unix timestamp (seconds) when the access token expires. */
  expiresAt: number;
  /** Refresh token (if offline_access scope was requested). */
  refreshToken?: string;
}

// ---------------------------------------------------------------------------
// Cookie names
// ---------------------------------------------------------------------------

const SESSION_COOKIE = "zitadel.session";
const PKCE_COOKIE = "zitadel.pkce";

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

function getCookieSecret(options: AuthOptions): string {
  const secret =
    options.cookieSecret ?? process.env.ZITADEL_COOKIE_SECRET;
  if (!secret || secret.length < 32) {
    throw new Error(
      "cookieSecret (or ZITADEL_COOKIE_SECRET env var) must be at least 32 characters",
    );
  }
  return secret;
}

async function discoverIssuer(issuerUrl: string) {
  const url = new URL(issuerUrl);
  const response = await oauth.discoveryRequest(url, {
    algorithm: "oidc",
  });
  return oauth.processDiscoveryResponse(url, response);
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

/**
 * Initiates an OIDC Authorization Code + PKCE sign-in flow.
 *
 * Generates PKCE parameters, stores them in an encrypted cookie,
 * and redirects the user to ZITADEL's authorization endpoint.
 */
export async function signIn(options: AuthOptions): Promise<never> {
  const secret = getCookieSecret(options);
  const issuerConfig = await discoverIssuer(options.issuerUrl);
  const scopes = options.scopes ?? ["openid", "profile", "email"];

  const { codeVerifier, codeChallenge } = await generatePKCE();
  const state = generateState();

  // Store PKCE state in an encrypted HTTP-only cookie
  const pkceData = JSON.stringify({ codeVerifier, state });
  const encrypted = await encrypt(pkceData, secret);

  const cookieStore = await cookies();
  cookieStore.set(PKCE_COOKIE, encrypted, {
    httpOnly: true,
    secure: true,
    sameSite: "lax",
    path: "/",
    maxAge: 600, // 10 minutes — generous for the auth flow
  });

  // Build the authorization URL
  const authUrl = new URL(
    issuerConfig.authorization_endpoint as string,
  );
  authUrl.searchParams.set("response_type", "code");
  authUrl.searchParams.set("client_id", options.clientId);
  authUrl.searchParams.set("redirect_uri", options.callbackUrl);
  authUrl.searchParams.set("scope", scopes.join(" "));
  authUrl.searchParams.set("state", state);
  authUrl.searchParams.set("code_challenge", codeChallenge);
  authUrl.searchParams.set("code_challenge_method", "S256");
  if (options.prompt) {
    authUrl.searchParams.set("prompt", options.prompt);
  }

  redirect(authUrl.toString());
  // redirect() throws NEXT_REDIRECT, this is unreachable
  throw new Error("redirect() failed");
}

/**
 * Handles the OIDC callback, exchanging the authorization code for tokens.
 *
 * Should be called from a Next.js Route Handler for the callback URL:
 * ```ts
 * // app/api/auth/callback/route.ts
 * import { handleCallback } from "@zitadel/nextjs";
 *
 * export async function GET(request: Request) {
 *   const result = await handleCallback(request, {
 *     issuerUrl: "https://my.zitadel.cloud",
 *     clientId: "my-client-id",
 *     callbackUrl: "http://localhost:3000/api/auth/callback",
 *   });
 *   // result contains accessToken, idToken, expiresAt
 * }
 * ```
 */
export async function handleCallback(
  request: Request,
  options: AuthOptions,
): Promise<CallbackResult> {
  const secret = getCookieSecret(options);
  const issuerConfig = await discoverIssuer(options.issuerUrl);

  // Retrieve and decrypt PKCE data
  const cookieStore = await cookies();
  const pkceCookie = cookieStore.get(PKCE_COOKIE);
  if (!pkceCookie?.value) {
    throw new Error("Missing PKCE cookie — was signIn() called first?");
  }

  const pkceJson = await decrypt(pkceCookie.value, secret);
  if (!pkceJson) {
    throw new Error("Failed to decrypt PKCE cookie");
  }
  const { codeVerifier, state } = JSON.parse(pkceJson) as {
    codeVerifier: string;
    state: string;
  };

  // Clear the PKCE cookie
  cookieStore.delete(PKCE_COOKIE);

  // Extract the authorization code from the callback URL
  const callbackUrl = new URL(request.url);

  // Validate state
  const returnedState = callbackUrl.searchParams.get("state");
  if (returnedState !== state) {
    throw new Error("State mismatch — possible CSRF attack");
  }

  // Exchange the authorization code for tokens
  const client: oauth.Client = {
    client_id: options.clientId,
  };

  const params = oauth.validateAuthResponse(
    issuerConfig,
    client,
    callbackUrl,
    state,
  );

  const tokenResponse = await oauth.authorizationCodeGrantRequest(
    issuerConfig,
    client,
    oauth.None(),
    params,
    options.callbackUrl,
    codeVerifier,
  );

  const tokenResult = await oauth.processAuthorizationCodeResponse(
    issuerConfig,
    client,
    tokenResponse,
  );

  const expiresAt = tokenResult.expires_in
    ? Math.floor(Date.now() / 1000) + tokenResult.expires_in
    : Math.floor(Date.now() / 1000) + 3600;

  const result: CallbackResult = {
    accessToken: tokenResult.access_token,
    idToken: tokenResult.id_token,
    expiresAt,
    refreshToken: tokenResult.refresh_token,
  };

  // Store the session in an encrypted cookie
  const sessionData = JSON.stringify(result);
  const encryptedSession = await encrypt(sessionData, secret);

  cookieStore.set(SESSION_COOKIE, encryptedSession, {
    httpOnly: true,
    secure: true,
    sameSite: "lax",
    path: "/",
    maxAge: tokenResult.expires_in ?? 3600,
  });

  return result;
}

/**
 * Signs out the current user by clearing the session cookie.
 *
 * Optionally redirects to ZITADEL's end_session endpoint for full logout.
 */
export async function signOut(options?: {
  /** The ZITADEL issuer URL. Required to redirect to the end_session endpoint. */
  issuerUrl?: string;
  /** URL to redirect to after logout. */
  postLogoutRedirectUri?: string;
  /** The cookie secret. Falls back to ZITADEL_COOKIE_SECRET env var. */
  cookieSecret?: string;
}): Promise<void> {
  const cookieStore = await cookies();

  // Retrieve the id_token_hint before clearing
  let idTokenHint: string | undefined;
  const sessionCookie = cookieStore.get(SESSION_COOKIE);
  if (sessionCookie?.value && options?.issuerUrl) {
    const secret =
      options.cookieSecret ?? process.env.ZITADEL_COOKIE_SECRET ?? "";
    if (secret.length >= 32) {
      const json = await decrypt(sessionCookie.value, secret);
      if (json) {
        const session = JSON.parse(json) as CallbackResult;
        idTokenHint = session.idToken;
      }
    }
  }

  // Clear session cookie
  cookieStore.delete(SESSION_COOKIE);

  // Redirect to end_session endpoint if issuerUrl is provided
  if (options?.issuerUrl) {
    const issuerConfig = await discoverIssuer(options.issuerUrl);
    const endSessionEndpoint = issuerConfig.end_session_endpoint;

    if (endSessionEndpoint) {
      const logoutUrl = new URL(endSessionEndpoint as string);
      if (idTokenHint) {
        logoutUrl.searchParams.set("id_token_hint", idTokenHint);
      }
      if (options.postLogoutRedirectUri) {
        logoutUrl.searchParams.set(
          "post_logout_redirect_uri",
          options.postLogoutRedirectUri,
        );
      }
      redirect(logoutUrl.toString());
    }
  }
}

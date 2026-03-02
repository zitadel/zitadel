/**
 * OIDC Authorization Code + PKCE flow for Next.js.
 *
 * Use this module when you want to **add "Login with ZITADEL"** to your app
 * via standard OIDC redirects. ZITADEL (or a custom login UI) handles
 * credential collection; your app receives tokens in the callback.
 *
 * If you need to **build your own login UI** that collects credentials
 * directly (password, passkey, TOTP, etc.), use `@zitadel/nextjs/auth/session`
 * instead — it wraps the ZITADEL Session API.
 *
 * @module
 */
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import {
  createOIDCAuthorizationUrl,
  createOIDCEndSessionUrl,
  discoverOIDCAuthorizationServer,
  exchangeOIDCAuthorizationCode,
  generatePKCE,
  generateState,
} from "@zitadel/zitadel-js/auth/oidc";
import { encrypt, decrypt } from "../crypto.js";

// ---------------------------------------------------------------------------
// Configuration types
// ---------------------------------------------------------------------------

/**
 * Options for the OIDC Authorization Code + PKCE flow.
 *
 * All connection parameters (`issuerUrl`, `clientId`, `callbackUrl`) can be
 * provided via environment variables so that `signIn()` can be called with
 * minimal (or zero) configuration.
 */
export interface OIDCOptions {
  /**
   * The ZITADEL issuer URL, e.g. `https://my.zitadel.cloud`.
   * Falls back to the `ZITADEL_ISSUER_URL` environment variable.
   */
  issuerUrl?: string;
  /**
   * The OIDC client ID registered in ZITADEL.
   * Falls back to the `ZITADEL_CLIENT_ID` environment variable.
   */
  clientId?: string;
  /**
   * The callback URL registered in ZITADEL (must match exactly).
   * Falls back to the `ZITADEL_CALLBACK_URL` environment variable.
   */
  callbackUrl?: string;
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

/**
 * Result of a successful OIDC callback token exchange.
 */
export interface OIDCCallbackResult {
  /** The OAuth 2.0 access token. */
  accessToken: string;
  /** The OIDC ID token (if returned). */
  idToken?: string;
  /** Unix timestamp (seconds) when the access token expires. */
  expiresAt: number;
  /** Refresh token (if offline_access scope was requested). */
  refreshToken?: string;
}

/**
 * Session data stored in the encrypted cookie after OIDC login.
 */
export interface OIDCSession {
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
// Cookie names — namespaced under `auth.oidc` to leave room for
// `auth.session.*` when the Session API module is added.
// ---------------------------------------------------------------------------

const SESSION_COOKIE = "zitadel.auth.oidc.session";
const PKCE_COOKIE = "zitadel.auth.oidc.pkce";

// ---------------------------------------------------------------------------
// Internal helpers — resolve options with env-var fallbacks
// ---------------------------------------------------------------------------

function getIssuerUrl(options?: Pick<OIDCOptions, "issuerUrl">): string {
  const value = options?.issuerUrl ?? process.env.ZITADEL_ISSUER_URL;
  if (!value) {
    throw new Error(
      "issuerUrl option or ZITADEL_ISSUER_URL environment variable is required",
    );
  }
  return value;
}

function getClientId(options?: Pick<OIDCOptions, "clientId">): string {
  const value = options?.clientId ?? process.env.ZITADEL_CLIENT_ID;
  if (!value) {
    throw new Error(
      "clientId option or ZITADEL_CLIENT_ID environment variable is required",
    );
  }
  return value;
}

function getCallbackUrl(options?: Pick<OIDCOptions, "callbackUrl">): string {
  const value = options?.callbackUrl ?? process.env.ZITADEL_CALLBACK_URL;
  if (!value) {
    throw new Error(
      "callbackUrl option or ZITADEL_CALLBACK_URL environment variable is required",
    );
  }
  return value;
}

function getCookieSecret(options?: Pick<OIDCOptions, "cookieSecret">): string {
  const secret =
    options?.cookieSecret ?? process.env.ZITADEL_COOKIE_SECRET;
  if (!secret || secret.length < 32) {
    throw new Error(
      "cookieSecret (or ZITADEL_COOKIE_SECRET env var) must be at least 32 characters",
    );
  }
  return secret;
}

// ---------------------------------------------------------------------------
// Public API — OIDC Authorization Code + PKCE
// ---------------------------------------------------------------------------

/**
 * Initiates an OIDC Authorization Code + PKCE sign-in flow.
 *
 * Generates PKCE parameters, stores them in an encrypted cookie,
 * and redirects the user to ZITADEL's authorization endpoint.
 *
 * All options can be omitted if the corresponding environment variables
 * are set (`ZITADEL_ISSUER_URL`, `ZITADEL_CLIENT_ID`, `ZITADEL_CALLBACK_URL`,
 * `ZITADEL_COOKIE_SECRET`).
 *
 * @example
 * ```ts
 * // app/api/auth/signin/route.ts
 * import { signIn } from "@zitadel/nextjs/auth/oidc";
 *
 * export async function GET() {
 *   // Zero-config when env vars are set
 *   await signIn();
 * }
 * ```
 */
export async function signIn(options?: OIDCOptions): Promise<never> {
  const issuerUrl = getIssuerUrl(options);
  const clientId = getClientId(options);
  const callbackUrl = getCallbackUrl(options);
  const secret = getCookieSecret(options);
  const scopes = options?.scopes ?? ["openid", "profile", "email"];

  const authorizationServer = await discoverOIDCAuthorizationServer(issuerUrl);

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

  const authUrl = createOIDCAuthorizationUrl({
    authorizationServer,
    clientId,
    redirectUri: callbackUrl,
    scopes,
    state,
    codeChallenge,
    prompt: options?.prompt,
  });

  redirect(authUrl.toString());
  // redirect() throws NEXT_REDIRECT — this line is unreachable
  throw new Error("redirect() failed");
}

/**
 * Handles the OIDC callback, exchanging the authorization code for tokens.
 *
 * Should be called from a Next.js Route Handler for the callback URL.
 *
 * @example
 * ```ts
 * // app/api/auth/callback/route.ts
 * import { handleCallback } from "@zitadel/nextjs/auth/oidc";
 *
 * export async function GET(request: Request) {
 *   const result = await handleCallback(request);
 *   // result contains accessToken, idToken, expiresAt
 *   return Response.redirect("/");
 * }
 * ```
 */
export async function handleCallback(
  request: Request,
  options?: OIDCOptions,
): Promise<OIDCCallbackResult> {
  const issuerUrl = getIssuerUrl(options);
  const clientId = getClientId(options);
  const callbackUrl = getCallbackUrl(options);
  const secret = getCookieSecret(options);

  const authorizationServer = await discoverOIDCAuthorizationServer(issuerUrl);

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
  const callbackRequestUrl = new URL(request.url);

  // Validate state
  const returnedState = callbackRequestUrl.searchParams.get("state");
  if (returnedState !== state) {
    throw new Error("State mismatch — possible CSRF attack");
  }

  const tokenResult = await exchangeOIDCAuthorizationCode({
    authorizationServer,
    clientId,
    callbackRequestUrl,
    callbackUrl,
    expectedState: state,
    codeVerifier,
  });

  const expiresAt = tokenResult.expiresIn
    ? Math.floor(Date.now() / 1000) + tokenResult.expiresIn
    : Math.floor(Date.now() / 1000) + 3600;

  const result: OIDCCallbackResult = {
    accessToken: tokenResult.accessToken,
    idToken: tokenResult.idToken,
    expiresAt,
    refreshToken: tokenResult.refreshToken,
  };

  // Store the session in an encrypted cookie
  const sessionData = JSON.stringify(result);
  const encryptedSession = await encrypt(sessionData, secret);

  cookieStore.set(SESSION_COOKIE, encryptedSession, {
    httpOnly: true,
    secure: true,
    sameSite: "lax",
    path: "/",
    maxAge: tokenResult.expiresIn ?? 3600,
  });

  return result;
}

/**
 * Reads the current OIDC session from the encrypted cookie.
 *
 * Returns `null` if there is no session or the cookie cannot be decrypted.
 */
export async function getOIDCSession(options?: {
  cookieSecret?: string;
}): Promise<OIDCSession | null> {
  const secret = getCookieSecret(options);
  const cookieStore = await cookies();
  const sessionCookie = cookieStore.get(SESSION_COOKIE);
  if (!sessionCookie?.value) {
    return null;
  }

  const json = await decrypt(sessionCookie.value, secret);
  if (!json) {
    return null;
  }

  const session = JSON.parse(json) as OIDCSession;

  // Check if the session has expired
  if (session.expiresAt < Math.floor(Date.now() / 1000)) {
    cookieStore.delete(SESSION_COOKIE);
    return null;
  }

  return session;
}

/**
 * Signs out the current user by clearing the OIDC session cookie.
 *
 * Optionally redirects to ZITADEL's end_session endpoint for full
 * RP-Initiated Logout (OpenID Connect RP-Initiated Logout 1.0).
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
  if (sessionCookie?.value) {
    const issuerUrl = options?.issuerUrl ?? process.env.ZITADEL_ISSUER_URL;
    if (issuerUrl) {
      const secret =
        options?.cookieSecret ?? process.env.ZITADEL_COOKIE_SECRET ?? "";
      if (secret.length >= 32) {
        const json = await decrypt(sessionCookie.value, secret);
        if (json) {
          const session = JSON.parse(json) as OIDCCallbackResult;
          idTokenHint = session.idToken;
        }
      }
    }
  }

  // Clear session cookie
  cookieStore.delete(SESSION_COOKIE);

  // Redirect to end_session endpoint if issuerUrl is available
  const issuerUrl = options?.issuerUrl ?? process.env.ZITADEL_ISSUER_URL;
  if (issuerUrl) {
    const authorizationServer =
      await discoverOIDCAuthorizationServer(issuerUrl);
    const logoutUrl = createOIDCEndSessionUrl({
      authorizationServer,
      idTokenHint,
      postLogoutRedirectUri: options?.postLogoutRedirectUri,
    });
    if (logoutUrl) {
      redirect(logoutUrl.toString());
    }
  }
}

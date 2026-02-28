import { cookies } from "next/headers";
import { isSessionValid } from "@zitadel/zitadel-js";
import { decrypt } from "./crypto.js";

export interface SessionData {
  /** The OAuth 2.0 access token. */
  accessToken: string;
  /** The OIDC ID token (if returned). */
  idToken?: string;
  /** Unix timestamp (seconds) when the access token expires. */
  expiresAt: number;
  /** Refresh token (if offline_access scope was requested). */
  refreshToken?: string;
}

const SESSION_COOKIE = "zitadel.auth.oidc.session";

/**
 * Retrieves the current session data from the encrypted session cookie.
 *
 * Returns `null` if no session exists or the session has expired.
 *
 * @param cookieSecret - Secret for decrypting the cookie.
 *   Falls back to the `ZITADEL_COOKIE_SECRET` env var.
 */
export async function getSession(
  cookieSecret?: string,
): Promise<SessionData | null> {
  const secret = cookieSecret ?? process.env.ZITADEL_COOKIE_SECRET;
  if (!secret || secret.length < 32) {
    return null;
  }

  const cookieStore = await cookies();
  const sessionCookie = cookieStore.get(SESSION_COOKIE);
  if (!sessionCookie?.value) {
    return null;
  }

  const json = await decrypt(sessionCookie.value, secret);
  if (!json) {
    return null;
  }

  const session = JSON.parse(json) as SessionData;

  // `expiresAt` is stored as epoch seconds in the cookie payload.
  if (!isSessionValid({ expiresAt: session.expiresAt * 1000 })) {
    return null;
  }

  return session;
}

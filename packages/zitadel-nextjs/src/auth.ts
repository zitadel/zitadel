export interface AuthOptions {
  /** The ZITADEL issuer URL. */
  issuerUrl: string;
  /** The OIDC client ID. */
  clientId: string;
  /** The callback URL registered in ZITADEL. */
  callbackUrl: string;
  /** Scopes to request. Defaults to ["openid", "profile", "email"]. */
  scopes?: string[];
}

export interface CallbackResult {
  accessToken: string;
  idToken?: string;
  expiresAt: number;
}

/**
 * Initiates an OIDC sign-in flow using PKCE.
 * Placeholder — to be implemented with Next.js redirect primitives.
 */
export async function signIn(_options: AuthOptions): Promise<void> {
  // TODO: implement PKCE flow with next/navigation redirect
}

/**
 * Signs out the current user by clearing the session.
 * Placeholder — to be implemented with Next.js cookie primitives.
 */
export async function signOut(): Promise<void> {
  // TODO: implement with next/headers cookies()
}

/**
 * Handles the OIDC callback, exchanging the authorization code for tokens.
 * Placeholder — to be implemented with Next.js request primitives.
 */
export async function handleCallback(_request: Request): Promise<CallbackResult | null> {
  // TODO: implement code exchange
  return null;
}

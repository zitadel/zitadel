export interface SessionData {
  accessToken: string;
  idToken?: string;
  expiresAt: number;
}

/**
 * Retrieves the current session data from the cookie store.
 * Placeholder — to be implemented with Next.js cookie primitives.
 */
export async function getSession(): Promise<SessionData | null> {
  // TODO: implement with next/headers cookies()
  return null;
}

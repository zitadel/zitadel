export interface SessionInfo {
  /** ISO 8601 expiration timestamp or epoch milliseconds. */
  expiresAt: string | number;
}

/**
 * Checks whether a session has expired based on its expiresAt field.
 */
export function isSessionExpired(session: SessionInfo): boolean {
  const expiresAt =
    typeof session.expiresAt === "string"
      ? new Date(session.expiresAt).getTime()
      : session.expiresAt;
  return Date.now() >= expiresAt;
}

/**
 * Checks whether a session is still valid (not expired).
 */
export function isSessionValid(session: SessionInfo): boolean {
  return !isSessionExpired(session);
}

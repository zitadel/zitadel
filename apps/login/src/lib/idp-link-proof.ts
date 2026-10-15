import { createSignedProof, isSignedProofValid } from "./signed-proof";

// covers the round trip to the external IdP, including a login (and MFA) there
export const IDP_LINK_PROOF_MAX_AGE_SECONDS = 900; // 15 minutes

/**
 * Creates the `linkFingerprint` passed through the IdP flow when a signed-in user links an IdP to
 * their account. It binds the session to link to the browser that started the flow and is signed
 * with a server-only secret, so it can't be computed by a client for another session.
 */
export function createIdpLinkProof(sessionId: string, fingerprintId: string, now: number = Date.now()): string {
  return createSignedProof("idp-link", [sessionId, fingerprintId], IDP_LINK_PROOF_MAX_AGE_SECONDS, now);
}

export function isIdpLinkProofValid(
  value: string | undefined,
  sessionId: string | undefined,
  fingerprintId: string | undefined,
  now: number = Date.now(),
): boolean {
  return isSignedProofValid("idp-link", value, [sessionId, fingerprintId], now);
}

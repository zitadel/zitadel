import "server-only";

import { createSignedProof, isSignedProofValid } from "@/lib/server/signed-proof";

// matches the default MaxIdPIntentLifetime and covers the complete-registration form
export const IDP_FLOW_PROOF_MAX_AGE_SECONDS = 3600; // 1 hour

/**
 * Creates the `flowFingerprint`, the one proof passed through every external IdP flow. It binds the
 * flow to the browser that started it (its binding id cookie) and, when a signed-in user links an
 * IdP, to that session, so a callback link opened in another browser (login CSRF) or a link proof
 * for another session is rejected. The "-" placeholder keeps the field list fixed and non-empty
 * for a plain sign-in, and an empty linkToSessionId counts as no session.
 */
export function createIdpFlowProof(bindingId: string, linkToSessionId?: string, now: number = Date.now()): string {
  return createSignedProof("idp-flow", [bindingId, linkToSessionId || "-"], IDP_FLOW_PROOF_MAX_AGE_SECONDS, now);
}

export function isIdpFlowProofValid(
  value: string | undefined,
  bindingId: string | undefined,
  linkToSessionId?: string,
  now: number = Date.now(),
): boolean {
  return isSignedProofValid("idp-flow", value, [bindingId, linkToSessionId || "-"], now);
}

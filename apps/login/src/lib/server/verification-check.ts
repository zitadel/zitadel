import "server-only";

import { getBindingIdCookie, getOrSetBindingId, setBrowserCookie } from "@/lib/server/fingerprint";
import { createSignedProof, isSignedProofValid } from "@/lib/server/signed-proof";
import { cookies } from "next/headers";

export const VERIFICATION_CHECK_COOKIE = "verificationCheck";
export const VERIFICATION_CHECK_MAX_AGE_SECONDS = 300; // 5 minutes

/**
 * Creates the value of the verification check cookie: `v1.<expiresAt>.<signature>`.
 * The signature is an HMAC over userId, the binding id of the browser and expiry, keyed with a server-only secret,
 * so it can neither be forged by a client nor be reused for another user, browser or after expiry.
 */
export function createVerificationCheck(userId: string, bindingId: string, now: number = Date.now()): string {
  return createSignedProof("verification-check", [userId, bindingId], VERIFICATION_CHECK_MAX_AGE_SECONDS, now);
}

export function isVerificationCheckValid(
  value: string | undefined,
  userId: string,
  bindingId: string | undefined,
  now: number = Date.now(),
): boolean {
  return isSignedProofValid("verification-check", value, [userId, bindingId], now);
}

/**
 * Marks the user as verified (email / invite code, or fresh registration) for this browser.
 * Must only be called after a genuine, server-side validated verification.
 */
export async function setUserVerificationCheckCookie(userId: string): Promise<void> {
  const bindingId = await getOrSetBindingId();

  await setBrowserCookie(
    VERIFICATION_CHECK_COOKIE,
    createVerificationCheck(userId, bindingId),
    VERIFICATION_CHECK_MAX_AGE_SECONDS,
  );
}

export async function hasValidUserVerificationCheckCookie(userId: string): Promise<boolean> {
  const cookiesList = await cookies();

  // only read cookie to prevent issues on page.tsx
  const bindingIdCookie = await getBindingIdCookie();

  return isVerificationCheckValid(cookiesList.get(VERIFICATION_CHECK_COOKIE)?.value, userId, bindingIdCookie?.value);
}

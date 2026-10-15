import { cookies } from "next/headers";
import { getFingerprintIdCookie, getOrSetFingerprintId } from "./fingerprint";
import { createSignedProof, isSignedProofValid } from "./signed-proof";

export const VERIFICATION_CHECK_COOKIE = "verificationCheck";
export const VERIFICATION_CHECK_MAX_AGE_SECONDS = 300; // 5 minutes

/**
 * Creates the value of the verification check cookie: `v1.<expiresAt>.<signature>`.
 * The signature is an HMAC over userId, fingerprintId and expiry, keyed with a server-only secret,
 * so it can neither be forged by a client nor be reused for another user, browser or after expiry.
 */
export function createVerificationCheck(userId: string, fingerprintId: string, now: number = Date.now()): string {
  return createSignedProof("verification-check", [userId, fingerprintId], VERIFICATION_CHECK_MAX_AGE_SECONDS, now);
}

export function isVerificationCheckValid(
  value: string | undefined,
  userId: string,
  fingerprintId: string | undefined,
  now: number = Date.now(),
): boolean {
  return isSignedProofValid("verification-check", value, [userId, fingerprintId], now);
}

/**
 * Marks the user as verified (email / invite code, or fresh registration) for this browser.
 * Must only be called after a genuine, server-side validated verification.
 */
export async function setUserVerificationCheckCookie(userId: string): Promise<void> {
  const cookiesList = await cookies();
  const fingerprintId = await getOrSetFingerprintId();

  await cookiesList.set({
    name: VERIFICATION_CHECK_COOKIE,
    value: createVerificationCheck(userId, fingerprintId),
    httpOnly: true,
    path: "/",
    maxAge: VERIFICATION_CHECK_MAX_AGE_SECONDS,
  });
}

export async function hasValidUserVerificationCheckCookie(userId: string): Promise<boolean> {
  const cookiesList = await cookies();

  // only read cookie to prevent issues on page.tsx
  const fingerprintCookie = await getFingerprintIdCookie();

  return isVerificationCheckValid(cookiesList.get(VERIFICATION_CHECK_COOKIE)?.value, userId, fingerprintCookie?.value);
}

import { cookies } from "next/headers";
import { getFingerprintIdCookie, getOrSetFingerprintId } from "./fingerprint";
import { createSignedProof, isSignedProofValid } from "./signed-proof";

// Binds a device authorization user code to the browser that entered it,
// so consent can only be given in that browser.
const COOKIE_NAME = "deviceCodeCheck";
export const DEVICE_CODE_CHECK_MAX_AGE_SECONDS = 300; // 5 minutes

/**
 * Creates the value of the device code check cookie: `v1.<expiresAt>.<signature>`.
 * The signature is an HMAC over the user code, fingerprintId and expiry, keyed with a server-only secret,
 * like the other proofs of the login (see signed-proof.ts).
 */
export function createDeviceCodeCheck(userCode: string, fingerprintId: string, now: number = Date.now()): string {
  return createSignedProof("device-code", [userCode, fingerprintId], DEVICE_CODE_CHECK_MAX_AGE_SECONDS, now);
}

export function isDeviceCodeCheckValid(
  value: string | undefined,
  userCode: string,
  fingerprintId: string | undefined,
  now: number = Date.now(),
): boolean {
  return isSignedProofValid("device-code", value, [userCode, fingerprintId], now);
}

export async function setDeviceCodeCheckCookie(userCode: string) {
  const fingerprintId = await getOrSetFingerprintId();
  const cookiesList = await cookies();

  cookiesList.set({
    name: COOKIE_NAME,
    value: createDeviceCodeCheck(userCode, fingerprintId),
    httpOnly: true,
    path: "/",
    maxAge: DEVICE_CODE_CHECK_MAX_AGE_SECONDS,
  });
}

export async function checkDeviceCodeCheckCookie(userCode: string): Promise<boolean> {
  const cookiesList = await cookies();

  // only read cookies, this is also called from a page
  const fingerprintCookie = await getFingerprintIdCookie();

  const valid = isDeviceCodeCheckValid(cookiesList.get(COOKIE_NAME)?.value, userCode, fingerprintCookie?.value);

  if (!valid) {
    console.warn("Device code check failed.");
  }

  return valid;
}

export async function deleteDeviceCodeCheckCookie() {
  const cookiesList = await cookies();
  cookiesList.delete(COOKIE_NAME);
}

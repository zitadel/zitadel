import "server-only";

import { getBindingIdCookie, getOrSetBindingId, setBrowserCookie } from "@/lib/server/fingerprint";
import { createSignedProof, isSignedProofValid } from "@/lib/server/signed-proof";
import { cookies } from "next/headers";

// Binds a device authorization user code to the browser that entered it,
// so consent can only be given in that browser.
const COOKIE_NAME = "deviceCodeCheck";
export const DEVICE_CODE_CHECK_MAX_AGE_SECONDS = 300; // 5 minutes

/**
 * Creates the value of the device code check cookie: `v1.<expiresAt>.<signature>`.
 * The signature is an HMAC over the user code, the binding id of the browser and expiry, keyed with a server-only secret,
 * like the other proofs of the login (see signed-proof.ts).
 */
export function createDeviceCodeCheck(userCode: string, bindingId: string, now: number = Date.now()): string {
  return createSignedProof("device-code", [userCode, bindingId], DEVICE_CODE_CHECK_MAX_AGE_SECONDS, now);
}

export function isDeviceCodeCheckValid(
  value: string | undefined,
  userCode: string,
  bindingId: string | undefined,
  now: number = Date.now(),
): boolean {
  return isSignedProofValid("device-code", value, [userCode, bindingId], now);
}

export async function setDeviceCodeCheckCookie(userCode: string) {
  const bindingId = await getOrSetBindingId();

  await setBrowserCookie(COOKIE_NAME, createDeviceCodeCheck(userCode, bindingId), DEVICE_CODE_CHECK_MAX_AGE_SECONDS);
}

export async function checkDeviceCodeCheckCookie(userCode: string): Promise<boolean> {
  const cookiesList = await cookies();

  // only read cookies, this is also called from a page
  const bindingIdCookie = await getBindingIdCookie();

  const valid = isDeviceCodeCheckValid(cookiesList.get(COOKIE_NAME)?.value, userCode, bindingIdCookie?.value);

  if (!valid) {
    console.warn("Device code check failed.");
  }

  return valid;
}

export async function deleteDeviceCodeCheckCookie() {
  const cookiesList = await cookies();
  cookiesList.delete(COOKIE_NAME);
}

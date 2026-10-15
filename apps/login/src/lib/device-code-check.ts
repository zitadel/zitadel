import crypto from "crypto";
import { cookies } from "next/headers";
import { getFingerprintIdCookie, getOrSetFingerprintId } from "./fingerprint";

// Binds a device authorization user code to the browser that entered it,
// so consent can only be given in that browser.
const COOKIE_NAME = "deviceCodeCheck";

function hash(userCode: string, fingerprintId: string) {
  return crypto.createHash("sha256").update(`${userCode}:${fingerprintId}`).digest("hex");
}

export async function setDeviceCodeCheckCookie(userCode: string) {
  const fingerprintId = await getOrSetFingerprintId();
  const cookiesList = await cookies();

  cookiesList.set({
    name: COOKIE_NAME,
    value: hash(userCode, fingerprintId),
    httpOnly: true,
    path: "/",
    maxAge: 300, // 5 minutes
  });
}

export async function checkDeviceCodeCheckCookie(userCode: string): Promise<boolean> {
  const cookiesList = await cookies();

  // only read cookies, this is also called from a page
  const fingerPrintCookie = await getFingerprintIdCookie();

  if (!fingerPrintCookie || !fingerPrintCookie.value) {
    return false;
  }

  const cookieValue = cookiesList.get(COOKIE_NAME)?.value;

  if (!cookieValue) {
    console.warn("Device code check cookie not found. Device code check failed.");
    return false;
  }

  if (cookieValue !== hash(userCode, fingerPrintCookie.value)) {
    console.warn("Device code check failed.");
    return false;
  }

  return true;
}

export async function deleteDeviceCodeCheckCookie() {
  const cookiesList = await cookies();
  cookiesList.delete(COOKIE_NAME);
}

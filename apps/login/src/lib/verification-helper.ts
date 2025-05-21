"use server";

import crypto from "crypto";
import { cookies } from "next/headers";
import { getOrSetFingerprintId } from "./fingerprint";

export async function checkUserVerification(userId: string): Promise<boolean> {
  // check if a verification was done earlier
  const cookiesList = await cookies();
  const userAgentId = await getOrSetFingerprintId();

  const verificationCheck = crypto
    .createHash("sha256")
    .update(`${userId}:${userAgentId}`)
    .digest("hex");

  const cookieValue = await cookiesList.get("verificationCheck")?.value;

  if (!cookieValue) {
    console.warn(
      "User verification check cookie not found. User verification check failed.",
    );
    return false;
  }

  if (cookieValue !== verificationCheck) {
    console.warn(
      `User verification check failed. Expected ${verificationCheck} but got ${cookieValue}`,
    );
    return false;
  }

  return true;
}

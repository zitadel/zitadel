import crypto from "crypto";

/**
 * Returns a server-side secret used to key the user-verification cookie.
 * It is derived from the deployment credentials that are already required to run the login app
 * (see service.ts), so the value is never exposed to the client.
 */
function getServerSecret(): string {
  const secret =
    process.env.SYSTEM_USER_PRIVATE_KEY ||
    process.env.ZITADEL_SERVICE_USER_TOKEN ||
    process.env.ZITADEL_LOGINCLIENT_KEYFILE ||
    process.env.SYSTEM_USER_PRIVATE_KEY_FILE;

  if (!secret) {
    throw new Error("No server secret available to sign the user verification check");
  }

  return secret;
}

/**
 * Computes the keyed hash bound to a user and browser fingerprint for the verificationCheck cookie.
 * Uses HMAC with a server-side secret so the value cannot be forged from the userId and fingerprint
 * alone, which are both known to the client.
 */
export function computeUserVerificationCheck(userId: string, fingerprintId: string): string {
  return crypto.createHmac("sha256", getServerSecret()).update(`${userId}:${fingerprintId}`).digest("hex");
}

import { createHmac } from "crypto";
import {
  deriveSigningKey,
  getPrimarySigningSecret,
  getSessionCookieSecrets,
  signaturesEqual,
} from "./session-cookie-signature";

const VERSION = "v1";

/**
 * Purposes a signed proof can be issued for. Each purpose gets its own HKDF-derived key, so a proof
 * issued for one purpose can never be replayed as a proof for another one.
 */
export type SignedProofPurpose = "verification-check" | "idp-link" | "device-code";

function computeSignature(secret: string, purpose: SignedProofPurpose, fields: string[], expiresAt: number): string {
  // JSON array encoding keeps the fields unambiguous (no delimiter injection through any field)
  const payload = JSON.stringify([VERSION, ...fields, expiresAt]);
  // Domain separation: all proofs share the secret of the session cookie signature
  // (ZITADEL_SESSION_COOKIE_SECRET or the API credential) but use their own derived key per purpose.
  return createHmac("sha256", deriveSigningKey(secret, `zitadel-login-${purpose}-${VERSION}`))
    .update(payload)
    .digest("base64url");
}

/**
 * Creates a proof `v1.<expiresAt>.<signature>` binding `fields` together until `ttlSeconds` from now.
 * The signature is an HMAC keyed with a server-only secret, so the proof can neither be forged by a
 * client nor be reused for other field values or after expiry. The value is cookie and URL safe.
 */
export function createSignedProof(
  purpose: SignedProofPurpose,
  fields: string[],
  ttlSeconds: number,
  now: number = Date.now(),
): string {
  if (fields.length === 0 || fields.some((field) => !field)) {
    throw new Error(`All fields are required to create a ${purpose} proof.`);
  }

  const secret = getPrimarySigningSecret();

  const expiresAt = Math.floor(now / 1000) + ttlSeconds;
  return `${VERSION}.${expiresAt}.${computeSignature(secret, purpose, fields, expiresAt)}`;
}

export function isSignedProofValid(
  purpose: SignedProofPurpose,
  value: string | undefined,
  fields: (string | undefined)[],
  now: number = Date.now(),
): boolean {
  if (!value || fields.length === 0 || fields.some((field) => !field)) {
    return false;
  }

  const parts = value.split(".");
  if (parts.length !== 3) {
    return false;
  }

  const [version, expiresAtRaw, signature] = parts;
  if (version !== VERSION || !/^\d{1,12}$/.test(expiresAtRaw) || !signature) {
    return false;
  }

  const expiresAt = Number(expiresAtRaw);

  // check every candidate secret (no early return) to keep timing independent of which one matches
  let valid = false;
  for (const secret of getSessionCookieSecrets()) {
    if (signaturesEqual(signature, computeSignature(secret, purpose, fields as string[], expiresAt))) {
      valid = true;
    }
  }

  // expiry is enforced server side, a cookie maxAge is only a hint for the browser
  return valid && expiresAt > Math.floor(now / 1000);
}

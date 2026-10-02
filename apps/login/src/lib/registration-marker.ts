import { getSessionCookieSecrets, hmacWithDerivedKey, signaturesEqual } from "./session-cookie-signature";

/**
 * A registration marker proves that this browser received it from `registerUser` when the user was
 * created on the passkey path (no credential yet), within the last 15 minutes. It lets that browser
 * resume the registration if it did not finish setting up the passkey.
 *
 * Format: `<userId>.<expiresAtUnixSeconds>.<signature>`, where the signature is an HMAC-SHA256 over
 * `registration:<userId>:<expiresAt>:<organization>` with a key derived from the session cookie
 * signing secrets (separate HKDF info, so the key is used for nothing else).
 */
export const REGISTRATION_MARKER_COOKIE_NAME = "zitadel-registration";
export const REGISTRATION_MARKER_LIFETIME_SECONDS = 15 * 60;

const HKDF_INFO = "zitadel-login-registration-v1";

// User ids are numeric in ZITADEL. Rejecting separators keeps the cookie value and the signed payload
// unambiguous.
function isUsableUserId(userId: string): boolean {
  return !!userId && !/[.:]/.test(userId);
}

function signedPayload(userId: string, expiresAt: string, organization: string): string {
  return `registration:${userId}:${expiresAt}:${organization}`;
}

/**
 * Returns a signed marker for the user, or undefined if no signing secret is configured.
 */
export function createRegistrationMarker(userId: string, organization: string, nowMs = Date.now()): string | undefined {
  if (!isUsableUserId(userId)) {
    return undefined;
  }

  const [secret] = getSessionCookieSecrets();
  if (!secret) {
    return undefined;
  }

  const expiresAt = `${Math.floor(nowMs / 1000) + REGISTRATION_MARKER_LIFETIME_SECONDS}`;
  return `${userId}.${expiresAt}.${hmacWithDerivedKey(secret, HKDF_INFO, signedPayload(userId, expiresAt, organization))}`;
}

/**
 * Returns the user id of `value` if it is a marker for this organization that has not expired and is
 * signed with one of the configured secrets, otherwise undefined. Always undefined if no signing secret
 * is configured.
 */
export function readRegistrationMarker(
  value: string | undefined,
  organization: string,
  nowMs = Date.now(),
): string | undefined {
  if (!value) {
    return undefined;
  }

  const parts = value.split(".");
  if (parts.length !== 3) {
    return undefined;
  }

  const [userId, expiresAt, signature] = parts;
  if (!isUsableUserId(userId) || !/^\d{1,12}$/.test(expiresAt) || !signature) {
    return undefined;
  }

  if (Number(expiresAt) <= Math.floor(nowMs / 1000)) {
    return undefined;
  }

  const payload = signedPayload(userId, expiresAt, organization);
  // check every candidate secret (no early return) to keep timing independent of which one matches
  let valid = false;
  for (const secret of getSessionCookieSecrets()) {
    if (signaturesEqual(signature, hmacWithDerivedKey(secret, HKDF_INFO, payload))) {
      valid = true;
    }
  }
  return valid ? userId : undefined;
}

/**
 * True only if `value` is a valid marker (see readRegistrationMarker) for exactly this user and organization.
 */
export function verifyRegistrationMarker(
  value: string | undefined,
  expected: { userId: string; organization: string },
  nowMs = Date.now(),
): boolean {
  return readRegistrationMarker(value, expected.organization, nowMs) === expected.userId;
}

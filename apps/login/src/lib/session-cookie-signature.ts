import { createHmac, hkdfSync, timingSafeEqual } from "crypto";
import { readFileSync } from "fs";

const SIGNATURE_VERSION = "v1";
// Domain separation: the derived key is only ever used for session cookie signatures,
// even when the underlying secret is an API credential (PAT / private key).
const HKDF_INFO = `zitadel-login-session-cookie-${SIGNATURE_VERSION}`;
const DERIVED_KEY_LENGTH = 32;

export type SignableSession = {
  id: string;
  token: string;
  sig?: string;
};

// Keys are only loaded once from disk per process (same as api.ts).
const fileCache = new Map<string, string | undefined>();
const derivedKeyCache = new Map<string, Buffer>();

function readSecretFile(path: string): string | undefined {
  if (!fileCache.has(path)) {
    try {
      fileCache.set(path, readFileSync(path, "utf-8").trim() || undefined);
    } catch {
      fileCache.set(path, undefined);
    }
  }
  return fileCache.get(path);
}

function deriveKey(secret: string): Buffer {
  let key = derivedKeyCache.get(secret);
  if (!key) {
    key = Buffer.from(hkdfSync("sha256", secret, "", HKDF_INFO, DERIVED_KEY_LENGTH));
    derivedKeyCache.set(secret, key);
  }
  return key;
}

/**
 * Minimum length of each ZITADEL_SESSION_COOKIE_SECRET value. API credentials carry enough entropy
 * on their own; a typed secret does not, and HKDF cannot add entropy to a weak input.
 */
export const MIN_SESSION_COOKIE_SECRET_LENGTH = 32;

function getDedicatedSecrets(): string[] {
  return (process.env.ZITADEL_SESSION_COOKIE_SECRET ?? "")
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);
}

function getCredentialSecrets(): string[] {
  return [
    process.env.ZITADEL_SERVICE_USER_TOKEN,
    process.env.SYSTEM_USER_PRIVATE_KEY,
    process.env.SYSTEM_USER_PRIVATE_KEY_FILE ? readSecretFile(process.env.SYSTEM_USER_PRIVATE_KEY_FILE) : undefined,
    process.env.ZITADEL_LOGINCLIENT_KEYFILE ? readSecretFile(process.env.ZITADEL_LOGINCLIENT_KEYFILE) : undefined,
  ].filter((s): s is string => !!s);
}

/**
 * Returns an error message if ZITADEL_SESSION_COOKIE_SECRET is set but not usable, otherwise undefined.
 */
export function getSessionCookieSecretConfigError(): string | undefined {
  const tooShort = getDedicatedSecrets().filter((s) => s.length < MIN_SESSION_COOKIE_SECRET_LENGTH).length;
  if (tooShort > 0) {
    return (
      `ZITADEL_SESSION_COOKIE_SECRET is invalid: every value must be at least ${MIN_SESSION_COOKIE_SECRET_LENGTH} characters long ` +
      `(${tooShort} value(s) are shorter). Generate one with \`openssl rand -base64 32\`.`
    );
  }
  return undefined;
}

/**
 * Returns all secrets that may be used to sign / verify session cookie entries, in priority order.
 *
 * - ZITADEL_SESSION_COOKIE_SECRET: dedicated secret(s), the recommended configuration. A comma-separated
 *   list is supported so a secret can be rotated without invalidating existing cookies: the first entry
 *   is used for signing, all entries are accepted for verification. Every entry must be at least
 *   MIN_SESSION_COOKIE_SECRET_LENGTH characters long; if any is shorter, no secret is returned at all
 *   (signing fails and /ready reports not ready) instead of silently falling back to the credential.
 * - Otherwise (deprecated fallback) the API credential of the login is used, so applying the security
 *   patch does not require new configuration. All configured credentials are accepted for verification,
 *   which keeps cookies valid while a deployment migrates from one credential type to another
 *   (e.g. PAT -> login client key). This fallback will be removed in a future major release.
 */
export function getSessionCookieSecrets(): string[] {
  if (getSessionCookieSecretConfigError()) {
    return [];
  }

  const dedicated = getDedicatedSecrets();
  const credentials = getCredentialSecrets();

  return Array.from(new Set([...dedicated, ...credentials]));
}

/**
 * True if cookies are signed with a key derived from the API credential because no dedicated
 * ZITADEL_SESSION_COOKIE_SECRET is configured (deprecated fallback).
 */
export function isUsingCredentialFallback(): boolean {
  return getDedicatedSecrets().length === 0 && getCredentialSecrets().length > 0;
}

/**
 * Describes problems with the session cookie signing configuration, to be logged once at startup.
 */
export function getSessionCookieSecretStartupNotice(): { level: "error" | "warn"; message: string } | undefined {
  const configError = getSessionCookieSecretConfigError();
  if (configError) {
    return { level: "error", message: `${configError} The Login UI cannot sign session cookies and reports not ready.` };
  }
  if (!hasSessionCookieSecret()) {
    return {
      level: "error",
      message:
        "No session cookie signing secret available. Set ZITADEL_SESSION_COOKIE_SECRET. The Login UI reports not ready until a secret is configured.",
    };
  }
  if (isUsingCredentialFallback()) {
    return {
      level: "warn",
      message:
        "ZITADEL_SESSION_COOKIE_SECRET is not set; the session cookie signing key is derived from the configured API credential. " +
        "Rotating that credential signs all users out of the Login UI. Set ZITADEL_SESSION_COOKIE_SECRET. " +
        "This fallback is deprecated and will be removed in a future release.",
    };
  }
  return undefined;
}

export function hasSessionCookieSecret(): boolean {
  return getSessionCookieSecrets().length > 0;
}

/**
 * Canonical representation of a cookie entry for signing. All fields (except `sig`) are covered,
 * so no part of an entry can be tampered with. Keys are sorted and undefined values are dropped to
 * make the payload independent of key order and of the JSON round-trip through the cookie
 * (JSON.stringify omits undefined values, so `organization: undefined` and a missing key are equal).
 */
function canonicalPayload(session: SignableSession): string {
  const { sig: _ignored, ...fields } = session;
  const sorted = Object.fromEntries(
    Object.keys(fields)
      .sort()
      .filter((key) => (fields as Record<string, unknown>)[key] !== undefined)
      .map((key) => [key, (fields as Record<string, unknown>)[key]]),
  );
  return JSON.stringify([SIGNATURE_VERSION, sorted]);
}

function computeSignature(session: SignableSession, secret: string): string {
  return createHmac("sha256", deriveKey(secret)).update(canonicalPayload(session)).digest("base64url");
}

function signaturesEqual(left: string, right: string): boolean {
  const leftBytes = Buffer.from(left);
  const rightBytes = Buffer.from(right);
  if (leftBytes.length !== rightBytes.length) {
    return false;
  }
  return timingSafeEqual(leftBytes, rightBytes);
}

export function signSession<T extends SignableSession>(session: T): T & { sig: string } {
  if (!session.id || !session.token) {
    throw new Error("Session cookie entries require a non-empty id and token to be signed.");
  }

  const [secret] = getSessionCookieSecrets();
  if (!secret) {
    throw new Error(
      getSessionCookieSecretConfigError() ??
        "Session cookie signing secret is not configured. Set ZITADEL_SESSION_COOKIE_SECRET (at least 32 characters, e.g. `openssl rand -base64 32`).",
    );
  }

  const { sig: _ignored, ...unsigned } = session;
  return {
    ...(unsigned as T),
    sig: computeSignature(unsigned, secret),
  };
}

export function verifySession(session: unknown): session is SignableSession & { sig: string } {
  if (!session || typeof session !== "object") {
    return false;
  }

  const candidate = session as SignableSession;
  if (
    typeof candidate.id !== "string" ||
    !candidate.id ||
    typeof candidate.token !== "string" ||
    !candidate.token ||
    typeof candidate.sig !== "string" ||
    !candidate.sig
  ) {
    return false;
  }

  const { sig, ...unsigned } = candidate;
  // check every candidate secret (no early return) to keep timing independent of which one matches
  let valid = false;
  for (const secret of getSessionCookieSecrets()) {
    if (signaturesEqual(sig, computeSignature(unsigned, secret))) {
      valid = true;
    }
  }
  return valid;
}

export function stripSessionSignature<T extends SignableSession>(session: T): Omit<T, "sig"> {
  const { sig: _ignored, ...unsigned } = session;
  return unsigned;
}

export function parseAndVerifySessions<T extends SignableSession>(value: string | undefined): T[] {
  if (!value) {
    return [];
  }

  try {
    const parsed = JSON.parse(value);
    if (!Array.isArray(parsed)) {
      return [];
    }

    return parsed.filter(verifySession).map((session) => stripSessionSignature(session as T) as T);
  } catch {
    return [];
  }
}

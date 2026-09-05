/**
 * Helpers for comparing user supplied login names with stored values.
 *
 * ZITADEL treats usernames as case insensitive when it writes them: the
 * eventstore lowercases the username unique constraint, so a single
 * organization can never hold both "user@example.com" and "User@example.com".
 * Lookups therefore have to be case insensitive too, otherwise an account that
 * demonstrably exists is reported as not found — for example when an upstream
 * IdP returns a different casing than what was stored, or when a mobile
 * keyboard capitalizes the first letter of the address.
 */

/**
 * Normalizes raw login name input.
 *
 * Only surrounding whitespace is removed: the casing is kept so it can still be
 * echoed back to the user, and inner characters are left untouched because a
 * username is not required to be an email address.
 */
export function normalizeLoginName(loginName: string): string {
  return loginName.trim();
}

/**
 * Compares a user supplied login name with a stored one, ignoring case and
 * surrounding whitespace.
 *
 * An absent or empty value never matches. The callers use this to check whether
 * the input identified the user by their login name, email or phone, and an
 * unset email must not be treated as a match for an empty input.
 */
export function loginNameEquals(a: string | undefined, b: string | undefined): boolean {
  if (!a || !b) {
    return false;
  }
  return a.trim().toLowerCase() === b.trim().toLowerCase();
}

import { LANGS } from "@/lib/i18n";

/**
 * Authentication utility functions that don't require server actions
 */

/**
 * Compare two identifiers (loginname, email, phone) case insensitively.
 *
 * User lookup is performed with EQUALS_IGNORE_CASE, so any follow up check that
 * verifies *which* identifier the user was found by has to ignore case as well.
 * Otherwise entering "Jackson@example.com" for the user "jackson@example.com"
 * resolves the user but is then rejected as unknown.
 */
export function equalsIgnoreCase(a: string | undefined, b: string | undefined): boolean {
  if (a === undefined || b === undefined) {
    return false;
  }
  return a.toLowerCase() === b.toLowerCase();
}

/**
 * Check if a language code is valid (supported by the login UI)
 */
export function isValidLanguage(code: string): boolean {
  const normalized = code.trim().toLowerCase();
  return LANGS.some((lang) => lang.code === normalized);
}

/**
 * Extract a valid language code from uiLocales array.
 * Returns the first valid language code (normalized to lowercase), or null if none found.
 */
export function getValidLocaleFromUILocales(uiLocales: string[] | undefined): string | null {
  if (!uiLocales || uiLocales.length === 0) {
    return null;
  }

  for (const locale of uiLocales) {
    const normalized = locale.trim().toLowerCase();

    // Check if the full locale is a valid language code (e.g., "de", "EN")
    if (isValidLanguage(normalized)) {
      return normalized;
    }

    // uiLocales may contain language tags like "en-US" or "de-CH"
    // Extract the language code (part before the hyphen)
    // Note: this strips any regional specifier
    // e.g., de-CH and de-AT both become just de
    // zh-Hans-CN (Simplified) and zh-Hant-TW (Traditional) both become zh
    // As of time of writing, this is expected behaviour, since there is only one translation for all languages
    const languageCode = normalized.split("-")[0];
    if (isValidLanguage(languageCode)) {
      return languageCode;
    }
  }

  return null;
}

/**
 * Validate authentication request parameters
 */
export function validateAuthRequest(searchParams: URLSearchParams): string | null {
  const oidcRequestId = searchParams.get("authRequest");
  const samlRequestId = searchParams.get("samlRequest");

  const requestId =
    searchParams.get("requestId") ??
    (oidcRequestId ? `oidc_${oidcRequestId}` : samlRequestId ? `saml_${samlRequestId}` : undefined);

  return requestId || null;
}

/**
 * Check if request is an RSC request
 */
export function isRSCRequest(searchParams: URLSearchParams): boolean {
  return searchParams.has("_rsc");
}

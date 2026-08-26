import { LANGS } from "@/lib/i18n";

/**
 * Authentication utility functions that don't require server actions
 */

/**
 * Resolve a language tag to a supported language code.
 * The comparison is case insensitive, the returned code is the one declared in
 * LANGS, so region qualified codes keep their canonical casing (e.g. "zh-TW").
 */
export function resolveLanguage(code: string): string | null {
  const normalized = code.trim().toLowerCase();
  return LANGS.find((lang) => lang.code.toLowerCase() === normalized)?.code ?? null;
}

/**
 * Check if a language code is valid (supported by the login UI)
 */
export function isValidLanguage(code: string): boolean {
  return resolveLanguage(code) !== null;
}

/**
 * Extract a valid language code from uiLocales array.
 * Returns the first supported language code, or null if none found.
 */
export function getValidLocaleFromUILocales(uiLocales: string[] | undefined): string | null {
  if (!uiLocales || uiLocales.length === 0) {
    return null;
  }

  for (const locale of uiLocales) {
    const normalized = locale.trim().toLowerCase();

    // Check if the full locale is a valid language code (e.g., "de", "EN", "zh-TW")
    const exactMatch = resolveLanguage(normalized);
    if (exactMatch) {
      return exactMatch;
    }

    // uiLocales may contain language tags like "en-US" or "de-CH"
    // Fall back to the language code (part before the hyphen)
    // e.g., de-CH and de-AT both become just de
    // Region qualified languages are only kept if they are supported themselves,
    // so zh-Hans-CN falls back to zh while zh-TW matches the Traditional Chinese translation
    const languageCode = resolveLanguage(normalized.split("-")[0]);
    if (languageCode) {
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

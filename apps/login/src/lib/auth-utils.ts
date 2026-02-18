import { LANGS } from "@/lib/i18n";

/**
 * Authentication utility functions that don't require server actions
 */

/**
 * Check if a language code is valid (supported by the login UI)
 */
export function isValidLanguage(code: string): boolean {
  return LANGS.some((lang) => lang.code === code);
}

/**
 * Extract a valid language code from uiLocales array.
 * Returns the first valid language code, or null if none found.
 */
export function getValidLocaleFromUILocales(uiLocales: string[] | undefined): string | null {
  if (!uiLocales || uiLocales.length === 0) {
    return null;
  }
  
  for (const locale of uiLocales) {
    // uiLocales may contain language tags like "en-US" or "de-CH"
    // Extract the language code (part before the hyphen)
    const languageCode = locale.split("-")[0].toLowerCase();
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
  
  const requestId = searchParams.get("requestId") ??
    (oidcRequestId ? `oidc_${oidcRequestId}` : samlRequestId ? `saml_${samlRequestId}` : undefined);
  
  return requestId || null;
}

/**
 * Check if request is an RSC request
 */
export function isRSCRequest(searchParams: URLSearchParams): boolean {
  return searchParams.has("_rsc");
}
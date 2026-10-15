import { LANGS } from "@/lib/i18n";

/**
 * Authentication utility functions that don't require server actions
 */

/**
 * Check if a language code is valid (supported by the login UI)
 */
export function isValidLanguage(code: string): boolean {
  return LANGS.some((lang) => lang.code.toLowerCase() === code.trim().toLowerCase());
}

export function resolveLanguage(code: string): string | null {
  const normalized = code.trim().toLowerCase();
  const exact = LANGS.find((lang) => lang.code.toLowerCase() === normalized);
  if (exact) return exact.code;

  try {
    const locale = new Intl.Locale(normalized);
    if (locale.language === "zh" && locale.maximize().script === "Hant") return "zh-TW";
  } catch {
    return null;
  }

  return LANGS.find((lang) => lang.code === normalized.split("-")[0])?.code ?? null;
}

/**
 * Extract a valid language code from uiLocales array.
 * Returns the first supported language code, or null if none is found.
 */
export function getValidLocaleFromUILocales(uiLocales: string[] | undefined): string | null {
  if (!uiLocales || uiLocales.length === 0) {
    return null;
  }

  for (const locale of uiLocales) {
    const language = resolveLanguage(locale);
    if (language) return language;
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

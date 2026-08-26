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
 * Find a supported language sharing the language and script of the given
 * locale, using the CLDR data the runtime already ships. Returns null if the
 * locale cannot be parsed or carries no script information.
 */
function matchLanguageByScript(locale: string): string | null {
  try {
    const wanted = new Intl.Locale(locale).maximize();
    if (!wanted.script) {
      return null;
    }
    return (
      LANGS.find((lang) => {
        const candidate = new Intl.Locale(lang.code).maximize();
        return candidate.language === wanted.language && candidate.script === wanted.script;
      })?.code ?? null
    );
  } catch {
    return null;
  }
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

    // Then a language with the same script, matching how the backends
    // language.NewMatcher resolves tags. This is what maps zh-HK (Hant) onto
    // zh-TW instead of letting it fall through to zh (Hans).
    const scriptMatch = matchLanguageByScript(normalized);
    if (scriptMatch) {
      return scriptMatch;
    }

    // uiLocales may contain language tags like "en-US" or "de-CH"
    // Fall back to the language code (part before the hyphen)
    // e.g., de-CH and de-AT both become just de
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

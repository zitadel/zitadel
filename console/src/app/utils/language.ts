export const supportedLanguages = [
  'de',
  'en',
  'es',
  'fr',
  'id',
  'it',
  'ja',
  'pl',
  'zh',
  'zh-TW',
  'bg',
  'pt',
  'mk',
  'cs',
  'ru',
  'nl',
  'sv',
  'hu',
  'ko',
  'ro',
  'tr',
  'uk',
  'ar',
];
export const fallbackLanguage: string = 'en';

/**
 * Resolves a locale (browser or user profile) to a supported language, in the
 * same order as the backends `language.NewMatcher`: an exact match first, then
 * a language sharing the same script, then the base language.
 *
 * The script step is what makes `zh-HK` resolve to `zh-TW` (both Hant) instead
 * of falling through to `zh` (Hans). Languages with a single supported variant
 * are unaffected, `de-CH` still resolves to `de`.
 */
export function resolveSupportedLanguage(locale: string | undefined | null): string | undefined {
  if (!locale) {
    return undefined;
  }
  const exactMatch = supportedLanguages.find((lang) => lang.toLowerCase() === locale.toLowerCase());
  if (exactMatch) {
    return exactMatch;
  }
  const scriptMatch = matchByScript(locale);
  if (scriptMatch) {
    return scriptMatch;
  }
  const baseLanguage = locale.split('-')[0].toLowerCase();
  return supportedLanguages.find((lang) => lang.toLowerCase() === baseLanguage);
}

/**
 * Finds a supported language with the same language and script as the locale,
 * using the CLDR data the browser already ships. Returns undefined if the
 * locale cannot be parsed or no script information is available.
 */
function matchByScript(locale: string): string | undefined {
  try {
    const wanted = new Intl.Locale(locale).maximize();
    if (!wanted.script) {
      return undefined;
    }
    return supportedLanguages.find((lang) => {
      const candidate = new Intl.Locale(lang).maximize();
      return candidate.language === wanted.language && candidate.script === wanted.script;
    });
  } catch {
    return undefined;
  }
}

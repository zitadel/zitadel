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
 * Resolves a locale (browser or user profile) to a supported language.
 * An exact match wins over the base language, so a `zh-TW` browser gets the
 * Traditional Chinese translation instead of the Simplified Chinese one.
 */
export function resolveSupportedLanguage(locale: string | undefined | null): string | undefined {
  if (!locale) {
    return undefined;
  }
  const exactMatch = supportedLanguages.find((lang) => lang.toLowerCase() === locale.toLowerCase());
  if (exactMatch) {
    return exactMatch;
  }
  const baseLanguage = locale.split('-')[0].toLowerCase();
  return supportedLanguages.find((lang) => lang.toLowerCase() === baseLanguage);
}

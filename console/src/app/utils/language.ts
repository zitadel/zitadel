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

export function resolveSupportedLanguage(value: string | undefined, fallback = fallbackLanguage): string {
  if (!value) return fallback;
  const normalized = value.trim().toLowerCase();
  const exact = supportedLanguages.find((lang) => lang.toLowerCase() === normalized);
  if (exact) return exact;

  try {
    const locale = new Intl.Locale(normalized);
    if (locale.language === 'zh' && locale.maximize().script === 'Hant') return 'zh-TW';
  } catch {
    return fallback;
  }

  return supportedLanguages.find((lang) => lang === normalized.split('-')[0]) ?? fallback;
}

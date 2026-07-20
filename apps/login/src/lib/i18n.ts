export interface Lang {
  name: string;
  code: string;
}

export const LANGS: Lang[] = [
  {
    name: "English",
    code: "en",
  },
  {
    name: "Deutsch",
    code: "de",
  },
  {
    name: "Italiano",
    code: "it",
  },
  {
    name: "Español",
    code: "es",
  },
  {
    name: "Français",
    code: "fr",
  },
  {
    name: "Nederlands",
    code: "nl",
  },
  {
    name: "Polski",
    code: "pl",
  },
  {
    name: "Slovenčina",
    code: "sk",
  },
  {
    name: "Português",
    code: "pt",
  },
  {
    name: "简体中文",
    code: "zh",
  },
  {
    name: "Русский",
    code: "ru",
  },
  {
    name: "Türkçe",
    code: "tr",
  },
  {
    name: "日本語",
    code: "ja",
  },
  {
    name: "Українська",
    code: "uk",
  },
  {
    name: "العربية",
    code: "ar",
  },
];

export const LANGUAGE_COOKIE_NAME = "NEXT_LOCALE";
export const LANGUAGE_HEADER_NAME = "accept-language";

/** Resolves a BCP 47 tag to a locale bundled with the Login UI. */
export function normalizeLanguageCode(code: string | undefined): string | undefined {
  if (!code) return undefined;
  const normalized = code.trim().toLowerCase();
  if (!normalized) return undefined;
  if (LANGS.some((language) => language.code === normalized)) return normalized;
  const primaryLanguage = normalized.split("-")[0];
  return LANGS.some((language) => language.code === primaryLanguage) ? primaryLanguage : undefined;
}

export function shouldUILocalesOverrideCookie(): boolean {
  return process.env.ZITADEL_UI_LOCALES_OVERRIDE_COOKIE === "true";
}

export function getLanguage(code: string): Lang {
  const normalizedCode = normalizeLanguageCode(code) || code;
  const lang = LANGS.find((l) => l.code === normalizedCode);
  if (lang) {
    return lang;
  }

  return {
    code: normalizedCode,
    name: new Intl.DisplayNames([normalizedCode], { type: "language" }).of(normalizedCode) || normalizedCode,
  };
}

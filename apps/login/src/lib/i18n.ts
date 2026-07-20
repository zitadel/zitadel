export interface Lang {
  name: string;
  code: string;
}

export const LANGS: Lang[] = [
  { name: "Čeština", code: "cs" },
  { name: "Deutsch", code: "de" },
  { name: "English", code: "en" },
  { name: "Magyar", code: "hu" },
  { name: "Italiano", code: "it" },
  { name: "Lietuvių", code: "lt" },
  { name: "Polski", code: "pl" },
  { name: "Slovenčina", code: "sk" },
  { name: "Українська", code: "uk" },
];

export const LANGUAGE_COOKIE_NAME = "NEXT_LOCALE";
export const LANGUAGE_HEADER_NAME = "accept-language";

export function normalizeLanguageCode(code: string | undefined): string | undefined {
  const normalized = code?.trim().toLowerCase();
  if (!normalized) return undefined;
  const language = normalized.split("-")[0];
  return LANGS.some((lang) => lang.code === language) ? language : undefined;
}

export function shouldUILocalesOverrideCookie(): boolean {
  return process.env.ZITADEL_UI_LOCALES_OVERRIDE_COOKIE === "true";
}

export function getLanguage(code: string): Lang {
  const lang = LANGS.find((l) => l.code === code);
  if (lang) {
    return lang;
  }

  return {
    code,
    name: new Intl.DisplayNames([code], { type: "language" }).of(code) || code,
  };
}

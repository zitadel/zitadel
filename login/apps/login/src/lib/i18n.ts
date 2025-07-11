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
    name: "Polski",
    code: "pl",
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
];

export const LANGUAGE_COOKIE_NAME = "NEXT_LOCALE";
export const LANGUAGE_HEADER_NAME = "accept-language";

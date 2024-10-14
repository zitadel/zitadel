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
];

export const LANGUAGE_COOKIE_NAME = "NEXT_LOCALE";

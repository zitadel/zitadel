import deepmerge from "deepmerge";
import { getRequestConfig } from "next-intl/server";
import { cookies } from "next/headers";

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

export default getRequestConfig(async () => {
  const fallback = "en";
  const cookiesList = cookies();
  const locale: string = cookiesList.get(LANGUAGE_COOKIE_NAME)?.value ?? "en";

  const userMessages = (await import(`../../locales/${locale}.json`)).default;
  const fallbackMessages = (await import(`../../locales/${fallback}.json`))
    .default;

  return {
    locale,
    messages: deepmerge(fallbackMessages, userMessages),
  };
});

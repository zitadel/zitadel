import { LANGS, LANGUAGE_COOKIE_NAME, LANGUAGE_HEADER_NAME } from "@/lib/i18n";
import deepmerge from "deepmerge";
import { getRequestConfig } from "next-intl/server";
import { cookies, headers } from "next/headers";

export default getRequestConfig(async () => {
  const fallback = "en";
  const cookiesList = await cookies();

  let locale: string = fallback;

  const languageHeader = await (await headers()).get(LANGUAGE_HEADER_NAME);
  if (languageHeader) {
    const headerLocale = languageHeader.split(",")[0].split("-")[0]; // Extract the language code
    if (LANGS.map((l) => l.code).includes(headerLocale)) {
      locale = headerLocale;
    }
  }

  const languageCookie = cookiesList?.get(LANGUAGE_COOKIE_NAME);
  if (languageCookie && languageCookie.value) {
    if (LANGS.map((l) => l.code).includes(languageCookie.value)) {
      locale = languageCookie.value;
    }
  }

  const userMessages = (await import(`../../locales/${locale}.json`)).default;
  const fallbackMessages = (await import(`../../locales/${fallback}.json`))
    .default;

  return {
    locale,
    messages: deepmerge(fallbackMessages, userMessages),
  };
});

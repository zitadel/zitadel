import { LANGS, LANGUAGE_COOKIE_NAME, LANGUAGE_HEADER_NAME } from "@/lib/i18n";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { getHostedLoginTranslation } from "@/lib/zitadel";
import { JsonObject } from "@zitadel/client";
import deepmerge from "deepmerge";
import { getRequestConfig } from "next-intl/server";
import { cookies, headers } from "next/headers";

export default getRequestConfig(async () => {
  const fallback = "en";
  const cookiesList = await cookies();

  let locale: string = fallback;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

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

  const i18nOrganization = _headers.get("x-zitadel-i18n-organization") || ""; // You may need to set this header in middleware

  let translations: JsonObject | {} = {};
  try {
    const i18nJSON = await getHostedLoginTranslation({
      serviceUrl,
      locale,
      organization: i18nOrganization,
    });

    if (i18nJSON) {
      translations = i18nJSON;
    }
  } catch (error) {
    console.warn("Error fetching custom translations:", error);
  }

  const customMessages = translations;
  const localeMessages = (await import(`../../locales/${locale}.json`)).default;
  const fallbackMessages = (await import(`../../locales/${fallback}.json`))
    .default;

  return {
    locale,
    messages: deepmerge.all([
      fallbackMessages,
      localeMessages,
      customMessages,
    ]) as Record<string, string>,
  };
});

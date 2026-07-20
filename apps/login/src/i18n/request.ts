import { LANGS, LANGUAGE_COOKIE_NAME, LANGUAGE_HEADER_NAME, normalizeLanguageCode } from "@/lib/i18n";
import { getServiceConfig } from "@/lib/service-url";
import { getAllowedLanguages, getHostedLoginTranslation } from "@/lib/zitadel";
import { JsonObject } from "@zitadel/client";
import deepmerge from "deepmerge";
import { getRequestConfig } from "next-intl/server";
import { cookies, headers } from "next/headers";

export default getRequestConfig(async () => {
  const fallback = "en";
  const cookiesList = await cookies();

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  let allowedLanguages = LANGS.map((l) => l.code);
  let defaultLanguage = fallback;

  try {
    const settings = await getAllowedLanguages({ serviceConfig });
    if (settings.allowedLanguages?.length) {
      allowedLanguages = [...new Set(settings.allowedLanguages.map(normalizeLanguageCode).filter((language): language is string => !!language))];
    }
    const configuredDefaultLanguage = normalizeLanguageCode(settings.defaultLanguage);
    if (configuredDefaultLanguage && allowedLanguages.includes(configuredDefaultLanguage)) {
      defaultLanguage = configuredDefaultLanguage;
    }
  } catch (e) {
    console.warn("Failed to load global settings", e);
  }

  let locale: string = defaultLanguage;

  const languageHeader = await (await headers()).get(LANGUAGE_HEADER_NAME);
  if (languageHeader) {
    const headerLocale = normalizeLanguageCode(languageHeader.split(",")[0]);
    if (headerLocale && allowedLanguages.includes(headerLocale)) {
      locale = headerLocale;
    }
  }

  const languageCookie = cookiesList?.get(LANGUAGE_COOKIE_NAME);
  const cookieLocale = normalizeLanguageCode(languageCookie?.value);
  if (languageCookie && languageCookie.value) {
    if (cookieLocale && allowedLanguages.includes(cookieLocale)) {
      locale = cookieLocale;
    } else {
      // If the cookie tells a language that is other than the supported ones, fall back to the default.
      locale = defaultLanguage;
    }
  }

  const i18nOrganization = _headers.get("x-zitadel-i18n-organization") || ""; // You may need to set this header in middleware

  let translations: JsonObject | Record<string, never> = {};
  try {
    const i18nJSON = await getHostedLoginTranslation({
      serviceConfig,
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

  // Load locale messages, fall back to default language messages if locale not found
  let localeMessages;
  try {
    localeMessages = (await import(`../../locales/${locale}.json`)).default;
  } catch {
    try {
      localeMessages = (await import(`../../locales/${defaultLanguage}.json`)).default;
    } catch {
      localeMessages = (await import(`../../locales/${fallback}.json`)).default;
    }
  }

  const fallbackMessages = (await import(`../../locales/${fallback}.json`)).default;

  return {
    locale,
    messages: deepmerge.all([fallbackMessages, localeMessages, customMessages]) as Record<string, string>,
  };
});

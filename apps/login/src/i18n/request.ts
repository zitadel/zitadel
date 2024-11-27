import { LANGUAGE_COOKIE_NAME } from "@/lib/i18n";
import deepmerge from "deepmerge";
import { getRequestConfig } from "next-intl/server";
import { cookies } from "next/headers";

export default getRequestConfig(async () => {
  const fallback = "en";
  const cookiesList = await cookies();
  const locale: string = cookiesList.get(LANGUAGE_COOKIE_NAME)?.value ?? "en";

  const userMessages = (await import(`../../locales/${locale}.json`)).default;
  const fallbackMessages = (await import(`../../locales/${fallback}.json`))
    .default;

  return {
    locale,
    messages: deepmerge(fallbackMessages, userMessages),
  };
});

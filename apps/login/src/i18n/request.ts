import deepmerge from "deepmerge";
import { getRequestConfig } from "next-intl/server";
import { cookies } from "next/headers";

export default getRequestConfig(async () => {
  const fallback = "en";
  const cookiesList = cookies();
  const locale: string = cookiesList.get("NEXT_LOCALE")?.value ?? "en";

  const userMessages = (await import(`../../locales/${locale}.json`)).default;
  const fallbackMessages = (await import(`../../locales/${fallback}.json`))
    .default;

  console.log("i18nRequest", locale);

  return {
    locale,
    messages: deepmerge(fallbackMessages, userMessages),
  };
});

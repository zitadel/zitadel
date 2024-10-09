import { getRequestConfig } from "next-intl/server";
import { cookies } from "next/headers";

export default getRequestConfig(async () => {
  // Provide a static locale, fetch a user setting,
  // read from `cookies()`, `headers()`, etc.

  const cookiesList = cookies();
  const locale = cookiesList.get("NEXT_LOCALE")?.value ?? "en";

  return {
    locale,
    messages: (await import(`../../locales/${locale}.json`)).default,
  };
});

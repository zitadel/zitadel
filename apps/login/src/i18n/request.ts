import { getRequestConfig } from "next-intl/server";
import { cookies } from "next/headers";

export default getRequestConfig(async () => {
  const cookiesList = cookies();
  const locale: string = cookiesList.get("NEXT_LOCALE")?.value ?? "en";

  console.log("i18nRequest", locale);

  return {
    locale,
    messages: (await import(`../../locales/${locale}.json`)).default,
  };
});

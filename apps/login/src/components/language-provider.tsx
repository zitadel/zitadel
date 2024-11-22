"use cache";

import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";
import { unstable_cacheLife as cacheLife } from "next/cache";
import { ReactNode } from "react";

export async function LanguageProvider({ children }: { children: ReactNode }) {
  cacheLife("hours");

  const locale = await getLocale();
  const messages = await getMessages();
  return (
    <NextIntlClientProvider messages={messages}>
      {children}
    </NextIntlClientProvider>
  );
}

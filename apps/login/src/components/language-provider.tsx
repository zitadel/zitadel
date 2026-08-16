import { Lang } from "@/lib/i18n";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { ReactNode } from "react";
import { LanguagesProvider } from "./languages-context";

/**
 * Server component that provides both next-intl messages and the allowed
 * languages list to the client component tree.
 */
export async function LanguageProvider({ languages, children }: { languages: Lang[]; children: ReactNode }) {
  const messages = await getMessages();

  return (
    <NextIntlClientProvider messages={messages}>
      <LanguagesProvider languages={languages}>{children}</LanguagesProvider>
    </NextIntlClientProvider>
  );
}

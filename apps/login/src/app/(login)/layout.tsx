import "@/styles/globals.scss";

import { LanguageSwitcher } from "@/components/language-switcher";
import { Theme } from "@/components/theme";
import { ThemeProvider } from "@/components/theme-provider";
import { Analytics } from "@vercel/analytics/react";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";
import { Lato } from "next/font/google";
import { ReactNode } from "react";

const lato = Lato({
  weight: ["400", "700", "900"],
  subsets: ["latin"],
});

export const revalidate = 60; // revalidate every minute

export default async function RootLayout({
  children,
}: {
  children: ReactNode;
}) {
  const locale = await getLocale();
  const messages = await getMessages();

  return (
    <html
      lang={locale}
      className={`${lato.className}`}
      suppressHydrationWarning
    >
      <head />
      <body>
        <ThemeProvider>
          <NextIntlClientProvider messages={messages}>
            <div
              className={`h-screen overflow-y-scroll bg-background-light-600 dark:bg-background-dark-600`}
            >
              <div className="absolute bottom-0 right-0 flex flex-row p-4 items-center space-x-4">
                <LanguageSwitcher />
                <Theme />
              </div>

              <div className={`pb-4 flex flex-col justify-center h-full`}>
                <div className="mx-auto max-w-[440px] space-y-8 pt-20 lg:py-8 w-full">
                  {children}
                </div>
              </div>
            </div>

            <Analytics />
          </NextIntlClientProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}

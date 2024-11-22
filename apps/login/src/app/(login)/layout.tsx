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
              className={`relative min-h-screen bg-background-light-600 dark:bg-background-dark-600 flex flex-col justify-center`}
            >
              <div className="relative mx-auto max-w-[440px] py-8 w-full ">
                {children}
                <div className="flex flex-row justify-end py-4 items-center space-x-4">
                  <LanguageSwitcher />
                  <Theme />
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

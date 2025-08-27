import "@/styles/globals.scss";

import { LanguageProvider } from "@/components/language-provider";
import { LanguageSwitcher } from "@/components/language-switcher";
import { Skeleton } from "@/components/skeleton";
import { Theme } from "@/components/theme";
import { ThemeProvider } from "@/components/theme-provider";
import * as Tooltip from "@radix-ui/react-tooltip";
import { Analytics } from "@vercel/analytics/react";
import { Lato } from "next/font/google";
import { ReactNode, Suspense } from "react";
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";

const lato = Lato({
  weight: ["400", "700", "900"],
  subsets: ["latin"],
});

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("common");
  return { title: t('title')};
}

export default async function RootLayout({
  children,
}: {
  children: ReactNode;
}) {
  
  return (
    <html className={`${lato.className}`} suppressHydrationWarning>
      <head />
      <body>
        <ThemeProvider>
          <Tooltip.Provider>
            <Suspense
              fallback={
                <div
                  className={`relative flex min-h-screen flex-col justify-center bg-background-light-600 dark:bg-background-dark-600`}
                >
                  <div className="relative mx-auto w-full max-w-[440px] py-8">
                    <Skeleton>
                      <div className="h-40"></div>
                    </Skeleton>
                    <div className="flex flex-row items-center justify-end space-x-4 py-4">
                      <Theme />
                    </div>
                  </div>
                </div>
              }
            >
              <LanguageProvider>
                <div
                  className={`relative flex min-h-screen flex-col justify-center bg-background-light-600 dark:bg-background-dark-600`}
                >
                  <div className="relative mx-auto w-full max-w-[440px] py-8">
                    {children}
                    <div className="flex flex-row items-center justify-end space-x-4 py-4">
                      <LanguageSwitcher />
                      <Theme />
                    </div>
                  </div>
                </div>
              </LanguageProvider>
            </Suspense>
          </Tooltip.Provider>
        </ThemeProvider>
        <Analytics />
      </body>
    </html>
  );
}

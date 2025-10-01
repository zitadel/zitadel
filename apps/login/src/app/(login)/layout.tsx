import "@/styles/globals.scss";

import { BackgroundWrapper } from "@/components/background-wrapper";
import { LanguageProvider } from "@/components/language-provider";
import { LanguageSwitcher } from "@/components/language-switcher";
import { Skeleton } from "@/components/skeleton";
import { ThemeProvider } from "@/components/theme-provider";
import * as Tooltip from "@radix-ui/react-tooltip";
import { Lato } from "next/font/google";
import { ReactNode, Suspense } from "react";
import ThemeSwitch from "@/components/theme-switch";
import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";

const lato = Lato({
  weight: ["400", "700", "900"],
  subsets: ["latin"],
});

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("common");
  return { title: t("title") };
}

export default async function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html className={`${lato.className}`} suppressHydrationWarning>
      <head />
      <body>
        <ThemeProvider>
          <Tooltip.Provider>
            <Suspense
              fallback={
                <BackgroundWrapper
                  className={`relative flex min-h-screen flex-col justify-center bg-background-light-600 dark:bg-background-dark-600`}
                >
                  <div className="relative mx-auto w-full max-w-[440px] py-8">
                    <Skeleton>
                      <div className="h-40"></div>
                    </Skeleton>
                    <div className="flex flex-row items-center justify-end space-x-4 py-4">
                      <ThemeSwitch />
                    </div>
                  </div>
                </BackgroundWrapper>
              }
            >
              <LanguageProvider>
                <BackgroundWrapper
                  className={`relative flex min-h-screen flex-col justify-center bg-background-light-600 dark:bg-background-dark-600`}
                >
                  <div className="relative mx-auto w-full max-w-[1100px] py-8">
                    <div>{children}</div>
                    <div className="flex flex-row items-center justify-end space-x-4 py-4 px-4 md:px-8 max-w-[440px] mx-auto md:max-w-full">
                      <LanguageSwitcher />
                      <ThemeSwitch />
                    </div>
                  </div>
                </BackgroundWrapper>
              </LanguageProvider>
            </Suspense>
          </Tooltip.Provider>
        </ThemeProvider>
      </body>
    </html>
  );
}

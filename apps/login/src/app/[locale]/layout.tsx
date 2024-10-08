import "@/styles/globals.scss";

import { AddressBar } from "@/components/address-bar";
import { GlobalNav } from "@/components/global-nav";
import { LanguageSwitcher } from "@/components/language-switcher";
import { Theme } from "@/components/theme";
import { ThemeProvider } from "@/components/theme-provider";
import { TranslationsProvider } from "@/components/translations-provider";
import { Analytics } from "@vercel/analytics/react";
import i18nConfig from "i18nConfig";
import { dir } from "i18next";
import { Lato } from "next/font/google";
import { ReactNode } from "react";
import initTranslations from "../i18n";

const lato = Lato({
  weight: ["400", "700", "900"],
  subsets: ["latin"],
});

// export const revalidate = 60; // revalidate every minute

export function generateStaticParams() {
  return i18nConfig.locales.map((locale) => ({ locale }));
}

export default async function RootLayout({
  children,
  params: { locale },
}: {
  children: ReactNode;
  params: { locale: string };
}) {
  const i18nNamespaces = ["loginname"];
  console.log("layout:", locale);
  const { t, resources } = await initTranslations(locale, i18nNamespaces);

  // later only shown with dev mode enabled
  const showNav = process.env.DEBUG === "true";

  let domain = process.env.ZITADEL_API_URL;
  domain = domain ? domain.replace("https://", "") : "acme.com";

  return (
    <html
      lang={locale}
      dir={dir(locale)}
      className={`${lato.className}`}
      suppressHydrationWarning
    >
      <head />
      <body>
        <ThemeProvider>
          <TranslationsProvider
            namespaces={i18nNamespaces}
            locale={locale}
            resources={resources}
          >
            <div
              className={`h-screen overflow-y-scroll bg-background-light-600 dark:bg-background-dark-600 ${
                showNav
                  ? "bg-[url('/grid-light.svg')] dark:bg-[url('/grid-dark.svg')]"
                  : ""
              }`}
            >
              {showNav ? (
                <GlobalNav />
              ) : (
                <div className="absolute bottom-0 right-0 flex flex-row p-4 items-center space-x-4">
                  <LanguageSwitcher locale={locale} />
                  <Theme />
                </div>
              )}

              <div
                className={`${
                  showNav ? "lg:pl-72" : ""
                } pb-4 flex flex-col justify-center h-full`}
              >
                <div className="mx-auto max-w-[440px] space-y-8 pt-20 lg:py-8 w-full">
                  {showNav && (
                    <div className="rounded-lg bg-vc-border-gradient dark:bg-dark-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20">
                      <div className="rounded-lg bg-background-light-400 dark:bg-background-dark-500">
                        <AddressBar domain={domain} />
                      </div>
                    </div>
                  )}

                  {children}
                </div>
              </div>
            </div>

            <Analytics />
          </TranslationsProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}

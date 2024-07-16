import "@/styles/globals.scss";
import { AddressBar } from "@/ui/AddressBar";
import { GlobalNav } from "@/ui/GlobalNav";
import { Lato } from "next/font/google";
import { LayoutProviders } from "@/ui/LayoutProviders";
import { Analytics } from "@vercel/analytics/react";
import ThemeWrapper from "@/ui/ThemeWrapper";
import { getBrandingSettings } from "@/lib/zitadel";
import ThemeProvider from "@/ui/ThemeProvider";
import Theme from "@/ui/Theme";

const lato = Lato({
  weight: ["400", "700", "900"],
  subsets: ["latin"],
});

export const revalidate = 60; // revalidate every minute

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // later only shown with dev mode enabled
  const showNav = process.env.DEBUG === "true";

  let domain = process.env.ZITADEL_API_URL;
  domain = domain ? domain.replace("https://", "") : "acme.com";

  return (
    <html lang="en" className={`${lato.className}`} suppressHydrationWarning>
      <head />
      <body>
        <ThemeProvider>
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
              <div className="absolute bottom-0 right-0 flex flex-row p-4">
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
        </ThemeProvider>
      </body>
    </html>
  );
}

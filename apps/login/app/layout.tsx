"use client";

import "#/styles/globals.css";
// include styles from the ui package
import "@zitadel/react/styles.css";
import { AddressBar } from "#/ui/AddressBar";
import { GlobalNav } from "#/ui/GlobalNav";
import { Lato } from "next/font/google";
import { ThemeProvider } from "next-themes";
import Byline from "#/ui/Byline";

const lato = Lato({
  weight: "400",
  subsets: ["latin"],
});

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={`${lato.className}`}>
      <head />
      <ThemeProvider
        attribute="class"
        defaultTheme="system"
        storageKey="cp-theme"
        // value={{ dark: "dark ui-dark" }}
        // [color-scheme:dark]
      >
        <body className="overflow-y-scroll bg-background-light-600 dark:bg-background-dark-600 bg-[url('/grid-light.svg')] dark:bg-[url('/grid-dark.svg')]">
          <GlobalNav />

          <div className="lg:pl-72">
            <div className="mx-auto max-w-xl space-y-8 px-2 pt-20 lg:py-8 lg:px-8">
              <div className="rounded-lg bg-vc-border-gradient dark:bg-dark-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20">
                <div className="rounded-lg bg-background-light-500 dark:bg-background-dark-600">
                  <AddressBar />
                </div>
              </div>

              <div className="rounded-lg bg-vc-border-gradient dark:bg-dark-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20">
                <div className="rounded-lg bg-background-light-500 dark:bg-background-dark-500 p-3.5 lg:p-8">
                  {children}
                </div>
              </div>

              <div className="rounded-lg bg-vc-border-gradient dark:bg-dark-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20">
                <div className="rounded-lg bg-background-light-500 dark:bg-background-dark-600">
                  <Byline />
                </div>
              </div>
            </div>
          </div>
        </body>
      </ThemeProvider>
    </html>
  );
}

// export const metadata = () => {
//   return (
//     <>
//       <DefaultTags />
//       <title>ZITADEL Login Playground</title>
//       <meta
//         name="description"
//         content="This is a ZITADEL Login Playground to get an understanding how the login API works."
//       />
//     </>
//   );
// };

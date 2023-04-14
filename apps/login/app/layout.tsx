import "#/styles/globals.css";
// include styles from the ui package
import "@zitadel/react/styles.css";
import { AddressBar } from "#/ui/AddressBar";
import { GlobalNav } from "#/ui/GlobalNav";
import { ZitadelLogo } from "#/ui/ZitadelLogo";
import { Lato } from "@next/font/google";

const lato = Lato({
  weight: "400",
  subsets: ["latin"],
});

const darkModeClasses = (d: boolean) =>
  d ? "dark [color-scheme:dark] ui-dark" : "";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={`${darkModeClasses(false)} ${lato.className}`}>
      <head />
      <body className="overflow-y-scroll bg-background-light-600 dark:bg-background-dark-600 bg-[url('/grid.svg')]">
        <GlobalNav />

        <div className="lg:pl-72">
          <div className="mx-auto max-w-xl space-y-8 px-2 pt-20 lg:py-8 lg:px-8">
            <div className="rounded-lg bg-vc-border-gradient dark:dark-bg-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20">
              <div className="rounded-lg bg-background-light-500 dark:bg-background-dark-600">
                <AddressBar />
              </div>
            </div>

            <div className="rounded-lg bg-vc-border-gradient dark:dark-bg-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20">
              <div className="rounded-lg bg-background-light-500 dark:bg-background-dark-500 p-3.5 lg:p-8">
                {children}
              </div>
            </div>

            <div className="rounded-lg bg-vc-border-gradient dark:dark-bg-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20">
              <div className="rounded-lg bg-background-light-500 dark:bg-background-dark-600">
                <Byline />
              </div>
            </div>
          </div>
        </div>
      </body>
    </html>
  );
}

function Byline() {
  return (
    <div className="flex items-center p-3.5 lg:px-5 lg:py-3">
      <div className="flex items-center space-x-1.5">
        <div className="text-sm text-gray-600">By</div>
        {/* <a href="https://zitadel.com" title="ZITADEL">
          <div className=" text-gray-300 hover:text-gray-50">
            <ZitadelLogo />
          </div>
        </a> */}
        <div className="text-sm font-semibold">ZITADEL</div>
      </div>
    </div>
  );
}

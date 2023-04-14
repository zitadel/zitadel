import "#/styles/globals.css";
// include styles from the ui package
import "@zitadel/react/styles.css";
import { AddressBar } from "#/ui/AddressBar";
import { GlobalNav } from "#/ui/GlobalNav";
import { VercelLogo } from "#/ui/VercelLogo";
import { ZitadelLogo } from "#/ui/ZitadelLogo";
import { Lato } from "@next/font/google";

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
    <html lang="en" className={`dark [color-scheme:dark] ${lato.className}`}>
      <head />
      <body className="overflow-y-scroll bg-background-dark-600 bg-[url('/grid.svg')]">
        <GlobalNav />

        <div className="lg:pl-72">
          <div className="mx-auto max-w-xl space-y-8 px-2 pt-20 lg:py-8 lg:px-8">
            <div className="rounded-lg bg-vc-border-gradient p-px shadow-lg shadow-black/20">
              <div className="rounded-lg bg-background-dark-600">
                <AddressBar />
              </div>
            </div>

            <div className="rounded-lg bg-vc-border-gradient p-px shadow-lg shadow-black/20">
              <div className="rounded-lg bg-background-dark-500 p-3.5 lg:p-8">
                {children}
              </div>
            </div>

            <div className="rounded-lg bg-vc-border-gradient p-px shadow-lg shadow-black/20">
              <div className="rounded-lg bg-background-dark-600">
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
    <div className="flex items-center justify-between space-x-4 p-3.5 lg:px-5 lg:py-3">
      <div className="flex items-center space-x-1.5">
        <div className="text-sm text-gray-600">By</div>
        <a href="https://vercel.com" title="Vercel">
          <div className="w-7 text-gray-300 hover:text-gray-50">
            <ZitadelLogo />
          </div>
        </a>
        <div className="text-sm font-semibold">ZITADEL</div>
      </div>
      <div className="flex items-center space-x-1.5">
        <div className="text-sm text-gray-600">Deployed on</div>
        <a href="https://vercel.com" title="Vercel">
          <div className="w-16 text-gray-300 hover:text-gray-50">
            <VercelLogo />
          </div>
        </a>
      </div>
    </div>
  );
}

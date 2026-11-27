"use client";

import { Boundary } from "@/components/boundary";
import { Button } from "@/components/button";
import { ThemeWrapper } from "@/components/theme-wrapper";

// global-error replaces the root layout, so nothing the layout provides is available here:
// no NextIntlClientProvider (useTranslations and <Translated> would throw and hide the
// original error behind the built-in Next.js error page) and no language. Keep it static.
export default function GlobalError({ error, reset }: { error: Error & { digest?: string }; reset: () => void }) {
  return (
    // global-error must include html and body tags
    <html lang="en">
      <body>
        <ThemeWrapper branding={undefined}>
          <Boundary labels={["Login Error"]} color="red">
            <div className="space-y-4">
              <div className="text-sm text-red-500 dark:text-red-500">
                <span className="font-bold">Error:</span> {error?.message}
              </div>
              <div>
                <Button onClick={() => reset()}>Try Again</Button>
              </div>
            </div>
          </Boundary>
        </ThemeWrapper>
      </body>
    </html>
  );
}

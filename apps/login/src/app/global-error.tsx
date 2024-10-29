"use client";

import { Boundary } from "@/components/boundary";
import { Button } from "@/components/button";
import { ThemeWrapper } from "@/components/theme-wrapper";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    // global-error must include html and body tags
    <html>
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

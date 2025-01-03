"use client";

import { Logo } from "@/components/logo";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { ReactNode } from "react";
import { ThemeWrapper } from "./theme-wrapper";

export function DynamicTheme({
  branding,
  children,
}: {
  children: ReactNode;
  branding?: BrandingSettings;
}) {
  return (
    <ThemeWrapper branding={branding}>
      <div className="rounded-lg bg-background-light-400 dark:bg-background-dark-500 px-8 py-12">
        <div className="mx-auto flex flex-col items-center space-y-4">
          <div className="relative">
            {branding && (
              <Logo
                lightSrc={branding.lightTheme?.logoUrl}
                darkSrc={branding.darkTheme?.logoUrl}
                height={150}
                width={150}
              />
            )}
          </div>

          <div className="w-full">{children}</div>
          <div className="flex flex-row justify-between"></div>
        </div>
      </div>
    </ThemeWrapper>
  );
}

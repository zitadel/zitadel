"use client";

import { Logo } from "@/ui/Logo";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import React from "react";
import ThemeWrapper from "./ThemeWrapper";

export default function DynamicTheme({
  branding,
  children,
}: {
  children: React.ReactNode;
  branding?: BrandingSettings;
}) {
  return (
    <ThemeWrapper branding={branding}>
      <div className="rounded-lg bg-vc-border-gradient dark:bg-dark-vc-border-gradient p-px shadow-lg shadow-black/5 dark:shadow-black/20 mb-10">
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
      </div>
    </ThemeWrapper>
  );
}

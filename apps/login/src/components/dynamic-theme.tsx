import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { ReactNode } from "react";
import { DynamicThemeLayout } from "./dynamic-theme-layout";
import { LanguageProvider } from "./language-provider";

// Server entry for the page frame. Every page renders inside it, so the messages are provided per request here.
export function DynamicTheme({ branding, children }: { children: ReactNode; branding?: BrandingSettings }) {
  return (
    // the provider in the root layout is not rendered again on client navigation and would serve stale messages
    <LanguageProvider>
      <DynamicThemeLayout branding={branding}>{children}</DynamicThemeLayout>
    </LanguageProvider>
  );
}

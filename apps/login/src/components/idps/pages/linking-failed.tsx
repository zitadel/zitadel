"use client";

import { LanguageProvider } from "@/components/language-provider";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { useTranslations } from "next-intl";
import { Alert, AlertType } from "../../alert";
import { DynamicTheme } from "../../dynamic-theme";

export function linkingFailed(branding?: BrandingSettings) {
  const t = useTranslations("idp");

  return (
    <LanguageProvider>
      <DynamicTheme branding={branding}>
        <div className="flex flex-col items-center space-y-4">
          <h1>{t("linkingError.title")}</h1>
          <div className="w-full">
            {
              <Alert type={AlertType.ALERT}>
                {t("linkingError.description")}
              </Alert>
            }
          </div>
        </div>
      </DynamicTheme>
    </LanguageProvider>
  );
}

import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { Alert, AlertType } from "../../alert";
import { DynamicTheme } from "../../dynamic-theme";

export async function linkingFailed(branding?: BrandingSettings) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  return (
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
  );
}

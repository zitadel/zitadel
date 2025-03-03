import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { DynamicTheme } from "../../dynamic-theme";
import { IdpSignin } from "../../idp-signin";

export async function loginSuccess(
  userId: string,
  idpIntent: { idpIntentId: string; idpIntentToken: string },
  requestId?: string,
  branding?: BrandingSettings,
) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("loginSuccess.title")}</h1>
        <p className="ztdl-p">{t("loginSuccess.description")}</p>

        <IdpSignin
          userId={userId}
          idpIntent={idpIntent}
          requestId={requestId}
        />
      </div>
    </DynamicTheme>
  );
}

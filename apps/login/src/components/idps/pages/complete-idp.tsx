import { RegisterFormIDPIncomplete } from "@/components/register-form-idp-incomplete";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { IDPInformation } from "@zitadel/proto/zitadel/user/v2/idp_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { DynamicTheme } from "../../dynamic-theme";

export async function completeIDP({
  userId,
  idpInformation,
  requestId,
  organization,
  branding,
  idpIntent,
}: {
  userId: string;
  idpInformation: IDPInformation;
  requestId?: string;
  organization?: string;
  branding?: BrandingSettings;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("completeRegister.title")}</h1>
        <p className="ztdl-p">{t("completeRegister.description")}</p>

        <RegisterFormIDPIncomplete
          userId={userId}
          idpInformation={idpInformation}
          requestId={requestId}
          organization={organization}
          idpIntent={idpIntent}
        />
      </div>
    </DynamicTheme>
  );
}

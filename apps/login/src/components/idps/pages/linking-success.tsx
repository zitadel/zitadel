"use client";

import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { useTranslations } from "next-intl";
import { DynamicTheme } from "../../dynamic-theme";
import { IdpSignin } from "../../idp-signin";

export function linkingSuccess(
  userId: string,
  idpIntent: { idpIntentId: string; idpIntentToken: string },
  authRequestId?: string,
  branding?: BrandingSettings,
) {
  const t = useTranslations("idp");

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("linkingSuccess.title")}</h1>
        <p className="ztdl-p">{t("linkingSuccess.description")}</p>

        <IdpSignin
          userId={userId}
          idpIntent={idpIntent}
          authRequestId={authRequestId}
        />
      </div>
    </DynamicTheme>
  );
}

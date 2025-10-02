import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { DynamicTheme } from "../../dynamic-theme";
import { IdpSignin } from "../../idp-signin";
import { Translated } from "../../translated";

export async function linkingSuccess(
  userId: string,
  idpIntent: { idpIntentId: string; idpIntentToken: string },
  requestId?: string,
  branding?: BrandingSettings,
) {
  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="linkingSuccess.title" namespace="idp" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="linkingSuccess.description" namespace="idp" />
        </p>
      </div>

      <div className="w-full">
        <IdpSignin userId={userId} idpIntent={idpIntent} requestId={requestId} />
      </div>
    </DynamicTheme>
  );
}

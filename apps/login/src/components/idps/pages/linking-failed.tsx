import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { Alert, AlertType } from "../../alert";
import { DynamicTheme } from "../../dynamic-theme";
import { Translated } from "../../translated";

export async function linkingFailed(branding?: BrandingSettings, error?: string) {
  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="linkingError.title" namespace="idp" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="linkingError.description" namespace="idp" />
        </p>
        {error && <div className="w-full">{<Alert type={AlertType.ALERT}>{error}</Alert>}</div>}
      </div>
      <div className="w-full"></div>
    </DynamicTheme>
  );
}

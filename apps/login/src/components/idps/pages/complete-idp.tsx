import { RegisterFormIDPIncomplete } from "@/components/register-form-idp-incomplete";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { AddHumanUserRequest } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { DynamicTheme } from "../../dynamic-theme";
import { Translated } from "../../translated";

export async function completeIDP({
  idpUserId,
  idpId,
  idpUserName,
  addHumanUser,
  requestId,
  organization,
  branding,
  idpIntent,
}: {
  idpUserId: string;
  idpId: string;
  idpUserName: string;
  addHumanUser?: AddHumanUserRequest;
  requestId?: string;
  organization: string;
  branding?: BrandingSettings;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
}) {
  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>
          <Translated i18nKey="completeRegister.title" namespace="idp" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="completeRegister.description" namespace="idp" />
        </p>

        <RegisterFormIDPIncomplete
          idpUserId={idpUserId}
          idpId={idpId}
          idpUserName={idpUserName}
          defaultValues={{
            email: addHumanUser?.email?.email || "",
            firstname: addHumanUser?.profile?.givenName || "",
            lastname: addHumanUser?.profile?.familyName || "",
          }}
          requestId={requestId}
          organization={organization}
          idpIntent={idpIntent}
        />
      </div>
    </DynamicTheme>
  );
}

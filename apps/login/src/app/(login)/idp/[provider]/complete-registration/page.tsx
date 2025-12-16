import { DynamicTheme } from "@/components/dynamic-theme";
import { RegisterFormIDPIncomplete } from "@/components/register-form-idp-incomplete";
import { Translated } from "@/components/translated";
import { getServiceConfig } from "@/lib/service-url";
import { getBrandingSettings } from "@/lib/zitadel";
import { headers } from "next/headers";

/**
 * Complete registration page - shown when manual user registration is required
 */
export default async function CompleteRegistrationPage(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const searchParams = await props.searchParams;
  const { id, token, requestId, organization, idpId, idpUserId, idpUserName, givenName, familyName, email } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const branding = await getBrandingSettings({ serviceConfig, organization,
  });

  if (!id || !token || !idpId || !organization || !idpUserId || !idpUserName) {
    throw new Error("Missing required parameters");
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="completeRegister.title" namespace="idp" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="completeRegister.description" namespace="idp" />
        </p>
      </div>

      <div className="w-full">
        <RegisterFormIDPIncomplete
          idpUserId={idpUserId}
          idpId={idpId}
          idpUserName={idpUserName}
          defaultValues={{
            email: email || "",
            firstname: givenName || "",
            lastname: familyName || "",
          }}
          requestId={requestId}
          organization={organization}
          idpIntent={{
            idpIntentId: id,
            idpIntentToken: token,
          }}
        />
      </div>
    </DynamicTheme>
  );
}

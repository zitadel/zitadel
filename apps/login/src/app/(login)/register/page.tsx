import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { RegisterForm } from "@/components/register-form";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { Translated } from "@/components/translated";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import {
  getActiveIdentityProviders,
  getBrandingSettings,
  getDefaultOrg,
  getLegalAndSupportSettings,
  getLoginSettings,
  getPasswordComplexitySettings,
} from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("register");
  return { title: t("title") };
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  let { firstname, lastname, email, organization, requestId } = searchParams;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  if (!organization) {
    const org: Organization | null = await getDefaultOrg({
      serviceUrl,
    });
    if (org) {
      organization = org.id;
    }
  }

  const legal = await getLegalAndSupportSettings({
    serviceUrl,
    organization,
  });
  const passwordComplexitySettings = await getPasswordComplexitySettings({
    serviceUrl,
    organization,
  });

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization,
  });

  const identityProviders = await getActiveIdentityProviders({
    serviceUrl,
    orgId: organization,
  }).then((resp) => {
    return resp.identityProviders.filter((idp) => {
      return idp.options?.isAutoCreation || idp.options?.isCreationAllowed; // check if IDP allows to create account automatically or manual creation is allowed
    });
  });

  if (!loginSettings?.allowRegister) {
    return (
      <DynamicTheme branding={branding}>
        <div className="flex flex-col items-center space-y-4">
          <h1>
            <Translated i18nKey="disabled.title" namespace="register" />
          </h1>
          <p className="ztdl-p">
            <Translated i18nKey="disabled.description" namespace="register" />
          </p>
        </div>
      </DynamicTheme>
    );
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="register" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="description" namespace="register" />
        </p>

        {!organization && (
          <Alert>
            <Translated i18nKey="unknownContext" namespace="error" />
          </Alert>
        )}

        {legal &&
          passwordComplexitySettings &&
          organization &&
          (loginSettings.allowUsernamePassword || loginSettings.passkeysType == PasskeysType.ALLOWED) && (
            <RegisterForm
              idpCount={!loginSettings?.allowExternalIdp ? 0 : identityProviders.length}
              legal={legal}
              organization={organization}
              firstname={firstname}
              lastname={lastname}
              email={email}
              requestId={requestId}
              loginSettings={loginSettings}
            ></RegisterForm>
          )}

        {loginSettings?.allowExternalIdp && !!identityProviders.length && (
          <>
            <SignInWithIdp
              identityProviders={identityProviders}
              requestId={requestId}
              organization={organization}
            ></SignInWithIdp>
          </>
        )}
      </div>
    </DynamicTheme>
  );
}

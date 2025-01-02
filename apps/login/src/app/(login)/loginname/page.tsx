import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { UsernameForm } from "@/components/username-form";
import {
  getActiveIdentityProviders,
  getBrandingSettings,
  getDefaultOrg,
  getLoginSettings,
} from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "loginname" });

  const loginName = searchParams?.loginName;
  const authRequestId = searchParams?.authRequestId;
  const organization = searchParams?.organization;
  const submit: boolean = searchParams?.submit === "true";

  let defaultOrganization;
  if (!organization) {
    const org: Organization | null = await getDefaultOrg();
    if (org) {
      defaultOrganization = org.id;
    }
  }

  const loginSettings = await getLoginSettings(
    organization ?? defaultOrganization,
  );

  const identityProviders = await getActiveIdentityProviders(
    organization ?? defaultOrganization,
  ).then((resp) => {
    return resp.identityProviders;
  });

  const branding = await getBrandingSettings(
    organization ?? defaultOrganization,
  );

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("title")}</h1>
        <p className="ztdl-p">{t("description")}</p>

        <UsernameForm
          loginName={loginName}
          authRequestId={authRequestId}
          organization={organization} // stick to "organization" as we still want to do user discovery based on the searchParams not the default organization, later the organization is determined by the found user
          submit={submit}
          allowRegister={!!loginSettings?.allowRegister}
        >
          {identityProviders && (
            <SignInWithIdp
              identityProviders={identityProviders}
              authRequestId={authRequestId}
              organization={organization}
            ></SignInWithIdp>
          )}
        </UsernameForm>
      </div>
    </DynamicTheme>
  );
}

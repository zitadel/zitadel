import { AutoSubmitForm } from "@/components/auto-submit-form";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { Translated } from "@/components/translated";
import { UsernameForm } from "@/components/username-form";
import { isSafeRedirectUri } from "@/lib/client-utils";
import { idpTypeToSlug } from "@/lib/idp";
import { getPublicHost } from "@/lib/server/host";
import { getServiceConfig } from "@/lib/service-url";
import {
  getActiveIdentityProviders,
  getBrandingSettings,
  getDefaultOrg,
  getLoginSettings,
  startIdentityProviderFlow,
} from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { IdentityProviderType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";
import { redirect } from "next/navigation";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("loginname");
  return { title: t("title") };
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  const loginName = searchParams?.loginName;
  const requestId = searchParams?.requestId;
  const organization = searchParams?.organization;
  const orgDomain = searchParams?.orgDomain;
  const submit: boolean = searchParams?.submit === "true";

  // With an org domain suffix the login name may only be the local part (the
  // form shows the suffix separately), so put it back together for the IdP
  // login hint the same way sendLoginname does for the username form.
  const idpLoginHint = loginName && orgDomain && !loginName.includes("@") ? `${loginName}@${orgDomain}` : loginName;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  let defaultOrganization;
  if (!organization) {
    const org: Organization | null = await getDefaultOrg({ serviceConfig });
    if (org) {
      defaultOrganization = org.id;
    }
  }

  const loginSettings = await getLoginSettings({ serviceConfig, organization: organization ?? defaultOrganization });

  const identityProviders = await getActiveIdentityProviders({
    serviceConfig,
    orgId: organization ?? defaultOrganization,
  }).then((resp) => {
    return resp.identityProviders;
  });

  const branding = await getBrandingSettings({ serviceConfig, organization: organization ?? defaultOrganization });

  // Like Login V1: with local authentication disabled and exactly one external IdP, there is nothing to choose, so start the IdP flow right away.
  if (!loginSettings?.allowLocalAuthentication && loginSettings?.allowExternalIdp && identityProviders?.length === 1) {
    const idp = identityProviders[0];
    const provider = idpTypeToSlug(idp.type);

    const params = new URLSearchParams();
    if (requestId) params.set("requestId", requestId);
    if (organization) params.set("organization", organization);
    params.set("postErrorRedirectUrl", "/loginname");

    // redirect to LDAP page where username and password is requested
    if (idp.type === IdentityProviderType.LDAP) {
      params.set("idpId", idp.id);
      redirect("/idp/ldap?" + params.toString());
    }

    const host = getPublicHost(_headers);
    const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";
    const origin = `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}`;

    // on failure, fall through and render the page with the single IdP button
    const response = await startIdentityProviderFlow({
      serviceConfig,
      idpId: idp.id,
      urls: {
        successUrl: `${origin}/idp/${provider}/process?${params.toString()}`,
        failureUrl: `${origin}/idp/${provider}/failure?${params.toString()}`,
        loginHint: idpLoginHint,
      },
    }).catch(() => null);

    if (response?.url && isSafeRedirectUri(response.url)) {
      if (response.fields) {
        return (
          <DynamicTheme branding={branding}>
            <AutoSubmitForm url={response.url} fields={response.fields} />
          </DynamicTheme>
        );
      }

      redirect(response.url);
    }
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="loginname" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="description" namespace="loginname" />
        </p>
      </div>

      <div className="w-full">
        {loginSettings?.allowLocalAuthentication && (
          <UsernameForm
            loginName={loginName}
            requestId={requestId}
            organization={organization} // stick to "organization" as we still want to do user discovery based on the searchParams not the default organization, later the organization is determined by the found user
            defaultOrganization={defaultOrganization}
            loginSettings={loginSettings}
            suffix={orgDomain}
            hideSuffix={branding?.hideLoginNameSuffix}
            submit={submit}
            allowRegister={!!loginSettings?.allowRegister}
          ></UsernameForm>
        )}

        {loginSettings?.allowExternalIdp && !!identityProviders?.length && (
          <div className="w-full pt-6 pb-4">
            <SignInWithIdp
              identityProviders={identityProviders}
              requestId={requestId}
              organization={organization}
              postErrorRedirectUrl="/loginname"
              loginHint={idpLoginHint}
              showLabel={loginSettings?.allowLocalAuthentication}
            ></SignInWithIdp>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}

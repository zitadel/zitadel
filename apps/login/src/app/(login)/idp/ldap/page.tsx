import { DynamicTheme } from "@/components/dynamic-theme";
import { LDAPUsernamePasswordForm } from "@/components/ldap-username-password-form";
import { Translated } from "@/components/translated";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { getBrandingSettings, getDefaultOrg } from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<{ provider: string }>;
}) {
  const searchParams = await props.searchParams;
  const { idpId, organization, link } = searchParams;

  if (!idpId) {
    throw new Error("No idpId provided in searchParams");
  }

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  let defaultOrganization;
  if (!organization) {
    const org: Organization | null = await getDefaultOrg({
      serviceUrl,
    });
    if (org) {
      defaultOrganization = org.id;
    }
  }

  const branding = await getBrandingSettings({
    serviceUrl,
    organization: organization ?? defaultOrganization,
  });

  // return login failed if no linking or creation is allowed and no user was found
  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="ldap" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="description" namespace="ldap" />
        </p>

        <LDAPUsernamePasswordForm
          idpId={idpId}
          link={link === "true"}
        ></LDAPUsernamePasswordForm>
      </div>
    </DynamicTheme>
  );
}

import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { Translated } from "@/components/translated";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { getActiveIdentityProviders, getBrandingSettings } from "@/lib/zitadel";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;

  const requestId = searchParams?.requestId;
  const organization = searchParams?.organization;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const identityProviders = await getActiveIdentityProviders({
    serviceUrl,
    orgId: organization,
  }).then((resp) => {
    return resp.identityProviders;
  });

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>
          <Translated i18nKey="title" namespace="idp" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="description" namespace="idp" />
        </p>

        {!!identityProviders?.length && (
          <SignInWithIdp
            identityProviders={identityProviders}
            requestId={requestId}
            organization={organization}
          ></SignInWithIdp>
        )}
      </div>
    </DynamicTheme>
  );
}

import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { getActiveIdentityProviders, getBrandingSettings } from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

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
        <h1>{t("title")}</h1>
        <p className="ztdl-p">{t("description")}</p>

        {identityProviders && (
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

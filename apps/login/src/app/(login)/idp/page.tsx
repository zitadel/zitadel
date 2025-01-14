import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { getActiveIdentityProviders, getBrandingSettings } from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  const authRequestId = searchParams?.authRequestId;
  const organization = searchParams?.organization;

  const identityProviders = await getActiveIdentityProviders(organization).then(
    (resp) => {
      return resp.identityProviders;
    },
  );

  const branding = await getBrandingSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("title")}</h1>
        <p className="ztdl-p">{t("description")}</p>

        {identityProviders && (
          <SignInWithIdp
            identityProviders={identityProviders}
            authRequestId={authRequestId}
            organization={organization}
          ></SignInWithIdp>
        )}
      </div>
    </DynamicTheme>
  );
}

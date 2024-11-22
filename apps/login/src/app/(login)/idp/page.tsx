import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { getBrandingSettings, settingsService } from "@/lib/zitadel";
import { makeReqCtx } from "@zitadel/client/v2";
import { getLocale, getTranslations } from "next-intl/server";

function getIdentityProviders(orgId?: string) {
  return settingsService
    .getActiveIdentityProviders({ ctx: makeReqCtx(orgId) }, {})
    .then((resp) => {
      return resp.identityProviders;
    });
}

export default async function Page(
  props: {
    searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  }
) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "idp" });

  const authRequestId = searchParams?.authRequestId;
  const organization = searchParams?.organization;

  const identityProviders = await getIdentityProviders(organization);

  const host = process.env.VERCEL_URL
    ? `https://${process.env.VERCEL_URL}`
    : "http://localhost:3000";

  const branding = await getBrandingSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("title")}</h1>
        <p className="ztdl-p">{t("description")}</p>

        {identityProviders && (
          <SignInWithIdp
            host={host}
            identityProviders={identityProviders}
            authRequestId={authRequestId}
            organization={organization}
          ></SignInWithIdp>
        )}
      </div>
    </DynamicTheme>
  );
}

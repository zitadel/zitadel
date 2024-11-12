import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { UsernameForm } from "@/components/username-form";
import {
  getBrandingSettings,
  getDefaultOrg,
  getLoginSettings,
  settingsService,
} from "@/lib/zitadel";
import { makeReqCtx } from "@zitadel/client/v2";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { getLocale, getTranslations } from "next-intl/server";

function getIdentityProviders(orgId?: string) {
  return settingsService
    .getActiveIdentityProviders({ ctx: makeReqCtx(orgId) }, {})
    .then((resp) => {
      return resp.identityProviders;
    });
}

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "loginname" });

  const loginName = searchParams?.loginName;
  const authRequestId = searchParams?.authRequestId;
  let organization = searchParams?.organization;
  const submit: boolean = searchParams?.submit === "true";

  if (!organization) {
    const org: Organization | null = await getDefaultOrg().catch((error) => {
      console.warn(error);
      return null;
    });
    if (!org) {
      console.warn("No default organization found");
    } else {
      organization = org.id;
    }
  }

  const loginSettings = await getLoginSettings(organization);

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

        <UsernameForm
          loginName={loginName}
          authRequestId={authRequestId}
          organization={organization}
          submit={submit}
          allowRegister={!!loginSettings?.allowRegister}
        >
          {identityProviders && process.env.ZITADEL_API_URL && (
            <SignInWithIdp
              host={host}
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

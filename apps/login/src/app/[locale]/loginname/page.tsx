import initTranslations from "@/app/i18n";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SignInWithIdp } from "@/components/sign-in-with-idp";
import { UsernameForm } from "@/components/username-form";
import {
  getBrandingSettings,
  getLegalAndSupportSettings,
  getLoginSettings,
  settingsService,
} from "@/lib/zitadel";
import { makeReqCtx } from "@zitadel/client/v2";

function getIdentityProviders(orgId?: string) {
  return settingsService
    .getActiveIdentityProviders({ ctx: makeReqCtx(orgId) }, {})
    .then((resp) => {
      return resp.identityProviders;
    });
}

export default async function Page({
  searchParams,
  params: { locale },
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: { locale: string };
}) {
  const { t } = await initTranslations(locale, ["loginname"]);

  const loginName = searchParams?.loginName;
  const authRequestId = searchParams?.authRequestId;
  const organization = searchParams?.organization;
  const submit: boolean = searchParams?.submit === "true";

  const loginSettings = await getLoginSettings(organization);
  const legal = await getLegalAndSupportSettings();

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
          {legal && identityProviders && process.env.ZITADEL_API_URL && (
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

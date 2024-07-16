import {
  getBrandingSettings,
  getLegalAndSupportSettings,
  getLoginSettings,
  settingsService,
} from "@/lib/zitadel";
import DynamicTheme from "@/ui/DynamicTheme";
import { SignInWithIDP } from "@/ui/SignInWithIDP";
import UsernameForm from "@/ui/UsernameForm";
import { makeReqCtx } from "@zitadel/client2/v2beta";

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
        <h1>Welcome back!</h1>
        <p className="ztdl-p">Enter your login data.</p>

        <UsernameForm
          loginSettings={loginSettings}
          loginName={loginName}
          authRequestId={authRequestId}
          organization={organization}
          submit={submit}
        />

        {legal && identityProviders && process.env.ZITADEL_API_URL && (
          <SignInWithIDP
            host={host}
            identityProviders={identityProviders}
            authRequestId={authRequestId}
            organization={organization}
          ></SignInWithIDP>
        )}
      </div>
    </DynamicTheme>
  );
}

import { getBrandingSettings, settingsService } from "@/lib/zitadel";
import DynamicTheme from "@/ui/DynamicTheme";
import { SignInWithIDP } from "@/ui/SignInWithIDP";
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
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
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
        <h1>Sign in with SSO</h1>
        <p className="ztdl-p">
          Select one of the following providers to sign in
        </p>

        {identityProviders && (
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

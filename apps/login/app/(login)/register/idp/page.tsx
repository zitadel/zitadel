import { getLegalAndSupportSettings, server } from "#/lib/zitadel";
import { SignInWithIDP } from "#/ui/SignInWithIDP";
import {
  GetActiveIdentityProvidersResponse,
  IdentityProvider,
  ZitadelServer,
  settings,
} from "@zitadel/server";

function getIdentityProviders(
  server: ZitadelServer,
  orgId?: string
): Promise<IdentityProvider[] | undefined> {
  const settingsService = settings.getSettings(server);
  console.log("req");
  return settingsService
    .getActiveIdentityProviders(
      orgId ? { ctx: { orgId } } : { ctx: { instance: true } },
      {}
    )
    .then((resp: GetActiveIdentityProvidersResponse) => {
      return resp.identityProviders;
    });
}

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const legal = await getLegalAndSupportSettings(server);

  const identityProviders = await getIdentityProviders(server, "");

  const host = process.env.VERCEL_URL
    ? `https://${process.env.VERCEL_URL}`
    : "http://localhost:3000";

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Register</h1>
      <p className="ztdl-p">
        Select one of the following providers to register
      </p>

      {legal && identityProviders && process.env.ZITADEL_API_URL && (
        <SignInWithIDP
          host={host}
          identityProviders={identityProviders}
        ></SignInWithIDP>
      )}
    </div>
  );
}

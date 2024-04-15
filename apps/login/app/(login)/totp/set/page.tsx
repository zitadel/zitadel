import {
  addMyAuthFactorOTP,
  getBrandingSettings,
  getLoginSettings,
  getSession,
  server,
} from "#/lib/zitadel";
import DynamicTheme from "#/ui/DynamicTheme";
import TOTPRegister from "#/ui/TOTPRegister";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, organization } = searchParams;

  const branding = await getBrandingSettings(server, organization);
  const auth = await getMostRecentCookieWithLoginname(
    loginName,
    organization
  ).then((cookie) => {
    if (cookie) {
      return addMyAuthFactorOTP(cookie.token);
    } else {
      throw new Error("No cookie found");
    }
  });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Register TOTP</h1>
        <p className="ztdl-p">
          Scan the QR Code or navigate to the URL manually.
        </p>

        <div>
          {auth && <div>{auth.url}</div>}
          <TOTPRegister></TOTPRegister>
        </div>
      </div>
    </DynamicTheme>
  );
}

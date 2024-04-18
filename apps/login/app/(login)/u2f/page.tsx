import { getBrandingSettings, getLoginSettings, server } from "#/lib/zitadel";
import DynamicTheme from "#/ui/DynamicTheme";
import LoginOTP from "#/ui/LoginOTP";
import VerifyU2F from "#/ui/VerifyU2F";

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, authRequestId, sessionId, organization, code, submit } =
    searchParams;

  const branding = await getBrandingSettings(server, organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Verify 2-Factor</h1>

        <p className="ztdl-p">Verify your account with your device.</p>

        <VerifyU2F
          loginName={loginName}
          sessionId={sessionId}
          authRequestId={authRequestId}
          organization={organization}
        ></VerifyU2F>
      </div>
    </DynamicTheme>
  );
}

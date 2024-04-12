import { getBrandingSettings, getLoginSettings, server } from "#/lib/zitadel";
import DynamicTheme from "#/ui/DynamicTheme";
import TOTPForm from "#/ui/TOTPForm";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, authRequestId, sessionId, organization, code, submit } =
    searchParams;

  const branding = await getBrandingSettings(server, organization);
  const loginSettings = await getLoginSettings(server, organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Verify 2-Factor</h1>
        <p className="ztdl-p">Enter the code from your authenticator app. </p>

        <div>
          {loginSettings?.secondFactors.map((factor) => {
            return (
              <div>
                {factor === 1 && <div>TOTP</div>}
                {factor === 2 && <div>U2F</div>}
                {factor === 3 && <div>OTP Email</div>}
                {factor === 4 && <div>OTP Sms</div>}
              </div>
            );
          })}
        </div>
      </div>
    </DynamicTheme>
  );
}

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

  const { method } = params;

  const branding = await getBrandingSettings(server, organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Verify 2-Factor</h1>
        {method === "time-based" && (
          <p className="ztdl-p">Enter the code from your authenticator app.</p>
        )}
        {method === "sms" && (
          <p className="ztdl-p">Enter the code you got on your phone.</p>
        )}
        {method === "email" && (
          <p className="ztdl-p">Enter the code you got via your email.</p>
        )}
        {method === "u2f" && (
          <p className="ztdl-p">Verify your account with your device.</p>
        )}

        {method && ["time-based", "sms", "email"].includes(method) ? (
          <LoginOTP
            loginName={loginName}
            sessionId={sessionId}
            authRequestId={authRequestId}
            organization={organization}
            method={method}
          ></LoginOTP>
        ) : (
          <VerifyU2F
            loginName={loginName}
            sessionId={sessionId}
            authRequestId={authRequestId}
            organization={organization}
            submit={submit === "true"}
          ></VerifyU2F>
        )}
      </div>
    </DynamicTheme>
  );
}

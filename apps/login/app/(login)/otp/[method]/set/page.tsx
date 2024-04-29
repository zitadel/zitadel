import {
  addOTPEmail,
  addOTPSMS,
  getBrandingSettings,
  getSession,
  registerTOTP,
  server,
} from "#/lib/zitadel";
import Alert from "#/ui/Alert";
import DynamicTheme from "#/ui/DynamicTheme";
import TOTPRegister from "#/ui/TOTPRegister";
import UserAvatar from "#/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, organization, sessionId, authRequestId } = searchParams;
  const { method } = params;

  const branding = await getBrandingSettings(server, organization);
  const { session, token } = await loadSession(loginName, organization);

  let totpResponse;
  if (session && session.factors?.user?.id) {
    if (method === "time-based") {
      // inconsistency with token: email works with machine token, totp works with session token
      totpResponse = await registerTOTP(session.factors.user.id, token);
    } else if (method === "sms") {
      // does not work
      await addOTPSMS(session.factors.user.id);
    } else if (method === "email") {
      // works
      await addOTPEmail(session.factors.user.id);
    } else {
      throw new Error("Invalid method");
    }
  } else {
    throw new Error("No session found");
  }

  async function loadSession(loginName?: string, organization?: string) {
    const recent = await getMostRecentCookieWithLoginname(
      loginName,
      organization
    );

    return getSession(server, recent.id, recent.token).then((response) => {
      return { session: response?.session, token: recent.token };
    });
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Register 2-factor</h1>
        {!session && (
          <div className="py-4">
            <Alert>
              Could not get the context of the user. Make sure to enter the
              username first or provide a loginName as searchParam.
            </Alert>
          </div>
        )}

        {session && (
          <UserAvatar
            loginName={loginName ?? session.factors?.user?.loginName}
            displayName={session.factors?.user?.displayName}
            showDropdown
          ></UserAvatar>
        )}

        {totpResponse && "uri" in totpResponse && "secret" in totpResponse ? (
          <>
            <p className="ztdl-p">
              Scan the QR Code or navigate to the URL manually.
            </p>
            <div>
              {/* {auth && <div>{auth.to}</div>} */}

              <TOTPRegister
                uri={totpResponse.uri as string}
                secret={totpResponse.secret as string}
                loginName={loginName}
                sessionId={sessionId}
                authRequestId={authRequestId}
                organization={organization}
              ></TOTPRegister>
            </div>{" "}
          </>
        ) : (
          <p className="ztdl-p">
            {method === "email"
              ? "Code via email was successfully added."
              : method === "sms"
              ? "Code via SMS was successfully added."
              : ""}
          </p>
        )}
      </div>
    </DynamicTheme>
  );
}

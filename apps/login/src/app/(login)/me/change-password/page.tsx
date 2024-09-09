import { getSessionCookieById } from "@/lib/cookies";
import {
  getBrandingSettings,
  getPasswordComplexitySettings,
} from "@/lib/zitadel";
import Alert from "@/ui/Alert";
import ChangePasswordForm from "@/ui/ChangePasswordForm";
import DynamicTheme from "@/ui/DynamicTheme";
import UserAvatar from "@/ui/UserAvatar";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { sessionId } = searchParams;

  if (!sessionId) {
    return (
      <div>
        <h1>Session ID not found</h1>
      </div>
    );
  }

  const session = await getSessionCookieById({
    sessionId,
  });

  const sessionFactors = await loadMostRecentSession({
    loginName,
    organization,
  });

  const passwordComplexitySettings = await getPasswordComplexitySettings(
    session.organization,
  );

  const branding = await getBrandingSettings(session.organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Set Password</h1>
        <p className="ztdl-p">Set the password for your account</p>

        {(!sessionFactors || !loginName) && (
          <div className="py-4">
            <Alert>
              Could not get the context of the user. Make sure to enter the
              username first or provide a loginName as searchParam.
            </Alert>
          </div>
        )}

        {sessionFactors && (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {passwordComplexitySettings && (
          <ChangePasswordForm
            passwordComplexitySettings={passwordComplexitySettings}
            userId={""}
            sessionId={sessionId}
          ></ChangePasswordForm>
        )}
      </div>
    </DynamicTheme>
  );
}

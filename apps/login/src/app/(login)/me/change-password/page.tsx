import { getSessionCookieById } from "@/lib/cookies";
import {
  getBrandingSettings,
  getPasswordComplexitySettings,
  getSession,
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

  const sessionCookie = await getSessionCookieById({
    sessionId,
  });

  const { session } = await getSession({
    sessionId: sessionCookie.id,
    sessionToken: sessionCookie.token,
  });

  const passwordComplexitySettings = await getPasswordComplexitySettings(
    session?.factors?.user?.organizationId,
  );

  const branding = await getBrandingSettings(
    session?.factors?.user?.organizationId,
  );

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Set Password</h1>
        <p className="ztdl-p">Set the password for your account</p>

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
            loginName={session.factors?.user?.loginName}
            displayName={session.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {passwordComplexitySettings && session?.factors?.user?.id && (
          <ChangePasswordForm
            passwordComplexitySettings={passwordComplexitySettings}
            userId={session.factors.user.id}
            sessionId={sessionId}
          ></ChangePasswordForm>
        )}
      </div>
    </DynamicTheme>
  );
}

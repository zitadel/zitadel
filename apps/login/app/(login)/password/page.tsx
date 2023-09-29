import { getSession, server } from "#/lib/zitadel";
import Alert from "#/ui/Alert";
import PasswordForm from "#/ui/PasswordForm";
import UserAvatar from "#/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, promptPasswordless, authRequestId, alt } = searchParams;
  const sessionFactors = await loadSession(loginName);

  async function loadSession(loginName?: string) {
    const recent = await getMostRecentCookieWithLoginname(loginName);

    return getSession(server, recent.id, recent.token).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>{sessionFactors?.factors?.user?.displayName ?? "Password"}</h1>
      <p className="ztdl-p mb-6 block">Enter your password.</p>

      {!sessionFactors && (
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
        ></UserAvatar>
      )}

      <PasswordForm
        loginName={loginName}
        authRequestId={authRequestId}
        promptPasswordless={promptPasswordless === "true"}
        isAlternative={alt === "true"}
      />
    </div>
  );
}

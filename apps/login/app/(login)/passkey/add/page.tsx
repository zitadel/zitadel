import { getSession, server } from "#/lib/zitadel";
import Alert, { AlertType } from "#/ui/Alert";
import RegisterPasskey from "#/ui/RegisterPasskey";
import UserAvatar from "#/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName } = searchParams;
  const sessionFactors = await loadSession(loginName);

  async function loadSession(loginName?: string) {
    const recent = await getMostRecentCookieWithLoginname(loginName);
    return getSession(server, recent.id, recent.token).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }

  console.log(sessionFactors);
  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Register Passkey</h1>
      <p className="ztdl-p mb-6 block">
        Setup your user to authenticate with passkeys.
      </p>

      <Alert type={AlertType.INFO}>
        A passkey is an authentication method on a device like your fingerprint,
        Apple FaceID or similar.
      </Alert>

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
          loginName={loginName ?? sessionFactors.factors?.user?.loginName ?? ""}
          displayName={sessionFactors.factors?.user?.displayName}
          showDropdown
        ></UserAvatar>
      )}

      {sessionFactors?.id && <RegisterPasskey sessionId={sessionFactors.id} />}
    </div>
  );
}

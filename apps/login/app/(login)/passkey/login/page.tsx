import { getSession, server } from "#/lib/zitadel";
import Alert from "#/ui/Alert";
import LoginPasskey from "#/ui/LoginPasskey";
import UserAvatar from "#/ui/UserAvatar";
import { getMostRecentCookieWithLoginname } from "#/utils/cookies";

const title = "Authenticate with a passkey";
const description =
  "Your device will ask for your fingerprint, face, or screen lock";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, altPassword, authRequestId } = searchParams;

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
      <h1>{title}</h1>

      {sessionFactors && (
        <UserAvatar
          loginName={loginName ?? sessionFactors.factors?.user?.loginName}
          displayName={sessionFactors.factors?.user?.displayName}
          showDropdown
        ></UserAvatar>
      )}
      <p className="ztdl-p mb-6 block">{description}</p>

      {!sessionFactors && <div className="py-4"></div>}

      {!loginName && (
        <Alert>Provide your active session as loginName param</Alert>
      )}

      {loginName && (
        <LoginPasskey
          loginName={loginName}
          authRequestId={authRequestId}
          altPassword={altPassword === "true"}
        />
      )}
    </div>
  );
}

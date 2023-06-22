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
  const { loginName, prompt } = searchParams;

  const sessionFactors = await loadSession(loginName);

  async function loadSession(loginName?: string) {
    const recent = await getMostRecentCookieWithLoginname(loginName);
    return getSession(server, recent.id, recent.token).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }
  const title = !!prompt
    ? "Authenticate with a passkey"
    : "Use your passkey to confirm it's really you";
  const description = !!prompt
    ? "When set up, you will be able to authenticate without a password."
    : "Your device will ask for your fingerprint, face, or screen lock";

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

      {!sessionFactors && (
        <div className="py-4">
          <Alert>
            Could not get the context of the user. Make sure to enter the
            username first or provide a loginName as searchParam.
          </Alert>
        </div>
      )}

      {sessionFactors?.id && (
        <LoginPasskey sessionId={sessionFactors.id} isPrompt={!!prompt} />
      )}
    </div>
  );
}

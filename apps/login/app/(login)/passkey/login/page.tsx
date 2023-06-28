import {
  getSession,
  listAuthenticationMethodTypes,
  server,
  setSession,
} from "#/lib/zitadel";
import Alert, { AlertType } from "#/ui/Alert";
import LoginPasskey from "#/ui/LoginPasskey";
import RegisterPasskey from "#/ui/RegisterPasskey";
import UserAvatar from "#/ui/UserAvatar";
import {
  SessionCookie,
  getMostRecentCookieWithLoginname,
  updateSessionCookie,
} from "#/utils/cookies";
import { ChallengeKind } from "@zitadel/server";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName } = searchParams;

  const session = await setSessionForPasskeyChallenge(loginName);

  const challenge = session?.challenges?.passkey;

  //   let methods = [];
  //   if (sessionFactors?.factors?.user?.id) {
  //     methods = await listAuthenticationMethodTypes(
  //       sessionFactors.factors.user.id
  //     );

  //     console.log(methods);
  //   }

  async function setSessionForPasskeyChallenge(loginName?: string) {
    const recent = await getMostRecentCookieWithLoginname(loginName);
    console.log(recent);
    return setSession(server, recent.id, recent.token, undefined, [
      ChallengeKind.CHALLENGE_KIND_PASSKEY,
    ]).then((session) => {
      const sessionCookie: SessionCookie = {
        id: recent.id,
        token: session.sessionToken,
        changeDate: session.changeDate?.toString() ?? "",
        loginName: session.factors?.user?.loginName ?? "",
      };

      return updateSessionCookie(sessionCookie.id, sessionCookie).then(() => {
        return session;
      });
    });
  }

  const title = "Authenticate with a passkey";
  const description =
    "Your device will ask for your fingerprint, face, or screen lock";

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>{title}</h1>

      {/* {sessionFactors && (
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
      )} */}

      {challenge && <LoginPasskey challenge={challenge} />}
    </div>
  );
}

import { getSessionCookieById } from "@/lib/cookies";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import Alert from "@/ui/Alert";
import BackButton from "@/ui/BackButton";
import ChooseSecondFactorToSetup from "@/ui/ChooseSecondFactorToSetup";
import DynamicTheme from "@/ui/DynamicTheme";
import UserAvatar from "@/ui/UserAvatar";
import { Timestamp, timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";

function isSessionValid(session: Partial<Session>): {
  valid: boolean;
  verifiedAt?: Timestamp;
} {
  const validPassword = session?.factors?.password?.verifiedAt;
  const validPasskey = session?.factors?.webAuthN?.verifiedAt;
  const stillValid = session.expirationDate
    ? timestampDate(session.expirationDate) > new Date()
    : true;

  const verifiedAt = validPassword || validPasskey;
  const valid = !!((validPassword || validPasskey) && stillValid);

  return { valid, verifiedAt };
}

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const {
    loginName,
    checkAfter,
    force,
    authRequestId,
    organization,
    sessionId,
  } = searchParams;

  const sessionWithData = sessionId
    ? await loadSessionById(sessionId, organization)
    : await loadSessionByLoginname(loginName, organization);

  async function getAuthMethodsAndUser(session?: Session) {
    const userId = session?.factors?.user?.id;

    if (!userId) {
      throw Error("Could not get user id from session");
    }

    return listAuthenticationMethodTypes(userId).then((methods) => {
      return getUserByID(userId).then((user) => {
        const humanUser =
          user.user?.type.case === "human" ? user.user?.type.value : undefined;

        return {
          factors: session?.factors,
          authMethods: methods.authMethodTypes ?? [],
          phoneVerified: humanUser?.phone?.isVerified ?? false,
          emailVerified: humanUser?.email?.isVerified ?? false,
          expirationDate: session?.expirationDate,
        };
      });
    });
  }

  async function loadSessionByLoginname(
    loginName?: string,
    organization?: string,
  ) {
    return loadMostRecentSession({
      loginName,
      organization,
    }).then((session) => {
      return getAuthMethodsAndUser(session);
    });
  }

  async function loadSessionById(sessionId: string, organization?: string) {
    const recent = await getSessionCookieById({ sessionId, organization });
    return getSession({
      sessionId: recent.id,
      sessionToken: recent.token,
    }).then((sessionResponse) => {
      return getAuthMethodsAndUser(sessionResponse.session);
    });
  }

  const branding = await getBrandingSettings(organization);
  const loginSettings = await getLoginSettings(
    sessionWithData.factors?.user?.organizationId,
  );

  const { valid } = isSessionValid(sessionWithData);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Set up 2-Factor</h1>

        <p className="ztdl-p">Choose one of the following second factors.</p>

        {sessionWithData && (
          <UserAvatar
            loginName={loginName ?? sessionWithData.factors?.user?.loginName}
            displayName={sessionWithData.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {!(loginName || sessionId) && (
          <Alert>Provide your active session as loginName param</Alert>
        )}

        {!valid && (
          <Alert>
            You need to have a valid session in order to set a second factor!
            {/* TODO: show reauth button */}
          </Alert>
        )}

        {isSessionValid(sessionWithData).valid &&
          loginSettings &&
          sessionWithData && (
            <ChooseSecondFactorToSetup
              loginName={loginName}
              sessionId={sessionId}
              authRequestId={authRequestId}
              organization={organization}
              loginSettings={loginSettings}
              userMethods={sessionWithData.authMethods ?? []}
              phoneVerified={sessionWithData.phoneVerified ?? false}
              emailVerified={sessionWithData.emailVerified ?? false}
              checkAfter={checkAfter === "true"}
            ></ChooseSecondFactorToSetup>
          )}

        <div className="mt-8 flex w-full flex-row items-center">
          <BackButton />
          <span className="flex-grow"></span>
        </div>
      </div>
    </DynamicTheme>
  );
}

import { Alert } from "@/components/alert";
import { BackButton } from "@/components/back-button";
import { ChooseSecondFactorToSetup } from "@/components/choose-second-factor-to-setup";
import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getLoginSettings,
  getSession,
  getUserByID,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { Timestamp, timestampDate } from "@zitadel/client";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { getLocale, getTranslations } from "next-intl/server";

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

export default async function Page(
  props: {
    searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  }
) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "mfa" });
  const tError = await getTranslations({ locale, namespace: "error" });

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
        <h1>{t("set.title")}</h1>

        <p className="ztdl-p">{t("set.description")}</p>

        {sessionWithData && (
          <UserAvatar
            loginName={loginName ?? sessionWithData.factors?.user?.loginName}
            displayName={sessionWithData.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {!(loginName || sessionId) && <Alert>{tError("unknownContext")}</Alert>}

        {!valid && <Alert>{tError("sessionExpired")}</Alert>}

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

import { Alert } from "@/components/alert";
import { BackButton } from "@/components/back-button";
import { ChooseAuthenticatorToSetup } from "@/components/choose-authenticator-to-setup";
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
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page(
  props: {
    searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  }
) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "authenticator" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const { loginName, authRequestId, organization, sessionId } = searchParams;

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

  const branding = await getBrandingSettings(
    sessionWithData.factors?.user?.organizationId,
  );

  const loginSettings = await getLoginSettings(
    sessionWithData.factors?.user?.organizationId,
  );

  const params = new URLSearchParams({
    initial: "true", // defines that a code is not required and is therefore not shown in the UI
  });

  if (loginName) {
    params.set("loginName", loginName);
  }

  if (organization) {
    params.set("organization", organization);
  }

  if (authRequestId) {
    params.set("authRequestId", authRequestId);
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("title")}</h1>

        <p className="ztdl-p">{t("description")}</p>

        {sessionWithData && (
          <UserAvatar
            loginName={loginName ?? sessionWithData.factors?.user?.loginName}
            displayName={sessionWithData.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {!(loginName || sessionId) && <Alert>{tError("unknownContext")}</Alert>}

        {loginSettings && sessionWithData && (
          <ChooseAuthenticatorToSetup
            authMethods={sessionWithData.authMethods}
            loginSettings={loginSettings}
            params={params}
          ></ChooseAuthenticatorToSetup>
        )}

        <div className="mt-8 flex w-full flex-row items-center">
          <BackButton />
          <span className="flex-grow"></span>
        </div>
      </div>
    </DynamicTheme>
  );
}

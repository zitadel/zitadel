import { Alert } from "@/components/alert";
import { BackButton } from "@/components/back-button";
import { ChooseSecondFactor } from "@/components/choose-second-factor";
import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getSession,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "mfa" });

  const { loginName, authRequestId, organization, sessionId } = searchParams;

  const sessionFactors = sessionId
    ? await loadSessionById(sessionId, organization)
    : await loadSessionByLoginname(loginName, organization);

  async function loadSessionByLoginname(
    loginName?: string,
    organization?: string,
  ) {
    return loadMostRecentSession({
      loginName,
      organization,
    }).then((session) => {
      if (session && session.factors?.user?.id) {
        return listAuthenticationMethodTypes(session.factors.user.id).then(
          (methods) => {
            return {
              factors: session?.factors,
              authMethods: methods.authMethodTypes ?? [],
            };
          },
        );
      }
    });
  }

  async function loadSessionById(sessionId: string, organization?: string) {
    const recent = await getSessionCookieById({ sessionId, organization });
    return getSession({
      sessionId: recent.id,
      sessionToken: recent.token,
    }).then((response) => {
      if (response?.session && response.session.factors?.user?.id) {
        return listAuthenticationMethodTypes(
          response.session.factors.user.id,
        ).then((methods) => {
          return {
            factors: response.session?.factors,
            authMethods: methods.authMethodTypes ?? [],
          };
        });
      }
    });
  }

  const branding = await getBrandingSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("verify.title")}</h1>

        <p className="ztdl-p">{t("verify.description")}</p>

        {sessionFactors && (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {!(loginName || sessionId) && (
          <Alert>{t("error:unknownContext")}</Alert>
        )}

        {sessionFactors ? (
          <ChooseSecondFactor
            loginName={loginName}
            sessionId={sessionId}
            authRequestId={authRequestId}
            organization={organization}
            userMethods={sessionFactors.authMethods ?? []}
          ></ChooseSecondFactor>
        ) : (
          <Alert>{t("verify.noResults")}</Alert>
        )}

        <div className="mt-8 flex w-full flex-row items-center">
          <BackButton />
          <span className="flex-grow"></span>
        </div>
      </div>
    </DynamicTheme>
  );
}

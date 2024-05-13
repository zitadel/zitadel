import {
  getBrandingSettings,
  getSession,
  listAuthenticationMethodTypes,
  server,
} from "@/lib/zitadel";
import Alert from "@/ui/Alert";
import ChooseSecondFactor from "@/ui/ChooseSecondFactor";
import DynamicTheme from "@/ui/DynamicTheme";
import UserAvatar from "@/ui/UserAvatar";
import {
  getMostRecentCookieWithLoginname,
  getSessionCookieById,
} from "@/utils/cookies";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, checkAfter, authRequestId, organization, sessionId } =
    searchParams;

  const sessionFactors = sessionId
    ? await loadSessionById(sessionId, organization)
    : await loadSessionByLoginname(loginName, organization);

  async function loadSessionByLoginname(
    loginName?: string,
    organization?: string,
  ) {
    const recent = await getMostRecentCookieWithLoginname(
      loginName,
      organization,
    );
    return getSession(server, recent.id, recent.token).then((response) => {
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

  async function loadSessionById(sessionId: string, organization?: string) {
    const recent = await getSessionCookieById(sessionId, organization);
    return getSession(server, recent.id, recent.token).then((response) => {
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

  const branding = await getBrandingSettings(server, organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Verify 2-Factor</h1>

        <p className="ztdl-p">Choose one of the following second factors.</p>

        {sessionFactors && (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {!(loginName || sessionId) && (
          <Alert>Provide your active session as loginName param</Alert>
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
          <Alert>No second factors available to setup.</Alert>
        )}
      </div>
    </DynamicTheme>
  );
}

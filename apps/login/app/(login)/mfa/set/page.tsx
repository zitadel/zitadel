import {
  getBrandingSettings,
  getLoginSettings,
  getSession,
  server,
} from "#/lib/zitadel";
import Alert from "#/ui/Alert";
import ChooseSecondFactorToSetup from "#/ui/ChooseSecondFactorToSetup";
import DynamicTheme from "#/ui/DynamicTheme";
import UserAvatar from "#/ui/UserAvatar";
import {
  getMostRecentCookieWithLoginname,
  getSessionCookieById,
} from "#/utils/cookies";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, altPassword, authRequestId, organization, sessionId } =
    searchParams;

  const sessionFactors = sessionId
    ? await loadSessionById(sessionId, organization)
    : await loadSessionByLoginname(loginName, organization);

  async function loadSessionByLoginname(
    loginName?: string,
    organization?: string
  ) {
    const recent = await getMostRecentCookieWithLoginname(
      loginName,
      organization
    );
    return getSession(server, recent.id, recent.token).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }

  async function loadSessionById(sessionId: string, organization?: string) {
    const recent = await getSessionCookieById(sessionId, organization);
    return getSession(server, recent.id, recent.token).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }

  const branding = await getBrandingSettings(server, organization);
  const loginSettings = await getLoginSettings(server, organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Set up 2-Factor</h1>

        <p className="ztdl-p">Choose one of the following second factors.</p>

        {sessionFactors && (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
          ></UserAvatar>
        )}

        {!sessionFactors && <div className="py-4"></div>}

        {!(loginName || sessionId) && (
          <Alert>Provide your active session as loginName param</Alert>
        )}

        {loginSettings ? (
          <ChooseSecondFactorToSetup
            loginName={loginName}
            sessionId={sessionId}
            authRequestId={authRequestId}
            organization={organization}
            loginSettings={loginSettings}
          ></ChooseSecondFactorToSetup>
        ) : (
          <Alert>No second factors available to setup.</Alert>
        )}
      </div>
    </DynamicTheme>
  );
}

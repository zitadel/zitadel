import { getBrandingSettings, getSession, server } from "@/lib/zitadel";
import Alert from "@/ui/Alert";
import DynamicTheme from "@/ui/DynamicTheme";
import LoginPasskey from "@/ui/LoginPasskey";
import UserAvatar from "@/ui/UserAvatar";
import {
  getMostRecentCookieWithLoginname,
  getSessionCookieById,
} from "@/utils/cookies";

const title = "Authenticate with a passkey";
const description =
  "Your device will ask for your fingerprint, face, or screen lock";

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
    organization?: string,
  ) {
    const recent = await getMostRecentCookieWithLoginname(
      loginName,
      organization,
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

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{title}</h1>

        {sessionFactors && (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}
        <p className="ztdl-p mb-6 block">{description}</p>

        {!(loginName || sessionId) && (
          <Alert>Provide your active session as loginName param</Alert>
        )}

        {(loginName || sessionId) && (
          <LoginPasskey
            loginName={loginName}
            sessionId={sessionId}
            authRequestId={authRequestId}
            altPassword={altPassword === "true"}
            organization={organization}
          />
        )}
      </div>
    </DynamicTheme>
  );
}

import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { LoginPasskey } from "@/components/login-passkey";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getSession } from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "passkey" });

  const { loginName, altPassword, authRequestId, organization, sessionId } =
    searchParams;

  const sessionFactors = sessionId
    ? await loadSessionById(sessionId, organization)
    : await loadMostRecentSession({ loginName, organization });

  async function loadSessionById(sessionId: string, organization?: string) {
    const recent = await getSessionCookieById({ sessionId, organization });
    return getSession({
      sessionId: recent.id,
      sessionToken: recent.token,
    }).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }

  const branding = await getBrandingSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("verify.title")}</h1>

        {sessionFactors && (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}
        <p className="ztdl-p mb-6 block">{t("verify.description")}</p>

        {!(loginName || sessionId) && (
          <Alert>{t("error:unknownContext")}</Alert>
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

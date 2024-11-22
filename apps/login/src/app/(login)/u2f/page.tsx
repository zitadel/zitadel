import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { LoginPasskey } from "@/components/login-passkey";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getSession } from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page(
  props: {
    searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  }
) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "u2f" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const { loginName, authRequestId, sessionId, organization } = searchParams;

  const branding = await getBrandingSettings(organization);

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

        {!(loginName || sessionId) && <Alert>{tError("unknownContext")}</Alert>}

        {(loginName || sessionId) && (
          <LoginPasskey
            loginName={loginName}
            sessionId={sessionId}
            authRequestId={authRequestId}
            altPassword={false}
            organization={organization}
            login={false} // this sets the userVerificationRequirement to discouraged as its used as second factor
          />
        )}
      </div>
    </DynamicTheme>
  );
}

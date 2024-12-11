import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { LoginOTP } from "@/components/login-otp";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getLoginSettings,
  getSession,
} from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  params: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const params = await props.params;
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "otp" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const {
    loginName, // send from password page
    userId, // send from email link
    authRequestId,
    sessionId,
    organization,
    code,
    submit,
  } = searchParams;

  const { method } = params;

  const session = sessionId
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

  // email links do not come with organization, thus we need to use the session's organization
  const branding = await getBrandingSettings(
    organization ?? session?.factors?.user?.organizationId,
  );

  const loginSettings = await getLoginSettings(
    organization ?? session?.factors?.user?.organizationId,
  );

  const host = (await headers()).get("host");

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("verify.title")}</h1>
        {method === "time-based" && (
          <p className="ztdl-p">{t("verify.totpDescription")}</p>
        )}
        {method === "sms" && (
          <p className="ztdl-p">{t("verify.smsDescription")}</p>
        )}
        {method === "email" && (
          <p className="ztdl-p">{t("verify.emailDescription")}</p>
        )}

        {!session && (
          <div className="py-4">
            <Alert>{tError("unknownContext")}</Alert>
          </div>
        )}

        {session && (
          <UserAvatar
            loginName={loginName ?? session.factors?.user?.loginName}
            displayName={session.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {method && session && (
          <LoginOTP
            loginName={loginName ?? session.factors?.user?.loginName}
            sessionId={sessionId}
            authRequestId={authRequestId}
            organization={
              organization ?? session?.factors?.user?.organizationId
            }
            method={method}
            loginSettings={loginSettings}
            host={host}
            code={code}
          ></LoginOTP>
        )}
      </div>
    </DynamicTheme>
  );
}

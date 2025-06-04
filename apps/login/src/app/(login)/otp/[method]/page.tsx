import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { LoginOTP } from "@/components/login-otp";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
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

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);
  const host = _headers.get("host");

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

  const {
    loginName, // send from password page
    userId, // send from email link
    requestId,
    sessionId,
    organization,
    code,
    submit,
  } = searchParams;

  const { method } = params;

  const session = sessionId
    ? await loadSessionById(serviceUrl, sessionId, organization)
    : await loadMostRecentSession({
        serviceUrl,
        sessionParams: { loginName, organization },
      });

  async function loadSessionById(
    host: string,
    sessionId: string,
    organization?: string,
  ) {
    const recent = await getSessionCookieById({ sessionId, organization });
    return getSession({
      serviceUrl,
      sessionId: recent.id,
      sessionToken: recent.token,
    }).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }

  // email links do not come with organization, thus we need to use the session's organization
  const branding = await getBrandingSettings({
    serviceUrl,
    organization: organization ?? session?.factors?.user?.organizationId,
  });

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: organization ?? session?.factors?.user?.organizationId,
  });

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
            requestId={requestId}
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

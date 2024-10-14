import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { LoginOTP } from "@/components/login-otp";
import { UserAvatar } from "@/components/user-avatar";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings } from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({
  searchParams,
  params,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
  params: Record<string | number | symbol, string | undefined>;
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "otp" });

  const { loginName, authRequestId, sessionId, organization, code, submit } =
    searchParams;

  const { method } = params;

  const session = await loadMostRecentSession({
    loginName,
    organization,
  });

  const branding = await getBrandingSettings(organization);

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
            <Alert>{t("error:unknownContext")}</Alert>
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

        {method && (
          <LoginOTP
            loginName={loginName}
            sessionId={sessionId}
            authRequestId={authRequestId}
            organization={organization}
            method={method}
          ></LoginOTP>
        )}
      </div>
    </DynamicTheme>
  );
}

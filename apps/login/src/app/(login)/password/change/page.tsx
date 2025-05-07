import { Alert } from "@/components/alert";
import { ChangePasswordForm } from "@/components/change-password-form";
import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getLoginSettings,
  getPasswordComplexitySettings,
} from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "password" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const { loginName, organization, requestId } = searchParams;

  // also allow no session to be found (ignoreUnkownUsername)
  const sessionFactors = await loadMostRecentSession({
    serviceUrl,
    sessionParams: {
      loginName,
      organization,
    },
  });

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  const passwordComplexity = await getPasswordComplexitySettings({
    serviceUrl,
    organization: sessionFactors?.factors?.user?.organizationId,
  });

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: sessionFactors?.factors?.user?.organizationId,
  });

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>
          {sessionFactors?.factors?.user?.displayName ?? t("change.title")}
        </h1>
        <p className="ztdl-p mb-6 block">{t("change.description")}</p>

        {/* show error only if usernames should be shown to be unknown */}
        {(!sessionFactors || !loginName) &&
          !loginSettings?.ignoreUnknownUsernames && (
            <div className="py-4">
              <Alert>{tError("unknownContext")}</Alert>
            </div>
          )}

        {sessionFactors && (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        )}

        {passwordComplexity &&
        loginName &&
        sessionFactors?.factors?.user?.id ? (
          <ChangePasswordForm
            sessionId={sessionFactors.id}
            loginName={loginName}
            requestId={requestId}
            organization={organization}
            passwordComplexitySettings={passwordComplexity}
          />
        ) : (
          <div className="py-4">
            <Alert>{tError("failedLoading")}</Alert>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}

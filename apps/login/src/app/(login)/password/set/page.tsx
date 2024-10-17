import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SetPasswordForm } from "@/components/set-password-form";
import { UserAvatar } from "@/components/user-avatar";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getLoginSettings,
  getPasswordComplexitySettings,
} from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "password" });

  const { loginName, organization, authRequestId, code } = searchParams;

  // also allow no session to be found (ignoreUnkownUsername)
  const sessionFactors = await loadMostRecentSession({
    loginName,
    organization,
  });

  const branding = await getBrandingSettings(organization);

  const passwordComplexity = await getPasswordComplexitySettings(
    sessionFactors?.factors?.user?.organizationId,
  );

  const loginSettings = await getLoginSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{sessionFactors?.factors?.user?.displayName ?? t("set.title")}</h1>
        <p className="ztdl-p mb-6 block">{t("set.description")}</p>

        {/* show error only if usernames should be shown to be unknown */}
        {(!sessionFactors || !loginName) &&
          !loginSettings?.ignoreUnknownUsernames && (
            <div className="py-4">
              <Alert>{t("error:unknownContext")}</Alert>
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

        <Alert type={AlertType.INFO}>{t("set.codeSent")}</Alert>

        {passwordComplexity &&
        loginName &&
        sessionFactors?.factors?.user?.id ? (
          <SetPasswordForm
            code={code}
            userId={sessionFactors.factors.user.id}
            loginName={loginName}
            authRequestId={authRequestId}
            organization={organization}
            passwordComplexitySettings={passwordComplexity}
          />
        ) : (
          <div className="py-4">
            <Alert>{t("error:failedLoading")}</Alert>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}

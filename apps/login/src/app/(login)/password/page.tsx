import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { PasswordForm } from "@/components/password-form";
import { UserAvatar } from "@/components/user-avatar";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getLoginSettings } from "@/lib/zitadel";
import { PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "password" });

  const { loginName, organization, authRequestId, alt } = searchParams;

  // also allow no session to be found (ignoreUnkownUsername)
  let sessionFactors;
  try {
    sessionFactors = await loadMostRecentSession({
      loginName,
      organization,
    });
  } catch (error) {
    // ignore error to continue to show the password form
    console.warn(error);
  }

  const branding = await getBrandingSettings(organization);
  const loginSettings = await getLoginSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{sessionFactors?.factors?.user?.displayName ?? t("title")}</h1>
        <p className="ztdl-p mb-6 block">{t("description")}</p>

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

        {loginName && (
          <PasswordForm
            loginName={loginName}
            authRequestId={authRequestId}
            organization={organization}
            loginSettings={loginSettings}
            promptPasswordless={
              loginSettings?.passkeysType === PasskeysType.ALLOWED
            }
            isAlternative={alt === "true"}
          />
        )}
      </div>
    </DynamicTheme>
  );
}

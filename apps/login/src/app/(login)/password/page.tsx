import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getLoginSettings } from "@/lib/zitadel";
import Alert from "@/ui/Alert";
import DynamicTheme from "@/ui/DynamicTheme";
import PasswordForm from "@/ui/PasswordForm";
import UserAvatar from "@/ui/UserAvatar";
import { PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, organization, authRequestId, alt } = searchParams;

  const sessionFactors = await loadMostRecentSession({
    loginName,
    organization,
  });

  const branding = await getBrandingSettings(organization);
  const loginSettings = await getLoginSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{sessionFactors?.factors?.user?.displayName ?? "Password"}</h1>
        <p className="ztdl-p mb-6 block">Enter your password.</p>

        {(!sessionFactors || !loginName) && (
          <div className="py-4">
            <Alert>
              Could not get the context of the user. Make sure to enter the
              username first or provide a loginName as searchParam.
            </Alert>
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

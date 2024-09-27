import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { RegisterU2f } from "@/components/register-u2f";
import { UserAvatar } from "@/components/user-avatar";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings } from "@/lib/zitadel";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { loginName, organization, authRequestId, checkAfter } = searchParams;

  const sessionFactors = await loadMostRecentSession({
    loginName,
    organization,
  });

  const title = "Use your passkey to confirm it's really you";
  const description =
    "Your device will ask for your fingerprint, face, or screen lock";

  const branding = await getBrandingSettings(organization);

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

        {!sessionFactors && (
          <div className="py-4">
            <Alert>
              Could not get the context of the user. Make sure to enter the
              username first or provide a loginName as searchParam.
            </Alert>
          </div>
        )}

        {sessionFactors?.id && (
          <RegisterU2f
            loginName={loginName}
            sessionId={sessionFactors.id}
            organization={organization}
            authRequestId={authRequestId}
            checkAfter={checkAfter === "true"}
          />
        )}
      </div>
    </DynamicTheme>
  );
}

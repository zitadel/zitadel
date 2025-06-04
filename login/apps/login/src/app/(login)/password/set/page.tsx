import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SetPasswordForm } from "@/components/set-password-form";
import { UserAvatar } from "@/components/user-avatar";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getLoginSettings,
  getPasswordComplexitySettings,
  getUserByID,
} from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "password" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const { userId, loginName, organization, requestId, code, initial } =
    searchParams;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  // also allow no session to be found (ignoreUnkownUsername)
  let session: Session | undefined;
  if (loginName) {
    session = await loadMostRecentSession({
      serviceUrl,
      sessionParams: {
        loginName,
        organization,
      },
    });
  }

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  const passwordComplexity = await getPasswordComplexitySettings({
    serviceUrl,
    organization: session?.factors?.user?.organizationId,
  });

  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization,
  });

  let user: User | undefined;
  let displayName: string | undefined;
  if (userId) {
    const userResponse = await getUserByID({
      serviceUrl,
      userId,
    });
    user = userResponse.user;

    if (user?.type.case === "human") {
      displayName = (user.type.value as HumanUser).profile?.displayName;
    }
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{session?.factors?.user?.displayName ?? t("set.title")}</h1>
        <p className="ztdl-p mb-6 block">{t("set.description")}</p>

        {/* show error only if usernames should be shown to be unknown */}
        {loginName && !session && !loginSettings?.ignoreUnknownUsernames && (
          <div className="py-4">
            <Alert>{tError("unknownContext")}</Alert>
          </div>
        )}

        {session ? (
          <UserAvatar
            loginName={loginName ?? session.factors?.user?.loginName}
            displayName={session.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        ) : user ? (
          <UserAvatar
            loginName={user?.preferredLoginName}
            displayName={displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        ) : null}

        {!initial && <Alert type={AlertType.INFO}>{t("set.codeSent")}</Alert>}

        {passwordComplexity &&
        (loginName ?? user?.preferredLoginName) &&
        (userId ?? session?.factors?.user?.id) ? (
          <SetPasswordForm
            code={code}
            userId={userId ?? (session?.factors?.user?.id as string)}
            loginName={loginName ?? (user?.preferredLoginName as string)}
            requestId={requestId}
            organization={organization}
            passwordComplexitySettings={passwordComplexity}
            codeRequired={!(initial === "true")}
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

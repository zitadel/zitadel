import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { SetPasswordForm } from "@/components/set-password-form";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getServiceConfig } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getLoginSettings,
  getPasswordComplexitySettings,
  getUserByID,
  searchUsers,
} from "@/lib/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { create } from "@zitadel/client";
import { LoginSettingsSchema } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("password");
  return { title: t("set.title") };
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  let { userId, loginName, organization, requestId, code, initial } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  // also allow no session to be found (ignoreUnkownUsername)
  let session: Session | undefined;
  if (loginName) {
    session = await loadMostRecentSession({
      serviceConfig,
      sessionParams: {
        loginName,
        organization,
      },
    });
  }

  const branding = await getBrandingSettings({ serviceConfig, organization });

  const passwordComplexity = await getPasswordComplexitySettings({
    serviceConfig,
    organization: organization ?? session?.factors?.user?.organizationId,
  });

  const loginSettings = await getLoginSettings({ serviceConfig, organization });

  let user: User | undefined;
  let displayName: string | undefined;
  if (userId) {
    const userResponse = await getUserByID({ serviceConfig, userId });
    user = userResponse.user;

    if (user?.type.case === "human") {
      displayName = (user.type.value as HumanUser).profile?.displayName;
    }
  } else if (loginName) {
    const users = await searchUsers({
      serviceConfig,
      searchValue: loginName,
      loginSettings:
        loginSettings ??
        create(LoginSettingsSchema, {
          ignoreUnknownUsernames: true,
          disableLoginWithEmail: false,
          disableLoginWithPhone: false,
          allowUsernamePassword: true,
          allowRegister: true,
          allowExternalIdp: true,
          forceMfa: false,
          hidePasswordReset: false,
        }), // Default settings if null
      organizationId: organization,
    });

    if (users.result && users.result.length === 1) {
      const foundUser = users.result[0];
      userId = foundUser.userId;
      user = foundUser;
      if (user.type.case === "human") {
        displayName = (user.type.value as HumanUser).profile?.displayName;
      }
    } else if (loginSettings?.ignoreUnknownUsernames) {
      // Prevent enumeration by pretending we found a user
      userId = "000000000000000000";
    }
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>{session?.factors?.user?.displayName ?? <Translated i18nKey="set.title" namespace="password" />}</h1>
        <p className="ztdl-p mb-6 block">
          <Translated i18nKey="set.description" namespace="password" />
        </p>

        {/* show error only if usernames should be shown to be unknown */}
        {loginName && !session && !loginSettings?.ignoreUnknownUsernames && (
          <div className="py-4">
            <Alert>
              <Translated i18nKey="unknownContext" namespace="error" />
            </Alert>
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
      </div>

      <div className="w-full">
        {!initial && (
          <Alert type={AlertType.INFO}>
            <Translated i18nKey="set.codeSent" namespace="password" />
          </Alert>
        )}

        {passwordComplexity &&
        (loginName ?? user?.preferredLoginName) &&
        (userId ?? session?.factors?.user?.id ?? (loginSettings?.ignoreUnknownUsernames ? "unknown" : undefined)) ? (
          <SetPasswordForm
            code={code}
            userId={userId ?? (session?.factors?.user?.id as string) ?? "unknown"}
            loginName={loginName ?? (user?.preferredLoginName as string)}
            requestId={requestId}
            organization={organization}
            passwordComplexitySettings={passwordComplexity}
            codeRequired={!(initial === "true")}
          />
        ) : (
          <div className="py-4">
            <Alert>
              <Translated i18nKey="failedLoading" namespace="error" />
            </Alert>
          </div>
        )}
      </div>
    </DynamicTheme>
  );
}

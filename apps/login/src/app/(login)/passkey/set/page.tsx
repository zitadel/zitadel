import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { RegisterPasskey } from "@/components/register-passkey";
import { RestartSignIn } from "@/components/restart-sign-in";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getSetupAuthorizationError, setupErrorTranslationKey } from "@/lib/server/enrollment-guard";
import { getServiceConfig } from "@/lib/server/service-url";
import { loadMostRecentSession } from "@/lib/server/session";
import { getBrandingSettings, getUserByID } from "@/lib/server/zitadel";
import { Session } from "@zitadel/proto/zitadel/session/v2/session_pb";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("passkey");
  return { title: t("set.title") };
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  const { userId, loginName, prompt, organization, requestId, code, codeId } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  // also allow no session to be found for userId-based flows
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

  // Defense in depth alongside registerPasskeyLink: don't offer the registration form for a
  // session that may not set up a passkey (unauthorized session or passkeys not allowed).
  // In the emailed code flow (userId + code) the code itself authorizes the registration, and
  // RegisterPasskey uses it even if a session is in the browser; otherwise registerPasskeyLink
  // applies the same check.
  const isCodeFlow = !!(userId && code && codeId);
  let setupError: string | null = null;
  if (session?.factors?.user?.id && !isCodeFlow) {
    setupError = await getSetupAuthorizationError({
      serviceConfig,
      session,
      userId: session.factors.user.id,
      method: "passkey",
    });
  }

  let user: User | undefined;
  let displayName: string | undefined;
  if (userId) {
    const userResponse = await getUserByID({ serviceConfig, userId });
    user = userResponse.user;

    if (user?.type.case === "human") {
      displayName = (user.type.value as HumanUser).profile?.displayName;
    }
  }

  // in the code flow the passkey is registered for the code's user, so show that user even if a
  // session (possibly of another user) is in the browser
  const showCodeUser = isCodeFlow && !!user;

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="set.title" namespace="passkey" />
        </h1>

        <p className="ztdl-p mb-6 block">
          <Translated i18nKey="set.description" namespace="passkey" />
        </p>

        {session && !showCodeUser ? (
          <UserAvatar
            loginName={loginName ?? session.factors?.user?.loginName}
            organizationId={session.factors?.user?.organizationId}
            hideLoginNameSuffix={branding?.hideLoginNameSuffix}
            displayName={session.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        ) : user ? (
          <UserAvatar
            loginName={user?.preferredLoginName}
            organizationId={user?.details?.resourceOwner}
            hideLoginNameSuffix={branding?.hideLoginNameSuffix}
            displayName={displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        ) : null}
      </div>

      <div className="w-full">
        <Alert type={AlertType.INFO}>
          <span>
            <Translated i18nKey="set.info.description" namespace="passkey" />
            <a
              className="text-primary-light-500 hover:text-primary-light-300 dark:text-primary-dark-500 hover:dark:text-primary-dark-300"
              target="_blank"
              href="https://zitadel.com/docs/guides/manage/user/reg-create-user#with-passwordless"
            >
              <Translated i18nKey="set.info.link" namespace="passkey" />
            </a>
          </span>
        </Alert>

        {!session && !user && (
          <div className="py-4">
            <Alert>
              {/* context was provided but the session could not be resolved (cookie missing, invalid or expired) */}
              <Translated i18nKey={loginName || userId ? "sessionExpired" : "unknownContext"} namespace="error" />
            </Alert>
            <RestartSignIn loginName={loginName} organization={organization} requestId={requestId} />
          </div>
        )}

        {setupError && (
          <div className="py-4">
            <Alert>
              <Translated i18nKey={setupErrorTranslationKey(setupError)} namespace="error" />
            </Alert>
          </div>
        )}

        {!setupError && (session?.id || userId) && (
          <RegisterPasskey
            sessionId={session?.id}
            userId={userId}
            isPrompt={!!prompt}
            organization={organization}
            requestId={requestId}
            code={code}
            codeId={codeId}
            loginName={
              showCodeUser
                ? user?.preferredLoginName
                : (loginName ?? session?.factors?.user?.loginName ?? user?.preferredLoginName)
            }
          />
        )}
      </div>
    </DynamicTheme>
  );
}

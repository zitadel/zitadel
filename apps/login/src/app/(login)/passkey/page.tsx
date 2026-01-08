import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { LoginPasskey } from "@/components/login-passkey";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { getServiceConfig } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getDefaultOrg, getLoginSettings, getSession, searchUsers } from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("passkey");
  return { title: t("verify.title") };
}

export default async function Page(props: { searchParams: Promise<Record<string | number | symbol, string | undefined>> }) {
  const searchParams = await props.searchParams;

  const { loginName, altPassword, requestId, organization, sessionId } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  let defaultOrganization;
  if (!organization) {
    const org: Organization | null = await getDefaultOrg({ serviceConfig });

    if (org) {
      defaultOrganization = org.id;
    }
  }

  let sessionFactors = sessionId ? await loadSessionById(sessionId, organization) : undefined;

  if (!sessionFactors && !sessionId) {
    sessionFactors = await loadMostRecentSession({
      serviceConfig,
      sessionParams: { loginName, organization },
    }).catch(() => {
      // ignore error
      return undefined;
    });
  }

  async function loadSessionById(sessionId: string, organization?: string) {
    const recent = await getSessionCookieById({ sessionId, organization });

    if (!recent) {
      return undefined;
    }

    return getSession({ serviceConfig, sessionId: recent.id, sessionToken: recent.token }).then((response) => {
      if (response?.session) {
        return response.session;
      }
    });
  }

  const branding = await getBrandingSettings({
    serviceConfig,
    organization: organization ?? sessionFactors?.factors?.user?.organizationId ?? defaultOrganization,
  });

  let user: User | undefined;
  let human: HumanUser | undefined;

  let loginSettings;
  if (!sessionFactors && loginName) {
    loginSettings = await getLoginSettings({ serviceConfig, organization: organization ?? defaultOrganization });

    if (loginSettings) {
      const users = await searchUsers({
        serviceConfig,
        searchValue: loginName,
        loginSettings: loginSettings,
        organizationId: organization,
      });

      if (users.result && users.result.length === 1) {
        const foundUser = users.result[0];
        user = foundUser;
        if (user.type.case === "human") {
          human = user.type.value as HumanUser;
        }
      }
    }
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="verify.title" namespace="passkey" />
        </h1>

        <p className="ztdl-p mb-6 block">
          <Translated i18nKey="verify.description" namespace="passkey" />
        </p>

        {sessionFactors ? (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        ) : (
          (user || loginName) && (
            <UserAvatar
              loginName={loginName ?? user?.preferredLoginName}
              displayName={
                !loginSettings?.ignoreUnknownUsernames
                  ? human?.profile?.displayName
                  : (loginName ?? user?.preferredLoginName)
              }
              showDropdown={false}
            />
          )
        )}
      </div>

      <div className="w-full">
        {!(loginName || sessionId) && (
          <Alert>
            <Translated i18nKey="unknownContext" namespace="error" />
          </Alert>
        )}

        {(loginName || sessionId) && (
          <LoginPasskey
            loginName={loginName}
            sessionId={sessionId}
            requestId={requestId}
            altPassword={altPassword === "true"}
            organization={organization}
          />
        )}
      </div>
    </DynamicTheme>
  );
}

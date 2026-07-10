import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { VerifyForm } from "@/components/verify-form";
import { UNKNOWN_USER_ID } from "@/lib/constants";
import { getServiceConfig } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getLoginSettings, getUserByID, searchUsers } from "@/lib/zitadel";
import { LoginSettings } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("verify");
  return { title: t("verify.title") };
}

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;

  const { userId, loginName, code, organization, requestId, invite } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const branding = await getBrandingSettings({ serviceConfig, organization });

  let sessionFactors;
  let user: User | undefined;
  let human: HumanUser | undefined;
  let id: string | undefined;
  let loginSettings: LoginSettings | undefined;

  const autoSubmitCode = process.env.NEXT_PUBLIC_AUTO_SUBMIT_CODE === "true";

  if ("loginName" in searchParams) {
    sessionFactors = await loadMostRecentSession({
      serviceConfig,
      sessionParams: {
        loginName,
        organization,
      },
    }).catch(async (error) => {
      loginSettings = await getLoginSettings({ serviceConfig, organization });
      if (!loginSettings?.ignoreUnknownUsernames) {
        console.error("loadMostRecentSession failed", error);
      }
      // ignore error, as we might not have a session yet
      return undefined;
    });
  } else if ("userId" in searchParams && userId) {
    const userResponse = await getUserByID({ serviceConfig, userId });
    if (userResponse) {
      user = userResponse.user;
      if (user?.type.case === "human") {
        human = user.type.value as HumanUser;
      }
    }
  }

  id = userId ?? sessionFactors?.factors?.user?.id;

  if (!id && loginName) {
    if (!loginSettings) {
      loginSettings = await getLoginSettings({ serviceConfig, organization });
    }

    if (!loginSettings) {
      console.error("loginSettings not found");
      return;
    }

    const users = await searchUsers({
      serviceConfig,
      searchValue: loginName,
      loginSettings: loginSettings,
      organizationId: organization,
    });

    if (users.result && users.result.length === 1) {
      const foundUser = users.result[0];
      id = foundUser.userId;
      user = foundUser;
      if (user.type.case === "human") {
        human = user.type.value as HumanUser;
      }
    } else if (loginSettings?.ignoreUnknownUsernames) {
      // Prevent enumeration by pretending we found a user
      id = UNKNOWN_USER_ID;
    }
  }

  const params = new URLSearchParams({
    userId: userId,
    initial: "true", // defines that a code is not required and is therefore not shown in the UI
  });

  if (loginName) {
    params.set("loginName", loginName);
  }

  if (organization) {
    params.set("organization", organization);
  }

  if (requestId) {
    params.set("requestId", requestId);
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>
          <Translated i18nKey="verify.title" namespace="verify" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="verify.description" namespace="verify" />
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
        {!id && (
          <div className="py-4">
            <Alert>
              <Translated i18nKey="unknownContext" namespace="error" />
            </Alert>
          </div>
        )}

        {id && (
          <VerifyForm
            loginName={loginName}
            organization={organization}
            userId={id}
            code={code}
            isInvite={invite === "true"}
            requestId={requestId}
            submit={autoSubmitCode}
          />
        )}
      </div>
    </DynamicTheme>
  );
}

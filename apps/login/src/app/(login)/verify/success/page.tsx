import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { getSessionCookieById } from "@/lib/cookies";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getLoginSettings,
  getSession,
  getUserByID,
} from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

async function loadSessionById(
  serviceUrl: string,
  sessionId: string,
  organization?: string,
) {
  const recent = await getSessionCookieById({ sessionId, organization });
  return getSession({
    serviceUrl,
    sessionId: recent.id,
    sessionToken: recent.token,
  }).then((response) => {
    if (response?.session) {
      return response.session;
    }
  });
}

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "verify" });

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const { loginName, requestId, organization, userId } = searchParams;

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  const sessionFactors = await loadMostRecentSession({
    serviceUrl,
    sessionParams: { loginName, organization },
  }).catch((error) => {
    console.warn("Error loading session:", error);
  });

  let loginSettings;
  if (!requestId) {
    loginSettings = await getLoginSettings({
      serviceUrl,
      organization,
    });
  }

  const id = userId ?? sessionFactors?.factors?.user?.id;

  if (!id) {
    throw Error("Failed to get user id");
  }

  const userResponse = await getUserByID({
    serviceUrl,
    userId: id,
  });

  let user: User | undefined;
  let human: HumanUser | undefined;

  if (userResponse) {
    user = userResponse.user;
    if (user?.type.case === "human") {
      human = user.type.value as HumanUser;
    }
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("successTitle")}</h1>
        <p className="ztdl-p mb-6 block">{t("successDescription")}</p>

        {sessionFactors ? (
          <UserAvatar
            loginName={loginName ?? sessionFactors.factors?.user?.loginName}
            displayName={sessionFactors.factors?.user?.displayName}
            showDropdown
            searchParams={searchParams}
          ></UserAvatar>
        ) : (
          user && (
            <UserAvatar
              loginName={user.preferredLoginName}
              displayName={human?.profile?.displayName}
              showDropdown={false}
            />
          )
        )}
      </div>
    </DynamicTheme>
  );
}

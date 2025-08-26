import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { VerifyForm } from "@/components/verify-form";
import { sendEmailCode, sendInviteEmailCode } from "@/lib/server/verify";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getUserByID } from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("verify");
  return { title: t('verify.title')};
}

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;

  const { userId, loginName, code, organization, requestId, invite, send } =
    searchParams;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const branding = await getBrandingSettings({
    serviceUrl,
    organization,
  });

  let sessionFactors;
  let user: User | undefined;
  let human: HumanUser | undefined;
  let id: string | undefined;

  const doSend = send === "true";

  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  async function sendEmail(userId: string) {
    const host = _headers.get("host");

    if (!host || typeof host !== "string") {
      throw new Error("No host found");
    }

    if (invite === "true") {
      await sendInviteEmailCode({
        userId,
        urlTemplate:
          `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&invite=true` +
          (requestId ? `&requestId=${requestId}` : ""),
      }).catch((error) => {
        console.error("Could not send invitation email", error);
        throw Error("Failed to send invitation email");
      });
    } else {
      await sendEmailCode({
        userId,
        urlTemplate:
          `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
          (requestId ? `&requestId=${requestId}` : ""),
      }).catch((error) => {
        console.error("Could not send verification email", error);
        throw Error("Failed to send verification email");
      });
    }
  }

  if ("loginName" in searchParams) {
    sessionFactors = await loadMostRecentSession({
      serviceUrl,
      sessionParams: {
        loginName,
        organization,
      },
    });

    if (doSend && sessionFactors?.factors?.user?.id) {
      await sendEmail(sessionFactors.factors.user.id);
    }
  } else if ("userId" in searchParams && userId) {
    if (doSend) {
      await sendEmail(userId);
    }

    const userResponse = await getUserByID({
      serviceUrl,
      userId,
    });
    if (userResponse) {
      user = userResponse.user;
      if (user?.type.case === "human") {
        human = user.type.value as HumanUser;
      }
    }
  }

  id = userId ?? sessionFactors?.factors?.user?.id;

  if (!id) {
    throw Error("Failed to get user id");
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
      <div className="flex flex-col items-center space-y-4">
        <h1>
          <Translated i18nKey="verify.title" namespace="verify" />
        </h1>
        <p className="ztdl-p mb-6 block">
          <Translated i18nKey="verify.description" namespace="verify" />
        </p>

        {!id && (
          <div className="py-4">
            <Alert>
              <Translated i18nKey="unknownContext" namespace="error" />
            </Alert>
          </div>
        )}

        {id && send && (
          <div className="w-full py-4">
            <Alert type={AlertType.INFO}>
              <Translated i18nKey="verify.codeSent" namespace="verify" />
            </Alert>
          </div>
        )}

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

        <VerifyForm
          loginName={loginName}
          organization={organization}
          userId={id}
          code={code}
          isInvite={invite === "true"}
          requestId={requestId}
        />
      </div>
    </DynamicTheme>
  );
}

import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { VerifyEmailForm } from "@/components/verify-email-form";
import { getPublicHostWithProtocol } from "@/lib/server/host";
import { sendEmailCode, sendInviteEmailCode } from "@/lib/server/verify";
import { getServiceConfig } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getUserByID } from "@/lib/zitadel";
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

  const { userId, loginName, code, organization, requestId, invite, send } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const branding = await getBrandingSettings({ serviceConfig, organization });

  let sessionFactors;
  let user: User | undefined;
  let human: HumanUser | undefined;
  let id: string | undefined;

  let error: string | undefined;

  const doSend = send === "true";

  const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? "";

  async function sendEmail(userId: string) {
    const hostWithProtocol = await getPublicHostWithProtocol(_headers);

    if (invite === "true") {
      await sendInviteEmailCode({
        userId,
        urlTemplate:
          `${hostWithProtocol}${basePath}/verify/email?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&invite=true` +
          (requestId ? `&requestId=${requestId}` : ""),
      }).catch((apiError) => {
        console.error("Could not send invitation email", apiError);
        error = "inviteSendFailed";
      });
    } else {
      await sendEmailCode({
        userId,
        urlTemplate:
          `${hostWithProtocol}${basePath}/verify/email?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
          (requestId ? `&requestId=${requestId}` : ""),
      }).catch((apiError) => {
        console.error("Could not send verification email", apiError);
        error = "emailSendFailed";
      });
    }
  }

  if ("loginName" in searchParams) {
    sessionFactors = await loadMostRecentSession({
      serviceConfig,
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

    const userResponse = await getUserByID({ serviceConfig, userId });
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
    initial: "true",
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
      <div className="flex flex-col space-y-4 w-full">
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
          user && (
            <UserAvatar loginName={user.preferredLoginName} displayName={human?.profile?.displayName} showDropdown={false} />
          )
        )}
      </div>

      <div className="w-full">
        {error && (
          <div className="py-4">
            <Alert>
              <Translated i18nKey={`errors.${error}`} namespace="verify" />
            </Alert>
          </div>
        )}

        {!id && (
          <div className="py-4">
            <Alert>
              <Translated i18nKey="unknownContext" namespace="error" />
            </Alert>
          </div>
        )}

        {id && send && (
          <div className="w-full">
            <Alert type={AlertType.INFO}>
              <Translated i18nKey="verify.codeSent" namespace="verify" />
            </Alert>
          </div>
        )}

        <VerifyEmailForm
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

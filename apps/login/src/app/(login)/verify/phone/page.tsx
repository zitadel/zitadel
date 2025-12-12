import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { Translated } from "@/components/translated";
import { UserAvatar } from "@/components/user-avatar";
import { VerifyPhoneForm } from "@/components/verify-phone-form";
import { getServiceConfig } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { getBrandingSettings, getUserByID, resendPhoneCode } from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { Metadata } from "next";
import { getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export async function generateMetadata(): Promise<Metadata> {
  const t = await getTranslations("verify");
  return { title: t("verifyPhone.title") };
}

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;

  const { userId, loginName, code, organization, requestId, send } = searchParams;

  const _headers = await headers();
  const { serviceConfig } = getServiceConfig(_headers);

  const branding = await getBrandingSettings({ serviceConfig, organization });

  let sessionFactors;
  let user: User | undefined;
  let human: HumanUser | undefined;
  let id: string | undefined;

  let error: string | undefined;

  const doSend = send === "true";

  async function sendPhone(userId: string) {
    await resendPhoneCode({ serviceConfig, userId }).catch((apiError) => {
      console.error("Could not send phone verification SMS", apiError);
      error = "phoneSendFailed";
    });
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
      await sendPhone(sessionFactors.factors.user.id);
    }
  } else if ("userId" in searchParams && userId) {
    if (doSend) {
      await sendPhone(userId);
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
          <Translated i18nKey="verifyPhone.title" namespace="verify" />
        </h1>
        <p className="ztdl-p">
          <Translated i18nKey="verifyPhone.description" namespace="verify" />
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
              <Translated i18nKey="verifyPhone.codeSent" namespace="verify" />
            </Alert>
          </div>
        )}

        <VerifyPhoneForm
          loginName={loginName}
          organization={organization}
          userId={id}
          code={code}
          requestId={requestId}
        />
      </div>
    </DynamicTheme>
  );
}

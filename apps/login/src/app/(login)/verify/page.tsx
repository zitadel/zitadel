import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { VerifyForm } from "@/components/verify-form";
import { VerifyRedirectButton } from "@/components/verify-redirect-button";
import { sendEmailCode } from "@/lib/server/verify";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { loadMostRecentSession } from "@/lib/session";
import { checkUserVerification } from "@/lib/verification-helper";
import {
  getBrandingSettings,
  getUserByID,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "verify" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const { userId, loginName, code, organization, requestId, invite, send } =
    searchParams;

  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);
  const host = _headers.get("host");

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

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

  if ("loginName" in searchParams) {
    sessionFactors = await loadMostRecentSession({
      serviceUrl,
      sessionParams: {
        loginName,
        organization,
      },
    });

    if (doSend && sessionFactors?.factors?.user?.id) {
      await sendEmailCode({
        serviceUrl,
        userId: sessionFactors?.factors?.user?.id,
        urlTemplate:
          `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
          (requestId ? `&requestId=${requestId}` : ""),
      }).catch((error) => {
        console.error("Could not resend verification email", error);
        throw Error("Failed to send verification email");
      });
    }
  } else if ("userId" in searchParams && userId) {
    if (doSend) {
      await sendEmailCode({
        serviceUrl,
        userId,
        urlTemplate:
          `${host.includes("localhost") ? "http://" : "https://"}${host}${basePath}/verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}` +
          (requestId ? `&requestId=${requestId}` : ""),
      }).catch((error) => {
        console.error("Could not resend verification email", error);
        throw Error("Failed to send verification email");
      });
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

  let authMethods: AuthenticationMethodType[] | null = null;
  if (human?.email?.isVerified) {
    const authMethodsResponse = await listAuthenticationMethodTypes({
      serviceUrl,
      userId,
    });
    if (authMethodsResponse.authMethodTypes) {
      authMethods = authMethodsResponse.authMethodTypes;
    }
  }

  const hasValidUserVerificationCheck = await checkUserVerification(id);

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
        <h1>{t("verify.title")}</h1>
        <p className="ztdl-p mb-6 block">{t("verify.description")}</p>

        {!id && (
          <>
            <h1>{t("verify.title")}</h1>
            <p className="ztdl-p mb-6 block">{t("verify.description")}</p>

            <div className="py-4">
              <Alert>{tError("unknownContext")}</Alert>
            </div>
          </>
        )}

        {id && send && (
          <div className="py-4 w-full">
            <Alert type={AlertType.INFO}>{t("verify.codeSent")}</Alert>
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

        {/* show a button to setup auth method for the user otherwise show the UI for reverifying */}
        {human?.email?.isVerified && hasValidUserVerificationCheck ? (
          // show page for already verified users
          <VerifyRedirectButton
            userId={id}
            loginName={loginName}
            organization={organization}
            requestId={requestId}
            authMethods={authMethods}
          />
        ) : (
          // check if auth methods are set
          <VerifyForm
            loginName={loginName}
            organization={organization}
            userId={id}
            code={code}
            isInvite={invite === "true"}
            requestId={requestId}
          />
        )}
      </div>
    </DynamicTheme>
  );
}

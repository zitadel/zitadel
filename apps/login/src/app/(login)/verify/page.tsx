import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { VerifyForm } from "@/components/verify-form";
import { VerifyRedirectButton } from "@/components/verify-redirect-button";
import { sendCode } from "@/lib/server/verify";
import { loadMostRecentSession } from "@/lib/session";
import {
  getBrandingSettings,
  getUserByID,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page(props: { searchParams: Promise<any> }) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "verify" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const {
    userId,
    loginName,
    code,
    organization,
    authRequestId,
    invite,
    skipsend,
  } = searchParams;

  const branding = await getBrandingSettings(organization);

  let sessionFactors;
  let user: User | undefined;
  let human: HumanUser | undefined;
  let id: string | undefined;

  if ("loginName" in searchParams) {
    sessionFactors = await loadMostRecentSession({
      loginName,
      organization,
    });

    if (!skipsend && sessionFactors?.factors?.user?.id) {
      await sendCode({
        userId: sessionFactors?.factors?.user?.id,
        isInvite: invite === "true",
      }).catch((error) => {
        console.error("Could not resend verification email", error);
        throw Error("Could not request email");
      });
    }
  } else if ("userId" in searchParams && userId) {
    if (!skipsend) {
      await sendCode({
        userId,
        isInvite: invite === "true",
      }).catch((error) => {
        console.error("Could not resend verification email", error);
        throw Error("Could not request email");
      });
    }

    const userResponse = await getUserByID(userId);
    if (userResponse) {
      user = userResponse.user;
      if (user?.type.case === "human") {
        human = user.type.value as HumanUser;
      }
    }
  }

  id = userId ?? sessionFactors?.factors?.user?.id;

  let authMethods: AuthenticationMethodType[] | null = null;
  if (human?.email?.isVerified) {
    const authMethodsResponse = await listAuthenticationMethodTypes(userId);
    if (authMethodsResponse.authMethodTypes) {
      authMethods = authMethodsResponse.authMethodTypes;
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

  if (authRequestId) {
    params.set("authRequestId", authRequestId);
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

        {id &&
          (human?.email?.isVerified ? (
            // show page for already verified users
            <VerifyRedirectButton
              userId={id}
              loginName={loginName}
              organization={organization}
              authRequestId={authRequestId}
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
              authRequestId={authRequestId}
            />
          ))}
      </div>
    </DynamicTheme>
  );
}

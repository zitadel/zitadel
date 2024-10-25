import { Alert, AlertType } from "@/components/alert";
import { BackButton } from "@/components/back-button";
import { Button, ButtonVariants } from "@/components/button";
import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { VerifyForm } from "@/components/verify-form";
import {
  getBrandingSettings,
  getUserByID,
  listAuthenticationMethodTypes,
} from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getLocale, getTranslations } from "next-intl/server";
import Link from "next/link";

export default async function Page({ searchParams }: { searchParams: any }) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "verify" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const { userId, loginName, code, organization, authRequestId, invite } =
    searchParams;

  const branding = await getBrandingSettings(organization);

  let user: User | undefined;
  let human: HumanUser | undefined;
  if (userId) {
    const userResponse = await getUserByID(userId);
    if (userResponse) {
      user = userResponse.user;
      if (user?.type.case === "human") {
        human = user.type.value as HumanUser;
      }
    }
  }

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

        {!userId && (
          <>
            <h1>{t("verify.title")}</h1>
            <p className="ztdl-p mb-6 block">{t("verify.description")}</p>

            <div className="py-4">
              <Alert>{tError("unknownContext")}</Alert>
            </div>
          </>
        )}

        {user && (
          <UserAvatar
            loginName={user.preferredLoginName}
            displayName={human?.profile?.displayName}
            showDropdown={false}
          />
        )}
        {human?.email?.isVerified ? (
          <>
            <Alert type={AlertType.INFO}>{t("success")}</Alert>

            <div className="mt-8 flex w-full flex-row items-center">
              <BackButton />
              <span className="flex-grow"></span>
              {authMethods?.length !== 0 && (
                <Link href={`/authenticator/set?+${params}`}>
                  <Button
                    type="submit"
                    className="self-end"
                    variant={ButtonVariants.Primary}
                  >
                    {t("setupAuthenticator")}
                  </Button>
                </Link>
              )}
            </div>
          </>
        ) : (
          // check if auth methods are set
          <VerifyForm
            userId={userId}
            code={code}
            isInvite={invite === "true"}
            params={params}
          />
        )}
      </div>
    </DynamicTheme>
  );
}

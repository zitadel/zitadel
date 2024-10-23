import { Alert } from "@/components/alert";
import { AuthenticatorMethods } from "@/components/authenticator-methods";
import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { VerifyForm } from "@/components/verify-form";
import { verifyUser } from "@/lib/server/email";
import { getBrandingSettings, getUserByID } from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({ searchParams }: { searchParams: any }) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "verify" });
  const tError = await getTranslations({ locale, namespace: "error" });

  const {
    userId,
    loginName,
    sessionId,
    code,
    organization,
    authRequestId,
    invite,
  } = searchParams;

  const branding = await getBrandingSettings(organization);

  let verifyResponse, error;
  if (code && userId) {
    verifyResponse = await verifyUser({
      code,
      userId,
      isInvite: invite === "true",
    }).catch(() => {
      error = "Could not verify user";
    });
  }

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

  const params = new URLSearchParams({
    userId: userId,
    initial: "true", // defines that a code is not required and is therefore not shown in the UI
  });

  if (organization) {
    params.set("organization", organization);
  }

  if (authRequestId) {
    params.set("authRequest", authRequestId);
  }

  if (sessionId) {
    params.set("sessionId", sessionId);
  }

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        {!userId && (
          <>
            <h1>{t("verify.title")}</h1>
            <p className="ztdl-p mb-6 block">{t("verify.description")}</p>

            <div className="py-4">
              <Alert>{tError("unknownContext")}</Alert>
            </div>
          </>
        )}
        {!verifyResponse || !verifyResponse.authMethodTypes ? (
          <VerifyForm
            userId={userId}
            loginName={loginName}
            code={!error ? code : ""}
            organization={organization}
            authRequestId={authRequestId}
            sessionId={sessionId}
            isInvite={invite === "true"}
          />
        ) : (
          <>
            <h1>{t("setup.title")}</h1>
            <p className="ztdl-p mb-6 block">{t("setup.description")}</p>
            {user && (
              <UserAvatar
                loginName={user.preferredLoginName}
                displayName={human?.profile?.displayName}
                showDropdown={false}
              />
            )}
            <AuthenticatorMethods
              authMethods={verifyResponse.authMethodTypes}
              params={params}
            />
          </>
        )}
      </div>
    </DynamicTheme>
  );
}

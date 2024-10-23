import { Alert } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { VerifyForm } from "@/components/verify-form";
import { getBrandingSettings, getUserByID } from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getLocale, getTranslations } from "next-intl/server";

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
        {!userId && (
          <>
            <h1>{t("verify.title")}</h1>
            <p className="ztdl-p mb-6 block">{t("verify.description")}</p>

            <div className="py-4">
              <Alert>{tError("unknownContext")}</Alert>
            </div>
          </>
        )}

        <VerifyForm
          userId={userId}
          code={code}
          isInvite={invite === "true"}
          params={params}
        />
      </div>
    </DynamicTheme>
  );
}

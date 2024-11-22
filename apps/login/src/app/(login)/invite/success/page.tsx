import { Alert, AlertType } from "@/components/alert";
import { Button, ButtonVariants } from "@/components/button";
import { DynamicTheme } from "@/components/dynamic-theme";
import { UserAvatar } from "@/components/user-avatar";
import { getBrandingSettings, getDefaultOrg, getUserByID } from "@/lib/zitadel";
import { HumanUser, User } from "@zitadel/proto/zitadel/user/v2/user_pb";
import { getLocale, getTranslations } from "next-intl/server";
import Link from "next/link";

export default async function Page(
  props: {
    searchParams: Promise<Record<string | number | symbol, string | undefined>>;
  }
) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "invite" });

  let { userId, organization } = searchParams;

  if (!organization) {
    const org = await getDefaultOrg();
    if (!org) {
      throw new Error("No default organization found");
    }

    organization = org.id;
  }

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

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("success.title")}</h1>
        <p className="ztdl-p">{t("success.description")}</p>
        {user && (
          <UserAvatar
            loginName={user.preferredLoginName}
            displayName={human?.profile?.displayName}
            showDropdown={false}
          />
        )}
        {human?.email?.isVerified ? (
          <Alert type={AlertType.INFO}>{t("success.verified")}</Alert>
        ) : (
          <Alert type={AlertType.INFO}>{t("success.notVerifiedYet")}</Alert>
        )}
        <div className="mt-8 flex w-full flex-row items-center justify-between">
          <span></span>
          <Link href="/invite">
            <Button type="submit" variant={ButtonVariants.Primary}>
              {t("success.submit")}
            </Button>
          </Link>
        </div>
      </div>
    </DynamicTheme>
  );
}

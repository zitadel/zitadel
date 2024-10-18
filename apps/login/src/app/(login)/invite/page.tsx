import { Alert, AlertType } from "@/components/alert";
import { DynamicTheme } from "@/components/dynamic-theme";
import { InviteForm } from "@/components/invite-form";
import {
  getBrandingSettings,
  getDefaultOrg,
  getPasswordComplexitySettings,
} from "@/lib/zitadel";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "invite" });

  let { firstname, lastname, email, organization } = searchParams;

  if (!organization) {
    const org = await getDefaultOrg();
    if (!org) {
      throw new Error("No default organization found");
    }

    organization = org.id;
  }

  const passwordComplexitySettings =
    await getPasswordComplexitySettings(organization);

  const branding = await getBrandingSettings(organization);

  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("title")}</h1>
        <p className="ztdl-p">{t("description")}</p>

        <Alert type={AlertType.INFO}>{t("info")}</Alert>

        {passwordComplexitySettings && (
          <InviteForm
            organization={organization}
            firstname={firstname}
            lastname={lastname}
            email={email}
          ></InviteForm>
        )}
      </div>
    </DynamicTheme>
  );
}

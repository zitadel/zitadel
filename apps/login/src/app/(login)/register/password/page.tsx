import { DynamicTheme } from "@/components/dynamic-theme";
import { SetRegisterPasswordForm } from "@/components/set-register-password-form";
import {
  getBrandingSettings,
  getDefaultOrg,
  getLegalAndSupportSettings,
  getLoginSettings,
  getPasswordComplexitySettings,
} from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { getLocale, getTranslations } from "next-intl/server";
import { headers } from "next/headers";

export default async function Page(props: {
  searchParams: Promise<Record<string | number | symbol, string | undefined>>;
}) {
  const searchParams = await props.searchParams;
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "register" });

  let { firstname, lastname, email, organization, authRequestId } =
    searchParams;

  const host = (await headers()).get("host");

  if (!host || typeof host !== "string") {
    throw new Error("No host found");
  }

  if (!organization) {
    const org: Organization | null = await getDefaultOrg({ host });
    if (org) {
      organization = org.id;
    }
  }

  const missingData = !firstname || !lastname || !email;

  const legal = await getLegalAndSupportSettings({ host, organization });
  const passwordComplexitySettings = await getPasswordComplexitySettings({
    host,
    organization,
  });

  const branding = await getBrandingSettings({ host, organization });

  const loginSettings = await getLoginSettings({ host, organization });

  return missingData ? (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("missingdata.title")}</h1>
        <p className="ztdl-p">{t("missingdata.description")}</p>
      </div>
    </DynamicTheme>
  ) : loginSettings?.allowRegister && loginSettings.allowUsernamePassword ? (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("password.title")}</h1>
        <p className="ztdl-p">{t("description")}</p>

        {legal && passwordComplexitySettings && (
          <SetRegisterPasswordForm
            passwordComplexitySettings={passwordComplexitySettings}
            email={email}
            firstname={firstname}
            lastname={lastname}
            organization={organization}
            authRequestId={authRequestId}
          ></SetRegisterPasswordForm>
        )}
      </div>
    </DynamicTheme>
  ) : (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>{t("disabled.title")}</h1>
        <p className="ztdl-p">{t("disabled.description")}</p>
      </div>
    </DynamicTheme>
  );
}

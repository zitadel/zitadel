import { DynamicTheme } from "@/components/dynamic-theme";
import { RegisterFormWithoutPassword } from "@/components/register-form-without-password";
import { SetRegisterPasswordForm } from "@/components/set-register-password-form";
import {
  getBrandingSettings,
  getDefaultOrg,
  getLegalAndSupportSettings,
  getPasswordComplexitySettings,
} from "@/lib/zitadel";
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";
import { getLocale, getTranslations } from "next-intl/server";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const locale = getLocale();
  const t = await getTranslations({ locale, namespace: "register" });

  let { firstname, lastname, email, organization, authRequestId } =
    searchParams;

  if (!organization) {
    const org: Organization | null = await getDefaultOrg();
    if (org) {
      organization = org.id;
    }
  }

  const setPassword = !!(firstname && lastname && email);

  const legal = await getLegalAndSupportSettings(organization);
  const passwordComplexitySettings =
    await getPasswordComplexitySettings(organization);

  const branding = await getBrandingSettings(organization);

  return setPassword ? (
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
        <h1>{t("title")}</h1>
        <p className="ztdl-p">{t("description")}</p>

        {legal && passwordComplexitySettings && (
          <RegisterFormWithoutPassword
            legal={legal}
            organization={organization}
            firstname={firstname}
            lastname={lastname}
            email={email}
            authRequestId={authRequestId}
          ></RegisterFormWithoutPassword>
        )}
      </div>
    </DynamicTheme>
  );
}

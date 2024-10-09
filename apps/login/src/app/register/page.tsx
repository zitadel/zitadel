import { DynamicTheme } from "@/components/dynamic-theme";
import { RegisterFormWithoutPassword } from "@/components/register-form-without-password";
import { SetPasswordForm } from "@/components/set-password-form";
import {
  getBrandingSettings,
  getLegalAndSupportSettings,
  getPasswordComplexitySettings,
} from "@/lib/zitadel";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { firstname, lastname, email, organization, authRequestId } =
    searchParams;

  if (!organization) {
    // TODO: get default organization
  }

  const setPassword = !!(firstname && lastname && email);

  const legal = await getLegalAndSupportSettings(organization);
  const passwordComplexitySettings =
    await getPasswordComplexitySettings(organization);

  const branding = await getBrandingSettings(organization);

  return setPassword ? (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Set Password</h1>
        <p className="ztdl-p">Set the password for your account</p>

        {legal && passwordComplexitySettings && (
          <SetPasswordForm
            passwordComplexitySettings={passwordComplexitySettings}
            email={email}
            firstname={firstname}
            lastname={lastname}
            organization={organization}
            authRequestId={authRequestId}
          ></SetPasswordForm>
        )}
      </div>
    </DynamicTheme>
  ) : (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col items-center space-y-4">
        <h1>Register</h1>
        <p className="ztdl-p">Create your ZITADEL account.</p>

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

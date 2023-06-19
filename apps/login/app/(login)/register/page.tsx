import {
  getLegalAndSupportSettings,
  getPasswordComplexitySettings,
  server,
} from "#/lib/zitadel";
import RegisterFormWithoutPassword from "#/ui/RegisterFormWithoutPassword";
import SetPasswordForm from "#/ui/SetPasswordForm";

export default async function Page({
  searchParams,
}: {
  searchParams: Record<string | number | symbol, string | undefined>;
}) {
  const { firstname, lastname, email } = searchParams;

  const setPassword = !!(firstname && lastname && email);

  const legal = await getLegalAndSupportSettings(server);
  const passwordComplexitySettings = await getPasswordComplexitySettings(
    server
  );

  return setPassword ? (
    <div className="flex flex-col items-center space-y-4">
      <h1>Set Password</h1>
      <p className="ztdl-p">Set the password for your account</p>

      {legal && passwordComplexitySettings && (
        <SetPasswordForm
          passwordComplexitySettings={passwordComplexitySettings}
          email={email}
          firstname={firstname}
          lastname={lastname}
        ></SetPasswordForm>
      )}
    </div>
  ) : (
    <div className="flex flex-col items-center space-y-4">
      <h1>Register</h1>
      <p className="ztdl-p">Create your ZITADEL account.</p>

      {legal && passwordComplexitySettings && (
        <RegisterFormWithoutPassword
          legal={legal}
        ></RegisterFormWithoutPassword>
      )}
    </div>
  );
}

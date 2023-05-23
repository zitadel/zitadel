import {
  getLegalAndSupportSettings,
  getPasswordComplexitySettings,
  server,
} from "#/lib/zitadel";
import RegisterForm from "#/ui/RegisterForm";

export default async function Page() {
  const legal = await getLegalAndSupportSettings(server);
  const passwordComplexitySettings = await getPasswordComplexitySettings(
    server
  );

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Register</h1>
      <p className="ztdl-p">Create your ZITADEL account.</p>

      {legal && passwordComplexitySettings && (
        <RegisterForm
          legal={legal}
          passwordComplexitySettings={passwordComplexitySettings}
        ></RegisterForm>
      )}
    </div>
  );
}

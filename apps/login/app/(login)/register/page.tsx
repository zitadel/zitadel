import {
  getPasswordComplexityPolicy,
  getPrivacyPolicy,
  server,
} from "#/lib/zitadel";
import RegisterForm from "#/ui/RegisterForm";

export default async function Page() {
  const privacyPolicy = await getPrivacyPolicy(server);
  const passwordComplexityPolicy = await getPasswordComplexityPolicy(server);

  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Register</h1>
      <p className="ztdl-p">Create your ZITADEL account.</p>

      {privacyPolicy && passwordComplexityPolicy && (
        <RegisterForm
          privacyPolicy={privacyPolicy}
          passwordComplexityPolicy={passwordComplexityPolicy}
        ></RegisterForm>
      )}
    </div>
  );
}

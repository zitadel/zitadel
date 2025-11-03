import { LoginSettings, PasskeysType } from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { Alert, AlertType } from "./alert";
import { PASSKEYS, PASSWORD } from "./auth-methods";
import { Translated } from "./translated";

type Props = {
  authMethods: AuthenticationMethodType[];
  params: URLSearchParams;
  loginSettings: LoginSettings;
};

export function ChooseAuthenticatorToSetup({ authMethods, params, loginSettings }: Props) {
  if (authMethods.length !== 0) {
    return (
      <Alert type={AlertType.ALERT}>
        <Translated i18nKey="allSetup" namespace="authenticator" />
      </Alert>
    );
  } else {
    return (
      <div className="grid w-full grid-cols-1 gap-5 pt-4">
        {!authMethods.includes(AuthenticationMethodType.PASSWORD) &&
          loginSettings.allowUsernamePassword &&
          PASSWORD(false, "/password/set?" + params)}
        {!authMethods.includes(AuthenticationMethodType.PASSKEY) &&
          loginSettings.passkeysType == PasskeysType.ALLOWED &&
          PASSKEYS(false, "/passkey/set?" + params)}
      </div>
    );
  }
}

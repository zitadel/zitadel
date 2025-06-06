import {
  LoginSettings,
  PasskeysType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { Alert, AlertType } from "./alert";
import { PASSKEYS, PASSWORD } from "./auth-methods";

type Props = {
  authMethods: AuthenticationMethodType[];
  params: URLSearchParams;
  loginSettings: LoginSettings;
};

export function ChooseAuthenticatorToSetup({
  authMethods,
  params,
  loginSettings,
}: Props) {
  const t = useTranslations("authenticator");

  if (authMethods.length !== 0) {
    return <Alert type={AlertType.ALERT}>{t("allSetup")}</Alert>;
  } else {
    return (
      <>
        {loginSettings.passkeysType == PasskeysType.NOT_ALLOWED &&
          !loginSettings.allowUsernamePassword && (
            <Alert type={AlertType.ALERT}>{t("noMethodsAvailable")}</Alert>
          )}

        <div className="grid grid-cols-1 gap-5 w-full pt-4">
          {!authMethods.includes(AuthenticationMethodType.PASSWORD) &&
            loginSettings.allowUsernamePassword &&
            PASSWORD(false, "/password/set?" + params)}
          {!authMethods.includes(AuthenticationMethodType.PASSKEY) &&
            loginSettings.passkeysType == PasskeysType.ALLOWED &&
            PASSKEYS(false, "/passkey/set?" + params)}
        </div>
      </>
    );
  }
}

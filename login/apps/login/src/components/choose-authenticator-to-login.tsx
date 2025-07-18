import {
  LoginSettings,
  PasskeysType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { PASSKEYS, PASSWORD } from "./auth-methods";
import { Translated } from "./translated";

type Props = {
  authMethods: AuthenticationMethodType[];
  params: URLSearchParams;
  loginSettings: LoginSettings | undefined;
};

export function ChooseAuthenticatorToLogin({
  authMethods,
  params,
  loginSettings,
}: Props) {
  return (
    <>
      {authMethods.includes(AuthenticationMethodType.PASSWORD) &&
        loginSettings?.allowUsernamePassword && (
          <div className="ztdl-p">
            <Translated i18nKey="chooseAlternativeMethod" namespace="idp" />
          </div>
        )}
      <div className="grid w-full grid-cols-1 gap-5 pt-4">
        {authMethods.includes(AuthenticationMethodType.PASSWORD) &&
          loginSettings?.allowUsernamePassword &&
          PASSWORD(false, "/password?" + params)}
        {authMethods.includes(AuthenticationMethodType.PASSKEY) &&
          loginSettings?.passkeysType == PasskeysType.ALLOWED &&
          PASSKEYS(false, "/passkey?" + params)}
      </div>
    </>
  );
}

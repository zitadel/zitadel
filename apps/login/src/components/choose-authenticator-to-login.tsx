import {
  LoginSettings,
  PasskeysType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { PASSKEYS, PASSWORD } from "./auth-methods";

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
  const t = useTranslations("idp");

  return (
    <>
      {authMethods.includes(AuthenticationMethodType.PASSWORD) &&
        loginSettings?.allowUsernamePassword && (
          <div className="ztdl-p">Choose an alternative method to login </div>
        )}
      <div className="grid grid-cols-1 gap-5 w-full pt-4">
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

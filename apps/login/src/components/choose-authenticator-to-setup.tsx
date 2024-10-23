import { Factors } from "@zitadel/proto/zitadel/session/v2/session_pb";
import {
  LoginSettings,
  PasskeysType,
} from "@zitadel/proto/zitadel/settings/v2/login_settings_pb";
import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { useTranslations } from "next-intl";
import { Alert, AlertType } from "./alert";
import { PASSKEYS, PASSWORD } from "./auth-methods";
import { UserAvatar } from "./user-avatar";

type Props = {
  authMethods: AuthenticationMethodType[];
  params: URLSearchParams;
  sessionFactors?: Factors;
  loginSettings: LoginSettings;
};

export function ChooseAuthenticatorToSetup({
  authMethods,
  params,
  sessionFactors,
  loginSettings,
}: Props) {
  const t = useTranslations("authenticator");

  return (
    <>
      {sessionFactors && (
        <UserAvatar
          loginName={sessionFactors.user?.loginName}
          displayName={sessionFactors.user?.displayName}
          showDropdown
        ></UserAvatar>
      )}

      {loginSettings.passkeysType === PasskeysType.ALLOWED &&
        !loginSettings.allowUsernamePassword && (
          <Alert type={AlertType.ALERT}>{t("noMethodsAvailable")}</Alert>
        )}

      <div className="grid grid-cols-1 gap-5 w-full pt-4">
        {!authMethods.includes(AuthenticationMethodType.PASSWORD) &&
          loginSettings.allowUsernamePassword &&
          PASSWORD(false, "/password/set?" + params)}
        {!authMethods.includes(AuthenticationMethodType.PASSKEY) &&
          loginSettings.passkeysType === PasskeysType.ALLOWED &&
          PASSKEYS(false, "/passkeys/set?" + params)}
      </div>
    </>
  );
}

import { AuthenticationMethodType } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { PASSKEYS, PASSWORD } from "./auth-methods";

type Props = {
  authMethods: AuthenticationMethodType[];
  params: URLSearchParams;
};

export function AuthenticatorMethods({ authMethods, params }: Props) {
  return (
    <div className="grid grid-cols-1 gap-5 w-full pt-4">
      {!authMethods.includes(AuthenticationMethodType.PASSWORD) &&
        PASSWORD(false, "/password/set?" + params)}
      {!authMethods.includes(AuthenticationMethodType.PASSKEY) &&
        PASSKEYS(false, "/passkeys/set?" + params)}
    </div>
  );
}

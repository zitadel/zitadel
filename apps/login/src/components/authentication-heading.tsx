import { Translated } from "@/components/translated";
import { getAuthRequest, type ServiceConfig } from "@/lib/zitadel";
import { Prompt } from "@zitadel/proto/zitadel/oidc/v2/authorization_pb";

export async function AuthenticationHeading({
  requestId,
  serviceConfig,
  namespace,
}: {
  requestId?: string;
  serviceConfig: ServiceConfig;
  namespace: "loginname" | "password";
}) {
  // Read the validated IdP request rather than trusting a display hint supplied
  // in the browser URL. This changes copy only, never authentication policy.
  const request = requestId?.startsWith("oidc_")
    ? await getAuthRequest({ serviceConfig, authRequestId: requestId.slice(5) })
    : undefined;
  const reauthenticate = request?.authRequest?.prompt.includes(Prompt.LOGIN);
  const messages = reauthenticate ? "reauthenticate" : namespace;
  const prefix = !reauthenticate && namespace === "password" ? "verify." : "";

  return (
    <>
      <h1>
        <Translated i18nKey={`${prefix}title`} namespace={messages} />
      </h1>
      <p className="ztdl-p">
        <Translated i18nKey={`${prefix}description`} namespace={messages} />
      </p>
    </>
  );
}

import { createClientFor } from "@zitadel/client";
import { IdentityProviderService } from "@zitadel/proto/zitadel/idp/v2/idp_service_pb";
import { OIDCService } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { SAMLService } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { systemAPIToken, loginServiceKeyToken } from "./api";
import { hasSystemUserCredentials, hasLoginServiceKey, hasServiceUserToken } from "./deployment";
import { createServerTransport, ServiceConfig } from "./zitadel";

type ServiceClass =
  | typeof IdentityProviderService
  | typeof UserService
  | typeof OrganizationService
  | typeof SessionService
  | typeof OIDCService
  | typeof SettingsService
  | typeof SAMLService;

export async function createServiceForHost<T extends ServiceClass>(service: T, serviceConfig: ServiceConfig) {
  let token;

  // Determine authentication method based on available credentials
  // Priority: system user JWT > login service key > service account token
  if (hasSystemUserCredentials()) {
    token = await systemAPIToken();
  } else if (hasLoginServiceKey()) {
    token = await loginServiceKeyToken();
  } else if (hasServiceUserToken()) {
    // Use service account token authentication (self-hosted)
    token = process.env.ZITADEL_SERVICE_USER_TOKEN;
  } else {
    throw new Error(
      "No authentication credentials found. Set either system user credentials (AUDIENCE, SYSTEM_USER_ID, SYSTEM_USER_PRIVATE_KEY), login service key (ZITADEL_LOGIN_SERVICE_KEY_FILE with ZITADEL_LOGIN_SYSTEM_USER_ID or SYSTEM_USER_ID), or ZITADEL_SERVICE_USER_TOKEN",
    );
  }

  if (!serviceConfig) {
    throw new Error("No service config found");
  }

  const transport = createServerTransport(token, serviceConfig);

  return createClientFor<T>(service)(transport);
}

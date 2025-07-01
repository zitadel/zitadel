import { createClientFor } from "@zitadel/client";
import { createServerTransport } from "./zitadel";
import { IdentityProviderService } from "@zitadel/proto/zitadel/idp/v2/idp_service_pb";
import { OIDCService } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { SAMLService } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { systemAPIToken } from "./api";

type ServiceClass =
  | typeof IdentityProviderService
  | typeof UserService
  | typeof OrganizationService
  | typeof SessionService
  | typeof OIDCService
  | typeof SettingsService
  | typeof SAMLService;

export async function createServiceForHost<T extends ServiceClass>(
  service: T,
  serviceUrl: string,
) {
  let token;

  // if we are running in a multitenancy context, use the system user token
  if (
    process.env.AUDIENCE &&
    process.env.SYSTEM_USER_ID &&
    process.env.SYSTEM_USER_PRIVATE_KEY
  ) {
    token = await systemAPIToken();
  } else if (process.env.ZITADEL_SERVICE_USER_TOKEN) {
    token = process.env.ZITADEL_SERVICE_USER_TOKEN;
  }

  if (!serviceUrl) {
    throw new Error("No instance url found");
  }

  if (!token) {
    throw new Error("No token found");
  }

  const transport = createServerTransport(token, {
    baseUrl: serviceUrl,
    interceptors: !process.env.CUSTOM_REQUEST_HEADERS
      ? undefined
      : [
          (next) => {
            return (req) => {
              process.env.CUSTOM_REQUEST_HEADERS!.split(",").forEach(
                (header) => {
                  const kv = header.split(":");
                  req.header.set(kv[0], kv[1]);
                },
              );
              return next(req);
            };
          },
        ],
  });

  return createClientFor<T>(service)(transport);
}

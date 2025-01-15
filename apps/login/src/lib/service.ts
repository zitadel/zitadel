import { createClientFor } from "@zitadel/client";
import { createServerTransport } from "@zitadel/client/node";
import { IdentityProviderService } from "@zitadel/proto/zitadel/idp/v2/idp_service_pb";
import { OIDCService } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getInstanceUrl, systemAPIToken } from "./api";

type ServiceClass =
  | typeof IdentityProviderService
  | typeof UserService
  | typeof OrganizationService
  | typeof SessionService
  | typeof OIDCService
  | typeof SettingsService;

export async function createServiceForHost<T extends ServiceClass>(
  service: T,
  host: string,
) {
  let instanceUrl;
  try {
    instanceUrl = await getInstanceUrl(host);
  } catch (error) {
    console.error(
      "Could not get instance url, fallback to ZITADEL_API_URL",
      error,
    );
    instanceUrl = process.env.ZITADEL_API_URL;
  }

  const systemToken = await systemAPIToken();

  const transport = createServerTransport(systemToken, {
    baseUrl: instanceUrl,
  });

  return createClientFor<T>(service)(transport);
}

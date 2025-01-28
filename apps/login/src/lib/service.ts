import { createClientFor } from "@zitadel/client";
import { createServerTransport } from "@zitadel/client/node";
import { IdentityProviderService } from "@zitadel/proto/zitadel/idp/v2/idp_service_pb";
import { OIDCService } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { ReadonlyHeaders } from "next/dist/server/web/spec-extension/adapters/headers";
import { systemAPIToken } from "./api";

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
  const token = await systemAPIToken();

  if (!host || !token) {
    throw new Error("No instance url or token found");
  }

  const transport = createServerTransport(token, {
    baseUrl: host,
  });

  return createClientFor<T>(service)(transport);
}

export function getApiUrlOfHeaders(headers: ReadonlyHeaders): string {
  let instanceUrl: string = process.env.ZITADEL_API_URL;

  if (headers.get("x-zitadel-forward-host")) {
    instanceUrl = headers.get("x-zitadel-forward-host") as string;
  } else {
    const host = headers.get("host");

    if (host) {
      const [hostname, port] = host.split(":");
      if (hostname !== "localhost") {
        instanceUrl = host;
      }
    }
  }

  return instanceUrl;
}

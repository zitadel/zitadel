import { FeatureService } from "@zitadel/proto/gen/zitadel/feature/v2/feature_service_connect";
import { RequestContext } from "@zitadel/proto/gen/zitadel/object/v2/object_pb";
import { OIDCService } from "@zitadel/proto/gen/zitadel/oidc/v2/oidc_service_connect";
import { OrganizationService } from "@zitadel/proto/gen/zitadel/org/v2/org_service_connect";
import { SessionService } from "@zitadel/proto/gen/zitadel/session/v2/session_service_connect";
import { SettingsService } from "@zitadel/proto/gen/zitadel/settings/v2/settings_service_connect";
import { UserService } from "@zitadel/proto/gen/zitadel/user/v2/user_service_connect";
import { createClientFor } from "./helpers";

export const createUserServiceClient = createClientFor(UserService);
export const createSettingsServiceClient = createClientFor(SettingsService);
export const createSessionServiceClient = createClientFor(SessionService);
export const createOIDCServiceClient = createClientFor(OIDCService);
export const createOrganizationServiceClient = createClientFor(OrganizationService);
export const createFeatureServiceClient = createClientFor(FeatureService);

export function makeReqCtx(orgId: string | undefined): Partial<RequestContext> {
  return {
    resourceOwner: orgId ? { case: "orgId", value: orgId } : { case: "instance", value: true },
  };
}

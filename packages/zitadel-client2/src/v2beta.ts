import type { PartialMessage } from "@bufbuild/protobuf";

import { createClientFor } from "./helpers";
import { UserService } from "@zitadel/proto/zitadel/user/v2beta/user_service_connect";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2beta/settings_service_connect";
import { SessionService } from "@zitadel/proto/zitadel/session/v2beta/session_service_connect";
import { OIDCService } from "@zitadel/proto/zitadel/oidc/v2beta/oidc_service_connect";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2beta/org_service_connect";
import { FeatureService } from "@zitadel/proto/zitadel/feature/v2beta/feature_service_connect";
import type { RequestContext } from "@zitadel/proto/zitadel/object/v2beta/object_pb";

export const createUserServiceClient = createClientFor(UserService);
export const createSettingsServiceClient = createClientFor(SettingsService);
export const createSessionServiceClient = createClientFor(SessionService);
export const createOIDCServiceClient = createClientFor(OIDCService);
export const createOrganizationServiceClient =
  createClientFor(OrganizationService);
export const createFeatureServiceClient = createClientFor(FeatureService);

export function makeReqCtx(
  orgId: string | undefined,
): PartialMessage<RequestContext> {
  return {
    resourceOwner: orgId
      ? { case: "orgId", value: orgId }
      : { case: "instance", value: true },
  };
}

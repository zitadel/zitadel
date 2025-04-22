import { MessageInitShape } from "@bufbuild/protobuf";
import { FeatureService } from "@zitadel/proto/zitadel/feature/v2/feature_service_pb.js";
import { IdentityProviderService } from "@zitadel/proto/zitadel/idp/v2/idp_service_pb.js";
import { RequestContextSchema } from "@zitadel/proto/zitadel/object/v2/object_pb.js";
import { OIDCService } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb.js";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb.js";
import { SAMLService } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb.js";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb.js";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb.js";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb.js";

import { createClientFor } from "./helpers.js";

export const createUserServiceClient: ReturnType<typeof createClientFor<typeof UserService>> = createClientFor(UserService);
export const createSettingsServiceClient: ReturnType<typeof createClientFor<typeof SettingsService>> =
  createClientFor(SettingsService);
export const createSessionServiceClient: ReturnType<typeof createClientFor<typeof SessionService>> =
  createClientFor(SessionService);
export const createOIDCServiceClient: ReturnType<typeof createClientFor<typeof OIDCService>> = createClientFor(OIDCService);
export const createSAMLServiceClient: ReturnType<typeof createClientFor<typeof SAMLService>> = createClientFor(SAMLService);
export const createOrganizationServiceClient: ReturnType<typeof createClientFor<typeof OrganizationService>> =
  createClientFor(OrganizationService);
export const createFeatureServiceClient: ReturnType<typeof createClientFor<typeof FeatureService>> =
  createClientFor(FeatureService);
export const createIdpServiceClient: ReturnType<typeof createClientFor<typeof IdentityProviderService>> =
  createClientFor(IdentityProviderService);

export function makeReqCtx(orgId: string | undefined): MessageInitShape<typeof RequestContextSchema> {
  return { resourceOwner: orgId ? { case: "orgId", value: orgId } : { case: "instance", value: true } };
}

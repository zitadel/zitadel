import {
  createClientFor
} from "./chunk-27KHKGT3.js";

// src/v2.ts
import { create } from "@bufbuild/protobuf";
import { FeatureService } from "@zitadel/proto/zitadel/feature/v2/feature_service_pb.js";
import { IdentityProviderService } from "@zitadel/proto/zitadel/idp/v2/idp_service_pb.js";
import { RequestContextSchema } from "@zitadel/proto/zitadel/object/v2/object_pb.js";
import { OIDCService } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb.js";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb.js";
import { SAMLService } from "@zitadel/proto/zitadel/saml/v2/saml_service_pb.js";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb.js";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb.js";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb.js";
var createUserServiceClient = createClientFor(UserService);
var createSettingsServiceClient = createClientFor(SettingsService);
var createSessionServiceClient = createClientFor(SessionService);
var createOIDCServiceClient = createClientFor(OIDCService);
var createSAMLServiceClient = createClientFor(SAMLService);
var createOrganizationServiceClient = createClientFor(OrganizationService);
var createFeatureServiceClient = createClientFor(FeatureService);
var createIdpServiceClient = createClientFor(IdentityProviderService);
function makeReqCtx(orgId) {
  return create(RequestContextSchema, {
    resourceOwner: orgId ? { case: "orgId", value: orgId } : { case: "instance", value: true }
  });
}
export {
  createFeatureServiceClient,
  createIdpServiceClient,
  createOIDCServiceClient,
  createOrganizationServiceClient,
  createSAMLServiceClient,
  createSessionServiceClient,
  createSettingsServiceClient,
  createUserServiceClient,
  makeReqCtx
};
//# sourceMappingURL=v2.js.map
/**
 * ZITADEL v2 API service client factories.
 *
 * Each export is a curried factory: pass a Connect transport to get a typed client.
 *
 * @example
 * ```ts
 * import { createUserServiceClient } from "@zitadel/zitadel-js/api/v2";
 * import { createGrpcTransport } from "@zitadel/zitadel-js";
 *
 * const transport = createGrpcTransport({ baseUrl: "https://my.zitadel.cloud" });
 * const users = createUserServiceClient(transport);
 * const user = await users.getUser({ userId: "123" });
 * ```
 */
import { ActionService } from "./generated/zitadel/action/v2/action_service_pb.js";
import { FeatureService } from "./generated/zitadel/feature/v2/feature_service_pb.js";
import { IdentityProviderService } from "./generated/zitadel/idp/v2/idp_service_pb.js";
import { OIDCService } from "./generated/zitadel/oidc/v2/oidc_service_pb.js";
import { OrganizationService } from "./generated/zitadel/org/v2/org_service_pb.js";
import { SAMLService } from "./generated/zitadel/saml/v2/saml_service_pb.js";
import { SessionService } from "./generated/zitadel/session/v2/session_service_pb.js";
import { SettingsService } from "./generated/zitadel/settings/v2/settings_service_pb.js";
import { UserService } from "./generated/zitadel/user/v2/user_service_pb.js";

import { createClientFor } from "./client.js";

export const createUserServiceClient = createClientFor(UserService);
export const createSettingsServiceClient = createClientFor(SettingsService);
export const createSessionServiceClient = createClientFor(SessionService);
export const createOIDCServiceClient = createClientFor(OIDCService);
export const createSAMLServiceClient = createClientFor(SAMLService);
export const createOrganizationServiceClient =
  createClientFor(OrganizationService);
export const createFeatureServiceClient = createClientFor(FeatureService);
export const createIdpServiceClient = createClientFor(IdentityProviderService);
export const createActionServiceClient = createClientFor(ActionService);

export type { RequestChallenges } from "./generated/zitadel/session/v2/challenge_pb.js";
export type { UserAgent } from "./generated/zitadel/session/v2/session_pb.js";
export type { Checks } from "./generated/zitadel/session/v2/session_service_pb.js";

export { makeReqCtx } from "./context.js";
export type { RequestContext } from "./context.js";

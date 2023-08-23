import * as settings from "./v2/settings";
import * as session from "./v2/session";
import * as user from "./v2/user";
import * as oidc from "./v2/oidc";
import * as management from "./management";

import * as login from "./proto/server/zitadel/settings/v2alpha/login_settings";
import * as password from "./proto/server/zitadel/settings/v2alpha/password_settings";
import * as legal from "./proto/server/zitadel/settings/v2alpha/legal_settings";

export {
  BrandingSettings,
  Theme,
} from "./proto/server/zitadel/settings/v2alpha/branding_settings";

export {
  LoginSettings,
  IdentityProvider,
  IdentityProviderType,
} from "./proto/server/zitadel/settings/v2alpha/login_settings";

export {
  RequestChallenges,
  Challenges,
  Challenges_WebAuthN,
} from "./proto/server/zitadel/session/v2alpha/challenge";

export {
  GetAuthRequestRequest,
  GetAuthRequestResponse,
  CreateCallbackRequest,
  CreateCallbackResponse,
} from "./proto/server/zitadel/oidc/v2alpha/oidc_service";

export {
  Session,
  Factors,
} from "./proto/server/zitadel/session/v2alpha/session";
export {
  IDPInformation,
  IDPLink,
} from "./proto/server/zitadel/user/v2alpha/idp";
export {
  ListSessionsResponse,
  GetSessionResponse,
  CreateSessionResponse,
  SetSessionResponse,
  SetSessionRequest,
  DeleteSessionResponse,
} from "./proto/server/zitadel/session/v2alpha/session_service";
export {
  GetPasswordComplexitySettingsResponse,
  GetBrandingSettingsResponse,
  GetLegalAndSupportSettingsResponse,
  GetGeneralSettingsResponse,
  GetLoginSettingsResponse,
  GetLoginSettingsRequest,
  GetActiveIdentityProvidersResponse,
  GetActiveIdentityProvidersRequest,
} from "./proto/server/zitadel/settings/v2alpha/settings_service";
export {
  AddHumanUserResponse,
  AddHumanUserRequest,
  VerifyEmailResponse,
  VerifyPasskeyRegistrationRequest,
  VerifyPasskeyRegistrationResponse,
  RegisterPasskeyRequest,
  RegisterPasskeyResponse,
  CreatePasskeyRegistrationLinkResponse,
  CreatePasskeyRegistrationLinkRequest,
  ListAuthenticationMethodTypesResponse,
  ListAuthenticationMethodTypesRequest,
  AuthenticationMethodType,
  StartIdentityProviderIntentRequest,
  StartIdentityProviderIntentResponse,
  RetrieveIdentityProviderIntentRequest,
  RetrieveIdentityProviderIntentResponse,
} from "./proto/server/zitadel/user/v2alpha/user_service";
export {
  SetHumanPasswordResponse,
  SetHumanPasswordRequest,
} from "./proto/server/zitadel/management";
export * from "./proto/server/zitadel/idp";
export { type LegalAndSupportSettings } from "./proto/server/zitadel/settings/v2alpha/legal_settings";
export { type PasswordComplexitySettings } from "./proto/server/zitadel/settings/v2alpha/password_settings";
export { type ResourceOwnerType } from "./proto/server/zitadel/settings/v2alpha/settings";

import {
  getServers,
  initializeServer,
  ZitadelServer,
  ZitadelServerOptions,
} from "./server";
export * from "./middleware";

export {
  getServers,
  ZitadelServer,
  type ZitadelServerOptions,
  initializeServer,
  user,
  management,
  session,
  settings,
  login,
  password,
  legal,
  oidc,
};

/* eslint-disable */
import Long from "long";
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";
import { Timestamp } from "../google/protobuf/timestamp";
import { AggregateType, Event, EventType } from "./event";
import {
  IDP,
  IDPFieldName,
  iDPFieldNameFromJSON,
  iDPFieldNameToJSON,
  IDPIDQuery,
  IDPLoginPolicyLink,
  IDPNameQuery,
  IDPStylingType,
  iDPStylingTypeFromJSON,
  iDPStylingTypeToJSON,
  IDPUserLink,
  LDAPAttributes,
  OIDCMappingField,
  oIDCMappingFieldFromJSON,
  oIDCMappingFieldToJSON,
  Options,
  Provider,
} from "./idp";
import {
  Domain,
  DomainFieldName,
  domainFieldNameFromJSON,
  domainFieldNameToJSON,
  DomainSearchQuery,
  InstanceDetail,
} from "./instance";
import {
  AddCustomLabelPolicyRequest,
  AddCustomLockoutPolicyRequest,
  AddCustomLoginPolicyRequest,
  AddCustomPasswordComplexityPolicyRequest,
  AddCustomPrivacyPolicyRequest,
  AddOrgMemberRequest,
  AddOrgRequest,
  AddProjectGrantMemberRequest,
  AddProjectMemberRequest,
  AddProjectRoleRequest,
  AddUserGrantRequest,
  SetCustomDomainClaimedMessageTextRequest,
  SetCustomInitMessageTextRequest,
  SetCustomLoginTextsRequest as SetCustomLoginTextsRequest2,
  SetCustomPasswordlessRegistrationMessageTextRequest,
  SetCustomPasswordResetMessageTextRequest,
  SetCustomVerifyEmailMessageTextRequest,
  SetCustomVerifyPhoneMessageTextRequest,
  SetTriggerActionsRequest,
  SetUserMetadataRequest,
} from "./management";
import { Member, SearchQuery } from "./member";
import { ListDetails, ListQuery, ObjectDetails } from "./object";
import { Domain as Domain3, Org, OrgFieldName, orgFieldNameFromJSON, orgFieldNameToJSON, OrgQuery } from "./org";
import {
  DomainPolicy,
  LabelPolicy,
  LockoutPolicy,
  LoginPolicy,
  MultiFactorType,
  multiFactorTypeFromJSON,
  multiFactorTypeToJSON,
  NotificationPolicy,
  OrgIAMPolicy,
  PasswordAgePolicy,
  PasswordComplexityPolicy,
  PasswordlessType,
  passwordlessTypeFromJSON,
  passwordlessTypeToJSON,
  PrivacyPolicy,
  SecondFactorType,
  secondFactorTypeFromJSON,
  secondFactorTypeToJSON,
} from "./policy";
import {
  DebugNotificationProvider,
  OIDCSettings,
  SecretGenerator,
  SecretGeneratorQuery,
  SecretGeneratorType,
  secretGeneratorTypeFromJSON,
  secretGeneratorTypeToJSON,
  SecurityPolicy,
  SMSProvider,
  SMTPConfig,
} from "./settings";
import {
  EmailVerificationDoneScreenText,
  EmailVerificationScreenText,
  ExternalRegistrationUserOverviewScreenText,
  ExternalUserNotFoundScreenText,
  FooterText,
  InitializeUserDoneScreenText,
  InitializeUserScreenText,
  InitMFADoneScreenText,
  InitMFAOTPScreenText,
  InitMFAPromptScreenText,
  InitMFAU2FScreenText,
  InitPasswordDoneScreenText,
  InitPasswordScreenText,
  LinkingUserDoneScreenText,
  LoginCustomText,
  LoginScreenText,
  LogoutDoneScreenText,
  MessageCustomText,
  MFAProvidersText,
  PasswordChangeDoneScreenText,
  PasswordChangeScreenText,
  PasswordlessPromptScreenText,
  PasswordlessRegistrationDoneScreenText,
  PasswordlessRegistrationScreenText,
  PasswordlessScreenText,
  PasswordResetDoneScreenText,
  PasswordScreenText,
  RegistrationOptionScreenText,
  RegistrationOrgScreenText,
  RegistrationUserScreenText,
  SelectAccountScreenText,
  SuccessLoginScreenText,
  UsernameChangeDoneScreenText,
  UsernameChangeScreenText,
  VerifyMFAOTPScreenText,
  VerifyMFAU2FScreenText,
} from "./text";
import { Gender, genderFromJSON, genderToJSON } from "./user";
import {
  DataAction,
  DataAPIApplication,
  DataAppKey,
  DataHumanUser,
  DataJWTIDP,
  DataMachineKey,
  DataMachineUser,
  DataOIDCApplication,
  DataOIDCIDP,
  DataProject,
  DataProjectGrant,
  ImportDataOrg as ImportDataOrg1,
} from "./v1";

export const protobufPackage = "zitadel.admin.v1";

/** This is an empty request */
export interface HealthzRequest {
}

/** This is an empty response */
export interface HealthzResponse {
}

/** This is an empty request */
export interface GetSupportedLanguagesRequest {
}

export interface GetSupportedLanguagesResponse {
  languages: string[];
}

export interface SetDefaultLanguageRequest {
  language: string;
}

export interface SetDefaultLanguageResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetDefaultLanguageRequest {
}

export interface GetDefaultLanguageResponse {
  language: string;
}

export interface SetDefaultOrgRequest {
  orgId: string;
}

export interface SetDefaultOrgResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetDefaultOrgRequest {
}

export interface GetDefaultOrgResponse {
  org: Org | undefined;
}

/** This is an empty request */
export interface GetMyInstanceRequest {
}

export interface GetMyInstanceResponse {
  instance: InstanceDetail | undefined;
}

export interface ListInstanceDomainsRequest {
  query:
    | ListQuery
    | undefined;
  /** the field the result is sorted */
  sortingColumn: DomainFieldName;
  /** criteria the client is looking for */
  queries: DomainSearchQuery[];
}

export interface ListInstanceDomainsResponse {
  details: ListDetails | undefined;
  sortingColumn: DomainFieldName;
  result: Domain[];
}

export interface ListSecretGeneratorsRequest {
  /** list limitations and ordering */
  query:
    | ListQuery
    | undefined;
  /** criteria the client is looking for */
  queries: SecretGeneratorQuery[];
}

export interface ListSecretGeneratorsResponse {
  details: ListDetails | undefined;
  result: SecretGenerator[];
}

export interface GetSecretGeneratorRequest {
  generatorType: SecretGeneratorType;
}

export interface GetSecretGeneratorResponse {
  secretGenerator: SecretGenerator | undefined;
}

export interface UpdateSecretGeneratorRequest {
  generatorType: SecretGeneratorType;
  length: number;
  expiry: Duration | undefined;
  includeLowerLetters: boolean;
  includeUpperLetters: boolean;
  includeDigits: boolean;
  includeSymbols: boolean;
}

export interface UpdateSecretGeneratorResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetSMTPConfigRequest {
}

export interface GetSMTPConfigResponse {
  smtpConfig: SMTPConfig | undefined;
}

export interface AddSMTPConfigRequest {
  senderAddress: string;
  senderName: string;
  tls: boolean;
  host: string;
  user: string;
  password: string;
}

export interface AddSMTPConfigResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateSMTPConfigRequest {
  senderAddress: string;
  senderName: string;
  tls: boolean;
  host: string;
  user: string;
}

export interface UpdateSMTPConfigResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateSMTPConfigPasswordRequest {
  password: string;
}

export interface UpdateSMTPConfigPasswordResponse {
  details: ObjectDetails | undefined;
}

/** this is an empty request */
export interface RemoveSMTPConfigRequest {
}

export interface RemoveSMTPConfigResponse {
  details: ObjectDetails | undefined;
}

export interface ListSMSProvidersRequest {
  /** list limitations and ordering */
  query: ListQuery | undefined;
}

export interface ListSMSProvidersResponse {
  details: ListDetails | undefined;
  result: SMSProvider[];
}

export interface GetSMSProviderRequest {
  id: string;
}

export interface GetSMSProviderResponse {
  config: SMSProvider | undefined;
}

export interface AddSMSProviderTwilioRequest {
  sid: string;
  token: string;
  senderNumber: string;
}

export interface AddSMSProviderTwilioResponse {
  details: ObjectDetails | undefined;
  id: string;
}

export interface UpdateSMSProviderTwilioRequest {
  id: string;
  sid: string;
  senderNumber: string;
}

export interface UpdateSMSProviderTwilioResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateSMSProviderTwilioTokenRequest {
  id: string;
  token: string;
}

export interface UpdateSMSProviderTwilioTokenResponse {
  details: ObjectDetails | undefined;
}

export interface ActivateSMSProviderRequest {
  id: string;
}

export interface ActivateSMSProviderResponse {
  details: ObjectDetails | undefined;
}

export interface DeactivateSMSProviderRequest {
  id: string;
}

export interface DeactivateSMSProviderResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveSMSProviderRequest {
  id: string;
}

export interface RemoveSMSProviderResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetFileSystemNotificationProviderRequest {
}

export interface GetFileSystemNotificationProviderResponse {
  provider: DebugNotificationProvider | undefined;
}

/** This is an empty request */
export interface GetLogNotificationProviderRequest {
}

export interface GetLogNotificationProviderResponse {
  provider: DebugNotificationProvider | undefined;
}

/** This is an empty request */
export interface GetOIDCSettingsRequest {
}

export interface GetOIDCSettingsResponse {
  settings: OIDCSettings | undefined;
}

export interface AddOIDCSettingsRequest {
  accessTokenLifetime: Duration | undefined;
  idTokenLifetime: Duration | undefined;
  refreshTokenIdleExpiration: Duration | undefined;
  refreshTokenExpiration: Duration | undefined;
}

export interface AddOIDCSettingsResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateOIDCSettingsRequest {
  accessTokenLifetime: Duration | undefined;
  idTokenLifetime: Duration | undefined;
  refreshTokenIdleExpiration: Duration | undefined;
  refreshTokenExpiration: Duration | undefined;
}

export interface UpdateOIDCSettingsResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetSecurityPolicyRequest {
}

export interface GetSecurityPolicyResponse {
  policy: SecurityPolicy | undefined;
}

export interface SetSecurityPolicyRequest {
  /** states if iframe embedding is enabled or disabled */
  enableIframeEmbedding: boolean;
  /** origins allowed loading ZITADEL in an iframe if enable_iframe_embedding is true */
  allowedOrigins: string[];
}

export interface SetSecurityPolicyResponse {
  details: ObjectDetails | undefined;
}

/**
 * if name or domain is already in use, org is not unique
 * at least one argument has to be provided
 */
export interface IsOrgUniqueRequest {
  name: string;
  domain: string;
}

export interface IsOrgUniqueResponse {
  isUnique: boolean;
}

export interface GetOrgByIDRequest {
  id: string;
}

export interface GetOrgByIDResponse {
  org: Org | undefined;
}

export interface ListOrgsRequest {
  /** list limitations and ordering */
  query:
    | ListQuery
    | undefined;
  /** the field the result is sorted */
  sortingColumn: OrgFieldName;
  /** criteria the client is looking for */
  queries: OrgQuery[];
}

export interface ListOrgsResponse {
  details: ListDetails | undefined;
  sortingColumn: OrgFieldName;
  result: Org[];
}

export interface SetUpOrgRequest {
  org:
    | SetUpOrgRequest_Org
    | undefined;
  /** oneof field for the user managing the organization */
  human?:
    | SetUpOrgRequest_Human
    | undefined;
  /** specify Org Member Roles for the provided user (default is ORG_OWNER if roles are empty) */
  roles: string[];
}

export interface SetUpOrgRequest_Org {
  name: string;
  domain: string;
}

export interface SetUpOrgRequest_Human {
  userName: string;
  profile: SetUpOrgRequest_Human_Profile | undefined;
  email: SetUpOrgRequest_Human_Email | undefined;
  phone: SetUpOrgRequest_Human_Phone | undefined;
  password: string;
}

export interface SetUpOrgRequest_Human_Profile {
  firstName: string;
  lastName: string;
  nickName: string;
  displayName: string;
  preferredLanguage: string;
  gender: Gender;
}

export interface SetUpOrgRequest_Human_Email {
  email: string;
  isEmailVerified: boolean;
}

export interface SetUpOrgRequest_Human_Phone {
  /** has to be a global number */
  phone: string;
  isPhoneVerified: boolean;
}

export interface SetUpOrgResponse {
  details: ObjectDetails | undefined;
  orgId: string;
  userId: string;
}

export interface RemoveOrgRequest {
  orgId: string;
}

export interface RemoveOrgResponse {
  details: ObjectDetails | undefined;
}

export interface GetIDPByIDRequest {
  id: string;
}

export interface GetIDPByIDResponse {
  idp: IDP | undefined;
}

export interface ListIDPsRequest {
  /** list limitations and ordering */
  query:
    | ListQuery
    | undefined;
  /** the field the result is sorted */
  sortingColumn: IDPFieldName;
  /** criteria the client is looking for */
  queries: IDPQuery[];
}

export interface IDPQuery {
  idpIdQuery?: IDPIDQuery | undefined;
  idpNameQuery?: IDPNameQuery | undefined;
}

export interface ListIDPsResponse {
  details: ListDetails | undefined;
  sortingColumn: IDPFieldName;
  result: IDP[];
}

export interface AddOIDCIDPRequest {
  name: string;
  stylingType: IDPStylingType;
  clientId: string;
  clientSecret: string;
  issuer: string;
  scopes: string[];
  displayNameMapping: OIDCMappingField;
  usernameMapping: OIDCMappingField;
  autoRegister: boolean;
}

export interface AddOIDCIDPResponse {
  details: ObjectDetails | undefined;
  idpId: string;
}

export interface AddJWTIDPRequest {
  name: string;
  stylingType: IDPStylingType;
  jwtEndpoint: string;
  issuer: string;
  keysEndpoint: string;
  headerName: string;
  autoRegister: boolean;
}

export interface AddJWTIDPResponse {
  details: ObjectDetails | undefined;
  idpId: string;
}

export interface UpdateIDPRequest {
  idpId: string;
  name: string;
  stylingType: IDPStylingType;
  autoRegister: boolean;
}

export interface UpdateIDPResponse {
  details: ObjectDetails | undefined;
}

export interface DeactivateIDPRequest {
  idpId: string;
}

export interface DeactivateIDPResponse {
  details: ObjectDetails | undefined;
}

export interface ReactivateIDPRequest {
  idpId: string;
}

export interface ReactivateIDPResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveIDPRequest {
  idpId: string;
}

export interface RemoveIDPResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateIDPOIDCConfigRequest {
  idpId: string;
  issuer: string;
  clientId: string;
  clientSecret: string;
  scopes: string[];
  displayNameMapping: OIDCMappingField;
  usernameMapping: OIDCMappingField;
}

export interface UpdateIDPOIDCConfigResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateIDPJWTConfigRequest {
  idpId: string;
  jwtEndpoint: string;
  issuer: string;
  keysEndpoint: string;
  headerName: string;
}

export interface UpdateIDPJWTConfigResponse {
  details: ObjectDetails | undefined;
}

export interface ListProvidersRequest {
  /** list limitations and ordering */
  query:
    | ListQuery
    | undefined;
  /** criteria the client is looking for */
  queries: ProviderQuery[];
}

export interface ProviderQuery {
  idpIdQuery?: IDPIDQuery | undefined;
  idpNameQuery?: IDPNameQuery | undefined;
}

export interface ListProvidersResponse {
  details: ListDetails | undefined;
  result: Provider[];
}

export interface GetProviderByIDRequest {
  id: string;
}

export interface GetProviderByIDResponse {
  idp: Provider | undefined;
}

export interface AddGenericOAuthProviderRequest {
  name: string;
  clientId: string;
  clientSecret: string;
  authorizationEndpoint: string;
  tokenEndpoint: string;
  userEndpoint: string;
  scopes: string[];
  /** identifying attribute of the user in the response of the user_endpoint */
  idAttribute: string;
  providerOptions: Options | undefined;
}

export interface AddGenericOAuthProviderResponse {
  details: ObjectDetails | undefined;
  id: string;
}

export interface UpdateGenericOAuthProviderRequest {
  id: string;
  name: string;
  clientId: string;
  /** client_secret will only be updated if provided */
  clientSecret: string;
  authorizationEndpoint: string;
  tokenEndpoint: string;
  userEndpoint: string;
  scopes: string[];
  /** identifying attribute of the user in the response of the user_endpoint */
  idAttribute: string;
  providerOptions: Options | undefined;
}

export interface UpdateGenericOAuthProviderResponse {
  details: ObjectDetails | undefined;
}

export interface AddGenericOIDCProviderRequest {
  name: string;
  issuer: string;
  clientId: string;
  clientSecret: string;
  scopes: string[];
  providerOptions: Options | undefined;
}

export interface AddGenericOIDCProviderResponse {
  details: ObjectDetails | undefined;
  id: string;
}

export interface UpdateGenericOIDCProviderRequest {
  id: string;
  name: string;
  issuer: string;
  clientId: string;
  /** client_secret will only be updated if provided */
  clientSecret: string;
  scopes: string[];
  providerOptions: Options | undefined;
}

export interface UpdateGenericOIDCProviderResponse {
  details: ObjectDetails | undefined;
}

export interface AddJWTProviderRequest {
  name: string;
  issuer: string;
  jwtEndpoint: string;
  keysEndpoint: string;
  headerName: string;
  providerOptions: Options | undefined;
}

export interface AddJWTProviderResponse {
  details: ObjectDetails | undefined;
  id: string;
}

export interface UpdateJWTProviderRequest {
  id: string;
  name: string;
  issuer: string;
  jwtEndpoint: string;
  keysEndpoint: string;
  headerName: string;
  providerOptions: Options | undefined;
}

export interface UpdateJWTProviderResponse {
  details: ObjectDetails | undefined;
}

export interface AddGitHubProviderRequest {
  /** GitHub will be used as default, if no name is provided */
  name: string;
  clientId: string;
  clientSecret: string;
  scopes: string[];
  providerOptions: Options | undefined;
}

export interface AddGitHubProviderResponse {
  details: ObjectDetails | undefined;
  id: string;
}

export interface UpdateGitHubProviderRequest {
  id: string;
  name: string;
  clientId: string;
  /** client_secret will only be updated if provided */
  clientSecret: string;
  scopes: string[];
  providerOptions: Options | undefined;
}

export interface UpdateGitHubProviderResponse {
  details: ObjectDetails | undefined;
}

export interface AddGitHubEnterpriseServerProviderRequest {
  clientId: string;
  name: string;
  clientSecret: string;
  authorizationEndpoint: string;
  tokenEndpoint: string;
  userEndpoint: string;
  scopes: string[];
  providerOptions: Options | undefined;
}

export interface AddGitHubEnterpriseServerProviderResponse {
  details: ObjectDetails | undefined;
  id: string;
}

export interface UpdateGitHubEnterpriseServerProviderRequest {
  id: string;
  name: string;
  clientId: string;
  /** client_secret will only be updated if provided */
  clientSecret: string;
  authorizationEndpoint: string;
  tokenEndpoint: string;
  userEndpoint: string;
  scopes: string[];
  providerOptions: Options | undefined;
}

export interface UpdateGitHubEnterpriseServerProviderResponse {
  details: ObjectDetails | undefined;
}

export interface AddGoogleProviderRequest {
  /** Google will be used as default, if no name is provided */
  name: string;
  clientId: string;
  clientSecret: string;
  scopes: string[];
  providerOptions: Options | undefined;
}

export interface AddGoogleProviderResponse {
  details: ObjectDetails | undefined;
  id: string;
}

export interface UpdateGoogleProviderRequest {
  id: string;
  name: string;
  clientId: string;
  /** client_secret will only be updated if provided */
  clientSecret: string;
  scopes: string[];
  providerOptions: Options | undefined;
}

export interface UpdateGoogleProviderResponse {
  details: ObjectDetails | undefined;
}

export interface AddLDAPProviderRequest {
  name: string;
  host: string;
  port: string;
  tls: boolean;
  baseDn: string;
  userObjectClass: string;
  userUniqueAttribute: string;
  admin: string;
  password: string;
  attributes: LDAPAttributes | undefined;
  providerOptions: Options | undefined;
}

export interface AddLDAPProviderResponse {
  details: ObjectDetails | undefined;
  id: string;
}

export interface UpdateLDAPProviderRequest {
  id: string;
  name: string;
  host: string;
  port: string;
  tls: boolean;
  baseDn: string;
  userObjectClass: string;
  userUniqueAttribute: string;
  admin: string;
  password: string;
  attributes: LDAPAttributes | undefined;
  providerOptions: Options | undefined;
}

export interface UpdateLDAPProviderResponse {
  details: ObjectDetails | undefined;
}

export interface DeleteProviderRequest {
  id: string;
}

export interface DeleteProviderResponse {
  details: ObjectDetails | undefined;
}

export interface GetOrgIAMPolicyRequest {
}

export interface GetOrgIAMPolicyResponse {
  policy: OrgIAMPolicy | undefined;
}

export interface UpdateOrgIAMPolicyRequest {
  userLoginMustBeDomain: boolean;
}

export interface UpdateOrgIAMPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface GetCustomOrgIAMPolicyRequest {
  orgId: string;
}

export interface GetCustomOrgIAMPolicyResponse {
  policy:
    | OrgIAMPolicy
    | undefined;
  /** deprecated: is_default is also defined in zitadel.policy.v1.OrgIAMPolicy */
  isDefault: boolean;
}

export interface AddCustomOrgIAMPolicyRequest {
  orgId: string;
  /** the username has to end with the domain of its organization (uniqueness is organization based) */
  userLoginMustBeDomain: boolean;
}

export interface AddCustomOrgIAMPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateCustomOrgIAMPolicyRequest {
  orgId: string;
  userLoginMustBeDomain: boolean;
}

export interface UpdateCustomOrgIAMPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomOrgIAMPolicyToDefaultRequest {
  orgId: string;
}

export interface ResetCustomOrgIAMPolicyToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface GetDomainPolicyRequest {
}

export interface GetDomainPolicyResponse {
  policy: DomainPolicy | undefined;
}

export interface UpdateDomainPolicyRequest {
  userLoginMustBeDomain: boolean;
  validateOrgDomains: boolean;
  smtpSenderAddressMatchesInstanceDomain: boolean;
}

export interface UpdateDomainPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface GetCustomDomainPolicyRequest {
  orgId: string;
}

export interface GetCustomDomainPolicyResponse {
  policy:
    | DomainPolicy
    | undefined;
  /** deprecated: is_default is also defined in zitadel.policy.v1.DomainPolicy */
  isDefault: boolean;
}

export interface AddCustomDomainPolicyRequest {
  orgId: string;
  /** the username has to end with the domain of its organization (uniqueness is organization based) */
  userLoginMustBeDomain: boolean;
  validateOrgDomains: boolean;
  smtpSenderAddressMatchesInstanceDomain: boolean;
}

export interface AddCustomDomainPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateCustomDomainPolicyRequest {
  orgId: string;
  userLoginMustBeDomain: boolean;
  validateOrgDomains: boolean;
  smtpSenderAddressMatchesInstanceDomain: boolean;
}

export interface UpdateCustomDomainPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomDomainPolicyToDefaultRequest {
  orgId: string;
}

export interface ResetCustomDomainPolicyToDefaultResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetLabelPolicyRequest {
}

export interface GetLabelPolicyResponse {
  policy: LabelPolicy | undefined;
}

/** This is an empty request */
export interface GetPreviewLabelPolicyRequest {
}

export interface GetPreviewLabelPolicyResponse {
  policy: LabelPolicy | undefined;
}

export interface UpdateLabelPolicyRequest {
  primaryColor: string;
  hideLoginNameSuffix: boolean;
  warnColor: string;
  backgroundColor: string;
  fontColor: string;
  primaryColorDark: string;
  backgroundColorDark: string;
  warnColorDark: string;
  fontColorDark: string;
  disableWatermark: boolean;
}

export interface UpdateLabelPolicyResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ActivateLabelPolicyRequest {
}

export interface ActivateLabelPolicyResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RemoveLabelPolicyLogoRequest {
}

export interface RemoveLabelPolicyLogoResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RemoveLabelPolicyLogoDarkRequest {
}

export interface RemoveLabelPolicyLogoDarkResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RemoveLabelPolicyIconRequest {
}

export interface RemoveLabelPolicyIconResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RemoveLabelPolicyIconDarkRequest {
}

export interface RemoveLabelPolicyIconDarkResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RemoveLabelPolicyFontRequest {
}

export interface RemoveLabelPolicyFontResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetLoginPolicyRequest {
}

export interface GetLoginPolicyResponse {
  policy: LoginPolicy | undefined;
}

export interface UpdateLoginPolicyRequest {
  allowUsernamePassword: boolean;
  allowRegister: boolean;
  allowExternalIdp: boolean;
  forceMfa: boolean;
  passwordlessType: PasswordlessType;
  hidePasswordReset: boolean;
  ignoreUnknownUsernames: boolean;
  defaultRedirectUri: string;
  passwordCheckLifetime: Duration | undefined;
  externalLoginCheckLifetime: Duration | undefined;
  mfaInitSkipLifetime: Duration | undefined;
  secondFactorCheckLifetime: Duration | undefined;
  multiFactorCheckLifetime:
    | Duration
    | undefined;
  /** If set to true, the suffix (@domain.com) of an unknown username input on the login screen will be matched against the org domains and will redirect to the registration of that organization on success. */
  allowDomainDiscovery: boolean;
  disableLoginWithEmail: boolean;
  disableLoginWithPhone: boolean;
}

export interface UpdateLoginPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface ListLoginPolicyIDPsRequest {
  /** list limitations and ordering */
  query: ListQuery | undefined;
}

export interface ListLoginPolicyIDPsResponse {
  details: ListDetails | undefined;
  result: IDPLoginPolicyLink[];
}

export interface AddIDPToLoginPolicyRequest {
  idpId: string;
}

export interface AddIDPToLoginPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveIDPFromLoginPolicyRequest {
  idpId: string;
}

export interface RemoveIDPFromLoginPolicyResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ListLoginPolicySecondFactorsRequest {
}

export interface ListLoginPolicySecondFactorsResponse {
  details: ListDetails | undefined;
  result: SecondFactorType[];
}

export interface AddSecondFactorToLoginPolicyRequest {
  type: SecondFactorType;
}

export interface AddSecondFactorToLoginPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveSecondFactorFromLoginPolicyRequest {
  type: SecondFactorType;
}

export interface RemoveSecondFactorFromLoginPolicyResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ListLoginPolicyMultiFactorsRequest {
}

export interface ListLoginPolicyMultiFactorsResponse {
  details: ListDetails | undefined;
  result: MultiFactorType[];
}

export interface AddMultiFactorToLoginPolicyRequest {
  type: MultiFactorType;
}

export interface AddMultiFactorToLoginPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveMultiFactorFromLoginPolicyRequest {
  type: MultiFactorType;
}

export interface RemoveMultiFactorFromLoginPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface GetPasswordComplexityPolicyRequest {
}

export interface GetPasswordComplexityPolicyResponse {
  policy: PasswordComplexityPolicy | undefined;
}

export interface UpdatePasswordComplexityPolicyRequest {
  minLength: number;
  hasUppercase: boolean;
  hasLowercase: boolean;
  hasNumber: boolean;
  hasSymbol: boolean;
}

export interface UpdatePasswordComplexityPolicyResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetPasswordAgePolicyRequest {
}

export interface GetPasswordAgePolicyResponse {
  policy: PasswordAgePolicy | undefined;
}

export interface UpdatePasswordAgePolicyRequest {
  maxAgeDays: number;
  expireWarnDays: number;
}

export interface UpdatePasswordAgePolicyResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetLockoutPolicyRequest {
}

export interface GetLockoutPolicyResponse {
  policy: LockoutPolicy | undefined;
}

export interface UpdateLockoutPolicyRequest {
  /** failed attempts until a user gets locked */
  maxPasswordAttempts: number;
}

export interface UpdateLockoutPolicyResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetPrivacyPolicyRequest {
}

export interface GetPrivacyPolicyResponse {
  policy: PrivacyPolicy | undefined;
}

export interface UpdatePrivacyPolicyRequest {
  tosLink: string;
  privacyLink: string;
  helpLink: string;
}

export interface UpdatePrivacyPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface AddNotificationPolicyRequest {
  passwordChange: boolean;
}

export interface AddNotificationPolicyResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetNotificationPolicyRequest {
}

export interface GetNotificationPolicyResponse {
  policy: NotificationPolicy | undefined;
}

export interface UpdateNotificationPolicyRequest {
  passwordChange: boolean;
}

export interface UpdateNotificationPolicyResponse {
  details: ObjectDetails | undefined;
}

export interface GetDefaultInitMessageTextRequest {
  language: string;
}

export interface GetDefaultInitMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface GetCustomInitMessageTextRequest {
  language: string;
}

export interface GetCustomInitMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface SetDefaultInitMessageTextRequest {
  language: string;
  title: string;
  preHeader: string;
  subject: string;
  greeting: string;
  text: string;
  buttonText: string;
  footerText: string;
}

export interface SetDefaultInitMessageTextResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomInitMessageTextToDefaultRequest {
  language: string;
}

export interface ResetCustomInitMessageTextToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface GetDefaultPasswordResetMessageTextRequest {
  language: string;
}

export interface GetDefaultPasswordResetMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface GetCustomPasswordResetMessageTextRequest {
  language: string;
}

export interface GetCustomPasswordResetMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface SetDefaultPasswordResetMessageTextRequest {
  language: string;
  title: string;
  preHeader: string;
  subject: string;
  greeting: string;
  text: string;
  buttonText: string;
  footerText: string;
}

export interface SetDefaultPasswordResetMessageTextResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomPasswordResetMessageTextToDefaultRequest {
  language: string;
}

export interface ResetCustomPasswordResetMessageTextToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface GetDefaultVerifyEmailMessageTextRequest {
  language: string;
}

export interface GetDefaultVerifyEmailMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface GetCustomVerifyEmailMessageTextRequest {
  language: string;
}

export interface GetCustomVerifyEmailMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface SetDefaultVerifyEmailMessageTextRequest {
  language: string;
  title: string;
  preHeader: string;
  subject: string;
  greeting: string;
  text: string;
  buttonText: string;
  footerText: string;
}

export interface SetDefaultVerifyEmailMessageTextResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomVerifyEmailMessageTextToDefaultRequest {
  language: string;
}

export interface ResetCustomVerifyEmailMessageTextToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface GetDefaultVerifyPhoneMessageTextRequest {
  language: string;
}

export interface GetDefaultVerifyPhoneMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface GetCustomVerifyPhoneMessageTextRequest {
  language: string;
}

export interface GetCustomVerifyPhoneMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface SetDefaultVerifyPhoneMessageTextRequest {
  language: string;
  title: string;
  preHeader: string;
  subject: string;
  greeting: string;
  text: string;
  buttonText: string;
  footerText: string;
}

export interface SetDefaultVerifyPhoneMessageTextResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomVerifyPhoneMessageTextToDefaultRequest {
  language: string;
}

export interface ResetCustomVerifyPhoneMessageTextToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface GetDefaultDomainClaimedMessageTextRequest {
  language: string;
}

export interface GetDefaultDomainClaimedMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface GetCustomDomainClaimedMessageTextRequest {
  language: string;
}

export interface GetCustomDomainClaimedMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface SetDefaultDomainClaimedMessageTextRequest {
  language: string;
  title: string;
  preHeader: string;
  subject: string;
  greeting: string;
  text: string;
  buttonText: string;
  footerText: string;
}

export interface SetDefaultDomainClaimedMessageTextResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomDomainClaimedMessageTextToDefaultRequest {
  language: string;
}

export interface ResetCustomDomainClaimedMessageTextToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface GetDefaultPasswordChangeMessageTextRequest {
  language: string;
}

export interface GetDefaultPasswordChangeMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface GetCustomPasswordChangeMessageTextRequest {
  language: string;
}

export interface GetCustomPasswordChangeMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface SetDefaultPasswordChangeMessageTextRequest {
  language: string;
  title: string;
  preHeader: string;
  subject: string;
  greeting: string;
  text: string;
  buttonText: string;
  footerText: string;
}

export interface SetDefaultPasswordChangeMessageTextResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomPasswordChangeMessageTextToDefaultRequest {
  language: string;
}

export interface ResetCustomPasswordChangeMessageTextToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface GetDefaultPasswordlessRegistrationMessageTextRequest {
  language: string;
}

export interface GetDefaultPasswordlessRegistrationMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface GetCustomPasswordlessRegistrationMessageTextRequest {
  language: string;
}

export interface GetCustomPasswordlessRegistrationMessageTextResponse {
  customText: MessageCustomText | undefined;
}

export interface SetDefaultPasswordlessRegistrationMessageTextRequest {
  language: string;
  title: string;
  preHeader: string;
  subject: string;
  greeting: string;
  text: string;
  buttonText: string;
  footerText: string;
}

export interface SetDefaultPasswordlessRegistrationMessageTextResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest {
  language: string;
}

export interface ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface GetDefaultLoginTextsRequest {
  language: string;
}

export interface GetDefaultLoginTextsResponse {
  customText: LoginCustomText | undefined;
}

export interface GetCustomLoginTextsRequest {
  language: string;
}

export interface GetCustomLoginTextsResponse {
  customText: LoginCustomText | undefined;
}

export interface SetCustomLoginTextsRequest {
  language: string;
  selectAccountText: SelectAccountScreenText | undefined;
  loginText: LoginScreenText | undefined;
  passwordText: PasswordScreenText | undefined;
  usernameChangeText: UsernameChangeScreenText | undefined;
  usernameChangeDoneText: UsernameChangeDoneScreenText | undefined;
  initPasswordText: InitPasswordScreenText | undefined;
  initPasswordDoneText: InitPasswordDoneScreenText | undefined;
  emailVerificationText: EmailVerificationScreenText | undefined;
  emailVerificationDoneText: EmailVerificationDoneScreenText | undefined;
  initializeUserText: InitializeUserScreenText | undefined;
  initializeDoneText: InitializeUserDoneScreenText | undefined;
  initMfaPromptText: InitMFAPromptScreenText | undefined;
  initMfaOtpText: InitMFAOTPScreenText | undefined;
  initMfaU2fText: InitMFAU2FScreenText | undefined;
  initMfaDoneText: InitMFADoneScreenText | undefined;
  mfaProvidersText: MFAProvidersText | undefined;
  verifyMfaOtpText: VerifyMFAOTPScreenText | undefined;
  verifyMfaU2fText: VerifyMFAU2FScreenText | undefined;
  passwordlessText: PasswordlessScreenText | undefined;
  passwordChangeText: PasswordChangeScreenText | undefined;
  passwordChangeDoneText: PasswordChangeDoneScreenText | undefined;
  passwordResetDoneText: PasswordResetDoneScreenText | undefined;
  registrationOptionText: RegistrationOptionScreenText | undefined;
  registrationUserText: RegistrationUserScreenText | undefined;
  registrationOrgText: RegistrationOrgScreenText | undefined;
  linkingUserDoneText: LinkingUserDoneScreenText | undefined;
  externalUserNotFoundText: ExternalUserNotFoundScreenText | undefined;
  successLoginText: SuccessLoginScreenText | undefined;
  logoutText: LogoutDoneScreenText | undefined;
  footerText: FooterText | undefined;
  passwordlessPromptText: PasswordlessPromptScreenText | undefined;
  passwordlessRegistrationText: PasswordlessRegistrationScreenText | undefined;
  passwordlessRegistrationDoneText: PasswordlessRegistrationDoneScreenText | undefined;
  externalRegistrationUserOverviewText: ExternalRegistrationUserOverviewScreenText | undefined;
}

export interface SetCustomLoginTextsResponse {
  details: ObjectDetails | undefined;
}

export interface ResetCustomLoginTextsToDefaultRequest {
  language: string;
}

export interface ResetCustomLoginTextsToDefaultResponse {
  details: ObjectDetails | undefined;
}

export interface AddIAMMemberRequest {
  userId: string;
  roles: string[];
}

export interface AddIAMMemberResponse {
  details: ObjectDetails | undefined;
}

export interface UpdateIAMMemberRequest {
  userId: string;
  roles: string[];
}

export interface UpdateIAMMemberResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveIAMMemberRequest {
  userId: string;
}

export interface RemoveIAMMemberResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ListIAMMemberRolesRequest {
}

export interface ListIAMMemberRolesResponse {
  details: ListDetails | undefined;
  roles: string[];
}

export interface ListIAMMembersRequest {
  /** list limitations and ordering */
  query:
    | ListQuery
    | undefined;
  /** criteria the client is looking for */
  queries: SearchQuery[];
}

export interface ListIAMMembersResponse {
  details: ListDetails | undefined;
  result: Member[];
}

/** This is an empty request */
export interface ListViewsRequest {
}

export interface ListViewsResponse {
  /** TODO: list details */
  result: View[];
}

/** This is an empty request */
export interface ListFailedEventsRequest {
}

export interface ListFailedEventsResponse {
  /** TODO: list details */
  result: FailedEvent[];
}

export interface RemoveFailedEventRequest {
  database: string;
  viewName: string;
  failedSequence: number;
}

/** This is an empty response */
export interface RemoveFailedEventResponse {
}

export interface View {
  database: string;
  viewName: string;
  processedSequence: number;
  /** The timestamp the event occurred */
  eventTimestamp: Date | undefined;
  lastSuccessfulSpoolerRun: Date | undefined;
}

export interface FailedEvent {
  database: string;
  viewName: string;
  failedSequence: number;
  failureCount: number;
  errorMessage: string;
  lastFailed: Date | undefined;
}

export interface ImportDataRequest {
  dataOrgs?: ImportDataOrg | undefined;
  dataOrgsv1?: ImportDataOrg1 | undefined;
  dataOrgsLocal?: ImportDataRequest_LocalInput | undefined;
  dataOrgsv1Local?: ImportDataRequest_LocalInput | undefined;
  dataOrgsS3?: ImportDataRequest_S3Input | undefined;
  dataOrgsv1S3?: ImportDataRequest_S3Input | undefined;
  dataOrgsGcs?: ImportDataRequest_GCSInput | undefined;
  dataOrgsv1Gcs?: ImportDataRequest_GCSInput | undefined;
  timeout: string;
}

export interface ImportDataRequest_LocalInput {
  path: string;
}

export interface ImportDataRequest_S3Input {
  path: string;
  endpoint: string;
  accessKeyId: string;
  secretAccessKey: string;
  ssl: boolean;
  bucket: string;
}

export interface ImportDataRequest_GCSInput {
  bucket: string;
  serviceaccountJson: string;
  path: string;
}

export interface ImportDataOrg {
  orgs: DataOrg[];
}

export interface DataOrg {
  orgId: string;
  org: AddOrgRequest | undefined;
  domainPolicy: AddCustomDomainPolicyRequest | undefined;
  labelPolicy: AddCustomLabelPolicyRequest | undefined;
  lockoutPolicy: AddCustomLockoutPolicyRequest | undefined;
  loginPolicy: AddCustomLoginPolicyRequest | undefined;
  passwordComplexityPolicy: AddCustomPasswordComplexityPolicyRequest | undefined;
  privacyPolicy: AddCustomPrivacyPolicyRequest | undefined;
  projects: DataProject[];
  projectRoles: AddProjectRoleRequest[];
  apiApps: DataAPIApplication[];
  oidcApps: DataOIDCApplication[];
  humanUsers: DataHumanUser[];
  machineUsers: DataMachineUser[];
  triggerActions: SetTriggerActionsRequest[];
  actions: DataAction[];
  projectGrants: DataProjectGrant[];
  userGrants: AddUserGrantRequest[];
  orgMembers: AddOrgMemberRequest[];
  projectMembers: AddProjectMemberRequest[];
  projectGrantMembers: AddProjectGrantMemberRequest[];
  userMetadata: SetUserMetadataRequest[];
  loginTexts: SetCustomLoginTextsRequest2[];
  initMessages: SetCustomInitMessageTextRequest[];
  passwordResetMessages: SetCustomPasswordResetMessageTextRequest[];
  verifyEmailMessages: SetCustomVerifyEmailMessageTextRequest[];
  verifyPhoneMessages: SetCustomVerifyPhoneMessageTextRequest[];
  domainClaimedMessages: SetCustomDomainClaimedMessageTextRequest[];
  passwordlessRegistrationMessages: SetCustomPasswordlessRegistrationMessageTextRequest[];
  oidcIdps: DataOIDCIDP[];
  jwtIdps: DataJWTIDP[];
  userLinks: IDPUserLink[];
  domains: Domain3[];
  appKeys: DataAppKey[];
  machineKeys: DataMachineKey[];
}

export interface ImportDataResponse {
  errors: ImportDataError[];
  success: ImportDataSuccess | undefined;
}

export interface ImportDataError {
  type: string;
  id: string;
  message: string;
}

export interface ImportDataSuccess {
  orgs: ImportDataSuccessOrg[];
}

export interface ImportDataSuccessOrg {
  orgId: string;
  projectIds: string[];
  projectRoles: string[];
  oidcAppIds: string[];
  apiAppIds: string[];
  humanUserIds: string[];
  machineUserIds: string[];
  actionIds: string[];
  triggerActions: SetTriggerActionsRequest[];
  projectGrants: ImportDataSuccessProjectGrant[];
  userGrants: ImportDataSuccessUserGrant[];
  orgMembers: string[];
  projectMembers: ImportDataSuccessProjectMember[];
  projectGrantMembers: ImportDataSuccessProjectGrantMember[];
  oidcIpds: string[];
  jwtIdps: string[];
  idpLinks: string[];
  userLinks: ImportDataSuccessUserLinks[];
  userMetadata: ImportDataSuccessUserMetadata[];
  domains: string[];
  appKeys: string[];
  machineKeys: string[];
}

export interface ImportDataSuccessProjectGrant {
  grantId: string;
  projectId: string;
  orgId: string;
}

export interface ImportDataSuccessUserGrant {
  projectId: string;
  userId: string;
}

export interface ImportDataSuccessProjectMember {
  projectId: string;
  userId: string;
}

export interface ImportDataSuccessProjectGrantMember {
  projectId: string;
  grantId: string;
  userId: string;
}

export interface ImportDataSuccessUserLinks {
  userId: string;
  externalUserId: string;
  displayName: string;
  idpId: string;
}

export interface ImportDataSuccessUserMetadata {
  userId: string;
  key: string;
}

export interface ExportDataRequest {
  orgIds: string[];
  excludedOrgIds: string[];
  withPasswords: boolean;
  withOtp: boolean;
  responseOutput: boolean;
  localOutput: ExportDataRequest_LocalOutput | undefined;
  s3Output: ExportDataRequest_S3Output | undefined;
  gcsOutput: ExportDataRequest_GCSOutput | undefined;
  timeout: string;
}

export interface ExportDataRequest_LocalOutput {
  path: string;
}

export interface ExportDataRequest_S3Output {
  path: string;
  endpoint: string;
  accessKeyId: string;
  secretAccessKey: string;
  ssl: boolean;
  bucket: string;
}

export interface ExportDataRequest_GCSOutput {
  bucket: string;
  serviceaccountJson: string;
  path: string;
}

export interface ExportDataResponse {
  orgs: DataOrg[];
}

export interface ListEventsRequest {
  sequence: number;
  limit: number;
  asc: boolean;
  editorUserId: string;
  eventTypes: string[];
  aggregateId: string;
  aggregateTypes: string[];
  resourceOwner: string;
  creationDate: Date | undefined;
}

export interface ListEventsResponse {
  events: Event[];
}

export interface ListEventTypesRequest {
}

export interface ListEventTypesResponse {
  eventTypes: EventType[];
}

export interface ListAggregateTypesRequest {
}

export interface ListAggregateTypesResponse {
  aggregateTypes: AggregateType[];
}

function createBaseHealthzRequest(): HealthzRequest {
  return {};
}

export const HealthzRequest = {
  encode(_: HealthzRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HealthzRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHealthzRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): HealthzRequest {
    return {};
  },

  toJSON(_: HealthzRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<HealthzRequest>): HealthzRequest {
    return HealthzRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<HealthzRequest>): HealthzRequest {
    const message = createBaseHealthzRequest();
    return message;
  },
};

function createBaseHealthzResponse(): HealthzResponse {
  return {};
}

export const HealthzResponse = {
  encode(_: HealthzResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HealthzResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHealthzResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): HealthzResponse {
    return {};
  },

  toJSON(_: HealthzResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<HealthzResponse>): HealthzResponse {
    return HealthzResponse.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<HealthzResponse>): HealthzResponse {
    const message = createBaseHealthzResponse();
    return message;
  },
};

function createBaseGetSupportedLanguagesRequest(): GetSupportedLanguagesRequest {
  return {};
}

export const GetSupportedLanguagesRequest = {
  encode(_: GetSupportedLanguagesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSupportedLanguagesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSupportedLanguagesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetSupportedLanguagesRequest {
    return {};
  },

  toJSON(_: GetSupportedLanguagesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetSupportedLanguagesRequest>): GetSupportedLanguagesRequest {
    return GetSupportedLanguagesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetSupportedLanguagesRequest>): GetSupportedLanguagesRequest {
    const message = createBaseGetSupportedLanguagesRequest();
    return message;
  },
};

function createBaseGetSupportedLanguagesResponse(): GetSupportedLanguagesResponse {
  return { languages: [] };
}

export const GetSupportedLanguagesResponse = {
  encode(message: GetSupportedLanguagesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.languages) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSupportedLanguagesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSupportedLanguagesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.languages.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetSupportedLanguagesResponse {
    return { languages: Array.isArray(object?.languages) ? object.languages.map((e: any) => String(e)) : [] };
  },

  toJSON(message: GetSupportedLanguagesResponse): unknown {
    const obj: any = {};
    if (message.languages) {
      obj.languages = message.languages.map((e) => e);
    } else {
      obj.languages = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GetSupportedLanguagesResponse>): GetSupportedLanguagesResponse {
    return GetSupportedLanguagesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSupportedLanguagesResponse>): GetSupportedLanguagesResponse {
    const message = createBaseGetSupportedLanguagesResponse();
    message.languages = object.languages?.map((e) => e) || [];
    return message;
  },
};

function createBaseSetDefaultLanguageRequest(): SetDefaultLanguageRequest {
  return { language: "" };
}

export const SetDefaultLanguageRequest = {
  encode(message: SetDefaultLanguageRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultLanguageRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultLanguageRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultLanguageRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: SetDefaultLanguageRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultLanguageRequest>): SetDefaultLanguageRequest {
    return SetDefaultLanguageRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultLanguageRequest>): SetDefaultLanguageRequest {
    const message = createBaseSetDefaultLanguageRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseSetDefaultLanguageResponse(): SetDefaultLanguageResponse {
  return { details: undefined };
}

export const SetDefaultLanguageResponse = {
  encode(message: SetDefaultLanguageResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultLanguageResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultLanguageResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultLanguageResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultLanguageResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultLanguageResponse>): SetDefaultLanguageResponse {
    return SetDefaultLanguageResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultLanguageResponse>): SetDefaultLanguageResponse {
    const message = createBaseSetDefaultLanguageResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultLanguageRequest(): GetDefaultLanguageRequest {
  return {};
}

export const GetDefaultLanguageRequest = {
  encode(_: GetDefaultLanguageRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultLanguageRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultLanguageRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetDefaultLanguageRequest {
    return {};
  },

  toJSON(_: GetDefaultLanguageRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetDefaultLanguageRequest>): GetDefaultLanguageRequest {
    return GetDefaultLanguageRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetDefaultLanguageRequest>): GetDefaultLanguageRequest {
    const message = createBaseGetDefaultLanguageRequest();
    return message;
  },
};

function createBaseGetDefaultLanguageResponse(): GetDefaultLanguageResponse {
  return { language: "" };
}

export const GetDefaultLanguageResponse = {
  encode(message: GetDefaultLanguageResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultLanguageResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultLanguageResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultLanguageResponse {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultLanguageResponse): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultLanguageResponse>): GetDefaultLanguageResponse {
    return GetDefaultLanguageResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultLanguageResponse>): GetDefaultLanguageResponse {
    const message = createBaseGetDefaultLanguageResponse();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseSetDefaultOrgRequest(): SetDefaultOrgRequest {
  return { orgId: "" };
}

export const SetDefaultOrgRequest = {
  encode(message: SetDefaultOrgRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultOrgRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultOrgRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultOrgRequest {
    return { orgId: isSet(object.orgId) ? String(object.orgId) : "" };
  },

  toJSON(message: SetDefaultOrgRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultOrgRequest>): SetDefaultOrgRequest {
    return SetDefaultOrgRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultOrgRequest>): SetDefaultOrgRequest {
    const message = createBaseSetDefaultOrgRequest();
    message.orgId = object.orgId ?? "";
    return message;
  },
};

function createBaseSetDefaultOrgResponse(): SetDefaultOrgResponse {
  return { details: undefined };
}

export const SetDefaultOrgResponse = {
  encode(message: SetDefaultOrgResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultOrgResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultOrgResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultOrgResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultOrgResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultOrgResponse>): SetDefaultOrgResponse {
    return SetDefaultOrgResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultOrgResponse>): SetDefaultOrgResponse {
    const message = createBaseSetDefaultOrgResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultOrgRequest(): GetDefaultOrgRequest {
  return {};
}

export const GetDefaultOrgRequest = {
  encode(_: GetDefaultOrgRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultOrgRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultOrgRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetDefaultOrgRequest {
    return {};
  },

  toJSON(_: GetDefaultOrgRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetDefaultOrgRequest>): GetDefaultOrgRequest {
    return GetDefaultOrgRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetDefaultOrgRequest>): GetDefaultOrgRequest {
    const message = createBaseGetDefaultOrgRequest();
    return message;
  },
};

function createBaseGetDefaultOrgResponse(): GetDefaultOrgResponse {
  return { org: undefined };
}

export const GetDefaultOrgResponse = {
  encode(message: GetDefaultOrgResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.org !== undefined) {
      Org.encode(message.org, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultOrgResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultOrgResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.org = Org.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultOrgResponse {
    return { org: isSet(object.org) ? Org.fromJSON(object.org) : undefined };
  },

  toJSON(message: GetDefaultOrgResponse): unknown {
    const obj: any = {};
    message.org !== undefined && (obj.org = message.org ? Org.toJSON(message.org) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultOrgResponse>): GetDefaultOrgResponse {
    return GetDefaultOrgResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultOrgResponse>): GetDefaultOrgResponse {
    const message = createBaseGetDefaultOrgResponse();
    message.org = (object.org !== undefined && object.org !== null) ? Org.fromPartial(object.org) : undefined;
    return message;
  },
};

function createBaseGetMyInstanceRequest(): GetMyInstanceRequest {
  return {};
}

export const GetMyInstanceRequest = {
  encode(_: GetMyInstanceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyInstanceRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyInstanceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyInstanceRequest {
    return {};
  },

  toJSON(_: GetMyInstanceRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyInstanceRequest>): GetMyInstanceRequest {
    return GetMyInstanceRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyInstanceRequest>): GetMyInstanceRequest {
    const message = createBaseGetMyInstanceRequest();
    return message;
  },
};

function createBaseGetMyInstanceResponse(): GetMyInstanceResponse {
  return { instance: undefined };
}

export const GetMyInstanceResponse = {
  encode(message: GetMyInstanceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.instance !== undefined) {
      InstanceDetail.encode(message.instance, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyInstanceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyInstanceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.instance = InstanceDetail.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyInstanceResponse {
    return { instance: isSet(object.instance) ? InstanceDetail.fromJSON(object.instance) : undefined };
  },

  toJSON(message: GetMyInstanceResponse): unknown {
    const obj: any = {};
    message.instance !== undefined &&
      (obj.instance = message.instance ? InstanceDetail.toJSON(message.instance) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyInstanceResponse>): GetMyInstanceResponse {
    return GetMyInstanceResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyInstanceResponse>): GetMyInstanceResponse {
    const message = createBaseGetMyInstanceResponse();
    message.instance = (object.instance !== undefined && object.instance !== null)
      ? InstanceDetail.fromPartial(object.instance)
      : undefined;
    return message;
  },
};

function createBaseListInstanceDomainsRequest(): ListInstanceDomainsRequest {
  return { query: undefined, sortingColumn: 0, queries: [] };
}

export const ListInstanceDomainsRequest = {
  encode(message: ListInstanceDomainsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.queries) {
      DomainSearchQuery.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListInstanceDomainsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListInstanceDomainsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.sortingColumn = reader.int32() as any;
          break;
        case 3:
          message.queries.push(DomainSearchQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListInstanceDomainsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? domainFieldNameFromJSON(object.sortingColumn) : 0,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => DomainSearchQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListInstanceDomainsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = domainFieldNameToJSON(message.sortingColumn));
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? DomainSearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListInstanceDomainsRequest>): ListInstanceDomainsRequest {
    return ListInstanceDomainsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListInstanceDomainsRequest>): ListInstanceDomainsRequest {
    const message = createBaseListInstanceDomainsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.queries = object.queries?.map((e) => DomainSearchQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListInstanceDomainsResponse(): ListInstanceDomainsResponse {
  return { details: undefined, sortingColumn: 0, result: [] };
}

export const ListInstanceDomainsResponse = {
  encode(message: ListInstanceDomainsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.result) {
      Domain.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListInstanceDomainsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListInstanceDomainsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.sortingColumn = reader.int32() as any;
          break;
        case 3:
          message.result.push(Domain.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListInstanceDomainsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? domainFieldNameFromJSON(object.sortingColumn) : 0,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Domain.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListInstanceDomainsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = domainFieldNameToJSON(message.sortingColumn));
    if (message.result) {
      obj.result = message.result.map((e) => e ? Domain.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListInstanceDomainsResponse>): ListInstanceDomainsResponse {
    return ListInstanceDomainsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListInstanceDomainsResponse>): ListInstanceDomainsResponse {
    const message = createBaseListInstanceDomainsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.result = object.result?.map((e) => Domain.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListSecretGeneratorsRequest(): ListSecretGeneratorsRequest {
  return { query: undefined, queries: [] };
}

export const ListSecretGeneratorsRequest = {
  encode(message: ListSecretGeneratorsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.queries) {
      SecretGeneratorQuery.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListSecretGeneratorsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListSecretGeneratorsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.queries.push(SecretGeneratorQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListSecretGeneratorsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => SecretGeneratorQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListSecretGeneratorsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? SecretGeneratorQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListSecretGeneratorsRequest>): ListSecretGeneratorsRequest {
    return ListSecretGeneratorsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListSecretGeneratorsRequest>): ListSecretGeneratorsRequest {
    const message = createBaseListSecretGeneratorsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.queries = object.queries?.map((e) => SecretGeneratorQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListSecretGeneratorsResponse(): ListSecretGeneratorsResponse {
  return { details: undefined, result: [] };
}

export const ListSecretGeneratorsResponse = {
  encode(message: ListSecretGeneratorsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      SecretGenerator.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListSecretGeneratorsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListSecretGeneratorsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 3:
          message.result.push(SecretGenerator.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListSecretGeneratorsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => SecretGenerator.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListSecretGeneratorsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? SecretGenerator.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListSecretGeneratorsResponse>): ListSecretGeneratorsResponse {
    return ListSecretGeneratorsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListSecretGeneratorsResponse>): ListSecretGeneratorsResponse {
    const message = createBaseListSecretGeneratorsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => SecretGenerator.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetSecretGeneratorRequest(): GetSecretGeneratorRequest {
  return { generatorType: 0 };
}

export const GetSecretGeneratorRequest = {
  encode(message: GetSecretGeneratorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.generatorType !== 0) {
      writer.uint32(8).int32(message.generatorType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSecretGeneratorRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSecretGeneratorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.generatorType = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetSecretGeneratorRequest {
    return { generatorType: isSet(object.generatorType) ? secretGeneratorTypeFromJSON(object.generatorType) : 0 };
  },

  toJSON(message: GetSecretGeneratorRequest): unknown {
    const obj: any = {};
    message.generatorType !== undefined && (obj.generatorType = secretGeneratorTypeToJSON(message.generatorType));
    return obj;
  },

  create(base?: DeepPartial<GetSecretGeneratorRequest>): GetSecretGeneratorRequest {
    return GetSecretGeneratorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSecretGeneratorRequest>): GetSecretGeneratorRequest {
    const message = createBaseGetSecretGeneratorRequest();
    message.generatorType = object.generatorType ?? 0;
    return message;
  },
};

function createBaseGetSecretGeneratorResponse(): GetSecretGeneratorResponse {
  return { secretGenerator: undefined };
}

export const GetSecretGeneratorResponse = {
  encode(message: GetSecretGeneratorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.secretGenerator !== undefined) {
      SecretGenerator.encode(message.secretGenerator, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSecretGeneratorResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSecretGeneratorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.secretGenerator = SecretGenerator.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetSecretGeneratorResponse {
    return {
      secretGenerator: isSet(object.secretGenerator) ? SecretGenerator.fromJSON(object.secretGenerator) : undefined,
    };
  },

  toJSON(message: GetSecretGeneratorResponse): unknown {
    const obj: any = {};
    message.secretGenerator !== undefined &&
      (obj.secretGenerator = message.secretGenerator ? SecretGenerator.toJSON(message.secretGenerator) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetSecretGeneratorResponse>): GetSecretGeneratorResponse {
    return GetSecretGeneratorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSecretGeneratorResponse>): GetSecretGeneratorResponse {
    const message = createBaseGetSecretGeneratorResponse();
    message.secretGenerator = (object.secretGenerator !== undefined && object.secretGenerator !== null)
      ? SecretGenerator.fromPartial(object.secretGenerator)
      : undefined;
    return message;
  },
};

function createBaseUpdateSecretGeneratorRequest(): UpdateSecretGeneratorRequest {
  return {
    generatorType: 0,
    length: 0,
    expiry: undefined,
    includeLowerLetters: false,
    includeUpperLetters: false,
    includeDigits: false,
    includeSymbols: false,
  };
}

export const UpdateSecretGeneratorRequest = {
  encode(message: UpdateSecretGeneratorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.generatorType !== 0) {
      writer.uint32(8).int32(message.generatorType);
    }
    if (message.length !== 0) {
      writer.uint32(16).uint32(message.length);
    }
    if (message.expiry !== undefined) {
      Duration.encode(message.expiry, writer.uint32(26).fork()).ldelim();
    }
    if (message.includeLowerLetters === true) {
      writer.uint32(32).bool(message.includeLowerLetters);
    }
    if (message.includeUpperLetters === true) {
      writer.uint32(40).bool(message.includeUpperLetters);
    }
    if (message.includeDigits === true) {
      writer.uint32(48).bool(message.includeDigits);
    }
    if (message.includeSymbols === true) {
      writer.uint32(56).bool(message.includeSymbols);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSecretGeneratorRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSecretGeneratorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.generatorType = reader.int32() as any;
          break;
        case 2:
          message.length = reader.uint32();
          break;
        case 3:
          message.expiry = Duration.decode(reader, reader.uint32());
          break;
        case 4:
          message.includeLowerLetters = reader.bool();
          break;
        case 5:
          message.includeUpperLetters = reader.bool();
          break;
        case 6:
          message.includeDigits = reader.bool();
          break;
        case 7:
          message.includeSymbols = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSecretGeneratorRequest {
    return {
      generatorType: isSet(object.generatorType) ? secretGeneratorTypeFromJSON(object.generatorType) : 0,
      length: isSet(object.length) ? Number(object.length) : 0,
      expiry: isSet(object.expiry) ? Duration.fromJSON(object.expiry) : undefined,
      includeLowerLetters: isSet(object.includeLowerLetters) ? Boolean(object.includeLowerLetters) : false,
      includeUpperLetters: isSet(object.includeUpperLetters) ? Boolean(object.includeUpperLetters) : false,
      includeDigits: isSet(object.includeDigits) ? Boolean(object.includeDigits) : false,
      includeSymbols: isSet(object.includeSymbols) ? Boolean(object.includeSymbols) : false,
    };
  },

  toJSON(message: UpdateSecretGeneratorRequest): unknown {
    const obj: any = {};
    message.generatorType !== undefined && (obj.generatorType = secretGeneratorTypeToJSON(message.generatorType));
    message.length !== undefined && (obj.length = Math.round(message.length));
    message.expiry !== undefined && (obj.expiry = message.expiry ? Duration.toJSON(message.expiry) : undefined);
    message.includeLowerLetters !== undefined && (obj.includeLowerLetters = message.includeLowerLetters);
    message.includeUpperLetters !== undefined && (obj.includeUpperLetters = message.includeUpperLetters);
    message.includeDigits !== undefined && (obj.includeDigits = message.includeDigits);
    message.includeSymbols !== undefined && (obj.includeSymbols = message.includeSymbols);
    return obj;
  },

  create(base?: DeepPartial<UpdateSecretGeneratorRequest>): UpdateSecretGeneratorRequest {
    return UpdateSecretGeneratorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSecretGeneratorRequest>): UpdateSecretGeneratorRequest {
    const message = createBaseUpdateSecretGeneratorRequest();
    message.generatorType = object.generatorType ?? 0;
    message.length = object.length ?? 0;
    message.expiry = (object.expiry !== undefined && object.expiry !== null)
      ? Duration.fromPartial(object.expiry)
      : undefined;
    message.includeLowerLetters = object.includeLowerLetters ?? false;
    message.includeUpperLetters = object.includeUpperLetters ?? false;
    message.includeDigits = object.includeDigits ?? false;
    message.includeSymbols = object.includeSymbols ?? false;
    return message;
  },
};

function createBaseUpdateSecretGeneratorResponse(): UpdateSecretGeneratorResponse {
  return { details: undefined };
}

export const UpdateSecretGeneratorResponse = {
  encode(message: UpdateSecretGeneratorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSecretGeneratorResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSecretGeneratorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSecretGeneratorResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateSecretGeneratorResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateSecretGeneratorResponse>): UpdateSecretGeneratorResponse {
    return UpdateSecretGeneratorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSecretGeneratorResponse>): UpdateSecretGeneratorResponse {
    const message = createBaseUpdateSecretGeneratorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetSMTPConfigRequest(): GetSMTPConfigRequest {
  return {};
}

export const GetSMTPConfigRequest = {
  encode(_: GetSMTPConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSMTPConfigRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSMTPConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetSMTPConfigRequest {
    return {};
  },

  toJSON(_: GetSMTPConfigRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetSMTPConfigRequest>): GetSMTPConfigRequest {
    return GetSMTPConfigRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetSMTPConfigRequest>): GetSMTPConfigRequest {
    const message = createBaseGetSMTPConfigRequest();
    return message;
  },
};

function createBaseGetSMTPConfigResponse(): GetSMTPConfigResponse {
  return { smtpConfig: undefined };
}

export const GetSMTPConfigResponse = {
  encode(message: GetSMTPConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.smtpConfig !== undefined) {
      SMTPConfig.encode(message.smtpConfig, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSMTPConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSMTPConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.smtpConfig = SMTPConfig.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetSMTPConfigResponse {
    return { smtpConfig: isSet(object.smtpConfig) ? SMTPConfig.fromJSON(object.smtpConfig) : undefined };
  },

  toJSON(message: GetSMTPConfigResponse): unknown {
    const obj: any = {};
    message.smtpConfig !== undefined &&
      (obj.smtpConfig = message.smtpConfig ? SMTPConfig.toJSON(message.smtpConfig) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetSMTPConfigResponse>): GetSMTPConfigResponse {
    return GetSMTPConfigResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSMTPConfigResponse>): GetSMTPConfigResponse {
    const message = createBaseGetSMTPConfigResponse();
    message.smtpConfig = (object.smtpConfig !== undefined && object.smtpConfig !== null)
      ? SMTPConfig.fromPartial(object.smtpConfig)
      : undefined;
    return message;
  },
};

function createBaseAddSMTPConfigRequest(): AddSMTPConfigRequest {
  return { senderAddress: "", senderName: "", tls: false, host: "", user: "", password: "" };
}

export const AddSMTPConfigRequest = {
  encode(message: AddSMTPConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.senderAddress !== "") {
      writer.uint32(10).string(message.senderAddress);
    }
    if (message.senderName !== "") {
      writer.uint32(18).string(message.senderName);
    }
    if (message.tls === true) {
      writer.uint32(24).bool(message.tls);
    }
    if (message.host !== "") {
      writer.uint32(34).string(message.host);
    }
    if (message.user !== "") {
      writer.uint32(42).string(message.user);
    }
    if (message.password !== "") {
      writer.uint32(50).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddSMTPConfigRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddSMTPConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.senderAddress = reader.string();
          break;
        case 2:
          message.senderName = reader.string();
          break;
        case 3:
          message.tls = reader.bool();
          break;
        case 4:
          message.host = reader.string();
          break;
        case 5:
          message.user = reader.string();
          break;
        case 6:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddSMTPConfigRequest {
    return {
      senderAddress: isSet(object.senderAddress) ? String(object.senderAddress) : "",
      senderName: isSet(object.senderName) ? String(object.senderName) : "",
      tls: isSet(object.tls) ? Boolean(object.tls) : false,
      host: isSet(object.host) ? String(object.host) : "",
      user: isSet(object.user) ? String(object.user) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: AddSMTPConfigRequest): unknown {
    const obj: any = {};
    message.senderAddress !== undefined && (obj.senderAddress = message.senderAddress);
    message.senderName !== undefined && (obj.senderName = message.senderName);
    message.tls !== undefined && (obj.tls = message.tls);
    message.host !== undefined && (obj.host = message.host);
    message.user !== undefined && (obj.user = message.user);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  create(base?: DeepPartial<AddSMTPConfigRequest>): AddSMTPConfigRequest {
    return AddSMTPConfigRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddSMTPConfigRequest>): AddSMTPConfigRequest {
    const message = createBaseAddSMTPConfigRequest();
    message.senderAddress = object.senderAddress ?? "";
    message.senderName = object.senderName ?? "";
    message.tls = object.tls ?? false;
    message.host = object.host ?? "";
    message.user = object.user ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseAddSMTPConfigResponse(): AddSMTPConfigResponse {
  return { details: undefined };
}

export const AddSMTPConfigResponse = {
  encode(message: AddSMTPConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddSMTPConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddSMTPConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddSMTPConfigResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddSMTPConfigResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddSMTPConfigResponse>): AddSMTPConfigResponse {
    return AddSMTPConfigResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddSMTPConfigResponse>): AddSMTPConfigResponse {
    const message = createBaseAddSMTPConfigResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateSMTPConfigRequest(): UpdateSMTPConfigRequest {
  return { senderAddress: "", senderName: "", tls: false, host: "", user: "" };
}

export const UpdateSMTPConfigRequest = {
  encode(message: UpdateSMTPConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.senderAddress !== "") {
      writer.uint32(10).string(message.senderAddress);
    }
    if (message.senderName !== "") {
      writer.uint32(18).string(message.senderName);
    }
    if (message.tls === true) {
      writer.uint32(24).bool(message.tls);
    }
    if (message.host !== "") {
      writer.uint32(34).string(message.host);
    }
    if (message.user !== "") {
      writer.uint32(42).string(message.user);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSMTPConfigRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSMTPConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.senderAddress = reader.string();
          break;
        case 2:
          message.senderName = reader.string();
          break;
        case 3:
          message.tls = reader.bool();
          break;
        case 4:
          message.host = reader.string();
          break;
        case 5:
          message.user = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSMTPConfigRequest {
    return {
      senderAddress: isSet(object.senderAddress) ? String(object.senderAddress) : "",
      senderName: isSet(object.senderName) ? String(object.senderName) : "",
      tls: isSet(object.tls) ? Boolean(object.tls) : false,
      host: isSet(object.host) ? String(object.host) : "",
      user: isSet(object.user) ? String(object.user) : "",
    };
  },

  toJSON(message: UpdateSMTPConfigRequest): unknown {
    const obj: any = {};
    message.senderAddress !== undefined && (obj.senderAddress = message.senderAddress);
    message.senderName !== undefined && (obj.senderName = message.senderName);
    message.tls !== undefined && (obj.tls = message.tls);
    message.host !== undefined && (obj.host = message.host);
    message.user !== undefined && (obj.user = message.user);
    return obj;
  },

  create(base?: DeepPartial<UpdateSMTPConfigRequest>): UpdateSMTPConfigRequest {
    return UpdateSMTPConfigRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSMTPConfigRequest>): UpdateSMTPConfigRequest {
    const message = createBaseUpdateSMTPConfigRequest();
    message.senderAddress = object.senderAddress ?? "";
    message.senderName = object.senderName ?? "";
    message.tls = object.tls ?? false;
    message.host = object.host ?? "";
    message.user = object.user ?? "";
    return message;
  },
};

function createBaseUpdateSMTPConfigResponse(): UpdateSMTPConfigResponse {
  return { details: undefined };
}

export const UpdateSMTPConfigResponse = {
  encode(message: UpdateSMTPConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSMTPConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSMTPConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSMTPConfigResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateSMTPConfigResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateSMTPConfigResponse>): UpdateSMTPConfigResponse {
    return UpdateSMTPConfigResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSMTPConfigResponse>): UpdateSMTPConfigResponse {
    const message = createBaseUpdateSMTPConfigResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateSMTPConfigPasswordRequest(): UpdateSMTPConfigPasswordRequest {
  return { password: "" };
}

export const UpdateSMTPConfigPasswordRequest = {
  encode(message: UpdateSMTPConfigPasswordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.password !== "") {
      writer.uint32(10).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSMTPConfigPasswordRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSMTPConfigPasswordRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSMTPConfigPasswordRequest {
    return { password: isSet(object.password) ? String(object.password) : "" };
  },

  toJSON(message: UpdateSMTPConfigPasswordRequest): unknown {
    const obj: any = {};
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  create(base?: DeepPartial<UpdateSMTPConfigPasswordRequest>): UpdateSMTPConfigPasswordRequest {
    return UpdateSMTPConfigPasswordRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSMTPConfigPasswordRequest>): UpdateSMTPConfigPasswordRequest {
    const message = createBaseUpdateSMTPConfigPasswordRequest();
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseUpdateSMTPConfigPasswordResponse(): UpdateSMTPConfigPasswordResponse {
  return { details: undefined };
}

export const UpdateSMTPConfigPasswordResponse = {
  encode(message: UpdateSMTPConfigPasswordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSMTPConfigPasswordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSMTPConfigPasswordResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSMTPConfigPasswordResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateSMTPConfigPasswordResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateSMTPConfigPasswordResponse>): UpdateSMTPConfigPasswordResponse {
    return UpdateSMTPConfigPasswordResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSMTPConfigPasswordResponse>): UpdateSMTPConfigPasswordResponse {
    const message = createBaseUpdateSMTPConfigPasswordResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveSMTPConfigRequest(): RemoveSMTPConfigRequest {
  return {};
}

export const RemoveSMTPConfigRequest = {
  encode(_: RemoveSMTPConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveSMTPConfigRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveSMTPConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveSMTPConfigRequest {
    return {};
  },

  toJSON(_: RemoveSMTPConfigRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveSMTPConfigRequest>): RemoveSMTPConfigRequest {
    return RemoveSMTPConfigRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveSMTPConfigRequest>): RemoveSMTPConfigRequest {
    const message = createBaseRemoveSMTPConfigRequest();
    return message;
  },
};

function createBaseRemoveSMTPConfigResponse(): RemoveSMTPConfigResponse {
  return { details: undefined };
}

export const RemoveSMTPConfigResponse = {
  encode(message: RemoveSMTPConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveSMTPConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveSMTPConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveSMTPConfigResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveSMTPConfigResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveSMTPConfigResponse>): RemoveSMTPConfigResponse {
    return RemoveSMTPConfigResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveSMTPConfigResponse>): RemoveSMTPConfigResponse {
    const message = createBaseRemoveSMTPConfigResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListSMSProvidersRequest(): ListSMSProvidersRequest {
  return { query: undefined };
}

export const ListSMSProvidersRequest = {
  encode(message: ListSMSProvidersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListSMSProvidersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListSMSProvidersRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListSMSProvidersRequest {
    return { query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined };
  },

  toJSON(message: ListSMSProvidersRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ListSMSProvidersRequest>): ListSMSProvidersRequest {
    return ListSMSProvidersRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListSMSProvidersRequest>): ListSMSProvidersRequest {
    const message = createBaseListSMSProvidersRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseListSMSProvidersResponse(): ListSMSProvidersResponse {
  return { details: undefined, result: [] };
}

export const ListSMSProvidersResponse = {
  encode(message: ListSMSProvidersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      SMSProvider.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListSMSProvidersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListSMSProvidersResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 3:
          message.result.push(SMSProvider.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListSMSProvidersResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => SMSProvider.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListSMSProvidersResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? SMSProvider.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListSMSProvidersResponse>): ListSMSProvidersResponse {
    return ListSMSProvidersResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListSMSProvidersResponse>): ListSMSProvidersResponse {
    const message = createBaseListSMSProvidersResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => SMSProvider.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetSMSProviderRequest(): GetSMSProviderRequest {
  return { id: "" };
}

export const GetSMSProviderRequest = {
  encode(message: GetSMSProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSMSProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSMSProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetSMSProviderRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: GetSMSProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<GetSMSProviderRequest>): GetSMSProviderRequest {
    return GetSMSProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSMSProviderRequest>): GetSMSProviderRequest {
    const message = createBaseGetSMSProviderRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseGetSMSProviderResponse(): GetSMSProviderResponse {
  return { config: undefined };
}

export const GetSMSProviderResponse = {
  encode(message: GetSMSProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.config !== undefined) {
      SMSProvider.encode(message.config, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSMSProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSMSProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.config = SMSProvider.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetSMSProviderResponse {
    return { config: isSet(object.config) ? SMSProvider.fromJSON(object.config) : undefined };
  },

  toJSON(message: GetSMSProviderResponse): unknown {
    const obj: any = {};
    message.config !== undefined && (obj.config = message.config ? SMSProvider.toJSON(message.config) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetSMSProviderResponse>): GetSMSProviderResponse {
    return GetSMSProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSMSProviderResponse>): GetSMSProviderResponse {
    const message = createBaseGetSMSProviderResponse();
    message.config = (object.config !== undefined && object.config !== null)
      ? SMSProvider.fromPartial(object.config)
      : undefined;
    return message;
  },
};

function createBaseAddSMSProviderTwilioRequest(): AddSMSProviderTwilioRequest {
  return { sid: "", token: "", senderNumber: "" };
}

export const AddSMSProviderTwilioRequest = {
  encode(message: AddSMSProviderTwilioRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sid !== "") {
      writer.uint32(10).string(message.sid);
    }
    if (message.token !== "") {
      writer.uint32(18).string(message.token);
    }
    if (message.senderNumber !== "") {
      writer.uint32(26).string(message.senderNumber);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddSMSProviderTwilioRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddSMSProviderTwilioRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sid = reader.string();
          break;
        case 2:
          message.token = reader.string();
          break;
        case 3:
          message.senderNumber = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddSMSProviderTwilioRequest {
    return {
      sid: isSet(object.sid) ? String(object.sid) : "",
      token: isSet(object.token) ? String(object.token) : "",
      senderNumber: isSet(object.senderNumber) ? String(object.senderNumber) : "",
    };
  },

  toJSON(message: AddSMSProviderTwilioRequest): unknown {
    const obj: any = {};
    message.sid !== undefined && (obj.sid = message.sid);
    message.token !== undefined && (obj.token = message.token);
    message.senderNumber !== undefined && (obj.senderNumber = message.senderNumber);
    return obj;
  },

  create(base?: DeepPartial<AddSMSProviderTwilioRequest>): AddSMSProviderTwilioRequest {
    return AddSMSProviderTwilioRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddSMSProviderTwilioRequest>): AddSMSProviderTwilioRequest {
    const message = createBaseAddSMSProviderTwilioRequest();
    message.sid = object.sid ?? "";
    message.token = object.token ?? "";
    message.senderNumber = object.senderNumber ?? "";
    return message;
  },
};

function createBaseAddSMSProviderTwilioResponse(): AddSMSProviderTwilioResponse {
  return { details: undefined, id: "" };
}

export const AddSMSProviderTwilioResponse = {
  encode(message: AddSMSProviderTwilioResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddSMSProviderTwilioResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddSMSProviderTwilioResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddSMSProviderTwilioResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
    };
  },

  toJSON(message: AddSMSProviderTwilioResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<AddSMSProviderTwilioResponse>): AddSMSProviderTwilioResponse {
    return AddSMSProviderTwilioResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddSMSProviderTwilioResponse>): AddSMSProviderTwilioResponse {
    const message = createBaseAddSMSProviderTwilioResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseUpdateSMSProviderTwilioRequest(): UpdateSMSProviderTwilioRequest {
  return { id: "", sid: "", senderNumber: "" };
}

export const UpdateSMSProviderTwilioRequest = {
  encode(message: UpdateSMSProviderTwilioRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.sid !== "") {
      writer.uint32(18).string(message.sid);
    }
    if (message.senderNumber !== "") {
      writer.uint32(26).string(message.senderNumber);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSMSProviderTwilioRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSMSProviderTwilioRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.sid = reader.string();
          break;
        case 3:
          message.senderNumber = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSMSProviderTwilioRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      sid: isSet(object.sid) ? String(object.sid) : "",
      senderNumber: isSet(object.senderNumber) ? String(object.senderNumber) : "",
    };
  },

  toJSON(message: UpdateSMSProviderTwilioRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.sid !== undefined && (obj.sid = message.sid);
    message.senderNumber !== undefined && (obj.senderNumber = message.senderNumber);
    return obj;
  },

  create(base?: DeepPartial<UpdateSMSProviderTwilioRequest>): UpdateSMSProviderTwilioRequest {
    return UpdateSMSProviderTwilioRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSMSProviderTwilioRequest>): UpdateSMSProviderTwilioRequest {
    const message = createBaseUpdateSMSProviderTwilioRequest();
    message.id = object.id ?? "";
    message.sid = object.sid ?? "";
    message.senderNumber = object.senderNumber ?? "";
    return message;
  },
};

function createBaseUpdateSMSProviderTwilioResponse(): UpdateSMSProviderTwilioResponse {
  return { details: undefined };
}

export const UpdateSMSProviderTwilioResponse = {
  encode(message: UpdateSMSProviderTwilioResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSMSProviderTwilioResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSMSProviderTwilioResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSMSProviderTwilioResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateSMSProviderTwilioResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateSMSProviderTwilioResponse>): UpdateSMSProviderTwilioResponse {
    return UpdateSMSProviderTwilioResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSMSProviderTwilioResponse>): UpdateSMSProviderTwilioResponse {
    const message = createBaseUpdateSMSProviderTwilioResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateSMSProviderTwilioTokenRequest(): UpdateSMSProviderTwilioTokenRequest {
  return { id: "", token: "" };
}

export const UpdateSMSProviderTwilioTokenRequest = {
  encode(message: UpdateSMSProviderTwilioTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.token !== "") {
      writer.uint32(18).string(message.token);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSMSProviderTwilioTokenRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSMSProviderTwilioTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.token = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSMSProviderTwilioTokenRequest {
    return { id: isSet(object.id) ? String(object.id) : "", token: isSet(object.token) ? String(object.token) : "" };
  },

  toJSON(message: UpdateSMSProviderTwilioTokenRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.token !== undefined && (obj.token = message.token);
    return obj;
  },

  create(base?: DeepPartial<UpdateSMSProviderTwilioTokenRequest>): UpdateSMSProviderTwilioTokenRequest {
    return UpdateSMSProviderTwilioTokenRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSMSProviderTwilioTokenRequest>): UpdateSMSProviderTwilioTokenRequest {
    const message = createBaseUpdateSMSProviderTwilioTokenRequest();
    message.id = object.id ?? "";
    message.token = object.token ?? "";
    return message;
  },
};

function createBaseUpdateSMSProviderTwilioTokenResponse(): UpdateSMSProviderTwilioTokenResponse {
  return { details: undefined };
}

export const UpdateSMSProviderTwilioTokenResponse = {
  encode(message: UpdateSMSProviderTwilioTokenResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSMSProviderTwilioTokenResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSMSProviderTwilioTokenResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateSMSProviderTwilioTokenResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateSMSProviderTwilioTokenResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateSMSProviderTwilioTokenResponse>): UpdateSMSProviderTwilioTokenResponse {
    return UpdateSMSProviderTwilioTokenResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateSMSProviderTwilioTokenResponse>): UpdateSMSProviderTwilioTokenResponse {
    const message = createBaseUpdateSMSProviderTwilioTokenResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseActivateSMSProviderRequest(): ActivateSMSProviderRequest {
  return { id: "" };
}

export const ActivateSMSProviderRequest = {
  encode(message: ActivateSMSProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActivateSMSProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActivateSMSProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ActivateSMSProviderRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: ActivateSMSProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<ActivateSMSProviderRequest>): ActivateSMSProviderRequest {
    return ActivateSMSProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ActivateSMSProviderRequest>): ActivateSMSProviderRequest {
    const message = createBaseActivateSMSProviderRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseActivateSMSProviderResponse(): ActivateSMSProviderResponse {
  return { details: undefined };
}

export const ActivateSMSProviderResponse = {
  encode(message: ActivateSMSProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActivateSMSProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActivateSMSProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ActivateSMSProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ActivateSMSProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ActivateSMSProviderResponse>): ActivateSMSProviderResponse {
    return ActivateSMSProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ActivateSMSProviderResponse>): ActivateSMSProviderResponse {
    const message = createBaseActivateSMSProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseDeactivateSMSProviderRequest(): DeactivateSMSProviderRequest {
  return { id: "" };
}

export const DeactivateSMSProviderRequest = {
  encode(message: DeactivateSMSProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeactivateSMSProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeactivateSMSProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeactivateSMSProviderRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: DeactivateSMSProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<DeactivateSMSProviderRequest>): DeactivateSMSProviderRequest {
    return DeactivateSMSProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeactivateSMSProviderRequest>): DeactivateSMSProviderRequest {
    const message = createBaseDeactivateSMSProviderRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseDeactivateSMSProviderResponse(): DeactivateSMSProviderResponse {
  return { details: undefined };
}

export const DeactivateSMSProviderResponse = {
  encode(message: DeactivateSMSProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeactivateSMSProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeactivateSMSProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeactivateSMSProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeactivateSMSProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeactivateSMSProviderResponse>): DeactivateSMSProviderResponse {
    return DeactivateSMSProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeactivateSMSProviderResponse>): DeactivateSMSProviderResponse {
    const message = createBaseDeactivateSMSProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveSMSProviderRequest(): RemoveSMSProviderRequest {
  return { id: "" };
}

export const RemoveSMSProviderRequest = {
  encode(message: RemoveSMSProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveSMSProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveSMSProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveSMSProviderRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: RemoveSMSProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<RemoveSMSProviderRequest>): RemoveSMSProviderRequest {
    return RemoveSMSProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveSMSProviderRequest>): RemoveSMSProviderRequest {
    const message = createBaseRemoveSMSProviderRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseRemoveSMSProviderResponse(): RemoveSMSProviderResponse {
  return { details: undefined };
}

export const RemoveSMSProviderResponse = {
  encode(message: RemoveSMSProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveSMSProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveSMSProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveSMSProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveSMSProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveSMSProviderResponse>): RemoveSMSProviderResponse {
    return RemoveSMSProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveSMSProviderResponse>): RemoveSMSProviderResponse {
    const message = createBaseRemoveSMSProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetFileSystemNotificationProviderRequest(): GetFileSystemNotificationProviderRequest {
  return {};
}

export const GetFileSystemNotificationProviderRequest = {
  encode(_: GetFileSystemNotificationProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetFileSystemNotificationProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetFileSystemNotificationProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetFileSystemNotificationProviderRequest {
    return {};
  },

  toJSON(_: GetFileSystemNotificationProviderRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetFileSystemNotificationProviderRequest>): GetFileSystemNotificationProviderRequest {
    return GetFileSystemNotificationProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetFileSystemNotificationProviderRequest>): GetFileSystemNotificationProviderRequest {
    const message = createBaseGetFileSystemNotificationProviderRequest();
    return message;
  },
};

function createBaseGetFileSystemNotificationProviderResponse(): GetFileSystemNotificationProviderResponse {
  return { provider: undefined };
}

export const GetFileSystemNotificationProviderResponse = {
  encode(message: GetFileSystemNotificationProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider !== undefined) {
      DebugNotificationProvider.encode(message.provider, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetFileSystemNotificationProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetFileSystemNotificationProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider = DebugNotificationProvider.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetFileSystemNotificationProviderResponse {
    return { provider: isSet(object.provider) ? DebugNotificationProvider.fromJSON(object.provider) : undefined };
  },

  toJSON(message: GetFileSystemNotificationProviderResponse): unknown {
    const obj: any = {};
    message.provider !== undefined &&
      (obj.provider = message.provider ? DebugNotificationProvider.toJSON(message.provider) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetFileSystemNotificationProviderResponse>): GetFileSystemNotificationProviderResponse {
    return GetFileSystemNotificationProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetFileSystemNotificationProviderResponse>,
  ): GetFileSystemNotificationProviderResponse {
    const message = createBaseGetFileSystemNotificationProviderResponse();
    message.provider = (object.provider !== undefined && object.provider !== null)
      ? DebugNotificationProvider.fromPartial(object.provider)
      : undefined;
    return message;
  },
};

function createBaseGetLogNotificationProviderRequest(): GetLogNotificationProviderRequest {
  return {};
}

export const GetLogNotificationProviderRequest = {
  encode(_: GetLogNotificationProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLogNotificationProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLogNotificationProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetLogNotificationProviderRequest {
    return {};
  },

  toJSON(_: GetLogNotificationProviderRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetLogNotificationProviderRequest>): GetLogNotificationProviderRequest {
    return GetLogNotificationProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetLogNotificationProviderRequest>): GetLogNotificationProviderRequest {
    const message = createBaseGetLogNotificationProviderRequest();
    return message;
  },
};

function createBaseGetLogNotificationProviderResponse(): GetLogNotificationProviderResponse {
  return { provider: undefined };
}

export const GetLogNotificationProviderResponse = {
  encode(message: GetLogNotificationProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider !== undefined) {
      DebugNotificationProvider.encode(message.provider, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLogNotificationProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLogNotificationProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider = DebugNotificationProvider.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetLogNotificationProviderResponse {
    return { provider: isSet(object.provider) ? DebugNotificationProvider.fromJSON(object.provider) : undefined };
  },

  toJSON(message: GetLogNotificationProviderResponse): unknown {
    const obj: any = {};
    message.provider !== undefined &&
      (obj.provider = message.provider ? DebugNotificationProvider.toJSON(message.provider) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetLogNotificationProviderResponse>): GetLogNotificationProviderResponse {
    return GetLogNotificationProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetLogNotificationProviderResponse>): GetLogNotificationProviderResponse {
    const message = createBaseGetLogNotificationProviderResponse();
    message.provider = (object.provider !== undefined && object.provider !== null)
      ? DebugNotificationProvider.fromPartial(object.provider)
      : undefined;
    return message;
  },
};

function createBaseGetOIDCSettingsRequest(): GetOIDCSettingsRequest {
  return {};
}

export const GetOIDCSettingsRequest = {
  encode(_: GetOIDCSettingsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOIDCSettingsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOIDCSettingsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetOIDCSettingsRequest {
    return {};
  },

  toJSON(_: GetOIDCSettingsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetOIDCSettingsRequest>): GetOIDCSettingsRequest {
    return GetOIDCSettingsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetOIDCSettingsRequest>): GetOIDCSettingsRequest {
    const message = createBaseGetOIDCSettingsRequest();
    return message;
  },
};

function createBaseGetOIDCSettingsResponse(): GetOIDCSettingsResponse {
  return { settings: undefined };
}

export const GetOIDCSettingsResponse = {
  encode(message: GetOIDCSettingsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.settings !== undefined) {
      OIDCSettings.encode(message.settings, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOIDCSettingsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOIDCSettingsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.settings = OIDCSettings.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetOIDCSettingsResponse {
    return { settings: isSet(object.settings) ? OIDCSettings.fromJSON(object.settings) : undefined };
  },

  toJSON(message: GetOIDCSettingsResponse): unknown {
    const obj: any = {};
    message.settings !== undefined &&
      (obj.settings = message.settings ? OIDCSettings.toJSON(message.settings) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetOIDCSettingsResponse>): GetOIDCSettingsResponse {
    return GetOIDCSettingsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetOIDCSettingsResponse>): GetOIDCSettingsResponse {
    const message = createBaseGetOIDCSettingsResponse();
    message.settings = (object.settings !== undefined && object.settings !== null)
      ? OIDCSettings.fromPartial(object.settings)
      : undefined;
    return message;
  },
};

function createBaseAddOIDCSettingsRequest(): AddOIDCSettingsRequest {
  return {
    accessTokenLifetime: undefined,
    idTokenLifetime: undefined,
    refreshTokenIdleExpiration: undefined,
    refreshTokenExpiration: undefined,
  };
}

export const AddOIDCSettingsRequest = {
  encode(message: AddOIDCSettingsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accessTokenLifetime !== undefined) {
      Duration.encode(message.accessTokenLifetime, writer.uint32(10).fork()).ldelim();
    }
    if (message.idTokenLifetime !== undefined) {
      Duration.encode(message.idTokenLifetime, writer.uint32(18).fork()).ldelim();
    }
    if (message.refreshTokenIdleExpiration !== undefined) {
      Duration.encode(message.refreshTokenIdleExpiration, writer.uint32(26).fork()).ldelim();
    }
    if (message.refreshTokenExpiration !== undefined) {
      Duration.encode(message.refreshTokenExpiration, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOIDCSettingsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOIDCSettingsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessTokenLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 2:
          message.idTokenLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 3:
          message.refreshTokenIdleExpiration = Duration.decode(reader, reader.uint32());
          break;
        case 4:
          message.refreshTokenExpiration = Duration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddOIDCSettingsRequest {
    return {
      accessTokenLifetime: isSet(object.accessTokenLifetime)
        ? Duration.fromJSON(object.accessTokenLifetime)
        : undefined,
      idTokenLifetime: isSet(object.idTokenLifetime) ? Duration.fromJSON(object.idTokenLifetime) : undefined,
      refreshTokenIdleExpiration: isSet(object.refreshTokenIdleExpiration)
        ? Duration.fromJSON(object.refreshTokenIdleExpiration)
        : undefined,
      refreshTokenExpiration: isSet(object.refreshTokenExpiration)
        ? Duration.fromJSON(object.refreshTokenExpiration)
        : undefined,
    };
  },

  toJSON(message: AddOIDCSettingsRequest): unknown {
    const obj: any = {};
    message.accessTokenLifetime !== undefined &&
      (obj.accessTokenLifetime = message.accessTokenLifetime
        ? Duration.toJSON(message.accessTokenLifetime)
        : undefined);
    message.idTokenLifetime !== undefined &&
      (obj.idTokenLifetime = message.idTokenLifetime ? Duration.toJSON(message.idTokenLifetime) : undefined);
    message.refreshTokenIdleExpiration !== undefined &&
      (obj.refreshTokenIdleExpiration = message.refreshTokenIdleExpiration
        ? Duration.toJSON(message.refreshTokenIdleExpiration)
        : undefined);
    message.refreshTokenExpiration !== undefined && (obj.refreshTokenExpiration = message.refreshTokenExpiration
      ? Duration.toJSON(message.refreshTokenExpiration)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddOIDCSettingsRequest>): AddOIDCSettingsRequest {
    return AddOIDCSettingsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOIDCSettingsRequest>): AddOIDCSettingsRequest {
    const message = createBaseAddOIDCSettingsRequest();
    message.accessTokenLifetime = (object.accessTokenLifetime !== undefined && object.accessTokenLifetime !== null)
      ? Duration.fromPartial(object.accessTokenLifetime)
      : undefined;
    message.idTokenLifetime = (object.idTokenLifetime !== undefined && object.idTokenLifetime !== null)
      ? Duration.fromPartial(object.idTokenLifetime)
      : undefined;
    message.refreshTokenIdleExpiration =
      (object.refreshTokenIdleExpiration !== undefined && object.refreshTokenIdleExpiration !== null)
        ? Duration.fromPartial(object.refreshTokenIdleExpiration)
        : undefined;
    message.refreshTokenExpiration =
      (object.refreshTokenExpiration !== undefined && object.refreshTokenExpiration !== null)
        ? Duration.fromPartial(object.refreshTokenExpiration)
        : undefined;
    return message;
  },
};

function createBaseAddOIDCSettingsResponse(): AddOIDCSettingsResponse {
  return { details: undefined };
}

export const AddOIDCSettingsResponse = {
  encode(message: AddOIDCSettingsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOIDCSettingsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOIDCSettingsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddOIDCSettingsResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddOIDCSettingsResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddOIDCSettingsResponse>): AddOIDCSettingsResponse {
    return AddOIDCSettingsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOIDCSettingsResponse>): AddOIDCSettingsResponse {
    const message = createBaseAddOIDCSettingsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateOIDCSettingsRequest(): UpdateOIDCSettingsRequest {
  return {
    accessTokenLifetime: undefined,
    idTokenLifetime: undefined,
    refreshTokenIdleExpiration: undefined,
    refreshTokenExpiration: undefined,
  };
}

export const UpdateOIDCSettingsRequest = {
  encode(message: UpdateOIDCSettingsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accessTokenLifetime !== undefined) {
      Duration.encode(message.accessTokenLifetime, writer.uint32(10).fork()).ldelim();
    }
    if (message.idTokenLifetime !== undefined) {
      Duration.encode(message.idTokenLifetime, writer.uint32(18).fork()).ldelim();
    }
    if (message.refreshTokenIdleExpiration !== undefined) {
      Duration.encode(message.refreshTokenIdleExpiration, writer.uint32(26).fork()).ldelim();
    }
    if (message.refreshTokenExpiration !== undefined) {
      Duration.encode(message.refreshTokenExpiration, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateOIDCSettingsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateOIDCSettingsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessTokenLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 2:
          message.idTokenLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 3:
          message.refreshTokenIdleExpiration = Duration.decode(reader, reader.uint32());
          break;
        case 4:
          message.refreshTokenExpiration = Duration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateOIDCSettingsRequest {
    return {
      accessTokenLifetime: isSet(object.accessTokenLifetime)
        ? Duration.fromJSON(object.accessTokenLifetime)
        : undefined,
      idTokenLifetime: isSet(object.idTokenLifetime) ? Duration.fromJSON(object.idTokenLifetime) : undefined,
      refreshTokenIdleExpiration: isSet(object.refreshTokenIdleExpiration)
        ? Duration.fromJSON(object.refreshTokenIdleExpiration)
        : undefined,
      refreshTokenExpiration: isSet(object.refreshTokenExpiration)
        ? Duration.fromJSON(object.refreshTokenExpiration)
        : undefined,
    };
  },

  toJSON(message: UpdateOIDCSettingsRequest): unknown {
    const obj: any = {};
    message.accessTokenLifetime !== undefined &&
      (obj.accessTokenLifetime = message.accessTokenLifetime
        ? Duration.toJSON(message.accessTokenLifetime)
        : undefined);
    message.idTokenLifetime !== undefined &&
      (obj.idTokenLifetime = message.idTokenLifetime ? Duration.toJSON(message.idTokenLifetime) : undefined);
    message.refreshTokenIdleExpiration !== undefined &&
      (obj.refreshTokenIdleExpiration = message.refreshTokenIdleExpiration
        ? Duration.toJSON(message.refreshTokenIdleExpiration)
        : undefined);
    message.refreshTokenExpiration !== undefined && (obj.refreshTokenExpiration = message.refreshTokenExpiration
      ? Duration.toJSON(message.refreshTokenExpiration)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateOIDCSettingsRequest>): UpdateOIDCSettingsRequest {
    return UpdateOIDCSettingsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateOIDCSettingsRequest>): UpdateOIDCSettingsRequest {
    const message = createBaseUpdateOIDCSettingsRequest();
    message.accessTokenLifetime = (object.accessTokenLifetime !== undefined && object.accessTokenLifetime !== null)
      ? Duration.fromPartial(object.accessTokenLifetime)
      : undefined;
    message.idTokenLifetime = (object.idTokenLifetime !== undefined && object.idTokenLifetime !== null)
      ? Duration.fromPartial(object.idTokenLifetime)
      : undefined;
    message.refreshTokenIdleExpiration =
      (object.refreshTokenIdleExpiration !== undefined && object.refreshTokenIdleExpiration !== null)
        ? Duration.fromPartial(object.refreshTokenIdleExpiration)
        : undefined;
    message.refreshTokenExpiration =
      (object.refreshTokenExpiration !== undefined && object.refreshTokenExpiration !== null)
        ? Duration.fromPartial(object.refreshTokenExpiration)
        : undefined;
    return message;
  },
};

function createBaseUpdateOIDCSettingsResponse(): UpdateOIDCSettingsResponse {
  return { details: undefined };
}

export const UpdateOIDCSettingsResponse = {
  encode(message: UpdateOIDCSettingsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateOIDCSettingsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateOIDCSettingsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateOIDCSettingsResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateOIDCSettingsResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateOIDCSettingsResponse>): UpdateOIDCSettingsResponse {
    return UpdateOIDCSettingsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateOIDCSettingsResponse>): UpdateOIDCSettingsResponse {
    const message = createBaseUpdateOIDCSettingsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetSecurityPolicyRequest(): GetSecurityPolicyRequest {
  return {};
}

export const GetSecurityPolicyRequest = {
  encode(_: GetSecurityPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSecurityPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSecurityPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetSecurityPolicyRequest {
    return {};
  },

  toJSON(_: GetSecurityPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetSecurityPolicyRequest>): GetSecurityPolicyRequest {
    return GetSecurityPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetSecurityPolicyRequest>): GetSecurityPolicyRequest {
    const message = createBaseGetSecurityPolicyRequest();
    return message;
  },
};

function createBaseGetSecurityPolicyResponse(): GetSecurityPolicyResponse {
  return { policy: undefined };
}

export const GetSecurityPolicyResponse = {
  encode(message: GetSecurityPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      SecurityPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSecurityPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSecurityPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = SecurityPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetSecurityPolicyResponse {
    return { policy: isSet(object.policy) ? SecurityPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetSecurityPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? SecurityPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetSecurityPolicyResponse>): GetSecurityPolicyResponse {
    return GetSecurityPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSecurityPolicyResponse>): GetSecurityPolicyResponse {
    const message = createBaseGetSecurityPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? SecurityPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseSetSecurityPolicyRequest(): SetSecurityPolicyRequest {
  return { enableIframeEmbedding: false, allowedOrigins: [] };
}

export const SetSecurityPolicyRequest = {
  encode(message: SetSecurityPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.enableIframeEmbedding === true) {
      writer.uint32(8).bool(message.enableIframeEmbedding);
    }
    for (const v of message.allowedOrigins) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetSecurityPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetSecurityPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.enableIframeEmbedding = reader.bool();
          break;
        case 2:
          message.allowedOrigins.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetSecurityPolicyRequest {
    return {
      enableIframeEmbedding: isSet(object.enableIframeEmbedding) ? Boolean(object.enableIframeEmbedding) : false,
      allowedOrigins: Array.isArray(object?.allowedOrigins) ? object.allowedOrigins.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: SetSecurityPolicyRequest): unknown {
    const obj: any = {};
    message.enableIframeEmbedding !== undefined && (obj.enableIframeEmbedding = message.enableIframeEmbedding);
    if (message.allowedOrigins) {
      obj.allowedOrigins = message.allowedOrigins.map((e) => e);
    } else {
      obj.allowedOrigins = [];
    }
    return obj;
  },

  create(base?: DeepPartial<SetSecurityPolicyRequest>): SetSecurityPolicyRequest {
    return SetSecurityPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetSecurityPolicyRequest>): SetSecurityPolicyRequest {
    const message = createBaseSetSecurityPolicyRequest();
    message.enableIframeEmbedding = object.enableIframeEmbedding ?? false;
    message.allowedOrigins = object.allowedOrigins?.map((e) => e) || [];
    return message;
  },
};

function createBaseSetSecurityPolicyResponse(): SetSecurityPolicyResponse {
  return { details: undefined };
}

export const SetSecurityPolicyResponse = {
  encode(message: SetSecurityPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetSecurityPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetSecurityPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetSecurityPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetSecurityPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetSecurityPolicyResponse>): SetSecurityPolicyResponse {
    return SetSecurityPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetSecurityPolicyResponse>): SetSecurityPolicyResponse {
    const message = createBaseSetSecurityPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseIsOrgUniqueRequest(): IsOrgUniqueRequest {
  return { name: "", domain: "" };
}

export const IsOrgUniqueRequest = {
  encode(message: IsOrgUniqueRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.domain !== "") {
      writer.uint32(18).string(message.domain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsOrgUniqueRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsOrgUniqueRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.domain = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsOrgUniqueRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      domain: isSet(object.domain) ? String(object.domain) : "",
    };
  },

  toJSON(message: IsOrgUniqueRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.domain !== undefined && (obj.domain = message.domain);
    return obj;
  },

  create(base?: DeepPartial<IsOrgUniqueRequest>): IsOrgUniqueRequest {
    return IsOrgUniqueRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IsOrgUniqueRequest>): IsOrgUniqueRequest {
    const message = createBaseIsOrgUniqueRequest();
    message.name = object.name ?? "";
    message.domain = object.domain ?? "";
    return message;
  },
};

function createBaseIsOrgUniqueResponse(): IsOrgUniqueResponse {
  return { isUnique: false };
}

export const IsOrgUniqueResponse = {
  encode(message: IsOrgUniqueResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.isUnique === true) {
      writer.uint32(8).bool(message.isUnique);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsOrgUniqueResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsOrgUniqueResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isUnique = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsOrgUniqueResponse {
    return { isUnique: isSet(object.isUnique) ? Boolean(object.isUnique) : false };
  },

  toJSON(message: IsOrgUniqueResponse): unknown {
    const obj: any = {};
    message.isUnique !== undefined && (obj.isUnique = message.isUnique);
    return obj;
  },

  create(base?: DeepPartial<IsOrgUniqueResponse>): IsOrgUniqueResponse {
    return IsOrgUniqueResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IsOrgUniqueResponse>): IsOrgUniqueResponse {
    const message = createBaseIsOrgUniqueResponse();
    message.isUnique = object.isUnique ?? false;
    return message;
  },
};

function createBaseGetOrgByIDRequest(): GetOrgByIDRequest {
  return { id: "" };
}

export const GetOrgByIDRequest = {
  encode(message: GetOrgByIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOrgByIDRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOrgByIDRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetOrgByIDRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: GetOrgByIDRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<GetOrgByIDRequest>): GetOrgByIDRequest {
    return GetOrgByIDRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetOrgByIDRequest>): GetOrgByIDRequest {
    const message = createBaseGetOrgByIDRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseGetOrgByIDResponse(): GetOrgByIDResponse {
  return { org: undefined };
}

export const GetOrgByIDResponse = {
  encode(message: GetOrgByIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.org !== undefined) {
      Org.encode(message.org, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOrgByIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOrgByIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.org = Org.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetOrgByIDResponse {
    return { org: isSet(object.org) ? Org.fromJSON(object.org) : undefined };
  },

  toJSON(message: GetOrgByIDResponse): unknown {
    const obj: any = {};
    message.org !== undefined && (obj.org = message.org ? Org.toJSON(message.org) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetOrgByIDResponse>): GetOrgByIDResponse {
    return GetOrgByIDResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetOrgByIDResponse>): GetOrgByIDResponse {
    const message = createBaseGetOrgByIDResponse();
    message.org = (object.org !== undefined && object.org !== null) ? Org.fromPartial(object.org) : undefined;
    return message;
  },
};

function createBaseListOrgsRequest(): ListOrgsRequest {
  return { query: undefined, sortingColumn: 0, queries: [] };
}

export const ListOrgsRequest = {
  encode(message: ListOrgsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.queries) {
      OrgQuery.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListOrgsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListOrgsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.sortingColumn = reader.int32() as any;
          break;
        case 3:
          message.queries.push(OrgQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListOrgsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? orgFieldNameFromJSON(object.sortingColumn) : 0,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => OrgQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListOrgsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = orgFieldNameToJSON(message.sortingColumn));
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? OrgQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListOrgsRequest>): ListOrgsRequest {
    return ListOrgsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListOrgsRequest>): ListOrgsRequest {
    const message = createBaseListOrgsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.queries = object.queries?.map((e) => OrgQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListOrgsResponse(): ListOrgsResponse {
  return { details: undefined, sortingColumn: 0, result: [] };
}

export const ListOrgsResponse = {
  encode(message: ListOrgsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.result) {
      Org.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListOrgsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListOrgsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.sortingColumn = reader.int32() as any;
          break;
        case 3:
          message.result.push(Org.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListOrgsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? orgFieldNameFromJSON(object.sortingColumn) : 0,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Org.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListOrgsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = orgFieldNameToJSON(message.sortingColumn));
    if (message.result) {
      obj.result = message.result.map((e) => e ? Org.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListOrgsResponse>): ListOrgsResponse {
    return ListOrgsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListOrgsResponse>): ListOrgsResponse {
    const message = createBaseListOrgsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.result = object.result?.map((e) => Org.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSetUpOrgRequest(): SetUpOrgRequest {
  return { org: undefined, human: undefined, roles: [] };
}

export const SetUpOrgRequest = {
  encode(message: SetUpOrgRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.org !== undefined) {
      SetUpOrgRequest_Org.encode(message.org, writer.uint32(10).fork()).ldelim();
    }
    if (message.human !== undefined) {
      SetUpOrgRequest_Human.encode(message.human, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.roles) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUpOrgRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUpOrgRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.org = SetUpOrgRequest_Org.decode(reader, reader.uint32());
          break;
        case 2:
          message.human = SetUpOrgRequest_Human.decode(reader, reader.uint32());
          break;
        case 3:
          message.roles.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetUpOrgRequest {
    return {
      org: isSet(object.org) ? SetUpOrgRequest_Org.fromJSON(object.org) : undefined,
      human: isSet(object.human) ? SetUpOrgRequest_Human.fromJSON(object.human) : undefined,
      roles: Array.isArray(object?.roles) ? object.roles.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: SetUpOrgRequest): unknown {
    const obj: any = {};
    message.org !== undefined && (obj.org = message.org ? SetUpOrgRequest_Org.toJSON(message.org) : undefined);
    message.human !== undefined &&
      (obj.human = message.human ? SetUpOrgRequest_Human.toJSON(message.human) : undefined);
    if (message.roles) {
      obj.roles = message.roles.map((e) => e);
    } else {
      obj.roles = [];
    }
    return obj;
  },

  create(base?: DeepPartial<SetUpOrgRequest>): SetUpOrgRequest {
    return SetUpOrgRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUpOrgRequest>): SetUpOrgRequest {
    const message = createBaseSetUpOrgRequest();
    message.org = (object.org !== undefined && object.org !== null)
      ? SetUpOrgRequest_Org.fromPartial(object.org)
      : undefined;
    message.human = (object.human !== undefined && object.human !== null)
      ? SetUpOrgRequest_Human.fromPartial(object.human)
      : undefined;
    message.roles = object.roles?.map((e) => e) || [];
    return message;
  },
};

function createBaseSetUpOrgRequest_Org(): SetUpOrgRequest_Org {
  return { name: "", domain: "" };
}

export const SetUpOrgRequest_Org = {
  encode(message: SetUpOrgRequest_Org, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.domain !== "") {
      writer.uint32(18).string(message.domain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUpOrgRequest_Org {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUpOrgRequest_Org();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.domain = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetUpOrgRequest_Org {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      domain: isSet(object.domain) ? String(object.domain) : "",
    };
  },

  toJSON(message: SetUpOrgRequest_Org): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.domain !== undefined && (obj.domain = message.domain);
    return obj;
  },

  create(base?: DeepPartial<SetUpOrgRequest_Org>): SetUpOrgRequest_Org {
    return SetUpOrgRequest_Org.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUpOrgRequest_Org>): SetUpOrgRequest_Org {
    const message = createBaseSetUpOrgRequest_Org();
    message.name = object.name ?? "";
    message.domain = object.domain ?? "";
    return message;
  },
};

function createBaseSetUpOrgRequest_Human(): SetUpOrgRequest_Human {
  return { userName: "", profile: undefined, email: undefined, phone: undefined, password: "" };
}

export const SetUpOrgRequest_Human = {
  encode(message: SetUpOrgRequest_Human, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userName !== "") {
      writer.uint32(10).string(message.userName);
    }
    if (message.profile !== undefined) {
      SetUpOrgRequest_Human_Profile.encode(message.profile, writer.uint32(18).fork()).ldelim();
    }
    if (message.email !== undefined) {
      SetUpOrgRequest_Human_Email.encode(message.email, writer.uint32(26).fork()).ldelim();
    }
    if (message.phone !== undefined) {
      SetUpOrgRequest_Human_Phone.encode(message.phone, writer.uint32(34).fork()).ldelim();
    }
    if (message.password !== "") {
      writer.uint32(42).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUpOrgRequest_Human {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUpOrgRequest_Human();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userName = reader.string();
          break;
        case 2:
          message.profile = SetUpOrgRequest_Human_Profile.decode(reader, reader.uint32());
          break;
        case 3:
          message.email = SetUpOrgRequest_Human_Email.decode(reader, reader.uint32());
          break;
        case 4:
          message.phone = SetUpOrgRequest_Human_Phone.decode(reader, reader.uint32());
          break;
        case 5:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetUpOrgRequest_Human {
    return {
      userName: isSet(object.userName) ? String(object.userName) : "",
      profile: isSet(object.profile) ? SetUpOrgRequest_Human_Profile.fromJSON(object.profile) : undefined,
      email: isSet(object.email) ? SetUpOrgRequest_Human_Email.fromJSON(object.email) : undefined,
      phone: isSet(object.phone) ? SetUpOrgRequest_Human_Phone.fromJSON(object.phone) : undefined,
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: SetUpOrgRequest_Human): unknown {
    const obj: any = {};
    message.userName !== undefined && (obj.userName = message.userName);
    message.profile !== undefined &&
      (obj.profile = message.profile ? SetUpOrgRequest_Human_Profile.toJSON(message.profile) : undefined);
    message.email !== undefined &&
      (obj.email = message.email ? SetUpOrgRequest_Human_Email.toJSON(message.email) : undefined);
    message.phone !== undefined &&
      (obj.phone = message.phone ? SetUpOrgRequest_Human_Phone.toJSON(message.phone) : undefined);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  create(base?: DeepPartial<SetUpOrgRequest_Human>): SetUpOrgRequest_Human {
    return SetUpOrgRequest_Human.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUpOrgRequest_Human>): SetUpOrgRequest_Human {
    const message = createBaseSetUpOrgRequest_Human();
    message.userName = object.userName ?? "";
    message.profile = (object.profile !== undefined && object.profile !== null)
      ? SetUpOrgRequest_Human_Profile.fromPartial(object.profile)
      : undefined;
    message.email = (object.email !== undefined && object.email !== null)
      ? SetUpOrgRequest_Human_Email.fromPartial(object.email)
      : undefined;
    message.phone = (object.phone !== undefined && object.phone !== null)
      ? SetUpOrgRequest_Human_Phone.fromPartial(object.phone)
      : undefined;
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseSetUpOrgRequest_Human_Profile(): SetUpOrgRequest_Human_Profile {
  return { firstName: "", lastName: "", nickName: "", displayName: "", preferredLanguage: "", gender: 0 };
}

export const SetUpOrgRequest_Human_Profile = {
  encode(message: SetUpOrgRequest_Human_Profile, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.firstName !== "") {
      writer.uint32(10).string(message.firstName);
    }
    if (message.lastName !== "") {
      writer.uint32(18).string(message.lastName);
    }
    if (message.nickName !== "") {
      writer.uint32(26).string(message.nickName);
    }
    if (message.displayName !== "") {
      writer.uint32(34).string(message.displayName);
    }
    if (message.preferredLanguage !== "") {
      writer.uint32(42).string(message.preferredLanguage);
    }
    if (message.gender !== 0) {
      writer.uint32(48).int32(message.gender);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUpOrgRequest_Human_Profile {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUpOrgRequest_Human_Profile();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.firstName = reader.string();
          break;
        case 2:
          message.lastName = reader.string();
          break;
        case 3:
          message.nickName = reader.string();
          break;
        case 4:
          message.displayName = reader.string();
          break;
        case 5:
          message.preferredLanguage = reader.string();
          break;
        case 6:
          message.gender = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetUpOrgRequest_Human_Profile {
    return {
      firstName: isSet(object.firstName) ? String(object.firstName) : "",
      lastName: isSet(object.lastName) ? String(object.lastName) : "",
      nickName: isSet(object.nickName) ? String(object.nickName) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      preferredLanguage: isSet(object.preferredLanguage) ? String(object.preferredLanguage) : "",
      gender: isSet(object.gender) ? genderFromJSON(object.gender) : 0,
    };
  },

  toJSON(message: SetUpOrgRequest_Human_Profile): unknown {
    const obj: any = {};
    message.firstName !== undefined && (obj.firstName = message.firstName);
    message.lastName !== undefined && (obj.lastName = message.lastName);
    message.nickName !== undefined && (obj.nickName = message.nickName);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.preferredLanguage !== undefined && (obj.preferredLanguage = message.preferredLanguage);
    message.gender !== undefined && (obj.gender = genderToJSON(message.gender));
    return obj;
  },

  create(base?: DeepPartial<SetUpOrgRequest_Human_Profile>): SetUpOrgRequest_Human_Profile {
    return SetUpOrgRequest_Human_Profile.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUpOrgRequest_Human_Profile>): SetUpOrgRequest_Human_Profile {
    const message = createBaseSetUpOrgRequest_Human_Profile();
    message.firstName = object.firstName ?? "";
    message.lastName = object.lastName ?? "";
    message.nickName = object.nickName ?? "";
    message.displayName = object.displayName ?? "";
    message.preferredLanguage = object.preferredLanguage ?? "";
    message.gender = object.gender ?? 0;
    return message;
  },
};

function createBaseSetUpOrgRequest_Human_Email(): SetUpOrgRequest_Human_Email {
  return { email: "", isEmailVerified: false };
}

export const SetUpOrgRequest_Human_Email = {
  encode(message: SetUpOrgRequest_Human_Email, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.email !== "") {
      writer.uint32(10).string(message.email);
    }
    if (message.isEmailVerified === true) {
      writer.uint32(16).bool(message.isEmailVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUpOrgRequest_Human_Email {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUpOrgRequest_Human_Email();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.email = reader.string();
          break;
        case 2:
          message.isEmailVerified = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetUpOrgRequest_Human_Email {
    return {
      email: isSet(object.email) ? String(object.email) : "",
      isEmailVerified: isSet(object.isEmailVerified) ? Boolean(object.isEmailVerified) : false,
    };
  },

  toJSON(message: SetUpOrgRequest_Human_Email): unknown {
    const obj: any = {};
    message.email !== undefined && (obj.email = message.email);
    message.isEmailVerified !== undefined && (obj.isEmailVerified = message.isEmailVerified);
    return obj;
  },

  create(base?: DeepPartial<SetUpOrgRequest_Human_Email>): SetUpOrgRequest_Human_Email {
    return SetUpOrgRequest_Human_Email.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUpOrgRequest_Human_Email>): SetUpOrgRequest_Human_Email {
    const message = createBaseSetUpOrgRequest_Human_Email();
    message.email = object.email ?? "";
    message.isEmailVerified = object.isEmailVerified ?? false;
    return message;
  },
};

function createBaseSetUpOrgRequest_Human_Phone(): SetUpOrgRequest_Human_Phone {
  return { phone: "", isPhoneVerified: false };
}

export const SetUpOrgRequest_Human_Phone = {
  encode(message: SetUpOrgRequest_Human_Phone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.phone !== "") {
      writer.uint32(10).string(message.phone);
    }
    if (message.isPhoneVerified === true) {
      writer.uint32(16).bool(message.isPhoneVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUpOrgRequest_Human_Phone {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUpOrgRequest_Human_Phone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.phone = reader.string();
          break;
        case 2:
          message.isPhoneVerified = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetUpOrgRequest_Human_Phone {
    return {
      phone: isSet(object.phone) ? String(object.phone) : "",
      isPhoneVerified: isSet(object.isPhoneVerified) ? Boolean(object.isPhoneVerified) : false,
    };
  },

  toJSON(message: SetUpOrgRequest_Human_Phone): unknown {
    const obj: any = {};
    message.phone !== undefined && (obj.phone = message.phone);
    message.isPhoneVerified !== undefined && (obj.isPhoneVerified = message.isPhoneVerified);
    return obj;
  },

  create(base?: DeepPartial<SetUpOrgRequest_Human_Phone>): SetUpOrgRequest_Human_Phone {
    return SetUpOrgRequest_Human_Phone.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUpOrgRequest_Human_Phone>): SetUpOrgRequest_Human_Phone {
    const message = createBaseSetUpOrgRequest_Human_Phone();
    message.phone = object.phone ?? "";
    message.isPhoneVerified = object.isPhoneVerified ?? false;
    return message;
  },
};

function createBaseSetUpOrgResponse(): SetUpOrgResponse {
  return { details: undefined, orgId: "", userId: "" };
}

export const SetUpOrgResponse = {
  encode(message: SetUpOrgResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.orgId !== "") {
      writer.uint32(18).string(message.orgId);
    }
    if (message.userId !== "") {
      writer.uint32(26).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUpOrgResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUpOrgResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.orgId = reader.string();
          break;
        case 3:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetUpOrgResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: SetUpOrgResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<SetUpOrgResponse>): SetUpOrgResponse {
    return SetUpOrgResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUpOrgResponse>): SetUpOrgResponse {
    const message = createBaseSetUpOrgResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.orgId = object.orgId ?? "";
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseRemoveOrgRequest(): RemoveOrgRequest {
  return { orgId: "" };
}

export const RemoveOrgRequest = {
  encode(message: RemoveOrgRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveOrgRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveOrgRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveOrgRequest {
    return { orgId: isSet(object.orgId) ? String(object.orgId) : "" };
  },

  toJSON(message: RemoveOrgRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    return obj;
  },

  create(base?: DeepPartial<RemoveOrgRequest>): RemoveOrgRequest {
    return RemoveOrgRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveOrgRequest>): RemoveOrgRequest {
    const message = createBaseRemoveOrgRequest();
    message.orgId = object.orgId ?? "";
    return message;
  },
};

function createBaseRemoveOrgResponse(): RemoveOrgResponse {
  return { details: undefined };
}

export const RemoveOrgResponse = {
  encode(message: RemoveOrgResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveOrgResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveOrgResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveOrgResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveOrgResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveOrgResponse>): RemoveOrgResponse {
    return RemoveOrgResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveOrgResponse>): RemoveOrgResponse {
    const message = createBaseRemoveOrgResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetIDPByIDRequest(): GetIDPByIDRequest {
  return { id: "" };
}

export const GetIDPByIDRequest = {
  encode(message: GetIDPByIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetIDPByIDRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetIDPByIDRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetIDPByIDRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: GetIDPByIDRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<GetIDPByIDRequest>): GetIDPByIDRequest {
    return GetIDPByIDRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetIDPByIDRequest>): GetIDPByIDRequest {
    const message = createBaseGetIDPByIDRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseGetIDPByIDResponse(): GetIDPByIDResponse {
  return { idp: undefined };
}

export const GetIDPByIDResponse = {
  encode(message: GetIDPByIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idp !== undefined) {
      IDP.encode(message.idp, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetIDPByIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetIDPByIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idp = IDP.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetIDPByIDResponse {
    return { idp: isSet(object.idp) ? IDP.fromJSON(object.idp) : undefined };
  },

  toJSON(message: GetIDPByIDResponse): unknown {
    const obj: any = {};
    message.idp !== undefined && (obj.idp = message.idp ? IDP.toJSON(message.idp) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetIDPByIDResponse>): GetIDPByIDResponse {
    return GetIDPByIDResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetIDPByIDResponse>): GetIDPByIDResponse {
    const message = createBaseGetIDPByIDResponse();
    message.idp = (object.idp !== undefined && object.idp !== null) ? IDP.fromPartial(object.idp) : undefined;
    return message;
  },
};

function createBaseListIDPsRequest(): ListIDPsRequest {
  return { query: undefined, sortingColumn: 0, queries: [] };
}

export const ListIDPsRequest = {
  encode(message: ListIDPsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.queries) {
      IDPQuery.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListIDPsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListIDPsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.sortingColumn = reader.int32() as any;
          break;
        case 3:
          message.queries.push(IDPQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListIDPsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? iDPFieldNameFromJSON(object.sortingColumn) : 0,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => IDPQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListIDPsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = iDPFieldNameToJSON(message.sortingColumn));
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? IDPQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListIDPsRequest>): ListIDPsRequest {
    return ListIDPsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListIDPsRequest>): ListIDPsRequest {
    const message = createBaseListIDPsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.queries = object.queries?.map((e) => IDPQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseIDPQuery(): IDPQuery {
  return { idpIdQuery: undefined, idpNameQuery: undefined };
}

export const IDPQuery = {
  encode(message: IDPQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpIdQuery !== undefined) {
      IDPIDQuery.encode(message.idpIdQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.idpNameQuery !== undefined) {
      IDPNameQuery.encode(message.idpNameQuery, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPQuery {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpIdQuery = IDPIDQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.idpNameQuery = IDPNameQuery.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IDPQuery {
    return {
      idpIdQuery: isSet(object.idpIdQuery) ? IDPIDQuery.fromJSON(object.idpIdQuery) : undefined,
      idpNameQuery: isSet(object.idpNameQuery) ? IDPNameQuery.fromJSON(object.idpNameQuery) : undefined,
    };
  },

  toJSON(message: IDPQuery): unknown {
    const obj: any = {};
    message.idpIdQuery !== undefined &&
      (obj.idpIdQuery = message.idpIdQuery ? IDPIDQuery.toJSON(message.idpIdQuery) : undefined);
    message.idpNameQuery !== undefined &&
      (obj.idpNameQuery = message.idpNameQuery ? IDPNameQuery.toJSON(message.idpNameQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<IDPQuery>): IDPQuery {
    return IDPQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPQuery>): IDPQuery {
    const message = createBaseIDPQuery();
    message.idpIdQuery = (object.idpIdQuery !== undefined && object.idpIdQuery !== null)
      ? IDPIDQuery.fromPartial(object.idpIdQuery)
      : undefined;
    message.idpNameQuery = (object.idpNameQuery !== undefined && object.idpNameQuery !== null)
      ? IDPNameQuery.fromPartial(object.idpNameQuery)
      : undefined;
    return message;
  },
};

function createBaseListIDPsResponse(): ListIDPsResponse {
  return { details: undefined, sortingColumn: 0, result: [] };
}

export const ListIDPsResponse = {
  encode(message: ListIDPsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.result) {
      IDP.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListIDPsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListIDPsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.sortingColumn = reader.int32() as any;
          break;
        case 3:
          message.result.push(IDP.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListIDPsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? iDPFieldNameFromJSON(object.sortingColumn) : 0,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => IDP.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListIDPsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = iDPFieldNameToJSON(message.sortingColumn));
    if (message.result) {
      obj.result = message.result.map((e) => e ? IDP.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListIDPsResponse>): ListIDPsResponse {
    return ListIDPsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListIDPsResponse>): ListIDPsResponse {
    const message = createBaseListIDPsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.result = object.result?.map((e) => IDP.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAddOIDCIDPRequest(): AddOIDCIDPRequest {
  return {
    name: "",
    stylingType: 0,
    clientId: "",
    clientSecret: "",
    issuer: "",
    scopes: [],
    displayNameMapping: 0,
    usernameMapping: 0,
    autoRegister: false,
  };
}

export const AddOIDCIDPRequest = {
  encode(message: AddOIDCIDPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.stylingType !== 0) {
      writer.uint32(16).int32(message.stylingType);
    }
    if (message.clientId !== "") {
      writer.uint32(26).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(34).string(message.clientSecret);
    }
    if (message.issuer !== "") {
      writer.uint32(42).string(message.issuer);
    }
    for (const v of message.scopes) {
      writer.uint32(50).string(v!);
    }
    if (message.displayNameMapping !== 0) {
      writer.uint32(56).int32(message.displayNameMapping);
    }
    if (message.usernameMapping !== 0) {
      writer.uint32(64).int32(message.usernameMapping);
    }
    if (message.autoRegister === true) {
      writer.uint32(72).bool(message.autoRegister);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOIDCIDPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOIDCIDPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.stylingType = reader.int32() as any;
          break;
        case 3:
          message.clientId = reader.string();
          break;
        case 4:
          message.clientSecret = reader.string();
          break;
        case 5:
          message.issuer = reader.string();
          break;
        case 6:
          message.scopes.push(reader.string());
          break;
        case 7:
          message.displayNameMapping = reader.int32() as any;
          break;
        case 8:
          message.usernameMapping = reader.int32() as any;
          break;
        case 9:
          message.autoRegister = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddOIDCIDPRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      stylingType: isSet(object.stylingType) ? iDPStylingTypeFromJSON(object.stylingType) : 0,
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      displayNameMapping: isSet(object.displayNameMapping) ? oIDCMappingFieldFromJSON(object.displayNameMapping) : 0,
      usernameMapping: isSet(object.usernameMapping) ? oIDCMappingFieldFromJSON(object.usernameMapping) : 0,
      autoRegister: isSet(object.autoRegister) ? Boolean(object.autoRegister) : false,
    };
  },

  toJSON(message: AddOIDCIDPRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.stylingType !== undefined && (obj.stylingType = iDPStylingTypeToJSON(message.stylingType));
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.displayNameMapping !== undefined &&
      (obj.displayNameMapping = oIDCMappingFieldToJSON(message.displayNameMapping));
    message.usernameMapping !== undefined && (obj.usernameMapping = oIDCMappingFieldToJSON(message.usernameMapping));
    message.autoRegister !== undefined && (obj.autoRegister = message.autoRegister);
    return obj;
  },

  create(base?: DeepPartial<AddOIDCIDPRequest>): AddOIDCIDPRequest {
    return AddOIDCIDPRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOIDCIDPRequest>): AddOIDCIDPRequest {
    const message = createBaseAddOIDCIDPRequest();
    message.name = object.name ?? "";
    message.stylingType = object.stylingType ?? 0;
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.issuer = object.issuer ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.displayNameMapping = object.displayNameMapping ?? 0;
    message.usernameMapping = object.usernameMapping ?? 0;
    message.autoRegister = object.autoRegister ?? false;
    return message;
  },
};

function createBaseAddOIDCIDPResponse(): AddOIDCIDPResponse {
  return { details: undefined, idpId: "" };
}

export const AddOIDCIDPResponse = {
  encode(message: AddOIDCIDPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.idpId !== "") {
      writer.uint32(18).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOIDCIDPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOIDCIDPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.idpId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddOIDCIDPResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
    };
  },

  toJSON(message: AddOIDCIDPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<AddOIDCIDPResponse>): AddOIDCIDPResponse {
    return AddOIDCIDPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOIDCIDPResponse>): AddOIDCIDPResponse {
    const message = createBaseAddOIDCIDPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseAddJWTIDPRequest(): AddJWTIDPRequest {
  return {
    name: "",
    stylingType: 0,
    jwtEndpoint: "",
    issuer: "",
    keysEndpoint: "",
    headerName: "",
    autoRegister: false,
  };
}

export const AddJWTIDPRequest = {
  encode(message: AddJWTIDPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.stylingType !== 0) {
      writer.uint32(16).int32(message.stylingType);
    }
    if (message.jwtEndpoint !== "") {
      writer.uint32(26).string(message.jwtEndpoint);
    }
    if (message.issuer !== "") {
      writer.uint32(34).string(message.issuer);
    }
    if (message.keysEndpoint !== "") {
      writer.uint32(42).string(message.keysEndpoint);
    }
    if (message.headerName !== "") {
      writer.uint32(50).string(message.headerName);
    }
    if (message.autoRegister === true) {
      writer.uint32(56).bool(message.autoRegister);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddJWTIDPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddJWTIDPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.stylingType = reader.int32() as any;
          break;
        case 3:
          message.jwtEndpoint = reader.string();
          break;
        case 4:
          message.issuer = reader.string();
          break;
        case 5:
          message.keysEndpoint = reader.string();
          break;
        case 6:
          message.headerName = reader.string();
          break;
        case 7:
          message.autoRegister = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddJWTIDPRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      stylingType: isSet(object.stylingType) ? iDPStylingTypeFromJSON(object.stylingType) : 0,
      jwtEndpoint: isSet(object.jwtEndpoint) ? String(object.jwtEndpoint) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      keysEndpoint: isSet(object.keysEndpoint) ? String(object.keysEndpoint) : "",
      headerName: isSet(object.headerName) ? String(object.headerName) : "",
      autoRegister: isSet(object.autoRegister) ? Boolean(object.autoRegister) : false,
    };
  },

  toJSON(message: AddJWTIDPRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.stylingType !== undefined && (obj.stylingType = iDPStylingTypeToJSON(message.stylingType));
    message.jwtEndpoint !== undefined && (obj.jwtEndpoint = message.jwtEndpoint);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.keysEndpoint !== undefined && (obj.keysEndpoint = message.keysEndpoint);
    message.headerName !== undefined && (obj.headerName = message.headerName);
    message.autoRegister !== undefined && (obj.autoRegister = message.autoRegister);
    return obj;
  },

  create(base?: DeepPartial<AddJWTIDPRequest>): AddJWTIDPRequest {
    return AddJWTIDPRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddJWTIDPRequest>): AddJWTIDPRequest {
    const message = createBaseAddJWTIDPRequest();
    message.name = object.name ?? "";
    message.stylingType = object.stylingType ?? 0;
    message.jwtEndpoint = object.jwtEndpoint ?? "";
    message.issuer = object.issuer ?? "";
    message.keysEndpoint = object.keysEndpoint ?? "";
    message.headerName = object.headerName ?? "";
    message.autoRegister = object.autoRegister ?? false;
    return message;
  },
};

function createBaseAddJWTIDPResponse(): AddJWTIDPResponse {
  return { details: undefined, idpId: "" };
}

export const AddJWTIDPResponse = {
  encode(message: AddJWTIDPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.idpId !== "") {
      writer.uint32(18).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddJWTIDPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddJWTIDPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.idpId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddJWTIDPResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
    };
  },

  toJSON(message: AddJWTIDPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<AddJWTIDPResponse>): AddJWTIDPResponse {
    return AddJWTIDPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddJWTIDPResponse>): AddJWTIDPResponse {
    const message = createBaseAddJWTIDPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseUpdateIDPRequest(): UpdateIDPRequest {
  return { idpId: "", name: "", stylingType: 0, autoRegister: false };
}

export const UpdateIDPRequest = {
  encode(message: UpdateIDPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.stylingType !== 0) {
      writer.uint32(24).int32(message.stylingType);
    }
    if (message.autoRegister === true) {
      writer.uint32(32).bool(message.autoRegister);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIDPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIDPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.stylingType = reader.int32() as any;
          break;
        case 4:
          message.autoRegister = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateIDPRequest {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      name: isSet(object.name) ? String(object.name) : "",
      stylingType: isSet(object.stylingType) ? iDPStylingTypeFromJSON(object.stylingType) : 0,
      autoRegister: isSet(object.autoRegister) ? Boolean(object.autoRegister) : false,
    };
  },

  toJSON(message: UpdateIDPRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.name !== undefined && (obj.name = message.name);
    message.stylingType !== undefined && (obj.stylingType = iDPStylingTypeToJSON(message.stylingType));
    message.autoRegister !== undefined && (obj.autoRegister = message.autoRegister);
    return obj;
  },

  create(base?: DeepPartial<UpdateIDPRequest>): UpdateIDPRequest {
    return UpdateIDPRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateIDPRequest>): UpdateIDPRequest {
    const message = createBaseUpdateIDPRequest();
    message.idpId = object.idpId ?? "";
    message.name = object.name ?? "";
    message.stylingType = object.stylingType ?? 0;
    message.autoRegister = object.autoRegister ?? false;
    return message;
  },
};

function createBaseUpdateIDPResponse(): UpdateIDPResponse {
  return { details: undefined };
}

export const UpdateIDPResponse = {
  encode(message: UpdateIDPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIDPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIDPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateIDPResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateIDPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateIDPResponse>): UpdateIDPResponse {
    return UpdateIDPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateIDPResponse>): UpdateIDPResponse {
    const message = createBaseUpdateIDPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseDeactivateIDPRequest(): DeactivateIDPRequest {
  return { idpId: "" };
}

export const DeactivateIDPRequest = {
  encode(message: DeactivateIDPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeactivateIDPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeactivateIDPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeactivateIDPRequest {
    return { idpId: isSet(object.idpId) ? String(object.idpId) : "" };
  },

  toJSON(message: DeactivateIDPRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<DeactivateIDPRequest>): DeactivateIDPRequest {
    return DeactivateIDPRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeactivateIDPRequest>): DeactivateIDPRequest {
    const message = createBaseDeactivateIDPRequest();
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseDeactivateIDPResponse(): DeactivateIDPResponse {
  return { details: undefined };
}

export const DeactivateIDPResponse = {
  encode(message: DeactivateIDPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeactivateIDPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeactivateIDPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeactivateIDPResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeactivateIDPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeactivateIDPResponse>): DeactivateIDPResponse {
    return DeactivateIDPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeactivateIDPResponse>): DeactivateIDPResponse {
    const message = createBaseDeactivateIDPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseReactivateIDPRequest(): ReactivateIDPRequest {
  return { idpId: "" };
}

export const ReactivateIDPRequest = {
  encode(message: ReactivateIDPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReactivateIDPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReactivateIDPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ReactivateIDPRequest {
    return { idpId: isSet(object.idpId) ? String(object.idpId) : "" };
  },

  toJSON(message: ReactivateIDPRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<ReactivateIDPRequest>): ReactivateIDPRequest {
    return ReactivateIDPRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ReactivateIDPRequest>): ReactivateIDPRequest {
    const message = createBaseReactivateIDPRequest();
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseReactivateIDPResponse(): ReactivateIDPResponse {
  return { details: undefined };
}

export const ReactivateIDPResponse = {
  encode(message: ReactivateIDPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReactivateIDPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReactivateIDPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ReactivateIDPResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ReactivateIDPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ReactivateIDPResponse>): ReactivateIDPResponse {
    return ReactivateIDPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ReactivateIDPResponse>): ReactivateIDPResponse {
    const message = createBaseReactivateIDPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveIDPRequest(): RemoveIDPRequest {
  return { idpId: "" };
}

export const RemoveIDPRequest = {
  encode(message: RemoveIDPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveIDPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveIDPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveIDPRequest {
    return { idpId: isSet(object.idpId) ? String(object.idpId) : "" };
  },

  toJSON(message: RemoveIDPRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<RemoveIDPRequest>): RemoveIDPRequest {
    return RemoveIDPRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveIDPRequest>): RemoveIDPRequest {
    const message = createBaseRemoveIDPRequest();
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseRemoveIDPResponse(): RemoveIDPResponse {
  return { details: undefined };
}

export const RemoveIDPResponse = {
  encode(message: RemoveIDPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveIDPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveIDPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveIDPResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveIDPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveIDPResponse>): RemoveIDPResponse {
    return RemoveIDPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveIDPResponse>): RemoveIDPResponse {
    const message = createBaseRemoveIDPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateIDPOIDCConfigRequest(): UpdateIDPOIDCConfigRequest {
  return {
    idpId: "",
    issuer: "",
    clientId: "",
    clientSecret: "",
    scopes: [],
    displayNameMapping: 0,
    usernameMapping: 0,
  };
}

export const UpdateIDPOIDCConfigRequest = {
  encode(message: UpdateIDPOIDCConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.issuer !== "") {
      writer.uint32(18).string(message.issuer);
    }
    if (message.clientId !== "") {
      writer.uint32(26).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(34).string(message.clientSecret);
    }
    for (const v of message.scopes) {
      writer.uint32(42).string(v!);
    }
    if (message.displayNameMapping !== 0) {
      writer.uint32(48).int32(message.displayNameMapping);
    }
    if (message.usernameMapping !== 0) {
      writer.uint32(56).int32(message.usernameMapping);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIDPOIDCConfigRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIDPOIDCConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        case 2:
          message.issuer = reader.string();
          break;
        case 3:
          message.clientId = reader.string();
          break;
        case 4:
          message.clientSecret = reader.string();
          break;
        case 5:
          message.scopes.push(reader.string());
          break;
        case 6:
          message.displayNameMapping = reader.int32() as any;
          break;
        case 7:
          message.usernameMapping = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateIDPOIDCConfigRequest {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      displayNameMapping: isSet(object.displayNameMapping) ? oIDCMappingFieldFromJSON(object.displayNameMapping) : 0,
      usernameMapping: isSet(object.usernameMapping) ? oIDCMappingFieldFromJSON(object.usernameMapping) : 0,
    };
  },

  toJSON(message: UpdateIDPOIDCConfigRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.displayNameMapping !== undefined &&
      (obj.displayNameMapping = oIDCMappingFieldToJSON(message.displayNameMapping));
    message.usernameMapping !== undefined && (obj.usernameMapping = oIDCMappingFieldToJSON(message.usernameMapping));
    return obj;
  },

  create(base?: DeepPartial<UpdateIDPOIDCConfigRequest>): UpdateIDPOIDCConfigRequest {
    return UpdateIDPOIDCConfigRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateIDPOIDCConfigRequest>): UpdateIDPOIDCConfigRequest {
    const message = createBaseUpdateIDPOIDCConfigRequest();
    message.idpId = object.idpId ?? "";
    message.issuer = object.issuer ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.displayNameMapping = object.displayNameMapping ?? 0;
    message.usernameMapping = object.usernameMapping ?? 0;
    return message;
  },
};

function createBaseUpdateIDPOIDCConfigResponse(): UpdateIDPOIDCConfigResponse {
  return { details: undefined };
}

export const UpdateIDPOIDCConfigResponse = {
  encode(message: UpdateIDPOIDCConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIDPOIDCConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIDPOIDCConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateIDPOIDCConfigResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateIDPOIDCConfigResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateIDPOIDCConfigResponse>): UpdateIDPOIDCConfigResponse {
    return UpdateIDPOIDCConfigResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateIDPOIDCConfigResponse>): UpdateIDPOIDCConfigResponse {
    const message = createBaseUpdateIDPOIDCConfigResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateIDPJWTConfigRequest(): UpdateIDPJWTConfigRequest {
  return { idpId: "", jwtEndpoint: "", issuer: "", keysEndpoint: "", headerName: "" };
}

export const UpdateIDPJWTConfigRequest = {
  encode(message: UpdateIDPJWTConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.jwtEndpoint !== "") {
      writer.uint32(18).string(message.jwtEndpoint);
    }
    if (message.issuer !== "") {
      writer.uint32(26).string(message.issuer);
    }
    if (message.keysEndpoint !== "") {
      writer.uint32(34).string(message.keysEndpoint);
    }
    if (message.headerName !== "") {
      writer.uint32(42).string(message.headerName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIDPJWTConfigRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIDPJWTConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        case 2:
          message.jwtEndpoint = reader.string();
          break;
        case 3:
          message.issuer = reader.string();
          break;
        case 4:
          message.keysEndpoint = reader.string();
          break;
        case 5:
          message.headerName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateIDPJWTConfigRequest {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      jwtEndpoint: isSet(object.jwtEndpoint) ? String(object.jwtEndpoint) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      keysEndpoint: isSet(object.keysEndpoint) ? String(object.keysEndpoint) : "",
      headerName: isSet(object.headerName) ? String(object.headerName) : "",
    };
  },

  toJSON(message: UpdateIDPJWTConfigRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.jwtEndpoint !== undefined && (obj.jwtEndpoint = message.jwtEndpoint);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.keysEndpoint !== undefined && (obj.keysEndpoint = message.keysEndpoint);
    message.headerName !== undefined && (obj.headerName = message.headerName);
    return obj;
  },

  create(base?: DeepPartial<UpdateIDPJWTConfigRequest>): UpdateIDPJWTConfigRequest {
    return UpdateIDPJWTConfigRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateIDPJWTConfigRequest>): UpdateIDPJWTConfigRequest {
    const message = createBaseUpdateIDPJWTConfigRequest();
    message.idpId = object.idpId ?? "";
    message.jwtEndpoint = object.jwtEndpoint ?? "";
    message.issuer = object.issuer ?? "";
    message.keysEndpoint = object.keysEndpoint ?? "";
    message.headerName = object.headerName ?? "";
    return message;
  },
};

function createBaseUpdateIDPJWTConfigResponse(): UpdateIDPJWTConfigResponse {
  return { details: undefined };
}

export const UpdateIDPJWTConfigResponse = {
  encode(message: UpdateIDPJWTConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIDPJWTConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIDPJWTConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateIDPJWTConfigResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateIDPJWTConfigResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateIDPJWTConfigResponse>): UpdateIDPJWTConfigResponse {
    return UpdateIDPJWTConfigResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateIDPJWTConfigResponse>): UpdateIDPJWTConfigResponse {
    const message = createBaseUpdateIDPJWTConfigResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListProvidersRequest(): ListProvidersRequest {
  return { query: undefined, queries: [] };
}

export const ListProvidersRequest = {
  encode(message: ListProvidersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.queries) {
      ProviderQuery.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListProvidersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListProvidersRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.queries.push(ProviderQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListProvidersRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => ProviderQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListProvidersRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? ProviderQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListProvidersRequest>): ListProvidersRequest {
    return ListProvidersRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListProvidersRequest>): ListProvidersRequest {
    const message = createBaseListProvidersRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.queries = object.queries?.map((e) => ProviderQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseProviderQuery(): ProviderQuery {
  return { idpIdQuery: undefined, idpNameQuery: undefined };
}

export const ProviderQuery = {
  encode(message: ProviderQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpIdQuery !== undefined) {
      IDPIDQuery.encode(message.idpIdQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.idpNameQuery !== undefined) {
      IDPNameQuery.encode(message.idpNameQuery, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProviderQuery {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProviderQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpIdQuery = IDPIDQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.idpNameQuery = IDPNameQuery.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProviderQuery {
    return {
      idpIdQuery: isSet(object.idpIdQuery) ? IDPIDQuery.fromJSON(object.idpIdQuery) : undefined,
      idpNameQuery: isSet(object.idpNameQuery) ? IDPNameQuery.fromJSON(object.idpNameQuery) : undefined,
    };
  },

  toJSON(message: ProviderQuery): unknown {
    const obj: any = {};
    message.idpIdQuery !== undefined &&
      (obj.idpIdQuery = message.idpIdQuery ? IDPIDQuery.toJSON(message.idpIdQuery) : undefined);
    message.idpNameQuery !== undefined &&
      (obj.idpNameQuery = message.idpNameQuery ? IDPNameQuery.toJSON(message.idpNameQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ProviderQuery>): ProviderQuery {
    return ProviderQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProviderQuery>): ProviderQuery {
    const message = createBaseProviderQuery();
    message.idpIdQuery = (object.idpIdQuery !== undefined && object.idpIdQuery !== null)
      ? IDPIDQuery.fromPartial(object.idpIdQuery)
      : undefined;
    message.idpNameQuery = (object.idpNameQuery !== undefined && object.idpNameQuery !== null)
      ? IDPNameQuery.fromPartial(object.idpNameQuery)
      : undefined;
    return message;
  },
};

function createBaseListProvidersResponse(): ListProvidersResponse {
  return { details: undefined, result: [] };
}

export const ListProvidersResponse = {
  encode(message: ListProvidersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      Provider.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListProvidersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListProvidersResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(Provider.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListProvidersResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Provider.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListProvidersResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? Provider.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListProvidersResponse>): ListProvidersResponse {
    return ListProvidersResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListProvidersResponse>): ListProvidersResponse {
    const message = createBaseListProvidersResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => Provider.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetProviderByIDRequest(): GetProviderByIDRequest {
  return { id: "" };
}

export const GetProviderByIDRequest = {
  encode(message: GetProviderByIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetProviderByIDRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetProviderByIDRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetProviderByIDRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: GetProviderByIDRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<GetProviderByIDRequest>): GetProviderByIDRequest {
    return GetProviderByIDRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetProviderByIDRequest>): GetProviderByIDRequest {
    const message = createBaseGetProviderByIDRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseGetProviderByIDResponse(): GetProviderByIDResponse {
  return { idp: undefined };
}

export const GetProviderByIDResponse = {
  encode(message: GetProviderByIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idp !== undefined) {
      Provider.encode(message.idp, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetProviderByIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetProviderByIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idp = Provider.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetProviderByIDResponse {
    return { idp: isSet(object.idp) ? Provider.fromJSON(object.idp) : undefined };
  },

  toJSON(message: GetProviderByIDResponse): unknown {
    const obj: any = {};
    message.idp !== undefined && (obj.idp = message.idp ? Provider.toJSON(message.idp) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetProviderByIDResponse>): GetProviderByIDResponse {
    return GetProviderByIDResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetProviderByIDResponse>): GetProviderByIDResponse {
    const message = createBaseGetProviderByIDResponse();
    message.idp = (object.idp !== undefined && object.idp !== null) ? Provider.fromPartial(object.idp) : undefined;
    return message;
  },
};

function createBaseAddGenericOAuthProviderRequest(): AddGenericOAuthProviderRequest {
  return {
    name: "",
    clientId: "",
    clientSecret: "",
    authorizationEndpoint: "",
    tokenEndpoint: "",
    userEndpoint: "",
    scopes: [],
    idAttribute: "",
    providerOptions: undefined,
  };
}

export const AddGenericOAuthProviderRequest = {
  encode(message: AddGenericOAuthProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.clientId !== "") {
      writer.uint32(18).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(26).string(message.clientSecret);
    }
    if (message.authorizationEndpoint !== "") {
      writer.uint32(34).string(message.authorizationEndpoint);
    }
    if (message.tokenEndpoint !== "") {
      writer.uint32(42).string(message.tokenEndpoint);
    }
    if (message.userEndpoint !== "") {
      writer.uint32(50).string(message.userEndpoint);
    }
    for (const v of message.scopes) {
      writer.uint32(58).string(v!);
    }
    if (message.idAttribute !== "") {
      writer.uint32(66).string(message.idAttribute);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGenericOAuthProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGenericOAuthProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.clientId = reader.string();
          break;
        case 3:
          message.clientSecret = reader.string();
          break;
        case 4:
          message.authorizationEndpoint = reader.string();
          break;
        case 5:
          message.tokenEndpoint = reader.string();
          break;
        case 6:
          message.userEndpoint = reader.string();
          break;
        case 7:
          message.scopes.push(reader.string());
          break;
        case 8:
          message.idAttribute = reader.string();
          break;
        case 9:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGenericOAuthProviderRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      authorizationEndpoint: isSet(object.authorizationEndpoint) ? String(object.authorizationEndpoint) : "",
      tokenEndpoint: isSet(object.tokenEndpoint) ? String(object.tokenEndpoint) : "",
      userEndpoint: isSet(object.userEndpoint) ? String(object.userEndpoint) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      idAttribute: isSet(object.idAttribute) ? String(object.idAttribute) : "",
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: AddGenericOAuthProviderRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    message.authorizationEndpoint !== undefined && (obj.authorizationEndpoint = message.authorizationEndpoint);
    message.tokenEndpoint !== undefined && (obj.tokenEndpoint = message.tokenEndpoint);
    message.userEndpoint !== undefined && (obj.userEndpoint = message.userEndpoint);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.idAttribute !== undefined && (obj.idAttribute = message.idAttribute);
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddGenericOAuthProviderRequest>): AddGenericOAuthProviderRequest {
    return AddGenericOAuthProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGenericOAuthProviderRequest>): AddGenericOAuthProviderRequest {
    const message = createBaseAddGenericOAuthProviderRequest();
    message.name = object.name ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.authorizationEndpoint = object.authorizationEndpoint ?? "";
    message.tokenEndpoint = object.tokenEndpoint ?? "";
    message.userEndpoint = object.userEndpoint ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.idAttribute = object.idAttribute ?? "";
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseAddGenericOAuthProviderResponse(): AddGenericOAuthProviderResponse {
  return { details: undefined, id: "" };
}

export const AddGenericOAuthProviderResponse = {
  encode(message: AddGenericOAuthProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGenericOAuthProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGenericOAuthProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGenericOAuthProviderResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
    };
  },

  toJSON(message: AddGenericOAuthProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<AddGenericOAuthProviderResponse>): AddGenericOAuthProviderResponse {
    return AddGenericOAuthProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGenericOAuthProviderResponse>): AddGenericOAuthProviderResponse {
    const message = createBaseAddGenericOAuthProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseUpdateGenericOAuthProviderRequest(): UpdateGenericOAuthProviderRequest {
  return {
    id: "",
    name: "",
    clientId: "",
    clientSecret: "",
    authorizationEndpoint: "",
    tokenEndpoint: "",
    userEndpoint: "",
    scopes: [],
    idAttribute: "",
    providerOptions: undefined,
  };
}

export const UpdateGenericOAuthProviderRequest = {
  encode(message: UpdateGenericOAuthProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.clientId !== "") {
      writer.uint32(26).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(34).string(message.clientSecret);
    }
    if (message.authorizationEndpoint !== "") {
      writer.uint32(42).string(message.authorizationEndpoint);
    }
    if (message.tokenEndpoint !== "") {
      writer.uint32(50).string(message.tokenEndpoint);
    }
    if (message.userEndpoint !== "") {
      writer.uint32(58).string(message.userEndpoint);
    }
    for (const v of message.scopes) {
      writer.uint32(66).string(v!);
    }
    if (message.idAttribute !== "") {
      writer.uint32(74).string(message.idAttribute);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(82).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGenericOAuthProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGenericOAuthProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.clientId = reader.string();
          break;
        case 4:
          message.clientSecret = reader.string();
          break;
        case 5:
          message.authorizationEndpoint = reader.string();
          break;
        case 6:
          message.tokenEndpoint = reader.string();
          break;
        case 7:
          message.userEndpoint = reader.string();
          break;
        case 8:
          message.scopes.push(reader.string());
          break;
        case 9:
          message.idAttribute = reader.string();
          break;
        case 10:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGenericOAuthProviderRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      authorizationEndpoint: isSet(object.authorizationEndpoint) ? String(object.authorizationEndpoint) : "",
      tokenEndpoint: isSet(object.tokenEndpoint) ? String(object.tokenEndpoint) : "",
      userEndpoint: isSet(object.userEndpoint) ? String(object.userEndpoint) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      idAttribute: isSet(object.idAttribute) ? String(object.idAttribute) : "",
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: UpdateGenericOAuthProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    message.authorizationEndpoint !== undefined && (obj.authorizationEndpoint = message.authorizationEndpoint);
    message.tokenEndpoint !== undefined && (obj.tokenEndpoint = message.tokenEndpoint);
    message.userEndpoint !== undefined && (obj.userEndpoint = message.userEndpoint);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.idAttribute !== undefined && (obj.idAttribute = message.idAttribute);
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGenericOAuthProviderRequest>): UpdateGenericOAuthProviderRequest {
    return UpdateGenericOAuthProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateGenericOAuthProviderRequest>): UpdateGenericOAuthProviderRequest {
    const message = createBaseUpdateGenericOAuthProviderRequest();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.authorizationEndpoint = object.authorizationEndpoint ?? "";
    message.tokenEndpoint = object.tokenEndpoint ?? "";
    message.userEndpoint = object.userEndpoint ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.idAttribute = object.idAttribute ?? "";
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseUpdateGenericOAuthProviderResponse(): UpdateGenericOAuthProviderResponse {
  return { details: undefined };
}

export const UpdateGenericOAuthProviderResponse = {
  encode(message: UpdateGenericOAuthProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGenericOAuthProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGenericOAuthProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGenericOAuthProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateGenericOAuthProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGenericOAuthProviderResponse>): UpdateGenericOAuthProviderResponse {
    return UpdateGenericOAuthProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateGenericOAuthProviderResponse>): UpdateGenericOAuthProviderResponse {
    const message = createBaseUpdateGenericOAuthProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddGenericOIDCProviderRequest(): AddGenericOIDCProviderRequest {
  return { name: "", issuer: "", clientId: "", clientSecret: "", scopes: [], providerOptions: undefined };
}

export const AddGenericOIDCProviderRequest = {
  encode(message: AddGenericOIDCProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.issuer !== "") {
      writer.uint32(18).string(message.issuer);
    }
    if (message.clientId !== "") {
      writer.uint32(26).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(34).string(message.clientSecret);
    }
    for (const v of message.scopes) {
      writer.uint32(42).string(v!);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGenericOIDCProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGenericOIDCProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.issuer = reader.string();
          break;
        case 3:
          message.clientId = reader.string();
          break;
        case 4:
          message.clientSecret = reader.string();
          break;
        case 5:
          message.scopes.push(reader.string());
          break;
        case 6:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGenericOIDCProviderRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: AddGenericOIDCProviderRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddGenericOIDCProviderRequest>): AddGenericOIDCProviderRequest {
    return AddGenericOIDCProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGenericOIDCProviderRequest>): AddGenericOIDCProviderRequest {
    const message = createBaseAddGenericOIDCProviderRequest();
    message.name = object.name ?? "";
    message.issuer = object.issuer ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseAddGenericOIDCProviderResponse(): AddGenericOIDCProviderResponse {
  return { details: undefined, id: "" };
}

export const AddGenericOIDCProviderResponse = {
  encode(message: AddGenericOIDCProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGenericOIDCProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGenericOIDCProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGenericOIDCProviderResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
    };
  },

  toJSON(message: AddGenericOIDCProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<AddGenericOIDCProviderResponse>): AddGenericOIDCProviderResponse {
    return AddGenericOIDCProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGenericOIDCProviderResponse>): AddGenericOIDCProviderResponse {
    const message = createBaseAddGenericOIDCProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseUpdateGenericOIDCProviderRequest(): UpdateGenericOIDCProviderRequest {
  return { id: "", name: "", issuer: "", clientId: "", clientSecret: "", scopes: [], providerOptions: undefined };
}

export const UpdateGenericOIDCProviderRequest = {
  encode(message: UpdateGenericOIDCProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.issuer !== "") {
      writer.uint32(26).string(message.issuer);
    }
    if (message.clientId !== "") {
      writer.uint32(34).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(42).string(message.clientSecret);
    }
    for (const v of message.scopes) {
      writer.uint32(50).string(v!);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGenericOIDCProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGenericOIDCProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.issuer = reader.string();
          break;
        case 4:
          message.clientId = reader.string();
          break;
        case 5:
          message.clientSecret = reader.string();
          break;
        case 6:
          message.scopes.push(reader.string());
          break;
        case 7:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGenericOIDCProviderRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: UpdateGenericOIDCProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGenericOIDCProviderRequest>): UpdateGenericOIDCProviderRequest {
    return UpdateGenericOIDCProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateGenericOIDCProviderRequest>): UpdateGenericOIDCProviderRequest {
    const message = createBaseUpdateGenericOIDCProviderRequest();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.issuer = object.issuer ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseUpdateGenericOIDCProviderResponse(): UpdateGenericOIDCProviderResponse {
  return { details: undefined };
}

export const UpdateGenericOIDCProviderResponse = {
  encode(message: UpdateGenericOIDCProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGenericOIDCProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGenericOIDCProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGenericOIDCProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateGenericOIDCProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGenericOIDCProviderResponse>): UpdateGenericOIDCProviderResponse {
    return UpdateGenericOIDCProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateGenericOIDCProviderResponse>): UpdateGenericOIDCProviderResponse {
    const message = createBaseUpdateGenericOIDCProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddJWTProviderRequest(): AddJWTProviderRequest {
  return { name: "", issuer: "", jwtEndpoint: "", keysEndpoint: "", headerName: "", providerOptions: undefined };
}

export const AddJWTProviderRequest = {
  encode(message: AddJWTProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.issuer !== "") {
      writer.uint32(18).string(message.issuer);
    }
    if (message.jwtEndpoint !== "") {
      writer.uint32(26).string(message.jwtEndpoint);
    }
    if (message.keysEndpoint !== "") {
      writer.uint32(34).string(message.keysEndpoint);
    }
    if (message.headerName !== "") {
      writer.uint32(42).string(message.headerName);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddJWTProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddJWTProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.issuer = reader.string();
          break;
        case 3:
          message.jwtEndpoint = reader.string();
          break;
        case 4:
          message.keysEndpoint = reader.string();
          break;
        case 5:
          message.headerName = reader.string();
          break;
        case 6:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddJWTProviderRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      jwtEndpoint: isSet(object.jwtEndpoint) ? String(object.jwtEndpoint) : "",
      keysEndpoint: isSet(object.keysEndpoint) ? String(object.keysEndpoint) : "",
      headerName: isSet(object.headerName) ? String(object.headerName) : "",
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: AddJWTProviderRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.jwtEndpoint !== undefined && (obj.jwtEndpoint = message.jwtEndpoint);
    message.keysEndpoint !== undefined && (obj.keysEndpoint = message.keysEndpoint);
    message.headerName !== undefined && (obj.headerName = message.headerName);
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddJWTProviderRequest>): AddJWTProviderRequest {
    return AddJWTProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddJWTProviderRequest>): AddJWTProviderRequest {
    const message = createBaseAddJWTProviderRequest();
    message.name = object.name ?? "";
    message.issuer = object.issuer ?? "";
    message.jwtEndpoint = object.jwtEndpoint ?? "";
    message.keysEndpoint = object.keysEndpoint ?? "";
    message.headerName = object.headerName ?? "";
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseAddJWTProviderResponse(): AddJWTProviderResponse {
  return { details: undefined, id: "" };
}

export const AddJWTProviderResponse = {
  encode(message: AddJWTProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddJWTProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddJWTProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddJWTProviderResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
    };
  },

  toJSON(message: AddJWTProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<AddJWTProviderResponse>): AddJWTProviderResponse {
    return AddJWTProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddJWTProviderResponse>): AddJWTProviderResponse {
    const message = createBaseAddJWTProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseUpdateJWTProviderRequest(): UpdateJWTProviderRequest {
  return {
    id: "",
    name: "",
    issuer: "",
    jwtEndpoint: "",
    keysEndpoint: "",
    headerName: "",
    providerOptions: undefined,
  };
}

export const UpdateJWTProviderRequest = {
  encode(message: UpdateJWTProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.issuer !== "") {
      writer.uint32(26).string(message.issuer);
    }
    if (message.jwtEndpoint !== "") {
      writer.uint32(34).string(message.jwtEndpoint);
    }
    if (message.keysEndpoint !== "") {
      writer.uint32(42).string(message.keysEndpoint);
    }
    if (message.headerName !== "") {
      writer.uint32(50).string(message.headerName);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateJWTProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateJWTProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.issuer = reader.string();
          break;
        case 4:
          message.jwtEndpoint = reader.string();
          break;
        case 5:
          message.keysEndpoint = reader.string();
          break;
        case 6:
          message.headerName = reader.string();
          break;
        case 7:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateJWTProviderRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      jwtEndpoint: isSet(object.jwtEndpoint) ? String(object.jwtEndpoint) : "",
      keysEndpoint: isSet(object.keysEndpoint) ? String(object.keysEndpoint) : "",
      headerName: isSet(object.headerName) ? String(object.headerName) : "",
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: UpdateJWTProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.jwtEndpoint !== undefined && (obj.jwtEndpoint = message.jwtEndpoint);
    message.keysEndpoint !== undefined && (obj.keysEndpoint = message.keysEndpoint);
    message.headerName !== undefined && (obj.headerName = message.headerName);
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateJWTProviderRequest>): UpdateJWTProviderRequest {
    return UpdateJWTProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateJWTProviderRequest>): UpdateJWTProviderRequest {
    const message = createBaseUpdateJWTProviderRequest();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.issuer = object.issuer ?? "";
    message.jwtEndpoint = object.jwtEndpoint ?? "";
    message.keysEndpoint = object.keysEndpoint ?? "";
    message.headerName = object.headerName ?? "";
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseUpdateJWTProviderResponse(): UpdateJWTProviderResponse {
  return { details: undefined };
}

export const UpdateJWTProviderResponse = {
  encode(message: UpdateJWTProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateJWTProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateJWTProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateJWTProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateJWTProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateJWTProviderResponse>): UpdateJWTProviderResponse {
    return UpdateJWTProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateJWTProviderResponse>): UpdateJWTProviderResponse {
    const message = createBaseUpdateJWTProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddGitHubProviderRequest(): AddGitHubProviderRequest {
  return { name: "", clientId: "", clientSecret: "", scopes: [], providerOptions: undefined };
}

export const AddGitHubProviderRequest = {
  encode(message: AddGitHubProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.clientId !== "") {
      writer.uint32(18).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(26).string(message.clientSecret);
    }
    for (const v of message.scopes) {
      writer.uint32(34).string(v!);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGitHubProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGitHubProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.clientId = reader.string();
          break;
        case 3:
          message.clientSecret = reader.string();
          break;
        case 4:
          message.scopes.push(reader.string());
          break;
        case 5:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGitHubProviderRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: AddGitHubProviderRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddGitHubProviderRequest>): AddGitHubProviderRequest {
    return AddGitHubProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGitHubProviderRequest>): AddGitHubProviderRequest {
    const message = createBaseAddGitHubProviderRequest();
    message.name = object.name ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseAddGitHubProviderResponse(): AddGitHubProviderResponse {
  return { details: undefined, id: "" };
}

export const AddGitHubProviderResponse = {
  encode(message: AddGitHubProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGitHubProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGitHubProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGitHubProviderResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
    };
  },

  toJSON(message: AddGitHubProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<AddGitHubProviderResponse>): AddGitHubProviderResponse {
    return AddGitHubProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGitHubProviderResponse>): AddGitHubProviderResponse {
    const message = createBaseAddGitHubProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseUpdateGitHubProviderRequest(): UpdateGitHubProviderRequest {
  return { id: "", name: "", clientId: "", clientSecret: "", scopes: [], providerOptions: undefined };
}

export const UpdateGitHubProviderRequest = {
  encode(message: UpdateGitHubProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.clientId !== "") {
      writer.uint32(26).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(34).string(message.clientSecret);
    }
    for (const v of message.scopes) {
      writer.uint32(42).string(v!);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGitHubProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGitHubProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.clientId = reader.string();
          break;
        case 4:
          message.clientSecret = reader.string();
          break;
        case 5:
          message.scopes.push(reader.string());
          break;
        case 6:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGitHubProviderRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: UpdateGitHubProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGitHubProviderRequest>): UpdateGitHubProviderRequest {
    return UpdateGitHubProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateGitHubProviderRequest>): UpdateGitHubProviderRequest {
    const message = createBaseUpdateGitHubProviderRequest();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseUpdateGitHubProviderResponse(): UpdateGitHubProviderResponse {
  return { details: undefined };
}

export const UpdateGitHubProviderResponse = {
  encode(message: UpdateGitHubProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGitHubProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGitHubProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGitHubProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateGitHubProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGitHubProviderResponse>): UpdateGitHubProviderResponse {
    return UpdateGitHubProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateGitHubProviderResponse>): UpdateGitHubProviderResponse {
    const message = createBaseUpdateGitHubProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddGitHubEnterpriseServerProviderRequest(): AddGitHubEnterpriseServerProviderRequest {
  return {
    clientId: "",
    name: "",
    clientSecret: "",
    authorizationEndpoint: "",
    tokenEndpoint: "",
    userEndpoint: "",
    scopes: [],
    providerOptions: undefined,
  };
}

export const AddGitHubEnterpriseServerProviderRequest = {
  encode(message: AddGitHubEnterpriseServerProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.clientSecret !== "") {
      writer.uint32(26).string(message.clientSecret);
    }
    if (message.authorizationEndpoint !== "") {
      writer.uint32(34).string(message.authorizationEndpoint);
    }
    if (message.tokenEndpoint !== "") {
      writer.uint32(42).string(message.tokenEndpoint);
    }
    if (message.userEndpoint !== "") {
      writer.uint32(50).string(message.userEndpoint);
    }
    for (const v of message.scopes) {
      writer.uint32(58).string(v!);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGitHubEnterpriseServerProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGitHubEnterpriseServerProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.clientId = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.clientSecret = reader.string();
          break;
        case 4:
          message.authorizationEndpoint = reader.string();
          break;
        case 5:
          message.tokenEndpoint = reader.string();
          break;
        case 6:
          message.userEndpoint = reader.string();
          break;
        case 7:
          message.scopes.push(reader.string());
          break;
        case 8:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGitHubEnterpriseServerProviderRequest {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      name: isSet(object.name) ? String(object.name) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      authorizationEndpoint: isSet(object.authorizationEndpoint) ? String(object.authorizationEndpoint) : "",
      tokenEndpoint: isSet(object.tokenEndpoint) ? String(object.tokenEndpoint) : "",
      userEndpoint: isSet(object.userEndpoint) ? String(object.userEndpoint) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: AddGitHubEnterpriseServerProviderRequest): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.name !== undefined && (obj.name = message.name);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    message.authorizationEndpoint !== undefined && (obj.authorizationEndpoint = message.authorizationEndpoint);
    message.tokenEndpoint !== undefined && (obj.tokenEndpoint = message.tokenEndpoint);
    message.userEndpoint !== undefined && (obj.userEndpoint = message.userEndpoint);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddGitHubEnterpriseServerProviderRequest>): AddGitHubEnterpriseServerProviderRequest {
    return AddGitHubEnterpriseServerProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGitHubEnterpriseServerProviderRequest>): AddGitHubEnterpriseServerProviderRequest {
    const message = createBaseAddGitHubEnterpriseServerProviderRequest();
    message.clientId = object.clientId ?? "";
    message.name = object.name ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.authorizationEndpoint = object.authorizationEndpoint ?? "";
    message.tokenEndpoint = object.tokenEndpoint ?? "";
    message.userEndpoint = object.userEndpoint ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseAddGitHubEnterpriseServerProviderResponse(): AddGitHubEnterpriseServerProviderResponse {
  return { details: undefined, id: "" };
}

export const AddGitHubEnterpriseServerProviderResponse = {
  encode(message: AddGitHubEnterpriseServerProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGitHubEnterpriseServerProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGitHubEnterpriseServerProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGitHubEnterpriseServerProviderResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
    };
  },

  toJSON(message: AddGitHubEnterpriseServerProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<AddGitHubEnterpriseServerProviderResponse>): AddGitHubEnterpriseServerProviderResponse {
    return AddGitHubEnterpriseServerProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<AddGitHubEnterpriseServerProviderResponse>,
  ): AddGitHubEnterpriseServerProviderResponse {
    const message = createBaseAddGitHubEnterpriseServerProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseUpdateGitHubEnterpriseServerProviderRequest(): UpdateGitHubEnterpriseServerProviderRequest {
  return {
    id: "",
    name: "",
    clientId: "",
    clientSecret: "",
    authorizationEndpoint: "",
    tokenEndpoint: "",
    userEndpoint: "",
    scopes: [],
    providerOptions: undefined,
  };
}

export const UpdateGitHubEnterpriseServerProviderRequest = {
  encode(message: UpdateGitHubEnterpriseServerProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.clientId !== "") {
      writer.uint32(26).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(34).string(message.clientSecret);
    }
    if (message.authorizationEndpoint !== "") {
      writer.uint32(42).string(message.authorizationEndpoint);
    }
    if (message.tokenEndpoint !== "") {
      writer.uint32(50).string(message.tokenEndpoint);
    }
    if (message.userEndpoint !== "") {
      writer.uint32(58).string(message.userEndpoint);
    }
    for (const v of message.scopes) {
      writer.uint32(66).string(v!);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGitHubEnterpriseServerProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGitHubEnterpriseServerProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.clientId = reader.string();
          break;
        case 4:
          message.clientSecret = reader.string();
          break;
        case 5:
          message.authorizationEndpoint = reader.string();
          break;
        case 6:
          message.tokenEndpoint = reader.string();
          break;
        case 7:
          message.userEndpoint = reader.string();
          break;
        case 8:
          message.scopes.push(reader.string());
          break;
        case 9:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGitHubEnterpriseServerProviderRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      authorizationEndpoint: isSet(object.authorizationEndpoint) ? String(object.authorizationEndpoint) : "",
      tokenEndpoint: isSet(object.tokenEndpoint) ? String(object.tokenEndpoint) : "",
      userEndpoint: isSet(object.userEndpoint) ? String(object.userEndpoint) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: UpdateGitHubEnterpriseServerProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    message.authorizationEndpoint !== undefined && (obj.authorizationEndpoint = message.authorizationEndpoint);
    message.tokenEndpoint !== undefined && (obj.tokenEndpoint = message.tokenEndpoint);
    message.userEndpoint !== undefined && (obj.userEndpoint = message.userEndpoint);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGitHubEnterpriseServerProviderRequest>): UpdateGitHubEnterpriseServerProviderRequest {
    return UpdateGitHubEnterpriseServerProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<UpdateGitHubEnterpriseServerProviderRequest>,
  ): UpdateGitHubEnterpriseServerProviderRequest {
    const message = createBaseUpdateGitHubEnterpriseServerProviderRequest();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.authorizationEndpoint = object.authorizationEndpoint ?? "";
    message.tokenEndpoint = object.tokenEndpoint ?? "";
    message.userEndpoint = object.userEndpoint ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseUpdateGitHubEnterpriseServerProviderResponse(): UpdateGitHubEnterpriseServerProviderResponse {
  return { details: undefined };
}

export const UpdateGitHubEnterpriseServerProviderResponse = {
  encode(message: UpdateGitHubEnterpriseServerProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGitHubEnterpriseServerProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGitHubEnterpriseServerProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGitHubEnterpriseServerProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateGitHubEnterpriseServerProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<UpdateGitHubEnterpriseServerProviderResponse>,
  ): UpdateGitHubEnterpriseServerProviderResponse {
    return UpdateGitHubEnterpriseServerProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<UpdateGitHubEnterpriseServerProviderResponse>,
  ): UpdateGitHubEnterpriseServerProviderResponse {
    const message = createBaseUpdateGitHubEnterpriseServerProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddGoogleProviderRequest(): AddGoogleProviderRequest {
  return { name: "", clientId: "", clientSecret: "", scopes: [], providerOptions: undefined };
}

export const AddGoogleProviderRequest = {
  encode(message: AddGoogleProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.clientId !== "") {
      writer.uint32(18).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(26).string(message.clientSecret);
    }
    for (const v of message.scopes) {
      writer.uint32(34).string(v!);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGoogleProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGoogleProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.clientId = reader.string();
          break;
        case 3:
          message.clientSecret = reader.string();
          break;
        case 4:
          message.scopes.push(reader.string());
          break;
        case 5:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGoogleProviderRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: AddGoogleProviderRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddGoogleProviderRequest>): AddGoogleProviderRequest {
    return AddGoogleProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGoogleProviderRequest>): AddGoogleProviderRequest {
    const message = createBaseAddGoogleProviderRequest();
    message.name = object.name ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseAddGoogleProviderResponse(): AddGoogleProviderResponse {
  return { details: undefined, id: "" };
}

export const AddGoogleProviderResponse = {
  encode(message: AddGoogleProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddGoogleProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddGoogleProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddGoogleProviderResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
    };
  },

  toJSON(message: AddGoogleProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<AddGoogleProviderResponse>): AddGoogleProviderResponse {
    return AddGoogleProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddGoogleProviderResponse>): AddGoogleProviderResponse {
    const message = createBaseAddGoogleProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseUpdateGoogleProviderRequest(): UpdateGoogleProviderRequest {
  return { id: "", name: "", clientId: "", clientSecret: "", scopes: [], providerOptions: undefined };
}

export const UpdateGoogleProviderRequest = {
  encode(message: UpdateGoogleProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.clientId !== "") {
      writer.uint32(26).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(34).string(message.clientSecret);
    }
    for (const v of message.scopes) {
      writer.uint32(42).string(v!);
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGoogleProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGoogleProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.clientId = reader.string();
          break;
        case 4:
          message.clientSecret = reader.string();
          break;
        case 5:
          message.scopes.push(reader.string());
          break;
        case 6:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGoogleProviderRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret) ? String(object.clientSecret) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: UpdateGoogleProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined && (obj.clientSecret = message.clientSecret);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGoogleProviderRequest>): UpdateGoogleProviderRequest {
    return UpdateGoogleProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateGoogleProviderRequest>): UpdateGoogleProviderRequest {
    const message = createBaseUpdateGoogleProviderRequest();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseUpdateGoogleProviderResponse(): UpdateGoogleProviderResponse {
  return { details: undefined };
}

export const UpdateGoogleProviderResponse = {
  encode(message: UpdateGoogleProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateGoogleProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateGoogleProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateGoogleProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateGoogleProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateGoogleProviderResponse>): UpdateGoogleProviderResponse {
    return UpdateGoogleProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateGoogleProviderResponse>): UpdateGoogleProviderResponse {
    const message = createBaseUpdateGoogleProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddLDAPProviderRequest(): AddLDAPProviderRequest {
  return {
    name: "",
    host: "",
    port: "",
    tls: false,
    baseDn: "",
    userObjectClass: "",
    userUniqueAttribute: "",
    admin: "",
    password: "",
    attributes: undefined,
    providerOptions: undefined,
  };
}

export const AddLDAPProviderRequest = {
  encode(message: AddLDAPProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.host !== "") {
      writer.uint32(18).string(message.host);
    }
    if (message.port !== "") {
      writer.uint32(26).string(message.port);
    }
    if (message.tls === true) {
      writer.uint32(32).bool(message.tls);
    }
    if (message.baseDn !== "") {
      writer.uint32(42).string(message.baseDn);
    }
    if (message.userObjectClass !== "") {
      writer.uint32(50).string(message.userObjectClass);
    }
    if (message.userUniqueAttribute !== "") {
      writer.uint32(58).string(message.userUniqueAttribute);
    }
    if (message.admin !== "") {
      writer.uint32(66).string(message.admin);
    }
    if (message.password !== "") {
      writer.uint32(74).string(message.password);
    }
    if (message.attributes !== undefined) {
      LDAPAttributes.encode(message.attributes, writer.uint32(82).fork()).ldelim();
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(90).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddLDAPProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddLDAPProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.host = reader.string();
          break;
        case 3:
          message.port = reader.string();
          break;
        case 4:
          message.tls = reader.bool();
          break;
        case 5:
          message.baseDn = reader.string();
          break;
        case 6:
          message.userObjectClass = reader.string();
          break;
        case 7:
          message.userUniqueAttribute = reader.string();
          break;
        case 8:
          message.admin = reader.string();
          break;
        case 9:
          message.password = reader.string();
          break;
        case 10:
          message.attributes = LDAPAttributes.decode(reader, reader.uint32());
          break;
        case 11:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddLDAPProviderRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      host: isSet(object.host) ? String(object.host) : "",
      port: isSet(object.port) ? String(object.port) : "",
      tls: isSet(object.tls) ? Boolean(object.tls) : false,
      baseDn: isSet(object.baseDn) ? String(object.baseDn) : "",
      userObjectClass: isSet(object.userObjectClass) ? String(object.userObjectClass) : "",
      userUniqueAttribute: isSet(object.userUniqueAttribute) ? String(object.userUniqueAttribute) : "",
      admin: isSet(object.admin) ? String(object.admin) : "",
      password: isSet(object.password) ? String(object.password) : "",
      attributes: isSet(object.attributes) ? LDAPAttributes.fromJSON(object.attributes) : undefined,
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: AddLDAPProviderRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.host !== undefined && (obj.host = message.host);
    message.port !== undefined && (obj.port = message.port);
    message.tls !== undefined && (obj.tls = message.tls);
    message.baseDn !== undefined && (obj.baseDn = message.baseDn);
    message.userObjectClass !== undefined && (obj.userObjectClass = message.userObjectClass);
    message.userUniqueAttribute !== undefined && (obj.userUniqueAttribute = message.userUniqueAttribute);
    message.admin !== undefined && (obj.admin = message.admin);
    message.password !== undefined && (obj.password = message.password);
    message.attributes !== undefined &&
      (obj.attributes = message.attributes ? LDAPAttributes.toJSON(message.attributes) : undefined);
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddLDAPProviderRequest>): AddLDAPProviderRequest {
    return AddLDAPProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddLDAPProviderRequest>): AddLDAPProviderRequest {
    const message = createBaseAddLDAPProviderRequest();
    message.name = object.name ?? "";
    message.host = object.host ?? "";
    message.port = object.port ?? "";
    message.tls = object.tls ?? false;
    message.baseDn = object.baseDn ?? "";
    message.userObjectClass = object.userObjectClass ?? "";
    message.userUniqueAttribute = object.userUniqueAttribute ?? "";
    message.admin = object.admin ?? "";
    message.password = object.password ?? "";
    message.attributes = (object.attributes !== undefined && object.attributes !== null)
      ? LDAPAttributes.fromPartial(object.attributes)
      : undefined;
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseAddLDAPProviderResponse(): AddLDAPProviderResponse {
  return { details: undefined, id: "" };
}

export const AddLDAPProviderResponse = {
  encode(message: AddLDAPProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddLDAPProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddLDAPProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddLDAPProviderResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
    };
  },

  toJSON(message: AddLDAPProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<AddLDAPProviderResponse>): AddLDAPProviderResponse {
    return AddLDAPProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddLDAPProviderResponse>): AddLDAPProviderResponse {
    const message = createBaseAddLDAPProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseUpdateLDAPProviderRequest(): UpdateLDAPProviderRequest {
  return {
    id: "",
    name: "",
    host: "",
    port: "",
    tls: false,
    baseDn: "",
    userObjectClass: "",
    userUniqueAttribute: "",
    admin: "",
    password: "",
    attributes: undefined,
    providerOptions: undefined,
  };
}

export const UpdateLDAPProviderRequest = {
  encode(message: UpdateLDAPProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.host !== "") {
      writer.uint32(26).string(message.host);
    }
    if (message.port !== "") {
      writer.uint32(34).string(message.port);
    }
    if (message.tls === true) {
      writer.uint32(40).bool(message.tls);
    }
    if (message.baseDn !== "") {
      writer.uint32(50).string(message.baseDn);
    }
    if (message.userObjectClass !== "") {
      writer.uint32(58).string(message.userObjectClass);
    }
    if (message.userUniqueAttribute !== "") {
      writer.uint32(66).string(message.userUniqueAttribute);
    }
    if (message.admin !== "") {
      writer.uint32(74).string(message.admin);
    }
    if (message.password !== "") {
      writer.uint32(82).string(message.password);
    }
    if (message.attributes !== undefined) {
      LDAPAttributes.encode(message.attributes, writer.uint32(90).fork()).ldelim();
    }
    if (message.providerOptions !== undefined) {
      Options.encode(message.providerOptions, writer.uint32(98).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateLDAPProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateLDAPProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.host = reader.string();
          break;
        case 4:
          message.port = reader.string();
          break;
        case 5:
          message.tls = reader.bool();
          break;
        case 6:
          message.baseDn = reader.string();
          break;
        case 7:
          message.userObjectClass = reader.string();
          break;
        case 8:
          message.userUniqueAttribute = reader.string();
          break;
        case 9:
          message.admin = reader.string();
          break;
        case 10:
          message.password = reader.string();
          break;
        case 11:
          message.attributes = LDAPAttributes.decode(reader, reader.uint32());
          break;
        case 12:
          message.providerOptions = Options.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateLDAPProviderRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      host: isSet(object.host) ? String(object.host) : "",
      port: isSet(object.port) ? String(object.port) : "",
      tls: isSet(object.tls) ? Boolean(object.tls) : false,
      baseDn: isSet(object.baseDn) ? String(object.baseDn) : "",
      userObjectClass: isSet(object.userObjectClass) ? String(object.userObjectClass) : "",
      userUniqueAttribute: isSet(object.userUniqueAttribute) ? String(object.userUniqueAttribute) : "",
      admin: isSet(object.admin) ? String(object.admin) : "",
      password: isSet(object.password) ? String(object.password) : "",
      attributes: isSet(object.attributes) ? LDAPAttributes.fromJSON(object.attributes) : undefined,
      providerOptions: isSet(object.providerOptions) ? Options.fromJSON(object.providerOptions) : undefined,
    };
  },

  toJSON(message: UpdateLDAPProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.host !== undefined && (obj.host = message.host);
    message.port !== undefined && (obj.port = message.port);
    message.tls !== undefined && (obj.tls = message.tls);
    message.baseDn !== undefined && (obj.baseDn = message.baseDn);
    message.userObjectClass !== undefined && (obj.userObjectClass = message.userObjectClass);
    message.userUniqueAttribute !== undefined && (obj.userUniqueAttribute = message.userUniqueAttribute);
    message.admin !== undefined && (obj.admin = message.admin);
    message.password !== undefined && (obj.password = message.password);
    message.attributes !== undefined &&
      (obj.attributes = message.attributes ? LDAPAttributes.toJSON(message.attributes) : undefined);
    message.providerOptions !== undefined &&
      (obj.providerOptions = message.providerOptions ? Options.toJSON(message.providerOptions) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateLDAPProviderRequest>): UpdateLDAPProviderRequest {
    return UpdateLDAPProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateLDAPProviderRequest>): UpdateLDAPProviderRequest {
    const message = createBaseUpdateLDAPProviderRequest();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.host = object.host ?? "";
    message.port = object.port ?? "";
    message.tls = object.tls ?? false;
    message.baseDn = object.baseDn ?? "";
    message.userObjectClass = object.userObjectClass ?? "";
    message.userUniqueAttribute = object.userUniqueAttribute ?? "";
    message.admin = object.admin ?? "";
    message.password = object.password ?? "";
    message.attributes = (object.attributes !== undefined && object.attributes !== null)
      ? LDAPAttributes.fromPartial(object.attributes)
      : undefined;
    message.providerOptions = (object.providerOptions !== undefined && object.providerOptions !== null)
      ? Options.fromPartial(object.providerOptions)
      : undefined;
    return message;
  },
};

function createBaseUpdateLDAPProviderResponse(): UpdateLDAPProviderResponse {
  return { details: undefined };
}

export const UpdateLDAPProviderResponse = {
  encode(message: UpdateLDAPProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateLDAPProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateLDAPProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateLDAPProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateLDAPProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateLDAPProviderResponse>): UpdateLDAPProviderResponse {
    return UpdateLDAPProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateLDAPProviderResponse>): UpdateLDAPProviderResponse {
    const message = createBaseUpdateLDAPProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseDeleteProviderRequest(): DeleteProviderRequest {
  return { id: "" };
}

export const DeleteProviderRequest = {
  encode(message: DeleteProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeleteProviderRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: DeleteProviderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<DeleteProviderRequest>): DeleteProviderRequest {
    return DeleteProviderRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteProviderRequest>): DeleteProviderRequest {
    const message = createBaseDeleteProviderRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseDeleteProviderResponse(): DeleteProviderResponse {
  return { details: undefined };
}

export const DeleteProviderResponse = {
  encode(message: DeleteProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeleteProviderResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeleteProviderResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeleteProviderResponse>): DeleteProviderResponse {
    return DeleteProviderResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteProviderResponse>): DeleteProviderResponse {
    const message = createBaseDeleteProviderResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetOrgIAMPolicyRequest(): GetOrgIAMPolicyRequest {
  return {};
}

export const GetOrgIAMPolicyRequest = {
  encode(_: GetOrgIAMPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOrgIAMPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOrgIAMPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetOrgIAMPolicyRequest {
    return {};
  },

  toJSON(_: GetOrgIAMPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetOrgIAMPolicyRequest>): GetOrgIAMPolicyRequest {
    return GetOrgIAMPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetOrgIAMPolicyRequest>): GetOrgIAMPolicyRequest {
    const message = createBaseGetOrgIAMPolicyRequest();
    return message;
  },
};

function createBaseGetOrgIAMPolicyResponse(): GetOrgIAMPolicyResponse {
  return { policy: undefined };
}

export const GetOrgIAMPolicyResponse = {
  encode(message: GetOrgIAMPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      OrgIAMPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOrgIAMPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOrgIAMPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = OrgIAMPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetOrgIAMPolicyResponse {
    return { policy: isSet(object.policy) ? OrgIAMPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetOrgIAMPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? OrgIAMPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetOrgIAMPolicyResponse>): GetOrgIAMPolicyResponse {
    return GetOrgIAMPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetOrgIAMPolicyResponse>): GetOrgIAMPolicyResponse {
    const message = createBaseGetOrgIAMPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? OrgIAMPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdateOrgIAMPolicyRequest(): UpdateOrgIAMPolicyRequest {
  return { userLoginMustBeDomain: false };
}

export const UpdateOrgIAMPolicyRequest = {
  encode(message: UpdateOrgIAMPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(8).bool(message.userLoginMustBeDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateOrgIAMPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateOrgIAMPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userLoginMustBeDomain = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateOrgIAMPolicyRequest {
    return {
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
    };
  },

  toJSON(message: UpdateOrgIAMPolicyRequest): unknown {
    const obj: any = {};
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    return obj;
  },

  create(base?: DeepPartial<UpdateOrgIAMPolicyRequest>): UpdateOrgIAMPolicyRequest {
    return UpdateOrgIAMPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateOrgIAMPolicyRequest>): UpdateOrgIAMPolicyRequest {
    const message = createBaseUpdateOrgIAMPolicyRequest();
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    return message;
  },
};

function createBaseUpdateOrgIAMPolicyResponse(): UpdateOrgIAMPolicyResponse {
  return { details: undefined };
}

export const UpdateOrgIAMPolicyResponse = {
  encode(message: UpdateOrgIAMPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateOrgIAMPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateOrgIAMPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateOrgIAMPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateOrgIAMPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateOrgIAMPolicyResponse>): UpdateOrgIAMPolicyResponse {
    return UpdateOrgIAMPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateOrgIAMPolicyResponse>): UpdateOrgIAMPolicyResponse {
    const message = createBaseUpdateOrgIAMPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetCustomOrgIAMPolicyRequest(): GetCustomOrgIAMPolicyRequest {
  return { orgId: "" };
}

export const GetCustomOrgIAMPolicyRequest = {
  encode(message: GetCustomOrgIAMPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomOrgIAMPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomOrgIAMPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomOrgIAMPolicyRequest {
    return { orgId: isSet(object.orgId) ? String(object.orgId) : "" };
  },

  toJSON(message: GetCustomOrgIAMPolicyRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    return obj;
  },

  create(base?: DeepPartial<GetCustomOrgIAMPolicyRequest>): GetCustomOrgIAMPolicyRequest {
    return GetCustomOrgIAMPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomOrgIAMPolicyRequest>): GetCustomOrgIAMPolicyRequest {
    const message = createBaseGetCustomOrgIAMPolicyRequest();
    message.orgId = object.orgId ?? "";
    return message;
  },
};

function createBaseGetCustomOrgIAMPolicyResponse(): GetCustomOrgIAMPolicyResponse {
  return { policy: undefined, isDefault: false };
}

export const GetCustomOrgIAMPolicyResponse = {
  encode(message: GetCustomOrgIAMPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      OrgIAMPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    if (message.isDefault === true) {
      writer.uint32(16).bool(message.isDefault);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomOrgIAMPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomOrgIAMPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = OrgIAMPolicy.decode(reader, reader.uint32());
          break;
        case 2:
          message.isDefault = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomOrgIAMPolicyResponse {
    return {
      policy: isSet(object.policy) ? OrgIAMPolicy.fromJSON(object.policy) : undefined,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
    };
  },

  toJSON(message: GetCustomOrgIAMPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? OrgIAMPolicy.toJSON(message.policy) : undefined);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    return obj;
  },

  create(base?: DeepPartial<GetCustomOrgIAMPolicyResponse>): GetCustomOrgIAMPolicyResponse {
    return GetCustomOrgIAMPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomOrgIAMPolicyResponse>): GetCustomOrgIAMPolicyResponse {
    const message = createBaseGetCustomOrgIAMPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? OrgIAMPolicy.fromPartial(object.policy)
      : undefined;
    message.isDefault = object.isDefault ?? false;
    return message;
  },
};

function createBaseAddCustomOrgIAMPolicyRequest(): AddCustomOrgIAMPolicyRequest {
  return { orgId: "", userLoginMustBeDomain: false };
}

export const AddCustomOrgIAMPolicyRequest = {
  encode(message: AddCustomOrgIAMPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(16).bool(message.userLoginMustBeDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddCustomOrgIAMPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddCustomOrgIAMPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        case 2:
          message.userLoginMustBeDomain = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddCustomOrgIAMPolicyRequest {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
    };
  },

  toJSON(message: AddCustomOrgIAMPolicyRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    return obj;
  },

  create(base?: DeepPartial<AddCustomOrgIAMPolicyRequest>): AddCustomOrgIAMPolicyRequest {
    return AddCustomOrgIAMPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddCustomOrgIAMPolicyRequest>): AddCustomOrgIAMPolicyRequest {
    const message = createBaseAddCustomOrgIAMPolicyRequest();
    message.orgId = object.orgId ?? "";
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    return message;
  },
};

function createBaseAddCustomOrgIAMPolicyResponse(): AddCustomOrgIAMPolicyResponse {
  return { details: undefined };
}

export const AddCustomOrgIAMPolicyResponse = {
  encode(message: AddCustomOrgIAMPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddCustomOrgIAMPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddCustomOrgIAMPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddCustomOrgIAMPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddCustomOrgIAMPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddCustomOrgIAMPolicyResponse>): AddCustomOrgIAMPolicyResponse {
    return AddCustomOrgIAMPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddCustomOrgIAMPolicyResponse>): AddCustomOrgIAMPolicyResponse {
    const message = createBaseAddCustomOrgIAMPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateCustomOrgIAMPolicyRequest(): UpdateCustomOrgIAMPolicyRequest {
  return { orgId: "", userLoginMustBeDomain: false };
}

export const UpdateCustomOrgIAMPolicyRequest = {
  encode(message: UpdateCustomOrgIAMPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(16).bool(message.userLoginMustBeDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateCustomOrgIAMPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateCustomOrgIAMPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        case 2:
          message.userLoginMustBeDomain = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateCustomOrgIAMPolicyRequest {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
    };
  },

  toJSON(message: UpdateCustomOrgIAMPolicyRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    return obj;
  },

  create(base?: DeepPartial<UpdateCustomOrgIAMPolicyRequest>): UpdateCustomOrgIAMPolicyRequest {
    return UpdateCustomOrgIAMPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateCustomOrgIAMPolicyRequest>): UpdateCustomOrgIAMPolicyRequest {
    const message = createBaseUpdateCustomOrgIAMPolicyRequest();
    message.orgId = object.orgId ?? "";
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    return message;
  },
};

function createBaseUpdateCustomOrgIAMPolicyResponse(): UpdateCustomOrgIAMPolicyResponse {
  return { details: undefined };
}

export const UpdateCustomOrgIAMPolicyResponse = {
  encode(message: UpdateCustomOrgIAMPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateCustomOrgIAMPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateCustomOrgIAMPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateCustomOrgIAMPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateCustomOrgIAMPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateCustomOrgIAMPolicyResponse>): UpdateCustomOrgIAMPolicyResponse {
    return UpdateCustomOrgIAMPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateCustomOrgIAMPolicyResponse>): UpdateCustomOrgIAMPolicyResponse {
    const message = createBaseUpdateCustomOrgIAMPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomOrgIAMPolicyToDefaultRequest(): ResetCustomOrgIAMPolicyToDefaultRequest {
  return { orgId: "" };
}

export const ResetCustomOrgIAMPolicyToDefaultRequest = {
  encode(message: ResetCustomOrgIAMPolicyToDefaultRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomOrgIAMPolicyToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomOrgIAMPolicyToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomOrgIAMPolicyToDefaultRequest {
    return { orgId: isSet(object.orgId) ? String(object.orgId) : "" };
  },

  toJSON(message: ResetCustomOrgIAMPolicyToDefaultRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    return obj;
  },

  create(base?: DeepPartial<ResetCustomOrgIAMPolicyToDefaultRequest>): ResetCustomOrgIAMPolicyToDefaultRequest {
    return ResetCustomOrgIAMPolicyToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetCustomOrgIAMPolicyToDefaultRequest>): ResetCustomOrgIAMPolicyToDefaultRequest {
    const message = createBaseResetCustomOrgIAMPolicyToDefaultRequest();
    message.orgId = object.orgId ?? "";
    return message;
  },
};

function createBaseResetCustomOrgIAMPolicyToDefaultResponse(): ResetCustomOrgIAMPolicyToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomOrgIAMPolicyToDefaultResponse = {
  encode(message: ResetCustomOrgIAMPolicyToDefaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomOrgIAMPolicyToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomOrgIAMPolicyToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomOrgIAMPolicyToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomOrgIAMPolicyToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResetCustomOrgIAMPolicyToDefaultResponse>): ResetCustomOrgIAMPolicyToDefaultResponse {
    return ResetCustomOrgIAMPolicyToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetCustomOrgIAMPolicyToDefaultResponse>): ResetCustomOrgIAMPolicyToDefaultResponse {
    const message = createBaseResetCustomOrgIAMPolicyToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDomainPolicyRequest(): GetDomainPolicyRequest {
  return {};
}

export const GetDomainPolicyRequest = {
  encode(_: GetDomainPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDomainPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDomainPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetDomainPolicyRequest {
    return {};
  },

  toJSON(_: GetDomainPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetDomainPolicyRequest>): GetDomainPolicyRequest {
    return GetDomainPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetDomainPolicyRequest>): GetDomainPolicyRequest {
    const message = createBaseGetDomainPolicyRequest();
    return message;
  },
};

function createBaseGetDomainPolicyResponse(): GetDomainPolicyResponse {
  return { policy: undefined };
}

export const GetDomainPolicyResponse = {
  encode(message: GetDomainPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      DomainPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDomainPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDomainPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = DomainPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDomainPolicyResponse {
    return { policy: isSet(object.policy) ? DomainPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetDomainPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? DomainPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDomainPolicyResponse>): GetDomainPolicyResponse {
    return GetDomainPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDomainPolicyResponse>): GetDomainPolicyResponse {
    const message = createBaseGetDomainPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? DomainPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdateDomainPolicyRequest(): UpdateDomainPolicyRequest {
  return { userLoginMustBeDomain: false, validateOrgDomains: false, smtpSenderAddressMatchesInstanceDomain: false };
}

export const UpdateDomainPolicyRequest = {
  encode(message: UpdateDomainPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(8).bool(message.userLoginMustBeDomain);
    }
    if (message.validateOrgDomains === true) {
      writer.uint32(16).bool(message.validateOrgDomains);
    }
    if (message.smtpSenderAddressMatchesInstanceDomain === true) {
      writer.uint32(24).bool(message.smtpSenderAddressMatchesInstanceDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateDomainPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateDomainPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userLoginMustBeDomain = reader.bool();
          break;
        case 2:
          message.validateOrgDomains = reader.bool();
          break;
        case 3:
          message.smtpSenderAddressMatchesInstanceDomain = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateDomainPolicyRequest {
    return {
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
      validateOrgDomains: isSet(object.validateOrgDomains) ? Boolean(object.validateOrgDomains) : false,
      smtpSenderAddressMatchesInstanceDomain: isSet(object.smtpSenderAddressMatchesInstanceDomain)
        ? Boolean(object.smtpSenderAddressMatchesInstanceDomain)
        : false,
    };
  },

  toJSON(message: UpdateDomainPolicyRequest): unknown {
    const obj: any = {};
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    message.validateOrgDomains !== undefined && (obj.validateOrgDomains = message.validateOrgDomains);
    message.smtpSenderAddressMatchesInstanceDomain !== undefined &&
      (obj.smtpSenderAddressMatchesInstanceDomain = message.smtpSenderAddressMatchesInstanceDomain);
    return obj;
  },

  create(base?: DeepPartial<UpdateDomainPolicyRequest>): UpdateDomainPolicyRequest {
    return UpdateDomainPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateDomainPolicyRequest>): UpdateDomainPolicyRequest {
    const message = createBaseUpdateDomainPolicyRequest();
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    message.validateOrgDomains = object.validateOrgDomains ?? false;
    message.smtpSenderAddressMatchesInstanceDomain = object.smtpSenderAddressMatchesInstanceDomain ?? false;
    return message;
  },
};

function createBaseUpdateDomainPolicyResponse(): UpdateDomainPolicyResponse {
  return { details: undefined };
}

export const UpdateDomainPolicyResponse = {
  encode(message: UpdateDomainPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateDomainPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateDomainPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateDomainPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateDomainPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateDomainPolicyResponse>): UpdateDomainPolicyResponse {
    return UpdateDomainPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateDomainPolicyResponse>): UpdateDomainPolicyResponse {
    const message = createBaseUpdateDomainPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetCustomDomainPolicyRequest(): GetCustomDomainPolicyRequest {
  return { orgId: "" };
}

export const GetCustomDomainPolicyRequest = {
  encode(message: GetCustomDomainPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomDomainPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomDomainPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomDomainPolicyRequest {
    return { orgId: isSet(object.orgId) ? String(object.orgId) : "" };
  },

  toJSON(message: GetCustomDomainPolicyRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    return obj;
  },

  create(base?: DeepPartial<GetCustomDomainPolicyRequest>): GetCustomDomainPolicyRequest {
    return GetCustomDomainPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomDomainPolicyRequest>): GetCustomDomainPolicyRequest {
    const message = createBaseGetCustomDomainPolicyRequest();
    message.orgId = object.orgId ?? "";
    return message;
  },
};

function createBaseGetCustomDomainPolicyResponse(): GetCustomDomainPolicyResponse {
  return { policy: undefined, isDefault: false };
}

export const GetCustomDomainPolicyResponse = {
  encode(message: GetCustomDomainPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      DomainPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    if (message.isDefault === true) {
      writer.uint32(16).bool(message.isDefault);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomDomainPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomDomainPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = DomainPolicy.decode(reader, reader.uint32());
          break;
        case 2:
          message.isDefault = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomDomainPolicyResponse {
    return {
      policy: isSet(object.policy) ? DomainPolicy.fromJSON(object.policy) : undefined,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
    };
  },

  toJSON(message: GetCustomDomainPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? DomainPolicy.toJSON(message.policy) : undefined);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    return obj;
  },

  create(base?: DeepPartial<GetCustomDomainPolicyResponse>): GetCustomDomainPolicyResponse {
    return GetCustomDomainPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomDomainPolicyResponse>): GetCustomDomainPolicyResponse {
    const message = createBaseGetCustomDomainPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? DomainPolicy.fromPartial(object.policy)
      : undefined;
    message.isDefault = object.isDefault ?? false;
    return message;
  },
};

function createBaseAddCustomDomainPolicyRequest(): AddCustomDomainPolicyRequest {
  return {
    orgId: "",
    userLoginMustBeDomain: false,
    validateOrgDomains: false,
    smtpSenderAddressMatchesInstanceDomain: false,
  };
}

export const AddCustomDomainPolicyRequest = {
  encode(message: AddCustomDomainPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(16).bool(message.userLoginMustBeDomain);
    }
    if (message.validateOrgDomains === true) {
      writer.uint32(24).bool(message.validateOrgDomains);
    }
    if (message.smtpSenderAddressMatchesInstanceDomain === true) {
      writer.uint32(32).bool(message.smtpSenderAddressMatchesInstanceDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddCustomDomainPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddCustomDomainPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        case 2:
          message.userLoginMustBeDomain = reader.bool();
          break;
        case 3:
          message.validateOrgDomains = reader.bool();
          break;
        case 4:
          message.smtpSenderAddressMatchesInstanceDomain = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddCustomDomainPolicyRequest {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
      validateOrgDomains: isSet(object.validateOrgDomains) ? Boolean(object.validateOrgDomains) : false,
      smtpSenderAddressMatchesInstanceDomain: isSet(object.smtpSenderAddressMatchesInstanceDomain)
        ? Boolean(object.smtpSenderAddressMatchesInstanceDomain)
        : false,
    };
  },

  toJSON(message: AddCustomDomainPolicyRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    message.validateOrgDomains !== undefined && (obj.validateOrgDomains = message.validateOrgDomains);
    message.smtpSenderAddressMatchesInstanceDomain !== undefined &&
      (obj.smtpSenderAddressMatchesInstanceDomain = message.smtpSenderAddressMatchesInstanceDomain);
    return obj;
  },

  create(base?: DeepPartial<AddCustomDomainPolicyRequest>): AddCustomDomainPolicyRequest {
    return AddCustomDomainPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddCustomDomainPolicyRequest>): AddCustomDomainPolicyRequest {
    const message = createBaseAddCustomDomainPolicyRequest();
    message.orgId = object.orgId ?? "";
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    message.validateOrgDomains = object.validateOrgDomains ?? false;
    message.smtpSenderAddressMatchesInstanceDomain = object.smtpSenderAddressMatchesInstanceDomain ?? false;
    return message;
  },
};

function createBaseAddCustomDomainPolicyResponse(): AddCustomDomainPolicyResponse {
  return { details: undefined };
}

export const AddCustomDomainPolicyResponse = {
  encode(message: AddCustomDomainPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddCustomDomainPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddCustomDomainPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddCustomDomainPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddCustomDomainPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddCustomDomainPolicyResponse>): AddCustomDomainPolicyResponse {
    return AddCustomDomainPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddCustomDomainPolicyResponse>): AddCustomDomainPolicyResponse {
    const message = createBaseAddCustomDomainPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateCustomDomainPolicyRequest(): UpdateCustomDomainPolicyRequest {
  return {
    orgId: "",
    userLoginMustBeDomain: false,
    validateOrgDomains: false,
    smtpSenderAddressMatchesInstanceDomain: false,
  };
}

export const UpdateCustomDomainPolicyRequest = {
  encode(message: UpdateCustomDomainPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(16).bool(message.userLoginMustBeDomain);
    }
    if (message.validateOrgDomains === true) {
      writer.uint32(24).bool(message.validateOrgDomains);
    }
    if (message.smtpSenderAddressMatchesInstanceDomain === true) {
      writer.uint32(32).bool(message.smtpSenderAddressMatchesInstanceDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateCustomDomainPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateCustomDomainPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        case 2:
          message.userLoginMustBeDomain = reader.bool();
          break;
        case 3:
          message.validateOrgDomains = reader.bool();
          break;
        case 4:
          message.smtpSenderAddressMatchesInstanceDomain = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateCustomDomainPolicyRequest {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
      validateOrgDomains: isSet(object.validateOrgDomains) ? Boolean(object.validateOrgDomains) : false,
      smtpSenderAddressMatchesInstanceDomain: isSet(object.smtpSenderAddressMatchesInstanceDomain)
        ? Boolean(object.smtpSenderAddressMatchesInstanceDomain)
        : false,
    };
  },

  toJSON(message: UpdateCustomDomainPolicyRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    message.validateOrgDomains !== undefined && (obj.validateOrgDomains = message.validateOrgDomains);
    message.smtpSenderAddressMatchesInstanceDomain !== undefined &&
      (obj.smtpSenderAddressMatchesInstanceDomain = message.smtpSenderAddressMatchesInstanceDomain);
    return obj;
  },

  create(base?: DeepPartial<UpdateCustomDomainPolicyRequest>): UpdateCustomDomainPolicyRequest {
    return UpdateCustomDomainPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateCustomDomainPolicyRequest>): UpdateCustomDomainPolicyRequest {
    const message = createBaseUpdateCustomDomainPolicyRequest();
    message.orgId = object.orgId ?? "";
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    message.validateOrgDomains = object.validateOrgDomains ?? false;
    message.smtpSenderAddressMatchesInstanceDomain = object.smtpSenderAddressMatchesInstanceDomain ?? false;
    return message;
  },
};

function createBaseUpdateCustomDomainPolicyResponse(): UpdateCustomDomainPolicyResponse {
  return { details: undefined };
}

export const UpdateCustomDomainPolicyResponse = {
  encode(message: UpdateCustomDomainPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateCustomDomainPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateCustomDomainPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateCustomDomainPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateCustomDomainPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateCustomDomainPolicyResponse>): UpdateCustomDomainPolicyResponse {
    return UpdateCustomDomainPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateCustomDomainPolicyResponse>): UpdateCustomDomainPolicyResponse {
    const message = createBaseUpdateCustomDomainPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomDomainPolicyToDefaultRequest(): ResetCustomDomainPolicyToDefaultRequest {
  return { orgId: "" };
}

export const ResetCustomDomainPolicyToDefaultRequest = {
  encode(message: ResetCustomDomainPolicyToDefaultRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomDomainPolicyToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomDomainPolicyToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomDomainPolicyToDefaultRequest {
    return { orgId: isSet(object.orgId) ? String(object.orgId) : "" };
  },

  toJSON(message: ResetCustomDomainPolicyToDefaultRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    return obj;
  },

  create(base?: DeepPartial<ResetCustomDomainPolicyToDefaultRequest>): ResetCustomDomainPolicyToDefaultRequest {
    return ResetCustomDomainPolicyToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetCustomDomainPolicyToDefaultRequest>): ResetCustomDomainPolicyToDefaultRequest {
    const message = createBaseResetCustomDomainPolicyToDefaultRequest();
    message.orgId = object.orgId ?? "";
    return message;
  },
};

function createBaseResetCustomDomainPolicyToDefaultResponse(): ResetCustomDomainPolicyToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomDomainPolicyToDefaultResponse = {
  encode(message: ResetCustomDomainPolicyToDefaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomDomainPolicyToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomDomainPolicyToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomDomainPolicyToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomDomainPolicyToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResetCustomDomainPolicyToDefaultResponse>): ResetCustomDomainPolicyToDefaultResponse {
    return ResetCustomDomainPolicyToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetCustomDomainPolicyToDefaultResponse>): ResetCustomDomainPolicyToDefaultResponse {
    const message = createBaseResetCustomDomainPolicyToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetLabelPolicyRequest(): GetLabelPolicyRequest {
  return {};
}

export const GetLabelPolicyRequest = {
  encode(_: GetLabelPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLabelPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLabelPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetLabelPolicyRequest {
    return {};
  },

  toJSON(_: GetLabelPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetLabelPolicyRequest>): GetLabelPolicyRequest {
    return GetLabelPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetLabelPolicyRequest>): GetLabelPolicyRequest {
    const message = createBaseGetLabelPolicyRequest();
    return message;
  },
};

function createBaseGetLabelPolicyResponse(): GetLabelPolicyResponse {
  return { policy: undefined };
}

export const GetLabelPolicyResponse = {
  encode(message: GetLabelPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      LabelPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLabelPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLabelPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = LabelPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetLabelPolicyResponse {
    return { policy: isSet(object.policy) ? LabelPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetLabelPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? LabelPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetLabelPolicyResponse>): GetLabelPolicyResponse {
    return GetLabelPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetLabelPolicyResponse>): GetLabelPolicyResponse {
    const message = createBaseGetLabelPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? LabelPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseGetPreviewLabelPolicyRequest(): GetPreviewLabelPolicyRequest {
  return {};
}

export const GetPreviewLabelPolicyRequest = {
  encode(_: GetPreviewLabelPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPreviewLabelPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPreviewLabelPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetPreviewLabelPolicyRequest {
    return {};
  },

  toJSON(_: GetPreviewLabelPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetPreviewLabelPolicyRequest>): GetPreviewLabelPolicyRequest {
    return GetPreviewLabelPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetPreviewLabelPolicyRequest>): GetPreviewLabelPolicyRequest {
    const message = createBaseGetPreviewLabelPolicyRequest();
    return message;
  },
};

function createBaseGetPreviewLabelPolicyResponse(): GetPreviewLabelPolicyResponse {
  return { policy: undefined };
}

export const GetPreviewLabelPolicyResponse = {
  encode(message: GetPreviewLabelPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      LabelPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPreviewLabelPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPreviewLabelPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = LabelPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetPreviewLabelPolicyResponse {
    return { policy: isSet(object.policy) ? LabelPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetPreviewLabelPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? LabelPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetPreviewLabelPolicyResponse>): GetPreviewLabelPolicyResponse {
    return GetPreviewLabelPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetPreviewLabelPolicyResponse>): GetPreviewLabelPolicyResponse {
    const message = createBaseGetPreviewLabelPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? LabelPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdateLabelPolicyRequest(): UpdateLabelPolicyRequest {
  return {
    primaryColor: "",
    hideLoginNameSuffix: false,
    warnColor: "",
    backgroundColor: "",
    fontColor: "",
    primaryColorDark: "",
    backgroundColorDark: "",
    warnColorDark: "",
    fontColorDark: "",
    disableWatermark: false,
  };
}

export const UpdateLabelPolicyRequest = {
  encode(message: UpdateLabelPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.primaryColor !== "") {
      writer.uint32(10).string(message.primaryColor);
    }
    if (message.hideLoginNameSuffix === true) {
      writer.uint32(24).bool(message.hideLoginNameSuffix);
    }
    if (message.warnColor !== "") {
      writer.uint32(34).string(message.warnColor);
    }
    if (message.backgroundColor !== "") {
      writer.uint32(42).string(message.backgroundColor);
    }
    if (message.fontColor !== "") {
      writer.uint32(50).string(message.fontColor);
    }
    if (message.primaryColorDark !== "") {
      writer.uint32(58).string(message.primaryColorDark);
    }
    if (message.backgroundColorDark !== "") {
      writer.uint32(66).string(message.backgroundColorDark);
    }
    if (message.warnColorDark !== "") {
      writer.uint32(74).string(message.warnColorDark);
    }
    if (message.fontColorDark !== "") {
      writer.uint32(82).string(message.fontColorDark);
    }
    if (message.disableWatermark === true) {
      writer.uint32(88).bool(message.disableWatermark);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateLabelPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateLabelPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.primaryColor = reader.string();
          break;
        case 3:
          message.hideLoginNameSuffix = reader.bool();
          break;
        case 4:
          message.warnColor = reader.string();
          break;
        case 5:
          message.backgroundColor = reader.string();
          break;
        case 6:
          message.fontColor = reader.string();
          break;
        case 7:
          message.primaryColorDark = reader.string();
          break;
        case 8:
          message.backgroundColorDark = reader.string();
          break;
        case 9:
          message.warnColorDark = reader.string();
          break;
        case 10:
          message.fontColorDark = reader.string();
          break;
        case 11:
          message.disableWatermark = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateLabelPolicyRequest {
    return {
      primaryColor: isSet(object.primaryColor) ? String(object.primaryColor) : "",
      hideLoginNameSuffix: isSet(object.hideLoginNameSuffix) ? Boolean(object.hideLoginNameSuffix) : false,
      warnColor: isSet(object.warnColor) ? String(object.warnColor) : "",
      backgroundColor: isSet(object.backgroundColor) ? String(object.backgroundColor) : "",
      fontColor: isSet(object.fontColor) ? String(object.fontColor) : "",
      primaryColorDark: isSet(object.primaryColorDark) ? String(object.primaryColorDark) : "",
      backgroundColorDark: isSet(object.backgroundColorDark) ? String(object.backgroundColorDark) : "",
      warnColorDark: isSet(object.warnColorDark) ? String(object.warnColorDark) : "",
      fontColorDark: isSet(object.fontColorDark) ? String(object.fontColorDark) : "",
      disableWatermark: isSet(object.disableWatermark) ? Boolean(object.disableWatermark) : false,
    };
  },

  toJSON(message: UpdateLabelPolicyRequest): unknown {
    const obj: any = {};
    message.primaryColor !== undefined && (obj.primaryColor = message.primaryColor);
    message.hideLoginNameSuffix !== undefined && (obj.hideLoginNameSuffix = message.hideLoginNameSuffix);
    message.warnColor !== undefined && (obj.warnColor = message.warnColor);
    message.backgroundColor !== undefined && (obj.backgroundColor = message.backgroundColor);
    message.fontColor !== undefined && (obj.fontColor = message.fontColor);
    message.primaryColorDark !== undefined && (obj.primaryColorDark = message.primaryColorDark);
    message.backgroundColorDark !== undefined && (obj.backgroundColorDark = message.backgroundColorDark);
    message.warnColorDark !== undefined && (obj.warnColorDark = message.warnColorDark);
    message.fontColorDark !== undefined && (obj.fontColorDark = message.fontColorDark);
    message.disableWatermark !== undefined && (obj.disableWatermark = message.disableWatermark);
    return obj;
  },

  create(base?: DeepPartial<UpdateLabelPolicyRequest>): UpdateLabelPolicyRequest {
    return UpdateLabelPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateLabelPolicyRequest>): UpdateLabelPolicyRequest {
    const message = createBaseUpdateLabelPolicyRequest();
    message.primaryColor = object.primaryColor ?? "";
    message.hideLoginNameSuffix = object.hideLoginNameSuffix ?? false;
    message.warnColor = object.warnColor ?? "";
    message.backgroundColor = object.backgroundColor ?? "";
    message.fontColor = object.fontColor ?? "";
    message.primaryColorDark = object.primaryColorDark ?? "";
    message.backgroundColorDark = object.backgroundColorDark ?? "";
    message.warnColorDark = object.warnColorDark ?? "";
    message.fontColorDark = object.fontColorDark ?? "";
    message.disableWatermark = object.disableWatermark ?? false;
    return message;
  },
};

function createBaseUpdateLabelPolicyResponse(): UpdateLabelPolicyResponse {
  return { details: undefined };
}

export const UpdateLabelPolicyResponse = {
  encode(message: UpdateLabelPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateLabelPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateLabelPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateLabelPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateLabelPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateLabelPolicyResponse>): UpdateLabelPolicyResponse {
    return UpdateLabelPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateLabelPolicyResponse>): UpdateLabelPolicyResponse {
    const message = createBaseUpdateLabelPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseActivateLabelPolicyRequest(): ActivateLabelPolicyRequest {
  return {};
}

export const ActivateLabelPolicyRequest = {
  encode(_: ActivateLabelPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActivateLabelPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActivateLabelPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ActivateLabelPolicyRequest {
    return {};
  },

  toJSON(_: ActivateLabelPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ActivateLabelPolicyRequest>): ActivateLabelPolicyRequest {
    return ActivateLabelPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ActivateLabelPolicyRequest>): ActivateLabelPolicyRequest {
    const message = createBaseActivateLabelPolicyRequest();
    return message;
  },
};

function createBaseActivateLabelPolicyResponse(): ActivateLabelPolicyResponse {
  return { details: undefined };
}

export const ActivateLabelPolicyResponse = {
  encode(message: ActivateLabelPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActivateLabelPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActivateLabelPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ActivateLabelPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ActivateLabelPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ActivateLabelPolicyResponse>): ActivateLabelPolicyResponse {
    return ActivateLabelPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ActivateLabelPolicyResponse>): ActivateLabelPolicyResponse {
    const message = createBaseActivateLabelPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveLabelPolicyLogoRequest(): RemoveLabelPolicyLogoRequest {
  return {};
}

export const RemoveLabelPolicyLogoRequest = {
  encode(_: RemoveLabelPolicyLogoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyLogoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyLogoRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveLabelPolicyLogoRequest {
    return {};
  },

  toJSON(_: RemoveLabelPolicyLogoRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyLogoRequest>): RemoveLabelPolicyLogoRequest {
    return RemoveLabelPolicyLogoRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveLabelPolicyLogoRequest>): RemoveLabelPolicyLogoRequest {
    const message = createBaseRemoveLabelPolicyLogoRequest();
    return message;
  },
};

function createBaseRemoveLabelPolicyLogoResponse(): RemoveLabelPolicyLogoResponse {
  return { details: undefined };
}

export const RemoveLabelPolicyLogoResponse = {
  encode(message: RemoveLabelPolicyLogoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyLogoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyLogoResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveLabelPolicyLogoResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveLabelPolicyLogoResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyLogoResponse>): RemoveLabelPolicyLogoResponse {
    return RemoveLabelPolicyLogoResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveLabelPolicyLogoResponse>): RemoveLabelPolicyLogoResponse {
    const message = createBaseRemoveLabelPolicyLogoResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveLabelPolicyLogoDarkRequest(): RemoveLabelPolicyLogoDarkRequest {
  return {};
}

export const RemoveLabelPolicyLogoDarkRequest = {
  encode(_: RemoveLabelPolicyLogoDarkRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyLogoDarkRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyLogoDarkRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveLabelPolicyLogoDarkRequest {
    return {};
  },

  toJSON(_: RemoveLabelPolicyLogoDarkRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyLogoDarkRequest>): RemoveLabelPolicyLogoDarkRequest {
    return RemoveLabelPolicyLogoDarkRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveLabelPolicyLogoDarkRequest>): RemoveLabelPolicyLogoDarkRequest {
    const message = createBaseRemoveLabelPolicyLogoDarkRequest();
    return message;
  },
};

function createBaseRemoveLabelPolicyLogoDarkResponse(): RemoveLabelPolicyLogoDarkResponse {
  return { details: undefined };
}

export const RemoveLabelPolicyLogoDarkResponse = {
  encode(message: RemoveLabelPolicyLogoDarkResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyLogoDarkResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyLogoDarkResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveLabelPolicyLogoDarkResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveLabelPolicyLogoDarkResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyLogoDarkResponse>): RemoveLabelPolicyLogoDarkResponse {
    return RemoveLabelPolicyLogoDarkResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveLabelPolicyLogoDarkResponse>): RemoveLabelPolicyLogoDarkResponse {
    const message = createBaseRemoveLabelPolicyLogoDarkResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveLabelPolicyIconRequest(): RemoveLabelPolicyIconRequest {
  return {};
}

export const RemoveLabelPolicyIconRequest = {
  encode(_: RemoveLabelPolicyIconRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyIconRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyIconRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveLabelPolicyIconRequest {
    return {};
  },

  toJSON(_: RemoveLabelPolicyIconRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyIconRequest>): RemoveLabelPolicyIconRequest {
    return RemoveLabelPolicyIconRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveLabelPolicyIconRequest>): RemoveLabelPolicyIconRequest {
    const message = createBaseRemoveLabelPolicyIconRequest();
    return message;
  },
};

function createBaseRemoveLabelPolicyIconResponse(): RemoveLabelPolicyIconResponse {
  return { details: undefined };
}

export const RemoveLabelPolicyIconResponse = {
  encode(message: RemoveLabelPolicyIconResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyIconResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyIconResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveLabelPolicyIconResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveLabelPolicyIconResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyIconResponse>): RemoveLabelPolicyIconResponse {
    return RemoveLabelPolicyIconResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveLabelPolicyIconResponse>): RemoveLabelPolicyIconResponse {
    const message = createBaseRemoveLabelPolicyIconResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveLabelPolicyIconDarkRequest(): RemoveLabelPolicyIconDarkRequest {
  return {};
}

export const RemoveLabelPolicyIconDarkRequest = {
  encode(_: RemoveLabelPolicyIconDarkRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyIconDarkRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyIconDarkRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveLabelPolicyIconDarkRequest {
    return {};
  },

  toJSON(_: RemoveLabelPolicyIconDarkRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyIconDarkRequest>): RemoveLabelPolicyIconDarkRequest {
    return RemoveLabelPolicyIconDarkRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveLabelPolicyIconDarkRequest>): RemoveLabelPolicyIconDarkRequest {
    const message = createBaseRemoveLabelPolicyIconDarkRequest();
    return message;
  },
};

function createBaseRemoveLabelPolicyIconDarkResponse(): RemoveLabelPolicyIconDarkResponse {
  return { details: undefined };
}

export const RemoveLabelPolicyIconDarkResponse = {
  encode(message: RemoveLabelPolicyIconDarkResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyIconDarkResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyIconDarkResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveLabelPolicyIconDarkResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveLabelPolicyIconDarkResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyIconDarkResponse>): RemoveLabelPolicyIconDarkResponse {
    return RemoveLabelPolicyIconDarkResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveLabelPolicyIconDarkResponse>): RemoveLabelPolicyIconDarkResponse {
    const message = createBaseRemoveLabelPolicyIconDarkResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveLabelPolicyFontRequest(): RemoveLabelPolicyFontRequest {
  return {};
}

export const RemoveLabelPolicyFontRequest = {
  encode(_: RemoveLabelPolicyFontRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyFontRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyFontRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveLabelPolicyFontRequest {
    return {};
  },

  toJSON(_: RemoveLabelPolicyFontRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyFontRequest>): RemoveLabelPolicyFontRequest {
    return RemoveLabelPolicyFontRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveLabelPolicyFontRequest>): RemoveLabelPolicyFontRequest {
    const message = createBaseRemoveLabelPolicyFontRequest();
    return message;
  },
};

function createBaseRemoveLabelPolicyFontResponse(): RemoveLabelPolicyFontResponse {
  return { details: undefined };
}

export const RemoveLabelPolicyFontResponse = {
  encode(message: RemoveLabelPolicyFontResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveLabelPolicyFontResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveLabelPolicyFontResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveLabelPolicyFontResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveLabelPolicyFontResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveLabelPolicyFontResponse>): RemoveLabelPolicyFontResponse {
    return RemoveLabelPolicyFontResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveLabelPolicyFontResponse>): RemoveLabelPolicyFontResponse {
    const message = createBaseRemoveLabelPolicyFontResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetLoginPolicyRequest(): GetLoginPolicyRequest {
  return {};
}

export const GetLoginPolicyRequest = {
  encode(_: GetLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetLoginPolicyRequest {
    return {};
  },

  toJSON(_: GetLoginPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetLoginPolicyRequest>): GetLoginPolicyRequest {
    return GetLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetLoginPolicyRequest>): GetLoginPolicyRequest {
    const message = createBaseGetLoginPolicyRequest();
    return message;
  },
};

function createBaseGetLoginPolicyResponse(): GetLoginPolicyResponse {
  return { policy: undefined };
}

export const GetLoginPolicyResponse = {
  encode(message: GetLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      LoginPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = LoginPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetLoginPolicyResponse {
    return { policy: isSet(object.policy) ? LoginPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetLoginPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? LoginPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetLoginPolicyResponse>): GetLoginPolicyResponse {
    return GetLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetLoginPolicyResponse>): GetLoginPolicyResponse {
    const message = createBaseGetLoginPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? LoginPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdateLoginPolicyRequest(): UpdateLoginPolicyRequest {
  return {
    allowUsernamePassword: false,
    allowRegister: false,
    allowExternalIdp: false,
    forceMfa: false,
    passwordlessType: 0,
    hidePasswordReset: false,
    ignoreUnknownUsernames: false,
    defaultRedirectUri: "",
    passwordCheckLifetime: undefined,
    externalLoginCheckLifetime: undefined,
    mfaInitSkipLifetime: undefined,
    secondFactorCheckLifetime: undefined,
    multiFactorCheckLifetime: undefined,
    allowDomainDiscovery: false,
    disableLoginWithEmail: false,
    disableLoginWithPhone: false,
  };
}

export const UpdateLoginPolicyRequest = {
  encode(message: UpdateLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.allowUsernamePassword === true) {
      writer.uint32(8).bool(message.allowUsernamePassword);
    }
    if (message.allowRegister === true) {
      writer.uint32(16).bool(message.allowRegister);
    }
    if (message.allowExternalIdp === true) {
      writer.uint32(24).bool(message.allowExternalIdp);
    }
    if (message.forceMfa === true) {
      writer.uint32(32).bool(message.forceMfa);
    }
    if (message.passwordlessType !== 0) {
      writer.uint32(40).int32(message.passwordlessType);
    }
    if (message.hidePasswordReset === true) {
      writer.uint32(48).bool(message.hidePasswordReset);
    }
    if (message.ignoreUnknownUsernames === true) {
      writer.uint32(56).bool(message.ignoreUnknownUsernames);
    }
    if (message.defaultRedirectUri !== "") {
      writer.uint32(66).string(message.defaultRedirectUri);
    }
    if (message.passwordCheckLifetime !== undefined) {
      Duration.encode(message.passwordCheckLifetime, writer.uint32(74).fork()).ldelim();
    }
    if (message.externalLoginCheckLifetime !== undefined) {
      Duration.encode(message.externalLoginCheckLifetime, writer.uint32(82).fork()).ldelim();
    }
    if (message.mfaInitSkipLifetime !== undefined) {
      Duration.encode(message.mfaInitSkipLifetime, writer.uint32(90).fork()).ldelim();
    }
    if (message.secondFactorCheckLifetime !== undefined) {
      Duration.encode(message.secondFactorCheckLifetime, writer.uint32(98).fork()).ldelim();
    }
    if (message.multiFactorCheckLifetime !== undefined) {
      Duration.encode(message.multiFactorCheckLifetime, writer.uint32(106).fork()).ldelim();
    }
    if (message.allowDomainDiscovery === true) {
      writer.uint32(112).bool(message.allowDomainDiscovery);
    }
    if (message.disableLoginWithEmail === true) {
      writer.uint32(120).bool(message.disableLoginWithEmail);
    }
    if (message.disableLoginWithPhone === true) {
      writer.uint32(128).bool(message.disableLoginWithPhone);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.allowUsernamePassword = reader.bool();
          break;
        case 2:
          message.allowRegister = reader.bool();
          break;
        case 3:
          message.allowExternalIdp = reader.bool();
          break;
        case 4:
          message.forceMfa = reader.bool();
          break;
        case 5:
          message.passwordlessType = reader.int32() as any;
          break;
        case 6:
          message.hidePasswordReset = reader.bool();
          break;
        case 7:
          message.ignoreUnknownUsernames = reader.bool();
          break;
        case 8:
          message.defaultRedirectUri = reader.string();
          break;
        case 9:
          message.passwordCheckLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 10:
          message.externalLoginCheckLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 11:
          message.mfaInitSkipLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 12:
          message.secondFactorCheckLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 13:
          message.multiFactorCheckLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 14:
          message.allowDomainDiscovery = reader.bool();
          break;
        case 15:
          message.disableLoginWithEmail = reader.bool();
          break;
        case 16:
          message.disableLoginWithPhone = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateLoginPolicyRequest {
    return {
      allowUsernamePassword: isSet(object.allowUsernamePassword) ? Boolean(object.allowUsernamePassword) : false,
      allowRegister: isSet(object.allowRegister) ? Boolean(object.allowRegister) : false,
      allowExternalIdp: isSet(object.allowExternalIdp) ? Boolean(object.allowExternalIdp) : false,
      forceMfa: isSet(object.forceMfa) ? Boolean(object.forceMfa) : false,
      passwordlessType: isSet(object.passwordlessType) ? passwordlessTypeFromJSON(object.passwordlessType) : 0,
      hidePasswordReset: isSet(object.hidePasswordReset) ? Boolean(object.hidePasswordReset) : false,
      ignoreUnknownUsernames: isSet(object.ignoreUnknownUsernames) ? Boolean(object.ignoreUnknownUsernames) : false,
      defaultRedirectUri: isSet(object.defaultRedirectUri) ? String(object.defaultRedirectUri) : "",
      passwordCheckLifetime: isSet(object.passwordCheckLifetime)
        ? Duration.fromJSON(object.passwordCheckLifetime)
        : undefined,
      externalLoginCheckLifetime: isSet(object.externalLoginCheckLifetime)
        ? Duration.fromJSON(object.externalLoginCheckLifetime)
        : undefined,
      mfaInitSkipLifetime: isSet(object.mfaInitSkipLifetime)
        ? Duration.fromJSON(object.mfaInitSkipLifetime)
        : undefined,
      secondFactorCheckLifetime: isSet(object.secondFactorCheckLifetime)
        ? Duration.fromJSON(object.secondFactorCheckLifetime)
        : undefined,
      multiFactorCheckLifetime: isSet(object.multiFactorCheckLifetime)
        ? Duration.fromJSON(object.multiFactorCheckLifetime)
        : undefined,
      allowDomainDiscovery: isSet(object.allowDomainDiscovery) ? Boolean(object.allowDomainDiscovery) : false,
      disableLoginWithEmail: isSet(object.disableLoginWithEmail) ? Boolean(object.disableLoginWithEmail) : false,
      disableLoginWithPhone: isSet(object.disableLoginWithPhone) ? Boolean(object.disableLoginWithPhone) : false,
    };
  },

  toJSON(message: UpdateLoginPolicyRequest): unknown {
    const obj: any = {};
    message.allowUsernamePassword !== undefined && (obj.allowUsernamePassword = message.allowUsernamePassword);
    message.allowRegister !== undefined && (obj.allowRegister = message.allowRegister);
    message.allowExternalIdp !== undefined && (obj.allowExternalIdp = message.allowExternalIdp);
    message.forceMfa !== undefined && (obj.forceMfa = message.forceMfa);
    message.passwordlessType !== undefined && (obj.passwordlessType = passwordlessTypeToJSON(message.passwordlessType));
    message.hidePasswordReset !== undefined && (obj.hidePasswordReset = message.hidePasswordReset);
    message.ignoreUnknownUsernames !== undefined && (obj.ignoreUnknownUsernames = message.ignoreUnknownUsernames);
    message.defaultRedirectUri !== undefined && (obj.defaultRedirectUri = message.defaultRedirectUri);
    message.passwordCheckLifetime !== undefined && (obj.passwordCheckLifetime = message.passwordCheckLifetime
      ? Duration.toJSON(message.passwordCheckLifetime)
      : undefined);
    message.externalLoginCheckLifetime !== undefined &&
      (obj.externalLoginCheckLifetime = message.externalLoginCheckLifetime
        ? Duration.toJSON(message.externalLoginCheckLifetime)
        : undefined);
    message.mfaInitSkipLifetime !== undefined &&
      (obj.mfaInitSkipLifetime = message.mfaInitSkipLifetime
        ? Duration.toJSON(message.mfaInitSkipLifetime)
        : undefined);
    message.secondFactorCheckLifetime !== undefined &&
      (obj.secondFactorCheckLifetime = message.secondFactorCheckLifetime
        ? Duration.toJSON(message.secondFactorCheckLifetime)
        : undefined);
    message.multiFactorCheckLifetime !== undefined && (obj.multiFactorCheckLifetime = message.multiFactorCheckLifetime
      ? Duration.toJSON(message.multiFactorCheckLifetime)
      : undefined);
    message.allowDomainDiscovery !== undefined && (obj.allowDomainDiscovery = message.allowDomainDiscovery);
    message.disableLoginWithEmail !== undefined && (obj.disableLoginWithEmail = message.disableLoginWithEmail);
    message.disableLoginWithPhone !== undefined && (obj.disableLoginWithPhone = message.disableLoginWithPhone);
    return obj;
  },

  create(base?: DeepPartial<UpdateLoginPolicyRequest>): UpdateLoginPolicyRequest {
    return UpdateLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateLoginPolicyRequest>): UpdateLoginPolicyRequest {
    const message = createBaseUpdateLoginPolicyRequest();
    message.allowUsernamePassword = object.allowUsernamePassword ?? false;
    message.allowRegister = object.allowRegister ?? false;
    message.allowExternalIdp = object.allowExternalIdp ?? false;
    message.forceMfa = object.forceMfa ?? false;
    message.passwordlessType = object.passwordlessType ?? 0;
    message.hidePasswordReset = object.hidePasswordReset ?? false;
    message.ignoreUnknownUsernames = object.ignoreUnknownUsernames ?? false;
    message.defaultRedirectUri = object.defaultRedirectUri ?? "";
    message.passwordCheckLifetime =
      (object.passwordCheckLifetime !== undefined && object.passwordCheckLifetime !== null)
        ? Duration.fromPartial(object.passwordCheckLifetime)
        : undefined;
    message.externalLoginCheckLifetime =
      (object.externalLoginCheckLifetime !== undefined && object.externalLoginCheckLifetime !== null)
        ? Duration.fromPartial(object.externalLoginCheckLifetime)
        : undefined;
    message.mfaInitSkipLifetime = (object.mfaInitSkipLifetime !== undefined && object.mfaInitSkipLifetime !== null)
      ? Duration.fromPartial(object.mfaInitSkipLifetime)
      : undefined;
    message.secondFactorCheckLifetime =
      (object.secondFactorCheckLifetime !== undefined && object.secondFactorCheckLifetime !== null)
        ? Duration.fromPartial(object.secondFactorCheckLifetime)
        : undefined;
    message.multiFactorCheckLifetime =
      (object.multiFactorCheckLifetime !== undefined && object.multiFactorCheckLifetime !== null)
        ? Duration.fromPartial(object.multiFactorCheckLifetime)
        : undefined;
    message.allowDomainDiscovery = object.allowDomainDiscovery ?? false;
    message.disableLoginWithEmail = object.disableLoginWithEmail ?? false;
    message.disableLoginWithPhone = object.disableLoginWithPhone ?? false;
    return message;
  },
};

function createBaseUpdateLoginPolicyResponse(): UpdateLoginPolicyResponse {
  return { details: undefined };
}

export const UpdateLoginPolicyResponse = {
  encode(message: UpdateLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateLoginPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateLoginPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateLoginPolicyResponse>): UpdateLoginPolicyResponse {
    return UpdateLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateLoginPolicyResponse>): UpdateLoginPolicyResponse {
    const message = createBaseUpdateLoginPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListLoginPolicyIDPsRequest(): ListLoginPolicyIDPsRequest {
  return { query: undefined };
}

export const ListLoginPolicyIDPsRequest = {
  encode(message: ListLoginPolicyIDPsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListLoginPolicyIDPsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListLoginPolicyIDPsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListLoginPolicyIDPsRequest {
    return { query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined };
  },

  toJSON(message: ListLoginPolicyIDPsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ListLoginPolicyIDPsRequest>): ListLoginPolicyIDPsRequest {
    return ListLoginPolicyIDPsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListLoginPolicyIDPsRequest>): ListLoginPolicyIDPsRequest {
    const message = createBaseListLoginPolicyIDPsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseListLoginPolicyIDPsResponse(): ListLoginPolicyIDPsResponse {
  return { details: undefined, result: [] };
}

export const ListLoginPolicyIDPsResponse = {
  encode(message: ListLoginPolicyIDPsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      IDPLoginPolicyLink.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListLoginPolicyIDPsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListLoginPolicyIDPsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(IDPLoginPolicyLink.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListLoginPolicyIDPsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => IDPLoginPolicyLink.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListLoginPolicyIDPsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? IDPLoginPolicyLink.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListLoginPolicyIDPsResponse>): ListLoginPolicyIDPsResponse {
    return ListLoginPolicyIDPsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListLoginPolicyIDPsResponse>): ListLoginPolicyIDPsResponse {
    const message = createBaseListLoginPolicyIDPsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => IDPLoginPolicyLink.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAddIDPToLoginPolicyRequest(): AddIDPToLoginPolicyRequest {
  return { idpId: "" };
}

export const AddIDPToLoginPolicyRequest = {
  encode(message: AddIDPToLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddIDPToLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddIDPToLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddIDPToLoginPolicyRequest {
    return { idpId: isSet(object.idpId) ? String(object.idpId) : "" };
  },

  toJSON(message: AddIDPToLoginPolicyRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<AddIDPToLoginPolicyRequest>): AddIDPToLoginPolicyRequest {
    return AddIDPToLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddIDPToLoginPolicyRequest>): AddIDPToLoginPolicyRequest {
    const message = createBaseAddIDPToLoginPolicyRequest();
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseAddIDPToLoginPolicyResponse(): AddIDPToLoginPolicyResponse {
  return { details: undefined };
}

export const AddIDPToLoginPolicyResponse = {
  encode(message: AddIDPToLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddIDPToLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddIDPToLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddIDPToLoginPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddIDPToLoginPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddIDPToLoginPolicyResponse>): AddIDPToLoginPolicyResponse {
    return AddIDPToLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddIDPToLoginPolicyResponse>): AddIDPToLoginPolicyResponse {
    const message = createBaseAddIDPToLoginPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveIDPFromLoginPolicyRequest(): RemoveIDPFromLoginPolicyRequest {
  return { idpId: "" };
}

export const RemoveIDPFromLoginPolicyRequest = {
  encode(message: RemoveIDPFromLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveIDPFromLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveIDPFromLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveIDPFromLoginPolicyRequest {
    return { idpId: isSet(object.idpId) ? String(object.idpId) : "" };
  },

  toJSON(message: RemoveIDPFromLoginPolicyRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<RemoveIDPFromLoginPolicyRequest>): RemoveIDPFromLoginPolicyRequest {
    return RemoveIDPFromLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveIDPFromLoginPolicyRequest>): RemoveIDPFromLoginPolicyRequest {
    const message = createBaseRemoveIDPFromLoginPolicyRequest();
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseRemoveIDPFromLoginPolicyResponse(): RemoveIDPFromLoginPolicyResponse {
  return { details: undefined };
}

export const RemoveIDPFromLoginPolicyResponse = {
  encode(message: RemoveIDPFromLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveIDPFromLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveIDPFromLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveIDPFromLoginPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveIDPFromLoginPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveIDPFromLoginPolicyResponse>): RemoveIDPFromLoginPolicyResponse {
    return RemoveIDPFromLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveIDPFromLoginPolicyResponse>): RemoveIDPFromLoginPolicyResponse {
    const message = createBaseRemoveIDPFromLoginPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListLoginPolicySecondFactorsRequest(): ListLoginPolicySecondFactorsRequest {
  return {};
}

export const ListLoginPolicySecondFactorsRequest = {
  encode(_: ListLoginPolicySecondFactorsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListLoginPolicySecondFactorsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListLoginPolicySecondFactorsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListLoginPolicySecondFactorsRequest {
    return {};
  },

  toJSON(_: ListLoginPolicySecondFactorsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListLoginPolicySecondFactorsRequest>): ListLoginPolicySecondFactorsRequest {
    return ListLoginPolicySecondFactorsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListLoginPolicySecondFactorsRequest>): ListLoginPolicySecondFactorsRequest {
    const message = createBaseListLoginPolicySecondFactorsRequest();
    return message;
  },
};

function createBaseListLoginPolicySecondFactorsResponse(): ListLoginPolicySecondFactorsResponse {
  return { details: undefined, result: [] };
}

export const ListLoginPolicySecondFactorsResponse = {
  encode(message: ListLoginPolicySecondFactorsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    writer.uint32(18).fork();
    for (const v of message.result) {
      writer.int32(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListLoginPolicySecondFactorsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListLoginPolicySecondFactorsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.result.push(reader.int32() as any);
            }
          } else {
            message.result.push(reader.int32() as any);
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListLoginPolicySecondFactorsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => secondFactorTypeFromJSON(e)) : [],
    };
  },

  toJSON(message: ListLoginPolicySecondFactorsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => secondFactorTypeToJSON(e));
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListLoginPolicySecondFactorsResponse>): ListLoginPolicySecondFactorsResponse {
    return ListLoginPolicySecondFactorsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListLoginPolicySecondFactorsResponse>): ListLoginPolicySecondFactorsResponse {
    const message = createBaseListLoginPolicySecondFactorsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => e) || [];
    return message;
  },
};

function createBaseAddSecondFactorToLoginPolicyRequest(): AddSecondFactorToLoginPolicyRequest {
  return { type: 0 };
}

export const AddSecondFactorToLoginPolicyRequest = {
  encode(message: AddSecondFactorToLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddSecondFactorToLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddSecondFactorToLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddSecondFactorToLoginPolicyRequest {
    return { type: isSet(object.type) ? secondFactorTypeFromJSON(object.type) : 0 };
  },

  toJSON(message: AddSecondFactorToLoginPolicyRequest): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = secondFactorTypeToJSON(message.type));
    return obj;
  },

  create(base?: DeepPartial<AddSecondFactorToLoginPolicyRequest>): AddSecondFactorToLoginPolicyRequest {
    return AddSecondFactorToLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddSecondFactorToLoginPolicyRequest>): AddSecondFactorToLoginPolicyRequest {
    const message = createBaseAddSecondFactorToLoginPolicyRequest();
    message.type = object.type ?? 0;
    return message;
  },
};

function createBaseAddSecondFactorToLoginPolicyResponse(): AddSecondFactorToLoginPolicyResponse {
  return { details: undefined };
}

export const AddSecondFactorToLoginPolicyResponse = {
  encode(message: AddSecondFactorToLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddSecondFactorToLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddSecondFactorToLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddSecondFactorToLoginPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddSecondFactorToLoginPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddSecondFactorToLoginPolicyResponse>): AddSecondFactorToLoginPolicyResponse {
    return AddSecondFactorToLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddSecondFactorToLoginPolicyResponse>): AddSecondFactorToLoginPolicyResponse {
    const message = createBaseAddSecondFactorToLoginPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveSecondFactorFromLoginPolicyRequest(): RemoveSecondFactorFromLoginPolicyRequest {
  return { type: 0 };
}

export const RemoveSecondFactorFromLoginPolicyRequest = {
  encode(message: RemoveSecondFactorFromLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveSecondFactorFromLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveSecondFactorFromLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveSecondFactorFromLoginPolicyRequest {
    return { type: isSet(object.type) ? secondFactorTypeFromJSON(object.type) : 0 };
  },

  toJSON(message: RemoveSecondFactorFromLoginPolicyRequest): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = secondFactorTypeToJSON(message.type));
    return obj;
  },

  create(base?: DeepPartial<RemoveSecondFactorFromLoginPolicyRequest>): RemoveSecondFactorFromLoginPolicyRequest {
    return RemoveSecondFactorFromLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveSecondFactorFromLoginPolicyRequest>): RemoveSecondFactorFromLoginPolicyRequest {
    const message = createBaseRemoveSecondFactorFromLoginPolicyRequest();
    message.type = object.type ?? 0;
    return message;
  },
};

function createBaseRemoveSecondFactorFromLoginPolicyResponse(): RemoveSecondFactorFromLoginPolicyResponse {
  return { details: undefined };
}

export const RemoveSecondFactorFromLoginPolicyResponse = {
  encode(message: RemoveSecondFactorFromLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveSecondFactorFromLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveSecondFactorFromLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveSecondFactorFromLoginPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveSecondFactorFromLoginPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveSecondFactorFromLoginPolicyResponse>): RemoveSecondFactorFromLoginPolicyResponse {
    return RemoveSecondFactorFromLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<RemoveSecondFactorFromLoginPolicyResponse>,
  ): RemoveSecondFactorFromLoginPolicyResponse {
    const message = createBaseRemoveSecondFactorFromLoginPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListLoginPolicyMultiFactorsRequest(): ListLoginPolicyMultiFactorsRequest {
  return {};
}

export const ListLoginPolicyMultiFactorsRequest = {
  encode(_: ListLoginPolicyMultiFactorsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListLoginPolicyMultiFactorsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListLoginPolicyMultiFactorsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListLoginPolicyMultiFactorsRequest {
    return {};
  },

  toJSON(_: ListLoginPolicyMultiFactorsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListLoginPolicyMultiFactorsRequest>): ListLoginPolicyMultiFactorsRequest {
    return ListLoginPolicyMultiFactorsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListLoginPolicyMultiFactorsRequest>): ListLoginPolicyMultiFactorsRequest {
    const message = createBaseListLoginPolicyMultiFactorsRequest();
    return message;
  },
};

function createBaseListLoginPolicyMultiFactorsResponse(): ListLoginPolicyMultiFactorsResponse {
  return { details: undefined, result: [] };
}

export const ListLoginPolicyMultiFactorsResponse = {
  encode(message: ListLoginPolicyMultiFactorsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    writer.uint32(18).fork();
    for (const v of message.result) {
      writer.int32(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListLoginPolicyMultiFactorsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListLoginPolicyMultiFactorsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.result.push(reader.int32() as any);
            }
          } else {
            message.result.push(reader.int32() as any);
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListLoginPolicyMultiFactorsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => multiFactorTypeFromJSON(e)) : [],
    };
  },

  toJSON(message: ListLoginPolicyMultiFactorsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => multiFactorTypeToJSON(e));
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListLoginPolicyMultiFactorsResponse>): ListLoginPolicyMultiFactorsResponse {
    return ListLoginPolicyMultiFactorsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListLoginPolicyMultiFactorsResponse>): ListLoginPolicyMultiFactorsResponse {
    const message = createBaseListLoginPolicyMultiFactorsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => e) || [];
    return message;
  },
};

function createBaseAddMultiFactorToLoginPolicyRequest(): AddMultiFactorToLoginPolicyRequest {
  return { type: 0 };
}

export const AddMultiFactorToLoginPolicyRequest = {
  encode(message: AddMultiFactorToLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMultiFactorToLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMultiFactorToLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddMultiFactorToLoginPolicyRequest {
    return { type: isSet(object.type) ? multiFactorTypeFromJSON(object.type) : 0 };
  },

  toJSON(message: AddMultiFactorToLoginPolicyRequest): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = multiFactorTypeToJSON(message.type));
    return obj;
  },

  create(base?: DeepPartial<AddMultiFactorToLoginPolicyRequest>): AddMultiFactorToLoginPolicyRequest {
    return AddMultiFactorToLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddMultiFactorToLoginPolicyRequest>): AddMultiFactorToLoginPolicyRequest {
    const message = createBaseAddMultiFactorToLoginPolicyRequest();
    message.type = object.type ?? 0;
    return message;
  },
};

function createBaseAddMultiFactorToLoginPolicyResponse(): AddMultiFactorToLoginPolicyResponse {
  return { details: undefined };
}

export const AddMultiFactorToLoginPolicyResponse = {
  encode(message: AddMultiFactorToLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMultiFactorToLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMultiFactorToLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddMultiFactorToLoginPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddMultiFactorToLoginPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddMultiFactorToLoginPolicyResponse>): AddMultiFactorToLoginPolicyResponse {
    return AddMultiFactorToLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddMultiFactorToLoginPolicyResponse>): AddMultiFactorToLoginPolicyResponse {
    const message = createBaseAddMultiFactorToLoginPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveMultiFactorFromLoginPolicyRequest(): RemoveMultiFactorFromLoginPolicyRequest {
  return { type: 0 };
}

export const RemoveMultiFactorFromLoginPolicyRequest = {
  encode(message: RemoveMultiFactorFromLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMultiFactorFromLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMultiFactorFromLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMultiFactorFromLoginPolicyRequest {
    return { type: isSet(object.type) ? multiFactorTypeFromJSON(object.type) : 0 };
  },

  toJSON(message: RemoveMultiFactorFromLoginPolicyRequest): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = multiFactorTypeToJSON(message.type));
    return obj;
  },

  create(base?: DeepPartial<RemoveMultiFactorFromLoginPolicyRequest>): RemoveMultiFactorFromLoginPolicyRequest {
    return RemoveMultiFactorFromLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMultiFactorFromLoginPolicyRequest>): RemoveMultiFactorFromLoginPolicyRequest {
    const message = createBaseRemoveMultiFactorFromLoginPolicyRequest();
    message.type = object.type ?? 0;
    return message;
  },
};

function createBaseRemoveMultiFactorFromLoginPolicyResponse(): RemoveMultiFactorFromLoginPolicyResponse {
  return { details: undefined };
}

export const RemoveMultiFactorFromLoginPolicyResponse = {
  encode(message: RemoveMultiFactorFromLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMultiFactorFromLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMultiFactorFromLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMultiFactorFromLoginPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMultiFactorFromLoginPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMultiFactorFromLoginPolicyResponse>): RemoveMultiFactorFromLoginPolicyResponse {
    return RemoveMultiFactorFromLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMultiFactorFromLoginPolicyResponse>): RemoveMultiFactorFromLoginPolicyResponse {
    const message = createBaseRemoveMultiFactorFromLoginPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetPasswordComplexityPolicyRequest(): GetPasswordComplexityPolicyRequest {
  return {};
}

export const GetPasswordComplexityPolicyRequest = {
  encode(_: GetPasswordComplexityPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPasswordComplexityPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPasswordComplexityPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetPasswordComplexityPolicyRequest {
    return {};
  },

  toJSON(_: GetPasswordComplexityPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetPasswordComplexityPolicyRequest>): GetPasswordComplexityPolicyRequest {
    return GetPasswordComplexityPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetPasswordComplexityPolicyRequest>): GetPasswordComplexityPolicyRequest {
    const message = createBaseGetPasswordComplexityPolicyRequest();
    return message;
  },
};

function createBaseGetPasswordComplexityPolicyResponse(): GetPasswordComplexityPolicyResponse {
  return { policy: undefined };
}

export const GetPasswordComplexityPolicyResponse = {
  encode(message: GetPasswordComplexityPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      PasswordComplexityPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPasswordComplexityPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPasswordComplexityPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = PasswordComplexityPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetPasswordComplexityPolicyResponse {
    return { policy: isSet(object.policy) ? PasswordComplexityPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetPasswordComplexityPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined &&
      (obj.policy = message.policy ? PasswordComplexityPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetPasswordComplexityPolicyResponse>): GetPasswordComplexityPolicyResponse {
    return GetPasswordComplexityPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetPasswordComplexityPolicyResponse>): GetPasswordComplexityPolicyResponse {
    const message = createBaseGetPasswordComplexityPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? PasswordComplexityPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdatePasswordComplexityPolicyRequest(): UpdatePasswordComplexityPolicyRequest {
  return { minLength: 0, hasUppercase: false, hasLowercase: false, hasNumber: false, hasSymbol: false };
}

export const UpdatePasswordComplexityPolicyRequest = {
  encode(message: UpdatePasswordComplexityPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.minLength !== 0) {
      writer.uint32(8).uint32(message.minLength);
    }
    if (message.hasUppercase === true) {
      writer.uint32(16).bool(message.hasUppercase);
    }
    if (message.hasLowercase === true) {
      writer.uint32(24).bool(message.hasLowercase);
    }
    if (message.hasNumber === true) {
      writer.uint32(32).bool(message.hasNumber);
    }
    if (message.hasSymbol === true) {
      writer.uint32(40).bool(message.hasSymbol);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePasswordComplexityPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePasswordComplexityPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.minLength = reader.uint32();
          break;
        case 2:
          message.hasUppercase = reader.bool();
          break;
        case 3:
          message.hasLowercase = reader.bool();
          break;
        case 4:
          message.hasNumber = reader.bool();
          break;
        case 5:
          message.hasSymbol = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdatePasswordComplexityPolicyRequest {
    return {
      minLength: isSet(object.minLength) ? Number(object.minLength) : 0,
      hasUppercase: isSet(object.hasUppercase) ? Boolean(object.hasUppercase) : false,
      hasLowercase: isSet(object.hasLowercase) ? Boolean(object.hasLowercase) : false,
      hasNumber: isSet(object.hasNumber) ? Boolean(object.hasNumber) : false,
      hasSymbol: isSet(object.hasSymbol) ? Boolean(object.hasSymbol) : false,
    };
  },

  toJSON(message: UpdatePasswordComplexityPolicyRequest): unknown {
    const obj: any = {};
    message.minLength !== undefined && (obj.minLength = Math.round(message.minLength));
    message.hasUppercase !== undefined && (obj.hasUppercase = message.hasUppercase);
    message.hasLowercase !== undefined && (obj.hasLowercase = message.hasLowercase);
    message.hasNumber !== undefined && (obj.hasNumber = message.hasNumber);
    message.hasSymbol !== undefined && (obj.hasSymbol = message.hasSymbol);
    return obj;
  },

  create(base?: DeepPartial<UpdatePasswordComplexityPolicyRequest>): UpdatePasswordComplexityPolicyRequest {
    return UpdatePasswordComplexityPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdatePasswordComplexityPolicyRequest>): UpdatePasswordComplexityPolicyRequest {
    const message = createBaseUpdatePasswordComplexityPolicyRequest();
    message.minLength = object.minLength ?? 0;
    message.hasUppercase = object.hasUppercase ?? false;
    message.hasLowercase = object.hasLowercase ?? false;
    message.hasNumber = object.hasNumber ?? false;
    message.hasSymbol = object.hasSymbol ?? false;
    return message;
  },
};

function createBaseUpdatePasswordComplexityPolicyResponse(): UpdatePasswordComplexityPolicyResponse {
  return { details: undefined };
}

export const UpdatePasswordComplexityPolicyResponse = {
  encode(message: UpdatePasswordComplexityPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePasswordComplexityPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePasswordComplexityPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdatePasswordComplexityPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdatePasswordComplexityPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdatePasswordComplexityPolicyResponse>): UpdatePasswordComplexityPolicyResponse {
    return UpdatePasswordComplexityPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdatePasswordComplexityPolicyResponse>): UpdatePasswordComplexityPolicyResponse {
    const message = createBaseUpdatePasswordComplexityPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetPasswordAgePolicyRequest(): GetPasswordAgePolicyRequest {
  return {};
}

export const GetPasswordAgePolicyRequest = {
  encode(_: GetPasswordAgePolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPasswordAgePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPasswordAgePolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetPasswordAgePolicyRequest {
    return {};
  },

  toJSON(_: GetPasswordAgePolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetPasswordAgePolicyRequest>): GetPasswordAgePolicyRequest {
    return GetPasswordAgePolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetPasswordAgePolicyRequest>): GetPasswordAgePolicyRequest {
    const message = createBaseGetPasswordAgePolicyRequest();
    return message;
  },
};

function createBaseGetPasswordAgePolicyResponse(): GetPasswordAgePolicyResponse {
  return { policy: undefined };
}

export const GetPasswordAgePolicyResponse = {
  encode(message: GetPasswordAgePolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      PasswordAgePolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPasswordAgePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPasswordAgePolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = PasswordAgePolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetPasswordAgePolicyResponse {
    return { policy: isSet(object.policy) ? PasswordAgePolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetPasswordAgePolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined &&
      (obj.policy = message.policy ? PasswordAgePolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetPasswordAgePolicyResponse>): GetPasswordAgePolicyResponse {
    return GetPasswordAgePolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetPasswordAgePolicyResponse>): GetPasswordAgePolicyResponse {
    const message = createBaseGetPasswordAgePolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? PasswordAgePolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdatePasswordAgePolicyRequest(): UpdatePasswordAgePolicyRequest {
  return { maxAgeDays: 0, expireWarnDays: 0 };
}

export const UpdatePasswordAgePolicyRequest = {
  encode(message: UpdatePasswordAgePolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.maxAgeDays !== 0) {
      writer.uint32(8).uint32(message.maxAgeDays);
    }
    if (message.expireWarnDays !== 0) {
      writer.uint32(16).uint32(message.expireWarnDays);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePasswordAgePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePasswordAgePolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.maxAgeDays = reader.uint32();
          break;
        case 2:
          message.expireWarnDays = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdatePasswordAgePolicyRequest {
    return {
      maxAgeDays: isSet(object.maxAgeDays) ? Number(object.maxAgeDays) : 0,
      expireWarnDays: isSet(object.expireWarnDays) ? Number(object.expireWarnDays) : 0,
    };
  },

  toJSON(message: UpdatePasswordAgePolicyRequest): unknown {
    const obj: any = {};
    message.maxAgeDays !== undefined && (obj.maxAgeDays = Math.round(message.maxAgeDays));
    message.expireWarnDays !== undefined && (obj.expireWarnDays = Math.round(message.expireWarnDays));
    return obj;
  },

  create(base?: DeepPartial<UpdatePasswordAgePolicyRequest>): UpdatePasswordAgePolicyRequest {
    return UpdatePasswordAgePolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdatePasswordAgePolicyRequest>): UpdatePasswordAgePolicyRequest {
    const message = createBaseUpdatePasswordAgePolicyRequest();
    message.maxAgeDays = object.maxAgeDays ?? 0;
    message.expireWarnDays = object.expireWarnDays ?? 0;
    return message;
  },
};

function createBaseUpdatePasswordAgePolicyResponse(): UpdatePasswordAgePolicyResponse {
  return { details: undefined };
}

export const UpdatePasswordAgePolicyResponse = {
  encode(message: UpdatePasswordAgePolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePasswordAgePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePasswordAgePolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdatePasswordAgePolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdatePasswordAgePolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdatePasswordAgePolicyResponse>): UpdatePasswordAgePolicyResponse {
    return UpdatePasswordAgePolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdatePasswordAgePolicyResponse>): UpdatePasswordAgePolicyResponse {
    const message = createBaseUpdatePasswordAgePolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetLockoutPolicyRequest(): GetLockoutPolicyRequest {
  return {};
}

export const GetLockoutPolicyRequest = {
  encode(_: GetLockoutPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLockoutPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLockoutPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetLockoutPolicyRequest {
    return {};
  },

  toJSON(_: GetLockoutPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetLockoutPolicyRequest>): GetLockoutPolicyRequest {
    return GetLockoutPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetLockoutPolicyRequest>): GetLockoutPolicyRequest {
    const message = createBaseGetLockoutPolicyRequest();
    return message;
  },
};

function createBaseGetLockoutPolicyResponse(): GetLockoutPolicyResponse {
  return { policy: undefined };
}

export const GetLockoutPolicyResponse = {
  encode(message: GetLockoutPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      LockoutPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLockoutPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLockoutPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = LockoutPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetLockoutPolicyResponse {
    return { policy: isSet(object.policy) ? LockoutPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetLockoutPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? LockoutPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetLockoutPolicyResponse>): GetLockoutPolicyResponse {
    return GetLockoutPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetLockoutPolicyResponse>): GetLockoutPolicyResponse {
    const message = createBaseGetLockoutPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? LockoutPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdateLockoutPolicyRequest(): UpdateLockoutPolicyRequest {
  return { maxPasswordAttempts: 0 };
}

export const UpdateLockoutPolicyRequest = {
  encode(message: UpdateLockoutPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.maxPasswordAttempts !== 0) {
      writer.uint32(8).uint32(message.maxPasswordAttempts);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateLockoutPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateLockoutPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.maxPasswordAttempts = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateLockoutPolicyRequest {
    return { maxPasswordAttempts: isSet(object.maxPasswordAttempts) ? Number(object.maxPasswordAttempts) : 0 };
  },

  toJSON(message: UpdateLockoutPolicyRequest): unknown {
    const obj: any = {};
    message.maxPasswordAttempts !== undefined && (obj.maxPasswordAttempts = Math.round(message.maxPasswordAttempts));
    return obj;
  },

  create(base?: DeepPartial<UpdateLockoutPolicyRequest>): UpdateLockoutPolicyRequest {
    return UpdateLockoutPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateLockoutPolicyRequest>): UpdateLockoutPolicyRequest {
    const message = createBaseUpdateLockoutPolicyRequest();
    message.maxPasswordAttempts = object.maxPasswordAttempts ?? 0;
    return message;
  },
};

function createBaseUpdateLockoutPolicyResponse(): UpdateLockoutPolicyResponse {
  return { details: undefined };
}

export const UpdateLockoutPolicyResponse = {
  encode(message: UpdateLockoutPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateLockoutPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateLockoutPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateLockoutPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateLockoutPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateLockoutPolicyResponse>): UpdateLockoutPolicyResponse {
    return UpdateLockoutPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateLockoutPolicyResponse>): UpdateLockoutPolicyResponse {
    const message = createBaseUpdateLockoutPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetPrivacyPolicyRequest(): GetPrivacyPolicyRequest {
  return {};
}

export const GetPrivacyPolicyRequest = {
  encode(_: GetPrivacyPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPrivacyPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPrivacyPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetPrivacyPolicyRequest {
    return {};
  },

  toJSON(_: GetPrivacyPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetPrivacyPolicyRequest>): GetPrivacyPolicyRequest {
    return GetPrivacyPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetPrivacyPolicyRequest>): GetPrivacyPolicyRequest {
    const message = createBaseGetPrivacyPolicyRequest();
    return message;
  },
};

function createBaseGetPrivacyPolicyResponse(): GetPrivacyPolicyResponse {
  return { policy: undefined };
}

export const GetPrivacyPolicyResponse = {
  encode(message: GetPrivacyPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      PrivacyPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPrivacyPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPrivacyPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = PrivacyPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetPrivacyPolicyResponse {
    return { policy: isSet(object.policy) ? PrivacyPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetPrivacyPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? PrivacyPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetPrivacyPolicyResponse>): GetPrivacyPolicyResponse {
    return GetPrivacyPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetPrivacyPolicyResponse>): GetPrivacyPolicyResponse {
    const message = createBaseGetPrivacyPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? PrivacyPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdatePrivacyPolicyRequest(): UpdatePrivacyPolicyRequest {
  return { tosLink: "", privacyLink: "", helpLink: "" };
}

export const UpdatePrivacyPolicyRequest = {
  encode(message: UpdatePrivacyPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tosLink !== "") {
      writer.uint32(10).string(message.tosLink);
    }
    if (message.privacyLink !== "") {
      writer.uint32(18).string(message.privacyLink);
    }
    if (message.helpLink !== "") {
      writer.uint32(26).string(message.helpLink);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePrivacyPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePrivacyPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tosLink = reader.string();
          break;
        case 2:
          message.privacyLink = reader.string();
          break;
        case 3:
          message.helpLink = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdatePrivacyPolicyRequest {
    return {
      tosLink: isSet(object.tosLink) ? String(object.tosLink) : "",
      privacyLink: isSet(object.privacyLink) ? String(object.privacyLink) : "",
      helpLink: isSet(object.helpLink) ? String(object.helpLink) : "",
    };
  },

  toJSON(message: UpdatePrivacyPolicyRequest): unknown {
    const obj: any = {};
    message.tosLink !== undefined && (obj.tosLink = message.tosLink);
    message.privacyLink !== undefined && (obj.privacyLink = message.privacyLink);
    message.helpLink !== undefined && (obj.helpLink = message.helpLink);
    return obj;
  },

  create(base?: DeepPartial<UpdatePrivacyPolicyRequest>): UpdatePrivacyPolicyRequest {
    return UpdatePrivacyPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdatePrivacyPolicyRequest>): UpdatePrivacyPolicyRequest {
    const message = createBaseUpdatePrivacyPolicyRequest();
    message.tosLink = object.tosLink ?? "";
    message.privacyLink = object.privacyLink ?? "";
    message.helpLink = object.helpLink ?? "";
    return message;
  },
};

function createBaseUpdatePrivacyPolicyResponse(): UpdatePrivacyPolicyResponse {
  return { details: undefined };
}

export const UpdatePrivacyPolicyResponse = {
  encode(message: UpdatePrivacyPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePrivacyPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePrivacyPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdatePrivacyPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdatePrivacyPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdatePrivacyPolicyResponse>): UpdatePrivacyPolicyResponse {
    return UpdatePrivacyPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdatePrivacyPolicyResponse>): UpdatePrivacyPolicyResponse {
    const message = createBaseUpdatePrivacyPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddNotificationPolicyRequest(): AddNotificationPolicyRequest {
  return { passwordChange: false };
}

export const AddNotificationPolicyRequest = {
  encode(message: AddNotificationPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.passwordChange === true) {
      writer.uint32(8).bool(message.passwordChange);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddNotificationPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddNotificationPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.passwordChange = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddNotificationPolicyRequest {
    return { passwordChange: isSet(object.passwordChange) ? Boolean(object.passwordChange) : false };
  },

  toJSON(message: AddNotificationPolicyRequest): unknown {
    const obj: any = {};
    message.passwordChange !== undefined && (obj.passwordChange = message.passwordChange);
    return obj;
  },

  create(base?: DeepPartial<AddNotificationPolicyRequest>): AddNotificationPolicyRequest {
    return AddNotificationPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddNotificationPolicyRequest>): AddNotificationPolicyRequest {
    const message = createBaseAddNotificationPolicyRequest();
    message.passwordChange = object.passwordChange ?? false;
    return message;
  },
};

function createBaseAddNotificationPolicyResponse(): AddNotificationPolicyResponse {
  return { details: undefined };
}

export const AddNotificationPolicyResponse = {
  encode(message: AddNotificationPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddNotificationPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddNotificationPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddNotificationPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddNotificationPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddNotificationPolicyResponse>): AddNotificationPolicyResponse {
    return AddNotificationPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddNotificationPolicyResponse>): AddNotificationPolicyResponse {
    const message = createBaseAddNotificationPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetNotificationPolicyRequest(): GetNotificationPolicyRequest {
  return {};
}

export const GetNotificationPolicyRequest = {
  encode(_: GetNotificationPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNotificationPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNotificationPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetNotificationPolicyRequest {
    return {};
  },

  toJSON(_: GetNotificationPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetNotificationPolicyRequest>): GetNotificationPolicyRequest {
    return GetNotificationPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetNotificationPolicyRequest>): GetNotificationPolicyRequest {
    const message = createBaseGetNotificationPolicyRequest();
    return message;
  },
};

function createBaseGetNotificationPolicyResponse(): GetNotificationPolicyResponse {
  return { policy: undefined };
}

export const GetNotificationPolicyResponse = {
  encode(message: GetNotificationPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      NotificationPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNotificationPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNotificationPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = NotificationPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetNotificationPolicyResponse {
    return { policy: isSet(object.policy) ? NotificationPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetNotificationPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined &&
      (obj.policy = message.policy ? NotificationPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetNotificationPolicyResponse>): GetNotificationPolicyResponse {
    return GetNotificationPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetNotificationPolicyResponse>): GetNotificationPolicyResponse {
    const message = createBaseGetNotificationPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? NotificationPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdateNotificationPolicyRequest(): UpdateNotificationPolicyRequest {
  return { passwordChange: false };
}

export const UpdateNotificationPolicyRequest = {
  encode(message: UpdateNotificationPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.passwordChange === true) {
      writer.uint32(8).bool(message.passwordChange);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateNotificationPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateNotificationPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.passwordChange = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateNotificationPolicyRequest {
    return { passwordChange: isSet(object.passwordChange) ? Boolean(object.passwordChange) : false };
  },

  toJSON(message: UpdateNotificationPolicyRequest): unknown {
    const obj: any = {};
    message.passwordChange !== undefined && (obj.passwordChange = message.passwordChange);
    return obj;
  },

  create(base?: DeepPartial<UpdateNotificationPolicyRequest>): UpdateNotificationPolicyRequest {
    return UpdateNotificationPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateNotificationPolicyRequest>): UpdateNotificationPolicyRequest {
    const message = createBaseUpdateNotificationPolicyRequest();
    message.passwordChange = object.passwordChange ?? false;
    return message;
  },
};

function createBaseUpdateNotificationPolicyResponse(): UpdateNotificationPolicyResponse {
  return { details: undefined };
}

export const UpdateNotificationPolicyResponse = {
  encode(message: UpdateNotificationPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateNotificationPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateNotificationPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateNotificationPolicyResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateNotificationPolicyResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateNotificationPolicyResponse>): UpdateNotificationPolicyResponse {
    return UpdateNotificationPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateNotificationPolicyResponse>): UpdateNotificationPolicyResponse {
    const message = createBaseUpdateNotificationPolicyResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultInitMessageTextRequest(): GetDefaultInitMessageTextRequest {
  return { language: "" };
}

export const GetDefaultInitMessageTextRequest = {
  encode(message: GetDefaultInitMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultInitMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultInitMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultInitMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultInitMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultInitMessageTextRequest>): GetDefaultInitMessageTextRequest {
    return GetDefaultInitMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultInitMessageTextRequest>): GetDefaultInitMessageTextRequest {
    const message = createBaseGetDefaultInitMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetDefaultInitMessageTextResponse(): GetDefaultInitMessageTextResponse {
  return { customText: undefined };
}

export const GetDefaultInitMessageTextResponse = {
  encode(message: GetDefaultInitMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultInitMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultInitMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultInitMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetDefaultInitMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultInitMessageTextResponse>): GetDefaultInitMessageTextResponse {
    return GetDefaultInitMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultInitMessageTextResponse>): GetDefaultInitMessageTextResponse {
    const message = createBaseGetDefaultInitMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseGetCustomInitMessageTextRequest(): GetCustomInitMessageTextRequest {
  return { language: "" };
}

export const GetCustomInitMessageTextRequest = {
  encode(message: GetCustomInitMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomInitMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomInitMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomInitMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetCustomInitMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetCustomInitMessageTextRequest>): GetCustomInitMessageTextRequest {
    return GetCustomInitMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomInitMessageTextRequest>): GetCustomInitMessageTextRequest {
    const message = createBaseGetCustomInitMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetCustomInitMessageTextResponse(): GetCustomInitMessageTextResponse {
  return { customText: undefined };
}

export const GetCustomInitMessageTextResponse = {
  encode(message: GetCustomInitMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomInitMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomInitMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomInitMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetCustomInitMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetCustomInitMessageTextResponse>): GetCustomInitMessageTextResponse {
    return GetCustomInitMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomInitMessageTextResponse>): GetCustomInitMessageTextResponse {
    const message = createBaseGetCustomInitMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseSetDefaultInitMessageTextRequest(): SetDefaultInitMessageTextRequest {
  return {
    language: "",
    title: "",
    preHeader: "",
    subject: "",
    greeting: "",
    text: "",
    buttonText: "",
    footerText: "",
  };
}

export const SetDefaultInitMessageTextRequest = {
  encode(message: SetDefaultInitMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.preHeader !== "") {
      writer.uint32(26).string(message.preHeader);
    }
    if (message.subject !== "") {
      writer.uint32(34).string(message.subject);
    }
    if (message.greeting !== "") {
      writer.uint32(42).string(message.greeting);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.buttonText !== "") {
      writer.uint32(58).string(message.buttonText);
    }
    if (message.footerText !== "") {
      writer.uint32(66).string(message.footerText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultInitMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultInitMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.preHeader = reader.string();
          break;
        case 4:
          message.subject = reader.string();
          break;
        case 5:
          message.greeting = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.buttonText = reader.string();
          break;
        case 8:
          message.footerText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultInitMessageTextRequest {
    return {
      language: isSet(object.language) ? String(object.language) : "",
      title: isSet(object.title) ? String(object.title) : "",
      preHeader: isSet(object.preHeader) ? String(object.preHeader) : "",
      subject: isSet(object.subject) ? String(object.subject) : "",
      greeting: isSet(object.greeting) ? String(object.greeting) : "",
      text: isSet(object.text) ? String(object.text) : "",
      buttonText: isSet(object.buttonText) ? String(object.buttonText) : "",
      footerText: isSet(object.footerText) ? String(object.footerText) : "",
    };
  },

  toJSON(message: SetDefaultInitMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    message.title !== undefined && (obj.title = message.title);
    message.preHeader !== undefined && (obj.preHeader = message.preHeader);
    message.subject !== undefined && (obj.subject = message.subject);
    message.greeting !== undefined && (obj.greeting = message.greeting);
    message.text !== undefined && (obj.text = message.text);
    message.buttonText !== undefined && (obj.buttonText = message.buttonText);
    message.footerText !== undefined && (obj.footerText = message.footerText);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultInitMessageTextRequest>): SetDefaultInitMessageTextRequest {
    return SetDefaultInitMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultInitMessageTextRequest>): SetDefaultInitMessageTextRequest {
    const message = createBaseSetDefaultInitMessageTextRequest();
    message.language = object.language ?? "";
    message.title = object.title ?? "";
    message.preHeader = object.preHeader ?? "";
    message.subject = object.subject ?? "";
    message.greeting = object.greeting ?? "";
    message.text = object.text ?? "";
    message.buttonText = object.buttonText ?? "";
    message.footerText = object.footerText ?? "";
    return message;
  },
};

function createBaseSetDefaultInitMessageTextResponse(): SetDefaultInitMessageTextResponse {
  return { details: undefined };
}

export const SetDefaultInitMessageTextResponse = {
  encode(message: SetDefaultInitMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultInitMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultInitMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultInitMessageTextResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultInitMessageTextResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultInitMessageTextResponse>): SetDefaultInitMessageTextResponse {
    return SetDefaultInitMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultInitMessageTextResponse>): SetDefaultInitMessageTextResponse {
    const message = createBaseSetDefaultInitMessageTextResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomInitMessageTextToDefaultRequest(): ResetCustomInitMessageTextToDefaultRequest {
  return { language: "" };
}

export const ResetCustomInitMessageTextToDefaultRequest = {
  encode(message: ResetCustomInitMessageTextToDefaultRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomInitMessageTextToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomInitMessageTextToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomInitMessageTextToDefaultRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: ResetCustomInitMessageTextToDefaultRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<ResetCustomInitMessageTextToDefaultRequest>): ResetCustomInitMessageTextToDefaultRequest {
    return ResetCustomInitMessageTextToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomInitMessageTextToDefaultRequest>,
  ): ResetCustomInitMessageTextToDefaultRequest {
    const message = createBaseResetCustomInitMessageTextToDefaultRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseResetCustomInitMessageTextToDefaultResponse(): ResetCustomInitMessageTextToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomInitMessageTextToDefaultResponse = {
  encode(message: ResetCustomInitMessageTextToDefaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomInitMessageTextToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomInitMessageTextToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomInitMessageTextToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomInitMessageTextToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResetCustomInitMessageTextToDefaultResponse>): ResetCustomInitMessageTextToDefaultResponse {
    return ResetCustomInitMessageTextToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomInitMessageTextToDefaultResponse>,
  ): ResetCustomInitMessageTextToDefaultResponse {
    const message = createBaseResetCustomInitMessageTextToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultPasswordResetMessageTextRequest(): GetDefaultPasswordResetMessageTextRequest {
  return { language: "" };
}

export const GetDefaultPasswordResetMessageTextRequest = {
  encode(message: GetDefaultPasswordResetMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultPasswordResetMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultPasswordResetMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultPasswordResetMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultPasswordResetMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultPasswordResetMessageTextRequest>): GetDefaultPasswordResetMessageTextRequest {
    return GetDefaultPasswordResetMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultPasswordResetMessageTextRequest>,
  ): GetDefaultPasswordResetMessageTextRequest {
    const message = createBaseGetDefaultPasswordResetMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetDefaultPasswordResetMessageTextResponse(): GetDefaultPasswordResetMessageTextResponse {
  return { customText: undefined };
}

export const GetDefaultPasswordResetMessageTextResponse = {
  encode(message: GetDefaultPasswordResetMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultPasswordResetMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultPasswordResetMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultPasswordResetMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetDefaultPasswordResetMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultPasswordResetMessageTextResponse>): GetDefaultPasswordResetMessageTextResponse {
    return GetDefaultPasswordResetMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultPasswordResetMessageTextResponse>,
  ): GetDefaultPasswordResetMessageTextResponse {
    const message = createBaseGetDefaultPasswordResetMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseGetCustomPasswordResetMessageTextRequest(): GetCustomPasswordResetMessageTextRequest {
  return { language: "" };
}

export const GetCustomPasswordResetMessageTextRequest = {
  encode(message: GetCustomPasswordResetMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomPasswordResetMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomPasswordResetMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomPasswordResetMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetCustomPasswordResetMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetCustomPasswordResetMessageTextRequest>): GetCustomPasswordResetMessageTextRequest {
    return GetCustomPasswordResetMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomPasswordResetMessageTextRequest>): GetCustomPasswordResetMessageTextRequest {
    const message = createBaseGetCustomPasswordResetMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetCustomPasswordResetMessageTextResponse(): GetCustomPasswordResetMessageTextResponse {
  return { customText: undefined };
}

export const GetCustomPasswordResetMessageTextResponse = {
  encode(message: GetCustomPasswordResetMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomPasswordResetMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomPasswordResetMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomPasswordResetMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetCustomPasswordResetMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetCustomPasswordResetMessageTextResponse>): GetCustomPasswordResetMessageTextResponse {
    return GetCustomPasswordResetMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetCustomPasswordResetMessageTextResponse>,
  ): GetCustomPasswordResetMessageTextResponse {
    const message = createBaseGetCustomPasswordResetMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseSetDefaultPasswordResetMessageTextRequest(): SetDefaultPasswordResetMessageTextRequest {
  return {
    language: "",
    title: "",
    preHeader: "",
    subject: "",
    greeting: "",
    text: "",
    buttonText: "",
    footerText: "",
  };
}

export const SetDefaultPasswordResetMessageTextRequest = {
  encode(message: SetDefaultPasswordResetMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.preHeader !== "") {
      writer.uint32(26).string(message.preHeader);
    }
    if (message.subject !== "") {
      writer.uint32(34).string(message.subject);
    }
    if (message.greeting !== "") {
      writer.uint32(42).string(message.greeting);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.buttonText !== "") {
      writer.uint32(58).string(message.buttonText);
    }
    if (message.footerText !== "") {
      writer.uint32(66).string(message.footerText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultPasswordResetMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultPasswordResetMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.preHeader = reader.string();
          break;
        case 4:
          message.subject = reader.string();
          break;
        case 5:
          message.greeting = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.buttonText = reader.string();
          break;
        case 8:
          message.footerText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultPasswordResetMessageTextRequest {
    return {
      language: isSet(object.language) ? String(object.language) : "",
      title: isSet(object.title) ? String(object.title) : "",
      preHeader: isSet(object.preHeader) ? String(object.preHeader) : "",
      subject: isSet(object.subject) ? String(object.subject) : "",
      greeting: isSet(object.greeting) ? String(object.greeting) : "",
      text: isSet(object.text) ? String(object.text) : "",
      buttonText: isSet(object.buttonText) ? String(object.buttonText) : "",
      footerText: isSet(object.footerText) ? String(object.footerText) : "",
    };
  },

  toJSON(message: SetDefaultPasswordResetMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    message.title !== undefined && (obj.title = message.title);
    message.preHeader !== undefined && (obj.preHeader = message.preHeader);
    message.subject !== undefined && (obj.subject = message.subject);
    message.greeting !== undefined && (obj.greeting = message.greeting);
    message.text !== undefined && (obj.text = message.text);
    message.buttonText !== undefined && (obj.buttonText = message.buttonText);
    message.footerText !== undefined && (obj.footerText = message.footerText);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultPasswordResetMessageTextRequest>): SetDefaultPasswordResetMessageTextRequest {
    return SetDefaultPasswordResetMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SetDefaultPasswordResetMessageTextRequest>,
  ): SetDefaultPasswordResetMessageTextRequest {
    const message = createBaseSetDefaultPasswordResetMessageTextRequest();
    message.language = object.language ?? "";
    message.title = object.title ?? "";
    message.preHeader = object.preHeader ?? "";
    message.subject = object.subject ?? "";
    message.greeting = object.greeting ?? "";
    message.text = object.text ?? "";
    message.buttonText = object.buttonText ?? "";
    message.footerText = object.footerText ?? "";
    return message;
  },
};

function createBaseSetDefaultPasswordResetMessageTextResponse(): SetDefaultPasswordResetMessageTextResponse {
  return { details: undefined };
}

export const SetDefaultPasswordResetMessageTextResponse = {
  encode(message: SetDefaultPasswordResetMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultPasswordResetMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultPasswordResetMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultPasswordResetMessageTextResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultPasswordResetMessageTextResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultPasswordResetMessageTextResponse>): SetDefaultPasswordResetMessageTextResponse {
    return SetDefaultPasswordResetMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SetDefaultPasswordResetMessageTextResponse>,
  ): SetDefaultPasswordResetMessageTextResponse {
    const message = createBaseSetDefaultPasswordResetMessageTextResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomPasswordResetMessageTextToDefaultRequest(): ResetCustomPasswordResetMessageTextToDefaultRequest {
  return { language: "" };
}

export const ResetCustomPasswordResetMessageTextToDefaultRequest = {
  encode(
    message: ResetCustomPasswordResetMessageTextToDefaultRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomPasswordResetMessageTextToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomPasswordResetMessageTextToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomPasswordResetMessageTextToDefaultRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: ResetCustomPasswordResetMessageTextToDefaultRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomPasswordResetMessageTextToDefaultRequest>,
  ): ResetCustomPasswordResetMessageTextToDefaultRequest {
    return ResetCustomPasswordResetMessageTextToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomPasswordResetMessageTextToDefaultRequest>,
  ): ResetCustomPasswordResetMessageTextToDefaultRequest {
    const message = createBaseResetCustomPasswordResetMessageTextToDefaultRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseResetCustomPasswordResetMessageTextToDefaultResponse(): ResetCustomPasswordResetMessageTextToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomPasswordResetMessageTextToDefaultResponse = {
  encode(
    message: ResetCustomPasswordResetMessageTextToDefaultResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomPasswordResetMessageTextToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomPasswordResetMessageTextToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomPasswordResetMessageTextToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomPasswordResetMessageTextToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomPasswordResetMessageTextToDefaultResponse>,
  ): ResetCustomPasswordResetMessageTextToDefaultResponse {
    return ResetCustomPasswordResetMessageTextToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomPasswordResetMessageTextToDefaultResponse>,
  ): ResetCustomPasswordResetMessageTextToDefaultResponse {
    const message = createBaseResetCustomPasswordResetMessageTextToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultVerifyEmailMessageTextRequest(): GetDefaultVerifyEmailMessageTextRequest {
  return { language: "" };
}

export const GetDefaultVerifyEmailMessageTextRequest = {
  encode(message: GetDefaultVerifyEmailMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultVerifyEmailMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultVerifyEmailMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultVerifyEmailMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultVerifyEmailMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultVerifyEmailMessageTextRequest>): GetDefaultVerifyEmailMessageTextRequest {
    return GetDefaultVerifyEmailMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultVerifyEmailMessageTextRequest>): GetDefaultVerifyEmailMessageTextRequest {
    const message = createBaseGetDefaultVerifyEmailMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetDefaultVerifyEmailMessageTextResponse(): GetDefaultVerifyEmailMessageTextResponse {
  return { customText: undefined };
}

export const GetDefaultVerifyEmailMessageTextResponse = {
  encode(message: GetDefaultVerifyEmailMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultVerifyEmailMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultVerifyEmailMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultVerifyEmailMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetDefaultVerifyEmailMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultVerifyEmailMessageTextResponse>): GetDefaultVerifyEmailMessageTextResponse {
    return GetDefaultVerifyEmailMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultVerifyEmailMessageTextResponse>): GetDefaultVerifyEmailMessageTextResponse {
    const message = createBaseGetDefaultVerifyEmailMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseGetCustomVerifyEmailMessageTextRequest(): GetCustomVerifyEmailMessageTextRequest {
  return { language: "" };
}

export const GetCustomVerifyEmailMessageTextRequest = {
  encode(message: GetCustomVerifyEmailMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomVerifyEmailMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomVerifyEmailMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomVerifyEmailMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetCustomVerifyEmailMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetCustomVerifyEmailMessageTextRequest>): GetCustomVerifyEmailMessageTextRequest {
    return GetCustomVerifyEmailMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomVerifyEmailMessageTextRequest>): GetCustomVerifyEmailMessageTextRequest {
    const message = createBaseGetCustomVerifyEmailMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetCustomVerifyEmailMessageTextResponse(): GetCustomVerifyEmailMessageTextResponse {
  return { customText: undefined };
}

export const GetCustomVerifyEmailMessageTextResponse = {
  encode(message: GetCustomVerifyEmailMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomVerifyEmailMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomVerifyEmailMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomVerifyEmailMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetCustomVerifyEmailMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetCustomVerifyEmailMessageTextResponse>): GetCustomVerifyEmailMessageTextResponse {
    return GetCustomVerifyEmailMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomVerifyEmailMessageTextResponse>): GetCustomVerifyEmailMessageTextResponse {
    const message = createBaseGetCustomVerifyEmailMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseSetDefaultVerifyEmailMessageTextRequest(): SetDefaultVerifyEmailMessageTextRequest {
  return {
    language: "",
    title: "",
    preHeader: "",
    subject: "",
    greeting: "",
    text: "",
    buttonText: "",
    footerText: "",
  };
}

export const SetDefaultVerifyEmailMessageTextRequest = {
  encode(message: SetDefaultVerifyEmailMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.preHeader !== "") {
      writer.uint32(26).string(message.preHeader);
    }
    if (message.subject !== "") {
      writer.uint32(34).string(message.subject);
    }
    if (message.greeting !== "") {
      writer.uint32(42).string(message.greeting);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.buttonText !== "") {
      writer.uint32(58).string(message.buttonText);
    }
    if (message.footerText !== "") {
      writer.uint32(66).string(message.footerText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultVerifyEmailMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultVerifyEmailMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.preHeader = reader.string();
          break;
        case 4:
          message.subject = reader.string();
          break;
        case 5:
          message.greeting = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.buttonText = reader.string();
          break;
        case 8:
          message.footerText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultVerifyEmailMessageTextRequest {
    return {
      language: isSet(object.language) ? String(object.language) : "",
      title: isSet(object.title) ? String(object.title) : "",
      preHeader: isSet(object.preHeader) ? String(object.preHeader) : "",
      subject: isSet(object.subject) ? String(object.subject) : "",
      greeting: isSet(object.greeting) ? String(object.greeting) : "",
      text: isSet(object.text) ? String(object.text) : "",
      buttonText: isSet(object.buttonText) ? String(object.buttonText) : "",
      footerText: isSet(object.footerText) ? String(object.footerText) : "",
    };
  },

  toJSON(message: SetDefaultVerifyEmailMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    message.title !== undefined && (obj.title = message.title);
    message.preHeader !== undefined && (obj.preHeader = message.preHeader);
    message.subject !== undefined && (obj.subject = message.subject);
    message.greeting !== undefined && (obj.greeting = message.greeting);
    message.text !== undefined && (obj.text = message.text);
    message.buttonText !== undefined && (obj.buttonText = message.buttonText);
    message.footerText !== undefined && (obj.footerText = message.footerText);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultVerifyEmailMessageTextRequest>): SetDefaultVerifyEmailMessageTextRequest {
    return SetDefaultVerifyEmailMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultVerifyEmailMessageTextRequest>): SetDefaultVerifyEmailMessageTextRequest {
    const message = createBaseSetDefaultVerifyEmailMessageTextRequest();
    message.language = object.language ?? "";
    message.title = object.title ?? "";
    message.preHeader = object.preHeader ?? "";
    message.subject = object.subject ?? "";
    message.greeting = object.greeting ?? "";
    message.text = object.text ?? "";
    message.buttonText = object.buttonText ?? "";
    message.footerText = object.footerText ?? "";
    return message;
  },
};

function createBaseSetDefaultVerifyEmailMessageTextResponse(): SetDefaultVerifyEmailMessageTextResponse {
  return { details: undefined };
}

export const SetDefaultVerifyEmailMessageTextResponse = {
  encode(message: SetDefaultVerifyEmailMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultVerifyEmailMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultVerifyEmailMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultVerifyEmailMessageTextResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultVerifyEmailMessageTextResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultVerifyEmailMessageTextResponse>): SetDefaultVerifyEmailMessageTextResponse {
    return SetDefaultVerifyEmailMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultVerifyEmailMessageTextResponse>): SetDefaultVerifyEmailMessageTextResponse {
    const message = createBaseSetDefaultVerifyEmailMessageTextResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomVerifyEmailMessageTextToDefaultRequest(): ResetCustomVerifyEmailMessageTextToDefaultRequest {
  return { language: "" };
}

export const ResetCustomVerifyEmailMessageTextToDefaultRequest = {
  encode(
    message: ResetCustomVerifyEmailMessageTextToDefaultRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomVerifyEmailMessageTextToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomVerifyEmailMessageTextToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomVerifyEmailMessageTextToDefaultRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: ResetCustomVerifyEmailMessageTextToDefaultRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomVerifyEmailMessageTextToDefaultRequest>,
  ): ResetCustomVerifyEmailMessageTextToDefaultRequest {
    return ResetCustomVerifyEmailMessageTextToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomVerifyEmailMessageTextToDefaultRequest>,
  ): ResetCustomVerifyEmailMessageTextToDefaultRequest {
    const message = createBaseResetCustomVerifyEmailMessageTextToDefaultRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseResetCustomVerifyEmailMessageTextToDefaultResponse(): ResetCustomVerifyEmailMessageTextToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomVerifyEmailMessageTextToDefaultResponse = {
  encode(
    message: ResetCustomVerifyEmailMessageTextToDefaultResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomVerifyEmailMessageTextToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomVerifyEmailMessageTextToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomVerifyEmailMessageTextToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomVerifyEmailMessageTextToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomVerifyEmailMessageTextToDefaultResponse>,
  ): ResetCustomVerifyEmailMessageTextToDefaultResponse {
    return ResetCustomVerifyEmailMessageTextToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomVerifyEmailMessageTextToDefaultResponse>,
  ): ResetCustomVerifyEmailMessageTextToDefaultResponse {
    const message = createBaseResetCustomVerifyEmailMessageTextToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultVerifyPhoneMessageTextRequest(): GetDefaultVerifyPhoneMessageTextRequest {
  return { language: "" };
}

export const GetDefaultVerifyPhoneMessageTextRequest = {
  encode(message: GetDefaultVerifyPhoneMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultVerifyPhoneMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultVerifyPhoneMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultVerifyPhoneMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultVerifyPhoneMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultVerifyPhoneMessageTextRequest>): GetDefaultVerifyPhoneMessageTextRequest {
    return GetDefaultVerifyPhoneMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultVerifyPhoneMessageTextRequest>): GetDefaultVerifyPhoneMessageTextRequest {
    const message = createBaseGetDefaultVerifyPhoneMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetDefaultVerifyPhoneMessageTextResponse(): GetDefaultVerifyPhoneMessageTextResponse {
  return { customText: undefined };
}

export const GetDefaultVerifyPhoneMessageTextResponse = {
  encode(message: GetDefaultVerifyPhoneMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultVerifyPhoneMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultVerifyPhoneMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultVerifyPhoneMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetDefaultVerifyPhoneMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultVerifyPhoneMessageTextResponse>): GetDefaultVerifyPhoneMessageTextResponse {
    return GetDefaultVerifyPhoneMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultVerifyPhoneMessageTextResponse>): GetDefaultVerifyPhoneMessageTextResponse {
    const message = createBaseGetDefaultVerifyPhoneMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseGetCustomVerifyPhoneMessageTextRequest(): GetCustomVerifyPhoneMessageTextRequest {
  return { language: "" };
}

export const GetCustomVerifyPhoneMessageTextRequest = {
  encode(message: GetCustomVerifyPhoneMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomVerifyPhoneMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomVerifyPhoneMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomVerifyPhoneMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetCustomVerifyPhoneMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetCustomVerifyPhoneMessageTextRequest>): GetCustomVerifyPhoneMessageTextRequest {
    return GetCustomVerifyPhoneMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomVerifyPhoneMessageTextRequest>): GetCustomVerifyPhoneMessageTextRequest {
    const message = createBaseGetCustomVerifyPhoneMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetCustomVerifyPhoneMessageTextResponse(): GetCustomVerifyPhoneMessageTextResponse {
  return { customText: undefined };
}

export const GetCustomVerifyPhoneMessageTextResponse = {
  encode(message: GetCustomVerifyPhoneMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomVerifyPhoneMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomVerifyPhoneMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomVerifyPhoneMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetCustomVerifyPhoneMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetCustomVerifyPhoneMessageTextResponse>): GetCustomVerifyPhoneMessageTextResponse {
    return GetCustomVerifyPhoneMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomVerifyPhoneMessageTextResponse>): GetCustomVerifyPhoneMessageTextResponse {
    const message = createBaseGetCustomVerifyPhoneMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseSetDefaultVerifyPhoneMessageTextRequest(): SetDefaultVerifyPhoneMessageTextRequest {
  return {
    language: "",
    title: "",
    preHeader: "",
    subject: "",
    greeting: "",
    text: "",
    buttonText: "",
    footerText: "",
  };
}

export const SetDefaultVerifyPhoneMessageTextRequest = {
  encode(message: SetDefaultVerifyPhoneMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.preHeader !== "") {
      writer.uint32(26).string(message.preHeader);
    }
    if (message.subject !== "") {
      writer.uint32(34).string(message.subject);
    }
    if (message.greeting !== "") {
      writer.uint32(42).string(message.greeting);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.buttonText !== "") {
      writer.uint32(58).string(message.buttonText);
    }
    if (message.footerText !== "") {
      writer.uint32(66).string(message.footerText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultVerifyPhoneMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultVerifyPhoneMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.preHeader = reader.string();
          break;
        case 4:
          message.subject = reader.string();
          break;
        case 5:
          message.greeting = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.buttonText = reader.string();
          break;
        case 8:
          message.footerText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultVerifyPhoneMessageTextRequest {
    return {
      language: isSet(object.language) ? String(object.language) : "",
      title: isSet(object.title) ? String(object.title) : "",
      preHeader: isSet(object.preHeader) ? String(object.preHeader) : "",
      subject: isSet(object.subject) ? String(object.subject) : "",
      greeting: isSet(object.greeting) ? String(object.greeting) : "",
      text: isSet(object.text) ? String(object.text) : "",
      buttonText: isSet(object.buttonText) ? String(object.buttonText) : "",
      footerText: isSet(object.footerText) ? String(object.footerText) : "",
    };
  },

  toJSON(message: SetDefaultVerifyPhoneMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    message.title !== undefined && (obj.title = message.title);
    message.preHeader !== undefined && (obj.preHeader = message.preHeader);
    message.subject !== undefined && (obj.subject = message.subject);
    message.greeting !== undefined && (obj.greeting = message.greeting);
    message.text !== undefined && (obj.text = message.text);
    message.buttonText !== undefined && (obj.buttonText = message.buttonText);
    message.footerText !== undefined && (obj.footerText = message.footerText);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultVerifyPhoneMessageTextRequest>): SetDefaultVerifyPhoneMessageTextRequest {
    return SetDefaultVerifyPhoneMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultVerifyPhoneMessageTextRequest>): SetDefaultVerifyPhoneMessageTextRequest {
    const message = createBaseSetDefaultVerifyPhoneMessageTextRequest();
    message.language = object.language ?? "";
    message.title = object.title ?? "";
    message.preHeader = object.preHeader ?? "";
    message.subject = object.subject ?? "";
    message.greeting = object.greeting ?? "";
    message.text = object.text ?? "";
    message.buttonText = object.buttonText ?? "";
    message.footerText = object.footerText ?? "";
    return message;
  },
};

function createBaseSetDefaultVerifyPhoneMessageTextResponse(): SetDefaultVerifyPhoneMessageTextResponse {
  return { details: undefined };
}

export const SetDefaultVerifyPhoneMessageTextResponse = {
  encode(message: SetDefaultVerifyPhoneMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultVerifyPhoneMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultVerifyPhoneMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultVerifyPhoneMessageTextResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultVerifyPhoneMessageTextResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultVerifyPhoneMessageTextResponse>): SetDefaultVerifyPhoneMessageTextResponse {
    return SetDefaultVerifyPhoneMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetDefaultVerifyPhoneMessageTextResponse>): SetDefaultVerifyPhoneMessageTextResponse {
    const message = createBaseSetDefaultVerifyPhoneMessageTextResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomVerifyPhoneMessageTextToDefaultRequest(): ResetCustomVerifyPhoneMessageTextToDefaultRequest {
  return { language: "" };
}

export const ResetCustomVerifyPhoneMessageTextToDefaultRequest = {
  encode(
    message: ResetCustomVerifyPhoneMessageTextToDefaultRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomVerifyPhoneMessageTextToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomVerifyPhoneMessageTextToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomVerifyPhoneMessageTextToDefaultRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: ResetCustomVerifyPhoneMessageTextToDefaultRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomVerifyPhoneMessageTextToDefaultRequest>,
  ): ResetCustomVerifyPhoneMessageTextToDefaultRequest {
    return ResetCustomVerifyPhoneMessageTextToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomVerifyPhoneMessageTextToDefaultRequest>,
  ): ResetCustomVerifyPhoneMessageTextToDefaultRequest {
    const message = createBaseResetCustomVerifyPhoneMessageTextToDefaultRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseResetCustomVerifyPhoneMessageTextToDefaultResponse(): ResetCustomVerifyPhoneMessageTextToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomVerifyPhoneMessageTextToDefaultResponse = {
  encode(
    message: ResetCustomVerifyPhoneMessageTextToDefaultResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomVerifyPhoneMessageTextToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomVerifyPhoneMessageTextToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomVerifyPhoneMessageTextToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomVerifyPhoneMessageTextToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomVerifyPhoneMessageTextToDefaultResponse>,
  ): ResetCustomVerifyPhoneMessageTextToDefaultResponse {
    return ResetCustomVerifyPhoneMessageTextToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomVerifyPhoneMessageTextToDefaultResponse>,
  ): ResetCustomVerifyPhoneMessageTextToDefaultResponse {
    const message = createBaseResetCustomVerifyPhoneMessageTextToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultDomainClaimedMessageTextRequest(): GetDefaultDomainClaimedMessageTextRequest {
  return { language: "" };
}

export const GetDefaultDomainClaimedMessageTextRequest = {
  encode(message: GetDefaultDomainClaimedMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultDomainClaimedMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultDomainClaimedMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultDomainClaimedMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultDomainClaimedMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultDomainClaimedMessageTextRequest>): GetDefaultDomainClaimedMessageTextRequest {
    return GetDefaultDomainClaimedMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultDomainClaimedMessageTextRequest>,
  ): GetDefaultDomainClaimedMessageTextRequest {
    const message = createBaseGetDefaultDomainClaimedMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetDefaultDomainClaimedMessageTextResponse(): GetDefaultDomainClaimedMessageTextResponse {
  return { customText: undefined };
}

export const GetDefaultDomainClaimedMessageTextResponse = {
  encode(message: GetDefaultDomainClaimedMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultDomainClaimedMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultDomainClaimedMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultDomainClaimedMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetDefaultDomainClaimedMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultDomainClaimedMessageTextResponse>): GetDefaultDomainClaimedMessageTextResponse {
    return GetDefaultDomainClaimedMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultDomainClaimedMessageTextResponse>,
  ): GetDefaultDomainClaimedMessageTextResponse {
    const message = createBaseGetDefaultDomainClaimedMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseGetCustomDomainClaimedMessageTextRequest(): GetCustomDomainClaimedMessageTextRequest {
  return { language: "" };
}

export const GetCustomDomainClaimedMessageTextRequest = {
  encode(message: GetCustomDomainClaimedMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomDomainClaimedMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomDomainClaimedMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomDomainClaimedMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetCustomDomainClaimedMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetCustomDomainClaimedMessageTextRequest>): GetCustomDomainClaimedMessageTextRequest {
    return GetCustomDomainClaimedMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomDomainClaimedMessageTextRequest>): GetCustomDomainClaimedMessageTextRequest {
    const message = createBaseGetCustomDomainClaimedMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetCustomDomainClaimedMessageTextResponse(): GetCustomDomainClaimedMessageTextResponse {
  return { customText: undefined };
}

export const GetCustomDomainClaimedMessageTextResponse = {
  encode(message: GetCustomDomainClaimedMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomDomainClaimedMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomDomainClaimedMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomDomainClaimedMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetCustomDomainClaimedMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetCustomDomainClaimedMessageTextResponse>): GetCustomDomainClaimedMessageTextResponse {
    return GetCustomDomainClaimedMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetCustomDomainClaimedMessageTextResponse>,
  ): GetCustomDomainClaimedMessageTextResponse {
    const message = createBaseGetCustomDomainClaimedMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseSetDefaultDomainClaimedMessageTextRequest(): SetDefaultDomainClaimedMessageTextRequest {
  return {
    language: "",
    title: "",
    preHeader: "",
    subject: "",
    greeting: "",
    text: "",
    buttonText: "",
    footerText: "",
  };
}

export const SetDefaultDomainClaimedMessageTextRequest = {
  encode(message: SetDefaultDomainClaimedMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.preHeader !== "") {
      writer.uint32(26).string(message.preHeader);
    }
    if (message.subject !== "") {
      writer.uint32(34).string(message.subject);
    }
    if (message.greeting !== "") {
      writer.uint32(42).string(message.greeting);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.buttonText !== "") {
      writer.uint32(58).string(message.buttonText);
    }
    if (message.footerText !== "") {
      writer.uint32(66).string(message.footerText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultDomainClaimedMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultDomainClaimedMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.preHeader = reader.string();
          break;
        case 4:
          message.subject = reader.string();
          break;
        case 5:
          message.greeting = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.buttonText = reader.string();
          break;
        case 8:
          message.footerText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultDomainClaimedMessageTextRequest {
    return {
      language: isSet(object.language) ? String(object.language) : "",
      title: isSet(object.title) ? String(object.title) : "",
      preHeader: isSet(object.preHeader) ? String(object.preHeader) : "",
      subject: isSet(object.subject) ? String(object.subject) : "",
      greeting: isSet(object.greeting) ? String(object.greeting) : "",
      text: isSet(object.text) ? String(object.text) : "",
      buttonText: isSet(object.buttonText) ? String(object.buttonText) : "",
      footerText: isSet(object.footerText) ? String(object.footerText) : "",
    };
  },

  toJSON(message: SetDefaultDomainClaimedMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    message.title !== undefined && (obj.title = message.title);
    message.preHeader !== undefined && (obj.preHeader = message.preHeader);
    message.subject !== undefined && (obj.subject = message.subject);
    message.greeting !== undefined && (obj.greeting = message.greeting);
    message.text !== undefined && (obj.text = message.text);
    message.buttonText !== undefined && (obj.buttonText = message.buttonText);
    message.footerText !== undefined && (obj.footerText = message.footerText);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultDomainClaimedMessageTextRequest>): SetDefaultDomainClaimedMessageTextRequest {
    return SetDefaultDomainClaimedMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SetDefaultDomainClaimedMessageTextRequest>,
  ): SetDefaultDomainClaimedMessageTextRequest {
    const message = createBaseSetDefaultDomainClaimedMessageTextRequest();
    message.language = object.language ?? "";
    message.title = object.title ?? "";
    message.preHeader = object.preHeader ?? "";
    message.subject = object.subject ?? "";
    message.greeting = object.greeting ?? "";
    message.text = object.text ?? "";
    message.buttonText = object.buttonText ?? "";
    message.footerText = object.footerText ?? "";
    return message;
  },
};

function createBaseSetDefaultDomainClaimedMessageTextResponse(): SetDefaultDomainClaimedMessageTextResponse {
  return { details: undefined };
}

export const SetDefaultDomainClaimedMessageTextResponse = {
  encode(message: SetDefaultDomainClaimedMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultDomainClaimedMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultDomainClaimedMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultDomainClaimedMessageTextResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultDomainClaimedMessageTextResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultDomainClaimedMessageTextResponse>): SetDefaultDomainClaimedMessageTextResponse {
    return SetDefaultDomainClaimedMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SetDefaultDomainClaimedMessageTextResponse>,
  ): SetDefaultDomainClaimedMessageTextResponse {
    const message = createBaseSetDefaultDomainClaimedMessageTextResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomDomainClaimedMessageTextToDefaultRequest(): ResetCustomDomainClaimedMessageTextToDefaultRequest {
  return { language: "" };
}

export const ResetCustomDomainClaimedMessageTextToDefaultRequest = {
  encode(
    message: ResetCustomDomainClaimedMessageTextToDefaultRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomDomainClaimedMessageTextToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomDomainClaimedMessageTextToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomDomainClaimedMessageTextToDefaultRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: ResetCustomDomainClaimedMessageTextToDefaultRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomDomainClaimedMessageTextToDefaultRequest>,
  ): ResetCustomDomainClaimedMessageTextToDefaultRequest {
    return ResetCustomDomainClaimedMessageTextToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomDomainClaimedMessageTextToDefaultRequest>,
  ): ResetCustomDomainClaimedMessageTextToDefaultRequest {
    const message = createBaseResetCustomDomainClaimedMessageTextToDefaultRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseResetCustomDomainClaimedMessageTextToDefaultResponse(): ResetCustomDomainClaimedMessageTextToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomDomainClaimedMessageTextToDefaultResponse = {
  encode(
    message: ResetCustomDomainClaimedMessageTextToDefaultResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomDomainClaimedMessageTextToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomDomainClaimedMessageTextToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomDomainClaimedMessageTextToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomDomainClaimedMessageTextToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomDomainClaimedMessageTextToDefaultResponse>,
  ): ResetCustomDomainClaimedMessageTextToDefaultResponse {
    return ResetCustomDomainClaimedMessageTextToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomDomainClaimedMessageTextToDefaultResponse>,
  ): ResetCustomDomainClaimedMessageTextToDefaultResponse {
    const message = createBaseResetCustomDomainClaimedMessageTextToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultPasswordChangeMessageTextRequest(): GetDefaultPasswordChangeMessageTextRequest {
  return { language: "" };
}

export const GetDefaultPasswordChangeMessageTextRequest = {
  encode(message: GetDefaultPasswordChangeMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultPasswordChangeMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultPasswordChangeMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultPasswordChangeMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultPasswordChangeMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultPasswordChangeMessageTextRequest>): GetDefaultPasswordChangeMessageTextRequest {
    return GetDefaultPasswordChangeMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultPasswordChangeMessageTextRequest>,
  ): GetDefaultPasswordChangeMessageTextRequest {
    const message = createBaseGetDefaultPasswordChangeMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetDefaultPasswordChangeMessageTextResponse(): GetDefaultPasswordChangeMessageTextResponse {
  return { customText: undefined };
}

export const GetDefaultPasswordChangeMessageTextResponse = {
  encode(message: GetDefaultPasswordChangeMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultPasswordChangeMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultPasswordChangeMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultPasswordChangeMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetDefaultPasswordChangeMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultPasswordChangeMessageTextResponse>): GetDefaultPasswordChangeMessageTextResponse {
    return GetDefaultPasswordChangeMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultPasswordChangeMessageTextResponse>,
  ): GetDefaultPasswordChangeMessageTextResponse {
    const message = createBaseGetDefaultPasswordChangeMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseGetCustomPasswordChangeMessageTextRequest(): GetCustomPasswordChangeMessageTextRequest {
  return { language: "" };
}

export const GetCustomPasswordChangeMessageTextRequest = {
  encode(message: GetCustomPasswordChangeMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomPasswordChangeMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomPasswordChangeMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomPasswordChangeMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetCustomPasswordChangeMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetCustomPasswordChangeMessageTextRequest>): GetCustomPasswordChangeMessageTextRequest {
    return GetCustomPasswordChangeMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetCustomPasswordChangeMessageTextRequest>,
  ): GetCustomPasswordChangeMessageTextRequest {
    const message = createBaseGetCustomPasswordChangeMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetCustomPasswordChangeMessageTextResponse(): GetCustomPasswordChangeMessageTextResponse {
  return { customText: undefined };
}

export const GetCustomPasswordChangeMessageTextResponse = {
  encode(message: GetCustomPasswordChangeMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomPasswordChangeMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomPasswordChangeMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomPasswordChangeMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetCustomPasswordChangeMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetCustomPasswordChangeMessageTextResponse>): GetCustomPasswordChangeMessageTextResponse {
    return GetCustomPasswordChangeMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetCustomPasswordChangeMessageTextResponse>,
  ): GetCustomPasswordChangeMessageTextResponse {
    const message = createBaseGetCustomPasswordChangeMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseSetDefaultPasswordChangeMessageTextRequest(): SetDefaultPasswordChangeMessageTextRequest {
  return {
    language: "",
    title: "",
    preHeader: "",
    subject: "",
    greeting: "",
    text: "",
    buttonText: "",
    footerText: "",
  };
}

export const SetDefaultPasswordChangeMessageTextRequest = {
  encode(message: SetDefaultPasswordChangeMessageTextRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.preHeader !== "") {
      writer.uint32(26).string(message.preHeader);
    }
    if (message.subject !== "") {
      writer.uint32(34).string(message.subject);
    }
    if (message.greeting !== "") {
      writer.uint32(42).string(message.greeting);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.buttonText !== "") {
      writer.uint32(58).string(message.buttonText);
    }
    if (message.footerText !== "") {
      writer.uint32(66).string(message.footerText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultPasswordChangeMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultPasswordChangeMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.preHeader = reader.string();
          break;
        case 4:
          message.subject = reader.string();
          break;
        case 5:
          message.greeting = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.buttonText = reader.string();
          break;
        case 8:
          message.footerText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultPasswordChangeMessageTextRequest {
    return {
      language: isSet(object.language) ? String(object.language) : "",
      title: isSet(object.title) ? String(object.title) : "",
      preHeader: isSet(object.preHeader) ? String(object.preHeader) : "",
      subject: isSet(object.subject) ? String(object.subject) : "",
      greeting: isSet(object.greeting) ? String(object.greeting) : "",
      text: isSet(object.text) ? String(object.text) : "",
      buttonText: isSet(object.buttonText) ? String(object.buttonText) : "",
      footerText: isSet(object.footerText) ? String(object.footerText) : "",
    };
  },

  toJSON(message: SetDefaultPasswordChangeMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    message.title !== undefined && (obj.title = message.title);
    message.preHeader !== undefined && (obj.preHeader = message.preHeader);
    message.subject !== undefined && (obj.subject = message.subject);
    message.greeting !== undefined && (obj.greeting = message.greeting);
    message.text !== undefined && (obj.text = message.text);
    message.buttonText !== undefined && (obj.buttonText = message.buttonText);
    message.footerText !== undefined && (obj.footerText = message.footerText);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultPasswordChangeMessageTextRequest>): SetDefaultPasswordChangeMessageTextRequest {
    return SetDefaultPasswordChangeMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SetDefaultPasswordChangeMessageTextRequest>,
  ): SetDefaultPasswordChangeMessageTextRequest {
    const message = createBaseSetDefaultPasswordChangeMessageTextRequest();
    message.language = object.language ?? "";
    message.title = object.title ?? "";
    message.preHeader = object.preHeader ?? "";
    message.subject = object.subject ?? "";
    message.greeting = object.greeting ?? "";
    message.text = object.text ?? "";
    message.buttonText = object.buttonText ?? "";
    message.footerText = object.footerText ?? "";
    return message;
  },
};

function createBaseSetDefaultPasswordChangeMessageTextResponse(): SetDefaultPasswordChangeMessageTextResponse {
  return { details: undefined };
}

export const SetDefaultPasswordChangeMessageTextResponse = {
  encode(message: SetDefaultPasswordChangeMessageTextResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultPasswordChangeMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultPasswordChangeMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultPasswordChangeMessageTextResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultPasswordChangeMessageTextResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetDefaultPasswordChangeMessageTextResponse>): SetDefaultPasswordChangeMessageTextResponse {
    return SetDefaultPasswordChangeMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SetDefaultPasswordChangeMessageTextResponse>,
  ): SetDefaultPasswordChangeMessageTextResponse {
    const message = createBaseSetDefaultPasswordChangeMessageTextResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomPasswordChangeMessageTextToDefaultRequest(): ResetCustomPasswordChangeMessageTextToDefaultRequest {
  return { language: "" };
}

export const ResetCustomPasswordChangeMessageTextToDefaultRequest = {
  encode(
    message: ResetCustomPasswordChangeMessageTextToDefaultRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomPasswordChangeMessageTextToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomPasswordChangeMessageTextToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomPasswordChangeMessageTextToDefaultRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: ResetCustomPasswordChangeMessageTextToDefaultRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomPasswordChangeMessageTextToDefaultRequest>,
  ): ResetCustomPasswordChangeMessageTextToDefaultRequest {
    return ResetCustomPasswordChangeMessageTextToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomPasswordChangeMessageTextToDefaultRequest>,
  ): ResetCustomPasswordChangeMessageTextToDefaultRequest {
    const message = createBaseResetCustomPasswordChangeMessageTextToDefaultRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseResetCustomPasswordChangeMessageTextToDefaultResponse(): ResetCustomPasswordChangeMessageTextToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomPasswordChangeMessageTextToDefaultResponse = {
  encode(
    message: ResetCustomPasswordChangeMessageTextToDefaultResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomPasswordChangeMessageTextToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomPasswordChangeMessageTextToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomPasswordChangeMessageTextToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomPasswordChangeMessageTextToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomPasswordChangeMessageTextToDefaultResponse>,
  ): ResetCustomPasswordChangeMessageTextToDefaultResponse {
    return ResetCustomPasswordChangeMessageTextToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomPasswordChangeMessageTextToDefaultResponse>,
  ): ResetCustomPasswordChangeMessageTextToDefaultResponse {
    const message = createBaseResetCustomPasswordChangeMessageTextToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultPasswordlessRegistrationMessageTextRequest(): GetDefaultPasswordlessRegistrationMessageTextRequest {
  return { language: "" };
}

export const GetDefaultPasswordlessRegistrationMessageTextRequest = {
  encode(
    message: GetDefaultPasswordlessRegistrationMessageTextRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultPasswordlessRegistrationMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultPasswordlessRegistrationMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultPasswordlessRegistrationMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultPasswordlessRegistrationMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(
    base?: DeepPartial<GetDefaultPasswordlessRegistrationMessageTextRequest>,
  ): GetDefaultPasswordlessRegistrationMessageTextRequest {
    return GetDefaultPasswordlessRegistrationMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultPasswordlessRegistrationMessageTextRequest>,
  ): GetDefaultPasswordlessRegistrationMessageTextRequest {
    const message = createBaseGetDefaultPasswordlessRegistrationMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetDefaultPasswordlessRegistrationMessageTextResponse(): GetDefaultPasswordlessRegistrationMessageTextResponse {
  return { customText: undefined };
}

export const GetDefaultPasswordlessRegistrationMessageTextResponse = {
  encode(
    message: GetDefaultPasswordlessRegistrationMessageTextResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultPasswordlessRegistrationMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultPasswordlessRegistrationMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultPasswordlessRegistrationMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetDefaultPasswordlessRegistrationMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<GetDefaultPasswordlessRegistrationMessageTextResponse>,
  ): GetDefaultPasswordlessRegistrationMessageTextResponse {
    return GetDefaultPasswordlessRegistrationMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultPasswordlessRegistrationMessageTextResponse>,
  ): GetDefaultPasswordlessRegistrationMessageTextResponse {
    const message = createBaseGetDefaultPasswordlessRegistrationMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseGetCustomPasswordlessRegistrationMessageTextRequest(): GetCustomPasswordlessRegistrationMessageTextRequest {
  return { language: "" };
}

export const GetCustomPasswordlessRegistrationMessageTextRequest = {
  encode(
    message: GetCustomPasswordlessRegistrationMessageTextRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomPasswordlessRegistrationMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomPasswordlessRegistrationMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomPasswordlessRegistrationMessageTextRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetCustomPasswordlessRegistrationMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(
    base?: DeepPartial<GetCustomPasswordlessRegistrationMessageTextRequest>,
  ): GetCustomPasswordlessRegistrationMessageTextRequest {
    return GetCustomPasswordlessRegistrationMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetCustomPasswordlessRegistrationMessageTextRequest>,
  ): GetCustomPasswordlessRegistrationMessageTextRequest {
    const message = createBaseGetCustomPasswordlessRegistrationMessageTextRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetCustomPasswordlessRegistrationMessageTextResponse(): GetCustomPasswordlessRegistrationMessageTextResponse {
  return { customText: undefined };
}

export const GetCustomPasswordlessRegistrationMessageTextResponse = {
  encode(
    message: GetCustomPasswordlessRegistrationMessageTextResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.customText !== undefined) {
      MessageCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomPasswordlessRegistrationMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomPasswordlessRegistrationMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = MessageCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomPasswordlessRegistrationMessageTextResponse {
    return { customText: isSet(object.customText) ? MessageCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetCustomPasswordlessRegistrationMessageTextResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? MessageCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<GetCustomPasswordlessRegistrationMessageTextResponse>,
  ): GetCustomPasswordlessRegistrationMessageTextResponse {
    return GetCustomPasswordlessRegistrationMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetCustomPasswordlessRegistrationMessageTextResponse>,
  ): GetCustomPasswordlessRegistrationMessageTextResponse {
    const message = createBaseGetCustomPasswordlessRegistrationMessageTextResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? MessageCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseSetDefaultPasswordlessRegistrationMessageTextRequest(): SetDefaultPasswordlessRegistrationMessageTextRequest {
  return {
    language: "",
    title: "",
    preHeader: "",
    subject: "",
    greeting: "",
    text: "",
    buttonText: "",
    footerText: "",
  };
}

export const SetDefaultPasswordlessRegistrationMessageTextRequest = {
  encode(
    message: SetDefaultPasswordlessRegistrationMessageTextRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.preHeader !== "") {
      writer.uint32(26).string(message.preHeader);
    }
    if (message.subject !== "") {
      writer.uint32(34).string(message.subject);
    }
    if (message.greeting !== "") {
      writer.uint32(42).string(message.greeting);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.buttonText !== "") {
      writer.uint32(58).string(message.buttonText);
    }
    if (message.footerText !== "") {
      writer.uint32(66).string(message.footerText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultPasswordlessRegistrationMessageTextRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultPasswordlessRegistrationMessageTextRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.preHeader = reader.string();
          break;
        case 4:
          message.subject = reader.string();
          break;
        case 5:
          message.greeting = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.buttonText = reader.string();
          break;
        case 8:
          message.footerText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultPasswordlessRegistrationMessageTextRequest {
    return {
      language: isSet(object.language) ? String(object.language) : "",
      title: isSet(object.title) ? String(object.title) : "",
      preHeader: isSet(object.preHeader) ? String(object.preHeader) : "",
      subject: isSet(object.subject) ? String(object.subject) : "",
      greeting: isSet(object.greeting) ? String(object.greeting) : "",
      text: isSet(object.text) ? String(object.text) : "",
      buttonText: isSet(object.buttonText) ? String(object.buttonText) : "",
      footerText: isSet(object.footerText) ? String(object.footerText) : "",
    };
  },

  toJSON(message: SetDefaultPasswordlessRegistrationMessageTextRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    message.title !== undefined && (obj.title = message.title);
    message.preHeader !== undefined && (obj.preHeader = message.preHeader);
    message.subject !== undefined && (obj.subject = message.subject);
    message.greeting !== undefined && (obj.greeting = message.greeting);
    message.text !== undefined && (obj.text = message.text);
    message.buttonText !== undefined && (obj.buttonText = message.buttonText);
    message.footerText !== undefined && (obj.footerText = message.footerText);
    return obj;
  },

  create(
    base?: DeepPartial<SetDefaultPasswordlessRegistrationMessageTextRequest>,
  ): SetDefaultPasswordlessRegistrationMessageTextRequest {
    return SetDefaultPasswordlessRegistrationMessageTextRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SetDefaultPasswordlessRegistrationMessageTextRequest>,
  ): SetDefaultPasswordlessRegistrationMessageTextRequest {
    const message = createBaseSetDefaultPasswordlessRegistrationMessageTextRequest();
    message.language = object.language ?? "";
    message.title = object.title ?? "";
    message.preHeader = object.preHeader ?? "";
    message.subject = object.subject ?? "";
    message.greeting = object.greeting ?? "";
    message.text = object.text ?? "";
    message.buttonText = object.buttonText ?? "";
    message.footerText = object.footerText ?? "";
    return message;
  },
};

function createBaseSetDefaultPasswordlessRegistrationMessageTextResponse(): SetDefaultPasswordlessRegistrationMessageTextResponse {
  return { details: undefined };
}

export const SetDefaultPasswordlessRegistrationMessageTextResponse = {
  encode(
    message: SetDefaultPasswordlessRegistrationMessageTextResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetDefaultPasswordlessRegistrationMessageTextResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetDefaultPasswordlessRegistrationMessageTextResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetDefaultPasswordlessRegistrationMessageTextResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetDefaultPasswordlessRegistrationMessageTextResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<SetDefaultPasswordlessRegistrationMessageTextResponse>,
  ): SetDefaultPasswordlessRegistrationMessageTextResponse {
    return SetDefaultPasswordlessRegistrationMessageTextResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<SetDefaultPasswordlessRegistrationMessageTextResponse>,
  ): SetDefaultPasswordlessRegistrationMessageTextResponse {
    const message = createBaseSetDefaultPasswordlessRegistrationMessageTextResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomPasswordlessRegistrationMessageTextToDefaultRequest(): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest {
  return { language: "" };
}

export const ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest = {
  encode(
    message: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number,
  ): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomPasswordlessRegistrationMessageTextToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest>,
  ): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest {
    return ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest>,
  ): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest {
    const message = createBaseResetCustomPasswordlessRegistrationMessageTextToDefaultRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseResetCustomPasswordlessRegistrationMessageTextToDefaultResponse(): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse = {
  encode(
    message: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number,
  ): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomPasswordlessRegistrationMessageTextToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(
    base?: DeepPartial<ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse>,
  ): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse {
    return ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse>,
  ): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse {
    const message = createBaseResetCustomPasswordlessRegistrationMessageTextToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetDefaultLoginTextsRequest(): GetDefaultLoginTextsRequest {
  return { language: "" };
}

export const GetDefaultLoginTextsRequest = {
  encode(message: GetDefaultLoginTextsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultLoginTextsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultLoginTextsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultLoginTextsRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetDefaultLoginTextsRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultLoginTextsRequest>): GetDefaultLoginTextsRequest {
    return GetDefaultLoginTextsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultLoginTextsRequest>): GetDefaultLoginTextsRequest {
    const message = createBaseGetDefaultLoginTextsRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetDefaultLoginTextsResponse(): GetDefaultLoginTextsResponse {
  return { customText: undefined };
}

export const GetDefaultLoginTextsResponse = {
  encode(message: GetDefaultLoginTextsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      LoginCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultLoginTextsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultLoginTextsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = LoginCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultLoginTextsResponse {
    return { customText: isSet(object.customText) ? LoginCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetDefaultLoginTextsResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? LoginCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetDefaultLoginTextsResponse>): GetDefaultLoginTextsResponse {
    return GetDefaultLoginTextsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultLoginTextsResponse>): GetDefaultLoginTextsResponse {
    const message = createBaseGetDefaultLoginTextsResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? LoginCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseGetCustomLoginTextsRequest(): GetCustomLoginTextsRequest {
  return { language: "" };
}

export const GetCustomLoginTextsRequest = {
  encode(message: GetCustomLoginTextsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomLoginTextsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomLoginTextsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomLoginTextsRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: GetCustomLoginTextsRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<GetCustomLoginTextsRequest>): GetCustomLoginTextsRequest {
    return GetCustomLoginTextsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomLoginTextsRequest>): GetCustomLoginTextsRequest {
    const message = createBaseGetCustomLoginTextsRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseGetCustomLoginTextsResponse(): GetCustomLoginTextsResponse {
  return { customText: undefined };
}

export const GetCustomLoginTextsResponse = {
  encode(message: GetCustomLoginTextsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.customText !== undefined) {
      LoginCustomText.encode(message.customText, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCustomLoginTextsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCustomLoginTextsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.customText = LoginCustomText.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetCustomLoginTextsResponse {
    return { customText: isSet(object.customText) ? LoginCustomText.fromJSON(object.customText) : undefined };
  },

  toJSON(message: GetCustomLoginTextsResponse): unknown {
    const obj: any = {};
    message.customText !== undefined &&
      (obj.customText = message.customText ? LoginCustomText.toJSON(message.customText) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetCustomLoginTextsResponse>): GetCustomLoginTextsResponse {
    return GetCustomLoginTextsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetCustomLoginTextsResponse>): GetCustomLoginTextsResponse {
    const message = createBaseGetCustomLoginTextsResponse();
    message.customText = (object.customText !== undefined && object.customText !== null)
      ? LoginCustomText.fromPartial(object.customText)
      : undefined;
    return message;
  },
};

function createBaseSetCustomLoginTextsRequest(): SetCustomLoginTextsRequest {
  return {
    language: "",
    selectAccountText: undefined,
    loginText: undefined,
    passwordText: undefined,
    usernameChangeText: undefined,
    usernameChangeDoneText: undefined,
    initPasswordText: undefined,
    initPasswordDoneText: undefined,
    emailVerificationText: undefined,
    emailVerificationDoneText: undefined,
    initializeUserText: undefined,
    initializeDoneText: undefined,
    initMfaPromptText: undefined,
    initMfaOtpText: undefined,
    initMfaU2fText: undefined,
    initMfaDoneText: undefined,
    mfaProvidersText: undefined,
    verifyMfaOtpText: undefined,
    verifyMfaU2fText: undefined,
    passwordlessText: undefined,
    passwordChangeText: undefined,
    passwordChangeDoneText: undefined,
    passwordResetDoneText: undefined,
    registrationOptionText: undefined,
    registrationUserText: undefined,
    registrationOrgText: undefined,
    linkingUserDoneText: undefined,
    externalUserNotFoundText: undefined,
    successLoginText: undefined,
    logoutText: undefined,
    footerText: undefined,
    passwordlessPromptText: undefined,
    passwordlessRegistrationText: undefined,
    passwordlessRegistrationDoneText: undefined,
    externalRegistrationUserOverviewText: undefined,
  };
}

export const SetCustomLoginTextsRequest = {
  encode(message: SetCustomLoginTextsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    if (message.selectAccountText !== undefined) {
      SelectAccountScreenText.encode(message.selectAccountText, writer.uint32(18).fork()).ldelim();
    }
    if (message.loginText !== undefined) {
      LoginScreenText.encode(message.loginText, writer.uint32(26).fork()).ldelim();
    }
    if (message.passwordText !== undefined) {
      PasswordScreenText.encode(message.passwordText, writer.uint32(34).fork()).ldelim();
    }
    if (message.usernameChangeText !== undefined) {
      UsernameChangeScreenText.encode(message.usernameChangeText, writer.uint32(42).fork()).ldelim();
    }
    if (message.usernameChangeDoneText !== undefined) {
      UsernameChangeDoneScreenText.encode(message.usernameChangeDoneText, writer.uint32(50).fork()).ldelim();
    }
    if (message.initPasswordText !== undefined) {
      InitPasswordScreenText.encode(message.initPasswordText, writer.uint32(58).fork()).ldelim();
    }
    if (message.initPasswordDoneText !== undefined) {
      InitPasswordDoneScreenText.encode(message.initPasswordDoneText, writer.uint32(66).fork()).ldelim();
    }
    if (message.emailVerificationText !== undefined) {
      EmailVerificationScreenText.encode(message.emailVerificationText, writer.uint32(74).fork()).ldelim();
    }
    if (message.emailVerificationDoneText !== undefined) {
      EmailVerificationDoneScreenText.encode(message.emailVerificationDoneText, writer.uint32(82).fork()).ldelim();
    }
    if (message.initializeUserText !== undefined) {
      InitializeUserScreenText.encode(message.initializeUserText, writer.uint32(90).fork()).ldelim();
    }
    if (message.initializeDoneText !== undefined) {
      InitializeUserDoneScreenText.encode(message.initializeDoneText, writer.uint32(98).fork()).ldelim();
    }
    if (message.initMfaPromptText !== undefined) {
      InitMFAPromptScreenText.encode(message.initMfaPromptText, writer.uint32(106).fork()).ldelim();
    }
    if (message.initMfaOtpText !== undefined) {
      InitMFAOTPScreenText.encode(message.initMfaOtpText, writer.uint32(114).fork()).ldelim();
    }
    if (message.initMfaU2fText !== undefined) {
      InitMFAU2FScreenText.encode(message.initMfaU2fText, writer.uint32(122).fork()).ldelim();
    }
    if (message.initMfaDoneText !== undefined) {
      InitMFADoneScreenText.encode(message.initMfaDoneText, writer.uint32(130).fork()).ldelim();
    }
    if (message.mfaProvidersText !== undefined) {
      MFAProvidersText.encode(message.mfaProvidersText, writer.uint32(138).fork()).ldelim();
    }
    if (message.verifyMfaOtpText !== undefined) {
      VerifyMFAOTPScreenText.encode(message.verifyMfaOtpText, writer.uint32(146).fork()).ldelim();
    }
    if (message.verifyMfaU2fText !== undefined) {
      VerifyMFAU2FScreenText.encode(message.verifyMfaU2fText, writer.uint32(154).fork()).ldelim();
    }
    if (message.passwordlessText !== undefined) {
      PasswordlessScreenText.encode(message.passwordlessText, writer.uint32(162).fork()).ldelim();
    }
    if (message.passwordChangeText !== undefined) {
      PasswordChangeScreenText.encode(message.passwordChangeText, writer.uint32(170).fork()).ldelim();
    }
    if (message.passwordChangeDoneText !== undefined) {
      PasswordChangeDoneScreenText.encode(message.passwordChangeDoneText, writer.uint32(178).fork()).ldelim();
    }
    if (message.passwordResetDoneText !== undefined) {
      PasswordResetDoneScreenText.encode(message.passwordResetDoneText, writer.uint32(186).fork()).ldelim();
    }
    if (message.registrationOptionText !== undefined) {
      RegistrationOptionScreenText.encode(message.registrationOptionText, writer.uint32(194).fork()).ldelim();
    }
    if (message.registrationUserText !== undefined) {
      RegistrationUserScreenText.encode(message.registrationUserText, writer.uint32(202).fork()).ldelim();
    }
    if (message.registrationOrgText !== undefined) {
      RegistrationOrgScreenText.encode(message.registrationOrgText, writer.uint32(210).fork()).ldelim();
    }
    if (message.linkingUserDoneText !== undefined) {
      LinkingUserDoneScreenText.encode(message.linkingUserDoneText, writer.uint32(218).fork()).ldelim();
    }
    if (message.externalUserNotFoundText !== undefined) {
      ExternalUserNotFoundScreenText.encode(message.externalUserNotFoundText, writer.uint32(226).fork()).ldelim();
    }
    if (message.successLoginText !== undefined) {
      SuccessLoginScreenText.encode(message.successLoginText, writer.uint32(234).fork()).ldelim();
    }
    if (message.logoutText !== undefined) {
      LogoutDoneScreenText.encode(message.logoutText, writer.uint32(242).fork()).ldelim();
    }
    if (message.footerText !== undefined) {
      FooterText.encode(message.footerText, writer.uint32(250).fork()).ldelim();
    }
    if (message.passwordlessPromptText !== undefined) {
      PasswordlessPromptScreenText.encode(message.passwordlessPromptText, writer.uint32(258).fork()).ldelim();
    }
    if (message.passwordlessRegistrationText !== undefined) {
      PasswordlessRegistrationScreenText.encode(message.passwordlessRegistrationText, writer.uint32(266).fork())
        .ldelim();
    }
    if (message.passwordlessRegistrationDoneText !== undefined) {
      PasswordlessRegistrationDoneScreenText.encode(message.passwordlessRegistrationDoneText, writer.uint32(274).fork())
        .ldelim();
    }
    if (message.externalRegistrationUserOverviewText !== undefined) {
      ExternalRegistrationUserOverviewScreenText.encode(
        message.externalRegistrationUserOverviewText,
        writer.uint32(282).fork(),
      ).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetCustomLoginTextsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetCustomLoginTextsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        case 2:
          message.selectAccountText = SelectAccountScreenText.decode(reader, reader.uint32());
          break;
        case 3:
          message.loginText = LoginScreenText.decode(reader, reader.uint32());
          break;
        case 4:
          message.passwordText = PasswordScreenText.decode(reader, reader.uint32());
          break;
        case 5:
          message.usernameChangeText = UsernameChangeScreenText.decode(reader, reader.uint32());
          break;
        case 6:
          message.usernameChangeDoneText = UsernameChangeDoneScreenText.decode(reader, reader.uint32());
          break;
        case 7:
          message.initPasswordText = InitPasswordScreenText.decode(reader, reader.uint32());
          break;
        case 8:
          message.initPasswordDoneText = InitPasswordDoneScreenText.decode(reader, reader.uint32());
          break;
        case 9:
          message.emailVerificationText = EmailVerificationScreenText.decode(reader, reader.uint32());
          break;
        case 10:
          message.emailVerificationDoneText = EmailVerificationDoneScreenText.decode(reader, reader.uint32());
          break;
        case 11:
          message.initializeUserText = InitializeUserScreenText.decode(reader, reader.uint32());
          break;
        case 12:
          message.initializeDoneText = InitializeUserDoneScreenText.decode(reader, reader.uint32());
          break;
        case 13:
          message.initMfaPromptText = InitMFAPromptScreenText.decode(reader, reader.uint32());
          break;
        case 14:
          message.initMfaOtpText = InitMFAOTPScreenText.decode(reader, reader.uint32());
          break;
        case 15:
          message.initMfaU2fText = InitMFAU2FScreenText.decode(reader, reader.uint32());
          break;
        case 16:
          message.initMfaDoneText = InitMFADoneScreenText.decode(reader, reader.uint32());
          break;
        case 17:
          message.mfaProvidersText = MFAProvidersText.decode(reader, reader.uint32());
          break;
        case 18:
          message.verifyMfaOtpText = VerifyMFAOTPScreenText.decode(reader, reader.uint32());
          break;
        case 19:
          message.verifyMfaU2fText = VerifyMFAU2FScreenText.decode(reader, reader.uint32());
          break;
        case 20:
          message.passwordlessText = PasswordlessScreenText.decode(reader, reader.uint32());
          break;
        case 21:
          message.passwordChangeText = PasswordChangeScreenText.decode(reader, reader.uint32());
          break;
        case 22:
          message.passwordChangeDoneText = PasswordChangeDoneScreenText.decode(reader, reader.uint32());
          break;
        case 23:
          message.passwordResetDoneText = PasswordResetDoneScreenText.decode(reader, reader.uint32());
          break;
        case 24:
          message.registrationOptionText = RegistrationOptionScreenText.decode(reader, reader.uint32());
          break;
        case 25:
          message.registrationUserText = RegistrationUserScreenText.decode(reader, reader.uint32());
          break;
        case 26:
          message.registrationOrgText = RegistrationOrgScreenText.decode(reader, reader.uint32());
          break;
        case 27:
          message.linkingUserDoneText = LinkingUserDoneScreenText.decode(reader, reader.uint32());
          break;
        case 28:
          message.externalUserNotFoundText = ExternalUserNotFoundScreenText.decode(reader, reader.uint32());
          break;
        case 29:
          message.successLoginText = SuccessLoginScreenText.decode(reader, reader.uint32());
          break;
        case 30:
          message.logoutText = LogoutDoneScreenText.decode(reader, reader.uint32());
          break;
        case 31:
          message.footerText = FooterText.decode(reader, reader.uint32());
          break;
        case 32:
          message.passwordlessPromptText = PasswordlessPromptScreenText.decode(reader, reader.uint32());
          break;
        case 33:
          message.passwordlessRegistrationText = PasswordlessRegistrationScreenText.decode(reader, reader.uint32());
          break;
        case 34:
          message.passwordlessRegistrationDoneText = PasswordlessRegistrationDoneScreenText.decode(
            reader,
            reader.uint32(),
          );
          break;
        case 35:
          message.externalRegistrationUserOverviewText = ExternalRegistrationUserOverviewScreenText.decode(
            reader,
            reader.uint32(),
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetCustomLoginTextsRequest {
    return {
      language: isSet(object.language) ? String(object.language) : "",
      selectAccountText: isSet(object.selectAccountText)
        ? SelectAccountScreenText.fromJSON(object.selectAccountText)
        : undefined,
      loginText: isSet(object.loginText) ? LoginScreenText.fromJSON(object.loginText) : undefined,
      passwordText: isSet(object.passwordText) ? PasswordScreenText.fromJSON(object.passwordText) : undefined,
      usernameChangeText: isSet(object.usernameChangeText)
        ? UsernameChangeScreenText.fromJSON(object.usernameChangeText)
        : undefined,
      usernameChangeDoneText: isSet(object.usernameChangeDoneText)
        ? UsernameChangeDoneScreenText.fromJSON(object.usernameChangeDoneText)
        : undefined,
      initPasswordText: isSet(object.initPasswordText)
        ? InitPasswordScreenText.fromJSON(object.initPasswordText)
        : undefined,
      initPasswordDoneText: isSet(object.initPasswordDoneText)
        ? InitPasswordDoneScreenText.fromJSON(object.initPasswordDoneText)
        : undefined,
      emailVerificationText: isSet(object.emailVerificationText)
        ? EmailVerificationScreenText.fromJSON(object.emailVerificationText)
        : undefined,
      emailVerificationDoneText: isSet(object.emailVerificationDoneText)
        ? EmailVerificationDoneScreenText.fromJSON(object.emailVerificationDoneText)
        : undefined,
      initializeUserText: isSet(object.initializeUserText)
        ? InitializeUserScreenText.fromJSON(object.initializeUserText)
        : undefined,
      initializeDoneText: isSet(object.initializeDoneText)
        ? InitializeUserDoneScreenText.fromJSON(object.initializeDoneText)
        : undefined,
      initMfaPromptText: isSet(object.initMfaPromptText)
        ? InitMFAPromptScreenText.fromJSON(object.initMfaPromptText)
        : undefined,
      initMfaOtpText: isSet(object.initMfaOtpText) ? InitMFAOTPScreenText.fromJSON(object.initMfaOtpText) : undefined,
      initMfaU2fText: isSet(object.initMfaU2fText) ? InitMFAU2FScreenText.fromJSON(object.initMfaU2fText) : undefined,
      initMfaDoneText: isSet(object.initMfaDoneText)
        ? InitMFADoneScreenText.fromJSON(object.initMfaDoneText)
        : undefined,
      mfaProvidersText: isSet(object.mfaProvidersText) ? MFAProvidersText.fromJSON(object.mfaProvidersText) : undefined,
      verifyMfaOtpText: isSet(object.verifyMfaOtpText)
        ? VerifyMFAOTPScreenText.fromJSON(object.verifyMfaOtpText)
        : undefined,
      verifyMfaU2fText: isSet(object.verifyMfaU2fText)
        ? VerifyMFAU2FScreenText.fromJSON(object.verifyMfaU2fText)
        : undefined,
      passwordlessText: isSet(object.passwordlessText)
        ? PasswordlessScreenText.fromJSON(object.passwordlessText)
        : undefined,
      passwordChangeText: isSet(object.passwordChangeText)
        ? PasswordChangeScreenText.fromJSON(object.passwordChangeText)
        : undefined,
      passwordChangeDoneText: isSet(object.passwordChangeDoneText)
        ? PasswordChangeDoneScreenText.fromJSON(object.passwordChangeDoneText)
        : undefined,
      passwordResetDoneText: isSet(object.passwordResetDoneText)
        ? PasswordResetDoneScreenText.fromJSON(object.passwordResetDoneText)
        : undefined,
      registrationOptionText: isSet(object.registrationOptionText)
        ? RegistrationOptionScreenText.fromJSON(object.registrationOptionText)
        : undefined,
      registrationUserText: isSet(object.registrationUserText)
        ? RegistrationUserScreenText.fromJSON(object.registrationUserText)
        : undefined,
      registrationOrgText: isSet(object.registrationOrgText)
        ? RegistrationOrgScreenText.fromJSON(object.registrationOrgText)
        : undefined,
      linkingUserDoneText: isSet(object.linkingUserDoneText)
        ? LinkingUserDoneScreenText.fromJSON(object.linkingUserDoneText)
        : undefined,
      externalUserNotFoundText: isSet(object.externalUserNotFoundText)
        ? ExternalUserNotFoundScreenText.fromJSON(object.externalUserNotFoundText)
        : undefined,
      successLoginText: isSet(object.successLoginText)
        ? SuccessLoginScreenText.fromJSON(object.successLoginText)
        : undefined,
      logoutText: isSet(object.logoutText) ? LogoutDoneScreenText.fromJSON(object.logoutText) : undefined,
      footerText: isSet(object.footerText) ? FooterText.fromJSON(object.footerText) : undefined,
      passwordlessPromptText: isSet(object.passwordlessPromptText)
        ? PasswordlessPromptScreenText.fromJSON(object.passwordlessPromptText)
        : undefined,
      passwordlessRegistrationText: isSet(object.passwordlessRegistrationText)
        ? PasswordlessRegistrationScreenText.fromJSON(object.passwordlessRegistrationText)
        : undefined,
      passwordlessRegistrationDoneText: isSet(object.passwordlessRegistrationDoneText)
        ? PasswordlessRegistrationDoneScreenText.fromJSON(object.passwordlessRegistrationDoneText)
        : undefined,
      externalRegistrationUserOverviewText: isSet(object.externalRegistrationUserOverviewText)
        ? ExternalRegistrationUserOverviewScreenText.fromJSON(object.externalRegistrationUserOverviewText)
        : undefined,
    };
  },

  toJSON(message: SetCustomLoginTextsRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    message.selectAccountText !== undefined && (obj.selectAccountText = message.selectAccountText
      ? SelectAccountScreenText.toJSON(message.selectAccountText)
      : undefined);
    message.loginText !== undefined &&
      (obj.loginText = message.loginText ? LoginScreenText.toJSON(message.loginText) : undefined);
    message.passwordText !== undefined &&
      (obj.passwordText = message.passwordText ? PasswordScreenText.toJSON(message.passwordText) : undefined);
    message.usernameChangeText !== undefined && (obj.usernameChangeText = message.usernameChangeText
      ? UsernameChangeScreenText.toJSON(message.usernameChangeText)
      : undefined);
    message.usernameChangeDoneText !== undefined && (obj.usernameChangeDoneText = message.usernameChangeDoneText
      ? UsernameChangeDoneScreenText.toJSON(message.usernameChangeDoneText)
      : undefined);
    message.initPasswordText !== undefined && (obj.initPasswordText = message.initPasswordText
      ? InitPasswordScreenText.toJSON(message.initPasswordText)
      : undefined);
    message.initPasswordDoneText !== undefined && (obj.initPasswordDoneText = message.initPasswordDoneText
      ? InitPasswordDoneScreenText.toJSON(message.initPasswordDoneText)
      : undefined);
    message.emailVerificationText !== undefined && (obj.emailVerificationText = message.emailVerificationText
      ? EmailVerificationScreenText.toJSON(message.emailVerificationText)
      : undefined);
    message.emailVerificationDoneText !== undefined &&
      (obj.emailVerificationDoneText = message.emailVerificationDoneText
        ? EmailVerificationDoneScreenText.toJSON(message.emailVerificationDoneText)
        : undefined);
    message.initializeUserText !== undefined && (obj.initializeUserText = message.initializeUserText
      ? InitializeUserScreenText.toJSON(message.initializeUserText)
      : undefined);
    message.initializeDoneText !== undefined && (obj.initializeDoneText = message.initializeDoneText
      ? InitializeUserDoneScreenText.toJSON(message.initializeDoneText)
      : undefined);
    message.initMfaPromptText !== undefined && (obj.initMfaPromptText = message.initMfaPromptText
      ? InitMFAPromptScreenText.toJSON(message.initMfaPromptText)
      : undefined);
    message.initMfaOtpText !== undefined &&
      (obj.initMfaOtpText = message.initMfaOtpText ? InitMFAOTPScreenText.toJSON(message.initMfaOtpText) : undefined);
    message.initMfaU2fText !== undefined &&
      (obj.initMfaU2fText = message.initMfaU2fText ? InitMFAU2FScreenText.toJSON(message.initMfaU2fText) : undefined);
    message.initMfaDoneText !== undefined &&
      (obj.initMfaDoneText = message.initMfaDoneText
        ? InitMFADoneScreenText.toJSON(message.initMfaDoneText)
        : undefined);
    message.mfaProvidersText !== undefined &&
      (obj.mfaProvidersText = message.mfaProvidersText ? MFAProvidersText.toJSON(message.mfaProvidersText) : undefined);
    message.verifyMfaOtpText !== undefined && (obj.verifyMfaOtpText = message.verifyMfaOtpText
      ? VerifyMFAOTPScreenText.toJSON(message.verifyMfaOtpText)
      : undefined);
    message.verifyMfaU2fText !== undefined && (obj.verifyMfaU2fText = message.verifyMfaU2fText
      ? VerifyMFAU2FScreenText.toJSON(message.verifyMfaU2fText)
      : undefined);
    message.passwordlessText !== undefined && (obj.passwordlessText = message.passwordlessText
      ? PasswordlessScreenText.toJSON(message.passwordlessText)
      : undefined);
    message.passwordChangeText !== undefined && (obj.passwordChangeText = message.passwordChangeText
      ? PasswordChangeScreenText.toJSON(message.passwordChangeText)
      : undefined);
    message.passwordChangeDoneText !== undefined && (obj.passwordChangeDoneText = message.passwordChangeDoneText
      ? PasswordChangeDoneScreenText.toJSON(message.passwordChangeDoneText)
      : undefined);
    message.passwordResetDoneText !== undefined && (obj.passwordResetDoneText = message.passwordResetDoneText
      ? PasswordResetDoneScreenText.toJSON(message.passwordResetDoneText)
      : undefined);
    message.registrationOptionText !== undefined && (obj.registrationOptionText = message.registrationOptionText
      ? RegistrationOptionScreenText.toJSON(message.registrationOptionText)
      : undefined);
    message.registrationUserText !== undefined && (obj.registrationUserText = message.registrationUserText
      ? RegistrationUserScreenText.toJSON(message.registrationUserText)
      : undefined);
    message.registrationOrgText !== undefined && (obj.registrationOrgText = message.registrationOrgText
      ? RegistrationOrgScreenText.toJSON(message.registrationOrgText)
      : undefined);
    message.linkingUserDoneText !== undefined && (obj.linkingUserDoneText = message.linkingUserDoneText
      ? LinkingUserDoneScreenText.toJSON(message.linkingUserDoneText)
      : undefined);
    message.externalUserNotFoundText !== undefined && (obj.externalUserNotFoundText = message.externalUserNotFoundText
      ? ExternalUserNotFoundScreenText.toJSON(message.externalUserNotFoundText)
      : undefined);
    message.successLoginText !== undefined && (obj.successLoginText = message.successLoginText
      ? SuccessLoginScreenText.toJSON(message.successLoginText)
      : undefined);
    message.logoutText !== undefined &&
      (obj.logoutText = message.logoutText ? LogoutDoneScreenText.toJSON(message.logoutText) : undefined);
    message.footerText !== undefined &&
      (obj.footerText = message.footerText ? FooterText.toJSON(message.footerText) : undefined);
    message.passwordlessPromptText !== undefined && (obj.passwordlessPromptText = message.passwordlessPromptText
      ? PasswordlessPromptScreenText.toJSON(message.passwordlessPromptText)
      : undefined);
    message.passwordlessRegistrationText !== undefined &&
      (obj.passwordlessRegistrationText = message.passwordlessRegistrationText
        ? PasswordlessRegistrationScreenText.toJSON(message.passwordlessRegistrationText)
        : undefined);
    message.passwordlessRegistrationDoneText !== undefined &&
      (obj.passwordlessRegistrationDoneText = message.passwordlessRegistrationDoneText
        ? PasswordlessRegistrationDoneScreenText.toJSON(message.passwordlessRegistrationDoneText)
        : undefined);
    message.externalRegistrationUserOverviewText !== undefined &&
      (obj.externalRegistrationUserOverviewText = message.externalRegistrationUserOverviewText
        ? ExternalRegistrationUserOverviewScreenText.toJSON(message.externalRegistrationUserOverviewText)
        : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetCustomLoginTextsRequest>): SetCustomLoginTextsRequest {
    return SetCustomLoginTextsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetCustomLoginTextsRequest>): SetCustomLoginTextsRequest {
    const message = createBaseSetCustomLoginTextsRequest();
    message.language = object.language ?? "";
    message.selectAccountText = (object.selectAccountText !== undefined && object.selectAccountText !== null)
      ? SelectAccountScreenText.fromPartial(object.selectAccountText)
      : undefined;
    message.loginText = (object.loginText !== undefined && object.loginText !== null)
      ? LoginScreenText.fromPartial(object.loginText)
      : undefined;
    message.passwordText = (object.passwordText !== undefined && object.passwordText !== null)
      ? PasswordScreenText.fromPartial(object.passwordText)
      : undefined;
    message.usernameChangeText = (object.usernameChangeText !== undefined && object.usernameChangeText !== null)
      ? UsernameChangeScreenText.fromPartial(object.usernameChangeText)
      : undefined;
    message.usernameChangeDoneText =
      (object.usernameChangeDoneText !== undefined && object.usernameChangeDoneText !== null)
        ? UsernameChangeDoneScreenText.fromPartial(object.usernameChangeDoneText)
        : undefined;
    message.initPasswordText = (object.initPasswordText !== undefined && object.initPasswordText !== null)
      ? InitPasswordScreenText.fromPartial(object.initPasswordText)
      : undefined;
    message.initPasswordDoneText = (object.initPasswordDoneText !== undefined && object.initPasswordDoneText !== null)
      ? InitPasswordDoneScreenText.fromPartial(object.initPasswordDoneText)
      : undefined;
    message.emailVerificationText =
      (object.emailVerificationText !== undefined && object.emailVerificationText !== null)
        ? EmailVerificationScreenText.fromPartial(object.emailVerificationText)
        : undefined;
    message.emailVerificationDoneText =
      (object.emailVerificationDoneText !== undefined && object.emailVerificationDoneText !== null)
        ? EmailVerificationDoneScreenText.fromPartial(object.emailVerificationDoneText)
        : undefined;
    message.initializeUserText = (object.initializeUserText !== undefined && object.initializeUserText !== null)
      ? InitializeUserScreenText.fromPartial(object.initializeUserText)
      : undefined;
    message.initializeDoneText = (object.initializeDoneText !== undefined && object.initializeDoneText !== null)
      ? InitializeUserDoneScreenText.fromPartial(object.initializeDoneText)
      : undefined;
    message.initMfaPromptText = (object.initMfaPromptText !== undefined && object.initMfaPromptText !== null)
      ? InitMFAPromptScreenText.fromPartial(object.initMfaPromptText)
      : undefined;
    message.initMfaOtpText = (object.initMfaOtpText !== undefined && object.initMfaOtpText !== null)
      ? InitMFAOTPScreenText.fromPartial(object.initMfaOtpText)
      : undefined;
    message.initMfaU2fText = (object.initMfaU2fText !== undefined && object.initMfaU2fText !== null)
      ? InitMFAU2FScreenText.fromPartial(object.initMfaU2fText)
      : undefined;
    message.initMfaDoneText = (object.initMfaDoneText !== undefined && object.initMfaDoneText !== null)
      ? InitMFADoneScreenText.fromPartial(object.initMfaDoneText)
      : undefined;
    message.mfaProvidersText = (object.mfaProvidersText !== undefined && object.mfaProvidersText !== null)
      ? MFAProvidersText.fromPartial(object.mfaProvidersText)
      : undefined;
    message.verifyMfaOtpText = (object.verifyMfaOtpText !== undefined && object.verifyMfaOtpText !== null)
      ? VerifyMFAOTPScreenText.fromPartial(object.verifyMfaOtpText)
      : undefined;
    message.verifyMfaU2fText = (object.verifyMfaU2fText !== undefined && object.verifyMfaU2fText !== null)
      ? VerifyMFAU2FScreenText.fromPartial(object.verifyMfaU2fText)
      : undefined;
    message.passwordlessText = (object.passwordlessText !== undefined && object.passwordlessText !== null)
      ? PasswordlessScreenText.fromPartial(object.passwordlessText)
      : undefined;
    message.passwordChangeText = (object.passwordChangeText !== undefined && object.passwordChangeText !== null)
      ? PasswordChangeScreenText.fromPartial(object.passwordChangeText)
      : undefined;
    message.passwordChangeDoneText =
      (object.passwordChangeDoneText !== undefined && object.passwordChangeDoneText !== null)
        ? PasswordChangeDoneScreenText.fromPartial(object.passwordChangeDoneText)
        : undefined;
    message.passwordResetDoneText =
      (object.passwordResetDoneText !== undefined && object.passwordResetDoneText !== null)
        ? PasswordResetDoneScreenText.fromPartial(object.passwordResetDoneText)
        : undefined;
    message.registrationOptionText =
      (object.registrationOptionText !== undefined && object.registrationOptionText !== null)
        ? RegistrationOptionScreenText.fromPartial(object.registrationOptionText)
        : undefined;
    message.registrationUserText = (object.registrationUserText !== undefined && object.registrationUserText !== null)
      ? RegistrationUserScreenText.fromPartial(object.registrationUserText)
      : undefined;
    message.registrationOrgText = (object.registrationOrgText !== undefined && object.registrationOrgText !== null)
      ? RegistrationOrgScreenText.fromPartial(object.registrationOrgText)
      : undefined;
    message.linkingUserDoneText = (object.linkingUserDoneText !== undefined && object.linkingUserDoneText !== null)
      ? LinkingUserDoneScreenText.fromPartial(object.linkingUserDoneText)
      : undefined;
    message.externalUserNotFoundText =
      (object.externalUserNotFoundText !== undefined && object.externalUserNotFoundText !== null)
        ? ExternalUserNotFoundScreenText.fromPartial(object.externalUserNotFoundText)
        : undefined;
    message.successLoginText = (object.successLoginText !== undefined && object.successLoginText !== null)
      ? SuccessLoginScreenText.fromPartial(object.successLoginText)
      : undefined;
    message.logoutText = (object.logoutText !== undefined && object.logoutText !== null)
      ? LogoutDoneScreenText.fromPartial(object.logoutText)
      : undefined;
    message.footerText = (object.footerText !== undefined && object.footerText !== null)
      ? FooterText.fromPartial(object.footerText)
      : undefined;
    message.passwordlessPromptText =
      (object.passwordlessPromptText !== undefined && object.passwordlessPromptText !== null)
        ? PasswordlessPromptScreenText.fromPartial(object.passwordlessPromptText)
        : undefined;
    message.passwordlessRegistrationText =
      (object.passwordlessRegistrationText !== undefined && object.passwordlessRegistrationText !== null)
        ? PasswordlessRegistrationScreenText.fromPartial(object.passwordlessRegistrationText)
        : undefined;
    message.passwordlessRegistrationDoneText =
      (object.passwordlessRegistrationDoneText !== undefined && object.passwordlessRegistrationDoneText !== null)
        ? PasswordlessRegistrationDoneScreenText.fromPartial(object.passwordlessRegistrationDoneText)
        : undefined;
    message.externalRegistrationUserOverviewText =
      (object.externalRegistrationUserOverviewText !== undefined &&
          object.externalRegistrationUserOverviewText !== null)
        ? ExternalRegistrationUserOverviewScreenText.fromPartial(object.externalRegistrationUserOverviewText)
        : undefined;
    return message;
  },
};

function createBaseSetCustomLoginTextsResponse(): SetCustomLoginTextsResponse {
  return { details: undefined };
}

export const SetCustomLoginTextsResponse = {
  encode(message: SetCustomLoginTextsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetCustomLoginTextsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetCustomLoginTextsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetCustomLoginTextsResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetCustomLoginTextsResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetCustomLoginTextsResponse>): SetCustomLoginTextsResponse {
    return SetCustomLoginTextsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetCustomLoginTextsResponse>): SetCustomLoginTextsResponse {
    const message = createBaseSetCustomLoginTextsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetCustomLoginTextsToDefaultRequest(): ResetCustomLoginTextsToDefaultRequest {
  return { language: "" };
}

export const ResetCustomLoginTextsToDefaultRequest = {
  encode(message: ResetCustomLoginTextsToDefaultRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.language !== "") {
      writer.uint32(10).string(message.language);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomLoginTextsToDefaultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomLoginTextsToDefaultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.language = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomLoginTextsToDefaultRequest {
    return { language: isSet(object.language) ? String(object.language) : "" };
  },

  toJSON(message: ResetCustomLoginTextsToDefaultRequest): unknown {
    const obj: any = {};
    message.language !== undefined && (obj.language = message.language);
    return obj;
  },

  create(base?: DeepPartial<ResetCustomLoginTextsToDefaultRequest>): ResetCustomLoginTextsToDefaultRequest {
    return ResetCustomLoginTextsToDefaultRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetCustomLoginTextsToDefaultRequest>): ResetCustomLoginTextsToDefaultRequest {
    const message = createBaseResetCustomLoginTextsToDefaultRequest();
    message.language = object.language ?? "";
    return message;
  },
};

function createBaseResetCustomLoginTextsToDefaultResponse(): ResetCustomLoginTextsToDefaultResponse {
  return { details: undefined };
}

export const ResetCustomLoginTextsToDefaultResponse = {
  encode(message: ResetCustomLoginTextsToDefaultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetCustomLoginTextsToDefaultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetCustomLoginTextsToDefaultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResetCustomLoginTextsToDefaultResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetCustomLoginTextsToDefaultResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResetCustomLoginTextsToDefaultResponse>): ResetCustomLoginTextsToDefaultResponse {
    return ResetCustomLoginTextsToDefaultResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetCustomLoginTextsToDefaultResponse>): ResetCustomLoginTextsToDefaultResponse {
    const message = createBaseResetCustomLoginTextsToDefaultResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddIAMMemberRequest(): AddIAMMemberRequest {
  return { userId: "", roles: [] };
}

export const AddIAMMemberRequest = {
  encode(message: AddIAMMemberRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    for (const v of message.roles) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddIAMMemberRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddIAMMemberRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userId = reader.string();
          break;
        case 2:
          message.roles.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddIAMMemberRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      roles: Array.isArray(object?.roles) ? object.roles.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: AddIAMMemberRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    if (message.roles) {
      obj.roles = message.roles.map((e) => e);
    } else {
      obj.roles = [];
    }
    return obj;
  },

  create(base?: DeepPartial<AddIAMMemberRequest>): AddIAMMemberRequest {
    return AddIAMMemberRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddIAMMemberRequest>): AddIAMMemberRequest {
    const message = createBaseAddIAMMemberRequest();
    message.userId = object.userId ?? "";
    message.roles = object.roles?.map((e) => e) || [];
    return message;
  },
};

function createBaseAddIAMMemberResponse(): AddIAMMemberResponse {
  return { details: undefined };
}

export const AddIAMMemberResponse = {
  encode(message: AddIAMMemberResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddIAMMemberResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddIAMMemberResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddIAMMemberResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddIAMMemberResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddIAMMemberResponse>): AddIAMMemberResponse {
    return AddIAMMemberResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddIAMMemberResponse>): AddIAMMemberResponse {
    const message = createBaseAddIAMMemberResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateIAMMemberRequest(): UpdateIAMMemberRequest {
  return { userId: "", roles: [] };
}

export const UpdateIAMMemberRequest = {
  encode(message: UpdateIAMMemberRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    for (const v of message.roles) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIAMMemberRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIAMMemberRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userId = reader.string();
          break;
        case 2:
          message.roles.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateIAMMemberRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      roles: Array.isArray(object?.roles) ? object.roles.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: UpdateIAMMemberRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    if (message.roles) {
      obj.roles = message.roles.map((e) => e);
    } else {
      obj.roles = [];
    }
    return obj;
  },

  create(base?: DeepPartial<UpdateIAMMemberRequest>): UpdateIAMMemberRequest {
    return UpdateIAMMemberRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateIAMMemberRequest>): UpdateIAMMemberRequest {
    const message = createBaseUpdateIAMMemberRequest();
    message.userId = object.userId ?? "";
    message.roles = object.roles?.map((e) => e) || [];
    return message;
  },
};

function createBaseUpdateIAMMemberResponse(): UpdateIAMMemberResponse {
  return { details: undefined };
}

export const UpdateIAMMemberResponse = {
  encode(message: UpdateIAMMemberResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIAMMemberResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIAMMemberResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateIAMMemberResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateIAMMemberResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateIAMMemberResponse>): UpdateIAMMemberResponse {
    return UpdateIAMMemberResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateIAMMemberResponse>): UpdateIAMMemberResponse {
    const message = createBaseUpdateIAMMemberResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveIAMMemberRequest(): RemoveIAMMemberRequest {
  return { userId: "" };
}

export const RemoveIAMMemberRequest = {
  encode(message: RemoveIAMMemberRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveIAMMemberRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveIAMMemberRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveIAMMemberRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: RemoveIAMMemberRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<RemoveIAMMemberRequest>): RemoveIAMMemberRequest {
    return RemoveIAMMemberRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveIAMMemberRequest>): RemoveIAMMemberRequest {
    const message = createBaseRemoveIAMMemberRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseRemoveIAMMemberResponse(): RemoveIAMMemberResponse {
  return { details: undefined };
}

export const RemoveIAMMemberResponse = {
  encode(message: RemoveIAMMemberResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveIAMMemberResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveIAMMemberResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveIAMMemberResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveIAMMemberResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveIAMMemberResponse>): RemoveIAMMemberResponse {
    return RemoveIAMMemberResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveIAMMemberResponse>): RemoveIAMMemberResponse {
    const message = createBaseRemoveIAMMemberResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListIAMMemberRolesRequest(): ListIAMMemberRolesRequest {
  return {};
}

export const ListIAMMemberRolesRequest = {
  encode(_: ListIAMMemberRolesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListIAMMemberRolesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListIAMMemberRolesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListIAMMemberRolesRequest {
    return {};
  },

  toJSON(_: ListIAMMemberRolesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListIAMMemberRolesRequest>): ListIAMMemberRolesRequest {
    return ListIAMMemberRolesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListIAMMemberRolesRequest>): ListIAMMemberRolesRequest {
    const message = createBaseListIAMMemberRolesRequest();
    return message;
  },
};

function createBaseListIAMMemberRolesResponse(): ListIAMMemberRolesResponse {
  return { details: undefined, roles: [] };
}

export const ListIAMMemberRolesResponse = {
  encode(message: ListIAMMemberRolesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.roles) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListIAMMemberRolesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListIAMMemberRolesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.roles.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListIAMMemberRolesResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      roles: Array.isArray(object?.roles) ? object.roles.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: ListIAMMemberRolesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.roles) {
      obj.roles = message.roles.map((e) => e);
    } else {
      obj.roles = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListIAMMemberRolesResponse>): ListIAMMemberRolesResponse {
    return ListIAMMemberRolesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListIAMMemberRolesResponse>): ListIAMMemberRolesResponse {
    const message = createBaseListIAMMemberRolesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.roles = object.roles?.map((e) => e) || [];
    return message;
  },
};

function createBaseListIAMMembersRequest(): ListIAMMembersRequest {
  return { query: undefined, queries: [] };
}

export const ListIAMMembersRequest = {
  encode(message: ListIAMMembersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.queries) {
      SearchQuery.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListIAMMembersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListIAMMembersRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.queries.push(SearchQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListIAMMembersRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => SearchQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListIAMMembersRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? SearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListIAMMembersRequest>): ListIAMMembersRequest {
    return ListIAMMembersRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListIAMMembersRequest>): ListIAMMembersRequest {
    const message = createBaseListIAMMembersRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.queries = object.queries?.map((e) => SearchQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListIAMMembersResponse(): ListIAMMembersResponse {
  return { details: undefined, result: [] };
}

export const ListIAMMembersResponse = {
  encode(message: ListIAMMembersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      Member.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListIAMMembersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListIAMMembersResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(Member.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListIAMMembersResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Member.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListIAMMembersResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? Member.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListIAMMembersResponse>): ListIAMMembersResponse {
    return ListIAMMembersResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListIAMMembersResponse>): ListIAMMembersResponse {
    const message = createBaseListIAMMembersResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => Member.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListViewsRequest(): ListViewsRequest {
  return {};
}

export const ListViewsRequest = {
  encode(_: ListViewsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListViewsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListViewsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListViewsRequest {
    return {};
  },

  toJSON(_: ListViewsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListViewsRequest>): ListViewsRequest {
    return ListViewsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListViewsRequest>): ListViewsRequest {
    const message = createBaseListViewsRequest();
    return message;
  },
};

function createBaseListViewsResponse(): ListViewsResponse {
  return { result: [] };
}

export const ListViewsResponse = {
  encode(message: ListViewsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.result) {
      View.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListViewsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListViewsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.result.push(View.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListViewsResponse {
    return { result: Array.isArray(object?.result) ? object.result.map((e: any) => View.fromJSON(e)) : [] };
  },

  toJSON(message: ListViewsResponse): unknown {
    const obj: any = {};
    if (message.result) {
      obj.result = message.result.map((e) => e ? View.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListViewsResponse>): ListViewsResponse {
    return ListViewsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListViewsResponse>): ListViewsResponse {
    const message = createBaseListViewsResponse();
    message.result = object.result?.map((e) => View.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListFailedEventsRequest(): ListFailedEventsRequest {
  return {};
}

export const ListFailedEventsRequest = {
  encode(_: ListFailedEventsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListFailedEventsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListFailedEventsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListFailedEventsRequest {
    return {};
  },

  toJSON(_: ListFailedEventsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListFailedEventsRequest>): ListFailedEventsRequest {
    return ListFailedEventsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListFailedEventsRequest>): ListFailedEventsRequest {
    const message = createBaseListFailedEventsRequest();
    return message;
  },
};

function createBaseListFailedEventsResponse(): ListFailedEventsResponse {
  return { result: [] };
}

export const ListFailedEventsResponse = {
  encode(message: ListFailedEventsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.result) {
      FailedEvent.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListFailedEventsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListFailedEventsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.result.push(FailedEvent.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListFailedEventsResponse {
    return { result: Array.isArray(object?.result) ? object.result.map((e: any) => FailedEvent.fromJSON(e)) : [] };
  },

  toJSON(message: ListFailedEventsResponse): unknown {
    const obj: any = {};
    if (message.result) {
      obj.result = message.result.map((e) => e ? FailedEvent.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListFailedEventsResponse>): ListFailedEventsResponse {
    return ListFailedEventsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListFailedEventsResponse>): ListFailedEventsResponse {
    const message = createBaseListFailedEventsResponse();
    message.result = object.result?.map((e) => FailedEvent.fromPartial(e)) || [];
    return message;
  },
};

function createBaseRemoveFailedEventRequest(): RemoveFailedEventRequest {
  return { database: "", viewName: "", failedSequence: 0 };
}

export const RemoveFailedEventRequest = {
  encode(message: RemoveFailedEventRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.database !== "") {
      writer.uint32(10).string(message.database);
    }
    if (message.viewName !== "") {
      writer.uint32(18).string(message.viewName);
    }
    if (message.failedSequence !== 0) {
      writer.uint32(24).uint64(message.failedSequence);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveFailedEventRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveFailedEventRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.database = reader.string();
          break;
        case 2:
          message.viewName = reader.string();
          break;
        case 3:
          message.failedSequence = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveFailedEventRequest {
    return {
      database: isSet(object.database) ? String(object.database) : "",
      viewName: isSet(object.viewName) ? String(object.viewName) : "",
      failedSequence: isSet(object.failedSequence) ? Number(object.failedSequence) : 0,
    };
  },

  toJSON(message: RemoveFailedEventRequest): unknown {
    const obj: any = {};
    message.database !== undefined && (obj.database = message.database);
    message.viewName !== undefined && (obj.viewName = message.viewName);
    message.failedSequence !== undefined && (obj.failedSequence = Math.round(message.failedSequence));
    return obj;
  },

  create(base?: DeepPartial<RemoveFailedEventRequest>): RemoveFailedEventRequest {
    return RemoveFailedEventRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveFailedEventRequest>): RemoveFailedEventRequest {
    const message = createBaseRemoveFailedEventRequest();
    message.database = object.database ?? "";
    message.viewName = object.viewName ?? "";
    message.failedSequence = object.failedSequence ?? 0;
    return message;
  },
};

function createBaseRemoveFailedEventResponse(): RemoveFailedEventResponse {
  return {};
}

export const RemoveFailedEventResponse = {
  encode(_: RemoveFailedEventResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveFailedEventResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveFailedEventResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveFailedEventResponse {
    return {};
  },

  toJSON(_: RemoveFailedEventResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveFailedEventResponse>): RemoveFailedEventResponse {
    return RemoveFailedEventResponse.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveFailedEventResponse>): RemoveFailedEventResponse {
    const message = createBaseRemoveFailedEventResponse();
    return message;
  },
};

function createBaseView(): View {
  return {
    database: "",
    viewName: "",
    processedSequence: 0,
    eventTimestamp: undefined,
    lastSuccessfulSpoolerRun: undefined,
  };
}

export const View = {
  encode(message: View, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.database !== "") {
      writer.uint32(10).string(message.database);
    }
    if (message.viewName !== "") {
      writer.uint32(18).string(message.viewName);
    }
    if (message.processedSequence !== 0) {
      writer.uint32(24).uint64(message.processedSequence);
    }
    if (message.eventTimestamp !== undefined) {
      Timestamp.encode(toTimestamp(message.eventTimestamp), writer.uint32(34).fork()).ldelim();
    }
    if (message.lastSuccessfulSpoolerRun !== undefined) {
      Timestamp.encode(toTimestamp(message.lastSuccessfulSpoolerRun), writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): View {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseView();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.database = reader.string();
          break;
        case 2:
          message.viewName = reader.string();
          break;
        case 3:
          message.processedSequence = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.eventTimestamp = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        case 5:
          message.lastSuccessfulSpoolerRun = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): View {
    return {
      database: isSet(object.database) ? String(object.database) : "",
      viewName: isSet(object.viewName) ? String(object.viewName) : "",
      processedSequence: isSet(object.processedSequence) ? Number(object.processedSequence) : 0,
      eventTimestamp: isSet(object.eventTimestamp) ? fromJsonTimestamp(object.eventTimestamp) : undefined,
      lastSuccessfulSpoolerRun: isSet(object.lastSuccessfulSpoolerRun)
        ? fromJsonTimestamp(object.lastSuccessfulSpoolerRun)
        : undefined,
    };
  },

  toJSON(message: View): unknown {
    const obj: any = {};
    message.database !== undefined && (obj.database = message.database);
    message.viewName !== undefined && (obj.viewName = message.viewName);
    message.processedSequence !== undefined && (obj.processedSequence = Math.round(message.processedSequence));
    message.eventTimestamp !== undefined && (obj.eventTimestamp = message.eventTimestamp.toISOString());
    message.lastSuccessfulSpoolerRun !== undefined &&
      (obj.lastSuccessfulSpoolerRun = message.lastSuccessfulSpoolerRun.toISOString());
    return obj;
  },

  create(base?: DeepPartial<View>): View {
    return View.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<View>): View {
    const message = createBaseView();
    message.database = object.database ?? "";
    message.viewName = object.viewName ?? "";
    message.processedSequence = object.processedSequence ?? 0;
    message.eventTimestamp = object.eventTimestamp ?? undefined;
    message.lastSuccessfulSpoolerRun = object.lastSuccessfulSpoolerRun ?? undefined;
    return message;
  },
};

function createBaseFailedEvent(): FailedEvent {
  return { database: "", viewName: "", failedSequence: 0, failureCount: 0, errorMessage: "", lastFailed: undefined };
}

export const FailedEvent = {
  encode(message: FailedEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.database !== "") {
      writer.uint32(10).string(message.database);
    }
    if (message.viewName !== "") {
      writer.uint32(18).string(message.viewName);
    }
    if (message.failedSequence !== 0) {
      writer.uint32(24).uint64(message.failedSequence);
    }
    if (message.failureCount !== 0) {
      writer.uint32(32).uint64(message.failureCount);
    }
    if (message.errorMessage !== "") {
      writer.uint32(42).string(message.errorMessage);
    }
    if (message.lastFailed !== undefined) {
      Timestamp.encode(toTimestamp(message.lastFailed), writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FailedEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFailedEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.database = reader.string();
          break;
        case 2:
          message.viewName = reader.string();
          break;
        case 3:
          message.failedSequence = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.failureCount = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.errorMessage = reader.string();
          break;
        case 6:
          message.lastFailed = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FailedEvent {
    return {
      database: isSet(object.database) ? String(object.database) : "",
      viewName: isSet(object.viewName) ? String(object.viewName) : "",
      failedSequence: isSet(object.failedSequence) ? Number(object.failedSequence) : 0,
      failureCount: isSet(object.failureCount) ? Number(object.failureCount) : 0,
      errorMessage: isSet(object.errorMessage) ? String(object.errorMessage) : "",
      lastFailed: isSet(object.lastFailed) ? fromJsonTimestamp(object.lastFailed) : undefined,
    };
  },

  toJSON(message: FailedEvent): unknown {
    const obj: any = {};
    message.database !== undefined && (obj.database = message.database);
    message.viewName !== undefined && (obj.viewName = message.viewName);
    message.failedSequence !== undefined && (obj.failedSequence = Math.round(message.failedSequence));
    message.failureCount !== undefined && (obj.failureCount = Math.round(message.failureCount));
    message.errorMessage !== undefined && (obj.errorMessage = message.errorMessage);
    message.lastFailed !== undefined && (obj.lastFailed = message.lastFailed.toISOString());
    return obj;
  },

  create(base?: DeepPartial<FailedEvent>): FailedEvent {
    return FailedEvent.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<FailedEvent>): FailedEvent {
    const message = createBaseFailedEvent();
    message.database = object.database ?? "";
    message.viewName = object.viewName ?? "";
    message.failedSequence = object.failedSequence ?? 0;
    message.failureCount = object.failureCount ?? 0;
    message.errorMessage = object.errorMessage ?? "";
    message.lastFailed = object.lastFailed ?? undefined;
    return message;
  },
};

function createBaseImportDataRequest(): ImportDataRequest {
  return {
    dataOrgs: undefined,
    dataOrgsv1: undefined,
    dataOrgsLocal: undefined,
    dataOrgsv1Local: undefined,
    dataOrgsS3: undefined,
    dataOrgsv1S3: undefined,
    dataOrgsGcs: undefined,
    dataOrgsv1Gcs: undefined,
    timeout: "",
  };
}

export const ImportDataRequest = {
  encode(message: ImportDataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.dataOrgs !== undefined) {
      ImportDataOrg.encode(message.dataOrgs, writer.uint32(10).fork()).ldelim();
    }
    if (message.dataOrgsv1 !== undefined) {
      ImportDataOrg1.encode(message.dataOrgsv1, writer.uint32(18).fork()).ldelim();
    }
    if (message.dataOrgsLocal !== undefined) {
      ImportDataRequest_LocalInput.encode(message.dataOrgsLocal, writer.uint32(26).fork()).ldelim();
    }
    if (message.dataOrgsv1Local !== undefined) {
      ImportDataRequest_LocalInput.encode(message.dataOrgsv1Local, writer.uint32(34).fork()).ldelim();
    }
    if (message.dataOrgsS3 !== undefined) {
      ImportDataRequest_S3Input.encode(message.dataOrgsS3, writer.uint32(42).fork()).ldelim();
    }
    if (message.dataOrgsv1S3 !== undefined) {
      ImportDataRequest_S3Input.encode(message.dataOrgsv1S3, writer.uint32(50).fork()).ldelim();
    }
    if (message.dataOrgsGcs !== undefined) {
      ImportDataRequest_GCSInput.encode(message.dataOrgsGcs, writer.uint32(58).fork()).ldelim();
    }
    if (message.dataOrgsv1Gcs !== undefined) {
      ImportDataRequest_GCSInput.encode(message.dataOrgsv1Gcs, writer.uint32(66).fork()).ldelim();
    }
    if (message.timeout !== "") {
      writer.uint32(74).string(message.timeout);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.dataOrgs = ImportDataOrg.decode(reader, reader.uint32());
          break;
        case 2:
          message.dataOrgsv1 = ImportDataOrg1.decode(reader, reader.uint32());
          break;
        case 3:
          message.dataOrgsLocal = ImportDataRequest_LocalInput.decode(reader, reader.uint32());
          break;
        case 4:
          message.dataOrgsv1Local = ImportDataRequest_LocalInput.decode(reader, reader.uint32());
          break;
        case 5:
          message.dataOrgsS3 = ImportDataRequest_S3Input.decode(reader, reader.uint32());
          break;
        case 6:
          message.dataOrgsv1S3 = ImportDataRequest_S3Input.decode(reader, reader.uint32());
          break;
        case 7:
          message.dataOrgsGcs = ImportDataRequest_GCSInput.decode(reader, reader.uint32());
          break;
        case 8:
          message.dataOrgsv1Gcs = ImportDataRequest_GCSInput.decode(reader, reader.uint32());
          break;
        case 9:
          message.timeout = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataRequest {
    return {
      dataOrgs: isSet(object.dataOrgs) ? ImportDataOrg.fromJSON(object.dataOrgs) : undefined,
      dataOrgsv1: isSet(object.dataOrgsv1) ? ImportDataOrg1.fromJSON(object.dataOrgsv1) : undefined,
      dataOrgsLocal: isSet(object.dataOrgsLocal)
        ? ImportDataRequest_LocalInput.fromJSON(object.dataOrgsLocal)
        : undefined,
      dataOrgsv1Local: isSet(object.dataOrgsv1Local)
        ? ImportDataRequest_LocalInput.fromJSON(object.dataOrgsv1Local)
        : undefined,
      dataOrgsS3: isSet(object.dataOrgsS3) ? ImportDataRequest_S3Input.fromJSON(object.dataOrgsS3) : undefined,
      dataOrgsv1S3: isSet(object.dataOrgsv1S3) ? ImportDataRequest_S3Input.fromJSON(object.dataOrgsv1S3) : undefined,
      dataOrgsGcs: isSet(object.dataOrgsGcs) ? ImportDataRequest_GCSInput.fromJSON(object.dataOrgsGcs) : undefined,
      dataOrgsv1Gcs: isSet(object.dataOrgsv1Gcs)
        ? ImportDataRequest_GCSInput.fromJSON(object.dataOrgsv1Gcs)
        : undefined,
      timeout: isSet(object.timeout) ? String(object.timeout) : "",
    };
  },

  toJSON(message: ImportDataRequest): unknown {
    const obj: any = {};
    message.dataOrgs !== undefined &&
      (obj.dataOrgs = message.dataOrgs ? ImportDataOrg.toJSON(message.dataOrgs) : undefined);
    message.dataOrgsv1 !== undefined &&
      (obj.dataOrgsv1 = message.dataOrgsv1 ? ImportDataOrg1.toJSON(message.dataOrgsv1) : undefined);
    message.dataOrgsLocal !== undefined && (obj.dataOrgsLocal = message.dataOrgsLocal
      ? ImportDataRequest_LocalInput.toJSON(message.dataOrgsLocal)
      : undefined);
    message.dataOrgsv1Local !== undefined && (obj.dataOrgsv1Local = message.dataOrgsv1Local
      ? ImportDataRequest_LocalInput.toJSON(message.dataOrgsv1Local)
      : undefined);
    message.dataOrgsS3 !== undefined &&
      (obj.dataOrgsS3 = message.dataOrgsS3 ? ImportDataRequest_S3Input.toJSON(message.dataOrgsS3) : undefined);
    message.dataOrgsv1S3 !== undefined &&
      (obj.dataOrgsv1S3 = message.dataOrgsv1S3 ? ImportDataRequest_S3Input.toJSON(message.dataOrgsv1S3) : undefined);
    message.dataOrgsGcs !== undefined &&
      (obj.dataOrgsGcs = message.dataOrgsGcs ? ImportDataRequest_GCSInput.toJSON(message.dataOrgsGcs) : undefined);
    message.dataOrgsv1Gcs !== undefined &&
      (obj.dataOrgsv1Gcs = message.dataOrgsv1Gcs
        ? ImportDataRequest_GCSInput.toJSON(message.dataOrgsv1Gcs)
        : undefined);
    message.timeout !== undefined && (obj.timeout = message.timeout);
    return obj;
  },

  create(base?: DeepPartial<ImportDataRequest>): ImportDataRequest {
    return ImportDataRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataRequest>): ImportDataRequest {
    const message = createBaseImportDataRequest();
    message.dataOrgs = (object.dataOrgs !== undefined && object.dataOrgs !== null)
      ? ImportDataOrg.fromPartial(object.dataOrgs)
      : undefined;
    message.dataOrgsv1 = (object.dataOrgsv1 !== undefined && object.dataOrgsv1 !== null)
      ? ImportDataOrg1.fromPartial(object.dataOrgsv1)
      : undefined;
    message.dataOrgsLocal = (object.dataOrgsLocal !== undefined && object.dataOrgsLocal !== null)
      ? ImportDataRequest_LocalInput.fromPartial(object.dataOrgsLocal)
      : undefined;
    message.dataOrgsv1Local = (object.dataOrgsv1Local !== undefined && object.dataOrgsv1Local !== null)
      ? ImportDataRequest_LocalInput.fromPartial(object.dataOrgsv1Local)
      : undefined;
    message.dataOrgsS3 = (object.dataOrgsS3 !== undefined && object.dataOrgsS3 !== null)
      ? ImportDataRequest_S3Input.fromPartial(object.dataOrgsS3)
      : undefined;
    message.dataOrgsv1S3 = (object.dataOrgsv1S3 !== undefined && object.dataOrgsv1S3 !== null)
      ? ImportDataRequest_S3Input.fromPartial(object.dataOrgsv1S3)
      : undefined;
    message.dataOrgsGcs = (object.dataOrgsGcs !== undefined && object.dataOrgsGcs !== null)
      ? ImportDataRequest_GCSInput.fromPartial(object.dataOrgsGcs)
      : undefined;
    message.dataOrgsv1Gcs = (object.dataOrgsv1Gcs !== undefined && object.dataOrgsv1Gcs !== null)
      ? ImportDataRequest_GCSInput.fromPartial(object.dataOrgsv1Gcs)
      : undefined;
    message.timeout = object.timeout ?? "";
    return message;
  },
};

function createBaseImportDataRequest_LocalInput(): ImportDataRequest_LocalInput {
  return { path: "" };
}

export const ImportDataRequest_LocalInput = {
  encode(message: ImportDataRequest_LocalInput, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataRequest_LocalInput {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataRequest_LocalInput();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataRequest_LocalInput {
    return { path: isSet(object.path) ? String(object.path) : "" };
  },

  toJSON(message: ImportDataRequest_LocalInput): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    return obj;
  },

  create(base?: DeepPartial<ImportDataRequest_LocalInput>): ImportDataRequest_LocalInput {
    return ImportDataRequest_LocalInput.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataRequest_LocalInput>): ImportDataRequest_LocalInput {
    const message = createBaseImportDataRequest_LocalInput();
    message.path = object.path ?? "";
    return message;
  },
};

function createBaseImportDataRequest_S3Input(): ImportDataRequest_S3Input {
  return { path: "", endpoint: "", accessKeyId: "", secretAccessKey: "", ssl: false, bucket: "" };
}

export const ImportDataRequest_S3Input = {
  encode(message: ImportDataRequest_S3Input, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    if (message.endpoint !== "") {
      writer.uint32(18).string(message.endpoint);
    }
    if (message.accessKeyId !== "") {
      writer.uint32(26).string(message.accessKeyId);
    }
    if (message.secretAccessKey !== "") {
      writer.uint32(34).string(message.secretAccessKey);
    }
    if (message.ssl === true) {
      writer.uint32(40).bool(message.ssl);
    }
    if (message.bucket !== "") {
      writer.uint32(50).string(message.bucket);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataRequest_S3Input {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataRequest_S3Input();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        case 2:
          message.endpoint = reader.string();
          break;
        case 3:
          message.accessKeyId = reader.string();
          break;
        case 4:
          message.secretAccessKey = reader.string();
          break;
        case 5:
          message.ssl = reader.bool();
          break;
        case 6:
          message.bucket = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataRequest_S3Input {
    return {
      path: isSet(object.path) ? String(object.path) : "",
      endpoint: isSet(object.endpoint) ? String(object.endpoint) : "",
      accessKeyId: isSet(object.accessKeyId) ? String(object.accessKeyId) : "",
      secretAccessKey: isSet(object.secretAccessKey) ? String(object.secretAccessKey) : "",
      ssl: isSet(object.ssl) ? Boolean(object.ssl) : false,
      bucket: isSet(object.bucket) ? String(object.bucket) : "",
    };
  },

  toJSON(message: ImportDataRequest_S3Input): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    message.endpoint !== undefined && (obj.endpoint = message.endpoint);
    message.accessKeyId !== undefined && (obj.accessKeyId = message.accessKeyId);
    message.secretAccessKey !== undefined && (obj.secretAccessKey = message.secretAccessKey);
    message.ssl !== undefined && (obj.ssl = message.ssl);
    message.bucket !== undefined && (obj.bucket = message.bucket);
    return obj;
  },

  create(base?: DeepPartial<ImportDataRequest_S3Input>): ImportDataRequest_S3Input {
    return ImportDataRequest_S3Input.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataRequest_S3Input>): ImportDataRequest_S3Input {
    const message = createBaseImportDataRequest_S3Input();
    message.path = object.path ?? "";
    message.endpoint = object.endpoint ?? "";
    message.accessKeyId = object.accessKeyId ?? "";
    message.secretAccessKey = object.secretAccessKey ?? "";
    message.ssl = object.ssl ?? false;
    message.bucket = object.bucket ?? "";
    return message;
  },
};

function createBaseImportDataRequest_GCSInput(): ImportDataRequest_GCSInput {
  return { bucket: "", serviceaccountJson: "", path: "" };
}

export const ImportDataRequest_GCSInput = {
  encode(message: ImportDataRequest_GCSInput, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bucket !== "") {
      writer.uint32(10).string(message.bucket);
    }
    if (message.serviceaccountJson !== "") {
      writer.uint32(18).string(message.serviceaccountJson);
    }
    if (message.path !== "") {
      writer.uint32(26).string(message.path);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataRequest_GCSInput {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataRequest_GCSInput();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.bucket = reader.string();
          break;
        case 2:
          message.serviceaccountJson = reader.string();
          break;
        case 3:
          message.path = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataRequest_GCSInput {
    return {
      bucket: isSet(object.bucket) ? String(object.bucket) : "",
      serviceaccountJson: isSet(object.serviceaccountJson) ? String(object.serviceaccountJson) : "",
      path: isSet(object.path) ? String(object.path) : "",
    };
  },

  toJSON(message: ImportDataRequest_GCSInput): unknown {
    const obj: any = {};
    message.bucket !== undefined && (obj.bucket = message.bucket);
    message.serviceaccountJson !== undefined && (obj.serviceaccountJson = message.serviceaccountJson);
    message.path !== undefined && (obj.path = message.path);
    return obj;
  },

  create(base?: DeepPartial<ImportDataRequest_GCSInput>): ImportDataRequest_GCSInput {
    return ImportDataRequest_GCSInput.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataRequest_GCSInput>): ImportDataRequest_GCSInput {
    const message = createBaseImportDataRequest_GCSInput();
    message.bucket = object.bucket ?? "";
    message.serviceaccountJson = object.serviceaccountJson ?? "";
    message.path = object.path ?? "";
    return message;
  },
};

function createBaseImportDataOrg(): ImportDataOrg {
  return { orgs: [] };
}

export const ImportDataOrg = {
  encode(message: ImportDataOrg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.orgs) {
      DataOrg.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataOrg {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataOrg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgs.push(DataOrg.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataOrg {
    return { orgs: Array.isArray(object?.orgs) ? object.orgs.map((e: any) => DataOrg.fromJSON(e)) : [] };
  },

  toJSON(message: ImportDataOrg): unknown {
    const obj: any = {};
    if (message.orgs) {
      obj.orgs = message.orgs.map((e) => e ? DataOrg.toJSON(e) : undefined);
    } else {
      obj.orgs = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ImportDataOrg>): ImportDataOrg {
    return ImportDataOrg.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataOrg>): ImportDataOrg {
    const message = createBaseImportDataOrg();
    message.orgs = object.orgs?.map((e) => DataOrg.fromPartial(e)) || [];
    return message;
  },
};

function createBaseDataOrg(): DataOrg {
  return {
    orgId: "",
    org: undefined,
    domainPolicy: undefined,
    labelPolicy: undefined,
    lockoutPolicy: undefined,
    loginPolicy: undefined,
    passwordComplexityPolicy: undefined,
    privacyPolicy: undefined,
    projects: [],
    projectRoles: [],
    apiApps: [],
    oidcApps: [],
    humanUsers: [],
    machineUsers: [],
    triggerActions: [],
    actions: [],
    projectGrants: [],
    userGrants: [],
    orgMembers: [],
    projectMembers: [],
    projectGrantMembers: [],
    userMetadata: [],
    loginTexts: [],
    initMessages: [],
    passwordResetMessages: [],
    verifyEmailMessages: [],
    verifyPhoneMessages: [],
    domainClaimedMessages: [],
    passwordlessRegistrationMessages: [],
    oidcIdps: [],
    jwtIdps: [],
    userLinks: [],
    domains: [],
    appKeys: [],
    machineKeys: [],
  };
}

export const DataOrg = {
  encode(message: DataOrg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    if (message.org !== undefined) {
      AddOrgRequest.encode(message.org, writer.uint32(26).fork()).ldelim();
    }
    if (message.domainPolicy !== undefined) {
      AddCustomDomainPolicyRequest.encode(message.domainPolicy, writer.uint32(34).fork()).ldelim();
    }
    if (message.labelPolicy !== undefined) {
      AddCustomLabelPolicyRequest.encode(message.labelPolicy, writer.uint32(42).fork()).ldelim();
    }
    if (message.lockoutPolicy !== undefined) {
      AddCustomLockoutPolicyRequest.encode(message.lockoutPolicy, writer.uint32(50).fork()).ldelim();
    }
    if (message.loginPolicy !== undefined) {
      AddCustomLoginPolicyRequest.encode(message.loginPolicy, writer.uint32(58).fork()).ldelim();
    }
    if (message.passwordComplexityPolicy !== undefined) {
      AddCustomPasswordComplexityPolicyRequest.encode(message.passwordComplexityPolicy, writer.uint32(66).fork())
        .ldelim();
    }
    if (message.privacyPolicy !== undefined) {
      AddCustomPrivacyPolicyRequest.encode(message.privacyPolicy, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.projects) {
      DataProject.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.projectRoles) {
      AddProjectRoleRequest.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.apiApps) {
      DataAPIApplication.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.oidcApps) {
      DataOIDCApplication.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    for (const v of message.humanUsers) {
      DataHumanUser.encode(v!, writer.uint32(114).fork()).ldelim();
    }
    for (const v of message.machineUsers) {
      DataMachineUser.encode(v!, writer.uint32(122).fork()).ldelim();
    }
    for (const v of message.triggerActions) {
      SetTriggerActionsRequest.encode(v!, writer.uint32(130).fork()).ldelim();
    }
    for (const v of message.actions) {
      DataAction.encode(v!, writer.uint32(138).fork()).ldelim();
    }
    for (const v of message.projectGrants) {
      DataProjectGrant.encode(v!, writer.uint32(146).fork()).ldelim();
    }
    for (const v of message.userGrants) {
      AddUserGrantRequest.encode(v!, writer.uint32(154).fork()).ldelim();
    }
    for (const v of message.orgMembers) {
      AddOrgMemberRequest.encode(v!, writer.uint32(162).fork()).ldelim();
    }
    for (const v of message.projectMembers) {
      AddProjectMemberRequest.encode(v!, writer.uint32(170).fork()).ldelim();
    }
    for (const v of message.projectGrantMembers) {
      AddProjectGrantMemberRequest.encode(v!, writer.uint32(178).fork()).ldelim();
    }
    for (const v of message.userMetadata) {
      SetUserMetadataRequest.encode(v!, writer.uint32(186).fork()).ldelim();
    }
    for (const v of message.loginTexts) {
      SetCustomLoginTextsRequest2.encode(v!, writer.uint32(194).fork()).ldelim();
    }
    for (const v of message.initMessages) {
      SetCustomInitMessageTextRequest.encode(v!, writer.uint32(202).fork()).ldelim();
    }
    for (const v of message.passwordResetMessages) {
      SetCustomPasswordResetMessageTextRequest.encode(v!, writer.uint32(210).fork()).ldelim();
    }
    for (const v of message.verifyEmailMessages) {
      SetCustomVerifyEmailMessageTextRequest.encode(v!, writer.uint32(218).fork()).ldelim();
    }
    for (const v of message.verifyPhoneMessages) {
      SetCustomVerifyPhoneMessageTextRequest.encode(v!, writer.uint32(226).fork()).ldelim();
    }
    for (const v of message.domainClaimedMessages) {
      SetCustomDomainClaimedMessageTextRequest.encode(v!, writer.uint32(234).fork()).ldelim();
    }
    for (const v of message.passwordlessRegistrationMessages) {
      SetCustomPasswordlessRegistrationMessageTextRequest.encode(v!, writer.uint32(242).fork()).ldelim();
    }
    for (const v of message.oidcIdps) {
      DataOIDCIDP.encode(v!, writer.uint32(250).fork()).ldelim();
    }
    for (const v of message.jwtIdps) {
      DataJWTIDP.encode(v!, writer.uint32(258).fork()).ldelim();
    }
    for (const v of message.userLinks) {
      IDPUserLink.encode(v!, writer.uint32(266).fork()).ldelim();
    }
    for (const v of message.domains) {
      Domain3.encode(v!, writer.uint32(274).fork()).ldelim();
    }
    for (const v of message.appKeys) {
      DataAppKey.encode(v!, writer.uint32(282).fork()).ldelim();
    }
    for (const v of message.machineKeys) {
      DataMachineKey.encode(v!, writer.uint32(290).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataOrg {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataOrg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        case 3:
          message.org = AddOrgRequest.decode(reader, reader.uint32());
          break;
        case 4:
          message.domainPolicy = AddCustomDomainPolicyRequest.decode(reader, reader.uint32());
          break;
        case 5:
          message.labelPolicy = AddCustomLabelPolicyRequest.decode(reader, reader.uint32());
          break;
        case 6:
          message.lockoutPolicy = AddCustomLockoutPolicyRequest.decode(reader, reader.uint32());
          break;
        case 7:
          message.loginPolicy = AddCustomLoginPolicyRequest.decode(reader, reader.uint32());
          break;
        case 8:
          message.passwordComplexityPolicy = AddCustomPasswordComplexityPolicyRequest.decode(reader, reader.uint32());
          break;
        case 9:
          message.privacyPolicy = AddCustomPrivacyPolicyRequest.decode(reader, reader.uint32());
          break;
        case 10:
          message.projects.push(DataProject.decode(reader, reader.uint32()));
          break;
        case 11:
          message.projectRoles.push(AddProjectRoleRequest.decode(reader, reader.uint32()));
          break;
        case 12:
          message.apiApps.push(DataAPIApplication.decode(reader, reader.uint32()));
          break;
        case 13:
          message.oidcApps.push(DataOIDCApplication.decode(reader, reader.uint32()));
          break;
        case 14:
          message.humanUsers.push(DataHumanUser.decode(reader, reader.uint32()));
          break;
        case 15:
          message.machineUsers.push(DataMachineUser.decode(reader, reader.uint32()));
          break;
        case 16:
          message.triggerActions.push(SetTriggerActionsRequest.decode(reader, reader.uint32()));
          break;
        case 17:
          message.actions.push(DataAction.decode(reader, reader.uint32()));
          break;
        case 18:
          message.projectGrants.push(DataProjectGrant.decode(reader, reader.uint32()));
          break;
        case 19:
          message.userGrants.push(AddUserGrantRequest.decode(reader, reader.uint32()));
          break;
        case 20:
          message.orgMembers.push(AddOrgMemberRequest.decode(reader, reader.uint32()));
          break;
        case 21:
          message.projectMembers.push(AddProjectMemberRequest.decode(reader, reader.uint32()));
          break;
        case 22:
          message.projectGrantMembers.push(AddProjectGrantMemberRequest.decode(reader, reader.uint32()));
          break;
        case 23:
          message.userMetadata.push(SetUserMetadataRequest.decode(reader, reader.uint32()));
          break;
        case 24:
          message.loginTexts.push(SetCustomLoginTextsRequest2.decode(reader, reader.uint32()));
          break;
        case 25:
          message.initMessages.push(SetCustomInitMessageTextRequest.decode(reader, reader.uint32()));
          break;
        case 26:
          message.passwordResetMessages.push(SetCustomPasswordResetMessageTextRequest.decode(reader, reader.uint32()));
          break;
        case 27:
          message.verifyEmailMessages.push(SetCustomVerifyEmailMessageTextRequest.decode(reader, reader.uint32()));
          break;
        case 28:
          message.verifyPhoneMessages.push(SetCustomVerifyPhoneMessageTextRequest.decode(reader, reader.uint32()));
          break;
        case 29:
          message.domainClaimedMessages.push(SetCustomDomainClaimedMessageTextRequest.decode(reader, reader.uint32()));
          break;
        case 30:
          message.passwordlessRegistrationMessages.push(
            SetCustomPasswordlessRegistrationMessageTextRequest.decode(reader, reader.uint32()),
          );
          break;
        case 31:
          message.oidcIdps.push(DataOIDCIDP.decode(reader, reader.uint32()));
          break;
        case 32:
          message.jwtIdps.push(DataJWTIDP.decode(reader, reader.uint32()));
          break;
        case 33:
          message.userLinks.push(IDPUserLink.decode(reader, reader.uint32()));
          break;
        case 34:
          message.domains.push(Domain3.decode(reader, reader.uint32()));
          break;
        case 35:
          message.appKeys.push(DataAppKey.decode(reader, reader.uint32()));
          break;
        case 36:
          message.machineKeys.push(DataMachineKey.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DataOrg {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      org: isSet(object.org) ? AddOrgRequest.fromJSON(object.org) : undefined,
      domainPolicy: isSet(object.domainPolicy) ? AddCustomDomainPolicyRequest.fromJSON(object.domainPolicy) : undefined,
      labelPolicy: isSet(object.labelPolicy) ? AddCustomLabelPolicyRequest.fromJSON(object.labelPolicy) : undefined,
      lockoutPolicy: isSet(object.lockoutPolicy)
        ? AddCustomLockoutPolicyRequest.fromJSON(object.lockoutPolicy)
        : undefined,
      loginPolicy: isSet(object.loginPolicy) ? AddCustomLoginPolicyRequest.fromJSON(object.loginPolicy) : undefined,
      passwordComplexityPolicy: isSet(object.passwordComplexityPolicy)
        ? AddCustomPasswordComplexityPolicyRequest.fromJSON(object.passwordComplexityPolicy)
        : undefined,
      privacyPolicy: isSet(object.privacyPolicy)
        ? AddCustomPrivacyPolicyRequest.fromJSON(object.privacyPolicy)
        : undefined,
      projects: Array.isArray(object?.projects) ? object.projects.map((e: any) => DataProject.fromJSON(e)) : [],
      projectRoles: Array.isArray(object?.projectRoles)
        ? object.projectRoles.map((e: any) => AddProjectRoleRequest.fromJSON(e))
        : [],
      apiApps: Array.isArray(object?.apiApps) ? object.apiApps.map((e: any) => DataAPIApplication.fromJSON(e)) : [],
      oidcApps: Array.isArray(object?.oidcApps) ? object.oidcApps.map((e: any) => DataOIDCApplication.fromJSON(e)) : [],
      humanUsers: Array.isArray(object?.humanUsers) ? object.humanUsers.map((e: any) => DataHumanUser.fromJSON(e)) : [],
      machineUsers: Array.isArray(object?.machineUsers)
        ? object.machineUsers.map((e: any) => DataMachineUser.fromJSON(e))
        : [],
      triggerActions: Array.isArray(object?.triggerActions)
        ? object.triggerActions.map((e: any) => SetTriggerActionsRequest.fromJSON(e))
        : [],
      actions: Array.isArray(object?.actions) ? object.actions.map((e: any) => DataAction.fromJSON(e)) : [],
      projectGrants: Array.isArray(object?.projectGrants)
        ? object.projectGrants.map((e: any) => DataProjectGrant.fromJSON(e))
        : [],
      userGrants: Array.isArray(object?.userGrants)
        ? object.userGrants.map((e: any) => AddUserGrantRequest.fromJSON(e))
        : [],
      orgMembers: Array.isArray(object?.orgMembers)
        ? object.orgMembers.map((e: any) => AddOrgMemberRequest.fromJSON(e))
        : [],
      projectMembers: Array.isArray(object?.projectMembers)
        ? object.projectMembers.map((e: any) => AddProjectMemberRequest.fromJSON(e))
        : [],
      projectGrantMembers: Array.isArray(object?.projectGrantMembers)
        ? object.projectGrantMembers.map((e: any) => AddProjectGrantMemberRequest.fromJSON(e))
        : [],
      userMetadata: Array.isArray(object?.userMetadata)
        ? object.userMetadata.map((e: any) => SetUserMetadataRequest.fromJSON(e))
        : [],
      loginTexts: Array.isArray(object?.loginTexts)
        ? object.loginTexts.map((e: any) => SetCustomLoginTextsRequest.fromJSON(e))
        : [],
      initMessages: Array.isArray(object?.initMessages)
        ? object.initMessages.map((e: any) => SetCustomInitMessageTextRequest.fromJSON(e))
        : [],
      passwordResetMessages: Array.isArray(object?.passwordResetMessages)
        ? object.passwordResetMessages.map((e: any) => SetCustomPasswordResetMessageTextRequest.fromJSON(e))
        : [],
      verifyEmailMessages: Array.isArray(object?.verifyEmailMessages)
        ? object.verifyEmailMessages.map((e: any) => SetCustomVerifyEmailMessageTextRequest.fromJSON(e))
        : [],
      verifyPhoneMessages: Array.isArray(object?.verifyPhoneMessages)
        ? object.verifyPhoneMessages.map((e: any) => SetCustomVerifyPhoneMessageTextRequest.fromJSON(e))
        : [],
      domainClaimedMessages: Array.isArray(object?.domainClaimedMessages)
        ? object.domainClaimedMessages.map((e: any) => SetCustomDomainClaimedMessageTextRequest.fromJSON(e))
        : [],
      passwordlessRegistrationMessages: Array.isArray(object?.passwordlessRegistrationMessages)
        ? object.passwordlessRegistrationMessages.map((e: any) =>
          SetCustomPasswordlessRegistrationMessageTextRequest.fromJSON(e)
        )
        : [],
      oidcIdps: Array.isArray(object?.oidcIdps)
        ? object.oidcIdps.map((e: any) => DataOIDCIDP.fromJSON(e))
        : [],
      jwtIdps: Array.isArray(object?.jwtIdps)
        ? object.jwtIdps.map((e: any) => DataJWTIDP.fromJSON(e))
        : [],
      userLinks: Array.isArray(object?.userLinks) ? object.userLinks.map((e: any) => IDPUserLink.fromJSON(e)) : [],
      domains: Array.isArray(object?.domains) ? object.domains.map((e: any) => Domain.fromJSON(e)) : [],
      appKeys: Array.isArray(object?.appKeys) ? object.appKeys.map((e: any) => DataAppKey.fromJSON(e)) : [],
      machineKeys: Array.isArray(object?.machineKeys)
        ? object.machineKeys.map((e: any) => DataMachineKey.fromJSON(e))
        : [],
    };
  },

  toJSON(message: DataOrg): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.org !== undefined && (obj.org = message.org ? AddOrgRequest.toJSON(message.org) : undefined);
    message.domainPolicy !== undefined &&
      (obj.domainPolicy = message.domainPolicy ? AddCustomDomainPolicyRequest.toJSON(message.domainPolicy) : undefined);
    message.labelPolicy !== undefined &&
      (obj.labelPolicy = message.labelPolicy ? AddCustomLabelPolicyRequest.toJSON(message.labelPolicy) : undefined);
    message.lockoutPolicy !== undefined && (obj.lockoutPolicy = message.lockoutPolicy
      ? AddCustomLockoutPolicyRequest.toJSON(message.lockoutPolicy)
      : undefined);
    message.loginPolicy !== undefined &&
      (obj.loginPolicy = message.loginPolicy ? AddCustomLoginPolicyRequest.toJSON(message.loginPolicy) : undefined);
    message.passwordComplexityPolicy !== undefined && (obj.passwordComplexityPolicy = message.passwordComplexityPolicy
      ? AddCustomPasswordComplexityPolicyRequest.toJSON(message.passwordComplexityPolicy)
      : undefined);
    message.privacyPolicy !== undefined && (obj.privacyPolicy = message.privacyPolicy
      ? AddCustomPrivacyPolicyRequest.toJSON(message.privacyPolicy)
      : undefined);
    if (message.projects) {
      obj.projects = message.projects.map((e) => e ? DataProject.toJSON(e) : undefined);
    } else {
      obj.projects = [];
    }
    if (message.projectRoles) {
      obj.projectRoles = message.projectRoles.map((e) => e ? AddProjectRoleRequest.toJSON(e) : undefined);
    } else {
      obj.projectRoles = [];
    }
    if (message.apiApps) {
      obj.apiApps = message.apiApps.map((e) => e ? DataAPIApplication.toJSON(e) : undefined);
    } else {
      obj.apiApps = [];
    }
    if (message.oidcApps) {
      obj.oidcApps = message.oidcApps.map((e) => e ? DataOIDCApplication.toJSON(e) : undefined);
    } else {
      obj.oidcApps = [];
    }
    if (message.humanUsers) {
      obj.humanUsers = message.humanUsers.map((e) => e ? DataHumanUser.toJSON(e) : undefined);
    } else {
      obj.humanUsers = [];
    }
    if (message.machineUsers) {
      obj.machineUsers = message.machineUsers.map((e) => e ? DataMachineUser.toJSON(e) : undefined);
    } else {
      obj.machineUsers = [];
    }
    if (message.triggerActions) {
      obj.triggerActions = message.triggerActions.map((e) => e ? SetTriggerActionsRequest.toJSON(e) : undefined);
    } else {
      obj.triggerActions = [];
    }
    if (message.actions) {
      obj.actions = message.actions.map((e) => e ? DataAction.toJSON(e) : undefined);
    } else {
      obj.actions = [];
    }
    if (message.projectGrants) {
      obj.projectGrants = message.projectGrants.map((e) => e ? DataProjectGrant.toJSON(e) : undefined);
    } else {
      obj.projectGrants = [];
    }
    if (message.userGrants) {
      obj.userGrants = message.userGrants.map((e) => e ? AddUserGrantRequest.toJSON(e) : undefined);
    } else {
      obj.userGrants = [];
    }
    if (message.orgMembers) {
      obj.orgMembers = message.orgMembers.map((e) => e ? AddOrgMemberRequest.toJSON(e) : undefined);
    } else {
      obj.orgMembers = [];
    }
    if (message.projectMembers) {
      obj.projectMembers = message.projectMembers.map((e) => e ? AddProjectMemberRequest.toJSON(e) : undefined);
    } else {
      obj.projectMembers = [];
    }
    if (message.projectGrantMembers) {
      obj.projectGrantMembers = message.projectGrantMembers.map((e) =>
        e ? AddProjectGrantMemberRequest.toJSON(e) : undefined
      );
    } else {
      obj.projectGrantMembers = [];
    }
    if (message.userMetadata) {
      obj.userMetadata = message.userMetadata.map((e) => e ? SetUserMetadataRequest.toJSON(e) : undefined);
    } else {
      obj.userMetadata = [];
    }
    if (message.loginTexts) {
      obj.loginTexts = message.loginTexts.map((e) => e ? SetCustomLoginTextsRequest2.toJSON(e) : undefined);
    } else {
      obj.loginTexts = [];
    }
    if (message.initMessages) {
      obj.initMessages = message.initMessages.map((e) => e ? SetCustomInitMessageTextRequest.toJSON(e) : undefined);
    } else {
      obj.initMessages = [];
    }
    if (message.passwordResetMessages) {
      obj.passwordResetMessages = message.passwordResetMessages.map((e) =>
        e ? SetCustomPasswordResetMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.passwordResetMessages = [];
    }
    if (message.verifyEmailMessages) {
      obj.verifyEmailMessages = message.verifyEmailMessages.map((e) =>
        e ? SetCustomVerifyEmailMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.verifyEmailMessages = [];
    }
    if (message.verifyPhoneMessages) {
      obj.verifyPhoneMessages = message.verifyPhoneMessages.map((e) =>
        e ? SetCustomVerifyPhoneMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.verifyPhoneMessages = [];
    }
    if (message.domainClaimedMessages) {
      obj.domainClaimedMessages = message.domainClaimedMessages.map((e) =>
        e ? SetCustomDomainClaimedMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.domainClaimedMessages = [];
    }
    if (message.passwordlessRegistrationMessages) {
      obj.passwordlessRegistrationMessages = message.passwordlessRegistrationMessages.map((e) =>
        e ? SetCustomPasswordlessRegistrationMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.passwordlessRegistrationMessages = [];
    }
    if (message.oidcIdps) {
      obj.oidcIdps = message.oidcIdps.map((e) => e ? DataOIDCIDP.toJSON(e) : undefined);
    } else {
      obj.oidcIdps = [];
    }
    if (message.jwtIdps) {
      obj.jwtIdps = message.jwtIdps.map((e) => e ? DataJWTIDP.toJSON(e) : undefined);
    } else {
      obj.jwtIdps = [];
    }
    if (message.userLinks) {
      obj.userLinks = message.userLinks.map((e) => e ? IDPUserLink.toJSON(e) : undefined);
    } else {
      obj.userLinks = [];
    }
    if (message.domains) {
      obj.domains = message.domains.map((e) => e ? Domain3.toJSON(e) : undefined);
    } else {
      obj.domains = [];
    }
    if (message.appKeys) {
      obj.appKeys = message.appKeys.map((e) => e ? DataAppKey.toJSON(e) : undefined);
    } else {
      obj.appKeys = [];
    }
    if (message.machineKeys) {
      obj.machineKeys = message.machineKeys.map((e) => e ? DataMachineKey.toJSON(e) : undefined);
    } else {
      obj.machineKeys = [];
    }
    return obj;
  },

  create(base?: DeepPartial<DataOrg>): DataOrg {
    return DataOrg.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataOrg>): DataOrg {
    const message = createBaseDataOrg();
    message.orgId = object.orgId ?? "";
    message.org = (object.org !== undefined && object.org !== null) ? AddOrgRequest.fromPartial(object.org) : undefined;
    message.domainPolicy = (object.domainPolicy !== undefined && object.domainPolicy !== null)
      ? AddCustomDomainPolicyRequest.fromPartial(object.domainPolicy)
      : undefined;
    message.labelPolicy = (object.labelPolicy !== undefined && object.labelPolicy !== null)
      ? AddCustomLabelPolicyRequest.fromPartial(object.labelPolicy)
      : undefined;
    message.lockoutPolicy = (object.lockoutPolicy !== undefined && object.lockoutPolicy !== null)
      ? AddCustomLockoutPolicyRequest.fromPartial(object.lockoutPolicy)
      : undefined;
    message.loginPolicy = (object.loginPolicy !== undefined && object.loginPolicy !== null)
      ? AddCustomLoginPolicyRequest.fromPartial(object.loginPolicy)
      : undefined;
    message.passwordComplexityPolicy =
      (object.passwordComplexityPolicy !== undefined && object.passwordComplexityPolicy !== null)
        ? AddCustomPasswordComplexityPolicyRequest.fromPartial(object.passwordComplexityPolicy)
        : undefined;
    message.privacyPolicy = (object.privacyPolicy !== undefined && object.privacyPolicy !== null)
      ? AddCustomPrivacyPolicyRequest.fromPartial(object.privacyPolicy)
      : undefined;
    message.projects = object.projects?.map((e) => DataProject.fromPartial(e)) || [];
    message.projectRoles = object.projectRoles?.map((e) => AddProjectRoleRequest.fromPartial(e)) || [];
    message.apiApps = object.apiApps?.map((e) => DataAPIApplication.fromPartial(e)) || [];
    message.oidcApps = object.oidcApps?.map((e) => DataOIDCApplication.fromPartial(e)) || [];
    message.humanUsers = object.humanUsers?.map((e) => DataHumanUser.fromPartial(e)) || [];
    message.machineUsers = object.machineUsers?.map((e) => DataMachineUser.fromPartial(e)) || [];
    message.triggerActions = object.triggerActions?.map((e) => SetTriggerActionsRequest.fromPartial(e)) || [];
    message.actions = object.actions?.map((e) => DataAction.fromPartial(e)) || [];
    message.projectGrants = object.projectGrants?.map((e) => DataProjectGrant.fromPartial(e)) || [];
    message.userGrants = object.userGrants?.map((e) => AddUserGrantRequest.fromPartial(e)) || [];
    message.orgMembers = object.orgMembers?.map((e) => AddOrgMemberRequest.fromPartial(e)) || [];
    message.projectMembers = object.projectMembers?.map((e) => AddProjectMemberRequest.fromPartial(e)) || [];
    message.projectGrantMembers = object.projectGrantMembers?.map((e) => AddProjectGrantMemberRequest.fromPartial(e)) ||
      [];
    message.userMetadata = object.userMetadata?.map((e) => SetUserMetadataRequest.fromPartial(e)) || [];
    message.loginTexts = object.loginTexts?.map((e) => SetCustomLoginTextsRequest2.fromPartial(e)) || [];
    message.initMessages = object.initMessages?.map((e) => SetCustomInitMessageTextRequest.fromPartial(e)) || [];
    message.passwordResetMessages =
      object.passwordResetMessages?.map((e) => SetCustomPasswordResetMessageTextRequest.fromPartial(e)) || [];
    message.verifyEmailMessages =
      object.verifyEmailMessages?.map((e) => SetCustomVerifyEmailMessageTextRequest.fromPartial(e)) || [];
    message.verifyPhoneMessages =
      object.verifyPhoneMessages?.map((e) => SetCustomVerifyPhoneMessageTextRequest.fromPartial(e)) || [];
    message.domainClaimedMessages =
      object.domainClaimedMessages?.map((e) => SetCustomDomainClaimedMessageTextRequest.fromPartial(e)) || [];
    message.passwordlessRegistrationMessages =
      object.passwordlessRegistrationMessages?.map((e) =>
        SetCustomPasswordlessRegistrationMessageTextRequest.fromPartial(e)
      ) || [];
    message.oidcIdps = object.oidcIdps?.map((e) => DataOIDCIDP.fromPartial(e)) || [];
    message.jwtIdps = object.jwtIdps?.map((e) => DataJWTIDP.fromPartial(e)) || [];
    message.userLinks = object.userLinks?.map((e) => IDPUserLink.fromPartial(e)) || [];
    message.domains = object.domains?.map((e) => Domain3.fromPartial(e)) || [];
    message.appKeys = object.appKeys?.map((e) => DataAppKey.fromPartial(e)) || [];
    message.machineKeys = object.machineKeys?.map((e) => DataMachineKey.fromPartial(e)) || [];
    return message;
  },
};

function createBaseImportDataResponse(): ImportDataResponse {
  return { errors: [], success: undefined };
}

export const ImportDataResponse = {
  encode(message: ImportDataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.errors) {
      ImportDataError.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.success !== undefined) {
      ImportDataSuccess.encode(message.success, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.errors.push(ImportDataError.decode(reader, reader.uint32()));
          break;
        case 2:
          message.success = ImportDataSuccess.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataResponse {
    return {
      errors: Array.isArray(object?.errors) ? object.errors.map((e: any) => ImportDataError.fromJSON(e)) : [],
      success: isSet(object.success) ? ImportDataSuccess.fromJSON(object.success) : undefined,
    };
  },

  toJSON(message: ImportDataResponse): unknown {
    const obj: any = {};
    if (message.errors) {
      obj.errors = message.errors.map((e) => e ? ImportDataError.toJSON(e) : undefined);
    } else {
      obj.errors = [];
    }
    message.success !== undefined &&
      (obj.success = message.success ? ImportDataSuccess.toJSON(message.success) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ImportDataResponse>): ImportDataResponse {
    return ImportDataResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataResponse>): ImportDataResponse {
    const message = createBaseImportDataResponse();
    message.errors = object.errors?.map((e) => ImportDataError.fromPartial(e)) || [];
    message.success = (object.success !== undefined && object.success !== null)
      ? ImportDataSuccess.fromPartial(object.success)
      : undefined;
    return message;
  },
};

function createBaseImportDataError(): ImportDataError {
  return { type: "", id: "", message: "" };
}

export const ImportDataError = {
  encode(message: ImportDataError, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    if (message.message !== "") {
      writer.uint32(26).string(message.message);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataError {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataError();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.type = reader.string();
          break;
        case 2:
          message.id = reader.string();
          break;
        case 3:
          message.message = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataError {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      id: isSet(object.id) ? String(object.id) : "",
      message: isSet(object.message) ? String(object.message) : "",
    };
  },

  toJSON(message: ImportDataError): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.id !== undefined && (obj.id = message.id);
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },

  create(base?: DeepPartial<ImportDataError>): ImportDataError {
    return ImportDataError.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataError>): ImportDataError {
    const message = createBaseImportDataError();
    message.type = object.type ?? "";
    message.id = object.id ?? "";
    message.message = object.message ?? "";
    return message;
  },
};

function createBaseImportDataSuccess(): ImportDataSuccess {
  return { orgs: [] };
}

export const ImportDataSuccess = {
  encode(message: ImportDataSuccess, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.orgs) {
      ImportDataSuccessOrg.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataSuccess {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataSuccess();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgs.push(ImportDataSuccessOrg.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataSuccess {
    return { orgs: Array.isArray(object?.orgs) ? object.orgs.map((e: any) => ImportDataSuccessOrg.fromJSON(e)) : [] };
  },

  toJSON(message: ImportDataSuccess): unknown {
    const obj: any = {};
    if (message.orgs) {
      obj.orgs = message.orgs.map((e) => e ? ImportDataSuccessOrg.toJSON(e) : undefined);
    } else {
      obj.orgs = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ImportDataSuccess>): ImportDataSuccess {
    return ImportDataSuccess.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataSuccess>): ImportDataSuccess {
    const message = createBaseImportDataSuccess();
    message.orgs = object.orgs?.map((e) => ImportDataSuccessOrg.fromPartial(e)) || [];
    return message;
  },
};

function createBaseImportDataSuccessOrg(): ImportDataSuccessOrg {
  return {
    orgId: "",
    projectIds: [],
    projectRoles: [],
    oidcAppIds: [],
    apiAppIds: [],
    humanUserIds: [],
    machineUserIds: [],
    actionIds: [],
    triggerActions: [],
    projectGrants: [],
    userGrants: [],
    orgMembers: [],
    projectMembers: [],
    projectGrantMembers: [],
    oidcIpds: [],
    jwtIdps: [],
    idpLinks: [],
    userLinks: [],
    userMetadata: [],
    domains: [],
    appKeys: [],
    machineKeys: [],
  };
}

export const ImportDataSuccessOrg = {
  encode(message: ImportDataSuccessOrg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    for (const v of message.projectIds) {
      writer.uint32(18).string(v!);
    }
    for (const v of message.projectRoles) {
      writer.uint32(26).string(v!);
    }
    for (const v of message.oidcAppIds) {
      writer.uint32(34).string(v!);
    }
    for (const v of message.apiAppIds) {
      writer.uint32(42).string(v!);
    }
    for (const v of message.humanUserIds) {
      writer.uint32(50).string(v!);
    }
    for (const v of message.machineUserIds) {
      writer.uint32(58).string(v!);
    }
    for (const v of message.actionIds) {
      writer.uint32(66).string(v!);
    }
    for (const v of message.triggerActions) {
      SetTriggerActionsRequest.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.projectGrants) {
      ImportDataSuccessProjectGrant.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.userGrants) {
      ImportDataSuccessUserGrant.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.orgMembers) {
      writer.uint32(98).string(v!);
    }
    for (const v of message.projectMembers) {
      ImportDataSuccessProjectMember.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    for (const v of message.projectGrantMembers) {
      ImportDataSuccessProjectGrantMember.encode(v!, writer.uint32(114).fork()).ldelim();
    }
    for (const v of message.oidcIpds) {
      writer.uint32(122).string(v!);
    }
    for (const v of message.jwtIdps) {
      writer.uint32(130).string(v!);
    }
    for (const v of message.idpLinks) {
      writer.uint32(138).string(v!);
    }
    for (const v of message.userLinks) {
      ImportDataSuccessUserLinks.encode(v!, writer.uint32(146).fork()).ldelim();
    }
    for (const v of message.userMetadata) {
      ImportDataSuccessUserMetadata.encode(v!, writer.uint32(154).fork()).ldelim();
    }
    for (const v of message.domains) {
      writer.uint32(162).string(v!);
    }
    for (const v of message.appKeys) {
      writer.uint32(170).string(v!);
    }
    for (const v of message.machineKeys) {
      writer.uint32(178).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataSuccessOrg {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataSuccessOrg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        case 2:
          message.projectIds.push(reader.string());
          break;
        case 3:
          message.projectRoles.push(reader.string());
          break;
        case 4:
          message.oidcAppIds.push(reader.string());
          break;
        case 5:
          message.apiAppIds.push(reader.string());
          break;
        case 6:
          message.humanUserIds.push(reader.string());
          break;
        case 7:
          message.machineUserIds.push(reader.string());
          break;
        case 8:
          message.actionIds.push(reader.string());
          break;
        case 9:
          message.triggerActions.push(SetTriggerActionsRequest.decode(reader, reader.uint32()));
          break;
        case 10:
          message.projectGrants.push(ImportDataSuccessProjectGrant.decode(reader, reader.uint32()));
          break;
        case 11:
          message.userGrants.push(ImportDataSuccessUserGrant.decode(reader, reader.uint32()));
          break;
        case 12:
          message.orgMembers.push(reader.string());
          break;
        case 13:
          message.projectMembers.push(ImportDataSuccessProjectMember.decode(reader, reader.uint32()));
          break;
        case 14:
          message.projectGrantMembers.push(ImportDataSuccessProjectGrantMember.decode(reader, reader.uint32()));
          break;
        case 15:
          message.oidcIpds.push(reader.string());
          break;
        case 16:
          message.jwtIdps.push(reader.string());
          break;
        case 17:
          message.idpLinks.push(reader.string());
          break;
        case 18:
          message.userLinks.push(ImportDataSuccessUserLinks.decode(reader, reader.uint32()));
          break;
        case 19:
          message.userMetadata.push(ImportDataSuccessUserMetadata.decode(reader, reader.uint32()));
          break;
        case 20:
          message.domains.push(reader.string());
          break;
        case 21:
          message.appKeys.push(reader.string());
          break;
        case 22:
          message.machineKeys.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataSuccessOrg {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      projectIds: Array.isArray(object?.projectIds) ? object.projectIds.map((e: any) => String(e)) : [],
      projectRoles: Array.isArray(object?.projectRoles) ? object.projectRoles.map((e: any) => String(e)) : [],
      oidcAppIds: Array.isArray(object?.oidcAppIds) ? object.oidcAppIds.map((e: any) => String(e)) : [],
      apiAppIds: Array.isArray(object?.apiAppIds) ? object.apiAppIds.map((e: any) => String(e)) : [],
      humanUserIds: Array.isArray(object?.humanUserIds) ? object.humanUserIds.map((e: any) => String(e)) : [],
      machineUserIds: Array.isArray(object?.machineUserIds) ? object.machineUserIds.map((e: any) => String(e)) : [],
      actionIds: Array.isArray(object?.actionIds) ? object.actionIds.map((e: any) => String(e)) : [],
      triggerActions: Array.isArray(object?.triggerActions)
        ? object.triggerActions.map((e: any) => SetTriggerActionsRequest.fromJSON(e))
        : [],
      projectGrants: Array.isArray(object?.projectGrants)
        ? object.projectGrants.map((e: any) => ImportDataSuccessProjectGrant.fromJSON(e))
        : [],
      userGrants: Array.isArray(object?.userGrants)
        ? object.userGrants.map((e: any) => ImportDataSuccessUserGrant.fromJSON(e))
        : [],
      orgMembers: Array.isArray(object?.orgMembers) ? object.orgMembers.map((e: any) => String(e)) : [],
      projectMembers: Array.isArray(object?.projectMembers)
        ? object.projectMembers.map((e: any) => ImportDataSuccessProjectMember.fromJSON(e))
        : [],
      projectGrantMembers: Array.isArray(object?.projectGrantMembers)
        ? object.projectGrantMembers.map((e: any) => ImportDataSuccessProjectGrantMember.fromJSON(e))
        : [],
      oidcIpds: Array.isArray(object?.oidcIpds) ? object.oidcIpds.map((e: any) => String(e)) : [],
      jwtIdps: Array.isArray(object?.jwtIdps) ? object.jwtIdps.map((e: any) => String(e)) : [],
      idpLinks: Array.isArray(object?.idpLinks) ? object.idpLinks.map((e: any) => String(e)) : [],
      userLinks: Array.isArray(object?.userLinks)
        ? object.userLinks.map((e: any) => ImportDataSuccessUserLinks.fromJSON(e))
        : [],
      userMetadata: Array.isArray(object?.userMetadata)
        ? object.userMetadata.map((e: any) => ImportDataSuccessUserMetadata.fromJSON(e))
        : [],
      domains: Array.isArray(object?.domains) ? object.domains.map((e: any) => String(e)) : [],
      appKeys: Array.isArray(object?.appKeys) ? object.appKeys.map((e: any) => String(e)) : [],
      machineKeys: Array.isArray(object?.machineKeys) ? object.machineKeys.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: ImportDataSuccessOrg): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    if (message.projectIds) {
      obj.projectIds = message.projectIds.map((e) => e);
    } else {
      obj.projectIds = [];
    }
    if (message.projectRoles) {
      obj.projectRoles = message.projectRoles.map((e) => e);
    } else {
      obj.projectRoles = [];
    }
    if (message.oidcAppIds) {
      obj.oidcAppIds = message.oidcAppIds.map((e) => e);
    } else {
      obj.oidcAppIds = [];
    }
    if (message.apiAppIds) {
      obj.apiAppIds = message.apiAppIds.map((e) => e);
    } else {
      obj.apiAppIds = [];
    }
    if (message.humanUserIds) {
      obj.humanUserIds = message.humanUserIds.map((e) => e);
    } else {
      obj.humanUserIds = [];
    }
    if (message.machineUserIds) {
      obj.machineUserIds = message.machineUserIds.map((e) => e);
    } else {
      obj.machineUserIds = [];
    }
    if (message.actionIds) {
      obj.actionIds = message.actionIds.map((e) => e);
    } else {
      obj.actionIds = [];
    }
    if (message.triggerActions) {
      obj.triggerActions = message.triggerActions.map((e) => e ? SetTriggerActionsRequest.toJSON(e) : undefined);
    } else {
      obj.triggerActions = [];
    }
    if (message.projectGrants) {
      obj.projectGrants = message.projectGrants.map((e) => e ? ImportDataSuccessProjectGrant.toJSON(e) : undefined);
    } else {
      obj.projectGrants = [];
    }
    if (message.userGrants) {
      obj.userGrants = message.userGrants.map((e) => e ? ImportDataSuccessUserGrant.toJSON(e) : undefined);
    } else {
      obj.userGrants = [];
    }
    if (message.orgMembers) {
      obj.orgMembers = message.orgMembers.map((e) => e);
    } else {
      obj.orgMembers = [];
    }
    if (message.projectMembers) {
      obj.projectMembers = message.projectMembers.map((e) => e ? ImportDataSuccessProjectMember.toJSON(e) : undefined);
    } else {
      obj.projectMembers = [];
    }
    if (message.projectGrantMembers) {
      obj.projectGrantMembers = message.projectGrantMembers.map((e) =>
        e ? ImportDataSuccessProjectGrantMember.toJSON(e) : undefined
      );
    } else {
      obj.projectGrantMembers = [];
    }
    if (message.oidcIpds) {
      obj.oidcIpds = message.oidcIpds.map((e) => e);
    } else {
      obj.oidcIpds = [];
    }
    if (message.jwtIdps) {
      obj.jwtIdps = message.jwtIdps.map((e) => e);
    } else {
      obj.jwtIdps = [];
    }
    if (message.idpLinks) {
      obj.idpLinks = message.idpLinks.map((e) => e);
    } else {
      obj.idpLinks = [];
    }
    if (message.userLinks) {
      obj.userLinks = message.userLinks.map((e) => e ? ImportDataSuccessUserLinks.toJSON(e) : undefined);
    } else {
      obj.userLinks = [];
    }
    if (message.userMetadata) {
      obj.userMetadata = message.userMetadata.map((e) => e ? ImportDataSuccessUserMetadata.toJSON(e) : undefined);
    } else {
      obj.userMetadata = [];
    }
    if (message.domains) {
      obj.domains = message.domains.map((e) => e);
    } else {
      obj.domains = [];
    }
    if (message.appKeys) {
      obj.appKeys = message.appKeys.map((e) => e);
    } else {
      obj.appKeys = [];
    }
    if (message.machineKeys) {
      obj.machineKeys = message.machineKeys.map((e) => e);
    } else {
      obj.machineKeys = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ImportDataSuccessOrg>): ImportDataSuccessOrg {
    return ImportDataSuccessOrg.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataSuccessOrg>): ImportDataSuccessOrg {
    const message = createBaseImportDataSuccessOrg();
    message.orgId = object.orgId ?? "";
    message.projectIds = object.projectIds?.map((e) => e) || [];
    message.projectRoles = object.projectRoles?.map((e) => e) || [];
    message.oidcAppIds = object.oidcAppIds?.map((e) => e) || [];
    message.apiAppIds = object.apiAppIds?.map((e) => e) || [];
    message.humanUserIds = object.humanUserIds?.map((e) => e) || [];
    message.machineUserIds = object.machineUserIds?.map((e) => e) || [];
    message.actionIds = object.actionIds?.map((e) => e) || [];
    message.triggerActions = object.triggerActions?.map((e) => SetTriggerActionsRequest.fromPartial(e)) || [];
    message.projectGrants = object.projectGrants?.map((e) => ImportDataSuccessProjectGrant.fromPartial(e)) || [];
    message.userGrants = object.userGrants?.map((e) => ImportDataSuccessUserGrant.fromPartial(e)) || [];
    message.orgMembers = object.orgMembers?.map((e) => e) || [];
    message.projectMembers = object.projectMembers?.map((e) => ImportDataSuccessProjectMember.fromPartial(e)) || [];
    message.projectGrantMembers =
      object.projectGrantMembers?.map((e) => ImportDataSuccessProjectGrantMember.fromPartial(e)) || [];
    message.oidcIpds = object.oidcIpds?.map((e) => e) || [];
    message.jwtIdps = object.jwtIdps?.map((e) => e) || [];
    message.idpLinks = object.idpLinks?.map((e) => e) || [];
    message.userLinks = object.userLinks?.map((e) => ImportDataSuccessUserLinks.fromPartial(e)) || [];
    message.userMetadata = object.userMetadata?.map((e) => ImportDataSuccessUserMetadata.fromPartial(e)) || [];
    message.domains = object.domains?.map((e) => e) || [];
    message.appKeys = object.appKeys?.map((e) => e) || [];
    message.machineKeys = object.machineKeys?.map((e) => e) || [];
    return message;
  },
};

function createBaseImportDataSuccessProjectGrant(): ImportDataSuccessProjectGrant {
  return { grantId: "", projectId: "", orgId: "" };
}

export const ImportDataSuccessProjectGrant = {
  encode(message: ImportDataSuccessProjectGrant, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.grantId !== "") {
      writer.uint32(10).string(message.grantId);
    }
    if (message.projectId !== "") {
      writer.uint32(18).string(message.projectId);
    }
    if (message.orgId !== "") {
      writer.uint32(26).string(message.orgId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataSuccessProjectGrant {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataSuccessProjectGrant();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.grantId = reader.string();
          break;
        case 2:
          message.projectId = reader.string();
          break;
        case 3:
          message.orgId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataSuccessProjectGrant {
    return {
      grantId: isSet(object.grantId) ? String(object.grantId) : "",
      projectId: isSet(object.projectId) ? String(object.projectId) : "",
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
    };
  },

  toJSON(message: ImportDataSuccessProjectGrant): unknown {
    const obj: any = {};
    message.grantId !== undefined && (obj.grantId = message.grantId);
    message.projectId !== undefined && (obj.projectId = message.projectId);
    message.orgId !== undefined && (obj.orgId = message.orgId);
    return obj;
  },

  create(base?: DeepPartial<ImportDataSuccessProjectGrant>): ImportDataSuccessProjectGrant {
    return ImportDataSuccessProjectGrant.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataSuccessProjectGrant>): ImportDataSuccessProjectGrant {
    const message = createBaseImportDataSuccessProjectGrant();
    message.grantId = object.grantId ?? "";
    message.projectId = object.projectId ?? "";
    message.orgId = object.orgId ?? "";
    return message;
  },
};

function createBaseImportDataSuccessUserGrant(): ImportDataSuccessUserGrant {
  return { projectId: "", userId: "" };
}

export const ImportDataSuccessUserGrant = {
  encode(message: ImportDataSuccessUserGrant, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.projectId !== "") {
      writer.uint32(10).string(message.projectId);
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataSuccessUserGrant {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataSuccessUserGrant();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.projectId = reader.string();
          break;
        case 2:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataSuccessUserGrant {
    return {
      projectId: isSet(object.projectId) ? String(object.projectId) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: ImportDataSuccessUserGrant): unknown {
    const obj: any = {};
    message.projectId !== undefined && (obj.projectId = message.projectId);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<ImportDataSuccessUserGrant>): ImportDataSuccessUserGrant {
    return ImportDataSuccessUserGrant.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataSuccessUserGrant>): ImportDataSuccessUserGrant {
    const message = createBaseImportDataSuccessUserGrant();
    message.projectId = object.projectId ?? "";
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseImportDataSuccessProjectMember(): ImportDataSuccessProjectMember {
  return { projectId: "", userId: "" };
}

export const ImportDataSuccessProjectMember = {
  encode(message: ImportDataSuccessProjectMember, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.projectId !== "") {
      writer.uint32(10).string(message.projectId);
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataSuccessProjectMember {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataSuccessProjectMember();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.projectId = reader.string();
          break;
        case 2:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataSuccessProjectMember {
    return {
      projectId: isSet(object.projectId) ? String(object.projectId) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: ImportDataSuccessProjectMember): unknown {
    const obj: any = {};
    message.projectId !== undefined && (obj.projectId = message.projectId);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<ImportDataSuccessProjectMember>): ImportDataSuccessProjectMember {
    return ImportDataSuccessProjectMember.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataSuccessProjectMember>): ImportDataSuccessProjectMember {
    const message = createBaseImportDataSuccessProjectMember();
    message.projectId = object.projectId ?? "";
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseImportDataSuccessProjectGrantMember(): ImportDataSuccessProjectGrantMember {
  return { projectId: "", grantId: "", userId: "" };
}

export const ImportDataSuccessProjectGrantMember = {
  encode(message: ImportDataSuccessProjectGrantMember, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.projectId !== "") {
      writer.uint32(10).string(message.projectId);
    }
    if (message.grantId !== "") {
      writer.uint32(18).string(message.grantId);
    }
    if (message.userId !== "") {
      writer.uint32(26).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataSuccessProjectGrantMember {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataSuccessProjectGrantMember();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.projectId = reader.string();
          break;
        case 2:
          message.grantId = reader.string();
          break;
        case 3:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataSuccessProjectGrantMember {
    return {
      projectId: isSet(object.projectId) ? String(object.projectId) : "",
      grantId: isSet(object.grantId) ? String(object.grantId) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: ImportDataSuccessProjectGrantMember): unknown {
    const obj: any = {};
    message.projectId !== undefined && (obj.projectId = message.projectId);
    message.grantId !== undefined && (obj.grantId = message.grantId);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<ImportDataSuccessProjectGrantMember>): ImportDataSuccessProjectGrantMember {
    return ImportDataSuccessProjectGrantMember.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataSuccessProjectGrantMember>): ImportDataSuccessProjectGrantMember {
    const message = createBaseImportDataSuccessProjectGrantMember();
    message.projectId = object.projectId ?? "";
    message.grantId = object.grantId ?? "";
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseImportDataSuccessUserLinks(): ImportDataSuccessUserLinks {
  return { userId: "", externalUserId: "", displayName: "", idpId: "" };
}

export const ImportDataSuccessUserLinks = {
  encode(message: ImportDataSuccessUserLinks, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.externalUserId !== "") {
      writer.uint32(18).string(message.externalUserId);
    }
    if (message.displayName !== "") {
      writer.uint32(26).string(message.displayName);
    }
    if (message.idpId !== "") {
      writer.uint32(34).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataSuccessUserLinks {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataSuccessUserLinks();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userId = reader.string();
          break;
        case 2:
          message.externalUserId = reader.string();
          break;
        case 3:
          message.displayName = reader.string();
          break;
        case 4:
          message.idpId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataSuccessUserLinks {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      externalUserId: isSet(object.externalUserId) ? String(object.externalUserId) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
    };
  },

  toJSON(message: ImportDataSuccessUserLinks): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.externalUserId !== undefined && (obj.externalUserId = message.externalUserId);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<ImportDataSuccessUserLinks>): ImportDataSuccessUserLinks {
    return ImportDataSuccessUserLinks.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataSuccessUserLinks>): ImportDataSuccessUserLinks {
    const message = createBaseImportDataSuccessUserLinks();
    message.userId = object.userId ?? "";
    message.externalUserId = object.externalUserId ?? "";
    message.displayName = object.displayName ?? "";
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseImportDataSuccessUserMetadata(): ImportDataSuccessUserMetadata {
  return { userId: "", key: "" };
}

export const ImportDataSuccessUserMetadata = {
  encode(message: ImportDataSuccessUserMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.key !== "") {
      writer.uint32(18).string(message.key);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataSuccessUserMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataSuccessUserMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userId = reader.string();
          break;
        case 2:
          message.key = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ImportDataSuccessUserMetadata {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      key: isSet(object.key) ? String(object.key) : "",
    };
  },

  toJSON(message: ImportDataSuccessUserMetadata): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.key !== undefined && (obj.key = message.key);
    return obj;
  },

  create(base?: DeepPartial<ImportDataSuccessUserMetadata>): ImportDataSuccessUserMetadata {
    return ImportDataSuccessUserMetadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataSuccessUserMetadata>): ImportDataSuccessUserMetadata {
    const message = createBaseImportDataSuccessUserMetadata();
    message.userId = object.userId ?? "";
    message.key = object.key ?? "";
    return message;
  },
};

function createBaseExportDataRequest(): ExportDataRequest {
  return {
    orgIds: [],
    excludedOrgIds: [],
    withPasswords: false,
    withOtp: false,
    responseOutput: false,
    localOutput: undefined,
    s3Output: undefined,
    gcsOutput: undefined,
    timeout: "",
  };
}

export const ExportDataRequest = {
  encode(message: ExportDataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.orgIds) {
      writer.uint32(10).string(v!);
    }
    for (const v of message.excludedOrgIds) {
      writer.uint32(18).string(v!);
    }
    if (message.withPasswords === true) {
      writer.uint32(24).bool(message.withPasswords);
    }
    if (message.withOtp === true) {
      writer.uint32(32).bool(message.withOtp);
    }
    if (message.responseOutput === true) {
      writer.uint32(40).bool(message.responseOutput);
    }
    if (message.localOutput !== undefined) {
      ExportDataRequest_LocalOutput.encode(message.localOutput, writer.uint32(50).fork()).ldelim();
    }
    if (message.s3Output !== undefined) {
      ExportDataRequest_S3Output.encode(message.s3Output, writer.uint32(58).fork()).ldelim();
    }
    if (message.gcsOutput !== undefined) {
      ExportDataRequest_GCSOutput.encode(message.gcsOutput, writer.uint32(66).fork()).ldelim();
    }
    if (message.timeout !== "") {
      writer.uint32(74).string(message.timeout);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportDataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportDataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgIds.push(reader.string());
          break;
        case 2:
          message.excludedOrgIds.push(reader.string());
          break;
        case 3:
          message.withPasswords = reader.bool();
          break;
        case 4:
          message.withOtp = reader.bool();
          break;
        case 5:
          message.responseOutput = reader.bool();
          break;
        case 6:
          message.localOutput = ExportDataRequest_LocalOutput.decode(reader, reader.uint32());
          break;
        case 7:
          message.s3Output = ExportDataRequest_S3Output.decode(reader, reader.uint32());
          break;
        case 8:
          message.gcsOutput = ExportDataRequest_GCSOutput.decode(reader, reader.uint32());
          break;
        case 9:
          message.timeout = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExportDataRequest {
    return {
      orgIds: Array.isArray(object?.orgIds) ? object.orgIds.map((e: any) => String(e)) : [],
      excludedOrgIds: Array.isArray(object?.excludedOrgIds) ? object.excludedOrgIds.map((e: any) => String(e)) : [],
      withPasswords: isSet(object.withPasswords) ? Boolean(object.withPasswords) : false,
      withOtp: isSet(object.withOtp) ? Boolean(object.withOtp) : false,
      responseOutput: isSet(object.responseOutput) ? Boolean(object.responseOutput) : false,
      localOutput: isSet(object.localOutput) ? ExportDataRequest_LocalOutput.fromJSON(object.localOutput) : undefined,
      s3Output: isSet(object.s3Output) ? ExportDataRequest_S3Output.fromJSON(object.s3Output) : undefined,
      gcsOutput: isSet(object.gcsOutput) ? ExportDataRequest_GCSOutput.fromJSON(object.gcsOutput) : undefined,
      timeout: isSet(object.timeout) ? String(object.timeout) : "",
    };
  },

  toJSON(message: ExportDataRequest): unknown {
    const obj: any = {};
    if (message.orgIds) {
      obj.orgIds = message.orgIds.map((e) => e);
    } else {
      obj.orgIds = [];
    }
    if (message.excludedOrgIds) {
      obj.excludedOrgIds = message.excludedOrgIds.map((e) => e);
    } else {
      obj.excludedOrgIds = [];
    }
    message.withPasswords !== undefined && (obj.withPasswords = message.withPasswords);
    message.withOtp !== undefined && (obj.withOtp = message.withOtp);
    message.responseOutput !== undefined && (obj.responseOutput = message.responseOutput);
    message.localOutput !== undefined &&
      (obj.localOutput = message.localOutput ? ExportDataRequest_LocalOutput.toJSON(message.localOutput) : undefined);
    message.s3Output !== undefined &&
      (obj.s3Output = message.s3Output ? ExportDataRequest_S3Output.toJSON(message.s3Output) : undefined);
    message.gcsOutput !== undefined &&
      (obj.gcsOutput = message.gcsOutput ? ExportDataRequest_GCSOutput.toJSON(message.gcsOutput) : undefined);
    message.timeout !== undefined && (obj.timeout = message.timeout);
    return obj;
  },

  create(base?: DeepPartial<ExportDataRequest>): ExportDataRequest {
    return ExportDataRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportDataRequest>): ExportDataRequest {
    const message = createBaseExportDataRequest();
    message.orgIds = object.orgIds?.map((e) => e) || [];
    message.excludedOrgIds = object.excludedOrgIds?.map((e) => e) || [];
    message.withPasswords = object.withPasswords ?? false;
    message.withOtp = object.withOtp ?? false;
    message.responseOutput = object.responseOutput ?? false;
    message.localOutput = (object.localOutput !== undefined && object.localOutput !== null)
      ? ExportDataRequest_LocalOutput.fromPartial(object.localOutput)
      : undefined;
    message.s3Output = (object.s3Output !== undefined && object.s3Output !== null)
      ? ExportDataRequest_S3Output.fromPartial(object.s3Output)
      : undefined;
    message.gcsOutput = (object.gcsOutput !== undefined && object.gcsOutput !== null)
      ? ExportDataRequest_GCSOutput.fromPartial(object.gcsOutput)
      : undefined;
    message.timeout = object.timeout ?? "";
    return message;
  },
};

function createBaseExportDataRequest_LocalOutput(): ExportDataRequest_LocalOutput {
  return { path: "" };
}

export const ExportDataRequest_LocalOutput = {
  encode(message: ExportDataRequest_LocalOutput, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportDataRequest_LocalOutput {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportDataRequest_LocalOutput();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExportDataRequest_LocalOutput {
    return { path: isSet(object.path) ? String(object.path) : "" };
  },

  toJSON(message: ExportDataRequest_LocalOutput): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    return obj;
  },

  create(base?: DeepPartial<ExportDataRequest_LocalOutput>): ExportDataRequest_LocalOutput {
    return ExportDataRequest_LocalOutput.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportDataRequest_LocalOutput>): ExportDataRequest_LocalOutput {
    const message = createBaseExportDataRequest_LocalOutput();
    message.path = object.path ?? "";
    return message;
  },
};

function createBaseExportDataRequest_S3Output(): ExportDataRequest_S3Output {
  return { path: "", endpoint: "", accessKeyId: "", secretAccessKey: "", ssl: false, bucket: "" };
}

export const ExportDataRequest_S3Output = {
  encode(message: ExportDataRequest_S3Output, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    if (message.endpoint !== "") {
      writer.uint32(18).string(message.endpoint);
    }
    if (message.accessKeyId !== "") {
      writer.uint32(26).string(message.accessKeyId);
    }
    if (message.secretAccessKey !== "") {
      writer.uint32(34).string(message.secretAccessKey);
    }
    if (message.ssl === true) {
      writer.uint32(40).bool(message.ssl);
    }
    if (message.bucket !== "") {
      writer.uint32(50).string(message.bucket);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportDataRequest_S3Output {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportDataRequest_S3Output();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        case 2:
          message.endpoint = reader.string();
          break;
        case 3:
          message.accessKeyId = reader.string();
          break;
        case 4:
          message.secretAccessKey = reader.string();
          break;
        case 5:
          message.ssl = reader.bool();
          break;
        case 6:
          message.bucket = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExportDataRequest_S3Output {
    return {
      path: isSet(object.path) ? String(object.path) : "",
      endpoint: isSet(object.endpoint) ? String(object.endpoint) : "",
      accessKeyId: isSet(object.accessKeyId) ? String(object.accessKeyId) : "",
      secretAccessKey: isSet(object.secretAccessKey) ? String(object.secretAccessKey) : "",
      ssl: isSet(object.ssl) ? Boolean(object.ssl) : false,
      bucket: isSet(object.bucket) ? String(object.bucket) : "",
    };
  },

  toJSON(message: ExportDataRequest_S3Output): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    message.endpoint !== undefined && (obj.endpoint = message.endpoint);
    message.accessKeyId !== undefined && (obj.accessKeyId = message.accessKeyId);
    message.secretAccessKey !== undefined && (obj.secretAccessKey = message.secretAccessKey);
    message.ssl !== undefined && (obj.ssl = message.ssl);
    message.bucket !== undefined && (obj.bucket = message.bucket);
    return obj;
  },

  create(base?: DeepPartial<ExportDataRequest_S3Output>): ExportDataRequest_S3Output {
    return ExportDataRequest_S3Output.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportDataRequest_S3Output>): ExportDataRequest_S3Output {
    const message = createBaseExportDataRequest_S3Output();
    message.path = object.path ?? "";
    message.endpoint = object.endpoint ?? "";
    message.accessKeyId = object.accessKeyId ?? "";
    message.secretAccessKey = object.secretAccessKey ?? "";
    message.ssl = object.ssl ?? false;
    message.bucket = object.bucket ?? "";
    return message;
  },
};

function createBaseExportDataRequest_GCSOutput(): ExportDataRequest_GCSOutput {
  return { bucket: "", serviceaccountJson: "", path: "" };
}

export const ExportDataRequest_GCSOutput = {
  encode(message: ExportDataRequest_GCSOutput, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bucket !== "") {
      writer.uint32(10).string(message.bucket);
    }
    if (message.serviceaccountJson !== "") {
      writer.uint32(18).string(message.serviceaccountJson);
    }
    if (message.path !== "") {
      writer.uint32(26).string(message.path);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportDataRequest_GCSOutput {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportDataRequest_GCSOutput();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.bucket = reader.string();
          break;
        case 2:
          message.serviceaccountJson = reader.string();
          break;
        case 3:
          message.path = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExportDataRequest_GCSOutput {
    return {
      bucket: isSet(object.bucket) ? String(object.bucket) : "",
      serviceaccountJson: isSet(object.serviceaccountJson) ? String(object.serviceaccountJson) : "",
      path: isSet(object.path) ? String(object.path) : "",
    };
  },

  toJSON(message: ExportDataRequest_GCSOutput): unknown {
    const obj: any = {};
    message.bucket !== undefined && (obj.bucket = message.bucket);
    message.serviceaccountJson !== undefined && (obj.serviceaccountJson = message.serviceaccountJson);
    message.path !== undefined && (obj.path = message.path);
    return obj;
  },

  create(base?: DeepPartial<ExportDataRequest_GCSOutput>): ExportDataRequest_GCSOutput {
    return ExportDataRequest_GCSOutput.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportDataRequest_GCSOutput>): ExportDataRequest_GCSOutput {
    const message = createBaseExportDataRequest_GCSOutput();
    message.bucket = object.bucket ?? "";
    message.serviceaccountJson = object.serviceaccountJson ?? "";
    message.path = object.path ?? "";
    return message;
  },
};

function createBaseExportDataResponse(): ExportDataResponse {
  return { orgs: [] };
}

export const ExportDataResponse = {
  encode(message: ExportDataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.orgs) {
      DataOrg.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportDataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportDataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgs.push(DataOrg.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExportDataResponse {
    return { orgs: Array.isArray(object?.orgs) ? object.orgs.map((e: any) => DataOrg.fromJSON(e)) : [] };
  },

  toJSON(message: ExportDataResponse): unknown {
    const obj: any = {};
    if (message.orgs) {
      obj.orgs = message.orgs.map((e) => e ? DataOrg.toJSON(e) : undefined);
    } else {
      obj.orgs = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ExportDataResponse>): ExportDataResponse {
    return ExportDataResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportDataResponse>): ExportDataResponse {
    const message = createBaseExportDataResponse();
    message.orgs = object.orgs?.map((e) => DataOrg.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListEventsRequest(): ListEventsRequest {
  return {
    sequence: 0,
    limit: 0,
    asc: false,
    editorUserId: "",
    eventTypes: [],
    aggregateId: "",
    aggregateTypes: [],
    resourceOwner: "",
    creationDate: undefined,
  };
}

export const ListEventsRequest = {
  encode(message: ListEventsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sequence !== 0) {
      writer.uint32(8).uint64(message.sequence);
    }
    if (message.limit !== 0) {
      writer.uint32(16).uint32(message.limit);
    }
    if (message.asc === true) {
      writer.uint32(24).bool(message.asc);
    }
    if (message.editorUserId !== "") {
      writer.uint32(34).string(message.editorUserId);
    }
    for (const v of message.eventTypes) {
      writer.uint32(42).string(v!);
    }
    if (message.aggregateId !== "") {
      writer.uint32(50).string(message.aggregateId);
    }
    for (const v of message.aggregateTypes) {
      writer.uint32(58).string(v!);
    }
    if (message.resourceOwner !== "") {
      writer.uint32(66).string(message.resourceOwner);
    }
    if (message.creationDate !== undefined) {
      Timestamp.encode(toTimestamp(message.creationDate), writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListEventsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListEventsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequence = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.limit = reader.uint32();
          break;
        case 3:
          message.asc = reader.bool();
          break;
        case 4:
          message.editorUserId = reader.string();
          break;
        case 5:
          message.eventTypes.push(reader.string());
          break;
        case 6:
          message.aggregateId = reader.string();
          break;
        case 7:
          message.aggregateTypes.push(reader.string());
          break;
        case 8:
          message.resourceOwner = reader.string();
          break;
        case 9:
          message.creationDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListEventsRequest {
    return {
      sequence: isSet(object.sequence) ? Number(object.sequence) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
      asc: isSet(object.asc) ? Boolean(object.asc) : false,
      editorUserId: isSet(object.editorUserId) ? String(object.editorUserId) : "",
      eventTypes: Array.isArray(object?.eventTypes) ? object.eventTypes.map((e: any) => String(e)) : [],
      aggregateId: isSet(object.aggregateId) ? String(object.aggregateId) : "",
      aggregateTypes: Array.isArray(object?.aggregateTypes) ? object.aggregateTypes.map((e: any) => String(e)) : [],
      resourceOwner: isSet(object.resourceOwner) ? String(object.resourceOwner) : "",
      creationDate: isSet(object.creationDate) ? fromJsonTimestamp(object.creationDate) : undefined,
    };
  },

  toJSON(message: ListEventsRequest): unknown {
    const obj: any = {};
    message.sequence !== undefined && (obj.sequence = Math.round(message.sequence));
    message.limit !== undefined && (obj.limit = Math.round(message.limit));
    message.asc !== undefined && (obj.asc = message.asc);
    message.editorUserId !== undefined && (obj.editorUserId = message.editorUserId);
    if (message.eventTypes) {
      obj.eventTypes = message.eventTypes.map((e) => e);
    } else {
      obj.eventTypes = [];
    }
    message.aggregateId !== undefined && (obj.aggregateId = message.aggregateId);
    if (message.aggregateTypes) {
      obj.aggregateTypes = message.aggregateTypes.map((e) => e);
    } else {
      obj.aggregateTypes = [];
    }
    message.resourceOwner !== undefined && (obj.resourceOwner = message.resourceOwner);
    message.creationDate !== undefined && (obj.creationDate = message.creationDate.toISOString());
    return obj;
  },

  create(base?: DeepPartial<ListEventsRequest>): ListEventsRequest {
    return ListEventsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListEventsRequest>): ListEventsRequest {
    const message = createBaseListEventsRequest();
    message.sequence = object.sequence ?? 0;
    message.limit = object.limit ?? 0;
    message.asc = object.asc ?? false;
    message.editorUserId = object.editorUserId ?? "";
    message.eventTypes = object.eventTypes?.map((e) => e) || [];
    message.aggregateId = object.aggregateId ?? "";
    message.aggregateTypes = object.aggregateTypes?.map((e) => e) || [];
    message.resourceOwner = object.resourceOwner ?? "";
    message.creationDate = object.creationDate ?? undefined;
    return message;
  },
};

function createBaseListEventsResponse(): ListEventsResponse {
  return { events: [] };
}

export const ListEventsResponse = {
  encode(message: ListEventsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.events) {
      Event.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListEventsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListEventsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.events.push(Event.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListEventsResponse {
    return { events: Array.isArray(object?.events) ? object.events.map((e: any) => Event.fromJSON(e)) : [] };
  },

  toJSON(message: ListEventsResponse): unknown {
    const obj: any = {};
    if (message.events) {
      obj.events = message.events.map((e) => e ? Event.toJSON(e) : undefined);
    } else {
      obj.events = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListEventsResponse>): ListEventsResponse {
    return ListEventsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListEventsResponse>): ListEventsResponse {
    const message = createBaseListEventsResponse();
    message.events = object.events?.map((e) => Event.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListEventTypesRequest(): ListEventTypesRequest {
  return {};
}

export const ListEventTypesRequest = {
  encode(_: ListEventTypesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListEventTypesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListEventTypesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListEventTypesRequest {
    return {};
  },

  toJSON(_: ListEventTypesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListEventTypesRequest>): ListEventTypesRequest {
    return ListEventTypesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListEventTypesRequest>): ListEventTypesRequest {
    const message = createBaseListEventTypesRequest();
    return message;
  },
};

function createBaseListEventTypesResponse(): ListEventTypesResponse {
  return { eventTypes: [] };
}

export const ListEventTypesResponse = {
  encode(message: ListEventTypesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.eventTypes) {
      EventType.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListEventTypesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListEventTypesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.eventTypes.push(EventType.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListEventTypesResponse {
    return {
      eventTypes: Array.isArray(object?.eventTypes) ? object.eventTypes.map((e: any) => EventType.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListEventTypesResponse): unknown {
    const obj: any = {};
    if (message.eventTypes) {
      obj.eventTypes = message.eventTypes.map((e) => e ? EventType.toJSON(e) : undefined);
    } else {
      obj.eventTypes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListEventTypesResponse>): ListEventTypesResponse {
    return ListEventTypesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListEventTypesResponse>): ListEventTypesResponse {
    const message = createBaseListEventTypesResponse();
    message.eventTypes = object.eventTypes?.map((e) => EventType.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListAggregateTypesRequest(): ListAggregateTypesRequest {
  return {};
}

export const ListAggregateTypesRequest = {
  encode(_: ListAggregateTypesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListAggregateTypesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListAggregateTypesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListAggregateTypesRequest {
    return {};
  },

  toJSON(_: ListAggregateTypesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListAggregateTypesRequest>): ListAggregateTypesRequest {
    return ListAggregateTypesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListAggregateTypesRequest>): ListAggregateTypesRequest {
    const message = createBaseListAggregateTypesRequest();
    return message;
  },
};

function createBaseListAggregateTypesResponse(): ListAggregateTypesResponse {
  return { aggregateTypes: [] };
}

export const ListAggregateTypesResponse = {
  encode(message: ListAggregateTypesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.aggregateTypes) {
      AggregateType.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListAggregateTypesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListAggregateTypesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.aggregateTypes.push(AggregateType.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListAggregateTypesResponse {
    return {
      aggregateTypes: Array.isArray(object?.aggregateTypes)
        ? object.aggregateTypes.map((e: any) => AggregateType.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ListAggregateTypesResponse): unknown {
    const obj: any = {};
    if (message.aggregateTypes) {
      obj.aggregateTypes = message.aggregateTypes.map((e) => e ? AggregateType.toJSON(e) : undefined);
    } else {
      obj.aggregateTypes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListAggregateTypesResponse>): ListAggregateTypesResponse {
    return ListAggregateTypesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListAggregateTypesResponse>): ListAggregateTypesResponse {
    const message = createBaseListAggregateTypesResponse();
    message.aggregateTypes = object.aggregateTypes?.map((e) => AggregateType.fromPartial(e)) || [];
    return message;
  },
};

export type AdminServiceDefinition = typeof AdminServiceDefinition;
export const AdminServiceDefinition = {
  name: "AdminService",
  fullName: "zitadel.admin.v1.AdminService",
  methods: {
    healthz: {
      name: "Healthz",
      requestType: HealthzRequest,
      requestStream: false,
      responseType: HealthzResponse,
      responseStream: false,
      options: {},
    },
    getSupportedLanguages: {
      name: "GetSupportedLanguages",
      requestType: GetSupportedLanguagesRequest,
      requestStream: false,
      responseType: GetSupportedLanguagesResponse,
      responseStream: false,
      options: {},
    },
    setDefaultLanguage: {
      name: "SetDefaultLanguage",
      requestType: SetDefaultLanguageRequest,
      requestStream: false,
      responseType: SetDefaultLanguageResponse,
      responseStream: false,
      options: {},
    },
    getDefaultLanguage: {
      name: "GetDefaultLanguage",
      requestType: GetDefaultLanguageRequest,
      requestStream: false,
      responseType: GetDefaultLanguageResponse,
      responseStream: false,
      options: {},
    },
    getMyInstance: {
      name: "GetMyInstance",
      requestType: GetMyInstanceRequest,
      requestStream: false,
      responseType: GetMyInstanceResponse,
      responseStream: false,
      options: {},
    },
    listInstanceDomains: {
      name: "ListInstanceDomains",
      requestType: ListInstanceDomainsRequest,
      requestStream: false,
      responseType: ListInstanceDomainsResponse,
      responseStream: false,
      options: {},
    },
    listSecretGenerators: {
      name: "ListSecretGenerators",
      requestType: ListSecretGeneratorsRequest,
      requestStream: false,
      responseType: ListSecretGeneratorsResponse,
      responseStream: false,
      options: {},
    },
    getSecretGenerator: {
      name: "GetSecretGenerator",
      requestType: GetSecretGeneratorRequest,
      requestStream: false,
      responseType: GetSecretGeneratorResponse,
      responseStream: false,
      options: {},
    },
    updateSecretGenerator: {
      name: "UpdateSecretGenerator",
      requestType: UpdateSecretGeneratorRequest,
      requestStream: false,
      responseType: UpdateSecretGeneratorResponse,
      responseStream: false,
      options: {},
    },
    getSMTPConfig: {
      name: "GetSMTPConfig",
      requestType: GetSMTPConfigRequest,
      requestStream: false,
      responseType: GetSMTPConfigResponse,
      responseStream: false,
      options: {},
    },
    addSMTPConfig: {
      name: "AddSMTPConfig",
      requestType: AddSMTPConfigRequest,
      requestStream: false,
      responseType: AddSMTPConfigResponse,
      responseStream: false,
      options: {},
    },
    updateSMTPConfig: {
      name: "UpdateSMTPConfig",
      requestType: UpdateSMTPConfigRequest,
      requestStream: false,
      responseType: UpdateSMTPConfigResponse,
      responseStream: false,
      options: {},
    },
    updateSMTPConfigPassword: {
      name: "UpdateSMTPConfigPassword",
      requestType: UpdateSMTPConfigPasswordRequest,
      requestStream: false,
      responseType: UpdateSMTPConfigPasswordResponse,
      responseStream: false,
      options: {},
    },
    removeSMTPConfig: {
      name: "RemoveSMTPConfig",
      requestType: RemoveSMTPConfigRequest,
      requestStream: false,
      responseType: RemoveSMTPConfigResponse,
      responseStream: false,
      options: {},
    },
    listSMSProviders: {
      name: "ListSMSProviders",
      requestType: ListSMSProvidersRequest,
      requestStream: false,
      responseType: ListSMSProvidersResponse,
      responseStream: false,
      options: {},
    },
    getSMSProvider: {
      name: "GetSMSProvider",
      requestType: GetSMSProviderRequest,
      requestStream: false,
      responseType: GetSMSProviderResponse,
      responseStream: false,
      options: {},
    },
    addSMSProviderTwilio: {
      name: "AddSMSProviderTwilio",
      requestType: AddSMSProviderTwilioRequest,
      requestStream: false,
      responseType: AddSMSProviderTwilioResponse,
      responseStream: false,
      options: {},
    },
    updateSMSProviderTwilio: {
      name: "UpdateSMSProviderTwilio",
      requestType: UpdateSMSProviderTwilioRequest,
      requestStream: false,
      responseType: UpdateSMSProviderTwilioResponse,
      responseStream: false,
      options: {},
    },
    updateSMSProviderTwilioToken: {
      name: "UpdateSMSProviderTwilioToken",
      requestType: UpdateSMSProviderTwilioTokenRequest,
      requestStream: false,
      responseType: UpdateSMSProviderTwilioTokenResponse,
      responseStream: false,
      options: {},
    },
    activateSMSProvider: {
      name: "ActivateSMSProvider",
      requestType: ActivateSMSProviderRequest,
      requestStream: false,
      responseType: ActivateSMSProviderResponse,
      responseStream: false,
      options: {},
    },
    deactivateSMSProvider: {
      name: "DeactivateSMSProvider",
      requestType: DeactivateSMSProviderRequest,
      requestStream: false,
      responseType: DeactivateSMSProviderResponse,
      responseStream: false,
      options: {},
    },
    removeSMSProvider: {
      name: "RemoveSMSProvider",
      requestType: RemoveSMSProviderRequest,
      requestStream: false,
      responseType: RemoveSMSProviderResponse,
      responseStream: false,
      options: {},
    },
    getOIDCSettings: {
      name: "GetOIDCSettings",
      requestType: GetOIDCSettingsRequest,
      requestStream: false,
      responseType: GetOIDCSettingsResponse,
      responseStream: false,
      options: {},
    },
    addOIDCSettings: {
      name: "AddOIDCSettings",
      requestType: AddOIDCSettingsRequest,
      requestStream: false,
      responseType: AddOIDCSettingsResponse,
      responseStream: false,
      options: {},
    },
    updateOIDCSettings: {
      name: "UpdateOIDCSettings",
      requestType: UpdateOIDCSettingsRequest,
      requestStream: false,
      responseType: UpdateOIDCSettingsResponse,
      responseStream: false,
      options: {},
    },
    getFileSystemNotificationProvider: {
      name: "GetFileSystemNotificationProvider",
      requestType: GetFileSystemNotificationProviderRequest,
      requestStream: false,
      responseType: GetFileSystemNotificationProviderResponse,
      responseStream: false,
      options: {},
    },
    getLogNotificationProvider: {
      name: "GetLogNotificationProvider",
      requestType: GetLogNotificationProviderRequest,
      requestStream: false,
      responseType: GetLogNotificationProviderResponse,
      responseStream: false,
      options: {},
    },
    getSecurityPolicy: {
      name: "GetSecurityPolicy",
      requestType: GetSecurityPolicyRequest,
      requestStream: false,
      responseType: GetSecurityPolicyResponse,
      responseStream: false,
      options: {},
    },
    setSecurityPolicy: {
      name: "SetSecurityPolicy",
      requestType: SetSecurityPolicyRequest,
      requestStream: false,
      responseType: SetSecurityPolicyResponse,
      responseStream: false,
      options: {},
    },
    getOrgByID: {
      name: "GetOrgByID",
      requestType: GetOrgByIDRequest,
      requestStream: false,
      responseType: GetOrgByIDResponse,
      responseStream: false,
      options: {},
    },
    isOrgUnique: {
      name: "IsOrgUnique",
      requestType: IsOrgUniqueRequest,
      requestStream: false,
      responseType: IsOrgUniqueResponse,
      responseStream: false,
      options: {},
    },
    setDefaultOrg: {
      name: "SetDefaultOrg",
      requestType: SetDefaultOrgRequest,
      requestStream: false,
      responseType: SetDefaultOrgResponse,
      responseStream: false,
      options: {},
    },
    getDefaultOrg: {
      name: "GetDefaultOrg",
      requestType: GetDefaultOrgRequest,
      requestStream: false,
      responseType: GetDefaultOrgResponse,
      responseStream: false,
      options: {},
    },
    listOrgs: {
      name: "ListOrgs",
      requestType: ListOrgsRequest,
      requestStream: false,
      responseType: ListOrgsResponse,
      responseStream: false,
      options: {},
    },
    setUpOrg: {
      name: "SetUpOrg",
      requestType: SetUpOrgRequest,
      requestStream: false,
      responseType: SetUpOrgResponse,
      responseStream: false,
      options: {},
    },
    removeOrg: {
      name: "RemoveOrg",
      requestType: RemoveOrgRequest,
      requestStream: false,
      responseType: RemoveOrgResponse,
      responseStream: false,
      options: {},
    },
    getIDPByID: {
      name: "GetIDPByID",
      requestType: GetIDPByIDRequest,
      requestStream: false,
      responseType: GetIDPByIDResponse,
      responseStream: false,
      options: {},
    },
    listIDPs: {
      name: "ListIDPs",
      requestType: ListIDPsRequest,
      requestStream: false,
      responseType: ListIDPsResponse,
      responseStream: false,
      options: {},
    },
    addOIDCIDP: {
      name: "AddOIDCIDP",
      requestType: AddOIDCIDPRequest,
      requestStream: false,
      responseType: AddOIDCIDPResponse,
      responseStream: false,
      options: {},
    },
    addJWTIDP: {
      name: "AddJWTIDP",
      requestType: AddJWTIDPRequest,
      requestStream: false,
      responseType: AddJWTIDPResponse,
      responseStream: false,
      options: {},
    },
    updateIDP: {
      name: "UpdateIDP",
      requestType: UpdateIDPRequest,
      requestStream: false,
      responseType: UpdateIDPResponse,
      responseStream: false,
      options: {},
    },
    deactivateIDP: {
      name: "DeactivateIDP",
      requestType: DeactivateIDPRequest,
      requestStream: false,
      responseType: DeactivateIDPResponse,
      responseStream: false,
      options: {},
    },
    reactivateIDP: {
      name: "ReactivateIDP",
      requestType: ReactivateIDPRequest,
      requestStream: false,
      responseType: ReactivateIDPResponse,
      responseStream: false,
      options: {},
    },
    removeIDP: {
      name: "RemoveIDP",
      requestType: RemoveIDPRequest,
      requestStream: false,
      responseType: RemoveIDPResponse,
      responseStream: false,
      options: {},
    },
    updateIDPOIDCConfig: {
      name: "UpdateIDPOIDCConfig",
      requestType: UpdateIDPOIDCConfigRequest,
      requestStream: false,
      responseType: UpdateIDPOIDCConfigResponse,
      responseStream: false,
      options: {},
    },
    updateIDPJWTConfig: {
      name: "UpdateIDPJWTConfig",
      requestType: UpdateIDPJWTConfigRequest,
      requestStream: false,
      responseType: UpdateIDPJWTConfigResponse,
      responseStream: false,
      options: {},
    },
    /**
     * Returns all identity providers, which match the query
     * Limit should always be set, there is a default limit set by the service
     */
    listProviders: {
      name: "ListProviders",
      requestType: ListProvidersRequest,
      requestStream: false,
      responseType: ListProvidersResponse,
      responseStream: false,
      options: {},
    },
    /** Returns an identity provider of the instance */
    getProviderByID: {
      name: "GetProviderByID",
      requestType: GetProviderByIDRequest,
      requestStream: false,
      responseType: GetProviderByIDResponse,
      responseStream: false,
      options: {},
    },
    /** Add a new OAuth2 identity provider on the instance */
    addGenericOAuthProvider: {
      name: "AddGenericOAuthProvider",
      requestType: AddGenericOAuthProviderRequest,
      requestStream: false,
      responseType: AddGenericOAuthProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Change an existing OAuth2 identity provider on the instance */
    updateGenericOAuthProvider: {
      name: "UpdateGenericOAuthProvider",
      requestType: UpdateGenericOAuthProviderRequest,
      requestStream: false,
      responseType: UpdateGenericOAuthProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Add a new OIDC identity provider on the instance */
    addGenericOIDCProvider: {
      name: "AddGenericOIDCProvider",
      requestType: AddGenericOIDCProviderRequest,
      requestStream: false,
      responseType: AddGenericOIDCProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Change an existing OIDC identity provider on the instance */
    updateGenericOIDCProvider: {
      name: "UpdateGenericOIDCProvider",
      requestType: UpdateGenericOIDCProviderRequest,
      requestStream: false,
      responseType: UpdateGenericOIDCProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Add a new JWT identity provider on the instance */
    addJWTProvider: {
      name: "AddJWTProvider",
      requestType: AddJWTProviderRequest,
      requestStream: false,
      responseType: AddJWTProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Change an existing JWT identity provider on the instance */
    updateJWTProvider: {
      name: "UpdateJWTProvider",
      requestType: UpdateJWTProviderRequest,
      requestStream: false,
      responseType: UpdateJWTProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Add a new GitHub identity provider on the instance */
    addGitHubProvider: {
      name: "AddGitHubProvider",
      requestType: AddGitHubProviderRequest,
      requestStream: false,
      responseType: AddGitHubProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Change an existing GitHub identity provider on the instance */
    updateGitHubProvider: {
      name: "UpdateGitHubProvider",
      requestType: UpdateGitHubProviderRequest,
      requestStream: false,
      responseType: UpdateGitHubProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Add a new GitHub Enterprise Server identity provider on the instance */
    addGitHubEnterpriseServerProvider: {
      name: "AddGitHubEnterpriseServerProvider",
      requestType: AddGitHubEnterpriseServerProviderRequest,
      requestStream: false,
      responseType: AddGitHubEnterpriseServerProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Change an existing GitHub Enterprise Server identity provider on the instance */
    updateGitHubEnterpriseServerProvider: {
      name: "UpdateGitHubEnterpriseServerProvider",
      requestType: UpdateGitHubEnterpriseServerProviderRequest,
      requestStream: false,
      responseType: UpdateGitHubEnterpriseServerProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Add a new Google identity provider on the instance */
    addGoogleProvider: {
      name: "AddGoogleProvider",
      requestType: AddGoogleProviderRequest,
      requestStream: false,
      responseType: AddGoogleProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Change an existing Google identity provider on the instance */
    updateGoogleProvider: {
      name: "UpdateGoogleProvider",
      requestType: UpdateGoogleProviderRequest,
      requestStream: false,
      responseType: UpdateGoogleProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Add a new LDAP identity provider on the instance */
    addLDAPProvider: {
      name: "AddLDAPProvider",
      requestType: AddLDAPProviderRequest,
      requestStream: false,
      responseType: AddLDAPProviderResponse,
      responseStream: false,
      options: {},
    },
    /** Change an existing LDAP identity provider on the instance */
    updateLDAPProvider: {
      name: "UpdateLDAPProvider",
      requestType: UpdateLDAPProviderRequest,
      requestStream: false,
      responseType: UpdateLDAPProviderResponse,
      responseStream: false,
      options: {},
    },
    /**
     * Remove an identity provider
     * Will remove all linked providers of this configuration on the users
     */
    deleteProvider: {
      name: "DeleteProvider",
      requestType: DeleteProviderRequest,
      requestStream: false,
      responseType: DeleteProviderResponse,
      responseStream: false,
      options: {},
    },
    getOrgIAMPolicy: {
      name: "GetOrgIAMPolicy",
      requestType: GetOrgIAMPolicyRequest,
      requestStream: false,
      responseType: GetOrgIAMPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateOrgIAMPolicy: {
      name: "UpdateOrgIAMPolicy",
      requestType: UpdateOrgIAMPolicyRequest,
      requestStream: false,
      responseType: UpdateOrgIAMPolicyResponse,
      responseStream: false,
      options: {},
    },
    getCustomOrgIAMPolicy: {
      name: "GetCustomOrgIAMPolicy",
      requestType: GetCustomOrgIAMPolicyRequest,
      requestStream: false,
      responseType: GetCustomOrgIAMPolicyResponse,
      responseStream: false,
      options: {},
    },
    addCustomOrgIAMPolicy: {
      name: "AddCustomOrgIAMPolicy",
      requestType: AddCustomOrgIAMPolicyRequest,
      requestStream: false,
      responseType: AddCustomOrgIAMPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateCustomOrgIAMPolicy: {
      name: "UpdateCustomOrgIAMPolicy",
      requestType: UpdateCustomOrgIAMPolicyRequest,
      requestStream: false,
      responseType: UpdateCustomOrgIAMPolicyResponse,
      responseStream: false,
      options: {},
    },
    resetCustomOrgIAMPolicyToDefault: {
      name: "ResetCustomOrgIAMPolicyToDefault",
      requestType: ResetCustomOrgIAMPolicyToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomOrgIAMPolicyToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getDomainPolicy: {
      name: "GetDomainPolicy",
      requestType: GetDomainPolicyRequest,
      requestStream: false,
      responseType: GetDomainPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateDomainPolicy: {
      name: "UpdateDomainPolicy",
      requestType: UpdateDomainPolicyRequest,
      requestStream: false,
      responseType: UpdateDomainPolicyResponse,
      responseStream: false,
      options: {},
    },
    getCustomDomainPolicy: {
      name: "GetCustomDomainPolicy",
      requestType: GetCustomDomainPolicyRequest,
      requestStream: false,
      responseType: GetCustomDomainPolicyResponse,
      responseStream: false,
      options: {},
    },
    addCustomDomainPolicy: {
      name: "AddCustomDomainPolicy",
      requestType: AddCustomDomainPolicyRequest,
      requestStream: false,
      responseType: AddCustomDomainPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateCustomDomainPolicy: {
      name: "UpdateCustomDomainPolicy",
      requestType: UpdateCustomDomainPolicyRequest,
      requestStream: false,
      responseType: UpdateCustomDomainPolicyResponse,
      responseStream: false,
      options: {},
    },
    resetCustomDomainPolicyToDefault: {
      name: "ResetCustomDomainPolicyToDefault",
      requestType: ResetCustomDomainPolicyToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomDomainPolicyToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getLabelPolicy: {
      name: "GetLabelPolicy",
      requestType: GetLabelPolicyRequest,
      requestStream: false,
      responseType: GetLabelPolicyResponse,
      responseStream: false,
      options: {},
    },
    getPreviewLabelPolicy: {
      name: "GetPreviewLabelPolicy",
      requestType: GetPreviewLabelPolicyRequest,
      requestStream: false,
      responseType: GetPreviewLabelPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateLabelPolicy: {
      name: "UpdateLabelPolicy",
      requestType: UpdateLabelPolicyRequest,
      requestStream: false,
      responseType: UpdateLabelPolicyResponse,
      responseStream: false,
      options: {},
    },
    activateLabelPolicy: {
      name: "ActivateLabelPolicy",
      requestType: ActivateLabelPolicyRequest,
      requestStream: false,
      responseType: ActivateLabelPolicyResponse,
      responseStream: false,
      options: {},
    },
    removeLabelPolicyLogo: {
      name: "RemoveLabelPolicyLogo",
      requestType: RemoveLabelPolicyLogoRequest,
      requestStream: false,
      responseType: RemoveLabelPolicyLogoResponse,
      responseStream: false,
      options: {},
    },
    removeLabelPolicyLogoDark: {
      name: "RemoveLabelPolicyLogoDark",
      requestType: RemoveLabelPolicyLogoDarkRequest,
      requestStream: false,
      responseType: RemoveLabelPolicyLogoDarkResponse,
      responseStream: false,
      options: {},
    },
    removeLabelPolicyIcon: {
      name: "RemoveLabelPolicyIcon",
      requestType: RemoveLabelPolicyIconRequest,
      requestStream: false,
      responseType: RemoveLabelPolicyIconResponse,
      responseStream: false,
      options: {},
    },
    removeLabelPolicyIconDark: {
      name: "RemoveLabelPolicyIconDark",
      requestType: RemoveLabelPolicyIconDarkRequest,
      requestStream: false,
      responseType: RemoveLabelPolicyIconDarkResponse,
      responseStream: false,
      options: {},
    },
    removeLabelPolicyFont: {
      name: "RemoveLabelPolicyFont",
      requestType: RemoveLabelPolicyFontRequest,
      requestStream: false,
      responseType: RemoveLabelPolicyFontResponse,
      responseStream: false,
      options: {},
    },
    getLoginPolicy: {
      name: "GetLoginPolicy",
      requestType: GetLoginPolicyRequest,
      requestStream: false,
      responseType: GetLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateLoginPolicy: {
      name: "UpdateLoginPolicy",
      requestType: UpdateLoginPolicyRequest,
      requestStream: false,
      responseType: UpdateLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
    listLoginPolicyIDPs: {
      name: "ListLoginPolicyIDPs",
      requestType: ListLoginPolicyIDPsRequest,
      requestStream: false,
      responseType: ListLoginPolicyIDPsResponse,
      responseStream: false,
      options: {},
    },
    addIDPToLoginPolicy: {
      name: "AddIDPToLoginPolicy",
      requestType: AddIDPToLoginPolicyRequest,
      requestStream: false,
      responseType: AddIDPToLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
    removeIDPFromLoginPolicy: {
      name: "RemoveIDPFromLoginPolicy",
      requestType: RemoveIDPFromLoginPolicyRequest,
      requestStream: false,
      responseType: RemoveIDPFromLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
    listLoginPolicySecondFactors: {
      name: "ListLoginPolicySecondFactors",
      requestType: ListLoginPolicySecondFactorsRequest,
      requestStream: false,
      responseType: ListLoginPolicySecondFactorsResponse,
      responseStream: false,
      options: {},
    },
    addSecondFactorToLoginPolicy: {
      name: "AddSecondFactorToLoginPolicy",
      requestType: AddSecondFactorToLoginPolicyRequest,
      requestStream: false,
      responseType: AddSecondFactorToLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
    removeSecondFactorFromLoginPolicy: {
      name: "RemoveSecondFactorFromLoginPolicy",
      requestType: RemoveSecondFactorFromLoginPolicyRequest,
      requestStream: false,
      responseType: RemoveSecondFactorFromLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
    listLoginPolicyMultiFactors: {
      name: "ListLoginPolicyMultiFactors",
      requestType: ListLoginPolicyMultiFactorsRequest,
      requestStream: false,
      responseType: ListLoginPolicyMultiFactorsResponse,
      responseStream: false,
      options: {},
    },
    addMultiFactorToLoginPolicy: {
      name: "AddMultiFactorToLoginPolicy",
      requestType: AddMultiFactorToLoginPolicyRequest,
      requestStream: false,
      responseType: AddMultiFactorToLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
    removeMultiFactorFromLoginPolicy: {
      name: "RemoveMultiFactorFromLoginPolicy",
      requestType: RemoveMultiFactorFromLoginPolicyRequest,
      requestStream: false,
      responseType: RemoveMultiFactorFromLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
    getPasswordComplexityPolicy: {
      name: "GetPasswordComplexityPolicy",
      requestType: GetPasswordComplexityPolicyRequest,
      requestStream: false,
      responseType: GetPasswordComplexityPolicyResponse,
      responseStream: false,
      options: {},
    },
    updatePasswordComplexityPolicy: {
      name: "UpdatePasswordComplexityPolicy",
      requestType: UpdatePasswordComplexityPolicyRequest,
      requestStream: false,
      responseType: UpdatePasswordComplexityPolicyResponse,
      responseStream: false,
      options: {},
    },
    getPasswordAgePolicy: {
      name: "GetPasswordAgePolicy",
      requestType: GetPasswordAgePolicyRequest,
      requestStream: false,
      responseType: GetPasswordAgePolicyResponse,
      responseStream: false,
      options: {},
    },
    updatePasswordAgePolicy: {
      name: "UpdatePasswordAgePolicy",
      requestType: UpdatePasswordAgePolicyRequest,
      requestStream: false,
      responseType: UpdatePasswordAgePolicyResponse,
      responseStream: false,
      options: {},
    },
    getLockoutPolicy: {
      name: "GetLockoutPolicy",
      requestType: GetLockoutPolicyRequest,
      requestStream: false,
      responseType: GetLockoutPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateLockoutPolicy: {
      name: "UpdateLockoutPolicy",
      requestType: UpdateLockoutPolicyRequest,
      requestStream: false,
      responseType: UpdateLockoutPolicyResponse,
      responseStream: false,
      options: {},
    },
    getPrivacyPolicy: {
      name: "GetPrivacyPolicy",
      requestType: GetPrivacyPolicyRequest,
      requestStream: false,
      responseType: GetPrivacyPolicyResponse,
      responseStream: false,
      options: {},
    },
    updatePrivacyPolicy: {
      name: "UpdatePrivacyPolicy",
      requestType: UpdatePrivacyPolicyRequest,
      requestStream: false,
      responseType: UpdatePrivacyPolicyResponse,
      responseStream: false,
      options: {},
    },
    addNotificationPolicy: {
      name: "AddNotificationPolicy",
      requestType: AddNotificationPolicyRequest,
      requestStream: false,
      responseType: AddNotificationPolicyResponse,
      responseStream: false,
      options: {},
    },
    getNotificationPolicy: {
      name: "GetNotificationPolicy",
      requestType: GetNotificationPolicyRequest,
      requestStream: false,
      responseType: GetNotificationPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateNotificationPolicy: {
      name: "UpdateNotificationPolicy",
      requestType: UpdateNotificationPolicyRequest,
      requestStream: false,
      responseType: UpdateNotificationPolicyResponse,
      responseStream: false,
      options: {},
    },
    getDefaultInitMessageText: {
      name: "GetDefaultInitMessageText",
      requestType: GetDefaultInitMessageTextRequest,
      requestStream: false,
      responseType: GetDefaultInitMessageTextResponse,
      responseStream: false,
      options: {},
    },
    getCustomInitMessageText: {
      name: "GetCustomInitMessageText",
      requestType: GetCustomInitMessageTextRequest,
      requestStream: false,
      responseType: GetCustomInitMessageTextResponse,
      responseStream: false,
      options: {},
    },
    setDefaultInitMessageText: {
      name: "SetDefaultInitMessageText",
      requestType: SetDefaultInitMessageTextRequest,
      requestStream: false,
      responseType: SetDefaultInitMessageTextResponse,
      responseStream: false,
      options: {},
    },
    resetCustomInitMessageTextToDefault: {
      name: "ResetCustomInitMessageTextToDefault",
      requestType: ResetCustomInitMessageTextToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomInitMessageTextToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getDefaultPasswordResetMessageText: {
      name: "GetDefaultPasswordResetMessageText",
      requestType: GetDefaultPasswordResetMessageTextRequest,
      requestStream: false,
      responseType: GetDefaultPasswordResetMessageTextResponse,
      responseStream: false,
      options: {},
    },
    getCustomPasswordResetMessageText: {
      name: "GetCustomPasswordResetMessageText",
      requestType: GetCustomPasswordResetMessageTextRequest,
      requestStream: false,
      responseType: GetCustomPasswordResetMessageTextResponse,
      responseStream: false,
      options: {},
    },
    setDefaultPasswordResetMessageText: {
      name: "SetDefaultPasswordResetMessageText",
      requestType: SetDefaultPasswordResetMessageTextRequest,
      requestStream: false,
      responseType: SetDefaultPasswordResetMessageTextResponse,
      responseStream: false,
      options: {},
    },
    resetCustomPasswordResetMessageTextToDefault: {
      name: "ResetCustomPasswordResetMessageTextToDefault",
      requestType: ResetCustomPasswordResetMessageTextToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomPasswordResetMessageTextToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getDefaultVerifyEmailMessageText: {
      name: "GetDefaultVerifyEmailMessageText",
      requestType: GetDefaultVerifyEmailMessageTextRequest,
      requestStream: false,
      responseType: GetDefaultVerifyEmailMessageTextResponse,
      responseStream: false,
      options: {},
    },
    getCustomVerifyEmailMessageText: {
      name: "GetCustomVerifyEmailMessageText",
      requestType: GetCustomVerifyEmailMessageTextRequest,
      requestStream: false,
      responseType: GetCustomVerifyEmailMessageTextResponse,
      responseStream: false,
      options: {},
    },
    setDefaultVerifyEmailMessageText: {
      name: "SetDefaultVerifyEmailMessageText",
      requestType: SetDefaultVerifyEmailMessageTextRequest,
      requestStream: false,
      responseType: SetDefaultVerifyEmailMessageTextResponse,
      responseStream: false,
      options: {},
    },
    resetCustomVerifyEmailMessageTextToDefault: {
      name: "ResetCustomVerifyEmailMessageTextToDefault",
      requestType: ResetCustomVerifyEmailMessageTextToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomVerifyEmailMessageTextToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getDefaultVerifyPhoneMessageText: {
      name: "GetDefaultVerifyPhoneMessageText",
      requestType: GetDefaultVerifyPhoneMessageTextRequest,
      requestStream: false,
      responseType: GetDefaultVerifyPhoneMessageTextResponse,
      responseStream: false,
      options: {},
    },
    getCustomVerifyPhoneMessageText: {
      name: "GetCustomVerifyPhoneMessageText",
      requestType: GetCustomVerifyPhoneMessageTextRequest,
      requestStream: false,
      responseType: GetCustomVerifyPhoneMessageTextResponse,
      responseStream: false,
      options: {},
    },
    setDefaultVerifyPhoneMessageText: {
      name: "SetDefaultVerifyPhoneMessageText",
      requestType: SetDefaultVerifyPhoneMessageTextRequest,
      requestStream: false,
      responseType: SetDefaultVerifyPhoneMessageTextResponse,
      responseStream: false,
      options: {},
    },
    resetCustomVerifyPhoneMessageTextToDefault: {
      name: "ResetCustomVerifyPhoneMessageTextToDefault",
      requestType: ResetCustomVerifyPhoneMessageTextToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomVerifyPhoneMessageTextToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getDefaultDomainClaimedMessageText: {
      name: "GetDefaultDomainClaimedMessageText",
      requestType: GetDefaultDomainClaimedMessageTextRequest,
      requestStream: false,
      responseType: GetDefaultDomainClaimedMessageTextResponse,
      responseStream: false,
      options: {},
    },
    getCustomDomainClaimedMessageText: {
      name: "GetCustomDomainClaimedMessageText",
      requestType: GetCustomDomainClaimedMessageTextRequest,
      requestStream: false,
      responseType: GetCustomDomainClaimedMessageTextResponse,
      responseStream: false,
      options: {},
    },
    setDefaultDomainClaimedMessageText: {
      name: "SetDefaultDomainClaimedMessageText",
      requestType: SetDefaultDomainClaimedMessageTextRequest,
      requestStream: false,
      responseType: SetDefaultDomainClaimedMessageTextResponse,
      responseStream: false,
      options: {},
    },
    resetCustomDomainClaimedMessageTextToDefault: {
      name: "ResetCustomDomainClaimedMessageTextToDefault",
      requestType: ResetCustomDomainClaimedMessageTextToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomDomainClaimedMessageTextToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getDefaultPasswordlessRegistrationMessageText: {
      name: "GetDefaultPasswordlessRegistrationMessageText",
      requestType: GetDefaultPasswordlessRegistrationMessageTextRequest,
      requestStream: false,
      responseType: GetDefaultPasswordlessRegistrationMessageTextResponse,
      responseStream: false,
      options: {},
    },
    getCustomPasswordlessRegistrationMessageText: {
      name: "GetCustomPasswordlessRegistrationMessageText",
      requestType: GetCustomPasswordlessRegistrationMessageTextRequest,
      requestStream: false,
      responseType: GetCustomPasswordlessRegistrationMessageTextResponse,
      responseStream: false,
      options: {},
    },
    setDefaultPasswordlessRegistrationMessageText: {
      name: "SetDefaultPasswordlessRegistrationMessageText",
      requestType: SetDefaultPasswordlessRegistrationMessageTextRequest,
      requestStream: false,
      responseType: SetDefaultPasswordlessRegistrationMessageTextResponse,
      responseStream: false,
      options: {},
    },
    resetCustomPasswordlessRegistrationMessageTextToDefault: {
      name: "ResetCustomPasswordlessRegistrationMessageTextToDefault",
      requestType: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getDefaultPasswordChangeMessageText: {
      name: "GetDefaultPasswordChangeMessageText",
      requestType: GetDefaultPasswordChangeMessageTextRequest,
      requestStream: false,
      responseType: GetDefaultPasswordChangeMessageTextResponse,
      responseStream: false,
      options: {},
    },
    getCustomPasswordChangeMessageText: {
      name: "GetCustomPasswordChangeMessageText",
      requestType: GetCustomPasswordChangeMessageTextRequest,
      requestStream: false,
      responseType: GetCustomPasswordChangeMessageTextResponse,
      responseStream: false,
      options: {},
    },
    setDefaultPasswordChangeMessageText: {
      name: "SetDefaultPasswordChangeMessageText",
      requestType: SetDefaultPasswordChangeMessageTextRequest,
      requestStream: false,
      responseType: SetDefaultPasswordChangeMessageTextResponse,
      responseStream: false,
      options: {},
    },
    resetCustomPasswordChangeMessageTextToDefault: {
      name: "ResetCustomPasswordChangeMessageTextToDefault",
      requestType: ResetCustomPasswordChangeMessageTextToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomPasswordChangeMessageTextToDefaultResponse,
      responseStream: false,
      options: {},
    },
    getDefaultLoginTexts: {
      name: "GetDefaultLoginTexts",
      requestType: GetDefaultLoginTextsRequest,
      requestStream: false,
      responseType: GetDefaultLoginTextsResponse,
      responseStream: false,
      options: {},
    },
    getCustomLoginTexts: {
      name: "GetCustomLoginTexts",
      requestType: GetCustomLoginTextsRequest,
      requestStream: false,
      responseType: GetCustomLoginTextsResponse,
      responseStream: false,
      options: {},
    },
    setCustomLoginText: {
      name: "SetCustomLoginText",
      requestType: SetCustomLoginTextsRequest,
      requestStream: false,
      responseType: SetCustomLoginTextsResponse,
      responseStream: false,
      options: {},
    },
    resetCustomLoginTextToDefault: {
      name: "ResetCustomLoginTextToDefault",
      requestType: ResetCustomLoginTextsToDefaultRequest,
      requestStream: false,
      responseType: ResetCustomLoginTextsToDefaultResponse,
      responseStream: false,
      options: {},
    },
    listIAMMemberRoles: {
      name: "ListIAMMemberRoles",
      requestType: ListIAMMemberRolesRequest,
      requestStream: false,
      responseType: ListIAMMemberRolesResponse,
      responseStream: false,
      options: {},
    },
    listIAMMembers: {
      name: "ListIAMMembers",
      requestType: ListIAMMembersRequest,
      requestStream: false,
      responseType: ListIAMMembersResponse,
      responseStream: false,
      options: {},
    },
    /**
     * Adds a user to the membership list of ZITADEL with the given roles
     * undefined roles will be dropped
     */
    addIAMMember: {
      name: "AddIAMMember",
      requestType: AddIAMMemberRequest,
      requestStream: false,
      responseType: AddIAMMemberResponse,
      responseStream: false,
      options: {},
    },
    updateIAMMember: {
      name: "UpdateIAMMember",
      requestType: UpdateIAMMemberRequest,
      requestStream: false,
      responseType: UpdateIAMMemberResponse,
      responseStream: false,
      options: {},
    },
    removeIAMMember: {
      name: "RemoveIAMMember",
      requestType: RemoveIAMMemberRequest,
      requestStream: false,
      responseType: RemoveIAMMemberResponse,
      responseStream: false,
      options: {},
    },
    listViews: {
      name: "ListViews",
      requestType: ListViewsRequest,
      requestStream: false,
      responseType: ListViewsResponse,
      responseStream: false,
      options: {},
    },
    listFailedEvents: {
      name: "ListFailedEvents",
      requestType: ListFailedEventsRequest,
      requestStream: false,
      responseType: ListFailedEventsResponse,
      responseStream: false,
      options: {},
    },
    removeFailedEvent: {
      name: "RemoveFailedEvent",
      requestType: RemoveFailedEventRequest,
      requestStream: false,
      responseType: RemoveFailedEventResponse,
      responseStream: false,
      options: {},
    },
    /** Imports data into an instance and creates different objects */
    importData: {
      name: "ImportData",
      requestType: ImportDataRequest,
      requestStream: false,
      responseType: ImportDataResponse,
      responseStream: false,
      options: {},
    },
    exportData: {
      name: "ExportData",
      requestType: ExportDataRequest,
      requestStream: false,
      responseType: ExportDataResponse,
      responseStream: false,
      options: {},
    },
    listEventTypes: {
      name: "ListEventTypes",
      requestType: ListEventTypesRequest,
      requestStream: false,
      responseType: ListEventTypesResponse,
      responseStream: false,
      options: {},
    },
    listEvents: {
      name: "ListEvents",
      requestType: ListEventsRequest,
      requestStream: false,
      responseType: ListEventsResponse,
      responseStream: false,
      options: {},
    },
    listAggregateTypes: {
      name: "ListAggregateTypes",
      requestType: ListAggregateTypesRequest,
      requestStream: false,
      responseType: ListAggregateTypesResponse,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface AdminServiceImplementation<CallContextExt = {}> {
  healthz(request: HealthzRequest, context: CallContext & CallContextExt): Promise<DeepPartial<HealthzResponse>>;
  getSupportedLanguages(
    request: GetSupportedLanguagesRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetSupportedLanguagesResponse>>;
  setDefaultLanguage(
    request: SetDefaultLanguageRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultLanguageResponse>>;
  getDefaultLanguage(
    request: GetDefaultLanguageRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultLanguageResponse>>;
  getMyInstance(
    request: GetMyInstanceRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyInstanceResponse>>;
  listInstanceDomains(
    request: ListInstanceDomainsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListInstanceDomainsResponse>>;
  listSecretGenerators(
    request: ListSecretGeneratorsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListSecretGeneratorsResponse>>;
  getSecretGenerator(
    request: GetSecretGeneratorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetSecretGeneratorResponse>>;
  updateSecretGenerator(
    request: UpdateSecretGeneratorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateSecretGeneratorResponse>>;
  getSMTPConfig(
    request: GetSMTPConfigRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetSMTPConfigResponse>>;
  addSMTPConfig(
    request: AddSMTPConfigRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddSMTPConfigResponse>>;
  updateSMTPConfig(
    request: UpdateSMTPConfigRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateSMTPConfigResponse>>;
  updateSMTPConfigPassword(
    request: UpdateSMTPConfigPasswordRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateSMTPConfigPasswordResponse>>;
  removeSMTPConfig(
    request: RemoveSMTPConfigRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveSMTPConfigResponse>>;
  listSMSProviders(
    request: ListSMSProvidersRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListSMSProvidersResponse>>;
  getSMSProvider(
    request: GetSMSProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetSMSProviderResponse>>;
  addSMSProviderTwilio(
    request: AddSMSProviderTwilioRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddSMSProviderTwilioResponse>>;
  updateSMSProviderTwilio(
    request: UpdateSMSProviderTwilioRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateSMSProviderTwilioResponse>>;
  updateSMSProviderTwilioToken(
    request: UpdateSMSProviderTwilioTokenRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateSMSProviderTwilioTokenResponse>>;
  activateSMSProvider(
    request: ActivateSMSProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ActivateSMSProviderResponse>>;
  deactivateSMSProvider(
    request: DeactivateSMSProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeactivateSMSProviderResponse>>;
  removeSMSProvider(
    request: RemoveSMSProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveSMSProviderResponse>>;
  getOIDCSettings(
    request: GetOIDCSettingsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetOIDCSettingsResponse>>;
  addOIDCSettings(
    request: AddOIDCSettingsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddOIDCSettingsResponse>>;
  updateOIDCSettings(
    request: UpdateOIDCSettingsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateOIDCSettingsResponse>>;
  getFileSystemNotificationProvider(
    request: GetFileSystemNotificationProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetFileSystemNotificationProviderResponse>>;
  getLogNotificationProvider(
    request: GetLogNotificationProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetLogNotificationProviderResponse>>;
  getSecurityPolicy(
    request: GetSecurityPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetSecurityPolicyResponse>>;
  setSecurityPolicy(
    request: SetSecurityPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetSecurityPolicyResponse>>;
  getOrgByID(
    request: GetOrgByIDRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetOrgByIDResponse>>;
  isOrgUnique(
    request: IsOrgUniqueRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<IsOrgUniqueResponse>>;
  setDefaultOrg(
    request: SetDefaultOrgRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultOrgResponse>>;
  getDefaultOrg(
    request: GetDefaultOrgRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultOrgResponse>>;
  listOrgs(request: ListOrgsRequest, context: CallContext & CallContextExt): Promise<DeepPartial<ListOrgsResponse>>;
  setUpOrg(request: SetUpOrgRequest, context: CallContext & CallContextExt): Promise<DeepPartial<SetUpOrgResponse>>;
  removeOrg(request: RemoveOrgRequest, context: CallContext & CallContextExt): Promise<DeepPartial<RemoveOrgResponse>>;
  getIDPByID(
    request: GetIDPByIDRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetIDPByIDResponse>>;
  listIDPs(request: ListIDPsRequest, context: CallContext & CallContextExt): Promise<DeepPartial<ListIDPsResponse>>;
  addOIDCIDP(
    request: AddOIDCIDPRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddOIDCIDPResponse>>;
  addJWTIDP(request: AddJWTIDPRequest, context: CallContext & CallContextExt): Promise<DeepPartial<AddJWTIDPResponse>>;
  updateIDP(request: UpdateIDPRequest, context: CallContext & CallContextExt): Promise<DeepPartial<UpdateIDPResponse>>;
  deactivateIDP(
    request: DeactivateIDPRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeactivateIDPResponse>>;
  reactivateIDP(
    request: ReactivateIDPRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ReactivateIDPResponse>>;
  removeIDP(request: RemoveIDPRequest, context: CallContext & CallContextExt): Promise<DeepPartial<RemoveIDPResponse>>;
  updateIDPOIDCConfig(
    request: UpdateIDPOIDCConfigRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateIDPOIDCConfigResponse>>;
  updateIDPJWTConfig(
    request: UpdateIDPJWTConfigRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateIDPJWTConfigResponse>>;
  /**
   * Returns all identity providers, which match the query
   * Limit should always be set, there is a default limit set by the service
   */
  listProviders(
    request: ListProvidersRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListProvidersResponse>>;
  /** Returns an identity provider of the instance */
  getProviderByID(
    request: GetProviderByIDRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetProviderByIDResponse>>;
  /** Add a new OAuth2 identity provider on the instance */
  addGenericOAuthProvider(
    request: AddGenericOAuthProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddGenericOAuthProviderResponse>>;
  /** Change an existing OAuth2 identity provider on the instance */
  updateGenericOAuthProvider(
    request: UpdateGenericOAuthProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateGenericOAuthProviderResponse>>;
  /** Add a new OIDC identity provider on the instance */
  addGenericOIDCProvider(
    request: AddGenericOIDCProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddGenericOIDCProviderResponse>>;
  /** Change an existing OIDC identity provider on the instance */
  updateGenericOIDCProvider(
    request: UpdateGenericOIDCProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateGenericOIDCProviderResponse>>;
  /** Add a new JWT identity provider on the instance */
  addJWTProvider(
    request: AddJWTProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddJWTProviderResponse>>;
  /** Change an existing JWT identity provider on the instance */
  updateJWTProvider(
    request: UpdateJWTProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateJWTProviderResponse>>;
  /** Add a new GitHub identity provider on the instance */
  addGitHubProvider(
    request: AddGitHubProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddGitHubProviderResponse>>;
  /** Change an existing GitHub identity provider on the instance */
  updateGitHubProvider(
    request: UpdateGitHubProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateGitHubProviderResponse>>;
  /** Add a new GitHub Enterprise Server identity provider on the instance */
  addGitHubEnterpriseServerProvider(
    request: AddGitHubEnterpriseServerProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddGitHubEnterpriseServerProviderResponse>>;
  /** Change an existing GitHub Enterprise Server identity provider on the instance */
  updateGitHubEnterpriseServerProvider(
    request: UpdateGitHubEnterpriseServerProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateGitHubEnterpriseServerProviderResponse>>;
  /** Add a new Google identity provider on the instance */
  addGoogleProvider(
    request: AddGoogleProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddGoogleProviderResponse>>;
  /** Change an existing Google identity provider on the instance */
  updateGoogleProvider(
    request: UpdateGoogleProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateGoogleProviderResponse>>;
  /** Add a new LDAP identity provider on the instance */
  addLDAPProvider(
    request: AddLDAPProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddLDAPProviderResponse>>;
  /** Change an existing LDAP identity provider on the instance */
  updateLDAPProvider(
    request: UpdateLDAPProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateLDAPProviderResponse>>;
  /**
   * Remove an identity provider
   * Will remove all linked providers of this configuration on the users
   */
  deleteProvider(
    request: DeleteProviderRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeleteProviderResponse>>;
  getOrgIAMPolicy(
    request: GetOrgIAMPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetOrgIAMPolicyResponse>>;
  updateOrgIAMPolicy(
    request: UpdateOrgIAMPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateOrgIAMPolicyResponse>>;
  getCustomOrgIAMPolicy(
    request: GetCustomOrgIAMPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomOrgIAMPolicyResponse>>;
  addCustomOrgIAMPolicy(
    request: AddCustomOrgIAMPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddCustomOrgIAMPolicyResponse>>;
  updateCustomOrgIAMPolicy(
    request: UpdateCustomOrgIAMPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateCustomOrgIAMPolicyResponse>>;
  resetCustomOrgIAMPolicyToDefault(
    request: ResetCustomOrgIAMPolicyToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomOrgIAMPolicyToDefaultResponse>>;
  getDomainPolicy(
    request: GetDomainPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDomainPolicyResponse>>;
  updateDomainPolicy(
    request: UpdateDomainPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateDomainPolicyResponse>>;
  getCustomDomainPolicy(
    request: GetCustomDomainPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomDomainPolicyResponse>>;
  addCustomDomainPolicy(
    request: AddCustomDomainPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddCustomDomainPolicyResponse>>;
  updateCustomDomainPolicy(
    request: UpdateCustomDomainPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateCustomDomainPolicyResponse>>;
  resetCustomDomainPolicyToDefault(
    request: ResetCustomDomainPolicyToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomDomainPolicyToDefaultResponse>>;
  getLabelPolicy(
    request: GetLabelPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetLabelPolicyResponse>>;
  getPreviewLabelPolicy(
    request: GetPreviewLabelPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetPreviewLabelPolicyResponse>>;
  updateLabelPolicy(
    request: UpdateLabelPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateLabelPolicyResponse>>;
  activateLabelPolicy(
    request: ActivateLabelPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ActivateLabelPolicyResponse>>;
  removeLabelPolicyLogo(
    request: RemoveLabelPolicyLogoRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveLabelPolicyLogoResponse>>;
  removeLabelPolicyLogoDark(
    request: RemoveLabelPolicyLogoDarkRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveLabelPolicyLogoDarkResponse>>;
  removeLabelPolicyIcon(
    request: RemoveLabelPolicyIconRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveLabelPolicyIconResponse>>;
  removeLabelPolicyIconDark(
    request: RemoveLabelPolicyIconDarkRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveLabelPolicyIconDarkResponse>>;
  removeLabelPolicyFont(
    request: RemoveLabelPolicyFontRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveLabelPolicyFontResponse>>;
  getLoginPolicy(
    request: GetLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetLoginPolicyResponse>>;
  updateLoginPolicy(
    request: UpdateLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateLoginPolicyResponse>>;
  listLoginPolicyIDPs(
    request: ListLoginPolicyIDPsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListLoginPolicyIDPsResponse>>;
  addIDPToLoginPolicy(
    request: AddIDPToLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddIDPToLoginPolicyResponse>>;
  removeIDPFromLoginPolicy(
    request: RemoveIDPFromLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveIDPFromLoginPolicyResponse>>;
  listLoginPolicySecondFactors(
    request: ListLoginPolicySecondFactorsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListLoginPolicySecondFactorsResponse>>;
  addSecondFactorToLoginPolicy(
    request: AddSecondFactorToLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddSecondFactorToLoginPolicyResponse>>;
  removeSecondFactorFromLoginPolicy(
    request: RemoveSecondFactorFromLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveSecondFactorFromLoginPolicyResponse>>;
  listLoginPolicyMultiFactors(
    request: ListLoginPolicyMultiFactorsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListLoginPolicyMultiFactorsResponse>>;
  addMultiFactorToLoginPolicy(
    request: AddMultiFactorToLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddMultiFactorToLoginPolicyResponse>>;
  removeMultiFactorFromLoginPolicy(
    request: RemoveMultiFactorFromLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveMultiFactorFromLoginPolicyResponse>>;
  getPasswordComplexityPolicy(
    request: GetPasswordComplexityPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetPasswordComplexityPolicyResponse>>;
  updatePasswordComplexityPolicy(
    request: UpdatePasswordComplexityPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdatePasswordComplexityPolicyResponse>>;
  getPasswordAgePolicy(
    request: GetPasswordAgePolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetPasswordAgePolicyResponse>>;
  updatePasswordAgePolicy(
    request: UpdatePasswordAgePolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdatePasswordAgePolicyResponse>>;
  getLockoutPolicy(
    request: GetLockoutPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetLockoutPolicyResponse>>;
  updateLockoutPolicy(
    request: UpdateLockoutPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateLockoutPolicyResponse>>;
  getPrivacyPolicy(
    request: GetPrivacyPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetPrivacyPolicyResponse>>;
  updatePrivacyPolicy(
    request: UpdatePrivacyPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdatePrivacyPolicyResponse>>;
  addNotificationPolicy(
    request: AddNotificationPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddNotificationPolicyResponse>>;
  getNotificationPolicy(
    request: GetNotificationPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetNotificationPolicyResponse>>;
  updateNotificationPolicy(
    request: UpdateNotificationPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateNotificationPolicyResponse>>;
  getDefaultInitMessageText(
    request: GetDefaultInitMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultInitMessageTextResponse>>;
  getCustomInitMessageText(
    request: GetCustomInitMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomInitMessageTextResponse>>;
  setDefaultInitMessageText(
    request: SetDefaultInitMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultInitMessageTextResponse>>;
  resetCustomInitMessageTextToDefault(
    request: ResetCustomInitMessageTextToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomInitMessageTextToDefaultResponse>>;
  getDefaultPasswordResetMessageText(
    request: GetDefaultPasswordResetMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultPasswordResetMessageTextResponse>>;
  getCustomPasswordResetMessageText(
    request: GetCustomPasswordResetMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomPasswordResetMessageTextResponse>>;
  setDefaultPasswordResetMessageText(
    request: SetDefaultPasswordResetMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultPasswordResetMessageTextResponse>>;
  resetCustomPasswordResetMessageTextToDefault(
    request: ResetCustomPasswordResetMessageTextToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomPasswordResetMessageTextToDefaultResponse>>;
  getDefaultVerifyEmailMessageText(
    request: GetDefaultVerifyEmailMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultVerifyEmailMessageTextResponse>>;
  getCustomVerifyEmailMessageText(
    request: GetCustomVerifyEmailMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomVerifyEmailMessageTextResponse>>;
  setDefaultVerifyEmailMessageText(
    request: SetDefaultVerifyEmailMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultVerifyEmailMessageTextResponse>>;
  resetCustomVerifyEmailMessageTextToDefault(
    request: ResetCustomVerifyEmailMessageTextToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomVerifyEmailMessageTextToDefaultResponse>>;
  getDefaultVerifyPhoneMessageText(
    request: GetDefaultVerifyPhoneMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultVerifyPhoneMessageTextResponse>>;
  getCustomVerifyPhoneMessageText(
    request: GetCustomVerifyPhoneMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomVerifyPhoneMessageTextResponse>>;
  setDefaultVerifyPhoneMessageText(
    request: SetDefaultVerifyPhoneMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultVerifyPhoneMessageTextResponse>>;
  resetCustomVerifyPhoneMessageTextToDefault(
    request: ResetCustomVerifyPhoneMessageTextToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomVerifyPhoneMessageTextToDefaultResponse>>;
  getDefaultDomainClaimedMessageText(
    request: GetDefaultDomainClaimedMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultDomainClaimedMessageTextResponse>>;
  getCustomDomainClaimedMessageText(
    request: GetCustomDomainClaimedMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomDomainClaimedMessageTextResponse>>;
  setDefaultDomainClaimedMessageText(
    request: SetDefaultDomainClaimedMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultDomainClaimedMessageTextResponse>>;
  resetCustomDomainClaimedMessageTextToDefault(
    request: ResetCustomDomainClaimedMessageTextToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomDomainClaimedMessageTextToDefaultResponse>>;
  getDefaultPasswordlessRegistrationMessageText(
    request: GetDefaultPasswordlessRegistrationMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultPasswordlessRegistrationMessageTextResponse>>;
  getCustomPasswordlessRegistrationMessageText(
    request: GetCustomPasswordlessRegistrationMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomPasswordlessRegistrationMessageTextResponse>>;
  setDefaultPasswordlessRegistrationMessageText(
    request: SetDefaultPasswordlessRegistrationMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultPasswordlessRegistrationMessageTextResponse>>;
  resetCustomPasswordlessRegistrationMessageTextToDefault(
    request: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse>>;
  getDefaultPasswordChangeMessageText(
    request: GetDefaultPasswordChangeMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultPasswordChangeMessageTextResponse>>;
  getCustomPasswordChangeMessageText(
    request: GetCustomPasswordChangeMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomPasswordChangeMessageTextResponse>>;
  setDefaultPasswordChangeMessageText(
    request: SetDefaultPasswordChangeMessageTextRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetDefaultPasswordChangeMessageTextResponse>>;
  resetCustomPasswordChangeMessageTextToDefault(
    request: ResetCustomPasswordChangeMessageTextToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomPasswordChangeMessageTextToDefaultResponse>>;
  getDefaultLoginTexts(
    request: GetDefaultLoginTextsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultLoginTextsResponse>>;
  getCustomLoginTexts(
    request: GetCustomLoginTextsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetCustomLoginTextsResponse>>;
  setCustomLoginText(
    request: SetCustomLoginTextsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetCustomLoginTextsResponse>>;
  resetCustomLoginTextToDefault(
    request: ResetCustomLoginTextsToDefaultRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResetCustomLoginTextsToDefaultResponse>>;
  listIAMMemberRoles(
    request: ListIAMMemberRolesRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListIAMMemberRolesResponse>>;
  listIAMMembers(
    request: ListIAMMembersRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListIAMMembersResponse>>;
  /**
   * Adds a user to the membership list of ZITADEL with the given roles
   * undefined roles will be dropped
   */
  addIAMMember(
    request: AddIAMMemberRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddIAMMemberResponse>>;
  updateIAMMember(
    request: UpdateIAMMemberRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateIAMMemberResponse>>;
  removeIAMMember(
    request: RemoveIAMMemberRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveIAMMemberResponse>>;
  listViews(request: ListViewsRequest, context: CallContext & CallContextExt): Promise<DeepPartial<ListViewsResponse>>;
  listFailedEvents(
    request: ListFailedEventsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListFailedEventsResponse>>;
  removeFailedEvent(
    request: RemoveFailedEventRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveFailedEventResponse>>;
  /** Imports data into an instance and creates different objects */
  importData(
    request: ImportDataRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ImportDataResponse>>;
  exportData(
    request: ExportDataRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ExportDataResponse>>;
  listEventTypes(
    request: ListEventTypesRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListEventTypesResponse>>;
  listEvents(
    request: ListEventsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListEventsResponse>>;
  listAggregateTypes(
    request: ListAggregateTypesRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListAggregateTypesResponse>>;
}

export interface AdminServiceClient<CallOptionsExt = {}> {
  healthz(request: DeepPartial<HealthzRequest>, options?: CallOptions & CallOptionsExt): Promise<HealthzResponse>;
  getSupportedLanguages(
    request: DeepPartial<GetSupportedLanguagesRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetSupportedLanguagesResponse>;
  setDefaultLanguage(
    request: DeepPartial<SetDefaultLanguageRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultLanguageResponse>;
  getDefaultLanguage(
    request: DeepPartial<GetDefaultLanguageRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultLanguageResponse>;
  getMyInstance(
    request: DeepPartial<GetMyInstanceRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyInstanceResponse>;
  listInstanceDomains(
    request: DeepPartial<ListInstanceDomainsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListInstanceDomainsResponse>;
  listSecretGenerators(
    request: DeepPartial<ListSecretGeneratorsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListSecretGeneratorsResponse>;
  getSecretGenerator(
    request: DeepPartial<GetSecretGeneratorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetSecretGeneratorResponse>;
  updateSecretGenerator(
    request: DeepPartial<UpdateSecretGeneratorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateSecretGeneratorResponse>;
  getSMTPConfig(
    request: DeepPartial<GetSMTPConfigRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetSMTPConfigResponse>;
  addSMTPConfig(
    request: DeepPartial<AddSMTPConfigRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddSMTPConfigResponse>;
  updateSMTPConfig(
    request: DeepPartial<UpdateSMTPConfigRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateSMTPConfigResponse>;
  updateSMTPConfigPassword(
    request: DeepPartial<UpdateSMTPConfigPasswordRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateSMTPConfigPasswordResponse>;
  removeSMTPConfig(
    request: DeepPartial<RemoveSMTPConfigRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveSMTPConfigResponse>;
  listSMSProviders(
    request: DeepPartial<ListSMSProvidersRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListSMSProvidersResponse>;
  getSMSProvider(
    request: DeepPartial<GetSMSProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetSMSProviderResponse>;
  addSMSProviderTwilio(
    request: DeepPartial<AddSMSProviderTwilioRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddSMSProviderTwilioResponse>;
  updateSMSProviderTwilio(
    request: DeepPartial<UpdateSMSProviderTwilioRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateSMSProviderTwilioResponse>;
  updateSMSProviderTwilioToken(
    request: DeepPartial<UpdateSMSProviderTwilioTokenRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateSMSProviderTwilioTokenResponse>;
  activateSMSProvider(
    request: DeepPartial<ActivateSMSProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ActivateSMSProviderResponse>;
  deactivateSMSProvider(
    request: DeepPartial<DeactivateSMSProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeactivateSMSProviderResponse>;
  removeSMSProvider(
    request: DeepPartial<RemoveSMSProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveSMSProviderResponse>;
  getOIDCSettings(
    request: DeepPartial<GetOIDCSettingsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetOIDCSettingsResponse>;
  addOIDCSettings(
    request: DeepPartial<AddOIDCSettingsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddOIDCSettingsResponse>;
  updateOIDCSettings(
    request: DeepPartial<UpdateOIDCSettingsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateOIDCSettingsResponse>;
  getFileSystemNotificationProvider(
    request: DeepPartial<GetFileSystemNotificationProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetFileSystemNotificationProviderResponse>;
  getLogNotificationProvider(
    request: DeepPartial<GetLogNotificationProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetLogNotificationProviderResponse>;
  getSecurityPolicy(
    request: DeepPartial<GetSecurityPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetSecurityPolicyResponse>;
  setSecurityPolicy(
    request: DeepPartial<SetSecurityPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetSecurityPolicyResponse>;
  getOrgByID(
    request: DeepPartial<GetOrgByIDRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetOrgByIDResponse>;
  isOrgUnique(
    request: DeepPartial<IsOrgUniqueRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<IsOrgUniqueResponse>;
  setDefaultOrg(
    request: DeepPartial<SetDefaultOrgRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultOrgResponse>;
  getDefaultOrg(
    request: DeepPartial<GetDefaultOrgRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultOrgResponse>;
  listOrgs(request: DeepPartial<ListOrgsRequest>, options?: CallOptions & CallOptionsExt): Promise<ListOrgsResponse>;
  setUpOrg(request: DeepPartial<SetUpOrgRequest>, options?: CallOptions & CallOptionsExt): Promise<SetUpOrgResponse>;
  removeOrg(request: DeepPartial<RemoveOrgRequest>, options?: CallOptions & CallOptionsExt): Promise<RemoveOrgResponse>;
  getIDPByID(
    request: DeepPartial<GetIDPByIDRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetIDPByIDResponse>;
  listIDPs(request: DeepPartial<ListIDPsRequest>, options?: CallOptions & CallOptionsExt): Promise<ListIDPsResponse>;
  addOIDCIDP(
    request: DeepPartial<AddOIDCIDPRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddOIDCIDPResponse>;
  addJWTIDP(request: DeepPartial<AddJWTIDPRequest>, options?: CallOptions & CallOptionsExt): Promise<AddJWTIDPResponse>;
  updateIDP(request: DeepPartial<UpdateIDPRequest>, options?: CallOptions & CallOptionsExt): Promise<UpdateIDPResponse>;
  deactivateIDP(
    request: DeepPartial<DeactivateIDPRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeactivateIDPResponse>;
  reactivateIDP(
    request: DeepPartial<ReactivateIDPRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ReactivateIDPResponse>;
  removeIDP(request: DeepPartial<RemoveIDPRequest>, options?: CallOptions & CallOptionsExt): Promise<RemoveIDPResponse>;
  updateIDPOIDCConfig(
    request: DeepPartial<UpdateIDPOIDCConfigRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateIDPOIDCConfigResponse>;
  updateIDPJWTConfig(
    request: DeepPartial<UpdateIDPJWTConfigRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateIDPJWTConfigResponse>;
  /**
   * Returns all identity providers, which match the query
   * Limit should always be set, there is a default limit set by the service
   */
  listProviders(
    request: DeepPartial<ListProvidersRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListProvidersResponse>;
  /** Returns an identity provider of the instance */
  getProviderByID(
    request: DeepPartial<GetProviderByIDRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetProviderByIDResponse>;
  /** Add a new OAuth2 identity provider on the instance */
  addGenericOAuthProvider(
    request: DeepPartial<AddGenericOAuthProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddGenericOAuthProviderResponse>;
  /** Change an existing OAuth2 identity provider on the instance */
  updateGenericOAuthProvider(
    request: DeepPartial<UpdateGenericOAuthProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateGenericOAuthProviderResponse>;
  /** Add a new OIDC identity provider on the instance */
  addGenericOIDCProvider(
    request: DeepPartial<AddGenericOIDCProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddGenericOIDCProviderResponse>;
  /** Change an existing OIDC identity provider on the instance */
  updateGenericOIDCProvider(
    request: DeepPartial<UpdateGenericOIDCProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateGenericOIDCProviderResponse>;
  /** Add a new JWT identity provider on the instance */
  addJWTProvider(
    request: DeepPartial<AddJWTProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddJWTProviderResponse>;
  /** Change an existing JWT identity provider on the instance */
  updateJWTProvider(
    request: DeepPartial<UpdateJWTProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateJWTProviderResponse>;
  /** Add a new GitHub identity provider on the instance */
  addGitHubProvider(
    request: DeepPartial<AddGitHubProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddGitHubProviderResponse>;
  /** Change an existing GitHub identity provider on the instance */
  updateGitHubProvider(
    request: DeepPartial<UpdateGitHubProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateGitHubProviderResponse>;
  /** Add a new GitHub Enterprise Server identity provider on the instance */
  addGitHubEnterpriseServerProvider(
    request: DeepPartial<AddGitHubEnterpriseServerProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddGitHubEnterpriseServerProviderResponse>;
  /** Change an existing GitHub Enterprise Server identity provider on the instance */
  updateGitHubEnterpriseServerProvider(
    request: DeepPartial<UpdateGitHubEnterpriseServerProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateGitHubEnterpriseServerProviderResponse>;
  /** Add a new Google identity provider on the instance */
  addGoogleProvider(
    request: DeepPartial<AddGoogleProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddGoogleProviderResponse>;
  /** Change an existing Google identity provider on the instance */
  updateGoogleProvider(
    request: DeepPartial<UpdateGoogleProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateGoogleProviderResponse>;
  /** Add a new LDAP identity provider on the instance */
  addLDAPProvider(
    request: DeepPartial<AddLDAPProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddLDAPProviderResponse>;
  /** Change an existing LDAP identity provider on the instance */
  updateLDAPProvider(
    request: DeepPartial<UpdateLDAPProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateLDAPProviderResponse>;
  /**
   * Remove an identity provider
   * Will remove all linked providers of this configuration on the users
   */
  deleteProvider(
    request: DeepPartial<DeleteProviderRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeleteProviderResponse>;
  getOrgIAMPolicy(
    request: DeepPartial<GetOrgIAMPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetOrgIAMPolicyResponse>;
  updateOrgIAMPolicy(
    request: DeepPartial<UpdateOrgIAMPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateOrgIAMPolicyResponse>;
  getCustomOrgIAMPolicy(
    request: DeepPartial<GetCustomOrgIAMPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomOrgIAMPolicyResponse>;
  addCustomOrgIAMPolicy(
    request: DeepPartial<AddCustomOrgIAMPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddCustomOrgIAMPolicyResponse>;
  updateCustomOrgIAMPolicy(
    request: DeepPartial<UpdateCustomOrgIAMPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateCustomOrgIAMPolicyResponse>;
  resetCustomOrgIAMPolicyToDefault(
    request: DeepPartial<ResetCustomOrgIAMPolicyToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomOrgIAMPolicyToDefaultResponse>;
  getDomainPolicy(
    request: DeepPartial<GetDomainPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDomainPolicyResponse>;
  updateDomainPolicy(
    request: DeepPartial<UpdateDomainPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateDomainPolicyResponse>;
  getCustomDomainPolicy(
    request: DeepPartial<GetCustomDomainPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomDomainPolicyResponse>;
  addCustomDomainPolicy(
    request: DeepPartial<AddCustomDomainPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddCustomDomainPolicyResponse>;
  updateCustomDomainPolicy(
    request: DeepPartial<UpdateCustomDomainPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateCustomDomainPolicyResponse>;
  resetCustomDomainPolicyToDefault(
    request: DeepPartial<ResetCustomDomainPolicyToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomDomainPolicyToDefaultResponse>;
  getLabelPolicy(
    request: DeepPartial<GetLabelPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetLabelPolicyResponse>;
  getPreviewLabelPolicy(
    request: DeepPartial<GetPreviewLabelPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetPreviewLabelPolicyResponse>;
  updateLabelPolicy(
    request: DeepPartial<UpdateLabelPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateLabelPolicyResponse>;
  activateLabelPolicy(
    request: DeepPartial<ActivateLabelPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ActivateLabelPolicyResponse>;
  removeLabelPolicyLogo(
    request: DeepPartial<RemoveLabelPolicyLogoRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveLabelPolicyLogoResponse>;
  removeLabelPolicyLogoDark(
    request: DeepPartial<RemoveLabelPolicyLogoDarkRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveLabelPolicyLogoDarkResponse>;
  removeLabelPolicyIcon(
    request: DeepPartial<RemoveLabelPolicyIconRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveLabelPolicyIconResponse>;
  removeLabelPolicyIconDark(
    request: DeepPartial<RemoveLabelPolicyIconDarkRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveLabelPolicyIconDarkResponse>;
  removeLabelPolicyFont(
    request: DeepPartial<RemoveLabelPolicyFontRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveLabelPolicyFontResponse>;
  getLoginPolicy(
    request: DeepPartial<GetLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetLoginPolicyResponse>;
  updateLoginPolicy(
    request: DeepPartial<UpdateLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateLoginPolicyResponse>;
  listLoginPolicyIDPs(
    request: DeepPartial<ListLoginPolicyIDPsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListLoginPolicyIDPsResponse>;
  addIDPToLoginPolicy(
    request: DeepPartial<AddIDPToLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddIDPToLoginPolicyResponse>;
  removeIDPFromLoginPolicy(
    request: DeepPartial<RemoveIDPFromLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveIDPFromLoginPolicyResponse>;
  listLoginPolicySecondFactors(
    request: DeepPartial<ListLoginPolicySecondFactorsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListLoginPolicySecondFactorsResponse>;
  addSecondFactorToLoginPolicy(
    request: DeepPartial<AddSecondFactorToLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddSecondFactorToLoginPolicyResponse>;
  removeSecondFactorFromLoginPolicy(
    request: DeepPartial<RemoveSecondFactorFromLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveSecondFactorFromLoginPolicyResponse>;
  listLoginPolicyMultiFactors(
    request: DeepPartial<ListLoginPolicyMultiFactorsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListLoginPolicyMultiFactorsResponse>;
  addMultiFactorToLoginPolicy(
    request: DeepPartial<AddMultiFactorToLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddMultiFactorToLoginPolicyResponse>;
  removeMultiFactorFromLoginPolicy(
    request: DeepPartial<RemoveMultiFactorFromLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveMultiFactorFromLoginPolicyResponse>;
  getPasswordComplexityPolicy(
    request: DeepPartial<GetPasswordComplexityPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetPasswordComplexityPolicyResponse>;
  updatePasswordComplexityPolicy(
    request: DeepPartial<UpdatePasswordComplexityPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdatePasswordComplexityPolicyResponse>;
  getPasswordAgePolicy(
    request: DeepPartial<GetPasswordAgePolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetPasswordAgePolicyResponse>;
  updatePasswordAgePolicy(
    request: DeepPartial<UpdatePasswordAgePolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdatePasswordAgePolicyResponse>;
  getLockoutPolicy(
    request: DeepPartial<GetLockoutPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetLockoutPolicyResponse>;
  updateLockoutPolicy(
    request: DeepPartial<UpdateLockoutPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateLockoutPolicyResponse>;
  getPrivacyPolicy(
    request: DeepPartial<GetPrivacyPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetPrivacyPolicyResponse>;
  updatePrivacyPolicy(
    request: DeepPartial<UpdatePrivacyPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdatePrivacyPolicyResponse>;
  addNotificationPolicy(
    request: DeepPartial<AddNotificationPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddNotificationPolicyResponse>;
  getNotificationPolicy(
    request: DeepPartial<GetNotificationPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetNotificationPolicyResponse>;
  updateNotificationPolicy(
    request: DeepPartial<UpdateNotificationPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateNotificationPolicyResponse>;
  getDefaultInitMessageText(
    request: DeepPartial<GetDefaultInitMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultInitMessageTextResponse>;
  getCustomInitMessageText(
    request: DeepPartial<GetCustomInitMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomInitMessageTextResponse>;
  setDefaultInitMessageText(
    request: DeepPartial<SetDefaultInitMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultInitMessageTextResponse>;
  resetCustomInitMessageTextToDefault(
    request: DeepPartial<ResetCustomInitMessageTextToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomInitMessageTextToDefaultResponse>;
  getDefaultPasswordResetMessageText(
    request: DeepPartial<GetDefaultPasswordResetMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultPasswordResetMessageTextResponse>;
  getCustomPasswordResetMessageText(
    request: DeepPartial<GetCustomPasswordResetMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomPasswordResetMessageTextResponse>;
  setDefaultPasswordResetMessageText(
    request: DeepPartial<SetDefaultPasswordResetMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultPasswordResetMessageTextResponse>;
  resetCustomPasswordResetMessageTextToDefault(
    request: DeepPartial<ResetCustomPasswordResetMessageTextToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomPasswordResetMessageTextToDefaultResponse>;
  getDefaultVerifyEmailMessageText(
    request: DeepPartial<GetDefaultVerifyEmailMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultVerifyEmailMessageTextResponse>;
  getCustomVerifyEmailMessageText(
    request: DeepPartial<GetCustomVerifyEmailMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomVerifyEmailMessageTextResponse>;
  setDefaultVerifyEmailMessageText(
    request: DeepPartial<SetDefaultVerifyEmailMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultVerifyEmailMessageTextResponse>;
  resetCustomVerifyEmailMessageTextToDefault(
    request: DeepPartial<ResetCustomVerifyEmailMessageTextToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomVerifyEmailMessageTextToDefaultResponse>;
  getDefaultVerifyPhoneMessageText(
    request: DeepPartial<GetDefaultVerifyPhoneMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultVerifyPhoneMessageTextResponse>;
  getCustomVerifyPhoneMessageText(
    request: DeepPartial<GetCustomVerifyPhoneMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomVerifyPhoneMessageTextResponse>;
  setDefaultVerifyPhoneMessageText(
    request: DeepPartial<SetDefaultVerifyPhoneMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultVerifyPhoneMessageTextResponse>;
  resetCustomVerifyPhoneMessageTextToDefault(
    request: DeepPartial<ResetCustomVerifyPhoneMessageTextToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomVerifyPhoneMessageTextToDefaultResponse>;
  getDefaultDomainClaimedMessageText(
    request: DeepPartial<GetDefaultDomainClaimedMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultDomainClaimedMessageTextResponse>;
  getCustomDomainClaimedMessageText(
    request: DeepPartial<GetCustomDomainClaimedMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomDomainClaimedMessageTextResponse>;
  setDefaultDomainClaimedMessageText(
    request: DeepPartial<SetDefaultDomainClaimedMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultDomainClaimedMessageTextResponse>;
  resetCustomDomainClaimedMessageTextToDefault(
    request: DeepPartial<ResetCustomDomainClaimedMessageTextToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomDomainClaimedMessageTextToDefaultResponse>;
  getDefaultPasswordlessRegistrationMessageText(
    request: DeepPartial<GetDefaultPasswordlessRegistrationMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultPasswordlessRegistrationMessageTextResponse>;
  getCustomPasswordlessRegistrationMessageText(
    request: DeepPartial<GetCustomPasswordlessRegistrationMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomPasswordlessRegistrationMessageTextResponse>;
  setDefaultPasswordlessRegistrationMessageText(
    request: DeepPartial<SetDefaultPasswordlessRegistrationMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultPasswordlessRegistrationMessageTextResponse>;
  resetCustomPasswordlessRegistrationMessageTextToDefault(
    request: DeepPartial<ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse>;
  getDefaultPasswordChangeMessageText(
    request: DeepPartial<GetDefaultPasswordChangeMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultPasswordChangeMessageTextResponse>;
  getCustomPasswordChangeMessageText(
    request: DeepPartial<GetCustomPasswordChangeMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomPasswordChangeMessageTextResponse>;
  setDefaultPasswordChangeMessageText(
    request: DeepPartial<SetDefaultPasswordChangeMessageTextRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetDefaultPasswordChangeMessageTextResponse>;
  resetCustomPasswordChangeMessageTextToDefault(
    request: DeepPartial<ResetCustomPasswordChangeMessageTextToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomPasswordChangeMessageTextToDefaultResponse>;
  getDefaultLoginTexts(
    request: DeepPartial<GetDefaultLoginTextsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultLoginTextsResponse>;
  getCustomLoginTexts(
    request: DeepPartial<GetCustomLoginTextsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetCustomLoginTextsResponse>;
  setCustomLoginText(
    request: DeepPartial<SetCustomLoginTextsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetCustomLoginTextsResponse>;
  resetCustomLoginTextToDefault(
    request: DeepPartial<ResetCustomLoginTextsToDefaultRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResetCustomLoginTextsToDefaultResponse>;
  listIAMMemberRoles(
    request: DeepPartial<ListIAMMemberRolesRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListIAMMemberRolesResponse>;
  listIAMMembers(
    request: DeepPartial<ListIAMMembersRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListIAMMembersResponse>;
  /**
   * Adds a user to the membership list of ZITADEL with the given roles
   * undefined roles will be dropped
   */
  addIAMMember(
    request: DeepPartial<AddIAMMemberRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddIAMMemberResponse>;
  updateIAMMember(
    request: DeepPartial<UpdateIAMMemberRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateIAMMemberResponse>;
  removeIAMMember(
    request: DeepPartial<RemoveIAMMemberRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveIAMMemberResponse>;
  listViews(request: DeepPartial<ListViewsRequest>, options?: CallOptions & CallOptionsExt): Promise<ListViewsResponse>;
  listFailedEvents(
    request: DeepPartial<ListFailedEventsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListFailedEventsResponse>;
  removeFailedEvent(
    request: DeepPartial<RemoveFailedEventRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveFailedEventResponse>;
  /** Imports data into an instance and creates different objects */
  importData(
    request: DeepPartial<ImportDataRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ImportDataResponse>;
  exportData(
    request: DeepPartial<ExportDataRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ExportDataResponse>;
  listEventTypes(
    request: DeepPartial<ListEventTypesRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListEventTypesResponse>;
  listEvents(
    request: DeepPartial<ListEventsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListEventsResponse>;
  listAggregateTypes(
    request: DeepPartial<ListAggregateTypesRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListAggregateTypesResponse>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

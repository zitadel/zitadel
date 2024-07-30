/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";
import { ObjectDetails, TextQueryMethod, textQueryMethodFromJSON, textQueryMethodToJSON } from "./object";

export const protobufPackage = "zitadel.idp.v1";

export enum IDPState {
  IDP_STATE_UNSPECIFIED = 0,
  IDP_STATE_ACTIVE = 1,
  IDP_STATE_INACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function iDPStateFromJSON(object: any): IDPState {
  switch (object) {
    case 0:
    case "IDP_STATE_UNSPECIFIED":
      return IDPState.IDP_STATE_UNSPECIFIED;
    case 1:
    case "IDP_STATE_ACTIVE":
      return IDPState.IDP_STATE_ACTIVE;
    case 2:
    case "IDP_STATE_INACTIVE":
      return IDPState.IDP_STATE_INACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IDPState.UNRECOGNIZED;
  }
}

export function iDPStateToJSON(object: IDPState): string {
  switch (object) {
    case IDPState.IDP_STATE_UNSPECIFIED:
      return "IDP_STATE_UNSPECIFIED";
    case IDPState.IDP_STATE_ACTIVE:
      return "IDP_STATE_ACTIVE";
    case IDPState.IDP_STATE_INACTIVE:
      return "IDP_STATE_INACTIVE";
    case IDPState.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum IDPStylingType {
  STYLING_TYPE_UNSPECIFIED = 0,
  STYLING_TYPE_GOOGLE = 1,
  UNRECOGNIZED = -1,
}

export function iDPStylingTypeFromJSON(object: any): IDPStylingType {
  switch (object) {
    case 0:
    case "STYLING_TYPE_UNSPECIFIED":
      return IDPStylingType.STYLING_TYPE_UNSPECIFIED;
    case 1:
    case "STYLING_TYPE_GOOGLE":
      return IDPStylingType.STYLING_TYPE_GOOGLE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IDPStylingType.UNRECOGNIZED;
  }
}

export function iDPStylingTypeToJSON(object: IDPStylingType): string {
  switch (object) {
    case IDPStylingType.STYLING_TYPE_UNSPECIFIED:
      return "STYLING_TYPE_UNSPECIFIED";
    case IDPStylingType.STYLING_TYPE_GOOGLE:
      return "STYLING_TYPE_GOOGLE";
    case IDPStylingType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** authorization framework of the identity provider */
export enum IDPType {
  IDP_TYPE_UNSPECIFIED = 0,
  IDP_TYPE_OIDC = 1,
  IDP_TYPE_JWT = 3,
  UNRECOGNIZED = -1,
}

export function iDPTypeFromJSON(object: any): IDPType {
  switch (object) {
    case 0:
    case "IDP_TYPE_UNSPECIFIED":
      return IDPType.IDP_TYPE_UNSPECIFIED;
    case 1:
    case "IDP_TYPE_OIDC":
      return IDPType.IDP_TYPE_OIDC;
    case 3:
    case "IDP_TYPE_JWT":
      return IDPType.IDP_TYPE_JWT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IDPType.UNRECOGNIZED;
  }
}

export function iDPTypeToJSON(object: IDPType): string {
  switch (object) {
    case IDPType.IDP_TYPE_UNSPECIFIED:
      return "IDP_TYPE_UNSPECIFIED";
    case IDPType.IDP_TYPE_OIDC:
      return "IDP_TYPE_OIDC";
    case IDPType.IDP_TYPE_JWT:
      return "IDP_TYPE_JWT";
    case IDPType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** the owner of the identity provider. */
export enum IDPOwnerType {
  IDP_OWNER_TYPE_UNSPECIFIED = 0,
  /** IDP_OWNER_TYPE_SYSTEM - system is managed by the ZITADEL administrators */
  IDP_OWNER_TYPE_SYSTEM = 1,
  /** IDP_OWNER_TYPE_ORG - org is managed by de organization administrators */
  IDP_OWNER_TYPE_ORG = 2,
  UNRECOGNIZED = -1,
}

export function iDPOwnerTypeFromJSON(object: any): IDPOwnerType {
  switch (object) {
    case 0:
    case "IDP_OWNER_TYPE_UNSPECIFIED":
      return IDPOwnerType.IDP_OWNER_TYPE_UNSPECIFIED;
    case 1:
    case "IDP_OWNER_TYPE_SYSTEM":
      return IDPOwnerType.IDP_OWNER_TYPE_SYSTEM;
    case 2:
    case "IDP_OWNER_TYPE_ORG":
      return IDPOwnerType.IDP_OWNER_TYPE_ORG;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IDPOwnerType.UNRECOGNIZED;
  }
}

export function iDPOwnerTypeToJSON(object: IDPOwnerType): string {
  switch (object) {
    case IDPOwnerType.IDP_OWNER_TYPE_UNSPECIFIED:
      return "IDP_OWNER_TYPE_UNSPECIFIED";
    case IDPOwnerType.IDP_OWNER_TYPE_SYSTEM:
      return "IDP_OWNER_TYPE_SYSTEM";
    case IDPOwnerType.IDP_OWNER_TYPE_ORG:
      return "IDP_OWNER_TYPE_ORG";
    case IDPOwnerType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum OIDCMappingField {
  OIDC_MAPPING_FIELD_UNSPECIFIED = 0,
  OIDC_MAPPING_FIELD_PREFERRED_USERNAME = 1,
  OIDC_MAPPING_FIELD_EMAIL = 2,
  UNRECOGNIZED = -1,
}

export function oIDCMappingFieldFromJSON(object: any): OIDCMappingField {
  switch (object) {
    case 0:
    case "OIDC_MAPPING_FIELD_UNSPECIFIED":
      return OIDCMappingField.OIDC_MAPPING_FIELD_UNSPECIFIED;
    case 1:
    case "OIDC_MAPPING_FIELD_PREFERRED_USERNAME":
      return OIDCMappingField.OIDC_MAPPING_FIELD_PREFERRED_USERNAME;
    case 2:
    case "OIDC_MAPPING_FIELD_EMAIL":
      return OIDCMappingField.OIDC_MAPPING_FIELD_EMAIL;
    case -1:
    case "UNRECOGNIZED":
    default:
      return OIDCMappingField.UNRECOGNIZED;
  }
}

export function oIDCMappingFieldToJSON(object: OIDCMappingField): string {
  switch (object) {
    case OIDCMappingField.OIDC_MAPPING_FIELD_UNSPECIFIED:
      return "OIDC_MAPPING_FIELD_UNSPECIFIED";
    case OIDCMappingField.OIDC_MAPPING_FIELD_PREFERRED_USERNAME:
      return "OIDC_MAPPING_FIELD_PREFERRED_USERNAME";
    case OIDCMappingField.OIDC_MAPPING_FIELD_EMAIL:
      return "OIDC_MAPPING_FIELD_EMAIL";
    case OIDCMappingField.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum IDPFieldName {
  IDP_FIELD_NAME_UNSPECIFIED = 0,
  IDP_FIELD_NAME_NAME = 1,
  UNRECOGNIZED = -1,
}

export function iDPFieldNameFromJSON(object: any): IDPFieldName {
  switch (object) {
    case 0:
    case "IDP_FIELD_NAME_UNSPECIFIED":
      return IDPFieldName.IDP_FIELD_NAME_UNSPECIFIED;
    case 1:
    case "IDP_FIELD_NAME_NAME":
      return IDPFieldName.IDP_FIELD_NAME_NAME;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IDPFieldName.UNRECOGNIZED;
  }
}

export function iDPFieldNameToJSON(object: IDPFieldName): string {
  switch (object) {
    case IDPFieldName.IDP_FIELD_NAME_UNSPECIFIED:
      return "IDP_FIELD_NAME_UNSPECIFIED";
    case IDPFieldName.IDP_FIELD_NAME_NAME:
      return "IDP_FIELD_NAME_NAME";
    case IDPFieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ProviderType {
  PROVIDER_TYPE_UNSPECIFIED = 0,
  PROVIDER_TYPE_OIDC = 1,
  PROVIDER_TYPE_JWT = 2,
  PROVIDER_TYPE_LDAP = 3,
  PROVIDER_TYPE_OAUTH = 4,
  PROVIDER_TYPE_AZURE_AD = 5,
  PROVIDER_TYPE_GITHUB = 6,
  PROVIDER_TYPE_GITHUB_ES = 7,
  PROVIDER_TYPE_GITLAB = 8,
  PROVIDER_TYPE_GITLAB_SELF_HOSTED = 9,
  PROVIDER_TYPE_GOOGLE = 10,
  PROVIDER_TYPE_APPLE = 11,
  PROVIDER_TYPE_SAML = 12,
  UNRECOGNIZED = -1,
}

export function providerTypeFromJSON(object: any): ProviderType {
  switch (object) {
    case 0:
    case "PROVIDER_TYPE_UNSPECIFIED":
      return ProviderType.PROVIDER_TYPE_UNSPECIFIED;
    case 1:
    case "PROVIDER_TYPE_OIDC":
      return ProviderType.PROVIDER_TYPE_OIDC;
    case 2:
    case "PROVIDER_TYPE_JWT":
      return ProviderType.PROVIDER_TYPE_JWT;
    case 3:
    case "PROVIDER_TYPE_LDAP":
      return ProviderType.PROVIDER_TYPE_LDAP;
    case 4:
    case "PROVIDER_TYPE_OAUTH":
      return ProviderType.PROVIDER_TYPE_OAUTH;
    case 5:
    case "PROVIDER_TYPE_AZURE_AD":
      return ProviderType.PROVIDER_TYPE_AZURE_AD;
    case 6:
    case "PROVIDER_TYPE_GITHUB":
      return ProviderType.PROVIDER_TYPE_GITHUB;
    case 7:
    case "PROVIDER_TYPE_GITHUB_ES":
      return ProviderType.PROVIDER_TYPE_GITHUB_ES;
    case 8:
    case "PROVIDER_TYPE_GITLAB":
      return ProviderType.PROVIDER_TYPE_GITLAB;
    case 9:
    case "PROVIDER_TYPE_GITLAB_SELF_HOSTED":
      return ProviderType.PROVIDER_TYPE_GITLAB_SELF_HOSTED;
    case 10:
    case "PROVIDER_TYPE_GOOGLE":
      return ProviderType.PROVIDER_TYPE_GOOGLE;
    case 11:
    case "PROVIDER_TYPE_APPLE":
      return ProviderType.PROVIDER_TYPE_APPLE;
    case 12:
    case "PROVIDER_TYPE_SAML":
      return ProviderType.PROVIDER_TYPE_SAML;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ProviderType.UNRECOGNIZED;
  }
}

export function providerTypeToJSON(object: ProviderType): string {
  switch (object) {
    case ProviderType.PROVIDER_TYPE_UNSPECIFIED:
      return "PROVIDER_TYPE_UNSPECIFIED";
    case ProviderType.PROVIDER_TYPE_OIDC:
      return "PROVIDER_TYPE_OIDC";
    case ProviderType.PROVIDER_TYPE_JWT:
      return "PROVIDER_TYPE_JWT";
    case ProviderType.PROVIDER_TYPE_LDAP:
      return "PROVIDER_TYPE_LDAP";
    case ProviderType.PROVIDER_TYPE_OAUTH:
      return "PROVIDER_TYPE_OAUTH";
    case ProviderType.PROVIDER_TYPE_AZURE_AD:
      return "PROVIDER_TYPE_AZURE_AD";
    case ProviderType.PROVIDER_TYPE_GITHUB:
      return "PROVIDER_TYPE_GITHUB";
    case ProviderType.PROVIDER_TYPE_GITHUB_ES:
      return "PROVIDER_TYPE_GITHUB_ES";
    case ProviderType.PROVIDER_TYPE_GITLAB:
      return "PROVIDER_TYPE_GITLAB";
    case ProviderType.PROVIDER_TYPE_GITLAB_SELF_HOSTED:
      return "PROVIDER_TYPE_GITLAB_SELF_HOSTED";
    case ProviderType.PROVIDER_TYPE_GOOGLE:
      return "PROVIDER_TYPE_GOOGLE";
    case ProviderType.PROVIDER_TYPE_APPLE:
      return "PROVIDER_TYPE_APPLE";
    case ProviderType.PROVIDER_TYPE_SAML:
      return "PROVIDER_TYPE_SAML";
    case ProviderType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum SAMLBinding {
  SAML_BINDING_UNSPECIFIED = 0,
  SAML_BINDING_POST = 1,
  SAML_BINDING_REDIRECT = 2,
  SAML_BINDING_ARTIFACT = 3,
  UNRECOGNIZED = -1,
}

export function sAMLBindingFromJSON(object: any): SAMLBinding {
  switch (object) {
    case 0:
    case "SAML_BINDING_UNSPECIFIED":
      return SAMLBinding.SAML_BINDING_UNSPECIFIED;
    case 1:
    case "SAML_BINDING_POST":
      return SAMLBinding.SAML_BINDING_POST;
    case 2:
    case "SAML_BINDING_REDIRECT":
      return SAMLBinding.SAML_BINDING_REDIRECT;
    case 3:
    case "SAML_BINDING_ARTIFACT":
      return SAMLBinding.SAML_BINDING_ARTIFACT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SAMLBinding.UNRECOGNIZED;
  }
}

export function sAMLBindingToJSON(object: SAMLBinding): string {
  switch (object) {
    case SAMLBinding.SAML_BINDING_UNSPECIFIED:
      return "SAML_BINDING_UNSPECIFIED";
    case SAMLBinding.SAML_BINDING_POST:
      return "SAML_BINDING_POST";
    case SAMLBinding.SAML_BINDING_REDIRECT:
      return "SAML_BINDING_REDIRECT";
    case SAMLBinding.SAML_BINDING_ARTIFACT:
      return "SAML_BINDING_ARTIFACT";
    case SAMLBinding.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum SAMLNameIDFormat {
  SAML_NAME_ID_FORMAT_UNSPECIFIED = 0,
  SAML_NAME_ID_FORMAT_EMAIL_ADDRESS = 1,
  SAML_NAME_ID_FORMAT_PERSISTENT = 2,
  SAML_NAME_ID_FORMAT_TRANSIENT = 3,
  UNRECOGNIZED = -1,
}

export function sAMLNameIDFormatFromJSON(object: any): SAMLNameIDFormat {
  switch (object) {
    case 0:
    case "SAML_NAME_ID_FORMAT_UNSPECIFIED":
      return SAMLNameIDFormat.SAML_NAME_ID_FORMAT_UNSPECIFIED;
    case 1:
    case "SAML_NAME_ID_FORMAT_EMAIL_ADDRESS":
      return SAMLNameIDFormat.SAML_NAME_ID_FORMAT_EMAIL_ADDRESS;
    case 2:
    case "SAML_NAME_ID_FORMAT_PERSISTENT":
      return SAMLNameIDFormat.SAML_NAME_ID_FORMAT_PERSISTENT;
    case 3:
    case "SAML_NAME_ID_FORMAT_TRANSIENT":
      return SAMLNameIDFormat.SAML_NAME_ID_FORMAT_TRANSIENT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SAMLNameIDFormat.UNRECOGNIZED;
  }
}

export function sAMLNameIDFormatToJSON(object: SAMLNameIDFormat): string {
  switch (object) {
    case SAMLNameIDFormat.SAML_NAME_ID_FORMAT_UNSPECIFIED:
      return "SAML_NAME_ID_FORMAT_UNSPECIFIED";
    case SAMLNameIDFormat.SAML_NAME_ID_FORMAT_EMAIL_ADDRESS:
      return "SAML_NAME_ID_FORMAT_EMAIL_ADDRESS";
    case SAMLNameIDFormat.SAML_NAME_ID_FORMAT_PERSISTENT:
      return "SAML_NAME_ID_FORMAT_PERSISTENT";
    case SAMLNameIDFormat.SAML_NAME_ID_FORMAT_TRANSIENT:
      return "SAML_NAME_ID_FORMAT_TRANSIENT";
    case SAMLNameIDFormat.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum AutoLinkingOption {
  /** AUTO_LINKING_OPTION_UNSPECIFIED - AUTO_LINKING_OPTION_UNSPECIFIED disables the auto linking prompt. */
  AUTO_LINKING_OPTION_UNSPECIFIED = 0,
  /** AUTO_LINKING_OPTION_USERNAME - AUTO_LINKING_OPTION_USERNAME will use the username of the external user to check for a corresponding ZITADEL user. */
  AUTO_LINKING_OPTION_USERNAME = 1,
  /**
   * AUTO_LINKING_OPTION_EMAIL - AUTO_LINKING_OPTION_EMAIL  will use the email of the external user to check for a corresponding ZITADEL user with the same verified email
   * Note that in case multiple users match, no prompt will be shown.
   */
  AUTO_LINKING_OPTION_EMAIL = 2,
  UNRECOGNIZED = -1,
}

export function autoLinkingOptionFromJSON(object: any): AutoLinkingOption {
  switch (object) {
    case 0:
    case "AUTO_LINKING_OPTION_UNSPECIFIED":
      return AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED;
    case 1:
    case "AUTO_LINKING_OPTION_USERNAME":
      return AutoLinkingOption.AUTO_LINKING_OPTION_USERNAME;
    case 2:
    case "AUTO_LINKING_OPTION_EMAIL":
      return AutoLinkingOption.AUTO_LINKING_OPTION_EMAIL;
    case -1:
    case "UNRECOGNIZED":
    default:
      return AutoLinkingOption.UNRECOGNIZED;
  }
}

export function autoLinkingOptionToJSON(object: AutoLinkingOption): string {
  switch (object) {
    case AutoLinkingOption.AUTO_LINKING_OPTION_UNSPECIFIED:
      return "AUTO_LINKING_OPTION_UNSPECIFIED";
    case AutoLinkingOption.AUTO_LINKING_OPTION_USERNAME:
      return "AUTO_LINKING_OPTION_USERNAME";
    case AutoLinkingOption.AUTO_LINKING_OPTION_EMAIL:
      return "AUTO_LINKING_OPTION_EMAIL";
    case AutoLinkingOption.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum AzureADTenantType {
  AZURE_AD_TENANT_TYPE_COMMON = 0,
  AZURE_AD_TENANT_TYPE_ORGANISATIONS = 1,
  AZURE_AD_TENANT_TYPE_CONSUMERS = 2,
  UNRECOGNIZED = -1,
}

export function azureADTenantTypeFromJSON(object: any): AzureADTenantType {
  switch (object) {
    case 0:
    case "AZURE_AD_TENANT_TYPE_COMMON":
      return AzureADTenantType.AZURE_AD_TENANT_TYPE_COMMON;
    case 1:
    case "AZURE_AD_TENANT_TYPE_ORGANISATIONS":
      return AzureADTenantType.AZURE_AD_TENANT_TYPE_ORGANISATIONS;
    case 2:
    case "AZURE_AD_TENANT_TYPE_CONSUMERS":
      return AzureADTenantType.AZURE_AD_TENANT_TYPE_CONSUMERS;
    case -1:
    case "UNRECOGNIZED":
    default:
      return AzureADTenantType.UNRECOGNIZED;
  }
}

export function azureADTenantTypeToJSON(object: AzureADTenantType): string {
  switch (object) {
    case AzureADTenantType.AZURE_AD_TENANT_TYPE_COMMON:
      return "AZURE_AD_TENANT_TYPE_COMMON";
    case AzureADTenantType.AZURE_AD_TENANT_TYPE_ORGANISATIONS:
      return "AZURE_AD_TENANT_TYPE_ORGANISATIONS";
    case AzureADTenantType.AZURE_AD_TENANT_TYPE_CONSUMERS:
      return "AZURE_AD_TENANT_TYPE_CONSUMERS";
    case AzureADTenantType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface IDP {
  id: string;
  details: ObjectDetails | undefined;
  state: IDPState;
  name: string;
  stylingType: IDPStylingType;
  owner: IDPOwnerType;
  oidcConfig?: OIDCConfig | undefined;
  jwtConfig?: JWTConfig | undefined;
  autoRegister: boolean;
}

export interface IDPUserLink {
  userId: string;
  idpId: string;
  idpName: string;
  providedUserId: string;
  providedUserName: string;
  idpType: IDPType;
}

export interface IDPLoginPolicyLink {
  idpId: string;
  idpName: string;
  idpType: IDPType;
}

export interface OIDCConfig {
  clientId: string;
  issuer: string;
  scopes: string[];
  displayNameMapping: OIDCMappingField;
  usernameMapping: OIDCMappingField;
}

export interface JWTConfig {
  jwtEndpoint: string;
  issuer: string;
  keysEndpoint: string;
  headerName: string;
}

export interface IDPIDQuery {
  id: string;
}

export interface IDPNameQuery {
  name: string;
  method: TextQueryMethod;
}

export interface IDPOwnerTypeQuery {
  ownerType: IDPOwnerType;
}

export interface Provider {
  id: string;
  details: ObjectDetails | undefined;
  state: IDPState;
  name: string;
  owner: IDPOwnerType;
  type: ProviderType;
  config: ProviderConfig | undefined;
}

export interface ProviderConfig {
  options: Options | undefined;
  ldap?: LDAPConfig | undefined;
  google?: GoogleConfig | undefined;
  oauth?: OAuthConfig | undefined;
  oidc?: GenericOIDCConfig | undefined;
  jwt?: JWTConfig | undefined;
  github?: GitHubConfig | undefined;
  githubEs?: GitHubEnterpriseServerConfig | undefined;
  gitlab?: GitLabConfig | undefined;
  gitlabSelfHosted?: GitLabSelfHostedConfig | undefined;
  azureAd?: AzureADConfig | undefined;
  apple?: AppleConfig | undefined;
  saml?: SAMLConfig | undefined;
}

export interface OAuthConfig {
  clientId: string;
  authorizationEndpoint: string;
  tokenEndpoint: string;
  userEndpoint: string;
  scopes: string[];
  idAttribute: string;
}

export interface GenericOIDCConfig {
  issuer: string;
  clientId: string;
  scopes: string[];
  isIdTokenMapping: boolean;
}

export interface GitHubConfig {
  clientId: string;
  scopes: string[];
}

export interface GitHubEnterpriseServerConfig {
  clientId: string;
  authorizationEndpoint: string;
  tokenEndpoint: string;
  userEndpoint: string;
  scopes: string[];
}

export interface GoogleConfig {
  clientId: string;
  scopes: string[];
}

export interface GitLabConfig {
  clientId: string;
  scopes: string[];
}

export interface GitLabSelfHostedConfig {
  issuer: string;
  clientId: string;
  scopes: string[];
}

export interface LDAPConfig {
  servers: string[];
  startTls: boolean;
  baseDn: string;
  bindDn: string;
  userBase: string;
  userObjectClasses: string[];
  userFilters: string[];
  timeout: Duration | undefined;
  attributes: LDAPAttributes | undefined;
}

export interface SAMLConfig {
  /** Metadata of the SAML identity provider. */
  metadataXml: Buffer;
  /** Binding which defines the type of communication with the identity provider. */
  binding: SAMLBinding;
  /** Boolean which defines if the authentication requests are signed. */
  withSignedRequest: boolean;
  /** `nameid-format` for the SAML Request. */
  nameIdFormat: SAMLNameIDFormat;
  /**
   * Optional name of the attribute, which will be used to map the user
   * in case the nameid-format returned is `urn:oasis:names:tc:SAML:2.0:nameid-format:transient`.
   */
  transientMappingAttributeName?: string | undefined;
}

export interface AzureADConfig {
  clientId: string;
  tenant: AzureADTenant | undefined;
  emailVerified: boolean;
  scopes: string[];
}

export interface Options {
  isLinkingAllowed: boolean;
  isCreationAllowed: boolean;
  isAutoCreation: boolean;
  isAutoUpdate: boolean;
  autoLinking: AutoLinkingOption;
}

export interface LDAPAttributes {
  idAttribute: string;
  firstNameAttribute: string;
  lastNameAttribute: string;
  displayNameAttribute: string;
  nickNameAttribute: string;
  preferredUsernameAttribute: string;
  emailAttribute: string;
  emailVerifiedAttribute: string;
  phoneAttribute: string;
  phoneVerifiedAttribute: string;
  preferredLanguageAttribute: string;
  avatarUrlAttribute: string;
  profileAttribute: string;
}

export interface AzureADTenant {
  tenantType?: AzureADTenantType | undefined;
  tenantId?: string | undefined;
}

export interface AppleConfig {
  clientId: string;
  teamId: string;
  keyId: string;
  scopes: string[];
}

function createBaseIDP(): IDP {
  return {
    id: "",
    details: undefined,
    state: 0,
    name: "",
    stylingType: 0,
    owner: 0,
    oidcConfig: undefined,
    jwtConfig: undefined,
    autoRegister: false,
  };
}

export const IDP = {
  encode(message: IDP, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.state !== 0) {
      writer.uint32(24).int32(message.state);
    }
    if (message.name !== "") {
      writer.uint32(34).string(message.name);
    }
    if (message.stylingType !== 0) {
      writer.uint32(40).int32(message.stylingType);
    }
    if (message.owner !== 0) {
      writer.uint32(48).int32(message.owner);
    }
    if (message.oidcConfig !== undefined) {
      OIDCConfig.encode(message.oidcConfig, writer.uint32(58).fork()).ldelim();
    }
    if (message.jwtConfig !== undefined) {
      JWTConfig.encode(message.jwtConfig, writer.uint32(74).fork()).ldelim();
    }
    if (message.autoRegister === true) {
      writer.uint32(64).bool(message.autoRegister);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDP {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDP();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.name = reader.string();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.stylingType = reader.int32() as any;
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.owner = reader.int32() as any;
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.oidcConfig = OIDCConfig.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.jwtConfig = JWTConfig.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag != 64) {
            break;
          }

          message.autoRegister = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDP {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      state: isSet(object.state) ? iDPStateFromJSON(object.state) : 0,
      name: isSet(object.name) ? String(object.name) : "",
      stylingType: isSet(object.stylingType) ? iDPStylingTypeFromJSON(object.stylingType) : 0,
      owner: isSet(object.owner) ? iDPOwnerTypeFromJSON(object.owner) : 0,
      oidcConfig: isSet(object.oidcConfig) ? OIDCConfig.fromJSON(object.oidcConfig) : undefined,
      jwtConfig: isSet(object.jwtConfig) ? JWTConfig.fromJSON(object.jwtConfig) : undefined,
      autoRegister: isSet(object.autoRegister) ? Boolean(object.autoRegister) : false,
    };
  },

  toJSON(message: IDP): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.state !== undefined && (obj.state = iDPStateToJSON(message.state));
    message.name !== undefined && (obj.name = message.name);
    message.stylingType !== undefined && (obj.stylingType = iDPStylingTypeToJSON(message.stylingType));
    message.owner !== undefined && (obj.owner = iDPOwnerTypeToJSON(message.owner));
    message.oidcConfig !== undefined &&
      (obj.oidcConfig = message.oidcConfig ? OIDCConfig.toJSON(message.oidcConfig) : undefined);
    message.jwtConfig !== undefined &&
      (obj.jwtConfig = message.jwtConfig ? JWTConfig.toJSON(message.jwtConfig) : undefined);
    message.autoRegister !== undefined && (obj.autoRegister = message.autoRegister);
    return obj;
  },

  create(base?: DeepPartial<IDP>): IDP {
    return IDP.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDP>): IDP {
    const message = createBaseIDP();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.state = object.state ?? 0;
    message.name = object.name ?? "";
    message.stylingType = object.stylingType ?? 0;
    message.owner = object.owner ?? 0;
    message.oidcConfig = (object.oidcConfig !== undefined && object.oidcConfig !== null)
      ? OIDCConfig.fromPartial(object.oidcConfig)
      : undefined;
    message.jwtConfig = (object.jwtConfig !== undefined && object.jwtConfig !== null)
      ? JWTConfig.fromPartial(object.jwtConfig)
      : undefined;
    message.autoRegister = object.autoRegister ?? false;
    return message;
  },
};

function createBaseIDPUserLink(): IDPUserLink {
  return { userId: "", idpId: "", idpName: "", providedUserId: "", providedUserName: "", idpType: 0 };
}

export const IDPUserLink = {
  encode(message: IDPUserLink, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.idpId !== "") {
      writer.uint32(18).string(message.idpId);
    }
    if (message.idpName !== "") {
      writer.uint32(26).string(message.idpName);
    }
    if (message.providedUserId !== "") {
      writer.uint32(34).string(message.providedUserId);
    }
    if (message.providedUserName !== "") {
      writer.uint32(42).string(message.providedUserName);
    }
    if (message.idpType !== 0) {
      writer.uint32(48).int32(message.idpType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPUserLink {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPUserLink();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idpId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.idpName = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.providedUserId = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.providedUserName = reader.string();
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.idpType = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPUserLink {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      idpName: isSet(object.idpName) ? String(object.idpName) : "",
      providedUserId: isSet(object.providedUserId) ? String(object.providedUserId) : "",
      providedUserName: isSet(object.providedUserName) ? String(object.providedUserName) : "",
      idpType: isSet(object.idpType) ? iDPTypeFromJSON(object.idpType) : 0,
    };
  },

  toJSON(message: IDPUserLink): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.idpName !== undefined && (obj.idpName = message.idpName);
    message.providedUserId !== undefined && (obj.providedUserId = message.providedUserId);
    message.providedUserName !== undefined && (obj.providedUserName = message.providedUserName);
    message.idpType !== undefined && (obj.idpType = iDPTypeToJSON(message.idpType));
    return obj;
  },

  create(base?: DeepPartial<IDPUserLink>): IDPUserLink {
    return IDPUserLink.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPUserLink>): IDPUserLink {
    const message = createBaseIDPUserLink();
    message.userId = object.userId ?? "";
    message.idpId = object.idpId ?? "";
    message.idpName = object.idpName ?? "";
    message.providedUserId = object.providedUserId ?? "";
    message.providedUserName = object.providedUserName ?? "";
    message.idpType = object.idpType ?? 0;
    return message;
  },
};

function createBaseIDPLoginPolicyLink(): IDPLoginPolicyLink {
  return { idpId: "", idpName: "", idpType: 0 };
}

export const IDPLoginPolicyLink = {
  encode(message: IDPLoginPolicyLink, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.idpName !== "") {
      writer.uint32(18).string(message.idpName);
    }
    if (message.idpType !== 0) {
      writer.uint32(24).int32(message.idpType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPLoginPolicyLink {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPLoginPolicyLink();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idpId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idpName = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.idpType = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPLoginPolicyLink {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      idpName: isSet(object.idpName) ? String(object.idpName) : "",
      idpType: isSet(object.idpType) ? iDPTypeFromJSON(object.idpType) : 0,
    };
  },

  toJSON(message: IDPLoginPolicyLink): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.idpName !== undefined && (obj.idpName = message.idpName);
    message.idpType !== undefined && (obj.idpType = iDPTypeToJSON(message.idpType));
    return obj;
  },

  create(base?: DeepPartial<IDPLoginPolicyLink>): IDPLoginPolicyLink {
    return IDPLoginPolicyLink.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPLoginPolicyLink>): IDPLoginPolicyLink {
    const message = createBaseIDPLoginPolicyLink();
    message.idpId = object.idpId ?? "";
    message.idpName = object.idpName ?? "";
    message.idpType = object.idpType ?? 0;
    return message;
  },
};

function createBaseOIDCConfig(): OIDCConfig {
  return { clientId: "", issuer: "", scopes: [], displayNameMapping: 0, usernameMapping: 0 };
}

export const OIDCConfig = {
  encode(message: OIDCConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    if (message.issuer !== "") {
      writer.uint32(18).string(message.issuer);
    }
    for (const v of message.scopes) {
      writer.uint32(26).string(v!);
    }
    if (message.displayNameMapping !== 0) {
      writer.uint32(32).int32(message.displayNameMapping);
    }
    if (message.usernameMapping !== 0) {
      writer.uint32(40).int32(message.usernameMapping);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OIDCConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOIDCConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.issuer = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.displayNameMapping = reader.int32() as any;
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.usernameMapping = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OIDCConfig {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      displayNameMapping: isSet(object.displayNameMapping) ? oIDCMappingFieldFromJSON(object.displayNameMapping) : 0,
      usernameMapping: isSet(object.usernameMapping) ? oIDCMappingFieldFromJSON(object.usernameMapping) : 0,
    };
  },

  toJSON(message: OIDCConfig): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.issuer !== undefined && (obj.issuer = message.issuer);
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

  create(base?: DeepPartial<OIDCConfig>): OIDCConfig {
    return OIDCConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OIDCConfig>): OIDCConfig {
    const message = createBaseOIDCConfig();
    message.clientId = object.clientId ?? "";
    message.issuer = object.issuer ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.displayNameMapping = object.displayNameMapping ?? 0;
    message.usernameMapping = object.usernameMapping ?? 0;
    return message;
  },
};

function createBaseJWTConfig(): JWTConfig {
  return { jwtEndpoint: "", issuer: "", keysEndpoint: "", headerName: "" };
}

export const JWTConfig = {
  encode(message: JWTConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.jwtEndpoint !== "") {
      writer.uint32(10).string(message.jwtEndpoint);
    }
    if (message.issuer !== "") {
      writer.uint32(18).string(message.issuer);
    }
    if (message.keysEndpoint !== "") {
      writer.uint32(26).string(message.keysEndpoint);
    }
    if (message.headerName !== "") {
      writer.uint32(34).string(message.headerName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): JWTConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseJWTConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.jwtEndpoint = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.issuer = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.keysEndpoint = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.headerName = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): JWTConfig {
    return {
      jwtEndpoint: isSet(object.jwtEndpoint) ? String(object.jwtEndpoint) : "",
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      keysEndpoint: isSet(object.keysEndpoint) ? String(object.keysEndpoint) : "",
      headerName: isSet(object.headerName) ? String(object.headerName) : "",
    };
  },

  toJSON(message: JWTConfig): unknown {
    const obj: any = {};
    message.jwtEndpoint !== undefined && (obj.jwtEndpoint = message.jwtEndpoint);
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.keysEndpoint !== undefined && (obj.keysEndpoint = message.keysEndpoint);
    message.headerName !== undefined && (obj.headerName = message.headerName);
    return obj;
  },

  create(base?: DeepPartial<JWTConfig>): JWTConfig {
    return JWTConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<JWTConfig>): JWTConfig {
    const message = createBaseJWTConfig();
    message.jwtEndpoint = object.jwtEndpoint ?? "";
    message.issuer = object.issuer ?? "";
    message.keysEndpoint = object.keysEndpoint ?? "";
    message.headerName = object.headerName ?? "";
    return message;
  },
};

function createBaseIDPIDQuery(): IDPIDQuery {
  return { id: "" };
}

export const IDPIDQuery = {
  encode(message: IDPIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPIDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPIDQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPIDQuery {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: IDPIDQuery): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<IDPIDQuery>): IDPIDQuery {
    return IDPIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPIDQuery>): IDPIDQuery {
    const message = createBaseIDPIDQuery();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseIDPNameQuery(): IDPNameQuery {
  return { name: "", method: 0 };
}

export const IDPNameQuery = {
  encode(message: IDPNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPNameQuery {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: IDPNameQuery): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<IDPNameQuery>): IDPNameQuery {
    return IDPNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPNameQuery>): IDPNameQuery {
    const message = createBaseIDPNameQuery();
    message.name = object.name ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseIDPOwnerTypeQuery(): IDPOwnerTypeQuery {
  return { ownerType: 0 };
}

export const IDPOwnerTypeQuery = {
  encode(message: IDPOwnerTypeQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ownerType !== 0) {
      writer.uint32(8).int32(message.ownerType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPOwnerTypeQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPOwnerTypeQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.ownerType = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPOwnerTypeQuery {
    return { ownerType: isSet(object.ownerType) ? iDPOwnerTypeFromJSON(object.ownerType) : 0 };
  },

  toJSON(message: IDPOwnerTypeQuery): unknown {
    const obj: any = {};
    message.ownerType !== undefined && (obj.ownerType = iDPOwnerTypeToJSON(message.ownerType));
    return obj;
  },

  create(base?: DeepPartial<IDPOwnerTypeQuery>): IDPOwnerTypeQuery {
    return IDPOwnerTypeQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPOwnerTypeQuery>): IDPOwnerTypeQuery {
    const message = createBaseIDPOwnerTypeQuery();
    message.ownerType = object.ownerType ?? 0;
    return message;
  },
};

function createBaseProvider(): Provider {
  return { id: "", details: undefined, state: 0, name: "", owner: 0, type: 0, config: undefined };
}

export const Provider = {
  encode(message: Provider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.state !== 0) {
      writer.uint32(24).int32(message.state);
    }
    if (message.name !== "") {
      writer.uint32(34).string(message.name);
    }
    if (message.owner !== 0) {
      writer.uint32(40).int32(message.owner);
    }
    if (message.type !== 0) {
      writer.uint32(48).int32(message.type);
    }
    if (message.config !== undefined) {
      ProviderConfig.encode(message.config, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Provider {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.name = reader.string();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.owner = reader.int32() as any;
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.config = ProviderConfig.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Provider {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      state: isSet(object.state) ? iDPStateFromJSON(object.state) : 0,
      name: isSet(object.name) ? String(object.name) : "",
      owner: isSet(object.owner) ? iDPOwnerTypeFromJSON(object.owner) : 0,
      type: isSet(object.type) ? providerTypeFromJSON(object.type) : 0,
      config: isSet(object.config) ? ProviderConfig.fromJSON(object.config) : undefined,
    };
  },

  toJSON(message: Provider): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.state !== undefined && (obj.state = iDPStateToJSON(message.state));
    message.name !== undefined && (obj.name = message.name);
    message.owner !== undefined && (obj.owner = iDPOwnerTypeToJSON(message.owner));
    message.type !== undefined && (obj.type = providerTypeToJSON(message.type));
    message.config !== undefined && (obj.config = message.config ? ProviderConfig.toJSON(message.config) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Provider>): Provider {
    return Provider.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Provider>): Provider {
    const message = createBaseProvider();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.state = object.state ?? 0;
    message.name = object.name ?? "";
    message.owner = object.owner ?? 0;
    message.type = object.type ?? 0;
    message.config = (object.config !== undefined && object.config !== null)
      ? ProviderConfig.fromPartial(object.config)
      : undefined;
    return message;
  },
};

function createBaseProviderConfig(): ProviderConfig {
  return {
    options: undefined,
    ldap: undefined,
    google: undefined,
    oauth: undefined,
    oidc: undefined,
    jwt: undefined,
    github: undefined,
    githubEs: undefined,
    gitlab: undefined,
    gitlabSelfHosted: undefined,
    azureAd: undefined,
    apple: undefined,
    saml: undefined,
  };
}

export const ProviderConfig = {
  encode(message: ProviderConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.options !== undefined) {
      Options.encode(message.options, writer.uint32(10).fork()).ldelim();
    }
    if (message.ldap !== undefined) {
      LDAPConfig.encode(message.ldap, writer.uint32(18).fork()).ldelim();
    }
    if (message.google !== undefined) {
      GoogleConfig.encode(message.google, writer.uint32(26).fork()).ldelim();
    }
    if (message.oauth !== undefined) {
      OAuthConfig.encode(message.oauth, writer.uint32(34).fork()).ldelim();
    }
    if (message.oidc !== undefined) {
      GenericOIDCConfig.encode(message.oidc, writer.uint32(42).fork()).ldelim();
    }
    if (message.jwt !== undefined) {
      JWTConfig.encode(message.jwt, writer.uint32(50).fork()).ldelim();
    }
    if (message.github !== undefined) {
      GitHubConfig.encode(message.github, writer.uint32(58).fork()).ldelim();
    }
    if (message.githubEs !== undefined) {
      GitHubEnterpriseServerConfig.encode(message.githubEs, writer.uint32(66).fork()).ldelim();
    }
    if (message.gitlab !== undefined) {
      GitLabConfig.encode(message.gitlab, writer.uint32(74).fork()).ldelim();
    }
    if (message.gitlabSelfHosted !== undefined) {
      GitLabSelfHostedConfig.encode(message.gitlabSelfHosted, writer.uint32(82).fork()).ldelim();
    }
    if (message.azureAd !== undefined) {
      AzureADConfig.encode(message.azureAd, writer.uint32(90).fork()).ldelim();
    }
    if (message.apple !== undefined) {
      AppleConfig.encode(message.apple, writer.uint32(98).fork()).ldelim();
    }
    if (message.saml !== undefined) {
      SAMLConfig.encode(message.saml, writer.uint32(106).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProviderConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProviderConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.options = Options.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.ldap = LDAPConfig.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.google = GoogleConfig.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.oauth = OAuthConfig.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.oidc = GenericOIDCConfig.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.jwt = JWTConfig.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.github = GitHubConfig.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.githubEs = GitHubEnterpriseServerConfig.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.gitlab = GitLabConfig.decode(reader, reader.uint32());
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          message.gitlabSelfHosted = GitLabSelfHostedConfig.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag != 90) {
            break;
          }

          message.azureAd = AzureADConfig.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag != 98) {
            break;
          }

          message.apple = AppleConfig.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag != 106) {
            break;
          }

          message.saml = SAMLConfig.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ProviderConfig {
    return {
      options: isSet(object.options) ? Options.fromJSON(object.options) : undefined,
      ldap: isSet(object.ldap) ? LDAPConfig.fromJSON(object.ldap) : undefined,
      google: isSet(object.google) ? GoogleConfig.fromJSON(object.google) : undefined,
      oauth: isSet(object.oauth) ? OAuthConfig.fromJSON(object.oauth) : undefined,
      oidc: isSet(object.oidc) ? GenericOIDCConfig.fromJSON(object.oidc) : undefined,
      jwt: isSet(object.jwt) ? JWTConfig.fromJSON(object.jwt) : undefined,
      github: isSet(object.github) ? GitHubConfig.fromJSON(object.github) : undefined,
      githubEs: isSet(object.githubEs) ? GitHubEnterpriseServerConfig.fromJSON(object.githubEs) : undefined,
      gitlab: isSet(object.gitlab) ? GitLabConfig.fromJSON(object.gitlab) : undefined,
      gitlabSelfHosted: isSet(object.gitlabSelfHosted)
        ? GitLabSelfHostedConfig.fromJSON(object.gitlabSelfHosted)
        : undefined,
      azureAd: isSet(object.azureAd) ? AzureADConfig.fromJSON(object.azureAd) : undefined,
      apple: isSet(object.apple) ? AppleConfig.fromJSON(object.apple) : undefined,
      saml: isSet(object.saml) ? SAMLConfig.fromJSON(object.saml) : undefined,
    };
  },

  toJSON(message: ProviderConfig): unknown {
    const obj: any = {};
    message.options !== undefined && (obj.options = message.options ? Options.toJSON(message.options) : undefined);
    message.ldap !== undefined && (obj.ldap = message.ldap ? LDAPConfig.toJSON(message.ldap) : undefined);
    message.google !== undefined && (obj.google = message.google ? GoogleConfig.toJSON(message.google) : undefined);
    message.oauth !== undefined && (obj.oauth = message.oauth ? OAuthConfig.toJSON(message.oauth) : undefined);
    message.oidc !== undefined && (obj.oidc = message.oidc ? GenericOIDCConfig.toJSON(message.oidc) : undefined);
    message.jwt !== undefined && (obj.jwt = message.jwt ? JWTConfig.toJSON(message.jwt) : undefined);
    message.github !== undefined && (obj.github = message.github ? GitHubConfig.toJSON(message.github) : undefined);
    message.githubEs !== undefined &&
      (obj.githubEs = message.githubEs ? GitHubEnterpriseServerConfig.toJSON(message.githubEs) : undefined);
    message.gitlab !== undefined && (obj.gitlab = message.gitlab ? GitLabConfig.toJSON(message.gitlab) : undefined);
    message.gitlabSelfHosted !== undefined && (obj.gitlabSelfHosted = message.gitlabSelfHosted
      ? GitLabSelfHostedConfig.toJSON(message.gitlabSelfHosted)
      : undefined);
    message.azureAd !== undefined &&
      (obj.azureAd = message.azureAd ? AzureADConfig.toJSON(message.azureAd) : undefined);
    message.apple !== undefined && (obj.apple = message.apple ? AppleConfig.toJSON(message.apple) : undefined);
    message.saml !== undefined && (obj.saml = message.saml ? SAMLConfig.toJSON(message.saml) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ProviderConfig>): ProviderConfig {
    return ProviderConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProviderConfig>): ProviderConfig {
    const message = createBaseProviderConfig();
    message.options = (object.options !== undefined && object.options !== null)
      ? Options.fromPartial(object.options)
      : undefined;
    message.ldap = (object.ldap !== undefined && object.ldap !== null)
      ? LDAPConfig.fromPartial(object.ldap)
      : undefined;
    message.google = (object.google !== undefined && object.google !== null)
      ? GoogleConfig.fromPartial(object.google)
      : undefined;
    message.oauth = (object.oauth !== undefined && object.oauth !== null)
      ? OAuthConfig.fromPartial(object.oauth)
      : undefined;
    message.oidc = (object.oidc !== undefined && object.oidc !== null)
      ? GenericOIDCConfig.fromPartial(object.oidc)
      : undefined;
    message.jwt = (object.jwt !== undefined && object.jwt !== null) ? JWTConfig.fromPartial(object.jwt) : undefined;
    message.github = (object.github !== undefined && object.github !== null)
      ? GitHubConfig.fromPartial(object.github)
      : undefined;
    message.githubEs = (object.githubEs !== undefined && object.githubEs !== null)
      ? GitHubEnterpriseServerConfig.fromPartial(object.githubEs)
      : undefined;
    message.gitlab = (object.gitlab !== undefined && object.gitlab !== null)
      ? GitLabConfig.fromPartial(object.gitlab)
      : undefined;
    message.gitlabSelfHosted = (object.gitlabSelfHosted !== undefined && object.gitlabSelfHosted !== null)
      ? GitLabSelfHostedConfig.fromPartial(object.gitlabSelfHosted)
      : undefined;
    message.azureAd = (object.azureAd !== undefined && object.azureAd !== null)
      ? AzureADConfig.fromPartial(object.azureAd)
      : undefined;
    message.apple = (object.apple !== undefined && object.apple !== null)
      ? AppleConfig.fromPartial(object.apple)
      : undefined;
    message.saml = (object.saml !== undefined && object.saml !== null)
      ? SAMLConfig.fromPartial(object.saml)
      : undefined;
    return message;
  },
};

function createBaseOAuthConfig(): OAuthConfig {
  return { clientId: "", authorizationEndpoint: "", tokenEndpoint: "", userEndpoint: "", scopes: [], idAttribute: "" };
}

export const OAuthConfig = {
  encode(message: OAuthConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    if (message.authorizationEndpoint !== "") {
      writer.uint32(18).string(message.authorizationEndpoint);
    }
    if (message.tokenEndpoint !== "") {
      writer.uint32(26).string(message.tokenEndpoint);
    }
    if (message.userEndpoint !== "") {
      writer.uint32(34).string(message.userEndpoint);
    }
    for (const v of message.scopes) {
      writer.uint32(42).string(v!);
    }
    if (message.idAttribute !== "") {
      writer.uint32(50).string(message.idAttribute);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OAuthConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOAuthConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.authorizationEndpoint = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.tokenEndpoint = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.userEndpoint = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.idAttribute = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OAuthConfig {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      authorizationEndpoint: isSet(object.authorizationEndpoint) ? String(object.authorizationEndpoint) : "",
      tokenEndpoint: isSet(object.tokenEndpoint) ? String(object.tokenEndpoint) : "",
      userEndpoint: isSet(object.userEndpoint) ? String(object.userEndpoint) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      idAttribute: isSet(object.idAttribute) ? String(object.idAttribute) : "",
    };
  },

  toJSON(message: OAuthConfig): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.authorizationEndpoint !== undefined && (obj.authorizationEndpoint = message.authorizationEndpoint);
    message.tokenEndpoint !== undefined && (obj.tokenEndpoint = message.tokenEndpoint);
    message.userEndpoint !== undefined && (obj.userEndpoint = message.userEndpoint);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.idAttribute !== undefined && (obj.idAttribute = message.idAttribute);
    return obj;
  },

  create(base?: DeepPartial<OAuthConfig>): OAuthConfig {
    return OAuthConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OAuthConfig>): OAuthConfig {
    const message = createBaseOAuthConfig();
    message.clientId = object.clientId ?? "";
    message.authorizationEndpoint = object.authorizationEndpoint ?? "";
    message.tokenEndpoint = object.tokenEndpoint ?? "";
    message.userEndpoint = object.userEndpoint ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.idAttribute = object.idAttribute ?? "";
    return message;
  },
};

function createBaseGenericOIDCConfig(): GenericOIDCConfig {
  return { issuer: "", clientId: "", scopes: [], isIdTokenMapping: false };
}

export const GenericOIDCConfig = {
  encode(message: GenericOIDCConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.clientId !== "") {
      writer.uint32(18).string(message.clientId);
    }
    for (const v of message.scopes) {
      writer.uint32(26).string(v!);
    }
    if (message.isIdTokenMapping === true) {
      writer.uint32(32).bool(message.isIdTokenMapping);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenericOIDCConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenericOIDCConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.issuer = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.isIdTokenMapping = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GenericOIDCConfig {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
      isIdTokenMapping: isSet(object.isIdTokenMapping) ? Boolean(object.isIdTokenMapping) : false,
    };
  },

  toJSON(message: GenericOIDCConfig): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    message.isIdTokenMapping !== undefined && (obj.isIdTokenMapping = message.isIdTokenMapping);
    return obj;
  },

  create(base?: DeepPartial<GenericOIDCConfig>): GenericOIDCConfig {
    return GenericOIDCConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GenericOIDCConfig>): GenericOIDCConfig {
    const message = createBaseGenericOIDCConfig();
    message.issuer = object.issuer ?? "";
    message.clientId = object.clientId ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.isIdTokenMapping = object.isIdTokenMapping ?? false;
    return message;
  },
};

function createBaseGitHubConfig(): GitHubConfig {
  return { clientId: "", scopes: [] };
}

export const GitHubConfig = {
  encode(message: GitHubConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    for (const v of message.scopes) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GitHubConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGitHubConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GitHubConfig {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: GitHubConfig): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GitHubConfig>): GitHubConfig {
    return GitHubConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GitHubConfig>): GitHubConfig {
    const message = createBaseGitHubConfig();
    message.clientId = object.clientId ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

function createBaseGitHubEnterpriseServerConfig(): GitHubEnterpriseServerConfig {
  return { clientId: "", authorizationEndpoint: "", tokenEndpoint: "", userEndpoint: "", scopes: [] };
}

export const GitHubEnterpriseServerConfig = {
  encode(message: GitHubEnterpriseServerConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    if (message.authorizationEndpoint !== "") {
      writer.uint32(18).string(message.authorizationEndpoint);
    }
    if (message.tokenEndpoint !== "") {
      writer.uint32(26).string(message.tokenEndpoint);
    }
    if (message.userEndpoint !== "") {
      writer.uint32(34).string(message.userEndpoint);
    }
    for (const v of message.scopes) {
      writer.uint32(42).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GitHubEnterpriseServerConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGitHubEnterpriseServerConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.authorizationEndpoint = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.tokenEndpoint = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.userEndpoint = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GitHubEnterpriseServerConfig {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      authorizationEndpoint: isSet(object.authorizationEndpoint) ? String(object.authorizationEndpoint) : "",
      tokenEndpoint: isSet(object.tokenEndpoint) ? String(object.tokenEndpoint) : "",
      userEndpoint: isSet(object.userEndpoint) ? String(object.userEndpoint) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: GitHubEnterpriseServerConfig): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.authorizationEndpoint !== undefined && (obj.authorizationEndpoint = message.authorizationEndpoint);
    message.tokenEndpoint !== undefined && (obj.tokenEndpoint = message.tokenEndpoint);
    message.userEndpoint !== undefined && (obj.userEndpoint = message.userEndpoint);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GitHubEnterpriseServerConfig>): GitHubEnterpriseServerConfig {
    return GitHubEnterpriseServerConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GitHubEnterpriseServerConfig>): GitHubEnterpriseServerConfig {
    const message = createBaseGitHubEnterpriseServerConfig();
    message.clientId = object.clientId ?? "";
    message.authorizationEndpoint = object.authorizationEndpoint ?? "";
    message.tokenEndpoint = object.tokenEndpoint ?? "";
    message.userEndpoint = object.userEndpoint ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

function createBaseGoogleConfig(): GoogleConfig {
  return { clientId: "", scopes: [] };
}

export const GoogleConfig = {
  encode(message: GoogleConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    for (const v of message.scopes) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GoogleConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGoogleConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GoogleConfig {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: GoogleConfig): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GoogleConfig>): GoogleConfig {
    return GoogleConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GoogleConfig>): GoogleConfig {
    const message = createBaseGoogleConfig();
    message.clientId = object.clientId ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

function createBaseGitLabConfig(): GitLabConfig {
  return { clientId: "", scopes: [] };
}

export const GitLabConfig = {
  encode(message: GitLabConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    for (const v of message.scopes) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GitLabConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGitLabConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GitLabConfig {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: GitLabConfig): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GitLabConfig>): GitLabConfig {
    return GitLabConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GitLabConfig>): GitLabConfig {
    const message = createBaseGitLabConfig();
    message.clientId = object.clientId ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

function createBaseGitLabSelfHostedConfig(): GitLabSelfHostedConfig {
  return { issuer: "", clientId: "", scopes: [] };
}

export const GitLabSelfHostedConfig = {
  encode(message: GitLabSelfHostedConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.issuer !== "") {
      writer.uint32(10).string(message.issuer);
    }
    if (message.clientId !== "") {
      writer.uint32(18).string(message.clientId);
    }
    for (const v of message.scopes) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GitLabSelfHostedConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGitLabSelfHostedConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.issuer = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GitLabSelfHostedConfig {
    return {
      issuer: isSet(object.issuer) ? String(object.issuer) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: GitLabSelfHostedConfig): unknown {
    const obj: any = {};
    message.issuer !== undefined && (obj.issuer = message.issuer);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GitLabSelfHostedConfig>): GitLabSelfHostedConfig {
    return GitLabSelfHostedConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GitLabSelfHostedConfig>): GitLabSelfHostedConfig {
    const message = createBaseGitLabSelfHostedConfig();
    message.issuer = object.issuer ?? "";
    message.clientId = object.clientId ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

function createBaseLDAPConfig(): LDAPConfig {
  return {
    servers: [],
    startTls: false,
    baseDn: "",
    bindDn: "",
    userBase: "",
    userObjectClasses: [],
    userFilters: [],
    timeout: undefined,
    attributes: undefined,
  };
}

export const LDAPConfig = {
  encode(message: LDAPConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.servers) {
      writer.uint32(10).string(v!);
    }
    if (message.startTls === true) {
      writer.uint32(16).bool(message.startTls);
    }
    if (message.baseDn !== "") {
      writer.uint32(26).string(message.baseDn);
    }
    if (message.bindDn !== "") {
      writer.uint32(34).string(message.bindDn);
    }
    if (message.userBase !== "") {
      writer.uint32(42).string(message.userBase);
    }
    for (const v of message.userObjectClasses) {
      writer.uint32(50).string(v!);
    }
    for (const v of message.userFilters) {
      writer.uint32(58).string(v!);
    }
    if (message.timeout !== undefined) {
      Duration.encode(message.timeout, writer.uint32(66).fork()).ldelim();
    }
    if (message.attributes !== undefined) {
      LDAPAttributes.encode(message.attributes, writer.uint32(74).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LDAPConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLDAPConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.servers.push(reader.string());
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.startTls = reader.bool();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.baseDn = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.bindDn = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.userBase = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.userObjectClasses.push(reader.string());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.userFilters.push(reader.string());
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.timeout = Duration.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.attributes = LDAPAttributes.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LDAPConfig {
    return {
      servers: Array.isArray(object?.servers) ? object.servers.map((e: any) => String(e)) : [],
      startTls: isSet(object.startTls) ? Boolean(object.startTls) : false,
      baseDn: isSet(object.baseDn) ? String(object.baseDn) : "",
      bindDn: isSet(object.bindDn) ? String(object.bindDn) : "",
      userBase: isSet(object.userBase) ? String(object.userBase) : "",
      userObjectClasses: Array.isArray(object?.userObjectClasses)
        ? object.userObjectClasses.map((e: any) => String(e))
        : [],
      userFilters: Array.isArray(object?.userFilters) ? object.userFilters.map((e: any) => String(e)) : [],
      timeout: isSet(object.timeout) ? Duration.fromJSON(object.timeout) : undefined,
      attributes: isSet(object.attributes) ? LDAPAttributes.fromJSON(object.attributes) : undefined,
    };
  },

  toJSON(message: LDAPConfig): unknown {
    const obj: any = {};
    if (message.servers) {
      obj.servers = message.servers.map((e) => e);
    } else {
      obj.servers = [];
    }
    message.startTls !== undefined && (obj.startTls = message.startTls);
    message.baseDn !== undefined && (obj.baseDn = message.baseDn);
    message.bindDn !== undefined && (obj.bindDn = message.bindDn);
    message.userBase !== undefined && (obj.userBase = message.userBase);
    if (message.userObjectClasses) {
      obj.userObjectClasses = message.userObjectClasses.map((e) => e);
    } else {
      obj.userObjectClasses = [];
    }
    if (message.userFilters) {
      obj.userFilters = message.userFilters.map((e) => e);
    } else {
      obj.userFilters = [];
    }
    message.timeout !== undefined && (obj.timeout = message.timeout ? Duration.toJSON(message.timeout) : undefined);
    message.attributes !== undefined &&
      (obj.attributes = message.attributes ? LDAPAttributes.toJSON(message.attributes) : undefined);
    return obj;
  },

  create(base?: DeepPartial<LDAPConfig>): LDAPConfig {
    return LDAPConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LDAPConfig>): LDAPConfig {
    const message = createBaseLDAPConfig();
    message.servers = object.servers?.map((e) => e) || [];
    message.startTls = object.startTls ?? false;
    message.baseDn = object.baseDn ?? "";
    message.bindDn = object.bindDn ?? "";
    message.userBase = object.userBase ?? "";
    message.userObjectClasses = object.userObjectClasses?.map((e) => e) || [];
    message.userFilters = object.userFilters?.map((e) => e) || [];
    message.timeout = (object.timeout !== undefined && object.timeout !== null)
      ? Duration.fromPartial(object.timeout)
      : undefined;
    message.attributes = (object.attributes !== undefined && object.attributes !== null)
      ? LDAPAttributes.fromPartial(object.attributes)
      : undefined;
    return message;
  },
};

function createBaseSAMLConfig(): SAMLConfig {
  return {
    metadataXml: Buffer.alloc(0),
    binding: 0,
    withSignedRequest: false,
    nameIdFormat: 0,
    transientMappingAttributeName: undefined,
  };
}

export const SAMLConfig = {
  encode(message: SAMLConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.metadataXml.length !== 0) {
      writer.uint32(10).bytes(message.metadataXml);
    }
    if (message.binding !== 0) {
      writer.uint32(16).int32(message.binding);
    }
    if (message.withSignedRequest === true) {
      writer.uint32(24).bool(message.withSignedRequest);
    }
    if (message.nameIdFormat !== 0) {
      writer.uint32(32).int32(message.nameIdFormat);
    }
    if (message.transientMappingAttributeName !== undefined) {
      writer.uint32(42).string(message.transientMappingAttributeName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SAMLConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSAMLConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.metadataXml = reader.bytes() as Buffer;
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.binding = reader.int32() as any;
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.withSignedRequest = reader.bool();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.nameIdFormat = reader.int32() as any;
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.transientMappingAttributeName = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SAMLConfig {
    return {
      metadataXml: isSet(object.metadataXml) ? Buffer.from(bytesFromBase64(object.metadataXml)) : Buffer.alloc(0),
      binding: isSet(object.binding) ? sAMLBindingFromJSON(object.binding) : 0,
      withSignedRequest: isSet(object.withSignedRequest) ? Boolean(object.withSignedRequest) : false,
      nameIdFormat: isSet(object.nameIdFormat) ? sAMLNameIDFormatFromJSON(object.nameIdFormat) : 0,
      transientMappingAttributeName: isSet(object.transientMappingAttributeName)
        ? String(object.transientMappingAttributeName)
        : undefined,
    };
  },

  toJSON(message: SAMLConfig): unknown {
    const obj: any = {};
    message.metadataXml !== undefined &&
      (obj.metadataXml = base64FromBytes(message.metadataXml !== undefined ? message.metadataXml : Buffer.alloc(0)));
    message.binding !== undefined && (obj.binding = sAMLBindingToJSON(message.binding));
    message.withSignedRequest !== undefined && (obj.withSignedRequest = message.withSignedRequest);
    message.nameIdFormat !== undefined && (obj.nameIdFormat = sAMLNameIDFormatToJSON(message.nameIdFormat));
    message.transientMappingAttributeName !== undefined &&
      (obj.transientMappingAttributeName = message.transientMappingAttributeName);
    return obj;
  },

  create(base?: DeepPartial<SAMLConfig>): SAMLConfig {
    return SAMLConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SAMLConfig>): SAMLConfig {
    const message = createBaseSAMLConfig();
    message.metadataXml = object.metadataXml ?? Buffer.alloc(0);
    message.binding = object.binding ?? 0;
    message.withSignedRequest = object.withSignedRequest ?? false;
    message.nameIdFormat = object.nameIdFormat ?? 0;
    message.transientMappingAttributeName = object.transientMappingAttributeName ?? undefined;
    return message;
  },
};

function createBaseAzureADConfig(): AzureADConfig {
  return { clientId: "", tenant: undefined, emailVerified: false, scopes: [] };
}

export const AzureADConfig = {
  encode(message: AzureADConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    if (message.tenant !== undefined) {
      AzureADTenant.encode(message.tenant, writer.uint32(18).fork()).ldelim();
    }
    if (message.emailVerified === true) {
      writer.uint32(24).bool(message.emailVerified);
    }
    for (const v of message.scopes) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AzureADConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAzureADConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.tenant = AzureADTenant.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.emailVerified = reader.bool();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AzureADConfig {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      tenant: isSet(object.tenant) ? AzureADTenant.fromJSON(object.tenant) : undefined,
      emailVerified: isSet(object.emailVerified) ? Boolean(object.emailVerified) : false,
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: AzureADConfig): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.tenant !== undefined && (obj.tenant = message.tenant ? AzureADTenant.toJSON(message.tenant) : undefined);
    message.emailVerified !== undefined && (obj.emailVerified = message.emailVerified);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<AzureADConfig>): AzureADConfig {
    return AzureADConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AzureADConfig>): AzureADConfig {
    const message = createBaseAzureADConfig();
    message.clientId = object.clientId ?? "";
    message.tenant = (object.tenant !== undefined && object.tenant !== null)
      ? AzureADTenant.fromPartial(object.tenant)
      : undefined;
    message.emailVerified = object.emailVerified ?? false;
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

function createBaseOptions(): Options {
  return {
    isLinkingAllowed: false,
    isCreationAllowed: false,
    isAutoCreation: false,
    isAutoUpdate: false,
    autoLinking: 0,
  };
}

export const Options = {
  encode(message: Options, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.isLinkingAllowed === true) {
      writer.uint32(8).bool(message.isLinkingAllowed);
    }
    if (message.isCreationAllowed === true) {
      writer.uint32(16).bool(message.isCreationAllowed);
    }
    if (message.isAutoCreation === true) {
      writer.uint32(24).bool(message.isAutoCreation);
    }
    if (message.isAutoUpdate === true) {
      writer.uint32(32).bool(message.isAutoUpdate);
    }
    if (message.autoLinking !== 0) {
      writer.uint32(40).int32(message.autoLinking);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Options {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOptions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.isLinkingAllowed = reader.bool();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.isCreationAllowed = reader.bool();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.isAutoCreation = reader.bool();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.isAutoUpdate = reader.bool();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.autoLinking = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Options {
    return {
      isLinkingAllowed: isSet(object.isLinkingAllowed) ? Boolean(object.isLinkingAllowed) : false,
      isCreationAllowed: isSet(object.isCreationAllowed) ? Boolean(object.isCreationAllowed) : false,
      isAutoCreation: isSet(object.isAutoCreation) ? Boolean(object.isAutoCreation) : false,
      isAutoUpdate: isSet(object.isAutoUpdate) ? Boolean(object.isAutoUpdate) : false,
      autoLinking: isSet(object.autoLinking) ? autoLinkingOptionFromJSON(object.autoLinking) : 0,
    };
  },

  toJSON(message: Options): unknown {
    const obj: any = {};
    message.isLinkingAllowed !== undefined && (obj.isLinkingAllowed = message.isLinkingAllowed);
    message.isCreationAllowed !== undefined && (obj.isCreationAllowed = message.isCreationAllowed);
    message.isAutoCreation !== undefined && (obj.isAutoCreation = message.isAutoCreation);
    message.isAutoUpdate !== undefined && (obj.isAutoUpdate = message.isAutoUpdate);
    message.autoLinking !== undefined && (obj.autoLinking = autoLinkingOptionToJSON(message.autoLinking));
    return obj;
  },

  create(base?: DeepPartial<Options>): Options {
    return Options.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Options>): Options {
    const message = createBaseOptions();
    message.isLinkingAllowed = object.isLinkingAllowed ?? false;
    message.isCreationAllowed = object.isCreationAllowed ?? false;
    message.isAutoCreation = object.isAutoCreation ?? false;
    message.isAutoUpdate = object.isAutoUpdate ?? false;
    message.autoLinking = object.autoLinking ?? 0;
    return message;
  },
};

function createBaseLDAPAttributes(): LDAPAttributes {
  return {
    idAttribute: "",
    firstNameAttribute: "",
    lastNameAttribute: "",
    displayNameAttribute: "",
    nickNameAttribute: "",
    preferredUsernameAttribute: "",
    emailAttribute: "",
    emailVerifiedAttribute: "",
    phoneAttribute: "",
    phoneVerifiedAttribute: "",
    preferredLanguageAttribute: "",
    avatarUrlAttribute: "",
    profileAttribute: "",
  };
}

export const LDAPAttributes = {
  encode(message: LDAPAttributes, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idAttribute !== "") {
      writer.uint32(10).string(message.idAttribute);
    }
    if (message.firstNameAttribute !== "") {
      writer.uint32(18).string(message.firstNameAttribute);
    }
    if (message.lastNameAttribute !== "") {
      writer.uint32(26).string(message.lastNameAttribute);
    }
    if (message.displayNameAttribute !== "") {
      writer.uint32(34).string(message.displayNameAttribute);
    }
    if (message.nickNameAttribute !== "") {
      writer.uint32(42).string(message.nickNameAttribute);
    }
    if (message.preferredUsernameAttribute !== "") {
      writer.uint32(50).string(message.preferredUsernameAttribute);
    }
    if (message.emailAttribute !== "") {
      writer.uint32(58).string(message.emailAttribute);
    }
    if (message.emailVerifiedAttribute !== "") {
      writer.uint32(66).string(message.emailVerifiedAttribute);
    }
    if (message.phoneAttribute !== "") {
      writer.uint32(74).string(message.phoneAttribute);
    }
    if (message.phoneVerifiedAttribute !== "") {
      writer.uint32(82).string(message.phoneVerifiedAttribute);
    }
    if (message.preferredLanguageAttribute !== "") {
      writer.uint32(90).string(message.preferredLanguageAttribute);
    }
    if (message.avatarUrlAttribute !== "") {
      writer.uint32(98).string(message.avatarUrlAttribute);
    }
    if (message.profileAttribute !== "") {
      writer.uint32(106).string(message.profileAttribute);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LDAPAttributes {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLDAPAttributes();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idAttribute = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.firstNameAttribute = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.lastNameAttribute = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.displayNameAttribute = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.nickNameAttribute = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.preferredUsernameAttribute = reader.string();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.emailAttribute = reader.string();
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.emailVerifiedAttribute = reader.string();
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.phoneAttribute = reader.string();
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          message.phoneVerifiedAttribute = reader.string();
          continue;
        case 11:
          if (tag != 90) {
            break;
          }

          message.preferredLanguageAttribute = reader.string();
          continue;
        case 12:
          if (tag != 98) {
            break;
          }

          message.avatarUrlAttribute = reader.string();
          continue;
        case 13:
          if (tag != 106) {
            break;
          }

          message.profileAttribute = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LDAPAttributes {
    return {
      idAttribute: isSet(object.idAttribute) ? String(object.idAttribute) : "",
      firstNameAttribute: isSet(object.firstNameAttribute) ? String(object.firstNameAttribute) : "",
      lastNameAttribute: isSet(object.lastNameAttribute) ? String(object.lastNameAttribute) : "",
      displayNameAttribute: isSet(object.displayNameAttribute) ? String(object.displayNameAttribute) : "",
      nickNameAttribute: isSet(object.nickNameAttribute) ? String(object.nickNameAttribute) : "",
      preferredUsernameAttribute: isSet(object.preferredUsernameAttribute)
        ? String(object.preferredUsernameAttribute)
        : "",
      emailAttribute: isSet(object.emailAttribute) ? String(object.emailAttribute) : "",
      emailVerifiedAttribute: isSet(object.emailVerifiedAttribute) ? String(object.emailVerifiedAttribute) : "",
      phoneAttribute: isSet(object.phoneAttribute) ? String(object.phoneAttribute) : "",
      phoneVerifiedAttribute: isSet(object.phoneVerifiedAttribute) ? String(object.phoneVerifiedAttribute) : "",
      preferredLanguageAttribute: isSet(object.preferredLanguageAttribute)
        ? String(object.preferredLanguageAttribute)
        : "",
      avatarUrlAttribute: isSet(object.avatarUrlAttribute) ? String(object.avatarUrlAttribute) : "",
      profileAttribute: isSet(object.profileAttribute) ? String(object.profileAttribute) : "",
    };
  },

  toJSON(message: LDAPAttributes): unknown {
    const obj: any = {};
    message.idAttribute !== undefined && (obj.idAttribute = message.idAttribute);
    message.firstNameAttribute !== undefined && (obj.firstNameAttribute = message.firstNameAttribute);
    message.lastNameAttribute !== undefined && (obj.lastNameAttribute = message.lastNameAttribute);
    message.displayNameAttribute !== undefined && (obj.displayNameAttribute = message.displayNameAttribute);
    message.nickNameAttribute !== undefined && (obj.nickNameAttribute = message.nickNameAttribute);
    message.preferredUsernameAttribute !== undefined &&
      (obj.preferredUsernameAttribute = message.preferredUsernameAttribute);
    message.emailAttribute !== undefined && (obj.emailAttribute = message.emailAttribute);
    message.emailVerifiedAttribute !== undefined && (obj.emailVerifiedAttribute = message.emailVerifiedAttribute);
    message.phoneAttribute !== undefined && (obj.phoneAttribute = message.phoneAttribute);
    message.phoneVerifiedAttribute !== undefined && (obj.phoneVerifiedAttribute = message.phoneVerifiedAttribute);
    message.preferredLanguageAttribute !== undefined &&
      (obj.preferredLanguageAttribute = message.preferredLanguageAttribute);
    message.avatarUrlAttribute !== undefined && (obj.avatarUrlAttribute = message.avatarUrlAttribute);
    message.profileAttribute !== undefined && (obj.profileAttribute = message.profileAttribute);
    return obj;
  },

  create(base?: DeepPartial<LDAPAttributes>): LDAPAttributes {
    return LDAPAttributes.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LDAPAttributes>): LDAPAttributes {
    const message = createBaseLDAPAttributes();
    message.idAttribute = object.idAttribute ?? "";
    message.firstNameAttribute = object.firstNameAttribute ?? "";
    message.lastNameAttribute = object.lastNameAttribute ?? "";
    message.displayNameAttribute = object.displayNameAttribute ?? "";
    message.nickNameAttribute = object.nickNameAttribute ?? "";
    message.preferredUsernameAttribute = object.preferredUsernameAttribute ?? "";
    message.emailAttribute = object.emailAttribute ?? "";
    message.emailVerifiedAttribute = object.emailVerifiedAttribute ?? "";
    message.phoneAttribute = object.phoneAttribute ?? "";
    message.phoneVerifiedAttribute = object.phoneVerifiedAttribute ?? "";
    message.preferredLanguageAttribute = object.preferredLanguageAttribute ?? "";
    message.avatarUrlAttribute = object.avatarUrlAttribute ?? "";
    message.profileAttribute = object.profileAttribute ?? "";
    return message;
  },
};

function createBaseAzureADTenant(): AzureADTenant {
  return { tenantType: undefined, tenantId: undefined };
}

export const AzureADTenant = {
  encode(message: AzureADTenant, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tenantType !== undefined) {
      writer.uint32(8).int32(message.tenantType);
    }
    if (message.tenantId !== undefined) {
      writer.uint32(18).string(message.tenantId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AzureADTenant {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAzureADTenant();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.tenantType = reader.int32() as any;
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.tenantId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AzureADTenant {
    return {
      tenantType: isSet(object.tenantType) ? azureADTenantTypeFromJSON(object.tenantType) : undefined,
      tenantId: isSet(object.tenantId) ? String(object.tenantId) : undefined,
    };
  },

  toJSON(message: AzureADTenant): unknown {
    const obj: any = {};
    message.tenantType !== undefined &&
      (obj.tenantType = message.tenantType !== undefined ? azureADTenantTypeToJSON(message.tenantType) : undefined);
    message.tenantId !== undefined && (obj.tenantId = message.tenantId);
    return obj;
  },

  create(base?: DeepPartial<AzureADTenant>): AzureADTenant {
    return AzureADTenant.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AzureADTenant>): AzureADTenant {
    const message = createBaseAzureADTenant();
    message.tenantType = object.tenantType ?? undefined;
    message.tenantId = object.tenantId ?? undefined;
    return message;
  },
};

function createBaseAppleConfig(): AppleConfig {
  return { clientId: "", teamId: "", keyId: "", scopes: [] };
}

export const AppleConfig = {
  encode(message: AppleConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clientId !== "") {
      writer.uint32(10).string(message.clientId);
    }
    if (message.teamId !== "") {
      writer.uint32(18).string(message.teamId);
    }
    if (message.keyId !== "") {
      writer.uint32(26).string(message.keyId);
    }
    for (const v of message.scopes) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AppleConfig {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAppleConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.teamId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.keyId = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AppleConfig {
    return {
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      teamId: isSet(object.teamId) ? String(object.teamId) : "",
      keyId: isSet(object.keyId) ? String(object.keyId) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: AppleConfig): unknown {
    const obj: any = {};
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.teamId !== undefined && (obj.teamId = message.teamId);
    message.keyId !== undefined && (obj.keyId = message.keyId);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<AppleConfig>): AppleConfig {
    return AppleConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AppleConfig>): AppleConfig {
    const message = createBaseAppleConfig();
    message.clientId = object.clientId ?? "";
    message.teamId = object.teamId ?? "";
    message.keyId = object.keyId ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

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

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

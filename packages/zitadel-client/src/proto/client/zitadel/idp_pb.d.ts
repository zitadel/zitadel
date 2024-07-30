import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"


export class IDP extends jspb.Message {
  getId(): string;
  setId(value: string): IDP;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): IDP;
  hasDetails(): boolean;
  clearDetails(): IDP;

  getState(): IDPState;
  setState(value: IDPState): IDP;

  getName(): string;
  setName(value: string): IDP;

  getStylingType(): IDPStylingType;
  setStylingType(value: IDPStylingType): IDP;

  getOwner(): IDPOwnerType;
  setOwner(value: IDPOwnerType): IDP;

  getOidcConfig(): OIDCConfig | undefined;
  setOidcConfig(value?: OIDCConfig): IDP;
  hasOidcConfig(): boolean;
  clearOidcConfig(): IDP;

  getJwtConfig(): JWTConfig | undefined;
  setJwtConfig(value?: JWTConfig): IDP;
  hasJwtConfig(): boolean;
  clearJwtConfig(): IDP;

  getAutoRegister(): boolean;
  setAutoRegister(value: boolean): IDP;

  getConfigCase(): IDP.ConfigCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDP.AsObject;
  static toObject(includeInstance: boolean, msg: IDP): IDP.AsObject;
  static serializeBinaryToWriter(message: IDP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDP;
  static deserializeBinaryFromReader(message: IDP, reader: jspb.BinaryReader): IDP;
}

export namespace IDP {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: IDPState,
    name: string,
    stylingType: IDPStylingType,
    owner: IDPOwnerType,
    oidcConfig?: OIDCConfig.AsObject,
    jwtConfig?: JWTConfig.AsObject,
    autoRegister: boolean,
  }

  export enum ConfigCase { 
    CONFIG_NOT_SET = 0,
    OIDC_CONFIG = 7,
    JWT_CONFIG = 9,
  }
}

export class IDPUserLink extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): IDPUserLink;

  getIdpId(): string;
  setIdpId(value: string): IDPUserLink;

  getIdpName(): string;
  setIdpName(value: string): IDPUserLink;

  getProvidedUserId(): string;
  setProvidedUserId(value: string): IDPUserLink;

  getProvidedUserName(): string;
  setProvidedUserName(value: string): IDPUserLink;

  getIdpType(): IDPType;
  setIdpType(value: IDPType): IDPUserLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPUserLink.AsObject;
  static toObject(includeInstance: boolean, msg: IDPUserLink): IDPUserLink.AsObject;
  static serializeBinaryToWriter(message: IDPUserLink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPUserLink;
  static deserializeBinaryFromReader(message: IDPUserLink, reader: jspb.BinaryReader): IDPUserLink;
}

export namespace IDPUserLink {
  export type AsObject = {
    userId: string,
    idpId: string,
    idpName: string,
    providedUserId: string,
    providedUserName: string,
    idpType: IDPType,
  }
}

export class IDPLoginPolicyLink extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): IDPLoginPolicyLink;

  getIdpName(): string;
  setIdpName(value: string): IDPLoginPolicyLink;

  getIdpType(): IDPType;
  setIdpType(value: IDPType): IDPLoginPolicyLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPLoginPolicyLink.AsObject;
  static toObject(includeInstance: boolean, msg: IDPLoginPolicyLink): IDPLoginPolicyLink.AsObject;
  static serializeBinaryToWriter(message: IDPLoginPolicyLink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPLoginPolicyLink;
  static deserializeBinaryFromReader(message: IDPLoginPolicyLink, reader: jspb.BinaryReader): IDPLoginPolicyLink;
}

export namespace IDPLoginPolicyLink {
  export type AsObject = {
    idpId: string,
    idpName: string,
    idpType: IDPType,
  }
}

export class OIDCConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): OIDCConfig;

  getIssuer(): string;
  setIssuer(value: string): OIDCConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): OIDCConfig;
  clearScopesList(): OIDCConfig;
  addScopes(value: string, index?: number): OIDCConfig;

  getDisplayNameMapping(): OIDCMappingField;
  setDisplayNameMapping(value: OIDCMappingField): OIDCConfig;

  getUsernameMapping(): OIDCMappingField;
  setUsernameMapping(value: OIDCMappingField): OIDCConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OIDCConfig.AsObject;
  static toObject(includeInstance: boolean, msg: OIDCConfig): OIDCConfig.AsObject;
  static serializeBinaryToWriter(message: OIDCConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OIDCConfig;
  static deserializeBinaryFromReader(message: OIDCConfig, reader: jspb.BinaryReader): OIDCConfig;
}

export namespace OIDCConfig {
  export type AsObject = {
    clientId: string,
    issuer: string,
    scopesList: Array<string>,
    displayNameMapping: OIDCMappingField,
    usernameMapping: OIDCMappingField,
  }
}

export class JWTConfig extends jspb.Message {
  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): JWTConfig;

  getIssuer(): string;
  setIssuer(value: string): JWTConfig;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): JWTConfig;

  getHeaderName(): string;
  setHeaderName(value: string): JWTConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JWTConfig.AsObject;
  static toObject(includeInstance: boolean, msg: JWTConfig): JWTConfig.AsObject;
  static serializeBinaryToWriter(message: JWTConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JWTConfig;
  static deserializeBinaryFromReader(message: JWTConfig, reader: jspb.BinaryReader): JWTConfig;
}

export namespace JWTConfig {
  export type AsObject = {
    jwtEndpoint: string,
    issuer: string,
    keysEndpoint: string,
    headerName: string,
  }
}

export class IDPIDQuery extends jspb.Message {
  getId(): string;
  setId(value: string): IDPIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IDPIDQuery): IDPIDQuery.AsObject;
  static serializeBinaryToWriter(message: IDPIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPIDQuery;
  static deserializeBinaryFromReader(message: IDPIDQuery, reader: jspb.BinaryReader): IDPIDQuery;
}

export namespace IDPIDQuery {
  export type AsObject = {
    id: string,
  }
}

export class IDPNameQuery extends jspb.Message {
  getName(): string;
  setName(value: string): IDPNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): IDPNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IDPNameQuery): IDPNameQuery.AsObject;
  static serializeBinaryToWriter(message: IDPNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPNameQuery;
  static deserializeBinaryFromReader(message: IDPNameQuery, reader: jspb.BinaryReader): IDPNameQuery;
}

export namespace IDPNameQuery {
  export type AsObject = {
    name: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class IDPOwnerTypeQuery extends jspb.Message {
  getOwnerType(): IDPOwnerType;
  setOwnerType(value: IDPOwnerType): IDPOwnerTypeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPOwnerTypeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IDPOwnerTypeQuery): IDPOwnerTypeQuery.AsObject;
  static serializeBinaryToWriter(message: IDPOwnerTypeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPOwnerTypeQuery;
  static deserializeBinaryFromReader(message: IDPOwnerTypeQuery, reader: jspb.BinaryReader): IDPOwnerTypeQuery;
}

export namespace IDPOwnerTypeQuery {
  export type AsObject = {
    ownerType: IDPOwnerType,
  }
}

export class Provider extends jspb.Message {
  getId(): string;
  setId(value: string): Provider;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Provider;
  hasDetails(): boolean;
  clearDetails(): Provider;

  getState(): IDPState;
  setState(value: IDPState): Provider;

  getName(): string;
  setName(value: string): Provider;

  getOwner(): IDPOwnerType;
  setOwner(value: IDPOwnerType): Provider;

  getType(): ProviderType;
  setType(value: ProviderType): Provider;

  getConfig(): ProviderConfig | undefined;
  setConfig(value?: ProviderConfig): Provider;
  hasConfig(): boolean;
  clearConfig(): Provider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Provider.AsObject;
  static toObject(includeInstance: boolean, msg: Provider): Provider.AsObject;
  static serializeBinaryToWriter(message: Provider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Provider;
  static deserializeBinaryFromReader(message: Provider, reader: jspb.BinaryReader): Provider;
}

export namespace Provider {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: IDPState,
    name: string,
    owner: IDPOwnerType,
    type: ProviderType,
    config?: ProviderConfig.AsObject,
  }
}

export class ProviderConfig extends jspb.Message {
  getOptions(): Options | undefined;
  setOptions(value?: Options): ProviderConfig;
  hasOptions(): boolean;
  clearOptions(): ProviderConfig;

  getLdap(): LDAPConfig | undefined;
  setLdap(value?: LDAPConfig): ProviderConfig;
  hasLdap(): boolean;
  clearLdap(): ProviderConfig;

  getGoogle(): GoogleConfig | undefined;
  setGoogle(value?: GoogleConfig): ProviderConfig;
  hasGoogle(): boolean;
  clearGoogle(): ProviderConfig;

  getOauth(): OAuthConfig | undefined;
  setOauth(value?: OAuthConfig): ProviderConfig;
  hasOauth(): boolean;
  clearOauth(): ProviderConfig;

  getOidc(): GenericOIDCConfig | undefined;
  setOidc(value?: GenericOIDCConfig): ProviderConfig;
  hasOidc(): boolean;
  clearOidc(): ProviderConfig;

  getJwt(): JWTConfig | undefined;
  setJwt(value?: JWTConfig): ProviderConfig;
  hasJwt(): boolean;
  clearJwt(): ProviderConfig;

  getGithub(): GitHubConfig | undefined;
  setGithub(value?: GitHubConfig): ProviderConfig;
  hasGithub(): boolean;
  clearGithub(): ProviderConfig;

  getGithubEs(): GitHubEnterpriseServerConfig | undefined;
  setGithubEs(value?: GitHubEnterpriseServerConfig): ProviderConfig;
  hasGithubEs(): boolean;
  clearGithubEs(): ProviderConfig;

  getGitlab(): GitLabConfig | undefined;
  setGitlab(value?: GitLabConfig): ProviderConfig;
  hasGitlab(): boolean;
  clearGitlab(): ProviderConfig;

  getGitlabSelfHosted(): GitLabSelfHostedConfig | undefined;
  setGitlabSelfHosted(value?: GitLabSelfHostedConfig): ProviderConfig;
  hasGitlabSelfHosted(): boolean;
  clearGitlabSelfHosted(): ProviderConfig;

  getAzureAd(): AzureADConfig | undefined;
  setAzureAd(value?: AzureADConfig): ProviderConfig;
  hasAzureAd(): boolean;
  clearAzureAd(): ProviderConfig;

  getApple(): AppleConfig | undefined;
  setApple(value?: AppleConfig): ProviderConfig;
  hasApple(): boolean;
  clearApple(): ProviderConfig;

  getSaml(): SAMLConfig | undefined;
  setSaml(value?: SAMLConfig): ProviderConfig;
  hasSaml(): boolean;
  clearSaml(): ProviderConfig;

  getConfigCase(): ProviderConfig.ConfigCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProviderConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ProviderConfig): ProviderConfig.AsObject;
  static serializeBinaryToWriter(message: ProviderConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProviderConfig;
  static deserializeBinaryFromReader(message: ProviderConfig, reader: jspb.BinaryReader): ProviderConfig;
}

export namespace ProviderConfig {
  export type AsObject = {
    options?: Options.AsObject,
    ldap?: LDAPConfig.AsObject,
    google?: GoogleConfig.AsObject,
    oauth?: OAuthConfig.AsObject,
    oidc?: GenericOIDCConfig.AsObject,
    jwt?: JWTConfig.AsObject,
    github?: GitHubConfig.AsObject,
    githubEs?: GitHubEnterpriseServerConfig.AsObject,
    gitlab?: GitLabConfig.AsObject,
    gitlabSelfHosted?: GitLabSelfHostedConfig.AsObject,
    azureAd?: AzureADConfig.AsObject,
    apple?: AppleConfig.AsObject,
    saml?: SAMLConfig.AsObject,
  }

  export enum ConfigCase { 
    CONFIG_NOT_SET = 0,
    LDAP = 2,
    GOOGLE = 3,
    OAUTH = 4,
    OIDC = 5,
    JWT = 6,
    GITHUB = 7,
    GITHUB_ES = 8,
    GITLAB = 9,
    GITLAB_SELF_HOSTED = 10,
    AZURE_AD = 11,
    APPLE = 12,
    SAML = 13,
  }
}

export class OAuthConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): OAuthConfig;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): OAuthConfig;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): OAuthConfig;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): OAuthConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): OAuthConfig;
  clearScopesList(): OAuthConfig;
  addScopes(value: string, index?: number): OAuthConfig;

  getIdAttribute(): string;
  setIdAttribute(value: string): OAuthConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OAuthConfig.AsObject;
  static toObject(includeInstance: boolean, msg: OAuthConfig): OAuthConfig.AsObject;
  static serializeBinaryToWriter(message: OAuthConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OAuthConfig;
  static deserializeBinaryFromReader(message: OAuthConfig, reader: jspb.BinaryReader): OAuthConfig;
}

export namespace OAuthConfig {
  export type AsObject = {
    clientId: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    idAttribute: string,
  }
}

export class GenericOIDCConfig extends jspb.Message {
  getIssuer(): string;
  setIssuer(value: string): GenericOIDCConfig;

  getClientId(): string;
  setClientId(value: string): GenericOIDCConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): GenericOIDCConfig;
  clearScopesList(): GenericOIDCConfig;
  addScopes(value: string, index?: number): GenericOIDCConfig;

  getIsIdTokenMapping(): boolean;
  setIsIdTokenMapping(value: boolean): GenericOIDCConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenericOIDCConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GenericOIDCConfig): GenericOIDCConfig.AsObject;
  static serializeBinaryToWriter(message: GenericOIDCConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenericOIDCConfig;
  static deserializeBinaryFromReader(message: GenericOIDCConfig, reader: jspb.BinaryReader): GenericOIDCConfig;
}

export namespace GenericOIDCConfig {
  export type AsObject = {
    issuer: string,
    clientId: string,
    scopesList: Array<string>,
    isIdTokenMapping: boolean,
  }
}

export class GitHubConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): GitHubConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): GitHubConfig;
  clearScopesList(): GitHubConfig;
  addScopes(value: string, index?: number): GitHubConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GitHubConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GitHubConfig): GitHubConfig.AsObject;
  static serializeBinaryToWriter(message: GitHubConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GitHubConfig;
  static deserializeBinaryFromReader(message: GitHubConfig, reader: jspb.BinaryReader): GitHubConfig;
}

export namespace GitHubConfig {
  export type AsObject = {
    clientId: string,
    scopesList: Array<string>,
  }
}

export class GitHubEnterpriseServerConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): GitHubEnterpriseServerConfig;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): GitHubEnterpriseServerConfig;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): GitHubEnterpriseServerConfig;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): GitHubEnterpriseServerConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): GitHubEnterpriseServerConfig;
  clearScopesList(): GitHubEnterpriseServerConfig;
  addScopes(value: string, index?: number): GitHubEnterpriseServerConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GitHubEnterpriseServerConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GitHubEnterpriseServerConfig): GitHubEnterpriseServerConfig.AsObject;
  static serializeBinaryToWriter(message: GitHubEnterpriseServerConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GitHubEnterpriseServerConfig;
  static deserializeBinaryFromReader(message: GitHubEnterpriseServerConfig, reader: jspb.BinaryReader): GitHubEnterpriseServerConfig;
}

export namespace GitHubEnterpriseServerConfig {
  export type AsObject = {
    clientId: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
  }
}

export class GoogleConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): GoogleConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): GoogleConfig;
  clearScopesList(): GoogleConfig;
  addScopes(value: string, index?: number): GoogleConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GoogleConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GoogleConfig): GoogleConfig.AsObject;
  static serializeBinaryToWriter(message: GoogleConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GoogleConfig;
  static deserializeBinaryFromReader(message: GoogleConfig, reader: jspb.BinaryReader): GoogleConfig;
}

export namespace GoogleConfig {
  export type AsObject = {
    clientId: string,
    scopesList: Array<string>,
  }
}

export class GitLabConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): GitLabConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): GitLabConfig;
  clearScopesList(): GitLabConfig;
  addScopes(value: string, index?: number): GitLabConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GitLabConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GitLabConfig): GitLabConfig.AsObject;
  static serializeBinaryToWriter(message: GitLabConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GitLabConfig;
  static deserializeBinaryFromReader(message: GitLabConfig, reader: jspb.BinaryReader): GitLabConfig;
}

export namespace GitLabConfig {
  export type AsObject = {
    clientId: string,
    scopesList: Array<string>,
  }
}

export class GitLabSelfHostedConfig extends jspb.Message {
  getIssuer(): string;
  setIssuer(value: string): GitLabSelfHostedConfig;

  getClientId(): string;
  setClientId(value: string): GitLabSelfHostedConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): GitLabSelfHostedConfig;
  clearScopesList(): GitLabSelfHostedConfig;
  addScopes(value: string, index?: number): GitLabSelfHostedConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GitLabSelfHostedConfig.AsObject;
  static toObject(includeInstance: boolean, msg: GitLabSelfHostedConfig): GitLabSelfHostedConfig.AsObject;
  static serializeBinaryToWriter(message: GitLabSelfHostedConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GitLabSelfHostedConfig;
  static deserializeBinaryFromReader(message: GitLabSelfHostedConfig, reader: jspb.BinaryReader): GitLabSelfHostedConfig;
}

export namespace GitLabSelfHostedConfig {
  export type AsObject = {
    issuer: string,
    clientId: string,
    scopesList: Array<string>,
  }
}

export class LDAPConfig extends jspb.Message {
  getServersList(): Array<string>;
  setServersList(value: Array<string>): LDAPConfig;
  clearServersList(): LDAPConfig;
  addServers(value: string, index?: number): LDAPConfig;

  getStartTls(): boolean;
  setStartTls(value: boolean): LDAPConfig;

  getBaseDn(): string;
  setBaseDn(value: string): LDAPConfig;

  getBindDn(): string;
  setBindDn(value: string): LDAPConfig;

  getUserBase(): string;
  setUserBase(value: string): LDAPConfig;

  getUserObjectClassesList(): Array<string>;
  setUserObjectClassesList(value: Array<string>): LDAPConfig;
  clearUserObjectClassesList(): LDAPConfig;
  addUserObjectClasses(value: string, index?: number): LDAPConfig;

  getUserFiltersList(): Array<string>;
  setUserFiltersList(value: Array<string>): LDAPConfig;
  clearUserFiltersList(): LDAPConfig;
  addUserFilters(value: string, index?: number): LDAPConfig;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): LDAPConfig;
  hasTimeout(): boolean;
  clearTimeout(): LDAPConfig;

  getAttributes(): LDAPAttributes | undefined;
  setAttributes(value?: LDAPAttributes): LDAPConfig;
  hasAttributes(): boolean;
  clearAttributes(): LDAPConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LDAPConfig.AsObject;
  static toObject(includeInstance: boolean, msg: LDAPConfig): LDAPConfig.AsObject;
  static serializeBinaryToWriter(message: LDAPConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LDAPConfig;
  static deserializeBinaryFromReader(message: LDAPConfig, reader: jspb.BinaryReader): LDAPConfig;
}

export namespace LDAPConfig {
  export type AsObject = {
    serversList: Array<string>,
    startTls: boolean,
    baseDn: string,
    bindDn: string,
    userBase: string,
    userObjectClassesList: Array<string>,
    userFiltersList: Array<string>,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    attributes?: LDAPAttributes.AsObject,
  }
}

export class SAMLConfig extends jspb.Message {
  getMetadataXml(): Uint8Array | string;
  getMetadataXml_asU8(): Uint8Array;
  getMetadataXml_asB64(): string;
  setMetadataXml(value: Uint8Array | string): SAMLConfig;

  getBinding(): SAMLBinding;
  setBinding(value: SAMLBinding): SAMLConfig;

  getWithSignedRequest(): boolean;
  setWithSignedRequest(value: boolean): SAMLConfig;

  getNameIdFormat(): SAMLNameIDFormat;
  setNameIdFormat(value: SAMLNameIDFormat): SAMLConfig;

  getTransientMappingAttributeName(): string;
  setTransientMappingAttributeName(value: string): SAMLConfig;
  hasTransientMappingAttributeName(): boolean;
  clearTransientMappingAttributeName(): SAMLConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SAMLConfig.AsObject;
  static toObject(includeInstance: boolean, msg: SAMLConfig): SAMLConfig.AsObject;
  static serializeBinaryToWriter(message: SAMLConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SAMLConfig;
  static deserializeBinaryFromReader(message: SAMLConfig, reader: jspb.BinaryReader): SAMLConfig;
}

export namespace SAMLConfig {
  export type AsObject = {
    metadataXml: Uint8Array | string,
    binding: SAMLBinding,
    withSignedRequest: boolean,
    nameIdFormat: SAMLNameIDFormat,
    transientMappingAttributeName?: string,
  }

  export enum TransientMappingAttributeNameCase { 
    _TRANSIENT_MAPPING_ATTRIBUTE_NAME_NOT_SET = 0,
    TRANSIENT_MAPPING_ATTRIBUTE_NAME = 5,
  }
}

export class AzureADConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): AzureADConfig;

  getTenant(): AzureADTenant | undefined;
  setTenant(value?: AzureADTenant): AzureADConfig;
  hasTenant(): boolean;
  clearTenant(): AzureADConfig;

  getEmailVerified(): boolean;
  setEmailVerified(value: boolean): AzureADConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AzureADConfig;
  clearScopesList(): AzureADConfig;
  addScopes(value: string, index?: number): AzureADConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AzureADConfig.AsObject;
  static toObject(includeInstance: boolean, msg: AzureADConfig): AzureADConfig.AsObject;
  static serializeBinaryToWriter(message: AzureADConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AzureADConfig;
  static deserializeBinaryFromReader(message: AzureADConfig, reader: jspb.BinaryReader): AzureADConfig;
}

export namespace AzureADConfig {
  export type AsObject = {
    clientId: string,
    tenant?: AzureADTenant.AsObject,
    emailVerified: boolean,
    scopesList: Array<string>,
  }
}

export class Options extends jspb.Message {
  getIsLinkingAllowed(): boolean;
  setIsLinkingAllowed(value: boolean): Options;

  getIsCreationAllowed(): boolean;
  setIsCreationAllowed(value: boolean): Options;

  getIsAutoCreation(): boolean;
  setIsAutoCreation(value: boolean): Options;

  getIsAutoUpdate(): boolean;
  setIsAutoUpdate(value: boolean): Options;

  getAutoLinking(): AutoLinkingOption;
  setAutoLinking(value: AutoLinkingOption): Options;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Options.AsObject;
  static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
  static serializeBinaryToWriter(message: Options, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Options;
  static deserializeBinaryFromReader(message: Options, reader: jspb.BinaryReader): Options;
}

export namespace Options {
  export type AsObject = {
    isLinkingAllowed: boolean,
    isCreationAllowed: boolean,
    isAutoCreation: boolean,
    isAutoUpdate: boolean,
    autoLinking: AutoLinkingOption,
  }
}

export class LDAPAttributes extends jspb.Message {
  getIdAttribute(): string;
  setIdAttribute(value: string): LDAPAttributes;

  getFirstNameAttribute(): string;
  setFirstNameAttribute(value: string): LDAPAttributes;

  getLastNameAttribute(): string;
  setLastNameAttribute(value: string): LDAPAttributes;

  getDisplayNameAttribute(): string;
  setDisplayNameAttribute(value: string): LDAPAttributes;

  getNickNameAttribute(): string;
  setNickNameAttribute(value: string): LDAPAttributes;

  getPreferredUsernameAttribute(): string;
  setPreferredUsernameAttribute(value: string): LDAPAttributes;

  getEmailAttribute(): string;
  setEmailAttribute(value: string): LDAPAttributes;

  getEmailVerifiedAttribute(): string;
  setEmailVerifiedAttribute(value: string): LDAPAttributes;

  getPhoneAttribute(): string;
  setPhoneAttribute(value: string): LDAPAttributes;

  getPhoneVerifiedAttribute(): string;
  setPhoneVerifiedAttribute(value: string): LDAPAttributes;

  getPreferredLanguageAttribute(): string;
  setPreferredLanguageAttribute(value: string): LDAPAttributes;

  getAvatarUrlAttribute(): string;
  setAvatarUrlAttribute(value: string): LDAPAttributes;

  getProfileAttribute(): string;
  setProfileAttribute(value: string): LDAPAttributes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LDAPAttributes.AsObject;
  static toObject(includeInstance: boolean, msg: LDAPAttributes): LDAPAttributes.AsObject;
  static serializeBinaryToWriter(message: LDAPAttributes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LDAPAttributes;
  static deserializeBinaryFromReader(message: LDAPAttributes, reader: jspb.BinaryReader): LDAPAttributes;
}

export namespace LDAPAttributes {
  export type AsObject = {
    idAttribute: string,
    firstNameAttribute: string,
    lastNameAttribute: string,
    displayNameAttribute: string,
    nickNameAttribute: string,
    preferredUsernameAttribute: string,
    emailAttribute: string,
    emailVerifiedAttribute: string,
    phoneAttribute: string,
    phoneVerifiedAttribute: string,
    preferredLanguageAttribute: string,
    avatarUrlAttribute: string,
    profileAttribute: string,
  }
}

export class AzureADTenant extends jspb.Message {
  getTenantType(): AzureADTenantType;
  setTenantType(value: AzureADTenantType): AzureADTenant;

  getTenantId(): string;
  setTenantId(value: string): AzureADTenant;

  getTypeCase(): AzureADTenant.TypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AzureADTenant.AsObject;
  static toObject(includeInstance: boolean, msg: AzureADTenant): AzureADTenant.AsObject;
  static serializeBinaryToWriter(message: AzureADTenant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AzureADTenant;
  static deserializeBinaryFromReader(message: AzureADTenant, reader: jspb.BinaryReader): AzureADTenant;
}

export namespace AzureADTenant {
  export type AsObject = {
    tenantType: AzureADTenantType,
    tenantId: string,
  }

  export enum TypeCase { 
    TYPE_NOT_SET = 0,
    TENANT_TYPE = 1,
    TENANT_ID = 2,
  }
}

export class AppleConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): AppleConfig;

  getTeamId(): string;
  setTeamId(value: string): AppleConfig;

  getKeyId(): string;
  setKeyId(value: string): AppleConfig;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AppleConfig;
  clearScopesList(): AppleConfig;
  addScopes(value: string, index?: number): AppleConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AppleConfig.AsObject;
  static toObject(includeInstance: boolean, msg: AppleConfig): AppleConfig.AsObject;
  static serializeBinaryToWriter(message: AppleConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AppleConfig;
  static deserializeBinaryFromReader(message: AppleConfig, reader: jspb.BinaryReader): AppleConfig;
}

export namespace AppleConfig {
  export type AsObject = {
    clientId: string,
    teamId: string,
    keyId: string,
    scopesList: Array<string>,
  }
}

export enum IDPState { 
  IDP_STATE_UNSPECIFIED = 0,
  IDP_STATE_ACTIVE = 1,
  IDP_STATE_INACTIVE = 2,
}
export enum IDPStylingType { 
  STYLING_TYPE_UNSPECIFIED = 0,
  STYLING_TYPE_GOOGLE = 1,
}
export enum IDPType { 
  IDP_TYPE_UNSPECIFIED = 0,
  IDP_TYPE_OIDC = 1,
  IDP_TYPE_JWT = 3,
}
export enum IDPOwnerType { 
  IDP_OWNER_TYPE_UNSPECIFIED = 0,
  IDP_OWNER_TYPE_SYSTEM = 1,
  IDP_OWNER_TYPE_ORG = 2,
}
export enum OIDCMappingField { 
  OIDC_MAPPING_FIELD_UNSPECIFIED = 0,
  OIDC_MAPPING_FIELD_PREFERRED_USERNAME = 1,
  OIDC_MAPPING_FIELD_EMAIL = 2,
}
export enum IDPFieldName { 
  IDP_FIELD_NAME_UNSPECIFIED = 0,
  IDP_FIELD_NAME_NAME = 1,
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
}
export enum SAMLBinding { 
  SAML_BINDING_UNSPECIFIED = 0,
  SAML_BINDING_POST = 1,
  SAML_BINDING_REDIRECT = 2,
  SAML_BINDING_ARTIFACT = 3,
}
export enum SAMLNameIDFormat { 
  SAML_NAME_ID_FORMAT_UNSPECIFIED = 0,
  SAML_NAME_ID_FORMAT_EMAIL_ADDRESS = 1,
  SAML_NAME_ID_FORMAT_PERSISTENT = 2,
  SAML_NAME_ID_FORMAT_TRANSIENT = 3,
}
export enum AutoLinkingOption { 
  AUTO_LINKING_OPTION_UNSPECIFIED = 0,
  AUTO_LINKING_OPTION_USERNAME = 1,
  AUTO_LINKING_OPTION_EMAIL = 2,
}
export enum AzureADTenantType { 
  AZURE_AD_TENANT_TYPE_COMMON = 0,
  AZURE_AD_TENANT_TYPE_ORGANISATIONS = 1,
  AZURE_AD_TENANT_TYPE_CONSUMERS = 2,
}

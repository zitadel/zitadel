import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as zitadel_message_pb from '../zitadel/message_pb'; // proto import: "zitadel/message.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class App extends jspb.Message {
  getId(): string;
  setId(value: string): App;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): App;
  hasDetails(): boolean;
  clearDetails(): App;

  getState(): AppState;
  setState(value: AppState): App;

  getName(): string;
  setName(value: string): App;

  getOidcConfig(): OIDCConfig | undefined;
  setOidcConfig(value?: OIDCConfig): App;
  hasOidcConfig(): boolean;
  clearOidcConfig(): App;

  getApiConfig(): APIConfig | undefined;
  setApiConfig(value?: APIConfig): App;
  hasApiConfig(): boolean;
  clearApiConfig(): App;

  getSamlConfig(): SAMLConfig | undefined;
  setSamlConfig(value?: SAMLConfig): App;
  hasSamlConfig(): boolean;
  clearSamlConfig(): App;

  getConfigCase(): App.ConfigCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): App.AsObject;
  static toObject(includeInstance: boolean, msg: App): App.AsObject;
  static serializeBinaryToWriter(message: App, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): App;
  static deserializeBinaryFromReader(message: App, reader: jspb.BinaryReader): App;
}

export namespace App {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: AppState,
    name: string,
    oidcConfig?: OIDCConfig.AsObject,
    apiConfig?: APIConfig.AsObject,
    samlConfig?: SAMLConfig.AsObject,
  }

  export enum ConfigCase { 
    CONFIG_NOT_SET = 0,
    OIDC_CONFIG = 5,
    API_CONFIG = 6,
    SAML_CONFIG = 7,
  }
}

export class AppQuery extends jspb.Message {
  getNameQuery(): AppNameQuery | undefined;
  setNameQuery(value?: AppNameQuery): AppQuery;
  hasNameQuery(): boolean;
  clearNameQuery(): AppQuery;

  getQueryCase(): AppQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AppQuery.AsObject;
  static toObject(includeInstance: boolean, msg: AppQuery): AppQuery.AsObject;
  static serializeBinaryToWriter(message: AppQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AppQuery;
  static deserializeBinaryFromReader(message: AppQuery, reader: jspb.BinaryReader): AppQuery;
}

export namespace AppQuery {
  export type AsObject = {
    nameQuery?: AppNameQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    NAME_QUERY = 1,
  }
}

export class AppNameQuery extends jspb.Message {
  getName(): string;
  setName(value: string): AppNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): AppNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AppNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: AppNameQuery): AppNameQuery.AsObject;
  static serializeBinaryToWriter(message: AppNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AppNameQuery;
  static deserializeBinaryFromReader(message: AppNameQuery, reader: jspb.BinaryReader): AppNameQuery;
}

export namespace AppNameQuery {
  export type AsObject = {
    name: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class OIDCConfig extends jspb.Message {
  getRedirectUrisList(): Array<string>;
  setRedirectUrisList(value: Array<string>): OIDCConfig;
  clearRedirectUrisList(): OIDCConfig;
  addRedirectUris(value: string, index?: number): OIDCConfig;

  getResponseTypesList(): Array<OIDCResponseType>;
  setResponseTypesList(value: Array<OIDCResponseType>): OIDCConfig;
  clearResponseTypesList(): OIDCConfig;
  addResponseTypes(value: OIDCResponseType, index?: number): OIDCConfig;

  getGrantTypesList(): Array<OIDCGrantType>;
  setGrantTypesList(value: Array<OIDCGrantType>): OIDCConfig;
  clearGrantTypesList(): OIDCConfig;
  addGrantTypes(value: OIDCGrantType, index?: number): OIDCConfig;

  getAppType(): OIDCAppType;
  setAppType(value: OIDCAppType): OIDCConfig;

  getClientId(): string;
  setClientId(value: string): OIDCConfig;

  getAuthMethodType(): OIDCAuthMethodType;
  setAuthMethodType(value: OIDCAuthMethodType): OIDCConfig;

  getPostLogoutRedirectUrisList(): Array<string>;
  setPostLogoutRedirectUrisList(value: Array<string>): OIDCConfig;
  clearPostLogoutRedirectUrisList(): OIDCConfig;
  addPostLogoutRedirectUris(value: string, index?: number): OIDCConfig;

  getVersion(): OIDCVersion;
  setVersion(value: OIDCVersion): OIDCConfig;

  getNoneCompliant(): boolean;
  setNoneCompliant(value: boolean): OIDCConfig;

  getComplianceProblemsList(): Array<zitadel_message_pb.LocalizedMessage>;
  setComplianceProblemsList(value: Array<zitadel_message_pb.LocalizedMessage>): OIDCConfig;
  clearComplianceProblemsList(): OIDCConfig;
  addComplianceProblems(value?: zitadel_message_pb.LocalizedMessage, index?: number): zitadel_message_pb.LocalizedMessage;

  getDevMode(): boolean;
  setDevMode(value: boolean): OIDCConfig;

  getAccessTokenType(): OIDCTokenType;
  setAccessTokenType(value: OIDCTokenType): OIDCConfig;

  getAccessTokenRoleAssertion(): boolean;
  setAccessTokenRoleAssertion(value: boolean): OIDCConfig;

  getIdTokenRoleAssertion(): boolean;
  setIdTokenRoleAssertion(value: boolean): OIDCConfig;

  getIdTokenUserinfoAssertion(): boolean;
  setIdTokenUserinfoAssertion(value: boolean): OIDCConfig;

  getClockSkew(): google_protobuf_duration_pb.Duration | undefined;
  setClockSkew(value?: google_protobuf_duration_pb.Duration): OIDCConfig;
  hasClockSkew(): boolean;
  clearClockSkew(): OIDCConfig;

  getAdditionalOriginsList(): Array<string>;
  setAdditionalOriginsList(value: Array<string>): OIDCConfig;
  clearAdditionalOriginsList(): OIDCConfig;
  addAdditionalOrigins(value: string, index?: number): OIDCConfig;

  getAllowedOriginsList(): Array<string>;
  setAllowedOriginsList(value: Array<string>): OIDCConfig;
  clearAllowedOriginsList(): OIDCConfig;
  addAllowedOrigins(value: string, index?: number): OIDCConfig;

  getSkipNativeAppSuccessPage(): boolean;
  setSkipNativeAppSuccessPage(value: boolean): OIDCConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OIDCConfig.AsObject;
  static toObject(includeInstance: boolean, msg: OIDCConfig): OIDCConfig.AsObject;
  static serializeBinaryToWriter(message: OIDCConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OIDCConfig;
  static deserializeBinaryFromReader(message: OIDCConfig, reader: jspb.BinaryReader): OIDCConfig;
}

export namespace OIDCConfig {
  export type AsObject = {
    redirectUrisList: Array<string>,
    responseTypesList: Array<OIDCResponseType>,
    grantTypesList: Array<OIDCGrantType>,
    appType: OIDCAppType,
    clientId: string,
    authMethodType: OIDCAuthMethodType,
    postLogoutRedirectUrisList: Array<string>,
    version: OIDCVersion,
    noneCompliant: boolean,
    complianceProblemsList: Array<zitadel_message_pb.LocalizedMessage.AsObject>,
    devMode: boolean,
    accessTokenType: OIDCTokenType,
    accessTokenRoleAssertion: boolean,
    idTokenRoleAssertion: boolean,
    idTokenUserinfoAssertion: boolean,
    clockSkew?: google_protobuf_duration_pb.Duration.AsObject,
    additionalOriginsList: Array<string>,
    allowedOriginsList: Array<string>,
    skipNativeAppSuccessPage: boolean,
  }
}

export class SAMLConfig extends jspb.Message {
  getMetadataXml(): Uint8Array | string;
  getMetadataXml_asU8(): Uint8Array;
  getMetadataXml_asB64(): string;
  setMetadataXml(value: Uint8Array | string): SAMLConfig;

  getMetadataUrl(): string;
  setMetadataUrl(value: string): SAMLConfig;

  getMetadataCase(): SAMLConfig.MetadataCase;

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
    metadataUrl: string,
  }

  export enum MetadataCase { 
    METADATA_NOT_SET = 0,
    METADATA_XML = 1,
    METADATA_URL = 2,
  }
}

export class APIConfig extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): APIConfig;

  getAuthMethodType(): APIAuthMethodType;
  setAuthMethodType(value: APIAuthMethodType): APIConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): APIConfig.AsObject;
  static toObject(includeInstance: boolean, msg: APIConfig): APIConfig.AsObject;
  static serializeBinaryToWriter(message: APIConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): APIConfig;
  static deserializeBinaryFromReader(message: APIConfig, reader: jspb.BinaryReader): APIConfig;
}

export namespace APIConfig {
  export type AsObject = {
    clientId: string,
    authMethodType: APIAuthMethodType,
  }
}

export enum AppState { 
  APP_STATE_UNSPECIFIED = 0,
  APP_STATE_ACTIVE = 1,
  APP_STATE_INACTIVE = 2,
}
export enum OIDCResponseType { 
  OIDC_RESPONSE_TYPE_CODE = 0,
  OIDC_RESPONSE_TYPE_ID_TOKEN = 1,
  OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN = 2,
}
export enum OIDCGrantType { 
  OIDC_GRANT_TYPE_AUTHORIZATION_CODE = 0,
  OIDC_GRANT_TYPE_IMPLICIT = 1,
  OIDC_GRANT_TYPE_REFRESH_TOKEN = 2,
  OIDC_GRANT_TYPE_DEVICE_CODE = 3,
  OIDC_GRANT_TYPE_TOKEN_EXCHANGE = 4,
}
export enum OIDCAppType { 
  OIDC_APP_TYPE_WEB = 0,
  OIDC_APP_TYPE_USER_AGENT = 1,
  OIDC_APP_TYPE_NATIVE = 2,
}
export enum OIDCAuthMethodType { 
  OIDC_AUTH_METHOD_TYPE_BASIC = 0,
  OIDC_AUTH_METHOD_TYPE_POST = 1,
  OIDC_AUTH_METHOD_TYPE_NONE = 2,
  OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT = 3,
}
export enum OIDCVersion { 
  OIDC_VERSION_1_0 = 0,
}
export enum OIDCTokenType { 
  OIDC_TOKEN_TYPE_BEARER = 0,
  OIDC_TOKEN_TYPE_JWT = 1,
}
export enum APIAuthMethodType { 
  API_AUTH_METHOD_TYPE_BASIC = 0,
  API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT = 1,
}

import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as zitadel_settings_v2beta_settings_pb from '../../../zitadel/settings/v2beta/settings_pb'; // proto import: "zitadel/settings/v2beta/settings.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"


export class LoginSettings extends jspb.Message {
  getAllowUsernamePassword(): boolean;
  setAllowUsernamePassword(value: boolean): LoginSettings;

  getAllowRegister(): boolean;
  setAllowRegister(value: boolean): LoginSettings;

  getAllowExternalIdp(): boolean;
  setAllowExternalIdp(value: boolean): LoginSettings;

  getForceMfa(): boolean;
  setForceMfa(value: boolean): LoginSettings;

  getPasskeysType(): PasskeysType;
  setPasskeysType(value: PasskeysType): LoginSettings;

  getHidePasswordReset(): boolean;
  setHidePasswordReset(value: boolean): LoginSettings;

  getIgnoreUnknownUsernames(): boolean;
  setIgnoreUnknownUsernames(value: boolean): LoginSettings;

  getDefaultRedirectUri(): string;
  setDefaultRedirectUri(value: string): LoginSettings;

  getPasswordCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setPasswordCheckLifetime(value?: google_protobuf_duration_pb.Duration): LoginSettings;
  hasPasswordCheckLifetime(): boolean;
  clearPasswordCheckLifetime(): LoginSettings;

  getExternalLoginCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setExternalLoginCheckLifetime(value?: google_protobuf_duration_pb.Duration): LoginSettings;
  hasExternalLoginCheckLifetime(): boolean;
  clearExternalLoginCheckLifetime(): LoginSettings;

  getMfaInitSkipLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMfaInitSkipLifetime(value?: google_protobuf_duration_pb.Duration): LoginSettings;
  hasMfaInitSkipLifetime(): boolean;
  clearMfaInitSkipLifetime(): LoginSettings;

  getSecondFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setSecondFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): LoginSettings;
  hasSecondFactorCheckLifetime(): boolean;
  clearSecondFactorCheckLifetime(): LoginSettings;

  getMultiFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMultiFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): LoginSettings;
  hasMultiFactorCheckLifetime(): boolean;
  clearMultiFactorCheckLifetime(): LoginSettings;

  getSecondFactorsList(): Array<SecondFactorType>;
  setSecondFactorsList(value: Array<SecondFactorType>): LoginSettings;
  clearSecondFactorsList(): LoginSettings;
  addSecondFactors(value: SecondFactorType, index?: number): LoginSettings;

  getMultiFactorsList(): Array<MultiFactorType>;
  setMultiFactorsList(value: Array<MultiFactorType>): LoginSettings;
  clearMultiFactorsList(): LoginSettings;
  addMultiFactors(value: MultiFactorType, index?: number): LoginSettings;

  getAllowDomainDiscovery(): boolean;
  setAllowDomainDiscovery(value: boolean): LoginSettings;

  getDisableLoginWithEmail(): boolean;
  setDisableLoginWithEmail(value: boolean): LoginSettings;

  getDisableLoginWithPhone(): boolean;
  setDisableLoginWithPhone(value: boolean): LoginSettings;

  getResourceOwnerType(): zitadel_settings_v2beta_settings_pb.ResourceOwnerType;
  setResourceOwnerType(value: zitadel_settings_v2beta_settings_pb.ResourceOwnerType): LoginSettings;

  getForceMfaLocalOnly(): boolean;
  setForceMfaLocalOnly(value: boolean): LoginSettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginSettings.AsObject;
  static toObject(includeInstance: boolean, msg: LoginSettings): LoginSettings.AsObject;
  static serializeBinaryToWriter(message: LoginSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginSettings;
  static deserializeBinaryFromReader(message: LoginSettings, reader: jspb.BinaryReader): LoginSettings;
}

export namespace LoginSettings {
  export type AsObject = {
    allowUsernamePassword: boolean,
    allowRegister: boolean,
    allowExternalIdp: boolean,
    forceMfa: boolean,
    passkeysType: PasskeysType,
    hidePasswordReset: boolean,
    ignoreUnknownUsernames: boolean,
    defaultRedirectUri: string,
    passwordCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    externalLoginCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    mfaInitSkipLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    secondFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    multiFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    secondFactorsList: Array<SecondFactorType>,
    multiFactorsList: Array<MultiFactorType>,
    allowDomainDiscovery: boolean,
    disableLoginWithEmail: boolean,
    disableLoginWithPhone: boolean,
    resourceOwnerType: zitadel_settings_v2beta_settings_pb.ResourceOwnerType,
    forceMfaLocalOnly: boolean,
  }
}

export class IdentityProvider extends jspb.Message {
  getId(): string;
  setId(value: string): IdentityProvider;

  getName(): string;
  setName(value: string): IdentityProvider;

  getType(): IdentityProviderType;
  setType(value: IdentityProviderType): IdentityProvider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IdentityProvider.AsObject;
  static toObject(includeInstance: boolean, msg: IdentityProvider): IdentityProvider.AsObject;
  static serializeBinaryToWriter(message: IdentityProvider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IdentityProvider;
  static deserializeBinaryFromReader(message: IdentityProvider, reader: jspb.BinaryReader): IdentityProvider;
}

export namespace IdentityProvider {
  export type AsObject = {
    id: string,
    name: string,
    type: IdentityProviderType,
  }
}

export enum SecondFactorType { 
  SECOND_FACTOR_TYPE_UNSPECIFIED = 0,
  SECOND_FACTOR_TYPE_OTP = 1,
  SECOND_FACTOR_TYPE_U2F = 2,
  SECOND_FACTOR_TYPE_OTP_EMAIL = 3,
  SECOND_FACTOR_TYPE_OTP_SMS = 4,
}
export enum MultiFactorType { 
  MULTI_FACTOR_TYPE_UNSPECIFIED = 0,
  MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION = 1,
}
export enum PasskeysType { 
  PASSKEYS_TYPE_NOT_ALLOWED = 0,
  PASSKEYS_TYPE_ALLOWED = 1,
}
export enum IdentityProviderType { 
  IDENTITY_PROVIDER_TYPE_UNSPECIFIED = 0,
  IDENTITY_PROVIDER_TYPE_OIDC = 1,
  IDENTITY_PROVIDER_TYPE_JWT = 2,
  IDENTITY_PROVIDER_TYPE_LDAP = 3,
  IDENTITY_PROVIDER_TYPE_OAUTH = 4,
  IDENTITY_PROVIDER_TYPE_AZURE_AD = 5,
  IDENTITY_PROVIDER_TYPE_GITHUB = 6,
  IDENTITY_PROVIDER_TYPE_GITHUB_ES = 7,
  IDENTITY_PROVIDER_TYPE_GITLAB = 8,
  IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED = 9,
  IDENTITY_PROVIDER_TYPE_GOOGLE = 10,
  IDENTITY_PROVIDER_TYPE_SAML = 11,
}

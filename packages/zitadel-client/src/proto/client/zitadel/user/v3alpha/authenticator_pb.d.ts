import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"


export class Authenticators extends jspb.Message {
  getUsernamesList(): Array<Username>;
  setUsernamesList(value: Array<Username>): Authenticators;
  clearUsernamesList(): Authenticators;
  addUsernames(value?: Username, index?: number): Username;

  getPassword(): Password | undefined;
  setPassword(value?: Password): Authenticators;
  hasPassword(): boolean;
  clearPassword(): Authenticators;

  getWebAuthNList(): Array<WebAuthN>;
  setWebAuthNList(value: Array<WebAuthN>): Authenticators;
  clearWebAuthNList(): Authenticators;
  addWebAuthN(value?: WebAuthN, index?: number): WebAuthN;

  getTotpsList(): Array<TOTP>;
  setTotpsList(value: Array<TOTP>): Authenticators;
  clearTotpsList(): Authenticators;
  addTotps(value?: TOTP, index?: number): TOTP;

  getOtpSmsList(): Array<OTPSMS>;
  setOtpSmsList(value: Array<OTPSMS>): Authenticators;
  clearOtpSmsList(): Authenticators;
  addOtpSms(value?: OTPSMS, index?: number): OTPSMS;

  getOtpEmailList(): Array<OTPEmail>;
  setOtpEmailList(value: Array<OTPEmail>): Authenticators;
  clearOtpEmailList(): Authenticators;
  addOtpEmail(value?: OTPEmail, index?: number): OTPEmail;

  getAuthenticationKeysList(): Array<AuthenticationKey>;
  setAuthenticationKeysList(value: Array<AuthenticationKey>): Authenticators;
  clearAuthenticationKeysList(): Authenticators;
  addAuthenticationKeys(value?: AuthenticationKey, index?: number): AuthenticationKey;

  getIdentityProvidersList(): Array<IdentityProvider>;
  setIdentityProvidersList(value: Array<IdentityProvider>): Authenticators;
  clearIdentityProvidersList(): Authenticators;
  addIdentityProviders(value?: IdentityProvider, index?: number): IdentityProvider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Authenticators.AsObject;
  static toObject(includeInstance: boolean, msg: Authenticators): Authenticators.AsObject;
  static serializeBinaryToWriter(message: Authenticators, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Authenticators;
  static deserializeBinaryFromReader(message: Authenticators, reader: jspb.BinaryReader): Authenticators;
}

export namespace Authenticators {
  export type AsObject = {
    usernamesList: Array<Username.AsObject>,
    password?: Password.AsObject,
    webAuthNList: Array<WebAuthN.AsObject>,
    totpsList: Array<TOTP.AsObject>,
    otpSmsList: Array<OTPSMS.AsObject>,
    otpEmailList: Array<OTPEmail.AsObject>,
    authenticationKeysList: Array<AuthenticationKey.AsObject>,
    identityProvidersList: Array<IdentityProvider.AsObject>,
  }
}

export class Username extends jspb.Message {
  getUsernameId(): string;
  setUsernameId(value: string): Username;

  getUsername(): string;
  setUsername(value: string): Username;

  getIsOrganizationSpecific(): boolean;
  setIsOrganizationSpecific(value: boolean): Username;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Username.AsObject;
  static toObject(includeInstance: boolean, msg: Username): Username.AsObject;
  static serializeBinaryToWriter(message: Username, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Username;
  static deserializeBinaryFromReader(message: Username, reader: jspb.BinaryReader): Username;
}

export namespace Username {
  export type AsObject = {
    usernameId: string,
    username: string,
    isOrganizationSpecific: boolean,
  }
}

export class SetUsername extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): SetUsername;

  getIsOrganizationSpecific(): boolean;
  setIsOrganizationSpecific(value: boolean): SetUsername;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUsername.AsObject;
  static toObject(includeInstance: boolean, msg: SetUsername): SetUsername.AsObject;
  static serializeBinaryToWriter(message: SetUsername, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUsername;
  static deserializeBinaryFromReader(message: SetUsername, reader: jspb.BinaryReader): SetUsername;
}

export namespace SetUsername {
  export type AsObject = {
    username: string,
    isOrganizationSpecific: boolean,
  }
}

export class Password extends jspb.Message {
  getLastChanged(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastChanged(value?: google_protobuf_timestamp_pb.Timestamp): Password;
  hasLastChanged(): boolean;
  clearLastChanged(): Password;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Password.AsObject;
  static toObject(includeInstance: boolean, msg: Password): Password.AsObject;
  static serializeBinaryToWriter(message: Password, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Password;
  static deserializeBinaryFromReader(message: Password, reader: jspb.BinaryReader): Password;
}

export namespace Password {
  export type AsObject = {
    lastChanged?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class WebAuthN extends jspb.Message {
  getWebAuthNId(): string;
  setWebAuthNId(value: string): WebAuthN;

  getName(): string;
  setName(value: string): WebAuthN;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): WebAuthN;

  getUserVerified(): boolean;
  setUserVerified(value: boolean): WebAuthN;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebAuthN.AsObject;
  static toObject(includeInstance: boolean, msg: WebAuthN): WebAuthN.AsObject;
  static serializeBinaryToWriter(message: WebAuthN, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebAuthN;
  static deserializeBinaryFromReader(message: WebAuthN, reader: jspb.BinaryReader): WebAuthN;
}

export namespace WebAuthN {
  export type AsObject = {
    webAuthNId: string,
    name: string,
    isVerified: boolean,
    userVerified: boolean,
  }
}

export class OTPSMS extends jspb.Message {
  getOtpSmsId(): string;
  setOtpSmsId(value: string): OTPSMS;

  getPhone(): string;
  setPhone(value: string): OTPSMS;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): OTPSMS;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OTPSMS.AsObject;
  static toObject(includeInstance: boolean, msg: OTPSMS): OTPSMS.AsObject;
  static serializeBinaryToWriter(message: OTPSMS, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OTPSMS;
  static deserializeBinaryFromReader(message: OTPSMS, reader: jspb.BinaryReader): OTPSMS;
}

export namespace OTPSMS {
  export type AsObject = {
    otpSmsId: string,
    phone: string,
    isVerified: boolean,
  }
}

export class OTPEmail extends jspb.Message {
  getOtpEmailId(): string;
  setOtpEmailId(value: string): OTPEmail;

  getAddress(): string;
  setAddress(value: string): OTPEmail;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): OTPEmail;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OTPEmail.AsObject;
  static toObject(includeInstance: boolean, msg: OTPEmail): OTPEmail.AsObject;
  static serializeBinaryToWriter(message: OTPEmail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OTPEmail;
  static deserializeBinaryFromReader(message: OTPEmail, reader: jspb.BinaryReader): OTPEmail;
}

export namespace OTPEmail {
  export type AsObject = {
    otpEmailId: string,
    address: string,
    isVerified: boolean,
  }
}

export class TOTP extends jspb.Message {
  getTotpId(): string;
  setTotpId(value: string): TOTP;

  getName(): string;
  setName(value: string): TOTP;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): TOTP;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TOTP.AsObject;
  static toObject(includeInstance: boolean, msg: TOTP): TOTP.AsObject;
  static serializeBinaryToWriter(message: TOTP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TOTP;
  static deserializeBinaryFromReader(message: TOTP, reader: jspb.BinaryReader): TOTP;
}

export namespace TOTP {
  export type AsObject = {
    totpId: string,
    name: string,
    isVerified: boolean,
  }
}

export class AuthenticationKey extends jspb.Message {
  getAuthenticationKeyId(): string;
  setAuthenticationKeyId(value: string): AuthenticationKey;

  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AuthenticationKey;
  hasDetails(): boolean;
  clearDetails(): AuthenticationKey;

  getType(): AuthNKeyType;
  setType(value: AuthNKeyType): AuthenticationKey;

  getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): AuthenticationKey;
  hasExpirationDate(): boolean;
  clearExpirationDate(): AuthenticationKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticationKey.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticationKey): AuthenticationKey.AsObject;
  static serializeBinaryToWriter(message: AuthenticationKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticationKey;
  static deserializeBinaryFromReader(message: AuthenticationKey, reader: jspb.BinaryReader): AuthenticationKey;
}

export namespace AuthenticationKey {
  export type AsObject = {
    authenticationKeyId: string,
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    type: AuthNKeyType,
    expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class IdentityProvider extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): IdentityProvider;

  getIdpName(): string;
  setIdpName(value: string): IdentityProvider;

  getUserId(): string;
  setUserId(value: string): IdentityProvider;

  getUsername(): string;
  setUsername(value: string): IdentityProvider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IdentityProvider.AsObject;
  static toObject(includeInstance: boolean, msg: IdentityProvider): IdentityProvider.AsObject;
  static serializeBinaryToWriter(message: IdentityProvider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IdentityProvider;
  static deserializeBinaryFromReader(message: IdentityProvider, reader: jspb.BinaryReader): IdentityProvider;
}

export namespace IdentityProvider {
  export type AsObject = {
    idpId: string,
    idpName: string,
    userId: string,
    username: string,
  }
}

export class SetAuthenticators extends jspb.Message {
  getUsernamesList(): Array<SetUsername>;
  setUsernamesList(value: Array<SetUsername>): SetAuthenticators;
  clearUsernamesList(): SetAuthenticators;
  addUsernames(value?: SetUsername, index?: number): SetUsername;

  getPassword(): SetPassword | undefined;
  setPassword(value?: SetPassword): SetAuthenticators;
  hasPassword(): boolean;
  clearPassword(): SetAuthenticators;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetAuthenticators.AsObject;
  static toObject(includeInstance: boolean, msg: SetAuthenticators): SetAuthenticators.AsObject;
  static serializeBinaryToWriter(message: SetAuthenticators, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetAuthenticators;
  static deserializeBinaryFromReader(message: SetAuthenticators, reader: jspb.BinaryReader): SetAuthenticators;
}

export namespace SetAuthenticators {
  export type AsObject = {
    usernamesList: Array<SetUsername.AsObject>,
    password?: SetPassword.AsObject,
  }
}

export class SetPassword extends jspb.Message {
  getPassword(): string;
  setPassword(value: string): SetPassword;

  getHash(): string;
  setHash(value: string): SetPassword;

  getChangeRequired(): boolean;
  setChangeRequired(value: boolean): SetPassword;

  getTypeCase(): SetPassword.TypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPassword.AsObject;
  static toObject(includeInstance: boolean, msg: SetPassword): SetPassword.AsObject;
  static serializeBinaryToWriter(message: SetPassword, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPassword;
  static deserializeBinaryFromReader(message: SetPassword, reader: jspb.BinaryReader): SetPassword;
}

export namespace SetPassword {
  export type AsObject = {
    password: string,
    hash: string,
    changeRequired: boolean,
  }

  export enum TypeCase { 
    TYPE_NOT_SET = 0,
    PASSWORD = 1,
    HASH = 2,
  }
}

export class SendPasswordResetEmail extends jspb.Message {
  getUrlTemplate(): string;
  setUrlTemplate(value: string): SendPasswordResetEmail;
  hasUrlTemplate(): boolean;
  clearUrlTemplate(): SendPasswordResetEmail;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendPasswordResetEmail.AsObject;
  static toObject(includeInstance: boolean, msg: SendPasswordResetEmail): SendPasswordResetEmail.AsObject;
  static serializeBinaryToWriter(message: SendPasswordResetEmail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendPasswordResetEmail;
  static deserializeBinaryFromReader(message: SendPasswordResetEmail, reader: jspb.BinaryReader): SendPasswordResetEmail;
}

export namespace SendPasswordResetEmail {
  export type AsObject = {
    urlTemplate?: string,
  }

  export enum UrlTemplateCase { 
    _URL_TEMPLATE_NOT_SET = 0,
    URL_TEMPLATE = 2,
  }
}

export class SendPasswordResetSMS extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendPasswordResetSMS.AsObject;
  static toObject(includeInstance: boolean, msg: SendPasswordResetSMS): SendPasswordResetSMS.AsObject;
  static serializeBinaryToWriter(message: SendPasswordResetSMS, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendPasswordResetSMS;
  static deserializeBinaryFromReader(message: SendPasswordResetSMS, reader: jspb.BinaryReader): SendPasswordResetSMS;
}

export namespace SendPasswordResetSMS {
  export type AsObject = {
  }
}

export class ReturnPasswordResetCode extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReturnPasswordResetCode.AsObject;
  static toObject(includeInstance: boolean, msg: ReturnPasswordResetCode): ReturnPasswordResetCode.AsObject;
  static serializeBinaryToWriter(message: ReturnPasswordResetCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReturnPasswordResetCode;
  static deserializeBinaryFromReader(message: ReturnPasswordResetCode, reader: jspb.BinaryReader): ReturnPasswordResetCode;
}

export namespace ReturnPasswordResetCode {
  export type AsObject = {
  }
}

export class AuthenticatorRegistrationCode extends jspb.Message {
  getId(): string;
  setId(value: string): AuthenticatorRegistrationCode;

  getCode(): string;
  setCode(value: string): AuthenticatorRegistrationCode;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthenticatorRegistrationCode.AsObject;
  static toObject(includeInstance: boolean, msg: AuthenticatorRegistrationCode): AuthenticatorRegistrationCode.AsObject;
  static serializeBinaryToWriter(message: AuthenticatorRegistrationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthenticatorRegistrationCode;
  static deserializeBinaryFromReader(message: AuthenticatorRegistrationCode, reader: jspb.BinaryReader): AuthenticatorRegistrationCode;
}

export namespace AuthenticatorRegistrationCode {
  export type AsObject = {
    id: string,
    code: string,
  }
}

export class SendWebAuthNRegistrationLink extends jspb.Message {
  getUrlTemplate(): string;
  setUrlTemplate(value: string): SendWebAuthNRegistrationLink;
  hasUrlTemplate(): boolean;
  clearUrlTemplate(): SendWebAuthNRegistrationLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendWebAuthNRegistrationLink.AsObject;
  static toObject(includeInstance: boolean, msg: SendWebAuthNRegistrationLink): SendWebAuthNRegistrationLink.AsObject;
  static serializeBinaryToWriter(message: SendWebAuthNRegistrationLink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendWebAuthNRegistrationLink;
  static deserializeBinaryFromReader(message: SendWebAuthNRegistrationLink, reader: jspb.BinaryReader): SendWebAuthNRegistrationLink;
}

export namespace SendWebAuthNRegistrationLink {
  export type AsObject = {
    urlTemplate?: string,
  }

  export enum UrlTemplateCase { 
    _URL_TEMPLATE_NOT_SET = 0,
    URL_TEMPLATE = 1,
  }
}

export class ReturnWebAuthNRegistrationCode extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReturnWebAuthNRegistrationCode.AsObject;
  static toObject(includeInstance: boolean, msg: ReturnWebAuthNRegistrationCode): ReturnWebAuthNRegistrationCode.AsObject;
  static serializeBinaryToWriter(message: ReturnWebAuthNRegistrationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReturnWebAuthNRegistrationCode;
  static deserializeBinaryFromReader(message: ReturnWebAuthNRegistrationCode, reader: jspb.BinaryReader): ReturnWebAuthNRegistrationCode;
}

export namespace ReturnWebAuthNRegistrationCode {
  export type AsObject = {
  }
}

export class RedirectURLs extends jspb.Message {
  getSuccessUrl(): string;
  setSuccessUrl(value: string): RedirectURLs;

  getFailureUrl(): string;
  setFailureUrl(value: string): RedirectURLs;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RedirectURLs.AsObject;
  static toObject(includeInstance: boolean, msg: RedirectURLs): RedirectURLs.AsObject;
  static serializeBinaryToWriter(message: RedirectURLs, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RedirectURLs;
  static deserializeBinaryFromReader(message: RedirectURLs, reader: jspb.BinaryReader): RedirectURLs;
}

export namespace RedirectURLs {
  export type AsObject = {
    successUrl: string,
    failureUrl: string,
  }
}

export class LDAPCredentials extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): LDAPCredentials;

  getPassword(): string;
  setPassword(value: string): LDAPCredentials;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LDAPCredentials.AsObject;
  static toObject(includeInstance: boolean, msg: LDAPCredentials): LDAPCredentials.AsObject;
  static serializeBinaryToWriter(message: LDAPCredentials, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LDAPCredentials;
  static deserializeBinaryFromReader(message: LDAPCredentials, reader: jspb.BinaryReader): LDAPCredentials;
}

export namespace LDAPCredentials {
  export type AsObject = {
    username: string,
    password: string,
  }
}

export class IdentityProviderIntent extends jspb.Message {
  getIdpIntentId(): string;
  setIdpIntentId(value: string): IdentityProviderIntent;

  getIdpIntentToken(): string;
  setIdpIntentToken(value: string): IdentityProviderIntent;

  getUserId(): string;
  setUserId(value: string): IdentityProviderIntent;
  hasUserId(): boolean;
  clearUserId(): IdentityProviderIntent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IdentityProviderIntent.AsObject;
  static toObject(includeInstance: boolean, msg: IdentityProviderIntent): IdentityProviderIntent.AsObject;
  static serializeBinaryToWriter(message: IdentityProviderIntent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IdentityProviderIntent;
  static deserializeBinaryFromReader(message: IdentityProviderIntent, reader: jspb.BinaryReader): IdentityProviderIntent;
}

export namespace IdentityProviderIntent {
  export type AsObject = {
    idpIntentId: string,
    idpIntentToken: string,
    userId?: string,
  }

  export enum UserIdCase { 
    _USER_ID_NOT_SET = 0,
    USER_ID = 3,
  }
}

export class IDPInformation extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): IDPInformation;

  getUserId(): string;
  setUserId(value: string): IDPInformation;

  getUserName(): string;
  setUserName(value: string): IDPInformation;

  getRawInformation(): google_protobuf_struct_pb.Struct | undefined;
  setRawInformation(value?: google_protobuf_struct_pb.Struct): IDPInformation;
  hasRawInformation(): boolean;
  clearRawInformation(): IDPInformation;

  getOauth(): IDPOAuthAccessInformation | undefined;
  setOauth(value?: IDPOAuthAccessInformation): IDPInformation;
  hasOauth(): boolean;
  clearOauth(): IDPInformation;

  getLdap(): IDPLDAPAccessInformation | undefined;
  setLdap(value?: IDPLDAPAccessInformation): IDPInformation;
  hasLdap(): boolean;
  clearLdap(): IDPInformation;

  getSaml(): IDPSAMLAccessInformation | undefined;
  setSaml(value?: IDPSAMLAccessInformation): IDPInformation;
  hasSaml(): boolean;
  clearSaml(): IDPInformation;

  getAccessCase(): IDPInformation.AccessCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPInformation.AsObject;
  static toObject(includeInstance: boolean, msg: IDPInformation): IDPInformation.AsObject;
  static serializeBinaryToWriter(message: IDPInformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPInformation;
  static deserializeBinaryFromReader(message: IDPInformation, reader: jspb.BinaryReader): IDPInformation;
}

export namespace IDPInformation {
  export type AsObject = {
    idpId: string,
    userId: string,
    userName: string,
    rawInformation?: google_protobuf_struct_pb.Struct.AsObject,
    oauth?: IDPOAuthAccessInformation.AsObject,
    ldap?: IDPLDAPAccessInformation.AsObject,
    saml?: IDPSAMLAccessInformation.AsObject,
  }

  export enum AccessCase { 
    ACCESS_NOT_SET = 0,
    OAUTH = 5,
    LDAP = 6,
    SAML = 7,
  }
}

export class IDPOAuthAccessInformation extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): IDPOAuthAccessInformation;

  getIdToken(): string;
  setIdToken(value: string): IDPOAuthAccessInformation;
  hasIdToken(): boolean;
  clearIdToken(): IDPOAuthAccessInformation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPOAuthAccessInformation.AsObject;
  static toObject(includeInstance: boolean, msg: IDPOAuthAccessInformation): IDPOAuthAccessInformation.AsObject;
  static serializeBinaryToWriter(message: IDPOAuthAccessInformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPOAuthAccessInformation;
  static deserializeBinaryFromReader(message: IDPOAuthAccessInformation, reader: jspb.BinaryReader): IDPOAuthAccessInformation;
}

export namespace IDPOAuthAccessInformation {
  export type AsObject = {
    accessToken: string,
    idToken?: string,
  }

  export enum IdTokenCase { 
    _ID_TOKEN_NOT_SET = 0,
    ID_TOKEN = 2,
  }
}

export class IDPLDAPAccessInformation extends jspb.Message {
  getAttributes(): google_protobuf_struct_pb.Struct | undefined;
  setAttributes(value?: google_protobuf_struct_pb.Struct): IDPLDAPAccessInformation;
  hasAttributes(): boolean;
  clearAttributes(): IDPLDAPAccessInformation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPLDAPAccessInformation.AsObject;
  static toObject(includeInstance: boolean, msg: IDPLDAPAccessInformation): IDPLDAPAccessInformation.AsObject;
  static serializeBinaryToWriter(message: IDPLDAPAccessInformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPLDAPAccessInformation;
  static deserializeBinaryFromReader(message: IDPLDAPAccessInformation, reader: jspb.BinaryReader): IDPLDAPAccessInformation;
}

export namespace IDPLDAPAccessInformation {
  export type AsObject = {
    attributes?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class IDPSAMLAccessInformation extends jspb.Message {
  getAssertion(): Uint8Array | string;
  getAssertion_asU8(): Uint8Array;
  getAssertion_asB64(): string;
  setAssertion(value: Uint8Array | string): IDPSAMLAccessInformation;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPSAMLAccessInformation.AsObject;
  static toObject(includeInstance: boolean, msg: IDPSAMLAccessInformation): IDPSAMLAccessInformation.AsObject;
  static serializeBinaryToWriter(message: IDPSAMLAccessInformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPSAMLAccessInformation;
  static deserializeBinaryFromReader(message: IDPSAMLAccessInformation, reader: jspb.BinaryReader): IDPSAMLAccessInformation;
}

export namespace IDPSAMLAccessInformation {
  export type AsObject = {
    assertion: Uint8Array | string,
  }
}

export class IDPAuthenticator extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): IDPAuthenticator;

  getUserId(): string;
  setUserId(value: string): IDPAuthenticator;

  getUserName(): string;
  setUserName(value: string): IDPAuthenticator;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPAuthenticator.AsObject;
  static toObject(includeInstance: boolean, msg: IDPAuthenticator): IDPAuthenticator.AsObject;
  static serializeBinaryToWriter(message: IDPAuthenticator, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPAuthenticator;
  static deserializeBinaryFromReader(message: IDPAuthenticator, reader: jspb.BinaryReader): IDPAuthenticator;
}

export namespace IDPAuthenticator {
  export type AsObject = {
    idpId: string,
    userId: string,
    userName: string,
  }
}

export enum AuthNKeyType { 
  AUTHN_KEY_TYPE_UNSPECIFIED = 0,
  AUTHN_KEY_TYPE_JSON = 1,
}
export enum WebAuthNAuthenticatorType { 
  WEB_AUTH_N_AUTHENTICATOR_UNSPECIFIED = 0,
  WEB_AUTH_N_AUTHENTICATOR_PLATFORM = 1,
  WEB_AUTH_N_AUTHENTICATOR_CROSS_PLATFORM = 2,
}

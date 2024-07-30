import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


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

export class IDPIntent extends jspb.Message {
  getIdpIntentId(): string;
  setIdpIntentId(value: string): IDPIntent;

  getIdpIntentToken(): string;
  setIdpIntentToken(value: string): IDPIntent;

  getUserId(): string;
  setUserId(value: string): IDPIntent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPIntent.AsObject;
  static toObject(includeInstance: boolean, msg: IDPIntent): IDPIntent.AsObject;
  static serializeBinaryToWriter(message: IDPIntent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPIntent;
  static deserializeBinaryFromReader(message: IDPIntent, reader: jspb.BinaryReader): IDPIntent;
}

export namespace IDPIntent {
  export type AsObject = {
    idpIntentId: string,
    idpIntentToken: string,
    userId: string,
  }
}

export class IDPInformation extends jspb.Message {
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
    oauth?: IDPOAuthAccessInformation.AsObject,
    ldap?: IDPLDAPAccessInformation.AsObject,
    saml?: IDPSAMLAccessInformation.AsObject,
    idpId: string,
    userId: string,
    userName: string,
    rawInformation?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export enum AccessCase { 
    ACCESS_NOT_SET = 0,
    OAUTH = 1,
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

export class IDPLink extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): IDPLink;

  getUserId(): string;
  setUserId(value: string): IDPLink;

  getUserName(): string;
  setUserName(value: string): IDPLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPLink.AsObject;
  static toObject(includeInstance: boolean, msg: IDPLink): IDPLink.AsObject;
  static serializeBinaryToWriter(message: IDPLink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPLink;
  static deserializeBinaryFromReader(message: IDPLink, reader: jspb.BinaryReader): IDPLink;
}

export namespace IDPLink {
  export type AsObject = {
    idpId: string,
    userId: string,
    userName: string,
  }
}


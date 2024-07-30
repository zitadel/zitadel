import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class SendPasskeyRegistrationLink extends jspb.Message {
  getUrlTemplate(): string;
  setUrlTemplate(value: string): SendPasskeyRegistrationLink;
  hasUrlTemplate(): boolean;
  clearUrlTemplate(): SendPasskeyRegistrationLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendPasskeyRegistrationLink.AsObject;
  static toObject(includeInstance: boolean, msg: SendPasskeyRegistrationLink): SendPasskeyRegistrationLink.AsObject;
  static serializeBinaryToWriter(message: SendPasskeyRegistrationLink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendPasskeyRegistrationLink;
  static deserializeBinaryFromReader(message: SendPasskeyRegistrationLink, reader: jspb.BinaryReader): SendPasskeyRegistrationLink;
}

export namespace SendPasskeyRegistrationLink {
  export type AsObject = {
    urlTemplate?: string,
  }

  export enum UrlTemplateCase { 
    _URL_TEMPLATE_NOT_SET = 0,
    URL_TEMPLATE = 1,
  }
}

export class ReturnPasskeyRegistrationCode extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReturnPasskeyRegistrationCode.AsObject;
  static toObject(includeInstance: boolean, msg: ReturnPasskeyRegistrationCode): ReturnPasskeyRegistrationCode.AsObject;
  static serializeBinaryToWriter(message: ReturnPasskeyRegistrationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReturnPasskeyRegistrationCode;
  static deserializeBinaryFromReader(message: ReturnPasskeyRegistrationCode, reader: jspb.BinaryReader): ReturnPasskeyRegistrationCode;
}

export namespace ReturnPasskeyRegistrationCode {
  export type AsObject = {
  }
}

export class PasskeyRegistrationCode extends jspb.Message {
  getId(): string;
  setId(value: string): PasskeyRegistrationCode;

  getCode(): string;
  setCode(value: string): PasskeyRegistrationCode;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasskeyRegistrationCode.AsObject;
  static toObject(includeInstance: boolean, msg: PasskeyRegistrationCode): PasskeyRegistrationCode.AsObject;
  static serializeBinaryToWriter(message: PasskeyRegistrationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasskeyRegistrationCode;
  static deserializeBinaryFromReader(message: PasskeyRegistrationCode, reader: jspb.BinaryReader): PasskeyRegistrationCode;
}

export namespace PasskeyRegistrationCode {
  export type AsObject = {
    id: string,
    code: string,
  }
}

export enum PasskeyAuthenticator { 
  PASSKEY_AUTHENTICATOR_UNSPECIFIED = 0,
  PASSKEY_AUTHENTICATOR_PLATFORM = 1,
  PASSKEY_AUTHENTICATOR_CROSS_PLATFORM = 2,
}

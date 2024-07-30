import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class SetHumanEmail extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): SetHumanEmail;

  getSendCode(): SendEmailVerificationCode | undefined;
  setSendCode(value?: SendEmailVerificationCode): SetHumanEmail;
  hasSendCode(): boolean;
  clearSendCode(): SetHumanEmail;

  getReturnCode(): ReturnEmailVerificationCode | undefined;
  setReturnCode(value?: ReturnEmailVerificationCode): SetHumanEmail;
  hasReturnCode(): boolean;
  clearReturnCode(): SetHumanEmail;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): SetHumanEmail;

  getVerificationCase(): SetHumanEmail.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetHumanEmail.AsObject;
  static toObject(includeInstance: boolean, msg: SetHumanEmail): SetHumanEmail.AsObject;
  static serializeBinaryToWriter(message: SetHumanEmail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetHumanEmail;
  static deserializeBinaryFromReader(message: SetHumanEmail, reader: jspb.BinaryReader): SetHumanEmail;
}

export namespace SetHumanEmail {
  export type AsObject = {
    email: string,
    sendCode?: SendEmailVerificationCode.AsObject,
    returnCode?: ReturnEmailVerificationCode.AsObject,
    isVerified: boolean,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    SEND_CODE = 2,
    RETURN_CODE = 3,
    IS_VERIFIED = 4,
  }
}

export class HumanEmail extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): HumanEmail;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): HumanEmail;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HumanEmail.AsObject;
  static toObject(includeInstance: boolean, msg: HumanEmail): HumanEmail.AsObject;
  static serializeBinaryToWriter(message: HumanEmail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HumanEmail;
  static deserializeBinaryFromReader(message: HumanEmail, reader: jspb.BinaryReader): HumanEmail;
}

export namespace HumanEmail {
  export type AsObject = {
    email: string,
    isVerified: boolean,
  }
}

export class SendEmailVerificationCode extends jspb.Message {
  getUrlTemplate(): string;
  setUrlTemplate(value: string): SendEmailVerificationCode;
  hasUrlTemplate(): boolean;
  clearUrlTemplate(): SendEmailVerificationCode;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendEmailVerificationCode.AsObject;
  static toObject(includeInstance: boolean, msg: SendEmailVerificationCode): SendEmailVerificationCode.AsObject;
  static serializeBinaryToWriter(message: SendEmailVerificationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendEmailVerificationCode;
  static deserializeBinaryFromReader(message: SendEmailVerificationCode, reader: jspb.BinaryReader): SendEmailVerificationCode;
}

export namespace SendEmailVerificationCode {
  export type AsObject = {
    urlTemplate?: string,
  }

  export enum UrlTemplateCase { 
    _URL_TEMPLATE_NOT_SET = 0,
    URL_TEMPLATE = 1,
  }
}

export class ReturnEmailVerificationCode extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReturnEmailVerificationCode.AsObject;
  static toObject(includeInstance: boolean, msg: ReturnEmailVerificationCode): ReturnEmailVerificationCode.AsObject;
  static serializeBinaryToWriter(message: ReturnEmailVerificationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReturnEmailVerificationCode;
  static deserializeBinaryFromReader(message: ReturnEmailVerificationCode, reader: jspb.BinaryReader): ReturnEmailVerificationCode;
}

export namespace ReturnEmailVerificationCode {
  export type AsObject = {
  }
}


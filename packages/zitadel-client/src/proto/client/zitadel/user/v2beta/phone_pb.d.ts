import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class SetHumanPhone extends jspb.Message {
  getPhone(): string;
  setPhone(value: string): SetHumanPhone;

  getSendCode(): SendPhoneVerificationCode | undefined;
  setSendCode(value?: SendPhoneVerificationCode): SetHumanPhone;
  hasSendCode(): boolean;
  clearSendCode(): SetHumanPhone;

  getReturnCode(): ReturnPhoneVerificationCode | undefined;
  setReturnCode(value?: ReturnPhoneVerificationCode): SetHumanPhone;
  hasReturnCode(): boolean;
  clearReturnCode(): SetHumanPhone;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): SetHumanPhone;

  getVerificationCase(): SetHumanPhone.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetHumanPhone.AsObject;
  static toObject(includeInstance: boolean, msg: SetHumanPhone): SetHumanPhone.AsObject;
  static serializeBinaryToWriter(message: SetHumanPhone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetHumanPhone;
  static deserializeBinaryFromReader(message: SetHumanPhone, reader: jspb.BinaryReader): SetHumanPhone;
}

export namespace SetHumanPhone {
  export type AsObject = {
    phone: string,
    sendCode?: SendPhoneVerificationCode.AsObject,
    returnCode?: ReturnPhoneVerificationCode.AsObject,
    isVerified: boolean,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    SEND_CODE = 2,
    RETURN_CODE = 3,
    IS_VERIFIED = 4,
  }
}

export class HumanPhone extends jspb.Message {
  getPhone(): string;
  setPhone(value: string): HumanPhone;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): HumanPhone;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HumanPhone.AsObject;
  static toObject(includeInstance: boolean, msg: HumanPhone): HumanPhone.AsObject;
  static serializeBinaryToWriter(message: HumanPhone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HumanPhone;
  static deserializeBinaryFromReader(message: HumanPhone, reader: jspb.BinaryReader): HumanPhone;
}

export namespace HumanPhone {
  export type AsObject = {
    phone: string,
    isVerified: boolean,
  }
}

export class SendPhoneVerificationCode extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendPhoneVerificationCode.AsObject;
  static toObject(includeInstance: boolean, msg: SendPhoneVerificationCode): SendPhoneVerificationCode.AsObject;
  static serializeBinaryToWriter(message: SendPhoneVerificationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendPhoneVerificationCode;
  static deserializeBinaryFromReader(message: SendPhoneVerificationCode, reader: jspb.BinaryReader): SendPhoneVerificationCode;
}

export namespace SendPhoneVerificationCode {
  export type AsObject = {
  }
}

export class ReturnPhoneVerificationCode extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReturnPhoneVerificationCode.AsObject;
  static toObject(includeInstance: boolean, msg: ReturnPhoneVerificationCode): ReturnPhoneVerificationCode.AsObject;
  static serializeBinaryToWriter(message: ReturnPhoneVerificationCode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReturnPhoneVerificationCode;
  static deserializeBinaryFromReader(message: ReturnPhoneVerificationCode, reader: jspb.BinaryReader): ReturnPhoneVerificationCode;
}

export namespace ReturnPhoneVerificationCode {
  export type AsObject = {
  }
}


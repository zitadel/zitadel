import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class RequestChallenges extends jspb.Message {
  getWebAuthN(): RequestChallenges.WebAuthN | undefined;
  setWebAuthN(value?: RequestChallenges.WebAuthN): RequestChallenges;
  hasWebAuthN(): boolean;
  clearWebAuthN(): RequestChallenges;

  getOtpSms(): RequestChallenges.OTPSMS | undefined;
  setOtpSms(value?: RequestChallenges.OTPSMS): RequestChallenges;
  hasOtpSms(): boolean;
  clearOtpSms(): RequestChallenges;

  getOtpEmail(): RequestChallenges.OTPEmail | undefined;
  setOtpEmail(value?: RequestChallenges.OTPEmail): RequestChallenges;
  hasOtpEmail(): boolean;
  clearOtpEmail(): RequestChallenges;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestChallenges.AsObject;
  static toObject(includeInstance: boolean, msg: RequestChallenges): RequestChallenges.AsObject;
  static serializeBinaryToWriter(message: RequestChallenges, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestChallenges;
  static deserializeBinaryFromReader(message: RequestChallenges, reader: jspb.BinaryReader): RequestChallenges;
}

export namespace RequestChallenges {
  export type AsObject = {
    webAuthN?: RequestChallenges.WebAuthN.AsObject,
    otpSms?: RequestChallenges.OTPSMS.AsObject,
    otpEmail?: RequestChallenges.OTPEmail.AsObject,
  }

  export class WebAuthN extends jspb.Message {
    getDomain(): string;
    setDomain(value: string): WebAuthN;

    getUserVerificationRequirement(): UserVerificationRequirement;
    setUserVerificationRequirement(value: UserVerificationRequirement): WebAuthN;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): WebAuthN.AsObject;
    static toObject(includeInstance: boolean, msg: WebAuthN): WebAuthN.AsObject;
    static serializeBinaryToWriter(message: WebAuthN, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): WebAuthN;
    static deserializeBinaryFromReader(message: WebAuthN, reader: jspb.BinaryReader): WebAuthN;
  }

  export namespace WebAuthN {
    export type AsObject = {
      domain: string,
      userVerificationRequirement: UserVerificationRequirement,
    }
  }


  export class OTPSMS extends jspb.Message {
    getReturnCode(): boolean;
    setReturnCode(value: boolean): OTPSMS;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OTPSMS.AsObject;
    static toObject(includeInstance: boolean, msg: OTPSMS): OTPSMS.AsObject;
    static serializeBinaryToWriter(message: OTPSMS, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OTPSMS;
    static deserializeBinaryFromReader(message: OTPSMS, reader: jspb.BinaryReader): OTPSMS;
  }

  export namespace OTPSMS {
    export type AsObject = {
      returnCode: boolean,
    }
  }


  export class OTPEmail extends jspb.Message {
    getSendCode(): RequestChallenges.OTPEmail.SendCode | undefined;
    setSendCode(value?: RequestChallenges.OTPEmail.SendCode): OTPEmail;
    hasSendCode(): boolean;
    clearSendCode(): OTPEmail;

    getReturnCode(): RequestChallenges.OTPEmail.ReturnCode | undefined;
    setReturnCode(value?: RequestChallenges.OTPEmail.ReturnCode): OTPEmail;
    hasReturnCode(): boolean;
    clearReturnCode(): OTPEmail;

    getDeliveryTypeCase(): OTPEmail.DeliveryTypeCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OTPEmail.AsObject;
    static toObject(includeInstance: boolean, msg: OTPEmail): OTPEmail.AsObject;
    static serializeBinaryToWriter(message: OTPEmail, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OTPEmail;
    static deserializeBinaryFromReader(message: OTPEmail, reader: jspb.BinaryReader): OTPEmail;
  }

  export namespace OTPEmail {
    export type AsObject = {
      sendCode?: RequestChallenges.OTPEmail.SendCode.AsObject,
      returnCode?: RequestChallenges.OTPEmail.ReturnCode.AsObject,
    }

    export class SendCode extends jspb.Message {
      getUrlTemplate(): string;
      setUrlTemplate(value: string): SendCode;
      hasUrlTemplate(): boolean;
      clearUrlTemplate(): SendCode;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): SendCode.AsObject;
      static toObject(includeInstance: boolean, msg: SendCode): SendCode.AsObject;
      static serializeBinaryToWriter(message: SendCode, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): SendCode;
      static deserializeBinaryFromReader(message: SendCode, reader: jspb.BinaryReader): SendCode;
    }

    export namespace SendCode {
      export type AsObject = {
        urlTemplate?: string,
      }

      export enum UrlTemplateCase { 
        _URL_TEMPLATE_NOT_SET = 0,
        URL_TEMPLATE = 1,
      }
    }


    export class ReturnCode extends jspb.Message {
      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ReturnCode.AsObject;
      static toObject(includeInstance: boolean, msg: ReturnCode): ReturnCode.AsObject;
      static serializeBinaryToWriter(message: ReturnCode, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ReturnCode;
      static deserializeBinaryFromReader(message: ReturnCode, reader: jspb.BinaryReader): ReturnCode;
    }

    export namespace ReturnCode {
      export type AsObject = {
      }
    }


    export enum DeliveryTypeCase { 
      DELIVERY_TYPE_NOT_SET = 0,
      SEND_CODE = 2,
      RETURN_CODE = 3,
    }
  }


  export enum WebAuthNCase { 
    _WEB_AUTH_N_NOT_SET = 0,
    WEB_AUTH_N = 1,
  }

  export enum OtpSmsCase { 
    _OTP_SMS_NOT_SET = 0,
    OTP_SMS = 2,
  }

  export enum OtpEmailCase { 
    _OTP_EMAIL_NOT_SET = 0,
    OTP_EMAIL = 3,
  }
}

export class Challenges extends jspb.Message {
  getWebAuthN(): Challenges.WebAuthN | undefined;
  setWebAuthN(value?: Challenges.WebAuthN): Challenges;
  hasWebAuthN(): boolean;
  clearWebAuthN(): Challenges;

  getOtpSms(): string;
  setOtpSms(value: string): Challenges;
  hasOtpSms(): boolean;
  clearOtpSms(): Challenges;

  getOtpEmail(): string;
  setOtpEmail(value: string): Challenges;
  hasOtpEmail(): boolean;
  clearOtpEmail(): Challenges;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Challenges.AsObject;
  static toObject(includeInstance: boolean, msg: Challenges): Challenges.AsObject;
  static serializeBinaryToWriter(message: Challenges, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Challenges;
  static deserializeBinaryFromReader(message: Challenges, reader: jspb.BinaryReader): Challenges;
}

export namespace Challenges {
  export type AsObject = {
    webAuthN?: Challenges.WebAuthN.AsObject,
    otpSms?: string,
    otpEmail?: string,
  }

  export class WebAuthN extends jspb.Message {
    getPublicKeyCredentialRequestOptions(): google_protobuf_struct_pb.Struct | undefined;
    setPublicKeyCredentialRequestOptions(value?: google_protobuf_struct_pb.Struct): WebAuthN;
    hasPublicKeyCredentialRequestOptions(): boolean;
    clearPublicKeyCredentialRequestOptions(): WebAuthN;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): WebAuthN.AsObject;
    static toObject(includeInstance: boolean, msg: WebAuthN): WebAuthN.AsObject;
    static serializeBinaryToWriter(message: WebAuthN, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): WebAuthN;
    static deserializeBinaryFromReader(message: WebAuthN, reader: jspb.BinaryReader): WebAuthN;
  }

  export namespace WebAuthN {
    export type AsObject = {
      publicKeyCredentialRequestOptions?: google_protobuf_struct_pb.Struct.AsObject,
    }
  }


  export enum WebAuthNCase { 
    _WEB_AUTH_N_NOT_SET = 0,
    WEB_AUTH_N = 1,
  }

  export enum OtpSmsCase { 
    _OTP_SMS_NOT_SET = 0,
    OTP_SMS = 2,
  }

  export enum OtpEmailCase { 
    _OTP_EMAIL_NOT_SET = 0,
    OTP_EMAIL = 3,
  }
}

export enum UserVerificationRequirement { 
  USER_VERIFICATION_REQUIREMENT_UNSPECIFIED = 0,
  USER_VERIFICATION_REQUIREMENT_REQUIRED = 1,
  USER_VERIFICATION_REQUIREMENT_PREFERRED = 2,
  USER_VERIFICATION_REQUIREMENT_DISCOURAGED = 3,
}

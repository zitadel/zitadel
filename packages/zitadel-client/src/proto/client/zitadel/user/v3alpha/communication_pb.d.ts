import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class Contact extends jspb.Message {
  getEmail(): Email | undefined;
  setEmail(value?: Email): Contact;
  hasEmail(): boolean;
  clearEmail(): Contact;

  getPhone(): Phone | undefined;
  setPhone(value?: Phone): Contact;
  hasPhone(): boolean;
  clearPhone(): Contact;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Contact.AsObject;
  static toObject(includeInstance: boolean, msg: Contact): Contact.AsObject;
  static serializeBinaryToWriter(message: Contact, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Contact;
  static deserializeBinaryFromReader(message: Contact, reader: jspb.BinaryReader): Contact;
}

export namespace Contact {
  export type AsObject = {
    email?: Email.AsObject,
    phone?: Phone.AsObject,
  }
}

export class Email extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): Email;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): Email;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Email.AsObject;
  static toObject(includeInstance: boolean, msg: Email): Email.AsObject;
  static serializeBinaryToWriter(message: Email, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Email;
  static deserializeBinaryFromReader(message: Email, reader: jspb.BinaryReader): Email;
}

export namespace Email {
  export type AsObject = {
    address: string,
    isVerified: boolean,
  }
}

export class Phone extends jspb.Message {
  getNumber(): string;
  setNumber(value: string): Phone;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): Phone;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Phone.AsObject;
  static toObject(includeInstance: boolean, msg: Phone): Phone.AsObject;
  static serializeBinaryToWriter(message: Phone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Phone;
  static deserializeBinaryFromReader(message: Phone, reader: jspb.BinaryReader): Phone;
}

export namespace Phone {
  export type AsObject = {
    number: string,
    isVerified: boolean,
  }
}

export class SetContact extends jspb.Message {
  getEmail(): SetEmail | undefined;
  setEmail(value?: SetEmail): SetContact;
  hasEmail(): boolean;
  clearEmail(): SetContact;

  getPhone(): SetPhone | undefined;
  setPhone(value?: SetPhone): SetContact;
  hasPhone(): boolean;
  clearPhone(): SetContact;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetContact.AsObject;
  static toObject(includeInstance: boolean, msg: SetContact): SetContact.AsObject;
  static serializeBinaryToWriter(message: SetContact, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetContact;
  static deserializeBinaryFromReader(message: SetContact, reader: jspb.BinaryReader): SetContact;
}

export namespace SetContact {
  export type AsObject = {
    email?: SetEmail.AsObject,
    phone?: SetPhone.AsObject,
  }

  export enum EmailCase { 
    _EMAIL_NOT_SET = 0,
    EMAIL = 1,
  }

  export enum PhoneCase { 
    _PHONE_NOT_SET = 0,
    PHONE = 2,
  }
}

export class SetEmail extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): SetEmail;

  getSendCode(): SendEmailVerificationCode | undefined;
  setSendCode(value?: SendEmailVerificationCode): SetEmail;
  hasSendCode(): boolean;
  clearSendCode(): SetEmail;

  getReturnCode(): ReturnEmailVerificationCode | undefined;
  setReturnCode(value?: ReturnEmailVerificationCode): SetEmail;
  hasReturnCode(): boolean;
  clearReturnCode(): SetEmail;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): SetEmail;

  getVerificationCase(): SetEmail.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetEmail.AsObject;
  static toObject(includeInstance: boolean, msg: SetEmail): SetEmail.AsObject;
  static serializeBinaryToWriter(message: SetEmail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetEmail;
  static deserializeBinaryFromReader(message: SetEmail, reader: jspb.BinaryReader): SetEmail;
}

export namespace SetEmail {
  export type AsObject = {
    address: string,
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

export class SetPhone extends jspb.Message {
  getNumber(): string;
  setNumber(value: string): SetPhone;

  getSendCode(): SendPhoneVerificationCode | undefined;
  setSendCode(value?: SendPhoneVerificationCode): SetPhone;
  hasSendCode(): boolean;
  clearSendCode(): SetPhone;

  getReturnCode(): ReturnPhoneVerificationCode | undefined;
  setReturnCode(value?: ReturnPhoneVerificationCode): SetPhone;
  hasReturnCode(): boolean;
  clearReturnCode(): SetPhone;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): SetPhone;

  getVerificationCase(): SetPhone.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPhone.AsObject;
  static toObject(includeInstance: boolean, msg: SetPhone): SetPhone.AsObject;
  static serializeBinaryToWriter(message: SetPhone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPhone;
  static deserializeBinaryFromReader(message: SetPhone, reader: jspb.BinaryReader): SetPhone;
}

export namespace SetPhone {
  export type AsObject = {
    number: string,
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


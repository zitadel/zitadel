import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class Password extends jspb.Message {
  getPassword(): string;
  setPassword(value: string): Password;

  getChangeRequired(): boolean;
  setChangeRequired(value: boolean): Password;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Password.AsObject;
  static toObject(includeInstance: boolean, msg: Password): Password.AsObject;
  static serializeBinaryToWriter(message: Password, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Password;
  static deserializeBinaryFromReader(message: Password, reader: jspb.BinaryReader): Password;
}

export namespace Password {
  export type AsObject = {
    password: string,
    changeRequired: boolean,
  }
}

export class HashedPassword extends jspb.Message {
  getHash(): string;
  setHash(value: string): HashedPassword;

  getChangeRequired(): boolean;
  setChangeRequired(value: boolean): HashedPassword;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HashedPassword.AsObject;
  static toObject(includeInstance: boolean, msg: HashedPassword): HashedPassword.AsObject;
  static serializeBinaryToWriter(message: HashedPassword, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HashedPassword;
  static deserializeBinaryFromReader(message: HashedPassword, reader: jspb.BinaryReader): HashedPassword;
}

export namespace HashedPassword {
  export type AsObject = {
    hash: string,
    changeRequired: boolean,
  }
}

export class SendPasswordResetLink extends jspb.Message {
  getNotificationType(): NotificationType;
  setNotificationType(value: NotificationType): SendPasswordResetLink;

  getUrlTemplate(): string;
  setUrlTemplate(value: string): SendPasswordResetLink;
  hasUrlTemplate(): boolean;
  clearUrlTemplate(): SendPasswordResetLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendPasswordResetLink.AsObject;
  static toObject(includeInstance: boolean, msg: SendPasswordResetLink): SendPasswordResetLink.AsObject;
  static serializeBinaryToWriter(message: SendPasswordResetLink, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendPasswordResetLink;
  static deserializeBinaryFromReader(message: SendPasswordResetLink, reader: jspb.BinaryReader): SendPasswordResetLink;
}

export namespace SendPasswordResetLink {
  export type AsObject = {
    notificationType: NotificationType,
    urlTemplate?: string,
  }

  export enum UrlTemplateCase { 
    _URL_TEMPLATE_NOT_SET = 0,
    URL_TEMPLATE = 2,
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

export class SetPassword extends jspb.Message {
  getPassword(): Password | undefined;
  setPassword(value?: Password): SetPassword;
  hasPassword(): boolean;
  clearPassword(): SetPassword;

  getHashedPassword(): HashedPassword | undefined;
  setHashedPassword(value?: HashedPassword): SetPassword;
  hasHashedPassword(): boolean;
  clearHashedPassword(): SetPassword;

  getCurrentPassword(): string;
  setCurrentPassword(value: string): SetPassword;

  getVerificationCode(): string;
  setVerificationCode(value: string): SetPassword;

  getPasswordTypeCase(): SetPassword.PasswordTypeCase;

  getVerificationCase(): SetPassword.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPassword.AsObject;
  static toObject(includeInstance: boolean, msg: SetPassword): SetPassword.AsObject;
  static serializeBinaryToWriter(message: SetPassword, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPassword;
  static deserializeBinaryFromReader(message: SetPassword, reader: jspb.BinaryReader): SetPassword;
}

export namespace SetPassword {
  export type AsObject = {
    password?: Password.AsObject,
    hashedPassword?: HashedPassword.AsObject,
    currentPassword: string,
    verificationCode: string,
  }

  export enum PasswordTypeCase { 
    PASSWORD_TYPE_NOT_SET = 0,
    PASSWORD = 1,
    HASHED_PASSWORD = 2,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    CURRENT_PASSWORD = 3,
    VERIFICATION_CODE = 4,
  }
}

export enum NotificationType { 
  NOTIFICATION_TYPE_UNSPECIFIED = 0,
  NOTIFICATION_TYPE_EMAIL = 1,
  NOTIFICATION_TYPE_SMS = 2,
}

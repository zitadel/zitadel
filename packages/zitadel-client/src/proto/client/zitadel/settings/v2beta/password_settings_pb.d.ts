import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as zitadel_settings_v2beta_settings_pb from '../../../zitadel/settings/v2beta/settings_pb'; // proto import: "zitadel/settings/v2beta/settings.proto"


export class PasswordComplexitySettings extends jspb.Message {
  getMinLength(): number;
  setMinLength(value: number): PasswordComplexitySettings;

  getRequiresUppercase(): boolean;
  setRequiresUppercase(value: boolean): PasswordComplexitySettings;

  getRequiresLowercase(): boolean;
  setRequiresLowercase(value: boolean): PasswordComplexitySettings;

  getRequiresNumber(): boolean;
  setRequiresNumber(value: boolean): PasswordComplexitySettings;

  getRequiresSymbol(): boolean;
  setRequiresSymbol(value: boolean): PasswordComplexitySettings;

  getResourceOwnerType(): zitadel_settings_v2beta_settings_pb.ResourceOwnerType;
  setResourceOwnerType(value: zitadel_settings_v2beta_settings_pb.ResourceOwnerType): PasswordComplexitySettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordComplexitySettings.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordComplexitySettings): PasswordComplexitySettings.AsObject;
  static serializeBinaryToWriter(message: PasswordComplexitySettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordComplexitySettings;
  static deserializeBinaryFromReader(message: PasswordComplexitySettings, reader: jspb.BinaryReader): PasswordComplexitySettings;
}

export namespace PasswordComplexitySettings {
  export type AsObject = {
    minLength: number,
    requiresUppercase: boolean,
    requiresLowercase: boolean,
    requiresNumber: boolean,
    requiresSymbol: boolean,
    resourceOwnerType: zitadel_settings_v2beta_settings_pb.ResourceOwnerType,
  }
}

export class PasswordExpirySettings extends jspb.Message {
  getMaxAgeDays(): number;
  setMaxAgeDays(value: number): PasswordExpirySettings;

  getExpireWarnDays(): number;
  setExpireWarnDays(value: number): PasswordExpirySettings;

  getResourceOwnerType(): zitadel_settings_v2beta_settings_pb.ResourceOwnerType;
  setResourceOwnerType(value: zitadel_settings_v2beta_settings_pb.ResourceOwnerType): PasswordExpirySettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordExpirySettings.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordExpirySettings): PasswordExpirySettings.AsObject;
  static serializeBinaryToWriter(message: PasswordExpirySettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordExpirySettings;
  static deserializeBinaryFromReader(message: PasswordExpirySettings, reader: jspb.BinaryReader): PasswordExpirySettings;
}

export namespace PasswordExpirySettings {
  export type AsObject = {
    maxAgeDays: number,
    expireWarnDays: number,
    resourceOwnerType: zitadel_settings_v2beta_settings_pb.ResourceOwnerType,
  }
}


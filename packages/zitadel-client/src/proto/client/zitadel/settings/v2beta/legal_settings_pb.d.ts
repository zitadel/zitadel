import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as zitadel_settings_v2beta_settings_pb from '../../../zitadel/settings/v2beta/settings_pb'; // proto import: "zitadel/settings/v2beta/settings.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class LegalAndSupportSettings extends jspb.Message {
  getTosLink(): string;
  setTosLink(value: string): LegalAndSupportSettings;

  getPrivacyPolicyLink(): string;
  setPrivacyPolicyLink(value: string): LegalAndSupportSettings;

  getHelpLink(): string;
  setHelpLink(value: string): LegalAndSupportSettings;

  getSupportEmail(): string;
  setSupportEmail(value: string): LegalAndSupportSettings;

  getResourceOwnerType(): zitadel_settings_v2beta_settings_pb.ResourceOwnerType;
  setResourceOwnerType(value: zitadel_settings_v2beta_settings_pb.ResourceOwnerType): LegalAndSupportSettings;

  getDocsLink(): string;
  setDocsLink(value: string): LegalAndSupportSettings;

  getCustomLink(): string;
  setCustomLink(value: string): LegalAndSupportSettings;

  getCustomLinkText(): string;
  setCustomLinkText(value: string): LegalAndSupportSettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LegalAndSupportSettings.AsObject;
  static toObject(includeInstance: boolean, msg: LegalAndSupportSettings): LegalAndSupportSettings.AsObject;
  static serializeBinaryToWriter(message: LegalAndSupportSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LegalAndSupportSettings;
  static deserializeBinaryFromReader(message: LegalAndSupportSettings, reader: jspb.BinaryReader): LegalAndSupportSettings;
}

export namespace LegalAndSupportSettings {
  export type AsObject = {
    tosLink: string,
    privacyPolicyLink: string,
    helpLink: string,
    supportEmail: string,
    resourceOwnerType: zitadel_settings_v2beta_settings_pb.ResourceOwnerType,
    docsLink: string,
    customLink: string,
    customLinkText: string,
  }
}


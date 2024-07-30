import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as zitadel_settings_v2beta_settings_pb from '../../../zitadel/settings/v2beta/settings_pb'; // proto import: "zitadel/settings/v2beta/settings.proto"


export class DomainSettings extends jspb.Message {
  getLoginNameIncludesDomain(): boolean;
  setLoginNameIncludesDomain(value: boolean): DomainSettings;

  getRequireOrgDomainVerification(): boolean;
  setRequireOrgDomainVerification(value: boolean): DomainSettings;

  getSmtpSenderAddressMatchesInstanceDomain(): boolean;
  setSmtpSenderAddressMatchesInstanceDomain(value: boolean): DomainSettings;

  getResourceOwnerType(): zitadel_settings_v2beta_settings_pb.ResourceOwnerType;
  setResourceOwnerType(value: zitadel_settings_v2beta_settings_pb.ResourceOwnerType): DomainSettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainSettings.AsObject;
  static toObject(includeInstance: boolean, msg: DomainSettings): DomainSettings.AsObject;
  static serializeBinaryToWriter(message: DomainSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainSettings;
  static deserializeBinaryFromReader(message: DomainSettings, reader: jspb.BinaryReader): DomainSettings;
}

export namespace DomainSettings {
  export type AsObject = {
    loginNameIncludesDomain: boolean,
    requireOrgDomainVerification: boolean,
    smtpSenderAddressMatchesInstanceDomain: boolean,
    resourceOwnerType: zitadel_settings_v2beta_settings_pb.ResourceOwnerType,
  }
}


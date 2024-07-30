import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_feature_v2beta_feature_pb from '../../../zitadel/feature/v2beta/feature_pb'; // proto import: "zitadel/feature/v2beta/feature.proto"


export class SetOrganizationFeaturesRequest extends jspb.Message {
  getOrganizationId(): string;
  setOrganizationId(value: string): SetOrganizationFeaturesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetOrganizationFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetOrganizationFeaturesRequest): SetOrganizationFeaturesRequest.AsObject;
  static serializeBinaryToWriter(message: SetOrganizationFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetOrganizationFeaturesRequest;
  static deserializeBinaryFromReader(message: SetOrganizationFeaturesRequest, reader: jspb.BinaryReader): SetOrganizationFeaturesRequest;
}

export namespace SetOrganizationFeaturesRequest {
  export type AsObject = {
    organizationId: string,
  }
}

export class SetOrganizationFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetOrganizationFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): SetOrganizationFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetOrganizationFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetOrganizationFeaturesResponse): SetOrganizationFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: SetOrganizationFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetOrganizationFeaturesResponse;
  static deserializeBinaryFromReader(message: SetOrganizationFeaturesResponse, reader: jspb.BinaryReader): SetOrganizationFeaturesResponse;
}

export namespace SetOrganizationFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ResetOrganizationFeaturesRequest extends jspb.Message {
  getOrganizationId(): string;
  setOrganizationId(value: string): ResetOrganizationFeaturesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetOrganizationFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetOrganizationFeaturesRequest): ResetOrganizationFeaturesRequest.AsObject;
  static serializeBinaryToWriter(message: ResetOrganizationFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetOrganizationFeaturesRequest;
  static deserializeBinaryFromReader(message: ResetOrganizationFeaturesRequest, reader: jspb.BinaryReader): ResetOrganizationFeaturesRequest;
}

export namespace ResetOrganizationFeaturesRequest {
  export type AsObject = {
    organizationId: string,
  }
}

export class ResetOrganizationFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ResetOrganizationFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): ResetOrganizationFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetOrganizationFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetOrganizationFeaturesResponse): ResetOrganizationFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: ResetOrganizationFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetOrganizationFeaturesResponse;
  static deserializeBinaryFromReader(message: ResetOrganizationFeaturesResponse, reader: jspb.BinaryReader): ResetOrganizationFeaturesResponse;
}

export namespace ResetOrganizationFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class GetOrganizationFeaturesRequest extends jspb.Message {
  getOrganizationId(): string;
  setOrganizationId(value: string): GetOrganizationFeaturesRequest;

  getInheritance(): boolean;
  setInheritance(value: boolean): GetOrganizationFeaturesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationFeaturesRequest): GetOrganizationFeaturesRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrganizationFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationFeaturesRequest;
  static deserializeBinaryFromReader(message: GetOrganizationFeaturesRequest, reader: jspb.BinaryReader): GetOrganizationFeaturesRequest;
}

export namespace GetOrganizationFeaturesRequest {
  export type AsObject = {
    organizationId: string,
    inheritance: boolean,
  }
}

export class GetOrganizationFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): GetOrganizationFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): GetOrganizationFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrganizationFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrganizationFeaturesResponse): GetOrganizationFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrganizationFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrganizationFeaturesResponse;
  static deserializeBinaryFromReader(message: GetOrganizationFeaturesResponse, reader: jspb.BinaryReader): GetOrganizationFeaturesResponse;
}

export namespace GetOrganizationFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}


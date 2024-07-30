import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_feature_v2beta_feature_pb from '../../../zitadel/feature/v2beta/feature_pb'; // proto import: "zitadel/feature/v2beta/feature.proto"


export class SetUserFeatureRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetUserFeatureRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUserFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetUserFeatureRequest): SetUserFeatureRequest.AsObject;
  static serializeBinaryToWriter(message: SetUserFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUserFeatureRequest;
  static deserializeBinaryFromReader(message: SetUserFeatureRequest, reader: jspb.BinaryReader): SetUserFeatureRequest;
}

export namespace SetUserFeatureRequest {
  export type AsObject = {
    userId: string,
  }
}

export class SetUserFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetUserFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): SetUserFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUserFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetUserFeaturesResponse): SetUserFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: SetUserFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUserFeaturesResponse;
  static deserializeBinaryFromReader(message: SetUserFeaturesResponse, reader: jspb.BinaryReader): SetUserFeaturesResponse;
}

export namespace SetUserFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ResetUserFeaturesRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ResetUserFeaturesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetUserFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetUserFeaturesRequest): ResetUserFeaturesRequest.AsObject;
  static serializeBinaryToWriter(message: ResetUserFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetUserFeaturesRequest;
  static deserializeBinaryFromReader(message: ResetUserFeaturesRequest, reader: jspb.BinaryReader): ResetUserFeaturesRequest;
}

export namespace ResetUserFeaturesRequest {
  export type AsObject = {
    userId: string,
  }
}

export class ResetUserFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ResetUserFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): ResetUserFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetUserFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetUserFeaturesResponse): ResetUserFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: ResetUserFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetUserFeaturesResponse;
  static deserializeBinaryFromReader(message: ResetUserFeaturesResponse, reader: jspb.BinaryReader): ResetUserFeaturesResponse;
}

export namespace ResetUserFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class GetUserFeaturesRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GetUserFeaturesRequest;

  getInheritance(): boolean;
  setInheritance(value: boolean): GetUserFeaturesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserFeaturesRequest): GetUserFeaturesRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserFeaturesRequest;
  static deserializeBinaryFromReader(message: GetUserFeaturesRequest, reader: jspb.BinaryReader): GetUserFeaturesRequest;
}

export namespace GetUserFeaturesRequest {
  export type AsObject = {
    userId: string,
    inheritance: boolean,
  }
}

export class GetUserFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): GetUserFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): GetUserFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserFeaturesResponse): GetUserFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserFeaturesResponse;
  static deserializeBinaryFromReader(message: GetUserFeaturesResponse, reader: jspb.BinaryReader): GetUserFeaturesResponse;
}

export namespace GetUserFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}


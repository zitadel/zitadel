import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"


export class Metadata extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Metadata;
  hasDetails(): boolean;
  clearDetails(): Metadata;

  getKey(): string;
  setKey(value: string): Metadata;

  getValue(): Uint8Array | string;
  getValue_asU8(): Uint8Array;
  getValue_asB64(): string;
  setValue(value: Uint8Array | string): Metadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Metadata.AsObject;
  static toObject(includeInstance: boolean, msg: Metadata): Metadata.AsObject;
  static serializeBinaryToWriter(message: Metadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Metadata;
  static deserializeBinaryFromReader(message: Metadata, reader: jspb.BinaryReader): Metadata;
}

export namespace Metadata {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    key: string,
    value: Uint8Array | string,
  }
}

export class MetadataQuery extends jspb.Message {
  getKeyQuery(): MetadataKeyQuery | undefined;
  setKeyQuery(value?: MetadataKeyQuery): MetadataQuery;
  hasKeyQuery(): boolean;
  clearKeyQuery(): MetadataQuery;

  getQueryCase(): MetadataQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetadataQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MetadataQuery): MetadataQuery.AsObject;
  static serializeBinaryToWriter(message: MetadataQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetadataQuery;
  static deserializeBinaryFromReader(message: MetadataQuery, reader: jspb.BinaryReader): MetadataQuery;
}

export namespace MetadataQuery {
  export type AsObject = {
    keyQuery?: MetadataKeyQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    KEY_QUERY = 1,
  }
}

export class MetadataKeyQuery extends jspb.Message {
  getKey(): string;
  setKey(value: string): MetadataKeyQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): MetadataKeyQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MetadataKeyQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MetadataKeyQuery): MetadataKeyQuery.AsObject;
  static serializeBinaryToWriter(message: MetadataKeyQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MetadataKeyQuery;
  static deserializeBinaryFromReader(message: MetadataKeyQuery, reader: jspb.BinaryReader): MetadataKeyQuery;
}

export namespace MetadataKeyQuery {
  export type AsObject = {
    key: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}


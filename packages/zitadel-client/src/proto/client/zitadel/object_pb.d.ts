import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class ObjectDetails extends jspb.Message {
  getSequence(): number;
  setSequence(value: number): ObjectDetails;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): ObjectDetails;
  hasCreationDate(): boolean;
  clearCreationDate(): ObjectDetails;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): ObjectDetails;
  hasChangeDate(): boolean;
  clearChangeDate(): ObjectDetails;

  getResourceOwner(): string;
  setResourceOwner(value: string): ObjectDetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ObjectDetails.AsObject;
  static toObject(includeInstance: boolean, msg: ObjectDetails): ObjectDetails.AsObject;
  static serializeBinaryToWriter(message: ObjectDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ObjectDetails;
  static deserializeBinaryFromReader(message: ObjectDetails, reader: jspb.BinaryReader): ObjectDetails;
}

export namespace ObjectDetails {
  export type AsObject = {
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    resourceOwner: string,
  }
}

export class ListQuery extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): ListQuery;

  getLimit(): number;
  setLimit(value: number): ListQuery;

  getAsc(): boolean;
  setAsc(value: boolean): ListQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ListQuery): ListQuery.AsObject;
  static serializeBinaryToWriter(message: ListQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListQuery;
  static deserializeBinaryFromReader(message: ListQuery, reader: jspb.BinaryReader): ListQuery;
}

export namespace ListQuery {
  export type AsObject = {
    offset: number,
    limit: number,
    asc: boolean,
  }
}

export class ListDetails extends jspb.Message {
  getTotalResult(): number;
  setTotalResult(value: number): ListDetails;

  getProcessedSequence(): number;
  setProcessedSequence(value: number): ListDetails;

  getViewTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setViewTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): ListDetails;
  hasViewTimestamp(): boolean;
  clearViewTimestamp(): ListDetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDetails.AsObject;
  static toObject(includeInstance: boolean, msg: ListDetails): ListDetails.AsObject;
  static serializeBinaryToWriter(message: ListDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDetails;
  static deserializeBinaryFromReader(message: ListDetails, reader: jspb.BinaryReader): ListDetails;
}

export namespace ListDetails {
  export type AsObject = {
    totalResult: number,
    processedSequence: number,
    viewTimestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export enum TextQueryMethod { 
  TEXT_QUERY_METHOD_EQUALS = 0,
  TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE = 1,
  TEXT_QUERY_METHOD_STARTS_WITH = 2,
  TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE = 3,
  TEXT_QUERY_METHOD_CONTAINS = 4,
  TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE = 5,
  TEXT_QUERY_METHOD_ENDS_WITH = 6,
  TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE = 7,
}
export enum ListQueryMethod { 
  LIST_QUERY_METHOD_IN = 0,
}
export enum TimestampQueryMethod { 
  TIMESTAMP_QUERY_METHOD_EQUALS = 0,
  TIMESTAMP_QUERY_METHOD_GREATER = 1,
  TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS = 2,
  TIMESTAMP_QUERY_METHOD_LESS = 3,
  TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS = 4,
}

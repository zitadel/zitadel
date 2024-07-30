import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../../../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class Milestone extends jspb.Message {
  getType(): MilestoneType;
  setType(value: MilestoneType): Milestone;

  getReachedDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setReachedDate(value?: google_protobuf_timestamp_pb.Timestamp): Milestone;
  hasReachedDate(): boolean;
  clearReachedDate(): Milestone;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Milestone.AsObject;
  static toObject(includeInstance: boolean, msg: Milestone): Milestone.AsObject;
  static serializeBinaryToWriter(message: Milestone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Milestone;
  static deserializeBinaryFromReader(message: Milestone, reader: jspb.BinaryReader): Milestone;
}

export namespace Milestone {
  export type AsObject = {
    type: MilestoneType,
    reachedDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class MilestoneQuery extends jspb.Message {
  getIsReachedQuery(): IsReachedQuery | undefined;
  setIsReachedQuery(value?: IsReachedQuery): MilestoneQuery;
  hasIsReachedQuery(): boolean;
  clearIsReachedQuery(): MilestoneQuery;

  getQueryCase(): MilestoneQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MilestoneQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MilestoneQuery): MilestoneQuery.AsObject;
  static serializeBinaryToWriter(message: MilestoneQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MilestoneQuery;
  static deserializeBinaryFromReader(message: MilestoneQuery, reader: jspb.BinaryReader): MilestoneQuery;
}

export namespace MilestoneQuery {
  export type AsObject = {
    isReachedQuery?: IsReachedQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    IS_REACHED_QUERY = 1,
  }
}

export class IsReachedQuery extends jspb.Message {
  getReached(): boolean;
  setReached(value: boolean): IsReachedQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsReachedQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IsReachedQuery): IsReachedQuery.AsObject;
  static serializeBinaryToWriter(message: IsReachedQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsReachedQuery;
  static deserializeBinaryFromReader(message: IsReachedQuery, reader: jspb.BinaryReader): IsReachedQuery;
}

export namespace IsReachedQuery {
  export type AsObject = {
    reached: boolean,
  }
}

export enum MilestoneType { 
  MILESTONE_TYPE_UNSPECIFIED = 0,
  MILESTONE_TYPE_INSTANCE_CREATED = 1,
  MILESTONE_TYPE_AUTHENTICATION_SUCCEEDED_ON_INSTANCE = 2,
  MILESTONE_TYPE_PROJECT_CREATED = 3,
  MILESTONE_TYPE_APPLICATION_CREATED = 4,
  MILESTONE_TYPE_AUTHENTICATION_SUCCEEDED_ON_APPLICATION = 5,
  MILESTONE_TYPE_INSTANCE_DELETED = 6,
}
export enum MilestoneFieldName { 
  MILESTONE_FIELD_NAME_UNSPECIFIED = 0,
  MILESTONE_FIELD_NAME_TYPE = 1,
  MILESTONE_FIELD_NAME_REACHED_DATE = 2,
}

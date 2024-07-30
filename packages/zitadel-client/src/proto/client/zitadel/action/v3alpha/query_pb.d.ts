import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_action_v3alpha_execution_pb from '../../../zitadel/action/v3alpha/execution_pb'; // proto import: "zitadel/action/v3alpha/execution.proto"


export class SearchQuery extends jspb.Message {
  getInConditionsQuery(): InConditionsQuery | undefined;
  setInConditionsQuery(value?: InConditionsQuery): SearchQuery;
  hasInConditionsQuery(): boolean;
  clearInConditionsQuery(): SearchQuery;

  getExecutionTypeQuery(): ExecutionTypeQuery | undefined;
  setExecutionTypeQuery(value?: ExecutionTypeQuery): SearchQuery;
  hasExecutionTypeQuery(): boolean;
  clearExecutionTypeQuery(): SearchQuery;

  getTargetQuery(): TargetQuery | undefined;
  setTargetQuery(value?: TargetQuery): SearchQuery;
  hasTargetQuery(): boolean;
  clearTargetQuery(): SearchQuery;

  getIncludeQuery(): IncludeQuery | undefined;
  setIncludeQuery(value?: IncludeQuery): SearchQuery;
  hasIncludeQuery(): boolean;
  clearIncludeQuery(): SearchQuery;

  getQueryCase(): SearchQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: SearchQuery): SearchQuery.AsObject;
  static serializeBinaryToWriter(message: SearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchQuery;
  static deserializeBinaryFromReader(message: SearchQuery, reader: jspb.BinaryReader): SearchQuery;
}

export namespace SearchQuery {
  export type AsObject = {
    inConditionsQuery?: InConditionsQuery.AsObject,
    executionTypeQuery?: ExecutionTypeQuery.AsObject,
    targetQuery?: TargetQuery.AsObject,
    includeQuery?: IncludeQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    IN_CONDITIONS_QUERY = 1,
    EXECUTION_TYPE_QUERY = 2,
    TARGET_QUERY = 3,
    INCLUDE_QUERY = 4,
  }
}

export class InConditionsQuery extends jspb.Message {
  getConditionsList(): Array<zitadel_action_v3alpha_execution_pb.Condition>;
  setConditionsList(value: Array<zitadel_action_v3alpha_execution_pb.Condition>): InConditionsQuery;
  clearConditionsList(): InConditionsQuery;
  addConditions(value?: zitadel_action_v3alpha_execution_pb.Condition, index?: number): zitadel_action_v3alpha_execution_pb.Condition;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InConditionsQuery.AsObject;
  static toObject(includeInstance: boolean, msg: InConditionsQuery): InConditionsQuery.AsObject;
  static serializeBinaryToWriter(message: InConditionsQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InConditionsQuery;
  static deserializeBinaryFromReader(message: InConditionsQuery, reader: jspb.BinaryReader): InConditionsQuery;
}

export namespace InConditionsQuery {
  export type AsObject = {
    conditionsList: Array<zitadel_action_v3alpha_execution_pb.Condition.AsObject>,
  }
}

export class ExecutionTypeQuery extends jspb.Message {
  getExecutionType(): ExecutionType;
  setExecutionType(value: ExecutionType): ExecutionTypeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecutionTypeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ExecutionTypeQuery): ExecutionTypeQuery.AsObject;
  static serializeBinaryToWriter(message: ExecutionTypeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecutionTypeQuery;
  static deserializeBinaryFromReader(message: ExecutionTypeQuery, reader: jspb.BinaryReader): ExecutionTypeQuery;
}

export namespace ExecutionTypeQuery {
  export type AsObject = {
    executionType: ExecutionType,
  }
}

export class TargetQuery extends jspb.Message {
  getTargetId(): string;
  setTargetId(value: string): TargetQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetQuery.AsObject;
  static toObject(includeInstance: boolean, msg: TargetQuery): TargetQuery.AsObject;
  static serializeBinaryToWriter(message: TargetQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetQuery;
  static deserializeBinaryFromReader(message: TargetQuery, reader: jspb.BinaryReader): TargetQuery;
}

export namespace TargetQuery {
  export type AsObject = {
    targetId: string,
  }
}

export class IncludeQuery extends jspb.Message {
  getInclude(): zitadel_action_v3alpha_execution_pb.Condition | undefined;
  setInclude(value?: zitadel_action_v3alpha_execution_pb.Condition): IncludeQuery;
  hasInclude(): boolean;
  clearInclude(): IncludeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IncludeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IncludeQuery): IncludeQuery.AsObject;
  static serializeBinaryToWriter(message: IncludeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IncludeQuery;
  static deserializeBinaryFromReader(message: IncludeQuery, reader: jspb.BinaryReader): IncludeQuery;
}

export namespace IncludeQuery {
  export type AsObject = {
    include?: zitadel_action_v3alpha_execution_pb.Condition.AsObject,
  }
}

export class TargetSearchQuery extends jspb.Message {
  getTargetNameQuery(): TargetNameQuery | undefined;
  setTargetNameQuery(value?: TargetNameQuery): TargetSearchQuery;
  hasTargetNameQuery(): boolean;
  clearTargetNameQuery(): TargetSearchQuery;

  getInTargetIdsQuery(): InTargetIDsQuery | undefined;
  setInTargetIdsQuery(value?: InTargetIDsQuery): TargetSearchQuery;
  hasInTargetIdsQuery(): boolean;
  clearInTargetIdsQuery(): TargetSearchQuery;

  getQueryCase(): TargetSearchQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: TargetSearchQuery): TargetSearchQuery.AsObject;
  static serializeBinaryToWriter(message: TargetSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetSearchQuery;
  static deserializeBinaryFromReader(message: TargetSearchQuery, reader: jspb.BinaryReader): TargetSearchQuery;
}

export namespace TargetSearchQuery {
  export type AsObject = {
    targetNameQuery?: TargetNameQuery.AsObject,
    inTargetIdsQuery?: InTargetIDsQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    TARGET_NAME_QUERY = 1,
    IN_TARGET_IDS_QUERY = 2,
  }
}

export class TargetNameQuery extends jspb.Message {
  getTargetName(): string;
  setTargetName(value: string): TargetNameQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): TargetNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: TargetNameQuery): TargetNameQuery.AsObject;
  static serializeBinaryToWriter(message: TargetNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetNameQuery;
  static deserializeBinaryFromReader(message: TargetNameQuery, reader: jspb.BinaryReader): TargetNameQuery;
}

export namespace TargetNameQuery {
  export type AsObject = {
    targetName: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class InTargetIDsQuery extends jspb.Message {
  getTargetIdsList(): Array<string>;
  setTargetIdsList(value: Array<string>): InTargetIDsQuery;
  clearTargetIdsList(): InTargetIDsQuery;
  addTargetIds(value: string, index?: number): InTargetIDsQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InTargetIDsQuery.AsObject;
  static toObject(includeInstance: boolean, msg: InTargetIDsQuery): InTargetIDsQuery.AsObject;
  static serializeBinaryToWriter(message: InTargetIDsQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InTargetIDsQuery;
  static deserializeBinaryFromReader(message: InTargetIDsQuery, reader: jspb.BinaryReader): InTargetIDsQuery;
}

export namespace InTargetIDsQuery {
  export type AsObject = {
    targetIdsList: Array<string>,
  }
}

export enum ExecutionType { 
  EXECUTION_TYPE_UNSPECIFIED = 0,
  EXECUTION_TYPE_REQUEST = 1,
  EXECUTION_TYPE_RESPONSE = 2,
  EXECUTION_TYPE_EVENT = 3,
  EXECUTION_TYPE_FUNCTION = 4,
}
export enum TargetFieldName { 
  FIELD_NAME_UNSPECIFIED = 0,
  FIELD_NAME_ID = 1,
  FIELD_NAME_CREATION_DATE = 2,
  FIELD_NAME_CHANGE_DATE = 3,
  FIELD_NAME_NAME = 4,
  FIELD_NAME_TARGET_TYPE = 5,
  FIELD_NAME_URL = 6,
  FIELD_NAME_TIMEOUT = 7,
  FIELD_NAME_ASYNC = 8,
  FIELD_NAME_INTERRUPT_ON_ERROR = 9,
}

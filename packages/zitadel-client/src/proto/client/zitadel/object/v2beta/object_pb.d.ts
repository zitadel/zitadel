import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class Organisation extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): Organisation;

  getOrgDomain(): string;
  setOrgDomain(value: string): Organisation;

  getOrgCase(): Organisation.OrgCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Organisation.AsObject;
  static toObject(includeInstance: boolean, msg: Organisation): Organisation.AsObject;
  static serializeBinaryToWriter(message: Organisation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Organisation;
  static deserializeBinaryFromReader(message: Organisation, reader: jspb.BinaryReader): Organisation;
}

export namespace Organisation {
  export type AsObject = {
    orgId: string,
    orgDomain: string,
  }

  export enum OrgCase { 
    ORG_NOT_SET = 0,
    ORG_ID = 1,
    ORG_DOMAIN = 2,
  }
}

export class Organization extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): Organization;

  getOrgDomain(): string;
  setOrgDomain(value: string): Organization;

  getOrgCase(): Organization.OrgCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Organization.AsObject;
  static toObject(includeInstance: boolean, msg: Organization): Organization.AsObject;
  static serializeBinaryToWriter(message: Organization, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Organization;
  static deserializeBinaryFromReader(message: Organization, reader: jspb.BinaryReader): Organization;
}

export namespace Organization {
  export type AsObject = {
    orgId: string,
    orgDomain: string,
  }

  export enum OrgCase { 
    ORG_NOT_SET = 0,
    ORG_ID = 1,
    ORG_DOMAIN = 2,
  }
}

export class RequestContext extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): RequestContext;

  getInstance(): boolean;
  setInstance(value: boolean): RequestContext;

  getResourceOwnerCase(): RequestContext.ResourceOwnerCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestContext.AsObject;
  static toObject(includeInstance: boolean, msg: RequestContext): RequestContext.AsObject;
  static serializeBinaryToWriter(message: RequestContext, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestContext;
  static deserializeBinaryFromReader(message: RequestContext, reader: jspb.BinaryReader): RequestContext;
}

export namespace RequestContext {
  export type AsObject = {
    orgId: string,
    instance: boolean,
  }

  export enum ResourceOwnerCase { 
    RESOURCE_OWNER_NOT_SET = 0,
    ORG_ID = 1,
    INSTANCE = 2,
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

export class Details extends jspb.Message {
  getSequence(): number;
  setSequence(value: number): Details;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): Details;
  hasChangeDate(): boolean;
  clearChangeDate(): Details;

  getResourceOwner(): string;
  setResourceOwner(value: string): Details;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Details.AsObject;
  static toObject(includeInstance: boolean, msg: Details): Details.AsObject;
  static serializeBinaryToWriter(message: Details, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Details;
  static deserializeBinaryFromReader(message: Details, reader: jspb.BinaryReader): Details;
}

export namespace Details {
  export type AsObject = {
    sequence: number,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    resourceOwner: string,
  }
}

export class ListDetails extends jspb.Message {
  getTotalResult(): number;
  setTotalResult(value: number): ListDetails;

  getProcessedSequence(): number;
  setProcessedSequence(value: number): ListDetails;

  getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): ListDetails;
  hasTimestamp(): boolean;
  clearTimestamp(): ListDetails;

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
    timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
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

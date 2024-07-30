/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Timestamp } from "../../../google/protobuf/timestamp";

export const protobufPackage = "zitadel.object.v2beta";

export enum TextQueryMethod {
  TEXT_QUERY_METHOD_EQUALS = 0,
  TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE = 1,
  TEXT_QUERY_METHOD_STARTS_WITH = 2,
  TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE = 3,
  TEXT_QUERY_METHOD_CONTAINS = 4,
  TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE = 5,
  TEXT_QUERY_METHOD_ENDS_WITH = 6,
  TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE = 7,
  UNRECOGNIZED = -1,
}

export function textQueryMethodFromJSON(object: any): TextQueryMethod {
  switch (object) {
    case 0:
    case "TEXT_QUERY_METHOD_EQUALS":
      return TextQueryMethod.TEXT_QUERY_METHOD_EQUALS;
    case 1:
    case "TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE":
      return TextQueryMethod.TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE;
    case 2:
    case "TEXT_QUERY_METHOD_STARTS_WITH":
      return TextQueryMethod.TEXT_QUERY_METHOD_STARTS_WITH;
    case 3:
    case "TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE":
      return TextQueryMethod.TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE;
    case 4:
    case "TEXT_QUERY_METHOD_CONTAINS":
      return TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS;
    case 5:
    case "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE":
      return TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE;
    case 6:
    case "TEXT_QUERY_METHOD_ENDS_WITH":
      return TextQueryMethod.TEXT_QUERY_METHOD_ENDS_WITH;
    case 7:
    case "TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE":
      return TextQueryMethod.TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return TextQueryMethod.UNRECOGNIZED;
  }
}

export function textQueryMethodToJSON(object: TextQueryMethod): string {
  switch (object) {
    case TextQueryMethod.TEXT_QUERY_METHOD_EQUALS:
      return "TEXT_QUERY_METHOD_EQUALS";
    case TextQueryMethod.TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE:
      return "TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE";
    case TextQueryMethod.TEXT_QUERY_METHOD_STARTS_WITH:
      return "TEXT_QUERY_METHOD_STARTS_WITH";
    case TextQueryMethod.TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE:
      return "TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE";
    case TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS:
      return "TEXT_QUERY_METHOD_CONTAINS";
    case TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE:
      return "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE";
    case TextQueryMethod.TEXT_QUERY_METHOD_ENDS_WITH:
      return "TEXT_QUERY_METHOD_ENDS_WITH";
    case TextQueryMethod.TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE:
      return "TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE";
    case TextQueryMethod.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ListQueryMethod {
  LIST_QUERY_METHOD_IN = 0,
  UNRECOGNIZED = -1,
}

export function listQueryMethodFromJSON(object: any): ListQueryMethod {
  switch (object) {
    case 0:
    case "LIST_QUERY_METHOD_IN":
      return ListQueryMethod.LIST_QUERY_METHOD_IN;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ListQueryMethod.UNRECOGNIZED;
  }
}

export function listQueryMethodToJSON(object: ListQueryMethod): string {
  switch (object) {
    case ListQueryMethod.LIST_QUERY_METHOD_IN:
      return "LIST_QUERY_METHOD_IN";
    case ListQueryMethod.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum TimestampQueryMethod {
  TIMESTAMP_QUERY_METHOD_EQUALS = 0,
  TIMESTAMP_QUERY_METHOD_GREATER = 1,
  TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS = 2,
  TIMESTAMP_QUERY_METHOD_LESS = 3,
  TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS = 4,
  UNRECOGNIZED = -1,
}

export function timestampQueryMethodFromJSON(object: any): TimestampQueryMethod {
  switch (object) {
    case 0:
    case "TIMESTAMP_QUERY_METHOD_EQUALS":
      return TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_EQUALS;
    case 1:
    case "TIMESTAMP_QUERY_METHOD_GREATER":
      return TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_GREATER;
    case 2:
    case "TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS":
      return TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS;
    case 3:
    case "TIMESTAMP_QUERY_METHOD_LESS":
      return TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_LESS;
    case 4:
    case "TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS":
      return TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS;
    case -1:
    case "UNRECOGNIZED":
    default:
      return TimestampQueryMethod.UNRECOGNIZED;
  }
}

export function timestampQueryMethodToJSON(object: TimestampQueryMethod): string {
  switch (object) {
    case TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_EQUALS:
      return "TIMESTAMP_QUERY_METHOD_EQUALS";
    case TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_GREATER:
      return "TIMESTAMP_QUERY_METHOD_GREATER";
    case TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS:
      return "TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS";
    case TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_LESS:
      return "TIMESTAMP_QUERY_METHOD_LESS";
    case TimestampQueryMethod.TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS:
      return "TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS";
    case TimestampQueryMethod.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** Deprecated: use Organization */
export interface Organisation {
  orgId?: string | undefined;
  orgDomain?: string | undefined;
}

export interface Organization {
  orgId?: string | undefined;
  orgDomain?: string | undefined;
}

export interface RequestContext {
  orgId?: string | undefined;
  instance?: boolean | undefined;
}

export interface ListQuery {
  offset: number;
  limit: number;
  asc: boolean;
}

export interface Details {
  /**
   * sequence represents the order of events. It's always counting
   *
   * on read: the sequence of the last event reduced by the projection
   *
   * on manipulation: the timestamp of the event(s) added by the manipulation
   */
  sequence: number;
  /**
   * change_date is the timestamp when the object was changed
   *
   * on read: the timestamp of the last event reduced by the projection
   *
   * on manipulation: the timestamp of the event(s) added by the manipulation
   */
  changeDate:
    | Date
    | undefined;
  /** resource_owner is the organization or instance_id an object belongs to */
  resourceOwner: string;
}

export interface ListDetails {
  totalResult: number;
  processedSequence: number;
  timestamp: Date | undefined;
}

function createBaseOrganisation(): Organisation {
  return { orgId: undefined, orgDomain: undefined };
}

export const Organisation = {
  encode(message: Organisation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== undefined) {
      writer.uint32(10).string(message.orgId);
    }
    if (message.orgDomain !== undefined) {
      writer.uint32(18).string(message.orgDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Organisation {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrganisation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.orgId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.orgDomain = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Organisation {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : undefined,
      orgDomain: isSet(object.orgDomain) ? String(object.orgDomain) : undefined,
    };
  },

  toJSON(message: Organisation): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.orgDomain !== undefined && (obj.orgDomain = message.orgDomain);
    return obj;
  },

  create(base?: DeepPartial<Organisation>): Organisation {
    return Organisation.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Organisation>): Organisation {
    const message = createBaseOrganisation();
    message.orgId = object.orgId ?? undefined;
    message.orgDomain = object.orgDomain ?? undefined;
    return message;
  },
};

function createBaseOrganization(): Organization {
  return { orgId: undefined, orgDomain: undefined };
}

export const Organization = {
  encode(message: Organization, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== undefined) {
      writer.uint32(10).string(message.orgId);
    }
    if (message.orgDomain !== undefined) {
      writer.uint32(18).string(message.orgDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Organization {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrganization();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.orgId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.orgDomain = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Organization {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : undefined,
      orgDomain: isSet(object.orgDomain) ? String(object.orgDomain) : undefined,
    };
  },

  toJSON(message: Organization): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.orgDomain !== undefined && (obj.orgDomain = message.orgDomain);
    return obj;
  },

  create(base?: DeepPartial<Organization>): Organization {
    return Organization.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Organization>): Organization {
    const message = createBaseOrganization();
    message.orgId = object.orgId ?? undefined;
    message.orgDomain = object.orgDomain ?? undefined;
    return message;
  },
};

function createBaseRequestContext(): RequestContext {
  return { orgId: undefined, instance: undefined };
}

export const RequestContext = {
  encode(message: RequestContext, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== undefined) {
      writer.uint32(10).string(message.orgId);
    }
    if (message.instance !== undefined) {
      writer.uint32(16).bool(message.instance);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestContext {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestContext();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.orgId = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.instance = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequestContext {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : undefined,
      instance: isSet(object.instance) ? Boolean(object.instance) : undefined,
    };
  },

  toJSON(message: RequestContext): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.instance !== undefined && (obj.instance = message.instance);
    return obj;
  },

  create(base?: DeepPartial<RequestContext>): RequestContext {
    return RequestContext.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RequestContext>): RequestContext {
    const message = createBaseRequestContext();
    message.orgId = object.orgId ?? undefined;
    message.instance = object.instance ?? undefined;
    return message;
  },
};

function createBaseListQuery(): ListQuery {
  return { offset: 0, limit: 0, asc: false };
}

export const ListQuery = {
  encode(message: ListQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.offset !== 0) {
      writer.uint32(8).uint64(message.offset);
    }
    if (message.limit !== 0) {
      writer.uint32(16).uint32(message.limit);
    }
    if (message.asc === true) {
      writer.uint32(24).bool(message.asc);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.offset = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.limit = reader.uint32();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.asc = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListQuery {
    return {
      offset: isSet(object.offset) ? Number(object.offset) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
      asc: isSet(object.asc) ? Boolean(object.asc) : false,
    };
  },

  toJSON(message: ListQuery): unknown {
    const obj: any = {};
    message.offset !== undefined && (obj.offset = Math.round(message.offset));
    message.limit !== undefined && (obj.limit = Math.round(message.limit));
    message.asc !== undefined && (obj.asc = message.asc);
    return obj;
  },

  create(base?: DeepPartial<ListQuery>): ListQuery {
    return ListQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListQuery>): ListQuery {
    const message = createBaseListQuery();
    message.offset = object.offset ?? 0;
    message.limit = object.limit ?? 0;
    message.asc = object.asc ?? false;
    return message;
  },
};

function createBaseDetails(): Details {
  return { sequence: 0, changeDate: undefined, resourceOwner: "" };
}

export const Details = {
  encode(message: Details, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sequence !== 0) {
      writer.uint32(8).uint64(message.sequence);
    }
    if (message.changeDate !== undefined) {
      Timestamp.encode(toTimestamp(message.changeDate), writer.uint32(18).fork()).ldelim();
    }
    if (message.resourceOwner !== "") {
      writer.uint32(26).string(message.resourceOwner);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Details {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.sequence = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.changeDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.resourceOwner = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Details {
    return {
      sequence: isSet(object.sequence) ? Number(object.sequence) : 0,
      changeDate: isSet(object.changeDate) ? fromJsonTimestamp(object.changeDate) : undefined,
      resourceOwner: isSet(object.resourceOwner) ? String(object.resourceOwner) : "",
    };
  },

  toJSON(message: Details): unknown {
    const obj: any = {};
    message.sequence !== undefined && (obj.sequence = Math.round(message.sequence));
    message.changeDate !== undefined && (obj.changeDate = message.changeDate.toISOString());
    message.resourceOwner !== undefined && (obj.resourceOwner = message.resourceOwner);
    return obj;
  },

  create(base?: DeepPartial<Details>): Details {
    return Details.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Details>): Details {
    const message = createBaseDetails();
    message.sequence = object.sequence ?? 0;
    message.changeDate = object.changeDate ?? undefined;
    message.resourceOwner = object.resourceOwner ?? "";
    return message;
  },
};

function createBaseListDetails(): ListDetails {
  return { totalResult: 0, processedSequence: 0, timestamp: undefined };
}

export const ListDetails = {
  encode(message: ListDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.totalResult !== 0) {
      writer.uint32(8).uint64(message.totalResult);
    }
    if (message.processedSequence !== 0) {
      writer.uint32(16).uint64(message.processedSequence);
    }
    if (message.timestamp !== undefined) {
      Timestamp.encode(toTimestamp(message.timestamp), writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListDetails {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListDetails();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.totalResult = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.processedSequence = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.timestamp = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListDetails {
    return {
      totalResult: isSet(object.totalResult) ? Number(object.totalResult) : 0,
      processedSequence: isSet(object.processedSequence) ? Number(object.processedSequence) : 0,
      timestamp: isSet(object.timestamp) ? fromJsonTimestamp(object.timestamp) : undefined,
    };
  },

  toJSON(message: ListDetails): unknown {
    const obj: any = {};
    message.totalResult !== undefined && (obj.totalResult = Math.round(message.totalResult));
    message.processedSequence !== undefined && (obj.processedSequence = Math.round(message.processedSequence));
    message.timestamp !== undefined && (obj.timestamp = message.timestamp.toISOString());
    return obj;
  },

  create(base?: DeepPartial<ListDetails>): ListDetails {
    return ListDetails.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListDetails>): ListDetails {
    const message = createBaseListDetails();
    message.totalResult = object.totalResult ?? 0;
    message.processedSequence = object.processedSequence ?? 0;
    message.timestamp = object.timestamp ?? undefined;
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

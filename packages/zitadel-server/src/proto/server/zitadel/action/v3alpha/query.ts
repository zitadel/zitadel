/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { TextQueryMethod, textQueryMethodFromJSON, textQueryMethodToJSON } from "../../object/v2beta/object";
import { Condition } from "./execution";

export const protobufPackage = "zitadel.action.v3alpha";

export enum ExecutionType {
  EXECUTION_TYPE_UNSPECIFIED = 0,
  EXECUTION_TYPE_REQUEST = 1,
  EXECUTION_TYPE_RESPONSE = 2,
  EXECUTION_TYPE_EVENT = 3,
  EXECUTION_TYPE_FUNCTION = 4,
  UNRECOGNIZED = -1,
}

export function executionTypeFromJSON(object: any): ExecutionType {
  switch (object) {
    case 0:
    case "EXECUTION_TYPE_UNSPECIFIED":
      return ExecutionType.EXECUTION_TYPE_UNSPECIFIED;
    case 1:
    case "EXECUTION_TYPE_REQUEST":
      return ExecutionType.EXECUTION_TYPE_REQUEST;
    case 2:
    case "EXECUTION_TYPE_RESPONSE":
      return ExecutionType.EXECUTION_TYPE_RESPONSE;
    case 3:
    case "EXECUTION_TYPE_EVENT":
      return ExecutionType.EXECUTION_TYPE_EVENT;
    case 4:
    case "EXECUTION_TYPE_FUNCTION":
      return ExecutionType.EXECUTION_TYPE_FUNCTION;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ExecutionType.UNRECOGNIZED;
  }
}

export function executionTypeToJSON(object: ExecutionType): string {
  switch (object) {
    case ExecutionType.EXECUTION_TYPE_UNSPECIFIED:
      return "EXECUTION_TYPE_UNSPECIFIED";
    case ExecutionType.EXECUTION_TYPE_REQUEST:
      return "EXECUTION_TYPE_REQUEST";
    case ExecutionType.EXECUTION_TYPE_RESPONSE:
      return "EXECUTION_TYPE_RESPONSE";
    case ExecutionType.EXECUTION_TYPE_EVENT:
      return "EXECUTION_TYPE_EVENT";
    case ExecutionType.EXECUTION_TYPE_FUNCTION:
      return "EXECUTION_TYPE_FUNCTION";
    case ExecutionType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
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
  UNRECOGNIZED = -1,
}

export function targetFieldNameFromJSON(object: any): TargetFieldName {
  switch (object) {
    case 0:
    case "FIELD_NAME_UNSPECIFIED":
      return TargetFieldName.FIELD_NAME_UNSPECIFIED;
    case 1:
    case "FIELD_NAME_ID":
      return TargetFieldName.FIELD_NAME_ID;
    case 2:
    case "FIELD_NAME_CREATION_DATE":
      return TargetFieldName.FIELD_NAME_CREATION_DATE;
    case 3:
    case "FIELD_NAME_CHANGE_DATE":
      return TargetFieldName.FIELD_NAME_CHANGE_DATE;
    case 4:
    case "FIELD_NAME_NAME":
      return TargetFieldName.FIELD_NAME_NAME;
    case 5:
    case "FIELD_NAME_TARGET_TYPE":
      return TargetFieldName.FIELD_NAME_TARGET_TYPE;
    case 6:
    case "FIELD_NAME_URL":
      return TargetFieldName.FIELD_NAME_URL;
    case 7:
    case "FIELD_NAME_TIMEOUT":
      return TargetFieldName.FIELD_NAME_TIMEOUT;
    case 8:
    case "FIELD_NAME_ASYNC":
      return TargetFieldName.FIELD_NAME_ASYNC;
    case 9:
    case "FIELD_NAME_INTERRUPT_ON_ERROR":
      return TargetFieldName.FIELD_NAME_INTERRUPT_ON_ERROR;
    case -1:
    case "UNRECOGNIZED":
    default:
      return TargetFieldName.UNRECOGNIZED;
  }
}

export function targetFieldNameToJSON(object: TargetFieldName): string {
  switch (object) {
    case TargetFieldName.FIELD_NAME_UNSPECIFIED:
      return "FIELD_NAME_UNSPECIFIED";
    case TargetFieldName.FIELD_NAME_ID:
      return "FIELD_NAME_ID";
    case TargetFieldName.FIELD_NAME_CREATION_DATE:
      return "FIELD_NAME_CREATION_DATE";
    case TargetFieldName.FIELD_NAME_CHANGE_DATE:
      return "FIELD_NAME_CHANGE_DATE";
    case TargetFieldName.FIELD_NAME_NAME:
      return "FIELD_NAME_NAME";
    case TargetFieldName.FIELD_NAME_TARGET_TYPE:
      return "FIELD_NAME_TARGET_TYPE";
    case TargetFieldName.FIELD_NAME_URL:
      return "FIELD_NAME_URL";
    case TargetFieldName.FIELD_NAME_TIMEOUT:
      return "FIELD_NAME_TIMEOUT";
    case TargetFieldName.FIELD_NAME_ASYNC:
      return "FIELD_NAME_ASYNC";
    case TargetFieldName.FIELD_NAME_INTERRUPT_ON_ERROR:
      return "FIELD_NAME_INTERRUPT_ON_ERROR";
    case TargetFieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface SearchQuery {
  inConditionsQuery?: InConditionsQuery | undefined;
  executionTypeQuery?: ExecutionTypeQuery | undefined;
  targetQuery?: TargetQuery | undefined;
  includeQuery?: IncludeQuery | undefined;
}

export interface InConditionsQuery {
  /** Defines the conditions to query for. */
  conditions: Condition[];
}

export interface ExecutionTypeQuery {
  /** Defines the type to query for. */
  executionType: ExecutionType;
}

export interface TargetQuery {
  /** Defines the id to query for. */
  targetId: string;
}

export interface IncludeQuery {
  /** Defines the include to query for. */
  include: Condition | undefined;
}

export interface TargetSearchQuery {
  targetNameQuery?: TargetNameQuery | undefined;
  inTargetIdsQuery?: InTargetIDsQuery | undefined;
}

export interface TargetNameQuery {
  /** Defines the name of the target to query for. */
  targetName: string;
  /** Defines which text comparison method used for the name query. */
  method: TextQueryMethod;
}

export interface InTargetIDsQuery {
  /** Defines the ids to query for. */
  targetIds: string[];
}

function createBaseSearchQuery(): SearchQuery {
  return {
    inConditionsQuery: undefined,
    executionTypeQuery: undefined,
    targetQuery: undefined,
    includeQuery: undefined,
  };
}

export const SearchQuery = {
  encode(message: SearchQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.inConditionsQuery !== undefined) {
      InConditionsQuery.encode(message.inConditionsQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.executionTypeQuery !== undefined) {
      ExecutionTypeQuery.encode(message.executionTypeQuery, writer.uint32(18).fork()).ldelim();
    }
    if (message.targetQuery !== undefined) {
      TargetQuery.encode(message.targetQuery, writer.uint32(26).fork()).ldelim();
    }
    if (message.includeQuery !== undefined) {
      IncludeQuery.encode(message.includeQuery, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SearchQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSearchQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.inConditionsQuery = InConditionsQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.executionTypeQuery = ExecutionTypeQuery.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.targetQuery = TargetQuery.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.includeQuery = IncludeQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SearchQuery {
    return {
      inConditionsQuery: isSet(object.inConditionsQuery)
        ? InConditionsQuery.fromJSON(object.inConditionsQuery)
        : undefined,
      executionTypeQuery: isSet(object.executionTypeQuery)
        ? ExecutionTypeQuery.fromJSON(object.executionTypeQuery)
        : undefined,
      targetQuery: isSet(object.targetQuery) ? TargetQuery.fromJSON(object.targetQuery) : undefined,
      includeQuery: isSet(object.includeQuery) ? IncludeQuery.fromJSON(object.includeQuery) : undefined,
    };
  },

  toJSON(message: SearchQuery): unknown {
    const obj: any = {};
    message.inConditionsQuery !== undefined && (obj.inConditionsQuery = message.inConditionsQuery
      ? InConditionsQuery.toJSON(message.inConditionsQuery)
      : undefined);
    message.executionTypeQuery !== undefined && (obj.executionTypeQuery = message.executionTypeQuery
      ? ExecutionTypeQuery.toJSON(message.executionTypeQuery)
      : undefined);
    message.targetQuery !== undefined &&
      (obj.targetQuery = message.targetQuery ? TargetQuery.toJSON(message.targetQuery) : undefined);
    message.includeQuery !== undefined &&
      (obj.includeQuery = message.includeQuery ? IncludeQuery.toJSON(message.includeQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SearchQuery>): SearchQuery {
    return SearchQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SearchQuery>): SearchQuery {
    const message = createBaseSearchQuery();
    message.inConditionsQuery = (object.inConditionsQuery !== undefined && object.inConditionsQuery !== null)
      ? InConditionsQuery.fromPartial(object.inConditionsQuery)
      : undefined;
    message.executionTypeQuery = (object.executionTypeQuery !== undefined && object.executionTypeQuery !== null)
      ? ExecutionTypeQuery.fromPartial(object.executionTypeQuery)
      : undefined;
    message.targetQuery = (object.targetQuery !== undefined && object.targetQuery !== null)
      ? TargetQuery.fromPartial(object.targetQuery)
      : undefined;
    message.includeQuery = (object.includeQuery !== undefined && object.includeQuery !== null)
      ? IncludeQuery.fromPartial(object.includeQuery)
      : undefined;
    return message;
  },
};

function createBaseInConditionsQuery(): InConditionsQuery {
  return { conditions: [] };
}

export const InConditionsQuery = {
  encode(message: InConditionsQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.conditions) {
      Condition.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InConditionsQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInConditionsQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.conditions.push(Condition.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): InConditionsQuery {
    return {
      conditions: Array.isArray(object?.conditions) ? object.conditions.map((e: any) => Condition.fromJSON(e)) : [],
    };
  },

  toJSON(message: InConditionsQuery): unknown {
    const obj: any = {};
    if (message.conditions) {
      obj.conditions = message.conditions.map((e) => e ? Condition.toJSON(e) : undefined);
    } else {
      obj.conditions = [];
    }
    return obj;
  },

  create(base?: DeepPartial<InConditionsQuery>): InConditionsQuery {
    return InConditionsQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InConditionsQuery>): InConditionsQuery {
    const message = createBaseInConditionsQuery();
    message.conditions = object.conditions?.map((e) => Condition.fromPartial(e)) || [];
    return message;
  },
};

function createBaseExecutionTypeQuery(): ExecutionTypeQuery {
  return { executionType: 0 };
}

export const ExecutionTypeQuery = {
  encode(message: ExecutionTypeQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.executionType !== 0) {
      writer.uint32(8).int32(message.executionType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExecutionTypeQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExecutionTypeQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.executionType = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExecutionTypeQuery {
    return { executionType: isSet(object.executionType) ? executionTypeFromJSON(object.executionType) : 0 };
  },

  toJSON(message: ExecutionTypeQuery): unknown {
    const obj: any = {};
    message.executionType !== undefined && (obj.executionType = executionTypeToJSON(message.executionType));
    return obj;
  },

  create(base?: DeepPartial<ExecutionTypeQuery>): ExecutionTypeQuery {
    return ExecutionTypeQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExecutionTypeQuery>): ExecutionTypeQuery {
    const message = createBaseExecutionTypeQuery();
    message.executionType = object.executionType ?? 0;
    return message;
  },
};

function createBaseTargetQuery(): TargetQuery {
  return { targetId: "" };
}

export const TargetQuery = {
  encode(message: TargetQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.targetId !== "") {
      writer.uint32(10).string(message.targetId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.targetId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TargetQuery {
    return { targetId: isSet(object.targetId) ? String(object.targetId) : "" };
  },

  toJSON(message: TargetQuery): unknown {
    const obj: any = {};
    message.targetId !== undefined && (obj.targetId = message.targetId);
    return obj;
  },

  create(base?: DeepPartial<TargetQuery>): TargetQuery {
    return TargetQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TargetQuery>): TargetQuery {
    const message = createBaseTargetQuery();
    message.targetId = object.targetId ?? "";
    return message;
  },
};

function createBaseIncludeQuery(): IncludeQuery {
  return { include: undefined };
}

export const IncludeQuery = {
  encode(message: IncludeQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.include !== undefined) {
      Condition.encode(message.include, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IncludeQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIncludeQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.include = Condition.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IncludeQuery {
    return { include: isSet(object.include) ? Condition.fromJSON(object.include) : undefined };
  },

  toJSON(message: IncludeQuery): unknown {
    const obj: any = {};
    message.include !== undefined && (obj.include = message.include ? Condition.toJSON(message.include) : undefined);
    return obj;
  },

  create(base?: DeepPartial<IncludeQuery>): IncludeQuery {
    return IncludeQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IncludeQuery>): IncludeQuery {
    const message = createBaseIncludeQuery();
    message.include = (object.include !== undefined && object.include !== null)
      ? Condition.fromPartial(object.include)
      : undefined;
    return message;
  },
};

function createBaseTargetSearchQuery(): TargetSearchQuery {
  return { targetNameQuery: undefined, inTargetIdsQuery: undefined };
}

export const TargetSearchQuery = {
  encode(message: TargetSearchQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.targetNameQuery !== undefined) {
      TargetNameQuery.encode(message.targetNameQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.inTargetIdsQuery !== undefined) {
      InTargetIDsQuery.encode(message.inTargetIdsQuery, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetSearchQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetSearchQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.targetNameQuery = TargetNameQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.inTargetIdsQuery = InTargetIDsQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TargetSearchQuery {
    return {
      targetNameQuery: isSet(object.targetNameQuery) ? TargetNameQuery.fromJSON(object.targetNameQuery) : undefined,
      inTargetIdsQuery: isSet(object.inTargetIdsQuery) ? InTargetIDsQuery.fromJSON(object.inTargetIdsQuery) : undefined,
    };
  },

  toJSON(message: TargetSearchQuery): unknown {
    const obj: any = {};
    message.targetNameQuery !== undefined &&
      (obj.targetNameQuery = message.targetNameQuery ? TargetNameQuery.toJSON(message.targetNameQuery) : undefined);
    message.inTargetIdsQuery !== undefined &&
      (obj.inTargetIdsQuery = message.inTargetIdsQuery ? InTargetIDsQuery.toJSON(message.inTargetIdsQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<TargetSearchQuery>): TargetSearchQuery {
    return TargetSearchQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TargetSearchQuery>): TargetSearchQuery {
    const message = createBaseTargetSearchQuery();
    message.targetNameQuery = (object.targetNameQuery !== undefined && object.targetNameQuery !== null)
      ? TargetNameQuery.fromPartial(object.targetNameQuery)
      : undefined;
    message.inTargetIdsQuery = (object.inTargetIdsQuery !== undefined && object.inTargetIdsQuery !== null)
      ? InTargetIDsQuery.fromPartial(object.inTargetIdsQuery)
      : undefined;
    return message;
  },
};

function createBaseTargetNameQuery(): TargetNameQuery {
  return { targetName: "", method: 0 };
}

export const TargetNameQuery = {
  encode(message: TargetNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.targetName !== "") {
      writer.uint32(10).string(message.targetName);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TargetNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTargetNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.targetName = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TargetNameQuery {
    return {
      targetName: isSet(object.targetName) ? String(object.targetName) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: TargetNameQuery): unknown {
    const obj: any = {};
    message.targetName !== undefined && (obj.targetName = message.targetName);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<TargetNameQuery>): TargetNameQuery {
    return TargetNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TargetNameQuery>): TargetNameQuery {
    const message = createBaseTargetNameQuery();
    message.targetName = object.targetName ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseInTargetIDsQuery(): InTargetIDsQuery {
  return { targetIds: [] };
}

export const InTargetIDsQuery = {
  encode(message: InTargetIDsQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.targetIds) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InTargetIDsQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInTargetIDsQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.targetIds.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): InTargetIDsQuery {
    return { targetIds: Array.isArray(object?.targetIds) ? object.targetIds.map((e: any) => String(e)) : [] };
  },

  toJSON(message: InTargetIDsQuery): unknown {
    const obj: any = {};
    if (message.targetIds) {
      obj.targetIds = message.targetIds.map((e) => e);
    } else {
      obj.targetIds = [];
    }
    return obj;
  },

  create(base?: DeepPartial<InTargetIDsQuery>): InTargetIDsQuery {
    return InTargetIDsQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InTargetIDsQuery>): InTargetIDsQuery {
    const message = createBaseInTargetIDsQuery();
    message.targetIds = object.targetIds?.map((e) => e) || [];
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

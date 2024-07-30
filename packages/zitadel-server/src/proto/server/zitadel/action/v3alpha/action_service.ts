/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Duration } from "../../../google/protobuf/duration";
import { Details, ListDetails, ListQuery } from "../../object/v2beta/object";
import { Condition, Execution, ExecutionTargetType } from "./execution";
import {
  SearchQuery,
  TargetFieldName,
  targetFieldNameFromJSON,
  targetFieldNameToJSON,
  TargetSearchQuery,
} from "./query";
import { SetRESTAsync, SetRESTCall, SetRESTWebhook, Target } from "./target";

export const protobufPackage = "zitadel.action.v3alpha";

export interface CreateTargetRequest {
  /** Unique name of the target. */
  name: string;
  restWebhook?: SetRESTWebhook | undefined;
  restCall?: SetRESTCall | undefined;
  restAsync?:
    | SetRESTAsync
    | undefined;
  /** Timeout defines the duration until ZITADEL cancels the execution. */
  timeout: Duration | undefined;
  endpoint: string;
}

export interface CreateTargetResponse {
  /** ID is the read-only unique identifier of the target. */
  id: string;
  /** Details provide some base information (such as the last change date) of the target. */
  details: Details | undefined;
}

export interface UpdateTargetRequest {
  /** unique identifier of the target. */
  targetId: string;
  /** Optionally change the unique name of the target. */
  name?: string | undefined;
  restWebhook?: SetRESTWebhook | undefined;
  restCall?: SetRESTCall | undefined;
  restAsync?:
    | SetRESTAsync
    | undefined;
  /** Optionally change the timeout, which defines the duration until ZITADEL cancels the execution. */
  timeout?: Duration | undefined;
  endpoint?: string | undefined;
}

export interface UpdateTargetResponse {
  /** Details provide some base information (such as the last change date) of the target. */
  details: Details | undefined;
}

export interface DeleteTargetRequest {
  /** unique identifier of the target. */
  targetId: string;
}

export interface DeleteTargetResponse {
  /** Details provide some base information (such as the last change date) of the target. */
  details: Details | undefined;
}

export interface ListTargetsRequest {
  /** list limitations and ordering. */
  query:
    | ListQuery
    | undefined;
  /** the field the result is sorted. */
  sortingColumn: TargetFieldName;
  /** Define the criteria to query for. */
  queries: TargetSearchQuery[];
}

export interface ListTargetsResponse {
  /** Details provides information about the returned result including total amount found. */
  details:
    | ListDetails
    | undefined;
  /** States by which field the results are sorted. */
  sortingColumn: TargetFieldName;
  /** The result contains the user schemas, which matched the queries. */
  result: Target[];
}

export interface GetTargetByIDRequest {
  /** unique identifier of the target. */
  targetId: string;
}

export interface GetTargetByIDResponse {
  target: Target | undefined;
}

export interface SetExecutionRequest {
  /** Defines the condition type and content of the condition for execution. */
  condition:
    | Condition
    | undefined;
  /** Ordered list of targets/includes called during the execution. */
  targets: ExecutionTargetType[];
}

export interface SetExecutionResponse {
  /** Details provide some base information (such as the last change date) of the execution. */
  details: Details | undefined;
}

export interface DeleteExecutionRequest {
  /** Unique identifier of the execution. */
  condition: Condition | undefined;
}

export interface DeleteExecutionResponse {
  /** Details provide some base information (such as the last change date) of the execution. */
  details: Details | undefined;
}

export interface ListExecutionsRequest {
  /** list limitations and ordering. */
  query:
    | ListQuery
    | undefined;
  /** Define the criteria to query for. */
  queries: SearchQuery[];
}

export interface ListExecutionsResponse {
  /** Details provides information about the returned result including total amount found. */
  details:
    | ListDetails
    | undefined;
  /** The result contains the executions, which matched the queries. */
  result: Execution[];
}

export interface ListExecutionFunctionsRequest {
}

export interface ListExecutionFunctionsResponse {
  /** All available methods */
  functions: string[];
}

export interface ListExecutionMethodsRequest {
}

export interface ListExecutionMethodsResponse {
  /** All available methods */
  methods: string[];
}

export interface ListExecutionServicesRequest {
}

export interface ListExecutionServicesResponse {
  /** All available methods */
  services: string[];
}

function createBaseCreateTargetRequest(): CreateTargetRequest {
  return {
    name: "",
    restWebhook: undefined,
    restCall: undefined,
    restAsync: undefined,
    timeout: undefined,
    endpoint: "",
  };
}

export const CreateTargetRequest = {
  encode(message: CreateTargetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.restWebhook !== undefined) {
      SetRESTWebhook.encode(message.restWebhook, writer.uint32(18).fork()).ldelim();
    }
    if (message.restCall !== undefined) {
      SetRESTCall.encode(message.restCall, writer.uint32(26).fork()).ldelim();
    }
    if (message.restAsync !== undefined) {
      SetRESTAsync.encode(message.restAsync, writer.uint32(34).fork()).ldelim();
    }
    if (message.timeout !== undefined) {
      Duration.encode(message.timeout, writer.uint32(42).fork()).ldelim();
    }
    if (message.endpoint !== "") {
      writer.uint32(50).string(message.endpoint);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateTargetRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTargetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.restWebhook = SetRESTWebhook.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.restCall = SetRESTCall.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.restAsync = SetRESTAsync.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.timeout = Duration.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.endpoint = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateTargetRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      restWebhook: isSet(object.restWebhook) ? SetRESTWebhook.fromJSON(object.restWebhook) : undefined,
      restCall: isSet(object.restCall) ? SetRESTCall.fromJSON(object.restCall) : undefined,
      restAsync: isSet(object.restAsync) ? SetRESTAsync.fromJSON(object.restAsync) : undefined,
      timeout: isSet(object.timeout) ? Duration.fromJSON(object.timeout) : undefined,
      endpoint: isSet(object.endpoint) ? String(object.endpoint) : "",
    };
  },

  toJSON(message: CreateTargetRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.restWebhook !== undefined &&
      (obj.restWebhook = message.restWebhook ? SetRESTWebhook.toJSON(message.restWebhook) : undefined);
    message.restCall !== undefined &&
      (obj.restCall = message.restCall ? SetRESTCall.toJSON(message.restCall) : undefined);
    message.restAsync !== undefined &&
      (obj.restAsync = message.restAsync ? SetRESTAsync.toJSON(message.restAsync) : undefined);
    message.timeout !== undefined && (obj.timeout = message.timeout ? Duration.toJSON(message.timeout) : undefined);
    message.endpoint !== undefined && (obj.endpoint = message.endpoint);
    return obj;
  },

  create(base?: DeepPartial<CreateTargetRequest>): CreateTargetRequest {
    return CreateTargetRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateTargetRequest>): CreateTargetRequest {
    const message = createBaseCreateTargetRequest();
    message.name = object.name ?? "";
    message.restWebhook = (object.restWebhook !== undefined && object.restWebhook !== null)
      ? SetRESTWebhook.fromPartial(object.restWebhook)
      : undefined;
    message.restCall = (object.restCall !== undefined && object.restCall !== null)
      ? SetRESTCall.fromPartial(object.restCall)
      : undefined;
    message.restAsync = (object.restAsync !== undefined && object.restAsync !== null)
      ? SetRESTAsync.fromPartial(object.restAsync)
      : undefined;
    message.timeout = (object.timeout !== undefined && object.timeout !== null)
      ? Duration.fromPartial(object.timeout)
      : undefined;
    message.endpoint = object.endpoint ?? "";
    return message;
  },
};

function createBaseCreateTargetResponse(): CreateTargetResponse {
  return { id: "", details: undefined };
}

export const CreateTargetResponse = {
  encode(message: CreateTargetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateTargetResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTargetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateTargetResponse {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
    };
  },

  toJSON(message: CreateTargetResponse): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<CreateTargetResponse>): CreateTargetResponse {
    return CreateTargetResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateTargetResponse>): CreateTargetResponse {
    const message = createBaseCreateTargetResponse();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateTargetRequest(): UpdateTargetRequest {
  return {
    targetId: "",
    name: undefined,
    restWebhook: undefined,
    restCall: undefined,
    restAsync: undefined,
    timeout: undefined,
    endpoint: undefined,
  };
}

export const UpdateTargetRequest = {
  encode(message: UpdateTargetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.targetId !== "") {
      writer.uint32(10).string(message.targetId);
    }
    if (message.name !== undefined) {
      writer.uint32(18).string(message.name);
    }
    if (message.restWebhook !== undefined) {
      SetRESTWebhook.encode(message.restWebhook, writer.uint32(26).fork()).ldelim();
    }
    if (message.restCall !== undefined) {
      SetRESTCall.encode(message.restCall, writer.uint32(34).fork()).ldelim();
    }
    if (message.restAsync !== undefined) {
      SetRESTAsync.encode(message.restAsync, writer.uint32(42).fork()).ldelim();
    }
    if (message.timeout !== undefined) {
      Duration.encode(message.timeout, writer.uint32(50).fork()).ldelim();
    }
    if (message.endpoint !== undefined) {
      writer.uint32(58).string(message.endpoint);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateTargetRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateTargetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.targetId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.restWebhook = SetRESTWebhook.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.restCall = SetRESTCall.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.restAsync = SetRESTAsync.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.timeout = Duration.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.endpoint = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateTargetRequest {
    return {
      targetId: isSet(object.targetId) ? String(object.targetId) : "",
      name: isSet(object.name) ? String(object.name) : undefined,
      restWebhook: isSet(object.restWebhook) ? SetRESTWebhook.fromJSON(object.restWebhook) : undefined,
      restCall: isSet(object.restCall) ? SetRESTCall.fromJSON(object.restCall) : undefined,
      restAsync: isSet(object.restAsync) ? SetRESTAsync.fromJSON(object.restAsync) : undefined,
      timeout: isSet(object.timeout) ? Duration.fromJSON(object.timeout) : undefined,
      endpoint: isSet(object.endpoint) ? String(object.endpoint) : undefined,
    };
  },

  toJSON(message: UpdateTargetRequest): unknown {
    const obj: any = {};
    message.targetId !== undefined && (obj.targetId = message.targetId);
    message.name !== undefined && (obj.name = message.name);
    message.restWebhook !== undefined &&
      (obj.restWebhook = message.restWebhook ? SetRESTWebhook.toJSON(message.restWebhook) : undefined);
    message.restCall !== undefined &&
      (obj.restCall = message.restCall ? SetRESTCall.toJSON(message.restCall) : undefined);
    message.restAsync !== undefined &&
      (obj.restAsync = message.restAsync ? SetRESTAsync.toJSON(message.restAsync) : undefined);
    message.timeout !== undefined && (obj.timeout = message.timeout ? Duration.toJSON(message.timeout) : undefined);
    message.endpoint !== undefined && (obj.endpoint = message.endpoint);
    return obj;
  },

  create(base?: DeepPartial<UpdateTargetRequest>): UpdateTargetRequest {
    return UpdateTargetRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateTargetRequest>): UpdateTargetRequest {
    const message = createBaseUpdateTargetRequest();
    message.targetId = object.targetId ?? "";
    message.name = object.name ?? undefined;
    message.restWebhook = (object.restWebhook !== undefined && object.restWebhook !== null)
      ? SetRESTWebhook.fromPartial(object.restWebhook)
      : undefined;
    message.restCall = (object.restCall !== undefined && object.restCall !== null)
      ? SetRESTCall.fromPartial(object.restCall)
      : undefined;
    message.restAsync = (object.restAsync !== undefined && object.restAsync !== null)
      ? SetRESTAsync.fromPartial(object.restAsync)
      : undefined;
    message.timeout = (object.timeout !== undefined && object.timeout !== null)
      ? Duration.fromPartial(object.timeout)
      : undefined;
    message.endpoint = object.endpoint ?? undefined;
    return message;
  },
};

function createBaseUpdateTargetResponse(): UpdateTargetResponse {
  return { details: undefined };
}

export const UpdateTargetResponse = {
  encode(message: UpdateTargetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateTargetResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateTargetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateTargetResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateTargetResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateTargetResponse>): UpdateTargetResponse {
    return UpdateTargetResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateTargetResponse>): UpdateTargetResponse {
    const message = createBaseUpdateTargetResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseDeleteTargetRequest(): DeleteTargetRequest {
  return { targetId: "" };
}

export const DeleteTargetRequest = {
  encode(message: DeleteTargetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.targetId !== "") {
      writer.uint32(10).string(message.targetId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteTargetRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteTargetRequest();
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

  fromJSON(object: any): DeleteTargetRequest {
    return { targetId: isSet(object.targetId) ? String(object.targetId) : "" };
  },

  toJSON(message: DeleteTargetRequest): unknown {
    const obj: any = {};
    message.targetId !== undefined && (obj.targetId = message.targetId);
    return obj;
  },

  create(base?: DeepPartial<DeleteTargetRequest>): DeleteTargetRequest {
    return DeleteTargetRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteTargetRequest>): DeleteTargetRequest {
    const message = createBaseDeleteTargetRequest();
    message.targetId = object.targetId ?? "";
    return message;
  },
};

function createBaseDeleteTargetResponse(): DeleteTargetResponse {
  return { details: undefined };
}

export const DeleteTargetResponse = {
  encode(message: DeleteTargetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteTargetResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteTargetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteTargetResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeleteTargetResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeleteTargetResponse>): DeleteTargetResponse {
    return DeleteTargetResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteTargetResponse>): DeleteTargetResponse {
    const message = createBaseDeleteTargetResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListTargetsRequest(): ListTargetsRequest {
  return { query: undefined, sortingColumn: 0, queries: [] };
}

export const ListTargetsRequest = {
  encode(message: ListTargetsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.queries) {
      TargetSearchQuery.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListTargetsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListTargetsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.query = ListQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.sortingColumn = reader.int32() as any;
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.queries.push(TargetSearchQuery.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListTargetsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? targetFieldNameFromJSON(object.sortingColumn) : 0,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => TargetSearchQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListTargetsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = targetFieldNameToJSON(message.sortingColumn));
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? TargetSearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListTargetsRequest>): ListTargetsRequest {
    return ListTargetsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListTargetsRequest>): ListTargetsRequest {
    const message = createBaseListTargetsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.queries = object.queries?.map((e) => TargetSearchQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListTargetsResponse(): ListTargetsResponse {
  return { details: undefined, sortingColumn: 0, result: [] };
}

export const ListTargetsResponse = {
  encode(message: ListTargetsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.result) {
      Target.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListTargetsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListTargetsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = ListDetails.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.sortingColumn = reader.int32() as any;
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.result.push(Target.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListTargetsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? targetFieldNameFromJSON(object.sortingColumn) : 0,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Target.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListTargetsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = targetFieldNameToJSON(message.sortingColumn));
    if (message.result) {
      obj.result = message.result.map((e) => e ? Target.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListTargetsResponse>): ListTargetsResponse {
    return ListTargetsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListTargetsResponse>): ListTargetsResponse {
    const message = createBaseListTargetsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.result = object.result?.map((e) => Target.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetTargetByIDRequest(): GetTargetByIDRequest {
  return { targetId: "" };
}

export const GetTargetByIDRequest = {
  encode(message: GetTargetByIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.targetId !== "") {
      writer.uint32(10).string(message.targetId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetTargetByIDRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetTargetByIDRequest();
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

  fromJSON(object: any): GetTargetByIDRequest {
    return { targetId: isSet(object.targetId) ? String(object.targetId) : "" };
  },

  toJSON(message: GetTargetByIDRequest): unknown {
    const obj: any = {};
    message.targetId !== undefined && (obj.targetId = message.targetId);
    return obj;
  },

  create(base?: DeepPartial<GetTargetByIDRequest>): GetTargetByIDRequest {
    return GetTargetByIDRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetTargetByIDRequest>): GetTargetByIDRequest {
    const message = createBaseGetTargetByIDRequest();
    message.targetId = object.targetId ?? "";
    return message;
  },
};

function createBaseGetTargetByIDResponse(): GetTargetByIDResponse {
  return { target: undefined };
}

export const GetTargetByIDResponse = {
  encode(message: GetTargetByIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.target !== undefined) {
      Target.encode(message.target, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetTargetByIDResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetTargetByIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.target = Target.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetTargetByIDResponse {
    return { target: isSet(object.target) ? Target.fromJSON(object.target) : undefined };
  },

  toJSON(message: GetTargetByIDResponse): unknown {
    const obj: any = {};
    message.target !== undefined && (obj.target = message.target ? Target.toJSON(message.target) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetTargetByIDResponse>): GetTargetByIDResponse {
    return GetTargetByIDResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetTargetByIDResponse>): GetTargetByIDResponse {
    const message = createBaseGetTargetByIDResponse();
    message.target = (object.target !== undefined && object.target !== null)
      ? Target.fromPartial(object.target)
      : undefined;
    return message;
  },
};

function createBaseSetExecutionRequest(): SetExecutionRequest {
  return { condition: undefined, targets: [] };
}

export const SetExecutionRequest = {
  encode(message: SetExecutionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.condition !== undefined) {
      Condition.encode(message.condition, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.targets) {
      ExecutionTargetType.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetExecutionRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetExecutionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.condition = Condition.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.targets.push(ExecutionTargetType.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetExecutionRequest {
    return {
      condition: isSet(object.condition) ? Condition.fromJSON(object.condition) : undefined,
      targets: Array.isArray(object?.targets) ? object.targets.map((e: any) => ExecutionTargetType.fromJSON(e)) : [],
    };
  },

  toJSON(message: SetExecutionRequest): unknown {
    const obj: any = {};
    message.condition !== undefined &&
      (obj.condition = message.condition ? Condition.toJSON(message.condition) : undefined);
    if (message.targets) {
      obj.targets = message.targets.map((e) => e ? ExecutionTargetType.toJSON(e) : undefined);
    } else {
      obj.targets = [];
    }
    return obj;
  },

  create(base?: DeepPartial<SetExecutionRequest>): SetExecutionRequest {
    return SetExecutionRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetExecutionRequest>): SetExecutionRequest {
    const message = createBaseSetExecutionRequest();
    message.condition = (object.condition !== undefined && object.condition !== null)
      ? Condition.fromPartial(object.condition)
      : undefined;
    message.targets = object.targets?.map((e) => ExecutionTargetType.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSetExecutionResponse(): SetExecutionResponse {
  return { details: undefined };
}

export const SetExecutionResponse = {
  encode(message: SetExecutionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetExecutionResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetExecutionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetExecutionResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetExecutionResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetExecutionResponse>): SetExecutionResponse {
    return SetExecutionResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetExecutionResponse>): SetExecutionResponse {
    const message = createBaseSetExecutionResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseDeleteExecutionRequest(): DeleteExecutionRequest {
  return { condition: undefined };
}

export const DeleteExecutionRequest = {
  encode(message: DeleteExecutionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.condition !== undefined) {
      Condition.encode(message.condition, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteExecutionRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteExecutionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.condition = Condition.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteExecutionRequest {
    return { condition: isSet(object.condition) ? Condition.fromJSON(object.condition) : undefined };
  },

  toJSON(message: DeleteExecutionRequest): unknown {
    const obj: any = {};
    message.condition !== undefined &&
      (obj.condition = message.condition ? Condition.toJSON(message.condition) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeleteExecutionRequest>): DeleteExecutionRequest {
    return DeleteExecutionRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteExecutionRequest>): DeleteExecutionRequest {
    const message = createBaseDeleteExecutionRequest();
    message.condition = (object.condition !== undefined && object.condition !== null)
      ? Condition.fromPartial(object.condition)
      : undefined;
    return message;
  },
};

function createBaseDeleteExecutionResponse(): DeleteExecutionResponse {
  return { details: undefined };
}

export const DeleteExecutionResponse = {
  encode(message: DeleteExecutionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteExecutionResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteExecutionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteExecutionResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeleteExecutionResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeleteExecutionResponse>): DeleteExecutionResponse {
    return DeleteExecutionResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteExecutionResponse>): DeleteExecutionResponse {
    const message = createBaseDeleteExecutionResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListExecutionsRequest(): ListExecutionsRequest {
  return { query: undefined, queries: [] };
}

export const ListExecutionsRequest = {
  encode(message: ListExecutionsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.queries) {
      SearchQuery.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListExecutionsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListExecutionsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.query = ListQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.queries.push(SearchQuery.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListExecutionsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => SearchQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListExecutionsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? SearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListExecutionsRequest>): ListExecutionsRequest {
    return ListExecutionsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListExecutionsRequest>): ListExecutionsRequest {
    const message = createBaseListExecutionsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.queries = object.queries?.map((e) => SearchQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListExecutionsResponse(): ListExecutionsResponse {
  return { details: undefined, result: [] };
}

export const ListExecutionsResponse = {
  encode(message: ListExecutionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      Execution.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListExecutionsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListExecutionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = ListDetails.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.result.push(Execution.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListExecutionsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Execution.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListExecutionsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? Execution.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListExecutionsResponse>): ListExecutionsResponse {
    return ListExecutionsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListExecutionsResponse>): ListExecutionsResponse {
    const message = createBaseListExecutionsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => Execution.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListExecutionFunctionsRequest(): ListExecutionFunctionsRequest {
  return {};
}

export const ListExecutionFunctionsRequest = {
  encode(_: ListExecutionFunctionsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListExecutionFunctionsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListExecutionFunctionsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): ListExecutionFunctionsRequest {
    return {};
  },

  toJSON(_: ListExecutionFunctionsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListExecutionFunctionsRequest>): ListExecutionFunctionsRequest {
    return ListExecutionFunctionsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListExecutionFunctionsRequest>): ListExecutionFunctionsRequest {
    const message = createBaseListExecutionFunctionsRequest();
    return message;
  },
};

function createBaseListExecutionFunctionsResponse(): ListExecutionFunctionsResponse {
  return { functions: [] };
}

export const ListExecutionFunctionsResponse = {
  encode(message: ListExecutionFunctionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.functions) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListExecutionFunctionsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListExecutionFunctionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.functions.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListExecutionFunctionsResponse {
    return { functions: Array.isArray(object?.functions) ? object.functions.map((e: any) => String(e)) : [] };
  },

  toJSON(message: ListExecutionFunctionsResponse): unknown {
    const obj: any = {};
    if (message.functions) {
      obj.functions = message.functions.map((e) => e);
    } else {
      obj.functions = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListExecutionFunctionsResponse>): ListExecutionFunctionsResponse {
    return ListExecutionFunctionsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListExecutionFunctionsResponse>): ListExecutionFunctionsResponse {
    const message = createBaseListExecutionFunctionsResponse();
    message.functions = object.functions?.map((e) => e) || [];
    return message;
  },
};

function createBaseListExecutionMethodsRequest(): ListExecutionMethodsRequest {
  return {};
}

export const ListExecutionMethodsRequest = {
  encode(_: ListExecutionMethodsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListExecutionMethodsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListExecutionMethodsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): ListExecutionMethodsRequest {
    return {};
  },

  toJSON(_: ListExecutionMethodsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListExecutionMethodsRequest>): ListExecutionMethodsRequest {
    return ListExecutionMethodsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListExecutionMethodsRequest>): ListExecutionMethodsRequest {
    const message = createBaseListExecutionMethodsRequest();
    return message;
  },
};

function createBaseListExecutionMethodsResponse(): ListExecutionMethodsResponse {
  return { methods: [] };
}

export const ListExecutionMethodsResponse = {
  encode(message: ListExecutionMethodsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.methods) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListExecutionMethodsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListExecutionMethodsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.methods.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListExecutionMethodsResponse {
    return { methods: Array.isArray(object?.methods) ? object.methods.map((e: any) => String(e)) : [] };
  },

  toJSON(message: ListExecutionMethodsResponse): unknown {
    const obj: any = {};
    if (message.methods) {
      obj.methods = message.methods.map((e) => e);
    } else {
      obj.methods = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListExecutionMethodsResponse>): ListExecutionMethodsResponse {
    return ListExecutionMethodsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListExecutionMethodsResponse>): ListExecutionMethodsResponse {
    const message = createBaseListExecutionMethodsResponse();
    message.methods = object.methods?.map((e) => e) || [];
    return message;
  },
};

function createBaseListExecutionServicesRequest(): ListExecutionServicesRequest {
  return {};
}

export const ListExecutionServicesRequest = {
  encode(_: ListExecutionServicesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListExecutionServicesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListExecutionServicesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): ListExecutionServicesRequest {
    return {};
  },

  toJSON(_: ListExecutionServicesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListExecutionServicesRequest>): ListExecutionServicesRequest {
    return ListExecutionServicesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListExecutionServicesRequest>): ListExecutionServicesRequest {
    const message = createBaseListExecutionServicesRequest();
    return message;
  },
};

function createBaseListExecutionServicesResponse(): ListExecutionServicesResponse {
  return { services: [] };
}

export const ListExecutionServicesResponse = {
  encode(message: ListExecutionServicesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.services) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListExecutionServicesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListExecutionServicesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.services.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListExecutionServicesResponse {
    return { services: Array.isArray(object?.services) ? object.services.map((e: any) => String(e)) : [] };
  },

  toJSON(message: ListExecutionServicesResponse): unknown {
    const obj: any = {};
    if (message.services) {
      obj.services = message.services.map((e) => e);
    } else {
      obj.services = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListExecutionServicesResponse>): ListExecutionServicesResponse {
    return ListExecutionServicesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListExecutionServicesResponse>): ListExecutionServicesResponse {
    const message = createBaseListExecutionServicesResponse();
    message.services = object.services?.map((e) => e) || [];
    return message;
  },
};

export type ActionServiceDefinition = typeof ActionServiceDefinition;
export const ActionServiceDefinition = {
  name: "ActionService",
  fullName: "zitadel.action.v3alpha.ActionService",
  methods: {
    /**
     * Create a target
     *
     * Create a new target, which can be used in executions.
     */
    createTarget: {
      name: "CreateTarget",
      requestType: CreateTargetRequest,
      requestStream: false,
      responseType: CreateTargetResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              85,
              74,
              83,
              10,
              3,
              50,
              48,
              49,
              18,
              76,
              10,
              27,
              84,
              97,
              114,
              103,
              101,
              116,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              99,
              114,
              101,
              97,
              116,
              101,
              100,
              18,
              45,
              10,
              43,
              26,
              41,
              35,
              47,
              100,
              101,
              102,
              105,
              110,
              105,
              116,
              105,
              111,
              110,
              115,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              67,
              114,
              101,
              97,
              116,
              101,
              84,
              97,
              114,
              103,
              101,
              116,
              82,
              101,
              115,
              112,
              111,
              110,
              115,
              101,
            ]),
          ],
          400010: [
            Buffer.from([
              31,
              10,
              24,
              10,
              22,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              46,
              116,
              97,
              114,
              103,
              101,
              116,
              46,
              119,
              114,
              105,
              116,
              101,
              18,
              3,
              8,
              201,
              1,
            ]),
          ],
          578365826: [
            Buffer.from([
              21,
              58,
              1,
              42,
              34,
              16,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              116,
              97,
              114,
              103,
              101,
              116,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * Update a target
     *
     * Update an existing target.
     */
    updateTarget: {
      name: "UpdateTarget",
      requestType: UpdateTargetRequest,
      requestStream: false,
      responseType: UpdateTargetResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              38,
              74,
              36,
              10,
              3,
              50,
              48,
              48,
              18,
              29,
              10,
              27,
              84,
              97,
              114,
              103,
              101,
              116,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              117,
              112,
              100,
              97,
              116,
              101,
              100,
            ]),
          ],
          400010: [
            Buffer.from([
              26,
              10,
              24,
              10,
              22,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              46,
              116,
              97,
              114,
              103,
              101,
              116,
              46,
              119,
              114,
              105,
              116,
              101,
            ]),
          ],
          578365826: [
            Buffer.from([
              33,
              58,
              1,
              42,
              26,
              28,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              116,
              97,
              114,
              103,
              101,
              116,
              115,
              47,
              123,
              116,
              97,
              114,
              103,
              101,
              116,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Delete a target
     *
     * Delete an existing target. This will remove it from any configured execution as well.
     */
    deleteTarget: {
      name: "DeleteTarget",
      requestType: DeleteTargetRequest,
      requestStream: false,
      responseType: DeleteTargetResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              38,
              74,
              36,
              10,
              3,
              50,
              48,
              48,
              18,
              29,
              10,
              27,
              84,
              97,
              114,
              103,
              101,
              116,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              100,
              101,
              108,
              101,
              116,
              101,
              100,
            ]),
          ],
          400010: [
            Buffer.from([
              27,
              10,
              25,
              10,
              23,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              46,
              116,
              97,
              114,
              103,
              101,
              116,
              46,
              100,
              101,
              108,
              101,
              116,
              101,
            ]),
          ],
          578365826: [
            Buffer.from([
              30,
              42,
              28,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              116,
              97,
              114,
              103,
              101,
              116,
              115,
              47,
              123,
              116,
              97,
              114,
              103,
              101,
              116,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * List targets
     *
     * List all matching targets. By default, we will return all targets of your instance.
     * Make sure to include a limit and sorting for pagination.
     */
    listTargets: {
      name: "ListTargets",
      requestType: ListTargetsRequest,
      requestStream: false,
      responseType: ListTargetsResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              109,
              74,
              49,
              10,
              3,
              50,
              48,
              48,
              18,
              42,
              10,
              40,
              65,
              32,
              108,
              105,
              115,
              116,
              32,
              111,
              102,
              32,
              97,
              108,
              108,
              32,
              116,
              97,
              114,
              103,
              101,
              116,
              115,
              32,
              109,
              97,
              116,
              99,
              104,
              105,
              110,
              103,
              32,
              116,
              104,
              101,
              32,
              113,
              117,
              101,
              114,
              121,
              74,
              56,
              10,
              3,
              52,
              48,
              48,
              18,
              49,
              10,
              18,
              105,
              110,
              118,
              97,
              108,
              105,
              100,
              32,
              108,
              105,
              115,
              116,
              32,
              113,
              117,
              101,
              114,
              121,
              18,
              27,
              10,
              25,
              26,
              23,
              35,
              47,
              100,
              101,
              102,
              105,
              110,
              105,
              116,
              105,
              111,
              110,
              115,
              47,
              114,
              112,
              99,
              83,
              116,
              97,
              116,
              117,
              115,
            ]),
          ],
          400010: [
            Buffer.from([
              25,
              10,
              23,
              10,
              21,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              46,
              116,
              97,
              114,
              103,
              101,
              116,
              46,
              114,
              101,
              97,
              100,
            ]),
          ],
          578365826: [
            Buffer.from([
              28,
              58,
              1,
              42,
              34,
              23,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              116,
              97,
              114,
              103,
              101,
              116,
              115,
              47,
              115,
              101,
              97,
              114,
              99,
              104,
            ]),
          ],
        },
      },
    },
    /**
     * Target by ID
     *
     * Returns the target identified by the requested ID.
     */
    getTargetByID: {
      name: "GetTargetByID",
      requestType: GetTargetByIDRequest,
      requestStream: false,
      responseType: GetTargetByIDResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              40,
              74,
              38,
              10,
              3,
              50,
              48,
              48,
              18,
              31,
              10,
              29,
              84,
              97,
              114,
              103,
              101,
              116,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              116,
              114,
              105,
              101,
              118,
              101,
              100,
            ]),
          ],
          400010: [
            Buffer.from([
              25,
              10,
              23,
              10,
              21,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              46,
              116,
              97,
              114,
              103,
              101,
              116,
              46,
              114,
              101,
              97,
              100,
            ]),
          ],
          578365826: [
            Buffer.from([
              30,
              18,
              28,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              116,
              97,
              114,
              103,
              101,
              116,
              115,
              47,
              123,
              116,
              97,
              114,
              103,
              101,
              116,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Set an execution
     *
     * Set an execution to call a previously defined target or include the targets of a previously defined execution.
     */
    setExecution: {
      name: "SetExecution",
      requestType: SetExecutionRequest,
      requestStream: false,
      responseType: SetExecutionResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              37,
              74,
              35,
              10,
              3,
              50,
              48,
              48,
              18,
              28,
              10,
              26,
              69,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              115,
              101,
              116,
            ]),
          ],
          400010: [
            Buffer.from([19, 10, 17, 10, 15, 101, 120, 101, 99, 117, 116, 105, 111, 110, 46, 119, 114, 105, 116, 101]),
          ],
          578365826: [
            Buffer.from([
              24,
              58,
              1,
              42,
              26,
              19,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * Delete an execution
     *
     * Delete an existing execution.
     */
    deleteExecution: {
      name: "DeleteExecution",
      requestType: DeleteExecutionRequest,
      requestStream: false,
      responseType: DeleteExecutionResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              41,
              74,
              39,
              10,
              3,
              50,
              48,
              48,
              18,
              32,
              10,
              30,
              69,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              100,
              101,
              108,
              101,
              116,
              101,
              100,
            ]),
          ],
          400010: [
            Buffer.from([
              20,
              10,
              18,
              10,
              16,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              46,
              100,
              101,
              108,
              101,
              116,
              101,
            ]),
          ],
          578365826: [
            Buffer.from([
              24,
              58,
              1,
              42,
              42,
              19,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * List executions
     *
     * List all matching executions. By default, we will return all executions of your instance.
     * Make sure to include a limit and sorting for pagination.
     */
    listExecutions: {
      name: "ListExecutions",
      requestType: ListExecutionsRequest,
      requestStream: false,
      responseType: ListExecutionsResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              112,
              74,
              52,
              10,
              3,
              50,
              48,
              48,
              18,
              45,
              10,
              43,
              65,
              32,
              108,
              105,
              115,
              116,
              32,
              111,
              102,
              32,
              97,
              108,
              108,
              32,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              115,
              32,
              109,
              97,
              116,
              99,
              104,
              105,
              110,
              103,
              32,
              116,
              104,
              101,
              32,
              113,
              117,
              101,
              114,
              121,
              74,
              56,
              10,
              3,
              52,
              48,
              48,
              18,
              49,
              10,
              18,
              105,
              110,
              118,
              97,
              108,
              105,
              100,
              32,
              108,
              105,
              115,
              116,
              32,
              113,
              117,
              101,
              114,
              121,
              18,
              27,
              10,
              25,
              26,
              23,
              35,
              47,
              100,
              101,
              102,
              105,
              110,
              105,
              116,
              105,
              111,
              110,
              115,
              47,
              114,
              112,
              99,
              83,
              116,
              97,
              116,
              117,
              115,
            ]),
          ],
          400010: [
            Buffer.from([18, 10, 16, 10, 14, 101, 120, 101, 99, 117, 116, 105, 111, 110, 46, 114, 101, 97, 100]),
          ],
          578365826: [
            Buffer.from([
              31,
              58,
              1,
              42,
              34,
              26,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              115,
              47,
              115,
              101,
              97,
              114,
              99,
              104,
            ]),
          ],
        },
      },
    },
    /**
     * List all available functions
     *
     * List all available functions which can be used as condition for executions.
     */
    listExecutionFunctions: {
      name: "ListExecutionFunctions",
      requestType: ListExecutionFunctionsRequest,
      requestStream: false,
      responseType: ListExecutionFunctionsResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              42,
              74,
              40,
              10,
              3,
              50,
              48,
              48,
              18,
              33,
              10,
              31,
              76,
              105,
              115,
              116,
              32,
              97,
              108,
              108,
              32,
              102,
              117,
              110,
              99,
              116,
              105,
              111,
              110,
              115,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
            ]),
          ],
          400010: [
            Buffer.from([18, 10, 16, 10, 14, 101, 120, 101, 99, 117, 116, 105, 111, 110, 46, 114, 101, 97, 100]),
          ],
          578365826: [
            Buffer.from([
              31,
              18,
              29,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              115,
              47,
              102,
              117,
              110,
              99,
              116,
              105,
              111,
              110,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * List all available methods
     *
     * List all available methods which can be used as condition for executions.
     */
    listExecutionMethods: {
      name: "ListExecutionMethods",
      requestType: ListExecutionMethodsRequest,
      requestStream: false,
      responseType: ListExecutionMethodsResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              40,
              74,
              38,
              10,
              3,
              50,
              48,
              48,
              18,
              31,
              10,
              29,
              76,
              105,
              115,
              116,
              32,
              97,
              108,
              108,
              32,
              109,
              101,
              116,
              104,
              111,
              100,
              115,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
            ]),
          ],
          400010: [
            Buffer.from([18, 10, 16, 10, 14, 101, 120, 101, 99, 117, 116, 105, 111, 110, 46, 114, 101, 97, 100]),
          ],
          578365826: [
            Buffer.from([
              29,
              18,
              27,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              115,
              47,
              109,
              101,
              116,
              104,
              111,
              100,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * List all available service
     *
     * List all available services which can be used as condition for executions.
     */
    listExecutionServices: {
      name: "ListExecutionServices",
      requestType: ListExecutionServicesRequest,
      requestStream: false,
      responseType: ListExecutionServicesResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              41,
              74,
              39,
              10,
              3,
              50,
              48,
              48,
              18,
              32,
              10,
              30,
              76,
              105,
              115,
              116,
              32,
              97,
              108,
              108,
              32,
              115,
              101,
              114,
              118,
              105,
              99,
              101,
              115,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
            ]),
          ],
          400010: [
            Buffer.from([18, 10, 16, 10, 14, 101, 120, 101, 99, 117, 116, 105, 111, 110, 46, 114, 101, 97, 100]),
          ],
          578365826: [
            Buffer.from([
              30,
              18,
              28,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              101,
              120,
              101,
              99,
              117,
              116,
              105,
              111,
              110,
              115,
              47,
              115,
              101,
              114,
              118,
              105,
              99,
              101,
              115,
            ]),
          ],
        },
      },
    },
  },
} as const;

export interface ActionServiceImplementation<CallContextExt = {}> {
  /**
   * Create a target
   *
   * Create a new target, which can be used in executions.
   */
  createTarget(
    request: CreateTargetRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<CreateTargetResponse>>;
  /**
   * Update a target
   *
   * Update an existing target.
   */
  updateTarget(
    request: UpdateTargetRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateTargetResponse>>;
  /**
   * Delete a target
   *
   * Delete an existing target. This will remove it from any configured execution as well.
   */
  deleteTarget(
    request: DeleteTargetRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeleteTargetResponse>>;
  /**
   * List targets
   *
   * List all matching targets. By default, we will return all targets of your instance.
   * Make sure to include a limit and sorting for pagination.
   */
  listTargets(
    request: ListTargetsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListTargetsResponse>>;
  /**
   * Target by ID
   *
   * Returns the target identified by the requested ID.
   */
  getTargetByID(
    request: GetTargetByIDRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetTargetByIDResponse>>;
  /**
   * Set an execution
   *
   * Set an execution to call a previously defined target or include the targets of a previously defined execution.
   */
  setExecution(
    request: SetExecutionRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetExecutionResponse>>;
  /**
   * Delete an execution
   *
   * Delete an existing execution.
   */
  deleteExecution(
    request: DeleteExecutionRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeleteExecutionResponse>>;
  /**
   * List executions
   *
   * List all matching executions. By default, we will return all executions of your instance.
   * Make sure to include a limit and sorting for pagination.
   */
  listExecutions(
    request: ListExecutionsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListExecutionsResponse>>;
  /**
   * List all available functions
   *
   * List all available functions which can be used as condition for executions.
   */
  listExecutionFunctions(
    request: ListExecutionFunctionsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListExecutionFunctionsResponse>>;
  /**
   * List all available methods
   *
   * List all available methods which can be used as condition for executions.
   */
  listExecutionMethods(
    request: ListExecutionMethodsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListExecutionMethodsResponse>>;
  /**
   * List all available service
   *
   * List all available services which can be used as condition for executions.
   */
  listExecutionServices(
    request: ListExecutionServicesRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListExecutionServicesResponse>>;
}

export interface ActionServiceClient<CallOptionsExt = {}> {
  /**
   * Create a target
   *
   * Create a new target, which can be used in executions.
   */
  createTarget(
    request: DeepPartial<CreateTargetRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<CreateTargetResponse>;
  /**
   * Update a target
   *
   * Update an existing target.
   */
  updateTarget(
    request: DeepPartial<UpdateTargetRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateTargetResponse>;
  /**
   * Delete a target
   *
   * Delete an existing target. This will remove it from any configured execution as well.
   */
  deleteTarget(
    request: DeepPartial<DeleteTargetRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeleteTargetResponse>;
  /**
   * List targets
   *
   * List all matching targets. By default, we will return all targets of your instance.
   * Make sure to include a limit and sorting for pagination.
   */
  listTargets(
    request: DeepPartial<ListTargetsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListTargetsResponse>;
  /**
   * Target by ID
   *
   * Returns the target identified by the requested ID.
   */
  getTargetByID(
    request: DeepPartial<GetTargetByIDRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetTargetByIDResponse>;
  /**
   * Set an execution
   *
   * Set an execution to call a previously defined target or include the targets of a previously defined execution.
   */
  setExecution(
    request: DeepPartial<SetExecutionRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetExecutionResponse>;
  /**
   * Delete an execution
   *
   * Delete an existing execution.
   */
  deleteExecution(
    request: DeepPartial<DeleteExecutionRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeleteExecutionResponse>;
  /**
   * List executions
   *
   * List all matching executions. By default, we will return all executions of your instance.
   * Make sure to include a limit and sorting for pagination.
   */
  listExecutions(
    request: DeepPartial<ListExecutionsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListExecutionsResponse>;
  /**
   * List all available functions
   *
   * List all available functions which can be used as condition for executions.
   */
  listExecutionFunctions(
    request: DeepPartial<ListExecutionFunctionsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListExecutionFunctionsResponse>;
  /**
   * List all available methods
   *
   * List all available methods which can be used as condition for executions.
   */
  listExecutionMethods(
    request: DeepPartial<ListExecutionMethodsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListExecutionMethodsResponse>;
  /**
   * List all available service
   *
   * List all available services which can be used as condition for executions.
   */
  listExecutionServices(
    request: DeepPartial<ListExecutionServicesRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListExecutionServicesResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

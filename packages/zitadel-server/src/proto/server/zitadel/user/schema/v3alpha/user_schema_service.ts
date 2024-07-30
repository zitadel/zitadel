/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Struct } from "../../../../google/protobuf/struct";
import { Details, ListDetails, ListQuery } from "../../../object/v2beta/object";
import {
  AuthenticatorType,
  authenticatorTypeFromJSON,
  authenticatorTypeToJSON,
  FieldName,
  fieldNameFromJSON,
  fieldNameToJSON,
  SearchQuery,
  UserSchema,
} from "./user_schema";

export const protobufPackage = "zitadel.user.schema.v3alpha";

export interface ListUserSchemasRequest {
  /** list limitations and ordering. */
  query:
    | ListQuery
    | undefined;
  /** the field the result is sorted. */
  sortingColumn: FieldName;
  /** Define the criteria to query for. */
  queries: SearchQuery[];
}

export interface ListUserSchemasResponse {
  /** Details provides information about the returned result including total amount found. */
  details:
    | ListDetails
    | undefined;
  /** States by which field the results are sorted. */
  sortingColumn: FieldName;
  /** The result contains the user schemas, which matched the queries. */
  result: UserSchema[];
}

export interface GetUserSchemaByIDRequest {
  /** unique identifier of the schema. */
  id: string;
}

export interface GetUserSchemaByIDResponse {
  schema: UserSchema | undefined;
}

export interface CreateUserSchemaRequest {
  /** Type is a human readable word describing the schema. */
  type: string;
  /** JSON schema representation defining the user. */
  schema?:
    | { [key: string]: any }
    | undefined;
  /** Defines the possible types of authenticators. */
  possibleAuthenticators: AuthenticatorType[];
}

export interface CreateUserSchemaResponse {
  /** ID is the read-only unique identifier of the schema. */
  id: string;
  /** Details provide some base information (such as the last change date) of the schema. */
  details: Details | undefined;
}

export interface UpdateUserSchemaRequest {
  /** unique identifier of the schema. */
  id: string;
  /** Type is a human readable word describing the schema. */
  type?:
    | string
    | undefined;
  /** JSON schema representation defining the user. */
  schema?:
    | { [key: string]: any }
    | undefined;
  /**
   * Defines the possible types of authenticators.
   *
   * Removal of an authenticator does not remove the authenticator on a user.
   */
  possibleAuthenticators: AuthenticatorType[];
}

export interface UpdateUserSchemaResponse {
  /** Details provide some base information (such as the last change date) of the schema. */
  details: Details | undefined;
}

export interface DeactivateUserSchemaRequest {
  /** unique identifier of the schema. */
  id: string;
}

export interface DeactivateUserSchemaResponse {
  /** Details provide some base information (such as the last change date) of the schema. */
  details: Details | undefined;
}

export interface ReactivateUserSchemaRequest {
  /** unique identifier of the schema. */
  id: string;
}

export interface ReactivateUserSchemaResponse {
  /** Details provide some base information (such as the last change date) of the schema. */
  details: Details | undefined;
}

export interface DeleteUserSchemaRequest {
  /** unique identifier of the schema. */
  id: string;
}

export interface DeleteUserSchemaResponse {
  /** Details provide some base information (such as the last change date) of the schema. */
  details: Details | undefined;
}

function createBaseListUserSchemasRequest(): ListUserSchemasRequest {
  return { query: undefined, sortingColumn: 0, queries: [] };
}

export const ListUserSchemasRequest = {
  encode(message: ListUserSchemasRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.queries) {
      SearchQuery.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListUserSchemasRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListUserSchemasRequest();
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

  fromJSON(object: any): ListUserSchemasRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? fieldNameFromJSON(object.sortingColumn) : 0,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => SearchQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListUserSchemasRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = fieldNameToJSON(message.sortingColumn));
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? SearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListUserSchemasRequest>): ListUserSchemasRequest {
    return ListUserSchemasRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListUserSchemasRequest>): ListUserSchemasRequest {
    const message = createBaseListUserSchemasRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.queries = object.queries?.map((e) => SearchQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListUserSchemasResponse(): ListUserSchemasResponse {
  return { details: undefined, sortingColumn: 0, result: [] };
}

export const ListUserSchemasResponse = {
  encode(message: ListUserSchemasResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.result) {
      UserSchema.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListUserSchemasResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListUserSchemasResponse();
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

          message.result.push(UserSchema.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListUserSchemasResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? fieldNameFromJSON(object.sortingColumn) : 0,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => UserSchema.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListUserSchemasResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = fieldNameToJSON(message.sortingColumn));
    if (message.result) {
      obj.result = message.result.map((e) => e ? UserSchema.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListUserSchemasResponse>): ListUserSchemasResponse {
    return ListUserSchemasResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListUserSchemasResponse>): ListUserSchemasResponse {
    const message = createBaseListUserSchemasResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.result = object.result?.map((e) => UserSchema.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetUserSchemaByIDRequest(): GetUserSchemaByIDRequest {
  return { id: "" };
}

export const GetUserSchemaByIDRequest = {
  encode(message: GetUserSchemaByIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserSchemaByIDRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserSchemaByIDRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetUserSchemaByIDRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: GetUserSchemaByIDRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<GetUserSchemaByIDRequest>): GetUserSchemaByIDRequest {
    return GetUserSchemaByIDRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetUserSchemaByIDRequest>): GetUserSchemaByIDRequest {
    const message = createBaseGetUserSchemaByIDRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseGetUserSchemaByIDResponse(): GetUserSchemaByIDResponse {
  return { schema: undefined };
}

export const GetUserSchemaByIDResponse = {
  encode(message: GetUserSchemaByIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.schema !== undefined) {
      UserSchema.encode(message.schema, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserSchemaByIDResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserSchemaByIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.schema = UserSchema.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetUserSchemaByIDResponse {
    return { schema: isSet(object.schema) ? UserSchema.fromJSON(object.schema) : undefined };
  },

  toJSON(message: GetUserSchemaByIDResponse): unknown {
    const obj: any = {};
    message.schema !== undefined && (obj.schema = message.schema ? UserSchema.toJSON(message.schema) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetUserSchemaByIDResponse>): GetUserSchemaByIDResponse {
    return GetUserSchemaByIDResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetUserSchemaByIDResponse>): GetUserSchemaByIDResponse {
    const message = createBaseGetUserSchemaByIDResponse();
    message.schema = (object.schema !== undefined && object.schema !== null)
      ? UserSchema.fromPartial(object.schema)
      : undefined;
    return message;
  },
};

function createBaseCreateUserSchemaRequest(): CreateUserSchemaRequest {
  return { type: "", schema: undefined, possibleAuthenticators: [] };
}

export const CreateUserSchemaRequest = {
  encode(message: CreateUserSchemaRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.schema !== undefined) {
      Struct.encode(Struct.wrap(message.schema), writer.uint32(18).fork()).ldelim();
    }
    writer.uint32(26).fork();
    for (const v of message.possibleAuthenticators) {
      writer.int32(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateUserSchemaRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateUserSchemaRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.type = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.schema = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag == 24) {
            message.possibleAuthenticators.push(reader.int32() as any);
            continue;
          }

          if (tag == 26) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.possibleAuthenticators.push(reader.int32() as any);
            }

            continue;
          }

          break;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateUserSchemaRequest {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      schema: isObject(object.schema) ? object.schema : undefined,
      possibleAuthenticators: Array.isArray(object?.possibleAuthenticators)
        ? object.possibleAuthenticators.map((e: any) => authenticatorTypeFromJSON(e))
        : [],
    };
  },

  toJSON(message: CreateUserSchemaRequest): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.schema !== undefined && (obj.schema = message.schema);
    if (message.possibleAuthenticators) {
      obj.possibleAuthenticators = message.possibleAuthenticators.map((e) => authenticatorTypeToJSON(e));
    } else {
      obj.possibleAuthenticators = [];
    }
    return obj;
  },

  create(base?: DeepPartial<CreateUserSchemaRequest>): CreateUserSchemaRequest {
    return CreateUserSchemaRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateUserSchemaRequest>): CreateUserSchemaRequest {
    const message = createBaseCreateUserSchemaRequest();
    message.type = object.type ?? "";
    message.schema = object.schema ?? undefined;
    message.possibleAuthenticators = object.possibleAuthenticators?.map((e) => e) || [];
    return message;
  },
};

function createBaseCreateUserSchemaResponse(): CreateUserSchemaResponse {
  return { id: "", details: undefined };
}

export const CreateUserSchemaResponse = {
  encode(message: CreateUserSchemaResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateUserSchemaResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateUserSchemaResponse();
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

  fromJSON(object: any): CreateUserSchemaResponse {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
    };
  },

  toJSON(message: CreateUserSchemaResponse): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<CreateUserSchemaResponse>): CreateUserSchemaResponse {
    return CreateUserSchemaResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateUserSchemaResponse>): CreateUserSchemaResponse {
    const message = createBaseCreateUserSchemaResponse();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUpdateUserSchemaRequest(): UpdateUserSchemaRequest {
  return { id: "", type: undefined, schema: undefined, possibleAuthenticators: [] };
}

export const UpdateUserSchemaRequest = {
  encode(message: UpdateUserSchemaRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.type !== undefined) {
      writer.uint32(18).string(message.type);
    }
    if (message.schema !== undefined) {
      Struct.encode(Struct.wrap(message.schema), writer.uint32(26).fork()).ldelim();
    }
    writer.uint32(34).fork();
    for (const v of message.possibleAuthenticators) {
      writer.int32(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateUserSchemaRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateUserSchemaRequest();
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

          message.type = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.schema = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag == 32) {
            message.possibleAuthenticators.push(reader.int32() as any);
            continue;
          }

          if (tag == 34) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.possibleAuthenticators.push(reader.int32() as any);
            }

            continue;
          }

          break;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateUserSchemaRequest {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      type: isSet(object.type) ? String(object.type) : undefined,
      schema: isObject(object.schema) ? object.schema : undefined,
      possibleAuthenticators: Array.isArray(object?.possibleAuthenticators)
        ? object.possibleAuthenticators.map((e: any) => authenticatorTypeFromJSON(e))
        : [],
    };
  },

  toJSON(message: UpdateUserSchemaRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.type !== undefined && (obj.type = message.type);
    message.schema !== undefined && (obj.schema = message.schema);
    if (message.possibleAuthenticators) {
      obj.possibleAuthenticators = message.possibleAuthenticators.map((e) => authenticatorTypeToJSON(e));
    } else {
      obj.possibleAuthenticators = [];
    }
    return obj;
  },

  create(base?: DeepPartial<UpdateUserSchemaRequest>): UpdateUserSchemaRequest {
    return UpdateUserSchemaRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateUserSchemaRequest>): UpdateUserSchemaRequest {
    const message = createBaseUpdateUserSchemaRequest();
    message.id = object.id ?? "";
    message.type = object.type ?? undefined;
    message.schema = object.schema ?? undefined;
    message.possibleAuthenticators = object.possibleAuthenticators?.map((e) => e) || [];
    return message;
  },
};

function createBaseUpdateUserSchemaResponse(): UpdateUserSchemaResponse {
  return { details: undefined };
}

export const UpdateUserSchemaResponse = {
  encode(message: UpdateUserSchemaResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateUserSchemaResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateUserSchemaResponse();
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

  fromJSON(object: any): UpdateUserSchemaResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateUserSchemaResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateUserSchemaResponse>): UpdateUserSchemaResponse {
    return UpdateUserSchemaResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateUserSchemaResponse>): UpdateUserSchemaResponse {
    const message = createBaseUpdateUserSchemaResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseDeactivateUserSchemaRequest(): DeactivateUserSchemaRequest {
  return { id: "" };
}

export const DeactivateUserSchemaRequest = {
  encode(message: DeactivateUserSchemaRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeactivateUserSchemaRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeactivateUserSchemaRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeactivateUserSchemaRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: DeactivateUserSchemaRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<DeactivateUserSchemaRequest>): DeactivateUserSchemaRequest {
    return DeactivateUserSchemaRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeactivateUserSchemaRequest>): DeactivateUserSchemaRequest {
    const message = createBaseDeactivateUserSchemaRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseDeactivateUserSchemaResponse(): DeactivateUserSchemaResponse {
  return { details: undefined };
}

export const DeactivateUserSchemaResponse = {
  encode(message: DeactivateUserSchemaResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeactivateUserSchemaResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeactivateUserSchemaResponse();
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

  fromJSON(object: any): DeactivateUserSchemaResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeactivateUserSchemaResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeactivateUserSchemaResponse>): DeactivateUserSchemaResponse {
    return DeactivateUserSchemaResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeactivateUserSchemaResponse>): DeactivateUserSchemaResponse {
    const message = createBaseDeactivateUserSchemaResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseReactivateUserSchemaRequest(): ReactivateUserSchemaRequest {
  return { id: "" };
}

export const ReactivateUserSchemaRequest = {
  encode(message: ReactivateUserSchemaRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReactivateUserSchemaRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReactivateUserSchemaRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReactivateUserSchemaRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: ReactivateUserSchemaRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<ReactivateUserSchemaRequest>): ReactivateUserSchemaRequest {
    return ReactivateUserSchemaRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ReactivateUserSchemaRequest>): ReactivateUserSchemaRequest {
    const message = createBaseReactivateUserSchemaRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseReactivateUserSchemaResponse(): ReactivateUserSchemaResponse {
  return { details: undefined };
}

export const ReactivateUserSchemaResponse = {
  encode(message: ReactivateUserSchemaResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReactivateUserSchemaResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReactivateUserSchemaResponse();
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

  fromJSON(object: any): ReactivateUserSchemaResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: ReactivateUserSchemaResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ReactivateUserSchemaResponse>): ReactivateUserSchemaResponse {
    return ReactivateUserSchemaResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ReactivateUserSchemaResponse>): ReactivateUserSchemaResponse {
    const message = createBaseReactivateUserSchemaResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseDeleteUserSchemaRequest(): DeleteUserSchemaRequest {
  return { id: "" };
}

export const DeleteUserSchemaRequest = {
  encode(message: DeleteUserSchemaRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteUserSchemaRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteUserSchemaRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteUserSchemaRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: DeleteUserSchemaRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<DeleteUserSchemaRequest>): DeleteUserSchemaRequest {
    return DeleteUserSchemaRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteUserSchemaRequest>): DeleteUserSchemaRequest {
    const message = createBaseDeleteUserSchemaRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseDeleteUserSchemaResponse(): DeleteUserSchemaResponse {
  return { details: undefined };
}

export const DeleteUserSchemaResponse = {
  encode(message: DeleteUserSchemaResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteUserSchemaResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteUserSchemaResponse();
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

  fromJSON(object: any): DeleteUserSchemaResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeleteUserSchemaResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeleteUserSchemaResponse>): DeleteUserSchemaResponse {
    return DeleteUserSchemaResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteUserSchemaResponse>): DeleteUserSchemaResponse {
    const message = createBaseDeleteUserSchemaResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

export type UserSchemaServiceDefinition = typeof UserSchemaServiceDefinition;
export const UserSchemaServiceDefinition = {
  name: "UserSchemaService",
  fullName: "zitadel.user.schema.v3alpha.UserSchemaService",
  methods: {
    /**
     * List user schemas
     *
     * List all matching user schemas. By default, we will return all user schema of your instance. Make sure to include a limit and sorting for pagination.
     */
    listUserSchemas: {
      name: "ListUserSchemas",
      requestType: ListUserSchemasRequest,
      requestStream: false,
      responseType: ListUserSchemasResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              113,
              74,
              53,
              10,
              3,
              50,
              48,
              48,
              18,
              46,
              10,
              44,
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
              117,
              115,
              101,
              114,
              32,
              115,
              99,
              104,
              101,
              109,
              97,
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
            Buffer.from([19, 10, 17, 10, 15, 117, 115, 101, 114, 115, 99, 104, 101, 109, 97, 46, 114, 101, 97, 100]),
          ],
          578365826: [
            Buffer.from([
              33,
              58,
              1,
              42,
              34,
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
              117,
              115,
              101,
              114,
              95,
              115,
              99,
              104,
              101,
              109,
              97,
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
     * User schema by ID
     *
     * Returns the user schema identified by the requested ID.
     */
    getUserSchemaByID: {
      name: "GetUserSchemaByID",
      requestType: GetUserSchemaByIDRequest,
      requestStream: false,
      responseType: GetUserSchemaByIDResponse,
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
              83,
              99,
              104,
              101,
              109,
              97,
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
            Buffer.from([19, 10, 17, 10, 15, 117, 115, 101, 114, 115, 99, 104, 101, 109, 97, 46, 114, 101, 97, 100]),
          ],
          578365826: [
            Buffer.from([
              28,
              18,
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
              117,
              115,
              101,
              114,
              95,
              115,
              99,
              104,
              101,
              109,
              97,
              115,
              47,
              123,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Create a user schema
     *
     * Create the first revision of a new user schema. The schema can then be used on users to store and validate their data.
     */
    createUserSchema: {
      name: "CreateUserSchema",
      requestType: CreateUserSchemaRequest,
      requestStream: false,
      responseType: CreateUserSchemaResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              89,
              74,
              87,
              10,
              3,
              50,
              48,
              49,
              18,
              80,
              10,
              27,
              83,
              99,
              104,
              101,
              109,
              97,
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
              49,
              10,
              47,
              26,
              45,
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
              85,
              115,
              101,
              114,
              83,
              99,
              104,
              101,
              109,
              97,
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
              25,
              10,
              18,
              10,
              16,
              117,
              115,
              101,
              114,
              115,
              99,
              104,
              101,
              109,
              97,
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
              26,
              58,
              1,
              42,
              34,
              21,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              95,
              115,
              99,
              104,
              101,
              109,
              97,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * Update a user schema
     *
     * Update an existing user schema to a new revision. Users based on the current revision will not be affected until they are updated.
     */
    updateUserSchema: {
      name: "UpdateUserSchema",
      requestType: UpdateUserSchemaRequest,
      requestStream: false,
      responseType: UpdateUserSchemaResponse,
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
              83,
              99,
              104,
              101,
              109,
              97,
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
              20,
              10,
              18,
              10,
              16,
              117,
              115,
              101,
              114,
              115,
              99,
              104,
              101,
              109,
              97,
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
              31,
              58,
              1,
              42,
              26,
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
              117,
              115,
              101,
              114,
              95,
              115,
              99,
              104,
              101,
              109,
              97,
              115,
              47,
              123,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Deactivate a user schema
     *
     * Deactivate an existing user schema and change it into a read-only state. Users based on this schema cannot be updated anymore, but are still able to authenticate.
     */
    deactivateUserSchema: {
      name: "DeactivateUserSchema",
      requestType: DeactivateUserSchemaRequest,
      requestStream: false,
      responseType: DeactivateUserSchemaResponse,
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
              83,
              99,
              104,
              101,
              109,
              97,
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
              97,
              99,
              116,
              105,
              118,
              97,
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
              117,
              115,
              101,
              114,
              115,
              99,
              104,
              101,
              109,
              97,
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
              39,
              34,
              37,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              95,
              115,
              99,
              104,
              101,
              109,
              97,
              115,
              47,
              123,
              105,
              100,
              125,
              47,
              100,
              101,
              97,
              99,
              116,
              105,
              118,
              97,
              116,
              101,
            ]),
          ],
        },
      },
    },
    /**
     * Reactivate a user schema
     *
     * Reactivate an previously deactivated user schema and change it into an active state again.
     */
    reactivateUserSchema: {
      name: "ReactivateUserSchema",
      requestType: ReactivateUserSchemaRequest,
      requestStream: false,
      responseType: ReactivateUserSchemaResponse,
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
              83,
              99,
              104,
              101,
              109,
              97,
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
              97,
              99,
              116,
              105,
              118,
              97,
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
              117,
              115,
              101,
              114,
              115,
              99,
              104,
              101,
              109,
              97,
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
              39,
              34,
              37,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              95,
              115,
              99,
              104,
              101,
              109,
              97,
              115,
              47,
              123,
              105,
              100,
              125,
              47,
              114,
              101,
              97,
              99,
              116,
              105,
              118,
              97,
              116,
              101,
            ]),
          ],
        },
      },
    },
    /**
     * Delete a user schema
     *
     * Delete an existing user schema. This operation is only allowed if there are no associated users to it.
     */
    deleteUserSchema: {
      name: "DeleteUserSchema",
      requestType: DeleteUserSchemaRequest,
      requestStream: false,
      responseType: DeleteUserSchemaResponse,
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
              83,
              99,
              104,
              101,
              109,
              97,
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
              21,
              10,
              19,
              10,
              17,
              117,
              115,
              101,
              114,
              115,
              99,
              104,
              101,
              109,
              97,
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
              28,
              42,
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
              117,
              115,
              101,
              114,
              95,
              115,
              99,
              104,
              101,
              109,
              97,
              115,
              47,
              123,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
  },
} as const;

export interface UserSchemaServiceImplementation<CallContextExt = {}> {
  /**
   * List user schemas
   *
   * List all matching user schemas. By default, we will return all user schema of your instance. Make sure to include a limit and sorting for pagination.
   */
  listUserSchemas(
    request: ListUserSchemasRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListUserSchemasResponse>>;
  /**
   * User schema by ID
   *
   * Returns the user schema identified by the requested ID.
   */
  getUserSchemaByID(
    request: GetUserSchemaByIDRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetUserSchemaByIDResponse>>;
  /**
   * Create a user schema
   *
   * Create the first revision of a new user schema. The schema can then be used on users to store and validate their data.
   */
  createUserSchema(
    request: CreateUserSchemaRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<CreateUserSchemaResponse>>;
  /**
   * Update a user schema
   *
   * Update an existing user schema to a new revision. Users based on the current revision will not be affected until they are updated.
   */
  updateUserSchema(
    request: UpdateUserSchemaRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateUserSchemaResponse>>;
  /**
   * Deactivate a user schema
   *
   * Deactivate an existing user schema and change it into a read-only state. Users based on this schema cannot be updated anymore, but are still able to authenticate.
   */
  deactivateUserSchema(
    request: DeactivateUserSchemaRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeactivateUserSchemaResponse>>;
  /**
   * Reactivate a user schema
   *
   * Reactivate an previously deactivated user schema and change it into an active state again.
   */
  reactivateUserSchema(
    request: ReactivateUserSchemaRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ReactivateUserSchemaResponse>>;
  /**
   * Delete a user schema
   *
   * Delete an existing user schema. This operation is only allowed if there are no associated users to it.
   */
  deleteUserSchema(
    request: DeleteUserSchemaRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeleteUserSchemaResponse>>;
}

export interface UserSchemaServiceClient<CallOptionsExt = {}> {
  /**
   * List user schemas
   *
   * List all matching user schemas. By default, we will return all user schema of your instance. Make sure to include a limit and sorting for pagination.
   */
  listUserSchemas(
    request: DeepPartial<ListUserSchemasRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListUserSchemasResponse>;
  /**
   * User schema by ID
   *
   * Returns the user schema identified by the requested ID.
   */
  getUserSchemaByID(
    request: DeepPartial<GetUserSchemaByIDRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetUserSchemaByIDResponse>;
  /**
   * Create a user schema
   *
   * Create the first revision of a new user schema. The schema can then be used on users to store and validate their data.
   */
  createUserSchema(
    request: DeepPartial<CreateUserSchemaRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<CreateUserSchemaResponse>;
  /**
   * Update a user schema
   *
   * Update an existing user schema to a new revision. Users based on the current revision will not be affected until they are updated.
   */
  updateUserSchema(
    request: DeepPartial<UpdateUserSchemaRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateUserSchemaResponse>;
  /**
   * Deactivate a user schema
   *
   * Deactivate an existing user schema and change it into a read-only state. Users based on this schema cannot be updated anymore, but are still able to authenticate.
   */
  deactivateUserSchema(
    request: DeepPartial<DeactivateUserSchemaRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeactivateUserSchemaResponse>;
  /**
   * Reactivate a user schema
   *
   * Reactivate an previously deactivated user schema and change it into an active state again.
   */
  reactivateUserSchema(
    request: DeepPartial<ReactivateUserSchemaRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ReactivateUserSchemaResponse>;
  /**
   * Delete a user schema
   *
   * Delete an existing user schema. This operation is only allowed if there are no associated users to it.
   */
  deleteUserSchema(
    request: DeepPartial<DeleteUserSchemaRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeleteUserSchemaResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

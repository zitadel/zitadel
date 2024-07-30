/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Struct } from "../../../../google/protobuf/struct";
import {
  Details,
  TextQueryMethod,
  textQueryMethodFromJSON,
  textQueryMethodToJSON,
} from "../../../object/v2beta/object";

export const protobufPackage = "zitadel.user.schema.v3alpha";

export enum FieldName {
  FIELD_NAME_UNSPECIFIED = 0,
  FIELD_NAME_TYPE = 1,
  FIELD_NAME_STATE = 2,
  FIELD_NAME_REVISION = 3,
  FIELD_NAME_CHANGE_DATE = 4,
  UNRECOGNIZED = -1,
}

export function fieldNameFromJSON(object: any): FieldName {
  switch (object) {
    case 0:
    case "FIELD_NAME_UNSPECIFIED":
      return FieldName.FIELD_NAME_UNSPECIFIED;
    case 1:
    case "FIELD_NAME_TYPE":
      return FieldName.FIELD_NAME_TYPE;
    case 2:
    case "FIELD_NAME_STATE":
      return FieldName.FIELD_NAME_STATE;
    case 3:
    case "FIELD_NAME_REVISION":
      return FieldName.FIELD_NAME_REVISION;
    case 4:
    case "FIELD_NAME_CHANGE_DATE":
      return FieldName.FIELD_NAME_CHANGE_DATE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FieldName.UNRECOGNIZED;
  }
}

export function fieldNameToJSON(object: FieldName): string {
  switch (object) {
    case FieldName.FIELD_NAME_UNSPECIFIED:
      return "FIELD_NAME_UNSPECIFIED";
    case FieldName.FIELD_NAME_TYPE:
      return "FIELD_NAME_TYPE";
    case FieldName.FIELD_NAME_STATE:
      return "FIELD_NAME_STATE";
    case FieldName.FIELD_NAME_REVISION:
      return "FIELD_NAME_REVISION";
    case FieldName.FIELD_NAME_CHANGE_DATE:
      return "FIELD_NAME_CHANGE_DATE";
    case FieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum State {
  STATE_UNSPECIFIED = 0,
  STATE_ACTIVE = 1,
  STATE_INACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function stateFromJSON(object: any): State {
  switch (object) {
    case 0:
    case "STATE_UNSPECIFIED":
      return State.STATE_UNSPECIFIED;
    case 1:
    case "STATE_ACTIVE":
      return State.STATE_ACTIVE;
    case 2:
    case "STATE_INACTIVE":
      return State.STATE_INACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return State.UNRECOGNIZED;
  }
}

export function stateToJSON(object: State): string {
  switch (object) {
    case State.STATE_UNSPECIFIED:
      return "STATE_UNSPECIFIED";
    case State.STATE_ACTIVE:
      return "STATE_ACTIVE";
    case State.STATE_INACTIVE:
      return "STATE_INACTIVE";
    case State.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum AuthenticatorType {
  AUTHENTICATOR_TYPE_UNSPECIFIED = 0,
  AUTHENTICATOR_TYPE_USERNAME = 1,
  AUTHENTICATOR_TYPE_PASSWORD = 2,
  AUTHENTICATOR_TYPE_WEBAUTHN = 3,
  AUTHENTICATOR_TYPE_TOTP = 4,
  AUTHENTICATOR_TYPE_OTP_EMAIL = 5,
  AUTHENTICATOR_TYPE_OTP_SMS = 6,
  AUTHENTICATOR_TYPE_AUTHENTICATION_KEY = 7,
  AUTHENTICATOR_TYPE_IDENTITY_PROVIDER = 8,
  UNRECOGNIZED = -1,
}

export function authenticatorTypeFromJSON(object: any): AuthenticatorType {
  switch (object) {
    case 0:
    case "AUTHENTICATOR_TYPE_UNSPECIFIED":
      return AuthenticatorType.AUTHENTICATOR_TYPE_UNSPECIFIED;
    case 1:
    case "AUTHENTICATOR_TYPE_USERNAME":
      return AuthenticatorType.AUTHENTICATOR_TYPE_USERNAME;
    case 2:
    case "AUTHENTICATOR_TYPE_PASSWORD":
      return AuthenticatorType.AUTHENTICATOR_TYPE_PASSWORD;
    case 3:
    case "AUTHENTICATOR_TYPE_WEBAUTHN":
      return AuthenticatorType.AUTHENTICATOR_TYPE_WEBAUTHN;
    case 4:
    case "AUTHENTICATOR_TYPE_TOTP":
      return AuthenticatorType.AUTHENTICATOR_TYPE_TOTP;
    case 5:
    case "AUTHENTICATOR_TYPE_OTP_EMAIL":
      return AuthenticatorType.AUTHENTICATOR_TYPE_OTP_EMAIL;
    case 6:
    case "AUTHENTICATOR_TYPE_OTP_SMS":
      return AuthenticatorType.AUTHENTICATOR_TYPE_OTP_SMS;
    case 7:
    case "AUTHENTICATOR_TYPE_AUTHENTICATION_KEY":
      return AuthenticatorType.AUTHENTICATOR_TYPE_AUTHENTICATION_KEY;
    case 8:
    case "AUTHENTICATOR_TYPE_IDENTITY_PROVIDER":
      return AuthenticatorType.AUTHENTICATOR_TYPE_IDENTITY_PROVIDER;
    case -1:
    case "UNRECOGNIZED":
    default:
      return AuthenticatorType.UNRECOGNIZED;
  }
}

export function authenticatorTypeToJSON(object: AuthenticatorType): string {
  switch (object) {
    case AuthenticatorType.AUTHENTICATOR_TYPE_UNSPECIFIED:
      return "AUTHENTICATOR_TYPE_UNSPECIFIED";
    case AuthenticatorType.AUTHENTICATOR_TYPE_USERNAME:
      return "AUTHENTICATOR_TYPE_USERNAME";
    case AuthenticatorType.AUTHENTICATOR_TYPE_PASSWORD:
      return "AUTHENTICATOR_TYPE_PASSWORD";
    case AuthenticatorType.AUTHENTICATOR_TYPE_WEBAUTHN:
      return "AUTHENTICATOR_TYPE_WEBAUTHN";
    case AuthenticatorType.AUTHENTICATOR_TYPE_TOTP:
      return "AUTHENTICATOR_TYPE_TOTP";
    case AuthenticatorType.AUTHENTICATOR_TYPE_OTP_EMAIL:
      return "AUTHENTICATOR_TYPE_OTP_EMAIL";
    case AuthenticatorType.AUTHENTICATOR_TYPE_OTP_SMS:
      return "AUTHENTICATOR_TYPE_OTP_SMS";
    case AuthenticatorType.AUTHENTICATOR_TYPE_AUTHENTICATION_KEY:
      return "AUTHENTICATOR_TYPE_AUTHENTICATION_KEY";
    case AuthenticatorType.AUTHENTICATOR_TYPE_IDENTITY_PROVIDER:
      return "AUTHENTICATOR_TYPE_IDENTITY_PROVIDER";
    case AuthenticatorType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface UserSchema {
  /** ID is the read-only unique identifier of the schema. */
  id: string;
  /** Details provide some base information (such as the last change date) of the schema. */
  details:
    | Details
    | undefined;
  /** Type is a human readable text describing the schema. */
  type: string;
  /** Current state of the schema. */
  state: State;
  /** Revision is a read only version of the schema, each update of the `schema`-field increases the revision. */
  revision: number;
  /** JSON schema representation defining the user. */
  schema:
    | { [key: string]: any }
    | undefined;
  /**
   * Defines the possible types of authenticators.
   * This allows creating different user types like human/machine without usage of actions to validate possible authenticators.
   * Removal of an authenticator does not remove the authenticator on a user.
   */
  possibleAuthenticators: AuthenticatorType[];
}

export interface SearchQuery {
  /** Union the results of each sub query ('OR'). */
  orQuery?:
    | OrQuery
    | undefined;
  /**
   * Limit the result to match all sub queries ('AND').
   * Note that if you specify multiple queries, they will be implicitly used as andQueries.
   * Use the andQuery in combination with orQuery and notQuery.
   */
  andQuery?:
    | AndQuery
    | undefined;
  /** Exclude / Negate the result of the sub query ('NOT'). */
  notQuery?:
    | NotQuery
    | undefined;
  /** Limit the result to a specific schema type. */
  typeQuery?:
    | TypeQuery
    | undefined;
  /** Limit the result to a specific state of the schema. */
  stateQuery?:
    | StateQuery
    | undefined;
  /** Limit the result to a specific schema ID. */
  idQuery?: IDQuery | undefined;
}

export interface OrQuery {
  queries: SearchQuery[];
}

export interface AndQuery {
  queries: SearchQuery[];
}

export interface NotQuery {
  query: SearchQuery | undefined;
}

export interface IDQuery {
  /** Defines the ID of the user schema to query for. */
  id: string;
  /** Defines which text comparison method used for the id query. */
  method: TextQueryMethod;
}

export interface TypeQuery {
  /** Defines which type to query for. */
  type: string;
  /** Defines which text comparison method used for the type query. */
  method: TextQueryMethod;
}

export interface StateQuery {
  /** Defines the state to query for. */
  state: State;
}

function createBaseUserSchema(): UserSchema {
  return { id: "", details: undefined, type: "", state: 0, revision: 0, schema: undefined, possibleAuthenticators: [] };
}

export const UserSchema = {
  encode(message: UserSchema, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.type !== "") {
      writer.uint32(26).string(message.type);
    }
    if (message.state !== 0) {
      writer.uint32(32).int32(message.state);
    }
    if (message.revision !== 0) {
      writer.uint32(40).uint32(message.revision);
    }
    if (message.schema !== undefined) {
      Struct.encode(Struct.wrap(message.schema), writer.uint32(50).fork()).ldelim();
    }
    writer.uint32(58).fork();
    for (const v of message.possibleAuthenticators) {
      writer.int32(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserSchema {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserSchema();
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
        case 3:
          if (tag != 26) {
            break;
          }

          message.type = reader.string();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.revision = reader.uint32();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.schema = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag == 56) {
            message.possibleAuthenticators.push(reader.int32() as any);
            continue;
          }

          if (tag == 58) {
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

  fromJSON(object: any): UserSchema {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      type: isSet(object.type) ? String(object.type) : "",
      state: isSet(object.state) ? stateFromJSON(object.state) : 0,
      revision: isSet(object.revision) ? Number(object.revision) : 0,
      schema: isObject(object.schema) ? object.schema : undefined,
      possibleAuthenticators: Array.isArray(object?.possibleAuthenticators)
        ? object.possibleAuthenticators.map((e: any) => authenticatorTypeFromJSON(e))
        : [],
    };
  },

  toJSON(message: UserSchema): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.type !== undefined && (obj.type = message.type);
    message.state !== undefined && (obj.state = stateToJSON(message.state));
    message.revision !== undefined && (obj.revision = Math.round(message.revision));
    message.schema !== undefined && (obj.schema = message.schema);
    if (message.possibleAuthenticators) {
      obj.possibleAuthenticators = message.possibleAuthenticators.map((e) => authenticatorTypeToJSON(e));
    } else {
      obj.possibleAuthenticators = [];
    }
    return obj;
  },

  create(base?: DeepPartial<UserSchema>): UserSchema {
    return UserSchema.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserSchema>): UserSchema {
    const message = createBaseUserSchema();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.type = object.type ?? "";
    message.state = object.state ?? 0;
    message.revision = object.revision ?? 0;
    message.schema = object.schema ?? undefined;
    message.possibleAuthenticators = object.possibleAuthenticators?.map((e) => e) || [];
    return message;
  },
};

function createBaseSearchQuery(): SearchQuery {
  return {
    orQuery: undefined,
    andQuery: undefined,
    notQuery: undefined,
    typeQuery: undefined,
    stateQuery: undefined,
    idQuery: undefined,
  };
}

export const SearchQuery = {
  encode(message: SearchQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orQuery !== undefined) {
      OrQuery.encode(message.orQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.andQuery !== undefined) {
      AndQuery.encode(message.andQuery, writer.uint32(18).fork()).ldelim();
    }
    if (message.notQuery !== undefined) {
      NotQuery.encode(message.notQuery, writer.uint32(26).fork()).ldelim();
    }
    if (message.typeQuery !== undefined) {
      TypeQuery.encode(message.typeQuery, writer.uint32(42).fork()).ldelim();
    }
    if (message.stateQuery !== undefined) {
      StateQuery.encode(message.stateQuery, writer.uint32(50).fork()).ldelim();
    }
    if (message.idQuery !== undefined) {
      IDQuery.encode(message.idQuery, writer.uint32(58).fork()).ldelim();
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

          message.orQuery = OrQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.andQuery = AndQuery.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.notQuery = NotQuery.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.typeQuery = TypeQuery.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.stateQuery = StateQuery.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.idQuery = IDQuery.decode(reader, reader.uint32());
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
      orQuery: isSet(object.orQuery) ? OrQuery.fromJSON(object.orQuery) : undefined,
      andQuery: isSet(object.andQuery) ? AndQuery.fromJSON(object.andQuery) : undefined,
      notQuery: isSet(object.notQuery) ? NotQuery.fromJSON(object.notQuery) : undefined,
      typeQuery: isSet(object.typeQuery) ? TypeQuery.fromJSON(object.typeQuery) : undefined,
      stateQuery: isSet(object.stateQuery) ? StateQuery.fromJSON(object.stateQuery) : undefined,
      idQuery: isSet(object.idQuery) ? IDQuery.fromJSON(object.idQuery) : undefined,
    };
  },

  toJSON(message: SearchQuery): unknown {
    const obj: any = {};
    message.orQuery !== undefined && (obj.orQuery = message.orQuery ? OrQuery.toJSON(message.orQuery) : undefined);
    message.andQuery !== undefined && (obj.andQuery = message.andQuery ? AndQuery.toJSON(message.andQuery) : undefined);
    message.notQuery !== undefined && (obj.notQuery = message.notQuery ? NotQuery.toJSON(message.notQuery) : undefined);
    message.typeQuery !== undefined &&
      (obj.typeQuery = message.typeQuery ? TypeQuery.toJSON(message.typeQuery) : undefined);
    message.stateQuery !== undefined &&
      (obj.stateQuery = message.stateQuery ? StateQuery.toJSON(message.stateQuery) : undefined);
    message.idQuery !== undefined && (obj.idQuery = message.idQuery ? IDQuery.toJSON(message.idQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SearchQuery>): SearchQuery {
    return SearchQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SearchQuery>): SearchQuery {
    const message = createBaseSearchQuery();
    message.orQuery = (object.orQuery !== undefined && object.orQuery !== null)
      ? OrQuery.fromPartial(object.orQuery)
      : undefined;
    message.andQuery = (object.andQuery !== undefined && object.andQuery !== null)
      ? AndQuery.fromPartial(object.andQuery)
      : undefined;
    message.notQuery = (object.notQuery !== undefined && object.notQuery !== null)
      ? NotQuery.fromPartial(object.notQuery)
      : undefined;
    message.typeQuery = (object.typeQuery !== undefined && object.typeQuery !== null)
      ? TypeQuery.fromPartial(object.typeQuery)
      : undefined;
    message.stateQuery = (object.stateQuery !== undefined && object.stateQuery !== null)
      ? StateQuery.fromPartial(object.stateQuery)
      : undefined;
    message.idQuery = (object.idQuery !== undefined && object.idQuery !== null)
      ? IDQuery.fromPartial(object.idQuery)
      : undefined;
    return message;
  },
};

function createBaseOrQuery(): OrQuery {
  return { queries: [] };
}

export const OrQuery = {
  encode(message: OrQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.queries) {
      SearchQuery.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): OrQuery {
    return { queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => SearchQuery.fromJSON(e)) : [] };
  },

  toJSON(message: OrQuery): unknown {
    const obj: any = {};
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? SearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<OrQuery>): OrQuery {
    return OrQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OrQuery>): OrQuery {
    const message = createBaseOrQuery();
    message.queries = object.queries?.map((e) => SearchQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAndQuery(): AndQuery {
  return { queries: [] };
}

export const AndQuery = {
  encode(message: AndQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.queries) {
      SearchQuery.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AndQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAndQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): AndQuery {
    return { queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => SearchQuery.fromJSON(e)) : [] };
  },

  toJSON(message: AndQuery): unknown {
    const obj: any = {};
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? SearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<AndQuery>): AndQuery {
    return AndQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AndQuery>): AndQuery {
    const message = createBaseAndQuery();
    message.queries = object.queries?.map((e) => SearchQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseNotQuery(): NotQuery {
  return { query: undefined };
}

export const NotQuery = {
  encode(message: NotQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      SearchQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NotQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNotQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.query = SearchQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): NotQuery {
    return { query: isSet(object.query) ? SearchQuery.fromJSON(object.query) : undefined };
  },

  toJSON(message: NotQuery): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? SearchQuery.toJSON(message.query) : undefined);
    return obj;
  },

  create(base?: DeepPartial<NotQuery>): NotQuery {
    return NotQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<NotQuery>): NotQuery {
    const message = createBaseNotQuery();
    message.query = (object.query !== undefined && object.query !== null)
      ? SearchQuery.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseIDQuery(): IDQuery {
  return { id: "", method: 0 };
}

export const IDQuery = {
  encode(message: IDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDQuery();
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

  fromJSON(object: any): IDQuery {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: IDQuery): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<IDQuery>): IDQuery {
    return IDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDQuery>): IDQuery {
    const message = createBaseIDQuery();
    message.id = object.id ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseTypeQuery(): TypeQuery {
  return { type: "", method: 0 };
}

export const TypeQuery = {
  encode(message: TypeQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TypeQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTypeQuery();
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

  fromJSON(object: any): TypeQuery {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: TypeQuery): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<TypeQuery>): TypeQuery {
    return TypeQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TypeQuery>): TypeQuery {
    const message = createBaseTypeQuery();
    message.type = object.type ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseStateQuery(): StateQuery {
  return { state: 0 };
}

export const StateQuery = {
  encode(message: StateQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.state !== 0) {
      writer.uint32(8).int32(message.state);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StateQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStateQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StateQuery {
    return { state: isSet(object.state) ? stateFromJSON(object.state) : 0 };
  },

  toJSON(message: StateQuery): unknown {
    const obj: any = {};
    message.state !== undefined && (obj.state = stateToJSON(message.state));
    return obj;
  },

  create(base?: DeepPartial<StateQuery>): StateQuery {
    return StateQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StateQuery>): StateQuery {
    const message = createBaseStateQuery();
    message.state = object.state ?? 0;
    return message;
  },
};

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

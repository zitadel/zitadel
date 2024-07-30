/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { TextQueryMethod, textQueryMethodFromJSON, textQueryMethodToJSON } from "../../object/v2beta/object";
import { State, stateFromJSON, stateToJSON } from "./user";

export const protobufPackage = "zitadel.user.v3alpha";

export enum FieldName {
  FIELD_NAME_UNSPECIFIED = 0,
  FIELD_NAME_ID = 1,
  FIELD_NAME_CREATION_DATE = 2,
  FIELD_NAME_CHANGE_DATE = 3,
  FIELD_NAME_EMAIL = 4,
  FIELD_NAME_PHONE = 5,
  FIELD_NAME_STATE = 6,
  FIELD_NAME_SCHEMA_ID = 7,
  FIELD_NAME_SCHEMA_TYPE = 8,
  UNRECOGNIZED = -1,
}

export function fieldNameFromJSON(object: any): FieldName {
  switch (object) {
    case 0:
    case "FIELD_NAME_UNSPECIFIED":
      return FieldName.FIELD_NAME_UNSPECIFIED;
    case 1:
    case "FIELD_NAME_ID":
      return FieldName.FIELD_NAME_ID;
    case 2:
    case "FIELD_NAME_CREATION_DATE":
      return FieldName.FIELD_NAME_CREATION_DATE;
    case 3:
    case "FIELD_NAME_CHANGE_DATE":
      return FieldName.FIELD_NAME_CHANGE_DATE;
    case 4:
    case "FIELD_NAME_EMAIL":
      return FieldName.FIELD_NAME_EMAIL;
    case 5:
    case "FIELD_NAME_PHONE":
      return FieldName.FIELD_NAME_PHONE;
    case 6:
    case "FIELD_NAME_STATE":
      return FieldName.FIELD_NAME_STATE;
    case 7:
    case "FIELD_NAME_SCHEMA_ID":
      return FieldName.FIELD_NAME_SCHEMA_ID;
    case 8:
    case "FIELD_NAME_SCHEMA_TYPE":
      return FieldName.FIELD_NAME_SCHEMA_TYPE;
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
    case FieldName.FIELD_NAME_ID:
      return "FIELD_NAME_ID";
    case FieldName.FIELD_NAME_CREATION_DATE:
      return "FIELD_NAME_CREATION_DATE";
    case FieldName.FIELD_NAME_CHANGE_DATE:
      return "FIELD_NAME_CHANGE_DATE";
    case FieldName.FIELD_NAME_EMAIL:
      return "FIELD_NAME_EMAIL";
    case FieldName.FIELD_NAME_PHONE:
      return "FIELD_NAME_PHONE";
    case FieldName.FIELD_NAME_STATE:
      return "FIELD_NAME_STATE";
    case FieldName.FIELD_NAME_SCHEMA_ID:
      return "FIELD_NAME_SCHEMA_ID";
    case FieldName.FIELD_NAME_SCHEMA_TYPE:
      return "FIELD_NAME_SCHEMA_TYPE";
    case FieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
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
  /** Limit the result to a specific user ID. */
  userIdQuery?:
    | UserIDQuery
    | undefined;
  /** Limit the result to a specific organization. */
  organizationIdQuery?:
    | OrganizationIDQuery
    | undefined;
  /** Limit the result to a specific username. */
  usernameQuery?:
    | UsernameQuery
    | undefined;
  /** Limit the result to a specific contact email. */
  emailQuery?:
    | EmailQuery
    | undefined;
  /** Limit the result to a specific contact phone. */
  phoneQuery?:
    | PhoneQuery
    | undefined;
  /** Limit the result to a specific state of the user. */
  stateQuery?:
    | StateQuery
    | undefined;
  /** Limit the result to a specific schema ID. */
  schemaIDQuery?:
    | SchemaIDQuery
    | undefined;
  /** Limit the result to a specific schema type. */
  schemaTypeQuery?: SchemaTypeQuery | undefined;
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

export interface UserIDQuery {
  /** Defines the ID of the user to query for. */
  id: string;
  /** Defines which text comparison method used for the id query. */
  method: TextQueryMethod;
}

export interface OrganizationIDQuery {
  /** Defines the ID of the organization to query for. */
  id: string;
  /** Defines which text comparison method used for the id query. */
  method: TextQueryMethod;
}

export interface UsernameQuery {
  /** Defines the username to query for. */
  username: string;
  /** Defines which text comparison method used for the username query. */
  method: TextQueryMethod;
  /** Defines that the username must only be unique in the organisation. */
  isOrganizationSpecific: boolean;
}

export interface EmailQuery {
  /** Defines the email of the user to query for. */
  address: string;
  /** Defines which text comparison method used for the email query. */
  method: TextQueryMethod;
}

export interface PhoneQuery {
  /** Defines the phone of the user to query for. */
  number: string;
  /** Defines which text comparison method used for the phone query. */
  method: TextQueryMethod;
}

export interface StateQuery {
  /** Defines the state to query for. */
  state: State;
}

export interface SchemaIDQuery {
  /** Defines the ID of the schema to query for. */
  id: string;
}

export interface SchemaTypeQuery {
  /** Defines which type to query for. */
  type: string;
  /** Defines which text comparison method used for the type query. */
  method: TextQueryMethod;
}

function createBaseSearchQuery(): SearchQuery {
  return {
    orQuery: undefined,
    andQuery: undefined,
    notQuery: undefined,
    userIdQuery: undefined,
    organizationIdQuery: undefined,
    usernameQuery: undefined,
    emailQuery: undefined,
    phoneQuery: undefined,
    stateQuery: undefined,
    schemaIDQuery: undefined,
    schemaTypeQuery: undefined,
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
    if (message.userIdQuery !== undefined) {
      UserIDQuery.encode(message.userIdQuery, writer.uint32(34).fork()).ldelim();
    }
    if (message.organizationIdQuery !== undefined) {
      OrganizationIDQuery.encode(message.organizationIdQuery, writer.uint32(42).fork()).ldelim();
    }
    if (message.usernameQuery !== undefined) {
      UsernameQuery.encode(message.usernameQuery, writer.uint32(50).fork()).ldelim();
    }
    if (message.emailQuery !== undefined) {
      EmailQuery.encode(message.emailQuery, writer.uint32(58).fork()).ldelim();
    }
    if (message.phoneQuery !== undefined) {
      PhoneQuery.encode(message.phoneQuery, writer.uint32(66).fork()).ldelim();
    }
    if (message.stateQuery !== undefined) {
      StateQuery.encode(message.stateQuery, writer.uint32(74).fork()).ldelim();
    }
    if (message.schemaIDQuery !== undefined) {
      SchemaIDQuery.encode(message.schemaIDQuery, writer.uint32(82).fork()).ldelim();
    }
    if (message.schemaTypeQuery !== undefined) {
      SchemaTypeQuery.encode(message.schemaTypeQuery, writer.uint32(90).fork()).ldelim();
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
        case 4:
          if (tag != 34) {
            break;
          }

          message.userIdQuery = UserIDQuery.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.organizationIdQuery = OrganizationIDQuery.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.usernameQuery = UsernameQuery.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.emailQuery = EmailQuery.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.phoneQuery = PhoneQuery.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.stateQuery = StateQuery.decode(reader, reader.uint32());
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          message.schemaIDQuery = SchemaIDQuery.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag != 90) {
            break;
          }

          message.schemaTypeQuery = SchemaTypeQuery.decode(reader, reader.uint32());
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
      userIdQuery: isSet(object.userIdQuery) ? UserIDQuery.fromJSON(object.userIdQuery) : undefined,
      organizationIdQuery: isSet(object.organizationIdQuery)
        ? OrganizationIDQuery.fromJSON(object.organizationIdQuery)
        : undefined,
      usernameQuery: isSet(object.usernameQuery) ? UsernameQuery.fromJSON(object.usernameQuery) : undefined,
      emailQuery: isSet(object.emailQuery) ? EmailQuery.fromJSON(object.emailQuery) : undefined,
      phoneQuery: isSet(object.phoneQuery) ? PhoneQuery.fromJSON(object.phoneQuery) : undefined,
      stateQuery: isSet(object.stateQuery) ? StateQuery.fromJSON(object.stateQuery) : undefined,
      schemaIDQuery: isSet(object.schemaIDQuery) ? SchemaIDQuery.fromJSON(object.schemaIDQuery) : undefined,
      schemaTypeQuery: isSet(object.schemaTypeQuery) ? SchemaTypeQuery.fromJSON(object.schemaTypeQuery) : undefined,
    };
  },

  toJSON(message: SearchQuery): unknown {
    const obj: any = {};
    message.orQuery !== undefined && (obj.orQuery = message.orQuery ? OrQuery.toJSON(message.orQuery) : undefined);
    message.andQuery !== undefined && (obj.andQuery = message.andQuery ? AndQuery.toJSON(message.andQuery) : undefined);
    message.notQuery !== undefined && (obj.notQuery = message.notQuery ? NotQuery.toJSON(message.notQuery) : undefined);
    message.userIdQuery !== undefined &&
      (obj.userIdQuery = message.userIdQuery ? UserIDQuery.toJSON(message.userIdQuery) : undefined);
    message.organizationIdQuery !== undefined && (obj.organizationIdQuery = message.organizationIdQuery
      ? OrganizationIDQuery.toJSON(message.organizationIdQuery)
      : undefined);
    message.usernameQuery !== undefined &&
      (obj.usernameQuery = message.usernameQuery ? UsernameQuery.toJSON(message.usernameQuery) : undefined);
    message.emailQuery !== undefined &&
      (obj.emailQuery = message.emailQuery ? EmailQuery.toJSON(message.emailQuery) : undefined);
    message.phoneQuery !== undefined &&
      (obj.phoneQuery = message.phoneQuery ? PhoneQuery.toJSON(message.phoneQuery) : undefined);
    message.stateQuery !== undefined &&
      (obj.stateQuery = message.stateQuery ? StateQuery.toJSON(message.stateQuery) : undefined);
    message.schemaIDQuery !== undefined &&
      (obj.schemaIDQuery = message.schemaIDQuery ? SchemaIDQuery.toJSON(message.schemaIDQuery) : undefined);
    message.schemaTypeQuery !== undefined &&
      (obj.schemaTypeQuery = message.schemaTypeQuery ? SchemaTypeQuery.toJSON(message.schemaTypeQuery) : undefined);
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
    message.userIdQuery = (object.userIdQuery !== undefined && object.userIdQuery !== null)
      ? UserIDQuery.fromPartial(object.userIdQuery)
      : undefined;
    message.organizationIdQuery = (object.organizationIdQuery !== undefined && object.organizationIdQuery !== null)
      ? OrganizationIDQuery.fromPartial(object.organizationIdQuery)
      : undefined;
    message.usernameQuery = (object.usernameQuery !== undefined && object.usernameQuery !== null)
      ? UsernameQuery.fromPartial(object.usernameQuery)
      : undefined;
    message.emailQuery = (object.emailQuery !== undefined && object.emailQuery !== null)
      ? EmailQuery.fromPartial(object.emailQuery)
      : undefined;
    message.phoneQuery = (object.phoneQuery !== undefined && object.phoneQuery !== null)
      ? PhoneQuery.fromPartial(object.phoneQuery)
      : undefined;
    message.stateQuery = (object.stateQuery !== undefined && object.stateQuery !== null)
      ? StateQuery.fromPartial(object.stateQuery)
      : undefined;
    message.schemaIDQuery = (object.schemaIDQuery !== undefined && object.schemaIDQuery !== null)
      ? SchemaIDQuery.fromPartial(object.schemaIDQuery)
      : undefined;
    message.schemaTypeQuery = (object.schemaTypeQuery !== undefined && object.schemaTypeQuery !== null)
      ? SchemaTypeQuery.fromPartial(object.schemaTypeQuery)
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

function createBaseUserIDQuery(): UserIDQuery {
  return { id: "", method: 0 };
}

export const UserIDQuery = {
  encode(message: UserIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserIDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserIDQuery();
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

  fromJSON(object: any): UserIDQuery {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: UserIDQuery): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<UserIDQuery>): UserIDQuery {
    return UserIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserIDQuery>): UserIDQuery {
    const message = createBaseUserIDQuery();
    message.id = object.id ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseOrganizationIDQuery(): OrganizationIDQuery {
  return { id: "", method: 0 };
}

export const OrganizationIDQuery = {
  encode(message: OrganizationIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrganizationIDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrganizationIDQuery();
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

  fromJSON(object: any): OrganizationIDQuery {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: OrganizationIDQuery): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<OrganizationIDQuery>): OrganizationIDQuery {
    return OrganizationIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OrganizationIDQuery>): OrganizationIDQuery {
    const message = createBaseOrganizationIDQuery();
    message.id = object.id ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseUsernameQuery(): UsernameQuery {
  return { username: "", method: 0, isOrganizationSpecific: false };
}

export const UsernameQuery = {
  encode(message: UsernameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    if (message.isOrganizationSpecific === true) {
      writer.uint32(24).bool(message.isOrganizationSpecific);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UsernameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUsernameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.username = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.isOrganizationSpecific = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UsernameQuery {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
      isOrganizationSpecific: isSet(object.isOrganizationSpecific) ? Boolean(object.isOrganizationSpecific) : false,
    };
  },

  toJSON(message: UsernameQuery): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    message.isOrganizationSpecific !== undefined && (obj.isOrganizationSpecific = message.isOrganizationSpecific);
    return obj;
  },

  create(base?: DeepPartial<UsernameQuery>): UsernameQuery {
    return UsernameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UsernameQuery>): UsernameQuery {
    const message = createBaseUsernameQuery();
    message.username = object.username ?? "";
    message.method = object.method ?? 0;
    message.isOrganizationSpecific = object.isOrganizationSpecific ?? false;
    return message;
  },
};

function createBaseEmailQuery(): EmailQuery {
  return { address: "", method: 0 };
}

export const EmailQuery = {
  encode(message: EmailQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EmailQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEmailQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.address = reader.string();
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

  fromJSON(object: any): EmailQuery {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: EmailQuery): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<EmailQuery>): EmailQuery {
    return EmailQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EmailQuery>): EmailQuery {
    const message = createBaseEmailQuery();
    message.address = object.address ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBasePhoneQuery(): PhoneQuery {
  return { number: "", method: 0 };
}

export const PhoneQuery = {
  encode(message: PhoneQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.number !== "") {
      writer.uint32(10).string(message.number);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PhoneQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePhoneQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.number = reader.string();
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

  fromJSON(object: any): PhoneQuery {
    return {
      number: isSet(object.number) ? String(object.number) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: PhoneQuery): unknown {
    const obj: any = {};
    message.number !== undefined && (obj.number = message.number);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<PhoneQuery>): PhoneQuery {
    return PhoneQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PhoneQuery>): PhoneQuery {
    const message = createBasePhoneQuery();
    message.number = object.number ?? "";
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

function createBaseSchemaIDQuery(): SchemaIDQuery {
  return { id: "" };
}

export const SchemaIDQuery = {
  encode(message: SchemaIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SchemaIDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSchemaIDQuery();
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

  fromJSON(object: any): SchemaIDQuery {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: SchemaIDQuery): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<SchemaIDQuery>): SchemaIDQuery {
    return SchemaIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SchemaIDQuery>): SchemaIDQuery {
    const message = createBaseSchemaIDQuery();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseSchemaTypeQuery(): SchemaTypeQuery {
  return { type: "", method: 0 };
}

export const SchemaTypeQuery = {
  encode(message: SchemaTypeQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SchemaTypeQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSchemaTypeQuery();
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

  fromJSON(object: any): SchemaTypeQuery {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: SchemaTypeQuery): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<SchemaTypeQuery>): SchemaTypeQuery {
    return SchemaTypeQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SchemaTypeQuery>): SchemaTypeQuery {
    const message = createBaseSchemaTypeQuery();
    message.type = object.type ?? "";
    message.method = object.method ?? 0;
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

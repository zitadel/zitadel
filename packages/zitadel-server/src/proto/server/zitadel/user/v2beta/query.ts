/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { TextQueryMethod, textQueryMethodFromJSON, textQueryMethodToJSON } from "../../object/v2beta/object";
import { UserState, userStateFromJSON, userStateToJSON } from "./user";

export const protobufPackage = "zitadel.user.v2beta";

export enum Type {
  TYPE_UNSPECIFIED = 0,
  TYPE_HUMAN = 1,
  TYPE_MACHINE = 2,
  UNRECOGNIZED = -1,
}

export function typeFromJSON(object: any): Type {
  switch (object) {
    case 0:
    case "TYPE_UNSPECIFIED":
      return Type.TYPE_UNSPECIFIED;
    case 1:
    case "TYPE_HUMAN":
      return Type.TYPE_HUMAN;
    case 2:
    case "TYPE_MACHINE":
      return Type.TYPE_MACHINE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Type.UNRECOGNIZED;
  }
}

export function typeToJSON(object: Type): string {
  switch (object) {
    case Type.TYPE_UNSPECIFIED:
      return "TYPE_UNSPECIFIED";
    case Type.TYPE_HUMAN:
      return "TYPE_HUMAN";
    case Type.TYPE_MACHINE:
      return "TYPE_MACHINE";
    case Type.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum UserFieldName {
  USER_FIELD_NAME_UNSPECIFIED = 0,
  USER_FIELD_NAME_USER_NAME = 1,
  USER_FIELD_NAME_FIRST_NAME = 2,
  USER_FIELD_NAME_LAST_NAME = 3,
  USER_FIELD_NAME_NICK_NAME = 4,
  USER_FIELD_NAME_DISPLAY_NAME = 5,
  USER_FIELD_NAME_EMAIL = 6,
  USER_FIELD_NAME_STATE = 7,
  USER_FIELD_NAME_TYPE = 8,
  USER_FIELD_NAME_CREATION_DATE = 9,
  UNRECOGNIZED = -1,
}

export function userFieldNameFromJSON(object: any): UserFieldName {
  switch (object) {
    case 0:
    case "USER_FIELD_NAME_UNSPECIFIED":
      return UserFieldName.USER_FIELD_NAME_UNSPECIFIED;
    case 1:
    case "USER_FIELD_NAME_USER_NAME":
      return UserFieldName.USER_FIELD_NAME_USER_NAME;
    case 2:
    case "USER_FIELD_NAME_FIRST_NAME":
      return UserFieldName.USER_FIELD_NAME_FIRST_NAME;
    case 3:
    case "USER_FIELD_NAME_LAST_NAME":
      return UserFieldName.USER_FIELD_NAME_LAST_NAME;
    case 4:
    case "USER_FIELD_NAME_NICK_NAME":
      return UserFieldName.USER_FIELD_NAME_NICK_NAME;
    case 5:
    case "USER_FIELD_NAME_DISPLAY_NAME":
      return UserFieldName.USER_FIELD_NAME_DISPLAY_NAME;
    case 6:
    case "USER_FIELD_NAME_EMAIL":
      return UserFieldName.USER_FIELD_NAME_EMAIL;
    case 7:
    case "USER_FIELD_NAME_STATE":
      return UserFieldName.USER_FIELD_NAME_STATE;
    case 8:
    case "USER_FIELD_NAME_TYPE":
      return UserFieldName.USER_FIELD_NAME_TYPE;
    case 9:
    case "USER_FIELD_NAME_CREATION_DATE":
      return UserFieldName.USER_FIELD_NAME_CREATION_DATE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return UserFieldName.UNRECOGNIZED;
  }
}

export function userFieldNameToJSON(object: UserFieldName): string {
  switch (object) {
    case UserFieldName.USER_FIELD_NAME_UNSPECIFIED:
      return "USER_FIELD_NAME_UNSPECIFIED";
    case UserFieldName.USER_FIELD_NAME_USER_NAME:
      return "USER_FIELD_NAME_USER_NAME";
    case UserFieldName.USER_FIELD_NAME_FIRST_NAME:
      return "USER_FIELD_NAME_FIRST_NAME";
    case UserFieldName.USER_FIELD_NAME_LAST_NAME:
      return "USER_FIELD_NAME_LAST_NAME";
    case UserFieldName.USER_FIELD_NAME_NICK_NAME:
      return "USER_FIELD_NAME_NICK_NAME";
    case UserFieldName.USER_FIELD_NAME_DISPLAY_NAME:
      return "USER_FIELD_NAME_DISPLAY_NAME";
    case UserFieldName.USER_FIELD_NAME_EMAIL:
      return "USER_FIELD_NAME_EMAIL";
    case UserFieldName.USER_FIELD_NAME_STATE:
      return "USER_FIELD_NAME_STATE";
    case UserFieldName.USER_FIELD_NAME_TYPE:
      return "USER_FIELD_NAME_TYPE";
    case UserFieldName.USER_FIELD_NAME_CREATION_DATE:
      return "USER_FIELD_NAME_CREATION_DATE";
    case UserFieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface SearchQuery {
  userNameQuery?: UserNameQuery | undefined;
  firstNameQuery?: FirstNameQuery | undefined;
  lastNameQuery?: LastNameQuery | undefined;
  nickNameQuery?: NickNameQuery | undefined;
  displayNameQuery?: DisplayNameQuery | undefined;
  emailQuery?: EmailQuery | undefined;
  stateQuery?: StateQuery | undefined;
  typeQuery?: TypeQuery | undefined;
  loginNameQuery?: LoginNameQuery | undefined;
  inUserIdsQuery?: InUserIDQuery | undefined;
  orQuery?: OrQuery | undefined;
  andQuery?: AndQuery | undefined;
  notQuery?: NotQuery | undefined;
  inUserEmailsQuery?: InUserEmailsQuery | undefined;
  organizationIdQuery?: OrganizationIdQuery | undefined;
}

/** Connect multiple sub-condition with and OR operator. */
export interface OrQuery {
  queries: SearchQuery[];
}

/** Connect multiple sub-condition with and AND operator. */
export interface AndQuery {
  queries: SearchQuery[];
}

/** Negate the sub-condition. */
export interface NotQuery {
  query: SearchQuery | undefined;
}

/** Query for users with ID in list of IDs. */
export interface InUserIDQuery {
  userIds: string[];
}

/** Query for users with a specific user name. */
export interface UserNameQuery {
  userName: string;
  method: TextQueryMethod;
}

/** Query for users with a specific first name. */
export interface FirstNameQuery {
  firstName: string;
  method: TextQueryMethod;
}

/** Query for users with a specific last name. */
export interface LastNameQuery {
  lastName: string;
  method: TextQueryMethod;
}

/** Query for users with a specific nickname. */
export interface NickNameQuery {
  nickName: string;
  method: TextQueryMethod;
}

/** Query for users with a specific display name. */
export interface DisplayNameQuery {
  displayName: string;
  method: TextQueryMethod;
}

/** Query for users with a specific email. */
export interface EmailQuery {
  emailAddress: string;
  method: TextQueryMethod;
}

/** Query for users with a specific state. */
export interface LoginNameQuery {
  loginName: string;
  method: TextQueryMethod;
}

/** Query for users with a specific state. */
export interface StateQuery {
  state: UserState;
}

/** Query for users with a specific type. */
export interface TypeQuery {
  type: Type;
}

/** Query for users with email in list of emails. */
export interface InUserEmailsQuery {
  userEmails: string[];
}

/** Query for users under a specific organization as resource owner. */
export interface OrganizationIdQuery {
  organizationId: string;
}

function createBaseSearchQuery(): SearchQuery {
  return {
    userNameQuery: undefined,
    firstNameQuery: undefined,
    lastNameQuery: undefined,
    nickNameQuery: undefined,
    displayNameQuery: undefined,
    emailQuery: undefined,
    stateQuery: undefined,
    typeQuery: undefined,
    loginNameQuery: undefined,
    inUserIdsQuery: undefined,
    orQuery: undefined,
    andQuery: undefined,
    notQuery: undefined,
    inUserEmailsQuery: undefined,
    organizationIdQuery: undefined,
  };
}

export const SearchQuery = {
  encode(message: SearchQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userNameQuery !== undefined) {
      UserNameQuery.encode(message.userNameQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.firstNameQuery !== undefined) {
      FirstNameQuery.encode(message.firstNameQuery, writer.uint32(18).fork()).ldelim();
    }
    if (message.lastNameQuery !== undefined) {
      LastNameQuery.encode(message.lastNameQuery, writer.uint32(26).fork()).ldelim();
    }
    if (message.nickNameQuery !== undefined) {
      NickNameQuery.encode(message.nickNameQuery, writer.uint32(34).fork()).ldelim();
    }
    if (message.displayNameQuery !== undefined) {
      DisplayNameQuery.encode(message.displayNameQuery, writer.uint32(42).fork()).ldelim();
    }
    if (message.emailQuery !== undefined) {
      EmailQuery.encode(message.emailQuery, writer.uint32(50).fork()).ldelim();
    }
    if (message.stateQuery !== undefined) {
      StateQuery.encode(message.stateQuery, writer.uint32(58).fork()).ldelim();
    }
    if (message.typeQuery !== undefined) {
      TypeQuery.encode(message.typeQuery, writer.uint32(66).fork()).ldelim();
    }
    if (message.loginNameQuery !== undefined) {
      LoginNameQuery.encode(message.loginNameQuery, writer.uint32(74).fork()).ldelim();
    }
    if (message.inUserIdsQuery !== undefined) {
      InUserIDQuery.encode(message.inUserIdsQuery, writer.uint32(82).fork()).ldelim();
    }
    if (message.orQuery !== undefined) {
      OrQuery.encode(message.orQuery, writer.uint32(90).fork()).ldelim();
    }
    if (message.andQuery !== undefined) {
      AndQuery.encode(message.andQuery, writer.uint32(98).fork()).ldelim();
    }
    if (message.notQuery !== undefined) {
      NotQuery.encode(message.notQuery, writer.uint32(106).fork()).ldelim();
    }
    if (message.inUserEmailsQuery !== undefined) {
      InUserEmailsQuery.encode(message.inUserEmailsQuery, writer.uint32(114).fork()).ldelim();
    }
    if (message.organizationIdQuery !== undefined) {
      OrganizationIdQuery.encode(message.organizationIdQuery, writer.uint32(122).fork()).ldelim();
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

          message.userNameQuery = UserNameQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.firstNameQuery = FirstNameQuery.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.lastNameQuery = LastNameQuery.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.nickNameQuery = NickNameQuery.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.displayNameQuery = DisplayNameQuery.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.emailQuery = EmailQuery.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.stateQuery = StateQuery.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.typeQuery = TypeQuery.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.loginNameQuery = LoginNameQuery.decode(reader, reader.uint32());
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          message.inUserIdsQuery = InUserIDQuery.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag != 90) {
            break;
          }

          message.orQuery = OrQuery.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag != 98) {
            break;
          }

          message.andQuery = AndQuery.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag != 106) {
            break;
          }

          message.notQuery = NotQuery.decode(reader, reader.uint32());
          continue;
        case 14:
          if (tag != 114) {
            break;
          }

          message.inUserEmailsQuery = InUserEmailsQuery.decode(reader, reader.uint32());
          continue;
        case 15:
          if (tag != 122) {
            break;
          }

          message.organizationIdQuery = OrganizationIdQuery.decode(reader, reader.uint32());
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
      userNameQuery: isSet(object.userNameQuery) ? UserNameQuery.fromJSON(object.userNameQuery) : undefined,
      firstNameQuery: isSet(object.firstNameQuery) ? FirstNameQuery.fromJSON(object.firstNameQuery) : undefined,
      lastNameQuery: isSet(object.lastNameQuery) ? LastNameQuery.fromJSON(object.lastNameQuery) : undefined,
      nickNameQuery: isSet(object.nickNameQuery) ? NickNameQuery.fromJSON(object.nickNameQuery) : undefined,
      displayNameQuery: isSet(object.displayNameQuery) ? DisplayNameQuery.fromJSON(object.displayNameQuery) : undefined,
      emailQuery: isSet(object.emailQuery) ? EmailQuery.fromJSON(object.emailQuery) : undefined,
      stateQuery: isSet(object.stateQuery) ? StateQuery.fromJSON(object.stateQuery) : undefined,
      typeQuery: isSet(object.typeQuery) ? TypeQuery.fromJSON(object.typeQuery) : undefined,
      loginNameQuery: isSet(object.loginNameQuery) ? LoginNameQuery.fromJSON(object.loginNameQuery) : undefined,
      inUserIdsQuery: isSet(object.inUserIdsQuery) ? InUserIDQuery.fromJSON(object.inUserIdsQuery) : undefined,
      orQuery: isSet(object.orQuery) ? OrQuery.fromJSON(object.orQuery) : undefined,
      andQuery: isSet(object.andQuery) ? AndQuery.fromJSON(object.andQuery) : undefined,
      notQuery: isSet(object.notQuery) ? NotQuery.fromJSON(object.notQuery) : undefined,
      inUserEmailsQuery: isSet(object.inUserEmailsQuery)
        ? InUserEmailsQuery.fromJSON(object.inUserEmailsQuery)
        : undefined,
      organizationIdQuery: isSet(object.organizationIdQuery)
        ? OrganizationIdQuery.fromJSON(object.organizationIdQuery)
        : undefined,
    };
  },

  toJSON(message: SearchQuery): unknown {
    const obj: any = {};
    message.userNameQuery !== undefined &&
      (obj.userNameQuery = message.userNameQuery ? UserNameQuery.toJSON(message.userNameQuery) : undefined);
    message.firstNameQuery !== undefined &&
      (obj.firstNameQuery = message.firstNameQuery ? FirstNameQuery.toJSON(message.firstNameQuery) : undefined);
    message.lastNameQuery !== undefined &&
      (obj.lastNameQuery = message.lastNameQuery ? LastNameQuery.toJSON(message.lastNameQuery) : undefined);
    message.nickNameQuery !== undefined &&
      (obj.nickNameQuery = message.nickNameQuery ? NickNameQuery.toJSON(message.nickNameQuery) : undefined);
    message.displayNameQuery !== undefined &&
      (obj.displayNameQuery = message.displayNameQuery ? DisplayNameQuery.toJSON(message.displayNameQuery) : undefined);
    message.emailQuery !== undefined &&
      (obj.emailQuery = message.emailQuery ? EmailQuery.toJSON(message.emailQuery) : undefined);
    message.stateQuery !== undefined &&
      (obj.stateQuery = message.stateQuery ? StateQuery.toJSON(message.stateQuery) : undefined);
    message.typeQuery !== undefined &&
      (obj.typeQuery = message.typeQuery ? TypeQuery.toJSON(message.typeQuery) : undefined);
    message.loginNameQuery !== undefined &&
      (obj.loginNameQuery = message.loginNameQuery ? LoginNameQuery.toJSON(message.loginNameQuery) : undefined);
    message.inUserIdsQuery !== undefined &&
      (obj.inUserIdsQuery = message.inUserIdsQuery ? InUserIDQuery.toJSON(message.inUserIdsQuery) : undefined);
    message.orQuery !== undefined && (obj.orQuery = message.orQuery ? OrQuery.toJSON(message.orQuery) : undefined);
    message.andQuery !== undefined && (obj.andQuery = message.andQuery ? AndQuery.toJSON(message.andQuery) : undefined);
    message.notQuery !== undefined && (obj.notQuery = message.notQuery ? NotQuery.toJSON(message.notQuery) : undefined);
    message.inUserEmailsQuery !== undefined && (obj.inUserEmailsQuery = message.inUserEmailsQuery
      ? InUserEmailsQuery.toJSON(message.inUserEmailsQuery)
      : undefined);
    message.organizationIdQuery !== undefined && (obj.organizationIdQuery = message.organizationIdQuery
      ? OrganizationIdQuery.toJSON(message.organizationIdQuery)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<SearchQuery>): SearchQuery {
    return SearchQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SearchQuery>): SearchQuery {
    const message = createBaseSearchQuery();
    message.userNameQuery = (object.userNameQuery !== undefined && object.userNameQuery !== null)
      ? UserNameQuery.fromPartial(object.userNameQuery)
      : undefined;
    message.firstNameQuery = (object.firstNameQuery !== undefined && object.firstNameQuery !== null)
      ? FirstNameQuery.fromPartial(object.firstNameQuery)
      : undefined;
    message.lastNameQuery = (object.lastNameQuery !== undefined && object.lastNameQuery !== null)
      ? LastNameQuery.fromPartial(object.lastNameQuery)
      : undefined;
    message.nickNameQuery = (object.nickNameQuery !== undefined && object.nickNameQuery !== null)
      ? NickNameQuery.fromPartial(object.nickNameQuery)
      : undefined;
    message.displayNameQuery = (object.displayNameQuery !== undefined && object.displayNameQuery !== null)
      ? DisplayNameQuery.fromPartial(object.displayNameQuery)
      : undefined;
    message.emailQuery = (object.emailQuery !== undefined && object.emailQuery !== null)
      ? EmailQuery.fromPartial(object.emailQuery)
      : undefined;
    message.stateQuery = (object.stateQuery !== undefined && object.stateQuery !== null)
      ? StateQuery.fromPartial(object.stateQuery)
      : undefined;
    message.typeQuery = (object.typeQuery !== undefined && object.typeQuery !== null)
      ? TypeQuery.fromPartial(object.typeQuery)
      : undefined;
    message.loginNameQuery = (object.loginNameQuery !== undefined && object.loginNameQuery !== null)
      ? LoginNameQuery.fromPartial(object.loginNameQuery)
      : undefined;
    message.inUserIdsQuery = (object.inUserIdsQuery !== undefined && object.inUserIdsQuery !== null)
      ? InUserIDQuery.fromPartial(object.inUserIdsQuery)
      : undefined;
    message.orQuery = (object.orQuery !== undefined && object.orQuery !== null)
      ? OrQuery.fromPartial(object.orQuery)
      : undefined;
    message.andQuery = (object.andQuery !== undefined && object.andQuery !== null)
      ? AndQuery.fromPartial(object.andQuery)
      : undefined;
    message.notQuery = (object.notQuery !== undefined && object.notQuery !== null)
      ? NotQuery.fromPartial(object.notQuery)
      : undefined;
    message.inUserEmailsQuery = (object.inUserEmailsQuery !== undefined && object.inUserEmailsQuery !== null)
      ? InUserEmailsQuery.fromPartial(object.inUserEmailsQuery)
      : undefined;
    message.organizationIdQuery = (object.organizationIdQuery !== undefined && object.organizationIdQuery !== null)
      ? OrganizationIdQuery.fromPartial(object.organizationIdQuery)
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

function createBaseInUserIDQuery(): InUserIDQuery {
  return { userIds: [] };
}

export const InUserIDQuery = {
  encode(message: InUserIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.userIds) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InUserIDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInUserIDQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userIds.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): InUserIDQuery {
    return { userIds: Array.isArray(object?.userIds) ? object.userIds.map((e: any) => String(e)) : [] };
  },

  toJSON(message: InUserIDQuery): unknown {
    const obj: any = {};
    if (message.userIds) {
      obj.userIds = message.userIds.map((e) => e);
    } else {
      obj.userIds = [];
    }
    return obj;
  },

  create(base?: DeepPartial<InUserIDQuery>): InUserIDQuery {
    return InUserIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InUserIDQuery>): InUserIDQuery {
    const message = createBaseInUserIDQuery();
    message.userIds = object.userIds?.map((e) => e) || [];
    return message;
  },
};

function createBaseUserNameQuery(): UserNameQuery {
  return { userName: "", method: 0 };
}

export const UserNameQuery = {
  encode(message: UserNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userName !== "") {
      writer.uint32(10).string(message.userName);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userName = reader.string();
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

  fromJSON(object: any): UserNameQuery {
    return {
      userName: isSet(object.userName) ? String(object.userName) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: UserNameQuery): unknown {
    const obj: any = {};
    message.userName !== undefined && (obj.userName = message.userName);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<UserNameQuery>): UserNameQuery {
    return UserNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserNameQuery>): UserNameQuery {
    const message = createBaseUserNameQuery();
    message.userName = object.userName ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseFirstNameQuery(): FirstNameQuery {
  return { firstName: "", method: 0 };
}

export const FirstNameQuery = {
  encode(message: FirstNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.firstName !== "") {
      writer.uint32(10).string(message.firstName);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FirstNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFirstNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.firstName = reader.string();
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

  fromJSON(object: any): FirstNameQuery {
    return {
      firstName: isSet(object.firstName) ? String(object.firstName) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: FirstNameQuery): unknown {
    const obj: any = {};
    message.firstName !== undefined && (obj.firstName = message.firstName);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<FirstNameQuery>): FirstNameQuery {
    return FirstNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<FirstNameQuery>): FirstNameQuery {
    const message = createBaseFirstNameQuery();
    message.firstName = object.firstName ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseLastNameQuery(): LastNameQuery {
  return { lastName: "", method: 0 };
}

export const LastNameQuery = {
  encode(message: LastNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.lastName !== "") {
      writer.uint32(10).string(message.lastName);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LastNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLastNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.lastName = reader.string();
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

  fromJSON(object: any): LastNameQuery {
    return {
      lastName: isSet(object.lastName) ? String(object.lastName) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: LastNameQuery): unknown {
    const obj: any = {};
    message.lastName !== undefined && (obj.lastName = message.lastName);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<LastNameQuery>): LastNameQuery {
    return LastNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LastNameQuery>): LastNameQuery {
    const message = createBaseLastNameQuery();
    message.lastName = object.lastName ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseNickNameQuery(): NickNameQuery {
  return { nickName: "", method: 0 };
}

export const NickNameQuery = {
  encode(message: NickNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nickName !== "") {
      writer.uint32(10).string(message.nickName);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NickNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNickNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.nickName = reader.string();
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

  fromJSON(object: any): NickNameQuery {
    return {
      nickName: isSet(object.nickName) ? String(object.nickName) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: NickNameQuery): unknown {
    const obj: any = {};
    message.nickName !== undefined && (obj.nickName = message.nickName);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<NickNameQuery>): NickNameQuery {
    return NickNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<NickNameQuery>): NickNameQuery {
    const message = createBaseNickNameQuery();
    message.nickName = object.nickName ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseDisplayNameQuery(): DisplayNameQuery {
  return { displayName: "", method: 0 };
}

export const DisplayNameQuery = {
  encode(message: DisplayNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.displayName !== "") {
      writer.uint32(10).string(message.displayName);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DisplayNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDisplayNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.displayName = reader.string();
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

  fromJSON(object: any): DisplayNameQuery {
    return {
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: DisplayNameQuery): unknown {
    const obj: any = {};
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<DisplayNameQuery>): DisplayNameQuery {
    return DisplayNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DisplayNameQuery>): DisplayNameQuery {
    const message = createBaseDisplayNameQuery();
    message.displayName = object.displayName ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseEmailQuery(): EmailQuery {
  return { emailAddress: "", method: 0 };
}

export const EmailQuery = {
  encode(message: EmailQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.emailAddress !== "") {
      writer.uint32(10).string(message.emailAddress);
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

          message.emailAddress = reader.string();
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
      emailAddress: isSet(object.emailAddress) ? String(object.emailAddress) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: EmailQuery): unknown {
    const obj: any = {};
    message.emailAddress !== undefined && (obj.emailAddress = message.emailAddress);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<EmailQuery>): EmailQuery {
    return EmailQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EmailQuery>): EmailQuery {
    const message = createBaseEmailQuery();
    message.emailAddress = object.emailAddress ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseLoginNameQuery(): LoginNameQuery {
  return { loginName: "", method: 0 };
}

export const LoginNameQuery = {
  encode(message: LoginNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.loginName !== "") {
      writer.uint32(10).string(message.loginName);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LoginNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.loginName = reader.string();
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

  fromJSON(object: any): LoginNameQuery {
    return {
      loginName: isSet(object.loginName) ? String(object.loginName) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: LoginNameQuery): unknown {
    const obj: any = {};
    message.loginName !== undefined && (obj.loginName = message.loginName);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<LoginNameQuery>): LoginNameQuery {
    return LoginNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LoginNameQuery>): LoginNameQuery {
    const message = createBaseLoginNameQuery();
    message.loginName = object.loginName ?? "";
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
    return { state: isSet(object.state) ? userStateFromJSON(object.state) : 0 };
  },

  toJSON(message: StateQuery): unknown {
    const obj: any = {};
    message.state !== undefined && (obj.state = userStateToJSON(message.state));
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

function createBaseTypeQuery(): TypeQuery {
  return { type: 0 };
}

export const TypeQuery = {
  encode(message: TypeQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
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
          if (tag != 8) {
            break;
          }

          message.type = reader.int32() as any;
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
    return { type: isSet(object.type) ? typeFromJSON(object.type) : 0 };
  },

  toJSON(message: TypeQuery): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = typeToJSON(message.type));
    return obj;
  },

  create(base?: DeepPartial<TypeQuery>): TypeQuery {
    return TypeQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TypeQuery>): TypeQuery {
    const message = createBaseTypeQuery();
    message.type = object.type ?? 0;
    return message;
  },
};

function createBaseInUserEmailsQuery(): InUserEmailsQuery {
  return { userEmails: [] };
}

export const InUserEmailsQuery = {
  encode(message: InUserEmailsQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.userEmails) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InUserEmailsQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInUserEmailsQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userEmails.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): InUserEmailsQuery {
    return { userEmails: Array.isArray(object?.userEmails) ? object.userEmails.map((e: any) => String(e)) : [] };
  },

  toJSON(message: InUserEmailsQuery): unknown {
    const obj: any = {};
    if (message.userEmails) {
      obj.userEmails = message.userEmails.map((e) => e);
    } else {
      obj.userEmails = [];
    }
    return obj;
  },

  create(base?: DeepPartial<InUserEmailsQuery>): InUserEmailsQuery {
    return InUserEmailsQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InUserEmailsQuery>): InUserEmailsQuery {
    const message = createBaseInUserEmailsQuery();
    message.userEmails = object.userEmails?.map((e) => e) || [];
    return message;
  },
};

function createBaseOrganizationIdQuery(): OrganizationIdQuery {
  return { organizationId: "" };
}

export const OrganizationIdQuery = {
  encode(message: OrganizationIdQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.organizationId !== "") {
      writer.uint32(10).string(message.organizationId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrganizationIdQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrganizationIdQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.organizationId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OrganizationIdQuery {
    return { organizationId: isSet(object.organizationId) ? String(object.organizationId) : "" };
  },

  toJSON(message: OrganizationIdQuery): unknown {
    const obj: any = {};
    message.organizationId !== undefined && (obj.organizationId = message.organizationId);
    return obj;
  },

  create(base?: DeepPartial<OrganizationIdQuery>): OrganizationIdQuery {
    return OrganizationIdQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OrganizationIdQuery>): OrganizationIdQuery {
    const message = createBaseOrganizationIdQuery();
    message.organizationId = object.organizationId ?? "";
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

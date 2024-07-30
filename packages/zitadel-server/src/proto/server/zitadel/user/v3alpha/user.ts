/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Struct } from "../../../google/protobuf/struct";
import { Details } from "../../object/v2beta/object";
import { Authenticators } from "./authenticator";
import { Contact } from "./communication";

export const protobufPackage = "zitadel.user.v3alpha";

export enum State {
  USER_STATE_UNSPECIFIED = 0,
  USER_STATE_ACTIVE = 1,
  USER_STATE_INACTIVE = 2,
  USER_STATE_DELETED = 3,
  USER_STATE_LOCKED = 4,
  UNRECOGNIZED = -1,
}

export function stateFromJSON(object: any): State {
  switch (object) {
    case 0:
    case "USER_STATE_UNSPECIFIED":
      return State.USER_STATE_UNSPECIFIED;
    case 1:
    case "USER_STATE_ACTIVE":
      return State.USER_STATE_ACTIVE;
    case 2:
    case "USER_STATE_INACTIVE":
      return State.USER_STATE_INACTIVE;
    case 3:
    case "USER_STATE_DELETED":
      return State.USER_STATE_DELETED;
    case 4:
    case "USER_STATE_LOCKED":
      return State.USER_STATE_LOCKED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return State.UNRECOGNIZED;
  }
}

export function stateToJSON(object: State): string {
  switch (object) {
    case State.USER_STATE_UNSPECIFIED:
      return "USER_STATE_UNSPECIFIED";
    case State.USER_STATE_ACTIVE:
      return "USER_STATE_ACTIVE";
    case State.USER_STATE_INACTIVE:
      return "USER_STATE_INACTIVE";
    case State.USER_STATE_DELETED:
      return "USER_STATE_DELETED";
    case State.USER_STATE_LOCKED:
      return "USER_STATE_LOCKED";
    case State.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface User {
  /** ID is the read-only unique identifier of the user. */
  userId: string;
  /** Details provide some base information (such as the last change date) of the user. */
  details:
    | Details
    | undefined;
  /**
   * The user's authenticators. They are used to identify and authenticate the user
   * during the authentication process.
   */
  authenticators:
    | Authenticators
    | undefined;
  /** Contact information for the user. ZITADEL will use this in case of internal notifications. */
  contact:
    | Contact
    | undefined;
  /** State of the user. */
  state: State;
  /** The schema the user and it's data is based on. */
  schema:
    | Schema
    | undefined;
  /** The user's data based on the provided schema. */
  data: { [key: string]: any } | undefined;
}

export interface Schema {
  /** The unique identifier of the user schema. */
  id: string;
  /** The human readable name of the user schema. */
  type: string;
  /** The revision the user's data is based on of the revision. */
  revision: number;
}

function createBaseUser(): User {
  return {
    userId: "",
    details: undefined,
    authenticators: undefined,
    contact: undefined,
    state: 0,
    schema: undefined,
    data: undefined,
  };
}

export const User = {
  encode(message: User, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.authenticators !== undefined) {
      Authenticators.encode(message.authenticators, writer.uint32(26).fork()).ldelim();
    }
    if (message.contact !== undefined) {
      Contact.encode(message.contact, writer.uint32(34).fork()).ldelim();
    }
    if (message.state !== 0) {
      writer.uint32(40).int32(message.state);
    }
    if (message.schema !== undefined) {
      Schema.encode(message.schema, writer.uint32(50).fork()).ldelim();
    }
    if (message.data !== undefined) {
      Struct.encode(Struct.wrap(message.data), writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): User {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUser();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
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

          message.authenticators = Authenticators.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.contact = Contact.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.schema = Schema.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.data = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): User {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      authenticators: isSet(object.authenticators) ? Authenticators.fromJSON(object.authenticators) : undefined,
      contact: isSet(object.contact) ? Contact.fromJSON(object.contact) : undefined,
      state: isSet(object.state) ? stateFromJSON(object.state) : 0,
      schema: isSet(object.schema) ? Schema.fromJSON(object.schema) : undefined,
      data: isObject(object.data) ? object.data : undefined,
    };
  },

  toJSON(message: User): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.authenticators !== undefined &&
      (obj.authenticators = message.authenticators ? Authenticators.toJSON(message.authenticators) : undefined);
    message.contact !== undefined && (obj.contact = message.contact ? Contact.toJSON(message.contact) : undefined);
    message.state !== undefined && (obj.state = stateToJSON(message.state));
    message.schema !== undefined && (obj.schema = message.schema ? Schema.toJSON(message.schema) : undefined);
    message.data !== undefined && (obj.data = message.data);
    return obj;
  },

  create(base?: DeepPartial<User>): User {
    return User.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<User>): User {
    const message = createBaseUser();
    message.userId = object.userId ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.authenticators = (object.authenticators !== undefined && object.authenticators !== null)
      ? Authenticators.fromPartial(object.authenticators)
      : undefined;
    message.contact = (object.contact !== undefined && object.contact !== null)
      ? Contact.fromPartial(object.contact)
      : undefined;
    message.state = object.state ?? 0;
    message.schema = (object.schema !== undefined && object.schema !== null)
      ? Schema.fromPartial(object.schema)
      : undefined;
    message.data = object.data ?? undefined;
    return message;
  },
};

function createBaseSchema(): Schema {
  return { id: "", type: "", revision: 0 };
}

export const Schema = {
  encode(message: Schema, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.type !== "") {
      writer.uint32(18).string(message.type);
    }
    if (message.revision !== 0) {
      writer.uint32(24).uint32(message.revision);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Schema {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSchema();
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
          if (tag != 24) {
            break;
          }

          message.revision = reader.uint32();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Schema {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      type: isSet(object.type) ? String(object.type) : "",
      revision: isSet(object.revision) ? Number(object.revision) : 0,
    };
  },

  toJSON(message: Schema): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.type !== undefined && (obj.type = message.type);
    message.revision !== undefined && (obj.revision = Math.round(message.revision));
    return obj;
  },

  create(base?: DeepPartial<Schema>): Schema {
    return Schema.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Schema>): Schema {
    const message = createBaseSchema();
    message.id = object.id ?? "";
    message.type = object.type ?? "";
    message.revision = object.revision ?? 0;
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

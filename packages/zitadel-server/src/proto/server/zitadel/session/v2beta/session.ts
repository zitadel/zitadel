/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { TimestampQueryMethod, timestampQueryMethodFromJSON, timestampQueryMethodToJSON } from "../../object";

export const protobufPackage = "zitadel.session.v2beta";

export enum SessionFieldName {
  SESSION_FIELD_NAME_UNSPECIFIED = 0,
  SESSION_FIELD_NAME_CREATION_DATE = 1,
  UNRECOGNIZED = -1,
}

export function sessionFieldNameFromJSON(object: any): SessionFieldName {
  switch (object) {
    case 0:
    case "SESSION_FIELD_NAME_UNSPECIFIED":
      return SessionFieldName.SESSION_FIELD_NAME_UNSPECIFIED;
    case 1:
    case "SESSION_FIELD_NAME_CREATION_DATE":
      return SessionFieldName.SESSION_FIELD_NAME_CREATION_DATE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SessionFieldName.UNRECOGNIZED;
  }
}

export function sessionFieldNameToJSON(object: SessionFieldName): string {
  switch (object) {
    case SessionFieldName.SESSION_FIELD_NAME_UNSPECIFIED:
      return "SESSION_FIELD_NAME_UNSPECIFIED";
    case SessionFieldName.SESSION_FIELD_NAME_CREATION_DATE:
      return "SESSION_FIELD_NAME_CREATION_DATE";
    case SessionFieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Session {
  id: string;
  creationDate: Date | undefined;
  changeDate: Date | undefined;
  sequence: number;
  factors: Factors | undefined;
  metadata: { [key: string]: Buffer };
  userAgent: UserAgent | undefined;
  expirationDate?: Date | undefined;
}

export interface Session_MetadataEntry {
  key: string;
  value: Buffer;
}

export interface Factors {
  user: UserFactor | undefined;
  password: PasswordFactor | undefined;
  webAuthN: WebAuthNFactor | undefined;
  intent: IntentFactor | undefined;
  totp: TOTPFactor | undefined;
  otpSms: OTPFactor | undefined;
  otpEmail: OTPFactor | undefined;
}

export interface UserFactor {
  verifiedAt: Date | undefined;
  id: string;
  loginName: string;
  displayName: string;
  organizationId: string;
}

export interface PasswordFactor {
  verifiedAt: Date | undefined;
}

export interface IntentFactor {
  verifiedAt: Date | undefined;
}

export interface WebAuthNFactor {
  verifiedAt: Date | undefined;
  userVerified: boolean;
}

export interface TOTPFactor {
  verifiedAt: Date | undefined;
}

export interface OTPFactor {
  verifiedAt: Date | undefined;
}

export interface SearchQuery {
  idsQuery?: IDsQuery | undefined;
  userIdQuery?: UserIDQuery | undefined;
  creationDateQuery?: CreationDateQuery | undefined;
}

export interface IDsQuery {
  ids: string[];
}

export interface UserIDQuery {
  id: string;
}

export interface CreationDateQuery {
  creationDate: Date | undefined;
  method: TimestampQueryMethod;
}

export interface UserAgent {
  fingerprintId?: string | undefined;
  ip?: string | undefined;
  description?: string | undefined;
  header: { [key: string]: UserAgent_HeaderValues };
}

/**
 * A header may have multiple values.
 * In Go, headers are defined
 * as map[string][]string, but protobuf
 * doesn't allow this scheme.
 */
export interface UserAgent_HeaderValues {
  values: string[];
}

export interface UserAgent_HeaderEntry {
  key: string;
  value: UserAgent_HeaderValues | undefined;
}

function createBaseSession(): Session {
  return {
    id: "",
    creationDate: undefined,
    changeDate: undefined,
    sequence: 0,
    factors: undefined,
    metadata: {},
    userAgent: undefined,
    expirationDate: undefined,
  };
}

export const Session = {
  encode(message: Session, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.creationDate !== undefined) {
      Timestamp.encode(toTimestamp(message.creationDate), writer.uint32(18).fork()).ldelim();
    }
    if (message.changeDate !== undefined) {
      Timestamp.encode(toTimestamp(message.changeDate), writer.uint32(26).fork()).ldelim();
    }
    if (message.sequence !== 0) {
      writer.uint32(32).uint64(message.sequence);
    }
    if (message.factors !== undefined) {
      Factors.encode(message.factors, writer.uint32(42).fork()).ldelim();
    }
    Object.entries(message.metadata).forEach(([key, value]) => {
      Session_MetadataEntry.encode({ key: key as any, value }, writer.uint32(50).fork()).ldelim();
    });
    if (message.userAgent !== undefined) {
      UserAgent.encode(message.userAgent, writer.uint32(58).fork()).ldelim();
    }
    if (message.expirationDate !== undefined) {
      Timestamp.encode(toTimestamp(message.expirationDate), writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Session {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSession();
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

          message.creationDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.changeDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.sequence = longToNumber(reader.uint64() as Long);
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.factors = Factors.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          const entry6 = Session_MetadataEntry.decode(reader, reader.uint32());
          if (entry6.value !== undefined) {
            message.metadata[entry6.key] = entry6.value;
          }
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.userAgent = UserAgent.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.expirationDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Session {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      creationDate: isSet(object.creationDate) ? fromJsonTimestamp(object.creationDate) : undefined,
      changeDate: isSet(object.changeDate) ? fromJsonTimestamp(object.changeDate) : undefined,
      sequence: isSet(object.sequence) ? Number(object.sequence) : 0,
      factors: isSet(object.factors) ? Factors.fromJSON(object.factors) : undefined,
      metadata: isObject(object.metadata)
        ? Object.entries(object.metadata).reduce<{ [key: string]: Buffer }>((acc, [key, value]) => {
          acc[key] = Buffer.from(bytesFromBase64(value as string));
          return acc;
        }, {})
        : {},
      userAgent: isSet(object.userAgent) ? UserAgent.fromJSON(object.userAgent) : undefined,
      expirationDate: isSet(object.expirationDate) ? fromJsonTimestamp(object.expirationDate) : undefined,
    };
  },

  toJSON(message: Session): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.creationDate !== undefined && (obj.creationDate = message.creationDate.toISOString());
    message.changeDate !== undefined && (obj.changeDate = message.changeDate.toISOString());
    message.sequence !== undefined && (obj.sequence = Math.round(message.sequence));
    message.factors !== undefined && (obj.factors = message.factors ? Factors.toJSON(message.factors) : undefined);
    obj.metadata = {};
    if (message.metadata) {
      Object.entries(message.metadata).forEach(([k, v]) => {
        obj.metadata[k] = base64FromBytes(v);
      });
    }
    message.userAgent !== undefined &&
      (obj.userAgent = message.userAgent ? UserAgent.toJSON(message.userAgent) : undefined);
    message.expirationDate !== undefined && (obj.expirationDate = message.expirationDate.toISOString());
    return obj;
  },

  create(base?: DeepPartial<Session>): Session {
    return Session.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Session>): Session {
    const message = createBaseSession();
    message.id = object.id ?? "";
    message.creationDate = object.creationDate ?? undefined;
    message.changeDate = object.changeDate ?? undefined;
    message.sequence = object.sequence ?? 0;
    message.factors = (object.factors !== undefined && object.factors !== null)
      ? Factors.fromPartial(object.factors)
      : undefined;
    message.metadata = Object.entries(object.metadata ?? {}).reduce<{ [key: string]: Buffer }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = value;
      }
      return acc;
    }, {});
    message.userAgent = (object.userAgent !== undefined && object.userAgent !== null)
      ? UserAgent.fromPartial(object.userAgent)
      : undefined;
    message.expirationDate = object.expirationDate ?? undefined;
    return message;
  },
};

function createBaseSession_MetadataEntry(): Session_MetadataEntry {
  return { key: "", value: Buffer.alloc(0) };
}

export const Session_MetadataEntry = {
  encode(message: Session_MetadataEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value.length !== 0) {
      writer.uint32(18).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Session_MetadataEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSession_MetadataEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = reader.bytes() as Buffer;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Session_MetadataEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Buffer.from(bytesFromBase64(object.value)) : Buffer.alloc(0),
    };
  },

  toJSON(message: Session_MetadataEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = base64FromBytes(message.value !== undefined ? message.value : Buffer.alloc(0)));
    return obj;
  },

  create(base?: DeepPartial<Session_MetadataEntry>): Session_MetadataEntry {
    return Session_MetadataEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Session_MetadataEntry>): Session_MetadataEntry {
    const message = createBaseSession_MetadataEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseFactors(): Factors {
  return {
    user: undefined,
    password: undefined,
    webAuthN: undefined,
    intent: undefined,
    totp: undefined,
    otpSms: undefined,
    otpEmail: undefined,
  };
}

export const Factors = {
  encode(message: Factors, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      UserFactor.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    if (message.password !== undefined) {
      PasswordFactor.encode(message.password, writer.uint32(18).fork()).ldelim();
    }
    if (message.webAuthN !== undefined) {
      WebAuthNFactor.encode(message.webAuthN, writer.uint32(26).fork()).ldelim();
    }
    if (message.intent !== undefined) {
      IntentFactor.encode(message.intent, writer.uint32(34).fork()).ldelim();
    }
    if (message.totp !== undefined) {
      TOTPFactor.encode(message.totp, writer.uint32(42).fork()).ldelim();
    }
    if (message.otpSms !== undefined) {
      OTPFactor.encode(message.otpSms, writer.uint32(50).fork()).ldelim();
    }
    if (message.otpEmail !== undefined) {
      OTPFactor.encode(message.otpEmail, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Factors {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFactors();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.user = UserFactor.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.password = PasswordFactor.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.webAuthN = WebAuthNFactor.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.intent = IntentFactor.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.totp = TOTPFactor.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.otpSms = OTPFactor.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.otpEmail = OTPFactor.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Factors {
    return {
      user: isSet(object.user) ? UserFactor.fromJSON(object.user) : undefined,
      password: isSet(object.password) ? PasswordFactor.fromJSON(object.password) : undefined,
      webAuthN: isSet(object.webAuthN) ? WebAuthNFactor.fromJSON(object.webAuthN) : undefined,
      intent: isSet(object.intent) ? IntentFactor.fromJSON(object.intent) : undefined,
      totp: isSet(object.totp) ? TOTPFactor.fromJSON(object.totp) : undefined,
      otpSms: isSet(object.otpSms) ? OTPFactor.fromJSON(object.otpSms) : undefined,
      otpEmail: isSet(object.otpEmail) ? OTPFactor.fromJSON(object.otpEmail) : undefined,
    };
  },

  toJSON(message: Factors): unknown {
    const obj: any = {};
    message.user !== undefined && (obj.user = message.user ? UserFactor.toJSON(message.user) : undefined);
    message.password !== undefined &&
      (obj.password = message.password ? PasswordFactor.toJSON(message.password) : undefined);
    message.webAuthN !== undefined &&
      (obj.webAuthN = message.webAuthN ? WebAuthNFactor.toJSON(message.webAuthN) : undefined);
    message.intent !== undefined && (obj.intent = message.intent ? IntentFactor.toJSON(message.intent) : undefined);
    message.totp !== undefined && (obj.totp = message.totp ? TOTPFactor.toJSON(message.totp) : undefined);
    message.otpSms !== undefined && (obj.otpSms = message.otpSms ? OTPFactor.toJSON(message.otpSms) : undefined);
    message.otpEmail !== undefined &&
      (obj.otpEmail = message.otpEmail ? OTPFactor.toJSON(message.otpEmail) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Factors>): Factors {
    return Factors.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Factors>): Factors {
    const message = createBaseFactors();
    message.user = (object.user !== undefined && object.user !== null)
      ? UserFactor.fromPartial(object.user)
      : undefined;
    message.password = (object.password !== undefined && object.password !== null)
      ? PasswordFactor.fromPartial(object.password)
      : undefined;
    message.webAuthN = (object.webAuthN !== undefined && object.webAuthN !== null)
      ? WebAuthNFactor.fromPartial(object.webAuthN)
      : undefined;
    message.intent = (object.intent !== undefined && object.intent !== null)
      ? IntentFactor.fromPartial(object.intent)
      : undefined;
    message.totp = (object.totp !== undefined && object.totp !== null)
      ? TOTPFactor.fromPartial(object.totp)
      : undefined;
    message.otpSms = (object.otpSms !== undefined && object.otpSms !== null)
      ? OTPFactor.fromPartial(object.otpSms)
      : undefined;
    message.otpEmail = (object.otpEmail !== undefined && object.otpEmail !== null)
      ? OTPFactor.fromPartial(object.otpEmail)
      : undefined;
    return message;
  },
};

function createBaseUserFactor(): UserFactor {
  return { verifiedAt: undefined, id: "", loginName: "", displayName: "", organizationId: "" };
}

export const UserFactor = {
  encode(message: UserFactor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.verifiedAt !== undefined) {
      Timestamp.encode(toTimestamp(message.verifiedAt), writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    if (message.loginName !== "") {
      writer.uint32(26).string(message.loginName);
    }
    if (message.displayName !== "") {
      writer.uint32(34).string(message.displayName);
    }
    if (message.organizationId !== "") {
      writer.uint32(50).string(message.organizationId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserFactor {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserFactor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.verifiedAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.id = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.loginName = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.displayName = reader.string();
          continue;
        case 6:
          if (tag != 50) {
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

  fromJSON(object: any): UserFactor {
    return {
      verifiedAt: isSet(object.verifiedAt) ? fromJsonTimestamp(object.verifiedAt) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
      loginName: isSet(object.loginName) ? String(object.loginName) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      organizationId: isSet(object.organizationId) ? String(object.organizationId) : "",
    };
  },

  toJSON(message: UserFactor): unknown {
    const obj: any = {};
    message.verifiedAt !== undefined && (obj.verifiedAt = message.verifiedAt.toISOString());
    message.id !== undefined && (obj.id = message.id);
    message.loginName !== undefined && (obj.loginName = message.loginName);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.organizationId !== undefined && (obj.organizationId = message.organizationId);
    return obj;
  },

  create(base?: DeepPartial<UserFactor>): UserFactor {
    return UserFactor.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserFactor>): UserFactor {
    const message = createBaseUserFactor();
    message.verifiedAt = object.verifiedAt ?? undefined;
    message.id = object.id ?? "";
    message.loginName = object.loginName ?? "";
    message.displayName = object.displayName ?? "";
    message.organizationId = object.organizationId ?? "";
    return message;
  },
};

function createBasePasswordFactor(): PasswordFactor {
  return { verifiedAt: undefined };
}

export const PasswordFactor = {
  encode(message: PasswordFactor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.verifiedAt !== undefined) {
      Timestamp.encode(toTimestamp(message.verifiedAt), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordFactor {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordFactor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.verifiedAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PasswordFactor {
    return { verifiedAt: isSet(object.verifiedAt) ? fromJsonTimestamp(object.verifiedAt) : undefined };
  },

  toJSON(message: PasswordFactor): unknown {
    const obj: any = {};
    message.verifiedAt !== undefined && (obj.verifiedAt = message.verifiedAt.toISOString());
    return obj;
  },

  create(base?: DeepPartial<PasswordFactor>): PasswordFactor {
    return PasswordFactor.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordFactor>): PasswordFactor {
    const message = createBasePasswordFactor();
    message.verifiedAt = object.verifiedAt ?? undefined;
    return message;
  },
};

function createBaseIntentFactor(): IntentFactor {
  return { verifiedAt: undefined };
}

export const IntentFactor = {
  encode(message: IntentFactor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.verifiedAt !== undefined) {
      Timestamp.encode(toTimestamp(message.verifiedAt), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IntentFactor {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIntentFactor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.verifiedAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IntentFactor {
    return { verifiedAt: isSet(object.verifiedAt) ? fromJsonTimestamp(object.verifiedAt) : undefined };
  },

  toJSON(message: IntentFactor): unknown {
    const obj: any = {};
    message.verifiedAt !== undefined && (obj.verifiedAt = message.verifiedAt.toISOString());
    return obj;
  },

  create(base?: DeepPartial<IntentFactor>): IntentFactor {
    return IntentFactor.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IntentFactor>): IntentFactor {
    const message = createBaseIntentFactor();
    message.verifiedAt = object.verifiedAt ?? undefined;
    return message;
  },
};

function createBaseWebAuthNFactor(): WebAuthNFactor {
  return { verifiedAt: undefined, userVerified: false };
}

export const WebAuthNFactor = {
  encode(message: WebAuthNFactor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.verifiedAt !== undefined) {
      Timestamp.encode(toTimestamp(message.verifiedAt), writer.uint32(10).fork()).ldelim();
    }
    if (message.userVerified === true) {
      writer.uint32(16).bool(message.userVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WebAuthNFactor {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWebAuthNFactor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.verifiedAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.userVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): WebAuthNFactor {
    return {
      verifiedAt: isSet(object.verifiedAt) ? fromJsonTimestamp(object.verifiedAt) : undefined,
      userVerified: isSet(object.userVerified) ? Boolean(object.userVerified) : false,
    };
  },

  toJSON(message: WebAuthNFactor): unknown {
    const obj: any = {};
    message.verifiedAt !== undefined && (obj.verifiedAt = message.verifiedAt.toISOString());
    message.userVerified !== undefined && (obj.userVerified = message.userVerified);
    return obj;
  },

  create(base?: DeepPartial<WebAuthNFactor>): WebAuthNFactor {
    return WebAuthNFactor.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<WebAuthNFactor>): WebAuthNFactor {
    const message = createBaseWebAuthNFactor();
    message.verifiedAt = object.verifiedAt ?? undefined;
    message.userVerified = object.userVerified ?? false;
    return message;
  },
};

function createBaseTOTPFactor(): TOTPFactor {
  return { verifiedAt: undefined };
}

export const TOTPFactor = {
  encode(message: TOTPFactor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.verifiedAt !== undefined) {
      Timestamp.encode(toTimestamp(message.verifiedAt), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TOTPFactor {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTOTPFactor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.verifiedAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TOTPFactor {
    return { verifiedAt: isSet(object.verifiedAt) ? fromJsonTimestamp(object.verifiedAt) : undefined };
  },

  toJSON(message: TOTPFactor): unknown {
    const obj: any = {};
    message.verifiedAt !== undefined && (obj.verifiedAt = message.verifiedAt.toISOString());
    return obj;
  },

  create(base?: DeepPartial<TOTPFactor>): TOTPFactor {
    return TOTPFactor.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TOTPFactor>): TOTPFactor {
    const message = createBaseTOTPFactor();
    message.verifiedAt = object.verifiedAt ?? undefined;
    return message;
  },
};

function createBaseOTPFactor(): OTPFactor {
  return { verifiedAt: undefined };
}

export const OTPFactor = {
  encode(message: OTPFactor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.verifiedAt !== undefined) {
      Timestamp.encode(toTimestamp(message.verifiedAt), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OTPFactor {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOTPFactor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.verifiedAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OTPFactor {
    return { verifiedAt: isSet(object.verifiedAt) ? fromJsonTimestamp(object.verifiedAt) : undefined };
  },

  toJSON(message: OTPFactor): unknown {
    const obj: any = {};
    message.verifiedAt !== undefined && (obj.verifiedAt = message.verifiedAt.toISOString());
    return obj;
  },

  create(base?: DeepPartial<OTPFactor>): OTPFactor {
    return OTPFactor.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OTPFactor>): OTPFactor {
    const message = createBaseOTPFactor();
    message.verifiedAt = object.verifiedAt ?? undefined;
    return message;
  },
};

function createBaseSearchQuery(): SearchQuery {
  return { idsQuery: undefined, userIdQuery: undefined, creationDateQuery: undefined };
}

export const SearchQuery = {
  encode(message: SearchQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idsQuery !== undefined) {
      IDsQuery.encode(message.idsQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.userIdQuery !== undefined) {
      UserIDQuery.encode(message.userIdQuery, writer.uint32(18).fork()).ldelim();
    }
    if (message.creationDateQuery !== undefined) {
      CreationDateQuery.encode(message.creationDateQuery, writer.uint32(26).fork()).ldelim();
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

          message.idsQuery = IDsQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.userIdQuery = UserIDQuery.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.creationDateQuery = CreationDateQuery.decode(reader, reader.uint32());
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
      idsQuery: isSet(object.idsQuery) ? IDsQuery.fromJSON(object.idsQuery) : undefined,
      userIdQuery: isSet(object.userIdQuery) ? UserIDQuery.fromJSON(object.userIdQuery) : undefined,
      creationDateQuery: isSet(object.creationDateQuery)
        ? CreationDateQuery.fromJSON(object.creationDateQuery)
        : undefined,
    };
  },

  toJSON(message: SearchQuery): unknown {
    const obj: any = {};
    message.idsQuery !== undefined && (obj.idsQuery = message.idsQuery ? IDsQuery.toJSON(message.idsQuery) : undefined);
    message.userIdQuery !== undefined &&
      (obj.userIdQuery = message.userIdQuery ? UserIDQuery.toJSON(message.userIdQuery) : undefined);
    message.creationDateQuery !== undefined && (obj.creationDateQuery = message.creationDateQuery
      ? CreationDateQuery.toJSON(message.creationDateQuery)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<SearchQuery>): SearchQuery {
    return SearchQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SearchQuery>): SearchQuery {
    const message = createBaseSearchQuery();
    message.idsQuery = (object.idsQuery !== undefined && object.idsQuery !== null)
      ? IDsQuery.fromPartial(object.idsQuery)
      : undefined;
    message.userIdQuery = (object.userIdQuery !== undefined && object.userIdQuery !== null)
      ? UserIDQuery.fromPartial(object.userIdQuery)
      : undefined;
    message.creationDateQuery = (object.creationDateQuery !== undefined && object.creationDateQuery !== null)
      ? CreationDateQuery.fromPartial(object.creationDateQuery)
      : undefined;
    return message;
  },
};

function createBaseIDsQuery(): IDsQuery {
  return { ids: [] };
}

export const IDsQuery = {
  encode(message: IDsQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.ids) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDsQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDsQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.ids.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDsQuery {
    return { ids: Array.isArray(object?.ids) ? object.ids.map((e: any) => String(e)) : [] };
  },

  toJSON(message: IDsQuery): unknown {
    const obj: any = {};
    if (message.ids) {
      obj.ids = message.ids.map((e) => e);
    } else {
      obj.ids = [];
    }
    return obj;
  },

  create(base?: DeepPartial<IDsQuery>): IDsQuery {
    return IDsQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDsQuery>): IDsQuery {
    const message = createBaseIDsQuery();
    message.ids = object.ids?.map((e) => e) || [];
    return message;
  },
};

function createBaseUserIDQuery(): UserIDQuery {
  return { id: "" };
}

export const UserIDQuery = {
  encode(message: UserIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
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
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UserIDQuery {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: UserIDQuery): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<UserIDQuery>): UserIDQuery {
    return UserIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserIDQuery>): UserIDQuery {
    const message = createBaseUserIDQuery();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseCreationDateQuery(): CreationDateQuery {
  return { creationDate: undefined, method: 0 };
}

export const CreationDateQuery = {
  encode(message: CreationDateQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creationDate !== undefined) {
      Timestamp.encode(toTimestamp(message.creationDate), writer.uint32(10).fork()).ldelim();
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreationDateQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreationDateQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.creationDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
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

  fromJSON(object: any): CreationDateQuery {
    return {
      creationDate: isSet(object.creationDate) ? fromJsonTimestamp(object.creationDate) : undefined,
      method: isSet(object.method) ? timestampQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: CreationDateQuery): unknown {
    const obj: any = {};
    message.creationDate !== undefined && (obj.creationDate = message.creationDate.toISOString());
    message.method !== undefined && (obj.method = timestampQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<CreationDateQuery>): CreationDateQuery {
    return CreationDateQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreationDateQuery>): CreationDateQuery {
    const message = createBaseCreationDateQuery();
    message.creationDate = object.creationDate ?? undefined;
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseUserAgent(): UserAgent {
  return { fingerprintId: undefined, ip: undefined, description: undefined, header: {} };
}

export const UserAgent = {
  encode(message: UserAgent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.fingerprintId !== undefined) {
      writer.uint32(10).string(message.fingerprintId);
    }
    if (message.ip !== undefined) {
      writer.uint32(18).string(message.ip);
    }
    if (message.description !== undefined) {
      writer.uint32(26).string(message.description);
    }
    Object.entries(message.header).forEach(([key, value]) => {
      UserAgent_HeaderEntry.encode({ key: key as any, value }, writer.uint32(34).fork()).ldelim();
    });
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserAgent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserAgent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.fingerprintId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.ip = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.description = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          const entry4 = UserAgent_HeaderEntry.decode(reader, reader.uint32());
          if (entry4.value !== undefined) {
            message.header[entry4.key] = entry4.value;
          }
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UserAgent {
    return {
      fingerprintId: isSet(object.fingerprintId) ? String(object.fingerprintId) : undefined,
      ip: isSet(object.ip) ? String(object.ip) : undefined,
      description: isSet(object.description) ? String(object.description) : undefined,
      header: isObject(object.header)
        ? Object.entries(object.header).reduce<{ [key: string]: UserAgent_HeaderValues }>((acc, [key, value]) => {
          acc[key] = UserAgent_HeaderValues.fromJSON(value);
          return acc;
        }, {})
        : {},
    };
  },

  toJSON(message: UserAgent): unknown {
    const obj: any = {};
    message.fingerprintId !== undefined && (obj.fingerprintId = message.fingerprintId);
    message.ip !== undefined && (obj.ip = message.ip);
    message.description !== undefined && (obj.description = message.description);
    obj.header = {};
    if (message.header) {
      Object.entries(message.header).forEach(([k, v]) => {
        obj.header[k] = UserAgent_HeaderValues.toJSON(v);
      });
    }
    return obj;
  },

  create(base?: DeepPartial<UserAgent>): UserAgent {
    return UserAgent.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserAgent>): UserAgent {
    const message = createBaseUserAgent();
    message.fingerprintId = object.fingerprintId ?? undefined;
    message.ip = object.ip ?? undefined;
    message.description = object.description ?? undefined;
    message.header = Object.entries(object.header ?? {}).reduce<{ [key: string]: UserAgent_HeaderValues }>(
      (acc, [key, value]) => {
        if (value !== undefined) {
          acc[key] = UserAgent_HeaderValues.fromPartial(value);
        }
        return acc;
      },
      {},
    );
    return message;
  },
};

function createBaseUserAgent_HeaderValues(): UserAgent_HeaderValues {
  return { values: [] };
}

export const UserAgent_HeaderValues = {
  encode(message: UserAgent_HeaderValues, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.values) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserAgent_HeaderValues {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserAgent_HeaderValues();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.values.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UserAgent_HeaderValues {
    return { values: Array.isArray(object?.values) ? object.values.map((e: any) => String(e)) : [] };
  },

  toJSON(message: UserAgent_HeaderValues): unknown {
    const obj: any = {};
    if (message.values) {
      obj.values = message.values.map((e) => e);
    } else {
      obj.values = [];
    }
    return obj;
  },

  create(base?: DeepPartial<UserAgent_HeaderValues>): UserAgent_HeaderValues {
    return UserAgent_HeaderValues.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserAgent_HeaderValues>): UserAgent_HeaderValues {
    const message = createBaseUserAgent_HeaderValues();
    message.values = object.values?.map((e) => e) || [];
    return message;
  },
};

function createBaseUserAgent_HeaderEntry(): UserAgent_HeaderEntry {
  return { key: "", value: undefined };
}

export const UserAgent_HeaderEntry = {
  encode(message: UserAgent_HeaderEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value !== undefined) {
      UserAgent_HeaderValues.encode(message.value, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserAgent_HeaderEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserAgent_HeaderEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = UserAgent_HeaderValues.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UserAgent_HeaderEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? UserAgent_HeaderValues.fromJSON(object.value) : undefined,
    };
  },

  toJSON(message: UserAgent_HeaderEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = message.value ? UserAgent_HeaderValues.toJSON(message.value) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UserAgent_HeaderEntry>): UserAgent_HeaderEntry {
    return UserAgent_HeaderEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserAgent_HeaderEntry>): UserAgent_HeaderEntry {
    const message = createBaseUserAgent_HeaderEntry();
    message.key = object.key ?? "";
    message.value = (object.value !== undefined && object.value !== null)
      ? UserAgent_HeaderValues.fromPartial(object.value)
      : undefined;
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

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

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

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

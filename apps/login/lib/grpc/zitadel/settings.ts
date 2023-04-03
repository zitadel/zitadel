/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";
import { ObjectDetails } from "./object";

export const protobufPackage = "zitadel.settings.v1";

export enum SecretGeneratorType {
  SECRET_GENERATOR_TYPE_UNSPECIFIED = 0,
  SECRET_GENERATOR_TYPE_INIT_CODE = 1,
  SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE = 2,
  SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE = 3,
  SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE = 4,
  SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE = 5,
  SECRET_GENERATOR_TYPE_APP_SECRET = 6,
  UNRECOGNIZED = -1,
}

export function secretGeneratorTypeFromJSON(object: any): SecretGeneratorType {
  switch (object) {
    case 0:
    case "SECRET_GENERATOR_TYPE_UNSPECIFIED":
      return SecretGeneratorType.SECRET_GENERATOR_TYPE_UNSPECIFIED;
    case 1:
    case "SECRET_GENERATOR_TYPE_INIT_CODE":
      return SecretGeneratorType.SECRET_GENERATOR_TYPE_INIT_CODE;
    case 2:
    case "SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE":
      return SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE;
    case 3:
    case "SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE":
      return SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE;
    case 4:
    case "SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE":
      return SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE;
    case 5:
    case "SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE":
      return SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE;
    case 6:
    case "SECRET_GENERATOR_TYPE_APP_SECRET":
      return SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SecretGeneratorType.UNRECOGNIZED;
  }
}

export function secretGeneratorTypeToJSON(object: SecretGeneratorType): string {
  switch (object) {
    case SecretGeneratorType.SECRET_GENERATOR_TYPE_UNSPECIFIED:
      return "SECRET_GENERATOR_TYPE_UNSPECIFIED";
    case SecretGeneratorType.SECRET_GENERATOR_TYPE_INIT_CODE:
      return "SECRET_GENERATOR_TYPE_INIT_CODE";
    case SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE:
      return "SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE";
    case SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE:
      return "SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE";
    case SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE:
      return "SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE";
    case SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE:
      return "SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE";
    case SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET:
      return "SECRET_GENERATOR_TYPE_APP_SECRET";
    case SecretGeneratorType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum SMSProviderConfigState {
  SMS_PROVIDER_CONFIG_STATE_UNSPECIFIED = 0,
  SMS_PROVIDER_CONFIG_ACTIVE = 1,
  SMS_PROVIDER_CONFIG_INACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function sMSProviderConfigStateFromJSON(object: any): SMSProviderConfigState {
  switch (object) {
    case 0:
    case "SMS_PROVIDER_CONFIG_STATE_UNSPECIFIED":
      return SMSProviderConfigState.SMS_PROVIDER_CONFIG_STATE_UNSPECIFIED;
    case 1:
    case "SMS_PROVIDER_CONFIG_ACTIVE":
      return SMSProviderConfigState.SMS_PROVIDER_CONFIG_ACTIVE;
    case 2:
    case "SMS_PROVIDER_CONFIG_INACTIVE":
      return SMSProviderConfigState.SMS_PROVIDER_CONFIG_INACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SMSProviderConfigState.UNRECOGNIZED;
  }
}

export function sMSProviderConfigStateToJSON(object: SMSProviderConfigState): string {
  switch (object) {
    case SMSProviderConfigState.SMS_PROVIDER_CONFIG_STATE_UNSPECIFIED:
      return "SMS_PROVIDER_CONFIG_STATE_UNSPECIFIED";
    case SMSProviderConfigState.SMS_PROVIDER_CONFIG_ACTIVE:
      return "SMS_PROVIDER_CONFIG_ACTIVE";
    case SMSProviderConfigState.SMS_PROVIDER_CONFIG_INACTIVE:
      return "SMS_PROVIDER_CONFIG_INACTIVE";
    case SMSProviderConfigState.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface SecretGenerator {
  generatorType: SecretGeneratorType;
  details: ObjectDetails | undefined;
  length: number;
  expiry: Duration | undefined;
  includeLowerLetters: boolean;
  includeUpperLetters: boolean;
  includeDigits: boolean;
  includeSymbols: boolean;
}

export interface SecretGeneratorQuery {
  typeQuery?: SecretGeneratorTypeQuery | undefined;
}

export interface SecretGeneratorTypeQuery {
  generatorType: SecretGeneratorType;
}

export interface SMTPConfig {
  details: ObjectDetails | undefined;
  senderAddress: string;
  senderName: string;
  tls: boolean;
  host: string;
  user: string;
}

export interface SMSProvider {
  details: ObjectDetails | undefined;
  id: string;
  state: SMSProviderConfigState;
  twilio?: TwilioConfig | undefined;
}

export interface TwilioConfig {
  sid: string;
  senderNumber: string;
}

export interface DebugNotificationProvider {
  details: ObjectDetails | undefined;
  compact: boolean;
}

export interface OIDCSettings {
  details: ObjectDetails | undefined;
  accessTokenLifetime: Duration | undefined;
  idTokenLifetime: Duration | undefined;
  refreshTokenIdleExpiration: Duration | undefined;
  refreshTokenExpiration: Duration | undefined;
}

export interface SecurityPolicy {
  details:
    | ObjectDetails
    | undefined;
  /** states if iframe embedding is enabled or disabled */
  enableIframeEmbedding: boolean;
  /** origins allowed loading ZITADEL in an iframe if enable_iframe_embedding is true */
  allowedOrigins: string[];
}

function createBaseSecretGenerator(): SecretGenerator {
  return {
    generatorType: 0,
    details: undefined,
    length: 0,
    expiry: undefined,
    includeLowerLetters: false,
    includeUpperLetters: false,
    includeDigits: false,
    includeSymbols: false,
  };
}

export const SecretGenerator = {
  encode(message: SecretGenerator, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.generatorType !== 0) {
      writer.uint32(8).int32(message.generatorType);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.length !== 0) {
      writer.uint32(24).uint32(message.length);
    }
    if (message.expiry !== undefined) {
      Duration.encode(message.expiry, writer.uint32(34).fork()).ldelim();
    }
    if (message.includeLowerLetters === true) {
      writer.uint32(40).bool(message.includeLowerLetters);
    }
    if (message.includeUpperLetters === true) {
      writer.uint32(48).bool(message.includeUpperLetters);
    }
    if (message.includeDigits === true) {
      writer.uint32(56).bool(message.includeDigits);
    }
    if (message.includeSymbols === true) {
      writer.uint32(64).bool(message.includeSymbols);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecretGenerator {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecretGenerator();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.generatorType = reader.int32() as any;
          break;
        case 2:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 3:
          message.length = reader.uint32();
          break;
        case 4:
          message.expiry = Duration.decode(reader, reader.uint32());
          break;
        case 5:
          message.includeLowerLetters = reader.bool();
          break;
        case 6:
          message.includeUpperLetters = reader.bool();
          break;
        case 7:
          message.includeDigits = reader.bool();
          break;
        case 8:
          message.includeSymbols = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SecretGenerator {
    return {
      generatorType: isSet(object.generatorType) ? secretGeneratorTypeFromJSON(object.generatorType) : 0,
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      length: isSet(object.length) ? Number(object.length) : 0,
      expiry: isSet(object.expiry) ? Duration.fromJSON(object.expiry) : undefined,
      includeLowerLetters: isSet(object.includeLowerLetters) ? Boolean(object.includeLowerLetters) : false,
      includeUpperLetters: isSet(object.includeUpperLetters) ? Boolean(object.includeUpperLetters) : false,
      includeDigits: isSet(object.includeDigits) ? Boolean(object.includeDigits) : false,
      includeSymbols: isSet(object.includeSymbols) ? Boolean(object.includeSymbols) : false,
    };
  },

  toJSON(message: SecretGenerator): unknown {
    const obj: any = {};
    message.generatorType !== undefined && (obj.generatorType = secretGeneratorTypeToJSON(message.generatorType));
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.length !== undefined && (obj.length = Math.round(message.length));
    message.expiry !== undefined && (obj.expiry = message.expiry ? Duration.toJSON(message.expiry) : undefined);
    message.includeLowerLetters !== undefined && (obj.includeLowerLetters = message.includeLowerLetters);
    message.includeUpperLetters !== undefined && (obj.includeUpperLetters = message.includeUpperLetters);
    message.includeDigits !== undefined && (obj.includeDigits = message.includeDigits);
    message.includeSymbols !== undefined && (obj.includeSymbols = message.includeSymbols);
    return obj;
  },

  create(base?: DeepPartial<SecretGenerator>): SecretGenerator {
    return SecretGenerator.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecretGenerator>): SecretGenerator {
    const message = createBaseSecretGenerator();
    message.generatorType = object.generatorType ?? 0;
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.length = object.length ?? 0;
    message.expiry = (object.expiry !== undefined && object.expiry !== null)
      ? Duration.fromPartial(object.expiry)
      : undefined;
    message.includeLowerLetters = object.includeLowerLetters ?? false;
    message.includeUpperLetters = object.includeUpperLetters ?? false;
    message.includeDigits = object.includeDigits ?? false;
    message.includeSymbols = object.includeSymbols ?? false;
    return message;
  },
};

function createBaseSecretGeneratorQuery(): SecretGeneratorQuery {
  return { typeQuery: undefined };
}

export const SecretGeneratorQuery = {
  encode(message: SecretGeneratorQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.typeQuery !== undefined) {
      SecretGeneratorTypeQuery.encode(message.typeQuery, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecretGeneratorQuery {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecretGeneratorQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.typeQuery = SecretGeneratorTypeQuery.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SecretGeneratorQuery {
    return { typeQuery: isSet(object.typeQuery) ? SecretGeneratorTypeQuery.fromJSON(object.typeQuery) : undefined };
  },

  toJSON(message: SecretGeneratorQuery): unknown {
    const obj: any = {};
    message.typeQuery !== undefined &&
      (obj.typeQuery = message.typeQuery ? SecretGeneratorTypeQuery.toJSON(message.typeQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SecretGeneratorQuery>): SecretGeneratorQuery {
    return SecretGeneratorQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecretGeneratorQuery>): SecretGeneratorQuery {
    const message = createBaseSecretGeneratorQuery();
    message.typeQuery = (object.typeQuery !== undefined && object.typeQuery !== null)
      ? SecretGeneratorTypeQuery.fromPartial(object.typeQuery)
      : undefined;
    return message;
  },
};

function createBaseSecretGeneratorTypeQuery(): SecretGeneratorTypeQuery {
  return { generatorType: 0 };
}

export const SecretGeneratorTypeQuery = {
  encode(message: SecretGeneratorTypeQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.generatorType !== 0) {
      writer.uint32(8).int32(message.generatorType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecretGeneratorTypeQuery {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecretGeneratorTypeQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.generatorType = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SecretGeneratorTypeQuery {
    return { generatorType: isSet(object.generatorType) ? secretGeneratorTypeFromJSON(object.generatorType) : 0 };
  },

  toJSON(message: SecretGeneratorTypeQuery): unknown {
    const obj: any = {};
    message.generatorType !== undefined && (obj.generatorType = secretGeneratorTypeToJSON(message.generatorType));
    return obj;
  },

  create(base?: DeepPartial<SecretGeneratorTypeQuery>): SecretGeneratorTypeQuery {
    return SecretGeneratorTypeQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecretGeneratorTypeQuery>): SecretGeneratorTypeQuery {
    const message = createBaseSecretGeneratorTypeQuery();
    message.generatorType = object.generatorType ?? 0;
    return message;
  },
};

function createBaseSMTPConfig(): SMTPConfig {
  return { details: undefined, senderAddress: "", senderName: "", tls: false, host: "", user: "" };
}

export const SMTPConfig = {
  encode(message: SMTPConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.senderAddress !== "") {
      writer.uint32(18).string(message.senderAddress);
    }
    if (message.senderName !== "") {
      writer.uint32(26).string(message.senderName);
    }
    if (message.tls === true) {
      writer.uint32(32).bool(message.tls);
    }
    if (message.host !== "") {
      writer.uint32(42).string(message.host);
    }
    if (message.user !== "") {
      writer.uint32(50).string(message.user);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SMTPConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSMTPConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.senderAddress = reader.string();
          break;
        case 3:
          message.senderName = reader.string();
          break;
        case 4:
          message.tls = reader.bool();
          break;
        case 5:
          message.host = reader.string();
          break;
        case 6:
          message.user = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SMTPConfig {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      senderAddress: isSet(object.senderAddress) ? String(object.senderAddress) : "",
      senderName: isSet(object.senderName) ? String(object.senderName) : "",
      tls: isSet(object.tls) ? Boolean(object.tls) : false,
      host: isSet(object.host) ? String(object.host) : "",
      user: isSet(object.user) ? String(object.user) : "",
    };
  },

  toJSON(message: SMTPConfig): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.senderAddress !== undefined && (obj.senderAddress = message.senderAddress);
    message.senderName !== undefined && (obj.senderName = message.senderName);
    message.tls !== undefined && (obj.tls = message.tls);
    message.host !== undefined && (obj.host = message.host);
    message.user !== undefined && (obj.user = message.user);
    return obj;
  },

  create(base?: DeepPartial<SMTPConfig>): SMTPConfig {
    return SMTPConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SMTPConfig>): SMTPConfig {
    const message = createBaseSMTPConfig();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.senderAddress = object.senderAddress ?? "";
    message.senderName = object.senderName ?? "";
    message.tls = object.tls ?? false;
    message.host = object.host ?? "";
    message.user = object.user ?? "";
    return message;
  },
};

function createBaseSMSProvider(): SMSProvider {
  return { details: undefined, id: "", state: 0, twilio: undefined };
}

export const SMSProvider = {
  encode(message: SMSProvider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    if (message.state !== 0) {
      writer.uint32(24).int32(message.state);
    }
    if (message.twilio !== undefined) {
      TwilioConfig.encode(message.twilio, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SMSProvider {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSMSProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.id = reader.string();
          break;
        case 3:
          message.state = reader.int32() as any;
          break;
        case 4:
          message.twilio = TwilioConfig.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SMSProvider {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      id: isSet(object.id) ? String(object.id) : "",
      state: isSet(object.state) ? sMSProviderConfigStateFromJSON(object.state) : 0,
      twilio: isSet(object.twilio) ? TwilioConfig.fromJSON(object.twilio) : undefined,
    };
  },

  toJSON(message: SMSProvider): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.id !== undefined && (obj.id = message.id);
    message.state !== undefined && (obj.state = sMSProviderConfigStateToJSON(message.state));
    message.twilio !== undefined && (obj.twilio = message.twilio ? TwilioConfig.toJSON(message.twilio) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SMSProvider>): SMSProvider {
    return SMSProvider.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SMSProvider>): SMSProvider {
    const message = createBaseSMSProvider();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.id = object.id ?? "";
    message.state = object.state ?? 0;
    message.twilio = (object.twilio !== undefined && object.twilio !== null)
      ? TwilioConfig.fromPartial(object.twilio)
      : undefined;
    return message;
  },
};

function createBaseTwilioConfig(): TwilioConfig {
  return { sid: "", senderNumber: "" };
}

export const TwilioConfig = {
  encode(message: TwilioConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sid !== "") {
      writer.uint32(10).string(message.sid);
    }
    if (message.senderNumber !== "") {
      writer.uint32(18).string(message.senderNumber);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TwilioConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTwilioConfig();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sid = reader.string();
          break;
        case 2:
          message.senderNumber = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TwilioConfig {
    return {
      sid: isSet(object.sid) ? String(object.sid) : "",
      senderNumber: isSet(object.senderNumber) ? String(object.senderNumber) : "",
    };
  },

  toJSON(message: TwilioConfig): unknown {
    const obj: any = {};
    message.sid !== undefined && (obj.sid = message.sid);
    message.senderNumber !== undefined && (obj.senderNumber = message.senderNumber);
    return obj;
  },

  create(base?: DeepPartial<TwilioConfig>): TwilioConfig {
    return TwilioConfig.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TwilioConfig>): TwilioConfig {
    const message = createBaseTwilioConfig();
    message.sid = object.sid ?? "";
    message.senderNumber = object.senderNumber ?? "";
    return message;
  },
};

function createBaseDebugNotificationProvider(): DebugNotificationProvider {
  return { details: undefined, compact: false };
}

export const DebugNotificationProvider = {
  encode(message: DebugNotificationProvider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.compact === true) {
      writer.uint32(16).bool(message.compact);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DebugNotificationProvider {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDebugNotificationProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.compact = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DebugNotificationProvider {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      compact: isSet(object.compact) ? Boolean(object.compact) : false,
    };
  },

  toJSON(message: DebugNotificationProvider): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.compact !== undefined && (obj.compact = message.compact);
    return obj;
  },

  create(base?: DeepPartial<DebugNotificationProvider>): DebugNotificationProvider {
    return DebugNotificationProvider.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DebugNotificationProvider>): DebugNotificationProvider {
    const message = createBaseDebugNotificationProvider();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.compact = object.compact ?? false;
    return message;
  },
};

function createBaseOIDCSettings(): OIDCSettings {
  return {
    details: undefined,
    accessTokenLifetime: undefined,
    idTokenLifetime: undefined,
    refreshTokenIdleExpiration: undefined,
    refreshTokenExpiration: undefined,
  };
}

export const OIDCSettings = {
  encode(message: OIDCSettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.accessTokenLifetime !== undefined) {
      Duration.encode(message.accessTokenLifetime, writer.uint32(18).fork()).ldelim();
    }
    if (message.idTokenLifetime !== undefined) {
      Duration.encode(message.idTokenLifetime, writer.uint32(26).fork()).ldelim();
    }
    if (message.refreshTokenIdleExpiration !== undefined) {
      Duration.encode(message.refreshTokenIdleExpiration, writer.uint32(34).fork()).ldelim();
    }
    if (message.refreshTokenExpiration !== undefined) {
      Duration.encode(message.refreshTokenExpiration, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OIDCSettings {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOIDCSettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.accessTokenLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 3:
          message.idTokenLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 4:
          message.refreshTokenIdleExpiration = Duration.decode(reader, reader.uint32());
          break;
        case 5:
          message.refreshTokenExpiration = Duration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): OIDCSettings {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      accessTokenLifetime: isSet(object.accessTokenLifetime)
        ? Duration.fromJSON(object.accessTokenLifetime)
        : undefined,
      idTokenLifetime: isSet(object.idTokenLifetime) ? Duration.fromJSON(object.idTokenLifetime) : undefined,
      refreshTokenIdleExpiration: isSet(object.refreshTokenIdleExpiration)
        ? Duration.fromJSON(object.refreshTokenIdleExpiration)
        : undefined,
      refreshTokenExpiration: isSet(object.refreshTokenExpiration)
        ? Duration.fromJSON(object.refreshTokenExpiration)
        : undefined,
    };
  },

  toJSON(message: OIDCSettings): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.accessTokenLifetime !== undefined &&
      (obj.accessTokenLifetime = message.accessTokenLifetime
        ? Duration.toJSON(message.accessTokenLifetime)
        : undefined);
    message.idTokenLifetime !== undefined &&
      (obj.idTokenLifetime = message.idTokenLifetime ? Duration.toJSON(message.idTokenLifetime) : undefined);
    message.refreshTokenIdleExpiration !== undefined &&
      (obj.refreshTokenIdleExpiration = message.refreshTokenIdleExpiration
        ? Duration.toJSON(message.refreshTokenIdleExpiration)
        : undefined);
    message.refreshTokenExpiration !== undefined && (obj.refreshTokenExpiration = message.refreshTokenExpiration
      ? Duration.toJSON(message.refreshTokenExpiration)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<OIDCSettings>): OIDCSettings {
    return OIDCSettings.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OIDCSettings>): OIDCSettings {
    const message = createBaseOIDCSettings();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.accessTokenLifetime = (object.accessTokenLifetime !== undefined && object.accessTokenLifetime !== null)
      ? Duration.fromPartial(object.accessTokenLifetime)
      : undefined;
    message.idTokenLifetime = (object.idTokenLifetime !== undefined && object.idTokenLifetime !== null)
      ? Duration.fromPartial(object.idTokenLifetime)
      : undefined;
    message.refreshTokenIdleExpiration =
      (object.refreshTokenIdleExpiration !== undefined && object.refreshTokenIdleExpiration !== null)
        ? Duration.fromPartial(object.refreshTokenIdleExpiration)
        : undefined;
    message.refreshTokenExpiration =
      (object.refreshTokenExpiration !== undefined && object.refreshTokenExpiration !== null)
        ? Duration.fromPartial(object.refreshTokenExpiration)
        : undefined;
    return message;
  },
};

function createBaseSecurityPolicy(): SecurityPolicy {
  return { details: undefined, enableIframeEmbedding: false, allowedOrigins: [] };
}

export const SecurityPolicy = {
  encode(message: SecurityPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.enableIframeEmbedding === true) {
      writer.uint32(16).bool(message.enableIframeEmbedding);
    }
    for (const v of message.allowedOrigins) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecurityPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecurityPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.enableIframeEmbedding = reader.bool();
          break;
        case 3:
          message.allowedOrigins.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SecurityPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      enableIframeEmbedding: isSet(object.enableIframeEmbedding) ? Boolean(object.enableIframeEmbedding) : false,
      allowedOrigins: Array.isArray(object?.allowedOrigins) ? object.allowedOrigins.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: SecurityPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.enableIframeEmbedding !== undefined && (obj.enableIframeEmbedding = message.enableIframeEmbedding);
    if (message.allowedOrigins) {
      obj.allowedOrigins = message.allowedOrigins.map((e) => e);
    } else {
      obj.allowedOrigins = [];
    }
    return obj;
  },

  create(base?: DeepPartial<SecurityPolicy>): SecurityPolicy {
    return SecurityPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecurityPolicy>): SecurityPolicy {
    const message = createBaseSecurityPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.enableIframeEmbedding = object.enableIframeEmbedding ?? false;
    message.allowedOrigins = object.allowedOrigins?.map((e) => e) || [];
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

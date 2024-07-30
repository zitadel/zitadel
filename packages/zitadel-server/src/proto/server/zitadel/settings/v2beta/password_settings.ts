/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { ResourceOwnerType, resourceOwnerTypeFromJSON, resourceOwnerTypeToJSON } from "./settings";

export const protobufPackage = "zitadel.settings.v2beta";

export interface PasswordComplexitySettings {
  minLength: number;
  requiresUppercase: boolean;
  requiresLowercase: boolean;
  requiresNumber: boolean;
  requiresSymbol: boolean;
  /** resource_owner_type returns if the settings is managed on the organization or on the instance */
  resourceOwnerType: ResourceOwnerType;
}

export interface PasswordExpirySettings {
  /** Amount of days after which a password will expire. The user will be forced to change the password on the following authentication. */
  maxAgeDays: number;
  /** Amount of days after which the user should be notified of the upcoming expiry. ZITADEL will not notify the user. */
  expireWarnDays: number;
  /** resource_owner_type returns if the settings is managed on the organization or on the instance */
  resourceOwnerType: ResourceOwnerType;
}

function createBasePasswordComplexitySettings(): PasswordComplexitySettings {
  return {
    minLength: 0,
    requiresUppercase: false,
    requiresLowercase: false,
    requiresNumber: false,
    requiresSymbol: false,
    resourceOwnerType: 0,
  };
}

export const PasswordComplexitySettings = {
  encode(message: PasswordComplexitySettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.minLength !== 0) {
      writer.uint32(8).uint64(message.minLength);
    }
    if (message.requiresUppercase === true) {
      writer.uint32(16).bool(message.requiresUppercase);
    }
    if (message.requiresLowercase === true) {
      writer.uint32(24).bool(message.requiresLowercase);
    }
    if (message.requiresNumber === true) {
      writer.uint32(32).bool(message.requiresNumber);
    }
    if (message.requiresSymbol === true) {
      writer.uint32(40).bool(message.requiresSymbol);
    }
    if (message.resourceOwnerType !== 0) {
      writer.uint32(48).int32(message.resourceOwnerType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordComplexitySettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordComplexitySettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.minLength = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.requiresUppercase = reader.bool();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.requiresLowercase = reader.bool();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.requiresNumber = reader.bool();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.requiresSymbol = reader.bool();
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.resourceOwnerType = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PasswordComplexitySettings {
    return {
      minLength: isSet(object.minLength) ? Number(object.minLength) : 0,
      requiresUppercase: isSet(object.requiresUppercase) ? Boolean(object.requiresUppercase) : false,
      requiresLowercase: isSet(object.requiresLowercase) ? Boolean(object.requiresLowercase) : false,
      requiresNumber: isSet(object.requiresNumber) ? Boolean(object.requiresNumber) : false,
      requiresSymbol: isSet(object.requiresSymbol) ? Boolean(object.requiresSymbol) : false,
      resourceOwnerType: isSet(object.resourceOwnerType) ? resourceOwnerTypeFromJSON(object.resourceOwnerType) : 0,
    };
  },

  toJSON(message: PasswordComplexitySettings): unknown {
    const obj: any = {};
    message.minLength !== undefined && (obj.minLength = Math.round(message.minLength));
    message.requiresUppercase !== undefined && (obj.requiresUppercase = message.requiresUppercase);
    message.requiresLowercase !== undefined && (obj.requiresLowercase = message.requiresLowercase);
    message.requiresNumber !== undefined && (obj.requiresNumber = message.requiresNumber);
    message.requiresSymbol !== undefined && (obj.requiresSymbol = message.requiresSymbol);
    message.resourceOwnerType !== undefined &&
      (obj.resourceOwnerType = resourceOwnerTypeToJSON(message.resourceOwnerType));
    return obj;
  },

  create(base?: DeepPartial<PasswordComplexitySettings>): PasswordComplexitySettings {
    return PasswordComplexitySettings.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordComplexitySettings>): PasswordComplexitySettings {
    const message = createBasePasswordComplexitySettings();
    message.minLength = object.minLength ?? 0;
    message.requiresUppercase = object.requiresUppercase ?? false;
    message.requiresLowercase = object.requiresLowercase ?? false;
    message.requiresNumber = object.requiresNumber ?? false;
    message.requiresSymbol = object.requiresSymbol ?? false;
    message.resourceOwnerType = object.resourceOwnerType ?? 0;
    return message;
  },
};

function createBasePasswordExpirySettings(): PasswordExpirySettings {
  return { maxAgeDays: 0, expireWarnDays: 0, resourceOwnerType: 0 };
}

export const PasswordExpirySettings = {
  encode(message: PasswordExpirySettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.maxAgeDays !== 0) {
      writer.uint32(8).uint64(message.maxAgeDays);
    }
    if (message.expireWarnDays !== 0) {
      writer.uint32(16).uint64(message.expireWarnDays);
    }
    if (message.resourceOwnerType !== 0) {
      writer.uint32(24).int32(message.resourceOwnerType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordExpirySettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordExpirySettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.maxAgeDays = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.expireWarnDays = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.resourceOwnerType = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PasswordExpirySettings {
    return {
      maxAgeDays: isSet(object.maxAgeDays) ? Number(object.maxAgeDays) : 0,
      expireWarnDays: isSet(object.expireWarnDays) ? Number(object.expireWarnDays) : 0,
      resourceOwnerType: isSet(object.resourceOwnerType) ? resourceOwnerTypeFromJSON(object.resourceOwnerType) : 0,
    };
  },

  toJSON(message: PasswordExpirySettings): unknown {
    const obj: any = {};
    message.maxAgeDays !== undefined && (obj.maxAgeDays = Math.round(message.maxAgeDays));
    message.expireWarnDays !== undefined && (obj.expireWarnDays = Math.round(message.expireWarnDays));
    message.resourceOwnerType !== undefined &&
      (obj.resourceOwnerType = resourceOwnerTypeToJSON(message.resourceOwnerType));
    return obj;
  },

  create(base?: DeepPartial<PasswordExpirySettings>): PasswordExpirySettings {
    return PasswordExpirySettings.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordExpirySettings>): PasswordExpirySettings {
    const message = createBasePasswordExpirySettings();
    message.maxAgeDays = object.maxAgeDays ?? 0;
    message.expireWarnDays = object.expireWarnDays ?? 0;
    message.resourceOwnerType = object.resourceOwnerType ?? 0;
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

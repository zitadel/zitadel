/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { ResourceOwnerType, resourceOwnerTypeFromJSON, resourceOwnerTypeToJSON } from "./settings";

export const protobufPackage = "zitadel.settings.v2beta";

export interface LockoutSettings {
  maxPasswordAttempts: number;
  /** resource_owner_type returns if the settings is managed on the organization or on the instance */
  resourceOwnerType: ResourceOwnerType;
  maxOtpAttempts: number;
}

function createBaseLockoutSettings(): LockoutSettings {
  return { maxPasswordAttempts: 0, resourceOwnerType: 0, maxOtpAttempts: 0 };
}

export const LockoutSettings = {
  encode(message: LockoutSettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.maxPasswordAttempts !== 0) {
      writer.uint32(8).uint64(message.maxPasswordAttempts);
    }
    if (message.resourceOwnerType !== 0) {
      writer.uint32(16).int32(message.resourceOwnerType);
    }
    if (message.maxOtpAttempts !== 0) {
      writer.uint32(24).uint64(message.maxOtpAttempts);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LockoutSettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLockoutSettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.maxPasswordAttempts = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.resourceOwnerType = reader.int32() as any;
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.maxOtpAttempts = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LockoutSettings {
    return {
      maxPasswordAttempts: isSet(object.maxPasswordAttempts) ? Number(object.maxPasswordAttempts) : 0,
      resourceOwnerType: isSet(object.resourceOwnerType) ? resourceOwnerTypeFromJSON(object.resourceOwnerType) : 0,
      maxOtpAttempts: isSet(object.maxOtpAttempts) ? Number(object.maxOtpAttempts) : 0,
    };
  },

  toJSON(message: LockoutSettings): unknown {
    const obj: any = {};
    message.maxPasswordAttempts !== undefined && (obj.maxPasswordAttempts = Math.round(message.maxPasswordAttempts));
    message.resourceOwnerType !== undefined &&
      (obj.resourceOwnerType = resourceOwnerTypeToJSON(message.resourceOwnerType));
    message.maxOtpAttempts !== undefined && (obj.maxOtpAttempts = Math.round(message.maxOtpAttempts));
    return obj;
  },

  create(base?: DeepPartial<LockoutSettings>): LockoutSettings {
    return LockoutSettings.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LockoutSettings>): LockoutSettings {
    const message = createBaseLockoutSettings();
    message.maxPasswordAttempts = object.maxPasswordAttempts ?? 0;
    message.resourceOwnerType = object.resourceOwnerType ?? 0;
    message.maxOtpAttempts = object.maxOtpAttempts ?? 0;
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

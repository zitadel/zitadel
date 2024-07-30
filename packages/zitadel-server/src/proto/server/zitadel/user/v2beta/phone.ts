/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.user.v2beta";

export interface SetHumanPhone {
  phone: string;
  sendCode?: SendPhoneVerificationCode | undefined;
  returnCode?: ReturnPhoneVerificationCode | undefined;
  isVerified?: boolean | undefined;
}

export interface HumanPhone {
  phone: string;
  isVerified: boolean;
}

export interface SendPhoneVerificationCode {
}

export interface ReturnPhoneVerificationCode {
}

function createBaseSetHumanPhone(): SetHumanPhone {
  return { phone: "", sendCode: undefined, returnCode: undefined, isVerified: undefined };
}

export const SetHumanPhone = {
  encode(message: SetHumanPhone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.phone !== "") {
      writer.uint32(10).string(message.phone);
    }
    if (message.sendCode !== undefined) {
      SendPhoneVerificationCode.encode(message.sendCode, writer.uint32(18).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      ReturnPhoneVerificationCode.encode(message.returnCode, writer.uint32(26).fork()).ldelim();
    }
    if (message.isVerified !== undefined) {
      writer.uint32(32).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetHumanPhone {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetHumanPhone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.phone = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sendCode = SendPhoneVerificationCode.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.returnCode = ReturnPhoneVerificationCode.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.isVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetHumanPhone {
    return {
      phone: isSet(object.phone) ? String(object.phone) : "",
      sendCode: isSet(object.sendCode) ? SendPhoneVerificationCode.fromJSON(object.sendCode) : undefined,
      returnCode: isSet(object.returnCode) ? ReturnPhoneVerificationCode.fromJSON(object.returnCode) : undefined,
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : undefined,
    };
  },

  toJSON(message: SetHumanPhone): unknown {
    const obj: any = {};
    message.phone !== undefined && (obj.phone = message.phone);
    message.sendCode !== undefined &&
      (obj.sendCode = message.sendCode ? SendPhoneVerificationCode.toJSON(message.sendCode) : undefined);
    message.returnCode !== undefined &&
      (obj.returnCode = message.returnCode ? ReturnPhoneVerificationCode.toJSON(message.returnCode) : undefined);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<SetHumanPhone>): SetHumanPhone {
    return SetHumanPhone.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetHumanPhone>): SetHumanPhone {
    const message = createBaseSetHumanPhone();
    message.phone = object.phone ?? "";
    message.sendCode = (object.sendCode !== undefined && object.sendCode !== null)
      ? SendPhoneVerificationCode.fromPartial(object.sendCode)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? ReturnPhoneVerificationCode.fromPartial(object.returnCode)
      : undefined;
    message.isVerified = object.isVerified ?? undefined;
    return message;
  },
};

function createBaseHumanPhone(): HumanPhone {
  return { phone: "", isVerified: false };
}

export const HumanPhone = {
  encode(message: HumanPhone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.phone !== "") {
      writer.uint32(10).string(message.phone);
    }
    if (message.isVerified === true) {
      writer.uint32(16).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HumanPhone {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHumanPhone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.phone = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.isVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HumanPhone {
    return {
      phone: isSet(object.phone) ? String(object.phone) : "",
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : false,
    };
  },

  toJSON(message: HumanPhone): unknown {
    const obj: any = {};
    message.phone !== undefined && (obj.phone = message.phone);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<HumanPhone>): HumanPhone {
    return HumanPhone.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<HumanPhone>): HumanPhone {
    const message = createBaseHumanPhone();
    message.phone = object.phone ?? "";
    message.isVerified = object.isVerified ?? false;
    return message;
  },
};

function createBaseSendPhoneVerificationCode(): SendPhoneVerificationCode {
  return {};
}

export const SendPhoneVerificationCode = {
  encode(_: SendPhoneVerificationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendPhoneVerificationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendPhoneVerificationCode();
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

  fromJSON(_: any): SendPhoneVerificationCode {
    return {};
  },

  toJSON(_: SendPhoneVerificationCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<SendPhoneVerificationCode>): SendPhoneVerificationCode {
    return SendPhoneVerificationCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<SendPhoneVerificationCode>): SendPhoneVerificationCode {
    const message = createBaseSendPhoneVerificationCode();
    return message;
  },
};

function createBaseReturnPhoneVerificationCode(): ReturnPhoneVerificationCode {
  return {};
}

export const ReturnPhoneVerificationCode = {
  encode(_: ReturnPhoneVerificationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReturnPhoneVerificationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReturnPhoneVerificationCode();
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

  fromJSON(_: any): ReturnPhoneVerificationCode {
    return {};
  },

  toJSON(_: ReturnPhoneVerificationCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ReturnPhoneVerificationCode>): ReturnPhoneVerificationCode {
    return ReturnPhoneVerificationCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ReturnPhoneVerificationCode>): ReturnPhoneVerificationCode {
    const message = createBaseReturnPhoneVerificationCode();
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

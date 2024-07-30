/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.user.v2beta";

export interface SetHumanEmail {
  email: string;
  sendCode?: SendEmailVerificationCode | undefined;
  returnCode?: ReturnEmailVerificationCode | undefined;
  isVerified?: boolean | undefined;
}

export interface HumanEmail {
  email: string;
  isVerified: boolean;
}

export interface SendEmailVerificationCode {
  urlTemplate?: string | undefined;
}

export interface ReturnEmailVerificationCode {
}

function createBaseSetHumanEmail(): SetHumanEmail {
  return { email: "", sendCode: undefined, returnCode: undefined, isVerified: undefined };
}

export const SetHumanEmail = {
  encode(message: SetHumanEmail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.email !== "") {
      writer.uint32(10).string(message.email);
    }
    if (message.sendCode !== undefined) {
      SendEmailVerificationCode.encode(message.sendCode, writer.uint32(18).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      ReturnEmailVerificationCode.encode(message.returnCode, writer.uint32(26).fork()).ldelim();
    }
    if (message.isVerified !== undefined) {
      writer.uint32(32).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetHumanEmail {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetHumanEmail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.email = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sendCode = SendEmailVerificationCode.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.returnCode = ReturnEmailVerificationCode.decode(reader, reader.uint32());
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

  fromJSON(object: any): SetHumanEmail {
    return {
      email: isSet(object.email) ? String(object.email) : "",
      sendCode: isSet(object.sendCode) ? SendEmailVerificationCode.fromJSON(object.sendCode) : undefined,
      returnCode: isSet(object.returnCode) ? ReturnEmailVerificationCode.fromJSON(object.returnCode) : undefined,
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : undefined,
    };
  },

  toJSON(message: SetHumanEmail): unknown {
    const obj: any = {};
    message.email !== undefined && (obj.email = message.email);
    message.sendCode !== undefined &&
      (obj.sendCode = message.sendCode ? SendEmailVerificationCode.toJSON(message.sendCode) : undefined);
    message.returnCode !== undefined &&
      (obj.returnCode = message.returnCode ? ReturnEmailVerificationCode.toJSON(message.returnCode) : undefined);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<SetHumanEmail>): SetHumanEmail {
    return SetHumanEmail.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetHumanEmail>): SetHumanEmail {
    const message = createBaseSetHumanEmail();
    message.email = object.email ?? "";
    message.sendCode = (object.sendCode !== undefined && object.sendCode !== null)
      ? SendEmailVerificationCode.fromPartial(object.sendCode)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? ReturnEmailVerificationCode.fromPartial(object.returnCode)
      : undefined;
    message.isVerified = object.isVerified ?? undefined;
    return message;
  },
};

function createBaseHumanEmail(): HumanEmail {
  return { email: "", isVerified: false };
}

export const HumanEmail = {
  encode(message: HumanEmail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.email !== "") {
      writer.uint32(10).string(message.email);
    }
    if (message.isVerified === true) {
      writer.uint32(16).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HumanEmail {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHumanEmail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.email = reader.string();
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

  fromJSON(object: any): HumanEmail {
    return {
      email: isSet(object.email) ? String(object.email) : "",
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : false,
    };
  },

  toJSON(message: HumanEmail): unknown {
    const obj: any = {};
    message.email !== undefined && (obj.email = message.email);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<HumanEmail>): HumanEmail {
    return HumanEmail.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<HumanEmail>): HumanEmail {
    const message = createBaseHumanEmail();
    message.email = object.email ?? "";
    message.isVerified = object.isVerified ?? false;
    return message;
  },
};

function createBaseSendEmailVerificationCode(): SendEmailVerificationCode {
  return { urlTemplate: undefined };
}

export const SendEmailVerificationCode = {
  encode(message: SendEmailVerificationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.urlTemplate !== undefined) {
      writer.uint32(10).string(message.urlTemplate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendEmailVerificationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendEmailVerificationCode();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.urlTemplate = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SendEmailVerificationCode {
    return { urlTemplate: isSet(object.urlTemplate) ? String(object.urlTemplate) : undefined };
  },

  toJSON(message: SendEmailVerificationCode): unknown {
    const obj: any = {};
    message.urlTemplate !== undefined && (obj.urlTemplate = message.urlTemplate);
    return obj;
  },

  create(base?: DeepPartial<SendEmailVerificationCode>): SendEmailVerificationCode {
    return SendEmailVerificationCode.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SendEmailVerificationCode>): SendEmailVerificationCode {
    const message = createBaseSendEmailVerificationCode();
    message.urlTemplate = object.urlTemplate ?? undefined;
    return message;
  },
};

function createBaseReturnEmailVerificationCode(): ReturnEmailVerificationCode {
  return {};
}

export const ReturnEmailVerificationCode = {
  encode(_: ReturnEmailVerificationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReturnEmailVerificationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReturnEmailVerificationCode();
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

  fromJSON(_: any): ReturnEmailVerificationCode {
    return {};
  },

  toJSON(_: ReturnEmailVerificationCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ReturnEmailVerificationCode>): ReturnEmailVerificationCode {
    return ReturnEmailVerificationCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ReturnEmailVerificationCode>): ReturnEmailVerificationCode {
    const message = createBaseReturnEmailVerificationCode();
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

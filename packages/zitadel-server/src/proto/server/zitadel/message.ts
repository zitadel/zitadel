/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.v1";

export interface ErrorDetail {
  id: string;
  message: string;
}

export interface LocalizedMessage {
  key: string;
  localizedMessage: string;
}

function createBaseErrorDetail(): ErrorDetail {
  return { id: "", message: "" };
}

export const ErrorDetail = {
  encode(message: ErrorDetail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ErrorDetail {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseErrorDetail();
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

          message.message = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ErrorDetail {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      message: isSet(object.message) ? String(object.message) : "",
    };
  },

  toJSON(message: ErrorDetail): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },

  create(base?: DeepPartial<ErrorDetail>): ErrorDetail {
    return ErrorDetail.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ErrorDetail>): ErrorDetail {
    const message = createBaseErrorDetail();
    message.id = object.id ?? "";
    message.message = object.message ?? "";
    return message;
  },
};

function createBaseLocalizedMessage(): LocalizedMessage {
  return { key: "", localizedMessage: "" };
}

export const LocalizedMessage = {
  encode(message: LocalizedMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.localizedMessage !== "") {
      writer.uint32(18).string(message.localizedMessage);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LocalizedMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLocalizedMessage();
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

          message.localizedMessage = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LocalizedMessage {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      localizedMessage: isSet(object.localizedMessage) ? String(object.localizedMessage) : "",
    };
  },

  toJSON(message: LocalizedMessage): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.localizedMessage !== undefined && (obj.localizedMessage = message.localizedMessage);
    return obj;
  },

  create(base?: DeepPartial<LocalizedMessage>): LocalizedMessage {
    return LocalizedMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LocalizedMessage>): LocalizedMessage {
    const message = createBaseLocalizedMessage();
    message.key = object.key ?? "";
    message.localizedMessage = object.localizedMessage ?? "";
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

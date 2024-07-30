/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.v1";

export interface AuthOption {
  permission: string;
  checkFieldName: string;
}

function createBaseAuthOption(): AuthOption {
  return { permission: "", checkFieldName: "" };
}

export const AuthOption = {
  encode(message: AuthOption, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.permission !== "") {
      writer.uint32(10).string(message.permission);
    }
    if (message.checkFieldName !== "") {
      writer.uint32(18).string(message.checkFieldName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthOption {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthOption();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.permission = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.checkFieldName = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AuthOption {
    return {
      permission: isSet(object.permission) ? String(object.permission) : "",
      checkFieldName: isSet(object.checkFieldName) ? String(object.checkFieldName) : "",
    };
  },

  toJSON(message: AuthOption): unknown {
    const obj: any = {};
    message.permission !== undefined && (obj.permission = message.permission);
    message.checkFieldName !== undefined && (obj.checkFieldName = message.checkFieldName);
    return obj;
  },

  create(base?: DeepPartial<AuthOption>): AuthOption {
    return AuthOption.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AuthOption>): AuthOption {
    const message = createBaseAuthOption();
    message.permission = object.permission ?? "";
    message.checkFieldName = object.checkFieldName ?? "";
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

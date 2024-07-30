/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.protoc_gen_zitadel.v2";

export interface Options {
  authOption: AuthOption | undefined;
  httpResponse: CustomHTTPResponse | undefined;
}

export interface AuthOption {
  permission: string;
  orgField: string;
}

export interface CustomHTTPResponse {
  successCode: number;
}

function createBaseOptions(): Options {
  return { authOption: undefined, httpResponse: undefined };
}

export const Options = {
  encode(message: Options, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authOption !== undefined) {
      AuthOption.encode(message.authOption, writer.uint32(10).fork()).ldelim();
    }
    if (message.httpResponse !== undefined) {
      CustomHTTPResponse.encode(message.httpResponse, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Options {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOptions();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.authOption = AuthOption.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.httpResponse = CustomHTTPResponse.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Options {
    return {
      authOption: isSet(object.authOption) ? AuthOption.fromJSON(object.authOption) : undefined,
      httpResponse: isSet(object.httpResponse) ? CustomHTTPResponse.fromJSON(object.httpResponse) : undefined,
    };
  },

  toJSON(message: Options): unknown {
    const obj: any = {};
    message.authOption !== undefined &&
      (obj.authOption = message.authOption ? AuthOption.toJSON(message.authOption) : undefined);
    message.httpResponse !== undefined &&
      (obj.httpResponse = message.httpResponse ? CustomHTTPResponse.toJSON(message.httpResponse) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Options>): Options {
    return Options.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Options>): Options {
    const message = createBaseOptions();
    message.authOption = (object.authOption !== undefined && object.authOption !== null)
      ? AuthOption.fromPartial(object.authOption)
      : undefined;
    message.httpResponse = (object.httpResponse !== undefined && object.httpResponse !== null)
      ? CustomHTTPResponse.fromPartial(object.httpResponse)
      : undefined;
    return message;
  },
};

function createBaseAuthOption(): AuthOption {
  return { permission: "", orgField: "" };
}

export const AuthOption = {
  encode(message: AuthOption, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.permission !== "") {
      writer.uint32(10).string(message.permission);
    }
    if (message.orgField !== "") {
      writer.uint32(26).string(message.orgField);
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
        case 3:
          if (tag != 26) {
            break;
          }

          message.orgField = reader.string();
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
      orgField: isSet(object.orgField) ? String(object.orgField) : "",
    };
  },

  toJSON(message: AuthOption): unknown {
    const obj: any = {};
    message.permission !== undefined && (obj.permission = message.permission);
    message.orgField !== undefined && (obj.orgField = message.orgField);
    return obj;
  },

  create(base?: DeepPartial<AuthOption>): AuthOption {
    return AuthOption.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AuthOption>): AuthOption {
    const message = createBaseAuthOption();
    message.permission = object.permission ?? "";
    message.orgField = object.orgField ?? "";
    return message;
  },
};

function createBaseCustomHTTPResponse(): CustomHTTPResponse {
  return { successCode: 0 };
}

export const CustomHTTPResponse = {
  encode(message: CustomHTTPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.successCode !== 0) {
      writer.uint32(8).int32(message.successCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CustomHTTPResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCustomHTTPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.successCode = reader.int32();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CustomHTTPResponse {
    return { successCode: isSet(object.successCode) ? Number(object.successCode) : 0 };
  },

  toJSON(message: CustomHTTPResponse): unknown {
    const obj: any = {};
    message.successCode !== undefined && (obj.successCode = Math.round(message.successCode));
    return obj;
  },

  create(base?: DeepPartial<CustomHTTPResponse>): CustomHTTPResponse {
    return CustomHTTPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CustomHTTPResponse>): CustomHTTPResponse {
    const message = createBaseCustomHTTPResponse();
    message.successCode = object.successCode ?? 0;
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

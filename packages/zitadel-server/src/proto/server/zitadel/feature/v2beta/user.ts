/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Details } from "../../object/v2beta/object";

export const protobufPackage = "zitadel.feature.v2beta";

export interface SetUserFeatureRequest {
  userId: string;
}

export interface SetUserFeaturesResponse {
  details: Details | undefined;
}

export interface ResetUserFeaturesRequest {
  userId: string;
}

export interface ResetUserFeaturesResponse {
  details: Details | undefined;
}

export interface GetUserFeaturesRequest {
  userId: string;
  inheritance: boolean;
}

export interface GetUserFeaturesResponse {
  details: Details | undefined;
}

function createBaseSetUserFeatureRequest(): SetUserFeatureRequest {
  return { userId: "" };
}

export const SetUserFeatureRequest = {
  encode(message: SetUserFeatureRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUserFeatureRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUserFeatureRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetUserFeatureRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: SetUserFeatureRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<SetUserFeatureRequest>): SetUserFeatureRequest {
    return SetUserFeatureRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUserFeatureRequest>): SetUserFeatureRequest {
    const message = createBaseSetUserFeatureRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseSetUserFeaturesResponse(): SetUserFeaturesResponse {
  return { details: undefined };
}

export const SetUserFeaturesResponse = {
  encode(message: SetUserFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUserFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUserFeaturesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetUserFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetUserFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetUserFeaturesResponse>): SetUserFeaturesResponse {
    return SetUserFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUserFeaturesResponse>): SetUserFeaturesResponse {
    const message = createBaseSetUserFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetUserFeaturesRequest(): ResetUserFeaturesRequest {
  return { userId: "" };
}

export const ResetUserFeaturesRequest = {
  encode(message: ResetUserFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetUserFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetUserFeaturesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ResetUserFeaturesRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: ResetUserFeaturesRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<ResetUserFeaturesRequest>): ResetUserFeaturesRequest {
    return ResetUserFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetUserFeaturesRequest>): ResetUserFeaturesRequest {
    const message = createBaseResetUserFeaturesRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseResetUserFeaturesResponse(): ResetUserFeaturesResponse {
  return { details: undefined };
}

export const ResetUserFeaturesResponse = {
  encode(message: ResetUserFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetUserFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetUserFeaturesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ResetUserFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetUserFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResetUserFeaturesResponse>): ResetUserFeaturesResponse {
    return ResetUserFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetUserFeaturesResponse>): ResetUserFeaturesResponse {
    const message = createBaseResetUserFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetUserFeaturesRequest(): GetUserFeaturesRequest {
  return { userId: "", inheritance: false };
}

export const GetUserFeaturesRequest = {
  encode(message: GetUserFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.inheritance === true) {
      writer.uint32(16).bool(message.inheritance);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserFeaturesRequest();
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
          if (tag != 16) {
            break;
          }

          message.inheritance = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetUserFeaturesRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      inheritance: isSet(object.inheritance) ? Boolean(object.inheritance) : false,
    };
  },

  toJSON(message: GetUserFeaturesRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.inheritance !== undefined && (obj.inheritance = message.inheritance);
    return obj;
  },

  create(base?: DeepPartial<GetUserFeaturesRequest>): GetUserFeaturesRequest {
    return GetUserFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetUserFeaturesRequest>): GetUserFeaturesRequest {
    const message = createBaseGetUserFeaturesRequest();
    message.userId = object.userId ?? "";
    message.inheritance = object.inheritance ?? false;
    return message;
  },
};

function createBaseGetUserFeaturesResponse(): GetUserFeaturesResponse {
  return { details: undefined };
}

export const GetUserFeaturesResponse = {
  encode(message: GetUserFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserFeaturesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetUserFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: GetUserFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetUserFeaturesResponse>): GetUserFeaturesResponse {
    return GetUserFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetUserFeaturesResponse>): GetUserFeaturesResponse {
    const message = createBaseGetUserFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
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

/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Details } from "../../object/v2beta/object";

export const protobufPackage = "zitadel.feature.v2beta";

export interface SetOrganizationFeaturesRequest {
  organizationId: string;
}

export interface SetOrganizationFeaturesResponse {
  details: Details | undefined;
}

export interface ResetOrganizationFeaturesRequest {
  organizationId: string;
}

export interface ResetOrganizationFeaturesResponse {
  details: Details | undefined;
}

export interface GetOrganizationFeaturesRequest {
  organizationId: string;
  inheritance: boolean;
}

export interface GetOrganizationFeaturesResponse {
  details: Details | undefined;
}

function createBaseSetOrganizationFeaturesRequest(): SetOrganizationFeaturesRequest {
  return { organizationId: "" };
}

export const SetOrganizationFeaturesRequest = {
  encode(message: SetOrganizationFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.organizationId !== "") {
      writer.uint32(10).string(message.organizationId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetOrganizationFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetOrganizationFeaturesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): SetOrganizationFeaturesRequest {
    return { organizationId: isSet(object.organizationId) ? String(object.organizationId) : "" };
  },

  toJSON(message: SetOrganizationFeaturesRequest): unknown {
    const obj: any = {};
    message.organizationId !== undefined && (obj.organizationId = message.organizationId);
    return obj;
  },

  create(base?: DeepPartial<SetOrganizationFeaturesRequest>): SetOrganizationFeaturesRequest {
    return SetOrganizationFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetOrganizationFeaturesRequest>): SetOrganizationFeaturesRequest {
    const message = createBaseSetOrganizationFeaturesRequest();
    message.organizationId = object.organizationId ?? "";
    return message;
  },
};

function createBaseSetOrganizationFeaturesResponse(): SetOrganizationFeaturesResponse {
  return { details: undefined };
}

export const SetOrganizationFeaturesResponse = {
  encode(message: SetOrganizationFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetOrganizationFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetOrganizationFeaturesResponse();
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

  fromJSON(object: any): SetOrganizationFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetOrganizationFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetOrganizationFeaturesResponse>): SetOrganizationFeaturesResponse {
    return SetOrganizationFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetOrganizationFeaturesResponse>): SetOrganizationFeaturesResponse {
    const message = createBaseSetOrganizationFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetOrganizationFeaturesRequest(): ResetOrganizationFeaturesRequest {
  return { organizationId: "" };
}

export const ResetOrganizationFeaturesRequest = {
  encode(message: ResetOrganizationFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.organizationId !== "") {
      writer.uint32(10).string(message.organizationId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetOrganizationFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetOrganizationFeaturesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): ResetOrganizationFeaturesRequest {
    return { organizationId: isSet(object.organizationId) ? String(object.organizationId) : "" };
  },

  toJSON(message: ResetOrganizationFeaturesRequest): unknown {
    const obj: any = {};
    message.organizationId !== undefined && (obj.organizationId = message.organizationId);
    return obj;
  },

  create(base?: DeepPartial<ResetOrganizationFeaturesRequest>): ResetOrganizationFeaturesRequest {
    return ResetOrganizationFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetOrganizationFeaturesRequest>): ResetOrganizationFeaturesRequest {
    const message = createBaseResetOrganizationFeaturesRequest();
    message.organizationId = object.organizationId ?? "";
    return message;
  },
};

function createBaseResetOrganizationFeaturesResponse(): ResetOrganizationFeaturesResponse {
  return { details: undefined };
}

export const ResetOrganizationFeaturesResponse = {
  encode(message: ResetOrganizationFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetOrganizationFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetOrganizationFeaturesResponse();
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

  fromJSON(object: any): ResetOrganizationFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetOrganizationFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResetOrganizationFeaturesResponse>): ResetOrganizationFeaturesResponse {
    return ResetOrganizationFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetOrganizationFeaturesResponse>): ResetOrganizationFeaturesResponse {
    const message = createBaseResetOrganizationFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetOrganizationFeaturesRequest(): GetOrganizationFeaturesRequest {
  return { organizationId: "", inheritance: false };
}

export const GetOrganizationFeaturesRequest = {
  encode(message: GetOrganizationFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.organizationId !== "") {
      writer.uint32(10).string(message.organizationId);
    }
    if (message.inheritance === true) {
      writer.uint32(16).bool(message.inheritance);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOrganizationFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOrganizationFeaturesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.organizationId = reader.string();
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

  fromJSON(object: any): GetOrganizationFeaturesRequest {
    return {
      organizationId: isSet(object.organizationId) ? String(object.organizationId) : "",
      inheritance: isSet(object.inheritance) ? Boolean(object.inheritance) : false,
    };
  },

  toJSON(message: GetOrganizationFeaturesRequest): unknown {
    const obj: any = {};
    message.organizationId !== undefined && (obj.organizationId = message.organizationId);
    message.inheritance !== undefined && (obj.inheritance = message.inheritance);
    return obj;
  },

  create(base?: DeepPartial<GetOrganizationFeaturesRequest>): GetOrganizationFeaturesRequest {
    return GetOrganizationFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetOrganizationFeaturesRequest>): GetOrganizationFeaturesRequest {
    const message = createBaseGetOrganizationFeaturesRequest();
    message.organizationId = object.organizationId ?? "";
    message.inheritance = object.inheritance ?? false;
    return message;
  },
};

function createBaseGetOrganizationFeaturesResponse(): GetOrganizationFeaturesResponse {
  return { details: undefined };
}

export const GetOrganizationFeaturesResponse = {
  encode(message: GetOrganizationFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetOrganizationFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetOrganizationFeaturesResponse();
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

  fromJSON(object: any): GetOrganizationFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: GetOrganizationFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetOrganizationFeaturesResponse>): GetOrganizationFeaturesResponse {
    return GetOrganizationFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetOrganizationFeaturesResponse>): GetOrganizationFeaturesResponse {
    const message = createBaseGetOrganizationFeaturesResponse();
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

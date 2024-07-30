/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Details } from "../../object/v2beta/object";
import {
  FeatureFlag,
  ImprovedPerformance,
  ImprovedPerformanceFeatureFlag,
  improvedPerformanceFromJSON,
  improvedPerformanceToJSON,
} from "./feature";

export const protobufPackage = "zitadel.feature.v2beta";

export interface SetSystemFeaturesRequest {
  loginDefaultOrg?: boolean | undefined;
  oidcTriggerIntrospectionProjections?: boolean | undefined;
  oidcLegacyIntrospection?: boolean | undefined;
  userSchema?: boolean | undefined;
  oidcTokenExchange?: boolean | undefined;
  actions?: boolean | undefined;
  improvedPerformance: ImprovedPerformance[];
}

export interface SetSystemFeaturesResponse {
  details: Details | undefined;
}

export interface ResetSystemFeaturesRequest {
}

export interface ResetSystemFeaturesResponse {
  details: Details | undefined;
}

export interface GetSystemFeaturesRequest {
}

export interface GetSystemFeaturesResponse {
  details: Details | undefined;
  loginDefaultOrg: FeatureFlag | undefined;
  oidcTriggerIntrospectionProjections: FeatureFlag | undefined;
  oidcLegacyIntrospection: FeatureFlag | undefined;
  userSchema: FeatureFlag | undefined;
  oidcTokenExchange: FeatureFlag | undefined;
  actions: FeatureFlag | undefined;
  improvedPerformance: ImprovedPerformanceFeatureFlag | undefined;
}

function createBaseSetSystemFeaturesRequest(): SetSystemFeaturesRequest {
  return {
    loginDefaultOrg: undefined,
    oidcTriggerIntrospectionProjections: undefined,
    oidcLegacyIntrospection: undefined,
    userSchema: undefined,
    oidcTokenExchange: undefined,
    actions: undefined,
    improvedPerformance: [],
  };
}

export const SetSystemFeaturesRequest = {
  encode(message: SetSystemFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.loginDefaultOrg !== undefined) {
      writer.uint32(8).bool(message.loginDefaultOrg);
    }
    if (message.oidcTriggerIntrospectionProjections !== undefined) {
      writer.uint32(16).bool(message.oidcTriggerIntrospectionProjections);
    }
    if (message.oidcLegacyIntrospection !== undefined) {
      writer.uint32(24).bool(message.oidcLegacyIntrospection);
    }
    if (message.userSchema !== undefined) {
      writer.uint32(32).bool(message.userSchema);
    }
    if (message.oidcTokenExchange !== undefined) {
      writer.uint32(40).bool(message.oidcTokenExchange);
    }
    if (message.actions !== undefined) {
      writer.uint32(48).bool(message.actions);
    }
    writer.uint32(58).fork();
    for (const v of message.improvedPerformance) {
      writer.int32(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetSystemFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetSystemFeaturesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.loginDefaultOrg = reader.bool();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.oidcTriggerIntrospectionProjections = reader.bool();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.oidcLegacyIntrospection = reader.bool();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.userSchema = reader.bool();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.oidcTokenExchange = reader.bool();
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.actions = reader.bool();
          continue;
        case 7:
          if (tag == 56) {
            message.improvedPerformance.push(reader.int32() as any);
            continue;
          }

          if (tag == 58) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.improvedPerformance.push(reader.int32() as any);
            }

            continue;
          }

          break;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetSystemFeaturesRequest {
    return {
      loginDefaultOrg: isSet(object.loginDefaultOrg) ? Boolean(object.loginDefaultOrg) : undefined,
      oidcTriggerIntrospectionProjections: isSet(object.oidcTriggerIntrospectionProjections)
        ? Boolean(object.oidcTriggerIntrospectionProjections)
        : undefined,
      oidcLegacyIntrospection: isSet(object.oidcLegacyIntrospection)
        ? Boolean(object.oidcLegacyIntrospection)
        : undefined,
      userSchema: isSet(object.userSchema) ? Boolean(object.userSchema) : undefined,
      oidcTokenExchange: isSet(object.oidcTokenExchange) ? Boolean(object.oidcTokenExchange) : undefined,
      actions: isSet(object.actions) ? Boolean(object.actions) : undefined,
      improvedPerformance: Array.isArray(object?.improvedPerformance)
        ? object.improvedPerformance.map((e: any) => improvedPerformanceFromJSON(e))
        : [],
    };
  },

  toJSON(message: SetSystemFeaturesRequest): unknown {
    const obj: any = {};
    message.loginDefaultOrg !== undefined && (obj.loginDefaultOrg = message.loginDefaultOrg);
    message.oidcTriggerIntrospectionProjections !== undefined &&
      (obj.oidcTriggerIntrospectionProjections = message.oidcTriggerIntrospectionProjections);
    message.oidcLegacyIntrospection !== undefined && (obj.oidcLegacyIntrospection = message.oidcLegacyIntrospection);
    message.userSchema !== undefined && (obj.userSchema = message.userSchema);
    message.oidcTokenExchange !== undefined && (obj.oidcTokenExchange = message.oidcTokenExchange);
    message.actions !== undefined && (obj.actions = message.actions);
    if (message.improvedPerformance) {
      obj.improvedPerformance = message.improvedPerformance.map((e) => improvedPerformanceToJSON(e));
    } else {
      obj.improvedPerformance = [];
    }
    return obj;
  },

  create(base?: DeepPartial<SetSystemFeaturesRequest>): SetSystemFeaturesRequest {
    return SetSystemFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetSystemFeaturesRequest>): SetSystemFeaturesRequest {
    const message = createBaseSetSystemFeaturesRequest();
    message.loginDefaultOrg = object.loginDefaultOrg ?? undefined;
    message.oidcTriggerIntrospectionProjections = object.oidcTriggerIntrospectionProjections ?? undefined;
    message.oidcLegacyIntrospection = object.oidcLegacyIntrospection ?? undefined;
    message.userSchema = object.userSchema ?? undefined;
    message.oidcTokenExchange = object.oidcTokenExchange ?? undefined;
    message.actions = object.actions ?? undefined;
    message.improvedPerformance = object.improvedPerformance?.map((e) => e) || [];
    return message;
  },
};

function createBaseSetSystemFeaturesResponse(): SetSystemFeaturesResponse {
  return { details: undefined };
}

export const SetSystemFeaturesResponse = {
  encode(message: SetSystemFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetSystemFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetSystemFeaturesResponse();
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

  fromJSON(object: any): SetSystemFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetSystemFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetSystemFeaturesResponse>): SetSystemFeaturesResponse {
    return SetSystemFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetSystemFeaturesResponse>): SetSystemFeaturesResponse {
    const message = createBaseSetSystemFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetSystemFeaturesRequest(): ResetSystemFeaturesRequest {
  return {};
}

export const ResetSystemFeaturesRequest = {
  encode(_: ResetSystemFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetSystemFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetSystemFeaturesRequest();
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

  fromJSON(_: any): ResetSystemFeaturesRequest {
    return {};
  },

  toJSON(_: ResetSystemFeaturesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ResetSystemFeaturesRequest>): ResetSystemFeaturesRequest {
    return ResetSystemFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ResetSystemFeaturesRequest>): ResetSystemFeaturesRequest {
    const message = createBaseResetSystemFeaturesRequest();
    return message;
  },
};

function createBaseResetSystemFeaturesResponse(): ResetSystemFeaturesResponse {
  return { details: undefined };
}

export const ResetSystemFeaturesResponse = {
  encode(message: ResetSystemFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetSystemFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetSystemFeaturesResponse();
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

  fromJSON(object: any): ResetSystemFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetSystemFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResetSystemFeaturesResponse>): ResetSystemFeaturesResponse {
    return ResetSystemFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetSystemFeaturesResponse>): ResetSystemFeaturesResponse {
    const message = createBaseResetSystemFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetSystemFeaturesRequest(): GetSystemFeaturesRequest {
  return {};
}

export const GetSystemFeaturesRequest = {
  encode(_: GetSystemFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSystemFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSystemFeaturesRequest();
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

  fromJSON(_: any): GetSystemFeaturesRequest {
    return {};
  },

  toJSON(_: GetSystemFeaturesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetSystemFeaturesRequest>): GetSystemFeaturesRequest {
    return GetSystemFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetSystemFeaturesRequest>): GetSystemFeaturesRequest {
    const message = createBaseGetSystemFeaturesRequest();
    return message;
  },
};

function createBaseGetSystemFeaturesResponse(): GetSystemFeaturesResponse {
  return {
    details: undefined,
    loginDefaultOrg: undefined,
    oidcTriggerIntrospectionProjections: undefined,
    oidcLegacyIntrospection: undefined,
    userSchema: undefined,
    oidcTokenExchange: undefined,
    actions: undefined,
    improvedPerformance: undefined,
  };
}

export const GetSystemFeaturesResponse = {
  encode(message: GetSystemFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.loginDefaultOrg !== undefined) {
      FeatureFlag.encode(message.loginDefaultOrg, writer.uint32(18).fork()).ldelim();
    }
    if (message.oidcTriggerIntrospectionProjections !== undefined) {
      FeatureFlag.encode(message.oidcTriggerIntrospectionProjections, writer.uint32(26).fork()).ldelim();
    }
    if (message.oidcLegacyIntrospection !== undefined) {
      FeatureFlag.encode(message.oidcLegacyIntrospection, writer.uint32(34).fork()).ldelim();
    }
    if (message.userSchema !== undefined) {
      FeatureFlag.encode(message.userSchema, writer.uint32(42).fork()).ldelim();
    }
    if (message.oidcTokenExchange !== undefined) {
      FeatureFlag.encode(message.oidcTokenExchange, writer.uint32(50).fork()).ldelim();
    }
    if (message.actions !== undefined) {
      FeatureFlag.encode(message.actions, writer.uint32(58).fork()).ldelim();
    }
    if (message.improvedPerformance !== undefined) {
      ImprovedPerformanceFeatureFlag.encode(message.improvedPerformance, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSystemFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSystemFeaturesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.loginDefaultOrg = FeatureFlag.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.oidcTriggerIntrospectionProjections = FeatureFlag.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.oidcLegacyIntrospection = FeatureFlag.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.userSchema = FeatureFlag.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.oidcTokenExchange = FeatureFlag.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.actions = FeatureFlag.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.improvedPerformance = ImprovedPerformanceFeatureFlag.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetSystemFeaturesResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      loginDefaultOrg: isSet(object.loginDefaultOrg) ? FeatureFlag.fromJSON(object.loginDefaultOrg) : undefined,
      oidcTriggerIntrospectionProjections: isSet(object.oidcTriggerIntrospectionProjections)
        ? FeatureFlag.fromJSON(object.oidcTriggerIntrospectionProjections)
        : undefined,
      oidcLegacyIntrospection: isSet(object.oidcLegacyIntrospection)
        ? FeatureFlag.fromJSON(object.oidcLegacyIntrospection)
        : undefined,
      userSchema: isSet(object.userSchema) ? FeatureFlag.fromJSON(object.userSchema) : undefined,
      oidcTokenExchange: isSet(object.oidcTokenExchange) ? FeatureFlag.fromJSON(object.oidcTokenExchange) : undefined,
      actions: isSet(object.actions) ? FeatureFlag.fromJSON(object.actions) : undefined,
      improvedPerformance: isSet(object.improvedPerformance)
        ? ImprovedPerformanceFeatureFlag.fromJSON(object.improvedPerformance)
        : undefined,
    };
  },

  toJSON(message: GetSystemFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.loginDefaultOrg !== undefined &&
      (obj.loginDefaultOrg = message.loginDefaultOrg ? FeatureFlag.toJSON(message.loginDefaultOrg) : undefined);
    message.oidcTriggerIntrospectionProjections !== undefined &&
      (obj.oidcTriggerIntrospectionProjections = message.oidcTriggerIntrospectionProjections
        ? FeatureFlag.toJSON(message.oidcTriggerIntrospectionProjections)
        : undefined);
    message.oidcLegacyIntrospection !== undefined && (obj.oidcLegacyIntrospection = message.oidcLegacyIntrospection
      ? FeatureFlag.toJSON(message.oidcLegacyIntrospection)
      : undefined);
    message.userSchema !== undefined &&
      (obj.userSchema = message.userSchema ? FeatureFlag.toJSON(message.userSchema) : undefined);
    message.oidcTokenExchange !== undefined &&
      (obj.oidcTokenExchange = message.oidcTokenExchange ? FeatureFlag.toJSON(message.oidcTokenExchange) : undefined);
    message.actions !== undefined && (obj.actions = message.actions ? FeatureFlag.toJSON(message.actions) : undefined);
    message.improvedPerformance !== undefined && (obj.improvedPerformance = message.improvedPerformance
      ? ImprovedPerformanceFeatureFlag.toJSON(message.improvedPerformance)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetSystemFeaturesResponse>): GetSystemFeaturesResponse {
    return GetSystemFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSystemFeaturesResponse>): GetSystemFeaturesResponse {
    const message = createBaseGetSystemFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.loginDefaultOrg = (object.loginDefaultOrg !== undefined && object.loginDefaultOrg !== null)
      ? FeatureFlag.fromPartial(object.loginDefaultOrg)
      : undefined;
    message.oidcTriggerIntrospectionProjections =
      (object.oidcTriggerIntrospectionProjections !== undefined && object.oidcTriggerIntrospectionProjections !== null)
        ? FeatureFlag.fromPartial(object.oidcTriggerIntrospectionProjections)
        : undefined;
    message.oidcLegacyIntrospection =
      (object.oidcLegacyIntrospection !== undefined && object.oidcLegacyIntrospection !== null)
        ? FeatureFlag.fromPartial(object.oidcLegacyIntrospection)
        : undefined;
    message.userSchema = (object.userSchema !== undefined && object.userSchema !== null)
      ? FeatureFlag.fromPartial(object.userSchema)
      : undefined;
    message.oidcTokenExchange = (object.oidcTokenExchange !== undefined && object.oidcTokenExchange !== null)
      ? FeatureFlag.fromPartial(object.oidcTokenExchange)
      : undefined;
    message.actions = (object.actions !== undefined && object.actions !== null)
      ? FeatureFlag.fromPartial(object.actions)
      : undefined;
    message.improvedPerformance = (object.improvedPerformance !== undefined && object.improvedPerformance !== null)
      ? ImprovedPerformanceFeatureFlag.fromPartial(object.improvedPerformance)
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

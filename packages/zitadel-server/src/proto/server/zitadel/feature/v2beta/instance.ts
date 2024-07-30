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

export interface SetInstanceFeaturesRequest {
  loginDefaultOrg?: boolean | undefined;
  oidcTriggerIntrospectionProjections?: boolean | undefined;
  oidcLegacyIntrospection?: boolean | undefined;
  userSchema?: boolean | undefined;
  oidcTokenExchange?: boolean | undefined;
  actions?: boolean | undefined;
  improvedPerformance: ImprovedPerformance[];
}

export interface SetInstanceFeaturesResponse {
  details: Details | undefined;
}

export interface ResetInstanceFeaturesRequest {
}

export interface ResetInstanceFeaturesResponse {
  details: Details | undefined;
}

export interface GetInstanceFeaturesRequest {
  inheritance: boolean;
}

export interface GetInstanceFeaturesResponse {
  details: Details | undefined;
  loginDefaultOrg: FeatureFlag | undefined;
  oidcTriggerIntrospectionProjections: FeatureFlag | undefined;
  oidcLegacyIntrospection: FeatureFlag | undefined;
  userSchema: FeatureFlag | undefined;
  oidcTokenExchange: FeatureFlag | undefined;
  actions: FeatureFlag | undefined;
  improvedPerformance: ImprovedPerformanceFeatureFlag | undefined;
}

function createBaseSetInstanceFeaturesRequest(): SetInstanceFeaturesRequest {
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

export const SetInstanceFeaturesRequest = {
  encode(message: SetInstanceFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): SetInstanceFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetInstanceFeaturesRequest();
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

  fromJSON(object: any): SetInstanceFeaturesRequest {
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

  toJSON(message: SetInstanceFeaturesRequest): unknown {
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

  create(base?: DeepPartial<SetInstanceFeaturesRequest>): SetInstanceFeaturesRequest {
    return SetInstanceFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetInstanceFeaturesRequest>): SetInstanceFeaturesRequest {
    const message = createBaseSetInstanceFeaturesRequest();
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

function createBaseSetInstanceFeaturesResponse(): SetInstanceFeaturesResponse {
  return { details: undefined };
}

export const SetInstanceFeaturesResponse = {
  encode(message: SetInstanceFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetInstanceFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetInstanceFeaturesResponse();
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

  fromJSON(object: any): SetInstanceFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetInstanceFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetInstanceFeaturesResponse>): SetInstanceFeaturesResponse {
    return SetInstanceFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetInstanceFeaturesResponse>): SetInstanceFeaturesResponse {
    const message = createBaseSetInstanceFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResetInstanceFeaturesRequest(): ResetInstanceFeaturesRequest {
  return {};
}

export const ResetInstanceFeaturesRequest = {
  encode(_: ResetInstanceFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetInstanceFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetInstanceFeaturesRequest();
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

  fromJSON(_: any): ResetInstanceFeaturesRequest {
    return {};
  },

  toJSON(_: ResetInstanceFeaturesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ResetInstanceFeaturesRequest>): ResetInstanceFeaturesRequest {
    return ResetInstanceFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ResetInstanceFeaturesRequest>): ResetInstanceFeaturesRequest {
    const message = createBaseResetInstanceFeaturesRequest();
    return message;
  },
};

function createBaseResetInstanceFeaturesResponse(): ResetInstanceFeaturesResponse {
  return { details: undefined };
}

export const ResetInstanceFeaturesResponse = {
  encode(message: ResetInstanceFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResetInstanceFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResetInstanceFeaturesResponse();
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

  fromJSON(object: any): ResetInstanceFeaturesResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResetInstanceFeaturesResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResetInstanceFeaturesResponse>): ResetInstanceFeaturesResponse {
    return ResetInstanceFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResetInstanceFeaturesResponse>): ResetInstanceFeaturesResponse {
    const message = createBaseResetInstanceFeaturesResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetInstanceFeaturesRequest(): GetInstanceFeaturesRequest {
  return { inheritance: false };
}

export const GetInstanceFeaturesRequest = {
  encode(message: GetInstanceFeaturesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.inheritance === true) {
      writer.uint32(8).bool(message.inheritance);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetInstanceFeaturesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetInstanceFeaturesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
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

  fromJSON(object: any): GetInstanceFeaturesRequest {
    return { inheritance: isSet(object.inheritance) ? Boolean(object.inheritance) : false };
  },

  toJSON(message: GetInstanceFeaturesRequest): unknown {
    const obj: any = {};
    message.inheritance !== undefined && (obj.inheritance = message.inheritance);
    return obj;
  },

  create(base?: DeepPartial<GetInstanceFeaturesRequest>): GetInstanceFeaturesRequest {
    return GetInstanceFeaturesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetInstanceFeaturesRequest>): GetInstanceFeaturesRequest {
    const message = createBaseGetInstanceFeaturesRequest();
    message.inheritance = object.inheritance ?? false;
    return message;
  },
};

function createBaseGetInstanceFeaturesResponse(): GetInstanceFeaturesResponse {
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

export const GetInstanceFeaturesResponse = {
  encode(message: GetInstanceFeaturesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): GetInstanceFeaturesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetInstanceFeaturesResponse();
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

  fromJSON(object: any): GetInstanceFeaturesResponse {
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

  toJSON(message: GetInstanceFeaturesResponse): unknown {
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

  create(base?: DeepPartial<GetInstanceFeaturesResponse>): GetInstanceFeaturesResponse {
    return GetInstanceFeaturesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetInstanceFeaturesResponse>): GetInstanceFeaturesResponse {
    const message = createBaseGetInstanceFeaturesResponse();
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

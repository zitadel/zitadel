/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Duration } from "../../../google/protobuf/duration";
import { Details } from "../../object/v2beta/object";

export const protobufPackage = "zitadel.action.v3alpha";

/** Wait for response but response body is ignored, status is checked, call is sent as post. */
export interface SetRESTWebhook {
  /** Define if any error stops the whole execution. By default the process continues as normal. */
  interruptOnError: boolean;
}

/** Wait for response and response body is used, status is checked, call is sent as post. */
export interface SetRESTCall {
  /** Define if any error stops the whole execution. By default the process continues as normal. */
  interruptOnError: boolean;
}

/** Call is executed in parallel to others, ZITADEL does not wait until the call is finished. The state is ignored, call is sent as post. */
export interface SetRESTAsync {
}

export interface Target {
  /** ID is the read-only unique identifier of the target. */
  targetId: string;
  /** Details provide some base information (such as the last change date) of the target. */
  details:
    | Details
    | undefined;
  /** Unique name of the target. */
  name: string;
  restWebhook?: SetRESTWebhook | undefined;
  restCall?: SetRESTCall | undefined;
  restAsync?:
    | SetRESTAsync
    | undefined;
  /** Timeout defines the duration until ZITADEL cancels the execution. */
  timeout: Duration | undefined;
  endpoint: string;
}

function createBaseSetRESTWebhook(): SetRESTWebhook {
  return { interruptOnError: false };
}

export const SetRESTWebhook = {
  encode(message: SetRESTWebhook, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.interruptOnError === true) {
      writer.uint32(8).bool(message.interruptOnError);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetRESTWebhook {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetRESTWebhook();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.interruptOnError = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetRESTWebhook {
    return { interruptOnError: isSet(object.interruptOnError) ? Boolean(object.interruptOnError) : false };
  },

  toJSON(message: SetRESTWebhook): unknown {
    const obj: any = {};
    message.interruptOnError !== undefined && (obj.interruptOnError = message.interruptOnError);
    return obj;
  },

  create(base?: DeepPartial<SetRESTWebhook>): SetRESTWebhook {
    return SetRESTWebhook.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetRESTWebhook>): SetRESTWebhook {
    const message = createBaseSetRESTWebhook();
    message.interruptOnError = object.interruptOnError ?? false;
    return message;
  },
};

function createBaseSetRESTCall(): SetRESTCall {
  return { interruptOnError: false };
}

export const SetRESTCall = {
  encode(message: SetRESTCall, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.interruptOnError === true) {
      writer.uint32(8).bool(message.interruptOnError);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetRESTCall {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetRESTCall();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.interruptOnError = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetRESTCall {
    return { interruptOnError: isSet(object.interruptOnError) ? Boolean(object.interruptOnError) : false };
  },

  toJSON(message: SetRESTCall): unknown {
    const obj: any = {};
    message.interruptOnError !== undefined && (obj.interruptOnError = message.interruptOnError);
    return obj;
  },

  create(base?: DeepPartial<SetRESTCall>): SetRESTCall {
    return SetRESTCall.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetRESTCall>): SetRESTCall {
    const message = createBaseSetRESTCall();
    message.interruptOnError = object.interruptOnError ?? false;
    return message;
  },
};

function createBaseSetRESTAsync(): SetRESTAsync {
  return {};
}

export const SetRESTAsync = {
  encode(_: SetRESTAsync, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetRESTAsync {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetRESTAsync();
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

  fromJSON(_: any): SetRESTAsync {
    return {};
  },

  toJSON(_: SetRESTAsync): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<SetRESTAsync>): SetRESTAsync {
    return SetRESTAsync.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<SetRESTAsync>): SetRESTAsync {
    const message = createBaseSetRESTAsync();
    return message;
  },
};

function createBaseTarget(): Target {
  return {
    targetId: "",
    details: undefined,
    name: "",
    restWebhook: undefined,
    restCall: undefined,
    restAsync: undefined,
    timeout: undefined,
    endpoint: "",
  };
}

export const Target = {
  encode(message: Target, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.targetId !== "") {
      writer.uint32(10).string(message.targetId);
    }
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    if (message.restWebhook !== undefined) {
      SetRESTWebhook.encode(message.restWebhook, writer.uint32(34).fork()).ldelim();
    }
    if (message.restCall !== undefined) {
      SetRESTCall.encode(message.restCall, writer.uint32(42).fork()).ldelim();
    }
    if (message.restAsync !== undefined) {
      SetRESTAsync.encode(message.restAsync, writer.uint32(50).fork()).ldelim();
    }
    if (message.timeout !== undefined) {
      Duration.encode(message.timeout, writer.uint32(58).fork()).ldelim();
    }
    if (message.endpoint !== "") {
      writer.uint32(66).string(message.endpoint);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Target {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTarget();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.targetId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.name = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.restWebhook = SetRESTWebhook.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.restCall = SetRESTCall.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.restAsync = SetRESTAsync.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.timeout = Duration.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.endpoint = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Target {
    return {
      targetId: isSet(object.targetId) ? String(object.targetId) : "",
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      name: isSet(object.name) ? String(object.name) : "",
      restWebhook: isSet(object.restWebhook) ? SetRESTWebhook.fromJSON(object.restWebhook) : undefined,
      restCall: isSet(object.restCall) ? SetRESTCall.fromJSON(object.restCall) : undefined,
      restAsync: isSet(object.restAsync) ? SetRESTAsync.fromJSON(object.restAsync) : undefined,
      timeout: isSet(object.timeout) ? Duration.fromJSON(object.timeout) : undefined,
      endpoint: isSet(object.endpoint) ? String(object.endpoint) : "",
    };
  },

  toJSON(message: Target): unknown {
    const obj: any = {};
    message.targetId !== undefined && (obj.targetId = message.targetId);
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.name !== undefined && (obj.name = message.name);
    message.restWebhook !== undefined &&
      (obj.restWebhook = message.restWebhook ? SetRESTWebhook.toJSON(message.restWebhook) : undefined);
    message.restCall !== undefined &&
      (obj.restCall = message.restCall ? SetRESTCall.toJSON(message.restCall) : undefined);
    message.restAsync !== undefined &&
      (obj.restAsync = message.restAsync ? SetRESTAsync.toJSON(message.restAsync) : undefined);
    message.timeout !== undefined && (obj.timeout = message.timeout ? Duration.toJSON(message.timeout) : undefined);
    message.endpoint !== undefined && (obj.endpoint = message.endpoint);
    return obj;
  },

  create(base?: DeepPartial<Target>): Target {
    return Target.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Target>): Target {
    const message = createBaseTarget();
    message.targetId = object.targetId ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.name = object.name ?? "";
    message.restWebhook = (object.restWebhook !== undefined && object.restWebhook !== null)
      ? SetRESTWebhook.fromPartial(object.restWebhook)
      : undefined;
    message.restCall = (object.restCall !== undefined && object.restCall !== null)
      ? SetRESTCall.fromPartial(object.restCall)
      : undefined;
    message.restAsync = (object.restAsync !== undefined && object.restAsync !== null)
      ? SetRESTAsync.fromPartial(object.restAsync)
      : undefined;
    message.timeout = (object.timeout !== undefined && object.timeout !== null)
      ? Duration.fromPartial(object.timeout)
      : undefined;
    message.endpoint = object.endpoint ?? "";
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

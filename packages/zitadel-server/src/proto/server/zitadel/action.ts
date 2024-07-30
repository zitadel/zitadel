/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";
import { LocalizedMessage } from "./message";
import { ObjectDetails, TextQueryMethod, textQueryMethodFromJSON, textQueryMethodToJSON } from "./object";

export const protobufPackage = "zitadel.action.v1";

export enum ActionState {
  ACTION_STATE_UNSPECIFIED = 0,
  ACTION_STATE_INACTIVE = 1,
  ACTION_STATE_ACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function actionStateFromJSON(object: any): ActionState {
  switch (object) {
    case 0:
    case "ACTION_STATE_UNSPECIFIED":
      return ActionState.ACTION_STATE_UNSPECIFIED;
    case 1:
    case "ACTION_STATE_INACTIVE":
      return ActionState.ACTION_STATE_INACTIVE;
    case 2:
    case "ACTION_STATE_ACTIVE":
      return ActionState.ACTION_STATE_ACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ActionState.UNRECOGNIZED;
  }
}

export function actionStateToJSON(object: ActionState): string {
  switch (object) {
    case ActionState.ACTION_STATE_UNSPECIFIED:
      return "ACTION_STATE_UNSPECIFIED";
    case ActionState.ACTION_STATE_INACTIVE:
      return "ACTION_STATE_INACTIVE";
    case ActionState.ACTION_STATE_ACTIVE:
      return "ACTION_STATE_ACTIVE";
    case ActionState.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ActionFieldName {
  ACTION_FIELD_NAME_UNSPECIFIED = 0,
  ACTION_FIELD_NAME_NAME = 1,
  ACTION_FIELD_NAME_ID = 2,
  ACTION_FIELD_NAME_STATE = 3,
  UNRECOGNIZED = -1,
}

export function actionFieldNameFromJSON(object: any): ActionFieldName {
  switch (object) {
    case 0:
    case "ACTION_FIELD_NAME_UNSPECIFIED":
      return ActionFieldName.ACTION_FIELD_NAME_UNSPECIFIED;
    case 1:
    case "ACTION_FIELD_NAME_NAME":
      return ActionFieldName.ACTION_FIELD_NAME_NAME;
    case 2:
    case "ACTION_FIELD_NAME_ID":
      return ActionFieldName.ACTION_FIELD_NAME_ID;
    case 3:
    case "ACTION_FIELD_NAME_STATE":
      return ActionFieldName.ACTION_FIELD_NAME_STATE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ActionFieldName.UNRECOGNIZED;
  }
}

export function actionFieldNameToJSON(object: ActionFieldName): string {
  switch (object) {
    case ActionFieldName.ACTION_FIELD_NAME_UNSPECIFIED:
      return "ACTION_FIELD_NAME_UNSPECIFIED";
    case ActionFieldName.ACTION_FIELD_NAME_NAME:
      return "ACTION_FIELD_NAME_NAME";
    case ActionFieldName.ACTION_FIELD_NAME_ID:
      return "ACTION_FIELD_NAME_ID";
    case ActionFieldName.ACTION_FIELD_NAME_STATE:
      return "ACTION_FIELD_NAME_STATE";
    case ActionFieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum FlowState {
  FLOW_STATE_UNSPECIFIED = 0,
  FLOW_STATE_INACTIVE = 1,
  FLOW_STATE_ACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function flowStateFromJSON(object: any): FlowState {
  switch (object) {
    case 0:
    case "FLOW_STATE_UNSPECIFIED":
      return FlowState.FLOW_STATE_UNSPECIFIED;
    case 1:
    case "FLOW_STATE_INACTIVE":
      return FlowState.FLOW_STATE_INACTIVE;
    case 2:
    case "FLOW_STATE_ACTIVE":
      return FlowState.FLOW_STATE_ACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FlowState.UNRECOGNIZED;
  }
}

export function flowStateToJSON(object: FlowState): string {
  switch (object) {
    case FlowState.FLOW_STATE_UNSPECIFIED:
      return "FLOW_STATE_UNSPECIFIED";
    case FlowState.FLOW_STATE_INACTIVE:
      return "FLOW_STATE_INACTIVE";
    case FlowState.FLOW_STATE_ACTIVE:
      return "FLOW_STATE_ACTIVE";
    case FlowState.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Action {
  id: string;
  details: ObjectDetails | undefined;
  state: ActionState;
  name: string;
  script: string;
  timeout: Duration | undefined;
  allowedToFail: boolean;
}

export interface ActionIDQuery {
  id: string;
}

export interface ActionNameQuery {
  name: string;
  method: TextQueryMethod;
}

/** ActionStateQuery always equals */
export interface ActionStateQuery {
  state: ActionState;
}

export interface Flow {
  /** id of the flow type */
  type: FlowType | undefined;
  details: ObjectDetails | undefined;
  state: FlowState;
  triggerActions: TriggerAction[];
}

export interface FlowType {
  /** identifier of the type */
  id: string;
  /** key and name of the type */
  name: LocalizedMessage | undefined;
}

export interface TriggerType {
  /** identifier of the type */
  id: string;
  /** key and name of the type */
  name: LocalizedMessage | undefined;
}

export interface TriggerAction {
  /** id of the trigger type */
  triggerType: TriggerType | undefined;
  actions: Action[];
}

function createBaseAction(): Action {
  return { id: "", details: undefined, state: 0, name: "", script: "", timeout: undefined, allowedToFail: false };
}

export const Action = {
  encode(message: Action, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.state !== 0) {
      writer.uint32(24).int32(message.state);
    }
    if (message.name !== "") {
      writer.uint32(34).string(message.name);
    }
    if (message.script !== "") {
      writer.uint32(42).string(message.script);
    }
    if (message.timeout !== undefined) {
      Duration.encode(message.timeout, writer.uint32(50).fork()).ldelim();
    }
    if (message.allowedToFail === true) {
      writer.uint32(56).bool(message.allowedToFail);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Action {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAction();
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

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.name = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.script = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.timeout = Duration.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 56) {
            break;
          }

          message.allowedToFail = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Action {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      state: isSet(object.state) ? actionStateFromJSON(object.state) : 0,
      name: isSet(object.name) ? String(object.name) : "",
      script: isSet(object.script) ? String(object.script) : "",
      timeout: isSet(object.timeout) ? Duration.fromJSON(object.timeout) : undefined,
      allowedToFail: isSet(object.allowedToFail) ? Boolean(object.allowedToFail) : false,
    };
  },

  toJSON(message: Action): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.state !== undefined && (obj.state = actionStateToJSON(message.state));
    message.name !== undefined && (obj.name = message.name);
    message.script !== undefined && (obj.script = message.script);
    message.timeout !== undefined && (obj.timeout = message.timeout ? Duration.toJSON(message.timeout) : undefined);
    message.allowedToFail !== undefined && (obj.allowedToFail = message.allowedToFail);
    return obj;
  },

  create(base?: DeepPartial<Action>): Action {
    return Action.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Action>): Action {
    const message = createBaseAction();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.state = object.state ?? 0;
    message.name = object.name ?? "";
    message.script = object.script ?? "";
    message.timeout = (object.timeout !== undefined && object.timeout !== null)
      ? Duration.fromPartial(object.timeout)
      : undefined;
    message.allowedToFail = object.allowedToFail ?? false;
    return message;
  },
};

function createBaseActionIDQuery(): ActionIDQuery {
  return { id: "" };
}

export const ActionIDQuery = {
  encode(message: ActionIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionIDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionIDQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionIDQuery {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: ActionIDQuery): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<ActionIDQuery>): ActionIDQuery {
    return ActionIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ActionIDQuery>): ActionIDQuery {
    const message = createBaseActionIDQuery();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseActionNameQuery(): ActionNameQuery {
  return { name: "", method: 0 };
}

export const ActionNameQuery = {
  encode(message: ActionNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionNameQuery {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: ActionNameQuery): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<ActionNameQuery>): ActionNameQuery {
    return ActionNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ActionNameQuery>): ActionNameQuery {
    const message = createBaseActionNameQuery();
    message.name = object.name ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseActionStateQuery(): ActionStateQuery {
  return { state: 0 };
}

export const ActionStateQuery = {
  encode(message: ActionStateQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.state !== 0) {
      writer.uint32(8).int32(message.state);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActionStateQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActionStateQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ActionStateQuery {
    return { state: isSet(object.state) ? actionStateFromJSON(object.state) : 0 };
  },

  toJSON(message: ActionStateQuery): unknown {
    const obj: any = {};
    message.state !== undefined && (obj.state = actionStateToJSON(message.state));
    return obj;
  },

  create(base?: DeepPartial<ActionStateQuery>): ActionStateQuery {
    return ActionStateQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ActionStateQuery>): ActionStateQuery {
    const message = createBaseActionStateQuery();
    message.state = object.state ?? 0;
    return message;
  },
};

function createBaseFlow(): Flow {
  return { type: undefined, details: undefined, state: 0, triggerActions: [] };
}

export const Flow = {
  encode(message: Flow, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== undefined) {
      FlowType.encode(message.type, writer.uint32(10).fork()).ldelim();
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.state !== 0) {
      writer.uint32(24).int32(message.state);
    }
    for (const v of message.triggerActions) {
      TriggerAction.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Flow {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFlow();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.type = FlowType.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.triggerActions.push(TriggerAction.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Flow {
    return {
      type: isSet(object.type) ? FlowType.fromJSON(object.type) : undefined,
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      state: isSet(object.state) ? flowStateFromJSON(object.state) : 0,
      triggerActions: Array.isArray(object?.triggerActions)
        ? object.triggerActions.map((e: any) => TriggerAction.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Flow): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type ? FlowType.toJSON(message.type) : undefined);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.state !== undefined && (obj.state = flowStateToJSON(message.state));
    if (message.triggerActions) {
      obj.triggerActions = message.triggerActions.map((e) => e ? TriggerAction.toJSON(e) : undefined);
    } else {
      obj.triggerActions = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Flow>): Flow {
    return Flow.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Flow>): Flow {
    const message = createBaseFlow();
    message.type = (object.type !== undefined && object.type !== null) ? FlowType.fromPartial(object.type) : undefined;
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.state = object.state ?? 0;
    message.triggerActions = object.triggerActions?.map((e) => TriggerAction.fromPartial(e)) || [];
    return message;
  },
};

function createBaseFlowType(): FlowType {
  return { id: "", name: undefined };
}

export const FlowType = {
  encode(message: FlowType, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== undefined) {
      LocalizedMessage.encode(message.name, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FlowType {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFlowType();
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

          message.name = LocalizedMessage.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FlowType {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? LocalizedMessage.fromJSON(object.name) : undefined,
    };
  },

  toJSON(message: FlowType): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name ? LocalizedMessage.toJSON(message.name) : undefined);
    return obj;
  },

  create(base?: DeepPartial<FlowType>): FlowType {
    return FlowType.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<FlowType>): FlowType {
    const message = createBaseFlowType();
    message.id = object.id ?? "";
    message.name = (object.name !== undefined && object.name !== null)
      ? LocalizedMessage.fromPartial(object.name)
      : undefined;
    return message;
  },
};

function createBaseTriggerType(): TriggerType {
  return { id: "", name: undefined };
}

export const TriggerType = {
  encode(message: TriggerType, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== undefined) {
      LocalizedMessage.encode(message.name, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TriggerType {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTriggerType();
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

          message.name = LocalizedMessage.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TriggerType {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? LocalizedMessage.fromJSON(object.name) : undefined,
    };
  },

  toJSON(message: TriggerType): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name ? LocalizedMessage.toJSON(message.name) : undefined);
    return obj;
  },

  create(base?: DeepPartial<TriggerType>): TriggerType {
    return TriggerType.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TriggerType>): TriggerType {
    const message = createBaseTriggerType();
    message.id = object.id ?? "";
    message.name = (object.name !== undefined && object.name !== null)
      ? LocalizedMessage.fromPartial(object.name)
      : undefined;
    return message;
  },
};

function createBaseTriggerAction(): TriggerAction {
  return { triggerType: undefined, actions: [] };
}

export const TriggerAction = {
  encode(message: TriggerAction, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.triggerType !== undefined) {
      TriggerType.encode(message.triggerType, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.actions) {
      Action.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TriggerAction {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTriggerAction();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.triggerType = TriggerType.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.actions.push(Action.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TriggerAction {
    return {
      triggerType: isSet(object.triggerType) ? TriggerType.fromJSON(object.triggerType) : undefined,
      actions: Array.isArray(object?.actions) ? object.actions.map((e: any) => Action.fromJSON(e)) : [],
    };
  },

  toJSON(message: TriggerAction): unknown {
    const obj: any = {};
    message.triggerType !== undefined &&
      (obj.triggerType = message.triggerType ? TriggerType.toJSON(message.triggerType) : undefined);
    if (message.actions) {
      obj.actions = message.actions.map((e) => e ? Action.toJSON(e) : undefined);
    } else {
      obj.actions = [];
    }
    return obj;
  },

  create(base?: DeepPartial<TriggerAction>): TriggerAction {
    return TriggerAction.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TriggerAction>): TriggerAction {
    const message = createBaseTriggerAction();
    message.triggerType = (object.triggerType !== undefined && object.triggerType !== null)
      ? TriggerType.fromPartial(object.triggerType)
      : undefined;
    message.actions = object.actions?.map((e) => Action.fromPartial(e)) || [];
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

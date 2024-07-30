import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as zitadel_message_pb from '../zitadel/message_pb'; // proto import: "zitadel/message.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class Action extends jspb.Message {
  getId(): string;
  setId(value: string): Action;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Action;
  hasDetails(): boolean;
  clearDetails(): Action;

  getState(): ActionState;
  setState(value: ActionState): Action;

  getName(): string;
  setName(value: string): Action;

  getScript(): string;
  setScript(value: string): Action;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): Action;
  hasTimeout(): boolean;
  clearTimeout(): Action;

  getAllowedToFail(): boolean;
  setAllowedToFail(value: boolean): Action;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Action.AsObject;
  static toObject(includeInstance: boolean, msg: Action): Action.AsObject;
  static serializeBinaryToWriter(message: Action, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Action;
  static deserializeBinaryFromReader(message: Action, reader: jspb.BinaryReader): Action;
}

export namespace Action {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: ActionState,
    name: string,
    script: string,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    allowedToFail: boolean,
  }
}

export class ActionIDQuery extends jspb.Message {
  getId(): string;
  setId(value: string): ActionIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActionIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ActionIDQuery): ActionIDQuery.AsObject;
  static serializeBinaryToWriter(message: ActionIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActionIDQuery;
  static deserializeBinaryFromReader(message: ActionIDQuery, reader: jspb.BinaryReader): ActionIDQuery;
}

export namespace ActionIDQuery {
  export type AsObject = {
    id: string,
  }
}

export class ActionNameQuery extends jspb.Message {
  getName(): string;
  setName(value: string): ActionNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): ActionNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActionNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ActionNameQuery): ActionNameQuery.AsObject;
  static serializeBinaryToWriter(message: ActionNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActionNameQuery;
  static deserializeBinaryFromReader(message: ActionNameQuery, reader: jspb.BinaryReader): ActionNameQuery;
}

export namespace ActionNameQuery {
  export type AsObject = {
    name: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class ActionStateQuery extends jspb.Message {
  getState(): ActionState;
  setState(value: ActionState): ActionStateQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActionStateQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ActionStateQuery): ActionStateQuery.AsObject;
  static serializeBinaryToWriter(message: ActionStateQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActionStateQuery;
  static deserializeBinaryFromReader(message: ActionStateQuery, reader: jspb.BinaryReader): ActionStateQuery;
}

export namespace ActionStateQuery {
  export type AsObject = {
    state: ActionState,
  }
}

export class Flow extends jspb.Message {
  getType(): FlowType | undefined;
  setType(value?: FlowType): Flow;
  hasType(): boolean;
  clearType(): Flow;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Flow;
  hasDetails(): boolean;
  clearDetails(): Flow;

  getState(): FlowState;
  setState(value: FlowState): Flow;

  getTriggerActionsList(): Array<TriggerAction>;
  setTriggerActionsList(value: Array<TriggerAction>): Flow;
  clearTriggerActionsList(): Flow;
  addTriggerActions(value?: TriggerAction, index?: number): TriggerAction;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Flow.AsObject;
  static toObject(includeInstance: boolean, msg: Flow): Flow.AsObject;
  static serializeBinaryToWriter(message: Flow, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Flow;
  static deserializeBinaryFromReader(message: Flow, reader: jspb.BinaryReader): Flow;
}

export namespace Flow {
  export type AsObject = {
    type?: FlowType.AsObject,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: FlowState,
    triggerActionsList: Array<TriggerAction.AsObject>,
  }
}

export class FlowType extends jspb.Message {
  getId(): string;
  setId(value: string): FlowType;

  getName(): zitadel_message_pb.LocalizedMessage | undefined;
  setName(value?: zitadel_message_pb.LocalizedMessage): FlowType;
  hasName(): boolean;
  clearName(): FlowType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FlowType.AsObject;
  static toObject(includeInstance: boolean, msg: FlowType): FlowType.AsObject;
  static serializeBinaryToWriter(message: FlowType, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FlowType;
  static deserializeBinaryFromReader(message: FlowType, reader: jspb.BinaryReader): FlowType;
}

export namespace FlowType {
  export type AsObject = {
    id: string,
    name?: zitadel_message_pb.LocalizedMessage.AsObject,
  }
}

export class TriggerType extends jspb.Message {
  getId(): string;
  setId(value: string): TriggerType;

  getName(): zitadel_message_pb.LocalizedMessage | undefined;
  setName(value?: zitadel_message_pb.LocalizedMessage): TriggerType;
  hasName(): boolean;
  clearName(): TriggerType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TriggerType.AsObject;
  static toObject(includeInstance: boolean, msg: TriggerType): TriggerType.AsObject;
  static serializeBinaryToWriter(message: TriggerType, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TriggerType;
  static deserializeBinaryFromReader(message: TriggerType, reader: jspb.BinaryReader): TriggerType;
}

export namespace TriggerType {
  export type AsObject = {
    id: string,
    name?: zitadel_message_pb.LocalizedMessage.AsObject,
  }
}

export class TriggerAction extends jspb.Message {
  getTriggerType(): TriggerType | undefined;
  setTriggerType(value?: TriggerType): TriggerAction;
  hasTriggerType(): boolean;
  clearTriggerType(): TriggerAction;

  getActionsList(): Array<Action>;
  setActionsList(value: Array<Action>): TriggerAction;
  clearActionsList(): TriggerAction;
  addActions(value?: Action, index?: number): Action;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TriggerAction.AsObject;
  static toObject(includeInstance: boolean, msg: TriggerAction): TriggerAction.AsObject;
  static serializeBinaryToWriter(message: TriggerAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TriggerAction;
  static deserializeBinaryFromReader(message: TriggerAction, reader: jspb.BinaryReader): TriggerAction;
}

export namespace TriggerAction {
  export type AsObject = {
    triggerType?: TriggerType.AsObject,
    actionsList: Array<Action.AsObject>,
  }
}

export enum ActionState { 
  ACTION_STATE_UNSPECIFIED = 0,
  ACTION_STATE_INACTIVE = 1,
  ACTION_STATE_ACTIVE = 2,
}
export enum ActionFieldName { 
  ACTION_FIELD_NAME_UNSPECIFIED = 0,
  ACTION_FIELD_NAME_NAME = 1,
  ACTION_FIELD_NAME_ID = 2,
  ACTION_FIELD_NAME_STATE = 3,
}
export enum FlowState { 
  FLOW_STATE_UNSPECIFIED = 0,
  FLOW_STATE_INACTIVE = 1,
  FLOW_STATE_ACTIVE = 2,
}

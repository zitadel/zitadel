/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Struct } from "../google/protobuf/struct";
import { Timestamp } from "../google/protobuf/timestamp";
import { LocalizedMessage } from "./message";

export const protobufPackage = "zitadel.event.v1";

export interface Event {
  editor: Editor | undefined;
  aggregate: Aggregate | undefined;
  sequence: number;
  creationDate: Date | undefined;
  payload: { [key: string]: any } | undefined;
  type: EventType | undefined;
}

export interface Editor {
  userId: string;
  displayName: string;
  service: string;
}

export interface Aggregate {
  id: string;
  type: AggregateType | undefined;
  resourceOwner: string;
}

export interface EventType {
  type: string;
  localized: LocalizedMessage | undefined;
}

export interface AggregateType {
  type: string;
  localized: LocalizedMessage | undefined;
}

function createBaseEvent(): Event {
  return {
    editor: undefined,
    aggregate: undefined,
    sequence: 0,
    creationDate: undefined,
    payload: undefined,
    type: undefined,
  };
}

export const Event = {
  encode(message: Event, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.editor !== undefined) {
      Editor.encode(message.editor, writer.uint32(10).fork()).ldelim();
    }
    if (message.aggregate !== undefined) {
      Aggregate.encode(message.aggregate, writer.uint32(18).fork()).ldelim();
    }
    if (message.sequence !== 0) {
      writer.uint32(24).uint64(message.sequence);
    }
    if (message.creationDate !== undefined) {
      Timestamp.encode(toTimestamp(message.creationDate), writer.uint32(34).fork()).ldelim();
    }
    if (message.payload !== undefined) {
      Struct.encode(Struct.wrap(message.payload), writer.uint32(42).fork()).ldelim();
    }
    if (message.type !== undefined) {
      EventType.encode(message.type, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Event {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.editor = Editor.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.aggregate = Aggregate.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.sequence = longToNumber(reader.uint64() as Long);
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.creationDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.payload = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.type = EventType.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Event {
    return {
      editor: isSet(object.editor) ? Editor.fromJSON(object.editor) : undefined,
      aggregate: isSet(object.aggregate) ? Aggregate.fromJSON(object.aggregate) : undefined,
      sequence: isSet(object.sequence) ? Number(object.sequence) : 0,
      creationDate: isSet(object.creationDate) ? fromJsonTimestamp(object.creationDate) : undefined,
      payload: isObject(object.payload) ? object.payload : undefined,
      type: isSet(object.type) ? EventType.fromJSON(object.type) : undefined,
    };
  },

  toJSON(message: Event): unknown {
    const obj: any = {};
    message.editor !== undefined && (obj.editor = message.editor ? Editor.toJSON(message.editor) : undefined);
    message.aggregate !== undefined &&
      (obj.aggregate = message.aggregate ? Aggregate.toJSON(message.aggregate) : undefined);
    message.sequence !== undefined && (obj.sequence = Math.round(message.sequence));
    message.creationDate !== undefined && (obj.creationDate = message.creationDate.toISOString());
    message.payload !== undefined && (obj.payload = message.payload);
    message.type !== undefined && (obj.type = message.type ? EventType.toJSON(message.type) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Event>): Event {
    return Event.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Event>): Event {
    const message = createBaseEvent();
    message.editor = (object.editor !== undefined && object.editor !== null)
      ? Editor.fromPartial(object.editor)
      : undefined;
    message.aggregate = (object.aggregate !== undefined && object.aggregate !== null)
      ? Aggregate.fromPartial(object.aggregate)
      : undefined;
    message.sequence = object.sequence ?? 0;
    message.creationDate = object.creationDate ?? undefined;
    message.payload = object.payload ?? undefined;
    message.type = (object.type !== undefined && object.type !== null) ? EventType.fromPartial(object.type) : undefined;
    return message;
  },
};

function createBaseEditor(): Editor {
  return { userId: "", displayName: "", service: "" };
}

export const Editor = {
  encode(message: Editor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.displayName !== "") {
      writer.uint32(18).string(message.displayName);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Editor {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEditor();
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
          if (tag != 18) {
            break;
          }

          message.displayName = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.service = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Editor {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      service: isSet(object.service) ? String(object.service) : "",
    };
  },

  toJSON(message: Editor): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.service !== undefined && (obj.service = message.service);
    return obj;
  },

  create(base?: DeepPartial<Editor>): Editor {
    return Editor.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Editor>): Editor {
    const message = createBaseEditor();
    message.userId = object.userId ?? "";
    message.displayName = object.displayName ?? "";
    message.service = object.service ?? "";
    return message;
  },
};

function createBaseAggregate(): Aggregate {
  return { id: "", type: undefined, resourceOwner: "" };
}

export const Aggregate = {
  encode(message: Aggregate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.type !== undefined) {
      AggregateType.encode(message.type, writer.uint32(18).fork()).ldelim();
    }
    if (message.resourceOwner !== "") {
      writer.uint32(26).string(message.resourceOwner);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Aggregate {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAggregate();
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

          message.type = AggregateType.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.resourceOwner = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Aggregate {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      type: isSet(object.type) ? AggregateType.fromJSON(object.type) : undefined,
      resourceOwner: isSet(object.resourceOwner) ? String(object.resourceOwner) : "",
    };
  },

  toJSON(message: Aggregate): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.type !== undefined && (obj.type = message.type ? AggregateType.toJSON(message.type) : undefined);
    message.resourceOwner !== undefined && (obj.resourceOwner = message.resourceOwner);
    return obj;
  },

  create(base?: DeepPartial<Aggregate>): Aggregate {
    return Aggregate.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Aggregate>): Aggregate {
    const message = createBaseAggregate();
    message.id = object.id ?? "";
    message.type = (object.type !== undefined && object.type !== null)
      ? AggregateType.fromPartial(object.type)
      : undefined;
    message.resourceOwner = object.resourceOwner ?? "";
    return message;
  },
};

function createBaseEventType(): EventType {
  return { type: "", localized: undefined };
}

export const EventType = {
  encode(message: EventType, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.localized !== undefined) {
      LocalizedMessage.encode(message.localized, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventType {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventType();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.type = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.localized = LocalizedMessage.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EventType {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      localized: isSet(object.localized) ? LocalizedMessage.fromJSON(object.localized) : undefined,
    };
  },

  toJSON(message: EventType): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.localized !== undefined &&
      (obj.localized = message.localized ? LocalizedMessage.toJSON(message.localized) : undefined);
    return obj;
  },

  create(base?: DeepPartial<EventType>): EventType {
    return EventType.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EventType>): EventType {
    const message = createBaseEventType();
    message.type = object.type ?? "";
    message.localized = (object.localized !== undefined && object.localized !== null)
      ? LocalizedMessage.fromPartial(object.localized)
      : undefined;
    return message;
  },
};

function createBaseAggregateType(): AggregateType {
  return { type: "", localized: undefined };
}

export const AggregateType = {
  encode(message: AggregateType, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== "") {
      writer.uint32(10).string(message.type);
    }
    if (message.localized !== undefined) {
      LocalizedMessage.encode(message.localized, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AggregateType {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAggregateType();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.type = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.localized = LocalizedMessage.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AggregateType {
    return {
      type: isSet(object.type) ? String(object.type) : "",
      localized: isSet(object.localized) ? LocalizedMessage.fromJSON(object.localized) : undefined,
    };
  },

  toJSON(message: AggregateType): unknown {
    const obj: any = {};
    message.type !== undefined && (obj.type = message.type);
    message.localized !== undefined &&
      (obj.localized = message.localized ? LocalizedMessage.toJSON(message.localized) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AggregateType>): AggregateType {
    return AggregateType.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AggregateType>): AggregateType {
    const message = createBaseAggregateType();
    message.type = object.type ?? "";
    message.localized = (object.localized !== undefined && object.localized !== null)
      ? LocalizedMessage.fromPartial(object.localized)
      : undefined;
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

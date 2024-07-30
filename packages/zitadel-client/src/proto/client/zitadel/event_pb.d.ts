import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as zitadel_message_pb from '../zitadel/message_pb'; // proto import: "zitadel/message.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class Event extends jspb.Message {
  getEditor(): Editor | undefined;
  setEditor(value?: Editor): Event;
  hasEditor(): boolean;
  clearEditor(): Event;

  getAggregate(): Aggregate | undefined;
  setAggregate(value?: Aggregate): Event;
  hasAggregate(): boolean;
  clearAggregate(): Event;

  getSequence(): number;
  setSequence(value: number): Event;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): Event;
  hasCreationDate(): boolean;
  clearCreationDate(): Event;

  getPayload(): google_protobuf_struct_pb.Struct | undefined;
  setPayload(value?: google_protobuf_struct_pb.Struct): Event;
  hasPayload(): boolean;
  clearPayload(): Event;

  getType(): EventType | undefined;
  setType(value?: EventType): Event;
  hasType(): boolean;
  clearType(): Event;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
  export type AsObject = {
    editor?: Editor.AsObject,
    aggregate?: Aggregate.AsObject,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    payload?: google_protobuf_struct_pb.Struct.AsObject,
    type?: EventType.AsObject,
  }
}

export class Editor extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): Editor;

  getDisplayName(): string;
  setDisplayName(value: string): Editor;

  getService(): string;
  setService(value: string): Editor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Editor.AsObject;
  static toObject(includeInstance: boolean, msg: Editor): Editor.AsObject;
  static serializeBinaryToWriter(message: Editor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Editor;
  static deserializeBinaryFromReader(message: Editor, reader: jspb.BinaryReader): Editor;
}

export namespace Editor {
  export type AsObject = {
    userId: string,
    displayName: string,
    service: string,
  }
}

export class Aggregate extends jspb.Message {
  getId(): string;
  setId(value: string): Aggregate;

  getType(): AggregateType | undefined;
  setType(value?: AggregateType): Aggregate;
  hasType(): boolean;
  clearType(): Aggregate;

  getResourceOwner(): string;
  setResourceOwner(value: string): Aggregate;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Aggregate.AsObject;
  static toObject(includeInstance: boolean, msg: Aggregate): Aggregate.AsObject;
  static serializeBinaryToWriter(message: Aggregate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Aggregate;
  static deserializeBinaryFromReader(message: Aggregate, reader: jspb.BinaryReader): Aggregate;
}

export namespace Aggregate {
  export type AsObject = {
    id: string,
    type?: AggregateType.AsObject,
    resourceOwner: string,
  }
}

export class EventType extends jspb.Message {
  getType(): string;
  setType(value: string): EventType;

  getLocalized(): zitadel_message_pb.LocalizedMessage | undefined;
  setLocalized(value?: zitadel_message_pb.LocalizedMessage): EventType;
  hasLocalized(): boolean;
  clearLocalized(): EventType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EventType.AsObject;
  static toObject(includeInstance: boolean, msg: EventType): EventType.AsObject;
  static serializeBinaryToWriter(message: EventType, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EventType;
  static deserializeBinaryFromReader(message: EventType, reader: jspb.BinaryReader): EventType;
}

export namespace EventType {
  export type AsObject = {
    type: string,
    localized?: zitadel_message_pb.LocalizedMessage.AsObject,
  }
}

export class AggregateType extends jspb.Message {
  getType(): string;
  setType(value: string): AggregateType;

  getLocalized(): zitadel_message_pb.LocalizedMessage | undefined;
  setLocalized(value?: zitadel_message_pb.LocalizedMessage): AggregateType;
  hasLocalized(): boolean;
  clearLocalized(): AggregateType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AggregateType.AsObject;
  static toObject(includeInstance: boolean, msg: AggregateType): AggregateType.AsObject;
  static serializeBinaryToWriter(message: AggregateType, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AggregateType;
  static deserializeBinaryFromReader(message: AggregateType, reader: jspb.BinaryReader): AggregateType;
}

export namespace AggregateType {
  export type AsObject = {
    type: string,
    localized?: zitadel_message_pb.LocalizedMessage.AsObject,
  }
}


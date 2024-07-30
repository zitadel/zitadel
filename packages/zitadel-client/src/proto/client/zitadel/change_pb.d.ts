import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as zitadel_message_pb from '../zitadel/message_pb'; // proto import: "zitadel/message.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class Change extends jspb.Message {
  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): Change;
  hasChangeDate(): boolean;
  clearChangeDate(): Change;

  getEventType(): zitadel_message_pb.LocalizedMessage | undefined;
  setEventType(value?: zitadel_message_pb.LocalizedMessage): Change;
  hasEventType(): boolean;
  clearEventType(): Change;

  getSequence(): number;
  setSequence(value: number): Change;

  getEditorId(): string;
  setEditorId(value: string): Change;

  getEditorDisplayName(): string;
  setEditorDisplayName(value: string): Change;

  getResourceOwnerId(): string;
  setResourceOwnerId(value: string): Change;

  getEditorPreferredLoginName(): string;
  setEditorPreferredLoginName(value: string): Change;

  getEditorAvatarUrl(): string;
  setEditorAvatarUrl(value: string): Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Change.AsObject;
  static toObject(includeInstance: boolean, msg: Change): Change.AsObject;
  static serializeBinaryToWriter(message: Change, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Change;
  static deserializeBinaryFromReader(message: Change, reader: jspb.BinaryReader): Change;
}

export namespace Change {
  export type AsObject = {
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    eventType?: zitadel_message_pb.LocalizedMessage.AsObject,
    sequence: number,
    editorId: string,
    editorDisplayName: string,
    resourceOwnerId: string,
    editorPreferredLoginName: string,
    editorAvatarUrl: string,
  }
}

export class ChangeQuery extends jspb.Message {
  getSequence(): number;
  setSequence(value: number): ChangeQuery;

  getLimit(): number;
  setLimit(value: number): ChangeQuery;

  getAsc(): boolean;
  setAsc(value: boolean): ChangeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeQuery): ChangeQuery.AsObject;
  static serializeBinaryToWriter(message: ChangeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeQuery;
  static deserializeBinaryFromReader(message: ChangeQuery, reader: jspb.BinaryReader): ChangeQuery;
}

export namespace ChangeQuery {
  export type AsObject = {
    sequence: number,
    limit: number,
    asc: boolean,
  }
}


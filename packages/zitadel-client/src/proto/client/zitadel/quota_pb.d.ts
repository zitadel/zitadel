import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"


export class Notification extends jspb.Message {
  getPercent(): number;
  setPercent(value: number): Notification;

  getRepeat(): boolean;
  setRepeat(value: boolean): Notification;

  getCallUrl(): string;
  setCallUrl(value: string): Notification;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Notification.AsObject;
  static toObject(includeInstance: boolean, msg: Notification): Notification.AsObject;
  static serializeBinaryToWriter(message: Notification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Notification;
  static deserializeBinaryFromReader(message: Notification, reader: jspb.BinaryReader): Notification;
}

export namespace Notification {
  export type AsObject = {
    percent: number,
    repeat: boolean,
    callUrl: string,
  }
}

export enum Unit { 
  UNIT_UNIMPLEMENTED = 0,
  UNIT_REQUESTS_ALL_AUTHENTICATED = 1,
  UNIT_ACTIONS_ALL_RUN_SECONDS = 2,
}

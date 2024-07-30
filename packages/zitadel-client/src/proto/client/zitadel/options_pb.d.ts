import * as jspb from 'google-protobuf'

import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb'; // proto import: "google/protobuf/descriptor.proto"


export class AuthOption extends jspb.Message {
  getPermission(): string;
  setPermission(value: string): AuthOption;

  getCheckFieldName(): string;
  setCheckFieldName(value: string): AuthOption;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthOption.AsObject;
  static toObject(includeInstance: boolean, msg: AuthOption): AuthOption.AsObject;
  static serializeBinaryToWriter(message: AuthOption, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthOption;
  static deserializeBinaryFromReader(message: AuthOption, reader: jspb.BinaryReader): AuthOption;
}

export namespace AuthOption {
  export type AsObject = {
    permission: string,
    checkFieldName: string,
  }
}


import * as jspb from 'google-protobuf'

import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb'; // proto import: "google/protobuf/descriptor.proto"


export class Options extends jspb.Message {
  getAuthOption(): AuthOption | undefined;
  setAuthOption(value?: AuthOption): Options;
  hasAuthOption(): boolean;
  clearAuthOption(): Options;

  getHttpResponse(): CustomHTTPResponse | undefined;
  setHttpResponse(value?: CustomHTTPResponse): Options;
  hasHttpResponse(): boolean;
  clearHttpResponse(): Options;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Options.AsObject;
  static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
  static serializeBinaryToWriter(message: Options, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Options;
  static deserializeBinaryFromReader(message: Options, reader: jspb.BinaryReader): Options;
}

export namespace Options {
  export type AsObject = {
    authOption?: AuthOption.AsObject,
    httpResponse?: CustomHTTPResponse.AsObject,
  }
}

export class AuthOption extends jspb.Message {
  getPermission(): string;
  setPermission(value: string): AuthOption;

  getOrgField(): string;
  setOrgField(value: string): AuthOption;

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
    orgField: string,
  }
}

export class CustomHTTPResponse extends jspb.Message {
  getSuccessCode(): number;
  setSuccessCode(value: number): CustomHTTPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CustomHTTPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CustomHTTPResponse): CustomHTTPResponse.AsObject;
  static serializeBinaryToWriter(message: CustomHTTPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CustomHTTPResponse;
  static deserializeBinaryFromReader(message: CustomHTTPResponse, reader: jspb.BinaryReader): CustomHTTPResponse;
}

export namespace CustomHTTPResponse {
  export type AsObject = {
    successCode: number,
  }
}


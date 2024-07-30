import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class SecuritySettings extends jspb.Message {
  getEmbeddedIframe(): EmbeddedIframeSettings | undefined;
  setEmbeddedIframe(value?: EmbeddedIframeSettings): SecuritySettings;
  hasEmbeddedIframe(): boolean;
  clearEmbeddedIframe(): SecuritySettings;

  getEnableImpersonation(): boolean;
  setEnableImpersonation(value: boolean): SecuritySettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SecuritySettings.AsObject;
  static toObject(includeInstance: boolean, msg: SecuritySettings): SecuritySettings.AsObject;
  static serializeBinaryToWriter(message: SecuritySettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SecuritySettings;
  static deserializeBinaryFromReader(message: SecuritySettings, reader: jspb.BinaryReader): SecuritySettings;
}

export namespace SecuritySettings {
  export type AsObject = {
    embeddedIframe?: EmbeddedIframeSettings.AsObject,
    enableImpersonation: boolean,
  }
}

export class EmbeddedIframeSettings extends jspb.Message {
  getEnabled(): boolean;
  setEnabled(value: boolean): EmbeddedIframeSettings;

  getAllowedOriginsList(): Array<string>;
  setAllowedOriginsList(value: Array<string>): EmbeddedIframeSettings;
  clearAllowedOriginsList(): EmbeddedIframeSettings;
  addAllowedOrigins(value: string, index?: number): EmbeddedIframeSettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmbeddedIframeSettings.AsObject;
  static toObject(includeInstance: boolean, msg: EmbeddedIframeSettings): EmbeddedIframeSettings.AsObject;
  static serializeBinaryToWriter(message: EmbeddedIframeSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmbeddedIframeSettings;
  static deserializeBinaryFromReader(message: EmbeddedIframeSettings, reader: jspb.BinaryReader): EmbeddedIframeSettings;
}

export namespace EmbeddedIframeSettings {
  export type AsObject = {
    enabled: boolean,
    allowedOriginsList: Array<string>,
  }
}


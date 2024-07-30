import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_protoc_gen_zitadel_v2_options_pb from '../../../zitadel/protoc_gen_zitadel/v2/options_pb'; // proto import: "zitadel/protoc_gen_zitadel/v2/options.proto"


export class SetRESTWebhook extends jspb.Message {
  getInterruptOnError(): boolean;
  setInterruptOnError(value: boolean): SetRESTWebhook;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRESTWebhook.AsObject;
  static toObject(includeInstance: boolean, msg: SetRESTWebhook): SetRESTWebhook.AsObject;
  static serializeBinaryToWriter(message: SetRESTWebhook, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRESTWebhook;
  static deserializeBinaryFromReader(message: SetRESTWebhook, reader: jspb.BinaryReader): SetRESTWebhook;
}

export namespace SetRESTWebhook {
  export type AsObject = {
    interruptOnError: boolean,
  }
}

export class SetRESTCall extends jspb.Message {
  getInterruptOnError(): boolean;
  setInterruptOnError(value: boolean): SetRESTCall;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRESTCall.AsObject;
  static toObject(includeInstance: boolean, msg: SetRESTCall): SetRESTCall.AsObject;
  static serializeBinaryToWriter(message: SetRESTCall, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRESTCall;
  static deserializeBinaryFromReader(message: SetRESTCall, reader: jspb.BinaryReader): SetRESTCall;
}

export namespace SetRESTCall {
  export type AsObject = {
    interruptOnError: boolean,
  }
}

export class SetRESTAsync extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRESTAsync.AsObject;
  static toObject(includeInstance: boolean, msg: SetRESTAsync): SetRESTAsync.AsObject;
  static serializeBinaryToWriter(message: SetRESTAsync, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRESTAsync;
  static deserializeBinaryFromReader(message: SetRESTAsync, reader: jspb.BinaryReader): SetRESTAsync;
}

export namespace SetRESTAsync {
  export type AsObject = {
  }
}

export class Target extends jspb.Message {
  getTargetId(): string;
  setTargetId(value: string): Target;

  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): Target;
  hasDetails(): boolean;
  clearDetails(): Target;

  getName(): string;
  setName(value: string): Target;

  getRestWebhook(): SetRESTWebhook | undefined;
  setRestWebhook(value?: SetRESTWebhook): Target;
  hasRestWebhook(): boolean;
  clearRestWebhook(): Target;

  getRestCall(): SetRESTCall | undefined;
  setRestCall(value?: SetRESTCall): Target;
  hasRestCall(): boolean;
  clearRestCall(): Target;

  getRestAsync(): SetRESTAsync | undefined;
  setRestAsync(value?: SetRESTAsync): Target;
  hasRestAsync(): boolean;
  clearRestAsync(): Target;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): Target;
  hasTimeout(): boolean;
  clearTimeout(): Target;

  getEndpoint(): string;
  setEndpoint(value: string): Target;

  getTargetTypeCase(): Target.TargetTypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Target.AsObject;
  static toObject(includeInstance: boolean, msg: Target): Target.AsObject;
  static serializeBinaryToWriter(message: Target, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Target;
  static deserializeBinaryFromReader(message: Target, reader: jspb.BinaryReader): Target;
}

export namespace Target {
  export type AsObject = {
    targetId: string,
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    name: string,
    restWebhook?: SetRESTWebhook.AsObject,
    restCall?: SetRESTCall.AsObject,
    restAsync?: SetRESTAsync.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    endpoint: string,
  }

  export enum TargetTypeCase { 
    TARGET_TYPE_NOT_SET = 0,
    REST_WEBHOOK = 4,
    REST_CALL = 5,
    REST_ASYNC = 6,
  }
}


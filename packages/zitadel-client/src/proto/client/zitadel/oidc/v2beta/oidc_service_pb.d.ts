import * as jspb from 'google-protobuf'

import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_protoc_gen_zitadel_v2_options_pb from '../../../zitadel/protoc_gen_zitadel/v2/options_pb'; // proto import: "zitadel/protoc_gen_zitadel/v2/options.proto"
import * as zitadel_oidc_v2beta_authorization_pb from '../../../zitadel/oidc/v2beta/authorization_pb'; // proto import: "zitadel/oidc/v2beta/authorization.proto"
import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class GetAuthRequestRequest extends jspb.Message {
  getAuthRequestId(): string;
  setAuthRequestId(value: string): GetAuthRequestRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthRequestRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAuthRequestRequest): GetAuthRequestRequest.AsObject;
  static serializeBinaryToWriter(message: GetAuthRequestRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthRequestRequest;
  static deserializeBinaryFromReader(message: GetAuthRequestRequest, reader: jspb.BinaryReader): GetAuthRequestRequest;
}

export namespace GetAuthRequestRequest {
  export type AsObject = {
    authRequestId: string,
  }
}

export class GetAuthRequestResponse extends jspb.Message {
  getAuthRequest(): zitadel_oidc_v2beta_authorization_pb.AuthRequest | undefined;
  setAuthRequest(value?: zitadel_oidc_v2beta_authorization_pb.AuthRequest): GetAuthRequestResponse;
  hasAuthRequest(): boolean;
  clearAuthRequest(): GetAuthRequestResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAuthRequestResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAuthRequestResponse): GetAuthRequestResponse.AsObject;
  static serializeBinaryToWriter(message: GetAuthRequestResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAuthRequestResponse;
  static deserializeBinaryFromReader(message: GetAuthRequestResponse, reader: jspb.BinaryReader): GetAuthRequestResponse;
}

export namespace GetAuthRequestResponse {
  export type AsObject = {
    authRequest?: zitadel_oidc_v2beta_authorization_pb.AuthRequest.AsObject,
  }
}

export class CreateCallbackRequest extends jspb.Message {
  getAuthRequestId(): string;
  setAuthRequestId(value: string): CreateCallbackRequest;

  getSession(): Session | undefined;
  setSession(value?: Session): CreateCallbackRequest;
  hasSession(): boolean;
  clearSession(): CreateCallbackRequest;

  getError(): zitadel_oidc_v2beta_authorization_pb.AuthorizationError | undefined;
  setError(value?: zitadel_oidc_v2beta_authorization_pb.AuthorizationError): CreateCallbackRequest;
  hasError(): boolean;
  clearError(): CreateCallbackRequest;

  getCallbackKindCase(): CreateCallbackRequest.CallbackKindCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateCallbackRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateCallbackRequest): CreateCallbackRequest.AsObject;
  static serializeBinaryToWriter(message: CreateCallbackRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateCallbackRequest;
  static deserializeBinaryFromReader(message: CreateCallbackRequest, reader: jspb.BinaryReader): CreateCallbackRequest;
}

export namespace CreateCallbackRequest {
  export type AsObject = {
    authRequestId: string,
    session?: Session.AsObject,
    error?: zitadel_oidc_v2beta_authorization_pb.AuthorizationError.AsObject,
  }

  export enum CallbackKindCase { 
    CALLBACK_KIND_NOT_SET = 0,
    SESSION = 2,
    ERROR = 3,
  }
}

export class Session extends jspb.Message {
  getSessionId(): string;
  setSessionId(value: string): Session;

  getSessionToken(): string;
  setSessionToken(value: string): Session;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Session.AsObject;
  static toObject(includeInstance: boolean, msg: Session): Session.AsObject;
  static serializeBinaryToWriter(message: Session, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Session;
  static deserializeBinaryFromReader(message: Session, reader: jspb.BinaryReader): Session;
}

export namespace Session {
  export type AsObject = {
    sessionId: string,
    sessionToken: string,
  }
}

export class CreateCallbackResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): CreateCallbackResponse;
  hasDetails(): boolean;
  clearDetails(): CreateCallbackResponse;

  getCallbackUrl(): string;
  setCallbackUrl(value: string): CreateCallbackResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateCallbackResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateCallbackResponse): CreateCallbackResponse.AsObject;
  static serializeBinaryToWriter(message: CreateCallbackResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateCallbackResponse;
  static deserializeBinaryFromReader(message: CreateCallbackResponse, reader: jspb.BinaryReader): CreateCallbackResponse;
}

export namespace CreateCallbackResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    callbackUrl: string,
  }
}


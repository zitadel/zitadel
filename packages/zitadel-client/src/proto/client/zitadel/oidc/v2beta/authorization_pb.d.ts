import * as jspb from 'google-protobuf'

import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class AuthRequest extends jspb.Message {
  getId(): string;
  setId(value: string): AuthRequest;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): AuthRequest;
  hasCreationDate(): boolean;
  clearCreationDate(): AuthRequest;

  getClientId(): string;
  setClientId(value: string): AuthRequest;

  getScopeList(): Array<string>;
  setScopeList(value: Array<string>): AuthRequest;
  clearScopeList(): AuthRequest;
  addScope(value: string, index?: number): AuthRequest;

  getRedirectUri(): string;
  setRedirectUri(value: string): AuthRequest;

  getPromptList(): Array<Prompt>;
  setPromptList(value: Array<Prompt>): AuthRequest;
  clearPromptList(): AuthRequest;
  addPrompt(value: Prompt, index?: number): AuthRequest;

  getUiLocalesList(): Array<string>;
  setUiLocalesList(value: Array<string>): AuthRequest;
  clearUiLocalesList(): AuthRequest;
  addUiLocales(value: string, index?: number): AuthRequest;

  getLoginHint(): string;
  setLoginHint(value: string): AuthRequest;
  hasLoginHint(): boolean;
  clearLoginHint(): AuthRequest;

  getMaxAge(): google_protobuf_duration_pb.Duration | undefined;
  setMaxAge(value?: google_protobuf_duration_pb.Duration): AuthRequest;
  hasMaxAge(): boolean;
  clearMaxAge(): AuthRequest;

  getHintUserId(): string;
  setHintUserId(value: string): AuthRequest;
  hasHintUserId(): boolean;
  clearHintUserId(): AuthRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AuthRequest): AuthRequest.AsObject;
  static serializeBinaryToWriter(message: AuthRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthRequest;
  static deserializeBinaryFromReader(message: AuthRequest, reader: jspb.BinaryReader): AuthRequest;
}

export namespace AuthRequest {
  export type AsObject = {
    id: string,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    clientId: string,
    scopeList: Array<string>,
    redirectUri: string,
    promptList: Array<Prompt>,
    uiLocalesList: Array<string>,
    loginHint?: string,
    maxAge?: google_protobuf_duration_pb.Duration.AsObject,
    hintUserId?: string,
  }

  export enum LoginHintCase { 
    _LOGIN_HINT_NOT_SET = 0,
    LOGIN_HINT = 8,
  }

  export enum MaxAgeCase { 
    _MAX_AGE_NOT_SET = 0,
    MAX_AGE = 9,
  }

  export enum HintUserIdCase { 
    _HINT_USER_ID_NOT_SET = 0,
    HINT_USER_ID = 10,
  }
}

export class AuthorizationError extends jspb.Message {
  getError(): ErrorReason;
  setError(value: ErrorReason): AuthorizationError;

  getErrorDescription(): string;
  setErrorDescription(value: string): AuthorizationError;
  hasErrorDescription(): boolean;
  clearErrorDescription(): AuthorizationError;

  getErrorUri(): string;
  setErrorUri(value: string): AuthorizationError;
  hasErrorUri(): boolean;
  clearErrorUri(): AuthorizationError;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthorizationError.AsObject;
  static toObject(includeInstance: boolean, msg: AuthorizationError): AuthorizationError.AsObject;
  static serializeBinaryToWriter(message: AuthorizationError, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthorizationError;
  static deserializeBinaryFromReader(message: AuthorizationError, reader: jspb.BinaryReader): AuthorizationError;
}

export namespace AuthorizationError {
  export type AsObject = {
    error: ErrorReason,
    errorDescription?: string,
    errorUri?: string,
  }

  export enum ErrorDescriptionCase { 
    _ERROR_DESCRIPTION_NOT_SET = 0,
    ERROR_DESCRIPTION = 2,
  }

  export enum ErrorUriCase { 
    _ERROR_URI_NOT_SET = 0,
    ERROR_URI = 3,
  }
}

export enum Prompt { 
  PROMPT_UNSPECIFIED = 0,
  PROMPT_NONE = 1,
  PROMPT_LOGIN = 2,
  PROMPT_CONSENT = 3,
  PROMPT_SELECT_ACCOUNT = 4,
  PROMPT_CREATE = 5,
}
export enum ErrorReason { 
  ERROR_REASON_UNSPECIFIED = 0,
  ERROR_REASON_INVALID_REQUEST = 1,
  ERROR_REASON_UNAUTHORIZED_CLIENT = 2,
  ERROR_REASON_ACCESS_DENIED = 3,
  ERROR_REASON_UNSUPPORTED_RESPONSE_TYPE = 4,
  ERROR_REASON_INVALID_SCOPE = 5,
  ERROR_REASON_SERVER_ERROR = 6,
  ERROR_REASON_TEMPORARY_UNAVAILABLE = 7,
  ERROR_REASON_INTERACTION_REQUIRED = 8,
  ERROR_REASON_LOGIN_REQUIRED = 9,
  ERROR_REASON_ACCOUNT_SELECTION_REQUIRED = 10,
  ERROR_REASON_CONSENT_REQUIRED = 11,
  ERROR_REASON_INVALID_REQUEST_URI = 12,
  ERROR_REASON_INVALID_REQUEST_OBJECT = 13,
  ERROR_REASON_REQUEST_NOT_SUPPORTED = 14,
  ERROR_REASON_REQUEST_URI_NOT_SUPPORTED = 15,
  ERROR_REASON_REGISTRATION_NOT_SUPPORTED = 16,
}

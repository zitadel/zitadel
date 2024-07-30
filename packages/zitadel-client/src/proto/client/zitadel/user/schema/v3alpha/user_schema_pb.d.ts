import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as validate_validate_pb from '../../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as zitadel_object_v2beta_object_pb from '../../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"


export class UserSchema extends jspb.Message {
  getId(): string;
  setId(value: string): UserSchema;

  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): UserSchema;
  hasDetails(): boolean;
  clearDetails(): UserSchema;

  getType(): string;
  setType(value: string): UserSchema;

  getState(): State;
  setState(value: State): UserSchema;

  getRevision(): number;
  setRevision(value: number): UserSchema;

  getSchema(): google_protobuf_struct_pb.Struct | undefined;
  setSchema(value?: google_protobuf_struct_pb.Struct): UserSchema;
  hasSchema(): boolean;
  clearSchema(): UserSchema;

  getPossibleAuthenticatorsList(): Array<AuthenticatorType>;
  setPossibleAuthenticatorsList(value: Array<AuthenticatorType>): UserSchema;
  clearPossibleAuthenticatorsList(): UserSchema;
  addPossibleAuthenticators(value: AuthenticatorType, index?: number): UserSchema;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSchema.AsObject;
  static toObject(includeInstance: boolean, msg: UserSchema): UserSchema.AsObject;
  static serializeBinaryToWriter(message: UserSchema, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSchema;
  static deserializeBinaryFromReader(message: UserSchema, reader: jspb.BinaryReader): UserSchema;
}

export namespace UserSchema {
  export type AsObject = {
    id: string,
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    type: string,
    state: State,
    revision: number,
    schema?: google_protobuf_struct_pb.Struct.AsObject,
    possibleAuthenticatorsList: Array<AuthenticatorType>,
  }
}

export class SearchQuery extends jspb.Message {
  getOrQuery(): OrQuery | undefined;
  setOrQuery(value?: OrQuery): SearchQuery;
  hasOrQuery(): boolean;
  clearOrQuery(): SearchQuery;

  getAndQuery(): AndQuery | undefined;
  setAndQuery(value?: AndQuery): SearchQuery;
  hasAndQuery(): boolean;
  clearAndQuery(): SearchQuery;

  getNotQuery(): NotQuery | undefined;
  setNotQuery(value?: NotQuery): SearchQuery;
  hasNotQuery(): boolean;
  clearNotQuery(): SearchQuery;

  getTypeQuery(): TypeQuery | undefined;
  setTypeQuery(value?: TypeQuery): SearchQuery;
  hasTypeQuery(): boolean;
  clearTypeQuery(): SearchQuery;

  getStateQuery(): StateQuery | undefined;
  setStateQuery(value?: StateQuery): SearchQuery;
  hasStateQuery(): boolean;
  clearStateQuery(): SearchQuery;

  getIdQuery(): IDQuery | undefined;
  setIdQuery(value?: IDQuery): SearchQuery;
  hasIdQuery(): boolean;
  clearIdQuery(): SearchQuery;

  getQueryCase(): SearchQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: SearchQuery): SearchQuery.AsObject;
  static serializeBinaryToWriter(message: SearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchQuery;
  static deserializeBinaryFromReader(message: SearchQuery, reader: jspb.BinaryReader): SearchQuery;
}

export namespace SearchQuery {
  export type AsObject = {
    orQuery?: OrQuery.AsObject,
    andQuery?: AndQuery.AsObject,
    notQuery?: NotQuery.AsObject,
    typeQuery?: TypeQuery.AsObject,
    stateQuery?: StateQuery.AsObject,
    idQuery?: IDQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    OR_QUERY = 1,
    AND_QUERY = 2,
    NOT_QUERY = 3,
    TYPE_QUERY = 5,
    STATE_QUERY = 6,
    ID_QUERY = 7,
  }
}

export class OrQuery extends jspb.Message {
  getQueriesList(): Array<SearchQuery>;
  setQueriesList(value: Array<SearchQuery>): OrQuery;
  clearQueriesList(): OrQuery;
  addQueries(value?: SearchQuery, index?: number): SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrQuery): OrQuery.AsObject;
  static serializeBinaryToWriter(message: OrQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrQuery;
  static deserializeBinaryFromReader(message: OrQuery, reader: jspb.BinaryReader): OrQuery;
}

export namespace OrQuery {
  export type AsObject = {
    queriesList: Array<SearchQuery.AsObject>,
  }
}

export class AndQuery extends jspb.Message {
  getQueriesList(): Array<SearchQuery>;
  setQueriesList(value: Array<SearchQuery>): AndQuery;
  clearQueriesList(): AndQuery;
  addQueries(value?: SearchQuery, index?: number): SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AndQuery.AsObject;
  static toObject(includeInstance: boolean, msg: AndQuery): AndQuery.AsObject;
  static serializeBinaryToWriter(message: AndQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AndQuery;
  static deserializeBinaryFromReader(message: AndQuery, reader: jspb.BinaryReader): AndQuery;
}

export namespace AndQuery {
  export type AsObject = {
    queriesList: Array<SearchQuery.AsObject>,
  }
}

export class NotQuery extends jspb.Message {
  getQuery(): SearchQuery | undefined;
  setQuery(value?: SearchQuery): NotQuery;
  hasQuery(): boolean;
  clearQuery(): NotQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotQuery.AsObject;
  static toObject(includeInstance: boolean, msg: NotQuery): NotQuery.AsObject;
  static serializeBinaryToWriter(message: NotQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotQuery;
  static deserializeBinaryFromReader(message: NotQuery, reader: jspb.BinaryReader): NotQuery;
}

export namespace NotQuery {
  export type AsObject = {
    query?: SearchQuery.AsObject,
  }
}

export class IDQuery extends jspb.Message {
  getId(): string;
  setId(value: string): IDQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): IDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IDQuery): IDQuery.AsObject;
  static serializeBinaryToWriter(message: IDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDQuery;
  static deserializeBinaryFromReader(message: IDQuery, reader: jspb.BinaryReader): IDQuery;
}

export namespace IDQuery {
  export type AsObject = {
    id: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class TypeQuery extends jspb.Message {
  getType(): string;
  setType(value: string): TypeQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): TypeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TypeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: TypeQuery): TypeQuery.AsObject;
  static serializeBinaryToWriter(message: TypeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TypeQuery;
  static deserializeBinaryFromReader(message: TypeQuery, reader: jspb.BinaryReader): TypeQuery;
}

export namespace TypeQuery {
  export type AsObject = {
    type: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class StateQuery extends jspb.Message {
  getState(): State;
  setState(value: State): StateQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StateQuery.AsObject;
  static toObject(includeInstance: boolean, msg: StateQuery): StateQuery.AsObject;
  static serializeBinaryToWriter(message: StateQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StateQuery;
  static deserializeBinaryFromReader(message: StateQuery, reader: jspb.BinaryReader): StateQuery;
}

export namespace StateQuery {
  export type AsObject = {
    state: State,
  }
}

export enum FieldName { 
  FIELD_NAME_UNSPECIFIED = 0,
  FIELD_NAME_TYPE = 1,
  FIELD_NAME_STATE = 2,
  FIELD_NAME_REVISION = 3,
  FIELD_NAME_CHANGE_DATE = 4,
}
export enum State { 
  STATE_UNSPECIFIED = 0,
  STATE_ACTIVE = 1,
  STATE_INACTIVE = 2,
}
export enum AuthenticatorType { 
  AUTHENTICATOR_TYPE_UNSPECIFIED = 0,
  AUTHENTICATOR_TYPE_USERNAME = 1,
  AUTHENTICATOR_TYPE_PASSWORD = 2,
  AUTHENTICATOR_TYPE_WEBAUTHN = 3,
  AUTHENTICATOR_TYPE_TOTP = 4,
  AUTHENTICATOR_TYPE_OTP_EMAIL = 5,
  AUTHENTICATOR_TYPE_OTP_SMS = 6,
  AUTHENTICATOR_TYPE_AUTHENTICATION_KEY = 7,
  AUTHENTICATOR_TYPE_IDENTITY_PROVIDER = 8,
}

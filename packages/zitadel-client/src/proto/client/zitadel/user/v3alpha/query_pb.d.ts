import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_user_v3alpha_user_pb from '../../../zitadel/user/v3alpha/user_pb'; // proto import: "zitadel/user/v3alpha/user.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"


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

  getUserIdQuery(): UserIDQuery | undefined;
  setUserIdQuery(value?: UserIDQuery): SearchQuery;
  hasUserIdQuery(): boolean;
  clearUserIdQuery(): SearchQuery;

  getOrganizationIdQuery(): OrganizationIDQuery | undefined;
  setOrganizationIdQuery(value?: OrganizationIDQuery): SearchQuery;
  hasOrganizationIdQuery(): boolean;
  clearOrganizationIdQuery(): SearchQuery;

  getUsernameQuery(): UsernameQuery | undefined;
  setUsernameQuery(value?: UsernameQuery): SearchQuery;
  hasUsernameQuery(): boolean;
  clearUsernameQuery(): SearchQuery;

  getEmailQuery(): EmailQuery | undefined;
  setEmailQuery(value?: EmailQuery): SearchQuery;
  hasEmailQuery(): boolean;
  clearEmailQuery(): SearchQuery;

  getPhoneQuery(): PhoneQuery | undefined;
  setPhoneQuery(value?: PhoneQuery): SearchQuery;
  hasPhoneQuery(): boolean;
  clearPhoneQuery(): SearchQuery;

  getStateQuery(): StateQuery | undefined;
  setStateQuery(value?: StateQuery): SearchQuery;
  hasStateQuery(): boolean;
  clearStateQuery(): SearchQuery;

  getSchemaIdQuery(): SchemaIDQuery | undefined;
  setSchemaIdQuery(value?: SchemaIDQuery): SearchQuery;
  hasSchemaIdQuery(): boolean;
  clearSchemaIdQuery(): SearchQuery;

  getSchemaTypeQuery(): SchemaTypeQuery | undefined;
  setSchemaTypeQuery(value?: SchemaTypeQuery): SearchQuery;
  hasSchemaTypeQuery(): boolean;
  clearSchemaTypeQuery(): SearchQuery;

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
    userIdQuery?: UserIDQuery.AsObject,
    organizationIdQuery?: OrganizationIDQuery.AsObject,
    usernameQuery?: UsernameQuery.AsObject,
    emailQuery?: EmailQuery.AsObject,
    phoneQuery?: PhoneQuery.AsObject,
    stateQuery?: StateQuery.AsObject,
    schemaIdQuery?: SchemaIDQuery.AsObject,
    schemaTypeQuery?: SchemaTypeQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    OR_QUERY = 1,
    AND_QUERY = 2,
    NOT_QUERY = 3,
    USER_ID_QUERY = 4,
    ORGANIZATION_ID_QUERY = 5,
    USERNAME_QUERY = 6,
    EMAIL_QUERY = 7,
    PHONE_QUERY = 8,
    STATE_QUERY = 9,
    SCHEMA_ID_QUERY = 10,
    SCHEMA_TYPE_QUERY = 11,
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

export class UserIDQuery extends jspb.Message {
  getId(): string;
  setId(value: string): UserIDQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): UserIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserIDQuery): UserIDQuery.AsObject;
  static serializeBinaryToWriter(message: UserIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserIDQuery;
  static deserializeBinaryFromReader(message: UserIDQuery, reader: jspb.BinaryReader): UserIDQuery;
}

export namespace UserIDQuery {
  export type AsObject = {
    id: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class OrganizationIDQuery extends jspb.Message {
  getId(): string;
  setId(value: string): OrganizationIDQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): OrganizationIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrganizationIDQuery): OrganizationIDQuery.AsObject;
  static serializeBinaryToWriter(message: OrganizationIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationIDQuery;
  static deserializeBinaryFromReader(message: OrganizationIDQuery, reader: jspb.BinaryReader): OrganizationIDQuery;
}

export namespace OrganizationIDQuery {
  export type AsObject = {
    id: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class UsernameQuery extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): UsernameQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): UsernameQuery;

  getIsOrganizationSpecific(): boolean;
  setIsOrganizationSpecific(value: boolean): UsernameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UsernameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UsernameQuery): UsernameQuery.AsObject;
  static serializeBinaryToWriter(message: UsernameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UsernameQuery;
  static deserializeBinaryFromReader(message: UsernameQuery, reader: jspb.BinaryReader): UsernameQuery;
}

export namespace UsernameQuery {
  export type AsObject = {
    username: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
    isOrganizationSpecific: boolean,
  }
}

export class EmailQuery extends jspb.Message {
  getAddress(): string;
  setAddress(value: string): EmailQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): EmailQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmailQuery.AsObject;
  static toObject(includeInstance: boolean, msg: EmailQuery): EmailQuery.AsObject;
  static serializeBinaryToWriter(message: EmailQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmailQuery;
  static deserializeBinaryFromReader(message: EmailQuery, reader: jspb.BinaryReader): EmailQuery;
}

export namespace EmailQuery {
  export type AsObject = {
    address: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class PhoneQuery extends jspb.Message {
  getNumber(): string;
  setNumber(value: string): PhoneQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): PhoneQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PhoneQuery.AsObject;
  static toObject(includeInstance: boolean, msg: PhoneQuery): PhoneQuery.AsObject;
  static serializeBinaryToWriter(message: PhoneQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PhoneQuery;
  static deserializeBinaryFromReader(message: PhoneQuery, reader: jspb.BinaryReader): PhoneQuery;
}

export namespace PhoneQuery {
  export type AsObject = {
    number: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class StateQuery extends jspb.Message {
  getState(): zitadel_user_v3alpha_user_pb.State;
  setState(value: zitadel_user_v3alpha_user_pb.State): StateQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StateQuery.AsObject;
  static toObject(includeInstance: boolean, msg: StateQuery): StateQuery.AsObject;
  static serializeBinaryToWriter(message: StateQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StateQuery;
  static deserializeBinaryFromReader(message: StateQuery, reader: jspb.BinaryReader): StateQuery;
}

export namespace StateQuery {
  export type AsObject = {
    state: zitadel_user_v3alpha_user_pb.State,
  }
}

export class SchemaIDQuery extends jspb.Message {
  getId(): string;
  setId(value: string): SchemaIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaIDQuery): SchemaIDQuery.AsObject;
  static serializeBinaryToWriter(message: SchemaIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaIDQuery;
  static deserializeBinaryFromReader(message: SchemaIDQuery, reader: jspb.BinaryReader): SchemaIDQuery;
}

export namespace SchemaIDQuery {
  export type AsObject = {
    id: string,
  }
}

export class SchemaTypeQuery extends jspb.Message {
  getType(): string;
  setType(value: string): SchemaTypeQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): SchemaTypeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SchemaTypeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: SchemaTypeQuery): SchemaTypeQuery.AsObject;
  static serializeBinaryToWriter(message: SchemaTypeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SchemaTypeQuery;
  static deserializeBinaryFromReader(message: SchemaTypeQuery, reader: jspb.BinaryReader): SchemaTypeQuery;
}

export namespace SchemaTypeQuery {
  export type AsObject = {
    type: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export enum FieldName { 
  FIELD_NAME_UNSPECIFIED = 0,
  FIELD_NAME_ID = 1,
  FIELD_NAME_CREATION_DATE = 2,
  FIELD_NAME_CHANGE_DATE = 3,
  FIELD_NAME_EMAIL = 4,
  FIELD_NAME_PHONE = 5,
  FIELD_NAME_STATE = 6,
  FIELD_NAME_SCHEMA_ID = 7,
  FIELD_NAME_SCHEMA_TYPE = 8,
}

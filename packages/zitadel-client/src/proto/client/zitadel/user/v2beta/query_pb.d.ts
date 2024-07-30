import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_user_v2beta_user_pb from '../../../zitadel/user/v2beta/user_pb'; // proto import: "zitadel/user/v2beta/user.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"


export class SearchQuery extends jspb.Message {
  getUserNameQuery(): UserNameQuery | undefined;
  setUserNameQuery(value?: UserNameQuery): SearchQuery;
  hasUserNameQuery(): boolean;
  clearUserNameQuery(): SearchQuery;

  getFirstNameQuery(): FirstNameQuery | undefined;
  setFirstNameQuery(value?: FirstNameQuery): SearchQuery;
  hasFirstNameQuery(): boolean;
  clearFirstNameQuery(): SearchQuery;

  getLastNameQuery(): LastNameQuery | undefined;
  setLastNameQuery(value?: LastNameQuery): SearchQuery;
  hasLastNameQuery(): boolean;
  clearLastNameQuery(): SearchQuery;

  getNickNameQuery(): NickNameQuery | undefined;
  setNickNameQuery(value?: NickNameQuery): SearchQuery;
  hasNickNameQuery(): boolean;
  clearNickNameQuery(): SearchQuery;

  getDisplayNameQuery(): DisplayNameQuery | undefined;
  setDisplayNameQuery(value?: DisplayNameQuery): SearchQuery;
  hasDisplayNameQuery(): boolean;
  clearDisplayNameQuery(): SearchQuery;

  getEmailQuery(): EmailQuery | undefined;
  setEmailQuery(value?: EmailQuery): SearchQuery;
  hasEmailQuery(): boolean;
  clearEmailQuery(): SearchQuery;

  getStateQuery(): StateQuery | undefined;
  setStateQuery(value?: StateQuery): SearchQuery;
  hasStateQuery(): boolean;
  clearStateQuery(): SearchQuery;

  getTypeQuery(): TypeQuery | undefined;
  setTypeQuery(value?: TypeQuery): SearchQuery;
  hasTypeQuery(): boolean;
  clearTypeQuery(): SearchQuery;

  getLoginNameQuery(): LoginNameQuery | undefined;
  setLoginNameQuery(value?: LoginNameQuery): SearchQuery;
  hasLoginNameQuery(): boolean;
  clearLoginNameQuery(): SearchQuery;

  getInUserIdsQuery(): InUserIDQuery | undefined;
  setInUserIdsQuery(value?: InUserIDQuery): SearchQuery;
  hasInUserIdsQuery(): boolean;
  clearInUserIdsQuery(): SearchQuery;

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

  getInUserEmailsQuery(): InUserEmailsQuery | undefined;
  setInUserEmailsQuery(value?: InUserEmailsQuery): SearchQuery;
  hasInUserEmailsQuery(): boolean;
  clearInUserEmailsQuery(): SearchQuery;

  getOrganizationIdQuery(): OrganizationIdQuery | undefined;
  setOrganizationIdQuery(value?: OrganizationIdQuery): SearchQuery;
  hasOrganizationIdQuery(): boolean;
  clearOrganizationIdQuery(): SearchQuery;

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
    userNameQuery?: UserNameQuery.AsObject,
    firstNameQuery?: FirstNameQuery.AsObject,
    lastNameQuery?: LastNameQuery.AsObject,
    nickNameQuery?: NickNameQuery.AsObject,
    displayNameQuery?: DisplayNameQuery.AsObject,
    emailQuery?: EmailQuery.AsObject,
    stateQuery?: StateQuery.AsObject,
    typeQuery?: TypeQuery.AsObject,
    loginNameQuery?: LoginNameQuery.AsObject,
    inUserIdsQuery?: InUserIDQuery.AsObject,
    orQuery?: OrQuery.AsObject,
    andQuery?: AndQuery.AsObject,
    notQuery?: NotQuery.AsObject,
    inUserEmailsQuery?: InUserEmailsQuery.AsObject,
    organizationIdQuery?: OrganizationIdQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    USER_NAME_QUERY = 1,
    FIRST_NAME_QUERY = 2,
    LAST_NAME_QUERY = 3,
    NICK_NAME_QUERY = 4,
    DISPLAY_NAME_QUERY = 5,
    EMAIL_QUERY = 6,
    STATE_QUERY = 7,
    TYPE_QUERY = 8,
    LOGIN_NAME_QUERY = 9,
    IN_USER_IDS_QUERY = 10,
    OR_QUERY = 11,
    AND_QUERY = 12,
    NOT_QUERY = 13,
    IN_USER_EMAILS_QUERY = 14,
    ORGANIZATION_ID_QUERY = 15,
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

export class InUserIDQuery extends jspb.Message {
  getUserIdsList(): Array<string>;
  setUserIdsList(value: Array<string>): InUserIDQuery;
  clearUserIdsList(): InUserIDQuery;
  addUserIds(value: string, index?: number): InUserIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InUserIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: InUserIDQuery): InUserIDQuery.AsObject;
  static serializeBinaryToWriter(message: InUserIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InUserIDQuery;
  static deserializeBinaryFromReader(message: InUserIDQuery, reader: jspb.BinaryReader): InUserIDQuery;
}

export namespace InUserIDQuery {
  export type AsObject = {
    userIdsList: Array<string>,
  }
}

export class UserNameQuery extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): UserNameQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): UserNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserNameQuery): UserNameQuery.AsObject;
  static serializeBinaryToWriter(message: UserNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserNameQuery;
  static deserializeBinaryFromReader(message: UserNameQuery, reader: jspb.BinaryReader): UserNameQuery;
}

export namespace UserNameQuery {
  export type AsObject = {
    userName: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class FirstNameQuery extends jspb.Message {
  getFirstName(): string;
  setFirstName(value: string): FirstNameQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): FirstNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FirstNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: FirstNameQuery): FirstNameQuery.AsObject;
  static serializeBinaryToWriter(message: FirstNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FirstNameQuery;
  static deserializeBinaryFromReader(message: FirstNameQuery, reader: jspb.BinaryReader): FirstNameQuery;
}

export namespace FirstNameQuery {
  export type AsObject = {
    firstName: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class LastNameQuery extends jspb.Message {
  getLastName(): string;
  setLastName(value: string): LastNameQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): LastNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LastNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: LastNameQuery): LastNameQuery.AsObject;
  static serializeBinaryToWriter(message: LastNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LastNameQuery;
  static deserializeBinaryFromReader(message: LastNameQuery, reader: jspb.BinaryReader): LastNameQuery;
}

export namespace LastNameQuery {
  export type AsObject = {
    lastName: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class NickNameQuery extends jspb.Message {
  getNickName(): string;
  setNickName(value: string): NickNameQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): NickNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NickNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: NickNameQuery): NickNameQuery.AsObject;
  static serializeBinaryToWriter(message: NickNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NickNameQuery;
  static deserializeBinaryFromReader(message: NickNameQuery, reader: jspb.BinaryReader): NickNameQuery;
}

export namespace NickNameQuery {
  export type AsObject = {
    nickName: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class DisplayNameQuery extends jspb.Message {
  getDisplayName(): string;
  setDisplayName(value: string): DisplayNameQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): DisplayNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisplayNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: DisplayNameQuery): DisplayNameQuery.AsObject;
  static serializeBinaryToWriter(message: DisplayNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisplayNameQuery;
  static deserializeBinaryFromReader(message: DisplayNameQuery, reader: jspb.BinaryReader): DisplayNameQuery;
}

export namespace DisplayNameQuery {
  export type AsObject = {
    displayName: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class EmailQuery extends jspb.Message {
  getEmailAddress(): string;
  setEmailAddress(value: string): EmailQuery;

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
    emailAddress: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class LoginNameQuery extends jspb.Message {
  getLoginName(): string;
  setLoginName(value: string): LoginNameQuery;

  getMethod(): zitadel_object_v2beta_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_v2beta_object_pb.TextQueryMethod): LoginNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: LoginNameQuery): LoginNameQuery.AsObject;
  static serializeBinaryToWriter(message: LoginNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginNameQuery;
  static deserializeBinaryFromReader(message: LoginNameQuery, reader: jspb.BinaryReader): LoginNameQuery;
}

export namespace LoginNameQuery {
  export type AsObject = {
    loginName: string,
    method: zitadel_object_v2beta_object_pb.TextQueryMethod,
  }
}

export class StateQuery extends jspb.Message {
  getState(): zitadel_user_v2beta_user_pb.UserState;
  setState(value: zitadel_user_v2beta_user_pb.UserState): StateQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StateQuery.AsObject;
  static toObject(includeInstance: boolean, msg: StateQuery): StateQuery.AsObject;
  static serializeBinaryToWriter(message: StateQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StateQuery;
  static deserializeBinaryFromReader(message: StateQuery, reader: jspb.BinaryReader): StateQuery;
}

export namespace StateQuery {
  export type AsObject = {
    state: zitadel_user_v2beta_user_pb.UserState,
  }
}

export class TypeQuery extends jspb.Message {
  getType(): Type;
  setType(value: Type): TypeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TypeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: TypeQuery): TypeQuery.AsObject;
  static serializeBinaryToWriter(message: TypeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TypeQuery;
  static deserializeBinaryFromReader(message: TypeQuery, reader: jspb.BinaryReader): TypeQuery;
}

export namespace TypeQuery {
  export type AsObject = {
    type: Type,
  }
}

export class InUserEmailsQuery extends jspb.Message {
  getUserEmailsList(): Array<string>;
  setUserEmailsList(value: Array<string>): InUserEmailsQuery;
  clearUserEmailsList(): InUserEmailsQuery;
  addUserEmails(value: string, index?: number): InUserEmailsQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InUserEmailsQuery.AsObject;
  static toObject(includeInstance: boolean, msg: InUserEmailsQuery): InUserEmailsQuery.AsObject;
  static serializeBinaryToWriter(message: InUserEmailsQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InUserEmailsQuery;
  static deserializeBinaryFromReader(message: InUserEmailsQuery, reader: jspb.BinaryReader): InUserEmailsQuery;
}

export namespace InUserEmailsQuery {
  export type AsObject = {
    userEmailsList: Array<string>,
  }
}

export class OrganizationIdQuery extends jspb.Message {
  getOrganizationId(): string;
  setOrganizationId(value: string): OrganizationIdQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrganizationIdQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrganizationIdQuery): OrganizationIdQuery.AsObject;
  static serializeBinaryToWriter(message: OrganizationIdQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrganizationIdQuery;
  static deserializeBinaryFromReader(message: OrganizationIdQuery, reader: jspb.BinaryReader): OrganizationIdQuery;
}

export namespace OrganizationIdQuery {
  export type AsObject = {
    organizationId: string,
  }
}

export enum Type { 
  TYPE_UNSPECIFIED = 0,
  TYPE_HUMAN = 1,
  TYPE_MACHINE = 2,
}
export enum UserFieldName { 
  USER_FIELD_NAME_UNSPECIFIED = 0,
  USER_FIELD_NAME_USER_NAME = 1,
  USER_FIELD_NAME_FIRST_NAME = 2,
  USER_FIELD_NAME_LAST_NAME = 3,
  USER_FIELD_NAME_NICK_NAME = 4,
  USER_FIELD_NAME_DISPLAY_NAME = 5,
  USER_FIELD_NAME_EMAIL = 6,
  USER_FIELD_NAME_STATE = 7,
  USER_FIELD_NAME_TYPE = 8,
  USER_FIELD_NAME_CREATION_DATE = 9,
}

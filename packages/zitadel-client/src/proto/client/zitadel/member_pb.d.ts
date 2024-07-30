import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as zitadel_user_pb from '../zitadel/user_pb'; // proto import: "zitadel/user.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class Member extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): Member;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Member;
  hasDetails(): boolean;
  clearDetails(): Member;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): Member;
  clearRolesList(): Member;
  addRoles(value: string, index?: number): Member;

  getPreferredLoginName(): string;
  setPreferredLoginName(value: string): Member;

  getEmail(): string;
  setEmail(value: string): Member;

  getFirstName(): string;
  setFirstName(value: string): Member;

  getLastName(): string;
  setLastName(value: string): Member;

  getDisplayName(): string;
  setDisplayName(value: string): Member;

  getAvatarUrl(): string;
  setAvatarUrl(value: string): Member;

  getUserType(): zitadel_user_pb.Type;
  setUserType(value: zitadel_user_pb.Type): Member;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Member.AsObject;
  static toObject(includeInstance: boolean, msg: Member): Member.AsObject;
  static serializeBinaryToWriter(message: Member, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Member;
  static deserializeBinaryFromReader(message: Member, reader: jspb.BinaryReader): Member;
}

export namespace Member {
  export type AsObject = {
    userId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    rolesList: Array<string>,
    preferredLoginName: string,
    email: string,
    firstName: string,
    lastName: string,
    displayName: string,
    avatarUrl: string,
    userType: zitadel_user_pb.Type,
  }
}

export class SearchQuery extends jspb.Message {
  getFirstNameQuery(): FirstNameQuery | undefined;
  setFirstNameQuery(value?: FirstNameQuery): SearchQuery;
  hasFirstNameQuery(): boolean;
  clearFirstNameQuery(): SearchQuery;

  getLastNameQuery(): LastNameQuery | undefined;
  setLastNameQuery(value?: LastNameQuery): SearchQuery;
  hasLastNameQuery(): boolean;
  clearLastNameQuery(): SearchQuery;

  getEmailQuery(): EmailQuery | undefined;
  setEmailQuery(value?: EmailQuery): SearchQuery;
  hasEmailQuery(): boolean;
  clearEmailQuery(): SearchQuery;

  getUserIdQuery(): UserIDQuery | undefined;
  setUserIdQuery(value?: UserIDQuery): SearchQuery;
  hasUserIdQuery(): boolean;
  clearUserIdQuery(): SearchQuery;

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
    firstNameQuery?: FirstNameQuery.AsObject,
    lastNameQuery?: LastNameQuery.AsObject,
    emailQuery?: EmailQuery.AsObject,
    userIdQuery?: UserIDQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    FIRST_NAME_QUERY = 1,
    LAST_NAME_QUERY = 2,
    EMAIL_QUERY = 3,
    USER_ID_QUERY = 4,
  }
}

export class FirstNameQuery extends jspb.Message {
  getFirstName(): string;
  setFirstName(value: string): FirstNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): FirstNameQuery;

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
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class LastNameQuery extends jspb.Message {
  getLastName(): string;
  setLastName(value: string): LastNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): LastNameQuery;

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
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class EmailQuery extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): EmailQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): EmailQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmailQuery.AsObject;
  static toObject(includeInstance: boolean, msg: EmailQuery): EmailQuery.AsObject;
  static serializeBinaryToWriter(message: EmailQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmailQuery;
  static deserializeBinaryFromReader(message: EmailQuery, reader: jspb.BinaryReader): EmailQuery;
}

export namespace EmailQuery {
  export type AsObject = {
    email: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserIDQuery extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UserIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserIDQuery): UserIDQuery.AsObject;
  static serializeBinaryToWriter(message: UserIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserIDQuery;
  static deserializeBinaryFromReader(message: UserIDQuery, reader: jspb.BinaryReader): UserIDQuery;
}

export namespace UserIDQuery {
  export type AsObject = {
    userId: string,
  }
}


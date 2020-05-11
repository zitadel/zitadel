import * as jspb from "google-protobuf"

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as authoption_options_pb from './authoption/options_pb';

export class OrgID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgID.AsObject;
  static toObject(includeInstance: boolean, msg: OrgID): OrgID.AsObject;
  static serializeBinaryToWriter(message: OrgID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgID;
  static deserializeBinaryFromReader(message: OrgID, reader: jspb.BinaryReader): OrgID;
}

export namespace OrgID {
  export type AsObject = {
    id: string,
  }
}

export class UniqueOrgRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDomain(): string;
  setDomain(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UniqueOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UniqueOrgRequest): UniqueOrgRequest.AsObject;
  static serializeBinaryToWriter(message: UniqueOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UniqueOrgRequest;
  static deserializeBinaryFromReader(message: UniqueOrgRequest, reader: jspb.BinaryReader): UniqueOrgRequest;
}

export namespace UniqueOrgRequest {
  export type AsObject = {
    name: string,
    domain: string,
  }
}

export class UniqueOrgResponse extends jspb.Message {
  getIsUnique(): boolean;
  setIsUnique(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UniqueOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UniqueOrgResponse): UniqueOrgResponse.AsObject;
  static serializeBinaryToWriter(message: UniqueOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UniqueOrgResponse;
  static deserializeBinaryFromReader(message: UniqueOrgResponse, reader: jspb.BinaryReader): UniqueOrgResponse;
}

export namespace UniqueOrgResponse {
  export type AsObject = {
    isUnique: boolean,
  }
}

export class Org extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getState(): OrgState;
  setState(value: OrgState): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getName(): string;
  setName(value: string): void;

  getDomain(): string;
  setDomain(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Org.AsObject;
  static toObject(includeInstance: boolean, msg: Org): Org.AsObject;
  static serializeBinaryToWriter(message: Org, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Org;
  static deserializeBinaryFromReader(message: Org, reader: jspb.BinaryReader): Org;
}

export namespace Org {
  export type AsObject = {
    id: string,
    state: OrgState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    name: string,
    domain: string,
  }
}

export class OrgSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getSortingColumn(): OrgSearchKey;
  setSortingColumn(value: OrgSearchKey): void;

  getAsc(): boolean;
  setAsc(value: boolean): void;

  getQueriesList(): Array<OrgSearchQuery>;
  setQueriesList(value: Array<OrgSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: OrgSearchQuery, index?: number): OrgSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: OrgSearchRequest): OrgSearchRequest.AsObject;
  static serializeBinaryToWriter(message: OrgSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgSearchRequest;
  static deserializeBinaryFromReader(message: OrgSearchRequest, reader: jspb.BinaryReader): OrgSearchRequest;
}

export namespace OrgSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    sortingColumn: OrgSearchKey,
    asc: boolean,
    queriesList: Array<OrgSearchQuery.AsObject>,
  }
}

export class OrgSearchQuery extends jspb.Message {
  getKey(): OrgSearchKey;
  setKey(value: OrgSearchKey): void;

  getMethod(): OrgSearchMethod;
  setMethod(value: OrgSearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrgSearchQuery): OrgSearchQuery.AsObject;
  static serializeBinaryToWriter(message: OrgSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgSearchQuery;
  static deserializeBinaryFromReader(message: OrgSearchQuery, reader: jspb.BinaryReader): OrgSearchQuery;
}

export namespace OrgSearchQuery {
  export type AsObject = {
    key: OrgSearchKey,
    method: OrgSearchMethod,
    value: string,
  }
}

export class OrgSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<Org>;
  setResultList(value: Array<Org>): void;
  clearResultList(): void;
  addResult(value?: Org, index?: number): Org;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: OrgSearchResponse): OrgSearchResponse.AsObject;
  static serializeBinaryToWriter(message: OrgSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgSearchResponse;
  static deserializeBinaryFromReader(message: OrgSearchResponse, reader: jspb.BinaryReader): OrgSearchResponse;
}

export namespace OrgSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<Org.AsObject>,
  }
}

export class OrgSetUpRequest extends jspb.Message {
  getOrg(): CreateOrgRequest | undefined;
  setOrg(value?: CreateOrgRequest): void;
  hasOrg(): boolean;
  clearOrg(): void;

  getUser(): RegisterUserRequest | undefined;
  setUser(value?: RegisterUserRequest): void;
  hasUser(): boolean;
  clearUser(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgSetUpRequest.AsObject;
  static toObject(includeInstance: boolean, msg: OrgSetUpRequest): OrgSetUpRequest.AsObject;
  static serializeBinaryToWriter(message: OrgSetUpRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgSetUpRequest;
  static deserializeBinaryFromReader(message: OrgSetUpRequest, reader: jspb.BinaryReader): OrgSetUpRequest;
}

export namespace OrgSetUpRequest {
  export type AsObject = {
    org?: CreateOrgRequest.AsObject,
    user?: RegisterUserRequest.AsObject,
  }
}

export class OrgSetUpResponse extends jspb.Message {
  getOrg(): Org | undefined;
  setOrg(value?: Org): void;
  hasOrg(): boolean;
  clearOrg(): void;

  getUser(): User | undefined;
  setUser(value?: User): void;
  hasUser(): boolean;
  clearUser(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgSetUpResponse.AsObject;
  static toObject(includeInstance: boolean, msg: OrgSetUpResponse): OrgSetUpResponse.AsObject;
  static serializeBinaryToWriter(message: OrgSetUpResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgSetUpResponse;
  static deserializeBinaryFromReader(message: OrgSetUpResponse, reader: jspb.BinaryReader): OrgSetUpResponse;
}

export namespace OrgSetUpResponse {
  export type AsObject = {
    org?: Org.AsObject,
    user?: User.AsObject,
  }
}

export class RegisterUserRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getNickName(): string;
  setNickName(value: string): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): void;

  getGender(): Gender;
  setGender(value: Gender): void;

  getPassword(): string;
  setPassword(value: string): void;

  getOrgId(): string;
  setOrgId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterUserRequest): RegisterUserRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterUserRequest;
  static deserializeBinaryFromReader(message: RegisterUserRequest, reader: jspb.BinaryReader): RegisterUserRequest;
}

export namespace RegisterUserRequest {
  export type AsObject = {
    email: string,
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: Gender,
    password: string,
    orgId: string,
  }
}

export class User extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getState(): UserState;
  setState(value: UserState): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getUserName(): string;
  setUserName(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getNickName(): string;
  setNickName(value: string): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): void;

  getGender(): Gender;
  setGender(value: Gender): void;

  getEmail(): string;
  setEmail(value: string): void;

  getIsemailverified(): boolean;
  setIsemailverified(value: boolean): void;

  getPhone(): string;
  setPhone(value: string): void;

  getIsphoneverified(): boolean;
  setIsphoneverified(value: boolean): void;

  getCountry(): string;
  setCountry(value: string): void;

  getLocality(): string;
  setLocality(value: string): void;

  getPostalCode(): string;
  setPostalCode(value: string): void;

  getRegion(): string;
  setRegion(value: string): void;

  getStreetAddress(): string;
  setStreetAddress(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    id: string,
    state: UserState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    userName: string,
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: Gender,
    email: string,
    isemailverified: boolean,
    phone: string,
    isphoneverified: boolean,
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
  }
}

export class CreateOrgRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDomain(): string;
  setDomain(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateOrgRequest): CreateOrgRequest.AsObject;
  static serializeBinaryToWriter(message: CreateOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateOrgRequest;
  static deserializeBinaryFromReader(message: CreateOrgRequest, reader: jspb.BinaryReader): CreateOrgRequest;
}

export namespace CreateOrgRequest {
  export type AsObject = {
    name: string,
    domain: string,
  }
}

export enum OrgState { 
  NONE = 0,
  ACTIVE = 1,
  INACTIVE = 2,
}
export enum OrgSearchKey { 
  UNKNOWN = 0,
  ORG_NAME = 1,
  DOMAIN = 2,
  STATE = 3,
}
export enum OrgSearchMethod { 
  EQUALS = 0,
  STARTS_WITH = 1,
  CONTAINS = 2,
}
export enum UserState { 
  USERSTATE_NONE = 0,
  USERSTATE_ACTIVE = 1,
  USERSTATE_INACTIVE = 2,
  USERSTATE_DELETED = 3,
  USERSTATE_LOCKED = 4,
  USERSTATE_SUSPEND = 5,
  USERSTATE_INITIAL = 6,
}
export enum Gender { 
  UNKNOWN_GENDER = 0,
  FEMALE = 1,
  MALE = 2,
  DIVERSE = 3,
}

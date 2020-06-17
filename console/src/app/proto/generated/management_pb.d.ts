import * as jspb from "google-protobuf"

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';
import * as authoption_options_pb from './authoption/options_pb';

export class Iam extends jspb.Message {
  getGlobalOrgId(): string;
  setGlobalOrgId(value: string): void;

  getIamProjectId(): string;
  setIamProjectId(value: string): void;

  getSetUpDone(): boolean;
  setSetUpDone(value: boolean): void;

  getSetUpStarted(): boolean;
  setSetUpStarted(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Iam.AsObject;
  static toObject(includeInstance: boolean, msg: Iam): Iam.AsObject;
  static serializeBinaryToWriter(message: Iam, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Iam;
  static deserializeBinaryFromReader(message: Iam, reader: jspb.BinaryReader): Iam;
}

export namespace Iam {
  export type AsObject = {
    globalOrgId: string,
    iamProjectId: string,
    setUpDone: boolean,
    setUpStarted: boolean,
  }
}

export class ChangeRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getSecId(): string;
  setSecId(value: string): void;

  getLimit(): number;
  setLimit(value: number): void;

  getSequenceOffset(): number;
  setSequenceOffset(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeRequest): ChangeRequest.AsObject;
  static serializeBinaryToWriter(message: ChangeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeRequest;
  static deserializeBinaryFromReader(message: ChangeRequest, reader: jspb.BinaryReader): ChangeRequest;
}

export namespace ChangeRequest {
  export type AsObject = {
    id: string,
    secId: string,
    limit: number,
    sequenceOffset: number,
  }
}

export class Changes extends jspb.Message {
  getChangesList(): Array<Change>;
  setChangesList(value: Array<Change>): void;
  clearChangesList(): void;
  addChanges(value?: Change, index?: number): Change;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Changes.AsObject;
  static toObject(includeInstance: boolean, msg: Changes): Changes.AsObject;
  static serializeBinaryToWriter(message: Changes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Changes;
  static deserializeBinaryFromReader(message: Changes, reader: jspb.BinaryReader): Changes;
}

export namespace Changes {
  export type AsObject = {
    changesList: Array<Change.AsObject>,
    offset: number,
    limit: number,
  }
}

export class Change extends jspb.Message {
  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getEventType(): string;
  setEventType(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  getEditor(): string;
  setEditor(value: string): void;

  getData(): google_protobuf_struct_pb.Struct | undefined;
  setData(value?: google_protobuf_struct_pb.Struct): void;
  hasData(): boolean;
  clearData(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Change.AsObject;
  static toObject(includeInstance: boolean, msg: Change): Change.AsObject;
  static serializeBinaryToWriter(message: Change, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Change;
  static deserializeBinaryFromReader(message: Change, reader: jspb.BinaryReader): Change;
}

export namespace Change {
  export type AsObject = {
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    eventType: string,
    sequence: number,
    editor: string,
    data?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class ApplicationID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationID.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationID): ApplicationID.AsObject;
  static serializeBinaryToWriter(message: ApplicationID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationID;
  static deserializeBinaryFromReader(message: ApplicationID, reader: jspb.BinaryReader): ApplicationID;
}

export namespace ApplicationID {
  export type AsObject = {
    id: string,
    projectId: string,
  }
}

export class ProjectID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectID.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectID): ProjectID.AsObject;
  static serializeBinaryToWriter(message: ProjectID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectID;
  static deserializeBinaryFromReader(message: ProjectID, reader: jspb.BinaryReader): ProjectID;
}

export namespace ProjectID {
  export type AsObject = {
    id: string,
  }
}

export class UserID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserID.AsObject;
  static toObject(includeInstance: boolean, msg: UserID): UserID.AsObject;
  static serializeBinaryToWriter(message: UserID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserID;
  static deserializeBinaryFromReader(message: UserID, reader: jspb.BinaryReader): UserID;
}

export namespace UserID {
  export type AsObject = {
    id: string,
  }
}

export class UserEmailID extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserEmailID.AsObject;
  static toObject(includeInstance: boolean, msg: UserEmailID): UserEmailID.AsObject;
  static serializeBinaryToWriter(message: UserEmailID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserEmailID;
  static deserializeBinaryFromReader(message: UserEmailID, reader: jspb.BinaryReader): UserEmailID;
}

export namespace UserEmailID {
  export type AsObject = {
    email: string,
  }
}

export class UniqueUserRequest extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UniqueUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UniqueUserRequest): UniqueUserRequest.AsObject;
  static serializeBinaryToWriter(message: UniqueUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UniqueUserRequest;
  static deserializeBinaryFromReader(message: UniqueUserRequest, reader: jspb.BinaryReader): UniqueUserRequest;
}

export namespace UniqueUserRequest {
  export type AsObject = {
    userName: string,
    email: string,
  }
}

export class UniqueUserResponse extends jspb.Message {
  getIsUnique(): boolean;
  setIsUnique(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UniqueUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UniqueUserResponse): UniqueUserResponse.AsObject;
  static serializeBinaryToWriter(message: UniqueUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UniqueUserResponse;
  static deserializeBinaryFromReader(message: UniqueUserResponse, reader: jspb.BinaryReader): UniqueUserResponse;
}

export namespace UniqueUserResponse {
  export type AsObject = {
    isUnique: boolean,
  }
}

export class CreateUserRequest extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getNickName(): string;
  setNickName(value: string): void;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): void;

  getGender(): Gender;
  setGender(value: Gender): void;

  getEmail(): string;
  setEmail(value: string): void;

  getIsEmailVerified(): boolean;
  setIsEmailVerified(value: boolean): void;

  getPhone(): string;
  setPhone(value: string): void;

  getIsPhoneVerified(): boolean;
  setIsPhoneVerified(value: boolean): void;

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

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUserRequest): CreateUserRequest.AsObject;
  static serializeBinaryToWriter(message: CreateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUserRequest;
  static deserializeBinaryFromReader(message: CreateUserRequest, reader: jspb.BinaryReader): CreateUserRequest;
}

export namespace CreateUserRequest {
  export type AsObject = {
    userName: string,
    firstName: string,
    lastName: string,
    nickName: string,
    preferredLanguage: string,
    gender: Gender,
    email: string,
    isEmailVerified: boolean,
    phone: string,
    isPhoneVerified: boolean,
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
    password: string,
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

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getNickName(): string;
  setNickName(value: string): void;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): void;

  getGender(): Gender;
  setGender(value: Gender): void;

  getEmail(): string;
  setEmail(value: string): void;

  getIsEmailVerified(): boolean;
  setIsEmailVerified(value: boolean): void;

  getPhone(): string;
  setPhone(value: string): void;

  getIsPhoneVerified(): boolean;
  setIsPhoneVerified(value: boolean): void;

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

  getSequence(): number;
  setSequence(value: number): void;

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
    displayName: string,
    nickName: string,
    preferredLanguage: string,
    gender: Gender,
    email: string,
    isEmailVerified: boolean,
    phone: string,
    isPhoneVerified: boolean,
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
    sequence: number,
  }
}

export class UserView extends jspb.Message {
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

  getLastLogin(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastLogin(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasLastLogin(): boolean;
  clearLastLogin(): void;

  getPasswordChanged(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setPasswordChanged(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasPasswordChanged(): boolean;
  clearPasswordChanged(): void;

  getUserName(): string;
  setUserName(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getNickName(): string;
  setNickName(value: string): void;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): void;

  getGender(): Gender;
  setGender(value: Gender): void;

  getEmail(): string;
  setEmail(value: string): void;

  getIsEmailVerified(): boolean;
  setIsEmailVerified(value: boolean): void;

  getPhone(): string;
  setPhone(value: string): void;

  getIsPhoneVerified(): boolean;
  setIsPhoneVerified(value: boolean): void;

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

  getSequence(): number;
  setSequence(value: number): void;

  getResourceOwner(): string;
  setResourceOwner(value: string): void;

  getLoginNamesList(): Array<string>;
  setLoginNamesList(value: Array<string>): void;
  clearLoginNamesList(): void;
  addLoginNames(value: string, index?: number): void;

  getPreferredLoginName(): string;
  setPreferredLoginName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserView.AsObject;
  static toObject(includeInstance: boolean, msg: UserView): UserView.AsObject;
  static serializeBinaryToWriter(message: UserView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserView;
  static deserializeBinaryFromReader(message: UserView, reader: jspb.BinaryReader): UserView;
}

export namespace UserView {
  export type AsObject = {
    id: string,
    state: UserState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastLogin?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    passwordChanged?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    userName: string,
    firstName: string,
    lastName: string,
    displayName: string,
    nickName: string,
    preferredLanguage: string,
    gender: Gender,
    email: string,
    isEmailVerified: boolean,
    phone: string,
    isPhoneVerified: boolean,
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
    sequence: number,
    resourceOwner: string,
    loginNamesList: Array<string>,
    preferredLoginName: string,
  }
}

export class UserSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getSortingColumn(): UserSearchKey;
  setSortingColumn(value: UserSearchKey): void;

  getAsc(): boolean;
  setAsc(value: boolean): void;

  getQueriesList(): Array<UserSearchQuery>;
  setQueriesList(value: Array<UserSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: UserSearchQuery, index?: number): UserSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UserSearchRequest): UserSearchRequest.AsObject;
  static serializeBinaryToWriter(message: UserSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSearchRequest;
  static deserializeBinaryFromReader(message: UserSearchRequest, reader: jspb.BinaryReader): UserSearchRequest;
}

export namespace UserSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    sortingColumn: UserSearchKey,
    asc: boolean,
    queriesList: Array<UserSearchQuery.AsObject>,
  }
}

export class UserSearchQuery extends jspb.Message {
  getKey(): UserSearchKey;
  setKey(value: UserSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserSearchQuery): UserSearchQuery.AsObject;
  static serializeBinaryToWriter(message: UserSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSearchQuery;
  static deserializeBinaryFromReader(message: UserSearchQuery, reader: jspb.BinaryReader): UserSearchQuery;
}

export namespace UserSearchQuery {
  export type AsObject = {
    key: UserSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class UserSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<UserView>;
  setResultList(value: Array<UserView>): void;
  clearResultList(): void;
  addResult(value?: UserView, index?: number): UserView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UserSearchResponse): UserSearchResponse.AsObject;
  static serializeBinaryToWriter(message: UserSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSearchResponse;
  static deserializeBinaryFromReader(message: UserSearchResponse, reader: jspb.BinaryReader): UserSearchResponse;
}

export namespace UserSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<UserView.AsObject>,
  }
}

export class UserProfile extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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

  getUserName(): string;
  setUserName(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserProfile.AsObject;
  static toObject(includeInstance: boolean, msg: UserProfile): UserProfile.AsObject;
  static serializeBinaryToWriter(message: UserProfile, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserProfile;
  static deserializeBinaryFromReader(message: UserProfile, reader: jspb.BinaryReader): UserProfile;
}

export namespace UserProfile {
  export type AsObject = {
    id: string,
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: Gender,
    userName: string,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UserProfileView extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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

  getUserName(): string;
  setUserName(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getLoginNamesList(): Array<string>;
  setLoginNamesList(value: Array<string>): void;
  clearLoginNamesList(): void;
  addLoginNames(value: string, index?: number): void;

  getPreferredLoginName(): string;
  setPreferredLoginName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserProfileView.AsObject;
  static toObject(includeInstance: boolean, msg: UserProfileView): UserProfileView.AsObject;
  static serializeBinaryToWriter(message: UserProfileView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserProfileView;
  static deserializeBinaryFromReader(message: UserProfileView, reader: jspb.BinaryReader): UserProfileView;
}

export namespace UserProfileView {
  export type AsObject = {
    id: string,
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: Gender,
    userName: string,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    loginNamesList: Array<string>,
    preferredLoginName: string,
  }
}

export class UpdateUserProfileRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getNickName(): string;
  setNickName(value: string): void;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): void;

  getGender(): Gender;
  setGender(value: Gender): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserProfileRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserProfileRequest): UpdateUserProfileRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserProfileRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserProfileRequest;
  static deserializeBinaryFromReader(message: UpdateUserProfileRequest, reader: jspb.BinaryReader): UpdateUserProfileRequest;
}

export namespace UpdateUserProfileRequest {
  export type AsObject = {
    id: string,
    firstName: string,
    lastName: string,
    nickName: string,
    preferredLanguage: string,
    gender: Gender,
  }
}

export class UserEmail extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getIsEmailVerified(): boolean;
  setIsEmailVerified(value: boolean): void;

  getSequence(): number;
  setSequence(value: number): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserEmail.AsObject;
  static toObject(includeInstance: boolean, msg: UserEmail): UserEmail.AsObject;
  static serializeBinaryToWriter(message: UserEmail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserEmail;
  static deserializeBinaryFromReader(message: UserEmail, reader: jspb.BinaryReader): UserEmail;
}

export namespace UserEmail {
  export type AsObject = {
    id: string,
    email: string,
    isEmailVerified: boolean,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UserEmailView extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getIsEmailVerified(): boolean;
  setIsEmailVerified(value: boolean): void;

  getSequence(): number;
  setSequence(value: number): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserEmailView.AsObject;
  static toObject(includeInstance: boolean, msg: UserEmailView): UserEmailView.AsObject;
  static serializeBinaryToWriter(message: UserEmailView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserEmailView;
  static deserializeBinaryFromReader(message: UserEmailView, reader: jspb.BinaryReader): UserEmailView;
}

export namespace UserEmailView {
  export type AsObject = {
    id: string,
    email: string,
    isEmailVerified: boolean,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UpdateUserEmailRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getIsEmailVerified(): boolean;
  setIsEmailVerified(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserEmailRequest): UpdateUserEmailRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserEmailRequest;
  static deserializeBinaryFromReader(message: UpdateUserEmailRequest, reader: jspb.BinaryReader): UpdateUserEmailRequest;
}

export namespace UpdateUserEmailRequest {
  export type AsObject = {
    id: string,
    email: string,
    isEmailVerified: boolean,
  }
}

export class UserPhone extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getPhone(): string;
  setPhone(value: string): void;

  getIsPhoneVerified(): boolean;
  setIsPhoneVerified(value: boolean): void;

  getSequence(): number;
  setSequence(value: number): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserPhone.AsObject;
  static toObject(includeInstance: boolean, msg: UserPhone): UserPhone.AsObject;
  static serializeBinaryToWriter(message: UserPhone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserPhone;
  static deserializeBinaryFromReader(message: UserPhone, reader: jspb.BinaryReader): UserPhone;
}

export namespace UserPhone {
  export type AsObject = {
    id: string,
    phone: string,
    isPhoneVerified: boolean,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UserPhoneView extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getPhone(): string;
  setPhone(value: string): void;

  getIsPhoneVerified(): boolean;
  setIsPhoneVerified(value: boolean): void;

  getSequence(): number;
  setSequence(value: number): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserPhoneView.AsObject;
  static toObject(includeInstance: boolean, msg: UserPhoneView): UserPhoneView.AsObject;
  static serializeBinaryToWriter(message: UserPhoneView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserPhoneView;
  static deserializeBinaryFromReader(message: UserPhoneView, reader: jspb.BinaryReader): UserPhoneView;
}

export namespace UserPhoneView {
  export type AsObject = {
    id: string,
    phone: string,
    isPhoneVerified: boolean,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UpdateUserPhoneRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getPhone(): string;
  setPhone(value: string): void;

  getIsPhoneVerified(): boolean;
  setIsPhoneVerified(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserPhoneRequest): UpdateUserPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserPhoneRequest;
  static deserializeBinaryFromReader(message: UpdateUserPhoneRequest, reader: jspb.BinaryReader): UpdateUserPhoneRequest;
}

export namespace UpdateUserPhoneRequest {
  export type AsObject = {
    id: string,
    phone: string,
    isPhoneVerified: boolean,
  }
}

export class UserAddress extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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

  getSequence(): number;
  setSequence(value: number): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAddress.AsObject;
  static toObject(includeInstance: boolean, msg: UserAddress): UserAddress.AsObject;
  static serializeBinaryToWriter(message: UserAddress, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserAddress;
  static deserializeBinaryFromReader(message: UserAddress, reader: jspb.BinaryReader): UserAddress;
}

export namespace UserAddress {
  export type AsObject = {
    id: string,
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UserAddressView extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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

  getSequence(): number;
  setSequence(value: number): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAddressView.AsObject;
  static toObject(includeInstance: boolean, msg: UserAddressView): UserAddressView.AsObject;
  static serializeBinaryToWriter(message: UserAddressView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserAddressView;
  static deserializeBinaryFromReader(message: UserAddressView, reader: jspb.BinaryReader): UserAddressView;
}

export namespace UserAddressView {
  export type AsObject = {
    id: string,
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UpdateUserAddressRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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
  toObject(includeInstance?: boolean): UpdateUserAddressRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserAddressRequest): UpdateUserAddressRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserAddressRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserAddressRequest;
  static deserializeBinaryFromReader(message: UpdateUserAddressRequest, reader: jspb.BinaryReader): UpdateUserAddressRequest;
}

export namespace UpdateUserAddressRequest {
  export type AsObject = {
    id: string,
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
  }
}

export class MultiFactors extends jspb.Message {
  getMfasList(): Array<MultiFactor>;
  setMfasList(value: Array<MultiFactor>): void;
  clearMfasList(): void;
  addMfas(value?: MultiFactor, index?: number): MultiFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MultiFactors.AsObject;
  static toObject(includeInstance: boolean, msg: MultiFactors): MultiFactors.AsObject;
  static serializeBinaryToWriter(message: MultiFactors, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MultiFactors;
  static deserializeBinaryFromReader(message: MultiFactors, reader: jspb.BinaryReader): MultiFactors;
}

export namespace MultiFactors {
  export type AsObject = {
    mfasList: Array<MultiFactor.AsObject>,
  }
}

export class MultiFactor extends jspb.Message {
  getType(): MfaType;
  setType(value: MfaType): void;

  getState(): MFAState;
  setState(value: MFAState): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MultiFactor.AsObject;
  static toObject(includeInstance: boolean, msg: MultiFactor): MultiFactor.AsObject;
  static serializeBinaryToWriter(message: MultiFactor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MultiFactor;
  static deserializeBinaryFromReader(message: MultiFactor, reader: jspb.BinaryReader): MultiFactor;
}

export namespace MultiFactor {
  export type AsObject = {
    type: MfaType,
    state: MFAState,
  }
}

export class PasswordID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordID.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordID): PasswordID.AsObject;
  static serializeBinaryToWriter(message: PasswordID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordID;
  static deserializeBinaryFromReader(message: PasswordID, reader: jspb.BinaryReader): PasswordID;
}

export namespace PasswordID {
  export type AsObject = {
    id: string,
  }
}

export class PasswordRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordRequest): PasswordRequest.AsObject;
  static serializeBinaryToWriter(message: PasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordRequest;
  static deserializeBinaryFromReader(message: PasswordRequest, reader: jspb.BinaryReader): PasswordRequest;
}

export namespace PasswordRequest {
  export type AsObject = {
    id: string,
    password: string,
  }
}

export class ResetPasswordRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPasswordRequest): ResetPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: ResetPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordRequest;
  static deserializeBinaryFromReader(message: ResetPasswordRequest, reader: jspb.BinaryReader): ResetPasswordRequest;
}

export namespace ResetPasswordRequest {
  export type AsObject = {
    id: string,
  }
}

export class SetPasswordNotificationRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getType(): NotificationType;
  setType(value: NotificationType): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPasswordNotificationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetPasswordNotificationRequest): SetPasswordNotificationRequest.AsObject;
  static serializeBinaryToWriter(message: SetPasswordNotificationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPasswordNotificationRequest;
  static deserializeBinaryFromReader(message: SetPasswordNotificationRequest, reader: jspb.BinaryReader): SetPasswordNotificationRequest;
}

export namespace SetPasswordNotificationRequest {
  export type AsObject = {
    id: string,
    type: NotificationType,
  }
}

export class PasswordComplexityPolicyID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordComplexityPolicyID.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordComplexityPolicyID): PasswordComplexityPolicyID.AsObject;
  static serializeBinaryToWriter(message: PasswordComplexityPolicyID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordComplexityPolicyID;
  static deserializeBinaryFromReader(message: PasswordComplexityPolicyID, reader: jspb.BinaryReader): PasswordComplexityPolicyID;
}

export namespace PasswordComplexityPolicyID {
  export type AsObject = {
    id: string,
  }
}

export class PasswordComplexityPolicy extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getState(): PolicyState;
  setState(value: PolicyState): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getMinLength(): number;
  setMinLength(value: number): void;

  getHasLowercase(): boolean;
  setHasLowercase(value: boolean): void;

  getHasUppercase(): boolean;
  setHasUppercase(value: boolean): void;

  getHasNumber(): boolean;
  setHasNumber(value: boolean): void;

  getHasSymbol(): boolean;
  setHasSymbol(value: boolean): void;

  getSequence(): number;
  setSequence(value: number): void;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordComplexityPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordComplexityPolicy): PasswordComplexityPolicy.AsObject;
  static serializeBinaryToWriter(message: PasswordComplexityPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordComplexityPolicy;
  static deserializeBinaryFromReader(message: PasswordComplexityPolicy, reader: jspb.BinaryReader): PasswordComplexityPolicy;
}

export namespace PasswordComplexityPolicy {
  export type AsObject = {
    id: string,
    description: string,
    state: PolicyState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    minLength: number,
    hasLowercase: boolean,
    hasUppercase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
    sequence: number,
    isDefault: boolean,
  }
}

export class PasswordComplexityPolicyCreate extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  getMinLength(): number;
  setMinLength(value: number): void;

  getHasLowercase(): boolean;
  setHasLowercase(value: boolean): void;

  getHasUppercase(): boolean;
  setHasUppercase(value: boolean): void;

  getHasNumber(): boolean;
  setHasNumber(value: boolean): void;

  getHasSymbol(): boolean;
  setHasSymbol(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordComplexityPolicyCreate.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordComplexityPolicyCreate): PasswordComplexityPolicyCreate.AsObject;
  static serializeBinaryToWriter(message: PasswordComplexityPolicyCreate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordComplexityPolicyCreate;
  static deserializeBinaryFromReader(message: PasswordComplexityPolicyCreate, reader: jspb.BinaryReader): PasswordComplexityPolicyCreate;
}

export namespace PasswordComplexityPolicyCreate {
  export type AsObject = {
    description: string,
    minLength: number,
    hasLowercase: boolean,
    hasUppercase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
  }
}

export class PasswordComplexityPolicyUpdate extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getMinLength(): number;
  setMinLength(value: number): void;

  getHasLowercase(): boolean;
  setHasLowercase(value: boolean): void;

  getHasUppercase(): boolean;
  setHasUppercase(value: boolean): void;

  getHasNumber(): boolean;
  setHasNumber(value: boolean): void;

  getHasSymbol(): boolean;
  setHasSymbol(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordComplexityPolicyUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordComplexityPolicyUpdate): PasswordComplexityPolicyUpdate.AsObject;
  static serializeBinaryToWriter(message: PasswordComplexityPolicyUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordComplexityPolicyUpdate;
  static deserializeBinaryFromReader(message: PasswordComplexityPolicyUpdate, reader: jspb.BinaryReader): PasswordComplexityPolicyUpdate;
}

export namespace PasswordComplexityPolicyUpdate {
  export type AsObject = {
    id: string,
    description: string,
    minLength: number,
    hasLowercase: boolean,
    hasUppercase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
  }
}

export class PasswordAgePolicyID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordAgePolicyID.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordAgePolicyID): PasswordAgePolicyID.AsObject;
  static serializeBinaryToWriter(message: PasswordAgePolicyID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordAgePolicyID;
  static deserializeBinaryFromReader(message: PasswordAgePolicyID, reader: jspb.BinaryReader): PasswordAgePolicyID;
}

export namespace PasswordAgePolicyID {
  export type AsObject = {
    id: string,
  }
}

export class PasswordAgePolicy extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getState(): PolicyState;
  setState(value: PolicyState): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getMaxAgeDays(): number;
  setMaxAgeDays(value: number): void;

  getExpireWarnDays(): number;
  setExpireWarnDays(value: number): void;

  getSequence(): number;
  setSequence(value: number): void;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordAgePolicy.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordAgePolicy): PasswordAgePolicy.AsObject;
  static serializeBinaryToWriter(message: PasswordAgePolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordAgePolicy;
  static deserializeBinaryFromReader(message: PasswordAgePolicy, reader: jspb.BinaryReader): PasswordAgePolicy;
}

export namespace PasswordAgePolicy {
  export type AsObject = {
    id: string,
    description: string,
    state: PolicyState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    maxAgeDays: number,
    expireWarnDays: number,
    sequence: number,
    isDefault: boolean,
  }
}

export class PasswordAgePolicyCreate extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  getMaxAgeDays(): number;
  setMaxAgeDays(value: number): void;

  getExpireWarnDays(): number;
  setExpireWarnDays(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordAgePolicyCreate.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordAgePolicyCreate): PasswordAgePolicyCreate.AsObject;
  static serializeBinaryToWriter(message: PasswordAgePolicyCreate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordAgePolicyCreate;
  static deserializeBinaryFromReader(message: PasswordAgePolicyCreate, reader: jspb.BinaryReader): PasswordAgePolicyCreate;
}

export namespace PasswordAgePolicyCreate {
  export type AsObject = {
    description: string,
    maxAgeDays: number,
    expireWarnDays: number,
  }
}

export class PasswordAgePolicyUpdate extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getMaxAgeDays(): number;
  setMaxAgeDays(value: number): void;

  getExpireWarnDays(): number;
  setExpireWarnDays(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordAgePolicyUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordAgePolicyUpdate): PasswordAgePolicyUpdate.AsObject;
  static serializeBinaryToWriter(message: PasswordAgePolicyUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordAgePolicyUpdate;
  static deserializeBinaryFromReader(message: PasswordAgePolicyUpdate, reader: jspb.BinaryReader): PasswordAgePolicyUpdate;
}

export namespace PasswordAgePolicyUpdate {
  export type AsObject = {
    id: string,
    description: string,
    maxAgeDays: number,
    expireWarnDays: number,
  }
}

export class PasswordLockoutPolicyID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordLockoutPolicyID.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordLockoutPolicyID): PasswordLockoutPolicyID.AsObject;
  static serializeBinaryToWriter(message: PasswordLockoutPolicyID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordLockoutPolicyID;
  static deserializeBinaryFromReader(message: PasswordLockoutPolicyID, reader: jspb.BinaryReader): PasswordLockoutPolicyID;
}

export namespace PasswordLockoutPolicyID {
  export type AsObject = {
    id: string,
  }
}

export class PasswordLockoutPolicy extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getState(): PolicyState;
  setState(value: PolicyState): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getMaxAttempts(): number;
  setMaxAttempts(value: number): void;

  getShowLockOutFailures(): boolean;
  setShowLockOutFailures(value: boolean): void;

  getSequence(): number;
  setSequence(value: number): void;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordLockoutPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordLockoutPolicy): PasswordLockoutPolicy.AsObject;
  static serializeBinaryToWriter(message: PasswordLockoutPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordLockoutPolicy;
  static deserializeBinaryFromReader(message: PasswordLockoutPolicy, reader: jspb.BinaryReader): PasswordLockoutPolicy;
}

export namespace PasswordLockoutPolicy {
  export type AsObject = {
    id: string,
    description: string,
    state: PolicyState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    maxAttempts: number,
    showLockOutFailures: boolean,
    sequence: number,
    isDefault: boolean,
  }
}

export class PasswordLockoutPolicyCreate extends jspb.Message {
  getDescription(): string;
  setDescription(value: string): void;

  getMaxAttempts(): number;
  setMaxAttempts(value: number): void;

  getShowLockOutFailures(): boolean;
  setShowLockOutFailures(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordLockoutPolicyCreate.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordLockoutPolicyCreate): PasswordLockoutPolicyCreate.AsObject;
  static serializeBinaryToWriter(message: PasswordLockoutPolicyCreate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordLockoutPolicyCreate;
  static deserializeBinaryFromReader(message: PasswordLockoutPolicyCreate, reader: jspb.BinaryReader): PasswordLockoutPolicyCreate;
}

export namespace PasswordLockoutPolicyCreate {
  export type AsObject = {
    description: string,
    maxAttempts: number,
    showLockOutFailures: boolean,
  }
}

export class PasswordLockoutPolicyUpdate extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getMaxAttempts(): number;
  setMaxAttempts(value: number): void;

  getShowLockOutFailures(): boolean;
  setShowLockOutFailures(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordLockoutPolicyUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordLockoutPolicyUpdate): PasswordLockoutPolicyUpdate.AsObject;
  static serializeBinaryToWriter(message: PasswordLockoutPolicyUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordLockoutPolicyUpdate;
  static deserializeBinaryFromReader(message: PasswordLockoutPolicyUpdate, reader: jspb.BinaryReader): PasswordLockoutPolicyUpdate;
}

export namespace PasswordLockoutPolicyUpdate {
  export type AsObject = {
    id: string,
    description: string,
    maxAttempts: number,
    showLockOutFailures: boolean,
  }
}

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

  getSequence(): number;
  setSequence(value: number): void;

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
    sequence: number,
  }
}

export class OrgDomains extends jspb.Message {
  getDomainsList(): Array<OrgDomain>;
  setDomainsList(value: Array<OrgDomain>): void;
  clearDomainsList(): void;
  addDomains(value?: OrgDomain, index?: number): OrgDomain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgDomains.AsObject;
  static toObject(includeInstance: boolean, msg: OrgDomains): OrgDomains.AsObject;
  static serializeBinaryToWriter(message: OrgDomains, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgDomains;
  static deserializeBinaryFromReader(message: OrgDomains, reader: jspb.BinaryReader): OrgDomains;
}

export namespace OrgDomains {
  export type AsObject = {
    domainsList: Array<OrgDomain.AsObject>,
  }
}

export class OrgDomain extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getDomain(): string;
  setDomain(value: string): void;

  getVerified(): boolean;
  setVerified(value: boolean): void;

  getPrimary(): boolean;
  setPrimary(value: boolean): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgDomain.AsObject;
  static toObject(includeInstance: boolean, msg: OrgDomain): OrgDomain.AsObject;
  static serializeBinaryToWriter(message: OrgDomain, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgDomain;
  static deserializeBinaryFromReader(message: OrgDomain, reader: jspb.BinaryReader): OrgDomain;
}

export namespace OrgDomain {
  export type AsObject = {
    orgId: string,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    domain: string,
    verified: boolean,
    primary: boolean,
    sequence: number,
  }
}

export class OrgDomainView extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getDomain(): string;
  setDomain(value: string): void;

  getVerified(): boolean;
  setVerified(value: boolean): void;

  getPrimary(): boolean;
  setPrimary(value: boolean): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgDomainView.AsObject;
  static toObject(includeInstance: boolean, msg: OrgDomainView): OrgDomainView.AsObject;
  static serializeBinaryToWriter(message: OrgDomainView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgDomainView;
  static deserializeBinaryFromReader(message: OrgDomainView, reader: jspb.BinaryReader): OrgDomainView;
}

export namespace OrgDomainView {
  export type AsObject = {
    orgId: string,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    domain: string,
    verified: boolean,
    primary: boolean,
    sequence: number,
  }
}

export class AddOrgDomainRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgDomainRequest): AddOrgDomainRequest.AsObject;
  static serializeBinaryToWriter(message: AddOrgDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgDomainRequest;
  static deserializeBinaryFromReader(message: AddOrgDomainRequest, reader: jspb.BinaryReader): AddOrgDomainRequest;
}

export namespace AddOrgDomainRequest {
  export type AsObject = {
    domain: string,
  }
}

export class RemoveOrgDomainRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgDomainRequest): RemoveOrgDomainRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgDomainRequest;
  static deserializeBinaryFromReader(message: RemoveOrgDomainRequest, reader: jspb.BinaryReader): RemoveOrgDomainRequest;
}

export namespace RemoveOrgDomainRequest {
  export type AsObject = {
    domain: string,
  }
}

export class OrgDomainSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<OrgDomainView>;
  setResultList(value: Array<OrgDomainView>): void;
  clearResultList(): void;
  addResult(value?: OrgDomainView, index?: number): OrgDomainView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgDomainSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: OrgDomainSearchResponse): OrgDomainSearchResponse.AsObject;
  static serializeBinaryToWriter(message: OrgDomainSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgDomainSearchResponse;
  static deserializeBinaryFromReader(message: OrgDomainSearchResponse, reader: jspb.BinaryReader): OrgDomainSearchResponse;
}

export namespace OrgDomainSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<OrgDomainView.AsObject>,
  }
}

export class OrgDomainSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<OrgDomainSearchQuery>;
  setQueriesList(value: Array<OrgDomainSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: OrgDomainSearchQuery, index?: number): OrgDomainSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgDomainSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: OrgDomainSearchRequest): OrgDomainSearchRequest.AsObject;
  static serializeBinaryToWriter(message: OrgDomainSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgDomainSearchRequest;
  static deserializeBinaryFromReader(message: OrgDomainSearchRequest, reader: jspb.BinaryReader): OrgDomainSearchRequest;
}

export namespace OrgDomainSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    queriesList: Array<OrgDomainSearchQuery.AsObject>,
  }
}

export class OrgDomainSearchQuery extends jspb.Message {
  getKey(): OrgDomainSearchKey;
  setKey(value: OrgDomainSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgDomainSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrgDomainSearchQuery): OrgDomainSearchQuery.AsObject;
  static serializeBinaryToWriter(message: OrgDomainSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgDomainSearchQuery;
  static deserializeBinaryFromReader(message: OrgDomainSearchQuery, reader: jspb.BinaryReader): OrgDomainSearchQuery;
}

export namespace OrgDomainSearchQuery {
  export type AsObject = {
    key: OrgDomainSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class OrgMemberRoles extends jspb.Message {
  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgMemberRoles.AsObject;
  static toObject(includeInstance: boolean, msg: OrgMemberRoles): OrgMemberRoles.AsObject;
  static serializeBinaryToWriter(message: OrgMemberRoles, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgMemberRoles;
  static deserializeBinaryFromReader(message: OrgMemberRoles, reader: jspb.BinaryReader): OrgMemberRoles;
}

export namespace OrgMemberRoles {
  export type AsObject = {
    rolesList: Array<string>,
  }
}

export class OrgMember extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgMember.AsObject;
  static toObject(includeInstance: boolean, msg: OrgMember): OrgMember.AsObject;
  static serializeBinaryToWriter(message: OrgMember, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgMember;
  static deserializeBinaryFromReader(message: OrgMember, reader: jspb.BinaryReader): OrgMember;
}

export namespace OrgMember {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class AddOrgMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgMemberRequest): AddOrgMemberRequest.AsObject;
  static serializeBinaryToWriter(message: AddOrgMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgMemberRequest;
  static deserializeBinaryFromReader(message: AddOrgMemberRequest, reader: jspb.BinaryReader): AddOrgMemberRequest;
}

export namespace AddOrgMemberRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
  }
}

export class ChangeOrgMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeOrgMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeOrgMemberRequest): ChangeOrgMemberRequest.AsObject;
  static serializeBinaryToWriter(message: ChangeOrgMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeOrgMemberRequest;
  static deserializeBinaryFromReader(message: ChangeOrgMemberRequest, reader: jspb.BinaryReader): ChangeOrgMemberRequest;
}

export namespace ChangeOrgMemberRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
  }
}

export class RemoveOrgMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgMemberRequest): RemoveOrgMemberRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgMemberRequest;
  static deserializeBinaryFromReader(message: RemoveOrgMemberRequest, reader: jspb.BinaryReader): RemoveOrgMemberRequest;
}

export namespace RemoveOrgMemberRequest {
  export type AsObject = {
    userId: string,
  }
}

export class OrgMemberSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<OrgMemberView>;
  setResultList(value: Array<OrgMemberView>): void;
  clearResultList(): void;
  addResult(value?: OrgMemberView, index?: number): OrgMemberView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgMemberSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: OrgMemberSearchResponse): OrgMemberSearchResponse.AsObject;
  static serializeBinaryToWriter(message: OrgMemberSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgMemberSearchResponse;
  static deserializeBinaryFromReader(message: OrgMemberSearchResponse, reader: jspb.BinaryReader): OrgMemberSearchResponse;
}

export namespace OrgMemberSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<OrgMemberView.AsObject>,
  }
}

export class OrgMemberView extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  getUserName(): string;
  setUserName(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgMemberView.AsObject;
  static toObject(includeInstance: boolean, msg: OrgMemberView): OrgMemberView.AsObject;
  static serializeBinaryToWriter(message: OrgMemberView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgMemberView;
  static deserializeBinaryFromReader(message: OrgMemberView, reader: jspb.BinaryReader): OrgMemberView;
}

export namespace OrgMemberView {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
    userName: string,
    email: string,
    firstName: string,
    lastName: string,
  }
}

export class OrgMemberSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<OrgMemberSearchQuery>;
  setQueriesList(value: Array<OrgMemberSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: OrgMemberSearchQuery, index?: number): OrgMemberSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgMemberSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: OrgMemberSearchRequest): OrgMemberSearchRequest.AsObject;
  static serializeBinaryToWriter(message: OrgMemberSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgMemberSearchRequest;
  static deserializeBinaryFromReader(message: OrgMemberSearchRequest, reader: jspb.BinaryReader): OrgMemberSearchRequest;
}

export namespace OrgMemberSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    queriesList: Array<OrgMemberSearchQuery.AsObject>,
  }
}

export class OrgMemberSearchQuery extends jspb.Message {
  getKey(): OrgMemberSearchKey;
  setKey(value: OrgMemberSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgMemberSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrgMemberSearchQuery): OrgMemberSearchQuery.AsObject;
  static serializeBinaryToWriter(message: OrgMemberSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgMemberSearchQuery;
  static deserializeBinaryFromReader(message: OrgMemberSearchQuery, reader: jspb.BinaryReader): OrgMemberSearchQuery;
}

export namespace OrgMemberSearchQuery {
  export type AsObject = {
    key: OrgMemberSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class ProjectCreateRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectCreateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectCreateRequest): ProjectCreateRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectCreateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectCreateRequest;
  static deserializeBinaryFromReader(message: ProjectCreateRequest, reader: jspb.BinaryReader): ProjectCreateRequest;
}

export namespace ProjectCreateRequest {
  export type AsObject = {
    name: string,
  }
}

export class ProjectUpdateRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectUpdateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectUpdateRequest): ProjectUpdateRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectUpdateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectUpdateRequest;
  static deserializeBinaryFromReader(message: ProjectUpdateRequest, reader: jspb.BinaryReader): ProjectUpdateRequest;
}

export namespace ProjectUpdateRequest {
  export type AsObject = {
    id: string,
    name: string,
  }
}

export class ProjectSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<ProjectView>;
  setResultList(value: Array<ProjectView>): void;
  clearResultList(): void;
  addResult(value?: ProjectView, index?: number): ProjectView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectSearchResponse): ProjectSearchResponse.AsObject;
  static serializeBinaryToWriter(message: ProjectSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectSearchResponse;
  static deserializeBinaryFromReader(message: ProjectSearchResponse, reader: jspb.BinaryReader): ProjectSearchResponse;
}

export namespace ProjectSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<ProjectView.AsObject>,
  }
}

export class ProjectView extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getState(): ProjectState;
  setState(value: ProjectState): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getResourceOwner(): string;
  setResourceOwner(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectView.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectView): ProjectView.AsObject;
  static serializeBinaryToWriter(message: ProjectView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectView;
  static deserializeBinaryFromReader(message: ProjectView, reader: jspb.BinaryReader): ProjectView;
}

export namespace ProjectView {
  export type AsObject = {
    projectId: string,
    name: string,
    state: ProjectState,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    resourceOwner: string,
    sequence: number,
  }
}

export class ProjectSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<ProjectSearchQuery>;
  setQueriesList(value: Array<ProjectSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: ProjectSearchQuery, index?: number): ProjectSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectSearchRequest): ProjectSearchRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectSearchRequest;
  static deserializeBinaryFromReader(message: ProjectSearchRequest, reader: jspb.BinaryReader): ProjectSearchRequest;
}

export namespace ProjectSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    queriesList: Array<ProjectSearchQuery.AsObject>,
  }
}

export class ProjectSearchQuery extends jspb.Message {
  getKey(): ProjectSearchKey;
  setKey(value: ProjectSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectSearchQuery): ProjectSearchQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectSearchQuery;
  static deserializeBinaryFromReader(message: ProjectSearchQuery, reader: jspb.BinaryReader): ProjectSearchQuery;
}

export namespace ProjectSearchQuery {
  export type AsObject = {
    key: ProjectSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class Projects extends jspb.Message {
  getProjectsList(): Array<Project>;
  setProjectsList(value: Array<Project>): void;
  clearProjectsList(): void;
  addProjects(value?: Project, index?: number): Project;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Projects.AsObject;
  static toObject(includeInstance: boolean, msg: Projects): Projects.AsObject;
  static serializeBinaryToWriter(message: Projects, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Projects;
  static deserializeBinaryFromReader(message: Projects, reader: jspb.BinaryReader): Projects;
}

export namespace Projects {
  export type AsObject = {
    projectsList: Array<Project.AsObject>,
  }
}

export class Project extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getState(): ProjectState;
  setState(value: ProjectState): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Project.AsObject;
  static toObject(includeInstance: boolean, msg: Project): Project.AsObject;
  static serializeBinaryToWriter(message: Project, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Project;
  static deserializeBinaryFromReader(message: Project, reader: jspb.BinaryReader): Project;
}

export namespace Project {
  export type AsObject = {
    id: string,
    name: string,
    state: ProjectState,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class ProjectMemberRoles extends jspb.Message {
  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMemberRoles.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMemberRoles): ProjectMemberRoles.AsObject;
  static serializeBinaryToWriter(message: ProjectMemberRoles, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMemberRoles;
  static deserializeBinaryFromReader(message: ProjectMemberRoles, reader: jspb.BinaryReader): ProjectMemberRoles;
}

export namespace ProjectMemberRoles {
  export type AsObject = {
    rolesList: Array<string>,
  }
}

export class ProjectMember extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMember.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMember): ProjectMember.AsObject;
  static serializeBinaryToWriter(message: ProjectMember, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMember;
  static deserializeBinaryFromReader(message: ProjectMember, reader: jspb.BinaryReader): ProjectMember;
}

export namespace ProjectMember {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class ProjectMemberAdd extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMemberAdd.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMemberAdd): ProjectMemberAdd.AsObject;
  static serializeBinaryToWriter(message: ProjectMemberAdd, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMemberAdd;
  static deserializeBinaryFromReader(message: ProjectMemberAdd, reader: jspb.BinaryReader): ProjectMemberAdd;
}

export namespace ProjectMemberAdd {
  export type AsObject = {
    id: string,
    userId: string,
    rolesList: Array<string>,
  }
}

export class ProjectMemberChange extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMemberChange.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMemberChange): ProjectMemberChange.AsObject;
  static serializeBinaryToWriter(message: ProjectMemberChange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMemberChange;
  static deserializeBinaryFromReader(message: ProjectMemberChange, reader: jspb.BinaryReader): ProjectMemberChange;
}

export namespace ProjectMemberChange {
  export type AsObject = {
    id: string,
    userId: string,
    rolesList: Array<string>,
  }
}

export class ProjectMemberRemove extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMemberRemove.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMemberRemove): ProjectMemberRemove.AsObject;
  static serializeBinaryToWriter(message: ProjectMemberRemove, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMemberRemove;
  static deserializeBinaryFromReader(message: ProjectMemberRemove, reader: jspb.BinaryReader): ProjectMemberRemove;
}

export namespace ProjectMemberRemove {
  export type AsObject = {
    id: string,
    userId: string,
  }
}

export class ProjectRoleAdd extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getGroup(): string;
  setGroup(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRoleAdd.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRoleAdd): ProjectRoleAdd.AsObject;
  static serializeBinaryToWriter(message: ProjectRoleAdd, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRoleAdd;
  static deserializeBinaryFromReader(message: ProjectRoleAdd, reader: jspb.BinaryReader): ProjectRoleAdd;
}

export namespace ProjectRoleAdd {
  export type AsObject = {
    id: string,
    key: string,
    displayName: string,
    group: string,
  }
}

export class ProjectRoleChange extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getGroup(): string;
  setGroup(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRoleChange.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRoleChange): ProjectRoleChange.AsObject;
  static serializeBinaryToWriter(message: ProjectRoleChange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRoleChange;
  static deserializeBinaryFromReader(message: ProjectRoleChange, reader: jspb.BinaryReader): ProjectRoleChange;
}

export namespace ProjectRoleChange {
  export type AsObject = {
    id: string,
    key: string,
    displayName: string,
    group: string,
  }
}

export class ProjectRole extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getGroup(): string;
  setGroup(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRole.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRole): ProjectRole.AsObject;
  static serializeBinaryToWriter(message: ProjectRole, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRole;
  static deserializeBinaryFromReader(message: ProjectRole, reader: jspb.BinaryReader): ProjectRole;
}

export namespace ProjectRole {
  export type AsObject = {
    projectId: string,
    key: string,
    displayName: string,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    group: string,
    sequence: number,
  }
}

export class ProjectRoleView extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  getDisplayName(): string;
  setDisplayName(value: string): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getGroup(): string;
  setGroup(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRoleView.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRoleView): ProjectRoleView.AsObject;
  static serializeBinaryToWriter(message: ProjectRoleView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRoleView;
  static deserializeBinaryFromReader(message: ProjectRoleView, reader: jspb.BinaryReader): ProjectRoleView;
}

export namespace ProjectRoleView {
  export type AsObject = {
    projectId: string,
    key: string,
    displayName: string,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    group: string,
    sequence: number,
  }
}

export class ProjectRoleRemove extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getKey(): string;
  setKey(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRoleRemove.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRoleRemove): ProjectRoleRemove.AsObject;
  static serializeBinaryToWriter(message: ProjectRoleRemove, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRoleRemove;
  static deserializeBinaryFromReader(message: ProjectRoleRemove, reader: jspb.BinaryReader): ProjectRoleRemove;
}

export namespace ProjectRoleRemove {
  export type AsObject = {
    id: string,
    key: string,
  }
}

export class ProjectRoleSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<ProjectRoleView>;
  setResultList(value: Array<ProjectRoleView>): void;
  clearResultList(): void;
  addResult(value?: ProjectRoleView, index?: number): ProjectRoleView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRoleSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRoleSearchResponse): ProjectRoleSearchResponse.AsObject;
  static serializeBinaryToWriter(message: ProjectRoleSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRoleSearchResponse;
  static deserializeBinaryFromReader(message: ProjectRoleSearchResponse, reader: jspb.BinaryReader): ProjectRoleSearchResponse;
}

export namespace ProjectRoleSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<ProjectRoleView.AsObject>,
  }
}

export class ProjectRoleSearchRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<ProjectRoleSearchQuery>;
  setQueriesList(value: Array<ProjectRoleSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: ProjectRoleSearchQuery, index?: number): ProjectRoleSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRoleSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRoleSearchRequest): ProjectRoleSearchRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectRoleSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRoleSearchRequest;
  static deserializeBinaryFromReader(message: ProjectRoleSearchRequest, reader: jspb.BinaryReader): ProjectRoleSearchRequest;
}

export namespace ProjectRoleSearchRequest {
  export type AsObject = {
    projectId: string,
    offset: number,
    limit: number,
    queriesList: Array<ProjectRoleSearchQuery.AsObject>,
  }
}

export class ProjectRoleSearchQuery extends jspb.Message {
  getKey(): ProjectRoleSearchKey;
  setKey(value: ProjectRoleSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRoleSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRoleSearchQuery): ProjectRoleSearchQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectRoleSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRoleSearchQuery;
  static deserializeBinaryFromReader(message: ProjectRoleSearchQuery, reader: jspb.BinaryReader): ProjectRoleSearchQuery;
}

export namespace ProjectRoleSearchQuery {
  export type AsObject = {
    key: ProjectRoleSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class ProjectMemberView extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getUserName(): string;
  setUserName(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMemberView.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMemberView): ProjectMemberView.AsObject;
  static serializeBinaryToWriter(message: ProjectMemberView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMemberView;
  static deserializeBinaryFromReader(message: ProjectMemberView, reader: jspb.BinaryReader): ProjectMemberView;
}

export namespace ProjectMemberView {
  export type AsObject = {
    userId: string,
    userName: string,
    email: string,
    firstName: string,
    lastName: string,
    rolesList: Array<string>,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class ProjectMemberSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<ProjectMemberView>;
  setResultList(value: Array<ProjectMemberView>): void;
  clearResultList(): void;
  addResult(value?: ProjectMemberView, index?: number): ProjectMemberView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMemberSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMemberSearchResponse): ProjectMemberSearchResponse.AsObject;
  static serializeBinaryToWriter(message: ProjectMemberSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMemberSearchResponse;
  static deserializeBinaryFromReader(message: ProjectMemberSearchResponse, reader: jspb.BinaryReader): ProjectMemberSearchResponse;
}

export namespace ProjectMemberSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<ProjectMemberView.AsObject>,
  }
}

export class ProjectMemberSearchRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<ProjectMemberSearchQuery>;
  setQueriesList(value: Array<ProjectMemberSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: ProjectMemberSearchQuery, index?: number): ProjectMemberSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMemberSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMemberSearchRequest): ProjectMemberSearchRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectMemberSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMemberSearchRequest;
  static deserializeBinaryFromReader(message: ProjectMemberSearchRequest, reader: jspb.BinaryReader): ProjectMemberSearchRequest;
}

export namespace ProjectMemberSearchRequest {
  export type AsObject = {
    projectId: string,
    offset: number,
    limit: number,
    queriesList: Array<ProjectMemberSearchQuery.AsObject>,
  }
}

export class ProjectMemberSearchQuery extends jspb.Message {
  getKey(): ProjectMemberSearchKey;
  setKey(value: ProjectMemberSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectMemberSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectMemberSearchQuery): ProjectMemberSearchQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectMemberSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectMemberSearchQuery;
  static deserializeBinaryFromReader(message: ProjectMemberSearchQuery, reader: jspb.BinaryReader): ProjectMemberSearchQuery;
}

export namespace ProjectMemberSearchQuery {
  export type AsObject = {
    key: ProjectMemberSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class Application extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getState(): AppState;
  setState(value: AppState): void;

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

  getOidcConfig(): OIDCConfig | undefined;
  setOidcConfig(value?: OIDCConfig): void;
  hasOidcConfig(): boolean;
  clearOidcConfig(): void;

  getSequence(): number;
  setSequence(value: number): void;

  getAppConfigCase(): Application.AppConfigCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Application.AsObject;
  static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
  static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Application;
  static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
}

export namespace Application {
  export type AsObject = {
    id: string,
    state: AppState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    name: string,
    oidcConfig?: OIDCConfig.AsObject,
    sequence: number,
  }

  export enum AppConfigCase { 
    APP_CONFIG_NOT_SET = 0,
    OIDC_CONFIG = 8,
  }
}

export class ApplicationUpdate extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationUpdate): ApplicationUpdate.AsObject;
  static serializeBinaryToWriter(message: ApplicationUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationUpdate;
  static deserializeBinaryFromReader(message: ApplicationUpdate, reader: jspb.BinaryReader): ApplicationUpdate;
}

export namespace ApplicationUpdate {
  export type AsObject = {
    projectId: string,
    id: string,
    name: string,
  }
}

export class OIDCConfig extends jspb.Message {
  getRedirectUrisList(): Array<string>;
  setRedirectUrisList(value: Array<string>): void;
  clearRedirectUrisList(): void;
  addRedirectUris(value: string, index?: number): void;

  getResponseTypesList(): Array<OIDCResponseType>;
  setResponseTypesList(value: Array<OIDCResponseType>): void;
  clearResponseTypesList(): void;
  addResponseTypes(value: OIDCResponseType, index?: number): void;

  getGrantTypesList(): Array<OIDCGrantType>;
  setGrantTypesList(value: Array<OIDCGrantType>): void;
  clearGrantTypesList(): void;
  addGrantTypes(value: OIDCGrantType, index?: number): void;

  getApplicationType(): OIDCApplicationType;
  setApplicationType(value: OIDCApplicationType): void;

  getClientId(): string;
  setClientId(value: string): void;

  getClientSecret(): string;
  setClientSecret(value: string): void;

  getAuthMethodType(): OIDCAuthMethodType;
  setAuthMethodType(value: OIDCAuthMethodType): void;

  getPostLogoutRedirectUrisList(): Array<string>;
  setPostLogoutRedirectUrisList(value: Array<string>): void;
  clearPostLogoutRedirectUrisList(): void;
  addPostLogoutRedirectUris(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OIDCConfig.AsObject;
  static toObject(includeInstance: boolean, msg: OIDCConfig): OIDCConfig.AsObject;
  static serializeBinaryToWriter(message: OIDCConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OIDCConfig;
  static deserializeBinaryFromReader(message: OIDCConfig, reader: jspb.BinaryReader): OIDCConfig;
}

export namespace OIDCConfig {
  export type AsObject = {
    redirectUrisList: Array<string>,
    responseTypesList: Array<OIDCResponseType>,
    grantTypesList: Array<OIDCGrantType>,
    applicationType: OIDCApplicationType,
    clientId: string,
    clientSecret: string,
    authMethodType: OIDCAuthMethodType,
    postLogoutRedirectUrisList: Array<string>,
  }
}

export class OIDCApplicationCreate extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getRedirectUrisList(): Array<string>;
  setRedirectUrisList(value: Array<string>): void;
  clearRedirectUrisList(): void;
  addRedirectUris(value: string, index?: number): void;

  getResponseTypesList(): Array<OIDCResponseType>;
  setResponseTypesList(value: Array<OIDCResponseType>): void;
  clearResponseTypesList(): void;
  addResponseTypes(value: OIDCResponseType, index?: number): void;

  getGrantTypesList(): Array<OIDCGrantType>;
  setGrantTypesList(value: Array<OIDCGrantType>): void;
  clearGrantTypesList(): void;
  addGrantTypes(value: OIDCGrantType, index?: number): void;

  getApplicationType(): OIDCApplicationType;
  setApplicationType(value: OIDCApplicationType): void;

  getAuthMethodType(): OIDCAuthMethodType;
  setAuthMethodType(value: OIDCAuthMethodType): void;

  getPostLogoutRedirectUrisList(): Array<string>;
  setPostLogoutRedirectUrisList(value: Array<string>): void;
  clearPostLogoutRedirectUrisList(): void;
  addPostLogoutRedirectUris(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OIDCApplicationCreate.AsObject;
  static toObject(includeInstance: boolean, msg: OIDCApplicationCreate): OIDCApplicationCreate.AsObject;
  static serializeBinaryToWriter(message: OIDCApplicationCreate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OIDCApplicationCreate;
  static deserializeBinaryFromReader(message: OIDCApplicationCreate, reader: jspb.BinaryReader): OIDCApplicationCreate;
}

export namespace OIDCApplicationCreate {
  export type AsObject = {
    projectId: string,
    name: string,
    redirectUrisList: Array<string>,
    responseTypesList: Array<OIDCResponseType>,
    grantTypesList: Array<OIDCGrantType>,
    applicationType: OIDCApplicationType,
    authMethodType: OIDCAuthMethodType,
    postLogoutRedirectUrisList: Array<string>,
  }
}

export class OIDCConfigUpdate extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getApplicationId(): string;
  setApplicationId(value: string): void;

  getRedirectUrisList(): Array<string>;
  setRedirectUrisList(value: Array<string>): void;
  clearRedirectUrisList(): void;
  addRedirectUris(value: string, index?: number): void;

  getResponseTypesList(): Array<OIDCResponseType>;
  setResponseTypesList(value: Array<OIDCResponseType>): void;
  clearResponseTypesList(): void;
  addResponseTypes(value: OIDCResponseType, index?: number): void;

  getGrantTypesList(): Array<OIDCGrantType>;
  setGrantTypesList(value: Array<OIDCGrantType>): void;
  clearGrantTypesList(): void;
  addGrantTypes(value: OIDCGrantType, index?: number): void;

  getApplicationType(): OIDCApplicationType;
  setApplicationType(value: OIDCApplicationType): void;

  getAuthMethodType(): OIDCAuthMethodType;
  setAuthMethodType(value: OIDCAuthMethodType): void;

  getPostLogoutRedirectUrisList(): Array<string>;
  setPostLogoutRedirectUrisList(value: Array<string>): void;
  clearPostLogoutRedirectUrisList(): void;
  addPostLogoutRedirectUris(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OIDCConfigUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: OIDCConfigUpdate): OIDCConfigUpdate.AsObject;
  static serializeBinaryToWriter(message: OIDCConfigUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OIDCConfigUpdate;
  static deserializeBinaryFromReader(message: OIDCConfigUpdate, reader: jspb.BinaryReader): OIDCConfigUpdate;
}

export namespace OIDCConfigUpdate {
  export type AsObject = {
    projectId: string,
    applicationId: string,
    redirectUrisList: Array<string>,
    responseTypesList: Array<OIDCResponseType>,
    grantTypesList: Array<OIDCGrantType>,
    applicationType: OIDCApplicationType,
    authMethodType: OIDCAuthMethodType,
    postLogoutRedirectUrisList: Array<string>,
  }
}

export class ClientSecret extends jspb.Message {
  getClientSecret(): string;
  setClientSecret(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClientSecret.AsObject;
  static toObject(includeInstance: boolean, msg: ClientSecret): ClientSecret.AsObject;
  static serializeBinaryToWriter(message: ClientSecret, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClientSecret;
  static deserializeBinaryFromReader(message: ClientSecret, reader: jspb.BinaryReader): ClientSecret;
}

export namespace ClientSecret {
  export type AsObject = {
    clientSecret: string,
  }
}

export class ApplicationView extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getState(): AppState;
  setState(value: AppState): void;

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

  getOidcConfig(): OIDCConfig | undefined;
  setOidcConfig(value?: OIDCConfig): void;
  hasOidcConfig(): boolean;
  clearOidcConfig(): void;

  getSequence(): number;
  setSequence(value: number): void;

  getAppConfigCase(): ApplicationView.AppConfigCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationView.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationView): ApplicationView.AsObject;
  static serializeBinaryToWriter(message: ApplicationView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationView;
  static deserializeBinaryFromReader(message: ApplicationView, reader: jspb.BinaryReader): ApplicationView;
}

export namespace ApplicationView {
  export type AsObject = {
    id: string,
    state: AppState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    name: string,
    oidcConfig?: OIDCConfig.AsObject,
    sequence: number,
  }

  export enum AppConfigCase { 
    APP_CONFIG_NOT_SET = 0,
    OIDC_CONFIG = 8,
  }
}

export class ApplicationSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<ApplicationView>;
  setResultList(value: Array<ApplicationView>): void;
  clearResultList(): void;
  addResult(value?: ApplicationView, index?: number): ApplicationView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationSearchResponse): ApplicationSearchResponse.AsObject;
  static serializeBinaryToWriter(message: ApplicationSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationSearchResponse;
  static deserializeBinaryFromReader(message: ApplicationSearchResponse, reader: jspb.BinaryReader): ApplicationSearchResponse;
}

export namespace ApplicationSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<ApplicationView.AsObject>,
  }
}

export class ApplicationSearchRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<ApplicationSearchQuery>;
  setQueriesList(value: Array<ApplicationSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: ApplicationSearchQuery, index?: number): ApplicationSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationSearchRequest): ApplicationSearchRequest.AsObject;
  static serializeBinaryToWriter(message: ApplicationSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationSearchRequest;
  static deserializeBinaryFromReader(message: ApplicationSearchRequest, reader: jspb.BinaryReader): ApplicationSearchRequest;
}

export namespace ApplicationSearchRequest {
  export type AsObject = {
    projectId: string,
    offset: number,
    limit: number,
    queriesList: Array<ApplicationSearchQuery.AsObject>,
  }
}

export class ApplicationSearchQuery extends jspb.Message {
  getKey(): ApplicationSearchKey;
  setKey(value: ApplicationSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationSearchQuery): ApplicationSearchQuery.AsObject;
  static serializeBinaryToWriter(message: ApplicationSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationSearchQuery;
  static deserializeBinaryFromReader(message: ApplicationSearchQuery, reader: jspb.BinaryReader): ApplicationSearchQuery;
}

export namespace ApplicationSearchQuery {
  export type AsObject = {
    key: ApplicationSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class ProjectGrant extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getGrantedOrgId(): string;
  setGrantedOrgId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  getState(): ProjectGrantState;
  setState(value: ProjectGrantState): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrant.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrant): ProjectGrant.AsObject;
  static serializeBinaryToWriter(message: ProjectGrant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrant;
  static deserializeBinaryFromReader(message: ProjectGrant, reader: jspb.BinaryReader): ProjectGrant;
}

export namespace ProjectGrant {
  export type AsObject = {
    id: string,
    projectId: string,
    grantedOrgId: string,
    roleKeysList: Array<string>,
    state: ProjectGrantState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class ProjectGrantCreate extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getGrantedOrgId(): string;
  setGrantedOrgId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantCreate.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantCreate): ProjectGrantCreate.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantCreate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantCreate;
  static deserializeBinaryFromReader(message: ProjectGrantCreate, reader: jspb.BinaryReader): ProjectGrantCreate;
}

export namespace ProjectGrantCreate {
  export type AsObject = {
    projectId: string,
    grantedOrgId: string,
    roleKeysList: Array<string>,
  }
}

export class ProjectGrantUpdate extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantUpdate): ProjectGrantUpdate.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantUpdate;
  static deserializeBinaryFromReader(message: ProjectGrantUpdate, reader: jspb.BinaryReader): ProjectGrantUpdate;
}

export namespace ProjectGrantUpdate {
  export type AsObject = {
    projectId: string,
    id: string,
    roleKeysList: Array<string>,
  }
}

export class ProjectGrantID extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantID.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantID): ProjectGrantID.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantID;
  static deserializeBinaryFromReader(message: ProjectGrantID, reader: jspb.BinaryReader): ProjectGrantID;
}

export namespace ProjectGrantID {
  export type AsObject = {
    projectId: string,
    id: string,
  }
}

export class ProjectGrantView extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getGrantedOrgId(): string;
  setGrantedOrgId(value: string): void;

  getGrantedOrgName(): string;
  setGrantedOrgName(value: string): void;

  getGrantedOrgDomain(): string;
  setGrantedOrgDomain(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  getState(): ProjectGrantState;
  setState(value: ProjectGrantState): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getProjectName(): string;
  setProjectName(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantView.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantView): ProjectGrantView.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantView;
  static deserializeBinaryFromReader(message: ProjectGrantView, reader: jspb.BinaryReader): ProjectGrantView;
}

export namespace ProjectGrantView {
  export type AsObject = {
    id: string,
    projectId: string,
    grantedOrgId: string,
    grantedOrgName: string,
    grantedOrgDomain: string,
    roleKeysList: Array<string>,
    state: ProjectGrantState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    projectName: string,
    sequence: number,
  }
}

export class ProjectGrantSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<ProjectGrantView>;
  setResultList(value: Array<ProjectGrantView>): void;
  clearResultList(): void;
  addResult(value?: ProjectGrantView, index?: number): ProjectGrantView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantSearchResponse): ProjectGrantSearchResponse.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantSearchResponse;
  static deserializeBinaryFromReader(message: ProjectGrantSearchResponse, reader: jspb.BinaryReader): ProjectGrantSearchResponse;
}

export namespace ProjectGrantSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<ProjectGrantView.AsObject>,
  }
}

export class ProjectGrantSearchRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantSearchRequest): ProjectGrantSearchRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantSearchRequest;
  static deserializeBinaryFromReader(message: ProjectGrantSearchRequest, reader: jspb.BinaryReader): ProjectGrantSearchRequest;
}

export namespace ProjectGrantSearchRequest {
  export type AsObject = {
    projectId: string,
    offset: number,
    limit: number,
  }
}

export class GrantedProjectSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<ProjectSearchQuery>;
  setQueriesList(value: Array<ProjectSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: ProjectSearchQuery, index?: number): ProjectSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantedProjectSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GrantedProjectSearchRequest): GrantedProjectSearchRequest.AsObject;
  static serializeBinaryToWriter(message: GrantedProjectSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantedProjectSearchRequest;
  static deserializeBinaryFromReader(message: GrantedProjectSearchRequest, reader: jspb.BinaryReader): GrantedProjectSearchRequest;
}

export namespace GrantedProjectSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    queriesList: Array<ProjectSearchQuery.AsObject>,
  }
}

export class ProjectGrantMemberRoles extends jspb.Message {
  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMemberRoles.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMemberRoles): ProjectGrantMemberRoles.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMemberRoles, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMemberRoles;
  static deserializeBinaryFromReader(message: ProjectGrantMemberRoles, reader: jspb.BinaryReader): ProjectGrantMemberRoles;
}

export namespace ProjectGrantMemberRoles {
  export type AsObject = {
    rolesList: Array<string>,
  }
}

export class ProjectGrantMember extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMember.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMember): ProjectGrantMember.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMember, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMember;
  static deserializeBinaryFromReader(message: ProjectGrantMember, reader: jspb.BinaryReader): ProjectGrantMember;
}

export namespace ProjectGrantMember {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class ProjectGrantMemberAdd extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getGrantId(): string;
  setGrantId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMemberAdd.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMemberAdd): ProjectGrantMemberAdd.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMemberAdd, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMemberAdd;
  static deserializeBinaryFromReader(message: ProjectGrantMemberAdd, reader: jspb.BinaryReader): ProjectGrantMemberAdd;
}

export namespace ProjectGrantMemberAdd {
  export type AsObject = {
    projectId: string,
    grantId: string,
    userId: string,
    rolesList: Array<string>,
  }
}

export class ProjectGrantMemberChange extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getGrantId(): string;
  setGrantId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMemberChange.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMemberChange): ProjectGrantMemberChange.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMemberChange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMemberChange;
  static deserializeBinaryFromReader(message: ProjectGrantMemberChange, reader: jspb.BinaryReader): ProjectGrantMemberChange;
}

export namespace ProjectGrantMemberChange {
  export type AsObject = {
    projectId: string,
    grantId: string,
    userId: string,
    rolesList: Array<string>,
  }
}

export class ProjectGrantMemberRemove extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getGrantId(): string;
  setGrantId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMemberRemove.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMemberRemove): ProjectGrantMemberRemove.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMemberRemove, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMemberRemove;
  static deserializeBinaryFromReader(message: ProjectGrantMemberRemove, reader: jspb.BinaryReader): ProjectGrantMemberRemove;
}

export namespace ProjectGrantMemberRemove {
  export type AsObject = {
    projectId: string,
    grantId: string,
    userId: string,
  }
}

export class ProjectGrantMemberView extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getUserName(): string;
  setUserName(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getFirstName(): string;
  setFirstName(value: string): void;

  getLastName(): string;
  setLastName(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMemberView.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMemberView): ProjectGrantMemberView.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMemberView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMemberView;
  static deserializeBinaryFromReader(message: ProjectGrantMemberView, reader: jspb.BinaryReader): ProjectGrantMemberView;
}

export namespace ProjectGrantMemberView {
  export type AsObject = {
    userId: string,
    userName: string,
    email: string,
    firstName: string,
    lastName: string,
    rolesList: Array<string>,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class ProjectGrantMemberSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<ProjectGrantMemberView>;
  setResultList(value: Array<ProjectGrantMemberView>): void;
  clearResultList(): void;
  addResult(value?: ProjectGrantMemberView, index?: number): ProjectGrantMemberView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMemberSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMemberSearchResponse): ProjectGrantMemberSearchResponse.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMemberSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMemberSearchResponse;
  static deserializeBinaryFromReader(message: ProjectGrantMemberSearchResponse, reader: jspb.BinaryReader): ProjectGrantMemberSearchResponse;
}

export namespace ProjectGrantMemberSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<ProjectGrantMemberView.AsObject>,
  }
}

export class ProjectGrantMemberSearchRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getGrantId(): string;
  setGrantId(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<ProjectGrantMemberSearchQuery>;
  setQueriesList(value: Array<ProjectGrantMemberSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: ProjectGrantMemberSearchQuery, index?: number): ProjectGrantMemberSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMemberSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMemberSearchRequest): ProjectGrantMemberSearchRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMemberSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMemberSearchRequest;
  static deserializeBinaryFromReader(message: ProjectGrantMemberSearchRequest, reader: jspb.BinaryReader): ProjectGrantMemberSearchRequest;
}

export namespace ProjectGrantMemberSearchRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
    offset: number,
    limit: number,
    queriesList: Array<ProjectGrantMemberSearchQuery.AsObject>,
  }
}

export class ProjectGrantMemberSearchQuery extends jspb.Message {
  getKey(): ProjectGrantMemberSearchKey;
  setKey(value: ProjectGrantMemberSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantMemberSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantMemberSearchQuery): ProjectGrantMemberSearchQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantMemberSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantMemberSearchQuery;
  static deserializeBinaryFromReader(message: ProjectGrantMemberSearchQuery, reader: jspb.BinaryReader): ProjectGrantMemberSearchQuery;
}

export namespace ProjectGrantMemberSearchQuery {
  export type AsObject = {
    key: ProjectGrantMemberSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class UserGrant extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getOrgId(): string;
  setOrgId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  getState(): UserGrantState;
  setState(value: UserGrantState): void;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasCreationDate(): boolean;
  clearCreationDate(): void;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasChangeDate(): boolean;
  clearChangeDate(): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrant.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrant): UserGrant.AsObject;
  static serializeBinaryToWriter(message: UserGrant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrant;
  static deserializeBinaryFromReader(message: UserGrant, reader: jspb.BinaryReader): UserGrant;
}

export namespace UserGrant {
  export type AsObject = {
    id: string,
    userId: string,
    orgId: string,
    projectId: string,
    roleKeysList: Array<string>,
    state: UserGrantState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class UserGrantCreate extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantCreate.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantCreate): UserGrantCreate.AsObject;
  static serializeBinaryToWriter(message: UserGrantCreate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantCreate;
  static deserializeBinaryFromReader(message: UserGrantCreate, reader: jspb.BinaryReader): UserGrantCreate;
}

export namespace UserGrantCreate {
  export type AsObject = {
    userId: string,
    projectId: string,
    roleKeysList: Array<string>,
  }
}

export class UserGrantUpdate extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantUpdate): UserGrantUpdate.AsObject;
  static serializeBinaryToWriter(message: UserGrantUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantUpdate;
  static deserializeBinaryFromReader(message: UserGrantUpdate, reader: jspb.BinaryReader): UserGrantUpdate;
}

export namespace UserGrantUpdate {
  export type AsObject = {
    userId: string,
    id: string,
    roleKeysList: Array<string>,
  }
}

export class UserGrantID extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantID.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantID): UserGrantID.AsObject;
  static serializeBinaryToWriter(message: UserGrantID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantID;
  static deserializeBinaryFromReader(message: UserGrantID, reader: jspb.BinaryReader): UserGrantID;
}

export namespace UserGrantID {
  export type AsObject = {
    userId: string,
    id: string,
  }
}

export class ProjectUserGrantID extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectUserGrantID.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectUserGrantID): ProjectUserGrantID.AsObject;
  static serializeBinaryToWriter(message: ProjectUserGrantID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectUserGrantID;
  static deserializeBinaryFromReader(message: ProjectUserGrantID, reader: jspb.BinaryReader): ProjectUserGrantID;
}

export namespace ProjectUserGrantID {
  export type AsObject = {
    projectId: string,
    userId: string,
    id: string,
  }
}

export class ProjectUserGrantUpdate extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectUserGrantUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectUserGrantUpdate): ProjectUserGrantUpdate.AsObject;
  static serializeBinaryToWriter(message: ProjectUserGrantUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectUserGrantUpdate;
  static deserializeBinaryFromReader(message: ProjectUserGrantUpdate, reader: jspb.BinaryReader): ProjectUserGrantUpdate;
}

export namespace ProjectUserGrantUpdate {
  export type AsObject = {
    projectId: string,
    userId: string,
    id: string,
    roleKeysList: Array<string>,
  }
}

export class ProjectGrantUserGrantID extends jspb.Message {
  getProjectGrantId(): string;
  setProjectGrantId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantUserGrantID.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantUserGrantID): ProjectGrantUserGrantID.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantUserGrantID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantUserGrantID;
  static deserializeBinaryFromReader(message: ProjectGrantUserGrantID, reader: jspb.BinaryReader): ProjectGrantUserGrantID;
}

export namespace ProjectGrantUserGrantID {
  export type AsObject = {
    projectGrantId: string,
    userId: string,
    id: string,
  }
}

export class ProjectGrantUserGrantCreate extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getOrgId(): string;
  setOrgId(value: string): void;

  getProjectGrantId(): string;
  setProjectGrantId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantUserGrantCreate.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantUserGrantCreate): ProjectGrantUserGrantCreate.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantUserGrantCreate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantUserGrantCreate;
  static deserializeBinaryFromReader(message: ProjectGrantUserGrantCreate, reader: jspb.BinaryReader): ProjectGrantUserGrantCreate;
}

export namespace ProjectGrantUserGrantCreate {
  export type AsObject = {
    userId: string,
    orgId: string,
    projectGrantId: string,
    projectId: string,
    roleKeysList: Array<string>,
  }
}

export class ProjectGrantUserGrantUpdate extends jspb.Message {
  getProjectGrantId(): string;
  setProjectGrantId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantUserGrantUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantUserGrantUpdate): ProjectGrantUserGrantUpdate.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantUserGrantUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantUserGrantUpdate;
  static deserializeBinaryFromReader(message: ProjectGrantUserGrantUpdate, reader: jspb.BinaryReader): ProjectGrantUserGrantUpdate;
}

export namespace ProjectGrantUserGrantUpdate {
  export type AsObject = {
    projectGrantId: string,
    userId: string,
    id: string,
    roleKeysList: Array<string>,
  }
}

export class UserGrantView extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getOrgId(): string;
  setOrgId(value: string): void;

  getProjectId(): string;
  setProjectId(value: string): void;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): void;
  clearRoleKeysList(): void;
  addRoleKeys(value: string, index?: number): void;

  getState(): UserGrantState;
  setState(value: UserGrantState): void;

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

  getEmail(): string;
  setEmail(value: string): void;

  getOrgName(): string;
  setOrgName(value: string): void;

  getOrgDomain(): string;
  setOrgDomain(value: string): void;

  getProjectName(): string;
  setProjectName(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  getResourceOwner(): string;
  setResourceOwner(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantView.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantView): UserGrantView.AsObject;
  static serializeBinaryToWriter(message: UserGrantView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantView;
  static deserializeBinaryFromReader(message: UserGrantView, reader: jspb.BinaryReader): UserGrantView;
}

export namespace UserGrantView {
  export type AsObject = {
    id: string,
    userId: string,
    orgId: string,
    projectId: string,
    roleKeysList: Array<string>,
    state: UserGrantState,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    userName: string,
    firstName: string,
    lastName: string,
    email: string,
    orgName: string,
    orgDomain: string,
    projectName: string,
    sequence: number,
    resourceOwner: string,
  }
}

export class UserGrantSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<UserGrantView>;
  setResultList(value: Array<UserGrantView>): void;
  clearResultList(): void;
  addResult(value?: UserGrantView, index?: number): UserGrantView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantSearchResponse): UserGrantSearchResponse.AsObject;
  static serializeBinaryToWriter(message: UserGrantSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantSearchResponse;
  static deserializeBinaryFromReader(message: UserGrantSearchResponse, reader: jspb.BinaryReader): UserGrantSearchResponse;
}

export namespace UserGrantSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<UserGrantView.AsObject>,
  }
}

export class UserGrantSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<UserGrantSearchQuery>;
  setQueriesList(value: Array<UserGrantSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: UserGrantSearchQuery, index?: number): UserGrantSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantSearchRequest): UserGrantSearchRequest.AsObject;
  static serializeBinaryToWriter(message: UserGrantSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantSearchRequest;
  static deserializeBinaryFromReader(message: UserGrantSearchRequest, reader: jspb.BinaryReader): UserGrantSearchRequest;
}

export namespace UserGrantSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    queriesList: Array<UserGrantSearchQuery.AsObject>,
  }
}

export class UserGrantSearchQuery extends jspb.Message {
  getKey(): UserGrantSearchKey;
  setKey(value: UserGrantSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantSearchQuery): UserGrantSearchQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantSearchQuery;
  static deserializeBinaryFromReader(message: UserGrantSearchQuery, reader: jspb.BinaryReader): UserGrantSearchQuery;
}

export namespace UserGrantSearchQuery {
  export type AsObject = {
    key: UserGrantSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class ProjectUserGrantSearchRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<UserGrantSearchQuery>;
  setQueriesList(value: Array<UserGrantSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: UserGrantSearchQuery, index?: number): UserGrantSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectUserGrantSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectUserGrantSearchRequest): ProjectUserGrantSearchRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectUserGrantSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectUserGrantSearchRequest;
  static deserializeBinaryFromReader(message: ProjectUserGrantSearchRequest, reader: jspb.BinaryReader): ProjectUserGrantSearchRequest;
}

export namespace ProjectUserGrantSearchRequest {
  export type AsObject = {
    projectId: string,
    offset: number,
    limit: number,
    queriesList: Array<UserGrantSearchQuery.AsObject>,
  }
}

export class ProjectGrantUserGrantSearchRequest extends jspb.Message {
  getProjectGrantId(): string;
  setProjectGrantId(value: string): void;

  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<UserGrantSearchQuery>;
  setQueriesList(value: Array<UserGrantSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: UserGrantSearchQuery, index?: number): UserGrantSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantUserGrantSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantUserGrantSearchRequest): ProjectGrantUserGrantSearchRequest.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantUserGrantSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantUserGrantSearchRequest;
  static deserializeBinaryFromReader(message: ProjectGrantUserGrantSearchRequest, reader: jspb.BinaryReader): ProjectGrantUserGrantSearchRequest;
}

export namespace ProjectGrantUserGrantSearchRequest {
  export type AsObject = {
    projectGrantId: string,
    offset: number,
    limit: number,
    queriesList: Array<UserGrantSearchQuery.AsObject>,
  }
}

export class AuthGrantSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getSortingColumn(): AuthGrantSearchKey;
  setSortingColumn(value: AuthGrantSearchKey): void;

  getAsc(): boolean;
  setAsc(value: boolean): void;

  getQueriesList(): Array<AuthGrantSearchQuery>;
  setQueriesList(value: Array<AuthGrantSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: AuthGrantSearchQuery, index?: number): AuthGrantSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthGrantSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AuthGrantSearchRequest): AuthGrantSearchRequest.AsObject;
  static serializeBinaryToWriter(message: AuthGrantSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthGrantSearchRequest;
  static deserializeBinaryFromReader(message: AuthGrantSearchRequest, reader: jspb.BinaryReader): AuthGrantSearchRequest;
}

export namespace AuthGrantSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    sortingColumn: AuthGrantSearchKey,
    asc: boolean,
    queriesList: Array<AuthGrantSearchQuery.AsObject>,
  }
}

export class AuthGrantSearchQuery extends jspb.Message {
  getKey(): AuthGrantSearchKey;
  setKey(value: AuthGrantSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthGrantSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: AuthGrantSearchQuery): AuthGrantSearchQuery.AsObject;
  static serializeBinaryToWriter(message: AuthGrantSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthGrantSearchQuery;
  static deserializeBinaryFromReader(message: AuthGrantSearchQuery, reader: jspb.BinaryReader): AuthGrantSearchQuery;
}

export namespace AuthGrantSearchQuery {
  export type AsObject = {
    key: AuthGrantSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class AuthGrantSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<AuthGrant>;
  setResultList(value: Array<AuthGrant>): void;
  clearResultList(): void;
  addResult(value?: AuthGrant, index?: number): AuthGrant;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthGrantSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AuthGrantSearchResponse): AuthGrantSearchResponse.AsObject;
  static serializeBinaryToWriter(message: AuthGrantSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthGrantSearchResponse;
  static deserializeBinaryFromReader(message: AuthGrantSearchResponse, reader: jspb.BinaryReader): AuthGrantSearchResponse;
}

export namespace AuthGrantSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<AuthGrant.AsObject>,
  }
}

export class AuthGrant extends jspb.Message {
  getOrgid(): string;
  setOrgid(value: string): void;

  getProjectid(): string;
  setProjectid(value: string): void;

  getUserid(): string;
  setUserid(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthGrant.AsObject;
  static toObject(includeInstance: boolean, msg: AuthGrant): AuthGrant.AsObject;
  static serializeBinaryToWriter(message: AuthGrant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthGrant;
  static deserializeBinaryFromReader(message: AuthGrant, reader: jspb.BinaryReader): AuthGrant;
}

export namespace AuthGrant {
  export type AsObject = {
    orgid: string,
    projectid: string,
    userid: string,
    rolesList: Array<string>,
  }
}

export enum UserState { 
  USERSTATE_UNSPECIFIED = 0,
  USERSTATE_ACTIVE = 1,
  USERSTATE_INACTIVE = 2,
  USERSTATE_DELETED = 3,
  USERSTATE_LOCKED = 4,
  USERSTATE_SUSPEND = 5,
  USERSTATE_INITIAL = 6,
}
export enum Gender { 
  GENDER_UNSPECIFIED = 0,
  GENDER_FEMALE = 1,
  GENDER_MALE = 2,
  GENDER_DIVERSE = 3,
}
export enum UserSearchKey { 
  USERSEARCHKEY_UNSPECIFIED = 0,
  USERSEARCHKEY_USER_NAME = 1,
  USERSEARCHKEY_FIRST_NAME = 2,
  USERSEARCHKEY_LAST_NAME = 3,
  USERSEARCHKEY_NICK_NAME = 4,
  USERSEARCHKEY_DISPLAY_NAME = 5,
  USERSEARCHKEY_EMAIL = 6,
  USERSEARCHKEY_STATE = 7,
}
export enum SearchMethod { 
  SEARCHMETHOD_EQUALS = 0,
  SEARCHMETHOD_STARTS_WITH = 1,
  SEARCHMETHOD_CONTAINS = 2,
  SEARCHMETHOD_EQUALS_IGNORE_CASE = 3,
  SEARCHMETHOD_STARTS_WITH_IGNORE_CASE = 4,
  SEARCHMETHOD_CONTAINS_IGNORE_CASE = 5,
}
export enum MfaType { 
  MFATYPE_UNSPECIFIED = 0,
  MFATYPE_SMS = 1,
  MFATYPE_OTP = 2,
}
export enum MFAState { 
  MFASTATE_UNSPECIFIED = 0,
  MFASTATE_NOT_READY = 1,
  MFASTATE_READY = 2,
  MFASTATE_REMOVED = 3,
}
export enum NotificationType { 
  NOTIFICATIONTYPE_EMAIL = 0,
  NOTIFICATIONTYPE_SMS = 1,
}
export enum PolicyState { 
  POLICYSTATE_UNSPECIFIED = 0,
  POLICYSTATE_ACTIVE = 1,
  POLICYSTATE_INACTIVE = 2,
  POLICYSTATE_DELETED = 3,
}
export enum OrgState { 
  ORGSTATE_UNSPECIFIED = 0,
  ORGSTATE_ACTIVE = 1,
  ORGSTATE_INACTIVE = 2,
}
export enum OrgDomainSearchKey { 
  ORGDOMAINSEARCHKEY_UNSPECIFIED = 0,
  ORGDOMAINSEARCHKEY_DOMAIN = 1,
}
export enum OrgMemberSearchKey { 
  ORGMEMBERSEARCHKEY_UNSPECIFIED = 0,
  ORGMEMBERSEARCHKEY_FIRST_NAME = 1,
  ORGMEMBERSEARCHKEY_LAST_NAME = 2,
  ORGMEMBERSEARCHKEY_EMAIL = 3,
  ORGMEMBERSEARCHKEY_USER_ID = 4,
}
export enum ProjectSearchKey { 
  PROJECTSEARCHKEY_UNSPECIFIED = 0,
  PROJECTSEARCHKEY_PROJECT_NAME = 1,
}
export enum ProjectState { 
  PROJECTSTATE_UNSPECIFIED = 0,
  PROJECTSTATE_ACTIVE = 1,
  PROJECTSTATE_INACTIVE = 2,
}
export enum ProjectType { 
  PROJECTTYPE_UNSPECIFIED = 0,
  PROJECTTYPE_OWNED = 1,
  PROJECTTYPE_GRANTED = 2,
}
export enum ProjectRoleSearchKey { 
  PROJECTROLESEARCHKEY_UNSPECIFIED = 0,
  PROJECTROLESEARCHKEY_KEY = 1,
  PROJECTROLESEARCHKEY_DISPLAY_NAME = 2,
}
export enum ProjectMemberSearchKey { 
  PROJECTMEMBERSEARCHKEY_UNSPECIFIED = 0,
  PROJECTMEMBERSEARCHKEY_FIRST_NAME = 1,
  PROJECTMEMBERSEARCHKEY_LAST_NAME = 2,
  PROJECTMEMBERSEARCHKEY_EMAIL = 3,
  PROJECTMEMBERSEARCHKEY_USER_ID = 4,
  PROJECTMEMBERSEARCHKEY_USER_NAME = 5,
}
export enum AppState { 
  APPSTATE_UNSPECIFIED = 0,
  APPSTATE_ACTIVE = 1,
  APPSTATE_INACTIVE = 2,
}
export enum OIDCResponseType { 
  OIDCRESPONSETYPE_CODE = 0,
  OIDCRESPONSETYPE_ID_TOKEN = 1,
  OIDCRESPONSETYPE_TOKEN = 2,
}
export enum OIDCGrantType { 
  OIDCGRANTTYPE_AUTHORIZATION_CODE = 0,
  OIDCGRANTTYPE_IMPLICIT = 1,
  OIDCGRANTTYPE_REFRESH_TOKEN = 2,
}
export enum OIDCApplicationType { 
  OIDCAPPLICATIONTYPE_WEB = 0,
  OIDCAPPLICATIONTYPE_USER_AGENT = 1,
  OIDCAPPLICATIONTYPE_NATIVE = 2,
}
export enum OIDCAuthMethodType { 
  OIDCAUTHMETHODTYPE_BASIC = 0,
  OIDCAUTHMETHODTYPE_POST = 1,
  OIDCAUTHMETHODTYPE_NONE = 2,
}
export enum ApplicationSearchKey { 
  APPLICATIONSERACHKEY_UNSPECIFIED = 0,
  APPLICATIONSEARCHKEY_APP_NAME = 1,
}
export enum ProjectGrantState { 
  PROJECTGRANTSTATE_UNSPECIFIED = 0,
  PROJECTGRANTSTATE_ACTIVE = 1,
  PROJECTGRANTSTATE_INACTIVE = 2,
}
export enum ProjectGrantMemberSearchKey { 
  PROJECTGRANTMEMBERSEARCHKEY_UNSPECIFIED = 0,
  PROJECTGRANTMEMBERSEARCHKEY_FIRST_NAME = 1,
  PROJECTGRANTMEMBERSEARCHKEY_LAST_NAME = 2,
  PROJECTGRANTMEMBERSEARCHKEY_EMAIL = 3,
  PROJECTGRANTMEMBERSEARCHKEY_USER_ID = 4,
  PROJECTGRANTMEMBERSEARCHKEY_USER_NAME = 5,
}
export enum UserGrantState { 
  USERGRANTSTATE_UNSPECIFIED = 0,
  USERGRANTSTATE_ACTIVE = 1,
  USERGRANTSTATE_INACTIVE = 2,
}
export enum UserGrantSearchKey { 
  USERGRANTSEARCHKEY_UNSPECIFIED = 0,
  USERGRANTSEARCHKEY_PROJECT_ID = 1,
  USERGRANTSEARCHKEY_USER_ID = 2,
  USERGRANTSEARCHKEY_ORG_ID = 3,
}
export enum AuthGrantSearchKey { 
  AUTHGRANTSEARCHKEY_UNSPECIFIED = 0,
  AUTHGRANTSEARCHKEY_ORG_ID = 1,
  AUTHGRANTSEARCHKEY_PROJECT_ID = 2,
  AUTHGRANTSEARCHKEY_USER_ID = 3,
}

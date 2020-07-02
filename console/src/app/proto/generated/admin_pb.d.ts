import * as jspb from "google-protobuf"

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
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

  getUser(): CreateUserRequest | undefined;
  setUser(value?: CreateUserRequest): void;
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
    user?: CreateUserRequest.AsObject,
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
    sequence: number,
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

export class OrgIamPolicy extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): void;

  getDefault(): boolean;
  setDefault(value: boolean): void;

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
  toObject(includeInstance?: boolean): OrgIamPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: OrgIamPolicy): OrgIamPolicy.AsObject;
  static serializeBinaryToWriter(message: OrgIamPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgIamPolicy;
  static deserializeBinaryFromReader(message: OrgIamPolicy, reader: jspb.BinaryReader): OrgIamPolicy;
}

export namespace OrgIamPolicy {
  export type AsObject = {
    orgId: string,
    description: string,
    userLoginMustBeDomain: boolean,
    pb_default: boolean,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class OrgIamPolicyRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgIamPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: OrgIamPolicyRequest): OrgIamPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: OrgIamPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgIamPolicyRequest;
  static deserializeBinaryFromReader(message: OrgIamPolicyRequest, reader: jspb.BinaryReader): OrgIamPolicyRequest;
}

export namespace OrgIamPolicyRequest {
  export type AsObject = {
    orgId: string,
    description: string,
    userLoginMustBeDomain: boolean,
  }
}

export class OrgIamPolicyID extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgIamPolicyID.AsObject;
  static toObject(includeInstance: boolean, msg: OrgIamPolicyID): OrgIamPolicyID.AsObject;
  static serializeBinaryToWriter(message: OrgIamPolicyID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgIamPolicyID;
  static deserializeBinaryFromReader(message: OrgIamPolicyID, reader: jspb.BinaryReader): OrgIamPolicyID;
}

export namespace OrgIamPolicyID {
  export type AsObject = {
    orgId: string,
  }
}

export class IamMemberRoles extends jspb.Message {
  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IamMemberRoles.AsObject;
  static toObject(includeInstance: boolean, msg: IamMemberRoles): IamMemberRoles.AsObject;
  static serializeBinaryToWriter(message: IamMemberRoles, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IamMemberRoles;
  static deserializeBinaryFromReader(message: IamMemberRoles, reader: jspb.BinaryReader): IamMemberRoles;
}

export namespace IamMemberRoles {
  export type AsObject = {
    rolesList: Array<string>,
  }
}

export class IamMember extends jspb.Message {
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
  toObject(includeInstance?: boolean): IamMember.AsObject;
  static toObject(includeInstance: boolean, msg: IamMember): IamMember.AsObject;
  static serializeBinaryToWriter(message: IamMember, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IamMember;
  static deserializeBinaryFromReader(message: IamMember, reader: jspb.BinaryReader): IamMember;
}

export namespace IamMember {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
  }
}

export class AddIamMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIamMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddIamMemberRequest): AddIamMemberRequest.AsObject;
  static serializeBinaryToWriter(message: AddIamMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIamMemberRequest;
  static deserializeBinaryFromReader(message: AddIamMemberRequest, reader: jspb.BinaryReader): AddIamMemberRequest;
}

export namespace AddIamMemberRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
  }
}

export class ChangeIamMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): void;
  clearRolesList(): void;
  addRoles(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeIamMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeIamMemberRequest): ChangeIamMemberRequest.AsObject;
  static serializeBinaryToWriter(message: ChangeIamMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeIamMemberRequest;
  static deserializeBinaryFromReader(message: ChangeIamMemberRequest, reader: jspb.BinaryReader): ChangeIamMemberRequest;
}

export namespace ChangeIamMemberRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
  }
}

export class RemoveIamMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIamMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIamMemberRequest): RemoveIamMemberRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveIamMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIamMemberRequest;
  static deserializeBinaryFromReader(message: RemoveIamMemberRequest, reader: jspb.BinaryReader): RemoveIamMemberRequest;
}

export namespace RemoveIamMemberRequest {
  export type AsObject = {
    userId: string,
  }
}

export class IamMemberSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<IamMemberView>;
  setResultList(value: Array<IamMemberView>): void;
  clearResultList(): void;
  addResult(value?: IamMemberView, index?: number): IamMemberView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IamMemberSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: IamMemberSearchResponse): IamMemberSearchResponse.AsObject;
  static serializeBinaryToWriter(message: IamMemberSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IamMemberSearchResponse;
  static deserializeBinaryFromReader(message: IamMemberSearchResponse, reader: jspb.BinaryReader): IamMemberSearchResponse;
}

export namespace IamMemberSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<IamMemberView.AsObject>,
  }
}

export class IamMemberView extends jspb.Message {
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
  toObject(includeInstance?: boolean): IamMemberView.AsObject;
  static toObject(includeInstance: boolean, msg: IamMemberView): IamMemberView.AsObject;
  static serializeBinaryToWriter(message: IamMemberView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IamMemberView;
  static deserializeBinaryFromReader(message: IamMemberView, reader: jspb.BinaryReader): IamMemberView;
}

export namespace IamMemberView {
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

export class IamMemberSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getQueriesList(): Array<IamMemberSearchQuery>;
  setQueriesList(value: Array<IamMemberSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: IamMemberSearchQuery, index?: number): IamMemberSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IamMemberSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: IamMemberSearchRequest): IamMemberSearchRequest.AsObject;
  static serializeBinaryToWriter(message: IamMemberSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IamMemberSearchRequest;
  static deserializeBinaryFromReader(message: IamMemberSearchRequest, reader: jspb.BinaryReader): IamMemberSearchRequest;
}

export namespace IamMemberSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    queriesList: Array<IamMemberSearchQuery.AsObject>,
  }
}

export class IamMemberSearchQuery extends jspb.Message {
  getKey(): IamMemberSearchKey;
  setKey(value: IamMemberSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IamMemberSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IamMemberSearchQuery): IamMemberSearchQuery.AsObject;
  static serializeBinaryToWriter(message: IamMemberSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IamMemberSearchQuery;
  static deserializeBinaryFromReader(message: IamMemberSearchQuery, reader: jspb.BinaryReader): IamMemberSearchQuery;
}

export namespace IamMemberSearchQuery {
  export type AsObject = {
    key: IamMemberSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class FailedEventID extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): void;

  getViewName(): string;
  setViewName(value: string): void;

  getFailedSequence(): number;
  setFailedSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FailedEventID.AsObject;
  static toObject(includeInstance: boolean, msg: FailedEventID): FailedEventID.AsObject;
  static serializeBinaryToWriter(message: FailedEventID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FailedEventID;
  static deserializeBinaryFromReader(message: FailedEventID, reader: jspb.BinaryReader): FailedEventID;
}

export namespace FailedEventID {
  export type AsObject = {
    database: string,
    viewName: string,
    failedSequence: number,
  }
}

export class FailedEvents extends jspb.Message {
  getFailedEventsList(): Array<FailedEvent>;
  setFailedEventsList(value: Array<FailedEvent>): void;
  clearFailedEventsList(): void;
  addFailedEvents(value?: FailedEvent, index?: number): FailedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FailedEvents.AsObject;
  static toObject(includeInstance: boolean, msg: FailedEvents): FailedEvents.AsObject;
  static serializeBinaryToWriter(message: FailedEvents, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FailedEvents;
  static deserializeBinaryFromReader(message: FailedEvents, reader: jspb.BinaryReader): FailedEvents;
}

export namespace FailedEvents {
  export type AsObject = {
    failedEventsList: Array<FailedEvent.AsObject>,
  }
}

export class FailedEvent extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): void;

  getViewName(): string;
  setViewName(value: string): void;

  getFailedSequence(): number;
  setFailedSequence(value: number): void;

  getFailureCount(): number;
  setFailureCount(value: number): void;

  getErrorMessage(): string;
  setErrorMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FailedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: FailedEvent): FailedEvent.AsObject;
  static serializeBinaryToWriter(message: FailedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FailedEvent;
  static deserializeBinaryFromReader(message: FailedEvent, reader: jspb.BinaryReader): FailedEvent;
}

export namespace FailedEvent {
  export type AsObject = {
    database: string,
    viewName: string,
    failedSequence: number,
    failureCount: number,
    errorMessage: string,
  }
}

export class ViewID extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): void;

  getViewName(): string;
  setViewName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ViewID.AsObject;
  static toObject(includeInstance: boolean, msg: ViewID): ViewID.AsObject;
  static serializeBinaryToWriter(message: ViewID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ViewID;
  static deserializeBinaryFromReader(message: ViewID, reader: jspb.BinaryReader): ViewID;
}

export namespace ViewID {
  export type AsObject = {
    database: string,
    viewName: string,
  }
}

export class Views extends jspb.Message {
  getViewsList(): Array<View>;
  setViewsList(value: Array<View>): void;
  clearViewsList(): void;
  addViews(value?: View, index?: number): View;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Views.AsObject;
  static toObject(includeInstance: boolean, msg: Views): Views.AsObject;
  static serializeBinaryToWriter(message: Views, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Views;
  static deserializeBinaryFromReader(message: Views, reader: jspb.BinaryReader): Views;
}

export namespace Views {
  export type AsObject = {
    viewsList: Array<View.AsObject>,
  }
}

export class View extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): void;

  getViewName(): string;
  setViewName(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): View.AsObject;
  static toObject(includeInstance: boolean, msg: View): View.AsObject;
  static serializeBinaryToWriter(message: View, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): View;
  static deserializeBinaryFromReader(message: View, reader: jspb.BinaryReader): View;
}

export namespace View {
  export type AsObject = {
    database: string,
    viewName: string,
    sequence: number,
  }
}

export enum OrgState { 
  ORGSTATE_UNSPECIFIED = 0,
  ORGSTATE_ACTIVE = 1,
  ORGSTATE_INACTIVE = 2,
}
export enum OrgSearchKey { 
  ORGSEARCHKEY_UNSPECIFIED = 0,
  ORGSEARCHKEY_ORG_NAME = 1,
  ORGSEARCHKEY_DOMAIN = 2,
  ORGSEARCHKEY_STATE = 3,
}
export enum OrgSearchMethod { 
  ORGSEARCHMETHOD_EQUALS = 0,
  ORGSEARCHMETHOD_STARTS_WITH = 1,
  ORGSEARCHMETHOD_CONTAINS = 2,
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
export enum IamMemberSearchKey { 
  IAMMEMBERSEARCHKEY_UNSPECIFIED = 0,
  IAMMEMBERSEARCHKEY_FIRST_NAME = 1,
  IAMMEMBERSEARCHKEY_LAST_NAME = 2,
  IAMMEMBERSEARCHKEY_EMAIL = 3,
  IAMMEMBERSEARCHKEY_USER_ID = 4,
}
export enum SearchMethod { 
  SEARCHMETHOD_EQUALS = 0,
  SEARCHMETHOD_STARTS_WITH = 1,
  SEARCHMETHOD_CONTAINS = 2,
  SEARCHMETHOD_EQUALS_IGNORE_CASE = 3,
  SEARCHMETHOD_STARTS_WITH_IGNORE_CASE = 4,
  SEARCHMETHOD_CONTAINS_IGNORE_CASE = 5,
  SEARCHMETHOD_NOT_EQUALS = 6,
  SEARCHMETHOD_GREATER_THAN = 7,
  SEARCHMETHOD_LESS_THAN = 8,
  SEARCHMETHOD_IS_ONE_OF = 9,
  SEARCHMETHOD_LIST_CONTAINS = 10,
}

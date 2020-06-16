import * as jspb from "google-protobuf"

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
import * as authoption_options_pb from './authoption/options_pb';

export class UserSessionViews extends jspb.Message {
  getUserSessionsList(): Array<UserSessionView>;
  setUserSessionsList(value: Array<UserSessionView>): void;
  clearUserSessionsList(): void;
  addUserSessions(value?: UserSessionView, index?: number): UserSessionView;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSessionViews.AsObject;
  static toObject(includeInstance: boolean, msg: UserSessionViews): UserSessionViews.AsObject;
  static serializeBinaryToWriter(message: UserSessionViews, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSessionViews;
  static deserializeBinaryFromReader(message: UserSessionViews, reader: jspb.BinaryReader): UserSessionViews;
}

export namespace UserSessionViews {
  export type AsObject = {
    userSessionsList: Array<UserSessionView.AsObject>,
  }
}

export class UserSessionView extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAgentId(): string;
  setAgentId(value: string): void;

  getAuthState(): UserSessionState;
  setAuthState(value: UserSessionState): void;

  getUserId(): string;
  setUserId(value: string): void;

  getUserName(): string;
  setUserName(value: string): void;

  getSequence(): number;
  setSequence(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSessionView.AsObject;
  static toObject(includeInstance: boolean, msg: UserSessionView): UserSessionView.AsObject;
  static serializeBinaryToWriter(message: UserSessionView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSessionView;
  static deserializeBinaryFromReader(message: UserSessionView, reader: jspb.BinaryReader): UserSessionView;
}

export namespace UserSessionView {
  export type AsObject = {
    id: string,
    agentId: string,
    authState: UserSessionState,
    userId: string,
    userName: string,
    sequence: number,
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

  getActivationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setActivationDate(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasActivationDate(): boolean;
  clearActivationDate(): void;

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

  getPasswordChangeRequired(): boolean;
  setPasswordChangeRequired(value: boolean): void;

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
    activationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastLogin?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    passwordChanged?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    userName: string,
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
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
    passwordChangeRequired: boolean,
    sequence: number,
  }
}

export class UserProfile extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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
    userName: string,
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: Gender,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class UpdateUserProfileRequest extends jspb.Message {
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

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserProfileRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserProfileRequest): UpdateUserProfileRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserProfileRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserProfileRequest;
  static deserializeBinaryFromReader(message: UpdateUserProfileRequest, reader: jspb.BinaryReader): UpdateUserProfileRequest;
}

export namespace UpdateUserProfileRequest {
  export type AsObject = {
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: Gender,
  }
}

export class UserEmail extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getEmail(): string;
  setEmail(value: string): void;

  getIsemailverified(): boolean;
  setIsemailverified(value: boolean): void;

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
    isemailverified: boolean,
    sequence: number,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class VerifyMyUserEmailRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyUserEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyUserEmailRequest): VerifyMyUserEmailRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyMyUserEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyUserEmailRequest;
  static deserializeBinaryFromReader(message: VerifyMyUserEmailRequest, reader: jspb.BinaryReader): VerifyMyUserEmailRequest;
}

export namespace VerifyMyUserEmailRequest {
  export type AsObject = {
    code: string,
  }
}

export class VerifyUserEmailRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getCode(): string;
  setCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyUserEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyUserEmailRequest): VerifyUserEmailRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyUserEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyUserEmailRequest;
  static deserializeBinaryFromReader(message: VerifyUserEmailRequest, reader: jspb.BinaryReader): VerifyUserEmailRequest;
}

export namespace VerifyUserEmailRequest {
  export type AsObject = {
    id: string,
    code: string,
  }
}

export class UpdateUserEmailRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserEmailRequest): UpdateUserEmailRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserEmailRequest;
  static deserializeBinaryFromReader(message: UpdateUserEmailRequest, reader: jspb.BinaryReader): UpdateUserEmailRequest;
}

export namespace UpdateUserEmailRequest {
  export type AsObject = {
    email: string,
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

export class UpdateUserPhoneRequest extends jspb.Message {
  getPhone(): string;
  setPhone(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserPhoneRequest): UpdateUserPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserPhoneRequest;
  static deserializeBinaryFromReader(message: UpdateUserPhoneRequest, reader: jspb.BinaryReader): UpdateUserPhoneRequest;
}

export namespace UpdateUserPhoneRequest {
  export type AsObject = {
    phone: string,
  }
}

export class VerifyUserPhoneRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyUserPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyUserPhoneRequest): VerifyUserPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyUserPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyUserPhoneRequest;
  static deserializeBinaryFromReader(message: VerifyUserPhoneRequest, reader: jspb.BinaryReader): VerifyUserPhoneRequest;
}

export namespace VerifyUserPhoneRequest {
  export type AsObject = {
    code: string,
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

export class UpdateUserAddressRequest extends jspb.Message {
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
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
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
    password: string,
  }
}

export class PasswordChange extends jspb.Message {
  getOldPassword(): string;
  setOldPassword(value: string): void;

  getNewPassword(): string;
  setNewPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordChange.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordChange): PasswordChange.AsObject;
  static serializeBinaryToWriter(message: PasswordChange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordChange;
  static deserializeBinaryFromReader(message: PasswordChange, reader: jspb.BinaryReader): PasswordChange;
}

export namespace PasswordChange {
  export type AsObject = {
    oldPassword: string,
    newPassword: string,
  }
}

export class VerifyMfaOtp extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMfaOtp.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMfaOtp): VerifyMfaOtp.AsObject;
  static serializeBinaryToWriter(message: VerifyMfaOtp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMfaOtp;
  static deserializeBinaryFromReader(message: VerifyMfaOtp, reader: jspb.BinaryReader): VerifyMfaOtp;
}

export namespace VerifyMfaOtp {
  export type AsObject = {
    code: string,
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

export class MfaOtpResponse extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getUrl(): string;
  setUrl(value: string): void;

  getSecret(): string;
  setSecret(value: string): void;

  getState(): MFAState;
  setState(value: MFAState): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MfaOtpResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MfaOtpResponse): MfaOtpResponse.AsObject;
  static serializeBinaryToWriter(message: MfaOtpResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MfaOtpResponse;
  static deserializeBinaryFromReader(message: MfaOtpResponse, reader: jspb.BinaryReader): MfaOtpResponse;
}

export namespace MfaOtpResponse {
  export type AsObject = {
    userId: string,
    url: string,
    secret: string,
    state: MFAState,
  }
}

export class OIDCClientAuth extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): void;

  getClientSecret(): string;
  setClientSecret(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OIDCClientAuth.AsObject;
  static toObject(includeInstance: boolean, msg: OIDCClientAuth): OIDCClientAuth.AsObject;
  static serializeBinaryToWriter(message: OIDCClientAuth, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OIDCClientAuth;
  static deserializeBinaryFromReader(message: OIDCClientAuth, reader: jspb.BinaryReader): OIDCClientAuth;
}

export namespace OIDCClientAuth {
  export type AsObject = {
    clientId: string,
    clientSecret: string,
  }
}

export class UserGrantSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getSortingColumn(): UserGrantSearchKey;
  setSortingColumn(value: UserGrantSearchKey): void;

  getAsc(): boolean;
  setAsc(value: boolean): void;

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
    sortingColumn: UserGrantSearchKey,
    asc: boolean,
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

export class UserGrantView extends jspb.Message {
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

  getOrgname(): string;
  setOrgname(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantView.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantView): UserGrantView.AsObject;
  static serializeBinaryToWriter(message: UserGrantView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantView;
  static deserializeBinaryFromReader(message: UserGrantView, reader: jspb.BinaryReader): UserGrantView;
}

export namespace UserGrantView {
  export type AsObject = {
    orgid: string,
    projectid: string,
    userid: string,
    rolesList: Array<string>,
    orgname: string,
  }
}

export class MyProjectOrgSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getAsc(): boolean;
  setAsc(value: boolean): void;

  getQueriesList(): Array<MyProjectOrgSearchQuery>;
  setQueriesList(value: Array<MyProjectOrgSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: MyProjectOrgSearchQuery, index?: number): MyProjectOrgSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MyProjectOrgSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MyProjectOrgSearchRequest): MyProjectOrgSearchRequest.AsObject;
  static serializeBinaryToWriter(message: MyProjectOrgSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MyProjectOrgSearchRequest;
  static deserializeBinaryFromReader(message: MyProjectOrgSearchRequest, reader: jspb.BinaryReader): MyProjectOrgSearchRequest;
}

export namespace MyProjectOrgSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    asc: boolean,
    queriesList: Array<MyProjectOrgSearchQuery.AsObject>,
  }
}

export class MyProjectOrgSearchQuery extends jspb.Message {
  getKey(): MyProjectOrgSearchKey;
  setKey(value: MyProjectOrgSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MyProjectOrgSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MyProjectOrgSearchQuery): MyProjectOrgSearchQuery.AsObject;
  static serializeBinaryToWriter(message: MyProjectOrgSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MyProjectOrgSearchQuery;
  static deserializeBinaryFromReader(message: MyProjectOrgSearchQuery, reader: jspb.BinaryReader): MyProjectOrgSearchQuery;
}

export namespace MyProjectOrgSearchQuery {
  export type AsObject = {
    key: MyProjectOrgSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class MyProjectOrgSearchResponse extends jspb.Message {
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
  toObject(includeInstance?: boolean): MyProjectOrgSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MyProjectOrgSearchResponse): MyProjectOrgSearchResponse.AsObject;
  static serializeBinaryToWriter(message: MyProjectOrgSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MyProjectOrgSearchResponse;
  static deserializeBinaryFromReader(message: MyProjectOrgSearchResponse, reader: jspb.BinaryReader): MyProjectOrgSearchResponse;
}

export namespace MyProjectOrgSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<Org.AsObject>,
  }
}

export class Org extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

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
    name: string,
  }
}

export class MyPermissions extends jspb.Message {
  getPermissionsList(): Array<string>;
  setPermissionsList(value: Array<string>): void;
  clearPermissionsList(): void;
  addPermissions(value: string, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MyPermissions.AsObject;
  static toObject(includeInstance: boolean, msg: MyPermissions): MyPermissions.AsObject;
  static serializeBinaryToWriter(message: MyPermissions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MyPermissions;
  static deserializeBinaryFromReader(message: MyPermissions, reader: jspb.BinaryReader): MyPermissions;
}

export namespace MyPermissions {
  export type AsObject = {
    permissionsList: Array<string>,
  }
}

export enum UserSessionState { 
  USERSESSIONSTATE_UNSPECIFIED = 0,
  USERSESSIONSTATE_ACTIVE = 1,
  USERSESSIONSTATE_TERMINATED = 2,
}
export enum OIDCResponseType { 
  OIDCRESPONSETYPE_CODE = 0,
  OIDCRESPONSETYPE_ID_TOKEN = 1,
  OIDCRESPONSETYPE_ID_TOKEN_TOKEN = 2,
}
export enum UserState { 
  USERSTATE_UNSPECIEFIED = 0,
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
export enum UserGrantSearchKey { 
  USERGRANTSEARCHKEY_UNKNOWN = 0,
  USERGRANTSEARCHKEY_ORG_ID = 1,
  USERGRANTSEARCHKEY_PROJECT_ID = 2,
}
export enum MyProjectOrgSearchKey { 
  MYPROJECTORGSEARCHKEY_UNSPECIFIED = 0,
  MYPROJECTORGSEARCHKEY_ORG_NAME = 1,
}
export enum SearchMethod { 
  SEARCHMETHOD_EQUALS = 0,
  SEARCHMETHOD_STARTS_WITH = 1,
  SEARCHMETHOD_CONTAINS = 2,
  SEARCHMETHOD_EQUALS_IGNORE_CASE = 3,
  SEARCHMETHOD_STARTS_WITH_IGNORE_CASE = 4,
  SEARCHMETHOD_CONTAINS_IGNORE_CASE = 5,
}

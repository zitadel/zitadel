import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class User extends jspb.Message {
  getId(): string;
  setId(value: string): User;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): User;
  hasDetails(): boolean;
  clearDetails(): User;

  getState(): UserState;
  setState(value: UserState): User;

  getUserName(): string;
  setUserName(value: string): User;

  getLoginNamesList(): Array<string>;
  setLoginNamesList(value: Array<string>): User;
  clearLoginNamesList(): User;
  addLoginNames(value: string, index?: number): User;

  getPreferredLoginName(): string;
  setPreferredLoginName(value: string): User;

  getHuman(): Human | undefined;
  setHuman(value?: Human): User;
  hasHuman(): boolean;
  clearHuman(): User;

  getMachine(): Machine | undefined;
  setMachine(value?: Machine): User;
  hasMachine(): boolean;
  clearMachine(): User;

  getTypeCase(): User.TypeCase;

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
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: UserState,
    userName: string,
    loginNamesList: Array<string>,
    preferredLoginName: string,
    human?: Human.AsObject,
    machine?: Machine.AsObject,
  }

  export enum TypeCase { 
    TYPE_NOT_SET = 0,
    HUMAN = 7,
    MACHINE = 8,
  }
}

export class Human extends jspb.Message {
  getProfile(): Profile | undefined;
  setProfile(value?: Profile): Human;
  hasProfile(): boolean;
  clearProfile(): Human;

  getEmail(): Email | undefined;
  setEmail(value?: Email): Human;
  hasEmail(): boolean;
  clearEmail(): Human;

  getPhone(): Phone | undefined;
  setPhone(value?: Phone): Human;
  hasPhone(): boolean;
  clearPhone(): Human;

  getPasswordChanged(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setPasswordChanged(value?: google_protobuf_timestamp_pb.Timestamp): Human;
  hasPasswordChanged(): boolean;
  clearPasswordChanged(): Human;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Human.AsObject;
  static toObject(includeInstance: boolean, msg: Human): Human.AsObject;
  static serializeBinaryToWriter(message: Human, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Human;
  static deserializeBinaryFromReader(message: Human, reader: jspb.BinaryReader): Human;
}

export namespace Human {
  export type AsObject = {
    profile?: Profile.AsObject,
    email?: Email.AsObject,
    phone?: Phone.AsObject,
    passwordChanged?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class Machine extends jspb.Message {
  getName(): string;
  setName(value: string): Machine;

  getDescription(): string;
  setDescription(value: string): Machine;

  getHasSecret(): boolean;
  setHasSecret(value: boolean): Machine;

  getAccessTokenType(): AccessTokenType;
  setAccessTokenType(value: AccessTokenType): Machine;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Machine.AsObject;
  static toObject(includeInstance: boolean, msg: Machine): Machine.AsObject;
  static serializeBinaryToWriter(message: Machine, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Machine;
  static deserializeBinaryFromReader(message: Machine, reader: jspb.BinaryReader): Machine;
}

export namespace Machine {
  export type AsObject = {
    name: string,
    description: string,
    hasSecret: boolean,
    accessTokenType: AccessTokenType,
  }
}

export class Profile extends jspb.Message {
  getFirstName(): string;
  setFirstName(value: string): Profile;

  getLastName(): string;
  setLastName(value: string): Profile;

  getNickName(): string;
  setNickName(value: string): Profile;

  getDisplayName(): string;
  setDisplayName(value: string): Profile;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): Profile;

  getGender(): Gender;
  setGender(value: Gender): Profile;

  getAvatarUrl(): string;
  setAvatarUrl(value: string): Profile;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Profile.AsObject;
  static toObject(includeInstance: boolean, msg: Profile): Profile.AsObject;
  static serializeBinaryToWriter(message: Profile, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Profile;
  static deserializeBinaryFromReader(message: Profile, reader: jspb.BinaryReader): Profile;
}

export namespace Profile {
  export type AsObject = {
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: Gender,
    avatarUrl: string,
  }
}

export class Email extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): Email;

  getIsEmailVerified(): boolean;
  setIsEmailVerified(value: boolean): Email;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Email.AsObject;
  static toObject(includeInstance: boolean, msg: Email): Email.AsObject;
  static serializeBinaryToWriter(message: Email, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Email;
  static deserializeBinaryFromReader(message: Email, reader: jspb.BinaryReader): Email;
}

export namespace Email {
  export type AsObject = {
    email: string,
    isEmailVerified: boolean,
  }
}

export class Phone extends jspb.Message {
  getPhone(): string;
  setPhone(value: string): Phone;

  getIsPhoneVerified(): boolean;
  setIsPhoneVerified(value: boolean): Phone;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Phone.AsObject;
  static toObject(includeInstance: boolean, msg: Phone): Phone.AsObject;
  static serializeBinaryToWriter(message: Phone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Phone;
  static deserializeBinaryFromReader(message: Phone, reader: jspb.BinaryReader): Phone;
}

export namespace Phone {
  export type AsObject = {
    phone: string,
    isPhoneVerified: boolean,
  }
}

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

export class UserNameQuery extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): UserNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserNameQuery;

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
    method: zitadel_object_pb.TextQueryMethod,
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

export class NickNameQuery extends jspb.Message {
  getNickName(): string;
  setNickName(value: string): NickNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): NickNameQuery;

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
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class DisplayNameQuery extends jspb.Message {
  getDisplayName(): string;
  setDisplayName(value: string): DisplayNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): DisplayNameQuery;

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
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class EmailQuery extends jspb.Message {
  getEmailAddress(): string;
  setEmailAddress(value: string): EmailQuery;

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
    emailAddress: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class LoginNameQuery extends jspb.Message {
  getLoginName(): string;
  setLoginName(value: string): LoginNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): LoginNameQuery;

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
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class StateQuery extends jspb.Message {
  getState(): UserState;
  setState(value: UserState): StateQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StateQuery.AsObject;
  static toObject(includeInstance: boolean, msg: StateQuery): StateQuery.AsObject;
  static serializeBinaryToWriter(message: StateQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StateQuery;
  static deserializeBinaryFromReader(message: StateQuery, reader: jspb.BinaryReader): StateQuery;
}

export namespace StateQuery {
  export type AsObject = {
    state: UserState,
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

export class AuthFactor extends jspb.Message {
  getState(): AuthFactorState;
  setState(value: AuthFactorState): AuthFactor;

  getOtp(): AuthFactorOTP | undefined;
  setOtp(value?: AuthFactorOTP): AuthFactor;
  hasOtp(): boolean;
  clearOtp(): AuthFactor;

  getU2f(): AuthFactorU2F | undefined;
  setU2f(value?: AuthFactorU2F): AuthFactor;
  hasU2f(): boolean;
  clearU2f(): AuthFactor;

  getOtpSms(): AuthFactorOTPSMS | undefined;
  setOtpSms(value?: AuthFactorOTPSMS): AuthFactor;
  hasOtpSms(): boolean;
  clearOtpSms(): AuthFactor;

  getOtpEmail(): AuthFactorOTPEmail | undefined;
  setOtpEmail(value?: AuthFactorOTPEmail): AuthFactor;
  hasOtpEmail(): boolean;
  clearOtpEmail(): AuthFactor;

  getTypeCase(): AuthFactor.TypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthFactor.AsObject;
  static toObject(includeInstance: boolean, msg: AuthFactor): AuthFactor.AsObject;
  static serializeBinaryToWriter(message: AuthFactor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthFactor;
  static deserializeBinaryFromReader(message: AuthFactor, reader: jspb.BinaryReader): AuthFactor;
}

export namespace AuthFactor {
  export type AsObject = {
    state: AuthFactorState,
    otp?: AuthFactorOTP.AsObject,
    u2f?: AuthFactorU2F.AsObject,
    otpSms?: AuthFactorOTPSMS.AsObject,
    otpEmail?: AuthFactorOTPEmail.AsObject,
  }

  export enum TypeCase { 
    TYPE_NOT_SET = 0,
    OTP = 2,
    U2F = 3,
    OTP_SMS = 4,
    OTP_EMAIL = 5,
  }
}

export class AuthFactorOTP extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthFactorOTP.AsObject;
  static toObject(includeInstance: boolean, msg: AuthFactorOTP): AuthFactorOTP.AsObject;
  static serializeBinaryToWriter(message: AuthFactorOTP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthFactorOTP;
  static deserializeBinaryFromReader(message: AuthFactorOTP, reader: jspb.BinaryReader): AuthFactorOTP;
}

export namespace AuthFactorOTP {
  export type AsObject = {
  }
}

export class AuthFactorOTPSMS extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthFactorOTPSMS.AsObject;
  static toObject(includeInstance: boolean, msg: AuthFactorOTPSMS): AuthFactorOTPSMS.AsObject;
  static serializeBinaryToWriter(message: AuthFactorOTPSMS, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthFactorOTPSMS;
  static deserializeBinaryFromReader(message: AuthFactorOTPSMS, reader: jspb.BinaryReader): AuthFactorOTPSMS;
}

export namespace AuthFactorOTPSMS {
  export type AsObject = {
  }
}

export class AuthFactorOTPEmail extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthFactorOTPEmail.AsObject;
  static toObject(includeInstance: boolean, msg: AuthFactorOTPEmail): AuthFactorOTPEmail.AsObject;
  static serializeBinaryToWriter(message: AuthFactorOTPEmail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthFactorOTPEmail;
  static deserializeBinaryFromReader(message: AuthFactorOTPEmail, reader: jspb.BinaryReader): AuthFactorOTPEmail;
}

export namespace AuthFactorOTPEmail {
  export type AsObject = {
  }
}

export class AuthFactorU2F extends jspb.Message {
  getId(): string;
  setId(value: string): AuthFactorU2F;

  getName(): string;
  setName(value: string): AuthFactorU2F;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthFactorU2F.AsObject;
  static toObject(includeInstance: boolean, msg: AuthFactorU2F): AuthFactorU2F.AsObject;
  static serializeBinaryToWriter(message: AuthFactorU2F, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthFactorU2F;
  static deserializeBinaryFromReader(message: AuthFactorU2F, reader: jspb.BinaryReader): AuthFactorU2F;
}

export namespace AuthFactorU2F {
  export type AsObject = {
    id: string,
    name: string,
  }
}

export class WebAuthNKey extends jspb.Message {
  getPublicKey(): Uint8Array | string;
  getPublicKey_asU8(): Uint8Array;
  getPublicKey_asB64(): string;
  setPublicKey(value: Uint8Array | string): WebAuthNKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebAuthNKey.AsObject;
  static toObject(includeInstance: boolean, msg: WebAuthNKey): WebAuthNKey.AsObject;
  static serializeBinaryToWriter(message: WebAuthNKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebAuthNKey;
  static deserializeBinaryFromReader(message: WebAuthNKey, reader: jspb.BinaryReader): WebAuthNKey;
}

export namespace WebAuthNKey {
  export type AsObject = {
    publicKey: Uint8Array | string,
  }
}

export class WebAuthNVerification extends jspb.Message {
  getPublicKeyCredential(): Uint8Array | string;
  getPublicKeyCredential_asU8(): Uint8Array;
  getPublicKeyCredential_asB64(): string;
  setPublicKeyCredential(value: Uint8Array | string): WebAuthNVerification;

  getTokenName(): string;
  setTokenName(value: string): WebAuthNVerification;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebAuthNVerification.AsObject;
  static toObject(includeInstance: boolean, msg: WebAuthNVerification): WebAuthNVerification.AsObject;
  static serializeBinaryToWriter(message: WebAuthNVerification, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebAuthNVerification;
  static deserializeBinaryFromReader(message: WebAuthNVerification, reader: jspb.BinaryReader): WebAuthNVerification;
}

export namespace WebAuthNVerification {
  export type AsObject = {
    publicKeyCredential: Uint8Array | string,
    tokenName: string,
  }
}

export class WebAuthNToken extends jspb.Message {
  getId(): string;
  setId(value: string): WebAuthNToken;

  getState(): AuthFactorState;
  setState(value: AuthFactorState): WebAuthNToken;

  getName(): string;
  setName(value: string): WebAuthNToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebAuthNToken.AsObject;
  static toObject(includeInstance: boolean, msg: WebAuthNToken): WebAuthNToken.AsObject;
  static serializeBinaryToWriter(message: WebAuthNToken, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebAuthNToken;
  static deserializeBinaryFromReader(message: WebAuthNToken, reader: jspb.BinaryReader): WebAuthNToken;
}

export namespace WebAuthNToken {
  export type AsObject = {
    id: string,
    state: AuthFactorState,
    name: string,
  }
}

export class Membership extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): Membership;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Membership;
  hasDetails(): boolean;
  clearDetails(): Membership;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): Membership;
  clearRolesList(): Membership;
  addRoles(value: string, index?: number): Membership;

  getDisplayName(): string;
  setDisplayName(value: string): Membership;

  getIam(): boolean;
  setIam(value: boolean): Membership;

  getOrgId(): string;
  setOrgId(value: string): Membership;

  getProjectId(): string;
  setProjectId(value: string): Membership;

  getProjectGrantId(): string;
  setProjectGrantId(value: string): Membership;

  getTypeCase(): Membership.TypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Membership.AsObject;
  static toObject(includeInstance: boolean, msg: Membership): Membership.AsObject;
  static serializeBinaryToWriter(message: Membership, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Membership;
  static deserializeBinaryFromReader(message: Membership, reader: jspb.BinaryReader): Membership;
}

export namespace Membership {
  export type AsObject = {
    userId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    rolesList: Array<string>,
    displayName: string,
    iam: boolean,
    orgId: string,
    projectId: string,
    projectGrantId: string,
  }

  export enum TypeCase { 
    TYPE_NOT_SET = 0,
    IAM = 5,
    ORG_ID = 6,
    PROJECT_ID = 7,
    PROJECT_GRANT_ID = 8,
  }
}

export class MembershipQuery extends jspb.Message {
  getOrgQuery(): MembershipOrgQuery | undefined;
  setOrgQuery(value?: MembershipOrgQuery): MembershipQuery;
  hasOrgQuery(): boolean;
  clearOrgQuery(): MembershipQuery;

  getProjectQuery(): MembershipProjectQuery | undefined;
  setProjectQuery(value?: MembershipProjectQuery): MembershipQuery;
  hasProjectQuery(): boolean;
  clearProjectQuery(): MembershipQuery;

  getProjectGrantQuery(): MembershipProjectGrantQuery | undefined;
  setProjectGrantQuery(value?: MembershipProjectGrantQuery): MembershipQuery;
  hasProjectGrantQuery(): boolean;
  clearProjectGrantQuery(): MembershipQuery;

  getIamQuery(): MembershipIAMQuery | undefined;
  setIamQuery(value?: MembershipIAMQuery): MembershipQuery;
  hasIamQuery(): boolean;
  clearIamQuery(): MembershipQuery;

  getQueryCase(): MembershipQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MembershipQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MembershipQuery): MembershipQuery.AsObject;
  static serializeBinaryToWriter(message: MembershipQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MembershipQuery;
  static deserializeBinaryFromReader(message: MembershipQuery, reader: jspb.BinaryReader): MembershipQuery;
}

export namespace MembershipQuery {
  export type AsObject = {
    orgQuery?: MembershipOrgQuery.AsObject,
    projectQuery?: MembershipProjectQuery.AsObject,
    projectGrantQuery?: MembershipProjectGrantQuery.AsObject,
    iamQuery?: MembershipIAMQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    ORG_QUERY = 1,
    PROJECT_QUERY = 2,
    PROJECT_GRANT_QUERY = 3,
    IAM_QUERY = 4,
  }
}

export class MembershipOrgQuery extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): MembershipOrgQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MembershipOrgQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MembershipOrgQuery): MembershipOrgQuery.AsObject;
  static serializeBinaryToWriter(message: MembershipOrgQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MembershipOrgQuery;
  static deserializeBinaryFromReader(message: MembershipOrgQuery, reader: jspb.BinaryReader): MembershipOrgQuery;
}

export namespace MembershipOrgQuery {
  export type AsObject = {
    orgId: string,
  }
}

export class MembershipProjectQuery extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): MembershipProjectQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MembershipProjectQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MembershipProjectQuery): MembershipProjectQuery.AsObject;
  static serializeBinaryToWriter(message: MembershipProjectQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MembershipProjectQuery;
  static deserializeBinaryFromReader(message: MembershipProjectQuery, reader: jspb.BinaryReader): MembershipProjectQuery;
}

export namespace MembershipProjectQuery {
  export type AsObject = {
    projectId: string,
  }
}

export class MembershipProjectGrantQuery extends jspb.Message {
  getProjectGrantId(): string;
  setProjectGrantId(value: string): MembershipProjectGrantQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MembershipProjectGrantQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MembershipProjectGrantQuery): MembershipProjectGrantQuery.AsObject;
  static serializeBinaryToWriter(message: MembershipProjectGrantQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MembershipProjectGrantQuery;
  static deserializeBinaryFromReader(message: MembershipProjectGrantQuery, reader: jspb.BinaryReader): MembershipProjectGrantQuery;
}

export namespace MembershipProjectGrantQuery {
  export type AsObject = {
    projectGrantId: string,
  }
}

export class MembershipIAMQuery extends jspb.Message {
  getIam(): boolean;
  setIam(value: boolean): MembershipIAMQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MembershipIAMQuery.AsObject;
  static toObject(includeInstance: boolean, msg: MembershipIAMQuery): MembershipIAMQuery.AsObject;
  static serializeBinaryToWriter(message: MembershipIAMQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MembershipIAMQuery;
  static deserializeBinaryFromReader(message: MembershipIAMQuery, reader: jspb.BinaryReader): MembershipIAMQuery;
}

export namespace MembershipIAMQuery {
  export type AsObject = {
    iam: boolean,
  }
}

export class Session extends jspb.Message {
  getSessionId(): string;
  setSessionId(value: string): Session;

  getAgentId(): string;
  setAgentId(value: string): Session;

  getAuthState(): SessionState;
  setAuthState(value: SessionState): Session;

  getUserId(): string;
  setUserId(value: string): Session;

  getUserName(): string;
  setUserName(value: string): Session;

  getLoginName(): string;
  setLoginName(value: string): Session;

  getDisplayName(): string;
  setDisplayName(value: string): Session;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Session;
  hasDetails(): boolean;
  clearDetails(): Session;

  getAvatarUrl(): string;
  setAvatarUrl(value: string): Session;

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
    agentId: string,
    authState: SessionState,
    userId: string,
    userName: string,
    loginName: string,
    displayName: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    avatarUrl: string,
  }
}

export class RefreshToken extends jspb.Message {
  getId(): string;
  setId(value: string): RefreshToken;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RefreshToken;
  hasDetails(): boolean;
  clearDetails(): RefreshToken;

  getClientId(): string;
  setClientId(value: string): RefreshToken;

  getAuthTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setAuthTime(value?: google_protobuf_timestamp_pb.Timestamp): RefreshToken;
  hasAuthTime(): boolean;
  clearAuthTime(): RefreshToken;

  getIdleExpiration(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setIdleExpiration(value?: google_protobuf_timestamp_pb.Timestamp): RefreshToken;
  hasIdleExpiration(): boolean;
  clearIdleExpiration(): RefreshToken;

  getExpiration(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpiration(value?: google_protobuf_timestamp_pb.Timestamp): RefreshToken;
  hasExpiration(): boolean;
  clearExpiration(): RefreshToken;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): RefreshToken;
  clearScopesList(): RefreshToken;
  addScopes(value: string, index?: number): RefreshToken;

  getAudienceList(): Array<string>;
  setAudienceList(value: Array<string>): RefreshToken;
  clearAudienceList(): RefreshToken;
  addAudience(value: string, index?: number): RefreshToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshToken.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshToken): RefreshToken.AsObject;
  static serializeBinaryToWriter(message: RefreshToken, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshToken;
  static deserializeBinaryFromReader(message: RefreshToken, reader: jspb.BinaryReader): RefreshToken;
}

export namespace RefreshToken {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    clientId: string,
    authTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    idleExpiration?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    expiration?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    scopesList: Array<string>,
    audienceList: Array<string>,
  }
}

export class PersonalAccessToken extends jspb.Message {
  getId(): string;
  setId(value: string): PersonalAccessToken;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): PersonalAccessToken;
  hasDetails(): boolean;
  clearDetails(): PersonalAccessToken;

  getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): PersonalAccessToken;
  hasExpirationDate(): boolean;
  clearExpirationDate(): PersonalAccessToken;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): PersonalAccessToken;
  clearScopesList(): PersonalAccessToken;
  addScopes(value: string, index?: number): PersonalAccessToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PersonalAccessToken.AsObject;
  static toObject(includeInstance: boolean, msg: PersonalAccessToken): PersonalAccessToken.AsObject;
  static serializeBinaryToWriter(message: PersonalAccessToken, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PersonalAccessToken;
  static deserializeBinaryFromReader(message: PersonalAccessToken, reader: jspb.BinaryReader): PersonalAccessToken;
}

export namespace PersonalAccessToken {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    scopesList: Array<string>,
  }
}

export class UserGrant extends jspb.Message {
  getId(): string;
  setId(value: string): UserGrant;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UserGrant;
  hasDetails(): boolean;
  clearDetails(): UserGrant;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): UserGrant;
  clearRoleKeysList(): UserGrant;
  addRoleKeys(value: string, index?: number): UserGrant;

  getState(): UserGrantState;
  setState(value: UserGrantState): UserGrant;

  getUserId(): string;
  setUserId(value: string): UserGrant;

  getUserName(): string;
  setUserName(value: string): UserGrant;

  getFirstName(): string;
  setFirstName(value: string): UserGrant;

  getLastName(): string;
  setLastName(value: string): UserGrant;

  getEmail(): string;
  setEmail(value: string): UserGrant;

  getDisplayName(): string;
  setDisplayName(value: string): UserGrant;

  getOrgId(): string;
  setOrgId(value: string): UserGrant;

  getOrgName(): string;
  setOrgName(value: string): UserGrant;

  getOrgDomain(): string;
  setOrgDomain(value: string): UserGrant;

  getProjectId(): string;
  setProjectId(value: string): UserGrant;

  getProjectName(): string;
  setProjectName(value: string): UserGrant;

  getProjectGrantId(): string;
  setProjectGrantId(value: string): UserGrant;

  getAvatarUrl(): string;
  setAvatarUrl(value: string): UserGrant;

  getPreferredLoginName(): string;
  setPreferredLoginName(value: string): UserGrant;

  getUserType(): Type;
  setUserType(value: Type): UserGrant;

  getGrantedOrgId(): string;
  setGrantedOrgId(value: string): UserGrant;

  getGrantedOrgName(): string;
  setGrantedOrgName(value: string): UserGrant;

  getGrantedOrgDomain(): string;
  setGrantedOrgDomain(value: string): UserGrant;

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
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    roleKeysList: Array<string>,
    state: UserGrantState,
    userId: string,
    userName: string,
    firstName: string,
    lastName: string,
    email: string,
    displayName: string,
    orgId: string,
    orgName: string,
    orgDomain: string,
    projectId: string,
    projectName: string,
    projectGrantId: string,
    avatarUrl: string,
    preferredLoginName: string,
    userType: Type,
    grantedOrgId: string,
    grantedOrgName: string,
    grantedOrgDomain: string,
  }
}

export class UserGrantQuery extends jspb.Message {
  getProjectIdQuery(): UserGrantProjectIDQuery | undefined;
  setProjectIdQuery(value?: UserGrantProjectIDQuery): UserGrantQuery;
  hasProjectIdQuery(): boolean;
  clearProjectIdQuery(): UserGrantQuery;

  getUserIdQuery(): UserGrantUserIDQuery | undefined;
  setUserIdQuery(value?: UserGrantUserIDQuery): UserGrantQuery;
  hasUserIdQuery(): boolean;
  clearUserIdQuery(): UserGrantQuery;

  getWithGrantedQuery(): UserGrantWithGrantedQuery | undefined;
  setWithGrantedQuery(value?: UserGrantWithGrantedQuery): UserGrantQuery;
  hasWithGrantedQuery(): boolean;
  clearWithGrantedQuery(): UserGrantQuery;

  getRoleKeyQuery(): UserGrantRoleKeyQuery | undefined;
  setRoleKeyQuery(value?: UserGrantRoleKeyQuery): UserGrantQuery;
  hasRoleKeyQuery(): boolean;
  clearRoleKeyQuery(): UserGrantQuery;

  getProjectGrantIdQuery(): UserGrantProjectGrantIDQuery | undefined;
  setProjectGrantIdQuery(value?: UserGrantProjectGrantIDQuery): UserGrantQuery;
  hasProjectGrantIdQuery(): boolean;
  clearProjectGrantIdQuery(): UserGrantQuery;

  getUserNameQuery(): UserGrantUserNameQuery | undefined;
  setUserNameQuery(value?: UserGrantUserNameQuery): UserGrantQuery;
  hasUserNameQuery(): boolean;
  clearUserNameQuery(): UserGrantQuery;

  getFirstNameQuery(): UserGrantFirstNameQuery | undefined;
  setFirstNameQuery(value?: UserGrantFirstNameQuery): UserGrantQuery;
  hasFirstNameQuery(): boolean;
  clearFirstNameQuery(): UserGrantQuery;

  getLastNameQuery(): UserGrantLastNameQuery | undefined;
  setLastNameQuery(value?: UserGrantLastNameQuery): UserGrantQuery;
  hasLastNameQuery(): boolean;
  clearLastNameQuery(): UserGrantQuery;

  getEmailQuery(): UserGrantEmailQuery | undefined;
  setEmailQuery(value?: UserGrantEmailQuery): UserGrantQuery;
  hasEmailQuery(): boolean;
  clearEmailQuery(): UserGrantQuery;

  getOrgNameQuery(): UserGrantOrgNameQuery | undefined;
  setOrgNameQuery(value?: UserGrantOrgNameQuery): UserGrantQuery;
  hasOrgNameQuery(): boolean;
  clearOrgNameQuery(): UserGrantQuery;

  getOrgDomainQuery(): UserGrantOrgDomainQuery | undefined;
  setOrgDomainQuery(value?: UserGrantOrgDomainQuery): UserGrantQuery;
  hasOrgDomainQuery(): boolean;
  clearOrgDomainQuery(): UserGrantQuery;

  getProjectNameQuery(): UserGrantProjectNameQuery | undefined;
  setProjectNameQuery(value?: UserGrantProjectNameQuery): UserGrantQuery;
  hasProjectNameQuery(): boolean;
  clearProjectNameQuery(): UserGrantQuery;

  getDisplayNameQuery(): UserGrantDisplayNameQuery | undefined;
  setDisplayNameQuery(value?: UserGrantDisplayNameQuery): UserGrantQuery;
  hasDisplayNameQuery(): boolean;
  clearDisplayNameQuery(): UserGrantQuery;

  getUserTypeQuery(): UserGrantUserTypeQuery | undefined;
  setUserTypeQuery(value?: UserGrantUserTypeQuery): UserGrantQuery;
  hasUserTypeQuery(): boolean;
  clearUserTypeQuery(): UserGrantQuery;

  getQueryCase(): UserGrantQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantQuery): UserGrantQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantQuery;
  static deserializeBinaryFromReader(message: UserGrantQuery, reader: jspb.BinaryReader): UserGrantQuery;
}

export namespace UserGrantQuery {
  export type AsObject = {
    projectIdQuery?: UserGrantProjectIDQuery.AsObject,
    userIdQuery?: UserGrantUserIDQuery.AsObject,
    withGrantedQuery?: UserGrantWithGrantedQuery.AsObject,
    roleKeyQuery?: UserGrantRoleKeyQuery.AsObject,
    projectGrantIdQuery?: UserGrantProjectGrantIDQuery.AsObject,
    userNameQuery?: UserGrantUserNameQuery.AsObject,
    firstNameQuery?: UserGrantFirstNameQuery.AsObject,
    lastNameQuery?: UserGrantLastNameQuery.AsObject,
    emailQuery?: UserGrantEmailQuery.AsObject,
    orgNameQuery?: UserGrantOrgNameQuery.AsObject,
    orgDomainQuery?: UserGrantOrgDomainQuery.AsObject,
    projectNameQuery?: UserGrantProjectNameQuery.AsObject,
    displayNameQuery?: UserGrantDisplayNameQuery.AsObject,
    userTypeQuery?: UserGrantUserTypeQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    PROJECT_ID_QUERY = 1,
    USER_ID_QUERY = 2,
    WITH_GRANTED_QUERY = 3,
    ROLE_KEY_QUERY = 4,
    PROJECT_GRANT_ID_QUERY = 5,
    USER_NAME_QUERY = 6,
    FIRST_NAME_QUERY = 7,
    LAST_NAME_QUERY = 8,
    EMAIL_QUERY = 9,
    ORG_NAME_QUERY = 10,
    ORG_DOMAIN_QUERY = 11,
    PROJECT_NAME_QUERY = 12,
    DISPLAY_NAME_QUERY = 13,
    USER_TYPE_QUERY = 14,
  }
}

export class UserGrantProjectIDQuery extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UserGrantProjectIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantProjectIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantProjectIDQuery): UserGrantProjectIDQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantProjectIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantProjectIDQuery;
  static deserializeBinaryFromReader(message: UserGrantProjectIDQuery, reader: jspb.BinaryReader): UserGrantProjectIDQuery;
}

export namespace UserGrantProjectIDQuery {
  export type AsObject = {
    projectId: string,
  }
}

export class UserGrantUserIDQuery extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UserGrantUserIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantUserIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantUserIDQuery): UserGrantUserIDQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantUserIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantUserIDQuery;
  static deserializeBinaryFromReader(message: UserGrantUserIDQuery, reader: jspb.BinaryReader): UserGrantUserIDQuery;
}

export namespace UserGrantUserIDQuery {
  export type AsObject = {
    userId: string,
  }
}

export class UserGrantWithGrantedQuery extends jspb.Message {
  getWithGranted(): boolean;
  setWithGranted(value: boolean): UserGrantWithGrantedQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantWithGrantedQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantWithGrantedQuery): UserGrantWithGrantedQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantWithGrantedQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantWithGrantedQuery;
  static deserializeBinaryFromReader(message: UserGrantWithGrantedQuery, reader: jspb.BinaryReader): UserGrantWithGrantedQuery;
}

export namespace UserGrantWithGrantedQuery {
  export type AsObject = {
    withGranted: boolean,
  }
}

export class UserGrantRoleKeyQuery extends jspb.Message {
  getRoleKey(): string;
  setRoleKey(value: string): UserGrantRoleKeyQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantRoleKeyQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantRoleKeyQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantRoleKeyQuery): UserGrantRoleKeyQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantRoleKeyQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantRoleKeyQuery;
  static deserializeBinaryFromReader(message: UserGrantRoleKeyQuery, reader: jspb.BinaryReader): UserGrantRoleKeyQuery;
}

export namespace UserGrantRoleKeyQuery {
  export type AsObject = {
    roleKey: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantProjectGrantIDQuery extends jspb.Message {
  getProjectGrantId(): string;
  setProjectGrantId(value: string): UserGrantProjectGrantIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantProjectGrantIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantProjectGrantIDQuery): UserGrantProjectGrantIDQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantProjectGrantIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantProjectGrantIDQuery;
  static deserializeBinaryFromReader(message: UserGrantProjectGrantIDQuery, reader: jspb.BinaryReader): UserGrantProjectGrantIDQuery;
}

export namespace UserGrantProjectGrantIDQuery {
  export type AsObject = {
    projectGrantId: string,
  }
}

export class UserGrantUserNameQuery extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): UserGrantUserNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantUserNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantUserNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantUserNameQuery): UserGrantUserNameQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantUserNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantUserNameQuery;
  static deserializeBinaryFromReader(message: UserGrantUserNameQuery, reader: jspb.BinaryReader): UserGrantUserNameQuery;
}

export namespace UserGrantUserNameQuery {
  export type AsObject = {
    userName: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantFirstNameQuery extends jspb.Message {
  getFirstName(): string;
  setFirstName(value: string): UserGrantFirstNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantFirstNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantFirstNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantFirstNameQuery): UserGrantFirstNameQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantFirstNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantFirstNameQuery;
  static deserializeBinaryFromReader(message: UserGrantFirstNameQuery, reader: jspb.BinaryReader): UserGrantFirstNameQuery;
}

export namespace UserGrantFirstNameQuery {
  export type AsObject = {
    firstName: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantLastNameQuery extends jspb.Message {
  getLastName(): string;
  setLastName(value: string): UserGrantLastNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantLastNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantLastNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantLastNameQuery): UserGrantLastNameQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantLastNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantLastNameQuery;
  static deserializeBinaryFromReader(message: UserGrantLastNameQuery, reader: jspb.BinaryReader): UserGrantLastNameQuery;
}

export namespace UserGrantLastNameQuery {
  export type AsObject = {
    lastName: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantEmailQuery extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): UserGrantEmailQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantEmailQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantEmailQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantEmailQuery): UserGrantEmailQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantEmailQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantEmailQuery;
  static deserializeBinaryFromReader(message: UserGrantEmailQuery, reader: jspb.BinaryReader): UserGrantEmailQuery;
}

export namespace UserGrantEmailQuery {
  export type AsObject = {
    email: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantOrgNameQuery extends jspb.Message {
  getOrgName(): string;
  setOrgName(value: string): UserGrantOrgNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantOrgNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantOrgNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantOrgNameQuery): UserGrantOrgNameQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantOrgNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantOrgNameQuery;
  static deserializeBinaryFromReader(message: UserGrantOrgNameQuery, reader: jspb.BinaryReader): UserGrantOrgNameQuery;
}

export namespace UserGrantOrgNameQuery {
  export type AsObject = {
    orgName: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantOrgDomainQuery extends jspb.Message {
  getOrgDomain(): string;
  setOrgDomain(value: string): UserGrantOrgDomainQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantOrgDomainQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantOrgDomainQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantOrgDomainQuery): UserGrantOrgDomainQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantOrgDomainQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantOrgDomainQuery;
  static deserializeBinaryFromReader(message: UserGrantOrgDomainQuery, reader: jspb.BinaryReader): UserGrantOrgDomainQuery;
}

export namespace UserGrantOrgDomainQuery {
  export type AsObject = {
    orgDomain: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantProjectNameQuery extends jspb.Message {
  getProjectName(): string;
  setProjectName(value: string): UserGrantProjectNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantProjectNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantProjectNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantProjectNameQuery): UserGrantProjectNameQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantProjectNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantProjectNameQuery;
  static deserializeBinaryFromReader(message: UserGrantProjectNameQuery, reader: jspb.BinaryReader): UserGrantProjectNameQuery;
}

export namespace UserGrantProjectNameQuery {
  export type AsObject = {
    projectName: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantDisplayNameQuery extends jspb.Message {
  getDisplayName(): string;
  setDisplayName(value: string): UserGrantDisplayNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): UserGrantDisplayNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantDisplayNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantDisplayNameQuery): UserGrantDisplayNameQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantDisplayNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantDisplayNameQuery;
  static deserializeBinaryFromReader(message: UserGrantDisplayNameQuery, reader: jspb.BinaryReader): UserGrantDisplayNameQuery;
}

export namespace UserGrantDisplayNameQuery {
  export type AsObject = {
    displayName: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class UserGrantUserTypeQuery extends jspb.Message {
  getType(): Type;
  setType(value: Type): UserGrantUserTypeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrantUserTypeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrantUserTypeQuery): UserGrantUserTypeQuery.AsObject;
  static serializeBinaryToWriter(message: UserGrantUserTypeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrantUserTypeQuery;
  static deserializeBinaryFromReader(message: UserGrantUserTypeQuery, reader: jspb.BinaryReader): UserGrantUserTypeQuery;
}

export namespace UserGrantUserTypeQuery {
  export type AsObject = {
    type: Type,
  }
}

export enum UserState { 
  USER_STATE_UNSPECIFIED = 0,
  USER_STATE_ACTIVE = 1,
  USER_STATE_INACTIVE = 2,
  USER_STATE_DELETED = 3,
  USER_STATE_LOCKED = 4,
  USER_STATE_SUSPEND = 5,
  USER_STATE_INITIAL = 6,
}
export enum Gender { 
  GENDER_UNSPECIFIED = 0,
  GENDER_FEMALE = 1,
  GENDER_MALE = 2,
  GENDER_DIVERSE = 3,
}
export enum AccessTokenType { 
  ACCESS_TOKEN_TYPE_BEARER = 0,
  ACCESS_TOKEN_TYPE_JWT = 1,
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
export enum AuthFactorState { 
  AUTH_FACTOR_STATE_UNSPECIFIED = 0,
  AUTH_FACTOR_STATE_NOT_READY = 1,
  AUTH_FACTOR_STATE_READY = 2,
  AUTH_FACTOR_STATE_REMOVED = 3,
}
export enum SessionState { 
  SESSION_STATE_UNSPECIFIED = 0,
  SESSION_STATE_ACTIVE = 1,
  SESSION_STATE_TERMINATED = 2,
}
export enum UserGrantState { 
  USER_GRANT_STATE_UNSPECIFIED = 0,
  USER_GRANT_STATE_ACTIVE = 1,
  USER_GRANT_STATE_INACTIVE = 2,
}

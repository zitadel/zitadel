import * as jspb from 'google-protobuf'

import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_protoc_gen_zitadel_v2_options_pb from '../../../zitadel/protoc_gen_zitadel/v2/options_pb'; // proto import: "zitadel/protoc_gen_zitadel/v2/options.proto"
import * as zitadel_session_v2beta_challenge_pb from '../../../zitadel/session/v2beta/challenge_pb'; // proto import: "zitadel/session/v2beta/challenge.proto"
import * as zitadel_session_v2beta_session_pb from '../../../zitadel/session/v2beta/session_pb'; // proto import: "zitadel/session/v2beta/session.proto"
import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class ListSessionsRequest extends jspb.Message {
  getQuery(): zitadel_object_v2beta_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_v2beta_object_pb.ListQuery): ListSessionsRequest;
  hasQuery(): boolean;
  clearQuery(): ListSessionsRequest;

  getQueriesList(): Array<zitadel_session_v2beta_session_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_session_v2beta_session_pb.SearchQuery>): ListSessionsRequest;
  clearQueriesList(): ListSessionsRequest;
  addQueries(value?: zitadel_session_v2beta_session_pb.SearchQuery, index?: number): zitadel_session_v2beta_session_pb.SearchQuery;

  getSortingColumn(): zitadel_session_v2beta_session_pb.SessionFieldName;
  setSortingColumn(value: zitadel_session_v2beta_session_pb.SessionFieldName): ListSessionsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSessionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSessionsRequest): ListSessionsRequest.AsObject;
  static serializeBinaryToWriter(message: ListSessionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSessionsRequest;
  static deserializeBinaryFromReader(message: ListSessionsRequest, reader: jspb.BinaryReader): ListSessionsRequest;
}

export namespace ListSessionsRequest {
  export type AsObject = {
    query?: zitadel_object_v2beta_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_session_v2beta_session_pb.SearchQuery.AsObject>,
    sortingColumn: zitadel_session_v2beta_session_pb.SessionFieldName,
  }
}

export class ListSessionsResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.ListDetails): ListSessionsResponse;
  hasDetails(): boolean;
  clearDetails(): ListSessionsResponse;

  getSessionsList(): Array<zitadel_session_v2beta_session_pb.Session>;
  setSessionsList(value: Array<zitadel_session_v2beta_session_pb.Session>): ListSessionsResponse;
  clearSessionsList(): ListSessionsResponse;
  addSessions(value?: zitadel_session_v2beta_session_pb.Session, index?: number): zitadel_session_v2beta_session_pb.Session;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSessionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSessionsResponse): ListSessionsResponse.AsObject;
  static serializeBinaryToWriter(message: ListSessionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSessionsResponse;
  static deserializeBinaryFromReader(message: ListSessionsResponse, reader: jspb.BinaryReader): ListSessionsResponse;
}

export namespace ListSessionsResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.ListDetails.AsObject,
    sessionsList: Array<zitadel_session_v2beta_session_pb.Session.AsObject>,
  }
}

export class GetSessionRequest extends jspb.Message {
  getSessionId(): string;
  setSessionId(value: string): GetSessionRequest;

  getSessionToken(): string;
  setSessionToken(value: string): GetSessionRequest;
  hasSessionToken(): boolean;
  clearSessionToken(): GetSessionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSessionRequest): GetSessionRequest.AsObject;
  static serializeBinaryToWriter(message: GetSessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSessionRequest;
  static deserializeBinaryFromReader(message: GetSessionRequest, reader: jspb.BinaryReader): GetSessionRequest;
}

export namespace GetSessionRequest {
  export type AsObject = {
    sessionId: string,
    sessionToken?: string,
  }

  export enum SessionTokenCase { 
    _SESSION_TOKEN_NOT_SET = 0,
    SESSION_TOKEN = 2,
  }
}

export class GetSessionResponse extends jspb.Message {
  getSession(): zitadel_session_v2beta_session_pb.Session | undefined;
  setSession(value?: zitadel_session_v2beta_session_pb.Session): GetSessionResponse;
  hasSession(): boolean;
  clearSession(): GetSessionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSessionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSessionResponse): GetSessionResponse.AsObject;
  static serializeBinaryToWriter(message: GetSessionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSessionResponse;
  static deserializeBinaryFromReader(message: GetSessionResponse, reader: jspb.BinaryReader): GetSessionResponse;
}

export namespace GetSessionResponse {
  export type AsObject = {
    session?: zitadel_session_v2beta_session_pb.Session.AsObject,
  }
}

export class CreateSessionRequest extends jspb.Message {
  getChecks(): Checks | undefined;
  setChecks(value?: Checks): CreateSessionRequest;
  hasChecks(): boolean;
  clearChecks(): CreateSessionRequest;

  getMetadataMap(): jspb.Map<string, Uint8Array | string>;
  clearMetadataMap(): CreateSessionRequest;

  getChallenges(): zitadel_session_v2beta_challenge_pb.RequestChallenges | undefined;
  setChallenges(value?: zitadel_session_v2beta_challenge_pb.RequestChallenges): CreateSessionRequest;
  hasChallenges(): boolean;
  clearChallenges(): CreateSessionRequest;

  getUserAgent(): zitadel_session_v2beta_session_pb.UserAgent | undefined;
  setUserAgent(value?: zitadel_session_v2beta_session_pb.UserAgent): CreateSessionRequest;
  hasUserAgent(): boolean;
  clearUserAgent(): CreateSessionRequest;

  getLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setLifetime(value?: google_protobuf_duration_pb.Duration): CreateSessionRequest;
  hasLifetime(): boolean;
  clearLifetime(): CreateSessionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSessionRequest): CreateSessionRequest.AsObject;
  static serializeBinaryToWriter(message: CreateSessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSessionRequest;
  static deserializeBinaryFromReader(message: CreateSessionRequest, reader: jspb.BinaryReader): CreateSessionRequest;
}

export namespace CreateSessionRequest {
  export type AsObject = {
    checks?: Checks.AsObject,
    metadataMap: Array<[string, Uint8Array | string]>,
    challenges?: zitadel_session_v2beta_challenge_pb.RequestChallenges.AsObject,
    userAgent?: zitadel_session_v2beta_session_pb.UserAgent.AsObject,
    lifetime?: google_protobuf_duration_pb.Duration.AsObject,
  }

  export enum LifetimeCase { 
    _LIFETIME_NOT_SET = 0,
    LIFETIME = 5,
  }
}

export class CreateSessionResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): CreateSessionResponse;
  hasDetails(): boolean;
  clearDetails(): CreateSessionResponse;

  getSessionId(): string;
  setSessionId(value: string): CreateSessionResponse;

  getSessionToken(): string;
  setSessionToken(value: string): CreateSessionResponse;

  getChallenges(): zitadel_session_v2beta_challenge_pb.Challenges | undefined;
  setChallenges(value?: zitadel_session_v2beta_challenge_pb.Challenges): CreateSessionResponse;
  hasChallenges(): boolean;
  clearChallenges(): CreateSessionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateSessionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateSessionResponse): CreateSessionResponse.AsObject;
  static serializeBinaryToWriter(message: CreateSessionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateSessionResponse;
  static deserializeBinaryFromReader(message: CreateSessionResponse, reader: jspb.BinaryReader): CreateSessionResponse;
}

export namespace CreateSessionResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    sessionId: string,
    sessionToken: string,
    challenges?: zitadel_session_v2beta_challenge_pb.Challenges.AsObject,
  }
}

export class SetSessionRequest extends jspb.Message {
  getSessionId(): string;
  setSessionId(value: string): SetSessionRequest;

  getSessionToken(): string;
  setSessionToken(value: string): SetSessionRequest;

  getChecks(): Checks | undefined;
  setChecks(value?: Checks): SetSessionRequest;
  hasChecks(): boolean;
  clearChecks(): SetSessionRequest;

  getMetadataMap(): jspb.Map<string, Uint8Array | string>;
  clearMetadataMap(): SetSessionRequest;

  getChallenges(): zitadel_session_v2beta_challenge_pb.RequestChallenges | undefined;
  setChallenges(value?: zitadel_session_v2beta_challenge_pb.RequestChallenges): SetSessionRequest;
  hasChallenges(): boolean;
  clearChallenges(): SetSessionRequest;

  getLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setLifetime(value?: google_protobuf_duration_pb.Duration): SetSessionRequest;
  hasLifetime(): boolean;
  clearLifetime(): SetSessionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetSessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetSessionRequest): SetSessionRequest.AsObject;
  static serializeBinaryToWriter(message: SetSessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetSessionRequest;
  static deserializeBinaryFromReader(message: SetSessionRequest, reader: jspb.BinaryReader): SetSessionRequest;
}

export namespace SetSessionRequest {
  export type AsObject = {
    sessionId: string,
    sessionToken: string,
    checks?: Checks.AsObject,
    metadataMap: Array<[string, Uint8Array | string]>,
    challenges?: zitadel_session_v2beta_challenge_pb.RequestChallenges.AsObject,
    lifetime?: google_protobuf_duration_pb.Duration.AsObject,
  }

  export enum LifetimeCase { 
    _LIFETIME_NOT_SET = 0,
    LIFETIME = 6,
  }
}

export class SetSessionResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetSessionResponse;
  hasDetails(): boolean;
  clearDetails(): SetSessionResponse;

  getSessionToken(): string;
  setSessionToken(value: string): SetSessionResponse;

  getChallenges(): zitadel_session_v2beta_challenge_pb.Challenges | undefined;
  setChallenges(value?: zitadel_session_v2beta_challenge_pb.Challenges): SetSessionResponse;
  hasChallenges(): boolean;
  clearChallenges(): SetSessionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetSessionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetSessionResponse): SetSessionResponse.AsObject;
  static serializeBinaryToWriter(message: SetSessionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetSessionResponse;
  static deserializeBinaryFromReader(message: SetSessionResponse, reader: jspb.BinaryReader): SetSessionResponse;
}

export namespace SetSessionResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    sessionToken: string,
    challenges?: zitadel_session_v2beta_challenge_pb.Challenges.AsObject,
  }
}

export class DeleteSessionRequest extends jspb.Message {
  getSessionId(): string;
  setSessionId(value: string): DeleteSessionRequest;

  getSessionToken(): string;
  setSessionToken(value: string): DeleteSessionRequest;
  hasSessionToken(): boolean;
  clearSessionToken(): DeleteSessionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSessionRequest): DeleteSessionRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteSessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSessionRequest;
  static deserializeBinaryFromReader(message: DeleteSessionRequest, reader: jspb.BinaryReader): DeleteSessionRequest;
}

export namespace DeleteSessionRequest {
  export type AsObject = {
    sessionId: string,
    sessionToken?: string,
  }

  export enum SessionTokenCase { 
    _SESSION_TOKEN_NOT_SET = 0,
    SESSION_TOKEN = 2,
  }
}

export class DeleteSessionResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): DeleteSessionResponse;
  hasDetails(): boolean;
  clearDetails(): DeleteSessionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteSessionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteSessionResponse): DeleteSessionResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteSessionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteSessionResponse;
  static deserializeBinaryFromReader(message: DeleteSessionResponse, reader: jspb.BinaryReader): DeleteSessionResponse;
}

export namespace DeleteSessionResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class Checks extends jspb.Message {
  getUser(): CheckUser | undefined;
  setUser(value?: CheckUser): Checks;
  hasUser(): boolean;
  clearUser(): Checks;

  getPassword(): CheckPassword | undefined;
  setPassword(value?: CheckPassword): Checks;
  hasPassword(): boolean;
  clearPassword(): Checks;

  getWebAuthN(): CheckWebAuthN | undefined;
  setWebAuthN(value?: CheckWebAuthN): Checks;
  hasWebAuthN(): boolean;
  clearWebAuthN(): Checks;

  getIdpIntent(): CheckIDPIntent | undefined;
  setIdpIntent(value?: CheckIDPIntent): Checks;
  hasIdpIntent(): boolean;
  clearIdpIntent(): Checks;

  getTotp(): CheckTOTP | undefined;
  setTotp(value?: CheckTOTP): Checks;
  hasTotp(): boolean;
  clearTotp(): Checks;

  getOtpSms(): CheckOTP | undefined;
  setOtpSms(value?: CheckOTP): Checks;
  hasOtpSms(): boolean;
  clearOtpSms(): Checks;

  getOtpEmail(): CheckOTP | undefined;
  setOtpEmail(value?: CheckOTP): Checks;
  hasOtpEmail(): boolean;
  clearOtpEmail(): Checks;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Checks.AsObject;
  static toObject(includeInstance: boolean, msg: Checks): Checks.AsObject;
  static serializeBinaryToWriter(message: Checks, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Checks;
  static deserializeBinaryFromReader(message: Checks, reader: jspb.BinaryReader): Checks;
}

export namespace Checks {
  export type AsObject = {
    user?: CheckUser.AsObject,
    password?: CheckPassword.AsObject,
    webAuthN?: CheckWebAuthN.AsObject,
    idpIntent?: CheckIDPIntent.AsObject,
    totp?: CheckTOTP.AsObject,
    otpSms?: CheckOTP.AsObject,
    otpEmail?: CheckOTP.AsObject,
  }

  export enum UserCase { 
    _USER_NOT_SET = 0,
    USER = 1,
  }

  export enum PasswordCase { 
    _PASSWORD_NOT_SET = 0,
    PASSWORD = 2,
  }

  export enum WebAuthNCase { 
    _WEB_AUTH_N_NOT_SET = 0,
    WEB_AUTH_N = 3,
  }

  export enum IdpIntentCase { 
    _IDP_INTENT_NOT_SET = 0,
    IDP_INTENT = 4,
  }

  export enum TotpCase { 
    _TOTP_NOT_SET = 0,
    TOTP = 5,
  }

  export enum OtpSmsCase { 
    _OTP_SMS_NOT_SET = 0,
    OTP_SMS = 6,
  }

  export enum OtpEmailCase { 
    _OTP_EMAIL_NOT_SET = 0,
    OTP_EMAIL = 7,
  }
}

export class CheckUser extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): CheckUser;

  getLoginName(): string;
  setLoginName(value: string): CheckUser;

  getSearchCase(): CheckUser.SearchCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckUser.AsObject;
  static toObject(includeInstance: boolean, msg: CheckUser): CheckUser.AsObject;
  static serializeBinaryToWriter(message: CheckUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckUser;
  static deserializeBinaryFromReader(message: CheckUser, reader: jspb.BinaryReader): CheckUser;
}

export namespace CheckUser {
  export type AsObject = {
    userId: string,
    loginName: string,
  }

  export enum SearchCase { 
    SEARCH_NOT_SET = 0,
    USER_ID = 1,
    LOGIN_NAME = 2,
  }
}

export class CheckPassword extends jspb.Message {
  getPassword(): string;
  setPassword(value: string): CheckPassword;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckPassword.AsObject;
  static toObject(includeInstance: boolean, msg: CheckPassword): CheckPassword.AsObject;
  static serializeBinaryToWriter(message: CheckPassword, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckPassword;
  static deserializeBinaryFromReader(message: CheckPassword, reader: jspb.BinaryReader): CheckPassword;
}

export namespace CheckPassword {
  export type AsObject = {
    password: string,
  }
}

export class CheckWebAuthN extends jspb.Message {
  getCredentialAssertionData(): google_protobuf_struct_pb.Struct | undefined;
  setCredentialAssertionData(value?: google_protobuf_struct_pb.Struct): CheckWebAuthN;
  hasCredentialAssertionData(): boolean;
  clearCredentialAssertionData(): CheckWebAuthN;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckWebAuthN.AsObject;
  static toObject(includeInstance: boolean, msg: CheckWebAuthN): CheckWebAuthN.AsObject;
  static serializeBinaryToWriter(message: CheckWebAuthN, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckWebAuthN;
  static deserializeBinaryFromReader(message: CheckWebAuthN, reader: jspb.BinaryReader): CheckWebAuthN;
}

export namespace CheckWebAuthN {
  export type AsObject = {
    credentialAssertionData?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class CheckIDPIntent extends jspb.Message {
  getIdpIntentId(): string;
  setIdpIntentId(value: string): CheckIDPIntent;

  getIdpIntentToken(): string;
  setIdpIntentToken(value: string): CheckIDPIntent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckIDPIntent.AsObject;
  static toObject(includeInstance: boolean, msg: CheckIDPIntent): CheckIDPIntent.AsObject;
  static serializeBinaryToWriter(message: CheckIDPIntent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckIDPIntent;
  static deserializeBinaryFromReader(message: CheckIDPIntent, reader: jspb.BinaryReader): CheckIDPIntent;
}

export namespace CheckIDPIntent {
  export type AsObject = {
    idpIntentId: string,
    idpIntentToken: string,
  }
}

export class CheckTOTP extends jspb.Message {
  getCode(): string;
  setCode(value: string): CheckTOTP;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckTOTP.AsObject;
  static toObject(includeInstance: boolean, msg: CheckTOTP): CheckTOTP.AsObject;
  static serializeBinaryToWriter(message: CheckTOTP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckTOTP;
  static deserializeBinaryFromReader(message: CheckTOTP, reader: jspb.BinaryReader): CheckTOTP;
}

export namespace CheckTOTP {
  export type AsObject = {
    code: string,
  }
}

export class CheckOTP extends jspb.Message {
  getCode(): string;
  setCode(value: string): CheckOTP;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckOTP.AsObject;
  static toObject(includeInstance: boolean, msg: CheckOTP): CheckOTP.AsObject;
  static serializeBinaryToWriter(message: CheckOTP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckOTP;
  static deserializeBinaryFromReader(message: CheckOTP, reader: jspb.BinaryReader): CheckOTP;
}

export namespace CheckOTP {
  export type AsObject = {
    code: string,
  }
}


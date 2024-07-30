import * as jspb from 'google-protobuf'

import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_protoc_gen_zitadel_v2_options_pb from '../../../zitadel/protoc_gen_zitadel/v2/options_pb'; // proto import: "zitadel/protoc_gen_zitadel/v2/options.proto"
import * as zitadel_user_v2beta_auth_pb from '../../../zitadel/user/v2beta/auth_pb'; // proto import: "zitadel/user/v2beta/auth.proto"
import * as zitadel_user_v2beta_email_pb from '../../../zitadel/user/v2beta/email_pb'; // proto import: "zitadel/user/v2beta/email.proto"
import * as zitadel_user_v2beta_phone_pb from '../../../zitadel/user/v2beta/phone_pb'; // proto import: "zitadel/user/v2beta/phone.proto"
import * as zitadel_user_v2beta_idp_pb from '../../../zitadel/user/v2beta/idp_pb'; // proto import: "zitadel/user/v2beta/idp.proto"
import * as zitadel_user_v2beta_password_pb from '../../../zitadel/user/v2beta/password_pb'; // proto import: "zitadel/user/v2beta/password.proto"
import * as zitadel_user_v2beta_user_pb from '../../../zitadel/user/v2beta/user_pb'; // proto import: "zitadel/user/v2beta/user.proto"
import * as zitadel_user_v2beta_query_pb from '../../../zitadel/user/v2beta/query_pb'; // proto import: "zitadel/user/v2beta/query.proto"
import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class AddHumanUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddHumanUserRequest;
  hasUserId(): boolean;
  clearUserId(): AddHumanUserRequest;

  getUsername(): string;
  setUsername(value: string): AddHumanUserRequest;
  hasUsername(): boolean;
  clearUsername(): AddHumanUserRequest;

  getOrganization(): zitadel_object_v2beta_object_pb.Organization | undefined;
  setOrganization(value?: zitadel_object_v2beta_object_pb.Organization): AddHumanUserRequest;
  hasOrganization(): boolean;
  clearOrganization(): AddHumanUserRequest;

  getProfile(): zitadel_user_v2beta_user_pb.SetHumanProfile | undefined;
  setProfile(value?: zitadel_user_v2beta_user_pb.SetHumanProfile): AddHumanUserRequest;
  hasProfile(): boolean;
  clearProfile(): AddHumanUserRequest;

  getEmail(): zitadel_user_v2beta_email_pb.SetHumanEmail | undefined;
  setEmail(value?: zitadel_user_v2beta_email_pb.SetHumanEmail): AddHumanUserRequest;
  hasEmail(): boolean;
  clearEmail(): AddHumanUserRequest;

  getPhone(): zitadel_user_v2beta_phone_pb.SetHumanPhone | undefined;
  setPhone(value?: zitadel_user_v2beta_phone_pb.SetHumanPhone): AddHumanUserRequest;
  hasPhone(): boolean;
  clearPhone(): AddHumanUserRequest;

  getMetadataList(): Array<zitadel_user_v2beta_user_pb.SetMetadataEntry>;
  setMetadataList(value: Array<zitadel_user_v2beta_user_pb.SetMetadataEntry>): AddHumanUserRequest;
  clearMetadataList(): AddHumanUserRequest;
  addMetadata(value?: zitadel_user_v2beta_user_pb.SetMetadataEntry, index?: number): zitadel_user_v2beta_user_pb.SetMetadataEntry;

  getPassword(): zitadel_user_v2beta_password_pb.Password | undefined;
  setPassword(value?: zitadel_user_v2beta_password_pb.Password): AddHumanUserRequest;
  hasPassword(): boolean;
  clearPassword(): AddHumanUserRequest;

  getHashedPassword(): zitadel_user_v2beta_password_pb.HashedPassword | undefined;
  setHashedPassword(value?: zitadel_user_v2beta_password_pb.HashedPassword): AddHumanUserRequest;
  hasHashedPassword(): boolean;
  clearHashedPassword(): AddHumanUserRequest;

  getIdpLinksList(): Array<zitadel_user_v2beta_idp_pb.IDPLink>;
  setIdpLinksList(value: Array<zitadel_user_v2beta_idp_pb.IDPLink>): AddHumanUserRequest;
  clearIdpLinksList(): AddHumanUserRequest;
  addIdpLinks(value?: zitadel_user_v2beta_idp_pb.IDPLink, index?: number): zitadel_user_v2beta_idp_pb.IDPLink;

  getTotpSecret(): string;
  setTotpSecret(value: string): AddHumanUserRequest;
  hasTotpSecret(): boolean;
  clearTotpSecret(): AddHumanUserRequest;

  getPasswordTypeCase(): AddHumanUserRequest.PasswordTypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddHumanUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddHumanUserRequest): AddHumanUserRequest.AsObject;
  static serializeBinaryToWriter(message: AddHumanUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddHumanUserRequest;
  static deserializeBinaryFromReader(message: AddHumanUserRequest, reader: jspb.BinaryReader): AddHumanUserRequest;
}

export namespace AddHumanUserRequest {
  export type AsObject = {
    userId?: string,
    username?: string,
    organization?: zitadel_object_v2beta_object_pb.Organization.AsObject,
    profile?: zitadel_user_v2beta_user_pb.SetHumanProfile.AsObject,
    email?: zitadel_user_v2beta_email_pb.SetHumanEmail.AsObject,
    phone?: zitadel_user_v2beta_phone_pb.SetHumanPhone.AsObject,
    metadataList: Array<zitadel_user_v2beta_user_pb.SetMetadataEntry.AsObject>,
    password?: zitadel_user_v2beta_password_pb.Password.AsObject,
    hashedPassword?: zitadel_user_v2beta_password_pb.HashedPassword.AsObject,
    idpLinksList: Array<zitadel_user_v2beta_idp_pb.IDPLink.AsObject>,
    totpSecret?: string,
  }

  export enum PasswordTypeCase { 
    PASSWORD_TYPE_NOT_SET = 0,
    PASSWORD = 7,
    HASHED_PASSWORD = 8,
  }

  export enum UserIdCase { 
    _USER_ID_NOT_SET = 0,
    USER_ID = 1,
  }

  export enum UsernameCase { 
    _USERNAME_NOT_SET = 0,
    USERNAME = 2,
  }

  export enum TotpSecretCase { 
    _TOTP_SECRET_NOT_SET = 0,
    TOTP_SECRET = 12,
  }
}

export class AddHumanUserResponse extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddHumanUserResponse;

  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddHumanUserResponse;
  hasDetails(): boolean;
  clearDetails(): AddHumanUserResponse;

  getEmailCode(): string;
  setEmailCode(value: string): AddHumanUserResponse;
  hasEmailCode(): boolean;
  clearEmailCode(): AddHumanUserResponse;

  getPhoneCode(): string;
  setPhoneCode(value: string): AddHumanUserResponse;
  hasPhoneCode(): boolean;
  clearPhoneCode(): AddHumanUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddHumanUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddHumanUserResponse): AddHumanUserResponse.AsObject;
  static serializeBinaryToWriter(message: AddHumanUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddHumanUserResponse;
  static deserializeBinaryFromReader(message: AddHumanUserResponse, reader: jspb.BinaryReader): AddHumanUserResponse;
}

export namespace AddHumanUserResponse {
  export type AsObject = {
    userId: string,
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    emailCode?: string,
    phoneCode?: string,
  }

  export enum EmailCodeCase { 
    _EMAIL_CODE_NOT_SET = 0,
    EMAIL_CODE = 3,
  }

  export enum PhoneCodeCase { 
    _PHONE_CODE_NOT_SET = 0,
    PHONE_CODE = 4,
  }
}

export class GetUserByIDRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GetUserByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserByIDRequest): GetUserByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserByIDRequest;
  static deserializeBinaryFromReader(message: GetUserByIDRequest, reader: jspb.BinaryReader): GetUserByIDRequest;
}

export namespace GetUserByIDRequest {
  export type AsObject = {
    userId: string,
  }
}

export class GetUserByIDResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): GetUserByIDResponse;
  hasDetails(): boolean;
  clearDetails(): GetUserByIDResponse;

  getUser(): zitadel_user_v2beta_user_pb.User | undefined;
  setUser(value?: zitadel_user_v2beta_user_pb.User): GetUserByIDResponse;
  hasUser(): boolean;
  clearUser(): GetUserByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserByIDResponse): GetUserByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserByIDResponse;
  static deserializeBinaryFromReader(message: GetUserByIDResponse, reader: jspb.BinaryReader): GetUserByIDResponse;
}

export namespace GetUserByIDResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    user?: zitadel_user_v2beta_user_pb.User.AsObject,
  }
}

export class ListUsersRequest extends jspb.Message {
  getQuery(): zitadel_object_v2beta_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_v2beta_object_pb.ListQuery): ListUsersRequest;
  hasQuery(): boolean;
  clearQuery(): ListUsersRequest;

  getSortingColumn(): zitadel_user_v2beta_query_pb.UserFieldName;
  setSortingColumn(value: zitadel_user_v2beta_query_pb.UserFieldName): ListUsersRequest;

  getQueriesList(): Array<zitadel_user_v2beta_query_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_user_v2beta_query_pb.SearchQuery>): ListUsersRequest;
  clearQueriesList(): ListUsersRequest;
  addQueries(value?: zitadel_user_v2beta_query_pb.SearchQuery, index?: number): zitadel_user_v2beta_query_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUsersRequest): ListUsersRequest.AsObject;
  static serializeBinaryToWriter(message: ListUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUsersRequest;
  static deserializeBinaryFromReader(message: ListUsersRequest, reader: jspb.BinaryReader): ListUsersRequest;
}

export namespace ListUsersRequest {
  export type AsObject = {
    query?: zitadel_object_v2beta_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_user_v2beta_query_pb.UserFieldName,
    queriesList: Array<zitadel_user_v2beta_query_pb.SearchQuery.AsObject>,
  }
}

export class ListUsersResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.ListDetails): ListUsersResponse;
  hasDetails(): boolean;
  clearDetails(): ListUsersResponse;

  getSortingColumn(): zitadel_user_v2beta_query_pb.UserFieldName;
  setSortingColumn(value: zitadel_user_v2beta_query_pb.UserFieldName): ListUsersResponse;

  getResultList(): Array<zitadel_user_v2beta_user_pb.User>;
  setResultList(value: Array<zitadel_user_v2beta_user_pb.User>): ListUsersResponse;
  clearResultList(): ListUsersResponse;
  addResult(value?: zitadel_user_v2beta_user_pb.User, index?: number): zitadel_user_v2beta_user_pb.User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUsersResponse): ListUsersResponse.AsObject;
  static serializeBinaryToWriter(message: ListUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUsersResponse;
  static deserializeBinaryFromReader(message: ListUsersResponse, reader: jspb.BinaryReader): ListUsersResponse;
}

export namespace ListUsersResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_user_v2beta_query_pb.UserFieldName,
    resultList: Array<zitadel_user_v2beta_user_pb.User.AsObject>,
  }
}

export class SetEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetEmailRequest;

  getEmail(): string;
  setEmail(value: string): SetEmailRequest;

  getSendCode(): zitadel_user_v2beta_email_pb.SendEmailVerificationCode | undefined;
  setSendCode(value?: zitadel_user_v2beta_email_pb.SendEmailVerificationCode): SetEmailRequest;
  hasSendCode(): boolean;
  clearSendCode(): SetEmailRequest;

  getReturnCode(): zitadel_user_v2beta_email_pb.ReturnEmailVerificationCode | undefined;
  setReturnCode(value?: zitadel_user_v2beta_email_pb.ReturnEmailVerificationCode): SetEmailRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): SetEmailRequest;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): SetEmailRequest;

  getVerificationCase(): SetEmailRequest.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetEmailRequest): SetEmailRequest.AsObject;
  static serializeBinaryToWriter(message: SetEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetEmailRequest;
  static deserializeBinaryFromReader(message: SetEmailRequest, reader: jspb.BinaryReader): SetEmailRequest;
}

export namespace SetEmailRequest {
  export type AsObject = {
    userId: string,
    email: string,
    sendCode?: zitadel_user_v2beta_email_pb.SendEmailVerificationCode.AsObject,
    returnCode?: zitadel_user_v2beta_email_pb.ReturnEmailVerificationCode.AsObject,
    isVerified: boolean,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    SEND_CODE = 3,
    RETURN_CODE = 4,
    IS_VERIFIED = 5,
  }
}

export class SetEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetEmailResponse;
  hasDetails(): boolean;
  clearDetails(): SetEmailResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): SetEmailResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): SetEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetEmailResponse): SetEmailResponse.AsObject;
  static serializeBinaryToWriter(message: SetEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetEmailResponse;
  static deserializeBinaryFromReader(message: SetEmailResponse, reader: jspb.BinaryReader): SetEmailResponse;
}

export namespace SetEmailResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 2,
  }
}

export class ResendEmailCodeRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ResendEmailCodeRequest;

  getSendCode(): zitadel_user_v2beta_email_pb.SendEmailVerificationCode | undefined;
  setSendCode(value?: zitadel_user_v2beta_email_pb.SendEmailVerificationCode): ResendEmailCodeRequest;
  hasSendCode(): boolean;
  clearSendCode(): ResendEmailCodeRequest;

  getReturnCode(): zitadel_user_v2beta_email_pb.ReturnEmailVerificationCode | undefined;
  setReturnCode(value?: zitadel_user_v2beta_email_pb.ReturnEmailVerificationCode): ResendEmailCodeRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): ResendEmailCodeRequest;

  getVerificationCase(): ResendEmailCodeRequest.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendEmailCodeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendEmailCodeRequest): ResendEmailCodeRequest.AsObject;
  static serializeBinaryToWriter(message: ResendEmailCodeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendEmailCodeRequest;
  static deserializeBinaryFromReader(message: ResendEmailCodeRequest, reader: jspb.BinaryReader): ResendEmailCodeRequest;
}

export namespace ResendEmailCodeRequest {
  export type AsObject = {
    userId: string,
    sendCode?: zitadel_user_v2beta_email_pb.SendEmailVerificationCode.AsObject,
    returnCode?: zitadel_user_v2beta_email_pb.ReturnEmailVerificationCode.AsObject,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    SEND_CODE = 2,
    RETURN_CODE = 3,
  }
}

export class ResendEmailCodeResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ResendEmailCodeResponse;
  hasDetails(): boolean;
  clearDetails(): ResendEmailCodeResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): ResendEmailCodeResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): ResendEmailCodeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendEmailCodeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendEmailCodeResponse): ResendEmailCodeResponse.AsObject;
  static serializeBinaryToWriter(message: ResendEmailCodeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendEmailCodeResponse;
  static deserializeBinaryFromReader(message: ResendEmailCodeResponse, reader: jspb.BinaryReader): ResendEmailCodeResponse;
}

export namespace ResendEmailCodeResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 2,
  }
}

export class VerifyEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyEmailRequest;

  getVerificationCode(): string;
  setVerificationCode(value: string): VerifyEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyEmailRequest): VerifyEmailRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyEmailRequest;
  static deserializeBinaryFromReader(message: VerifyEmailRequest, reader: jspb.BinaryReader): VerifyEmailRequest;
}

export namespace VerifyEmailRequest {
  export type AsObject = {
    userId: string,
    verificationCode: string,
  }
}

export class VerifyEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyEmailResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyEmailResponse): VerifyEmailResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyEmailResponse;
  static deserializeBinaryFromReader(message: VerifyEmailResponse, reader: jspb.BinaryReader): VerifyEmailResponse;
}

export namespace VerifyEmailResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class SetPhoneRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetPhoneRequest;

  getPhone(): string;
  setPhone(value: string): SetPhoneRequest;

  getSendCode(): zitadel_user_v2beta_phone_pb.SendPhoneVerificationCode | undefined;
  setSendCode(value?: zitadel_user_v2beta_phone_pb.SendPhoneVerificationCode): SetPhoneRequest;
  hasSendCode(): boolean;
  clearSendCode(): SetPhoneRequest;

  getReturnCode(): zitadel_user_v2beta_phone_pb.ReturnPhoneVerificationCode | undefined;
  setReturnCode(value?: zitadel_user_v2beta_phone_pb.ReturnPhoneVerificationCode): SetPhoneRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): SetPhoneRequest;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): SetPhoneRequest;

  getVerificationCase(): SetPhoneRequest.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetPhoneRequest): SetPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: SetPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPhoneRequest;
  static deserializeBinaryFromReader(message: SetPhoneRequest, reader: jspb.BinaryReader): SetPhoneRequest;
}

export namespace SetPhoneRequest {
  export type AsObject = {
    userId: string,
    phone: string,
    sendCode?: zitadel_user_v2beta_phone_pb.SendPhoneVerificationCode.AsObject,
    returnCode?: zitadel_user_v2beta_phone_pb.ReturnPhoneVerificationCode.AsObject,
    isVerified: boolean,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    SEND_CODE = 3,
    RETURN_CODE = 4,
    IS_VERIFIED = 5,
  }
}

export class SetPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): SetPhoneResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): SetPhoneResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): SetPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetPhoneResponse): SetPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: SetPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPhoneResponse;
  static deserializeBinaryFromReader(message: SetPhoneResponse, reader: jspb.BinaryReader): SetPhoneResponse;
}

export namespace SetPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 2,
  }
}

export class ResendPhoneCodeRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ResendPhoneCodeRequest;

  getSendCode(): zitadel_user_v2beta_phone_pb.SendPhoneVerificationCode | undefined;
  setSendCode(value?: zitadel_user_v2beta_phone_pb.SendPhoneVerificationCode): ResendPhoneCodeRequest;
  hasSendCode(): boolean;
  clearSendCode(): ResendPhoneCodeRequest;

  getReturnCode(): zitadel_user_v2beta_phone_pb.ReturnPhoneVerificationCode | undefined;
  setReturnCode(value?: zitadel_user_v2beta_phone_pb.ReturnPhoneVerificationCode): ResendPhoneCodeRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): ResendPhoneCodeRequest;

  getVerificationCase(): ResendPhoneCodeRequest.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendPhoneCodeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendPhoneCodeRequest): ResendPhoneCodeRequest.AsObject;
  static serializeBinaryToWriter(message: ResendPhoneCodeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendPhoneCodeRequest;
  static deserializeBinaryFromReader(message: ResendPhoneCodeRequest, reader: jspb.BinaryReader): ResendPhoneCodeRequest;
}

export namespace ResendPhoneCodeRequest {
  export type AsObject = {
    userId: string,
    sendCode?: zitadel_user_v2beta_phone_pb.SendPhoneVerificationCode.AsObject,
    returnCode?: zitadel_user_v2beta_phone_pb.ReturnPhoneVerificationCode.AsObject,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    SEND_CODE = 3,
    RETURN_CODE = 4,
  }
}

export class ResendPhoneCodeResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ResendPhoneCodeResponse;
  hasDetails(): boolean;
  clearDetails(): ResendPhoneCodeResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): ResendPhoneCodeResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): ResendPhoneCodeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendPhoneCodeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendPhoneCodeResponse): ResendPhoneCodeResponse.AsObject;
  static serializeBinaryToWriter(message: ResendPhoneCodeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendPhoneCodeResponse;
  static deserializeBinaryFromReader(message: ResendPhoneCodeResponse, reader: jspb.BinaryReader): ResendPhoneCodeResponse;
}

export namespace ResendPhoneCodeResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 2,
  }
}

export class VerifyPhoneRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyPhoneRequest;

  getVerificationCode(): string;
  setVerificationCode(value: string): VerifyPhoneRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyPhoneRequest): VerifyPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyPhoneRequest;
  static deserializeBinaryFromReader(message: VerifyPhoneRequest, reader: jspb.BinaryReader): VerifyPhoneRequest;
}

export namespace VerifyPhoneRequest {
  export type AsObject = {
    userId: string,
    verificationCode: string,
  }
}

export class VerifyPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyPhoneResponse): VerifyPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyPhoneResponse;
  static deserializeBinaryFromReader(message: VerifyPhoneResponse, reader: jspb.BinaryReader): VerifyPhoneResponse;
}

export namespace VerifyPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class DeleteUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): DeleteUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUserRequest): DeleteUserRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUserRequest;
  static deserializeBinaryFromReader(message: DeleteUserRequest, reader: jspb.BinaryReader): DeleteUserRequest;
}

export namespace DeleteUserRequest {
  export type AsObject = {
    userId: string,
  }
}

export class DeleteUserResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): DeleteUserResponse;
  hasDetails(): boolean;
  clearDetails(): DeleteUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUserResponse): DeleteUserResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUserResponse;
  static deserializeBinaryFromReader(message: DeleteUserResponse, reader: jspb.BinaryReader): DeleteUserResponse;
}

export namespace DeleteUserResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class UpdateHumanUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateHumanUserRequest;

  getUsername(): string;
  setUsername(value: string): UpdateHumanUserRequest;
  hasUsername(): boolean;
  clearUsername(): UpdateHumanUserRequest;

  getProfile(): zitadel_user_v2beta_user_pb.SetHumanProfile | undefined;
  setProfile(value?: zitadel_user_v2beta_user_pb.SetHumanProfile): UpdateHumanUserRequest;
  hasProfile(): boolean;
  clearProfile(): UpdateHumanUserRequest;

  getEmail(): zitadel_user_v2beta_email_pb.SetHumanEmail | undefined;
  setEmail(value?: zitadel_user_v2beta_email_pb.SetHumanEmail): UpdateHumanUserRequest;
  hasEmail(): boolean;
  clearEmail(): UpdateHumanUserRequest;

  getPhone(): zitadel_user_v2beta_phone_pb.SetHumanPhone | undefined;
  setPhone(value?: zitadel_user_v2beta_phone_pb.SetHumanPhone): UpdateHumanUserRequest;
  hasPhone(): boolean;
  clearPhone(): UpdateHumanUserRequest;

  getPassword(): zitadel_user_v2beta_password_pb.SetPassword | undefined;
  setPassword(value?: zitadel_user_v2beta_password_pb.SetPassword): UpdateHumanUserRequest;
  hasPassword(): boolean;
  clearPassword(): UpdateHumanUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateHumanUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateHumanUserRequest): UpdateHumanUserRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateHumanUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateHumanUserRequest;
  static deserializeBinaryFromReader(message: UpdateHumanUserRequest, reader: jspb.BinaryReader): UpdateHumanUserRequest;
}

export namespace UpdateHumanUserRequest {
  export type AsObject = {
    userId: string,
    username?: string,
    profile?: zitadel_user_v2beta_user_pb.SetHumanProfile.AsObject,
    email?: zitadel_user_v2beta_email_pb.SetHumanEmail.AsObject,
    phone?: zitadel_user_v2beta_phone_pb.SetHumanPhone.AsObject,
    password?: zitadel_user_v2beta_password_pb.SetPassword.AsObject,
  }

  export enum UsernameCase { 
    _USERNAME_NOT_SET = 0,
    USERNAME = 2,
  }

  export enum ProfileCase { 
    _PROFILE_NOT_SET = 0,
    PROFILE = 3,
  }

  export enum EmailCase { 
    _EMAIL_NOT_SET = 0,
    EMAIL = 4,
  }

  export enum PhoneCase { 
    _PHONE_NOT_SET = 0,
    PHONE = 5,
  }

  export enum PasswordCase { 
    _PASSWORD_NOT_SET = 0,
    PASSWORD = 6,
  }
}

export class UpdateHumanUserResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): UpdateHumanUserResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateHumanUserResponse;

  getEmailCode(): string;
  setEmailCode(value: string): UpdateHumanUserResponse;
  hasEmailCode(): boolean;
  clearEmailCode(): UpdateHumanUserResponse;

  getPhoneCode(): string;
  setPhoneCode(value: string): UpdateHumanUserResponse;
  hasPhoneCode(): boolean;
  clearPhoneCode(): UpdateHumanUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateHumanUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateHumanUserResponse): UpdateHumanUserResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateHumanUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateHumanUserResponse;
  static deserializeBinaryFromReader(message: UpdateHumanUserResponse, reader: jspb.BinaryReader): UpdateHumanUserResponse;
}

export namespace UpdateHumanUserResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    emailCode?: string,
    phoneCode?: string,
  }

  export enum EmailCodeCase { 
    _EMAIL_CODE_NOT_SET = 0,
    EMAIL_CODE = 2,
  }

  export enum PhoneCodeCase { 
    _PHONE_CODE_NOT_SET = 0,
    PHONE_CODE = 3,
  }
}

export class DeactivateUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): DeactivateUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateUserRequest): DeactivateUserRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateUserRequest;
  static deserializeBinaryFromReader(message: DeactivateUserRequest, reader: jspb.BinaryReader): DeactivateUserRequest;
}

export namespace DeactivateUserRequest {
  export type AsObject = {
    userId: string,
  }
}

export class DeactivateUserResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): DeactivateUserResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateUserResponse): DeactivateUserResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateUserResponse;
  static deserializeBinaryFromReader(message: DeactivateUserResponse, reader: jspb.BinaryReader): DeactivateUserResponse;
}

export namespace DeactivateUserResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ReactivateUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ReactivateUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateUserRequest): ReactivateUserRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateUserRequest;
  static deserializeBinaryFromReader(message: ReactivateUserRequest, reader: jspb.BinaryReader): ReactivateUserRequest;
}

export namespace ReactivateUserRequest {
  export type AsObject = {
    userId: string,
  }
}

export class ReactivateUserResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ReactivateUserResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateUserResponse): ReactivateUserResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateUserResponse;
  static deserializeBinaryFromReader(message: ReactivateUserResponse, reader: jspb.BinaryReader): ReactivateUserResponse;
}

export namespace ReactivateUserResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class LockUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): LockUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LockUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LockUserRequest): LockUserRequest.AsObject;
  static serializeBinaryToWriter(message: LockUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LockUserRequest;
  static deserializeBinaryFromReader(message: LockUserRequest, reader: jspb.BinaryReader): LockUserRequest;
}

export namespace LockUserRequest {
  export type AsObject = {
    userId: string,
  }
}

export class LockUserResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): LockUserResponse;
  hasDetails(): boolean;
  clearDetails(): LockUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LockUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LockUserResponse): LockUserResponse.AsObject;
  static serializeBinaryToWriter(message: LockUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LockUserResponse;
  static deserializeBinaryFromReader(message: LockUserResponse, reader: jspb.BinaryReader): LockUserResponse;
}

export namespace LockUserResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class UnlockUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UnlockUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnlockUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UnlockUserRequest): UnlockUserRequest.AsObject;
  static serializeBinaryToWriter(message: UnlockUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnlockUserRequest;
  static deserializeBinaryFromReader(message: UnlockUserRequest, reader: jspb.BinaryReader): UnlockUserRequest;
}

export namespace UnlockUserRequest {
  export type AsObject = {
    userId: string,
  }
}

export class UnlockUserResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): UnlockUserResponse;
  hasDetails(): boolean;
  clearDetails(): UnlockUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnlockUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UnlockUserResponse): UnlockUserResponse.AsObject;
  static serializeBinaryToWriter(message: UnlockUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnlockUserResponse;
  static deserializeBinaryFromReader(message: UnlockUserResponse, reader: jspb.BinaryReader): UnlockUserResponse;
}

export namespace UnlockUserResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RegisterPasskeyRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RegisterPasskeyRequest;

  getCode(): zitadel_user_v2beta_auth_pb.PasskeyRegistrationCode | undefined;
  setCode(value?: zitadel_user_v2beta_auth_pb.PasskeyRegistrationCode): RegisterPasskeyRequest;
  hasCode(): boolean;
  clearCode(): RegisterPasskeyRequest;

  getAuthenticator(): zitadel_user_v2beta_auth_pb.PasskeyAuthenticator;
  setAuthenticator(value: zitadel_user_v2beta_auth_pb.PasskeyAuthenticator): RegisterPasskeyRequest;

  getDomain(): string;
  setDomain(value: string): RegisterPasskeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterPasskeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterPasskeyRequest): RegisterPasskeyRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterPasskeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterPasskeyRequest;
  static deserializeBinaryFromReader(message: RegisterPasskeyRequest, reader: jspb.BinaryReader): RegisterPasskeyRequest;
}

export namespace RegisterPasskeyRequest {
  export type AsObject = {
    userId: string,
    code?: zitadel_user_v2beta_auth_pb.PasskeyRegistrationCode.AsObject,
    authenticator: zitadel_user_v2beta_auth_pb.PasskeyAuthenticator,
    domain: string,
  }

  export enum CodeCase { 
    _CODE_NOT_SET = 0,
    CODE = 2,
  }
}

export class RegisterPasskeyResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RegisterPasskeyResponse;
  hasDetails(): boolean;
  clearDetails(): RegisterPasskeyResponse;

  getPasskeyId(): string;
  setPasskeyId(value: string): RegisterPasskeyResponse;

  getPublicKeyCredentialCreationOptions(): google_protobuf_struct_pb.Struct | undefined;
  setPublicKeyCredentialCreationOptions(value?: google_protobuf_struct_pb.Struct): RegisterPasskeyResponse;
  hasPublicKeyCredentialCreationOptions(): boolean;
  clearPublicKeyCredentialCreationOptions(): RegisterPasskeyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterPasskeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterPasskeyResponse): RegisterPasskeyResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterPasskeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterPasskeyResponse;
  static deserializeBinaryFromReader(message: RegisterPasskeyResponse, reader: jspb.BinaryReader): RegisterPasskeyResponse;
}

export namespace RegisterPasskeyResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    passkeyId: string,
    publicKeyCredentialCreationOptions?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class VerifyPasskeyRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyPasskeyRegistrationRequest;

  getPasskeyId(): string;
  setPasskeyId(value: string): VerifyPasskeyRegistrationRequest;

  getPublicKeyCredential(): google_protobuf_struct_pb.Struct | undefined;
  setPublicKeyCredential(value?: google_protobuf_struct_pb.Struct): VerifyPasskeyRegistrationRequest;
  hasPublicKeyCredential(): boolean;
  clearPublicKeyCredential(): VerifyPasskeyRegistrationRequest;

  getPasskeyName(): string;
  setPasskeyName(value: string): VerifyPasskeyRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyPasskeyRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyPasskeyRegistrationRequest): VerifyPasskeyRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyPasskeyRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyPasskeyRegistrationRequest;
  static deserializeBinaryFromReader(message: VerifyPasskeyRegistrationRequest, reader: jspb.BinaryReader): VerifyPasskeyRegistrationRequest;
}

export namespace VerifyPasskeyRegistrationRequest {
  export type AsObject = {
    userId: string,
    passkeyId: string,
    publicKeyCredential?: google_protobuf_struct_pb.Struct.AsObject,
    passkeyName: string,
  }
}

export class VerifyPasskeyRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyPasskeyRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyPasskeyRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyPasskeyRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyPasskeyRegistrationResponse): VerifyPasskeyRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyPasskeyRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyPasskeyRegistrationResponse;
  static deserializeBinaryFromReader(message: VerifyPasskeyRegistrationResponse, reader: jspb.BinaryReader): VerifyPasskeyRegistrationResponse;
}

export namespace VerifyPasskeyRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RegisterU2FRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RegisterU2FRequest;

  getDomain(): string;
  setDomain(value: string): RegisterU2FRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterU2FRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterU2FRequest): RegisterU2FRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterU2FRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterU2FRequest;
  static deserializeBinaryFromReader(message: RegisterU2FRequest, reader: jspb.BinaryReader): RegisterU2FRequest;
}

export namespace RegisterU2FRequest {
  export type AsObject = {
    userId: string,
    domain: string,
  }
}

export class RegisterU2FResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RegisterU2FResponse;
  hasDetails(): boolean;
  clearDetails(): RegisterU2FResponse;

  getU2fId(): string;
  setU2fId(value: string): RegisterU2FResponse;

  getPublicKeyCredentialCreationOptions(): google_protobuf_struct_pb.Struct | undefined;
  setPublicKeyCredentialCreationOptions(value?: google_protobuf_struct_pb.Struct): RegisterU2FResponse;
  hasPublicKeyCredentialCreationOptions(): boolean;
  clearPublicKeyCredentialCreationOptions(): RegisterU2FResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterU2FResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterU2FResponse): RegisterU2FResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterU2FResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterU2FResponse;
  static deserializeBinaryFromReader(message: RegisterU2FResponse, reader: jspb.BinaryReader): RegisterU2FResponse;
}

export namespace RegisterU2FResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    u2fId: string,
    publicKeyCredentialCreationOptions?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class VerifyU2FRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyU2FRegistrationRequest;

  getU2fId(): string;
  setU2fId(value: string): VerifyU2FRegistrationRequest;

  getPublicKeyCredential(): google_protobuf_struct_pb.Struct | undefined;
  setPublicKeyCredential(value?: google_protobuf_struct_pb.Struct): VerifyU2FRegistrationRequest;
  hasPublicKeyCredential(): boolean;
  clearPublicKeyCredential(): VerifyU2FRegistrationRequest;

  getTokenName(): string;
  setTokenName(value: string): VerifyU2FRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyU2FRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyU2FRegistrationRequest): VerifyU2FRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyU2FRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyU2FRegistrationRequest;
  static deserializeBinaryFromReader(message: VerifyU2FRegistrationRequest, reader: jspb.BinaryReader): VerifyU2FRegistrationRequest;
}

export namespace VerifyU2FRegistrationRequest {
  export type AsObject = {
    userId: string,
    u2fId: string,
    publicKeyCredential?: google_protobuf_struct_pb.Struct.AsObject,
    tokenName: string,
  }
}

export class VerifyU2FRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyU2FRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyU2FRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyU2FRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyU2FRegistrationResponse): VerifyU2FRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyU2FRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyU2FRegistrationResponse;
  static deserializeBinaryFromReader(message: VerifyU2FRegistrationResponse, reader: jspb.BinaryReader): VerifyU2FRegistrationResponse;
}

export namespace VerifyU2FRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RegisterTOTPRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RegisterTOTPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterTOTPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterTOTPRequest): RegisterTOTPRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterTOTPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterTOTPRequest;
  static deserializeBinaryFromReader(message: RegisterTOTPRequest, reader: jspb.BinaryReader): RegisterTOTPRequest;
}

export namespace RegisterTOTPRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RegisterTOTPResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RegisterTOTPResponse;
  hasDetails(): boolean;
  clearDetails(): RegisterTOTPResponse;

  getUri(): string;
  setUri(value: string): RegisterTOTPResponse;

  getSecret(): string;
  setSecret(value: string): RegisterTOTPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterTOTPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterTOTPResponse): RegisterTOTPResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterTOTPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterTOTPResponse;
  static deserializeBinaryFromReader(message: RegisterTOTPResponse, reader: jspb.BinaryReader): RegisterTOTPResponse;
}

export namespace RegisterTOTPResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    uri: string,
    secret: string,
  }
}

export class VerifyTOTPRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyTOTPRegistrationRequest;

  getCode(): string;
  setCode(value: string): VerifyTOTPRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyTOTPRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyTOTPRegistrationRequest): VerifyTOTPRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyTOTPRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyTOTPRegistrationRequest;
  static deserializeBinaryFromReader(message: VerifyTOTPRegistrationRequest, reader: jspb.BinaryReader): VerifyTOTPRegistrationRequest;
}

export namespace VerifyTOTPRegistrationRequest {
  export type AsObject = {
    userId: string,
    code: string,
  }
}

export class VerifyTOTPRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyTOTPRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyTOTPRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyTOTPRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyTOTPRegistrationResponse): VerifyTOTPRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyTOTPRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyTOTPRegistrationResponse;
  static deserializeBinaryFromReader(message: VerifyTOTPRegistrationResponse, reader: jspb.BinaryReader): VerifyTOTPRegistrationResponse;
}

export namespace VerifyTOTPRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RemoveTOTPRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveTOTPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveTOTPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveTOTPRequest): RemoveTOTPRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveTOTPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveTOTPRequest;
  static deserializeBinaryFromReader(message: RemoveTOTPRequest, reader: jspb.BinaryReader): RemoveTOTPRequest;
}

export namespace RemoveTOTPRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveTOTPResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveTOTPResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveTOTPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveTOTPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveTOTPResponse): RemoveTOTPResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveTOTPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveTOTPResponse;
  static deserializeBinaryFromReader(message: RemoveTOTPResponse, reader: jspb.BinaryReader): RemoveTOTPResponse;
}

export namespace RemoveTOTPResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class AddOTPSMSRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddOTPSMSRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOTPSMSRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOTPSMSRequest): AddOTPSMSRequest.AsObject;
  static serializeBinaryToWriter(message: AddOTPSMSRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOTPSMSRequest;
  static deserializeBinaryFromReader(message: AddOTPSMSRequest, reader: jspb.BinaryReader): AddOTPSMSRequest;
}

export namespace AddOTPSMSRequest {
  export type AsObject = {
    userId: string,
  }
}

export class AddOTPSMSResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddOTPSMSResponse;
  hasDetails(): boolean;
  clearDetails(): AddOTPSMSResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOTPSMSResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOTPSMSResponse): AddOTPSMSResponse.AsObject;
  static serializeBinaryToWriter(message: AddOTPSMSResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOTPSMSResponse;
  static deserializeBinaryFromReader(message: AddOTPSMSResponse, reader: jspb.BinaryReader): AddOTPSMSResponse;
}

export namespace AddOTPSMSResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RemoveOTPSMSRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveOTPSMSRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOTPSMSRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOTPSMSRequest): RemoveOTPSMSRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOTPSMSRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOTPSMSRequest;
  static deserializeBinaryFromReader(message: RemoveOTPSMSRequest, reader: jspb.BinaryReader): RemoveOTPSMSRequest;
}

export namespace RemoveOTPSMSRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveOTPSMSResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveOTPSMSResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOTPSMSResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOTPSMSResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOTPSMSResponse): RemoveOTPSMSResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOTPSMSResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOTPSMSResponse;
  static deserializeBinaryFromReader(message: RemoveOTPSMSResponse, reader: jspb.BinaryReader): RemoveOTPSMSResponse;
}

export namespace RemoveOTPSMSResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class AddOTPEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddOTPEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOTPEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOTPEmailRequest): AddOTPEmailRequest.AsObject;
  static serializeBinaryToWriter(message: AddOTPEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOTPEmailRequest;
  static deserializeBinaryFromReader(message: AddOTPEmailRequest, reader: jspb.BinaryReader): AddOTPEmailRequest;
}

export namespace AddOTPEmailRequest {
  export type AsObject = {
    userId: string,
  }
}

export class AddOTPEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddOTPEmailResponse;
  hasDetails(): boolean;
  clearDetails(): AddOTPEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOTPEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOTPEmailResponse): AddOTPEmailResponse.AsObject;
  static serializeBinaryToWriter(message: AddOTPEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOTPEmailResponse;
  static deserializeBinaryFromReader(message: AddOTPEmailResponse, reader: jspb.BinaryReader): AddOTPEmailResponse;
}

export namespace AddOTPEmailResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RemoveOTPEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveOTPEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOTPEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOTPEmailRequest): RemoveOTPEmailRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOTPEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOTPEmailRequest;
  static deserializeBinaryFromReader(message: RemoveOTPEmailRequest, reader: jspb.BinaryReader): RemoveOTPEmailRequest;
}

export namespace RemoveOTPEmailRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveOTPEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveOTPEmailResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOTPEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOTPEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOTPEmailResponse): RemoveOTPEmailResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOTPEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOTPEmailResponse;
  static deserializeBinaryFromReader(message: RemoveOTPEmailResponse, reader: jspb.BinaryReader): RemoveOTPEmailResponse;
}

export namespace RemoveOTPEmailResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class CreatePasskeyRegistrationLinkRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): CreatePasskeyRegistrationLinkRequest;

  getSendLink(): zitadel_user_v2beta_auth_pb.SendPasskeyRegistrationLink | undefined;
  setSendLink(value?: zitadel_user_v2beta_auth_pb.SendPasskeyRegistrationLink): CreatePasskeyRegistrationLinkRequest;
  hasSendLink(): boolean;
  clearSendLink(): CreatePasskeyRegistrationLinkRequest;

  getReturnCode(): zitadel_user_v2beta_auth_pb.ReturnPasskeyRegistrationCode | undefined;
  setReturnCode(value?: zitadel_user_v2beta_auth_pb.ReturnPasskeyRegistrationCode): CreatePasskeyRegistrationLinkRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): CreatePasskeyRegistrationLinkRequest;

  getMediumCase(): CreatePasskeyRegistrationLinkRequest.MediumCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePasskeyRegistrationLinkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePasskeyRegistrationLinkRequest): CreatePasskeyRegistrationLinkRequest.AsObject;
  static serializeBinaryToWriter(message: CreatePasskeyRegistrationLinkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePasskeyRegistrationLinkRequest;
  static deserializeBinaryFromReader(message: CreatePasskeyRegistrationLinkRequest, reader: jspb.BinaryReader): CreatePasskeyRegistrationLinkRequest;
}

export namespace CreatePasskeyRegistrationLinkRequest {
  export type AsObject = {
    userId: string,
    sendLink?: zitadel_user_v2beta_auth_pb.SendPasskeyRegistrationLink.AsObject,
    returnCode?: zitadel_user_v2beta_auth_pb.ReturnPasskeyRegistrationCode.AsObject,
  }

  export enum MediumCase { 
    MEDIUM_NOT_SET = 0,
    SEND_LINK = 2,
    RETURN_CODE = 3,
  }
}

export class CreatePasskeyRegistrationLinkResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): CreatePasskeyRegistrationLinkResponse;
  hasDetails(): boolean;
  clearDetails(): CreatePasskeyRegistrationLinkResponse;

  getCode(): zitadel_user_v2beta_auth_pb.PasskeyRegistrationCode | undefined;
  setCode(value?: zitadel_user_v2beta_auth_pb.PasskeyRegistrationCode): CreatePasskeyRegistrationLinkResponse;
  hasCode(): boolean;
  clearCode(): CreatePasskeyRegistrationLinkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreatePasskeyRegistrationLinkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreatePasskeyRegistrationLinkResponse): CreatePasskeyRegistrationLinkResponse.AsObject;
  static serializeBinaryToWriter(message: CreatePasskeyRegistrationLinkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreatePasskeyRegistrationLinkResponse;
  static deserializeBinaryFromReader(message: CreatePasskeyRegistrationLinkResponse, reader: jspb.BinaryReader): CreatePasskeyRegistrationLinkResponse;
}

export namespace CreatePasskeyRegistrationLinkResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    code?: zitadel_user_v2beta_auth_pb.PasskeyRegistrationCode.AsObject,
  }

  export enum CodeCase { 
    _CODE_NOT_SET = 0,
    CODE = 2,
  }
}

export class StartIdentityProviderIntentRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): StartIdentityProviderIntentRequest;

  getUrls(): zitadel_user_v2beta_idp_pb.RedirectURLs | undefined;
  setUrls(value?: zitadel_user_v2beta_idp_pb.RedirectURLs): StartIdentityProviderIntentRequest;
  hasUrls(): boolean;
  clearUrls(): StartIdentityProviderIntentRequest;

  getLdap(): zitadel_user_v2beta_idp_pb.LDAPCredentials | undefined;
  setLdap(value?: zitadel_user_v2beta_idp_pb.LDAPCredentials): StartIdentityProviderIntentRequest;
  hasLdap(): boolean;
  clearLdap(): StartIdentityProviderIntentRequest;

  getContentCase(): StartIdentityProviderIntentRequest.ContentCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartIdentityProviderIntentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartIdentityProviderIntentRequest): StartIdentityProviderIntentRequest.AsObject;
  static serializeBinaryToWriter(message: StartIdentityProviderIntentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartIdentityProviderIntentRequest;
  static deserializeBinaryFromReader(message: StartIdentityProviderIntentRequest, reader: jspb.BinaryReader): StartIdentityProviderIntentRequest;
}

export namespace StartIdentityProviderIntentRequest {
  export type AsObject = {
    idpId: string,
    urls?: zitadel_user_v2beta_idp_pb.RedirectURLs.AsObject,
    ldap?: zitadel_user_v2beta_idp_pb.LDAPCredentials.AsObject,
  }

  export enum ContentCase { 
    CONTENT_NOT_SET = 0,
    URLS = 2,
    LDAP = 3,
  }
}

export class StartIdentityProviderIntentResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): StartIdentityProviderIntentResponse;
  hasDetails(): boolean;
  clearDetails(): StartIdentityProviderIntentResponse;

  getAuthUrl(): string;
  setAuthUrl(value: string): StartIdentityProviderIntentResponse;

  getIdpIntent(): zitadel_user_v2beta_idp_pb.IDPIntent | undefined;
  setIdpIntent(value?: zitadel_user_v2beta_idp_pb.IDPIntent): StartIdentityProviderIntentResponse;
  hasIdpIntent(): boolean;
  clearIdpIntent(): StartIdentityProviderIntentResponse;

  getPostForm(): Uint8Array | string;
  getPostForm_asU8(): Uint8Array;
  getPostForm_asB64(): string;
  setPostForm(value: Uint8Array | string): StartIdentityProviderIntentResponse;

  getNextStepCase(): StartIdentityProviderIntentResponse.NextStepCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartIdentityProviderIntentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StartIdentityProviderIntentResponse): StartIdentityProviderIntentResponse.AsObject;
  static serializeBinaryToWriter(message: StartIdentityProviderIntentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartIdentityProviderIntentResponse;
  static deserializeBinaryFromReader(message: StartIdentityProviderIntentResponse, reader: jspb.BinaryReader): StartIdentityProviderIntentResponse;
}

export namespace StartIdentityProviderIntentResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    authUrl: string,
    idpIntent?: zitadel_user_v2beta_idp_pb.IDPIntent.AsObject,
    postForm: Uint8Array | string,
  }

  export enum NextStepCase { 
    NEXT_STEP_NOT_SET = 0,
    AUTH_URL = 2,
    IDP_INTENT = 3,
    POST_FORM = 4,
  }
}

export class RetrieveIdentityProviderIntentRequest extends jspb.Message {
  getIdpIntentId(): string;
  setIdpIntentId(value: string): RetrieveIdentityProviderIntentRequest;

  getIdpIntentToken(): string;
  setIdpIntentToken(value: string): RetrieveIdentityProviderIntentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetrieveIdentityProviderIntentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RetrieveIdentityProviderIntentRequest): RetrieveIdentityProviderIntentRequest.AsObject;
  static serializeBinaryToWriter(message: RetrieveIdentityProviderIntentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetrieveIdentityProviderIntentRequest;
  static deserializeBinaryFromReader(message: RetrieveIdentityProviderIntentRequest, reader: jspb.BinaryReader): RetrieveIdentityProviderIntentRequest;
}

export namespace RetrieveIdentityProviderIntentRequest {
  export type AsObject = {
    idpIntentId: string,
    idpIntentToken: string,
  }
}

export class RetrieveIdentityProviderIntentResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RetrieveIdentityProviderIntentResponse;
  hasDetails(): boolean;
  clearDetails(): RetrieveIdentityProviderIntentResponse;

  getIdpInformation(): zitadel_user_v2beta_idp_pb.IDPInformation | undefined;
  setIdpInformation(value?: zitadel_user_v2beta_idp_pb.IDPInformation): RetrieveIdentityProviderIntentResponse;
  hasIdpInformation(): boolean;
  clearIdpInformation(): RetrieveIdentityProviderIntentResponse;

  getUserId(): string;
  setUserId(value: string): RetrieveIdentityProviderIntentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RetrieveIdentityProviderIntentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RetrieveIdentityProviderIntentResponse): RetrieveIdentityProviderIntentResponse.AsObject;
  static serializeBinaryToWriter(message: RetrieveIdentityProviderIntentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RetrieveIdentityProviderIntentResponse;
  static deserializeBinaryFromReader(message: RetrieveIdentityProviderIntentResponse, reader: jspb.BinaryReader): RetrieveIdentityProviderIntentResponse;
}

export namespace RetrieveIdentityProviderIntentResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    idpInformation?: zitadel_user_v2beta_idp_pb.IDPInformation.AsObject,
    userId: string,
  }
}

export class AddIDPLinkRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddIDPLinkRequest;

  getIdpLink(): zitadel_user_v2beta_idp_pb.IDPLink | undefined;
  setIdpLink(value?: zitadel_user_v2beta_idp_pb.IDPLink): AddIDPLinkRequest;
  hasIdpLink(): boolean;
  clearIdpLink(): AddIDPLinkRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIDPLinkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddIDPLinkRequest): AddIDPLinkRequest.AsObject;
  static serializeBinaryToWriter(message: AddIDPLinkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIDPLinkRequest;
  static deserializeBinaryFromReader(message: AddIDPLinkRequest, reader: jspb.BinaryReader): AddIDPLinkRequest;
}

export namespace AddIDPLinkRequest {
  export type AsObject = {
    userId: string,
    idpLink?: zitadel_user_v2beta_idp_pb.IDPLink.AsObject,
  }
}

export class AddIDPLinkResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddIDPLinkResponse;
  hasDetails(): boolean;
  clearDetails(): AddIDPLinkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIDPLinkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddIDPLinkResponse): AddIDPLinkResponse.AsObject;
  static serializeBinaryToWriter(message: AddIDPLinkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIDPLinkResponse;
  static deserializeBinaryFromReader(message: AddIDPLinkResponse, reader: jspb.BinaryReader): AddIDPLinkResponse;
}

export namespace AddIDPLinkResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class PasswordResetRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): PasswordResetRequest;

  getSendLink(): zitadel_user_v2beta_password_pb.SendPasswordResetLink | undefined;
  setSendLink(value?: zitadel_user_v2beta_password_pb.SendPasswordResetLink): PasswordResetRequest;
  hasSendLink(): boolean;
  clearSendLink(): PasswordResetRequest;

  getReturnCode(): zitadel_user_v2beta_password_pb.ReturnPasswordResetCode | undefined;
  setReturnCode(value?: zitadel_user_v2beta_password_pb.ReturnPasswordResetCode): PasswordResetRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): PasswordResetRequest;

  getMediumCase(): PasswordResetRequest.MediumCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordResetRequest.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordResetRequest): PasswordResetRequest.AsObject;
  static serializeBinaryToWriter(message: PasswordResetRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordResetRequest;
  static deserializeBinaryFromReader(message: PasswordResetRequest, reader: jspb.BinaryReader): PasswordResetRequest;
}

export namespace PasswordResetRequest {
  export type AsObject = {
    userId: string,
    sendLink?: zitadel_user_v2beta_password_pb.SendPasswordResetLink.AsObject,
    returnCode?: zitadel_user_v2beta_password_pb.ReturnPasswordResetCode.AsObject,
  }

  export enum MediumCase { 
    MEDIUM_NOT_SET = 0,
    SEND_LINK = 2,
    RETURN_CODE = 3,
  }
}

export class PasswordResetResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): PasswordResetResponse;
  hasDetails(): boolean;
  clearDetails(): PasswordResetResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): PasswordResetResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): PasswordResetResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordResetResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordResetResponse): PasswordResetResponse.AsObject;
  static serializeBinaryToWriter(message: PasswordResetResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordResetResponse;
  static deserializeBinaryFromReader(message: PasswordResetResponse, reader: jspb.BinaryReader): PasswordResetResponse;
}

export namespace PasswordResetResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 2,
  }
}

export class SetPasswordRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetPasswordRequest;

  getNewPassword(): zitadel_user_v2beta_password_pb.Password | undefined;
  setNewPassword(value?: zitadel_user_v2beta_password_pb.Password): SetPasswordRequest;
  hasNewPassword(): boolean;
  clearNewPassword(): SetPasswordRequest;

  getCurrentPassword(): string;
  setCurrentPassword(value: string): SetPasswordRequest;

  getVerificationCode(): string;
  setVerificationCode(value: string): SetPasswordRequest;

  getVerificationCase(): SetPasswordRequest.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetPasswordRequest): SetPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: SetPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPasswordRequest;
  static deserializeBinaryFromReader(message: SetPasswordRequest, reader: jspb.BinaryReader): SetPasswordRequest;
}

export namespace SetPasswordRequest {
  export type AsObject = {
    userId: string,
    newPassword?: zitadel_user_v2beta_password_pb.Password.AsObject,
    currentPassword: string,
    verificationCode: string,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    CURRENT_PASSWORD = 3,
    VERIFICATION_CODE = 4,
  }
}

export class SetPasswordResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetPasswordResponse;
  hasDetails(): boolean;
  clearDetails(): SetPasswordResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetPasswordResponse): SetPasswordResponse.AsObject;
  static serializeBinaryToWriter(message: SetPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPasswordResponse;
  static deserializeBinaryFromReader(message: SetPasswordResponse, reader: jspb.BinaryReader): SetPasswordResponse;
}

export namespace SetPasswordResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ListAuthenticationMethodTypesRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ListAuthenticationMethodTypesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAuthenticationMethodTypesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAuthenticationMethodTypesRequest): ListAuthenticationMethodTypesRequest.AsObject;
  static serializeBinaryToWriter(message: ListAuthenticationMethodTypesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAuthenticationMethodTypesRequest;
  static deserializeBinaryFromReader(message: ListAuthenticationMethodTypesRequest, reader: jspb.BinaryReader): ListAuthenticationMethodTypesRequest;
}

export namespace ListAuthenticationMethodTypesRequest {
  export type AsObject = {
    userId: string,
  }
}

export class ListAuthenticationMethodTypesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.ListDetails): ListAuthenticationMethodTypesResponse;
  hasDetails(): boolean;
  clearDetails(): ListAuthenticationMethodTypesResponse;

  getAuthMethodTypesList(): Array<AuthenticationMethodType>;
  setAuthMethodTypesList(value: Array<AuthenticationMethodType>): ListAuthenticationMethodTypesResponse;
  clearAuthMethodTypesList(): ListAuthenticationMethodTypesResponse;
  addAuthMethodTypes(value: AuthenticationMethodType, index?: number): ListAuthenticationMethodTypesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAuthenticationMethodTypesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAuthenticationMethodTypesResponse): ListAuthenticationMethodTypesResponse.AsObject;
  static serializeBinaryToWriter(message: ListAuthenticationMethodTypesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAuthenticationMethodTypesResponse;
  static deserializeBinaryFromReader(message: ListAuthenticationMethodTypesResponse, reader: jspb.BinaryReader): ListAuthenticationMethodTypesResponse;
}

export namespace ListAuthenticationMethodTypesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.ListDetails.AsObject,
    authMethodTypesList: Array<AuthenticationMethodType>,
  }
}

export enum AuthenticationMethodType { 
  AUTHENTICATION_METHOD_TYPE_UNSPECIFIED = 0,
  AUTHENTICATION_METHOD_TYPE_PASSWORD = 1,
  AUTHENTICATION_METHOD_TYPE_PASSKEY = 2,
  AUTHENTICATION_METHOD_TYPE_IDP = 3,
  AUTHENTICATION_METHOD_TYPE_TOTP = 4,
  AUTHENTICATION_METHOD_TYPE_U2F = 5,
  AUTHENTICATION_METHOD_TYPE_OTP_SMS = 6,
  AUTHENTICATION_METHOD_TYPE_OTP_EMAIL = 7,
}

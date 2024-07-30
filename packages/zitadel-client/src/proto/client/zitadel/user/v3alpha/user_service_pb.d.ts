import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_protoc_gen_zitadel_v2_options_pb from '../../../zitadel/protoc_gen_zitadel/v2/options_pb'; // proto import: "zitadel/protoc_gen_zitadel/v2/options.proto"
import * as zitadel_user_v3alpha_authenticator_pb from '../../../zitadel/user/v3alpha/authenticator_pb'; // proto import: "zitadel/user/v3alpha/authenticator.proto"
import * as zitadel_user_v3alpha_communication_pb from '../../../zitadel/user/v3alpha/communication_pb'; // proto import: "zitadel/user/v3alpha/communication.proto"
import * as zitadel_user_v3alpha_query_pb from '../../../zitadel/user/v3alpha/query_pb'; // proto import: "zitadel/user/v3alpha/query.proto"
import * as zitadel_user_v3alpha_user_pb from '../../../zitadel/user/v3alpha/user_pb'; // proto import: "zitadel/user/v3alpha/user.proto"


export class ListUsersRequest extends jspb.Message {
  getQuery(): zitadel_object_v2beta_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_v2beta_object_pb.ListQuery): ListUsersRequest;
  hasQuery(): boolean;
  clearQuery(): ListUsersRequest;

  getSortingColumn(): zitadel_user_v3alpha_query_pb.FieldName;
  setSortingColumn(value: zitadel_user_v3alpha_query_pb.FieldName): ListUsersRequest;

  getQueriesList(): Array<zitadel_user_v3alpha_query_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_user_v3alpha_query_pb.SearchQuery>): ListUsersRequest;
  clearQueriesList(): ListUsersRequest;
  addQueries(value?: zitadel_user_v3alpha_query_pb.SearchQuery, index?: number): zitadel_user_v3alpha_query_pb.SearchQuery;

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
    sortingColumn: zitadel_user_v3alpha_query_pb.FieldName,
    queriesList: Array<zitadel_user_v3alpha_query_pb.SearchQuery.AsObject>,
  }
}

export class ListUsersResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.ListDetails): ListUsersResponse;
  hasDetails(): boolean;
  clearDetails(): ListUsersResponse;

  getSortingColumn(): zitadel_user_v3alpha_query_pb.FieldName;
  setSortingColumn(value: zitadel_user_v3alpha_query_pb.FieldName): ListUsersResponse;

  getResultList(): Array<zitadel_user_v3alpha_user_pb.User>;
  setResultList(value: Array<zitadel_user_v3alpha_user_pb.User>): ListUsersResponse;
  clearResultList(): ListUsersResponse;
  addResult(value?: zitadel_user_v3alpha_user_pb.User, index?: number): zitadel_user_v3alpha_user_pb.User;

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
    sortingColumn: zitadel_user_v3alpha_query_pb.FieldName,
    resultList: Array<zitadel_user_v3alpha_user_pb.User.AsObject>,
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
  getUser(): zitadel_user_v3alpha_user_pb.User | undefined;
  setUser(value?: zitadel_user_v3alpha_user_pb.User): GetUserByIDResponse;
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
    user?: zitadel_user_v3alpha_user_pb.User.AsObject,
  }
}

export class CreateUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): CreateUserRequest;
  hasUserId(): boolean;
  clearUserId(): CreateUserRequest;

  getOrganization(): zitadel_object_v2beta_object_pb.Organization | undefined;
  setOrganization(value?: zitadel_object_v2beta_object_pb.Organization): CreateUserRequest;
  hasOrganization(): boolean;
  clearOrganization(): CreateUserRequest;

  getAuthenticators(): zitadel_user_v3alpha_authenticator_pb.SetAuthenticators | undefined;
  setAuthenticators(value?: zitadel_user_v3alpha_authenticator_pb.SetAuthenticators): CreateUserRequest;
  hasAuthenticators(): boolean;
  clearAuthenticators(): CreateUserRequest;

  getContact(): zitadel_user_v3alpha_communication_pb.SetContact | undefined;
  setContact(value?: zitadel_user_v3alpha_communication_pb.SetContact): CreateUserRequest;
  hasContact(): boolean;
  clearContact(): CreateUserRequest;

  getSchemaId(): string;
  setSchemaId(value: string): CreateUserRequest;

  getData(): google_protobuf_struct_pb.Struct | undefined;
  setData(value?: google_protobuf_struct_pb.Struct): CreateUserRequest;
  hasData(): boolean;
  clearData(): CreateUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUserRequest): CreateUserRequest.AsObject;
  static serializeBinaryToWriter(message: CreateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUserRequest;
  static deserializeBinaryFromReader(message: CreateUserRequest, reader: jspb.BinaryReader): CreateUserRequest;
}

export namespace CreateUserRequest {
  export type AsObject = {
    userId?: string,
    organization?: zitadel_object_v2beta_object_pb.Organization.AsObject,
    authenticators?: zitadel_user_v3alpha_authenticator_pb.SetAuthenticators.AsObject,
    contact?: zitadel_user_v3alpha_communication_pb.SetContact.AsObject,
    schemaId: string,
    data?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export enum UserIdCase { 
    _USER_ID_NOT_SET = 0,
    USER_ID = 1,
  }
}

export class CreateUserResponse extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): CreateUserResponse;

  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): CreateUserResponse;
  hasDetails(): boolean;
  clearDetails(): CreateUserResponse;

  getEmailCode(): string;
  setEmailCode(value: string): CreateUserResponse;
  hasEmailCode(): boolean;
  clearEmailCode(): CreateUserResponse;

  getPhoneCode(): string;
  setPhoneCode(value: string): CreateUserResponse;
  hasPhoneCode(): boolean;
  clearPhoneCode(): CreateUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUserResponse): CreateUserResponse.AsObject;
  static serializeBinaryToWriter(message: CreateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUserResponse;
  static deserializeBinaryFromReader(message: CreateUserResponse, reader: jspb.BinaryReader): CreateUserResponse;
}

export namespace CreateUserResponse {
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

export class UpdateUserRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateUserRequest;

  getContact(): zitadel_user_v3alpha_communication_pb.SetContact | undefined;
  setContact(value?: zitadel_user_v3alpha_communication_pb.SetContact): UpdateUserRequest;
  hasContact(): boolean;
  clearContact(): UpdateUserRequest;

  getSchemaId(): string;
  setSchemaId(value: string): UpdateUserRequest;
  hasSchemaId(): boolean;
  clearSchemaId(): UpdateUserRequest;

  getData(): google_protobuf_struct_pb.Struct | undefined;
  setData(value?: google_protobuf_struct_pb.Struct): UpdateUserRequest;
  hasData(): boolean;
  clearData(): UpdateUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserRequest): UpdateUserRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserRequest;
  static deserializeBinaryFromReader(message: UpdateUserRequest, reader: jspb.BinaryReader): UpdateUserRequest;
}

export namespace UpdateUserRequest {
  export type AsObject = {
    userId: string,
    contact?: zitadel_user_v3alpha_communication_pb.SetContact.AsObject,
    schemaId?: string,
    data?: google_protobuf_struct_pb.Struct.AsObject,
  }

  export enum ContactCase { 
    _CONTACT_NOT_SET = 0,
    CONTACT = 4,
  }

  export enum SchemaIdCase { 
    _SCHEMA_ID_NOT_SET = 0,
    SCHEMA_ID = 5,
  }

  export enum DataCase { 
    _DATA_NOT_SET = 0,
    DATA = 6,
  }
}

export class UpdateUserResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): UpdateUserResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateUserResponse;

  getEmailCode(): string;
  setEmailCode(value: string): UpdateUserResponse;
  hasEmailCode(): boolean;
  clearEmailCode(): UpdateUserResponse;

  getPhoneCode(): string;
  setPhoneCode(value: string): UpdateUserResponse;
  hasPhoneCode(): boolean;
  clearPhoneCode(): UpdateUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserResponse): UpdateUserResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserResponse;
  static deserializeBinaryFromReader(message: UpdateUserResponse, reader: jspb.BinaryReader): UpdateUserResponse;
}

export namespace UpdateUserResponse {
  export type AsObject = {
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

export class SetContactEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetContactEmailRequest;

  getEmail(): zitadel_user_v3alpha_communication_pb.SetEmail | undefined;
  setEmail(value?: zitadel_user_v3alpha_communication_pb.SetEmail): SetContactEmailRequest;
  hasEmail(): boolean;
  clearEmail(): SetContactEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetContactEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetContactEmailRequest): SetContactEmailRequest.AsObject;
  static serializeBinaryToWriter(message: SetContactEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetContactEmailRequest;
  static deserializeBinaryFromReader(message: SetContactEmailRequest, reader: jspb.BinaryReader): SetContactEmailRequest;
}

export namespace SetContactEmailRequest {
  export type AsObject = {
    userId: string,
    email?: zitadel_user_v3alpha_communication_pb.SetEmail.AsObject,
  }
}

export class SetContactEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetContactEmailResponse;
  hasDetails(): boolean;
  clearDetails(): SetContactEmailResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): SetContactEmailResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): SetContactEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetContactEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetContactEmailResponse): SetContactEmailResponse.AsObject;
  static serializeBinaryToWriter(message: SetContactEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetContactEmailResponse;
  static deserializeBinaryFromReader(message: SetContactEmailResponse, reader: jspb.BinaryReader): SetContactEmailResponse;
}

export namespace SetContactEmailResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 3,
  }
}

export class VerifyContactEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyContactEmailRequest;

  getVerificationCode(): string;
  setVerificationCode(value: string): VerifyContactEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyContactEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyContactEmailRequest): VerifyContactEmailRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyContactEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyContactEmailRequest;
  static deserializeBinaryFromReader(message: VerifyContactEmailRequest, reader: jspb.BinaryReader): VerifyContactEmailRequest;
}

export namespace VerifyContactEmailRequest {
  export type AsObject = {
    userId: string,
    verificationCode: string,
  }
}

export class VerifyContactEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyContactEmailResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyContactEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyContactEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyContactEmailResponse): VerifyContactEmailResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyContactEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyContactEmailResponse;
  static deserializeBinaryFromReader(message: VerifyContactEmailResponse, reader: jspb.BinaryReader): VerifyContactEmailResponse;
}

export namespace VerifyContactEmailResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ResendContactEmailCodeRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ResendContactEmailCodeRequest;

  getSendCode(): zitadel_user_v3alpha_communication_pb.SendEmailVerificationCode | undefined;
  setSendCode(value?: zitadel_user_v3alpha_communication_pb.SendEmailVerificationCode): ResendContactEmailCodeRequest;
  hasSendCode(): boolean;
  clearSendCode(): ResendContactEmailCodeRequest;

  getReturnCode(): zitadel_user_v3alpha_communication_pb.ReturnEmailVerificationCode | undefined;
  setReturnCode(value?: zitadel_user_v3alpha_communication_pb.ReturnEmailVerificationCode): ResendContactEmailCodeRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): ResendContactEmailCodeRequest;

  getVerificationCase(): ResendContactEmailCodeRequest.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendContactEmailCodeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendContactEmailCodeRequest): ResendContactEmailCodeRequest.AsObject;
  static serializeBinaryToWriter(message: ResendContactEmailCodeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendContactEmailCodeRequest;
  static deserializeBinaryFromReader(message: ResendContactEmailCodeRequest, reader: jspb.BinaryReader): ResendContactEmailCodeRequest;
}

export namespace ResendContactEmailCodeRequest {
  export type AsObject = {
    userId: string,
    sendCode?: zitadel_user_v3alpha_communication_pb.SendEmailVerificationCode.AsObject,
    returnCode?: zitadel_user_v3alpha_communication_pb.ReturnEmailVerificationCode.AsObject,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    SEND_CODE = 2,
    RETURN_CODE = 3,
  }
}

export class ResendContactEmailCodeResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ResendContactEmailCodeResponse;
  hasDetails(): boolean;
  clearDetails(): ResendContactEmailCodeResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): ResendContactEmailCodeResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): ResendContactEmailCodeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendContactEmailCodeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendContactEmailCodeResponse): ResendContactEmailCodeResponse.AsObject;
  static serializeBinaryToWriter(message: ResendContactEmailCodeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendContactEmailCodeResponse;
  static deserializeBinaryFromReader(message: ResendContactEmailCodeResponse, reader: jspb.BinaryReader): ResendContactEmailCodeResponse;
}

export namespace ResendContactEmailCodeResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 2,
  }
}

export class SetContactPhoneRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetContactPhoneRequest;

  getPhone(): zitadel_user_v3alpha_communication_pb.SetPhone | undefined;
  setPhone(value?: zitadel_user_v3alpha_communication_pb.SetPhone): SetContactPhoneRequest;
  hasPhone(): boolean;
  clearPhone(): SetContactPhoneRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetContactPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetContactPhoneRequest): SetContactPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: SetContactPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetContactPhoneRequest;
  static deserializeBinaryFromReader(message: SetContactPhoneRequest, reader: jspb.BinaryReader): SetContactPhoneRequest;
}

export namespace SetContactPhoneRequest {
  export type AsObject = {
    userId: string,
    phone?: zitadel_user_v3alpha_communication_pb.SetPhone.AsObject,
  }
}

export class SetContactPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetContactPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): SetContactPhoneResponse;

  getEmailCode(): string;
  setEmailCode(value: string): SetContactPhoneResponse;
  hasEmailCode(): boolean;
  clearEmailCode(): SetContactPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetContactPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetContactPhoneResponse): SetContactPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: SetContactPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetContactPhoneResponse;
  static deserializeBinaryFromReader(message: SetContactPhoneResponse, reader: jspb.BinaryReader): SetContactPhoneResponse;
}

export namespace SetContactPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    emailCode?: string,
  }

  export enum EmailCodeCase { 
    _EMAIL_CODE_NOT_SET = 0,
    EMAIL_CODE = 3,
  }
}

export class VerifyContactPhoneRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyContactPhoneRequest;

  getVerificationCode(): string;
  setVerificationCode(value: string): VerifyContactPhoneRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyContactPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyContactPhoneRequest): VerifyContactPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyContactPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyContactPhoneRequest;
  static deserializeBinaryFromReader(message: VerifyContactPhoneRequest, reader: jspb.BinaryReader): VerifyContactPhoneRequest;
}

export namespace VerifyContactPhoneRequest {
  export type AsObject = {
    userId: string,
    verificationCode: string,
  }
}

export class VerifyContactPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyContactPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyContactPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyContactPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyContactPhoneResponse): VerifyContactPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyContactPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyContactPhoneResponse;
  static deserializeBinaryFromReader(message: VerifyContactPhoneResponse, reader: jspb.BinaryReader): VerifyContactPhoneResponse;
}

export namespace VerifyContactPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ResendContactPhoneCodeRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ResendContactPhoneCodeRequest;

  getSendCode(): zitadel_user_v3alpha_communication_pb.SendPhoneVerificationCode | undefined;
  setSendCode(value?: zitadel_user_v3alpha_communication_pb.SendPhoneVerificationCode): ResendContactPhoneCodeRequest;
  hasSendCode(): boolean;
  clearSendCode(): ResendContactPhoneCodeRequest;

  getReturnCode(): zitadel_user_v3alpha_communication_pb.ReturnPhoneVerificationCode | undefined;
  setReturnCode(value?: zitadel_user_v3alpha_communication_pb.ReturnPhoneVerificationCode): ResendContactPhoneCodeRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): ResendContactPhoneCodeRequest;

  getVerificationCase(): ResendContactPhoneCodeRequest.VerificationCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendContactPhoneCodeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendContactPhoneCodeRequest): ResendContactPhoneCodeRequest.AsObject;
  static serializeBinaryToWriter(message: ResendContactPhoneCodeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendContactPhoneCodeRequest;
  static deserializeBinaryFromReader(message: ResendContactPhoneCodeRequest, reader: jspb.BinaryReader): ResendContactPhoneCodeRequest;
}

export namespace ResendContactPhoneCodeRequest {
  export type AsObject = {
    userId: string,
    sendCode?: zitadel_user_v3alpha_communication_pb.SendPhoneVerificationCode.AsObject,
    returnCode?: zitadel_user_v3alpha_communication_pb.ReturnPhoneVerificationCode.AsObject,
  }

  export enum VerificationCase { 
    VERIFICATION_NOT_SET = 0,
    SEND_CODE = 2,
    RETURN_CODE = 3,
  }
}

export class ResendContactPhoneCodeResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ResendContactPhoneCodeResponse;
  hasDetails(): boolean;
  clearDetails(): ResendContactPhoneCodeResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): ResendContactPhoneCodeResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): ResendContactPhoneCodeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendContactPhoneCodeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendContactPhoneCodeResponse): ResendContactPhoneCodeResponse.AsObject;
  static serializeBinaryToWriter(message: ResendContactPhoneCodeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendContactPhoneCodeResponse;
  static deserializeBinaryFromReader(message: ResendContactPhoneCodeResponse, reader: jspb.BinaryReader): ResendContactPhoneCodeResponse;
}

export namespace ResendContactPhoneCodeResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 2,
  }
}

export class AddUsernameRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddUsernameRequest;

  getUsername(): zitadel_user_v3alpha_authenticator_pb.SetUsername | undefined;
  setUsername(value?: zitadel_user_v3alpha_authenticator_pb.SetUsername): AddUsernameRequest;
  hasUsername(): boolean;
  clearUsername(): AddUsernameRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddUsernameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddUsernameRequest): AddUsernameRequest.AsObject;
  static serializeBinaryToWriter(message: AddUsernameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddUsernameRequest;
  static deserializeBinaryFromReader(message: AddUsernameRequest, reader: jspb.BinaryReader): AddUsernameRequest;
}

export namespace AddUsernameRequest {
  export type AsObject = {
    userId: string,
    username?: zitadel_user_v3alpha_authenticator_pb.SetUsername.AsObject,
  }
}

export class AddUsernameResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddUsernameResponse;
  hasDetails(): boolean;
  clearDetails(): AddUsernameResponse;

  getUsernameId(): string;
  setUsernameId(value: string): AddUsernameResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddUsernameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddUsernameResponse): AddUsernameResponse.AsObject;
  static serializeBinaryToWriter(message: AddUsernameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddUsernameResponse;
  static deserializeBinaryFromReader(message: AddUsernameResponse, reader: jspb.BinaryReader): AddUsernameResponse;
}

export namespace AddUsernameResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    usernameId: string,
  }
}

export class RemoveUsernameRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveUsernameRequest;

  getUsernameId(): string;
  setUsernameId(value: string): RemoveUsernameRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUsernameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUsernameRequest): RemoveUsernameRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveUsernameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUsernameRequest;
  static deserializeBinaryFromReader(message: RemoveUsernameRequest, reader: jspb.BinaryReader): RemoveUsernameRequest;
}

export namespace RemoveUsernameRequest {
  export type AsObject = {
    userId: string,
    usernameId: string,
  }
}

export class RemoveUsernameResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveUsernameResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveUsernameResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUsernameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUsernameResponse): RemoveUsernameResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveUsernameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUsernameResponse;
  static deserializeBinaryFromReader(message: RemoveUsernameResponse, reader: jspb.BinaryReader): RemoveUsernameResponse;
}

export namespace RemoveUsernameResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class SetPasswordRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetPasswordRequest;

  getNewPassword(): zitadel_user_v3alpha_authenticator_pb.SetPassword | undefined;
  setNewPassword(value?: zitadel_user_v3alpha_authenticator_pb.SetPassword): SetPasswordRequest;
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
    newPassword?: zitadel_user_v3alpha_authenticator_pb.SetPassword.AsObject,
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

export class RequestPasswordResetRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RequestPasswordResetRequest;

  getSendEmail(): zitadel_user_v3alpha_authenticator_pb.SendPasswordResetEmail | undefined;
  setSendEmail(value?: zitadel_user_v3alpha_authenticator_pb.SendPasswordResetEmail): RequestPasswordResetRequest;
  hasSendEmail(): boolean;
  clearSendEmail(): RequestPasswordResetRequest;

  getSendSms(): zitadel_user_v3alpha_authenticator_pb.SendPasswordResetSMS | undefined;
  setSendSms(value?: zitadel_user_v3alpha_authenticator_pb.SendPasswordResetSMS): RequestPasswordResetRequest;
  hasSendSms(): boolean;
  clearSendSms(): RequestPasswordResetRequest;

  getReturnCode(): zitadel_user_v3alpha_authenticator_pb.ReturnPasswordResetCode | undefined;
  setReturnCode(value?: zitadel_user_v3alpha_authenticator_pb.ReturnPasswordResetCode): RequestPasswordResetRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): RequestPasswordResetRequest;

  getMediumCase(): RequestPasswordResetRequest.MediumCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestPasswordResetRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RequestPasswordResetRequest): RequestPasswordResetRequest.AsObject;
  static serializeBinaryToWriter(message: RequestPasswordResetRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestPasswordResetRequest;
  static deserializeBinaryFromReader(message: RequestPasswordResetRequest, reader: jspb.BinaryReader): RequestPasswordResetRequest;
}

export namespace RequestPasswordResetRequest {
  export type AsObject = {
    userId: string,
    sendEmail?: zitadel_user_v3alpha_authenticator_pb.SendPasswordResetEmail.AsObject,
    sendSms?: zitadel_user_v3alpha_authenticator_pb.SendPasswordResetSMS.AsObject,
    returnCode?: zitadel_user_v3alpha_authenticator_pb.ReturnPasswordResetCode.AsObject,
  }

  export enum MediumCase { 
    MEDIUM_NOT_SET = 0,
    SEND_EMAIL = 2,
    SEND_SMS = 3,
    RETURN_CODE = 4,
  }
}

export class RequestPasswordResetResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RequestPasswordResetResponse;
  hasDetails(): boolean;
  clearDetails(): RequestPasswordResetResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): RequestPasswordResetResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): RequestPasswordResetResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestPasswordResetResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RequestPasswordResetResponse): RequestPasswordResetResponse.AsObject;
  static serializeBinaryToWriter(message: RequestPasswordResetResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestPasswordResetResponse;
  static deserializeBinaryFromReader(message: RequestPasswordResetResponse, reader: jspb.BinaryReader): RequestPasswordResetResponse;
}

export namespace RequestPasswordResetResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 2,
  }
}

export class StartWebAuthNRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): StartWebAuthNRegistrationRequest;

  getDomain(): string;
  setDomain(value: string): StartWebAuthNRegistrationRequest;

  getAuthenticatorType(): zitadel_user_v3alpha_authenticator_pb.WebAuthNAuthenticatorType;
  setAuthenticatorType(value: zitadel_user_v3alpha_authenticator_pb.WebAuthNAuthenticatorType): StartWebAuthNRegistrationRequest;

  getCode(): zitadel_user_v3alpha_authenticator_pb.AuthenticatorRegistrationCode | undefined;
  setCode(value?: zitadel_user_v3alpha_authenticator_pb.AuthenticatorRegistrationCode): StartWebAuthNRegistrationRequest;
  hasCode(): boolean;
  clearCode(): StartWebAuthNRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartWebAuthNRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartWebAuthNRegistrationRequest): StartWebAuthNRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: StartWebAuthNRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartWebAuthNRegistrationRequest;
  static deserializeBinaryFromReader(message: StartWebAuthNRegistrationRequest, reader: jspb.BinaryReader): StartWebAuthNRegistrationRequest;
}

export namespace StartWebAuthNRegistrationRequest {
  export type AsObject = {
    userId: string,
    domain: string,
    authenticatorType: zitadel_user_v3alpha_authenticator_pb.WebAuthNAuthenticatorType,
    code?: zitadel_user_v3alpha_authenticator_pb.AuthenticatorRegistrationCode.AsObject,
  }

  export enum CodeCase { 
    _CODE_NOT_SET = 0,
    CODE = 2,
  }
}

export class StartWebAuthNRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): StartWebAuthNRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): StartWebAuthNRegistrationResponse;

  getWebAuthNId(): string;
  setWebAuthNId(value: string): StartWebAuthNRegistrationResponse;

  getPublicKeyCredentialCreationOptions(): google_protobuf_struct_pb.Struct | undefined;
  setPublicKeyCredentialCreationOptions(value?: google_protobuf_struct_pb.Struct): StartWebAuthNRegistrationResponse;
  hasPublicKeyCredentialCreationOptions(): boolean;
  clearPublicKeyCredentialCreationOptions(): StartWebAuthNRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartWebAuthNRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StartWebAuthNRegistrationResponse): StartWebAuthNRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: StartWebAuthNRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartWebAuthNRegistrationResponse;
  static deserializeBinaryFromReader(message: StartWebAuthNRegistrationResponse, reader: jspb.BinaryReader): StartWebAuthNRegistrationResponse;
}

export namespace StartWebAuthNRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    webAuthNId: string,
    publicKeyCredentialCreationOptions?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class VerifyWebAuthNRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyWebAuthNRegistrationRequest;

  getWebAuthNId(): string;
  setWebAuthNId(value: string): VerifyWebAuthNRegistrationRequest;

  getPublicKeyCredential(): google_protobuf_struct_pb.Struct | undefined;
  setPublicKeyCredential(value?: google_protobuf_struct_pb.Struct): VerifyWebAuthNRegistrationRequest;
  hasPublicKeyCredential(): boolean;
  clearPublicKeyCredential(): VerifyWebAuthNRegistrationRequest;

  getWebAuthNName(): string;
  setWebAuthNName(value: string): VerifyWebAuthNRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyWebAuthNRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyWebAuthNRegistrationRequest): VerifyWebAuthNRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyWebAuthNRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyWebAuthNRegistrationRequest;
  static deserializeBinaryFromReader(message: VerifyWebAuthNRegistrationRequest, reader: jspb.BinaryReader): VerifyWebAuthNRegistrationRequest;
}

export namespace VerifyWebAuthNRegistrationRequest {
  export type AsObject = {
    userId: string,
    webAuthNId: string,
    publicKeyCredential?: google_protobuf_struct_pb.Struct.AsObject,
    webAuthNName: string,
  }
}

export class VerifyWebAuthNRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyWebAuthNRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyWebAuthNRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyWebAuthNRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyWebAuthNRegistrationResponse): VerifyWebAuthNRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyWebAuthNRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyWebAuthNRegistrationResponse;
  static deserializeBinaryFromReader(message: VerifyWebAuthNRegistrationResponse, reader: jspb.BinaryReader): VerifyWebAuthNRegistrationResponse;
}

export namespace VerifyWebAuthNRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class CreateWebAuthNRegistrationLinkRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): CreateWebAuthNRegistrationLinkRequest;

  getSendLink(): zitadel_user_v3alpha_authenticator_pb.SendWebAuthNRegistrationLink | undefined;
  setSendLink(value?: zitadel_user_v3alpha_authenticator_pb.SendWebAuthNRegistrationLink): CreateWebAuthNRegistrationLinkRequest;
  hasSendLink(): boolean;
  clearSendLink(): CreateWebAuthNRegistrationLinkRequest;

  getReturnCode(): zitadel_user_v3alpha_authenticator_pb.ReturnWebAuthNRegistrationCode | undefined;
  setReturnCode(value?: zitadel_user_v3alpha_authenticator_pb.ReturnWebAuthNRegistrationCode): CreateWebAuthNRegistrationLinkRequest;
  hasReturnCode(): boolean;
  clearReturnCode(): CreateWebAuthNRegistrationLinkRequest;

  getMediumCase(): CreateWebAuthNRegistrationLinkRequest.MediumCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateWebAuthNRegistrationLinkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateWebAuthNRegistrationLinkRequest): CreateWebAuthNRegistrationLinkRequest.AsObject;
  static serializeBinaryToWriter(message: CreateWebAuthNRegistrationLinkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateWebAuthNRegistrationLinkRequest;
  static deserializeBinaryFromReader(message: CreateWebAuthNRegistrationLinkRequest, reader: jspb.BinaryReader): CreateWebAuthNRegistrationLinkRequest;
}

export namespace CreateWebAuthNRegistrationLinkRequest {
  export type AsObject = {
    userId: string,
    sendLink?: zitadel_user_v3alpha_authenticator_pb.SendWebAuthNRegistrationLink.AsObject,
    returnCode?: zitadel_user_v3alpha_authenticator_pb.ReturnWebAuthNRegistrationCode.AsObject,
  }

  export enum MediumCase { 
    MEDIUM_NOT_SET = 0,
    SEND_LINK = 2,
    RETURN_CODE = 3,
  }
}

export class CreateWebAuthNRegistrationLinkResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): CreateWebAuthNRegistrationLinkResponse;
  hasDetails(): boolean;
  clearDetails(): CreateWebAuthNRegistrationLinkResponse;

  getCode(): zitadel_user_v3alpha_authenticator_pb.AuthenticatorRegistrationCode | undefined;
  setCode(value?: zitadel_user_v3alpha_authenticator_pb.AuthenticatorRegistrationCode): CreateWebAuthNRegistrationLinkResponse;
  hasCode(): boolean;
  clearCode(): CreateWebAuthNRegistrationLinkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateWebAuthNRegistrationLinkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateWebAuthNRegistrationLinkResponse): CreateWebAuthNRegistrationLinkResponse.AsObject;
  static serializeBinaryToWriter(message: CreateWebAuthNRegistrationLinkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateWebAuthNRegistrationLinkResponse;
  static deserializeBinaryFromReader(message: CreateWebAuthNRegistrationLinkResponse, reader: jspb.BinaryReader): CreateWebAuthNRegistrationLinkResponse;
}

export namespace CreateWebAuthNRegistrationLinkResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    code?: zitadel_user_v3alpha_authenticator_pb.AuthenticatorRegistrationCode.AsObject,
  }

  export enum CodeCase { 
    _CODE_NOT_SET = 0,
    CODE = 2,
  }
}

export class RemoveWebAuthNAuthenticatorRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveWebAuthNAuthenticatorRequest;

  getWebAuthNId(): string;
  setWebAuthNId(value: string): RemoveWebAuthNAuthenticatorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveWebAuthNAuthenticatorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveWebAuthNAuthenticatorRequest): RemoveWebAuthNAuthenticatorRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveWebAuthNAuthenticatorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveWebAuthNAuthenticatorRequest;
  static deserializeBinaryFromReader(message: RemoveWebAuthNAuthenticatorRequest, reader: jspb.BinaryReader): RemoveWebAuthNAuthenticatorRequest;
}

export namespace RemoveWebAuthNAuthenticatorRequest {
  export type AsObject = {
    userId: string,
    webAuthNId: string,
  }
}

export class RemoveWebAuthNAuthenticatorResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveWebAuthNAuthenticatorResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveWebAuthNAuthenticatorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveWebAuthNAuthenticatorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveWebAuthNAuthenticatorResponse): RemoveWebAuthNAuthenticatorResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveWebAuthNAuthenticatorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveWebAuthNAuthenticatorResponse;
  static deserializeBinaryFromReader(message: RemoveWebAuthNAuthenticatorResponse, reader: jspb.BinaryReader): RemoveWebAuthNAuthenticatorResponse;
}

export namespace RemoveWebAuthNAuthenticatorResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class StartTOTPRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): StartTOTPRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartTOTPRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartTOTPRegistrationRequest): StartTOTPRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: StartTOTPRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartTOTPRegistrationRequest;
  static deserializeBinaryFromReader(message: StartTOTPRegistrationRequest, reader: jspb.BinaryReader): StartTOTPRegistrationRequest;
}

export namespace StartTOTPRegistrationRequest {
  export type AsObject = {
    userId: string,
  }
}

export class StartTOTPRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): StartTOTPRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): StartTOTPRegistrationResponse;

  getTotpId(): string;
  setTotpId(value: string): StartTOTPRegistrationResponse;

  getUri(): string;
  setUri(value: string): StartTOTPRegistrationResponse;

  getSecret(): string;
  setSecret(value: string): StartTOTPRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartTOTPRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: StartTOTPRegistrationResponse): StartTOTPRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: StartTOTPRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartTOTPRegistrationResponse;
  static deserializeBinaryFromReader(message: StartTOTPRegistrationResponse, reader: jspb.BinaryReader): StartTOTPRegistrationResponse;
}

export namespace StartTOTPRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    totpId: string,
    uri: string,
    secret: string,
  }
}

export class VerifyTOTPRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyTOTPRegistrationRequest;

  getTotpId(): string;
  setTotpId(value: string): VerifyTOTPRegistrationRequest;

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
    totpId: string,
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

export class RemoveTOTPAuthenticatorRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveTOTPAuthenticatorRequest;

  getTotpId(): string;
  setTotpId(value: string): RemoveTOTPAuthenticatorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveTOTPAuthenticatorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveTOTPAuthenticatorRequest): RemoveTOTPAuthenticatorRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveTOTPAuthenticatorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveTOTPAuthenticatorRequest;
  static deserializeBinaryFromReader(message: RemoveTOTPAuthenticatorRequest, reader: jspb.BinaryReader): RemoveTOTPAuthenticatorRequest;
}

export namespace RemoveTOTPAuthenticatorRequest {
  export type AsObject = {
    userId: string,
    totpId: string,
  }
}

export class RemoveTOTPAuthenticatorResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveTOTPAuthenticatorResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveTOTPAuthenticatorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveTOTPAuthenticatorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveTOTPAuthenticatorResponse): RemoveTOTPAuthenticatorResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveTOTPAuthenticatorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveTOTPAuthenticatorResponse;
  static deserializeBinaryFromReader(message: RemoveTOTPAuthenticatorResponse, reader: jspb.BinaryReader): RemoveTOTPAuthenticatorResponse;
}

export namespace RemoveTOTPAuthenticatorResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class AddOTPSMSAuthenticatorRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddOTPSMSAuthenticatorRequest;

  getPhone(): zitadel_user_v3alpha_communication_pb.SetPhone | undefined;
  setPhone(value?: zitadel_user_v3alpha_communication_pb.SetPhone): AddOTPSMSAuthenticatorRequest;
  hasPhone(): boolean;
  clearPhone(): AddOTPSMSAuthenticatorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOTPSMSAuthenticatorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOTPSMSAuthenticatorRequest): AddOTPSMSAuthenticatorRequest.AsObject;
  static serializeBinaryToWriter(message: AddOTPSMSAuthenticatorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOTPSMSAuthenticatorRequest;
  static deserializeBinaryFromReader(message: AddOTPSMSAuthenticatorRequest, reader: jspb.BinaryReader): AddOTPSMSAuthenticatorRequest;
}

export namespace AddOTPSMSAuthenticatorRequest {
  export type AsObject = {
    userId: string,
    phone?: zitadel_user_v3alpha_communication_pb.SetPhone.AsObject,
  }
}

export class AddOTPSMSAuthenticatorResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddOTPSMSAuthenticatorResponse;
  hasDetails(): boolean;
  clearDetails(): AddOTPSMSAuthenticatorResponse;

  getOtpSmsId(): string;
  setOtpSmsId(value: string): AddOTPSMSAuthenticatorResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): AddOTPSMSAuthenticatorResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): AddOTPSMSAuthenticatorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOTPSMSAuthenticatorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOTPSMSAuthenticatorResponse): AddOTPSMSAuthenticatorResponse.AsObject;
  static serializeBinaryToWriter(message: AddOTPSMSAuthenticatorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOTPSMSAuthenticatorResponse;
  static deserializeBinaryFromReader(message: AddOTPSMSAuthenticatorResponse, reader: jspb.BinaryReader): AddOTPSMSAuthenticatorResponse;
}

export namespace AddOTPSMSAuthenticatorResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    otpSmsId: string,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 3,
  }
}

export class VerifyOTPSMSRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyOTPSMSRegistrationRequest;

  getOtpSmsId(): string;
  setOtpSmsId(value: string): VerifyOTPSMSRegistrationRequest;

  getCode(): string;
  setCode(value: string): VerifyOTPSMSRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyOTPSMSRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyOTPSMSRegistrationRequest): VerifyOTPSMSRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyOTPSMSRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyOTPSMSRegistrationRequest;
  static deserializeBinaryFromReader(message: VerifyOTPSMSRegistrationRequest, reader: jspb.BinaryReader): VerifyOTPSMSRegistrationRequest;
}

export namespace VerifyOTPSMSRegistrationRequest {
  export type AsObject = {
    userId: string,
    otpSmsId: string,
    code: string,
  }
}

export class VerifyOTPSMSRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyOTPSMSRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyOTPSMSRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyOTPSMSRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyOTPSMSRegistrationResponse): VerifyOTPSMSRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyOTPSMSRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyOTPSMSRegistrationResponse;
  static deserializeBinaryFromReader(message: VerifyOTPSMSRegistrationResponse, reader: jspb.BinaryReader): VerifyOTPSMSRegistrationResponse;
}

export namespace VerifyOTPSMSRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RemoveOTPSMSAuthenticatorRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveOTPSMSAuthenticatorRequest;

  getOtpSmsId(): string;
  setOtpSmsId(value: string): RemoveOTPSMSAuthenticatorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOTPSMSAuthenticatorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOTPSMSAuthenticatorRequest): RemoveOTPSMSAuthenticatorRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOTPSMSAuthenticatorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOTPSMSAuthenticatorRequest;
  static deserializeBinaryFromReader(message: RemoveOTPSMSAuthenticatorRequest, reader: jspb.BinaryReader): RemoveOTPSMSAuthenticatorRequest;
}

export namespace RemoveOTPSMSAuthenticatorRequest {
  export type AsObject = {
    userId: string,
    otpSmsId: string,
  }
}

export class RemoveOTPSMSAuthenticatorResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveOTPSMSAuthenticatorResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOTPSMSAuthenticatorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOTPSMSAuthenticatorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOTPSMSAuthenticatorResponse): RemoveOTPSMSAuthenticatorResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOTPSMSAuthenticatorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOTPSMSAuthenticatorResponse;
  static deserializeBinaryFromReader(message: RemoveOTPSMSAuthenticatorResponse, reader: jspb.BinaryReader): RemoveOTPSMSAuthenticatorResponse;
}

export namespace RemoveOTPSMSAuthenticatorResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class AddOTPEmailAuthenticatorRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddOTPEmailAuthenticatorRequest;

  getEmail(): zitadel_user_v3alpha_communication_pb.SetEmail | undefined;
  setEmail(value?: zitadel_user_v3alpha_communication_pb.SetEmail): AddOTPEmailAuthenticatorRequest;
  hasEmail(): boolean;
  clearEmail(): AddOTPEmailAuthenticatorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOTPEmailAuthenticatorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOTPEmailAuthenticatorRequest): AddOTPEmailAuthenticatorRequest.AsObject;
  static serializeBinaryToWriter(message: AddOTPEmailAuthenticatorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOTPEmailAuthenticatorRequest;
  static deserializeBinaryFromReader(message: AddOTPEmailAuthenticatorRequest, reader: jspb.BinaryReader): AddOTPEmailAuthenticatorRequest;
}

export namespace AddOTPEmailAuthenticatorRequest {
  export type AsObject = {
    userId: string,
    email?: zitadel_user_v3alpha_communication_pb.SetEmail.AsObject,
  }
}

export class AddOTPEmailAuthenticatorResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddOTPEmailAuthenticatorResponse;
  hasDetails(): boolean;
  clearDetails(): AddOTPEmailAuthenticatorResponse;

  getOtpEmailId(): string;
  setOtpEmailId(value: string): AddOTPEmailAuthenticatorResponse;

  getVerificationCode(): string;
  setVerificationCode(value: string): AddOTPEmailAuthenticatorResponse;
  hasVerificationCode(): boolean;
  clearVerificationCode(): AddOTPEmailAuthenticatorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOTPEmailAuthenticatorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOTPEmailAuthenticatorResponse): AddOTPEmailAuthenticatorResponse.AsObject;
  static serializeBinaryToWriter(message: AddOTPEmailAuthenticatorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOTPEmailAuthenticatorResponse;
  static deserializeBinaryFromReader(message: AddOTPEmailAuthenticatorResponse, reader: jspb.BinaryReader): AddOTPEmailAuthenticatorResponse;
}

export namespace AddOTPEmailAuthenticatorResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    otpEmailId: string,
    verificationCode?: string,
  }

  export enum VerificationCodeCase { 
    _VERIFICATION_CODE_NOT_SET = 0,
    VERIFICATION_CODE = 3,
  }
}

export class VerifyOTPEmailRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): VerifyOTPEmailRegistrationRequest;

  getOtpEmailId(): string;
  setOtpEmailId(value: string): VerifyOTPEmailRegistrationRequest;

  getCode(): string;
  setCode(value: string): VerifyOTPEmailRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyOTPEmailRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyOTPEmailRegistrationRequest): VerifyOTPEmailRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyOTPEmailRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyOTPEmailRegistrationRequest;
  static deserializeBinaryFromReader(message: VerifyOTPEmailRegistrationRequest, reader: jspb.BinaryReader): VerifyOTPEmailRegistrationRequest;
}

export namespace VerifyOTPEmailRegistrationRequest {
  export type AsObject = {
    userId: string,
    otpEmailId: string,
    code: string,
  }
}

export class VerifyOTPEmailRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): VerifyOTPEmailRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyOTPEmailRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyOTPEmailRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyOTPEmailRegistrationResponse): VerifyOTPEmailRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyOTPEmailRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyOTPEmailRegistrationResponse;
  static deserializeBinaryFromReader(message: VerifyOTPEmailRegistrationResponse, reader: jspb.BinaryReader): VerifyOTPEmailRegistrationResponse;
}

export namespace VerifyOTPEmailRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RemoveOTPEmailAuthenticatorRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveOTPEmailAuthenticatorRequest;

  getOtpEmailId(): string;
  setOtpEmailId(value: string): RemoveOTPEmailAuthenticatorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOTPEmailAuthenticatorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOTPEmailAuthenticatorRequest): RemoveOTPEmailAuthenticatorRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOTPEmailAuthenticatorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOTPEmailAuthenticatorRequest;
  static deserializeBinaryFromReader(message: RemoveOTPEmailAuthenticatorRequest, reader: jspb.BinaryReader): RemoveOTPEmailAuthenticatorRequest;
}

export namespace RemoveOTPEmailAuthenticatorRequest {
  export type AsObject = {
    userId: string,
    otpEmailId: string,
  }
}

export class RemoveOTPEmailAuthenticatorResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveOTPEmailAuthenticatorResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOTPEmailAuthenticatorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOTPEmailAuthenticatorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOTPEmailAuthenticatorResponse): RemoveOTPEmailAuthenticatorResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOTPEmailAuthenticatorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOTPEmailAuthenticatorResponse;
  static deserializeBinaryFromReader(message: RemoveOTPEmailAuthenticatorResponse, reader: jspb.BinaryReader): RemoveOTPEmailAuthenticatorResponse;
}

export namespace RemoveOTPEmailAuthenticatorResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class StartIdentityProviderIntentRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): StartIdentityProviderIntentRequest;

  getUrls(): zitadel_user_v3alpha_authenticator_pb.RedirectURLs | undefined;
  setUrls(value?: zitadel_user_v3alpha_authenticator_pb.RedirectURLs): StartIdentityProviderIntentRequest;
  hasUrls(): boolean;
  clearUrls(): StartIdentityProviderIntentRequest;

  getLdap(): zitadel_user_v3alpha_authenticator_pb.LDAPCredentials | undefined;
  setLdap(value?: zitadel_user_v3alpha_authenticator_pb.LDAPCredentials): StartIdentityProviderIntentRequest;
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
    urls?: zitadel_user_v3alpha_authenticator_pb.RedirectURLs.AsObject,
    ldap?: zitadel_user_v3alpha_authenticator_pb.LDAPCredentials.AsObject,
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

  getIdpIntent(): zitadel_user_v3alpha_authenticator_pb.IdentityProviderIntent | undefined;
  setIdpIntent(value?: zitadel_user_v3alpha_authenticator_pb.IdentityProviderIntent): StartIdentityProviderIntentResponse;
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
    idpIntent?: zitadel_user_v3alpha_authenticator_pb.IdentityProviderIntent.AsObject,
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

  getIdpInformation(): zitadel_user_v3alpha_authenticator_pb.IDPInformation | undefined;
  setIdpInformation(value?: zitadel_user_v3alpha_authenticator_pb.IDPInformation): RetrieveIdentityProviderIntentResponse;
  hasIdpInformation(): boolean;
  clearIdpInformation(): RetrieveIdentityProviderIntentResponse;

  getUserId(): string;
  setUserId(value: string): RetrieveIdentityProviderIntentResponse;
  hasUserId(): boolean;
  clearUserId(): RetrieveIdentityProviderIntentResponse;

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
    idpInformation?: zitadel_user_v3alpha_authenticator_pb.IDPInformation.AsObject,
    userId?: string,
  }

  export enum UserIdCase { 
    _USER_ID_NOT_SET = 0,
    USER_ID = 3,
  }
}

export class AddIDPAuthenticatorRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddIDPAuthenticatorRequest;

  getIdpAuthenticator(): zitadel_user_v3alpha_authenticator_pb.IDPAuthenticator | undefined;
  setIdpAuthenticator(value?: zitadel_user_v3alpha_authenticator_pb.IDPAuthenticator): AddIDPAuthenticatorRequest;
  hasIdpAuthenticator(): boolean;
  clearIdpAuthenticator(): AddIDPAuthenticatorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIDPAuthenticatorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddIDPAuthenticatorRequest): AddIDPAuthenticatorRequest.AsObject;
  static serializeBinaryToWriter(message: AddIDPAuthenticatorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIDPAuthenticatorRequest;
  static deserializeBinaryFromReader(message: AddIDPAuthenticatorRequest, reader: jspb.BinaryReader): AddIDPAuthenticatorRequest;
}

export namespace AddIDPAuthenticatorRequest {
  export type AsObject = {
    userId: string,
    idpAuthenticator?: zitadel_user_v3alpha_authenticator_pb.IDPAuthenticator.AsObject,
  }
}

export class AddIDPAuthenticatorResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddIDPAuthenticatorResponse;
  hasDetails(): boolean;
  clearDetails(): AddIDPAuthenticatorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIDPAuthenticatorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddIDPAuthenticatorResponse): AddIDPAuthenticatorResponse.AsObject;
  static serializeBinaryToWriter(message: AddIDPAuthenticatorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIDPAuthenticatorResponse;
  static deserializeBinaryFromReader(message: AddIDPAuthenticatorResponse, reader: jspb.BinaryReader): AddIDPAuthenticatorResponse;
}

export namespace AddIDPAuthenticatorResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class RemoveIDPAuthenticatorRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveIDPAuthenticatorRequest;

  getIdpId(): string;
  setIdpId(value: string): RemoveIDPAuthenticatorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIDPAuthenticatorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIDPAuthenticatorRequest): RemoveIDPAuthenticatorRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveIDPAuthenticatorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIDPAuthenticatorRequest;
  static deserializeBinaryFromReader(message: RemoveIDPAuthenticatorRequest, reader: jspb.BinaryReader): RemoveIDPAuthenticatorRequest;
}

export namespace RemoveIDPAuthenticatorRequest {
  export type AsObject = {
    userId: string,
    idpId: string,
  }
}

export class RemoveIDPAuthenticatorResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): RemoveIDPAuthenticatorResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveIDPAuthenticatorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIDPAuthenticatorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIDPAuthenticatorResponse): RemoveIDPAuthenticatorResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveIDPAuthenticatorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIDPAuthenticatorResponse;
  static deserializeBinaryFromReader(message: RemoveIDPAuthenticatorResponse, reader: jspb.BinaryReader): RemoveIDPAuthenticatorResponse;
}

export namespace RemoveIDPAuthenticatorResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}


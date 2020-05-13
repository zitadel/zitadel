import * as grpcWeb from 'grpc-web';

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
import * as authoption_options_pb from './authoption/options_pb';

import {
  IsAdminResponse,
  MfaOtpResponse,
  MultiFactors,
  MyPermissions,
  MyProjectOrgSearchRequest,
  MyProjectOrgSearchResponse,
  PasswordChange,
  PasswordRequest,
  UpdateUserAddressRequest,
  UpdateUserEmailRequest,
  UpdateUserPhoneRequest,
  UpdateUserProfileRequest,
  UserAddress,
  UserEmail,
  UserPhone,
  UserProfile,
  UserSessionViews,
  VerifyMfaOtp,
  VerifyMyUserEmailRequest,
  VerifyUserPhoneRequest} from './auth_pb';

export class AuthServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  healthz(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  ready(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  validate(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_struct_pb.Struct) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_struct_pb.Struct>;

  getMyUserSessions(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserSessionViews) => void
  ): grpcWeb.ClientReadableStream<UserSessionViews>;

  getMyUserProfile(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserProfile) => void
  ): grpcWeb.ClientReadableStream<UserProfile>;

  updateMyUserProfile(
    request: UpdateUserProfileRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserProfile) => void
  ): grpcWeb.ClientReadableStream<UserProfile>;

  getMyUserEmail(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserEmail) => void
  ): grpcWeb.ClientReadableStream<UserEmail>;

  changeMyUserEmail(
    request: UpdateUserEmailRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserEmail) => void
  ): grpcWeb.ClientReadableStream<UserEmail>;

  verifyMyUserEmail(
    request: VerifyMyUserEmailRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  resendMyEmailVerificationMail(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getMyUserPhone(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserPhone) => void
  ): grpcWeb.ClientReadableStream<UserPhone>;

  changeMyUserPhone(
    request: UpdateUserPhoneRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserPhone) => void
  ): grpcWeb.ClientReadableStream<UserPhone>;

  verifyMyUserPhone(
    request: VerifyUserPhoneRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  resendMyPhoneVerificationCode(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getMyUserAddress(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserAddress) => void
  ): grpcWeb.ClientReadableStream<UserAddress>;

  updateMyUserAddress(
    request: UpdateUserAddressRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserAddress) => void
  ): grpcWeb.ClientReadableStream<UserAddress>;

  getMyMfas(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MultiFactors) => void
  ): grpcWeb.ClientReadableStream<MultiFactors>;

  setMyPassword(
    request: PasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  changeMyPassword(
    request: PasswordChange,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  addMfaOTP(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MfaOtpResponse) => void
  ): grpcWeb.ClientReadableStream<MfaOtpResponse>;

  verifyMfaOTP(
    request: VerifyMfaOtp,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MfaOtpResponse) => void
  ): grpcWeb.ClientReadableStream<MfaOtpResponse>;

  removeMfaOTP(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  searchMyProjectOrgs(
    request: MyProjectOrgSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MyProjectOrgSearchResponse) => void
  ): grpcWeb.ClientReadableStream<MyProjectOrgSearchResponse>;

  isIamAdmin(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: IsAdminResponse) => void
  ): grpcWeb.ClientReadableStream<IsAdminResponse>;

  getMyZitadelPermissions(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MyPermissions) => void
  ): grpcWeb.ClientReadableStream<MyPermissions>;

}

export class AuthServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  healthz(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  ready(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  validate(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_struct_pb.Struct>;

  getMyUserSessions(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<UserSessionViews>;

  getMyUserProfile(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<UserProfile>;

  updateMyUserProfile(
    request: UpdateUserProfileRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserProfile>;

  getMyUserEmail(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<UserEmail>;

  changeMyUserEmail(
    request: UpdateUserEmailRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserEmail>;

  verifyMyUserEmail(
    request: VerifyMyUserEmailRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  resendMyEmailVerificationMail(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getMyUserPhone(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<UserPhone>;

  changeMyUserPhone(
    request: UpdateUserPhoneRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserPhone>;

  verifyMyUserPhone(
    request: VerifyUserPhoneRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  resendMyPhoneVerificationCode(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getMyUserAddress(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<UserAddress>;

  updateMyUserAddress(
    request: UpdateUserAddressRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserAddress>;

  getMyMfas(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<MultiFactors>;

  setMyPassword(
    request: PasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  changeMyPassword(
    request: PasswordChange,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  addMfaOTP(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<MfaOtpResponse>;

  verifyMfaOTP(
    request: VerifyMfaOtp,
    metadata?: grpcWeb.Metadata
  ): Promise<MfaOtpResponse>;

  removeMfaOTP(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  searchMyProjectOrgs(
    request: MyProjectOrgSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<MyProjectOrgSearchResponse>;

  isIamAdmin(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<IsAdminResponse>;

  getMyZitadelPermissions(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<MyPermissions>;

}


import * as grpcWeb from 'grpc-web';

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as authoption_options_pb from './authoption/options_pb';

import {
  Application,
  ApplicationAuthorizeRequest,
  ApplicationID,
  ApplicationSearchRequest,
  ApplicationSearchResponse,
  AuthSessionCreation,
  AuthSessionID,
  AuthSessionResponse,
  AuthSessionView,
  CreateTokenRequest,
  GrantSearchRequest,
  GrantSearchResponse,
  IsAdminResponse,
  MfaOtpResponse,
  MultiFactors,
  MyPermissions,
  MyProjectOrgSearchRequest,
  MyProjectOrgSearchResponse,
  PasswordChange,
  PasswordRequest,
  RegisterUserExternalIDPRequest,
  RegisterUserRequest,
  ResetPassword,
  ResetPasswordRequest,
  SelectUserRequest,
  SkipMfaInitRequest,
  Token,
  TokenID,
  UniqueUserRequest,
  UniqueUserResponse,
  UpdateUserAddressRequest,
  UpdateUserEmailRequest,
  UpdateUserPhoneRequest,
  UpdateUserProfileRequest,
  User,
  UserAddress,
  UserAgent,
  UserAgentCreation,
  UserAgentID,
  UserEmail,
  UserID,
  UserPhone,
  UserProfile,
  UserSession,
  UserSessionID,
  UserSessionViews,
  UserSessions,
  VerifyMfaOtp,
  VerifyMfaRequest,
  VerifyMyUserEmailRequest,
  VerifyPasswordRequest,
  VerifyUserEmailRequest,
  VerifyUserInitRequest,
  VerifyUserPhoneRequest,
  VerifyUserRequest} from './auth_pb';

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

  getUserAgent(
    request: UserAgentID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserAgent) => void
  ): grpcWeb.ClientReadableStream<UserAgent>;

  createUserAgent(
    request: UserAgentCreation,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserAgent) => void
  ): grpcWeb.ClientReadableStream<UserAgent>;

  revokeUserAgent(
    request: UserAgentID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserAgent) => void
  ): grpcWeb.ClientReadableStream<UserAgent>;

  createAuthSession(
    request: AuthSessionCreation,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthSessionResponse) => void
  ): grpcWeb.ClientReadableStream<AuthSessionResponse>;

  getAuthSession(
    request: AuthSessionID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthSessionResponse) => void
  ): grpcWeb.ClientReadableStream<AuthSessionResponse>;

  getAuthSessionByTokenID(
    request: TokenID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthSessionView) => void
  ): grpcWeb.ClientReadableStream<AuthSessionView>;

  selectUser(
    request: SelectUserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthSessionResponse) => void
  ): grpcWeb.ClientReadableStream<AuthSessionResponse>;

  verifyUser(
    request: VerifyUserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthSessionResponse) => void
  ): grpcWeb.ClientReadableStream<AuthSessionResponse>;

  verifyPassword(
    request: VerifyPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthSessionResponse) => void
  ): grpcWeb.ClientReadableStream<AuthSessionResponse>;

  verifyMfa(
    request: VerifyMfaRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthSessionResponse) => void
  ): grpcWeb.ClientReadableStream<AuthSessionResponse>;

  getUserAgentSessions(
    request: UserAgentID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserSessions) => void
  ): grpcWeb.ClientReadableStream<UserSessions>;

  getUserSession(
    request: UserSessionID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserSession) => void
  ): grpcWeb.ClientReadableStream<UserSession>;

  getMyUserSessions(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserSessionViews) => void
  ): grpcWeb.ClientReadableStream<UserSessionViews>;

  terminateUserSession(
    request: UserSessionID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  createToken(
    request: CreateTokenRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Token) => void
  ): grpcWeb.ClientReadableStream<Token>;

  isUserUnique(
    request: UniqueUserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UniqueUserResponse) => void
  ): grpcWeb.ClientReadableStream<UniqueUserResponse>;

  registerUser(
    request: RegisterUserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: User) => void
  ): grpcWeb.ClientReadableStream<User>;

  registerUserWithExternal(
    request: RegisterUserExternalIDPRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: User) => void
  ): grpcWeb.ClientReadableStream<User>;

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

  verifyUserEmail(
    request: VerifyUserEmailRequest,
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

  resendEmailVerificationMail(
    request: UserID,
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

  requestPasswordReset(
    request: ResetPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  passwordReset(
    request: ResetPassword,
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

  skipMfaInit(
    request: SkipMfaInitRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  verifyUserInit(
    request: VerifyUserInitRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  resendUserInitMail(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getUserByID(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: User) => void
  ): grpcWeb.ClientReadableStream<User>;

  getApplicationByID(
    request: ApplicationID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Application) => void
  ): grpcWeb.ClientReadableStream<Application>;

  searchApplications(
    request: ApplicationSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ApplicationSearchResponse) => void
  ): grpcWeb.ClientReadableStream<ApplicationSearchResponse>;

  authorizeApplication(
    request: ApplicationAuthorizeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Application) => void
  ): grpcWeb.ClientReadableStream<Application>;

  searchGrant(
    request: GrantSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: GrantSearchResponse) => void
  ): grpcWeb.ClientReadableStream<GrantSearchResponse>;

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

  getMyCitadelPermissions(
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

  getUserAgent(
    request: UserAgentID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserAgent>;

  createUserAgent(
    request: UserAgentCreation,
    metadata?: grpcWeb.Metadata
  ): Promise<UserAgent>;

  revokeUserAgent(
    request: UserAgentID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserAgent>;

  createAuthSession(
    request: AuthSessionCreation,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthSessionResponse>;

  getAuthSession(
    request: AuthSessionID,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthSessionResponse>;

  getAuthSessionByTokenID(
    request: TokenID,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthSessionView>;

  selectUser(
    request: SelectUserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthSessionResponse>;

  verifyUser(
    request: VerifyUserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthSessionResponse>;

  verifyPassword(
    request: VerifyPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthSessionResponse>;

  verifyMfa(
    request: VerifyMfaRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthSessionResponse>;

  getUserAgentSessions(
    request: UserAgentID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserSessions>;

  getUserSession(
    request: UserSessionID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserSession>;

  getMyUserSessions(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<UserSessionViews>;

  terminateUserSession(
    request: UserSessionID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  createToken(
    request: CreateTokenRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<Token>;

  isUserUnique(
    request: UniqueUserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UniqueUserResponse>;

  registerUser(
    request: RegisterUserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<User>;

  registerUserWithExternal(
    request: RegisterUserExternalIDPRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<User>;

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

  verifyUserEmail(
    request: VerifyUserEmailRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  resendMyEmailVerificationMail(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  resendEmailVerificationMail(
    request: UserID,
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

  requestPasswordReset(
    request: ResetPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  passwordReset(
    request: ResetPassword,
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

  skipMfaInit(
    request: SkipMfaInitRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  verifyUserInit(
    request: VerifyUserInitRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  resendUserInitMail(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getUserByID(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<User>;

  getApplicationByID(
    request: ApplicationID,
    metadata?: grpcWeb.Metadata
  ): Promise<Application>;

  searchApplications(
    request: ApplicationSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ApplicationSearchResponse>;

  authorizeApplication(
    request: ApplicationAuthorizeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<Application>;

  searchGrant(
    request: GrantSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<GrantSearchResponse>;

  searchMyProjectOrgs(
    request: MyProjectOrgSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<MyProjectOrgSearchResponse>;

  isIamAdmin(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<IsAdminResponse>;

  getMyCitadelPermissions(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<MyPermissions>;

}


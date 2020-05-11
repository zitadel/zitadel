import * as jspb from "google-protobuf"

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as authoption_options_pb from './authoption/options_pb';

export class SessionRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SessionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SessionRequest): SessionRequest.AsObject;
  static serializeBinaryToWriter(message: SessionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SessionRequest;
  static deserializeBinaryFromReader(message: SessionRequest, reader: jspb.BinaryReader): SessionRequest;
}

export namespace SessionRequest {
  export type AsObject = {
    userId: string,
    browserInfo?: BrowserInformation.AsObject,
  }
}

export class UserAgent extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  getState(): UserAgentState;
  setState(value: UserAgentState): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAgent.AsObject;
  static toObject(includeInstance: boolean, msg: UserAgent): UserAgent.AsObject;
  static serializeBinaryToWriter(message: UserAgent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserAgent;
  static deserializeBinaryFromReader(message: UserAgent, reader: jspb.BinaryReader): UserAgent;
}

export namespace UserAgent {
  export type AsObject = {
    id: string,
    browserInfo?: BrowserInformation.AsObject,
    state: UserAgentState,
  }
}

export class UserAgentID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAgentID.AsObject;
  static toObject(includeInstance: boolean, msg: UserAgentID): UserAgentID.AsObject;
  static serializeBinaryToWriter(message: UserAgentID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserAgentID;
  static deserializeBinaryFromReader(message: UserAgentID, reader: jspb.BinaryReader): UserAgentID;
}

export namespace UserAgentID {
  export type AsObject = {
    id: string,
  }
}

export class UserAgentCreation extends jspb.Message {
  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAgentCreation.AsObject;
  static toObject(includeInstance: boolean, msg: UserAgentCreation): UserAgentCreation.AsObject;
  static serializeBinaryToWriter(message: UserAgentCreation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserAgentCreation;
  static deserializeBinaryFromReader(message: UserAgentCreation, reader: jspb.BinaryReader): UserAgentCreation;
}

export namespace UserAgentCreation {
  export type AsObject = {
    browserInfo?: BrowserInformation.AsObject,
  }
}

export class UserAgents extends jspb.Message {
  getSessionsList(): Array<UserAgent>;
  setSessionsList(value: Array<UserAgent>): void;
  clearSessionsList(): void;
  addSessions(value?: UserAgent, index?: number): UserAgent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAgents.AsObject;
  static toObject(includeInstance: boolean, msg: UserAgents): UserAgents.AsObject;
  static serializeBinaryToWriter(message: UserAgents, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserAgents;
  static deserializeBinaryFromReader(message: UserAgents, reader: jspb.BinaryReader): UserAgents;
}

export namespace UserAgents {
  export type AsObject = {
    sessionsList: Array<UserAgent.AsObject>,
  }
}

export class AuthSessionCreation extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): void;

  getType(): AuthSessionType;
  setType(value: AuthSessionType): void;

  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  getClientId(): string;
  setClientId(value: string): void;

  getRedirectUri(): string;
  setRedirectUri(value: string): void;

  getState(): string;
  setState(value: string): void;

  getPrompt(): Prompt;
  setPrompt(value: Prompt): void;

  getAuthContextClassReferenceList(): Array<string>;
  setAuthContextClassReferenceList(value: Array<string>): void;
  clearAuthContextClassReferenceList(): void;
  addAuthContextClassReference(value: string, index?: number): void;

  getUiLocalesList(): Array<string>;
  setUiLocalesList(value: Array<string>): void;
  clearUiLocalesList(): void;
  addUiLocales(value: string, index?: number): void;

  getLoginHint(): string;
  setLoginHint(value: string): void;

  getMaxAge(): number;
  setMaxAge(value: number): void;

  getOidc(): AuthRequestOIDC | undefined;
  setOidc(value?: AuthRequestOIDC): void;
  hasOidc(): boolean;
  clearOidc(): void;

  getPreselectedUserId(): string;
  setPreselectedUserId(value: string): void;

  getTypeInfoCase(): AuthSessionCreation.TypeInfoCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthSessionCreation.AsObject;
  static toObject(includeInstance: boolean, msg: AuthSessionCreation): AuthSessionCreation.AsObject;
  static serializeBinaryToWriter(message: AuthSessionCreation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthSessionCreation;
  static deserializeBinaryFromReader(message: AuthSessionCreation, reader: jspb.BinaryReader): AuthSessionCreation;
}

export namespace AuthSessionCreation {
  export type AsObject = {
    agentId: string,
    type: AuthSessionType,
    browserInfo?: BrowserInformation.AsObject,
    clientId: string,
    redirectUri: string,
    state: string,
    prompt: Prompt,
    authContextClassReferenceList: Array<string>,
    uiLocalesList: Array<string>,
    loginHint: string,
    maxAge: number,
    oidc?: AuthRequestOIDC.AsObject,
    preselectedUserId: string,
  }

  export enum TypeInfoCase { 
    TYPE_INFO_NOT_SET = 0,
    OIDC = 12,
  }
}

export class AuthSessionResponse extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): void;

  getId(): string;
  setId(value: string): void;

  getType(): AuthSessionType;
  setType(value: AuthSessionType): void;

  getClientId(): string;
  setClientId(value: string): void;

  getRedirectUri(): string;
  setRedirectUri(value: string): void;

  getState(): string;
  setState(value: string): void;

  getPrompt(): Prompt;
  setPrompt(value: Prompt): void;

  getAuthContextClassReferenceList(): Array<string>;
  setAuthContextClassReferenceList(value: Array<string>): void;
  clearAuthContextClassReferenceList(): void;
  addAuthContextClassReference(value: string, index?: number): void;

  getUiLocalesList(): Array<string>;
  setUiLocalesList(value: Array<string>): void;
  clearUiLocalesList(): void;
  addUiLocales(value: string, index?: number): void;

  getLoginHint(): string;
  setLoginHint(value: string): void;

  getMaxAge(): number;
  setMaxAge(value: number): void;

  getOidc(): AuthRequestOIDC | undefined;
  setOidc(value?: AuthRequestOIDC): void;
  hasOidc(): boolean;
  clearOidc(): void;

  getPossibleStepsList(): Array<NextStep>;
  setPossibleStepsList(value: Array<NextStep>): void;
  clearPossibleStepsList(): void;
  addPossibleSteps(value?: NextStep, index?: number): NextStep;

  getProjectClientIdsList(): Array<string>;
  setProjectClientIdsList(value: Array<string>): void;
  clearProjectClientIdsList(): void;
  addProjectClientIds(value: string, index?: number): void;

  getUserSession(): UserSession | undefined;
  setUserSession(value?: UserSession): void;
  hasUserSession(): boolean;
  clearUserSession(): void;

  getTypeInfoCase(): AuthSessionResponse.TypeInfoCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthSessionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AuthSessionResponse): AuthSessionResponse.AsObject;
  static serializeBinaryToWriter(message: AuthSessionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthSessionResponse;
  static deserializeBinaryFromReader(message: AuthSessionResponse, reader: jspb.BinaryReader): AuthSessionResponse;
}

export namespace AuthSessionResponse {
  export type AsObject = {
    agentId: string,
    id: string,
    type: AuthSessionType,
    clientId: string,
    redirectUri: string,
    state: string,
    prompt: Prompt,
    authContextClassReferenceList: Array<string>,
    uiLocalesList: Array<string>,
    loginHint: string,
    maxAge: number,
    oidc?: AuthRequestOIDC.AsObject,
    possibleStepsList: Array<NextStep.AsObject>,
    projectClientIdsList: Array<string>,
    userSession?: UserSession.AsObject,
  }

  export enum TypeInfoCase { 
    TYPE_INFO_NOT_SET = 0,
    OIDC = 12,
  }
}

export class AuthSessionView extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): void;

  getAuthSessionId(): string;
  setAuthSessionId(value: string): void;

  getType(): AuthSessionType;
  setType(value: AuthSessionType): void;

  getClientId(): string;
  setClientId(value: string): void;

  getUserSessionId(): string;
  setUserSessionId(value: string): void;

  getProjectClientIdsList(): Array<string>;
  setProjectClientIdsList(value: Array<string>): void;
  clearProjectClientIdsList(): void;
  addProjectClientIds(value: string, index?: number): void;

  getTokenId(): string;
  setTokenId(value: string): void;

  getTokenExpiration(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setTokenExpiration(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasTokenExpiration(): boolean;
  clearTokenExpiration(): void;

  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthSessionView.AsObject;
  static toObject(includeInstance: boolean, msg: AuthSessionView): AuthSessionView.AsObject;
  static serializeBinaryToWriter(message: AuthSessionView, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthSessionView;
  static deserializeBinaryFromReader(message: AuthSessionView, reader: jspb.BinaryReader): AuthSessionView;
}

export namespace AuthSessionView {
  export type AsObject = {
    agentId: string,
    authSessionId: string,
    type: AuthSessionType,
    clientId: string,
    userSessionId: string,
    projectClientIdsList: Array<string>,
    tokenId: string,
    tokenExpiration?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    userId: string,
  }
}

export class TokenID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TokenID.AsObject;
  static toObject(includeInstance: boolean, msg: TokenID): TokenID.AsObject;
  static serializeBinaryToWriter(message: TokenID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TokenID;
  static deserializeBinaryFromReader(message: TokenID, reader: jspb.BinaryReader): TokenID;
}

export namespace TokenID {
  export type AsObject = {
    id: string,
  }
}

export class UserSessionID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAgentId(): string;
  setAgentId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSessionID.AsObject;
  static toObject(includeInstance: boolean, msg: UserSessionID): UserSessionID.AsObject;
  static serializeBinaryToWriter(message: UserSessionID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSessionID;
  static deserializeBinaryFromReader(message: UserSessionID, reader: jspb.BinaryReader): UserSessionID;
}

export namespace UserSessionID {
  export type AsObject = {
    id: string,
    agentId: string,
  }
}

export class UserSessions extends jspb.Message {
  getUserSessionsList(): Array<UserSession>;
  setUserSessionsList(value: Array<UserSession>): void;
  clearUserSessionsList(): void;
  addUserSessions(value?: UserSession, index?: number): UserSession;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSessions.AsObject;
  static toObject(includeInstance: boolean, msg: UserSessions): UserSessions.AsObject;
  static serializeBinaryToWriter(message: UserSessions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSessions;
  static deserializeBinaryFromReader(message: UserSessions, reader: jspb.BinaryReader): UserSessions;
}

export namespace UserSessions {
  export type AsObject = {
    userSessionsList: Array<UserSession.AsObject>,
  }
}

export class UserSession extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAgentId(): string;
  setAgentId(value: string): void;

  getAuthState(): UserSessionState;
  setAuthState(value: UserSessionState): void;

  getUser(): AuthUser | undefined;
  setUser(value?: AuthUser): void;
  hasUser(): boolean;
  clearUser(): void;

  getPasswordVerified(): boolean;
  setPasswordVerified(value: boolean): void;

  getMfa(): MfaType;
  setMfa(value: MfaType): void;

  getMfaVerified(): boolean;
  setMfaVerified(value: boolean): void;

  getAuthTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setAuthTime(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasAuthTime(): boolean;
  clearAuthTime(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserSession.AsObject;
  static toObject(includeInstance: boolean, msg: UserSession): UserSession.AsObject;
  static serializeBinaryToWriter(message: UserSession, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserSession;
  static deserializeBinaryFromReader(message: UserSession, reader: jspb.BinaryReader): UserSession;
}

export namespace UserSession {
  export type AsObject = {
    id: string,
    agentId: string,
    authState: UserSessionState,
    user?: AuthUser.AsObject,
    passwordVerified: boolean,
    mfa: MfaType,
    mfaVerified: boolean,
    authTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

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
  }
}

export class AuthUser extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  getUserName(): string;
  setUserName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthUser.AsObject;
  static toObject(includeInstance: boolean, msg: AuthUser): AuthUser.AsObject;
  static serializeBinaryToWriter(message: AuthUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthUser;
  static deserializeBinaryFromReader(message: AuthUser, reader: jspb.BinaryReader): AuthUser;
}

export namespace AuthUser {
  export type AsObject = {
    userId: string,
    userName: string,
  }
}

export class AuthSessionID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAgentId(): string;
  setAgentId(value: string): void;

  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthSessionID.AsObject;
  static toObject(includeInstance: boolean, msg: AuthSessionID): AuthSessionID.AsObject;
  static serializeBinaryToWriter(message: AuthSessionID, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthSessionID;
  static deserializeBinaryFromReader(message: AuthSessionID, reader: jspb.BinaryReader): AuthSessionID;
}

export namespace AuthSessionID {
  export type AsObject = {
    id: string,
    agentId: string,
    browserInfo?: BrowserInformation.AsObject,
  }
}

export class SelectUserRequest extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): void;

  getAuthSessionId(): string;
  setAuthSessionId(value: string): void;

  getUserSessionId(): string;
  setUserSessionId(value: string): void;

  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SelectUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SelectUserRequest): SelectUserRequest.AsObject;
  static serializeBinaryToWriter(message: SelectUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SelectUserRequest;
  static deserializeBinaryFromReader(message: SelectUserRequest, reader: jspb.BinaryReader): SelectUserRequest;
}

export namespace SelectUserRequest {
  export type AsObject = {
    agentId: string,
    authSessionId: string,
    userSessionId: string,
    browserInfo?: BrowserInformation.AsObject,
  }
}

export class VerifyUserRequest extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): void;

  getAuthSessionId(): string;
  setAuthSessionId(value: string): void;

  getUserName(): string;
  setUserName(value: string): void;

  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyUserRequest): VerifyUserRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyUserRequest;
  static deserializeBinaryFromReader(message: VerifyUserRequest, reader: jspb.BinaryReader): VerifyUserRequest;
}

export namespace VerifyUserRequest {
  export type AsObject = {
    agentId: string,
    authSessionId: string,
    userName: string,
    browserInfo?: BrowserInformation.AsObject,
  }
}

export class VerifyPasswordRequest extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): void;

  getAuthSessionId(): string;
  setAuthSessionId(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyPasswordRequest): VerifyPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyPasswordRequest;
  static deserializeBinaryFromReader(message: VerifyPasswordRequest, reader: jspb.BinaryReader): VerifyPasswordRequest;
}

export namespace VerifyPasswordRequest {
  export type AsObject = {
    agentId: string,
    authSessionId: string,
    password: string,
    browserInfo?: BrowserInformation.AsObject,
  }
}

export class VerifyMfaRequest extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): void;

  getAuthSessionId(): string;
  setAuthSessionId(value: string): void;

  getBrowserInfo(): BrowserInformation | undefined;
  setBrowserInfo(value?: BrowserInformation): void;
  hasBrowserInfo(): boolean;
  clearBrowserInfo(): void;

  getOtp(): AuthSessionMultiFactorOTP | undefined;
  setOtp(value?: AuthSessionMultiFactorOTP): void;
  hasOtp(): boolean;
  clearOtp(): void;

  getMfaCase(): VerifyMfaRequest.MfaCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMfaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMfaRequest): VerifyMfaRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyMfaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMfaRequest;
  static deserializeBinaryFromReader(message: VerifyMfaRequest, reader: jspb.BinaryReader): VerifyMfaRequest;
}

export namespace VerifyMfaRequest {
  export type AsObject = {
    agentId: string,
    authSessionId: string,
    browserInfo?: BrowserInformation.AsObject,
    otp?: AuthSessionMultiFactorOTP.AsObject,
  }

  export enum MfaCase { 
    MFA_NOT_SET = 0,
    OTP = 4,
  }
}

export class AuthSessionMultiFactorOTP extends jspb.Message {
  getCode(): string;
  setCode(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthSessionMultiFactorOTP.AsObject;
  static toObject(includeInstance: boolean, msg: AuthSessionMultiFactorOTP): AuthSessionMultiFactorOTP.AsObject;
  static serializeBinaryToWriter(message: AuthSessionMultiFactorOTP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthSessionMultiFactorOTP;
  static deserializeBinaryFromReader(message: AuthSessionMultiFactorOTP, reader: jspb.BinaryReader): AuthSessionMultiFactorOTP;
}

export namespace AuthSessionMultiFactorOTP {
  export type AsObject = {
    code: string,
  }
}

export class NextStep extends jspb.Message {
  getType(): NextStepType;
  setType(value: NextStepType): void;

  getLogin(): LoginData | undefined;
  setLogin(value?: LoginData): void;
  hasLogin(): boolean;
  clearLogin(): void;

  getPassword(): PasswordData | undefined;
  setPassword(value?: PasswordData): void;
  hasPassword(): boolean;
  clearPassword(): void;

  getMfaVerify(): MfaVerifyData | undefined;
  setMfaVerify(value?: MfaVerifyData): void;
  hasMfaVerify(): boolean;
  clearMfaVerify(): void;

  getMfaPrompt(): MfaPromptData | undefined;
  setMfaPrompt(value?: MfaPromptData): void;
  hasMfaPrompt(): boolean;
  clearMfaPrompt(): void;

  getChooseUser(): ChooseUserData | undefined;
  setChooseUser(value?: ChooseUserData): void;
  hasChooseUser(): boolean;
  clearChooseUser(): void;

  getDataCase(): NextStep.DataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NextStep.AsObject;
  static toObject(includeInstance: boolean, msg: NextStep): NextStep.AsObject;
  static serializeBinaryToWriter(message: NextStep, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NextStep;
  static deserializeBinaryFromReader(message: NextStep, reader: jspb.BinaryReader): NextStep;
}

export namespace NextStep {
  export type AsObject = {
    type: NextStepType,
    login?: LoginData.AsObject,
    password?: PasswordData.AsObject,
    mfaVerify?: MfaVerifyData.AsObject,
    mfaPrompt?: MfaPromptData.AsObject,
    chooseUser?: ChooseUserData.AsObject,
  }

  export enum DataCase { 
    DATA_NOT_SET = 0,
    LOGIN = 2,
    PASSWORD = 3,
    MFA_VERIFY = 4,
    MFA_PROMPT = 5,
    CHOOSE_USER = 6,
  }
}

export class LoginData extends jspb.Message {
  getErrMsg(): string;
  setErrMsg(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginData.AsObject;
  static toObject(includeInstance: boolean, msg: LoginData): LoginData.AsObject;
  static serializeBinaryToWriter(message: LoginData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginData;
  static deserializeBinaryFromReader(message: LoginData, reader: jspb.BinaryReader): LoginData;
}

export namespace LoginData {
  export type AsObject = {
    errMsg: string,
  }
}

export class PasswordData extends jspb.Message {
  getErrMsg(): string;
  setErrMsg(value: string): void;

  getFailureCount(): number;
  setFailureCount(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordData.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordData): PasswordData.AsObject;
  static serializeBinaryToWriter(message: PasswordData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordData;
  static deserializeBinaryFromReader(message: PasswordData, reader: jspb.BinaryReader): PasswordData;
}

export namespace PasswordData {
  export type AsObject = {
    errMsg: string,
    failureCount: number,
  }
}

export class MfaVerifyData extends jspb.Message {
  getErrMsg(): string;
  setErrMsg(value: string): void;

  getFailureCount(): number;
  setFailureCount(value: number): void;

  getMfaProvidersList(): Array<MfaType>;
  setMfaProvidersList(value: Array<MfaType>): void;
  clearMfaProvidersList(): void;
  addMfaProviders(value: MfaType, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MfaVerifyData.AsObject;
  static toObject(includeInstance: boolean, msg: MfaVerifyData): MfaVerifyData.AsObject;
  static serializeBinaryToWriter(message: MfaVerifyData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MfaVerifyData;
  static deserializeBinaryFromReader(message: MfaVerifyData, reader: jspb.BinaryReader): MfaVerifyData;
}

export namespace MfaVerifyData {
  export type AsObject = {
    errMsg: string,
    failureCount: number,
    mfaProvidersList: Array<MfaType>,
  }
}

export class MfaPromptData extends jspb.Message {
  getRequired(): boolean;
  setRequired(value: boolean): void;

  getMfaProvidersList(): Array<MfaType>;
  setMfaProvidersList(value: Array<MfaType>): void;
  clearMfaProvidersList(): void;
  addMfaProviders(value: MfaType, index?: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MfaPromptData.AsObject;
  static toObject(includeInstance: boolean, msg: MfaPromptData): MfaPromptData.AsObject;
  static serializeBinaryToWriter(message: MfaPromptData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MfaPromptData;
  static deserializeBinaryFromReader(message: MfaPromptData, reader: jspb.BinaryReader): MfaPromptData;
}

export namespace MfaPromptData {
  export type AsObject = {
    required: boolean,
    mfaProvidersList: Array<MfaType>,
  }
}

export class ChooseUserData extends jspb.Message {
  getUsersList(): Array<ChooseUser>;
  setUsersList(value: Array<ChooseUser>): void;
  clearUsersList(): void;
  addUsers(value?: ChooseUser, index?: number): ChooseUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChooseUserData.AsObject;
  static toObject(includeInstance: boolean, msg: ChooseUserData): ChooseUserData.AsObject;
  static serializeBinaryToWriter(message: ChooseUserData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChooseUserData;
  static deserializeBinaryFromReader(message: ChooseUserData, reader: jspb.BinaryReader): ChooseUserData;
}

export namespace ChooseUserData {
  export type AsObject = {
    usersList: Array<ChooseUser.AsObject>,
  }
}

export class ChooseUser extends jspb.Message {
  getUserSessionId(): string;
  setUserSessionId(value: string): void;

  getUserId(): string;
  setUserId(value: string): void;

  getUserName(): string;
  setUserName(value: string): void;

  getUserSessionState(): UserSessionState;
  setUserSessionState(value: UserSessionState): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChooseUser.AsObject;
  static toObject(includeInstance: boolean, msg: ChooseUser): ChooseUser.AsObject;
  static serializeBinaryToWriter(message: ChooseUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChooseUser;
  static deserializeBinaryFromReader(message: ChooseUser, reader: jspb.BinaryReader): ChooseUser;
}

export namespace ChooseUser {
  export type AsObject = {
    userSessionId: string,
    userId: string,
    userName: string,
    userSessionState: UserSessionState,
  }
}

export class SkipMfaInitRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SkipMfaInitRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SkipMfaInitRequest): SkipMfaInitRequest.AsObject;
  static serializeBinaryToWriter(message: SkipMfaInitRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SkipMfaInitRequest;
  static deserializeBinaryFromReader(message: SkipMfaInitRequest, reader: jspb.BinaryReader): SkipMfaInitRequest;
}

export namespace SkipMfaInitRequest {
  export type AsObject = {
    userId: string,
  }
}

export class BrowserInformation extends jspb.Message {
  getUserAgent(): string;
  setUserAgent(value: string): void;

  getRemoteIp(): IP | undefined;
  setRemoteIp(value?: IP): void;
  hasRemoteIp(): boolean;
  clearRemoteIp(): void;

  getAcceptLanguage(): string;
  setAcceptLanguage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BrowserInformation.AsObject;
  static toObject(includeInstance: boolean, msg: BrowserInformation): BrowserInformation.AsObject;
  static serializeBinaryToWriter(message: BrowserInformation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BrowserInformation;
  static deserializeBinaryFromReader(message: BrowserInformation, reader: jspb.BinaryReader): BrowserInformation;
}

export namespace BrowserInformation {
  export type AsObject = {
    userAgent: string,
    remoteIp?: IP.AsObject,
    acceptLanguage: string,
  }
}

export class IP extends jspb.Message {
  getV4(): string;
  setV4(value: string): void;

  getV6(): string;
  setV6(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IP.AsObject;
  static toObject(includeInstance: boolean, msg: IP): IP.AsObject;
  static serializeBinaryToWriter(message: IP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IP;
  static deserializeBinaryFromReader(message: IP, reader: jspb.BinaryReader): IP;
}

export namespace IP {
  export type AsObject = {
    v4: string,
    v6: string,
  }
}

export class AuthRequestOIDC extends jspb.Message {
  getScopeList(): Array<string>;
  setScopeList(value: Array<string>): void;
  clearScopeList(): void;
  addScope(value: string, index?: number): void;

  getResponseType(): OIDCResponseType;
  setResponseType(value: OIDCResponseType): void;

  getNonce(): string;
  setNonce(value: string): void;

  getCodeChallenge(): CodeChallenge | undefined;
  setCodeChallenge(value?: CodeChallenge): void;
  hasCodeChallenge(): boolean;
  clearCodeChallenge(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthRequestOIDC.AsObject;
  static toObject(includeInstance: boolean, msg: AuthRequestOIDC): AuthRequestOIDC.AsObject;
  static serializeBinaryToWriter(message: AuthRequestOIDC, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthRequestOIDC;
  static deserializeBinaryFromReader(message: AuthRequestOIDC, reader: jspb.BinaryReader): AuthRequestOIDC;
}

export namespace AuthRequestOIDC {
  export type AsObject = {
    scopeList: Array<string>,
    responseType: OIDCResponseType,
    nonce: string,
    codeChallenge?: CodeChallenge.AsObject,
  }
}

export class CodeChallenge extends jspb.Message {
  getChallenge(): string;
  setChallenge(value: string): void;

  getMethod(): CodeChallengeMethod;
  setMethod(value: CodeChallengeMethod): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CodeChallenge.AsObject;
  static toObject(includeInstance: boolean, msg: CodeChallenge): CodeChallenge.AsObject;
  static serializeBinaryToWriter(message: CodeChallenge, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CodeChallenge;
  static deserializeBinaryFromReader(message: CodeChallenge, reader: jspb.BinaryReader): CodeChallenge;
}

export namespace CodeChallenge {
  export type AsObject = {
    challenge: string,
    method: CodeChallengeMethod,
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

export class RegisterUserExternalIDPRequest extends jspb.Message {
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

  getIdpProvider(): IDPProvider | undefined;
  setIdpProvider(value?: IDPProvider): void;
  hasIdpProvider(): boolean;
  clearIdpProvider(): void;

  getOrgId(): string;
  setOrgId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterUserExternalIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterUserExternalIDPRequest): RegisterUserExternalIDPRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterUserExternalIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterUserExternalIDPRequest;
  static deserializeBinaryFromReader(message: RegisterUserExternalIDPRequest, reader: jspb.BinaryReader): RegisterUserExternalIDPRequest;
}

export namespace RegisterUserExternalIDPRequest {
  export type AsObject = {
    email: string,
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: Gender,
    idpProvider?: IDPProvider.AsObject,
    orgId: string,
  }
}

export class IDPProvider extends jspb.Message {
  getProvider(): string;
  setProvider(value: string): void;

  getExternalidpid(): string;
  setExternalidpid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPProvider.AsObject;
  static toObject(includeInstance: boolean, msg: IDPProvider): IDPProvider.AsObject;
  static serializeBinaryToWriter(message: IDPProvider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPProvider;
  static deserializeBinaryFromReader(message: IDPProvider, reader: jspb.BinaryReader): IDPProvider;
}

export namespace IDPProvider {
  export type AsObject = {
    provider: string,
    externalidpid: string,
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

  getPasswordChangeRequired(): boolean;
  setPasswordChangeRequired(value: boolean): void;

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
    isemailverified: boolean,
    phone: string,
    isphoneverified: boolean,
    country: string,
    locality: string,
    postalCode: string,
    region: string,
    streetAddress: string,
    passwordChangeRequired: boolean,
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

  getIsphoneverified(): boolean;
  setIsphoneverified(value: boolean): void;

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
    isphoneverified: boolean,
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

export class ResetPasswordRequest extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): void;

  getType(): NotificationType;
  setType(value: NotificationType): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPasswordRequest): ResetPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: ResetPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordRequest;
  static deserializeBinaryFromReader(message: ResetPasswordRequest, reader: jspb.BinaryReader): ResetPasswordRequest;
}

export namespace ResetPasswordRequest {
  export type AsObject = {
    userName: string,
    type: NotificationType,
  }
}

export class ResetPassword extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getCode(): string;
  setCode(value: string): void;

  getNewPassword(): string;
  setNewPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPassword.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPassword): ResetPassword.AsObject;
  static serializeBinaryToWriter(message: ResetPassword, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPassword;
  static deserializeBinaryFromReader(message: ResetPassword, reader: jspb.BinaryReader): ResetPassword;
}

export namespace ResetPassword {
  export type AsObject = {
    id: string,
    code: string,
    newPassword: string,
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

export class ApplicationID extends jspb.Message {
  getId(): string;
  setId(value: string): void;

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
  }

  export enum AppConfigCase { 
    APP_CONFIG_NOT_SET = 0,
    OIDC_CONFIG = 8,
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

  getClientSecret(): string;
  setClientSecret(value: string): void;

  getClientId(): string;
  setClientId(value: string): void;

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
    clientSecret: string,
    clientId: string,
    authMethodType: OIDCAuthMethodType,
    postLogoutRedirectUrisList: Array<string>,
  }
}

export class ApplicationSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getSortingColumn(): ApplicationSearchKey;
  setSortingColumn(value: ApplicationSearchKey): void;

  getAsc(): boolean;
  setAsc(value: boolean): void;

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
    offset: number,
    limit: number,
    sortingColumn: ApplicationSearchKey,
    asc: boolean,
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

export class ApplicationSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<Application>;
  setResultList(value: Array<Application>): void;
  clearResultList(): void;
  addResult(value?: Application, index?: number): Application;

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
    resultList: Array<Application.AsObject>,
  }
}

export class ApplicationAuthorizeRequest extends jspb.Message {
  getOidcClientAuth(): OIDCClientAuth | undefined;
  setOidcClientAuth(value?: OIDCClientAuth): void;
  hasOidcClientAuth(): boolean;
  clearOidcClientAuth(): void;

  getAuthCase(): ApplicationAuthorizeRequest.AuthCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationAuthorizeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationAuthorizeRequest): ApplicationAuthorizeRequest.AsObject;
  static serializeBinaryToWriter(message: ApplicationAuthorizeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationAuthorizeRequest;
  static deserializeBinaryFromReader(message: ApplicationAuthorizeRequest, reader: jspb.BinaryReader): ApplicationAuthorizeRequest;
}

export namespace ApplicationAuthorizeRequest {
  export type AsObject = {
    oidcClientAuth?: OIDCClientAuth.AsObject,
  }

  export enum AuthCase { 
    AUTH_NOT_SET = 0,
    OIDC_CLIENT_AUTH = 1,
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

export class GrantSearchRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getSortingColumn(): GrantSearchKey;
  setSortingColumn(value: GrantSearchKey): void;

  getAsc(): boolean;
  setAsc(value: boolean): void;

  getQueriesList(): Array<GrantSearchQuery>;
  setQueriesList(value: Array<GrantSearchQuery>): void;
  clearQueriesList(): void;
  addQueries(value?: GrantSearchQuery, index?: number): GrantSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantSearchRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GrantSearchRequest): GrantSearchRequest.AsObject;
  static serializeBinaryToWriter(message: GrantSearchRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantSearchRequest;
  static deserializeBinaryFromReader(message: GrantSearchRequest, reader: jspb.BinaryReader): GrantSearchRequest;
}

export namespace GrantSearchRequest {
  export type AsObject = {
    offset: number,
    limit: number,
    sortingColumn: GrantSearchKey,
    asc: boolean,
    queriesList: Array<GrantSearchQuery.AsObject>,
  }
}

export class GrantSearchQuery extends jspb.Message {
  getKey(): GrantSearchKey;
  setKey(value: GrantSearchKey): void;

  getMethod(): SearchMethod;
  setMethod(value: SearchMethod): void;

  getValue(): string;
  setValue(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: GrantSearchQuery): GrantSearchQuery.AsObject;
  static serializeBinaryToWriter(message: GrantSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantSearchQuery;
  static deserializeBinaryFromReader(message: GrantSearchQuery, reader: jspb.BinaryReader): GrantSearchQuery;
}

export namespace GrantSearchQuery {
  export type AsObject = {
    key: GrantSearchKey,
    method: SearchMethod,
    value: string,
  }
}

export class GrantSearchResponse extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): void;

  getLimit(): number;
  setLimit(value: number): void;

  getTotalResult(): number;
  setTotalResult(value: number): void;

  getResultList(): Array<Grant>;
  setResultList(value: Array<Grant>): void;
  clearResultList(): void;
  addResult(value?: Grant, index?: number): Grant;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantSearchResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GrantSearchResponse): GrantSearchResponse.AsObject;
  static serializeBinaryToWriter(message: GrantSearchResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantSearchResponse;
  static deserializeBinaryFromReader(message: GrantSearchResponse, reader: jspb.BinaryReader): GrantSearchResponse;
}

export namespace GrantSearchResponse {
  export type AsObject = {
    offset: number,
    limit: number,
    totalResult: number,
    resultList: Array<Grant.AsObject>,
  }
}

export class Grant extends jspb.Message {
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
  toObject(includeInstance?: boolean): Grant.AsObject;
  static toObject(includeInstance: boolean, msg: Grant): Grant.AsObject;
  static serializeBinaryToWriter(message: Grant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Grant;
  static deserializeBinaryFromReader(message: Grant, reader: jspb.BinaryReader): Grant;
}

export namespace Grant {
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

export class IsAdminResponse extends jspb.Message {
  getIsAdmin(): boolean;
  setIsAdmin(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsAdminResponse.AsObject;
  static toObject(includeInstance: boolean, msg: IsAdminResponse): IsAdminResponse.AsObject;
  static serializeBinaryToWriter(message: IsAdminResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsAdminResponse;
  static deserializeBinaryFromReader(message: IsAdminResponse, reader: jspb.BinaryReader): IsAdminResponse;
}

export namespace IsAdminResponse {
  export type AsObject = {
    isAdmin: boolean,
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

export class CreateTokenRequest extends jspb.Message {
  getAgentId(): string;
  setAgentId(value: string): void;

  getAuthSessionId(): string;
  setAuthSessionId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateTokenRequest): CreateTokenRequest.AsObject;
  static serializeBinaryToWriter(message: CreateTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateTokenRequest;
  static deserializeBinaryFromReader(message: CreateTokenRequest, reader: jspb.BinaryReader): CreateTokenRequest;
}

export namespace CreateTokenRequest {
  export type AsObject = {
    agentId: string,
    authSessionId: string,
  }
}

export class Token extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getExpiration(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpiration(value?: google_protobuf_timestamp_pb.Timestamp): void;
  hasExpiration(): boolean;
  clearExpiration(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Token.AsObject;
  static toObject(includeInstance: boolean, msg: Token): Token.AsObject;
  static serializeBinaryToWriter(message: Token, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Token;
  static deserializeBinaryFromReader(message: Token, reader: jspb.BinaryReader): Token;
}

export namespace Token {
  export type AsObject = {
    id: string,
    expiration?: google_protobuf_timestamp_pb.Timestamp.AsObject,
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

export class VerifyUserInitRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getCode(): string;
  setCode(value: string): void;

  getPassword(): string;
  setPassword(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyUserInitRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyUserInitRequest): VerifyUserInitRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyUserInitRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyUserInitRequest;
  static deserializeBinaryFromReader(message: VerifyUserInitRequest, reader: jspb.BinaryReader): VerifyUserInitRequest;
}

export namespace VerifyUserInitRequest {
  export type AsObject = {
    id: string,
    code: string,
    password: string,
  }
}

export enum UserAgentState { 
  NO_STATE = 0,
  ACTIVE_SESSION = 1,
  TERMINATED_SESSION = 2,
}
export enum UserSessionState { 
  USER_SESSION_STATE_UNKNOWN = 0,
  USER_SESSION_STATE_ACTIVE = 1,
  USER_SESSION_STATE_TERMINATED = 2,
}
export enum NextStepType { 
  NEXT_STEP_UNSPECIFIED = 0,
  NEXT_STEP_LOGIN = 1,
  NEXT_STEP_PASSWORD = 2,
  NEXT_STEP_CHANGE_PASSWORD = 3,
  NEXT_STEP_MFA_PROMPT = 4,
  NEXT_STEP_MFA_INIT_CHOICE = 5,
  NEXT_STEP_MFA_INIT_CREATE = 6,
  NEXT_STEP_MFA_INIT_VERIFY = 7,
  NEXT_STEP_MFA_INIT_DONE = 8,
  NEXT_STEP_MFA_VERIFY = 9,
  NEXT_STEP_MFA_VERIFY_ASYNC = 10,
  NEXT_STEP_VERIFY_EMAIL = 11,
  NEXT_STEP_REDIRECT_TO_CALLBACK = 12,
  NEXT_STEP_INIT_PASSWORD = 13,
  NEXT_STEP_CHOOSE_USER = 14,
}
export enum AuthSessionType { 
  TYPE_UNKNOWN = 0,
  TYPE_OIDC = 1,
  TYPE_SAML = 2,
}
export enum Prompt { 
  NO_PROMPT = 0,
  PROMPT_NONE = 1,
  PROMPT_LOGIN = 2,
  PROMPT_CONSENT = 3,
  PROMPT_SELECT_ACCOUNT = 4,
}
export enum OIDCResponseType { 
  CODE = 0,
  ID_TOKEN = 1,
  ID_TOKEN_TOKEN = 2,
}
export enum CodeChallengeMethod { 
  PLAIN = 0,
  S256 = 1,
}
export enum UserState { 
  NONE = 0,
  ACTIVE = 1,
  INACTIVE = 2,
  DELETED = 3,
  LOCKED = 4,
  SUSPEND = 5,
  INITIAL = 6,
}
export enum Gender { 
  UNKNOWN_GENDER = 0,
  FEMALE = 1,
  MALE = 2,
  DIVERSE = 3,
}
export enum NotificationType { 
  EMAIL = 0,
  SMS = 1,
}
export enum MfaType { 
  NO_MFA = 0,
  MFA_SMS = 1,
  MFA_OTP = 2,
}
export enum MFAState { 
  MFASTATE_NO = 0,
  NOT_READY = 1,
  READY = 2,
  REMOVED = 3,
}
export enum AppState { 
  NONE_APP = 0,
  ACTIVE_APP = 1,
  INACTIVE_APP = 2,
  DELETED_APP = 3,
}
export enum OIDCGrantType { 
  AUTHORIZATION_CODE = 0,
  GRANT_TYPE_NONE = 1,
  REFRESH_TOKEN = 2,
}
export enum OIDCApplicationType { 
  WEB = 0,
  USER_AGENT = 1,
  NATIVE = 2,
}
export enum OIDCAuthMethodType { 
  AUTH_TYPE_BASIC = 0,
  AUTH_TYPE_POST = 1,
  AUTH_TYPE_NONE = 2,
}
export enum ApplicationSearchKey { 
  UNKNOWN = 0,
  APP_TYPE = 1,
  STATE = 2,
  CLIENT_ID = 3,
  APP_NAME = 4,
  PROJECT_ID = 5,
}
export enum SearchMethod { 
  EQUALS = 0,
  STARTS_WITH = 1,
  CONTAINS = 2,
}
export enum GrantSearchKey { 
  GRANTSEARCHKEY_UNKNOWN = 0,
  GRANTSEARCHKEY_ORG_ID = 1,
  GRANTSEARCHKEY_PROJECT_ID = 2,
  GRANTSEARCHKEY_USER_ID = 3,
}
export enum MyProjectOrgSearchKey { 
  MYPROJECTORGKEY_UNKNOWN = 0,
  MYPROJECTORGKEY_ORG_NAME = 1,
}

/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";
import { Timestamp } from "../google/protobuf/timestamp";
import { Change, ChangeQuery } from "./change";
import { IDPUserLink } from "./idp";
import { Metadata, MetadataQuery } from "./metadata";
import { ListDetails, ListQuery, ObjectDetails } from "./object";
import { Org, OrgQuery } from "./org";
import { LabelPolicy, LoginPolicy, PasswordComplexityPolicy, PrivacyPolicy } from "./policy";
import {
  AuthFactor,
  Email,
  Gender,
  genderFromJSON,
  genderToJSON,
  Membership,
  MembershipQuery,
  Phone,
  Profile,
  RefreshToken,
  Session,
  User,
  WebAuthNKey,
  WebAuthNToken,
  WebAuthNVerification,
} from "./user";

export const protobufPackage = "zitadel.auth.v1";

/** This is an empty request */
export interface HealthzRequest {
}

/** This is an empty response */
export interface HealthzResponse {
}

/** This is an empty request */
export interface GetSupportedLanguagesRequest {
}

/** This is an empty response */
export interface GetSupportedLanguagesResponse {
  languages: string[];
}

/**
 * This is an empty request
 * the request parameters are read from the token-header
 */
export interface GetMyUserRequest {
}

export interface GetMyUserResponse {
  user: User | undefined;
  lastLogin: Date | undefined;
}

/**
 * This is an empty request
 * the request parameters are read from the token-header
 */
export interface RemoveMyUserRequest {
}

export interface RemoveMyUserResponse {
  details: ObjectDetails | undefined;
}

export interface ListMyUserChangesRequest {
  query: ChangeQuery | undefined;
}

export interface ListMyUserChangesResponse {
  /** zitadel.v1.ListDetails details = 1; was always returned empty (as we cannot get the necessary info) */
  result: Change[];
}

/** This is an empty request */
export interface ListMyUserSessionsRequest {
}

export interface ListMyUserSessionsResponse {
  result: Session[];
}

export interface ListMyMetadataRequest {
  query: ListQuery | undefined;
  queries: MetadataQuery[];
}

export interface ListMyMetadataResponse {
  details: ListDetails | undefined;
  result: Metadata[];
}

export interface GetMyMetadataRequest {
  key: string;
}

export interface GetMyMetadataResponse {
  metadata: Metadata | undefined;
}

export interface SetMyMetadataRequest {
  key: string;
  value: Buffer;
}

export interface SetMyMetadataResponse {
  details: ObjectDetails | undefined;
}

export interface BulkSetMyMetadataRequest {
  metadata: BulkSetMyMetadataRequest_Metadata[];
}

export interface BulkSetMyMetadataRequest_Metadata {
  key: string;
  value: Buffer;
}

export interface BulkSetMyMetadataResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveMyMetadataRequest {
  key: string;
}

export interface RemoveMyMetadataResponse {
  details: ObjectDetails | undefined;
}

export interface BulkRemoveMyMetadataRequest {
  keys: string[];
}

export interface BulkRemoveMyMetadataResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ListMyRefreshTokensRequest {
}

export interface ListMyRefreshTokensResponse {
  details: ListDetails | undefined;
  result: RefreshToken[];
}

export interface RevokeMyRefreshTokenRequest {
  id: string;
}

export interface RevokeMyRefreshTokenResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RevokeAllMyRefreshTokensRequest {
}

/** This is an empty response */
export interface RevokeAllMyRefreshTokensResponse {
}

export interface UpdateMyUserNameRequest {
  userName: string;
}

export interface UpdateMyUserNameResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetMyPasswordComplexityPolicyRequest {
}

export interface GetMyPasswordComplexityPolicyResponse {
  policy: PasswordComplexityPolicy | undefined;
}

export interface UpdateMyPasswordRequest {
  oldPassword: string;
  newPassword: string;
}

export interface UpdateMyPasswordResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetMyProfileRequest {
}

export interface GetMyProfileResponse {
  details: ObjectDetails | undefined;
  profile: Profile | undefined;
}

export interface UpdateMyProfileRequest {
  firstName: string;
  lastName: string;
  nickName: string;
  displayName: string;
  preferredLanguage: string;
  gender: Gender;
}

export interface UpdateMyProfileResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetMyEmailRequest {
}

export interface GetMyEmailResponse {
  details: ObjectDetails | undefined;
  email: Email | undefined;
}

export interface SetMyEmailRequest {
  email: string;
}

export interface SetMyEmailResponse {
  details: ObjectDetails | undefined;
}

export interface VerifyMyEmailRequest {
  code: string;
}

export interface VerifyMyEmailResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ResendMyEmailVerificationRequest {
}

export interface ResendMyEmailVerificationResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface GetMyPhoneRequest {
}

export interface GetMyPhoneResponse {
  details: ObjectDetails | undefined;
  phone: Phone | undefined;
}

export interface SetMyPhoneRequest {
  phone: string;
}

export interface SetMyPhoneResponse {
  details: ObjectDetails | undefined;
}

export interface VerifyMyPhoneRequest {
  code: string;
}

export interface VerifyMyPhoneResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ResendMyPhoneVerificationRequest {
}

export interface ResendMyPhoneVerificationResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RemoveMyPhoneRequest {
}

export interface RemoveMyPhoneResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RemoveMyAvatarRequest {
}

export interface RemoveMyAvatarResponse {
  details: ObjectDetails | undefined;
}

export interface ListMyLinkedIDPsRequest {
  /** list limitations and ordering */
  query: ListQuery | undefined;
}

export interface ListMyLinkedIDPsResponse {
  details: ListDetails | undefined;
  result: IDPUserLink[];
}

export interface RemoveMyLinkedIDPRequest {
  idpId: string;
  linkedUserId: string;
}

export interface RemoveMyLinkedIDPResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ListMyAuthFactorsRequest {
}

export interface ListMyAuthFactorsResponse {
  result: AuthFactor[];
}

/** This is an empty request */
export interface AddMyAuthFactorU2FRequest {
}

export interface AddMyAuthFactorU2FResponse {
  key: WebAuthNKey | undefined;
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface AddMyAuthFactorOTPRequest {
}

export interface AddMyAuthFactorOTPResponse {
  url: string;
  secret: string;
  details: ObjectDetails | undefined;
}

export interface VerifyMyAuthFactorOTPRequest {
  code: string;
}

export interface VerifyMyAuthFactorOTPResponse {
  details: ObjectDetails | undefined;
}

export interface VerifyMyAuthFactorU2FRequest {
  verification: WebAuthNVerification | undefined;
}

export interface VerifyMyAuthFactorU2FResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface RemoveMyAuthFactorOTPRequest {
}

export interface RemoveMyAuthFactorOTPResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveMyAuthFactorU2FRequest {
  tokenId: string;
}

export interface RemoveMyAuthFactorU2FResponse {
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface ListMyPasswordlessRequest {
}

export interface ListMyPasswordlessResponse {
  result: WebAuthNToken[];
}

/** This is an empty request */
export interface AddMyPasswordlessRequest {
}

export interface AddMyPasswordlessResponse {
  key: WebAuthNKey | undefined;
  details: ObjectDetails | undefined;
}

/** This is an empty request */
export interface AddMyPasswordlessLinkRequest {
}

export interface AddMyPasswordlessLinkResponse {
  details: ObjectDetails | undefined;
  link: string;
  expiration: Duration | undefined;
}

/** This is an empty request */
export interface SendMyPasswordlessLinkRequest {
}

export interface SendMyPasswordlessLinkResponse {
  details: ObjectDetails | undefined;
}

export interface VerifyMyPasswordlessRequest {
  verification: WebAuthNVerification | undefined;
}

export interface VerifyMyPasswordlessResponse {
  details: ObjectDetails | undefined;
}

export interface RemoveMyPasswordlessRequest {
  tokenId: string;
}

export interface RemoveMyPasswordlessResponse {
  details: ObjectDetails | undefined;
}

export interface ListMyUserGrantsRequest {
  /** list limitations and ordering */
  query: ListQuery | undefined;
}

export interface ListMyUserGrantsResponse {
  details: ListDetails | undefined;
  result: UserGrant[];
}

export interface UserGrant {
  orgId: string;
  projectId: string;
  userId: string;
  /** Deprecated: user role_keys */
  roles: string[];
  orgName: string;
  grantId: string;
  details: ObjectDetails | undefined;
  orgDomain: string;
  projectName: string;
  projectGrantId: string;
  roleKeys: string[];
}

export interface ListMyProjectOrgsRequest {
  /** list limitations and ordering */
  query:
    | ListQuery
    | undefined;
  /** criteria the client is looking for */
  queries: OrgQuery[];
}

export interface ListMyProjectOrgsResponse {
  details: ListDetails | undefined;
  result: Org[];
}

/** This is an empty request */
export interface ListMyZitadelPermissionsRequest {
}

export interface ListMyZitadelPermissionsResponse {
  result: string[];
}

/** This is an empty request */
export interface ListMyProjectPermissionsRequest {
}

export interface ListMyProjectPermissionsResponse {
  result: string[];
}

export interface ListMyMembershipsRequest {
  /** the field the result is sorted */
  query:
    | ListQuery
    | undefined;
  /** criteria the client is looking for */
  queries: MembershipQuery[];
}

export interface ListMyMembershipsResponse {
  details: ListDetails | undefined;
  result: Membership[];
}

/** This is an empty request */
export interface GetMyLabelPolicyRequest {
}

export interface GetMyLabelPolicyResponse {
  policy: LabelPolicy | undefined;
}

/** This is an empty request */
export interface GetMyPrivacyPolicyRequest {
}

export interface GetMyPrivacyPolicyResponse {
  policy: PrivacyPolicy | undefined;
}

/** This is an empty request */
export interface GetMyLoginPolicyRequest {
}

export interface GetMyLoginPolicyResponse {
  policy: LoginPolicy | undefined;
}

function createBaseHealthzRequest(): HealthzRequest {
  return {};
}

export const HealthzRequest = {
  encode(_: HealthzRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HealthzRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHealthzRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): HealthzRequest {
    return {};
  },

  toJSON(_: HealthzRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<HealthzRequest>): HealthzRequest {
    return HealthzRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<HealthzRequest>): HealthzRequest {
    const message = createBaseHealthzRequest();
    return message;
  },
};

function createBaseHealthzResponse(): HealthzResponse {
  return {};
}

export const HealthzResponse = {
  encode(_: HealthzResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HealthzResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHealthzResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): HealthzResponse {
    return {};
  },

  toJSON(_: HealthzResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<HealthzResponse>): HealthzResponse {
    return HealthzResponse.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<HealthzResponse>): HealthzResponse {
    const message = createBaseHealthzResponse();
    return message;
  },
};

function createBaseGetSupportedLanguagesRequest(): GetSupportedLanguagesRequest {
  return {};
}

export const GetSupportedLanguagesRequest = {
  encode(_: GetSupportedLanguagesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSupportedLanguagesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSupportedLanguagesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetSupportedLanguagesRequest {
    return {};
  },

  toJSON(_: GetSupportedLanguagesRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetSupportedLanguagesRequest>): GetSupportedLanguagesRequest {
    return GetSupportedLanguagesRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetSupportedLanguagesRequest>): GetSupportedLanguagesRequest {
    const message = createBaseGetSupportedLanguagesRequest();
    return message;
  },
};

function createBaseGetSupportedLanguagesResponse(): GetSupportedLanguagesResponse {
  return { languages: [] };
}

export const GetSupportedLanguagesResponse = {
  encode(message: GetSupportedLanguagesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.languages) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSupportedLanguagesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSupportedLanguagesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.languages.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetSupportedLanguagesResponse {
    return { languages: Array.isArray(object?.languages) ? object.languages.map((e: any) => String(e)) : [] };
  },

  toJSON(message: GetSupportedLanguagesResponse): unknown {
    const obj: any = {};
    if (message.languages) {
      obj.languages = message.languages.map((e) => e);
    } else {
      obj.languages = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GetSupportedLanguagesResponse>): GetSupportedLanguagesResponse {
    return GetSupportedLanguagesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSupportedLanguagesResponse>): GetSupportedLanguagesResponse {
    const message = createBaseGetSupportedLanguagesResponse();
    message.languages = object.languages?.map((e) => e) || [];
    return message;
  },
};

function createBaseGetMyUserRequest(): GetMyUserRequest {
  return {};
}

export const GetMyUserRequest = {
  encode(_: GetMyUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyUserRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyUserRequest {
    return {};
  },

  toJSON(_: GetMyUserRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyUserRequest>): GetMyUserRequest {
    return GetMyUserRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyUserRequest>): GetMyUserRequest {
    const message = createBaseGetMyUserRequest();
    return message;
  },
};

function createBaseGetMyUserResponse(): GetMyUserResponse {
  return { user: undefined, lastLogin: undefined };
}

export const GetMyUserResponse = {
  encode(message: GetMyUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    if (message.lastLogin !== undefined) {
      Timestamp.encode(toTimestamp(message.lastLogin), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyUserResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.user = User.decode(reader, reader.uint32());
          break;
        case 2:
          message.lastLogin = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyUserResponse {
    return {
      user: isSet(object.user) ? User.fromJSON(object.user) : undefined,
      lastLogin: isSet(object.lastLogin) ? fromJsonTimestamp(object.lastLogin) : undefined,
    };
  },

  toJSON(message: GetMyUserResponse): unknown {
    const obj: any = {};
    message.user !== undefined && (obj.user = message.user ? User.toJSON(message.user) : undefined);
    message.lastLogin !== undefined && (obj.lastLogin = message.lastLogin.toISOString());
    return obj;
  },

  create(base?: DeepPartial<GetMyUserResponse>): GetMyUserResponse {
    return GetMyUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyUserResponse>): GetMyUserResponse {
    const message = createBaseGetMyUserResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    message.lastLogin = object.lastLogin ?? undefined;
    return message;
  },
};

function createBaseRemoveMyUserRequest(): RemoveMyUserRequest {
  return {};
}

export const RemoveMyUserRequest = {
  encode(_: RemoveMyUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyUserRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveMyUserRequest {
    return {};
  },

  toJSON(_: RemoveMyUserRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveMyUserRequest>): RemoveMyUserRequest {
    return RemoveMyUserRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveMyUserRequest>): RemoveMyUserRequest {
    const message = createBaseRemoveMyUserRequest();
    return message;
  },
};

function createBaseRemoveMyUserResponse(): RemoveMyUserResponse {
  return { details: undefined };
}

export const RemoveMyUserResponse = {
  encode(message: RemoveMyUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyUserResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyUserResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMyUserResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyUserResponse>): RemoveMyUserResponse {
    return RemoveMyUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyUserResponse>): RemoveMyUserResponse {
    const message = createBaseRemoveMyUserResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListMyUserChangesRequest(): ListMyUserChangesRequest {
  return { query: undefined };
}

export const ListMyUserChangesRequest = {
  encode(message: ListMyUserChangesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ChangeQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyUserChangesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyUserChangesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ChangeQuery.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyUserChangesRequest {
    return { query: isSet(object.query) ? ChangeQuery.fromJSON(object.query) : undefined };
  },

  toJSON(message: ListMyUserChangesRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ChangeQuery.toJSON(message.query) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ListMyUserChangesRequest>): ListMyUserChangesRequest {
    return ListMyUserChangesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyUserChangesRequest>): ListMyUserChangesRequest {
    const message = createBaseListMyUserChangesRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ChangeQuery.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseListMyUserChangesResponse(): ListMyUserChangesResponse {
  return { result: [] };
}

export const ListMyUserChangesResponse = {
  encode(message: ListMyUserChangesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.result) {
      Change.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyUserChangesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyUserChangesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.result.push(Change.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyUserChangesResponse {
    return { result: Array.isArray(object?.result) ? object.result.map((e: any) => Change.fromJSON(e)) : [] };
  },

  toJSON(message: ListMyUserChangesResponse): unknown {
    const obj: any = {};
    if (message.result) {
      obj.result = message.result.map((e) => e ? Change.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyUserChangesResponse>): ListMyUserChangesResponse {
    return ListMyUserChangesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyUserChangesResponse>): ListMyUserChangesResponse {
    const message = createBaseListMyUserChangesResponse();
    message.result = object.result?.map((e) => Change.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListMyUserSessionsRequest(): ListMyUserSessionsRequest {
  return {};
}

export const ListMyUserSessionsRequest = {
  encode(_: ListMyUserSessionsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyUserSessionsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyUserSessionsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListMyUserSessionsRequest {
    return {};
  },

  toJSON(_: ListMyUserSessionsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListMyUserSessionsRequest>): ListMyUserSessionsRequest {
    return ListMyUserSessionsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListMyUserSessionsRequest>): ListMyUserSessionsRequest {
    const message = createBaseListMyUserSessionsRequest();
    return message;
  },
};

function createBaseListMyUserSessionsResponse(): ListMyUserSessionsResponse {
  return { result: [] };
}

export const ListMyUserSessionsResponse = {
  encode(message: ListMyUserSessionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.result) {
      Session.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyUserSessionsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyUserSessionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.result.push(Session.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyUserSessionsResponse {
    return { result: Array.isArray(object?.result) ? object.result.map((e: any) => Session.fromJSON(e)) : [] };
  },

  toJSON(message: ListMyUserSessionsResponse): unknown {
    const obj: any = {};
    if (message.result) {
      obj.result = message.result.map((e) => e ? Session.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyUserSessionsResponse>): ListMyUserSessionsResponse {
    return ListMyUserSessionsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyUserSessionsResponse>): ListMyUserSessionsResponse {
    const message = createBaseListMyUserSessionsResponse();
    message.result = object.result?.map((e) => Session.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListMyMetadataRequest(): ListMyMetadataRequest {
  return { query: undefined, queries: [] };
}

export const ListMyMetadataRequest = {
  encode(message: ListMyMetadataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.queries) {
      MetadataQuery.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyMetadataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyMetadataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.queries.push(MetadataQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyMetadataRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => MetadataQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyMetadataRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? MetadataQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyMetadataRequest>): ListMyMetadataRequest {
    return ListMyMetadataRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyMetadataRequest>): ListMyMetadataRequest {
    const message = createBaseListMyMetadataRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.queries = object.queries?.map((e) => MetadataQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListMyMetadataResponse(): ListMyMetadataResponse {
  return { details: undefined, result: [] };
}

export const ListMyMetadataResponse = {
  encode(message: ListMyMetadataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      Metadata.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyMetadataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyMetadataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(Metadata.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyMetadataResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Metadata.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyMetadataResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? Metadata.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyMetadataResponse>): ListMyMetadataResponse {
    return ListMyMetadataResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyMetadataResponse>): ListMyMetadataResponse {
    const message = createBaseListMyMetadataResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => Metadata.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetMyMetadataRequest(): GetMyMetadataRequest {
  return { key: "" };
}

export const GetMyMetadataRequest = {
  encode(message: GetMyMetadataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyMetadataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyMetadataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyMetadataRequest {
    return { key: isSet(object.key) ? String(object.key) : "" };
  },

  toJSON(message: GetMyMetadataRequest): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    return obj;
  },

  create(base?: DeepPartial<GetMyMetadataRequest>): GetMyMetadataRequest {
    return GetMyMetadataRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyMetadataRequest>): GetMyMetadataRequest {
    const message = createBaseGetMyMetadataRequest();
    message.key = object.key ?? "";
    return message;
  },
};

function createBaseGetMyMetadataResponse(): GetMyMetadataResponse {
  return { metadata: undefined };
}

export const GetMyMetadataResponse = {
  encode(message: GetMyMetadataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.metadata !== undefined) {
      Metadata.encode(message.metadata, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyMetadataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyMetadataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.metadata = Metadata.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyMetadataResponse {
    return { metadata: isSet(object.metadata) ? Metadata.fromJSON(object.metadata) : undefined };
  },

  toJSON(message: GetMyMetadataResponse): unknown {
    const obj: any = {};
    message.metadata !== undefined && (obj.metadata = message.metadata ? Metadata.toJSON(message.metadata) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyMetadataResponse>): GetMyMetadataResponse {
    return GetMyMetadataResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyMetadataResponse>): GetMyMetadataResponse {
    const message = createBaseGetMyMetadataResponse();
    message.metadata = (object.metadata !== undefined && object.metadata !== null)
      ? Metadata.fromPartial(object.metadata)
      : undefined;
    return message;
  },
};

function createBaseSetMyMetadataRequest(): SetMyMetadataRequest {
  return { key: "", value: Buffer.alloc(0) };
}

export const SetMyMetadataRequest = {
  encode(message: SetMyMetadataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value.length !== 0) {
      writer.uint32(18).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetMyMetadataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetMyMetadataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.value = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetMyMetadataRequest {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Buffer.from(bytesFromBase64(object.value)) : Buffer.alloc(0),
    };
  },

  toJSON(message: SetMyMetadataRequest): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = base64FromBytes(message.value !== undefined ? message.value : Buffer.alloc(0)));
    return obj;
  },

  create(base?: DeepPartial<SetMyMetadataRequest>): SetMyMetadataRequest {
    return SetMyMetadataRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetMyMetadataRequest>): SetMyMetadataRequest {
    const message = createBaseSetMyMetadataRequest();
    message.key = object.key ?? "";
    message.value = object.value ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseSetMyMetadataResponse(): SetMyMetadataResponse {
  return { details: undefined };
}

export const SetMyMetadataResponse = {
  encode(message: SetMyMetadataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetMyMetadataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetMyMetadataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetMyMetadataResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetMyMetadataResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetMyMetadataResponse>): SetMyMetadataResponse {
    return SetMyMetadataResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetMyMetadataResponse>): SetMyMetadataResponse {
    const message = createBaseSetMyMetadataResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseBulkSetMyMetadataRequest(): BulkSetMyMetadataRequest {
  return { metadata: [] };
}

export const BulkSetMyMetadataRequest = {
  encode(message: BulkSetMyMetadataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.metadata) {
      BulkSetMyMetadataRequest_Metadata.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BulkSetMyMetadataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBulkSetMyMetadataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.metadata.push(BulkSetMyMetadataRequest_Metadata.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BulkSetMyMetadataRequest {
    return {
      metadata: Array.isArray(object?.metadata)
        ? object.metadata.map((e: any) => BulkSetMyMetadataRequest_Metadata.fromJSON(e))
        : [],
    };
  },

  toJSON(message: BulkSetMyMetadataRequest): unknown {
    const obj: any = {};
    if (message.metadata) {
      obj.metadata = message.metadata.map((e) => e ? BulkSetMyMetadataRequest_Metadata.toJSON(e) : undefined);
    } else {
      obj.metadata = [];
    }
    return obj;
  },

  create(base?: DeepPartial<BulkSetMyMetadataRequest>): BulkSetMyMetadataRequest {
    return BulkSetMyMetadataRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<BulkSetMyMetadataRequest>): BulkSetMyMetadataRequest {
    const message = createBaseBulkSetMyMetadataRequest();
    message.metadata = object.metadata?.map((e) => BulkSetMyMetadataRequest_Metadata.fromPartial(e)) || [];
    return message;
  },
};

function createBaseBulkSetMyMetadataRequest_Metadata(): BulkSetMyMetadataRequest_Metadata {
  return { key: "", value: Buffer.alloc(0) };
}

export const BulkSetMyMetadataRequest_Metadata = {
  encode(message: BulkSetMyMetadataRequest_Metadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value.length !== 0) {
      writer.uint32(18).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BulkSetMyMetadataRequest_Metadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBulkSetMyMetadataRequest_Metadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        case 2:
          message.value = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BulkSetMyMetadataRequest_Metadata {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Buffer.from(bytesFromBase64(object.value)) : Buffer.alloc(0),
    };
  },

  toJSON(message: BulkSetMyMetadataRequest_Metadata): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = base64FromBytes(message.value !== undefined ? message.value : Buffer.alloc(0)));
    return obj;
  },

  create(base?: DeepPartial<BulkSetMyMetadataRequest_Metadata>): BulkSetMyMetadataRequest_Metadata {
    return BulkSetMyMetadataRequest_Metadata.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<BulkSetMyMetadataRequest_Metadata>): BulkSetMyMetadataRequest_Metadata {
    const message = createBaseBulkSetMyMetadataRequest_Metadata();
    message.key = object.key ?? "";
    message.value = object.value ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseBulkSetMyMetadataResponse(): BulkSetMyMetadataResponse {
  return { details: undefined };
}

export const BulkSetMyMetadataResponse = {
  encode(message: BulkSetMyMetadataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BulkSetMyMetadataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBulkSetMyMetadataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BulkSetMyMetadataResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: BulkSetMyMetadataResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<BulkSetMyMetadataResponse>): BulkSetMyMetadataResponse {
    return BulkSetMyMetadataResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<BulkSetMyMetadataResponse>): BulkSetMyMetadataResponse {
    const message = createBaseBulkSetMyMetadataResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveMyMetadataRequest(): RemoveMyMetadataRequest {
  return { key: "" };
}

export const RemoveMyMetadataRequest = {
  encode(message: RemoveMyMetadataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyMetadataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyMetadataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyMetadataRequest {
    return { key: isSet(object.key) ? String(object.key) : "" };
  },

  toJSON(message: RemoveMyMetadataRequest): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyMetadataRequest>): RemoveMyMetadataRequest {
    return RemoveMyMetadataRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyMetadataRequest>): RemoveMyMetadataRequest {
    const message = createBaseRemoveMyMetadataRequest();
    message.key = object.key ?? "";
    return message;
  },
};

function createBaseRemoveMyMetadataResponse(): RemoveMyMetadataResponse {
  return { details: undefined };
}

export const RemoveMyMetadataResponse = {
  encode(message: RemoveMyMetadataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyMetadataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyMetadataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyMetadataResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMyMetadataResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyMetadataResponse>): RemoveMyMetadataResponse {
    return RemoveMyMetadataResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyMetadataResponse>): RemoveMyMetadataResponse {
    const message = createBaseRemoveMyMetadataResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseBulkRemoveMyMetadataRequest(): BulkRemoveMyMetadataRequest {
  return { keys: [] };
}

export const BulkRemoveMyMetadataRequest = {
  encode(message: BulkRemoveMyMetadataRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.keys) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BulkRemoveMyMetadataRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBulkRemoveMyMetadataRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.keys.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BulkRemoveMyMetadataRequest {
    return { keys: Array.isArray(object?.keys) ? object.keys.map((e: any) => String(e)) : [] };
  },

  toJSON(message: BulkRemoveMyMetadataRequest): unknown {
    const obj: any = {};
    if (message.keys) {
      obj.keys = message.keys.map((e) => e);
    } else {
      obj.keys = [];
    }
    return obj;
  },

  create(base?: DeepPartial<BulkRemoveMyMetadataRequest>): BulkRemoveMyMetadataRequest {
    return BulkRemoveMyMetadataRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<BulkRemoveMyMetadataRequest>): BulkRemoveMyMetadataRequest {
    const message = createBaseBulkRemoveMyMetadataRequest();
    message.keys = object.keys?.map((e) => e) || [];
    return message;
  },
};

function createBaseBulkRemoveMyMetadataResponse(): BulkRemoveMyMetadataResponse {
  return { details: undefined };
}

export const BulkRemoveMyMetadataResponse = {
  encode(message: BulkRemoveMyMetadataResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BulkRemoveMyMetadataResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBulkRemoveMyMetadataResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): BulkRemoveMyMetadataResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: BulkRemoveMyMetadataResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<BulkRemoveMyMetadataResponse>): BulkRemoveMyMetadataResponse {
    return BulkRemoveMyMetadataResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<BulkRemoveMyMetadataResponse>): BulkRemoveMyMetadataResponse {
    const message = createBaseBulkRemoveMyMetadataResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListMyRefreshTokensRequest(): ListMyRefreshTokensRequest {
  return {};
}

export const ListMyRefreshTokensRequest = {
  encode(_: ListMyRefreshTokensRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyRefreshTokensRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyRefreshTokensRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListMyRefreshTokensRequest {
    return {};
  },

  toJSON(_: ListMyRefreshTokensRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListMyRefreshTokensRequest>): ListMyRefreshTokensRequest {
    return ListMyRefreshTokensRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListMyRefreshTokensRequest>): ListMyRefreshTokensRequest {
    const message = createBaseListMyRefreshTokensRequest();
    return message;
  },
};

function createBaseListMyRefreshTokensResponse(): ListMyRefreshTokensResponse {
  return { details: undefined, result: [] };
}

export const ListMyRefreshTokensResponse = {
  encode(message: ListMyRefreshTokensResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      RefreshToken.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyRefreshTokensResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyRefreshTokensResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(RefreshToken.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyRefreshTokensResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => RefreshToken.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyRefreshTokensResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? RefreshToken.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyRefreshTokensResponse>): ListMyRefreshTokensResponse {
    return ListMyRefreshTokensResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyRefreshTokensResponse>): ListMyRefreshTokensResponse {
    const message = createBaseListMyRefreshTokensResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => RefreshToken.fromPartial(e)) || [];
    return message;
  },
};

function createBaseRevokeMyRefreshTokenRequest(): RevokeMyRefreshTokenRequest {
  return { id: "" };
}

export const RevokeMyRefreshTokenRequest = {
  encode(message: RevokeMyRefreshTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RevokeMyRefreshTokenRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRevokeMyRefreshTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RevokeMyRefreshTokenRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: RevokeMyRefreshTokenRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<RevokeMyRefreshTokenRequest>): RevokeMyRefreshTokenRequest {
    return RevokeMyRefreshTokenRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RevokeMyRefreshTokenRequest>): RevokeMyRefreshTokenRequest {
    const message = createBaseRevokeMyRefreshTokenRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseRevokeMyRefreshTokenResponse(): RevokeMyRefreshTokenResponse {
  return { details: undefined };
}

export const RevokeMyRefreshTokenResponse = {
  encode(message: RevokeMyRefreshTokenResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RevokeMyRefreshTokenResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRevokeMyRefreshTokenResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RevokeMyRefreshTokenResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RevokeMyRefreshTokenResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RevokeMyRefreshTokenResponse>): RevokeMyRefreshTokenResponse {
    return RevokeMyRefreshTokenResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RevokeMyRefreshTokenResponse>): RevokeMyRefreshTokenResponse {
    const message = createBaseRevokeMyRefreshTokenResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRevokeAllMyRefreshTokensRequest(): RevokeAllMyRefreshTokensRequest {
  return {};
}

export const RevokeAllMyRefreshTokensRequest = {
  encode(_: RevokeAllMyRefreshTokensRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RevokeAllMyRefreshTokensRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRevokeAllMyRefreshTokensRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RevokeAllMyRefreshTokensRequest {
    return {};
  },

  toJSON(_: RevokeAllMyRefreshTokensRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RevokeAllMyRefreshTokensRequest>): RevokeAllMyRefreshTokensRequest {
    return RevokeAllMyRefreshTokensRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RevokeAllMyRefreshTokensRequest>): RevokeAllMyRefreshTokensRequest {
    const message = createBaseRevokeAllMyRefreshTokensRequest();
    return message;
  },
};

function createBaseRevokeAllMyRefreshTokensResponse(): RevokeAllMyRefreshTokensResponse {
  return {};
}

export const RevokeAllMyRefreshTokensResponse = {
  encode(_: RevokeAllMyRefreshTokensResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RevokeAllMyRefreshTokensResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRevokeAllMyRefreshTokensResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RevokeAllMyRefreshTokensResponse {
    return {};
  },

  toJSON(_: RevokeAllMyRefreshTokensResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RevokeAllMyRefreshTokensResponse>): RevokeAllMyRefreshTokensResponse {
    return RevokeAllMyRefreshTokensResponse.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RevokeAllMyRefreshTokensResponse>): RevokeAllMyRefreshTokensResponse {
    const message = createBaseRevokeAllMyRefreshTokensResponse();
    return message;
  },
};

function createBaseUpdateMyUserNameRequest(): UpdateMyUserNameRequest {
  return { userName: "" };
}

export const UpdateMyUserNameRequest = {
  encode(message: UpdateMyUserNameRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userName !== "") {
      writer.uint32(10).string(message.userName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateMyUserNameRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateMyUserNameRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateMyUserNameRequest {
    return { userName: isSet(object.userName) ? String(object.userName) : "" };
  },

  toJSON(message: UpdateMyUserNameRequest): unknown {
    const obj: any = {};
    message.userName !== undefined && (obj.userName = message.userName);
    return obj;
  },

  create(base?: DeepPartial<UpdateMyUserNameRequest>): UpdateMyUserNameRequest {
    return UpdateMyUserNameRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateMyUserNameRequest>): UpdateMyUserNameRequest {
    const message = createBaseUpdateMyUserNameRequest();
    message.userName = object.userName ?? "";
    return message;
  },
};

function createBaseUpdateMyUserNameResponse(): UpdateMyUserNameResponse {
  return { details: undefined };
}

export const UpdateMyUserNameResponse = {
  encode(message: UpdateMyUserNameResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateMyUserNameResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateMyUserNameResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateMyUserNameResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateMyUserNameResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateMyUserNameResponse>): UpdateMyUserNameResponse {
    return UpdateMyUserNameResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateMyUserNameResponse>): UpdateMyUserNameResponse {
    const message = createBaseUpdateMyUserNameResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetMyPasswordComplexityPolicyRequest(): GetMyPasswordComplexityPolicyRequest {
  return {};
}

export const GetMyPasswordComplexityPolicyRequest = {
  encode(_: GetMyPasswordComplexityPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyPasswordComplexityPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyPasswordComplexityPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyPasswordComplexityPolicyRequest {
    return {};
  },

  toJSON(_: GetMyPasswordComplexityPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyPasswordComplexityPolicyRequest>): GetMyPasswordComplexityPolicyRequest {
    return GetMyPasswordComplexityPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyPasswordComplexityPolicyRequest>): GetMyPasswordComplexityPolicyRequest {
    const message = createBaseGetMyPasswordComplexityPolicyRequest();
    return message;
  },
};

function createBaseGetMyPasswordComplexityPolicyResponse(): GetMyPasswordComplexityPolicyResponse {
  return { policy: undefined };
}

export const GetMyPasswordComplexityPolicyResponse = {
  encode(message: GetMyPasswordComplexityPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      PasswordComplexityPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyPasswordComplexityPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyPasswordComplexityPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = PasswordComplexityPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyPasswordComplexityPolicyResponse {
    return { policy: isSet(object.policy) ? PasswordComplexityPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetMyPasswordComplexityPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined &&
      (obj.policy = message.policy ? PasswordComplexityPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyPasswordComplexityPolicyResponse>): GetMyPasswordComplexityPolicyResponse {
    return GetMyPasswordComplexityPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyPasswordComplexityPolicyResponse>): GetMyPasswordComplexityPolicyResponse {
    const message = createBaseGetMyPasswordComplexityPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? PasswordComplexityPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseUpdateMyPasswordRequest(): UpdateMyPasswordRequest {
  return { oldPassword: "", newPassword: "" };
}

export const UpdateMyPasswordRequest = {
  encode(message: UpdateMyPasswordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.oldPassword !== "") {
      writer.uint32(10).string(message.oldPassword);
    }
    if (message.newPassword !== "") {
      writer.uint32(18).string(message.newPassword);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateMyPasswordRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateMyPasswordRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.oldPassword = reader.string();
          break;
        case 2:
          message.newPassword = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateMyPasswordRequest {
    return {
      oldPassword: isSet(object.oldPassword) ? String(object.oldPassword) : "",
      newPassword: isSet(object.newPassword) ? String(object.newPassword) : "",
    };
  },

  toJSON(message: UpdateMyPasswordRequest): unknown {
    const obj: any = {};
    message.oldPassword !== undefined && (obj.oldPassword = message.oldPassword);
    message.newPassword !== undefined && (obj.newPassword = message.newPassword);
    return obj;
  },

  create(base?: DeepPartial<UpdateMyPasswordRequest>): UpdateMyPasswordRequest {
    return UpdateMyPasswordRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateMyPasswordRequest>): UpdateMyPasswordRequest {
    const message = createBaseUpdateMyPasswordRequest();
    message.oldPassword = object.oldPassword ?? "";
    message.newPassword = object.newPassword ?? "";
    return message;
  },
};

function createBaseUpdateMyPasswordResponse(): UpdateMyPasswordResponse {
  return { details: undefined };
}

export const UpdateMyPasswordResponse = {
  encode(message: UpdateMyPasswordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateMyPasswordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateMyPasswordResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateMyPasswordResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateMyPasswordResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateMyPasswordResponse>): UpdateMyPasswordResponse {
    return UpdateMyPasswordResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateMyPasswordResponse>): UpdateMyPasswordResponse {
    const message = createBaseUpdateMyPasswordResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetMyProfileRequest(): GetMyProfileRequest {
  return {};
}

export const GetMyProfileRequest = {
  encode(_: GetMyProfileRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyProfileRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyProfileRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyProfileRequest {
    return {};
  },

  toJSON(_: GetMyProfileRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyProfileRequest>): GetMyProfileRequest {
    return GetMyProfileRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyProfileRequest>): GetMyProfileRequest {
    const message = createBaseGetMyProfileRequest();
    return message;
  },
};

function createBaseGetMyProfileResponse(): GetMyProfileResponse {
  return { details: undefined, profile: undefined };
}

export const GetMyProfileResponse = {
  encode(message: GetMyProfileResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.profile !== undefined) {
      Profile.encode(message.profile, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyProfileResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyProfileResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.profile = Profile.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyProfileResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      profile: isSet(object.profile) ? Profile.fromJSON(object.profile) : undefined,
    };
  },

  toJSON(message: GetMyProfileResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.profile !== undefined && (obj.profile = message.profile ? Profile.toJSON(message.profile) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyProfileResponse>): GetMyProfileResponse {
    return GetMyProfileResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyProfileResponse>): GetMyProfileResponse {
    const message = createBaseGetMyProfileResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.profile = (object.profile !== undefined && object.profile !== null)
      ? Profile.fromPartial(object.profile)
      : undefined;
    return message;
  },
};

function createBaseUpdateMyProfileRequest(): UpdateMyProfileRequest {
  return { firstName: "", lastName: "", nickName: "", displayName: "", preferredLanguage: "", gender: 0 };
}

export const UpdateMyProfileRequest = {
  encode(message: UpdateMyProfileRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.firstName !== "") {
      writer.uint32(10).string(message.firstName);
    }
    if (message.lastName !== "") {
      writer.uint32(18).string(message.lastName);
    }
    if (message.nickName !== "") {
      writer.uint32(26).string(message.nickName);
    }
    if (message.displayName !== "") {
      writer.uint32(34).string(message.displayName);
    }
    if (message.preferredLanguage !== "") {
      writer.uint32(42).string(message.preferredLanguage);
    }
    if (message.gender !== 0) {
      writer.uint32(48).int32(message.gender);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateMyProfileRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateMyProfileRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.firstName = reader.string();
          break;
        case 2:
          message.lastName = reader.string();
          break;
        case 3:
          message.nickName = reader.string();
          break;
        case 4:
          message.displayName = reader.string();
          break;
        case 5:
          message.preferredLanguage = reader.string();
          break;
        case 6:
          message.gender = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateMyProfileRequest {
    return {
      firstName: isSet(object.firstName) ? String(object.firstName) : "",
      lastName: isSet(object.lastName) ? String(object.lastName) : "",
      nickName: isSet(object.nickName) ? String(object.nickName) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      preferredLanguage: isSet(object.preferredLanguage) ? String(object.preferredLanguage) : "",
      gender: isSet(object.gender) ? genderFromJSON(object.gender) : 0,
    };
  },

  toJSON(message: UpdateMyProfileRequest): unknown {
    const obj: any = {};
    message.firstName !== undefined && (obj.firstName = message.firstName);
    message.lastName !== undefined && (obj.lastName = message.lastName);
    message.nickName !== undefined && (obj.nickName = message.nickName);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.preferredLanguage !== undefined && (obj.preferredLanguage = message.preferredLanguage);
    message.gender !== undefined && (obj.gender = genderToJSON(message.gender));
    return obj;
  },

  create(base?: DeepPartial<UpdateMyProfileRequest>): UpdateMyProfileRequest {
    return UpdateMyProfileRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateMyProfileRequest>): UpdateMyProfileRequest {
    const message = createBaseUpdateMyProfileRequest();
    message.firstName = object.firstName ?? "";
    message.lastName = object.lastName ?? "";
    message.nickName = object.nickName ?? "";
    message.displayName = object.displayName ?? "";
    message.preferredLanguage = object.preferredLanguage ?? "";
    message.gender = object.gender ?? 0;
    return message;
  },
};

function createBaseUpdateMyProfileResponse(): UpdateMyProfileResponse {
  return { details: undefined };
}

export const UpdateMyProfileResponse = {
  encode(message: UpdateMyProfileResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateMyProfileResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateMyProfileResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateMyProfileResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: UpdateMyProfileResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UpdateMyProfileResponse>): UpdateMyProfileResponse {
    return UpdateMyProfileResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateMyProfileResponse>): UpdateMyProfileResponse {
    const message = createBaseUpdateMyProfileResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetMyEmailRequest(): GetMyEmailRequest {
  return {};
}

export const GetMyEmailRequest = {
  encode(_: GetMyEmailRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyEmailRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyEmailRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyEmailRequest {
    return {};
  },

  toJSON(_: GetMyEmailRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyEmailRequest>): GetMyEmailRequest {
    return GetMyEmailRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyEmailRequest>): GetMyEmailRequest {
    const message = createBaseGetMyEmailRequest();
    return message;
  },
};

function createBaseGetMyEmailResponse(): GetMyEmailResponse {
  return { details: undefined, email: undefined };
}

export const GetMyEmailResponse = {
  encode(message: GetMyEmailResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.email !== undefined) {
      Email.encode(message.email, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyEmailResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyEmailResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.email = Email.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyEmailResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      email: isSet(object.email) ? Email.fromJSON(object.email) : undefined,
    };
  },

  toJSON(message: GetMyEmailResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.email !== undefined && (obj.email = message.email ? Email.toJSON(message.email) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyEmailResponse>): GetMyEmailResponse {
    return GetMyEmailResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyEmailResponse>): GetMyEmailResponse {
    const message = createBaseGetMyEmailResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.email = (object.email !== undefined && object.email !== null) ? Email.fromPartial(object.email) : undefined;
    return message;
  },
};

function createBaseSetMyEmailRequest(): SetMyEmailRequest {
  return { email: "" };
}

export const SetMyEmailRequest = {
  encode(message: SetMyEmailRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.email !== "") {
      writer.uint32(10).string(message.email);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetMyEmailRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetMyEmailRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.email = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetMyEmailRequest {
    return { email: isSet(object.email) ? String(object.email) : "" };
  },

  toJSON(message: SetMyEmailRequest): unknown {
    const obj: any = {};
    message.email !== undefined && (obj.email = message.email);
    return obj;
  },

  create(base?: DeepPartial<SetMyEmailRequest>): SetMyEmailRequest {
    return SetMyEmailRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetMyEmailRequest>): SetMyEmailRequest {
    const message = createBaseSetMyEmailRequest();
    message.email = object.email ?? "";
    return message;
  },
};

function createBaseSetMyEmailResponse(): SetMyEmailResponse {
  return { details: undefined };
}

export const SetMyEmailResponse = {
  encode(message: SetMyEmailResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetMyEmailResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetMyEmailResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetMyEmailResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetMyEmailResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetMyEmailResponse>): SetMyEmailResponse {
    return SetMyEmailResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetMyEmailResponse>): SetMyEmailResponse {
    const message = createBaseSetMyEmailResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseVerifyMyEmailRequest(): VerifyMyEmailRequest {
  return { code: "" };
}

export const VerifyMyEmailRequest = {
  encode(message: VerifyMyEmailRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.code !== "") {
      writer.uint32(10).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyEmailRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyEmailRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.code = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyEmailRequest {
    return { code: isSet(object.code) ? String(object.code) : "" };
  },

  toJSON(message: VerifyMyEmailRequest): unknown {
    const obj: any = {};
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyEmailRequest>): VerifyMyEmailRequest {
    return VerifyMyEmailRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyEmailRequest>): VerifyMyEmailRequest {
    const message = createBaseVerifyMyEmailRequest();
    message.code = object.code ?? "";
    return message;
  },
};

function createBaseVerifyMyEmailResponse(): VerifyMyEmailResponse {
  return { details: undefined };
}

export const VerifyMyEmailResponse = {
  encode(message: VerifyMyEmailResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyEmailResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyEmailResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyEmailResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyMyEmailResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyEmailResponse>): VerifyMyEmailResponse {
    return VerifyMyEmailResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyEmailResponse>): VerifyMyEmailResponse {
    const message = createBaseVerifyMyEmailResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResendMyEmailVerificationRequest(): ResendMyEmailVerificationRequest {
  return {};
}

export const ResendMyEmailVerificationRequest = {
  encode(_: ResendMyEmailVerificationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResendMyEmailVerificationRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResendMyEmailVerificationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ResendMyEmailVerificationRequest {
    return {};
  },

  toJSON(_: ResendMyEmailVerificationRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ResendMyEmailVerificationRequest>): ResendMyEmailVerificationRequest {
    return ResendMyEmailVerificationRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ResendMyEmailVerificationRequest>): ResendMyEmailVerificationRequest {
    const message = createBaseResendMyEmailVerificationRequest();
    return message;
  },
};

function createBaseResendMyEmailVerificationResponse(): ResendMyEmailVerificationResponse {
  return { details: undefined };
}

export const ResendMyEmailVerificationResponse = {
  encode(message: ResendMyEmailVerificationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResendMyEmailVerificationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResendMyEmailVerificationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResendMyEmailVerificationResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResendMyEmailVerificationResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResendMyEmailVerificationResponse>): ResendMyEmailVerificationResponse {
    return ResendMyEmailVerificationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResendMyEmailVerificationResponse>): ResendMyEmailVerificationResponse {
    const message = createBaseResendMyEmailVerificationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseGetMyPhoneRequest(): GetMyPhoneRequest {
  return {};
}

export const GetMyPhoneRequest = {
  encode(_: GetMyPhoneRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyPhoneRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyPhoneRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyPhoneRequest {
    return {};
  },

  toJSON(_: GetMyPhoneRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyPhoneRequest>): GetMyPhoneRequest {
    return GetMyPhoneRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyPhoneRequest>): GetMyPhoneRequest {
    const message = createBaseGetMyPhoneRequest();
    return message;
  },
};

function createBaseGetMyPhoneResponse(): GetMyPhoneResponse {
  return { details: undefined, phone: undefined };
}

export const GetMyPhoneResponse = {
  encode(message: GetMyPhoneResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.phone !== undefined) {
      Phone.encode(message.phone, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyPhoneResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyPhoneResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.phone = Phone.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyPhoneResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      phone: isSet(object.phone) ? Phone.fromJSON(object.phone) : undefined,
    };
  },

  toJSON(message: GetMyPhoneResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.phone !== undefined && (obj.phone = message.phone ? Phone.toJSON(message.phone) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyPhoneResponse>): GetMyPhoneResponse {
    return GetMyPhoneResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyPhoneResponse>): GetMyPhoneResponse {
    const message = createBaseGetMyPhoneResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.phone = (object.phone !== undefined && object.phone !== null) ? Phone.fromPartial(object.phone) : undefined;
    return message;
  },
};

function createBaseSetMyPhoneRequest(): SetMyPhoneRequest {
  return { phone: "" };
}

export const SetMyPhoneRequest = {
  encode(message: SetMyPhoneRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.phone !== "") {
      writer.uint32(10).string(message.phone);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetMyPhoneRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetMyPhoneRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.phone = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetMyPhoneRequest {
    return { phone: isSet(object.phone) ? String(object.phone) : "" };
  },

  toJSON(message: SetMyPhoneRequest): unknown {
    const obj: any = {};
    message.phone !== undefined && (obj.phone = message.phone);
    return obj;
  },

  create(base?: DeepPartial<SetMyPhoneRequest>): SetMyPhoneRequest {
    return SetMyPhoneRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetMyPhoneRequest>): SetMyPhoneRequest {
    const message = createBaseSetMyPhoneRequest();
    message.phone = object.phone ?? "";
    return message;
  },
};

function createBaseSetMyPhoneResponse(): SetMyPhoneResponse {
  return { details: undefined };
}

export const SetMyPhoneResponse = {
  encode(message: SetMyPhoneResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetMyPhoneResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetMyPhoneResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetMyPhoneResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetMyPhoneResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetMyPhoneResponse>): SetMyPhoneResponse {
    return SetMyPhoneResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetMyPhoneResponse>): SetMyPhoneResponse {
    const message = createBaseSetMyPhoneResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseVerifyMyPhoneRequest(): VerifyMyPhoneRequest {
  return { code: "" };
}

export const VerifyMyPhoneRequest = {
  encode(message: VerifyMyPhoneRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.code !== "") {
      writer.uint32(10).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyPhoneRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyPhoneRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.code = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyPhoneRequest {
    return { code: isSet(object.code) ? String(object.code) : "" };
  },

  toJSON(message: VerifyMyPhoneRequest): unknown {
    const obj: any = {};
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyPhoneRequest>): VerifyMyPhoneRequest {
    return VerifyMyPhoneRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyPhoneRequest>): VerifyMyPhoneRequest {
    const message = createBaseVerifyMyPhoneRequest();
    message.code = object.code ?? "";
    return message;
  },
};

function createBaseVerifyMyPhoneResponse(): VerifyMyPhoneResponse {
  return { details: undefined };
}

export const VerifyMyPhoneResponse = {
  encode(message: VerifyMyPhoneResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyPhoneResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyPhoneResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyPhoneResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyMyPhoneResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyPhoneResponse>): VerifyMyPhoneResponse {
    return VerifyMyPhoneResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyPhoneResponse>): VerifyMyPhoneResponse {
    const message = createBaseVerifyMyPhoneResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResendMyPhoneVerificationRequest(): ResendMyPhoneVerificationRequest {
  return {};
}

export const ResendMyPhoneVerificationRequest = {
  encode(_: ResendMyPhoneVerificationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResendMyPhoneVerificationRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResendMyPhoneVerificationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ResendMyPhoneVerificationRequest {
    return {};
  },

  toJSON(_: ResendMyPhoneVerificationRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ResendMyPhoneVerificationRequest>): ResendMyPhoneVerificationRequest {
    return ResendMyPhoneVerificationRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ResendMyPhoneVerificationRequest>): ResendMyPhoneVerificationRequest {
    const message = createBaseResendMyPhoneVerificationRequest();
    return message;
  },
};

function createBaseResendMyPhoneVerificationResponse(): ResendMyPhoneVerificationResponse {
  return { details: undefined };
}

export const ResendMyPhoneVerificationResponse = {
  encode(message: ResendMyPhoneVerificationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResendMyPhoneVerificationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResendMyPhoneVerificationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ResendMyPhoneVerificationResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: ResendMyPhoneVerificationResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResendMyPhoneVerificationResponse>): ResendMyPhoneVerificationResponse {
    return ResendMyPhoneVerificationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResendMyPhoneVerificationResponse>): ResendMyPhoneVerificationResponse {
    const message = createBaseResendMyPhoneVerificationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveMyPhoneRequest(): RemoveMyPhoneRequest {
  return {};
}

export const RemoveMyPhoneRequest = {
  encode(_: RemoveMyPhoneRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyPhoneRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyPhoneRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveMyPhoneRequest {
    return {};
  },

  toJSON(_: RemoveMyPhoneRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveMyPhoneRequest>): RemoveMyPhoneRequest {
    return RemoveMyPhoneRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveMyPhoneRequest>): RemoveMyPhoneRequest {
    const message = createBaseRemoveMyPhoneRequest();
    return message;
  },
};

function createBaseRemoveMyPhoneResponse(): RemoveMyPhoneResponse {
  return { details: undefined };
}

export const RemoveMyPhoneResponse = {
  encode(message: RemoveMyPhoneResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyPhoneResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyPhoneResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyPhoneResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMyPhoneResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyPhoneResponse>): RemoveMyPhoneResponse {
    return RemoveMyPhoneResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyPhoneResponse>): RemoveMyPhoneResponse {
    const message = createBaseRemoveMyPhoneResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveMyAvatarRequest(): RemoveMyAvatarRequest {
  return {};
}

export const RemoveMyAvatarRequest = {
  encode(_: RemoveMyAvatarRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyAvatarRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyAvatarRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveMyAvatarRequest {
    return {};
  },

  toJSON(_: RemoveMyAvatarRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveMyAvatarRequest>): RemoveMyAvatarRequest {
    return RemoveMyAvatarRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveMyAvatarRequest>): RemoveMyAvatarRequest {
    const message = createBaseRemoveMyAvatarRequest();
    return message;
  },
};

function createBaseRemoveMyAvatarResponse(): RemoveMyAvatarResponse {
  return { details: undefined };
}

export const RemoveMyAvatarResponse = {
  encode(message: RemoveMyAvatarResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyAvatarResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyAvatarResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyAvatarResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMyAvatarResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyAvatarResponse>): RemoveMyAvatarResponse {
    return RemoveMyAvatarResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyAvatarResponse>): RemoveMyAvatarResponse {
    const message = createBaseRemoveMyAvatarResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListMyLinkedIDPsRequest(): ListMyLinkedIDPsRequest {
  return { query: undefined };
}

export const ListMyLinkedIDPsRequest = {
  encode(message: ListMyLinkedIDPsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyLinkedIDPsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyLinkedIDPsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyLinkedIDPsRequest {
    return { query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined };
  },

  toJSON(message: ListMyLinkedIDPsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ListMyLinkedIDPsRequest>): ListMyLinkedIDPsRequest {
    return ListMyLinkedIDPsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyLinkedIDPsRequest>): ListMyLinkedIDPsRequest {
    const message = createBaseListMyLinkedIDPsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseListMyLinkedIDPsResponse(): ListMyLinkedIDPsResponse {
  return { details: undefined, result: [] };
}

export const ListMyLinkedIDPsResponse = {
  encode(message: ListMyLinkedIDPsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      IDPUserLink.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyLinkedIDPsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyLinkedIDPsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(IDPUserLink.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyLinkedIDPsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => IDPUserLink.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyLinkedIDPsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? IDPUserLink.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyLinkedIDPsResponse>): ListMyLinkedIDPsResponse {
    return ListMyLinkedIDPsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyLinkedIDPsResponse>): ListMyLinkedIDPsResponse {
    const message = createBaseListMyLinkedIDPsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => IDPUserLink.fromPartial(e)) || [];
    return message;
  },
};

function createBaseRemoveMyLinkedIDPRequest(): RemoveMyLinkedIDPRequest {
  return { idpId: "", linkedUserId: "" };
}

export const RemoveMyLinkedIDPRequest = {
  encode(message: RemoveMyLinkedIDPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.linkedUserId !== "") {
      writer.uint32(18).string(message.linkedUserId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyLinkedIDPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyLinkedIDPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.idpId = reader.string();
          break;
        case 2:
          message.linkedUserId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyLinkedIDPRequest {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      linkedUserId: isSet(object.linkedUserId) ? String(object.linkedUserId) : "",
    };
  },

  toJSON(message: RemoveMyLinkedIDPRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.linkedUserId !== undefined && (obj.linkedUserId = message.linkedUserId);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyLinkedIDPRequest>): RemoveMyLinkedIDPRequest {
    return RemoveMyLinkedIDPRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyLinkedIDPRequest>): RemoveMyLinkedIDPRequest {
    const message = createBaseRemoveMyLinkedIDPRequest();
    message.idpId = object.idpId ?? "";
    message.linkedUserId = object.linkedUserId ?? "";
    return message;
  },
};

function createBaseRemoveMyLinkedIDPResponse(): RemoveMyLinkedIDPResponse {
  return { details: undefined };
}

export const RemoveMyLinkedIDPResponse = {
  encode(message: RemoveMyLinkedIDPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyLinkedIDPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyLinkedIDPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyLinkedIDPResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMyLinkedIDPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyLinkedIDPResponse>): RemoveMyLinkedIDPResponse {
    return RemoveMyLinkedIDPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyLinkedIDPResponse>): RemoveMyLinkedIDPResponse {
    const message = createBaseRemoveMyLinkedIDPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListMyAuthFactorsRequest(): ListMyAuthFactorsRequest {
  return {};
}

export const ListMyAuthFactorsRequest = {
  encode(_: ListMyAuthFactorsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyAuthFactorsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyAuthFactorsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListMyAuthFactorsRequest {
    return {};
  },

  toJSON(_: ListMyAuthFactorsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListMyAuthFactorsRequest>): ListMyAuthFactorsRequest {
    return ListMyAuthFactorsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListMyAuthFactorsRequest>): ListMyAuthFactorsRequest {
    const message = createBaseListMyAuthFactorsRequest();
    return message;
  },
};

function createBaseListMyAuthFactorsResponse(): ListMyAuthFactorsResponse {
  return { result: [] };
}

export const ListMyAuthFactorsResponse = {
  encode(message: ListMyAuthFactorsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.result) {
      AuthFactor.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyAuthFactorsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyAuthFactorsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.result.push(AuthFactor.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyAuthFactorsResponse {
    return { result: Array.isArray(object?.result) ? object.result.map((e: any) => AuthFactor.fromJSON(e)) : [] };
  },

  toJSON(message: ListMyAuthFactorsResponse): unknown {
    const obj: any = {};
    if (message.result) {
      obj.result = message.result.map((e) => e ? AuthFactor.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyAuthFactorsResponse>): ListMyAuthFactorsResponse {
    return ListMyAuthFactorsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyAuthFactorsResponse>): ListMyAuthFactorsResponse {
    const message = createBaseListMyAuthFactorsResponse();
    message.result = object.result?.map((e) => AuthFactor.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAddMyAuthFactorU2FRequest(): AddMyAuthFactorU2FRequest {
  return {};
}

export const AddMyAuthFactorU2FRequest = {
  encode(_: AddMyAuthFactorU2FRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMyAuthFactorU2FRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMyAuthFactorU2FRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): AddMyAuthFactorU2FRequest {
    return {};
  },

  toJSON(_: AddMyAuthFactorU2FRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<AddMyAuthFactorU2FRequest>): AddMyAuthFactorU2FRequest {
    return AddMyAuthFactorU2FRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<AddMyAuthFactorU2FRequest>): AddMyAuthFactorU2FRequest {
    const message = createBaseAddMyAuthFactorU2FRequest();
    return message;
  },
};

function createBaseAddMyAuthFactorU2FResponse(): AddMyAuthFactorU2FResponse {
  return { key: undefined, details: undefined };
}

export const AddMyAuthFactorU2FResponse = {
  encode(message: AddMyAuthFactorU2FResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== undefined) {
      WebAuthNKey.encode(message.key, writer.uint32(10).fork()).ldelim();
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMyAuthFactorU2FResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMyAuthFactorU2FResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = WebAuthNKey.decode(reader, reader.uint32());
          break;
        case 2:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddMyAuthFactorU2FResponse {
    return {
      key: isSet(object.key) ? WebAuthNKey.fromJSON(object.key) : undefined,
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
    };
  },

  toJSON(message: AddMyAuthFactorU2FResponse): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key ? WebAuthNKey.toJSON(message.key) : undefined);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddMyAuthFactorU2FResponse>): AddMyAuthFactorU2FResponse {
    return AddMyAuthFactorU2FResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddMyAuthFactorU2FResponse>): AddMyAuthFactorU2FResponse {
    const message = createBaseAddMyAuthFactorU2FResponse();
    message.key = (object.key !== undefined && object.key !== null) ? WebAuthNKey.fromPartial(object.key) : undefined;
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddMyAuthFactorOTPRequest(): AddMyAuthFactorOTPRequest {
  return {};
}

export const AddMyAuthFactorOTPRequest = {
  encode(_: AddMyAuthFactorOTPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMyAuthFactorOTPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMyAuthFactorOTPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): AddMyAuthFactorOTPRequest {
    return {};
  },

  toJSON(_: AddMyAuthFactorOTPRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<AddMyAuthFactorOTPRequest>): AddMyAuthFactorOTPRequest {
    return AddMyAuthFactorOTPRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<AddMyAuthFactorOTPRequest>): AddMyAuthFactorOTPRequest {
    const message = createBaseAddMyAuthFactorOTPRequest();
    return message;
  },
};

function createBaseAddMyAuthFactorOTPResponse(): AddMyAuthFactorOTPResponse {
  return { url: "", secret: "", details: undefined };
}

export const AddMyAuthFactorOTPResponse = {
  encode(message: AddMyAuthFactorOTPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.url !== "") {
      writer.uint32(10).string(message.url);
    }
    if (message.secret !== "") {
      writer.uint32(18).string(message.secret);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMyAuthFactorOTPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMyAuthFactorOTPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.url = reader.string();
          break;
        case 2:
          message.secret = reader.string();
          break;
        case 3:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddMyAuthFactorOTPResponse {
    return {
      url: isSet(object.url) ? String(object.url) : "",
      secret: isSet(object.secret) ? String(object.secret) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
    };
  },

  toJSON(message: AddMyAuthFactorOTPResponse): unknown {
    const obj: any = {};
    message.url !== undefined && (obj.url = message.url);
    message.secret !== undefined && (obj.secret = message.secret);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddMyAuthFactorOTPResponse>): AddMyAuthFactorOTPResponse {
    return AddMyAuthFactorOTPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddMyAuthFactorOTPResponse>): AddMyAuthFactorOTPResponse {
    const message = createBaseAddMyAuthFactorOTPResponse();
    message.url = object.url ?? "";
    message.secret = object.secret ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseVerifyMyAuthFactorOTPRequest(): VerifyMyAuthFactorOTPRequest {
  return { code: "" };
}

export const VerifyMyAuthFactorOTPRequest = {
  encode(message: VerifyMyAuthFactorOTPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.code !== "") {
      writer.uint32(10).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyAuthFactorOTPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyAuthFactorOTPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.code = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyAuthFactorOTPRequest {
    return { code: isSet(object.code) ? String(object.code) : "" };
  },

  toJSON(message: VerifyMyAuthFactorOTPRequest): unknown {
    const obj: any = {};
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyAuthFactorOTPRequest>): VerifyMyAuthFactorOTPRequest {
    return VerifyMyAuthFactorOTPRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyAuthFactorOTPRequest>): VerifyMyAuthFactorOTPRequest {
    const message = createBaseVerifyMyAuthFactorOTPRequest();
    message.code = object.code ?? "";
    return message;
  },
};

function createBaseVerifyMyAuthFactorOTPResponse(): VerifyMyAuthFactorOTPResponse {
  return { details: undefined };
}

export const VerifyMyAuthFactorOTPResponse = {
  encode(message: VerifyMyAuthFactorOTPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyAuthFactorOTPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyAuthFactorOTPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyAuthFactorOTPResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyMyAuthFactorOTPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyAuthFactorOTPResponse>): VerifyMyAuthFactorOTPResponse {
    return VerifyMyAuthFactorOTPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyAuthFactorOTPResponse>): VerifyMyAuthFactorOTPResponse {
    const message = createBaseVerifyMyAuthFactorOTPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseVerifyMyAuthFactorU2FRequest(): VerifyMyAuthFactorU2FRequest {
  return { verification: undefined };
}

export const VerifyMyAuthFactorU2FRequest = {
  encode(message: VerifyMyAuthFactorU2FRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.verification !== undefined) {
      WebAuthNVerification.encode(message.verification, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyAuthFactorU2FRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyAuthFactorU2FRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.verification = WebAuthNVerification.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyAuthFactorU2FRequest {
    return {
      verification: isSet(object.verification) ? WebAuthNVerification.fromJSON(object.verification) : undefined,
    };
  },

  toJSON(message: VerifyMyAuthFactorU2FRequest): unknown {
    const obj: any = {};
    message.verification !== undefined &&
      (obj.verification = message.verification ? WebAuthNVerification.toJSON(message.verification) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyAuthFactorU2FRequest>): VerifyMyAuthFactorU2FRequest {
    return VerifyMyAuthFactorU2FRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyAuthFactorU2FRequest>): VerifyMyAuthFactorU2FRequest {
    const message = createBaseVerifyMyAuthFactorU2FRequest();
    message.verification = (object.verification !== undefined && object.verification !== null)
      ? WebAuthNVerification.fromPartial(object.verification)
      : undefined;
    return message;
  },
};

function createBaseVerifyMyAuthFactorU2FResponse(): VerifyMyAuthFactorU2FResponse {
  return { details: undefined };
}

export const VerifyMyAuthFactorU2FResponse = {
  encode(message: VerifyMyAuthFactorU2FResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyAuthFactorU2FResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyAuthFactorU2FResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyAuthFactorU2FResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyMyAuthFactorU2FResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyAuthFactorU2FResponse>): VerifyMyAuthFactorU2FResponse {
    return VerifyMyAuthFactorU2FResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyAuthFactorU2FResponse>): VerifyMyAuthFactorU2FResponse {
    const message = createBaseVerifyMyAuthFactorU2FResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveMyAuthFactorOTPRequest(): RemoveMyAuthFactorOTPRequest {
  return {};
}

export const RemoveMyAuthFactorOTPRequest = {
  encode(_: RemoveMyAuthFactorOTPRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyAuthFactorOTPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyAuthFactorOTPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveMyAuthFactorOTPRequest {
    return {};
  },

  toJSON(_: RemoveMyAuthFactorOTPRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RemoveMyAuthFactorOTPRequest>): RemoveMyAuthFactorOTPRequest {
    return RemoveMyAuthFactorOTPRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RemoveMyAuthFactorOTPRequest>): RemoveMyAuthFactorOTPRequest {
    const message = createBaseRemoveMyAuthFactorOTPRequest();
    return message;
  },
};

function createBaseRemoveMyAuthFactorOTPResponse(): RemoveMyAuthFactorOTPResponse {
  return { details: undefined };
}

export const RemoveMyAuthFactorOTPResponse = {
  encode(message: RemoveMyAuthFactorOTPResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyAuthFactorOTPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyAuthFactorOTPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyAuthFactorOTPResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMyAuthFactorOTPResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyAuthFactorOTPResponse>): RemoveMyAuthFactorOTPResponse {
    return RemoveMyAuthFactorOTPResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyAuthFactorOTPResponse>): RemoveMyAuthFactorOTPResponse {
    const message = createBaseRemoveMyAuthFactorOTPResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveMyAuthFactorU2FRequest(): RemoveMyAuthFactorU2FRequest {
  return { tokenId: "" };
}

export const RemoveMyAuthFactorU2FRequest = {
  encode(message: RemoveMyAuthFactorU2FRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tokenId !== "") {
      writer.uint32(10).string(message.tokenId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyAuthFactorU2FRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyAuthFactorU2FRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tokenId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyAuthFactorU2FRequest {
    return { tokenId: isSet(object.tokenId) ? String(object.tokenId) : "" };
  },

  toJSON(message: RemoveMyAuthFactorU2FRequest): unknown {
    const obj: any = {};
    message.tokenId !== undefined && (obj.tokenId = message.tokenId);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyAuthFactorU2FRequest>): RemoveMyAuthFactorU2FRequest {
    return RemoveMyAuthFactorU2FRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyAuthFactorU2FRequest>): RemoveMyAuthFactorU2FRequest {
    const message = createBaseRemoveMyAuthFactorU2FRequest();
    message.tokenId = object.tokenId ?? "";
    return message;
  },
};

function createBaseRemoveMyAuthFactorU2FResponse(): RemoveMyAuthFactorU2FResponse {
  return { details: undefined };
}

export const RemoveMyAuthFactorU2FResponse = {
  encode(message: RemoveMyAuthFactorU2FResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyAuthFactorU2FResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyAuthFactorU2FResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyAuthFactorU2FResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMyAuthFactorU2FResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyAuthFactorU2FResponse>): RemoveMyAuthFactorU2FResponse {
    return RemoveMyAuthFactorU2FResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyAuthFactorU2FResponse>): RemoveMyAuthFactorU2FResponse {
    const message = createBaseRemoveMyAuthFactorU2FResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListMyPasswordlessRequest(): ListMyPasswordlessRequest {
  return {};
}

export const ListMyPasswordlessRequest = {
  encode(_: ListMyPasswordlessRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyPasswordlessRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyPasswordlessRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListMyPasswordlessRequest {
    return {};
  },

  toJSON(_: ListMyPasswordlessRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListMyPasswordlessRequest>): ListMyPasswordlessRequest {
    return ListMyPasswordlessRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListMyPasswordlessRequest>): ListMyPasswordlessRequest {
    const message = createBaseListMyPasswordlessRequest();
    return message;
  },
};

function createBaseListMyPasswordlessResponse(): ListMyPasswordlessResponse {
  return { result: [] };
}

export const ListMyPasswordlessResponse = {
  encode(message: ListMyPasswordlessResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.result) {
      WebAuthNToken.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyPasswordlessResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyPasswordlessResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.result.push(WebAuthNToken.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyPasswordlessResponse {
    return { result: Array.isArray(object?.result) ? object.result.map((e: any) => WebAuthNToken.fromJSON(e)) : [] };
  },

  toJSON(message: ListMyPasswordlessResponse): unknown {
    const obj: any = {};
    if (message.result) {
      obj.result = message.result.map((e) => e ? WebAuthNToken.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyPasswordlessResponse>): ListMyPasswordlessResponse {
    return ListMyPasswordlessResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyPasswordlessResponse>): ListMyPasswordlessResponse {
    const message = createBaseListMyPasswordlessResponse();
    message.result = object.result?.map((e) => WebAuthNToken.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAddMyPasswordlessRequest(): AddMyPasswordlessRequest {
  return {};
}

export const AddMyPasswordlessRequest = {
  encode(_: AddMyPasswordlessRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMyPasswordlessRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMyPasswordlessRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): AddMyPasswordlessRequest {
    return {};
  },

  toJSON(_: AddMyPasswordlessRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<AddMyPasswordlessRequest>): AddMyPasswordlessRequest {
    return AddMyPasswordlessRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<AddMyPasswordlessRequest>): AddMyPasswordlessRequest {
    const message = createBaseAddMyPasswordlessRequest();
    return message;
  },
};

function createBaseAddMyPasswordlessResponse(): AddMyPasswordlessResponse {
  return { key: undefined, details: undefined };
}

export const AddMyPasswordlessResponse = {
  encode(message: AddMyPasswordlessResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== undefined) {
      WebAuthNKey.encode(message.key, writer.uint32(10).fork()).ldelim();
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMyPasswordlessResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMyPasswordlessResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.key = WebAuthNKey.decode(reader, reader.uint32());
          break;
        case 2:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddMyPasswordlessResponse {
    return {
      key: isSet(object.key) ? WebAuthNKey.fromJSON(object.key) : undefined,
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
    };
  },

  toJSON(message: AddMyPasswordlessResponse): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key ? WebAuthNKey.toJSON(message.key) : undefined);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddMyPasswordlessResponse>): AddMyPasswordlessResponse {
    return AddMyPasswordlessResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddMyPasswordlessResponse>): AddMyPasswordlessResponse {
    const message = createBaseAddMyPasswordlessResponse();
    message.key = (object.key !== undefined && object.key !== null) ? WebAuthNKey.fromPartial(object.key) : undefined;
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddMyPasswordlessLinkRequest(): AddMyPasswordlessLinkRequest {
  return {};
}

export const AddMyPasswordlessLinkRequest = {
  encode(_: AddMyPasswordlessLinkRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMyPasswordlessLinkRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMyPasswordlessLinkRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): AddMyPasswordlessLinkRequest {
    return {};
  },

  toJSON(_: AddMyPasswordlessLinkRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<AddMyPasswordlessLinkRequest>): AddMyPasswordlessLinkRequest {
    return AddMyPasswordlessLinkRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<AddMyPasswordlessLinkRequest>): AddMyPasswordlessLinkRequest {
    const message = createBaseAddMyPasswordlessLinkRequest();
    return message;
  },
};

function createBaseAddMyPasswordlessLinkResponse(): AddMyPasswordlessLinkResponse {
  return { details: undefined, link: "", expiration: undefined };
}

export const AddMyPasswordlessLinkResponse = {
  encode(message: AddMyPasswordlessLinkResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.link !== "") {
      writer.uint32(18).string(message.link);
    }
    if (message.expiration !== undefined) {
      Duration.encode(message.expiration, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddMyPasswordlessLinkResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddMyPasswordlessLinkResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.link = reader.string();
          break;
        case 3:
          message.expiration = Duration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddMyPasswordlessLinkResponse {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      link: isSet(object.link) ? String(object.link) : "",
      expiration: isSet(object.expiration) ? Duration.fromJSON(object.expiration) : undefined,
    };
  },

  toJSON(message: AddMyPasswordlessLinkResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.link !== undefined && (obj.link = message.link);
    message.expiration !== undefined &&
      (obj.expiration = message.expiration ? Duration.toJSON(message.expiration) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddMyPasswordlessLinkResponse>): AddMyPasswordlessLinkResponse {
    return AddMyPasswordlessLinkResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddMyPasswordlessLinkResponse>): AddMyPasswordlessLinkResponse {
    const message = createBaseAddMyPasswordlessLinkResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.link = object.link ?? "";
    message.expiration = (object.expiration !== undefined && object.expiration !== null)
      ? Duration.fromPartial(object.expiration)
      : undefined;
    return message;
  },
};

function createBaseSendMyPasswordlessLinkRequest(): SendMyPasswordlessLinkRequest {
  return {};
}

export const SendMyPasswordlessLinkRequest = {
  encode(_: SendMyPasswordlessLinkRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendMyPasswordlessLinkRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendMyPasswordlessLinkRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): SendMyPasswordlessLinkRequest {
    return {};
  },

  toJSON(_: SendMyPasswordlessLinkRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<SendMyPasswordlessLinkRequest>): SendMyPasswordlessLinkRequest {
    return SendMyPasswordlessLinkRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<SendMyPasswordlessLinkRequest>): SendMyPasswordlessLinkRequest {
    const message = createBaseSendMyPasswordlessLinkRequest();
    return message;
  },
};

function createBaseSendMyPasswordlessLinkResponse(): SendMyPasswordlessLinkResponse {
  return { details: undefined };
}

export const SendMyPasswordlessLinkResponse = {
  encode(message: SendMyPasswordlessLinkResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendMyPasswordlessLinkResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendMyPasswordlessLinkResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SendMyPasswordlessLinkResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: SendMyPasswordlessLinkResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SendMyPasswordlessLinkResponse>): SendMyPasswordlessLinkResponse {
    return SendMyPasswordlessLinkResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SendMyPasswordlessLinkResponse>): SendMyPasswordlessLinkResponse {
    const message = createBaseSendMyPasswordlessLinkResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseVerifyMyPasswordlessRequest(): VerifyMyPasswordlessRequest {
  return { verification: undefined };
}

export const VerifyMyPasswordlessRequest = {
  encode(message: VerifyMyPasswordlessRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.verification !== undefined) {
      WebAuthNVerification.encode(message.verification, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyPasswordlessRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyPasswordlessRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.verification = WebAuthNVerification.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyPasswordlessRequest {
    return {
      verification: isSet(object.verification) ? WebAuthNVerification.fromJSON(object.verification) : undefined,
    };
  },

  toJSON(message: VerifyMyPasswordlessRequest): unknown {
    const obj: any = {};
    message.verification !== undefined &&
      (obj.verification = message.verification ? WebAuthNVerification.toJSON(message.verification) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyPasswordlessRequest>): VerifyMyPasswordlessRequest {
    return VerifyMyPasswordlessRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyPasswordlessRequest>): VerifyMyPasswordlessRequest {
    const message = createBaseVerifyMyPasswordlessRequest();
    message.verification = (object.verification !== undefined && object.verification !== null)
      ? WebAuthNVerification.fromPartial(object.verification)
      : undefined;
    return message;
  },
};

function createBaseVerifyMyPasswordlessResponse(): VerifyMyPasswordlessResponse {
  return { details: undefined };
}

export const VerifyMyPasswordlessResponse = {
  encode(message: VerifyMyPasswordlessResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMyPasswordlessResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMyPasswordlessResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMyPasswordlessResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyMyPasswordlessResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyMyPasswordlessResponse>): VerifyMyPasswordlessResponse {
    return VerifyMyPasswordlessResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMyPasswordlessResponse>): VerifyMyPasswordlessResponse {
    const message = createBaseVerifyMyPasswordlessResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveMyPasswordlessRequest(): RemoveMyPasswordlessRequest {
  return { tokenId: "" };
}

export const RemoveMyPasswordlessRequest = {
  encode(message: RemoveMyPasswordlessRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tokenId !== "") {
      writer.uint32(10).string(message.tokenId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyPasswordlessRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyPasswordlessRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tokenId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyPasswordlessRequest {
    return { tokenId: isSet(object.tokenId) ? String(object.tokenId) : "" };
  },

  toJSON(message: RemoveMyPasswordlessRequest): unknown {
    const obj: any = {};
    message.tokenId !== undefined && (obj.tokenId = message.tokenId);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyPasswordlessRequest>): RemoveMyPasswordlessRequest {
    return RemoveMyPasswordlessRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyPasswordlessRequest>): RemoveMyPasswordlessRequest {
    const message = createBaseRemoveMyPasswordlessRequest();
    message.tokenId = object.tokenId ?? "";
    return message;
  },
};

function createBaseRemoveMyPasswordlessResponse(): RemoveMyPasswordlessResponse {
  return { details: undefined };
}

export const RemoveMyPasswordlessResponse = {
  encode(message: RemoveMyPasswordlessResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveMyPasswordlessResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveMyPasswordlessResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveMyPasswordlessResponse {
    return { details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveMyPasswordlessResponse): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveMyPasswordlessResponse>): RemoveMyPasswordlessResponse {
    return RemoveMyPasswordlessResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveMyPasswordlessResponse>): RemoveMyPasswordlessResponse {
    const message = createBaseRemoveMyPasswordlessResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseListMyUserGrantsRequest(): ListMyUserGrantsRequest {
  return { query: undefined };
}

export const ListMyUserGrantsRequest = {
  encode(message: ListMyUserGrantsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyUserGrantsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyUserGrantsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyUserGrantsRequest {
    return { query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined };
  },

  toJSON(message: ListMyUserGrantsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ListMyUserGrantsRequest>): ListMyUserGrantsRequest {
    return ListMyUserGrantsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyUserGrantsRequest>): ListMyUserGrantsRequest {
    const message = createBaseListMyUserGrantsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    return message;
  },
};

function createBaseListMyUserGrantsResponse(): ListMyUserGrantsResponse {
  return { details: undefined, result: [] };
}

export const ListMyUserGrantsResponse = {
  encode(message: ListMyUserGrantsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      UserGrant.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyUserGrantsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyUserGrantsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(UserGrant.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyUserGrantsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => UserGrant.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyUserGrantsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? UserGrant.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyUserGrantsResponse>): ListMyUserGrantsResponse {
    return ListMyUserGrantsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyUserGrantsResponse>): ListMyUserGrantsResponse {
    const message = createBaseListMyUserGrantsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => UserGrant.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUserGrant(): UserGrant {
  return {
    orgId: "",
    projectId: "",
    userId: "",
    roles: [],
    orgName: "",
    grantId: "",
    details: undefined,
    orgDomain: "",
    projectName: "",
    projectGrantId: "",
    roleKeys: [],
  };
}

export const UserGrant = {
  encode(message: UserGrant, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    if (message.projectId !== "") {
      writer.uint32(18).string(message.projectId);
    }
    if (message.userId !== "") {
      writer.uint32(26).string(message.userId);
    }
    for (const v of message.roles) {
      writer.uint32(34).string(v!);
    }
    if (message.orgName !== "") {
      writer.uint32(42).string(message.orgName);
    }
    if (message.grantId !== "") {
      writer.uint32(50).string(message.grantId);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(58).fork()).ldelim();
    }
    if (message.orgDomain !== "") {
      writer.uint32(66).string(message.orgDomain);
    }
    if (message.projectName !== "") {
      writer.uint32(74).string(message.projectName);
    }
    if (message.projectGrantId !== "") {
      writer.uint32(82).string(message.projectGrantId);
    }
    for (const v of message.roleKeys) {
      writer.uint32(90).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserGrant {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserGrant();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orgId = reader.string();
          break;
        case 2:
          message.projectId = reader.string();
          break;
        case 3:
          message.userId = reader.string();
          break;
        case 4:
          message.roles.push(reader.string());
          break;
        case 5:
          message.orgName = reader.string();
          break;
        case 6:
          message.grantId = reader.string();
          break;
        case 7:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 8:
          message.orgDomain = reader.string();
          break;
        case 9:
          message.projectName = reader.string();
          break;
        case 10:
          message.projectGrantId = reader.string();
          break;
        case 11:
          message.roleKeys.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserGrant {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      projectId: isSet(object.projectId) ? String(object.projectId) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
      roles: Array.isArray(object?.roles) ? object.roles.map((e: any) => String(e)) : [],
      orgName: isSet(object.orgName) ? String(object.orgName) : "",
      grantId: isSet(object.grantId) ? String(object.grantId) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      orgDomain: isSet(object.orgDomain) ? String(object.orgDomain) : "",
      projectName: isSet(object.projectName) ? String(object.projectName) : "",
      projectGrantId: isSet(object.projectGrantId) ? String(object.projectGrantId) : "",
      roleKeys: Array.isArray(object?.roleKeys) ? object.roleKeys.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: UserGrant): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.projectId !== undefined && (obj.projectId = message.projectId);
    message.userId !== undefined && (obj.userId = message.userId);
    if (message.roles) {
      obj.roles = message.roles.map((e) => e);
    } else {
      obj.roles = [];
    }
    message.orgName !== undefined && (obj.orgName = message.orgName);
    message.grantId !== undefined && (obj.grantId = message.grantId);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.orgDomain !== undefined && (obj.orgDomain = message.orgDomain);
    message.projectName !== undefined && (obj.projectName = message.projectName);
    message.projectGrantId !== undefined && (obj.projectGrantId = message.projectGrantId);
    if (message.roleKeys) {
      obj.roleKeys = message.roleKeys.map((e) => e);
    } else {
      obj.roleKeys = [];
    }
    return obj;
  },

  create(base?: DeepPartial<UserGrant>): UserGrant {
    return UserGrant.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UserGrant>): UserGrant {
    const message = createBaseUserGrant();
    message.orgId = object.orgId ?? "";
    message.projectId = object.projectId ?? "";
    message.userId = object.userId ?? "";
    message.roles = object.roles?.map((e) => e) || [];
    message.orgName = object.orgName ?? "";
    message.grantId = object.grantId ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.orgDomain = object.orgDomain ?? "";
    message.projectName = object.projectName ?? "";
    message.projectGrantId = object.projectGrantId ?? "";
    message.roleKeys = object.roleKeys?.map((e) => e) || [];
    return message;
  },
};

function createBaseListMyProjectOrgsRequest(): ListMyProjectOrgsRequest {
  return { query: undefined, queries: [] };
}

export const ListMyProjectOrgsRequest = {
  encode(message: ListMyProjectOrgsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.queries) {
      OrgQuery.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyProjectOrgsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyProjectOrgsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.queries.push(OrgQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyProjectOrgsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => OrgQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyProjectOrgsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? OrgQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyProjectOrgsRequest>): ListMyProjectOrgsRequest {
    return ListMyProjectOrgsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyProjectOrgsRequest>): ListMyProjectOrgsRequest {
    const message = createBaseListMyProjectOrgsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.queries = object.queries?.map((e) => OrgQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListMyProjectOrgsResponse(): ListMyProjectOrgsResponse {
  return { details: undefined, result: [] };
}

export const ListMyProjectOrgsResponse = {
  encode(message: ListMyProjectOrgsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      Org.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyProjectOrgsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyProjectOrgsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(Org.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyProjectOrgsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Org.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyProjectOrgsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? Org.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyProjectOrgsResponse>): ListMyProjectOrgsResponse {
    return ListMyProjectOrgsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyProjectOrgsResponse>): ListMyProjectOrgsResponse {
    const message = createBaseListMyProjectOrgsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => Org.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListMyZitadelPermissionsRequest(): ListMyZitadelPermissionsRequest {
  return {};
}

export const ListMyZitadelPermissionsRequest = {
  encode(_: ListMyZitadelPermissionsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyZitadelPermissionsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyZitadelPermissionsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListMyZitadelPermissionsRequest {
    return {};
  },

  toJSON(_: ListMyZitadelPermissionsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListMyZitadelPermissionsRequest>): ListMyZitadelPermissionsRequest {
    return ListMyZitadelPermissionsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListMyZitadelPermissionsRequest>): ListMyZitadelPermissionsRequest {
    const message = createBaseListMyZitadelPermissionsRequest();
    return message;
  },
};

function createBaseListMyZitadelPermissionsResponse(): ListMyZitadelPermissionsResponse {
  return { result: [] };
}

export const ListMyZitadelPermissionsResponse = {
  encode(message: ListMyZitadelPermissionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.result) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyZitadelPermissionsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyZitadelPermissionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.result.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyZitadelPermissionsResponse {
    return { result: Array.isArray(object?.result) ? object.result.map((e: any) => String(e)) : [] };
  },

  toJSON(message: ListMyZitadelPermissionsResponse): unknown {
    const obj: any = {};
    if (message.result) {
      obj.result = message.result.map((e) => e);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyZitadelPermissionsResponse>): ListMyZitadelPermissionsResponse {
    return ListMyZitadelPermissionsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyZitadelPermissionsResponse>): ListMyZitadelPermissionsResponse {
    const message = createBaseListMyZitadelPermissionsResponse();
    message.result = object.result?.map((e) => e) || [];
    return message;
  },
};

function createBaseListMyProjectPermissionsRequest(): ListMyProjectPermissionsRequest {
  return {};
}

export const ListMyProjectPermissionsRequest = {
  encode(_: ListMyProjectPermissionsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyProjectPermissionsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyProjectPermissionsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ListMyProjectPermissionsRequest {
    return {};
  },

  toJSON(_: ListMyProjectPermissionsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ListMyProjectPermissionsRequest>): ListMyProjectPermissionsRequest {
    return ListMyProjectPermissionsRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ListMyProjectPermissionsRequest>): ListMyProjectPermissionsRequest {
    const message = createBaseListMyProjectPermissionsRequest();
    return message;
  },
};

function createBaseListMyProjectPermissionsResponse(): ListMyProjectPermissionsResponse {
  return { result: [] };
}

export const ListMyProjectPermissionsResponse = {
  encode(message: ListMyProjectPermissionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.result) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyProjectPermissionsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyProjectPermissionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.result.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyProjectPermissionsResponse {
    return { result: Array.isArray(object?.result) ? object.result.map((e: any) => String(e)) : [] };
  },

  toJSON(message: ListMyProjectPermissionsResponse): unknown {
    const obj: any = {};
    if (message.result) {
      obj.result = message.result.map((e) => e);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyProjectPermissionsResponse>): ListMyProjectPermissionsResponse {
    return ListMyProjectPermissionsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyProjectPermissionsResponse>): ListMyProjectPermissionsResponse {
    const message = createBaseListMyProjectPermissionsResponse();
    message.result = object.result?.map((e) => e) || [];
    return message;
  },
};

function createBaseListMyMembershipsRequest(): ListMyMembershipsRequest {
  return { query: undefined, queries: [] };
}

export const ListMyMembershipsRequest = {
  encode(message: ListMyMembershipsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.queries) {
      MembershipQuery.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyMembershipsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyMembershipsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.query = ListQuery.decode(reader, reader.uint32());
          break;
        case 2:
          message.queries.push(MembershipQuery.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyMembershipsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => MembershipQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyMembershipsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? MembershipQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyMembershipsRequest>): ListMyMembershipsRequest {
    return ListMyMembershipsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyMembershipsRequest>): ListMyMembershipsRequest {
    const message = createBaseListMyMembershipsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.queries = object.queries?.map((e) => MembershipQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListMyMembershipsResponse(): ListMyMembershipsResponse {
  return { details: undefined, result: [] };
}

export const ListMyMembershipsResponse = {
  encode(message: ListMyMembershipsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.result) {
      Membership.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListMyMembershipsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListMyMembershipsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ListDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.result.push(Membership.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListMyMembershipsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => Membership.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListMyMembershipsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.result) {
      obj.result = message.result.map((e) => e ? Membership.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListMyMembershipsResponse>): ListMyMembershipsResponse {
    return ListMyMembershipsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListMyMembershipsResponse>): ListMyMembershipsResponse {
    const message = createBaseListMyMembershipsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.result = object.result?.map((e) => Membership.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetMyLabelPolicyRequest(): GetMyLabelPolicyRequest {
  return {};
}

export const GetMyLabelPolicyRequest = {
  encode(_: GetMyLabelPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyLabelPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyLabelPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyLabelPolicyRequest {
    return {};
  },

  toJSON(_: GetMyLabelPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyLabelPolicyRequest>): GetMyLabelPolicyRequest {
    return GetMyLabelPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyLabelPolicyRequest>): GetMyLabelPolicyRequest {
    const message = createBaseGetMyLabelPolicyRequest();
    return message;
  },
};

function createBaseGetMyLabelPolicyResponse(): GetMyLabelPolicyResponse {
  return { policy: undefined };
}

export const GetMyLabelPolicyResponse = {
  encode(message: GetMyLabelPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      LabelPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyLabelPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyLabelPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = LabelPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyLabelPolicyResponse {
    return { policy: isSet(object.policy) ? LabelPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetMyLabelPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? LabelPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyLabelPolicyResponse>): GetMyLabelPolicyResponse {
    return GetMyLabelPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyLabelPolicyResponse>): GetMyLabelPolicyResponse {
    const message = createBaseGetMyLabelPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? LabelPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseGetMyPrivacyPolicyRequest(): GetMyPrivacyPolicyRequest {
  return {};
}

export const GetMyPrivacyPolicyRequest = {
  encode(_: GetMyPrivacyPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyPrivacyPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyPrivacyPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyPrivacyPolicyRequest {
    return {};
  },

  toJSON(_: GetMyPrivacyPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyPrivacyPolicyRequest>): GetMyPrivacyPolicyRequest {
    return GetMyPrivacyPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyPrivacyPolicyRequest>): GetMyPrivacyPolicyRequest {
    const message = createBaseGetMyPrivacyPolicyRequest();
    return message;
  },
};

function createBaseGetMyPrivacyPolicyResponse(): GetMyPrivacyPolicyResponse {
  return { policy: undefined };
}

export const GetMyPrivacyPolicyResponse = {
  encode(message: GetMyPrivacyPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      PrivacyPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyPrivacyPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyPrivacyPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = PrivacyPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyPrivacyPolicyResponse {
    return { policy: isSet(object.policy) ? PrivacyPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetMyPrivacyPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? PrivacyPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyPrivacyPolicyResponse>): GetMyPrivacyPolicyResponse {
    return GetMyPrivacyPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyPrivacyPolicyResponse>): GetMyPrivacyPolicyResponse {
    const message = createBaseGetMyPrivacyPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? PrivacyPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseGetMyLoginPolicyRequest(): GetMyLoginPolicyRequest {
  return {};
}

export const GetMyLoginPolicyRequest = {
  encode(_: GetMyLoginPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyLoginPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyLoginPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): GetMyLoginPolicyRequest {
    return {};
  },

  toJSON(_: GetMyLoginPolicyRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<GetMyLoginPolicyRequest>): GetMyLoginPolicyRequest {
    return GetMyLoginPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<GetMyLoginPolicyRequest>): GetMyLoginPolicyRequest {
    const message = createBaseGetMyLoginPolicyRequest();
    return message;
  },
};

function createBaseGetMyLoginPolicyResponse(): GetMyLoginPolicyResponse {
  return { policy: undefined };
}

export const GetMyLoginPolicyResponse = {
  encode(message: GetMyLoginPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      LoginPolicy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMyLoginPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMyLoginPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = LoginPolicy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetMyLoginPolicyResponse {
    return { policy: isSet(object.policy) ? LoginPolicy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetMyLoginPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined && (obj.policy = message.policy ? LoginPolicy.toJSON(message.policy) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetMyLoginPolicyResponse>): GetMyLoginPolicyResponse {
    return GetMyLoginPolicyResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetMyLoginPolicyResponse>): GetMyLoginPolicyResponse {
    const message = createBaseGetMyLoginPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? LoginPolicy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

export type AuthServiceDefinition = typeof AuthServiceDefinition;
export const AuthServiceDefinition = {
  name: "AuthService",
  fullName: "zitadel.auth.v1.AuthService",
  methods: {
    healthz: {
      name: "Healthz",
      requestType: HealthzRequest,
      requestStream: false,
      responseType: HealthzResponse,
      responseStream: false,
      options: {},
    },
    getSupportedLanguages: {
      name: "GetSupportedLanguages",
      requestType: GetSupportedLanguagesRequest,
      requestStream: false,
      responseType: GetSupportedLanguagesResponse,
      responseStream: false,
      options: {},
    },
    getMyUser: {
      name: "GetMyUser",
      requestType: GetMyUserRequest,
      requestStream: false,
      responseType: GetMyUserResponse,
      responseStream: false,
      options: {},
    },
    removeMyUser: {
      name: "RemoveMyUser",
      requestType: RemoveMyUserRequest,
      requestStream: false,
      responseType: RemoveMyUserResponse,
      responseStream: false,
      options: {},
    },
    listMyUserChanges: {
      name: "ListMyUserChanges",
      requestType: ListMyUserChangesRequest,
      requestStream: false,
      responseType: ListMyUserChangesResponse,
      responseStream: false,
      options: {},
    },
    listMyUserSessions: {
      name: "ListMyUserSessions",
      requestType: ListMyUserSessionsRequest,
      requestStream: false,
      responseType: ListMyUserSessionsResponse,
      responseStream: false,
      options: {},
    },
    listMyMetadata: {
      name: "ListMyMetadata",
      requestType: ListMyMetadataRequest,
      requestStream: false,
      responseType: ListMyMetadataResponse,
      responseStream: false,
      options: {},
    },
    getMyMetadata: {
      name: "GetMyMetadata",
      requestType: GetMyMetadataRequest,
      requestStream: false,
      responseType: GetMyMetadataResponse,
      responseStream: false,
      options: {},
    },
    listMyRefreshTokens: {
      name: "ListMyRefreshTokens",
      requestType: ListMyRefreshTokensRequest,
      requestStream: false,
      responseType: ListMyRefreshTokensResponse,
      responseStream: false,
      options: {},
    },
    revokeMyRefreshToken: {
      name: "RevokeMyRefreshToken",
      requestType: RevokeMyRefreshTokenRequest,
      requestStream: false,
      responseType: RevokeMyRefreshTokenResponse,
      responseStream: false,
      options: {},
    },
    revokeAllMyRefreshTokens: {
      name: "RevokeAllMyRefreshTokens",
      requestType: RevokeAllMyRefreshTokensRequest,
      requestStream: false,
      responseType: RevokeAllMyRefreshTokensResponse,
      responseStream: false,
      options: {},
    },
    updateMyUserName: {
      name: "UpdateMyUserName",
      requestType: UpdateMyUserNameRequest,
      requestStream: false,
      responseType: UpdateMyUserNameResponse,
      responseStream: false,
      options: {},
    },
    getMyPasswordComplexityPolicy: {
      name: "GetMyPasswordComplexityPolicy",
      requestType: GetMyPasswordComplexityPolicyRequest,
      requestStream: false,
      responseType: GetMyPasswordComplexityPolicyResponse,
      responseStream: false,
      options: {},
    },
    updateMyPassword: {
      name: "UpdateMyPassword",
      requestType: UpdateMyPasswordRequest,
      requestStream: false,
      responseType: UpdateMyPasswordResponse,
      responseStream: false,
      options: {},
    },
    getMyProfile: {
      name: "GetMyProfile",
      requestType: GetMyProfileRequest,
      requestStream: false,
      responseType: GetMyProfileResponse,
      responseStream: false,
      options: {},
    },
    updateMyProfile: {
      name: "UpdateMyProfile",
      requestType: UpdateMyProfileRequest,
      requestStream: false,
      responseType: UpdateMyProfileResponse,
      responseStream: false,
      options: {},
    },
    getMyEmail: {
      name: "GetMyEmail",
      requestType: GetMyEmailRequest,
      requestStream: false,
      responseType: GetMyEmailResponse,
      responseStream: false,
      options: {},
    },
    setMyEmail: {
      name: "SetMyEmail",
      requestType: SetMyEmailRequest,
      requestStream: false,
      responseType: SetMyEmailResponse,
      responseStream: false,
      options: {},
    },
    verifyMyEmail: {
      name: "VerifyMyEmail",
      requestType: VerifyMyEmailRequest,
      requestStream: false,
      responseType: VerifyMyEmailResponse,
      responseStream: false,
      options: {},
    },
    resendMyEmailVerification: {
      name: "ResendMyEmailVerification",
      requestType: ResendMyEmailVerificationRequest,
      requestStream: false,
      responseType: ResendMyEmailVerificationResponse,
      responseStream: false,
      options: {},
    },
    getMyPhone: {
      name: "GetMyPhone",
      requestType: GetMyPhoneRequest,
      requestStream: false,
      responseType: GetMyPhoneResponse,
      responseStream: false,
      options: {},
    },
    setMyPhone: {
      name: "SetMyPhone",
      requestType: SetMyPhoneRequest,
      requestStream: false,
      responseType: SetMyPhoneResponse,
      responseStream: false,
      options: {},
    },
    verifyMyPhone: {
      name: "VerifyMyPhone",
      requestType: VerifyMyPhoneRequest,
      requestStream: false,
      responseType: VerifyMyPhoneResponse,
      responseStream: false,
      options: {},
    },
    /** Resends an sms to the last given phone number, to verify it */
    resendMyPhoneVerification: {
      name: "ResendMyPhoneVerification",
      requestType: ResendMyPhoneVerificationRequest,
      requestStream: false,
      responseType: ResendMyPhoneVerificationResponse,
      responseStream: false,
      options: {},
    },
    removeMyPhone: {
      name: "RemoveMyPhone",
      requestType: RemoveMyPhoneRequest,
      requestStream: false,
      responseType: RemoveMyPhoneResponse,
      responseStream: false,
      options: {},
    },
    removeMyAvatar: {
      name: "RemoveMyAvatar",
      requestType: RemoveMyAvatarRequest,
      requestStream: false,
      responseType: RemoveMyAvatarResponse,
      responseStream: false,
      options: {},
    },
    listMyLinkedIDPs: {
      name: "ListMyLinkedIDPs",
      requestType: ListMyLinkedIDPsRequest,
      requestStream: false,
      responseType: ListMyLinkedIDPsResponse,
      responseStream: false,
      options: {},
    },
    removeMyLinkedIDP: {
      name: "RemoveMyLinkedIDP",
      requestType: RemoveMyLinkedIDPRequest,
      requestStream: false,
      responseType: RemoveMyLinkedIDPResponse,
      responseStream: false,
      options: {},
    },
    listMyAuthFactors: {
      name: "ListMyAuthFactors",
      requestType: ListMyAuthFactorsRequest,
      requestStream: false,
      responseType: ListMyAuthFactorsResponse,
      responseStream: false,
      options: {},
    },
    addMyAuthFactorOTP: {
      name: "AddMyAuthFactorOTP",
      requestType: AddMyAuthFactorOTPRequest,
      requestStream: false,
      responseType: AddMyAuthFactorOTPResponse,
      responseStream: false,
      options: {},
    },
    verifyMyAuthFactorOTP: {
      name: "VerifyMyAuthFactorOTP",
      requestType: VerifyMyAuthFactorOTPRequest,
      requestStream: false,
      responseType: VerifyMyAuthFactorOTPResponse,
      responseStream: false,
      options: {},
    },
    removeMyAuthFactorOTP: {
      name: "RemoveMyAuthFactorOTP",
      requestType: RemoveMyAuthFactorOTPRequest,
      requestStream: false,
      responseType: RemoveMyAuthFactorOTPResponse,
      responseStream: false,
      options: {},
    },
    addMyAuthFactorU2F: {
      name: "AddMyAuthFactorU2F",
      requestType: AddMyAuthFactorU2FRequest,
      requestStream: false,
      responseType: AddMyAuthFactorU2FResponse,
      responseStream: false,
      options: {},
    },
    verifyMyAuthFactorU2F: {
      name: "VerifyMyAuthFactorU2F",
      requestType: VerifyMyAuthFactorU2FRequest,
      requestStream: false,
      responseType: VerifyMyAuthFactorU2FResponse,
      responseStream: false,
      options: {},
    },
    removeMyAuthFactorU2F: {
      name: "RemoveMyAuthFactorU2F",
      requestType: RemoveMyAuthFactorU2FRequest,
      requestStream: false,
      responseType: RemoveMyAuthFactorU2FResponse,
      responseStream: false,
      options: {},
    },
    listMyPasswordless: {
      name: "ListMyPasswordless",
      requestType: ListMyPasswordlessRequest,
      requestStream: false,
      responseType: ListMyPasswordlessResponse,
      responseStream: false,
      options: {},
    },
    addMyPasswordless: {
      name: "AddMyPasswordless",
      requestType: AddMyPasswordlessRequest,
      requestStream: false,
      responseType: AddMyPasswordlessResponse,
      responseStream: false,
      options: {},
    },
    addMyPasswordlessLink: {
      name: "AddMyPasswordlessLink",
      requestType: AddMyPasswordlessLinkRequest,
      requestStream: false,
      responseType: AddMyPasswordlessLinkResponse,
      responseStream: false,
      options: {},
    },
    sendMyPasswordlessLink: {
      name: "SendMyPasswordlessLink",
      requestType: SendMyPasswordlessLinkRequest,
      requestStream: false,
      responseType: SendMyPasswordlessLinkResponse,
      responseStream: false,
      options: {},
    },
    verifyMyPasswordless: {
      name: "VerifyMyPasswordless",
      requestType: VerifyMyPasswordlessRequest,
      requestStream: false,
      responseType: VerifyMyPasswordlessResponse,
      responseStream: false,
      options: {},
    },
    removeMyPasswordless: {
      name: "RemoveMyPasswordless",
      requestType: RemoveMyPasswordlessRequest,
      requestStream: false,
      responseType: RemoveMyPasswordlessResponse,
      responseStream: false,
      options: {},
    },
    listMyUserGrants: {
      name: "ListMyUserGrants",
      requestType: ListMyUserGrantsRequest,
      requestStream: false,
      responseType: ListMyUserGrantsResponse,
      responseStream: false,
      options: {},
    },
    listMyProjectOrgs: {
      name: "ListMyProjectOrgs",
      requestType: ListMyProjectOrgsRequest,
      requestStream: false,
      responseType: ListMyProjectOrgsResponse,
      responseStream: false,
      options: {},
    },
    listMyZitadelPermissions: {
      name: "ListMyZitadelPermissions",
      requestType: ListMyZitadelPermissionsRequest,
      requestStream: false,
      responseType: ListMyZitadelPermissionsResponse,
      responseStream: false,
      options: {},
    },
    listMyProjectPermissions: {
      name: "ListMyProjectPermissions",
      requestType: ListMyProjectPermissionsRequest,
      requestStream: false,
      responseType: ListMyProjectPermissionsResponse,
      responseStream: false,
      options: {},
    },
    listMyMemberships: {
      name: "ListMyMemberships",
      requestType: ListMyMembershipsRequest,
      requestStream: false,
      responseType: ListMyMembershipsResponse,
      responseStream: false,
      options: {},
    },
    getMyLabelPolicy: {
      name: "GetMyLabelPolicy",
      requestType: GetMyLabelPolicyRequest,
      requestStream: false,
      responseType: GetMyLabelPolicyResponse,
      responseStream: false,
      options: {},
    },
    getMyPrivacyPolicy: {
      name: "GetMyPrivacyPolicy",
      requestType: GetMyPrivacyPolicyRequest,
      requestStream: false,
      responseType: GetMyPrivacyPolicyResponse,
      responseStream: false,
      options: {},
    },
    getMyLoginPolicy: {
      name: "GetMyLoginPolicy",
      requestType: GetMyLoginPolicyRequest,
      requestStream: false,
      responseType: GetMyLoginPolicyResponse,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface AuthServiceImplementation<CallContextExt = {}> {
  healthz(request: HealthzRequest, context: CallContext & CallContextExt): Promise<DeepPartial<HealthzResponse>>;
  getSupportedLanguages(
    request: GetSupportedLanguagesRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetSupportedLanguagesResponse>>;
  getMyUser(request: GetMyUserRequest, context: CallContext & CallContextExt): Promise<DeepPartial<GetMyUserResponse>>;
  removeMyUser(
    request: RemoveMyUserRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveMyUserResponse>>;
  listMyUserChanges(
    request: ListMyUserChangesRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyUserChangesResponse>>;
  listMyUserSessions(
    request: ListMyUserSessionsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyUserSessionsResponse>>;
  listMyMetadata(
    request: ListMyMetadataRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyMetadataResponse>>;
  getMyMetadata(
    request: GetMyMetadataRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyMetadataResponse>>;
  listMyRefreshTokens(
    request: ListMyRefreshTokensRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyRefreshTokensResponse>>;
  revokeMyRefreshToken(
    request: RevokeMyRefreshTokenRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RevokeMyRefreshTokenResponse>>;
  revokeAllMyRefreshTokens(
    request: RevokeAllMyRefreshTokensRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RevokeAllMyRefreshTokensResponse>>;
  updateMyUserName(
    request: UpdateMyUserNameRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateMyUserNameResponse>>;
  getMyPasswordComplexityPolicy(
    request: GetMyPasswordComplexityPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyPasswordComplexityPolicyResponse>>;
  updateMyPassword(
    request: UpdateMyPasswordRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateMyPasswordResponse>>;
  getMyProfile(
    request: GetMyProfileRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyProfileResponse>>;
  updateMyProfile(
    request: UpdateMyProfileRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateMyProfileResponse>>;
  getMyEmail(
    request: GetMyEmailRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyEmailResponse>>;
  setMyEmail(
    request: SetMyEmailRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetMyEmailResponse>>;
  verifyMyEmail(
    request: VerifyMyEmailRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyMyEmailResponse>>;
  resendMyEmailVerification(
    request: ResendMyEmailVerificationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResendMyEmailVerificationResponse>>;
  getMyPhone(
    request: GetMyPhoneRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyPhoneResponse>>;
  setMyPhone(
    request: SetMyPhoneRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetMyPhoneResponse>>;
  verifyMyPhone(
    request: VerifyMyPhoneRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyMyPhoneResponse>>;
  /** Resends an sms to the last given phone number, to verify it */
  resendMyPhoneVerification(
    request: ResendMyPhoneVerificationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResendMyPhoneVerificationResponse>>;
  removeMyPhone(
    request: RemoveMyPhoneRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveMyPhoneResponse>>;
  removeMyAvatar(
    request: RemoveMyAvatarRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveMyAvatarResponse>>;
  listMyLinkedIDPs(
    request: ListMyLinkedIDPsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyLinkedIDPsResponse>>;
  removeMyLinkedIDP(
    request: RemoveMyLinkedIDPRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveMyLinkedIDPResponse>>;
  listMyAuthFactors(
    request: ListMyAuthFactorsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyAuthFactorsResponse>>;
  addMyAuthFactorOTP(
    request: AddMyAuthFactorOTPRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddMyAuthFactorOTPResponse>>;
  verifyMyAuthFactorOTP(
    request: VerifyMyAuthFactorOTPRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyMyAuthFactorOTPResponse>>;
  removeMyAuthFactorOTP(
    request: RemoveMyAuthFactorOTPRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveMyAuthFactorOTPResponse>>;
  addMyAuthFactorU2F(
    request: AddMyAuthFactorU2FRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddMyAuthFactorU2FResponse>>;
  verifyMyAuthFactorU2F(
    request: VerifyMyAuthFactorU2FRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyMyAuthFactorU2FResponse>>;
  removeMyAuthFactorU2F(
    request: RemoveMyAuthFactorU2FRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveMyAuthFactorU2FResponse>>;
  listMyPasswordless(
    request: ListMyPasswordlessRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyPasswordlessResponse>>;
  addMyPasswordless(
    request: AddMyPasswordlessRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddMyPasswordlessResponse>>;
  addMyPasswordlessLink(
    request: AddMyPasswordlessLinkRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddMyPasswordlessLinkResponse>>;
  sendMyPasswordlessLink(
    request: SendMyPasswordlessLinkRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SendMyPasswordlessLinkResponse>>;
  verifyMyPasswordless(
    request: VerifyMyPasswordlessRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyMyPasswordlessResponse>>;
  removeMyPasswordless(
    request: RemoveMyPasswordlessRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveMyPasswordlessResponse>>;
  listMyUserGrants(
    request: ListMyUserGrantsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyUserGrantsResponse>>;
  listMyProjectOrgs(
    request: ListMyProjectOrgsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyProjectOrgsResponse>>;
  listMyZitadelPermissions(
    request: ListMyZitadelPermissionsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyZitadelPermissionsResponse>>;
  listMyProjectPermissions(
    request: ListMyProjectPermissionsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyProjectPermissionsResponse>>;
  listMyMemberships(
    request: ListMyMembershipsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListMyMembershipsResponse>>;
  getMyLabelPolicy(
    request: GetMyLabelPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyLabelPolicyResponse>>;
  getMyPrivacyPolicy(
    request: GetMyPrivacyPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyPrivacyPolicyResponse>>;
  getMyLoginPolicy(
    request: GetMyLoginPolicyRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetMyLoginPolicyResponse>>;
}

export interface AuthServiceClient<CallOptionsExt = {}> {
  healthz(request: DeepPartial<HealthzRequest>, options?: CallOptions & CallOptionsExt): Promise<HealthzResponse>;
  getSupportedLanguages(
    request: DeepPartial<GetSupportedLanguagesRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetSupportedLanguagesResponse>;
  getMyUser(request: DeepPartial<GetMyUserRequest>, options?: CallOptions & CallOptionsExt): Promise<GetMyUserResponse>;
  removeMyUser(
    request: DeepPartial<RemoveMyUserRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveMyUserResponse>;
  listMyUserChanges(
    request: DeepPartial<ListMyUserChangesRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyUserChangesResponse>;
  listMyUserSessions(
    request: DeepPartial<ListMyUserSessionsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyUserSessionsResponse>;
  listMyMetadata(
    request: DeepPartial<ListMyMetadataRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyMetadataResponse>;
  getMyMetadata(
    request: DeepPartial<GetMyMetadataRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyMetadataResponse>;
  listMyRefreshTokens(
    request: DeepPartial<ListMyRefreshTokensRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyRefreshTokensResponse>;
  revokeMyRefreshToken(
    request: DeepPartial<RevokeMyRefreshTokenRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RevokeMyRefreshTokenResponse>;
  revokeAllMyRefreshTokens(
    request: DeepPartial<RevokeAllMyRefreshTokensRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RevokeAllMyRefreshTokensResponse>;
  updateMyUserName(
    request: DeepPartial<UpdateMyUserNameRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateMyUserNameResponse>;
  getMyPasswordComplexityPolicy(
    request: DeepPartial<GetMyPasswordComplexityPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyPasswordComplexityPolicyResponse>;
  updateMyPassword(
    request: DeepPartial<UpdateMyPasswordRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateMyPasswordResponse>;
  getMyProfile(
    request: DeepPartial<GetMyProfileRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyProfileResponse>;
  updateMyProfile(
    request: DeepPartial<UpdateMyProfileRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateMyProfileResponse>;
  getMyEmail(
    request: DeepPartial<GetMyEmailRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyEmailResponse>;
  setMyEmail(
    request: DeepPartial<SetMyEmailRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetMyEmailResponse>;
  verifyMyEmail(
    request: DeepPartial<VerifyMyEmailRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyMyEmailResponse>;
  resendMyEmailVerification(
    request: DeepPartial<ResendMyEmailVerificationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResendMyEmailVerificationResponse>;
  getMyPhone(
    request: DeepPartial<GetMyPhoneRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyPhoneResponse>;
  setMyPhone(
    request: DeepPartial<SetMyPhoneRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetMyPhoneResponse>;
  verifyMyPhone(
    request: DeepPartial<VerifyMyPhoneRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyMyPhoneResponse>;
  /** Resends an sms to the last given phone number, to verify it */
  resendMyPhoneVerification(
    request: DeepPartial<ResendMyPhoneVerificationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResendMyPhoneVerificationResponse>;
  removeMyPhone(
    request: DeepPartial<RemoveMyPhoneRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveMyPhoneResponse>;
  removeMyAvatar(
    request: DeepPartial<RemoveMyAvatarRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveMyAvatarResponse>;
  listMyLinkedIDPs(
    request: DeepPartial<ListMyLinkedIDPsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyLinkedIDPsResponse>;
  removeMyLinkedIDP(
    request: DeepPartial<RemoveMyLinkedIDPRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveMyLinkedIDPResponse>;
  listMyAuthFactors(
    request: DeepPartial<ListMyAuthFactorsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyAuthFactorsResponse>;
  addMyAuthFactorOTP(
    request: DeepPartial<AddMyAuthFactorOTPRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddMyAuthFactorOTPResponse>;
  verifyMyAuthFactorOTP(
    request: DeepPartial<VerifyMyAuthFactorOTPRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyMyAuthFactorOTPResponse>;
  removeMyAuthFactorOTP(
    request: DeepPartial<RemoveMyAuthFactorOTPRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveMyAuthFactorOTPResponse>;
  addMyAuthFactorU2F(
    request: DeepPartial<AddMyAuthFactorU2FRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddMyAuthFactorU2FResponse>;
  verifyMyAuthFactorU2F(
    request: DeepPartial<VerifyMyAuthFactorU2FRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyMyAuthFactorU2FResponse>;
  removeMyAuthFactorU2F(
    request: DeepPartial<RemoveMyAuthFactorU2FRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveMyAuthFactorU2FResponse>;
  listMyPasswordless(
    request: DeepPartial<ListMyPasswordlessRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyPasswordlessResponse>;
  addMyPasswordless(
    request: DeepPartial<AddMyPasswordlessRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddMyPasswordlessResponse>;
  addMyPasswordlessLink(
    request: DeepPartial<AddMyPasswordlessLinkRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddMyPasswordlessLinkResponse>;
  sendMyPasswordlessLink(
    request: DeepPartial<SendMyPasswordlessLinkRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SendMyPasswordlessLinkResponse>;
  verifyMyPasswordless(
    request: DeepPartial<VerifyMyPasswordlessRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyMyPasswordlessResponse>;
  removeMyPasswordless(
    request: DeepPartial<RemoveMyPasswordlessRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveMyPasswordlessResponse>;
  listMyUserGrants(
    request: DeepPartial<ListMyUserGrantsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyUserGrantsResponse>;
  listMyProjectOrgs(
    request: DeepPartial<ListMyProjectOrgsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyProjectOrgsResponse>;
  listMyZitadelPermissions(
    request: DeepPartial<ListMyZitadelPermissionsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyZitadelPermissionsResponse>;
  listMyProjectPermissions(
    request: DeepPartial<ListMyProjectPermissionsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyProjectPermissionsResponse>;
  listMyMemberships(
    request: DeepPartial<ListMyMembershipsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListMyMembershipsResponse>;
  getMyLabelPolicy(
    request: DeepPartial<GetMyLabelPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyLabelPolicyResponse>;
  getMyPrivacyPolicy(
    request: DeepPartial<GetMyPrivacyPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyPrivacyPolicyResponse>;
  getMyLoginPolicy(
    request: DeepPartial<GetMyLoginPolicyRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetMyLoginPolicyResponse>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

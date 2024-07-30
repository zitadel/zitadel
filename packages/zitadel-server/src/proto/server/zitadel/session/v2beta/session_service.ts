/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Duration } from "../../../google/protobuf/duration";
import { Struct } from "../../../google/protobuf/struct";
import { Details, ListDetails, ListQuery } from "../../object/v2beta/object";
import { Challenges, RequestChallenges } from "./challenge";
import {
  SearchQuery,
  Session,
  SessionFieldName,
  sessionFieldNameFromJSON,
  sessionFieldNameToJSON,
  UserAgent,
} from "./session";

export const protobufPackage = "zitadel.session.v2beta";

export interface ListSessionsRequest {
  query: ListQuery | undefined;
  queries: SearchQuery[];
  sortingColumn: SessionFieldName;
}

export interface ListSessionsResponse {
  details: ListDetails | undefined;
  sessions: Session[];
}

export interface GetSessionRequest {
  sessionId: string;
  sessionToken?: string | undefined;
}

export interface GetSessionResponse {
  session: Session | undefined;
}

export interface CreateSessionRequest {
  checks: Checks | undefined;
  metadata: { [key: string]: Buffer };
  challenges: RequestChallenges | undefined;
  userAgent: UserAgent | undefined;
  lifetime?: Duration | undefined;
}

export interface CreateSessionRequest_MetadataEntry {
  key: string;
  value: Buffer;
}

export interface CreateSessionResponse {
  details: Details | undefined;
  sessionId: string;
  sessionToken: string;
  challenges: Challenges | undefined;
}

export interface SetSessionRequest {
  sessionId: string;
  sessionToken: string;
  checks: Checks | undefined;
  metadata: { [key: string]: Buffer };
  challenges: RequestChallenges | undefined;
  lifetime?: Duration | undefined;
}

export interface SetSessionRequest_MetadataEntry {
  key: string;
  value: Buffer;
}

export interface SetSessionResponse {
  details: Details | undefined;
  sessionToken: string;
  challenges: Challenges | undefined;
}

export interface DeleteSessionRequest {
  sessionId: string;
  sessionToken?: string | undefined;
}

export interface DeleteSessionResponse {
  details: Details | undefined;
}

export interface Checks {
  user?: CheckUser | undefined;
  password?: CheckPassword | undefined;
  webAuthN?: CheckWebAuthN | undefined;
  idpIntent?: CheckIDPIntent | undefined;
  totp?: CheckTOTP | undefined;
  otpSms?: CheckOTP | undefined;
  otpEmail?: CheckOTP | undefined;
}

export interface CheckUser {
  userId?: string | undefined;
  loginName?: string | undefined;
}

export interface CheckPassword {
  password: string;
}

export interface CheckWebAuthN {
  credentialAssertionData: { [key: string]: any } | undefined;
}

export interface CheckIDPIntent {
  idpIntentId: string;
  idpIntentToken: string;
}

export interface CheckTOTP {
  code: string;
}

export interface CheckOTP {
  code: string;
}

function createBaseListSessionsRequest(): ListSessionsRequest {
  return { query: undefined, queries: [], sortingColumn: 0 };
}

export const ListSessionsRequest = {
  encode(message: ListSessionsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.queries) {
      SearchQuery.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(24).int32(message.sortingColumn);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListSessionsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListSessionsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.query = ListQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.queries.push(SearchQuery.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.sortingColumn = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListSessionsRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => SearchQuery.fromJSON(e)) : [],
      sortingColumn: isSet(object.sortingColumn) ? sessionFieldNameFromJSON(object.sortingColumn) : 0,
    };
  },

  toJSON(message: ListSessionsRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? SearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    message.sortingColumn !== undefined && (obj.sortingColumn = sessionFieldNameToJSON(message.sortingColumn));
    return obj;
  },

  create(base?: DeepPartial<ListSessionsRequest>): ListSessionsRequest {
    return ListSessionsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListSessionsRequest>): ListSessionsRequest {
    const message = createBaseListSessionsRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.queries = object.queries?.map((e) => SearchQuery.fromPartial(e)) || [];
    message.sortingColumn = object.sortingColumn ?? 0;
    return message;
  },
};

function createBaseListSessionsResponse(): ListSessionsResponse {
  return { details: undefined, sessions: [] };
}

export const ListSessionsResponse = {
  encode(message: ListSessionsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.sessions) {
      Session.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListSessionsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListSessionsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = ListDetails.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sessions.push(Session.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListSessionsResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      sessions: Array.isArray(object?.sessions) ? object.sessions.map((e: any) => Session.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListSessionsResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    if (message.sessions) {
      obj.sessions = message.sessions.map((e) => e ? Session.toJSON(e) : undefined);
    } else {
      obj.sessions = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListSessionsResponse>): ListSessionsResponse {
    return ListSessionsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListSessionsResponse>): ListSessionsResponse {
    const message = createBaseListSessionsResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.sessions = object.sessions?.map((e) => Session.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetSessionRequest(): GetSessionRequest {
  return { sessionId: "", sessionToken: undefined };
}

export const GetSessionRequest = {
  encode(message: GetSessionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sessionId !== "") {
      writer.uint32(10).string(message.sessionId);
    }
    if (message.sessionToken !== undefined) {
      writer.uint32(18).string(message.sessionToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSessionRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSessionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.sessionId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sessionToken = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetSessionRequest {
    return {
      sessionId: isSet(object.sessionId) ? String(object.sessionId) : "",
      sessionToken: isSet(object.sessionToken) ? String(object.sessionToken) : undefined,
    };
  },

  toJSON(message: GetSessionRequest): unknown {
    const obj: any = {};
    message.sessionId !== undefined && (obj.sessionId = message.sessionId);
    message.sessionToken !== undefined && (obj.sessionToken = message.sessionToken);
    return obj;
  },

  create(base?: DeepPartial<GetSessionRequest>): GetSessionRequest {
    return GetSessionRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSessionRequest>): GetSessionRequest {
    const message = createBaseGetSessionRequest();
    message.sessionId = object.sessionId ?? "";
    message.sessionToken = object.sessionToken ?? undefined;
    return message;
  },
};

function createBaseGetSessionResponse(): GetSessionResponse {
  return { session: undefined };
}

export const GetSessionResponse = {
  encode(message: GetSessionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.session !== undefined) {
      Session.encode(message.session, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSessionResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSessionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.session = Session.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetSessionResponse {
    return { session: isSet(object.session) ? Session.fromJSON(object.session) : undefined };
  },

  toJSON(message: GetSessionResponse): unknown {
    const obj: any = {};
    message.session !== undefined && (obj.session = message.session ? Session.toJSON(message.session) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetSessionResponse>): GetSessionResponse {
    return GetSessionResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetSessionResponse>): GetSessionResponse {
    const message = createBaseGetSessionResponse();
    message.session = (object.session !== undefined && object.session !== null)
      ? Session.fromPartial(object.session)
      : undefined;
    return message;
  },
};

function createBaseCreateSessionRequest(): CreateSessionRequest {
  return { checks: undefined, metadata: {}, challenges: undefined, userAgent: undefined, lifetime: undefined };
}

export const CreateSessionRequest = {
  encode(message: CreateSessionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.checks !== undefined) {
      Checks.encode(message.checks, writer.uint32(10).fork()).ldelim();
    }
    Object.entries(message.metadata).forEach(([key, value]) => {
      CreateSessionRequest_MetadataEntry.encode({ key: key as any, value }, writer.uint32(18).fork()).ldelim();
    });
    if (message.challenges !== undefined) {
      RequestChallenges.encode(message.challenges, writer.uint32(26).fork()).ldelim();
    }
    if (message.userAgent !== undefined) {
      UserAgent.encode(message.userAgent, writer.uint32(34).fork()).ldelim();
    }
    if (message.lifetime !== undefined) {
      Duration.encode(message.lifetime, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateSessionRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateSessionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.checks = Checks.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          const entry2 = CreateSessionRequest_MetadataEntry.decode(reader, reader.uint32());
          if (entry2.value !== undefined) {
            message.metadata[entry2.key] = entry2.value;
          }
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.challenges = RequestChallenges.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.userAgent = UserAgent.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.lifetime = Duration.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateSessionRequest {
    return {
      checks: isSet(object.checks) ? Checks.fromJSON(object.checks) : undefined,
      metadata: isObject(object.metadata)
        ? Object.entries(object.metadata).reduce<{ [key: string]: Buffer }>((acc, [key, value]) => {
          acc[key] = Buffer.from(bytesFromBase64(value as string));
          return acc;
        }, {})
        : {},
      challenges: isSet(object.challenges) ? RequestChallenges.fromJSON(object.challenges) : undefined,
      userAgent: isSet(object.userAgent) ? UserAgent.fromJSON(object.userAgent) : undefined,
      lifetime: isSet(object.lifetime) ? Duration.fromJSON(object.lifetime) : undefined,
    };
  },

  toJSON(message: CreateSessionRequest): unknown {
    const obj: any = {};
    message.checks !== undefined && (obj.checks = message.checks ? Checks.toJSON(message.checks) : undefined);
    obj.metadata = {};
    if (message.metadata) {
      Object.entries(message.metadata).forEach(([k, v]) => {
        obj.metadata[k] = base64FromBytes(v);
      });
    }
    message.challenges !== undefined &&
      (obj.challenges = message.challenges ? RequestChallenges.toJSON(message.challenges) : undefined);
    message.userAgent !== undefined &&
      (obj.userAgent = message.userAgent ? UserAgent.toJSON(message.userAgent) : undefined);
    message.lifetime !== undefined && (obj.lifetime = message.lifetime ? Duration.toJSON(message.lifetime) : undefined);
    return obj;
  },

  create(base?: DeepPartial<CreateSessionRequest>): CreateSessionRequest {
    return CreateSessionRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateSessionRequest>): CreateSessionRequest {
    const message = createBaseCreateSessionRequest();
    message.checks = (object.checks !== undefined && object.checks !== null)
      ? Checks.fromPartial(object.checks)
      : undefined;
    message.metadata = Object.entries(object.metadata ?? {}).reduce<{ [key: string]: Buffer }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = value;
      }
      return acc;
    }, {});
    message.challenges = (object.challenges !== undefined && object.challenges !== null)
      ? RequestChallenges.fromPartial(object.challenges)
      : undefined;
    message.userAgent = (object.userAgent !== undefined && object.userAgent !== null)
      ? UserAgent.fromPartial(object.userAgent)
      : undefined;
    message.lifetime = (object.lifetime !== undefined && object.lifetime !== null)
      ? Duration.fromPartial(object.lifetime)
      : undefined;
    return message;
  },
};

function createBaseCreateSessionRequest_MetadataEntry(): CreateSessionRequest_MetadataEntry {
  return { key: "", value: Buffer.alloc(0) };
}

export const CreateSessionRequest_MetadataEntry = {
  encode(message: CreateSessionRequest_MetadataEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value.length !== 0) {
      writer.uint32(18).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateSessionRequest_MetadataEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateSessionRequest_MetadataEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = reader.bytes() as Buffer;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateSessionRequest_MetadataEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Buffer.from(bytesFromBase64(object.value)) : Buffer.alloc(0),
    };
  },

  toJSON(message: CreateSessionRequest_MetadataEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = base64FromBytes(message.value !== undefined ? message.value : Buffer.alloc(0)));
    return obj;
  },

  create(base?: DeepPartial<CreateSessionRequest_MetadataEntry>): CreateSessionRequest_MetadataEntry {
    return CreateSessionRequest_MetadataEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateSessionRequest_MetadataEntry>): CreateSessionRequest_MetadataEntry {
    const message = createBaseCreateSessionRequest_MetadataEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseCreateSessionResponse(): CreateSessionResponse {
  return { details: undefined, sessionId: "", sessionToken: "", challenges: undefined };
}

export const CreateSessionResponse = {
  encode(message: CreateSessionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.sessionId !== "") {
      writer.uint32(18).string(message.sessionId);
    }
    if (message.sessionToken !== "") {
      writer.uint32(26).string(message.sessionToken);
    }
    if (message.challenges !== undefined) {
      Challenges.encode(message.challenges, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateSessionResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateSessionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sessionId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.sessionToken = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.challenges = Challenges.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateSessionResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      sessionId: isSet(object.sessionId) ? String(object.sessionId) : "",
      sessionToken: isSet(object.sessionToken) ? String(object.sessionToken) : "",
      challenges: isSet(object.challenges) ? Challenges.fromJSON(object.challenges) : undefined,
    };
  },

  toJSON(message: CreateSessionResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.sessionId !== undefined && (obj.sessionId = message.sessionId);
    message.sessionToken !== undefined && (obj.sessionToken = message.sessionToken);
    message.challenges !== undefined &&
      (obj.challenges = message.challenges ? Challenges.toJSON(message.challenges) : undefined);
    return obj;
  },

  create(base?: DeepPartial<CreateSessionResponse>): CreateSessionResponse {
    return CreateSessionResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateSessionResponse>): CreateSessionResponse {
    const message = createBaseCreateSessionResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.sessionId = object.sessionId ?? "";
    message.sessionToken = object.sessionToken ?? "";
    message.challenges = (object.challenges !== undefined && object.challenges !== null)
      ? Challenges.fromPartial(object.challenges)
      : undefined;
    return message;
  },
};

function createBaseSetSessionRequest(): SetSessionRequest {
  return {
    sessionId: "",
    sessionToken: "",
    checks: undefined,
    metadata: {},
    challenges: undefined,
    lifetime: undefined,
  };
}

export const SetSessionRequest = {
  encode(message: SetSessionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sessionId !== "") {
      writer.uint32(10).string(message.sessionId);
    }
    if (message.sessionToken !== "") {
      writer.uint32(18).string(message.sessionToken);
    }
    if (message.checks !== undefined) {
      Checks.encode(message.checks, writer.uint32(26).fork()).ldelim();
    }
    Object.entries(message.metadata).forEach(([key, value]) => {
      SetSessionRequest_MetadataEntry.encode({ key: key as any, value }, writer.uint32(34).fork()).ldelim();
    });
    if (message.challenges !== undefined) {
      RequestChallenges.encode(message.challenges, writer.uint32(42).fork()).ldelim();
    }
    if (message.lifetime !== undefined) {
      Duration.encode(message.lifetime, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetSessionRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetSessionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.sessionId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sessionToken = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.checks = Checks.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          const entry4 = SetSessionRequest_MetadataEntry.decode(reader, reader.uint32());
          if (entry4.value !== undefined) {
            message.metadata[entry4.key] = entry4.value;
          }
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.challenges = RequestChallenges.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.lifetime = Duration.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetSessionRequest {
    return {
      sessionId: isSet(object.sessionId) ? String(object.sessionId) : "",
      sessionToken: isSet(object.sessionToken) ? String(object.sessionToken) : "",
      checks: isSet(object.checks) ? Checks.fromJSON(object.checks) : undefined,
      metadata: isObject(object.metadata)
        ? Object.entries(object.metadata).reduce<{ [key: string]: Buffer }>((acc, [key, value]) => {
          acc[key] = Buffer.from(bytesFromBase64(value as string));
          return acc;
        }, {})
        : {},
      challenges: isSet(object.challenges) ? RequestChallenges.fromJSON(object.challenges) : undefined,
      lifetime: isSet(object.lifetime) ? Duration.fromJSON(object.lifetime) : undefined,
    };
  },

  toJSON(message: SetSessionRequest): unknown {
    const obj: any = {};
    message.sessionId !== undefined && (obj.sessionId = message.sessionId);
    message.sessionToken !== undefined && (obj.sessionToken = message.sessionToken);
    message.checks !== undefined && (obj.checks = message.checks ? Checks.toJSON(message.checks) : undefined);
    obj.metadata = {};
    if (message.metadata) {
      Object.entries(message.metadata).forEach(([k, v]) => {
        obj.metadata[k] = base64FromBytes(v);
      });
    }
    message.challenges !== undefined &&
      (obj.challenges = message.challenges ? RequestChallenges.toJSON(message.challenges) : undefined);
    message.lifetime !== undefined && (obj.lifetime = message.lifetime ? Duration.toJSON(message.lifetime) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetSessionRequest>): SetSessionRequest {
    return SetSessionRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetSessionRequest>): SetSessionRequest {
    const message = createBaseSetSessionRequest();
    message.sessionId = object.sessionId ?? "";
    message.sessionToken = object.sessionToken ?? "";
    message.checks = (object.checks !== undefined && object.checks !== null)
      ? Checks.fromPartial(object.checks)
      : undefined;
    message.metadata = Object.entries(object.metadata ?? {}).reduce<{ [key: string]: Buffer }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[key] = value;
      }
      return acc;
    }, {});
    message.challenges = (object.challenges !== undefined && object.challenges !== null)
      ? RequestChallenges.fromPartial(object.challenges)
      : undefined;
    message.lifetime = (object.lifetime !== undefined && object.lifetime !== null)
      ? Duration.fromPartial(object.lifetime)
      : undefined;
    return message;
  },
};

function createBaseSetSessionRequest_MetadataEntry(): SetSessionRequest_MetadataEntry {
  return { key: "", value: Buffer.alloc(0) };
}

export const SetSessionRequest_MetadataEntry = {
  encode(message: SetSessionRequest_MetadataEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.value.length !== 0) {
      writer.uint32(18).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetSessionRequest_MetadataEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetSessionRequest_MetadataEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.value = reader.bytes() as Buffer;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetSessionRequest_MetadataEntry {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? Buffer.from(bytesFromBase64(object.value)) : Buffer.alloc(0),
    };
  },

  toJSON(message: SetSessionRequest_MetadataEntry): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = base64FromBytes(message.value !== undefined ? message.value : Buffer.alloc(0)));
    return obj;
  },

  create(base?: DeepPartial<SetSessionRequest_MetadataEntry>): SetSessionRequest_MetadataEntry {
    return SetSessionRequest_MetadataEntry.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetSessionRequest_MetadataEntry>): SetSessionRequest_MetadataEntry {
    const message = createBaseSetSessionRequest_MetadataEntry();
    message.key = object.key ?? "";
    message.value = object.value ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseSetSessionResponse(): SetSessionResponse {
  return { details: undefined, sessionToken: "", challenges: undefined };
}

export const SetSessionResponse = {
  encode(message: SetSessionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.sessionToken !== "") {
      writer.uint32(18).string(message.sessionToken);
    }
    if (message.challenges !== undefined) {
      Challenges.encode(message.challenges, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetSessionResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetSessionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sessionToken = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.challenges = Challenges.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetSessionResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      sessionToken: isSet(object.sessionToken) ? String(object.sessionToken) : "",
      challenges: isSet(object.challenges) ? Challenges.fromJSON(object.challenges) : undefined,
    };
  },

  toJSON(message: SetSessionResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.sessionToken !== undefined && (obj.sessionToken = message.sessionToken);
    message.challenges !== undefined &&
      (obj.challenges = message.challenges ? Challenges.toJSON(message.challenges) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetSessionResponse>): SetSessionResponse {
    return SetSessionResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetSessionResponse>): SetSessionResponse {
    const message = createBaseSetSessionResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.sessionToken = object.sessionToken ?? "";
    message.challenges = (object.challenges !== undefined && object.challenges !== null)
      ? Challenges.fromPartial(object.challenges)
      : undefined;
    return message;
  },
};

function createBaseDeleteSessionRequest(): DeleteSessionRequest {
  return { sessionId: "", sessionToken: undefined };
}

export const DeleteSessionRequest = {
  encode(message: DeleteSessionRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sessionId !== "") {
      writer.uint32(10).string(message.sessionId);
    }
    if (message.sessionToken !== undefined) {
      writer.uint32(18).string(message.sessionToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteSessionRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteSessionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.sessionId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sessionToken = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteSessionRequest {
    return {
      sessionId: isSet(object.sessionId) ? String(object.sessionId) : "",
      sessionToken: isSet(object.sessionToken) ? String(object.sessionToken) : undefined,
    };
  },

  toJSON(message: DeleteSessionRequest): unknown {
    const obj: any = {};
    message.sessionId !== undefined && (obj.sessionId = message.sessionId);
    message.sessionToken !== undefined && (obj.sessionToken = message.sessionToken);
    return obj;
  },

  create(base?: DeepPartial<DeleteSessionRequest>): DeleteSessionRequest {
    return DeleteSessionRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteSessionRequest>): DeleteSessionRequest {
    const message = createBaseDeleteSessionRequest();
    message.sessionId = object.sessionId ?? "";
    message.sessionToken = object.sessionToken ?? undefined;
    return message;
  },
};

function createBaseDeleteSessionResponse(): DeleteSessionResponse {
  return { details: undefined };
}

export const DeleteSessionResponse = {
  encode(message: DeleteSessionResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteSessionResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteSessionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteSessionResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeleteSessionResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeleteSessionResponse>): DeleteSessionResponse {
    return DeleteSessionResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteSessionResponse>): DeleteSessionResponse {
    const message = createBaseDeleteSessionResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseChecks(): Checks {
  return {
    user: undefined,
    password: undefined,
    webAuthN: undefined,
    idpIntent: undefined,
    totp: undefined,
    otpSms: undefined,
    otpEmail: undefined,
  };
}

export const Checks = {
  encode(message: Checks, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      CheckUser.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    if (message.password !== undefined) {
      CheckPassword.encode(message.password, writer.uint32(18).fork()).ldelim();
    }
    if (message.webAuthN !== undefined) {
      CheckWebAuthN.encode(message.webAuthN, writer.uint32(26).fork()).ldelim();
    }
    if (message.idpIntent !== undefined) {
      CheckIDPIntent.encode(message.idpIntent, writer.uint32(34).fork()).ldelim();
    }
    if (message.totp !== undefined) {
      CheckTOTP.encode(message.totp, writer.uint32(42).fork()).ldelim();
    }
    if (message.otpSms !== undefined) {
      CheckOTP.encode(message.otpSms, writer.uint32(50).fork()).ldelim();
    }
    if (message.otpEmail !== undefined) {
      CheckOTP.encode(message.otpEmail, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Checks {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChecks();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.user = CheckUser.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.password = CheckPassword.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.webAuthN = CheckWebAuthN.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.idpIntent = CheckIDPIntent.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.totp = CheckTOTP.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.otpSms = CheckOTP.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.otpEmail = CheckOTP.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Checks {
    return {
      user: isSet(object.user) ? CheckUser.fromJSON(object.user) : undefined,
      password: isSet(object.password) ? CheckPassword.fromJSON(object.password) : undefined,
      webAuthN: isSet(object.webAuthN) ? CheckWebAuthN.fromJSON(object.webAuthN) : undefined,
      idpIntent: isSet(object.idpIntent) ? CheckIDPIntent.fromJSON(object.idpIntent) : undefined,
      totp: isSet(object.totp) ? CheckTOTP.fromJSON(object.totp) : undefined,
      otpSms: isSet(object.otpSms) ? CheckOTP.fromJSON(object.otpSms) : undefined,
      otpEmail: isSet(object.otpEmail) ? CheckOTP.fromJSON(object.otpEmail) : undefined,
    };
  },

  toJSON(message: Checks): unknown {
    const obj: any = {};
    message.user !== undefined && (obj.user = message.user ? CheckUser.toJSON(message.user) : undefined);
    message.password !== undefined &&
      (obj.password = message.password ? CheckPassword.toJSON(message.password) : undefined);
    message.webAuthN !== undefined &&
      (obj.webAuthN = message.webAuthN ? CheckWebAuthN.toJSON(message.webAuthN) : undefined);
    message.idpIntent !== undefined &&
      (obj.idpIntent = message.idpIntent ? CheckIDPIntent.toJSON(message.idpIntent) : undefined);
    message.totp !== undefined && (obj.totp = message.totp ? CheckTOTP.toJSON(message.totp) : undefined);
    message.otpSms !== undefined && (obj.otpSms = message.otpSms ? CheckOTP.toJSON(message.otpSms) : undefined);
    message.otpEmail !== undefined && (obj.otpEmail = message.otpEmail ? CheckOTP.toJSON(message.otpEmail) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Checks>): Checks {
    return Checks.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Checks>): Checks {
    const message = createBaseChecks();
    message.user = (object.user !== undefined && object.user !== null) ? CheckUser.fromPartial(object.user) : undefined;
    message.password = (object.password !== undefined && object.password !== null)
      ? CheckPassword.fromPartial(object.password)
      : undefined;
    message.webAuthN = (object.webAuthN !== undefined && object.webAuthN !== null)
      ? CheckWebAuthN.fromPartial(object.webAuthN)
      : undefined;
    message.idpIntent = (object.idpIntent !== undefined && object.idpIntent !== null)
      ? CheckIDPIntent.fromPartial(object.idpIntent)
      : undefined;
    message.totp = (object.totp !== undefined && object.totp !== null) ? CheckTOTP.fromPartial(object.totp) : undefined;
    message.otpSms = (object.otpSms !== undefined && object.otpSms !== null)
      ? CheckOTP.fromPartial(object.otpSms)
      : undefined;
    message.otpEmail = (object.otpEmail !== undefined && object.otpEmail !== null)
      ? CheckOTP.fromPartial(object.otpEmail)
      : undefined;
    return message;
  },
};

function createBaseCheckUser(): CheckUser {
  return { userId: undefined, loginName: undefined };
}

export const CheckUser = {
  encode(message: CheckUser, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== undefined) {
      writer.uint32(10).string(message.userId);
    }
    if (message.loginName !== undefined) {
      writer.uint32(18).string(message.loginName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckUser {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckUser();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.loginName = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckUser {
    return {
      userId: isSet(object.userId) ? String(object.userId) : undefined,
      loginName: isSet(object.loginName) ? String(object.loginName) : undefined,
    };
  },

  toJSON(message: CheckUser): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.loginName !== undefined && (obj.loginName = message.loginName);
    return obj;
  },

  create(base?: DeepPartial<CheckUser>): CheckUser {
    return CheckUser.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CheckUser>): CheckUser {
    const message = createBaseCheckUser();
    message.userId = object.userId ?? undefined;
    message.loginName = object.loginName ?? undefined;
    return message;
  },
};

function createBaseCheckPassword(): CheckPassword {
  return { password: "" };
}

export const CheckPassword = {
  encode(message: CheckPassword, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.password !== "") {
      writer.uint32(10).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckPassword {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckPassword();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.password = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckPassword {
    return { password: isSet(object.password) ? String(object.password) : "" };
  },

  toJSON(message: CheckPassword): unknown {
    const obj: any = {};
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  create(base?: DeepPartial<CheckPassword>): CheckPassword {
    return CheckPassword.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CheckPassword>): CheckPassword {
    const message = createBaseCheckPassword();
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseCheckWebAuthN(): CheckWebAuthN {
  return { credentialAssertionData: undefined };
}

export const CheckWebAuthN = {
  encode(message: CheckWebAuthN, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.credentialAssertionData !== undefined) {
      Struct.encode(Struct.wrap(message.credentialAssertionData), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckWebAuthN {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckWebAuthN();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.credentialAssertionData = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckWebAuthN {
    return {
      credentialAssertionData: isObject(object.credentialAssertionData) ? object.credentialAssertionData : undefined,
    };
  },

  toJSON(message: CheckWebAuthN): unknown {
    const obj: any = {};
    message.credentialAssertionData !== undefined && (obj.credentialAssertionData = message.credentialAssertionData);
    return obj;
  },

  create(base?: DeepPartial<CheckWebAuthN>): CheckWebAuthN {
    return CheckWebAuthN.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CheckWebAuthN>): CheckWebAuthN {
    const message = createBaseCheckWebAuthN();
    message.credentialAssertionData = object.credentialAssertionData ?? undefined;
    return message;
  },
};

function createBaseCheckIDPIntent(): CheckIDPIntent {
  return { idpIntentId: "", idpIntentToken: "" };
}

export const CheckIDPIntent = {
  encode(message: CheckIDPIntent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpIntentId !== "") {
      writer.uint32(10).string(message.idpIntentId);
    }
    if (message.idpIntentToken !== "") {
      writer.uint32(18).string(message.idpIntentToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckIDPIntent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckIDPIntent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idpIntentId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idpIntentToken = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckIDPIntent {
    return {
      idpIntentId: isSet(object.idpIntentId) ? String(object.idpIntentId) : "",
      idpIntentToken: isSet(object.idpIntentToken) ? String(object.idpIntentToken) : "",
    };
  },

  toJSON(message: CheckIDPIntent): unknown {
    const obj: any = {};
    message.idpIntentId !== undefined && (obj.idpIntentId = message.idpIntentId);
    message.idpIntentToken !== undefined && (obj.idpIntentToken = message.idpIntentToken);
    return obj;
  },

  create(base?: DeepPartial<CheckIDPIntent>): CheckIDPIntent {
    return CheckIDPIntent.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CheckIDPIntent>): CheckIDPIntent {
    const message = createBaseCheckIDPIntent();
    message.idpIntentId = object.idpIntentId ?? "";
    message.idpIntentToken = object.idpIntentToken ?? "";
    return message;
  },
};

function createBaseCheckTOTP(): CheckTOTP {
  return { code: "" };
}

export const CheckTOTP = {
  encode(message: CheckTOTP, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.code !== "") {
      writer.uint32(10).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckTOTP {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckTOTP();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.code = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckTOTP {
    return { code: isSet(object.code) ? String(object.code) : "" };
  },

  toJSON(message: CheckTOTP): unknown {
    const obj: any = {};
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<CheckTOTP>): CheckTOTP {
    return CheckTOTP.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CheckTOTP>): CheckTOTP {
    const message = createBaseCheckTOTP();
    message.code = object.code ?? "";
    return message;
  },
};

function createBaseCheckOTP(): CheckOTP {
  return { code: "" };
}

export const CheckOTP = {
  encode(message: CheckOTP, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.code !== "") {
      writer.uint32(10).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckOTP {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckOTP();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.code = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckOTP {
    return { code: isSet(object.code) ? String(object.code) : "" };
  },

  toJSON(message: CheckOTP): unknown {
    const obj: any = {};
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<CheckOTP>): CheckOTP {
    return CheckOTP.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CheckOTP>): CheckOTP {
    const message = createBaseCheckOTP();
    message.code = object.code ?? "";
    return message;
  },
};

export type SessionServiceDefinition = typeof SessionServiceDefinition;
export const SessionServiceDefinition = {
  name: "SessionService",
  fullName: "zitadel.session.v2beta.SessionService",
  methods: {
    /** Search sessions */
    listSessions: {
      name: "ListSessions",
      requestType: ListSessionsRequest,
      requestStream: false,
      responseType: ListSessionsResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              109,
              18,
              15,
              83,
              101,
              97,
              114,
              99,
              104,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              115,
              26,
              19,
              83,
              101,
              97,
              114,
              99,
              104,
              32,
              102,
              111,
              114,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              115,
              74,
              11,
              10,
              3,
              50,
              48,
              48,
              18,
              4,
              10,
              2,
              79,
              75,
              74,
              56,
              10,
              3,
              52,
              48,
              48,
              18,
              49,
              10,
              18,
              105,
              110,
              118,
              97,
              108,
              105,
              100,
              32,
              108,
              105,
              115,
              116,
              32,
              113,
              117,
              101,
              114,
              121,
              18,
              27,
              10,
              25,
              26,
              23,
              35,
              47,
              100,
              101,
              102,
              105,
              110,
              105,
              116,
              105,
              111,
              110,
              115,
              47,
              114,
              112,
              99,
              83,
              116,
              97,
              116,
              117,
              115,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              28,
              58,
              1,
              42,
              34,
              23,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              115,
              47,
              115,
              101,
              97,
              114,
              99,
              104,
            ]),
          ],
        },
      },
    },
    /** GetSession a session */
    getSession: {
      name: "GetSession",
      requestType: GetSessionRequest,
      requestStream: false,
      responseType: GetSessionResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              118,
              18,
              13,
              71,
              101,
              116,
              32,
              97,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              26,
              88,
              71,
              101,
              116,
              32,
              97,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              32,
              97,
              110,
              100,
              32,
              97,
              108,
              108,
              32,
              105,
              116,
              115,
              32,
              105,
              110,
              102,
              111,
              114,
              109,
              97,
              116,
              105,
              111,
              110,
              32,
              108,
              105,
              107,
              101,
              32,
              116,
              104,
              101,
              32,
              116,
              105,
              109,
              101,
              32,
              111,
              102,
              32,
              116,
              104,
              101,
              32,
              117,
              115,
              101,
              114,
              32,
              111,
              114,
              32,
              112,
              97,
              115,
              115,
              119,
              111,
              114,
              100,
              32,
              118,
              101,
              114,
              105,
              102,
              105,
              99,
              97,
              116,
              105,
              111,
              110,
              74,
              11,
              10,
              3,
              50,
              48,
              48,
              18,
              4,
              10,
              2,
              79,
              75,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              31,
              18,
              29,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              115,
              47,
              123,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /** Create a new session */
    createSession: {
      name: "CreateSession",
      requestType: CreateSessionRequest,
      requestStream: false,
      responseType: CreateSessionResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              138,
              1,
              18,
              20,
              67,
              114,
              101,
              97,
              116,
              101,
              32,
              97,
              32,
              110,
              101,
              119,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              26,
              101,
              67,
              114,
              101,
              97,
              116,
              101,
              32,
              97,
              32,
              110,
              101,
              119,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              46,
              32,
              65,
              32,
              116,
              111,
              107,
              101,
              110,
              32,
              119,
              105,
              108,
              108,
              32,
              98,
              101,
              32,
              114,
              101,
              116,
              117,
              114,
              110,
              101,
              100,
              44,
              32,
              119,
              104,
              105,
              99,
              104,
              32,
              105,
              115,
              32,
              114,
              101,
              113,
              117,
              105,
              114,
              101,
              100,
              32,
              102,
              111,
              114,
              32,
              102,
              117,
              114,
              116,
              104,
              101,
              114,
              32,
              117,
              112,
              100,
              97,
              116,
              101,
              115,
              32,
              111,
              102,
              32,
              116,
              104,
              101,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              46,
              74,
              11,
              10,
              3,
              50,
              48,
              48,
              18,
              4,
              10,
              2,
              79,
              75,
            ]),
          ],
          400010: [
            Buffer.from([
              22,
              10,
              15,
              10,
              13,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              101,
              100,
              18,
              3,
              8,
              201,
              1,
            ]),
          ],
          578365826: [
            Buffer.from([
              21,
              58,
              1,
              42,
              34,
              16,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              115,
            ]),
          ],
        },
      },
    },
    /** Update a session */
    setSession: {
      name: "SetSession",
      requestType: SetSessionRequest,
      requestStream: false,
      responseType: SetSessionResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              91,
              18,
              26,
              85,
              112,
              100,
              97,
              116,
              101,
              32,
              97,
              110,
              32,
              101,
              120,
              105,
              115,
              116,
              105,
              110,
              103,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              26,
              48,
              85,
              112,
              100,
              97,
              116,
              101,
              32,
              97,
              110,
              32,
              101,
              120,
              105,
              115,
              116,
              105,
              110,
              103,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              32,
              119,
              105,
              116,
              104,
              32,
              110,
              101,
              119,
              32,
              105,
              110,
              102,
              111,
              114,
              109,
              97,
              116,
              105,
              111,
              110,
              46,
              74,
              11,
              10,
              3,
              50,
              48,
              48,
              18,
              4,
              10,
              2,
              79,
              75,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              34,
              58,
              1,
              42,
              50,
              29,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              115,
              47,
              123,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /** Terminate a session */
    deleteSession: {
      name: "DeleteSession",
      requestType: DeleteSessionRequest,
      requestStream: false,
      responseType: DeleteSessionResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              105,
              18,
              29,
              84,
              101,
              114,
              109,
              105,
              110,
              97,
              116,
              101,
              32,
              97,
              110,
              32,
              101,
              120,
              105,
              115,
              116,
              105,
              110,
              103,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              26,
              59,
              84,
              101,
              114,
              109,
              105,
              110,
              97,
              116,
              101,
              32,
              121,
              111,
              117,
              114,
              32,
              111,
              119,
              110,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              32,
              111,
              114,
              32,
              105,
              102,
              32,
              103,
              114,
              97,
              110,
              116,
              101,
              100,
              32,
              97,
              110,
              121,
              32,
              111,
              116,
              104,
              101,
              114,
              32,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              46,
              74,
              11,
              10,
              3,
              50,
              48,
              48,
              18,
              4,
              10,
              2,
              79,
              75,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              34,
              58,
              1,
              42,
              42,
              29,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              115,
              47,
              123,
              115,
              101,
              115,
              115,
              105,
              111,
              110,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
  },
} as const;

export interface SessionServiceImplementation<CallContextExt = {}> {
  /** Search sessions */
  listSessions(
    request: ListSessionsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ListSessionsResponse>>;
  /** GetSession a session */
  getSession(
    request: GetSessionRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetSessionResponse>>;
  /** Create a new session */
  createSession(
    request: CreateSessionRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<CreateSessionResponse>>;
  /** Update a session */
  setSession(
    request: SetSessionRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetSessionResponse>>;
  /** Terminate a session */
  deleteSession(
    request: DeleteSessionRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeleteSessionResponse>>;
}

export interface SessionServiceClient<CallOptionsExt = {}> {
  /** Search sessions */
  listSessions(
    request: DeepPartial<ListSessionsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ListSessionsResponse>;
  /** GetSession a session */
  getSession(
    request: DeepPartial<GetSessionRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetSessionResponse>;
  /** Create a new session */
  createSession(
    request: DeepPartial<CreateSessionRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<CreateSessionResponse>;
  /** Update a session */
  setSession(
    request: DeepPartial<SetSessionRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetSessionResponse>;
  /** Terminate a session */
  deleteSession(
    request: DeepPartial<DeleteSessionRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeleteSessionResponse>;
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

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Details } from "../../object/v2beta/object";
import { AuthorizationError, AuthRequest } from "./authorization";

export const protobufPackage = "zitadel.oidc.v2beta";

export interface GetAuthRequestRequest {
  authRequestId: string;
}

export interface GetAuthRequestResponse {
  authRequest: AuthRequest | undefined;
}

export interface CreateCallbackRequest {
  authRequestId: string;
  session?: Session | undefined;
  error?: AuthorizationError | undefined;
}

export interface Session {
  sessionId: string;
  sessionToken: string;
}

export interface CreateCallbackResponse {
  details: Details | undefined;
  callbackUrl: string;
}

function createBaseGetAuthRequestRequest(): GetAuthRequestRequest {
  return { authRequestId: "" };
}

export const GetAuthRequestRequest = {
  encode(message: GetAuthRequestRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authRequestId !== "") {
      writer.uint32(10).string(message.authRequestId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAuthRequestRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAuthRequestRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.authRequestId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetAuthRequestRequest {
    return { authRequestId: isSet(object.authRequestId) ? String(object.authRequestId) : "" };
  },

  toJSON(message: GetAuthRequestRequest): unknown {
    const obj: any = {};
    message.authRequestId !== undefined && (obj.authRequestId = message.authRequestId);
    return obj;
  },

  create(base?: DeepPartial<GetAuthRequestRequest>): GetAuthRequestRequest {
    return GetAuthRequestRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetAuthRequestRequest>): GetAuthRequestRequest {
    const message = createBaseGetAuthRequestRequest();
    message.authRequestId = object.authRequestId ?? "";
    return message;
  },
};

function createBaseGetAuthRequestResponse(): GetAuthRequestResponse {
  return { authRequest: undefined };
}

export const GetAuthRequestResponse = {
  encode(message: GetAuthRequestResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authRequest !== undefined) {
      AuthRequest.encode(message.authRequest, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAuthRequestResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAuthRequestResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.authRequest = AuthRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetAuthRequestResponse {
    return { authRequest: isSet(object.authRequest) ? AuthRequest.fromJSON(object.authRequest) : undefined };
  },

  toJSON(message: GetAuthRequestResponse): unknown {
    const obj: any = {};
    message.authRequest !== undefined &&
      (obj.authRequest = message.authRequest ? AuthRequest.toJSON(message.authRequest) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetAuthRequestResponse>): GetAuthRequestResponse {
    return GetAuthRequestResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetAuthRequestResponse>): GetAuthRequestResponse {
    const message = createBaseGetAuthRequestResponse();
    message.authRequest = (object.authRequest !== undefined && object.authRequest !== null)
      ? AuthRequest.fromPartial(object.authRequest)
      : undefined;
    return message;
  },
};

function createBaseCreateCallbackRequest(): CreateCallbackRequest {
  return { authRequestId: "", session: undefined, error: undefined };
}

export const CreateCallbackRequest = {
  encode(message: CreateCallbackRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authRequestId !== "") {
      writer.uint32(10).string(message.authRequestId);
    }
    if (message.session !== undefined) {
      Session.encode(message.session, writer.uint32(18).fork()).ldelim();
    }
    if (message.error !== undefined) {
      AuthorizationError.encode(message.error, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateCallbackRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateCallbackRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.authRequestId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.session = Session.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.error = AuthorizationError.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateCallbackRequest {
    return {
      authRequestId: isSet(object.authRequestId) ? String(object.authRequestId) : "",
      session: isSet(object.session) ? Session.fromJSON(object.session) : undefined,
      error: isSet(object.error) ? AuthorizationError.fromJSON(object.error) : undefined,
    };
  },

  toJSON(message: CreateCallbackRequest): unknown {
    const obj: any = {};
    message.authRequestId !== undefined && (obj.authRequestId = message.authRequestId);
    message.session !== undefined && (obj.session = message.session ? Session.toJSON(message.session) : undefined);
    message.error !== undefined && (obj.error = message.error ? AuthorizationError.toJSON(message.error) : undefined);
    return obj;
  },

  create(base?: DeepPartial<CreateCallbackRequest>): CreateCallbackRequest {
    return CreateCallbackRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateCallbackRequest>): CreateCallbackRequest {
    const message = createBaseCreateCallbackRequest();
    message.authRequestId = object.authRequestId ?? "";
    message.session = (object.session !== undefined && object.session !== null)
      ? Session.fromPartial(object.session)
      : undefined;
    message.error = (object.error !== undefined && object.error !== null)
      ? AuthorizationError.fromPartial(object.error)
      : undefined;
    return message;
  },
};

function createBaseSession(): Session {
  return { sessionId: "", sessionToken: "" };
}

export const Session = {
  encode(message: Session, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sessionId !== "") {
      writer.uint32(10).string(message.sessionId);
    }
    if (message.sessionToken !== "") {
      writer.uint32(18).string(message.sessionToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Session {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSession();
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

  fromJSON(object: any): Session {
    return {
      sessionId: isSet(object.sessionId) ? String(object.sessionId) : "",
      sessionToken: isSet(object.sessionToken) ? String(object.sessionToken) : "",
    };
  },

  toJSON(message: Session): unknown {
    const obj: any = {};
    message.sessionId !== undefined && (obj.sessionId = message.sessionId);
    message.sessionToken !== undefined && (obj.sessionToken = message.sessionToken);
    return obj;
  },

  create(base?: DeepPartial<Session>): Session {
    return Session.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Session>): Session {
    const message = createBaseSession();
    message.sessionId = object.sessionId ?? "";
    message.sessionToken = object.sessionToken ?? "";
    return message;
  },
};

function createBaseCreateCallbackResponse(): CreateCallbackResponse {
  return { details: undefined, callbackUrl: "" };
}

export const CreateCallbackResponse = {
  encode(message: CreateCallbackResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.callbackUrl !== "") {
      writer.uint32(18).string(message.callbackUrl);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateCallbackResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateCallbackResponse();
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

          message.callbackUrl = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateCallbackResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      callbackUrl: isSet(object.callbackUrl) ? String(object.callbackUrl) : "",
    };
  },

  toJSON(message: CreateCallbackResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.callbackUrl !== undefined && (obj.callbackUrl = message.callbackUrl);
    return obj;
  },

  create(base?: DeepPartial<CreateCallbackResponse>): CreateCallbackResponse {
    return CreateCallbackResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateCallbackResponse>): CreateCallbackResponse {
    const message = createBaseCreateCallbackResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.callbackUrl = object.callbackUrl ?? "";
    return message;
  },
};

export type OIDCServiceDefinition = typeof OIDCServiceDefinition;
export const OIDCServiceDefinition = {
  name: "OIDCService",
  fullName: "zitadel.oidc.v2beta.OIDCService",
  methods: {
    getAuthRequest: {
      name: "GetAuthRequest",
      requestType: GetAuthRequestRequest,
      requestStream: false,
      responseType: GetAuthRequestResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              184,
              1,
              18,
              29,
              71,
              101,
              116,
              32,
              79,
              73,
              68,
              67,
              32,
              65,
              117,
              116,
              104,
              32,
              82,
              101,
              113,
              117,
              101,
              115,
              116,
              32,
              100,
              101,
              116,
              97,
              105,
              108,
              115,
              26,
              137,
              1,
              71,
              101,
              116,
              32,
              79,
              73,
              68,
              67,
              32,
              65,
              117,
              116,
              104,
              32,
              82,
              101,
              113,
              117,
              101,
              115,
              116,
              32,
              100,
              101,
              116,
              97,
              105,
              108,
              115,
              32,
              98,
              121,
              32,
              73,
              68,
              44,
              32,
              111,
              98,
              116,
              97,
              105,
              110,
              101,
              100,
              32,
              102,
              114,
              111,
              109,
              32,
              116,
              104,
              101,
              32,
              114,
              101,
              100,
              105,
              114,
              101,
              99,
              116,
              32,
              85,
              82,
              76,
              46,
              32,
              82,
              101,
              116,
              117,
              114,
              110,
              115,
              32,
              100,
              101,
              116,
              97,
              105,
              108,
              115,
              32,
              116,
              104,
              97,
              116,
              32,
              97,
              114,
              101,
              32,
              112,
              97,
              114,
              115,
              101,
              100,
              32,
              102,
              114,
              111,
              109,
              32,
              116,
              104,
              101,
              32,
              97,
              112,
              112,
              108,
              105,
              99,
              97,
              116,
              105,
              111,
              110,
              39,
              115,
              32,
              65,
              117,
              116,
              104,
              32,
              82,
              101,
              113,
              117,
              101,
              115,
              116,
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
              46,
              18,
              44,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              111,
              105,
              100,
              99,
              47,
              97,
              117,
              116,
              104,
              95,
              114,
              101,
              113,
              117,
              101,
              115,
              116,
              115,
              47,
              123,
              97,
              117,
              116,
              104,
              95,
              114,
              101,
              113,
              117,
              101,
              115,
              116,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    createCallback: {
      name: "CreateCallback",
      requestType: CreateCallbackRequest,
      requestStream: false,
      responseType: CreateCallbackResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              253,
              2,
              18,
              50,
              70,
              105,
              110,
              97,
              108,
              105,
              122,
              101,
              32,
              97,
              110,
              32,
              65,
              117,
              116,
              104,
              32,
              82,
              101,
              113,
              117,
              101,
              115,
              116,
              32,
              97,
              110,
              100,
              32,
              103,
              101,
              116,
              32,
              116,
              104,
              101,
              32,
              99,
              97,
              108,
              108,
              98,
              97,
              99,
              107,
              32,
              85,
              82,
              76,
              46,
              26,
              185,
              2,
              70,
              105,
              110,
              97,
              108,
              105,
              122,
              101,
              32,
              97,
              110,
              32,
              65,
              117,
              116,
              104,
              32,
              82,
              101,
              113,
              117,
              101,
              115,
              116,
              32,
              97,
              110,
              100,
              32,
              103,
              101,
              116,
              32,
              116,
              104,
              101,
              32,
              99,
              97,
              108,
              108,
              98,
              97,
              99,
              107,
              32,
              85,
              82,
              76,
              32,
              102,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              32,
              111,
              114,
              32,
              102,
              97,
              105,
              108,
              117,
              114,
              101,
              46,
              32,
              84,
              104,
              101,
              32,
              117,
              115,
              101,
              114,
              32,
              109,
              117,
              115,
              116,
              32,
              98,
              101,
              32,
              114,
              101,
              100,
              105,
              114,
              101,
              99,
              116,
              101,
              100,
              32,
              116,
              111,
              32,
              116,
              104,
              101,
              32,
              85,
              82,
              76,
              32,
              105,
              110,
              32,
              111,
              114,
              100,
              101,
              114,
              32,
              116,
              111,
              32,
              105,
              110,
              102,
              111,
              114,
              109,
              32,
              116,
              104,
              101,
              32,
              97,
              112,
              112,
              108,
              105,
              99,
              97,
              116,
              105,
              111,
              110,
              32,
              97,
              98,
              111,
              117,
              116,
              32,
              116,
              104,
              101,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              32,
              111,
              114,
              32,
              102,
              97,
              105,
              108,
              117,
              114,
              101,
              46,
              32,
              79,
              110,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              44,
              32,
              116,
              104,
              101,
              32,
              85,
              82,
              76,
              32,
              99,
              111,
              110,
              116,
              97,
              105,
              110,
              115,
              32,
              100,
              101,
              116,
              97,
              105,
              108,
              115,
              32,
              102,
              111,
              114,
              32,
              116,
              104,
              101,
              32,
              97,
              112,
              112,
              108,
              105,
              99,
              97,
              116,
              105,
              111,
              110,
              32,
              116,
              111,
              32,
              111,
              98,
              116,
              97,
              105,
              110,
              32,
              116,
              104,
              101,
              32,
              116,
              111,
              107,
              101,
              110,
              115,
              46,
              32,
              84,
              104,
              105,
              115,
              32,
              109,
              101,
              116,
              104,
              111,
              100,
              32,
              99,
              97,
              110,
              32,
              111,
              110,
              108,
              121,
              32,
              98,
              101,
              32,
              99,
              97,
              108,
              108,
              101,
              100,
              32,
              111,
              110,
              99,
              101,
              32,
              102,
              111,
              114,
              32,
              97,
              110,
              32,
              65,
              117,
              116,
              104,
              32,
              114,
              101,
              113,
              117,
              101,
              115,
              116,
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
              49,
              58,
              1,
              42,
              34,
              44,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              111,
              105,
              100,
              99,
              47,
              97,
              117,
              116,
              104,
              95,
              114,
              101,
              113,
              117,
              101,
              115,
              116,
              115,
              47,
              123,
              97,
              117,
              116,
              104,
              95,
              114,
              101,
              113,
              117,
              101,
              115,
              116,
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

export interface OIDCServiceImplementation<CallContextExt = {}> {
  getAuthRequest(
    request: GetAuthRequestRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetAuthRequestResponse>>;
  createCallback(
    request: CreateCallbackRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<CreateCallbackResponse>>;
}

export interface OIDCServiceClient<CallOptionsExt = {}> {
  getAuthRequest(
    request: DeepPartial<GetAuthRequestRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetAuthRequestResponse>;
  createCallback(
    request: DeepPartial<CreateCallbackRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<CreateCallbackResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

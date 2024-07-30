/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Details } from "../../object/v2beta/object";
import { AddHumanUserRequest } from "../../user/v2beta/user_service";

export const protobufPackage = "zitadel.org.v2beta";

export interface AddOrganizationRequest {
  name: string;
  admins: AddOrganizationRequest_Admin[];
}

export interface AddOrganizationRequest_Admin {
  userId?: string | undefined;
  human?:
    | AddHumanUserRequest
    | undefined;
  /** specify Org Member Roles for the provided user (default is ORG_OWNER if roles are empty) */
  roles: string[];
}

export interface AddOrganizationResponse {
  details: Details | undefined;
  organizationId: string;
  createdAdmins: AddOrganizationResponse_CreatedAdmin[];
}

export interface AddOrganizationResponse_CreatedAdmin {
  userId: string;
  emailCode?: string | undefined;
  phoneCode?: string | undefined;
}

function createBaseAddOrganizationRequest(): AddOrganizationRequest {
  return { name: "", admins: [] };
}

export const AddOrganizationRequest = {
  encode(message: AddOrganizationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.admins) {
      AddOrganizationRequest_Admin.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOrganizationRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOrganizationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.admins.push(AddOrganizationRequest_Admin.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddOrganizationRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      admins: Array.isArray(object?.admins)
        ? object.admins.map((e: any) => AddOrganizationRequest_Admin.fromJSON(e))
        : [],
    };
  },

  toJSON(message: AddOrganizationRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    if (message.admins) {
      obj.admins = message.admins.map((e) => e ? AddOrganizationRequest_Admin.toJSON(e) : undefined);
    } else {
      obj.admins = [];
    }
    return obj;
  },

  create(base?: DeepPartial<AddOrganizationRequest>): AddOrganizationRequest {
    return AddOrganizationRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOrganizationRequest>): AddOrganizationRequest {
    const message = createBaseAddOrganizationRequest();
    message.name = object.name ?? "";
    message.admins = object.admins?.map((e) => AddOrganizationRequest_Admin.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAddOrganizationRequest_Admin(): AddOrganizationRequest_Admin {
  return { userId: undefined, human: undefined, roles: [] };
}

export const AddOrganizationRequest_Admin = {
  encode(message: AddOrganizationRequest_Admin, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== undefined) {
      writer.uint32(10).string(message.userId);
    }
    if (message.human !== undefined) {
      AddHumanUserRequest.encode(message.human, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.roles) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOrganizationRequest_Admin {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOrganizationRequest_Admin();
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

          message.human = AddHumanUserRequest.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.roles.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddOrganizationRequest_Admin {
    return {
      userId: isSet(object.userId) ? String(object.userId) : undefined,
      human: isSet(object.human) ? AddHumanUserRequest.fromJSON(object.human) : undefined,
      roles: Array.isArray(object?.roles) ? object.roles.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: AddOrganizationRequest_Admin): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.human !== undefined && (obj.human = message.human ? AddHumanUserRequest.toJSON(message.human) : undefined);
    if (message.roles) {
      obj.roles = message.roles.map((e) => e);
    } else {
      obj.roles = [];
    }
    return obj;
  },

  create(base?: DeepPartial<AddOrganizationRequest_Admin>): AddOrganizationRequest_Admin {
    return AddOrganizationRequest_Admin.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOrganizationRequest_Admin>): AddOrganizationRequest_Admin {
    const message = createBaseAddOrganizationRequest_Admin();
    message.userId = object.userId ?? undefined;
    message.human = (object.human !== undefined && object.human !== null)
      ? AddHumanUserRequest.fromPartial(object.human)
      : undefined;
    message.roles = object.roles?.map((e) => e) || [];
    return message;
  },
};

function createBaseAddOrganizationResponse(): AddOrganizationResponse {
  return { details: undefined, organizationId: "", createdAdmins: [] };
}

export const AddOrganizationResponse = {
  encode(message: AddOrganizationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.organizationId !== "") {
      writer.uint32(18).string(message.organizationId);
    }
    for (const v of message.createdAdmins) {
      AddOrganizationResponse_CreatedAdmin.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOrganizationResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOrganizationResponse();
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

          message.organizationId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.createdAdmins.push(AddOrganizationResponse_CreatedAdmin.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddOrganizationResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      organizationId: isSet(object.organizationId) ? String(object.organizationId) : "",
      createdAdmins: Array.isArray(object?.createdAdmins)
        ? object.createdAdmins.map((e: any) => AddOrganizationResponse_CreatedAdmin.fromJSON(e))
        : [],
    };
  },

  toJSON(message: AddOrganizationResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.organizationId !== undefined && (obj.organizationId = message.organizationId);
    if (message.createdAdmins) {
      obj.createdAdmins = message.createdAdmins.map((e) =>
        e ? AddOrganizationResponse_CreatedAdmin.toJSON(e) : undefined
      );
    } else {
      obj.createdAdmins = [];
    }
    return obj;
  },

  create(base?: DeepPartial<AddOrganizationResponse>): AddOrganizationResponse {
    return AddOrganizationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOrganizationResponse>): AddOrganizationResponse {
    const message = createBaseAddOrganizationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.organizationId = object.organizationId ?? "";
    message.createdAdmins = object.createdAdmins?.map((e) => AddOrganizationResponse_CreatedAdmin.fromPartial(e)) || [];
    return message;
  },
};

function createBaseAddOrganizationResponse_CreatedAdmin(): AddOrganizationResponse_CreatedAdmin {
  return { userId: "", emailCode: undefined, phoneCode: undefined };
}

export const AddOrganizationResponse_CreatedAdmin = {
  encode(message: AddOrganizationResponse_CreatedAdmin, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.emailCode !== undefined) {
      writer.uint32(18).string(message.emailCode);
    }
    if (message.phoneCode !== undefined) {
      writer.uint32(26).string(message.phoneCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOrganizationResponse_CreatedAdmin {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOrganizationResponse_CreatedAdmin();
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

          message.emailCode = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.phoneCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddOrganizationResponse_CreatedAdmin {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      emailCode: isSet(object.emailCode) ? String(object.emailCode) : undefined,
      phoneCode: isSet(object.phoneCode) ? String(object.phoneCode) : undefined,
    };
  },

  toJSON(message: AddOrganizationResponse_CreatedAdmin): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.emailCode !== undefined && (obj.emailCode = message.emailCode);
    message.phoneCode !== undefined && (obj.phoneCode = message.phoneCode);
    return obj;
  },

  create(base?: DeepPartial<AddOrganizationResponse_CreatedAdmin>): AddOrganizationResponse_CreatedAdmin {
    return AddOrganizationResponse_CreatedAdmin.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOrganizationResponse_CreatedAdmin>): AddOrganizationResponse_CreatedAdmin {
    const message = createBaseAddOrganizationResponse_CreatedAdmin();
    message.userId = object.userId ?? "";
    message.emailCode = object.emailCode ?? undefined;
    message.phoneCode = object.phoneCode ?? undefined;
    return message;
  },
};

export type OrganizationServiceDefinition = typeof OrganizationServiceDefinition;
export const OrganizationServiceDefinition = {
  name: "OrganizationService",
  fullName: "zitadel.org.v2beta.OrganizationService",
  methods: {
    /** Create a new organization and grant the user(s) permission to manage it */
    addOrganization: {
      name: "AddOrganization",
      requestType: AddOrganizationRequest,
      requestStream: false,
      responseType: AddOrganizationResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              180,
              1,
              18,
              22,
              67,
              114,
              101,
              97,
              116,
              101,
              32,
              97,
              110,
              32,
              79,
              114,
              103,
              97,
              110,
              105,
              122,
              97,
              116,
              105,
              111,
              110,
              26,
              140,
              1,
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
              111,
              114,
              103,
              97,
              110,
              105,
              122,
              97,
              116,
              105,
              111,
              110,
              32,
              119,
              105,
              116,
              104,
              32,
              97,
              110,
              32,
              97,
              100,
              109,
              105,
              110,
              105,
              115,
              116,
              114,
              97,
              116,
              105,
              118,
              101,
              32,
              117,
              115,
              101,
              114,
              46,
              32,
              73,
              102,
              32,
              110,
              111,
              32,
              115,
              112,
              101,
              99,
              105,
              102,
              105,
              99,
              32,
              114,
              111,
              108,
              101,
              115,
              32,
              97,
              114,
              101,
              32,
              115,
              101,
              110,
              116,
              32,
              102,
              111,
              114,
              32,
              116,
              104,
              101,
              32,
              117,
              115,
              101,
              114,
              115,
              44,
              32,
              116,
              104,
              101,
              121,
              32,
              119,
              105,
              108,
              108,
              32,
              98,
              101,
              32,
              103,
              114,
              97,
              110,
              116,
              101,
              100,
              32,
              116,
              104,
              101,
              32,
              114,
              111,
              108,
              101,
              32,
              79,
              82,
              71,
              95,
              79,
              87,
              78,
              69,
              82,
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
          400010: [Buffer.from([19, 10, 12, 10, 10, 111, 114, 103, 46, 99, 114, 101, 97, 116, 101, 18, 3, 8, 201, 1])],
          578365826: [
            Buffer.from([
              26,
              58,
              1,
              42,
              34,
              21,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              111,
              114,
              103,
              97,
              110,
              105,
              122,
              97,
              116,
              105,
              111,
              110,
              115,
            ]),
          ],
        },
      },
    },
  },
} as const;

export interface OrganizationServiceImplementation<CallContextExt = {}> {
  /** Create a new organization and grant the user(s) permission to manage it */
  addOrganization(
    request: AddOrganizationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddOrganizationResponse>>;
}

export interface OrganizationServiceClient<CallOptionsExt = {}> {
  /** Create a new organization and grant the user(s) permission to manage it */
  addOrganization(
    request: DeepPartial<AddOrganizationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddOrganizationResponse>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

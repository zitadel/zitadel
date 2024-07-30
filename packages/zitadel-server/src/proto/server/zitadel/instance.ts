/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { ObjectDetails, TextQueryMethod, textQueryMethodFromJSON, textQueryMethodToJSON } from "./object";

export const protobufPackage = "zitadel.instance.v1";

export enum State {
  STATE_UNSPECIFIED = 0,
  STATE_CREATING = 1,
  STATE_RUNNING = 2,
  STATE_STOPPING = 3,
  STATE_STOPPED = 4,
  UNRECOGNIZED = -1,
}

export function stateFromJSON(object: any): State {
  switch (object) {
    case 0:
    case "STATE_UNSPECIFIED":
      return State.STATE_UNSPECIFIED;
    case 1:
    case "STATE_CREATING":
      return State.STATE_CREATING;
    case 2:
    case "STATE_RUNNING":
      return State.STATE_RUNNING;
    case 3:
    case "STATE_STOPPING":
      return State.STATE_STOPPING;
    case 4:
    case "STATE_STOPPED":
      return State.STATE_STOPPED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return State.UNRECOGNIZED;
  }
}

export function stateToJSON(object: State): string {
  switch (object) {
    case State.STATE_UNSPECIFIED:
      return "STATE_UNSPECIFIED";
    case State.STATE_CREATING:
      return "STATE_CREATING";
    case State.STATE_RUNNING:
      return "STATE_RUNNING";
    case State.STATE_STOPPING:
      return "STATE_STOPPING";
    case State.STATE_STOPPED:
      return "STATE_STOPPED";
    case State.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum FieldName {
  FIELD_NAME_UNSPECIFIED = 0,
  FIELD_NAME_ID = 1,
  FIELD_NAME_NAME = 2,
  FIELD_NAME_CREATION_DATE = 3,
  UNRECOGNIZED = -1,
}

export function fieldNameFromJSON(object: any): FieldName {
  switch (object) {
    case 0:
    case "FIELD_NAME_UNSPECIFIED":
      return FieldName.FIELD_NAME_UNSPECIFIED;
    case 1:
    case "FIELD_NAME_ID":
      return FieldName.FIELD_NAME_ID;
    case 2:
    case "FIELD_NAME_NAME":
      return FieldName.FIELD_NAME_NAME;
    case 3:
    case "FIELD_NAME_CREATION_DATE":
      return FieldName.FIELD_NAME_CREATION_DATE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FieldName.UNRECOGNIZED;
  }
}

export function fieldNameToJSON(object: FieldName): string {
  switch (object) {
    case FieldName.FIELD_NAME_UNSPECIFIED:
      return "FIELD_NAME_UNSPECIFIED";
    case FieldName.FIELD_NAME_ID:
      return "FIELD_NAME_ID";
    case FieldName.FIELD_NAME_NAME:
      return "FIELD_NAME_NAME";
    case FieldName.FIELD_NAME_CREATION_DATE:
      return "FIELD_NAME_CREATION_DATE";
    case FieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum DomainFieldName {
  DOMAIN_FIELD_NAME_UNSPECIFIED = 0,
  DOMAIN_FIELD_NAME_DOMAIN = 1,
  DOMAIN_FIELD_NAME_PRIMARY = 2,
  DOMAIN_FIELD_NAME_GENERATED = 3,
  DOMAIN_FIELD_NAME_CREATION_DATE = 4,
  UNRECOGNIZED = -1,
}

export function domainFieldNameFromJSON(object: any): DomainFieldName {
  switch (object) {
    case 0:
    case "DOMAIN_FIELD_NAME_UNSPECIFIED":
      return DomainFieldName.DOMAIN_FIELD_NAME_UNSPECIFIED;
    case 1:
    case "DOMAIN_FIELD_NAME_DOMAIN":
      return DomainFieldName.DOMAIN_FIELD_NAME_DOMAIN;
    case 2:
    case "DOMAIN_FIELD_NAME_PRIMARY":
      return DomainFieldName.DOMAIN_FIELD_NAME_PRIMARY;
    case 3:
    case "DOMAIN_FIELD_NAME_GENERATED":
      return DomainFieldName.DOMAIN_FIELD_NAME_GENERATED;
    case 4:
    case "DOMAIN_FIELD_NAME_CREATION_DATE":
      return DomainFieldName.DOMAIN_FIELD_NAME_CREATION_DATE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return DomainFieldName.UNRECOGNIZED;
  }
}

export function domainFieldNameToJSON(object: DomainFieldName): string {
  switch (object) {
    case DomainFieldName.DOMAIN_FIELD_NAME_UNSPECIFIED:
      return "DOMAIN_FIELD_NAME_UNSPECIFIED";
    case DomainFieldName.DOMAIN_FIELD_NAME_DOMAIN:
      return "DOMAIN_FIELD_NAME_DOMAIN";
    case DomainFieldName.DOMAIN_FIELD_NAME_PRIMARY:
      return "DOMAIN_FIELD_NAME_PRIMARY";
    case DomainFieldName.DOMAIN_FIELD_NAME_GENERATED:
      return "DOMAIN_FIELD_NAME_GENERATED";
    case DomainFieldName.DOMAIN_FIELD_NAME_CREATION_DATE:
      return "DOMAIN_FIELD_NAME_CREATION_DATE";
    case DomainFieldName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Instance {
  id: string;
  details: ObjectDetails | undefined;
  state: State;
  name: string;
  version: string;
  domains: Domain[];
}

export interface InstanceDetail {
  id: string;
  details: ObjectDetails | undefined;
  state: State;
  name: string;
  version: string;
  domains: Domain[];
}

export interface Query {
  idQuery?: IdsQuery | undefined;
  domainQuery?: DomainsQuery | undefined;
}

/** IdQuery always equals */
export interface IdsQuery {
  ids: string[];
}

export interface DomainsQuery {
  domains: string[];
}

export interface Domain {
  details: ObjectDetails | undefined;
  domain: string;
  primary: boolean;
  generated: boolean;
}

export interface DomainSearchQuery {
  domainQuery?: DomainQuery | undefined;
  generatedQuery?: DomainGeneratedQuery | undefined;
  primaryQuery?: DomainPrimaryQuery | undefined;
}

export interface DomainQuery {
  domain: string;
  method: TextQueryMethod;
}

/** DomainGeneratedQuery is always equals */
export interface DomainGeneratedQuery {
  generated: boolean;
}

/** DomainPrimaryQuery is always equals */
export interface DomainPrimaryQuery {
  primary: boolean;
}

function createBaseInstance(): Instance {
  return { id: "", details: undefined, state: 0, name: "", version: "", domains: [] };
}

export const Instance = {
  encode(message: Instance, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.state !== 0) {
      writer.uint32(24).int32(message.state);
    }
    if (message.name !== "") {
      writer.uint32(34).string(message.name);
    }
    if (message.version !== "") {
      writer.uint32(42).string(message.version);
    }
    for (const v of message.domains) {
      Domain.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Instance {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstance();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.name = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.version = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.domains.push(Domain.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Instance {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      state: isSet(object.state) ? stateFromJSON(object.state) : 0,
      name: isSet(object.name) ? String(object.name) : "",
      version: isSet(object.version) ? String(object.version) : "",
      domains: Array.isArray(object?.domains) ? object.domains.map((e: any) => Domain.fromJSON(e)) : [],
    };
  },

  toJSON(message: Instance): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.state !== undefined && (obj.state = stateToJSON(message.state));
    message.name !== undefined && (obj.name = message.name);
    message.version !== undefined && (obj.version = message.version);
    if (message.domains) {
      obj.domains = message.domains.map((e) => e ? Domain.toJSON(e) : undefined);
    } else {
      obj.domains = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Instance>): Instance {
    return Instance.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Instance>): Instance {
    const message = createBaseInstance();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.state = object.state ?? 0;
    message.name = object.name ?? "";
    message.version = object.version ?? "";
    message.domains = object.domains?.map((e) => Domain.fromPartial(e)) || [];
    return message;
  },
};

function createBaseInstanceDetail(): InstanceDetail {
  return { id: "", details: undefined, state: 0, name: "", version: "", domains: [] };
}

export const InstanceDetail = {
  encode(message: InstanceDetail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.state !== 0) {
      writer.uint32(24).int32(message.state);
    }
    if (message.name !== "") {
      writer.uint32(34).string(message.name);
    }
    if (message.version !== "") {
      writer.uint32(42).string(message.version);
    }
    for (const v of message.domains) {
      Domain.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InstanceDetail {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInstanceDetail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.name = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.version = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.domains.push(Domain.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): InstanceDetail {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      state: isSet(object.state) ? stateFromJSON(object.state) : 0,
      name: isSet(object.name) ? String(object.name) : "",
      version: isSet(object.version) ? String(object.version) : "",
      domains: Array.isArray(object?.domains) ? object.domains.map((e: any) => Domain.fromJSON(e)) : [],
    };
  },

  toJSON(message: InstanceDetail): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.state !== undefined && (obj.state = stateToJSON(message.state));
    message.name !== undefined && (obj.name = message.name);
    message.version !== undefined && (obj.version = message.version);
    if (message.domains) {
      obj.domains = message.domains.map((e) => e ? Domain.toJSON(e) : undefined);
    } else {
      obj.domains = [];
    }
    return obj;
  },

  create(base?: DeepPartial<InstanceDetail>): InstanceDetail {
    return InstanceDetail.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InstanceDetail>): InstanceDetail {
    const message = createBaseInstanceDetail();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.state = object.state ?? 0;
    message.name = object.name ?? "";
    message.version = object.version ?? "";
    message.domains = object.domains?.map((e) => Domain.fromPartial(e)) || [];
    return message;
  },
};

function createBaseQuery(): Query {
  return { idQuery: undefined, domainQuery: undefined };
}

export const Query = {
  encode(message: Query, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idQuery !== undefined) {
      IdsQuery.encode(message.idQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.domainQuery !== undefined) {
      DomainsQuery.encode(message.domainQuery, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Query {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idQuery = IdsQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.domainQuery = DomainsQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Query {
    return {
      idQuery: isSet(object.idQuery) ? IdsQuery.fromJSON(object.idQuery) : undefined,
      domainQuery: isSet(object.domainQuery) ? DomainsQuery.fromJSON(object.domainQuery) : undefined,
    };
  },

  toJSON(message: Query): unknown {
    const obj: any = {};
    message.idQuery !== undefined && (obj.idQuery = message.idQuery ? IdsQuery.toJSON(message.idQuery) : undefined);
    message.domainQuery !== undefined &&
      (obj.domainQuery = message.domainQuery ? DomainsQuery.toJSON(message.domainQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Query>): Query {
    return Query.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Query>): Query {
    const message = createBaseQuery();
    message.idQuery = (object.idQuery !== undefined && object.idQuery !== null)
      ? IdsQuery.fromPartial(object.idQuery)
      : undefined;
    message.domainQuery = (object.domainQuery !== undefined && object.domainQuery !== null)
      ? DomainsQuery.fromPartial(object.domainQuery)
      : undefined;
    return message;
  },
};

function createBaseIdsQuery(): IdsQuery {
  return { ids: [] };
}

export const IdsQuery = {
  encode(message: IdsQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.ids) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IdsQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdsQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.ids.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IdsQuery {
    return { ids: Array.isArray(object?.ids) ? object.ids.map((e: any) => String(e)) : [] };
  },

  toJSON(message: IdsQuery): unknown {
    const obj: any = {};
    if (message.ids) {
      obj.ids = message.ids.map((e) => e);
    } else {
      obj.ids = [];
    }
    return obj;
  },

  create(base?: DeepPartial<IdsQuery>): IdsQuery {
    return IdsQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IdsQuery>): IdsQuery {
    const message = createBaseIdsQuery();
    message.ids = object.ids?.map((e) => e) || [];
    return message;
  },
};

function createBaseDomainsQuery(): DomainsQuery {
  return { domains: [] };
}

export const DomainsQuery = {
  encode(message: DomainsQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.domains) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DomainsQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDomainsQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.domains.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DomainsQuery {
    return { domains: Array.isArray(object?.domains) ? object.domains.map((e: any) => String(e)) : [] };
  },

  toJSON(message: DomainsQuery): unknown {
    const obj: any = {};
    if (message.domains) {
      obj.domains = message.domains.map((e) => e);
    } else {
      obj.domains = [];
    }
    return obj;
  },

  create(base?: DeepPartial<DomainsQuery>): DomainsQuery {
    return DomainsQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DomainsQuery>): DomainsQuery {
    const message = createBaseDomainsQuery();
    message.domains = object.domains?.map((e) => e) || [];
    return message;
  },
};

function createBaseDomain(): Domain {
  return { details: undefined, domain: "", primary: false, generated: false };
}

export const Domain = {
  encode(message: Domain, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.domain !== "") {
      writer.uint32(18).string(message.domain);
    }
    if (message.primary === true) {
      writer.uint32(24).bool(message.primary);
    }
    if (message.generated === true) {
      writer.uint32(32).bool(message.generated);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Domain {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDomain();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.domain = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.primary = reader.bool();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.generated = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Domain {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      domain: isSet(object.domain) ? String(object.domain) : "",
      primary: isSet(object.primary) ? Boolean(object.primary) : false,
      generated: isSet(object.generated) ? Boolean(object.generated) : false,
    };
  },

  toJSON(message: Domain): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.domain !== undefined && (obj.domain = message.domain);
    message.primary !== undefined && (obj.primary = message.primary);
    message.generated !== undefined && (obj.generated = message.generated);
    return obj;
  },

  create(base?: DeepPartial<Domain>): Domain {
    return Domain.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Domain>): Domain {
    const message = createBaseDomain();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.domain = object.domain ?? "";
    message.primary = object.primary ?? false;
    message.generated = object.generated ?? false;
    return message;
  },
};

function createBaseDomainSearchQuery(): DomainSearchQuery {
  return { domainQuery: undefined, generatedQuery: undefined, primaryQuery: undefined };
}

export const DomainSearchQuery = {
  encode(message: DomainSearchQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.domainQuery !== undefined) {
      DomainQuery.encode(message.domainQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.generatedQuery !== undefined) {
      DomainGeneratedQuery.encode(message.generatedQuery, writer.uint32(18).fork()).ldelim();
    }
    if (message.primaryQuery !== undefined) {
      DomainPrimaryQuery.encode(message.primaryQuery, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DomainSearchQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDomainSearchQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.domainQuery = DomainQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.generatedQuery = DomainGeneratedQuery.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.primaryQuery = DomainPrimaryQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DomainSearchQuery {
    return {
      domainQuery: isSet(object.domainQuery) ? DomainQuery.fromJSON(object.domainQuery) : undefined,
      generatedQuery: isSet(object.generatedQuery) ? DomainGeneratedQuery.fromJSON(object.generatedQuery) : undefined,
      primaryQuery: isSet(object.primaryQuery) ? DomainPrimaryQuery.fromJSON(object.primaryQuery) : undefined,
    };
  },

  toJSON(message: DomainSearchQuery): unknown {
    const obj: any = {};
    message.domainQuery !== undefined &&
      (obj.domainQuery = message.domainQuery ? DomainQuery.toJSON(message.domainQuery) : undefined);
    message.generatedQuery !== undefined &&
      (obj.generatedQuery = message.generatedQuery ? DomainGeneratedQuery.toJSON(message.generatedQuery) : undefined);
    message.primaryQuery !== undefined &&
      (obj.primaryQuery = message.primaryQuery ? DomainPrimaryQuery.toJSON(message.primaryQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DomainSearchQuery>): DomainSearchQuery {
    return DomainSearchQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DomainSearchQuery>): DomainSearchQuery {
    const message = createBaseDomainSearchQuery();
    message.domainQuery = (object.domainQuery !== undefined && object.domainQuery !== null)
      ? DomainQuery.fromPartial(object.domainQuery)
      : undefined;
    message.generatedQuery = (object.generatedQuery !== undefined && object.generatedQuery !== null)
      ? DomainGeneratedQuery.fromPartial(object.generatedQuery)
      : undefined;
    message.primaryQuery = (object.primaryQuery !== undefined && object.primaryQuery !== null)
      ? DomainPrimaryQuery.fromPartial(object.primaryQuery)
      : undefined;
    return message;
  },
};

function createBaseDomainQuery(): DomainQuery {
  return { domain: "", method: 0 };
}

export const DomainQuery = {
  encode(message: DomainQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.domain !== "") {
      writer.uint32(10).string(message.domain);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DomainQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDomainQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.domain = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DomainQuery {
    return {
      domain: isSet(object.domain) ? String(object.domain) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: DomainQuery): unknown {
    const obj: any = {};
    message.domain !== undefined && (obj.domain = message.domain);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<DomainQuery>): DomainQuery {
    return DomainQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DomainQuery>): DomainQuery {
    const message = createBaseDomainQuery();
    message.domain = object.domain ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseDomainGeneratedQuery(): DomainGeneratedQuery {
  return { generated: false };
}

export const DomainGeneratedQuery = {
  encode(message: DomainGeneratedQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.generated === true) {
      writer.uint32(8).bool(message.generated);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DomainGeneratedQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDomainGeneratedQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.generated = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DomainGeneratedQuery {
    return { generated: isSet(object.generated) ? Boolean(object.generated) : false };
  },

  toJSON(message: DomainGeneratedQuery): unknown {
    const obj: any = {};
    message.generated !== undefined && (obj.generated = message.generated);
    return obj;
  },

  create(base?: DeepPartial<DomainGeneratedQuery>): DomainGeneratedQuery {
    return DomainGeneratedQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DomainGeneratedQuery>): DomainGeneratedQuery {
    const message = createBaseDomainGeneratedQuery();
    message.generated = object.generated ?? false;
    return message;
  },
};

function createBaseDomainPrimaryQuery(): DomainPrimaryQuery {
  return { primary: false };
}

export const DomainPrimaryQuery = {
  encode(message: DomainPrimaryQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.primary === true) {
      writer.uint32(8).bool(message.primary);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DomainPrimaryQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDomainPrimaryQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.primary = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DomainPrimaryQuery {
    return { primary: isSet(object.primary) ? Boolean(object.primary) : false };
  },

  toJSON(message: DomainPrimaryQuery): unknown {
    const obj: any = {};
    message.primary !== undefined && (obj.primary = message.primary);
    return obj;
  },

  create(base?: DeepPartial<DomainPrimaryQuery>): DomainPrimaryQuery {
    return DomainPrimaryQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DomainPrimaryQuery>): DomainPrimaryQuery {
    const message = createBaseDomainPrimaryQuery();
    message.primary = object.primary ?? false;
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

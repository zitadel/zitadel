import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class Instance extends jspb.Message {
  getId(): string;
  setId(value: string): Instance;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Instance;
  hasDetails(): boolean;
  clearDetails(): Instance;

  getState(): State;
  setState(value: State): Instance;

  getName(): string;
  setName(value: string): Instance;

  getVersion(): string;
  setVersion(value: string): Instance;

  getDomainsList(): Array<Domain>;
  setDomainsList(value: Array<Domain>): Instance;
  clearDomainsList(): Instance;
  addDomains(value?: Domain, index?: number): Domain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Instance.AsObject;
  static toObject(includeInstance: boolean, msg: Instance): Instance.AsObject;
  static serializeBinaryToWriter(message: Instance, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Instance;
  static deserializeBinaryFromReader(message: Instance, reader: jspb.BinaryReader): Instance;
}

export namespace Instance {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: State,
    name: string,
    version: string,
    domainsList: Array<Domain.AsObject>,
  }
}

export class InstanceDetail extends jspb.Message {
  getId(): string;
  setId(value: string): InstanceDetail;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): InstanceDetail;
  hasDetails(): boolean;
  clearDetails(): InstanceDetail;

  getState(): State;
  setState(value: State): InstanceDetail;

  getName(): string;
  setName(value: string): InstanceDetail;

  getVersion(): string;
  setVersion(value: string): InstanceDetail;

  getDomainsList(): Array<Domain>;
  setDomainsList(value: Array<Domain>): InstanceDetail;
  clearDomainsList(): InstanceDetail;
  addDomains(value?: Domain, index?: number): Domain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InstanceDetail.AsObject;
  static toObject(includeInstance: boolean, msg: InstanceDetail): InstanceDetail.AsObject;
  static serializeBinaryToWriter(message: InstanceDetail, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InstanceDetail;
  static deserializeBinaryFromReader(message: InstanceDetail, reader: jspb.BinaryReader): InstanceDetail;
}

export namespace InstanceDetail {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: State,
    name: string,
    version: string,
    domainsList: Array<Domain.AsObject>,
  }
}

export class Query extends jspb.Message {
  getIdQuery(): IdsQuery | undefined;
  setIdQuery(value?: IdsQuery): Query;
  hasIdQuery(): boolean;
  clearIdQuery(): Query;

  getDomainQuery(): DomainsQuery | undefined;
  setDomainQuery(value?: DomainsQuery): Query;
  hasDomainQuery(): boolean;
  clearDomainQuery(): Query;

  getQueryCase(): Query.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Query.AsObject;
  static toObject(includeInstance: boolean, msg: Query): Query.AsObject;
  static serializeBinaryToWriter(message: Query, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Query;
  static deserializeBinaryFromReader(message: Query, reader: jspb.BinaryReader): Query;
}

export namespace Query {
  export type AsObject = {
    idQuery?: IdsQuery.AsObject,
    domainQuery?: DomainsQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    ID_QUERY = 1,
    DOMAIN_QUERY = 2,
  }
}

export class IdsQuery extends jspb.Message {
  getIdsList(): Array<string>;
  setIdsList(value: Array<string>): IdsQuery;
  clearIdsList(): IdsQuery;
  addIds(value: string, index?: number): IdsQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IdsQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IdsQuery): IdsQuery.AsObject;
  static serializeBinaryToWriter(message: IdsQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IdsQuery;
  static deserializeBinaryFromReader(message: IdsQuery, reader: jspb.BinaryReader): IdsQuery;
}

export namespace IdsQuery {
  export type AsObject = {
    idsList: Array<string>,
  }
}

export class DomainsQuery extends jspb.Message {
  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): DomainsQuery;
  clearDomainsList(): DomainsQuery;
  addDomains(value: string, index?: number): DomainsQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainsQuery.AsObject;
  static toObject(includeInstance: boolean, msg: DomainsQuery): DomainsQuery.AsObject;
  static serializeBinaryToWriter(message: DomainsQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainsQuery;
  static deserializeBinaryFromReader(message: DomainsQuery, reader: jspb.BinaryReader): DomainsQuery;
}

export namespace DomainsQuery {
  export type AsObject = {
    domainsList: Array<string>,
  }
}

export class Domain extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Domain;
  hasDetails(): boolean;
  clearDetails(): Domain;

  getDomain(): string;
  setDomain(value: string): Domain;

  getPrimary(): boolean;
  setPrimary(value: boolean): Domain;

  getGenerated(): boolean;
  setGenerated(value: boolean): Domain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Domain.AsObject;
  static toObject(includeInstance: boolean, msg: Domain): Domain.AsObject;
  static serializeBinaryToWriter(message: Domain, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Domain;
  static deserializeBinaryFromReader(message: Domain, reader: jspb.BinaryReader): Domain;
}

export namespace Domain {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    domain: string,
    primary: boolean,
    generated: boolean,
  }
}

export class DomainSearchQuery extends jspb.Message {
  getDomainQuery(): DomainQuery | undefined;
  setDomainQuery(value?: DomainQuery): DomainSearchQuery;
  hasDomainQuery(): boolean;
  clearDomainQuery(): DomainSearchQuery;

  getGeneratedQuery(): DomainGeneratedQuery | undefined;
  setGeneratedQuery(value?: DomainGeneratedQuery): DomainSearchQuery;
  hasGeneratedQuery(): boolean;
  clearGeneratedQuery(): DomainSearchQuery;

  getPrimaryQuery(): DomainPrimaryQuery | undefined;
  setPrimaryQuery(value?: DomainPrimaryQuery): DomainSearchQuery;
  hasPrimaryQuery(): boolean;
  clearPrimaryQuery(): DomainSearchQuery;

  getQueryCase(): DomainSearchQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainSearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: DomainSearchQuery): DomainSearchQuery.AsObject;
  static serializeBinaryToWriter(message: DomainSearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainSearchQuery;
  static deserializeBinaryFromReader(message: DomainSearchQuery, reader: jspb.BinaryReader): DomainSearchQuery;
}

export namespace DomainSearchQuery {
  export type AsObject = {
    domainQuery?: DomainQuery.AsObject,
    generatedQuery?: DomainGeneratedQuery.AsObject,
    primaryQuery?: DomainPrimaryQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    DOMAIN_QUERY = 1,
    GENERATED_QUERY = 2,
    PRIMARY_QUERY = 3,
  }
}

export class DomainQuery extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): DomainQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): DomainQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainQuery.AsObject;
  static toObject(includeInstance: boolean, msg: DomainQuery): DomainQuery.AsObject;
  static serializeBinaryToWriter(message: DomainQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainQuery;
  static deserializeBinaryFromReader(message: DomainQuery, reader: jspb.BinaryReader): DomainQuery;
}

export namespace DomainQuery {
  export type AsObject = {
    domain: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class DomainGeneratedQuery extends jspb.Message {
  getGenerated(): boolean;
  setGenerated(value: boolean): DomainGeneratedQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainGeneratedQuery.AsObject;
  static toObject(includeInstance: boolean, msg: DomainGeneratedQuery): DomainGeneratedQuery.AsObject;
  static serializeBinaryToWriter(message: DomainGeneratedQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainGeneratedQuery;
  static deserializeBinaryFromReader(message: DomainGeneratedQuery, reader: jspb.BinaryReader): DomainGeneratedQuery;
}

export namespace DomainGeneratedQuery {
  export type AsObject = {
    generated: boolean,
  }
}

export class DomainPrimaryQuery extends jspb.Message {
  getPrimary(): boolean;
  setPrimary(value: boolean): DomainPrimaryQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainPrimaryQuery.AsObject;
  static toObject(includeInstance: boolean, msg: DomainPrimaryQuery): DomainPrimaryQuery.AsObject;
  static serializeBinaryToWriter(message: DomainPrimaryQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainPrimaryQuery;
  static deserializeBinaryFromReader(message: DomainPrimaryQuery, reader: jspb.BinaryReader): DomainPrimaryQuery;
}

export namespace DomainPrimaryQuery {
  export type AsObject = {
    primary: boolean,
  }
}

export enum State { 
  STATE_UNSPECIFIED = 0,
  STATE_CREATING = 1,
  STATE_RUNNING = 2,
  STATE_STOPPING = 3,
  STATE_STOPPED = 4,
}
export enum FieldName { 
  FIELD_NAME_UNSPECIFIED = 0,
  FIELD_NAME_ID = 1,
  FIELD_NAME_NAME = 2,
  FIELD_NAME_CREATION_DATE = 3,
}
export enum DomainFieldName { 
  DOMAIN_FIELD_NAME_UNSPECIFIED = 0,
  DOMAIN_FIELD_NAME_DOMAIN = 1,
  DOMAIN_FIELD_NAME_PRIMARY = 2,
  DOMAIN_FIELD_NAME_GENERATED = 3,
  DOMAIN_FIELD_NAME_CREATION_DATE = 4,
}

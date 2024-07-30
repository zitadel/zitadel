import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class Org extends jspb.Message {
  getId(): string;
  setId(value: string): Org;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Org;
  hasDetails(): boolean;
  clearDetails(): Org;

  getState(): OrgState;
  setState(value: OrgState): Org;

  getName(): string;
  setName(value: string): Org;

  getPrimaryDomain(): string;
  setPrimaryDomain(value: string): Org;

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
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    state: OrgState,
    name: string,
    primaryDomain: string,
  }
}

export class Domain extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): Domain;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Domain;
  hasDetails(): boolean;
  clearDetails(): Domain;

  getDomainName(): string;
  setDomainName(value: string): Domain;

  getIsVerified(): boolean;
  setIsVerified(value: boolean): Domain;

  getIsPrimary(): boolean;
  setIsPrimary(value: boolean): Domain;

  getValidationType(): DomainValidationType;
  setValidationType(value: DomainValidationType): Domain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Domain.AsObject;
  static toObject(includeInstance: boolean, msg: Domain): Domain.AsObject;
  static serializeBinaryToWriter(message: Domain, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Domain;
  static deserializeBinaryFromReader(message: Domain, reader: jspb.BinaryReader): Domain;
}

export namespace Domain {
  export type AsObject = {
    orgId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    domainName: string,
    isVerified: boolean,
    isPrimary: boolean,
    validationType: DomainValidationType,
  }
}

export class OrgQuery extends jspb.Message {
  getNameQuery(): OrgNameQuery | undefined;
  setNameQuery(value?: OrgNameQuery): OrgQuery;
  hasNameQuery(): boolean;
  clearNameQuery(): OrgQuery;

  getDomainQuery(): OrgDomainQuery | undefined;
  setDomainQuery(value?: OrgDomainQuery): OrgQuery;
  hasDomainQuery(): boolean;
  clearDomainQuery(): OrgQuery;

  getStateQuery(): OrgStateQuery | undefined;
  setStateQuery(value?: OrgStateQuery): OrgQuery;
  hasStateQuery(): boolean;
  clearStateQuery(): OrgQuery;

  getQueryCase(): OrgQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrgQuery): OrgQuery.AsObject;
  static serializeBinaryToWriter(message: OrgQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgQuery;
  static deserializeBinaryFromReader(message: OrgQuery, reader: jspb.BinaryReader): OrgQuery;
}

export namespace OrgQuery {
  export type AsObject = {
    nameQuery?: OrgNameQuery.AsObject,
    domainQuery?: OrgDomainQuery.AsObject,
    stateQuery?: OrgStateQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    NAME_QUERY = 1,
    DOMAIN_QUERY = 2,
    STATE_QUERY = 3,
  }
}

export class OrgNameQuery extends jspb.Message {
  getName(): string;
  setName(value: string): OrgNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): OrgNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrgNameQuery): OrgNameQuery.AsObject;
  static serializeBinaryToWriter(message: OrgNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgNameQuery;
  static deserializeBinaryFromReader(message: OrgNameQuery, reader: jspb.BinaryReader): OrgNameQuery;
}

export namespace OrgNameQuery {
  export type AsObject = {
    name: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class OrgDomainQuery extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): OrgDomainQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): OrgDomainQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgDomainQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrgDomainQuery): OrgDomainQuery.AsObject;
  static serializeBinaryToWriter(message: OrgDomainQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgDomainQuery;
  static deserializeBinaryFromReader(message: OrgDomainQuery, reader: jspb.BinaryReader): OrgDomainQuery;
}

export namespace OrgDomainQuery {
  export type AsObject = {
    domain: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class OrgStateQuery extends jspb.Message {
  getState(): OrgState;
  setState(value: OrgState): OrgStateQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgStateQuery.AsObject;
  static toObject(includeInstance: boolean, msg: OrgStateQuery): OrgStateQuery.AsObject;
  static serializeBinaryToWriter(message: OrgStateQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgStateQuery;
  static deserializeBinaryFromReader(message: OrgStateQuery, reader: jspb.BinaryReader): OrgStateQuery;
}

export namespace OrgStateQuery {
  export type AsObject = {
    state: OrgState,
  }
}

export class DomainSearchQuery extends jspb.Message {
  getDomainNameQuery(): DomainNameQuery | undefined;
  setDomainNameQuery(value?: DomainNameQuery): DomainSearchQuery;
  hasDomainNameQuery(): boolean;
  clearDomainNameQuery(): DomainSearchQuery;

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
    domainNameQuery?: DomainNameQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    DOMAIN_NAME_QUERY = 1,
  }
}

export class DomainNameQuery extends jspb.Message {
  getName(): string;
  setName(value: string): DomainNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): DomainNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: DomainNameQuery): DomainNameQuery.AsObject;
  static serializeBinaryToWriter(message: DomainNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainNameQuery;
  static deserializeBinaryFromReader(message: DomainNameQuery, reader: jspb.BinaryReader): DomainNameQuery;
}

export namespace DomainNameQuery {
  export type AsObject = {
    name: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export enum OrgState { 
  ORG_STATE_UNSPECIFIED = 0,
  ORG_STATE_ACTIVE = 1,
  ORG_STATE_INACTIVE = 2,
  ORG_STATE_REMOVED = 3,
}
export enum DomainValidationType { 
  DOMAIN_VALIDATION_TYPE_UNSPECIFIED = 0,
  DOMAIN_VALIDATION_TYPE_HTTP = 1,
  DOMAIN_VALIDATION_TYPE_DNS = 2,
}
export enum OrgFieldName { 
  ORG_FIELD_NAME_UNSPECIFIED = 0,
  ORG_FIELD_NAME_NAME = 1,
}

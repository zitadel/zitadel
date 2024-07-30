import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../../../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class Session extends jspb.Message {
  getId(): string;
  setId(value: string): Session;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): Session;
  hasCreationDate(): boolean;
  clearCreationDate(): Session;

  getChangeDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setChangeDate(value?: google_protobuf_timestamp_pb.Timestamp): Session;
  hasChangeDate(): boolean;
  clearChangeDate(): Session;

  getSequence(): number;
  setSequence(value: number): Session;

  getFactors(): Factors | undefined;
  setFactors(value?: Factors): Session;
  hasFactors(): boolean;
  clearFactors(): Session;

  getMetadataMap(): jspb.Map<string, Uint8Array | string>;
  clearMetadataMap(): Session;

  getUserAgent(): UserAgent | undefined;
  setUserAgent(value?: UserAgent): Session;
  hasUserAgent(): boolean;
  clearUserAgent(): Session;

  getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): Session;
  hasExpirationDate(): boolean;
  clearExpirationDate(): Session;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Session.AsObject;
  static toObject(includeInstance: boolean, msg: Session): Session.AsObject;
  static serializeBinaryToWriter(message: Session, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Session;
  static deserializeBinaryFromReader(message: Session, reader: jspb.BinaryReader): Session;
}

export namespace Session {
  export type AsObject = {
    id: string,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    changeDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    sequence: number,
    factors?: Factors.AsObject,
    metadataMap: Array<[string, Uint8Array | string]>,
    userAgent?: UserAgent.AsObject,
    expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }

  export enum ExpirationDateCase { 
    _EXPIRATION_DATE_NOT_SET = 0,
    EXPIRATION_DATE = 8,
  }
}

export class Factors extends jspb.Message {
  getUser(): UserFactor | undefined;
  setUser(value?: UserFactor): Factors;
  hasUser(): boolean;
  clearUser(): Factors;

  getPassword(): PasswordFactor | undefined;
  setPassword(value?: PasswordFactor): Factors;
  hasPassword(): boolean;
  clearPassword(): Factors;

  getWebAuthN(): WebAuthNFactor | undefined;
  setWebAuthN(value?: WebAuthNFactor): Factors;
  hasWebAuthN(): boolean;
  clearWebAuthN(): Factors;

  getIntent(): IntentFactor | undefined;
  setIntent(value?: IntentFactor): Factors;
  hasIntent(): boolean;
  clearIntent(): Factors;

  getTotp(): TOTPFactor | undefined;
  setTotp(value?: TOTPFactor): Factors;
  hasTotp(): boolean;
  clearTotp(): Factors;

  getOtpSms(): OTPFactor | undefined;
  setOtpSms(value?: OTPFactor): Factors;
  hasOtpSms(): boolean;
  clearOtpSms(): Factors;

  getOtpEmail(): OTPFactor | undefined;
  setOtpEmail(value?: OTPFactor): Factors;
  hasOtpEmail(): boolean;
  clearOtpEmail(): Factors;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Factors.AsObject;
  static toObject(includeInstance: boolean, msg: Factors): Factors.AsObject;
  static serializeBinaryToWriter(message: Factors, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Factors;
  static deserializeBinaryFromReader(message: Factors, reader: jspb.BinaryReader): Factors;
}

export namespace Factors {
  export type AsObject = {
    user?: UserFactor.AsObject,
    password?: PasswordFactor.AsObject,
    webAuthN?: WebAuthNFactor.AsObject,
    intent?: IntentFactor.AsObject,
    totp?: TOTPFactor.AsObject,
    otpSms?: OTPFactor.AsObject,
    otpEmail?: OTPFactor.AsObject,
  }
}

export class UserFactor extends jspb.Message {
  getVerifiedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setVerifiedAt(value?: google_protobuf_timestamp_pb.Timestamp): UserFactor;
  hasVerifiedAt(): boolean;
  clearVerifiedAt(): UserFactor;

  getId(): string;
  setId(value: string): UserFactor;

  getLoginName(): string;
  setLoginName(value: string): UserFactor;

  getDisplayName(): string;
  setDisplayName(value: string): UserFactor;

  getOrganizationId(): string;
  setOrganizationId(value: string): UserFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserFactor.AsObject;
  static toObject(includeInstance: boolean, msg: UserFactor): UserFactor.AsObject;
  static serializeBinaryToWriter(message: UserFactor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserFactor;
  static deserializeBinaryFromReader(message: UserFactor, reader: jspb.BinaryReader): UserFactor;
}

export namespace UserFactor {
  export type AsObject = {
    verifiedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    id: string,
    loginName: string,
    displayName: string,
    organizationId: string,
  }
}

export class PasswordFactor extends jspb.Message {
  getVerifiedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setVerifiedAt(value?: google_protobuf_timestamp_pb.Timestamp): PasswordFactor;
  hasVerifiedAt(): boolean;
  clearVerifiedAt(): PasswordFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordFactor.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordFactor): PasswordFactor.AsObject;
  static serializeBinaryToWriter(message: PasswordFactor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordFactor;
  static deserializeBinaryFromReader(message: PasswordFactor, reader: jspb.BinaryReader): PasswordFactor;
}

export namespace PasswordFactor {
  export type AsObject = {
    verifiedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class IntentFactor extends jspb.Message {
  getVerifiedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setVerifiedAt(value?: google_protobuf_timestamp_pb.Timestamp): IntentFactor;
  hasVerifiedAt(): boolean;
  clearVerifiedAt(): IntentFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IntentFactor.AsObject;
  static toObject(includeInstance: boolean, msg: IntentFactor): IntentFactor.AsObject;
  static serializeBinaryToWriter(message: IntentFactor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IntentFactor;
  static deserializeBinaryFromReader(message: IntentFactor, reader: jspb.BinaryReader): IntentFactor;
}

export namespace IntentFactor {
  export type AsObject = {
    verifiedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class WebAuthNFactor extends jspb.Message {
  getVerifiedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setVerifiedAt(value?: google_protobuf_timestamp_pb.Timestamp): WebAuthNFactor;
  hasVerifiedAt(): boolean;
  clearVerifiedAt(): WebAuthNFactor;

  getUserVerified(): boolean;
  setUserVerified(value: boolean): WebAuthNFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebAuthNFactor.AsObject;
  static toObject(includeInstance: boolean, msg: WebAuthNFactor): WebAuthNFactor.AsObject;
  static serializeBinaryToWriter(message: WebAuthNFactor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebAuthNFactor;
  static deserializeBinaryFromReader(message: WebAuthNFactor, reader: jspb.BinaryReader): WebAuthNFactor;
}

export namespace WebAuthNFactor {
  export type AsObject = {
    verifiedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    userVerified: boolean,
  }
}

export class TOTPFactor extends jspb.Message {
  getVerifiedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setVerifiedAt(value?: google_protobuf_timestamp_pb.Timestamp): TOTPFactor;
  hasVerifiedAt(): boolean;
  clearVerifiedAt(): TOTPFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TOTPFactor.AsObject;
  static toObject(includeInstance: boolean, msg: TOTPFactor): TOTPFactor.AsObject;
  static serializeBinaryToWriter(message: TOTPFactor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TOTPFactor;
  static deserializeBinaryFromReader(message: TOTPFactor, reader: jspb.BinaryReader): TOTPFactor;
}

export namespace TOTPFactor {
  export type AsObject = {
    verifiedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class OTPFactor extends jspb.Message {
  getVerifiedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setVerifiedAt(value?: google_protobuf_timestamp_pb.Timestamp): OTPFactor;
  hasVerifiedAt(): boolean;
  clearVerifiedAt(): OTPFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OTPFactor.AsObject;
  static toObject(includeInstance: boolean, msg: OTPFactor): OTPFactor.AsObject;
  static serializeBinaryToWriter(message: OTPFactor, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OTPFactor;
  static deserializeBinaryFromReader(message: OTPFactor, reader: jspb.BinaryReader): OTPFactor;
}

export namespace OTPFactor {
  export type AsObject = {
    verifiedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class SearchQuery extends jspb.Message {
  getIdsQuery(): IDsQuery | undefined;
  setIdsQuery(value?: IDsQuery): SearchQuery;
  hasIdsQuery(): boolean;
  clearIdsQuery(): SearchQuery;

  getUserIdQuery(): UserIDQuery | undefined;
  setUserIdQuery(value?: UserIDQuery): SearchQuery;
  hasUserIdQuery(): boolean;
  clearUserIdQuery(): SearchQuery;

  getCreationDateQuery(): CreationDateQuery | undefined;
  setCreationDateQuery(value?: CreationDateQuery): SearchQuery;
  hasCreationDateQuery(): boolean;
  clearCreationDateQuery(): SearchQuery;

  getQueryCase(): SearchQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SearchQuery.AsObject;
  static toObject(includeInstance: boolean, msg: SearchQuery): SearchQuery.AsObject;
  static serializeBinaryToWriter(message: SearchQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SearchQuery;
  static deserializeBinaryFromReader(message: SearchQuery, reader: jspb.BinaryReader): SearchQuery;
}

export namespace SearchQuery {
  export type AsObject = {
    idsQuery?: IDsQuery.AsObject,
    userIdQuery?: UserIDQuery.AsObject,
    creationDateQuery?: CreationDateQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    IDS_QUERY = 1,
    USER_ID_QUERY = 2,
    CREATION_DATE_QUERY = 3,
  }
}

export class IDsQuery extends jspb.Message {
  getIdsList(): Array<string>;
  setIdsList(value: Array<string>): IDsQuery;
  clearIdsList(): IDsQuery;
  addIds(value: string, index?: number): IDsQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDsQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IDsQuery): IDsQuery.AsObject;
  static serializeBinaryToWriter(message: IDsQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDsQuery;
  static deserializeBinaryFromReader(message: IDsQuery, reader: jspb.BinaryReader): IDsQuery;
}

export namespace IDsQuery {
  export type AsObject = {
    idsList: Array<string>,
  }
}

export class UserIDQuery extends jspb.Message {
  getId(): string;
  setId(value: string): UserIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: UserIDQuery): UserIDQuery.AsObject;
  static serializeBinaryToWriter(message: UserIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserIDQuery;
  static deserializeBinaryFromReader(message: UserIDQuery, reader: jspb.BinaryReader): UserIDQuery;
}

export namespace UserIDQuery {
  export type AsObject = {
    id: string,
  }
}

export class CreationDateQuery extends jspb.Message {
  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): CreationDateQuery;
  hasCreationDate(): boolean;
  clearCreationDate(): CreationDateQuery;

  getMethod(): zitadel_object_pb.TimestampQueryMethod;
  setMethod(value: zitadel_object_pb.TimestampQueryMethod): CreationDateQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreationDateQuery.AsObject;
  static toObject(includeInstance: boolean, msg: CreationDateQuery): CreationDateQuery.AsObject;
  static serializeBinaryToWriter(message: CreationDateQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreationDateQuery;
  static deserializeBinaryFromReader(message: CreationDateQuery, reader: jspb.BinaryReader): CreationDateQuery;
}

export namespace CreationDateQuery {
  export type AsObject = {
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    method: zitadel_object_pb.TimestampQueryMethod,
  }
}

export class UserAgent extends jspb.Message {
  getFingerprintId(): string;
  setFingerprintId(value: string): UserAgent;
  hasFingerprintId(): boolean;
  clearFingerprintId(): UserAgent;

  getIp(): string;
  setIp(value: string): UserAgent;
  hasIp(): boolean;
  clearIp(): UserAgent;

  getDescription(): string;
  setDescription(value: string): UserAgent;
  hasDescription(): boolean;
  clearDescription(): UserAgent;

  getHeaderMap(): jspb.Map<string, UserAgent.HeaderValues>;
  clearHeaderMap(): UserAgent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserAgent.AsObject;
  static toObject(includeInstance: boolean, msg: UserAgent): UserAgent.AsObject;
  static serializeBinaryToWriter(message: UserAgent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserAgent;
  static deserializeBinaryFromReader(message: UserAgent, reader: jspb.BinaryReader): UserAgent;
}

export namespace UserAgent {
  export type AsObject = {
    fingerprintId?: string,
    ip?: string,
    description?: string,
    headerMap: Array<[string, UserAgent.HeaderValues.AsObject]>,
  }

  export class HeaderValues extends jspb.Message {
    getValuesList(): Array<string>;
    setValuesList(value: Array<string>): HeaderValues;
    clearValuesList(): HeaderValues;
    addValues(value: string, index?: number): HeaderValues;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HeaderValues.AsObject;
    static toObject(includeInstance: boolean, msg: HeaderValues): HeaderValues.AsObject;
    static serializeBinaryToWriter(message: HeaderValues, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HeaderValues;
    static deserializeBinaryFromReader(message: HeaderValues, reader: jspb.BinaryReader): HeaderValues;
  }

  export namespace HeaderValues {
    export type AsObject = {
      valuesList: Array<string>,
    }
  }


  export enum FingerprintIdCase { 
    _FINGERPRINT_ID_NOT_SET = 0,
    FINGERPRINT_ID = 1,
  }

  export enum IpCase { 
    _IP_NOT_SET = 0,
    IP = 2,
  }

  export enum DescriptionCase { 
    _DESCRIPTION_NOT_SET = 0,
    DESCRIPTION = 3,
  }
}

export enum SessionFieldName { 
  SESSION_FIELD_NAME_UNSPECIFIED = 0,
  SESSION_FIELD_NAME_CREATION_DATE = 1,
}

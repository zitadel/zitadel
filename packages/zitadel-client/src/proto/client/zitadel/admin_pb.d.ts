import * as jspb from 'google-protobuf'

import * as zitadel_idp_pb from '../zitadel/idp_pb'; // proto import: "zitadel/idp.proto"
import * as zitadel_instance_pb from '../zitadel/instance_pb'; // proto import: "zitadel/instance.proto"
import * as zitadel_user_pb from '../zitadel/user_pb'; // proto import: "zitadel/user.proto"
import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as zitadel_options_pb from '../zitadel/options_pb'; // proto import: "zitadel/options.proto"
import * as zitadel_org_pb from '../zitadel/org_pb'; // proto import: "zitadel/org.proto"
import * as zitadel_policy_pb from '../zitadel/policy_pb'; // proto import: "zitadel/policy.proto"
import * as zitadel_settings_pb from '../zitadel/settings_pb'; // proto import: "zitadel/settings.proto"
import * as zitadel_text_pb from '../zitadel/text_pb'; // proto import: "zitadel/text.proto"
import * as zitadel_member_pb from '../zitadel/member_pb'; // proto import: "zitadel/member.proto"
import * as zitadel_event_pb from '../zitadel/event_pb'; // proto import: "zitadel/event.proto"
import * as zitadel_management_pb from '../zitadel/management_pb'; // proto import: "zitadel/management.proto"
import * as zitadel_v1_pb from '../zitadel/v1_pb'; // proto import: "zitadel/v1.proto"
import * as zitadel_message_pb from '../zitadel/message_pb'; // proto import: "zitadel/message.proto"
import * as zitadel_milestone_v1_milestone_pb from '../zitadel/milestone/v1/milestone_pb'; // proto import: "zitadel/milestone/v1/milestone.proto"
import * as google_api_annotations_pb from '../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"


export class HealthzRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthzRequest.AsObject;
  static toObject(includeInstance: boolean, msg: HealthzRequest): HealthzRequest.AsObject;
  static serializeBinaryToWriter(message: HealthzRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthzRequest;
  static deserializeBinaryFromReader(message: HealthzRequest, reader: jspb.BinaryReader): HealthzRequest;
}

export namespace HealthzRequest {
  export type AsObject = {
  }
}

export class HealthzResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthzResponse.AsObject;
  static toObject(includeInstance: boolean, msg: HealthzResponse): HealthzResponse.AsObject;
  static serializeBinaryToWriter(message: HealthzResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthzResponse;
  static deserializeBinaryFromReader(message: HealthzResponse, reader: jspb.BinaryReader): HealthzResponse;
}

export namespace HealthzResponse {
  export type AsObject = {
  }
}

export class GetSupportedLanguagesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupportedLanguagesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupportedLanguagesRequest): GetSupportedLanguagesRequest.AsObject;
  static serializeBinaryToWriter(message: GetSupportedLanguagesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupportedLanguagesRequest;
  static deserializeBinaryFromReader(message: GetSupportedLanguagesRequest, reader: jspb.BinaryReader): GetSupportedLanguagesRequest;
}

export namespace GetSupportedLanguagesRequest {
  export type AsObject = {
  }
}

export class GetSupportedLanguagesResponse extends jspb.Message {
  getLanguagesList(): Array<string>;
  setLanguagesList(value: Array<string>): GetSupportedLanguagesResponse;
  clearLanguagesList(): GetSupportedLanguagesResponse;
  addLanguages(value: string, index?: number): GetSupportedLanguagesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupportedLanguagesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupportedLanguagesResponse): GetSupportedLanguagesResponse.AsObject;
  static serializeBinaryToWriter(message: GetSupportedLanguagesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupportedLanguagesResponse;
  static deserializeBinaryFromReader(message: GetSupportedLanguagesResponse, reader: jspb.BinaryReader): GetSupportedLanguagesResponse;
}

export namespace GetSupportedLanguagesResponse {
  export type AsObject = {
    languagesList: Array<string>,
  }
}

export class GetAllowedLanguagesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllowedLanguagesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllowedLanguagesRequest): GetAllowedLanguagesRequest.AsObject;
  static serializeBinaryToWriter(message: GetAllowedLanguagesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllowedLanguagesRequest;
  static deserializeBinaryFromReader(message: GetAllowedLanguagesRequest, reader: jspb.BinaryReader): GetAllowedLanguagesRequest;
}

export namespace GetAllowedLanguagesRequest {
  export type AsObject = {
  }
}

export class GetAllowedLanguagesResponse extends jspb.Message {
  getLanguagesList(): Array<string>;
  setLanguagesList(value: Array<string>): GetAllowedLanguagesResponse;
  clearLanguagesList(): GetAllowedLanguagesResponse;
  addLanguages(value: string, index?: number): GetAllowedLanguagesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAllowedLanguagesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAllowedLanguagesResponse): GetAllowedLanguagesResponse.AsObject;
  static serializeBinaryToWriter(message: GetAllowedLanguagesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAllowedLanguagesResponse;
  static deserializeBinaryFromReader(message: GetAllowedLanguagesResponse, reader: jspb.BinaryReader): GetAllowedLanguagesResponse;
}

export namespace GetAllowedLanguagesResponse {
  export type AsObject = {
    languagesList: Array<string>,
  }
}

export class SetDefaultLanguageRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultLanguageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultLanguageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultLanguageRequest): SetDefaultLanguageRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultLanguageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultLanguageRequest;
  static deserializeBinaryFromReader(message: SetDefaultLanguageRequest, reader: jspb.BinaryReader): SetDefaultLanguageRequest;
}

export namespace SetDefaultLanguageRequest {
  export type AsObject = {
    language: string,
  }
}

export class SetDefaultLanguageResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultLanguageResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultLanguageResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultLanguageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultLanguageResponse): SetDefaultLanguageResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultLanguageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultLanguageResponse;
  static deserializeBinaryFromReader(message: SetDefaultLanguageResponse, reader: jspb.BinaryReader): SetDefaultLanguageResponse;
}

export namespace SetDefaultLanguageResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultLanguageRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLanguageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLanguageRequest): GetDefaultLanguageRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLanguageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLanguageRequest;
  static deserializeBinaryFromReader(message: GetDefaultLanguageRequest, reader: jspb.BinaryReader): GetDefaultLanguageRequest;
}

export namespace GetDefaultLanguageRequest {
  export type AsObject = {
  }
}

export class GetDefaultLanguageResponse extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultLanguageResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLanguageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLanguageResponse): GetDefaultLanguageResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLanguageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLanguageResponse;
  static deserializeBinaryFromReader(message: GetDefaultLanguageResponse, reader: jspb.BinaryReader): GetDefaultLanguageResponse;
}

export namespace GetDefaultLanguageResponse {
  export type AsObject = {
    language: string,
  }
}

export class SetDefaultOrgRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): SetDefaultOrgRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultOrgRequest): SetDefaultOrgRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultOrgRequest;
  static deserializeBinaryFromReader(message: SetDefaultOrgRequest, reader: jspb.BinaryReader): SetDefaultOrgRequest;
}

export namespace SetDefaultOrgRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class SetDefaultOrgResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultOrgResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultOrgResponse): SetDefaultOrgResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultOrgResponse;
  static deserializeBinaryFromReader(message: SetDefaultOrgResponse, reader: jspb.BinaryReader): SetDefaultOrgResponse;
}

export namespace SetDefaultOrgResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultOrgRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultOrgRequest): GetDefaultOrgRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultOrgRequest;
  static deserializeBinaryFromReader(message: GetDefaultOrgRequest, reader: jspb.BinaryReader): GetDefaultOrgRequest;
}

export namespace GetDefaultOrgRequest {
  export type AsObject = {
  }
}

export class GetDefaultOrgResponse extends jspb.Message {
  getOrg(): zitadel_org_pb.Org | undefined;
  setOrg(value?: zitadel_org_pb.Org): GetDefaultOrgResponse;
  hasOrg(): boolean;
  clearOrg(): GetDefaultOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultOrgResponse): GetDefaultOrgResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultOrgResponse;
  static deserializeBinaryFromReader(message: GetDefaultOrgResponse, reader: jspb.BinaryReader): GetDefaultOrgResponse;
}

export namespace GetDefaultOrgResponse {
  export type AsObject = {
    org?: zitadel_org_pb.Org.AsObject,
  }
}

export class GetMyInstanceRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyInstanceRequest): GetMyInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyInstanceRequest;
  static deserializeBinaryFromReader(message: GetMyInstanceRequest, reader: jspb.BinaryReader): GetMyInstanceRequest;
}

export namespace GetMyInstanceRequest {
  export type AsObject = {
  }
}

export class GetMyInstanceResponse extends jspb.Message {
  getInstance(): zitadel_instance_pb.InstanceDetail | undefined;
  setInstance(value?: zitadel_instance_pb.InstanceDetail): GetMyInstanceResponse;
  hasInstance(): boolean;
  clearInstance(): GetMyInstanceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyInstanceResponse): GetMyInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyInstanceResponse;
  static deserializeBinaryFromReader(message: GetMyInstanceResponse, reader: jspb.BinaryReader): GetMyInstanceResponse;
}

export namespace GetMyInstanceResponse {
  export type AsObject = {
    instance?: zitadel_instance_pb.InstanceDetail.AsObject,
  }
}

export class ListInstanceDomainsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListInstanceDomainsRequest;
  hasQuery(): boolean;
  clearQuery(): ListInstanceDomainsRequest;

  getSortingColumn(): zitadel_instance_pb.DomainFieldName;
  setSortingColumn(value: zitadel_instance_pb.DomainFieldName): ListInstanceDomainsRequest;

  getQueriesList(): Array<zitadel_instance_pb.DomainSearchQuery>;
  setQueriesList(value: Array<zitadel_instance_pb.DomainSearchQuery>): ListInstanceDomainsRequest;
  clearQueriesList(): ListInstanceDomainsRequest;
  addQueries(value?: zitadel_instance_pb.DomainSearchQuery, index?: number): zitadel_instance_pb.DomainSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListInstanceDomainsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListInstanceDomainsRequest): ListInstanceDomainsRequest.AsObject;
  static serializeBinaryToWriter(message: ListInstanceDomainsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListInstanceDomainsRequest;
  static deserializeBinaryFromReader(message: ListInstanceDomainsRequest, reader: jspb.BinaryReader): ListInstanceDomainsRequest;
}

export namespace ListInstanceDomainsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_instance_pb.DomainFieldName,
    queriesList: Array<zitadel_instance_pb.DomainSearchQuery.AsObject>,
  }
}

export class ListInstanceDomainsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListInstanceDomainsResponse;
  hasDetails(): boolean;
  clearDetails(): ListInstanceDomainsResponse;

  getSortingColumn(): zitadel_instance_pb.DomainFieldName;
  setSortingColumn(value: zitadel_instance_pb.DomainFieldName): ListInstanceDomainsResponse;

  getResultList(): Array<zitadel_instance_pb.Domain>;
  setResultList(value: Array<zitadel_instance_pb.Domain>): ListInstanceDomainsResponse;
  clearResultList(): ListInstanceDomainsResponse;
  addResult(value?: zitadel_instance_pb.Domain, index?: number): zitadel_instance_pb.Domain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListInstanceDomainsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListInstanceDomainsResponse): ListInstanceDomainsResponse.AsObject;
  static serializeBinaryToWriter(message: ListInstanceDomainsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListInstanceDomainsResponse;
  static deserializeBinaryFromReader(message: ListInstanceDomainsResponse, reader: jspb.BinaryReader): ListInstanceDomainsResponse;
}

export namespace ListInstanceDomainsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_instance_pb.DomainFieldName,
    resultList: Array<zitadel_instance_pb.Domain.AsObject>,
  }
}

export class ListSecretGeneratorsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListSecretGeneratorsRequest;
  hasQuery(): boolean;
  clearQuery(): ListSecretGeneratorsRequest;

  getQueriesList(): Array<zitadel_settings_pb.SecretGeneratorQuery>;
  setQueriesList(value: Array<zitadel_settings_pb.SecretGeneratorQuery>): ListSecretGeneratorsRequest;
  clearQueriesList(): ListSecretGeneratorsRequest;
  addQueries(value?: zitadel_settings_pb.SecretGeneratorQuery, index?: number): zitadel_settings_pb.SecretGeneratorQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSecretGeneratorsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSecretGeneratorsRequest): ListSecretGeneratorsRequest.AsObject;
  static serializeBinaryToWriter(message: ListSecretGeneratorsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSecretGeneratorsRequest;
  static deserializeBinaryFromReader(message: ListSecretGeneratorsRequest, reader: jspb.BinaryReader): ListSecretGeneratorsRequest;
}

export namespace ListSecretGeneratorsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_settings_pb.SecretGeneratorQuery.AsObject>,
  }
}

export class ListSecretGeneratorsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListSecretGeneratorsResponse;
  hasDetails(): boolean;
  clearDetails(): ListSecretGeneratorsResponse;

  getResultList(): Array<zitadel_settings_pb.SecretGenerator>;
  setResultList(value: Array<zitadel_settings_pb.SecretGenerator>): ListSecretGeneratorsResponse;
  clearResultList(): ListSecretGeneratorsResponse;
  addResult(value?: zitadel_settings_pb.SecretGenerator, index?: number): zitadel_settings_pb.SecretGenerator;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSecretGeneratorsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSecretGeneratorsResponse): ListSecretGeneratorsResponse.AsObject;
  static serializeBinaryToWriter(message: ListSecretGeneratorsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSecretGeneratorsResponse;
  static deserializeBinaryFromReader(message: ListSecretGeneratorsResponse, reader: jspb.BinaryReader): ListSecretGeneratorsResponse;
}

export namespace ListSecretGeneratorsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_settings_pb.SecretGenerator.AsObject>,
  }
}

export class GetSecretGeneratorRequest extends jspb.Message {
  getGeneratorType(): zitadel_settings_pb.SecretGeneratorType;
  setGeneratorType(value: zitadel_settings_pb.SecretGeneratorType): GetSecretGeneratorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSecretGeneratorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSecretGeneratorRequest): GetSecretGeneratorRequest.AsObject;
  static serializeBinaryToWriter(message: GetSecretGeneratorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSecretGeneratorRequest;
  static deserializeBinaryFromReader(message: GetSecretGeneratorRequest, reader: jspb.BinaryReader): GetSecretGeneratorRequest;
}

export namespace GetSecretGeneratorRequest {
  export type AsObject = {
    generatorType: zitadel_settings_pb.SecretGeneratorType,
  }
}

export class GetSecretGeneratorResponse extends jspb.Message {
  getSecretGenerator(): zitadel_settings_pb.SecretGenerator | undefined;
  setSecretGenerator(value?: zitadel_settings_pb.SecretGenerator): GetSecretGeneratorResponse;
  hasSecretGenerator(): boolean;
  clearSecretGenerator(): GetSecretGeneratorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSecretGeneratorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSecretGeneratorResponse): GetSecretGeneratorResponse.AsObject;
  static serializeBinaryToWriter(message: GetSecretGeneratorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSecretGeneratorResponse;
  static deserializeBinaryFromReader(message: GetSecretGeneratorResponse, reader: jspb.BinaryReader): GetSecretGeneratorResponse;
}

export namespace GetSecretGeneratorResponse {
  export type AsObject = {
    secretGenerator?: zitadel_settings_pb.SecretGenerator.AsObject,
  }
}

export class UpdateSecretGeneratorRequest extends jspb.Message {
  getGeneratorType(): zitadel_settings_pb.SecretGeneratorType;
  setGeneratorType(value: zitadel_settings_pb.SecretGeneratorType): UpdateSecretGeneratorRequest;

  getLength(): number;
  setLength(value: number): UpdateSecretGeneratorRequest;

  getExpiry(): google_protobuf_duration_pb.Duration | undefined;
  setExpiry(value?: google_protobuf_duration_pb.Duration): UpdateSecretGeneratorRequest;
  hasExpiry(): boolean;
  clearExpiry(): UpdateSecretGeneratorRequest;

  getIncludeLowerLetters(): boolean;
  setIncludeLowerLetters(value: boolean): UpdateSecretGeneratorRequest;

  getIncludeUpperLetters(): boolean;
  setIncludeUpperLetters(value: boolean): UpdateSecretGeneratorRequest;

  getIncludeDigits(): boolean;
  setIncludeDigits(value: boolean): UpdateSecretGeneratorRequest;

  getIncludeSymbols(): boolean;
  setIncludeSymbols(value: boolean): UpdateSecretGeneratorRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSecretGeneratorRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSecretGeneratorRequest): UpdateSecretGeneratorRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateSecretGeneratorRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSecretGeneratorRequest;
  static deserializeBinaryFromReader(message: UpdateSecretGeneratorRequest, reader: jspb.BinaryReader): UpdateSecretGeneratorRequest;
}

export namespace UpdateSecretGeneratorRequest {
  export type AsObject = {
    generatorType: zitadel_settings_pb.SecretGeneratorType,
    length: number,
    expiry?: google_protobuf_duration_pb.Duration.AsObject,
    includeLowerLetters: boolean,
    includeUpperLetters: boolean,
    includeDigits: boolean,
    includeSymbols: boolean,
  }
}

export class UpdateSecretGeneratorResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateSecretGeneratorResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateSecretGeneratorResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSecretGeneratorResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSecretGeneratorResponse): UpdateSecretGeneratorResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateSecretGeneratorResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSecretGeneratorResponse;
  static deserializeBinaryFromReader(message: UpdateSecretGeneratorResponse, reader: jspb.BinaryReader): UpdateSecretGeneratorResponse;
}

export namespace UpdateSecretGeneratorResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetSMTPConfigRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSMTPConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSMTPConfigRequest): GetSMTPConfigRequest.AsObject;
  static serializeBinaryToWriter(message: GetSMTPConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSMTPConfigRequest;
  static deserializeBinaryFromReader(message: GetSMTPConfigRequest, reader: jspb.BinaryReader): GetSMTPConfigRequest;
}

export namespace GetSMTPConfigRequest {
  export type AsObject = {
  }
}

export class GetSMTPConfigResponse extends jspb.Message {
  getSmtpConfig(): zitadel_settings_pb.SMTPConfig | undefined;
  setSmtpConfig(value?: zitadel_settings_pb.SMTPConfig): GetSMTPConfigResponse;
  hasSmtpConfig(): boolean;
  clearSmtpConfig(): GetSMTPConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSMTPConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSMTPConfigResponse): GetSMTPConfigResponse.AsObject;
  static serializeBinaryToWriter(message: GetSMTPConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSMTPConfigResponse;
  static deserializeBinaryFromReader(message: GetSMTPConfigResponse, reader: jspb.BinaryReader): GetSMTPConfigResponse;
}

export namespace GetSMTPConfigResponse {
  export type AsObject = {
    smtpConfig?: zitadel_settings_pb.SMTPConfig.AsObject,
  }
}

export class GetSMTPConfigByIdRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetSMTPConfigByIdRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSMTPConfigByIdRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSMTPConfigByIdRequest): GetSMTPConfigByIdRequest.AsObject;
  static serializeBinaryToWriter(message: GetSMTPConfigByIdRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSMTPConfigByIdRequest;
  static deserializeBinaryFromReader(message: GetSMTPConfigByIdRequest, reader: jspb.BinaryReader): GetSMTPConfigByIdRequest;
}

export namespace GetSMTPConfigByIdRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetSMTPConfigByIdResponse extends jspb.Message {
  getSmtpConfig(): zitadel_settings_pb.SMTPConfig | undefined;
  setSmtpConfig(value?: zitadel_settings_pb.SMTPConfig): GetSMTPConfigByIdResponse;
  hasSmtpConfig(): boolean;
  clearSmtpConfig(): GetSMTPConfigByIdResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSMTPConfigByIdResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSMTPConfigByIdResponse): GetSMTPConfigByIdResponse.AsObject;
  static serializeBinaryToWriter(message: GetSMTPConfigByIdResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSMTPConfigByIdResponse;
  static deserializeBinaryFromReader(message: GetSMTPConfigByIdResponse, reader: jspb.BinaryReader): GetSMTPConfigByIdResponse;
}

export namespace GetSMTPConfigByIdResponse {
  export type AsObject = {
    smtpConfig?: zitadel_settings_pb.SMTPConfig.AsObject,
  }
}

export class ListSMTPConfigsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListSMTPConfigsRequest;
  hasQuery(): boolean;
  clearQuery(): ListSMTPConfigsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSMTPConfigsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSMTPConfigsRequest): ListSMTPConfigsRequest.AsObject;
  static serializeBinaryToWriter(message: ListSMTPConfigsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSMTPConfigsRequest;
  static deserializeBinaryFromReader(message: ListSMTPConfigsRequest, reader: jspb.BinaryReader): ListSMTPConfigsRequest;
}

export namespace ListSMTPConfigsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListSMTPConfigsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListSMTPConfigsResponse;
  hasDetails(): boolean;
  clearDetails(): ListSMTPConfigsResponse;

  getResultList(): Array<zitadel_settings_pb.SMTPConfig>;
  setResultList(value: Array<zitadel_settings_pb.SMTPConfig>): ListSMTPConfigsResponse;
  clearResultList(): ListSMTPConfigsResponse;
  addResult(value?: zitadel_settings_pb.SMTPConfig, index?: number): zitadel_settings_pb.SMTPConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSMTPConfigsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSMTPConfigsResponse): ListSMTPConfigsResponse.AsObject;
  static serializeBinaryToWriter(message: ListSMTPConfigsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSMTPConfigsResponse;
  static deserializeBinaryFromReader(message: ListSMTPConfigsResponse, reader: jspb.BinaryReader): ListSMTPConfigsResponse;
}

export namespace ListSMTPConfigsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_settings_pb.SMTPConfig.AsObject>,
  }
}

export class AddSMTPConfigRequest extends jspb.Message {
  getSenderAddress(): string;
  setSenderAddress(value: string): AddSMTPConfigRequest;

  getSenderName(): string;
  setSenderName(value: string): AddSMTPConfigRequest;

  getTls(): boolean;
  setTls(value: boolean): AddSMTPConfigRequest;

  getHost(): string;
  setHost(value: string): AddSMTPConfigRequest;

  getUser(): string;
  setUser(value: string): AddSMTPConfigRequest;

  getPassword(): string;
  setPassword(value: string): AddSMTPConfigRequest;

  getReplyToAddress(): string;
  setReplyToAddress(value: string): AddSMTPConfigRequest;

  getDescription(): string;
  setDescription(value: string): AddSMTPConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSMTPConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddSMTPConfigRequest): AddSMTPConfigRequest.AsObject;
  static serializeBinaryToWriter(message: AddSMTPConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSMTPConfigRequest;
  static deserializeBinaryFromReader(message: AddSMTPConfigRequest, reader: jspb.BinaryReader): AddSMTPConfigRequest;
}

export namespace AddSMTPConfigRequest {
  export type AsObject = {
    senderAddress: string,
    senderName: string,
    tls: boolean,
    host: string,
    user: string,
    password: string,
    replyToAddress: string,
    description: string,
  }
}

export class AddSMTPConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddSMTPConfigResponse;
  hasDetails(): boolean;
  clearDetails(): AddSMTPConfigResponse;

  getId(): string;
  setId(value: string): AddSMTPConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSMTPConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSMTPConfigResponse): AddSMTPConfigResponse.AsObject;
  static serializeBinaryToWriter(message: AddSMTPConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSMTPConfigResponse;
  static deserializeBinaryFromReader(message: AddSMTPConfigResponse, reader: jspb.BinaryReader): AddSMTPConfigResponse;
}

export namespace AddSMTPConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateSMTPConfigRequest extends jspb.Message {
  getSenderAddress(): string;
  setSenderAddress(value: string): UpdateSMTPConfigRequest;

  getSenderName(): string;
  setSenderName(value: string): UpdateSMTPConfigRequest;

  getTls(): boolean;
  setTls(value: boolean): UpdateSMTPConfigRequest;

  getHost(): string;
  setHost(value: string): UpdateSMTPConfigRequest;

  getUser(): string;
  setUser(value: string): UpdateSMTPConfigRequest;

  getReplyToAddress(): string;
  setReplyToAddress(value: string): UpdateSMTPConfigRequest;

  getPassword(): string;
  setPassword(value: string): UpdateSMTPConfigRequest;

  getDescription(): string;
  setDescription(value: string): UpdateSMTPConfigRequest;

  getId(): string;
  setId(value: string): UpdateSMTPConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSMTPConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSMTPConfigRequest): UpdateSMTPConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateSMTPConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSMTPConfigRequest;
  static deserializeBinaryFromReader(message: UpdateSMTPConfigRequest, reader: jspb.BinaryReader): UpdateSMTPConfigRequest;
}

export namespace UpdateSMTPConfigRequest {
  export type AsObject = {
    senderAddress: string,
    senderName: string,
    tls: boolean,
    host: string,
    user: string,
    replyToAddress: string,
    password: string,
    description: string,
    id: string,
  }
}

export class UpdateSMTPConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateSMTPConfigResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateSMTPConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSMTPConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSMTPConfigResponse): UpdateSMTPConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateSMTPConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSMTPConfigResponse;
  static deserializeBinaryFromReader(message: UpdateSMTPConfigResponse, reader: jspb.BinaryReader): UpdateSMTPConfigResponse;
}

export namespace UpdateSMTPConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateSMTPConfigPasswordRequest extends jspb.Message {
  getPassword(): string;
  setPassword(value: string): UpdateSMTPConfigPasswordRequest;

  getId(): string;
  setId(value: string): UpdateSMTPConfigPasswordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSMTPConfigPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSMTPConfigPasswordRequest): UpdateSMTPConfigPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateSMTPConfigPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSMTPConfigPasswordRequest;
  static deserializeBinaryFromReader(message: UpdateSMTPConfigPasswordRequest, reader: jspb.BinaryReader): UpdateSMTPConfigPasswordRequest;
}

export namespace UpdateSMTPConfigPasswordRequest {
  export type AsObject = {
    password: string,
    id: string,
  }
}

export class UpdateSMTPConfigPasswordResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateSMTPConfigPasswordResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateSMTPConfigPasswordResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSMTPConfigPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSMTPConfigPasswordResponse): UpdateSMTPConfigPasswordResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateSMTPConfigPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSMTPConfigPasswordResponse;
  static deserializeBinaryFromReader(message: UpdateSMTPConfigPasswordResponse, reader: jspb.BinaryReader): UpdateSMTPConfigPasswordResponse;
}

export namespace UpdateSMTPConfigPasswordResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ActivateSMTPConfigRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ActivateSMTPConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateSMTPConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateSMTPConfigRequest): ActivateSMTPConfigRequest.AsObject;
  static serializeBinaryToWriter(message: ActivateSMTPConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateSMTPConfigRequest;
  static deserializeBinaryFromReader(message: ActivateSMTPConfigRequest, reader: jspb.BinaryReader): ActivateSMTPConfigRequest;
}

export namespace ActivateSMTPConfigRequest {
  export type AsObject = {
    id: string,
  }
}

export class ActivateSMTPConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ActivateSMTPConfigResponse;
  hasDetails(): boolean;
  clearDetails(): ActivateSMTPConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateSMTPConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateSMTPConfigResponse): ActivateSMTPConfigResponse.AsObject;
  static serializeBinaryToWriter(message: ActivateSMTPConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateSMTPConfigResponse;
  static deserializeBinaryFromReader(message: ActivateSMTPConfigResponse, reader: jspb.BinaryReader): ActivateSMTPConfigResponse;
}

export namespace ActivateSMTPConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateSMTPConfigRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeactivateSMTPConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateSMTPConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateSMTPConfigRequest): DeactivateSMTPConfigRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateSMTPConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateSMTPConfigRequest;
  static deserializeBinaryFromReader(message: DeactivateSMTPConfigRequest, reader: jspb.BinaryReader): DeactivateSMTPConfigRequest;
}

export namespace DeactivateSMTPConfigRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeactivateSMTPConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateSMTPConfigResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateSMTPConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateSMTPConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateSMTPConfigResponse): DeactivateSMTPConfigResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateSMTPConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateSMTPConfigResponse;
  static deserializeBinaryFromReader(message: DeactivateSMTPConfigResponse, reader: jspb.BinaryReader): DeactivateSMTPConfigResponse;
}

export namespace DeactivateSMTPConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveSMTPConfigRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RemoveSMTPConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveSMTPConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveSMTPConfigRequest): RemoveSMTPConfigRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveSMTPConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveSMTPConfigRequest;
  static deserializeBinaryFromReader(message: RemoveSMTPConfigRequest, reader: jspb.BinaryReader): RemoveSMTPConfigRequest;
}

export namespace RemoveSMTPConfigRequest {
  export type AsObject = {
    id: string,
  }
}

export class RemoveSMTPConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveSMTPConfigResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveSMTPConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveSMTPConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveSMTPConfigResponse): RemoveSMTPConfigResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveSMTPConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveSMTPConfigResponse;
  static deserializeBinaryFromReader(message: RemoveSMTPConfigResponse, reader: jspb.BinaryReader): RemoveSMTPConfigResponse;
}

export namespace RemoveSMTPConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class TestSMTPConfigByIdRequest extends jspb.Message {
  getId(): string;
  setId(value: string): TestSMTPConfigByIdRequest;

  getReceiverAddress(): string;
  setReceiverAddress(value: string): TestSMTPConfigByIdRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TestSMTPConfigByIdRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TestSMTPConfigByIdRequest): TestSMTPConfigByIdRequest.AsObject;
  static serializeBinaryToWriter(message: TestSMTPConfigByIdRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TestSMTPConfigByIdRequest;
  static deserializeBinaryFromReader(message: TestSMTPConfigByIdRequest, reader: jspb.BinaryReader): TestSMTPConfigByIdRequest;
}

export namespace TestSMTPConfigByIdRequest {
  export type AsObject = {
    id: string,
    receiverAddress: string,
  }
}

export class TestSMTPConfigByIdResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TestSMTPConfigByIdResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TestSMTPConfigByIdResponse): TestSMTPConfigByIdResponse.AsObject;
  static serializeBinaryToWriter(message: TestSMTPConfigByIdResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TestSMTPConfigByIdResponse;
  static deserializeBinaryFromReader(message: TestSMTPConfigByIdResponse, reader: jspb.BinaryReader): TestSMTPConfigByIdResponse;
}

export namespace TestSMTPConfigByIdResponse {
  export type AsObject = {
  }
}

export class TestSMTPConfigRequest extends jspb.Message {
  getSenderAddress(): string;
  setSenderAddress(value: string): TestSMTPConfigRequest;

  getSenderName(): string;
  setSenderName(value: string): TestSMTPConfigRequest;

  getTls(): boolean;
  setTls(value: boolean): TestSMTPConfigRequest;

  getHost(): string;
  setHost(value: string): TestSMTPConfigRequest;

  getUser(): string;
  setUser(value: string): TestSMTPConfigRequest;

  getPassword(): string;
  setPassword(value: string): TestSMTPConfigRequest;

  getReceiverAddress(): string;
  setReceiverAddress(value: string): TestSMTPConfigRequest;

  getId(): string;
  setId(value: string): TestSMTPConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TestSMTPConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TestSMTPConfigRequest): TestSMTPConfigRequest.AsObject;
  static serializeBinaryToWriter(message: TestSMTPConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TestSMTPConfigRequest;
  static deserializeBinaryFromReader(message: TestSMTPConfigRequest, reader: jspb.BinaryReader): TestSMTPConfigRequest;
}

export namespace TestSMTPConfigRequest {
  export type AsObject = {
    senderAddress: string,
    senderName: string,
    tls: boolean,
    host: string,
    user: string,
    password: string,
    receiverAddress: string,
    id: string,
  }
}

export class TestSMTPConfigResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TestSMTPConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TestSMTPConfigResponse): TestSMTPConfigResponse.AsObject;
  static serializeBinaryToWriter(message: TestSMTPConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TestSMTPConfigResponse;
  static deserializeBinaryFromReader(message: TestSMTPConfigResponse, reader: jspb.BinaryReader): TestSMTPConfigResponse;
}

export namespace TestSMTPConfigResponse {
  export type AsObject = {
  }
}

export class ListSMSProvidersRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListSMSProvidersRequest;
  hasQuery(): boolean;
  clearQuery(): ListSMSProvidersRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSMSProvidersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListSMSProvidersRequest): ListSMSProvidersRequest.AsObject;
  static serializeBinaryToWriter(message: ListSMSProvidersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSMSProvidersRequest;
  static deserializeBinaryFromReader(message: ListSMSProvidersRequest, reader: jspb.BinaryReader): ListSMSProvidersRequest;
}

export namespace ListSMSProvidersRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListSMSProvidersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListSMSProvidersResponse;
  hasDetails(): boolean;
  clearDetails(): ListSMSProvidersResponse;

  getResultList(): Array<zitadel_settings_pb.SMSProvider>;
  setResultList(value: Array<zitadel_settings_pb.SMSProvider>): ListSMSProvidersResponse;
  clearResultList(): ListSMSProvidersResponse;
  addResult(value?: zitadel_settings_pb.SMSProvider, index?: number): zitadel_settings_pb.SMSProvider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListSMSProvidersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListSMSProvidersResponse): ListSMSProvidersResponse.AsObject;
  static serializeBinaryToWriter(message: ListSMSProvidersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListSMSProvidersResponse;
  static deserializeBinaryFromReader(message: ListSMSProvidersResponse, reader: jspb.BinaryReader): ListSMSProvidersResponse;
}

export namespace ListSMSProvidersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_settings_pb.SMSProvider.AsObject>,
  }
}

export class GetSMSProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetSMSProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSMSProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSMSProviderRequest): GetSMSProviderRequest.AsObject;
  static serializeBinaryToWriter(message: GetSMSProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSMSProviderRequest;
  static deserializeBinaryFromReader(message: GetSMSProviderRequest, reader: jspb.BinaryReader): GetSMSProviderRequest;
}

export namespace GetSMSProviderRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetSMSProviderResponse extends jspb.Message {
  getConfig(): zitadel_settings_pb.SMSProvider | undefined;
  setConfig(value?: zitadel_settings_pb.SMSProvider): GetSMSProviderResponse;
  hasConfig(): boolean;
  clearConfig(): GetSMSProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSMSProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSMSProviderResponse): GetSMSProviderResponse.AsObject;
  static serializeBinaryToWriter(message: GetSMSProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSMSProviderResponse;
  static deserializeBinaryFromReader(message: GetSMSProviderResponse, reader: jspb.BinaryReader): GetSMSProviderResponse;
}

export namespace GetSMSProviderResponse {
  export type AsObject = {
    config?: zitadel_settings_pb.SMSProvider.AsObject,
  }
}

export class AddSMSProviderTwilioRequest extends jspb.Message {
  getSid(): string;
  setSid(value: string): AddSMSProviderTwilioRequest;

  getToken(): string;
  setToken(value: string): AddSMSProviderTwilioRequest;

  getSenderNumber(): string;
  setSenderNumber(value: string): AddSMSProviderTwilioRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSMSProviderTwilioRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddSMSProviderTwilioRequest): AddSMSProviderTwilioRequest.AsObject;
  static serializeBinaryToWriter(message: AddSMSProviderTwilioRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSMSProviderTwilioRequest;
  static deserializeBinaryFromReader(message: AddSMSProviderTwilioRequest, reader: jspb.BinaryReader): AddSMSProviderTwilioRequest;
}

export namespace AddSMSProviderTwilioRequest {
  export type AsObject = {
    sid: string,
    token: string,
    senderNumber: string,
  }
}

export class AddSMSProviderTwilioResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddSMSProviderTwilioResponse;
  hasDetails(): boolean;
  clearDetails(): AddSMSProviderTwilioResponse;

  getId(): string;
  setId(value: string): AddSMSProviderTwilioResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSMSProviderTwilioResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSMSProviderTwilioResponse): AddSMSProviderTwilioResponse.AsObject;
  static serializeBinaryToWriter(message: AddSMSProviderTwilioResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSMSProviderTwilioResponse;
  static deserializeBinaryFromReader(message: AddSMSProviderTwilioResponse, reader: jspb.BinaryReader): AddSMSProviderTwilioResponse;
}

export namespace AddSMSProviderTwilioResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateSMSProviderTwilioRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateSMSProviderTwilioRequest;

  getSid(): string;
  setSid(value: string): UpdateSMSProviderTwilioRequest;

  getSenderNumber(): string;
  setSenderNumber(value: string): UpdateSMSProviderTwilioRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSMSProviderTwilioRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSMSProviderTwilioRequest): UpdateSMSProviderTwilioRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateSMSProviderTwilioRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSMSProviderTwilioRequest;
  static deserializeBinaryFromReader(message: UpdateSMSProviderTwilioRequest, reader: jspb.BinaryReader): UpdateSMSProviderTwilioRequest;
}

export namespace UpdateSMSProviderTwilioRequest {
  export type AsObject = {
    id: string,
    sid: string,
    senderNumber: string,
  }
}

export class UpdateSMSProviderTwilioResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateSMSProviderTwilioResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateSMSProviderTwilioResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSMSProviderTwilioResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSMSProviderTwilioResponse): UpdateSMSProviderTwilioResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateSMSProviderTwilioResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSMSProviderTwilioResponse;
  static deserializeBinaryFromReader(message: UpdateSMSProviderTwilioResponse, reader: jspb.BinaryReader): UpdateSMSProviderTwilioResponse;
}

export namespace UpdateSMSProviderTwilioResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateSMSProviderTwilioTokenRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateSMSProviderTwilioTokenRequest;

  getToken(): string;
  setToken(value: string): UpdateSMSProviderTwilioTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSMSProviderTwilioTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSMSProviderTwilioTokenRequest): UpdateSMSProviderTwilioTokenRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateSMSProviderTwilioTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSMSProviderTwilioTokenRequest;
  static deserializeBinaryFromReader(message: UpdateSMSProviderTwilioTokenRequest, reader: jspb.BinaryReader): UpdateSMSProviderTwilioTokenRequest;
}

export namespace UpdateSMSProviderTwilioTokenRequest {
  export type AsObject = {
    id: string,
    token: string,
  }
}

export class UpdateSMSProviderTwilioTokenResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateSMSProviderTwilioTokenResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateSMSProviderTwilioTokenResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSMSProviderTwilioTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSMSProviderTwilioTokenResponse): UpdateSMSProviderTwilioTokenResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateSMSProviderTwilioTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSMSProviderTwilioTokenResponse;
  static deserializeBinaryFromReader(message: UpdateSMSProviderTwilioTokenResponse, reader: jspb.BinaryReader): UpdateSMSProviderTwilioTokenResponse;
}

export namespace UpdateSMSProviderTwilioTokenResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ActivateSMSProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ActivateSMSProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateSMSProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateSMSProviderRequest): ActivateSMSProviderRequest.AsObject;
  static serializeBinaryToWriter(message: ActivateSMSProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateSMSProviderRequest;
  static deserializeBinaryFromReader(message: ActivateSMSProviderRequest, reader: jspb.BinaryReader): ActivateSMSProviderRequest;
}

export namespace ActivateSMSProviderRequest {
  export type AsObject = {
    id: string,
  }
}

export class ActivateSMSProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ActivateSMSProviderResponse;
  hasDetails(): boolean;
  clearDetails(): ActivateSMSProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateSMSProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateSMSProviderResponse): ActivateSMSProviderResponse.AsObject;
  static serializeBinaryToWriter(message: ActivateSMSProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateSMSProviderResponse;
  static deserializeBinaryFromReader(message: ActivateSMSProviderResponse, reader: jspb.BinaryReader): ActivateSMSProviderResponse;
}

export namespace ActivateSMSProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateSMSProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeactivateSMSProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateSMSProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateSMSProviderRequest): DeactivateSMSProviderRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateSMSProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateSMSProviderRequest;
  static deserializeBinaryFromReader(message: DeactivateSMSProviderRequest, reader: jspb.BinaryReader): DeactivateSMSProviderRequest;
}

export namespace DeactivateSMSProviderRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeactivateSMSProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateSMSProviderResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateSMSProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateSMSProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateSMSProviderResponse): DeactivateSMSProviderResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateSMSProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateSMSProviderResponse;
  static deserializeBinaryFromReader(message: DeactivateSMSProviderResponse, reader: jspb.BinaryReader): DeactivateSMSProviderResponse;
}

export namespace DeactivateSMSProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveSMSProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RemoveSMSProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveSMSProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveSMSProviderRequest): RemoveSMSProviderRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveSMSProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveSMSProviderRequest;
  static deserializeBinaryFromReader(message: RemoveSMSProviderRequest, reader: jspb.BinaryReader): RemoveSMSProviderRequest;
}

export namespace RemoveSMSProviderRequest {
  export type AsObject = {
    id: string,
  }
}

export class RemoveSMSProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveSMSProviderResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveSMSProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveSMSProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveSMSProviderResponse): RemoveSMSProviderResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveSMSProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveSMSProviderResponse;
  static deserializeBinaryFromReader(message: RemoveSMSProviderResponse, reader: jspb.BinaryReader): RemoveSMSProviderResponse;
}

export namespace RemoveSMSProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetFileSystemNotificationProviderRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFileSystemNotificationProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFileSystemNotificationProviderRequest): GetFileSystemNotificationProviderRequest.AsObject;
  static serializeBinaryToWriter(message: GetFileSystemNotificationProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFileSystemNotificationProviderRequest;
  static deserializeBinaryFromReader(message: GetFileSystemNotificationProviderRequest, reader: jspb.BinaryReader): GetFileSystemNotificationProviderRequest;
}

export namespace GetFileSystemNotificationProviderRequest {
  export type AsObject = {
  }
}

export class GetFileSystemNotificationProviderResponse extends jspb.Message {
  getProvider(): zitadel_settings_pb.DebugNotificationProvider | undefined;
  setProvider(value?: zitadel_settings_pb.DebugNotificationProvider): GetFileSystemNotificationProviderResponse;
  hasProvider(): boolean;
  clearProvider(): GetFileSystemNotificationProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFileSystemNotificationProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFileSystemNotificationProviderResponse): GetFileSystemNotificationProviderResponse.AsObject;
  static serializeBinaryToWriter(message: GetFileSystemNotificationProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFileSystemNotificationProviderResponse;
  static deserializeBinaryFromReader(message: GetFileSystemNotificationProviderResponse, reader: jspb.BinaryReader): GetFileSystemNotificationProviderResponse;
}

export namespace GetFileSystemNotificationProviderResponse {
  export type AsObject = {
    provider?: zitadel_settings_pb.DebugNotificationProvider.AsObject,
  }
}

export class GetLogNotificationProviderRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLogNotificationProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLogNotificationProviderRequest): GetLogNotificationProviderRequest.AsObject;
  static serializeBinaryToWriter(message: GetLogNotificationProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLogNotificationProviderRequest;
  static deserializeBinaryFromReader(message: GetLogNotificationProviderRequest, reader: jspb.BinaryReader): GetLogNotificationProviderRequest;
}

export namespace GetLogNotificationProviderRequest {
  export type AsObject = {
  }
}

export class GetLogNotificationProviderResponse extends jspb.Message {
  getProvider(): zitadel_settings_pb.DebugNotificationProvider | undefined;
  setProvider(value?: zitadel_settings_pb.DebugNotificationProvider): GetLogNotificationProviderResponse;
  hasProvider(): boolean;
  clearProvider(): GetLogNotificationProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLogNotificationProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLogNotificationProviderResponse): GetLogNotificationProviderResponse.AsObject;
  static serializeBinaryToWriter(message: GetLogNotificationProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLogNotificationProviderResponse;
  static deserializeBinaryFromReader(message: GetLogNotificationProviderResponse, reader: jspb.BinaryReader): GetLogNotificationProviderResponse;
}

export namespace GetLogNotificationProviderResponse {
  export type AsObject = {
    provider?: zitadel_settings_pb.DebugNotificationProvider.AsObject,
  }
}

export class GetOIDCSettingsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOIDCSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOIDCSettingsRequest): GetOIDCSettingsRequest.AsObject;
  static serializeBinaryToWriter(message: GetOIDCSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOIDCSettingsRequest;
  static deserializeBinaryFromReader(message: GetOIDCSettingsRequest, reader: jspb.BinaryReader): GetOIDCSettingsRequest;
}

export namespace GetOIDCSettingsRequest {
  export type AsObject = {
  }
}

export class GetOIDCSettingsResponse extends jspb.Message {
  getSettings(): zitadel_settings_pb.OIDCSettings | undefined;
  setSettings(value?: zitadel_settings_pb.OIDCSettings): GetOIDCSettingsResponse;
  hasSettings(): boolean;
  clearSettings(): GetOIDCSettingsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOIDCSettingsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOIDCSettingsResponse): GetOIDCSettingsResponse.AsObject;
  static serializeBinaryToWriter(message: GetOIDCSettingsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOIDCSettingsResponse;
  static deserializeBinaryFromReader(message: GetOIDCSettingsResponse, reader: jspb.BinaryReader): GetOIDCSettingsResponse;
}

export namespace GetOIDCSettingsResponse {
  export type AsObject = {
    settings?: zitadel_settings_pb.OIDCSettings.AsObject,
  }
}

export class AddOIDCSettingsRequest extends jspb.Message {
  getAccessTokenLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setAccessTokenLifetime(value?: google_protobuf_duration_pb.Duration): AddOIDCSettingsRequest;
  hasAccessTokenLifetime(): boolean;
  clearAccessTokenLifetime(): AddOIDCSettingsRequest;

  getIdTokenLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setIdTokenLifetime(value?: google_protobuf_duration_pb.Duration): AddOIDCSettingsRequest;
  hasIdTokenLifetime(): boolean;
  clearIdTokenLifetime(): AddOIDCSettingsRequest;

  getRefreshTokenIdleExpiration(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshTokenIdleExpiration(value?: google_protobuf_duration_pb.Duration): AddOIDCSettingsRequest;
  hasRefreshTokenIdleExpiration(): boolean;
  clearRefreshTokenIdleExpiration(): AddOIDCSettingsRequest;

  getRefreshTokenExpiration(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshTokenExpiration(value?: google_protobuf_duration_pb.Duration): AddOIDCSettingsRequest;
  hasRefreshTokenExpiration(): boolean;
  clearRefreshTokenExpiration(): AddOIDCSettingsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOIDCSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOIDCSettingsRequest): AddOIDCSettingsRequest.AsObject;
  static serializeBinaryToWriter(message: AddOIDCSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOIDCSettingsRequest;
  static deserializeBinaryFromReader(message: AddOIDCSettingsRequest, reader: jspb.BinaryReader): AddOIDCSettingsRequest;
}

export namespace AddOIDCSettingsRequest {
  export type AsObject = {
    accessTokenLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    idTokenLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    refreshTokenIdleExpiration?: google_protobuf_duration_pb.Duration.AsObject,
    refreshTokenExpiration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class AddOIDCSettingsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddOIDCSettingsResponse;
  hasDetails(): boolean;
  clearDetails(): AddOIDCSettingsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOIDCSettingsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOIDCSettingsResponse): AddOIDCSettingsResponse.AsObject;
  static serializeBinaryToWriter(message: AddOIDCSettingsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOIDCSettingsResponse;
  static deserializeBinaryFromReader(message: AddOIDCSettingsResponse, reader: jspb.BinaryReader): AddOIDCSettingsResponse;
}

export namespace AddOIDCSettingsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateOIDCSettingsRequest extends jspb.Message {
  getAccessTokenLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setAccessTokenLifetime(value?: google_protobuf_duration_pb.Duration): UpdateOIDCSettingsRequest;
  hasAccessTokenLifetime(): boolean;
  clearAccessTokenLifetime(): UpdateOIDCSettingsRequest;

  getIdTokenLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setIdTokenLifetime(value?: google_protobuf_duration_pb.Duration): UpdateOIDCSettingsRequest;
  hasIdTokenLifetime(): boolean;
  clearIdTokenLifetime(): UpdateOIDCSettingsRequest;

  getRefreshTokenIdleExpiration(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshTokenIdleExpiration(value?: google_protobuf_duration_pb.Duration): UpdateOIDCSettingsRequest;
  hasRefreshTokenIdleExpiration(): boolean;
  clearRefreshTokenIdleExpiration(): UpdateOIDCSettingsRequest;

  getRefreshTokenExpiration(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshTokenExpiration(value?: google_protobuf_duration_pb.Duration): UpdateOIDCSettingsRequest;
  hasRefreshTokenExpiration(): boolean;
  clearRefreshTokenExpiration(): UpdateOIDCSettingsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOIDCSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOIDCSettingsRequest): UpdateOIDCSettingsRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateOIDCSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOIDCSettingsRequest;
  static deserializeBinaryFromReader(message: UpdateOIDCSettingsRequest, reader: jspb.BinaryReader): UpdateOIDCSettingsRequest;
}

export namespace UpdateOIDCSettingsRequest {
  export type AsObject = {
    accessTokenLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    idTokenLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    refreshTokenIdleExpiration?: google_protobuf_duration_pb.Duration.AsObject,
    refreshTokenExpiration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class UpdateOIDCSettingsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateOIDCSettingsResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateOIDCSettingsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOIDCSettingsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOIDCSettingsResponse): UpdateOIDCSettingsResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateOIDCSettingsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOIDCSettingsResponse;
  static deserializeBinaryFromReader(message: UpdateOIDCSettingsResponse, reader: jspb.BinaryReader): UpdateOIDCSettingsResponse;
}

export namespace UpdateOIDCSettingsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetSecurityPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSecurityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSecurityPolicyRequest): GetSecurityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetSecurityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSecurityPolicyRequest;
  static deserializeBinaryFromReader(message: GetSecurityPolicyRequest, reader: jspb.BinaryReader): GetSecurityPolicyRequest;
}

export namespace GetSecurityPolicyRequest {
  export type AsObject = {
  }
}

export class GetSecurityPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_settings_pb.SecurityPolicy | undefined;
  setPolicy(value?: zitadel_settings_pb.SecurityPolicy): GetSecurityPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetSecurityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSecurityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSecurityPolicyResponse): GetSecurityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetSecurityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSecurityPolicyResponse;
  static deserializeBinaryFromReader(message: GetSecurityPolicyResponse, reader: jspb.BinaryReader): GetSecurityPolicyResponse;
}

export namespace GetSecurityPolicyResponse {
  export type AsObject = {
    policy?: zitadel_settings_pb.SecurityPolicy.AsObject,
  }
}

export class SetSecurityPolicyRequest extends jspb.Message {
  getEnableIframeEmbedding(): boolean;
  setEnableIframeEmbedding(value: boolean): SetSecurityPolicyRequest;

  getAllowedOriginsList(): Array<string>;
  setAllowedOriginsList(value: Array<string>): SetSecurityPolicyRequest;
  clearAllowedOriginsList(): SetSecurityPolicyRequest;
  addAllowedOrigins(value: string, index?: number): SetSecurityPolicyRequest;

  getEnableImpersonation(): boolean;
  setEnableImpersonation(value: boolean): SetSecurityPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetSecurityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetSecurityPolicyRequest): SetSecurityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: SetSecurityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetSecurityPolicyRequest;
  static deserializeBinaryFromReader(message: SetSecurityPolicyRequest, reader: jspb.BinaryReader): SetSecurityPolicyRequest;
}

export namespace SetSecurityPolicyRequest {
  export type AsObject = {
    enableIframeEmbedding: boolean,
    allowedOriginsList: Array<string>,
    enableImpersonation: boolean,
  }
}

export class SetSecurityPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetSecurityPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): SetSecurityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetSecurityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetSecurityPolicyResponse): SetSecurityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: SetSecurityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetSecurityPolicyResponse;
  static deserializeBinaryFromReader(message: SetSecurityPolicyResponse, reader: jspb.BinaryReader): SetSecurityPolicyResponse;
}

export namespace SetSecurityPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class IsOrgUniqueRequest extends jspb.Message {
  getName(): string;
  setName(value: string): IsOrgUniqueRequest;

  getDomain(): string;
  setDomain(value: string): IsOrgUniqueRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsOrgUniqueRequest.AsObject;
  static toObject(includeInstance: boolean, msg: IsOrgUniqueRequest): IsOrgUniqueRequest.AsObject;
  static serializeBinaryToWriter(message: IsOrgUniqueRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsOrgUniqueRequest;
  static deserializeBinaryFromReader(message: IsOrgUniqueRequest, reader: jspb.BinaryReader): IsOrgUniqueRequest;
}

export namespace IsOrgUniqueRequest {
  export type AsObject = {
    name: string,
    domain: string,
  }
}

export class IsOrgUniqueResponse extends jspb.Message {
  getIsUnique(): boolean;
  setIsUnique(value: boolean): IsOrgUniqueResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsOrgUniqueResponse.AsObject;
  static toObject(includeInstance: boolean, msg: IsOrgUniqueResponse): IsOrgUniqueResponse.AsObject;
  static serializeBinaryToWriter(message: IsOrgUniqueResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsOrgUniqueResponse;
  static deserializeBinaryFromReader(message: IsOrgUniqueResponse, reader: jspb.BinaryReader): IsOrgUniqueResponse;
}

export namespace IsOrgUniqueResponse {
  export type AsObject = {
    isUnique: boolean,
  }
}

export class GetOrgByIDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetOrgByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgByIDRequest): GetOrgByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrgByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgByIDRequest;
  static deserializeBinaryFromReader(message: GetOrgByIDRequest, reader: jspb.BinaryReader): GetOrgByIDRequest;
}

export namespace GetOrgByIDRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetOrgByIDResponse extends jspb.Message {
  getOrg(): zitadel_org_pb.Org | undefined;
  setOrg(value?: zitadel_org_pb.Org): GetOrgByIDResponse;
  hasOrg(): boolean;
  clearOrg(): GetOrgByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgByIDResponse): GetOrgByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrgByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgByIDResponse;
  static deserializeBinaryFromReader(message: GetOrgByIDResponse, reader: jspb.BinaryReader): GetOrgByIDResponse;
}

export namespace GetOrgByIDResponse {
  export type AsObject = {
    org?: zitadel_org_pb.Org.AsObject,
  }
}

export class ListOrgsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListOrgsRequest;
  hasQuery(): boolean;
  clearQuery(): ListOrgsRequest;

  getSortingColumn(): zitadel_org_pb.OrgFieldName;
  setSortingColumn(value: zitadel_org_pb.OrgFieldName): ListOrgsRequest;

  getQueriesList(): Array<zitadel_org_pb.OrgQuery>;
  setQueriesList(value: Array<zitadel_org_pb.OrgQuery>): ListOrgsRequest;
  clearQueriesList(): ListOrgsRequest;
  addQueries(value?: zitadel_org_pb.OrgQuery, index?: number): zitadel_org_pb.OrgQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgsRequest): ListOrgsRequest.AsObject;
  static serializeBinaryToWriter(message: ListOrgsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgsRequest;
  static deserializeBinaryFromReader(message: ListOrgsRequest, reader: jspb.BinaryReader): ListOrgsRequest;
}

export namespace ListOrgsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_org_pb.OrgFieldName,
    queriesList: Array<zitadel_org_pb.OrgQuery.AsObject>,
  }
}

export class ListOrgsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListOrgsResponse;
  hasDetails(): boolean;
  clearDetails(): ListOrgsResponse;

  getSortingColumn(): zitadel_org_pb.OrgFieldName;
  setSortingColumn(value: zitadel_org_pb.OrgFieldName): ListOrgsResponse;

  getResultList(): Array<zitadel_org_pb.Org>;
  setResultList(value: Array<zitadel_org_pb.Org>): ListOrgsResponse;
  clearResultList(): ListOrgsResponse;
  addResult(value?: zitadel_org_pb.Org, index?: number): zitadel_org_pb.Org;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgsResponse): ListOrgsResponse.AsObject;
  static serializeBinaryToWriter(message: ListOrgsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgsResponse;
  static deserializeBinaryFromReader(message: ListOrgsResponse, reader: jspb.BinaryReader): ListOrgsResponse;
}

export namespace ListOrgsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_org_pb.OrgFieldName,
    resultList: Array<zitadel_org_pb.Org.AsObject>,
  }
}

export class SetUpOrgRequest extends jspb.Message {
  getOrg(): SetUpOrgRequest.Org | undefined;
  setOrg(value?: SetUpOrgRequest.Org): SetUpOrgRequest;
  hasOrg(): boolean;
  clearOrg(): SetUpOrgRequest;

  getHuman(): SetUpOrgRequest.Human | undefined;
  setHuman(value?: SetUpOrgRequest.Human): SetUpOrgRequest;
  hasHuman(): boolean;
  clearHuman(): SetUpOrgRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): SetUpOrgRequest;
  clearRolesList(): SetUpOrgRequest;
  addRoles(value: string, index?: number): SetUpOrgRequest;

  getUserCase(): SetUpOrgRequest.UserCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUpOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetUpOrgRequest): SetUpOrgRequest.AsObject;
  static serializeBinaryToWriter(message: SetUpOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUpOrgRequest;
  static deserializeBinaryFromReader(message: SetUpOrgRequest, reader: jspb.BinaryReader): SetUpOrgRequest;
}

export namespace SetUpOrgRequest {
  export type AsObject = {
    org?: SetUpOrgRequest.Org.AsObject,
    human?: SetUpOrgRequest.Human.AsObject,
    rolesList: Array<string>,
  }

  export class Org extends jspb.Message {
    getName(): string;
    setName(value: string): Org;

    getDomain(): string;
    setDomain(value: string): Org;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Org.AsObject;
    static toObject(includeInstance: boolean, msg: Org): Org.AsObject;
    static serializeBinaryToWriter(message: Org, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Org;
    static deserializeBinaryFromReader(message: Org, reader: jspb.BinaryReader): Org;
  }

  export namespace Org {
    export type AsObject = {
      name: string,
      domain: string,
    }
  }


  export class Human extends jspb.Message {
    getUserName(): string;
    setUserName(value: string): Human;

    getProfile(): SetUpOrgRequest.Human.Profile | undefined;
    setProfile(value?: SetUpOrgRequest.Human.Profile): Human;
    hasProfile(): boolean;
    clearProfile(): Human;

    getEmail(): SetUpOrgRequest.Human.Email | undefined;
    setEmail(value?: SetUpOrgRequest.Human.Email): Human;
    hasEmail(): boolean;
    clearEmail(): Human;

    getPhone(): SetUpOrgRequest.Human.Phone | undefined;
    setPhone(value?: SetUpOrgRequest.Human.Phone): Human;
    hasPhone(): boolean;
    clearPhone(): Human;

    getPassword(): string;
    setPassword(value: string): Human;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Human.AsObject;
    static toObject(includeInstance: boolean, msg: Human): Human.AsObject;
    static serializeBinaryToWriter(message: Human, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Human;
    static deserializeBinaryFromReader(message: Human, reader: jspb.BinaryReader): Human;
  }

  export namespace Human {
    export type AsObject = {
      userName: string,
      profile?: SetUpOrgRequest.Human.Profile.AsObject,
      email?: SetUpOrgRequest.Human.Email.AsObject,
      phone?: SetUpOrgRequest.Human.Phone.AsObject,
      password: string,
    }

    export class Profile extends jspb.Message {
      getFirstName(): string;
      setFirstName(value: string): Profile;

      getLastName(): string;
      setLastName(value: string): Profile;

      getNickName(): string;
      setNickName(value: string): Profile;

      getDisplayName(): string;
      setDisplayName(value: string): Profile;

      getPreferredLanguage(): string;
      setPreferredLanguage(value: string): Profile;

      getGender(): zitadel_user_pb.Gender;
      setGender(value: zitadel_user_pb.Gender): Profile;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Profile.AsObject;
      static toObject(includeInstance: boolean, msg: Profile): Profile.AsObject;
      static serializeBinaryToWriter(message: Profile, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Profile;
      static deserializeBinaryFromReader(message: Profile, reader: jspb.BinaryReader): Profile;
    }

    export namespace Profile {
      export type AsObject = {
        firstName: string,
        lastName: string,
        nickName: string,
        displayName: string,
        preferredLanguage: string,
        gender: zitadel_user_pb.Gender,
      }
    }


    export class Email extends jspb.Message {
      getEmail(): string;
      setEmail(value: string): Email;

      getIsEmailVerified(): boolean;
      setIsEmailVerified(value: boolean): Email;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Email.AsObject;
      static toObject(includeInstance: boolean, msg: Email): Email.AsObject;
      static serializeBinaryToWriter(message: Email, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Email;
      static deserializeBinaryFromReader(message: Email, reader: jspb.BinaryReader): Email;
    }

    export namespace Email {
      export type AsObject = {
        email: string,
        isEmailVerified: boolean,
      }
    }


    export class Phone extends jspb.Message {
      getPhone(): string;
      setPhone(value: string): Phone;

      getIsPhoneVerified(): boolean;
      setIsPhoneVerified(value: boolean): Phone;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Phone.AsObject;
      static toObject(includeInstance: boolean, msg: Phone): Phone.AsObject;
      static serializeBinaryToWriter(message: Phone, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Phone;
      static deserializeBinaryFromReader(message: Phone, reader: jspb.BinaryReader): Phone;
    }

    export namespace Phone {
      export type AsObject = {
        phone: string,
        isPhoneVerified: boolean,
      }
    }

  }


  export enum UserCase { 
    USER_NOT_SET = 0,
    HUMAN = 2,
  }
}

export class SetUpOrgResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetUpOrgResponse;
  hasDetails(): boolean;
  clearDetails(): SetUpOrgResponse;

  getOrgId(): string;
  setOrgId(value: string): SetUpOrgResponse;

  getUserId(): string;
  setUserId(value: string): SetUpOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUpOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetUpOrgResponse): SetUpOrgResponse.AsObject;
  static serializeBinaryToWriter(message: SetUpOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUpOrgResponse;
  static deserializeBinaryFromReader(message: SetUpOrgResponse, reader: jspb.BinaryReader): SetUpOrgResponse;
}

export namespace SetUpOrgResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    orgId: string,
    userId: string,
  }
}

export class RemoveOrgRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): RemoveOrgRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgRequest): RemoveOrgRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgRequest;
  static deserializeBinaryFromReader(message: RemoveOrgRequest, reader: jspb.BinaryReader): RemoveOrgRequest;
}

export namespace RemoveOrgRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class RemoveOrgResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveOrgResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgResponse): RemoveOrgResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgResponse;
  static deserializeBinaryFromReader(message: RemoveOrgResponse, reader: jspb.BinaryReader): RemoveOrgResponse;
}

export namespace RemoveOrgResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetIDPByIDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetIDPByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetIDPByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetIDPByIDRequest): GetIDPByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetIDPByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetIDPByIDRequest;
  static deserializeBinaryFromReader(message: GetIDPByIDRequest, reader: jspb.BinaryReader): GetIDPByIDRequest;
}

export namespace GetIDPByIDRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetIDPByIDResponse extends jspb.Message {
  getIdp(): zitadel_idp_pb.IDP | undefined;
  setIdp(value?: zitadel_idp_pb.IDP): GetIDPByIDResponse;
  hasIdp(): boolean;
  clearIdp(): GetIDPByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetIDPByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetIDPByIDResponse): GetIDPByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetIDPByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetIDPByIDResponse;
  static deserializeBinaryFromReader(message: GetIDPByIDResponse, reader: jspb.BinaryReader): GetIDPByIDResponse;
}

export namespace GetIDPByIDResponse {
  export type AsObject = {
    idp?: zitadel_idp_pb.IDP.AsObject,
  }
}

export class ListIDPsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListIDPsRequest;
  hasQuery(): boolean;
  clearQuery(): ListIDPsRequest;

  getSortingColumn(): zitadel_idp_pb.IDPFieldName;
  setSortingColumn(value: zitadel_idp_pb.IDPFieldName): ListIDPsRequest;

  getQueriesList(): Array<IDPQuery>;
  setQueriesList(value: Array<IDPQuery>): ListIDPsRequest;
  clearQueriesList(): ListIDPsRequest;
  addQueries(value?: IDPQuery, index?: number): IDPQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListIDPsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListIDPsRequest): ListIDPsRequest.AsObject;
  static serializeBinaryToWriter(message: ListIDPsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListIDPsRequest;
  static deserializeBinaryFromReader(message: ListIDPsRequest, reader: jspb.BinaryReader): ListIDPsRequest;
}

export namespace ListIDPsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_idp_pb.IDPFieldName,
    queriesList: Array<IDPQuery.AsObject>,
  }
}

export class IDPQuery extends jspb.Message {
  getIdpIdQuery(): zitadel_idp_pb.IDPIDQuery | undefined;
  setIdpIdQuery(value?: zitadel_idp_pb.IDPIDQuery): IDPQuery;
  hasIdpIdQuery(): boolean;
  clearIdpIdQuery(): IDPQuery;

  getIdpNameQuery(): zitadel_idp_pb.IDPNameQuery | undefined;
  setIdpNameQuery(value?: zitadel_idp_pb.IDPNameQuery): IDPQuery;
  hasIdpNameQuery(): boolean;
  clearIdpNameQuery(): IDPQuery;

  getQueryCase(): IDPQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IDPQuery): IDPQuery.AsObject;
  static serializeBinaryToWriter(message: IDPQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPQuery;
  static deserializeBinaryFromReader(message: IDPQuery, reader: jspb.BinaryReader): IDPQuery;
}

export namespace IDPQuery {
  export type AsObject = {
    idpIdQuery?: zitadel_idp_pb.IDPIDQuery.AsObject,
    idpNameQuery?: zitadel_idp_pb.IDPNameQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    IDP_ID_QUERY = 1,
    IDP_NAME_QUERY = 2,
  }
}

export class ListIDPsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListIDPsResponse;
  hasDetails(): boolean;
  clearDetails(): ListIDPsResponse;

  getSortingColumn(): zitadel_idp_pb.IDPFieldName;
  setSortingColumn(value: zitadel_idp_pb.IDPFieldName): ListIDPsResponse;

  getResultList(): Array<zitadel_idp_pb.IDP>;
  setResultList(value: Array<zitadel_idp_pb.IDP>): ListIDPsResponse;
  clearResultList(): ListIDPsResponse;
  addResult(value?: zitadel_idp_pb.IDP, index?: number): zitadel_idp_pb.IDP;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListIDPsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListIDPsResponse): ListIDPsResponse.AsObject;
  static serializeBinaryToWriter(message: ListIDPsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListIDPsResponse;
  static deserializeBinaryFromReader(message: ListIDPsResponse, reader: jspb.BinaryReader): ListIDPsResponse;
}

export namespace ListIDPsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_idp_pb.IDPFieldName,
    resultList: Array<zitadel_idp_pb.IDP.AsObject>,
  }
}

export class AddOIDCIDPRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddOIDCIDPRequest;

  getStylingType(): zitadel_idp_pb.IDPStylingType;
  setStylingType(value: zitadel_idp_pb.IDPStylingType): AddOIDCIDPRequest;

  getClientId(): string;
  setClientId(value: string): AddOIDCIDPRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddOIDCIDPRequest;

  getIssuer(): string;
  setIssuer(value: string): AddOIDCIDPRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddOIDCIDPRequest;
  clearScopesList(): AddOIDCIDPRequest;
  addScopes(value: string, index?: number): AddOIDCIDPRequest;

  getDisplayNameMapping(): zitadel_idp_pb.OIDCMappingField;
  setDisplayNameMapping(value: zitadel_idp_pb.OIDCMappingField): AddOIDCIDPRequest;

  getUsernameMapping(): zitadel_idp_pb.OIDCMappingField;
  setUsernameMapping(value: zitadel_idp_pb.OIDCMappingField): AddOIDCIDPRequest;

  getAutoRegister(): boolean;
  setAutoRegister(value: boolean): AddOIDCIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOIDCIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOIDCIDPRequest): AddOIDCIDPRequest.AsObject;
  static serializeBinaryToWriter(message: AddOIDCIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOIDCIDPRequest;
  static deserializeBinaryFromReader(message: AddOIDCIDPRequest, reader: jspb.BinaryReader): AddOIDCIDPRequest;
}

export namespace AddOIDCIDPRequest {
  export type AsObject = {
    name: string,
    stylingType: zitadel_idp_pb.IDPStylingType,
    clientId: string,
    clientSecret: string,
    issuer: string,
    scopesList: Array<string>,
    displayNameMapping: zitadel_idp_pb.OIDCMappingField,
    usernameMapping: zitadel_idp_pb.OIDCMappingField,
    autoRegister: boolean,
  }
}

export class AddOIDCIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddOIDCIDPResponse;
  hasDetails(): boolean;
  clearDetails(): AddOIDCIDPResponse;

  getIdpId(): string;
  setIdpId(value: string): AddOIDCIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOIDCIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOIDCIDPResponse): AddOIDCIDPResponse.AsObject;
  static serializeBinaryToWriter(message: AddOIDCIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOIDCIDPResponse;
  static deserializeBinaryFromReader(message: AddOIDCIDPResponse, reader: jspb.BinaryReader): AddOIDCIDPResponse;
}

export namespace AddOIDCIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    idpId: string,
  }
}

export class AddJWTIDPRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddJWTIDPRequest;

  getStylingType(): zitadel_idp_pb.IDPStylingType;
  setStylingType(value: zitadel_idp_pb.IDPStylingType): AddJWTIDPRequest;

  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): AddJWTIDPRequest;

  getIssuer(): string;
  setIssuer(value: string): AddJWTIDPRequest;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): AddJWTIDPRequest;

  getHeaderName(): string;
  setHeaderName(value: string): AddJWTIDPRequest;

  getAutoRegister(): boolean;
  setAutoRegister(value: boolean): AddJWTIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddJWTIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddJWTIDPRequest): AddJWTIDPRequest.AsObject;
  static serializeBinaryToWriter(message: AddJWTIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddJWTIDPRequest;
  static deserializeBinaryFromReader(message: AddJWTIDPRequest, reader: jspb.BinaryReader): AddJWTIDPRequest;
}

export namespace AddJWTIDPRequest {
  export type AsObject = {
    name: string,
    stylingType: zitadel_idp_pb.IDPStylingType,
    jwtEndpoint: string,
    issuer: string,
    keysEndpoint: string,
    headerName: string,
    autoRegister: boolean,
  }
}

export class AddJWTIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddJWTIDPResponse;
  hasDetails(): boolean;
  clearDetails(): AddJWTIDPResponse;

  getIdpId(): string;
  setIdpId(value: string): AddJWTIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddJWTIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddJWTIDPResponse): AddJWTIDPResponse.AsObject;
  static serializeBinaryToWriter(message: AddJWTIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddJWTIDPResponse;
  static deserializeBinaryFromReader(message: AddJWTIDPResponse, reader: jspb.BinaryReader): AddJWTIDPResponse;
}

export namespace AddJWTIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    idpId: string,
  }
}

export class UpdateIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): UpdateIDPRequest;

  getName(): string;
  setName(value: string): UpdateIDPRequest;

  getStylingType(): zitadel_idp_pb.IDPStylingType;
  setStylingType(value: zitadel_idp_pb.IDPStylingType): UpdateIDPRequest;

  getAutoRegister(): boolean;
  setAutoRegister(value: boolean): UpdateIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateIDPRequest): UpdateIDPRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateIDPRequest;
  static deserializeBinaryFromReader(message: UpdateIDPRequest, reader: jspb.BinaryReader): UpdateIDPRequest;
}

export namespace UpdateIDPRequest {
  export type AsObject = {
    idpId: string,
    name: string,
    stylingType: zitadel_idp_pb.IDPStylingType,
    autoRegister: boolean,
  }
}

export class UpdateIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateIDPResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateIDPResponse): UpdateIDPResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateIDPResponse;
  static deserializeBinaryFromReader(message: UpdateIDPResponse, reader: jspb.BinaryReader): UpdateIDPResponse;
}

export namespace UpdateIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): DeactivateIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateIDPRequest): DeactivateIDPRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateIDPRequest;
  static deserializeBinaryFromReader(message: DeactivateIDPRequest, reader: jspb.BinaryReader): DeactivateIDPRequest;
}

export namespace DeactivateIDPRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class DeactivateIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateIDPResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateIDPResponse): DeactivateIDPResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateIDPResponse;
  static deserializeBinaryFromReader(message: DeactivateIDPResponse, reader: jspb.BinaryReader): DeactivateIDPResponse;
}

export namespace DeactivateIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): ReactivateIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateIDPRequest): ReactivateIDPRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateIDPRequest;
  static deserializeBinaryFromReader(message: ReactivateIDPRequest, reader: jspb.BinaryReader): ReactivateIDPRequest;
}

export namespace ReactivateIDPRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class ReactivateIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateIDPResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateIDPResponse): ReactivateIDPResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateIDPResponse;
  static deserializeBinaryFromReader(message: ReactivateIDPResponse, reader: jspb.BinaryReader): ReactivateIDPResponse;
}

export namespace ReactivateIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): RemoveIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIDPRequest): RemoveIDPRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIDPRequest;
  static deserializeBinaryFromReader(message: RemoveIDPRequest, reader: jspb.BinaryReader): RemoveIDPRequest;
}

export namespace RemoveIDPRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class RemoveIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveIDPResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIDPResponse): RemoveIDPResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIDPResponse;
  static deserializeBinaryFromReader(message: RemoveIDPResponse, reader: jspb.BinaryReader): RemoveIDPResponse;
}

export namespace RemoveIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateIDPOIDCConfigRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): UpdateIDPOIDCConfigRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateIDPOIDCConfigRequest;

  getClientId(): string;
  setClientId(value: string): UpdateIDPOIDCConfigRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateIDPOIDCConfigRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateIDPOIDCConfigRequest;
  clearScopesList(): UpdateIDPOIDCConfigRequest;
  addScopes(value: string, index?: number): UpdateIDPOIDCConfigRequest;

  getDisplayNameMapping(): zitadel_idp_pb.OIDCMappingField;
  setDisplayNameMapping(value: zitadel_idp_pb.OIDCMappingField): UpdateIDPOIDCConfigRequest;

  getUsernameMapping(): zitadel_idp_pb.OIDCMappingField;
  setUsernameMapping(value: zitadel_idp_pb.OIDCMappingField): UpdateIDPOIDCConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateIDPOIDCConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateIDPOIDCConfigRequest): UpdateIDPOIDCConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateIDPOIDCConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateIDPOIDCConfigRequest;
  static deserializeBinaryFromReader(message: UpdateIDPOIDCConfigRequest, reader: jspb.BinaryReader): UpdateIDPOIDCConfigRequest;
}

export namespace UpdateIDPOIDCConfigRequest {
  export type AsObject = {
    idpId: string,
    issuer: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    displayNameMapping: zitadel_idp_pb.OIDCMappingField,
    usernameMapping: zitadel_idp_pb.OIDCMappingField,
  }
}

export class UpdateIDPOIDCConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateIDPOIDCConfigResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateIDPOIDCConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateIDPOIDCConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateIDPOIDCConfigResponse): UpdateIDPOIDCConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateIDPOIDCConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateIDPOIDCConfigResponse;
  static deserializeBinaryFromReader(message: UpdateIDPOIDCConfigResponse, reader: jspb.BinaryReader): UpdateIDPOIDCConfigResponse;
}

export namespace UpdateIDPOIDCConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateIDPJWTConfigRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): UpdateIDPJWTConfigRequest;

  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): UpdateIDPJWTConfigRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateIDPJWTConfigRequest;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): UpdateIDPJWTConfigRequest;

  getHeaderName(): string;
  setHeaderName(value: string): UpdateIDPJWTConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateIDPJWTConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateIDPJWTConfigRequest): UpdateIDPJWTConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateIDPJWTConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateIDPJWTConfigRequest;
  static deserializeBinaryFromReader(message: UpdateIDPJWTConfigRequest, reader: jspb.BinaryReader): UpdateIDPJWTConfigRequest;
}

export namespace UpdateIDPJWTConfigRequest {
  export type AsObject = {
    idpId: string,
    jwtEndpoint: string,
    issuer: string,
    keysEndpoint: string,
    headerName: string,
  }
}

export class UpdateIDPJWTConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateIDPJWTConfigResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateIDPJWTConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateIDPJWTConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateIDPJWTConfigResponse): UpdateIDPJWTConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateIDPJWTConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateIDPJWTConfigResponse;
  static deserializeBinaryFromReader(message: UpdateIDPJWTConfigResponse, reader: jspb.BinaryReader): UpdateIDPJWTConfigResponse;
}

export namespace UpdateIDPJWTConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListProvidersRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListProvidersRequest;
  hasQuery(): boolean;
  clearQuery(): ListProvidersRequest;

  getQueriesList(): Array<ProviderQuery>;
  setQueriesList(value: Array<ProviderQuery>): ListProvidersRequest;
  clearQueriesList(): ListProvidersRequest;
  addQueries(value?: ProviderQuery, index?: number): ProviderQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProvidersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProvidersRequest): ListProvidersRequest.AsObject;
  static serializeBinaryToWriter(message: ListProvidersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProvidersRequest;
  static deserializeBinaryFromReader(message: ListProvidersRequest, reader: jspb.BinaryReader): ListProvidersRequest;
}

export namespace ListProvidersRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<ProviderQuery.AsObject>,
  }
}

export class ProviderQuery extends jspb.Message {
  getIdpIdQuery(): zitadel_idp_pb.IDPIDQuery | undefined;
  setIdpIdQuery(value?: zitadel_idp_pb.IDPIDQuery): ProviderQuery;
  hasIdpIdQuery(): boolean;
  clearIdpIdQuery(): ProviderQuery;

  getIdpNameQuery(): zitadel_idp_pb.IDPNameQuery | undefined;
  setIdpNameQuery(value?: zitadel_idp_pb.IDPNameQuery): ProviderQuery;
  hasIdpNameQuery(): boolean;
  clearIdpNameQuery(): ProviderQuery;

  getQueryCase(): ProviderQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProviderQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProviderQuery): ProviderQuery.AsObject;
  static serializeBinaryToWriter(message: ProviderQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProviderQuery;
  static deserializeBinaryFromReader(message: ProviderQuery, reader: jspb.BinaryReader): ProviderQuery;
}

export namespace ProviderQuery {
  export type AsObject = {
    idpIdQuery?: zitadel_idp_pb.IDPIDQuery.AsObject,
    idpNameQuery?: zitadel_idp_pb.IDPNameQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    IDP_ID_QUERY = 1,
    IDP_NAME_QUERY = 2,
  }
}

export class ListProvidersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProvidersResponse;
  hasDetails(): boolean;
  clearDetails(): ListProvidersResponse;

  getResultList(): Array<zitadel_idp_pb.Provider>;
  setResultList(value: Array<zitadel_idp_pb.Provider>): ListProvidersResponse;
  clearResultList(): ListProvidersResponse;
  addResult(value?: zitadel_idp_pb.Provider, index?: number): zitadel_idp_pb.Provider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProvidersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProvidersResponse): ListProvidersResponse.AsObject;
  static serializeBinaryToWriter(message: ListProvidersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProvidersResponse;
  static deserializeBinaryFromReader(message: ListProvidersResponse, reader: jspb.BinaryReader): ListProvidersResponse;
}

export namespace ListProvidersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_idp_pb.Provider.AsObject>,
  }
}

export class GetProviderByIDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetProviderByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProviderByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProviderByIDRequest): GetProviderByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetProviderByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProviderByIDRequest;
  static deserializeBinaryFromReader(message: GetProviderByIDRequest, reader: jspb.BinaryReader): GetProviderByIDRequest;
}

export namespace GetProviderByIDRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetProviderByIDResponse extends jspb.Message {
  getIdp(): zitadel_idp_pb.Provider | undefined;
  setIdp(value?: zitadel_idp_pb.Provider): GetProviderByIDResponse;
  hasIdp(): boolean;
  clearIdp(): GetProviderByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProviderByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProviderByIDResponse): GetProviderByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetProviderByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProviderByIDResponse;
  static deserializeBinaryFromReader(message: GetProviderByIDResponse, reader: jspb.BinaryReader): GetProviderByIDResponse;
}

export namespace GetProviderByIDResponse {
  export type AsObject = {
    idp?: zitadel_idp_pb.Provider.AsObject,
  }
}

export class AddGenericOAuthProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGenericOAuthProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGenericOAuthProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGenericOAuthProviderRequest;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): AddGenericOAuthProviderRequest;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): AddGenericOAuthProviderRequest;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): AddGenericOAuthProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGenericOAuthProviderRequest;
  clearScopesList(): AddGenericOAuthProviderRequest;
  addScopes(value: string, index?: number): AddGenericOAuthProviderRequest;

  getIdAttribute(): string;
  setIdAttribute(value: string): AddGenericOAuthProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGenericOAuthProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGenericOAuthProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGenericOAuthProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGenericOAuthProviderRequest): AddGenericOAuthProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGenericOAuthProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGenericOAuthProviderRequest;
  static deserializeBinaryFromReader(message: AddGenericOAuthProviderRequest, reader: jspb.BinaryReader): AddGenericOAuthProviderRequest;
}

export namespace AddGenericOAuthProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    idAttribute: string,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGenericOAuthProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGenericOAuthProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGenericOAuthProviderResponse;

  getId(): string;
  setId(value: string): AddGenericOAuthProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGenericOAuthProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGenericOAuthProviderResponse): AddGenericOAuthProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGenericOAuthProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGenericOAuthProviderResponse;
  static deserializeBinaryFromReader(message: AddGenericOAuthProviderResponse, reader: jspb.BinaryReader): AddGenericOAuthProviderResponse;
}

export namespace AddGenericOAuthProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGenericOAuthProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGenericOAuthProviderRequest;

  getName(): string;
  setName(value: string): UpdateGenericOAuthProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGenericOAuthProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGenericOAuthProviderRequest;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): UpdateGenericOAuthProviderRequest;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): UpdateGenericOAuthProviderRequest;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): UpdateGenericOAuthProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGenericOAuthProviderRequest;
  clearScopesList(): UpdateGenericOAuthProviderRequest;
  addScopes(value: string, index?: number): UpdateGenericOAuthProviderRequest;

  getIdAttribute(): string;
  setIdAttribute(value: string): UpdateGenericOAuthProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGenericOAuthProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGenericOAuthProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGenericOAuthProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGenericOAuthProviderRequest): UpdateGenericOAuthProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGenericOAuthProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGenericOAuthProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGenericOAuthProviderRequest, reader: jspb.BinaryReader): UpdateGenericOAuthProviderRequest;
}

export namespace UpdateGenericOAuthProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    idAttribute: string,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGenericOAuthProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGenericOAuthProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGenericOAuthProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGenericOAuthProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGenericOAuthProviderResponse): UpdateGenericOAuthProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGenericOAuthProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGenericOAuthProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGenericOAuthProviderResponse, reader: jspb.BinaryReader): UpdateGenericOAuthProviderResponse;
}

export namespace UpdateGenericOAuthProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGenericOIDCProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGenericOIDCProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): AddGenericOIDCProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGenericOIDCProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGenericOIDCProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGenericOIDCProviderRequest;
  clearScopesList(): AddGenericOIDCProviderRequest;
  addScopes(value: string, index?: number): AddGenericOIDCProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGenericOIDCProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGenericOIDCProviderRequest;

  getIsIdTokenMapping(): boolean;
  setIsIdTokenMapping(value: boolean): AddGenericOIDCProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGenericOIDCProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGenericOIDCProviderRequest): AddGenericOIDCProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGenericOIDCProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGenericOIDCProviderRequest;
  static deserializeBinaryFromReader(message: AddGenericOIDCProviderRequest, reader: jspb.BinaryReader): AddGenericOIDCProviderRequest;
}

export namespace AddGenericOIDCProviderRequest {
  export type AsObject = {
    name: string,
    issuer: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
    isIdTokenMapping: boolean,
  }
}

export class AddGenericOIDCProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGenericOIDCProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGenericOIDCProviderResponse;

  getId(): string;
  setId(value: string): AddGenericOIDCProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGenericOIDCProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGenericOIDCProviderResponse): AddGenericOIDCProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGenericOIDCProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGenericOIDCProviderResponse;
  static deserializeBinaryFromReader(message: AddGenericOIDCProviderResponse, reader: jspb.BinaryReader): AddGenericOIDCProviderResponse;
}

export namespace AddGenericOIDCProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGenericOIDCProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGenericOIDCProviderRequest;

  getName(): string;
  setName(value: string): UpdateGenericOIDCProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateGenericOIDCProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGenericOIDCProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGenericOIDCProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGenericOIDCProviderRequest;
  clearScopesList(): UpdateGenericOIDCProviderRequest;
  addScopes(value: string, index?: number): UpdateGenericOIDCProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGenericOIDCProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGenericOIDCProviderRequest;

  getIsIdTokenMapping(): boolean;
  setIsIdTokenMapping(value: boolean): UpdateGenericOIDCProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGenericOIDCProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGenericOIDCProviderRequest): UpdateGenericOIDCProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGenericOIDCProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGenericOIDCProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGenericOIDCProviderRequest, reader: jspb.BinaryReader): UpdateGenericOIDCProviderRequest;
}

export namespace UpdateGenericOIDCProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    issuer: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
    isIdTokenMapping: boolean,
  }
}

export class UpdateGenericOIDCProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGenericOIDCProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGenericOIDCProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGenericOIDCProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGenericOIDCProviderResponse): UpdateGenericOIDCProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGenericOIDCProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGenericOIDCProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGenericOIDCProviderResponse, reader: jspb.BinaryReader): UpdateGenericOIDCProviderResponse;
}

export namespace UpdateGenericOIDCProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class MigrateGenericOIDCProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): MigrateGenericOIDCProviderRequest;

  getAzure(): AddAzureADProviderRequest | undefined;
  setAzure(value?: AddAzureADProviderRequest): MigrateGenericOIDCProviderRequest;
  hasAzure(): boolean;
  clearAzure(): MigrateGenericOIDCProviderRequest;

  getGoogle(): AddGoogleProviderRequest | undefined;
  setGoogle(value?: AddGoogleProviderRequest): MigrateGenericOIDCProviderRequest;
  hasGoogle(): boolean;
  clearGoogle(): MigrateGenericOIDCProviderRequest;

  getTemplateCase(): MigrateGenericOIDCProviderRequest.TemplateCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrateGenericOIDCProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MigrateGenericOIDCProviderRequest): MigrateGenericOIDCProviderRequest.AsObject;
  static serializeBinaryToWriter(message: MigrateGenericOIDCProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrateGenericOIDCProviderRequest;
  static deserializeBinaryFromReader(message: MigrateGenericOIDCProviderRequest, reader: jspb.BinaryReader): MigrateGenericOIDCProviderRequest;
}

export namespace MigrateGenericOIDCProviderRequest {
  export type AsObject = {
    id: string,
    azure?: AddAzureADProviderRequest.AsObject,
    google?: AddGoogleProviderRequest.AsObject,
  }

  export enum TemplateCase { 
    TEMPLATE_NOT_SET = 0,
    AZURE = 2,
    GOOGLE = 3,
  }
}

export class MigrateGenericOIDCProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): MigrateGenericOIDCProviderResponse;
  hasDetails(): boolean;
  clearDetails(): MigrateGenericOIDCProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrateGenericOIDCProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MigrateGenericOIDCProviderResponse): MigrateGenericOIDCProviderResponse.AsObject;
  static serializeBinaryToWriter(message: MigrateGenericOIDCProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrateGenericOIDCProviderResponse;
  static deserializeBinaryFromReader(message: MigrateGenericOIDCProviderResponse, reader: jspb.BinaryReader): MigrateGenericOIDCProviderResponse;
}

export namespace MigrateGenericOIDCProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddJWTProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddJWTProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): AddJWTProviderRequest;

  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): AddJWTProviderRequest;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): AddJWTProviderRequest;

  getHeaderName(): string;
  setHeaderName(value: string): AddJWTProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddJWTProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddJWTProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddJWTProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddJWTProviderRequest): AddJWTProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddJWTProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddJWTProviderRequest;
  static deserializeBinaryFromReader(message: AddJWTProviderRequest, reader: jspb.BinaryReader): AddJWTProviderRequest;
}

export namespace AddJWTProviderRequest {
  export type AsObject = {
    name: string,
    issuer: string,
    jwtEndpoint: string,
    keysEndpoint: string,
    headerName: string,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddJWTProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddJWTProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddJWTProviderResponse;

  getId(): string;
  setId(value: string): AddJWTProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddJWTProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddJWTProviderResponse): AddJWTProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddJWTProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddJWTProviderResponse;
  static deserializeBinaryFromReader(message: AddJWTProviderResponse, reader: jspb.BinaryReader): AddJWTProviderResponse;
}

export namespace AddJWTProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateJWTProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateJWTProviderRequest;

  getName(): string;
  setName(value: string): UpdateJWTProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateJWTProviderRequest;

  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): UpdateJWTProviderRequest;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): UpdateJWTProviderRequest;

  getHeaderName(): string;
  setHeaderName(value: string): UpdateJWTProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateJWTProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateJWTProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateJWTProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateJWTProviderRequest): UpdateJWTProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateJWTProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateJWTProviderRequest;
  static deserializeBinaryFromReader(message: UpdateJWTProviderRequest, reader: jspb.BinaryReader): UpdateJWTProviderRequest;
}

export namespace UpdateJWTProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    issuer: string,
    jwtEndpoint: string,
    keysEndpoint: string,
    headerName: string,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateJWTProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateJWTProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateJWTProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateJWTProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateJWTProviderResponse): UpdateJWTProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateJWTProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateJWTProviderResponse;
  static deserializeBinaryFromReader(message: UpdateJWTProviderResponse, reader: jspb.BinaryReader): UpdateJWTProviderResponse;
}

export namespace UpdateJWTProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddAzureADProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddAzureADProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddAzureADProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddAzureADProviderRequest;

  getTenant(): zitadel_idp_pb.AzureADTenant | undefined;
  setTenant(value?: zitadel_idp_pb.AzureADTenant): AddAzureADProviderRequest;
  hasTenant(): boolean;
  clearTenant(): AddAzureADProviderRequest;

  getEmailVerified(): boolean;
  setEmailVerified(value: boolean): AddAzureADProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddAzureADProviderRequest;
  clearScopesList(): AddAzureADProviderRequest;
  addScopes(value: string, index?: number): AddAzureADProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddAzureADProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddAzureADProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAzureADProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddAzureADProviderRequest): AddAzureADProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddAzureADProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAzureADProviderRequest;
  static deserializeBinaryFromReader(message: AddAzureADProviderRequest, reader: jspb.BinaryReader): AddAzureADProviderRequest;
}

export namespace AddAzureADProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    tenant?: zitadel_idp_pb.AzureADTenant.AsObject,
    emailVerified: boolean,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddAzureADProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddAzureADProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddAzureADProviderResponse;

  getId(): string;
  setId(value: string): AddAzureADProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAzureADProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddAzureADProviderResponse): AddAzureADProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddAzureADProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAzureADProviderResponse;
  static deserializeBinaryFromReader(message: AddAzureADProviderResponse, reader: jspb.BinaryReader): AddAzureADProviderResponse;
}

export namespace AddAzureADProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateAzureADProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateAzureADProviderRequest;

  getName(): string;
  setName(value: string): UpdateAzureADProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateAzureADProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateAzureADProviderRequest;

  getTenant(): zitadel_idp_pb.AzureADTenant | undefined;
  setTenant(value?: zitadel_idp_pb.AzureADTenant): UpdateAzureADProviderRequest;
  hasTenant(): boolean;
  clearTenant(): UpdateAzureADProviderRequest;

  getEmailVerified(): boolean;
  setEmailVerified(value: boolean): UpdateAzureADProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateAzureADProviderRequest;
  clearScopesList(): UpdateAzureADProviderRequest;
  addScopes(value: string, index?: number): UpdateAzureADProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateAzureADProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateAzureADProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAzureADProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAzureADProviderRequest): UpdateAzureADProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAzureADProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAzureADProviderRequest;
  static deserializeBinaryFromReader(message: UpdateAzureADProviderRequest, reader: jspb.BinaryReader): UpdateAzureADProviderRequest;
}

export namespace UpdateAzureADProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    tenant?: zitadel_idp_pb.AzureADTenant.AsObject,
    emailVerified: boolean,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateAzureADProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateAzureADProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateAzureADProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAzureADProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAzureADProviderResponse): UpdateAzureADProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateAzureADProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAzureADProviderResponse;
  static deserializeBinaryFromReader(message: UpdateAzureADProviderResponse, reader: jspb.BinaryReader): UpdateAzureADProviderResponse;
}

export namespace UpdateAzureADProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGitHubProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGitHubProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGitHubProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGitHubProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGitHubProviderRequest;
  clearScopesList(): AddGitHubProviderRequest;
  addScopes(value: string, index?: number): AddGitHubProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGitHubProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGitHubProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitHubProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitHubProviderRequest): AddGitHubProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGitHubProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitHubProviderRequest;
  static deserializeBinaryFromReader(message: AddGitHubProviderRequest, reader: jspb.BinaryReader): AddGitHubProviderRequest;
}

export namespace AddGitHubProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGitHubProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGitHubProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGitHubProviderResponse;

  getId(): string;
  setId(value: string): AddGitHubProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitHubProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitHubProviderResponse): AddGitHubProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGitHubProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitHubProviderResponse;
  static deserializeBinaryFromReader(message: AddGitHubProviderResponse, reader: jspb.BinaryReader): AddGitHubProviderResponse;
}

export namespace AddGitHubProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGitHubProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGitHubProviderRequest;

  getName(): string;
  setName(value: string): UpdateGitHubProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGitHubProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGitHubProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGitHubProviderRequest;
  clearScopesList(): UpdateGitHubProviderRequest;
  addScopes(value: string, index?: number): UpdateGitHubProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGitHubProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGitHubProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitHubProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitHubProviderRequest): UpdateGitHubProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGitHubProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitHubProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGitHubProviderRequest, reader: jspb.BinaryReader): UpdateGitHubProviderRequest;
}

export namespace UpdateGitHubProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGitHubProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGitHubProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGitHubProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitHubProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitHubProviderResponse): UpdateGitHubProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGitHubProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitHubProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGitHubProviderResponse, reader: jspb.BinaryReader): UpdateGitHubProviderResponse;
}

export namespace UpdateGitHubProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGitHubEnterpriseServerProviderRequest extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): AddGitHubEnterpriseServerProviderRequest;

  getName(): string;
  setName(value: string): AddGitHubEnterpriseServerProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGitHubEnterpriseServerProviderRequest;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): AddGitHubEnterpriseServerProviderRequest;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): AddGitHubEnterpriseServerProviderRequest;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): AddGitHubEnterpriseServerProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGitHubEnterpriseServerProviderRequest;
  clearScopesList(): AddGitHubEnterpriseServerProviderRequest;
  addScopes(value: string, index?: number): AddGitHubEnterpriseServerProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGitHubEnterpriseServerProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGitHubEnterpriseServerProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitHubEnterpriseServerProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitHubEnterpriseServerProviderRequest): AddGitHubEnterpriseServerProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGitHubEnterpriseServerProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitHubEnterpriseServerProviderRequest;
  static deserializeBinaryFromReader(message: AddGitHubEnterpriseServerProviderRequest, reader: jspb.BinaryReader): AddGitHubEnterpriseServerProviderRequest;
}

export namespace AddGitHubEnterpriseServerProviderRequest {
  export type AsObject = {
    clientId: string,
    name: string,
    clientSecret: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGitHubEnterpriseServerProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGitHubEnterpriseServerProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGitHubEnterpriseServerProviderResponse;

  getId(): string;
  setId(value: string): AddGitHubEnterpriseServerProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitHubEnterpriseServerProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitHubEnterpriseServerProviderResponse): AddGitHubEnterpriseServerProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGitHubEnterpriseServerProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitHubEnterpriseServerProviderResponse;
  static deserializeBinaryFromReader(message: AddGitHubEnterpriseServerProviderResponse, reader: jspb.BinaryReader): AddGitHubEnterpriseServerProviderResponse;
}

export namespace AddGitHubEnterpriseServerProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGitHubEnterpriseServerProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getName(): string;
  setName(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGitHubEnterpriseServerProviderRequest;
  clearScopesList(): UpdateGitHubEnterpriseServerProviderRequest;
  addScopes(value: string, index?: number): UpdateGitHubEnterpriseServerProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGitHubEnterpriseServerProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGitHubEnterpriseServerProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitHubEnterpriseServerProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitHubEnterpriseServerProviderRequest): UpdateGitHubEnterpriseServerProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGitHubEnterpriseServerProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitHubEnterpriseServerProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGitHubEnterpriseServerProviderRequest, reader: jspb.BinaryReader): UpdateGitHubEnterpriseServerProviderRequest;
}

export namespace UpdateGitHubEnterpriseServerProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGitHubEnterpriseServerProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGitHubEnterpriseServerProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGitHubEnterpriseServerProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitHubEnterpriseServerProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitHubEnterpriseServerProviderResponse): UpdateGitHubEnterpriseServerProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGitHubEnterpriseServerProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitHubEnterpriseServerProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGitHubEnterpriseServerProviderResponse, reader: jspb.BinaryReader): UpdateGitHubEnterpriseServerProviderResponse;
}

export namespace UpdateGitHubEnterpriseServerProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGitLabProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGitLabProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGitLabProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGitLabProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGitLabProviderRequest;
  clearScopesList(): AddGitLabProviderRequest;
  addScopes(value: string, index?: number): AddGitLabProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGitLabProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGitLabProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitLabProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitLabProviderRequest): AddGitLabProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGitLabProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitLabProviderRequest;
  static deserializeBinaryFromReader(message: AddGitLabProviderRequest, reader: jspb.BinaryReader): AddGitLabProviderRequest;
}

export namespace AddGitLabProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGitLabProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGitLabProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGitLabProviderResponse;

  getId(): string;
  setId(value: string): AddGitLabProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitLabProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitLabProviderResponse): AddGitLabProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGitLabProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitLabProviderResponse;
  static deserializeBinaryFromReader(message: AddGitLabProviderResponse, reader: jspb.BinaryReader): AddGitLabProviderResponse;
}

export namespace AddGitLabProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGitLabProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGitLabProviderRequest;

  getName(): string;
  setName(value: string): UpdateGitLabProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGitLabProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGitLabProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGitLabProviderRequest;
  clearScopesList(): UpdateGitLabProviderRequest;
  addScopes(value: string, index?: number): UpdateGitLabProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGitLabProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGitLabProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitLabProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitLabProviderRequest): UpdateGitLabProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGitLabProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitLabProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGitLabProviderRequest, reader: jspb.BinaryReader): UpdateGitLabProviderRequest;
}

export namespace UpdateGitLabProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGitLabProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGitLabProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGitLabProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitLabProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitLabProviderResponse): UpdateGitLabProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGitLabProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitLabProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGitLabProviderResponse, reader: jspb.BinaryReader): UpdateGitLabProviderResponse;
}

export namespace UpdateGitLabProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGitLabSelfHostedProviderRequest extends jspb.Message {
  getIssuer(): string;
  setIssuer(value: string): AddGitLabSelfHostedProviderRequest;

  getName(): string;
  setName(value: string): AddGitLabSelfHostedProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGitLabSelfHostedProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGitLabSelfHostedProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGitLabSelfHostedProviderRequest;
  clearScopesList(): AddGitLabSelfHostedProviderRequest;
  addScopes(value: string, index?: number): AddGitLabSelfHostedProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGitLabSelfHostedProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGitLabSelfHostedProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitLabSelfHostedProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitLabSelfHostedProviderRequest): AddGitLabSelfHostedProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGitLabSelfHostedProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitLabSelfHostedProviderRequest;
  static deserializeBinaryFromReader(message: AddGitLabSelfHostedProviderRequest, reader: jspb.BinaryReader): AddGitLabSelfHostedProviderRequest;
}

export namespace AddGitLabSelfHostedProviderRequest {
  export type AsObject = {
    issuer: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGitLabSelfHostedProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGitLabSelfHostedProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGitLabSelfHostedProviderResponse;

  getId(): string;
  setId(value: string): AddGitLabSelfHostedProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitLabSelfHostedProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitLabSelfHostedProviderResponse): AddGitLabSelfHostedProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGitLabSelfHostedProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitLabSelfHostedProviderResponse;
  static deserializeBinaryFromReader(message: AddGitLabSelfHostedProviderResponse, reader: jspb.BinaryReader): AddGitLabSelfHostedProviderResponse;
}

export namespace AddGitLabSelfHostedProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGitLabSelfHostedProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGitLabSelfHostedProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateGitLabSelfHostedProviderRequest;

  getName(): string;
  setName(value: string): UpdateGitLabSelfHostedProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGitLabSelfHostedProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGitLabSelfHostedProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGitLabSelfHostedProviderRequest;
  clearScopesList(): UpdateGitLabSelfHostedProviderRequest;
  addScopes(value: string, index?: number): UpdateGitLabSelfHostedProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGitLabSelfHostedProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGitLabSelfHostedProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitLabSelfHostedProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitLabSelfHostedProviderRequest): UpdateGitLabSelfHostedProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGitLabSelfHostedProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitLabSelfHostedProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGitLabSelfHostedProviderRequest, reader: jspb.BinaryReader): UpdateGitLabSelfHostedProviderRequest;
}

export namespace UpdateGitLabSelfHostedProviderRequest {
  export type AsObject = {
    id: string,
    issuer: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGitLabSelfHostedProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGitLabSelfHostedProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGitLabSelfHostedProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitLabSelfHostedProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitLabSelfHostedProviderResponse): UpdateGitLabSelfHostedProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGitLabSelfHostedProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitLabSelfHostedProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGitLabSelfHostedProviderResponse, reader: jspb.BinaryReader): UpdateGitLabSelfHostedProviderResponse;
}

export namespace UpdateGitLabSelfHostedProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGoogleProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGoogleProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGoogleProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGoogleProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGoogleProviderRequest;
  clearScopesList(): AddGoogleProviderRequest;
  addScopes(value: string, index?: number): AddGoogleProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGoogleProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGoogleProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGoogleProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGoogleProviderRequest): AddGoogleProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGoogleProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGoogleProviderRequest;
  static deserializeBinaryFromReader(message: AddGoogleProviderRequest, reader: jspb.BinaryReader): AddGoogleProviderRequest;
}

export namespace AddGoogleProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGoogleProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGoogleProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGoogleProviderResponse;

  getId(): string;
  setId(value: string): AddGoogleProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGoogleProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGoogleProviderResponse): AddGoogleProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGoogleProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGoogleProviderResponse;
  static deserializeBinaryFromReader(message: AddGoogleProviderResponse, reader: jspb.BinaryReader): AddGoogleProviderResponse;
}

export namespace AddGoogleProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGoogleProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGoogleProviderRequest;

  getName(): string;
  setName(value: string): UpdateGoogleProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGoogleProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGoogleProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGoogleProviderRequest;
  clearScopesList(): UpdateGoogleProviderRequest;
  addScopes(value: string, index?: number): UpdateGoogleProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGoogleProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGoogleProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGoogleProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGoogleProviderRequest): UpdateGoogleProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGoogleProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGoogleProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGoogleProviderRequest, reader: jspb.BinaryReader): UpdateGoogleProviderRequest;
}

export namespace UpdateGoogleProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGoogleProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGoogleProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGoogleProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGoogleProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGoogleProviderResponse): UpdateGoogleProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGoogleProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGoogleProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGoogleProviderResponse, reader: jspb.BinaryReader): UpdateGoogleProviderResponse;
}

export namespace UpdateGoogleProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddLDAPProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddLDAPProviderRequest;

  getServersList(): Array<string>;
  setServersList(value: Array<string>): AddLDAPProviderRequest;
  clearServersList(): AddLDAPProviderRequest;
  addServers(value: string, index?: number): AddLDAPProviderRequest;

  getStartTls(): boolean;
  setStartTls(value: boolean): AddLDAPProviderRequest;

  getBaseDn(): string;
  setBaseDn(value: string): AddLDAPProviderRequest;

  getBindDn(): string;
  setBindDn(value: string): AddLDAPProviderRequest;

  getBindPassword(): string;
  setBindPassword(value: string): AddLDAPProviderRequest;

  getUserBase(): string;
  setUserBase(value: string): AddLDAPProviderRequest;

  getUserObjectClassesList(): Array<string>;
  setUserObjectClassesList(value: Array<string>): AddLDAPProviderRequest;
  clearUserObjectClassesList(): AddLDAPProviderRequest;
  addUserObjectClasses(value: string, index?: number): AddLDAPProviderRequest;

  getUserFiltersList(): Array<string>;
  setUserFiltersList(value: Array<string>): AddLDAPProviderRequest;
  clearUserFiltersList(): AddLDAPProviderRequest;
  addUserFilters(value: string, index?: number): AddLDAPProviderRequest;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): AddLDAPProviderRequest;
  hasTimeout(): boolean;
  clearTimeout(): AddLDAPProviderRequest;

  getAttributes(): zitadel_idp_pb.LDAPAttributes | undefined;
  setAttributes(value?: zitadel_idp_pb.LDAPAttributes): AddLDAPProviderRequest;
  hasAttributes(): boolean;
  clearAttributes(): AddLDAPProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddLDAPProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddLDAPProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddLDAPProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddLDAPProviderRequest): AddLDAPProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddLDAPProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddLDAPProviderRequest;
  static deserializeBinaryFromReader(message: AddLDAPProviderRequest, reader: jspb.BinaryReader): AddLDAPProviderRequest;
}

export namespace AddLDAPProviderRequest {
  export type AsObject = {
    name: string,
    serversList: Array<string>,
    startTls: boolean,
    baseDn: string,
    bindDn: string,
    bindPassword: string,
    userBase: string,
    userObjectClassesList: Array<string>,
    userFiltersList: Array<string>,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    attributes?: zitadel_idp_pb.LDAPAttributes.AsObject,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddLDAPProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddLDAPProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddLDAPProviderResponse;

  getId(): string;
  setId(value: string): AddLDAPProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddLDAPProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddLDAPProviderResponse): AddLDAPProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddLDAPProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddLDAPProviderResponse;
  static deserializeBinaryFromReader(message: AddLDAPProviderResponse, reader: jspb.BinaryReader): AddLDAPProviderResponse;
}

export namespace AddLDAPProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateLDAPProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateLDAPProviderRequest;

  getName(): string;
  setName(value: string): UpdateLDAPProviderRequest;

  getServersList(): Array<string>;
  setServersList(value: Array<string>): UpdateLDAPProviderRequest;
  clearServersList(): UpdateLDAPProviderRequest;
  addServers(value: string, index?: number): UpdateLDAPProviderRequest;

  getStartTls(): boolean;
  setStartTls(value: boolean): UpdateLDAPProviderRequest;

  getBaseDn(): string;
  setBaseDn(value: string): UpdateLDAPProviderRequest;

  getBindDn(): string;
  setBindDn(value: string): UpdateLDAPProviderRequest;

  getBindPassword(): string;
  setBindPassword(value: string): UpdateLDAPProviderRequest;

  getUserBase(): string;
  setUserBase(value: string): UpdateLDAPProviderRequest;

  getUserObjectClassesList(): Array<string>;
  setUserObjectClassesList(value: Array<string>): UpdateLDAPProviderRequest;
  clearUserObjectClassesList(): UpdateLDAPProviderRequest;
  addUserObjectClasses(value: string, index?: number): UpdateLDAPProviderRequest;

  getUserFiltersList(): Array<string>;
  setUserFiltersList(value: Array<string>): UpdateLDAPProviderRequest;
  clearUserFiltersList(): UpdateLDAPProviderRequest;
  addUserFilters(value: string, index?: number): UpdateLDAPProviderRequest;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): UpdateLDAPProviderRequest;
  hasTimeout(): boolean;
  clearTimeout(): UpdateLDAPProviderRequest;

  getAttributes(): zitadel_idp_pb.LDAPAttributes | undefined;
  setAttributes(value?: zitadel_idp_pb.LDAPAttributes): UpdateLDAPProviderRequest;
  hasAttributes(): boolean;
  clearAttributes(): UpdateLDAPProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateLDAPProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateLDAPProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLDAPProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLDAPProviderRequest): UpdateLDAPProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateLDAPProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLDAPProviderRequest;
  static deserializeBinaryFromReader(message: UpdateLDAPProviderRequest, reader: jspb.BinaryReader): UpdateLDAPProviderRequest;
}

export namespace UpdateLDAPProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    serversList: Array<string>,
    startTls: boolean,
    baseDn: string,
    bindDn: string,
    bindPassword: string,
    userBase: string,
    userObjectClassesList: Array<string>,
    userFiltersList: Array<string>,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    attributes?: zitadel_idp_pb.LDAPAttributes.AsObject,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateLDAPProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateLDAPProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateLDAPProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLDAPProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLDAPProviderResponse): UpdateLDAPProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateLDAPProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLDAPProviderResponse;
  static deserializeBinaryFromReader(message: UpdateLDAPProviderResponse, reader: jspb.BinaryReader): UpdateLDAPProviderResponse;
}

export namespace UpdateLDAPProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddAppleProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddAppleProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddAppleProviderRequest;

  getTeamId(): string;
  setTeamId(value: string): AddAppleProviderRequest;

  getKeyId(): string;
  setKeyId(value: string): AddAppleProviderRequest;

  getPrivateKey(): Uint8Array | string;
  getPrivateKey_asU8(): Uint8Array;
  getPrivateKey_asB64(): string;
  setPrivateKey(value: Uint8Array | string): AddAppleProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddAppleProviderRequest;
  clearScopesList(): AddAppleProviderRequest;
  addScopes(value: string, index?: number): AddAppleProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddAppleProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddAppleProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAppleProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddAppleProviderRequest): AddAppleProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddAppleProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAppleProviderRequest;
  static deserializeBinaryFromReader(message: AddAppleProviderRequest, reader: jspb.BinaryReader): AddAppleProviderRequest;
}

export namespace AddAppleProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    teamId: string,
    keyId: string,
    privateKey: Uint8Array | string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddAppleProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddAppleProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddAppleProviderResponse;

  getId(): string;
  setId(value: string): AddAppleProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAppleProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddAppleProviderResponse): AddAppleProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddAppleProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAppleProviderResponse;
  static deserializeBinaryFromReader(message: AddAppleProviderResponse, reader: jspb.BinaryReader): AddAppleProviderResponse;
}

export namespace AddAppleProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateAppleProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateAppleProviderRequest;

  getName(): string;
  setName(value: string): UpdateAppleProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateAppleProviderRequest;

  getTeamId(): string;
  setTeamId(value: string): UpdateAppleProviderRequest;

  getKeyId(): string;
  setKeyId(value: string): UpdateAppleProviderRequest;

  getPrivateKey(): Uint8Array | string;
  getPrivateKey_asU8(): Uint8Array;
  getPrivateKey_asB64(): string;
  setPrivateKey(value: Uint8Array | string): UpdateAppleProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateAppleProviderRequest;
  clearScopesList(): UpdateAppleProviderRequest;
  addScopes(value: string, index?: number): UpdateAppleProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateAppleProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateAppleProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAppleProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAppleProviderRequest): UpdateAppleProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAppleProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAppleProviderRequest;
  static deserializeBinaryFromReader(message: UpdateAppleProviderRequest, reader: jspb.BinaryReader): UpdateAppleProviderRequest;
}

export namespace UpdateAppleProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    teamId: string,
    keyId: string,
    privateKey: Uint8Array | string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateAppleProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateAppleProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateAppleProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAppleProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAppleProviderResponse): UpdateAppleProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateAppleProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAppleProviderResponse;
  static deserializeBinaryFromReader(message: UpdateAppleProviderResponse, reader: jspb.BinaryReader): UpdateAppleProviderResponse;
}

export namespace UpdateAppleProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddSAMLProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddSAMLProviderRequest;

  getMetadataXml(): Uint8Array | string;
  getMetadataXml_asU8(): Uint8Array;
  getMetadataXml_asB64(): string;
  setMetadataXml(value: Uint8Array | string): AddSAMLProviderRequest;

  getMetadataUrl(): string;
  setMetadataUrl(value: string): AddSAMLProviderRequest;

  getBinding(): zitadel_idp_pb.SAMLBinding;
  setBinding(value: zitadel_idp_pb.SAMLBinding): AddSAMLProviderRequest;

  getWithSignedRequest(): boolean;
  setWithSignedRequest(value: boolean): AddSAMLProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddSAMLProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddSAMLProviderRequest;

  getNameIdFormat(): zitadel_idp_pb.SAMLNameIDFormat;
  setNameIdFormat(value: zitadel_idp_pb.SAMLNameIDFormat): AddSAMLProviderRequest;
  hasNameIdFormat(): boolean;
  clearNameIdFormat(): AddSAMLProviderRequest;

  getTransientMappingAttributeName(): string;
  setTransientMappingAttributeName(value: string): AddSAMLProviderRequest;
  hasTransientMappingAttributeName(): boolean;
  clearTransientMappingAttributeName(): AddSAMLProviderRequest;

  getMetadataCase(): AddSAMLProviderRequest.MetadataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSAMLProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddSAMLProviderRequest): AddSAMLProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddSAMLProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSAMLProviderRequest;
  static deserializeBinaryFromReader(message: AddSAMLProviderRequest, reader: jspb.BinaryReader): AddSAMLProviderRequest;
}

export namespace AddSAMLProviderRequest {
  export type AsObject = {
    name: string,
    metadataXml: Uint8Array | string,
    metadataUrl: string,
    binding: zitadel_idp_pb.SAMLBinding,
    withSignedRequest: boolean,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
    nameIdFormat?: zitadel_idp_pb.SAMLNameIDFormat,
    transientMappingAttributeName?: string,
  }

  export enum MetadataCase { 
    METADATA_NOT_SET = 0,
    METADATA_XML = 2,
    METADATA_URL = 3,
  }

  export enum NameIdFormatCase { 
    _NAME_ID_FORMAT_NOT_SET = 0,
    NAME_ID_FORMAT = 7,
  }

  export enum TransientMappingAttributeNameCase { 
    _TRANSIENT_MAPPING_ATTRIBUTE_NAME_NOT_SET = 0,
    TRANSIENT_MAPPING_ATTRIBUTE_NAME = 8,
  }
}

export class AddSAMLProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddSAMLProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddSAMLProviderResponse;

  getId(): string;
  setId(value: string): AddSAMLProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSAMLProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSAMLProviderResponse): AddSAMLProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddSAMLProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSAMLProviderResponse;
  static deserializeBinaryFromReader(message: AddSAMLProviderResponse, reader: jspb.BinaryReader): AddSAMLProviderResponse;
}

export namespace AddSAMLProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateSAMLProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateSAMLProviderRequest;

  getName(): string;
  setName(value: string): UpdateSAMLProviderRequest;

  getMetadataXml(): Uint8Array | string;
  getMetadataXml_asU8(): Uint8Array;
  getMetadataXml_asB64(): string;
  setMetadataXml(value: Uint8Array | string): UpdateSAMLProviderRequest;

  getMetadataUrl(): string;
  setMetadataUrl(value: string): UpdateSAMLProviderRequest;

  getBinding(): zitadel_idp_pb.SAMLBinding;
  setBinding(value: zitadel_idp_pb.SAMLBinding): UpdateSAMLProviderRequest;

  getWithSignedRequest(): boolean;
  setWithSignedRequest(value: boolean): UpdateSAMLProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateSAMLProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateSAMLProviderRequest;

  getNameIdFormat(): zitadel_idp_pb.SAMLNameIDFormat;
  setNameIdFormat(value: zitadel_idp_pb.SAMLNameIDFormat): UpdateSAMLProviderRequest;
  hasNameIdFormat(): boolean;
  clearNameIdFormat(): UpdateSAMLProviderRequest;

  getTransientMappingAttributeName(): string;
  setTransientMappingAttributeName(value: string): UpdateSAMLProviderRequest;
  hasTransientMappingAttributeName(): boolean;
  clearTransientMappingAttributeName(): UpdateSAMLProviderRequest;

  getMetadataCase(): UpdateSAMLProviderRequest.MetadataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSAMLProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSAMLProviderRequest): UpdateSAMLProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateSAMLProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSAMLProviderRequest;
  static deserializeBinaryFromReader(message: UpdateSAMLProviderRequest, reader: jspb.BinaryReader): UpdateSAMLProviderRequest;
}

export namespace UpdateSAMLProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    metadataXml: Uint8Array | string,
    metadataUrl: string,
    binding: zitadel_idp_pb.SAMLBinding,
    withSignedRequest: boolean,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
    nameIdFormat?: zitadel_idp_pb.SAMLNameIDFormat,
    transientMappingAttributeName?: string,
  }

  export enum MetadataCase { 
    METADATA_NOT_SET = 0,
    METADATA_XML = 3,
    METADATA_URL = 4,
  }

  export enum NameIdFormatCase { 
    _NAME_ID_FORMAT_NOT_SET = 0,
    NAME_ID_FORMAT = 8,
  }

  export enum TransientMappingAttributeNameCase { 
    _TRANSIENT_MAPPING_ATTRIBUTE_NAME_NOT_SET = 0,
    TRANSIENT_MAPPING_ATTRIBUTE_NAME = 9,
  }
}

export class UpdateSAMLProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateSAMLProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateSAMLProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSAMLProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSAMLProviderResponse): UpdateSAMLProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateSAMLProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSAMLProviderResponse;
  static deserializeBinaryFromReader(message: UpdateSAMLProviderResponse, reader: jspb.BinaryReader): UpdateSAMLProviderResponse;
}

export namespace UpdateSAMLProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RegenerateSAMLProviderCertificateRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RegenerateSAMLProviderCertificateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegenerateSAMLProviderCertificateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegenerateSAMLProviderCertificateRequest): RegenerateSAMLProviderCertificateRequest.AsObject;
  static serializeBinaryToWriter(message: RegenerateSAMLProviderCertificateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegenerateSAMLProviderCertificateRequest;
  static deserializeBinaryFromReader(message: RegenerateSAMLProviderCertificateRequest, reader: jspb.BinaryReader): RegenerateSAMLProviderCertificateRequest;
}

export namespace RegenerateSAMLProviderCertificateRequest {
  export type AsObject = {
    id: string,
  }
}

export class RegenerateSAMLProviderCertificateResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RegenerateSAMLProviderCertificateResponse;
  hasDetails(): boolean;
  clearDetails(): RegenerateSAMLProviderCertificateResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegenerateSAMLProviderCertificateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegenerateSAMLProviderCertificateResponse): RegenerateSAMLProviderCertificateResponse.AsObject;
  static serializeBinaryToWriter(message: RegenerateSAMLProviderCertificateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegenerateSAMLProviderCertificateResponse;
  static deserializeBinaryFromReader(message: RegenerateSAMLProviderCertificateResponse, reader: jspb.BinaryReader): RegenerateSAMLProviderCertificateResponse;
}

export namespace RegenerateSAMLProviderCertificateResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeleteProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProviderRequest): DeleteProviderRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProviderRequest;
  static deserializeBinaryFromReader(message: DeleteProviderRequest, reader: jspb.BinaryReader): DeleteProviderRequest;
}

export namespace DeleteProviderRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeleteProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeleteProviderResponse;
  hasDetails(): boolean;
  clearDetails(): DeleteProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProviderResponse): DeleteProviderResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProviderResponse;
  static deserializeBinaryFromReader(message: DeleteProviderResponse, reader: jspb.BinaryReader): DeleteProviderResponse;
}

export namespace DeleteProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetOrgIAMPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgIAMPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgIAMPolicyRequest): GetOrgIAMPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrgIAMPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgIAMPolicyRequest;
  static deserializeBinaryFromReader(message: GetOrgIAMPolicyRequest, reader: jspb.BinaryReader): GetOrgIAMPolicyRequest;
}

export namespace GetOrgIAMPolicyRequest {
  export type AsObject = {
  }
}

export class GetOrgIAMPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.OrgIAMPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.OrgIAMPolicy): GetOrgIAMPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetOrgIAMPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgIAMPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgIAMPolicyResponse): GetOrgIAMPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrgIAMPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgIAMPolicyResponse;
  static deserializeBinaryFromReader(message: GetOrgIAMPolicyResponse, reader: jspb.BinaryReader): GetOrgIAMPolicyResponse;
}

export namespace GetOrgIAMPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.OrgIAMPolicy.AsObject,
  }
}

export class UpdateOrgIAMPolicyRequest extends jspb.Message {
  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): UpdateOrgIAMPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgIAMPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgIAMPolicyRequest): UpdateOrgIAMPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgIAMPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgIAMPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateOrgIAMPolicyRequest, reader: jspb.BinaryReader): UpdateOrgIAMPolicyRequest;
}

export namespace UpdateOrgIAMPolicyRequest {
  export type AsObject = {
    userLoginMustBeDomain: boolean,
  }
}

export class UpdateOrgIAMPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateOrgIAMPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateOrgIAMPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgIAMPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgIAMPolicyResponse): UpdateOrgIAMPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgIAMPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgIAMPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateOrgIAMPolicyResponse, reader: jspb.BinaryReader): UpdateOrgIAMPolicyResponse;
}

export namespace UpdateOrgIAMPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomOrgIAMPolicyRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): GetCustomOrgIAMPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomOrgIAMPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomOrgIAMPolicyRequest): GetCustomOrgIAMPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomOrgIAMPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomOrgIAMPolicyRequest;
  static deserializeBinaryFromReader(message: GetCustomOrgIAMPolicyRequest, reader: jspb.BinaryReader): GetCustomOrgIAMPolicyRequest;
}

export namespace GetCustomOrgIAMPolicyRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class GetCustomOrgIAMPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.OrgIAMPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.OrgIAMPolicy): GetCustomOrgIAMPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetCustomOrgIAMPolicyResponse;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): GetCustomOrgIAMPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomOrgIAMPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomOrgIAMPolicyResponse): GetCustomOrgIAMPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomOrgIAMPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomOrgIAMPolicyResponse;
  static deserializeBinaryFromReader(message: GetCustomOrgIAMPolicyResponse, reader: jspb.BinaryReader): GetCustomOrgIAMPolicyResponse;
}

export namespace GetCustomOrgIAMPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.OrgIAMPolicy.AsObject,
    isDefault: boolean,
  }
}

export class AddCustomOrgIAMPolicyRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): AddCustomOrgIAMPolicyRequest;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): AddCustomOrgIAMPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomOrgIAMPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomOrgIAMPolicyRequest): AddCustomOrgIAMPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomOrgIAMPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomOrgIAMPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomOrgIAMPolicyRequest, reader: jspb.BinaryReader): AddCustomOrgIAMPolicyRequest;
}

export namespace AddCustomOrgIAMPolicyRequest {
  export type AsObject = {
    orgId: string,
    userLoginMustBeDomain: boolean,
  }
}

export class AddCustomOrgIAMPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomOrgIAMPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomOrgIAMPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomOrgIAMPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomOrgIAMPolicyResponse): AddCustomOrgIAMPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomOrgIAMPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomOrgIAMPolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomOrgIAMPolicyResponse, reader: jspb.BinaryReader): AddCustomOrgIAMPolicyResponse;
}

export namespace AddCustomOrgIAMPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomOrgIAMPolicyRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): UpdateCustomOrgIAMPolicyRequest;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): UpdateCustomOrgIAMPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomOrgIAMPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomOrgIAMPolicyRequest): UpdateCustomOrgIAMPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomOrgIAMPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomOrgIAMPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomOrgIAMPolicyRequest, reader: jspb.BinaryReader): UpdateCustomOrgIAMPolicyRequest;
}

export namespace UpdateCustomOrgIAMPolicyRequest {
  export type AsObject = {
    orgId: string,
    userLoginMustBeDomain: boolean,
  }
}

export class UpdateCustomOrgIAMPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomOrgIAMPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomOrgIAMPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomOrgIAMPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomOrgIAMPolicyResponse): UpdateCustomOrgIAMPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomOrgIAMPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomOrgIAMPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomOrgIAMPolicyResponse, reader: jspb.BinaryReader): UpdateCustomOrgIAMPolicyResponse;
}

export namespace UpdateCustomOrgIAMPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomOrgIAMPolicyToDefaultRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): ResetCustomOrgIAMPolicyToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomOrgIAMPolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomOrgIAMPolicyToDefaultRequest): ResetCustomOrgIAMPolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomOrgIAMPolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomOrgIAMPolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomOrgIAMPolicyToDefaultRequest, reader: jspb.BinaryReader): ResetCustomOrgIAMPolicyToDefaultRequest;
}

export namespace ResetCustomOrgIAMPolicyToDefaultRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class ResetCustomOrgIAMPolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomOrgIAMPolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomOrgIAMPolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomOrgIAMPolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomOrgIAMPolicyToDefaultResponse): ResetCustomOrgIAMPolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomOrgIAMPolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomOrgIAMPolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomOrgIAMPolicyToDefaultResponse, reader: jspb.BinaryReader): ResetCustomOrgIAMPolicyToDefaultResponse;
}

export namespace ResetCustomOrgIAMPolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDomainPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDomainPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDomainPolicyRequest): GetDomainPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDomainPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDomainPolicyRequest;
  static deserializeBinaryFromReader(message: GetDomainPolicyRequest, reader: jspb.BinaryReader): GetDomainPolicyRequest;
}

export namespace GetDomainPolicyRequest {
  export type AsObject = {
  }
}

export class GetDomainPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.DomainPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.DomainPolicy): GetDomainPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDomainPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDomainPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDomainPolicyResponse): GetDomainPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDomainPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDomainPolicyResponse;
  static deserializeBinaryFromReader(message: GetDomainPolicyResponse, reader: jspb.BinaryReader): GetDomainPolicyResponse;
}

export namespace GetDomainPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.DomainPolicy.AsObject,
  }
}

export class UpdateDomainPolicyRequest extends jspb.Message {
  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): UpdateDomainPolicyRequest;

  getValidateOrgDomains(): boolean;
  setValidateOrgDomains(value: boolean): UpdateDomainPolicyRequest;

  getSmtpSenderAddressMatchesInstanceDomain(): boolean;
  setSmtpSenderAddressMatchesInstanceDomain(value: boolean): UpdateDomainPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateDomainPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateDomainPolicyRequest): UpdateDomainPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateDomainPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateDomainPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateDomainPolicyRequest, reader: jspb.BinaryReader): UpdateDomainPolicyRequest;
}

export namespace UpdateDomainPolicyRequest {
  export type AsObject = {
    userLoginMustBeDomain: boolean,
    validateOrgDomains: boolean,
    smtpSenderAddressMatchesInstanceDomain: boolean,
  }
}

export class UpdateDomainPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateDomainPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateDomainPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateDomainPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateDomainPolicyResponse): UpdateDomainPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateDomainPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateDomainPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateDomainPolicyResponse, reader: jspb.BinaryReader): UpdateDomainPolicyResponse;
}

export namespace UpdateDomainPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomDomainPolicyRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): GetCustomDomainPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomDomainPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomDomainPolicyRequest): GetCustomDomainPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomDomainPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomDomainPolicyRequest;
  static deserializeBinaryFromReader(message: GetCustomDomainPolicyRequest, reader: jspb.BinaryReader): GetCustomDomainPolicyRequest;
}

export namespace GetCustomDomainPolicyRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class GetCustomDomainPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.DomainPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.DomainPolicy): GetCustomDomainPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetCustomDomainPolicyResponse;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): GetCustomDomainPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomDomainPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomDomainPolicyResponse): GetCustomDomainPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomDomainPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomDomainPolicyResponse;
  static deserializeBinaryFromReader(message: GetCustomDomainPolicyResponse, reader: jspb.BinaryReader): GetCustomDomainPolicyResponse;
}

export namespace GetCustomDomainPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.DomainPolicy.AsObject,
    isDefault: boolean,
  }
}

export class AddCustomDomainPolicyRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): AddCustomDomainPolicyRequest;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): AddCustomDomainPolicyRequest;

  getValidateOrgDomains(): boolean;
  setValidateOrgDomains(value: boolean): AddCustomDomainPolicyRequest;

  getSmtpSenderAddressMatchesInstanceDomain(): boolean;
  setSmtpSenderAddressMatchesInstanceDomain(value: boolean): AddCustomDomainPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomDomainPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomDomainPolicyRequest): AddCustomDomainPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomDomainPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomDomainPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomDomainPolicyRequest, reader: jspb.BinaryReader): AddCustomDomainPolicyRequest;
}

export namespace AddCustomDomainPolicyRequest {
  export type AsObject = {
    orgId: string,
    userLoginMustBeDomain: boolean,
    validateOrgDomains: boolean,
    smtpSenderAddressMatchesInstanceDomain: boolean,
  }
}

export class AddCustomDomainPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomDomainPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomDomainPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomDomainPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomDomainPolicyResponse): AddCustomDomainPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomDomainPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomDomainPolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomDomainPolicyResponse, reader: jspb.BinaryReader): AddCustomDomainPolicyResponse;
}

export namespace AddCustomDomainPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomDomainPolicyRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): UpdateCustomDomainPolicyRequest;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): UpdateCustomDomainPolicyRequest;

  getValidateOrgDomains(): boolean;
  setValidateOrgDomains(value: boolean): UpdateCustomDomainPolicyRequest;

  getSmtpSenderAddressMatchesInstanceDomain(): boolean;
  setSmtpSenderAddressMatchesInstanceDomain(value: boolean): UpdateCustomDomainPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomDomainPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomDomainPolicyRequest): UpdateCustomDomainPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomDomainPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomDomainPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomDomainPolicyRequest, reader: jspb.BinaryReader): UpdateCustomDomainPolicyRequest;
}

export namespace UpdateCustomDomainPolicyRequest {
  export type AsObject = {
    orgId: string,
    userLoginMustBeDomain: boolean,
    validateOrgDomains: boolean,
    smtpSenderAddressMatchesInstanceDomain: boolean,
  }
}

export class UpdateCustomDomainPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomDomainPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomDomainPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomDomainPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomDomainPolicyResponse): UpdateCustomDomainPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomDomainPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomDomainPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomDomainPolicyResponse, reader: jspb.BinaryReader): UpdateCustomDomainPolicyResponse;
}

export namespace UpdateCustomDomainPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomDomainPolicyToDefaultRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): ResetCustomDomainPolicyToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomDomainPolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomDomainPolicyToDefaultRequest): ResetCustomDomainPolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomDomainPolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomDomainPolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomDomainPolicyToDefaultRequest, reader: jspb.BinaryReader): ResetCustomDomainPolicyToDefaultRequest;
}

export namespace ResetCustomDomainPolicyToDefaultRequest {
  export type AsObject = {
    orgId: string,
  }
}

export class ResetCustomDomainPolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomDomainPolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomDomainPolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomDomainPolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomDomainPolicyToDefaultResponse): ResetCustomDomainPolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomDomainPolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomDomainPolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomDomainPolicyToDefaultResponse, reader: jspb.BinaryReader): ResetCustomDomainPolicyToDefaultResponse;
}

export namespace ResetCustomDomainPolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetLabelPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLabelPolicyRequest): GetLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLabelPolicyRequest;
  static deserializeBinaryFromReader(message: GetLabelPolicyRequest, reader: jspb.BinaryReader): GetLabelPolicyRequest;
}

export namespace GetLabelPolicyRequest {
  export type AsObject = {
  }
}

export class GetLabelPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LabelPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LabelPolicy): GetLabelPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLabelPolicyResponse): GetLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLabelPolicyResponse;
  static deserializeBinaryFromReader(message: GetLabelPolicyResponse, reader: jspb.BinaryReader): GetLabelPolicyResponse;
}

export namespace GetLabelPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LabelPolicy.AsObject,
  }
}

export class GetPreviewLabelPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPreviewLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPreviewLabelPolicyRequest): GetPreviewLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetPreviewLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPreviewLabelPolicyRequest;
  static deserializeBinaryFromReader(message: GetPreviewLabelPolicyRequest, reader: jspb.BinaryReader): GetPreviewLabelPolicyRequest;
}

export namespace GetPreviewLabelPolicyRequest {
  export type AsObject = {
  }
}

export class GetPreviewLabelPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LabelPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LabelPolicy): GetPreviewLabelPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetPreviewLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPreviewLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPreviewLabelPolicyResponse): GetPreviewLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetPreviewLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPreviewLabelPolicyResponse;
  static deserializeBinaryFromReader(message: GetPreviewLabelPolicyResponse, reader: jspb.BinaryReader): GetPreviewLabelPolicyResponse;
}

export namespace GetPreviewLabelPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LabelPolicy.AsObject,
  }
}

export class UpdateLabelPolicyRequest extends jspb.Message {
  getPrimaryColor(): string;
  setPrimaryColor(value: string): UpdateLabelPolicyRequest;

  getHideLoginNameSuffix(): boolean;
  setHideLoginNameSuffix(value: boolean): UpdateLabelPolicyRequest;

  getWarnColor(): string;
  setWarnColor(value: string): UpdateLabelPolicyRequest;

  getBackgroundColor(): string;
  setBackgroundColor(value: string): UpdateLabelPolicyRequest;

  getFontColor(): string;
  setFontColor(value: string): UpdateLabelPolicyRequest;

  getPrimaryColorDark(): string;
  setPrimaryColorDark(value: string): UpdateLabelPolicyRequest;

  getBackgroundColorDark(): string;
  setBackgroundColorDark(value: string): UpdateLabelPolicyRequest;

  getWarnColorDark(): string;
  setWarnColorDark(value: string): UpdateLabelPolicyRequest;

  getFontColorDark(): string;
  setFontColorDark(value: string): UpdateLabelPolicyRequest;

  getDisableWatermark(): boolean;
  setDisableWatermark(value: boolean): UpdateLabelPolicyRequest;

  getThemeMode(): zitadel_policy_pb.ThemeMode;
  setThemeMode(value: zitadel_policy_pb.ThemeMode): UpdateLabelPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLabelPolicyRequest): UpdateLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLabelPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateLabelPolicyRequest, reader: jspb.BinaryReader): UpdateLabelPolicyRequest;
}

export namespace UpdateLabelPolicyRequest {
  export type AsObject = {
    primaryColor: string,
    hideLoginNameSuffix: boolean,
    warnColor: string,
    backgroundColor: string,
    fontColor: string,
    primaryColorDark: string,
    backgroundColorDark: string,
    warnColorDark: string,
    fontColorDark: string,
    disableWatermark: boolean,
    themeMode: zitadel_policy_pb.ThemeMode,
  }
}

export class UpdateLabelPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateLabelPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLabelPolicyResponse): UpdateLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLabelPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateLabelPolicyResponse, reader: jspb.BinaryReader): UpdateLabelPolicyResponse;
}

export namespace UpdateLabelPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ActivateLabelPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateLabelPolicyRequest): ActivateLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: ActivateLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateLabelPolicyRequest;
  static deserializeBinaryFromReader(message: ActivateLabelPolicyRequest, reader: jspb.BinaryReader): ActivateLabelPolicyRequest;
}

export namespace ActivateLabelPolicyRequest {
  export type AsObject = {
  }
}

export class ActivateLabelPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ActivateLabelPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): ActivateLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateLabelPolicyResponse): ActivateLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: ActivateLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateLabelPolicyResponse;
  static deserializeBinaryFromReader(message: ActivateLabelPolicyResponse, reader: jspb.BinaryReader): ActivateLabelPolicyResponse;
}

export namespace ActivateLabelPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveLabelPolicyLogoRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyLogoRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyLogoRequest): RemoveLabelPolicyLogoRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyLogoRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyLogoRequest;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyLogoRequest, reader: jspb.BinaryReader): RemoveLabelPolicyLogoRequest;
}

export namespace RemoveLabelPolicyLogoRequest {
  export type AsObject = {
  }
}

export class RemoveLabelPolicyLogoResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveLabelPolicyLogoResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveLabelPolicyLogoResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyLogoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyLogoResponse): RemoveLabelPolicyLogoResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyLogoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyLogoResponse;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyLogoResponse, reader: jspb.BinaryReader): RemoveLabelPolicyLogoResponse;
}

export namespace RemoveLabelPolicyLogoResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveLabelPolicyLogoDarkRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyLogoDarkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyLogoDarkRequest): RemoveLabelPolicyLogoDarkRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyLogoDarkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyLogoDarkRequest;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyLogoDarkRequest, reader: jspb.BinaryReader): RemoveLabelPolicyLogoDarkRequest;
}

export namespace RemoveLabelPolicyLogoDarkRequest {
  export type AsObject = {
  }
}

export class RemoveLabelPolicyLogoDarkResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveLabelPolicyLogoDarkResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveLabelPolicyLogoDarkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyLogoDarkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyLogoDarkResponse): RemoveLabelPolicyLogoDarkResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyLogoDarkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyLogoDarkResponse;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyLogoDarkResponse, reader: jspb.BinaryReader): RemoveLabelPolicyLogoDarkResponse;
}

export namespace RemoveLabelPolicyLogoDarkResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveLabelPolicyIconRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyIconRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyIconRequest): RemoveLabelPolicyIconRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyIconRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyIconRequest;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyIconRequest, reader: jspb.BinaryReader): RemoveLabelPolicyIconRequest;
}

export namespace RemoveLabelPolicyIconRequest {
  export type AsObject = {
  }
}

export class RemoveLabelPolicyIconResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveLabelPolicyIconResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveLabelPolicyIconResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyIconResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyIconResponse): RemoveLabelPolicyIconResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyIconResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyIconResponse;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyIconResponse, reader: jspb.BinaryReader): RemoveLabelPolicyIconResponse;
}

export namespace RemoveLabelPolicyIconResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveLabelPolicyIconDarkRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyIconDarkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyIconDarkRequest): RemoveLabelPolicyIconDarkRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyIconDarkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyIconDarkRequest;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyIconDarkRequest, reader: jspb.BinaryReader): RemoveLabelPolicyIconDarkRequest;
}

export namespace RemoveLabelPolicyIconDarkRequest {
  export type AsObject = {
  }
}

export class RemoveLabelPolicyIconDarkResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveLabelPolicyIconDarkResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveLabelPolicyIconDarkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyIconDarkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyIconDarkResponse): RemoveLabelPolicyIconDarkResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyIconDarkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyIconDarkResponse;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyIconDarkResponse, reader: jspb.BinaryReader): RemoveLabelPolicyIconDarkResponse;
}

export namespace RemoveLabelPolicyIconDarkResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveLabelPolicyFontRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyFontRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyFontRequest): RemoveLabelPolicyFontRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyFontRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyFontRequest;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyFontRequest, reader: jspb.BinaryReader): RemoveLabelPolicyFontRequest;
}

export namespace RemoveLabelPolicyFontRequest {
  export type AsObject = {
  }
}

export class RemoveLabelPolicyFontResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveLabelPolicyFontResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveLabelPolicyFontResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveLabelPolicyFontResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveLabelPolicyFontResponse): RemoveLabelPolicyFontResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveLabelPolicyFontResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveLabelPolicyFontResponse;
  static deserializeBinaryFromReader(message: RemoveLabelPolicyFontResponse, reader: jspb.BinaryReader): RemoveLabelPolicyFontResponse;
}

export namespace RemoveLabelPolicyFontResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetLoginPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLoginPolicyRequest): GetLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLoginPolicyRequest;
  static deserializeBinaryFromReader(message: GetLoginPolicyRequest, reader: jspb.BinaryReader): GetLoginPolicyRequest;
}

export namespace GetLoginPolicyRequest {
  export type AsObject = {
  }
}

export class GetLoginPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LoginPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LoginPolicy): GetLoginPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLoginPolicyResponse): GetLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLoginPolicyResponse;
  static deserializeBinaryFromReader(message: GetLoginPolicyResponse, reader: jspb.BinaryReader): GetLoginPolicyResponse;
}

export namespace GetLoginPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LoginPolicy.AsObject,
  }
}

export class UpdateLoginPolicyRequest extends jspb.Message {
  getAllowUsernamePassword(): boolean;
  setAllowUsernamePassword(value: boolean): UpdateLoginPolicyRequest;

  getAllowRegister(): boolean;
  setAllowRegister(value: boolean): UpdateLoginPolicyRequest;

  getAllowExternalIdp(): boolean;
  setAllowExternalIdp(value: boolean): UpdateLoginPolicyRequest;

  getForceMfa(): boolean;
  setForceMfa(value: boolean): UpdateLoginPolicyRequest;

  getPasswordlessType(): zitadel_policy_pb.PasswordlessType;
  setPasswordlessType(value: zitadel_policy_pb.PasswordlessType): UpdateLoginPolicyRequest;

  getHidePasswordReset(): boolean;
  setHidePasswordReset(value: boolean): UpdateLoginPolicyRequest;

  getIgnoreUnknownUsernames(): boolean;
  setIgnoreUnknownUsernames(value: boolean): UpdateLoginPolicyRequest;

  getDefaultRedirectUri(): string;
  setDefaultRedirectUri(value: string): UpdateLoginPolicyRequest;

  getPasswordCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setPasswordCheckLifetime(value?: google_protobuf_duration_pb.Duration): UpdateLoginPolicyRequest;
  hasPasswordCheckLifetime(): boolean;
  clearPasswordCheckLifetime(): UpdateLoginPolicyRequest;

  getExternalLoginCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setExternalLoginCheckLifetime(value?: google_protobuf_duration_pb.Duration): UpdateLoginPolicyRequest;
  hasExternalLoginCheckLifetime(): boolean;
  clearExternalLoginCheckLifetime(): UpdateLoginPolicyRequest;

  getMfaInitSkipLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMfaInitSkipLifetime(value?: google_protobuf_duration_pb.Duration): UpdateLoginPolicyRequest;
  hasMfaInitSkipLifetime(): boolean;
  clearMfaInitSkipLifetime(): UpdateLoginPolicyRequest;

  getSecondFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setSecondFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): UpdateLoginPolicyRequest;
  hasSecondFactorCheckLifetime(): boolean;
  clearSecondFactorCheckLifetime(): UpdateLoginPolicyRequest;

  getMultiFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMultiFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): UpdateLoginPolicyRequest;
  hasMultiFactorCheckLifetime(): boolean;
  clearMultiFactorCheckLifetime(): UpdateLoginPolicyRequest;

  getAllowDomainDiscovery(): boolean;
  setAllowDomainDiscovery(value: boolean): UpdateLoginPolicyRequest;

  getDisableLoginWithEmail(): boolean;
  setDisableLoginWithEmail(value: boolean): UpdateLoginPolicyRequest;

  getDisableLoginWithPhone(): boolean;
  setDisableLoginWithPhone(value: boolean): UpdateLoginPolicyRequest;

  getForceMfaLocalOnly(): boolean;
  setForceMfaLocalOnly(value: boolean): UpdateLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLoginPolicyRequest): UpdateLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLoginPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateLoginPolicyRequest, reader: jspb.BinaryReader): UpdateLoginPolicyRequest;
}

export namespace UpdateLoginPolicyRequest {
  export type AsObject = {
    allowUsernamePassword: boolean,
    allowRegister: boolean,
    allowExternalIdp: boolean,
    forceMfa: boolean,
    passwordlessType: zitadel_policy_pb.PasswordlessType,
    hidePasswordReset: boolean,
    ignoreUnknownUsernames: boolean,
    defaultRedirectUri: string,
    passwordCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    externalLoginCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    mfaInitSkipLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    secondFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    multiFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    allowDomainDiscovery: boolean,
    disableLoginWithEmail: boolean,
    disableLoginWithPhone: boolean,
    forceMfaLocalOnly: boolean,
  }
}

export class UpdateLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLoginPolicyResponse): UpdateLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLoginPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateLoginPolicyResponse, reader: jspb.BinaryReader): UpdateLoginPolicyResponse;
}

export namespace UpdateLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListLoginPolicyIDPsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListLoginPolicyIDPsRequest;
  hasQuery(): boolean;
  clearQuery(): ListLoginPolicyIDPsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicyIDPsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicyIDPsRequest): ListLoginPolicyIDPsRequest.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicyIDPsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicyIDPsRequest;
  static deserializeBinaryFromReader(message: ListLoginPolicyIDPsRequest, reader: jspb.BinaryReader): ListLoginPolicyIDPsRequest;
}

export namespace ListLoginPolicyIDPsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListLoginPolicyIDPsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListLoginPolicyIDPsResponse;
  hasDetails(): boolean;
  clearDetails(): ListLoginPolicyIDPsResponse;

  getResultList(): Array<zitadel_idp_pb.IDPLoginPolicyLink>;
  setResultList(value: Array<zitadel_idp_pb.IDPLoginPolicyLink>): ListLoginPolicyIDPsResponse;
  clearResultList(): ListLoginPolicyIDPsResponse;
  addResult(value?: zitadel_idp_pb.IDPLoginPolicyLink, index?: number): zitadel_idp_pb.IDPLoginPolicyLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicyIDPsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicyIDPsResponse): ListLoginPolicyIDPsResponse.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicyIDPsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicyIDPsResponse;
  static deserializeBinaryFromReader(message: ListLoginPolicyIDPsResponse, reader: jspb.BinaryReader): ListLoginPolicyIDPsResponse;
}

export namespace ListLoginPolicyIDPsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_idp_pb.IDPLoginPolicyLink.AsObject>,
  }
}

export class AddIDPToLoginPolicyRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): AddIDPToLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIDPToLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddIDPToLoginPolicyRequest): AddIDPToLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddIDPToLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIDPToLoginPolicyRequest;
  static deserializeBinaryFromReader(message: AddIDPToLoginPolicyRequest, reader: jspb.BinaryReader): AddIDPToLoginPolicyRequest;
}

export namespace AddIDPToLoginPolicyRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class AddIDPToLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddIDPToLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddIDPToLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIDPToLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddIDPToLoginPolicyResponse): AddIDPToLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddIDPToLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIDPToLoginPolicyResponse;
  static deserializeBinaryFromReader(message: AddIDPToLoginPolicyResponse, reader: jspb.BinaryReader): AddIDPToLoginPolicyResponse;
}

export namespace AddIDPToLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveIDPFromLoginPolicyRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): RemoveIDPFromLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIDPFromLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIDPFromLoginPolicyRequest): RemoveIDPFromLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveIDPFromLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIDPFromLoginPolicyRequest;
  static deserializeBinaryFromReader(message: RemoveIDPFromLoginPolicyRequest, reader: jspb.BinaryReader): RemoveIDPFromLoginPolicyRequest;
}

export namespace RemoveIDPFromLoginPolicyRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class RemoveIDPFromLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveIDPFromLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveIDPFromLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIDPFromLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIDPFromLoginPolicyResponse): RemoveIDPFromLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveIDPFromLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIDPFromLoginPolicyResponse;
  static deserializeBinaryFromReader(message: RemoveIDPFromLoginPolicyResponse, reader: jspb.BinaryReader): RemoveIDPFromLoginPolicyResponse;
}

export namespace RemoveIDPFromLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListLoginPolicySecondFactorsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicySecondFactorsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicySecondFactorsRequest): ListLoginPolicySecondFactorsRequest.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicySecondFactorsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicySecondFactorsRequest;
  static deserializeBinaryFromReader(message: ListLoginPolicySecondFactorsRequest, reader: jspb.BinaryReader): ListLoginPolicySecondFactorsRequest;
}

export namespace ListLoginPolicySecondFactorsRequest {
  export type AsObject = {
  }
}

export class ListLoginPolicySecondFactorsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListLoginPolicySecondFactorsResponse;
  hasDetails(): boolean;
  clearDetails(): ListLoginPolicySecondFactorsResponse;

  getResultList(): Array<zitadel_policy_pb.SecondFactorType>;
  setResultList(value: Array<zitadel_policy_pb.SecondFactorType>): ListLoginPolicySecondFactorsResponse;
  clearResultList(): ListLoginPolicySecondFactorsResponse;
  addResult(value: zitadel_policy_pb.SecondFactorType, index?: number): ListLoginPolicySecondFactorsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicySecondFactorsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicySecondFactorsResponse): ListLoginPolicySecondFactorsResponse.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicySecondFactorsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicySecondFactorsResponse;
  static deserializeBinaryFromReader(message: ListLoginPolicySecondFactorsResponse, reader: jspb.BinaryReader): ListLoginPolicySecondFactorsResponse;
}

export namespace ListLoginPolicySecondFactorsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_policy_pb.SecondFactorType>,
  }
}

export class AddSecondFactorToLoginPolicyRequest extends jspb.Message {
  getType(): zitadel_policy_pb.SecondFactorType;
  setType(value: zitadel_policy_pb.SecondFactorType): AddSecondFactorToLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSecondFactorToLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddSecondFactorToLoginPolicyRequest): AddSecondFactorToLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddSecondFactorToLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSecondFactorToLoginPolicyRequest;
  static deserializeBinaryFromReader(message: AddSecondFactorToLoginPolicyRequest, reader: jspb.BinaryReader): AddSecondFactorToLoginPolicyRequest;
}

export namespace AddSecondFactorToLoginPolicyRequest {
  export type AsObject = {
    type: zitadel_policy_pb.SecondFactorType,
  }
}

export class AddSecondFactorToLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddSecondFactorToLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddSecondFactorToLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSecondFactorToLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSecondFactorToLoginPolicyResponse): AddSecondFactorToLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddSecondFactorToLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSecondFactorToLoginPolicyResponse;
  static deserializeBinaryFromReader(message: AddSecondFactorToLoginPolicyResponse, reader: jspb.BinaryReader): AddSecondFactorToLoginPolicyResponse;
}

export namespace AddSecondFactorToLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveSecondFactorFromLoginPolicyRequest extends jspb.Message {
  getType(): zitadel_policy_pb.SecondFactorType;
  setType(value: zitadel_policy_pb.SecondFactorType): RemoveSecondFactorFromLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveSecondFactorFromLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveSecondFactorFromLoginPolicyRequest): RemoveSecondFactorFromLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveSecondFactorFromLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveSecondFactorFromLoginPolicyRequest;
  static deserializeBinaryFromReader(message: RemoveSecondFactorFromLoginPolicyRequest, reader: jspb.BinaryReader): RemoveSecondFactorFromLoginPolicyRequest;
}

export namespace RemoveSecondFactorFromLoginPolicyRequest {
  export type AsObject = {
    type: zitadel_policy_pb.SecondFactorType,
  }
}

export class RemoveSecondFactorFromLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveSecondFactorFromLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveSecondFactorFromLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveSecondFactorFromLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveSecondFactorFromLoginPolicyResponse): RemoveSecondFactorFromLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveSecondFactorFromLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveSecondFactorFromLoginPolicyResponse;
  static deserializeBinaryFromReader(message: RemoveSecondFactorFromLoginPolicyResponse, reader: jspb.BinaryReader): RemoveSecondFactorFromLoginPolicyResponse;
}

export namespace RemoveSecondFactorFromLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListLoginPolicyMultiFactorsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicyMultiFactorsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicyMultiFactorsRequest): ListLoginPolicyMultiFactorsRequest.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicyMultiFactorsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicyMultiFactorsRequest;
  static deserializeBinaryFromReader(message: ListLoginPolicyMultiFactorsRequest, reader: jspb.BinaryReader): ListLoginPolicyMultiFactorsRequest;
}

export namespace ListLoginPolicyMultiFactorsRequest {
  export type AsObject = {
  }
}

export class ListLoginPolicyMultiFactorsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListLoginPolicyMultiFactorsResponse;
  hasDetails(): boolean;
  clearDetails(): ListLoginPolicyMultiFactorsResponse;

  getResultList(): Array<zitadel_policy_pb.MultiFactorType>;
  setResultList(value: Array<zitadel_policy_pb.MultiFactorType>): ListLoginPolicyMultiFactorsResponse;
  clearResultList(): ListLoginPolicyMultiFactorsResponse;
  addResult(value: zitadel_policy_pb.MultiFactorType, index?: number): ListLoginPolicyMultiFactorsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicyMultiFactorsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicyMultiFactorsResponse): ListLoginPolicyMultiFactorsResponse.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicyMultiFactorsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicyMultiFactorsResponse;
  static deserializeBinaryFromReader(message: ListLoginPolicyMultiFactorsResponse, reader: jspb.BinaryReader): ListLoginPolicyMultiFactorsResponse;
}

export namespace ListLoginPolicyMultiFactorsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_policy_pb.MultiFactorType>,
  }
}

export class AddMultiFactorToLoginPolicyRequest extends jspb.Message {
  getType(): zitadel_policy_pb.MultiFactorType;
  setType(value: zitadel_policy_pb.MultiFactorType): AddMultiFactorToLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMultiFactorToLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMultiFactorToLoginPolicyRequest): AddMultiFactorToLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddMultiFactorToLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMultiFactorToLoginPolicyRequest;
  static deserializeBinaryFromReader(message: AddMultiFactorToLoginPolicyRequest, reader: jspb.BinaryReader): AddMultiFactorToLoginPolicyRequest;
}

export namespace AddMultiFactorToLoginPolicyRequest {
  export type AsObject = {
    type: zitadel_policy_pb.MultiFactorType,
  }
}

export class AddMultiFactorToLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMultiFactorToLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddMultiFactorToLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMultiFactorToLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMultiFactorToLoginPolicyResponse): AddMultiFactorToLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddMultiFactorToLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMultiFactorToLoginPolicyResponse;
  static deserializeBinaryFromReader(message: AddMultiFactorToLoginPolicyResponse, reader: jspb.BinaryReader): AddMultiFactorToLoginPolicyResponse;
}

export namespace AddMultiFactorToLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMultiFactorFromLoginPolicyRequest extends jspb.Message {
  getType(): zitadel_policy_pb.MultiFactorType;
  setType(value: zitadel_policy_pb.MultiFactorType): RemoveMultiFactorFromLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMultiFactorFromLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMultiFactorFromLoginPolicyRequest): RemoveMultiFactorFromLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMultiFactorFromLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMultiFactorFromLoginPolicyRequest;
  static deserializeBinaryFromReader(message: RemoveMultiFactorFromLoginPolicyRequest, reader: jspb.BinaryReader): RemoveMultiFactorFromLoginPolicyRequest;
}

export namespace RemoveMultiFactorFromLoginPolicyRequest {
  export type AsObject = {
    type: zitadel_policy_pb.MultiFactorType,
  }
}

export class RemoveMultiFactorFromLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMultiFactorFromLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMultiFactorFromLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMultiFactorFromLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMultiFactorFromLoginPolicyResponse): RemoveMultiFactorFromLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMultiFactorFromLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMultiFactorFromLoginPolicyResponse;
  static deserializeBinaryFromReader(message: RemoveMultiFactorFromLoginPolicyResponse, reader: jspb.BinaryReader): RemoveMultiFactorFromLoginPolicyResponse;
}

export namespace RemoveMultiFactorFromLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetPasswordComplexityPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPasswordComplexityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPasswordComplexityPolicyRequest): GetPasswordComplexityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetPasswordComplexityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPasswordComplexityPolicyRequest;
  static deserializeBinaryFromReader(message: GetPasswordComplexityPolicyRequest, reader: jspb.BinaryReader): GetPasswordComplexityPolicyRequest;
}

export namespace GetPasswordComplexityPolicyRequest {
  export type AsObject = {
  }
}

export class GetPasswordComplexityPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PasswordComplexityPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PasswordComplexityPolicy): GetPasswordComplexityPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetPasswordComplexityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPasswordComplexityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPasswordComplexityPolicyResponse): GetPasswordComplexityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetPasswordComplexityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPasswordComplexityPolicyResponse;
  static deserializeBinaryFromReader(message: GetPasswordComplexityPolicyResponse, reader: jspb.BinaryReader): GetPasswordComplexityPolicyResponse;
}

export namespace GetPasswordComplexityPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PasswordComplexityPolicy.AsObject,
  }
}

export class UpdatePasswordComplexityPolicyRequest extends jspb.Message {
  getMinLength(): number;
  setMinLength(value: number): UpdatePasswordComplexityPolicyRequest;

  getHasUppercase(): boolean;
  setHasUppercase(value: boolean): UpdatePasswordComplexityPolicyRequest;

  getHasLowercase(): boolean;
  setHasLowercase(value: boolean): UpdatePasswordComplexityPolicyRequest;

  getHasNumber(): boolean;
  setHasNumber(value: boolean): UpdatePasswordComplexityPolicyRequest;

  getHasSymbol(): boolean;
  setHasSymbol(value: boolean): UpdatePasswordComplexityPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePasswordComplexityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePasswordComplexityPolicyRequest): UpdatePasswordComplexityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdatePasswordComplexityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePasswordComplexityPolicyRequest;
  static deserializeBinaryFromReader(message: UpdatePasswordComplexityPolicyRequest, reader: jspb.BinaryReader): UpdatePasswordComplexityPolicyRequest;
}

export namespace UpdatePasswordComplexityPolicyRequest {
  export type AsObject = {
    minLength: number,
    hasUppercase: boolean,
    hasLowercase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
  }
}

export class UpdatePasswordComplexityPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdatePasswordComplexityPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdatePasswordComplexityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePasswordComplexityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePasswordComplexityPolicyResponse): UpdatePasswordComplexityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdatePasswordComplexityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePasswordComplexityPolicyResponse;
  static deserializeBinaryFromReader(message: UpdatePasswordComplexityPolicyResponse, reader: jspb.BinaryReader): UpdatePasswordComplexityPolicyResponse;
}

export namespace UpdatePasswordComplexityPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetPasswordAgePolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPasswordAgePolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPasswordAgePolicyRequest): GetPasswordAgePolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetPasswordAgePolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPasswordAgePolicyRequest;
  static deserializeBinaryFromReader(message: GetPasswordAgePolicyRequest, reader: jspb.BinaryReader): GetPasswordAgePolicyRequest;
}

export namespace GetPasswordAgePolicyRequest {
  export type AsObject = {
  }
}

export class GetPasswordAgePolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PasswordAgePolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PasswordAgePolicy): GetPasswordAgePolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetPasswordAgePolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPasswordAgePolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPasswordAgePolicyResponse): GetPasswordAgePolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetPasswordAgePolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPasswordAgePolicyResponse;
  static deserializeBinaryFromReader(message: GetPasswordAgePolicyResponse, reader: jspb.BinaryReader): GetPasswordAgePolicyResponse;
}

export namespace GetPasswordAgePolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PasswordAgePolicy.AsObject,
  }
}

export class UpdatePasswordAgePolicyRequest extends jspb.Message {
  getMaxAgeDays(): number;
  setMaxAgeDays(value: number): UpdatePasswordAgePolicyRequest;

  getExpireWarnDays(): number;
  setExpireWarnDays(value: number): UpdatePasswordAgePolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePasswordAgePolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePasswordAgePolicyRequest): UpdatePasswordAgePolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdatePasswordAgePolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePasswordAgePolicyRequest;
  static deserializeBinaryFromReader(message: UpdatePasswordAgePolicyRequest, reader: jspb.BinaryReader): UpdatePasswordAgePolicyRequest;
}

export namespace UpdatePasswordAgePolicyRequest {
  export type AsObject = {
    maxAgeDays: number,
    expireWarnDays: number,
  }
}

export class UpdatePasswordAgePolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdatePasswordAgePolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdatePasswordAgePolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePasswordAgePolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePasswordAgePolicyResponse): UpdatePasswordAgePolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdatePasswordAgePolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePasswordAgePolicyResponse;
  static deserializeBinaryFromReader(message: UpdatePasswordAgePolicyResponse, reader: jspb.BinaryReader): UpdatePasswordAgePolicyResponse;
}

export namespace UpdatePasswordAgePolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetLockoutPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLockoutPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLockoutPolicyRequest): GetLockoutPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetLockoutPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLockoutPolicyRequest;
  static deserializeBinaryFromReader(message: GetLockoutPolicyRequest, reader: jspb.BinaryReader): GetLockoutPolicyRequest;
}

export namespace GetLockoutPolicyRequest {
  export type AsObject = {
  }
}

export class GetLockoutPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LockoutPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LockoutPolicy): GetLockoutPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetLockoutPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLockoutPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLockoutPolicyResponse): GetLockoutPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetLockoutPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLockoutPolicyResponse;
  static deserializeBinaryFromReader(message: GetLockoutPolicyResponse, reader: jspb.BinaryReader): GetLockoutPolicyResponse;
}

export namespace GetLockoutPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LockoutPolicy.AsObject,
  }
}

export class UpdateLockoutPolicyRequest extends jspb.Message {
  getMaxPasswordAttempts(): number;
  setMaxPasswordAttempts(value: number): UpdateLockoutPolicyRequest;

  getMaxOtpAttempts(): number;
  setMaxOtpAttempts(value: number): UpdateLockoutPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLockoutPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLockoutPolicyRequest): UpdateLockoutPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateLockoutPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLockoutPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateLockoutPolicyRequest, reader: jspb.BinaryReader): UpdateLockoutPolicyRequest;
}

export namespace UpdateLockoutPolicyRequest {
  export type AsObject = {
    maxPasswordAttempts: number,
    maxOtpAttempts: number,
  }
}

export class UpdateLockoutPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateLockoutPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateLockoutPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLockoutPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLockoutPolicyResponse): UpdateLockoutPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateLockoutPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLockoutPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateLockoutPolicyResponse, reader: jspb.BinaryReader): UpdateLockoutPolicyResponse;
}

export namespace UpdateLockoutPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetPrivacyPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPrivacyPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPrivacyPolicyRequest): GetPrivacyPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetPrivacyPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPrivacyPolicyRequest;
  static deserializeBinaryFromReader(message: GetPrivacyPolicyRequest, reader: jspb.BinaryReader): GetPrivacyPolicyRequest;
}

export namespace GetPrivacyPolicyRequest {
  export type AsObject = {
  }
}

export class GetPrivacyPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PrivacyPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PrivacyPolicy): GetPrivacyPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetPrivacyPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPrivacyPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPrivacyPolicyResponse): GetPrivacyPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetPrivacyPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPrivacyPolicyResponse;
  static deserializeBinaryFromReader(message: GetPrivacyPolicyResponse, reader: jspb.BinaryReader): GetPrivacyPolicyResponse;
}

export namespace GetPrivacyPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PrivacyPolicy.AsObject,
  }
}

export class UpdatePrivacyPolicyRequest extends jspb.Message {
  getTosLink(): string;
  setTosLink(value: string): UpdatePrivacyPolicyRequest;

  getPrivacyLink(): string;
  setPrivacyLink(value: string): UpdatePrivacyPolicyRequest;

  getHelpLink(): string;
  setHelpLink(value: string): UpdatePrivacyPolicyRequest;

  getSupportEmail(): string;
  setSupportEmail(value: string): UpdatePrivacyPolicyRequest;

  getDocsLink(): string;
  setDocsLink(value: string): UpdatePrivacyPolicyRequest;

  getCustomLink(): string;
  setCustomLink(value: string): UpdatePrivacyPolicyRequest;

  getCustomLinkText(): string;
  setCustomLinkText(value: string): UpdatePrivacyPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePrivacyPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePrivacyPolicyRequest): UpdatePrivacyPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdatePrivacyPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePrivacyPolicyRequest;
  static deserializeBinaryFromReader(message: UpdatePrivacyPolicyRequest, reader: jspb.BinaryReader): UpdatePrivacyPolicyRequest;
}

export namespace UpdatePrivacyPolicyRequest {
  export type AsObject = {
    tosLink: string,
    privacyLink: string,
    helpLink: string,
    supportEmail: string,
    docsLink: string,
    customLink: string,
    customLinkText: string,
  }
}

export class UpdatePrivacyPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdatePrivacyPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdatePrivacyPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePrivacyPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePrivacyPolicyResponse): UpdatePrivacyPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdatePrivacyPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePrivacyPolicyResponse;
  static deserializeBinaryFromReader(message: UpdatePrivacyPolicyResponse, reader: jspb.BinaryReader): UpdatePrivacyPolicyResponse;
}

export namespace UpdatePrivacyPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddNotificationPolicyRequest extends jspb.Message {
  getPasswordChange(): boolean;
  setPasswordChange(value: boolean): AddNotificationPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddNotificationPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddNotificationPolicyRequest): AddNotificationPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddNotificationPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddNotificationPolicyRequest;
  static deserializeBinaryFromReader(message: AddNotificationPolicyRequest, reader: jspb.BinaryReader): AddNotificationPolicyRequest;
}

export namespace AddNotificationPolicyRequest {
  export type AsObject = {
    passwordChange: boolean,
  }
}

export class AddNotificationPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddNotificationPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddNotificationPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddNotificationPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddNotificationPolicyResponse): AddNotificationPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddNotificationPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddNotificationPolicyResponse;
  static deserializeBinaryFromReader(message: AddNotificationPolicyResponse, reader: jspb.BinaryReader): AddNotificationPolicyResponse;
}

export namespace AddNotificationPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetNotificationPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNotificationPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetNotificationPolicyRequest): GetNotificationPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetNotificationPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNotificationPolicyRequest;
  static deserializeBinaryFromReader(message: GetNotificationPolicyRequest, reader: jspb.BinaryReader): GetNotificationPolicyRequest;
}

export namespace GetNotificationPolicyRequest {
  export type AsObject = {
  }
}

export class GetNotificationPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.NotificationPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.NotificationPolicy): GetNotificationPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetNotificationPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNotificationPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetNotificationPolicyResponse): GetNotificationPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetNotificationPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNotificationPolicyResponse;
  static deserializeBinaryFromReader(message: GetNotificationPolicyResponse, reader: jspb.BinaryReader): GetNotificationPolicyResponse;
}

export namespace GetNotificationPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.NotificationPolicy.AsObject,
  }
}

export class UpdateNotificationPolicyRequest extends jspb.Message {
  getPasswordChange(): boolean;
  setPasswordChange(value: boolean): UpdateNotificationPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateNotificationPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateNotificationPolicyRequest): UpdateNotificationPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateNotificationPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateNotificationPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateNotificationPolicyRequest, reader: jspb.BinaryReader): UpdateNotificationPolicyRequest;
}

export namespace UpdateNotificationPolicyRequest {
  export type AsObject = {
    passwordChange: boolean,
  }
}

export class UpdateNotificationPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateNotificationPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateNotificationPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateNotificationPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateNotificationPolicyResponse): UpdateNotificationPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateNotificationPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateNotificationPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateNotificationPolicyResponse, reader: jspb.BinaryReader): UpdateNotificationPolicyResponse;
}

export namespace UpdateNotificationPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultInitMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultInitMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultInitMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultInitMessageTextRequest): GetDefaultInitMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultInitMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultInitMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultInitMessageTextRequest, reader: jspb.BinaryReader): GetDefaultInitMessageTextRequest;
}

export namespace GetDefaultInitMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultInitMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultInitMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultInitMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultInitMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultInitMessageTextResponse): GetDefaultInitMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultInitMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultInitMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultInitMessageTextResponse, reader: jspb.BinaryReader): GetDefaultInitMessageTextResponse;
}

export namespace GetDefaultInitMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetCustomInitMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomInitMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomInitMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomInitMessageTextRequest): GetCustomInitMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomInitMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomInitMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomInitMessageTextRequest, reader: jspb.BinaryReader): GetCustomInitMessageTextRequest;
}

export namespace GetCustomInitMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomInitMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomInitMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomInitMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomInitMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomInitMessageTextResponse): GetCustomInitMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomInitMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomInitMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomInitMessageTextResponse, reader: jspb.BinaryReader): GetCustomInitMessageTextResponse;
}

export namespace GetCustomInitMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultInitMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultInitMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetDefaultInitMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetDefaultInitMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetDefaultInitMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetDefaultInitMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultInitMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetDefaultInitMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetDefaultInitMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultInitMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultInitMessageTextRequest): SetDefaultInitMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultInitMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultInitMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultInitMessageTextRequest, reader: jspb.BinaryReader): SetDefaultInitMessageTextRequest;
}

export namespace SetDefaultInitMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetDefaultInitMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultInitMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultInitMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultInitMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultInitMessageTextResponse): SetDefaultInitMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultInitMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultInitMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultInitMessageTextResponse, reader: jspb.BinaryReader): SetDefaultInitMessageTextResponse;
}

export namespace SetDefaultInitMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomInitMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomInitMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomInitMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomInitMessageTextToDefaultRequest): ResetCustomInitMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomInitMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomInitMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomInitMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomInitMessageTextToDefaultRequest;
}

export namespace ResetCustomInitMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomInitMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomInitMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomInitMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomInitMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomInitMessageTextToDefaultResponse): ResetCustomInitMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomInitMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomInitMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomInitMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomInitMessageTextToDefaultResponse;
}

export namespace ResetCustomInitMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultPasswordResetMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultPasswordResetMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordResetMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordResetMessageTextRequest): GetDefaultPasswordResetMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordResetMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordResetMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultPasswordResetMessageTextRequest, reader: jspb.BinaryReader): GetDefaultPasswordResetMessageTextRequest;
}

export namespace GetDefaultPasswordResetMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultPasswordResetMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultPasswordResetMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultPasswordResetMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordResetMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordResetMessageTextResponse): GetDefaultPasswordResetMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordResetMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordResetMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultPasswordResetMessageTextResponse, reader: jspb.BinaryReader): GetDefaultPasswordResetMessageTextResponse;
}

export namespace GetDefaultPasswordResetMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetCustomPasswordResetMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomPasswordResetMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordResetMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordResetMessageTextRequest): GetCustomPasswordResetMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordResetMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordResetMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomPasswordResetMessageTextRequest, reader: jspb.BinaryReader): GetCustomPasswordResetMessageTextRequest;
}

export namespace GetCustomPasswordResetMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomPasswordResetMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomPasswordResetMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomPasswordResetMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordResetMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordResetMessageTextResponse): GetCustomPasswordResetMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordResetMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordResetMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomPasswordResetMessageTextResponse, reader: jspb.BinaryReader): GetCustomPasswordResetMessageTextResponse;
}

export namespace GetCustomPasswordResetMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultPasswordResetMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultPasswordResetMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetDefaultPasswordResetMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetDefaultPasswordResetMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetDefaultPasswordResetMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetDefaultPasswordResetMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultPasswordResetMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetDefaultPasswordResetMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetDefaultPasswordResetMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultPasswordResetMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultPasswordResetMessageTextRequest): SetDefaultPasswordResetMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultPasswordResetMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultPasswordResetMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultPasswordResetMessageTextRequest, reader: jspb.BinaryReader): SetDefaultPasswordResetMessageTextRequest;
}

export namespace SetDefaultPasswordResetMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetDefaultPasswordResetMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultPasswordResetMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultPasswordResetMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultPasswordResetMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultPasswordResetMessageTextResponse): SetDefaultPasswordResetMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultPasswordResetMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultPasswordResetMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultPasswordResetMessageTextResponse, reader: jspb.BinaryReader): SetDefaultPasswordResetMessageTextResponse;
}

export namespace SetDefaultPasswordResetMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomPasswordResetMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomPasswordResetMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordResetMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordResetMessageTextToDefaultRequest): ResetCustomPasswordResetMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordResetMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordResetMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomPasswordResetMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomPasswordResetMessageTextToDefaultRequest;
}

export namespace ResetCustomPasswordResetMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomPasswordResetMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomPasswordResetMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomPasswordResetMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordResetMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordResetMessageTextToDefaultResponse): ResetCustomPasswordResetMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordResetMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordResetMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomPasswordResetMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomPasswordResetMessageTextToDefaultResponse;
}

export namespace ResetCustomPasswordResetMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultVerifyEmailMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultVerifyEmailMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyEmailMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyEmailMessageTextRequest): GetDefaultVerifyEmailMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyEmailMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyEmailMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultVerifyEmailMessageTextRequest, reader: jspb.BinaryReader): GetDefaultVerifyEmailMessageTextRequest;
}

export namespace GetDefaultVerifyEmailMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultVerifyEmailMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultVerifyEmailMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultVerifyEmailMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyEmailMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyEmailMessageTextResponse): GetDefaultVerifyEmailMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyEmailMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyEmailMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultVerifyEmailMessageTextResponse, reader: jspb.BinaryReader): GetDefaultVerifyEmailMessageTextResponse;
}

export namespace GetDefaultVerifyEmailMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetCustomVerifyEmailMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomVerifyEmailMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyEmailMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyEmailMessageTextRequest): GetCustomVerifyEmailMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyEmailMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyEmailMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomVerifyEmailMessageTextRequest, reader: jspb.BinaryReader): GetCustomVerifyEmailMessageTextRequest;
}

export namespace GetCustomVerifyEmailMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomVerifyEmailMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomVerifyEmailMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomVerifyEmailMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyEmailMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyEmailMessageTextResponse): GetCustomVerifyEmailMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyEmailMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyEmailMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomVerifyEmailMessageTextResponse, reader: jspb.BinaryReader): GetCustomVerifyEmailMessageTextResponse;
}

export namespace GetCustomVerifyEmailMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultVerifyEmailMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultVerifyEmailMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetDefaultVerifyEmailMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetDefaultVerifyEmailMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetDefaultVerifyEmailMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetDefaultVerifyEmailMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultVerifyEmailMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetDefaultVerifyEmailMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetDefaultVerifyEmailMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultVerifyEmailMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultVerifyEmailMessageTextRequest): SetDefaultVerifyEmailMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultVerifyEmailMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultVerifyEmailMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultVerifyEmailMessageTextRequest, reader: jspb.BinaryReader): SetDefaultVerifyEmailMessageTextRequest;
}

export namespace SetDefaultVerifyEmailMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetDefaultVerifyEmailMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultVerifyEmailMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultVerifyEmailMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultVerifyEmailMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultVerifyEmailMessageTextResponse): SetDefaultVerifyEmailMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultVerifyEmailMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultVerifyEmailMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultVerifyEmailMessageTextResponse, reader: jspb.BinaryReader): SetDefaultVerifyEmailMessageTextResponse;
}

export namespace SetDefaultVerifyEmailMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomVerifyEmailMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomVerifyEmailMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyEmailMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyEmailMessageTextToDefaultRequest): ResetCustomVerifyEmailMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyEmailMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyEmailMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomVerifyEmailMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomVerifyEmailMessageTextToDefaultRequest;
}

export namespace ResetCustomVerifyEmailMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomVerifyEmailMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomVerifyEmailMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomVerifyEmailMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyEmailMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyEmailMessageTextToDefaultResponse): ResetCustomVerifyEmailMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyEmailMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyEmailMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomVerifyEmailMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomVerifyEmailMessageTextToDefaultResponse;
}

export namespace ResetCustomVerifyEmailMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultVerifyPhoneMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultVerifyPhoneMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyPhoneMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyPhoneMessageTextRequest): GetDefaultVerifyPhoneMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyPhoneMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyPhoneMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultVerifyPhoneMessageTextRequest, reader: jspb.BinaryReader): GetDefaultVerifyPhoneMessageTextRequest;
}

export namespace GetDefaultVerifyPhoneMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultVerifyPhoneMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultVerifyPhoneMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultVerifyPhoneMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyPhoneMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyPhoneMessageTextResponse): GetDefaultVerifyPhoneMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyPhoneMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyPhoneMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultVerifyPhoneMessageTextResponse, reader: jspb.BinaryReader): GetDefaultVerifyPhoneMessageTextResponse;
}

export namespace GetDefaultVerifyPhoneMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetCustomVerifyPhoneMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomVerifyPhoneMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyPhoneMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyPhoneMessageTextRequest): GetCustomVerifyPhoneMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyPhoneMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyPhoneMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomVerifyPhoneMessageTextRequest, reader: jspb.BinaryReader): GetCustomVerifyPhoneMessageTextRequest;
}

export namespace GetCustomVerifyPhoneMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomVerifyPhoneMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomVerifyPhoneMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomVerifyPhoneMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyPhoneMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyPhoneMessageTextResponse): GetCustomVerifyPhoneMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyPhoneMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyPhoneMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomVerifyPhoneMessageTextResponse, reader: jspb.BinaryReader): GetCustomVerifyPhoneMessageTextResponse;
}

export namespace GetCustomVerifyPhoneMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultVerifyPhoneMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultVerifyPhoneMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetDefaultVerifyPhoneMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetDefaultVerifyPhoneMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetDefaultVerifyPhoneMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetDefaultVerifyPhoneMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultVerifyPhoneMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetDefaultVerifyPhoneMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetDefaultVerifyPhoneMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultVerifyPhoneMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultVerifyPhoneMessageTextRequest): SetDefaultVerifyPhoneMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultVerifyPhoneMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultVerifyPhoneMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultVerifyPhoneMessageTextRequest, reader: jspb.BinaryReader): SetDefaultVerifyPhoneMessageTextRequest;
}

export namespace SetDefaultVerifyPhoneMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetDefaultVerifyPhoneMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultVerifyPhoneMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultVerifyPhoneMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultVerifyPhoneMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultVerifyPhoneMessageTextResponse): SetDefaultVerifyPhoneMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultVerifyPhoneMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultVerifyPhoneMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultVerifyPhoneMessageTextResponse, reader: jspb.BinaryReader): SetDefaultVerifyPhoneMessageTextResponse;
}

export namespace SetDefaultVerifyPhoneMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomVerifyPhoneMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomVerifyPhoneMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyPhoneMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyPhoneMessageTextToDefaultRequest): ResetCustomVerifyPhoneMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyPhoneMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyPhoneMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomVerifyPhoneMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomVerifyPhoneMessageTextToDefaultRequest;
}

export namespace ResetCustomVerifyPhoneMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomVerifyPhoneMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomVerifyPhoneMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomVerifyPhoneMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyPhoneMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyPhoneMessageTextToDefaultResponse): ResetCustomVerifyPhoneMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyPhoneMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyPhoneMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomVerifyPhoneMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomVerifyPhoneMessageTextToDefaultResponse;
}

export namespace ResetCustomVerifyPhoneMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomVerifySMSOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomVerifySMSOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifySMSOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifySMSOTPMessageTextRequest): GetCustomVerifySMSOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifySMSOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifySMSOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomVerifySMSOTPMessageTextRequest, reader: jspb.BinaryReader): GetCustomVerifySMSOTPMessageTextRequest;
}

export namespace GetCustomVerifySMSOTPMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomVerifySMSOTPMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomVerifySMSOTPMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomVerifySMSOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifySMSOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifySMSOTPMessageTextResponse): GetCustomVerifySMSOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifySMSOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifySMSOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomVerifySMSOTPMessageTextResponse, reader: jspb.BinaryReader): GetCustomVerifySMSOTPMessageTextResponse;
}

export namespace GetCustomVerifySMSOTPMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultVerifySMSOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultVerifySMSOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifySMSOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifySMSOTPMessageTextRequest): GetDefaultVerifySMSOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifySMSOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifySMSOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultVerifySMSOTPMessageTextRequest, reader: jspb.BinaryReader): GetDefaultVerifySMSOTPMessageTextRequest;
}

export namespace GetDefaultVerifySMSOTPMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultVerifySMSOTPMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultVerifySMSOTPMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultVerifySMSOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifySMSOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifySMSOTPMessageTextResponse): GetDefaultVerifySMSOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifySMSOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifySMSOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultVerifySMSOTPMessageTextResponse, reader: jspb.BinaryReader): GetDefaultVerifySMSOTPMessageTextResponse;
}

export namespace GetDefaultVerifySMSOTPMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultVerifySMSOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultVerifySMSOTPMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultVerifySMSOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultVerifySMSOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultVerifySMSOTPMessageTextRequest): SetDefaultVerifySMSOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultVerifySMSOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultVerifySMSOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultVerifySMSOTPMessageTextRequest, reader: jspb.BinaryReader): SetDefaultVerifySMSOTPMessageTextRequest;
}

export namespace SetDefaultVerifySMSOTPMessageTextRequest {
  export type AsObject = {
    language: string,
    text: string,
  }
}

export class SetDefaultVerifySMSOTPMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultVerifySMSOTPMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultVerifySMSOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultVerifySMSOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultVerifySMSOTPMessageTextResponse): SetDefaultVerifySMSOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultVerifySMSOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultVerifySMSOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultVerifySMSOTPMessageTextResponse, reader: jspb.BinaryReader): SetDefaultVerifySMSOTPMessageTextResponse;
}

export namespace SetDefaultVerifySMSOTPMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomVerifySMSOTPMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomVerifySMSOTPMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifySMSOTPMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifySMSOTPMessageTextToDefaultRequest): ResetCustomVerifySMSOTPMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifySMSOTPMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifySMSOTPMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomVerifySMSOTPMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomVerifySMSOTPMessageTextToDefaultRequest;
}

export namespace ResetCustomVerifySMSOTPMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomVerifySMSOTPMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomVerifySMSOTPMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomVerifySMSOTPMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifySMSOTPMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifySMSOTPMessageTextToDefaultResponse): ResetCustomVerifySMSOTPMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifySMSOTPMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifySMSOTPMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomVerifySMSOTPMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomVerifySMSOTPMessageTextToDefaultResponse;
}

export namespace ResetCustomVerifySMSOTPMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomVerifyEmailOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomVerifyEmailOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyEmailOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyEmailOTPMessageTextRequest): GetCustomVerifyEmailOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyEmailOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyEmailOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomVerifyEmailOTPMessageTextRequest, reader: jspb.BinaryReader): GetCustomVerifyEmailOTPMessageTextRequest;
}

export namespace GetCustomVerifyEmailOTPMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomVerifyEmailOTPMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomVerifyEmailOTPMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomVerifyEmailOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyEmailOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyEmailOTPMessageTextResponse): GetCustomVerifyEmailOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyEmailOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyEmailOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomVerifyEmailOTPMessageTextResponse, reader: jspb.BinaryReader): GetCustomVerifyEmailOTPMessageTextResponse;
}

export namespace GetCustomVerifyEmailOTPMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultVerifyEmailOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultVerifyEmailOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyEmailOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyEmailOTPMessageTextRequest): GetDefaultVerifyEmailOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyEmailOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyEmailOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultVerifyEmailOTPMessageTextRequest, reader: jspb.BinaryReader): GetDefaultVerifyEmailOTPMessageTextRequest;
}

export namespace GetDefaultVerifyEmailOTPMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultVerifyEmailOTPMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultVerifyEmailOTPMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultVerifyEmailOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyEmailOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyEmailOTPMessageTextResponse): GetDefaultVerifyEmailOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyEmailOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyEmailOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultVerifyEmailOTPMessageTextResponse, reader: jspb.BinaryReader): GetDefaultVerifyEmailOTPMessageTextResponse;
}

export namespace GetDefaultVerifyEmailOTPMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultVerifyEmailOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultVerifyEmailOTPMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetDefaultVerifyEmailOTPMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetDefaultVerifyEmailOTPMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetDefaultVerifyEmailOTPMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetDefaultVerifyEmailOTPMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultVerifyEmailOTPMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetDefaultVerifyEmailOTPMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetDefaultVerifyEmailOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultVerifyEmailOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultVerifyEmailOTPMessageTextRequest): SetDefaultVerifyEmailOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultVerifyEmailOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultVerifyEmailOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultVerifyEmailOTPMessageTextRequest, reader: jspb.BinaryReader): SetDefaultVerifyEmailOTPMessageTextRequest;
}

export namespace SetDefaultVerifyEmailOTPMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetDefaultVerifyEmailOTPMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultVerifyEmailOTPMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultVerifyEmailOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultVerifyEmailOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultVerifyEmailOTPMessageTextResponse): SetDefaultVerifyEmailOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultVerifyEmailOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultVerifyEmailOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultVerifyEmailOTPMessageTextResponse, reader: jspb.BinaryReader): SetDefaultVerifyEmailOTPMessageTextResponse;
}

export namespace SetDefaultVerifyEmailOTPMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomVerifyEmailOTPMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyEmailOTPMessageTextToDefaultRequest): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyEmailOTPMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomVerifyEmailOTPMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest;
}

export namespace ResetCustomVerifyEmailOTPMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomVerifyEmailOTPMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyEmailOTPMessageTextToDefaultResponse): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyEmailOTPMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomVerifyEmailOTPMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse;
}

export namespace ResetCustomVerifyEmailOTPMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultDomainClaimedMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultDomainClaimedMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultDomainClaimedMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultDomainClaimedMessageTextRequest): GetDefaultDomainClaimedMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultDomainClaimedMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultDomainClaimedMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultDomainClaimedMessageTextRequest, reader: jspb.BinaryReader): GetDefaultDomainClaimedMessageTextRequest;
}

export namespace GetDefaultDomainClaimedMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultDomainClaimedMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultDomainClaimedMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultDomainClaimedMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultDomainClaimedMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultDomainClaimedMessageTextResponse): GetDefaultDomainClaimedMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultDomainClaimedMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultDomainClaimedMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultDomainClaimedMessageTextResponse, reader: jspb.BinaryReader): GetDefaultDomainClaimedMessageTextResponse;
}

export namespace GetDefaultDomainClaimedMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetCustomDomainClaimedMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomDomainClaimedMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomDomainClaimedMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomDomainClaimedMessageTextRequest): GetCustomDomainClaimedMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomDomainClaimedMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomDomainClaimedMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomDomainClaimedMessageTextRequest, reader: jspb.BinaryReader): GetCustomDomainClaimedMessageTextRequest;
}

export namespace GetCustomDomainClaimedMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomDomainClaimedMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomDomainClaimedMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomDomainClaimedMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomDomainClaimedMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomDomainClaimedMessageTextResponse): GetCustomDomainClaimedMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomDomainClaimedMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomDomainClaimedMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomDomainClaimedMessageTextResponse, reader: jspb.BinaryReader): GetCustomDomainClaimedMessageTextResponse;
}

export namespace GetCustomDomainClaimedMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultDomainClaimedMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultDomainClaimedMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetDefaultDomainClaimedMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetDefaultDomainClaimedMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetDefaultDomainClaimedMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetDefaultDomainClaimedMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultDomainClaimedMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetDefaultDomainClaimedMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetDefaultDomainClaimedMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultDomainClaimedMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultDomainClaimedMessageTextRequest): SetDefaultDomainClaimedMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultDomainClaimedMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultDomainClaimedMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultDomainClaimedMessageTextRequest, reader: jspb.BinaryReader): SetDefaultDomainClaimedMessageTextRequest;
}

export namespace SetDefaultDomainClaimedMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetDefaultDomainClaimedMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultDomainClaimedMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultDomainClaimedMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultDomainClaimedMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultDomainClaimedMessageTextResponse): SetDefaultDomainClaimedMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultDomainClaimedMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultDomainClaimedMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultDomainClaimedMessageTextResponse, reader: jspb.BinaryReader): SetDefaultDomainClaimedMessageTextResponse;
}

export namespace SetDefaultDomainClaimedMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomDomainClaimedMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomDomainClaimedMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomDomainClaimedMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomDomainClaimedMessageTextToDefaultRequest): ResetCustomDomainClaimedMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomDomainClaimedMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomDomainClaimedMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomDomainClaimedMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomDomainClaimedMessageTextToDefaultRequest;
}

export namespace ResetCustomDomainClaimedMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomDomainClaimedMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomDomainClaimedMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomDomainClaimedMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomDomainClaimedMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomDomainClaimedMessageTextToDefaultResponse): ResetCustomDomainClaimedMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomDomainClaimedMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomDomainClaimedMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomDomainClaimedMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomDomainClaimedMessageTextToDefaultResponse;
}

export namespace ResetCustomDomainClaimedMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultPasswordChangeMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultPasswordChangeMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordChangeMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordChangeMessageTextRequest): GetDefaultPasswordChangeMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordChangeMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordChangeMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultPasswordChangeMessageTextRequest, reader: jspb.BinaryReader): GetDefaultPasswordChangeMessageTextRequest;
}

export namespace GetDefaultPasswordChangeMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultPasswordChangeMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultPasswordChangeMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultPasswordChangeMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordChangeMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordChangeMessageTextResponse): GetDefaultPasswordChangeMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordChangeMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordChangeMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultPasswordChangeMessageTextResponse, reader: jspb.BinaryReader): GetDefaultPasswordChangeMessageTextResponse;
}

export namespace GetDefaultPasswordChangeMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetCustomPasswordChangeMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomPasswordChangeMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordChangeMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordChangeMessageTextRequest): GetCustomPasswordChangeMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordChangeMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordChangeMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomPasswordChangeMessageTextRequest, reader: jspb.BinaryReader): GetCustomPasswordChangeMessageTextRequest;
}

export namespace GetCustomPasswordChangeMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomPasswordChangeMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomPasswordChangeMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomPasswordChangeMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordChangeMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordChangeMessageTextResponse): GetCustomPasswordChangeMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordChangeMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordChangeMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomPasswordChangeMessageTextResponse, reader: jspb.BinaryReader): GetCustomPasswordChangeMessageTextResponse;
}

export namespace GetCustomPasswordChangeMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultPasswordChangeMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultPasswordChangeMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetDefaultPasswordChangeMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetDefaultPasswordChangeMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetDefaultPasswordChangeMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetDefaultPasswordChangeMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultPasswordChangeMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetDefaultPasswordChangeMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetDefaultPasswordChangeMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultPasswordChangeMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultPasswordChangeMessageTextRequest): SetDefaultPasswordChangeMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultPasswordChangeMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultPasswordChangeMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultPasswordChangeMessageTextRequest, reader: jspb.BinaryReader): SetDefaultPasswordChangeMessageTextRequest;
}

export namespace SetDefaultPasswordChangeMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetDefaultPasswordChangeMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultPasswordChangeMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultPasswordChangeMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultPasswordChangeMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultPasswordChangeMessageTextResponse): SetDefaultPasswordChangeMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultPasswordChangeMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultPasswordChangeMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultPasswordChangeMessageTextResponse, reader: jspb.BinaryReader): SetDefaultPasswordChangeMessageTextResponse;
}

export namespace SetDefaultPasswordChangeMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomPasswordChangeMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomPasswordChangeMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordChangeMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordChangeMessageTextToDefaultRequest): ResetCustomPasswordChangeMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordChangeMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordChangeMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomPasswordChangeMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomPasswordChangeMessageTextToDefaultRequest;
}

export namespace ResetCustomPasswordChangeMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomPasswordChangeMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomPasswordChangeMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomPasswordChangeMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordChangeMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordChangeMessageTextToDefaultResponse): ResetCustomPasswordChangeMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordChangeMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordChangeMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomPasswordChangeMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomPasswordChangeMessageTextToDefaultResponse;
}

export namespace ResetCustomPasswordChangeMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultPasswordlessRegistrationMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultPasswordlessRegistrationMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordlessRegistrationMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordlessRegistrationMessageTextRequest): GetDefaultPasswordlessRegistrationMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordlessRegistrationMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordlessRegistrationMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultPasswordlessRegistrationMessageTextRequest, reader: jspb.BinaryReader): GetDefaultPasswordlessRegistrationMessageTextRequest;
}

export namespace GetDefaultPasswordlessRegistrationMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultPasswordlessRegistrationMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultPasswordlessRegistrationMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultPasswordlessRegistrationMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordlessRegistrationMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordlessRegistrationMessageTextResponse): GetDefaultPasswordlessRegistrationMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordlessRegistrationMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordlessRegistrationMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultPasswordlessRegistrationMessageTextResponse, reader: jspb.BinaryReader): GetDefaultPasswordlessRegistrationMessageTextResponse;
}

export namespace GetDefaultPasswordlessRegistrationMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetCustomPasswordlessRegistrationMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomPasswordlessRegistrationMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordlessRegistrationMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordlessRegistrationMessageTextRequest): GetCustomPasswordlessRegistrationMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordlessRegistrationMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordlessRegistrationMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomPasswordlessRegistrationMessageTextRequest, reader: jspb.BinaryReader): GetCustomPasswordlessRegistrationMessageTextRequest;
}

export namespace GetCustomPasswordlessRegistrationMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomPasswordlessRegistrationMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomPasswordlessRegistrationMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomPasswordlessRegistrationMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordlessRegistrationMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordlessRegistrationMessageTextResponse): GetCustomPasswordlessRegistrationMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordlessRegistrationMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordlessRegistrationMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomPasswordlessRegistrationMessageTextResponse, reader: jspb.BinaryReader): GetCustomPasswordlessRegistrationMessageTextResponse;
}

export namespace GetCustomPasswordlessRegistrationMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetDefaultPasswordlessRegistrationMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetDefaultPasswordlessRegistrationMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetDefaultPasswordlessRegistrationMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetDefaultPasswordlessRegistrationMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetDefaultPasswordlessRegistrationMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetDefaultPasswordlessRegistrationMessageTextRequest;

  getText(): string;
  setText(value: string): SetDefaultPasswordlessRegistrationMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetDefaultPasswordlessRegistrationMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetDefaultPasswordlessRegistrationMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultPasswordlessRegistrationMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultPasswordlessRegistrationMessageTextRequest): SetDefaultPasswordlessRegistrationMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetDefaultPasswordlessRegistrationMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultPasswordlessRegistrationMessageTextRequest;
  static deserializeBinaryFromReader(message: SetDefaultPasswordlessRegistrationMessageTextRequest, reader: jspb.BinaryReader): SetDefaultPasswordlessRegistrationMessageTextRequest;
}

export namespace SetDefaultPasswordlessRegistrationMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetDefaultPasswordlessRegistrationMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetDefaultPasswordlessRegistrationMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetDefaultPasswordlessRegistrationMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetDefaultPasswordlessRegistrationMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetDefaultPasswordlessRegistrationMessageTextResponse): SetDefaultPasswordlessRegistrationMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetDefaultPasswordlessRegistrationMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetDefaultPasswordlessRegistrationMessageTextResponse;
  static deserializeBinaryFromReader(message: SetDefaultPasswordlessRegistrationMessageTextResponse, reader: jspb.BinaryReader): SetDefaultPasswordlessRegistrationMessageTextResponse;
}

export namespace SetDefaultPasswordlessRegistrationMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest;
}

export namespace ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse;
}

export namespace ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultLoginTextsRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultLoginTextsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLoginTextsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLoginTextsRequest): GetDefaultLoginTextsRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLoginTextsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLoginTextsRequest;
  static deserializeBinaryFromReader(message: GetDefaultLoginTextsRequest, reader: jspb.BinaryReader): GetDefaultLoginTextsRequest;
}

export namespace GetDefaultLoginTextsRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultLoginTextsResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.LoginCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.LoginCustomText): GetDefaultLoginTextsResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultLoginTextsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLoginTextsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLoginTextsResponse): GetDefaultLoginTextsResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLoginTextsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLoginTextsResponse;
  static deserializeBinaryFromReader(message: GetDefaultLoginTextsResponse, reader: jspb.BinaryReader): GetDefaultLoginTextsResponse;
}

export namespace GetDefaultLoginTextsResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.LoginCustomText.AsObject,
  }
}

export class GetCustomLoginTextsRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomLoginTextsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomLoginTextsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomLoginTextsRequest): GetCustomLoginTextsRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomLoginTextsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomLoginTextsRequest;
  static deserializeBinaryFromReader(message: GetCustomLoginTextsRequest, reader: jspb.BinaryReader): GetCustomLoginTextsRequest;
}

export namespace GetCustomLoginTextsRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomLoginTextsResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.LoginCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.LoginCustomText): GetCustomLoginTextsResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomLoginTextsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomLoginTextsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomLoginTextsResponse): GetCustomLoginTextsResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomLoginTextsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomLoginTextsResponse;
  static deserializeBinaryFromReader(message: GetCustomLoginTextsResponse, reader: jspb.BinaryReader): GetCustomLoginTextsResponse;
}

export namespace GetCustomLoginTextsResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.LoginCustomText.AsObject,
  }
}

export class SetCustomLoginTextsRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomLoginTextsRequest;

  getSelectAccountText(): zitadel_text_pb.SelectAccountScreenText | undefined;
  setSelectAccountText(value?: zitadel_text_pb.SelectAccountScreenText): SetCustomLoginTextsRequest;
  hasSelectAccountText(): boolean;
  clearSelectAccountText(): SetCustomLoginTextsRequest;

  getLoginText(): zitadel_text_pb.LoginScreenText | undefined;
  setLoginText(value?: zitadel_text_pb.LoginScreenText): SetCustomLoginTextsRequest;
  hasLoginText(): boolean;
  clearLoginText(): SetCustomLoginTextsRequest;

  getPasswordText(): zitadel_text_pb.PasswordScreenText | undefined;
  setPasswordText(value?: zitadel_text_pb.PasswordScreenText): SetCustomLoginTextsRequest;
  hasPasswordText(): boolean;
  clearPasswordText(): SetCustomLoginTextsRequest;

  getUsernameChangeText(): zitadel_text_pb.UsernameChangeScreenText | undefined;
  setUsernameChangeText(value?: zitadel_text_pb.UsernameChangeScreenText): SetCustomLoginTextsRequest;
  hasUsernameChangeText(): boolean;
  clearUsernameChangeText(): SetCustomLoginTextsRequest;

  getUsernameChangeDoneText(): zitadel_text_pb.UsernameChangeDoneScreenText | undefined;
  setUsernameChangeDoneText(value?: zitadel_text_pb.UsernameChangeDoneScreenText): SetCustomLoginTextsRequest;
  hasUsernameChangeDoneText(): boolean;
  clearUsernameChangeDoneText(): SetCustomLoginTextsRequest;

  getInitPasswordText(): zitadel_text_pb.InitPasswordScreenText | undefined;
  setInitPasswordText(value?: zitadel_text_pb.InitPasswordScreenText): SetCustomLoginTextsRequest;
  hasInitPasswordText(): boolean;
  clearInitPasswordText(): SetCustomLoginTextsRequest;

  getInitPasswordDoneText(): zitadel_text_pb.InitPasswordDoneScreenText | undefined;
  setInitPasswordDoneText(value?: zitadel_text_pb.InitPasswordDoneScreenText): SetCustomLoginTextsRequest;
  hasInitPasswordDoneText(): boolean;
  clearInitPasswordDoneText(): SetCustomLoginTextsRequest;

  getEmailVerificationText(): zitadel_text_pb.EmailVerificationScreenText | undefined;
  setEmailVerificationText(value?: zitadel_text_pb.EmailVerificationScreenText): SetCustomLoginTextsRequest;
  hasEmailVerificationText(): boolean;
  clearEmailVerificationText(): SetCustomLoginTextsRequest;

  getEmailVerificationDoneText(): zitadel_text_pb.EmailVerificationDoneScreenText | undefined;
  setEmailVerificationDoneText(value?: zitadel_text_pb.EmailVerificationDoneScreenText): SetCustomLoginTextsRequest;
  hasEmailVerificationDoneText(): boolean;
  clearEmailVerificationDoneText(): SetCustomLoginTextsRequest;

  getInitializeUserText(): zitadel_text_pb.InitializeUserScreenText | undefined;
  setInitializeUserText(value?: zitadel_text_pb.InitializeUserScreenText): SetCustomLoginTextsRequest;
  hasInitializeUserText(): boolean;
  clearInitializeUserText(): SetCustomLoginTextsRequest;

  getInitializeDoneText(): zitadel_text_pb.InitializeUserDoneScreenText | undefined;
  setInitializeDoneText(value?: zitadel_text_pb.InitializeUserDoneScreenText): SetCustomLoginTextsRequest;
  hasInitializeDoneText(): boolean;
  clearInitializeDoneText(): SetCustomLoginTextsRequest;

  getInitMfaPromptText(): zitadel_text_pb.InitMFAPromptScreenText | undefined;
  setInitMfaPromptText(value?: zitadel_text_pb.InitMFAPromptScreenText): SetCustomLoginTextsRequest;
  hasInitMfaPromptText(): boolean;
  clearInitMfaPromptText(): SetCustomLoginTextsRequest;

  getInitMfaOtpText(): zitadel_text_pb.InitMFAOTPScreenText | undefined;
  setInitMfaOtpText(value?: zitadel_text_pb.InitMFAOTPScreenText): SetCustomLoginTextsRequest;
  hasInitMfaOtpText(): boolean;
  clearInitMfaOtpText(): SetCustomLoginTextsRequest;

  getInitMfaU2fText(): zitadel_text_pb.InitMFAU2FScreenText | undefined;
  setInitMfaU2fText(value?: zitadel_text_pb.InitMFAU2FScreenText): SetCustomLoginTextsRequest;
  hasInitMfaU2fText(): boolean;
  clearInitMfaU2fText(): SetCustomLoginTextsRequest;

  getInitMfaDoneText(): zitadel_text_pb.InitMFADoneScreenText | undefined;
  setInitMfaDoneText(value?: zitadel_text_pb.InitMFADoneScreenText): SetCustomLoginTextsRequest;
  hasInitMfaDoneText(): boolean;
  clearInitMfaDoneText(): SetCustomLoginTextsRequest;

  getMfaProvidersText(): zitadel_text_pb.MFAProvidersText | undefined;
  setMfaProvidersText(value?: zitadel_text_pb.MFAProvidersText): SetCustomLoginTextsRequest;
  hasMfaProvidersText(): boolean;
  clearMfaProvidersText(): SetCustomLoginTextsRequest;

  getVerifyMfaOtpText(): zitadel_text_pb.VerifyMFAOTPScreenText | undefined;
  setVerifyMfaOtpText(value?: zitadel_text_pb.VerifyMFAOTPScreenText): SetCustomLoginTextsRequest;
  hasVerifyMfaOtpText(): boolean;
  clearVerifyMfaOtpText(): SetCustomLoginTextsRequest;

  getVerifyMfaU2fText(): zitadel_text_pb.VerifyMFAU2FScreenText | undefined;
  setVerifyMfaU2fText(value?: zitadel_text_pb.VerifyMFAU2FScreenText): SetCustomLoginTextsRequest;
  hasVerifyMfaU2fText(): boolean;
  clearVerifyMfaU2fText(): SetCustomLoginTextsRequest;

  getPasswordlessText(): zitadel_text_pb.PasswordlessScreenText | undefined;
  setPasswordlessText(value?: zitadel_text_pb.PasswordlessScreenText): SetCustomLoginTextsRequest;
  hasPasswordlessText(): boolean;
  clearPasswordlessText(): SetCustomLoginTextsRequest;

  getPasswordChangeText(): zitadel_text_pb.PasswordChangeScreenText | undefined;
  setPasswordChangeText(value?: zitadel_text_pb.PasswordChangeScreenText): SetCustomLoginTextsRequest;
  hasPasswordChangeText(): boolean;
  clearPasswordChangeText(): SetCustomLoginTextsRequest;

  getPasswordChangeDoneText(): zitadel_text_pb.PasswordChangeDoneScreenText | undefined;
  setPasswordChangeDoneText(value?: zitadel_text_pb.PasswordChangeDoneScreenText): SetCustomLoginTextsRequest;
  hasPasswordChangeDoneText(): boolean;
  clearPasswordChangeDoneText(): SetCustomLoginTextsRequest;

  getPasswordResetDoneText(): zitadel_text_pb.PasswordResetDoneScreenText | undefined;
  setPasswordResetDoneText(value?: zitadel_text_pb.PasswordResetDoneScreenText): SetCustomLoginTextsRequest;
  hasPasswordResetDoneText(): boolean;
  clearPasswordResetDoneText(): SetCustomLoginTextsRequest;

  getRegistrationOptionText(): zitadel_text_pb.RegistrationOptionScreenText | undefined;
  setRegistrationOptionText(value?: zitadel_text_pb.RegistrationOptionScreenText): SetCustomLoginTextsRequest;
  hasRegistrationOptionText(): boolean;
  clearRegistrationOptionText(): SetCustomLoginTextsRequest;

  getRegistrationUserText(): zitadel_text_pb.RegistrationUserScreenText | undefined;
  setRegistrationUserText(value?: zitadel_text_pb.RegistrationUserScreenText): SetCustomLoginTextsRequest;
  hasRegistrationUserText(): boolean;
  clearRegistrationUserText(): SetCustomLoginTextsRequest;

  getRegistrationOrgText(): zitadel_text_pb.RegistrationOrgScreenText | undefined;
  setRegistrationOrgText(value?: zitadel_text_pb.RegistrationOrgScreenText): SetCustomLoginTextsRequest;
  hasRegistrationOrgText(): boolean;
  clearRegistrationOrgText(): SetCustomLoginTextsRequest;

  getLinkingUserDoneText(): zitadel_text_pb.LinkingUserDoneScreenText | undefined;
  setLinkingUserDoneText(value?: zitadel_text_pb.LinkingUserDoneScreenText): SetCustomLoginTextsRequest;
  hasLinkingUserDoneText(): boolean;
  clearLinkingUserDoneText(): SetCustomLoginTextsRequest;

  getExternalUserNotFoundText(): zitadel_text_pb.ExternalUserNotFoundScreenText | undefined;
  setExternalUserNotFoundText(value?: zitadel_text_pb.ExternalUserNotFoundScreenText): SetCustomLoginTextsRequest;
  hasExternalUserNotFoundText(): boolean;
  clearExternalUserNotFoundText(): SetCustomLoginTextsRequest;

  getSuccessLoginText(): zitadel_text_pb.SuccessLoginScreenText | undefined;
  setSuccessLoginText(value?: zitadel_text_pb.SuccessLoginScreenText): SetCustomLoginTextsRequest;
  hasSuccessLoginText(): boolean;
  clearSuccessLoginText(): SetCustomLoginTextsRequest;

  getLogoutText(): zitadel_text_pb.LogoutDoneScreenText | undefined;
  setLogoutText(value?: zitadel_text_pb.LogoutDoneScreenText): SetCustomLoginTextsRequest;
  hasLogoutText(): boolean;
  clearLogoutText(): SetCustomLoginTextsRequest;

  getFooterText(): zitadel_text_pb.FooterText | undefined;
  setFooterText(value?: zitadel_text_pb.FooterText): SetCustomLoginTextsRequest;
  hasFooterText(): boolean;
  clearFooterText(): SetCustomLoginTextsRequest;

  getPasswordlessPromptText(): zitadel_text_pb.PasswordlessPromptScreenText | undefined;
  setPasswordlessPromptText(value?: zitadel_text_pb.PasswordlessPromptScreenText): SetCustomLoginTextsRequest;
  hasPasswordlessPromptText(): boolean;
  clearPasswordlessPromptText(): SetCustomLoginTextsRequest;

  getPasswordlessRegistrationText(): zitadel_text_pb.PasswordlessRegistrationScreenText | undefined;
  setPasswordlessRegistrationText(value?: zitadel_text_pb.PasswordlessRegistrationScreenText): SetCustomLoginTextsRequest;
  hasPasswordlessRegistrationText(): boolean;
  clearPasswordlessRegistrationText(): SetCustomLoginTextsRequest;

  getPasswordlessRegistrationDoneText(): zitadel_text_pb.PasswordlessRegistrationDoneScreenText | undefined;
  setPasswordlessRegistrationDoneText(value?: zitadel_text_pb.PasswordlessRegistrationDoneScreenText): SetCustomLoginTextsRequest;
  hasPasswordlessRegistrationDoneText(): boolean;
  clearPasswordlessRegistrationDoneText(): SetCustomLoginTextsRequest;

  getExternalRegistrationUserOverviewText(): zitadel_text_pb.ExternalRegistrationUserOverviewScreenText | undefined;
  setExternalRegistrationUserOverviewText(value?: zitadel_text_pb.ExternalRegistrationUserOverviewScreenText): SetCustomLoginTextsRequest;
  hasExternalRegistrationUserOverviewText(): boolean;
  clearExternalRegistrationUserOverviewText(): SetCustomLoginTextsRequest;

  getLinkingUserPromptText(): zitadel_text_pb.LinkingUserPromptScreenText | undefined;
  setLinkingUserPromptText(value?: zitadel_text_pb.LinkingUserPromptScreenText): SetCustomLoginTextsRequest;
  hasLinkingUserPromptText(): boolean;
  clearLinkingUserPromptText(): SetCustomLoginTextsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomLoginTextsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomLoginTextsRequest): SetCustomLoginTextsRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomLoginTextsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomLoginTextsRequest;
  static deserializeBinaryFromReader(message: SetCustomLoginTextsRequest, reader: jspb.BinaryReader): SetCustomLoginTextsRequest;
}

export namespace SetCustomLoginTextsRequest {
  export type AsObject = {
    language: string,
    selectAccountText?: zitadel_text_pb.SelectAccountScreenText.AsObject,
    loginText?: zitadel_text_pb.LoginScreenText.AsObject,
    passwordText?: zitadel_text_pb.PasswordScreenText.AsObject,
    usernameChangeText?: zitadel_text_pb.UsernameChangeScreenText.AsObject,
    usernameChangeDoneText?: zitadel_text_pb.UsernameChangeDoneScreenText.AsObject,
    initPasswordText?: zitadel_text_pb.InitPasswordScreenText.AsObject,
    initPasswordDoneText?: zitadel_text_pb.InitPasswordDoneScreenText.AsObject,
    emailVerificationText?: zitadel_text_pb.EmailVerificationScreenText.AsObject,
    emailVerificationDoneText?: zitadel_text_pb.EmailVerificationDoneScreenText.AsObject,
    initializeUserText?: zitadel_text_pb.InitializeUserScreenText.AsObject,
    initializeDoneText?: zitadel_text_pb.InitializeUserDoneScreenText.AsObject,
    initMfaPromptText?: zitadel_text_pb.InitMFAPromptScreenText.AsObject,
    initMfaOtpText?: zitadel_text_pb.InitMFAOTPScreenText.AsObject,
    initMfaU2fText?: zitadel_text_pb.InitMFAU2FScreenText.AsObject,
    initMfaDoneText?: zitadel_text_pb.InitMFADoneScreenText.AsObject,
    mfaProvidersText?: zitadel_text_pb.MFAProvidersText.AsObject,
    verifyMfaOtpText?: zitadel_text_pb.VerifyMFAOTPScreenText.AsObject,
    verifyMfaU2fText?: zitadel_text_pb.VerifyMFAU2FScreenText.AsObject,
    passwordlessText?: zitadel_text_pb.PasswordlessScreenText.AsObject,
    passwordChangeText?: zitadel_text_pb.PasswordChangeScreenText.AsObject,
    passwordChangeDoneText?: zitadel_text_pb.PasswordChangeDoneScreenText.AsObject,
    passwordResetDoneText?: zitadel_text_pb.PasswordResetDoneScreenText.AsObject,
    registrationOptionText?: zitadel_text_pb.RegistrationOptionScreenText.AsObject,
    registrationUserText?: zitadel_text_pb.RegistrationUserScreenText.AsObject,
    registrationOrgText?: zitadel_text_pb.RegistrationOrgScreenText.AsObject,
    linkingUserDoneText?: zitadel_text_pb.LinkingUserDoneScreenText.AsObject,
    externalUserNotFoundText?: zitadel_text_pb.ExternalUserNotFoundScreenText.AsObject,
    successLoginText?: zitadel_text_pb.SuccessLoginScreenText.AsObject,
    logoutText?: zitadel_text_pb.LogoutDoneScreenText.AsObject,
    footerText?: zitadel_text_pb.FooterText.AsObject,
    passwordlessPromptText?: zitadel_text_pb.PasswordlessPromptScreenText.AsObject,
    passwordlessRegistrationText?: zitadel_text_pb.PasswordlessRegistrationScreenText.AsObject,
    passwordlessRegistrationDoneText?: zitadel_text_pb.PasswordlessRegistrationDoneScreenText.AsObject,
    externalRegistrationUserOverviewText?: zitadel_text_pb.ExternalRegistrationUserOverviewScreenText.AsObject,
    linkingUserPromptText?: zitadel_text_pb.LinkingUserPromptScreenText.AsObject,
  }
}

export class SetCustomLoginTextsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomLoginTextsResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomLoginTextsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomLoginTextsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomLoginTextsResponse): SetCustomLoginTextsResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomLoginTextsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomLoginTextsResponse;
  static deserializeBinaryFromReader(message: SetCustomLoginTextsResponse, reader: jspb.BinaryReader): SetCustomLoginTextsResponse;
}

export namespace SetCustomLoginTextsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomLoginTextsToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomLoginTextsToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomLoginTextsToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomLoginTextsToDefaultRequest): ResetCustomLoginTextsToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomLoginTextsToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomLoginTextsToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomLoginTextsToDefaultRequest, reader: jspb.BinaryReader): ResetCustomLoginTextsToDefaultRequest;
}

export namespace ResetCustomLoginTextsToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomLoginTextsToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomLoginTextsToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomLoginTextsToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomLoginTextsToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomLoginTextsToDefaultResponse): ResetCustomLoginTextsToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomLoginTextsToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomLoginTextsToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomLoginTextsToDefaultResponse, reader: jspb.BinaryReader): ResetCustomLoginTextsToDefaultResponse;
}

export namespace ResetCustomLoginTextsToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddIAMMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddIAMMemberRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): AddIAMMemberRequest;
  clearRolesList(): AddIAMMemberRequest;
  addRoles(value: string, index?: number): AddIAMMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIAMMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddIAMMemberRequest): AddIAMMemberRequest.AsObject;
  static serializeBinaryToWriter(message: AddIAMMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIAMMemberRequest;
  static deserializeBinaryFromReader(message: AddIAMMemberRequest, reader: jspb.BinaryReader): AddIAMMemberRequest;
}

export namespace AddIAMMemberRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
  }
}

export class AddIAMMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddIAMMemberResponse;
  hasDetails(): boolean;
  clearDetails(): AddIAMMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIAMMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddIAMMemberResponse): AddIAMMemberResponse.AsObject;
  static serializeBinaryToWriter(message: AddIAMMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIAMMemberResponse;
  static deserializeBinaryFromReader(message: AddIAMMemberResponse, reader: jspb.BinaryReader): AddIAMMemberResponse;
}

export namespace AddIAMMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateIAMMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateIAMMemberRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): UpdateIAMMemberRequest;
  clearRolesList(): UpdateIAMMemberRequest;
  addRoles(value: string, index?: number): UpdateIAMMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateIAMMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateIAMMemberRequest): UpdateIAMMemberRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateIAMMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateIAMMemberRequest;
  static deserializeBinaryFromReader(message: UpdateIAMMemberRequest, reader: jspb.BinaryReader): UpdateIAMMemberRequest;
}

export namespace UpdateIAMMemberRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
  }
}

export class UpdateIAMMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateIAMMemberResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateIAMMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateIAMMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateIAMMemberResponse): UpdateIAMMemberResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateIAMMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateIAMMemberResponse;
  static deserializeBinaryFromReader(message: UpdateIAMMemberResponse, reader: jspb.BinaryReader): UpdateIAMMemberResponse;
}

export namespace UpdateIAMMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveIAMMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveIAMMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIAMMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIAMMemberRequest): RemoveIAMMemberRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveIAMMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIAMMemberRequest;
  static deserializeBinaryFromReader(message: RemoveIAMMemberRequest, reader: jspb.BinaryReader): RemoveIAMMemberRequest;
}

export namespace RemoveIAMMemberRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveIAMMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveIAMMemberResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveIAMMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIAMMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIAMMemberResponse): RemoveIAMMemberResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveIAMMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIAMMemberResponse;
  static deserializeBinaryFromReader(message: RemoveIAMMemberResponse, reader: jspb.BinaryReader): RemoveIAMMemberResponse;
}

export namespace RemoveIAMMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListIAMMemberRolesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListIAMMemberRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListIAMMemberRolesRequest): ListIAMMemberRolesRequest.AsObject;
  static serializeBinaryToWriter(message: ListIAMMemberRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListIAMMemberRolesRequest;
  static deserializeBinaryFromReader(message: ListIAMMemberRolesRequest, reader: jspb.BinaryReader): ListIAMMemberRolesRequest;
}

export namespace ListIAMMemberRolesRequest {
  export type AsObject = {
  }
}

export class ListIAMMemberRolesResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListIAMMemberRolesResponse;
  hasDetails(): boolean;
  clearDetails(): ListIAMMemberRolesResponse;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): ListIAMMemberRolesResponse;
  clearRolesList(): ListIAMMemberRolesResponse;
  addRoles(value: string, index?: number): ListIAMMemberRolesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListIAMMemberRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListIAMMemberRolesResponse): ListIAMMemberRolesResponse.AsObject;
  static serializeBinaryToWriter(message: ListIAMMemberRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListIAMMemberRolesResponse;
  static deserializeBinaryFromReader(message: ListIAMMemberRolesResponse, reader: jspb.BinaryReader): ListIAMMemberRolesResponse;
}

export namespace ListIAMMemberRolesResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    rolesList: Array<string>,
  }
}

export class ListIAMMembersRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListIAMMembersRequest;
  hasQuery(): boolean;
  clearQuery(): ListIAMMembersRequest;

  getQueriesList(): Array<zitadel_member_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_member_pb.SearchQuery>): ListIAMMembersRequest;
  clearQueriesList(): ListIAMMembersRequest;
  addQueries(value?: zitadel_member_pb.SearchQuery, index?: number): zitadel_member_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListIAMMembersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListIAMMembersRequest): ListIAMMembersRequest.AsObject;
  static serializeBinaryToWriter(message: ListIAMMembersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListIAMMembersRequest;
  static deserializeBinaryFromReader(message: ListIAMMembersRequest, reader: jspb.BinaryReader): ListIAMMembersRequest;
}

export namespace ListIAMMembersRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_member_pb.SearchQuery.AsObject>,
  }
}

export class ListIAMMembersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListIAMMembersResponse;
  hasDetails(): boolean;
  clearDetails(): ListIAMMembersResponse;

  getResultList(): Array<zitadel_member_pb.Member>;
  setResultList(value: Array<zitadel_member_pb.Member>): ListIAMMembersResponse;
  clearResultList(): ListIAMMembersResponse;
  addResult(value?: zitadel_member_pb.Member, index?: number): zitadel_member_pb.Member;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListIAMMembersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListIAMMembersResponse): ListIAMMembersResponse.AsObject;
  static serializeBinaryToWriter(message: ListIAMMembersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListIAMMembersResponse;
  static deserializeBinaryFromReader(message: ListIAMMembersResponse, reader: jspb.BinaryReader): ListIAMMembersResponse;
}

export namespace ListIAMMembersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_member_pb.Member.AsObject>,
  }
}

export class ListViewsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListViewsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListViewsRequest): ListViewsRequest.AsObject;
  static serializeBinaryToWriter(message: ListViewsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListViewsRequest;
  static deserializeBinaryFromReader(message: ListViewsRequest, reader: jspb.BinaryReader): ListViewsRequest;
}

export namespace ListViewsRequest {
  export type AsObject = {
  }
}

export class ListViewsResponse extends jspb.Message {
  getResultList(): Array<View>;
  setResultList(value: Array<View>): ListViewsResponse;
  clearResultList(): ListViewsResponse;
  addResult(value?: View, index?: number): View;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListViewsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListViewsResponse): ListViewsResponse.AsObject;
  static serializeBinaryToWriter(message: ListViewsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListViewsResponse;
  static deserializeBinaryFromReader(message: ListViewsResponse, reader: jspb.BinaryReader): ListViewsResponse;
}

export namespace ListViewsResponse {
  export type AsObject = {
    resultList: Array<View.AsObject>,
  }
}

export class ListFailedEventsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFailedEventsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFailedEventsRequest): ListFailedEventsRequest.AsObject;
  static serializeBinaryToWriter(message: ListFailedEventsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFailedEventsRequest;
  static deserializeBinaryFromReader(message: ListFailedEventsRequest, reader: jspb.BinaryReader): ListFailedEventsRequest;
}

export namespace ListFailedEventsRequest {
  export type AsObject = {
  }
}

export class ListFailedEventsResponse extends jspb.Message {
  getResultList(): Array<FailedEvent>;
  setResultList(value: Array<FailedEvent>): ListFailedEventsResponse;
  clearResultList(): ListFailedEventsResponse;
  addResult(value?: FailedEvent, index?: number): FailedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFailedEventsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFailedEventsResponse): ListFailedEventsResponse.AsObject;
  static serializeBinaryToWriter(message: ListFailedEventsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFailedEventsResponse;
  static deserializeBinaryFromReader(message: ListFailedEventsResponse, reader: jspb.BinaryReader): ListFailedEventsResponse;
}

export namespace ListFailedEventsResponse {
  export type AsObject = {
    resultList: Array<FailedEvent.AsObject>,
  }
}

export class RemoveFailedEventRequest extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): RemoveFailedEventRequest;

  getViewName(): string;
  setViewName(value: string): RemoveFailedEventRequest;

  getFailedSequence(): number;
  setFailedSequence(value: number): RemoveFailedEventRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveFailedEventRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveFailedEventRequest): RemoveFailedEventRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveFailedEventRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveFailedEventRequest;
  static deserializeBinaryFromReader(message: RemoveFailedEventRequest, reader: jspb.BinaryReader): RemoveFailedEventRequest;
}

export namespace RemoveFailedEventRequest {
  export type AsObject = {
    database: string,
    viewName: string,
    failedSequence: number,
  }
}

export class RemoveFailedEventResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveFailedEventResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveFailedEventResponse): RemoveFailedEventResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveFailedEventResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveFailedEventResponse;
  static deserializeBinaryFromReader(message: RemoveFailedEventResponse, reader: jspb.BinaryReader): RemoveFailedEventResponse;
}

export namespace RemoveFailedEventResponse {
  export type AsObject = {
  }
}

export class View extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): View;

  getViewName(): string;
  setViewName(value: string): View;

  getProcessedSequence(): number;
  setProcessedSequence(value: number): View;

  getEventTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setEventTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): View;
  hasEventTimestamp(): boolean;
  clearEventTimestamp(): View;

  getLastSuccessfulSpoolerRun(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastSuccessfulSpoolerRun(value?: google_protobuf_timestamp_pb.Timestamp): View;
  hasLastSuccessfulSpoolerRun(): boolean;
  clearLastSuccessfulSpoolerRun(): View;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): View.AsObject;
  static toObject(includeInstance: boolean, msg: View): View.AsObject;
  static serializeBinaryToWriter(message: View, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): View;
  static deserializeBinaryFromReader(message: View, reader: jspb.BinaryReader): View;
}

export namespace View {
  export type AsObject = {
    database: string,
    viewName: string,
    processedSequence: number,
    eventTimestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastSuccessfulSpoolerRun?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class FailedEvent extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): FailedEvent;

  getViewName(): string;
  setViewName(value: string): FailedEvent;

  getFailedSequence(): number;
  setFailedSequence(value: number): FailedEvent;

  getFailureCount(): number;
  setFailureCount(value: number): FailedEvent;

  getErrorMessage(): string;
  setErrorMessage(value: string): FailedEvent;

  getLastFailed(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastFailed(value?: google_protobuf_timestamp_pb.Timestamp): FailedEvent;
  hasLastFailed(): boolean;
  clearLastFailed(): FailedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FailedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: FailedEvent): FailedEvent.AsObject;
  static serializeBinaryToWriter(message: FailedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FailedEvent;
  static deserializeBinaryFromReader(message: FailedEvent, reader: jspb.BinaryReader): FailedEvent;
}

export namespace FailedEvent {
  export type AsObject = {
    database: string,
    viewName: string,
    failedSequence: number,
    failureCount: number,
    errorMessage: string,
    lastFailed?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class ImportDataRequest extends jspb.Message {
  getDataOrgs(): ImportDataOrg | undefined;
  setDataOrgs(value?: ImportDataOrg): ImportDataRequest;
  hasDataOrgs(): boolean;
  clearDataOrgs(): ImportDataRequest;

  getDataOrgsv1(): zitadel_v1_pb.ImportDataOrg | undefined;
  setDataOrgsv1(value?: zitadel_v1_pb.ImportDataOrg): ImportDataRequest;
  hasDataOrgsv1(): boolean;
  clearDataOrgsv1(): ImportDataRequest;

  getDataOrgsLocal(): ImportDataRequest.LocalInput | undefined;
  setDataOrgsLocal(value?: ImportDataRequest.LocalInput): ImportDataRequest;
  hasDataOrgsLocal(): boolean;
  clearDataOrgsLocal(): ImportDataRequest;

  getDataOrgsv1Local(): ImportDataRequest.LocalInput | undefined;
  setDataOrgsv1Local(value?: ImportDataRequest.LocalInput): ImportDataRequest;
  hasDataOrgsv1Local(): boolean;
  clearDataOrgsv1Local(): ImportDataRequest;

  getDataOrgsS3(): ImportDataRequest.S3Input | undefined;
  setDataOrgsS3(value?: ImportDataRequest.S3Input): ImportDataRequest;
  hasDataOrgsS3(): boolean;
  clearDataOrgsS3(): ImportDataRequest;

  getDataOrgsv1S3(): ImportDataRequest.S3Input | undefined;
  setDataOrgsv1S3(value?: ImportDataRequest.S3Input): ImportDataRequest;
  hasDataOrgsv1S3(): boolean;
  clearDataOrgsv1S3(): ImportDataRequest;

  getDataOrgsGcs(): ImportDataRequest.GCSInput | undefined;
  setDataOrgsGcs(value?: ImportDataRequest.GCSInput): ImportDataRequest;
  hasDataOrgsGcs(): boolean;
  clearDataOrgsGcs(): ImportDataRequest;

  getDataOrgsv1Gcs(): ImportDataRequest.GCSInput | undefined;
  setDataOrgsv1Gcs(value?: ImportDataRequest.GCSInput): ImportDataRequest;
  hasDataOrgsv1Gcs(): boolean;
  clearDataOrgsv1Gcs(): ImportDataRequest;

  getTimeout(): string;
  setTimeout(value: string): ImportDataRequest;

  getDataCase(): ImportDataRequest.DataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataRequest): ImportDataRequest.AsObject;
  static serializeBinaryToWriter(message: ImportDataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataRequest;
  static deserializeBinaryFromReader(message: ImportDataRequest, reader: jspb.BinaryReader): ImportDataRequest;
}

export namespace ImportDataRequest {
  export type AsObject = {
    dataOrgs?: ImportDataOrg.AsObject,
    dataOrgsv1?: zitadel_v1_pb.ImportDataOrg.AsObject,
    dataOrgsLocal?: ImportDataRequest.LocalInput.AsObject,
    dataOrgsv1Local?: ImportDataRequest.LocalInput.AsObject,
    dataOrgsS3?: ImportDataRequest.S3Input.AsObject,
    dataOrgsv1S3?: ImportDataRequest.S3Input.AsObject,
    dataOrgsGcs?: ImportDataRequest.GCSInput.AsObject,
    dataOrgsv1Gcs?: ImportDataRequest.GCSInput.AsObject,
    timeout: string,
  }

  export class LocalInput extends jspb.Message {
    getPath(): string;
    setPath(value: string): LocalInput;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LocalInput.AsObject;
    static toObject(includeInstance: boolean, msg: LocalInput): LocalInput.AsObject;
    static serializeBinaryToWriter(message: LocalInput, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LocalInput;
    static deserializeBinaryFromReader(message: LocalInput, reader: jspb.BinaryReader): LocalInput;
  }

  export namespace LocalInput {
    export type AsObject = {
      path: string,
    }
  }


  export class S3Input extends jspb.Message {
    getPath(): string;
    setPath(value: string): S3Input;

    getEndpoint(): string;
    setEndpoint(value: string): S3Input;

    getAccessKeyId(): string;
    setAccessKeyId(value: string): S3Input;

    getSecretAccessKey(): string;
    setSecretAccessKey(value: string): S3Input;

    getSsl(): boolean;
    setSsl(value: boolean): S3Input;

    getBucket(): string;
    setBucket(value: string): S3Input;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): S3Input.AsObject;
    static toObject(includeInstance: boolean, msg: S3Input): S3Input.AsObject;
    static serializeBinaryToWriter(message: S3Input, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): S3Input;
    static deserializeBinaryFromReader(message: S3Input, reader: jspb.BinaryReader): S3Input;
  }

  export namespace S3Input {
    export type AsObject = {
      path: string,
      endpoint: string,
      accessKeyId: string,
      secretAccessKey: string,
      ssl: boolean,
      bucket: string,
    }
  }


  export class GCSInput extends jspb.Message {
    getBucket(): string;
    setBucket(value: string): GCSInput;

    getServiceaccountJson(): string;
    setServiceaccountJson(value: string): GCSInput;

    getPath(): string;
    setPath(value: string): GCSInput;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GCSInput.AsObject;
    static toObject(includeInstance: boolean, msg: GCSInput): GCSInput.AsObject;
    static serializeBinaryToWriter(message: GCSInput, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GCSInput;
    static deserializeBinaryFromReader(message: GCSInput, reader: jspb.BinaryReader): GCSInput;
  }

  export namespace GCSInput {
    export type AsObject = {
      bucket: string,
      serviceaccountJson: string,
      path: string,
    }
  }


  export enum DataCase { 
    DATA_NOT_SET = 0,
    DATA_ORGS = 1,
    DATA_ORGSV1 = 2,
    DATA_ORGS_LOCAL = 3,
    DATA_ORGSV1_LOCAL = 4,
    DATA_ORGS_S3 = 5,
    DATA_ORGSV1_S3 = 6,
    DATA_ORGS_GCS = 7,
    DATA_ORGSV1_GCS = 8,
  }
}

export class ImportDataOrg extends jspb.Message {
  getOrgsList(): Array<DataOrg>;
  setOrgsList(value: Array<DataOrg>): ImportDataOrg;
  clearOrgsList(): ImportDataOrg;
  addOrgs(value?: DataOrg, index?: number): DataOrg;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataOrg.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataOrg): ImportDataOrg.AsObject;
  static serializeBinaryToWriter(message: ImportDataOrg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataOrg;
  static deserializeBinaryFromReader(message: ImportDataOrg, reader: jspb.BinaryReader): ImportDataOrg;
}

export namespace ImportDataOrg {
  export type AsObject = {
    orgsList: Array<DataOrg.AsObject>,
  }
}

export class DataOrg extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): DataOrg;

  getOrg(): zitadel_management_pb.AddOrgRequest | undefined;
  setOrg(value?: zitadel_management_pb.AddOrgRequest): DataOrg;
  hasOrg(): boolean;
  clearOrg(): DataOrg;

  getDomainPolicy(): AddCustomDomainPolicyRequest | undefined;
  setDomainPolicy(value?: AddCustomDomainPolicyRequest): DataOrg;
  hasDomainPolicy(): boolean;
  clearDomainPolicy(): DataOrg;

  getLabelPolicy(): zitadel_management_pb.AddCustomLabelPolicyRequest | undefined;
  setLabelPolicy(value?: zitadel_management_pb.AddCustomLabelPolicyRequest): DataOrg;
  hasLabelPolicy(): boolean;
  clearLabelPolicy(): DataOrg;

  getLockoutPolicy(): zitadel_management_pb.AddCustomLockoutPolicyRequest | undefined;
  setLockoutPolicy(value?: zitadel_management_pb.AddCustomLockoutPolicyRequest): DataOrg;
  hasLockoutPolicy(): boolean;
  clearLockoutPolicy(): DataOrg;

  getLoginPolicy(): zitadel_management_pb.AddCustomLoginPolicyRequest | undefined;
  setLoginPolicy(value?: zitadel_management_pb.AddCustomLoginPolicyRequest): DataOrg;
  hasLoginPolicy(): boolean;
  clearLoginPolicy(): DataOrg;

  getPasswordComplexityPolicy(): zitadel_management_pb.AddCustomPasswordComplexityPolicyRequest | undefined;
  setPasswordComplexityPolicy(value?: zitadel_management_pb.AddCustomPasswordComplexityPolicyRequest): DataOrg;
  hasPasswordComplexityPolicy(): boolean;
  clearPasswordComplexityPolicy(): DataOrg;

  getPrivacyPolicy(): zitadel_management_pb.AddCustomPrivacyPolicyRequest | undefined;
  setPrivacyPolicy(value?: zitadel_management_pb.AddCustomPrivacyPolicyRequest): DataOrg;
  hasPrivacyPolicy(): boolean;
  clearPrivacyPolicy(): DataOrg;

  getProjectsList(): Array<zitadel_v1_pb.DataProject>;
  setProjectsList(value: Array<zitadel_v1_pb.DataProject>): DataOrg;
  clearProjectsList(): DataOrg;
  addProjects(value?: zitadel_v1_pb.DataProject, index?: number): zitadel_v1_pb.DataProject;

  getProjectRolesList(): Array<zitadel_management_pb.AddProjectRoleRequest>;
  setProjectRolesList(value: Array<zitadel_management_pb.AddProjectRoleRequest>): DataOrg;
  clearProjectRolesList(): DataOrg;
  addProjectRoles(value?: zitadel_management_pb.AddProjectRoleRequest, index?: number): zitadel_management_pb.AddProjectRoleRequest;

  getApiAppsList(): Array<zitadel_v1_pb.DataAPIApplication>;
  setApiAppsList(value: Array<zitadel_v1_pb.DataAPIApplication>): DataOrg;
  clearApiAppsList(): DataOrg;
  addApiApps(value?: zitadel_v1_pb.DataAPIApplication, index?: number): zitadel_v1_pb.DataAPIApplication;

  getOidcAppsList(): Array<zitadel_v1_pb.DataOIDCApplication>;
  setOidcAppsList(value: Array<zitadel_v1_pb.DataOIDCApplication>): DataOrg;
  clearOidcAppsList(): DataOrg;
  addOidcApps(value?: zitadel_v1_pb.DataOIDCApplication, index?: number): zitadel_v1_pb.DataOIDCApplication;

  getHumanUsersList(): Array<zitadel_v1_pb.DataHumanUser>;
  setHumanUsersList(value: Array<zitadel_v1_pb.DataHumanUser>): DataOrg;
  clearHumanUsersList(): DataOrg;
  addHumanUsers(value?: zitadel_v1_pb.DataHumanUser, index?: number): zitadel_v1_pb.DataHumanUser;

  getMachineUsersList(): Array<zitadel_v1_pb.DataMachineUser>;
  setMachineUsersList(value: Array<zitadel_v1_pb.DataMachineUser>): DataOrg;
  clearMachineUsersList(): DataOrg;
  addMachineUsers(value?: zitadel_v1_pb.DataMachineUser, index?: number): zitadel_v1_pb.DataMachineUser;

  getTriggerActionsList(): Array<zitadel_management_pb.SetTriggerActionsRequest>;
  setTriggerActionsList(value: Array<zitadel_management_pb.SetTriggerActionsRequest>): DataOrg;
  clearTriggerActionsList(): DataOrg;
  addTriggerActions(value?: zitadel_management_pb.SetTriggerActionsRequest, index?: number): zitadel_management_pb.SetTriggerActionsRequest;

  getActionsList(): Array<zitadel_v1_pb.DataAction>;
  setActionsList(value: Array<zitadel_v1_pb.DataAction>): DataOrg;
  clearActionsList(): DataOrg;
  addActions(value?: zitadel_v1_pb.DataAction, index?: number): zitadel_v1_pb.DataAction;

  getProjectGrantsList(): Array<zitadel_v1_pb.DataProjectGrant>;
  setProjectGrantsList(value: Array<zitadel_v1_pb.DataProjectGrant>): DataOrg;
  clearProjectGrantsList(): DataOrg;
  addProjectGrants(value?: zitadel_v1_pb.DataProjectGrant, index?: number): zitadel_v1_pb.DataProjectGrant;

  getUserGrantsList(): Array<zitadel_management_pb.AddUserGrantRequest>;
  setUserGrantsList(value: Array<zitadel_management_pb.AddUserGrantRequest>): DataOrg;
  clearUserGrantsList(): DataOrg;
  addUserGrants(value?: zitadel_management_pb.AddUserGrantRequest, index?: number): zitadel_management_pb.AddUserGrantRequest;

  getOrgMembersList(): Array<zitadel_management_pb.AddOrgMemberRequest>;
  setOrgMembersList(value: Array<zitadel_management_pb.AddOrgMemberRequest>): DataOrg;
  clearOrgMembersList(): DataOrg;
  addOrgMembers(value?: zitadel_management_pb.AddOrgMemberRequest, index?: number): zitadel_management_pb.AddOrgMemberRequest;

  getProjectMembersList(): Array<zitadel_management_pb.AddProjectMemberRequest>;
  setProjectMembersList(value: Array<zitadel_management_pb.AddProjectMemberRequest>): DataOrg;
  clearProjectMembersList(): DataOrg;
  addProjectMembers(value?: zitadel_management_pb.AddProjectMemberRequest, index?: number): zitadel_management_pb.AddProjectMemberRequest;

  getProjectGrantMembersList(): Array<zitadel_management_pb.AddProjectGrantMemberRequest>;
  setProjectGrantMembersList(value: Array<zitadel_management_pb.AddProjectGrantMemberRequest>): DataOrg;
  clearProjectGrantMembersList(): DataOrg;
  addProjectGrantMembers(value?: zitadel_management_pb.AddProjectGrantMemberRequest, index?: number): zitadel_management_pb.AddProjectGrantMemberRequest;

  getUserMetadataList(): Array<zitadel_management_pb.SetUserMetadataRequest>;
  setUserMetadataList(value: Array<zitadel_management_pb.SetUserMetadataRequest>): DataOrg;
  clearUserMetadataList(): DataOrg;
  addUserMetadata(value?: zitadel_management_pb.SetUserMetadataRequest, index?: number): zitadel_management_pb.SetUserMetadataRequest;

  getLoginTextsList(): Array<zitadel_management_pb.SetCustomLoginTextsRequest>;
  setLoginTextsList(value: Array<zitadel_management_pb.SetCustomLoginTextsRequest>): DataOrg;
  clearLoginTextsList(): DataOrg;
  addLoginTexts(value?: zitadel_management_pb.SetCustomLoginTextsRequest, index?: number): zitadel_management_pb.SetCustomLoginTextsRequest;

  getInitMessagesList(): Array<zitadel_management_pb.SetCustomInitMessageTextRequest>;
  setInitMessagesList(value: Array<zitadel_management_pb.SetCustomInitMessageTextRequest>): DataOrg;
  clearInitMessagesList(): DataOrg;
  addInitMessages(value?: zitadel_management_pb.SetCustomInitMessageTextRequest, index?: number): zitadel_management_pb.SetCustomInitMessageTextRequest;

  getPasswordResetMessagesList(): Array<zitadel_management_pb.SetCustomPasswordResetMessageTextRequest>;
  setPasswordResetMessagesList(value: Array<zitadel_management_pb.SetCustomPasswordResetMessageTextRequest>): DataOrg;
  clearPasswordResetMessagesList(): DataOrg;
  addPasswordResetMessages(value?: zitadel_management_pb.SetCustomPasswordResetMessageTextRequest, index?: number): zitadel_management_pb.SetCustomPasswordResetMessageTextRequest;

  getVerifyEmailMessagesList(): Array<zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest>;
  setVerifyEmailMessagesList(value: Array<zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest>): DataOrg;
  clearVerifyEmailMessagesList(): DataOrg;
  addVerifyEmailMessages(value?: zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest, index?: number): zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest;

  getVerifyPhoneMessagesList(): Array<zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest>;
  setVerifyPhoneMessagesList(value: Array<zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest>): DataOrg;
  clearVerifyPhoneMessagesList(): DataOrg;
  addVerifyPhoneMessages(value?: zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest, index?: number): zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest;

  getDomainClaimedMessagesList(): Array<zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest>;
  setDomainClaimedMessagesList(value: Array<zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest>): DataOrg;
  clearDomainClaimedMessagesList(): DataOrg;
  addDomainClaimedMessages(value?: zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest, index?: number): zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest;

  getPasswordlessRegistrationMessagesList(): Array<zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest>;
  setPasswordlessRegistrationMessagesList(value: Array<zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest>): DataOrg;
  clearPasswordlessRegistrationMessagesList(): DataOrg;
  addPasswordlessRegistrationMessages(value?: zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest, index?: number): zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest;

  getOidcIdpsList(): Array<zitadel_v1_pb.DataOIDCIDP>;
  setOidcIdpsList(value: Array<zitadel_v1_pb.DataOIDCIDP>): DataOrg;
  clearOidcIdpsList(): DataOrg;
  addOidcIdps(value?: zitadel_v1_pb.DataOIDCIDP, index?: number): zitadel_v1_pb.DataOIDCIDP;

  getJwtIdpsList(): Array<zitadel_v1_pb.DataJWTIDP>;
  setJwtIdpsList(value: Array<zitadel_v1_pb.DataJWTIDP>): DataOrg;
  clearJwtIdpsList(): DataOrg;
  addJwtIdps(value?: zitadel_v1_pb.DataJWTIDP, index?: number): zitadel_v1_pb.DataJWTIDP;

  getUserLinksList(): Array<zitadel_idp_pb.IDPUserLink>;
  setUserLinksList(value: Array<zitadel_idp_pb.IDPUserLink>): DataOrg;
  clearUserLinksList(): DataOrg;
  addUserLinks(value?: zitadel_idp_pb.IDPUserLink, index?: number): zitadel_idp_pb.IDPUserLink;

  getDomainsList(): Array<zitadel_org_pb.Domain>;
  setDomainsList(value: Array<zitadel_org_pb.Domain>): DataOrg;
  clearDomainsList(): DataOrg;
  addDomains(value?: zitadel_org_pb.Domain, index?: number): zitadel_org_pb.Domain;

  getAppKeysList(): Array<zitadel_v1_pb.DataAppKey>;
  setAppKeysList(value: Array<zitadel_v1_pb.DataAppKey>): DataOrg;
  clearAppKeysList(): DataOrg;
  addAppKeys(value?: zitadel_v1_pb.DataAppKey, index?: number): zitadel_v1_pb.DataAppKey;

  getMachineKeysList(): Array<zitadel_v1_pb.DataMachineKey>;
  setMachineKeysList(value: Array<zitadel_v1_pb.DataMachineKey>): DataOrg;
  clearMachineKeysList(): DataOrg;
  addMachineKeys(value?: zitadel_v1_pb.DataMachineKey, index?: number): zitadel_v1_pb.DataMachineKey;

  getVerifySmsOtpMessagesList(): Array<zitadel_management_pb.SetCustomVerifySMSOTPMessageTextRequest>;
  setVerifySmsOtpMessagesList(value: Array<zitadel_management_pb.SetCustomVerifySMSOTPMessageTextRequest>): DataOrg;
  clearVerifySmsOtpMessagesList(): DataOrg;
  addVerifySmsOtpMessages(value?: zitadel_management_pb.SetCustomVerifySMSOTPMessageTextRequest, index?: number): zitadel_management_pb.SetCustomVerifySMSOTPMessageTextRequest;

  getVerifyEmailOtpMessagesList(): Array<zitadel_management_pb.SetCustomVerifyEmailOTPMessageTextRequest>;
  setVerifyEmailOtpMessagesList(value: Array<zitadel_management_pb.SetCustomVerifyEmailOTPMessageTextRequest>): DataOrg;
  clearVerifyEmailOtpMessagesList(): DataOrg;
  addVerifyEmailOtpMessages(value?: zitadel_management_pb.SetCustomVerifyEmailOTPMessageTextRequest, index?: number): zitadel_management_pb.SetCustomVerifyEmailOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataOrg.AsObject;
  static toObject(includeInstance: boolean, msg: DataOrg): DataOrg.AsObject;
  static serializeBinaryToWriter(message: DataOrg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataOrg;
  static deserializeBinaryFromReader(message: DataOrg, reader: jspb.BinaryReader): DataOrg;
}

export namespace DataOrg {
  export type AsObject = {
    orgId: string,
    org?: zitadel_management_pb.AddOrgRequest.AsObject,
    domainPolicy?: AddCustomDomainPolicyRequest.AsObject,
    labelPolicy?: zitadel_management_pb.AddCustomLabelPolicyRequest.AsObject,
    lockoutPolicy?: zitadel_management_pb.AddCustomLockoutPolicyRequest.AsObject,
    loginPolicy?: zitadel_management_pb.AddCustomLoginPolicyRequest.AsObject,
    passwordComplexityPolicy?: zitadel_management_pb.AddCustomPasswordComplexityPolicyRequest.AsObject,
    privacyPolicy?: zitadel_management_pb.AddCustomPrivacyPolicyRequest.AsObject,
    projectsList: Array<zitadel_v1_pb.DataProject.AsObject>,
    projectRolesList: Array<zitadel_management_pb.AddProjectRoleRequest.AsObject>,
    apiAppsList: Array<zitadel_v1_pb.DataAPIApplication.AsObject>,
    oidcAppsList: Array<zitadel_v1_pb.DataOIDCApplication.AsObject>,
    humanUsersList: Array<zitadel_v1_pb.DataHumanUser.AsObject>,
    machineUsersList: Array<zitadel_v1_pb.DataMachineUser.AsObject>,
    triggerActionsList: Array<zitadel_management_pb.SetTriggerActionsRequest.AsObject>,
    actionsList: Array<zitadel_v1_pb.DataAction.AsObject>,
    projectGrantsList: Array<zitadel_v1_pb.DataProjectGrant.AsObject>,
    userGrantsList: Array<zitadel_management_pb.AddUserGrantRequest.AsObject>,
    orgMembersList: Array<zitadel_management_pb.AddOrgMemberRequest.AsObject>,
    projectMembersList: Array<zitadel_management_pb.AddProjectMemberRequest.AsObject>,
    projectGrantMembersList: Array<zitadel_management_pb.AddProjectGrantMemberRequest.AsObject>,
    userMetadataList: Array<zitadel_management_pb.SetUserMetadataRequest.AsObject>,
    loginTextsList: Array<zitadel_management_pb.SetCustomLoginTextsRequest.AsObject>,
    initMessagesList: Array<zitadel_management_pb.SetCustomInitMessageTextRequest.AsObject>,
    passwordResetMessagesList: Array<zitadel_management_pb.SetCustomPasswordResetMessageTextRequest.AsObject>,
    verifyEmailMessagesList: Array<zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest.AsObject>,
    verifyPhoneMessagesList: Array<zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest.AsObject>,
    domainClaimedMessagesList: Array<zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest.AsObject>,
    passwordlessRegistrationMessagesList: Array<zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest.AsObject>,
    oidcIdpsList: Array<zitadel_v1_pb.DataOIDCIDP.AsObject>,
    jwtIdpsList: Array<zitadel_v1_pb.DataJWTIDP.AsObject>,
    userLinksList: Array<zitadel_idp_pb.IDPUserLink.AsObject>,
    domainsList: Array<zitadel_org_pb.Domain.AsObject>,
    appKeysList: Array<zitadel_v1_pb.DataAppKey.AsObject>,
    machineKeysList: Array<zitadel_v1_pb.DataMachineKey.AsObject>,
    verifySmsOtpMessagesList: Array<zitadel_management_pb.SetCustomVerifySMSOTPMessageTextRequest.AsObject>,
    verifyEmailOtpMessagesList: Array<zitadel_management_pb.SetCustomVerifyEmailOTPMessageTextRequest.AsObject>,
  }
}

export class ImportDataResponse extends jspb.Message {
  getErrorsList(): Array<ImportDataError>;
  setErrorsList(value: Array<ImportDataError>): ImportDataResponse;
  clearErrorsList(): ImportDataResponse;
  addErrors(value?: ImportDataError, index?: number): ImportDataError;

  getSuccess(): ImportDataSuccess | undefined;
  setSuccess(value?: ImportDataSuccess): ImportDataResponse;
  hasSuccess(): boolean;
  clearSuccess(): ImportDataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataResponse): ImportDataResponse.AsObject;
  static serializeBinaryToWriter(message: ImportDataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataResponse;
  static deserializeBinaryFromReader(message: ImportDataResponse, reader: jspb.BinaryReader): ImportDataResponse;
}

export namespace ImportDataResponse {
  export type AsObject = {
    errorsList: Array<ImportDataError.AsObject>,
    success?: ImportDataSuccess.AsObject,
  }
}

export class ImportDataError extends jspb.Message {
  getType(): string;
  setType(value: string): ImportDataError;

  getId(): string;
  setId(value: string): ImportDataError;

  getMessage(): string;
  setMessage(value: string): ImportDataError;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataError.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataError): ImportDataError.AsObject;
  static serializeBinaryToWriter(message: ImportDataError, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataError;
  static deserializeBinaryFromReader(message: ImportDataError, reader: jspb.BinaryReader): ImportDataError;
}

export namespace ImportDataError {
  export type AsObject = {
    type: string,
    id: string,
    message: string,
  }
}

export class ImportDataSuccess extends jspb.Message {
  getOrgsList(): Array<ImportDataSuccessOrg>;
  setOrgsList(value: Array<ImportDataSuccessOrg>): ImportDataSuccess;
  clearOrgsList(): ImportDataSuccess;
  addOrgs(value?: ImportDataSuccessOrg, index?: number): ImportDataSuccessOrg;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataSuccess.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataSuccess): ImportDataSuccess.AsObject;
  static serializeBinaryToWriter(message: ImportDataSuccess, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataSuccess;
  static deserializeBinaryFromReader(message: ImportDataSuccess, reader: jspb.BinaryReader): ImportDataSuccess;
}

export namespace ImportDataSuccess {
  export type AsObject = {
    orgsList: Array<ImportDataSuccessOrg.AsObject>,
  }
}

export class ImportDataSuccessOrg extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): ImportDataSuccessOrg;

  getProjectIdsList(): Array<string>;
  setProjectIdsList(value: Array<string>): ImportDataSuccessOrg;
  clearProjectIdsList(): ImportDataSuccessOrg;
  addProjectIds(value: string, index?: number): ImportDataSuccessOrg;

  getProjectRolesList(): Array<string>;
  setProjectRolesList(value: Array<string>): ImportDataSuccessOrg;
  clearProjectRolesList(): ImportDataSuccessOrg;
  addProjectRoles(value: string, index?: number): ImportDataSuccessOrg;

  getOidcAppIdsList(): Array<string>;
  setOidcAppIdsList(value: Array<string>): ImportDataSuccessOrg;
  clearOidcAppIdsList(): ImportDataSuccessOrg;
  addOidcAppIds(value: string, index?: number): ImportDataSuccessOrg;

  getApiAppIdsList(): Array<string>;
  setApiAppIdsList(value: Array<string>): ImportDataSuccessOrg;
  clearApiAppIdsList(): ImportDataSuccessOrg;
  addApiAppIds(value: string, index?: number): ImportDataSuccessOrg;

  getHumanUserIdsList(): Array<string>;
  setHumanUserIdsList(value: Array<string>): ImportDataSuccessOrg;
  clearHumanUserIdsList(): ImportDataSuccessOrg;
  addHumanUserIds(value: string, index?: number): ImportDataSuccessOrg;

  getMachineUserIdsList(): Array<string>;
  setMachineUserIdsList(value: Array<string>): ImportDataSuccessOrg;
  clearMachineUserIdsList(): ImportDataSuccessOrg;
  addMachineUserIds(value: string, index?: number): ImportDataSuccessOrg;

  getActionIdsList(): Array<string>;
  setActionIdsList(value: Array<string>): ImportDataSuccessOrg;
  clearActionIdsList(): ImportDataSuccessOrg;
  addActionIds(value: string, index?: number): ImportDataSuccessOrg;

  getTriggerActionsList(): Array<zitadel_management_pb.SetTriggerActionsRequest>;
  setTriggerActionsList(value: Array<zitadel_management_pb.SetTriggerActionsRequest>): ImportDataSuccessOrg;
  clearTriggerActionsList(): ImportDataSuccessOrg;
  addTriggerActions(value?: zitadel_management_pb.SetTriggerActionsRequest, index?: number): zitadel_management_pb.SetTriggerActionsRequest;

  getProjectGrantsList(): Array<ImportDataSuccessProjectGrant>;
  setProjectGrantsList(value: Array<ImportDataSuccessProjectGrant>): ImportDataSuccessOrg;
  clearProjectGrantsList(): ImportDataSuccessOrg;
  addProjectGrants(value?: ImportDataSuccessProjectGrant, index?: number): ImportDataSuccessProjectGrant;

  getUserGrantsList(): Array<ImportDataSuccessUserGrant>;
  setUserGrantsList(value: Array<ImportDataSuccessUserGrant>): ImportDataSuccessOrg;
  clearUserGrantsList(): ImportDataSuccessOrg;
  addUserGrants(value?: ImportDataSuccessUserGrant, index?: number): ImportDataSuccessUserGrant;

  getOrgMembersList(): Array<string>;
  setOrgMembersList(value: Array<string>): ImportDataSuccessOrg;
  clearOrgMembersList(): ImportDataSuccessOrg;
  addOrgMembers(value: string, index?: number): ImportDataSuccessOrg;

  getProjectMembersList(): Array<ImportDataSuccessProjectMember>;
  setProjectMembersList(value: Array<ImportDataSuccessProjectMember>): ImportDataSuccessOrg;
  clearProjectMembersList(): ImportDataSuccessOrg;
  addProjectMembers(value?: ImportDataSuccessProjectMember, index?: number): ImportDataSuccessProjectMember;

  getProjectGrantMembersList(): Array<ImportDataSuccessProjectGrantMember>;
  setProjectGrantMembersList(value: Array<ImportDataSuccessProjectGrantMember>): ImportDataSuccessOrg;
  clearProjectGrantMembersList(): ImportDataSuccessOrg;
  addProjectGrantMembers(value?: ImportDataSuccessProjectGrantMember, index?: number): ImportDataSuccessProjectGrantMember;

  getOidcIpdsList(): Array<string>;
  setOidcIpdsList(value: Array<string>): ImportDataSuccessOrg;
  clearOidcIpdsList(): ImportDataSuccessOrg;
  addOidcIpds(value: string, index?: number): ImportDataSuccessOrg;

  getJwtIdpsList(): Array<string>;
  setJwtIdpsList(value: Array<string>): ImportDataSuccessOrg;
  clearJwtIdpsList(): ImportDataSuccessOrg;
  addJwtIdps(value: string, index?: number): ImportDataSuccessOrg;

  getIdpLinksList(): Array<string>;
  setIdpLinksList(value: Array<string>): ImportDataSuccessOrg;
  clearIdpLinksList(): ImportDataSuccessOrg;
  addIdpLinks(value: string, index?: number): ImportDataSuccessOrg;

  getUserLinksList(): Array<ImportDataSuccessUserLinks>;
  setUserLinksList(value: Array<ImportDataSuccessUserLinks>): ImportDataSuccessOrg;
  clearUserLinksList(): ImportDataSuccessOrg;
  addUserLinks(value?: ImportDataSuccessUserLinks, index?: number): ImportDataSuccessUserLinks;

  getUserMetadataList(): Array<ImportDataSuccessUserMetadata>;
  setUserMetadataList(value: Array<ImportDataSuccessUserMetadata>): ImportDataSuccessOrg;
  clearUserMetadataList(): ImportDataSuccessOrg;
  addUserMetadata(value?: ImportDataSuccessUserMetadata, index?: number): ImportDataSuccessUserMetadata;

  getDomainsList(): Array<string>;
  setDomainsList(value: Array<string>): ImportDataSuccessOrg;
  clearDomainsList(): ImportDataSuccessOrg;
  addDomains(value: string, index?: number): ImportDataSuccessOrg;

  getAppKeysList(): Array<string>;
  setAppKeysList(value: Array<string>): ImportDataSuccessOrg;
  clearAppKeysList(): ImportDataSuccessOrg;
  addAppKeys(value: string, index?: number): ImportDataSuccessOrg;

  getMachineKeysList(): Array<string>;
  setMachineKeysList(value: Array<string>): ImportDataSuccessOrg;
  clearMachineKeysList(): ImportDataSuccessOrg;
  addMachineKeys(value: string, index?: number): ImportDataSuccessOrg;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataSuccessOrg.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataSuccessOrg): ImportDataSuccessOrg.AsObject;
  static serializeBinaryToWriter(message: ImportDataSuccessOrg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataSuccessOrg;
  static deserializeBinaryFromReader(message: ImportDataSuccessOrg, reader: jspb.BinaryReader): ImportDataSuccessOrg;
}

export namespace ImportDataSuccessOrg {
  export type AsObject = {
    orgId: string,
    projectIdsList: Array<string>,
    projectRolesList: Array<string>,
    oidcAppIdsList: Array<string>,
    apiAppIdsList: Array<string>,
    humanUserIdsList: Array<string>,
    machineUserIdsList: Array<string>,
    actionIdsList: Array<string>,
    triggerActionsList: Array<zitadel_management_pb.SetTriggerActionsRequest.AsObject>,
    projectGrantsList: Array<ImportDataSuccessProjectGrant.AsObject>,
    userGrantsList: Array<ImportDataSuccessUserGrant.AsObject>,
    orgMembersList: Array<string>,
    projectMembersList: Array<ImportDataSuccessProjectMember.AsObject>,
    projectGrantMembersList: Array<ImportDataSuccessProjectGrantMember.AsObject>,
    oidcIpdsList: Array<string>,
    jwtIdpsList: Array<string>,
    idpLinksList: Array<string>,
    userLinksList: Array<ImportDataSuccessUserLinks.AsObject>,
    userMetadataList: Array<ImportDataSuccessUserMetadata.AsObject>,
    domainsList: Array<string>,
    appKeysList: Array<string>,
    machineKeysList: Array<string>,
  }
}

export class ImportDataSuccessProjectGrant extends jspb.Message {
  getGrantId(): string;
  setGrantId(value: string): ImportDataSuccessProjectGrant;

  getProjectId(): string;
  setProjectId(value: string): ImportDataSuccessProjectGrant;

  getOrgId(): string;
  setOrgId(value: string): ImportDataSuccessProjectGrant;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataSuccessProjectGrant.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataSuccessProjectGrant): ImportDataSuccessProjectGrant.AsObject;
  static serializeBinaryToWriter(message: ImportDataSuccessProjectGrant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataSuccessProjectGrant;
  static deserializeBinaryFromReader(message: ImportDataSuccessProjectGrant, reader: jspb.BinaryReader): ImportDataSuccessProjectGrant;
}

export namespace ImportDataSuccessProjectGrant {
  export type AsObject = {
    grantId: string,
    projectId: string,
    orgId: string,
  }
}

export class ImportDataSuccessUserGrant extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ImportDataSuccessUserGrant;

  getUserId(): string;
  setUserId(value: string): ImportDataSuccessUserGrant;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataSuccessUserGrant.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataSuccessUserGrant): ImportDataSuccessUserGrant.AsObject;
  static serializeBinaryToWriter(message: ImportDataSuccessUserGrant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataSuccessUserGrant;
  static deserializeBinaryFromReader(message: ImportDataSuccessUserGrant, reader: jspb.BinaryReader): ImportDataSuccessUserGrant;
}

export namespace ImportDataSuccessUserGrant {
  export type AsObject = {
    projectId: string,
    userId: string,
  }
}

export class ImportDataSuccessProjectMember extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ImportDataSuccessProjectMember;

  getUserId(): string;
  setUserId(value: string): ImportDataSuccessProjectMember;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataSuccessProjectMember.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataSuccessProjectMember): ImportDataSuccessProjectMember.AsObject;
  static serializeBinaryToWriter(message: ImportDataSuccessProjectMember, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataSuccessProjectMember;
  static deserializeBinaryFromReader(message: ImportDataSuccessProjectMember, reader: jspb.BinaryReader): ImportDataSuccessProjectMember;
}

export namespace ImportDataSuccessProjectMember {
  export type AsObject = {
    projectId: string,
    userId: string,
  }
}

export class ImportDataSuccessProjectGrantMember extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ImportDataSuccessProjectGrantMember;

  getGrantId(): string;
  setGrantId(value: string): ImportDataSuccessProjectGrantMember;

  getUserId(): string;
  setUserId(value: string): ImportDataSuccessProjectGrantMember;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataSuccessProjectGrantMember.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataSuccessProjectGrantMember): ImportDataSuccessProjectGrantMember.AsObject;
  static serializeBinaryToWriter(message: ImportDataSuccessProjectGrantMember, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataSuccessProjectGrantMember;
  static deserializeBinaryFromReader(message: ImportDataSuccessProjectGrantMember, reader: jspb.BinaryReader): ImportDataSuccessProjectGrantMember;
}

export namespace ImportDataSuccessProjectGrantMember {
  export type AsObject = {
    projectId: string,
    grantId: string,
    userId: string,
  }
}

export class ImportDataSuccessUserLinks extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ImportDataSuccessUserLinks;

  getExternalUserId(): string;
  setExternalUserId(value: string): ImportDataSuccessUserLinks;

  getDisplayName(): string;
  setDisplayName(value: string): ImportDataSuccessUserLinks;

  getIdpId(): string;
  setIdpId(value: string): ImportDataSuccessUserLinks;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataSuccessUserLinks.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataSuccessUserLinks): ImportDataSuccessUserLinks.AsObject;
  static serializeBinaryToWriter(message: ImportDataSuccessUserLinks, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataSuccessUserLinks;
  static deserializeBinaryFromReader(message: ImportDataSuccessUserLinks, reader: jspb.BinaryReader): ImportDataSuccessUserLinks;
}

export namespace ImportDataSuccessUserLinks {
  export type AsObject = {
    userId: string,
    externalUserId: string,
    displayName: string,
    idpId: string,
  }
}

export class ImportDataSuccessUserMetadata extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ImportDataSuccessUserMetadata;

  getKey(): string;
  setKey(value: string): ImportDataSuccessUserMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataSuccessUserMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataSuccessUserMetadata): ImportDataSuccessUserMetadata.AsObject;
  static serializeBinaryToWriter(message: ImportDataSuccessUserMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataSuccessUserMetadata;
  static deserializeBinaryFromReader(message: ImportDataSuccessUserMetadata, reader: jspb.BinaryReader): ImportDataSuccessUserMetadata;
}

export namespace ImportDataSuccessUserMetadata {
  export type AsObject = {
    userId: string,
    key: string,
  }
}

export class ExportDataRequest extends jspb.Message {
  getOrgIdsList(): Array<string>;
  setOrgIdsList(value: Array<string>): ExportDataRequest;
  clearOrgIdsList(): ExportDataRequest;
  addOrgIds(value: string, index?: number): ExportDataRequest;

  getExcludedOrgIdsList(): Array<string>;
  setExcludedOrgIdsList(value: Array<string>): ExportDataRequest;
  clearExcludedOrgIdsList(): ExportDataRequest;
  addExcludedOrgIds(value: string, index?: number): ExportDataRequest;

  getWithPasswords(): boolean;
  setWithPasswords(value: boolean): ExportDataRequest;

  getWithOtp(): boolean;
  setWithOtp(value: boolean): ExportDataRequest;

  getResponseOutput(): boolean;
  setResponseOutput(value: boolean): ExportDataRequest;

  getLocalOutput(): ExportDataRequest.LocalOutput | undefined;
  setLocalOutput(value?: ExportDataRequest.LocalOutput): ExportDataRequest;
  hasLocalOutput(): boolean;
  clearLocalOutput(): ExportDataRequest;

  getS3Output(): ExportDataRequest.S3Output | undefined;
  setS3Output(value?: ExportDataRequest.S3Output): ExportDataRequest;
  hasS3Output(): boolean;
  clearS3Output(): ExportDataRequest;

  getGcsOutput(): ExportDataRequest.GCSOutput | undefined;
  setGcsOutput(value?: ExportDataRequest.GCSOutput): ExportDataRequest;
  hasGcsOutput(): boolean;
  clearGcsOutput(): ExportDataRequest;

  getTimeout(): string;
  setTimeout(value: string): ExportDataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExportDataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExportDataRequest): ExportDataRequest.AsObject;
  static serializeBinaryToWriter(message: ExportDataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExportDataRequest;
  static deserializeBinaryFromReader(message: ExportDataRequest, reader: jspb.BinaryReader): ExportDataRequest;
}

export namespace ExportDataRequest {
  export type AsObject = {
    orgIdsList: Array<string>,
    excludedOrgIdsList: Array<string>,
    withPasswords: boolean,
    withOtp: boolean,
    responseOutput: boolean,
    localOutput?: ExportDataRequest.LocalOutput.AsObject,
    s3Output?: ExportDataRequest.S3Output.AsObject,
    gcsOutput?: ExportDataRequest.GCSOutput.AsObject,
    timeout: string,
  }

  export class LocalOutput extends jspb.Message {
    getPath(): string;
    setPath(value: string): LocalOutput;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LocalOutput.AsObject;
    static toObject(includeInstance: boolean, msg: LocalOutput): LocalOutput.AsObject;
    static serializeBinaryToWriter(message: LocalOutput, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LocalOutput;
    static deserializeBinaryFromReader(message: LocalOutput, reader: jspb.BinaryReader): LocalOutput;
  }

  export namespace LocalOutput {
    export type AsObject = {
      path: string,
    }
  }


  export class S3Output extends jspb.Message {
    getPath(): string;
    setPath(value: string): S3Output;

    getEndpoint(): string;
    setEndpoint(value: string): S3Output;

    getAccessKeyId(): string;
    setAccessKeyId(value: string): S3Output;

    getSecretAccessKey(): string;
    setSecretAccessKey(value: string): S3Output;

    getSsl(): boolean;
    setSsl(value: boolean): S3Output;

    getBucket(): string;
    setBucket(value: string): S3Output;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): S3Output.AsObject;
    static toObject(includeInstance: boolean, msg: S3Output): S3Output.AsObject;
    static serializeBinaryToWriter(message: S3Output, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): S3Output;
    static deserializeBinaryFromReader(message: S3Output, reader: jspb.BinaryReader): S3Output;
  }

  export namespace S3Output {
    export type AsObject = {
      path: string,
      endpoint: string,
      accessKeyId: string,
      secretAccessKey: string,
      ssl: boolean,
      bucket: string,
    }
  }


  export class GCSOutput extends jspb.Message {
    getBucket(): string;
    setBucket(value: string): GCSOutput;

    getServiceaccountJson(): string;
    setServiceaccountJson(value: string): GCSOutput;

    getPath(): string;
    setPath(value: string): GCSOutput;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GCSOutput.AsObject;
    static toObject(includeInstance: boolean, msg: GCSOutput): GCSOutput.AsObject;
    static serializeBinaryToWriter(message: GCSOutput, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GCSOutput;
    static deserializeBinaryFromReader(message: GCSOutput, reader: jspb.BinaryReader): GCSOutput;
  }

  export namespace GCSOutput {
    export type AsObject = {
      bucket: string,
      serviceaccountJson: string,
      path: string,
    }
  }

}

export class ExportDataResponse extends jspb.Message {
  getOrgsList(): Array<DataOrg>;
  setOrgsList(value: Array<DataOrg>): ExportDataResponse;
  clearOrgsList(): ExportDataResponse;
  addOrgs(value?: DataOrg, index?: number): DataOrg;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExportDataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExportDataResponse): ExportDataResponse.AsObject;
  static serializeBinaryToWriter(message: ExportDataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExportDataResponse;
  static deserializeBinaryFromReader(message: ExportDataResponse, reader: jspb.BinaryReader): ExportDataResponse;
}

export namespace ExportDataResponse {
  export type AsObject = {
    orgsList: Array<DataOrg.AsObject>,
  }
}

export class ListEventsRequest extends jspb.Message {
  getSequence(): number;
  setSequence(value: number): ListEventsRequest;

  getLimit(): number;
  setLimit(value: number): ListEventsRequest;

  getAsc(): boolean;
  setAsc(value: boolean): ListEventsRequest;

  getEditorUserId(): string;
  setEditorUserId(value: string): ListEventsRequest;

  getEventTypesList(): Array<string>;
  setEventTypesList(value: Array<string>): ListEventsRequest;
  clearEventTypesList(): ListEventsRequest;
  addEventTypes(value: string, index?: number): ListEventsRequest;

  getAggregateId(): string;
  setAggregateId(value: string): ListEventsRequest;

  getAggregateTypesList(): Array<string>;
  setAggregateTypesList(value: Array<string>): ListEventsRequest;
  clearAggregateTypesList(): ListEventsRequest;
  addAggregateTypes(value: string, index?: number): ListEventsRequest;

  getResourceOwner(): string;
  setResourceOwner(value: string): ListEventsRequest;

  getCreationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreationDate(value?: google_protobuf_timestamp_pb.Timestamp): ListEventsRequest;
  hasCreationDate(): boolean;
  clearCreationDate(): ListEventsRequest;

  getRange(): ListEventsRequest.creation_date_range | undefined;
  setRange(value?: ListEventsRequest.creation_date_range): ListEventsRequest;
  hasRange(): boolean;
  clearRange(): ListEventsRequest;

  getFrom(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setFrom(value?: google_protobuf_timestamp_pb.Timestamp): ListEventsRequest;
  hasFrom(): boolean;
  clearFrom(): ListEventsRequest;

  getCreationDateFilterCase(): ListEventsRequest.CreationDateFilterCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEventsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEventsRequest): ListEventsRequest.AsObject;
  static serializeBinaryToWriter(message: ListEventsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEventsRequest;
  static deserializeBinaryFromReader(message: ListEventsRequest, reader: jspb.BinaryReader): ListEventsRequest;
}

export namespace ListEventsRequest {
  export type AsObject = {
    sequence: number,
    limit: number,
    asc: boolean,
    editorUserId: string,
    eventTypesList: Array<string>,
    aggregateId: string,
    aggregateTypesList: Array<string>,
    resourceOwner: string,
    creationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    range?: ListEventsRequest.creation_date_range.AsObject,
    from?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }

  export class creation_date_range extends jspb.Message {
    getSince(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setSince(value?: google_protobuf_timestamp_pb.Timestamp): creation_date_range;
    hasSince(): boolean;
    clearSince(): creation_date_range;

    getUntil(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setUntil(value?: google_protobuf_timestamp_pb.Timestamp): creation_date_range;
    hasUntil(): boolean;
    clearUntil(): creation_date_range;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): creation_date_range.AsObject;
    static toObject(includeInstance: boolean, msg: creation_date_range): creation_date_range.AsObject;
    static serializeBinaryToWriter(message: creation_date_range, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): creation_date_range;
    static deserializeBinaryFromReader(message: creation_date_range, reader: jspb.BinaryReader): creation_date_range;
  }

  export namespace creation_date_range {
    export type AsObject = {
      since?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      until?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
  }


  export enum CreationDateFilterCase { 
    CREATION_DATE_FILTER_NOT_SET = 0,
    RANGE = 10,
    FROM = 11,
  }
}

export class ListEventsResponse extends jspb.Message {
  getEventsList(): Array<zitadel_event_pb.Event>;
  setEventsList(value: Array<zitadel_event_pb.Event>): ListEventsResponse;
  clearEventsList(): ListEventsResponse;
  addEvents(value?: zitadel_event_pb.Event, index?: number): zitadel_event_pb.Event;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEventsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEventsResponse): ListEventsResponse.AsObject;
  static serializeBinaryToWriter(message: ListEventsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEventsResponse;
  static deserializeBinaryFromReader(message: ListEventsResponse, reader: jspb.BinaryReader): ListEventsResponse;
}

export namespace ListEventsResponse {
  export type AsObject = {
    eventsList: Array<zitadel_event_pb.Event.AsObject>,
  }
}

export class ListEventTypesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEventTypesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEventTypesRequest): ListEventTypesRequest.AsObject;
  static serializeBinaryToWriter(message: ListEventTypesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEventTypesRequest;
  static deserializeBinaryFromReader(message: ListEventTypesRequest, reader: jspb.BinaryReader): ListEventTypesRequest;
}

export namespace ListEventTypesRequest {
  export type AsObject = {
  }
}

export class ListEventTypesResponse extends jspb.Message {
  getEventTypesList(): Array<zitadel_event_pb.EventType>;
  setEventTypesList(value: Array<zitadel_event_pb.EventType>): ListEventTypesResponse;
  clearEventTypesList(): ListEventTypesResponse;
  addEventTypes(value?: zitadel_event_pb.EventType, index?: number): zitadel_event_pb.EventType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEventTypesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEventTypesResponse): ListEventTypesResponse.AsObject;
  static serializeBinaryToWriter(message: ListEventTypesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEventTypesResponse;
  static deserializeBinaryFromReader(message: ListEventTypesResponse, reader: jspb.BinaryReader): ListEventTypesResponse;
}

export namespace ListEventTypesResponse {
  export type AsObject = {
    eventTypesList: Array<zitadel_event_pb.EventType.AsObject>,
  }
}

export class ListAggregateTypesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAggregateTypesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAggregateTypesRequest): ListAggregateTypesRequest.AsObject;
  static serializeBinaryToWriter(message: ListAggregateTypesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAggregateTypesRequest;
  static deserializeBinaryFromReader(message: ListAggregateTypesRequest, reader: jspb.BinaryReader): ListAggregateTypesRequest;
}

export namespace ListAggregateTypesRequest {
  export type AsObject = {
  }
}

export class ListAggregateTypesResponse extends jspb.Message {
  getAggregateTypesList(): Array<zitadel_event_pb.AggregateType>;
  setAggregateTypesList(value: Array<zitadel_event_pb.AggregateType>): ListAggregateTypesResponse;
  clearAggregateTypesList(): ListAggregateTypesResponse;
  addAggregateTypes(value?: zitadel_event_pb.AggregateType, index?: number): zitadel_event_pb.AggregateType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAggregateTypesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAggregateTypesResponse): ListAggregateTypesResponse.AsObject;
  static serializeBinaryToWriter(message: ListAggregateTypesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAggregateTypesResponse;
  static deserializeBinaryFromReader(message: ListAggregateTypesResponse, reader: jspb.BinaryReader): ListAggregateTypesResponse;
}

export namespace ListAggregateTypesResponse {
  export type AsObject = {
    aggregateTypesList: Array<zitadel_event_pb.AggregateType.AsObject>,
  }
}

export class ActivateFeatureLoginDefaultOrgRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateFeatureLoginDefaultOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateFeatureLoginDefaultOrgRequest): ActivateFeatureLoginDefaultOrgRequest.AsObject;
  static serializeBinaryToWriter(message: ActivateFeatureLoginDefaultOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateFeatureLoginDefaultOrgRequest;
  static deserializeBinaryFromReader(message: ActivateFeatureLoginDefaultOrgRequest, reader: jspb.BinaryReader): ActivateFeatureLoginDefaultOrgRequest;
}

export namespace ActivateFeatureLoginDefaultOrgRequest {
  export type AsObject = {
  }
}

export class ActivateFeatureLoginDefaultOrgResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ActivateFeatureLoginDefaultOrgResponse;
  hasDetails(): boolean;
  clearDetails(): ActivateFeatureLoginDefaultOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateFeatureLoginDefaultOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateFeatureLoginDefaultOrgResponse): ActivateFeatureLoginDefaultOrgResponse.AsObject;
  static serializeBinaryToWriter(message: ActivateFeatureLoginDefaultOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateFeatureLoginDefaultOrgResponse;
  static deserializeBinaryFromReader(message: ActivateFeatureLoginDefaultOrgResponse, reader: jspb.BinaryReader): ActivateFeatureLoginDefaultOrgResponse;
}

export namespace ActivateFeatureLoginDefaultOrgResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListMilestonesRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListMilestonesRequest;
  hasQuery(): boolean;
  clearQuery(): ListMilestonesRequest;

  getSortingColumn(): zitadel_milestone_v1_milestone_pb.MilestoneFieldName;
  setSortingColumn(value: zitadel_milestone_v1_milestone_pb.MilestoneFieldName): ListMilestonesRequest;

  getQueriesList(): Array<zitadel_milestone_v1_milestone_pb.MilestoneQuery>;
  setQueriesList(value: Array<zitadel_milestone_v1_milestone_pb.MilestoneQuery>): ListMilestonesRequest;
  clearQueriesList(): ListMilestonesRequest;
  addQueries(value?: zitadel_milestone_v1_milestone_pb.MilestoneQuery, index?: number): zitadel_milestone_v1_milestone_pb.MilestoneQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMilestonesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMilestonesRequest): ListMilestonesRequest.AsObject;
  static serializeBinaryToWriter(message: ListMilestonesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMilestonesRequest;
  static deserializeBinaryFromReader(message: ListMilestonesRequest, reader: jspb.BinaryReader): ListMilestonesRequest;
}

export namespace ListMilestonesRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_milestone_v1_milestone_pb.MilestoneFieldName,
    queriesList: Array<zitadel_milestone_v1_milestone_pb.MilestoneQuery.AsObject>,
  }
}

export class ListMilestonesResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListMilestonesResponse;
  hasDetails(): boolean;
  clearDetails(): ListMilestonesResponse;

  getResultList(): Array<zitadel_milestone_v1_milestone_pb.Milestone>;
  setResultList(value: Array<zitadel_milestone_v1_milestone_pb.Milestone>): ListMilestonesResponse;
  clearResultList(): ListMilestonesResponse;
  addResult(value?: zitadel_milestone_v1_milestone_pb.Milestone, index?: number): zitadel_milestone_v1_milestone_pb.Milestone;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMilestonesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMilestonesResponse): ListMilestonesResponse.AsObject;
  static serializeBinaryToWriter(message: ListMilestonesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMilestonesResponse;
  static deserializeBinaryFromReader(message: ListMilestonesResponse, reader: jspb.BinaryReader): ListMilestonesResponse;
}

export namespace ListMilestonesResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_milestone_v1_milestone_pb.Milestone.AsObject>,
  }
}

export class SetRestrictionsRequest extends jspb.Message {
  getDisallowPublicOrgRegistration(): boolean;
  setDisallowPublicOrgRegistration(value: boolean): SetRestrictionsRequest;
  hasDisallowPublicOrgRegistration(): boolean;
  clearDisallowPublicOrgRegistration(): SetRestrictionsRequest;

  getAllowedLanguages(): SelectLanguages | undefined;
  setAllowedLanguages(value?: SelectLanguages): SetRestrictionsRequest;
  hasAllowedLanguages(): boolean;
  clearAllowedLanguages(): SetRestrictionsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRestrictionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetRestrictionsRequest): SetRestrictionsRequest.AsObject;
  static serializeBinaryToWriter(message: SetRestrictionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRestrictionsRequest;
  static deserializeBinaryFromReader(message: SetRestrictionsRequest, reader: jspb.BinaryReader): SetRestrictionsRequest;
}

export namespace SetRestrictionsRequest {
  export type AsObject = {
    disallowPublicOrgRegistration?: boolean,
    allowedLanguages?: SelectLanguages.AsObject,
  }

  export enum DisallowPublicOrgRegistrationCase { 
    _DISALLOW_PUBLIC_ORG_REGISTRATION_NOT_SET = 0,
    DISALLOW_PUBLIC_ORG_REGISTRATION = 1,
  }

  export enum AllowedLanguagesCase { 
    _ALLOWED_LANGUAGES_NOT_SET = 0,
    ALLOWED_LANGUAGES = 2,
  }
}

export class SelectLanguages extends jspb.Message {
  getListList(): Array<string>;
  setListList(value: Array<string>): SelectLanguages;
  clearListList(): SelectLanguages;
  addList(value: string, index?: number): SelectLanguages;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SelectLanguages.AsObject;
  static toObject(includeInstance: boolean, msg: SelectLanguages): SelectLanguages.AsObject;
  static serializeBinaryToWriter(message: SelectLanguages, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SelectLanguages;
  static deserializeBinaryFromReader(message: SelectLanguages, reader: jspb.BinaryReader): SelectLanguages;
}

export namespace SelectLanguages {
  export type AsObject = {
    listList: Array<string>,
  }
}

export class SetRestrictionsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetRestrictionsResponse;
  hasDetails(): boolean;
  clearDetails(): SetRestrictionsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetRestrictionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetRestrictionsResponse): SetRestrictionsResponse.AsObject;
  static serializeBinaryToWriter(message: SetRestrictionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetRestrictionsResponse;
  static deserializeBinaryFromReader(message: SetRestrictionsResponse, reader: jspb.BinaryReader): SetRestrictionsResponse;
}

export namespace SetRestrictionsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetRestrictionsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRestrictionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetRestrictionsRequest): GetRestrictionsRequest.AsObject;
  static serializeBinaryToWriter(message: GetRestrictionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRestrictionsRequest;
  static deserializeBinaryFromReader(message: GetRestrictionsRequest, reader: jspb.BinaryReader): GetRestrictionsRequest;
}

export namespace GetRestrictionsRequest {
  export type AsObject = {
  }
}

export class GetRestrictionsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GetRestrictionsResponse;
  hasDetails(): boolean;
  clearDetails(): GetRestrictionsResponse;

  getDisallowPublicOrgRegistration(): boolean;
  setDisallowPublicOrgRegistration(value: boolean): GetRestrictionsResponse;

  getAllowedLanguagesList(): Array<string>;
  setAllowedLanguagesList(value: Array<string>): GetRestrictionsResponse;
  clearAllowedLanguagesList(): GetRestrictionsResponse;
  addAllowedLanguages(value: string, index?: number): GetRestrictionsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetRestrictionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetRestrictionsResponse): GetRestrictionsResponse.AsObject;
  static serializeBinaryToWriter(message: GetRestrictionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetRestrictionsResponse;
  static deserializeBinaryFromReader(message: GetRestrictionsResponse, reader: jspb.BinaryReader): GetRestrictionsResponse;
}

export namespace GetRestrictionsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    disallowPublicOrgRegistration: boolean,
    allowedLanguagesList: Array<string>,
  }
}


import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from '../../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_protoc_gen_zitadel_v2_options_pb from '../../../../zitadel/protoc_gen_zitadel/v2/options_pb'; // proto import: "zitadel/protoc_gen_zitadel/v2/options.proto"
import * as zitadel_user_schema_v3alpha_user_schema_pb from '../../../../zitadel/user/schema/v3alpha/user_schema_pb'; // proto import: "zitadel/user/schema/v3alpha/user_schema.proto"


export class ListUserSchemasRequest extends jspb.Message {
  getQuery(): zitadel_object_v2beta_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_v2beta_object_pb.ListQuery): ListUserSchemasRequest;
  hasQuery(): boolean;
  clearQuery(): ListUserSchemasRequest;

  getSortingColumn(): zitadel_user_schema_v3alpha_user_schema_pb.FieldName;
  setSortingColumn(value: zitadel_user_schema_v3alpha_user_schema_pb.FieldName): ListUserSchemasRequest;

  getQueriesList(): Array<zitadel_user_schema_v3alpha_user_schema_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_user_schema_v3alpha_user_schema_pb.SearchQuery>): ListUserSchemasRequest;
  clearQueriesList(): ListUserSchemasRequest;
  addQueries(value?: zitadel_user_schema_v3alpha_user_schema_pb.SearchQuery, index?: number): zitadel_user_schema_v3alpha_user_schema_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserSchemasRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserSchemasRequest): ListUserSchemasRequest.AsObject;
  static serializeBinaryToWriter(message: ListUserSchemasRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserSchemasRequest;
  static deserializeBinaryFromReader(message: ListUserSchemasRequest, reader: jspb.BinaryReader): ListUserSchemasRequest;
}

export namespace ListUserSchemasRequest {
  export type AsObject = {
    query?: zitadel_object_v2beta_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_user_schema_v3alpha_user_schema_pb.FieldName,
    queriesList: Array<zitadel_user_schema_v3alpha_user_schema_pb.SearchQuery.AsObject>,
  }
}

export class ListUserSchemasResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.ListDetails): ListUserSchemasResponse;
  hasDetails(): boolean;
  clearDetails(): ListUserSchemasResponse;

  getSortingColumn(): zitadel_user_schema_v3alpha_user_schema_pb.FieldName;
  setSortingColumn(value: zitadel_user_schema_v3alpha_user_schema_pb.FieldName): ListUserSchemasResponse;

  getResultList(): Array<zitadel_user_schema_v3alpha_user_schema_pb.UserSchema>;
  setResultList(value: Array<zitadel_user_schema_v3alpha_user_schema_pb.UserSchema>): ListUserSchemasResponse;
  clearResultList(): ListUserSchemasResponse;
  addResult(value?: zitadel_user_schema_v3alpha_user_schema_pb.UserSchema, index?: number): zitadel_user_schema_v3alpha_user_schema_pb.UserSchema;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserSchemasResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserSchemasResponse): ListUserSchemasResponse.AsObject;
  static serializeBinaryToWriter(message: ListUserSchemasResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserSchemasResponse;
  static deserializeBinaryFromReader(message: ListUserSchemasResponse, reader: jspb.BinaryReader): ListUserSchemasResponse;
}

export namespace ListUserSchemasResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_user_schema_v3alpha_user_schema_pb.FieldName,
    resultList: Array<zitadel_user_schema_v3alpha_user_schema_pb.UserSchema.AsObject>,
  }
}

export class GetUserSchemaByIDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetUserSchemaByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserSchemaByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserSchemaByIDRequest): GetUserSchemaByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserSchemaByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserSchemaByIDRequest;
  static deserializeBinaryFromReader(message: GetUserSchemaByIDRequest, reader: jspb.BinaryReader): GetUserSchemaByIDRequest;
}

export namespace GetUserSchemaByIDRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetUserSchemaByIDResponse extends jspb.Message {
  getSchema(): zitadel_user_schema_v3alpha_user_schema_pb.UserSchema | undefined;
  setSchema(value?: zitadel_user_schema_v3alpha_user_schema_pb.UserSchema): GetUserSchemaByIDResponse;
  hasSchema(): boolean;
  clearSchema(): GetUserSchemaByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserSchemaByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserSchemaByIDResponse): GetUserSchemaByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserSchemaByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserSchemaByIDResponse;
  static deserializeBinaryFromReader(message: GetUserSchemaByIDResponse, reader: jspb.BinaryReader): GetUserSchemaByIDResponse;
}

export namespace GetUserSchemaByIDResponse {
  export type AsObject = {
    schema?: zitadel_user_schema_v3alpha_user_schema_pb.UserSchema.AsObject,
  }
}

export class CreateUserSchemaRequest extends jspb.Message {
  getType(): string;
  setType(value: string): CreateUserSchemaRequest;

  getSchema(): google_protobuf_struct_pb.Struct | undefined;
  setSchema(value?: google_protobuf_struct_pb.Struct): CreateUserSchemaRequest;
  hasSchema(): boolean;
  clearSchema(): CreateUserSchemaRequest;

  getPossibleAuthenticatorsList(): Array<zitadel_user_schema_v3alpha_user_schema_pb.AuthenticatorType>;
  setPossibleAuthenticatorsList(value: Array<zitadel_user_schema_v3alpha_user_schema_pb.AuthenticatorType>): CreateUserSchemaRequest;
  clearPossibleAuthenticatorsList(): CreateUserSchemaRequest;
  addPossibleAuthenticators(value: zitadel_user_schema_v3alpha_user_schema_pb.AuthenticatorType, index?: number): CreateUserSchemaRequest;

  getDataTypeCase(): CreateUserSchemaRequest.DataTypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUserSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUserSchemaRequest): CreateUserSchemaRequest.AsObject;
  static serializeBinaryToWriter(message: CreateUserSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUserSchemaRequest;
  static deserializeBinaryFromReader(message: CreateUserSchemaRequest, reader: jspb.BinaryReader): CreateUserSchemaRequest;
}

export namespace CreateUserSchemaRequest {
  export type AsObject = {
    type: string,
    schema?: google_protobuf_struct_pb.Struct.AsObject,
    possibleAuthenticatorsList: Array<zitadel_user_schema_v3alpha_user_schema_pb.AuthenticatorType>,
  }

  export enum DataTypeCase { 
    DATA_TYPE_NOT_SET = 0,
    SCHEMA = 2,
  }
}

export class CreateUserSchemaResponse extends jspb.Message {
  getId(): string;
  setId(value: string): CreateUserSchemaResponse;

  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): CreateUserSchemaResponse;
  hasDetails(): boolean;
  clearDetails(): CreateUserSchemaResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateUserSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateUserSchemaResponse): CreateUserSchemaResponse.AsObject;
  static serializeBinaryToWriter(message: CreateUserSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateUserSchemaResponse;
  static deserializeBinaryFromReader(message: CreateUserSchemaResponse, reader: jspb.BinaryReader): CreateUserSchemaResponse;
}

export namespace CreateUserSchemaResponse {
  export type AsObject = {
    id: string,
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class UpdateUserSchemaRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateUserSchemaRequest;

  getType(): string;
  setType(value: string): UpdateUserSchemaRequest;
  hasType(): boolean;
  clearType(): UpdateUserSchemaRequest;

  getSchema(): google_protobuf_struct_pb.Struct | undefined;
  setSchema(value?: google_protobuf_struct_pb.Struct): UpdateUserSchemaRequest;
  hasSchema(): boolean;
  clearSchema(): UpdateUserSchemaRequest;

  getPossibleAuthenticatorsList(): Array<zitadel_user_schema_v3alpha_user_schema_pb.AuthenticatorType>;
  setPossibleAuthenticatorsList(value: Array<zitadel_user_schema_v3alpha_user_schema_pb.AuthenticatorType>): UpdateUserSchemaRequest;
  clearPossibleAuthenticatorsList(): UpdateUserSchemaRequest;
  addPossibleAuthenticators(value: zitadel_user_schema_v3alpha_user_schema_pb.AuthenticatorType, index?: number): UpdateUserSchemaRequest;

  getDataTypeCase(): UpdateUserSchemaRequest.DataTypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserSchemaRequest): UpdateUserSchemaRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserSchemaRequest;
  static deserializeBinaryFromReader(message: UpdateUserSchemaRequest, reader: jspb.BinaryReader): UpdateUserSchemaRequest;
}

export namespace UpdateUserSchemaRequest {
  export type AsObject = {
    id: string,
    type?: string,
    schema?: google_protobuf_struct_pb.Struct.AsObject,
    possibleAuthenticatorsList: Array<zitadel_user_schema_v3alpha_user_schema_pb.AuthenticatorType>,
  }

  export enum DataTypeCase { 
    DATA_TYPE_NOT_SET = 0,
    SCHEMA = 3,
  }

  export enum TypeCase { 
    _TYPE_NOT_SET = 0,
    TYPE = 2,
  }
}

export class UpdateUserSchemaResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): UpdateUserSchemaResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateUserSchemaResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserSchemaResponse): UpdateUserSchemaResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateUserSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserSchemaResponse;
  static deserializeBinaryFromReader(message: UpdateUserSchemaResponse, reader: jspb.BinaryReader): UpdateUserSchemaResponse;
}

export namespace UpdateUserSchemaResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class DeactivateUserSchemaRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeactivateUserSchemaRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateUserSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateUserSchemaRequest): DeactivateUserSchemaRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateUserSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateUserSchemaRequest;
  static deserializeBinaryFromReader(message: DeactivateUserSchemaRequest, reader: jspb.BinaryReader): DeactivateUserSchemaRequest;
}

export namespace DeactivateUserSchemaRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeactivateUserSchemaResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): DeactivateUserSchemaResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateUserSchemaResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateUserSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateUserSchemaResponse): DeactivateUserSchemaResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateUserSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateUserSchemaResponse;
  static deserializeBinaryFromReader(message: DeactivateUserSchemaResponse, reader: jspb.BinaryReader): DeactivateUserSchemaResponse;
}

export namespace DeactivateUserSchemaResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ReactivateUserSchemaRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ReactivateUserSchemaRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateUserSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateUserSchemaRequest): ReactivateUserSchemaRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateUserSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateUserSchemaRequest;
  static deserializeBinaryFromReader(message: ReactivateUserSchemaRequest, reader: jspb.BinaryReader): ReactivateUserSchemaRequest;
}

export namespace ReactivateUserSchemaRequest {
  export type AsObject = {
    id: string,
  }
}

export class ReactivateUserSchemaResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ReactivateUserSchemaResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateUserSchemaResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateUserSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateUserSchemaResponse): ReactivateUserSchemaResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateUserSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateUserSchemaResponse;
  static deserializeBinaryFromReader(message: ReactivateUserSchemaResponse, reader: jspb.BinaryReader): ReactivateUserSchemaResponse;
}

export namespace ReactivateUserSchemaResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class DeleteUserSchemaRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteUserSchemaRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUserSchemaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUserSchemaRequest): DeleteUserSchemaRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteUserSchemaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUserSchemaRequest;
  static deserializeBinaryFromReader(message: DeleteUserSchemaRequest, reader: jspb.BinaryReader): DeleteUserSchemaRequest;
}

export namespace DeleteUserSchemaRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeleteUserSchemaResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): DeleteUserSchemaResponse;
  hasDetails(): boolean;
  clearDetails(): DeleteUserSchemaResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteUserSchemaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteUserSchemaResponse): DeleteUserSchemaResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteUserSchemaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteUserSchemaResponse;
  static deserializeBinaryFromReader(message: DeleteUserSchemaResponse, reader: jspb.BinaryReader): DeleteUserSchemaResponse;
}

export namespace DeleteUserSchemaResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}


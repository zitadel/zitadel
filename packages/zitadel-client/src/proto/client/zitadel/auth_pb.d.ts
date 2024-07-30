import * as jspb from 'google-protobuf'

import * as zitadel_user_pb from '../zitadel/user_pb'; // proto import: "zitadel/user.proto"
import * as zitadel_org_pb from '../zitadel/org_pb'; // proto import: "zitadel/org.proto"
import * as zitadel_change_pb from '../zitadel/change_pb'; // proto import: "zitadel/change.proto"
import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as zitadel_options_pb from '../zitadel/options_pb'; // proto import: "zitadel/options.proto"
import * as zitadel_policy_pb from '../zitadel/policy_pb'; // proto import: "zitadel/policy.proto"
import * as zitadel_idp_pb from '../zitadel/idp_pb'; // proto import: "zitadel/idp.proto"
import * as zitadel_metadata_pb from '../zitadel/metadata_pb'; // proto import: "zitadel/metadata.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as google_api_annotations_pb from '../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


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

export class GetMyUserRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyUserRequest): GetMyUserRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyUserRequest;
  static deserializeBinaryFromReader(message: GetMyUserRequest, reader: jspb.BinaryReader): GetMyUserRequest;
}

export namespace GetMyUserRequest {
  export type AsObject = {
  }
}

export class GetMyUserResponse extends jspb.Message {
  getUser(): zitadel_user_pb.User | undefined;
  setUser(value?: zitadel_user_pb.User): GetMyUserResponse;
  hasUser(): boolean;
  clearUser(): GetMyUserResponse;

  getLastLogin(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastLogin(value?: google_protobuf_timestamp_pb.Timestamp): GetMyUserResponse;
  hasLastLogin(): boolean;
  clearLastLogin(): GetMyUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyUserResponse): GetMyUserResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyUserResponse;
  static deserializeBinaryFromReader(message: GetMyUserResponse, reader: jspb.BinaryReader): GetMyUserResponse;
}

export namespace GetMyUserResponse {
  export type AsObject = {
    user?: zitadel_user_pb.User.AsObject,
    lastLogin?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class RemoveMyUserRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyUserRequest): RemoveMyUserRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyUserRequest;
  static deserializeBinaryFromReader(message: RemoveMyUserRequest, reader: jspb.BinaryReader): RemoveMyUserRequest;
}

export namespace RemoveMyUserRequest {
  export type AsObject = {
  }
}

export class RemoveMyUserResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyUserResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyUserResponse): RemoveMyUserResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyUserResponse;
  static deserializeBinaryFromReader(message: RemoveMyUserResponse, reader: jspb.BinaryReader): RemoveMyUserResponse;
}

export namespace RemoveMyUserResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListMyUserChangesRequest extends jspb.Message {
  getQuery(): zitadel_change_pb.ChangeQuery | undefined;
  setQuery(value?: zitadel_change_pb.ChangeQuery): ListMyUserChangesRequest;
  hasQuery(): boolean;
  clearQuery(): ListMyUserChangesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyUserChangesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyUserChangesRequest): ListMyUserChangesRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyUserChangesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyUserChangesRequest;
  static deserializeBinaryFromReader(message: ListMyUserChangesRequest, reader: jspb.BinaryReader): ListMyUserChangesRequest;
}

export namespace ListMyUserChangesRequest {
  export type AsObject = {
    query?: zitadel_change_pb.ChangeQuery.AsObject,
  }
}

export class ListMyUserChangesResponse extends jspb.Message {
  getResultList(): Array<zitadel_change_pb.Change>;
  setResultList(value: Array<zitadel_change_pb.Change>): ListMyUserChangesResponse;
  clearResultList(): ListMyUserChangesResponse;
  addResult(value?: zitadel_change_pb.Change, index?: number): zitadel_change_pb.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyUserChangesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyUserChangesResponse): ListMyUserChangesResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyUserChangesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyUserChangesResponse;
  static deserializeBinaryFromReader(message: ListMyUserChangesResponse, reader: jspb.BinaryReader): ListMyUserChangesResponse;
}

export namespace ListMyUserChangesResponse {
  export type AsObject = {
    resultList: Array<zitadel_change_pb.Change.AsObject>,
  }
}

export class ListMyUserSessionsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyUserSessionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyUserSessionsRequest): ListMyUserSessionsRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyUserSessionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyUserSessionsRequest;
  static deserializeBinaryFromReader(message: ListMyUserSessionsRequest, reader: jspb.BinaryReader): ListMyUserSessionsRequest;
}

export namespace ListMyUserSessionsRequest {
  export type AsObject = {
  }
}

export class ListMyUserSessionsResponse extends jspb.Message {
  getResultList(): Array<zitadel_user_pb.Session>;
  setResultList(value: Array<zitadel_user_pb.Session>): ListMyUserSessionsResponse;
  clearResultList(): ListMyUserSessionsResponse;
  addResult(value?: zitadel_user_pb.Session, index?: number): zitadel_user_pb.Session;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyUserSessionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyUserSessionsResponse): ListMyUserSessionsResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyUserSessionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyUserSessionsResponse;
  static deserializeBinaryFromReader(message: ListMyUserSessionsResponse, reader: jspb.BinaryReader): ListMyUserSessionsResponse;
}

export namespace ListMyUserSessionsResponse {
  export type AsObject = {
    resultList: Array<zitadel_user_pb.Session.AsObject>,
  }
}

export class ListMyMetadataRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListMyMetadataRequest;
  hasQuery(): boolean;
  clearQuery(): ListMyMetadataRequest;

  getQueriesList(): Array<zitadel_metadata_pb.MetadataQuery>;
  setQueriesList(value: Array<zitadel_metadata_pb.MetadataQuery>): ListMyMetadataRequest;
  clearQueriesList(): ListMyMetadataRequest;
  addQueries(value?: zitadel_metadata_pb.MetadataQuery, index?: number): zitadel_metadata_pb.MetadataQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyMetadataRequest): ListMyMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyMetadataRequest;
  static deserializeBinaryFromReader(message: ListMyMetadataRequest, reader: jspb.BinaryReader): ListMyMetadataRequest;
}

export namespace ListMyMetadataRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_metadata_pb.MetadataQuery.AsObject>,
  }
}

export class ListMyMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListMyMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): ListMyMetadataResponse;

  getResultList(): Array<zitadel_metadata_pb.Metadata>;
  setResultList(value: Array<zitadel_metadata_pb.Metadata>): ListMyMetadataResponse;
  clearResultList(): ListMyMetadataResponse;
  addResult(value?: zitadel_metadata_pb.Metadata, index?: number): zitadel_metadata_pb.Metadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyMetadataResponse): ListMyMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyMetadataResponse;
  static deserializeBinaryFromReader(message: ListMyMetadataResponse, reader: jspb.BinaryReader): ListMyMetadataResponse;
}

export namespace ListMyMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_metadata_pb.Metadata.AsObject>,
  }
}

export class GetMyMetadataRequest extends jspb.Message {
  getKey(): string;
  setKey(value: string): GetMyMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyMetadataRequest): GetMyMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyMetadataRequest;
  static deserializeBinaryFromReader(message: GetMyMetadataRequest, reader: jspb.BinaryReader): GetMyMetadataRequest;
}

export namespace GetMyMetadataRequest {
  export type AsObject = {
    key: string,
  }
}

export class GetMyMetadataResponse extends jspb.Message {
  getMetadata(): zitadel_metadata_pb.Metadata | undefined;
  setMetadata(value?: zitadel_metadata_pb.Metadata): GetMyMetadataResponse;
  hasMetadata(): boolean;
  clearMetadata(): GetMyMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyMetadataResponse): GetMyMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyMetadataResponse;
  static deserializeBinaryFromReader(message: GetMyMetadataResponse, reader: jspb.BinaryReader): GetMyMetadataResponse;
}

export namespace GetMyMetadataResponse {
  export type AsObject = {
    metadata?: zitadel_metadata_pb.Metadata.AsObject,
  }
}

export class SetMyMetadataRequest extends jspb.Message {
  getKey(): string;
  setKey(value: string): SetMyMetadataRequest;

  getValue(): Uint8Array | string;
  getValue_asU8(): Uint8Array;
  getValue_asB64(): string;
  setValue(value: Uint8Array | string): SetMyMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetMyMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetMyMetadataRequest): SetMyMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: SetMyMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetMyMetadataRequest;
  static deserializeBinaryFromReader(message: SetMyMetadataRequest, reader: jspb.BinaryReader): SetMyMetadataRequest;
}

export namespace SetMyMetadataRequest {
  export type AsObject = {
    key: string,
    value: Uint8Array | string,
  }
}

export class SetMyMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetMyMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): SetMyMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetMyMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetMyMetadataResponse): SetMyMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: SetMyMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetMyMetadataResponse;
  static deserializeBinaryFromReader(message: SetMyMetadataResponse, reader: jspb.BinaryReader): SetMyMetadataResponse;
}

export namespace SetMyMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkSetMyMetadataRequest extends jspb.Message {
  getMetadataList(): Array<BulkSetMyMetadataRequest.Metadata>;
  setMetadataList(value: Array<BulkSetMyMetadataRequest.Metadata>): BulkSetMyMetadataRequest;
  clearMetadataList(): BulkSetMyMetadataRequest;
  addMetadata(value?: BulkSetMyMetadataRequest.Metadata, index?: number): BulkSetMyMetadataRequest.Metadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSetMyMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSetMyMetadataRequest): BulkSetMyMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: BulkSetMyMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSetMyMetadataRequest;
  static deserializeBinaryFromReader(message: BulkSetMyMetadataRequest, reader: jspb.BinaryReader): BulkSetMyMetadataRequest;
}

export namespace BulkSetMyMetadataRequest {
  export type AsObject = {
    metadataList: Array<BulkSetMyMetadataRequest.Metadata.AsObject>,
  }

  export class Metadata extends jspb.Message {
    getKey(): string;
    setKey(value: string): Metadata;

    getValue(): Uint8Array | string;
    getValue_asU8(): Uint8Array;
    getValue_asB64(): string;
    setValue(value: Uint8Array | string): Metadata;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Metadata.AsObject;
    static toObject(includeInstance: boolean, msg: Metadata): Metadata.AsObject;
    static serializeBinaryToWriter(message: Metadata, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Metadata;
    static deserializeBinaryFromReader(message: Metadata, reader: jspb.BinaryReader): Metadata;
  }

  export namespace Metadata {
    export type AsObject = {
      key: string,
      value: Uint8Array | string,
    }
  }

}

export class BulkSetMyMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): BulkSetMyMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): BulkSetMyMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSetMyMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSetMyMetadataResponse): BulkSetMyMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: BulkSetMyMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSetMyMetadataResponse;
  static deserializeBinaryFromReader(message: BulkSetMyMetadataResponse, reader: jspb.BinaryReader): BulkSetMyMetadataResponse;
}

export namespace BulkSetMyMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMyMetadataRequest extends jspb.Message {
  getKey(): string;
  setKey(value: string): RemoveMyMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyMetadataRequest): RemoveMyMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyMetadataRequest;
  static deserializeBinaryFromReader(message: RemoveMyMetadataRequest, reader: jspb.BinaryReader): RemoveMyMetadataRequest;
}

export namespace RemoveMyMetadataRequest {
  export type AsObject = {
    key: string,
  }
}

export class RemoveMyMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyMetadataResponse): RemoveMyMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyMetadataResponse;
  static deserializeBinaryFromReader(message: RemoveMyMetadataResponse, reader: jspb.BinaryReader): RemoveMyMetadataResponse;
}

export namespace RemoveMyMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkRemoveMyMetadataRequest extends jspb.Message {
  getKeysList(): Array<string>;
  setKeysList(value: Array<string>): BulkRemoveMyMetadataRequest;
  clearKeysList(): BulkRemoveMyMetadataRequest;
  addKeys(value: string, index?: number): BulkRemoveMyMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkRemoveMyMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkRemoveMyMetadataRequest): BulkRemoveMyMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: BulkRemoveMyMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkRemoveMyMetadataRequest;
  static deserializeBinaryFromReader(message: BulkRemoveMyMetadataRequest, reader: jspb.BinaryReader): BulkRemoveMyMetadataRequest;
}

export namespace BulkRemoveMyMetadataRequest {
  export type AsObject = {
    keysList: Array<string>,
  }
}

export class BulkRemoveMyMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): BulkRemoveMyMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): BulkRemoveMyMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkRemoveMyMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkRemoveMyMetadataResponse): BulkRemoveMyMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: BulkRemoveMyMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkRemoveMyMetadataResponse;
  static deserializeBinaryFromReader(message: BulkRemoveMyMetadataResponse, reader: jspb.BinaryReader): BulkRemoveMyMetadataResponse;
}

export namespace BulkRemoveMyMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListMyRefreshTokensRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyRefreshTokensRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyRefreshTokensRequest): ListMyRefreshTokensRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyRefreshTokensRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyRefreshTokensRequest;
  static deserializeBinaryFromReader(message: ListMyRefreshTokensRequest, reader: jspb.BinaryReader): ListMyRefreshTokensRequest;
}

export namespace ListMyRefreshTokensRequest {
  export type AsObject = {
  }
}

export class ListMyRefreshTokensResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListMyRefreshTokensResponse;
  hasDetails(): boolean;
  clearDetails(): ListMyRefreshTokensResponse;

  getResultList(): Array<zitadel_user_pb.RefreshToken>;
  setResultList(value: Array<zitadel_user_pb.RefreshToken>): ListMyRefreshTokensResponse;
  clearResultList(): ListMyRefreshTokensResponse;
  addResult(value?: zitadel_user_pb.RefreshToken, index?: number): zitadel_user_pb.RefreshToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyRefreshTokensResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyRefreshTokensResponse): ListMyRefreshTokensResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyRefreshTokensResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyRefreshTokensResponse;
  static deserializeBinaryFromReader(message: ListMyRefreshTokensResponse, reader: jspb.BinaryReader): ListMyRefreshTokensResponse;
}

export namespace ListMyRefreshTokensResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_user_pb.RefreshToken.AsObject>,
  }
}

export class RevokeMyRefreshTokenRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RevokeMyRefreshTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeMyRefreshTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeMyRefreshTokenRequest): RevokeMyRefreshTokenRequest.AsObject;
  static serializeBinaryToWriter(message: RevokeMyRefreshTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeMyRefreshTokenRequest;
  static deserializeBinaryFromReader(message: RevokeMyRefreshTokenRequest, reader: jspb.BinaryReader): RevokeMyRefreshTokenRequest;
}

export namespace RevokeMyRefreshTokenRequest {
  export type AsObject = {
    id: string,
  }
}

export class RevokeMyRefreshTokenResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RevokeMyRefreshTokenResponse;
  hasDetails(): boolean;
  clearDetails(): RevokeMyRefreshTokenResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeMyRefreshTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeMyRefreshTokenResponse): RevokeMyRefreshTokenResponse.AsObject;
  static serializeBinaryToWriter(message: RevokeMyRefreshTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeMyRefreshTokenResponse;
  static deserializeBinaryFromReader(message: RevokeMyRefreshTokenResponse, reader: jspb.BinaryReader): RevokeMyRefreshTokenResponse;
}

export namespace RevokeMyRefreshTokenResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RevokeAllMyRefreshTokensRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeAllMyRefreshTokensRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeAllMyRefreshTokensRequest): RevokeAllMyRefreshTokensRequest.AsObject;
  static serializeBinaryToWriter(message: RevokeAllMyRefreshTokensRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeAllMyRefreshTokensRequest;
  static deserializeBinaryFromReader(message: RevokeAllMyRefreshTokensRequest, reader: jspb.BinaryReader): RevokeAllMyRefreshTokensRequest;
}

export namespace RevokeAllMyRefreshTokensRequest {
  export type AsObject = {
  }
}

export class RevokeAllMyRefreshTokensResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RevokeAllMyRefreshTokensResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RevokeAllMyRefreshTokensResponse): RevokeAllMyRefreshTokensResponse.AsObject;
  static serializeBinaryToWriter(message: RevokeAllMyRefreshTokensResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RevokeAllMyRefreshTokensResponse;
  static deserializeBinaryFromReader(message: RevokeAllMyRefreshTokensResponse, reader: jspb.BinaryReader): RevokeAllMyRefreshTokensResponse;
}

export namespace RevokeAllMyRefreshTokensResponse {
  export type AsObject = {
  }
}

export class UpdateMyUserNameRequest extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): UpdateMyUserNameRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMyUserNameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMyUserNameRequest): UpdateMyUserNameRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateMyUserNameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMyUserNameRequest;
  static deserializeBinaryFromReader(message: UpdateMyUserNameRequest, reader: jspb.BinaryReader): UpdateMyUserNameRequest;
}

export namespace UpdateMyUserNameRequest {
  export type AsObject = {
    userName: string,
  }
}

export class UpdateMyUserNameResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateMyUserNameResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateMyUserNameResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMyUserNameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMyUserNameResponse): UpdateMyUserNameResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateMyUserNameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMyUserNameResponse;
  static deserializeBinaryFromReader(message: UpdateMyUserNameResponse, reader: jspb.BinaryReader): UpdateMyUserNameResponse;
}

export namespace UpdateMyUserNameResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetMyPasswordComplexityPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyPasswordComplexityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyPasswordComplexityPolicyRequest): GetMyPasswordComplexityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyPasswordComplexityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyPasswordComplexityPolicyRequest;
  static deserializeBinaryFromReader(message: GetMyPasswordComplexityPolicyRequest, reader: jspb.BinaryReader): GetMyPasswordComplexityPolicyRequest;
}

export namespace GetMyPasswordComplexityPolicyRequest {
  export type AsObject = {
  }
}

export class GetMyPasswordComplexityPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PasswordComplexityPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PasswordComplexityPolicy): GetMyPasswordComplexityPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetMyPasswordComplexityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyPasswordComplexityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyPasswordComplexityPolicyResponse): GetMyPasswordComplexityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyPasswordComplexityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyPasswordComplexityPolicyResponse;
  static deserializeBinaryFromReader(message: GetMyPasswordComplexityPolicyResponse, reader: jspb.BinaryReader): GetMyPasswordComplexityPolicyResponse;
}

export namespace GetMyPasswordComplexityPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PasswordComplexityPolicy.AsObject,
  }
}

export class UpdateMyPasswordRequest extends jspb.Message {
  getOldPassword(): string;
  setOldPassword(value: string): UpdateMyPasswordRequest;

  getNewPassword(): string;
  setNewPassword(value: string): UpdateMyPasswordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMyPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMyPasswordRequest): UpdateMyPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateMyPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMyPasswordRequest;
  static deserializeBinaryFromReader(message: UpdateMyPasswordRequest, reader: jspb.BinaryReader): UpdateMyPasswordRequest;
}

export namespace UpdateMyPasswordRequest {
  export type AsObject = {
    oldPassword: string,
    newPassword: string,
  }
}

export class UpdateMyPasswordResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateMyPasswordResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateMyPasswordResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMyPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMyPasswordResponse): UpdateMyPasswordResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateMyPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMyPasswordResponse;
  static deserializeBinaryFromReader(message: UpdateMyPasswordResponse, reader: jspb.BinaryReader): UpdateMyPasswordResponse;
}

export namespace UpdateMyPasswordResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetMyProfileRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyProfileRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyProfileRequest): GetMyProfileRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyProfileRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyProfileRequest;
  static deserializeBinaryFromReader(message: GetMyProfileRequest, reader: jspb.BinaryReader): GetMyProfileRequest;
}

export namespace GetMyProfileRequest {
  export type AsObject = {
  }
}

export class GetMyProfileResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GetMyProfileResponse;
  hasDetails(): boolean;
  clearDetails(): GetMyProfileResponse;

  getProfile(): zitadel_user_pb.Profile | undefined;
  setProfile(value?: zitadel_user_pb.Profile): GetMyProfileResponse;
  hasProfile(): boolean;
  clearProfile(): GetMyProfileResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyProfileResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyProfileResponse): GetMyProfileResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyProfileResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyProfileResponse;
  static deserializeBinaryFromReader(message: GetMyProfileResponse, reader: jspb.BinaryReader): GetMyProfileResponse;
}

export namespace GetMyProfileResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    profile?: zitadel_user_pb.Profile.AsObject,
  }
}

export class UpdateMyProfileRequest extends jspb.Message {
  getFirstName(): string;
  setFirstName(value: string): UpdateMyProfileRequest;

  getLastName(): string;
  setLastName(value: string): UpdateMyProfileRequest;

  getNickName(): string;
  setNickName(value: string): UpdateMyProfileRequest;

  getDisplayName(): string;
  setDisplayName(value: string): UpdateMyProfileRequest;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): UpdateMyProfileRequest;

  getGender(): zitadel_user_pb.Gender;
  setGender(value: zitadel_user_pb.Gender): UpdateMyProfileRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMyProfileRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMyProfileRequest): UpdateMyProfileRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateMyProfileRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMyProfileRequest;
  static deserializeBinaryFromReader(message: UpdateMyProfileRequest, reader: jspb.BinaryReader): UpdateMyProfileRequest;
}

export namespace UpdateMyProfileRequest {
  export type AsObject = {
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: zitadel_user_pb.Gender,
  }
}

export class UpdateMyProfileResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateMyProfileResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateMyProfileResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMyProfileResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMyProfileResponse): UpdateMyProfileResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateMyProfileResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMyProfileResponse;
  static deserializeBinaryFromReader(message: UpdateMyProfileResponse, reader: jspb.BinaryReader): UpdateMyProfileResponse;
}

export namespace UpdateMyProfileResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetMyEmailRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyEmailRequest): GetMyEmailRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyEmailRequest;
  static deserializeBinaryFromReader(message: GetMyEmailRequest, reader: jspb.BinaryReader): GetMyEmailRequest;
}

export namespace GetMyEmailRequest {
  export type AsObject = {
  }
}

export class GetMyEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GetMyEmailResponse;
  hasDetails(): boolean;
  clearDetails(): GetMyEmailResponse;

  getEmail(): zitadel_user_pb.Email | undefined;
  setEmail(value?: zitadel_user_pb.Email): GetMyEmailResponse;
  hasEmail(): boolean;
  clearEmail(): GetMyEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyEmailResponse): GetMyEmailResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyEmailResponse;
  static deserializeBinaryFromReader(message: GetMyEmailResponse, reader: jspb.BinaryReader): GetMyEmailResponse;
}

export namespace GetMyEmailResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    email?: zitadel_user_pb.Email.AsObject,
  }
}

export class SetMyEmailRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): SetMyEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetMyEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetMyEmailRequest): SetMyEmailRequest.AsObject;
  static serializeBinaryToWriter(message: SetMyEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetMyEmailRequest;
  static deserializeBinaryFromReader(message: SetMyEmailRequest, reader: jspb.BinaryReader): SetMyEmailRequest;
}

export namespace SetMyEmailRequest {
  export type AsObject = {
    email: string,
  }
}

export class SetMyEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetMyEmailResponse;
  hasDetails(): boolean;
  clearDetails(): SetMyEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetMyEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetMyEmailResponse): SetMyEmailResponse.AsObject;
  static serializeBinaryToWriter(message: SetMyEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetMyEmailResponse;
  static deserializeBinaryFromReader(message: SetMyEmailResponse, reader: jspb.BinaryReader): SetMyEmailResponse;
}

export namespace SetMyEmailResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class VerifyMyEmailRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): VerifyMyEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyEmailRequest): VerifyMyEmailRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyMyEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyEmailRequest;
  static deserializeBinaryFromReader(message: VerifyMyEmailRequest, reader: jspb.BinaryReader): VerifyMyEmailRequest;
}

export namespace VerifyMyEmailRequest {
  export type AsObject = {
    code: string,
  }
}

export class VerifyMyEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): VerifyMyEmailResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyMyEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyEmailResponse): VerifyMyEmailResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyMyEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyEmailResponse;
  static deserializeBinaryFromReader(message: VerifyMyEmailResponse, reader: jspb.BinaryReader): VerifyMyEmailResponse;
}

export namespace VerifyMyEmailResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResendMyEmailVerificationRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendMyEmailVerificationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendMyEmailVerificationRequest): ResendMyEmailVerificationRequest.AsObject;
  static serializeBinaryToWriter(message: ResendMyEmailVerificationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendMyEmailVerificationRequest;
  static deserializeBinaryFromReader(message: ResendMyEmailVerificationRequest, reader: jspb.BinaryReader): ResendMyEmailVerificationRequest;
}

export namespace ResendMyEmailVerificationRequest {
  export type AsObject = {
  }
}

export class ResendMyEmailVerificationResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResendMyEmailVerificationResponse;
  hasDetails(): boolean;
  clearDetails(): ResendMyEmailVerificationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendMyEmailVerificationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendMyEmailVerificationResponse): ResendMyEmailVerificationResponse.AsObject;
  static serializeBinaryToWriter(message: ResendMyEmailVerificationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendMyEmailVerificationResponse;
  static deserializeBinaryFromReader(message: ResendMyEmailVerificationResponse, reader: jspb.BinaryReader): ResendMyEmailVerificationResponse;
}

export namespace ResendMyEmailVerificationResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetMyPhoneRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyPhoneRequest): GetMyPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyPhoneRequest;
  static deserializeBinaryFromReader(message: GetMyPhoneRequest, reader: jspb.BinaryReader): GetMyPhoneRequest;
}

export namespace GetMyPhoneRequest {
  export type AsObject = {
  }
}

export class GetMyPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GetMyPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): GetMyPhoneResponse;

  getPhone(): zitadel_user_pb.Phone | undefined;
  setPhone(value?: zitadel_user_pb.Phone): GetMyPhoneResponse;
  hasPhone(): boolean;
  clearPhone(): GetMyPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyPhoneResponse): GetMyPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyPhoneResponse;
  static deserializeBinaryFromReader(message: GetMyPhoneResponse, reader: jspb.BinaryReader): GetMyPhoneResponse;
}

export namespace GetMyPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    phone?: zitadel_user_pb.Phone.AsObject,
  }
}

export class SetMyPhoneRequest extends jspb.Message {
  getPhone(): string;
  setPhone(value: string): SetMyPhoneRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetMyPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetMyPhoneRequest): SetMyPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: SetMyPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetMyPhoneRequest;
  static deserializeBinaryFromReader(message: SetMyPhoneRequest, reader: jspb.BinaryReader): SetMyPhoneRequest;
}

export namespace SetMyPhoneRequest {
  export type AsObject = {
    phone: string,
  }
}

export class SetMyPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetMyPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): SetMyPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetMyPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetMyPhoneResponse): SetMyPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: SetMyPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetMyPhoneResponse;
  static deserializeBinaryFromReader(message: SetMyPhoneResponse, reader: jspb.BinaryReader): SetMyPhoneResponse;
}

export namespace SetMyPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class VerifyMyPhoneRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): VerifyMyPhoneRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyPhoneRequest): VerifyMyPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyMyPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyPhoneRequest;
  static deserializeBinaryFromReader(message: VerifyMyPhoneRequest, reader: jspb.BinaryReader): VerifyMyPhoneRequest;
}

export namespace VerifyMyPhoneRequest {
  export type AsObject = {
    code: string,
  }
}

export class VerifyMyPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): VerifyMyPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyMyPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyPhoneResponse): VerifyMyPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyMyPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyPhoneResponse;
  static deserializeBinaryFromReader(message: VerifyMyPhoneResponse, reader: jspb.BinaryReader): VerifyMyPhoneResponse;
}

export namespace VerifyMyPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResendMyPhoneVerificationRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendMyPhoneVerificationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendMyPhoneVerificationRequest): ResendMyPhoneVerificationRequest.AsObject;
  static serializeBinaryToWriter(message: ResendMyPhoneVerificationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendMyPhoneVerificationRequest;
  static deserializeBinaryFromReader(message: ResendMyPhoneVerificationRequest, reader: jspb.BinaryReader): ResendMyPhoneVerificationRequest;
}

export namespace ResendMyPhoneVerificationRequest {
  export type AsObject = {
  }
}

export class ResendMyPhoneVerificationResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResendMyPhoneVerificationResponse;
  hasDetails(): boolean;
  clearDetails(): ResendMyPhoneVerificationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendMyPhoneVerificationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendMyPhoneVerificationResponse): ResendMyPhoneVerificationResponse.AsObject;
  static serializeBinaryToWriter(message: ResendMyPhoneVerificationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendMyPhoneVerificationResponse;
  static deserializeBinaryFromReader(message: ResendMyPhoneVerificationResponse, reader: jspb.BinaryReader): ResendMyPhoneVerificationResponse;
}

export namespace ResendMyPhoneVerificationResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMyPhoneRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyPhoneRequest): RemoveMyPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyPhoneRequest;
  static deserializeBinaryFromReader(message: RemoveMyPhoneRequest, reader: jspb.BinaryReader): RemoveMyPhoneRequest;
}

export namespace RemoveMyPhoneRequest {
  export type AsObject = {
  }
}

export class RemoveMyPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyPhoneResponse): RemoveMyPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyPhoneResponse;
  static deserializeBinaryFromReader(message: RemoveMyPhoneResponse, reader: jspb.BinaryReader): RemoveMyPhoneResponse;
}

export namespace RemoveMyPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMyAvatarRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAvatarRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAvatarRequest): RemoveMyAvatarRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAvatarRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAvatarRequest;
  static deserializeBinaryFromReader(message: RemoveMyAvatarRequest, reader: jspb.BinaryReader): RemoveMyAvatarRequest;
}

export namespace RemoveMyAvatarRequest {
  export type AsObject = {
  }
}

export class RemoveMyAvatarResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyAvatarResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyAvatarResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAvatarResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAvatarResponse): RemoveMyAvatarResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAvatarResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAvatarResponse;
  static deserializeBinaryFromReader(message: RemoveMyAvatarResponse, reader: jspb.BinaryReader): RemoveMyAvatarResponse;
}

export namespace RemoveMyAvatarResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListMyLinkedIDPsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListMyLinkedIDPsRequest;
  hasQuery(): boolean;
  clearQuery(): ListMyLinkedIDPsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyLinkedIDPsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyLinkedIDPsRequest): ListMyLinkedIDPsRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyLinkedIDPsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyLinkedIDPsRequest;
  static deserializeBinaryFromReader(message: ListMyLinkedIDPsRequest, reader: jspb.BinaryReader): ListMyLinkedIDPsRequest;
}

export namespace ListMyLinkedIDPsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListMyLinkedIDPsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListMyLinkedIDPsResponse;
  hasDetails(): boolean;
  clearDetails(): ListMyLinkedIDPsResponse;

  getResultList(): Array<zitadel_idp_pb.IDPUserLink>;
  setResultList(value: Array<zitadel_idp_pb.IDPUserLink>): ListMyLinkedIDPsResponse;
  clearResultList(): ListMyLinkedIDPsResponse;
  addResult(value?: zitadel_idp_pb.IDPUserLink, index?: number): zitadel_idp_pb.IDPUserLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyLinkedIDPsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyLinkedIDPsResponse): ListMyLinkedIDPsResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyLinkedIDPsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyLinkedIDPsResponse;
  static deserializeBinaryFromReader(message: ListMyLinkedIDPsResponse, reader: jspb.BinaryReader): ListMyLinkedIDPsResponse;
}

export namespace ListMyLinkedIDPsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_idp_pb.IDPUserLink.AsObject>,
  }
}

export class RemoveMyLinkedIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): RemoveMyLinkedIDPRequest;

  getLinkedUserId(): string;
  setLinkedUserId(value: string): RemoveMyLinkedIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyLinkedIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyLinkedIDPRequest): RemoveMyLinkedIDPRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyLinkedIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyLinkedIDPRequest;
  static deserializeBinaryFromReader(message: RemoveMyLinkedIDPRequest, reader: jspb.BinaryReader): RemoveMyLinkedIDPRequest;
}

export namespace RemoveMyLinkedIDPRequest {
  export type AsObject = {
    idpId: string,
    linkedUserId: string,
  }
}

export class RemoveMyLinkedIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyLinkedIDPResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyLinkedIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyLinkedIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyLinkedIDPResponse): RemoveMyLinkedIDPResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyLinkedIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyLinkedIDPResponse;
  static deserializeBinaryFromReader(message: RemoveMyLinkedIDPResponse, reader: jspb.BinaryReader): RemoveMyLinkedIDPResponse;
}

export namespace RemoveMyLinkedIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListMyAuthFactorsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyAuthFactorsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyAuthFactorsRequest): ListMyAuthFactorsRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyAuthFactorsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyAuthFactorsRequest;
  static deserializeBinaryFromReader(message: ListMyAuthFactorsRequest, reader: jspb.BinaryReader): ListMyAuthFactorsRequest;
}

export namespace ListMyAuthFactorsRequest {
  export type AsObject = {
  }
}

export class ListMyAuthFactorsResponse extends jspb.Message {
  getResultList(): Array<zitadel_user_pb.AuthFactor>;
  setResultList(value: Array<zitadel_user_pb.AuthFactor>): ListMyAuthFactorsResponse;
  clearResultList(): ListMyAuthFactorsResponse;
  addResult(value?: zitadel_user_pb.AuthFactor, index?: number): zitadel_user_pb.AuthFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyAuthFactorsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyAuthFactorsResponse): ListMyAuthFactorsResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyAuthFactorsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyAuthFactorsResponse;
  static deserializeBinaryFromReader(message: ListMyAuthFactorsResponse, reader: jspb.BinaryReader): ListMyAuthFactorsResponse;
}

export namespace ListMyAuthFactorsResponse {
  export type AsObject = {
    resultList: Array<zitadel_user_pb.AuthFactor.AsObject>,
  }
}

export class AddMyAuthFactorU2FRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyAuthFactorU2FRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyAuthFactorU2FRequest): AddMyAuthFactorU2FRequest.AsObject;
  static serializeBinaryToWriter(message: AddMyAuthFactorU2FRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyAuthFactorU2FRequest;
  static deserializeBinaryFromReader(message: AddMyAuthFactorU2FRequest, reader: jspb.BinaryReader): AddMyAuthFactorU2FRequest;
}

export namespace AddMyAuthFactorU2FRequest {
  export type AsObject = {
  }
}

export class AddMyAuthFactorU2FResponse extends jspb.Message {
  getKey(): zitadel_user_pb.WebAuthNKey | undefined;
  setKey(value?: zitadel_user_pb.WebAuthNKey): AddMyAuthFactorU2FResponse;
  hasKey(): boolean;
  clearKey(): AddMyAuthFactorU2FResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMyAuthFactorU2FResponse;
  hasDetails(): boolean;
  clearDetails(): AddMyAuthFactorU2FResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyAuthFactorU2FResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyAuthFactorU2FResponse): AddMyAuthFactorU2FResponse.AsObject;
  static serializeBinaryToWriter(message: AddMyAuthFactorU2FResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyAuthFactorU2FResponse;
  static deserializeBinaryFromReader(message: AddMyAuthFactorU2FResponse, reader: jspb.BinaryReader): AddMyAuthFactorU2FResponse;
}

export namespace AddMyAuthFactorU2FResponse {
  export type AsObject = {
    key?: zitadel_user_pb.WebAuthNKey.AsObject,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddMyAuthFactorOTPRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyAuthFactorOTPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyAuthFactorOTPRequest): AddMyAuthFactorOTPRequest.AsObject;
  static serializeBinaryToWriter(message: AddMyAuthFactorOTPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyAuthFactorOTPRequest;
  static deserializeBinaryFromReader(message: AddMyAuthFactorOTPRequest, reader: jspb.BinaryReader): AddMyAuthFactorOTPRequest;
}

export namespace AddMyAuthFactorOTPRequest {
  export type AsObject = {
  }
}

export class AddMyAuthFactorOTPResponse extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): AddMyAuthFactorOTPResponse;

  getSecret(): string;
  setSecret(value: string): AddMyAuthFactorOTPResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMyAuthFactorOTPResponse;
  hasDetails(): boolean;
  clearDetails(): AddMyAuthFactorOTPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyAuthFactorOTPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyAuthFactorOTPResponse): AddMyAuthFactorOTPResponse.AsObject;
  static serializeBinaryToWriter(message: AddMyAuthFactorOTPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyAuthFactorOTPResponse;
  static deserializeBinaryFromReader(message: AddMyAuthFactorOTPResponse, reader: jspb.BinaryReader): AddMyAuthFactorOTPResponse;
}

export namespace AddMyAuthFactorOTPResponse {
  export type AsObject = {
    url: string,
    secret: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class VerifyMyAuthFactorOTPRequest extends jspb.Message {
  getCode(): string;
  setCode(value: string): VerifyMyAuthFactorOTPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyAuthFactorOTPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyAuthFactorOTPRequest): VerifyMyAuthFactorOTPRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyMyAuthFactorOTPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyAuthFactorOTPRequest;
  static deserializeBinaryFromReader(message: VerifyMyAuthFactorOTPRequest, reader: jspb.BinaryReader): VerifyMyAuthFactorOTPRequest;
}

export namespace VerifyMyAuthFactorOTPRequest {
  export type AsObject = {
    code: string,
  }
}

export class VerifyMyAuthFactorOTPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): VerifyMyAuthFactorOTPResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyMyAuthFactorOTPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyAuthFactorOTPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyAuthFactorOTPResponse): VerifyMyAuthFactorOTPResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyMyAuthFactorOTPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyAuthFactorOTPResponse;
  static deserializeBinaryFromReader(message: VerifyMyAuthFactorOTPResponse, reader: jspb.BinaryReader): VerifyMyAuthFactorOTPResponse;
}

export namespace VerifyMyAuthFactorOTPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class VerifyMyAuthFactorU2FRequest extends jspb.Message {
  getVerification(): zitadel_user_pb.WebAuthNVerification | undefined;
  setVerification(value?: zitadel_user_pb.WebAuthNVerification): VerifyMyAuthFactorU2FRequest;
  hasVerification(): boolean;
  clearVerification(): VerifyMyAuthFactorU2FRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyAuthFactorU2FRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyAuthFactorU2FRequest): VerifyMyAuthFactorU2FRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyMyAuthFactorU2FRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyAuthFactorU2FRequest;
  static deserializeBinaryFromReader(message: VerifyMyAuthFactorU2FRequest, reader: jspb.BinaryReader): VerifyMyAuthFactorU2FRequest;
}

export namespace VerifyMyAuthFactorU2FRequest {
  export type AsObject = {
    verification?: zitadel_user_pb.WebAuthNVerification.AsObject,
  }
}

export class VerifyMyAuthFactorU2FResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): VerifyMyAuthFactorU2FResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyMyAuthFactorU2FResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyAuthFactorU2FResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyAuthFactorU2FResponse): VerifyMyAuthFactorU2FResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyMyAuthFactorU2FResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyAuthFactorU2FResponse;
  static deserializeBinaryFromReader(message: VerifyMyAuthFactorU2FResponse, reader: jspb.BinaryReader): VerifyMyAuthFactorU2FResponse;
}

export namespace VerifyMyAuthFactorU2FResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMyAuthFactorOTPRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAuthFactorOTPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAuthFactorOTPRequest): RemoveMyAuthFactorOTPRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAuthFactorOTPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAuthFactorOTPRequest;
  static deserializeBinaryFromReader(message: RemoveMyAuthFactorOTPRequest, reader: jspb.BinaryReader): RemoveMyAuthFactorOTPRequest;
}

export namespace RemoveMyAuthFactorOTPRequest {
  export type AsObject = {
  }
}

export class RemoveMyAuthFactorOTPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyAuthFactorOTPResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyAuthFactorOTPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAuthFactorOTPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAuthFactorOTPResponse): RemoveMyAuthFactorOTPResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAuthFactorOTPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAuthFactorOTPResponse;
  static deserializeBinaryFromReader(message: RemoveMyAuthFactorOTPResponse, reader: jspb.BinaryReader): RemoveMyAuthFactorOTPResponse;
}

export namespace RemoveMyAuthFactorOTPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddMyAuthFactorOTPSMSRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyAuthFactorOTPSMSRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyAuthFactorOTPSMSRequest): AddMyAuthFactorOTPSMSRequest.AsObject;
  static serializeBinaryToWriter(message: AddMyAuthFactorOTPSMSRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyAuthFactorOTPSMSRequest;
  static deserializeBinaryFromReader(message: AddMyAuthFactorOTPSMSRequest, reader: jspb.BinaryReader): AddMyAuthFactorOTPSMSRequest;
}

export namespace AddMyAuthFactorOTPSMSRequest {
  export type AsObject = {
  }
}

export class AddMyAuthFactorOTPSMSResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMyAuthFactorOTPSMSResponse;
  hasDetails(): boolean;
  clearDetails(): AddMyAuthFactorOTPSMSResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyAuthFactorOTPSMSResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyAuthFactorOTPSMSResponse): AddMyAuthFactorOTPSMSResponse.AsObject;
  static serializeBinaryToWriter(message: AddMyAuthFactorOTPSMSResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyAuthFactorOTPSMSResponse;
  static deserializeBinaryFromReader(message: AddMyAuthFactorOTPSMSResponse, reader: jspb.BinaryReader): AddMyAuthFactorOTPSMSResponse;
}

export namespace AddMyAuthFactorOTPSMSResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMyAuthFactorOTPSMSRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAuthFactorOTPSMSRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAuthFactorOTPSMSRequest): RemoveMyAuthFactorOTPSMSRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAuthFactorOTPSMSRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAuthFactorOTPSMSRequest;
  static deserializeBinaryFromReader(message: RemoveMyAuthFactorOTPSMSRequest, reader: jspb.BinaryReader): RemoveMyAuthFactorOTPSMSRequest;
}

export namespace RemoveMyAuthFactorOTPSMSRequest {
  export type AsObject = {
  }
}

export class RemoveMyAuthFactorOTPSMSResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyAuthFactorOTPSMSResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyAuthFactorOTPSMSResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAuthFactorOTPSMSResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAuthFactorOTPSMSResponse): RemoveMyAuthFactorOTPSMSResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAuthFactorOTPSMSResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAuthFactorOTPSMSResponse;
  static deserializeBinaryFromReader(message: RemoveMyAuthFactorOTPSMSResponse, reader: jspb.BinaryReader): RemoveMyAuthFactorOTPSMSResponse;
}

export namespace RemoveMyAuthFactorOTPSMSResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddMyAuthFactorOTPEmailRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyAuthFactorOTPEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyAuthFactorOTPEmailRequest): AddMyAuthFactorOTPEmailRequest.AsObject;
  static serializeBinaryToWriter(message: AddMyAuthFactorOTPEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyAuthFactorOTPEmailRequest;
  static deserializeBinaryFromReader(message: AddMyAuthFactorOTPEmailRequest, reader: jspb.BinaryReader): AddMyAuthFactorOTPEmailRequest;
}

export namespace AddMyAuthFactorOTPEmailRequest {
  export type AsObject = {
  }
}

export class AddMyAuthFactorOTPEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMyAuthFactorOTPEmailResponse;
  hasDetails(): boolean;
  clearDetails(): AddMyAuthFactorOTPEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyAuthFactorOTPEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyAuthFactorOTPEmailResponse): AddMyAuthFactorOTPEmailResponse.AsObject;
  static serializeBinaryToWriter(message: AddMyAuthFactorOTPEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyAuthFactorOTPEmailResponse;
  static deserializeBinaryFromReader(message: AddMyAuthFactorOTPEmailResponse, reader: jspb.BinaryReader): AddMyAuthFactorOTPEmailResponse;
}

export namespace AddMyAuthFactorOTPEmailResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMyAuthFactorOTPEmailRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAuthFactorOTPEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAuthFactorOTPEmailRequest): RemoveMyAuthFactorOTPEmailRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAuthFactorOTPEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAuthFactorOTPEmailRequest;
  static deserializeBinaryFromReader(message: RemoveMyAuthFactorOTPEmailRequest, reader: jspb.BinaryReader): RemoveMyAuthFactorOTPEmailRequest;
}

export namespace RemoveMyAuthFactorOTPEmailRequest {
  export type AsObject = {
  }
}

export class RemoveMyAuthFactorOTPEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyAuthFactorOTPEmailResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyAuthFactorOTPEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAuthFactorOTPEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAuthFactorOTPEmailResponse): RemoveMyAuthFactorOTPEmailResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAuthFactorOTPEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAuthFactorOTPEmailResponse;
  static deserializeBinaryFromReader(message: RemoveMyAuthFactorOTPEmailResponse, reader: jspb.BinaryReader): RemoveMyAuthFactorOTPEmailResponse;
}

export namespace RemoveMyAuthFactorOTPEmailResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMyAuthFactorU2FRequest extends jspb.Message {
  getTokenId(): string;
  setTokenId(value: string): RemoveMyAuthFactorU2FRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAuthFactorU2FRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAuthFactorU2FRequest): RemoveMyAuthFactorU2FRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAuthFactorU2FRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAuthFactorU2FRequest;
  static deserializeBinaryFromReader(message: RemoveMyAuthFactorU2FRequest, reader: jspb.BinaryReader): RemoveMyAuthFactorU2FRequest;
}

export namespace RemoveMyAuthFactorU2FRequest {
  export type AsObject = {
    tokenId: string,
  }
}

export class RemoveMyAuthFactorU2FResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyAuthFactorU2FResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyAuthFactorU2FResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyAuthFactorU2FResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyAuthFactorU2FResponse): RemoveMyAuthFactorU2FResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyAuthFactorU2FResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyAuthFactorU2FResponse;
  static deserializeBinaryFromReader(message: RemoveMyAuthFactorU2FResponse, reader: jspb.BinaryReader): RemoveMyAuthFactorU2FResponse;
}

export namespace RemoveMyAuthFactorU2FResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListMyPasswordlessRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyPasswordlessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyPasswordlessRequest): ListMyPasswordlessRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyPasswordlessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyPasswordlessRequest;
  static deserializeBinaryFromReader(message: ListMyPasswordlessRequest, reader: jspb.BinaryReader): ListMyPasswordlessRequest;
}

export namespace ListMyPasswordlessRequest {
  export type AsObject = {
  }
}

export class ListMyPasswordlessResponse extends jspb.Message {
  getResultList(): Array<zitadel_user_pb.WebAuthNToken>;
  setResultList(value: Array<zitadel_user_pb.WebAuthNToken>): ListMyPasswordlessResponse;
  clearResultList(): ListMyPasswordlessResponse;
  addResult(value?: zitadel_user_pb.WebAuthNToken, index?: number): zitadel_user_pb.WebAuthNToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyPasswordlessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyPasswordlessResponse): ListMyPasswordlessResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyPasswordlessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyPasswordlessResponse;
  static deserializeBinaryFromReader(message: ListMyPasswordlessResponse, reader: jspb.BinaryReader): ListMyPasswordlessResponse;
}

export namespace ListMyPasswordlessResponse {
  export type AsObject = {
    resultList: Array<zitadel_user_pb.WebAuthNToken.AsObject>,
  }
}

export class AddMyPasswordlessRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyPasswordlessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyPasswordlessRequest): AddMyPasswordlessRequest.AsObject;
  static serializeBinaryToWriter(message: AddMyPasswordlessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyPasswordlessRequest;
  static deserializeBinaryFromReader(message: AddMyPasswordlessRequest, reader: jspb.BinaryReader): AddMyPasswordlessRequest;
}

export namespace AddMyPasswordlessRequest {
  export type AsObject = {
  }
}

export class AddMyPasswordlessResponse extends jspb.Message {
  getKey(): zitadel_user_pb.WebAuthNKey | undefined;
  setKey(value?: zitadel_user_pb.WebAuthNKey): AddMyPasswordlessResponse;
  hasKey(): boolean;
  clearKey(): AddMyPasswordlessResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMyPasswordlessResponse;
  hasDetails(): boolean;
  clearDetails(): AddMyPasswordlessResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyPasswordlessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyPasswordlessResponse): AddMyPasswordlessResponse.AsObject;
  static serializeBinaryToWriter(message: AddMyPasswordlessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyPasswordlessResponse;
  static deserializeBinaryFromReader(message: AddMyPasswordlessResponse, reader: jspb.BinaryReader): AddMyPasswordlessResponse;
}

export namespace AddMyPasswordlessResponse {
  export type AsObject = {
    key?: zitadel_user_pb.WebAuthNKey.AsObject,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddMyPasswordlessLinkRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyPasswordlessLinkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyPasswordlessLinkRequest): AddMyPasswordlessLinkRequest.AsObject;
  static serializeBinaryToWriter(message: AddMyPasswordlessLinkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyPasswordlessLinkRequest;
  static deserializeBinaryFromReader(message: AddMyPasswordlessLinkRequest, reader: jspb.BinaryReader): AddMyPasswordlessLinkRequest;
}

export namespace AddMyPasswordlessLinkRequest {
  export type AsObject = {
  }
}

export class AddMyPasswordlessLinkResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMyPasswordlessLinkResponse;
  hasDetails(): boolean;
  clearDetails(): AddMyPasswordlessLinkResponse;

  getLink(): string;
  setLink(value: string): AddMyPasswordlessLinkResponse;

  getExpiration(): google_protobuf_duration_pb.Duration | undefined;
  setExpiration(value?: google_protobuf_duration_pb.Duration): AddMyPasswordlessLinkResponse;
  hasExpiration(): boolean;
  clearExpiration(): AddMyPasswordlessLinkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMyPasswordlessLinkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMyPasswordlessLinkResponse): AddMyPasswordlessLinkResponse.AsObject;
  static serializeBinaryToWriter(message: AddMyPasswordlessLinkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMyPasswordlessLinkResponse;
  static deserializeBinaryFromReader(message: AddMyPasswordlessLinkResponse, reader: jspb.BinaryReader): AddMyPasswordlessLinkResponse;
}

export namespace AddMyPasswordlessLinkResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    link: string,
    expiration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class SendMyPasswordlessLinkRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendMyPasswordlessLinkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SendMyPasswordlessLinkRequest): SendMyPasswordlessLinkRequest.AsObject;
  static serializeBinaryToWriter(message: SendMyPasswordlessLinkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendMyPasswordlessLinkRequest;
  static deserializeBinaryFromReader(message: SendMyPasswordlessLinkRequest, reader: jspb.BinaryReader): SendMyPasswordlessLinkRequest;
}

export namespace SendMyPasswordlessLinkRequest {
  export type AsObject = {
  }
}

export class SendMyPasswordlessLinkResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SendMyPasswordlessLinkResponse;
  hasDetails(): boolean;
  clearDetails(): SendMyPasswordlessLinkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendMyPasswordlessLinkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SendMyPasswordlessLinkResponse): SendMyPasswordlessLinkResponse.AsObject;
  static serializeBinaryToWriter(message: SendMyPasswordlessLinkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendMyPasswordlessLinkResponse;
  static deserializeBinaryFromReader(message: SendMyPasswordlessLinkResponse, reader: jspb.BinaryReader): SendMyPasswordlessLinkResponse;
}

export namespace SendMyPasswordlessLinkResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class VerifyMyPasswordlessRequest extends jspb.Message {
  getVerification(): zitadel_user_pb.WebAuthNVerification | undefined;
  setVerification(value?: zitadel_user_pb.WebAuthNVerification): VerifyMyPasswordlessRequest;
  hasVerification(): boolean;
  clearVerification(): VerifyMyPasswordlessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyPasswordlessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyPasswordlessRequest): VerifyMyPasswordlessRequest.AsObject;
  static serializeBinaryToWriter(message: VerifyMyPasswordlessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyPasswordlessRequest;
  static deserializeBinaryFromReader(message: VerifyMyPasswordlessRequest, reader: jspb.BinaryReader): VerifyMyPasswordlessRequest;
}

export namespace VerifyMyPasswordlessRequest {
  export type AsObject = {
    verification?: zitadel_user_pb.WebAuthNVerification.AsObject,
  }
}

export class VerifyMyPasswordlessResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): VerifyMyPasswordlessResponse;
  hasDetails(): boolean;
  clearDetails(): VerifyMyPasswordlessResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMyPasswordlessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMyPasswordlessResponse): VerifyMyPasswordlessResponse.AsObject;
  static serializeBinaryToWriter(message: VerifyMyPasswordlessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMyPasswordlessResponse;
  static deserializeBinaryFromReader(message: VerifyMyPasswordlessResponse, reader: jspb.BinaryReader): VerifyMyPasswordlessResponse;
}

export namespace VerifyMyPasswordlessResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMyPasswordlessRequest extends jspb.Message {
  getTokenId(): string;
  setTokenId(value: string): RemoveMyPasswordlessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyPasswordlessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyPasswordlessRequest): RemoveMyPasswordlessRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMyPasswordlessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyPasswordlessRequest;
  static deserializeBinaryFromReader(message: RemoveMyPasswordlessRequest, reader: jspb.BinaryReader): RemoveMyPasswordlessRequest;
}

export namespace RemoveMyPasswordlessRequest {
  export type AsObject = {
    tokenId: string,
  }
}

export class RemoveMyPasswordlessResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMyPasswordlessResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMyPasswordlessResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMyPasswordlessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMyPasswordlessResponse): RemoveMyPasswordlessResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMyPasswordlessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMyPasswordlessResponse;
  static deserializeBinaryFromReader(message: RemoveMyPasswordlessResponse, reader: jspb.BinaryReader): RemoveMyPasswordlessResponse;
}

export namespace RemoveMyPasswordlessResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListMyUserGrantsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListMyUserGrantsRequest;
  hasQuery(): boolean;
  clearQuery(): ListMyUserGrantsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyUserGrantsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyUserGrantsRequest): ListMyUserGrantsRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyUserGrantsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyUserGrantsRequest;
  static deserializeBinaryFromReader(message: ListMyUserGrantsRequest, reader: jspb.BinaryReader): ListMyUserGrantsRequest;
}

export namespace ListMyUserGrantsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListMyUserGrantsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListMyUserGrantsResponse;
  hasDetails(): boolean;
  clearDetails(): ListMyUserGrantsResponse;

  getResultList(): Array<UserGrant>;
  setResultList(value: Array<UserGrant>): ListMyUserGrantsResponse;
  clearResultList(): ListMyUserGrantsResponse;
  addResult(value?: UserGrant, index?: number): UserGrant;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyUserGrantsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyUserGrantsResponse): ListMyUserGrantsResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyUserGrantsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyUserGrantsResponse;
  static deserializeBinaryFromReader(message: ListMyUserGrantsResponse, reader: jspb.BinaryReader): ListMyUserGrantsResponse;
}

export namespace ListMyUserGrantsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<UserGrant.AsObject>,
  }
}

export class UserGrant extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): UserGrant;

  getProjectId(): string;
  setProjectId(value: string): UserGrant;

  getUserId(): string;
  setUserId(value: string): UserGrant;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): UserGrant;
  clearRolesList(): UserGrant;
  addRoles(value: string, index?: number): UserGrant;

  getOrgName(): string;
  setOrgName(value: string): UserGrant;

  getGrantId(): string;
  setGrantId(value: string): UserGrant;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UserGrant;
  hasDetails(): boolean;
  clearDetails(): UserGrant;

  getOrgDomain(): string;
  setOrgDomain(value: string): UserGrant;

  getProjectName(): string;
  setProjectName(value: string): UserGrant;

  getProjectGrantId(): string;
  setProjectGrantId(value: string): UserGrant;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): UserGrant;
  clearRoleKeysList(): UserGrant;
  addRoleKeys(value: string, index?: number): UserGrant;

  getUserType(): zitadel_user_pb.Type;
  setUserType(value: zitadel_user_pb.Type): UserGrant;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserGrant.AsObject;
  static toObject(includeInstance: boolean, msg: UserGrant): UserGrant.AsObject;
  static serializeBinaryToWriter(message: UserGrant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserGrant;
  static deserializeBinaryFromReader(message: UserGrant, reader: jspb.BinaryReader): UserGrant;
}

export namespace UserGrant {
  export type AsObject = {
    orgId: string,
    projectId: string,
    userId: string,
    rolesList: Array<string>,
    orgName: string,
    grantId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    orgDomain: string,
    projectName: string,
    projectGrantId: string,
    roleKeysList: Array<string>,
    userType: zitadel_user_pb.Type,
  }
}

export class ListMyProjectOrgsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListMyProjectOrgsRequest;
  hasQuery(): boolean;
  clearQuery(): ListMyProjectOrgsRequest;

  getQueriesList(): Array<zitadel_org_pb.OrgQuery>;
  setQueriesList(value: Array<zitadel_org_pb.OrgQuery>): ListMyProjectOrgsRequest;
  clearQueriesList(): ListMyProjectOrgsRequest;
  addQueries(value?: zitadel_org_pb.OrgQuery, index?: number): zitadel_org_pb.OrgQuery;

  getSortingColumn(): zitadel_org_pb.OrgFieldName;
  setSortingColumn(value: zitadel_org_pb.OrgFieldName): ListMyProjectOrgsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyProjectOrgsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyProjectOrgsRequest): ListMyProjectOrgsRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyProjectOrgsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyProjectOrgsRequest;
  static deserializeBinaryFromReader(message: ListMyProjectOrgsRequest, reader: jspb.BinaryReader): ListMyProjectOrgsRequest;
}

export namespace ListMyProjectOrgsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_org_pb.OrgQuery.AsObject>,
    sortingColumn: zitadel_org_pb.OrgFieldName,
  }
}

export class ListMyProjectOrgsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListMyProjectOrgsResponse;
  hasDetails(): boolean;
  clearDetails(): ListMyProjectOrgsResponse;

  getResultList(): Array<zitadel_org_pb.Org>;
  setResultList(value: Array<zitadel_org_pb.Org>): ListMyProjectOrgsResponse;
  clearResultList(): ListMyProjectOrgsResponse;
  addResult(value?: zitadel_org_pb.Org, index?: number): zitadel_org_pb.Org;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyProjectOrgsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyProjectOrgsResponse): ListMyProjectOrgsResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyProjectOrgsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyProjectOrgsResponse;
  static deserializeBinaryFromReader(message: ListMyProjectOrgsResponse, reader: jspb.BinaryReader): ListMyProjectOrgsResponse;
}

export namespace ListMyProjectOrgsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_org_pb.Org.AsObject>,
  }
}

export class ListMyZitadelPermissionsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyZitadelPermissionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyZitadelPermissionsRequest): ListMyZitadelPermissionsRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyZitadelPermissionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyZitadelPermissionsRequest;
  static deserializeBinaryFromReader(message: ListMyZitadelPermissionsRequest, reader: jspb.BinaryReader): ListMyZitadelPermissionsRequest;
}

export namespace ListMyZitadelPermissionsRequest {
  export type AsObject = {
  }
}

export class ListMyZitadelPermissionsResponse extends jspb.Message {
  getResultList(): Array<string>;
  setResultList(value: Array<string>): ListMyZitadelPermissionsResponse;
  clearResultList(): ListMyZitadelPermissionsResponse;
  addResult(value: string, index?: number): ListMyZitadelPermissionsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyZitadelPermissionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyZitadelPermissionsResponse): ListMyZitadelPermissionsResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyZitadelPermissionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyZitadelPermissionsResponse;
  static deserializeBinaryFromReader(message: ListMyZitadelPermissionsResponse, reader: jspb.BinaryReader): ListMyZitadelPermissionsResponse;
}

export namespace ListMyZitadelPermissionsResponse {
  export type AsObject = {
    resultList: Array<string>,
  }
}

export class ListMyProjectPermissionsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyProjectPermissionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyProjectPermissionsRequest): ListMyProjectPermissionsRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyProjectPermissionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyProjectPermissionsRequest;
  static deserializeBinaryFromReader(message: ListMyProjectPermissionsRequest, reader: jspb.BinaryReader): ListMyProjectPermissionsRequest;
}

export namespace ListMyProjectPermissionsRequest {
  export type AsObject = {
  }
}

export class ListMyProjectPermissionsResponse extends jspb.Message {
  getResultList(): Array<string>;
  setResultList(value: Array<string>): ListMyProjectPermissionsResponse;
  clearResultList(): ListMyProjectPermissionsResponse;
  addResult(value: string, index?: number): ListMyProjectPermissionsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyProjectPermissionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyProjectPermissionsResponse): ListMyProjectPermissionsResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyProjectPermissionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyProjectPermissionsResponse;
  static deserializeBinaryFromReader(message: ListMyProjectPermissionsResponse, reader: jspb.BinaryReader): ListMyProjectPermissionsResponse;
}

export namespace ListMyProjectPermissionsResponse {
  export type AsObject = {
    resultList: Array<string>,
  }
}

export class ListMyMembershipsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListMyMembershipsRequest;
  hasQuery(): boolean;
  clearQuery(): ListMyMembershipsRequest;

  getQueriesList(): Array<zitadel_user_pb.MembershipQuery>;
  setQueriesList(value: Array<zitadel_user_pb.MembershipQuery>): ListMyMembershipsRequest;
  clearQueriesList(): ListMyMembershipsRequest;
  addQueries(value?: zitadel_user_pb.MembershipQuery, index?: number): zitadel_user_pb.MembershipQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyMembershipsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyMembershipsRequest): ListMyMembershipsRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyMembershipsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyMembershipsRequest;
  static deserializeBinaryFromReader(message: ListMyMembershipsRequest, reader: jspb.BinaryReader): ListMyMembershipsRequest;
}

export namespace ListMyMembershipsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_user_pb.MembershipQuery.AsObject>,
  }
}

export class ListMyMembershipsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListMyMembershipsResponse;
  hasDetails(): boolean;
  clearDetails(): ListMyMembershipsResponse;

  getResultList(): Array<zitadel_user_pb.Membership>;
  setResultList(value: Array<zitadel_user_pb.Membership>): ListMyMembershipsResponse;
  clearResultList(): ListMyMembershipsResponse;
  addResult(value?: zitadel_user_pb.Membership, index?: number): zitadel_user_pb.Membership;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyMembershipsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyMembershipsResponse): ListMyMembershipsResponse.AsObject;
  static serializeBinaryToWriter(message: ListMyMembershipsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyMembershipsResponse;
  static deserializeBinaryFromReader(message: ListMyMembershipsResponse, reader: jspb.BinaryReader): ListMyMembershipsResponse;
}

export namespace ListMyMembershipsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_user_pb.Membership.AsObject>,
  }
}

export class GetMyLabelPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyLabelPolicyRequest): GetMyLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyLabelPolicyRequest;
  static deserializeBinaryFromReader(message: GetMyLabelPolicyRequest, reader: jspb.BinaryReader): GetMyLabelPolicyRequest;
}

export namespace GetMyLabelPolicyRequest {
  export type AsObject = {
  }
}

export class GetMyLabelPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LabelPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LabelPolicy): GetMyLabelPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetMyLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyLabelPolicyResponse): GetMyLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyLabelPolicyResponse;
  static deserializeBinaryFromReader(message: GetMyLabelPolicyResponse, reader: jspb.BinaryReader): GetMyLabelPolicyResponse;
}

export namespace GetMyLabelPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LabelPolicy.AsObject,
  }
}

export class GetMyPrivacyPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyPrivacyPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyPrivacyPolicyRequest): GetMyPrivacyPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyPrivacyPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyPrivacyPolicyRequest;
  static deserializeBinaryFromReader(message: GetMyPrivacyPolicyRequest, reader: jspb.BinaryReader): GetMyPrivacyPolicyRequest;
}

export namespace GetMyPrivacyPolicyRequest {
  export type AsObject = {
  }
}

export class GetMyPrivacyPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PrivacyPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PrivacyPolicy): GetMyPrivacyPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetMyPrivacyPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyPrivacyPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyPrivacyPolicyResponse): GetMyPrivacyPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyPrivacyPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyPrivacyPolicyResponse;
  static deserializeBinaryFromReader(message: GetMyPrivacyPolicyResponse, reader: jspb.BinaryReader): GetMyPrivacyPolicyResponse;
}

export namespace GetMyPrivacyPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PrivacyPolicy.AsObject,
  }
}

export class GetMyLoginPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyLoginPolicyRequest): GetMyLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyLoginPolicyRequest;
  static deserializeBinaryFromReader(message: GetMyLoginPolicyRequest, reader: jspb.BinaryReader): GetMyLoginPolicyRequest;
}

export namespace GetMyLoginPolicyRequest {
  export type AsObject = {
  }
}

export class GetMyLoginPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LoginPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LoginPolicy): GetMyLoginPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetMyLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyLoginPolicyResponse): GetMyLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyLoginPolicyResponse;
  static deserializeBinaryFromReader(message: GetMyLoginPolicyResponse, reader: jspb.BinaryReader): GetMyLoginPolicyResponse;
}

export namespace GetMyLoginPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LoginPolicy.AsObject,
  }
}


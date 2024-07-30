import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_feature_v2beta_feature_pb from '../../../zitadel/feature/v2beta/feature_pb'; // proto import: "zitadel/feature/v2beta/feature.proto"


export class SetInstanceFeaturesRequest extends jspb.Message {
  getLoginDefaultOrg(): boolean;
  setLoginDefaultOrg(value: boolean): SetInstanceFeaturesRequest;
  hasLoginDefaultOrg(): boolean;
  clearLoginDefaultOrg(): SetInstanceFeaturesRequest;

  getOidcTriggerIntrospectionProjections(): boolean;
  setOidcTriggerIntrospectionProjections(value: boolean): SetInstanceFeaturesRequest;
  hasOidcTriggerIntrospectionProjections(): boolean;
  clearOidcTriggerIntrospectionProjections(): SetInstanceFeaturesRequest;

  getOidcLegacyIntrospection(): boolean;
  setOidcLegacyIntrospection(value: boolean): SetInstanceFeaturesRequest;
  hasOidcLegacyIntrospection(): boolean;
  clearOidcLegacyIntrospection(): SetInstanceFeaturesRequest;

  getUserSchema(): boolean;
  setUserSchema(value: boolean): SetInstanceFeaturesRequest;
  hasUserSchema(): boolean;
  clearUserSchema(): SetInstanceFeaturesRequest;

  getOidcTokenExchange(): boolean;
  setOidcTokenExchange(value: boolean): SetInstanceFeaturesRequest;
  hasOidcTokenExchange(): boolean;
  clearOidcTokenExchange(): SetInstanceFeaturesRequest;

  getActions(): boolean;
  setActions(value: boolean): SetInstanceFeaturesRequest;
  hasActions(): boolean;
  clearActions(): SetInstanceFeaturesRequest;

  getImprovedPerformanceList(): Array<zitadel_feature_v2beta_feature_pb.ImprovedPerformance>;
  setImprovedPerformanceList(value: Array<zitadel_feature_v2beta_feature_pb.ImprovedPerformance>): SetInstanceFeaturesRequest;
  clearImprovedPerformanceList(): SetInstanceFeaturesRequest;
  addImprovedPerformance(value: zitadel_feature_v2beta_feature_pb.ImprovedPerformance, index?: number): SetInstanceFeaturesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetInstanceFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetInstanceFeaturesRequest): SetInstanceFeaturesRequest.AsObject;
  static serializeBinaryToWriter(message: SetInstanceFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetInstanceFeaturesRequest;
  static deserializeBinaryFromReader(message: SetInstanceFeaturesRequest, reader: jspb.BinaryReader): SetInstanceFeaturesRequest;
}

export namespace SetInstanceFeaturesRequest {
  export type AsObject = {
    loginDefaultOrg?: boolean,
    oidcTriggerIntrospectionProjections?: boolean,
    oidcLegacyIntrospection?: boolean,
    userSchema?: boolean,
    oidcTokenExchange?: boolean,
    actions?: boolean,
    improvedPerformanceList: Array<zitadel_feature_v2beta_feature_pb.ImprovedPerformance>,
  }

  export enum LoginDefaultOrgCase { 
    _LOGIN_DEFAULT_ORG_NOT_SET = 0,
    LOGIN_DEFAULT_ORG = 1,
  }

  export enum OidcTriggerIntrospectionProjectionsCase { 
    _OIDC_TRIGGER_INTROSPECTION_PROJECTIONS_NOT_SET = 0,
    OIDC_TRIGGER_INTROSPECTION_PROJECTIONS = 2,
  }

  export enum OidcLegacyIntrospectionCase { 
    _OIDC_LEGACY_INTROSPECTION_NOT_SET = 0,
    OIDC_LEGACY_INTROSPECTION = 3,
  }

  export enum UserSchemaCase { 
    _USER_SCHEMA_NOT_SET = 0,
    USER_SCHEMA = 4,
  }

  export enum OidcTokenExchangeCase { 
    _OIDC_TOKEN_EXCHANGE_NOT_SET = 0,
    OIDC_TOKEN_EXCHANGE = 5,
  }

  export enum ActionsCase { 
    _ACTIONS_NOT_SET = 0,
    ACTIONS = 6,
  }
}

export class SetInstanceFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetInstanceFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): SetInstanceFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetInstanceFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetInstanceFeaturesResponse): SetInstanceFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: SetInstanceFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetInstanceFeaturesResponse;
  static deserializeBinaryFromReader(message: SetInstanceFeaturesResponse, reader: jspb.BinaryReader): SetInstanceFeaturesResponse;
}

export namespace SetInstanceFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ResetInstanceFeaturesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetInstanceFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetInstanceFeaturesRequest): ResetInstanceFeaturesRequest.AsObject;
  static serializeBinaryToWriter(message: ResetInstanceFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetInstanceFeaturesRequest;
  static deserializeBinaryFromReader(message: ResetInstanceFeaturesRequest, reader: jspb.BinaryReader): ResetInstanceFeaturesRequest;
}

export namespace ResetInstanceFeaturesRequest {
  export type AsObject = {
  }
}

export class ResetInstanceFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): ResetInstanceFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): ResetInstanceFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetInstanceFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetInstanceFeaturesResponse): ResetInstanceFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: ResetInstanceFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetInstanceFeaturesResponse;
  static deserializeBinaryFromReader(message: ResetInstanceFeaturesResponse, reader: jspb.BinaryReader): ResetInstanceFeaturesResponse;
}

export namespace ResetInstanceFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class GetInstanceFeaturesRequest extends jspb.Message {
  getInheritance(): boolean;
  setInheritance(value: boolean): GetInstanceFeaturesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInstanceFeaturesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetInstanceFeaturesRequest): GetInstanceFeaturesRequest.AsObject;
  static serializeBinaryToWriter(message: GetInstanceFeaturesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInstanceFeaturesRequest;
  static deserializeBinaryFromReader(message: GetInstanceFeaturesRequest, reader: jspb.BinaryReader): GetInstanceFeaturesRequest;
}

export namespace GetInstanceFeaturesRequest {
  export type AsObject = {
    inheritance: boolean,
  }
}

export class GetInstanceFeaturesResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): GetInstanceFeaturesResponse;
  hasDetails(): boolean;
  clearDetails(): GetInstanceFeaturesResponse;

  getLoginDefaultOrg(): zitadel_feature_v2beta_feature_pb.FeatureFlag | undefined;
  setLoginDefaultOrg(value?: zitadel_feature_v2beta_feature_pb.FeatureFlag): GetInstanceFeaturesResponse;
  hasLoginDefaultOrg(): boolean;
  clearLoginDefaultOrg(): GetInstanceFeaturesResponse;

  getOidcTriggerIntrospectionProjections(): zitadel_feature_v2beta_feature_pb.FeatureFlag | undefined;
  setOidcTriggerIntrospectionProjections(value?: zitadel_feature_v2beta_feature_pb.FeatureFlag): GetInstanceFeaturesResponse;
  hasOidcTriggerIntrospectionProjections(): boolean;
  clearOidcTriggerIntrospectionProjections(): GetInstanceFeaturesResponse;

  getOidcLegacyIntrospection(): zitadel_feature_v2beta_feature_pb.FeatureFlag | undefined;
  setOidcLegacyIntrospection(value?: zitadel_feature_v2beta_feature_pb.FeatureFlag): GetInstanceFeaturesResponse;
  hasOidcLegacyIntrospection(): boolean;
  clearOidcLegacyIntrospection(): GetInstanceFeaturesResponse;

  getUserSchema(): zitadel_feature_v2beta_feature_pb.FeatureFlag | undefined;
  setUserSchema(value?: zitadel_feature_v2beta_feature_pb.FeatureFlag): GetInstanceFeaturesResponse;
  hasUserSchema(): boolean;
  clearUserSchema(): GetInstanceFeaturesResponse;

  getOidcTokenExchange(): zitadel_feature_v2beta_feature_pb.FeatureFlag | undefined;
  setOidcTokenExchange(value?: zitadel_feature_v2beta_feature_pb.FeatureFlag): GetInstanceFeaturesResponse;
  hasOidcTokenExchange(): boolean;
  clearOidcTokenExchange(): GetInstanceFeaturesResponse;

  getActions(): zitadel_feature_v2beta_feature_pb.FeatureFlag | undefined;
  setActions(value?: zitadel_feature_v2beta_feature_pb.FeatureFlag): GetInstanceFeaturesResponse;
  hasActions(): boolean;
  clearActions(): GetInstanceFeaturesResponse;

  getImprovedPerformance(): zitadel_feature_v2beta_feature_pb.ImprovedPerformanceFeatureFlag | undefined;
  setImprovedPerformance(value?: zitadel_feature_v2beta_feature_pb.ImprovedPerformanceFeatureFlag): GetInstanceFeaturesResponse;
  hasImprovedPerformance(): boolean;
  clearImprovedPerformance(): GetInstanceFeaturesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInstanceFeaturesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetInstanceFeaturesResponse): GetInstanceFeaturesResponse.AsObject;
  static serializeBinaryToWriter(message: GetInstanceFeaturesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInstanceFeaturesResponse;
  static deserializeBinaryFromReader(message: GetInstanceFeaturesResponse, reader: jspb.BinaryReader): GetInstanceFeaturesResponse;
}

export namespace GetInstanceFeaturesResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    loginDefaultOrg?: zitadel_feature_v2beta_feature_pb.FeatureFlag.AsObject,
    oidcTriggerIntrospectionProjections?: zitadel_feature_v2beta_feature_pb.FeatureFlag.AsObject,
    oidcLegacyIntrospection?: zitadel_feature_v2beta_feature_pb.FeatureFlag.AsObject,
    userSchema?: zitadel_feature_v2beta_feature_pb.FeatureFlag.AsObject,
    oidcTokenExchange?: zitadel_feature_v2beta_feature_pb.FeatureFlag.AsObject,
    actions?: zitadel_feature_v2beta_feature_pb.FeatureFlag.AsObject,
    improvedPerformance?: zitadel_feature_v2beta_feature_pb.ImprovedPerformanceFeatureFlag.AsObject,
  }
}


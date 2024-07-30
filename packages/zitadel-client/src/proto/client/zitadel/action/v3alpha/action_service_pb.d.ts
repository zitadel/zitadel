import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_action_v3alpha_target_pb from '../../../zitadel/action/v3alpha/target_pb'; // proto import: "zitadel/action/v3alpha/target.proto"
import * as zitadel_action_v3alpha_execution_pb from '../../../zitadel/action/v3alpha/execution_pb'; // proto import: "zitadel/action/v3alpha/execution.proto"
import * as zitadel_action_v3alpha_query_pb from '../../../zitadel/action/v3alpha/query_pb'; // proto import: "zitadel/action/v3alpha/query.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_protoc_gen_zitadel_v2_options_pb from '../../../zitadel/protoc_gen_zitadel/v2/options_pb'; // proto import: "zitadel/protoc_gen_zitadel/v2/options.proto"


export class CreateTargetRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateTargetRequest;

  getRestWebhook(): zitadel_action_v3alpha_target_pb.SetRESTWebhook | undefined;
  setRestWebhook(value?: zitadel_action_v3alpha_target_pb.SetRESTWebhook): CreateTargetRequest;
  hasRestWebhook(): boolean;
  clearRestWebhook(): CreateTargetRequest;

  getRestCall(): zitadel_action_v3alpha_target_pb.SetRESTCall | undefined;
  setRestCall(value?: zitadel_action_v3alpha_target_pb.SetRESTCall): CreateTargetRequest;
  hasRestCall(): boolean;
  clearRestCall(): CreateTargetRequest;

  getRestAsync(): zitadel_action_v3alpha_target_pb.SetRESTAsync | undefined;
  setRestAsync(value?: zitadel_action_v3alpha_target_pb.SetRESTAsync): CreateTargetRequest;
  hasRestAsync(): boolean;
  clearRestAsync(): CreateTargetRequest;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): CreateTargetRequest;
  hasTimeout(): boolean;
  clearTimeout(): CreateTargetRequest;

  getEndpoint(): string;
  setEndpoint(value: string): CreateTargetRequest;

  getTargetTypeCase(): CreateTargetRequest.TargetTypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTargetRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateTargetRequest): CreateTargetRequest.AsObject;
  static serializeBinaryToWriter(message: CreateTargetRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateTargetRequest;
  static deserializeBinaryFromReader(message: CreateTargetRequest, reader: jspb.BinaryReader): CreateTargetRequest;
}

export namespace CreateTargetRequest {
  export type AsObject = {
    name: string,
    restWebhook?: zitadel_action_v3alpha_target_pb.SetRESTWebhook.AsObject,
    restCall?: zitadel_action_v3alpha_target_pb.SetRESTCall.AsObject,
    restAsync?: zitadel_action_v3alpha_target_pb.SetRESTAsync.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    endpoint: string,
  }

  export enum TargetTypeCase { 
    TARGET_TYPE_NOT_SET = 0,
    REST_WEBHOOK = 2,
    REST_CALL = 3,
    REST_ASYNC = 4,
  }
}

export class CreateTargetResponse extends jspb.Message {
  getId(): string;
  setId(value: string): CreateTargetResponse;

  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): CreateTargetResponse;
  hasDetails(): boolean;
  clearDetails(): CreateTargetResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateTargetResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateTargetResponse): CreateTargetResponse.AsObject;
  static serializeBinaryToWriter(message: CreateTargetResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateTargetResponse;
  static deserializeBinaryFromReader(message: CreateTargetResponse, reader: jspb.BinaryReader): CreateTargetResponse;
}

export namespace CreateTargetResponse {
  export type AsObject = {
    id: string,
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class UpdateTargetRequest extends jspb.Message {
  getTargetId(): string;
  setTargetId(value: string): UpdateTargetRequest;

  getName(): string;
  setName(value: string): UpdateTargetRequest;
  hasName(): boolean;
  clearName(): UpdateTargetRequest;

  getRestWebhook(): zitadel_action_v3alpha_target_pb.SetRESTWebhook | undefined;
  setRestWebhook(value?: zitadel_action_v3alpha_target_pb.SetRESTWebhook): UpdateTargetRequest;
  hasRestWebhook(): boolean;
  clearRestWebhook(): UpdateTargetRequest;

  getRestCall(): zitadel_action_v3alpha_target_pb.SetRESTCall | undefined;
  setRestCall(value?: zitadel_action_v3alpha_target_pb.SetRESTCall): UpdateTargetRequest;
  hasRestCall(): boolean;
  clearRestCall(): UpdateTargetRequest;

  getRestAsync(): zitadel_action_v3alpha_target_pb.SetRESTAsync | undefined;
  setRestAsync(value?: zitadel_action_v3alpha_target_pb.SetRESTAsync): UpdateTargetRequest;
  hasRestAsync(): boolean;
  clearRestAsync(): UpdateTargetRequest;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): UpdateTargetRequest;
  hasTimeout(): boolean;
  clearTimeout(): UpdateTargetRequest;

  getEndpoint(): string;
  setEndpoint(value: string): UpdateTargetRequest;
  hasEndpoint(): boolean;
  clearEndpoint(): UpdateTargetRequest;

  getTargetTypeCase(): UpdateTargetRequest.TargetTypeCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateTargetRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateTargetRequest): UpdateTargetRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateTargetRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateTargetRequest;
  static deserializeBinaryFromReader(message: UpdateTargetRequest, reader: jspb.BinaryReader): UpdateTargetRequest;
}

export namespace UpdateTargetRequest {
  export type AsObject = {
    targetId: string,
    name?: string,
    restWebhook?: zitadel_action_v3alpha_target_pb.SetRESTWebhook.AsObject,
    restCall?: zitadel_action_v3alpha_target_pb.SetRESTCall.AsObject,
    restAsync?: zitadel_action_v3alpha_target_pb.SetRESTAsync.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    endpoint?: string,
  }

  export enum TargetTypeCase { 
    TARGET_TYPE_NOT_SET = 0,
    REST_WEBHOOK = 3,
    REST_CALL = 4,
    REST_ASYNC = 5,
  }

  export enum NameCase { 
    _NAME_NOT_SET = 0,
    NAME = 2,
  }

  export enum TimeoutCase { 
    _TIMEOUT_NOT_SET = 0,
    TIMEOUT = 6,
  }

  export enum EndpointCase { 
    _ENDPOINT_NOT_SET = 0,
    ENDPOINT = 7,
  }
}

export class UpdateTargetResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): UpdateTargetResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateTargetResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateTargetResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateTargetResponse): UpdateTargetResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateTargetResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateTargetResponse;
  static deserializeBinaryFromReader(message: UpdateTargetResponse, reader: jspb.BinaryReader): UpdateTargetResponse;
}

export namespace UpdateTargetResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class DeleteTargetRequest extends jspb.Message {
  getTargetId(): string;
  setTargetId(value: string): DeleteTargetRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteTargetRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteTargetRequest): DeleteTargetRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteTargetRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteTargetRequest;
  static deserializeBinaryFromReader(message: DeleteTargetRequest, reader: jspb.BinaryReader): DeleteTargetRequest;
}

export namespace DeleteTargetRequest {
  export type AsObject = {
    targetId: string,
  }
}

export class DeleteTargetResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): DeleteTargetResponse;
  hasDetails(): boolean;
  clearDetails(): DeleteTargetResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteTargetResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteTargetResponse): DeleteTargetResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteTargetResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteTargetResponse;
  static deserializeBinaryFromReader(message: DeleteTargetResponse, reader: jspb.BinaryReader): DeleteTargetResponse;
}

export namespace DeleteTargetResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ListTargetsRequest extends jspb.Message {
  getQuery(): zitadel_object_v2beta_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_v2beta_object_pb.ListQuery): ListTargetsRequest;
  hasQuery(): boolean;
  clearQuery(): ListTargetsRequest;

  getSortingColumn(): zitadel_action_v3alpha_query_pb.TargetFieldName;
  setSortingColumn(value: zitadel_action_v3alpha_query_pb.TargetFieldName): ListTargetsRequest;

  getQueriesList(): Array<zitadel_action_v3alpha_query_pb.TargetSearchQuery>;
  setQueriesList(value: Array<zitadel_action_v3alpha_query_pb.TargetSearchQuery>): ListTargetsRequest;
  clearQueriesList(): ListTargetsRequest;
  addQueries(value?: zitadel_action_v3alpha_query_pb.TargetSearchQuery, index?: number): zitadel_action_v3alpha_query_pb.TargetSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTargetsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListTargetsRequest): ListTargetsRequest.AsObject;
  static serializeBinaryToWriter(message: ListTargetsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTargetsRequest;
  static deserializeBinaryFromReader(message: ListTargetsRequest, reader: jspb.BinaryReader): ListTargetsRequest;
}

export namespace ListTargetsRequest {
  export type AsObject = {
    query?: zitadel_object_v2beta_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_action_v3alpha_query_pb.TargetFieldName,
    queriesList: Array<zitadel_action_v3alpha_query_pb.TargetSearchQuery.AsObject>,
  }
}

export class ListTargetsResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.ListDetails): ListTargetsResponse;
  hasDetails(): boolean;
  clearDetails(): ListTargetsResponse;

  getSortingColumn(): zitadel_action_v3alpha_query_pb.TargetFieldName;
  setSortingColumn(value: zitadel_action_v3alpha_query_pb.TargetFieldName): ListTargetsResponse;

  getResultList(): Array<zitadel_action_v3alpha_target_pb.Target>;
  setResultList(value: Array<zitadel_action_v3alpha_target_pb.Target>): ListTargetsResponse;
  clearResultList(): ListTargetsResponse;
  addResult(value?: zitadel_action_v3alpha_target_pb.Target, index?: number): zitadel_action_v3alpha_target_pb.Target;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListTargetsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListTargetsResponse): ListTargetsResponse.AsObject;
  static serializeBinaryToWriter(message: ListTargetsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListTargetsResponse;
  static deserializeBinaryFromReader(message: ListTargetsResponse, reader: jspb.BinaryReader): ListTargetsResponse;
}

export namespace ListTargetsResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_action_v3alpha_query_pb.TargetFieldName,
    resultList: Array<zitadel_action_v3alpha_target_pb.Target.AsObject>,
  }
}

export class GetTargetByIDRequest extends jspb.Message {
  getTargetId(): string;
  setTargetId(value: string): GetTargetByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTargetByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetTargetByIDRequest): GetTargetByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetTargetByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTargetByIDRequest;
  static deserializeBinaryFromReader(message: GetTargetByIDRequest, reader: jspb.BinaryReader): GetTargetByIDRequest;
}

export namespace GetTargetByIDRequest {
  export type AsObject = {
    targetId: string,
  }
}

export class GetTargetByIDResponse extends jspb.Message {
  getTarget(): zitadel_action_v3alpha_target_pb.Target | undefined;
  setTarget(value?: zitadel_action_v3alpha_target_pb.Target): GetTargetByIDResponse;
  hasTarget(): boolean;
  clearTarget(): GetTargetByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTargetByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetTargetByIDResponse): GetTargetByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetTargetByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTargetByIDResponse;
  static deserializeBinaryFromReader(message: GetTargetByIDResponse, reader: jspb.BinaryReader): GetTargetByIDResponse;
}

export namespace GetTargetByIDResponse {
  export type AsObject = {
    target?: zitadel_action_v3alpha_target_pb.Target.AsObject,
  }
}

export class SetExecutionRequest extends jspb.Message {
  getCondition(): zitadel_action_v3alpha_execution_pb.Condition | undefined;
  setCondition(value?: zitadel_action_v3alpha_execution_pb.Condition): SetExecutionRequest;
  hasCondition(): boolean;
  clearCondition(): SetExecutionRequest;

  getTargetsList(): Array<zitadel_action_v3alpha_execution_pb.ExecutionTargetType>;
  setTargetsList(value: Array<zitadel_action_v3alpha_execution_pb.ExecutionTargetType>): SetExecutionRequest;
  clearTargetsList(): SetExecutionRequest;
  addTargets(value?: zitadel_action_v3alpha_execution_pb.ExecutionTargetType, index?: number): zitadel_action_v3alpha_execution_pb.ExecutionTargetType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetExecutionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetExecutionRequest): SetExecutionRequest.AsObject;
  static serializeBinaryToWriter(message: SetExecutionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetExecutionRequest;
  static deserializeBinaryFromReader(message: SetExecutionRequest, reader: jspb.BinaryReader): SetExecutionRequest;
}

export namespace SetExecutionRequest {
  export type AsObject = {
    condition?: zitadel_action_v3alpha_execution_pb.Condition.AsObject,
    targetsList: Array<zitadel_action_v3alpha_execution_pb.ExecutionTargetType.AsObject>,
  }
}

export class SetExecutionResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): SetExecutionResponse;
  hasDetails(): boolean;
  clearDetails(): SetExecutionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetExecutionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetExecutionResponse): SetExecutionResponse.AsObject;
  static serializeBinaryToWriter(message: SetExecutionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetExecutionResponse;
  static deserializeBinaryFromReader(message: SetExecutionResponse, reader: jspb.BinaryReader): SetExecutionResponse;
}

export namespace SetExecutionResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class DeleteExecutionRequest extends jspb.Message {
  getCondition(): zitadel_action_v3alpha_execution_pb.Condition | undefined;
  setCondition(value?: zitadel_action_v3alpha_execution_pb.Condition): DeleteExecutionRequest;
  hasCondition(): boolean;
  clearCondition(): DeleteExecutionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteExecutionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteExecutionRequest): DeleteExecutionRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteExecutionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteExecutionRequest;
  static deserializeBinaryFromReader(message: DeleteExecutionRequest, reader: jspb.BinaryReader): DeleteExecutionRequest;
}

export namespace DeleteExecutionRequest {
  export type AsObject = {
    condition?: zitadel_action_v3alpha_execution_pb.Condition.AsObject,
  }
}

export class DeleteExecutionResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): DeleteExecutionResponse;
  hasDetails(): boolean;
  clearDetails(): DeleteExecutionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteExecutionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteExecutionResponse): DeleteExecutionResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteExecutionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteExecutionResponse;
  static deserializeBinaryFromReader(message: DeleteExecutionResponse, reader: jspb.BinaryReader): DeleteExecutionResponse;
}

export namespace DeleteExecutionResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
  }
}

export class ListExecutionsRequest extends jspb.Message {
  getQuery(): zitadel_object_v2beta_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_v2beta_object_pb.ListQuery): ListExecutionsRequest;
  hasQuery(): boolean;
  clearQuery(): ListExecutionsRequest;

  getQueriesList(): Array<zitadel_action_v3alpha_query_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_action_v3alpha_query_pb.SearchQuery>): ListExecutionsRequest;
  clearQueriesList(): ListExecutionsRequest;
  addQueries(value?: zitadel_action_v3alpha_query_pb.SearchQuery, index?: number): zitadel_action_v3alpha_query_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExecutionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListExecutionsRequest): ListExecutionsRequest.AsObject;
  static serializeBinaryToWriter(message: ListExecutionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExecutionsRequest;
  static deserializeBinaryFromReader(message: ListExecutionsRequest, reader: jspb.BinaryReader): ListExecutionsRequest;
}

export namespace ListExecutionsRequest {
  export type AsObject = {
    query?: zitadel_object_v2beta_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_action_v3alpha_query_pb.SearchQuery.AsObject>,
  }
}

export class ListExecutionsResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.ListDetails): ListExecutionsResponse;
  hasDetails(): boolean;
  clearDetails(): ListExecutionsResponse;

  getResultList(): Array<zitadel_action_v3alpha_execution_pb.Execution>;
  setResultList(value: Array<zitadel_action_v3alpha_execution_pb.Execution>): ListExecutionsResponse;
  clearResultList(): ListExecutionsResponse;
  addResult(value?: zitadel_action_v3alpha_execution_pb.Execution, index?: number): zitadel_action_v3alpha_execution_pb.Execution;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExecutionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListExecutionsResponse): ListExecutionsResponse.AsObject;
  static serializeBinaryToWriter(message: ListExecutionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExecutionsResponse;
  static deserializeBinaryFromReader(message: ListExecutionsResponse, reader: jspb.BinaryReader): ListExecutionsResponse;
}

export namespace ListExecutionsResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_action_v3alpha_execution_pb.Execution.AsObject>,
  }
}

export class ListExecutionFunctionsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExecutionFunctionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListExecutionFunctionsRequest): ListExecutionFunctionsRequest.AsObject;
  static serializeBinaryToWriter(message: ListExecutionFunctionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExecutionFunctionsRequest;
  static deserializeBinaryFromReader(message: ListExecutionFunctionsRequest, reader: jspb.BinaryReader): ListExecutionFunctionsRequest;
}

export namespace ListExecutionFunctionsRequest {
  export type AsObject = {
  }
}

export class ListExecutionFunctionsResponse extends jspb.Message {
  getFunctionsList(): Array<string>;
  setFunctionsList(value: Array<string>): ListExecutionFunctionsResponse;
  clearFunctionsList(): ListExecutionFunctionsResponse;
  addFunctions(value: string, index?: number): ListExecutionFunctionsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExecutionFunctionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListExecutionFunctionsResponse): ListExecutionFunctionsResponse.AsObject;
  static serializeBinaryToWriter(message: ListExecutionFunctionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExecutionFunctionsResponse;
  static deserializeBinaryFromReader(message: ListExecutionFunctionsResponse, reader: jspb.BinaryReader): ListExecutionFunctionsResponse;
}

export namespace ListExecutionFunctionsResponse {
  export type AsObject = {
    functionsList: Array<string>,
  }
}

export class ListExecutionMethodsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExecutionMethodsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListExecutionMethodsRequest): ListExecutionMethodsRequest.AsObject;
  static serializeBinaryToWriter(message: ListExecutionMethodsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExecutionMethodsRequest;
  static deserializeBinaryFromReader(message: ListExecutionMethodsRequest, reader: jspb.BinaryReader): ListExecutionMethodsRequest;
}

export namespace ListExecutionMethodsRequest {
  export type AsObject = {
  }
}

export class ListExecutionMethodsResponse extends jspb.Message {
  getMethodsList(): Array<string>;
  setMethodsList(value: Array<string>): ListExecutionMethodsResponse;
  clearMethodsList(): ListExecutionMethodsResponse;
  addMethods(value: string, index?: number): ListExecutionMethodsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExecutionMethodsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListExecutionMethodsResponse): ListExecutionMethodsResponse.AsObject;
  static serializeBinaryToWriter(message: ListExecutionMethodsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExecutionMethodsResponse;
  static deserializeBinaryFromReader(message: ListExecutionMethodsResponse, reader: jspb.BinaryReader): ListExecutionMethodsResponse;
}

export namespace ListExecutionMethodsResponse {
  export type AsObject = {
    methodsList: Array<string>,
  }
}

export class ListExecutionServicesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExecutionServicesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListExecutionServicesRequest): ListExecutionServicesRequest.AsObject;
  static serializeBinaryToWriter(message: ListExecutionServicesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExecutionServicesRequest;
  static deserializeBinaryFromReader(message: ListExecutionServicesRequest, reader: jspb.BinaryReader): ListExecutionServicesRequest;
}

export namespace ListExecutionServicesRequest {
  export type AsObject = {
  }
}

export class ListExecutionServicesResponse extends jspb.Message {
  getServicesList(): Array<string>;
  setServicesList(value: Array<string>): ListExecutionServicesResponse;
  clearServicesList(): ListExecutionServicesResponse;
  addServices(value: string, index?: number): ListExecutionServicesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListExecutionServicesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListExecutionServicesResponse): ListExecutionServicesResponse.AsObject;
  static serializeBinaryToWriter(message: ListExecutionServicesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListExecutionServicesResponse;
  static deserializeBinaryFromReader(message: ListExecutionServicesResponse, reader: jspb.BinaryReader): ListExecutionServicesResponse;
}

export namespace ListExecutionServicesResponse {
  export type AsObject = {
    servicesList: Array<string>,
  }
}


import * as grpcWeb from 'grpc-web';

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
import * as authoption_options_pb from './authoption/options_pb';

import {
  Org,
  OrgID,
  OrgSearchRequest,
  OrgSearchResponse,
  OrgSetUpRequest,
  OrgSetUpResponse,
  UniqueOrgRequest,
  UniqueOrgResponse} from './admin_pb';

export class AdminServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  healthz(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  ready(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  validate(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_struct_pb.Struct) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_struct_pb.Struct>;

  isOrgUnique(
    request: UniqueOrgRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UniqueOrgResponse) => void
  ): grpcWeb.ClientReadableStream<UniqueOrgResponse>;

  getOrgByID(
    request: OrgID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Org) => void
  ): grpcWeb.ClientReadableStream<Org>;

  searchOrgs(
    request: OrgSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgSearchResponse) => void
  ): grpcWeb.ClientReadableStream<OrgSearchResponse>;

  setUpOrg(
    request: OrgSetUpRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgSetUpResponse) => void
  ): grpcWeb.ClientReadableStream<OrgSetUpResponse>;

}

export class AdminServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  healthz(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  ready(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  validate(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_struct_pb.Struct>;

  isOrgUnique(
    request: UniqueOrgRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UniqueOrgResponse>;

  getOrgByID(
    request: OrgID,
    metadata?: grpcWeb.Metadata
  ): Promise<Org>;

  searchOrgs(
    request: OrgSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgSearchResponse>;

  setUpOrg(
    request: OrgSetUpRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgSetUpResponse>;

}


import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class FeatureFlag extends jspb.Message {
  getEnabled(): boolean;
  setEnabled(value: boolean): FeatureFlag;

  getSource(): Source;
  setSource(value: Source): FeatureFlag;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FeatureFlag.AsObject;
  static toObject(includeInstance: boolean, msg: FeatureFlag): FeatureFlag.AsObject;
  static serializeBinaryToWriter(message: FeatureFlag, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FeatureFlag;
  static deserializeBinaryFromReader(message: FeatureFlag, reader: jspb.BinaryReader): FeatureFlag;
}

export namespace FeatureFlag {
  export type AsObject = {
    enabled: boolean,
    source: Source,
  }
}

export class ImprovedPerformanceFeatureFlag extends jspb.Message {
  getExecutionPathsList(): Array<ImprovedPerformance>;
  setExecutionPathsList(value: Array<ImprovedPerformance>): ImprovedPerformanceFeatureFlag;
  clearExecutionPathsList(): ImprovedPerformanceFeatureFlag;
  addExecutionPaths(value: ImprovedPerformance, index?: number): ImprovedPerformanceFeatureFlag;

  getSource(): Source;
  setSource(value: Source): ImprovedPerformanceFeatureFlag;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImprovedPerformanceFeatureFlag.AsObject;
  static toObject(includeInstance: boolean, msg: ImprovedPerformanceFeatureFlag): ImprovedPerformanceFeatureFlag.AsObject;
  static serializeBinaryToWriter(message: ImprovedPerformanceFeatureFlag, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImprovedPerformanceFeatureFlag;
  static deserializeBinaryFromReader(message: ImprovedPerformanceFeatureFlag, reader: jspb.BinaryReader): ImprovedPerformanceFeatureFlag;
}

export namespace ImprovedPerformanceFeatureFlag {
  export type AsObject = {
    executionPathsList: Array<ImprovedPerformance>,
    source: Source,
  }
}

export enum Source { 
  SOURCE_UNSPECIFIED = 0,
  SOURCE_SYSTEM = 2,
  SOURCE_INSTANCE = 3,
  SOURCE_ORGANIZATION = 4,
  SOURCE_PROJECT = 5,
  SOURCE_APP = 6,
  SOURCE_USER = 7,
}
export enum ImprovedPerformance { 
  IMPROVED_PERFORMANCE_UNSPECIFIED = 0,
  IMPROVED_PERFORMANCE_ORG_BY_ID = 1,
  IMPROVED_PERFORMANCE_PROJECT_GRANT = 2,
  IMPROVED_PERFORMANCE_PROJECT = 3,
  IMPROVED_PERFORMANCE_USER_GRANT = 4,
  IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED = 5,
}

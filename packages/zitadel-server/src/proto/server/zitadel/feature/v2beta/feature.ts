/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.feature.v2beta";

export enum Source {
  SOURCE_UNSPECIFIED = 0,
  SOURCE_SYSTEM = 2,
  SOURCE_INSTANCE = 3,
  SOURCE_ORGANIZATION = 4,
  /** SOURCE_PROJECT - reserved for future use */
  SOURCE_PROJECT = 5,
  /** SOURCE_APP - reserved for future use */
  SOURCE_APP = 6,
  SOURCE_USER = 7,
  UNRECOGNIZED = -1,
}

export function sourceFromJSON(object: any): Source {
  switch (object) {
    case 0:
    case "SOURCE_UNSPECIFIED":
      return Source.SOURCE_UNSPECIFIED;
    case 2:
    case "SOURCE_SYSTEM":
      return Source.SOURCE_SYSTEM;
    case 3:
    case "SOURCE_INSTANCE":
      return Source.SOURCE_INSTANCE;
    case 4:
    case "SOURCE_ORGANIZATION":
      return Source.SOURCE_ORGANIZATION;
    case 5:
    case "SOURCE_PROJECT":
      return Source.SOURCE_PROJECT;
    case 6:
    case "SOURCE_APP":
      return Source.SOURCE_APP;
    case 7:
    case "SOURCE_USER":
      return Source.SOURCE_USER;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Source.UNRECOGNIZED;
  }
}

export function sourceToJSON(object: Source): string {
  switch (object) {
    case Source.SOURCE_UNSPECIFIED:
      return "SOURCE_UNSPECIFIED";
    case Source.SOURCE_SYSTEM:
      return "SOURCE_SYSTEM";
    case Source.SOURCE_INSTANCE:
      return "SOURCE_INSTANCE";
    case Source.SOURCE_ORGANIZATION:
      return "SOURCE_ORGANIZATION";
    case Source.SOURCE_PROJECT:
      return "SOURCE_PROJECT";
    case Source.SOURCE_APP:
      return "SOURCE_APP";
    case Source.SOURCE_USER:
      return "SOURCE_USER";
    case Source.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ImprovedPerformance {
  IMPROVED_PERFORMANCE_UNSPECIFIED = 0,
  /**
   * IMPROVED_PERFORMANCE_ORG_BY_ID - Uses the eventstore to query the org by id
   * instead of the sql table.
   */
  IMPROVED_PERFORMANCE_ORG_BY_ID = 1,
  /**
   * IMPROVED_PERFORMANCE_PROJECT_GRANT - Improves performance on write side by using
   * optimized processes to query data to determine
   * correctnes of data.
   */
  IMPROVED_PERFORMANCE_PROJECT_GRANT = 2,
  IMPROVED_PERFORMANCE_PROJECT = 3,
  IMPROVED_PERFORMANCE_USER_GRANT = 4,
  /**
   * IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED - Improve performance on write side when
   * users are checked against verified domains
   * from other organizations.
   */
  IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED = 5,
  UNRECOGNIZED = -1,
}

export function improvedPerformanceFromJSON(object: any): ImprovedPerformance {
  switch (object) {
    case 0:
    case "IMPROVED_PERFORMANCE_UNSPECIFIED":
      return ImprovedPerformance.IMPROVED_PERFORMANCE_UNSPECIFIED;
    case 1:
    case "IMPROVED_PERFORMANCE_ORG_BY_ID":
      return ImprovedPerformance.IMPROVED_PERFORMANCE_ORG_BY_ID;
    case 2:
    case "IMPROVED_PERFORMANCE_PROJECT_GRANT":
      return ImprovedPerformance.IMPROVED_PERFORMANCE_PROJECT_GRANT;
    case 3:
    case "IMPROVED_PERFORMANCE_PROJECT":
      return ImprovedPerformance.IMPROVED_PERFORMANCE_PROJECT;
    case 4:
    case "IMPROVED_PERFORMANCE_USER_GRANT":
      return ImprovedPerformance.IMPROVED_PERFORMANCE_USER_GRANT;
    case 5:
    case "IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED":
      return ImprovedPerformance.IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ImprovedPerformance.UNRECOGNIZED;
  }
}

export function improvedPerformanceToJSON(object: ImprovedPerformance): string {
  switch (object) {
    case ImprovedPerformance.IMPROVED_PERFORMANCE_UNSPECIFIED:
      return "IMPROVED_PERFORMANCE_UNSPECIFIED";
    case ImprovedPerformance.IMPROVED_PERFORMANCE_ORG_BY_ID:
      return "IMPROVED_PERFORMANCE_ORG_BY_ID";
    case ImprovedPerformance.IMPROVED_PERFORMANCE_PROJECT_GRANT:
      return "IMPROVED_PERFORMANCE_PROJECT_GRANT";
    case ImprovedPerformance.IMPROVED_PERFORMANCE_PROJECT:
      return "IMPROVED_PERFORMANCE_PROJECT";
    case ImprovedPerformance.IMPROVED_PERFORMANCE_USER_GRANT:
      return "IMPROVED_PERFORMANCE_USER_GRANT";
    case ImprovedPerformance.IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED:
      return "IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED";
    case ImprovedPerformance.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** FeatureFlag is a simple boolean Feature setting, without further payload. */
export interface FeatureFlag {
  enabled: boolean;
  source: Source;
}

export interface ImprovedPerformanceFeatureFlag {
  executionPaths: ImprovedPerformance[];
  source: Source;
}

function createBaseFeatureFlag(): FeatureFlag {
  return { enabled: false, source: 0 };
}

export const FeatureFlag = {
  encode(message: FeatureFlag, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.enabled === true) {
      writer.uint32(8).bool(message.enabled);
    }
    if (message.source !== 0) {
      writer.uint32(16).int32(message.source);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FeatureFlag {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFeatureFlag();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.enabled = reader.bool();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.source = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FeatureFlag {
    return {
      enabled: isSet(object.enabled) ? Boolean(object.enabled) : false,
      source: isSet(object.source) ? sourceFromJSON(object.source) : 0,
    };
  },

  toJSON(message: FeatureFlag): unknown {
    const obj: any = {};
    message.enabled !== undefined && (obj.enabled = message.enabled);
    message.source !== undefined && (obj.source = sourceToJSON(message.source));
    return obj;
  },

  create(base?: DeepPartial<FeatureFlag>): FeatureFlag {
    return FeatureFlag.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<FeatureFlag>): FeatureFlag {
    const message = createBaseFeatureFlag();
    message.enabled = object.enabled ?? false;
    message.source = object.source ?? 0;
    return message;
  },
};

function createBaseImprovedPerformanceFeatureFlag(): ImprovedPerformanceFeatureFlag {
  return { executionPaths: [], source: 0 };
}

export const ImprovedPerformanceFeatureFlag = {
  encode(message: ImprovedPerformanceFeatureFlag, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();
    for (const v of message.executionPaths) {
      writer.int32(v);
    }
    writer.ldelim();
    if (message.source !== 0) {
      writer.uint32(16).int32(message.source);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImprovedPerformanceFeatureFlag {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImprovedPerformanceFeatureFlag();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag == 8) {
            message.executionPaths.push(reader.int32() as any);
            continue;
          }

          if (tag == 10) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.executionPaths.push(reader.int32() as any);
            }

            continue;
          }

          break;
        case 2:
          if (tag != 16) {
            break;
          }

          message.source = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ImprovedPerformanceFeatureFlag {
    return {
      executionPaths: Array.isArray(object?.executionPaths)
        ? object.executionPaths.map((e: any) => improvedPerformanceFromJSON(e))
        : [],
      source: isSet(object.source) ? sourceFromJSON(object.source) : 0,
    };
  },

  toJSON(message: ImprovedPerformanceFeatureFlag): unknown {
    const obj: any = {};
    if (message.executionPaths) {
      obj.executionPaths = message.executionPaths.map((e) => improvedPerformanceToJSON(e));
    } else {
      obj.executionPaths = [];
    }
    message.source !== undefined && (obj.source = sourceToJSON(message.source));
    return obj;
  },

  create(base?: DeepPartial<ImprovedPerformanceFeatureFlag>): ImprovedPerformanceFeatureFlag {
    return ImprovedPerformanceFeatureFlag.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImprovedPerformanceFeatureFlag>): ImprovedPerformanceFeatureFlag {
    const message = createBaseImprovedPerformanceFeatureFlag();
    message.executionPaths = object.executionPaths?.map((e) => e) || [];
    message.source = object.source ?? 0;
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

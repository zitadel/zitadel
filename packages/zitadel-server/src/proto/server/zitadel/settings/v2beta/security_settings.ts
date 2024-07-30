/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.settings.v2beta";

export interface SecuritySettings {
  embeddedIframe: EmbeddedIframeSettings | undefined;
  enableImpersonation: boolean;
}

export interface EmbeddedIframeSettings {
  enabled: boolean;
  allowedOrigins: string[];
}

function createBaseSecuritySettings(): SecuritySettings {
  return { embeddedIframe: undefined, enableImpersonation: false };
}

export const SecuritySettings = {
  encode(message: SecuritySettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.embeddedIframe !== undefined) {
      EmbeddedIframeSettings.encode(message.embeddedIframe, writer.uint32(10).fork()).ldelim();
    }
    if (message.enableImpersonation === true) {
      writer.uint32(16).bool(message.enableImpersonation);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SecuritySettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSecuritySettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.embeddedIframe = EmbeddedIframeSettings.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.enableImpersonation = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SecuritySettings {
    return {
      embeddedIframe: isSet(object.embeddedIframe) ? EmbeddedIframeSettings.fromJSON(object.embeddedIframe) : undefined,
      enableImpersonation: isSet(object.enableImpersonation) ? Boolean(object.enableImpersonation) : false,
    };
  },

  toJSON(message: SecuritySettings): unknown {
    const obj: any = {};
    message.embeddedIframe !== undefined &&
      (obj.embeddedIframe = message.embeddedIframe ? EmbeddedIframeSettings.toJSON(message.embeddedIframe) : undefined);
    message.enableImpersonation !== undefined && (obj.enableImpersonation = message.enableImpersonation);
    return obj;
  },

  create(base?: DeepPartial<SecuritySettings>): SecuritySettings {
    return SecuritySettings.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SecuritySettings>): SecuritySettings {
    const message = createBaseSecuritySettings();
    message.embeddedIframe = (object.embeddedIframe !== undefined && object.embeddedIframe !== null)
      ? EmbeddedIframeSettings.fromPartial(object.embeddedIframe)
      : undefined;
    message.enableImpersonation = object.enableImpersonation ?? false;
    return message;
  },
};

function createBaseEmbeddedIframeSettings(): EmbeddedIframeSettings {
  return { enabled: false, allowedOrigins: [] };
}

export const EmbeddedIframeSettings = {
  encode(message: EmbeddedIframeSettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.enabled === true) {
      writer.uint32(8).bool(message.enabled);
    }
    for (const v of message.allowedOrigins) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EmbeddedIframeSettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEmbeddedIframeSettings();
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
          if (tag != 18) {
            break;
          }

          message.allowedOrigins.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EmbeddedIframeSettings {
    return {
      enabled: isSet(object.enabled) ? Boolean(object.enabled) : false,
      allowedOrigins: Array.isArray(object?.allowedOrigins) ? object.allowedOrigins.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: EmbeddedIframeSettings): unknown {
    const obj: any = {};
    message.enabled !== undefined && (obj.enabled = message.enabled);
    if (message.allowedOrigins) {
      obj.allowedOrigins = message.allowedOrigins.map((e) => e);
    } else {
      obj.allowedOrigins = [];
    }
    return obj;
  },

  create(base?: DeepPartial<EmbeddedIframeSettings>): EmbeddedIframeSettings {
    return EmbeddedIframeSettings.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EmbeddedIframeSettings>): EmbeddedIframeSettings {
    const message = createBaseEmbeddedIframeSettings();
    message.enabled = object.enabled ?? false;
    message.allowedOrigins = object.allowedOrigins?.map((e) => e) || [];
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

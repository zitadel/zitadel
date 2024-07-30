/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { ResourceOwnerType, resourceOwnerTypeFromJSON, resourceOwnerTypeToJSON } from "./settings";

export const protobufPackage = "zitadel.settings.v2beta";

export interface DomainSettings {
  loginNameIncludesDomain: boolean;
  requireOrgDomainVerification: boolean;
  smtpSenderAddressMatchesInstanceDomain: boolean;
  /** resource_owner_type returns if the setting is managed on the organization or on the instance */
  resourceOwnerType: ResourceOwnerType;
}

function createBaseDomainSettings(): DomainSettings {
  return {
    loginNameIncludesDomain: false,
    requireOrgDomainVerification: false,
    smtpSenderAddressMatchesInstanceDomain: false,
    resourceOwnerType: 0,
  };
}

export const DomainSettings = {
  encode(message: DomainSettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.loginNameIncludesDomain === true) {
      writer.uint32(8).bool(message.loginNameIncludesDomain);
    }
    if (message.requireOrgDomainVerification === true) {
      writer.uint32(16).bool(message.requireOrgDomainVerification);
    }
    if (message.smtpSenderAddressMatchesInstanceDomain === true) {
      writer.uint32(24).bool(message.smtpSenderAddressMatchesInstanceDomain);
    }
    if (message.resourceOwnerType !== 0) {
      writer.uint32(48).int32(message.resourceOwnerType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DomainSettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDomainSettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.loginNameIncludesDomain = reader.bool();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.requireOrgDomainVerification = reader.bool();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.smtpSenderAddressMatchesInstanceDomain = reader.bool();
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.resourceOwnerType = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DomainSettings {
    return {
      loginNameIncludesDomain: isSet(object.loginNameIncludesDomain) ? Boolean(object.loginNameIncludesDomain) : false,
      requireOrgDomainVerification: isSet(object.requireOrgDomainVerification)
        ? Boolean(object.requireOrgDomainVerification)
        : false,
      smtpSenderAddressMatchesInstanceDomain: isSet(object.smtpSenderAddressMatchesInstanceDomain)
        ? Boolean(object.smtpSenderAddressMatchesInstanceDomain)
        : false,
      resourceOwnerType: isSet(object.resourceOwnerType) ? resourceOwnerTypeFromJSON(object.resourceOwnerType) : 0,
    };
  },

  toJSON(message: DomainSettings): unknown {
    const obj: any = {};
    message.loginNameIncludesDomain !== undefined && (obj.loginNameIncludesDomain = message.loginNameIncludesDomain);
    message.requireOrgDomainVerification !== undefined &&
      (obj.requireOrgDomainVerification = message.requireOrgDomainVerification);
    message.smtpSenderAddressMatchesInstanceDomain !== undefined &&
      (obj.smtpSenderAddressMatchesInstanceDomain = message.smtpSenderAddressMatchesInstanceDomain);
    message.resourceOwnerType !== undefined &&
      (obj.resourceOwnerType = resourceOwnerTypeToJSON(message.resourceOwnerType));
    return obj;
  },

  create(base?: DeepPartial<DomainSettings>): DomainSettings {
    return DomainSettings.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DomainSettings>): DomainSettings {
    const message = createBaseDomainSettings();
    message.loginNameIncludesDomain = object.loginNameIncludesDomain ?? false;
    message.requireOrgDomainVerification = object.requireOrgDomainVerification ?? false;
    message.smtpSenderAddressMatchesInstanceDomain = object.smtpSenderAddressMatchesInstanceDomain ?? false;
    message.resourceOwnerType = object.resourceOwnerType ?? 0;
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

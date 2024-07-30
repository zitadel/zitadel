/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Duration } from "../../../google/protobuf/duration";
import { ResourceOwnerType, resourceOwnerTypeFromJSON, resourceOwnerTypeToJSON } from "./settings";

export const protobufPackage = "zitadel.settings.v2beta";

export enum SecondFactorType {
  SECOND_FACTOR_TYPE_UNSPECIFIED = 0,
  /** SECOND_FACTOR_TYPE_OTP - This is the type for TOTP */
  SECOND_FACTOR_TYPE_OTP = 1,
  SECOND_FACTOR_TYPE_U2F = 2,
  SECOND_FACTOR_TYPE_OTP_EMAIL = 3,
  SECOND_FACTOR_TYPE_OTP_SMS = 4,
  UNRECOGNIZED = -1,
}

export function secondFactorTypeFromJSON(object: any): SecondFactorType {
  switch (object) {
    case 0:
    case "SECOND_FACTOR_TYPE_UNSPECIFIED":
      return SecondFactorType.SECOND_FACTOR_TYPE_UNSPECIFIED;
    case 1:
    case "SECOND_FACTOR_TYPE_OTP":
      return SecondFactorType.SECOND_FACTOR_TYPE_OTP;
    case 2:
    case "SECOND_FACTOR_TYPE_U2F":
      return SecondFactorType.SECOND_FACTOR_TYPE_U2F;
    case 3:
    case "SECOND_FACTOR_TYPE_OTP_EMAIL":
      return SecondFactorType.SECOND_FACTOR_TYPE_OTP_EMAIL;
    case 4:
    case "SECOND_FACTOR_TYPE_OTP_SMS":
      return SecondFactorType.SECOND_FACTOR_TYPE_OTP_SMS;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SecondFactorType.UNRECOGNIZED;
  }
}

export function secondFactorTypeToJSON(object: SecondFactorType): string {
  switch (object) {
    case SecondFactorType.SECOND_FACTOR_TYPE_UNSPECIFIED:
      return "SECOND_FACTOR_TYPE_UNSPECIFIED";
    case SecondFactorType.SECOND_FACTOR_TYPE_OTP:
      return "SECOND_FACTOR_TYPE_OTP";
    case SecondFactorType.SECOND_FACTOR_TYPE_U2F:
      return "SECOND_FACTOR_TYPE_U2F";
    case SecondFactorType.SECOND_FACTOR_TYPE_OTP_EMAIL:
      return "SECOND_FACTOR_TYPE_OTP_EMAIL";
    case SecondFactorType.SECOND_FACTOR_TYPE_OTP_SMS:
      return "SECOND_FACTOR_TYPE_OTP_SMS";
    case SecondFactorType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum MultiFactorType {
  MULTI_FACTOR_TYPE_UNSPECIFIED = 0,
  MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION = 1,
  UNRECOGNIZED = -1,
}

export function multiFactorTypeFromJSON(object: any): MultiFactorType {
  switch (object) {
    case 0:
    case "MULTI_FACTOR_TYPE_UNSPECIFIED":
      return MultiFactorType.MULTI_FACTOR_TYPE_UNSPECIFIED;
    case 1:
    case "MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION":
      return MultiFactorType.MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION;
    case -1:
    case "UNRECOGNIZED":
    default:
      return MultiFactorType.UNRECOGNIZED;
  }
}

export function multiFactorTypeToJSON(object: MultiFactorType): string {
  switch (object) {
    case MultiFactorType.MULTI_FACTOR_TYPE_UNSPECIFIED:
      return "MULTI_FACTOR_TYPE_UNSPECIFIED";
    case MultiFactorType.MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION:
      return "MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION";
    case MultiFactorType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum PasskeysType {
  PASSKEYS_TYPE_NOT_ALLOWED = 0,
  PASSKEYS_TYPE_ALLOWED = 1,
  UNRECOGNIZED = -1,
}

export function passkeysTypeFromJSON(object: any): PasskeysType {
  switch (object) {
    case 0:
    case "PASSKEYS_TYPE_NOT_ALLOWED":
      return PasskeysType.PASSKEYS_TYPE_NOT_ALLOWED;
    case 1:
    case "PASSKEYS_TYPE_ALLOWED":
      return PasskeysType.PASSKEYS_TYPE_ALLOWED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return PasskeysType.UNRECOGNIZED;
  }
}

export function passkeysTypeToJSON(object: PasskeysType): string {
  switch (object) {
    case PasskeysType.PASSKEYS_TYPE_NOT_ALLOWED:
      return "PASSKEYS_TYPE_NOT_ALLOWED";
    case PasskeysType.PASSKEYS_TYPE_ALLOWED:
      return "PASSKEYS_TYPE_ALLOWED";
    case PasskeysType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum IdentityProviderType {
  IDENTITY_PROVIDER_TYPE_UNSPECIFIED = 0,
  IDENTITY_PROVIDER_TYPE_OIDC = 1,
  IDENTITY_PROVIDER_TYPE_JWT = 2,
  IDENTITY_PROVIDER_TYPE_LDAP = 3,
  IDENTITY_PROVIDER_TYPE_OAUTH = 4,
  IDENTITY_PROVIDER_TYPE_AZURE_AD = 5,
  IDENTITY_PROVIDER_TYPE_GITHUB = 6,
  IDENTITY_PROVIDER_TYPE_GITHUB_ES = 7,
  IDENTITY_PROVIDER_TYPE_GITLAB = 8,
  IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED = 9,
  IDENTITY_PROVIDER_TYPE_GOOGLE = 10,
  IDENTITY_PROVIDER_TYPE_SAML = 11,
  UNRECOGNIZED = -1,
}

export function identityProviderTypeFromJSON(object: any): IdentityProviderType {
  switch (object) {
    case 0:
    case "IDENTITY_PROVIDER_TYPE_UNSPECIFIED":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_UNSPECIFIED;
    case 1:
    case "IDENTITY_PROVIDER_TYPE_OIDC":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_OIDC;
    case 2:
    case "IDENTITY_PROVIDER_TYPE_JWT":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_JWT;
    case 3:
    case "IDENTITY_PROVIDER_TYPE_LDAP":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_LDAP;
    case 4:
    case "IDENTITY_PROVIDER_TYPE_OAUTH":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_OAUTH;
    case 5:
    case "IDENTITY_PROVIDER_TYPE_AZURE_AD":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_AZURE_AD;
    case 6:
    case "IDENTITY_PROVIDER_TYPE_GITHUB":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB;
    case 7:
    case "IDENTITY_PROVIDER_TYPE_GITHUB_ES":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB_ES;
    case 8:
    case "IDENTITY_PROVIDER_TYPE_GITLAB":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB;
    case 9:
    case "IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED;
    case 10:
    case "IDENTITY_PROVIDER_TYPE_GOOGLE":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_GOOGLE;
    case 11:
    case "IDENTITY_PROVIDER_TYPE_SAML":
      return IdentityProviderType.IDENTITY_PROVIDER_TYPE_SAML;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IdentityProviderType.UNRECOGNIZED;
  }
}

export function identityProviderTypeToJSON(object: IdentityProviderType): string {
  switch (object) {
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_UNSPECIFIED:
      return "IDENTITY_PROVIDER_TYPE_UNSPECIFIED";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_OIDC:
      return "IDENTITY_PROVIDER_TYPE_OIDC";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_JWT:
      return "IDENTITY_PROVIDER_TYPE_JWT";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_LDAP:
      return "IDENTITY_PROVIDER_TYPE_LDAP";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_OAUTH:
      return "IDENTITY_PROVIDER_TYPE_OAUTH";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_AZURE_AD:
      return "IDENTITY_PROVIDER_TYPE_AZURE_AD";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB:
      return "IDENTITY_PROVIDER_TYPE_GITHUB";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITHUB_ES:
      return "IDENTITY_PROVIDER_TYPE_GITHUB_ES";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB:
      return "IDENTITY_PROVIDER_TYPE_GITLAB";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED:
      return "IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_GOOGLE:
      return "IDENTITY_PROVIDER_TYPE_GOOGLE";
    case IdentityProviderType.IDENTITY_PROVIDER_TYPE_SAML:
      return "IDENTITY_PROVIDER_TYPE_SAML";
    case IdentityProviderType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface LoginSettings {
  allowUsernamePassword: boolean;
  allowRegister: boolean;
  allowExternalIdp: boolean;
  forceMfa: boolean;
  passkeysType: PasskeysType;
  hidePasswordReset: boolean;
  ignoreUnknownUsernames: boolean;
  defaultRedirectUri: string;
  passwordCheckLifetime: Duration | undefined;
  externalLoginCheckLifetime: Duration | undefined;
  mfaInitSkipLifetime: Duration | undefined;
  secondFactorCheckLifetime: Duration | undefined;
  multiFactorCheckLifetime: Duration | undefined;
  secondFactors: SecondFactorType[];
  multiFactors: MultiFactorType[];
  /** If set to true, the suffix (@domain.com) of an unknown username input on the login screen will be matched against the org domains and will redirect to the registration of that organization on success. */
  allowDomainDiscovery: boolean;
  disableLoginWithEmail: boolean;
  disableLoginWithPhone: boolean;
  /** resource_owner_type returns if the settings is managed on the organization or on the instance */
  resourceOwnerType: ResourceOwnerType;
  forceMfaLocalOnly: boolean;
}

export interface IdentityProvider {
  id: string;
  name: string;
  type: IdentityProviderType;
}

function createBaseLoginSettings(): LoginSettings {
  return {
    allowUsernamePassword: false,
    allowRegister: false,
    allowExternalIdp: false,
    forceMfa: false,
    passkeysType: 0,
    hidePasswordReset: false,
    ignoreUnknownUsernames: false,
    defaultRedirectUri: "",
    passwordCheckLifetime: undefined,
    externalLoginCheckLifetime: undefined,
    mfaInitSkipLifetime: undefined,
    secondFactorCheckLifetime: undefined,
    multiFactorCheckLifetime: undefined,
    secondFactors: [],
    multiFactors: [],
    allowDomainDiscovery: false,
    disableLoginWithEmail: false,
    disableLoginWithPhone: false,
    resourceOwnerType: 0,
    forceMfaLocalOnly: false,
  };
}

export const LoginSettings = {
  encode(message: LoginSettings, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.allowUsernamePassword === true) {
      writer.uint32(8).bool(message.allowUsernamePassword);
    }
    if (message.allowRegister === true) {
      writer.uint32(16).bool(message.allowRegister);
    }
    if (message.allowExternalIdp === true) {
      writer.uint32(24).bool(message.allowExternalIdp);
    }
    if (message.forceMfa === true) {
      writer.uint32(32).bool(message.forceMfa);
    }
    if (message.passkeysType !== 0) {
      writer.uint32(40).int32(message.passkeysType);
    }
    if (message.hidePasswordReset === true) {
      writer.uint32(48).bool(message.hidePasswordReset);
    }
    if (message.ignoreUnknownUsernames === true) {
      writer.uint32(56).bool(message.ignoreUnknownUsernames);
    }
    if (message.defaultRedirectUri !== "") {
      writer.uint32(66).string(message.defaultRedirectUri);
    }
    if (message.passwordCheckLifetime !== undefined) {
      Duration.encode(message.passwordCheckLifetime, writer.uint32(74).fork()).ldelim();
    }
    if (message.externalLoginCheckLifetime !== undefined) {
      Duration.encode(message.externalLoginCheckLifetime, writer.uint32(82).fork()).ldelim();
    }
    if (message.mfaInitSkipLifetime !== undefined) {
      Duration.encode(message.mfaInitSkipLifetime, writer.uint32(90).fork()).ldelim();
    }
    if (message.secondFactorCheckLifetime !== undefined) {
      Duration.encode(message.secondFactorCheckLifetime, writer.uint32(98).fork()).ldelim();
    }
    if (message.multiFactorCheckLifetime !== undefined) {
      Duration.encode(message.multiFactorCheckLifetime, writer.uint32(106).fork()).ldelim();
    }
    writer.uint32(114).fork();
    for (const v of message.secondFactors) {
      writer.int32(v);
    }
    writer.ldelim();
    writer.uint32(122).fork();
    for (const v of message.multiFactors) {
      writer.int32(v);
    }
    writer.ldelim();
    if (message.allowDomainDiscovery === true) {
      writer.uint32(128).bool(message.allowDomainDiscovery);
    }
    if (message.disableLoginWithEmail === true) {
      writer.uint32(136).bool(message.disableLoginWithEmail);
    }
    if (message.disableLoginWithPhone === true) {
      writer.uint32(144).bool(message.disableLoginWithPhone);
    }
    if (message.resourceOwnerType !== 0) {
      writer.uint32(152).int32(message.resourceOwnerType);
    }
    if (message.forceMfaLocalOnly === true) {
      writer.uint32(176).bool(message.forceMfaLocalOnly);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LoginSettings {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginSettings();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.allowUsernamePassword = reader.bool();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.allowRegister = reader.bool();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.allowExternalIdp = reader.bool();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.forceMfa = reader.bool();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.passkeysType = reader.int32() as any;
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.hidePasswordReset = reader.bool();
          continue;
        case 7:
          if (tag != 56) {
            break;
          }

          message.ignoreUnknownUsernames = reader.bool();
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.defaultRedirectUri = reader.string();
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.passwordCheckLifetime = Duration.decode(reader, reader.uint32());
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          message.externalLoginCheckLifetime = Duration.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag != 90) {
            break;
          }

          message.mfaInitSkipLifetime = Duration.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag != 98) {
            break;
          }

          message.secondFactorCheckLifetime = Duration.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag != 106) {
            break;
          }

          message.multiFactorCheckLifetime = Duration.decode(reader, reader.uint32());
          continue;
        case 14:
          if (tag == 112) {
            message.secondFactors.push(reader.int32() as any);
            continue;
          }

          if (tag == 114) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.secondFactors.push(reader.int32() as any);
            }

            continue;
          }

          break;
        case 15:
          if (tag == 120) {
            message.multiFactors.push(reader.int32() as any);
            continue;
          }

          if (tag == 122) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.multiFactors.push(reader.int32() as any);
            }

            continue;
          }

          break;
        case 16:
          if (tag != 128) {
            break;
          }

          message.allowDomainDiscovery = reader.bool();
          continue;
        case 17:
          if (tag != 136) {
            break;
          }

          message.disableLoginWithEmail = reader.bool();
          continue;
        case 18:
          if (tag != 144) {
            break;
          }

          message.disableLoginWithPhone = reader.bool();
          continue;
        case 19:
          if (tag != 152) {
            break;
          }

          message.resourceOwnerType = reader.int32() as any;
          continue;
        case 22:
          if (tag != 176) {
            break;
          }

          message.forceMfaLocalOnly = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LoginSettings {
    return {
      allowUsernamePassword: isSet(object.allowUsernamePassword) ? Boolean(object.allowUsernamePassword) : false,
      allowRegister: isSet(object.allowRegister) ? Boolean(object.allowRegister) : false,
      allowExternalIdp: isSet(object.allowExternalIdp) ? Boolean(object.allowExternalIdp) : false,
      forceMfa: isSet(object.forceMfa) ? Boolean(object.forceMfa) : false,
      passkeysType: isSet(object.passkeysType) ? passkeysTypeFromJSON(object.passkeysType) : 0,
      hidePasswordReset: isSet(object.hidePasswordReset) ? Boolean(object.hidePasswordReset) : false,
      ignoreUnknownUsernames: isSet(object.ignoreUnknownUsernames) ? Boolean(object.ignoreUnknownUsernames) : false,
      defaultRedirectUri: isSet(object.defaultRedirectUri) ? String(object.defaultRedirectUri) : "",
      passwordCheckLifetime: isSet(object.passwordCheckLifetime)
        ? Duration.fromJSON(object.passwordCheckLifetime)
        : undefined,
      externalLoginCheckLifetime: isSet(object.externalLoginCheckLifetime)
        ? Duration.fromJSON(object.externalLoginCheckLifetime)
        : undefined,
      mfaInitSkipLifetime: isSet(object.mfaInitSkipLifetime)
        ? Duration.fromJSON(object.mfaInitSkipLifetime)
        : undefined,
      secondFactorCheckLifetime: isSet(object.secondFactorCheckLifetime)
        ? Duration.fromJSON(object.secondFactorCheckLifetime)
        : undefined,
      multiFactorCheckLifetime: isSet(object.multiFactorCheckLifetime)
        ? Duration.fromJSON(object.multiFactorCheckLifetime)
        : undefined,
      secondFactors: Array.isArray(object?.secondFactors)
        ? object.secondFactors.map((e: any) => secondFactorTypeFromJSON(e))
        : [],
      multiFactors: Array.isArray(object?.multiFactors)
        ? object.multiFactors.map((e: any) => multiFactorTypeFromJSON(e))
        : [],
      allowDomainDiscovery: isSet(object.allowDomainDiscovery) ? Boolean(object.allowDomainDiscovery) : false,
      disableLoginWithEmail: isSet(object.disableLoginWithEmail) ? Boolean(object.disableLoginWithEmail) : false,
      disableLoginWithPhone: isSet(object.disableLoginWithPhone) ? Boolean(object.disableLoginWithPhone) : false,
      resourceOwnerType: isSet(object.resourceOwnerType) ? resourceOwnerTypeFromJSON(object.resourceOwnerType) : 0,
      forceMfaLocalOnly: isSet(object.forceMfaLocalOnly) ? Boolean(object.forceMfaLocalOnly) : false,
    };
  },

  toJSON(message: LoginSettings): unknown {
    const obj: any = {};
    message.allowUsernamePassword !== undefined && (obj.allowUsernamePassword = message.allowUsernamePassword);
    message.allowRegister !== undefined && (obj.allowRegister = message.allowRegister);
    message.allowExternalIdp !== undefined && (obj.allowExternalIdp = message.allowExternalIdp);
    message.forceMfa !== undefined && (obj.forceMfa = message.forceMfa);
    message.passkeysType !== undefined && (obj.passkeysType = passkeysTypeToJSON(message.passkeysType));
    message.hidePasswordReset !== undefined && (obj.hidePasswordReset = message.hidePasswordReset);
    message.ignoreUnknownUsernames !== undefined && (obj.ignoreUnknownUsernames = message.ignoreUnknownUsernames);
    message.defaultRedirectUri !== undefined && (obj.defaultRedirectUri = message.defaultRedirectUri);
    message.passwordCheckLifetime !== undefined && (obj.passwordCheckLifetime = message.passwordCheckLifetime
      ? Duration.toJSON(message.passwordCheckLifetime)
      : undefined);
    message.externalLoginCheckLifetime !== undefined &&
      (obj.externalLoginCheckLifetime = message.externalLoginCheckLifetime
        ? Duration.toJSON(message.externalLoginCheckLifetime)
        : undefined);
    message.mfaInitSkipLifetime !== undefined &&
      (obj.mfaInitSkipLifetime = message.mfaInitSkipLifetime
        ? Duration.toJSON(message.mfaInitSkipLifetime)
        : undefined);
    message.secondFactorCheckLifetime !== undefined &&
      (obj.secondFactorCheckLifetime = message.secondFactorCheckLifetime
        ? Duration.toJSON(message.secondFactorCheckLifetime)
        : undefined);
    message.multiFactorCheckLifetime !== undefined && (obj.multiFactorCheckLifetime = message.multiFactorCheckLifetime
      ? Duration.toJSON(message.multiFactorCheckLifetime)
      : undefined);
    if (message.secondFactors) {
      obj.secondFactors = message.secondFactors.map((e) => secondFactorTypeToJSON(e));
    } else {
      obj.secondFactors = [];
    }
    if (message.multiFactors) {
      obj.multiFactors = message.multiFactors.map((e) => multiFactorTypeToJSON(e));
    } else {
      obj.multiFactors = [];
    }
    message.allowDomainDiscovery !== undefined && (obj.allowDomainDiscovery = message.allowDomainDiscovery);
    message.disableLoginWithEmail !== undefined && (obj.disableLoginWithEmail = message.disableLoginWithEmail);
    message.disableLoginWithPhone !== undefined && (obj.disableLoginWithPhone = message.disableLoginWithPhone);
    message.resourceOwnerType !== undefined &&
      (obj.resourceOwnerType = resourceOwnerTypeToJSON(message.resourceOwnerType));
    message.forceMfaLocalOnly !== undefined && (obj.forceMfaLocalOnly = message.forceMfaLocalOnly);
    return obj;
  },

  create(base?: DeepPartial<LoginSettings>): LoginSettings {
    return LoginSettings.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LoginSettings>): LoginSettings {
    const message = createBaseLoginSettings();
    message.allowUsernamePassword = object.allowUsernamePassword ?? false;
    message.allowRegister = object.allowRegister ?? false;
    message.allowExternalIdp = object.allowExternalIdp ?? false;
    message.forceMfa = object.forceMfa ?? false;
    message.passkeysType = object.passkeysType ?? 0;
    message.hidePasswordReset = object.hidePasswordReset ?? false;
    message.ignoreUnknownUsernames = object.ignoreUnknownUsernames ?? false;
    message.defaultRedirectUri = object.defaultRedirectUri ?? "";
    message.passwordCheckLifetime =
      (object.passwordCheckLifetime !== undefined && object.passwordCheckLifetime !== null)
        ? Duration.fromPartial(object.passwordCheckLifetime)
        : undefined;
    message.externalLoginCheckLifetime =
      (object.externalLoginCheckLifetime !== undefined && object.externalLoginCheckLifetime !== null)
        ? Duration.fromPartial(object.externalLoginCheckLifetime)
        : undefined;
    message.mfaInitSkipLifetime = (object.mfaInitSkipLifetime !== undefined && object.mfaInitSkipLifetime !== null)
      ? Duration.fromPartial(object.mfaInitSkipLifetime)
      : undefined;
    message.secondFactorCheckLifetime =
      (object.secondFactorCheckLifetime !== undefined && object.secondFactorCheckLifetime !== null)
        ? Duration.fromPartial(object.secondFactorCheckLifetime)
        : undefined;
    message.multiFactorCheckLifetime =
      (object.multiFactorCheckLifetime !== undefined && object.multiFactorCheckLifetime !== null)
        ? Duration.fromPartial(object.multiFactorCheckLifetime)
        : undefined;
    message.secondFactors = object.secondFactors?.map((e) => e) || [];
    message.multiFactors = object.multiFactors?.map((e) => e) || [];
    message.allowDomainDiscovery = object.allowDomainDiscovery ?? false;
    message.disableLoginWithEmail = object.disableLoginWithEmail ?? false;
    message.disableLoginWithPhone = object.disableLoginWithPhone ?? false;
    message.resourceOwnerType = object.resourceOwnerType ?? 0;
    message.forceMfaLocalOnly = object.forceMfaLocalOnly ?? false;
    return message;
  },
};

function createBaseIdentityProvider(): IdentityProvider {
  return { id: "", name: "", type: 0 };
}

export const IdentityProvider = {
  encode(message: IdentityProvider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.type !== 0) {
      writer.uint32(24).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IdentityProvider {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentityProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IdentityProvider {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      type: isSet(object.type) ? identityProviderTypeFromJSON(object.type) : 0,
    };
  },

  toJSON(message: IdentityProvider): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.type !== undefined && (obj.type = identityProviderTypeToJSON(message.type));
    return obj;
  },

  create(base?: DeepPartial<IdentityProvider>): IdentityProvider {
    return IdentityProvider.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IdentityProvider>): IdentityProvider {
    const message = createBaseIdentityProvider();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.type = object.type ?? 0;
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

/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Duration } from "../google/protobuf/duration";
import { IDPLoginPolicyLink } from "./idp";
import { ObjectDetails } from "./object";

export const protobufPackage = "zitadel.policy.v1";

export enum SecondFactorType {
  SECOND_FACTOR_TYPE_UNSPECIFIED = 0,
  SECOND_FACTOR_TYPE_OTP = 1,
  SECOND_FACTOR_TYPE_U2F = 2,
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

export enum PasswordlessType {
  PASSWORDLESS_TYPE_NOT_ALLOWED = 0,
  /** PASSWORDLESS_TYPE_ALLOWED - PLANNED: PASSWORDLESS_TYPE_WITH_CERT */
  PASSWORDLESS_TYPE_ALLOWED = 1,
  UNRECOGNIZED = -1,
}

export function passwordlessTypeFromJSON(object: any): PasswordlessType {
  switch (object) {
    case 0:
    case "PASSWORDLESS_TYPE_NOT_ALLOWED":
      return PasswordlessType.PASSWORDLESS_TYPE_NOT_ALLOWED;
    case 1:
    case "PASSWORDLESS_TYPE_ALLOWED":
      return PasswordlessType.PASSWORDLESS_TYPE_ALLOWED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return PasswordlessType.UNRECOGNIZED;
  }
}

export function passwordlessTypeToJSON(object: PasswordlessType): string {
  switch (object) {
    case PasswordlessType.PASSWORDLESS_TYPE_NOT_ALLOWED:
      return "PASSWORDLESS_TYPE_NOT_ALLOWED";
    case PasswordlessType.PASSWORDLESS_TYPE_ALLOWED:
      return "PASSWORDLESS_TYPE_ALLOWED";
    case PasswordlessType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** deprecated: please use DomainPolicy instead */
export interface OrgIAMPolicy {
  details: ObjectDetails | undefined;
  userLoginMustBeDomain: boolean;
  isDefault: boolean;
}

export interface DomainPolicy {
  details: ObjectDetails | undefined;
  userLoginMustBeDomain: boolean;
  isDefault: boolean;
  validateOrgDomains: boolean;
  smtpSenderAddressMatchesInstanceDomain: boolean;
}

export interface LabelPolicy {
  details:
    | ObjectDetails
    | undefined;
  /** hex value for primary color */
  primaryColor: string;
  /** defines if the organization's admin changed the policy */
  isDefault: boolean;
  /** hides the org suffix on the login form if the scope \"urn:zitadel:iam:org:domain:primary:{domainname}\" is set */
  hideLoginNameSuffix: boolean;
  /** hex value for secondary color */
  warnColor: string;
  /** hex value for background color */
  backgroundColor: string;
  /** hex value for font color */
  fontColor: string;
  /** hex value for primary color dark theme */
  primaryColorDark: string;
  /** hex value for background color dark theme */
  backgroundColorDark: string;
  /** hex value for warning color dark theme */
  warnColorDark: string;
  /** hex value for font color dark theme */
  fontColorDark: string;
  disableWatermark: boolean;
  logoUrl: string;
  iconUrl: string;
  logoUrlDark: string;
  iconUrlDark: string;
  fontUrl: string;
}

export interface LoginPolicy {
  details: ObjectDetails | undefined;
  allowUsernamePassword: boolean;
  allowRegister: boolean;
  allowExternalIdp: boolean;
  forceMfa: boolean;
  passwordlessType: PasswordlessType;
  isDefault: boolean;
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
  idps: IDPLoginPolicyLink[];
  /** If set to true, the suffix (@domain.com) of an unknown username input on the login screen will be matched against the org domains and will redirect to the registration of that organization on success. */
  allowDomainDiscovery: boolean;
  disableLoginWithEmail: boolean;
  disableLoginWithPhone: boolean;
}

export interface PasswordComplexityPolicy {
  details: ObjectDetails | undefined;
  minLength: number;
  hasUppercase: boolean;
  hasLowercase: boolean;
  hasNumber: boolean;
  hasSymbol: boolean;
  isDefault: boolean;
}

export interface PasswordAgePolicy {
  details: ObjectDetails | undefined;
  maxAgeDays: number;
  expireWarnDays: number;
  isDefault: boolean;
}

export interface LockoutPolicy {
  details: ObjectDetails | undefined;
  maxPasswordAttempts: number;
  isDefault: boolean;
}

export interface PrivacyPolicy {
  details: ObjectDetails | undefined;
  tosLink: string;
  privacyLink: string;
  isDefault: boolean;
  helpLink: string;
}

export interface NotificationPolicy {
  details: ObjectDetails | undefined;
  isDefault: boolean;
  passwordChange: boolean;
}

function createBaseOrgIAMPolicy(): OrgIAMPolicy {
  return { details: undefined, userLoginMustBeDomain: false, isDefault: false };
}

export const OrgIAMPolicy = {
  encode(message: OrgIAMPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(16).bool(message.userLoginMustBeDomain);
    }
    if (message.isDefault === true) {
      writer.uint32(24).bool(message.isDefault);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrgIAMPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrgIAMPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.userLoginMustBeDomain = reader.bool();
          break;
        case 3:
          message.isDefault = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): OrgIAMPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
    };
  },

  toJSON(message: OrgIAMPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    return obj;
  },

  create(base?: DeepPartial<OrgIAMPolicy>): OrgIAMPolicy {
    return OrgIAMPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OrgIAMPolicy>): OrgIAMPolicy {
    const message = createBaseOrgIAMPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    message.isDefault = object.isDefault ?? false;
    return message;
  },
};

function createBaseDomainPolicy(): DomainPolicy {
  return {
    details: undefined,
    userLoginMustBeDomain: false,
    isDefault: false,
    validateOrgDomains: false,
    smtpSenderAddressMatchesInstanceDomain: false,
  };
}

export const DomainPolicy = {
  encode(message: DomainPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(16).bool(message.userLoginMustBeDomain);
    }
    if (message.isDefault === true) {
      writer.uint32(24).bool(message.isDefault);
    }
    if (message.validateOrgDomains === true) {
      writer.uint32(32).bool(message.validateOrgDomains);
    }
    if (message.smtpSenderAddressMatchesInstanceDomain === true) {
      writer.uint32(40).bool(message.smtpSenderAddressMatchesInstanceDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DomainPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDomainPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.userLoginMustBeDomain = reader.bool();
          break;
        case 3:
          message.isDefault = reader.bool();
          break;
        case 4:
          message.validateOrgDomains = reader.bool();
          break;
        case 5:
          message.smtpSenderAddressMatchesInstanceDomain = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DomainPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
      validateOrgDomains: isSet(object.validateOrgDomains) ? Boolean(object.validateOrgDomains) : false,
      smtpSenderAddressMatchesInstanceDomain: isSet(object.smtpSenderAddressMatchesInstanceDomain)
        ? Boolean(object.smtpSenderAddressMatchesInstanceDomain)
        : false,
    };
  },

  toJSON(message: DomainPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    message.validateOrgDomains !== undefined && (obj.validateOrgDomains = message.validateOrgDomains);
    message.smtpSenderAddressMatchesInstanceDomain !== undefined &&
      (obj.smtpSenderAddressMatchesInstanceDomain = message.smtpSenderAddressMatchesInstanceDomain);
    return obj;
  },

  create(base?: DeepPartial<DomainPolicy>): DomainPolicy {
    return DomainPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DomainPolicy>): DomainPolicy {
    const message = createBaseDomainPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    message.isDefault = object.isDefault ?? false;
    message.validateOrgDomains = object.validateOrgDomains ?? false;
    message.smtpSenderAddressMatchesInstanceDomain = object.smtpSenderAddressMatchesInstanceDomain ?? false;
    return message;
  },
};

function createBaseLabelPolicy(): LabelPolicy {
  return {
    details: undefined,
    primaryColor: "",
    isDefault: false,
    hideLoginNameSuffix: false,
    warnColor: "",
    backgroundColor: "",
    fontColor: "",
    primaryColorDark: "",
    backgroundColorDark: "",
    warnColorDark: "",
    fontColorDark: "",
    disableWatermark: false,
    logoUrl: "",
    iconUrl: "",
    logoUrlDark: "",
    iconUrlDark: "",
    fontUrl: "",
  };
}

export const LabelPolicy = {
  encode(message: LabelPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.primaryColor !== "") {
      writer.uint32(18).string(message.primaryColor);
    }
    if (message.isDefault === true) {
      writer.uint32(32).bool(message.isDefault);
    }
    if (message.hideLoginNameSuffix === true) {
      writer.uint32(40).bool(message.hideLoginNameSuffix);
    }
    if (message.warnColor !== "") {
      writer.uint32(50).string(message.warnColor);
    }
    if (message.backgroundColor !== "") {
      writer.uint32(58).string(message.backgroundColor);
    }
    if (message.fontColor !== "") {
      writer.uint32(66).string(message.fontColor);
    }
    if (message.primaryColorDark !== "") {
      writer.uint32(74).string(message.primaryColorDark);
    }
    if (message.backgroundColorDark !== "") {
      writer.uint32(82).string(message.backgroundColorDark);
    }
    if (message.warnColorDark !== "") {
      writer.uint32(90).string(message.warnColorDark);
    }
    if (message.fontColorDark !== "") {
      writer.uint32(98).string(message.fontColorDark);
    }
    if (message.disableWatermark === true) {
      writer.uint32(104).bool(message.disableWatermark);
    }
    if (message.logoUrl !== "") {
      writer.uint32(114).string(message.logoUrl);
    }
    if (message.iconUrl !== "") {
      writer.uint32(122).string(message.iconUrl);
    }
    if (message.logoUrlDark !== "") {
      writer.uint32(130).string(message.logoUrlDark);
    }
    if (message.iconUrlDark !== "") {
      writer.uint32(138).string(message.iconUrlDark);
    }
    if (message.fontUrl !== "") {
      writer.uint32(146).string(message.fontUrl);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LabelPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLabelPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.primaryColor = reader.string();
          break;
        case 4:
          message.isDefault = reader.bool();
          break;
        case 5:
          message.hideLoginNameSuffix = reader.bool();
          break;
        case 6:
          message.warnColor = reader.string();
          break;
        case 7:
          message.backgroundColor = reader.string();
          break;
        case 8:
          message.fontColor = reader.string();
          break;
        case 9:
          message.primaryColorDark = reader.string();
          break;
        case 10:
          message.backgroundColorDark = reader.string();
          break;
        case 11:
          message.warnColorDark = reader.string();
          break;
        case 12:
          message.fontColorDark = reader.string();
          break;
        case 13:
          message.disableWatermark = reader.bool();
          break;
        case 14:
          message.logoUrl = reader.string();
          break;
        case 15:
          message.iconUrl = reader.string();
          break;
        case 16:
          message.logoUrlDark = reader.string();
          break;
        case 17:
          message.iconUrlDark = reader.string();
          break;
        case 18:
          message.fontUrl = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LabelPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      primaryColor: isSet(object.primaryColor) ? String(object.primaryColor) : "",
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
      hideLoginNameSuffix: isSet(object.hideLoginNameSuffix) ? Boolean(object.hideLoginNameSuffix) : false,
      warnColor: isSet(object.warnColor) ? String(object.warnColor) : "",
      backgroundColor: isSet(object.backgroundColor) ? String(object.backgroundColor) : "",
      fontColor: isSet(object.fontColor) ? String(object.fontColor) : "",
      primaryColorDark: isSet(object.primaryColorDark) ? String(object.primaryColorDark) : "",
      backgroundColorDark: isSet(object.backgroundColorDark) ? String(object.backgroundColorDark) : "",
      warnColorDark: isSet(object.warnColorDark) ? String(object.warnColorDark) : "",
      fontColorDark: isSet(object.fontColorDark) ? String(object.fontColorDark) : "",
      disableWatermark: isSet(object.disableWatermark) ? Boolean(object.disableWatermark) : false,
      logoUrl: isSet(object.logoUrl) ? String(object.logoUrl) : "",
      iconUrl: isSet(object.iconUrl) ? String(object.iconUrl) : "",
      logoUrlDark: isSet(object.logoUrlDark) ? String(object.logoUrlDark) : "",
      iconUrlDark: isSet(object.iconUrlDark) ? String(object.iconUrlDark) : "",
      fontUrl: isSet(object.fontUrl) ? String(object.fontUrl) : "",
    };
  },

  toJSON(message: LabelPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.primaryColor !== undefined && (obj.primaryColor = message.primaryColor);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    message.hideLoginNameSuffix !== undefined && (obj.hideLoginNameSuffix = message.hideLoginNameSuffix);
    message.warnColor !== undefined && (obj.warnColor = message.warnColor);
    message.backgroundColor !== undefined && (obj.backgroundColor = message.backgroundColor);
    message.fontColor !== undefined && (obj.fontColor = message.fontColor);
    message.primaryColorDark !== undefined && (obj.primaryColorDark = message.primaryColorDark);
    message.backgroundColorDark !== undefined && (obj.backgroundColorDark = message.backgroundColorDark);
    message.warnColorDark !== undefined && (obj.warnColorDark = message.warnColorDark);
    message.fontColorDark !== undefined && (obj.fontColorDark = message.fontColorDark);
    message.disableWatermark !== undefined && (obj.disableWatermark = message.disableWatermark);
    message.logoUrl !== undefined && (obj.logoUrl = message.logoUrl);
    message.iconUrl !== undefined && (obj.iconUrl = message.iconUrl);
    message.logoUrlDark !== undefined && (obj.logoUrlDark = message.logoUrlDark);
    message.iconUrlDark !== undefined && (obj.iconUrlDark = message.iconUrlDark);
    message.fontUrl !== undefined && (obj.fontUrl = message.fontUrl);
    return obj;
  },

  create(base?: DeepPartial<LabelPolicy>): LabelPolicy {
    return LabelPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LabelPolicy>): LabelPolicy {
    const message = createBaseLabelPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.primaryColor = object.primaryColor ?? "";
    message.isDefault = object.isDefault ?? false;
    message.hideLoginNameSuffix = object.hideLoginNameSuffix ?? false;
    message.warnColor = object.warnColor ?? "";
    message.backgroundColor = object.backgroundColor ?? "";
    message.fontColor = object.fontColor ?? "";
    message.primaryColorDark = object.primaryColorDark ?? "";
    message.backgroundColorDark = object.backgroundColorDark ?? "";
    message.warnColorDark = object.warnColorDark ?? "";
    message.fontColorDark = object.fontColorDark ?? "";
    message.disableWatermark = object.disableWatermark ?? false;
    message.logoUrl = object.logoUrl ?? "";
    message.iconUrl = object.iconUrl ?? "";
    message.logoUrlDark = object.logoUrlDark ?? "";
    message.iconUrlDark = object.iconUrlDark ?? "";
    message.fontUrl = object.fontUrl ?? "";
    return message;
  },
};

function createBaseLoginPolicy(): LoginPolicy {
  return {
    details: undefined,
    allowUsernamePassword: false,
    allowRegister: false,
    allowExternalIdp: false,
    forceMfa: false,
    passwordlessType: 0,
    isDefault: false,
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
    idps: [],
    allowDomainDiscovery: false,
    disableLoginWithEmail: false,
    disableLoginWithPhone: false,
  };
}

export const LoginPolicy = {
  encode(message: LoginPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.allowUsernamePassword === true) {
      writer.uint32(16).bool(message.allowUsernamePassword);
    }
    if (message.allowRegister === true) {
      writer.uint32(24).bool(message.allowRegister);
    }
    if (message.allowExternalIdp === true) {
      writer.uint32(32).bool(message.allowExternalIdp);
    }
    if (message.forceMfa === true) {
      writer.uint32(40).bool(message.forceMfa);
    }
    if (message.passwordlessType !== 0) {
      writer.uint32(48).int32(message.passwordlessType);
    }
    if (message.isDefault === true) {
      writer.uint32(56).bool(message.isDefault);
    }
    if (message.hidePasswordReset === true) {
      writer.uint32(64).bool(message.hidePasswordReset);
    }
    if (message.ignoreUnknownUsernames === true) {
      writer.uint32(72).bool(message.ignoreUnknownUsernames);
    }
    if (message.defaultRedirectUri !== "") {
      writer.uint32(82).string(message.defaultRedirectUri);
    }
    if (message.passwordCheckLifetime !== undefined) {
      Duration.encode(message.passwordCheckLifetime, writer.uint32(90).fork()).ldelim();
    }
    if (message.externalLoginCheckLifetime !== undefined) {
      Duration.encode(message.externalLoginCheckLifetime, writer.uint32(98).fork()).ldelim();
    }
    if (message.mfaInitSkipLifetime !== undefined) {
      Duration.encode(message.mfaInitSkipLifetime, writer.uint32(106).fork()).ldelim();
    }
    if (message.secondFactorCheckLifetime !== undefined) {
      Duration.encode(message.secondFactorCheckLifetime, writer.uint32(114).fork()).ldelim();
    }
    if (message.multiFactorCheckLifetime !== undefined) {
      Duration.encode(message.multiFactorCheckLifetime, writer.uint32(122).fork()).ldelim();
    }
    writer.uint32(130).fork();
    for (const v of message.secondFactors) {
      writer.int32(v);
    }
    writer.ldelim();
    writer.uint32(138).fork();
    for (const v of message.multiFactors) {
      writer.int32(v);
    }
    writer.ldelim();
    for (const v of message.idps) {
      IDPLoginPolicyLink.encode(v!, writer.uint32(146).fork()).ldelim();
    }
    if (message.allowDomainDiscovery === true) {
      writer.uint32(152).bool(message.allowDomainDiscovery);
    }
    if (message.disableLoginWithEmail === true) {
      writer.uint32(160).bool(message.disableLoginWithEmail);
    }
    if (message.disableLoginWithPhone === true) {
      writer.uint32(168).bool(message.disableLoginWithPhone);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LoginPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.allowUsernamePassword = reader.bool();
          break;
        case 3:
          message.allowRegister = reader.bool();
          break;
        case 4:
          message.allowExternalIdp = reader.bool();
          break;
        case 5:
          message.forceMfa = reader.bool();
          break;
        case 6:
          message.passwordlessType = reader.int32() as any;
          break;
        case 7:
          message.isDefault = reader.bool();
          break;
        case 8:
          message.hidePasswordReset = reader.bool();
          break;
        case 9:
          message.ignoreUnknownUsernames = reader.bool();
          break;
        case 10:
          message.defaultRedirectUri = reader.string();
          break;
        case 11:
          message.passwordCheckLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 12:
          message.externalLoginCheckLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 13:
          message.mfaInitSkipLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 14:
          message.secondFactorCheckLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 15:
          message.multiFactorCheckLifetime = Duration.decode(reader, reader.uint32());
          break;
        case 16:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.secondFactors.push(reader.int32() as any);
            }
          } else {
            message.secondFactors.push(reader.int32() as any);
          }
          break;
        case 17:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.multiFactors.push(reader.int32() as any);
            }
          } else {
            message.multiFactors.push(reader.int32() as any);
          }
          break;
        case 18:
          message.idps.push(IDPLoginPolicyLink.decode(reader, reader.uint32()));
          break;
        case 19:
          message.allowDomainDiscovery = reader.bool();
          break;
        case 20:
          message.disableLoginWithEmail = reader.bool();
          break;
        case 21:
          message.disableLoginWithPhone = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LoginPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      allowUsernamePassword: isSet(object.allowUsernamePassword) ? Boolean(object.allowUsernamePassword) : false,
      allowRegister: isSet(object.allowRegister) ? Boolean(object.allowRegister) : false,
      allowExternalIdp: isSet(object.allowExternalIdp) ? Boolean(object.allowExternalIdp) : false,
      forceMfa: isSet(object.forceMfa) ? Boolean(object.forceMfa) : false,
      passwordlessType: isSet(object.passwordlessType) ? passwordlessTypeFromJSON(object.passwordlessType) : 0,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
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
      idps: Array.isArray(object?.idps) ? object.idps.map((e: any) => IDPLoginPolicyLink.fromJSON(e)) : [],
      allowDomainDiscovery: isSet(object.allowDomainDiscovery) ? Boolean(object.allowDomainDiscovery) : false,
      disableLoginWithEmail: isSet(object.disableLoginWithEmail) ? Boolean(object.disableLoginWithEmail) : false,
      disableLoginWithPhone: isSet(object.disableLoginWithPhone) ? Boolean(object.disableLoginWithPhone) : false,
    };
  },

  toJSON(message: LoginPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.allowUsernamePassword !== undefined && (obj.allowUsernamePassword = message.allowUsernamePassword);
    message.allowRegister !== undefined && (obj.allowRegister = message.allowRegister);
    message.allowExternalIdp !== undefined && (obj.allowExternalIdp = message.allowExternalIdp);
    message.forceMfa !== undefined && (obj.forceMfa = message.forceMfa);
    message.passwordlessType !== undefined && (obj.passwordlessType = passwordlessTypeToJSON(message.passwordlessType));
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
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
    if (message.idps) {
      obj.idps = message.idps.map((e) => e ? IDPLoginPolicyLink.toJSON(e) : undefined);
    } else {
      obj.idps = [];
    }
    message.allowDomainDiscovery !== undefined && (obj.allowDomainDiscovery = message.allowDomainDiscovery);
    message.disableLoginWithEmail !== undefined && (obj.disableLoginWithEmail = message.disableLoginWithEmail);
    message.disableLoginWithPhone !== undefined && (obj.disableLoginWithPhone = message.disableLoginWithPhone);
    return obj;
  },

  create(base?: DeepPartial<LoginPolicy>): LoginPolicy {
    return LoginPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LoginPolicy>): LoginPolicy {
    const message = createBaseLoginPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.allowUsernamePassword = object.allowUsernamePassword ?? false;
    message.allowRegister = object.allowRegister ?? false;
    message.allowExternalIdp = object.allowExternalIdp ?? false;
    message.forceMfa = object.forceMfa ?? false;
    message.passwordlessType = object.passwordlessType ?? 0;
    message.isDefault = object.isDefault ?? false;
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
    message.idps = object.idps?.map((e) => IDPLoginPolicyLink.fromPartial(e)) || [];
    message.allowDomainDiscovery = object.allowDomainDiscovery ?? false;
    message.disableLoginWithEmail = object.disableLoginWithEmail ?? false;
    message.disableLoginWithPhone = object.disableLoginWithPhone ?? false;
    return message;
  },
};

function createBasePasswordComplexityPolicy(): PasswordComplexityPolicy {
  return {
    details: undefined,
    minLength: 0,
    hasUppercase: false,
    hasLowercase: false,
    hasNumber: false,
    hasSymbol: false,
    isDefault: false,
  };
}

export const PasswordComplexityPolicy = {
  encode(message: PasswordComplexityPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.minLength !== 0) {
      writer.uint32(16).uint64(message.minLength);
    }
    if (message.hasUppercase === true) {
      writer.uint32(24).bool(message.hasUppercase);
    }
    if (message.hasLowercase === true) {
      writer.uint32(32).bool(message.hasLowercase);
    }
    if (message.hasNumber === true) {
      writer.uint32(40).bool(message.hasNumber);
    }
    if (message.hasSymbol === true) {
      writer.uint32(48).bool(message.hasSymbol);
    }
    if (message.isDefault === true) {
      writer.uint32(56).bool(message.isDefault);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordComplexityPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordComplexityPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.minLength = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.hasUppercase = reader.bool();
          break;
        case 4:
          message.hasLowercase = reader.bool();
          break;
        case 5:
          message.hasNumber = reader.bool();
          break;
        case 6:
          message.hasSymbol = reader.bool();
          break;
        case 7:
          message.isDefault = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordComplexityPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      minLength: isSet(object.minLength) ? Number(object.minLength) : 0,
      hasUppercase: isSet(object.hasUppercase) ? Boolean(object.hasUppercase) : false,
      hasLowercase: isSet(object.hasLowercase) ? Boolean(object.hasLowercase) : false,
      hasNumber: isSet(object.hasNumber) ? Boolean(object.hasNumber) : false,
      hasSymbol: isSet(object.hasSymbol) ? Boolean(object.hasSymbol) : false,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
    };
  },

  toJSON(message: PasswordComplexityPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.minLength !== undefined && (obj.minLength = Math.round(message.minLength));
    message.hasUppercase !== undefined && (obj.hasUppercase = message.hasUppercase);
    message.hasLowercase !== undefined && (obj.hasLowercase = message.hasLowercase);
    message.hasNumber !== undefined && (obj.hasNumber = message.hasNumber);
    message.hasSymbol !== undefined && (obj.hasSymbol = message.hasSymbol);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    return obj;
  },

  create(base?: DeepPartial<PasswordComplexityPolicy>): PasswordComplexityPolicy {
    return PasswordComplexityPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordComplexityPolicy>): PasswordComplexityPolicy {
    const message = createBasePasswordComplexityPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.minLength = object.minLength ?? 0;
    message.hasUppercase = object.hasUppercase ?? false;
    message.hasLowercase = object.hasLowercase ?? false;
    message.hasNumber = object.hasNumber ?? false;
    message.hasSymbol = object.hasSymbol ?? false;
    message.isDefault = object.isDefault ?? false;
    return message;
  },
};

function createBasePasswordAgePolicy(): PasswordAgePolicy {
  return { details: undefined, maxAgeDays: 0, expireWarnDays: 0, isDefault: false };
}

export const PasswordAgePolicy = {
  encode(message: PasswordAgePolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.maxAgeDays !== 0) {
      writer.uint32(16).uint64(message.maxAgeDays);
    }
    if (message.expireWarnDays !== 0) {
      writer.uint32(24).uint64(message.expireWarnDays);
    }
    if (message.isDefault === true) {
      writer.uint32(32).bool(message.isDefault);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordAgePolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordAgePolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.maxAgeDays = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.expireWarnDays = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.isDefault = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordAgePolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      maxAgeDays: isSet(object.maxAgeDays) ? Number(object.maxAgeDays) : 0,
      expireWarnDays: isSet(object.expireWarnDays) ? Number(object.expireWarnDays) : 0,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
    };
  },

  toJSON(message: PasswordAgePolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.maxAgeDays !== undefined && (obj.maxAgeDays = Math.round(message.maxAgeDays));
    message.expireWarnDays !== undefined && (obj.expireWarnDays = Math.round(message.expireWarnDays));
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    return obj;
  },

  create(base?: DeepPartial<PasswordAgePolicy>): PasswordAgePolicy {
    return PasswordAgePolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordAgePolicy>): PasswordAgePolicy {
    const message = createBasePasswordAgePolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.maxAgeDays = object.maxAgeDays ?? 0;
    message.expireWarnDays = object.expireWarnDays ?? 0;
    message.isDefault = object.isDefault ?? false;
    return message;
  },
};

function createBaseLockoutPolicy(): LockoutPolicy {
  return { details: undefined, maxPasswordAttempts: 0, isDefault: false };
}

export const LockoutPolicy = {
  encode(message: LockoutPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.maxPasswordAttempts !== 0) {
      writer.uint32(16).uint64(message.maxPasswordAttempts);
    }
    if (message.isDefault === true) {
      writer.uint32(32).bool(message.isDefault);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LockoutPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLockoutPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.maxPasswordAttempts = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.isDefault = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LockoutPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      maxPasswordAttempts: isSet(object.maxPasswordAttempts) ? Number(object.maxPasswordAttempts) : 0,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
    };
  },

  toJSON(message: LockoutPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.maxPasswordAttempts !== undefined && (obj.maxPasswordAttempts = Math.round(message.maxPasswordAttempts));
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    return obj;
  },

  create(base?: DeepPartial<LockoutPolicy>): LockoutPolicy {
    return LockoutPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LockoutPolicy>): LockoutPolicy {
    const message = createBaseLockoutPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.maxPasswordAttempts = object.maxPasswordAttempts ?? 0;
    message.isDefault = object.isDefault ?? false;
    return message;
  },
};

function createBasePrivacyPolicy(): PrivacyPolicy {
  return { details: undefined, tosLink: "", privacyLink: "", isDefault: false, helpLink: "" };
}

export const PrivacyPolicy = {
  encode(message: PrivacyPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.tosLink !== "") {
      writer.uint32(18).string(message.tosLink);
    }
    if (message.privacyLink !== "") {
      writer.uint32(26).string(message.privacyLink);
    }
    if (message.isDefault === true) {
      writer.uint32(32).bool(message.isDefault);
    }
    if (message.helpLink !== "") {
      writer.uint32(42).string(message.helpLink);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PrivacyPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePrivacyPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.tosLink = reader.string();
          break;
        case 3:
          message.privacyLink = reader.string();
          break;
        case 4:
          message.isDefault = reader.bool();
          break;
        case 5:
          message.helpLink = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PrivacyPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      tosLink: isSet(object.tosLink) ? String(object.tosLink) : "",
      privacyLink: isSet(object.privacyLink) ? String(object.privacyLink) : "",
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
      helpLink: isSet(object.helpLink) ? String(object.helpLink) : "",
    };
  },

  toJSON(message: PrivacyPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.tosLink !== undefined && (obj.tosLink = message.tosLink);
    message.privacyLink !== undefined && (obj.privacyLink = message.privacyLink);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    message.helpLink !== undefined && (obj.helpLink = message.helpLink);
    return obj;
  },

  create(base?: DeepPartial<PrivacyPolicy>): PrivacyPolicy {
    return PrivacyPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PrivacyPolicy>): PrivacyPolicy {
    const message = createBasePrivacyPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.tosLink = object.tosLink ?? "";
    message.privacyLink = object.privacyLink ?? "";
    message.isDefault = object.isDefault ?? false;
    message.helpLink = object.helpLink ?? "";
    return message;
  },
};

function createBaseNotificationPolicy(): NotificationPolicy {
  return { details: undefined, isDefault: false, passwordChange: false };
}

export const NotificationPolicy = {
  encode(message: NotificationPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.isDefault === true) {
      writer.uint32(16).bool(message.isDefault);
    }
    if (message.passwordChange === true) {
      writer.uint32(24).bool(message.passwordChange);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NotificationPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNotificationPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.isDefault = reader.bool();
          break;
        case 3:
          message.passwordChange = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): NotificationPolicy {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
      passwordChange: isSet(object.passwordChange) ? Boolean(object.passwordChange) : false,
    };
  },

  toJSON(message: NotificationPolicy): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    message.passwordChange !== undefined && (obj.passwordChange = message.passwordChange);
    return obj;
  },

  create(base?: DeepPartial<NotificationPolicy>): NotificationPolicy {
    return NotificationPolicy.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<NotificationPolicy>): NotificationPolicy {
    const message = createBaseNotificationPolicy();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.isDefault = object.isDefault ?? false;
    message.passwordChange = object.passwordChange ?? false;
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

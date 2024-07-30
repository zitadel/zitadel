import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as zitadel_idp_pb from '../zitadel/idp_pb'; // proto import: "zitadel/idp.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"


export class OrgIAMPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): OrgIAMPolicy;
  hasDetails(): boolean;
  clearDetails(): OrgIAMPolicy;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): OrgIAMPolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): OrgIAMPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrgIAMPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: OrgIAMPolicy): OrgIAMPolicy.AsObject;
  static serializeBinaryToWriter(message: OrgIAMPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrgIAMPolicy;
  static deserializeBinaryFromReader(message: OrgIAMPolicy, reader: jspb.BinaryReader): OrgIAMPolicy;
}

export namespace OrgIAMPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    userLoginMustBeDomain: boolean,
    isDefault: boolean,
  }
}

export class DomainPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DomainPolicy;
  hasDetails(): boolean;
  clearDetails(): DomainPolicy;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): DomainPolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): DomainPolicy;

  getValidateOrgDomains(): boolean;
  setValidateOrgDomains(value: boolean): DomainPolicy;

  getSmtpSenderAddressMatchesInstanceDomain(): boolean;
  setSmtpSenderAddressMatchesInstanceDomain(value: boolean): DomainPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DomainPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: DomainPolicy): DomainPolicy.AsObject;
  static serializeBinaryToWriter(message: DomainPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DomainPolicy;
  static deserializeBinaryFromReader(message: DomainPolicy, reader: jspb.BinaryReader): DomainPolicy;
}

export namespace DomainPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    userLoginMustBeDomain: boolean,
    isDefault: boolean,
    validateOrgDomains: boolean,
    smtpSenderAddressMatchesInstanceDomain: boolean,
  }
}

export class LabelPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): LabelPolicy;
  hasDetails(): boolean;
  clearDetails(): LabelPolicy;

  getPrimaryColor(): string;
  setPrimaryColor(value: string): LabelPolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): LabelPolicy;

  getHideLoginNameSuffix(): boolean;
  setHideLoginNameSuffix(value: boolean): LabelPolicy;

  getWarnColor(): string;
  setWarnColor(value: string): LabelPolicy;

  getBackgroundColor(): string;
  setBackgroundColor(value: string): LabelPolicy;

  getFontColor(): string;
  setFontColor(value: string): LabelPolicy;

  getPrimaryColorDark(): string;
  setPrimaryColorDark(value: string): LabelPolicy;

  getBackgroundColorDark(): string;
  setBackgroundColorDark(value: string): LabelPolicy;

  getWarnColorDark(): string;
  setWarnColorDark(value: string): LabelPolicy;

  getFontColorDark(): string;
  setFontColorDark(value: string): LabelPolicy;

  getDisableWatermark(): boolean;
  setDisableWatermark(value: boolean): LabelPolicy;

  getLogoUrl(): string;
  setLogoUrl(value: string): LabelPolicy;

  getIconUrl(): string;
  setIconUrl(value: string): LabelPolicy;

  getLogoUrlDark(): string;
  setLogoUrlDark(value: string): LabelPolicy;

  getIconUrlDark(): string;
  setIconUrlDark(value: string): LabelPolicy;

  getFontUrl(): string;
  setFontUrl(value: string): LabelPolicy;

  getThemeMode(): ThemeMode;
  setThemeMode(value: ThemeMode): LabelPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LabelPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: LabelPolicy): LabelPolicy.AsObject;
  static serializeBinaryToWriter(message: LabelPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LabelPolicy;
  static deserializeBinaryFromReader(message: LabelPolicy, reader: jspb.BinaryReader): LabelPolicy;
}

export namespace LabelPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    primaryColor: string,
    isDefault: boolean,
    hideLoginNameSuffix: boolean,
    warnColor: string,
    backgroundColor: string,
    fontColor: string,
    primaryColorDark: string,
    backgroundColorDark: string,
    warnColorDark: string,
    fontColorDark: string,
    disableWatermark: boolean,
    logoUrl: string,
    iconUrl: string,
    logoUrlDark: string,
    iconUrlDark: string,
    fontUrl: string,
    themeMode: ThemeMode,
  }
}

export class LoginPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): LoginPolicy;
  hasDetails(): boolean;
  clearDetails(): LoginPolicy;

  getAllowUsernamePassword(): boolean;
  setAllowUsernamePassword(value: boolean): LoginPolicy;

  getAllowRegister(): boolean;
  setAllowRegister(value: boolean): LoginPolicy;

  getAllowExternalIdp(): boolean;
  setAllowExternalIdp(value: boolean): LoginPolicy;

  getForceMfa(): boolean;
  setForceMfa(value: boolean): LoginPolicy;

  getPasswordlessType(): PasswordlessType;
  setPasswordlessType(value: PasswordlessType): LoginPolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): LoginPolicy;

  getHidePasswordReset(): boolean;
  setHidePasswordReset(value: boolean): LoginPolicy;

  getIgnoreUnknownUsernames(): boolean;
  setIgnoreUnknownUsernames(value: boolean): LoginPolicy;

  getDefaultRedirectUri(): string;
  setDefaultRedirectUri(value: string): LoginPolicy;

  getPasswordCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setPasswordCheckLifetime(value?: google_protobuf_duration_pb.Duration): LoginPolicy;
  hasPasswordCheckLifetime(): boolean;
  clearPasswordCheckLifetime(): LoginPolicy;

  getExternalLoginCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setExternalLoginCheckLifetime(value?: google_protobuf_duration_pb.Duration): LoginPolicy;
  hasExternalLoginCheckLifetime(): boolean;
  clearExternalLoginCheckLifetime(): LoginPolicy;

  getMfaInitSkipLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMfaInitSkipLifetime(value?: google_protobuf_duration_pb.Duration): LoginPolicy;
  hasMfaInitSkipLifetime(): boolean;
  clearMfaInitSkipLifetime(): LoginPolicy;

  getSecondFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setSecondFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): LoginPolicy;
  hasSecondFactorCheckLifetime(): boolean;
  clearSecondFactorCheckLifetime(): LoginPolicy;

  getMultiFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMultiFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): LoginPolicy;
  hasMultiFactorCheckLifetime(): boolean;
  clearMultiFactorCheckLifetime(): LoginPolicy;

  getSecondFactorsList(): Array<SecondFactorType>;
  setSecondFactorsList(value: Array<SecondFactorType>): LoginPolicy;
  clearSecondFactorsList(): LoginPolicy;
  addSecondFactors(value: SecondFactorType, index?: number): LoginPolicy;

  getMultiFactorsList(): Array<MultiFactorType>;
  setMultiFactorsList(value: Array<MultiFactorType>): LoginPolicy;
  clearMultiFactorsList(): LoginPolicy;
  addMultiFactors(value: MultiFactorType, index?: number): LoginPolicy;

  getIdpsList(): Array<zitadel_idp_pb.IDPLoginPolicyLink>;
  setIdpsList(value: Array<zitadel_idp_pb.IDPLoginPolicyLink>): LoginPolicy;
  clearIdpsList(): LoginPolicy;
  addIdps(value?: zitadel_idp_pb.IDPLoginPolicyLink, index?: number): zitadel_idp_pb.IDPLoginPolicyLink;

  getAllowDomainDiscovery(): boolean;
  setAllowDomainDiscovery(value: boolean): LoginPolicy;

  getDisableLoginWithEmail(): boolean;
  setDisableLoginWithEmail(value: boolean): LoginPolicy;

  getDisableLoginWithPhone(): boolean;
  setDisableLoginWithPhone(value: boolean): LoginPolicy;

  getForceMfaLocalOnly(): boolean;
  setForceMfaLocalOnly(value: boolean): LoginPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: LoginPolicy): LoginPolicy.AsObject;
  static serializeBinaryToWriter(message: LoginPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginPolicy;
  static deserializeBinaryFromReader(message: LoginPolicy, reader: jspb.BinaryReader): LoginPolicy;
}

export namespace LoginPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    allowUsernamePassword: boolean,
    allowRegister: boolean,
    allowExternalIdp: boolean,
    forceMfa: boolean,
    passwordlessType: PasswordlessType,
    isDefault: boolean,
    hidePasswordReset: boolean,
    ignoreUnknownUsernames: boolean,
    defaultRedirectUri: string,
    passwordCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    externalLoginCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    mfaInitSkipLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    secondFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    multiFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    secondFactorsList: Array<SecondFactorType>,
    multiFactorsList: Array<MultiFactorType>,
    idpsList: Array<zitadel_idp_pb.IDPLoginPolicyLink.AsObject>,
    allowDomainDiscovery: boolean,
    disableLoginWithEmail: boolean,
    disableLoginWithPhone: boolean,
    forceMfaLocalOnly: boolean,
  }
}

export class PasswordComplexityPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): PasswordComplexityPolicy;
  hasDetails(): boolean;
  clearDetails(): PasswordComplexityPolicy;

  getMinLength(): number;
  setMinLength(value: number): PasswordComplexityPolicy;

  getHasUppercase(): boolean;
  setHasUppercase(value: boolean): PasswordComplexityPolicy;

  getHasLowercase(): boolean;
  setHasLowercase(value: boolean): PasswordComplexityPolicy;

  getHasNumber(): boolean;
  setHasNumber(value: boolean): PasswordComplexityPolicy;

  getHasSymbol(): boolean;
  setHasSymbol(value: boolean): PasswordComplexityPolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): PasswordComplexityPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordComplexityPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordComplexityPolicy): PasswordComplexityPolicy.AsObject;
  static serializeBinaryToWriter(message: PasswordComplexityPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordComplexityPolicy;
  static deserializeBinaryFromReader(message: PasswordComplexityPolicy, reader: jspb.BinaryReader): PasswordComplexityPolicy;
}

export namespace PasswordComplexityPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    minLength: number,
    hasUppercase: boolean,
    hasLowercase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
    isDefault: boolean,
  }
}

export class PasswordAgePolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): PasswordAgePolicy;
  hasDetails(): boolean;
  clearDetails(): PasswordAgePolicy;

  getMaxAgeDays(): number;
  setMaxAgeDays(value: number): PasswordAgePolicy;

  getExpireWarnDays(): number;
  setExpireWarnDays(value: number): PasswordAgePolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): PasswordAgePolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordAgePolicy.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordAgePolicy): PasswordAgePolicy.AsObject;
  static serializeBinaryToWriter(message: PasswordAgePolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordAgePolicy;
  static deserializeBinaryFromReader(message: PasswordAgePolicy, reader: jspb.BinaryReader): PasswordAgePolicy;
}

export namespace PasswordAgePolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    maxAgeDays: number,
    expireWarnDays: number,
    isDefault: boolean,
  }
}

export class LockoutPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): LockoutPolicy;
  hasDetails(): boolean;
  clearDetails(): LockoutPolicy;

  getMaxPasswordAttempts(): number;
  setMaxPasswordAttempts(value: number): LockoutPolicy;

  getMaxOtpAttempts(): number;
  setMaxOtpAttempts(value: number): LockoutPolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): LockoutPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LockoutPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: LockoutPolicy): LockoutPolicy.AsObject;
  static serializeBinaryToWriter(message: LockoutPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LockoutPolicy;
  static deserializeBinaryFromReader(message: LockoutPolicy, reader: jspb.BinaryReader): LockoutPolicy;
}

export namespace LockoutPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    maxPasswordAttempts: number,
    maxOtpAttempts: number,
    isDefault: boolean,
  }
}

export class PrivacyPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): PrivacyPolicy;
  hasDetails(): boolean;
  clearDetails(): PrivacyPolicy;

  getTosLink(): string;
  setTosLink(value: string): PrivacyPolicy;

  getPrivacyLink(): string;
  setPrivacyLink(value: string): PrivacyPolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): PrivacyPolicy;

  getHelpLink(): string;
  setHelpLink(value: string): PrivacyPolicy;

  getSupportEmail(): string;
  setSupportEmail(value: string): PrivacyPolicy;

  getDocsLink(): string;
  setDocsLink(value: string): PrivacyPolicy;

  getCustomLink(): string;
  setCustomLink(value: string): PrivacyPolicy;

  getCustomLinkText(): string;
  setCustomLinkText(value: string): PrivacyPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PrivacyPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: PrivacyPolicy): PrivacyPolicy.AsObject;
  static serializeBinaryToWriter(message: PrivacyPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PrivacyPolicy;
  static deserializeBinaryFromReader(message: PrivacyPolicy, reader: jspb.BinaryReader): PrivacyPolicy;
}

export namespace PrivacyPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    tosLink: string,
    privacyLink: string,
    isDefault: boolean,
    helpLink: string,
    supportEmail: string,
    docsLink: string,
    customLink: string,
    customLinkText: string,
  }
}

export class NotificationPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): NotificationPolicy;
  hasDetails(): boolean;
  clearDetails(): NotificationPolicy;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): NotificationPolicy;

  getPasswordChange(): boolean;
  setPasswordChange(value: boolean): NotificationPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationPolicy): NotificationPolicy.AsObject;
  static serializeBinaryToWriter(message: NotificationPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationPolicy;
  static deserializeBinaryFromReader(message: NotificationPolicy, reader: jspb.BinaryReader): NotificationPolicy;
}

export namespace NotificationPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    isDefault: boolean,
    passwordChange: boolean,
  }
}

export enum ThemeMode { 
  THEME_MODE_UNSPECIFIED = 0,
  THEME_MODE_AUTO = 1,
  THEME_MODE_DARK = 2,
  THEME_MODE_LIGHT = 3,
}
export enum SecondFactorType { 
  SECOND_FACTOR_TYPE_UNSPECIFIED = 0,
  SECOND_FACTOR_TYPE_OTP = 1,
  SECOND_FACTOR_TYPE_U2F = 2,
  SECOND_FACTOR_TYPE_OTP_EMAIL = 3,
  SECOND_FACTOR_TYPE_OTP_SMS = 4,
}
export enum MultiFactorType { 
  MULTI_FACTOR_TYPE_UNSPECIFIED = 0,
  MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION = 1,
}
export enum PasswordlessType { 
  PASSWORDLESS_TYPE_NOT_ALLOWED = 0,
  PASSWORDLESS_TYPE_ALLOWED = 1,
}

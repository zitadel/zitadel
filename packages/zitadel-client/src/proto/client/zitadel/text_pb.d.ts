import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"


export class MessageCustomText extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): MessageCustomText;
  hasDetails(): boolean;
  clearDetails(): MessageCustomText;

  getTitle(): string;
  setTitle(value: string): MessageCustomText;

  getPreHeader(): string;
  setPreHeader(value: string): MessageCustomText;

  getSubject(): string;
  setSubject(value: string): MessageCustomText;

  getGreeting(): string;
  setGreeting(value: string): MessageCustomText;

  getText(): string;
  setText(value: string): MessageCustomText;

  getButtonText(): string;
  setButtonText(value: string): MessageCustomText;

  getFooterText(): string;
  setFooterText(value: string): MessageCustomText;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): MessageCustomText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MessageCustomText.AsObject;
  static toObject(includeInstance: boolean, msg: MessageCustomText): MessageCustomText.AsObject;
  static serializeBinaryToWriter(message: MessageCustomText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MessageCustomText;
  static deserializeBinaryFromReader(message: MessageCustomText, reader: jspb.BinaryReader): MessageCustomText;
}

export namespace MessageCustomText {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
    isDefault: boolean,
  }
}

export class LoginCustomText extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): LoginCustomText;
  hasDetails(): boolean;
  clearDetails(): LoginCustomText;

  getSelectAccountText(): SelectAccountScreenText | undefined;
  setSelectAccountText(value?: SelectAccountScreenText): LoginCustomText;
  hasSelectAccountText(): boolean;
  clearSelectAccountText(): LoginCustomText;

  getLoginText(): LoginScreenText | undefined;
  setLoginText(value?: LoginScreenText): LoginCustomText;
  hasLoginText(): boolean;
  clearLoginText(): LoginCustomText;

  getPasswordText(): PasswordScreenText | undefined;
  setPasswordText(value?: PasswordScreenText): LoginCustomText;
  hasPasswordText(): boolean;
  clearPasswordText(): LoginCustomText;

  getUsernameChangeText(): UsernameChangeScreenText | undefined;
  setUsernameChangeText(value?: UsernameChangeScreenText): LoginCustomText;
  hasUsernameChangeText(): boolean;
  clearUsernameChangeText(): LoginCustomText;

  getUsernameChangeDoneText(): UsernameChangeDoneScreenText | undefined;
  setUsernameChangeDoneText(value?: UsernameChangeDoneScreenText): LoginCustomText;
  hasUsernameChangeDoneText(): boolean;
  clearUsernameChangeDoneText(): LoginCustomText;

  getInitPasswordText(): InitPasswordScreenText | undefined;
  setInitPasswordText(value?: InitPasswordScreenText): LoginCustomText;
  hasInitPasswordText(): boolean;
  clearInitPasswordText(): LoginCustomText;

  getInitPasswordDoneText(): InitPasswordDoneScreenText | undefined;
  setInitPasswordDoneText(value?: InitPasswordDoneScreenText): LoginCustomText;
  hasInitPasswordDoneText(): boolean;
  clearInitPasswordDoneText(): LoginCustomText;

  getEmailVerificationText(): EmailVerificationScreenText | undefined;
  setEmailVerificationText(value?: EmailVerificationScreenText): LoginCustomText;
  hasEmailVerificationText(): boolean;
  clearEmailVerificationText(): LoginCustomText;

  getEmailVerificationDoneText(): EmailVerificationDoneScreenText | undefined;
  setEmailVerificationDoneText(value?: EmailVerificationDoneScreenText): LoginCustomText;
  hasEmailVerificationDoneText(): boolean;
  clearEmailVerificationDoneText(): LoginCustomText;

  getInitializeUserText(): InitializeUserScreenText | undefined;
  setInitializeUserText(value?: InitializeUserScreenText): LoginCustomText;
  hasInitializeUserText(): boolean;
  clearInitializeUserText(): LoginCustomText;

  getInitializeDoneText(): InitializeUserDoneScreenText | undefined;
  setInitializeDoneText(value?: InitializeUserDoneScreenText): LoginCustomText;
  hasInitializeDoneText(): boolean;
  clearInitializeDoneText(): LoginCustomText;

  getInitMfaPromptText(): InitMFAPromptScreenText | undefined;
  setInitMfaPromptText(value?: InitMFAPromptScreenText): LoginCustomText;
  hasInitMfaPromptText(): boolean;
  clearInitMfaPromptText(): LoginCustomText;

  getInitMfaOtpText(): InitMFAOTPScreenText | undefined;
  setInitMfaOtpText(value?: InitMFAOTPScreenText): LoginCustomText;
  hasInitMfaOtpText(): boolean;
  clearInitMfaOtpText(): LoginCustomText;

  getInitMfaU2fText(): InitMFAU2FScreenText | undefined;
  setInitMfaU2fText(value?: InitMFAU2FScreenText): LoginCustomText;
  hasInitMfaU2fText(): boolean;
  clearInitMfaU2fText(): LoginCustomText;

  getInitMfaDoneText(): InitMFADoneScreenText | undefined;
  setInitMfaDoneText(value?: InitMFADoneScreenText): LoginCustomText;
  hasInitMfaDoneText(): boolean;
  clearInitMfaDoneText(): LoginCustomText;

  getMfaProvidersText(): MFAProvidersText | undefined;
  setMfaProvidersText(value?: MFAProvidersText): LoginCustomText;
  hasMfaProvidersText(): boolean;
  clearMfaProvidersText(): LoginCustomText;

  getVerifyMfaOtpText(): VerifyMFAOTPScreenText | undefined;
  setVerifyMfaOtpText(value?: VerifyMFAOTPScreenText): LoginCustomText;
  hasVerifyMfaOtpText(): boolean;
  clearVerifyMfaOtpText(): LoginCustomText;

  getVerifyMfaU2fText(): VerifyMFAU2FScreenText | undefined;
  setVerifyMfaU2fText(value?: VerifyMFAU2FScreenText): LoginCustomText;
  hasVerifyMfaU2fText(): boolean;
  clearVerifyMfaU2fText(): LoginCustomText;

  getPasswordlessText(): PasswordlessScreenText | undefined;
  setPasswordlessText(value?: PasswordlessScreenText): LoginCustomText;
  hasPasswordlessText(): boolean;
  clearPasswordlessText(): LoginCustomText;

  getPasswordChangeText(): PasswordChangeScreenText | undefined;
  setPasswordChangeText(value?: PasswordChangeScreenText): LoginCustomText;
  hasPasswordChangeText(): boolean;
  clearPasswordChangeText(): LoginCustomText;

  getPasswordChangeDoneText(): PasswordChangeDoneScreenText | undefined;
  setPasswordChangeDoneText(value?: PasswordChangeDoneScreenText): LoginCustomText;
  hasPasswordChangeDoneText(): boolean;
  clearPasswordChangeDoneText(): LoginCustomText;

  getPasswordResetDoneText(): PasswordResetDoneScreenText | undefined;
  setPasswordResetDoneText(value?: PasswordResetDoneScreenText): LoginCustomText;
  hasPasswordResetDoneText(): boolean;
  clearPasswordResetDoneText(): LoginCustomText;

  getRegistrationOptionText(): RegistrationOptionScreenText | undefined;
  setRegistrationOptionText(value?: RegistrationOptionScreenText): LoginCustomText;
  hasRegistrationOptionText(): boolean;
  clearRegistrationOptionText(): LoginCustomText;

  getRegistrationUserText(): RegistrationUserScreenText | undefined;
  setRegistrationUserText(value?: RegistrationUserScreenText): LoginCustomText;
  hasRegistrationUserText(): boolean;
  clearRegistrationUserText(): LoginCustomText;

  getRegistrationOrgText(): RegistrationOrgScreenText | undefined;
  setRegistrationOrgText(value?: RegistrationOrgScreenText): LoginCustomText;
  hasRegistrationOrgText(): boolean;
  clearRegistrationOrgText(): LoginCustomText;

  getLinkingUserDoneText(): LinkingUserDoneScreenText | undefined;
  setLinkingUserDoneText(value?: LinkingUserDoneScreenText): LoginCustomText;
  hasLinkingUserDoneText(): boolean;
  clearLinkingUserDoneText(): LoginCustomText;

  getExternalUserNotFoundText(): ExternalUserNotFoundScreenText | undefined;
  setExternalUserNotFoundText(value?: ExternalUserNotFoundScreenText): LoginCustomText;
  hasExternalUserNotFoundText(): boolean;
  clearExternalUserNotFoundText(): LoginCustomText;

  getSuccessLoginText(): SuccessLoginScreenText | undefined;
  setSuccessLoginText(value?: SuccessLoginScreenText): LoginCustomText;
  hasSuccessLoginText(): boolean;
  clearSuccessLoginText(): LoginCustomText;

  getLogoutText(): LogoutDoneScreenText | undefined;
  setLogoutText(value?: LogoutDoneScreenText): LoginCustomText;
  hasLogoutText(): boolean;
  clearLogoutText(): LoginCustomText;

  getFooterText(): FooterText | undefined;
  setFooterText(value?: FooterText): LoginCustomText;
  hasFooterText(): boolean;
  clearFooterText(): LoginCustomText;

  getPasswordlessPromptText(): PasswordlessPromptScreenText | undefined;
  setPasswordlessPromptText(value?: PasswordlessPromptScreenText): LoginCustomText;
  hasPasswordlessPromptText(): boolean;
  clearPasswordlessPromptText(): LoginCustomText;

  getPasswordlessRegistrationText(): PasswordlessRegistrationScreenText | undefined;
  setPasswordlessRegistrationText(value?: PasswordlessRegistrationScreenText): LoginCustomText;
  hasPasswordlessRegistrationText(): boolean;
  clearPasswordlessRegistrationText(): LoginCustomText;

  getPasswordlessRegistrationDoneText(): PasswordlessRegistrationDoneScreenText | undefined;
  setPasswordlessRegistrationDoneText(value?: PasswordlessRegistrationDoneScreenText): LoginCustomText;
  hasPasswordlessRegistrationDoneText(): boolean;
  clearPasswordlessRegistrationDoneText(): LoginCustomText;

  getExternalRegistrationUserOverviewText(): ExternalRegistrationUserOverviewScreenText | undefined;
  setExternalRegistrationUserOverviewText(value?: ExternalRegistrationUserOverviewScreenText): LoginCustomText;
  hasExternalRegistrationUserOverviewText(): boolean;
  clearExternalRegistrationUserOverviewText(): LoginCustomText;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): LoginCustomText;

  getLinkingUserPromptText(): LinkingUserPromptScreenText | undefined;
  setLinkingUserPromptText(value?: LinkingUserPromptScreenText): LoginCustomText;
  hasLinkingUserPromptText(): boolean;
  clearLinkingUserPromptText(): LoginCustomText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginCustomText.AsObject;
  static toObject(includeInstance: boolean, msg: LoginCustomText): LoginCustomText.AsObject;
  static serializeBinaryToWriter(message: LoginCustomText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginCustomText;
  static deserializeBinaryFromReader(message: LoginCustomText, reader: jspb.BinaryReader): LoginCustomText;
}

export namespace LoginCustomText {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    selectAccountText?: SelectAccountScreenText.AsObject,
    loginText?: LoginScreenText.AsObject,
    passwordText?: PasswordScreenText.AsObject,
    usernameChangeText?: UsernameChangeScreenText.AsObject,
    usernameChangeDoneText?: UsernameChangeDoneScreenText.AsObject,
    initPasswordText?: InitPasswordScreenText.AsObject,
    initPasswordDoneText?: InitPasswordDoneScreenText.AsObject,
    emailVerificationText?: EmailVerificationScreenText.AsObject,
    emailVerificationDoneText?: EmailVerificationDoneScreenText.AsObject,
    initializeUserText?: InitializeUserScreenText.AsObject,
    initializeDoneText?: InitializeUserDoneScreenText.AsObject,
    initMfaPromptText?: InitMFAPromptScreenText.AsObject,
    initMfaOtpText?: InitMFAOTPScreenText.AsObject,
    initMfaU2fText?: InitMFAU2FScreenText.AsObject,
    initMfaDoneText?: InitMFADoneScreenText.AsObject,
    mfaProvidersText?: MFAProvidersText.AsObject,
    verifyMfaOtpText?: VerifyMFAOTPScreenText.AsObject,
    verifyMfaU2fText?: VerifyMFAU2FScreenText.AsObject,
    passwordlessText?: PasswordlessScreenText.AsObject,
    passwordChangeText?: PasswordChangeScreenText.AsObject,
    passwordChangeDoneText?: PasswordChangeDoneScreenText.AsObject,
    passwordResetDoneText?: PasswordResetDoneScreenText.AsObject,
    registrationOptionText?: RegistrationOptionScreenText.AsObject,
    registrationUserText?: RegistrationUserScreenText.AsObject,
    registrationOrgText?: RegistrationOrgScreenText.AsObject,
    linkingUserDoneText?: LinkingUserDoneScreenText.AsObject,
    externalUserNotFoundText?: ExternalUserNotFoundScreenText.AsObject,
    successLoginText?: SuccessLoginScreenText.AsObject,
    logoutText?: LogoutDoneScreenText.AsObject,
    footerText?: FooterText.AsObject,
    passwordlessPromptText?: PasswordlessPromptScreenText.AsObject,
    passwordlessRegistrationText?: PasswordlessRegistrationScreenText.AsObject,
    passwordlessRegistrationDoneText?: PasswordlessRegistrationDoneScreenText.AsObject,
    externalRegistrationUserOverviewText?: ExternalRegistrationUserOverviewScreenText.AsObject,
    isDefault: boolean,
    linkingUserPromptText?: LinkingUserPromptScreenText.AsObject,
  }
}

export class SelectAccountScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): SelectAccountScreenText;

  getDescription(): string;
  setDescription(value: string): SelectAccountScreenText;

  getTitleLinkingProcess(): string;
  setTitleLinkingProcess(value: string): SelectAccountScreenText;

  getDescriptionLinkingProcess(): string;
  setDescriptionLinkingProcess(value: string): SelectAccountScreenText;

  getOtherUser(): string;
  setOtherUser(value: string): SelectAccountScreenText;

  getSessionStateActive(): string;
  setSessionStateActive(value: string): SelectAccountScreenText;

  getSessionStateInactive(): string;
  setSessionStateInactive(value: string): SelectAccountScreenText;

  getUserMustBeMemberOfOrg(): string;
  setUserMustBeMemberOfOrg(value: string): SelectAccountScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SelectAccountScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: SelectAccountScreenText): SelectAccountScreenText.AsObject;
  static serializeBinaryToWriter(message: SelectAccountScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SelectAccountScreenText;
  static deserializeBinaryFromReader(message: SelectAccountScreenText, reader: jspb.BinaryReader): SelectAccountScreenText;
}

export namespace SelectAccountScreenText {
  export type AsObject = {
    title: string,
    description: string,
    titleLinkingProcess: string,
    descriptionLinkingProcess: string,
    otherUser: string,
    sessionStateActive: string,
    sessionStateInactive: string,
    userMustBeMemberOfOrg: string,
  }
}

export class LoginScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): LoginScreenText;

  getDescription(): string;
  setDescription(value: string): LoginScreenText;

  getTitleLinkingProcess(): string;
  setTitleLinkingProcess(value: string): LoginScreenText;

  getDescriptionLinkingProcess(): string;
  setDescriptionLinkingProcess(value: string): LoginScreenText;

  getUserMustBeMemberOfOrg(): string;
  setUserMustBeMemberOfOrg(value: string): LoginScreenText;

  getLoginNameLabel(): string;
  setLoginNameLabel(value: string): LoginScreenText;

  getRegisterButtonText(): string;
  setRegisterButtonText(value: string): LoginScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): LoginScreenText;

  getExternalUserDescription(): string;
  setExternalUserDescription(value: string): LoginScreenText;

  getUserNamePlaceholder(): string;
  setUserNamePlaceholder(value: string): LoginScreenText;

  getLoginNamePlaceholder(): string;
  setLoginNamePlaceholder(value: string): LoginScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: LoginScreenText): LoginScreenText.AsObject;
  static serializeBinaryToWriter(message: LoginScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginScreenText;
  static deserializeBinaryFromReader(message: LoginScreenText, reader: jspb.BinaryReader): LoginScreenText;
}

export namespace LoginScreenText {
  export type AsObject = {
    title: string,
    description: string,
    titleLinkingProcess: string,
    descriptionLinkingProcess: string,
    userMustBeMemberOfOrg: string,
    loginNameLabel: string,
    registerButtonText: string,
    nextButtonText: string,
    externalUserDescription: string,
    userNamePlaceholder: string,
    loginNamePlaceholder: string,
  }
}

export class PasswordScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): PasswordScreenText;

  getDescription(): string;
  setDescription(value: string): PasswordScreenText;

  getPasswordLabel(): string;
  setPasswordLabel(value: string): PasswordScreenText;

  getResetLinkText(): string;
  setResetLinkText(value: string): PasswordScreenText;

  getBackButtonText(): string;
  setBackButtonText(value: string): PasswordScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): PasswordScreenText;

  getMinLength(): string;
  setMinLength(value: string): PasswordScreenText;

  getHasUppercase(): string;
  setHasUppercase(value: string): PasswordScreenText;

  getHasLowercase(): string;
  setHasLowercase(value: string): PasswordScreenText;

  getHasNumber(): string;
  setHasNumber(value: string): PasswordScreenText;

  getHasSymbol(): string;
  setHasSymbol(value: string): PasswordScreenText;

  getConfirmation(): string;
  setConfirmation(value: string): PasswordScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordScreenText): PasswordScreenText.AsObject;
  static serializeBinaryToWriter(message: PasswordScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordScreenText;
  static deserializeBinaryFromReader(message: PasswordScreenText, reader: jspb.BinaryReader): PasswordScreenText;
}

export namespace PasswordScreenText {
  export type AsObject = {
    title: string,
    description: string,
    passwordLabel: string,
    resetLinkText: string,
    backButtonText: string,
    nextButtonText: string,
    minLength: string,
    hasUppercase: string,
    hasLowercase: string,
    hasNumber: string,
    hasSymbol: string,
    confirmation: string,
  }
}

export class UsernameChangeScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): UsernameChangeScreenText;

  getDescription(): string;
  setDescription(value: string): UsernameChangeScreenText;

  getUsernameLabel(): string;
  setUsernameLabel(value: string): UsernameChangeScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): UsernameChangeScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): UsernameChangeScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UsernameChangeScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: UsernameChangeScreenText): UsernameChangeScreenText.AsObject;
  static serializeBinaryToWriter(message: UsernameChangeScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UsernameChangeScreenText;
  static deserializeBinaryFromReader(message: UsernameChangeScreenText, reader: jspb.BinaryReader): UsernameChangeScreenText;
}

export namespace UsernameChangeScreenText {
  export type AsObject = {
    title: string,
    description: string,
    usernameLabel: string,
    cancelButtonText: string,
    nextButtonText: string,
  }
}

export class UsernameChangeDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): UsernameChangeDoneScreenText;

  getDescription(): string;
  setDescription(value: string): UsernameChangeDoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): UsernameChangeDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UsernameChangeDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: UsernameChangeDoneScreenText): UsernameChangeDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: UsernameChangeDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UsernameChangeDoneScreenText;
  static deserializeBinaryFromReader(message: UsernameChangeDoneScreenText, reader: jspb.BinaryReader): UsernameChangeDoneScreenText;
}

export namespace UsernameChangeDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    nextButtonText: string,
  }
}

export class InitPasswordScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): InitPasswordScreenText;

  getDescription(): string;
  setDescription(value: string): InitPasswordScreenText;

  getCodeLabel(): string;
  setCodeLabel(value: string): InitPasswordScreenText;

  getNewPasswordLabel(): string;
  setNewPasswordLabel(value: string): InitPasswordScreenText;

  getNewPasswordConfirmLabel(): string;
  setNewPasswordConfirmLabel(value: string): InitPasswordScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): InitPasswordScreenText;

  getResendButtonText(): string;
  setResendButtonText(value: string): InitPasswordScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitPasswordScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: InitPasswordScreenText): InitPasswordScreenText.AsObject;
  static serializeBinaryToWriter(message: InitPasswordScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitPasswordScreenText;
  static deserializeBinaryFromReader(message: InitPasswordScreenText, reader: jspb.BinaryReader): InitPasswordScreenText;
}

export namespace InitPasswordScreenText {
  export type AsObject = {
    title: string,
    description: string,
    codeLabel: string,
    newPasswordLabel: string,
    newPasswordConfirmLabel: string,
    nextButtonText: string,
    resendButtonText: string,
  }
}

export class InitPasswordDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): InitPasswordDoneScreenText;

  getDescription(): string;
  setDescription(value: string): InitPasswordDoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): InitPasswordDoneScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): InitPasswordDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitPasswordDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: InitPasswordDoneScreenText): InitPasswordDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: InitPasswordDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitPasswordDoneScreenText;
  static deserializeBinaryFromReader(message: InitPasswordDoneScreenText, reader: jspb.BinaryReader): InitPasswordDoneScreenText;
}

export namespace InitPasswordDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    nextButtonText: string,
    cancelButtonText: string,
  }
}

export class EmailVerificationScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): EmailVerificationScreenText;

  getDescription(): string;
  setDescription(value: string): EmailVerificationScreenText;

  getCodeLabel(): string;
  setCodeLabel(value: string): EmailVerificationScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): EmailVerificationScreenText;

  getResendButtonText(): string;
  setResendButtonText(value: string): EmailVerificationScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmailVerificationScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: EmailVerificationScreenText): EmailVerificationScreenText.AsObject;
  static serializeBinaryToWriter(message: EmailVerificationScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmailVerificationScreenText;
  static deserializeBinaryFromReader(message: EmailVerificationScreenText, reader: jspb.BinaryReader): EmailVerificationScreenText;
}

export namespace EmailVerificationScreenText {
  export type AsObject = {
    title: string,
    description: string,
    codeLabel: string,
    nextButtonText: string,
    resendButtonText: string,
  }
}

export class EmailVerificationDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): EmailVerificationDoneScreenText;

  getDescription(): string;
  setDescription(value: string): EmailVerificationDoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): EmailVerificationDoneScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): EmailVerificationDoneScreenText;

  getLoginButtonText(): string;
  setLoginButtonText(value: string): EmailVerificationDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmailVerificationDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: EmailVerificationDoneScreenText): EmailVerificationDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: EmailVerificationDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmailVerificationDoneScreenText;
  static deserializeBinaryFromReader(message: EmailVerificationDoneScreenText, reader: jspb.BinaryReader): EmailVerificationDoneScreenText;
}

export namespace EmailVerificationDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    nextButtonText: string,
    cancelButtonText: string,
    loginButtonText: string,
  }
}

export class InitializeUserScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): InitializeUserScreenText;

  getDescription(): string;
  setDescription(value: string): InitializeUserScreenText;

  getCodeLabel(): string;
  setCodeLabel(value: string): InitializeUserScreenText;

  getNewPasswordLabel(): string;
  setNewPasswordLabel(value: string): InitializeUserScreenText;

  getNewPasswordConfirmLabel(): string;
  setNewPasswordConfirmLabel(value: string): InitializeUserScreenText;

  getResendButtonText(): string;
  setResendButtonText(value: string): InitializeUserScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): InitializeUserScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitializeUserScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: InitializeUserScreenText): InitializeUserScreenText.AsObject;
  static serializeBinaryToWriter(message: InitializeUserScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitializeUserScreenText;
  static deserializeBinaryFromReader(message: InitializeUserScreenText, reader: jspb.BinaryReader): InitializeUserScreenText;
}

export namespace InitializeUserScreenText {
  export type AsObject = {
    title: string,
    description: string,
    codeLabel: string,
    newPasswordLabel: string,
    newPasswordConfirmLabel: string,
    resendButtonText: string,
    nextButtonText: string,
  }
}

export class InitializeUserDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): InitializeUserDoneScreenText;

  getDescription(): string;
  setDescription(value: string): InitializeUserDoneScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): InitializeUserDoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): InitializeUserDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitializeUserDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: InitializeUserDoneScreenText): InitializeUserDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: InitializeUserDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitializeUserDoneScreenText;
  static deserializeBinaryFromReader(message: InitializeUserDoneScreenText, reader: jspb.BinaryReader): InitializeUserDoneScreenText;
}

export namespace InitializeUserDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    cancelButtonText: string,
    nextButtonText: string,
  }
}

export class InitMFAPromptScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): InitMFAPromptScreenText;

  getDescription(): string;
  setDescription(value: string): InitMFAPromptScreenText;

  getOtpOption(): string;
  setOtpOption(value: string): InitMFAPromptScreenText;

  getU2fOption(): string;
  setU2fOption(value: string): InitMFAPromptScreenText;

  getSkipButtonText(): string;
  setSkipButtonText(value: string): InitMFAPromptScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): InitMFAPromptScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitMFAPromptScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: InitMFAPromptScreenText): InitMFAPromptScreenText.AsObject;
  static serializeBinaryToWriter(message: InitMFAPromptScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitMFAPromptScreenText;
  static deserializeBinaryFromReader(message: InitMFAPromptScreenText, reader: jspb.BinaryReader): InitMFAPromptScreenText;
}

export namespace InitMFAPromptScreenText {
  export type AsObject = {
    title: string,
    description: string,
    otpOption: string,
    u2fOption: string,
    skipButtonText: string,
    nextButtonText: string,
  }
}

export class InitMFAOTPScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): InitMFAOTPScreenText;

  getDescription(): string;
  setDescription(value: string): InitMFAOTPScreenText;

  getDescriptionOtp(): string;
  setDescriptionOtp(value: string): InitMFAOTPScreenText;

  getSecretLabel(): string;
  setSecretLabel(value: string): InitMFAOTPScreenText;

  getCodeLabel(): string;
  setCodeLabel(value: string): InitMFAOTPScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): InitMFAOTPScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): InitMFAOTPScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitMFAOTPScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: InitMFAOTPScreenText): InitMFAOTPScreenText.AsObject;
  static serializeBinaryToWriter(message: InitMFAOTPScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitMFAOTPScreenText;
  static deserializeBinaryFromReader(message: InitMFAOTPScreenText, reader: jspb.BinaryReader): InitMFAOTPScreenText;
}

export namespace InitMFAOTPScreenText {
  export type AsObject = {
    title: string,
    description: string,
    descriptionOtp: string,
    secretLabel: string,
    codeLabel: string,
    nextButtonText: string,
    cancelButtonText: string,
  }
}

export class InitMFAU2FScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): InitMFAU2FScreenText;

  getDescription(): string;
  setDescription(value: string): InitMFAU2FScreenText;

  getTokenNameLabel(): string;
  setTokenNameLabel(value: string): InitMFAU2FScreenText;

  getNotSupported(): string;
  setNotSupported(value: string): InitMFAU2FScreenText;

  getRegisterTokenButtonText(): string;
  setRegisterTokenButtonText(value: string): InitMFAU2FScreenText;

  getErrorRetry(): string;
  setErrorRetry(value: string): InitMFAU2FScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitMFAU2FScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: InitMFAU2FScreenText): InitMFAU2FScreenText.AsObject;
  static serializeBinaryToWriter(message: InitMFAU2FScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitMFAU2FScreenText;
  static deserializeBinaryFromReader(message: InitMFAU2FScreenText, reader: jspb.BinaryReader): InitMFAU2FScreenText;
}

export namespace InitMFAU2FScreenText {
  export type AsObject = {
    title: string,
    description: string,
    tokenNameLabel: string,
    notSupported: string,
    registerTokenButtonText: string,
    errorRetry: string,
  }
}

export class InitMFADoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): InitMFADoneScreenText;

  getDescription(): string;
  setDescription(value: string): InitMFADoneScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): InitMFADoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): InitMFADoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitMFADoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: InitMFADoneScreenText): InitMFADoneScreenText.AsObject;
  static serializeBinaryToWriter(message: InitMFADoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitMFADoneScreenText;
  static deserializeBinaryFromReader(message: InitMFADoneScreenText, reader: jspb.BinaryReader): InitMFADoneScreenText;
}

export namespace InitMFADoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    cancelButtonText: string,
    nextButtonText: string,
  }
}

export class MFAProvidersText extends jspb.Message {
  getChooseOther(): string;
  setChooseOther(value: string): MFAProvidersText;

  getOtp(): string;
  setOtp(value: string): MFAProvidersText;

  getU2f(): string;
  setU2f(value: string): MFAProvidersText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MFAProvidersText.AsObject;
  static toObject(includeInstance: boolean, msg: MFAProvidersText): MFAProvidersText.AsObject;
  static serializeBinaryToWriter(message: MFAProvidersText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MFAProvidersText;
  static deserializeBinaryFromReader(message: MFAProvidersText, reader: jspb.BinaryReader): MFAProvidersText;
}

export namespace MFAProvidersText {
  export type AsObject = {
    chooseOther: string,
    otp: string,
    u2f: string,
  }
}

export class VerifyMFAOTPScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): VerifyMFAOTPScreenText;

  getDescription(): string;
  setDescription(value: string): VerifyMFAOTPScreenText;

  getCodeLabel(): string;
  setCodeLabel(value: string): VerifyMFAOTPScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): VerifyMFAOTPScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMFAOTPScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMFAOTPScreenText): VerifyMFAOTPScreenText.AsObject;
  static serializeBinaryToWriter(message: VerifyMFAOTPScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMFAOTPScreenText;
  static deserializeBinaryFromReader(message: VerifyMFAOTPScreenText, reader: jspb.BinaryReader): VerifyMFAOTPScreenText;
}

export namespace VerifyMFAOTPScreenText {
  export type AsObject = {
    title: string,
    description: string,
    codeLabel: string,
    nextButtonText: string,
  }
}

export class VerifyMFAU2FScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): VerifyMFAU2FScreenText;

  getDescription(): string;
  setDescription(value: string): VerifyMFAU2FScreenText;

  getValidateTokenText(): string;
  setValidateTokenText(value: string): VerifyMFAU2FScreenText;

  getNotSupported(): string;
  setNotSupported(value: string): VerifyMFAU2FScreenText;

  getErrorRetry(): string;
  setErrorRetry(value: string): VerifyMFAU2FScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VerifyMFAU2FScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: VerifyMFAU2FScreenText): VerifyMFAU2FScreenText.AsObject;
  static serializeBinaryToWriter(message: VerifyMFAU2FScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VerifyMFAU2FScreenText;
  static deserializeBinaryFromReader(message: VerifyMFAU2FScreenText, reader: jspb.BinaryReader): VerifyMFAU2FScreenText;
}

export namespace VerifyMFAU2FScreenText {
  export type AsObject = {
    title: string,
    description: string,
    validateTokenText: string,
    notSupported: string,
    errorRetry: string,
  }
}

export class PasswordlessScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): PasswordlessScreenText;

  getDescription(): string;
  setDescription(value: string): PasswordlessScreenText;

  getLoginWithPwButtonText(): string;
  setLoginWithPwButtonText(value: string): PasswordlessScreenText;

  getValidateTokenButtonText(): string;
  setValidateTokenButtonText(value: string): PasswordlessScreenText;

  getNotSupported(): string;
  setNotSupported(value: string): PasswordlessScreenText;

  getErrorRetry(): string;
  setErrorRetry(value: string): PasswordlessScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordlessScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordlessScreenText): PasswordlessScreenText.AsObject;
  static serializeBinaryToWriter(message: PasswordlessScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordlessScreenText;
  static deserializeBinaryFromReader(message: PasswordlessScreenText, reader: jspb.BinaryReader): PasswordlessScreenText;
}

export namespace PasswordlessScreenText {
  export type AsObject = {
    title: string,
    description: string,
    loginWithPwButtonText: string,
    validateTokenButtonText: string,
    notSupported: string,
    errorRetry: string,
  }
}

export class PasswordChangeScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): PasswordChangeScreenText;

  getDescription(): string;
  setDescription(value: string): PasswordChangeScreenText;

  getOldPasswordLabel(): string;
  setOldPasswordLabel(value: string): PasswordChangeScreenText;

  getNewPasswordLabel(): string;
  setNewPasswordLabel(value: string): PasswordChangeScreenText;

  getNewPasswordConfirmLabel(): string;
  setNewPasswordConfirmLabel(value: string): PasswordChangeScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): PasswordChangeScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): PasswordChangeScreenText;

  getExpiredDescription(): string;
  setExpiredDescription(value: string): PasswordChangeScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordChangeScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordChangeScreenText): PasswordChangeScreenText.AsObject;
  static serializeBinaryToWriter(message: PasswordChangeScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordChangeScreenText;
  static deserializeBinaryFromReader(message: PasswordChangeScreenText, reader: jspb.BinaryReader): PasswordChangeScreenText;
}

export namespace PasswordChangeScreenText {
  export type AsObject = {
    title: string,
    description: string,
    oldPasswordLabel: string,
    newPasswordLabel: string,
    newPasswordConfirmLabel: string,
    cancelButtonText: string,
    nextButtonText: string,
    expiredDescription: string,
  }
}

export class PasswordChangeDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): PasswordChangeDoneScreenText;

  getDescription(): string;
  setDescription(value: string): PasswordChangeDoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): PasswordChangeDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordChangeDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordChangeDoneScreenText): PasswordChangeDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: PasswordChangeDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordChangeDoneScreenText;
  static deserializeBinaryFromReader(message: PasswordChangeDoneScreenText, reader: jspb.BinaryReader): PasswordChangeDoneScreenText;
}

export namespace PasswordChangeDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    nextButtonText: string,
  }
}

export class PasswordResetDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): PasswordResetDoneScreenText;

  getDescription(): string;
  setDescription(value: string): PasswordResetDoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): PasswordResetDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordResetDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordResetDoneScreenText): PasswordResetDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: PasswordResetDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordResetDoneScreenText;
  static deserializeBinaryFromReader(message: PasswordResetDoneScreenText, reader: jspb.BinaryReader): PasswordResetDoneScreenText;
}

export namespace PasswordResetDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    nextButtonText: string,
  }
}

export class RegistrationOptionScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): RegistrationOptionScreenText;

  getDescription(): string;
  setDescription(value: string): RegistrationOptionScreenText;

  getUserNameButtonText(): string;
  setUserNameButtonText(value: string): RegistrationOptionScreenText;

  getExternalLoginDescription(): string;
  setExternalLoginDescription(value: string): RegistrationOptionScreenText;

  getLoginButtonText(): string;
  setLoginButtonText(value: string): RegistrationOptionScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegistrationOptionScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: RegistrationOptionScreenText): RegistrationOptionScreenText.AsObject;
  static serializeBinaryToWriter(message: RegistrationOptionScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegistrationOptionScreenText;
  static deserializeBinaryFromReader(message: RegistrationOptionScreenText, reader: jspb.BinaryReader): RegistrationOptionScreenText;
}

export namespace RegistrationOptionScreenText {
  export type AsObject = {
    title: string,
    description: string,
    userNameButtonText: string,
    externalLoginDescription: string,
    loginButtonText: string,
  }
}

export class RegistrationUserScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): RegistrationUserScreenText;

  getDescription(): string;
  setDescription(value: string): RegistrationUserScreenText;

  getDescriptionOrgRegister(): string;
  setDescriptionOrgRegister(value: string): RegistrationUserScreenText;

  getFirstnameLabel(): string;
  setFirstnameLabel(value: string): RegistrationUserScreenText;

  getLastnameLabel(): string;
  setLastnameLabel(value: string): RegistrationUserScreenText;

  getEmailLabel(): string;
  setEmailLabel(value: string): RegistrationUserScreenText;

  getUsernameLabel(): string;
  setUsernameLabel(value: string): RegistrationUserScreenText;

  getLanguageLabel(): string;
  setLanguageLabel(value: string): RegistrationUserScreenText;

  getGenderLabel(): string;
  setGenderLabel(value: string): RegistrationUserScreenText;

  getPasswordLabel(): string;
  setPasswordLabel(value: string): RegistrationUserScreenText;

  getPasswordConfirmLabel(): string;
  setPasswordConfirmLabel(value: string): RegistrationUserScreenText;

  getTosAndPrivacyLabel(): string;
  setTosAndPrivacyLabel(value: string): RegistrationUserScreenText;

  getTosConfirm(): string;
  setTosConfirm(value: string): RegistrationUserScreenText;

  getTosLinkText(): string;
  setTosLinkText(value: string): RegistrationUserScreenText;

  getPrivacyConfirm(): string;
  setPrivacyConfirm(value: string): RegistrationUserScreenText;

  getPrivacyLinkText(): string;
  setPrivacyLinkText(value: string): RegistrationUserScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): RegistrationUserScreenText;

  getBackButtonText(): string;
  setBackButtonText(value: string): RegistrationUserScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegistrationUserScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: RegistrationUserScreenText): RegistrationUserScreenText.AsObject;
  static serializeBinaryToWriter(message: RegistrationUserScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegistrationUserScreenText;
  static deserializeBinaryFromReader(message: RegistrationUserScreenText, reader: jspb.BinaryReader): RegistrationUserScreenText;
}

export namespace RegistrationUserScreenText {
  export type AsObject = {
    title: string,
    description: string,
    descriptionOrgRegister: string,
    firstnameLabel: string,
    lastnameLabel: string,
    emailLabel: string,
    usernameLabel: string,
    languageLabel: string,
    genderLabel: string,
    passwordLabel: string,
    passwordConfirmLabel: string,
    tosAndPrivacyLabel: string,
    tosConfirm: string,
    tosLinkText: string,
    privacyConfirm: string,
    privacyLinkText: string,
    nextButtonText: string,
    backButtonText: string,
  }
}

export class ExternalRegistrationUserOverviewScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): ExternalRegistrationUserOverviewScreenText;

  getDescription(): string;
  setDescription(value: string): ExternalRegistrationUserOverviewScreenText;

  getEmailLabel(): string;
  setEmailLabel(value: string): ExternalRegistrationUserOverviewScreenText;

  getUsernameLabel(): string;
  setUsernameLabel(value: string): ExternalRegistrationUserOverviewScreenText;

  getFirstnameLabel(): string;
  setFirstnameLabel(value: string): ExternalRegistrationUserOverviewScreenText;

  getLastnameLabel(): string;
  setLastnameLabel(value: string): ExternalRegistrationUserOverviewScreenText;

  getNicknameLabel(): string;
  setNicknameLabel(value: string): ExternalRegistrationUserOverviewScreenText;

  getLanguageLabel(): string;
  setLanguageLabel(value: string): ExternalRegistrationUserOverviewScreenText;

  getPhoneLabel(): string;
  setPhoneLabel(value: string): ExternalRegistrationUserOverviewScreenText;

  getTosAndPrivacyLabel(): string;
  setTosAndPrivacyLabel(value: string): ExternalRegistrationUserOverviewScreenText;

  getTosConfirm(): string;
  setTosConfirm(value: string): ExternalRegistrationUserOverviewScreenText;

  getTosLinkText(): string;
  setTosLinkText(value: string): ExternalRegistrationUserOverviewScreenText;

  getPrivacyLinkText(): string;
  setPrivacyLinkText(value: string): ExternalRegistrationUserOverviewScreenText;

  getBackButtonText(): string;
  setBackButtonText(value: string): ExternalRegistrationUserOverviewScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): ExternalRegistrationUserOverviewScreenText;

  getPrivacyConfirm(): string;
  setPrivacyConfirm(value: string): ExternalRegistrationUserOverviewScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExternalRegistrationUserOverviewScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: ExternalRegistrationUserOverviewScreenText): ExternalRegistrationUserOverviewScreenText.AsObject;
  static serializeBinaryToWriter(message: ExternalRegistrationUserOverviewScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExternalRegistrationUserOverviewScreenText;
  static deserializeBinaryFromReader(message: ExternalRegistrationUserOverviewScreenText, reader: jspb.BinaryReader): ExternalRegistrationUserOverviewScreenText;
}

export namespace ExternalRegistrationUserOverviewScreenText {
  export type AsObject = {
    title: string,
    description: string,
    emailLabel: string,
    usernameLabel: string,
    firstnameLabel: string,
    lastnameLabel: string,
    nicknameLabel: string,
    languageLabel: string,
    phoneLabel: string,
    tosAndPrivacyLabel: string,
    tosConfirm: string,
    tosLinkText: string,
    privacyLinkText: string,
    backButtonText: string,
    nextButtonText: string,
    privacyConfirm: string,
  }
}

export class RegistrationOrgScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): RegistrationOrgScreenText;

  getDescription(): string;
  setDescription(value: string): RegistrationOrgScreenText;

  getOrgnameLabel(): string;
  setOrgnameLabel(value: string): RegistrationOrgScreenText;

  getFirstnameLabel(): string;
  setFirstnameLabel(value: string): RegistrationOrgScreenText;

  getLastnameLabel(): string;
  setLastnameLabel(value: string): RegistrationOrgScreenText;

  getUsernameLabel(): string;
  setUsernameLabel(value: string): RegistrationOrgScreenText;

  getEmailLabel(): string;
  setEmailLabel(value: string): RegistrationOrgScreenText;

  getPasswordLabel(): string;
  setPasswordLabel(value: string): RegistrationOrgScreenText;

  getPasswordConfirmLabel(): string;
  setPasswordConfirmLabel(value: string): RegistrationOrgScreenText;

  getTosAndPrivacyLabel(): string;
  setTosAndPrivacyLabel(value: string): RegistrationOrgScreenText;

  getTosConfirm(): string;
  setTosConfirm(value: string): RegistrationOrgScreenText;

  getTosLinkText(): string;
  setTosLinkText(value: string): RegistrationOrgScreenText;

  getPrivacyConfirm(): string;
  setPrivacyConfirm(value: string): RegistrationOrgScreenText;

  getPrivacyLinkText(): string;
  setPrivacyLinkText(value: string): RegistrationOrgScreenText;

  getSaveButtonText(): string;
  setSaveButtonText(value: string): RegistrationOrgScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegistrationOrgScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: RegistrationOrgScreenText): RegistrationOrgScreenText.AsObject;
  static serializeBinaryToWriter(message: RegistrationOrgScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegistrationOrgScreenText;
  static deserializeBinaryFromReader(message: RegistrationOrgScreenText, reader: jspb.BinaryReader): RegistrationOrgScreenText;
}

export namespace RegistrationOrgScreenText {
  export type AsObject = {
    title: string,
    description: string,
    orgnameLabel: string,
    firstnameLabel: string,
    lastnameLabel: string,
    usernameLabel: string,
    emailLabel: string,
    passwordLabel: string,
    passwordConfirmLabel: string,
    tosAndPrivacyLabel: string,
    tosConfirm: string,
    tosLinkText: string,
    privacyConfirm: string,
    privacyLinkText: string,
    saveButtonText: string,
  }
}

export class LinkingUserPromptScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): LinkingUserPromptScreenText;

  getDescription(): string;
  setDescription(value: string): LinkingUserPromptScreenText;

  getLinkButtonText(): string;
  setLinkButtonText(value: string): LinkingUserPromptScreenText;

  getOtherButtonText(): string;
  setOtherButtonText(value: string): LinkingUserPromptScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LinkingUserPromptScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: LinkingUserPromptScreenText): LinkingUserPromptScreenText.AsObject;
  static serializeBinaryToWriter(message: LinkingUserPromptScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LinkingUserPromptScreenText;
  static deserializeBinaryFromReader(message: LinkingUserPromptScreenText, reader: jspb.BinaryReader): LinkingUserPromptScreenText;
}

export namespace LinkingUserPromptScreenText {
  export type AsObject = {
    title: string,
    description: string,
    linkButtonText: string,
    otherButtonText: string,
  }
}

export class LinkingUserDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): LinkingUserDoneScreenText;

  getDescription(): string;
  setDescription(value: string): LinkingUserDoneScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): LinkingUserDoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): LinkingUserDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LinkingUserDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: LinkingUserDoneScreenText): LinkingUserDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: LinkingUserDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LinkingUserDoneScreenText;
  static deserializeBinaryFromReader(message: LinkingUserDoneScreenText, reader: jspb.BinaryReader): LinkingUserDoneScreenText;
}

export namespace LinkingUserDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    cancelButtonText: string,
    nextButtonText: string,
  }
}

export class ExternalUserNotFoundScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): ExternalUserNotFoundScreenText;

  getDescription(): string;
  setDescription(value: string): ExternalUserNotFoundScreenText;

  getLinkButtonText(): string;
  setLinkButtonText(value: string): ExternalUserNotFoundScreenText;

  getAutoRegisterButtonText(): string;
  setAutoRegisterButtonText(value: string): ExternalUserNotFoundScreenText;

  getTosAndPrivacyLabel(): string;
  setTosAndPrivacyLabel(value: string): ExternalUserNotFoundScreenText;

  getTosConfirm(): string;
  setTosConfirm(value: string): ExternalUserNotFoundScreenText;

  getTosLinkText(): string;
  setTosLinkText(value: string): ExternalUserNotFoundScreenText;

  getPrivacyLinkText(): string;
  setPrivacyLinkText(value: string): ExternalUserNotFoundScreenText;

  getPrivacyConfirm(): string;
  setPrivacyConfirm(value: string): ExternalUserNotFoundScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExternalUserNotFoundScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: ExternalUserNotFoundScreenText): ExternalUserNotFoundScreenText.AsObject;
  static serializeBinaryToWriter(message: ExternalUserNotFoundScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExternalUserNotFoundScreenText;
  static deserializeBinaryFromReader(message: ExternalUserNotFoundScreenText, reader: jspb.BinaryReader): ExternalUserNotFoundScreenText;
}

export namespace ExternalUserNotFoundScreenText {
  export type AsObject = {
    title: string,
    description: string,
    linkButtonText: string,
    autoRegisterButtonText: string,
    tosAndPrivacyLabel: string,
    tosConfirm: string,
    tosLinkText: string,
    privacyLinkText: string,
    privacyConfirm: string,
  }
}

export class SuccessLoginScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): SuccessLoginScreenText;

  getAutoRedirectDescription(): string;
  setAutoRedirectDescription(value: string): SuccessLoginScreenText;

  getRedirectedDescription(): string;
  setRedirectedDescription(value: string): SuccessLoginScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): SuccessLoginScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SuccessLoginScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: SuccessLoginScreenText): SuccessLoginScreenText.AsObject;
  static serializeBinaryToWriter(message: SuccessLoginScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SuccessLoginScreenText;
  static deserializeBinaryFromReader(message: SuccessLoginScreenText, reader: jspb.BinaryReader): SuccessLoginScreenText;
}

export namespace SuccessLoginScreenText {
  export type AsObject = {
    title: string,
    autoRedirectDescription: string,
    redirectedDescription: string,
    nextButtonText: string,
  }
}

export class LogoutDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): LogoutDoneScreenText;

  getDescription(): string;
  setDescription(value: string): LogoutDoneScreenText;

  getLoginButtonText(): string;
  setLoginButtonText(value: string): LogoutDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogoutDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: LogoutDoneScreenText): LogoutDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: LogoutDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogoutDoneScreenText;
  static deserializeBinaryFromReader(message: LogoutDoneScreenText, reader: jspb.BinaryReader): LogoutDoneScreenText;
}

export namespace LogoutDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    loginButtonText: string,
  }
}

export class FooterText extends jspb.Message {
  getTos(): string;
  setTos(value: string): FooterText;

  getPrivacyPolicy(): string;
  setPrivacyPolicy(value: string): FooterText;

  getHelp(): string;
  setHelp(value: string): FooterText;

  getSupportEmail(): string;
  setSupportEmail(value: string): FooterText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FooterText.AsObject;
  static toObject(includeInstance: boolean, msg: FooterText): FooterText.AsObject;
  static serializeBinaryToWriter(message: FooterText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FooterText;
  static deserializeBinaryFromReader(message: FooterText, reader: jspb.BinaryReader): FooterText;
}

export namespace FooterText {
  export type AsObject = {
    tos: string,
    privacyPolicy: string,
    help: string,
    supportEmail: string,
  }
}

export class PasswordlessPromptScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): PasswordlessPromptScreenText;

  getDescription(): string;
  setDescription(value: string): PasswordlessPromptScreenText;

  getDescriptionInit(): string;
  setDescriptionInit(value: string): PasswordlessPromptScreenText;

  getPasswordlessButtonText(): string;
  setPasswordlessButtonText(value: string): PasswordlessPromptScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): PasswordlessPromptScreenText;

  getSkipButtonText(): string;
  setSkipButtonText(value: string): PasswordlessPromptScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordlessPromptScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordlessPromptScreenText): PasswordlessPromptScreenText.AsObject;
  static serializeBinaryToWriter(message: PasswordlessPromptScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordlessPromptScreenText;
  static deserializeBinaryFromReader(message: PasswordlessPromptScreenText, reader: jspb.BinaryReader): PasswordlessPromptScreenText;
}

export namespace PasswordlessPromptScreenText {
  export type AsObject = {
    title: string,
    description: string,
    descriptionInit: string,
    passwordlessButtonText: string,
    nextButtonText: string,
    skipButtonText: string,
  }
}

export class PasswordlessRegistrationScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): PasswordlessRegistrationScreenText;

  getDescription(): string;
  setDescription(value: string): PasswordlessRegistrationScreenText;

  getTokenNameLabel(): string;
  setTokenNameLabel(value: string): PasswordlessRegistrationScreenText;

  getNotSupported(): string;
  setNotSupported(value: string): PasswordlessRegistrationScreenText;

  getRegisterTokenButtonText(): string;
  setRegisterTokenButtonText(value: string): PasswordlessRegistrationScreenText;

  getErrorRetry(): string;
  setErrorRetry(value: string): PasswordlessRegistrationScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordlessRegistrationScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordlessRegistrationScreenText): PasswordlessRegistrationScreenText.AsObject;
  static serializeBinaryToWriter(message: PasswordlessRegistrationScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordlessRegistrationScreenText;
  static deserializeBinaryFromReader(message: PasswordlessRegistrationScreenText, reader: jspb.BinaryReader): PasswordlessRegistrationScreenText;
}

export namespace PasswordlessRegistrationScreenText {
  export type AsObject = {
    title: string,
    description: string,
    tokenNameLabel: string,
    notSupported: string,
    registerTokenButtonText: string,
    errorRetry: string,
  }
}

export class PasswordlessRegistrationDoneScreenText extends jspb.Message {
  getTitle(): string;
  setTitle(value: string): PasswordlessRegistrationDoneScreenText;

  getDescription(): string;
  setDescription(value: string): PasswordlessRegistrationDoneScreenText;

  getNextButtonText(): string;
  setNextButtonText(value: string): PasswordlessRegistrationDoneScreenText;

  getCancelButtonText(): string;
  setCancelButtonText(value: string): PasswordlessRegistrationDoneScreenText;

  getDescriptionClose(): string;
  setDescriptionClose(value: string): PasswordlessRegistrationDoneScreenText;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PasswordlessRegistrationDoneScreenText.AsObject;
  static toObject(includeInstance: boolean, msg: PasswordlessRegistrationDoneScreenText): PasswordlessRegistrationDoneScreenText.AsObject;
  static serializeBinaryToWriter(message: PasswordlessRegistrationDoneScreenText, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PasswordlessRegistrationDoneScreenText;
  static deserializeBinaryFromReader(message: PasswordlessRegistrationDoneScreenText, reader: jspb.BinaryReader): PasswordlessRegistrationDoneScreenText;
}

export namespace PasswordlessRegistrationDoneScreenText {
  export type AsObject = {
    title: string,
    description: string,
    nextButtonText: string,
    cancelButtonText: string,
    descriptionClose: string,
  }
}


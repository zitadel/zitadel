/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { ObjectDetails } from "./object";

export const protobufPackage = "zitadel.text.v1";

export interface MessageCustomText {
  details: ObjectDetails | undefined;
  title: string;
  preHeader: string;
  subject: string;
  greeting: string;
  text: string;
  buttonText: string;
  footerText: string;
  isDefault: boolean;
}

export interface LoginCustomText {
  details: ObjectDetails | undefined;
  selectAccountText: SelectAccountScreenText | undefined;
  loginText: LoginScreenText | undefined;
  passwordText: PasswordScreenText | undefined;
  usernameChangeText: UsernameChangeScreenText | undefined;
  usernameChangeDoneText: UsernameChangeDoneScreenText | undefined;
  initPasswordText: InitPasswordScreenText | undefined;
  initPasswordDoneText: InitPasswordDoneScreenText | undefined;
  emailVerificationText: EmailVerificationScreenText | undefined;
  emailVerificationDoneText: EmailVerificationDoneScreenText | undefined;
  initializeUserText: InitializeUserScreenText | undefined;
  initializeDoneText: InitializeUserDoneScreenText | undefined;
  initMfaPromptText: InitMFAPromptScreenText | undefined;
  initMfaOtpText: InitMFAOTPScreenText | undefined;
  initMfaU2fText: InitMFAU2FScreenText | undefined;
  initMfaDoneText: InitMFADoneScreenText | undefined;
  mfaProvidersText: MFAProvidersText | undefined;
  verifyMfaOtpText: VerifyMFAOTPScreenText | undefined;
  verifyMfaU2fText: VerifyMFAU2FScreenText | undefined;
  passwordlessText: PasswordlessScreenText | undefined;
  passwordChangeText: PasswordChangeScreenText | undefined;
  passwordChangeDoneText: PasswordChangeDoneScreenText | undefined;
  passwordResetDoneText: PasswordResetDoneScreenText | undefined;
  registrationOptionText: RegistrationOptionScreenText | undefined;
  registrationUserText: RegistrationUserScreenText | undefined;
  registrationOrgText: RegistrationOrgScreenText | undefined;
  linkingUserDoneText: LinkingUserDoneScreenText | undefined;
  externalUserNotFoundText: ExternalUserNotFoundScreenText | undefined;
  successLoginText: SuccessLoginScreenText | undefined;
  logoutText: LogoutDoneScreenText | undefined;
  footerText: FooterText | undefined;
  passwordlessPromptText: PasswordlessPromptScreenText | undefined;
  passwordlessRegistrationText: PasswordlessRegistrationScreenText | undefined;
  passwordlessRegistrationDoneText: PasswordlessRegistrationDoneScreenText | undefined;
  externalRegistrationUserOverviewText: ExternalRegistrationUserOverviewScreenText | undefined;
  isDefault: boolean;
}

export interface SelectAccountScreenText {
  title: string;
  description: string;
  titleLinkingProcess: string;
  descriptionLinkingProcess: string;
  otherUser: string;
  sessionStateActive: string;
  sessionStateInactive: string;
  userMustBeMemberOfOrg: string;
}

export interface LoginScreenText {
  title: string;
  description: string;
  titleLinkingProcess: string;
  descriptionLinkingProcess: string;
  userMustBeMemberOfOrg: string;
  loginNameLabel: string;
  registerButtonText: string;
  nextButtonText: string;
  externalUserDescription: string;
  userNamePlaceholder: string;
  loginNamePlaceholder: string;
}

export interface PasswordScreenText {
  title: string;
  description: string;
  passwordLabel: string;
  resetLinkText: string;
  backButtonText: string;
  nextButtonText: string;
  minLength: string;
  hasUppercase: string;
  hasLowercase: string;
  hasNumber: string;
  hasSymbol: string;
  confirmation: string;
}

export interface UsernameChangeScreenText {
  title: string;
  description: string;
  usernameLabel: string;
  cancelButtonText: string;
  nextButtonText: string;
}

export interface UsernameChangeDoneScreenText {
  title: string;
  description: string;
  nextButtonText: string;
}

export interface InitPasswordScreenText {
  title: string;
  description: string;
  codeLabel: string;
  newPasswordLabel: string;
  newPasswordConfirmLabel: string;
  nextButtonText: string;
  resendButtonText: string;
}

export interface InitPasswordDoneScreenText {
  title: string;
  description: string;
  nextButtonText: string;
  cancelButtonText: string;
}

export interface EmailVerificationScreenText {
  title: string;
  description: string;
  codeLabel: string;
  nextButtonText: string;
  resendButtonText: string;
}

export interface EmailVerificationDoneScreenText {
  title: string;
  description: string;
  nextButtonText: string;
  cancelButtonText: string;
  loginButtonText: string;
}

export interface InitializeUserScreenText {
  title: string;
  description: string;
  codeLabel: string;
  newPasswordLabel: string;
  newPasswordConfirmLabel: string;
  resendButtonText: string;
  nextButtonText: string;
}

export interface InitializeUserDoneScreenText {
  title: string;
  description: string;
  cancelButtonText: string;
  nextButtonText: string;
}

export interface InitMFAPromptScreenText {
  title: string;
  description: string;
  otpOption: string;
  u2fOption: string;
  skipButtonText: string;
  nextButtonText: string;
}

export interface InitMFAOTPScreenText {
  title: string;
  description: string;
  descriptionOtp: string;
  secretLabel: string;
  codeLabel: string;
  nextButtonText: string;
  cancelButtonText: string;
}

export interface InitMFAU2FScreenText {
  title: string;
  description: string;
  tokenNameLabel: string;
  notSupported: string;
  registerTokenButtonText: string;
  errorRetry: string;
}

export interface InitMFADoneScreenText {
  title: string;
  description: string;
  cancelButtonText: string;
  nextButtonText: string;
}

export interface MFAProvidersText {
  chooseOther: string;
  otp: string;
  u2f: string;
}

export interface VerifyMFAOTPScreenText {
  title: string;
  description: string;
  codeLabel: string;
  nextButtonText: string;
}

export interface VerifyMFAU2FScreenText {
  title: string;
  description: string;
  validateTokenText: string;
  notSupported: string;
  errorRetry: string;
}

export interface PasswordlessScreenText {
  title: string;
  description: string;
  loginWithPwButtonText: string;
  validateTokenButtonText: string;
  notSupported: string;
  errorRetry: string;
}

export interface PasswordChangeScreenText {
  title: string;
  description: string;
  oldPasswordLabel: string;
  newPasswordLabel: string;
  newPasswordConfirmLabel: string;
  cancelButtonText: string;
  nextButtonText: string;
}

export interface PasswordChangeDoneScreenText {
  title: string;
  description: string;
  nextButtonText: string;
}

export interface PasswordResetDoneScreenText {
  title: string;
  description: string;
  nextButtonText: string;
}

export interface RegistrationOptionScreenText {
  title: string;
  description: string;
  userNameButtonText: string;
  externalLoginDescription: string;
  loginButtonText: string;
}

export interface RegistrationUserScreenText {
  title: string;
  description: string;
  descriptionOrgRegister: string;
  firstnameLabel: string;
  lastnameLabel: string;
  emailLabel: string;
  usernameLabel: string;
  languageLabel: string;
  genderLabel: string;
  passwordLabel: string;
  passwordConfirmLabel: string;
  tosAndPrivacyLabel: string;
  tosConfirm: string;
  tosLinkText: string;
  privacyConfirm: string;
  privacyLinkText: string;
  nextButtonText: string;
  backButtonText: string;
}

export interface ExternalRegistrationUserOverviewScreenText {
  title: string;
  description: string;
  emailLabel: string;
  usernameLabel: string;
  firstnameLabel: string;
  lastnameLabel: string;
  nicknameLabel: string;
  languageLabel: string;
  phoneLabel: string;
  tosAndPrivacyLabel: string;
  tosConfirm: string;
  tosLinkText: string;
  privacyLinkText: string;
  backButtonText: string;
  nextButtonText: string;
  privacyConfirm: string;
}

export interface RegistrationOrgScreenText {
  title: string;
  description: string;
  orgnameLabel: string;
  firstnameLabel: string;
  lastnameLabel: string;
  usernameLabel: string;
  emailLabel: string;
  passwordLabel: string;
  passwordConfirmLabel: string;
  tosAndPrivacyLabel: string;
  tosConfirm: string;
  tosLinkText: string;
  privacyConfirm: string;
  privacyLinkText: string;
  saveButtonText: string;
}

export interface LinkingUserDoneScreenText {
  title: string;
  description: string;
  cancelButtonText: string;
  nextButtonText: string;
}

export interface ExternalUserNotFoundScreenText {
  title: string;
  description: string;
  linkButtonText: string;
  autoRegisterButtonText: string;
  tosAndPrivacyLabel: string;
  tosConfirm: string;
  tosLinkText: string;
  privacyLinkText: string;
  privacyConfirm: string;
}

export interface SuccessLoginScreenText {
  title: string;
  /** Text to describe that auto-redirect should happen after successful login */
  autoRedirectDescription: string;
  /** Text to describe that the window can be closed after redirect */
  redirectedDescription: string;
  nextButtonText: string;
}

export interface LogoutDoneScreenText {
  title: string;
  description: string;
  loginButtonText: string;
}

export interface FooterText {
  tos: string;
  privacyPolicy: string;
  help: string;
}

export interface PasswordlessPromptScreenText {
  title: string;
  description: string;
  descriptionInit: string;
  passwordlessButtonText: string;
  nextButtonText: string;
  skipButtonText: string;
}

export interface PasswordlessRegistrationScreenText {
  title: string;
  description: string;
  tokenNameLabel: string;
  notSupported: string;
  registerTokenButtonText: string;
  errorRetry: string;
}

export interface PasswordlessRegistrationDoneScreenText {
  title: string;
  description: string;
  nextButtonText: string;
  cancelButtonText: string;
  descriptionClose: string;
}

function createBaseMessageCustomText(): MessageCustomText {
  return {
    details: undefined,
    title: "",
    preHeader: "",
    subject: "",
    greeting: "",
    text: "",
    buttonText: "",
    footerText: "",
    isDefault: false,
  };
}

export const MessageCustomText = {
  encode(message: MessageCustomText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.preHeader !== "") {
      writer.uint32(26).string(message.preHeader);
    }
    if (message.subject !== "") {
      writer.uint32(34).string(message.subject);
    }
    if (message.greeting !== "") {
      writer.uint32(42).string(message.greeting);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.buttonText !== "") {
      writer.uint32(58).string(message.buttonText);
    }
    if (message.footerText !== "") {
      writer.uint32(66).string(message.footerText);
    }
    if (message.isDefault === true) {
      writer.uint32(72).bool(message.isDefault);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MessageCustomText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMessageCustomText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.preHeader = reader.string();
          break;
        case 4:
          message.subject = reader.string();
          break;
        case 5:
          message.greeting = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.buttonText = reader.string();
          break;
        case 8:
          message.footerText = reader.string();
          break;
        case 9:
          message.isDefault = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MessageCustomText {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      title: isSet(object.title) ? String(object.title) : "",
      preHeader: isSet(object.preHeader) ? String(object.preHeader) : "",
      subject: isSet(object.subject) ? String(object.subject) : "",
      greeting: isSet(object.greeting) ? String(object.greeting) : "",
      text: isSet(object.text) ? String(object.text) : "",
      buttonText: isSet(object.buttonText) ? String(object.buttonText) : "",
      footerText: isSet(object.footerText) ? String(object.footerText) : "",
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
    };
  },

  toJSON(message: MessageCustomText): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.title !== undefined && (obj.title = message.title);
    message.preHeader !== undefined && (obj.preHeader = message.preHeader);
    message.subject !== undefined && (obj.subject = message.subject);
    message.greeting !== undefined && (obj.greeting = message.greeting);
    message.text !== undefined && (obj.text = message.text);
    message.buttonText !== undefined && (obj.buttonText = message.buttonText);
    message.footerText !== undefined && (obj.footerText = message.footerText);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    return obj;
  },

  create(base?: DeepPartial<MessageCustomText>): MessageCustomText {
    return MessageCustomText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<MessageCustomText>): MessageCustomText {
    const message = createBaseMessageCustomText();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.title = object.title ?? "";
    message.preHeader = object.preHeader ?? "";
    message.subject = object.subject ?? "";
    message.greeting = object.greeting ?? "";
    message.text = object.text ?? "";
    message.buttonText = object.buttonText ?? "";
    message.footerText = object.footerText ?? "";
    message.isDefault = object.isDefault ?? false;
    return message;
  },
};

function createBaseLoginCustomText(): LoginCustomText {
  return {
    details: undefined,
    selectAccountText: undefined,
    loginText: undefined,
    passwordText: undefined,
    usernameChangeText: undefined,
    usernameChangeDoneText: undefined,
    initPasswordText: undefined,
    initPasswordDoneText: undefined,
    emailVerificationText: undefined,
    emailVerificationDoneText: undefined,
    initializeUserText: undefined,
    initializeDoneText: undefined,
    initMfaPromptText: undefined,
    initMfaOtpText: undefined,
    initMfaU2fText: undefined,
    initMfaDoneText: undefined,
    mfaProvidersText: undefined,
    verifyMfaOtpText: undefined,
    verifyMfaU2fText: undefined,
    passwordlessText: undefined,
    passwordChangeText: undefined,
    passwordChangeDoneText: undefined,
    passwordResetDoneText: undefined,
    registrationOptionText: undefined,
    registrationUserText: undefined,
    registrationOrgText: undefined,
    linkingUserDoneText: undefined,
    externalUserNotFoundText: undefined,
    successLoginText: undefined,
    logoutText: undefined,
    footerText: undefined,
    passwordlessPromptText: undefined,
    passwordlessRegistrationText: undefined,
    passwordlessRegistrationDoneText: undefined,
    externalRegistrationUserOverviewText: undefined,
    isDefault: false,
  };
}

export const LoginCustomText = {
  encode(message: LoginCustomText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.selectAccountText !== undefined) {
      SelectAccountScreenText.encode(message.selectAccountText, writer.uint32(18).fork()).ldelim();
    }
    if (message.loginText !== undefined) {
      LoginScreenText.encode(message.loginText, writer.uint32(26).fork()).ldelim();
    }
    if (message.passwordText !== undefined) {
      PasswordScreenText.encode(message.passwordText, writer.uint32(34).fork()).ldelim();
    }
    if (message.usernameChangeText !== undefined) {
      UsernameChangeScreenText.encode(message.usernameChangeText, writer.uint32(42).fork()).ldelim();
    }
    if (message.usernameChangeDoneText !== undefined) {
      UsernameChangeDoneScreenText.encode(message.usernameChangeDoneText, writer.uint32(50).fork()).ldelim();
    }
    if (message.initPasswordText !== undefined) {
      InitPasswordScreenText.encode(message.initPasswordText, writer.uint32(58).fork()).ldelim();
    }
    if (message.initPasswordDoneText !== undefined) {
      InitPasswordDoneScreenText.encode(message.initPasswordDoneText, writer.uint32(66).fork()).ldelim();
    }
    if (message.emailVerificationText !== undefined) {
      EmailVerificationScreenText.encode(message.emailVerificationText, writer.uint32(74).fork()).ldelim();
    }
    if (message.emailVerificationDoneText !== undefined) {
      EmailVerificationDoneScreenText.encode(message.emailVerificationDoneText, writer.uint32(82).fork()).ldelim();
    }
    if (message.initializeUserText !== undefined) {
      InitializeUserScreenText.encode(message.initializeUserText, writer.uint32(90).fork()).ldelim();
    }
    if (message.initializeDoneText !== undefined) {
      InitializeUserDoneScreenText.encode(message.initializeDoneText, writer.uint32(98).fork()).ldelim();
    }
    if (message.initMfaPromptText !== undefined) {
      InitMFAPromptScreenText.encode(message.initMfaPromptText, writer.uint32(106).fork()).ldelim();
    }
    if (message.initMfaOtpText !== undefined) {
      InitMFAOTPScreenText.encode(message.initMfaOtpText, writer.uint32(114).fork()).ldelim();
    }
    if (message.initMfaU2fText !== undefined) {
      InitMFAU2FScreenText.encode(message.initMfaU2fText, writer.uint32(122).fork()).ldelim();
    }
    if (message.initMfaDoneText !== undefined) {
      InitMFADoneScreenText.encode(message.initMfaDoneText, writer.uint32(130).fork()).ldelim();
    }
    if (message.mfaProvidersText !== undefined) {
      MFAProvidersText.encode(message.mfaProvidersText, writer.uint32(138).fork()).ldelim();
    }
    if (message.verifyMfaOtpText !== undefined) {
      VerifyMFAOTPScreenText.encode(message.verifyMfaOtpText, writer.uint32(146).fork()).ldelim();
    }
    if (message.verifyMfaU2fText !== undefined) {
      VerifyMFAU2FScreenText.encode(message.verifyMfaU2fText, writer.uint32(154).fork()).ldelim();
    }
    if (message.passwordlessText !== undefined) {
      PasswordlessScreenText.encode(message.passwordlessText, writer.uint32(162).fork()).ldelim();
    }
    if (message.passwordChangeText !== undefined) {
      PasswordChangeScreenText.encode(message.passwordChangeText, writer.uint32(170).fork()).ldelim();
    }
    if (message.passwordChangeDoneText !== undefined) {
      PasswordChangeDoneScreenText.encode(message.passwordChangeDoneText, writer.uint32(178).fork()).ldelim();
    }
    if (message.passwordResetDoneText !== undefined) {
      PasswordResetDoneScreenText.encode(message.passwordResetDoneText, writer.uint32(186).fork()).ldelim();
    }
    if (message.registrationOptionText !== undefined) {
      RegistrationOptionScreenText.encode(message.registrationOptionText, writer.uint32(194).fork()).ldelim();
    }
    if (message.registrationUserText !== undefined) {
      RegistrationUserScreenText.encode(message.registrationUserText, writer.uint32(202).fork()).ldelim();
    }
    if (message.registrationOrgText !== undefined) {
      RegistrationOrgScreenText.encode(message.registrationOrgText, writer.uint32(210).fork()).ldelim();
    }
    if (message.linkingUserDoneText !== undefined) {
      LinkingUserDoneScreenText.encode(message.linkingUserDoneText, writer.uint32(218).fork()).ldelim();
    }
    if (message.externalUserNotFoundText !== undefined) {
      ExternalUserNotFoundScreenText.encode(message.externalUserNotFoundText, writer.uint32(226).fork()).ldelim();
    }
    if (message.successLoginText !== undefined) {
      SuccessLoginScreenText.encode(message.successLoginText, writer.uint32(234).fork()).ldelim();
    }
    if (message.logoutText !== undefined) {
      LogoutDoneScreenText.encode(message.logoutText, writer.uint32(242).fork()).ldelim();
    }
    if (message.footerText !== undefined) {
      FooterText.encode(message.footerText, writer.uint32(250).fork()).ldelim();
    }
    if (message.passwordlessPromptText !== undefined) {
      PasswordlessPromptScreenText.encode(message.passwordlessPromptText, writer.uint32(258).fork()).ldelim();
    }
    if (message.passwordlessRegistrationText !== undefined) {
      PasswordlessRegistrationScreenText.encode(message.passwordlessRegistrationText, writer.uint32(266).fork())
        .ldelim();
    }
    if (message.passwordlessRegistrationDoneText !== undefined) {
      PasswordlessRegistrationDoneScreenText.encode(message.passwordlessRegistrationDoneText, writer.uint32(274).fork())
        .ldelim();
    }
    if (message.externalRegistrationUserOverviewText !== undefined) {
      ExternalRegistrationUserOverviewScreenText.encode(
        message.externalRegistrationUserOverviewText,
        writer.uint32(282).fork(),
      ).ldelim();
    }
    if (message.isDefault === true) {
      writer.uint32(288).bool(message.isDefault);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LoginCustomText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginCustomText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.details = ObjectDetails.decode(reader, reader.uint32());
          break;
        case 2:
          message.selectAccountText = SelectAccountScreenText.decode(reader, reader.uint32());
          break;
        case 3:
          message.loginText = LoginScreenText.decode(reader, reader.uint32());
          break;
        case 4:
          message.passwordText = PasswordScreenText.decode(reader, reader.uint32());
          break;
        case 5:
          message.usernameChangeText = UsernameChangeScreenText.decode(reader, reader.uint32());
          break;
        case 6:
          message.usernameChangeDoneText = UsernameChangeDoneScreenText.decode(reader, reader.uint32());
          break;
        case 7:
          message.initPasswordText = InitPasswordScreenText.decode(reader, reader.uint32());
          break;
        case 8:
          message.initPasswordDoneText = InitPasswordDoneScreenText.decode(reader, reader.uint32());
          break;
        case 9:
          message.emailVerificationText = EmailVerificationScreenText.decode(reader, reader.uint32());
          break;
        case 10:
          message.emailVerificationDoneText = EmailVerificationDoneScreenText.decode(reader, reader.uint32());
          break;
        case 11:
          message.initializeUserText = InitializeUserScreenText.decode(reader, reader.uint32());
          break;
        case 12:
          message.initializeDoneText = InitializeUserDoneScreenText.decode(reader, reader.uint32());
          break;
        case 13:
          message.initMfaPromptText = InitMFAPromptScreenText.decode(reader, reader.uint32());
          break;
        case 14:
          message.initMfaOtpText = InitMFAOTPScreenText.decode(reader, reader.uint32());
          break;
        case 15:
          message.initMfaU2fText = InitMFAU2FScreenText.decode(reader, reader.uint32());
          break;
        case 16:
          message.initMfaDoneText = InitMFADoneScreenText.decode(reader, reader.uint32());
          break;
        case 17:
          message.mfaProvidersText = MFAProvidersText.decode(reader, reader.uint32());
          break;
        case 18:
          message.verifyMfaOtpText = VerifyMFAOTPScreenText.decode(reader, reader.uint32());
          break;
        case 19:
          message.verifyMfaU2fText = VerifyMFAU2FScreenText.decode(reader, reader.uint32());
          break;
        case 20:
          message.passwordlessText = PasswordlessScreenText.decode(reader, reader.uint32());
          break;
        case 21:
          message.passwordChangeText = PasswordChangeScreenText.decode(reader, reader.uint32());
          break;
        case 22:
          message.passwordChangeDoneText = PasswordChangeDoneScreenText.decode(reader, reader.uint32());
          break;
        case 23:
          message.passwordResetDoneText = PasswordResetDoneScreenText.decode(reader, reader.uint32());
          break;
        case 24:
          message.registrationOptionText = RegistrationOptionScreenText.decode(reader, reader.uint32());
          break;
        case 25:
          message.registrationUserText = RegistrationUserScreenText.decode(reader, reader.uint32());
          break;
        case 26:
          message.registrationOrgText = RegistrationOrgScreenText.decode(reader, reader.uint32());
          break;
        case 27:
          message.linkingUserDoneText = LinkingUserDoneScreenText.decode(reader, reader.uint32());
          break;
        case 28:
          message.externalUserNotFoundText = ExternalUserNotFoundScreenText.decode(reader, reader.uint32());
          break;
        case 29:
          message.successLoginText = SuccessLoginScreenText.decode(reader, reader.uint32());
          break;
        case 30:
          message.logoutText = LogoutDoneScreenText.decode(reader, reader.uint32());
          break;
        case 31:
          message.footerText = FooterText.decode(reader, reader.uint32());
          break;
        case 32:
          message.passwordlessPromptText = PasswordlessPromptScreenText.decode(reader, reader.uint32());
          break;
        case 33:
          message.passwordlessRegistrationText = PasswordlessRegistrationScreenText.decode(reader, reader.uint32());
          break;
        case 34:
          message.passwordlessRegistrationDoneText = PasswordlessRegistrationDoneScreenText.decode(
            reader,
            reader.uint32(),
          );
          break;
        case 35:
          message.externalRegistrationUserOverviewText = ExternalRegistrationUserOverviewScreenText.decode(
            reader,
            reader.uint32(),
          );
          break;
        case 36:
          message.isDefault = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LoginCustomText {
    return {
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      selectAccountText: isSet(object.selectAccountText)
        ? SelectAccountScreenText.fromJSON(object.selectAccountText)
        : undefined,
      loginText: isSet(object.loginText) ? LoginScreenText.fromJSON(object.loginText) : undefined,
      passwordText: isSet(object.passwordText) ? PasswordScreenText.fromJSON(object.passwordText) : undefined,
      usernameChangeText: isSet(object.usernameChangeText)
        ? UsernameChangeScreenText.fromJSON(object.usernameChangeText)
        : undefined,
      usernameChangeDoneText: isSet(object.usernameChangeDoneText)
        ? UsernameChangeDoneScreenText.fromJSON(object.usernameChangeDoneText)
        : undefined,
      initPasswordText: isSet(object.initPasswordText)
        ? InitPasswordScreenText.fromJSON(object.initPasswordText)
        : undefined,
      initPasswordDoneText: isSet(object.initPasswordDoneText)
        ? InitPasswordDoneScreenText.fromJSON(object.initPasswordDoneText)
        : undefined,
      emailVerificationText: isSet(object.emailVerificationText)
        ? EmailVerificationScreenText.fromJSON(object.emailVerificationText)
        : undefined,
      emailVerificationDoneText: isSet(object.emailVerificationDoneText)
        ? EmailVerificationDoneScreenText.fromJSON(object.emailVerificationDoneText)
        : undefined,
      initializeUserText: isSet(object.initializeUserText)
        ? InitializeUserScreenText.fromJSON(object.initializeUserText)
        : undefined,
      initializeDoneText: isSet(object.initializeDoneText)
        ? InitializeUserDoneScreenText.fromJSON(object.initializeDoneText)
        : undefined,
      initMfaPromptText: isSet(object.initMfaPromptText)
        ? InitMFAPromptScreenText.fromJSON(object.initMfaPromptText)
        : undefined,
      initMfaOtpText: isSet(object.initMfaOtpText) ? InitMFAOTPScreenText.fromJSON(object.initMfaOtpText) : undefined,
      initMfaU2fText: isSet(object.initMfaU2fText) ? InitMFAU2FScreenText.fromJSON(object.initMfaU2fText) : undefined,
      initMfaDoneText: isSet(object.initMfaDoneText)
        ? InitMFADoneScreenText.fromJSON(object.initMfaDoneText)
        : undefined,
      mfaProvidersText: isSet(object.mfaProvidersText) ? MFAProvidersText.fromJSON(object.mfaProvidersText) : undefined,
      verifyMfaOtpText: isSet(object.verifyMfaOtpText)
        ? VerifyMFAOTPScreenText.fromJSON(object.verifyMfaOtpText)
        : undefined,
      verifyMfaU2fText: isSet(object.verifyMfaU2fText)
        ? VerifyMFAU2FScreenText.fromJSON(object.verifyMfaU2fText)
        : undefined,
      passwordlessText: isSet(object.passwordlessText)
        ? PasswordlessScreenText.fromJSON(object.passwordlessText)
        : undefined,
      passwordChangeText: isSet(object.passwordChangeText)
        ? PasswordChangeScreenText.fromJSON(object.passwordChangeText)
        : undefined,
      passwordChangeDoneText: isSet(object.passwordChangeDoneText)
        ? PasswordChangeDoneScreenText.fromJSON(object.passwordChangeDoneText)
        : undefined,
      passwordResetDoneText: isSet(object.passwordResetDoneText)
        ? PasswordResetDoneScreenText.fromJSON(object.passwordResetDoneText)
        : undefined,
      registrationOptionText: isSet(object.registrationOptionText)
        ? RegistrationOptionScreenText.fromJSON(object.registrationOptionText)
        : undefined,
      registrationUserText: isSet(object.registrationUserText)
        ? RegistrationUserScreenText.fromJSON(object.registrationUserText)
        : undefined,
      registrationOrgText: isSet(object.registrationOrgText)
        ? RegistrationOrgScreenText.fromJSON(object.registrationOrgText)
        : undefined,
      linkingUserDoneText: isSet(object.linkingUserDoneText)
        ? LinkingUserDoneScreenText.fromJSON(object.linkingUserDoneText)
        : undefined,
      externalUserNotFoundText: isSet(object.externalUserNotFoundText)
        ? ExternalUserNotFoundScreenText.fromJSON(object.externalUserNotFoundText)
        : undefined,
      successLoginText: isSet(object.successLoginText)
        ? SuccessLoginScreenText.fromJSON(object.successLoginText)
        : undefined,
      logoutText: isSet(object.logoutText) ? LogoutDoneScreenText.fromJSON(object.logoutText) : undefined,
      footerText: isSet(object.footerText) ? FooterText.fromJSON(object.footerText) : undefined,
      passwordlessPromptText: isSet(object.passwordlessPromptText)
        ? PasswordlessPromptScreenText.fromJSON(object.passwordlessPromptText)
        : undefined,
      passwordlessRegistrationText: isSet(object.passwordlessRegistrationText)
        ? PasswordlessRegistrationScreenText.fromJSON(object.passwordlessRegistrationText)
        : undefined,
      passwordlessRegistrationDoneText: isSet(object.passwordlessRegistrationDoneText)
        ? PasswordlessRegistrationDoneScreenText.fromJSON(object.passwordlessRegistrationDoneText)
        : undefined,
      externalRegistrationUserOverviewText: isSet(object.externalRegistrationUserOverviewText)
        ? ExternalRegistrationUserOverviewScreenText.fromJSON(object.externalRegistrationUserOverviewText)
        : undefined,
      isDefault: isSet(object.isDefault) ? Boolean(object.isDefault) : false,
    };
  },

  toJSON(message: LoginCustomText): unknown {
    const obj: any = {};
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.selectAccountText !== undefined && (obj.selectAccountText = message.selectAccountText
      ? SelectAccountScreenText.toJSON(message.selectAccountText)
      : undefined);
    message.loginText !== undefined &&
      (obj.loginText = message.loginText ? LoginScreenText.toJSON(message.loginText) : undefined);
    message.passwordText !== undefined &&
      (obj.passwordText = message.passwordText ? PasswordScreenText.toJSON(message.passwordText) : undefined);
    message.usernameChangeText !== undefined && (obj.usernameChangeText = message.usernameChangeText
      ? UsernameChangeScreenText.toJSON(message.usernameChangeText)
      : undefined);
    message.usernameChangeDoneText !== undefined && (obj.usernameChangeDoneText = message.usernameChangeDoneText
      ? UsernameChangeDoneScreenText.toJSON(message.usernameChangeDoneText)
      : undefined);
    message.initPasswordText !== undefined && (obj.initPasswordText = message.initPasswordText
      ? InitPasswordScreenText.toJSON(message.initPasswordText)
      : undefined);
    message.initPasswordDoneText !== undefined && (obj.initPasswordDoneText = message.initPasswordDoneText
      ? InitPasswordDoneScreenText.toJSON(message.initPasswordDoneText)
      : undefined);
    message.emailVerificationText !== undefined && (obj.emailVerificationText = message.emailVerificationText
      ? EmailVerificationScreenText.toJSON(message.emailVerificationText)
      : undefined);
    message.emailVerificationDoneText !== undefined &&
      (obj.emailVerificationDoneText = message.emailVerificationDoneText
        ? EmailVerificationDoneScreenText.toJSON(message.emailVerificationDoneText)
        : undefined);
    message.initializeUserText !== undefined && (obj.initializeUserText = message.initializeUserText
      ? InitializeUserScreenText.toJSON(message.initializeUserText)
      : undefined);
    message.initializeDoneText !== undefined && (obj.initializeDoneText = message.initializeDoneText
      ? InitializeUserDoneScreenText.toJSON(message.initializeDoneText)
      : undefined);
    message.initMfaPromptText !== undefined && (obj.initMfaPromptText = message.initMfaPromptText
      ? InitMFAPromptScreenText.toJSON(message.initMfaPromptText)
      : undefined);
    message.initMfaOtpText !== undefined &&
      (obj.initMfaOtpText = message.initMfaOtpText ? InitMFAOTPScreenText.toJSON(message.initMfaOtpText) : undefined);
    message.initMfaU2fText !== undefined &&
      (obj.initMfaU2fText = message.initMfaU2fText ? InitMFAU2FScreenText.toJSON(message.initMfaU2fText) : undefined);
    message.initMfaDoneText !== undefined &&
      (obj.initMfaDoneText = message.initMfaDoneText
        ? InitMFADoneScreenText.toJSON(message.initMfaDoneText)
        : undefined);
    message.mfaProvidersText !== undefined &&
      (obj.mfaProvidersText = message.mfaProvidersText ? MFAProvidersText.toJSON(message.mfaProvidersText) : undefined);
    message.verifyMfaOtpText !== undefined && (obj.verifyMfaOtpText = message.verifyMfaOtpText
      ? VerifyMFAOTPScreenText.toJSON(message.verifyMfaOtpText)
      : undefined);
    message.verifyMfaU2fText !== undefined && (obj.verifyMfaU2fText = message.verifyMfaU2fText
      ? VerifyMFAU2FScreenText.toJSON(message.verifyMfaU2fText)
      : undefined);
    message.passwordlessText !== undefined && (obj.passwordlessText = message.passwordlessText
      ? PasswordlessScreenText.toJSON(message.passwordlessText)
      : undefined);
    message.passwordChangeText !== undefined && (obj.passwordChangeText = message.passwordChangeText
      ? PasswordChangeScreenText.toJSON(message.passwordChangeText)
      : undefined);
    message.passwordChangeDoneText !== undefined && (obj.passwordChangeDoneText = message.passwordChangeDoneText
      ? PasswordChangeDoneScreenText.toJSON(message.passwordChangeDoneText)
      : undefined);
    message.passwordResetDoneText !== undefined && (obj.passwordResetDoneText = message.passwordResetDoneText
      ? PasswordResetDoneScreenText.toJSON(message.passwordResetDoneText)
      : undefined);
    message.registrationOptionText !== undefined && (obj.registrationOptionText = message.registrationOptionText
      ? RegistrationOptionScreenText.toJSON(message.registrationOptionText)
      : undefined);
    message.registrationUserText !== undefined && (obj.registrationUserText = message.registrationUserText
      ? RegistrationUserScreenText.toJSON(message.registrationUserText)
      : undefined);
    message.registrationOrgText !== undefined && (obj.registrationOrgText = message.registrationOrgText
      ? RegistrationOrgScreenText.toJSON(message.registrationOrgText)
      : undefined);
    message.linkingUserDoneText !== undefined && (obj.linkingUserDoneText = message.linkingUserDoneText
      ? LinkingUserDoneScreenText.toJSON(message.linkingUserDoneText)
      : undefined);
    message.externalUserNotFoundText !== undefined && (obj.externalUserNotFoundText = message.externalUserNotFoundText
      ? ExternalUserNotFoundScreenText.toJSON(message.externalUserNotFoundText)
      : undefined);
    message.successLoginText !== undefined && (obj.successLoginText = message.successLoginText
      ? SuccessLoginScreenText.toJSON(message.successLoginText)
      : undefined);
    message.logoutText !== undefined &&
      (obj.logoutText = message.logoutText ? LogoutDoneScreenText.toJSON(message.logoutText) : undefined);
    message.footerText !== undefined &&
      (obj.footerText = message.footerText ? FooterText.toJSON(message.footerText) : undefined);
    message.passwordlessPromptText !== undefined && (obj.passwordlessPromptText = message.passwordlessPromptText
      ? PasswordlessPromptScreenText.toJSON(message.passwordlessPromptText)
      : undefined);
    message.passwordlessRegistrationText !== undefined &&
      (obj.passwordlessRegistrationText = message.passwordlessRegistrationText
        ? PasswordlessRegistrationScreenText.toJSON(message.passwordlessRegistrationText)
        : undefined);
    message.passwordlessRegistrationDoneText !== undefined &&
      (obj.passwordlessRegistrationDoneText = message.passwordlessRegistrationDoneText
        ? PasswordlessRegistrationDoneScreenText.toJSON(message.passwordlessRegistrationDoneText)
        : undefined);
    message.externalRegistrationUserOverviewText !== undefined &&
      (obj.externalRegistrationUserOverviewText = message.externalRegistrationUserOverviewText
        ? ExternalRegistrationUserOverviewScreenText.toJSON(message.externalRegistrationUserOverviewText)
        : undefined);
    message.isDefault !== undefined && (obj.isDefault = message.isDefault);
    return obj;
  },

  create(base?: DeepPartial<LoginCustomText>): LoginCustomText {
    return LoginCustomText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LoginCustomText>): LoginCustomText {
    const message = createBaseLoginCustomText();
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.selectAccountText = (object.selectAccountText !== undefined && object.selectAccountText !== null)
      ? SelectAccountScreenText.fromPartial(object.selectAccountText)
      : undefined;
    message.loginText = (object.loginText !== undefined && object.loginText !== null)
      ? LoginScreenText.fromPartial(object.loginText)
      : undefined;
    message.passwordText = (object.passwordText !== undefined && object.passwordText !== null)
      ? PasswordScreenText.fromPartial(object.passwordText)
      : undefined;
    message.usernameChangeText = (object.usernameChangeText !== undefined && object.usernameChangeText !== null)
      ? UsernameChangeScreenText.fromPartial(object.usernameChangeText)
      : undefined;
    message.usernameChangeDoneText =
      (object.usernameChangeDoneText !== undefined && object.usernameChangeDoneText !== null)
        ? UsernameChangeDoneScreenText.fromPartial(object.usernameChangeDoneText)
        : undefined;
    message.initPasswordText = (object.initPasswordText !== undefined && object.initPasswordText !== null)
      ? InitPasswordScreenText.fromPartial(object.initPasswordText)
      : undefined;
    message.initPasswordDoneText = (object.initPasswordDoneText !== undefined && object.initPasswordDoneText !== null)
      ? InitPasswordDoneScreenText.fromPartial(object.initPasswordDoneText)
      : undefined;
    message.emailVerificationText =
      (object.emailVerificationText !== undefined && object.emailVerificationText !== null)
        ? EmailVerificationScreenText.fromPartial(object.emailVerificationText)
        : undefined;
    message.emailVerificationDoneText =
      (object.emailVerificationDoneText !== undefined && object.emailVerificationDoneText !== null)
        ? EmailVerificationDoneScreenText.fromPartial(object.emailVerificationDoneText)
        : undefined;
    message.initializeUserText = (object.initializeUserText !== undefined && object.initializeUserText !== null)
      ? InitializeUserScreenText.fromPartial(object.initializeUserText)
      : undefined;
    message.initializeDoneText = (object.initializeDoneText !== undefined && object.initializeDoneText !== null)
      ? InitializeUserDoneScreenText.fromPartial(object.initializeDoneText)
      : undefined;
    message.initMfaPromptText = (object.initMfaPromptText !== undefined && object.initMfaPromptText !== null)
      ? InitMFAPromptScreenText.fromPartial(object.initMfaPromptText)
      : undefined;
    message.initMfaOtpText = (object.initMfaOtpText !== undefined && object.initMfaOtpText !== null)
      ? InitMFAOTPScreenText.fromPartial(object.initMfaOtpText)
      : undefined;
    message.initMfaU2fText = (object.initMfaU2fText !== undefined && object.initMfaU2fText !== null)
      ? InitMFAU2FScreenText.fromPartial(object.initMfaU2fText)
      : undefined;
    message.initMfaDoneText = (object.initMfaDoneText !== undefined && object.initMfaDoneText !== null)
      ? InitMFADoneScreenText.fromPartial(object.initMfaDoneText)
      : undefined;
    message.mfaProvidersText = (object.mfaProvidersText !== undefined && object.mfaProvidersText !== null)
      ? MFAProvidersText.fromPartial(object.mfaProvidersText)
      : undefined;
    message.verifyMfaOtpText = (object.verifyMfaOtpText !== undefined && object.verifyMfaOtpText !== null)
      ? VerifyMFAOTPScreenText.fromPartial(object.verifyMfaOtpText)
      : undefined;
    message.verifyMfaU2fText = (object.verifyMfaU2fText !== undefined && object.verifyMfaU2fText !== null)
      ? VerifyMFAU2FScreenText.fromPartial(object.verifyMfaU2fText)
      : undefined;
    message.passwordlessText = (object.passwordlessText !== undefined && object.passwordlessText !== null)
      ? PasswordlessScreenText.fromPartial(object.passwordlessText)
      : undefined;
    message.passwordChangeText = (object.passwordChangeText !== undefined && object.passwordChangeText !== null)
      ? PasswordChangeScreenText.fromPartial(object.passwordChangeText)
      : undefined;
    message.passwordChangeDoneText =
      (object.passwordChangeDoneText !== undefined && object.passwordChangeDoneText !== null)
        ? PasswordChangeDoneScreenText.fromPartial(object.passwordChangeDoneText)
        : undefined;
    message.passwordResetDoneText =
      (object.passwordResetDoneText !== undefined && object.passwordResetDoneText !== null)
        ? PasswordResetDoneScreenText.fromPartial(object.passwordResetDoneText)
        : undefined;
    message.registrationOptionText =
      (object.registrationOptionText !== undefined && object.registrationOptionText !== null)
        ? RegistrationOptionScreenText.fromPartial(object.registrationOptionText)
        : undefined;
    message.registrationUserText = (object.registrationUserText !== undefined && object.registrationUserText !== null)
      ? RegistrationUserScreenText.fromPartial(object.registrationUserText)
      : undefined;
    message.registrationOrgText = (object.registrationOrgText !== undefined && object.registrationOrgText !== null)
      ? RegistrationOrgScreenText.fromPartial(object.registrationOrgText)
      : undefined;
    message.linkingUserDoneText = (object.linkingUserDoneText !== undefined && object.linkingUserDoneText !== null)
      ? LinkingUserDoneScreenText.fromPartial(object.linkingUserDoneText)
      : undefined;
    message.externalUserNotFoundText =
      (object.externalUserNotFoundText !== undefined && object.externalUserNotFoundText !== null)
        ? ExternalUserNotFoundScreenText.fromPartial(object.externalUserNotFoundText)
        : undefined;
    message.successLoginText = (object.successLoginText !== undefined && object.successLoginText !== null)
      ? SuccessLoginScreenText.fromPartial(object.successLoginText)
      : undefined;
    message.logoutText = (object.logoutText !== undefined && object.logoutText !== null)
      ? LogoutDoneScreenText.fromPartial(object.logoutText)
      : undefined;
    message.footerText = (object.footerText !== undefined && object.footerText !== null)
      ? FooterText.fromPartial(object.footerText)
      : undefined;
    message.passwordlessPromptText =
      (object.passwordlessPromptText !== undefined && object.passwordlessPromptText !== null)
        ? PasswordlessPromptScreenText.fromPartial(object.passwordlessPromptText)
        : undefined;
    message.passwordlessRegistrationText =
      (object.passwordlessRegistrationText !== undefined && object.passwordlessRegistrationText !== null)
        ? PasswordlessRegistrationScreenText.fromPartial(object.passwordlessRegistrationText)
        : undefined;
    message.passwordlessRegistrationDoneText =
      (object.passwordlessRegistrationDoneText !== undefined && object.passwordlessRegistrationDoneText !== null)
        ? PasswordlessRegistrationDoneScreenText.fromPartial(object.passwordlessRegistrationDoneText)
        : undefined;
    message.externalRegistrationUserOverviewText =
      (object.externalRegistrationUserOverviewText !== undefined &&
          object.externalRegistrationUserOverviewText !== null)
        ? ExternalRegistrationUserOverviewScreenText.fromPartial(object.externalRegistrationUserOverviewText)
        : undefined;
    message.isDefault = object.isDefault ?? false;
    return message;
  },
};

function createBaseSelectAccountScreenText(): SelectAccountScreenText {
  return {
    title: "",
    description: "",
    titleLinkingProcess: "",
    descriptionLinkingProcess: "",
    otherUser: "",
    sessionStateActive: "",
    sessionStateInactive: "",
    userMustBeMemberOfOrg: "",
  };
}

export const SelectAccountScreenText = {
  encode(message: SelectAccountScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.titleLinkingProcess !== "") {
      writer.uint32(26).string(message.titleLinkingProcess);
    }
    if (message.descriptionLinkingProcess !== "") {
      writer.uint32(34).string(message.descriptionLinkingProcess);
    }
    if (message.otherUser !== "") {
      writer.uint32(42).string(message.otherUser);
    }
    if (message.sessionStateActive !== "") {
      writer.uint32(50).string(message.sessionStateActive);
    }
    if (message.sessionStateInactive !== "") {
      writer.uint32(58).string(message.sessionStateInactive);
    }
    if (message.userMustBeMemberOfOrg !== "") {
      writer.uint32(66).string(message.userMustBeMemberOfOrg);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SelectAccountScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSelectAccountScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.titleLinkingProcess = reader.string();
          break;
        case 4:
          message.descriptionLinkingProcess = reader.string();
          break;
        case 5:
          message.otherUser = reader.string();
          break;
        case 6:
          message.sessionStateActive = reader.string();
          break;
        case 7:
          message.sessionStateInactive = reader.string();
          break;
        case 8:
          message.userMustBeMemberOfOrg = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SelectAccountScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      titleLinkingProcess: isSet(object.titleLinkingProcess) ? String(object.titleLinkingProcess) : "",
      descriptionLinkingProcess: isSet(object.descriptionLinkingProcess)
        ? String(object.descriptionLinkingProcess)
        : "",
      otherUser: isSet(object.otherUser) ? String(object.otherUser) : "",
      sessionStateActive: isSet(object.sessionStateActive) ? String(object.sessionStateActive) : "",
      sessionStateInactive: isSet(object.sessionStateInactive) ? String(object.sessionStateInactive) : "",
      userMustBeMemberOfOrg: isSet(object.userMustBeMemberOfOrg) ? String(object.userMustBeMemberOfOrg) : "",
    };
  },

  toJSON(message: SelectAccountScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.titleLinkingProcess !== undefined && (obj.titleLinkingProcess = message.titleLinkingProcess);
    message.descriptionLinkingProcess !== undefined &&
      (obj.descriptionLinkingProcess = message.descriptionLinkingProcess);
    message.otherUser !== undefined && (obj.otherUser = message.otherUser);
    message.sessionStateActive !== undefined && (obj.sessionStateActive = message.sessionStateActive);
    message.sessionStateInactive !== undefined && (obj.sessionStateInactive = message.sessionStateInactive);
    message.userMustBeMemberOfOrg !== undefined && (obj.userMustBeMemberOfOrg = message.userMustBeMemberOfOrg);
    return obj;
  },

  create(base?: DeepPartial<SelectAccountScreenText>): SelectAccountScreenText {
    return SelectAccountScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SelectAccountScreenText>): SelectAccountScreenText {
    const message = createBaseSelectAccountScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.titleLinkingProcess = object.titleLinkingProcess ?? "";
    message.descriptionLinkingProcess = object.descriptionLinkingProcess ?? "";
    message.otherUser = object.otherUser ?? "";
    message.sessionStateActive = object.sessionStateActive ?? "";
    message.sessionStateInactive = object.sessionStateInactive ?? "";
    message.userMustBeMemberOfOrg = object.userMustBeMemberOfOrg ?? "";
    return message;
  },
};

function createBaseLoginScreenText(): LoginScreenText {
  return {
    title: "",
    description: "",
    titleLinkingProcess: "",
    descriptionLinkingProcess: "",
    userMustBeMemberOfOrg: "",
    loginNameLabel: "",
    registerButtonText: "",
    nextButtonText: "",
    externalUserDescription: "",
    userNamePlaceholder: "",
    loginNamePlaceholder: "",
  };
}

export const LoginScreenText = {
  encode(message: LoginScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.titleLinkingProcess !== "") {
      writer.uint32(26).string(message.titleLinkingProcess);
    }
    if (message.descriptionLinkingProcess !== "") {
      writer.uint32(34).string(message.descriptionLinkingProcess);
    }
    if (message.userMustBeMemberOfOrg !== "") {
      writer.uint32(42).string(message.userMustBeMemberOfOrg);
    }
    if (message.loginNameLabel !== "") {
      writer.uint32(50).string(message.loginNameLabel);
    }
    if (message.registerButtonText !== "") {
      writer.uint32(58).string(message.registerButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(66).string(message.nextButtonText);
    }
    if (message.externalUserDescription !== "") {
      writer.uint32(74).string(message.externalUserDescription);
    }
    if (message.userNamePlaceholder !== "") {
      writer.uint32(82).string(message.userNamePlaceholder);
    }
    if (message.loginNamePlaceholder !== "") {
      writer.uint32(90).string(message.loginNamePlaceholder);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LoginScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.titleLinkingProcess = reader.string();
          break;
        case 4:
          message.descriptionLinkingProcess = reader.string();
          break;
        case 5:
          message.userMustBeMemberOfOrg = reader.string();
          break;
        case 6:
          message.loginNameLabel = reader.string();
          break;
        case 7:
          message.registerButtonText = reader.string();
          break;
        case 8:
          message.nextButtonText = reader.string();
          break;
        case 9:
          message.externalUserDescription = reader.string();
          break;
        case 10:
          message.userNamePlaceholder = reader.string();
          break;
        case 11:
          message.loginNamePlaceholder = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LoginScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      titleLinkingProcess: isSet(object.titleLinkingProcess) ? String(object.titleLinkingProcess) : "",
      descriptionLinkingProcess: isSet(object.descriptionLinkingProcess)
        ? String(object.descriptionLinkingProcess)
        : "",
      userMustBeMemberOfOrg: isSet(object.userMustBeMemberOfOrg) ? String(object.userMustBeMemberOfOrg) : "",
      loginNameLabel: isSet(object.loginNameLabel) ? String(object.loginNameLabel) : "",
      registerButtonText: isSet(object.registerButtonText) ? String(object.registerButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      externalUserDescription: isSet(object.externalUserDescription) ? String(object.externalUserDescription) : "",
      userNamePlaceholder: isSet(object.userNamePlaceholder) ? String(object.userNamePlaceholder) : "",
      loginNamePlaceholder: isSet(object.loginNamePlaceholder) ? String(object.loginNamePlaceholder) : "",
    };
  },

  toJSON(message: LoginScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.titleLinkingProcess !== undefined && (obj.titleLinkingProcess = message.titleLinkingProcess);
    message.descriptionLinkingProcess !== undefined &&
      (obj.descriptionLinkingProcess = message.descriptionLinkingProcess);
    message.userMustBeMemberOfOrg !== undefined && (obj.userMustBeMemberOfOrg = message.userMustBeMemberOfOrg);
    message.loginNameLabel !== undefined && (obj.loginNameLabel = message.loginNameLabel);
    message.registerButtonText !== undefined && (obj.registerButtonText = message.registerButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.externalUserDescription !== undefined && (obj.externalUserDescription = message.externalUserDescription);
    message.userNamePlaceholder !== undefined && (obj.userNamePlaceholder = message.userNamePlaceholder);
    message.loginNamePlaceholder !== undefined && (obj.loginNamePlaceholder = message.loginNamePlaceholder);
    return obj;
  },

  create(base?: DeepPartial<LoginScreenText>): LoginScreenText {
    return LoginScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LoginScreenText>): LoginScreenText {
    const message = createBaseLoginScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.titleLinkingProcess = object.titleLinkingProcess ?? "";
    message.descriptionLinkingProcess = object.descriptionLinkingProcess ?? "";
    message.userMustBeMemberOfOrg = object.userMustBeMemberOfOrg ?? "";
    message.loginNameLabel = object.loginNameLabel ?? "";
    message.registerButtonText = object.registerButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.externalUserDescription = object.externalUserDescription ?? "";
    message.userNamePlaceholder = object.userNamePlaceholder ?? "";
    message.loginNamePlaceholder = object.loginNamePlaceholder ?? "";
    return message;
  },
};

function createBasePasswordScreenText(): PasswordScreenText {
  return {
    title: "",
    description: "",
    passwordLabel: "",
    resetLinkText: "",
    backButtonText: "",
    nextButtonText: "",
    minLength: "",
    hasUppercase: "",
    hasLowercase: "",
    hasNumber: "",
    hasSymbol: "",
    confirmation: "",
  };
}

export const PasswordScreenText = {
  encode(message: PasswordScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.passwordLabel !== "") {
      writer.uint32(26).string(message.passwordLabel);
    }
    if (message.resetLinkText !== "") {
      writer.uint32(34).string(message.resetLinkText);
    }
    if (message.backButtonText !== "") {
      writer.uint32(42).string(message.backButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(50).string(message.nextButtonText);
    }
    if (message.minLength !== "") {
      writer.uint32(58).string(message.minLength);
    }
    if (message.hasUppercase !== "") {
      writer.uint32(66).string(message.hasUppercase);
    }
    if (message.hasLowercase !== "") {
      writer.uint32(74).string(message.hasLowercase);
    }
    if (message.hasNumber !== "") {
      writer.uint32(82).string(message.hasNumber);
    }
    if (message.hasSymbol !== "") {
      writer.uint32(90).string(message.hasSymbol);
    }
    if (message.confirmation !== "") {
      writer.uint32(98).string(message.confirmation);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.passwordLabel = reader.string();
          break;
        case 4:
          message.resetLinkText = reader.string();
          break;
        case 5:
          message.backButtonText = reader.string();
          break;
        case 6:
          message.nextButtonText = reader.string();
          break;
        case 7:
          message.minLength = reader.string();
          break;
        case 8:
          message.hasUppercase = reader.string();
          break;
        case 9:
          message.hasLowercase = reader.string();
          break;
        case 10:
          message.hasNumber = reader.string();
          break;
        case 11:
          message.hasSymbol = reader.string();
          break;
        case 12:
          message.confirmation = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      passwordLabel: isSet(object.passwordLabel) ? String(object.passwordLabel) : "",
      resetLinkText: isSet(object.resetLinkText) ? String(object.resetLinkText) : "",
      backButtonText: isSet(object.backButtonText) ? String(object.backButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      minLength: isSet(object.minLength) ? String(object.minLength) : "",
      hasUppercase: isSet(object.hasUppercase) ? String(object.hasUppercase) : "",
      hasLowercase: isSet(object.hasLowercase) ? String(object.hasLowercase) : "",
      hasNumber: isSet(object.hasNumber) ? String(object.hasNumber) : "",
      hasSymbol: isSet(object.hasSymbol) ? String(object.hasSymbol) : "",
      confirmation: isSet(object.confirmation) ? String(object.confirmation) : "",
    };
  },

  toJSON(message: PasswordScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.passwordLabel !== undefined && (obj.passwordLabel = message.passwordLabel);
    message.resetLinkText !== undefined && (obj.resetLinkText = message.resetLinkText);
    message.backButtonText !== undefined && (obj.backButtonText = message.backButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.minLength !== undefined && (obj.minLength = message.minLength);
    message.hasUppercase !== undefined && (obj.hasUppercase = message.hasUppercase);
    message.hasLowercase !== undefined && (obj.hasLowercase = message.hasLowercase);
    message.hasNumber !== undefined && (obj.hasNumber = message.hasNumber);
    message.hasSymbol !== undefined && (obj.hasSymbol = message.hasSymbol);
    message.confirmation !== undefined && (obj.confirmation = message.confirmation);
    return obj;
  },

  create(base?: DeepPartial<PasswordScreenText>): PasswordScreenText {
    return PasswordScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordScreenText>): PasswordScreenText {
    const message = createBasePasswordScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.passwordLabel = object.passwordLabel ?? "";
    message.resetLinkText = object.resetLinkText ?? "";
    message.backButtonText = object.backButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.minLength = object.minLength ?? "";
    message.hasUppercase = object.hasUppercase ?? "";
    message.hasLowercase = object.hasLowercase ?? "";
    message.hasNumber = object.hasNumber ?? "";
    message.hasSymbol = object.hasSymbol ?? "";
    message.confirmation = object.confirmation ?? "";
    return message;
  },
};

function createBaseUsernameChangeScreenText(): UsernameChangeScreenText {
  return { title: "", description: "", usernameLabel: "", cancelButtonText: "", nextButtonText: "" };
}

export const UsernameChangeScreenText = {
  encode(message: UsernameChangeScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.usernameLabel !== "") {
      writer.uint32(26).string(message.usernameLabel);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(34).string(message.cancelButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(42).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UsernameChangeScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUsernameChangeScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.usernameLabel = reader.string();
          break;
        case 4:
          message.cancelButtonText = reader.string();
          break;
        case 5:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UsernameChangeScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      usernameLabel: isSet(object.usernameLabel) ? String(object.usernameLabel) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: UsernameChangeScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.usernameLabel !== undefined && (obj.usernameLabel = message.usernameLabel);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<UsernameChangeScreenText>): UsernameChangeScreenText {
    return UsernameChangeScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UsernameChangeScreenText>): UsernameChangeScreenText {
    const message = createBaseUsernameChangeScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.usernameLabel = object.usernameLabel ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseUsernameChangeDoneScreenText(): UsernameChangeDoneScreenText {
  return { title: "", description: "", nextButtonText: "" };
}

export const UsernameChangeDoneScreenText = {
  encode(message: UsernameChangeDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(26).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UsernameChangeDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUsernameChangeDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UsernameChangeDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: UsernameChangeDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<UsernameChangeDoneScreenText>): UsernameChangeDoneScreenText {
    return UsernameChangeDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UsernameChangeDoneScreenText>): UsernameChangeDoneScreenText {
    const message = createBaseUsernameChangeDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseInitPasswordScreenText(): InitPasswordScreenText {
  return {
    title: "",
    description: "",
    codeLabel: "",
    newPasswordLabel: "",
    newPasswordConfirmLabel: "",
    nextButtonText: "",
    resendButtonText: "",
  };
}

export const InitPasswordScreenText = {
  encode(message: InitPasswordScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.codeLabel !== "") {
      writer.uint32(26).string(message.codeLabel);
    }
    if (message.newPasswordLabel !== "") {
      writer.uint32(34).string(message.newPasswordLabel);
    }
    if (message.newPasswordConfirmLabel !== "") {
      writer.uint32(42).string(message.newPasswordConfirmLabel);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(50).string(message.nextButtonText);
    }
    if (message.resendButtonText !== "") {
      writer.uint32(58).string(message.resendButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InitPasswordScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInitPasswordScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.codeLabel = reader.string();
          break;
        case 4:
          message.newPasswordLabel = reader.string();
          break;
        case 5:
          message.newPasswordConfirmLabel = reader.string();
          break;
        case 6:
          message.nextButtonText = reader.string();
          break;
        case 7:
          message.resendButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InitPasswordScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      codeLabel: isSet(object.codeLabel) ? String(object.codeLabel) : "",
      newPasswordLabel: isSet(object.newPasswordLabel) ? String(object.newPasswordLabel) : "",
      newPasswordConfirmLabel: isSet(object.newPasswordConfirmLabel) ? String(object.newPasswordConfirmLabel) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      resendButtonText: isSet(object.resendButtonText) ? String(object.resendButtonText) : "",
    };
  },

  toJSON(message: InitPasswordScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.codeLabel !== undefined && (obj.codeLabel = message.codeLabel);
    message.newPasswordLabel !== undefined && (obj.newPasswordLabel = message.newPasswordLabel);
    message.newPasswordConfirmLabel !== undefined && (obj.newPasswordConfirmLabel = message.newPasswordConfirmLabel);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.resendButtonText !== undefined && (obj.resendButtonText = message.resendButtonText);
    return obj;
  },

  create(base?: DeepPartial<InitPasswordScreenText>): InitPasswordScreenText {
    return InitPasswordScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InitPasswordScreenText>): InitPasswordScreenText {
    const message = createBaseInitPasswordScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.codeLabel = object.codeLabel ?? "";
    message.newPasswordLabel = object.newPasswordLabel ?? "";
    message.newPasswordConfirmLabel = object.newPasswordConfirmLabel ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.resendButtonText = object.resendButtonText ?? "";
    return message;
  },
};

function createBaseInitPasswordDoneScreenText(): InitPasswordDoneScreenText {
  return { title: "", description: "", nextButtonText: "", cancelButtonText: "" };
}

export const InitPasswordDoneScreenText = {
  encode(message: InitPasswordDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(26).string(message.nextButtonText);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(34).string(message.cancelButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InitPasswordDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInitPasswordDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.nextButtonText = reader.string();
          break;
        case 4:
          message.cancelButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InitPasswordDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
    };
  },

  toJSON(message: InitPasswordDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    return obj;
  },

  create(base?: DeepPartial<InitPasswordDoneScreenText>): InitPasswordDoneScreenText {
    return InitPasswordDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InitPasswordDoneScreenText>): InitPasswordDoneScreenText {
    const message = createBaseInitPasswordDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    return message;
  },
};

function createBaseEmailVerificationScreenText(): EmailVerificationScreenText {
  return { title: "", description: "", codeLabel: "", nextButtonText: "", resendButtonText: "" };
}

export const EmailVerificationScreenText = {
  encode(message: EmailVerificationScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.codeLabel !== "") {
      writer.uint32(26).string(message.codeLabel);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(34).string(message.nextButtonText);
    }
    if (message.resendButtonText !== "") {
      writer.uint32(42).string(message.resendButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EmailVerificationScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEmailVerificationScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.codeLabel = reader.string();
          break;
        case 4:
          message.nextButtonText = reader.string();
          break;
        case 5:
          message.resendButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EmailVerificationScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      codeLabel: isSet(object.codeLabel) ? String(object.codeLabel) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      resendButtonText: isSet(object.resendButtonText) ? String(object.resendButtonText) : "",
    };
  },

  toJSON(message: EmailVerificationScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.codeLabel !== undefined && (obj.codeLabel = message.codeLabel);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.resendButtonText !== undefined && (obj.resendButtonText = message.resendButtonText);
    return obj;
  },

  create(base?: DeepPartial<EmailVerificationScreenText>): EmailVerificationScreenText {
    return EmailVerificationScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EmailVerificationScreenText>): EmailVerificationScreenText {
    const message = createBaseEmailVerificationScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.codeLabel = object.codeLabel ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.resendButtonText = object.resendButtonText ?? "";
    return message;
  },
};

function createBaseEmailVerificationDoneScreenText(): EmailVerificationDoneScreenText {
  return { title: "", description: "", nextButtonText: "", cancelButtonText: "", loginButtonText: "" };
}

export const EmailVerificationDoneScreenText = {
  encode(message: EmailVerificationDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(26).string(message.nextButtonText);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(34).string(message.cancelButtonText);
    }
    if (message.loginButtonText !== "") {
      writer.uint32(42).string(message.loginButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EmailVerificationDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEmailVerificationDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.nextButtonText = reader.string();
          break;
        case 4:
          message.cancelButtonText = reader.string();
          break;
        case 5:
          message.loginButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EmailVerificationDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
      loginButtonText: isSet(object.loginButtonText) ? String(object.loginButtonText) : "",
    };
  },

  toJSON(message: EmailVerificationDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    message.loginButtonText !== undefined && (obj.loginButtonText = message.loginButtonText);
    return obj;
  },

  create(base?: DeepPartial<EmailVerificationDoneScreenText>): EmailVerificationDoneScreenText {
    return EmailVerificationDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EmailVerificationDoneScreenText>): EmailVerificationDoneScreenText {
    const message = createBaseEmailVerificationDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    message.loginButtonText = object.loginButtonText ?? "";
    return message;
  },
};

function createBaseInitializeUserScreenText(): InitializeUserScreenText {
  return {
    title: "",
    description: "",
    codeLabel: "",
    newPasswordLabel: "",
    newPasswordConfirmLabel: "",
    resendButtonText: "",
    nextButtonText: "",
  };
}

export const InitializeUserScreenText = {
  encode(message: InitializeUserScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.codeLabel !== "") {
      writer.uint32(26).string(message.codeLabel);
    }
    if (message.newPasswordLabel !== "") {
      writer.uint32(34).string(message.newPasswordLabel);
    }
    if (message.newPasswordConfirmLabel !== "") {
      writer.uint32(42).string(message.newPasswordConfirmLabel);
    }
    if (message.resendButtonText !== "") {
      writer.uint32(50).string(message.resendButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(58).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InitializeUserScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInitializeUserScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.codeLabel = reader.string();
          break;
        case 4:
          message.newPasswordLabel = reader.string();
          break;
        case 5:
          message.newPasswordConfirmLabel = reader.string();
          break;
        case 6:
          message.resendButtonText = reader.string();
          break;
        case 7:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InitializeUserScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      codeLabel: isSet(object.codeLabel) ? String(object.codeLabel) : "",
      newPasswordLabel: isSet(object.newPasswordLabel) ? String(object.newPasswordLabel) : "",
      newPasswordConfirmLabel: isSet(object.newPasswordConfirmLabel) ? String(object.newPasswordConfirmLabel) : "",
      resendButtonText: isSet(object.resendButtonText) ? String(object.resendButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: InitializeUserScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.codeLabel !== undefined && (obj.codeLabel = message.codeLabel);
    message.newPasswordLabel !== undefined && (obj.newPasswordLabel = message.newPasswordLabel);
    message.newPasswordConfirmLabel !== undefined && (obj.newPasswordConfirmLabel = message.newPasswordConfirmLabel);
    message.resendButtonText !== undefined && (obj.resendButtonText = message.resendButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<InitializeUserScreenText>): InitializeUserScreenText {
    return InitializeUserScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InitializeUserScreenText>): InitializeUserScreenText {
    const message = createBaseInitializeUserScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.codeLabel = object.codeLabel ?? "";
    message.newPasswordLabel = object.newPasswordLabel ?? "";
    message.newPasswordConfirmLabel = object.newPasswordConfirmLabel ?? "";
    message.resendButtonText = object.resendButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseInitializeUserDoneScreenText(): InitializeUserDoneScreenText {
  return { title: "", description: "", cancelButtonText: "", nextButtonText: "" };
}

export const InitializeUserDoneScreenText = {
  encode(message: InitializeUserDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(26).string(message.cancelButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(34).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InitializeUserDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInitializeUserDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.cancelButtonText = reader.string();
          break;
        case 4:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InitializeUserDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: InitializeUserDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<InitializeUserDoneScreenText>): InitializeUserDoneScreenText {
    return InitializeUserDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InitializeUserDoneScreenText>): InitializeUserDoneScreenText {
    const message = createBaseInitializeUserDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseInitMFAPromptScreenText(): InitMFAPromptScreenText {
  return { title: "", description: "", otpOption: "", u2fOption: "", skipButtonText: "", nextButtonText: "" };
}

export const InitMFAPromptScreenText = {
  encode(message: InitMFAPromptScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.otpOption !== "") {
      writer.uint32(26).string(message.otpOption);
    }
    if (message.u2fOption !== "") {
      writer.uint32(34).string(message.u2fOption);
    }
    if (message.skipButtonText !== "") {
      writer.uint32(42).string(message.skipButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(50).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InitMFAPromptScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInitMFAPromptScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.otpOption = reader.string();
          break;
        case 4:
          message.u2fOption = reader.string();
          break;
        case 5:
          message.skipButtonText = reader.string();
          break;
        case 6:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InitMFAPromptScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      otpOption: isSet(object.otpOption) ? String(object.otpOption) : "",
      u2fOption: isSet(object.u2fOption) ? String(object.u2fOption) : "",
      skipButtonText: isSet(object.skipButtonText) ? String(object.skipButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: InitMFAPromptScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.otpOption !== undefined && (obj.otpOption = message.otpOption);
    message.u2fOption !== undefined && (obj.u2fOption = message.u2fOption);
    message.skipButtonText !== undefined && (obj.skipButtonText = message.skipButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<InitMFAPromptScreenText>): InitMFAPromptScreenText {
    return InitMFAPromptScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InitMFAPromptScreenText>): InitMFAPromptScreenText {
    const message = createBaseInitMFAPromptScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.otpOption = object.otpOption ?? "";
    message.u2fOption = object.u2fOption ?? "";
    message.skipButtonText = object.skipButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseInitMFAOTPScreenText(): InitMFAOTPScreenText {
  return {
    title: "",
    description: "",
    descriptionOtp: "",
    secretLabel: "",
    codeLabel: "",
    nextButtonText: "",
    cancelButtonText: "",
  };
}

export const InitMFAOTPScreenText = {
  encode(message: InitMFAOTPScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.descriptionOtp !== "") {
      writer.uint32(26).string(message.descriptionOtp);
    }
    if (message.secretLabel !== "") {
      writer.uint32(34).string(message.secretLabel);
    }
    if (message.codeLabel !== "") {
      writer.uint32(42).string(message.codeLabel);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(50).string(message.nextButtonText);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(58).string(message.cancelButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InitMFAOTPScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInitMFAOTPScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.descriptionOtp = reader.string();
          break;
        case 4:
          message.secretLabel = reader.string();
          break;
        case 5:
          message.codeLabel = reader.string();
          break;
        case 6:
          message.nextButtonText = reader.string();
          break;
        case 7:
          message.cancelButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InitMFAOTPScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      descriptionOtp: isSet(object.descriptionOtp) ? String(object.descriptionOtp) : "",
      secretLabel: isSet(object.secretLabel) ? String(object.secretLabel) : "",
      codeLabel: isSet(object.codeLabel) ? String(object.codeLabel) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
    };
  },

  toJSON(message: InitMFAOTPScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.descriptionOtp !== undefined && (obj.descriptionOtp = message.descriptionOtp);
    message.secretLabel !== undefined && (obj.secretLabel = message.secretLabel);
    message.codeLabel !== undefined && (obj.codeLabel = message.codeLabel);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    return obj;
  },

  create(base?: DeepPartial<InitMFAOTPScreenText>): InitMFAOTPScreenText {
    return InitMFAOTPScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InitMFAOTPScreenText>): InitMFAOTPScreenText {
    const message = createBaseInitMFAOTPScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.descriptionOtp = object.descriptionOtp ?? "";
    message.secretLabel = object.secretLabel ?? "";
    message.codeLabel = object.codeLabel ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    return message;
  },
};

function createBaseInitMFAU2FScreenText(): InitMFAU2FScreenText {
  return {
    title: "",
    description: "",
    tokenNameLabel: "",
    notSupported: "",
    registerTokenButtonText: "",
    errorRetry: "",
  };
}

export const InitMFAU2FScreenText = {
  encode(message: InitMFAU2FScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.tokenNameLabel !== "") {
      writer.uint32(26).string(message.tokenNameLabel);
    }
    if (message.notSupported !== "") {
      writer.uint32(34).string(message.notSupported);
    }
    if (message.registerTokenButtonText !== "") {
      writer.uint32(42).string(message.registerTokenButtonText);
    }
    if (message.errorRetry !== "") {
      writer.uint32(50).string(message.errorRetry);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InitMFAU2FScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInitMFAU2FScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.tokenNameLabel = reader.string();
          break;
        case 4:
          message.notSupported = reader.string();
          break;
        case 5:
          message.registerTokenButtonText = reader.string();
          break;
        case 6:
          message.errorRetry = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InitMFAU2FScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      tokenNameLabel: isSet(object.tokenNameLabel) ? String(object.tokenNameLabel) : "",
      notSupported: isSet(object.notSupported) ? String(object.notSupported) : "",
      registerTokenButtonText: isSet(object.registerTokenButtonText) ? String(object.registerTokenButtonText) : "",
      errorRetry: isSet(object.errorRetry) ? String(object.errorRetry) : "",
    };
  },

  toJSON(message: InitMFAU2FScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.tokenNameLabel !== undefined && (obj.tokenNameLabel = message.tokenNameLabel);
    message.notSupported !== undefined && (obj.notSupported = message.notSupported);
    message.registerTokenButtonText !== undefined && (obj.registerTokenButtonText = message.registerTokenButtonText);
    message.errorRetry !== undefined && (obj.errorRetry = message.errorRetry);
    return obj;
  },

  create(base?: DeepPartial<InitMFAU2FScreenText>): InitMFAU2FScreenText {
    return InitMFAU2FScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InitMFAU2FScreenText>): InitMFAU2FScreenText {
    const message = createBaseInitMFAU2FScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.tokenNameLabel = object.tokenNameLabel ?? "";
    message.notSupported = object.notSupported ?? "";
    message.registerTokenButtonText = object.registerTokenButtonText ?? "";
    message.errorRetry = object.errorRetry ?? "";
    return message;
  },
};

function createBaseInitMFADoneScreenText(): InitMFADoneScreenText {
  return { title: "", description: "", cancelButtonText: "", nextButtonText: "" };
}

export const InitMFADoneScreenText = {
  encode(message: InitMFADoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(26).string(message.cancelButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(34).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InitMFADoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInitMFADoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.cancelButtonText = reader.string();
          break;
        case 4:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InitMFADoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: InitMFADoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<InitMFADoneScreenText>): InitMFADoneScreenText {
    return InitMFADoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<InitMFADoneScreenText>): InitMFADoneScreenText {
    const message = createBaseInitMFADoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseMFAProvidersText(): MFAProvidersText {
  return { chooseOther: "", otp: "", u2f: "" };
}

export const MFAProvidersText = {
  encode(message: MFAProvidersText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chooseOther !== "") {
      writer.uint32(10).string(message.chooseOther);
    }
    if (message.otp !== "") {
      writer.uint32(18).string(message.otp);
    }
    if (message.u2f !== "") {
      writer.uint32(26).string(message.u2f);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MFAProvidersText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMFAProvidersText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chooseOther = reader.string();
          break;
        case 2:
          message.otp = reader.string();
          break;
        case 3:
          message.u2f = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MFAProvidersText {
    return {
      chooseOther: isSet(object.chooseOther) ? String(object.chooseOther) : "",
      otp: isSet(object.otp) ? String(object.otp) : "",
      u2f: isSet(object.u2f) ? String(object.u2f) : "",
    };
  },

  toJSON(message: MFAProvidersText): unknown {
    const obj: any = {};
    message.chooseOther !== undefined && (obj.chooseOther = message.chooseOther);
    message.otp !== undefined && (obj.otp = message.otp);
    message.u2f !== undefined && (obj.u2f = message.u2f);
    return obj;
  },

  create(base?: DeepPartial<MFAProvidersText>): MFAProvidersText {
    return MFAProvidersText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<MFAProvidersText>): MFAProvidersText {
    const message = createBaseMFAProvidersText();
    message.chooseOther = object.chooseOther ?? "";
    message.otp = object.otp ?? "";
    message.u2f = object.u2f ?? "";
    return message;
  },
};

function createBaseVerifyMFAOTPScreenText(): VerifyMFAOTPScreenText {
  return { title: "", description: "", codeLabel: "", nextButtonText: "" };
}

export const VerifyMFAOTPScreenText = {
  encode(message: VerifyMFAOTPScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.codeLabel !== "") {
      writer.uint32(26).string(message.codeLabel);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(34).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMFAOTPScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMFAOTPScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.codeLabel = reader.string();
          break;
        case 4:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMFAOTPScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      codeLabel: isSet(object.codeLabel) ? String(object.codeLabel) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: VerifyMFAOTPScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.codeLabel !== undefined && (obj.codeLabel = message.codeLabel);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<VerifyMFAOTPScreenText>): VerifyMFAOTPScreenText {
    return VerifyMFAOTPScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMFAOTPScreenText>): VerifyMFAOTPScreenText {
    const message = createBaseVerifyMFAOTPScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.codeLabel = object.codeLabel ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseVerifyMFAU2FScreenText(): VerifyMFAU2FScreenText {
  return { title: "", description: "", validateTokenText: "", notSupported: "", errorRetry: "" };
}

export const VerifyMFAU2FScreenText = {
  encode(message: VerifyMFAU2FScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.validateTokenText !== "") {
      writer.uint32(26).string(message.validateTokenText);
    }
    if (message.notSupported !== "") {
      writer.uint32(34).string(message.notSupported);
    }
    if (message.errorRetry !== "") {
      writer.uint32(42).string(message.errorRetry);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyMFAU2FScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyMFAU2FScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.validateTokenText = reader.string();
          break;
        case 4:
          message.notSupported = reader.string();
          break;
        case 5:
          message.errorRetry = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyMFAU2FScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      validateTokenText: isSet(object.validateTokenText) ? String(object.validateTokenText) : "",
      notSupported: isSet(object.notSupported) ? String(object.notSupported) : "",
      errorRetry: isSet(object.errorRetry) ? String(object.errorRetry) : "",
    };
  },

  toJSON(message: VerifyMFAU2FScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.validateTokenText !== undefined && (obj.validateTokenText = message.validateTokenText);
    message.notSupported !== undefined && (obj.notSupported = message.notSupported);
    message.errorRetry !== undefined && (obj.errorRetry = message.errorRetry);
    return obj;
  },

  create(base?: DeepPartial<VerifyMFAU2FScreenText>): VerifyMFAU2FScreenText {
    return VerifyMFAU2FScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyMFAU2FScreenText>): VerifyMFAU2FScreenText {
    const message = createBaseVerifyMFAU2FScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.validateTokenText = object.validateTokenText ?? "";
    message.notSupported = object.notSupported ?? "";
    message.errorRetry = object.errorRetry ?? "";
    return message;
  },
};

function createBasePasswordlessScreenText(): PasswordlessScreenText {
  return {
    title: "",
    description: "",
    loginWithPwButtonText: "",
    validateTokenButtonText: "",
    notSupported: "",
    errorRetry: "",
  };
}

export const PasswordlessScreenText = {
  encode(message: PasswordlessScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.loginWithPwButtonText !== "") {
      writer.uint32(26).string(message.loginWithPwButtonText);
    }
    if (message.validateTokenButtonText !== "") {
      writer.uint32(34).string(message.validateTokenButtonText);
    }
    if (message.notSupported !== "") {
      writer.uint32(42).string(message.notSupported);
    }
    if (message.errorRetry !== "") {
      writer.uint32(50).string(message.errorRetry);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordlessScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordlessScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.loginWithPwButtonText = reader.string();
          break;
        case 4:
          message.validateTokenButtonText = reader.string();
          break;
        case 5:
          message.notSupported = reader.string();
          break;
        case 6:
          message.errorRetry = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordlessScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      loginWithPwButtonText: isSet(object.loginWithPwButtonText) ? String(object.loginWithPwButtonText) : "",
      validateTokenButtonText: isSet(object.validateTokenButtonText) ? String(object.validateTokenButtonText) : "",
      notSupported: isSet(object.notSupported) ? String(object.notSupported) : "",
      errorRetry: isSet(object.errorRetry) ? String(object.errorRetry) : "",
    };
  },

  toJSON(message: PasswordlessScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.loginWithPwButtonText !== undefined && (obj.loginWithPwButtonText = message.loginWithPwButtonText);
    message.validateTokenButtonText !== undefined && (obj.validateTokenButtonText = message.validateTokenButtonText);
    message.notSupported !== undefined && (obj.notSupported = message.notSupported);
    message.errorRetry !== undefined && (obj.errorRetry = message.errorRetry);
    return obj;
  },

  create(base?: DeepPartial<PasswordlessScreenText>): PasswordlessScreenText {
    return PasswordlessScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordlessScreenText>): PasswordlessScreenText {
    const message = createBasePasswordlessScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.loginWithPwButtonText = object.loginWithPwButtonText ?? "";
    message.validateTokenButtonText = object.validateTokenButtonText ?? "";
    message.notSupported = object.notSupported ?? "";
    message.errorRetry = object.errorRetry ?? "";
    return message;
  },
};

function createBasePasswordChangeScreenText(): PasswordChangeScreenText {
  return {
    title: "",
    description: "",
    oldPasswordLabel: "",
    newPasswordLabel: "",
    newPasswordConfirmLabel: "",
    cancelButtonText: "",
    nextButtonText: "",
  };
}

export const PasswordChangeScreenText = {
  encode(message: PasswordChangeScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.oldPasswordLabel !== "") {
      writer.uint32(26).string(message.oldPasswordLabel);
    }
    if (message.newPasswordLabel !== "") {
      writer.uint32(34).string(message.newPasswordLabel);
    }
    if (message.newPasswordConfirmLabel !== "") {
      writer.uint32(42).string(message.newPasswordConfirmLabel);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(50).string(message.cancelButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(58).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordChangeScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordChangeScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.oldPasswordLabel = reader.string();
          break;
        case 4:
          message.newPasswordLabel = reader.string();
          break;
        case 5:
          message.newPasswordConfirmLabel = reader.string();
          break;
        case 6:
          message.cancelButtonText = reader.string();
          break;
        case 7:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordChangeScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      oldPasswordLabel: isSet(object.oldPasswordLabel) ? String(object.oldPasswordLabel) : "",
      newPasswordLabel: isSet(object.newPasswordLabel) ? String(object.newPasswordLabel) : "",
      newPasswordConfirmLabel: isSet(object.newPasswordConfirmLabel) ? String(object.newPasswordConfirmLabel) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: PasswordChangeScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.oldPasswordLabel !== undefined && (obj.oldPasswordLabel = message.oldPasswordLabel);
    message.newPasswordLabel !== undefined && (obj.newPasswordLabel = message.newPasswordLabel);
    message.newPasswordConfirmLabel !== undefined && (obj.newPasswordConfirmLabel = message.newPasswordConfirmLabel);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<PasswordChangeScreenText>): PasswordChangeScreenText {
    return PasswordChangeScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordChangeScreenText>): PasswordChangeScreenText {
    const message = createBasePasswordChangeScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.oldPasswordLabel = object.oldPasswordLabel ?? "";
    message.newPasswordLabel = object.newPasswordLabel ?? "";
    message.newPasswordConfirmLabel = object.newPasswordConfirmLabel ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBasePasswordChangeDoneScreenText(): PasswordChangeDoneScreenText {
  return { title: "", description: "", nextButtonText: "" };
}

export const PasswordChangeDoneScreenText = {
  encode(message: PasswordChangeDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(26).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordChangeDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordChangeDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordChangeDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: PasswordChangeDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<PasswordChangeDoneScreenText>): PasswordChangeDoneScreenText {
    return PasswordChangeDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordChangeDoneScreenText>): PasswordChangeDoneScreenText {
    const message = createBasePasswordChangeDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBasePasswordResetDoneScreenText(): PasswordResetDoneScreenText {
  return { title: "", description: "", nextButtonText: "" };
}

export const PasswordResetDoneScreenText = {
  encode(message: PasswordResetDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(26).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordResetDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordResetDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordResetDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: PasswordResetDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<PasswordResetDoneScreenText>): PasswordResetDoneScreenText {
    return PasswordResetDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordResetDoneScreenText>): PasswordResetDoneScreenText {
    const message = createBasePasswordResetDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseRegistrationOptionScreenText(): RegistrationOptionScreenText {
  return { title: "", description: "", userNameButtonText: "", externalLoginDescription: "", loginButtonText: "" };
}

export const RegistrationOptionScreenText = {
  encode(message: RegistrationOptionScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.userNameButtonText !== "") {
      writer.uint32(26).string(message.userNameButtonText);
    }
    if (message.externalLoginDescription !== "") {
      writer.uint32(34).string(message.externalLoginDescription);
    }
    if (message.loginButtonText !== "") {
      writer.uint32(42).string(message.loginButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegistrationOptionScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegistrationOptionScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.userNameButtonText = reader.string();
          break;
        case 4:
          message.externalLoginDescription = reader.string();
          break;
        case 5:
          message.loginButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RegistrationOptionScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      userNameButtonText: isSet(object.userNameButtonText) ? String(object.userNameButtonText) : "",
      externalLoginDescription: isSet(object.externalLoginDescription) ? String(object.externalLoginDescription) : "",
      loginButtonText: isSet(object.loginButtonText) ? String(object.loginButtonText) : "",
    };
  },

  toJSON(message: RegistrationOptionScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.userNameButtonText !== undefined && (obj.userNameButtonText = message.userNameButtonText);
    message.externalLoginDescription !== undefined && (obj.externalLoginDescription = message.externalLoginDescription);
    message.loginButtonText !== undefined && (obj.loginButtonText = message.loginButtonText);
    return obj;
  },

  create(base?: DeepPartial<RegistrationOptionScreenText>): RegistrationOptionScreenText {
    return RegistrationOptionScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RegistrationOptionScreenText>): RegistrationOptionScreenText {
    const message = createBaseRegistrationOptionScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.userNameButtonText = object.userNameButtonText ?? "";
    message.externalLoginDescription = object.externalLoginDescription ?? "";
    message.loginButtonText = object.loginButtonText ?? "";
    return message;
  },
};

function createBaseRegistrationUserScreenText(): RegistrationUserScreenText {
  return {
    title: "",
    description: "",
    descriptionOrgRegister: "",
    firstnameLabel: "",
    lastnameLabel: "",
    emailLabel: "",
    usernameLabel: "",
    languageLabel: "",
    genderLabel: "",
    passwordLabel: "",
    passwordConfirmLabel: "",
    tosAndPrivacyLabel: "",
    tosConfirm: "",
    tosLinkText: "",
    privacyConfirm: "",
    privacyLinkText: "",
    nextButtonText: "",
    backButtonText: "",
  };
}

export const RegistrationUserScreenText = {
  encode(message: RegistrationUserScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.descriptionOrgRegister !== "") {
      writer.uint32(26).string(message.descriptionOrgRegister);
    }
    if (message.firstnameLabel !== "") {
      writer.uint32(34).string(message.firstnameLabel);
    }
    if (message.lastnameLabel !== "") {
      writer.uint32(42).string(message.lastnameLabel);
    }
    if (message.emailLabel !== "") {
      writer.uint32(50).string(message.emailLabel);
    }
    if (message.usernameLabel !== "") {
      writer.uint32(58).string(message.usernameLabel);
    }
    if (message.languageLabel !== "") {
      writer.uint32(66).string(message.languageLabel);
    }
    if (message.genderLabel !== "") {
      writer.uint32(74).string(message.genderLabel);
    }
    if (message.passwordLabel !== "") {
      writer.uint32(82).string(message.passwordLabel);
    }
    if (message.passwordConfirmLabel !== "") {
      writer.uint32(90).string(message.passwordConfirmLabel);
    }
    if (message.tosAndPrivacyLabel !== "") {
      writer.uint32(98).string(message.tosAndPrivacyLabel);
    }
    if (message.tosConfirm !== "") {
      writer.uint32(106).string(message.tosConfirm);
    }
    if (message.tosLinkText !== "") {
      writer.uint32(122).string(message.tosLinkText);
    }
    if (message.privacyConfirm !== "") {
      writer.uint32(130).string(message.privacyConfirm);
    }
    if (message.privacyLinkText !== "") {
      writer.uint32(146).string(message.privacyLinkText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(162).string(message.nextButtonText);
    }
    if (message.backButtonText !== "") {
      writer.uint32(170).string(message.backButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegistrationUserScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegistrationUserScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.descriptionOrgRegister = reader.string();
          break;
        case 4:
          message.firstnameLabel = reader.string();
          break;
        case 5:
          message.lastnameLabel = reader.string();
          break;
        case 6:
          message.emailLabel = reader.string();
          break;
        case 7:
          message.usernameLabel = reader.string();
          break;
        case 8:
          message.languageLabel = reader.string();
          break;
        case 9:
          message.genderLabel = reader.string();
          break;
        case 10:
          message.passwordLabel = reader.string();
          break;
        case 11:
          message.passwordConfirmLabel = reader.string();
          break;
        case 12:
          message.tosAndPrivacyLabel = reader.string();
          break;
        case 13:
          message.tosConfirm = reader.string();
          break;
        case 15:
          message.tosLinkText = reader.string();
          break;
        case 16:
          message.privacyConfirm = reader.string();
          break;
        case 18:
          message.privacyLinkText = reader.string();
          break;
        case 20:
          message.nextButtonText = reader.string();
          break;
        case 21:
          message.backButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RegistrationUserScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      descriptionOrgRegister: isSet(object.descriptionOrgRegister) ? String(object.descriptionOrgRegister) : "",
      firstnameLabel: isSet(object.firstnameLabel) ? String(object.firstnameLabel) : "",
      lastnameLabel: isSet(object.lastnameLabel) ? String(object.lastnameLabel) : "",
      emailLabel: isSet(object.emailLabel) ? String(object.emailLabel) : "",
      usernameLabel: isSet(object.usernameLabel) ? String(object.usernameLabel) : "",
      languageLabel: isSet(object.languageLabel) ? String(object.languageLabel) : "",
      genderLabel: isSet(object.genderLabel) ? String(object.genderLabel) : "",
      passwordLabel: isSet(object.passwordLabel) ? String(object.passwordLabel) : "",
      passwordConfirmLabel: isSet(object.passwordConfirmLabel) ? String(object.passwordConfirmLabel) : "",
      tosAndPrivacyLabel: isSet(object.tosAndPrivacyLabel) ? String(object.tosAndPrivacyLabel) : "",
      tosConfirm: isSet(object.tosConfirm) ? String(object.tosConfirm) : "",
      tosLinkText: isSet(object.tosLinkText) ? String(object.tosLinkText) : "",
      privacyConfirm: isSet(object.privacyConfirm) ? String(object.privacyConfirm) : "",
      privacyLinkText: isSet(object.privacyLinkText) ? String(object.privacyLinkText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      backButtonText: isSet(object.backButtonText) ? String(object.backButtonText) : "",
    };
  },

  toJSON(message: RegistrationUserScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.descriptionOrgRegister !== undefined && (obj.descriptionOrgRegister = message.descriptionOrgRegister);
    message.firstnameLabel !== undefined && (obj.firstnameLabel = message.firstnameLabel);
    message.lastnameLabel !== undefined && (obj.lastnameLabel = message.lastnameLabel);
    message.emailLabel !== undefined && (obj.emailLabel = message.emailLabel);
    message.usernameLabel !== undefined && (obj.usernameLabel = message.usernameLabel);
    message.languageLabel !== undefined && (obj.languageLabel = message.languageLabel);
    message.genderLabel !== undefined && (obj.genderLabel = message.genderLabel);
    message.passwordLabel !== undefined && (obj.passwordLabel = message.passwordLabel);
    message.passwordConfirmLabel !== undefined && (obj.passwordConfirmLabel = message.passwordConfirmLabel);
    message.tosAndPrivacyLabel !== undefined && (obj.tosAndPrivacyLabel = message.tosAndPrivacyLabel);
    message.tosConfirm !== undefined && (obj.tosConfirm = message.tosConfirm);
    message.tosLinkText !== undefined && (obj.tosLinkText = message.tosLinkText);
    message.privacyConfirm !== undefined && (obj.privacyConfirm = message.privacyConfirm);
    message.privacyLinkText !== undefined && (obj.privacyLinkText = message.privacyLinkText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.backButtonText !== undefined && (obj.backButtonText = message.backButtonText);
    return obj;
  },

  create(base?: DeepPartial<RegistrationUserScreenText>): RegistrationUserScreenText {
    return RegistrationUserScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RegistrationUserScreenText>): RegistrationUserScreenText {
    const message = createBaseRegistrationUserScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.descriptionOrgRegister = object.descriptionOrgRegister ?? "";
    message.firstnameLabel = object.firstnameLabel ?? "";
    message.lastnameLabel = object.lastnameLabel ?? "";
    message.emailLabel = object.emailLabel ?? "";
    message.usernameLabel = object.usernameLabel ?? "";
    message.languageLabel = object.languageLabel ?? "";
    message.genderLabel = object.genderLabel ?? "";
    message.passwordLabel = object.passwordLabel ?? "";
    message.passwordConfirmLabel = object.passwordConfirmLabel ?? "";
    message.tosAndPrivacyLabel = object.tosAndPrivacyLabel ?? "";
    message.tosConfirm = object.tosConfirm ?? "";
    message.tosLinkText = object.tosLinkText ?? "";
    message.privacyConfirm = object.privacyConfirm ?? "";
    message.privacyLinkText = object.privacyLinkText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.backButtonText = object.backButtonText ?? "";
    return message;
  },
};

function createBaseExternalRegistrationUserOverviewScreenText(): ExternalRegistrationUserOverviewScreenText {
  return {
    title: "",
    description: "",
    emailLabel: "",
    usernameLabel: "",
    firstnameLabel: "",
    lastnameLabel: "",
    nicknameLabel: "",
    languageLabel: "",
    phoneLabel: "",
    tosAndPrivacyLabel: "",
    tosConfirm: "",
    tosLinkText: "",
    privacyLinkText: "",
    backButtonText: "",
    nextButtonText: "",
    privacyConfirm: "",
  };
}

export const ExternalRegistrationUserOverviewScreenText = {
  encode(message: ExternalRegistrationUserOverviewScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.emailLabel !== "") {
      writer.uint32(26).string(message.emailLabel);
    }
    if (message.usernameLabel !== "") {
      writer.uint32(34).string(message.usernameLabel);
    }
    if (message.firstnameLabel !== "") {
      writer.uint32(42).string(message.firstnameLabel);
    }
    if (message.lastnameLabel !== "") {
      writer.uint32(50).string(message.lastnameLabel);
    }
    if (message.nicknameLabel !== "") {
      writer.uint32(58).string(message.nicknameLabel);
    }
    if (message.languageLabel !== "") {
      writer.uint32(66).string(message.languageLabel);
    }
    if (message.phoneLabel !== "") {
      writer.uint32(74).string(message.phoneLabel);
    }
    if (message.tosAndPrivacyLabel !== "") {
      writer.uint32(82).string(message.tosAndPrivacyLabel);
    }
    if (message.tosConfirm !== "") {
      writer.uint32(90).string(message.tosConfirm);
    }
    if (message.tosLinkText !== "") {
      writer.uint32(98).string(message.tosLinkText);
    }
    if (message.privacyLinkText !== "") {
      writer.uint32(114).string(message.privacyLinkText);
    }
    if (message.backButtonText !== "") {
      writer.uint32(122).string(message.backButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(130).string(message.nextButtonText);
    }
    if (message.privacyConfirm !== "") {
      writer.uint32(138).string(message.privacyConfirm);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExternalRegistrationUserOverviewScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExternalRegistrationUserOverviewScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.emailLabel = reader.string();
          break;
        case 4:
          message.usernameLabel = reader.string();
          break;
        case 5:
          message.firstnameLabel = reader.string();
          break;
        case 6:
          message.lastnameLabel = reader.string();
          break;
        case 7:
          message.nicknameLabel = reader.string();
          break;
        case 8:
          message.languageLabel = reader.string();
          break;
        case 9:
          message.phoneLabel = reader.string();
          break;
        case 10:
          message.tosAndPrivacyLabel = reader.string();
          break;
        case 11:
          message.tosConfirm = reader.string();
          break;
        case 12:
          message.tosLinkText = reader.string();
          break;
        case 14:
          message.privacyLinkText = reader.string();
          break;
        case 15:
          message.backButtonText = reader.string();
          break;
        case 16:
          message.nextButtonText = reader.string();
          break;
        case 17:
          message.privacyConfirm = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExternalRegistrationUserOverviewScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      emailLabel: isSet(object.emailLabel) ? String(object.emailLabel) : "",
      usernameLabel: isSet(object.usernameLabel) ? String(object.usernameLabel) : "",
      firstnameLabel: isSet(object.firstnameLabel) ? String(object.firstnameLabel) : "",
      lastnameLabel: isSet(object.lastnameLabel) ? String(object.lastnameLabel) : "",
      nicknameLabel: isSet(object.nicknameLabel) ? String(object.nicknameLabel) : "",
      languageLabel: isSet(object.languageLabel) ? String(object.languageLabel) : "",
      phoneLabel: isSet(object.phoneLabel) ? String(object.phoneLabel) : "",
      tosAndPrivacyLabel: isSet(object.tosAndPrivacyLabel) ? String(object.tosAndPrivacyLabel) : "",
      tosConfirm: isSet(object.tosConfirm) ? String(object.tosConfirm) : "",
      tosLinkText: isSet(object.tosLinkText) ? String(object.tosLinkText) : "",
      privacyLinkText: isSet(object.privacyLinkText) ? String(object.privacyLinkText) : "",
      backButtonText: isSet(object.backButtonText) ? String(object.backButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      privacyConfirm: isSet(object.privacyConfirm) ? String(object.privacyConfirm) : "",
    };
  },

  toJSON(message: ExternalRegistrationUserOverviewScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.emailLabel !== undefined && (obj.emailLabel = message.emailLabel);
    message.usernameLabel !== undefined && (obj.usernameLabel = message.usernameLabel);
    message.firstnameLabel !== undefined && (obj.firstnameLabel = message.firstnameLabel);
    message.lastnameLabel !== undefined && (obj.lastnameLabel = message.lastnameLabel);
    message.nicknameLabel !== undefined && (obj.nicknameLabel = message.nicknameLabel);
    message.languageLabel !== undefined && (obj.languageLabel = message.languageLabel);
    message.phoneLabel !== undefined && (obj.phoneLabel = message.phoneLabel);
    message.tosAndPrivacyLabel !== undefined && (obj.tosAndPrivacyLabel = message.tosAndPrivacyLabel);
    message.tosConfirm !== undefined && (obj.tosConfirm = message.tosConfirm);
    message.tosLinkText !== undefined && (obj.tosLinkText = message.tosLinkText);
    message.privacyLinkText !== undefined && (obj.privacyLinkText = message.privacyLinkText);
    message.backButtonText !== undefined && (obj.backButtonText = message.backButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.privacyConfirm !== undefined && (obj.privacyConfirm = message.privacyConfirm);
    return obj;
  },

  create(base?: DeepPartial<ExternalRegistrationUserOverviewScreenText>): ExternalRegistrationUserOverviewScreenText {
    return ExternalRegistrationUserOverviewScreenText.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<ExternalRegistrationUserOverviewScreenText>,
  ): ExternalRegistrationUserOverviewScreenText {
    const message = createBaseExternalRegistrationUserOverviewScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.emailLabel = object.emailLabel ?? "";
    message.usernameLabel = object.usernameLabel ?? "";
    message.firstnameLabel = object.firstnameLabel ?? "";
    message.lastnameLabel = object.lastnameLabel ?? "";
    message.nicknameLabel = object.nicknameLabel ?? "";
    message.languageLabel = object.languageLabel ?? "";
    message.phoneLabel = object.phoneLabel ?? "";
    message.tosAndPrivacyLabel = object.tosAndPrivacyLabel ?? "";
    message.tosConfirm = object.tosConfirm ?? "";
    message.tosLinkText = object.tosLinkText ?? "";
    message.privacyLinkText = object.privacyLinkText ?? "";
    message.backButtonText = object.backButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.privacyConfirm = object.privacyConfirm ?? "";
    return message;
  },
};

function createBaseRegistrationOrgScreenText(): RegistrationOrgScreenText {
  return {
    title: "",
    description: "",
    orgnameLabel: "",
    firstnameLabel: "",
    lastnameLabel: "",
    usernameLabel: "",
    emailLabel: "",
    passwordLabel: "",
    passwordConfirmLabel: "",
    tosAndPrivacyLabel: "",
    tosConfirm: "",
    tosLinkText: "",
    privacyConfirm: "",
    privacyLinkText: "",
    saveButtonText: "",
  };
}

export const RegistrationOrgScreenText = {
  encode(message: RegistrationOrgScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.orgnameLabel !== "") {
      writer.uint32(26).string(message.orgnameLabel);
    }
    if (message.firstnameLabel !== "") {
      writer.uint32(34).string(message.firstnameLabel);
    }
    if (message.lastnameLabel !== "") {
      writer.uint32(42).string(message.lastnameLabel);
    }
    if (message.usernameLabel !== "") {
      writer.uint32(50).string(message.usernameLabel);
    }
    if (message.emailLabel !== "") {
      writer.uint32(58).string(message.emailLabel);
    }
    if (message.passwordLabel !== "") {
      writer.uint32(74).string(message.passwordLabel);
    }
    if (message.passwordConfirmLabel !== "") {
      writer.uint32(82).string(message.passwordConfirmLabel);
    }
    if (message.tosAndPrivacyLabel !== "") {
      writer.uint32(90).string(message.tosAndPrivacyLabel);
    }
    if (message.tosConfirm !== "") {
      writer.uint32(98).string(message.tosConfirm);
    }
    if (message.tosLinkText !== "") {
      writer.uint32(114).string(message.tosLinkText);
    }
    if (message.privacyConfirm !== "") {
      writer.uint32(122).string(message.privacyConfirm);
    }
    if (message.privacyLinkText !== "") {
      writer.uint32(138).string(message.privacyLinkText);
    }
    if (message.saveButtonText !== "") {
      writer.uint32(154).string(message.saveButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegistrationOrgScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegistrationOrgScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.orgnameLabel = reader.string();
          break;
        case 4:
          message.firstnameLabel = reader.string();
          break;
        case 5:
          message.lastnameLabel = reader.string();
          break;
        case 6:
          message.usernameLabel = reader.string();
          break;
        case 7:
          message.emailLabel = reader.string();
          break;
        case 9:
          message.passwordLabel = reader.string();
          break;
        case 10:
          message.passwordConfirmLabel = reader.string();
          break;
        case 11:
          message.tosAndPrivacyLabel = reader.string();
          break;
        case 12:
          message.tosConfirm = reader.string();
          break;
        case 14:
          message.tosLinkText = reader.string();
          break;
        case 15:
          message.privacyConfirm = reader.string();
          break;
        case 17:
          message.privacyLinkText = reader.string();
          break;
        case 19:
          message.saveButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RegistrationOrgScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      orgnameLabel: isSet(object.orgnameLabel) ? String(object.orgnameLabel) : "",
      firstnameLabel: isSet(object.firstnameLabel) ? String(object.firstnameLabel) : "",
      lastnameLabel: isSet(object.lastnameLabel) ? String(object.lastnameLabel) : "",
      usernameLabel: isSet(object.usernameLabel) ? String(object.usernameLabel) : "",
      emailLabel: isSet(object.emailLabel) ? String(object.emailLabel) : "",
      passwordLabel: isSet(object.passwordLabel) ? String(object.passwordLabel) : "",
      passwordConfirmLabel: isSet(object.passwordConfirmLabel) ? String(object.passwordConfirmLabel) : "",
      tosAndPrivacyLabel: isSet(object.tosAndPrivacyLabel) ? String(object.tosAndPrivacyLabel) : "",
      tosConfirm: isSet(object.tosConfirm) ? String(object.tosConfirm) : "",
      tosLinkText: isSet(object.tosLinkText) ? String(object.tosLinkText) : "",
      privacyConfirm: isSet(object.privacyConfirm) ? String(object.privacyConfirm) : "",
      privacyLinkText: isSet(object.privacyLinkText) ? String(object.privacyLinkText) : "",
      saveButtonText: isSet(object.saveButtonText) ? String(object.saveButtonText) : "",
    };
  },

  toJSON(message: RegistrationOrgScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.orgnameLabel !== undefined && (obj.orgnameLabel = message.orgnameLabel);
    message.firstnameLabel !== undefined && (obj.firstnameLabel = message.firstnameLabel);
    message.lastnameLabel !== undefined && (obj.lastnameLabel = message.lastnameLabel);
    message.usernameLabel !== undefined && (obj.usernameLabel = message.usernameLabel);
    message.emailLabel !== undefined && (obj.emailLabel = message.emailLabel);
    message.passwordLabel !== undefined && (obj.passwordLabel = message.passwordLabel);
    message.passwordConfirmLabel !== undefined && (obj.passwordConfirmLabel = message.passwordConfirmLabel);
    message.tosAndPrivacyLabel !== undefined && (obj.tosAndPrivacyLabel = message.tosAndPrivacyLabel);
    message.tosConfirm !== undefined && (obj.tosConfirm = message.tosConfirm);
    message.tosLinkText !== undefined && (obj.tosLinkText = message.tosLinkText);
    message.privacyConfirm !== undefined && (obj.privacyConfirm = message.privacyConfirm);
    message.privacyLinkText !== undefined && (obj.privacyLinkText = message.privacyLinkText);
    message.saveButtonText !== undefined && (obj.saveButtonText = message.saveButtonText);
    return obj;
  },

  create(base?: DeepPartial<RegistrationOrgScreenText>): RegistrationOrgScreenText {
    return RegistrationOrgScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RegistrationOrgScreenText>): RegistrationOrgScreenText {
    const message = createBaseRegistrationOrgScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.orgnameLabel = object.orgnameLabel ?? "";
    message.firstnameLabel = object.firstnameLabel ?? "";
    message.lastnameLabel = object.lastnameLabel ?? "";
    message.usernameLabel = object.usernameLabel ?? "";
    message.emailLabel = object.emailLabel ?? "";
    message.passwordLabel = object.passwordLabel ?? "";
    message.passwordConfirmLabel = object.passwordConfirmLabel ?? "";
    message.tosAndPrivacyLabel = object.tosAndPrivacyLabel ?? "";
    message.tosConfirm = object.tosConfirm ?? "";
    message.tosLinkText = object.tosLinkText ?? "";
    message.privacyConfirm = object.privacyConfirm ?? "";
    message.privacyLinkText = object.privacyLinkText ?? "";
    message.saveButtonText = object.saveButtonText ?? "";
    return message;
  },
};

function createBaseLinkingUserDoneScreenText(): LinkingUserDoneScreenText {
  return { title: "", description: "", cancelButtonText: "", nextButtonText: "" };
}

export const LinkingUserDoneScreenText = {
  encode(message: LinkingUserDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(26).string(message.cancelButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(34).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LinkingUserDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLinkingUserDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.cancelButtonText = reader.string();
          break;
        case 4:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LinkingUserDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: LinkingUserDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<LinkingUserDoneScreenText>): LinkingUserDoneScreenText {
    return LinkingUserDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LinkingUserDoneScreenText>): LinkingUserDoneScreenText {
    const message = createBaseLinkingUserDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseExternalUserNotFoundScreenText(): ExternalUserNotFoundScreenText {
  return {
    title: "",
    description: "",
    linkButtonText: "",
    autoRegisterButtonText: "",
    tosAndPrivacyLabel: "",
    tosConfirm: "",
    tosLinkText: "",
    privacyLinkText: "",
    privacyConfirm: "",
  };
}

export const ExternalUserNotFoundScreenText = {
  encode(message: ExternalUserNotFoundScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.linkButtonText !== "") {
      writer.uint32(26).string(message.linkButtonText);
    }
    if (message.autoRegisterButtonText !== "") {
      writer.uint32(34).string(message.autoRegisterButtonText);
    }
    if (message.tosAndPrivacyLabel !== "") {
      writer.uint32(42).string(message.tosAndPrivacyLabel);
    }
    if (message.tosConfirm !== "") {
      writer.uint32(50).string(message.tosConfirm);
    }
    if (message.tosLinkText !== "") {
      writer.uint32(58).string(message.tosLinkText);
    }
    if (message.privacyLinkText !== "") {
      writer.uint32(66).string(message.privacyLinkText);
    }
    if (message.privacyConfirm !== "") {
      writer.uint32(82).string(message.privacyConfirm);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExternalUserNotFoundScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExternalUserNotFoundScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.linkButtonText = reader.string();
          break;
        case 4:
          message.autoRegisterButtonText = reader.string();
          break;
        case 5:
          message.tosAndPrivacyLabel = reader.string();
          break;
        case 6:
          message.tosConfirm = reader.string();
          break;
        case 7:
          message.tosLinkText = reader.string();
          break;
        case 8:
          message.privacyLinkText = reader.string();
          break;
        case 10:
          message.privacyConfirm = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExternalUserNotFoundScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      linkButtonText: isSet(object.linkButtonText) ? String(object.linkButtonText) : "",
      autoRegisterButtonText: isSet(object.autoRegisterButtonText) ? String(object.autoRegisterButtonText) : "",
      tosAndPrivacyLabel: isSet(object.tosAndPrivacyLabel) ? String(object.tosAndPrivacyLabel) : "",
      tosConfirm: isSet(object.tosConfirm) ? String(object.tosConfirm) : "",
      tosLinkText: isSet(object.tosLinkText) ? String(object.tosLinkText) : "",
      privacyLinkText: isSet(object.privacyLinkText) ? String(object.privacyLinkText) : "",
      privacyConfirm: isSet(object.privacyConfirm) ? String(object.privacyConfirm) : "",
    };
  },

  toJSON(message: ExternalUserNotFoundScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.linkButtonText !== undefined && (obj.linkButtonText = message.linkButtonText);
    message.autoRegisterButtonText !== undefined && (obj.autoRegisterButtonText = message.autoRegisterButtonText);
    message.tosAndPrivacyLabel !== undefined && (obj.tosAndPrivacyLabel = message.tosAndPrivacyLabel);
    message.tosConfirm !== undefined && (obj.tosConfirm = message.tosConfirm);
    message.tosLinkText !== undefined && (obj.tosLinkText = message.tosLinkText);
    message.privacyLinkText !== undefined && (obj.privacyLinkText = message.privacyLinkText);
    message.privacyConfirm !== undefined && (obj.privacyConfirm = message.privacyConfirm);
    return obj;
  },

  create(base?: DeepPartial<ExternalUserNotFoundScreenText>): ExternalUserNotFoundScreenText {
    return ExternalUserNotFoundScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExternalUserNotFoundScreenText>): ExternalUserNotFoundScreenText {
    const message = createBaseExternalUserNotFoundScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.linkButtonText = object.linkButtonText ?? "";
    message.autoRegisterButtonText = object.autoRegisterButtonText ?? "";
    message.tosAndPrivacyLabel = object.tosAndPrivacyLabel ?? "";
    message.tosConfirm = object.tosConfirm ?? "";
    message.tosLinkText = object.tosLinkText ?? "";
    message.privacyLinkText = object.privacyLinkText ?? "";
    message.privacyConfirm = object.privacyConfirm ?? "";
    return message;
  },
};

function createBaseSuccessLoginScreenText(): SuccessLoginScreenText {
  return { title: "", autoRedirectDescription: "", redirectedDescription: "", nextButtonText: "" };
}

export const SuccessLoginScreenText = {
  encode(message: SuccessLoginScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.autoRedirectDescription !== "") {
      writer.uint32(18).string(message.autoRedirectDescription);
    }
    if (message.redirectedDescription !== "") {
      writer.uint32(26).string(message.redirectedDescription);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(34).string(message.nextButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SuccessLoginScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSuccessLoginScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.autoRedirectDescription = reader.string();
          break;
        case 3:
          message.redirectedDescription = reader.string();
          break;
        case 4:
          message.nextButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SuccessLoginScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      autoRedirectDescription: isSet(object.autoRedirectDescription) ? String(object.autoRedirectDescription) : "",
      redirectedDescription: isSet(object.redirectedDescription) ? String(object.redirectedDescription) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
    };
  },

  toJSON(message: SuccessLoginScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.autoRedirectDescription !== undefined && (obj.autoRedirectDescription = message.autoRedirectDescription);
    message.redirectedDescription !== undefined && (obj.redirectedDescription = message.redirectedDescription);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    return obj;
  },

  create(base?: DeepPartial<SuccessLoginScreenText>): SuccessLoginScreenText {
    return SuccessLoginScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SuccessLoginScreenText>): SuccessLoginScreenText {
    const message = createBaseSuccessLoginScreenText();
    message.title = object.title ?? "";
    message.autoRedirectDescription = object.autoRedirectDescription ?? "";
    message.redirectedDescription = object.redirectedDescription ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    return message;
  },
};

function createBaseLogoutDoneScreenText(): LogoutDoneScreenText {
  return { title: "", description: "", loginButtonText: "" };
}

export const LogoutDoneScreenText = {
  encode(message: LogoutDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.loginButtonText !== "") {
      writer.uint32(26).string(message.loginButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LogoutDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLogoutDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.loginButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LogoutDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      loginButtonText: isSet(object.loginButtonText) ? String(object.loginButtonText) : "",
    };
  },

  toJSON(message: LogoutDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.loginButtonText !== undefined && (obj.loginButtonText = message.loginButtonText);
    return obj;
  },

  create(base?: DeepPartial<LogoutDoneScreenText>): LogoutDoneScreenText {
    return LogoutDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LogoutDoneScreenText>): LogoutDoneScreenText {
    const message = createBaseLogoutDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.loginButtonText = object.loginButtonText ?? "";
    return message;
  },
};

function createBaseFooterText(): FooterText {
  return { tos: "", privacyPolicy: "", help: "" };
}

export const FooterText = {
  encode(message: FooterText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tos !== "") {
      writer.uint32(10).string(message.tos);
    }
    if (message.privacyPolicy !== "") {
      writer.uint32(26).string(message.privacyPolicy);
    }
    if (message.help !== "") {
      writer.uint32(42).string(message.help);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FooterText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFooterText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tos = reader.string();
          break;
        case 3:
          message.privacyPolicy = reader.string();
          break;
        case 5:
          message.help = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FooterText {
    return {
      tos: isSet(object.tos) ? String(object.tos) : "",
      privacyPolicy: isSet(object.privacyPolicy) ? String(object.privacyPolicy) : "",
      help: isSet(object.help) ? String(object.help) : "",
    };
  },

  toJSON(message: FooterText): unknown {
    const obj: any = {};
    message.tos !== undefined && (obj.tos = message.tos);
    message.privacyPolicy !== undefined && (obj.privacyPolicy = message.privacyPolicy);
    message.help !== undefined && (obj.help = message.help);
    return obj;
  },

  create(base?: DeepPartial<FooterText>): FooterText {
    return FooterText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<FooterText>): FooterText {
    const message = createBaseFooterText();
    message.tos = object.tos ?? "";
    message.privacyPolicy = object.privacyPolicy ?? "";
    message.help = object.help ?? "";
    return message;
  },
};

function createBasePasswordlessPromptScreenText(): PasswordlessPromptScreenText {
  return {
    title: "",
    description: "",
    descriptionInit: "",
    passwordlessButtonText: "",
    nextButtonText: "",
    skipButtonText: "",
  };
}

export const PasswordlessPromptScreenText = {
  encode(message: PasswordlessPromptScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.descriptionInit !== "") {
      writer.uint32(26).string(message.descriptionInit);
    }
    if (message.passwordlessButtonText !== "") {
      writer.uint32(34).string(message.passwordlessButtonText);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(42).string(message.nextButtonText);
    }
    if (message.skipButtonText !== "") {
      writer.uint32(50).string(message.skipButtonText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordlessPromptScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordlessPromptScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.descriptionInit = reader.string();
          break;
        case 4:
          message.passwordlessButtonText = reader.string();
          break;
        case 5:
          message.nextButtonText = reader.string();
          break;
        case 6:
          message.skipButtonText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordlessPromptScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      descriptionInit: isSet(object.descriptionInit) ? String(object.descriptionInit) : "",
      passwordlessButtonText: isSet(object.passwordlessButtonText) ? String(object.passwordlessButtonText) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      skipButtonText: isSet(object.skipButtonText) ? String(object.skipButtonText) : "",
    };
  },

  toJSON(message: PasswordlessPromptScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.descriptionInit !== undefined && (obj.descriptionInit = message.descriptionInit);
    message.passwordlessButtonText !== undefined && (obj.passwordlessButtonText = message.passwordlessButtonText);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.skipButtonText !== undefined && (obj.skipButtonText = message.skipButtonText);
    return obj;
  },

  create(base?: DeepPartial<PasswordlessPromptScreenText>): PasswordlessPromptScreenText {
    return PasswordlessPromptScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordlessPromptScreenText>): PasswordlessPromptScreenText {
    const message = createBasePasswordlessPromptScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.descriptionInit = object.descriptionInit ?? "";
    message.passwordlessButtonText = object.passwordlessButtonText ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.skipButtonText = object.skipButtonText ?? "";
    return message;
  },
};

function createBasePasswordlessRegistrationScreenText(): PasswordlessRegistrationScreenText {
  return {
    title: "",
    description: "",
    tokenNameLabel: "",
    notSupported: "",
    registerTokenButtonText: "",
    errorRetry: "",
  };
}

export const PasswordlessRegistrationScreenText = {
  encode(message: PasswordlessRegistrationScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.tokenNameLabel !== "") {
      writer.uint32(26).string(message.tokenNameLabel);
    }
    if (message.notSupported !== "") {
      writer.uint32(34).string(message.notSupported);
    }
    if (message.registerTokenButtonText !== "") {
      writer.uint32(42).string(message.registerTokenButtonText);
    }
    if (message.errorRetry !== "") {
      writer.uint32(50).string(message.errorRetry);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordlessRegistrationScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordlessRegistrationScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.tokenNameLabel = reader.string();
          break;
        case 4:
          message.notSupported = reader.string();
          break;
        case 5:
          message.registerTokenButtonText = reader.string();
          break;
        case 6:
          message.errorRetry = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordlessRegistrationScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      tokenNameLabel: isSet(object.tokenNameLabel) ? String(object.tokenNameLabel) : "",
      notSupported: isSet(object.notSupported) ? String(object.notSupported) : "",
      registerTokenButtonText: isSet(object.registerTokenButtonText) ? String(object.registerTokenButtonText) : "",
      errorRetry: isSet(object.errorRetry) ? String(object.errorRetry) : "",
    };
  },

  toJSON(message: PasswordlessRegistrationScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.tokenNameLabel !== undefined && (obj.tokenNameLabel = message.tokenNameLabel);
    message.notSupported !== undefined && (obj.notSupported = message.notSupported);
    message.registerTokenButtonText !== undefined && (obj.registerTokenButtonText = message.registerTokenButtonText);
    message.errorRetry !== undefined && (obj.errorRetry = message.errorRetry);
    return obj;
  },

  create(base?: DeepPartial<PasswordlessRegistrationScreenText>): PasswordlessRegistrationScreenText {
    return PasswordlessRegistrationScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordlessRegistrationScreenText>): PasswordlessRegistrationScreenText {
    const message = createBasePasswordlessRegistrationScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.tokenNameLabel = object.tokenNameLabel ?? "";
    message.notSupported = object.notSupported ?? "";
    message.registerTokenButtonText = object.registerTokenButtonText ?? "";
    message.errorRetry = object.errorRetry ?? "";
    return message;
  },
};

function createBasePasswordlessRegistrationDoneScreenText(): PasswordlessRegistrationDoneScreenText {
  return { title: "", description: "", nextButtonText: "", cancelButtonText: "", descriptionClose: "" };
}

export const PasswordlessRegistrationDoneScreenText = {
  encode(message: PasswordlessRegistrationDoneScreenText, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.nextButtonText !== "") {
      writer.uint32(26).string(message.nextButtonText);
    }
    if (message.cancelButtonText !== "") {
      writer.uint32(34).string(message.cancelButtonText);
    }
    if (message.descriptionClose !== "") {
      writer.uint32(42).string(message.descriptionClose);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PasswordlessRegistrationDoneScreenText {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePasswordlessRegistrationDoneScreenText();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.nextButtonText = reader.string();
          break;
        case 4:
          message.cancelButtonText = reader.string();
          break;
        case 5:
          message.descriptionClose = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PasswordlessRegistrationDoneScreenText {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      nextButtonText: isSet(object.nextButtonText) ? String(object.nextButtonText) : "",
      cancelButtonText: isSet(object.cancelButtonText) ? String(object.cancelButtonText) : "",
      descriptionClose: isSet(object.descriptionClose) ? String(object.descriptionClose) : "",
    };
  },

  toJSON(message: PasswordlessRegistrationDoneScreenText): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.nextButtonText !== undefined && (obj.nextButtonText = message.nextButtonText);
    message.cancelButtonText !== undefined && (obj.cancelButtonText = message.cancelButtonText);
    message.descriptionClose !== undefined && (obj.descriptionClose = message.descriptionClose);
    return obj;
  },

  create(base?: DeepPartial<PasswordlessRegistrationDoneScreenText>): PasswordlessRegistrationDoneScreenText {
    return PasswordlessRegistrationDoneScreenText.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<PasswordlessRegistrationDoneScreenText>): PasswordlessRegistrationDoneScreenText {
    const message = createBasePasswordlessRegistrationDoneScreenText();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.nextButtonText = object.nextButtonText ?? "";
    message.cancelButtonText = object.cancelButtonText ?? "";
    message.descriptionClose = object.descriptionClose ?? "";
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

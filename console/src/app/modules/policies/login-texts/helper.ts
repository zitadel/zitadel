import { SetCustomLoginTextsRequest as AdminSetCustomLoginTextsRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { SetCustomLoginTextsRequest } from 'src/app/proto/generated/zitadel/management_pb';
import {
  EmailVerificationDoneScreenText,
  EmailVerificationScreenText,
  ExternalUserNotFoundScreenText,
  FooterText,
  InitializeUserDoneScreenText,
  InitializeUserScreenText,
  InitMFADoneScreenText,
  InitMFAOTPScreenText,
  InitMFAPromptScreenText,
  InitMFAU2FScreenText,
  InitPasswordDoneScreenText,
  InitPasswordScreenText,
  LinkingUserDoneScreenText,
  LoginScreenText,
  LogoutDoneScreenText,
  MFAProvidersText,
  PasswordChangeDoneScreenText,
  PasswordChangeScreenText,
  PasswordlessScreenText,
  PasswordResetDoneScreenText,
  PasswordScreenText,
  RegistrationOptionScreenText,
  RegistrationOrgScreenText,
  RegistrationUserScreenText,
  SelectAccountScreenText,
  SuccessLoginScreenText,
  UsernameChangeDoneScreenText,
  UsernameChangeScreenText,
  VerifyMFAOTPScreenText,
  VerifyMFAU2FScreenText,
} from 'src/app/proto/generated/zitadel/text_pb';

type Req = AdminSetCustomLoginTextsRequest | SetCustomLoginTextsRequest;
type Map = AdminSetCustomLoginTextsRequest.AsObject | SetCustomLoginTextsRequest.AsObject;

export function mapRequestValues(map: Partial<Map>, req: Req): Req {
  if (!!map.emailVerificationDoneText) {
    const r = new EmailVerificationDoneScreenText();
    r.setCancelButtonText(map.emailVerificationDoneText?.cancelButtonText ?? '');
    r.setDescription(map.emailVerificationDoneText?.description ?? '');
    r.setLoginButtonText(map.emailVerificationDoneText?.loginButtonText ?? '');
    r.setNextButtonText(map.emailVerificationDoneText?.nextButtonText ?? '');
    r.setTitle(map.emailVerificationDoneText?.title ?? '');

    req.setEmailVerificationDoneText(r);
  }


  if (!!map.emailVerificationText) {
    const r = new EmailVerificationScreenText();
    r.setCodeLabel(map.emailVerificationText?.codeLabel ?? '');
    r.setDescription(map.emailVerificationText?.description ?? '');
    r.setNextButtonText(map.emailVerificationText?.nextButtonText ?? '');
    r.setResendButtonText(map.emailVerificationText?.resendButtonText ?? '');
    r.setTitle(map.emailVerificationText?.title ?? '');

    req.setEmailVerificationText(r);
  }


  if (!!map.externalUserNotFoundText) {
    const r = new ExternalUserNotFoundScreenText();
    r.setAutoRegisterButtonText(map.externalUserNotFoundText?.autoRegisterButtonText ?? '');
    r.setDescription(map.externalUserNotFoundText?.description ?? '');
    r.setLinkButtonText(map.externalUserNotFoundText?.linkButtonText ?? '');
    r.setTitle(map.externalUserNotFoundText?.title ?? '');

    req.setExternalUserNotFoundText(r);
  }

  if (!!map.footerText) {
    const r = new FooterText();
    r.setHelp(map.footerText?.help ?? '');
    r.setHelpLink(map.footerText?.helpLink ?? '');
    r.setPrivacyPolicy(map.footerText?.privacyPolicy ?? '');
    r.setPrivacyPolicyLink(map.footerText?.privacyPolicyLink ?? '');
    r.setTos(map.footerText?.tos ?? '');
    r.setTosLink(map.footerText?.tosLink ?? '');

    req.setFooterText(r);
  }

  if (!!map.initMfaDoneText) {
    const r = new InitMFADoneScreenText();
    r.setCancelButtonText(map.initMfaDoneText?.cancelButtonText ?? '');
    r.setDescription(map.initMfaDoneText?.description ?? '');
    r.setNextButtonText(map.initMfaDoneText?.nextButtonText ?? '');
    r.setTitle(map.initMfaDoneText?.title ?? '');

    req.setInitMfaDoneText(r);
  }

  if (!!map.initMfaOtpText) {
    const r = new InitMFAOTPScreenText();
    r.setCancelButtonText(map.initMfaOtpText?.cancelButtonText ?? '');
    r.setCodeLabel(map.initMfaOtpText?.codeLabel ?? '');
    r.setDescription(map.initMfaOtpText?.description ?? '');
    r.setDescriptionOtp(map.initMfaOtpText?.descriptionOtp ?? '');
    r.setNextButtonText(map.initMfaOtpText?.nextButtonText ?? '');
    r.setSecretLabel(map.initMfaOtpText?.secretLabel ?? '');
    r.setTitle(map.initMfaOtpText?.title ?? '');

    req.setInitMfaOtpText(r);
  }

  if (!!map.initMfaPromptText) {
    const r = new InitMFAPromptScreenText();
    r.setDescription(map.initMfaPromptText?.description ?? '');
    r.setNextButtonText(map.initMfaPromptText?.nextButtonText ?? '');
    r.setOtpOption(map.initMfaPromptText?.otpOption ?? '');
    r.setSkipButtonText(map.initMfaPromptText?.skipButtonText ?? '');
    r.setTitle(map.initMfaPromptText?.title ?? '');
    r.setU2fOption(map.initMfaPromptText?.otpOption ?? '');

    req.setInitMfaPromptText(r);
  }


  if (!!map.initMfaU2fText) {
    const r = new InitMFAU2FScreenText();
    r.setDescription(map.initMfaU2fText?.description ?? '');
    r.setErrorRetry(map.initMfaU2fText?.errorRetry ?? '');
    r.setNotSupported(map.initMfaU2fText?.notSupported ?? '');
    r.setRegisterTokenButtonText(map.initMfaU2fText?.registerTokenButtonText ?? '');
    r.setTitle(map.initMfaU2fText?.title ?? '');
    r.setTokenNameLabel(map.initMfaU2fText?.tokenNameLabel ?? '');

    req.setInitMfaU2fText(r);
  }

  if (!!map.initPasswordDoneText) {
    const r = new InitPasswordDoneScreenText();
    r.setCancelButtonText(map.initPasswordDoneText?.cancelButtonText ?? '');
    r.setDescription(map.initPasswordDoneText?.description ?? '');
    r.setNextButtonText(map.initPasswordDoneText?.nextButtonText ?? '');
    r.setTitle(map.initPasswordDoneText?.title ?? '');

    req.setInitPasswordDoneText(r);
  }

  if (!!map.initPasswordText) {
    const r = new InitPasswordScreenText();
    r.setCodeLabel(map.initPasswordText?.description ?? '');
    r.setDescription(map.initPasswordText?.description ?? '');
    r.setNewPasswordConfirmLabel(map.initPasswordText?.newPasswordConfirmLabel ?? '');
    r.setNewPasswordLabel(map.initPasswordText?.newPasswordLabel ?? '');
    r.setNextButtonText(map.initPasswordText?.nextButtonText ?? '');
    r.setResendButtonText(map.initPasswordText?.resendButtonText ?? '');
    r.setTitle(map.initPasswordText?.title ?? '');

    req.setInitPasswordText(r);
  }


  if (!!map.initializeDoneText) {
    const r = new InitializeUserDoneScreenText();
    r.setCancelButtonText(map.initializeDoneText?.cancelButtonText ?? '');
    r.setDescription(map.initializeDoneText?.description ?? '');
    r.setNextButtonText(map.initializeDoneText?.nextButtonText ?? '');
    r.setTitle(map.initializeDoneText?.title ?? '');

    req.setInitializeDoneText(r);
  }

  if (!!map.initializeUserText) {
    const initializeUserTextRequest = new InitializeUserScreenText();
    initializeUserTextRequest.setCodeLabel(map.initializeUserText?.codeLabel ?? '');
    initializeUserTextRequest.setDescription(map.initializeUserText?.description ?? '');
    initializeUserTextRequest.setNewPasswordConfirmLabel(map.initializeUserText?.newPasswordConfirmLabel ?? '');
    initializeUserTextRequest.setNewPasswordLabel(map.initializeUserText?.newPasswordLabel ?? '');
    initializeUserTextRequest.setNextButtonText(map.initializeUserText?.nextButtonText ?? '');
    initializeUserTextRequest.setResendButtonText(map.initializeUserText?.resendButtonText ?? '');
    initializeUserTextRequest.setTitle(map.initializeUserText?.title ?? '');

    req.setInitializeUserText(initializeUserTextRequest);
  }

  if (!!map.linkingUserDoneText) {
    const r = new LinkingUserDoneScreenText();
    r.setCancelButtonText(map.linkingUserDoneText?.cancelButtonText ?? '');
    r.setDescription(map.linkingUserDoneText?.description ?? '');
    r.setNextButtonText(map.linkingUserDoneText?.nextButtonText ?? '');
    r.setTitle(map.linkingUserDoneText?.title ?? '');

    req.setLinkingUserDoneText(r);
  }

  if (!!map.loginText) {
    const r = new LoginScreenText();
    r.setDescription(map.loginText?.description ?? '');
    r.setDescriptionLinkingProcess(map.loginText?.descriptionLinkingProcess ?? '');
    r.setExternalUserDescription(map.loginText?.externalUserDescription ?? '');
    r.setLoginNameLabel(map.loginText?.loginNameLabel ?? '');
    r.setLoginNamePlaceholder(map.loginText?.loginNamePlaceholder ?? '');
    r.setNextButtonText(map.loginText?.nextButtonText ?? '');
    r.setRegisterButtonText(map.loginText?.registerButtonText ?? '');
    r.setTitle(map.loginText?.title ?? '');
    r.setTitleLinkingProcess(map.loginText?.titleLinkingProcess ?? '');
    r.setUserMustBeMemberOfOrg(map.loginText?.userMustBeMemberOfOrg ?? '');
    r.setUserNamePlaceholder(map.loginText?.userNamePlaceholder ?? '');

    req.setLoginText(r);
  }

  if (!!map.logoutText) {
    const r = new LogoutDoneScreenText();
    r.setDescription(map.logoutText?.description ?? '');
    r.setLoginButtonText(map.logoutText?.loginButtonText ?? '');
    r.setTitle(map.logoutText?.title ?? '');

    req.setLogoutText(r);
  }

  if (!!map.mfaProvidersText) {
    const r = new MFAProvidersText();
    r.setChooseOther(map.mfaProvidersText?.chooseOther ?? '');
    r.setOtp(map.mfaProvidersText?.otp ?? '');
    r.setU2f(map.mfaProvidersText?.u2f ?? '');

    req.setMfaProvidersText(r);
  }

  if (!!map.passwordChangeDoneText) {
    const r = new PasswordChangeDoneScreenText();
    r.setDescription(map.passwordChangeDoneText?.description ?? '');
    r.setNextButtonText(map.passwordChangeDoneText?.nextButtonText ?? '');
    r.setTitle(map.passwordChangeDoneText?.title ?? '');

    req.setPasswordChangeDoneText(r);
  }

  if (!!map.passwordChangeText) {
    const r = new PasswordChangeScreenText();
    r.setDescription(map.passwordChangeText?.description ?? '');
    r.setNextButtonText(map.passwordChangeText?.nextButtonText ?? '');
    r.setTitle(map.passwordChangeText?.title ?? '');

    req.setPasswordChangeText(r);
  }

  if (!!map.passwordResetDoneText) {
    const r = new PasswordResetDoneScreenText();
    r.setDescription(map.passwordResetDoneText?.description ?? '');
    r.setNextButtonText(map.passwordResetDoneText?.nextButtonText ?? '');
    r.setTitle(map.passwordResetDoneText?.title ?? '');

    req.setPasswordResetDoneText(r);
  }

  if (!!map.passwordText) {
    const r = new PasswordScreenText();
    r.setBackButtonText(map.passwordText?.backButtonText ?? '');
    r.setConfirmation(map.passwordText?.confirmation ?? '');
    r.setDescription(map.passwordText?.description ?? '');
    r.setHasLowercase(map.passwordText?.hasLowercase ?? '');
    r.setHasNumber(map.passwordText?.hasNumber ?? '');
    r.setHasSymbol(map.passwordText?.hasSymbol ?? '');
    r.setHasUppercase(map.passwordText?.hasUppercase ?? '');
    r.setMinLength(map.passwordText?.minLength ?? '');
    r.setNextButtonText(map.passwordText?.nextButtonText ?? '');
    r.setPasswordLabel(map.passwordText?.passwordLabel ?? '');
    r.setResetLinkText(map.passwordText?.resetLinkText ?? '');
    r.setTitle(map.passwordText?.title ?? '');

    req.setPasswordText(r);
  }

  if (!!map.passwordlessText) {
    const r = new PasswordlessScreenText();
    r.setDescription(map.passwordlessText?.description ?? '');
    r.setErrorRetry(map.passwordlessText?.errorRetry ?? '');
    r.setLoginWithPwButtonText(map.passwordlessText?.loginWithPwButtonText ?? '');
    r.setNotSupported(map.passwordlessText?.notSupported ?? '');
    r.setTitle(map.passwordlessText?.title ?? '');
    r.setValidateTokenButtonText(map.passwordlessText?.validateTokenButtonText ?? '');

    req.setPasswordlessText(r);
  }

  if (!!map.registrationOptionText) {
    const r = new RegistrationOptionScreenText();
    r.setDescription(map.registrationOptionText?.description ?? '');
    r.setExternalLoginDescription(map.registrationOptionText?.externalLoginDescription ?? '');
    r.setTitle(map.registrationOptionText?.title ?? '');
    r.setUserNameButtonText(map.registrationOptionText?.userNameButtonText ?? '');

    req.setRegistrationOptionText(r);
  }

  if (!!map.registrationOrgText) {
    const r = new RegistrationOrgScreenText();
    r.setDescription(map.registrationOrgText?.description ?? '');
    r.setEmailLabel(map.registrationOrgText?.emailLabel ?? '');
    r.setExternalLoginDescription(map.registrationOrgText?.externalLoginDescription ?? '');
    r.setFirstnameLabel(map.registrationOrgText?.firstnameLabel ?? '');
    r.setLastnameLabel(map.registrationOrgText?.lastnameLabel ?? '');
    r.setOrgnameLabel(map.registrationOrgText?.orgnameLabel ?? '');
    r.setPasswordConfirmLabel(map.registrationOrgText?.passwordConfirmLabel ?? '');
    r.setPasswordLabel(map.registrationOrgText?.passwordLabel ?? '');
    r.setPrivacyConfirm(map.registrationOrgText?.privacyConfirm ?? '');
    r.setPrivacyLink(map.registrationOrgText?.privacyLink ?? '');
    r.setPrivacyLinkText(map.registrationOrgText?.privacyLinkText ?? '');
    r.setSaveButtonText(map.registrationOrgText?.saveButtonText ?? '');
    r.setTitle(map.registrationOrgText?.title ?? '');
    r.setTosAndPrivacyLabel(map.registrationOrgText?.tosAndPrivacyLabel ?? '');
    r.setTosConfirm(map.registrationOrgText?.tosConfirm ?? '');
    r.setTosLink(map.registrationOrgText?.tosLink ?? '');
    r.setTosLinkText(map.registrationOrgText?.tosLinkText ?? '');
    r.setUsernameLabel(map.registrationOrgText?.usernameLabel ?? '');

    req.setRegistrationOrgText(r);
  }

  if (!!map.registrationUserText) {
    const r = new RegistrationUserScreenText();
    r.setBackButtonText(map.registrationUserText?.backButtonText ?? '');
    r.setDescription(map.registrationUserText?.description ?? '');
    r.setDescriptionOrgRegister(map.registrationUserText?.descriptionOrgRegister ?? '');
    r.setEmailLabel(map.registrationUserText?.emailLabel ?? '');
    r.setExternalLoginDescription(map.registrationUserText?.externalLoginDescription ?? '');
    r.setFirstnameLabel(map.registrationUserText?.firstnameLabel ?? '');
    r.setGenderLabel(map.registrationUserText?.genderLabel ?? '');
    r.setLanguageLabel(map.registrationUserText?.languageLabel ?? '');
    r.setLastnameLabel(map.registrationUserText?.lastnameLabel ?? '');
    r.setNextButtonText(map.registrationUserText?.nextButtonText ?? '');
    r.setPasswordConfirmLabel(map.registrationUserText?.passwordConfirmLabel ?? '');
    r.setPasswordLabel(map.registrationUserText?.passwordLabel ?? '');
    r.setPrivacyConfirm(map.registrationUserText?.privacyConfirm ?? '');
    r.setPrivacyLink(map.registrationUserText?.privacyLink ?? '');
    r.setPrivacyLinkText(map.registrationUserText?.privacyLinkText ?? '');
    r.setTitle(map.registrationUserText?.title ?? '');
    r.setTosAndPrivacyLabel(map.registrationUserText?.tosAndPrivacyLabel ?? '');
    r.setTosConfirm(map.registrationUserText?.tosConfirm ?? '');
    r.setTosLink(map.registrationUserText?.tosLink ?? '');
    r.setTosLinkText(map.registrationUserText?.tosLinkText ?? '');
    r.setUsernameLabel(map.registrationUserText?.usernameLabel ?? '');

    req.setRegistrationUserText(r);
  }

  if (!!map.selectAccountText) {
    const r = new SelectAccountScreenText();
    r.setDescription(map.selectAccountText?.description ?? '');
    r.setDescriptionLinkingProcess(map.selectAccountText?.descriptionLinkingProcess ?? '');
    r.setOtherUser(map.selectAccountText?.otherUser ?? '');
    r.setSessionStateActive(map.selectAccountText?.sessionStateActive ?? '');
    r.setSessionStateInactive(map.selectAccountText?.sessionStateInactive ?? '');
    r.setTitle(map.selectAccountText?.title ?? '');
    r.setTitleLinkingProcess(map.selectAccountText?.titleLinkingProcess ?? '');
    r.setUserMustBeMemberOfOrg(map.selectAccountText?.userMustBeMemberOfOrg ?? '');

    req.setSelectAccountText(r);
  }


  if (!!map.successLoginText) {
    const r = new SuccessLoginScreenText();
    r.setAutoRedirectDescription(map.successLoginText?.autoRedirectDescription ?? '');
    r.setNextButtonText(map.successLoginText?.nextButtonText ?? '');
    r.setRedirectedDescription(map.successLoginText?.redirectedDescription ?? '');
    r.setTitle(map.successLoginText?.title ?? '');

    req.setSuccessLoginText(r);
  }


  if (!!map.usernameChangeDoneText) {
    const r = new UsernameChangeDoneScreenText();
    r.setDescription(map.usernameChangeDoneText?.description ?? '');
    r.setNextButtonText(map.usernameChangeDoneText?.nextButtonText ?? '');
    r.setTitle(map.usernameChangeDoneText?.title ?? '');

    req.setUsernameChangeDoneText(r);
  }


  if (!!map.usernameChangeText) {
    const r = new UsernameChangeScreenText();
    r.setCancelButtonText(map.usernameChangeText?.cancelButtonText ?? '');
    r.setDescription(map.usernameChangeText?.description ?? '');
    r.setNextButtonText(map.usernameChangeText?.nextButtonText ?? '');
    r.setTitle(map.usernameChangeText?.title ?? '');
    r.setUsernameLabel(map.usernameChangeText?.usernameLabel ?? '');

    req.setUsernameChangeText(r);
  }

  if (!!map.verifyMfaOtpText) {
    const r = new VerifyMFAOTPScreenText();
    r.setCodeLabel(map.verifyMfaOtpText?.codeLabel ?? '');
    r.setDescription(map.verifyMfaOtpText?.description ?? '');
    r.setNextButtonText(map.verifyMfaOtpText?.nextButtonText ?? '');
    r.setTitle(map.verifyMfaOtpText?.title ?? '');

    req.setVerifyMfaOtpText(r);
  }


  if (!!map.verifyMfaU2fText) {
    const r = new VerifyMFAU2FScreenText();
    r.setDescription(map.verifyMfaU2fText?.description ?? '');
    r.setErrorRetry(map.verifyMfaU2fText?.errorRetry ?? '');
    r.setNotSupported(map.verifyMfaU2fText?.notSupported ?? '');
    r.setTitle(map.verifyMfaU2fText?.title ?? '');
    r.setValidateTokenText(map.verifyMfaU2fText?.validateTokenText ?? '');

    req.setVerifyMfaU2fText(r);
  }

  return req;
}

import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, from, Observable, of, Subscription } from 'rxjs';
import { map, switchMap } from 'rxjs/operators';
import {
  GetDefaultLoginTextsRequest as AdminGetDefaultLoginTextsRequest,
  SetCustomLoginTextsRequest as AdminSetCustomLoginTextsRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  GetCustomLoginTextsRequest,
  GetDefaultLoginTextsRequest,
  SetCustomLoginTextsRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
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
  PasswordResetDoneScreenText,
  PasswordScreenText,
} from 'src/app/proto/generated/zitadel/text_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';

import { CnslLinks } from '../../links/links.component';
import {
  IAM_COMPLEXITY_LINK,
  IAM_POLICY_LINK,
  IAM_PRIVATELABEL_LINK,
  ORG_COMPLEXITY_LINK,
  ORG_IAM_POLICY_LINK,
  ORG_PRIVATELABEL_LINK,
} from '../../policy-grid/policy-links';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

type ArgumentTypes<F extends Function> = F extends (...args: infer A) => any ? A : never;
type Parameters<T> = T extends (...args: infer T) => any ? T : never;

type ParameterNames = Parameters<ManagementService["setCustomLoginText"]>;
// type KeyNames = keyof SetCustomLoginTextsRequest.AsObject;
// const KeyNamesArray = [
//   'setEmailVerificationDoneText',
//   'setEmailVerificationText',
//   'setExternalUserNotFoundText',
//   'setFooterText',
//   'setInitMfaDoneText',
//   'setInitMfaDoneText',
//   'setInitMfaOtpText',
//   'setInitMfaPromptText',
//   'setInitMfaU2fText',
//   'setInitPasswordDoneText',
//   'setInitPasswordText',
//   'setInitializeDoneText',
//   'setInitializeUserText',
//   'setLinkingUserDoneText',
//   'setLoginText',
//   'setLogoutText',
//   'setMfaProvidersText',
//   'setPasswordChangeDoneText',
//   'setPasswordChangeText',
//   'setPasswordResetDoneText',
//   'setPasswordText',
//   'setPasswordlessText',
//   'setRegistrationOptionText',
//   'setRegistrationOrgText',
//   'setRegistrationUserText',
//   'setSelectAccountText',
//   'setSuccessLoginText',
//   'setUsernameChangeDoneText',
//   'setUsernameChangeText',
//   'setVerifyMfaOtpText',
//   'setVerifyMfaU2fText'
// ];
const KeyNamesArray = [
  'emailVerificationDoneText',
  'emailVerificationText',
  'externalUserNotFoundText',
  'footerText',
  'initMfaDoneText',
  'initMfaDoneText',
  'initMfaOtpText',
  'initMfaPromptText',
  'initMfaU2fText',
  'initPasswordDoneText',
  'initPasswordText',
  'initializeDoneText',
  'initializeUserText',
  'linkingUserDoneText',
  'loginText',
  'logoutText',
  'mfaProvidersText',
  'passwordChangeDoneText',
  'passwordChangeText',
  'passwordResetDoneText',
  'passwordText',
  'passwordlessText',
  'registrationOptionText',
  'registrationOrgText',
  'registrationUserText',
  'selectAccountText',
  'successLoginText',
  'usernameChangeDoneText',
  'usernameChangeText',
  'verifyMfaOtpText',
  'verifyMfaU2fText'
];
type KeyName = keyof typeof KeyNamesArray;

const REQUESTMAP = {
  [PolicyComponentServiceType.MGMT]: {
    get: new GetCustomLoginTextsRequest(),
    set: new SetCustomLoginTextsRequest(),
    getDefault: new GetDefaultLoginTextsRequest(),
    setFcn: (map: Partial<SetCustomLoginTextsRequest.AsObject>): SetCustomLoginTextsRequest => {
      const req = new SetCustomLoginTextsRequest();
      // req.setEmailVerificationDoneText(map.emailVerificationDoneText ?? '');
      // req.setEmailVerificationText(map.emailVerificationText ?? '');
      // req.setExternalUserNotFoundText(map.externalUserNotFoundText ?? '');
      // req.setFooterText(map.footerText ?? '');
      // req.setInitMfaDoneText(map.initMfaDoneText ?? '');
      // req.setInitMfaOtpText(map.initMfaOtpText ?? '');
      // req.setInitMfaPromptText(map.initMfaPromptText ?? '');
      // req.setInitMfaU2fText(map.initMfaU2fText ?? '');
      // req.setInitPasswordDoneText(map.initPasswordDoneText ?? '');
      // req.setInitPasswordText(map.initPasswordText ?? '');
      // req.setInitializeDoneText(map.initializeDoneText ?? '');
      // req.setInitializeUserText(map.initializeUserText ?? '');
      req.setLanguage(map.language ?? '');
      // req.setLinkingUserDoneText(map.linkingUserDoneText ?? '');
      // req.setLoginText(map.loginText ?? '');
      // req.setLogoutText(map.logoutText ?? '');
      // req.setMfaProvidersText(map.mfaProvidersText ?? '');
      // req.setPasswordChangeDoneText(map.passwordChangeDoneText ?? '');
      // req.setPasswordChangeText(map.passwordChangeText ?? '');
      // req.setPasswordResetDoneText(map.passwordResetDoneText ?? '');
      // req.setPasswordText(map.passwordText ?? '');
      // req.setPasswordlessText(map.passwordlessText ?? '');
      // req.setRegistrationOptionText(map.registrationOptionText ?? '');
      // req.setRegistrationOrgText(map.registrationOrgText ?? '');
      // req.setRegistrationUserText(map.registrationUserText ?? '');
      // req.setSelectAccountText(map.selectAccountText ?? '');
      // req.setSuccessLoginText(map.successLoginText ?? '');
      // req.setUsernameChangeDoneText(map.usernameChangeDoneText ?? '');
      // req.setUsernameChangeText(map.usernameChangeText ?? '');
      // req.setVerifyMfaOtpText(map.verifyMfaOtpText ?? '');
      // req.setVerifyMfaU2fText(map.verifyMfaU2fText ?? '');

      return req;
    }
  },
  [PolicyComponentServiceType.ADMIN]: {
    get: new AdminGetDefaultLoginTextsRequest(),
    set: new AdminSetCustomLoginTextsRequest(),
    setFcn: (map: Partial<AdminSetCustomLoginTextsRequest.AsObject>): AdminSetCustomLoginTextsRequest => {
      const req = new AdminSetCustomLoginTextsRequest();
      req.setLanguage(map.language ?? '');

      map.emailVerificationDoneText ? () => {
        const r = new EmailVerificationDoneScreenText();
        r.setCancelButtonText(map.emailVerificationDoneText?.cancelButtonText ?? '');
        r.setDescription(map.emailVerificationDoneText?.description ?? '');
        r.setLoginButtonText(map.emailVerificationDoneText?.loginButtonText ?? '');
        r.setNextButtonText(map.initializeUserText?.nextButtonText ?? '');
        r.setTitle(map.initializeUserText?.title ?? '');

        req.setEmailVerificationDoneText(r);
      } : null;

      map.emailVerificationText ? () => {
        const r = new EmailVerificationScreenText();
        r.setCodeLabel(map.emailVerificationText?.codeLabel ?? '');
        r.setDescription(map.emailVerificationText?.description ?? '');
        r.setNextButtonText(map.emailVerificationText?.nextButtonText ?? '');
        r.setResendButtonText(map.emailVerificationText?.resendButtonText ?? '');
        r.setTitle(map.emailVerificationText?.title ?? '');

        req.setEmailVerificationText(r);
      } : null;

      map.externalUserNotFoundText ? () => {
        const r = new ExternalUserNotFoundScreenText();
        r.setAutoRegisterButtonText(map.externalUserNotFoundText?.autoRegisterButtonText ?? '');
        r.setDescription(map.externalUserNotFoundText?.description ?? '');
        r.setLinkButtonText(map.externalUserNotFoundText?.linkButtonText ?? '');
        r.setTitle(map.externalUserNotFoundText?.title ?? '');

        req.setExternalUserNotFoundText(r);
      } : null;

      map.footerText ? () => {
        const r = new FooterText();
        r.setHelp(map.footerText?.help ?? '');
        r.setHelpLink(map.footerText?.helpLink ?? '');
        r.setPrivacyPolicy(map.footerText?.privacyPolicy ?? '');
        r.setPrivacyPolicyLink(map.footerText?.privacyPolicyLink ?? '');
        r.setTos(map.footerText?.tos ?? '');
        r.setTosLink(map.footerText?.tosLink ?? '');

        req.setFooterText(r);
      } : null;

      map.initMfaDoneText ? () => {
        const r = new InitMFADoneScreenText();
        r.setCancelButtonText(map.initMfaDoneText?.cancelButtonText ?? '');
        r.setDescription(map.initMfaDoneText?.description ?? '');
        r.setNextButtonText(map.initMfaDoneText?.nextButtonText ?? '');
        r.setTitle(map.initMfaDoneText?.title ?? '');

        req.setInitMfaDoneText(r);
      } : null;

      map.initMfaOtpText ? () => {
        const r = new InitMFAOTPScreenText();
        r.setCancelButtonText(map.initMfaOtpText?.cancelButtonText ?? '');
        r.setCodeLabel(map.initMfaOtpText?.codeLabel ?? '');
        r.setDescription(map.initMfaOtpText?.description ?? '');
        r.setDescriptionOtp(map.initMfaOtpText?.descriptionOtp ?? '');
        r.setNextButtonText(map.initMfaOtpText?.nextButtonText ?? '');
        r.setSecretLabel(map.initMfaOtpText?.secretLabel ?? '');
        r.setTitle(map.initMfaOtpText?.title ?? '');

        req.setInitMfaOtpText(r);
      } : null;

      map.initMfaPromptText ? () => {
        const r = new InitMFAPromptScreenText();
        r.setDescription(map.initMfaPromptText?.description ?? '');
        r.setNextButtonText(map.initMfaPromptText?.nextButtonText ?? '');
        r.setOtpOption(map.initMfaPromptText?.otpOption ?? '');
        r.setSkipButtonText(map.initMfaPromptText?.skipButtonText ?? '');
        r.setTitle(map.initMfaPromptText?.title ?? '');
        r.setU2fOption(map.initMfaPromptText?.otpOption ?? '');

        req.setInitMfaPromptText(r);
      } : null;

      map.initMfaU2fText ? () => {
        const r = new InitMFAU2FScreenText();
        r.setDescription(map.initMfaU2fText?.description ?? '');
        r.setErrorRetry(map.initMfaU2fText?.errorRetry ?? '');
        r.setNotSupported(map.initMfaU2fText?.notSupported ?? '');
        r.setRegisterTokenButtonText(map.initMfaU2fText?.registerTokenButtonText ?? '');
        r.setTitle(map.initMfaU2fText?.title ?? '');
        r.setTokenNameLabel(map.initMfaU2fText?.tokenNameLabel ?? '');

        req.setInitMfaU2fText(r);
      } : null;

      map.initPasswordDoneText ? () => {
        const r = new InitPasswordDoneScreenText();
        r.setCancelButtonText(map.initPasswordDoneText?.cancelButtonText ?? '');
        r.setDescription(map.initPasswordDoneText?.description ?? '');
        r.setNextButtonText(map.initPasswordDoneText?.nextButtonText ?? '');
        r.setTitle(map.initPasswordDoneText?.title ?? '');

        req.setInitPasswordDoneText(r);
      } : null;

      map.initPasswordText ? () => {
        const r = new InitPasswordScreenText();
        r.setCodeLabel(map.initPasswordText?.description ?? '');
        r.setDescription(map.initPasswordText?.description ?? '');
        r.setNewPasswordConfirmLabel(map.initPasswordText?.newPasswordConfirmLabel ?? '');
        r.setNewPasswordLabel(map.initPasswordText?.newPasswordLabel ?? '');
        r.setNextButtonText(map.initPasswordText?.nextButtonText ?? '');
        r.setResendButtonText(map.initPasswordText?.resendButtonText ?? '');
        r.setTitle(map.initPasswordText?.title ?? '');

        req.setInitPasswordText(r);
      } : null;

      map.initializeDoneText ? () => {
        const r = new InitializeUserDoneScreenText();
        r.setCancelButtonText(map.initializeDoneText?.cancelButtonText ?? '');
        r.setDescription(map.initializeDoneText?.description ?? '');
        r.setNextButtonText(map.initializeDoneText?.nextButtonText ?? '');
        r.setTitle(map.initializeDoneText?.title ?? '');

        req.setInitializeDoneText(r);
      } : null;

      map.initializeUserText ? () => {
        const initializeUserTextRequest = new InitializeUserScreenText();
        initializeUserTextRequest.setCodeLabel(map.initializeUserText?.codeLabel ?? '');
        initializeUserTextRequest.setDescription(map.initializeUserText?.description ?? '');
        initializeUserTextRequest.setNewPasswordConfirmLabel(map.initializeUserText?.newPasswordConfirmLabel ?? '');
        initializeUserTextRequest.setNewPasswordLabel(map.initializeUserText?.newPasswordLabel ?? '');
        initializeUserTextRequest.setNextButtonText(map.initializeUserText?.nextButtonText ?? '');
        initializeUserTextRequest.setResendButtonText(map.initializeUserText?.resendButtonText ?? '');
        initializeUserTextRequest.setTitle(map.initializeUserText?.title ?? '');

        req.setInitializeUserText(initializeUserTextRequest);
      } : null;

      map.linkingUserDoneText ? () => {
        const r = new LinkingUserDoneScreenText();
        r.setCancelButtonText(map.linkingUserDoneText?.cancelButtonText ?? '');
        r.setDescription(map.linkingUserDoneText?.description ?? '');
        r.setNextButtonText(map.linkingUserDoneText?.nextButtonText ?? '');
        r.setTitle(map.linkingUserDoneText?.title ?? '');

        req.setLinkingUserDoneText(r);
      } : null;

      map.loginText ? () => {
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
      } : null;

      map.logoutText ? () => {
        const r = new LogoutDoneScreenText();
        r.setDescription(map.logoutText?.description ?? '');
        r.setLoginButtonText(map.logoutText?.loginButtonText ?? '');
        r.setTitle(map.logoutText?.title ?? '');

        req.setLogoutText(r);
      } : null;

      map.mfaProvidersText ? () => {
        const r = new MFAProvidersText();
        r.setChooseOther(map.mfaProvidersText?.chooseOther ?? '');
        r.setOtp(map.mfaProvidersText?.otp ?? '');
        r.setU2f(map.mfaProvidersText?.u2f ?? '');

        req.setMfaProvidersText(r);
      } : null;

      map.passwordChangeDoneText ? () => {
        const r = new PasswordChangeDoneScreenText();
        r.setDescription(map.passwordChangeDoneText?.description ?? '');
        r.setNextButtonText(map.passwordChangeDoneText?.nextButtonText ?? '');
        r.setTitle(map.passwordChangeDoneText?.title ?? '');

        req.setPasswordChangeDoneText(r);
      } : null;

      map.passwordChangeText ? () => {
        const r = new PasswordChangeScreenText();
        r.setDescription(map.passwordChangeText?.description ?? '');
        r.setNextButtonText(map.passwordChangeText?.nextButtonText ?? '');
        r.setTitle(map.passwordChangeText?.title ?? '');

        req.setPasswordChangeText(r);
      } : null;

      map.passwordResetDoneText ? () => {
        const r = new PasswordResetDoneScreenText();
        r.setDescription(map.passwordResetDoneText?.description ?? '');
        r.setNextButtonText(map.passwordResetDoneText?.nextButtonText ?? '');
        r.setTitle(map.passwordResetDoneText?.title ?? '');

        req.setPasswordResetDoneText(r);
      } : null;

      map.passwordText ? () => {
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
      } : null;

      // req.setPasswordlessText(map.passwordlessText ?? '');
      // req.setRegistrationOptionText(map.registrationOptionText ?? '');
      // req.setRegistrationOrgText(map.registrationOrgText ?? '');
      // req.setRegistrationUserText(map.registrationUserText ?? '');
      // req.setSelectAccountText(map.selectAccountText ?? '');
      // req.setSuccessLoginText(map.successLoginText ?? '');
      // req.setUsernameChangeDoneText(map.usernameChangeDoneText ?? '');
      // req.setUsernameChangeText(map.usernameChangeText ?? '');
      // req.setVerifyMfaOtpText(map.verifyMfaOtpText ?? '');
      // req.setVerifyMfaU2fText(map.verifyMfaU2fText ?? '');

      return req;
    }
  },
};
@Component({
  selector: 'app-login-texts',
  templateUrl: './login-texts.component.html',
  styleUrls: ['./login-texts.component.scss'],
})
export class LoginTextsComponent implements OnDestroy {
  public getDefaultInitMessageTextMap$: Observable<{ [key: string]: string; }> = of({});
  public getCustomInitMessageTextMap$: BehaviorSubject<{ [key: string]: string; }> = new BehaviorSubject({});

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public nextLinks: CnslLinks[] = [];

  public currentSubMap: string = 'emailVerificationDoneText';

  public KeyNamesArray: string[] = KeyNamesArray;
  private sub: Subscription = new Subscription();
  constructor(
    private route: ActivatedRoute,
    private injector: Injector,
    private translate: TranslateService,
  ) {
    this.sub = this.route.data.pipe(switchMap(data => {
      this.serviceType = data.serviceType;
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          this.service = this.injector.get(ManagementService as Type<ManagementService>);
          this.nextLinks = [
            ORG_COMPLEXITY_LINK,
            ORG_IAM_POLICY_LINK,
            ORG_PRIVATELABEL_LINK,
          ];

          // this.setCurrentType('emailVerificationDoneText');
          this.loadData();
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);
          this.nextLinks = [
            IAM_COMPLEXITY_LINK,
            IAM_POLICY_LINK,
            IAM_PRIVATELABEL_LINK,
          ];
          break;
      }

      return this.route.params;
    })).subscribe(() => {

    });
  }

  public getDefaultValues(req: any): Promise<any> {
    return this.stripDetails((this.service).getDefaultLoginTexts(req));
  }

  public getCurrentValues(req: any): Promise<any> {
    return this.stripDetails((this.service as ManagementService).getCustomLoginTexts(req));
  }

  public async loadData() {
    const lang = this.translate.currentLang ?? 'en';
    if (this.serviceType == PolicyComponentServiceType.MGMT) {
      const reqDefaultInit = REQUESTMAP[this.serviceType].getDefault;


      reqDefaultInit.setLanguage(lang);
      this.getDefaultInitMessageTextMap$ = from(
        this.getDefaultValues(reqDefaultInit)
      ).pipe(map(m => m[this.currentSubMap]));
    }

    const reqCustomInit = REQUESTMAP[this.serviceType].get.setLanguage(lang);
    this.getCustomInitMessageTextMap$.next(
      (await this.getCurrentValues(reqCustomInit))[this.currentSubMap]
    );
  }

  public updateCurrentValues(values: { [key: string]: string; }): void {
    const req = REQUESTMAP[this.serviceType].setFcn;
    const mappedValues = req({ [this.currentSubMap]: values });

    console.log(mappedValues.toObject());
  }

  public saveCurrentMessage(): void {
    console.log('save');
  }

  private stripDetails(prom: Promise<any>): Promise<any> {
    return prom.then(res => {
      if (res.customText) {
        delete res.customText.details;
        console.log(Object.assign({}, res.customText));
        return Object.assign({}, res.customText);
      } else {
        return {};
      }
    });
  }
  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public async setCurrentType(key: string): Promise<void> {
    this.currentSubMap = key;

    this.loadData();
  }
};

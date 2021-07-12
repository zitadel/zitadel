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
    const mappedValues = req(values);

    console.log(mappedValues.toObject());
  }

  public saveCurrentMessage(): void {
    console.log('save');
  }

  private stripDetails(prom: Promise<any>): Promise<any> {
    return prom.then(res => {
      if (res.customText) {
        delete res.customText.details;
        return Object.assign({}, res.customText as unknown as { [key: string]: string; });
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

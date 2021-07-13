import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
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
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { mapRequestValues } from './helper';

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
// type KeyName = keyof typeof KeyNamesArray;

const REQUESTMAP = {
  [PolicyComponentServiceType.MGMT]: {
    get: new GetCustomLoginTextsRequest(),
    set: new SetCustomLoginTextsRequest(),
    getDefault: new GetDefaultLoginTextsRequest(),
    setFcn: (map: Partial<SetCustomLoginTextsRequest.AsObject>): SetCustomLoginTextsRequest => {
      console.log(map);
      let req = new SetCustomLoginTextsRequest();
      req.setLanguage(map.language ?? '');
      req = mapRequestValues(map, req);
      return req;
    }
  },
  [PolicyComponentServiceType.ADMIN]: {
    get: new AdminGetDefaultLoginTextsRequest(),
    set: new AdminSetCustomLoginTextsRequest(),
    setFcn: (map: Partial<AdminSetCustomLoginTextsRequest.AsObject>): AdminSetCustomLoginTextsRequest => {
      let req = new AdminSetCustomLoginTextsRequest();
      req.setLanguage(map.language ?? '');
      req = mapRequestValues(map, req);
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
    private dialog: MatDialog,
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
    console.log(values);
    const req = REQUESTMAP[this.serviceType].setFcn;
    const mappedValues = req({ [this.currentSubMap]: values });

    console.log(mappedValues.toObject());
  }

  public saveCurrentMessage(): void {
    console.log('save');
  }

  public resetDefault(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        icon: 'las la-history',
        confirmKey: 'ACTIONS.RESTORE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'POLICY.LOGIN_TEXTS.RESET_TITLE',
        descriptionKey: 'POLICY.LOGIN_TEXTS.RESET_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe(resp => {
      if (resp) {

      }
    });
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
}
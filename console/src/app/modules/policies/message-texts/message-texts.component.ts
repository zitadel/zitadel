import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, from, Observable, of, Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import {
  GetCustomPasswordResetMessageTextRequest as AdminGetCustomPasswordResetMessageTextRequest,
  GetDefaultInitMessageTextRequest as AdminGetDefaultInitMessageTextRequest,
  GetDefaultVerifyEmailMessageTextRequest as AdminGetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextRequest as AdminGetDefaultVerifyPhoneMessageTextRequest,
  SetDefaultDomainClaimedMessageTextRequest,
  SetDefaultInitMessageTextRequest,
  SetDefaultPasswordResetMessageTextRequest,
  SetDefaultVerifyEmailMessageTextRequest,
  SetDefaultVerifyPhoneMessageTextRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  GetCustomDomainClaimedMessageTextRequest,
  GetCustomPasswordResetMessageTextRequest,
  GetCustomVerifyEmailMessageTextRequest,
  GetCustomVerifyPhoneMessageTextRequest,
  GetDefaultDomainClaimedMessageTextRequest,
  GetDefaultInitMessageTextRequest,
  GetDefaultPasswordResetMessageTextRequest,
  GetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextRequest,
  SetCustomDomainClaimedMessageTextRequest,
  SetCustomInitMessageTextRequest,
  SetCustomPasswordResetMessageTextRequest,
  SetCustomVerifyEmailMessageTextRequest,
  SetCustomVerifyPhoneMessageTextRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { MessageCustomText } from 'src/app/proto/generated/zitadel/text_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

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

enum MESSAGETYPES {
  INIT = 'INIT',
  VERIFYPHONE = 'VP',
  VERIFYEMAIL = 'VE',
  PASSWORDRESET = 'PR',
  DOMAINCLAIMED = 'DC'
};

const REQUESTMAP = {
  [PolicyComponentServiceType.MGMT]: {
    [MESSAGETYPES.INIT]: {
      get: new GetDefaultInitMessageTextRequest(),
      set: new SetCustomInitMessageTextRequest(),
      getDefault: new GetDefaultInitMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomInitMessageTextRequest => {
        const req = new SetCustomInitMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
    [MESSAGETYPES.VERIFYEMAIL]: {
      get: new GetCustomVerifyEmailMessageTextRequest(),
      set: new SetCustomVerifyEmailMessageTextRequest(),
      getDefault: new GetDefaultVerifyEmailMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomVerifyEmailMessageTextRequest => {
        const req = new SetCustomVerifyEmailMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
    [MESSAGETYPES.VERIFYPHONE]: {
      get: new GetCustomVerifyPhoneMessageTextRequest(),
      set: new SetCustomVerifyPhoneMessageTextRequest(),
      getDefault: new GetDefaultVerifyPhoneMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetCustomVerifyPhoneMessageTextRequest => {
        const req = new SetCustomVerifyPhoneMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
    [MESSAGETYPES.PASSWORDRESET]: {
      get: new GetCustomPasswordResetMessageTextRequest(),
      set: new SetCustomPasswordResetMessageTextRequest(),
      getDefault: new GetDefaultPasswordResetMessageTextRequest(),
      setFcn: (map: Partial<SetCustomPasswordResetMessageTextRequest.AsObject>): SetCustomPasswordResetMessageTextRequest => {
        const req = new SetCustomPasswordResetMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
    [MESSAGETYPES.DOMAINCLAIMED]: {
      get: new GetCustomDomainClaimedMessageTextRequest(),
      set: new SetCustomDomainClaimedMessageTextRequest(),
      getDefault: new GetDefaultDomainClaimedMessageTextRequest(),
      setFcn: (map: Partial<SetCustomDomainClaimedMessageTextRequest.AsObject>): SetCustomDomainClaimedMessageTextRequest => {
        const req = new SetCustomDomainClaimedMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
  },
  [PolicyComponentServiceType.ADMIN]: {
    [MESSAGETYPES.INIT]: {
      get: new AdminGetDefaultInitMessageTextRequest(),
      set: new SetDefaultInitMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultInitMessageTextRequest => {
        const req = new SetDefaultInitMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
    [MESSAGETYPES.VERIFYEMAIL]: {
      get: new AdminGetDefaultVerifyEmailMessageTextRequest(),
      set: new SetDefaultVerifyEmailMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultVerifyEmailMessageTextRequest => {
        const req = new SetDefaultVerifyEmailMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
    [MESSAGETYPES.VERIFYPHONE]: {
      get: new AdminGetDefaultVerifyPhoneMessageTextRequest(),
      set: new SetDefaultVerifyPhoneMessageTextRequest(),
      setFcn: (map: Partial<MessageCustomText.AsObject>): SetDefaultVerifyPhoneMessageTextRequest => {
        const req = new SetDefaultVerifyPhoneMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      },
    },
    [MESSAGETYPES.PASSWORDRESET]: {
      get: new AdminGetCustomPasswordResetMessageTextRequest(),
      set: new SetDefaultPasswordResetMessageTextRequest(),
      setFcn: (map: Partial<SetDefaultPasswordResetMessageTextRequest.AsObject>): SetDefaultPasswordResetMessageTextRequest => {
        const req = new SetDefaultPasswordResetMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
    [MESSAGETYPES.DOMAINCLAIMED]: {
      get: new GetDefaultDomainClaimedMessageTextRequest(),
      set: new SetDefaultDomainClaimedMessageTextRequest(),
      setFcn: (map: Partial<SetDefaultDomainClaimedMessageTextRequest.AsObject>): SetDefaultDomainClaimedMessageTextRequest => {
        const req = new SetDefaultDomainClaimedMessageTextRequest();
        req.setButtonText(map.buttonText ?? '');
        req.setFooterText(map.footerText ?? '');
        req.setGreeting(map.greeting ?? '');
        req.setPreHeader(map.preHeader ?? '');
        req.setSubject(map.subject ?? '');
        req.setText(map.text ?? '');
        req.setTitle(map.title ?? '');

        return req;
      }
    },
  },
};
@Component({
  selector: 'app-message-texts',
  templateUrl: './message-texts.component.html',
  styleUrls: ['./message-texts.component.scss'],
})
export class MessageTextsComponent implements OnDestroy {
  public getDefaultInitMessageTextMap$: Observable<{ [key: string]: string; }> = of({});
  public getCustomInitMessageTextMap$: BehaviorSubject<{ [key: string]: string; }> = new BehaviorSubject({});

  public currentType: MESSAGETYPES = MESSAGETYPES.INIT;

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public nextLinks: CnslLinks[] = [];
  public MESSAGETYPES: any = MESSAGETYPES;

  public updateRequest!: SetCustomInitMessageTextRequest | SetDefaultInitMessageTextRequest;

  public chips: any[] = [
    { key: 'POLICY.MESSAGE_TEXTS.CHIPS.firstname', value: '{{.FirstName}}' },
    { key: 'POLICY.MESSAGE_TEXTS.CHIPS.lastname', value: '{{.Lastname}}' },
    { key: 'POLICY.MESSAGE_TEXTS.CHIPS.code', value: '{{.Code}}' },
    { key: 'POLICY.MESSAGE_TEXTS.CHIPS.preferredLoginName', value: '{{.PreferredLoginName}}' },
  ];

  private sub: Subscription = new Subscription();

  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
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

          this.loadData(this.currentType);

          // this.defaultInitMsg = of(req);
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

  public getDefaultValues(type: MESSAGETYPES, req: any): Promise<any> {
    switch (type) {
      case MESSAGETYPES.INIT:
        return this.stripDetails((this.service).getDefaultInitMessageText(req));
      case MESSAGETYPES.VERIFYPHONE:
        return this.stripDetails((this.service).getDefaultVerifyPhoneMessageText(req));
      case MESSAGETYPES.VERIFYEMAIL:
        return this.stripDetails((this.service).getDefaultVerifyEmailMessageText(req));
      case MESSAGETYPES.PASSWORDRESET:
        return this.stripDetails((this.service).getDefaultPasswordResetMessageText(req));
      case MESSAGETYPES.DOMAINCLAIMED:
        return this.stripDetails((this.service).getDefaultDomainClaimedMessageText(req));
    }
  }

  public getCurrentValues(type: MESSAGETYPES, req: any): Promise<any> {
    switch (type) {
      case MESSAGETYPES.INIT:
        return this.stripDetails((this.service as ManagementService).getCustomInitMessageText(req));
      case MESSAGETYPES.VERIFYPHONE:
        return this.stripDetails((this.service as ManagementService).getCustomVerifyPhoneMessageText(req));
      case MESSAGETYPES.VERIFYEMAIL:
        return this.stripDetails((this.service as ManagementService).getCustomVerifyEmailMessageText(req));
      case MESSAGETYPES.PASSWORDRESET:
        return this.stripDetails((this.service as ManagementService).getCustomPasswordResetMessageText(req));
      case MESSAGETYPES.DOMAINCLAIMED:
        return this.stripDetails((this.service as ManagementService).getCustomDomainClaimedMessageText(req));
    }
  }

  public async loadData(type: MESSAGETYPES) {

    console.log(this.serviceType, type);
    if (this.serviceType == PolicyComponentServiceType.MGMT) {
      const reqDefaultInit = REQUESTMAP[this.serviceType][type].getDefault;

      reqDefaultInit.setLanguage(this.translate.currentLang);
      this.getDefaultInitMessageTextMap$ = from(
        this.getDefaultValues(type, reqDefaultInit)
      );
    }

    const reqCustomInit = REQUESTMAP[this.serviceType][type].get.setLanguage(this.translate.currentLang);
    this.getCustomInitMessageTextMap$.next(
      await this.getCurrentValues(type, reqCustomInit)
    );
  }

  public updateCurrentValues(values: { [key: string]: string; }): void {
    const req = REQUESTMAP[this.serviceType][this.currentType].setFcn;
    const mappedValues = req(values);
    this.updateRequest = mappedValues;
    this.updateRequest.setLanguage(this.translate.currentLang);
    console.log(mappedValues.toObject());
  }

  public saveCurrentMessage(): void {
    if (this.serviceType == PolicyComponentServiceType.MGMT) {
      (this.service as ManagementService).setCustomInitMessageText(this.updateRequest).then(() => {
        this.toast.showInfo('POLICY.MESSAGE_TEXTS.TOAST.UPDATED', true);
      }).catch(error => this.toast.showError(error));
    } else if (this.serviceType == PolicyComponentServiceType.ADMIN) {
      (this.service as AdminService).setDefaultInitMessageText(this.updateRequest).then(() => {
        this.toast.showInfo('POLICY.MESSAGE_TEXTS.TOAST.UPDATED', true);
      }).catch(error => this.toast.showError(error));
    }
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
        if (this.serviceType == PolicyComponentServiceType.MGMT) {
          const handler = (prom: Promise<any>): Promise<any> => {
            return prom.then(() => {
              setTimeout(() => {
                this.loadData(this.currentType);
              }, 1000);
            }).catch(error => {
              this.toast.showError(error);
            });
          };

          const lang = this.translate.currentLang ?? 'en';

          switch (this.currentType) {
            case MESSAGETYPES.INIT:
              return handler((this.service as ManagementService).resetCustomInitMessageTextToDefault(lang));
            case MESSAGETYPES.VERIFYPHONE:
              return handler((this.service as ManagementService).resetCustomVerifyPhoneMessageTextToDefault(lang));
            case MESSAGETYPES.VERIFYEMAIL:
              return handler((this.service as ManagementService).resetCustomVerifyEmailMessageTextToDefault(lang));
            case MESSAGETYPES.PASSWORDRESET:
              return handler((this.service as ManagementService).resetCustomPasswordResetMessageTextToDefault(lang));
            case MESSAGETYPES.DOMAINCLAIMED:
              return handler((this.service as ManagementService).resetCustomDomainClaimedMessageTextToDefault(lang));
          }

        }
      }
    });
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

  public setCurrentType(key: MESSAGETYPES): void {
    this.currentType = key;
    this.loadData(this.currentType);
  }
}

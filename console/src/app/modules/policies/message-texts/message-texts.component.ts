import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { from, Observable, of, Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import {
  GetDefaultInitMessageTextRequest as AdminGetDefaultInitMessageTextRequest,
  GetDefaultVerifyEmailMessageTextRequest as AdminGetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextRequest as AdminGetDefaultVerifyPhoneMessageTextRequest,
  SetDefaultInitMessageTextRequest,
  SetDefaultVerifyEmailMessageTextRequest,
  SetDefaultVerifyPhoneMessageTextRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  GetCustomVerifyEmailMessageTextRequest,
  GetCustomVerifyPhoneMessageTextRequest,
  GetDefaultInitMessageTextRequest,
  GetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextRequest,
  SetCustomInitMessageTextRequest,
  SetCustomVerifyEmailMessageTextRequest,
  SetCustomVerifyPhoneMessageTextRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { MessageCustomText } from 'src/app/proto/generated/zitadel/text_pb';
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

enum MESSAGETYPES {
  INIT = 'INIT',
  VERIFYPHONE = 'VP',
  VERIFYEMAIL = 'VE',
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
  public getCustomInitMessageTextMap$: Observable<{ [key: string]: string; }> = of({});

  public currentType: MESSAGETYPES = MESSAGETYPES.INIT;

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public nextLinks: CnslLinks[] = [];
  public MESSAGETYPES: any = MESSAGETYPES;

  private sub: Subscription = new Subscription();

  constructor(
    private route: ActivatedRoute,
    private injector: Injector,
    private translate: TranslateService,
  ) {
    console.log(this.MESSAGETYPES);
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
    }
  }

  public loadData(type: MESSAGETYPES) {

    console.log(this.serviceType, type);
    if (this.serviceType == PolicyComponentServiceType.MGMT) {
      const reqDefaultInit = REQUESTMAP[this.serviceType][type].getDefault;

      reqDefaultInit.setLanguage(this.translate.currentLang);
      this.getDefaultInitMessageTextMap$ = from(
        this.getDefaultValues(type, reqDefaultInit)
      );
    }

    const reqCustomInit = REQUESTMAP[this.serviceType][type].get.setLanguage(this.translate.currentLang);
    this.getCustomInitMessageTextMap$ = from(
      this.getCurrentValues(type, reqCustomInit)
    );
  }

  public updateCurrentValues(values: { [key: string]: string; }): void {
    const req = REQUESTMAP[this.serviceType][MESSAGETYPES.INIT].setFcn;
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

  public setCurrentType(key: MESSAGETYPES): void {
    this.currentType = key;
    this.loadData(this.currentType);
  }
}

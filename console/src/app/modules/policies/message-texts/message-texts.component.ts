import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Observable, of, Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { GetDefaultInitMessageTextRequest } from 'src/app/proto/generated/zitadel/admin_pb';
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

@Component({
  selector: 'app-message-texts',
  templateUrl: './message-texts.component.html',
  styleUrls: ['./message-texts.component.scss'],
})
export class MessageTextsComponent implements OnDestroy {
  public getDefaultInitMessageTextMap$: Observable<{ [key: string]: string; }> = of({});
  public getCustomInitMessageTextMap$: Observable<{ [key: string]: string; }> = of({});

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public nextLinks: CnslLinks[] = [];

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

          const req = new GetDefaultInitMessageTextRequest().setLanguage(this.translate.currentLang);
          // Object.keys((new GetDefaultInitMessageTextResponse()).toObject().customText).keys(key => {})
          this.getDefaultInitMessageTextMap$ = of(
            this.service.getDefaultInitMessageText(req).then(res => {
              if (res.customText) {
                delete res.customText.details;
                return Object.assign({}, res.customText as unknown as { [key: string]: string; });
              } else {
                return {};
              }
            })
          );
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

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }
}

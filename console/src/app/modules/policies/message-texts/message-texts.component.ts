import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { LoginMethodComponentType } from 'src/app/modules/mfa-table/mfa-table.component';
import {
  GetDefaultDomainClaimedMessageTextRequest,
  GetDefaultInitMessageTextRequest,
  GetDefaultPasswordResetMessageTextRequest,
  GetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { IDPStylingType } from 'src/app/proto/generated/zitadel/idp_pb';
import {
  GetCustomDomainClaimedMessageTextRequest,
  GetCustomInitMessageTextRequest,
  GetCustomPasswordResetMessageTextRequest,
  GetCustomVerifyEmailMessageTextRequest,
  GetCustomVerifyPhoneMessageTextRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { PasswordlessType } from 'src/app/proto/generated/zitadel/policy_pb';
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
  public LoginMethodComponentType: any = LoginMethodComponentType;
  public passwordlessTypes: Array<PasswordlessType> = [];

  public initMsg!: GetCustomInitMessageTextRequest.AsObject | GetDefaultInitMessageTextRequest.AsObject;
  public verifyEmailMsg!: GetCustomVerifyEmailMessageTextRequest.AsObject | GetDefaultVerifyEmailMessageTextRequest.AsObject;
  public verifyPhoneMsg!: GetCustomVerifyPhoneMessageTextRequest.AsObject | GetDefaultVerifyPhoneMessageTextRequest.AsObject;
  public passwordResetMsg!: GetCustomPasswordResetMessageTextRequest.AsObject | GetDefaultPasswordResetMessageTextRequest.AsObject;
  public domainClaimed!: GetCustomDomainClaimedMessageTextRequest.AsObject | GetDefaultDomainClaimedMessageTextRequest.AsObject;

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public loading: boolean = false;
  public disabled: boolean = true;

  public IDPStylingType: any = IDPStylingType;
  public nextLinks: CnslLinks[] = [];

  private sub: Subscription = new Subscription();

  constructor(
    private route: ActivatedRoute,
    // private toast: ToastService,
    private injector: Injector,
  ) {
    this.sub = this.route.data.pipe(switchMap(data => {
      this.serviceType = data.serviceType;
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          this.service = this.injector.get(ManagementService as Type<ManagementService>);
          this.passwordlessTypes = [
            PasswordlessType.PASSWORDLESS_TYPE_ALLOWED,
            PasswordlessType.PASSWORDLESS_TYPE_NOT_ALLOWED,
          ];
          this.nextLinks = [
            ORG_COMPLEXITY_LINK,
            ORG_IAM_POLICY_LINK,
            ORG_PRIVATELABEL_LINK,
          ];
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);
          this.passwordlessTypes = [
            PasswordlessType.PASSWORDLESS_TYPE_ALLOWED,
            PasswordlessType.PASSWORDLESS_TYPE_NOT_ALLOWED,
          ];
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

import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import {
    GetLoginPolicyResponse as AdminGetLoginPolicyResponse,
    UpdateLoginPolicyRequest,
    UpdateLoginPolicyResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
    AddCustomLoginPolicyRequest,
    GetLoginPolicyResponse as MgmtGetLoginPolicyResponse,
} from 'src/app/proto/generated/zitadel/management_pb';
import { LoginPolicy, PasswordlessType } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { LoginMethodComponentType } from './mfa-table/mfa-table.component';

@Component({
  selector: 'cnsl-login-policy',
  templateUrl: './login-policy.component.html',
  styleUrls: ['./login-policy.component.scss'],
})
export class LoginPolicyComponent implements OnInit {
  public LoginMethodComponentType: any = LoginMethodComponentType;
  public passwordlessTypes: Array<PasswordlessType> = [
    PasswordlessType.PASSWORDLESS_TYPE_NOT_ALLOWED,
    PasswordlessType.PASSWORDLESS_TYPE_ALLOWED,
  ];
  public loginData!: LoginPolicy.AsObject;

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public loading: boolean = false;
  public InfoSectionType: any = InfoSectionType;
  public PasswordlessType: any = PasswordlessType;

  constructor(private toast: ToastService, private injector: Injector) {}

  private fetchData(): void {
    this.getData().then((resp) => {
      if (resp.policy) {
        this.loginData = resp.policy;
        this.loading = false;
      }
    });
  }

  public ngOnInit(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        this.service = this.injector.get(ManagementService as Type<ManagementService>);
        this.passwordlessTypes = [
          PasswordlessType.PASSWORDLESS_TYPE_ALLOWED,
          PasswordlessType.PASSWORDLESS_TYPE_NOT_ALLOWED,
        ];
        break;
      case PolicyComponentServiceType.ADMIN:
        this.service = this.injector.get(AdminService as Type<AdminService>);
        this.passwordlessTypes = [
          PasswordlessType.PASSWORDLESS_TYPE_ALLOWED,
          PasswordlessType.PASSWORDLESS_TYPE_NOT_ALLOWED,
        ];
        break;
    }
    this.fetchData();
  }

  private async getData(): Promise<AdminGetLoginPolicyResponse.AsObject | MgmtGetLoginPolicyResponse.AsObject> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getLoginPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getLoginPolicy();
    }
  }

  private async updateData(): Promise<UpdateLoginPolicyResponse.AsObject> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        const mgmtreq = new AddCustomLoginPolicyRequest();
        mgmtreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
        mgmtreq.setAllowRegister(this.loginData.allowRegister);
        mgmtreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
        mgmtreq.setForceMfa(this.loginData.forceMfa);
        mgmtreq.setPasswordlessType(this.loginData.passwordlessType);
        mgmtreq.setHidePasswordReset(this.loginData.hidePasswordReset);
        // if(this.loginData.passwordCheckLifetime) {
        // mgmtreq.setPasswordCheckLifetime(this.loginData.passwordCheckLifetime);
        // }

        if ((this.loginData as LoginPolicy.AsObject).isDefault) {
          return (this.service as ManagementService).addCustomLoginPolicy(mgmtreq);
        } else {
          return (this.service as ManagementService).updateCustomLoginPolicy(mgmtreq);
        }
      case PolicyComponentServiceType.ADMIN:
        const adminreq = new UpdateLoginPolicyRequest();
        adminreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
        adminreq.setAllowRegister(this.loginData.allowRegister);
        adminreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
        adminreq.setForceMfa(this.loginData.forceMfa);
        adminreq.setPasswordlessType(this.loginData.passwordlessType);
        adminreq.setHidePasswordReset(this.loginData.hidePasswordReset);
        // adminreq.setPasswordCheckLifetime(this.loginData.passwordCheckLifetime);

        return (this.service as AdminService).updateLoginPolicy(adminreq);
    }
  }

  public savePolicy(): void {
    this.updateData()
      .then(() => {
        this.toast.showInfo('POLICY.LOGIN_POLICY.SAVED', true);
        this.loading = true;
        setTimeout(() => {
          this.fetchData();
        }, 2000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public removePolicy(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      (this.service as ManagementService)
        .resetLoginPolicyToDefault()
        .then(() => {
          this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
          this.loading = true;
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public get isDefault(): boolean {
    if (this.loginData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.loginData as LoginPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }
}

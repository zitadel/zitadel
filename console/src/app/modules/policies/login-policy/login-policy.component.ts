import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
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
import { GridPolicy, LOGIN_POLICY } from '../../policy-grid/policies';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { LoginMethodComponentType } from './mfa-table/mfa-table.component';

@Component({
  selector: 'cnsl-login-policy',
  templateUrl: './login-policy.component.html',
  styleUrls: ['./login-policy.component.scss'],
})
export class LoginPolicyComponent implements OnDestroy {
  public LoginMethodComponentType: any = LoginMethodComponentType;
  public passwordlessTypes: Array<PasswordlessType> = [];
  public loginData!: LoginPolicy.AsObject;

  private sub: Subscription = new Subscription();
  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public loading: boolean = false;
  public disabled: boolean = true;

  public currentPolicy: GridPolicy = LOGIN_POLICY;
  public InfoSectionType: any = InfoSectionType;
  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
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
          break;
        case PolicyComponentServiceType.ADMIN:
          this.service = this.injector.get(AdminService as Type<AdminService>);
          this.passwordlessTypes = [
            PasswordlessType.PASSWORDLESS_TYPE_ALLOWED,
            PasswordlessType.PASSWORDLESS_TYPE_NOT_ALLOWED,
          ];
          break;
      }

      return this.route.params;
    })).subscribe(() => {
      this.fetchData();
    });
  }

  private fetchData(): void {
    this.getData().then(resp => {
      if (resp.policy) {
        this.loginData = resp.policy;
        this.loading = false;
        this.disabled = this.isDefault;
      }
    });

  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  private async getData():
    Promise<AdminGetLoginPolicyResponse.AsObject | MgmtGetLoginPolicyResponse.AsObject> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getLoginPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getLoginPolicy();
    }
  }

  private async updateData():
    Promise<UpdateLoginPolicyResponse.AsObject> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        const mgmtreq = new AddCustomLoginPolicyRequest();
        mgmtreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
        mgmtreq.setAllowRegister(this.loginData.allowRegister);
        mgmtreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
        mgmtreq.setForceMfa(this.loginData.forceMfa);
        mgmtreq.setPasswordlessType(this.loginData.passwordlessType);
        mgmtreq.setHidePasswordReset(this.loginData.hidePasswordReset);
        mgmtreq.setIgnoreUnknownUsernames(this.loginData.ignoreUnknownUsernames);
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
        adminreq.setIgnoreUnknownUsernames(this.loginData.ignoreUnknownUsernames);

        return (this.service as AdminService).updateLoginPolicy(adminreq);
    }
  }

  public savePolicy(): void {
    this.updateData().then(() => {
      this.toast.showInfo('POLICY.LOGIN_POLICY.SAVED', true);
      this.loading = true;
      setTimeout(() => {
        this.fetchData();
      }, 2000);
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public removePolicy(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      (this.service as ManagementService).resetLoginPolicyToDefault().then(() => {
        this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
        this.loading = true;
        setTimeout(() => {
          this.fetchData();
        }, 2000);
      }).catch(error => {
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

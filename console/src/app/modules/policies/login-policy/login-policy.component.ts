import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
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
  public lifetimeForm!: FormGroup;
  constructor(private toast: ToastService, private injector: Injector, private fb: FormBuilder) {
    this.lifetimeForm = this.fb.group({
      passwordCheckLifetime: [240, [Validators.required]],
      externalLoginCheckLifetime: [12, [Validators.required]],
      mfaInitSkipLifetime: [720, [Validators.required]],
      secondFactorCheckLifetime: [12, [Validators.required]],
      multiFactorCheckLifetime: [12, [Validators.required]],
    });
  }

  private fetchData(): void {
    this.getData()
      .then((resp) => {
        console.log(resp);

        if (resp.policy) {
          this.loginData = resp.policy;
          this.loading = false;

          this.passwordCheckLifetime?.setValue(
            this.loginData.passwordCheckLifetime?.seconds ? this.loginData.passwordCheckLifetime?.seconds / 60 / 60 : 240,
          );

          this.externalLoginCheckLifetime?.setValue(
            this.loginData.externalLoginCheckLifetime?.seconds
              ? this.loginData.externalLoginCheckLifetime?.seconds / 60 / 60
              : 12,
          );

          this.mfaInitSkipLifetime?.setValue(
            this.loginData.mfaInitSkipLifetime?.seconds ? this.loginData.mfaInitSkipLifetime?.seconds / 60 / 60 : 720,
          );

          this.secondFactorCheckLifetime?.setValue(
            this.loginData.secondFactorCheckLifetime?.seconds
              ? this.loginData.secondFactorCheckLifetime?.seconds / 60 / 60
              : 12,
          );

          this.multiFactorCheckLifetime?.setValue(
            this.loginData.multiFactorCheckLifetime?.seconds
              ? this.loginData.multiFactorCheckLifetime?.seconds / 60 / 60
              : 12,
          );
        }
      })
      .catch(this.toast.showError);
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

        const pcl = new Duration().setSeconds((this.passwordCheckLifetime?.value ?? 240) * 60 * 60);
        mgmtreq.setPasswordCheckLifetime(pcl);

        const elcl = new Duration().setSeconds((this.externalLoginCheckLifetime?.value ?? 12) * 60 * 60);
        mgmtreq.setExternalLoginCheckLifetime(elcl);

        const misl = new Duration().setSeconds((this.mfaInitSkipLifetime?.value ?? 720) * 60 * 60);
        mgmtreq.setMfaInitSkipLifetime(misl);

        const sfcl = new Duration().setSeconds((this.secondFactorCheckLifetime?.value ?? 12) * 60 * 60);
        mgmtreq.setSecondFactorCheckLifetime(sfcl);

        const mficl = new Duration().setSeconds((this.multiFactorCheckLifetime?.value ?? 12) * 60 * 60);
        mgmtreq.setMultiFactorCheckLifetime(mficl);

        mgmtreq.setIgnoreUnknownUsernames(this.loginData.ignoreUnknownUsernames);
        mgmtreq.setDefaultRedirectUri(this.loginData.defaultRedirectUri);

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

        const admin_pcl = new Duration().setSeconds((this.passwordCheckLifetime?.value ?? 240) * 60 * 60);
        adminreq.setPasswordCheckLifetime(admin_pcl);

        const admin_elcl = new Duration().setSeconds((this.externalLoginCheckLifetime?.value ?? 12) * 60 * 60);
        adminreq.setExternalLoginCheckLifetime(admin_elcl);

        const admin_misl = new Duration().setSeconds((this.mfaInitSkipLifetime?.value ?? 720) * 60 * 60);
        adminreq.setMfaInitSkipLifetime(admin_misl);

        const admin_sfcl = new Duration().setSeconds((this.secondFactorCheckLifetime?.value ?? 12) * 60 * 60);
        adminreq.setSecondFactorCheckLifetime(admin_sfcl);

        const admin_mficl = new Duration().setSeconds((this.multiFactorCheckLifetime?.value ?? 12) * 60 * 60);
        adminreq.setMultiFactorCheckLifetime(admin_mficl);
        adminreq.setIgnoreUnknownUsernames(this.loginData.ignoreUnknownUsernames);
        adminreq.setDefaultRedirectUri(this.loginData.defaultRedirectUri);
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

  public get passwordCheckLifetime(): AbstractControl | null {
    return this.lifetimeForm.get('passwordCheckLifetime');
  }

  public get externalLoginCheckLifetime(): AbstractControl | null {
    return this.lifetimeForm.get('externalLoginCheckLifetime');
  }

  public get mfaInitSkipLifetime(): AbstractControl | null {
    return this.lifetimeForm.get('mfaInitSkipLifetime');
  }

  public get secondFactorCheckLifetime(): AbstractControl | null {
    return this.lifetimeForm.get('secondFactorCheckLifetime');
  }

  public get multiFactorCheckLifetime(): AbstractControl | null {
    return this.lifetimeForm.get('multiFactorCheckLifetime');
  }
}

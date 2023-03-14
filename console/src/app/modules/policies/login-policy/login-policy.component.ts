import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { take } from 'rxjs';
import {
  GetLoginPolicyResponse as AdminGetLoginPolicyResponse,
  UpdateLoginPolicyRequest,
  UpdateLoginPolicyResponse,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  AddCustomLoginPolicyRequest,
  GetLoginPolicyResponse as MgmtGetLoginPolicyResponse,
  UpdateCustomLoginPolicyRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { LoginPolicy, PasswordlessType } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../../form-field/validators/validators';

import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { LoginMethodComponentType } from './factor-table/factor-table.component';

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
  public loginData?: LoginPolicy.AsObject;

  public service!: ManagementService | AdminService;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public loading: boolean = false;
  public InfoSectionType: any = InfoSectionType;
  public PasswordlessType: any = PasswordlessType;
  public lifetimeForm: UntypedFormGroup = this.fb.group({
    passwordCheckLifetime: [{ disabled: true }, [requiredValidator]],
    externalLoginCheckLifetime: [{ disabled: true }, [requiredValidator]],
    mfaInitSkipLifetime: [{ disabled: true }, [requiredValidator]],
    secondFactorCheckLifetime: [{ disabled: true }, [requiredValidator]],
    multiFactorCheckLifetime: [{ disabled: true }, [requiredValidator]],
  });
  constructor(
    private toast: ToastService,
    private injector: Injector,
    private fb: UntypedFormBuilder,
    private authService: GrpcAuthService,
    private dialog: MatDialog,
  ) {}

  public fetchData(): void {
    this.getData()
      .then((resp) => {
        if (resp.policy) {
          this.loginData = resp.policy;
          this.loading = false;

          this.passwordCheckLifetime?.setValue(
            this.loginData.passwordCheckLifetime?.seconds ? this.loginData.passwordCheckLifetime?.seconds / 60 / 60 : 0,
          );

          this.externalLoginCheckLifetime?.setValue(
            this.loginData.externalLoginCheckLifetime?.seconds
              ? this.loginData.externalLoginCheckLifetime?.seconds / 60 / 60
              : 0,
          );

          this.mfaInitSkipLifetime?.setValue(
            this.loginData.mfaInitSkipLifetime?.seconds ? this.loginData.mfaInitSkipLifetime?.seconds / 60 / 60 : 0,
          );

          this.secondFactorCheckLifetime?.setValue(
            this.loginData.secondFactorCheckLifetime?.seconds
              ? this.loginData.secondFactorCheckLifetime?.seconds / 60 / 60
              : 0,
          );

          this.multiFactorCheckLifetime?.setValue(
            this.loginData.multiFactorCheckLifetime?.seconds
              ? this.loginData.multiFactorCheckLifetime?.seconds / 60 / 60
              : 0,
          );
        }
      })
      .catch((error) => {
        this.toast.showError(error);
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
    this.authService
      .isAllowed(
        this.serviceType === PolicyComponentServiceType.ADMIN
          ? ['iam.policy.write']
          : this.serviceType === PolicyComponentServiceType.MGMT
          ? ['policy.write']
          : [],
      )
      .pipe(take(1))
      .subscribe((allowed) => {
        if (allowed) {
          this.lifetimeForm.enable();
        }
      });
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
    if (this.loginData) {
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          if (this.isDefault) {
            const mgmtreq = new AddCustomLoginPolicyRequest();
            mgmtreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
            mgmtreq.setAllowRegister(this.loginData.allowRegister);
            mgmtreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
            mgmtreq.setForceMfa(this.loginData.forceMfa);
            mgmtreq.setPasswordlessType(this.loginData.passwordlessType);
            mgmtreq.setHidePasswordReset(this.loginData.hidePasswordReset);
            mgmtreq.setMultiFactorsList(this.loginData.multiFactorsList);
            mgmtreq.setSecondFactorsList(this.loginData.secondFactorsList);
            mgmtreq.setDisableLoginWithEmail(this.loginData.disableLoginWithEmail);
            mgmtreq.setDisableLoginWithPhone(this.loginData.disableLoginWithPhone);

            const pcl = new Duration().setSeconds((this.passwordCheckLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setPasswordCheckLifetime(pcl);

            const elcl = new Duration().setSeconds((this.externalLoginCheckLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setExternalLoginCheckLifetime(elcl);

            const misl = new Duration().setSeconds((this.mfaInitSkipLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setMfaInitSkipLifetime(misl);

            const sfcl = new Duration().setSeconds((this.secondFactorCheckLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setSecondFactorCheckLifetime(sfcl);

            const mficl = new Duration().setSeconds((this.multiFactorCheckLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setMultiFactorCheckLifetime(mficl);

            mgmtreq.setAllowDomainDiscovery(this.loginData.allowDomainDiscovery);
            mgmtreq.setIgnoreUnknownUsernames(this.loginData.ignoreUnknownUsernames);
            mgmtreq.setDefaultRedirectUri(this.loginData.defaultRedirectUri);

            return (this.service as ManagementService).addCustomLoginPolicy(mgmtreq);
          } else {
            const mgmtreq = new UpdateCustomLoginPolicyRequest();
            mgmtreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
            mgmtreq.setAllowRegister(this.loginData.allowRegister);
            mgmtreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
            mgmtreq.setForceMfa(this.loginData.forceMfa);
            mgmtreq.setPasswordlessType(this.loginData.passwordlessType);
            mgmtreq.setHidePasswordReset(this.loginData.hidePasswordReset);
            mgmtreq.setDisableLoginWithEmail(this.loginData.disableLoginWithEmail);
            mgmtreq.setDisableLoginWithPhone(this.loginData.disableLoginWithPhone);

            const pcl = new Duration().setSeconds((this.passwordCheckLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setPasswordCheckLifetime(pcl);

            const elcl = new Duration().setSeconds((this.externalLoginCheckLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setExternalLoginCheckLifetime(elcl);

            const misl = new Duration().setSeconds((this.mfaInitSkipLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setMfaInitSkipLifetime(misl);

            const sfcl = new Duration().setSeconds((this.secondFactorCheckLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setSecondFactorCheckLifetime(sfcl);

            const mficl = new Duration().setSeconds((this.multiFactorCheckLifetime?.value ?? 0) * 60 * 60);
            mgmtreq.setMultiFactorCheckLifetime(mficl);

            mgmtreq.setAllowDomainDiscovery(this.loginData.allowDomainDiscovery);
            mgmtreq.setIgnoreUnknownUsernames(this.loginData.ignoreUnknownUsernames);
            mgmtreq.setDefaultRedirectUri(this.loginData.defaultRedirectUri);

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
          adminreq.setDisableLoginWithEmail(this.loginData.disableLoginWithEmail);
          adminreq.setDisableLoginWithPhone(this.loginData.disableLoginWithPhone);

          const admin_pcl = new Duration().setSeconds((this.passwordCheckLifetime?.value ?? 0) * 60 * 60);
          adminreq.setPasswordCheckLifetime(admin_pcl);

          const admin_elcl = new Duration().setSeconds((this.externalLoginCheckLifetime?.value ?? 0) * 60 * 60);
          adminreq.setExternalLoginCheckLifetime(admin_elcl);

          const admin_misl = new Duration().setSeconds((this.mfaInitSkipLifetime?.value ?? 0) * 60 * 60);
          adminreq.setMfaInitSkipLifetime(admin_misl);

          const admin_sfcl = new Duration().setSeconds((this.secondFactorCheckLifetime?.value ?? 0) * 60 * 60);
          adminreq.setSecondFactorCheckLifetime(admin_sfcl);

          const admin_mficl = new Duration().setSeconds((this.multiFactorCheckLifetime?.value ?? 0) * 60 * 60);
          adminreq.setMultiFactorCheckLifetime(admin_mficl);
          adminreq.setAllowDomainDiscovery(this.loginData.allowDomainDiscovery);
          adminreq.setIgnoreUnknownUsernames(this.loginData.ignoreUnknownUsernames);
          adminreq.setDefaultRedirectUri(this.loginData.defaultRedirectUri);

          return (this.service as AdminService).updateLoginPolicy(adminreq);
      }
    } else {
      return Promise.reject();
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
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.RESET',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'SETTING.DIALOG.RESET.DEFAULTTITLE',
          descriptionKey: 'SETTING.DIALOG.RESET.DEFAULTDESCRIPTION',
          warnSectionKey: 'SETTING.DIALOG.RESET.LOGINPOLICY_DESCRIPTION',
        },
        width: '400px',
      });

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
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
      });
    }
  }

  public removeFactor(request: Promise<unknown>): void {
    // create policy before types can be removed
    if (this.isDefault) {
      this.updateData()
        .then(() => {
          return request;
        })
        .then(() => {
          this.toast.showInfo('MFA.TOAST.DELETED', true);
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else {
      request
        .then(() => {
          this.toast.showInfo('MFA.TOAST.DELETED', true);
          setTimeout(() => {
            this.fetchData();
          }, 2000);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public addFactor(request: Promise<unknown>): void {
    // create policy before types can be added
    const task: Promise<unknown> = this.isDefault
      ? this.updateData().then(() => {
          return request;
        })
      : request;

    task
      .then(() => {
        this.toast.showInfo('MFA.TOAST.ADDED', true);
        setTimeout(() => {
          this.fetchData();
        }, 2000);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
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

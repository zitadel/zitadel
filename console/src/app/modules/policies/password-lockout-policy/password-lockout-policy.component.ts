import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { UntypedFormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { GetLockoutPolicyResponse as AdminGetPasswordLockoutPolicyResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { GetLockoutPolicyResponse as MgmtGetPasswordLockoutPolicyResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { LockoutPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'cnsl-password-lockout-policy',
  templateUrl: './password-lockout-policy.component.html',
  styleUrls: ['./password-lockout-policy.component.scss'],
})
export class PasswordLockoutPolicyComponent implements OnInit {
  @Input() public service!: ManagementService | AdminService;
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public lockoutForm!: UntypedFormGroup;
  public lockoutData?: LockoutPolicy.AsObject;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public InfoSectionType: any = InfoSectionType;

  constructor(
    private toast: ToastService,
    private injector: Injector,
    private dialog: MatDialog,
  ) {}

  public ngOnInit(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        this.service = this.injector.get(ManagementService as Type<ManagementService>);
        break;
      case PolicyComponentServiceType.ADMIN:
        this.service = this.injector.get(AdminService as Type<AdminService>);
        break;
    }
    this.fetchData();
  }

  private fetchData(): void {
    this.getData().then((resp) => {
      if (resp.policy) {
        this.lockoutData = resp.policy;
      }
    });
  }

  private getData(): Promise<
    AdminGetPasswordLockoutPolicyResponse.AsObject | MgmtGetPasswordLockoutPolicyResponse.AsObject
  > {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getLockoutPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getLockoutPolicy();
    }
  }

  public resetPolicy(): void {
    if (this.service instanceof ManagementService) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.RESET',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'SETTING.DIALOG.RESET.DEFAULTTITLE',
          descriptionKey: 'SETTING.DIALOG.RESET.DEFAULTDESCRIPTION',
        },
        width: '400px',
      });

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          (this.service as ManagementService)
            .resetLockoutPolicyToDefault()
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    }
  }

  public incrementMaxAttempts(): void {
    if (this.lockoutData?.maxPasswordAttempts !== undefined) {
      this.lockoutData.maxPasswordAttempts++;
    }
  }

  public decrementMaxAttempts(): void {
    if (this.lockoutData?.maxPasswordAttempts && this.lockoutData?.maxPasswordAttempts > 0) {
      this.lockoutData.maxPasswordAttempts--;
    }
  }

  public savePolicy(): void {
    let promise: Promise<any>;
    if (this.lockoutData) {
      if (this.service instanceof AdminService) {
        promise = this.service
          .updateLockoutPolicy(this.lockoutData.maxPasswordAttempts)
          .then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
            this.fetchData();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      } else {
        if ((this.lockoutData as LockoutPolicy.AsObject).isDefault) {
          promise = (this.service as ManagementService)
            .addCustomLockoutPolicy(this.lockoutData.maxPasswordAttempts)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else {
          promise = (this.service as ManagementService)
            .updateCustomLockoutPolicy(this.lockoutData.maxPasswordAttempts)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      }
    }
  }

  public get isDefault(): boolean {
    if (this.lockoutData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.lockoutData as LockoutPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }
}

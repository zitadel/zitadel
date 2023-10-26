import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { GetPasswordComplexityPolicyResponse as AdminGetPasswordComplexityPolicyResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { GetPasswordComplexityPolicyResponse as MgmtGetPasswordComplexityPolicyResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'cnsl-password-complexity-policy',
  templateUrl: './password-complexity-policy.component.html',
  styleUrls: ['./password-complexity-policy.component.scss'],
})
export class PasswordComplexityPolicyComponent implements OnInit {
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  public service!: ManagementService | AdminService;

  public complexityData?: PasswordComplexityPolicy.AsObject;

  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  public loading: boolean = false;
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

  public fetchData(): void {
    this.loading = true;

    this.getData().then((data) => {
      if (data.policy) {
        this.complexityData = data.policy;
        this.loading = false;
      }
    });
  }

  private async getData(): Promise<
    MgmtGetPasswordComplexityPolicyResponse.AsObject | AdminGetPasswordComplexityPolicyResponse.AsObject
  > {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getPasswordComplexityPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getPasswordComplexityPolicy();
    }
  }

  public removePolicy(): void {
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
            .resetPasswordComplexityPolicyToDefault()
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
              setTimeout(() => {
                this.fetchData();
              }, 1000);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    }
  }

  public incrementLength(): void {
    if (this.complexityData?.minLength !== undefined && this.complexityData?.minLength <= 72) {
      this.complexityData.minLength++;
    }
  }

  public decrementLength(): void {
    if (this.complexityData?.minLength && this.complexityData?.minLength > 1) {
      this.complexityData.minLength--;
    }
  }

  public savePolicy(): void {
    if (this.complexityData) {
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          if ((this.complexityData as PasswordComplexityPolicy.AsObject).isDefault) {
            (this.service as ManagementService)
              .addCustomPasswordComplexityPolicy(
                this.complexityData.hasLowercase,
                this.complexityData.hasUppercase,
                this.complexityData.hasNumber,
                this.complexityData.hasSymbol,
                this.complexityData.minLength,
              )
              .then(() => {
                this.toast.showInfo('POLICY.TOAST.SET', true);
              })
              .catch((error) => {
                this.toast.showError(error);
              });
          } else {
            (this.service as ManagementService)
              .updateCustomPasswordComplexityPolicy(
                this.complexityData.hasLowercase,
                this.complexityData.hasUppercase,
                this.complexityData.hasNumber,
                this.complexityData.hasSymbol,
                this.complexityData.minLength,
              )
              .then(() => {
                this.toast.showInfo('POLICY.TOAST.SET', true);
              })
              .catch((error) => {
                this.toast.showError(error);
              });
          }
          break;
        case PolicyComponentServiceType.ADMIN:
          (this.service as AdminService)
            .updatePasswordComplexityPolicy(
              this.complexityData.hasLowercase,
              this.complexityData.hasUppercase,
              this.complexityData.hasNumber,
              this.complexityData.hasSymbol,
              this.complexityData.minLength,
            )
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
          break;
      }
    }
  }

  public get isDefault(): boolean {
    if (this.complexityData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.complexityData as PasswordComplexityPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }
}

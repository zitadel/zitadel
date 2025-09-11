import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { GetPasswordAgePolicyResponse as AdminGetPasswordAgePolicyResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { GetPasswordAgePolicyResponse as MgmtGetPasswordAgePolicyResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { PasswordAgePolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { requiredValidator } from '../../form-field/validators/validators';
import { Observable } from 'rxjs';
import { GrpcAuthService } from '../../../services/grpc-auth.service';
import { take } from 'rxjs/operators';

@Component({
  selector: 'cnsl-password-age-policy',
  templateUrl: './password-age-policy.component.html',
  styleUrls: ['./password-age-policy.component.scss'],
})
export class PasswordAgePolicyComponent implements OnInit {
  @Input() public service!: ManagementService | AdminService;
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

  public passwordAgeForm!: UntypedFormGroup;
  public passwordAgeData?: PasswordAgePolicy.AsObject;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public InfoSectionType: any = InfoSectionType;
  public loading: boolean = false;

  public canWrite$: Observable<boolean> = this.authService.isAllowed([
    this.serviceType === PolicyComponentServiceType.ADMIN
      ? 'iam.policy.write'
      : this.serviceType === PolicyComponentServiceType.MGMT
        ? 'policy.write'
        : '',
  ]);

  constructor(
    private authService: GrpcAuthService,
    private toast: ToastService,
    private injector: Injector,
    private dialog: MatDialog,
    private fb: UntypedFormBuilder,
  ) {
    this.passwordAgeForm = this.fb.group({
      maxAgeDays: ['', []],
      expireWarnDays: ['', []],
    });

    this.canWrite$.pipe(take(1)).subscribe((canWrite) => {
      if (canWrite) {
        this.passwordAgeForm.enable();
      } else {
        this.passwordAgeForm.disable();
      }
    });
  }

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
    this.loading = true;

    this.getData().then((resp) => {
      if (resp.policy) {
        this.passwordAgeData = resp.policy;
        this.passwordAgeForm.patchValue(this.passwordAgeData);
        this.loading = false;
      }
    });
  }

  private getData(): Promise<AdminGetPasswordAgePolicyResponse.AsObject | MgmtGetPasswordAgePolicyResponse.AsObject> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getPasswordAgePolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getPasswordAgePolicy();
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
            .resetPasswordAgePolicyToDefault()
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

  public savePolicy(): void {
    let promise: Promise<any>;
    if (this.passwordAgeData) {
      if (this.service instanceof AdminService) {
        promise = this.service
          .updatePasswordAgePolicy(this.maxAgeDays?.value ?? 0, this.expireWarnDays?.value ?? 0)
          .then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
            this.fetchData();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      } else {
        if ((this.passwordAgeData as PasswordAgePolicy.AsObject).isDefault) {
          promise = (this.service as ManagementService)
            .addCustomPasswordAgePolicy(this.maxAgeDays?.value ?? 0, this.expireWarnDays?.value ?? 0)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        } else {
          promise = (this.service as ManagementService)
            .updateCustomPasswordAgePolicy(this.maxAgeDays?.value ?? 0, this.expireWarnDays?.value ?? 0)
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
    if (this.passwordAgeData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.passwordAgeData as PasswordAgePolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }

  public get maxAgeDays(): AbstractControl | null {
    return this.passwordAgeForm.get('maxAgeDays');
  }

  public get expireWarnDays(): AbstractControl | null {
    return this.passwordAgeForm.get('expireWarnDays');
  }
}

import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import {
  AddNotificationPolicyRequest,
  GetNotificationPolicyResponse as AdminGetNotificationPolicyResponse,
  UpdateNotificationPolicyRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
  AddCustomNotificationPolicyRequest,
  GetNotificationPolicyResponse as MgmtGetNotificationPolicyResponse,
} from 'src/app/proto/generated/zitadel/management_pb';
import { NotificationPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from '../../info-section/info-section.component';
import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'cnsl-notification-policy',
  templateUrl: './notification-policy.component.html',
  styleUrls: ['./notification-policy.component.scss'],
})
export class NotificationPolicyComponent implements OnInit {
  @Input() public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
  public service!: ManagementService | AdminService;

  public notificationData?: NotificationPolicy.AsObject = { isDefault: false, passwordChange: false };

  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  public loading: boolean = false;
  public InfoSectionType: any = InfoSectionType;

  public isDefault: boolean = false;
  private hasNotificationPolicy: boolean = false;
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

    this.getData()
      .then((data) => {
        if (data.policy) {
          this.hasNotificationPolicy = true;
          this.notificationData = data.policy;
          this.isDefault = data.policy.isDefault;
          this.loading = false;
        }
      })
      .catch((error) => {
        this.loading = false;
        if (error && error.code === 5) {
          console.log(error);
          this.hasNotificationPolicy = false;
        }
      });
  }

  private async getData(): Promise<
    MgmtGetNotificationPolicyResponse.AsObject | AdminGetNotificationPolicyResponse.AsObject
  > {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return (this.service as ManagementService).getNotificationPolicy();
      case PolicyComponentServiceType.ADMIN:
        return (this.service as AdminService).getNotificationPolicy();
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
            .resetNotificationPolicyToDefault()
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
              this.isDefault = true;
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

  public savePolicy(): void {
    if (this.notificationData) {
      switch (this.serviceType) {
        case PolicyComponentServiceType.MGMT:
          if ((this.notificationData as NotificationPolicy.AsObject).isDefault) {
            const req = new AddCustomNotificationPolicyRequest();
            req.setPasswordChange(this.notificationData.passwordChange);
            (this.service as ManagementService)
              .addCustomNotificationPolicy(req)
              .then(() => {
                this.isDefault = false;
                this.toast.showInfo('POLICY.TOAST.SET', true);
              })
              .catch((error) => {
                this.toast.showError(error);
              });
          } else {
            const req = new UpdateNotificationPolicyRequest();
            req.setPasswordChange(this.notificationData.passwordChange);
            (this.service as ManagementService)
              .updateCustomNotificationPolicy(req)
              .then(() => {
                this.isDefault = false;
                this.toast.showInfo('POLICY.TOAST.SET', true);
              })
              .catch((error) => {
                this.toast.showError(error);
              });
          }
          break;
        case PolicyComponentServiceType.ADMIN:
          if (this.hasNotificationPolicy) {
            const req = new UpdateNotificationPolicyRequest();
            req.setPasswordChange(this.notificationData.passwordChange);
            (this.service as AdminService)
              .updateNotificationPolicy(req)
              .then(() => {
                this.isDefault = false;
                this.toast.showInfo('POLICY.TOAST.SET', true);
              })
              .catch((error) => {
                this.toast.showError(error);
              });
          } else {
            const req = new AddNotificationPolicyRequest();
            req.setPasswordChange(this.notificationData.passwordChange);
            (this.service as AdminService)
              .addNotificationPolicy(req)
              .then(() => {
                this.isDefault = false;
                this.hasNotificationPolicy = true;
                this.toast.showInfo('POLICY.TOAST.SET', true);
              })
              .catch((error) => {
                this.toast.showError(error);
              });
          }
          break;
      }
    }
  }
}

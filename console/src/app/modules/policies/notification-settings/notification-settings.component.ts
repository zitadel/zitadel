import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { SetDefaultLanguageResponse, UpdateSMTPConfigRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { SMTPConfig } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'cnsl-notification-settings',
  templateUrl: './notification-settings.component.html',
  styleUrls: ['./notification-settings.component.scss'],
})
export class NotificationSettingsComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;
  public service!: ManagementService | AdminService;

  public smtpConfig!: SMTPConfig.AsObject;

  public loading: boolean = false;
  constructor(private injector: Injector, private toast: ToastService) {}

  ngOnInit(): void {
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
    if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      (this.service as AdminService)
        .getSMTPConfig()
        .then((smtpConfig) => {
          if (smtpConfig.smtpConfig) {
            this.smtpConfig = smtpConfig.smtpConfig;
          }
        })
        .catch((error) => {
          if (error && error.code === 5) {
            console.log(error);
          }
        });
    }
  }

  private updateData(): Promise<SetDefaultLanguageResponse.AsObject> | void {
    const req = new UpdateSMTPConfigRequest();
    req.setHost(this.smtpConfig.host);
    req.setSenderAddress(this.smtpConfig.senderAddress);
    req.setSenderName(this.smtpConfig.senderName);
    req.setTls(this.smtpConfig.tls);
    req.setUser(this.smtpConfig.user);

    return (this.service as AdminService).updateSMTPConfig(req);
  }

  public savePolicy(): void {
    const prom = this.updateData();
    if (prom) {
      prom
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
}

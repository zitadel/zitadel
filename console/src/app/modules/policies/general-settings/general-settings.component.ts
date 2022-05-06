import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { SetDefaultLanguageResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'cnsl-general-settings',
  templateUrl: './general-settings.component.html',
  styleUrls: ['./general-settings.component.scss'],
})
export class GeneralSettingsComponent implements OnInit {
  @Input() public serviceType!: PolicyComponentServiceType;
  public service!: ManagementService | AdminService;

  public defaultLanguage: string = 'en';
  public defaultLanguageOptions: string[] = [];

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
      (this.service as AdminService).getDefaultLanguage().then((langResp) => {
        this.defaultLanguage = langResp.language;
      });
      (this.service as AdminService).getSupportedLanguages().then((supportedResp) => {
        this.defaultLanguageOptions = supportedResp.languagesList;
      });
    }
  }

  private updateData(): Promise<SetDefaultLanguageResponse.AsObject> | void {
    if (this.serviceType === PolicyComponentServiceType.ADMIN) {
      return (this.service as AdminService).setDefaultLanguage(this.defaultLanguage);
    } else {
      return;
    }
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

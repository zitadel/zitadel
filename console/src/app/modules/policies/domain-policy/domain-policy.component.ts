import { Component, DestroyRef, Injector, Input, OnInit, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import {
  AddCustomDomainPolicyRequest,
  GetCustomOrgIAMPolicyResponse,
  UpdateDomainPolicyRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { GetOrgIAMPolicyResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { DomainPolicy, OrgIAMPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { GrpcAuthService } from '../../../services/grpc-auth.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'cnsl-domain-policy',
  templateUrl: './domain-policy.component.html',
  styleUrls: ['./domain-policy.component.scss'],
  standalone: false,
})
export class DomainPolicyComponent implements OnInit {
  private managementService!: ManagementService;
  @Input() public serviceType!: PolicyComponentServiceType;

  public domainData!: DomainPolicy.AsObject;

  public loading: boolean = false;

  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  constructor(
    private dialog: MatDialog,
    private toast: ToastService,
    private injector: Injector,
    private adminService: AdminService,
    private readonly authService: GrpcAuthService,
    private readonly destroyRef: DestroyRef,
  ) {}

  ngOnInit(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      this.managementService = this.injector.get(ManagementService as Type<ManagementService>);
    }
    this.fetchData().then();
  }

  public async fetchData(): Promise<void> {
    this.loading = true;
    try {
      const resp = await this.getData();
      if (resp?.policy) {
        this.domainData = resp.policy;
      }
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.loading = false;
    }
  }

  private async getData(): Promise<GetCustomOrgIAMPolicyResponse.AsObject | GetOrgIAMPolicyResponse.AsObject | any> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return this.managementService.getDomainPolicy();
      case PolicyComponentServiceType.ADMIN:
        return this.adminService.getDomainPolicy();
      default:
        return Promise.reject();
    }
  }

  public async savePolicy(): Promise<void> {
    const org = await this.authService.getActiveOrg();
    if (!org) {
      console.log('No active organization found. Cannot save domain policy.');
      return;
    }
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if ((this.domainData as OrgIAMPolicy.AsObject).isDefault) {
          const req = new AddCustomDomainPolicyRequest();
          req.setOrgId(org.id);
          req.setUserLoginMustBeDomain(this.domainData.userLoginMustBeDomain);
          req.setValidateOrgDomains(this.domainData.validateOrgDomains);
          req.setSmtpSenderAddressMatchesInstanceDomain(this.domainData.smtpSenderAddressMatchesInstanceDomain);

          this.adminService
            .addCustomDomainPolicy(req)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
          break;
        } else {
          const req = new AddCustomDomainPolicyRequest();
          req.setOrgId(org.id);
          req.setUserLoginMustBeDomain(this.domainData.userLoginMustBeDomain);
          req.setValidateOrgDomains(this.domainData.validateOrgDomains);
          req.setSmtpSenderAddressMatchesInstanceDomain(this.domainData.smtpSenderAddressMatchesInstanceDomain);

          this.adminService
            .updateCustomDomainPolicy(req)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);
            })
            .catch((error) => {
              this.toast.showError(error);
            });
          break;
        }
      case PolicyComponentServiceType.ADMIN:
        const req = new UpdateDomainPolicyRequest();
        req.setUserLoginMustBeDomain(this.domainData.userLoginMustBeDomain);
        req.setValidateOrgDomains(this.domainData.validateOrgDomains);
        req.setSmtpSenderAddressMatchesInstanceDomain(this.domainData.smtpSenderAddressMatchesInstanceDomain);

        this.adminService
          .updateDomainPolicy(req)
          .then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
        break;
    }
  }

  public removePolicy(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.RESET',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'SETTING.DIALOG.RESET.DEFAULTTITLE',
          descriptionKey: 'SETTING.DIALOG.RESET.DEFAULTDESCRIPTION',
        },
        width: '400px',
      });

      dialogRef
        .afterClosed()
        .pipe(takeUntilDestroyed(this.destroyRef))
        .subscribe(async (resp) => {
          if (!resp) {
            return;
          }

          try {
            const org = await this.authService.getActiveOrg();
            if (!org) {
              console.log('No active organization found. Cannot reset domain policy.');
              return;
            }
            await this.adminService.resetCustomDomainPolicyToDefault(org.id);
            this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
            await new Promise((res) => setTimeout(res, 1000));
            await this.fetchData();
          } catch (error) {
            this.toast.showError(error);
          }
        });
    }
  }

  public get isDefault(): boolean {
    if (this.domainData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.domainData as OrgIAMPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }
}

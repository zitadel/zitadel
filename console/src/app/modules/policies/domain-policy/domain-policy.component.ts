import { Component, Injector, Input, OnDestroy, OnInit, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Subscription } from 'rxjs';
import {
  AddCustomDomainPolicyRequest,
  GetCustomOrgIAMPolicyResponse,
  UpdateDomainPolicyRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import { GetOrgIAMPolicyResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { DomainPolicy, OrgIAMPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

import { WarnDialogComponent } from '../../warn-dialog/warn-dialog.component';
import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { NewOrganizationService } from '../../../services/new-organization.service';

@Component({
  selector: 'cnsl-domain-policy',
  templateUrl: './domain-policy.component.html',
  styleUrls: ['./domain-policy.component.scss'],
})
export class DomainPolicyComponent implements OnInit, OnDestroy {
  private managementService!: ManagementService;
  @Input() public serviceType!: PolicyComponentServiceType;

  public domainData!: DomainPolicy.AsObject;

  public loading: boolean = false;
  private sub: Subscription = new Subscription();
  private orgId = this.newOrganizationService.getOrgId();

  public PolicyComponentServiceType: any = PolicyComponentServiceType;

  constructor(
    private dialog: MatDialog,
    private toast: ToastService,
    private injector: Injector,
    private adminService: AdminService,
    private readonly newOrganizationService: NewOrganizationService,
  ) {}

  ngOnInit(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      this.managementService = this.injector.get(ManagementService as Type<ManagementService>);
    }
    this.fetchData();
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public fetchData(): void {
    this.loading = true;
    this.getData()
      .then((resp) => {
        this.loading = false;
        if (resp?.policy) {
          this.domainData = resp.policy;
        }
      })
      .catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
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

  public savePolicy(): void {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if ((this.domainData as OrgIAMPolicy.AsObject).isDefault) {
          const req = new AddCustomDomainPolicyRequest();
          req.setOrgId(this.orgId());
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
          req.setOrgId(this.orgId());
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

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.adminService
            .resetCustomDomainPolicyToDefault(this.orgId())
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

  public get isDefault(): boolean {
    if (this.domainData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.domainData as OrgIAMPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }
}

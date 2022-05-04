import { Component, Injector, Input, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { GetCustomOrgIAMPolicyResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { GetOrgIAMPolicyResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { OrgIAMPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageKey, StorageLocation, StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

import { GridPolicy, IAM_POLICY } from '../../policy-grid/policies';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
  selector: 'cnsl-org-iam-policy',
  templateUrl: './org-iam-policy.component.html',
  styleUrls: ['./org-iam-policy.component.scss'],
})
export class OrgIamPolicyComponent implements OnDestroy {
  private managementService!: ManagementService;
  @Input() public serviceType!: PolicyComponentServiceType;

  public iamData!: OrgIAMPolicy.AsObject;

  private sub: Subscription = new Subscription();
  private org!: Org.AsObject;

  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public currentPolicy: GridPolicy = IAM_POLICY;
  public orgName: string = '';

  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
    private storage: StorageService,
    private injector: Injector,
    private adminService: AdminService,
    private storageService: StorageService,
    breadcrumbService: BreadcrumbService,
  ) {
    const temporg = this.storage.getItem(StorageKey.organization, StorageLocation.session) as Org.AsObject;
    if (temporg) {
      this.org = temporg;
    }
    this.sub = this.route.data
      .pipe(
        switchMap((data) => {
          this.serviceType = data.serviceType;
          if (this.serviceType === PolicyComponentServiceType.MGMT) {
            const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
            if (org && org.id) {
              this.orgName = org.name;
            }
            this.managementService = this.injector.get(ManagementService as Type<ManagementService>);

            const iambread = new Breadcrumb({
              type: BreadcrumbType.INSTANCE,
              name: 'Instance',
              routerLink: ['/instance'],
            });
            const bread: Breadcrumb = {
              type: BreadcrumbType.ORG,
              routerLink: ['/org'],
            };
            breadcrumbService.setBreadcrumb([iambread, bread]);
          }
          return this.route.params;
        }),
      )
      .subscribe((_) => {
        this.fetchData();
      });
  }

  public ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  public fetchData(): void {
    this.getData().then((resp) => {
      if (resp?.policy) {
        this.iamData = resp.policy;
      }
    });
  }

  private async getData(): Promise<GetCustomOrgIAMPolicyResponse.AsObject | GetOrgIAMPolicyResponse.AsObject | any> {
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        return this.managementService.getOrgIAMPolicy();
      case PolicyComponentServiceType.ADMIN:
        if (this.org?.id) {
          return this.adminService.getCustomOrgIAMPolicy(this.org.id);
        }
        break;
      default:
        return Promise.reject();
    }
  }

  public savePolicy(): void {
    console.log(this.iamData);
    switch (this.serviceType) {
      case PolicyComponentServiceType.MGMT:
        if ((this.iamData as OrgIAMPolicy.AsObject).isDefault) {
          this.adminService
            .addCustomOrgIAMPolicy(this.org.id, this.iamData.userLoginMustBeDomain)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
          break;
        } else {
          this.adminService
            .updateCustomOrgIAMPolicy(this.org.id, this.iamData.userLoginMustBeDomain)
            .then(() => {
              this.toast.showInfo('POLICY.TOAST.SET', true);
              this.fetchData();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
          break;
        }
      case PolicyComponentServiceType.ADMIN:
        // update Default org iam policy?
        this.adminService
          .updateOrgIAMPolicy(this.iamData.userLoginMustBeDomain)
          .then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
            this.fetchData();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
        break;
    }
  }

  public removePolicy(): void {
    if (this.serviceType === PolicyComponentServiceType.MGMT) {
      this.adminService
        .resetCustomOrgIAMPolicyToDefault(this.org.id)
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
  }

  public get isDefault(): boolean {
    if (this.iamData && this.serviceType === PolicyComponentServiceType.MGMT) {
      return (this.iamData as OrgIAMPolicy.AsObject).isDefault;
    } else {
      return false;
    }
  }
}

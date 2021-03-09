import { Component, Injector, Input, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { GetCustomOrgIAMPolicyResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { GetOrgIAMPolicyResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { OrgIAMPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

import { CnslLinks } from '../../links/links.component';
import {
    IAM_COMPLEXITY_LINK,
    IAM_LABEL_LINK,
    IAM_LOGIN_POLICY_LINK,
    ORG_COMPLEXITY_LINK,
    ORG_LOGIN_POLICY_LINK,
} from '../../policy-grid/policy-links';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
    selector: 'app-org-iam-policy',
    templateUrl: './org-iam-policy.component.html',
    styleUrls: ['./org-iam-policy.component.scss'],
})
export class OrgIamPolicyComponent implements OnDestroy {
    @Input() service!: AdminService;
    private managementService!: ManagementService;
    public serviceType!: PolicyComponentServiceType;

    public iamData!: OrgIAMPolicy.AsObject;

    private sub: Subscription = new Subscription();
    private org!: Org.AsObject;

    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    public nextLinks: Array<CnslLinks> = [];
    constructor(
        private route: ActivatedRoute,
        private toast: ToastService,
        private sessionStorage: StorageService,
        private injector: Injector,
        private adminService: AdminService,
    ) {
        const temporg = this.sessionStorage.getItem('organization') as Org.AsObject;
        if (temporg) {
            this.org = temporg;
        }
        this.sub = this.route.data.pipe(switchMap(data => {
            this.serviceType = data.serviceType;
            if (this.serviceType === PolicyComponentServiceType.MGMT) {
                this.managementService = this.injector.get(ManagementService as Type<ManagementService>);
                this.nextLinks = [
                    ORG_COMPLEXITY_LINK,
                    ORG_LOGIN_POLICY_LINK,
                ];
            } else {
                this.nextLinks = [
                    IAM_COMPLEXITY_LINK,
                    IAM_LOGIN_POLICY_LINK,
                    IAM_LABEL_LINK,
                ];
            }
            return this.route.params;
        })).subscribe(_ => {
            this.fetchData();
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    public fetchData(): void {
        this.getData().then(resp => {
            if (resp?.policy) {
                this.iamData = resp.policy;
            }
        });
    }

    private async getData(): Promise<GetCustomOrgIAMPolicyResponse.AsObject | GetOrgIAMPolicyResponse.AsObject | undefined> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return this.managementService.getOrgIAMPolicy();
            case PolicyComponentServiceType.ADMIN:
                if (this.org?.id) {
                    return this.adminService.getCustomOrgIAMPolicy(this.org.id);
                }
                break;
        }
    }

    public savePolicy(): void {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                if ((this.iamData as OrgIAMPolicy.AsObject).isDefault) {
                    this.adminService.addCustomOrgIAMPolicy(
                        this.org.id,
                        this.iamData.userLoginMustBeDomain,
                    ).then(() => {
                        this.toast.showInfo('POLICY.TOAST.SET', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                    break;
                } else {
                    this.adminService.updateCustomOrgIAMPolicy(
                        this.org.id,
                        this.iamData.userLoginMustBeDomain,
                    ).then(() => {
                        this.toast.showInfo('POLICY.TOAST.SET', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                    break;
                }
            case PolicyComponentServiceType.ADMIN:
                // update Default org iam policy?
                this.adminService.updateOrgIAMPolicy(
                    this.iamData.userLoginMustBeDomain,
                ).then(() => {
                    this.toast.showInfo('POLICY.TOAST.SET', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
                break;
        }
    }

    public removePolicy(): void {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            this.adminService.resetCustomOrgIAMPolicyToDefault(this.org.id).then(() => {
                this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
                setTimeout(() => {
                    this.fetchData();
                }, 1000);
            }).catch(error => {
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

import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { OrgIamPolicyView as AdminOrgIamPolicyView } from 'src/app/proto/generated/admin_pb';
import { OrgIamPolicyView as MgmtOrgIamPolicyView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentAction } from '../policy-component-action.enum';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
    selector: 'app-password-iam-policy',
    templateUrl: './password-iam-policy.component.html',
    styleUrls: ['./password-iam-policy.component.scss'],
})
export class PasswordIamPolicyComponent implements OnDestroy {
    public title: string = '';
    public desc: string = '';
    private managementService!: ManagementService;
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;

    public PolicyComponentAction: any = PolicyComponentAction;

    public iamData!: AdminOrgIamPolicyView.AsObject | MgmtOrgIamPolicyView.AsObject;

    private sub: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private toast: ToastService,
        private sessionStorage: StorageService,
        private injector: Injector,
        private adminService: AdminService,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.serviceType = data.serviceType;

            if (this.serviceType === PolicyComponentServiceType.MGMT) {
                this.managementService = this.injector.get(ManagementService as Type<ManagementService>);
            }
            return this.route.params;
        })).subscribe(_ => {
            this.title = 'ORG.POLICY.IAM_POLICY.TITLECREATE';
            this.desc = 'ORG.POLICY.IAM_POLICY.DESCRIPTIONCREATE';

            this.getData().then(data => {
                if (data) {
                    this.iamData = data.toObject();
                }
            });
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData(): Promise<AdminOrgIamPolicyView | MgmtOrgIamPolicyView | undefined> {
        this.title = 'ORG.POLICY.IAM_POLICY.TITLECREATE';
        this.desc = 'ORG.POLICY.IAM_POLICY.DESCRIPTIONCREATE';

        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return this.managementService.GetMyOrgIamPolicy();
            case PolicyComponentServiceType.ADMIN:
                const orgId = this.sessionStorage.getItem('organization');
                if (orgId) {
                    return this.adminService.GetOrgIamPolicy(orgId);
                }
                break;
        }
    }

    public savePolicy(): void {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                if ((this.iamData as MgmtOrgIamPolicyView.AsObject)) {
                    this.adminService.CreateOrgIamPolicy(
                        '',
                        this.iamData.userLoginMustBeDomain,
                    ).then(() => {
                        this.router.navigate(['org']);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                    break;
                } else {
                    this.adminService.UpdateOrgIamPolicy(
                        '',
                        this.iamData.userLoginMustBeDomain,
                    ).then(() => {
                        this.router.navigate(['org']);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                    break;
                }
            case PolicyComponentServiceType.ADMIN:
                this.adminService.UpdateOrgIamPolicy(
                    '',
                    this.iamData.userLoginMustBeDomain,
                ).then(() => {
                    this.router.navigate(['org']);
                }).catch(error => {
                    this.toast.showError(error);
                });
                break;
        }
    }
}

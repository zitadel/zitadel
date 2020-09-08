import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import {
    OrgIamPolicy,
    PasswordAgePolicy,
    PasswordComplexityPolicy,
    PasswordLockoutPolicy,
} from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentAction } from '../policy-component-action.enum';

@Component({
    selector: 'app-password-iam-policy',
    templateUrl: './password-iam-policy.component.html',
    styleUrls: ['./password-iam-policy.component.scss'],
})
export class PasswordIamPolicyComponent implements OnDestroy {
    public title: string = '';
    public desc: string = '';

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    public PolicyComponentAction: any = PolicyComponentAction;

    public iamData!: OrgIamPolicy.AsObject;

    private sub: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private mgmtService: ManagementService,
        private adminService: AdminService,
        private router: Router,
        private toast: ToastService,
        private sessionStorage: StorageService,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.componentAction = data.action;
            console.log(data.action);
            return this.route.params;
        })).subscribe(params => {
            this.title = 'ORG.POLICY.IAM_POLICY.TITLECREATE';
            this.desc = 'ORG.POLICY.IAM_POLICY.DESCRIPTIONCREATE';

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData(params).then(data => {
                    if (data) {
                        this.iamData = data.toObject() as OrgIamPolicy.AsObject;
                    }
                });
            }
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData(params: any):
        Promise<PasswordLockoutPolicy | PasswordAgePolicy | PasswordComplexityPolicy | OrgIamPolicy | undefined> {

        this.title = 'ORG.POLICY.IAM_POLICY.TITLECREATE';
        this.desc = 'ORG.POLICY.IAM_POLICY.DESCRIPTIONCREATE';
        return this.mgmtService.GetMyOrgIamPolicy();
    }

    public savePolicy(): void {
        if (this.componentAction === PolicyComponentAction.CREATE) {
            const orgId = this.sessionStorage.getItem('organization');
            if (orgId) {
                this.adminService.CreateOrgIamPolicy(
                    orgId,
                    this.iamData.description,
                    this.iamData.userLoginMustBeDomain,
                ).then(() => {
                    this.router.navigate(['org']);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        } else if (this.componentAction === PolicyComponentAction.MODIFY) {
            const orgId = this.sessionStorage.getItem('organization');
            if (orgId) {
                this.adminService.UpdateOrgIamPolicy(
                    orgId,
                    this.iamData.description,
                    this.iamData.userLoginMustBeDomain,
                ).then(() => {
                    this.router.navigate(['org']);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        }
    }
}

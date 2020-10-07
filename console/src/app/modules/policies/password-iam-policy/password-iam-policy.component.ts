import { Component, Input, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { OrgIamPolicyView as AdminOrgIamPolicyView } from 'src/app/proto/generated/admin_pb';
import { OrgIamPolicyView as MgmtOrgIamPolicyView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentAction } from '../policy-component-action.enum';

@Component({
    selector: 'app-password-iam-policy',
    templateUrl: './password-iam-policy.component.html',
    styleUrls: ['./password-iam-policy.component.scss'],
})
export class PasswordIamPolicyComponent implements OnDestroy {
    @Input() service!: AdminService;
    public title: string = '';
    public desc: string = '';

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    public PolicyComponentAction: any = PolicyComponentAction;

    public iamData!: AdminOrgIamPolicyView.AsObject | MgmtOrgIamPolicyView.AsObject;

    private sub: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private toast: ToastService,
        private sessionStorage: StorageService,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.componentAction = data.action;
            console.log(data.action);
            return this.route.params;
        })).subscribe(_ => {
            this.title = 'ORG.POLICY.IAM_POLICY.TITLECREATE';
            this.desc = 'ORG.POLICY.IAM_POLICY.DESCRIPTIONCREATE';

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData().then(data => {
                    if (data) {
                        this.iamData = data.toObject();
                    }
                });
            }
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData(): Promise<AdminOrgIamPolicyView | MgmtOrgIamPolicyView | undefined> {
        this.title = 'ORG.POLICY.IAM_POLICY.TITLECREATE';
        this.desc = 'ORG.POLICY.IAM_POLICY.DESCRIPTIONCREATE';
        const orgId = this.sessionStorage.getItem('organization');
        if (orgId) {
            return this.service.GetOrgIamPolicy(orgId);
        }
    }

    public savePolicy(): void {
        const orgId = this.sessionStorage.getItem('organization');
        if (this.componentAction === PolicyComponentAction.CREATE && orgId) {
            this.service.CreateOrgIamPolicy(
                orgId,
                this.iamData.userLoginMustBeDomain,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (this.componentAction === PolicyComponentAction.MODIFY && orgId) {
            if (this.service instanceof AdminService) {
                this.service.UpdateOrgIamPolicy(
                    orgId,
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

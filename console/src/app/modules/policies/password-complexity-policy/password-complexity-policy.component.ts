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
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentAction } from '../policy-component-action.enum';

@Component({
    selector: 'app-password-policy',
    templateUrl: './password-complexity-policy.component.html',
    styleUrls: ['./password-complexity-policy.component.scss'],
})
export class PasswordComplexityPolicyComponent implements OnDestroy {
    public title: string = '';
    public desc: string = '';

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    public PolicyComponentAction: any = PolicyComponentAction;

    public complexityData!: PasswordComplexityPolicy.AsObject;

    private sub: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private mgmtService: ManagementService,
        private router: Router,
        private toast: ToastService,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.componentAction = data.action;
            return this.route.params;
        })).subscribe(params => {
            this.title = 'ORG.POLICY.PWD_COMPLEXITY.TITLECREATE';
            this.desc = 'ORG.POLICY.PWD_COMPLEXITY.DESCRIPTIONCREATE';

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData(params).then(data => {
                    if (data) {
                        this.complexityData = data.toObject() as PasswordComplexityPolicy.AsObject;
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
        this.title = 'ORG.POLICY.PWD_COMPLEXITY.TITLE';
        this.desc = 'ORG.POLICY.PWD_COMPLEXITY.DESCRIPTION';
        return this.mgmtService.GetPasswordComplexityPolicy();
    }

    public deletePolicy(): void {
        this.mgmtService.DeletePasswordComplexityPolicy(this.complexityData.id).then(() => {
            this.toast.showInfo('Successfully deleted');
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public incrementLength(): void {
        if (this.complexityData?.minLength !== undefined && this.complexityData?.minLength <= 72) {
            this.complexityData.minLength++;
        }
    }

    public decrementLength(): void {
        if (this.complexityData?.minLength && this.complexityData?.minLength > 1) {
            this.complexityData.minLength--;
        }
    }

    public savePolicy(): void {
        if (this.componentAction === PolicyComponentAction.CREATE) {

            this.mgmtService.CreatePasswordComplexityPolicy(
                this.complexityData.description,
                this.complexityData.hasLowercase,
                this.complexityData.hasUppercase,
                this.complexityData.hasNumber,
                this.complexityData.hasSymbol,
                this.complexityData.minLength,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });

        } else if (this.componentAction === PolicyComponentAction.MODIFY) {

            this.mgmtService.UpdatePasswordComplexityPolicy(
                this.complexityData.description,
                this.complexityData.hasLowercase,
                this.complexityData.hasUppercase,
                this.complexityData.hasNumber,
                this.complexityData.hasSymbol,
                this.complexityData.minLength,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }
}

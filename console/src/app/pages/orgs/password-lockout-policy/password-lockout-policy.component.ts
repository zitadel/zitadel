import { Component, OnDestroy } from '@angular/core';
import { FormGroup } from '@angular/forms';
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

export enum PolicyComponentAction {
    CREATE = 'create',
    MODIFY = 'modify',
}

@Component({
    selector: 'app-password-lockout-policy',
    templateUrl: './password-lockout-policy.component.html',
    styleUrls: ['./password-lockout-policy.component.scss'],
})
export class PasswordLockoutPolicyComponent implements OnDestroy {
    public title: string = '';
    public desc: string = '';

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    public PolicyComponentAction: any = PolicyComponentAction;

    public lockoutForm!: FormGroup;
    public lockoutData!: PasswordLockoutPolicy.AsObject;
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
            this.title = 'ORG.POLICY.PWD_LOCKOUT.TITLECREATE';
            this.desc = 'ORG.POLICY.PWD_LOCKOUT.DESCRIPTIONCREATE';

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData(params).then(data => {
                    if (data) {
                        this.lockoutData = data.toObject() as PasswordLockoutPolicy.AsObject;
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

        this.title = 'ORG.POLICY.PWD_LOCKOUT.TITLE';
        this.desc = 'ORG.POLICY.PWD_LOCKOUT.DESCRIPTION';
        return this.mgmtService.GetPasswordLockoutPolicy();
    }

    public deletePolicy(): void {
        this.mgmtService.DeletePasswordLockoutPolicy(this.lockoutData.id).then(() => {
            this.toast.showInfo('Successfully deleted');
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public incrementMaxAttempts(): void {
        if (this.lockoutData?.maxAttempts !== undefined) {
            this.lockoutData.maxAttempts++;
        }
    }

    public decrementMaxAttempts(): void {
        if (this.lockoutData?.maxAttempts && this.lockoutData?.maxAttempts > 0) {
            this.lockoutData.maxAttempts--;
        }
    }

    public savePolicy(): void {
        if (this.componentAction === PolicyComponentAction.CREATE) {
            this.mgmtService.CreatePasswordLockoutPolicy(
                this.lockoutData.description,
                this.lockoutData.maxAttempts,
                this.lockoutData.showLockOutFailures,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (this.componentAction === PolicyComponentAction.MODIFY) {

            this.mgmtService.UpdatePasswordLockoutPolicy(
                this.lockoutData.description,
                this.lockoutData.maxAttempts,
                this.lockoutData.showLockOutFailures,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }
}

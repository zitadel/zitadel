import { Component, Injector, Input, OnDestroy, Type } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { DefaultPasswordLockoutPolicyView } from 'src/app/proto/generated/admin_pb';
import { PasswordLockoutPolicy, PasswordLockoutPolicyView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentAction } from '../policy-component-action.enum';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
    selector: 'app-password-lockout-policy',
    templateUrl: './password-lockout-policy.component.html',
    styleUrls: ['./password-lockout-policy.component.scss'],
})
export class PasswordLockoutPolicyComponent implements OnDestroy {
    @Input() public service!: ManagementService | AdminService;
    public title: string = '';
    public desc: string = '';

    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    public PolicyComponentAction: any = PolicyComponentAction;

    public lockoutForm!: FormGroup;
    public lockoutData!: PasswordLockoutPolicy.AsObject;
    private sub: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private toast: ToastService,
        private injector: Injector,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.serviceType = data.serviceType;
            this.componentAction = data.action;

            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    this.service = this.injector.get(ManagementService as Type<ManagementService>);
                    break;
                case PolicyComponentServiceType.ADMIN:
                    this.service = this.injector.get(AdminService as Type<AdminService>);
                    break;
            }

            return this.route.params;
        })).subscribe(() => {
            this.title = 'ORG.POLICY.PWD_LOCKOUT.TITLECREATE';
            this.desc = 'ORG.POLICY.PWD_LOCKOUT.DESCRIPTIONCREATE';

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData().then(data => {
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

    private getData(): Promise<PasswordLockoutPolicyView | DefaultPasswordLockoutPolicyView> {

        this.title = 'ORG.POLICY.PWD_LOCKOUT.TITLE';
        this.desc = 'ORG.POLICY.PWD_LOCKOUT.DESCRIPTION';
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).GetPasswordLockoutPolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultPasswordLockoutPolicy();
        }
    }

    public deletePolicy(): void {
        if (this.service instanceof ManagementService) {
            this.service.RemovePasswordLockoutPolicy().then(() => {
                this.toast.showInfo('Successfully deleted');
            }).catch(error => {
                this.toast.showError(error);
            });
        }
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
        let promise: Promise<any>;
        if (this.service instanceof AdminService) {
            promise = this.service.UpdateDefaultPasswordLockoutPolicy(
                this.lockoutData.maxAttempts,
                this.lockoutData.showLockoutFailure,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });
        } else {
            if (this.componentAction === PolicyComponentAction.CREATE) {
                promise = this.service.CreatePasswordLockoutPolicy(
                    this.lockoutData.maxAttempts,
                    this.lockoutData.showLockoutFailure,
                ).then(() => {
                    this.router.navigate(['org']);
                }).catch(error => {
                    this.toast.showError(error);
                });
            } else if (this.componentAction === PolicyComponentAction.MODIFY) {
                promise = this.service.UpdatePasswordLockoutPolicy(
                    this.lockoutData.maxAttempts,
                    this.lockoutData.showLockoutFailure,
                ).then(() => {
                    this.router.navigate(['org']);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        }
    }
}

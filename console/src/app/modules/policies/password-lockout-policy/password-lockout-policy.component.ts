import { Component, Injector, Input, OnDestroy, Type } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { DefaultPasswordLockoutPolicyView } from 'src/app/proto/generated/zitadel/admin_pb';
import { PasswordLockoutPolicyView } from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
    selector: 'app-password-lockout-policy',
    templateUrl: './password-lockout-policy.component.html',
    styleUrls: ['./password-lockout-policy.component.scss'],
})
export class PasswordLockoutPolicyComponent implements OnDestroy {
    @Input() public service!: ManagementService | AdminService;
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;


    public lockoutForm!: FormGroup;
    public lockoutData!: PasswordLockoutPolicyView.AsObject;
    private sub: Subscription = new Subscription();
    public PolicyComponentServiceType: any = PolicyComponentServiceType;

    constructor(
        private route: ActivatedRoute,
        private toast: ToastService,
        private injector: Injector,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.serviceType = data.serviceType;

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
            this.fetchData();
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private fetchData(): void {
        this.getData().then(data => {
            if (data) {
                this.lockoutData = data.toObject() as PasswordLockoutPolicyView.AsObject;
            }
        });
    }

    private getData(): Promise<PasswordLockoutPolicyView | DefaultPasswordLockoutPolicyView> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).GetPasswordLockoutPolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultPasswordLockoutPolicy();
        }
    }

    public removePolicy(): void {
        if (this.service instanceof ManagementService) {
            this.service.RemovePasswordLockoutPolicy().then(() => {
                this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
                this.fetchData();
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
                this.toast.showInfo('POLICY.TOAST.SET', true);
            }).catch(error => {
                this.toast.showError(error);
            });
        } else {
            if ((this.lockoutData as PasswordLockoutPolicyView.AsObject).pb_default) {
                promise = this.service.CreatePasswordLockoutPolicy(
                    this.lockoutData.maxAttempts,
                    this.lockoutData.showLockoutFailure,
                ).then(() => {
                    this.toast.showInfo('POLICY.TOAST.SET', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            } else {
                promise = this.service.UpdatePasswordLockoutPolicy(
                    this.lockoutData.maxAttempts,
                    this.lockoutData.showLockoutFailure,
                ).then(() => {
                    this.toast.showInfo('POLICY.TOAST.SET', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        }
    }

    public get isDefault(): boolean {
        if (this.lockoutData && this.serviceType === PolicyComponentServiceType.MGMT) {
            return (this.lockoutData as PasswordLockoutPolicyView.AsObject).pb_default;
        } else {
            return false;
        }
    }
}

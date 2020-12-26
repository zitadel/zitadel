import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { DefaultPasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/admin_pb';
import { PasswordComplexityPolicyView } from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
    selector: 'app-password-policy',
    templateUrl: './password-complexity-policy.component.html',
    styleUrls: ['./password-complexity-policy.component.scss'],
})
export class PasswordComplexityPolicyComponent implements OnDestroy {
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    public service!: ManagementService | AdminService;

    public complexityData!: PasswordComplexityPolicyView.AsObject | DefaultPasswordComplexityPolicy.AsObject;

    private sub: Subscription = new Subscription();
    public PolicyComponentServiceType: any = PolicyComponentServiceType;

    public loading: boolean = false;
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

    public fetchData(): void {
        this.loading = true;

        this.getData().then(data => {
            if (data) {
                this.complexityData = data.toObject();
                this.loading = false;
            }
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData():
        Promise<PasswordComplexityPolicyView | DefaultPasswordComplexityPolicy> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).GetPasswordComplexityPolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultPasswordComplexityPolicy();
        }
    }

    public removePolicy(): void {
        if (this.service instanceof ManagementService) {
            this.service.removePasswordComplexityPolicy().then(() => {
                this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
                setTimeout(() => {
                    this.fetchData();
                }, 1000);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
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
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                if ((this.complexityData as PasswordComplexityPolicyView.AsObject).pb_default) {
                    (this.service as ManagementService).CreatePasswordComplexityPolicy(

                        this.complexityData.hasLowercase,
                        this.complexityData.hasUppercase,
                        this.complexityData.hasNumber,
                        this.complexityData.hasSymbol,
                        this.complexityData.minLength,
                    ).then(() => {
                        this.toast.showInfo('POLICY.TOAST.SET', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                } else {
                    (this.service as ManagementService).UpdatePasswordComplexityPolicy(
                        this.complexityData.hasLowercase,
                        this.complexityData.hasUppercase,
                        this.complexityData.hasNumber,
                        this.complexityData.hasSymbol,
                        this.complexityData.minLength,
                    ).then(() => {
                        this.toast.showInfo('POLICY.TOAST.SET', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
                break;
            case PolicyComponentServiceType.ADMIN:
                (this.service as AdminService).UpdateDefaultPasswordComplexityPolicy(
                    this.complexityData.hasLowercase,
                    this.complexityData.hasUppercase,
                    this.complexityData.hasNumber,
                    this.complexityData.hasSymbol,
                    this.complexityData.minLength,
                ).then(() => {
                    this.toast.showInfo('POLICY.TOAST.SET', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
                break;
        }
    }

    public get isDefault(): boolean {
        if (this.complexityData && this.serviceType === PolicyComponentServiceType.MGMT) {
            return (this.complexityData as PasswordComplexityPolicyView.AsObject).pb_default;
        } else {
            return false;
        }
    }
}

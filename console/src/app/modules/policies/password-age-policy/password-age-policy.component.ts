import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { DefaultPasswordAgePolicyView } from 'src/app/proto/generated/admin_pb';
import { PasswordAgePolicyView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policy-component-types.enum';


@Component({
    selector: 'app-password-age-policy',
    templateUrl: './password-age-policy.component.html',
    styleUrls: ['./password-age-policy.component.scss'],
})
export class PasswordAgePolicyComponent implements OnDestroy {
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    public service!: AdminService | ManagementService;

    public ageData!: PasswordAgePolicyView.AsObject | DefaultPasswordAgePolicyView.AsObject;

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
            this.getData().then(data => {
                if (data) {
                    this.ageData = data.toObject();
                }
            });
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData():
        Promise<PasswordAgePolicyView | DefaultPasswordAgePolicyView> {

        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).GetPasswordAgePolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultPasswordAgePolicy();
        }
    }

    public removePolicy(): void {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            (this.service as ManagementService).RemovePasswordAgePolicy().then(() => {
                this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
                setTimeout(() => {
                    this.getData();
                }, 1000);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public incrementExpireWarnDays(): void {
        if (this.ageData?.expireWarnDays !== undefined) {
            this.ageData.expireWarnDays++;
        }
    }

    public decrementExpireWarnDays(): void {
        if (this.ageData?.expireWarnDays && this.ageData?.expireWarnDays > 0) {
            this.ageData.expireWarnDays--;
        }
    }

    public incrementMaxAgeDays(): void {
        if (this.ageData?.maxAgeDays !== undefined) {
            this.ageData.maxAgeDays++;
        }
    }

    public decrementMaxAgeDays(): void {
        if (this.ageData?.maxAgeDays && this.ageData?.maxAgeDays > 0) {
            this.ageData.maxAgeDays--;
        }
    }

    public savePolicy(): void {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                if ((this.ageData as PasswordAgePolicyView.AsObject).pb_default) {
                    (this.service as ManagementService).CreatePasswordAgePolicy(
                        this.ageData.maxAgeDays,
                        this.ageData.expireWarnDays,
                    ).then(() => {
                        this.toast.showInfo('POLICY.TOAST.SET', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                } else {
                    (this.service as ManagementService).UpdatePasswordAgePolicy(
                        this.ageData.maxAgeDays,
                        this.ageData.expireWarnDays,
                    ).then(() => {
                        this.toast.showInfo('POLICY.TOAST.SET', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
                break;
            case PolicyComponentServiceType.ADMIN:
                (this.service as AdminService).UpdateDefaultPasswordAgePolicy(
                    this.ageData.maxAgeDays,
                    this.ageData.expireWarnDays,
                ).then(() => {
                    this.toast.showInfo('POLICY.TOAST.SET', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
                break;
        }
    }

    public get isDefault(): boolean {
        if (this.ageData && this.serviceType === PolicyComponentServiceType.MGMT) {
            return (this.ageData as PasswordAgePolicyView.AsObject).pb_default;
        } else {
            return false;
        }
    }
}

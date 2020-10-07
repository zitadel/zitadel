import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { DefaultPasswordAgePolicyView } from 'src/app/proto/generated/admin_pb';
import { PasswordAgePolicyView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentAction } from '../policy-component-action.enum';
import { PolicyComponentServiceType } from '../policy-component-types.enum';


@Component({
    selector: 'app-password-age-policy',
    templateUrl: './password-age-policy.component.html',
    styleUrls: ['./password-age-policy.component.scss'],
})
export class PasswordAgePolicyComponent implements OnDestroy {
    public title: string = '';
    public desc: string = '';

    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;
    public service!: AdminService | ManagementService;

    public PolicyComponentAction: any = PolicyComponentAction;

    public ageData!: PasswordAgePolicyView.AsObject | DefaultPasswordAgePolicyView.AsObject;

    private sub: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private router: Router,
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
            this.title = 'ORG.POLICY.PWD_AGE.TITLECREATE';
            this.desc = 'ORG.POLICY.PWD_AGE.DESCRIPTIONCREATE';

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
        this.title = 'ORG.POLICY.PWD_AGE.TITLE';
        this.desc = 'ORG.POLICY.PWD_AGE.DESCRIPTION';
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).GetPasswordAgePolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultPasswordAgePolicy();
        }
    }

    public deletePolicy(): void {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            (this.service as ManagementService).RemovePasswordAgePolicy().then(() => {
                this.toast.showInfo('Successfully deleted');
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
                if (this.componentAction === PolicyComponentAction.CREATE) {
                    (this.service as ManagementService).CreatePasswordAgePolicy(
                        this.ageData.maxAgeDays,
                        this.ageData.expireWarnDays,
                    ).then(() => {
                        this.router.navigate(['/org']);
                    }).catch(error => {
                        this.toast.showError(error);
                    });

                } else if (this.componentAction === PolicyComponentAction.MODIFY) {
                    (this.service as ManagementService).UpdatePasswordAgePolicy(
                        this.ageData.maxAgeDays,
                        this.ageData.expireWarnDays,
                    ).then(() => {
                        this.router.navigate(['/org']);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
                break;
            case PolicyComponentServiceType.ADMIN:
                if (this.componentAction === PolicyComponentAction.MODIFY) {
                    (this.service as AdminService).UpdateDefaultPasswordAgePolicy(
                        this.ageData.maxAgeDays,
                        this.ageData.expireWarnDays,
                    ).then(() => {
                        this.router.navigate(['/iam']);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
                break;
        }

    }
}

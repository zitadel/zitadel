import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { PasswordAgePolicy, PasswordAgePolicyView } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentAction } from '../policy-component-action.enum';


@Component({
    selector: 'app-password-age-policy',
    templateUrl: './password-age-policy.component.html',
    styleUrls: ['./password-age-policy.component.scss'],
})
export class PasswordAgePolicyComponent implements OnDestroy {
    public title: string = '';
    public desc: string = '';

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    public PolicyComponentAction: any = PolicyComponentAction;

    public ageData!: PasswordAgePolicy.AsObject;

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
        })).subscribe(_ => {
            this.title = 'ORG.POLICY.PWD_AGE.TITLECREATE';
            this.desc = 'ORG.POLICY.PWD_AGE.DESCRIPTIONCREATE';

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData().then(data => {
                    if (data) {
                        this.ageData = data.toObject() as PasswordAgePolicy.AsObject;
                    }
                });
            }
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData():
        Promise<PasswordAgePolicyView> {
        this.title = 'ORG.POLICY.PWD_AGE.TITLE';
        this.desc = 'ORG.POLICY.PWD_AGE.DESCRIPTION';
        return this.mgmtService.GetPasswordAgePolicy();
    }

    public deletePolicy(): void {
        this.mgmtService.RemovePasswordAgePolicy().then(() => {
            this.toast.showInfo('Successfully deleted');
        }).catch(error => {
            this.toast.showError(error);
        });
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
        if (this.componentAction === PolicyComponentAction.CREATE) {

            this.mgmtService.CreatePasswordAgePolicy(
                this.ageData.maxAgeDays,
                this.ageData.expireWarnDays,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });

        } else if (this.componentAction === PolicyComponentAction.MODIFY) {

            this.mgmtService.UpdatePasswordAgePolicy(
                this.ageData.maxAgeDays,
                this.ageData.expireWarnDays,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }
}

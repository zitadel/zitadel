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
    selector: 'app-password-age-policy',
    templateUrl: './password-age-policy.component.html',
    styleUrls: ['./password-age-policy.component.scss'],
})
export class PasswordAgePolicyComponent implements OnDestroy {
    public title: string = '';
    public desc: string = '';

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    public PolicyComponentAction: any = PolicyComponentAction;

    public ageForm!: FormGroup;
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
        })).subscribe(params => {
            this.title = 'ORG.POLICY.PWD_AGE.TITLECREATE';
            this.desc = 'ORG.POLICY.PWD_AGE.DESCRIPTIONCREATE';

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData(params).then(data => {
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

    private async getData(params: any):
        Promise<PasswordLockoutPolicy | PasswordAgePolicy | PasswordComplexityPolicy | OrgIamPolicy | undefined> {
        this.title = 'ORG.POLICY.PWD_AGE.TITLE';
        this.desc = 'ORG.POLICY.PWD_AGE.DESCRIPTION';
        return this.mgmtService.GetPasswordAgePolicy();
    }

    public deletePolicy(): void {
        this.mgmtService.DeletePasswordAgePolicy(this.ageData.id).then(() => {
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
                this.ageData.description,
                this.ageData.maxAgeDays,
                this.ageData.expireWarnDays,
            ).then(() => {
                this.router.navigate(['org']);
            }).catch(error => {
                this.toast.showError(error);
            });

        } else if (this.componentAction === PolicyComponentAction.MODIFY) {

            this.mgmtService.UpdatePasswordAgePolicy(
                this.ageData.description,
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

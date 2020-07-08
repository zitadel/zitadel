import { Component, OnDestroy, OnInit } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject, Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import {
    OrgIamPolicy,
    PasswordAgePolicy,
    PasswordComplexityPolicy,
    PasswordLockoutPolicy,
} from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { OrgService } from 'src/app/services/org.service';
import { StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

export enum PolicyComponentAction {
    CREATE = 'create',
    MODIFY = 'modify',
}

export enum PolicyComponentType {
    LOCKOUT = 'lockout',
    AGE = 'age',
    COMPLEXITY = 'complexity',
    IAM_POLICY = 'iam_policy',
}

@Component({
    selector: 'app-password-policy',
    templateUrl: './password-policy.component.html',
    styleUrls: ['./password-policy.component.scss'],
})
export class PasswordPolicyComponent implements OnInit, OnDestroy {
    titleSub: BehaviorSubject<string> = new BehaviorSubject('');
    descSub: BehaviorSubject<string> = new BehaviorSubject('');

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    policyData!: PasswordLockoutPolicy.AsObject |
        PasswordAgePolicy.AsObject |
        PasswordComplexityPolicy.AsObject |
        OrgIamPolicy.AsObject;
    policyType: PolicyComponentType = PolicyComponentType.COMPLEXITY;

    public PolicyComponentType: any = PolicyComponentType;
    public PolicyComponentAction: any = PolicyComponentAction;

    public lockoutForm!: FormGroup;
    public ageForm!: FormGroup;

    public complexityData: any = {
        minLength: 8,
        description: '',
        hasNumber: true,
        hasSymbol: true,
        hasLowercase: true,
        hasUppercase: true,
    };

    public lockoutData: any = {
        description: '',
        maxAttempts: 5,
        showLockOutFailures: true,
    };

    public ageData: any = {
        description: '',
        expireWarnDays: 80,
        maxAgeDays: 90,
    };

    public iamData: any = {
        description: '',
        userLoginMustBeDomain: false,
    };

    private sub: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private adminService: AdminService,
        private orgService: OrgService,
        private router: Router,
        private toast: ToastService,
        private sessionStorage: StorageService,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.componentAction = data.action;
            return this.route.params;
        })).subscribe(params => {
            this.policyType = params.policytype;

            switch (params.policytype) {
                case PolicyComponentType.LOCKOUT:
                    this.titleSub.next('ORG.POLICY.PWD_LOCKOUT.TITLECREATE');
                    this.descSub.next('ORG.POLICY.PWD_LOCKOUT.DESCRIPTIONCREATE');
                    break;
                case PolicyComponentType.AGE:
                    this.titleSub.next('ORG.POLICY.PWD_AGE.TITLECREATE');
                    this.descSub.next('ORG.POLICY.PWD_AGE.DESCRIPTIONCREATE');
                    break;
                case PolicyComponentType.COMPLEXITY:
                    this.titleSub.next('ORG.POLICY.PWD_COMPLEXITY.TITLECREATE');
                    this.descSub.next('ORG.POLICY.PWD_COMPLEXITY.DESCRIPTIONCREATE');
                    break;
                case PolicyComponentType.IAM_POLICY:
                    this.titleSub.next('ORG.POLICY.IAM_POLICY.TITLECREATE');
                    this.descSub.next('ORG.POLICY.IAM_POLICY.DESCRIPTIONCREATE');
                    break;
            }

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData(params).then(data => {
                    switch (this.policyType) {
                        case PolicyComponentType.LOCKOUT:
                            this.lockoutData = data.toObject();
                            break;
                        case PolicyComponentType.AGE:
                            this.ageData = data.toObject();
                            break;
                        case PolicyComponentType.COMPLEXITY:
                            this.complexityData = data.toObject();
                            break;
                        case PolicyComponentType.IAM_POLICY:
                            this.iamData = data.toObject();
                            break;
                    }
                });
            }
        });
    }

    ngOnInit(): void {
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData(params: any): Promise<any> {
        switch (params.policytype) {
            case PolicyComponentType.LOCKOUT:
                this.titleSub.next('ORG.POLICY.PWD_LOCKOUT.TITLE');
                this.descSub.next('ORG.POLICY.PWD_LOCKOUT.DESCRIPTION');
                return this.orgService.GetPasswordLockoutPolicy();
            case PolicyComponentType.AGE:
                this.titleSub.next('ORG.POLICY.PWD_AGE.TITLE');
                this.descSub.next('ORG.POLICY.PWD_AGE.DESCRIPTION');
                return this.orgService.GetPasswordAgePolicy();
            case PolicyComponentType.COMPLEXITY:
                this.titleSub.next('ORG.POLICY.PWD_COMPLEXITY.TITLE');
                this.descSub.next('ORG.POLICY.PWD_COMPLEXITY.DESCRIPTION');
                return this.orgService.GetPasswordComplexityPolicy();
            case PolicyComponentType.IAM_POLICY:
                this.titleSub.next('ORG.POLICY.IAM_POLICY.TITLECREATE');
                this.descSub.next('ORG.POLICY.IAM_POLICY.DESCRIPTIONCREATE');
                return this.orgService.GetMyOrgIamPolicy();
        }
    }

    public incrementLength(): void {
        if (this.complexityData?.minLength !== undefined) {
            this.complexityData.minLength++;
        }
    }

    public decrementLength(): void {
        if (this.complexityData?.minLength && this.complexityData?.minLength > 0) {
            this.complexityData.minLength--;
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
            switch (this.policyType) {
                case PolicyComponentType.LOCKOUT:
                    this.orgService.CreatePasswordLockoutPolicy(
                        this.lockoutData.description,
                        this.lockoutData.maxAttempts,
                        this.lockoutData.showLockOutFailures,
                    ).then(() => {
                        this.router.navigate(['org']);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });

                    break;
                case PolicyComponentType.AGE:
                    this.orgService.CreatePasswordAgePolicy(
                        this.ageData.description,
                        this.ageData.maxAgeDays,
                        this.ageData.expireWarnDays,
                    ).then(() => {
                        this.router.navigate(['org']);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });

                    break;
                case PolicyComponentType.COMPLEXITY:
                    this.orgService.CreatePasswordComplexityPolicy(
                        this.complexityData.description,
                        this.complexityData.hasLowercase,
                        this.complexityData.hasUppercase,
                        this.complexityData.hasNumber,
                        this.complexityData.hasSymbol,
                        this.complexityData.minLength,
                    ).then(() => {
                        this.router.navigate(['org']);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
                    break;

                case PolicyComponentType.IAM_POLICY:
                    const orgId = this.sessionStorage.getItem('organization');
                    if (orgId) {
                        this.adminService.CreateOrgIamPolicy(
                            orgId,
                            this.complexityData.description,
                            this.complexityData.userLoginMustBeDomain,
                        ).then(() => {
                            this.router.navigate(['org']);
                        }).catch(error => {
                            this.toast.showError(error.message);
                        });
                    }
                    break;
            }
        } else if (this.componentAction === PolicyComponentAction.MODIFY) {
            switch (this.policyType) {
                case PolicyComponentType.LOCKOUT:
                    this.orgService.UpdatePasswordLockoutPolicy(
                        this.lockoutData.description,
                        this.lockoutData.maxAttempts,
                        this.lockoutData.showLockOutFailures,
                    ).then(() => {
                        this.router.navigate(['org']);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });

                    break;
                case PolicyComponentType.AGE:
                    this.orgService.UpdatePasswordAgePolicy(
                        this.ageData.description,
                        this.ageData.maxAgeDays,
                        this.ageData.expireWarnDays,
                    ).then(() => {
                        this.router.navigate(['org']);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });

                    break;
                case PolicyComponentType.COMPLEXITY:
                    this.orgService.UpdatePasswordComplexityPolicy(
                        this.complexityData.description,
                        this.complexityData.hasLowercase,
                        this.complexityData.hasUppercase,
                        this.complexityData.hasNumber,
                        this.complexityData.hasSymbol,
                        this.complexityData.minLength,
                    ).then(() => {
                        this.router.navigate(['org']);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
                    break;

                case PolicyComponentType.IAM_POLICY:
                    const orgId = this.sessionStorage.getItem('organization');
                    if (orgId) {
                        this.adminService.UpdateOrgIamPolicy(
                            orgId,
                            this.complexityData.description,
                            this.complexityData.userLoginMustBeDomain,
                        ).then(() => {
                            this.router.navigate(['org']);
                        }).catch(error => {
                            this.toast.showError(error.message);
                        });
                    }
                    break;
            }
        }
    }
}

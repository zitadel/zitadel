import { Component, OnDestroy, OnInit } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject, Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { PasswordAgePolicy, PasswordComplexityPolicy, PasswordLockoutPolicy } from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

export enum PolicyComponentAction {
    CREATE = 'create',
    MODIFY = 'modify',
}

export enum PolicyComponentType {
    LOCKOUT = 'lockout',
    AGE = 'age',
    COMPLEXITY = 'complexity',
}

@Component({
    selector: 'app-password-policy',
    templateUrl: './password-policy.component.html',
    styleUrls: ['./password-policy.component.scss'],
})
export class PasswordPolicyComponent implements OnInit, OnDestroy {
    public orgId: string = '';
    titleSub: BehaviorSubject<string> = new BehaviorSubject('');
    descSub: BehaviorSubject<string> = new BehaviorSubject('');

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    policyData!: PasswordLockoutPolicy.AsObject | PasswordAgePolicy.AsObject | PasswordComplexityPolicy.AsObject;
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

    private sub: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private orgService: OrgService,
        private router: Router,
        private toast: ToastService,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.componentAction = data.action;
            return this.route.params;
        })).subscribe(params => {
            this.orgId = params.id;
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
        }
    }

    public incrementLength(): void {
        if (this.complexityData?.minLength) {
            this.complexityData.minLength++;
        }
    }

    public decrementLength(): void {
        if (this.complexityData?.minLength && this.complexityData?.minLength > 0) {
            this.complexityData.minLength--;
        }
    }

    public incrementMaxAttempts(): void {
        if (this.lockoutData?.maxAttempts) {
            this.lockoutData.maxAttempts++;
        }
    }

    public decrementMaxAttempts(): void {
        if (this.lockoutData?.maxAttempts && this.lockoutData?.maxAttempts > 0) {
            this.lockoutData.maxAttempts--;
        }
    }

    public incrementExpireWarnDays(): void {
        if (this.ageData?.expireWarnDays) {
            this.ageData.expireWarnDays++;
        }
    }

    public decrementExpireWarnDays(): void {
        if (this.ageData?.expireWarnDays && this.ageData?.expireWarnDays > 0) {
            this.ageData.expireWarnDays--;
        }
    }

    public incrementMaxAgeDays(): void {
        if (this.ageData?.maxAgeDays) {
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
                        this.router.navigate(['orgs', this.orgId]);
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
                        this.router.navigate(['orgs', this.orgId]);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });

                    break;
                case PolicyComponentType.COMPLEXITY:
                    console.log(this.complexityData);
                    this.orgService.CreatePasswordComplexityPolicy(
                        this.complexityData.description,
                        this.complexityData.hasLowercase,
                        this.complexityData.hasUppercase,
                        this.complexityData.hasNumber,
                        this.complexityData.hasSymbol,
                        this.complexityData.minLength,
                    ).then(() => {
                        this.router.navigate(['orgs', this.orgId]);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
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
                        this.router.navigate(['orgs', this.orgId]);
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
                        this.router.navigate(['orgs', this.orgId]);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });

                    break;
                case PolicyComponentType.COMPLEXITY:
                    console.log(this.complexityData);
                    this.orgService.UpdatePasswordComplexityPolicy(
                        this.complexityData.description,
                        this.complexityData.hasLowercase,
                        this.complexityData.hasUppercase,
                        this.complexityData.hasNumber,
                        this.complexityData.hasSymbol,
                        this.complexityData.minLength,
                    ).then(() => {
                        this.router.navigate(['orgs', this.orgId]);
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
                    break;
            }
        }
    }
}

import { animate, style, transition, trigger } from '@angular/animations';
import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, ValidatorFn, Validators } from '@angular/forms';
import { MatSlideToggleChange } from '@angular/material/slide-toggle';
import { Router } from '@angular/router';
import { take } from 'rxjs/operators';
import { lowerCaseValidator, numberValidator, symbolValidator, upperCaseValidator } from 'src/app/pages/validators';
import { CreateHumanRequest, CreateOrgRequest, Gender, OrgSetUpResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { PasswordComplexityPolicy as MgmtPasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

function passwordConfirmValidator(c: AbstractControl): any {
    if (!c.parent || !c) {
        return;
    }
    const pwd = c.parent.get('password');
    const cpwd = c.parent.get('confirmPassword');

    if (!pwd || !cpwd) {
        return;
    }
    if (pwd.value !== cpwd.value) {
        return {
            invalid: true,
            notequal: {
                valid: false,
            },
        };
    }
}

@Component({
    selector: 'app-org-create',
    templateUrl: './org-create.component.html',
    styleUrls: ['./org-create.component.scss'],
    animations: [
        trigger('openClose', [
            transition(':enter', [
                style({ height: '0', opacity: 0 }),
                animate('150ms ease-in-out', style({ height: '*', opacity: 1 })),
            ]),
            transition(':leave', [
                animate('150ms ease-in-out', style({ height: '0', opacity: 0 })),
            ]),
        ]),
    ],
})
export class OrgCreateComponent {
    public orgForm!: FormGroup;
    public userForm!: FormGroup;
    public pwdForm!: FormGroup;

    public genders: Gender[] = [Gender.GENDER_FEMALE, Gender.GENDER_MALE, Gender.GENDER_UNSPECIFIED];
    public languages: string[] = ['de', 'en'];

    public policy!: MgmtPasswordComplexityPolicy.AsObject;
    public usePassword: boolean = false;

    public forSelf: boolean = true;

    constructor(
        private router: Router,
        private toast: ToastService,
        private adminService: AdminService,
        private _location: Location,
        private fb: FormBuilder,
        private mgmtService: ManagementService,
        private authService: GrpcAuthService,
    ) {
        this.authService.isAllowed(['iam.write']).pipe(take(1)).subscribe((allowed) => {
            if (allowed) {
                this.forSelf = false;
            }
        });

        this.orgForm = this.fb.group({
            name: ['', [Validators.required]],
            domain: [''],
        });

        this.initForm();
    }

    public createSteps: number = 2;
    public currentCreateStep: number = 1;

    public finish(): void {
        const createOrgRequest: CreateOrgRequest = new CreateOrgRequest();
        createOrgRequest.setName(this.name?.value);
        createOrgRequest.setDomain(this.domain?.value);

        const humanRequest: CreateHumanRequest = new CreateHumanRequest();
        humanRequest.setEmail(this.email?.value);
        humanRequest.setFirstName(this.firstName?.value);
        humanRequest.setLastName(this.lastName?.value);
        humanRequest.setNickName(this.nickName?.value);
        humanRequest.setGender(this.gender?.value);
        humanRequest.setPreferredLanguage(this.preferredLanguage?.value);

        if (this.usePassword && this.password) {
            humanRequest.setPassword(this.password?.value);
        }

        this.adminService
            .SetUpOrg(createOrgRequest, humanRequest)
            .then((org: OrgSetUpResponse) => {
                this.router.navigate(['/org/overview']);
                // const orgResp = org.getOrg();
                // if (orgResp) {
                //     this.authService.setActiveOrg(orgResp.toObject());
                //     this.router.navigate(['/org']);
                // } else {
                //     this.router.navigate(['/org', 'overview']);
                // }
            })
            .catch(error => {
                this.toast.showError(error);
            });
    }

    public next(): void {
        this.currentCreateStep++;
    }

    public previous(): void {
        this.currentCreateStep--;
    }

    private initForm(): void {
        this.userForm = this.fb.group({
            userName: ['', [Validators.required]],
            firstName: ['', [Validators.required]],
            lastName: ['', [Validators.required]],
            email: ['', [Validators.required]],
            gender: [''],
            nickName: [''],
            preferredLanguage: [''],
        });
    }

    public initPwdValidators(): void {
        const validators: Validators[] = [Validators.required];

        if (this.usePassword) {
            this.mgmtService.GetDefaultPasswordComplexityPolicy().then(data => {
                this.policy = data.toObject();

                if (this.policy.minLength) {
                    validators.push(Validators.minLength(this.policy.minLength));
                }
                if (this.policy.hasLowercase) {
                    validators.push(lowerCaseValidator);
                }
                if (this.policy.hasUppercase) {
                    validators.push(upperCaseValidator);
                }
                if (this.policy.hasNumber) {
                    validators.push(numberValidator);
                }
                if (this.policy.hasSymbol) {
                    validators.push(symbolValidator);
                }

                const pwdValidators = [...validators] as ValidatorFn[];
                const confirmPwdValidators = [...validators, passwordConfirmValidator] as ValidatorFn[];
                this.pwdForm = this.fb.group({
                    password: ['', pwdValidators],
                    confirmPassword: ['', confirmPwdValidators],
                });
            });
        } else {
            this.pwdForm = this.fb.group({
                password: ['', []],
                confirmPassword: ['', []],
            });
        }
    }

    public changeSelf(change: MatSlideToggleChange): void {
        if (change.checked) {
            this.createSteps = 1;

            this.orgForm = this.fb.group({
                name: ['', [Validators.required]],
            });
        } else {
            this.createSteps = 2;

            this.orgForm = this.fb.group({
                name: ['', [Validators.required]],
                domain: [''],
            });
        }
    }

    public createOrgForSelf(): void {
        if (this.name && this.name.value) {
            this.mgmtService.CreateOrg(this.name.value).then((org) => {
                this.router.navigate(['/org/overview']);
                // const newOrg = org.toObject();
                // setTimeout(() => {
                //     this.authService.setActiveOrg(newOrg);
                //     this.router.navigate(['/org']);
                // }, 1000);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public get name(): AbstractControl | null {
        return this.orgForm.get('name');
    }

    public get domain(): AbstractControl | null {
        return this.orgForm.get('domain');
    }

    public get userName(): AbstractControl | null {
        return this.userForm.get('userName');
    }

    public get firstName(): AbstractControl | null {
        return this.userForm.get('firstName');
    }

    public get lastName(): AbstractControl | null {
        return this.userForm.get('lastName');
    }

    public get email(): AbstractControl | null {
        return this.userForm.get('email');
    }

    public get nickName(): AbstractControl | null {
        return this.userForm.get('nickName');
    }

    public get preferredLanguage(): AbstractControl | null {
        return this.userForm.get('preferredLanguage');
    }

    public get gender(): AbstractControl | null {
        return this.userForm.get('gender');
    }

    public get password(): AbstractControl | null {
        return this.pwdForm.get('password');
    }

    public get confirmPassword(): AbstractControl | null {
        return this.pwdForm.get('confirmPassword');
    }

    public close(): void {
        this._location.back();
    }
}

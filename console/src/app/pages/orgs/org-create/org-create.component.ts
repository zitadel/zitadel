import { animate, style, transition, trigger } from '@angular/animations';
import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, ValidatorFn, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { lowerCaseValidator, numberValidator, symbolValidator, upperCaseValidator } from 'src/app/pages/validators';
import { CreateOrgRequest, CreateUserRequest, Gender, OrgSetUpResponse } from 'src/app/proto/generated/admin_pb';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/auth_pb';
import { AdminService } from 'src/app/services/admin.service';
import { OrgService } from 'src/app/services/org.service';
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

    public policy!: PasswordComplexityPolicy.AsObject;
    public usePassword: boolean = false;
    constructor(
        private router: Router,
        private toast: ToastService,
        private adminService: AdminService,
        private _location: Location,
        private fb: FormBuilder,
        private orgService: OrgService,
    ) {
        const validators: Validators[] = [];

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

        const registerUserRequest: CreateUserRequest = new CreateUserRequest();
        registerUserRequest.setUserName(this.userName?.value);
        registerUserRequest.setEmail(this.email?.value);
        registerUserRequest.setFirstName(this.firstName?.value);
        registerUserRequest.setLastName(this.lastName?.value);
        registerUserRequest.setNickName(this.nickName?.value);
        registerUserRequest.setGender(this.gender?.value);
        registerUserRequest.setPassword(this.password?.value);
        registerUserRequest.setPreferredLanguage(this.preferredLanguage?.value);
        this.adminService
            .SetUpOrg(createOrgRequest, registerUserRequest)
            .then((data: OrgSetUpResponse) => {
                this.router.navigate(['orgs', data.toObject().org?.id]);
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
        console.log('init without pwd');
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

        console.log(this.usePassword);
        if (this.usePassword) {
            this.orgService.GetDefaultPasswordComplexityPolicy().then(data => {
                this.policy = data.toObject();
                console.log(this.policy);

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

                // this.initForm(validators);
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

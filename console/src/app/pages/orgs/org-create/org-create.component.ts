import { animate, style, transition, trigger } from '@angular/animations';
import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { CreateOrgRequest, Gender, OrgSetUpResponse } from 'src/app/proto/generated/admin_pb';
import { RegisterUserRequest } from 'src/app/proto/generated/auth_pb';
import { AdminService } from 'src/app/services/admin.service';
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
        return { invalid: true };
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

    public genders: Gender[] = [Gender.FEMALE, Gender.MALE, Gender.UNKNOWN_GENDER];
    public languages: string[] = ['de', 'en'];
    constructor(
        private router: Router,
        private toast: ToastService,
        private adminService: AdminService,
        private _location: Location,
        private fb: FormBuilder,
    ) {
        this.orgForm = this.fb.group({
            name: ['', [Validators.required]],
            domain: ['', [Validators.required]],
        });

        this.userForm = this.fb.group({
            firstName: ['', [Validators.required]],
            lastName: ['', [Validators.required]],
            displayName: [''],
            email: ['', [Validators.required]],
            gender: [''],
            nickName: [''],
            preferredLanguage: [''],
            password: ['', [Validators.required]],
            confirmPassword: ['', [Validators.required, passwordConfirmValidator]],
        });
    }

    public createSteps: number = 2;
    public currentCreateStep: number = 1;

    public finish(): void {
        const createOrgRequest: CreateOrgRequest = new CreateOrgRequest();
        createOrgRequest.setName(this.name?.value);
        createOrgRequest.setDomain(this.domain?.value);

        const registerUserRequest: RegisterUserRequest = new RegisterUserRequest();
        registerUserRequest.setEmail(this.email?.value);
        registerUserRequest.setFirstName(this.firstName?.value);
        registerUserRequest.setLastName(this.lastName?.value);
        registerUserRequest.setNickName(this.nickName?.value);
        registerUserRequest.setDisplayName(this.displayName?.value);
        registerUserRequest.setGender(this.gender?.value);
        registerUserRequest.setPassword(this.password?.value);
        registerUserRequest.setPreferredLanguage(this.preferredLanguage?.value);

        this.adminService
            .SetUpOrg(createOrgRequest, registerUserRequest)
            .then((data: OrgSetUpResponse) => {
                this.router.navigate(['orgs', data.toObject().org?.id]);
            })
            .catch(data => {
                this.toast.showError(data.message);
            });
    }

    public next(): void {
        this.currentCreateStep++;
    }

    public previous(): void {
        this.currentCreateStep--;
    }

    public get name(): AbstractControl | null {
        return this.orgForm.get('name');
    }

    public get domain(): AbstractControl | null {
        return this.orgForm.get('domain');
    }


    public get firstName(): AbstractControl | null {
        return this.userForm.get('firstName');
    }

    public get lastName(): AbstractControl | null {
        return this.userForm.get('lastName');
    }

    public get displayName(): AbstractControl | null {
        return this.userForm.get('displayName');
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
        return this.userForm.get('password');
    }

    public get confirmPassword(): AbstractControl | null {
        return this.userForm.get('confirmPassword');
    }

    public close(): void {
        this._location.back();
    }
}

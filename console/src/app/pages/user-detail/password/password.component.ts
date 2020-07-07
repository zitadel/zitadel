import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/management_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

import { lowerCaseValidator, numberValidator, symbolValidator, upperCaseValidator } from '../validators';

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
        return { invalid: true, notequal: 'Password is not equal' };
    }
}

@Component({
    selector: 'app-password',
    templateUrl: './password.component.html',
    styleUrls: ['./password.component.scss'],
})
export class PasswordComponent implements OnInit {
    userId: string = '';

    public policy!: PasswordComplexityPolicy.AsObject;
    public passwordForm!: FormGroup;

    constructor(
        private orgService: OrgService,
        activatedRoute: ActivatedRoute,
        private fb: FormBuilder,
        private userService: AuthUserService,
        private mgmtUserService: MgmtUserService,
        private toast: ToastService,
    ) {

        activatedRoute.params.subscribe(data => {
            const { id } = data;
            if (id) {
                this.userId = id;
            }

            const validators: Validators[] = [Validators.required];
            this.orgService.GetPasswordComplexityPolicy().then(complexity => {
                this.policy = complexity.toObject();
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

                this.setupForm(validators);
            }).catch(error => {
                this.setupForm(validators);
            });
        });
    }

    ngOnInit(): void {
    }

    setupForm(validators: Validators[]): void {
        if (this.userId) {
            this.passwordForm = this.fb.group({
                password: ['', validators],
                confirmPassword: ['', [...validators, passwordConfirmValidator]],
            });
        } else {
            this.passwordForm = this.fb.group({
                currentPassword: ['', validators],
                newPassword: ['', validators],
                confirmPassword: ['', [...validators, passwordConfirmValidator]],
            });
        }
    }

    public setInitialPassword(userId: string): void {
        if (this.passwordForm.valid && this.password && this.password.value) {
            this.mgmtUserService.SetInitialPassword(userId, this.password.value).then((data: any) => {
                this.toast.showInfo('Set initial Password');
                window.history.back();
            }).catch(data => {
                this.toast.showError(data.message);
            });
        }
    }

    public setPassword(): void {
        if (this.passwordForm.valid && this.currentPassword &&
            this.currentPassword.value &&
            this.newPassword && this.newPassword.value && this.newPassword.valid) {
            this.userService
                .ChangeMyPassword(this.currentPassword.value, this.newPassword.value).then((data: any) => {
                    this.toast.showInfo('Password Set');
                    window.history.back();
                }).catch(data => {
                    this.toast.showError(data.message);
                });
        }
    }

    public get password(): AbstractControl | null {
        return this.passwordForm.get('password');
    }

    public get newPassword(): AbstractControl | null {
        return this.passwordForm.get('newPassword');
    }

    public get currentPassword(): AbstractControl | null {
        return this.passwordForm.get('newPassword');
    }

    public get confirmPassword(): AbstractControl | null {
        return this.passwordForm.get('confirmPassword');
    }
}

import { Component, OnDestroy } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { Gender, UserAddress, UserEmail, UserPhone, UserProfile } from 'src/app/proto/generated/auth_pb';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/management_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

import { CodeDialogComponent } from '../code-dialog/code-dialog.component';

function passwordConfirmValidator(c: AbstractControl): any {
    if (!c.parent || !c) {
        return;
    }
    const pwd = c.parent.get('newPassword');
    const cpwd = c.parent.get('confirmPassword');

    if (!pwd || !cpwd) {
        return;
    }
    if (pwd.value !== cpwd.value) {
        return { invalid: true, notequal: 'Password is not equal' };
    }
}

@Component({
    selector: 'app-auth-user-detail',
    templateUrl: './auth-user-detail.component.html',
    styleUrls: ['./auth-user-detail.component.scss'],
})
export class AuthUserDetailComponent implements OnDestroy {
    public profile!: UserProfile.AsObject;
    public email: UserEmail.AsObject = { email: '' } as any;
    public phone: UserPhone.AsObject = { phone: '' } as any;
    public address: UserAddress.AsObject = { id: '' } as any;
    public genders: Gender[] = [Gender.GENDER_MALE, Gender.GENDER_FEMALE, Gender.GENDER_DIVERSE];
    public languages: string[] = ['de', 'en'];

    public passwordForm!: FormGroup;
    public addressForm!: FormGroup;
    private subscription: Subscription = new Subscription();

    public emailEditState: boolean = false;
    public phoneEditState: boolean = false;

    public loading: boolean = false;

    public policy!: PasswordComplexityPolicy.AsObject;

    constructor(
        public translate: TranslateService,
        private toast: ToastService,
        private mgmtUserService: MgmtUserService,
        private userService: AuthUserService,
        private fb: FormBuilder,
        private dialog: MatDialog,
        private orgService: OrgService,
    ) {
        const validators: Validators[] = [Validators.required];
        this.orgService.GetPasswordComplexityPolicy().then(data => {
            this.policy = data.toObject();
            if (this.policy.minLength) {
                validators.push(Validators.minLength(this.policy.minLength));
            }
            if (this.policy.hasLowercase) {
                validators.push(Validators.pattern(/[a-z]/g));
            }
            if (this.policy.hasUppercase) {
                validators.push(Validators.pattern(/[A-Z]/g));
            }
            if (this.policy.hasNumber) {
                validators.push(Validators.pattern(/[0-9]/g));
            }
            if (this.policy.hasSymbol) {
                // All characters that are not a digit or an English letter \W or a whitespace \S
                validators.push(Validators.pattern(/[\W\S]/));
            }

            this.passwordForm = this.fb.group({
                currentPassword: ['', []],
                newPassword: ['', validators],
                confirmPassword: ['', [...validators, passwordConfirmValidator]],
            });
        }).catch(error => {
            console.log('no password complexity policy defined!');
            this.passwordForm = this.fb.group({
                currentPassword: ['', []],
                newPassword: ['', []],
                confirmPassword: ['', [passwordConfirmValidator]],
            });
        });

        this.addressForm = this.fb.group({
            streetAddress: [''],
            postalCode: [''],
            locality: [''],
            region: [''],
            country: [''],
        });

        this.loading = true;
        this.getData().then(() => {
            this.loading = false;
        }).catch(error => {
            this.loading = false;
        });
    }

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }

    public saveProfile(profileData: UserProfile.AsObject): void {
        console.log(profileData);
        this.profile.firstName = profileData.firstName;
        this.profile.lastName = profileData.lastName;
        this.profile.nickName = profileData.nickName;
        this.profile.displayName = profileData.displayName;
        this.profile.gender = profileData.gender;
        this.profile.preferredLanguage = profileData.preferredLanguage;
        console.log(this.profile);
        this.userService
            .SaveMyUserProfile(this.profile as UserProfile.AsObject)
            .then((data: UserProfile) => {
                this.toast.showInfo('Saved Profile');
                this.profile = data.toObject();
            })
            .catch(data => {
                this.toast.showError(data.message);
            });
    }

    public setPassword(): void {
        if (this.passwordForm.valid && this.currentPassword &&
            this.currentPassword.value &&
            this.newPassword && this.newPassword.value && this.newPassword.valid) {
            this.userService
                .ChangeMyPassword(this.currentPassword.value, this.newPassword.value).then((data: any) => {
                    this.toast.showInfo('Password Set');
                    this.email = data.toObject();
                }).catch(data => {
                    this.toast.showError(data.message);
                });
        }
    }

    public saveEmail(): void {
        this.emailEditState = false;

        this.userService
            .SaveMyUserEmail(this.email).then((data: UserEmail) => {
                this.toast.showInfo('Saved Email');
                this.email = data.toObject();
                this.emailEditState = false;
            }).catch(data => {
                this.toast.showError(data.message);
                this.emailEditState = false;
            });
    }

    public deletePhone(): void {
        this.phone.phone = '';
        this.savePhone();
    }

    public enterCode(): void {
        const dialogRef = this.dialog.open(CodeDialogComponent, {
            data: {
                number: this.phone.phone,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(code => {
            if (code) {
                console.log(code);
                this.userService.VerifyMyUserPhone(code).then(() => {
                    this.toast.showInfo('Verified Phone');
                }).catch(error => {
                    this.toast.showError(error.message);
                });
            }
        });
    }

    public changedLanguage(language: string): void {
        this.translate.use(language);
    }

    public resendVerification(): void {
        this.userService.ResendEmailVerification().then((data: any) => {
            this.toast.showInfo('Saved Email');
            this.email = data.toObject();
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public resendPhoneVerification(): void {
        console.log('resendverification');
        this.userService.ResendPhoneVerification().then((data: any) => {
            console.log(data);
            this.toast.showInfo('Phoneverification was successfully sent!');
            this.email = data.toObject();
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public savePhone(): void {
        this.phoneEditState = false;
        if (!this.phone.id) {
            this.phone.id = this.profile.id;
        }
        this.userService
            .SaveMyUserPhone(this.phone).then((data: UserPhone) => {
                this.toast.showInfo('Saved Phone');
                this.phone = data.toObject();
                this.phoneEditState = false;
            }).catch(data => {
                this.toast.showError(data.message);
                this.phoneEditState = false;
            });
    }

    public saveAddress(): void {
        this.address = this.addressForm.value;
        this.userService
            .SaveMyUserAddress(this.address as UserAddress.AsObject).then((data: UserAddress) => {
                this.toast.showInfo('Saved Address');
                this.address = data.toObject();
            }).catch(data => {
                this.toast.showError(data.message);
            });
    }

    public get streetAddress(): AbstractControl | null {
        return this.addressForm.get('streetAddress');
    }
    public get postalCode(): AbstractControl | null {
        return this.addressForm.get('postalCode');
    }
    public get locality(): AbstractControl | null {
        return this.addressForm.get('locality');
    }
    public get region(): AbstractControl | null {
        return this.addressForm.get('region');
    }
    public get country(): AbstractControl | null {
        return this.addressForm.get('country');
    }

    public get currentPassword(): AbstractControl | null {
        return this.passwordForm.get('currentPassword');
    }
    public get newPassword(): AbstractControl | null {
        return this.passwordForm.get('newPassword');
    }
    public get confirmPassword(): AbstractControl | null {
        return this.passwordForm.get('confirmPassword');
    }

    private async getData(): Promise<void> {
        this.profile = (await this.userService.GetMyUserProfile()).toObject();
        this.email = (await this.userService.GetMyUserEmail()).toObject();
        this.phone = (await this.userService.GetMyUserPhone()).toObject();
        this.address = (await this.userService.GetMyUserAddress()).toObject();

        console.log(this.profile);
        this.addressForm.patchValue(this.address);
    }
}

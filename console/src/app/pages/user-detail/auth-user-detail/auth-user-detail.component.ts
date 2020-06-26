import { Component, OnDestroy } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { Gender, UserAddress, UserEmail, UserPhone, UserProfile, UserView } from 'src/app/proto/generated/auth_pb';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/management_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

import { CodeDialogComponent } from '../code-dialog/code-dialog.component';
import { lowerCaseValidator, numberValidator, symbolValidator, upperCaseValidator } from '../validators';

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
        return {
            invalid: true,
            notequal: {
                valid: false,
            },
        };
    }
}

@Component({
    selector: 'app-auth-user-detail',
    templateUrl: './auth-user-detail.component.html',
    styleUrls: ['./auth-user-detail.component.scss'],
})
export class AuthUserDetailComponent implements OnDestroy {
    public user!: UserView.AsObject;
    public address: UserAddress.AsObject = { id: '' } as any;
    public genders: Gender[] = [Gender.GENDER_MALE, Gender.GENDER_FEMALE, Gender.GENDER_DIVERSE];
    public languages: string[] = ['de', 'en'];

    public passwordForm!: FormGroup;
    // public addressForm!: FormGroup;
    private subscription: Subscription = new Subscription();

    public emailEditState: boolean = false;
    public phoneEditState: boolean = false;

    public loading: boolean = false;

    public policy!: PasswordComplexityPolicy.AsObject;
    public copied: string = '';

    public ChangeType: any = ChangeType;
    public userLoginMustBeDomain: boolean = false;

    constructor(
        public translate: TranslateService,
        private toast: ToastService,
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

            this.passwordForm = this.fb.group({
                currentPassword: ['', []],
                newPassword: ['', validators],
                confirmPassword: ['', [...validators, passwordConfirmValidator]],
            });
        }).catch(error => {
            this.toast.showError(error.message);
            console.error(error.message);
            this.passwordForm = this.fb.group({
                currentPassword: ['', []],
                newPassword: ['', validators],
                confirmPassword: ['', [...validators, passwordConfirmValidator]],
            });
        });

        // this.addressForm = this.fb.group({
        //     streetAddress: [''],
        //     postalCode: [''],
        //     locality: [''],
        //     region: [''],
        //     country: [''],
        // });

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
        this.user.firstName = profileData.firstName;
        this.user.lastName = profileData.lastName;
        this.user.nickName = profileData.nickName;
        this.user.displayName = profileData.displayName;
        this.user.gender = profileData.gender;
        this.user.preferredLanguage = profileData.preferredLanguage;
        this.userService
            .SaveMyUserProfile(
                this.user.firstName,
                this.user.lastName,
                this.user.nickName,
                this.user.preferredLanguage,
                this.user.gender,
            )
            .then((data: UserProfile) => {
                this.toast.showInfo('Saved Profile');
                this.user = Object.assign(this.user, data.toObject());
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
                }).catch(data => {
                    this.toast.showError(data.message);
                });
        }
    }

    public saveEmail(): void {
        this.emailEditState = false;

        this.userService
            .SaveMyUserEmail(this.user.email).then((data: UserEmail) => {
                this.toast.showInfo('Saved Email');
                this.user.email = data.toObject().email;
                this.emailEditState = false;
            }).catch(data => {
                this.toast.showError(data.message);
                this.emailEditState = false;
            });
    }

    public deletePhone(): void {
        this.user.phone = '';
        this.savePhone();
    }

    public enterCode(): void {
        const dialogRef = this.dialog.open(CodeDialogComponent, {
            data: {
                number: this.user.phone,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(code => {
            if (code) {
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
        this.userService.ResendEmailVerification().then(() => {
            this.toast.showInfo('Saved Email');
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public resendPhoneVerification(): void {
        this.userService.ResendPhoneVerification().then(() => {
            this.toast.showInfo('Phoneverification was successfully sent!');
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public savePhone(): void {
        this.phoneEditState = false;
        this.userService
            .SaveMyUserPhone(this.user.phone).then((data: UserPhone) => {
                this.toast.showInfo('Saved Phone');
                this.user.phone = data.toObject().phone;
                this.phoneEditState = false;
            }).catch(data => {
                this.toast.showError(data.message);
                this.phoneEditState = false;
            });
    }

    // public saveAddress(): void {
    //     this.address = this.addressForm.value;
    //     this.userService
    //         .SaveMyUserAddress(this.address as UserAddress.AsObject).then((data: UserAddress) => {
    //             this.toast.showInfo('Saved Address');
    //             this.address = data.toObject();
    //         }).catch(data => {
    //             this.toast.showError(data.message);
    //         });
    // }

    // public get streetAddress(): AbstractControl | null {
    //     return this.addressForm.get('streetAddress');
    // }
    // public get postalCode(): AbstractControl | null {
    //     return this.addressForm.get('postalCode');
    // }
    // public get locality(): AbstractControl | null {
    //     return this.addressForm.get('locality');
    // }
    // public get region(): AbstractControl | null {
    //     return this.addressForm.get('region');
    // }
    // public get country(): AbstractControl | null {
    //     return this.addressForm.get('country');
    // }

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
        this.userService.GetMyUser().then(user => {
            this.user = user.toObject();
        }).catch(err => {
            console.error(err);
        });

        // this.address = (await this.userService.GetMyUserAddress()).toObject();
        // this.addressForm.patchValue(this.address);
    }

    public copytoclipboard(value: string): void {
        const selBox = document.createElement('textarea');
        selBox.style.position = 'fixed';
        selBox.style.left = '0';
        selBox.style.top = '0';
        selBox.style.opacity = '0';
        selBox.value = value;
        document.body.appendChild(selBox);
        selBox.focus();
        selBox.select();
        document.execCommand('copy');
        document.body.removeChild(selBox);
        this.copied = value;
        setTimeout(() => {
            this.copied = '';
        }, 3000);
    }
}

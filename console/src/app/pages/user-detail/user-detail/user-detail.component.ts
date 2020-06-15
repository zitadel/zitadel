import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import {
    Gender,
    NotificationType,
    PasswordComplexityPolicy,
    UserAddress,
    UserEmail,
    UserPhone,
    UserProfile,
} from 'src/app/proto/generated/management_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

import { CodeDialogComponent } from '../code-dialog/code-dialog.component';

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
    selector: 'app-user-detail',
    templateUrl: './user-detail.component.html',
    styleUrls: ['./user-detail.component.scss'],
})
export class UserDetailComponent implements OnInit, OnDestroy {
    public profile!: UserProfile.AsObject;
    public email!: UserEmail.AsObject;
    public phone!: UserPhone.AsObject;
    public address!: UserAddress.AsObject;
    public genders: Gender[] = [Gender.GENDER_MALE, Gender.GENDER_FEMALE, Gender.GENDER_DIVERSE];
    public languages: string[] = ['de', 'en'];

    public passwordForm!: FormGroup;
    public addressForm!: FormGroup;

    public isMgmt: boolean = false;
    private subscription: Subscription = new Subscription();
    public emailEditState: boolean = false;
    public phoneEditState: boolean = false;

    public ChangeType: any = ChangeType;
    public loading: boolean = false;


    public minLengthPassword: any = {
        value: 0,
    };
    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private mgmtUserService: MgmtUserService,
        private fb: FormBuilder,
        private _location: Location,
        private dialog: MatDialog,
        private orgService: OrgService,
        public authUserService: AuthUserService,
    ) {
        const validators: Validators[] = [Validators.required];
        this.orgService.GetPasswordComplexityPolicy().then(data => {
            const policy: PasswordComplexityPolicy.AsObject = data.toObject();
            this.minLengthPassword.value = data.toObject().minLength;
            if (policy.minLength) {
                validators.push(Validators.minLength(policy.minLength));
            }
            if (policy.hasLowercase) {
                validators.push(Validators.pattern(/[a-z]/g));
            }
            if (policy.hasUppercase) {
                validators.push(Validators.pattern(/[A-Z]/g));
            }
            if (policy.hasNumber) {
                validators.push(Validators.pattern(/[0-9]/g));
            }
            if (policy.hasSymbol) {
                // All characters that are not a digit or an English letter \W or a whitespace \S
                validators.push(Validators.pattern(/[\W\S]/));
            }

            this.passwordForm = this.fb.group({
                password: ['', validators],
                confirmPassword: ['', [...validators, passwordConfirmValidator]],
            });
            // TODO custom validator for pattern
        }).catch(error => {
            console.log('no password complexity policy defined!');
            this.passwordForm = this.fb.group({
                password: ['', []],
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
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => {
            this.loading = true;
            this.getData(params).then(() => {
                this.loading = false;
            }).catch(error => {
                this.loading = false;
            });
        });
    }

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
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
        });

        dialogRef.afterClosed().subscribe(code => {
            if (code) {
                this.toast.showInfo('TODO: implement service');
            }
        });
    }

    public saveProfile(profileData: UserProfile.AsObject): void {
        this.profile.firstName = profileData.firstName;
        this.profile.lastName = profileData.lastName;
        this.profile.nickName = profileData.nickName;
        this.profile.displayName = profileData.displayName;
        this.profile.gender = profileData.gender;
        this.profile.preferredLanguage = profileData.preferredLanguage;
        this.mgmtUserService
            .SaveUserProfile(this.profile)
            .then((data: UserProfile) => {
                this.toast.showInfo('Saved Profile');
                this.profile = data.toObject();
            })
            .catch(data => {
                this.toast.showError(data.message);
            });
    }

    public resendVerification(): void {
        console.log('resendverification');
        this.mgmtUserService.ResendEmailVerification(this.profile.id).then((data: any) => {
            this.toast.showInfo('Email was successfully sent!');
            this.email = data.toObject();
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public resendPhoneVerification(): void {
        this.mgmtUserService.ResendPhoneVerification(this.profile.id).then((data: any) => {
            this.toast.showInfo('Phoneverification was successfully sent!');
            this.email = data.toObject();
        }).catch(data => {
            this.toast.showError(data.message);
        });
    }

    public setInitialPassword(): void {
        if (this.passwordForm.valid && this.password && this.password.value) {
            this.mgmtUserService.SetInitialPassword(this.profile.id, this.password.value).then((data: any) => {
                this.toast.showInfo('Set initial Password');
                this.email = data.toObject();
            }).catch(data => {
                this.toast.showError(data.message);
            });
        }
    }

    public sendSetPasswordNotification(): void {
        this.mgmtUserService.SendSetPasswordNotification(this.profile.id, NotificationType.NOTIFICATIONTYPE_EMAIL)
            .then((data: any) => {
                this.toast.showInfo('Set initial Password');
                this.email = data.toObject();
            }).catch(data => {
                this.toast.showError(data.message);
            });
    }

    public saveEmail(): void {
        this.emailEditState = false;
        this.mgmtUserService
            .SaveUserEmail(this.email).then((data: UserEmail) => {
                this.toast.showInfo('Saved Email');
                this.email = data.toObject();
            }).catch(data => {
                this.toast.showError(data.message);
            });
    }

    public savePhone(): void {
        this.phoneEditState = false;
        this.mgmtUserService
            .SaveUserPhone(this.phone).then((data: UserPhone) => {
                this.toast.showInfo('Saved Phone');
                this.phone = data.toObject();
            }).catch(data => {
                this.toast.showError(data.message);
            });
    }

    public saveAddress(): void {
        this.address.streetAddress = this.streetAddress?.value;
        this.address.postalCode = this.postalCode?.value;
        this.address.locality = this.locality?.value;
        this.address.region = this.region?.value;
        this.address.country = this.country?.value;

        this.mgmtUserService
            .SaveUserAddress(this.address as UserAddress.AsObject).then((data: UserAddress) => {
                this.toast.showInfo('Saved Address');
                this.address = data.toObject();
            }).catch(data => {
                this.toast.showError(data.message);
            });
    }

    public navigateBack(): void {
        this._location.back();
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

    public get password(): AbstractControl | null {
        return this.passwordForm.get('password');
    }
    public get confirmPassword(): AbstractControl | null {
        return this.passwordForm.get('confirmPassword');
    }

    private async getData({ id }: Params): Promise<void> {
        this.isMgmt = true;
        this.profile = (await this.mgmtUserService.GetUserProfile(id)).toObject();
        this.email = (await this.mgmtUserService.GetUserEmail(id)).toObject();
        this.phone = (await this.mgmtUserService.GetUserPhone(id)).toObject();
        this.address = (await this.mgmtUserService.GetUserAddress(id)).toObject();
        this.addressForm.patchValue(this.address);
    }
}

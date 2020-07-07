import { Component, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { Gender, UserAddress, UserEmail, UserPhone, UserProfile, UserView } from 'src/app/proto/generated/auth_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { ToastService } from 'src/app/services/toast.service';

import { CodeDialogComponent } from './code-dialog/code-dialog.component';

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

    private subscription: Subscription = new Subscription();

    public emailEditState: boolean = false;
    public phoneEditState: boolean = false;

    public loading: boolean = false;

    public copied: string = '';

    public ChangeType: any = ChangeType;
    public userLoginMustBeDomain: boolean = false;

    constructor(
        public translate: TranslateService,
        private toast: ToastService,
        private userService: AuthUserService,
        private dialog: MatDialog,
    ) {
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

    public deletePhone(): void {
        // this.userService.rem;
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

    private async getData(): Promise<void> {
        this.userService.GetMyUser().then(user => {
            this.user = user.toObject();
        }).catch(err => {
            console.error(err);
        });
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

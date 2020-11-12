import { Component, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import {
    Gender,
    UserAddress,
    UserEmail,
    UserPhone,
    UserProfile,
    UserState,
    UserView,
} from 'src/app/proto/generated/auth_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

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

    public loading: boolean = false;

    public copied: string = '';

    public ChangeType: any = ChangeType;
    public userLoginMustBeDomain: boolean = false;
    public UserState: any = UserState;

    constructor(
        public translate: TranslateService,
        private toast: ToastService,
        public userService: GrpcAuthService,
        private dialog: MatDialog,
    ) {
        this.loading = true;
        this.userService.GetMyUser().then(user => {
            this.user = user.toObject();
            this.loading = false;
        }).catch(error => {
            this.toast.showError(error);
            this.loading = false;
        });
    }

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }

    public saveProfile(profileData: UserProfile.AsObject): void {
        if (this.user.human) {
            this.user.human.firstName = profileData.firstName;
            this.user.human.lastName = profileData.lastName;
            this.user.human.nickName = profileData.nickName;
            this.user.human.displayName = profileData.displayName;
            this.user.human.gender = profileData.gender;
            this.user.human.preferredLanguage = profileData.preferredLanguage;

            this.userService
                .SaveMyUserProfile(
                    this.user.human.firstName,
                    this.user.human.lastName,
                    this.user.human.nickName,
                    this.user.human.preferredLanguage,
                    this.user.human.gender,
                )
                .then((data: UserProfile) => {
                    this.toast.showInfo('USER.TOAST.SAVED', true);
                    this.user = Object.assign(this.user, data.toObject());
                })
                .catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public saveEmail(email: string): void {
        this.userService
            .SaveMyUserEmail(email).then((data: UserEmail) => {
                this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
                if (this.user.human) {
                    this.user.human.email = data.toObject().email;
                }
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public enteredPhoneCode(code: string): void {
        this.userService.VerifyMyUserPhone(code).then(() => {
            this.toast.showInfo('USER.TOAST.PHONESAVED', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public changedLanguage(language: string): void {
        this.translate.use(language);
    }

    public resendPhoneVerification(): void {
        this.userService.ResendPhoneVerification().then(() => {
            this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public resendEmailVerification(): void {
        this.userService.ResendMyEmailVerificationMail().then(() => {
            this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public deletePhone(): void {
        this.userService.RemoveMyUserPhone().then(() => {
            this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
            if (this.user.human) {
                this.user.human.phone = '';
            }
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public savePhone(phone: string): void {
        if (this.user.human) {
            this.userService
                .SaveMyUserPhone(phone).then((data: UserPhone) => {
                    this.toast.showInfo('USER.TOAST.PHONESAVED', true);
                    if (this.user.human) {
                        this.user.human.phone = data.toObject().phone;
                    }
                }).catch(error => {
                    this.toast.showError(error);
                });
        }
    }
}

import { Component, EventEmitter, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import {
    Gender,
    UserAddress,
    UserEmail,
    UserPhone,
    UserProfile,
    UserState,
    UserView,
} from 'src/app/proto/generated/zitadel/auth_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { EditDialogType } from '../user-detail/user-detail.component';
import { EditDialogComponent } from './edit-dialog/edit-dialog.component';

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

    public USERGRANTCONTEXT: UserGrantContext = UserGrantContext.USER;
    public refreshChanges$: EventEmitter<void> = new EventEmitter();

    constructor(
        public translate: TranslateService,
        private toast: ToastService,
        public userService: GrpcAuthService,
        private dialog: MatDialog,
    ) {
        this.loading = true;
        this.refreshUser();
    }

    refreshUser(): void {
        this.refreshChanges$.emit();
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
                    this.refreshChanges$.emit();
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
                    this.refreshUser();
                }
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public enteredPhoneCode(code: string): void {
        this.userService.VerifyMyUserPhone(code).then(() => {
            this.toast.showInfo('USER.TOAST.PHONESAVED', true);
            this.refreshUser();
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
            this.refreshChanges$.emit();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public resendEmailVerification(): void {
        this.userService.ResendMyEmailVerificationMail().then(() => {
            this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
            this.refreshChanges$.emit();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public deletePhone(): void {
        this.userService.RemoveMyUserPhone().then(() => {
            this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
            if (this.user.human) {
                this.user.human.phone = '';
                this.refreshUser();
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
                        this.refreshUser();
                    }
                }).catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public openEditDialog(type: EditDialogType): void {
        switch (type) {
            case EditDialogType.PHONE:
                const dialogRefPhone = this.dialog.open(EditDialogComponent, {
                    data: {
                        confirmKey: 'ACTIONS.SAVE',
                        cancelKey: 'ACTIONS.CANCEL',
                        labelKey: 'ACTIONS.NEWVALUE',
                        titleKey: 'USER.LOGINMETHODS.PHONE.EDITTITLE',
                        descriptionKey: 'USER.LOGINMETHODS.PHONE.EDITDESC',
                        value: this.user.human?.phone,
                    },
                    width: '400px',
                });

                dialogRefPhone.afterClosed().subscribe(resp => {
                    if (resp) {
                        this.savePhone(resp);
                    }
                });
                break;
            case EditDialogType.EMAIL:
                const dialogRefEmail = this.dialog.open(EditDialogComponent, {
                    data: {
                        confirmKey: 'ACTIONS.SAVE',
                        cancelKey: 'ACTIONS.CANCEL',
                        labelKey: 'ACTIONS.NEWVALUE',
                        titleKey: 'USER.LOGINMETHODS.EMAIL.EDITTITLE',
                        descriptionKey: 'USER.LOGINMETHODS.EMAIL.EDITDESC',
                        value: this.user.human?.email,
                    },
                    width: '400px',
                });

                dialogRefEmail.afterClosed().subscribe(resp => {
                    if (resp) {
                        this.saveEmail(resp);
                    }
                });
                break;
        }
    }
}

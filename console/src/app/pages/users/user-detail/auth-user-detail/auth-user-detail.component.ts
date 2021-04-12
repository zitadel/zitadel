import { Component, EventEmitter, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { Email, Gender, Phone, Profile, User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { EditDialogComponent, EditDialogType } from './edit-dialog/edit-dialog.component';

@Component({
    selector: 'app-auth-user-detail',
    templateUrl: './auth-user-detail.component.html',
    styleUrls: ['./auth-user-detail.component.scss'],
})
export class AuthUserDetailComponent implements OnDestroy {
    public user!: User.AsObject;
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
        this.userService.getMyUser().then(resp => {
            if (resp.user) {
                this.user = resp.user;
            }
            this.loading = false;
        }).catch(error => {
            this.toast.showError(error);
            this.loading = false;
        });
    }

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }

    public saveProfile(profileData: Profile.AsObject): void {
        if (this.user.human) {
            this.user.human.profile = profileData;

            this.userService
                .updateMyProfile(
                    this.user.human.profile?.firstName,
                    this.user.human.profile?.lastName,
                    this.user.human.profile?.nickName,
                    this.user.human.profile?.displayName,
                    this.user.human.profile?.preferredLanguage,
                    this.user.human.profile?.gender,
                )
                .then(() => {
                    this.toast.showInfo('USER.TOAST.SAVED', true);
                    this.refreshChanges$.emit();
                })
                .catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public saveEmail(email: string): void {
        this.userService
            .setMyEmail(email).then(() => {
                this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
                if (this.user.human) {
                    const mailToSet = new Email();
                    mailToSet.setEmail(email);
                    this.user.human.email = mailToSet.toObject();
                    this.refreshUser();
                }
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public enteredPhoneCode(code: string): void {
        this.userService.verifyMyPhone(code).then(() => {
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
        this.userService.resendMyPhoneVerification().then(() => {
            this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
            this.refreshChanges$.emit();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public resendEmailVerification(): void {
        this.userService.resendMyEmailVerification().then(() => {
            this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
            this.refreshChanges$.emit();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public deletePhone(): void {
        this.userService.removeMyPhone().then(() => {
            this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
            if (this.user.human?.phone) {
                const phone = new Phone();
                this.user.human.phone = phone.toObject();
                this.refreshUser();
            }
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public savePhone(phone: string): void {
        if (this.user.human) {
            this.userService
                .setMyPhone(phone).then(() => {
                    this.toast.showInfo('USER.TOAST.PHONESAVED', true);
                    if (this.user.human) {
                        const phoneToSet = new Phone();
                        phoneToSet.setPhone(phone);
                        this.user.human.phone = phoneToSet.toObject();
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
                        value: this.user.human?.phone?.phone,
                        type: type,
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
                        value: this.user.human?.email?.email,
                        type: type,
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

import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import {
    Gender,
    MachineResponse,
    MachineView,
    NotificationType,
    UserEmail,
    UserPhone,
    UserProfile,
    UserState,
    UserView,
} from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-user-detail',
    templateUrl: './user-detail.component.html',
    styleUrls: ['./user-detail.component.scss'],
})
export class UserDetailComponent implements OnInit, OnDestroy {
    public user!: UserView.AsObject;
    public genders: Gender[] = [Gender.GENDER_MALE, Gender.GENDER_FEMALE, Gender.GENDER_DIVERSE];
    public languages: string[] = ['de', 'en'];

    private subscription: Subscription = new Subscription();

    public ChangeType: any = ChangeType;
    public loading: boolean = false;

    public UserState: any = UserState;
    public copied: string = '';
    public USERGRANTCONTEXT: UserGrantContext = UserGrantContext.USER;

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        public mgmtUserService: ManagementService,
        private _location: Location,
        private dialog: MatDialog,
    ) { }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => {
            const { id } = params;
            this.mgmtUserService.GetUserByID(id).then(user => {
                this.user = user.toObject();
            }).catch(err => {
                console.error(err);
            });
        });
    }

    public ngOnDestroy(): void {
        this.subscription.unsubscribe();
    }

    public changeState(newState: UserState): void {
        if (newState === UserState.USERSTATE_ACTIVE) {
            this.mgmtUserService.ReactivateUser(this.user.id).then(() => {
                this.toast.showInfo('USER.TOAST.REACTIVATED', true);
                this.user.state = newState;
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (newState === UserState.USERSTATE_INACTIVE) {
            this.mgmtUserService.DeactivateUser(this.user.id).then(() => {
                this.toast.showInfo('USER.TOAST.DEACTIVATED', true);
                this.user.state = newState;
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public saveProfile(profileData: UserProfile.AsObject): void {
        if (this.user.human) {
            this.user.human.firstName = profileData.firstName;
            this.user.human.lastName = profileData.lastName;
            this.user.human.nickName = profileData.nickName;
            this.user.human.displayName = profileData.displayName;
            this.user.human.gender = profileData.gender;
            this.user.human.preferredLanguage = profileData.preferredLanguage;
            this.mgmtUserService
                .SaveUserProfile(
                    this.user.id,
                    this.user.human.firstName,
                    this.user.human.lastName,
                    this.user.human.nickName,
                    this.user.human.preferredLanguage,
                    this.user.human.gender)
                .then((data: UserProfile) => {
                    this.toast.showInfo('USER.TOAST.SAVED', true);
                    this.user = Object.assign(this.user, data.toObject());
                })
                .catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public saveMachine(machineData: MachineView.AsObject): void {
        if (this.user.machine) {
            this.user.machine.name = machineData.name;
            this.user.machine.description = machineData.description;

            this.mgmtUserService
                .UpdateUserMachine(
                    this.user.id,
                    this.user.machine.description)
                .then((data: MachineResponse) => {
                    this.toast.showInfo('USER.TOAST.SAVED', true);
                    this.user = Object.assign(this.user, data.toObject());
                })
                .catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public resendEmailVerification(): void {
        this.mgmtUserService.ResendEmailVerification(this.user.id).then(() => {
            this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public resendPhoneVerification(): void {
        console.log('resend phone ver', this.user.id);
        this.mgmtUserService.ResendPhoneVerification(this.user.id).then(() => {
            this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public deletePhone(): void {
        this.mgmtUserService.RemoveUserPhone(this.user.id).then(() => {
            this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
            if (this.user.human) {
                this.user.human.phone = '';
            }
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public saveEmail(email: string): void {
        if (this.user.id && email) {
            this.mgmtUserService.SaveUserEmail(this.user.id, email).then((data: UserEmail) => {
                this.toast.showInfo('USER.TOAST.EMAILSENT', true);
                if (this.user.human) {
                    this.user.human.email = data.toObject().email;
                }
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public savePhone(phone: string): void {
        if (this.user.id && phone) {
            this.mgmtUserService
                .SaveUserPhone(this.user.id, phone).then((data: UserPhone) => {
                    this.toast.showInfo('USER.TOAST.PHONESAVED', true);
                    if (this.user.human) {
                        this.user.human.phone = data.toObject().phone;
                    }
                }).catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public navigateBack(): void {
        this._location.back();
    }

    public sendSetPasswordNotification(): void {
        this.mgmtUserService.SendSetPasswordNotification(this.user.id, NotificationType.NOTIFICATIONTYPE_EMAIL)
            .then(() => {
                this.toast.showInfo('USER.TOAST.PASSWORDNOTIFICATIONSENT', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public deleteUser(): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.DIALOG.DELETE_TITLE',
                descriptionKey: 'USER.DIALOG.DELETE_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                this.mgmtUserService.DeleteUser(this.user.id).then(() => {
                    this.navigateBack();
                    this.toast.showInfo('USER.TOAST.DELETED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }
}

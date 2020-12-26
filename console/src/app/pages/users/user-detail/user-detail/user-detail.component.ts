import { Location } from '@angular/common';
import { Component, EventEmitter, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
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
} from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { EditDialogComponent } from '../auth-user-detail/edit-dialog/edit-dialog.component';
import { ResendEmailDialogComponent } from '../auth-user-detail/resend-email-dialog/resend-email-dialog.component';

export enum EditDialogType {
    PHONE = 1,
    EMAIL = 2,
}

@Component({
    selector: 'app-user-detail',
    templateUrl: './user-detail.component.html',
    styleUrls: ['./user-detail.component.scss'],
})
export class UserDetailComponent implements OnInit {
    public user!: UserView.AsObject;
    public genders: Gender[] = [Gender.GENDER_MALE, Gender.GENDER_FEMALE, Gender.GENDER_DIVERSE];
    public languages: string[] = ['de', 'en'];

    public ChangeType: any = ChangeType;
    public loading: boolean = false;

    public UserState: any = UserState;
    public copied: string = '';
    public USERGRANTCONTEXT: UserGrantContext = UserGrantContext.USER;

    public EditDialogType: any = EditDialogType;
    public refreshChanges$: EventEmitter<void> = new EventEmitter();

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        public mgmtUserService: ManagementService,
        private _location: Location,
        private dialog: MatDialog,
        private router: Router,
    ) { }

    refreshUser(): void {
        this.refreshChanges$.emit();
        this.route.params.pipe(take(1)).subscribe(params => {
            const { id } = params;
            this.mgmtUserService.GetUserByID(id).then(user => {
                this.user = user.toObject();
            }).catch(err => {
                console.error(err);
            });
        });
    }

    public ngOnInit(): void {
        this.refreshUser();
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
                    this.refreshChanges$.emit();
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
                    this.refreshChanges$.emit();
                })
                .catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public resendEmailVerification(): void {
        this.mgmtUserService.ResendEmailVerification(this.user.id).then(() => {
            this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
            this.refreshChanges$.emit();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public resendPhoneVerification(): void {
        this.mgmtUserService.ResendPhoneVerification(this.user.id).then(() => {
            this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
            this.refreshChanges$.emit();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public deletePhone(): void {
        this.mgmtUserService.RemoveUserPhone(this.user.id).then(() => {
            this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
            if (this.user.human) {
                this.user.human.phone = '';
                this.refreshUser();
            }
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public saveEmail(email: string): void {
        if (this.user.id && email) {
            this.mgmtUserService.SaveUserEmail(this.user.id, email).then((data: UserEmail) => {
                this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
                if (this.user.human) {
                    this.user.human.email = data.toObject().email;
                    this.refreshUser();
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
                        this.refreshUser();
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
                this.refreshChanges$.emit();
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
                    const params: Params = {
                        'deferredReload': true,
                    };
                    this.router.navigate(['/users/list', this.user.human ? 'humans' : 'machines'], { queryParams: params });
                    this.toast.showInfo('USER.TOAST.DELETED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    public resendInitEmail(): void {
        const dialogRef = this.dialog.open(ResendEmailDialogComponent, {
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp.send && this.user.id) {
                this.mgmtUserService.ResendInitialMail(this.user.id, resp.email ?? '').then(() => {
                    this.toast.showInfo('USER.TOAST.INITEMAILSENT', true);
                    this.refreshChanges$.emit();
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
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

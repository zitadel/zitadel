import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
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
    public emailEditState: boolean = false;
    public phoneEditState: boolean = false;

    public ChangeType: any = ChangeType;
    public loading: boolean = false;

    public UserState: any = UserState;
    public copied: string = '';

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private mgmtUserService: ManagementService,
        private _location: Location,
        public mgmtService: ManagementService,
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

    public resendVerification(): void {
        this.mgmtUserService.ResendEmailVerification(this.user.id).then(() => {
            this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public resendPhoneVerification(): void {
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
            this.phoneEditState = false;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public saveEmail(): void {
        this.emailEditState = false;
        if (this.user && this.user.human?.email) {
            this.mgmtUserService
                .SaveUserEmail(this.user.id, this.user.human.email).then((data: UserEmail) => {
                    this.toast.showInfo('USER.TOAST.EMAILSENT', true);
                    if (this.user.human) {
                        this.user.human.email = data.toObject().email;
                    }
                }).catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public savePhone(): void {
        this.phoneEditState = false;
        if (this.user && this.user.human?.phone) {
            this.mgmtUserService
                .SaveUserPhone(this.user.id, this.user.human.phone).then((data: UserPhone) => {
                    this.toast.showInfo('USER.TOAST.PHONESAVED', true);
                    if (this.user.human) {
                        this.user.human.phone = data.toObject().phone;
                    }
                    this.phoneEditState = false;
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
}

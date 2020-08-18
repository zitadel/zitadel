import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import {
    Gender,
    NotificationType,
    UserEmail,
    UserPhone,
    UserProfile,
    UserState,
    UserView,
} from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ProjectService } from 'src/app/services/project.service';
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
        private mgmtUserService: MgmtUserService,
        private _location: Location,
        public projectService: ProjectService,
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
        this.user.firstName = profileData.firstName;
        this.user.lastName = profileData.lastName;
        this.user.nickName = profileData.nickName;
        this.user.displayName = profileData.displayName;
        this.user.gender = profileData.gender;
        this.user.preferredLanguage = profileData.preferredLanguage;
        this.mgmtUserService
            .SaveUserProfile(
                this.user.id,
                this.user.firstName,
                this.user.lastName,
                this.user.nickName,
                this.user.preferredLanguage,
                this.user.gender)
            .then((data: UserProfile) => {
                this.toast.showInfo('USER.TOAST.SAVED', true);
                this.user = Object.assign(this.user, data.toObject());
            })
            .catch(error => {
                this.toast.showError(error);
            });
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
            this.user.phone = '';
            this.phoneEditState = false;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public saveEmail(): void {
        this.emailEditState = false;
        this.mgmtUserService
            .SaveUserEmail(this.user.id, this.user.email).then((data: UserEmail) => {
                this.toast.showInfo('USER.TOAST.EMAILSENT', true);
                this.user.email = data.toObject().email;
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public savePhone(): void {
        this.phoneEditState = false;
        this.mgmtUserService
            .SaveUserPhone(this.user.id, this.user.phone).then((data: UserPhone) => {
                this.toast.showInfo('USER.TOAST.PHONESAVED', true);
                this.user.phone = data.toObject().phone;
                this.phoneEditState = false;
            }).catch(error => {
                this.toast.showError(error);
            });
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

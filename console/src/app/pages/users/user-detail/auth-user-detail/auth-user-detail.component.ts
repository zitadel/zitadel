import { MediaMatcher } from '@angular/cdk/layout';
import { Component, EventEmitter, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Email, Gender, Phone, Profile, User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { BreadcrumbService } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { EditDialogComponent, EditDialogType } from './edit-dialog/edit-dialog.component';

interface UserSetting {
  id: string;
  i18nKey: string;
  featureRequired: string[] | false;
}

@Component({
  selector: 'cnsl-auth-user-detail',
  templateUrl: './auth-user-detail.component.html',
  styleUrls: ['./auth-user-detail.component.scss'],
})
export class AuthUserDetailComponent implements OnDestroy {
  public user!: User.AsObject;
  public genders: Gender[] = [Gender.GENDER_MALE, Gender.GENDER_FEMALE, Gender.GENDER_DIVERSE];
  public languages: string[] = ['de', 'en', 'it'];

  private subscription: Subscription = new Subscription();

  public loading: boolean = false;

  public ChangeType: any = ChangeType;
  public userLoginMustBeDomain: boolean = false;
  public UserState: any = UserState;

  public USERGRANTCONTEXT: UserGrantContext = UserGrantContext.USER;
  public refreshChanges$: EventEmitter<void> = new EventEmitter();

  public settingsList: UserSetting[] = [
    { id: 'general', i18nKey: 'USER.SETTINGS.GENERAL', featureRequired: false },
    { id: 'idp', i18nKey: 'USER.SETTINGS.IDP', featureRequired: false },
    { id: 'passwordless', i18nKey: 'USER.SETTINGS.PASSWORDLESS', featureRequired: false },
    { id: 'mfa', i18nKey: 'USER.SETTINGS.MFA', featureRequired: false },
    { id: 'grants', i18nKey: 'USER.SETTINGS.USERGRANTS', featureRequired: false },
    { id: 'memberships', i18nKey: 'USER.SETTINGS.MEMBERSHIPS', featureRequired: false },
    { id: 'metadata', i18nKey: 'USER.SETTINGS.METADATA', featureRequired: ['metadata.user'] },
  ];
  public currentSetting: UserSetting | undefined = this.settingsList[0];

  constructor(
    public translate: TranslateService,
    private toast: ToastService,
    public userService: GrpcAuthService,
    private dialog: MatDialog,
    private auth: AuthenticationService,
    breadcrumbService: BreadcrumbService,
    private mediaMatcher: MediaMatcher,
  ) {
    const mediaq: string = '(max-width: 500px)';
    const small = this.mediaMatcher.matchMedia(mediaq).matches;
    if (small) {
      this.changeSelection(small);
    }
    this.mediaMatcher.matchMedia(mediaq).onchange = (small) => {
      this.changeSelection(small.matches);
    };

    breadcrumbService.setBreadcrumb([]);

    this.loading = true;
    this.refreshUser();

    this.userService.getSupportedLanguages().then((lang) => {
      this.languages = lang.languagesList;
    });
  }

  private changeSelection(small: boolean): void {
    if (small) {
      this.currentSetting = undefined;
    } else {
      this.currentSetting = this.currentSetting === undefined ? this.settingsList[0] : this.currentSetting;
    }
  }

  refreshUser(): void {
    this.refreshChanges$.emit();
    this.userService
      .getMyUser()
      .then((resp) => {
        if (resp.user) {
          this.user = resp.user;
        }
        this.loading = false;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }

  public changeUsername(): void {
    const dialogRef = this.dialog.open(EditDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CHANGE',
        cancelKey: 'ACTIONS.CANCEL',
        labelKey: 'ACTIONS.NEWVALUE',
        titleKey: 'USER.PROFILE.CHANGEUSERNAME_TITLE',
        descriptionKey: 'USER.PROFILE.CHANGEUSERNAME_DESC',
        value: this.user.userName,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp: { value: string }) => {
      if (resp && resp.value && resp.value !== this.user.userName) {
        this.userService
          .updateMyUserName(resp.value)
          .then(() => {
            this.toast.showInfo('USER.TOAST.USERNAMECHANGED', true);
            this.refreshUser();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
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
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public saveEmail(email: string): void {
    this.userService
      .setMyEmail(email)
      .then(() => {
        this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
        if (this.user.human) {
          const mailToSet = new Email();
          mailToSet.setEmail(email);
          this.user.human.email = mailToSet.toObject();
          this.refreshUser();
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public enteredPhoneCode(code: string): void {
    this.userService
      .verifyMyPhone(code)
      .then(() => {
        this.toast.showInfo('USER.TOAST.PHONESAVED', true);
        this.refreshUser();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public changedLanguage(language: string): void {
    this.translate.use(language);
  }

  public resendPhoneVerification(): void {
    this.userService
      .resendMyPhoneVerification()
      .then(() => {
        this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public resendEmailVerification(): void {
    this.userService
      .resendMyEmailVerification()
      .then(() => {
        this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deletePhone(): void {
    this.userService
      .removeMyPhone()
      .then(() => {
        this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
        if (this.user.human?.phone) {
          const phone = new Phone();
          this.user.human.phone = phone.toObject();
          this.refreshUser();
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public savePhone(phone: string): void {
    if (this.user.human) {
      this.userService
        .setMyPhone(phone)
        .then(() => {
          this.toast.showInfo('USER.TOAST.PHONESAVED', true);
          if (this.user.human) {
            const phoneToSet = new Phone();
            phoneToSet.setPhone(phone);
            this.user.human.phone = phoneToSet.toObject();
            this.refreshUser();
          }
        })
        .catch((error) => {
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

        dialogRefPhone.afterClosed().subscribe((resp: { value: string; isVerified: boolean }) => {
          if (resp && resp.value) {
            this.savePhone(resp.value);
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

        dialogRefEmail.afterClosed().subscribe((resp: { value: string; isVerified: boolean }) => {
          if (resp && resp.value) {
            this.saveEmail(resp.value);
          }
        });
        break;
    }
  }

  public deleteAccount(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'USER.DIALOG.DELETE_BTN',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.DIALOG.DELETE_TITLE',
        descriptionKey: 'USER.DIALOG.DELETE_AUTH_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.userService
          .RemoveMyUser()
          .then(() => {
            this.toast.showInfo('USER.PAGES.DELETEACCOUNT_SUCCESS', true);
            this.auth.signout();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }
}

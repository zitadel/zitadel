import { MediaMatcher } from '@angular/cdk/layout';
import { Location } from '@angular/common';
import { Component, EventEmitter, OnInit } from '@angular/core';
import { Validators } from '@angular/forms';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Buffer } from 'buffer';
import { take } from 'rxjs/operators';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { phoneValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { MetadataDialogComponent } from 'src/app/modules/metadata/metadata-dialog/metadata-dialog.component';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { SendHumanResetPasswordNotificationRequest, UnlockUserRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';
import { LoginPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { Email, Gender, Machine, Phone, Profile, User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { formatPhone } from 'src/app/utils/formatPhone';
import { EditDialogComponent, EditDialogType } from '../auth-user-detail/edit-dialog/edit-dialog.component';
import { ResendEmailDialogComponent } from '../auth-user-detail/resend-email-dialog/resend-email-dialog.component';
import { MachineSecretDialogComponent } from './machine-secret-dialog/machine-secret-dialog.component';

const GENERAL: SidenavSetting = { id: 'general', i18nKey: 'USER.SETTINGS.GENERAL' };
const GRANTS: SidenavSetting = { id: 'grants', i18nKey: 'USER.SETTINGS.USERGRANTS' };
const METADATA: SidenavSetting = { id: 'metadata', i18nKey: 'USER.SETTINGS.METADATA' };
const IDP: SidenavSetting = { id: 'idp', i18nKey: 'USER.SETTINGS.IDP' };
const SECURITY: SidenavSetting = { id: 'security', i18nKey: 'USER.SETTINGS.SECURITY' };
const PERSONALACCESSTOKEN: SidenavSetting = { id: 'pat', i18nKey: 'USER.SETTINGS.PAT' };
const KEYS: SidenavSetting = { id: 'keys', i18nKey: 'USER.SETTINGS.KEYS' };
const MEMBERSHIPS: SidenavSetting = { id: 'memberships', i18nKey: 'USER.SETTINGS.MEMBERSHIPS' };

@Component({
  selector: 'cnsl-user-detail',
  templateUrl: './user-detail.component.html',
  styleUrls: ['./user-detail.component.scss'],
})
export class UserDetailComponent implements OnInit {
  public user!: User.AsObject;
  public metadata: Metadata.AsObject[] = [];
  public genders: Gender[] = [Gender.GENDER_MALE, Gender.GENDER_FEMALE, Gender.GENDER_DIVERSE];
  public languages: string[] = ['de', 'en', 'it', 'fr', 'pl', 'zh'];

  public ChangeType: any = ChangeType;

  public loading: boolean = true;
  public loadingMetadata: boolean = true;

  public UserState: any = UserState;
  public copied: string = '';
  public USERGRANTCONTEXT: UserGrantContext = UserGrantContext.USER;

  public EditDialogType: any = EditDialogType;
  public refreshChanges$: EventEmitter<void> = new EventEmitter();
  public InfoSectionType: any = InfoSectionType;

  public error: string = '';

  public settingsList: SidenavSetting[] = [GENERAL, GRANTS, MEMBERSHIPS, METADATA];
  public currentSetting: string | undefined = 'general';
  public loginPolicy?: LoginPolicy.AsObject;

  constructor(
    public translate: TranslateService,
    private route: ActivatedRoute,
    private toast: ToastService,
    public mgmtUserService: ManagementService,
    private _location: Location,
    private dialog: MatDialog,
    private router: Router,
    activatedRoute: ActivatedRoute,
    private mediaMatcher: MediaMatcher,
    breadcrumbService: BreadcrumbService,
  ) {
    activatedRoute.queryParams.pipe(take(1)).subscribe((params: Params) => {
      const { id } = params;
      if (id) {
        this.currentSetting = id;
      }
    });

    breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);

    const mediaq: string = '(max-width: 500px)';
    const small = this.mediaMatcher.matchMedia(mediaq).matches;
    if (small) {
      this.changeSelection(small);
    }
    this.mediaMatcher.matchMedia(mediaq).onchange = (small) => {
      this.changeSelection(small.matches);
    };

    this.mgmtUserService.getSupportedLanguages().then((lang) => {
      this.languages = lang.languagesList;
    });
  }

  private changeSelection(small: boolean): void {
    if (small) {
      this.currentSetting = undefined;
    } else {
      this.currentSetting = this.currentSetting === undefined ? 'general' : this.currentSetting;
    }
  }

  refreshUser(): void {
    this.refreshChanges$.emit();
    this.route.params.pipe(take(1)).subscribe((params) => {
      this.loading = true;
      const { id } = params;
      this.mgmtUserService
        .getUserByID(id)
        .then((resp) => {
          this.loadMetadata(id);
          this.loading = false;
          if (resp.user) {
            this.user = resp.user;

            if (this.user.human) {
              this.settingsList = [GENERAL, SECURITY, IDP, GRANTS, MEMBERSHIPS, METADATA];
            } else if (this.user.machine) {
              this.settingsList = [GENERAL, GRANTS, MEMBERSHIPS, PERSONALACCESSTOKEN, KEYS, METADATA];
            }
          }
        })
        .catch((err) => {
          this.error = err.message ?? '';
          this.loading = false;
          this.toast.showError(err);
        });
    });
  }

  public ngOnInit(): void {
    this.refreshUser();

    this.mgmtUserService.getLoginPolicy().then((policy) => {
      if (policy.policy) {
        this.loginPolicy = policy.policy;
      }
    });
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
      if (resp.value && resp.value !== this.user.userName) {
        this.mgmtUserService
          .updateUserName(this.user.id, resp.value)
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

  public unlockUser(): void {
    const req = new UnlockUserRequest();
    req.setId(this.user.id);
    this.mgmtUserService
      .unlockUser(req)
      .then(() => {
        this.toast.showInfo('USER.TOAST.UNLOCKED', true);
        this.refreshUser();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public generateMachineSecret(): void {
    this.mgmtUserService
      .generateMachineSecret(this.user.id)
      .then((resp) => {
        this.toast.showInfo('USER.TOAST.SECRETGENERATED', true);
        console.log(resp.clientSecret);
        this.dialog.open(MachineSecretDialogComponent, {
          data: {
            clientId: resp.clientId,
            clientSecret: resp.clientSecret,
          },
          width: '400px',
        });
        this.refreshUser();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public removeMachineSecret(): void {
    this.mgmtUserService
      .removeMachineSecret(this.user.id)
      .then((resp) => {
        this.toast.showInfo('USER.TOAST.SECRETREMOVED', true);
        this.refreshUser();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public changeState(newState: UserState): void {
    if (newState === UserState.USER_STATE_ACTIVE) {
      this.mgmtUserService
        .reactivateUser(this.user.id)
        .then(() => {
          this.toast.showInfo('USER.TOAST.REACTIVATED', true);
          this.user.state = newState;
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (newState === UserState.USER_STATE_INACTIVE) {
      this.mgmtUserService
        .deactivateUser(this.user.id)
        .then(() => {
          this.toast.showInfo('USER.TOAST.DEACTIVATED', true);
          this.user.state = newState;
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public saveProfile(profileData: Profile.AsObject): void {
    if (this.user.human) {
      this.user.human.profile = profileData;
      this.mgmtUserService
        .updateHumanProfile(
          this.user.id,
          this.user.human.profile.firstName,
          this.user.human.profile.lastName,
          this.user.human.profile.nickName,
          this.user.human.profile.displayName,
          this.user.human.profile.preferredLanguage,
          this.user.human.profile.gender,
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

  public saveMachine(machineData: Machine.AsObject): void {
    if (this.user.machine) {
      this.user.machine.name = machineData.name;
      this.user.machine.description = machineData.description;
      this.user.machine.accessTokenType = machineData.accessTokenType;

      this.mgmtUserService
        .updateMachine(
          this.user.id,
          this.user.machine.name,
          this.user.machine.description,
          this.user.machine.accessTokenType,
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

  public resendEmailVerification(): void {
    this.mgmtUserService
      .resendHumanEmailVerification(this.user.id)
      .then(() => {
        this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public resendPhoneVerification(): void {
    this.mgmtUserService
      .resendHumanPhoneVerification(this.user.id)
      .then(() => {
        this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deletePhone(): void {
    this.mgmtUserService
      .removeHumanPhone(this.user.id)
      .then(() => {
        this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
        if (this.user.human) {
          this.user.human.phone = new Phone().setPhone('').toObject();
          this.refreshUser();
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public saveEmail(email: string, isVerified: boolean): void {
    if (this.user.id && email) {
      this.mgmtUserService
        .updateHumanEmail(this.user.id, email, isVerified)
        .then(() => {
          this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
          if (this.user.state === UserState.USER_STATE_INITIAL) {
            this.mgmtUserService
              .resendHumanInitialization(this.user.id, email ?? '')
              .then(() => {
                this.toast.showInfo('USER.TOAST.INITEMAILSENT', true);
                this.refreshChanges$.emit();
              })
              .catch((error) => {
                this.toast.showError(error);
              });
          }
          if (this.user.human) {
            this.user.human.email = new Email().setEmail(email).toObject();
            this.refreshUser();
          }
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public savePhone(phone: string): void {
    if (this.user.id && phone) {
      // Format phone before save (add +)
      phone = formatPhone(phone).phone;

      this.mgmtUserService
        .updateHumanPhone(this.user.id, phone)
        .then(() => {
          this.toast.showInfo('USER.TOAST.PHONESAVED', true);
          if (this.user.human) {
            this.user.human.phone = new Phone().setPhone(phone).toObject();
            this.refreshUser();
          }
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public navigateBack(): void {
    this._location.back();
  }

  public sendSetPasswordNotification(): void {
    this.mgmtUserService
      .sendHumanResetPasswordNotification(this.user.id, SendHumanResetPasswordNotificationRequest.Type.TYPE_EMAIL)
      .then(() => {
        this.toast.showInfo('USER.TOAST.PASSWORDNOTIFICATIONSENT', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
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

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtUserService
          .removeUser(this.user.id)
          .then(() => {
            const params: Params = {
              deferredReload: true,
              type: this.user.human ? 'humans' : 'machines',
            };
            this.router.navigate(['/users'], { queryParams: params });
            this.toast.showInfo('USER.TOAST.DELETED', true);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public resendInitEmail(): void {
    const dialogRef = this.dialog.open(ResendEmailDialogComponent, {
      width: '400px',
      data: {
        email: this.user.human?.email?.email ?? '',
      },
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp.send && this.user.id) {
        this.mgmtUserService
          .resendHumanInitialization(this.user.id, resp.email ?? '')
          .then(() => {
            this.toast.showInfo('USER.TOAST.INITEMAILSENT', true);
            this.refreshChanges$.emit();
          })
          .catch((error) => {
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
            value: this.user.human?.phone?.phone,
            type: EditDialogType.PHONE,
            validator: Validators.compose([phoneValidator, requiredValidator]),
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
            isVerifiedTextKey: 'USER.LOGINMETHODS.EMAIL.ISVERIFIED',
            isVerifiedTextDescKey: 'USER.LOGINMETHODS.EMAIL.ISVERIFIEDDESC',
            value: this.user.human?.email?.email,
            type: EditDialogType.EMAIL,
          },
          width: '400px',
        });

        dialogRefEmail.afterClosed().subscribe((resp: { value: string; isVerified: boolean }) => {
          if (resp && resp.value) {
            this.saveEmail(resp.value, resp.isVerified);
          }
        });
        break;
    }
  }

  public loadMetadata(id: string): Promise<any> | void {
    this.loadingMetadata = true;
    return this.mgmtUserService
      .listUserMetadata(id)
      .then((resp) => {
        this.loadingMetadata = false;
        this.metadata = resp.resultList.map((md) => {
          return {
            key: md.key,
            value: Buffer.from(md.value as string, 'base64').toString('ascii'),
          };
        });
      })
      .catch((error) => {
        this.loadingMetadata = false;
        this.toast.showError(error);
      });
  }

  public editMetadata(): void {
    if (this.user) {
      const setFcn = (key: string, value: string): Promise<any> =>
        this.mgmtUserService.setUserMetadata(key, Buffer.from(value).toString('base64'), this.user.id);
      const removeFcn = (key: string): Promise<any> => this.mgmtUserService.removeUserMetadata(key, this.user.id);

      const dialogRef = this.dialog.open(MetadataDialogComponent, {
        data: {
          metadata: this.metadata,
          setFcn: setFcn,
          removeFcn: removeFcn,
        },
      });

      dialogRef.afterClosed().subscribe(() => {
        this.loadMetadata(this.user.id);
      });
    }
  }
}

import { Location } from '@angular/common';
import { Component, EventEmitter, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { SendHumanResetPasswordNotificationRequest, UnlockUserRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Email, Gender, Machine, Phone, Profile, User, UserState } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { EditDialogComponent, EditDialogType } from '../auth-user-detail/edit-dialog/edit-dialog.component';
import { ResendEmailDialogComponent } from '../auth-user-detail/resend-email-dialog/resend-email-dialog.component';

@Component({
  selector: 'app-user-detail',
  templateUrl: './user-detail.component.html',
  styleUrls: ['./user-detail.component.scss'],
})
export class UserDetailComponent implements OnInit {
  public user!: User.AsObject;
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
      this.mgmtUserService.getUserByID(id).then(resp => {
        if (resp.user) {
          this.user = resp.user;
        }
      }).catch(err => {
        console.error(err);
      });
    });
  }

  public ngOnInit(): void {
    this.refreshUser();
  }

  public unlockUser(): void {
    const req = new UnlockUserRequest();
    req.setId(this.user.id);
    this.mgmtUserService.unlockUser(req).then(() => {
      this.toast.showInfo('USER.TOAST.UNLOCKED', true);
      this.refreshUser();
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public changeState(newState: UserState): void {
    if (newState === UserState.USER_STATE_ACTIVE) {
      this.mgmtUserService.reactivateUser(this.user.id).then(() => {
        this.toast.showInfo('USER.TOAST.REACTIVATED', true);
        this.user.state = newState;
      }).catch(error => {
        this.toast.showError(error);
      });
    } else if (newState === UserState.USER_STATE_INACTIVE) {
      this.mgmtUserService.deactivateUser(this.user.id).then(() => {
        this.toast.showInfo('USER.TOAST.DEACTIVATED', true);
        this.user.state = newState;
      }).catch(error => {
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
          this.user.human.profile.gender)
        .then(() => {
          this.toast.showInfo('USER.TOAST.SAVED', true);
          this.refreshChanges$.emit();
        })
        .catch(error => {
          this.toast.showError(error);
        });
    }
  }

  public saveMachine(machineData: Machine.AsObject): void {
    if (this.user.machine) {
      this.user.machine.name = machineData.name;
      this.user.machine.description = machineData.description;

      this.mgmtUserService
        .updateMachine(
          this.user.id,
          this.user.machine.name,
          this.user.machine.description)
        .then(() => {
          this.toast.showInfo('USER.TOAST.SAVED', true);
          this.refreshChanges$.emit();
        })
        .catch(error => {
          this.toast.showError(error);
        });
    }
  }

  public resendEmailVerification(): void {
    this.mgmtUserService.resendHumanEmailVerification(this.user.id).then(() => {
      this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
      this.refreshChanges$.emit();
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public resendPhoneVerification(): void {
    this.mgmtUserService.resendHumanPhoneVerification(this.user.id).then(() => {
      this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
      this.refreshChanges$.emit();
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public deletePhone(): void {
    this.mgmtUserService.removeHumanPhone(this.user.id).then(() => {
      this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
      if (this.user.human) {
        this.user.human.phone = new Phone().setPhone('').toObject();
        this.refreshUser();
      }
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public saveEmail(email: string): void {
    if (this.user.id && email) {
      this.mgmtUserService.updateHumanEmail(this.user.id, email).then(() => {
        this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
        if (this.user.state === UserState.USER_STATE_INITIAL) {
          this.mgmtUserService.resendHumanInitialization(this.user.id, email ?? '').then(() => {
            this.toast.showInfo('USER.TOAST.INITEMAILSENT', true);
            this.refreshChanges$.emit();
          }).catch(error => {
            this.toast.showError(error);
          });
        }
        if (this.user.human) {
          this.user.human.email = new Email().setEmail(email).toObject();
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
        .updateHumanPhone(this.user.id, phone).then(() => {
          this.toast.showInfo('USER.TOAST.PHONESAVED', true);
          if (this.user.human) {
            this.user.human.phone = new Phone().setPhone(phone).toObject();
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
    this.mgmtUserService.sendHumanResetPasswordNotification(
      this.user.id,
      SendHumanResetPasswordNotificationRequest.Type.TYPE_EMAIL,
    ).then(() => {
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
        this.mgmtUserService.removeUser(this.user.id).then(() => {
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
        this.mgmtUserService.resendHumanInitialization(this.user.id, resp.email ?? '').then(() => {
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
            value: this.user.human?.phone?.phone,
            type: EditDialogType.PHONE,
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
            type: EditDialogType.EMAIL,
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

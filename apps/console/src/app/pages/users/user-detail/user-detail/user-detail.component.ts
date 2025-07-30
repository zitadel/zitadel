import { Location } from '@angular/common';
import { Component, DestroyRef, EventEmitter, OnInit, signal } from '@angular/core';
import { Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { catchError, filter, map, startWith, take } from 'rxjs/operators';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { phoneValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import {
  MetadataDialogComponent,
  MetadataDialogData,
} from 'src/app/modules/metadata/metadata-dialog/metadata-dialog.component';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import { formatPhone } from 'src/app/utils/formatPhone';
import {
  EditDialogData,
  EditDialogResult,
  EditDialogComponent,
  EditDialogType,
} from '../auth-user-detail/edit-dialog/edit-dialog.component';
import {
  ResendEmailDialogComponent,
  ResendEmailDialogData,
  ResendEmailDialogResult,
} from '../auth-user-detail/resend-email-dialog/resend-email-dialog.component';
import { MachineSecretDialogComponent } from './machine-secret-dialog/machine-secret-dialog.component';
import { LanguagesService } from 'src/app/services/languages.service';
import { UserService } from 'src/app/services/user.service';
import { Gender, HumanProfile, HumanUser, User as UserV2, UserState } from '@zitadel/proto/zitadel/user/v2/user_pb';
import {
  combineLatestWith,
  defer,
  EMPTY,
  identity,
  mergeWith,
  Observable,
  ObservedValueOf,
  of,
  shareReplay,
  Subject,
  switchMap,
} from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { DetailFormMachineComponent } from '../detail-form-machine/detail-form-machine.component';
import { NewMgmtService } from 'src/app/services/new-mgmt.service';
import { LoginPolicy } from '@zitadel/proto/zitadel/policy_pb';
import { SendHumanResetPasswordNotificationRequest_Type } from '@zitadel/proto/zitadel/management_pb';
import { pairwiseStartWith } from 'src/app/utils/pairwiseStartWith';
import { Metadata } from '@zitadel/proto/zitadel/metadata_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

const GENERAL: SidenavSetting = { id: 'general', i18nKey: 'USER.SETTINGS.GENERAL' };
const GRANTS: SidenavSetting = { id: 'grants', i18nKey: 'USER.SETTINGS.USERGRANTS' };
const METADATA: SidenavSetting = { id: 'metadata', i18nKey: 'USER.SETTINGS.METADATA' };
const IDP: SidenavSetting = { id: 'idp', i18nKey: 'USER.SETTINGS.IDP' };
const SECURITY: SidenavSetting = { id: 'security', i18nKey: 'USER.SETTINGS.SECURITY' };
const PERSONALACCESSTOKEN: SidenavSetting = { id: 'pat', i18nKey: 'USER.SETTINGS.PAT' };
const KEYS: SidenavSetting = { id: 'keys', i18nKey: 'USER.SETTINGS.KEYS' };
const MEMBERSHIPS: SidenavSetting = { id: 'memberships', i18nKey: 'USER.SETTINGS.MEMBERSHIPS' };

type UserQuery =
  | { state: 'success'; value: UserV2 }
  | { state: 'error'; value: string }
  | { state: 'loading'; value?: UserV2 }
  | { state: 'notfound' };

type MetadataQuery =
  | { state: 'success'; value: Metadata[] }
  | { state: 'loading'; value: Metadata[] }
  | { state: 'error'; value: string };

type UserWithHumanType = Omit<UserV2, 'type'> & { type: { case: 'human'; value: HumanUser } };

// todo: figure out why media matcher is needed
@Component({
  selector: 'cnsl-user-detail',
  templateUrl: './user-detail.component.html',
  styleUrls: ['./user-detail.component.scss'],
})
export class UserDetailComponent implements OnInit {
  public user$: Observable<UserQuery>;
  public genders: Gender[] = [Gender.MALE, Gender.FEMALE, Gender.DIVERSE];

  public ChangeType: any = ChangeType;

  public UserState = UserState;
  public copied: string = '';
  public USERGRANTCONTEXT: UserGrantContext = UserGrantContext.USER;

  public EditDialogType: any = EditDialogType;
  public refreshChanges$: EventEmitter<void> = new EventEmitter();
  public InfoSectionType: any = InfoSectionType;

  public currentSetting$ = signal<SidenavSetting>(GENERAL);
  public settingsList$: Observable<SidenavSetting[]>;
  public metadata$: Observable<MetadataQuery>;
  public loginPolicy$: Observable<LoginPolicy>;
  public refreshMetadata$ = new Subject<true>();

  constructor(
    public translate: TranslateService,
    private readonly route: ActivatedRoute,
    private toast: ToastService,
    private _location: Location,
    private dialog: MatDialog,
    private router: Router,
    public langSvc: LanguagesService,
    private readonly userService: UserService,
    private readonly newMgmtService: NewMgmtService,
    public readonly mgmtService: ManagementService,
    breadcrumbService: BreadcrumbService,
    private readonly destroyRef: DestroyRef,
  ) {
    breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);

    this.user$ = this.getUser$().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.settingsList$ = this.getSettingsList$(this.user$).pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.metadata$ = this.getMetadata$(this.user$).pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.loginPolicy$ = defer(() => this.newMgmtService.getLoginPolicy()).pipe(
      catchError(() => EMPTY),
      map(({ policy }) => policy),
      filter(Boolean),
    );
  }

  private getId$(): Observable<string> {
    return this.route.paramMap.pipe(
      map((params) => params.get('id')),
      filter(Boolean),
    );
  }

  private getUser$(): Observable<UserQuery> {
    return this.getId$().pipe(
      combineLatestWith(this.refreshChanges$.pipe(startWith(undefined))),
      switchMap(([id]) => this.getUserById(id)),
      pairwiseStartWith(undefined),
      map(([prev, curr]) => {
        if (prev?.state === 'success' && curr.state === 'loading') {
          return { state: 'loading', value: prev.value } as const;
        }
        return curr;
      }),
    );
  }

  private getSettingsList$(user$: Observable<UserQuery>): Observable<SidenavSetting[]> {
    return user$.pipe(
      switchMap((user) => {
        if (user.state !== 'success') {
          return EMPTY;
        }

        if (user.value.type.case === 'human') {
          return of([GENERAL, SECURITY, IDP, GRANTS, MEMBERSHIPS, METADATA]);
        } else if (user.value.type.case === 'machine') {
          return of([GENERAL, GRANTS, MEMBERSHIPS, PERSONALACCESSTOKEN, KEYS, METADATA]);
        }
        return EMPTY;
      }),
      startWith([GENERAL, GRANTS, MEMBERSHIPS, METADATA]),
    );
  }

  private getUserById(userId: string): Observable<UserQuery> {
    return defer(() => this.userService.getUserById(userId)).pipe(
      map(({ user }) => {
        if (user) {
          return { state: 'success', value: user } as const;
        }
        return { state: 'notfound' } as const;
      }),
      catchError((error) => of({ state: 'error', value: error.message ?? '' } as const)),
      startWith({ state: 'loading' } as const),
    );
  }

  getMetadata$(user$: Observable<UserQuery>): Observable<MetadataQuery> {
    return this.refreshMetadata$.pipe(
      startWith(true),
      combineLatestWith(user$),
      switchMap(([_, user]) => {
        if (!(user.state === 'success' || user.state === 'loading')) {
          return EMPTY;
        }
        if (!user.value) {
          return EMPTY;
        }
        return this.getMetadataById(user.value.userId);
      }),
      pairwiseStartWith(undefined),
      map(([prev, curr]) => {
        if (prev?.state === 'success' && curr.state === 'loading') {
          return { state: 'loading', value: prev.value } as const;
        }
        return curr;
      }),
    );
  }

  public ngOnInit(): void {
    this.user$.pipe(mergeWith(this.metadata$), takeUntilDestroyed(this.destroyRef)).subscribe((query) => {
      if (query.state == 'error') {
        this.toast.showError(query.value);
      }
    });

    const param = this.route.snapshot.queryParamMap.get('id');
    if (!param) {
      return;
    }

    this.settingsList$
      .pipe(
        takeUntilDestroyed(this.destroyRef),
        map((settings) => settings.find(({ id }) => id === param)),
        filter(Boolean),
        take(1),
      )
      .subscribe((setting) => this.currentSetting$.set(setting));
  }

  public changeUsername(user: UserV2): void {
    const dialogRef = this.dialog.open(EditDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CHANGE',
        cancelKey: 'ACTIONS.CANCEL',
        labelKey: 'ACTIONS.NEWVALUE',
        titleKey: 'USER.PROFILE.CHANGEUSERNAME_TITLE',
        descriptionKey: 'USER.PROFILE.CHANGEUSERNAME_DESC',
        value: user.username,
      },
      width: '400px',
    });

    dialogRef
      .afterClosed()
      .pipe(
        map(({ value }: { value?: string }) => value),
        filter(Boolean),
        filter((value) => user.username != value),
        switchMap((username) => this.userService.updateUser({ userId: user.userId, username })),
      )
      .subscribe({
        next: () => {
          this.toast.showInfo('USER.TOAST.USERNAMECHANGED', true);
          this.refreshChanges$.emit();
        },
        error: (error) => {
          this.toast.showError(error);
        },
      });
  }

  public unlockUser(user: UserV2): void {
    this.userService
      .unlockUser(user.userId)
      .then(() => {
        this.toast.showInfo('USER.TOAST.UNLOCKED', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public generateMachineSecret(user: UserV2): void {
    this.newMgmtService
      .generateMachineSecret(user.userId)
      .then((resp) => {
        this.toast.showInfo('USER.TOAST.SECRETGENERATED', true);
        this.dialog.open(MachineSecretDialogComponent, {
          data: {
            clientId: resp.clientId,
            clientSecret: resp.clientSecret,
          },
          width: '400px',
        });
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public removeMachineSecret(user: UserV2): void {
    this.newMgmtService
      .removeMachineSecret(user.userId)
      .then(() => {
        this.toast.showInfo('USER.TOAST.SECRETREMOVED', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public changeState(user: UserV2, newState: UserState): void {
    if (newState === UserState.ACTIVE) {
      this.userService
        .reactivateUser(user.userId)
        .then(() => {
          this.toast.showInfo('USER.TOAST.REACTIVATED', true);
          this.refreshChanges$.emit();
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (newState === UserState.INACTIVE) {
      this.userService
        .deactivateUser(user.userId)
        .then(() => {
          this.toast.showInfo('USER.TOAST.DEACTIVATED', true);
          this.refreshChanges$.emit();
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public saveProfile(user: UserV2, profile: HumanProfile): void {
    this.userService
      .updateUser({
        userId: user.userId,
        profile: {
          givenName: profile.givenName,
          familyName: profile.familyName,
          nickName: profile.nickName,
          displayName: profile.displayName,
          preferredLanguage: profile.preferredLanguage,
          gender: profile.gender,
        },
      })
      .then(() => {
        this.toast.showInfo('USER.TOAST.SAVED', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public saveMachine(user: UserV2, form: ObservedValueOf<DetailFormMachineComponent['submitData']>): void {
    this.newMgmtService
      .updateMachine({
        userId: user.userId,
        name: form.name,
        description: form.description,
        accessTokenType: form.accessTokenType,
      })
      .then(() => {
        this.toast.showInfo('USER.TOAST.SAVED', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public resendEmailVerification(user: UserV2): void {
    this.newMgmtService
      .resendHumanEmailVerification(user.userId)
      .then(() => {
        this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public resendPhoneVerification(user: UserV2): void {
    this.newMgmtService
      .resendHumanPhoneVerification(user.userId)
      .then(() => {
        this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deletePhone(user: UserV2): void {
    this.userService
      .removePhone(user.userId)
      .then(() => {
        this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public navigateBack(): void {
    this._location.back();
  }

  public sendSetPasswordNotification(user: UserV2): void {
    this.newMgmtService
      .sendHumanResetPasswordNotification({
        userId: user.userId,
        type: SendHumanResetPasswordNotificationRequest_Type.EMAIL,
      })
      .then(() => {
        this.toast.showInfo('USER.TOAST.PASSWORDNOTIFICATIONSENT', true);
        this.refreshChanges$.emit();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deleteUser(user: UserV2): void {
    const data = {
      confirmKey: 'ACTIONS.DELETE',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'USER.DIALOG.DELETE_TITLE',
      descriptionKey: 'USER.DIALOG.DELETE_DESCRIPTION',
    };

    const dialogRef = this.dialog.open<WarnDialogComponent, typeof data, boolean>(WarnDialogComponent, {
      data,
      width: '400px',
    });

    dialogRef
      .afterClosed()
      .pipe(
        filter(Boolean),
        switchMap(() => this.userService.deleteUser(user.userId)),
      )
      .subscribe({
        next: () => {
          const params: Params = {
            deferredReload: true,
            type: user.type.case === 'human' ? 'humans' : 'machines',
          };
          this.router.navigate(['/users'], { queryParams: params }).then();
          this.toast.showInfo('USER.TOAST.DELETED', true);
        },
        error: (error) => this.toast.showError(error),
      });
  }

  public resendInitEmail(user: UserV2): void {
    const dialogRef = this.dialog.open<ResendEmailDialogComponent, ResendEmailDialogData, ResendEmailDialogResult>(
      ResendEmailDialogComponent,
      {
        width: '400px',
        data: {
          email: user.type.case === 'human' ? (user.type.value.email?.email ?? '') : '',
        },
      },
    );

    dialogRef
      .afterClosed()
      .pipe(
        filter((resp): resp is { send: true; email: string } => !!resp?.send && !!user.userId),
        switchMap(({ email }) => this.newMgmtService.resendHumanInitialization(user.userId, email)),
      )
      .subscribe({
        next: () => {
          this.toast.showInfo('USER.TOAST.INITEMAILSENT', true);
          this.refreshChanges$.emit();
        },
        error: (error) => this.toast.showError(error),
      });
  }

  public openEditDialog(user: UserWithHumanType, type: EditDialogType): void {
    switch (type) {
      case EditDialogType.PHONE:
        this.openEditPhoneDialog(user);
        return;
      case EditDialogType.EMAIL:
        this.openEditEmailDialog(user);
        return;
    }
  }

  private openEditEmailDialog(user: UserWithHumanType) {
    const data: EditDialogData = {
      confirmKey: 'ACTIONS.SAVE',
      cancelKey: 'ACTIONS.CANCEL',
      labelKey: 'ACTIONS.NEWVALUE',
      titleKey: 'USER.LOGINMETHODS.EMAIL.EDITTITLE',
      descriptionKey: 'USER.LOGINMETHODS.EMAIL.EDITDESC',
      isVerifiedTextKey: 'USER.LOGINMETHODS.EMAIL.ISVERIFIED',
      isVerifiedTextDescKey: 'USER.LOGINMETHODS.EMAIL.ISVERIFIEDDESC',
      value: user.type.value?.email?.email,
      type: EditDialogType.EMAIL,
    } as const;

    const dialogRefEmail = this.dialog.open<EditDialogComponent, EditDialogData, EditDialogResult>(EditDialogComponent, {
      data,
      width: '400px',
    });

    dialogRefEmail
      .afterClosed()
      .pipe(
        filter((resp): resp is Required<EditDialogResult> => !!resp?.value),
        switchMap(({ value, isVerified }) =>
          this.userService.setEmail({
            userId: user.userId,
            email: value,
            verification: isVerified ? { case: 'isVerified', value: isVerified } : { case: undefined },
          }),
        ),
        switchMap(() => {
          this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
          this.refreshChanges$.emit();
          if (user.state !== UserState.INITIAL) {
            return EMPTY;
          }
          return this.userService.resendInviteCode(user.userId);
        }),
      )
      .subscribe({
        next: () => this.toast.showInfo('USER.TOAST.INITEMAILSENT', true),
        error: (error) => this.toast.showError(error),
      });
  }

  private openEditPhoneDialog(user: UserWithHumanType) {
    const data = {
      confirmKey: 'ACTIONS.SAVE',
      cancelKey: 'ACTIONS.CANCEL',
      labelKey: 'ACTIONS.NEWVALUE',
      titleKey: 'USER.LOGINMETHODS.PHONE.EDITTITLE',
      descriptionKey: 'USER.LOGINMETHODS.PHONE.EDITDESC',
      value: user.type.value.phone?.phone,
      type: EditDialogType.PHONE,
      validator: Validators.compose([phoneValidator, requiredValidator]),
    };
    const dialogRefPhone = this.dialog.open<EditDialogComponent, typeof data, { value: string; isVerified: boolean }>(
      EditDialogComponent,
      { data, width: '400px' },
    );

    dialogRefPhone
      .afterClosed()
      .pipe(
        map((resp) => formatPhone(resp?.value)),
        filter(Boolean),
        switchMap(({ phone }) => this.userService.setPhone({ userId: user.userId, phone })),
      )
      .subscribe({
        next: () => {
          this.toast.showInfo('USER.TOAST.PHONESAVED', true);
          this.refreshChanges$.emit();
        },
        error: (error) => {
          this.toast.showError(error);
        },
      });
  }

  private getMetadataById(userId: string): Observable<MetadataQuery> {
    return defer(() => this.newMgmtService.listUserMetadata(userId)).pipe(
      map((metadata) => ({ state: 'success', value: metadata.result }) as const),
      startWith({ state: 'loading', value: [] as Metadata[] } as const),
      catchError((err) => of({ state: 'error', value: err.message ?? '' } as const)),
    );
  }

  public editMetadata(user: UserV2, metadata: Metadata[]): void {
    const setFcn = (key: string, value: string) =>
      this.newMgmtService.setUserMetadata({
        key,
        value: new TextEncoder().encode(value),
        id: user.userId,
      });
    const removeFcn = (key: string): Promise<any> => this.newMgmtService.removeUserMetadata({ key, id: user.userId });

    const dialogRef = this.dialog.open<MetadataDialogComponent, MetadataDialogData>(MetadataDialogComponent, {
      data: {
        metadata: [...metadata],
        setFcn: setFcn,
        removeFcn: removeFcn,
      },
    });

    dialogRef
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => {
        this.refreshMetadata$.next(true);
      });
  }

  public humanUser(user: UserV2): UserWithHumanType | undefined {
    if (user.type.case === 'human') {
      return { ...user, type: user.type };
    }
    return undefined;
  }
}

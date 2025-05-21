import { MediaMatcher } from '@angular/cdk/layout';
import { Component, DestroyRef, EventEmitter, OnInit, signal } from '@angular/core';
import { Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Buffer } from 'buffer';
import { defer, EMPTY, mergeWith, Observable, of, shareReplay, Subject, switchMap, take } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { phoneValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { InfoDialogComponent } from 'src/app/modules/info-dialog/info-dialog.component';
import {
  MetadataDialogComponent,
  MetadataDialogData,
} from 'src/app/modules/metadata/metadata-dialog/metadata-dialog.component';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { formatPhone } from 'src/app/utils/formatPhone';
import { EditDialogComponent, EditDialogData, EditDialogResult, EditDialogType } from './edit-dialog/edit-dialog.component';
import { LanguagesService } from 'src/app/services/languages.service';
import { Gender, HumanProfile, HumanUser, User, UserState } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { catchError, filter, map, startWith } from 'rxjs/operators';
import { pairwiseStartWith } from 'src/app/utils/pairwiseStartWith';
import { NewAuthService } from 'src/app/services/new-auth.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { NewMgmtService } from 'src/app/services/new-mgmt.service';
import { Metadata } from '@zitadel/proto/zitadel/metadata_pb';
import { UserService } from 'src/app/services/user.service';
import { LoginPolicy } from '@zitadel/proto/zitadel/policy_pb';
import { query } from '@angular/animations';

type UserQuery = { state: 'success'; value: User } | { state: 'error'; error: any } | { state: 'loading'; value?: User };

type MetadataQuery =
  | { state: 'success'; value: Metadata[] }
  | { state: 'loading'; value: Metadata[] }
  | { state: 'error'; error: any };

type UserWithHumanType = Omit<User, 'type'> & { type: { case: 'human'; value: HumanUser } };

@Component({
  selector: 'cnsl-auth-user-detail',
  templateUrl: './auth-user-detail.component.html',
  styleUrls: ['./auth-user-detail.component.scss'],
})
export class AuthUserDetailComponent implements OnInit {
  protected readonly genders: Gender[] = [Gender.MALE, Gender.FEMALE, Gender.DIVERSE];

  protected readonly ChangeType = ChangeType;
  public userLoginMustBeDomain: boolean = false;
  protected readonly UserState = UserState;

  protected USERGRANTCONTEXT: UserGrantContext = UserGrantContext.AUTHUSER;
  protected readonly refreshChanges$: EventEmitter<void> = new EventEmitter();
  protected readonly refreshMetadata$ = new Subject<true>();

  protected readonly settingsList: SidenavSetting[] = [
    { id: 'general', i18nKey: 'USER.SETTINGS.GENERAL' },
    { id: 'security', i18nKey: 'USER.SETTINGS.SECURITY' },
    { id: 'idp', i18nKey: 'USER.SETTINGS.IDP' },
    { id: 'grants', i18nKey: 'USER.SETTINGS.USERGRANTS' },
    { id: 'memberships', i18nKey: 'USER.SETTINGS.MEMBERSHIPS' },
    {
      id: 'metadata',
      i18nKey: 'USER.SETTINGS.METADATA',
      requiredRoles: { [PolicyComponentServiceType.MGMT]: ['user.read'] },
    },
  ];
  protected readonly user$: Observable<UserQuery>;
  protected readonly metadata$: Observable<MetadataQuery>;
  private readonly savedLanguage$: Observable<string>;
  protected readonly currentSetting$ = signal<SidenavSetting>(this.settingsList[0]);
  protected readonly loginPolicy$: Observable<LoginPolicy>;
  protected readonly userName$: Observable<string>;

  constructor(
    private translate: TranslateService,
    private toast: ToastService,
    protected grpcAuthService: GrpcAuthService,
    private dialog: MatDialog,
    private auth: AuthenticationService,
    private breadcrumbService: BreadcrumbService,
    public langSvc: LanguagesService,
    private readonly route: ActivatedRoute,
    private readonly newAuthService: NewAuthService,
    private readonly newMgmtService: NewMgmtService,
    private readonly userService: UserService,
    private readonly destroyRef: DestroyRef,
    private readonly router: Router,
  ) {
    this.user$ = this.getUser$().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.userName$ = this.getUserName(this.user$);
    this.savedLanguage$ = this.getSavedLanguage$(this.user$);
    this.metadata$ = this.getMetadata$().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.loginPolicy$ = defer(() => this.newMgmtService.getLoginPolicy()).pipe(
      catchError(() => EMPTY),
      map(({ policy }) => policy),
      filter(Boolean),
    );
  }

  getUserName(user$: Observable<UserQuery>) {
    return user$.pipe(
      map((query) => {
        const user = this.user(query);
        if (!user) {
          return '';
        }
        if (user.type.case === 'human') {
          return user.type.value.profile?.displayName ?? '';
        }
        if (user.type.case === 'machine') {
          return user.type.value.name;
        }
        return '';
      }),
    );
  }

  getSavedLanguage$(user$: Observable<UserQuery>) {
    return user$.pipe(
      switchMap((query) => {
        if (query.state !== 'success' || query.value.type.case !== 'human') {
          return EMPTY;
        }
        return query.value.type.value.profile?.preferredLanguage ?? EMPTY;
      }),
      startWith(this.translate.defaultLang),
    );
  }

  ngOnInit(): void {
    this.user$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((query) => {
      if ((query.state === 'loading' || query.state === 'success') && query.value?.type.case === 'human') {
        this.breadcrumbService.setBreadcrumb([
          new Breadcrumb({
            type: BreadcrumbType.AUTHUSER,
            name: query.value.type.value.profile?.displayName,
            routerLink: ['/users', 'me'],
          }),
        ]);
      }
    });

    this.user$.pipe(mergeWith(this.metadata$), takeUntilDestroyed(this.destroyRef)).subscribe((query) => {
      if (query.state == 'error') {
        this.toast.showError(query.error);
      }
    });

    this.savedLanguage$
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((savedLanguage) => this.translate.use(savedLanguage));

    const param = this.route.snapshot.queryParamMap.get('id');
    if (!param) {
      return;
    }
    const setting = this.settingsList.find(({ id }) => id === param);
    if (!setting) {
      return;
    }
    this.currentSetting$.set(setting);
  }

  private getUser$(): Observable<UserQuery> {
    return this.refreshChanges$.pipe(
      startWith(true),
      switchMap(() => this.getMyUser()),
      pairwiseStartWith(undefined),
      map(([prev, curr]) => {
        if (prev?.state === 'success' && curr.state === 'loading') {
          return { state: 'loading', value: prev.value } as const;
        }
        return curr;
      }),
    );
  }

  private getMyUser(): Observable<UserQuery> {
    return this.userService.user$.pipe(
      map((user) => ({ state: 'success' as const, value: user })),
      catchError((error) => of({ state: 'error', error } as const)),
      startWith({ state: 'loading' } as const),
    );
  }

  getMetadata$(): Observable<MetadataQuery> {
    return this.refreshMetadata$.pipe(
      startWith(true),
      switchMap(() => this.getMetadata()),
      pairwiseStartWith(undefined),
      map(([prev, curr]) => {
        if (prev?.state === 'success' && curr.state === 'loading') {
          return { state: 'loading', value: prev.value } as const;
        }
        return curr;
      }),
    );
  }

  private getMetadata(): Observable<MetadataQuery> {
    return defer(() => this.newAuthService.listMyMetadata()).pipe(
      map((metadata) => ({ state: 'success', value: metadata.result }) as const),
      startWith({ state: 'loading', value: [] as Metadata[] } as const),
      catchError((error) => of({ state: 'error', error } as const)),
    );
  }

  public changeUsername(user: User): void {
    const data = {
      confirmKey: 'ACTIONS.CHANGE' as const,
      cancelKey: 'ACTIONS.CANCEL' as const,
      labelKey: 'ACTIONS.NEWVALUE' as const,
      titleKey: 'USER.PROFILE.CHANGEUSERNAME_TITLE' as const,
      descriptionKey: 'USER.PROFILE.CHANGEUSERNAME_DESC' as const,
      value: user.username,
    };
    const dialogRef = this.dialog.open<EditDialogComponent, typeof data, EditDialogResult>(EditDialogComponent, {
      data,
      width: '400px',
    });

    dialogRef
      .afterClosed()
      .pipe(
        map((value) => value?.value),
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

  public saveProfile(user: User, profile: HumanProfile): void {
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

  public enteredPhoneCode(code: string): void {
    this.newAuthService
      .verifyMyPhone(code)
      .then(() => {
        this.toast.showInfo('USER.TOAST.PHONESAVED', true);
        this.refreshChanges$.emit();
        this.promptSetupforSMSOTP();
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public promptSetupforSMSOTP(): void {
    const dialogRef = this.dialog.open(InfoDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CONTINUE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.MFA.OTPSMS',
        descriptionKey: 'USER.MFA.SETUPOTPSMSDESCRIPTION',
      },
      width: '400px',
    });

    dialogRef
      .afterClosed()
      .pipe(
        filter(Boolean),
        switchMap(() => this.newAuthService.addMyAuthFactorOTPSMS()),
        switchMap(() => this.translate.get('USER.MFA.OTPSMSSUCCESS').pipe(take(1))),
      )
      .subscribe({
        next: (msg) => this.toast.showInfo(msg),
        error: (err) => this.toast.showError(err),
      });
  }

  public changedLanguage(language: string): void {
    this.translate.use(language);
  }

  public resendEmailVerification(user: User): void {
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

  public resendPhoneVerification(user: User): void {
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

  public deletePhone(user: User): void {
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
      )
      .subscribe({
        next: () => {
          this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
          this.refreshChanges$.emit();
        },
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

  public deleteUser(user: User): void {
    const data = {
      confirmKey: 'USER.DIALOG.DELETE_BTN',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'USER.DIALOG.DELETE_TITLE',
      descriptionKey: 'USER.DIALOG.DELETE_AUTH_DESCRIPTION',
    };

    const dialogRef = this.dialog.open<WarnDialogComponent, typeof data, boolean>(WarnDialogComponent, {
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
          this.toast.showInfo('USER.PAGES.DELETEACCOUNT_SUCCESS', true);
          this.auth.signout();
        },
        error: (error) => this.toast.showError(error),
      });
  }

  public editMetadata(user: User, metadata: Metadata[]): void {
    const setFcn = (key: string, value: string) =>
      this.newMgmtService.setUserMetadata({
        key,
        value: Buffer.from(value),
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

  protected readonly query = query;

  protected user(user: UserQuery): User | undefined {
    if (user.state === 'success' || user.state === 'loading') {
      return user.value;
    }
    return;
  }

  public async goToSetting(setting: string) {
    await this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { id: setting },
      queryParamsHandling: 'merge',
      skipLocationChange: true,
    });
  }

  public humanUser(userQuery: UserQuery): UserWithHumanType | undefined {
    const user = this.user(userQuery);
    if (user?.type.case === 'human') {
      return { ...user, type: user.type };
    }
    return;
  }
}

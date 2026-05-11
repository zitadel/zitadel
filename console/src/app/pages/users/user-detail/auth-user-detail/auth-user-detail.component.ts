import { Component, computed, DestroyRef, effect, OnInit, signal } from '@angular/core';
import { Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Buffer } from 'buffer';
import { defer, EMPTY, lastValueFrom, Observable, of, shareReplay, Subject, switchMap } from 'rxjs';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { phoneValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { InfoDialogComponent, InfoDialogData, InfoDialogResult } from 'src/app/modules/info-dialog/info-dialog.component';
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
import { QueryClient } from '@tanstack/angular-query-experimental';

type MetadataQuery =
  | { state: 'success'; value: Metadata[] }
  | { state: 'loading'; value: Metadata[] }
  | { state: 'error'; error: any };

type UserWithHumanType = Omit<User, 'type'> & { type: { case: 'human'; value: HumanUser } };

@Component({
  selector: 'cnsl-auth-user-detail',
  templateUrl: './auth-user-detail.component.html',
  styleUrls: ['./auth-user-detail.component.scss'],
  standalone: false,
})
export class AuthUserDetailComponent implements OnInit {
  protected readonly genders: Gender[] = [Gender.MALE, Gender.FEMALE, Gender.DIVERSE];

  protected readonly ChangeType = ChangeType;
  public userLoginMustBeDomain: boolean = false;
  protected readonly UserState = UserState;

  protected USERGRANTCONTEXT: UserGrantContext = UserGrantContext.AUTHUSER;
  protected readonly refreshMetadata$ = new Subject<true>();

  protected readonly settingsList: SidenavSetting[] = [
    { id: 'general', i18nKey: 'USER.SETTINGS.GENERAL' },
    { id: 'security', i18nKey: 'USER.SETTINGS.SECURITY' },
    { id: 'idp', i18nKey: 'USER.SETTINGS.IDP' },
    { id: 'grants', i18nKey: 'USER.SETTINGS.ROLEASSIGNMENTS' },
    { id: 'memberships', i18nKey: 'USER.SETTINGS.MEMBERSHIPS' },
    {
      id: 'metadata',
      i18nKey: 'USER.SETTINGS.METADATA',
      requiredRoles: { [PolicyComponentServiceType.MGMT]: ['user.read'] },
    },
  ];
  protected readonly metadata$: Observable<MetadataQuery>;
  protected readonly currentSetting$ = signal<SidenavSetting>(this.settingsList[0]);
  protected readonly loginPolicy$: Observable<LoginPolicy>;
  protected readonly user = this.userService.userQuery();
  protected readonly refreshChanges$ = new Subject<void>();

  protected readonly userName = computed(() => {
    const user = this.user.data();
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
  });

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
    private readonly queryClient: QueryClient,
  ) {
    this.metadata$ = this.getMetadata$().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.loginPolicy$ = defer(() => this.newMgmtService.getLoginPolicy()).pipe(
      catchError(() => EMPTY),
      map(({ policy }) => policy),
      filter(Boolean),
    );

    effect(() => {
      const user = this.user.data();
      if (!user || user.type.case !== 'human') {
        return;
      }

      this.breadcrumbService.setBreadcrumb([
        new Breadcrumb({
          type: BreadcrumbType.AUTHUSER,
          name: user.type.value.profile?.displayName,
          routerLink: ['/users', 'me'],
        }),
      ]);
    });

    effect(() => {
      const error = this.user.error();
      if (error) {
        this.toast.showError(error);
      }
    });
  }

  ngOnInit(): void {
    this.metadata$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((query) => {
      if (query.state == 'error') {
        this.toast.showError(query.error);
      }
    });

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

  protected invalidateUser() {
    this.refreshChanges$.next();
    return this.queryClient.invalidateQueries({
      queryKey: this.userService.userQueryOptions().queryKey,
    });
  }

  protected async changeUsername(user: User) {
    const data = {
      confirmKey: 'ACTIONS.CHANGE',
      cancelKey: 'ACTIONS.CANCEL',
      labelKey: 'ACTIONS.NEWVALUE',
      titleKey: 'USER.PROFILE.CHANGEUSERNAME_TITLE',
      descriptionKey: 'USER.PROFILE.CHANGEUSERNAME_DESC',
      value: user.username,
      type: EditDialogType.GENERIC,
    } as const satisfies EditDialogData;

    const dialogRef = this.dialog.open<EditDialogComponent, EditDialogData, EditDialogResult>(EditDialogComponent, {
      data,
      width: '400px',
    });

    const { value: username } = (await lastValueFrom(dialogRef.afterClosed())) ?? {};
    if (!username || user.username === username) {
      // no changes made
      return;
    }

    try {
      await this.userService.updateUser({ userId: user.userId, username });
      this.toast.showInfo('USER.TOAST.USERNAMECHANGED', true);
      await this.invalidateUser();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async saveProfile(user: User, profile: HumanProfile) {
    try {
      await this.userService.updateUser({
        userId: user.userId,
        profile: {
          givenName: profile.givenName,
          familyName: profile.familyName,
          nickName: profile.nickName,
          displayName: profile.displayName,
          preferredLanguage: profile.preferredLanguage,
          gender: profile.gender,
        },
      });
      this.toast.showInfo('USER.TOAST.SAVED', true);
      await this.invalidateUser();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async enteredPhoneCode(code: string) {
    try {
      await this.newAuthService.verifyMyPhone(code);
      this.toast.showInfo('USER.TOAST.PHONESAVED', true);
      await this.invalidateUser();
      await this.promptSetupforSMSOTP();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async promptSetupforSMSOTP() {
    const dialogRef = this.dialog.open<InfoDialogComponent, InfoDialogData, InfoDialogResult>(InfoDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.CONTINUE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'USER.MFA.OTPSMS',
        descriptionKey: 'USER.MFA.SETUPOTPSMSDESCRIPTION',
      },
      width: '400px',
    });

    const confirmed = await lastValueFrom(dialogRef.afterClosed());
    if (!confirmed) {
      return;
    }

    try {
      await this.newAuthService.addMyAuthFactorOTPSMS();
      const msg = await lastValueFrom(this.translate.get('USER.MFA.OTPSMSSUCCESS'));
      this.toast.showInfo(msg);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public changedLanguage(language: string): void {
    this.translate.use(language);
  }

  public async resendEmailVerification(user: User) {
    try {
      await this.newMgmtService.resendHumanEmailVerification(user.userId);
      this.toast.showInfo('USER.TOAST.EMAILVERIFICATIONSENT', true);
      await this.invalidateUser();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async resendPhoneVerification(user: User) {
    try {
      await this.newMgmtService.resendHumanPhoneVerification(user.userId);
      this.toast.showInfo('USER.TOAST.PHONEVERIFICATIONSENT', true);
      await this.invalidateUser();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async deletePhone(user: User) {
    try {
      await this.userService.removePhone(user.userId);
      this.toast.showInfo('USER.TOAST.PHONEREMOVED', true);
      await this.invalidateUser();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async openEditDialog(user: UserWithHumanType, type: EditDialogType) {
    switch (type) {
      case EditDialogType.PHONE:
        await this.openEditPhoneDialog(user);
        return;
      case EditDialogType.EMAIL:
        await this.openEditEmailDialog(user);
        return;
    }
  }

  private async openEditEmailDialog(user: UserWithHumanType) {
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

    const { value, isVerified } = (await lastValueFrom(dialogRefEmail.afterClosed())) ?? {};
    if (!value) {
      return;
    }

    try {
      await this.userService.setEmail({
        userId: user.userId,
        email: value,
        verification: isVerified ? { case: 'isVerified', value: isVerified } : { case: undefined },
      });
      this.toast.showInfo('USER.TOAST.EMAILSAVED', true);
      await this.invalidateUser();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  private async openEditPhoneDialog(user: UserWithHumanType) {
    const data = {
      confirmKey: 'ACTIONS.SAVE',
      cancelKey: 'ACTIONS.CANCEL',
      labelKey: 'ACTIONS.NEWVALUE',
      titleKey: 'USER.LOGINMETHODS.PHONE.EDITTITLE',
      descriptionKey: 'USER.LOGINMETHODS.PHONE.EDITDESC',
      value: user.type.value.phone?.phone,
      type: EditDialogType.PHONE,
      validator: Validators.compose([phoneValidator, requiredValidator]) ?? undefined,
    } as const satisfies EditDialogData;

    const dialogRefPhone = this.dialog.open<EditDialogComponent, EditDialogData, EditDialogResult>(EditDialogComponent, {
      data,
      width: '400px',
    });

    const { value } = (await lastValueFrom(dialogRefPhone.afterClosed())) ?? {};
    const formatted = formatPhone(value);
    if (!formatted) {
      return;
    }

    try {
      await this.userService.setPhone({ userId: user.userId, phone: formatted.phone });
      this.toast.showInfo('USER.TOAST.PHONESAVED', true);
      await this.invalidateUser();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async deleteUser(user: User) {
    const data = {
      confirmKey: 'USER.DIALOG.DELETE_BTN',
      cancelKey: 'ACTIONS.CANCEL',
      titleKey: 'USER.DIALOG.DELETE_TITLE',
      descriptionKey: 'USER.DIALOG.DELETE_AUTH_DESCRIPTION',
    };

    const dialogRef = this.dialog.open<WarnDialogComponent, typeof data, boolean>(WarnDialogComponent, {
      data,
      width: '400px',
    });

    const confirmed = await lastValueFrom(dialogRef.afterClosed());
    if (!confirmed) {
      return;
    }

    try {
      await this.userService.deleteUser(user.userId);
      this.toast.showInfo('USER.PAGES.DELETEACCOUNT_SUCCESS', true);
      this.auth.signout();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  public async editMetadata(user: User, metadata: Metadata[]) {
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

    await lastValueFrom(dialogRef.afterClosed());
    this.refreshMetadata$.next(true);
  }

  protected readonly query = query;

  public humanUser(user: User | undefined): UserWithHumanType | undefined {
    if (user?.type.case === 'human') {
      return { ...user, type: user.type };
    }
    return;
  }
}

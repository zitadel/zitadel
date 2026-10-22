import { Component, OnInit } from '@angular/core';
import { Validators } from '@angular/forms';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { firstValueFrom, from, timeout } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { User, Gender } from 'src/app/proto/generated/zitadel/user_pb';
import { create } from '@bufbuild/protobuf';
import { HumanProfile, HumanProfileSchema, Gender as ProfileGender } from '@zitadel/proto/zitadel/user/v2/user_pb';
import {
  EditDialogComponent,
  EditDialogData,
  EditDialogResult,
  EditDialogType,
} from 'src/app/pages/users/user-detail/auth-user-detail/edit-dialog/edit-dialog.component';
import { phoneValidator, emailValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { AccountAuthService } from './auth.service';
import { AccountBranding } from './branding';
import { AccountIdentity } from './identity';
import { AccountPreferences } from './preferences';

@Component({ selector: 'cnsl-account', templateUrl: './shell.html', standalone: false })
export class AccountComponent implements OnInit {
  helpUrl = '';
  privacyUrl = '';
  readonly nav = [
    { path: '/', key: 'HOME' },
    { path: '/perfil', key: 'PROFILE' },
    { path: '/seguranca', key: 'SECURITY' },
    { path: '/aplicativos', key: 'APPS' },
  ];
  constructor(
    public auth: AccountAuthService,
    public branding: AccountBranding,
    public identity: AccountIdentity,
    public preferences: AccountPreferences,
    public native: GrpcAuthService,
    public translate: TranslateService,
    public router: Router,
  ) {
    this.preferences.initialize();
  }
  get breadcrumb(): string {
    const path = this.router.url.split('?')[0];
    return (
      (
        {
          '/': 'HOME',
          '/perfil': 'PROFILE',
          '/aplicativos': 'APPS',
          '/atividade': 'ACTIVITY',
          '/seguranca': 'SECURITY',
          '/seguranca/senha': 'PASSWORD',
          '/seguranca/passkeys': 'PASSKEYS',
          '/seguranca/mfa': 'MFA',
          '/signedout': 'SIGNED_OUT',
          '/auth/error': 'AUTH_ERROR',
        } as Record<string, string>
      )[path] ?? 'NOT_FOUND'
    );
  }
  async ngOnInit(): Promise<void> {
    if (!this.auth.authenticated) return;
    try {
      const policy = await this.native.getMyPrivacyPolicy(this.auth.environment.loginOrgId);
      this.helpUrl = this.externalLink(policy.helpLink);
      this.privacyUrl = this.externalLink(policy.privacyLink);
    } catch {
      /* Links without an effective native policy are omitted. */
    }
  }
  private externalLink(value: string): string {
    try {
      const url = new URL(value);
      return url.protocol === 'https:' && !url.username && !url.password ? url.href : '';
    } catch {
      return '';
    }
  }
  focusContent(event: Event): void {
    event.preventDefault();
    document.getElementById('content')?.focus();
  }
  language(value: string): void {
    this.preferences.language(value);
  }
}
@Component({ selector: 'cnsl-account-page', templateUrl: './page.html', standalone: false })
export class AccountPageComponent implements OnInit {
  user?: User.AsObject;
  loading = false;
  failed = false;
  saving = false;
  readonly genders = [ProfileGender.UNSPECIFIED, ProfileGender.FEMALE, ProfileGender.MALE, ProfileGender.DIVERSE];
  constructor(
    public router: Router,
    public auth: AccountAuthService,
    public native: GrpcAuthService,
    private dialog: MatDialog,
    private toast: ToastService,
    private identity: AccountIdentity,
    private translate: TranslateService,
    private preferences: AccountPreferences,
  ) {}
  get path(): string {
    return this.router.url.split('?')[0];
  }
  get title(): string {
    return (
      (
        {
          '/': 'HOME',
          '/perfil': 'PROFILE',
          '/aplicativos': 'APPS',
          '/atividade': 'ACTIVITY',
          '/seguranca': 'SECURITY',
          '/seguranca/mfa': 'MFA',
          '/seguranca/passkeys': 'PASSKEYS',
          '/auth/error': 'AUTH_ERROR',
          '/signedout': 'SIGNED_OUT',
        } as Record<string, string>
      )[this.path] ?? 'NOT_FOUND'
    );
  }
  get crossOrigin(): boolean {
    return location.origin !== this.auth.environment.issuer;
  }
  get human() {
    return this.user?.human;
  }
  get profileLanguages(): string[] {
    const current = this.profile?.preferredLanguage ?? '';
    return [...new Set([current ? 'und' : '', 'pt', 'en', 'es', current])];
  }
  profile?: HumanProfile;
  ngOnInit(): void {
    if (
      this.auth.authenticated &&
      ['/', '/perfil', '/seguranca', '/seguranca/mfa', '/seguranca/passkeys'].includes(this.path)
    )
      void this.refresh();
  }
  async refresh(): Promise<void> {
    this.loading = true;
    this.failed = false;
    try {
      this.user = (await firstValueFrom(from(this.native.getMyUser()).pipe(timeout(15000)))).user;
      const p = this.human?.profile;
      this.profile = p
        ? create(HumanProfileSchema, {
            givenName: p.firstName,
            familyName: p.lastName,
            nickName: p.nickName,
            displayName: p.displayName,
            preferredLanguage: p.preferredLanguage,
            avatarUrl: p.avatarUrl,
            gender: Number(p.gender),
          })
        : undefined;
      if (!this.human || !this.user) throw new Error('Human account required');
      this.identity.update(this.user);
    } catch {
      this.failed = true;
    } finally {
      this.loading = false;
    }
  }
  async perform(operation: () => Promise<unknown>): Promise<boolean> {
    if (this.saving || !this.auth.authenticated) return false;
    this.saving = true;
    try {
      await operation();
      await this.refresh();
      this.toast.showInfo('USER.TOAST.SAVED', true);
      return true;
    } catch (error) {
      this.toast.showError(error);
      return false;
    } finally {
      this.saving = false;
    }
  }
  async save(profile: HumanProfile): Promise<void> {
    const genders = [Gender.GENDER_UNSPECIFIED, Gender.GENDER_FEMALE, Gender.GENDER_MALE, Gender.GENDER_DIVERSE];
    const saved = await this.perform(() =>
      this.native.updateMyProfile(
        profile.givenName,
        profile.familyName,
        profile.nickName,
        profile.displayName,
        profile.preferredLanguage,
        genders[profile.gender ?? 0],
      ),
    );
    const language = profile.preferredLanguage;
    if (saved && language && language !== 'und') this.preferences.language(language);
  }
  async edit(type: EditDialogType): Promise<void> {
    const email = type === EditDialogType.EMAIL;
    const data: EditDialogData = {
      confirmKey: 'ACTIONS.SAVE',
      cancelKey: 'ACTIONS.CANCEL',
      labelKey: 'ACTIONS.NEWVALUE',
      type,
      titleKey: email ? 'USER.LOGINMETHODS.EMAIL.EDITTITLE' : 'USER.LOGINMETHODS.PHONE.EDITTITLE',
      descriptionKey: email ? 'USER.LOGINMETHODS.EMAIL.EDITDESC' : 'USER.LOGINMETHODS.PHONE.EDITDESC',
      value: email ? this.human?.email?.email : this.human?.phone?.phone,
      validator: email
        ? (Validators.compose([requiredValidator, emailValidator]) ?? undefined)
        : (Validators.compose([requiredValidator, phoneValidator]) ?? undefined),
    };
    const result = await firstValueFrom(
      this.dialog
        .open<EditDialogComponent, EditDialogData, EditDialogResult>(EditDialogComponent, { data, width: '400px' })
        .afterClosed(),
    );
    if (result?.value) {
      const value = result.value;
      await this.perform(() => (email ? this.native.setMyEmail(value) : this.native.setMyPhone(value)));
    }
  }
  async editUsername(): Promise<void> {
    const data: EditDialogData = {
      confirmKey: 'ACTIONS.CHANGE',
      cancelKey: 'ACTIONS.CANCEL',
      labelKey: 'ACTIONS.NEWVALUE',
      titleKey: 'USER.PROFILE.CHANGEUSERNAME_TITLE',
      descriptionKey: 'USER.PROFILE.CHANGEUSERNAME_DESC',
      value: this.user?.userName,
      type: EditDialogType.GENERIC,
    };
    const result = await firstValueFrom(
      this.dialog
        .open<EditDialogComponent, EditDialogData, EditDialogResult>(EditDialogComponent, { data, width: '400px' })
        .afterClosed(),
    );
    if (result?.value) {
      const value = result.value;
      await this.perform(() => this.native.updateMyUserName(value));
    }
  }
  emailVerification(): void {
    void this.perform(() => this.native.resendMyEmailVerification());
  }
  phoneVerification(): void {
    void this.perform(() => this.native.resendMyPhoneVerification());
  }
  verifyPhone(code: string): void {
    void this.perform(() => this.native.verifyMyPhone(code));
  }
  deletePhone(): void {
    void this.perform(() => this.native.removeMyPhone());
  }
}

@Component({
  selector: 'cnsl-account-password',
  standalone: false,
  template: `<a class="password-back" routerLink="/seguranca">‹ {{ 'ACCOUNT.SECURITY' | translate }}</a>
    <cnsl-password class="account-password-form"></cnsl-password>`,
})
export class AccountPasswordPageComponent {}

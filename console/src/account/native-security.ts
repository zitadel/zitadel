import { Component } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { OAuthService } from 'angular-oauth2-oidc';
import { of } from 'rxjs';
import { SecondFactorType } from '@zitadel/proto/zitadel/policy_pb';
import { AuthPasswordlessComponent } from 'src/app/pages/users/user-detail/auth-user-detail/auth-passwordless/auth-passwordless.component';
import { AuthUserMfaComponent } from 'src/app/pages/users/user-detail/auth-user-detail/auth-user-mfa/auth-user-mfa.component';
import {
  AuthFactorDialogComponent,
  AddAuthFactorDialogData,
} from 'src/app/pages/users/user-detail/auth-user-detail/auth-factor-dialog/auth-factor-dialog.component';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { NewAuthService } from 'src/app/services/new-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { AccountAuthService } from './auth.service';

/** Self-service only: use Auth API My operations, never the administrative V2 registration-link RPC. */
@Component({
  selector: 'cnsl-account-passkey-dialog',
  standalone: false,
  template: ` <h2 mat-dialog-title>{{ 'USER.PASSWORDLESS.DIALOG.ADD_TITLE' | translate }}</h2>
    <div mat-dialog-content>
      <p>{{ 'ACCOUNT.PASSKEY_ORIGIN' | translate: { issuer: auth.environment.issuer } }}</p>
      @if (busy) {
        <p role="status">{{ 'ACCOUNT.LOADING' | translate }}</p>
      }
      @if (link) {
        <p>{{ 'ACCOUNT.PASSKEY_RETURN' | translate }}</p>
        <a class="passkey-handoff" mat-raised-button color="primary" [href]="link" target="_blank" rel="noopener noreferrer">
          {{ 'ACCOUNT.PASSKEY_CONTINUE' | translate }} ↗
        </a>
      }
    </div>
    <div mat-dialog-actions class="passkey-actions">
      <button mat-stroked-button [disabled]="busy" (click)="dialog.close()">{{ 'ACTIONS.CLOSE' | translate }}</button>
      @if (!link) {
        <button mat-raised-button color="primary" [disabled]="busy" (click)="enroll('device')">
          {{ 'ACCOUNT.PASSKEY_DEVICE' | translate }}
        </button>
      }
    </div>
    @if (!link) {
      <div class="passkey-alternative">
        <button class="text-action" [disabled]="busy" (click)="enroll('email')">
          {{ 'ACCOUNT.PASSKEY_EMAIL_ALTERNATIVE' | translate }}
        </button>
      </div>
    }`,
})
export class AccountPasskeyDialog {
  busy = false;
  link = '';
  constructor(
    public dialog: MatDialogRef<AccountPasskeyDialog>,
    private native: GrpcAuthService,
    public auth: AccountAuthService,
    private oauth: OAuthService,
    private toast: ToastService,
  ) {}
  async enroll(mode: 'email' | 'device'): Promise<void> {
    if (this.busy || !this.auth.authenticated) return;
    this.busy = true;
    try {
      if (mode === 'email') {
        await this.native.sendMyPasswordlessLink();
        this.toast.showInfo('USER.TOAST.PASSWORDLESSREGISTRATIONSENT', true);
        this.dialog.close();
      } else {
        const response = await this.native.addMyPasswordlessLink();
        const source = new URL(response.link);
        // Validate the native self-registration link before adapting its documented query
        // names to the official Login V2 page. Never log/persist the one-time enrollment code.
        const issuer = this.auth.environment.issuer;
        const user = source.searchParams.get('userID');
        const org = source.searchParams.get('orgID');
        const codeId = source.searchParams.get('codeID');
        const code = source.searchParams.get('code');
        if (
          source.origin !== issuer ||
          source.pathname !== '/ui/login/login/passwordless/init' ||
          !user ||
          user !== this.oauth.getIdentityClaims()?.['sub'] ||
          org !== this.auth.environment.loginOrgId ||
          !codeId ||
          !code
        ) {
          throw new Error('Invalid native self-registration link');
        }
        const target = new URL('/ui/v2/login/passkey/set', issuer);
        target.search = new URLSearchParams({ userId: user, organization: org, codeId, code }).toString();
        this.link = target.href;
      }
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.busy = false;
    }
  }
}

/** Reuses the native table/removal UI; cross-origin enrollment is handed to Login V2. */
@Component({
  selector: 'cnsl-account-passkeys',
  standalone: false,
  templateUrl: '../app/pages/users/user-detail/auth-user-detail/auth-passwordless/auth-passwordless.component.html',
  styleUrls: ['../app/pages/users/user-detail/auth-user-detail/auth-passwordless/auth-passwordless.component.scss'],
})
export class AccountPasskeysComponent extends AuthPasswordlessComponent {
  private enrollmentPending = false;
  get loadError(): boolean {
    return !!this.error;
  }
  override getPasswordless(): void {
    this.error = '';
    super.getPasswordless();
  }
  constructor(
    native: GrpcAuthService,
    toast: ToastService,
    private accountDialog: MatDialog,
    private accountAuth: AccountAuthService,
  ) {
    super(native, toast, accountDialog);
  }
  override addPasswordless(): void {
    if (location.origin === this.accountAuth.environment.issuer) {
      super.addPasswordless();
      return;
    }
    if (this.enrollmentPending) return;
    this.enrollmentPending = true;
    this.accountDialog
      .open(AccountPasskeyDialog, { width: '480px' })
      .afterClosed()
      .subscribe(() => {
        this.enrollmentPending = false;
        this.getPasswordless();
      });
  }
}

/** Native MFA operations. Only offer WebAuthn enrollment on the origin accepted by the native API. */
@Component({
  selector: 'cnsl-account-mfa',
  standalone: false,
  templateUrl: '../app/pages/users/user-detail/auth-user-detail/auth-user-mfa/auth-user-mfa.component.html',
  styleUrls: ['../app/pages/users/user-detail/auth-user-detail/auth-user-mfa/auth-user-mfa.component.scss'],
})
export class AccountMfaComponent extends AuthUserMfaComponent {
  private policyFailed = false;
  get loadError(): boolean {
    return !!this.error || this.policyFailed;
  }
  override getMFAs(): void {
    this.error = '';
    super.getMFAs();
  }
  constructor(
    private accountPolicy: NewAuthService,
    private accountToast: ToastService,
    private accountDialog: MatDialog,
    private accountAuth: AccountAuthService,
  ) {
    super(accountPolicy, accountToast, accountDialog);
  }
  override applyOrgPolicy(): void {
    this.policyFailed = false;
    this.accountPolicy
      .getMyLoginPolicy()
      .then(({ policy }) => {
        if (!policy) throw new Error('Missing policy');
        const factors = policy.secondFactors;
        this.otpAvailable$.next(factors.includes(SecondFactorType.OTP));
        this.u2fAvailable$.next(
          factors.includes(SecondFactorType.U2F) && location.origin === this.accountAuth.environment.issuer,
        );
        this.otpSmsAvailable$.next(factors.includes(SecondFactorType.OTP_SMS));
        this.otpEmailAvailable$.next(factors.includes(SecondFactorType.OTP_EMAIL));
      })
      .catch(() => {
        this.policyFailed = true;
      });
  }
  override addAuthFactor(): void {
    const data: AddAuthFactorDialogData = {
      otp$: this.otpAvailable$,
      u2f$: location.origin === this.accountAuth.environment.issuer ? this.u2fAvailable$ : of(false),
      otpSms$: this.otpSmsAvailable$,
      otpEmail$: this.otpEmailAvailable$,
      otpDisabled$: this.otpDisabled$,
      otpSmsDisabled$: this.otpSmsDisabled$,
      otpEmailDisabled$: this.otpEmailDisabled$,
      phoneVerified: this.phoneVerified,
    };
    this.accountDialog
      .open(AuthFactorDialogComponent, { data })
      .afterClosed()
      .subscribe(() => this.getMFAs());
  }
}

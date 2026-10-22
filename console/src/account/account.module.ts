import { CommonModule } from '@angular/common';
import { HttpClient, provideHttpClient } from '@angular/common/http';
import { inject, NgModule, provideAppInitializer } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { Router, RouterModule, Routes, CanActivateFn } from '@angular/router';
import { MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule } from '@angular/material/sort';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatSelectModule } from '@angular/material/select';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { TranslateModule, TranslateLoader, TranslationObject } from '@ngx-translate/core';
import { OAuthService, OAuthModule, OAuthStorage, ValidationHandler } from 'angular-oauth2-oidc';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { QRCodeComponent } from 'angularx-qrcode';
import { forkJoin, map } from 'rxjs';
import { CountryCallingCodesService } from 'src/app/services/country-calling-codes.service';
import * as countries from 'i18n-iso-countries';
import ptCountries from 'i18n-iso-countries/langs/pt.json';
import enCountries from 'i18n-iso-countries/langs/en.json';
import { ExternalIdpsComponent } from 'src/app/pages/users/user-detail/external-idps/external-idps.component';
import { GrpcService } from 'src/app/services/grpc.service';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { StorageService } from 'src/app/services/storage.service';
import { ThemeService } from 'src/app/services/theme.service';
import { AccountAuthService } from './auth.service';
import { AccountStorage } from './storage';
import { AccountSignatureHandler } from './signature';
import { AccountComponent, AccountPageComponent, AccountPasswordPageComponent } from './shell';
import { AccountBranding } from './branding';
import { PasswordComponent } from 'src/app/pages/users/user-detail/password/password.component';
import { AuthUserMfaComponent } from 'src/app/pages/users/user-detail/auth-user-detail/auth-user-mfa/auth-user-mfa.component';
import { AuthPasswordlessComponent } from 'src/app/pages/users/user-detail/auth-user-detail/auth-passwordless/auth-passwordless.component';
import { DialogPasswordlessComponent } from 'src/app/pages/users/user-detail/auth-user-detail/auth-passwordless/dialog-passwordless/dialog-passwordless.component';
import { AuthFactorDialogComponent } from 'src/app/pages/users/user-detail/auth-user-detail/auth-factor-dialog/auth-factor-dialog.component';
import { DialogU2FComponent } from 'src/app/pages/users/user-detail/auth-user-detail/dialog-u2f/dialog-u2f.component';
import { ContactComponent } from 'src/app/pages/users/user-detail/contact/contact.component';
import { PhoneDetailComponent } from 'src/app/pages/users/user-detail/phone-detail/phone-detail.component';
import { CodeDialogComponent } from 'src/app/pages/users/user-detail/auth-user-detail/code-dialog/code-dialog.component';
import { EditDialogComponent } from 'src/app/pages/users/user-detail/auth-user-detail/edit-dialog/edit-dialog.component';
import { CardModule } from 'src/app/modules/card/card.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { PasswordComplexityViewModule } from 'src/app/modules/password-complexity-view/password-complexity-view.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { DetailFormModule } from 'src/app/pages/users/user-detail/detail-form/detail-form.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';

countries.registerLocale(ptCountries);
countries.registerLocale(enCountries);

export class AccountTranslations implements TranslateLoader {
  constructor(private http: HttpClient) {}
  getTranslation(lang: string) {
    const locale = lang === 'pt' ? 'pt' : 'en';
    return forkJoin([
      this.http.get<TranslationObject>(`/assets/native-i18n/${locale}.json`),
      this.http.get<TranslationObject>(`/assets/i18n/${locale}.json`),
    ]).pipe(map(([native, account]) => ({ ...native, ...account })));
  }
}
const authenticated: CanActivateFn = async (_route, state) => {
  const auth = inject(AccountAuthService);
  const router = inject(Router);
  if (auth.error) return router.parseUrl('/auth/error');
  if (auth.authenticated) return true;
  if (!auth.error) auth.login(state.url);
  return auth.error ? router.parseUrl('/auth/error') : false;
};
export const accountRoutes: Routes = [
  ...['', 'perfil', 'seguranca', 'seguranca/passkeys', 'seguranca/mfa'].map((path) => ({
    path,
    pathMatch: 'full' as const,
    component: AccountPageComponent,
    canActivate: [authenticated],
  })),
  { path: 'seguranca/senha', component: AccountPasswordPageComponent, canActivate: [authenticated] },
  { path: 'auth/callback', component: AccountPageComponent },
  { path: 'auth/error', component: AccountPageComponent },
  { path: 'signedout', component: AccountPageComponent },
  { path: '**', component: AccountPageComponent },
];
@NgModule({
  declarations: [
    AccountComponent,
    AccountPageComponent,
    AccountPasswordPageComponent,
    ExternalIdpsComponent,
    PasswordComponent,
    AuthUserMfaComponent,
    AuthPasswordlessComponent,
    DialogPasswordlessComponent,
    AuthFactorDialogComponent,
    DialogU2FComponent,
    ContactComponent,
    PhoneDetailComponent,
    CodeDialogComponent,
    EditDialogComponent,
  ],
  imports: [
    BrowserModule,
    CommonModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forRoot(accountRoutes, { scrollPositionRestoration: 'enabled' }),
    OAuthModule.forRoot(),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useFactory: (http: HttpClient) => new AccountTranslations(http),
        deps: [HttpClient],
      },
    }),
    MatDialogModule,
    MatButtonModule,
    MatIconModule,
    MatTableModule,
    MatSortModule,
    MatTooltipModule,
    MatProgressSpinnerModule,
    MatCheckboxModule,
    MatSelectModule,
    MatSnackBarModule,
    QRCodeComponent,
    CardModule,
    RefreshTableModule,
    TableActionsModule,
    InfoSectionModule,
    DetailLayoutModule,
    PasswordComplexityViewModule,
    InputModule,
    WarnDialogModule,
    AvatarModule,
    PaginatorModule,
    DetailFormModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    CopyToClipboardModule,
  ],
  providers: [
    CountryCallingCodesService,
    provideHttpClient(),
    provideTanStackQuery(new QueryClient()),
    ThemeService,
    { provide: OAuthStorage, useExisting: AccountStorage },
    { provide: StorageService, useExisting: AccountStorage },
    { provide: AuthenticationService, useExisting: AccountAuthService },
    { provide: ValidationHandler, useClass: AccountSignatureHandler },
    { provide: 'windowObject', useValue: window },
    provideAppInitializer(() => {
      const grpc = inject(GrpcService),
        auth = inject(AccountAuthService),
        branding = inject(AccountBranding),
        oauth = inject(OAuthService);
      return grpc
        .loadAppEnvironment()
        .then(() => auth.initialize())
        .then(async () => {
          if (!auth.authenticated) return;
          const { user } = await grpc.authNew.getMyUser({});
          if (
            !user ||
            user.id !== oauth.getIdentityClaims()?.['sub'] ||
            user.details?.resourceOwner !== auth.environment.loginOrgId
          )
            throw new Error('Identity context mismatch');
          await branding.load();
        })
        .catch(() => {
          auth.error = true;
          oauth.logOut(true);
          auth.authenticationChanged.next(false);
          history.replaceState(null, '', '/auth/error');
        });
    }),
  ],
  bootstrap: [AccountComponent],
})
export class AccountModule {}

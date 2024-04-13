import { PlatformLocation } from '@angular/common';
import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { AuthConfig } from 'angular-oauth2-oidc';
import { catchError, switchMap, tap, throwError } from 'rxjs';

import { AdminServiceClient } from '../proto/generated/zitadel/AdminServiceClientPb';
import { AuthServiceClient } from '../proto/generated/zitadel/AuthServiceClientPb';
import { ManagementServiceClient } from '../proto/generated/zitadel/ManagementServiceClientPb';
import { fallbackLanguage, supportedLanguagesRegexp } from '../utils/language';
import { AuthenticationService } from './authentication.service';
import { EnvironmentService } from './environment.service';
import { ExhaustedService } from './exhausted.service';
import { AuthInterceptor } from './interceptors/auth.interceptor';
import { ExhaustedGrpcInterceptor } from './interceptors/exhausted.grpc.interceptor';
import { I18nInterceptor } from './interceptors/i18n.interceptor';
import { OrgInterceptor } from './interceptors/org.interceptor';
import { StorageService } from './storage.service';
import { FeatureServiceClient } from '../proto/generated/zitadel/feature/v2beta/Feature_serviceServiceClientPb';

@Injectable({
  providedIn: 'root',
})
export class GrpcService {
  public auth!: AuthServiceClient;
  public mgmt!: ManagementServiceClient;
  public admin!: AdminServiceClient;
  public feature!: FeatureServiceClient;

  constructor(
    private envService: EnvironmentService,
    private platformLocation: PlatformLocation,
    private authenticationService: AuthenticationService,
    private storageService: StorageService,
    private dialog: MatDialog,
    private translate: TranslateService,
    private exhaustedService: ExhaustedService,
  ) {}

  public loadAppEnvironment(): Promise<any> {
    // We use the browser language until we can make API requests to get the users configured language.

    const browserLanguage = this.translate.getBrowserLang();
    const language = browserLanguage?.match(supportedLanguagesRegexp) ? browserLanguage : fallbackLanguage;
    return this.translate
      .use(language || this.translate.defaultLang)
      .pipe(
        switchMap(() => this.envService.env),
        tap((env) => {
          if (!env?.api || !env?.issuer) {
            return;
          }
          const interceptors = {
            unaryInterceptors: [
              new ExhaustedGrpcInterceptor(this.exhaustedService, this.envService),
              new OrgInterceptor(this.storageService),
              new AuthInterceptor(this.authenticationService, this.storageService, this.dialog),
              new I18nInterceptor(this.translate),
            ],
          };

          this.auth = new AuthServiceClient(
            env.api,
            null,
            // @ts-ignore
            interceptors,
          );
          this.mgmt = new ManagementServiceClient(
            env.api,
            null,
            // @ts-ignore
            interceptors,
          );
          this.admin = new AdminServiceClient(
            env.api,
            null,
            // @ts-ignore
            interceptors,
          );
          this.feature = new FeatureServiceClient(
            env.api,
            null,
            // @ts-ignore
            interceptors,
          );

          const authConfig: AuthConfig = {
            scope: 'openid profile email',
            responseType: 'code',
            oidc: true,
            clientId: env.clientid,
            issuer: env.issuer,
            redirectUri: window.location.origin + this.platformLocation.getBaseHrefFromDOM() + 'auth/callback',
            postLogoutRedirectUri: window.location.origin + this.platformLocation.getBaseHrefFromDOM() + 'signedout',
            requireHttps: false,
          };

          this.authenticationService.initConfig(authConfig);
        }),
        catchError((err) => {
          console.error('Failed to load environment from assets', err);
          return throwError(() => err);
        }),
      )
      .toPromise();
  }
}

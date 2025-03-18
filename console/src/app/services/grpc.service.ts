import { PlatformLocation } from '@angular/common';
import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { AuthConfig } from 'angular-oauth2-oidc';
import { catchError, firstValueFrom, switchMap, tap } from 'rxjs';

import { AdminServiceClient } from '../proto/generated/zitadel/AdminServiceClientPb';
import { AuthServiceClient } from '../proto/generated/zitadel/AuthServiceClientPb';
import { ManagementServiceClient } from '../proto/generated/zitadel/ManagementServiceClientPb';
import { fallbackLanguage, supportedLanguagesRegexp } from '../utils/language';
import { AuthenticationService } from './authentication.service';
import { EnvironmentService } from './environment.service';
import { ExhaustedService } from './exhausted.service';
import { AuthInterceptor, AuthInterceptorProvider, NewConnectWebAuthInterceptor } from './interceptors/auth.interceptor';
import { ExhaustedGrpcInterceptor } from './interceptors/exhausted.grpc.interceptor';
import { I18nInterceptor } from './interceptors/i18n.interceptor';
import { NewConnectWebOrgInterceptor, OrgInterceptor, OrgInterceptorProvider } from './interceptors/org.interceptor';
import { StorageService } from './storage.service';
import { UserServiceClient } from '../proto/generated/zitadel/user/v2/User_serviceServiceClientPb';
//@ts-ignore
import { createUserServiceClient } from '@zitadel/client/v2';
//@ts-ignore
import { createAuthServiceClient, createManagementServiceClient } from '@zitadel/client/v1';
import { createGrpcWebTransport } from '@connectrpc/connect-web';
import { FeatureServiceClient } from '../proto/generated/zitadel/feature/v2/Feature_serviceServiceClientPb';

@Injectable({
  providedIn: 'root',
})
export class GrpcService {
  public auth!: AuthServiceClient;
  public mgmt!: ManagementServiceClient;
  public admin!: AdminServiceClient;
  public feature!: FeatureServiceClient;
  public user!: UserServiceClient;
  public userNew!: ReturnType<typeof createUserServiceClient>;
  public mgmtNew!: ReturnType<typeof createManagementServiceClient>;
  public authNew!: ReturnType<typeof createAuthServiceClient>;

  constructor(
    private readonly envService: EnvironmentService,
    private readonly platformLocation: PlatformLocation,
    private readonly authenticationService: AuthenticationService,
    private readonly storageService: StorageService,
    private readonly translate: TranslateService,
    private readonly exhaustedService: ExhaustedService,
    private readonly authInterceptor: AuthInterceptor,
    private readonly authInterceptorProvider: AuthInterceptorProvider,
    private readonly orgInterceptorProvider: OrgInterceptorProvider,
  ) {}

  public loadAppEnvironment(): Promise<any> {
    // We use the browser language until we can make API requests to get the users configured language.

    const browserLanguage = this.translate.getBrowserLang();
    const language = browserLanguage?.match(supportedLanguagesRegexp) ? browserLanguage : fallbackLanguage;
    const init = this.translate.use(language || this.translate.defaultLang).pipe(
      switchMap(() => this.envService.env),
      tap((env) => {
        if (!env?.api || !env?.issuer) {
          return;
        }
        const interceptors = {
          unaryInterceptors: [
            new ExhaustedGrpcInterceptor(this.exhaustedService, this.envService),
            new OrgInterceptor(this.orgInterceptorProvider),
            this.authInterceptor,
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
        this.user = new UserServiceClient(
          env.api,
          null,
          // @ts-ignore
          interceptors,
        );

        const transport = createGrpcWebTransport({
          baseUrl: env.api,
          interceptors: [NewConnectWebAuthInterceptor(this.authInterceptorProvider)],
        });
        const transportOldAPIs = createGrpcWebTransport({
          baseUrl: env.api,
          interceptors: [
            NewConnectWebAuthInterceptor(this.authInterceptorProvider),
            NewConnectWebOrgInterceptor(this.orgInterceptorProvider),
          ],
        });
        this.userNew = createUserServiceClient(transport);
        this.mgmtNew = createManagementServiceClient(transportOldAPIs);
        this.authNew = createAuthServiceClient(transport);

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
        throw err;
      }),
    );

    return firstValueFrom(init);
  }
}

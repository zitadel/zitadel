import { PlatformLocation } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { TranslateService } from '@ngx-translate/core';
import { AuthConfig } from 'angular-oauth2-oidc';

import { AdminServiceClient } from '../proto/generated/zitadel/AdminServiceClientPb';
import { AuthServiceClient } from '../proto/generated/zitadel/AuthServiceClientPb';
import { ManagementServiceClient } from '../proto/generated/zitadel/ManagementServiceClientPb';
import { AuthenticationService } from './authentication.service';
import { ExhaustedService } from './exhausted.service';
import { AuthInterceptor } from './interceptors/auth.interceptor';
import { ExhaustedInterceptor } from './interceptors/exhausted.interceptor';
import { I18nInterceptor } from './interceptors/i18n.interceptor';
import { OrgInterceptor } from './interceptors/org.interceptor';
import { StorageService } from './storage.service';

@Injectable({
  providedIn: 'root',
})
export class GrpcService {
  public auth!: AuthServiceClient;
  public mgmt!: ManagementServiceClient;
  public admin!: AdminServiceClient;

  constructor(
    private http: HttpClient,
    private platformLocation: PlatformLocation,
    private authenticationService: AuthenticationService,
    private storageService: StorageService,
    private dialog: MatDialog,
    private translate: TranslateService,
    private exhaustedService: ExhaustedService,
  ) {}

  public async loadAppEnvironment(): Promise<any> {
    return this.http
      .get('./assets/environment.json')
      .toPromise()
      .then((data: any) => {
        if (data && data.api && data.issuer) {
          const interceptors = {
            unaryInterceptors: [
              new ExhaustedInterceptor(this.exhaustedService),
              new OrgInterceptor(this.storageService),
              new AuthInterceptor(this.authenticationService, this.storageService, this.dialog),
              new I18nInterceptor(this.translate),
            ],
          };

          this.auth = new AuthServiceClient(
            data.api,
            null,
            // @ts-ignore
            interceptors,
          );
          this.mgmt = new ManagementServiceClient(
            data.api,
            null,
            // @ts-ignore
            interceptors,
          );
          this.admin = new AdminServiceClient(
            data.api,
            null,
            // @ts-ignore
            interceptors,
          );

          const authConfig: AuthConfig = {
            scope: 'openid profile email',
            responseType: 'code',
            oidc: true,
            clientId: data.clientid,
            issuer: data.issuer,
            redirectUri: window.location.origin + this.platformLocation.getBaseHrefFromDOM() + 'auth/callback',
            postLogoutRedirectUri: window.location.origin + this.platformLocation.getBaseHrefFromDOM() + 'signedout',
            requireHttps: false,
          };

          this.authenticationService.initConfig(authConfig);
        }
        return Promise.resolve(data);
      })
      .catch(() => {
        console.error('Failed to load environment from assets');
      });
  }
}

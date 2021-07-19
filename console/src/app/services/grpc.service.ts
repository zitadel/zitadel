import { PlatformLocation } from '@angular/common';
import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { AuthConfig } from 'angular-oauth2-oidc';

import { AdminServiceClient } from '../proto/generated/zitadel/AdminServiceClientPb';
import { AuthServiceClient } from '../proto/generated/zitadel/AuthServiceClientPb';
import { ManagementServiceClient } from '../proto/generated/zitadel/ManagementServiceClientPb';
import { AuthenticationService } from './authentication.service';
import { AuthInterceptor } from './interceptors/auth.interceptor';
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
    private platformLocation: PlatformLocation,
    private authenticationService: AuthenticationService,
    private storageService: StorageService,
    private dialog: MatDialog,
  ) { }

  public initializeGrpc(data: any): Promise<void> {
    return new Promise((resolve, reject) => {
      if (data && data.authServiceUrl && data.mgmtServiceUrl && data.issuer) {
        const interceptors = {
          unaryInterceptors: [
            new OrgInterceptor(this.storageService),
            new AuthInterceptor(this.authenticationService, this.storageService, this.dialog),
            new I18nInterceptor(),
          ],
        };

        this.auth = new AuthServiceClient(
          data.authServiceUrl,
          null,
          // @ts-ignore
          interceptors,
        );
        this.mgmt = new ManagementServiceClient(
          data.mgmtServiceUrl,
          null,
          // @ts-ignore
          interceptors,
        );
        this.admin = new AdminServiceClient(
          // TODO: replace with service url
          data.mgmtServiceUrl,
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
        };

        this.authenticationService.initConfig(authConfig);

        return resolve();
      } else {
        return reject();
      }
    });
  }
}
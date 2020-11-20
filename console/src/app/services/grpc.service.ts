import { PlatformLocation } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { AuthConfig } from 'angular-oauth2-oidc';

import { AdminServicePromiseClient } from '../proto/generated/admin_grpc_web_pb';
import { AuthServicePromiseClient } from '../proto/generated/auth_grpc_web_pb';
import { ManagementServicePromiseClient } from '../proto/generated/management_grpc_web_pb';
import { AuthenticationService } from './authentication.service';
import { AuthInterceptor } from './interceptors/auth.interceptor';
import { I18nInterceptor } from './interceptors/i18n.interceptor';
import { OrgInterceptor } from './interceptors/org.interceptor';
import { StorageService } from './storage.service';

@Injectable({
    providedIn: 'root',
})
export class GrpcService {
    public auth!: AuthServicePromiseClient;
    public mgmt!: ManagementServicePromiseClient;
    public admin!: AdminServicePromiseClient;

    constructor(
        private http: HttpClient,
        private platformLocation: PlatformLocation,
        private authenticationService: AuthenticationService,
        private storageService: StorageService,
        private dialog: MatDialog,
        // private toast: ToastService,
    ) { }

    public async loadAppEnvironment(): Promise<any> {
        return this.http.get('./assets/environment.json')
            .toPromise().then((data: any) => {
                if (data && data.authServiceUrl && data.mgmtServiceUrl && data.issuer) {
                    const interceptors = {
                        unaryInterceptors: [
                            new OrgInterceptor(this.storageService),
                            new AuthInterceptor(this.authenticationService, this.storageService, this.dialog),
                            new I18nInterceptor(),
                        ],
                    };

                    this.auth = new AuthServicePromiseClient(
                        data.authServiceUrl,
                        null,
                        // @ts-ignore
                        interceptors,
                    );
                    this.mgmt = new ManagementServicePromiseClient(
                        data.mgmtServiceUrl,
                        null,
                        // @ts-ignore
                        interceptors,
                    );
                    this.admin = new AdminServicePromiseClient(
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
                }
                return Promise.resolve(data);
            }).catch(() => {
                console.error('Failed to load environment from assets');
            });
    }
}

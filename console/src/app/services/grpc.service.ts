import { PlatformLocation } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { AdminServicePromiseClient } from '../proto/generated/admin_grpc_web_pb';
import { AuthServicePromiseClient } from '../proto/generated/auth_grpc_web_pb';
import { ManagementServicePromiseClient } from '../proto/generated/management_grpc_web_pb';
import { AuthenticationService } from './authentication.service';
import { AuthInterceptor } from './interceptors/auth.interceptor';
import { OrgInterceptor } from './interceptors/org.interceptor';
import { StorageService } from './storage.service';

@Injectable({
    providedIn: 'root',
})
export class GrpcService {
    public issuer: string = '';
    public clientid: string = '';
    public redirectUri: string = '';
    public postLogoutRedirectUri: string = '';

    public auth!: AuthServicePromiseClient;
    public mgmt!: ManagementServicePromiseClient;
    public admin!: AdminServicePromiseClient;

    constructor(
        private http: HttpClient,
        private platformLocation: PlatformLocation,
        private authenticationService: AuthenticationService,
        private storageService: StorageService,
    ) { }

    public async loadAppEnvironment(): Promise<any> {
        return this.http.get('./assets/environment.json')
            .toPromise().then((data: any) => {
                if (data && data.authServiceUrl && data.mgmtServiceUrl && data.issuer) {
                    const interceptors = {
                        'unaryInterceptors': [
                            new AuthInterceptor(this.authenticationService, this.storageService),
                            new OrgInterceptor(this.storageService),
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
                        data.adminServiceUrl,
                        null,
                        // @ts-ignore
                        interceptors,
                    );

                    this.issuer = data.issuer;
                    if (data.clientid) {
                        this.clientid = data.clientid;
                        this.redirectUri = window.location.origin + this.platformLocation.getBaseHrefFromDOM() + 'auth/callback';
                        this.postLogoutRedirectUri = window.location.origin + this.platformLocation.getBaseHrefFromDOM() + 'signedout';
                    }
                }
                return Promise.resolve(data);
            }).catch(() => {
                console.log('Failed to load environment from assets');
            });
    }
}

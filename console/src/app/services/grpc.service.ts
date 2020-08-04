import { PlatformLocation } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { AdminServicePromiseClient } from '../proto/generated/admin_grpc_web_pb';
import { AuthServicePromiseClient } from '../proto/generated/auth_grpc_web_pb';
import { ManagementServicePromiseClient } from '../proto/generated/management_grpc_web_pb';
import { GrpcRequestFn } from './grpc-handler';

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
    ) { }

    public async loadAppEnvironment(): Promise<any> {
        return this.http.get('./assets/environment.json')
            .toPromise().then((data: any) => {
                console.log('init grpc');

                if (data && data.authServiceUrl && data.mgmtServiceUrl && data.issuer) {
                    console.log('init grpc promiseclients');
                    this.auth = new AuthServicePromiseClient(data.authServiceUrl);
                    this.mgmt = new ManagementServicePromiseClient(data.mgmtServiceUrl);
                    this.admin = new AdminServicePromiseClient(data.adminServiceUrl);

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

export type RequestFactory<TClient, TReq, TResp> = (
    client: TClient,
) => GrpcRequestFn<TReq, TResp>;

export type ResponseMapper<TResp, TMappedResp> = (resp: TResp) => TMappedResp;

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

    public auth!: AuthServicePromiseClient;
    public mgmt!: ManagementServicePromiseClient;
    public admin!: AdminServicePromiseClient;

    constructor(
        private http: HttpClient,
    ) { }

    public async loadAppEnvironment(): Promise<any> {
        return this.http.get('/assets/environment.json')
            .toPromise().then((data: any) => {
                if (data && data.authServiceUrl && data.mgmtServiceUrl && data.issuer) {
                    this.auth = new AuthServicePromiseClient(data.authServiceUrl);
                    this.mgmt = new ManagementServicePromiseClient(data.mgmtServiceUrl);
                    this.admin = new AdminServicePromiseClient(data.adminServiceUrl);

                    this.issuer = data.issuer;
                    if (data.clientid) {
                        console.log(data.clientid);
                        this.clientid = data.clientid;
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

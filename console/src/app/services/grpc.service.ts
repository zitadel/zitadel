import { Injectable } from '@angular/core';

import { EnvironmentDep } from '../app.module';
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
    public auth!: AuthServicePromiseClient;
    public mgmt!: ManagementServicePromiseClient;
    public admin!: AdminServicePromiseClient;

    constructor(
        private authenticationService: AuthenticationService,
        private storageService: StorageService,
    ) { }

    public grpcInit = (
        envDeps: /*() => */ Promise<EnvironmentDep>,
    ): () => Promise<any> => {
        return (): Promise<any> => {
            return envDeps.then(data => {
                console.log(data);

                const interceptors = {
                    'unaryInterceptors': [
                        new AuthInterceptor(this.authenticationService),
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
            });
        };
    };
}

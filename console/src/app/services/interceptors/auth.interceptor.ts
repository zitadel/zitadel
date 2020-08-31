import { Injectable } from '@angular/core';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { filter, first } from 'rxjs/operators';

import { AuthenticationService } from '../authentication.service';
import { StorageService } from '../storage.service';


const authorizationKey = 'Authorization';
const bearerPrefix = 'Bearer';
const accessTokenStorageKey = 'access_token';
@Injectable({ providedIn: 'root' })
export class AuthInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
    constructor(private authenticationService: AuthenticationService, private storageService: StorageService) { }

    public async intercept(request: Request<TReq, TResp>, invoker: any): Promise<UnaryResponse<TReq, TResp>> {
        await this.authenticationService.authenticationChanged.pipe(
            filter((authed) => !!authed),
            first(),
        ).toPromise();

        const metadata = request.getMetadata();
        const accessToken = this.storageService.getItem(accessTokenStorageKey);
        metadata[authorizationKey] = `${bearerPrefix} ${accessToken}`;

        return invoker(request).then((response: any) => {
            // const message = response.getResponseMessage();
            const respMetadata = response.getMetadata();

            // TODO: intercept unauthenticated an authenticate

            // const status = respMetadata['grpc-status'];
            // console.log(respMetadata, status);
            // if (status?.code === 16) {
            //     this.authenticationService.authenticate();
            // }
            return response;
        });
    }
}

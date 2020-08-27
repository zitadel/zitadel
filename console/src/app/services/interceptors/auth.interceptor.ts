import { Injectable, InjectionToken, Injector } from '@angular/core';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { filter, first } from 'rxjs/operators';

import { AuthenticationService } from '../authentication.service';


const authorizationKey = 'Authorization';
const bearerPrefix = 'Bearer';
const accessTokenStorageField = 'access_token';

export const GRPC_INTERCEPTORS = new InjectionToken<Array<UnaryInterceptor<any, any>>>(
    'GRPC_INTERCEPTORS',
);

@Injectable({ providedIn: 'root' })
export class AuthInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
    constructor(private injector: Injector) { }

    public async intercept(request: Request<TReq, TResp>, invoker: any): Promise<UnaryResponse<TReq, TResp>> {

        const auth = this.injector.get(AuthenticationService);

        const accessToken = await auth.authenticationChanged.pipe(
            filter((authed) => !!authed),
            first(),
        ).toPromise();

        const metadata = request.getMetadata();
        metadata[authorizationKey] = `${bearerPrefix} ${accessToken}`;

        return invoker(request).then((response: any) => {
            // const message = response.getResponseMessage();
            const respMetadata = response.getMetadata();
            console.log(respMetadata['grpc-status]']);
            return response;
        });
    }
}

import { Injectable } from '@angular/core';
import { UnaryInterceptor } from 'grpc-web';

import { StorageService } from '../storage.service';


const authorizationKey = 'Authorization';
const bearerPrefix = 'Bearer ';
const accessTokenStorageField = 'access_token';

@Injectable({ providedIn: 'root' })
export class AuthInterceptor implements UnaryInterceptor<any, any> {
    constructor(private readonly authStorage: StorageService) { }

    public intercept(request: any, invoker: any): any {
        console.log('authinterceptor');
        // Update the request message before the RPC.
        const reqMsg = request.getRequestMessage();
        console.log(request.getMetadata());

        const accessToken = this.authStorage.getItem(accessTokenStorageField);

        if (accessToken) {
            const metadata = { 'Authorization': bearerPrefix + accessToken };
            console.log(metadata);
            reqMsg.metadata = metadata;
        }

        console.log(request.getMetadata());

        // After the RPC returns successfully, update the response.
        return invoker(request).then((response: any) => {
            return response;
        });
    }
}

import { Injectable } from '@angular/core';
import { UnaryInterceptor } from 'grpc-web';
import { Org } from 'src/app/proto/generated/auth_pb';

import { StorageService } from '../storage.service';

@Injectable({ providedIn: 'root' })
export class OrgInterceptor implements UnaryInterceptor<any, any> {
    constructor(private readonly storageService: StorageService) { }

    public intercept(request: any, invoker: any): any {
        console.log('orginterceptor');
        const reqMsg = request.getRequestMessage();

        const org: Org.AsObject | null = (this.storageService.getItem('organization'));
        console.log(org);

        if (org) {
            reqMsg.metadata = { 'x-zitadel-orgid': org.id };
        }

        // After the RPC returns successfully, update the response.
        return invoker(request).then((response: any) => {
            // You can also do something with response metadata here.
            console.log(response.getMetadata());

            // Update the response message.
            return response;
        });
    }
}

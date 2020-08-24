import { Injectable } from '@angular/core';
import { Org } from 'src/app/proto/generated/auth_pb';

import { StorageService } from '../storage.service';

@Injectable({ providedIn: 'root' })
export class OrgInterceptor {
    constructor(private readonly storageService: StorageService) { }

    public intercept(request: any, invoker: any): any {
        console.log('orginterceptor');
        // Update the request message before the RPC.
        console.log(request);
        const reqMsg = request.getRequestMessage();
        reqMsg.setMessage('[Intercept request]' + reqMsg.getMessage());

        // After the RPC returns successfully, update the response.
        return invoker(request).then((response: any) => {
            // You can also do something with response metadata here.
            console.log(response.getMetadata());

            // Update the response message.
            const responseMsg = response.getResponseMessage();

            const org: Org.AsObject | null = (this.storageService.getItem('organization'));
            console.log(org);
            // if (!response.setMetadata([orgKey] && org) {
            //     metadata[orgKey] = org.id ?? '';
            // }

            responseMsg.setMessage('[Intercept response]' + responseMsg.getMessage());

            return response;
        });
    };
}

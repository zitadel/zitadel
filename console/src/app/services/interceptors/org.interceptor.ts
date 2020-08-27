import { Injectable } from '@angular/core';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { Org } from 'src/app/proto/generated/auth_pb';

import { StorageService } from '../storage.service';


const orgKey = 'x-zitadel-orgid';
@Injectable({ providedIn: 'root' })
export class OrgInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
    constructor(private readonly storageService: StorageService) { }

    public intercept(request: Request<TReq, TResp>, invoker: any): Promise<UnaryResponse<TReq, TResp>> {
        const metadata = request.getMetadata();

        const org: Org.AsObject | null = (this.storageService.getItem('organization'));

        if (org) {
            metadata[orgKey] = `${org.id}`;
        }

        return invoker(request).then((response: any) => {
            return response;
        });
    }
}

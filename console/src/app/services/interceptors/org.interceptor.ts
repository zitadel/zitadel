import { Injectable } from '@angular/core';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';

import { StorageService } from '../storage.service';


const ORG_HEADER_KEY = 'x-zitadel-orgid';
const ORG_STORAGE_KEY = 'organization';
@Injectable({ providedIn: 'root' })
export class OrgInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
    constructor(private readonly storageService: StorageService) { }

    public intercept(request: Request<TReq, TResp>, invoker: any): Promise<UnaryResponse<TReq, TResp>> {
        const metadata = request.getMetadata();

        const org: Org.AsObject | null = (this.storageService.getItem(ORG_STORAGE_KEY));

        if (org) {
            metadata[ORG_HEADER_KEY] = `${org.id}`;
        }

        return invoker(request).then((response: any) => {
            return response;
        }).catch((error: any) => {
            if (error.code === 7 && error.message.startsWith('Organisation doesn\'t exist')) {
                this.storageService.removeItem(ORG_STORAGE_KEY);
            }
            return Promise.reject(error);
        });
    }
}

import { Injectable } from '@angular/core';
import { Metadata } from 'grpc-web';

import { Org } from '../proto/generated/auth_pb';
import { GrpcHandler } from './grpc-handler';
import { GrpcInterceptor } from './grpc-interceptor';
import { StorageService } from './storage.service';

const orgKey = 'x-zitadel-orgid';

@Injectable({ providedIn: 'root' })
export class GrpcOrgInterceptor implements GrpcInterceptor {
    constructor(private readonly storageService: StorageService) { }

    public async intercept(
        req: unknown,
        metadata: Metadata,
        next: GrpcHandler,
    ): Promise<any> {
        const org: Org.AsObject | null = (this.storageService.getItem('organization'));
        if (!metadata[orgKey] && org) {
            metadata[orgKey] = org.id ?? '';
        }

        // metadata['x-zitadel-login'] = 'true';
        // metadata['x-zitadel-userid'] = '58922557365027097';
        // metadata['x-zitadel-orgid'] = '58922556878487833';
        return await next.handle(req, metadata);
    }
}

import { InjectionToken } from '@angular/core';
import { Metadata } from 'grpc-web';

import { GrpcHandler } from '../grpc-handler';

export const GRPC_INTERCEPTORS = new InjectionToken<GrpcInterceptor[]>(
    'GRPC_INTERCEPTORS',
);

export interface GrpcInterceptor {
    intercept(req: unknown, metadata: Metadata, next: GrpcHandler): Promise<any>;
}

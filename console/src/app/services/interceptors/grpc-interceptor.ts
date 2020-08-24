import { InjectionToken } from '@angular/core';

export const GRPC_INTERCEPTORS = new InjectionToken<GrpcInterceptor[]>(
    'GRPC_INTERCEPTORS',
);

export interface GrpcInterceptor {
    intercept(req: unknown, invoker: any): any;
}

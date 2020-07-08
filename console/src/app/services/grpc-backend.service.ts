import { Inject, Injectable } from '@angular/core';
import { Metadata } from 'grpc-web';

import { GrpcDefaultHandler, GrpcHandler, GrpcInterceptorHandler, GrpcRequestFn } from './grpc-handler';
import { GRPC_INTERCEPTORS, GrpcInterceptor } from './interceptors/grpc-interceptor';

@Injectable({
    providedIn: 'root',
})
export class GrpcBackendService {
    constructor(
        @Inject(GRPC_INTERCEPTORS) private readonly interceptors: GrpcInterceptor[],
    ) { }

    public async runRequest<TReq, TResp>(
        requestFn: GrpcRequestFn<TReq, TResp>,
        request: TReq,
        metadata: Metadata = {},
    ): Promise<TResp> {
        const defaultHandler: GrpcHandler = new GrpcDefaultHandler(requestFn);
        const chain = this.interceptors.reduceRight(
            (next, interceptor) => new GrpcInterceptorHandler(next, interceptor),
            defaultHandler,
        );
        return chain.handle(request, metadata);
    }
}

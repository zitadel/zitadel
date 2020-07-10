import { Metadata } from 'grpc-web';

import { GrpcInterceptor } from './interceptors/grpc-interceptor';

export type GrpcRequestFn<TReq, TResp> = (
    request: TReq,
    metadata?: Metadata,
) => Promise<TResp>;

export interface GrpcHandler {
    handle(req: unknown, metadata: Metadata): Promise<any>;
}

export class GrpcInterceptorHandler implements GrpcHandler {
    constructor(
        private readonly next: GrpcHandler,
        private interceptor: GrpcInterceptor,
    ) { }

    public handle(req: unknown, metadata: Metadata): Promise<unknown> {
        return this.interceptor.intercept(req, metadata, this.next);
    }
}

export class GrpcDefaultHandler<TReq, TResp> implements GrpcHandler {
    constructor(private readonly requestFn: GrpcRequestFn<TReq, TResp>) { }

    public handle(req: TReq, metadata: Metadata): Promise<TResp> {
        return this.requestFn(req, metadata);
    }
}

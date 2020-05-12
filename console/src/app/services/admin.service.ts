import { Injectable } from '@angular/core';
import { Metadata } from 'grpc-web';

import { AdminServicePromiseClient } from '../proto/generated/admin_grpc_web_pb';
import { CreateOrgRequest, OrgSetUpRequest, OrgSetUpResponse, RegisterUserRequest } from '../proto/generated/admin_pb';
import { GrpcBackendService } from './grpc-backend.service';
import { GrpcService, RequestFactory, ResponseMapper } from './grpc.service';

@Injectable({
    providedIn: 'root',
})
export class AdminService {
    constructor(private readonly grpcService: GrpcService, private grpcBackendService: GrpcBackendService) { }

    public async request<TReq, TResp, TMappedResp>(
        requestFn: RequestFactory<AdminServicePromiseClient, TReq, TResp>,
        request: TReq,
        responseMapper: ResponseMapper<TResp, TMappedResp>,
        metadata?: Metadata,
    ): Promise<TMappedResp> {
        const mappedRequestFn = requestFn(this.grpcService.admin).bind(this.grpcService.mgmt);
        const response = await this.grpcBackendService.runRequest(
            mappedRequestFn,
            request,
            metadata,
        );
        return responseMapper(response);
    }

    public async SetUpOrg(
        createOrgRequest: CreateOrgRequest,
        registerUserRequest: RegisterUserRequest,
    ): Promise<OrgSetUpResponse> {
        const req: OrgSetUpRequest = new OrgSetUpRequest();

        req.setOrg(createOrgRequest);
        req.setUser(registerUserRequest);

        return await this.request(
            c => c.setUpOrg,
            req,
            f => f,
        );
    }
}

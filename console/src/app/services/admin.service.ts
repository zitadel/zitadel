import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Metadata } from 'grpc-web';

import { AdminServicePromiseClient } from '../proto/generated/admin_grpc_web_pb';
import {
    CreateOrgRequest,
    CreateUserRequest,
    OrgIamPolicy,
    OrgIamPolicyID,
    OrgIamPolicyRequest,
    OrgSetUpRequest,
    OrgSetUpResponse,
} from '../proto/generated/admin_pb';
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
        registerUserRequest: CreateUserRequest,
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

    // public async deleteOrgIamPolicy(
    //     orgId: string,
    // ): Promise<Empty> {
    //     const req = new OrgIamPolicyID();
    //     req.setOrgId(orgId);
    //     return await this.request(
    //         c => c.,
    //         req,
    //         f => f,
    //     );
    // }

    public async GetOrgIamPolicy(orgId: string): Promise<OrgIamPolicy> {
        const req = new OrgIamPolicyID();
        req.setOrgId(orgId);

        return await this.request(
            c => c.getOrgIamPolicy,
            req,
            f => f,
        );
    }

    public async CreateOrgIamPolicy(
        orgId: string,
        description: string,
        userLoginMustBeDomain: boolean): Promise<OrgIamPolicy> {
        const req = new OrgIamPolicyRequest();
        req.setOrgId(orgId);
        req.setDescription(description);
        req.setUserLoginMustBeDomain(userLoginMustBeDomain);

        return await this.request(
            c => c.createOrgIamPolicy,
            req,
            f => f,
        );
    }

    public async UpdateOrgIamPolicy(
        orgId: string,
        description: string,
        userLoginMustBeDomain: boolean): Promise<OrgIamPolicy> {
        const req = new OrgIamPolicyRequest();
        req.setOrgId(orgId);
        req.setDescription(description);
        req.setUserLoginMustBeDomain(userLoginMustBeDomain);

        return await this.request(
            c => c.updateOrgIamPolicy,
            req,
            f => f,
        );
    }

    public async deleteOrgIamPolicy(
        orgId: string,
    ): Promise<Empty> {
        const req = new OrgIamPolicyID();
        req.setOrgId(orgId);
        return await this.request(
            c => c.deleteOrgIamPolicy,
            req,
            f => f,
        );
    }
}

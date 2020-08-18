import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Metadata } from 'grpc-web';

import { AdminServicePromiseClient } from '../proto/generated/admin_grpc_web_pb';
import {
    AddIamMemberRequest,
    ChangeIamMemberRequest,
    CreateOrgRequest,
    CreateUserRequest,
    FailedEventID,
    FailedEvents,
    IamMember,
    IamMemberRoles,
    IamMemberSearchQuery,
    IamMemberSearchRequest,
    IamMemberSearchResponse,
    OrgIamPolicy,
    OrgIamPolicyID,
    OrgIamPolicyRequest,
    OrgSetUpRequest,
    OrgSetUpResponse,
    RemoveIamMemberRequest,
    ViewID,
    Views,
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

    public async GetIamMemberRoles(): Promise<IamMemberRoles> {
        return await this.request(
            c => c.getIamMemberRoles,
            new Empty(),
            f => f,
        );
    }

    public async GetViews(): Promise<Views> {
        return await this.request(
            c => c.getViews,
            new Empty(),
            f => f,
        );
    }

    public async GetFailedEvents(): Promise<FailedEvents> {
        return await this.request(
            c => c.getFailedEvents,
            new Empty(),
            f => f,
        );
    }

    public async ClearView(viewname: string, db: string): Promise<Empty> {
        const req: ViewID = new ViewID();
        req.setDatabase(db);
        req.setViewName(viewname);
        return await this.request(
            c => c.clearView,
            req,
            f => f,
        );
    }

    public async RemoveFailedEvent(viewname: string, db: string, sequence: number): Promise<Empty> {
        const req: FailedEventID = new FailedEventID();
        req.setDatabase(db);
        req.setViewName(viewname);
        req.setFailedSequence(sequence);
        return await this.request(
            c => c.removeFailedEvent,
            req,
            f => f,
        );
    }

    public async SearchIamMembers(
        limit: number,
        offset: number,
        queryList?: IamMemberSearchQuery[],
    ): Promise<IamMemberSearchResponse> {
        const req = new IamMemberSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchIamMembers,
            req,
            f => f,
        );
    }

    public async RemoveIamMember(
        userId: string,
    ): Promise<Empty> {
        const req = new RemoveIamMemberRequest();
        req.setUserId(userId);

        return await this.request(
            c => c.removeIamMember,
            req,
            f => f,
        );
    }

    public async AddIamMember(
        userId: string,
        rolesList: string[],
    ): Promise<IamMember> {
        const req = new AddIamMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);

        return await this.request(
            c => c.addIamMember,
            req,
            f => f,
        );
    }

    public async ChangeIamMember(
        userId: string,
        rolesList: string[],
    ): Promise<IamMember> {
        const req = new ChangeIamMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);

        return await this.request(
            c => c.changeIamMember,
            req,
            f => f,
        );
    }

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

import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

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
import { GrpcService } from './grpc.service';

@Injectable({
    providedIn: 'root',
})
export class AdminService {
    constructor(private readonly grpcService: GrpcService) { }

    public async SetUpOrg(
        createOrgRequest: CreateOrgRequest,
        registerUserRequest: CreateUserRequest,
    ): Promise<OrgSetUpResponse> {
        const req: OrgSetUpRequest = new OrgSetUpRequest();

        req.setOrg(createOrgRequest);
        req.setUser(registerUserRequest);

        return this.grpcService.admin.setUpOrg(req);
    }

    public async GetIamMemberRoles(): Promise<IamMemberRoles> {
        const req = new Empty();
        return this.grpcService.admin.getIamMemberRoles(req);
    }

    public async GetViews(): Promise<Views> {
        const req = new Empty();
        return this.grpcService.admin.getViews(req);
    }

    public async GetFailedEvents(): Promise<FailedEvents> {
        const req = new Empty();
        return this.grpcService.admin.getFailedEvents(req);
    }

    public async ClearView(viewname: string, db: string): Promise<Empty> {
        const req: ViewID = new ViewID();
        req.setDatabase(db);
        req.setViewName(viewname);
        return this.grpcService.admin.clearView(req);
    }

    public async RemoveFailedEvent(viewname: string, db: string, sequence: number): Promise<Empty> {
        const req: FailedEventID = new FailedEventID();
        req.setDatabase(db);
        req.setViewName(viewname);
        req.setFailedSequence(sequence);
        return this.grpcService.admin.removeFailedEvent(req);
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
        return this.grpcService.admin.searchIamMembers(req);
    }

    public async RemoveIamMember(
        userId: string,
    ): Promise<Empty> {
        const req = new RemoveIamMemberRequest();
        req.setUserId(userId);

        return this.grpcService.admin.removeIamMember(req);
    }

    public async AddIamMember(
        userId: string,
        rolesList: string[],
    ): Promise<IamMember> {
        const req = new AddIamMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);

        return this.grpcService.admin.addIamMember(req);
    }

    public async ChangeIamMember(
        userId: string,
        rolesList: string[],
    ): Promise<IamMember> {
        const req = new ChangeIamMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);

        return this.grpcService.admin.changeIamMember(req);
    }

    public async GetOrgIamPolicy(orgId: string): Promise<OrgIamPolicy> {
        const req = new OrgIamPolicyID();
        req.setOrgId(orgId);

        return this.grpcService.admin.getOrgIamPolicy(req);
    }

    public async CreateOrgIamPolicy(
        orgId: string,
        description: string,
        userLoginMustBeDomain: boolean): Promise<OrgIamPolicy> {
        const req = new OrgIamPolicyRequest();
        req.setOrgId(orgId);
        req.setDescription(description);
        req.setUserLoginMustBeDomain(userLoginMustBeDomain);

        return this.grpcService.admin.createOrgIamPolicy(req);
    }

    public async UpdateOrgIamPolicy(
        orgId: string,
        description: string,
        userLoginMustBeDomain: boolean): Promise<OrgIamPolicy> {
        const req = new OrgIamPolicyRequest();
        req.setOrgId(orgId);
        req.setDescription(description);
        req.setUserLoginMustBeDomain(userLoginMustBeDomain);
        return this.grpcService.admin.updateOrgIamPolicy(req);
    }

    public async deleteOrgIamPolicy(
        orgId: string,
    ): Promise<Empty> {
        const req = new OrgIamPolicyID();
        req.setOrgId(orgId);
        return this.grpcService.admin.deleteOrgIamPolicy(req);
    }
}

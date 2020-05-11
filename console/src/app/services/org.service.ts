import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Metadata } from 'grpc-web';

import { ManagementServicePromiseClient } from '../proto/generated/management_grpc_web_pb';
import {
    AddOrgMemberRequest,
    Org,
    OrgDomain,
    OrgID,
    OrgMemberRoles,
    OrgMemberSearchRequest,
    OrgMemberSearchResponse,
    PasswordAgePolicy,
    PasswordAgePolicyCreate,
    PasswordAgePolicyID,
    PasswordAgePolicyUpdate,
    PasswordComplexityPolicy,
    PasswordComplexityPolicyCreate,
    PasswordComplexityPolicyID,
    PasswordComplexityPolicyUpdate,
    PasswordLockoutPolicy,
    PasswordLockoutPolicyCreate,
    PasswordLockoutPolicyID,
    PasswordLockoutPolicyUpdate,
    ProjectGrant,
    ProjectGrantCreate,
    RemoveOrgMemberRequest,
} from '../proto/generated/management_pb';
import { GrpcBackendService } from './grpc-backend.service';
import { GrpcService, RequestFactory, ResponseMapper } from './grpc.service';

@Injectable({
    providedIn: 'root',
})
export class OrgService {
    constructor(private readonly grpcService: GrpcService, private grpcBackendService: GrpcBackendService) { }

    public async request<TReq, TResp, TMappedResp>(
        requestFn: RequestFactory<ManagementServicePromiseClient, TReq, TResp>,
        request: TReq,
        responseMapper: ResponseMapper<TResp, TMappedResp>,
        metadata?: Metadata,
    ): Promise<TMappedResp> {
        const mappedRequestFn = requestFn(this.grpcService.mgmt).bind(this.grpcService.mgmt);
        const response = await this.grpcBackendService.runRequest(
            mappedRequestFn,
            request,
            metadata,
        );
        return responseMapper(response);
    }

    public async GetOrgById(orgId: string): Promise<Org> {
        const req: OrgID = new OrgID();
        req.setId(orgId);
        return await this.request(
            c => c.getOrgByID,
            req,
            f => f,
        );
    }

    public async SearchOrgMembers(orgId: string, limit: number, offset: number): Promise<OrgMemberSearchResponse> {
        const req = new OrgMemberSearchRequest();
        req.setOrgId(orgId);
        req.setLimit(limit);
        req.setOffset(offset);
        return await this.request(
            c => c.searchOrgMembers,
            req,
            f => f,
        );
    }

    public async getOrgByDomainGlobal(domain: string): Promise<Org> {
        const req = new OrgDomain();
        req.setDomain(domain);
        return await this.request(
            c => c.getOrgByDomainGlobal,
            req,
            f => f,
        );
    }

    public async AddOrgMember(orgId: string, userId: string, rolesList: string[]): Promise<Empty> {
        const req = new AddOrgMemberRequest();
        req.setOrgId(orgId);
        req.setUserId(userId);
        if (rolesList) {
            req.setRolesList(rolesList);
        }
        return await this.request(
            c => c.addOrgMember,
            req,
            f => f,
        );
    }

    public async RemoveOrgMember(orgId: string, userId: string): Promise<Empty> {
        const req = new RemoveOrgMemberRequest();
        req.setOrgId(orgId);
        req.setUserId(userId);
        return await this.request(
            c => c.removeOrgMember,
            req,
            f => f,
        );
    }

    public async DeactivateOrg(id: string): Promise<Org> {
        const req = new OrgID();
        req.setId(id);
        return await this.request(
            c => c.deactivateOrg,
            req,
            f => f,
        );
    }

    public async ReactivateOrg(id: string): Promise<Org> {
        const req = new OrgID();
        req.setId(id);
        return await this.request(
            c => c.reactivateOrg,
            req,
            f => f,
        );
    }

    public async CreateProjectGrant(
        projectId: string,
        orgId: string,
        roleNamesList: string[],
    ): Promise<ProjectGrant> {
        const req = new ProjectGrantCreate();
        req.setProjectId(projectId);
        req.setGrantedOrgId(orgId);
        req.setRoleNamesList(roleNamesList);
        return await this.request(
            c => c.createProjectGrant,
            req,
            f => f,
        );
    }

    public async GetOrgMemberRoles(): Promise<OrgMemberRoles> {
        const req = new Empty();
        return await this.request(
            c => c.getOrgMemberRoles,
            req,
            f => f,
        );
    }

    // Policy

    public async GetPasswordAgePolicy(): Promise<PasswordAgePolicy> {
        const req = new Empty();

        return await this.request(
            c => c.getPasswordAgePolicy,
            req,
            f => f,
        );
    }

    public async CreatePasswordAgePolicy(
        description: string,
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<PasswordAgePolicy> {
        const req = new PasswordAgePolicyCreate();
        req.setDescription(description);
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);

        return await this.request(
            c => c.createPasswordAgePolicy,
            req,
            f => f,
        );
    }

    public async DeletePasswordAgePolicy(id: string): Promise<Empty> {
        const req = new PasswordAgePolicyID();
        req.setId(id);
        return await this.request(
            c => c.deletePasswordAgePolicy,
            req,
            f => f,
        );
    }

    public async UpdatePasswordAgePolicy(
        description: string,
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<PasswordAgePolicy> {
        const req = new PasswordAgePolicyUpdate();
        req.setDescription(description);
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);
        return await this.request(
            c => c.updatePasswordAgePolicy,
            req,
            f => f,
        );
    }

    public async GetPasswordComplexityPolicy(): Promise<PasswordComplexityPolicy> {
        const req = new Empty();

        return await this.request(
            c => c.getPasswordComplexityPolicy,
            req,
            f => f,
        );
    }

    public async CreatePasswordComplexityPolicy(
        description: string,
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<PasswordComplexityPolicy> {
        const req = new PasswordComplexityPolicyCreate();
        req.setDescription(description);
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return await this.request(
            c => c.createPasswordComplexityPolicy,
            req,
            f => f,
        );
    }

    public async DeletePasswordComplexityPolicy(id: string): Promise<Empty> {
        const req = new PasswordComplexityPolicyID();
        req.setId(id);
        return await this.request(
            c => c.deletePasswordComplexityPolicy,
            req,
            f => f,
        );
    }

    public async UpdatePasswordComplexityPolicy(
        description: string,
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<PasswordComplexityPolicy> {
        const req = new PasswordComplexityPolicyUpdate();
        req.setDescription(description);
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return await this.request(
            c => c.updatePasswordComplexityPolicy,
            req,
            f => f,
        );
    }

    public async GetPasswordLockoutPolicy(): Promise<PasswordLockoutPolicy> {
        const req = new Empty();

        return await this.request(
            c => c.getPasswordLockoutPolicy,
            req,
            f => f,
        );
    }

    public async CreatePasswordLockoutPolicy(
        description: string,
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<PasswordLockoutPolicy> {
        const req = new PasswordLockoutPolicyCreate();
        req.setDescription(description);
        req.setMaxAttempts(maxAttempts);
        req.setShowLockOutFailures(showLockoutFailures);

        return await this.request(
            c => c.createPasswordLockoutPolicy,
            req,
            f => f,
        );
    }

    public async DeletePasswordLockoutPolicy(id: string): Promise<Empty> {
        const req = new PasswordLockoutPolicyID();
        req.setId(id);

        return await this.request(
            c => c.deletePasswordLockoutPolicy,
            req,
            f => f,
        );
    }

    public async UpdatePasswordLockoutPolicy(
        description: string,
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<PasswordLockoutPolicy> {
        const req = new PasswordLockoutPolicyUpdate();
        req.setDescription(description);
        req.setMaxAttempts(maxAttempts);
        req.setShowLockOutFailures(showLockoutFailures);
        return await this.request(
            c => c.updatePasswordLockoutPolicy,
            req,
            f => f,
        );
    }
}

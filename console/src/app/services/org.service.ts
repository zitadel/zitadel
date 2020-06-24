import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Metadata } from 'grpc-web';

import { ManagementServicePromiseClient } from '../proto/generated/management_grpc_web_pb';
import {
    AddOrgDomainRequest,
    AddOrgMemberRequest,
    Iam,
    Org,
    OrgDomain,
    OrgDomainSearchQuery,
    OrgDomainSearchRequest,
    OrgDomainSearchResponse,
    OrgIamPolicy,
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
    RemoveOrgDomainRequest,
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

    public async GetIam(): Promise<Iam> {
        const req: Empty = new Empty();
        return await this.request(
            c => c.getIam,
            req,
            f => f,
        );
    }

    public async GetMyOrg(): Promise<Org> {
        return await this.request(
            c => c.getMyOrg,
            new Empty(),
            f => f,
        );
    }

    public async AddMyOrgDomain(domain: string): Promise<OrgDomain> {
        const req: AddOrgDomainRequest = new AddOrgDomainRequest();
        req.setDomain(domain);
        return await this.request(
            c => c.addMyOrgDomain,
            req,
            f => f,
        );
    }

    public async RemoveMyOrgDomain(domain: string): Promise<Empty> {
        const req: RemoveOrgDomainRequest = new RemoveOrgDomainRequest();
        req.setDomain(domain);
        return await this.request(
            c => c.removeMyOrgDomain,
            req,
            f => f,
        );
    }

    public async SearchMyOrgDomains(offset: number, limit: number, queryList?: OrgDomainSearchQuery[]):
        Promise<OrgDomainSearchResponse> {
        const req: OrgDomainSearchRequest = new OrgDomainSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }

        return await this.request(
            c => c.searchMyOrgDomains,
            req,
            f => f,
        );
    }

    public async SearchMyOrgMembers(limit: number, offset: number): Promise<OrgMemberSearchResponse> {
        const req = new OrgMemberSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        return await this.request(
            c => c.searchMyOrgMembers,
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

    public async AddMyOrgMember(userId: string, rolesList: string[]): Promise<Empty> {
        const req = new AddOrgMemberRequest();
        req.setUserId(userId);
        if (rolesList) {
            req.setRolesList(rolesList);
        }
        return await this.request(
            c => c.addMyOrgMember,
            req,
            f => f,
        );
    }

    public async RemoveMyOrgMember(userId: string): Promise<Empty> {
        const req = new RemoveOrgMemberRequest();
        req.setUserId(userId);
        return await this.request(
            c => c.removeMyOrgMember,
            req,
            f => f,
        );
    }

    public async DeactivateMyOrg(): Promise<Org> {
        return await this.request(
            c => c.deactivateMyOrg,
            new Empty(),
            f => f,
        );
    }

    public async ReactivateMyOrg(): Promise<Org> {
        const req = new OrgID();
        return await this.request(
            c => c.reactivateMyOrg,
            new Empty(),
            f => f,
        );
    }

    public async CreateProjectGrant(
        projectId: string,
        orgId: string,
        roleKeysList: string[],
    ): Promise<ProjectGrant> {
        const req = new ProjectGrantCreate();
        req.setProjectId(projectId);
        req.setGrantedOrgId(orgId);
        req.setRoleKeysList(roleKeysList);
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

    public async GetMyOrgIamPolicy(): Promise<OrgIamPolicy> {
        return await this.request(
            c => c.getMyOrgIamPolicy,
            new Empty(),
            f => f,
        );
    }

    public async CreateMyOrgIamPolicy(
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

    public async DeleteMyOrgIamPolicy(id: string): Promise<Empty> {
        const req = new PasswordAgePolicyID();
        req.setId(id);
        return await this.request(
            c => c.deletePasswordAgePolicy,
            req,
            f => f,
        );
    }

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
        return await this.request(
            c => c.getPasswordComplexityPolicy,
            new Empty(),
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

    public getLocalizedComplexityPolicyPatternErrorString(policy: PasswordComplexityPolicy.AsObject): string {
        if (policy.hasNumber && policy.hasSymbol) {
            return 'ORG.POLICY.PWD_COMPLEXITY.SYMBOLANDNUMBERERROR';
        } else if (policy.hasNumber) {
            return 'ORG.POLICY.PWD_COMPLEXITY.NUMBERERROR';
        } else if (policy.hasSymbol) {
            return 'ORG.POLICY.PWD_COMPLEXITY.SYMBOLERROR';
        } else {
            return 'ORG.POLICY.PWD_COMPLEXITY.PATTERNERROR';
        }
    }
}

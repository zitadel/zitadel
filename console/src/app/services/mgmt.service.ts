import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject } from 'rxjs';

import { IDPQuery } from '../proto/generated/zitadel/idp_pb';
import {
    AddMachineKeyRequest,
    AddMachineKeyResponse,
    AddOrgDomainRequest,
    AddOrgMemberRequest,
    ListHumanPasswordlessRequest,
    ListHumanPasswordlessResponse,
    ListLoginPolicyMultiFactorsResponse,
    ListOrgIDPsRequest,
    ListOrgIDPsResponse,
    RemoveHumanPasswordlessRequest,
    RemoveHumanPasswordlessResponse,
    RemoveOrgDomainRequest,
    RemoveOrgMemberRequest,
    UpdateMachineRequest,
    ListLoginPolicyMultiFactorsRequest,
    ValidateOrgDomainRequest,
    AddMultiFactorToLoginPolicyRequest,
    AddMultiFactorToLoginPolicyResponse,
    RemoveMultiFactorFromLoginPolicyRequest,
    RemoveMultiFactorFromLoginPolicyResponse,
    ListLoginPolicySecondFactorsResponse,
    AddSecondFactorToLoginPolicyResponse,
    AddSecondFactorToLoginPolicyRequest,
    UpdateOrgIDPResponse,
    AddOrgOIDCIDPRequest,
    RemoveSecondFactorFromLoginPolicyRequest,
    RemoveSecondFactorFromLoginPolicyResponse,
    GetLoginPolicyResponse,
    GetLoginPolicyRequest,
    UpdateCustomLoginPolicyRequest,
    UpdateCustomLoginPolicyResponse,
    GetOrgIDPByIDRequest,
    GetOrgIDPByIDResponse,
    AddCustomLoginPolicyRequest,
    AddCustomLoginPolicyResponse,
    ListMachineKeysRequest,
    ListMachineKeysResponse,
    ResetLoginPolicyToDefaultRequest,
    ResetLoginPolicyToDefaultResponse,
    AddIDPToLoginPolicyRequest,
    AddIDPToLoginPolicyResponse,
    RemoveIDPFromLoginPolicyRequest,
    ListLoginPolicyIDPsRequest,
    ListLoginPolicyIDPsResponse,
    UpdateOrgIDPRequest,
    AddOrgOIDCIDPResponse,
    UpdateOrgIDPOIDCConfigRequest,
    RemoveOrgIDPRequest,
    UpdateOrgIDPOIDCConfigResponse,
    RemoveOrgIDPResponse,
    ReactivateOrgIDPRequest,
    DeactivateOrgIDPRequest,
    AddHumanUserResponse,
    AddHumanUserRequest,
    AddMachineUserRequest,
    AddMachineUserResponse,
    UpdateMachineResponse,
    RemoveMachineKeyRequest,
    RemoveMachineKeyResponse,
    RemoveUserIDPRequest,
    RemoveUserIDPResponse,
    ListUserIDPsRequest,
    ListUserIDPsResponse,
    GetIAMResponse,
    GetIAMRequest,
    GetDefaultPasswordComplexityPolicyResponse,
    GetDefaultPasswordComplexityPolicyRequest,
    GetMyOrgRequest,
    GetMyOrgResponse,
    AddOrgDomainResponse,
    RemoveOrgDomainResponse,
    ListOrgDomainsRequest,
    ListOrgDomainsResponse,
    SetPrimaryOrgDomainRequest,
    SetPrimaryOrgDomainResponse,
    GenerateOrgDomainValidationResponse,
    GenerateOrgDomainValidationRequest,
    ValidateOrgDomainResponse,
    ListOrgMembersRequest,
    ListOrgMembersResponse,
    GetOrgByDomainGlobalResponse,
    GetOrgByDomainGlobalRequest,
    AddOrgResponse,
    AddOrgRequest,
    UpdateOrgMemberResponse,
    UpdateOrgMemberRequest,
    RemoveOrgMemberResponse,
    DeactivateOrgResponse,
    DeactivateOrgRequest,
    ReactivateOrgResponse,
    ReactivateOrgRequest,
    AddProjectGrantRequest,
    AddProjectGrantResponse,
    ListOrgMemberRolesResponse,
    ListOrgMemberRolesRequest,
    GetOrgIAMPolicyRequest,
    GetOrgIAMPolicyResponse,
    GetPasswordAgePolicyResponse,
    GetPasswordAgePolicyRequest,
    AddCustomPasswordAgePolicyRequest,
    AddCustomPasswordAgePolicyResponse,
    ResetPasswordAgePolicyToDefaultRequest,
    ResetPasswordAgePolicyToDefaultResponse,
    UpdateCustomPasswordAgePolicyRequest,
    UpdateCustomPasswordAgePolicyResponse,
    GetPasswordComplexityPolicyResponse,
    GetPasswordComplexityPolicyRequest,
    AddCustomPasswordComplexityPolicyRequest,
    AddCustomPasswordComplexityPolicyResponse,
    ResetPasswordComplexityPolicyToDefaultResponse,
    ResetPasswordComplexityPolicyToDefaultRequest,
    UpdateCustomPasswordComplexityPolicyResponse,
    UpdateCustomPasswordComplexityPolicyRequest,
    GetPasswordLockoutPolicyResponse,
    GetPasswordLockoutPolicyRequest,
    AddCustomPasswordLockoutPolicyRequest,
    AddCustomPasswordLockoutPolicyResponse,
    ResetPasswordLockoutPolicyToDefaultRequest,
    ResetPasswordLockoutPolicyToDefaultResponse,
    UpdateCustomPasswordLockoutPolicyResponse,
    UpdateCustomPasswordLockoutPolicyRequest,
    GetUserByIDRequest,
    GetUserByIDResponse,
    RemoveUserRequest,
    RemoveUserResponse,
    ListProjectMembersRequest,
    ListProjectMembersResponse,
    ListUserMembershipsRequest,
    ListUserMembershipsResponse,
    GetHumanProfileResponse,
    GetHumanProfileRequest,
    ListUserMultiFactorsResponse,
    ListUserMultiFactorsRequest,
    RemoveHumanMultiFactorOTPResponse,
    RemoveHumanMultiFactorOTPRequest,
    RemoveHumanMultiFactorU2FRequest,
    RemoveHumanMultiFactorU2FResponse,
    UpdateHumanProfileRequest,
    UpdateHumanProfileResponse,
    GetHumanEmailResponse,
    GetHumanEmailRequest,
    UpdateHumanEmailResponse,
    UpdateHumanEmailRequest,
    GetHumanPhoneResponse,
    GetHumanPhoneRequest,
    UpdateHumanPhoneResponse,
    UpdateHumanPhoneRequest,
    RemoveHumanPhoneRequest,
    DeactivateUserRequest,
    DeactivateUserResponse,
    AddUserGrantRequest,
    AddUserGrantResponse,
    ReactivateUserResponse,
    ReactivateUserRequest,
    AddProjectRoleRequest
} from '../proto/generated/zitadel/management_pb';
import { KeyType } from '../proto/generated/zitadel/auth_n_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { GrpcService } from './grpc.service';
import { Metadata } from 'grpc';
import { DomainSearchQuery, DomainValidationType } from '../proto/generated/zitadel/org_pb';
import { PasswordComplexityPolicy } from '../proto/generated/management_pb';
import { Member, SearchQuery } from '../proto/generated/zitadel/member_pb';
import { Gender, MembershipQuery } from '../proto/generated/zitadel/user_pb';

export type ResponseMapper<TResp, TMappedResp> = (resp: TResp) => TMappedResp;

@Injectable({
    providedIn: 'root',
})
export class ManagementService {
    public ownedProjectsCount: BehaviorSubject<number> = new BehaviorSubject(0);
    public grantedProjectsCount: BehaviorSubject<number> = new BehaviorSubject(0);

    constructor(private readonly grpcService: GrpcService) { }

    public listOrgIDPs(
        limit?: number,
        offset?: number,
        queryList?: IDPQuery[],
    ): Promise<ListOrgIDPsResponse> {
        const req = new ListOrgIDPsRequest();
        const metadata = new ListQuery();

        if (limit) {
            metadata.setLimit(limit);
        }
        if (offset) {
            metadata.setOffset(offset);
        }
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.listOrgIDPs(req);
    }

    public listHumanPasswordless(userId: string): Promise<ListHumanPasswordlessResponse> {
        const req = new ListHumanPasswordlessRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.listHumanPasswordless(req);
    }

    public removeHumanPasswordless(tokenId: string, userId: string): Promise<RemoveHumanPasswordlessResponse> {
        const req = new RemoveHumanPasswordlessRequest();
        req.setTokenId(tokenId);
        req.setUserId(userId);
        return this.grpcService.mgmt.removeHumanPasswordless(req);
    }

    public listLoginPolicyMultiFactors(): Promise<ListLoginPolicyMultiFactorsResponse> {
        const req = new ListLoginPolicyMultiFactorsRequest();
        return this.grpcService.mgmt.listLoginPolicyMultiFactors(req);
    }

    public addMultiFactorToLoginPolicy(req: AddMultiFactorToLoginPolicyRequest): Promise<AddMultiFactorToLoginPolicyResponse> {
        return this.grpcService.mgmt.addMultiFactorToLoginPolicy(req);
    }

    public removeMultiFactorFromLoginPolicy(req: RemoveMultiFactorFromLoginPolicyRequest): Promise<RemoveMultiFactorFromLoginPolicyResponse> {
        return this.grpcService.mgmt.removeMultiFactorFromLoginPolicy(req);
    }

    public listLoginPolicySecondFactors(): Promise<ListLoginPolicySecondFactorsResponse> {
        const req = new Empty();
        return this.grpcService.mgmt.listLoginPolicySecondFactors(req);
    }

    public addSecondFactorToLoginPolicy(req: AddSecondFactorToLoginPolicyRequest): Promise<AddSecondFactorToLoginPolicyResponse> {
        return this.grpcService.mgmt.addSecondFactorToLoginPolicy(req);
    }

    public removeSecondFactorFromLoginPolicy(req: RemoveSecondFactorFromLoginPolicyRequest): Promise<RemoveSecondFactorFromLoginPolicyResponse> {
        return this.grpcService.mgmt.removeSecondFactorFromLoginPolicy(req);
    }

    public getLoginPolicy(): Promise<GetLoginPolicyResponse> {
        const req = new GetLoginPolicyRequest();
        return this.grpcService.mgmt.getLoginPolicy(req);
    }

    public updateCustomLoginPolicy(req: UpdateCustomLoginPolicyRequest): Promise<UpdateCustomLoginPolicyResponse> {
        return this.grpcService.mgmt.updateCustomLoginPolicy(req);
    }

    public addCustomLoginPolicy(req: AddCustomLoginPolicyRequest): Promise<AddCustomLoginPolicyResponse> {
        return this.grpcService.mgmt.addCustomLoginPolicy(req);
    }

    public resetLoginPolicyToDefault(): Promise<ResetLoginPolicyToDefaultResponse> {
        return this.grpcService.mgmt.resetLoginPolicyToDefault(new ResetLoginPolicyToDefaultRequest());
    }

    public addIDPToLoginPolicy(idpId: string): Promise<AddIDPToLoginPolicyResponse> {
        const req = new AddIDPToLoginPolicyRequest();
        req.setIdpId(idpId);
        return this.grpcService.mgmt.addIDPToLoginPolicy(req);
    }

    public removeIDPFromLoginPolicy(idpId: string): Promise<Empty> {
        const req = new RemoveIDPFromLoginPolicyRequest();
        req.setIdpId(idpId);
        return this.grpcService.mgmt.removeIDPFromLoginPolicy(req);
    }

    public listLoginPolicyIDPs(limit?: number, offset?: number): Promise<ListLoginPolicyIDPsResponse> {
        const req = new ListLoginPolicyIDPsRequest();
        const metadata = new ListQuery();
        if (limit) {
            metadata.setLimit(limit);
        }
        if (offset) {
            metadata.setOffset(offset);
        }
        return this.grpcService.mgmt.listLoginPolicyIDPs(req);
    }

    public getOrgIDPByID(
        id: string,
    ): Promise<GetOrgIDPByIDResponse> {
        const req = new GetOrgIDPByIDRequest();
        req.setId(id);
        return this.grpcService.mgmt.getOrgIDPByID(req);
    }

    public updateOrgIDP(
        req: UpdateOrgIDPRequest,
    ): Promise<UpdateOrgIDPResponse> {
        return this.grpcService.mgmt.updateOrgIDP(req);
    }

    public addOrgOIDCIDP(
        req: AddOrgOIDCIDPRequest,
    ): Promise<AddOrgOIDCIDPResponse> {
        return this.grpcService.mgmt.addOrgOIDCIDP(req);
    }

    public updateOrgIDPOIDCConfig(
        req: UpdateOrgIDPOIDCConfigRequest,
    ): Promise<UpdateOrgIDPOIDCConfigResponse> {
        return this.grpcService.mgmt.updateOrgIDPOIDCConfig(req);
    }

    public removeOrgIDP(
        idpId: string,
    ): Promise<RemoveOrgIDPResponse> {
        const req = new RemoveOrgIDPRequest();
        req.setIdpId(idpId);
        return this.grpcService.mgmt.removeOrgIDP(req);
    }

    public deactivateOrgIDP(
        idpId: string,
    ): Promise<Empty> {
        const req = new DeactivateOrgIDPRequest();
        req.setIdpId(idpId);
        return this.grpcService.mgmt.deactivateOrgIDP(req);
    }

    public reactivateOrgIDP(
        idpId: string,
    ): Promise<Empty> {
        const req = new ReactivateOrgIDPRequest();
        req.setIdpId(idpId);
        return this.grpcService.mgmt.reactivateOrgIDP(req);
    }

    public addHumanUser(request: AddHumanUserRequest): Promise<AddHumanUserResponse> {
        return this.grpcService.mgmt.addHumanUser(request);
    }

    public addMachineUser(request: AddMachineUserRequest): Promise<AddMachineUserResponse> {
        return this.grpcService.mgmt.addMachineUser(request);
    }

    public updateMachine(
        userId: string,
        name?: string,
        description?: string,
    ): Promise<UpdateMachineResponse> {
        const req = new UpdateMachineRequest();
        req.setUserId(userId);
        if (name) {
            req.setName(name);
        }
        if (description) {
            req.setDescription(description);
        }
        return this.grpcService.mgmt.updateMachine(req);
    }

    public addMachineKey(
        userId: string,
        type: KeyType,
        date?: Timestamp,
    ): Promise<AddMachineKeyResponse> {
        const req = new AddMachineKeyRequest();
        req.setType(type);
        req.setUserId(userId);
        if (date) {
            req.setExpirationDate(date);
        }
        return this.grpcService.mgmt.addMachineKey(req);
    }

    public removeMachineKey(
        keyId: string,
        userId: string,
    ): Promise<RemoveMachineKeyResponse> {
        const req = new RemoveMachineKeyRequest();
        req.setKeyId(keyId);
        req.setUserId(userId);

        return this.grpcService.mgmt.removeMachineKey(req);
    }

    public listMachineKeys(
        userId: string,
        limit?: number,
        offset?: number,
        asc?: boolean,
    ): Promise<ListMachineKeysResponse> {
        const req = new ListMachineKeysRequest();
        const metadata = new ListQuery();
        req.setUserId(userId);
        if (limit) {
            metadata.setLimit(limit);
        }
        if (offset) {
            metadata.setOffset(offset);
        }
        if (asc) {
            metadata.setAsc(asc);
        }
        req.setMetaData(metadata);
        return this.grpcService.mgmt.listMachineKeys(req);
    }

    public removeUserIDP(
        idpId: string,
        userId: string,
        linkedUserId: string,
    ): Promise<RemoveUserIDPResponse> {
        const req = new RemoveUserIDPRequest();
        req.setUserId(userId);
        req.setIdpId(idpId);
        req.setUserId(userId);
        req.setLinkedUserId(linkedUserId);
        return this.grpcService.mgmt.removeUserIDP(req);
    }

    public listUserIDPs(
        userId: string,
        limit?: number,
        offset?: number,
    ): Promise<ListUserIDPsResponse> {
        const req = new ListUserIDPsRequest();
        const metadata = new ListQuery();
        req.setUserId(userId);
        if (limit) {
            metadata.setLimit(limit);
        }
        if (offset) {
            metadata.setOffset(offset);
        }
        req.setMetaData(metadata);
        return this.grpcService.mgmt.listUserIDPs(req);
    }

    public getIAM(): Promise<GetIAMResponse> {
        const req = new GetIAMRequest();
        return this.grpcService.mgmt.getIAM(req);
    }

    public getDefaultPasswordComplexityPolicy(): Promise<GetDefaultPasswordComplexityPolicyResponse> {
        const req = new GetDefaultPasswordComplexityPolicyRequest();
        return this.grpcService.mgmt.getDefaultPasswordComplexityPolicy(req);
    }

    public getMyOrg(): Promise<GetMyOrgResponse> {
        const req = new GetMyOrgRequest();
        return this.grpcService.mgmt.getMyOrg(req);
    }

    public addOrgDomain(domain: string): Promise<AddOrgDomainResponse> {
        const req = new AddOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.addOrgDomain(req);
    }

    public removeOrgDomain(domain: string): Promise<RemoveOrgDomainResponse> {
        const req = new RemoveOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.removeOrgDomain(req);
    }

    public listOrgDomains(queryList?: DomainSearchQuery[]):
        Promise<ListOrgDomainsResponse> {
        const req: ListOrgDomainsRequest = new ListOrgDomainsRequest();
        // const metadata= new ListQuery();
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.listOrgDomains(req);
    }

    public setPrimaryOrgDomain(domain: string): Promise<SetPrimaryOrgDomainResponse> {
        const req = new SetPrimaryOrgDomainRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.setPrimaryOrgDomain(req);
    }

    public generateOrgDomainValidation(domain: string, type: DomainValidationType):
        Promise<GenerateOrgDomainValidationResponse> {
        const req: GenerateOrgDomainValidationRequest = new GenerateOrgDomainValidationRequest();
        req.setDomain(domain);
        req.setType(type);

        return this.grpcService.mgmt.generateOrgDomainValidation(req);
    }

    public validateOrgDomain(domain: string):
        Promise<ValidateOrgDomainResponse> {
        const req = new ValidateOrgDomainRequest();
        req.setDomain(domain);

        return this.grpcService.mgmt.validateOrgDomain(req);
    }

    public listOrgMembers(limit: number, offset: number): Promise<ListOrgMembersResponse> {
        const req = new ListOrgMembersRequest();
        const query = new ListQuery();
        if (limit) {
            query.setLimit(limit);
        }
        if (offset) {
            query.setOffset(offset);
        }
        req.setMetaData(query);

        return this.grpcService.mgmt.listOrgMembers(req);
    }

    public getOrgByDomainGlobal(domain: string): Promise<GetOrgByDomainGlobalResponse> {
        const req = new GetOrgByDomainGlobalRequest();
        req.setDomain(domain);
        return this.grpcService.mgmt.getOrgByDomainGlobal(req);
    }

    public addOrg(name: string): Promise<AddOrgResponse> {
        const req = new AddOrgRequest();
        req.setName(name);
        return this.grpcService.mgmt.addOrg(req);
    }

    public addOrgMember(userId: string, rolesList: string[]): Promise<Empty> {
        const req = new AddOrgMemberRequest();
        req.setUserId(userId);
        if (rolesList) {
            req.setRolesList(rolesList);
        }
        return this.grpcService.mgmt.addOrgMember(req);
    }

    public updateOrgMember(userId: string, rolesList: string[]): Promise<UpdateOrgMemberResponse> {
        const req = new UpdateOrgMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.updateOrgMember(req);
    }


    public removeOrgMember(userId: string): Promise<RemoveOrgMemberResponse> {
        const req = new RemoveOrgMemberRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.removeOrgMember(req);
    }

    public deactivateOrg(): Promise<DeactivateOrgResponse> {
        const req = new DeactivateOrgRequest();
        return this.grpcService.mgmt.deactivateOrg(req);
    }

    public reactivateOrg(): Promise<ReactivateOrgResponse> {
        const req = new ReactivateOrgRequest();
        return this.grpcService.mgmt.reactivateOrg(req);
    }

    public addProjectGrant(
        orgId: string,
        projectId: string,
        roleKeysList: string[],
    ): Promise<AddProjectGrantResponse> {
        const req = new AddProjectGrantRequest();
        req.setProjectId(projectId);
        req.setGrantedOrgId(orgId);
        req.setRoleKeysList(roleKeysList);
        return this.grpcService.mgmt.addProjectGrant(req);
    }

    public listOrgMemberRoles(): Promise<ListOrgMemberRolesResponse> {
        const req = new ListOrgMemberRolesRequest();
        return this.grpcService.mgmt.listOrgMemberRoles(req);
    }

    // Policy

    public getOrgIAMPolicy(): Promise<GetOrgIAMPolicyResponse> {
        const req = new GetOrgIAMPolicyRequest();
        return this.grpcService.mgmt.getOrgIAMPolicy(req);
    }

    public GetPasswordAgePolicy(): Promise<GetPasswordAgePolicyResponse> {
        const req = new GetPasswordAgePolicyRequest();
        return this.grpcService.mgmt.getPasswordAgePolicy(req);
    }

    public addCustomPasswordAgePolicy(
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<AddCustomPasswordAgePolicyResponse> {
        const req = new AddCustomPasswordAgePolicyRequest();
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);

        return this.grpcService.mgmt.addCustomPasswordAgePolicy(req);
    }

    public resetPasswordAgePolicyToDefault(): Promise<ResetPasswordAgePolicyToDefaultResponse> {
        const req = new ResetPasswordAgePolicyToDefaultRequest();
        return this.grpcService.mgmt.resetPasswordAgePolicyToDefault(req);
    }

    public updateCustomPasswordAgePolicy(
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<UpdateCustomPasswordAgePolicyResponse> {
        const req = new UpdateCustomPasswordAgePolicyRequest();
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);
        return this.grpcService.mgmt.updateCustomPasswordAgePolicy(req);
    }

    public GetPasswordComplexityPolicy(): Promise<GetPasswordComplexityPolicyResponse> {
        const req = new GetPasswordComplexityPolicyRequest();
        return this.grpcService.mgmt.getPasswordComplexityPolicy(req);
    }

    public addCustomPasswordComplexityPolicy(
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<AddCustomPasswordComplexityPolicyResponse> {
        const req = new AddCustomPasswordComplexityPolicyRequest();
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return this.grpcService.mgmt.addCustomPasswordComplexityPolicy(req);
    }

    public resetPasswordComplexityPolicyToDefault(): Promise<ResetPasswordComplexityPolicyToDefaultResponse> {
        const req = new ResetPasswordComplexityPolicyToDefaultRequest();
        return this.grpcService.mgmt.resetPasswordComplexityPolicyToDefault(req);
    }

    public updateCustomPasswordComplexityPolicy(
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<UpdateCustomPasswordComplexityPolicyResponse> {
        const req = new UpdateCustomPasswordComplexityPolicyRequest();
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return this.grpcService.mgmt.updateCustomPasswordComplexityPolicy(req);
    }

    public getPasswordLockoutPolicy(): Promise<GetPasswordLockoutPolicyResponse> {
        const req = new GetPasswordLockoutPolicyRequest();

        return this.grpcService.mgmt.getPasswordLockoutPolicy(req);
    }

    public addCustomPasswordLockoutPolicy(
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<AddCustomPasswordLockoutPolicyResponse> {
        const req = new AddCustomPasswordLockoutPolicyRequest();
        req.setMaxAttempts(maxAttempts);
        req.setShowLockoutFailure(showLockoutFailures);

        return this.grpcService.mgmt.addCustomPasswordLockoutPolicy(req);
    }

    public resetPasswordLockoutPolicyToDefault(): Promise<ResetPasswordLockoutPolicyToDefaultResponse> {
        const req = new ResetPasswordLockoutPolicyToDefaultRequest();
        return this.grpcService.mgmt.resetPasswordLockoutPolicyToDefault(req);
    }

    public updateCustomPasswordLockoutPolicy(
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<UpdateCustomPasswordLockoutPolicyResponse> {
        const req = new UpdateCustomPasswordLockoutPolicyRequest();
        req.setMaxAttempts(maxAttempts);
        req.setShowLockoutFailure(showLockoutFailures);
        return this.grpcService.mgmt.updateCustomPasswordLockoutPolicy(req);
    }

    public getLocalizedComplexityPolicyPatternErrorString(policy: PasswordComplexityPolicy.AsObject): string {
        if (policy.hasNumber && policy.hasSymbol) {
            return 'POLICY.PWD_COMPLEXITY.SYMBOLANDNUMBERERROR';
        } else if (policy.hasNumber) {
            return 'POLICY.PWD_COMPLEXITY.NUMBERERROR';
        } else if (policy.hasSymbol) {
            return 'POLICY.PWD_COMPLEXITY.SYMBOLERROR';
        } else {
            return 'POLICY.PWD_COMPLEXITY.PATTERNERROR';
        }
    }

    public getUserByID(id: string): Promise<GetUserByIDResponse> {
        const req = new GetUserByIDRequest();
        req.setId(id);
        return this.grpcService.mgmt.getUserByID(req);
    }

    public removeUser(id: string): Promise<RemoveUserResponse> {
        const req = new RemoveUserRequest();
        req.setId(id);
        return this.grpcService.mgmt.removeUser(req);
    }

    public listProjectMembers(
        projectId: string,
        limit: number,
        offset: number,
        queryList?: SearchQuery[],
    ): Promise<ListProjectMembersResponse> {
        const req = new ListProjectMembersRequest();
        const query = new ListQuery();
        req.setMetaData(query);
        req.setProjectId(projectId);
        if (limit) {
            query.setLimit(limit);
        }
        if (offset) {
            query.setOffset(offset);
        }
        if (queryList) {
            req.setQueriesList(queryList);
        }
        req.setMetaData(query);
        return this.grpcService.mgmt.listProjectMembers(req);
    }

    public listUserMemberships(userId: string,
        limit: number, offset: number,
        queryList?: MembershipQuery[],
    ): Promise<ListUserMembershipsResponse> {
        const req = new ListUserMembershipsRequest();
        req.setUserId(userId);
        const metadata = new ListQuery();
        if (limit) {
            metadata.setLimit(limit);
        }
        if (offset) {
            metadata.setOffset(offset);
        }
        if (queryList) {
            req.setQueriesList(queryList);
        }
        req.setMetaData(metadata);
        return this.grpcService.mgmt.listUserMemberships(req);
    }

    public GetUserProfile(userId: string): Promise<GetHumanProfileResponse> {
        const req = new GetHumanProfileRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.getHumanProfile(req);
    }

    public listUserMultiFactors(userId: string): Promise<ListUserMultiFactorsResponse> {
        const req = new ListUserMultiFactorsRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.listUserMultiFactors(req);
    }

    public removeHumanMultiFactorOTP(userId: string): Promise<RemoveHumanMultiFactorOTPResponse> {
        const req = new RemoveHumanMultiFactorOTPRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.removeHumanMultiFactorOTP(req);
    }

    public removeHumanMultiFactorU2F(userId: string, id: string): Promise<RemoveHumanMultiFactorU2FResponse> {
        const req = new RemoveHumanMultiFactorU2FRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.removeHumanMultiFactorU2F(req);
    }

    public updateHumanProfile(
        userId: string,
        firstName?: string,
        lastName?: string,
        nickName?: string,
        preferredLanguage?: string,
        gender?: Gender,
    ): Promise<UpdateHumanProfileResponse> {
        const req = new UpdateHumanProfileRequest();
        req.setUserId(userId);
        if (firstName) {
            req.setFirstName(firstName);
        }
        if (lastName) {
            req.setLastName(lastName);
        }
        if (nickName) {
            req.setNickName(nickName);
        }
        if (gender) {
            req.setGender(gender);
        }
        if (preferredLanguage) {
            req.setPreferredLanguage(preferredLanguage);
        }
        return this.grpcService.mgmt.updateHumanProfile(req);
    }

    public getHumanEmail(id: string): Promise<GetHumanEmailResponse> {
        const req = new GetHumanEmailRequest();
        req.setUserId(id);
        return this.grpcService.mgmt.getHumanEmail(req);
    }

    public updateHumanEmail(userId: string, email: string): Promise<UpdateHumanEmailResponse> {
        const req = new UpdateHumanEmailRequest();
        req.setUserId(userId);
        req.setEmail(email);
        return this.grpcService.mgmt.updateHumanEmail(req);
    }

    public getHumanPhone(userId: string): Promise<GetHumanPhoneResponse> {
        const req = new GetHumanPhoneRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.getHumanPhone(req);
    }

    public updateHumanPhone(userId: string, phone: string): Promise<UpdateHumanPhoneResponse> {
        const req = new UpdateHumanPhoneRequest();
        req.setUserId(userId);
        req.setPhone(phone);
        return this.grpcService.mgmt.updateHumanPhone(req);
    }

    public removeHumanPhone(userId: string): Promise<Empty> {
        const req = new RemoveHumanPhoneRequest();
        req.setUserId(userId);
        return this.grpcService.mgmt.removeHumanPhone(req);
    }

    public deactivateUser(id: string): Promise<DeactivateUserResponse> {
        const req = new DeactivateUserRequest();
        req.setId(id);
        return this.grpcService.mgmt.deactivateUser(req);
    }

    public addUserGrant(
        userId: string,
        roleNamesList: string[],
        projectId?: string,
    ): Promise<AddUserGrantResponse> {
        const req = new AddUserGrantRequest();
        if (projectId) {
            req.setProjectId(projectId);
        }
        req.setUserId(userId);
        req.setRoleKeysList(roleNamesList);

        return this.grpcService.mgmt.addUserGrant(req);
    }

    public reactivateUser(id: string): Promise<ReactivateUserResponse> {
        const req = new ReactivateUserRequest();
        req.setId(id);
        return this.grpcService.mgmt.reactivateUser(req);
    }

    public AddRole(projectId: string, roleKey: string, displayName: string, group: string): Promise<AddProjectRoleResponse> {
        const req = new AddProjectRoleRequest();
        req.setProjectId(projectId);
        req.setRoleKey(roleKey);
        if (displayName) {
            req.setDisplayName(displayName);
        }
        req.setGroup(group);
        return this.grpcService.mgmt.addProjectRole(req);
    }

    public GetUserAddress(id: string): Promise<UserAddress> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.getUserAddress(req);
    }

    public ResendEmailVerification(id: string): Promise<any> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.resendEmailVerificationMail(req);
    }

    public ResendInitialMail(userId: string, newemail: string): Promise<Empty> {
        const req = new InitialMailRequest();
        if (newemail) {
            req.setEmail(newemail);
        }
        req.setId(userId);

        return this.grpcService.mgmt.resendInitialMail(req);
    }

    public ResendPhoneVerification(id: string): Promise<any> {
        const req = new UserID();
        req.setId(id);
        return this.grpcService.mgmt.resendPhoneVerificationCode(req);
    }

    public SetInitialPassword(id: string, password: string): Promise<any> {
        const req = new PasswordRequest();
        req.setId(id);
        req.setPassword(password);
        return this.grpcService.mgmt.setInitialPassword(req);
    }

    public SendSetPasswordNotification(id: string, type: NotificationType): Promise<any> {
        const req = new SetPasswordNotificationRequest();
        req.setId(id);
        req.setType(type);
        return this.grpcService.mgmt.sendSetPasswordNotification(req);
    }

    public SaveUserAddress(address: UserAddress.AsObject): Promise<UserAddress> {
        const req = new UpdateUserAddressRequest();
        req.setId(address.id);
        req.setStreetAddress(address.streetAddress);
        req.setPostalCode(address.postalCode);
        req.setLocality(address.locality);
        req.setRegion(address.region);
        req.setCountry(address.country);
        return this.grpcService.mgmt.updateUserAddress(req);
    }

    public SearchUsers(limit: number, offset: number, queryList?: UserSearchQuery[]): Promise<UserSearchResponse> {
        const req = new UserSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchUsers(req);
    }

    public GetUserByLoginNameGlobal(loginname: string): Promise<UserView> {
        const req = new LoginName();
        req.setLoginName(loginname);
        return this.grpcService.mgmt.getUserByLoginNameGlobal(req);
    }

    // USER GRANTS

    public SearchUserGrants(
        limit?: number,
        offset?: number,
        queryList?: UserGrantSearchQuery[],
    ): Promise<UserGrantSearchResponse> {
        const req = new UserGrantSearchRequest();
        if (limit) {
            req.setLimit(limit);
        }
        if (offset) {
            req.setOffset(offset);
        }
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchUserGrants(req);
    }


    public UserGrantByID(
        id: string,
        userId: string,
    ): Promise<UserGrantView> {
        const req = new UserGrantID();
        req.setId(id);
        req.setUserId(userId);

        return this.grpcService.mgmt.userGrantByID(req);
    }

    public UpdateUserGrant(
        id: string,
        userId: string,
        roleKeysList: string[],
    ): Promise<UserGrant> {
        const req = new UserGrantUpdate();
        req.setId(id);
        req.setRoleKeysList(roleKeysList);
        req.setUserId(userId);

        return this.grpcService.mgmt.updateUserGrant(req);
    }

    public RemoveUserGrant(
        id: string,
        userId: string,
    ): Promise<Empty> {
        const req = new UserGrantID();
        req.setId(id);
        req.setUserId(userId);

        return this.grpcService.mgmt.removeUserGrant(req);
    }

    public BulkRemoveUserGrant(
        idsList: string[],
    ): Promise<Empty> {
        const req = new UserGrantRemoveBulk();
        req.setIdsList(idsList);

        return this.grpcService.mgmt.bulkRemoveUserGrant(req);
    }

    //

    public ApplicationChanges(id: string, secId: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setSecId(secId);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.applicationChanges(req);
    }

    public OrgChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.orgChanges(req);
    }

    public ProjectChanges(id: string, limit: number, offset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(offset);
        return this.grpcService.mgmt.projectChanges(req);
    }

    public UserChanges(id: string, limit: number, sequenceoffset: number): Promise<Changes> {
        const req = new ChangeRequest();
        req.setId(id);
        req.setLimit(limit);
        req.setSequenceOffset(sequenceoffset);
        return this.grpcService.mgmt.userChanges(req);
    }

    // project

    public SearchProjects(
        limit?: number, offset?: number, queryList?: ProjectSearchQuery[]): Promise<ProjectSearchResponse> {
        const req = new ProjectSearchRequest();
        if (limit) {
            req.setLimit(limit);
        }
        if (offset) {
            req.setOffset(offset);
        }

        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchProjects(req).then(value => {
            const count = value.toObject().resultList.length;
            if (count >= 0) {
                this.ownedProjectsCount.next(count);
            }

            return value;
        });
    }

    public SearchGrantedProjects(
        limit: number, offset: number, queryList?: ProjectSearchQuery[]): Promise<ProjectGrantSearchResponse> {
        const req = new GrantedProjectSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchGrantedProjects(req).then(value => {
            this.grantedProjectsCount.next(value.toObject().resultList.length);
            return value;
        });
    }

    public GetZitadelDocs(): Promise<ZitadelDocs> {
        const req = new Empty();
        return this.grpcService.mgmt.getZitadelDocs(req);
    }

    public GetProjectById(projectId: string): Promise<ProjectView> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.projectByID(req);
    }

    public GetGrantedProjectByID(projectId: string, id: string): Promise<ProjectGrantView> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.getGrantedProjectByID(req);
    }

    public CreateProject(project: ProjectCreateRequest.AsObject): Promise<Project> {
        const req = new ProjectCreateRequest();
        req.setName(project.name);
        return this.grpcService.mgmt.createProject(req).then(value => {
            const current = this.ownedProjectsCount.getValue();
            this.ownedProjectsCount.next(current + 1);
            return value;
        });
    }

    public UpdateProject(id: string, projectView: ProjectView.AsObject): Promise<Project> {
        const req = new ProjectUpdateRequest();
        req.setId(id);
        req.setName(projectView.name);
        req.setProjectRoleAssertion(projectView.projectRoleAssertion);
        req.setProjectRoleCheck(projectView.projectRoleCheck);
        return this.grpcService.mgmt.updateProject(req);
    }

    public UpdateProjectGrant(id: string, projectId: string, rolesList: string[]): Promise<ProjectGrant> {
        const req = new ProjectGrantUpdate();
        req.setRoleKeysList(rolesList);
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.updateProjectGrant(req);
    }

    public RemoveProjectGrant(id: string, projectId: string): Promise<Empty> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.removeProjectGrant(req);
    }

    public DeactivateProject(projectId: string): Promise<Project> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.deactivateProject(req);
    }

    public ReactivateProject(projectId: string): Promise<Project> {
        const req = new ProjectID();
        req.setId(projectId);
        return this.grpcService.mgmt.reactivateProject(req);
    }

    public SearchProjectGrants(projectId: string, limit: number, offset: number): Promise<ProjectGrantSearchResponse> {
        const req = new ProjectGrantSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        return this.grpcService.mgmt.searchProjectGrants(req);
    }

    public GetProjectGrantMemberRoles(): Promise<ProjectGrantMemberRoles> {
        const req = new Empty();
        return this.grpcService.mgmt.getProjectGrantMemberRoles(req);
    }

    public AddProjectMember(id: string, userId: string, rolesList: string[]): Promise<Empty> {
        const req = new ProjectMemberAdd();
        req.setId(id);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.addProjectMember(req);
    }

    public ChangeProjectMember(id: string, userId: string, rolesList: string[]): Promise<ProjectMember> {
        const req = new ProjectMemberChange();
        req.setId(id);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.changeProjectMember(req);
    }

    public AddProjectGrantMember(
        projectId: string,
        grantId: string,
        userId: string,
        rolesList: string[],
    ): Promise<Empty> {
        const req = new ProjectGrantMemberAdd();
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.addProjectGrantMember(req);
    }

    public ChangeProjectGrantMember(
        projectId: string,
        grantId: string,
        userId: string,
        rolesList: string[],
    ): Promise<ProjectGrantMember> {
        const req = new ProjectGrantMemberChange();
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return this.grpcService.mgmt.changeProjectGrantMember(req);
    }

    public SearchProjectGrantMembers(
        projectId: string,
        grantId: string,
        limit: number,
        offset: number,
        queryList?: ProjectGrantMemberSearchQuery[],
    ): Promise<ProjectMemberSearchResponse> {
        const req = new ProjectGrantMemberSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        return this.grpcService.mgmt.searchProjectGrantMembers(req);
    }

    public RemoveProjectGrantMember(
        projectId: string,
        grantId: string,
        userId: string,
    ): Promise<Empty> {
        const req = new ProjectGrantMemberRemove();
        req.setGrantId(grantId);
        req.setUserId(userId);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.removeProjectGrantMember(req);
    }

    public ReactivateApplication(projectId: string, appId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setId(appId);
        req.setProjectId(projectId);

        return this.grpcService.mgmt.reactivateApplication(req);
    }

    public DeactivateApplication(projectId: string, appId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setId(appId);
        req.setProjectId(projectId);

        return this.grpcService.mgmt.deactivateApplication(req);
    }

    public RegenerateOIDCClientSecret(id: string, projectId: string): Promise<any> {
        const req = new ApplicationID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.regenerateOIDCClientSecret(req);
    }

    public SearchProjectRoles(
        projectId: string,
        limit: number,
        offset: number,
        queryList?: ProjectRoleSearchQuery[],
    ): Promise<ProjectRoleSearchResponse> {
        const req = new ProjectRoleSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchProjectRoles(req);
    }

    public AddProjectRole(role: ProjectRoleAdd.AsObject): Promise<Empty> {
        const req = new ProjectRoleAdd();
        req.setId(role.id);
        if (role.displayName) {
            req.setDisplayName(role.displayName);
        }
        req.setKey(role.key);
        req.setGroup(role.group);
        return this.grpcService.mgmt.addProjectRole(req);
    }

    public BulkAddProjectRole(
        id: string,
        rolesList: ProjectRoleAdd[],
    ): Promise<Empty> {
        const req = new ProjectRoleAddBulk();
        req.setId(id);
        req.setProjectRolesList(rolesList);
        return this.grpcService.mgmt.bulkAddProjectRole(req);
    }

    public RemoveProjectRole(projectId: string, key: string): Promise<Empty> {
        const req = new ProjectRoleRemove();
        req.setId(projectId);
        req.setKey(key);
        return this.grpcService.mgmt.removeProjectRole(req);
    }


    public ChangeProjectRole(projectId: string, key: string, displayName: string, group: string):
        Promise<ProjectRole> {
        const req = new ProjectRoleChange();
        req.setId(projectId);
        req.setKey(key);
        req.setGroup(group);
        req.setDisplayName(displayName);
        return this.grpcService.mgmt.changeProjectRole(req);
    }


    public RemoveProjectMember(id: string, userId: string): Promise<Empty> {
        const req = new ProjectMemberRemove();
        req.setId(id);
        req.setUserId(userId);
        return this.grpcService.mgmt.removeProjectMember(req);
    }

    public SearchApplications(
        projectId: string,
        limit: number,
        offset: number,
        queryList?: ApplicationSearchQuery[]): Promise<ApplicationSearchResponse> {
        const req = new ApplicationSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.mgmt.searchApplications(req);
    }

    public GetApplicationById(projectId: string, applicationId: string): Promise<ApplicationView> {
        const req = new ApplicationID();
        req.setProjectId(projectId);
        req.setId(applicationId);
        return this.grpcService.mgmt.applicationByID(req);
    }

    public GetProjectMemberRoles(): Promise<ProjectMemberRoles> {
        const req = new Empty();
        return this.grpcService.mgmt.getProjectMemberRoles(req);
    }

    public ProjectGrantByID(id: string, projectId: string): Promise<ProjectGrantView> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.projectGrantByID(req);
    }

    public RemoveProject(id: string): Promise<Empty> {
        const req = new ProjectID();
        req.setId(id);
        return this.grpcService.mgmt.removeProject(req).then(value => {
            const current = this.ownedProjectsCount.getValue();
            this.ownedProjectsCount.next(current > 0 ? current - 1 : 0);
            return value;
        });
    }


    public DeactivateProjectGrant(id: string, projectId: string): Promise<ProjectGrant> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.deactivateProjectGrant(req);
    }

    public ReactivateProjectGrant(id: string, projectId: string): Promise<ProjectGrant> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.reactivateProjectGrant(req);
    }

    public CreateOIDCApp(app: OIDCApplicationCreate.AsObject): Promise<Application> {
        const req = new OIDCApplicationCreate();
        req.setProjectId(app.projectId);
        req.setName(app.name);
        req.setRedirectUrisList(app.redirectUrisList);
        req.setResponseTypesList(app.responseTypesList);
        req.setGrantTypesList(app.grantTypesList);
        req.setApplicationType(app.applicationType);
        req.setAuthMethodType(app.authMethodType);
        req.setPostLogoutRedirectUrisList(app.postLogoutRedirectUrisList);

        return this.grpcService.mgmt.createOIDCApplication(req);
    }

    public UpdateApplication(projectId: string, appId: string, name: string): Promise<Application> {
        const req = new ApplicationUpdate();
        req.setId(appId);
        req.setName(name);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.updateApplication(req);
    }

    public UpdateOIDCAppConfig(req: OIDCConfigUpdate): Promise<OIDCConfig> {
        return this.grpcService.mgmt.updateApplicationOIDCConfig(req);
    }

    public RemoveApplication(projectId: string, appId: string): Promise<Empty> {
        const req = new ApplicationID();
        req.setId(appId);
        req.setProjectId(projectId);
        return this.grpcService.mgmt.removeApplication(req);
    }
}

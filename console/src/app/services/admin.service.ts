import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

import {
    AddIamMemberRequest,
    ChangeIamMemberRequest,
    CreateHumanRequest,
    CreateOrgRequest,
    CreateUserRequest,
    DefaultLabelPolicy,
    DefaultLabelPolicyUpdate,
    DefaultLabelPolicyView,
    DefaultLoginPolicy,
    DefaultLoginPolicyRequest,
    DefaultLoginPolicyView,
    DefaultPasswordAgePolicyRequest,
    DefaultPasswordAgePolicyView,
    DefaultPasswordComplexityPolicy,
    DefaultPasswordComplexityPolicyRequest,
    DefaultPasswordComplexityPolicyView,
    DefaultPasswordLockoutPolicy,
    DefaultPasswordLockoutPolicyRequest,
    DefaultPasswordLockoutPolicyView,
    FailedEventID,
    FailedEvents,
    IamMember,
    IamMemberRoles,
    IamMemberSearchQuery,
    IamMemberSearchRequest,
    IamMemberSearchResponse,
    Idp,
    IdpID,
    IdpProviderID,
    IdpProviderSearchRequest,
    IdpProviderSearchResponse,
    IdpSearchQuery,
    IdpSearchRequest,
    IdpSearchResponse,
    IdpView,
    MultiFactor,
    MultiFactorsResult,
    OidcIdpConfig,
    OidcIdpConfigCreate,
    OidcIdpConfigUpdate,
    OrgIamPolicy,
    OrgIamPolicyID,
    OrgIamPolicyRequest,
    OrgIamPolicyView,
    OrgSetUpRequest,
    OrgSetUpResponse,
    RemoveIamMemberRequest,
    SecondFactor,
    SecondFactorsResult,
    ViewID,
    Views,
} from '../proto/generated/admin_pb';
import { IdpUpdate } from '../proto/generated/management_pb';
import { GrpcService } from './grpc.service';

@Injectable({
    providedIn: 'root',
})
export class AdminService {
    constructor(private readonly grpcService: GrpcService) { }

    public SetUpOrg(
        createOrgRequest: CreateOrgRequest,
        humanRequest: CreateHumanRequest,
    ): Promise<OrgSetUpResponse> {
        const req: OrgSetUpRequest = new OrgSetUpRequest();
        const userReq: CreateUserRequest = new CreateUserRequest();

        userReq.setHuman(humanRequest);

        req.setOrg(createOrgRequest);
        req.setUser(userReq);

        return this.grpcService.admin.setUpOrg(req);
    }

    public getDefaultLoginPolicyMultiFactors(): Promise<MultiFactorsResult> {
        const req = new Empty();
        return this.grpcService.admin.getDefaultLoginPolicyMultiFactors(req);
    }

    public addMultiFactorToDefaultLoginPolicy(req: MultiFactor): Promise<MultiFactor> {
        return this.grpcService.admin.addMultiFactorToDefaultLoginPolicy(req);
    }

    public RemoveMultiFactorFromDefaultLoginPolicy(req: MultiFactor): Promise<Empty> {
        return this.grpcService.admin.removeMultiFactorFromDefaultLoginPolicy(req);
    }

    public GetDefaultLoginPolicySecondFactors(): Promise<SecondFactorsResult> {
        const req = new Empty();
        return this.grpcService.admin.getDefaultLoginPolicySecondFactors(req);
    }

    public AddSecondFactorToDefaultLoginPolicy(req: SecondFactor): Promise<SecondFactor> {
        return this.grpcService.admin.addSecondFactorToDefaultLoginPolicy(req);
    }

    public RemoveSecondFactorFromDefaultLoginPolicy(req: SecondFactor): Promise<Empty> {
        return this.grpcService.admin.removeSecondFactorFromDefaultLoginPolicy(req);
    }

    public GetIamMemberRoles(): Promise<IamMemberRoles> {
        const req = new Empty();
        return this.grpcService.admin.getIamMemberRoles(req);
    }

    public GetViews(): Promise<Views> {
        const req = new Empty();
        return this.grpcService.admin.getViews(req);
    }

    public GetFailedEvents(): Promise<FailedEvents> {
        const req = new Empty();
        return this.grpcService.admin.getFailedEvents(req);
    }

    public ClearView(viewname: string, db: string): Promise<Empty> {
        const req: ViewID = new ViewID();
        req.setDatabase(db);
        req.setViewName(viewname);
        return this.grpcService.admin.clearView(req);
    }

    public RemoveFailedEvent(viewname: string, db: string, sequence: number): Promise<Empty> {
        const req: FailedEventID = new FailedEventID();
        req.setDatabase(db);
        req.setViewName(viewname);
        req.setFailedSequence(sequence);
        return this.grpcService.admin.removeFailedEvent(req);
    }

    /* Policies */

    /* complexity */

    public GetDefaultPasswordComplexityPolicy(): Promise<DefaultPasswordComplexityPolicyView> {
        const req = new Empty();
        return this.grpcService.admin.getDefaultPasswordComplexityPolicy(req);
    }

    public UpdateDefaultPasswordComplexityPolicy(
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<DefaultPasswordComplexityPolicy> {
        const req = new DefaultPasswordComplexityPolicyRequest();
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return this.grpcService.admin.updateDefaultPasswordComplexityPolicy(req);
    }

    /* age */

    public GetDefaultPasswordAgePolicy(): Promise<DefaultPasswordAgePolicyView> {
        const req = new Empty();

        return this.grpcService.admin.getDefaultPasswordAgePolicy(req);
    }

    public UpdateDefaultPasswordAgePolicy(
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<DefaultPasswordAgePolicyView> {
        const req = new DefaultPasswordAgePolicyRequest();
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);

        return this.grpcService.admin.updateDefaultPasswordAgePolicy(req);
    }

    /* lockout */

    public GetDefaultPasswordLockoutPolicy(): Promise<DefaultPasswordLockoutPolicyView> {
        const req = new Empty();
        return this.grpcService.admin.getDefaultPasswordLockoutPolicy(req);
    }

    public UpdateDefaultPasswordLockoutPolicy(
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<DefaultPasswordLockoutPolicy> {
        const req = new DefaultPasswordLockoutPolicyRequest();
        req.setMaxAttempts(maxAttempts);
        req.setShowLockoutFailure(showLockoutFailures);

        return this.grpcService.admin.updateDefaultPasswordLockoutPolicy(req);
    }

    /* label */

    public GetDefaultLabelPolicy(): Promise<DefaultLabelPolicyView> {
        const req = new Empty();
        return this.grpcService.admin.getDefaultLabelPolicy(req);
    }

    public UpdateDefaultLabelPolicy(req: DefaultLabelPolicyUpdate): Promise<DefaultLabelPolicy> {
        return this.grpcService.admin.updateDefaultLabelPolicy(req);
    }

    /* login */

    public GetDefaultLoginPolicy(
    ): Promise<DefaultLoginPolicyView> {
        const req = new Empty();
        return this.grpcService.admin.getDefaultLoginPolicy(req);
    }

    public UpdateDefaultLoginPolicy(req: DefaultLoginPolicyRequest): Promise<DefaultLoginPolicy> {
        return this.grpcService.admin.updateDefaultLoginPolicy(req);
    }

    /* org iam */

    public GetOrgIamPolicy(orgId: string): Promise<OrgIamPolicyView> {
        const req = new OrgIamPolicyID();
        req.setOrgId(orgId);
        return this.grpcService.admin.getOrgIamPolicy(req);
    }

    public CreateOrgIamPolicy(
        orgId: string,
        userLoginMustBeDomain: boolean): Promise<OrgIamPolicy> {
        const req = new OrgIamPolicyRequest();
        req.setOrgId(orgId);
        req.setUserLoginMustBeDomain(userLoginMustBeDomain);

        return this.grpcService.admin.createOrgIamPolicy(req);
    }

    public UpdateOrgIamPolicy(
        orgId: string,
        userLoginMustBeDomain: boolean): Promise<OrgIamPolicy> {
        const req = new OrgIamPolicyRequest();
        req.setOrgId(orgId);
        req.setUserLoginMustBeDomain(userLoginMustBeDomain);
        return this.grpcService.admin.updateOrgIamPolicy(req);
    }

    public RemoveOrgIamPolicy(
        orgId: string,
    ): Promise<Empty> {
        const req = new OrgIamPolicyID();
        req.setOrgId(orgId);
        return this.grpcService.admin.removeOrgIamPolicy(req);
    }

    /* admin iam */

    public GetDefaultOrgIamPolicy(): Promise<OrgIamPolicyView> {
        const req = new Empty();
        return this.grpcService.admin.getDefaultOrgIamPolicy(req);
    }

    /* policies end */

    public AddIdpProviderToDefaultLoginPolicy(configId: string): Promise<IdpProviderID> {
        const req = new IdpProviderID();
        req.setIdpConfigId(configId);
        return this.grpcService.admin.addIdpProviderToDefaultLoginPolicy(req);
    }

    public RemoveIdpProviderFromDefaultLoginPolicy(configId: string): Promise<Empty> {
        const req = new IdpProviderID();
        req.setIdpConfigId(configId);
        return this.grpcService.admin.removeIdpProviderFromDefaultLoginPolicy(req);
    }

    public GetDefaultLoginPolicyIdpProviders(limit?: number, offset?: number): Promise<IdpProviderSearchResponse> {
        const req = new IdpProviderSearchRequest();
        if (limit) {
            req.setLimit(limit);
        }
        if (offset) {
            req.setOffset(offset);
        }
        return this.grpcService.admin.getDefaultLoginPolicyIdpProviders(req);
    }

    public SearchIdps(
        limit?: number,
        offset?: number,
        queryList?: IdpSearchQuery[],
    ): Promise<IdpSearchResponse> {
        const req = new IdpSearchRequest();
        if (limit) {
            req.setLimit(limit);
        }
        if (offset) {
            req.setOffset(offset);
        }
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return this.grpcService.admin.searchIdps(req);
    }

    public IdpByID(
        id: string,
    ): Promise<IdpView> {
        const req = new IdpID();
        req.setId(id);
        return this.grpcService.admin.idpByID(req);
    }

    public UpdateIdp(
        req: IdpUpdate,
    ): Promise<Idp> {
        return this.grpcService.admin.updateIdpConfig(req);
    }

    public CreateOidcIdp(
        req: OidcIdpConfigCreate,
    ): Promise<Idp> {
        return this.grpcService.admin.createOidcIdp(req);
    }

    public UpdateOidcIdpConfig(
        req: OidcIdpConfigUpdate,
    ): Promise<OidcIdpConfig> {
        return this.grpcService.admin.updateOidcIdpConfig(req);
    }

    public RemoveIdpConfig(
        id: string,
    ): Promise<Empty> {
        const req = new IdpID;
        req.setId(id);
        return this.grpcService.admin.removeIdpConfig(req);
    }

    public DeactivateIdpConfig(
        id: string,
    ): Promise<Empty> {
        const req = new IdpID;
        req.setId(id);
        return this.grpcService.admin.deactivateIdpConfig(req);
    }

    public ReactivateIdpConfig(
        id: string,
    ): Promise<Empty> {
        const req = new IdpID;
        req.setId(id);
        return this.grpcService.admin.reactivateIdpConfig(req);
    }

    public SearchIamMembers(
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

    public RemoveIamMember(
        userId: string,
    ): Promise<Empty> {
        const req = new RemoveIamMemberRequest();
        req.setUserId(userId);

        return this.grpcService.admin.removeIamMember(req);
    }

    public AddIamMember(
        userId: string,
        rolesList: string[],
    ): Promise<IamMember> {
        const req = new AddIamMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);

        return this.grpcService.admin.addIamMember(req);
    }

    public ChangeIamMember(
        userId: string,
        rolesList: string[],
    ): Promise<IamMember> {
        const req = new ChangeIamMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);

        return this.grpcService.admin.changeIamMember(req);
    }
}

import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

import {
    AddCustomOrgIAMPolicyRequest,
    AddCustomOrgIAMPolicyResponse,
    AddIAMMemberRequest,
    AddIAMMemberResponse,
    AddIDPToDefaultLoginPolicyRequest,
    AddIDPToDefaultLoginPolicyResponse,
    AddMultiFactorToDefaultLoginPolicyRequest,
    AddMultiFactorToDefaultLoginPolicyResponse,
    AddOIDCIDPRequest,
    AddOIDCIDPResponse,
    AddSecondFactorToDefaultLoginPolicyRequest,
    AddSecondFactorToDefaultLoginPolicyResponse,
    ClearViewRequest,
    ClearViewResponse,
    DeactivateIDPRequest,
    DeactivateIDPResponse,
    GetDefaultLabelPolicyRequest,
    GetDefaultLabelPolicyResponse,
    GetDefaultLoginPolicyRequest,
    GetDefaultLoginPolicyResponse,
    GetDefaultOrgIAMPolicyRequest,
    GetDefaultOrgIAMPolicyResponse,
    GetDefaultPasswordAgePolicyResponse,
    GetDefaultPasswordComplexityPolicyRequest,
    GetDefaultPasswordComplexityPolicyResponse,
    GetDefaultPasswordLockoutPolicyRequest,
    GetDefaultPasswordLockoutPolicyResponse,
    GetIDPByIDRequest,
    GetIDPByIDResponse,
    GetOrgIAMPolicyRequest,
    GetOrgIAMPolicyResponse,
    ListDefaultLoginPolicyIDPsRequest,
    ListDefaultLoginPolicyIDPsResponse,
    ListDefaultLoginPolicyMultiFactorsRequest,
    ListDefaultLoginPolicyMultiFactorsResponse,
    ListDefaultLoginPolicySecondFactorsResponse,
    ListFailedEventsRequest,
    ListFailedEventsResponse,
    ListIAMMemberRolesRequest,
    ListIAMMemberRolesResponse,
    ListIAMMembersRequest,
    ListIAMMembersResponse,
    ListIDPsRequest,
    ListIDPsResponse,
    ListViewsRequest,
    ListViewsResponse,
    ReactivateIDPRequest,
    ReactivateIDPResponse,
    RemoveFailedEventRequest,
    RemoveFailedEventResponse,
    RemoveIAMMemberRequest,
    RemoveIAMMemberResponse,
    RemoveIDPFromDefaultLoginPolicyRequest,
    RemoveIDPFromDefaultLoginPolicyResponse,
    RemoveIDPRequest,
    RemoveIDPResponse,
    RemoveMultiFactorFromDefaultLoginPolicyRequest,
    RemoveMultiFactorFromDefaultLoginPolicyResponse,
    RemoveSecondFactorFromDefaultLoginPolicyRequest,
    RemoveSecondFactorFromDefaultLoginPolicyResponse,
    ResetOrgIAMPolicyToDefaultRequest,
    ResetOrgIAMPolicyToDefaultResponse,
    SetUpOrgRequest,
    SetUpOrgResponse,
    UpdateCustomOrgIAMPolicyRequest,
    UpdateCustomOrgIAMPolicyResponse,
    UpdateDefaultLabelPolicyRequest,
    UpdateDefaultLabelPolicyResponse,
    UpdateDefaultLoginPolicyRequest,
    UpdateDefaultLoginPolicyResponse,
    UpdateDefaultPasswordAgePolicyRequest,
    UpdateDefaultPasswordAgePolicyResponse,
    UpdateDefaultPasswordComplexityPolicyRequest,
    UpdateDefaultPasswordComplexityPolicyResponse,
    UpdateDefaultPasswordLockoutPolicyRequest,
    UpdateDefaultPasswordLockoutPolicyResponse,
    UpdateIAMMemberRequest,
    UpdateIAMMemberResponse,
    UpdateIDPOIDCConfigRequest,
    UpdateIDPOIDCConfigResponse,
    UpdateIDPRequest,
    UpdateIDPResponse,
} from '../proto/generated/zitadel/admin_pb';
import { IDPQuery } from '../proto/generated/zitadel/idp_pb';
import { SearchQuery } from '../proto/generated/zitadel/member_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { GrpcService } from './grpc.service';

@Injectable({
    providedIn: 'root',
})
export class AdminService {
    constructor(private readonly grpcService: GrpcService) { }

    public SetUpOrg(
        org: SetUpOrgRequest.Org,
        human: SetUpOrgRequest.Human,
    ): Promise<SetUpOrgResponse> {
        const req = new SetUpOrgRequest();

        req.setOrg(org);
        req.setHuman(human);

        return this.grpcService.admin.setUpOrg(req);
    }

    public listDefaultLoginPolicyMultiFactors(): Promise<ListDefaultLoginPolicyMultiFactorsResponse> {
        const req = new ListDefaultLoginPolicyMultiFactorsRequest();
        return this.grpcService.admin.listDefaultLoginPolicyMultiFactors(req);
    }

    public addMultiFactorToDefaultLoginPolicy(req: AddMultiFactorToDefaultLoginPolicyRequest): Promise<AddMultiFactorToDefaultLoginPolicyResponse> {
        return this.grpcService.admin.addMultiFactorToDefaultLoginPolicy(req);
    }

    public removeMultiFactorFromDefaultLoginPolicy(req: RemoveMultiFactorFromDefaultLoginPolicyRequest): Promise<RemoveMultiFactorFromDefaultLoginPolicyResponse> {
        return this.grpcService.admin.removeMultiFactorFromDefaultLoginPolicy(req);
    }

    public listDefaultLoginPolicySecondFactors(): Promise<ListDefaultLoginPolicySecondFactorsResponse> {
        const req = new ListDefaultLoginPolicyMultiFactorsRequest();
        return this.grpcService.admin.listDefaultLoginPolicySecondFactors(req);
    }

    public AddSecondFactorToDefaultLoginPolicy(req: AddSecondFactorToDefaultLoginPolicyRequest): Promise<AddSecondFactorToDefaultLoginPolicyResponse> {
        return this.grpcService.admin.addSecondFactorToDefaultLoginPolicy(req);
    }

    public removeSecondFactorFromDefaultLoginPolicy(req: RemoveSecondFactorFromDefaultLoginPolicyRequest): Promise<RemoveSecondFactorFromDefaultLoginPolicyResponse> {
        return this.grpcService.admin.removeSecondFactorFromDefaultLoginPolicy(req);
    }

    public listIAMMemberRoles(): Promise<ListIAMMemberRolesResponse> {
        const req = new ListIAMMemberRolesRequest();
        return this.grpcService.admin.listIAMMemberRoles(req);
    }

    public listViews(): Promise<ListViewsResponse> {
        const req = new ListViewsRequest();
        return this.grpcService.admin.listViews(req);
    }

    public listFailedEvents(): Promise<ListFailedEventsResponse> {
        const req = new ListFailedEventsRequest();
        return this.grpcService.admin.listFailedEvents(req);
    }

    public clearView(viewname: string, db: string): Promise<ClearViewResponse> {
        const req = new ClearViewRequest();
        req.setDatabase(db);
        req.setViewName(viewname);
        return this.grpcService.admin.clearView(req);
    }

    public removeFailedEvent(viewname: string, db: string, sequence: number): Promise<RemoveFailedEventResponse> {
        const req = new RemoveFailedEventRequest();
        req.setDatabase(db);
        req.setViewName(viewname);
        req.setFailedSequence(sequence);
        return this.grpcService.admin.removeFailedEvent(req);
    }

    /* Policies */

    /* complexity */

    public getDefaultPasswordComplexityPolicy(): Promise<GetDefaultPasswordComplexityPolicyResponse> {
        const req = new GetDefaultPasswordComplexityPolicyRequest();
        return this.grpcService.admin.getDefaultPasswordComplexityPolicy(req);
    }

    public updateDefaultPasswordComplexityPolicy(
        hasLowerCase: boolean,
        hasUpperCase: boolean,
        hasNumber: boolean,
        hasSymbol: boolean,
        minLength: number,
    ): Promise<UpdateDefaultPasswordComplexityPolicyResponse> {
        const req = new UpdateDefaultPasswordComplexityPolicyRequest();
        req.setHasLowercase(hasLowerCase);
        req.setHasUppercase(hasUpperCase);
        req.setHasNumber(hasNumber);
        req.setHasSymbol(hasSymbol);
        req.setMinLength(minLength);
        return this.grpcService.admin.updateDefaultPasswordComplexityPolicy(req);
    }

    /* age */

    public getDefaultPasswordAgePolicy(): Promise<GetDefaultPasswordAgePolicyResponse> {
        const req = new Empty();

        return this.grpcService.admin.getDefaultPasswordAgePolicy(req);
    }

    public updateDefaultPasswordAgePolicy(
        maxAgeDays: number,
        expireWarnDays: number,
    ): Promise<UpdateDefaultPasswordAgePolicyResponse> {
        const req = new UpdateDefaultPasswordAgePolicyRequest();
        req.setMaxAgeDays(maxAgeDays);
        req.setExpireWarnDays(expireWarnDays);

        return this.grpcService.admin.updateDefaultPasswordAgePolicy(req);
    }

    /* lockout */

    public getDefaultPasswordLockoutPolicy(): Promise<GetDefaultPasswordLockoutPolicyResponse> {
        const req = new GetDefaultPasswordLockoutPolicyRequest();
        return this.grpcService.admin.getDefaultPasswordLockoutPolicy(req);
    }

    public UpdateDefaultPasswordLockoutPolicy(
        maxAttempts: number,
        showLockoutFailures: boolean,
    ): Promise<UpdateDefaultPasswordLockoutPolicyResponse> {
        const req = new UpdateDefaultPasswordLockoutPolicyRequest();
        req.setMaxAttempts(maxAttempts);
        req.setShowLockoutFailure(showLockoutFailures);

        return this.grpcService.admin.updateDefaultPasswordLockoutPolicy(req);
    }

    /* label */

    public getDefaultLabelPolicy(): Promise<GetDefaultLabelPolicyResponse> {
        const req = new GetDefaultLabelPolicyRequest();
        return this.grpcService.admin.getDefaultLabelPolicy(req);
    }

    public UpdateDefaultLabelPolicy(req: UpdateDefaultLabelPolicyRequest): Promise<UpdateDefaultLabelPolicyResponse> {
        return this.grpcService.admin.updateDefaultLabelPolicy(req);
    }

    /* login */

    public getDefaultLoginPolicy(
    ): Promise<GetDefaultLoginPolicyResponse> {
        const req = new GetDefaultLoginPolicyRequest();
        return this.grpcService.admin.getDefaultLoginPolicy(req);
    }

    public UpdateDefaultLoginPolicy(req: UpdateDefaultLoginPolicyRequest): Promise<UpdateDefaultLoginPolicyResponse> {
        return this.grpcService.admin.updateDefaultLoginPolicy(req);
    }

    /* org iam */

    public getOrgIAMPolicy(orgId: string): Promise<GetOrgIAMPolicyResponse> {
        const req = new GetOrgIAMPolicyRequest();
        req.setOrgId(orgId);
        return this.grpcService.admin.getOrgIAMPolicy(req);
    }

    public addCustomOrgIAMPolicy(
        orgId: string,
        userLoginMustBeDomain: boolean): Promise<AddCustomOrgIAMPolicyResponse> {
        const req = new AddCustomOrgIAMPolicyRequest();
        req.setOrgId(orgId);
        req.setUserLoginMustBeDomain(userLoginMustBeDomain);

        return this.grpcService.admin.addCustomOrgIAMPolicy(req);
    }

    public updateCustomOrgIAMPolicy(
        orgId: string,
        userLoginMustBeDomain: boolean): Promise<UpdateCustomOrgIAMPolicyResponse> {
        const req = new UpdateCustomOrgIAMPolicyRequest();
        req.setOrgId(orgId);
        req.setUserLoginMustBeDomain(userLoginMustBeDomain);
        return this.grpcService.admin.updateCustomOrgIAMPolicy(req);
    }

    public resetOrgIAMPolicyToDefault(
        orgId: string,
    ): Promise<ResetOrgIAMPolicyToDefaultResponse> {
        const req = new ResetOrgIAMPolicyToDefaultRequest();
        req.setOrgId(orgId);
        return this.grpcService.admin.resetOrgIAMPolicyToDefault(req);
    }

    /* admin iam */

    public getDefaultOrgIAMPolicy(): Promise<GetDefaultOrgIAMPolicyResponse> {
        const req = new GetDefaultOrgIAMPolicyRequest();
        return this.grpcService.admin.getDefaultOrgIAMPolicy(req);
    }

    /* policies end */

    public addIDPToDefaultLoginPolicy(idpId: string): Promise<AddIDPToDefaultLoginPolicyResponse> {
        const req = new AddIDPToDefaultLoginPolicyRequest();
        req.setIdpId(idpId);
        return this.grpcService.admin.addIDPToDefaultLoginPolicy(req);
    }

    public removeIDPFromDefaultLoginPolicy(idpId: string): Promise<RemoveIDPFromDefaultLoginPolicyResponse> {
        const req = new RemoveIDPFromDefaultLoginPolicyRequest();
        req.setIdpId(idpId);
        return this.grpcService.admin.removeIDPFromDefaultLoginPolicy(req);
    }

    public listDefaultLoginPolicyIDPs(limit?: number, offset?: number): Promise<ListDefaultLoginPolicyIDPsResponse> {
        const req = new ListDefaultLoginPolicyIDPsRequest();
        const query = new ListQuery();
        if (limit) {
            query.setLimit(limit);
        }
        if (offset) {
            query.setOffset(offset);
        }
        req.setMetaData(query);
        return this.grpcService.admin.listDefaultLoginPolicyIDPs(req);
    }

    public listIDPs(
        limit?: number,
        offset?: number,
        queriesList?: IDPQuery[],
    ): Promise<ListIDPsResponse> {
        const req = new ListIDPsRequest();
        const query = new ListQuery();

        if (limit) {
            query.setLimit(limit);
        }
        if (offset) {
            query.setOffset(offset);
        }
        if (queriesList) {
            req.setQueriesList(queriesList);
        }
        req.setMetaData(query);
        return this.grpcService.admin.listIDPs(req);
    }

    public getIDPByID(
        id: string,
    ): Promise<GetIDPByIDResponse> {
        const req = new GetIDPByIDRequest();
        req.setId(id);
        return this.grpcService.admin.getIDPByID(req);
    }

    public updateIDP(
        req: UpdateIDPRequest,
    ): Promise<UpdateIDPResponse> {
        return this.grpcService.admin.updateIDP(req);
    }

    public addOIDCIDP(
        req: AddOIDCIDPRequest,
    ): Promise<AddOIDCIDPResponse> {
        return this.grpcService.admin.addOIDCIDP(req);
    }

    public updateIDPOIDCConfig(
        req: UpdateIDPOIDCConfigRequest,
    ): Promise<UpdateIDPOIDCConfigResponse> {
        return this.grpcService.admin.updateIDPOIDCConfig(req);
    }

    public removeIDP(
        id: string,
    ): Promise<RemoveIDPResponse> {
        const req = new RemoveIDPRequest;
        req.setId(id);
        return this.grpcService.admin.removeIDP(req);
    }

    public deactivateIDP(
        id: string,
    ): Promise<DeactivateIDPResponse> {
        const req = new DeactivateIDPRequest;
        req.setId(id);
        return this.grpcService.admin.deactivateIDP(req);
    }

    public reactivateIDP(
        id: string,
    ): Promise<ReactivateIDPResponse> {
        const req = new ReactivateIDPRequest;
        req.setId(id);
        return this.grpcService.admin.reactivateIDP(req);
    }

    public listIAMMembers(
        limit: number,
        offset: number,
        query?: SearchQuery,
    ): Promise<ListIAMMembersResponse> {
        const req = new ListIAMMembersRequest();
        const metadata = new ListQuery();
        if (limit) {
            metadata.setLimit(limit);
        }
        if (offset) {
            metadata.setOffset(offset);
        }
        if (query) {
            req.setQuery(query);
        }
        req.setMetaData(metadata);

        return this.grpcService.admin.listIAMMembers(req);
    }

    public removeIAMMember(
        userId: string,
    ): Promise<RemoveIAMMemberResponse> {
        const req = new RemoveIAMMemberRequest();
        req.setUserId(userId);
        return this.grpcService.admin.removeIAMMember(req);
    }

    public addIAMMember(
        userId: string,
        rolesList: string[],
    ): Promise<AddIAMMemberResponse> {
        const req = new AddIAMMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);

        return this.grpcService.admin.addIAMMember(req);
    }

    public updateIAMMember(
        userId: string,
        rolesList: string[],
    ): Promise<UpdateIAMMemberResponse> {
        const req = new UpdateIAMMemberRequest();
        req.setUserId(userId);
        req.setRolesList(rolesList);

        return this.grpcService.admin.updateIAMMember(req);
    }
}

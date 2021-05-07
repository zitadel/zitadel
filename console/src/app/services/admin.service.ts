import { Injectable } from '@angular/core';

import {
  AddCustomOrgIAMPolicyRequest,
  AddCustomOrgIAMPolicyResponse,
  AddIAMMemberRequest,
  AddIAMMemberResponse,
  AddIDPToLoginPolicyRequest,
  AddIDPToLoginPolicyResponse,
  AddMultiFactorToLoginPolicyRequest,
  AddMultiFactorToLoginPolicyResponse,
  AddOIDCIDPRequest,
  AddOIDCIDPResponse,
  AddSecondFactorToLoginPolicyRequest,
  AddSecondFactorToLoginPolicyResponse,
  ClearViewRequest,
  ClearViewResponse,
  DeactivateIDPRequest,
  DeactivateIDPResponse,
  GetCustomOrgIAMPolicyRequest,
  GetCustomOrgIAMPolicyResponse,
  GetDefaultFeaturesRequest,
  GetDefaultFeaturesResponse,
  GetIDPByIDRequest,
  GetIDPByIDResponse,
  GetLabelPolicyRequest,
  GetLabelPolicyResponse,
  GetLoginPolicyRequest,
  GetLoginPolicyResponse,
  GetOrgFeaturesRequest,
  GetOrgFeaturesResponse,
  GetOrgIAMPolicyRequest,
  GetOrgIAMPolicyResponse,
  GetPasswordAgePolicyRequest,
  GetPasswordAgePolicyResponse,
  GetPasswordComplexityPolicyRequest,
  GetPasswordComplexityPolicyResponse,
  GetPasswordLockoutPolicyRequest,
  GetPasswordLockoutPolicyResponse,
  GetPreviewLabelPolicyRequest,
  GetPreviewLabelPolicyResponse,
  IDPQuery,
  ListFailedEventsRequest,
  ListFailedEventsResponse,
  ListIAMMemberRolesRequest,
  ListIAMMemberRolesResponse,
  ListIAMMembersRequest,
  ListIAMMembersResponse,
  ListIDPsRequest,
  ListIDPsResponse,
  ListLoginPolicyIDPsRequest,
  ListLoginPolicyIDPsResponse,
  ListLoginPolicyMultiFactorsRequest,
  ListLoginPolicyMultiFactorsResponse,
  ListLoginPolicySecondFactorsRequest,
  ListLoginPolicySecondFactorsResponse,
  ListViewsRequest,
  ListViewsResponse,
  ReactivateIDPRequest,
  ReactivateIDPResponse,
  RemoveFailedEventRequest,
  RemoveFailedEventResponse,
  RemoveIAMMemberRequest,
  RemoveIAMMemberResponse,
  RemoveIDPFromLoginPolicyRequest,
  RemoveIDPFromLoginPolicyResponse,
  RemoveIDPRequest,
  RemoveIDPResponse,
  RemoveMultiFactorFromLoginPolicyRequest,
  RemoveMultiFactorFromLoginPolicyResponse,
  RemoveSecondFactorFromLoginPolicyRequest,
  RemoveSecondFactorFromLoginPolicyResponse,
  ResetCustomOrgIAMPolicyToDefaultRequest,
  ResetCustomOrgIAMPolicyToDefaultResponse,
  ResetOrgFeaturesRequest,
  ResetOrgFeaturesResponse,
  SetDefaultFeaturesRequest,
  SetDefaultFeaturesResponse,
  SetOrgFeaturesRequest,
  SetOrgFeaturesResponse,
  SetUpOrgRequest,
  SetUpOrgResponse,
  UpdateCustomOrgIAMPolicyRequest,
  UpdateCustomOrgIAMPolicyResponse,
  UpdateIAMMemberRequest,
  UpdateIAMMemberResponse,
  UpdateIDPOIDCConfigRequest,
  UpdateIDPOIDCConfigResponse,
  UpdateIDPRequest,
  UpdateIDPResponse,
  UpdateLabelPolicyRequest,
  UpdateLabelPolicyResponse,
  UpdateLoginPolicyRequest,
  UpdateLoginPolicyResponse,
  UpdateOrgIAMPolicyRequest,
  UpdateOrgIAMPolicyResponse,
  UpdatePasswordAgePolicyRequest,
  UpdatePasswordAgePolicyResponse,
  UpdatePasswordComplexityPolicyRequest,
  UpdatePasswordComplexityPolicyResponse,
  UpdatePasswordLockoutPolicyRequest,
  UpdatePasswordLockoutPolicyResponse,
} from '../proto/generated/zitadel/admin_pb';
import {
  ActivateCustomLabelPolicyRequest,
  ActivateCustomLabelPolicyResponse,
} from '../proto/generated/zitadel/management_pb';
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
  ): Promise<SetUpOrgResponse.AsObject> {
    const req = new SetUpOrgRequest();

    req.setOrg(org);
    req.setHuman(human);

    return this.grpcService.admin.setUpOrg(req, null).then(resp => resp.toObject());
  }

  public listLoginPolicyMultiFactors(): Promise<ListLoginPolicyMultiFactorsResponse.AsObject> {
    const req = new ListLoginPolicyMultiFactorsRequest();
    return this.grpcService.admin.listLoginPolicyMultiFactors(req, null).then(resp => resp.toObject());
  }

  public addMultiFactorToLoginPolicy(req: AddMultiFactorToLoginPolicyRequest):
    Promise<AddMultiFactorToLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.addMultiFactorToLoginPolicy(req, null).then(resp => resp.toObject());
  }

  public removeMultiFactorFromLoginPolicy(req: RemoveMultiFactorFromLoginPolicyRequest):
    Promise<RemoveMultiFactorFromLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.removeMultiFactorFromLoginPolicy(req, null).then(resp => resp.toObject());
  }

  public listLoginPolicySecondFactors(): Promise<ListLoginPolicySecondFactorsResponse.AsObject> {
    const req = new ListLoginPolicySecondFactorsRequest();
    return this.grpcService.admin.listLoginPolicySecondFactors(req, null).then(resp => resp.toObject());
  }

  public addSecondFactorToLoginPolicy(req: AddSecondFactorToLoginPolicyRequest):
    Promise<AddSecondFactorToLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.addSecondFactorToLoginPolicy(req, null).then(resp => resp.toObject());
  }

  public removeSecondFactorFromLoginPolicy(req: RemoveSecondFactorFromLoginPolicyRequest):
    Promise<RemoveSecondFactorFromLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.removeSecondFactorFromLoginPolicy(req, null).then(resp => resp.toObject());
  }

  public listIAMMemberRoles(): Promise<ListIAMMemberRolesResponse.AsObject> {
    const req = new ListIAMMemberRolesRequest();
    return this.grpcService.admin.listIAMMemberRoles(req, null).then(resp => resp.toObject());
  }

  public listViews(): Promise<ListViewsResponse.AsObject> {
    const req = new ListViewsRequest();
    return this.grpcService.admin.listViews(req, null).then(resp => resp.toObject());
  }

  public listFailedEvents(): Promise<ListFailedEventsResponse.AsObject> {
    const req = new ListFailedEventsRequest();
    return this.grpcService.admin.listFailedEvents(req, null).then(resp => resp.toObject());
  }

  public clearView(viewname: string, db: string): Promise<ClearViewResponse.AsObject> {
    const req = new ClearViewRequest();
    req.setDatabase(db);
    req.setViewName(viewname);
    return this.grpcService.admin.clearView(req, null).then(resp => resp.toObject());
  }

  public removeFailedEvent(viewname: string, db: string, sequence: number): Promise<RemoveFailedEventResponse.AsObject> {
    const req = new RemoveFailedEventRequest();
    req.setDatabase(db);
    req.setViewName(viewname);
    req.setFailedSequence(sequence);
    return this.grpcService.admin.removeFailedEvent(req, null).then(resp => resp.toObject());
  }

  // Features

  public getOrgFeatures(orgId: string): Promise<GetOrgFeaturesResponse.AsObject> {
    const req = new GetOrgFeaturesRequest();
    req.setOrgId(orgId);
    return this.grpcService.admin.getOrgFeatures(req, null).then(resp => resp.toObject());
  }

  public setOrgFeatures(req: SetOrgFeaturesRequest): Promise<SetOrgFeaturesResponse.AsObject> {
    return this.grpcService.admin.setOrgFeatures(req, null).then(resp => resp.toObject());
  }

  public resetOrgFeatures(orgId: string): Promise<ResetOrgFeaturesResponse.AsObject> {
    const req = new ResetOrgFeaturesRequest();
    req.setOrgId(orgId);
    return this.grpcService.admin.resetOrgFeatures(req, null).then(resp => resp.toObject());
  }

  public getDefaultFeatures(): Promise<GetDefaultFeaturesResponse.AsObject> {
    const req = new GetDefaultFeaturesRequest();
    return this.grpcService.admin.getDefaultFeatures(req, null).then(resp => resp.toObject());
  }

  public setDefaultFeatures(req: SetDefaultFeaturesRequest): Promise<SetDefaultFeaturesResponse.AsObject> {
    return this.grpcService.admin.setDefaultFeatures(req, null).then(resp => resp.toObject());
  }

  /* Policies */

  /* complexity */

  public getPasswordComplexityPolicy(): Promise<GetPasswordComplexityPolicyResponse.AsObject> {
    const req = new GetPasswordComplexityPolicyRequest();
    return this.grpcService.admin.getPasswordComplexityPolicy(req, null).then(resp => resp.toObject());
  }

  public updatePasswordComplexityPolicy(
    hasLowerCase: boolean,
    hasUpperCase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
    minLength: number,
  ): Promise<UpdatePasswordComplexityPolicyResponse.AsObject> {
    const req = new UpdatePasswordComplexityPolicyRequest();
    req.setHasLowercase(hasLowerCase);
    req.setHasUppercase(hasUpperCase);
    req.setHasNumber(hasNumber);
    req.setHasSymbol(hasSymbol);
    req.setMinLength(minLength);
    return this.grpcService.admin.updatePasswordComplexityPolicy(req, null).then(resp => resp.toObject());
  }

  /* age */

  public getPasswordAgePolicy(): Promise<GetPasswordAgePolicyResponse.AsObject> {
    const req = new GetPasswordAgePolicyRequest();

    return this.grpcService.admin.getPasswordAgePolicy(req, null).then(resp => resp.toObject());
  }

  public updatePasswordAgePolicy(
    maxAgeDays: number,
    expireWarnDays: number,
  ): Promise<UpdatePasswordAgePolicyResponse.AsObject> {
    const req = new UpdatePasswordAgePolicyRequest();
    req.setMaxAgeDays(maxAgeDays);
    req.setExpireWarnDays(expireWarnDays);

    return this.grpcService.admin.updatePasswordAgePolicy(req, null).then(resp => resp.toObject());
  }

  /* lockout */

  public getPasswordLockoutPolicy(): Promise<GetPasswordLockoutPolicyResponse.AsObject> {
    const req = new GetPasswordLockoutPolicyRequest();
    return this.grpcService.admin.getPasswordLockoutPolicy(req, null).then(resp => resp.toObject());
  }

  public updatePasswordLockoutPolicy(
    maxAttempts: number,
    showLockoutFailures: boolean,
  ): Promise<UpdatePasswordLockoutPolicyResponse.AsObject> {
    const req = new UpdatePasswordLockoutPolicyRequest();
    req.setMaxAttempts(maxAttempts);
    req.setShowLockoutFailure(showLockoutFailures);

    return this.grpcService.admin.updatePasswordLockoutPolicy(req, null).then(resp => resp.toObject());
  }

  /* label */

  public getLabelPolicy(): Promise<GetLabelPolicyResponse.AsObject> {
    const req = new GetLabelPolicyRequest();
    return this.grpcService.admin.getLabelPolicy(req, null).then(resp => resp.toObject());
  }

  public updateLabelPolicy(req: UpdateLabelPolicyRequest): Promise<UpdateLabelPolicyResponse.AsObject> {
    return this.grpcService.admin.updateLabelPolicy(req, null).then(resp => resp.toObject());
  }

  public getPreviewLabelPolicy(req: GetPreviewLabelPolicyRequest): Promise<GetPreviewLabelPolicyResponse.AsObject> {
    return this.grpcService.admin.getPreviewLabelPolicy(req, null).then(resp => resp.toObject());
  }

  public activateCustomLabelPolicy(req: ActivateCustomLabelPolicyRequest): Promise<ActivateCustomLabelPolicyResponse.AsObject> {
    return this.grpcService.admin.activateCustomLabelPolicy(req, null).then(resp => resp.toObject());
  }

  /* login */

  public getLoginPolicy(
  ): Promise<GetLoginPolicyResponse.AsObject> {
    const req = new GetLoginPolicyRequest();
    return this.grpcService.admin.getLoginPolicy(req, null).then(resp => resp.toObject());
  }

  public updateLoginPolicy(req: UpdateLoginPolicyRequest): Promise<UpdateLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.updateLoginPolicy(req, null).then(resp => resp.toObject());
  }

  /* org iam */

  public getCustomOrgIAMPolicy(orgId: string): Promise<GetCustomOrgIAMPolicyResponse.AsObject> {
    const req = new GetCustomOrgIAMPolicyRequest();
    req.setOrgId(orgId);
    return this.grpcService.admin.getCustomOrgIAMPolicy(req, null).then(resp => resp.toObject());
  }

  public addCustomOrgIAMPolicy(
    orgId: string,
    userLoginMustBeDomain: boolean): Promise<AddCustomOrgIAMPolicyResponse.AsObject> {
    const req = new AddCustomOrgIAMPolicyRequest();
    req.setOrgId(orgId);
    req.setUserLoginMustBeDomain(userLoginMustBeDomain);

    return this.grpcService.admin.addCustomOrgIAMPolicy(req, null).then(resp => resp.toObject());
  }

  public updateCustomOrgIAMPolicy(
    orgId: string,
    userLoginMustBeDomain: boolean): Promise<UpdateCustomOrgIAMPolicyResponse.AsObject> {
    const req = new UpdateCustomOrgIAMPolicyRequest();
    req.setOrgId(orgId);
    req.setUserLoginMustBeDomain(userLoginMustBeDomain);
    return this.grpcService.admin.updateCustomOrgIAMPolicy(req, null).then(resp => resp.toObject());
  }

  public resetCustomOrgIAMPolicyToDefault(
    orgId: string,
  ): Promise<ResetCustomOrgIAMPolicyToDefaultResponse.AsObject> {
    const req = new ResetCustomOrgIAMPolicyToDefaultRequest();
    req.setOrgId(orgId);
    return this.grpcService.admin.resetCustomOrgIAMPolicyToDefault(req, null).then(resp => resp.toObject());
  }

  /* admin iam */

  public getOrgIAMPolicy(): Promise<GetOrgIAMPolicyResponse.AsObject> {
    const req = new GetOrgIAMPolicyRequest();
    return this.grpcService.admin.getOrgIAMPolicy(req, null).then(resp => resp.toObject());
  }

  public updateOrgIAMPolicy(userLoginMustBeDomain: boolean): Promise<UpdateOrgIAMPolicyResponse.AsObject> {
    const req = new UpdateOrgIAMPolicyRequest();
    req.setUserLoginMustBeDomain(userLoginMustBeDomain);
    return this.grpcService.admin.updateOrgIAMPolicy(req, null).then(resp => resp.toObject());
  }

  /* policies end */

  public addIDPToLoginPolicy(idpId: string): Promise<AddIDPToLoginPolicyResponse.AsObject> {
    const req = new AddIDPToLoginPolicyRequest();
    req.setIdpId(idpId);
    return this.grpcService.admin.addIDPToLoginPolicy(req, null).then(resp => resp.toObject());
  }

  public removeIDPFromLoginPolicy(idpId: string): Promise<RemoveIDPFromLoginPolicyResponse.AsObject> {
    const req = new RemoveIDPFromLoginPolicyRequest();
    req.setIdpId(idpId);
    return this.grpcService.admin.removeIDPFromLoginPolicy(req, null).then(resp => resp.toObject());
  }

  public listLoginPolicyIDPs(limit?: number, offset?: number): Promise<ListLoginPolicyIDPsResponse.AsObject> {
    const req = new ListLoginPolicyIDPsRequest();
    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }
    req.setQuery(query);
    return this.grpcService.admin.listLoginPolicyIDPs(req, null).then(resp => resp.toObject());
  }

  public listIDPs(
    limit?: number,
    offset?: number,
    queriesList?: IDPQuery[],
  ): Promise<ListIDPsResponse.AsObject> {
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
    req.setQuery(query);
    return this.grpcService.admin.listIDPs(req, null).then(resp => resp.toObject());
  }

  public getIDPByID(
    id: string,
  ): Promise<GetIDPByIDResponse.AsObject> {
    const req = new GetIDPByIDRequest();
    req.setId(id);
    return this.grpcService.admin.getIDPByID(req, null).then(resp => resp.toObject());
  }

  public updateIDP(
    req: UpdateIDPRequest,
  ): Promise<UpdateIDPResponse.AsObject> {
    return this.grpcService.admin.updateIDP(req, null).then(resp => resp.toObject());
  }

  public addOIDCIDP(
    req: AddOIDCIDPRequest,
  ): Promise<AddOIDCIDPResponse.AsObject> {
    return this.grpcService.admin.addOIDCIDP(req, null).then(resp => resp.toObject());
  }

  public updateIDPOIDCConfig(
    req: UpdateIDPOIDCConfigRequest,
  ): Promise<UpdateIDPOIDCConfigResponse.AsObject> {
    return this.grpcService.admin.updateIDPOIDCConfig(req, null).then(resp => resp.toObject());
  }

  public removeIDP(
    id: string,
  ): Promise<RemoveIDPResponse.AsObject> {
    const req = new RemoveIDPRequest;
    req.setIdpId(id);
    return this.grpcService.admin.removeIDP(req, null).then(resp => resp.toObject());
  }

  public deactivateIDP(
    id: string,
  ): Promise<DeactivateIDPResponse.AsObject> {
    const req = new DeactivateIDPRequest;
    req.setIdpId(id);
    return this.grpcService.admin.deactivateIDP(req, null).then(resp => resp.toObject());
  }

  public reactivateIDP(
    id: string,
  ): Promise<ReactivateIDPResponse.AsObject> {
    const req = new ReactivateIDPRequest;
    req.setIdpId(id);
    return this.grpcService.admin.reactivateIDP(req, null).then(resp => resp.toObject());
  }

  public listIAMMembers(
    limit: number,
    offset: number,
    queriesList?: SearchQuery[],
  ): Promise<ListIAMMembersResponse.AsObject> {
    const req = new ListIAMMembersRequest();
    const metadata = new ListQuery();
    if (limit) {
      metadata.setLimit(limit);
    }
    if (offset) {
      metadata.setOffset(offset);
    }
    if (queriesList) {
      req.setQueriesList(queriesList);
    }
    req.setQuery(metadata);

    return this.grpcService.admin.listIAMMembers(req, null).then(resp => resp.toObject());
  }

  public removeIAMMember(
    userId: string,
  ): Promise<RemoveIAMMemberResponse.AsObject> {
    const req = new RemoveIAMMemberRequest();
    req.setUserId(userId);
    return this.grpcService.admin.removeIAMMember(req, null).then(resp => resp.toObject());
  }

  public addIAMMember(
    userId: string,
    rolesList: string[],
  ): Promise<AddIAMMemberResponse.AsObject> {
    const req = new AddIAMMemberRequest();
    req.setUserId(userId);
    req.setRolesList(rolesList);

    return this.grpcService.admin.addIAMMember(req, null).then(resp => resp.toObject());
  }

  public updateIAMMember(
    userId: string,
    rolesList: string[],
  ): Promise<UpdateIAMMemberResponse.AsObject> {
    const req = new UpdateIAMMemberRequest();
    req.setUserId(userId);
    req.setRolesList(rolesList);

    return this.grpcService.admin.updateIAMMember(req, null).then(resp => resp.toObject());
  }
}

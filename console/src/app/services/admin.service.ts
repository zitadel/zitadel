import { Injectable } from '@angular/core';

import {
  ActivateLabelPolicyRequest,
  ActivateLabelPolicyResponse,
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
  GetCustomDomainClaimedMessageTextRequest,
  GetCustomDomainClaimedMessageTextResponse,
  GetCustomInitMessageTextRequest,
  GetCustomInitMessageTextResponse,
  GetCustomLoginTextsRequest,
  GetCustomLoginTextsResponse,
  GetCustomOrgIAMPolicyRequest,
  GetCustomOrgIAMPolicyResponse,
  GetCustomPasswordResetMessageTextRequest,
  GetCustomPasswordResetMessageTextResponse,
  GetCustomVerifyEmailMessageTextRequest,
  GetCustomVerifyEmailMessageTextResponse,
  GetCustomVerifyPhoneMessageTextRequest,
  GetCustomVerifyPhoneMessageTextResponse,
  GetDefaultDomainClaimedMessageTextRequest,
  GetDefaultDomainClaimedMessageTextResponse,
  GetDefaultFeaturesRequest,
  GetDefaultFeaturesResponse,
  GetDefaultInitMessageTextRequest,
  GetDefaultInitMessageTextResponse,
  GetDefaultLoginTextsRequest,
  GetDefaultLoginTextsResponse,
  GetDefaultPasswordResetMessageTextRequest,
  GetDefaultPasswordResetMessageTextResponse,
  GetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyEmailMessageTextResponse,
  GetDefaultVerifyPhoneMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextResponse,
  GetIDPByIDRequest,
  GetIDPByIDResponse,
  GetLabelPolicyRequest,
  GetLabelPolicyResponse,
  GetLockoutPolicyRequest,
  GetLockoutPolicyResponse,
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
  GetPreviewLabelPolicyRequest,
  GetPreviewLabelPolicyResponse,
  GetPrivacyPolicyRequest,
  GetPrivacyPolicyResponse,
  GetSupportedLanguagesRequest,
  GetSupportedLanguagesResponse,
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
  RemoveLabelPolicyFontRequest,
  RemoveLabelPolicyFontResponse,
  RemoveLabelPolicyIconDarkRequest,
  RemoveLabelPolicyIconDarkResponse,
  RemoveLabelPolicyIconRequest,
  RemoveLabelPolicyIconResponse,
  RemoveLabelPolicyLogoDarkRequest,
  RemoveLabelPolicyLogoDarkResponse,
  RemoveLabelPolicyLogoRequest,
  RemoveLabelPolicyLogoResponse,
  RemoveMultiFactorFromLoginPolicyRequest,
  RemoveMultiFactorFromLoginPolicyResponse,
  RemoveSecondFactorFromLoginPolicyRequest,
  RemoveSecondFactorFromLoginPolicyResponse,
  ResetCustomLoginTextsToDefaultRequest,
  ResetCustomLoginTextsToDefaultResponse,
  ResetCustomOrgIAMPolicyToDefaultRequest,
  ResetCustomOrgIAMPolicyToDefaultResponse,
  ResetOrgFeaturesRequest,
  ResetOrgFeaturesResponse,
  SetCustomLoginTextsRequest,
  SetCustomLoginTextsResponse,
  SetDefaultDomainClaimedMessageTextRequest,
  SetDefaultDomainClaimedMessageTextResponse,
  SetDefaultFeaturesRequest,
  SetDefaultFeaturesResponse,
  SetDefaultInitMessageTextRequest,
  SetDefaultInitMessageTextResponse,
  SetDefaultPasswordResetMessageTextRequest,
  SetDefaultPasswordResetMessageTextResponse,
  SetDefaultVerifyEmailMessageTextRequest,
  SetDefaultVerifyEmailMessageTextResponse,
  SetDefaultVerifyPhoneMessageTextRequest,
  SetDefaultVerifyPhoneMessageTextResponse,
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
  UpdateLockoutPolicyRequest,
  UpdateLockoutPolicyResponse,
  UpdateLoginPolicyRequest,
  UpdateLoginPolicyResponse,
  UpdateOrgIAMPolicyRequest,
  UpdateOrgIAMPolicyResponse,
  UpdatePasswordAgePolicyRequest,
  UpdatePasswordAgePolicyResponse,
  UpdatePasswordComplexityPolicyRequest,
  UpdatePasswordComplexityPolicyResponse,
  UpdatePrivacyPolicyRequest,
  UpdatePrivacyPolicyResponse,
} from '../proto/generated/zitadel/admin_pb';
import { SearchQuery } from '../proto/generated/zitadel/member_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { GrpcService } from './grpc.service';

@Injectable({
  providedIn: 'root',
})
export class AdminService {
  constructor(private readonly grpcService: GrpcService) { }


  public getSupportedLanguages(): Promise<GetSupportedLanguagesResponse.AsObject> {
    const req = new GetSupportedLanguagesRequest();
    return this.grpcService.admin.getSupportedLanguages(req, null).then(resp => resp.toObject());
  }

  public getDefaultLoginTexts(req: GetDefaultLoginTextsRequest):
    Promise<GetDefaultLoginTextsResponse.AsObject> {
    return this.grpcService.admin.getDefaultLoginTexts(req, null).then(resp => resp.toObject());
  }

  public getCustomLoginTexts(req: GetCustomLoginTextsRequest):
    Promise<GetCustomLoginTextsResponse.AsObject> {
    return this.grpcService.admin.getCustomLoginTexts(req, null).then(resp => resp.toObject());
  }

  public setCustomLoginText(req: SetCustomLoginTextsRequest):
    Promise<SetCustomLoginTextsResponse.AsObject> {
    return this.grpcService.admin.setCustomLoginText(req, null).then(resp => resp.toObject());
  }

  public resetCustomLoginTextToDefault(lang: string): Promise<ResetCustomLoginTextsToDefaultResponse.AsObject> {
    const req = new ResetCustomLoginTextsToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomLoginTextToDefault(req, null).then(resp => resp.toObject());
  }

  // message texts

  public getDefaultInitMessageText(req: GetDefaultInitMessageTextRequest):
    Promise<GetDefaultInitMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultInitMessageText(req, null).then(resp => resp.toObject());
  }

  public getCustomInitMessageText(req: GetCustomInitMessageTextRequest):
    Promise<GetCustomInitMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomInitMessageText(req, null).then(resp => resp.toObject());
  }

  public setDefaultInitMessageText(req: SetDefaultInitMessageTextRequest):
    Promise<SetDefaultInitMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultInitMessageText(req, null).then(resp => resp.toObject());
  }


  public getDefaultVerifyEmailMessageText(req: GetDefaultVerifyEmailMessageTextRequest):
    Promise<GetDefaultVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultVerifyEmailMessageText(req, null).then(resp => resp.toObject());
  }

  public getCustomVerifyEmailMessageText(req: GetCustomVerifyEmailMessageTextRequest):
    Promise<GetCustomVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomVerifyEmailMessageText(req, null).then(resp => resp.toObject());
  }

  public setDefaultVerifyEmailMessageText(req: SetDefaultVerifyEmailMessageTextRequest):
    Promise<SetDefaultVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultVerifyEmailMessageText(req, null).then(resp => resp.toObject());
  }


  public getDefaultVerifyPhoneMessageText(req: GetDefaultVerifyPhoneMessageTextRequest):
    Promise<GetDefaultVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultVerifyPhoneMessageText(req, null).then(resp => resp.toObject());
  }

  public getCustomVerifyPhoneMessageText(req: GetCustomVerifyPhoneMessageTextRequest):
    Promise<GetCustomVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomVerifyPhoneMessageText(req, null).then(resp => resp.toObject());
  }

  public setDefaultVerifyPhoneMessageText(req: SetDefaultVerifyPhoneMessageTextRequest):
    Promise<SetDefaultVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultVerifyPhoneMessageText(req, null).then(resp => resp.toObject());
  }


  public getDefaultPasswordResetMessageText(req: GetDefaultPasswordResetMessageTextRequest):
    Promise<GetDefaultPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultPasswordResetMessageText(req, null).then(resp => resp.toObject());
  }

  public getCustomPasswordResetMessageText(req: GetCustomPasswordResetMessageTextRequest):
    Promise<GetCustomPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomPasswordResetMessageText(req, null).then(resp => resp.toObject());
  }

  public setDefaultPasswordResetMessageText(req: SetDefaultPasswordResetMessageTextRequest):
    Promise<SetDefaultPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultPasswordResetMessageText(req, null).then(resp => resp.toObject());
  }


  public getDefaultDomainClaimedMessageText(req: GetDefaultDomainClaimedMessageTextRequest):
    Promise<GetDefaultDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultDomainClaimedMessageText(req, null).then(resp => resp.toObject());
  }

  public getCustomDomainClaimedMessageText(req: GetCustomDomainClaimedMessageTextRequest):
    Promise<GetCustomDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomDomainClaimedMessageText(req, null).then(resp => resp.toObject());
  }

  public setDefaultDomainClaimedMessageText(req: SetDefaultDomainClaimedMessageTextRequest):
    Promise<SetDefaultDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultDomainClaimedMessageText(req, null).then(resp => resp.toObject());
  }

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

  public getPrivacyPolicy():
    Promise<GetPrivacyPolicyResponse.AsObject> {
    const req = new GetPrivacyPolicyRequest();
    return this.grpcService.admin.getPrivacyPolicy(req, null).then(resp => resp.toObject());
  }

  public updatePrivacyPolicy(req: UpdatePrivacyPolicyRequest):
    Promise<UpdatePrivacyPolicyResponse.AsObject> {
    return this.grpcService.admin.updatePrivacyPolicy(req, null).then(resp => resp.toObject());
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

  public getLockoutPolicy(): Promise<GetLockoutPolicyResponse.AsObject> {
    const req = new GetLockoutPolicyRequest();
    return this.grpcService.admin.getLockoutPolicy(req, null).then(resp => resp.toObject());
  }

  public updateLockoutPolicy(
    maxAttempts: number,
  ): Promise<UpdateLockoutPolicyResponse.AsObject> {
    const req = new UpdateLockoutPolicyRequest();
    req.setMaxPasswordAttempts(maxAttempts);

    return this.grpcService.admin.updateLockoutPolicy(req, null).then(resp => resp.toObject());
  }

  /* label */

  public getLabelPolicy(): Promise<GetLabelPolicyResponse.AsObject> {
    const req = new GetLabelPolicyRequest();
    return this.grpcService.admin.getLabelPolicy(req, null).then(resp => resp.toObject());
  }

  public updateLabelPolicy(req: UpdateLabelPolicyRequest): Promise<UpdateLabelPolicyResponse.AsObject> {
    return this.grpcService.admin.updateLabelPolicy(req, null).then(resp => resp.toObject());
  }

  public getPreviewLabelPolicy(): Promise<GetPreviewLabelPolicyResponse.AsObject> {
    const req = new GetPreviewLabelPolicyRequest();
    return this.grpcService.admin.getPreviewLabelPolicy(req, null).then(resp => resp.toObject());
  }

  public activateLabelPolicy():
    Promise<ActivateLabelPolicyResponse.AsObject> {
    const req = new ActivateLabelPolicyRequest();
    return this.grpcService.admin.activateLabelPolicy(req, null).then(resp => resp.toObject());
  }

  public removeLabelPolicyFont():
    Promise<RemoveLabelPolicyFontResponse.AsObject> {
    const req = new RemoveLabelPolicyFontRequest();
    return this.grpcService.admin.removeLabelPolicyFont(req, null).then(resp => resp.toObject());
  }

  public removeLabelPolicyIcon():
    Promise<RemoveLabelPolicyIconResponse.AsObject> {
    const req = new RemoveLabelPolicyIconRequest();
    return this.grpcService.admin.removeLabelPolicyIcon(req, null).then(resp => resp.toObject());
  }

  public removeLabelPolicyIconDark():
    Promise<RemoveLabelPolicyIconDarkResponse.AsObject> {
    const req = new RemoveLabelPolicyIconDarkRequest();
    return this.grpcService.admin.removeLabelPolicyIconDark(req, null).then(resp => resp.toObject());
  }

  public removeLabelPolicyLogo():
    Promise<RemoveLabelPolicyLogoResponse.AsObject> {
    const req = new RemoveLabelPolicyLogoRequest();
    return this.grpcService.admin.removeLabelPolicyLogo(req, null).then(resp => resp.toObject());
  }

  public removeLabelPolicyLogoDark():
    Promise<RemoveLabelPolicyLogoDarkResponse.AsObject> {
    const req = new RemoveLabelPolicyLogoDarkRequest();
    return this.grpcService.admin.removeLabelPolicyLogoDark(req, null).then(resp => resp.toObject());
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

import { Injectable } from '@angular/core';
import { SortDirection } from '@angular/material/sort';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject } from 'rxjs';
import { AppQuery } from '../proto/generated/zitadel/app_pb';
import { KeyType } from '../proto/generated/zitadel/auth_n_key_pb';
import { ChangeQuery } from '../proto/generated/zitadel/change_pb';
import { IDPOwnerType } from '../proto/generated/zitadel/idp_pb';
import {
  ActionQuery,
  ActivateCustomLabelPolicyRequest,
  ActivateCustomLabelPolicyResponse,
  AddAPIAppRequest,
  AddAPIAppResponse,
  AddAppKeyRequest,
  AddAppKeyResponse,
  AddAppleProviderRequest,
  AddAppleProviderResponse,
  AddAzureADProviderRequest,
  AddAzureADProviderResponse,
  AddCustomLabelPolicyRequest,
  AddCustomLabelPolicyResponse,
  AddCustomLockoutPolicyRequest,
  AddCustomLockoutPolicyResponse,
  AddCustomLoginPolicyRequest,
  AddCustomLoginPolicyResponse,
  AddCustomNotificationPolicyRequest,
  AddCustomNotificationPolicyResponse,
  AddCustomPasswordAgePolicyRequest,
  AddCustomPasswordAgePolicyResponse,
  AddCustomPasswordComplexityPolicyRequest,
  AddCustomPasswordComplexityPolicyResponse,
  AddCustomPrivacyPolicyRequest,
  AddCustomPrivacyPolicyResponse,
  AddGenericOAuthProviderRequest,
  AddGenericOAuthProviderResponse,
  AddGenericOIDCProviderRequest,
  AddGenericOIDCProviderResponse,
  AddGitHubEnterpriseServerProviderRequest,
  AddGitHubEnterpriseServerProviderResponse,
  AddGitHubProviderRequest,
  AddGitHubProviderResponse,
  AddGitLabProviderRequest,
  AddGitLabProviderResponse,
  AddGitLabSelfHostedProviderRequest,
  AddGitLabSelfHostedProviderResponse,
  AddGoogleProviderRequest,
  AddGoogleProviderResponse,
  AddHumanUserRequest,
  AddHumanUserResponse,
  AddIDPToLoginPolicyRequest,
  AddIDPToLoginPolicyResponse,
  AddJWTProviderRequest,
  AddJWTProviderResponse,
  AddLDAPProviderRequest,
  AddLDAPProviderResponse,
  AddMachineKeyRequest,
  AddMachineKeyResponse,
  AddMachineUserRequest,
  AddMachineUserResponse,
  AddMultiFactorToLoginPolicyRequest,
  AddMultiFactorToLoginPolicyResponse,
  AddOIDCAppRequest,
  AddOIDCAppResponse,
  AddOrgDomainRequest,
  AddOrgDomainResponse,
  AddOrgMemberRequest,
  AddOrgMemberResponse,
  AddOrgRequest,
  AddOrgResponse,
  AddPersonalAccessTokenRequest,
  AddPersonalAccessTokenResponse,
  AddProjectGrantMemberRequest,
  AddProjectGrantMemberResponse,
  AddProjectGrantRequest,
  AddProjectGrantResponse,
  AddProjectMemberRequest,
  AddProjectMemberResponse,
  AddProjectRequest,
  AddProjectResponse,
  AddProjectRoleRequest,
  AddProjectRoleResponse,
  AddSAMLAppRequest,
  AddSAMLAppResponse,
  AddSAMLProviderRequest,
  AddSAMLProviderResponse,
  AddSecondFactorToLoginPolicyRequest,
  AddSecondFactorToLoginPolicyResponse,
  AddUserGrantRequest,
  AddUserGrantResponse,
  BulkAddProjectRolesRequest,
  BulkAddProjectRolesResponse,
  BulkRemoveUserGrantRequest,
  BulkRemoveUserGrantResponse,
  BulkSetUserMetadataRequest,
  BulkSetUserMetadataResponse,
  ClearFlowRequest,
  ClearFlowResponse,
  CreateActionRequest,
  CreateActionResponse,
  DeactivateActionRequest,
  DeactivateActionResponse,
  DeactivateAppRequest,
  DeactivateAppResponse,
  DeactivateOrgIDPRequest,
  DeactivateOrgIDPResponse,
  DeactivateOrgRequest,
  DeactivateOrgResponse,
  DeactivateProjectGrantRequest,
  DeactivateProjectGrantResponse,
  DeactivateProjectRequest,
  DeactivateProjectResponse,
  DeactivateUserRequest,
  DeactivateUserResponse,
  DeleteActionRequest,
  DeleteActionResponse,
  DeleteProviderRequest,
  DeleteProviderResponse,
  GenerateMachineSecretRequest,
  GenerateMachineSecretResponse,
  GenerateOrgDomainValidationRequest,
  GenerateOrgDomainValidationResponse,
  GetActionRequest,
  GetActionResponse,
  GetAppByIDRequest,
  GetAppByIDResponse,
  GetCustomDomainClaimedMessageTextRequest,
  GetCustomDomainClaimedMessageTextResponse,
  GetCustomInitMessageTextRequest,
  GetCustomInitMessageTextResponse,
  GetCustomLoginTextsRequest,
  GetCustomLoginTextsResponse,
  GetCustomPasswordChangeMessageTextRequest,
  GetCustomPasswordChangeMessageTextResponse,
  GetCustomInviteUserMessageTextRequest,
  GetCustomInviteUserMessageTextResponse,
  GetCustomPasswordlessRegistrationMessageTextRequest,
  GetCustomPasswordlessRegistrationMessageTextResponse,
  GetCustomPasswordResetMessageTextRequest,
  GetCustomPasswordResetMessageTextResponse,
  GetCustomVerifyEmailMessageTextRequest,
  GetCustomVerifyEmailMessageTextResponse,
  GetCustomVerifyEmailOTPMessageTextRequest,
  GetCustomVerifyEmailOTPMessageTextResponse,
  GetCustomVerifyPhoneMessageTextRequest,
  GetCustomVerifyPhoneMessageTextResponse,
  GetCustomVerifySMSOTPMessageTextRequest,
  GetCustomVerifySMSOTPMessageTextResponse,
  GetDefaultDomainClaimedMessageTextRequest,
  GetDefaultDomainClaimedMessageTextResponse,
  GetDefaultInitMessageTextRequest,
  GetDefaultInitMessageTextResponse,
  GetDefaultLabelPolicyRequest,
  GetDefaultLabelPolicyResponse,
  GetDefaultLoginTextsRequest,
  GetDefaultLoginTextsResponse,
  GetDefaultPasswordChangeMessageTextRequest,
  GetDefaultPasswordChangeMessageTextResponse,
  GetDefaultInviteUserMessageTextRequest,
  GetDefaultInviteUserMessageTextResponse,
  GetDefaultPasswordComplexityPolicyRequest,
  GetDefaultPasswordComplexityPolicyResponse,
  GetDefaultPasswordlessRegistrationMessageTextRequest,
  GetDefaultPasswordlessRegistrationMessageTextResponse,
  GetDefaultPasswordResetMessageTextRequest,
  GetDefaultPasswordResetMessageTextResponse,
  GetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyEmailMessageTextResponse,
  GetDefaultVerifyEmailOTPMessageTextRequest,
  GetDefaultVerifyEmailOTPMessageTextResponse,
  GetDefaultVerifyPhoneMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextResponse,
  GetDefaultVerifySMSOTPMessageTextRequest,
  GetDefaultVerifySMSOTPMessageTextResponse,
  GetDomainPolicyRequest,
  GetDomainPolicyResponse,
  GetFlowRequest,
  GetFlowResponse,
  GetGrantedProjectByIDRequest,
  GetGrantedProjectByIDResponse,
  GetHumanEmailRequest,
  GetHumanEmailResponse,
  GetHumanPhoneRequest,
  GetHumanPhoneResponse,
  GetHumanProfileRequest,
  GetHumanProfileResponse,
  GetIAMRequest,
  GetIAMResponse,
  GetLabelPolicyRequest,
  GetLabelPolicyResponse,
  GetLockoutPolicyRequest,
  GetLockoutPolicyResponse,
  GetLoginPolicyRequest,
  GetLoginPolicyResponse,
  GetMyOrgRequest,
  GetMyOrgResponse,
  GetNotificationPolicyRequest,
  GetNotificationPolicyResponse,
  GetOIDCInformationRequest,
  GetOIDCInformationResponse,
  GetOrgByDomainGlobalRequest,
  GetOrgByDomainGlobalResponse,
  GetPasswordAgePolicyRequest,
  GetPasswordAgePolicyResponse,
  GetPasswordComplexityPolicyRequest,
  GetPasswordComplexityPolicyResponse,
  GetPreviewLabelPolicyRequest,
  GetPreviewLabelPolicyResponse,
  GetPrivacyPolicyRequest,
  GetPrivacyPolicyResponse,
  GetProjectByIDRequest,
  GetProjectByIDResponse,
  GetProjectGrantByIDRequest,
  GetProjectGrantByIDResponse,
  GetProviderByIDRequest,
  GetProviderByIDResponse,
  GetSupportedLanguagesRequest,
  GetSupportedLanguagesResponse,
  GetUserByIDRequest,
  GetUserByIDResponse,
  GetUserByLoginNameGlobalRequest,
  GetUserByLoginNameGlobalResponse,
  GetUserGrantByIDRequest,
  GetUserGrantByIDResponse,
  GetUserMetadataRequest,
  GetUserMetadataResponse,
  ListActionsRequest,
  ListActionsResponse,
  ListAppChangesRequest,
  ListAppChangesResponse,
  ListAppKeysRequest,
  ListAppKeysResponse,
  ListAppsRequest,
  ListAppsResponse,
  ListFlowTriggerTypesRequest,
  ListFlowTriggerTypesResponse,
  ListFlowTypesRequest,
  ListFlowTypesResponse,
  ListGrantedProjectRolesRequest,
  ListGrantedProjectRolesResponse,
  ListGrantedProjectsRequest,
  ListGrantedProjectsResponse,
  ListHumanAuthFactorsRequest,
  ListHumanAuthFactorsResponse,
  ListHumanLinkedIDPsRequest,
  ListHumanLinkedIDPsResponse,
  ListHumanPasswordlessRequest,
  ListHumanPasswordlessResponse,
  ListLoginPolicyIDPsRequest,
  ListLoginPolicyIDPsResponse,
  ListLoginPolicyMultiFactorsRequest,
  ListLoginPolicyMultiFactorsResponse,
  ListLoginPolicySecondFactorsResponse,
  ListMachineKeysRequest,
  ListMachineKeysResponse,
  ListOrgChangesRequest,
  ListOrgChangesResponse,
  ListOrgDomainsRequest,
  ListOrgDomainsResponse,
  ListOrgMemberRolesRequest,
  ListOrgMemberRolesResponse,
  ListOrgMembersRequest,
  ListOrgMembersResponse,
  ListOrgMetadataRequest,
  ListOrgMetadataResponse,
  ListPersonalAccessTokensRequest,
  ListPersonalAccessTokensResponse,
  ListProjectChangesRequest,
  ListProjectChangesResponse,
  ListProjectGrantChangesRequest,
  ListProjectGrantChangesResponse,
  ListProjectGrantMemberRolesRequest,
  ListProjectGrantMemberRolesResponse,
  ListProjectGrantMembersRequest,
  ListProjectGrantMembersResponse,
  ListProjectGrantsRequest,
  ListProjectGrantsResponse,
  ListProjectMemberRolesRequest,
  ListProjectMemberRolesResponse,
  ListProjectMembersRequest,
  ListProjectMembersResponse,
  ListProjectRolesRequest,
  ListProjectRolesResponse,
  ListProjectsRequest,
  ListProjectsResponse,
  ListProvidersRequest,
  ListProvidersResponse,
  ListUserChangesRequest,
  ListUserChangesResponse,
  ListUserGrantRequest,
  ListUserGrantResponse,
  ListUserMembershipsRequest,
  ListUserMembershipsResponse,
  ListUserMetadataRequest,
  ListUserMetadataResponse,
  ListUsersRequest,
  ListUsersResponse,
  ReactivateActionRequest,
  ReactivateActionResponse,
  ReactivateAppRequest,
  ReactivateAppResponse,
  ReactivateOrgIDPRequest,
  ReactivateOrgIDPResponse,
  ReactivateOrgRequest,
  ReactivateOrgResponse,
  ReactivateProjectGrantRequest,
  ReactivateProjectGrantResponse,
  ReactivateProjectRequest,
  ReactivateProjectResponse,
  ReactivateUserRequest,
  ReactivateUserResponse,
  RegenerateAPIClientSecretRequest,
  RegenerateAPIClientSecretResponse,
  RegenerateOIDCClientSecretRequest,
  RegenerateOIDCClientSecretResponse,
  RemoveAppKeyRequest,
  RemoveAppKeyResponse,
  RemoveAppRequest,
  RemoveAppResponse,
  RemoveCustomLabelPolicyFontRequest,
  RemoveCustomLabelPolicyFontResponse,
  RemoveCustomLabelPolicyIconDarkRequest,
  RemoveCustomLabelPolicyIconDarkResponse,
  RemoveCustomLabelPolicyIconRequest,
  RemoveCustomLabelPolicyIconResponse,
  RemoveCustomLabelPolicyLogoDarkRequest,
  RemoveCustomLabelPolicyLogoDarkResponse,
  RemoveCustomLabelPolicyLogoRequest,
  RemoveCustomLabelPolicyLogoResponse,
  RemoveHumanAuthFactorOTPEmailRequest,
  RemoveHumanAuthFactorOTPEmailResponse,
  RemoveHumanAuthFactorOTPRequest,
  RemoveHumanAuthFactorOTPResponse,
  RemoveHumanAuthFactorOTPSMSRequest,
  RemoveHumanAuthFactorOTPSMSResponse,
  RemoveHumanAuthFactorU2FRequest,
  RemoveHumanAuthFactorU2FResponse,
  RemoveHumanLinkedIDPRequest,
  RemoveHumanLinkedIDPResponse,
  RemoveHumanPasswordlessRequest,
  RemoveHumanPasswordlessResponse,
  RemoveHumanPhoneRequest,
  RemoveHumanPhoneResponse,
  RemoveIDPFromLoginPolicyRequest,
  RemoveIDPFromLoginPolicyResponse,
  RemoveMachineKeyRequest,
  RemoveMachineKeyResponse,
  RemoveMachineSecretRequest,
  RemoveMachineSecretResponse,
  RemoveMultiFactorFromLoginPolicyRequest,
  RemoveMultiFactorFromLoginPolicyResponse,
  RemoveOrgDomainRequest,
  RemoveOrgDomainResponse,
  RemoveOrgMemberRequest,
  RemoveOrgMemberResponse,
  RemoveOrgMetadataRequest,
  RemoveOrgMetadataResponse,
  RemoveOrgRequest,
  RemoveOrgResponse,
  RemovePersonalAccessTokenRequest,
  RemovePersonalAccessTokenResponse,
  RemoveProjectGrantMemberRequest,
  RemoveProjectGrantMemberResponse,
  RemoveProjectGrantRequest,
  RemoveProjectGrantResponse,
  RemoveProjectMemberRequest,
  RemoveProjectMemberResponse,
  RemoveProjectRequest,
  RemoveProjectResponse,
  RemoveProjectRoleRequest,
  RemoveProjectRoleResponse,
  RemoveSecondFactorFromLoginPolicyRequest,
  RemoveSecondFactorFromLoginPolicyResponse,
  RemoveUserGrantRequest,
  RemoveUserGrantResponse,
  RemoveUserMetadataRequest,
  RemoveUserMetadataResponse,
  RemoveUserRequest,
  RemoveUserResponse,
  ResendHumanEmailVerificationRequest,
  ResendHumanInitializationRequest,
  ResendHumanInitializationResponse,
  ResendHumanPhoneVerificationRequest,
  ResetCustomDomainClaimedMessageTextToDefaultRequest,
  ResetCustomDomainClaimedMessageTextToDefaultResponse,
  ResetCustomInitMessageTextToDefaultRequest,
  ResetCustomInitMessageTextToDefaultResponse,
  ResetCustomLoginTextsToDefaultRequest,
  ResetCustomLoginTextsToDefaultResponse,
  ResetCustomPasswordChangeMessageTextToDefaultRequest,
  ResetCustomPasswordChangeMessageTextToDefaultResponse,
  ResetCustomInviteUserMessageTextToDefaultRequest,
  ResetCustomInviteUserMessageTextToDefaultResponse,
  ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest,
  ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse,
  ResetCustomPasswordResetMessageTextToDefaultRequest,
  ResetCustomPasswordResetMessageTextToDefaultResponse,
  ResetCustomVerifyEmailMessageTextToDefaultRequest,
  ResetCustomVerifyEmailMessageTextToDefaultResponse,
  ResetCustomVerifyEmailOTPMessageTextToDefaultRequest,
  ResetCustomVerifyEmailOTPMessageTextToDefaultResponse,
  ResetCustomVerifyPhoneMessageTextToDefaultRequest,
  ResetCustomVerifyPhoneMessageTextToDefaultResponse,
  ResetCustomVerifySMSOTPMessageTextToDefaultRequest,
  ResetCustomVerifySMSOTPMessageTextToDefaultResponse,
  ResetLabelPolicyToDefaultRequest,
  ResetLabelPolicyToDefaultResponse,
  ResetLockoutPolicyToDefaultRequest,
  ResetLockoutPolicyToDefaultResponse,
  ResetLoginPolicyToDefaultRequest,
  ResetLoginPolicyToDefaultResponse,
  ResetNotificationPolicyToDefaultRequest,
  ResetNotificationPolicyToDefaultResponse,
  ResetPasswordAgePolicyToDefaultRequest,
  ResetPasswordAgePolicyToDefaultResponse,
  ResetPasswordComplexityPolicyToDefaultRequest,
  ResetPasswordComplexityPolicyToDefaultResponse,
  ResetPrivacyPolicyToDefaultRequest,
  ResetPrivacyPolicyToDefaultResponse,
  SendHumanResetPasswordNotificationRequest,
  SendPasswordlessRegistrationRequest,
  SendPasswordlessRegistrationResponse,
  SetCustomDomainClaimedMessageTextRequest,
  SetCustomDomainClaimedMessageTextResponse,
  SetCustomInitMessageTextRequest,
  SetCustomInitMessageTextResponse,
  SetCustomLoginTextsRequest,
  SetCustomLoginTextsResponse,
  SetCustomPasswordChangeMessageTextRequest,
  SetCustomPasswordChangeMessageTextResponse,
  SetCustomInviteUserMessageTextRequest,
  SetCustomInviteUserMessageTextResponse,
  SetCustomPasswordlessRegistrationMessageTextRequest,
  SetCustomPasswordlessRegistrationMessageTextResponse,
  SetCustomPasswordResetMessageTextRequest,
  SetCustomPasswordResetMessageTextResponse,
  SetCustomVerifyEmailMessageTextRequest,
  SetCustomVerifyEmailMessageTextResponse,
  SetCustomVerifyEmailOTPMessageTextRequest,
  SetCustomVerifyEmailOTPMessageTextResponse,
  SetCustomVerifyPhoneMessageTextRequest,
  SetCustomVerifyPhoneMessageTextResponse,
  SetCustomVerifySMSOTPMessageTextRequest,
  SetCustomVerifySMSOTPMessageTextResponse,
  SetHumanInitialPasswordRequest,
  SetOrgMetadataRequest,
  SetOrgMetadataResponse,
  SetPrimaryOrgDomainRequest,
  SetPrimaryOrgDomainResponse,
  SetTriggerActionsRequest,
  SetTriggerActionsResponse,
  SetUserMetadataRequest,
  SetUserMetadataResponse,
  UnlockUserRequest,
  UnlockUserResponse,
  UpdateActionRequest,
  UpdateActionResponse,
  UpdateAPIAppConfigRequest,
  UpdateAPIAppConfigResponse,
  UpdateAppleProviderRequest,
  UpdateAppleProviderResponse,
  UpdateAppRequest,
  UpdateAppResponse,
  UpdateAzureADProviderRequest,
  UpdateAzureADProviderResponse,
  UpdateCustomLabelPolicyRequest,
  UpdateCustomLabelPolicyResponse,
  UpdateCustomLockoutPolicyRequest,
  UpdateCustomLockoutPolicyResponse,
  UpdateCustomLoginPolicyRequest,
  UpdateCustomLoginPolicyResponse,
  UpdateCustomNotificationPolicyRequest,
  UpdateCustomNotificationPolicyResponse,
  UpdateCustomPasswordAgePolicyRequest,
  UpdateCustomPasswordAgePolicyResponse,
  UpdateCustomPasswordComplexityPolicyRequest,
  UpdateCustomPasswordComplexityPolicyResponse,
  UpdateCustomPrivacyPolicyRequest,
  UpdateCustomPrivacyPolicyResponse,
  UpdateGenericOAuthProviderRequest,
  UpdateGenericOAuthProviderResponse,
  UpdateGenericOIDCProviderRequest,
  UpdateGenericOIDCProviderResponse,
  UpdateGitHubEnterpriseServerProviderRequest,
  UpdateGitHubEnterpriseServerProviderResponse,
  UpdateGitHubProviderRequest,
  UpdateGitHubProviderResponse,
  UpdateGitLabProviderRequest,
  UpdateGitLabProviderResponse,
  UpdateGitLabSelfHostedProviderRequest,
  UpdateGitLabSelfHostedProviderResponse,
  UpdateGoogleProviderRequest,
  UpdateGoogleProviderResponse,
  UpdateHumanEmailRequest,
  UpdateHumanEmailResponse,
  UpdateHumanPhoneRequest,
  UpdateHumanPhoneResponse,
  UpdateHumanProfileRequest,
  UpdateHumanProfileResponse,
  UpdateJWTProviderRequest,
  UpdateJWTProviderResponse,
  UpdateLDAPProviderRequest,
  UpdateLDAPProviderResponse,
  UpdateMachineRequest,
  UpdateMachineResponse,
  UpdateOIDCAppConfigRequest,
  UpdateOIDCAppConfigResponse,
  UpdateOrgMemberRequest,
  UpdateOrgMemberResponse,
  UpdateOrgRequest,
  UpdateOrgResponse,
  UpdateProjectGrantMemberRequest,
  UpdateProjectGrantMemberResponse,
  UpdateProjectGrantRequest,
  UpdateProjectGrantResponse,
  UpdateProjectMemberRequest,
  UpdateProjectMemberResponse,
  UpdateProjectRequest,
  UpdateProjectResponse,
  UpdateProjectRoleRequest,
  UpdateProjectRoleResponse,
  UpdateSAMLAppConfigRequest,
  UpdateSAMLAppConfigResponse,
  UpdateSAMLProviderRequest,
  UpdateSAMLProviderResponse,
  UpdateUserGrantRequest,
  UpdateUserGrantResponse,
  UpdateUserNameRequest,
  UpdateUserNameResponse,
  ValidateOrgDomainRequest,
  ValidateOrgDomainResponse,
} from '../proto/generated/zitadel/management_pb';
import { SearchQuery } from '../proto/generated/zitadel/member_pb';
import { MetadataQuery } from '../proto/generated/zitadel/metadata_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { DomainSearchQuery, DomainValidationType } from '../proto/generated/zitadel/org_pb';
import { PasswordComplexityPolicy } from '../proto/generated/zitadel/policy_pb';
import { GrantedProject, Project, ProjectQuery, RoleQuery } from '../proto/generated/zitadel/project_pb';
import {
  AccessTokenType,
  Gender,
  MembershipQuery,
  SearchQuery as UserSearchQuery,
  UserFieldName,
  UserGrantQuery,
} from '../proto/generated/zitadel/user_pb';
import { GrpcService } from './grpc.service';

export type ResponseMapper<TResp, TMappedResp> = (resp: TResp) => TMappedResp;

@Injectable({
  providedIn: 'root',
})
export class ManagementService {
  public ownedProjects: BehaviorSubject<Project.AsObject[]> = new BehaviorSubject<Project.AsObject[]>([]);
  public ownedProjectsCount: BehaviorSubject<number> = new BehaviorSubject(0);
  public grantedProjects: BehaviorSubject<GrantedProject.AsObject[]> = new BehaviorSubject<GrantedProject.AsObject[]>([]);
  public grantedProjectsCount: BehaviorSubject<number> = new BehaviorSubject(0);

  constructor(private readonly grpcService: GrpcService) {}

  public getDefaultLoginTexts(req: GetDefaultLoginTextsRequest): Promise<GetDefaultLoginTextsResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultLoginTexts(req, null).then((resp) => resp.toObject());
  }

  public getCustomLoginTexts(req: GetCustomLoginTextsRequest): Promise<GetCustomLoginTextsResponse.AsObject> {
    return this.grpcService.mgmt.getCustomLoginTexts(req, null).then((resp) => resp.toObject());
  }

  public setCustomLoginText(req: SetCustomLoginTextsRequest): Promise<SetCustomLoginTextsResponse.AsObject> {
    return this.grpcService.mgmt.setCustomLoginText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomLoginTextToDefault(lang: string): Promise<ResetCustomLoginTextsToDefaultResponse.AsObject> {
    const req = new ResetCustomLoginTextsToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomLoginTextToDefault(req, null).then((resp) => resp.toObject());
  }

  // message texts

  public getDefaultInitMessageText(
    req: GetDefaultInitMessageTextRequest,
  ): Promise<GetDefaultInitMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultInitMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomInitMessageText(req: GetCustomInitMessageTextRequest): Promise<GetCustomInitMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomInitMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomInitMessageText(req: SetCustomInitMessageTextRequest): Promise<SetCustomInitMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomInitMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomInitMessageTextToDefault(lang: string): Promise<ResetCustomInitMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomInitMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomInitMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultVerifyEmailMessageText(
    req: GetDefaultVerifyEmailMessageTextRequest,
  ): Promise<GetDefaultVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultVerifyEmailMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomVerifyEmailMessageText(
    req: GetCustomVerifyEmailMessageTextRequest,
  ): Promise<GetCustomVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomVerifyEmailMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomVerifyEmailMessageText(
    req: SetCustomVerifyEmailMessageTextRequest,
  ): Promise<SetCustomVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomVerifyEmailMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomVerifyEmailMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomVerifyEmailMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomVerifyEmailMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomVerifyEmailMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultVerifyPhoneMessageText(
    req: GetDefaultVerifyPhoneMessageTextRequest,
  ): Promise<GetDefaultVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultVerifyPhoneMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomVerifyPhoneMessageText(
    req: GetCustomVerifyPhoneMessageTextRequest,
  ): Promise<GetCustomVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomVerifyPhoneMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomVerifyPhoneMessageText(
    req: SetCustomVerifyPhoneMessageTextRequest,
  ): Promise<SetCustomVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomVerifyPhoneMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomVerifyPhoneMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomVerifyPhoneMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomVerifyPhoneMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomVerifyPhoneMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultVerifySMSOTPMessageText(
    req: GetDefaultVerifySMSOTPMessageTextRequest,
  ): Promise<GetDefaultVerifySMSOTPMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultVerifySMSOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomVerifySMSOTPMessageText(
    req: GetCustomVerifySMSOTPMessageTextRequest,
  ): Promise<GetCustomVerifySMSOTPMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomVerifySMSOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomVerifySMSOTPMessageText(
    req: SetCustomVerifySMSOTPMessageTextRequest,
  ): Promise<SetCustomVerifySMSOTPMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomVerifySMSOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomVerifySMSOTPMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomVerifySMSOTPMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomVerifySMSOTPMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomVerifySMSOTPMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultVerifyEmailOTPMessageText(
    req: GetDefaultVerifyEmailOTPMessageTextRequest,
  ): Promise<GetDefaultVerifyEmailOTPMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultVerifyEmailOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomVerifyEmailOTPMessageText(
    req: GetCustomVerifyEmailOTPMessageTextRequest,
  ): Promise<GetCustomVerifyEmailOTPMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomVerifyEmailOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomVerifyEmailOTPMessageText(
    req: SetCustomVerifyEmailOTPMessageTextRequest,
  ): Promise<SetCustomVerifyEmailOTPMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomVerifyEmailOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomVerifyEmailOTPMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomVerifyEmailOTPMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomVerifyEmailOTPMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomVerifyEmailOTPMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultPasswordResetMessageText(
    req: GetDefaultPasswordResetMessageTextRequest,
  ): Promise<GetDefaultPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultPasswordResetMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomPasswordResetMessageText(
    req: GetCustomPasswordResetMessageTextRequest,
  ): Promise<GetCustomPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomPasswordResetMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomPasswordResetMessageText(
    req: SetCustomPasswordResetMessageTextRequest,
  ): Promise<SetCustomPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomPasswordResetMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomPasswordResetMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomPasswordResetMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomPasswordResetMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomPasswordResetMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultDomainClaimedMessageText(
    req: GetDefaultDomainClaimedMessageTextRequest,
  ): Promise<GetDefaultDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultDomainClaimedMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomDomainClaimedMessageText(
    req: GetCustomDomainClaimedMessageTextRequest,
  ): Promise<GetCustomDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomDomainClaimedMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomDomainClaimedMessageCustomText(
    req: SetCustomDomainClaimedMessageTextRequest,
  ): Promise<SetCustomDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomDomainClaimedMessageCustomText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomDomainClaimedMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomDomainClaimedMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomDomainClaimedMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomDomainClaimedMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultPasswordlessRegistrationMessageText(
    req: GetDefaultPasswordlessRegistrationMessageTextRequest,
  ): Promise<GetDefaultPasswordlessRegistrationMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultPasswordlessRegistrationMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomPasswordlessRegistrationMessageText(
    req: GetCustomPasswordlessRegistrationMessageTextRequest,
  ): Promise<GetCustomPasswordlessRegistrationMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomPasswordlessRegistrationMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomPasswordlessRegistrationMessageCustomText(
    req: SetCustomPasswordlessRegistrationMessageTextRequest,
  ): Promise<SetCustomPasswordlessRegistrationMessageTextResponse.AsObject> {
    return this.grpcService.mgmt
      .setCustomPasswordlessRegistrationMessageCustomText(req, null)
      .then((resp) => resp.toObject());
  }

  public resetCustomPasswordlessRegistrationMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt
      .resetCustomPasswordlessRegistrationMessageTextToDefault(req, null)
      .then((resp) => resp.toObject());
  }

  public getDefaultPasswordChangeMessageText(
    req: GetDefaultPasswordChangeMessageTextRequest,
  ): Promise<GetDefaultPasswordChangeMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultPasswordChangeMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomPasswordChangeMessageText(
    req: GetCustomPasswordChangeMessageTextRequest,
  ): Promise<GetCustomPasswordChangeMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomPasswordChangeMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomPasswordChangeMessageText(
    req: SetCustomPasswordChangeMessageTextRequest,
  ): Promise<SetCustomPasswordChangeMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomPasswordChangeMessageCustomText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomPasswordChangeMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomPasswordChangeMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomPasswordChangeMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomPasswordChangeMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultInviteUserMessageText(
    req: GetDefaultInviteUserMessageTextRequest,
  ): Promise<GetDefaultInviteUserMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultInviteUserMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomInviteUserMessageText(
    req: GetCustomInviteUserMessageTextRequest,
  ): Promise<GetCustomInviteUserMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.getCustomInviteUserMessageText(req, null).then((resp) => resp.toObject());
  }

  public setCustomInviteUserMessageText(
    req: SetCustomInviteUserMessageTextRequest,
  ): Promise<SetCustomInviteUserMessageTextResponse.AsObject> {
    return this.grpcService.mgmt.setCustomInviteUserMessageCustomText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomInviteUserMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomInviteUserMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomInviteUserMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.mgmt.resetCustomInviteUserMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public updateUserName(userId: string, username: string): Promise<UpdateUserNameResponse.AsObject> {
    const req = new UpdateUserNameRequest();
    req.setUserId(userId);
    req.setUserName(username);
    return this.grpcService.mgmt.updateUserName(req, null).then((resp) => resp.toObject());
  }

  public unlockUser(req: UnlockUserRequest): Promise<UnlockUserResponse.AsObject> {
    return this.grpcService.mgmt.unlockUser(req, null).then((resp) => resp.toObject());
  }

  public generateMachineSecret(userId: string): Promise<GenerateMachineSecretResponse.AsObject> {
    const req = new GenerateMachineSecretRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.generateMachineSecret(req, null).then((resp) => resp.toObject());
  }

  public removeMachineSecret(userId: string): Promise<RemoveMachineSecretResponse.AsObject> {
    const req = new RemoveMachineSecretRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.removeMachineSecret(req, null).then((resp) => resp.toObject());
  }

  public getPrivacyPolicy(): Promise<GetPrivacyPolicyResponse.AsObject> {
    const req = new GetPrivacyPolicyRequest();
    return this.grpcService.mgmt.getPrivacyPolicy(req, null).then((resp) => resp.toObject());
  }

  public addCustomPrivacyPolicy(req: AddCustomPrivacyPolicyRequest): Promise<AddCustomPrivacyPolicyResponse.AsObject> {
    return this.grpcService.mgmt.addCustomPrivacyPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateCustomPrivacyPolicy(
    req: UpdateCustomPrivacyPolicyRequest,
  ): Promise<UpdateCustomPrivacyPolicyResponse.AsObject> {
    return this.grpcService.mgmt.updateCustomPrivacyPolicy(req, null).then((resp) => resp.toObject());
  }

  public resetPrivacyPolicyToDefault(): Promise<ResetPrivacyPolicyToDefaultResponse.AsObject> {
    const req = new ResetPrivacyPolicyToDefaultRequest();
    return this.grpcService.mgmt.resetPrivacyPolicyToDefault(req, null).then((resp) => resp.toObject());
  }

  public listHumanPasswordless(userId: string): Promise<ListHumanPasswordlessResponse.AsObject> {
    const req = new ListHumanPasswordlessRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.listHumanPasswordless(req, null).then((resp) => resp.toObject());
  }

  public removeHumanPasswordless(tokenId: string, userId: string): Promise<RemoveHumanPasswordlessResponse.AsObject> {
    const req = new RemoveHumanPasswordlessRequest();
    req.setTokenId(tokenId);
    req.setUserId(userId);
    return this.grpcService.mgmt.removeHumanPasswordless(req, null).then((resp) => resp.toObject());
  }

  public sendPasswordlessRegistration(userId: string): Promise<SendPasswordlessRegistrationResponse.AsObject> {
    const req = new SendPasswordlessRegistrationRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.sendPasswordlessRegistration(req, null).then((resp) => resp.toObject());
  }

  public listLoginPolicyMultiFactors(): Promise<ListLoginPolicyMultiFactorsResponse.AsObject> {
    const req = new ListLoginPolicyMultiFactorsRequest();
    return this.grpcService.mgmt.listLoginPolicyMultiFactors(req, null).then((resp) => resp.toObject());
  }

  public addMultiFactorToLoginPolicy(
    req: AddMultiFactorToLoginPolicyRequest,
  ): Promise<AddMultiFactorToLoginPolicyResponse.AsObject> {
    return this.grpcService.mgmt.addMultiFactorToLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public removeMultiFactorFromLoginPolicy(
    req: RemoveMultiFactorFromLoginPolicyRequest,
  ): Promise<RemoveMultiFactorFromLoginPolicyResponse.AsObject> {
    return this.grpcService.mgmt.removeMultiFactorFromLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public listLoginPolicySecondFactors(): Promise<ListLoginPolicySecondFactorsResponse.AsObject> {
    const req = new Empty();
    return this.grpcService.mgmt.listLoginPolicySecondFactors(req, null).then((resp) => resp.toObject());
  }

  public addSecondFactorToLoginPolicy(
    req: AddSecondFactorToLoginPolicyRequest,
  ): Promise<AddSecondFactorToLoginPolicyResponse.AsObject> {
    return this.grpcService.mgmt.addSecondFactorToLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public removeSecondFactorFromLoginPolicy(
    req: RemoveSecondFactorFromLoginPolicyRequest,
  ): Promise<RemoveSecondFactorFromLoginPolicyResponse.AsObject> {
    return this.grpcService.mgmt.removeSecondFactorFromLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public getLoginPolicy(): Promise<GetLoginPolicyResponse.AsObject> {
    const req = new GetLoginPolicyRequest();
    return this.grpcService.mgmt.getLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateCustomLoginPolicy(req: UpdateCustomLoginPolicyRequest): Promise<UpdateCustomLoginPolicyResponse.AsObject> {
    return this.grpcService.mgmt.updateCustomLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public addCustomLoginPolicy(req: AddCustomLoginPolicyRequest): Promise<AddCustomLoginPolicyResponse.AsObject> {
    return this.grpcService.mgmt.addCustomLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public resetLoginPolicyToDefault(): Promise<ResetLoginPolicyToDefaultResponse.AsObject> {
    const req = new ResetLoginPolicyToDefaultRequest();
    return this.grpcService.mgmt.resetLoginPolicyToDefault(req, null).then((resp) => resp.toObject());
  }

  public addIDPToLoginPolicy(idpId: string, ownerType: IDPOwnerType): Promise<AddIDPToLoginPolicyResponse.AsObject> {
    const req = new AddIDPToLoginPolicyRequest();
    req.setIdpId(idpId);
    req.setOwnertype(ownerType);
    return this.grpcService.mgmt.addIDPToLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public removeIDPFromLoginPolicy(idpId: string): Promise<RemoveIDPFromLoginPolicyResponse.AsObject> {
    const req = new RemoveIDPFromLoginPolicyRequest();
    req.setIdpId(idpId);
    return this.grpcService.mgmt.removeIDPFromLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public listLoginPolicyIDPs(limit?: number, offset?: number): Promise<ListLoginPolicyIDPsResponse.AsObject> {
    const req = new ListLoginPolicyIDPsRequest();
    const metadata = new ListQuery();
    if (limit) {
      metadata.setLimit(limit);
    }
    if (offset) {
      metadata.setOffset(offset);
    }
    return this.grpcService.mgmt.listLoginPolicyIDPs(req, null).then((resp) => resp.toObject());
  }

  public deactivateOrgIDP(idpId: string): Promise<DeactivateOrgIDPResponse.AsObject> {
    const req = new DeactivateOrgIDPRequest();
    req.setIdpId(idpId);
    return this.grpcService.mgmt.deactivateOrgIDP(req, null).then((resp) => resp.toObject());
  }

  public reactivateOrgIDP(idpId: string): Promise<ReactivateOrgIDPResponse.AsObject> {
    const req = new ReactivateOrgIDPRequest();
    req.setIdpId(idpId);
    return this.grpcService.mgmt.reactivateOrgIDP(req, null).then((resp) => resp.toObject());
  }

  //   idp templates

  public addAzureADProvider(req: AddAzureADProviderRequest): Promise<AddAzureADProviderResponse.AsObject> {
    return this.grpcService.mgmt.addAzureADProvider(req, null).then((resp) => resp.toObject());
  }

  public updateAzureADProvider(req: UpdateAzureADProviderRequest): Promise<UpdateAzureADProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateAzureADProvider(req, null).then((resp) => resp.toObject());
  }

  public addGoogleProvider(req: AddGoogleProviderRequest): Promise<AddGoogleProviderResponse.AsObject> {
    return this.grpcService.mgmt.addGoogleProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGoogleProvider(req: UpdateGoogleProviderRequest): Promise<UpdateGoogleProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateGoogleProvider(req, null).then((resp) => resp.toObject());
  }

  public addLDAPProvider(req: AddLDAPProviderRequest): Promise<AddLDAPProviderResponse.AsObject> {
    return this.grpcService.mgmt.addLDAPProvider(req, null).then((resp) => resp.toObject());
  }

  public updateLDAPProvider(req: UpdateLDAPProviderRequest): Promise<UpdateLDAPProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateLDAPProvider(req, null).then((resp) => resp.toObject());
  }

  public addGitLabProvider(req: AddGitLabProviderRequest): Promise<AddGitLabProviderResponse.AsObject> {
    return this.grpcService.mgmt.addGitLabProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGitLabProvider(req: UpdateGitLabProviderRequest): Promise<UpdateGitLabProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateGitLabProvider(req, null).then((resp) => resp.toObject());
  }

  public addGitLabSelfHostedProvider(
    req: AddGitLabSelfHostedProviderRequest,
  ): Promise<AddGitLabSelfHostedProviderResponse.AsObject> {
    return this.grpcService.mgmt.addGitLabSelfHostedProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGitLabSelfHostedProvider(
    req: UpdateGitLabSelfHostedProviderRequest,
  ): Promise<UpdateGitLabSelfHostedProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateGitLabSelfHostedProvider(req, null).then((resp) => resp.toObject());
  }

  public addGitHubProvider(req: AddGitHubProviderRequest): Promise<AddGitHubProviderResponse.AsObject> {
    return this.grpcService.mgmt.addGitHubProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGitHubProvider(req: UpdateGitHubProviderRequest): Promise<UpdateGitHubProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateGitHubProvider(req, null).then((resp) => resp.toObject());
  }

  public addGenericOIDCProvider(req: AddGenericOIDCProviderRequest): Promise<AddGenericOIDCProviderResponse.AsObject> {
    return this.grpcService.mgmt.addGenericOIDCProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGenericOIDCProvider(
    req: UpdateGenericOIDCProviderRequest,
  ): Promise<UpdateGenericOIDCProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateGenericOIDCProvider(req, null).then((resp) => resp.toObject());
  }

  public addGenericOAuthProvider(req: AddGenericOAuthProviderRequest): Promise<AddGenericOAuthProviderResponse.AsObject> {
    return this.grpcService.mgmt.addGenericOAuthProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGenericOAuthProvider(
    req: UpdateGenericOAuthProviderRequest,
  ): Promise<UpdateGenericOAuthProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateGenericOAuthProvider(req, null).then((resp) => resp.toObject());
  }

  public addJWTProvider(req: AddJWTProviderRequest): Promise<AddJWTProviderResponse.AsObject> {
    return this.grpcService.mgmt.addJWTProvider(req, null).then((resp) => resp.toObject());
  }

  public updateJWTProvider(req: UpdateJWTProviderRequest): Promise<UpdateJWTProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateJWTProvider(req, null).then((resp) => resp.toObject());
  }

  public addSAMLProvider(req: AddSAMLProviderRequest): Promise<AddSAMLProviderResponse.AsObject> {
    return this.grpcService.mgmt.addSAMLProvider(req, null).then((resp) => resp.toObject());
  }

  public updateSAMLProvider(req: UpdateSAMLProviderRequest): Promise<UpdateSAMLProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateSAMLProvider(req, null).then((resp) => resp.toObject());
  }

  public addGitHubEnterpriseServerProvider(
    req: AddGitHubEnterpriseServerProviderRequest,
  ): Promise<AddGitHubEnterpriseServerProviderResponse.AsObject> {
    return this.grpcService.mgmt.addGitHubEnterpriseServerProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGitHubEnterpriseServerProvider(
    req: UpdateGitHubEnterpriseServerProviderRequest,
  ): Promise<UpdateGitHubEnterpriseServerProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateGitHubEnterpriseServerProvider(req, null).then((resp) => resp.toObject());
  }

  public addAppleProvider(req: AddAppleProviderRequest): Promise<AddAppleProviderResponse.AsObject> {
    return this.grpcService.mgmt.addAppleProvider(req, null).then((resp) => resp.toObject());
  }

  public updateAppleProvider(req: UpdateAppleProviderRequest): Promise<UpdateAppleProviderResponse.AsObject> {
    return this.grpcService.mgmt.updateAppleProvider(req, null).then((resp) => resp.toObject());
  }

  public deleteProvider(id: string): Promise<DeleteProviderResponse.AsObject> {
    const req = new DeleteProviderRequest();
    req.setId(id);
    return this.grpcService.mgmt.deleteProvider(req, null).then((resp) => resp.toObject());
  }

  public listProviders(req: ListProvidersRequest): Promise<ListProvidersResponse.AsObject> {
    return this.grpcService.mgmt.listProviders(req, null).then((resp) => resp.toObject());
  }

  public getProviderByID(req: GetProviderByIDRequest): Promise<GetProviderByIDResponse.AsObject> {
    return this.grpcService.mgmt.getProviderByID(req, null).then((resp) => resp.toObject());
  }

  public getProviderID(req: GetProviderByIDRequest): Promise<GetProviderByIDResponse> {
    return this.grpcService.mgmt.getProviderByID(req, null);
  }

  public addHumanUser(req: AddHumanUserRequest): Promise<AddHumanUserResponse.AsObject> {
    return this.grpcService.mgmt.addHumanUser(req, null).then((resp) => resp.toObject());
  }

  public addMachineUser(req: AddMachineUserRequest): Promise<AddMachineUserResponse.AsObject> {
    return this.grpcService.mgmt.addMachineUser(req, null).then((resp) => resp.toObject());
  }

  public updateMachine(
    userId: string,
    name?: string,
    description?: string,
    accessTokenType?: AccessTokenType,
  ): Promise<UpdateMachineResponse.AsObject> {
    const req = new UpdateMachineRequest();
    req.setUserId(userId);
    if (name) {
      req.setName(name);
    }
    if (description) {
      req.setDescription(description);
    }
    if (accessTokenType !== undefined) {
      req.setAccessTokenType(accessTokenType);
    }
    return this.grpcService.mgmt.updateMachine(req, null).then((resp) => resp.toObject());
  }

  public addMachineKey(userId: string, type: KeyType, date?: Timestamp): Promise<AddMachineKeyResponse.AsObject> {
    const req = new AddMachineKeyRequest();
    req.setType(type);
    req.setUserId(userId);
    if (date) {
      req.setExpirationDate(date);
    }
    return this.grpcService.mgmt.addMachineKey(req, null).then((resp) => resp.toObject());
  }

  public removeMachineKey(keyId: string, userId: string): Promise<RemoveMachineKeyResponse.AsObject> {
    const req = new RemoveMachineKeyRequest();
    req.setKeyId(keyId);
    req.setUserId(userId);

    return this.grpcService.mgmt.removeMachineKey(req, null).then((resp) => resp.toObject());
  }

  public listMachineKeys(
    userId: string,
    limit?: number,
    offset?: number,
    asc?: boolean,
  ): Promise<ListMachineKeysResponse.AsObject> {
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
    req.setQuery(metadata);
    return this.grpcService.mgmt.listMachineKeys(req, null).then((resp) => resp.toObject());
  }

  public removeHumanLinkedIDP(
    idpId: string,
    linkedUserId: string,
    userId: string,
  ): Promise<RemoveHumanLinkedIDPResponse.AsObject> {
    const req = new RemoveHumanLinkedIDPRequest();
    req.setIdpId(idpId);
    req.setUserId(userId);
    req.setLinkedUserId(linkedUserId);
    return this.grpcService.mgmt.removeHumanLinkedIDP(req, null).then((resp) => resp.toObject());
  }

  public listHumanLinkedIDPs(
    userId: string,
    limit?: number,
    offset?: number,
  ): Promise<ListHumanLinkedIDPsResponse.AsObject> {
    const req = new ListHumanLinkedIDPsRequest();
    const metadata = new ListQuery();
    req.setUserId(userId);
    if (limit) {
      metadata.setLimit(limit);
    }
    if (offset) {
      metadata.setOffset(offset);
    }
    req.setQuery(metadata);
    return this.grpcService.mgmt.listHumanLinkedIDPs(req, null).then((resp) => resp.toObject());
  }

  public getAction(id: string): Promise<GetActionResponse.AsObject> {
    const req = new GetActionRequest();
    req.setId(id);
    return this.grpcService.mgmt.getAction(req, null).then((resp) => resp.toObject());
  }

  public createAction(req: CreateActionRequest): Promise<CreateActionResponse.AsObject> {
    return this.grpcService.mgmt.createAction(req, null).then((resp) => resp.toObject());
  }

  public updateAction(req: UpdateActionRequest): Promise<UpdateActionResponse.AsObject> {
    return this.grpcService.mgmt.updateAction(req, null).then((resp) => resp.toObject());
  }

  public deleteAction(id: string): Promise<DeleteActionResponse.AsObject> {
    const req = new DeleteActionRequest();
    req.setId(id);
    return this.grpcService.mgmt.deleteAction(req, null).then((resp) => resp.toObject());
  }

  public deactivateAction(id: string): Promise<DeactivateActionResponse.AsObject> {
    const req = new DeactivateActionRequest();
    req.setId(id);
    return this.grpcService.mgmt.deactivateAction(req, null).then((resp) => resp.toObject());
  }

  public reactivateAction(id: string): Promise<ReactivateActionResponse.AsObject> {
    const req = new ReactivateActionRequest();
    req.setId(id);
    return this.grpcService.mgmt.reactivateAction(req, null).then((resp) => resp.toObject());
  }

  public listActions(
    limit?: number,
    offset?: number,
    asc?: boolean,
    queryList?: ActionQuery[],
  ): Promise<ListActionsResponse.AsObject> {
    const req = new ListActionsRequest();
    const metadata = new ListQuery();
    if (queryList) {
      req.setQueriesList(queryList);
    }

    if (limit) {
      metadata.setLimit(limit);
    }
    if (offset) {
      metadata.setOffset(offset);
    }
    if (asc) {
      metadata.setAsc(asc);
    }
    req.setQuery(metadata);
    return this.grpcService.mgmt.listActions(req, null).then((resp) => resp.toObject());
  }

  public listFlowTypes(): Promise<ListFlowTypesResponse.AsObject> {
    const req = new ListFlowTypesRequest();
    return this.grpcService.mgmt.listFlowTypes(req, null).then((resp) => resp.toObject());
  }

  public listFlowTriggerTypes(type: string): Promise<ListFlowTriggerTypesResponse.AsObject> {
    const req = new ListFlowTriggerTypesRequest();
    req.setType(type);
    return this.grpcService.mgmt.listFlowTriggerTypes(req, null).then((resp) => resp.toObject());
  }

  public getFlow(type: string): Promise<GetFlowResponse.AsObject> {
    const req = new GetFlowRequest();
    req.setType(type);
    return this.grpcService.mgmt.getFlow(req, null).then((resp) => resp.toObject());
  }

  public clearFlow(type: string): Promise<ClearFlowResponse.AsObject> {
    const req = new ClearFlowRequest();
    req.setType(type);
    return this.grpcService.mgmt.clearFlow(req, null).then((resp) => resp.toObject());
  }

  public setTriggerActions(
    actionIdsList: string[],
    type: string,
    triggerType: string,
  ): Promise<SetTriggerActionsResponse.AsObject> {
    const req = new SetTriggerActionsRequest();
    req.setActionIdsList(actionIdsList);
    req.setFlowType(type);
    req.setTriggerType(triggerType);
    return this.grpcService.mgmt.setTriggerActions(req, null).then((resp) => resp.toObject());
  }

  public addPersonalAccessToken(userId: string, date?: Timestamp): Promise<AddPersonalAccessTokenResponse.AsObject> {
    const req = new AddPersonalAccessTokenRequest();
    req.setUserId(userId);
    if (date) {
      req.setExpirationDate(date);
    }
    return this.grpcService.mgmt.addPersonalAccessToken(req, null).then((resp) => resp.toObject());
  }

  public removePersonalAccessToken(tokenId: string, userId: string): Promise<RemovePersonalAccessTokenResponse.AsObject> {
    const req = new RemovePersonalAccessTokenRequest();
    req.setTokenId(tokenId);
    req.setUserId(userId);
    return this.grpcService.mgmt.removePersonalAccessToken(req, null).then((resp) => resp.toObject());
  }

  public listPersonalAccessTokens(
    userId: string,
    limit?: number,
    offset?: number,
    asc?: boolean,
  ): Promise<ListPersonalAccessTokensResponse.AsObject> {
    const req = new ListPersonalAccessTokensRequest();
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
    req.setQuery(metadata);
    return this.grpcService.mgmt.listPersonalAccessTokens(req, null).then((resp) => resp.toObject());
  }

  public getIAM(): Promise<GetIAMResponse.AsObject> {
    const req = new GetIAMRequest();
    return this.grpcService.mgmt.getIAM(req, null).then((resp) => resp.toObject());
  }

  public getDefaultPasswordComplexityPolicy(): Promise<GetDefaultPasswordComplexityPolicyResponse.AsObject> {
    const req = new GetDefaultPasswordComplexityPolicyRequest();
    return this.grpcService.mgmt.getDefaultPasswordComplexityPolicy(req, null).then((resp) => resp.toObject());
  }

  public getMyOrg(): Promise<GetMyOrgResponse.AsObject> {
    const req = new GetMyOrgRequest();
    return this.grpcService.mgmt.getMyOrg(req, null).then((resp) => resp.toObject());
  }

  public addOrgDomain(domain: string): Promise<AddOrgDomainResponse.AsObject> {
    const req = new AddOrgDomainRequest();
    req.setDomain(domain);
    return this.grpcService.mgmt.addOrgDomain(req, null).then((resp) => resp.toObject());
  }

  public removeOrgDomain(domain: string): Promise<RemoveOrgDomainResponse.AsObject> {
    const req = new RemoveOrgDomainRequest();
    req.setDomain(domain);
    return this.grpcService.mgmt.removeOrgDomain(req, null).then((resp) => resp.toObject());
  }

  public listOrgDomains(queryList?: DomainSearchQuery[]): Promise<ListOrgDomainsResponse.AsObject> {
    const req: ListOrgDomainsRequest = new ListOrgDomainsRequest();
    // const metadata= new ListQuery();
    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listOrgDomains(req, null).then((resp) => resp.toObject());
  }

  public setPrimaryOrgDomain(domain: string): Promise<SetPrimaryOrgDomainResponse.AsObject> {
    const req = new SetPrimaryOrgDomainRequest();
    req.setDomain(domain);
    return this.grpcService.mgmt.setPrimaryOrgDomain(req, null).then((resp) => resp.toObject());
  }

  public generateOrgDomainValidation(
    domain: string,
    type: DomainValidationType,
  ): Promise<GenerateOrgDomainValidationResponse.AsObject> {
    const req: GenerateOrgDomainValidationRequest = new GenerateOrgDomainValidationRequest();
    req.setDomain(domain);
    req.setType(type);

    return this.grpcService.mgmt.generateOrgDomainValidation(req, null).then((resp) => resp.toObject());
  }

  public validateOrgDomain(domain: string): Promise<ValidateOrgDomainResponse.AsObject> {
    const req = new ValidateOrgDomainRequest();
    req.setDomain(domain);

    return this.grpcService.mgmt.validateOrgDomain(req, null).then((resp) => resp.toObject());
  }

  public listOrgMembers(limit: number, offset: number): Promise<ListOrgMembersResponse.AsObject> {
    const req = new ListOrgMembersRequest();
    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }
    req.setQuery(query);

    return this.grpcService.mgmt.listOrgMembers(req, null).then((resp) => resp.toObject());
  }

  public getOrgByDomainGlobal(domain: string): Promise<GetOrgByDomainGlobalResponse.AsObject> {
    const req = new GetOrgByDomainGlobalRequest();
    req.setDomain(domain);
    return this.grpcService.mgmt.getOrgByDomainGlobal(req, null).then((resp) => resp.toObject());
  }

  public addOrg(name: string): Promise<AddOrgResponse.AsObject> {
    const req = new AddOrgRequest();
    req.setName(name);
    return this.grpcService.mgmt.addOrg(req, null).then((resp) => resp.toObject());
  }

  public addOrgMember(userId: string, rolesList: string[]): Promise<AddOrgMemberResponse.AsObject> {
    const req = new AddOrgMemberRequest();
    req.setUserId(userId);
    if (rolesList) {
      req.setRolesList(rolesList);
    }
    return this.grpcService.mgmt.addOrgMember(req, null).then((resp) => resp.toObject());
  }

  public updateOrgMember(userId: string, rolesList: string[]): Promise<UpdateOrgMemberResponse.AsObject> {
    const req = new UpdateOrgMemberRequest();
    req.setUserId(userId);
    req.setRolesList(rolesList);
    return this.grpcService.mgmt.updateOrgMember(req, null).then((resp) => resp.toObject());
  }

  public removeOrgMember(userId: string): Promise<RemoveOrgMemberResponse.AsObject> {
    const req = new RemoveOrgMemberRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.removeOrgMember(req, null).then((resp) => resp.toObject());
  }

  public deactivateOrg(): Promise<DeactivateOrgResponse.AsObject> {
    const req = new DeactivateOrgRequest();
    return this.grpcService.mgmt.deactivateOrg(req, null).then((resp) => resp.toObject());
  }

  public reactivateOrg(): Promise<ReactivateOrgResponse.AsObject> {
    const req = new ReactivateOrgRequest();
    return this.grpcService.mgmt.reactivateOrg(req, null).then((resp) => resp.toObject());
  }

  public addProjectGrant(
    orgId: string,
    projectId: string,
    roleKeysList: string[],
  ): Promise<AddProjectGrantResponse.AsObject> {
    const req = new AddProjectGrantRequest();
    req.setProjectId(projectId);
    req.setGrantedOrgId(orgId);
    req.setRoleKeysList(roleKeysList);
    return this.grpcService.mgmt.addProjectGrant(req, null).then((resp) => resp.toObject());
  }

  public listOrgMemberRoles(): Promise<ListOrgMemberRolesResponse.AsObject> {
    const req = new ListOrgMemberRolesRequest();
    return this.grpcService.mgmt.listOrgMemberRoles(req, null).then((resp) => resp.toObject());
  }

  // Policy

  public getLabelPolicy(): Promise<GetLabelPolicyResponse.AsObject> {
    const req = new GetLabelPolicyRequest();
    return this.grpcService.mgmt.getLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateCustomLabelPolicy(req: UpdateCustomLabelPolicyRequest): Promise<UpdateCustomLabelPolicyResponse.AsObject> {
    return this.grpcService.mgmt.updateCustomLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public resetLabelPolicyToDefault(): Promise<ResetLabelPolicyToDefaultResponse.AsObject> {
    const req = new ResetLabelPolicyToDefaultRequest();
    return this.grpcService.mgmt.resetLabelPolicyToDefault(req, null).then((resp) => resp.toObject());
  }

  public addCustomLabelPolicy(req: AddCustomLabelPolicyRequest): Promise<AddCustomLabelPolicyResponse.AsObject> {
    return this.grpcService.mgmt.addCustomLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public getDefaultLabelPolicy(req: GetDefaultLabelPolicyRequest): Promise<GetDefaultLabelPolicyResponse.AsObject> {
    return this.grpcService.mgmt.getDefaultLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public getPreviewLabelPolicy(): Promise<GetPreviewLabelPolicyResponse.AsObject> {
    const req = new GetPreviewLabelPolicyRequest();
    return this.grpcService.mgmt.getPreviewLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public activateCustomLabelPolicy(): Promise<ActivateCustomLabelPolicyResponse.AsObject> {
    const req = new ActivateCustomLabelPolicyRequest();
    return this.grpcService.mgmt.activateCustomLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyFont(): Promise<RemoveCustomLabelPolicyFontResponse.AsObject> {
    const req = new RemoveCustomLabelPolicyFontRequest();
    return this.grpcService.mgmt.removeCustomLabelPolicyFont(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyIcon(): Promise<RemoveCustomLabelPolicyIconResponse.AsObject> {
    const req = new RemoveCustomLabelPolicyIconRequest();
    return this.grpcService.mgmt.removeCustomLabelPolicyIcon(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyIconDark(): Promise<RemoveCustomLabelPolicyIconDarkResponse.AsObject> {
    const req = new RemoveCustomLabelPolicyIconDarkRequest();
    return this.grpcService.mgmt.removeCustomLabelPolicyIconDark(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyLogo(): Promise<RemoveCustomLabelPolicyLogoResponse.AsObject> {
    const req = new RemoveCustomLabelPolicyLogoRequest();
    return this.grpcService.mgmt.removeCustomLabelPolicyLogo(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyLogoDark(): Promise<RemoveCustomLabelPolicyLogoDarkResponse.AsObject> {
    const req = new RemoveCustomLabelPolicyLogoDarkRequest();
    return this.grpcService.mgmt.removeCustomLabelPolicyLogoDark(req, null).then((resp) => resp.toObject());
  }

  public getDomainPolicy(): Promise<GetDomainPolicyResponse.AsObject> {
    const req = new GetDomainPolicyRequest();
    return this.grpcService.mgmt.getDomainPolicy(req, null).then((resp) => resp.toObject());
  }

  public getPasswordAgePolicy(): Promise<GetPasswordAgePolicyResponse.AsObject> {
    const req = new GetPasswordAgePolicyRequest();
    return this.grpcService.mgmt.getPasswordAgePolicy(req, null).then((resp) => resp.toObject());
  }

  public addCustomPasswordAgePolicy(
    maxAgeDays: number,
    expireWarnDays: number,
  ): Promise<AddCustomPasswordAgePolicyResponse.AsObject> {
    const req = new AddCustomPasswordAgePolicyRequest();
    req.setMaxAgeDays(maxAgeDays);
    req.setExpireWarnDays(expireWarnDays);

    return this.grpcService.mgmt.addCustomPasswordAgePolicy(req, null).then((resp) => resp.toObject());
  }

  public resetPasswordAgePolicyToDefault(): Promise<ResetPasswordAgePolicyToDefaultResponse.AsObject> {
    const req = new ResetPasswordAgePolicyToDefaultRequest();
    return this.grpcService.mgmt.resetPasswordAgePolicyToDefault(req, null).then((resp) => resp.toObject());
  }

  public updateCustomPasswordAgePolicy(
    maxAgeDays: number,
    expireWarnDays: number,
  ): Promise<UpdateCustomPasswordAgePolicyResponse.AsObject> {
    const req = new UpdateCustomPasswordAgePolicyRequest();
    req.setMaxAgeDays(maxAgeDays);
    req.setExpireWarnDays(expireWarnDays);
    return this.grpcService.mgmt.updateCustomPasswordAgePolicy(req, null).then((resp) => resp.toObject());
  }

  public getPasswordComplexityPolicy(): Promise<GetPasswordComplexityPolicyResponse.AsObject> {
    const req = new GetPasswordComplexityPolicyRequest();
    return this.grpcService.mgmt.getPasswordComplexityPolicy(req, null).then((resp) => resp.toObject());
  }

  public addCustomPasswordComplexityPolicy(
    hasLowerCase: boolean,
    hasUpperCase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
    minLength: number,
  ): Promise<AddCustomPasswordComplexityPolicyResponse.AsObject> {
    const req = new AddCustomPasswordComplexityPolicyRequest();
    req.setHasLowercase(hasLowerCase);
    req.setHasUppercase(hasUpperCase);
    req.setHasNumber(hasNumber);
    req.setHasSymbol(hasSymbol);
    req.setMinLength(minLength);
    return this.grpcService.mgmt.addCustomPasswordComplexityPolicy(req, null).then((resp) => resp.toObject());
  }

  public resetPasswordComplexityPolicyToDefault(): Promise<ResetPasswordComplexityPolicyToDefaultResponse.AsObject> {
    const req = new ResetPasswordComplexityPolicyToDefaultRequest();
    return this.grpcService.mgmt.resetPasswordComplexityPolicyToDefault(req, null).then((resp) => resp.toObject());
  }

  public updateCustomPasswordComplexityPolicy(
    hasLowerCase: boolean,
    hasUpperCase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
    minLength: number,
  ): Promise<UpdateCustomPasswordComplexityPolicyResponse.AsObject> {
    const req = new UpdateCustomPasswordComplexityPolicyRequest();
    req.setHasLowercase(hasLowerCase);
    req.setHasUppercase(hasUpperCase);
    req.setHasNumber(hasNumber);
    req.setHasSymbol(hasSymbol);
    req.setMinLength(minLength);
    return this.grpcService.mgmt.updateCustomPasswordComplexityPolicy(req, null).then((resp) => resp.toObject());
  }

  public getLockoutPolicy(): Promise<GetLockoutPolicyResponse.AsObject> {
    const req = new GetLockoutPolicyRequest();
    return this.grpcService.mgmt.getLockoutPolicy(req, null).then((resp) => resp.toObject());
  }

  public addCustomLockoutPolicy(
    maxPasswordAttempts: number,
    maxOTPAttempts: number,
  ): Promise<AddCustomLockoutPolicyResponse.AsObject> {
    const req = new AddCustomLockoutPolicyRequest();
    req.setMaxPasswordAttempts(maxPasswordAttempts);
    req.setMaxOtpAttempts(maxOTPAttempts);

    return this.grpcService.mgmt.addCustomLockoutPolicy(req, null).then((resp) => resp.toObject());
  }

  public resetLockoutPolicyToDefault(): Promise<ResetLockoutPolicyToDefaultResponse.AsObject> {
    const req = new ResetLockoutPolicyToDefaultRequest();
    return this.grpcService.mgmt.resetLockoutPolicyToDefault(req, null).then((resp) => resp.toObject());
  }

  public updateCustomLockoutPolicy(
    maxPasswordAttempts: number,
    maxOTPAttempts: number,
  ): Promise<UpdateCustomLockoutPolicyResponse.AsObject> {
    const req = new UpdateCustomLockoutPolicyRequest();
    req.setMaxPasswordAttempts(maxPasswordAttempts);
    req.setMaxOtpAttempts(maxOTPAttempts);

    return this.grpcService.mgmt.updateCustomLockoutPolicy(req, null).then((resp) => resp.toObject());
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

  /* notification policy */

  public getNotificationPolicy(): Promise<GetNotificationPolicyResponse.AsObject> {
    const req = new GetNotificationPolicyRequest();
    return this.grpcService.mgmt.getNotificationPolicy(req, null).then((resp) => resp.toObject());
  }

  public resetNotificationPolicyToDefault(): Promise<ResetNotificationPolicyToDefaultResponse.AsObject> {
    const req = new ResetNotificationPolicyToDefaultRequest();
    return this.grpcService.mgmt.resetNotificationPolicyToDefault(req, null).then((resp) => resp.toObject());
  }

  public addCustomNotificationPolicy(
    req: AddCustomNotificationPolicyRequest,
  ): Promise<AddCustomNotificationPolicyResponse.AsObject> {
    return this.grpcService.mgmt.addCustomNotificationPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateCustomNotificationPolicy(
    req: UpdateCustomNotificationPolicyRequest,
  ): Promise<UpdateCustomNotificationPolicyResponse.AsObject> {
    return this.grpcService.mgmt.updateCustomNotificationPolicy(req, null).then((resp) => resp.toObject());
  }

  public getUserByID(id: string): Promise<GetUserByIDResponse.AsObject> {
    const req = new GetUserByIDRequest();
    req.setId(id);
    return this.grpcService.mgmt.getUserByID(req, null).then((resp) => resp.toObject());
  }

  public listUserMetadata(
    userId: string,
    offset?: number,
    limit?: number,
    queryList?: MetadataQuery[],
  ): Promise<ListUserMetadataResponse.AsObject> {
    const req = new ListUserMetadataRequest();

    req.setId(userId);
    const metadata = new ListQuery();
    if (offset) {
      metadata.setOffset(offset);
    }
    if (limit) {
      metadata.setLimit(limit);
    }
    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listUserMetadata(req, null).then((resp) => resp.toObject());
  }

  public listOrgMetadata(
    offset?: number,
    limit?: number,
    queryList?: MetadataQuery[],
  ): Promise<ListOrgMetadataResponse.AsObject> {
    const req = new ListOrgMetadataRequest();

    const metadata = new ListQuery();
    if (offset) {
      metadata.setOffset(offset);
    }
    if (limit) {
      metadata.setLimit(limit);
    }
    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listOrgMetadata(req, null).then((resp) => resp.toObject());
  }

  public getUserMetadata(userId: string, key: string): Promise<GetUserMetadataResponse.AsObject> {
    const req = new GetUserMetadataRequest();
    req.setId(userId);
    req.setKey(key);
    return this.grpcService.mgmt.getUserMetadata(req, null).then((resp) => resp.toObject());
  }

  public setUserMetadata(key: string, value: string, userId: string): Promise<SetUserMetadataResponse.AsObject> {
    const req = new SetUserMetadataRequest();
    req.setKey(key);
    req.setValue(value);
    req.setId(userId);
    return this.grpcService.mgmt.setUserMetadata(req, null).then((resp) => resp.toObject());
  }

  public setOrgMetadata(key: string, value: string): Promise<SetOrgMetadataResponse.AsObject> {
    const req = new SetOrgMetadataRequest();
    req.setKey(key);
    req.setValue(value);
    return this.grpcService.mgmt.setOrgMetadata(req, null).then((resp) => resp.toObject());
  }

  public bulkSetUserMetadata(
    list: BulkSetUserMetadataRequest.Metadata[],
    userId: string,
  ): Promise<BulkSetUserMetadataResponse.AsObject> {
    const req = new BulkSetUserMetadataRequest();
    req.setMetadataList(list);
    req.setId(userId);
    return this.grpcService.mgmt.bulkSetUserMetadata(req, null).then((resp) => resp.toObject());
  }

  public removeUserMetadata(key: string, userId: string): Promise<RemoveUserMetadataResponse.AsObject> {
    const req = new RemoveUserMetadataRequest();
    req.setKey(key);
    req.setId(userId);
    return this.grpcService.mgmt.removeUserMetadata(req, null).then((resp) => resp.toObject());
  }

  public removeOrgMetadata(key: string): Promise<RemoveOrgMetadataResponse.AsObject> {
    const req = new RemoveOrgMetadataRequest();
    req.setKey(key);
    return this.grpcService.mgmt.removeOrgMetadata(req, null).then((resp) => resp.toObject());
  }

  public removeUser(id: string): Promise<RemoveUserResponse.AsObject> {
    const req = new RemoveUserRequest();
    req.setId(id);
    return this.grpcService.mgmt.removeUser(req, null).then((resp) => resp.toObject());
  }

  public removeOrg(): Promise<RemoveOrgResponse.AsObject> {
    const req = new RemoveOrgRequest();
    return this.grpcService.mgmt.removeOrg(req, null).then((resp) => resp.toObject());
  }

  public listProjectMembers(
    projectId: string,
    limit: number,
    offset: number,
    queryList?: SearchQuery[],
  ): Promise<ListProjectMembersResponse.AsObject> {
    const req = new ListProjectMembersRequest();
    const query = new ListQuery();
    req.setQuery(query);
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
    req.setQuery(query);
    return this.grpcService.mgmt.listProjectMembers(req, null).then((resp) => resp.toObject());
  }

  public listUserMemberships(
    userId: string,
    limit: number,
    offset: number,
    queryList?: MembershipQuery[],
  ): Promise<ListUserMembershipsResponse.AsObject> {
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
    req.setQuery(metadata);
    return this.grpcService.mgmt.listUserMemberships(req, null).then((resp) => resp.toObject());
  }

  public getHumanProfile(userId: string): Promise<GetHumanProfileResponse.AsObject> {
    const req = new GetHumanProfileRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.getHumanProfile(req, null).then((resp) => resp.toObject());
  }

  public listHumanMultiFactors(userId: string): Promise<ListHumanAuthFactorsResponse.AsObject> {
    const req = new ListHumanAuthFactorsRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.listHumanAuthFactors(req, null).then((resp) => resp.toObject());
  }

  public removeHumanMultiFactorOTP(userId: string): Promise<RemoveHumanAuthFactorOTPResponse.AsObject> {
    const req = new RemoveHumanAuthFactorOTPRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.removeHumanAuthFactorOTP(req, null).then((resp) => resp.toObject());
  }

  public removeHumanAuthFactorU2F(userId: string, tokenId: string): Promise<RemoveHumanAuthFactorU2FResponse.AsObject> {
    const req = new RemoveHumanAuthFactorU2FRequest();
    req.setUserId(userId);
    req.setTokenId(tokenId);
    return this.grpcService.mgmt.removeHumanAuthFactorU2F(req, null).then((resp) => resp.toObject());
  }

  public removeHumanAuthFactorOTPSMS(userId: string): Promise<RemoveHumanAuthFactorOTPSMSResponse.AsObject> {
    const req = new RemoveHumanAuthFactorOTPSMSRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.removeHumanAuthFactorOTPSMS(req, null).then((resp) => resp.toObject());
  }

  public removeHumanAuthFactorOTPEmail(userId: string): Promise<RemoveHumanAuthFactorOTPEmailResponse.AsObject> {
    const req = new RemoveHumanAuthFactorOTPEmailRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.removeHumanAuthFactorOTPEmail(req, null).then((resp) => resp.toObject());
  }

  public updateHumanProfile(
    userId: string,
    firstName?: string,
    lastName?: string,
    nickName?: string,
    displayName?: string,
    preferredLanguage?: string,
    gender?: Gender,
  ): Promise<UpdateHumanProfileResponse.AsObject> {
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
    if (displayName) {
      req.setDisplayName(displayName);
    }
    if (gender) {
      req.setGender(gender);
    }
    if (preferredLanguage) {
      req.setPreferredLanguage(preferredLanguage);
    }
    return this.grpcService.mgmt.updateHumanProfile(req, null).then((resp) => resp.toObject());
  }

  public getHumanEmail(id: string): Promise<GetHumanEmailResponse.AsObject> {
    const req = new GetHumanEmailRequest();
    req.setUserId(id);
    return this.grpcService.mgmt.getHumanEmail(req, null).then((resp) => resp.toObject());
  }

  public updateHumanEmail(userId: string, email: string, isVerified?: boolean): Promise<UpdateHumanEmailResponse.AsObject> {
    const req = new UpdateHumanEmailRequest();
    req.setUserId(userId);
    req.setEmail(email);
    if (isVerified) {
      req.setIsEmailVerified(isVerified);
    }
    return this.grpcService.mgmt.updateHumanEmail(req, null).then((resp) => resp.toObject());
  }

  public getHumanPhone(userId: string): Promise<GetHumanPhoneResponse.AsObject> {
    const req = new GetHumanPhoneRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.getHumanPhone(req, null).then((resp) => resp.toObject());
  }

  public updateHumanPhone(userId: string, phone: string): Promise<UpdateHumanPhoneResponse.AsObject> {
    const req = new UpdateHumanPhoneRequest();
    req.setUserId(userId);
    req.setPhone(phone);
    return this.grpcService.mgmt.updateHumanPhone(req, null).then((resp) => resp.toObject());
  }

  public removeHumanPhone(userId: string): Promise<RemoveHumanPhoneResponse.AsObject> {
    const req = new RemoveHumanPhoneRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.removeHumanPhone(req, null).then((resp) => resp.toObject());
  }

  public deactivateUser(id: string): Promise<DeactivateUserResponse.AsObject> {
    const req = new DeactivateUserRequest();
    req.setId(id);
    return this.grpcService.mgmt.deactivateUser(req, null).then((resp) => resp.toObject());
  }

  public addUserGrant(
    userId: string,
    roleNamesList: string[],
    projectId?: string,
    projectGrantId?: string,
  ): Promise<AddUserGrantResponse.AsObject> {
    const req = new AddUserGrantRequest();
    if (projectId) {
      req.setProjectId(projectId);
    }
    if (projectGrantId) {
      req.setProjectGrantId(projectGrantId);
    }
    req.setUserId(userId);
    req.setRoleKeysList(roleNamesList);

    return this.grpcService.mgmt.addUserGrant(req, null).then((resp) => resp.toObject());
  }

  public reactivateUser(id: string): Promise<ReactivateUserResponse.AsObject> {
    const req = new ReactivateUserRequest();
    req.setId(id);
    return this.grpcService.mgmt.reactivateUser(req, null).then((resp) => resp.toObject());
  }

  public addProjectRole(
    projectId: string,
    roleKey: string,
    displayName: string,
    group: string,
  ): Promise<AddProjectRoleResponse.AsObject> {
    const req = new AddProjectRoleRequest();
    req.setProjectId(projectId);
    req.setRoleKey(roleKey);
    if (displayName) {
      req.setDisplayName(displayName);
    }
    req.setGroup(group);
    return this.grpcService.mgmt.addProjectRole(req, null).then((resp) => resp.toObject());
  }

  public resendHumanEmailVerification(userId: string): Promise<any> {
    const req = new ResendHumanEmailVerificationRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.resendHumanEmailVerification(req, null).then((resp) => resp.toObject());
  }

  public resendHumanInitialization(userId: string, newemail: string): Promise<ResendHumanInitializationResponse.AsObject> {
    const req = new ResendHumanInitializationRequest();
    if (newemail) {
      req.setEmail(newemail);
    }
    req.setUserId(userId);

    return this.grpcService.mgmt.resendHumanInitialization(req, null).then((resp) => resp.toObject());
  }

  public resendHumanPhoneVerification(userId: string): Promise<any> {
    const req = new ResendHumanPhoneVerificationRequest();
    req.setUserId(userId);
    return this.grpcService.mgmt.resendHumanPhoneVerification(req, null).then((resp) => resp.toObject());
  }

  public setHumanInitialPassword(id: string, password: string): Promise<any> {
    const req = new SetHumanInitialPasswordRequest();
    req.setUserId(id);
    req.setPassword(password);
    return this.grpcService.mgmt.setHumanInitialPassword(req, null).then((resp) => resp.toObject());
  }

  public sendHumanResetPasswordNotification(id: string, type: SendHumanResetPasswordNotificationRequest.Type): Promise<any> {
    const req = new SendHumanResetPasswordNotificationRequest();
    req.setUserId(id);
    req.setType(type);
    return this.grpcService.mgmt.sendHumanResetPasswordNotification(req, null).then((resp) => resp.toObject());
  }

  public listUsers(
    limit: number,
    offset: number,
    queriesList?: UserSearchQuery[],
    sortingColumn?: UserFieldName,
    sortingDirection?: SortDirection,
  ): Promise<ListUsersResponse.AsObject> {
    const req = new ListUsersRequest();
    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }
    if (sortingDirection) {
      query.setAsc(sortingDirection === 'asc');
    }
    req.setQuery(query);
    if (sortingColumn) {
      req.setSortingColumn(sortingColumn);
    }

    if (queriesList) {
      req.setQueriesList(queriesList);
    }
    return this.grpcService.mgmt.listUsers(req, null).then((resp) => resp.toObject());
  }

  public getUserByLoginNameGlobal(loginname: string): Promise<GetUserByLoginNameGlobalResponse.AsObject> {
    const req = new GetUserByLoginNameGlobalRequest();
    req.setLoginName(loginname);
    return this.grpcService.mgmt.getUserByLoginNameGlobal(req, null).then((resp) => resp.toObject());
  }

  // USER GRANTS

  public listUserGrants(
    limit?: number,
    offset?: number,
    queryList?: UserGrantQuery[],
  ): Promise<ListUserGrantResponse.AsObject> {
    const req = new ListUserGrantRequest();
    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }
    req.setQuery(query);

    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listUserGrants(req, null).then((resp) => resp.toObject());
  }

  public getUserGrantByID(grantId: string, userId: string): Promise<GetUserGrantByIDResponse.AsObject> {
    const req = new GetUserGrantByIDRequest();
    req.setGrantId(grantId);
    req.setUserId(userId);

    return this.grpcService.mgmt.getUserGrantByID(req, null).then((resp) => resp.toObject());
  }

  public updateUserGrant(
    grantId: string,
    userId: string,
    roleKeysList: string[],
  ): Promise<UpdateUserGrantResponse.AsObject> {
    const req = new UpdateUserGrantRequest();
    req.setGrantId(grantId);
    req.setRoleKeysList(roleKeysList);
    req.setUserId(userId);

    return this.grpcService.mgmt.updateUserGrant(req, null).then((resp) => resp.toObject());
  }

  public removeUserGrant(grantId: string, userId: string): Promise<RemoveUserGrantResponse.AsObject> {
    const req = new RemoveUserGrantRequest();
    req.setGrantId(grantId);
    req.setUserId(userId);

    return this.grpcService.mgmt.removeUserGrant(req, null).then((resp) => resp.toObject());
  }

  public bulkRemoveUserGrant(grantIdsList: string[]): Promise<BulkRemoveUserGrantResponse.AsObject> {
    const req = new BulkRemoveUserGrantRequest();
    req.setGrantIdList(grantIdsList);

    return this.grpcService.mgmt.bulkRemoveUserGrant(req, null).then((resp) => resp.toObject());
  }

  public listAppChanges(
    appId: string,
    projectId: string,
    limit: number,
    sequence: number,
  ): Promise<ListAppChangesResponse.AsObject> {
    const req = new ListAppChangesRequest();
    const query = new ChangeQuery();
    req.setAppId(appId);
    req.setProjectId(projectId);

    if (limit) {
      query.setLimit(limit);
    }
    if (sequence) {
      query.setSequence(sequence);
    }
    req.setQuery(query);
    return this.grpcService.mgmt.listAppChanges(req, null).then((resp) => resp.toObject());
  }

  public listOrgChanges(limit: number, sequence: number): Promise<ListOrgChangesResponse.AsObject> {
    const req = new ListOrgChangesRequest();
    const query = new ChangeQuery();

    if (limit) {
      query.setLimit(limit);
    }
    if (sequence) {
      query.setSequence(sequence);
    }

    req.setQuery(query);
    return this.grpcService.mgmt.listOrgChanges(req, null).then((resp) => resp.toObject());
  }

  public listProjectChanges(
    projectId: string,
    limit: number,
    sequence: number,
  ): Promise<ListProjectChangesResponse.AsObject> {
    const req = new ListProjectChangesRequest();
    req.setProjectId(projectId);
    const query = new ChangeQuery();

    if (limit) {
      query.setLimit(limit);
    }
    if (sequence) {
      query.setSequence(sequence);
    }

    req.setQuery(query);
    return this.grpcService.mgmt.listProjectChanges(req, null).then((resp) => resp.toObject());
  }

  public listProjectGrantChanges(
    projectId: string,
    grantId: string,
    limit: number,
    sequence: number,
  ): Promise<ListProjectGrantChangesResponse.AsObject> {
    const req = new ListProjectGrantChangesRequest();
    req.setProjectId(projectId);
    req.setGrantId(grantId);
    const query = new ChangeQuery();

    if (limit) {
      query.setLimit(limit);
    }
    if (sequence) {
      query.setSequence(sequence);
    }

    req.setQuery(query);
    return this.grpcService.mgmt.listProjectGrantChanges(req, null).then((resp) => resp.toObject());
  }

  public listUserChanges(userId: string, limit: number, sequence: number): Promise<ListUserChangesResponse.AsObject> {
    const req = new ListUserChangesRequest();
    req.setUserId(userId);
    const query = new ChangeQuery();

    if (limit) {
      query.setLimit(limit);
    }
    if (sequence) {
      query.setSequence(sequence);
    }

    req.setQuery(query);
    return this.grpcService.mgmt.listUserChanges(req, null).then((resp) => resp.toObject());
  }

  // project

  public listProjects(limit?: number, offset?: number, queryList?: ProjectQuery[]): Promise<ListProjectsResponse.AsObject> {
    const req = new ListProjectsRequest();
    const query = new ListQuery();

    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }

    req.setQuery(query);

    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listProjects(req, null).then((value) => {
      const obj = value.toObject();
      const count = obj.resultList.length;
      if (count >= 0) {
        this.ownedProjects.next(obj.resultList);
        this.ownedProjectsCount.next(count);
      }

      return obj;
    });
  }

  public listGrantedProjects(
    limit?: number,
    offset?: number,
    queryList?: ProjectQuery[],
  ): Promise<ListGrantedProjectsResponse.AsObject> {
    const req = new ListGrantedProjectsRequest();
    const query = new ListQuery();

    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }

    req.setQuery(query);
    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listGrantedProjects(req, null).then((value) => {
      const obj = value.toObject();
      this.grantedProjects.next(obj.resultList);
      this.grantedProjectsCount.next(obj.resultList.length);
      return obj;
    });
  }

  public getOIDCInformation(): Promise<GetOIDCInformationResponse.AsObject> {
    const req = new GetOIDCInformationRequest();
    return this.grpcService.mgmt.getOIDCInformation(req, null).then((resp) => resp.toObject());
  }

  public getProjectByID(projectId: string): Promise<GetProjectByIDResponse.AsObject> {
    const req = new GetProjectByIDRequest();
    req.setId(projectId);
    return this.grpcService.mgmt.getProjectByID(req, null).then((resp) => resp.toObject());
  }

  public getGrantedProjectByID(projectId: string, grantId: string): Promise<GetGrantedProjectByIDResponse.AsObject> {
    const req = new GetGrantedProjectByIDRequest();
    req.setGrantId(grantId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.getGrantedProjectByID(req, null).then((resp) => resp.toObject());
  }

  public addProject(project: AddProjectRequest.AsObject): Promise<AddProjectResponse.AsObject> {
    const req = new AddProjectRequest();
    req.setName(project.name);
    return this.grpcService.mgmt.addProject(req, null).then((value) => {
      const current = this.ownedProjectsCount.getValue();
      this.ownedProjectsCount.next(current + 1);
      return value.toObject();
    });
  }

  public updateProject(req: UpdateProjectRequest): Promise<UpdateProjectResponse.AsObject> {
    return this.grpcService.mgmt.updateProject(req, null).then((resp) => resp.toObject());
  }

  public updateProjectGrant(
    grantId: string,
    projectId: string,
    rolesList: string[],
  ): Promise<UpdateProjectGrantResponse.AsObject> {
    const req = new UpdateProjectGrantRequest();
    req.setRoleKeysList(rolesList);
    req.setGrantId(grantId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.updateProjectGrant(req, null).then((resp) => resp.toObject());
  }

  public removeProjectGrant(grantId: string, projectId: string): Promise<RemoveProjectGrantResponse.AsObject> {
    const req = new RemoveProjectGrantRequest();
    req.setGrantId(grantId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.removeProjectGrant(req, null).then((resp) => resp.toObject());
  }

  public deactivateProject(projectId: string): Promise<DeactivateProjectResponse.AsObject> {
    const req = new DeactivateProjectRequest();
    req.setId(projectId);
    return this.grpcService.mgmt.deactivateProject(req, null).then((resp) => resp.toObject());
  }

  public reactivateProject(projectId: string): Promise<ReactivateProjectResponse.AsObject> {
    const req = new ReactivateProjectRequest();
    req.setId(projectId);
    return this.grpcService.mgmt.reactivateProject(req, null).then((resp) => resp.toObject());
  }

  public listProjectGrants(projectId: string, limit: number, offset: number): Promise<ListProjectGrantsResponse.AsObject> {
    const req = new ListProjectGrantsRequest();
    req.setProjectId(projectId);
    const query = new ListQuery();

    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }

    req.setQuery(query);
    return this.grpcService.mgmt.listProjectGrants(req, null).then((resp) => resp.toObject());
  }

  public listProjectGrantMemberRoles(): Promise<ListProjectGrantMemberRolesResponse.AsObject> {
    const req = new ListProjectGrantMemberRolesRequest();
    return this.grpcService.mgmt.listProjectGrantMemberRoles(req, null).then((resp) => resp.toObject());
  }

  public addProjectMember(
    projectId: string,
    userId: string,
    rolesList: string[],
  ): Promise<AddProjectMemberResponse.AsObject> {
    const req = new AddProjectMemberRequest();
    req.setProjectId(projectId);
    req.setUserId(userId);
    req.setRolesList(rolesList);
    return this.grpcService.mgmt.addProjectMember(req, null).then((resp) => resp.toObject());
  }

  public updateProjectMember(
    projectId: string,
    userId: string,
    rolesList: string[],
  ): Promise<UpdateProjectMemberResponse.AsObject> {
    const req = new UpdateProjectMemberRequest();
    req.setProjectId(projectId);
    req.setUserId(userId);
    req.setRolesList(rolesList);
    return this.grpcService.mgmt.updateProjectMember(req, null).then((resp) => resp.toObject());
  }

  public addProjectGrantMember(
    projectId: string,
    grantId: string,
    userId: string,
    rolesList: string[],
  ): Promise<AddProjectGrantMemberResponse.AsObject> {
    const req = new AddProjectGrantMemberRequest();
    req.setProjectId(projectId);
    req.setGrantId(grantId);
    req.setUserId(userId);
    req.setRolesList(rolesList);
    return this.grpcService.mgmt.addProjectGrantMember(req, null).then((resp) => resp.toObject());
  }

  public updateProjectGrantMember(
    projectId: string,
    grantId: string,
    userId: string,
    rolesList: string[],
  ): Promise<UpdateProjectGrantMemberResponse.AsObject> {
    const req = new UpdateProjectGrantMemberRequest();
    req.setProjectId(projectId);
    req.setGrantId(grantId);
    req.setUserId(userId);
    req.setRolesList(rolesList);
    return this.grpcService.mgmt.updateProjectGrantMember(req, null).then((resp) => resp.toObject());
  }

  public listProjectGrantMembers(
    projectId: string,
    grantId: string,
    limit: number,
    offset: number,
    queryList?: SearchQuery[],
  ): Promise<ListProjectGrantMembersResponse.AsObject> {
    const req = new ListProjectGrantMembersRequest();
    req.setProjectId(projectId);
    req.setGrantId(grantId);

    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }

    req.setQuery(query);
    if (queryList) {
      req.setQueriesList(queryList);
    }

    return this.grpcService.mgmt.listProjectGrantMembers(req, null).then((resp) => resp.toObject());
  }

  public removeProjectGrantMember(
    projectId: string,
    grantId: string,
    userId: string,
  ): Promise<RemoveProjectGrantMemberResponse.AsObject> {
    const req = new RemoveProjectGrantMemberRequest();
    req.setGrantId(grantId);
    req.setUserId(userId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.removeProjectGrantMember(req, null).then((resp) => resp.toObject());
  }

  public reactivateApp(projectId: string, appId: string): Promise<ReactivateAppResponse.AsObject> {
    const req = new ReactivateAppRequest();
    req.setAppId(appId);
    req.setProjectId(projectId);

    return this.grpcService.mgmt.reactivateApp(req, null).then((resp) => resp.toObject());
  }

  public deactivateApp(projectId: string, appId: string): Promise<DeactivateAppResponse.AsObject> {
    const req = new DeactivateAppRequest();
    req.setAppId(appId);
    req.setProjectId(projectId);

    return this.grpcService.mgmt.deactivateApp(req, null).then((resp) => resp.toObject());
  }

  public regenerateOIDCClientSecret(appId: string, projectId: string): Promise<RegenerateOIDCClientSecretResponse.AsObject> {
    const req = new RegenerateOIDCClientSecretRequest();
    req.setAppId(appId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.regenerateOIDCClientSecret(req, null).then((resp) => resp.toObject());
  }

  public listAppKeys(
    projectId: string,
    appId: string,
    limit: number,
    offset: number,
  ): Promise<ListAppKeysResponse.AsObject> {
    const req = new ListAppKeysRequest();
    req.setProjectId(projectId);
    req.setAppId(appId);
    const metaData = new ListQuery();
    if (limit) {
      metaData.setLimit(limit);
    }
    if (offset) {
      metaData.setOffset(offset);
    }
    req.setQuery(metaData);
    return this.grpcService.mgmt.listAppKeys(req, null).then((resp) => resp.toObject());
  }

  public addAppKey(
    projectId: string,
    appId: string,
    type: KeyType,
    expirationDate?: Timestamp,
  ): Promise<AddAppKeyResponse.AsObject> {
    const req = new AddAppKeyRequest();
    req.setProjectId(projectId);
    req.setAppId(appId);
    req.setType(type);
    if (expirationDate) {
      req.setExpirationDate(expirationDate);
    }
    return this.grpcService.mgmt.addAppKey(req, null).then((resp) => resp.toObject());
  }

  public removeAppKey(projectId: string, appId: string, keyId: string): Promise<RemoveAppKeyResponse.AsObject> {
    const req = new RemoveAppKeyRequest();
    req.setAppId(appId);
    req.setKeyId(keyId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.removeAppKey(req, null).then((resp) => resp.toObject());
  }

  public listProjectRoles(
    projectId: string,
    limit: number,
    offset: number,
    queryList?: RoleQuery[],
  ): Promise<ListProjectRolesResponse.AsObject> {
    const req = new ListProjectRolesRequest();
    req.setProjectId(projectId);

    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }

    req.setQuery(query);
    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listProjectRoles(req, null).then((resp) => resp.toObject());
  }

  public listGrantedProjectRoles(
    projectId: string,
    grantId: string,
    limit?: number,
    offset?: number,
    queryList?: RoleQuery[],
  ): Promise<ListGrantedProjectRolesResponse.AsObject> {
    const req = new ListGrantedProjectRolesRequest();
    req.setProjectId(projectId);
    req.setGrantId(grantId);

    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }

    req.setQuery(query);
    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listGrantedProjectRoles(req, null).then((resp) => resp.toObject());
  }

  public bulkAddProjectRoles(
    projectId: string,
    rolesList: BulkAddProjectRolesRequest.Role[],
  ): Promise<BulkAddProjectRolesResponse.AsObject> {
    const req = new BulkAddProjectRolesRequest();
    req.setProjectId(projectId);
    req.setRolesList(rolesList);
    return this.grpcService.mgmt.bulkAddProjectRoles(req, null).then((resp) => resp.toObject());
  }

  public removeProjectRole(projectId: string, roleKey: string): Promise<RemoveProjectRoleResponse.AsObject> {
    const req = new RemoveProjectRoleRequest();
    req.setProjectId(projectId);
    req.setRoleKey(roleKey);
    return this.grpcService.mgmt.removeProjectRole(req, null).then((resp) => resp.toObject());
  }

  public updateProjectRole(
    projectId: string,
    roleKey: string,
    displayName: string,
    group: string,
  ): Promise<UpdateProjectRoleResponse.AsObject> {
    const req = new UpdateProjectRoleRequest();
    req.setProjectId(projectId);
    req.setRoleKey(roleKey);
    req.setGroup(group);
    req.setDisplayName(displayName);
    return this.grpcService.mgmt.updateProjectRole(req, null).then((resp) => resp.toObject());
  }

  public removeProjectMember(projectId: string, userId: string): Promise<RemoveProjectMemberResponse.AsObject> {
    const req = new RemoveProjectMemberRequest();
    req.setProjectId(projectId);
    req.setUserId(userId);
    return this.grpcService.mgmt.removeProjectMember(req, null).then((resp) => resp.toObject());
  }

  public listApps(
    projectId: string,
    limit: number,
    offset: number,
    queryList?: AppQuery[],
  ): Promise<ListAppsResponse.AsObject> {
    const req = new ListAppsRequest();
    req.setProjectId(projectId);
    const query = new ListQuery();
    if (limit) {
      query.setLimit(limit);
    }
    if (offset) {
      query.setOffset(offset);
    }
    req.setQuery(query);
    if (queryList) {
      req.setQueriesList(queryList);
    }
    return this.grpcService.mgmt.listApps(req, null).then((resp) => resp.toObject());
  }

  public getAppByID(projectId: string, appId: string): Promise<GetAppByIDResponse.AsObject> {
    const req = new GetAppByIDRequest();
    req.setProjectId(projectId);
    req.setAppId(appId);
    return this.grpcService.mgmt.getAppByID(req, null).then((resp) => resp.toObject());
  }

  public listProjectMemberRoles(): Promise<ListProjectMemberRolesResponse.AsObject> {
    const req = new ListProjectMemberRolesRequest();
    return this.grpcService.mgmt.listProjectMemberRoles(req, null).then((resp) => resp.toObject());
  }

  public getProjectGrantByID(grantId: string, projectId: string): Promise<GetProjectGrantByIDResponse.AsObject> {
    const req = new GetProjectGrantByIDRequest();
    req.setGrantId(grantId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.getProjectGrantByID(req, null).then((resp) => resp.toObject());
  }

  public removeProject(id: string): Promise<RemoveProjectResponse.AsObject> {
    const req = new RemoveProjectRequest();
    req.setId(id);
    return this.grpcService.mgmt.removeProject(req, null).then((value) => {
      const current = this.ownedProjectsCount.getValue();
      this.ownedProjectsCount.next(current > 0 ? current - 1 : 0);
      return value.toObject();
    });
  }

  public deactivateProjectGrant(grantId: string, projectId: string): Promise<DeactivateProjectGrantResponse.AsObject> {
    const req = new DeactivateProjectGrantRequest();
    req.setGrantId(grantId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.deactivateProjectGrant(req, null).then((resp) => resp.toObject());
  }

  public reactivateProjectGrant(grantId: string, projectId: string): Promise<ReactivateProjectGrantResponse.AsObject> {
    const req = new ReactivateProjectGrantRequest();
    req.setGrantId(grantId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.reactivateProjectGrant(req, null).then((resp) => resp.toObject());
  }

  public addOIDCApp(req: AddOIDCAppRequest): Promise<AddOIDCAppResponse.AsObject> {
    return this.grpcService.mgmt.addOIDCApp(req, null).then((resp) => resp.toObject());
  }

  public addAPIApp(req: AddAPIAppRequest): Promise<AddAPIAppResponse.AsObject> {
    return this.grpcService.mgmt.addAPIApp(req, null).then((resp) => resp.toObject());
  }

  public addSAMLApp(req: AddSAMLAppRequest): Promise<AddSAMLAppResponse.AsObject> {
    return this.grpcService.mgmt.addSAMLApp(req, null).then((resp) => resp.toObject());
  }

  public regenerateAPIClientSecret(appId: string, projectId: string): Promise<RegenerateAPIClientSecretResponse.AsObject> {
    const req = new RegenerateAPIClientSecretRequest();
    req.setAppId(appId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.regenerateAPIClientSecret(req, null).then((resp) => resp.toObject());
  }

  public updateApp(projectId: string, appId: string, name: string): Promise<UpdateAppResponse.AsObject> {
    const req = new UpdateAppRequest();
    req.setAppId(appId);
    req.setName(name);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.updateApp(req, null).then((resp) => resp.toObject());
  }

  public updateOrg(name: string): Promise<UpdateOrgResponse.AsObject> {
    const req = new UpdateOrgRequest();
    req.setName(name);
    return this.grpcService.mgmt.updateOrg(req, null).then((resp) => resp.toObject());
  }

  public updateOIDCAppConfig(req: UpdateOIDCAppConfigRequest): Promise<UpdateOIDCAppConfigResponse.AsObject> {
    return this.grpcService.mgmt.updateOIDCAppConfig(req, null).then((resp) => resp.toObject());
  }

  public updateAPIAppConfig(req: UpdateAPIAppConfigRequest): Promise<UpdateAPIAppConfigResponse.AsObject> {
    return this.grpcService.mgmt.updateAPIAppConfig(req, null).then((resp) => resp.toObject());
  }

  public updateSAMLAppConfig(req: UpdateSAMLAppConfigRequest): Promise<UpdateSAMLAppConfigResponse.AsObject> {
    return this.grpcService.mgmt.updateSAMLAppConfig(req, null).then((resp) => resp.toObject());
  }

  public removeApp(projectId: string, appId: string): Promise<RemoveAppResponse.AsObject> {
    const req = new RemoveAppRequest();
    req.setAppId(appId);
    req.setProjectId(projectId);
    return this.grpcService.mgmt.removeApp(req, null).then((resp) => resp.toObject());
  }
}

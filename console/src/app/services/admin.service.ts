import { Injectable } from '@angular/core';
import { BehaviorSubject, catchError, finalize, from, map, Observable, of, Subject, switchMap, tap } from 'rxjs';

import {
  ActivateLabelPolicyRequest,
  ActivateLabelPolicyResponse,
  ActivateSMSProviderRequest,
  ActivateSMSProviderResponse,
  ActivateSMTPConfigRequest,
  ActivateSMTPConfigResponse,
  AddAppleProviderRequest,
  AddAppleProviderResponse,
  AddAzureADProviderRequest,
  AddAzureADProviderResponse,
  AddCustomDomainPolicyRequest,
  AddCustomOrgIAMPolicyResponse,
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
  AddIAMMemberRequest,
  AddIAMMemberResponse,
  AddIDPToLoginPolicyRequest,
  AddIDPToLoginPolicyResponse,
  AddJWTProviderRequest,
  AddJWTProviderResponse,
  AddLDAPProviderRequest,
  AddLDAPProviderResponse,
  AddMultiFactorToLoginPolicyRequest,
  AddMultiFactorToLoginPolicyResponse,
  AddNotificationPolicyRequest,
  AddNotificationPolicyResponse,
  AddOIDCSettingsRequest,
  AddOIDCSettingsResponse,
  AddSAMLProviderRequest,
  AddSAMLProviderResponse,
  AddSecondFactorToLoginPolicyRequest,
  AddSecondFactorToLoginPolicyResponse,
  AddSMSProviderTwilioRequest,
  AddSMSProviderTwilioResponse,
  AddSMTPConfigRequest,
  AddSMTPConfigResponse,
  DeactivateIDPRequest,
  DeactivateIDPResponse,
  DeactivateSMSProviderRequest,
  DeactivateSMSProviderResponse,
  DeactivateSMTPConfigRequest,
  DeactivateSMTPConfigResponse,
  DeleteProviderRequest,
  DeleteProviderResponse,
  GetAllowedLanguagesRequest,
  GetAllowedLanguagesResponse,
  GetCustomDomainClaimedMessageTextRequest,
  GetCustomDomainClaimedMessageTextResponse,
  GetCustomDomainPolicyRequest,
  GetCustomDomainPolicyResponse,
  GetCustomInitMessageTextRequest,
  GetCustomInitMessageTextResponse,
  GetCustomLoginTextsRequest,
  GetCustomLoginTextsResponse,
  GetCustomPasswordChangeMessageTextRequest,
  GetCustomPasswordChangeMessageTextResponse,
  GetCustomPasswordlessRegistrationMessageTextRequest,
  GetCustomPasswordlessRegistrationMessageTextResponse,
  GetCustomInviteUserMessageTextRequest,
  GetCustomInviteUserMessageTextResponse,
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
  GetDefaultLanguageRequest,
  GetDefaultLanguageResponse,
  GetDefaultLoginTextsRequest,
  GetDefaultLoginTextsResponse,
  GetDefaultOrgRequest,
  GetDefaultOrgResponse,
  GetDefaultPasswordChangeMessageTextRequest,
  GetDefaultPasswordChangeMessageTextResponse,
  GetDefaultPasswordlessRegistrationMessageTextRequest,
  GetDefaultPasswordlessRegistrationMessageTextResponse,
  GetDefaultInviteUserMessageTextRequest,
  GetDefaultInviteUserMessageTextResponse,
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
  GetFileSystemNotificationProviderRequest,
  GetFileSystemNotificationProviderResponse,
  GetLabelPolicyRequest,
  GetLabelPolicyResponse,
  GetLockoutPolicyRequest,
  GetLockoutPolicyResponse,
  GetLoginPolicyRequest,
  GetLoginPolicyResponse,
  GetLogNotificationProviderRequest,
  GetLogNotificationProviderResponse,
  GetMyInstanceRequest,
  GetMyInstanceResponse,
  GetNotificationPolicyRequest,
  GetNotificationPolicyResponse,
  GetOIDCSettingsRequest,
  GetOIDCSettingsResponse,
  GetPasswordAgePolicyRequest,
  GetPasswordAgePolicyResponse,
  GetPasswordComplexityPolicyRequest,
  GetPasswordComplexityPolicyResponse,
  GetPreviewLabelPolicyRequest,
  GetPreviewLabelPolicyResponse,
  GetPrivacyPolicyRequest,
  GetPrivacyPolicyResponse,
  GetProviderByIDRequest,
  GetProviderByIDResponse,
  GetRestrictionsResponse,
  GetSecretGeneratorRequest,
  GetSecretGeneratorResponse,
  GetSecurityPolicyRequest,
  GetSecurityPolicyResponse,
  GetSMSProviderRequest,
  GetSMSProviderResponse,
  GetSMTPConfigByIdRequest,
  GetSMTPConfigByIdResponse,
  GetSMTPConfigRequest,
  GetSMTPConfigResponse,
  GetSupportedLanguagesRequest,
  GetSupportedLanguagesResponse,
  ListAggregateTypesRequest,
  ListAggregateTypesResponse,
  ListEventsRequest,
  ListEventsResponse,
  ListEventTypesRequest,
  ListEventTypesResponse,
  ListFailedEventsRequest,
  ListFailedEventsResponse,
  ListIAMMemberRolesRequest,
  ListIAMMemberRolesResponse,
  ListIAMMembersRequest,
  ListIAMMembersResponse,
  ListLoginPolicyIDPsRequest,
  ListLoginPolicyIDPsResponse,
  ListLoginPolicyMultiFactorsRequest,
  ListLoginPolicyMultiFactorsResponse,
  ListLoginPolicySecondFactorsRequest,
  ListLoginPolicySecondFactorsResponse,
  ListMilestonesRequest,
  ListMilestonesResponse,
  ListOrgsRequest,
  ListOrgsResponse,
  ListProvidersRequest,
  ListProvidersResponse,
  ListSecretGeneratorsRequest,
  ListSecretGeneratorsResponse,
  ListSMSProvidersRequest,
  ListSMSProvidersResponse,
  ListSMTPConfigsRequest,
  ListSMTPConfigsResponse,
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
  RemoveSMSProviderRequest,
  RemoveSMSProviderResponse,
  RemoveSMTPConfigRequest,
  RemoveSMTPConfigResponse,
  ResetCustomDomainPolicyToDefaultRequest,
  ResetCustomDomainPolicyToDefaultResponse,
  ResetCustomLoginTextsToDefaultRequest,
  ResetCustomLoginTextsToDefaultResponse,
  SelectLanguages,
  SetCustomLoginTextsRequest,
  SetCustomLoginTextsResponse,
  SetDefaultDomainClaimedMessageTextRequest,
  SetDefaultDomainClaimedMessageTextResponse,
  SetDefaultInitMessageTextRequest,
  SetDefaultInitMessageTextResponse,
  SetDefaultLanguageRequest,
  SetDefaultLanguageResponse,
  SetDefaultOrgRequest,
  SetDefaultOrgResponse,
  SetDefaultPasswordChangeMessageTextRequest,
  SetDefaultPasswordChangeMessageTextResponse,
  SetDefaultPasswordlessRegistrationMessageTextRequest,
  SetDefaultPasswordlessRegistrationMessageTextResponse,
  SetDefaultInviteUserMessageTextRequest,
  SetDefaultInviteUserMessageTextResponse,
  SetDefaultPasswordResetMessageTextRequest,
  SetDefaultPasswordResetMessageTextResponse,
  SetDefaultVerifyEmailMessageTextRequest,
  SetDefaultVerifyEmailMessageTextResponse,
  SetDefaultVerifyEmailOTPMessageTextRequest,
  SetDefaultVerifyEmailOTPMessageTextResponse,
  SetDefaultVerifyPhoneMessageTextRequest,
  SetDefaultVerifyPhoneMessageTextResponse,
  SetDefaultVerifySMSOTPMessageTextRequest,
  SetDefaultVerifySMSOTPMessageTextResponse,
  SetRestrictionsRequest,
  SetRestrictionsResponse,
  SetSecurityPolicyRequest,
  SetSecurityPolicyResponse,
  SetUpOrgRequest,
  SetUpOrgResponse,
  TestSMTPConfigByIdRequest,
  TestSMTPConfigByIdResponse,
  TestSMTPConfigRequest,
  TestSMTPConfigResponse,
  UpdateAppleProviderRequest,
  UpdateAppleProviderResponse,
  UpdateAzureADProviderRequest,
  UpdateAzureADProviderResponse,
  UpdateCustomDomainPolicyRequest,
  UpdateCustomDomainPolicyResponse,
  UpdateDomainPolicyRequest,
  UpdateDomainPolicyResponse,
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
  UpdateIAMMemberRequest,
  UpdateIAMMemberResponse,
  UpdateJWTProviderRequest,
  UpdateJWTProviderResponse,
  UpdateLabelPolicyRequest,
  UpdateLabelPolicyResponse,
  UpdateLDAPProviderRequest,
  UpdateLDAPProviderResponse,
  UpdateLockoutPolicyRequest,
  UpdateLockoutPolicyResponse,
  UpdateLoginPolicyRequest,
  UpdateLoginPolicyResponse,
  UpdateNotificationPolicyRequest,
  UpdateNotificationPolicyResponse,
  UpdateOIDCSettingsRequest,
  UpdateOIDCSettingsResponse,
  UpdatePasswordAgePolicyRequest,
  UpdatePasswordAgePolicyResponse,
  UpdatePasswordComplexityPolicyRequest,
  UpdatePasswordComplexityPolicyResponse,
  UpdatePrivacyPolicyRequest,
  UpdatePrivacyPolicyResponse,
  UpdateSAMLProviderRequest,
  UpdateSAMLProviderResponse,
  UpdateSecretGeneratorRequest,
  UpdateSecretGeneratorResponse,
  UpdateSMSProviderTwilioRequest,
  UpdateSMSProviderTwilioResponse,
  UpdateSMSProviderTwilioTokenRequest,
  UpdateSMSProviderTwilioTokenResponse,
  UpdateSMTPConfigPasswordRequest,
  UpdateSMTPConfigPasswordResponse,
  UpdateSMTPConfigRequest,
  UpdateSMTPConfigResponse,
} from '../proto/generated/zitadel/admin_pb';
import { Event } from '../proto/generated/zitadel/event_pb';
import {
  ResetCustomDomainClaimedMessageTextToDefaultRequest,
  ResetCustomDomainClaimedMessageTextToDefaultResponse,
  ResetCustomInitMessageTextToDefaultRequest,
  ResetCustomInitMessageTextToDefaultResponse,
  ResetCustomPasswordChangeMessageTextToDefaultRequest,
  ResetCustomPasswordChangeMessageTextToDefaultResponse,
  ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest,
  ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse,
  ResetCustomInviteUserMessageTextToDefaultRequest,
  ResetCustomInviteUserMessageTextToDefaultResponse,
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
} from '../proto/generated/zitadel/management_pb';
import { SearchQuery } from '../proto/generated/zitadel/member_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { GrpcService } from './grpc.service';
import { StorageLocation, StorageService } from './storage.service';
import {
  IsReachedQuery,
  Milestone,
  MilestoneQuery,
  MilestoneType,
} from '../proto/generated/zitadel/milestone/v1/milestone_pb';
import { OrgFieldName, OrgQuery } from '../proto/generated/zitadel/org_pb';
import { SortDirection } from '@angular/material/sort';
import { SMTPConfig } from '../proto/generated/zitadel/settings_pb';

export interface OnboardingActions {
  order: number;
  milestoneType: MilestoneType;
  link: string;
  externalLink?: boolean;
  fragment?: string | undefined;
  iconClasses?: string;
  darkcolor: string;
  lightcolor: string;
}

type OnboardingMilestone = {
  order: number;
  link: string;
  externalLink?: boolean;
  fragment: string | undefined;
  reached: Milestone.AsObject | undefined;
  iconClasses?: string;
  darkcolor: string;
  lightcolor: string;
};
type OnboardingMilestoneEntries = Array<[string, OnboardingMilestone]> | [];

@Injectable({
  providedIn: 'root',
})
export class AdminService {
  private readonly milestoneTypePrefixLength = 'MILESTONE_TYPE_'.length;
  public hideOnboarding: boolean = false;
  public loadMilestones: Subject<OnboardingActions[]> = new Subject();
  public onboardingLoading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public progressMilestones$: Observable<OnboardingMilestoneEntries> = this.loadMilestones.pipe(
    tap(() => this.onboardingLoading.next(true)),
    switchMap((actions) => {
      const milestonesListQuery = new ListQuery();
      milestonesListQuery.setAsc(true);
      milestonesListQuery.setLimit(20);
      const milestoneIsReachedQuery = new IsReachedQuery().setReached(true);
      const milestonesQuery = new MilestoneQuery().setIsReachedQuery(milestoneIsReachedQuery);
      const milestonesReq = new ListMilestonesRequest().setQuery(milestonesListQuery).setQueriesList([milestonesQuery]);
      return from(this.listMilestones(milestonesReq)).pipe(
        map((reachedMilestones) => {
          let obj: { [type: string]: OnboardingMilestone } = {};
          actions.map((action) => {
            obj[Object.keys(MilestoneType)[action.milestoneType].substring(this.milestoneTypePrefixLength)] = {
              order: action.order,
              link: action.link,
              externalLink: action.externalLink,
              fragment: action.fragment,
              iconClasses: action.iconClasses,
              darkcolor: action.darkcolor,
              lightcolor: action.lightcolor,
              reached: reachedMilestones.resultList.find((reached) => {
                return reached.type.valueOf() == action.milestoneType;
              }),
            };
          });
          const toArray = Object.entries(obj).sort(([key0, a], [key1, b]) => a.order - b.order);

          const toDo = toArray.filter(([key, value]) => value.reached === undefined);
          const done = toArray.filter(([key, value]) => !!value.reached);

          return [...toDo, ...done];
        }),
        tap((milestones) => {
          const total = milestones.length;
          const done = milestones.map(([type, value]) => value.reached !== undefined).filter((res) => !!res).length;
          const percentage = Math.round((done / total) * 100);
          this.progressDone.next(done);
          this.progressTotal.next(total);
          this.progressPercentage.next(percentage);
          this.progressAllDone.next(done === total);
        }),
        catchError((error) => {
          console.error(error);
          return of([]);
        }),
        finalize(() => this.onboardingLoading.next(false)),
      );
    }),
  );

  public progressMilestones: BehaviorSubject<OnboardingMilestoneEntries> = new BehaviorSubject<OnboardingMilestoneEntries>(
    [],
  );
  public progressPercentage: BehaviorSubject<number> = new BehaviorSubject(0);
  public progressDone: BehaviorSubject<number> = new BehaviorSubject(0);
  public progressTotal: BehaviorSubject<number> = new BehaviorSubject(0);
  public progressAllDone: BehaviorSubject<boolean> = new BehaviorSubject(true);

  constructor(
    private readonly grpcService: GrpcService,
    private storageService: StorageService,
  ) {
    this.progressMilestones$.subscribe(this.progressMilestones);

    this.hideOnboarding =
      this.storageService.getItem('onboarding-dismissed', StorageLocation.local) === 'true' ? true : false;
  }

  public getDefaultOrg(): Promise<GetDefaultOrgResponse.AsObject> {
    const req = new GetDefaultOrgRequest();
    return this.grpcService.admin.getDefaultOrg(req, null).then((resp) => resp.toObject());
  }

  public setDefaultOrg(orgId: string): Promise<SetDefaultOrgResponse.AsObject> {
    const req = new SetDefaultOrgRequest();
    req.setOrgId(orgId);

    return this.grpcService.admin.setDefaultOrg(req, null).then((resp) => resp.toObject());
  }

  public listEvents(req: ListEventsRequest): Promise<ListEventsResponse> {
    return this.grpcService.admin.listEvents(req, null).then((resp) => resp);
  }

  public listEventTypes(req: ListEventTypesRequest): Promise<ListEventTypesResponse.AsObject> {
    return this.grpcService.admin.listEventTypes(req, null).then((resp) => resp.toObject());
  }

  public listAggregateTypes(req: ListAggregateTypesRequest): Promise<ListAggregateTypesResponse.AsObject> {
    return this.grpcService.admin.listAggregateTypes(req, null).then((resp) => resp.toObject());
  }

  public getSupportedLanguages(): Promise<GetSupportedLanguagesResponse.AsObject> {
    const req = new GetSupportedLanguagesRequest();
    return this.grpcService.admin.getSupportedLanguages(req, null).then((resp) => resp.toObject());
  }

  public getAllowedLanguages(): Promise<GetAllowedLanguagesResponse.AsObject> {
    const req = new GetAllowedLanguagesRequest();
    return this.grpcService.admin.getAllowedLanguages(req, null).then((resp) => resp.toObject());
  }

  public getDefaultLoginTexts(req: GetDefaultLoginTextsRequest): Promise<GetDefaultLoginTextsResponse.AsObject> {
    return this.grpcService.admin.getDefaultLoginTexts(req, null).then((resp) => resp.toObject());
  }

  public getCustomLoginTexts(req: GetCustomLoginTextsRequest): Promise<GetCustomLoginTextsResponse.AsObject> {
    return this.grpcService.admin.getCustomLoginTexts(req, null).then((resp) => resp.toObject());
  }

  public setCustomLoginText(req: SetCustomLoginTextsRequest): Promise<SetCustomLoginTextsResponse.AsObject> {
    return this.grpcService.admin.setCustomLoginText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomLoginTextToDefault(lang: string): Promise<ResetCustomLoginTextsToDefaultResponse.AsObject> {
    const req = new ResetCustomLoginTextsToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomLoginTextToDefault(req, null).then((resp) => resp.toObject());
  }

  // message texts

  public getDefaultInitMessageText(
    req: GetDefaultInitMessageTextRequest,
  ): Promise<GetDefaultInitMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultInitMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomInitMessageText(req: GetCustomInitMessageTextRequest): Promise<GetCustomInitMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomInitMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomInitMessageTextToDefault(lang: string): Promise<ResetCustomInitMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomInitMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomInitMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public setDefaultInitMessageText(
    req: SetDefaultInitMessageTextRequest,
  ): Promise<SetDefaultInitMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultInitMessageText(req, null).then((resp) => resp.toObject());
  }

  public getDefaultVerifyEmailMessageText(
    req: GetDefaultVerifyEmailMessageTextRequest,
  ): Promise<GetDefaultVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultVerifyEmailMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomVerifyEmailMessageText(
    req: GetCustomVerifyEmailMessageTextRequest,
  ): Promise<GetCustomVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomVerifyEmailMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultVerifyEmailMessageText(
    req: SetDefaultVerifyEmailMessageTextRequest,
  ): Promise<SetDefaultVerifyEmailMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultVerifyEmailMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomVerifyEmailMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomVerifyEmailMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomVerifyEmailMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomVerifyEmailMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultVerifyPhoneMessageText(
    req: GetDefaultVerifyPhoneMessageTextRequest,
  ): Promise<GetDefaultVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultVerifyPhoneMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomVerifyPhoneMessageText(
    req: GetCustomVerifyPhoneMessageTextRequest,
  ): Promise<GetCustomVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomVerifyPhoneMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultVerifyPhoneMessageText(
    req: SetDefaultVerifyPhoneMessageTextRequest,
  ): Promise<SetDefaultVerifyPhoneMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultVerifyPhoneMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomVerifyPhoneMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomVerifyPhoneMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomVerifyPhoneMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomVerifyPhoneMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultVerifySMSOTPMessageText(
    req: GetDefaultVerifySMSOTPMessageTextRequest,
  ): Promise<GetDefaultVerifySMSOTPMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultVerifySMSOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomVerifySMSOTPMessageText(
    req: GetCustomVerifySMSOTPMessageTextRequest,
  ): Promise<GetCustomVerifySMSOTPMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomVerifySMSOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultVerifySMSOTPMessageText(
    req: SetDefaultVerifySMSOTPMessageTextRequest,
  ): Promise<SetDefaultVerifySMSOTPMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultVerifySMSOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomVerifySMSOTPMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomVerifySMSOTPMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomVerifySMSOTPMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomVerifySMSOTPMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultVerifyEmailOTPMessageText(
    req: GetDefaultVerifyEmailOTPMessageTextRequest,
  ): Promise<GetDefaultVerifyEmailOTPMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultVerifyEmailOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomVerifyEmailOTPMessageText(
    req: GetCustomVerifyEmailOTPMessageTextRequest,
  ): Promise<GetCustomVerifyEmailOTPMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomVerifyEmailOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultVerifyEmailOTPMessageText(
    req: SetDefaultVerifyEmailOTPMessageTextRequest,
  ): Promise<SetDefaultVerifyEmailOTPMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultVerifyEmailOTPMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomVerifyEmailOTPMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomVerifyEmailOTPMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomVerifyEmailOTPMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomVerifyEmailOTPMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultPasswordResetMessageText(
    req: GetDefaultPasswordResetMessageTextRequest,
  ): Promise<GetDefaultPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultPasswordResetMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomPasswordResetMessageText(
    req: GetCustomPasswordResetMessageTextRequest,
  ): Promise<GetCustomPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomPasswordResetMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultPasswordResetMessageText(
    req: SetDefaultPasswordResetMessageTextRequest,
  ): Promise<SetDefaultPasswordResetMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultPasswordResetMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomPasswordResetMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomPasswordResetMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomPasswordResetMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomPasswordResetMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultDomainClaimedMessageText(
    req: GetDefaultDomainClaimedMessageTextRequest,
  ): Promise<GetDefaultDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultDomainClaimedMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomDomainClaimedMessageText(
    req: GetCustomDomainClaimedMessageTextRequest,
  ): Promise<GetCustomDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomDomainClaimedMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultDomainClaimedMessageText(
    req: SetDefaultDomainClaimedMessageTextRequest,
  ): Promise<SetDefaultDomainClaimedMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultDomainClaimedMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomDomainClaimedMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomDomainClaimedMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomDomainClaimedMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomDomainClaimedMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultPasswordlessRegistrationMessageText(
    req: GetDefaultPasswordlessRegistrationMessageTextRequest,
  ): Promise<GetDefaultPasswordlessRegistrationMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultPasswordlessRegistrationMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomPasswordlessRegistrationMessageText(
    req: GetCustomPasswordlessRegistrationMessageTextRequest,
  ): Promise<GetCustomPasswordlessRegistrationMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomPasswordlessRegistrationMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultPasswordlessRegistrationMessageText(
    req: SetDefaultPasswordlessRegistrationMessageTextRequest,
  ): Promise<SetDefaultPasswordlessRegistrationMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultPasswordlessRegistrationMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomPasswordlessRegistrationMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin
      .resetCustomPasswordlessRegistrationMessageTextToDefault(req, null)
      .then((resp) => resp.toObject());
  }

  public getDefaultPasswordChangeMessageText(
    req: GetDefaultPasswordChangeMessageTextRequest,
  ): Promise<GetDefaultPasswordChangeMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultPasswordChangeMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomPasswordChangeMessageText(
    req: GetCustomPasswordChangeMessageTextRequest,
  ): Promise<GetCustomPasswordChangeMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomPasswordChangeMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultPasswordChangeMessageText(
    req: SetDefaultPasswordChangeMessageTextRequest,
  ): Promise<SetDefaultPasswordChangeMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultPasswordChangeMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomPasswordChangeMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomPasswordChangeMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomPasswordChangeMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomPasswordChangeMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public getDefaultInviteUserMessageText(
    req: GetDefaultInviteUserMessageTextRequest,
  ): Promise<GetDefaultInviteUserMessageTextResponse.AsObject> {
    return this.grpcService.admin.getDefaultInviteUserMessageText(req, null).then((resp) => resp.toObject());
  }

  public getCustomInviteUserMessageText(
    req: GetCustomInviteUserMessageTextRequest,
  ): Promise<GetCustomInviteUserMessageTextResponse.AsObject> {
    return this.grpcService.admin.getCustomInviteUserMessageText(req, null).then((resp) => resp.toObject());
  }

  public setDefaultInviteUserMessageText(
    req: SetDefaultInviteUserMessageTextRequest,
  ): Promise<SetDefaultInviteUserMessageTextResponse.AsObject> {
    return this.grpcService.admin.setDefaultInviteUserMessageText(req, null).then((resp) => resp.toObject());
  }

  public resetCustomInviteUserMessageTextToDefault(
    lang: string,
  ): Promise<ResetCustomInviteUserMessageTextToDefaultResponse.AsObject> {
    const req = new ResetCustomInviteUserMessageTextToDefaultRequest();
    req.setLanguage(lang);
    return this.grpcService.admin.resetCustomInviteUserMessageTextToDefault(req, null).then((resp) => resp.toObject());
  }

  public SetUpOrg(org: SetUpOrgRequest.Org, human: SetUpOrgRequest.Human): Promise<SetUpOrgResponse.AsObject> {
    const req = new SetUpOrgRequest();

    req.setOrg(org);
    req.setHuman(human);

    return this.grpcService.admin.setUpOrg(req, null).then((resp) => resp.toObject());
  }

  public listLoginPolicyMultiFactors(): Promise<ListLoginPolicyMultiFactorsResponse.AsObject> {
    const req = new ListLoginPolicyMultiFactorsRequest();
    return this.grpcService.admin.listLoginPolicyMultiFactors(req, null).then((resp) => resp.toObject());
  }

  public addMultiFactorToLoginPolicy(
    req: AddMultiFactorToLoginPolicyRequest,
  ): Promise<AddMultiFactorToLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.addMultiFactorToLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public removeMultiFactorFromLoginPolicy(
    req: RemoveMultiFactorFromLoginPolicyRequest,
  ): Promise<RemoveMultiFactorFromLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.removeMultiFactorFromLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public listLoginPolicySecondFactors(): Promise<ListLoginPolicySecondFactorsResponse.AsObject> {
    const req = new ListLoginPolicySecondFactorsRequest();
    return this.grpcService.admin.listLoginPolicySecondFactors(req, null).then((resp) => resp.toObject());
  }

  public addSecondFactorToLoginPolicy(
    req: AddSecondFactorToLoginPolicyRequest,
  ): Promise<AddSecondFactorToLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.addSecondFactorToLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public removeSecondFactorFromLoginPolicy(
    req: RemoveSecondFactorFromLoginPolicyRequest,
  ): Promise<RemoveSecondFactorFromLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.removeSecondFactorFromLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public listIAMMemberRoles(): Promise<ListIAMMemberRolesResponse.AsObject> {
    const req = new ListIAMMemberRolesRequest();
    return this.grpcService.admin.listIAMMemberRoles(req, null).then((resp) => resp.toObject());
  }

  public listViews(): Promise<ListViewsResponse.AsObject> {
    const req = new ListViewsRequest();
    return this.grpcService.admin.listViews(req, null).then((resp) => resp.toObject());
  }

  public listFailedEvents(): Promise<ListFailedEventsResponse.AsObject> {
    const req = new ListFailedEventsRequest();
    return this.grpcService.admin.listFailedEvents(req, null).then((resp) => resp.toObject());
  }

  public removeFailedEvent(viewname: string, db: string, sequence: number): Promise<RemoveFailedEventResponse.AsObject> {
    const req = new RemoveFailedEventRequest();
    req.setDatabase(db);
    req.setViewName(viewname);
    req.setFailedSequence(sequence);
    return this.grpcService.admin.removeFailedEvent(req, null).then((resp) => resp.toObject());
  }

  public getMyInstance(): Promise<GetMyInstanceResponse.AsObject> {
    const req = new GetMyInstanceRequest();
    return this.grpcService.admin.getMyInstance(req, null).then((resp) => resp.toObject());
  }

  public getPrivacyPolicy(): Promise<GetPrivacyPolicyResponse.AsObject> {
    const req = new GetPrivacyPolicyRequest();
    return this.grpcService.admin.getPrivacyPolicy(req, null).then((resp) => resp.toObject());
  }

  public updatePrivacyPolicy(req: UpdatePrivacyPolicyRequest): Promise<UpdatePrivacyPolicyResponse.AsObject> {
    return this.grpcService.admin.updatePrivacyPolicy(req, null).then((resp) => resp.toObject());
  }

  /* Policies */

  /* complexity */

  public getPasswordComplexityPolicy(): Promise<GetPasswordComplexityPolicyResponse.AsObject> {
    const req = new GetPasswordComplexityPolicyRequest();
    return this.grpcService.admin.getPasswordComplexityPolicy(req, null).then((resp) => resp.toObject());
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
    return this.grpcService.admin.updatePasswordComplexityPolicy(req, null).then((resp) => resp.toObject());
  }

  /* age */

  public getPasswordAgePolicy(): Promise<GetPasswordAgePolicyResponse.AsObject> {
    const req = new GetPasswordAgePolicyRequest();

    return this.grpcService.admin.getPasswordAgePolicy(req, null).then((resp) => resp.toObject());
  }

  public updatePasswordAgePolicy(
    maxAgeDays: number,
    expireWarnDays: number,
  ): Promise<UpdatePasswordAgePolicyResponse.AsObject> {
    const req = new UpdatePasswordAgePolicyRequest();
    req.setMaxAgeDays(maxAgeDays);
    req.setExpireWarnDays(expireWarnDays);

    return this.grpcService.admin.updatePasswordAgePolicy(req, null).then((resp) => resp.toObject());
  }

  /* default language */

  public getDefaultLanguage(): Promise<GetDefaultLanguageResponse.AsObject> {
    const req = new GetDefaultLanguageRequest();
    return this.grpcService.admin.getDefaultLanguage(req, null).then((resp) => resp.toObject());
  }

  public setDefaultLanguage(language: string): Promise<SetDefaultLanguageResponse.AsObject> {
    const req = new SetDefaultLanguageRequest();
    req.setLanguage(language);

    return this.grpcService.admin.setDefaultLanguage(req, null).then((resp) => resp.toObject());
  }

  /* restrictions */

  public getRestrictions(): Promise<GetRestrictionsResponse.AsObject> {
    const req = new GetDefaultLanguageRequest();
    return this.grpcService.admin.getRestrictions(req, null).then((resp) => resp.toObject());
  }

  public setRestrictions(
    disallowPublicOrgRegistration?: boolean,
    allowedLanguages?: string[],
  ): Promise<SetRestrictionsResponse.AsObject> {
    const req = new SetRestrictionsRequest();
    if (disallowPublicOrgRegistration !== undefined) {
      req.setDisallowPublicOrgRegistration(disallowPublicOrgRegistration);
    }
    if (allowedLanguages !== undefined) {
      const langs = new SelectLanguages();
      langs.setListList(allowedLanguages);
      req.setAllowedLanguages(langs);
    }
    return this.grpcService.admin.setRestrictions(req, null).then((resp) => resp.toObject());
  }

  /* notification policy */

  public getNotificationPolicy(): Promise<GetNotificationPolicyResponse.AsObject> {
    const req = new GetNotificationPolicyRequest();
    return this.grpcService.admin.getNotificationPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateNotificationPolicy(req: UpdateNotificationPolicyRequest): Promise<UpdateNotificationPolicyResponse.AsObject> {
    return this.grpcService.admin.updateNotificationPolicy(req, null).then((resp) => resp.toObject());
  }

  public addNotificationPolicy(req: AddNotificationPolicyRequest): Promise<AddNotificationPolicyResponse.AsObject> {
    return this.grpcService.admin.addNotificationPolicy(req, null).then((resp) => resp.toObject());
  }

  /* security policy */

  public getSecurityPolicy(): Promise<GetSecurityPolicyResponse.AsObject> {
    const req = new GetSecurityPolicyRequest();
    return this.grpcService.admin.getSecurityPolicy(req, null).then((resp) => resp.toObject());
  }

  public setSecurityPolicy(req: SetSecurityPolicyRequest): Promise<SetSecurityPolicyResponse.AsObject> {
    return this.grpcService.admin.setSecurityPolicy(req, null).then((resp) => resp.toObject());
  }

  /* notification settings */

  public getSMTPConfig(): Promise<GetSMTPConfigResponse.AsObject> {
    const req = new GetSMTPConfigRequest();
    return this.grpcService.admin.getSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  public getSMTPConfigById(id: string): Promise<GetSMTPConfigByIdResponse.AsObject> {
    const req = new GetSMTPConfigByIdRequest();
    req.setId(id);
    return this.grpcService.admin.getSMTPConfigById(req, null).then((resp) => resp.toObject());
  }

  public listSMTPConfigs(): Promise<ListSMTPConfigsResponse.AsObject> {
    const req = new ListSMTPConfigsRequest();
    return this.grpcService.admin.listSMTPConfigs(req, null).then((resp) => resp.toObject());
  }

  public addSMTPConfig(req: AddSMTPConfigRequest): Promise<AddSMTPConfigResponse.AsObject> {
    return this.grpcService.admin.addSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  public updateSMTPConfig(req: UpdateSMTPConfigRequest): Promise<UpdateSMTPConfigResponse.AsObject> {
    return this.grpcService.admin.updateSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  public updateSMTPConfigPassword(req: UpdateSMTPConfigPasswordRequest): Promise<UpdateSMTPConfigPasswordResponse.AsObject> {
    return this.grpcService.admin.updateSMTPConfigPassword(req, null).then((resp) => resp.toObject());
  }

  public activateSMTPConfig(id: string): Promise<ActivateSMTPConfigResponse.AsObject> {
    const req = new ActivateSMTPConfigRequest();
    req.setId(id);
    return this.grpcService.admin.activateSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  public deactivateSMTPConfig(id: string): Promise<DeactivateSMTPConfigResponse.AsObject> {
    const req = new DeactivateSMTPConfigRequest();
    req.setId(id);
    return this.grpcService.admin.deactivateSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  public testSMTPConfigById(req: TestSMTPConfigByIdRequest): Promise<TestSMTPConfigByIdResponse.AsObject> {
    return this.grpcService.admin.testSMTPConfigById(req, null).then((resp) => resp.toObject());
  }

  public testSMTPConfig(req: TestSMTPConfigRequest): Promise<TestSMTPConfigResponse.AsObject> {
    return this.grpcService.admin.testSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  public removeSMTPConfig(id: string): Promise<RemoveSMTPConfigResponse.AsObject> {
    const req = new RemoveSMTPConfigRequest();
    req.setId(id);
    return this.grpcService.admin.removeSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  /* sms */

  public listSMSProviders(): Promise<ListSMSProvidersResponse.AsObject> {
    const req = new ListSMSProvidersRequest();
    return this.grpcService.admin.listSMSProviders(req, null).then((resp) => resp.toObject());
  }

  public getSMSProvider(): Promise<GetSMSProviderResponse.AsObject> {
    const req = new GetSMSProviderRequest();
    return this.grpcService.admin.getSMSProvider(req, null).then((resp) => resp.toObject());
  }

  public addSMSProviderTwilio(req: AddSMSProviderTwilioRequest): Promise<AddSMSProviderTwilioResponse.AsObject> {
    return this.grpcService.admin.addSMSProviderTwilio(req, null).then((resp) => resp.toObject());
  }

  public updateSMSProviderTwilio(req: UpdateSMSProviderTwilioRequest): Promise<UpdateSMSProviderTwilioResponse.AsObject> {
    return this.grpcService.admin.updateSMSProviderTwilio(req, null).then((resp) => resp.toObject());
  }

  public removeSMSProvider(id: string): Promise<RemoveSMSProviderResponse.AsObject> {
    const req = new RemoveSMSProviderRequest();
    req.setId(id);
    return this.grpcService.admin.removeSMSProvider(req, null).then((resp) => resp.toObject());
  }

  public activateSMSProvider(id: string): Promise<ActivateSMSProviderResponse.AsObject> {
    const req = new ActivateSMSProviderRequest();
    req.setId(id);
    return this.grpcService.admin.activateSMSProvider(req, null).then((resp) => resp.toObject());
  }

  public deactivateSMSProvider(id: string): Promise<DeactivateSMSProviderResponse.AsObject> {
    const req = new DeactivateSMSProviderRequest();
    req.setId(id);
    return this.grpcService.admin.deactivateSMSProvider(req, null).then((resp) => resp.toObject());
  }

  public updateSMSProviderTwilioToken(
    req: UpdateSMSProviderTwilioTokenRequest,
  ): Promise<UpdateSMSProviderTwilioTokenResponse.AsObject> {
    return this.grpcService.admin.updateSMSProviderTwilioToken(req, null).then((resp) => resp.toObject());
  }

  /* lockout */

  public getLockoutPolicy(): Promise<GetLockoutPolicyResponse.AsObject> {
    const req = new GetLockoutPolicyRequest();
    return this.grpcService.admin.getLockoutPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateLockoutPolicy(
    maxPasswordAttempts: number,
    maxOTPAttempts: number,
  ): Promise<UpdateLockoutPolicyResponse.AsObject> {
    const req = new UpdateLockoutPolicyRequest();
    req.setMaxPasswordAttempts(maxPasswordAttempts);
    req.setMaxOtpAttempts(maxOTPAttempts);

    return this.grpcService.admin.updateLockoutPolicy(req, null).then((resp) => resp.toObject());
  }

  /* label */

  public getLabelPolicy(): Promise<GetLabelPolicyResponse.AsObject> {
    const req = new GetLabelPolicyRequest();
    return this.grpcService.admin.getLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateLabelPolicy(req: UpdateLabelPolicyRequest): Promise<UpdateLabelPolicyResponse.AsObject> {
    return this.grpcService.admin.updateLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public getPreviewLabelPolicy(): Promise<GetPreviewLabelPolicyResponse.AsObject> {
    const req = new GetPreviewLabelPolicyRequest();
    return this.grpcService.admin.getPreviewLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public activateLabelPolicy(): Promise<ActivateLabelPolicyResponse.AsObject> {
    const req = new ActivateLabelPolicyRequest();
    return this.grpcService.admin.activateLabelPolicy(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyFont(): Promise<RemoveLabelPolicyFontResponse.AsObject> {
    const req = new RemoveLabelPolicyFontRequest();
    return this.grpcService.admin.removeLabelPolicyFont(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyIcon(): Promise<RemoveLabelPolicyIconResponse.AsObject> {
    const req = new RemoveLabelPolicyIconRequest();
    return this.grpcService.admin.removeLabelPolicyIcon(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyIconDark(): Promise<RemoveLabelPolicyIconDarkResponse.AsObject> {
    const req = new RemoveLabelPolicyIconDarkRequest();
    return this.grpcService.admin.removeLabelPolicyIconDark(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyLogo(): Promise<RemoveLabelPolicyLogoResponse.AsObject> {
    const req = new RemoveLabelPolicyLogoRequest();
    return this.grpcService.admin.removeLabelPolicyLogo(req, null).then((resp) => resp.toObject());
  }

  public removeLabelPolicyLogoDark(): Promise<RemoveLabelPolicyLogoDarkResponse.AsObject> {
    const req = new RemoveLabelPolicyLogoDarkRequest();
    return this.grpcService.admin.removeLabelPolicyLogoDark(req, null).then((resp) => resp.toObject());
  }

  /* login */

  public getLoginPolicy(): Promise<GetLoginPolicyResponse.AsObject> {
    const req = new GetLoginPolicyRequest();
    return this.grpcService.admin.getLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateLoginPolicy(req: UpdateLoginPolicyRequest): Promise<UpdateLoginPolicyResponse.AsObject> {
    return this.grpcService.admin.updateLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  /* OIDC Configuration */

  public getOIDCSettings(): Promise<GetOIDCSettingsResponse.AsObject> {
    const req = new GetOIDCSettingsRequest();
    return this.grpcService.admin.getOIDCSettings(req, null).then((resp) => resp.toObject());
  }

  public updateOIDCSettings(req: UpdateOIDCSettingsRequest): Promise<UpdateOIDCSettingsResponse.AsObject> {
    return this.grpcService.admin.updateOIDCSettings(req, null).then((resp) => resp.toObject());
  }

  public addOIDCSettings(req: AddOIDCSettingsRequest): Promise<AddOIDCSettingsResponse.AsObject> {
    return this.grpcService.admin.addOIDCSettings(req, null).then((resp) => resp.toObject());
  }

  /* LOG and FILE Notifications */

  public getLogNotificationProvider(): Promise<GetLogNotificationProviderResponse.AsObject> {
    const req = new GetLogNotificationProviderRequest();
    return this.grpcService.admin.getLogNotificationProvider(req, null).then((resp) => resp.toObject());
  }

  public getFileSystemNotificationProvider(): Promise<GetFileSystemNotificationProviderResponse.AsObject> {
    const req = new GetFileSystemNotificationProviderRequest();
    return this.grpcService.admin.getFileSystemNotificationProvider(req, null).then((resp) => resp.toObject());
  }

  /* secrets generator */

  public listSecretGenerators(): Promise<ListSecretGeneratorsResponse.AsObject> {
    const req = new ListSecretGeneratorsRequest();
    return this.grpcService.admin.listSecretGenerators(req, null).then((resp) => resp.toObject());
  }

  public getSecretGenerator(req: GetSecretGeneratorRequest): Promise<GetSecretGeneratorResponse.AsObject> {
    return this.grpcService.admin.getSecretGenerator(req, null).then((resp) => resp.toObject());
  }

  public updateSecretGenerator(req: UpdateSecretGeneratorRequest): Promise<UpdateSecretGeneratorResponse.AsObject> {
    return this.grpcService.admin.updateSecretGenerator(req, null).then((resp) => resp.toObject());
  }

  /* org domain policy */

  public getDomainPolicy(): Promise<GetDomainPolicyResponse.AsObject> {
    const req = new GetDomainPolicyRequest();
    return this.grpcService.admin.getDomainPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateDomainPolicy(req: UpdateDomainPolicyRequest): Promise<UpdateDomainPolicyResponse.AsObject> {
    return this.grpcService.admin.updateDomainPolicy(req, null).then((resp) => resp.toObject());
  }

  public getCustomDomainPolicy(orgId: string): Promise<GetCustomDomainPolicyResponse.AsObject> {
    const req = new GetCustomDomainPolicyRequest();
    req.setOrgId(orgId);
    return this.grpcService.admin.getCustomDomainPolicy(req, null).then((resp) => resp.toObject());
  }

  public addCustomDomainPolicy(req: AddCustomDomainPolicyRequest): Promise<AddCustomOrgIAMPolicyResponse.AsObject> {
    return this.grpcService.admin.addCustomDomainPolicy(req, null).then((resp) => resp.toObject());
  }

  public updateCustomDomainPolicy(req: UpdateCustomDomainPolicyRequest): Promise<UpdateCustomDomainPolicyResponse.AsObject> {
    return this.grpcService.admin.updateCustomDomainPolicy(req, null).then((resp) => resp.toObject());
  }

  public resetCustomDomainPolicyToDefault(orgId: string): Promise<ResetCustomDomainPolicyToDefaultResponse.AsObject> {
    const req = new ResetCustomDomainPolicyToDefaultRequest();
    req.setOrgId(orgId);
    return this.grpcService.admin.resetCustomDomainPolicyToDefault(req, null).then((resp) => resp.toObject());
  }

  /* policies end */

  public addIDPToLoginPolicy(idpId: string): Promise<AddIDPToLoginPolicyResponse.AsObject> {
    const req = new AddIDPToLoginPolicyRequest();
    req.setIdpId(idpId);
    return this.grpcService.admin.addIDPToLoginPolicy(req, null).then((resp) => resp.toObject());
  }

  public removeIDPFromLoginPolicy(idpId: string): Promise<RemoveIDPFromLoginPolicyResponse.AsObject> {
    const req = new RemoveIDPFromLoginPolicyRequest();
    req.setIdpId(idpId);
    return this.grpcService.admin.removeIDPFromLoginPolicy(req, null).then((resp) => resp.toObject());
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
    return this.grpcService.admin.listLoginPolicyIDPs(req, null).then((resp) => resp.toObject());
  }

  public deactivateIDP(id: string): Promise<DeactivateIDPResponse.AsObject> {
    const req = new DeactivateIDPRequest();
    req.setIdpId(id);
    return this.grpcService.admin.deactivateIDP(req, null).then((resp) => resp.toObject());
  }

  public reactivateIDP(id: string): Promise<ReactivateIDPResponse.AsObject> {
    const req = new ReactivateIDPRequest();
    req.setIdpId(id);
    return this.grpcService.admin.reactivateIDP(req, null).then((resp) => resp.toObject());
  }

  //   idp templates

  public addAzureADProvider(req: AddAzureADProviderRequest): Promise<AddAzureADProviderResponse.AsObject> {
    return this.grpcService.admin.addAzureADProvider(req, null).then((resp) => resp.toObject());
  }

  public updateAzureADProvider(req: UpdateAzureADProviderRequest): Promise<UpdateAzureADProviderResponse.AsObject> {
    return this.grpcService.admin.updateAzureADProvider(req, null).then((resp) => resp.toObject());
  }

  public addGoogleProvider(req: AddGoogleProviderRequest): Promise<AddGoogleProviderResponse.AsObject> {
    return this.grpcService.admin.addGoogleProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGoogleProvider(req: UpdateGoogleProviderRequest): Promise<UpdateGoogleProviderResponse.AsObject> {
    return this.grpcService.admin.updateGoogleProvider(req, null).then((resp) => resp.toObject());
  }

  public addLDAPProvider(req: AddLDAPProviderRequest): Promise<AddLDAPProviderResponse.AsObject> {
    return this.grpcService.admin.addLDAPProvider(req, null).then((resp) => resp.toObject());
  }

  public updateLDAPProvider(req: UpdateLDAPProviderRequest): Promise<UpdateLDAPProviderResponse.AsObject> {
    return this.grpcService.admin.updateLDAPProvider(req, null).then((resp) => resp.toObject());
  }

  public addGitLabProvider(req: AddGitLabProviderRequest): Promise<AddGitLabProviderResponse.AsObject> {
    return this.grpcService.admin.addGitLabProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGitLabProvider(req: UpdateGitLabProviderRequest): Promise<UpdateGitLabProviderResponse.AsObject> {
    return this.grpcService.admin.updateGitLabProvider(req, null).then((resp) => resp.toObject());
  }

  public addGitLabSelfHostedProvider(
    req: AddGitLabSelfHostedProviderRequest,
  ): Promise<AddGitLabSelfHostedProviderResponse.AsObject> {
    return this.grpcService.admin.addGitLabSelfHostedProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGitLabSelfHostedProvider(
    req: UpdateGitLabSelfHostedProviderRequest,
  ): Promise<UpdateGitLabSelfHostedProviderResponse.AsObject> {
    return this.grpcService.admin.updateGitLabSelfHostedProvider(req, null).then((resp) => resp.toObject());
  }

  public addGitHubProvider(req: AddGitHubProviderRequest): Promise<AddGitHubProviderResponse.AsObject> {
    return this.grpcService.admin.addGitHubProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGitHubProvider(req: UpdateGitHubProviderRequest): Promise<UpdateGitHubProviderResponse.AsObject> {
    return this.grpcService.admin.updateGitHubProvider(req, null).then((resp) => resp.toObject());
  }

  public addGenericOIDCProvider(req: AddGenericOIDCProviderRequest): Promise<AddGenericOIDCProviderResponse.AsObject> {
    return this.grpcService.admin.addGenericOIDCProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGenericOIDCProvider(
    req: UpdateGenericOIDCProviderRequest,
  ): Promise<UpdateGenericOIDCProviderResponse.AsObject> {
    return this.grpcService.admin.updateGenericOIDCProvider(req, null).then((resp) => resp.toObject());
  }

  public addGenericOAuthProvider(req: AddGenericOAuthProviderRequest): Promise<AddGenericOAuthProviderResponse.AsObject> {
    return this.grpcService.admin.addGenericOAuthProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGenericOAuthProvider(
    req: UpdateGenericOAuthProviderRequest,
  ): Promise<UpdateGenericOAuthProviderResponse.AsObject> {
    return this.grpcService.admin.updateGenericOAuthProvider(req, null).then((resp) => resp.toObject());
  }

  public addJWTProvider(req: AddJWTProviderRequest): Promise<AddJWTProviderResponse.AsObject> {
    return this.grpcService.admin.addJWTProvider(req, null).then((resp) => resp.toObject());
  }

  public updateJWTProvider(req: UpdateJWTProviderRequest): Promise<UpdateJWTProviderResponse.AsObject> {
    return this.grpcService.admin.updateJWTProvider(req, null).then((resp) => resp.toObject());
  }

  public addSAMLProvider(req: AddSAMLProviderRequest): Promise<AddSAMLProviderResponse.AsObject> {
    return this.grpcService.admin.addSAMLProvider(req, null).then((resp) => resp.toObject());
  }

  public updateSAMLProvider(req: UpdateSAMLProviderRequest): Promise<UpdateSAMLProviderResponse.AsObject> {
    return this.grpcService.admin.updateSAMLProvider(req, null).then((resp) => resp.toObject());
  }

  public addGitHubEnterpriseServerProvider(
    req: AddGitHubEnterpriseServerProviderRequest,
  ): Promise<AddGitHubEnterpriseServerProviderResponse.AsObject> {
    return this.grpcService.admin.addGitHubEnterpriseServerProvider(req, null).then((resp) => resp.toObject());
  }

  public updateGitHubEnterpriseServerProvider(
    req: UpdateGitHubEnterpriseServerProviderRequest,
  ): Promise<UpdateGitHubEnterpriseServerProviderResponse.AsObject> {
    return this.grpcService.admin.updateGitHubEnterpriseServerProvider(req, null).then((resp) => resp.toObject());
  }

  public addAppleProvider(req: AddAppleProviderRequest): Promise<AddAppleProviderResponse.AsObject> {
    return this.grpcService.admin.addAppleProvider(req, null).then((resp) => resp.toObject());
  }

  public updateAppleProvider(req: UpdateAppleProviderRequest): Promise<UpdateAppleProviderResponse.AsObject> {
    return this.grpcService.admin.updateAppleProvider(req, null).then((resp) => resp.toObject());
  }

  public deleteProvider(id: string): Promise<DeleteProviderResponse.AsObject> {
    const req = new DeleteProviderRequest();
    req.setId(id);
    return this.grpcService.admin.deleteProvider(req, null).then((resp) => resp.toObject());
  }

  public listProviders(req: ListProvidersRequest): Promise<ListProvidersResponse.AsObject> {
    return this.grpcService.admin.listProviders(req, null).then((resp) => resp.toObject());
  }

  public getProviderByID(req: GetProviderByIDRequest): Promise<GetProviderByIDResponse.AsObject> {
    return this.grpcService.admin.getProviderByID(req, null).then((resp) => resp.toObject());
  }

  public getProviderID(req: GetProviderByIDRequest): Promise<GetProviderByIDResponse> {
    return this.grpcService.admin.getProviderByID(req, null);
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

    return this.grpcService.admin.listIAMMembers(req, null).then((resp) => resp.toObject());
  }

  public removeIAMMember(userId: string): Promise<RemoveIAMMemberResponse.AsObject> {
    const req = new RemoveIAMMemberRequest();
    req.setUserId(userId);
    return this.grpcService.admin.removeIAMMember(req, null).then((resp) => resp.toObject());
  }

  public addIAMMember(userId: string, rolesList: string[]): Promise<AddIAMMemberResponse.AsObject> {
    const req = new AddIAMMemberRequest();
    req.setUserId(userId);
    req.setRolesList(rolesList);

    return this.grpcService.admin.addIAMMember(req, null).then((resp) => resp.toObject());
  }

  public updateIAMMember(userId: string, rolesList: string[]): Promise<UpdateIAMMemberResponse.AsObject> {
    const req = new UpdateIAMMemberRequest();
    req.setUserId(userId);
    req.setRolesList(rolesList);

    return this.grpcService.admin.updateIAMMember(req, null).then((resp) => resp.toObject());
  }

  public listMilestones(req: ListMilestonesRequest): Promise<ListMilestonesResponse.AsObject> {
    return this.grpcService.admin.listMilestones(req, null).then((resp) => resp.toObject());
  }

  public listOrgs(
    limit: number,
    offset: number,
    queriesList?: OrgQuery[],
    sortingColumn?: OrgFieldName,
    sortingDirection?: SortDirection,
  ): Promise<ListOrgsResponse.AsObject> {
    const req = new ListOrgsRequest();
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
    return this.grpcService.admin.listOrgs(req, null).then((resp) => resp.toObject());
  }
}

import { Injectable } from '@angular/core';
import { BehaviorSubject, catchError, from, map, Observable, of, Subject, switchMap, tap } from 'rxjs';

import {
  ActivateLabelPolicyRequest,
  ActivateLabelPolicyResponse,
  ActivateSMSProviderRequest,
  ActivateSMSProviderResponse,
  AddCustomDomainPolicyRequest,
  AddCustomOrgIAMPolicyResponse,
  AddIAMMemberRequest,
  AddIAMMemberResponse,
  AddIDPToLoginPolicyRequest,
  AddIDPToLoginPolicyResponse,
  AddJWTIDPRequest,
  AddJWTIDPResponse,
  AddMultiFactorToLoginPolicyRequest,
  AddMultiFactorToLoginPolicyResponse,
  AddOIDCIDPRequest,
  AddOIDCIDPResponse,
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
  GetCustomDomainClaimedMessageTextRequest,
  GetCustomDomainClaimedMessageTextResponse,
  GetCustomDomainPolicyRequest,
  GetCustomDomainPolicyResponse,
  GetCustomInitMessageTextRequest,
  GetCustomInitMessageTextResponse,
  GetCustomLoginTextsRequest,
  GetCustomLoginTextsResponse,
  GetCustomPasswordlessRegistrationMessageTextRequest,
  GetCustomPasswordlessRegistrationMessageTextResponse,
  GetCustomPasswordResetMessageTextRequest,
  GetCustomPasswordResetMessageTextResponse,
  GetCustomVerifyEmailMessageTextRequest,
  GetCustomVerifyEmailMessageTextResponse,
  GetCustomVerifyPhoneMessageTextRequest,
  GetCustomVerifyPhoneMessageTextResponse,
  GetDefaultDomainClaimedMessageTextRequest,
  GetDefaultDomainClaimedMessageTextResponse,
  GetDefaultInitMessageTextRequest,
  GetDefaultInitMessageTextResponse,
  GetDefaultLanguageRequest,
  GetDefaultLanguageResponse,
  GetDefaultLoginTextsRequest,
  GetDefaultLoginTextsResponse,
  GetDefaultPasswordlessRegistrationMessageTextRequest,
  GetDefaultPasswordlessRegistrationMessageTextResponse,
  GetDefaultPasswordResetMessageTextRequest,
  GetDefaultPasswordResetMessageTextResponse,
  GetDefaultVerifyEmailMessageTextRequest,
  GetDefaultVerifyEmailMessageTextResponse,
  GetDefaultVerifyPhoneMessageTextRequest,
  GetDefaultVerifyPhoneMessageTextResponse,
  GetDomainPolicyRequest,
  GetDomainPolicyResponse,
  GetFileSystemNotificationProviderRequest,
  GetFileSystemNotificationProviderResponse,
  GetIDPByIDRequest,
  GetIDPByIDResponse,
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
  GetSecretGeneratorRequest,
  GetSecretGeneratorResponse,
  GetSMSProviderRequest,
  GetSMSProviderResponse,
  GetSMTPConfigRequest,
  GetSMTPConfigResponse,
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
  ListSecretGeneratorsRequest,
  ListSecretGeneratorsResponse,
  ListSMSProvidersRequest,
  ListSMSProvidersResponse,
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
  RemoveSMSProviderRequest,
  RemoveSMSProviderResponse,
  ResetCustomDomainPolicyToDefaultRequest,
  ResetCustomDomainPolicyToDefaultResponse,
  ResetCustomLoginTextsToDefaultRequest,
  ResetCustomLoginTextsToDefaultResponse,
  SetCustomLoginTextsRequest,
  SetCustomLoginTextsResponse,
  SetDefaultDomainClaimedMessageTextRequest,
  SetDefaultDomainClaimedMessageTextResponse,
  SetDefaultInitMessageTextRequest,
  SetDefaultInitMessageTextResponse,
  SetDefaultLanguageRequest,
  SetDefaultLanguageResponse,
  SetDefaultPasswordlessRegistrationMessageTextRequest,
  SetDefaultPasswordlessRegistrationMessageTextResponse,
  SetDefaultPasswordResetMessageTextRequest,
  SetDefaultPasswordResetMessageTextResponse,
  SetDefaultVerifyEmailMessageTextRequest,
  SetDefaultVerifyEmailMessageTextResponse,
  SetDefaultVerifyPhoneMessageTextRequest,
  SetDefaultVerifyPhoneMessageTextResponse,
  SetUpOrgRequest,
  SetUpOrgResponse,
  UpdateCustomDomainPolicyRequest,
  UpdateCustomDomainPolicyResponse,
  UpdateDomainPolicyRequest,
  UpdateDomainPolicyResponse,
  UpdateIAMMemberRequest,
  UpdateIAMMemberResponse,
  UpdateIDPJWTConfigRequest,
  UpdateIDPJWTConfigResponse,
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
  AddOIDCSettingsRequest,
  AddOIDCSettingsResponse,
  UpdateOIDCSettingsRequest,
  UpdateOIDCSettingsResponse,
  UpdatePasswordAgePolicyRequest,
  UpdatePasswordAgePolicyResponse,
  UpdatePasswordComplexityPolicyRequest,
  UpdatePasswordComplexityPolicyResponse,
  UpdatePrivacyPolicyRequest,
  UpdatePrivacyPolicyResponse,
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
  GetSecurityPolicyRequest,
  GetSecurityPolicyResponse,
  SetSecurityPolicyRequest,
  SetSecurityPolicyResponse,
  ListEventsResponse,
  ListEventsRequest,
  ListEventTypesRequest,
  ListEventTypesResponse,
  ListAggregateTypesRequest,
  ListAggregateTypesResponse,
  GetNotificationPolicyRequest,
  GetNotificationPolicyResponse,
  UpdateNotificationPolicyRequest,
  UpdateNotificationPolicyResponse,
  GetDefaultPasswordChangeMessageTextResponse,
  GetDefaultPasswordChangeMessageTextRequest,
  GetCustomPasswordChangeMessageTextResponse,
  SetDefaultPasswordChangeMessageTextRequest,
  SetDefaultPasswordChangeMessageTextResponse,
  GetCustomPasswordChangeMessageTextRequest,
  AddNotificationPolicyRequest,
  AddNotificationPolicyResponse,
  SetDefaultOrgRequest,
  SetDefaultOrgResponse,
} from '../proto/generated/zitadel/admin_pb';
import { Event } from '../proto/generated/zitadel/event_pb';
import { SearchQuery } from '../proto/generated/zitadel/member_pb';
import { ListQuery } from '../proto/generated/zitadel/object_pb';
import { GrpcService } from './grpc.service';

interface OnboardingActions {
  order: number;
  eventType: string;
  link: string | string[];
  fragment?: string | undefined;
}

type OnboardingEvent = { order: number; link: string; fragment: string | undefined; event: Event.AsObject | undefined };
type OnboardingEventEntries = Array<[string, OnboardingEvent]> | [];

const ONBOARDING_EVENTS: OnboardingActions[] = [
  { order: 0, eventType: 'instance.policy.label.added', link: ['/settings'], fragment: 'branding' },
  { order: 1, eventType: 'project.added', link: ['/projects/create'] },
  { order: 2, eventType: 'project.application.added', link: ['/projects/app-create'] },
  { order: 3, eventType: 'user.human.added', link: ['/users/create'] },
  { order: 4, eventType: 'instance.policy.notification.added', link: ['/settings'], fragment: 'notifications' },
  { order: 5, eventType: 'user.grant.added', link: ['/grant-create'] },
  //   { eventType: 'org.policy.notification.added', link: '/settings?id=notifications' },
];

@Injectable({
  providedIn: 'root',
})
export class AdminService {
  private loadEvents: Subject<string[]> = new Subject();

  public progressEvents$: Observable<OnboardingEventEntries> = this.loadEvents.pipe(
    switchMap((types) => {
      const eventsReq = new ListEventsRequest().setEventTypesList(types).setAsc(false);
      return from(this.listEvents(eventsReq)).pipe(
        map((events) => {
          const eventList = events.getEventsList().map((event) => event.toObject());

          let obj: { [type: string]: OnboardingEvent } = {};
          types.map((type) => {
            const filtered = eventList.filter((event) => event.type?.type === type);
            const el = ONBOARDING_EVENTS.find((oe) => oe.eventType === type);
            (obj as any)[type] = filtered.length
              ? { order: el?.order, link: el?.link, fragment: el?.fragment, event: filtered[0] }
              : { order: el?.order, link: el?.link, fragment: el?.fragment, event: undefined };
          });

          const toArray = Object.entries(obj).sort(([key0, a], [key1, b]) => a.order - b.order);

          const toDo = toArray.filter(([key, value]) => value.event === undefined);
          const done = toArray.filter(([key, value]) => !!value.event);

          return [...toDo, ...done];
        }),
        tap((events) => {
          const total = events.length;
          const done = events.map(([type, value]) => value.event !== undefined).filter((res) => !!res).length;
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
      );
    }),
  );

  public progressEvents: BehaviorSubject<OnboardingEventEntries> = new BehaviorSubject<OnboardingEventEntries>([]);
  public progressPercentage: BehaviorSubject<number> = new BehaviorSubject(0);
  public progressDone: BehaviorSubject<number> = new BehaviorSubject(0);
  public progressTotal: BehaviorSubject<number> = new BehaviorSubject(0);
  public progressAllDone: BehaviorSubject<boolean> = new BehaviorSubject(false);

  constructor(private readonly grpcService: GrpcService) {
    this.progressEvents$.subscribe(this.progressEvents);
    this.initOnboarding();
  }

  public initOnboarding(): void {
    console.log('init onboarding');
    this.loadEvents.next(ONBOARDING_EVENTS.map((oe) => oe.eventType));
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

  public addSMTPConfig(req: AddSMTPConfigRequest): Promise<AddSMTPConfigResponse.AsObject> {
    return this.grpcService.admin.addSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  public updateSMTPConfig(req: UpdateSMTPConfigRequest): Promise<UpdateSMTPConfigResponse.AsObject> {
    return this.grpcService.admin.updateSMTPConfig(req, null).then((resp) => resp.toObject());
  }

  public updateSMTPConfigPassword(req: UpdateSMTPConfigPasswordRequest): Promise<UpdateSMTPConfigPasswordResponse.AsObject> {
    return this.grpcService.admin.updateSMTPConfigPassword(req, null).then((resp) => resp.toObject());
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

  public updateLockoutPolicy(maxAttempts: number): Promise<UpdateLockoutPolicyResponse.AsObject> {
    const req = new UpdateLockoutPolicyRequest();
    req.setMaxPasswordAttempts(maxAttempts);

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

  public listIDPs(limit?: number, offset?: number, queriesList?: IDPQuery[]): Promise<ListIDPsResponse.AsObject> {
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
    return this.grpcService.admin.listIDPs(req, null).then((resp) => resp.toObject());
  }

  public getIDPByID(id: string): Promise<GetIDPByIDResponse.AsObject> {
    const req = new GetIDPByIDRequest();
    req.setId(id);
    return this.grpcService.admin.getIDPByID(req, null).then((resp) => resp.toObject());
  }

  public updateIDP(req: UpdateIDPRequest): Promise<UpdateIDPResponse.AsObject> {
    return this.grpcService.admin.updateIDP(req, null).then((resp) => resp.toObject());
  }

  public addOIDCIDP(req: AddOIDCIDPRequest): Promise<AddOIDCIDPResponse.AsObject> {
    return this.grpcService.admin.addOIDCIDP(req, null).then((resp) => resp.toObject());
  }

  public updateIDPOIDCConfig(req: UpdateIDPOIDCConfigRequest): Promise<UpdateIDPOIDCConfigResponse.AsObject> {
    return this.grpcService.admin.updateIDPOIDCConfig(req, null).then((resp) => resp.toObject());
  }

  public removeIDP(id: string): Promise<RemoveIDPResponse.AsObject> {
    const req = new RemoveIDPRequest();
    req.setIdpId(id);
    return this.grpcService.admin.removeIDP(req, null).then((resp) => resp.toObject());
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

  public addJWTIDP(req: AddJWTIDPRequest): Promise<AddJWTIDPResponse.AsObject> {
    return this.grpcService.admin.addJWTIDP(req, null).then((resp) => resp.toObject());
  }

  public updateIDPJWTConfig(req: UpdateIDPJWTConfigRequest): Promise<UpdateIDPJWTConfigResponse.AsObject> {
    return this.grpcService.admin.updateIDPJWTConfig(req, null).then((resp) => resp.toObject());
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
}

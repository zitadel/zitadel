---
title: Admin API
description: …
---

### Administration aka Admin

This API is intended to configure and manage the IAM itself.

| Service | URI                                                                                                                             |
|:--------|:--------------------------------------------------------------------------------------------------------------------------------|
| REST    | [https://api.zitadel.ch/admin/v1/](https://api.zitadel.ch/admin/v1/)                                                            |
| GRPC    | [https://api.zitadel.ch/caos.zitadel.admin.api.v1.AdminService/](https://api.zitadel.ch/caos.zitadel.admin.api.v1.AdminService) |

[Latest API Version](https://github.com/caos/zitadel/blob/main/proto/zitadel/admin.proto)


### Table of Contents


- Services
  - [AdminService](#zitadeladminv1adminservice)



- Messages
  - [AddCustomOrgIAMPolicyRequest](#addcustomorgiampolicyrequest)
  - [AddCustomOrgIAMPolicyResponse](#addcustomorgiampolicyresponse)
  - [AddIAMMemberRequest](#addiammemberrequest)
  - [AddIAMMemberResponse](#addiammemberresponse)
  - [AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest)
  - [AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)
  - [AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest)
  - [AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)
  - [AddOIDCIDPRequest](#addoidcidprequest)
  - [AddOIDCIDPResponse](#addoidcidpresponse)
  - [AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest)
  - [AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)
  - [ClearViewRequest](#clearviewrequest)
  - [ClearViewResponse](#clearviewresponse)
  - [DeactivateIDPRequest](#deactivateidprequest)
  - [DeactivateIDPResponse](#deactivateidpresponse)
  - [FailedEvent](#failedevent)
  - [GetCustomOrgIAMPolicyRequest](#getcustomorgiampolicyrequest)
  - [GetCustomOrgIAMPolicyResponse](#getcustomorgiampolicyresponse)
  - [GetDefaultFeaturesRequest](#getdefaultfeaturesrequest)
  - [GetDefaultFeaturesResponse](#getdefaultfeaturesresponse)
  - [GetIDPByIDRequest](#getidpbyidrequest)
  - [GetIDPByIDResponse](#getidpbyidresponse)
  - [GetLabelPolicyRequest](#getlabelpolicyrequest)
  - [GetLabelPolicyResponse](#getlabelpolicyresponse)
  - [GetLoginPolicyRequest](#getloginpolicyrequest)
  - [GetLoginPolicyResponse](#getloginpolicyresponse)
  - [GetOrgByIDRequest](#getorgbyidrequest)
  - [GetOrgByIDResponse](#getorgbyidresponse)
  - [GetOrgFeaturesRequest](#getorgfeaturesrequest)
  - [GetOrgFeaturesResponse](#getorgfeaturesresponse)
  - [GetOrgIAMPolicyRequest](#getorgiampolicyrequest)
  - [GetOrgIAMPolicyResponse](#getorgiampolicyresponse)
  - [GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest)
  - [GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)
  - [GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest)
  - [GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)
  - [GetPasswordLockoutPolicyRequest](#getpasswordlockoutpolicyrequest)
  - [GetPasswordLockoutPolicyResponse](#getpasswordlockoutpolicyresponse)
  - [HealthzRequest](#healthzrequest)
  - [HealthzResponse](#healthzresponse)
  - [IDPQuery](#idpquery)
  - [IsOrgUniqueRequest](#isorguniquerequest)
  - [IsOrgUniqueResponse](#isorguniqueresponse)
  - [ListFailedEventsRequest](#listfailedeventsrequest)
  - [ListFailedEventsResponse](#listfailedeventsresponse)
  - [ListIAMMemberRolesRequest](#listiammemberrolesrequest)
  - [ListIAMMemberRolesResponse](#listiammemberrolesresponse)
  - [ListIAMMembersRequest](#listiammembersrequest)
  - [ListIAMMembersResponse](#listiammembersresponse)
  - [ListIDPsRequest](#listidpsrequest)
  - [ListIDPsResponse](#listidpsresponse)
  - [ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest)
  - [ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)
  - [ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest)
  - [ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)
  - [ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest)
  - [ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)
  - [ListOrgsRequest](#listorgsrequest)
  - [ListOrgsResponse](#listorgsresponse)
  - [ListViewsRequest](#listviewsrequest)
  - [ListViewsResponse](#listviewsresponse)
  - [ReactivateIDPRequest](#reactivateidprequest)
  - [ReactivateIDPResponse](#reactivateidpresponse)
  - [RemoveFailedEventRequest](#removefailedeventrequest)
  - [RemoveFailedEventResponse](#removefailedeventresponse)
  - [RemoveIAMMemberRequest](#removeiammemberrequest)
  - [RemoveIAMMemberResponse](#removeiammemberresponse)
  - [RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest)
  - [RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)
  - [RemoveIDPRequest](#removeidprequest)
  - [RemoveIDPResponse](#removeidpresponse)
  - [RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest)
  - [RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)
  - [RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest)
  - [RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)
  - [ResetCustomOrgIAMPolicyToDefaultRequest](#resetcustomorgiampolicytodefaultrequest)
  - [ResetCustomOrgIAMPolicyToDefaultResponse](#resetcustomorgiampolicytodefaultresponse)
  - [ResetOrgFeaturesRequest](#resetorgfeaturesrequest)
  - [ResetOrgFeaturesResponse](#resetorgfeaturesresponse)
  - [SetDefaultFeaturesRequest](#setdefaultfeaturesrequest)
  - [SetDefaultFeaturesResponse](#setdefaultfeaturesresponse)
  - [SetOrgFeaturesRequest](#setorgfeaturesrequest)
  - [SetOrgFeaturesResponse](#setorgfeaturesresponse)
  - [SetUpOrgRequest](#setuporgrequest)
  - [SetUpOrgRequest.Human](#setuporgrequesthuman)
  - [SetUpOrgRequest.Human.Email](#setuporgrequesthumanemail)
  - [SetUpOrgRequest.Human.Phone](#setuporgrequesthumanphone)
  - [SetUpOrgRequest.Human.Profile](#setuporgrequesthumanprofile)
  - [SetUpOrgRequest.Org](#setuporgrequestorg)
  - [SetUpOrgResponse](#setuporgresponse)
  - [UpdateCustomOrgIAMPolicyRequest](#updatecustomorgiampolicyrequest)
  - [UpdateCustomOrgIAMPolicyResponse](#updatecustomorgiampolicyresponse)
  - [UpdateIAMMemberRequest](#updateiammemberrequest)
  - [UpdateIAMMemberResponse](#updateiammemberresponse)
  - [UpdateIDPOIDCConfigRequest](#updateidpoidcconfigrequest)
  - [UpdateIDPOIDCConfigResponse](#updateidpoidcconfigresponse)
  - [UpdateIDPRequest](#updateidprequest)
  - [UpdateIDPResponse](#updateidpresponse)
  - [UpdateLabelPolicyRequest](#updatelabelpolicyrequest)
  - [UpdateLabelPolicyResponse](#updatelabelpolicyresponse)
  - [UpdateLoginPolicyRequest](#updateloginpolicyrequest)
  - [UpdateLoginPolicyResponse](#updateloginpolicyresponse)
  - [UpdateOrgIAMPolicyRequest](#updateorgiampolicyrequest)
  - [UpdateOrgIAMPolicyResponse](#updateorgiampolicyresponse)
  - [UpdatePasswordAgePolicyRequest](#updatepasswordagepolicyrequest)
  - [UpdatePasswordAgePolicyResponse](#updatepasswordagepolicyresponse)
  - [UpdatePasswordComplexityPolicyRequest](#updatepasswordcomplexitypolicyrequest)
  - [UpdatePasswordComplexityPolicyResponse](#updatepasswordcomplexitypolicyresponse)
  - [UpdatePasswordLockoutPolicyRequest](#updatepasswordlockoutpolicyrequest)
  - [UpdatePasswordLockoutPolicyResponse](#updatepasswordlockoutpolicyresponse)
  - [View](#view)






- Messages
  - [APIConfig](#apiconfig)
  - [App](#app)
  - [AppNameQuery](#appnamequery)
  - [AppQuery](#appquery)
  - [OIDCConfig](#oidcconfig)



- Enums
  - [APIAuthMethodType](#apiauthmethodtype)
  - [AppState](#appstate)
  - [OIDCAppType](#oidcapptype)
  - [OIDCAuthMethodType](#oidcauthmethodtype)
  - [OIDCGrantType](#oidcgranttype)
  - [OIDCResponseType](#oidcresponsetype)
  - [OIDCTokenType](#oidctokentype)
  - [OIDCVersion](#oidcversion)




- Services
  - [AuthService](#zitadelauthv1authservice)



- Messages
  - [AddMyAuthFactorOTPRequest](#addmyauthfactorotprequest)
  - [AddMyAuthFactorOTPResponse](#addmyauthfactorotpresponse)
  - [AddMyAuthFactorU2FRequest](#addmyauthfactoru2frequest)
  - [AddMyAuthFactorU2FResponse](#addmyauthfactoru2fresponse)
  - [AddMyPasswordlessRequest](#addmypasswordlessrequest)
  - [AddMyPasswordlessResponse](#addmypasswordlessresponse)
  - [GetMyEmailRequest](#getmyemailrequest)
  - [GetMyEmailResponse](#getmyemailresponse)
  - [GetMyPasswordComplexityPolicyRequest](#getmypasswordcomplexitypolicyrequest)
  - [GetMyPasswordComplexityPolicyResponse](#getmypasswordcomplexitypolicyresponse)
  - [GetMyPhoneRequest](#getmyphonerequest)
  - [GetMyPhoneResponse](#getmyphoneresponse)
  - [GetMyProfileRequest](#getmyprofilerequest)
  - [GetMyProfileResponse](#getmyprofileresponse)
  - [GetMyUserRequest](#getmyuserrequest)
  - [GetMyUserResponse](#getmyuserresponse)
  - [HealthzRequest](#healthzrequest)
  - [HealthzResponse](#healthzresponse)
  - [ListMyAuthFactorsRequest](#listmyauthfactorsrequest)
  - [ListMyAuthFactorsResponse](#listmyauthfactorsresponse)
  - [ListMyLinkedIDPsRequest](#listmylinkedidpsrequest)
  - [ListMyLinkedIDPsResponse](#listmylinkedidpsresponse)
  - [ListMyPasswordlessRequest](#listmypasswordlessrequest)
  - [ListMyPasswordlessResponse](#listmypasswordlessresponse)
  - [ListMyProjectOrgsRequest](#listmyprojectorgsrequest)
  - [ListMyProjectOrgsResponse](#listmyprojectorgsresponse)
  - [ListMyProjectPermissionsRequest](#listmyprojectpermissionsrequest)
  - [ListMyProjectPermissionsResponse](#listmyprojectpermissionsresponse)
  - [ListMyUserChangesRequest](#listmyuserchangesrequest)
  - [ListMyUserChangesResponse](#listmyuserchangesresponse)
  - [ListMyUserGrantsRequest](#listmyusergrantsrequest)
  - [ListMyUserGrantsResponse](#listmyusergrantsresponse)
  - [ListMyUserSessionsRequest](#listmyusersessionsrequest)
  - [ListMyUserSessionsResponse](#listmyusersessionsresponse)
  - [ListMyZitadelFeaturesRequest](#listmyzitadelfeaturesrequest)
  - [ListMyZitadelFeaturesResponse](#listmyzitadelfeaturesresponse)
  - [ListMyZitadelPermissionsRequest](#listmyzitadelpermissionsrequest)
  - [ListMyZitadelPermissionsResponse](#listmyzitadelpermissionsresponse)
  - [RemoveMyAuthFactorOTPRequest](#removemyauthfactorotprequest)
  - [RemoveMyAuthFactorOTPResponse](#removemyauthfactorotpresponse)
  - [RemoveMyAuthFactorU2FRequest](#removemyauthfactoru2frequest)
  - [RemoveMyAuthFactorU2FResponse](#removemyauthfactoru2fresponse)
  - [RemoveMyLinkedIDPRequest](#removemylinkedidprequest)
  - [RemoveMyLinkedIDPResponse](#removemylinkedidpresponse)
  - [RemoveMyPasswordlessRequest](#removemypasswordlessrequest)
  - [RemoveMyPasswordlessResponse](#removemypasswordlessresponse)
  - [RemoveMyPhoneRequest](#removemyphonerequest)
  - [RemoveMyPhoneResponse](#removemyphoneresponse)
  - [ResendMyEmailVerificationRequest](#resendmyemailverificationrequest)
  - [ResendMyEmailVerificationResponse](#resendmyemailverificationresponse)
  - [ResendMyPhoneVerificationRequest](#resendmyphoneverificationrequest)
  - [ResendMyPhoneVerificationResponse](#resendmyphoneverificationresponse)
  - [SetMyEmailRequest](#setmyemailrequest)
  - [SetMyEmailResponse](#setmyemailresponse)
  - [SetMyPhoneRequest](#setmyphonerequest)
  - [SetMyPhoneResponse](#setmyphoneresponse)
  - [UpdateMyPasswordRequest](#updatemypasswordrequest)
  - [UpdateMyPasswordResponse](#updatemypasswordresponse)
  - [UpdateMyProfileRequest](#updatemyprofilerequest)
  - [UpdateMyProfileResponse](#updatemyprofileresponse)
  - [UpdateMyUserNameRequest](#updatemyusernamerequest)
  - [UpdateMyUserNameResponse](#updatemyusernameresponse)
  - [UserGrant](#usergrant)
  - [VerifyMyAuthFactorOTPRequest](#verifymyauthfactorotprequest)
  - [VerifyMyAuthFactorOTPResponse](#verifymyauthfactorotpresponse)
  - [VerifyMyAuthFactorU2FRequest](#verifymyauthfactoru2frequest)
  - [VerifyMyAuthFactorU2FResponse](#verifymyauthfactoru2fresponse)
  - [VerifyMyEmailRequest](#verifymyemailrequest)
  - [VerifyMyEmailResponse](#verifymyemailresponse)
  - [VerifyMyPasswordlessRequest](#verifymypasswordlessrequest)
  - [VerifyMyPasswordlessResponse](#verifymypasswordlessresponse)
  - [VerifyMyPhoneRequest](#verifymyphonerequest)
  - [VerifyMyPhoneResponse](#verifymyphoneresponse)






- Messages
  - [Key](#key)



- Enums
  - [KeyType](#keytype)





- Messages
  - [Change](#change)
  - [ChangeQuery](#changequery)






- Messages
  - [FeatureTier](#featuretier)
  - [Features](#features)



- Enums
  - [FeaturesState](#featuresstate)





- Messages
  - [IDP](#idp)
  - [IDPIDQuery](#idpidquery)
  - [IDPLoginPolicyLink](#idploginpolicylink)
  - [IDPNameQuery](#idpnamequery)
  - [IDPOwnerTypeQuery](#idpownertypequery)
  - [IDPUserLink](#idpuserlink)
  - [OIDCConfig](#oidcconfig)



- Enums
  - [IDPFieldName](#idpfieldname)
  - [IDPOwnerType](#idpownertype)
  - [IDPState](#idpstate)
  - [IDPStylingType](#idpstylingtype)
  - [IDPType](#idptype)
  - [OIDCMappingField](#oidcmappingfield)




- Services
  - [ManagementService](#zitadelmanagementv1managementservice)



- Messages
  - [AddAPIAppRequest](#addapiapprequest)
  - [AddAPIAppResponse](#addapiappresponse)
  - [AddAppKeyRequest](#addappkeyrequest)
  - [AddAppKeyResponse](#addappkeyresponse)
  - [AddCustomLabelPolicyRequest](#addcustomlabelpolicyrequest)
  - [AddCustomLabelPolicyResponse](#addcustomlabelpolicyresponse)
  - [AddCustomLoginPolicyRequest](#addcustomloginpolicyrequest)
  - [AddCustomLoginPolicyResponse](#addcustomloginpolicyresponse)
  - [AddCustomPasswordAgePolicyRequest](#addcustompasswordagepolicyrequest)
  - [AddCustomPasswordAgePolicyResponse](#addcustompasswordagepolicyresponse)
  - [AddCustomPasswordComplexityPolicyRequest](#addcustompasswordcomplexitypolicyrequest)
  - [AddCustomPasswordComplexityPolicyResponse](#addcustompasswordcomplexitypolicyresponse)
  - [AddCustomPasswordLockoutPolicyRequest](#addcustompasswordlockoutpolicyrequest)
  - [AddCustomPasswordLockoutPolicyResponse](#addcustompasswordlockoutpolicyresponse)
  - [AddHumanUserRequest](#addhumanuserrequest)
  - [AddHumanUserRequest.Email](#addhumanuserrequestemail)
  - [AddHumanUserRequest.Phone](#addhumanuserrequestphone)
  - [AddHumanUserRequest.Profile](#addhumanuserrequestprofile)
  - [AddHumanUserResponse](#addhumanuserresponse)
  - [AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest)
  - [AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)
  - [AddMachineKeyRequest](#addmachinekeyrequest)
  - [AddMachineKeyResponse](#addmachinekeyresponse)
  - [AddMachineUserRequest](#addmachineuserrequest)
  - [AddMachineUserResponse](#addmachineuserresponse)
  - [AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest)
  - [AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)
  - [AddOIDCAppRequest](#addoidcapprequest)
  - [AddOIDCAppResponse](#addoidcappresponse)
  - [AddOrgDomainRequest](#addorgdomainrequest)
  - [AddOrgDomainResponse](#addorgdomainresponse)
  - [AddOrgMemberRequest](#addorgmemberrequest)
  - [AddOrgMemberResponse](#addorgmemberresponse)
  - [AddOrgOIDCIDPRequest](#addorgoidcidprequest)
  - [AddOrgOIDCIDPResponse](#addorgoidcidpresponse)
  - [AddOrgRequest](#addorgrequest)
  - [AddOrgResponse](#addorgresponse)
  - [AddProjectGrantMemberRequest](#addprojectgrantmemberrequest)
  - [AddProjectGrantMemberResponse](#addprojectgrantmemberresponse)
  - [AddProjectGrantRequest](#addprojectgrantrequest)
  - [AddProjectGrantResponse](#addprojectgrantresponse)
  - [AddProjectMemberRequest](#addprojectmemberrequest)
  - [AddProjectMemberResponse](#addprojectmemberresponse)
  - [AddProjectRequest](#addprojectrequest)
  - [AddProjectResponse](#addprojectresponse)
  - [AddProjectRoleRequest](#addprojectrolerequest)
  - [AddProjectRoleResponse](#addprojectroleresponse)
  - [AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest)
  - [AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)
  - [AddUserGrantRequest](#addusergrantrequest)
  - [AddUserGrantResponse](#addusergrantresponse)
  - [BulkAddProjectRolesRequest](#bulkaddprojectrolesrequest)
  - [BulkAddProjectRolesRequest.Role](#bulkaddprojectrolesrequestrole)
  - [BulkAddProjectRolesResponse](#bulkaddprojectrolesresponse)
  - [BulkRemoveUserGrantRequest](#bulkremoveusergrantrequest)
  - [BulkRemoveUserGrantResponse](#bulkremoveusergrantresponse)
  - [DeactivateAppRequest](#deactivateapprequest)
  - [DeactivateAppResponse](#deactivateappresponse)
  - [DeactivateOrgIDPRequest](#deactivateorgidprequest)
  - [DeactivateOrgIDPResponse](#deactivateorgidpresponse)
  - [DeactivateOrgRequest](#deactivateorgrequest)
  - [DeactivateOrgResponse](#deactivateorgresponse)
  - [DeactivateProjectGrantRequest](#deactivateprojectgrantrequest)
  - [DeactivateProjectGrantResponse](#deactivateprojectgrantresponse)
  - [DeactivateProjectRequest](#deactivateprojectrequest)
  - [DeactivateProjectResponse](#deactivateprojectresponse)
  - [DeactivateUserGrantRequest](#deactivateusergrantrequest)
  - [DeactivateUserGrantResponse](#deactivateusergrantresponse)
  - [DeactivateUserRequest](#deactivateuserrequest)
  - [DeactivateUserResponse](#deactivateuserresponse)
  - [GenerateOrgDomainValidationRequest](#generateorgdomainvalidationrequest)
  - [GenerateOrgDomainValidationResponse](#generateorgdomainvalidationresponse)
  - [GetAppByIDRequest](#getappbyidrequest)
  - [GetAppByIDResponse](#getappbyidresponse)
  - [GetAppKeyRequest](#getappkeyrequest)
  - [GetAppKeyResponse](#getappkeyresponse)
  - [GetDefaultLabelPolicyRequest](#getdefaultlabelpolicyrequest)
  - [GetDefaultLabelPolicyResponse](#getdefaultlabelpolicyresponse)
  - [GetDefaultLoginPolicyRequest](#getdefaultloginpolicyrequest)
  - [GetDefaultLoginPolicyResponse](#getdefaultloginpolicyresponse)
  - [GetDefaultPasswordAgePolicyRequest](#getdefaultpasswordagepolicyrequest)
  - [GetDefaultPasswordAgePolicyResponse](#getdefaultpasswordagepolicyresponse)
  - [GetDefaultPasswordComplexityPolicyRequest](#getdefaultpasswordcomplexitypolicyrequest)
  - [GetDefaultPasswordComplexityPolicyResponse](#getdefaultpasswordcomplexitypolicyresponse)
  - [GetDefaultPasswordLockoutPolicyRequest](#getdefaultpasswordlockoutpolicyrequest)
  - [GetDefaultPasswordLockoutPolicyResponse](#getdefaultpasswordlockoutpolicyresponse)
  - [GetFeaturesRequest](#getfeaturesrequest)
  - [GetFeaturesResponse](#getfeaturesresponse)
  - [GetGrantedProjectByIDRequest](#getgrantedprojectbyidrequest)
  - [GetGrantedProjectByIDResponse](#getgrantedprojectbyidresponse)
  - [GetHumanEmailRequest](#gethumanemailrequest)
  - [GetHumanEmailResponse](#gethumanemailresponse)
  - [GetHumanPhoneRequest](#gethumanphonerequest)
  - [GetHumanPhoneResponse](#gethumanphoneresponse)
  - [GetHumanProfileRequest](#gethumanprofilerequest)
  - [GetHumanProfileResponse](#gethumanprofileresponse)
  - [GetIAMRequest](#getiamrequest)
  - [GetIAMResponse](#getiamresponse)
  - [GetLabelPolicyRequest](#getlabelpolicyrequest)
  - [GetLabelPolicyResponse](#getlabelpolicyresponse)
  - [GetLoginPolicyRequest](#getloginpolicyrequest)
  - [GetLoginPolicyResponse](#getloginpolicyresponse)
  - [GetMachineKeyByIDsRequest](#getmachinekeybyidsrequest)
  - [GetMachineKeyByIDsResponse](#getmachinekeybyidsresponse)
  - [GetMyOrgRequest](#getmyorgrequest)
  - [GetMyOrgResponse](#getmyorgresponse)
  - [GetOIDCInformationRequest](#getoidcinformationrequest)
  - [GetOIDCInformationResponse](#getoidcinformationresponse)
  - [GetOrgByDomainGlobalRequest](#getorgbydomainglobalrequest)
  - [GetOrgByDomainGlobalResponse](#getorgbydomainglobalresponse)
  - [GetOrgIAMPolicyRequest](#getorgiampolicyrequest)
  - [GetOrgIAMPolicyResponse](#getorgiampolicyresponse)
  - [GetOrgIDPByIDRequest](#getorgidpbyidrequest)
  - [GetOrgIDPByIDResponse](#getorgidpbyidresponse)
  - [GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest)
  - [GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)
  - [GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest)
  - [GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)
  - [GetPasswordLockoutPolicyRequest](#getpasswordlockoutpolicyrequest)
  - [GetPasswordLockoutPolicyResponse](#getpasswordlockoutpolicyresponse)
  - [GetProjectByIDRequest](#getprojectbyidrequest)
  - [GetProjectByIDResponse](#getprojectbyidresponse)
  - [GetProjectGrantByIDRequest](#getprojectgrantbyidrequest)
  - [GetProjectGrantByIDResponse](#getprojectgrantbyidresponse)
  - [GetUserByIDRequest](#getuserbyidrequest)
  - [GetUserByIDResponse](#getuserbyidresponse)
  - [GetUserByLoginNameGlobalRequest](#getuserbyloginnameglobalrequest)
  - [GetUserByLoginNameGlobalResponse](#getuserbyloginnameglobalresponse)
  - [GetUserGrantByIDRequest](#getusergrantbyidrequest)
  - [GetUserGrantByIDResponse](#getusergrantbyidresponse)
  - [HealthzRequest](#healthzrequest)
  - [HealthzResponse](#healthzresponse)
  - [IDPQuery](#idpquery)
  - [ImportHumanUserRequest](#importhumanuserrequest)
  - [ImportHumanUserRequest.Email](#importhumanuserrequestemail)
  - [ImportHumanUserRequest.Phone](#importhumanuserrequestphone)
  - [ImportHumanUserRequest.Profile](#importhumanuserrequestprofile)
  - [ImportHumanUserResponse](#importhumanuserresponse)
  - [IsUserUniqueRequest](#isuseruniquerequest)
  - [IsUserUniqueResponse](#isuseruniqueresponse)
  - [ListAppChangesRequest](#listappchangesrequest)
  - [ListAppChangesResponse](#listappchangesresponse)
  - [ListAppKeysRequest](#listappkeysrequest)
  - [ListAppKeysResponse](#listappkeysresponse)
  - [ListAppsRequest](#listappsrequest)
  - [ListAppsResponse](#listappsresponse)
  - [ListGrantedProjectRolesRequest](#listgrantedprojectrolesrequest)
  - [ListGrantedProjectRolesResponse](#listgrantedprojectrolesresponse)
  - [ListGrantedProjectsRequest](#listgrantedprojectsrequest)
  - [ListGrantedProjectsResponse](#listgrantedprojectsresponse)
  - [ListHumanAuthFactorsRequest](#listhumanauthfactorsrequest)
  - [ListHumanAuthFactorsResponse](#listhumanauthfactorsresponse)
  - [ListHumanLinkedIDPsRequest](#listhumanlinkedidpsrequest)
  - [ListHumanLinkedIDPsResponse](#listhumanlinkedidpsresponse)
  - [ListHumanPasswordlessRequest](#listhumanpasswordlessrequest)
  - [ListHumanPasswordlessResponse](#listhumanpasswordlessresponse)
  - [ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest)
  - [ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)
  - [ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest)
  - [ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)
  - [ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest)
  - [ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)
  - [ListMachineKeysRequest](#listmachinekeysrequest)
  - [ListMachineKeysResponse](#listmachinekeysresponse)
  - [ListOrgChangesRequest](#listorgchangesrequest)
  - [ListOrgChangesResponse](#listorgchangesresponse)
  - [ListOrgDomainsRequest](#listorgdomainsrequest)
  - [ListOrgDomainsResponse](#listorgdomainsresponse)
  - [ListOrgIDPsRequest](#listorgidpsrequest)
  - [ListOrgIDPsResponse](#listorgidpsresponse)
  - [ListOrgMemberRolesRequest](#listorgmemberrolesrequest)
  - [ListOrgMemberRolesResponse](#listorgmemberrolesresponse)
  - [ListOrgMembersRequest](#listorgmembersrequest)
  - [ListOrgMembersResponse](#listorgmembersresponse)
  - [ListProjectChangesRequest](#listprojectchangesrequest)
  - [ListProjectChangesResponse](#listprojectchangesresponse)
  - [ListProjectGrantMemberRolesRequest](#listprojectgrantmemberrolesrequest)
  - [ListProjectGrantMemberRolesResponse](#listprojectgrantmemberrolesresponse)
  - [ListProjectGrantMembersRequest](#listprojectgrantmembersrequest)
  - [ListProjectGrantMembersResponse](#listprojectgrantmembersresponse)
  - [ListProjectGrantsRequest](#listprojectgrantsrequest)
  - [ListProjectGrantsResponse](#listprojectgrantsresponse)
  - [ListProjectMemberRolesRequest](#listprojectmemberrolesrequest)
  - [ListProjectMemberRolesResponse](#listprojectmemberrolesresponse)
  - [ListProjectMembersRequest](#listprojectmembersrequest)
  - [ListProjectMembersResponse](#listprojectmembersresponse)
  - [ListProjectRolesRequest](#listprojectrolesrequest)
  - [ListProjectRolesResponse](#listprojectrolesresponse)
  - [ListProjectsRequest](#listprojectsrequest)
  - [ListProjectsResponse](#listprojectsresponse)
  - [ListUserChangesRequest](#listuserchangesrequest)
  - [ListUserChangesResponse](#listuserchangesresponse)
  - [ListUserGrantRequest](#listusergrantrequest)
  - [ListUserGrantResponse](#listusergrantresponse)
  - [ListUserMembershipsRequest](#listusermembershipsrequest)
  - [ListUserMembershipsResponse](#listusermembershipsresponse)
  - [ListUsersRequest](#listusersrequest)
  - [ListUsersResponse](#listusersresponse)
  - [LockUserRequest](#lockuserrequest)
  - [LockUserResponse](#lockuserresponse)
  - [ReactivateAppRequest](#reactivateapprequest)
  - [ReactivateAppResponse](#reactivateappresponse)
  - [ReactivateOrgIDPRequest](#reactivateorgidprequest)
  - [ReactivateOrgIDPResponse](#reactivateorgidpresponse)
  - [ReactivateOrgRequest](#reactivateorgrequest)
  - [ReactivateOrgResponse](#reactivateorgresponse)
  - [ReactivateProjectGrantRequest](#reactivateprojectgrantrequest)
  - [ReactivateProjectGrantResponse](#reactivateprojectgrantresponse)
  - [ReactivateProjectRequest](#reactivateprojectrequest)
  - [ReactivateProjectResponse](#reactivateprojectresponse)
  - [ReactivateUserGrantRequest](#reactivateusergrantrequest)
  - [ReactivateUserGrantResponse](#reactivateusergrantresponse)
  - [ReactivateUserRequest](#reactivateuserrequest)
  - [ReactivateUserResponse](#reactivateuserresponse)
  - [RegenerateAPIClientSecretRequest](#regenerateapiclientsecretrequest)
  - [RegenerateAPIClientSecretResponse](#regenerateapiclientsecretresponse)
  - [RegenerateOIDCClientSecretRequest](#regenerateoidcclientsecretrequest)
  - [RegenerateOIDCClientSecretResponse](#regenerateoidcclientsecretresponse)
  - [RemoveAppKeyRequest](#removeappkeyrequest)
  - [RemoveAppKeyResponse](#removeappkeyresponse)
  - [RemoveAppRequest](#removeapprequest)
  - [RemoveAppResponse](#removeappresponse)
  - [RemoveHumanAuthFactorOTPRequest](#removehumanauthfactorotprequest)
  - [RemoveHumanAuthFactorOTPResponse](#removehumanauthfactorotpresponse)
  - [RemoveHumanAuthFactorU2FRequest](#removehumanauthfactoru2frequest)
  - [RemoveHumanAuthFactorU2FResponse](#removehumanauthfactoru2fresponse)
  - [RemoveHumanLinkedIDPRequest](#removehumanlinkedidprequest)
  - [RemoveHumanLinkedIDPResponse](#removehumanlinkedidpresponse)
  - [RemoveHumanPasswordlessRequest](#removehumanpasswordlessrequest)
  - [RemoveHumanPasswordlessResponse](#removehumanpasswordlessresponse)
  - [RemoveHumanPhoneRequest](#removehumanphonerequest)
  - [RemoveHumanPhoneResponse](#removehumanphoneresponse)
  - [RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest)
  - [RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)
  - [RemoveMachineKeyRequest](#removemachinekeyrequest)
  - [RemoveMachineKeyResponse](#removemachinekeyresponse)
  - [RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest)
  - [RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)
  - [RemoveOrgDomainRequest](#removeorgdomainrequest)
  - [RemoveOrgDomainResponse](#removeorgdomainresponse)
  - [RemoveOrgIDPRequest](#removeorgidprequest)
  - [RemoveOrgIDPResponse](#removeorgidpresponse)
  - [RemoveOrgMemberRequest](#removeorgmemberrequest)
  - [RemoveOrgMemberResponse](#removeorgmemberresponse)
  - [RemoveProjectGrantMemberRequest](#removeprojectgrantmemberrequest)
  - [RemoveProjectGrantMemberResponse](#removeprojectgrantmemberresponse)
  - [RemoveProjectGrantRequest](#removeprojectgrantrequest)
  - [RemoveProjectGrantResponse](#removeprojectgrantresponse)
  - [RemoveProjectMemberRequest](#removeprojectmemberrequest)
  - [RemoveProjectMemberResponse](#removeprojectmemberresponse)
  - [RemoveProjectRequest](#removeprojectrequest)
  - [RemoveProjectResponse](#removeprojectresponse)
  - [RemoveProjectRoleRequest](#removeprojectrolerequest)
  - [RemoveProjectRoleResponse](#removeprojectroleresponse)
  - [RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest)
  - [RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)
  - [RemoveUserGrantRequest](#removeusergrantrequest)
  - [RemoveUserGrantResponse](#removeusergrantresponse)
  - [RemoveUserRequest](#removeuserrequest)
  - [RemoveUserResponse](#removeuserresponse)
  - [ResendHumanEmailVerificationRequest](#resendhumanemailverificationrequest)
  - [ResendHumanEmailVerificationResponse](#resendhumanemailverificationresponse)
  - [ResendHumanInitializationRequest](#resendhumaninitializationrequest)
  - [ResendHumanInitializationResponse](#resendhumaninitializationresponse)
  - [ResendHumanPhoneVerificationRequest](#resendhumanphoneverificationrequest)
  - [ResendHumanPhoneVerificationResponse](#resendhumanphoneverificationresponse)
  - [ResetLabelPolicyToDefaultRequest](#resetlabelpolicytodefaultrequest)
  - [ResetLabelPolicyToDefaultResponse](#resetlabelpolicytodefaultresponse)
  - [ResetLoginPolicyToDefaultRequest](#resetloginpolicytodefaultrequest)
  - [ResetLoginPolicyToDefaultResponse](#resetloginpolicytodefaultresponse)
  - [ResetPasswordAgePolicyToDefaultRequest](#resetpasswordagepolicytodefaultrequest)
  - [ResetPasswordAgePolicyToDefaultResponse](#resetpasswordagepolicytodefaultresponse)
  - [ResetPasswordComplexityPolicyToDefaultRequest](#resetpasswordcomplexitypolicytodefaultrequest)
  - [ResetPasswordComplexityPolicyToDefaultResponse](#resetpasswordcomplexitypolicytodefaultresponse)
  - [ResetPasswordLockoutPolicyToDefaultRequest](#resetpasswordlockoutpolicytodefaultrequest)
  - [ResetPasswordLockoutPolicyToDefaultResponse](#resetpasswordlockoutpolicytodefaultresponse)
  - [SendHumanResetPasswordNotificationRequest](#sendhumanresetpasswordnotificationrequest)
  - [SendHumanResetPasswordNotificationResponse](#sendhumanresetpasswordnotificationresponse)
  - [SetHumanInitialPasswordRequest](#sethumaninitialpasswordrequest)
  - [SetHumanInitialPasswordResponse](#sethumaninitialpasswordresponse)
  - [SetPrimaryOrgDomainRequest](#setprimaryorgdomainrequest)
  - [SetPrimaryOrgDomainResponse](#setprimaryorgdomainresponse)
  - [UnlockUserRequest](#unlockuserrequest)
  - [UnlockUserResponse](#unlockuserresponse)
  - [UpdateAPIAppConfigRequest](#updateapiappconfigrequest)
  - [UpdateAPIAppConfigResponse](#updateapiappconfigresponse)
  - [UpdateAppRequest](#updateapprequest)
  - [UpdateAppResponse](#updateappresponse)
  - [UpdateCustomLabelPolicyRequest](#updatecustomlabelpolicyrequest)
  - [UpdateCustomLabelPolicyResponse](#updatecustomlabelpolicyresponse)
  - [UpdateCustomLoginPolicyRequest](#updatecustomloginpolicyrequest)
  - [UpdateCustomLoginPolicyResponse](#updatecustomloginpolicyresponse)
  - [UpdateCustomPasswordAgePolicyRequest](#updatecustompasswordagepolicyrequest)
  - [UpdateCustomPasswordAgePolicyResponse](#updatecustompasswordagepolicyresponse)
  - [UpdateCustomPasswordComplexityPolicyRequest](#updatecustompasswordcomplexitypolicyrequest)
  - [UpdateCustomPasswordComplexityPolicyResponse](#updatecustompasswordcomplexitypolicyresponse)
  - [UpdateCustomPasswordLockoutPolicyRequest](#updatecustompasswordlockoutpolicyrequest)
  - [UpdateCustomPasswordLockoutPolicyResponse](#updatecustompasswordlockoutpolicyresponse)
  - [UpdateHumanEmailRequest](#updatehumanemailrequest)
  - [UpdateHumanEmailResponse](#updatehumanemailresponse)
  - [UpdateHumanPhoneRequest](#updatehumanphonerequest)
  - [UpdateHumanPhoneResponse](#updatehumanphoneresponse)
  - [UpdateHumanProfileRequest](#updatehumanprofilerequest)
  - [UpdateHumanProfileResponse](#updatehumanprofileresponse)
  - [UpdateMachineRequest](#updatemachinerequest)
  - [UpdateMachineResponse](#updatemachineresponse)
  - [UpdateOIDCAppConfigRequest](#updateoidcappconfigrequest)
  - [UpdateOIDCAppConfigResponse](#updateoidcappconfigresponse)
  - [UpdateOrgIDPOIDCConfigRequest](#updateorgidpoidcconfigrequest)
  - [UpdateOrgIDPOIDCConfigResponse](#updateorgidpoidcconfigresponse)
  - [UpdateOrgIDPRequest](#updateorgidprequest)
  - [UpdateOrgIDPResponse](#updateorgidpresponse)
  - [UpdateOrgMemberRequest](#updateorgmemberrequest)
  - [UpdateOrgMemberResponse](#updateorgmemberresponse)
  - [UpdateProjectGrantMemberRequest](#updateprojectgrantmemberrequest)
  - [UpdateProjectGrantMemberResponse](#updateprojectgrantmemberresponse)
  - [UpdateProjectGrantRequest](#updateprojectgrantrequest)
  - [UpdateProjectGrantResponse](#updateprojectgrantresponse)
  - [UpdateProjectMemberRequest](#updateprojectmemberrequest)
  - [UpdateProjectMemberResponse](#updateprojectmemberresponse)
  - [UpdateProjectRequest](#updateprojectrequest)
  - [UpdateProjectResponse](#updateprojectresponse)
  - [UpdateProjectRoleRequest](#updateprojectrolerequest)
  - [UpdateProjectRoleResponse](#updateprojectroleresponse)
  - [UpdateUserGrantRequest](#updateusergrantrequest)
  - [UpdateUserGrantResponse](#updateusergrantresponse)
  - [UpdateUserNameRequest](#updateusernamerequest)
  - [UpdateUserNameResponse](#updateusernameresponse)
  - [ValidateOrgDomainRequest](#validateorgdomainrequest)
  - [ValidateOrgDomainResponse](#validateorgdomainresponse)






- Messages
  - [EmailQuery](#emailquery)
  - [FirstNameQuery](#firstnamequery)
  - [LastNameQuery](#lastnamequery)
  - [Member](#member)
  - [SearchQuery](#searchquery)
  - [UserIDQuery](#useridquery)






- Messages
  - [ErrorDetail](#errordetail)
  - [LocalizedMessage](#localizedmessage)






- Messages
  - [ListDetails](#listdetails)
  - [ListQuery](#listquery)
  - [ObjectDetails](#objectdetails)



- Enums
  - [TextQueryMethod](#textquerymethod)





- Messages
  - [AuthOption](#authoption)






- Messages
  - [Domain](#domain)
  - [DomainNameQuery](#domainnamequery)
  - [DomainSearchQuery](#domainsearchquery)
  - [Org](#org)
  - [OrgDomainQuery](#orgdomainquery)
  - [OrgNameQuery](#orgnamequery)
  - [OrgQuery](#orgquery)



- Enums
  - [DomainValidationType](#domainvalidationtype)
  - [OrgFieldName](#orgfieldname)
  - [OrgState](#orgstate)





- Messages
  - [LabelPolicy](#labelpolicy)
  - [LoginPolicy](#loginpolicy)
  - [OrgIAMPolicy](#orgiampolicy)
  - [PasswordAgePolicy](#passwordagepolicy)
  - [PasswordComplexityPolicy](#passwordcomplexitypolicy)
  - [PasswordLockoutPolicy](#passwordlockoutpolicy)



- Enums
  - [MultiFactorType](#multifactortype)
  - [PasswordlessType](#passwordlesstype)
  - [SecondFactorType](#secondfactortype)





- Messages
  - [GrantProjectNameQuery](#grantprojectnamequery)
  - [GrantRoleKeyQuery](#grantrolekeyquery)
  - [GrantedProject](#grantedproject)
  - [Project](#project)
  - [ProjectGrantQuery](#projectgrantquery)
  - [ProjectNameQuery](#projectnamequery)
  - [ProjectQuery](#projectquery)
  - [Role](#role)
  - [RoleDisplayNameQuery](#roledisplaynamequery)
  - [RoleKeyQuery](#rolekeyquery)
  - [RoleQuery](#rolequery)



- Enums
  - [ProjectGrantState](#projectgrantstate)
  - [ProjectState](#projectstate)





- Messages
  - [AuthFactor](#authfactor)
  - [AuthFactorOTP](#authfactorotp)
  - [AuthFactorU2F](#authfactoru2f)
  - [DisplayNameQuery](#displaynamequery)
  - [Email](#email)
  - [EmailQuery](#emailquery)
  - [FirstNameQuery](#firstnamequery)
  - [Human](#human)
  - [LastNameQuery](#lastnamequery)
  - [Machine](#machine)
  - [Membership](#membership)
  - [MembershipIAMQuery](#membershipiamquery)
  - [MembershipOrgQuery](#membershiporgquery)
  - [MembershipProjectGrantQuery](#membershipprojectgrantquery)
  - [MembershipProjectQuery](#membershipprojectquery)
  - [MembershipQuery](#membershipquery)
  - [NickNameQuery](#nicknamequery)
  - [Phone](#phone)
  - [Profile](#profile)
  - [SearchQuery](#searchquery)
  - [Session](#session)
  - [StateQuery](#statequery)
  - [TypeQuery](#typequery)
  - [User](#user)
  - [UserGrant](#usergrant)
  - [UserGrantDisplayNameQuery](#usergrantdisplaynamequery)
  - [UserGrantEmailQuery](#usergrantemailquery)
  - [UserGrantFirstNameQuery](#usergrantfirstnamequery)
  - [UserGrantLastNameQuery](#usergrantlastnamequery)
  - [UserGrantOrgDomainQuery](#usergrantorgdomainquery)
  - [UserGrantOrgNameQuery](#usergrantorgnamequery)
  - [UserGrantProjectGrantIDQuery](#usergrantprojectgrantidquery)
  - [UserGrantProjectIDQuery](#usergrantprojectidquery)
  - [UserGrantProjectNameQuery](#usergrantprojectnamequery)
  - [UserGrantQuery](#usergrantquery)
  - [UserGrantRoleKeyQuery](#usergrantrolekeyquery)
  - [UserGrantUserIDQuery](#usergrantuseridquery)
  - [UserGrantUserNameQuery](#usergrantusernamequery)
  - [UserGrantWithGrantedQuery](#usergrantwithgrantedquery)
  - [UserNameQuery](#usernamequery)
  - [WebAuthNKey](#webauthnkey)
  - [WebAuthNToken](#webauthntoken)
  - [WebAuthNVerification](#webauthnverification)



- Enums
  - [AuthFactorState](#authfactorstate)
  - [Gender](#gender)
  - [SessionState](#sessionstate)
  - [Type](#type)
  - [UserFieldName](#userfieldname)
  - [UserGrantState](#usergrantstate)
  - [UserState](#userstate)



- [Scalar Value Types](#scalar-value-types)



### AdminService {#zitadeladminv1adminservice}


#### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)

Indicates if ZITADEL is running.
It respondes as soon as ZITADEL started


    Key:google.api.http
    Value:{[{GET /healthz }]}
 <!-- end options -->

#### IsOrgUnique

> **rpc** IsOrgUnique([IsOrgUniqueRequest](#isorguniquerequest))
[IsOrgUniqueResponse](#isorguniqueresponse)

Checks whether an organisation exists by the given parameters


    Key:google.api.http
    Value:{[{GET /orgs/_is_unique }]}
 <!-- end options -->

#### GetOrgByID

> **rpc** GetOrgByID([GetOrgByIDRequest](#getorgbyidrequest))
[GetOrgByIDResponse](#getorgbyidresponse)




    Key:google.api.http
    Value:{[{GET /orgs/{id} }]}
 <!-- end options -->

#### ListOrgs

> **rpc** ListOrgs([ListOrgsRequest](#listorgsrequest))
[ListOrgsResponse](#listorgsresponse)

Returns all organisations matching the request
all queries need to match (ANDed)


    Key:google.api.http
    Value:{[{POST /orgs/_search *}]}
 <!-- end options -->

#### SetUpOrg

> **rpc** SetUpOrg([SetUpOrgRequest](#setuporgrequest))
[SetUpOrgResponse](#setuporgresponse)

Creates a new org and user 
and adds the user to the orgs members as ORG_OWNER


    Key:google.api.http
    Value:{[{POST /orgs/_setup *}]}
 <!-- end options -->

#### GetIDPByID

> **rpc** GetIDPByID([GetIDPByIDRequest](#getidpbyidrequest))
[GetIDPByIDResponse](#getidpbyidresponse)




    Key:google.api.http
    Value:{[{GET /idps/{id} }]}
 <!-- end options -->

#### ListIDPs

> **rpc** ListIDPs([ListIDPsRequest](#listidpsrequest))
[ListIDPsResponse](#listidpsresponse)




    Key:google.api.http
    Value:{[{POST /idps/_search *}]}
 <!-- end options -->

#### AddOIDCIDP

> **rpc** AddOIDCIDP([AddOIDCIDPRequest](#addoidcidprequest))
[AddOIDCIDPResponse](#addoidcidpresponse)




    Key:google.api.http
    Value:{[{POST /idps/oidc *}]}
 <!-- end options -->

#### UpdateIDP

> **rpc** UpdateIDP([UpdateIDPRequest](#updateidprequest))
[UpdateIDPResponse](#updateidpresponse)

Updates the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.


    Key:google.api.http
    Value:{[{PUT /idps/{idp_id} *}]}
 <!-- end options -->

#### DeactivateIDP

> **rpc** DeactivateIDP([DeactivateIDPRequest](#deactivateidprequest))
[DeactivateIDPResponse](#deactivateidpresponse)

Sets the state of the idp to IDP_STATE_INACTIVE
the state MUST be IDP_STATE_ACTIVE for this call


    Key:google.api.http
    Value:{[{POST /idps/{idp_id}/_deactivate }]}
 <!-- end options -->

#### ReactivateIDP

> **rpc** ReactivateIDP([ReactivateIDPRequest](#reactivateidprequest))
[ReactivateIDPResponse](#reactivateidpresponse)

Sets the state of the idp to IDP_STATE_ACTIVE
the state MUST be IDP_STATE_INACTIVE for this call


    Key:google.api.http
    Value:{[{POST /idps/{idp_id}/_reactivate }]}
 <!-- end options -->

#### RemoveIDP

> **rpc** RemoveIDP([RemoveIDPRequest](#removeidprequest))
[RemoveIDPResponse](#removeidpresponse)

RemoveIDP deletes the IDP permanetly


    Key:google.api.http
    Value:{[{DELETE /idps/{idp_id} }]}
 <!-- end options -->

#### UpdateIDPOIDCConfig

> **rpc** UpdateIDPOIDCConfig([UpdateIDPOIDCConfigRequest](#updateidpoidcconfigrequest))
[UpdateIDPOIDCConfigResponse](#updateidpoidcconfigresponse)

Updates the oidc configuration of the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.


    Key:google.api.http
    Value:{[{PUT /idps/{idp_id}/oidc_config *}]}
 <!-- end options -->

#### GetDefaultFeatures

> **rpc** GetDefaultFeatures([GetDefaultFeaturesRequest](#getdefaultfeaturesrequest))
[GetDefaultFeaturesResponse](#getdefaultfeaturesresponse)




    Key:google.api.http
    Value:{[{GET /features }]}
 <!-- end options -->

#### SetDefaultFeatures

> **rpc** SetDefaultFeatures([SetDefaultFeaturesRequest](#setdefaultfeaturesrequest))
[SetDefaultFeaturesResponse](#setdefaultfeaturesresponse)




    Key:google.api.http
    Value:{[{PUT /features *}]}
 <!-- end options -->

#### GetOrgFeatures

> **rpc** GetOrgFeatures([GetOrgFeaturesRequest](#getorgfeaturesrequest))
[GetOrgFeaturesResponse](#getorgfeaturesresponse)




    Key:google.api.http
    Value:{[{GET /orgs/{org_id}/features }]}
 <!-- end options -->

#### SetOrgFeatures

> **rpc** SetOrgFeatures([SetOrgFeaturesRequest](#setorgfeaturesrequest))
[SetOrgFeaturesResponse](#setorgfeaturesresponse)




    Key:google.api.http
    Value:{[{PUT /orgs/{org_id}/features *}]}
 <!-- end options -->

#### ResetOrgFeatures

> **rpc** ResetOrgFeatures([ResetOrgFeaturesRequest](#resetorgfeaturesrequest))
[ResetOrgFeaturesResponse](#resetorgfeaturesresponse)




    Key:google.api.http
    Value:{[{DELETE /orgs/{org_id}/features }]}
 <!-- end options -->

#### GetOrgIAMPolicy

> **rpc** GetOrgIAMPolicy([GetOrgIAMPolicyRequest](#getorgiampolicyrequest))
[GetOrgIAMPolicyResponse](#getorgiampolicyresponse)

Returns the IAM policy defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{GET /policies/orgiam }]}
 <!-- end options -->

#### UpdateOrgIAMPolicy

> **rpc** UpdateOrgIAMPolicy([UpdateOrgIAMPolicyRequest](#updateorgiampolicyrequest))
[UpdateOrgIAMPolicyResponse](#updateorgiampolicyresponse)

Updates the default IAM policy.
it impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{PUT /policies/orgiam *}]}
 <!-- end options -->

#### GetCustomOrgIAMPolicy

> **rpc** GetCustomOrgIAMPolicy([GetCustomOrgIAMPolicyRequest](#getcustomorgiampolicyrequest))
[GetCustomOrgIAMPolicyResponse](#getcustomorgiampolicyresponse)

Returns the customised policy or the default if not customised


    Key:google.api.http
    Value:{[{GET /orgs/{org_id}/policies/orgiam }]}
 <!-- end options -->

#### AddCustomOrgIAMPolicy

> **rpc** AddCustomOrgIAMPolicy([AddCustomOrgIAMPolicyRequest](#addcustomorgiampolicyrequest))
[AddCustomOrgIAMPolicyResponse](#addcustomorgiampolicyresponse)

Defines a custom ORGIAM policy as specified


    Key:google.api.http
    Value:{[{POST /orgs/{org_id}/policies/orgiam *}]}
 <!-- end options -->

#### UpdateCustomOrgIAMPolicy

> **rpc** UpdateCustomOrgIAMPolicy([UpdateCustomOrgIAMPolicyRequest](#updatecustomorgiampolicyrequest))
[UpdateCustomOrgIAMPolicyResponse](#updatecustomorgiampolicyresponse)

Updates a custom ORGIAM policy as specified


    Key:google.api.http
    Value:{[{PUT /orgs/{org_id}/policies/orgiam *}]}
 <!-- end options -->

#### ResetCustomOrgIAMPolicyToDefault

> **rpc** ResetCustomOrgIAMPolicyToDefault([ResetCustomOrgIAMPolicyToDefaultRequest](#resetcustomorgiampolicytodefaultrequest))
[ResetCustomOrgIAMPolicyToDefaultResponse](#resetcustomorgiampolicytodefaultresponse)

Resets the org iam policy of the organisation to default
ZITADEL will fallback to the default policy defined by the ZITADEL administrators


    Key:google.api.http
    Value:{[{DELETE /orgs/{org_id}/policies/orgiam }]}
 <!-- end options -->

#### GetLabelPolicy

> **rpc** GetLabelPolicy([GetLabelPolicyRequest](#getlabelpolicyrequest))
[GetLabelPolicyResponse](#getlabelpolicyresponse)

Returns the label policy defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{GET /policies/label }]}
 <!-- end options -->

#### UpdateLabelPolicy

> **rpc** UpdateLabelPolicy([UpdateLabelPolicyRequest](#updatelabelpolicyrequest))
[UpdateLabelPolicyResponse](#updatelabelpolicyresponse)

Updates the default label policy of ZITADEL
it impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{PUT /policies/label *}]}
 <!-- end options -->

#### GetLoginPolicy

> **rpc** GetLoginPolicy([GetLoginPolicyRequest](#getloginpolicyrequest))
[GetLoginPolicyResponse](#getloginpolicyresponse)

Returns the login policy defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{GET /policies/login }]}
 <!-- end options -->

#### UpdateLoginPolicy

> **rpc** UpdateLoginPolicy([UpdateLoginPolicyRequest](#updateloginpolicyrequest))
[UpdateLoginPolicyResponse](#updateloginpolicyresponse)

Updates the default login policy of ZITADEL
it impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{PUT /policies/login *}]}
 <!-- end options -->

#### ListLoginPolicyIDPs

> **rpc** ListLoginPolicyIDPs([ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest))
[ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)

Returns the idps linked to the default login policy,
defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{POST /policies/login/idps/_search *}]}
 <!-- end options -->

#### AddIDPToLoginPolicy

> **rpc** AddIDPToLoginPolicy([AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest))
[AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)

Adds the povided idp to the default login policy.
It impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{POST /policies/login/idps *}]}
 <!-- end options -->

#### RemoveIDPFromLoginPolicy

> **rpc** RemoveIDPFromLoginPolicy([RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest))
[RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)

Removes the povided idp from the default login policy.
It impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{DELETE /policies/login/idps/{idp_id} }]}
 <!-- end options -->

#### ListLoginPolicySecondFactors

> **rpc** ListLoginPolicySecondFactors([ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest))
[ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)

Returns the available second factors defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{POST /policies/login/second_factors/_search }]}
 <!-- end options -->

#### AddSecondFactorToLoginPolicy

> **rpc** AddSecondFactorToLoginPolicy([AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest))
[AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)

Adds a second factor to the default login policy.
It impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{POST /policies/login/second_factors *}]}
 <!-- end options -->

#### RemoveSecondFactorFromLoginPolicy

> **rpc** RemoveSecondFactorFromLoginPolicy([RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest))
[RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)

Removes a second factor from the default login policy.
It impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{DELETE /policies/login/second_factors/{type} }]}
 <!-- end options -->

#### ListLoginPolicyMultiFactors

> **rpc** ListLoginPolicyMultiFactors([ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest))
[ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)

Returns the available multi factors defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{POST /policies/login/multi_factors/_search }]}
 <!-- end options -->

#### AddMultiFactorToLoginPolicy

> **rpc** AddMultiFactorToLoginPolicy([AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest))
[AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)

Adds a multi factor to the default login policy.
It impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{POST /policies/login/multi_factors *}]}
 <!-- end options -->

#### RemoveMultiFactorFromLoginPolicy

> **rpc** RemoveMultiFactorFromLoginPolicy([RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest))
[RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)

Removes a multi factor from the default login policy.
It impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{DELETE /policies/login/multi_factors/{type} }]}
 <!-- end options -->

#### GetPasswordComplexityPolicy

> **rpc** GetPasswordComplexityPolicy([GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest))
[GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)

Returns the password complexity policy defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{GET /policies/password/complexity }]}
 <!-- end options -->

#### UpdatePasswordComplexityPolicy

> **rpc** UpdatePasswordComplexityPolicy([UpdatePasswordComplexityPolicyRequest](#updatepasswordcomplexitypolicyrequest))
[UpdatePasswordComplexityPolicyResponse](#updatepasswordcomplexitypolicyresponse)

Updates the default password complexity policy of ZITADEL
it impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{PUT /policies/password/complexity *}]}
 <!-- end options -->

#### GetPasswordAgePolicy

> **rpc** GetPasswordAgePolicy([GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest))
[GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)

Returns the password age policy defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{GET /policies/password/age }]}
 <!-- end options -->

#### UpdatePasswordAgePolicy

> **rpc** UpdatePasswordAgePolicy([UpdatePasswordAgePolicyRequest](#updatepasswordagepolicyrequest))
[UpdatePasswordAgePolicyResponse](#updatepasswordagepolicyresponse)

Updates the default password age policy of ZITADEL
it impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{PUT /policies/password/age *}]}
 <!-- end options -->

#### GetPasswordLockoutPolicy

> **rpc** GetPasswordLockoutPolicy([GetPasswordLockoutPolicyRequest](#getpasswordlockoutpolicyrequest))
[GetPasswordLockoutPolicyResponse](#getpasswordlockoutpolicyresponse)

Returns the password lockout policy defined by the administrators of ZITADEL


    Key:google.api.http
    Value:{[{GET /policies/password/lockout }]}
 <!-- end options -->

#### UpdatePasswordLockoutPolicy

> **rpc** UpdatePasswordLockoutPolicy([UpdatePasswordLockoutPolicyRequest](#updatepasswordlockoutpolicyrequest))
[UpdatePasswordLockoutPolicyResponse](#updatepasswordlockoutpolicyresponse)

Updates the default password lockout policy of ZITADEL
it impacts all organisations without a customised policy


    Key:google.api.http
    Value:{[{PUT /policies/password/lockout *}]}
 <!-- end options -->

#### ListIAMMemberRoles

> **rpc** ListIAMMemberRoles([ListIAMMemberRolesRequest](#listiammemberrolesrequest))
[ListIAMMemberRolesResponse](#listiammemberrolesresponse)

Returns the IAM roles visible for the requested user


    Key:google.api.http
    Value:{[{POST /members/roles/_search }]}
 <!-- end options -->

#### ListIAMMembers

> **rpc** ListIAMMembers([ListIAMMembersRequest](#listiammembersrequest))
[ListIAMMembersResponse](#listiammembersresponse)

Returns all members matching the request
all queries need to match (ANDed)


    Key:google.api.http
    Value:{[{POST /members/_search *}]}
 <!-- end options -->

#### AddIAMMember

> **rpc** AddIAMMember([AddIAMMemberRequest](#addiammemberrequest))
[AddIAMMemberResponse](#addiammemberresponse)

Adds a user to the membership list of ZITADEL with the given roles
undefined roles will be dropped


    Key:google.api.http
    Value:{[{POST /members *}]}
 <!-- end options -->

#### UpdateIAMMember

> **rpc** UpdateIAMMember([UpdateIAMMemberRequest](#updateiammemberrequest))
[UpdateIAMMemberResponse](#updateiammemberresponse)

Sets the given roles on a member.
The member has only roles provided by this call


    Key:google.api.http
    Value:{[{PUT /members/{user_id} *}]}
 <!-- end options -->

#### RemoveIAMMember

> **rpc** RemoveIAMMember([RemoveIAMMemberRequest](#removeiammemberrequest))
[RemoveIAMMemberResponse](#removeiammemberresponse)

Removes the user from the membership list of ZITADEL


    Key:google.api.http
    Value:{[{DELETE /members/{user_id} }]}
 <!-- end options -->

#### ListViews

> **rpc** ListViews([ListViewsRequest](#listviewsrequest))
[ListViewsResponse](#listviewsresponse)

Returns all stored read models of ZITADEL
views are used for search optimisation and optimise request latencies
they represent the delta of the event happend on the objects


    Key:google.api.http
    Value:{[{POST /views/_search }]}
 <!-- end options -->

#### ClearView

> **rpc** ClearView([ClearViewRequest](#clearviewrequest))
[ClearViewResponse](#clearviewresponse)

Truncates the delta of the change stream
be carefull with this function because ZITADEL has to 
recompute the deltas after they got cleared. 
Search requests will return wrong results until all deltas are recomputed


    Key:google.api.http
    Value:{[{POST /views/{database}/{view_name} }]}
 <!-- end options -->

#### ListFailedEvents

> **rpc** ListFailedEvents([ListFailedEventsRequest](#listfailedeventsrequest))
[ListFailedEventsResponse](#listfailedeventsresponse)

Returns event descriptions which cannot be processed.
It's possible that some events need some retries. 
For example if the SMTP-API wasn't able to send an email at the first time


    Key:google.api.http
    Value:{[{POST /failedevents/_search }]}
 <!-- end options -->

#### RemoveFailedEvent

> **rpc** RemoveFailedEvent([RemoveFailedEventRequest](#removefailedeventrequest))
[RemoveFailedEventResponse](#removefailedeventresponse)

Deletes the event from failed events view.
the event is not removed from the change stream
This call is usefull if the system was able to process the event later. 
e.g. if the second try of sending an email was successful. the first try produced a
failed event. You can find out if it worked on the `failure_count`


    Key:google.api.http
    Value:{[{DELETE /failedevents/{database}/{view_name}/{failed_sequence} }]}
 <!-- end options -->

 <!-- end methods -->
 <!-- end services -->

### Messages


#### AddCustomOrgIAMPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
| user_login_must_be_domain | [ bool](#bool) | the username has to end with the domain of it's organisation (uniqueness is organisation based) |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomOrgIAMPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddIAMMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | if no roles provided the user won't have any rights |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddIAMMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddIDPToLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | Id of the predefined idp configuration |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddIDPToLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMultiFactorToLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ zitadel.policy.v1.MultiFactorType](#zitadelpolicyv1multifactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMultiFactorToLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOIDCIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| styling_type | [ zitadel.idp.v1.IDPStylingType](#zitadelidpv1idpstylingtype) | - |
| client_id | [ string](#string) | - |
| client_secret | [ string](#string) | - |
| issuer | [ string](#string) | - |
| scopes | [repeated string](#string) | - |
| display_name_mapping | [ zitadel.idp.v1.OIDCMappingField](#zitadelidpv1oidcmappingfield) | - |
| username_mapping | [ zitadel.idp.v1.OIDCMappingField](#zitadelidpv1oidcmappingfield) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOIDCIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddSecondFactorToLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ zitadel.policy.v1.SecondFactorType](#zitadelpolicyv1secondfactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddSecondFactorToLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ClearViewRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| database | [ string](#string) | - |
| view_name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ClearViewResponse


 <!-- end HasFields -->


#### DeactivateIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### FailedEvent



| Field | Type | Description |
| ----- | ---- | ----------- |
| database | [ string](#string) | - |
| view_name | [ string](#string) | - |
| failed_sequence | [ uint64](#uint64) | - |
| failure_count | [ uint64](#uint64) | - |
| error_message | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetCustomOrgIAMPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetCustomOrgIAMPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.OrgIAMPolicy](#zitadelpolicyv1orgiampolicy) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetDefaultFeaturesRequest


 <!-- end HasFields -->


#### GetDefaultFeaturesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| features | [ zitadel.features.v1.Features](#zitadelfeaturesv1features) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetIDPByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetIDPByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp | [ zitadel.idp.v1.IDP](#zitadelidpv1idp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetLabelPolicyRequest


 <!-- end HasFields -->


#### GetLabelPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.LabelPolicy](#zitadelpolicyv1labelpolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetLoginPolicyRequest


 <!-- end HasFields -->


#### GetLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.LoginPolicy](#zitadelpolicyv1loginpolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| org | [ zitadel.org.v1.Org](#zitadelorgv1org) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgFeaturesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgFeaturesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| features | [ zitadel.features.v1.Features](#zitadelfeaturesv1features) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgIAMPolicyRequest


 <!-- end HasFields -->


#### GetOrgIAMPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.OrgIAMPolicy](#zitadelpolicyv1orgiampolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetPasswordAgePolicyRequest


 <!-- end HasFields -->


#### GetPasswordAgePolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordAgePolicy](#zitadelpolicyv1passwordagepolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetPasswordComplexityPolicyRequest


 <!-- end HasFields -->


#### GetPasswordComplexityPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordComplexityPolicy](#zitadelpolicyv1passwordcomplexitypolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetPasswordLockoutPolicyRequest


 <!-- end HasFields -->


#### GetPasswordLockoutPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordLockoutPolicy](#zitadelpolicyv1passwordlockoutpolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### HealthzRequest


 <!-- end HasFields -->


#### HealthzResponse


 <!-- end HasFields -->


#### IDPQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_id_query | [ zitadel.idp.v1.IDPIDQuery](#zitadelidpv1idpidquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_name_query | [ zitadel.idp.v1.IDPNameQuery](#zitadelidpv1idpnamequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IsOrgUniqueRequest
parameters are ORed


| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| domain | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IsOrgUniqueResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| is_unique | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListFailedEventsRequest


 <!-- end HasFields -->


#### ListFailedEventsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated FailedEvent](#failedevent) | TODO: list details |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListIAMMemberRolesRequest


 <!-- end HasFields -->


#### ListIAMMemberRolesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| roles | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListIAMMembersRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.member.v1.SearchQuery](#zitadelmemberv1searchquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListIAMMembersResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.member.v1.Member](#zitadelmemberv1member) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListIDPsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| sorting_column | [ zitadel.idp.v1.IDPFieldName](#zitadelidpv1idpfieldname) | the field the result is sorted |
| queries | [repeated IDPQuery](#idpquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListIDPsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| sorting_column | [ zitadel.idp.v1.IDPFieldName](#zitadelidpv1idpfieldname) | - |
| result | [repeated zitadel.idp.v1.IDP](#zitadelidpv1idp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListLoginPolicyIDPsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListLoginPolicyIDPsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.idp.v1.IDPLoginPolicyLink](#zitadelidpv1idploginpolicylink) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListLoginPolicyMultiFactorsRequest


 <!-- end HasFields -->


#### ListLoginPolicyMultiFactorsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.policy.v1.MultiFactorType](#zitadelpolicyv1multifactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListLoginPolicySecondFactorsRequest


 <!-- end HasFields -->


#### ListLoginPolicySecondFactorsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.policy.v1.SecondFactorType](#zitadelpolicyv1secondfactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| sorting_column | [ zitadel.org.v1.OrgFieldName](#zitadelorgv1orgfieldname) | the field the result is sorted |
| queries | [repeated zitadel.org.v1.OrgQuery](#zitadelorgv1orgquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| sorting_column | [ zitadel.org.v1.OrgFieldName](#zitadelorgv1orgfieldname) | - |
| result | [repeated zitadel.org.v1.Org](#zitadelorgv1org) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListViewsRequest


 <!-- end HasFields -->


#### ListViewsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated View](#view) | TODO: list details |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveFailedEventRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| database | [ string](#string) | - |
| view_name | [ string](#string) | - |
| failed_sequence | [ uint64](#uint64) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveFailedEventResponse


 <!-- end HasFields -->


#### RemoveIAMMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveIAMMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveIDPFromLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveIDPFromLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMultiFactorFromLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ zitadel.policy.v1.MultiFactorType](#zitadelpolicyv1multifactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMultiFactorFromLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveSecondFactorFromLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ zitadel.policy.v1.SecondFactorType](#zitadelpolicyv1secondfactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveSecondFactorFromLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetCustomOrgIAMPolicyToDefaultRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetCustomOrgIAMPolicyToDefaultResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetOrgFeaturesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetOrgFeaturesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetDefaultFeaturesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| tier_name | [ string](#string) | - |
| description | [ string](#string) | - |
| audit_log_retention | [ google.protobuf.Duration](#googleprotobufduration) | - |
| login_policy_username_login | [ bool](#bool) | - |
| login_policy_registration | [ bool](#bool) | - |
| login_policy_idp | [ bool](#bool) | - |
| login_policy_factors | [ bool](#bool) | - |
| login_policy_passwordless | [ bool](#bool) | - |
| password_complexity_policy | [ bool](#bool) | - |
| label_policy | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetDefaultFeaturesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetOrgFeaturesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
| tier_name | [ string](#string) | - |
| description | [ string](#string) | - |
| state | [ zitadel.features.v1.FeaturesState](#zitadelfeaturesv1featuresstate) | - |
| state_description | [ string](#string) | - |
| audit_log_retention | [ google.protobuf.Duration](#googleprotobufduration) | - |
| login_policy_username_login | [ bool](#bool) | - |
| login_policy_registration | [ bool](#bool) | - |
| login_policy_idp | [ bool](#bool) | - |
| login_policy_factors | [ bool](#bool) | - |
| login_policy_passwordless | [ bool](#bool) | - |
| password_complexity_policy | [ bool](#bool) | - |
| label_policy | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetOrgFeaturesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetUpOrgRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| org | [ SetUpOrgRequest.Org](#setuporgrequestorg) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) user.human | [ SetUpOrgRequest.Human](#setuporgrequesthuman) | oneof field for the user managing the organisation |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetUpOrgRequest.Human



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name | [ string](#string) | - |
| profile | [ SetUpOrgRequest.Human.Profile](#setuporgrequesthumanprofile) | - |
| email | [ SetUpOrgRequest.Human.Email](#setuporgrequesthumanemail) | - |
| phone | [ SetUpOrgRequest.Human.Phone](#setuporgrequesthumanphone) | - |
| password | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetUpOrgRequest.Human.Email



| Field | Type | Description |
| ----- | ---- | ----------- |
| email | [ string](#string) | TODO: check if no value is allowed |
| is_email_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetUpOrgRequest.Human.Phone



| Field | Type | Description |
| ----- | ---- | ----------- |
| phone | [ string](#string) | has to be a global number |
| is_phone_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetUpOrgRequest.Human.Profile



| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name | [ string](#string) | - |
| last_name | [ string](#string) | - |
| nick_name | [ string](#string) | - |
| display_name | [ string](#string) | - |
| preferred_language | [ string](#string) | - |
| gender | [ zitadel.user.v1.Gender](#zitadeluserv1gender) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetUpOrgRequest.Org



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| domain | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetUpOrgResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| org_id | [ string](#string) | - |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomOrgIAMPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
| user_login_must_be_domain | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomOrgIAMPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateIAMMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | if no roles provided the user won't have any rights |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateIAMMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateIDPOIDCConfigRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
| issuer | [ string](#string) | - |
| client_id | [ string](#string) | - |
| client_secret | [ string](#string) | - |
| scopes | [repeated string](#string) | - |
| display_name_mapping | [ zitadel.idp.v1.OIDCMappingField](#zitadelidpv1oidcmappingfield) | - |
| username_mapping | [ zitadel.idp.v1.OIDCMappingField](#zitadelidpv1oidcmappingfield) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateIDPOIDCConfigResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
| name | [ string](#string) | - |
| styling_type | [ zitadel.idp.v1.IDPStylingType](#zitadelidpv1idpstylingtype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateLabelPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| primary_color | [ string](#string) | - |
| secondary_color | [ string](#string) | - |
| hide_login_name_suffix | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateLabelPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| allow_username_password | [ bool](#bool) | - |
| allow_register | [ bool](#bool) | - |
| allow_external_idp | [ bool](#bool) | - |
| force_mfa | [ bool](#bool) | - |
| passwordless_type | [ zitadel.policy.v1.PasswordlessType](#zitadelpolicyv1passwordlesstype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOrgIAMPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_login_must_be_domain | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOrgIAMPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdatePasswordAgePolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| max_age_days | [ uint32](#uint32) | - |
| expire_warn_days | [ uint32](#uint32) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdatePasswordAgePolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdatePasswordComplexityPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| min_length | [ uint32](#uint32) | - |
| has_uppercase | [ bool](#bool) | - |
| has_lowercase | [ bool](#bool) | - |
| has_number | [ bool](#bool) | - |
| has_symbol | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdatePasswordComplexityPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdatePasswordLockoutPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| max_attempts | [ uint32](#uint32) | failed attempts until a user gets locked |
| show_lockout_failure | [ bool](#bool) | If an error should be displayed during a lockout or not |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdatePasswordLockoutPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### View



| Field | Type | Description |
| ----- | ---- | ----------- |
| database | [ string](#string) | - |
| view_name | [ string](#string) | - |
| processed_sequence | [ uint64](#uint64) | - |
| event_timestamp | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | The timestamp the event occured |
| last_successful_spooler_run | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->

 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### APIConfig



| Field | Type | Description |
| ----- | ---- | ----------- |
| client_id | [ string](#string) | - |
| client_secret | [ string](#string) | - |
| auth_method_type | [ APIAuthMethodType](#apiauthmethodtype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### App



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| state | [ AppState](#appstate) | - |
| name | [ string](#string) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.oidc_config | [ OIDCConfig](#oidcconfig) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.api_config | [ APIConfig](#apiconfig) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AppNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AppQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.name_query | [ AppNameQuery](#appnamequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### OIDCConfig



| Field | Type | Description |
| ----- | ---- | ----------- |
| redirect_uris | [repeated string](#string) | - |
| response_types | [repeated OIDCResponseType](#oidcresponsetype) | - |
| grant_types | [repeated OIDCGrantType](#oidcgranttype) | - |
| app_type | [ OIDCAppType](#oidcapptype) | - |
| client_id | [ string](#string) | - |
| client_secret | [ string](#string) | - |
| auth_method_type | [ OIDCAuthMethodType](#oidcauthmethodtype) | - |
| post_logout_redirect_uris | [repeated string](#string) | - |
| version | [ OIDCVersion](#oidcversion) | - |
| none_compliant | [ bool](#bool) | - |
| compliance_problems | [repeated zitadel.v1.LocalizedMessage](#zitadelv1localizedmessage) | - |
| dev_mode | [ bool](#bool) | - |
| access_token_type | [ OIDCTokenType](#oidctokentype) | - |
| access_token_role_assertion | [ bool](#bool) | - |
| id_token_role_assertion | [ bool](#bool) | - |
| id_token_userinfo_assertion | [ bool](#bool) | - |
| clock_skew | [ google.protobuf.Duration](#googleprotobufduration) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### APIAuthMethodType {#apiauthmethodtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| API_AUTH_METHOD_TYPE_BASIC | 0 | - |
| API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT | 1 | - |




#### AppState {#appstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| APP_STATE_UNSPECIFIED | 0 | - |
| APP_STATE_ACTIVE | 1 | - |
| APP_STATE_INACTIVE | 2 | - |




#### OIDCAppType {#oidcapptype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_APP_TYPE_WEB | 0 | - |
| OIDC_APP_TYPE_USER_AGENT | 1 | - |
| OIDC_APP_TYPE_NATIVE | 2 | - |




#### OIDCAuthMethodType {#oidcauthmethodtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_AUTH_METHOD_TYPE_BASIC | 0 | - |
| OIDC_AUTH_METHOD_TYPE_POST | 1 | - |
| OIDC_AUTH_METHOD_TYPE_NONE | 2 | - |
| OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT | 3 | - |




#### OIDCGrantType {#oidcgranttype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_GRANT_TYPE_AUTHORIZATION_CODE | 0 | - |
| OIDC_GRANT_TYPE_IMPLICIT | 1 | - |
| OIDC_GRANT_TYPE_REFRESH_TOKEN | 2 | - |




#### OIDCResponseType {#oidcresponsetype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_RESPONSE_TYPE_CODE | 0 | - |
| OIDC_RESPONSE_TYPE_ID_TOKEN | 1 | - |
| OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN | 2 | - |




#### OIDCTokenType {#oidctokentype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_TOKEN_TYPE_BEARER | 0 | - |
| OIDC_TOKEN_TYPE_JWT | 1 | - |




#### OIDCVersion {#oidcversion}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_VERSION_1_0 | 0 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


### AuthService {#zitadelauthv1authservice}


#### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)




    Key:google.api.http
    Value:{[{GET /healthz }]}
 <!-- end options -->

#### GetMyUser

> **rpc** GetMyUser([GetMyUserRequest](#getmyuserrequest))
[GetMyUserResponse](#getmyuserresponse)




    Key:google.api.http
    Value:{[{GET /users/me }]}
 <!-- end options -->

#### ListMyUserChanges

> **rpc** ListMyUserChanges([ListMyUserChangesRequest](#listmyuserchangesrequest))
[ListMyUserChangesResponse](#listmyuserchangesresponse)




    Key:google.api.http
    Value:{[{POST /users/me/changes/_search }]}
 <!-- end options -->

#### ListMyUserSessions

> **rpc** ListMyUserSessions([ListMyUserSessionsRequest](#listmyusersessionsrequest))
[ListMyUserSessionsResponse](#listmyusersessionsresponse)




    Key:google.api.http
    Value:{[{POST /users/me/sessions/_search }]}
 <!-- end options -->

#### UpdateMyUserName

> **rpc** UpdateMyUserName([UpdateMyUserNameRequest](#updatemyusernamerequest))
[UpdateMyUserNameResponse](#updatemyusernameresponse)




    Key:google.api.http
    Value:{[{PUT /users/me/username *}]}
 <!-- end options -->

#### GetMyPasswordComplexityPolicy

> **rpc** GetMyPasswordComplexityPolicy([GetMyPasswordComplexityPolicyRequest](#getmypasswordcomplexitypolicyrequest))
[GetMyPasswordComplexityPolicyResponse](#getmypasswordcomplexitypolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/passwords/complexity }]}
 <!-- end options -->

#### UpdateMyPassword

> **rpc** UpdateMyPassword([UpdateMyPasswordRequest](#updatemypasswordrequest))
[UpdateMyPasswordResponse](#updatemypasswordresponse)




    Key:google.api.http
    Value:{[{PUT /users/me/password *}]}
 <!-- end options -->

#### GetMyProfile

> **rpc** GetMyProfile([GetMyProfileRequest](#getmyprofilerequest))
[GetMyProfileResponse](#getmyprofileresponse)




    Key:google.api.http
    Value:{[{GET /users/me/profile }]}
 <!-- end options -->

#### UpdateMyProfile

> **rpc** UpdateMyProfile([UpdateMyProfileRequest](#updatemyprofilerequest))
[UpdateMyProfileResponse](#updatemyprofileresponse)




    Key:google.api.http
    Value:{[{PUT /users/me/profile *}]}
 <!-- end options -->

#### GetMyEmail

> **rpc** GetMyEmail([GetMyEmailRequest](#getmyemailrequest))
[GetMyEmailResponse](#getmyemailresponse)




    Key:google.api.http
    Value:{[{GET /users/me/email }]}
 <!-- end options -->

#### SetMyEmail

> **rpc** SetMyEmail([SetMyEmailRequest](#setmyemailrequest))
[SetMyEmailResponse](#setmyemailresponse)




    Key:google.api.http
    Value:{[{PUT /users/me/email *}]}
 <!-- end options -->

#### VerifyMyEmail

> **rpc** VerifyMyEmail([VerifyMyEmailRequest](#verifymyemailrequest))
[VerifyMyEmailResponse](#verifymyemailresponse)




    Key:google.api.http
    Value:{[{POST /users/me/email/_verify *}]}
 <!-- end options -->

#### ResendMyEmailVerification

> **rpc** ResendMyEmailVerification([ResendMyEmailVerificationRequest](#resendmyemailverificationrequest))
[ResendMyEmailVerificationResponse](#resendmyemailverificationresponse)




    Key:google.api.http
    Value:{[{POST /users/me/email/_resend_verification *}]}
 <!-- end options -->

#### GetMyPhone

> **rpc** GetMyPhone([GetMyPhoneRequest](#getmyphonerequest))
[GetMyPhoneResponse](#getmyphoneresponse)




    Key:google.api.http
    Value:{[{GET /users/me/phone }]}
 <!-- end options -->

#### SetMyPhone

> **rpc** SetMyPhone([SetMyPhoneRequest](#setmyphonerequest))
[SetMyPhoneResponse](#setmyphoneresponse)




    Key:google.api.http
    Value:{[{PUT /users/me/phone *}]}
 <!-- end options -->

#### VerifyMyPhone

> **rpc** VerifyMyPhone([VerifyMyPhoneRequest](#verifymyphonerequest))
[VerifyMyPhoneResponse](#verifymyphoneresponse)




    Key:google.api.http
    Value:{[{POST /users/me/phone/_verify *}]}
 <!-- end options -->

#### ResendMyPhoneVerification

> **rpc** ResendMyPhoneVerification([ResendMyPhoneVerificationRequest](#resendmyphoneverificationrequest))
[ResendMyPhoneVerificationResponse](#resendmyphoneverificationresponse)




    Key:google.api.http
    Value:{[{POST /users/me/phone/_resend_verification *}]}
 <!-- end options -->

#### RemoveMyPhone

> **rpc** RemoveMyPhone([RemoveMyPhoneRequest](#removemyphonerequest))
[RemoveMyPhoneResponse](#removemyphoneresponse)




    Key:google.api.http
    Value:{[{DELETE /users/me/phone }]}
 <!-- end options -->

#### ListMyLinkedIDPs

> **rpc** ListMyLinkedIDPs([ListMyLinkedIDPsRequest](#listmylinkedidpsrequest))
[ListMyLinkedIDPsResponse](#listmylinkedidpsresponse)




    Key:google.api.http
    Value:{[{POST /users/me/idps/_search *}]}
 <!-- end options -->

#### RemoveMyLinkedIDP

> **rpc** RemoveMyLinkedIDP([RemoveMyLinkedIDPRequest](#removemylinkedidprequest))
[RemoveMyLinkedIDPResponse](#removemylinkedidpresponse)




    Key:google.api.http
    Value:{[{DELETE /users/me/idps/{idp_id}/{linked_user_id} }]}
 <!-- end options -->

#### ListMyAuthFactors

> **rpc** ListMyAuthFactors([ListMyAuthFactorsRequest](#listmyauthfactorsrequest))
[ListMyAuthFactorsResponse](#listmyauthfactorsresponse)




    Key:google.api.http
    Value:{[{POST /users/me/auth_factors/_search }]}
 <!-- end options -->

#### AddMyAuthFactorOTP

> **rpc** AddMyAuthFactorOTP([AddMyAuthFactorOTPRequest](#addmyauthfactorotprequest))
[AddMyAuthFactorOTPResponse](#addmyauthfactorotpresponse)




    Key:google.api.http
    Value:{[{POST /users/me/auth_factors/otp *}]}
 <!-- end options -->

#### VerifyMyAuthFactorOTP

> **rpc** VerifyMyAuthFactorOTP([VerifyMyAuthFactorOTPRequest](#verifymyauthfactorotprequest))
[VerifyMyAuthFactorOTPResponse](#verifymyauthfactorotpresponse)




    Key:google.api.http
    Value:{[{POST /users/me/auth_factors/otp/_verify *}]}
 <!-- end options -->

#### RemoveMyAuthFactorOTP

> **rpc** RemoveMyAuthFactorOTP([RemoveMyAuthFactorOTPRequest](#removemyauthfactorotprequest))
[RemoveMyAuthFactorOTPResponse](#removemyauthfactorotpresponse)




    Key:google.api.http
    Value:{[{DELETE /users/me/auth_factors/otp }]}
 <!-- end options -->

#### AddMyAuthFactorU2F

> **rpc** AddMyAuthFactorU2F([AddMyAuthFactorU2FRequest](#addmyauthfactoru2frequest))
[AddMyAuthFactorU2FResponse](#addmyauthfactoru2fresponse)




    Key:google.api.http
    Value:{[{POST /users/me/auth_factors/u2f *}]}
 <!-- end options -->

#### VerifyMyAuthFactorU2F

> **rpc** VerifyMyAuthFactorU2F([VerifyMyAuthFactorU2FRequest](#verifymyauthfactoru2frequest))
[VerifyMyAuthFactorU2FResponse](#verifymyauthfactoru2fresponse)




    Key:google.api.http
    Value:{[{POST /users/me/auth_factors/u2f/_verify *}]}
 <!-- end options -->

#### RemoveMyAuthFactorU2F

> **rpc** RemoveMyAuthFactorU2F([RemoveMyAuthFactorU2FRequest](#removemyauthfactoru2frequest))
[RemoveMyAuthFactorU2FResponse](#removemyauthfactoru2fresponse)




    Key:google.api.http
    Value:{[{DELETE /users/me/auth_factors/u2f/{token_id} }]}
 <!-- end options -->

#### ListMyPasswordless

> **rpc** ListMyPasswordless([ListMyPasswordlessRequest](#listmypasswordlessrequest))
[ListMyPasswordlessResponse](#listmypasswordlessresponse)




    Key:google.api.http
    Value:{[{POST /users/me/passwordless/_search }]}
 <!-- end options -->

#### AddMyPasswordless

> **rpc** AddMyPasswordless([AddMyPasswordlessRequest](#addmypasswordlessrequest))
[AddMyPasswordlessResponse](#addmypasswordlessresponse)




    Key:google.api.http
    Value:{[{POST /users/me/passwordless *}]}
 <!-- end options -->

#### VerifyMyPasswordless

> **rpc** VerifyMyPasswordless([VerifyMyPasswordlessRequest](#verifymypasswordlessrequest))
[VerifyMyPasswordlessResponse](#verifymypasswordlessresponse)




    Key:google.api.http
    Value:{[{POST /users/me/passwordless/_verify *}]}
 <!-- end options -->

#### RemoveMyPasswordless

> **rpc** RemoveMyPasswordless([RemoveMyPasswordlessRequest](#removemypasswordlessrequest))
[RemoveMyPasswordlessResponse](#removemypasswordlessresponse)




    Key:google.api.http
    Value:{[{DELETE /users/me/passwordless/{token_id} }]}
 <!-- end options -->

#### ListMyUserGrants

> **rpc** ListMyUserGrants([ListMyUserGrantsRequest](#listmyusergrantsrequest))
[ListMyUserGrantsResponse](#listmyusergrantsresponse)




    Key:google.api.http
    Value:{[{POST /usergrants/me/_search *}]}
 <!-- end options -->

#### ListMyProjectOrgs

> **rpc** ListMyProjectOrgs([ListMyProjectOrgsRequest](#listmyprojectorgsrequest))
[ListMyProjectOrgsResponse](#listmyprojectorgsresponse)




    Key:google.api.http
    Value:{[{POST /global/projectorgs/_search *}]}
 <!-- end options -->

#### ListMyZitadelFeatures

> **rpc** ListMyZitadelFeatures([ListMyZitadelFeaturesRequest](#listmyzitadelfeaturesrequest))
[ListMyZitadelFeaturesResponse](#listmyzitadelfeaturesresponse)




    Key:google.api.http
    Value:{[{POST /features/zitadel/me/_search }]}
 <!-- end options -->

#### ListMyZitadelPermissions

> **rpc** ListMyZitadelPermissions([ListMyZitadelPermissionsRequest](#listmyzitadelpermissionsrequest))
[ListMyZitadelPermissionsResponse](#listmyzitadelpermissionsresponse)




    Key:google.api.http
    Value:{[{POST /permissions/zitadel/me/_search }]}
 <!-- end options -->

#### ListMyProjectPermissions

> **rpc** ListMyProjectPermissions([ListMyProjectPermissionsRequest](#listmyprojectpermissionsrequest))
[ListMyProjectPermissionsResponse](#listmyprojectpermissionsresponse)




    Key:google.api.http
    Value:{[{POST /permissions/me/_search }]}
 <!-- end options -->

 <!-- end methods -->
 <!-- end services -->

### Messages


#### AddMyAuthFactorOTPRequest


 <!-- end HasFields -->


#### AddMyAuthFactorOTPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| url | [ string](#string) | - |
| secret | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMyAuthFactorU2FRequest


 <!-- end HasFields -->


#### AddMyAuthFactorU2FResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ zitadel.user.v1.WebAuthNKey](#zitadeluserv1webauthnkey) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMyPasswordlessRequest


 <!-- end HasFields -->


#### AddMyPasswordlessResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ zitadel.user.v1.WebAuthNKey](#zitadeluserv1webauthnkey) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetMyEmailRequest


 <!-- end HasFields -->


#### GetMyEmailResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| email | [ zitadel.user.v1.Email](#zitadeluserv1email) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetMyPasswordComplexityPolicyRequest


 <!-- end HasFields -->


#### GetMyPasswordComplexityPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordComplexityPolicy](#zitadelpolicyv1passwordcomplexitypolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetMyPhoneRequest


 <!-- end HasFields -->


#### GetMyPhoneResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| phone | [ zitadel.user.v1.Phone](#zitadeluserv1phone) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetMyProfileRequest


 <!-- end HasFields -->


#### GetMyProfileResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| profile | [ zitadel.user.v1.Profile](#zitadeluserv1profile) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetMyUserRequest
GetMyUserRequest is an empty request
the request parameters are read from the token-header

 <!-- end HasFields -->


#### GetMyUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| user | [ zitadel.user.v1.User](#zitadeluserv1user) | - |
| last_login | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### HealthzRequest


 <!-- end HasFields -->


#### HealthzResponse


 <!-- end HasFields -->


#### ListMyAuthFactorsRequest


 <!-- end HasFields -->


#### ListMyAuthFactorsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated zitadel.user.v1.AuthFactor](#zitadeluserv1authfactor) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyLinkedIDPsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering

PLANNED: queries for idp name and login name |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyLinkedIDPsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.idp.v1.IDPUserLink](#zitadelidpv1idpuserlink) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyPasswordlessRequest


 <!-- end HasFields -->


#### ListMyPasswordlessResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated zitadel.user.v1.WebAuthNToken](#zitadeluserv1webauthntoken) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyProjectOrgsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.org.v1.OrgQuery](#zitadelorgv1orgquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyProjectOrgsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.org.v1.Org](#zitadelorgv1org) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyProjectPermissionsRequest


 <!-- end HasFields -->


#### ListMyProjectPermissionsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyUserChangesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.change.v1.ChangeQuery](#zitadelchangev1changequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyUserChangesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.change.v1.Change](#zitadelchangev1change) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyUserGrantsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyUserGrantsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated UserGrant](#usergrant) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyUserSessionsRequest


 <!-- end HasFields -->


#### ListMyUserSessionsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated zitadel.user.v1.Session](#zitadeluserv1session) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyZitadelFeaturesRequest


 <!-- end HasFields -->


#### ListMyZitadelFeaturesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMyZitadelPermissionsRequest


 <!-- end HasFields -->


#### ListMyZitadelPermissionsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMyAuthFactorOTPRequest


 <!-- end HasFields -->


#### RemoveMyAuthFactorOTPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMyAuthFactorU2FRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| token_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMyAuthFactorU2FResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMyLinkedIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
| linked_user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMyLinkedIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMyPasswordlessRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| token_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMyPasswordlessResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMyPhoneRequest


 <!-- end HasFields -->


#### RemoveMyPhoneResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResendMyEmailVerificationRequest


 <!-- end HasFields -->


#### ResendMyEmailVerificationResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResendMyPhoneVerificationRequest


 <!-- end HasFields -->


#### ResendMyPhoneVerificationResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetMyEmailRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| email | [ string](#string) | TODO: check if no value is allowed |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetMyEmailResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetMyPhoneRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| phone | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetMyPhoneResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateMyPasswordRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| old_password | [ string](#string) | - |
| new_password | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateMyPasswordResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateMyProfileRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name | [ string](#string) | - |
| last_name | [ string](#string) | - |
| nick_name | [ string](#string) | - |
| display_name | [ string](#string) | - |
| preferred_language | [ string](#string) | - |
| gender | [ zitadel.user.v1.Gender](#zitadeluserv1gender) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateMyProfileResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateMyUserNameRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateMyUserNameResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrant



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
| project_id | [ string](#string) | - |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | - |
| org_name | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyAuthFactorOTPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| code | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyAuthFactorOTPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyAuthFactorU2FRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| verification | [ zitadel.user.v1.WebAuthNVerification](#zitadeluserv1webauthnverification) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyAuthFactorU2FResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyEmailRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| code | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyEmailResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyPasswordlessRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| verification | [ zitadel.user.v1.WebAuthNVerification](#zitadeluserv1webauthnverification) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyPasswordlessResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyPhoneRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| code | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### VerifyMyPhoneResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->

 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### Key



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| type | [ KeyType](#keytype) | - |
| expiration_date | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### KeyType {#keytype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| KEY_TYPE_UNSPECIFIED | 0 | - |
| KEY_TYPE_JSON | 1 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### Change



| Field | Type | Description |
| ----- | ---- | ----------- |
| change_date | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | - |
| event_type | [ zitadel.v1.LocalizedMessage](#zitadelv1localizedmessage) | - |
| sequence | [ uint64](#uint64) | - |
| editor_id | [ string](#string) | - |
| editor_display_name | [ string](#string) | - |
| resource_owner_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ChangeQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| sequence | [ uint64](#uint64) | sequence represents the order of events. It's always upcounting |
| limit | [ uint32](#uint32) | - |
| asc | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->

 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### FeatureTier



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| description | [ string](#string) | - |
| state | [ FeaturesState](#featuresstate) | - |
| status_info | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Features



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| tier | [ FeatureTier](#featuretier) | - |
| is_default | [ bool](#bool) | - |
| audit_log_retention | [ google.protobuf.Duration](#googleprotobufduration) | - |
| login_policy_username_login | [ bool](#bool) | - |
| login_policy_registration | [ bool](#bool) | - |
| login_policy_idp | [ bool](#bool) | - |
| login_policy_factors | [ bool](#bool) | - |
| login_policy_passwordless | [ bool](#bool) | - |
| password_complexity_policy | [ bool](#bool) | - |
| label_policy | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### FeaturesState {#featuresstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| FEATURES_STATE_ACTIVE | 0 | - |
| FEATURES_STATE_ACTION_REQUIRED | 1 | - |
| FEATURES_STATE_CANCELED | 2 | - |
| FEATURES_STATE_GRANDFATHERED | 3 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### IDP



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| state | [ IDPState](#idpstate) | - |
| name | [ string](#string) | - |
| styling_type | [ IDPStylingType](#idpstylingtype) | - |
| owner | [ IDPOwnerType](#idpownertype) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) config.oidc_config | [ OIDCConfig](#oidcconfig) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IDPIDQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IDPLoginPolicyLink



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
| idp_name | [ string](#string) | - |
| idp_type | [ IDPType](#idptype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IDPNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IDPOwnerTypeQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| owner_type | [ IDPOwnerType](#idpownertype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IDPUserLink



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| idp_id | [ string](#string) | - |
| idp_name | [ string](#string) | - |
| provided_user_id | [ string](#string) | - |
| provided_user_name | [ string](#string) | - |
| idp_type | [ IDPType](#idptype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### OIDCConfig



| Field | Type | Description |
| ----- | ---- | ----------- |
| client_id | [ string](#string) | - |
| issuer | [ string](#string) | - |
| scopes | [repeated string](#string) | - |
| display_name_mapping | [ OIDCMappingField](#oidcmappingfield) | - |
| username_mapping | [ OIDCMappingField](#oidcmappingfield) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### IDPFieldName {#idpfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_FIELD_NAME_UNSPECIFIED | 0 | - |
| IDP_FIELD_NAME_NAME | 1 | - |




#### IDPOwnerType {#idpownertype}
the owner of the identity provider.

| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_OWNER_TYPE_UNSPECIFIED | 0 | - |
| IDP_OWNER_TYPE_SYSTEM | 1 | system is managed by the ZITADEL administrators |
| IDP_OWNER_TYPE_ORG | 2 | org is managed by de organisation administrators |




#### IDPState {#idpstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_STATE_UNSPECIFIED | 0 | - |
| IDP_STATE_ACTIVE | 1 | - |
| IDP_STATE_INACTIVE | 2 | - |




#### IDPStylingType {#idpstylingtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| STYLING_TYPE_UNSPECIFIED | 0 | - |
| STYLING_TYPE_GOOGLE | 1 | - |




#### IDPType {#idptype}
authorization framework of the identity provider

| Name | Number | Description |
| ---- | ------ | ----------- |
| IDP_TYPE_UNSPECIFIED | 0 | - |
| IDP_TYPE_OIDC | 1 | PLANNED: IDP_TYPE_SAML |




#### OIDCMappingField {#oidcmappingfield}


| Name | Number | Description |
| ---- | ------ | ----------- |
| OIDC_MAPPING_FIELD_UNSPECIFIED | 0 | - |
| OIDC_MAPPING_FIELD_PREFERRED_USERNAME | 1 | - |
| OIDC_MAPPING_FIELD_EMAIL | 2 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


### ManagementService {#zitadelmanagementv1managementservice}


#### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)




    Key:google.api.http
    Value:{[{GET /healthz }]}
 <!-- end options -->

#### GetOIDCInformation

> **rpc** GetOIDCInformation([GetOIDCInformationRequest](#getoidcinformationrequest))
[GetOIDCInformationResponse](#getoidcinformationresponse)




    Key:google.api.http
    Value:{[{GET /zitadel/docs }]}
 <!-- end options -->

#### GetIAM

> **rpc** GetIAM([GetIAMRequest](#getiamrequest))
[GetIAMResponse](#getiamresponse)

GetIam returns some needed settings of the iam (Global Organisation ID, Zitadel Project ID)


    Key:google.api.http
    Value:{[{GET /iam }]}
 <!-- end options -->

#### GetUserByID

> **rpc** GetUserByID([GetUserByIDRequest](#getuserbyidrequest))
[GetUserByIDResponse](#getuserbyidresponse)




    Key:google.api.http
    Value:{[{GET /users/{id} }]}
 <!-- end options -->

#### GetUserByLoginNameGlobal

> **rpc** GetUserByLoginNameGlobal([GetUserByLoginNameGlobalRequest](#getuserbyloginnameglobalrequest))
[GetUserByLoginNameGlobalResponse](#getuserbyloginnameglobalresponse)

GetUserByLoginNameGlobal searches a user over all organisations
the login name has to match exactly


    Key:google.api.http
    Value:{[{GET /global/users/_by_login_name }]}
 <!-- end options -->

#### ListUsers

> **rpc** ListUsers([ListUsersRequest](#listusersrequest))
[ListUsersResponse](#listusersresponse)

Limit should always be set, there is a default limit set by the service


    Key:google.api.http
    Value:{[{POST /users/_search *}]}
 <!-- end options -->

#### ListUserChanges

> **rpc** ListUserChanges([ListUserChangesRequest](#listuserchangesrequest))
[ListUserChangesResponse](#listuserchangesresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/changes/_search *}]}
 <!-- end options -->

#### IsUserUnique

> **rpc** IsUserUnique([IsUserUniqueRequest](#isuseruniquerequest))
[IsUserUniqueResponse](#isuseruniqueresponse)




    Key:google.api.http
    Value:{[{GET /users/_is_unique }]}
 <!-- end options -->

#### AddHumanUser

> **rpc** AddHumanUser([AddHumanUserRequest](#addhumanuserrequest))
[AddHumanUserResponse](#addhumanuserresponse)




    Key:google.api.http
    Value:{[{POST /users/human *}]}
 <!-- end options -->

#### ImportHumanUser

> **rpc** ImportHumanUser([ImportHumanUserRequest](#importhumanuserrequest))
[ImportHumanUserResponse](#importhumanuserresponse)




    Key:google.api.http
    Value:{[{POST /users/human/_import *}]}
 <!-- end options -->

#### AddMachineUser

> **rpc** AddMachineUser([AddMachineUserRequest](#addmachineuserrequest))
[AddMachineUserResponse](#addmachineuserresponse)




    Key:google.api.http
    Value:{[{POST /users/machine *}]}
 <!-- end options -->

#### DeactivateUser

> **rpc** DeactivateUser([DeactivateUserRequest](#deactivateuserrequest))
[DeactivateUserResponse](#deactivateuserresponse)




    Key:google.api.http
    Value:{[{POST /users/{id}/_deactivate *}]}
 <!-- end options -->

#### ReactivateUser

> **rpc** ReactivateUser([ReactivateUserRequest](#reactivateuserrequest))
[ReactivateUserResponse](#reactivateuserresponse)




    Key:google.api.http
    Value:{[{POST /users/{id}/_reactivate *}]}
 <!-- end options -->

#### LockUser

> **rpc** LockUser([LockUserRequest](#lockuserrequest))
[LockUserResponse](#lockuserresponse)




    Key:google.api.http
    Value:{[{POST /users/{id}/_lock *}]}
 <!-- end options -->

#### UnlockUser

> **rpc** UnlockUser([UnlockUserRequest](#unlockuserrequest))
[UnlockUserResponse](#unlockuserresponse)




    Key:google.api.http
    Value:{[{POST /users/{id}/_unlock *}]}
 <!-- end options -->

#### RemoveUser

> **rpc** RemoveUser([RemoveUserRequest](#removeuserrequest))
[RemoveUserResponse](#removeuserresponse)




    Key:google.api.http
    Value:{[{DELETE /users/{id} }]}
 <!-- end options -->

#### UpdateUserName

> **rpc** UpdateUserName([UpdateUserNameRequest](#updateusernamerequest))
[UpdateUserNameResponse](#updateusernameresponse)




    Key:google.api.http
    Value:{[{GET /users/{user_id}/username }]}
 <!-- end options -->

#### GetHumanProfile

> **rpc** GetHumanProfile([GetHumanProfileRequest](#gethumanprofilerequest))
[GetHumanProfileResponse](#gethumanprofileresponse)




    Key:google.api.http
    Value:{[{GET /users/{user_id}/profile }]}
 <!-- end options -->

#### UpdateHumanProfile

> **rpc** UpdateHumanProfile([UpdateHumanProfileRequest](#updatehumanprofilerequest))
[UpdateHumanProfileResponse](#updatehumanprofileresponse)




    Key:google.api.http
    Value:{[{PUT /users/{user_id}/profile *}]}
 <!-- end options -->

#### GetHumanEmail

> **rpc** GetHumanEmail([GetHumanEmailRequest](#gethumanemailrequest))
[GetHumanEmailResponse](#gethumanemailresponse)




    Key:google.api.http
    Value:{[{GET /users/{user_id}/email }]}
 <!-- end options -->

#### UpdateHumanEmail

> **rpc** UpdateHumanEmail([UpdateHumanEmailRequest](#updatehumanemailrequest))
[UpdateHumanEmailResponse](#updatehumanemailresponse)




    Key:google.api.http
    Value:{[{PUT /users/{user_id}/email *}]}
 <!-- end options -->

#### ResendHumanInitialization

> **rpc** ResendHumanInitialization([ResendHumanInitializationRequest](#resendhumaninitializationrequest))
[ResendHumanInitializationResponse](#resendhumaninitializationresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/_resend_initialization *}]}
 <!-- end options -->

#### ResendHumanEmailVerification

> **rpc** ResendHumanEmailVerification([ResendHumanEmailVerificationRequest](#resendhumanemailverificationrequest))
[ResendHumanEmailVerificationResponse](#resendhumanemailverificationresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/email/_resend_verification *}]}
 <!-- end options -->

#### GetHumanPhone

> **rpc** GetHumanPhone([GetHumanPhoneRequest](#gethumanphonerequest))
[GetHumanPhoneResponse](#gethumanphoneresponse)




    Key:google.api.http
    Value:{[{GET /users/{user_id}/phone }]}
 <!-- end options -->

#### UpdateHumanPhone

> **rpc** UpdateHumanPhone([UpdateHumanPhoneRequest](#updatehumanphonerequest))
[UpdateHumanPhoneResponse](#updatehumanphoneresponse)




    Key:google.api.http
    Value:{[{PUT /users/{user_id}/phone *}]}
 <!-- end options -->

#### RemoveHumanPhone

> **rpc** RemoveHumanPhone([RemoveHumanPhoneRequest](#removehumanphonerequest))
[RemoveHumanPhoneResponse](#removehumanphoneresponse)




    Key:google.api.http
    Value:{[{DELETE /users/{user_id}/phone }]}
 <!-- end options -->

#### ResendHumanPhoneVerification

> **rpc** ResendHumanPhoneVerification([ResendHumanPhoneVerificationRequest](#resendhumanphoneverificationrequest))
[ResendHumanPhoneVerificationResponse](#resendhumanphoneverificationresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/phone/_resend_verification *}]}
 <!-- end options -->

#### SetHumanInitialPassword

> **rpc** SetHumanInitialPassword([SetHumanInitialPasswordRequest](#sethumaninitialpasswordrequest))
[SetHumanInitialPasswordResponse](#sethumaninitialpasswordresponse)

A Manager is only allowed to set an initial password, on the next login the user has to change his password


    Key:google.api.http
    Value:{[{POST /users/{user_id}/password/_initialize *}]}
 <!-- end options -->

#### SendHumanResetPasswordNotification

> **rpc** SendHumanResetPasswordNotification([SendHumanResetPasswordNotificationRequest](#sendhumanresetpasswordnotificationrequest))
[SendHumanResetPasswordNotificationResponse](#sendhumanresetpasswordnotificationresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/password/_reset *}]}
 <!-- end options -->

#### ListHumanAuthFactors

> **rpc** ListHumanAuthFactors([ListHumanAuthFactorsRequest](#listhumanauthfactorsrequest))
[ListHumanAuthFactorsResponse](#listhumanauthfactorsresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/auth_factors/_search }]}
 <!-- end options -->

#### RemoveHumanAuthFactorOTP

> **rpc** RemoveHumanAuthFactorOTP([RemoveHumanAuthFactorOTPRequest](#removehumanauthfactorotprequest))
[RemoveHumanAuthFactorOTPResponse](#removehumanauthfactorotpresponse)




    Key:google.api.http
    Value:{[{DELETE /users/{user_id}/auth_factors/otp }]}
 <!-- end options -->

#### RemoveHumanAuthFactorU2F

> **rpc** RemoveHumanAuthFactorU2F([RemoveHumanAuthFactorU2FRequest](#removehumanauthfactoru2frequest))
[RemoveHumanAuthFactorU2FResponse](#removehumanauthfactoru2fresponse)




    Key:google.api.http
    Value:{[{DELETE /users/{user_id}/auth_factors/u2f/{token_id} }]}
 <!-- end options -->

#### ListHumanPasswordless

> **rpc** ListHumanPasswordless([ListHumanPasswordlessRequest](#listhumanpasswordlessrequest))
[ListHumanPasswordlessResponse](#listhumanpasswordlessresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/passwordless/_search }]}
 <!-- end options -->

#### RemoveHumanPasswordless

> **rpc** RemoveHumanPasswordless([RemoveHumanPasswordlessRequest](#removehumanpasswordlessrequest))
[RemoveHumanPasswordlessResponse](#removehumanpasswordlessresponse)




    Key:google.api.http
    Value:{[{DELETE /users/{user_id}/passwordless/{token_id} }]}
 <!-- end options -->

#### UpdateMachine

> **rpc** UpdateMachine([UpdateMachineRequest](#updatemachinerequest))
[UpdateMachineResponse](#updatemachineresponse)




    Key:google.api.http
    Value:{[{PUT /users/{user_id}/machine *}]}
 <!-- end options -->

#### GetMachineKeyByIDs

> **rpc** GetMachineKeyByIDs([GetMachineKeyByIDsRequest](#getmachinekeybyidsrequest))
[GetMachineKeyByIDsResponse](#getmachinekeybyidsresponse)




    Key:google.api.http
    Value:{[{GET /users/{user_id}/keys/{key_id} }]}
 <!-- end options -->

#### ListMachineKeys

> **rpc** ListMachineKeys([ListMachineKeysRequest](#listmachinekeysrequest))
[ListMachineKeysResponse](#listmachinekeysresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/keys/_search *}]}
 <!-- end options -->

#### AddMachineKey

> **rpc** AddMachineKey([AddMachineKeyRequest](#addmachinekeyrequest))
[AddMachineKeyResponse](#addmachinekeyresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/keys *}]}
 <!-- end options -->

#### RemoveMachineKey

> **rpc** RemoveMachineKey([RemoveMachineKeyRequest](#removemachinekeyrequest))
[RemoveMachineKeyResponse](#removemachinekeyresponse)




    Key:google.api.http
    Value:{[{DELETE /users/{user_id}/keys/{key_id} }]}
 <!-- end options -->

#### ListHumanLinkedIDPs

> **rpc** ListHumanLinkedIDPs([ListHumanLinkedIDPsRequest](#listhumanlinkedidpsrequest))
[ListHumanLinkedIDPsResponse](#listhumanlinkedidpsresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/idps/_search *}]}
 <!-- end options -->

#### RemoveHumanLinkedIDP

> **rpc** RemoveHumanLinkedIDP([RemoveHumanLinkedIDPRequest](#removehumanlinkedidprequest))
[RemoveHumanLinkedIDPResponse](#removehumanlinkedidpresponse)




    Key:google.api.http
    Value:{[{DELETE /users/{user_id}/idps/{idp_id}/{linked_user_id} }]}
 <!-- end options -->

#### ListUserMemberships

> **rpc** ListUserMemberships([ListUserMembershipsRequest](#listusermembershipsrequest))
[ListUserMembershipsResponse](#listusermembershipsresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/memberships/_search *}]}
 <!-- end options -->

#### GetMyOrg

> **rpc** GetMyOrg([GetMyOrgRequest](#getmyorgrequest))
[GetMyOrgResponse](#getmyorgresponse)




    Key:google.api.http
    Value:{[{GET /orgs/me }]}
 <!-- end options -->

#### GetOrgByDomainGlobal

> **rpc** GetOrgByDomainGlobal([GetOrgByDomainGlobalRequest](#getorgbydomainglobalrequest))
[GetOrgByDomainGlobalResponse](#getorgbydomainglobalresponse)




    Key:google.api.http
    Value:{[{GET /global/orgs/_by_domain }]}
 <!-- end options -->

#### ListOrgChanges

> **rpc** ListOrgChanges([ListOrgChangesRequest](#listorgchangesrequest))
[ListOrgChangesResponse](#listorgchangesresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/changes/_search *}]}
 <!-- end options -->

#### AddOrg

> **rpc** AddOrg([AddOrgRequest](#addorgrequest))
[AddOrgResponse](#addorgresponse)




    Key:google.api.http
    Value:{[{POST /orgs *}]}
 <!-- end options -->

#### DeactivateOrg

> **rpc** DeactivateOrg([DeactivateOrgRequest](#deactivateorgrequest))
[DeactivateOrgResponse](#deactivateorgresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/_deactivate *}]}
 <!-- end options -->

#### ReactivateOrg

> **rpc** ReactivateOrg([ReactivateOrgRequest](#reactivateorgrequest))
[ReactivateOrgResponse](#reactivateorgresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/_reactivate *}]}
 <!-- end options -->

#### ListOrgDomains

> **rpc** ListOrgDomains([ListOrgDomainsRequest](#listorgdomainsrequest))
[ListOrgDomainsResponse](#listorgdomainsresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/domains/_search *}]}
 <!-- end options -->

#### AddOrgDomain

> **rpc** AddOrgDomain([AddOrgDomainRequest](#addorgdomainrequest))
[AddOrgDomainResponse](#addorgdomainresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/domains *}]}
 <!-- end options -->

#### RemoveOrgDomain

> **rpc** RemoveOrgDomain([RemoveOrgDomainRequest](#removeorgdomainrequest))
[RemoveOrgDomainResponse](#removeorgdomainresponse)




    Key:google.api.http
    Value:{[{DELETE /orgs/me/domains/{domain} }]}
 <!-- end options -->

#### GenerateOrgDomainValidation

> **rpc** GenerateOrgDomainValidation([GenerateOrgDomainValidationRequest](#generateorgdomainvalidationrequest))
[GenerateOrgDomainValidationResponse](#generateorgdomainvalidationresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/domains/{domain}/validation/_generate *}]}
 <!-- end options -->

#### ValidateOrgDomain

> **rpc** ValidateOrgDomain([ValidateOrgDomainRequest](#validateorgdomainrequest))
[ValidateOrgDomainResponse](#validateorgdomainresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/domains/{domain}/validation/_validate *}]}
 <!-- end options -->

#### SetPrimaryOrgDomain

> **rpc** SetPrimaryOrgDomain([SetPrimaryOrgDomainRequest](#setprimaryorgdomainrequest))
[SetPrimaryOrgDomainResponse](#setprimaryorgdomainresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/domains/{domain}/_set_primary }]}
 <!-- end options -->

#### ListOrgMemberRoles

> **rpc** ListOrgMemberRoles([ListOrgMemberRolesRequest](#listorgmemberrolesrequest))
[ListOrgMemberRolesResponse](#listorgmemberrolesresponse)




    Key:google.api.http
    Value:{[{POST /orgs/members/roles/_search }]}
 <!-- end options -->

#### ListOrgMembers

> **rpc** ListOrgMembers([ListOrgMembersRequest](#listorgmembersrequest))
[ListOrgMembersResponse](#listorgmembersresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/members/_search *}]}
 <!-- end options -->

#### AddOrgMember

> **rpc** AddOrgMember([AddOrgMemberRequest](#addorgmemberrequest))
[AddOrgMemberResponse](#addorgmemberresponse)




    Key:google.api.http
    Value:{[{POST /orgs/me/members *}]}
 <!-- end options -->

#### UpdateOrgMember

> **rpc** UpdateOrgMember([UpdateOrgMemberRequest](#updateorgmemberrequest))
[UpdateOrgMemberResponse](#updateorgmemberresponse)




    Key:google.api.http
    Value:{[{PUT /orgs/me/members/{user_id} *}]}
 <!-- end options -->

#### RemoveOrgMember

> **rpc** RemoveOrgMember([RemoveOrgMemberRequest](#removeorgmemberrequest))
[RemoveOrgMemberResponse](#removeorgmemberresponse)




    Key:google.api.http
    Value:{[{DELETE /orgs/me/members/{user_id} }]}
 <!-- end options -->

#### GetProjectByID

> **rpc** GetProjectByID([GetProjectByIDRequest](#getprojectbyidrequest))
[GetProjectByIDResponse](#getprojectbyidresponse)




    Key:google.api.http
    Value:{[{GET /projects/{id} }]}
 <!-- end options -->

#### GetGrantedProjectByID

> **rpc** GetGrantedProjectByID([GetGrantedProjectByIDRequest](#getgrantedprojectbyidrequest))
[GetGrantedProjectByIDResponse](#getgrantedprojectbyidresponse)

returns a project my organisation got granted from another organisation


    Key:google.api.http
    Value:{[{GET /granted_projects/{project_id}/grants/{grant_id} }]}
 <!-- end options -->

#### ListProjects

> **rpc** ListProjects([ListProjectsRequest](#listprojectsrequest))
[ListProjectsResponse](#listprojectsresponse)




    Key:google.api.http
    Value:{[{POST /projects/_search *}]}
 <!-- end options -->

#### ListGrantedProjects

> **rpc** ListGrantedProjects([ListGrantedProjectsRequest](#listgrantedprojectsrequest))
[ListGrantedProjectsResponse](#listgrantedprojectsresponse)

returns all projects my organisation got granted from another organisation


    Key:google.api.http
    Value:{[{POST /granted_projects/_search *}]}
 <!-- end options -->

#### ListGrantedProjectRoles

> **rpc** ListGrantedProjectRoles([ListGrantedProjectRolesRequest](#listgrantedprojectrolesrequest))
[ListGrantedProjectRolesResponse](#listgrantedprojectrolesresponse)

returns all roles of a project grant


    Key:google.api.http
    Value:{[{GET /granted_projects/{project_id}/grants/{grant_id}/roles/_search }]}
 <!-- end options -->

#### ListProjectChanges

> **rpc** ListProjectChanges([ListProjectChangesRequest](#listprojectchangesrequest))
[ListProjectChangesResponse](#listprojectchangesresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/changes/_search }]}
 <!-- end options -->

#### AddProject

> **rpc** AddProject([AddProjectRequest](#addprojectrequest))
[AddProjectResponse](#addprojectresponse)




    Key:google.api.http
    Value:{[{POST /projects *}]}
 <!-- end options -->

#### UpdateProject

> **rpc** UpdateProject([UpdateProjectRequest](#updateprojectrequest))
[UpdateProjectResponse](#updateprojectresponse)




    Key:google.api.http
    Value:{[{PUT /projects/{id} *}]}
 <!-- end options -->

#### DeactivateProject

> **rpc** DeactivateProject([DeactivateProjectRequest](#deactivateprojectrequest))
[DeactivateProjectResponse](#deactivateprojectresponse)




    Key:google.api.http
    Value:{[{POST /projects/{id}/_deactivate *}]}
 <!-- end options -->

#### ReactivateProject

> **rpc** ReactivateProject([ReactivateProjectRequest](#reactivateprojectrequest))
[ReactivateProjectResponse](#reactivateprojectresponse)




    Key:google.api.http
    Value:{[{POST /projects/{id}/_reactivate *}]}
 <!-- end options -->

#### RemoveProject

> **rpc** RemoveProject([RemoveProjectRequest](#removeprojectrequest))
[RemoveProjectResponse](#removeprojectresponse)




    Key:google.api.http
    Value:{[{DELETE /projects/{id} }]}
 <!-- end options -->

#### ListProjectRoles

> **rpc** ListProjectRoles([ListProjectRolesRequest](#listprojectrolesrequest))
[ListProjectRolesResponse](#listprojectrolesresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/roles/_search *}]}
 <!-- end options -->

#### AddProjectRole

> **rpc** AddProjectRole([AddProjectRoleRequest](#addprojectrolerequest))
[AddProjectRoleResponse](#addprojectroleresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/roles *}]}
 <!-- end options -->

#### BulkAddProjectRoles

> **rpc** BulkAddProjectRoles([BulkAddProjectRolesRequest](#bulkaddprojectrolesrequest))
[BulkAddProjectRolesResponse](#bulkaddprojectrolesresponse)

add a list of project roles in one request


    Key:google.api.http
    Value:{[{POST /projects/{project_id}/roles/_bulk *}]}
 <!-- end options -->

#### UpdateProjectRole

> **rpc** UpdateProjectRole([UpdateProjectRoleRequest](#updateprojectrolerequest))
[UpdateProjectRoleResponse](#updateprojectroleresponse)




    Key:google.api.http
    Value:{[{PUT /projects/{project_id}/roles/{role_key} *}]}
 <!-- end options -->

#### RemoveProjectRole

> **rpc** RemoveProjectRole([RemoveProjectRoleRequest](#removeprojectrolerequest))
[RemoveProjectRoleResponse](#removeprojectroleresponse)

RemoveProjectRole removes role from UserGrants, ProjectGrants and from Project


    Key:google.api.http
    Value:{[{DELETE /projects/{project_id}/roles/{role_key} }]}
 <!-- end options -->

#### ListProjectMemberRoles

> **rpc** ListProjectMemberRoles([ListProjectMemberRolesRequest](#listprojectmemberrolesrequest))
[ListProjectMemberRolesResponse](#listprojectmemberrolesresponse)




    Key:google.api.http
    Value:{[{POST /projects/members/roles/_search }]}
 <!-- end options -->

#### ListProjectMembers

> **rpc** ListProjectMembers([ListProjectMembersRequest](#listprojectmembersrequest))
[ListProjectMembersResponse](#listprojectmembersresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/members/_search *}]}
 <!-- end options -->

#### AddProjectMember

> **rpc** AddProjectMember([AddProjectMemberRequest](#addprojectmemberrequest))
[AddProjectMemberResponse](#addprojectmemberresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/members *}]}
 <!-- end options -->

#### UpdateProjectMember

> **rpc** UpdateProjectMember([UpdateProjectMemberRequest](#updateprojectmemberrequest))
[UpdateProjectMemberResponse](#updateprojectmemberresponse)




    Key:google.api.http
    Value:{[{PUT /projects/{project_id}/members/{user_id} *}]}
 <!-- end options -->

#### RemoveProjectMember

> **rpc** RemoveProjectMember([RemoveProjectMemberRequest](#removeprojectmemberrequest))
[RemoveProjectMemberResponse](#removeprojectmemberresponse)




    Key:google.api.http
    Value:{[{DELETE /projects/{project_id}/members/{user_id} }]}
 <!-- end options -->

#### GetAppByID

> **rpc** GetAppByID([GetAppByIDRequest](#getappbyidrequest))
[GetAppByIDResponse](#getappbyidresponse)




    Key:google.api.http
    Value:{[{GET /projects/{project_id}/apps/{app_id} }]}
 <!-- end options -->

#### ListApps

> **rpc** ListApps([ListAppsRequest](#listappsrequest))
[ListAppsResponse](#listappsresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/_search *}]}
 <!-- end options -->

#### ListAppChanges

> **rpc** ListAppChanges([ListAppChangesRequest](#listappchangesrequest))
[ListAppChangesResponse](#listappchangesresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/{app_id}/changes/_search }]}
 <!-- end options -->

#### AddOIDCApp

> **rpc** AddOIDCApp([AddOIDCAppRequest](#addoidcapprequest))
[AddOIDCAppResponse](#addoidcappresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/oidc *}]}
 <!-- end options -->

#### AddAPIApp

> **rpc** AddAPIApp([AddAPIAppRequest](#addapiapprequest))
[AddAPIAppResponse](#addapiappresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/api *}]}
 <!-- end options -->

#### UpdateApp

> **rpc** UpdateApp([UpdateAppRequest](#updateapprequest))
[UpdateAppResponse](#updateappresponse)




    Key:google.api.http
    Value:{[{PUT /projects/{project_id}/apps/{app_id} *}]}
 <!-- end options -->

#### UpdateOIDCAppConfig

> **rpc** UpdateOIDCAppConfig([UpdateOIDCAppConfigRequest](#updateoidcappconfigrequest))
[UpdateOIDCAppConfigResponse](#updateoidcappconfigresponse)




    Key:google.api.http
    Value:{[{PUT /projects/{project_id}/apps/{app_id}/oidc_config *}]}
 <!-- end options -->

#### UpdateAPIAppConfig

> **rpc** UpdateAPIAppConfig([UpdateAPIAppConfigRequest](#updateapiappconfigrequest))
[UpdateAPIAppConfigResponse](#updateapiappconfigresponse)




    Key:google.api.http
    Value:{[{PUT /projects/{project_id}/apps/{app_id}/api_config *}]}
 <!-- end options -->

#### DeactivateApp

> **rpc** DeactivateApp([DeactivateAppRequest](#deactivateapprequest))
[DeactivateAppResponse](#deactivateappresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/{app_id}/_deactivate *}]}
 <!-- end options -->

#### ReactivateApp

> **rpc** ReactivateApp([ReactivateAppRequest](#reactivateapprequest))
[ReactivateAppResponse](#reactivateappresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/{app_id}/_reactivate *}]}
 <!-- end options -->

#### RemoveApp

> **rpc** RemoveApp([RemoveAppRequest](#removeapprequest))
[RemoveAppResponse](#removeappresponse)




    Key:google.api.http
    Value:{[{DELETE /projects/{project_id}/apps/{app_id} }]}
 <!-- end options -->

#### RegenerateOIDCClientSecret

> **rpc** RegenerateOIDCClientSecret([RegenerateOIDCClientSecretRequest](#regenerateoidcclientsecretrequest))
[RegenerateOIDCClientSecretResponse](#regenerateoidcclientsecretresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/{app_id}/oidc_config/_generate_client_secret *}]}
 <!-- end options -->

#### RegenerateAPIClientSecret

> **rpc** RegenerateAPIClientSecret([RegenerateAPIClientSecretRequest](#regenerateapiclientsecretrequest))
[RegenerateAPIClientSecretResponse](#regenerateapiclientsecretresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/{app_id}/api_config/_generate_client_secret *}]}
 <!-- end options -->

#### GetAppKey

> **rpc** GetAppKey([GetAppKeyRequest](#getappkeyrequest))
[GetAppKeyResponse](#getappkeyresponse)




    Key:google.api.http
    Value:{[{GET /projects/{project_id}/apps/{app_id}/keys/{key_id} }]}
 <!-- end options -->

#### ListAppKeys

> **rpc** ListAppKeys([ListAppKeysRequest](#listappkeysrequest))
[ListAppKeysResponse](#listappkeysresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/{app_id}/keys/_search *}]}
 <!-- end options -->

#### AddAppKey

> **rpc** AddAppKey([AddAppKeyRequest](#addappkeyrequest))
[AddAppKeyResponse](#addappkeyresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/apps/{app_id}/keys *}]}
 <!-- end options -->

#### RemoveAppKey

> **rpc** RemoveAppKey([RemoveAppKeyRequest](#removeappkeyrequest))
[RemoveAppKeyResponse](#removeappkeyresponse)




    Key:google.api.http
    Value:{[{DELETE /projects/{project_id}/apps/{app_id}/keys/{key_id} }]}
 <!-- end options -->

#### GetProjectGrantByID

> **rpc** GetProjectGrantByID([GetProjectGrantByIDRequest](#getprojectgrantbyidrequest))
[GetProjectGrantByIDResponse](#getprojectgrantbyidresponse)




    Key:google.api.http
    Value:{[{GET /projects/{project_id}/grants/{grant_id} }]}
 <!-- end options -->

#### ListProjectGrants

> **rpc** ListProjectGrants([ListProjectGrantsRequest](#listprojectgrantsrequest))
[ListProjectGrantsResponse](#listprojectgrantsresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/grants/_search *}]}
 <!-- end options -->

#### AddProjectGrant

> **rpc** AddProjectGrant([AddProjectGrantRequest](#addprojectgrantrequest))
[AddProjectGrantResponse](#addprojectgrantresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/grants *}]}
 <!-- end options -->

#### UpdateProjectGrant

> **rpc** UpdateProjectGrant([UpdateProjectGrantRequest](#updateprojectgrantrequest))
[UpdateProjectGrantResponse](#updateprojectgrantresponse)




    Key:google.api.http
    Value:{[{PUT /projects/{project_id}/grants/{grant_id} *}]}
 <!-- end options -->

#### DeactivateProjectGrant

> **rpc** DeactivateProjectGrant([DeactivateProjectGrantRequest](#deactivateprojectgrantrequest))
[DeactivateProjectGrantResponse](#deactivateprojectgrantresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/grants/{grant_id}/_deactivate *}]}
 <!-- end options -->

#### ReactivateProjectGrant

> **rpc** ReactivateProjectGrant([ReactivateProjectGrantRequest](#reactivateprojectgrantrequest))
[ReactivateProjectGrantResponse](#reactivateprojectgrantresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/grants/{grant_id}/_reactivate *}]}
 <!-- end options -->

#### RemoveProjectGrant

> **rpc** RemoveProjectGrant([RemoveProjectGrantRequest](#removeprojectgrantrequest))
[RemoveProjectGrantResponse](#removeprojectgrantresponse)

RemoveProjectGrant removes project grant and all user grants for this project grant


    Key:google.api.http
    Value:{[{DELETE /projects/{project_id}/grants/{grant_id} }]}
 <!-- end options -->

#### ListProjectGrantMemberRoles

> **rpc** ListProjectGrantMemberRoles([ListProjectGrantMemberRolesRequest](#listprojectgrantmemberrolesrequest))
[ListProjectGrantMemberRolesResponse](#listprojectgrantmemberrolesresponse)




    Key:google.api.http
    Value:{[{POST /projects/grants/members/roles/_search }]}
 <!-- end options -->

#### ListProjectGrantMembers

> **rpc** ListProjectGrantMembers([ListProjectGrantMembersRequest](#listprojectgrantmembersrequest))
[ListProjectGrantMembersResponse](#listprojectgrantmembersresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/grants/{grant_id}/members/_search *}]}
 <!-- end options -->

#### AddProjectGrantMember

> **rpc** AddProjectGrantMember([AddProjectGrantMemberRequest](#addprojectgrantmemberrequest))
[AddProjectGrantMemberResponse](#addprojectgrantmemberresponse)




    Key:google.api.http
    Value:{[{POST /projects/{project_id}/grants/{grant_id}/members *}]}
 <!-- end options -->

#### UpdateProjectGrantMember

> **rpc** UpdateProjectGrantMember([UpdateProjectGrantMemberRequest](#updateprojectgrantmemberrequest))
[UpdateProjectGrantMemberResponse](#updateprojectgrantmemberresponse)




    Key:google.api.http
    Value:{[{PUT /projects/{project_id}/grants/{grant_id}/members/{user_id} *}]}
 <!-- end options -->

#### RemoveProjectGrantMember

> **rpc** RemoveProjectGrantMember([RemoveProjectGrantMemberRequest](#removeprojectgrantmemberrequest))
[RemoveProjectGrantMemberResponse](#removeprojectgrantmemberresponse)




    Key:google.api.http
    Value:{[{DELETE /projects/{project_id}/grants/{grant_id}/members/{user_id} }]}
 <!-- end options -->

#### GetUserGrantByID

> **rpc** GetUserGrantByID([GetUserGrantByIDRequest](#getusergrantbyidrequest))
[GetUserGrantByIDResponse](#getusergrantbyidresponse)




    Key:google.api.http
    Value:{[{GET /users/{user_id}/grants/{grant_id} }]}
 <!-- end options -->

#### ListUserGrants

> **rpc** ListUserGrants([ListUserGrantRequest](#listusergrantrequest))
[ListUserGrantResponse](#listusergrantresponse)




    Key:google.api.http
    Value:{[{POST /users/grants/_search *}]}
 <!-- end options -->

#### AddUserGrant

> **rpc** AddUserGrant([AddUserGrantRequest](#addusergrantrequest))
[AddUserGrantResponse](#addusergrantresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/grants *}]}
 <!-- end options -->

#### UpdateUserGrant

> **rpc** UpdateUserGrant([UpdateUserGrantRequest](#updateusergrantrequest))
[UpdateUserGrantResponse](#updateusergrantresponse)




    Key:google.api.http
    Value:{[{PUT /users/{user_id}/grants/{grant_id} *}]}
 <!-- end options -->

#### DeactivateUserGrant

> **rpc** DeactivateUserGrant([DeactivateUserGrantRequest](#deactivateusergrantrequest))
[DeactivateUserGrantResponse](#deactivateusergrantresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/grants/{grant_id}/_deactivate *}]}
 <!-- end options -->

#### ReactivateUserGrant

> **rpc** ReactivateUserGrant([ReactivateUserGrantRequest](#reactivateusergrantrequest))
[ReactivateUserGrantResponse](#reactivateusergrantresponse)




    Key:google.api.http
    Value:{[{POST /users/{user_id}/grants/{grant_id}/_reactivate *}]}
 <!-- end options -->

#### RemoveUserGrant

> **rpc** RemoveUserGrant([RemoveUserGrantRequest](#removeusergrantrequest))
[RemoveUserGrantResponse](#removeusergrantresponse)




    Key:google.api.http
    Value:{[{DELETE /users/{user_id}/grants/{grant_id} }]}
 <!-- end options -->

#### BulkRemoveUserGrant

> **rpc** BulkRemoveUserGrant([BulkRemoveUserGrantRequest](#bulkremoveusergrantrequest))
[BulkRemoveUserGrantResponse](#bulkremoveusergrantresponse)

remove a list of user grants in one request


    Key:google.api.http
    Value:{[{DELETE /user_grants/_bulk *}]}
 <!-- end options -->

#### GetFeatures

> **rpc** GetFeatures([GetFeaturesRequest](#getfeaturesrequest))
[GetFeaturesResponse](#getfeaturesresponse)




    Key:google.api.http
    Value:{[{GET /features }]}
 <!-- end options -->

#### GetOrgIAMPolicy

> **rpc** GetOrgIAMPolicy([GetOrgIAMPolicyRequest](#getorgiampolicyrequest))
[GetOrgIAMPolicyResponse](#getorgiampolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/orgiam }]}
 <!-- end options -->

#### GetLoginPolicy

> **rpc** GetLoginPolicy([GetLoginPolicyRequest](#getloginpolicyrequest))
[GetLoginPolicyResponse](#getloginpolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/login }]}
 <!-- end options -->

#### GetDefaultLoginPolicy

> **rpc** GetDefaultLoginPolicy([GetDefaultLoginPolicyRequest](#getdefaultloginpolicyrequest))
[GetDefaultLoginPolicyResponse](#getdefaultloginpolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/default/login }]}
 <!-- end options -->

#### AddCustomLoginPolicy

> **rpc** AddCustomLoginPolicy([AddCustomLoginPolicyRequest](#addcustomloginpolicyrequest))
[AddCustomLoginPolicyResponse](#addcustomloginpolicyresponse)




    Key:google.api.http
    Value:{[{POST /policies/login *}]}
 <!-- end options -->

#### UpdateCustomLoginPolicy

> **rpc** UpdateCustomLoginPolicy([UpdateCustomLoginPolicyRequest](#updatecustomloginpolicyrequest))
[UpdateCustomLoginPolicyResponse](#updatecustomloginpolicyresponse)




    Key:google.api.http
    Value:{[{PUT /policies/login *}]}
 <!-- end options -->

#### ResetLoginPolicyToDefault

> **rpc** ResetLoginPolicyToDefault([ResetLoginPolicyToDefaultRequest](#resetloginpolicytodefaultrequest))
[ResetLoginPolicyToDefaultResponse](#resetloginpolicytodefaultresponse)




    Key:google.api.http
    Value:{[{DELETE /policies/login }]}
 <!-- end options -->

#### ListLoginPolicyIDPs

> **rpc** ListLoginPolicyIDPs([ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest))
[ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)




    Key:google.api.http
    Value:{[{POST /policies/login/idps/_search *}]}
 <!-- end options -->

#### AddIDPToLoginPolicy

> **rpc** AddIDPToLoginPolicy([AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest))
[AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)




    Key:google.api.http
    Value:{[{POST /policies/login/idps *}]}
 <!-- end options -->

#### RemoveIDPFromLoginPolicy

> **rpc** RemoveIDPFromLoginPolicy([RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest))
[RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)




    Key:google.api.http
    Value:{[{DELETE /policies/login/idps/{idp_id} }]}
 <!-- end options -->

#### ListLoginPolicySecondFactors

> **rpc** ListLoginPolicySecondFactors([ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest))
[ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)




    Key:google.api.http
    Value:{[{POST /policies/login/second_factors/_search }]}
 <!-- end options -->

#### AddSecondFactorToLoginPolicy

> **rpc** AddSecondFactorToLoginPolicy([AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest))
[AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)




    Key:google.api.http
    Value:{[{POST /policies/login/second_factors *}]}
 <!-- end options -->

#### RemoveSecondFactorFromLoginPolicy

> **rpc** RemoveSecondFactorFromLoginPolicy([RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest))
[RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)




    Key:google.api.http
    Value:{[{DELETE /policies/login/second_factors/{type} }]}
 <!-- end options -->

#### ListLoginPolicyMultiFactors

> **rpc** ListLoginPolicyMultiFactors([ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest))
[ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)




    Key:google.api.http
    Value:{[{POST /policies/login/auth_factors/_search }]}
 <!-- end options -->

#### AddMultiFactorToLoginPolicy

> **rpc** AddMultiFactorToLoginPolicy([AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest))
[AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)




    Key:google.api.http
    Value:{[{POST /policies/login/multi_factors *}]}
 <!-- end options -->

#### RemoveMultiFactorFromLoginPolicy

> **rpc** RemoveMultiFactorFromLoginPolicy([RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest))
[RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)




    Key:google.api.http
    Value:{[{DELETE /policies/login/multi_factors/{type} }]}
 <!-- end options -->

#### GetPasswordComplexityPolicy

> **rpc** GetPasswordComplexityPolicy([GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest))
[GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/password/complexity }]}
 <!-- end options -->

#### GetDefaultPasswordComplexityPolicy

> **rpc** GetDefaultPasswordComplexityPolicy([GetDefaultPasswordComplexityPolicyRequest](#getdefaultpasswordcomplexitypolicyrequest))
[GetDefaultPasswordComplexityPolicyResponse](#getdefaultpasswordcomplexitypolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/default/password/complexity }]}
 <!-- end options -->

#### AddCustomPasswordComplexityPolicy

> **rpc** AddCustomPasswordComplexityPolicy([AddCustomPasswordComplexityPolicyRequest](#addcustompasswordcomplexitypolicyrequest))
[AddCustomPasswordComplexityPolicyResponse](#addcustompasswordcomplexitypolicyresponse)




    Key:google.api.http
    Value:{[{POST /policies/password/complexity *}]}
 <!-- end options -->

#### UpdateCustomPasswordComplexityPolicy

> **rpc** UpdateCustomPasswordComplexityPolicy([UpdateCustomPasswordComplexityPolicyRequest](#updatecustompasswordcomplexitypolicyrequest))
[UpdateCustomPasswordComplexityPolicyResponse](#updatecustompasswordcomplexitypolicyresponse)




    Key:google.api.http
    Value:{[{PUT /policies/password/complexity *}]}
 <!-- end options -->

#### ResetPasswordComplexityPolicyToDefault

> **rpc** ResetPasswordComplexityPolicyToDefault([ResetPasswordComplexityPolicyToDefaultRequest](#resetpasswordcomplexitypolicytodefaultrequest))
[ResetPasswordComplexityPolicyToDefaultResponse](#resetpasswordcomplexitypolicytodefaultresponse)




    Key:google.api.http
    Value:{[{DELETE /policies/password/complexity }]}
 <!-- end options -->

#### GetPasswordAgePolicy

> **rpc** GetPasswordAgePolicy([GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest))
[GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/password/age }]}
 <!-- end options -->

#### GetDefaultPasswordAgePolicy

> **rpc** GetDefaultPasswordAgePolicy([GetDefaultPasswordAgePolicyRequest](#getdefaultpasswordagepolicyrequest))
[GetDefaultPasswordAgePolicyResponse](#getdefaultpasswordagepolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/default/password/age }]}
 <!-- end options -->

#### AddCustomPasswordAgePolicy

> **rpc** AddCustomPasswordAgePolicy([AddCustomPasswordAgePolicyRequest](#addcustompasswordagepolicyrequest))
[AddCustomPasswordAgePolicyResponse](#addcustompasswordagepolicyresponse)




    Key:google.api.http
    Value:{[{POST /policies/password/age *}]}
 <!-- end options -->

#### UpdateCustomPasswordAgePolicy

> **rpc** UpdateCustomPasswordAgePolicy([UpdateCustomPasswordAgePolicyRequest](#updatecustompasswordagepolicyrequest))
[UpdateCustomPasswordAgePolicyResponse](#updatecustompasswordagepolicyresponse)




    Key:google.api.http
    Value:{[{PUT /policies/password/age *}]}
 <!-- end options -->

#### ResetPasswordAgePolicyToDefault

> **rpc** ResetPasswordAgePolicyToDefault([ResetPasswordAgePolicyToDefaultRequest](#resetpasswordagepolicytodefaultrequest))
[ResetPasswordAgePolicyToDefaultResponse](#resetpasswordagepolicytodefaultresponse)




    Key:google.api.http
    Value:{[{DELETE /policies/password/age }]}
 <!-- end options -->

#### GetPasswordLockoutPolicy

> **rpc** GetPasswordLockoutPolicy([GetPasswordLockoutPolicyRequest](#getpasswordlockoutpolicyrequest))
[GetPasswordLockoutPolicyResponse](#getpasswordlockoutpolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/password/lockout }]}
 <!-- end options -->

#### GetDefaultPasswordLockoutPolicy

> **rpc** GetDefaultPasswordLockoutPolicy([GetDefaultPasswordLockoutPolicyRequest](#getdefaultpasswordlockoutpolicyrequest))
[GetDefaultPasswordLockoutPolicyResponse](#getdefaultpasswordlockoutpolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/default/password/lockout }]}
 <!-- end options -->

#### AddCustomPasswordLockoutPolicy

> **rpc** AddCustomPasswordLockoutPolicy([AddCustomPasswordLockoutPolicyRequest](#addcustompasswordlockoutpolicyrequest))
[AddCustomPasswordLockoutPolicyResponse](#addcustompasswordlockoutpolicyresponse)




    Key:google.api.http
    Value:{[{POST /policies/password/lockout *}]}
 <!-- end options -->

#### UpdateCustomPasswordLockoutPolicy

> **rpc** UpdateCustomPasswordLockoutPolicy([UpdateCustomPasswordLockoutPolicyRequest](#updatecustompasswordlockoutpolicyrequest))
[UpdateCustomPasswordLockoutPolicyResponse](#updatecustompasswordlockoutpolicyresponse)




    Key:google.api.http
    Value:{[{PUT /policies/password/lockout *}]}
 <!-- end options -->

#### ResetPasswordLockoutPolicyToDefault

> **rpc** ResetPasswordLockoutPolicyToDefault([ResetPasswordLockoutPolicyToDefaultRequest](#resetpasswordlockoutpolicytodefaultrequest))
[ResetPasswordLockoutPolicyToDefaultResponse](#resetpasswordlockoutpolicytodefaultresponse)




    Key:google.api.http
    Value:{[{DELETE /policies/password/lockout }]}
 <!-- end options -->

#### GetLabelPolicy

> **rpc** GetLabelPolicy([GetLabelPolicyRequest](#getlabelpolicyrequest))
[GetLabelPolicyResponse](#getlabelpolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/label }]}
 <!-- end options -->

#### GetDefaultLabelPolicy

> **rpc** GetDefaultLabelPolicy([GetDefaultLabelPolicyRequest](#getdefaultlabelpolicyrequest))
[GetDefaultLabelPolicyResponse](#getdefaultlabelpolicyresponse)




    Key:google.api.http
    Value:{[{GET /policies/default/label }]}
 <!-- end options -->

#### AddCustomLabelPolicy

> **rpc** AddCustomLabelPolicy([AddCustomLabelPolicyRequest](#addcustomlabelpolicyrequest))
[AddCustomLabelPolicyResponse](#addcustomlabelpolicyresponse)




    Key:google.api.http
    Value:{[{POST /policies/label *}]}
 <!-- end options -->

#### UpdateCustomLabelPolicy

> **rpc** UpdateCustomLabelPolicy([UpdateCustomLabelPolicyRequest](#updatecustomlabelpolicyrequest))
[UpdateCustomLabelPolicyResponse](#updatecustomlabelpolicyresponse)




    Key:google.api.http
    Value:{[{PUT /policies/label *}]}
 <!-- end options -->

#### ResetLabelPolicyToDefault

> **rpc** ResetLabelPolicyToDefault([ResetLabelPolicyToDefaultRequest](#resetlabelpolicytodefaultrequest))
[ResetLabelPolicyToDefaultResponse](#resetlabelpolicytodefaultresponse)




    Key:google.api.http
    Value:{[{DELETE /policies/label }]}
 <!-- end options -->

#### GetOrgIDPByID

> **rpc** GetOrgIDPByID([GetOrgIDPByIDRequest](#getorgidpbyidrequest))
[GetOrgIDPByIDResponse](#getorgidpbyidresponse)




    Key:google.api.http
    Value:{[{GET /idps/{id} }]}
 <!-- end options -->

#### ListOrgIDPs

> **rpc** ListOrgIDPs([ListOrgIDPsRequest](#listorgidpsrequest))
[ListOrgIDPsResponse](#listorgidpsresponse)




    Key:google.api.http
    Value:{[{POST /idps/_search *}]}
 <!-- end options -->

#### AddOrgOIDCIDP

> **rpc** AddOrgOIDCIDP([AddOrgOIDCIDPRequest](#addorgoidcidprequest))
[AddOrgOIDCIDPResponse](#addorgoidcidpresponse)




    Key:google.api.http
    Value:{[{POST /idps/oidc *}]}
 <!-- end options -->

#### DeactivateOrgIDP

> **rpc** DeactivateOrgIDP([DeactivateOrgIDPRequest](#deactivateorgidprequest))
[DeactivateOrgIDPResponse](#deactivateorgidpresponse)




    Key:google.api.http
    Value:{[{POST /idps/{idp_id}/_deactivate *}]}
 <!-- end options -->

#### ReactivateOrgIDP

> **rpc** ReactivateOrgIDP([ReactivateOrgIDPRequest](#reactivateorgidprequest))
[ReactivateOrgIDPResponse](#reactivateorgidpresponse)




    Key:google.api.http
    Value:{[{POST /idps/{idp_id}/_reactivate *}]}
 <!-- end options -->

#### RemoveOrgIDP

> **rpc** RemoveOrgIDP([RemoveOrgIDPRequest](#removeorgidprequest))
[RemoveOrgIDPResponse](#removeorgidpresponse)




    Key:google.api.http
    Value:{[{DELETE /idps/{idp_id} }]}
 <!-- end options -->

#### UpdateOrgIDP

> **rpc** UpdateOrgIDP([UpdateOrgIDPRequest](#updateorgidprequest))
[UpdateOrgIDPResponse](#updateorgidpresponse)




    Key:google.api.http
    Value:{[{PUT /idps/{idp_id} *}]}
 <!-- end options -->

#### UpdateOrgIDPOIDCConfig

> **rpc** UpdateOrgIDPOIDCConfig([UpdateOrgIDPOIDCConfigRequest](#updateorgidpoidcconfigrequest))
[UpdateOrgIDPOIDCConfigResponse](#updateorgidpoidcconfigresponse)




    Key:google.api.http
    Value:{[{PUT /idps/{idp_id}/oidc_config *}]}
 <!-- end options -->

 <!-- end methods -->
 <!-- end services -->

### Messages


#### AddAPIAppRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| name | [ string](#string) | - |
| auth_method_type | [ zitadel.app.v1.APIAuthMethodType](#zitadelappv1apiauthmethodtype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddAPIAppResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| app_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| client_id | [ string](#string) | - |
| client_secret | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddAppKeyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
| type | [ zitadel.authn.v1.KeyType](#zitadelauthnv1keytype) | - |
| expiration_date | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddAppKeyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| key_details | [ bytes](#bytes) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomLabelPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| primary_color | [ string](#string) | - |
| secondary_color | [ string](#string) | - |
| hide_login_name_suffix | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomLabelPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| allow_username_password | [ bool](#bool) | - |
| allow_register | [ bool](#bool) | - |
| allow_external_idp | [ bool](#bool) | - |
| force_mfa | [ bool](#bool) | - |
| passwordless_type | [ zitadel.policy.v1.PasswordlessType](#zitadelpolicyv1passwordlesstype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomPasswordAgePolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| max_age_days | [ uint32](#uint32) | - |
| expire_warn_days | [ uint32](#uint32) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomPasswordAgePolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomPasswordComplexityPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| min_length | [ uint64](#uint64) | - |
| has_uppercase | [ bool](#bool) | - |
| has_lowercase | [ bool](#bool) | - |
| has_number | [ bool](#bool) | - |
| has_symbol | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomPasswordComplexityPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomPasswordLockoutPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| max_attempts | [ uint32](#uint32) | - |
| show_lockout_failure | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddCustomPasswordLockoutPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddHumanUserRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name | [ string](#string) | - |
| profile | [ AddHumanUserRequest.Profile](#addhumanuserrequestprofile) | - |
| email | [ AddHumanUserRequest.Email](#addhumanuserrequestemail) | - |
| phone | [ AddHumanUserRequest.Phone](#addhumanuserrequestphone) | - |
| initial_password | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddHumanUserRequest.Email



| Field | Type | Description |
| ----- | ---- | ----------- |
| email | [ string](#string) | TODO: check if no value is allowed |
| is_email_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddHumanUserRequest.Phone



| Field | Type | Description |
| ----- | ---- | ----------- |
| phone | [ string](#string) | has to be a global number |
| is_phone_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddHumanUserRequest.Profile



| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name | [ string](#string) | - |
| last_name | [ string](#string) | - |
| nick_name | [ string](#string) | - |
| display_name | [ string](#string) | - |
| preferred_language | [ string](#string) | - |
| gender | [ zitadel.user.v1.Gender](#zitadeluserv1gender) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddHumanUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddIDPToLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
| ownerType | [ zitadel.idp.v1.IDPOwnerType](#zitadelidpv1idpownertype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddIDPToLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMachineKeyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| type | [ zitadel.authn.v1.KeyType](#zitadelauthnv1keytype) | - |
| expiration_date | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMachineKeyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| key_id | [ string](#string) | - |
| key_details | [ bytes](#bytes) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMachineUserRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name | [ string](#string) | - |
| name | [ string](#string) | - |
| description | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMachineUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMultiFactorToLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ zitadel.policy.v1.MultiFactorType](#zitadelpolicyv1multifactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddMultiFactorToLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOIDCAppRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| name | [ string](#string) | - |
| redirect_uris | [repeated string](#string) | - |
| response_types | [repeated zitadel.app.v1.OIDCResponseType](#zitadelappv1oidcresponsetype) | - |
| grant_types | [repeated zitadel.app.v1.OIDCGrantType](#zitadelappv1oidcgranttype) | - |
| app_type | [ zitadel.app.v1.OIDCAppType](#zitadelappv1oidcapptype) | - |
| auth_method_type | [ zitadel.app.v1.OIDCAuthMethodType](#zitadelappv1oidcauthmethodtype) | - |
| post_logout_redirect_uris | [repeated string](#string) | - |
| version | [ zitadel.app.v1.OIDCVersion](#zitadelappv1oidcversion) | - |
| dev_mode | [ bool](#bool) | - |
| access_token_type | [ zitadel.app.v1.OIDCTokenType](#zitadelappv1oidctokentype) | - |
| access_token_role_assertion | [ bool](#bool) | - |
| id_token_role_assertion | [ bool](#bool) | - |
| id_token_userinfo_assertion | [ bool](#bool) | - |
| clock_skew | [ google.protobuf.Duration](#googleprotobufduration) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOIDCAppResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| app_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| client_id | [ string](#string) | - |
| client_secret | [ string](#string) | - |
| none_compliant | [ bool](#bool) | - |
| compliance_problems | [repeated zitadel.v1.LocalizedMessage](#zitadelv1localizedmessage) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOrgDomainRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| domain | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOrgDomainResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOrgMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOrgMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOrgOIDCIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| styling_type | [ zitadel.idp.v1.IDPStylingType](#zitadelidpv1idpstylingtype) | - |
| client_id | [ string](#string) | - |
| client_secret | [ string](#string) | - |
| issuer | [ string](#string) | - |
| scopes | [repeated string](#string) | - |
| display_name_mapping | [ zitadel.idp.v1.OIDCMappingField](#zitadelidpv1oidcmappingfield) | - |
| username_mapping | [ zitadel.idp.v1.OIDCMappingField](#zitadelidpv1oidcmappingfield) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOrgOIDCIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOrgRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddOrgResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectGrantMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectGrantMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| granted_org_id | [ string](#string) | - |
| role_keys | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| grant_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| project_role_assertion | [ bool](#bool) | - |
| project_role_check | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectRoleRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| role_key | [ string](#string) | - |
| display_name | [ string](#string) | - |
| group | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddProjectRoleResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddSecondFactorToLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ zitadel.policy.v1.SecondFactorType](#zitadelpolicyv1secondfactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddSecondFactorToLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddUserGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| project_id | [ string](#string) | - |
| project_grant_id | [ string](#string) | - |
| role_keys | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AddUserGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_grant_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### BulkAddProjectRolesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| roles | [repeated BulkAddProjectRolesRequest.Role](#bulkaddprojectrolesrequestrole) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### BulkAddProjectRolesRequest.Role



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ string](#string) | - |
| display_name | [ string](#string) | - |
| group | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### BulkAddProjectRolesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### BulkRemoveUserGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| grant_id | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### BulkRemoveUserGrantResponse


 <!-- end HasFields -->


#### DeactivateAppRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateAppResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateOrgIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateOrgIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateOrgRequest


 <!-- end HasFields -->


#### DeactivateOrgResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateProjectGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateProjectGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateProjectRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateProjectResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateUserGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateUserGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateUserRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DeactivateUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GenerateOrgDomainValidationRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| domain | [ string](#string) | - |
| type | [ zitadel.org.v1.DomainValidationType](#zitadelorgv1domainvalidationtype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GenerateOrgDomainValidationResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| token | [ string](#string) | - |
| url | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetAppByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetAppByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| app | [ zitadel.app.v1.App](#zitadelappv1app) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetAppKeyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
| key_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetAppKeyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ zitadel.authn.v1.Key](#zitadelauthnv1key) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetDefaultLabelPolicyRequest


 <!-- end HasFields -->


#### GetDefaultLabelPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.LabelPolicy](#zitadelpolicyv1labelpolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetDefaultLoginPolicyRequest


 <!-- end HasFields -->


#### GetDefaultLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.LoginPolicy](#zitadelpolicyv1loginpolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetDefaultPasswordAgePolicyRequest


 <!-- end HasFields -->


#### GetDefaultPasswordAgePolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordAgePolicy](#zitadelpolicyv1passwordagepolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetDefaultPasswordComplexityPolicyRequest


 <!-- end HasFields -->


#### GetDefaultPasswordComplexityPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordComplexityPolicy](#zitadelpolicyv1passwordcomplexitypolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetDefaultPasswordLockoutPolicyRequest


 <!-- end HasFields -->


#### GetDefaultPasswordLockoutPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordLockoutPolicy](#zitadelpolicyv1passwordlockoutpolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetFeaturesRequest


 <!-- end HasFields -->


#### GetFeaturesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| features | [ zitadel.features.v1.Features](#zitadelfeaturesv1features) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetGrantedProjectByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetGrantedProjectByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| granted_project | [ zitadel.project.v1.GrantedProject](#zitadelprojectv1grantedproject) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetHumanEmailRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetHumanEmailResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| email | [ zitadel.user.v1.Email](#zitadeluserv1email) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetHumanPhoneRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetHumanPhoneResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| phone | [ zitadel.user.v1.Phone](#zitadeluserv1phone) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetHumanProfileRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetHumanProfileResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| profile | [ zitadel.user.v1.Profile](#zitadeluserv1profile) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetIAMRequest


 <!-- end HasFields -->


#### GetIAMResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| global_org_id | [ string](#string) | - |
| iam_project_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetLabelPolicyRequest


 <!-- end HasFields -->


#### GetLabelPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.LabelPolicy](#zitadelpolicyv1labelpolicy) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetLoginPolicyRequest


 <!-- end HasFields -->


#### GetLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.LoginPolicy](#zitadelpolicyv1loginpolicy) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetMachineKeyByIDsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| key_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetMachineKeyByIDsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ zitadel.authn.v1.Key](#zitadelauthnv1key) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetMyOrgRequest


 <!-- end HasFields -->


#### GetMyOrgResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| org | [ zitadel.org.v1.Org](#zitadelorgv1org) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOIDCInformationRequest


 <!-- end HasFields -->


#### GetOIDCInformationResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| issuer | [ string](#string) | - |
| discovery_endpoint | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgByDomainGlobalRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| domain | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgByDomainGlobalResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| org | [ zitadel.org.v1.Org](#zitadelorgv1org) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgIAMPolicyRequest


 <!-- end HasFields -->


#### GetOrgIAMPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.OrgIAMPolicy](#zitadelpolicyv1orgiampolicy) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgIDPByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetOrgIDPByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp | [ zitadel.idp.v1.IDP](#zitadelidpv1idp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetPasswordAgePolicyRequest


 <!-- end HasFields -->


#### GetPasswordAgePolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordAgePolicy](#zitadelpolicyv1passwordagepolicy) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetPasswordComplexityPolicyRequest


 <!-- end HasFields -->


#### GetPasswordComplexityPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordComplexityPolicy](#zitadelpolicyv1passwordcomplexitypolicy) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetPasswordLockoutPolicyRequest


 <!-- end HasFields -->


#### GetPasswordLockoutPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| policy | [ zitadel.policy.v1.PasswordLockoutPolicy](#zitadelpolicyv1passwordlockoutpolicy) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetProjectByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetProjectByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| project | [ zitadel.project.v1.Project](#zitadelprojectv1project) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetProjectGrantByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetProjectGrantByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_grant | [ zitadel.project.v1.GrantedProject](#zitadelprojectv1grantedproject) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetUserByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetUserByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| user | [ zitadel.user.v1.User](#zitadeluserv1user) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetUserByLoginNameGlobalRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| login_name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetUserByLoginNameGlobalResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| user | [ zitadel.user.v1.User](#zitadeluserv1user) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetUserGrantByIDRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GetUserGrantByIDResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_grant | [ zitadel.user.v1.UserGrant](#zitadeluserv1usergrant) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### HealthzRequest


 <!-- end HasFields -->


#### HealthzResponse


 <!-- end HasFields -->


#### IDPQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_id_query | [ zitadel.idp.v1.IDPIDQuery](#zitadelidpv1idpidquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_name_query | [ zitadel.idp.v1.IDPNameQuery](#zitadelidpv1idpnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.owner_type_query | [ zitadel.idp.v1.IDPOwnerTypeQuery](#zitadelidpv1idpownertypequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ImportHumanUserRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name | [ string](#string) | - |
| profile | [ ImportHumanUserRequest.Profile](#importhumanuserrequestprofile) | - |
| email | [ ImportHumanUserRequest.Email](#importhumanuserrequestemail) | - |
| phone | [ ImportHumanUserRequest.Phone](#importhumanuserrequestphone) | - |
| password | [ string](#string) | - |
| password_change_required | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ImportHumanUserRequest.Email



| Field | Type | Description |
| ----- | ---- | ----------- |
| email | [ string](#string) | TODO: check if no value is allowed |
| is_email_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ImportHumanUserRequest.Phone



| Field | Type | Description |
| ----- | ---- | ----------- |
| phone | [ string](#string) | has to be a global number |
| is_phone_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ImportHumanUserRequest.Profile



| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name | [ string](#string) | - |
| last_name | [ string](#string) | - |
| nick_name | [ string](#string) | - |
| display_name | [ string](#string) | - |
| preferred_language | [ string](#string) | - |
| gender | [ zitadel.user.v1.Gender](#zitadeluserv1gender) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ImportHumanUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IsUserUniqueRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name | [ string](#string) | - |
| email | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### IsUserUniqueResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| is_unique | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListAppChangesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.change.v1.ChangeQuery](#zitadelchangev1changequery) | list limitations and ordering |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListAppChangesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.change.v1.Change](#zitadelchangev1change) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListAppKeysRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| app_id | [ string](#string) | - |
| project_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListAppKeysResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.authn.v1.Key](#zitadelauthnv1key) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListAppsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.app.v1.AppQuery](#zitadelappv1appquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListAppsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.app.v1.App](#zitadelappv1app) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListGrantedProjectRolesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.project.v1.RoleQuery](#zitadelprojectv1rolequery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListGrantedProjectRolesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.project.v1.Role](#zitadelprojectv1role) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListGrantedProjectsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.project.v1.ProjectQuery](#zitadelprojectv1projectquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListGrantedProjectsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.project.v1.GrantedProject](#zitadelprojectv1grantedproject) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListHumanAuthFactorsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListHumanAuthFactorsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated zitadel.user.v1.AuthFactor](#zitadeluserv1authfactor) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListHumanLinkedIDPsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListHumanLinkedIDPsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.idp.v1.IDPUserLink](#zitadelidpv1idpuserlink) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListHumanPasswordlessRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListHumanPasswordlessResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated zitadel.user.v1.WebAuthNToken](#zitadeluserv1webauthntoken) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListLoginPolicyIDPsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListLoginPolicyIDPsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.idp.v1.IDPLoginPolicyLink](#zitadelidpv1idploginpolicylink) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListLoginPolicyMultiFactorsRequest


 <!-- end HasFields -->


#### ListLoginPolicyMultiFactorsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.policy.v1.MultiFactorType](#zitadelpolicyv1multifactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListLoginPolicySecondFactorsRequest


 <!-- end HasFields -->


#### ListLoginPolicySecondFactorsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.policy.v1.SecondFactorType](#zitadelpolicyv1secondfactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMachineKeysRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListMachineKeysResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.authn.v1.Key](#zitadelauthnv1key) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgChangesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.change.v1.ChangeQuery](#zitadelchangev1changequery) | list limitations and ordering |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgChangesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.change.v1.Change](#zitadelchangev1change) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgDomainsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.org.v1.DomainSearchQuery](#zitadelorgv1domainsearchquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgDomainsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.org.v1.Domain](#zitadelorgv1domain) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgIDPsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| sorting_column | [ zitadel.idp.v1.IDPFieldName](#zitadelidpv1idpfieldname) | the field the result is sorted |
| queries | [repeated IDPQuery](#idpquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgIDPsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| sorting_column | [ zitadel.idp.v1.IDPFieldName](#zitadelidpv1idpfieldname) | - |
| result | [repeated zitadel.idp.v1.IDP](#zitadelidpv1idp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgMemberRolesRequest


 <!-- end HasFields -->


#### ListOrgMemberRolesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| result | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgMembersRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.member.v1.SearchQuery](#zitadelmemberv1searchquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListOrgMembersResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | list limitations and ordering |
| result | [repeated zitadel.member.v1.Member](#zitadelmemberv1member) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectChangesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.change.v1.ChangeQuery](#zitadelchangev1changequery) | list limitations and ordering |
| project_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectChangesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.change.v1.Change](#zitadelchangev1change) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectGrantMemberRolesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | - |
| result | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectGrantMemberRolesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectGrantMembersRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.member.v1.SearchQuery](#zitadelmemberv1searchquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectGrantMembersResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.member.v1.Member](#zitadelmemberv1member) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectGrantsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.project.v1.ProjectGrantQuery](#zitadelprojectv1projectgrantquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectGrantsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.project.v1.GrantedProject](#zitadelprojectv1grantedproject) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectMemberRolesRequest


 <!-- end HasFields -->


#### ListProjectMemberRolesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectMembersRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.member.v1.SearchQuery](#zitadelmemberv1searchquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectMembersResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.member.v1.Member](#zitadelmemberv1member) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectRolesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.project.v1.RoleQuery](#zitadelprojectv1rolequery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectRolesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.project.v1.Role](#zitadelprojectv1role) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.project.v1.ProjectQuery](#zitadelprojectv1projectquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListProjectsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.project.v1.Project](#zitadelprojectv1project) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListUserChangesRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.change.v1.ChangeQuery](#zitadelchangev1changequery) | list limitations and ordering |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListUserChangesResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.change.v1.Change](#zitadelchangev1change) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListUserGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| queries | [repeated zitadel.user.v1.UserGrantQuery](#zitadeluserv1usergrantquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListUserGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.user.v1.UserGrant](#zitadeluserv1usergrant) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListUserMembershipsRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | list limitations and ordering |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | the field the result is sorted |
| queries | [repeated zitadel.user.v1.MembershipQuery](#zitadeluserv1membershipquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListUserMembershipsResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| result | [repeated zitadel.user.v1.Membership](#zitadeluserv1membership) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListUsersRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| query | [ zitadel.v1.ListQuery](#zitadelv1listquery) | list limitations and ordering |
| sorting_column | [ zitadel.user.v1.UserFieldName](#zitadeluserv1userfieldname) | the field the result is sorted |
| queries | [repeated zitadel.user.v1.SearchQuery](#zitadeluserv1searchquery) | criterias the client is looking for |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListUsersResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ListDetails](#zitadelv1listdetails) | - |
| sorting_column | [ zitadel.user.v1.UserFieldName](#zitadeluserv1userfieldname) | - |
| result | [repeated zitadel.user.v1.User](#zitadeluserv1user) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### LockUserRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### LockUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateAppRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateAppResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateOrgIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateOrgIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateOrgRequest


 <!-- end HasFields -->


#### ReactivateOrgResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateProjectGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateProjectGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateProjectRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateProjectResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateUserGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateUserGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateUserRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ReactivateUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RegenerateAPIClientSecretRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RegenerateAPIClientSecretResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| client_secret | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RegenerateOIDCClientSecretRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RegenerateOIDCClientSecretResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| client_secret | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveAppKeyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
| key_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveAppKeyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveAppRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveAppResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanAuthFactorOTPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanAuthFactorOTPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanAuthFactorU2FRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| token_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanAuthFactorU2FResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanLinkedIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| idp_id | [ string](#string) | - |
| linked_user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanLinkedIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanPasswordlessRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| token_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanPasswordlessResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanPhoneRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveHumanPhoneResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveIDPFromLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveIDPFromLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMachineKeyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| key_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMachineKeyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMultiFactorFromLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ zitadel.policy.v1.MultiFactorType](#zitadelpolicyv1multifactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveMultiFactorFromLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveOrgDomainRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| domain | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveOrgDomainResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveOrgIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveOrgIDPResponse


 <!-- end HasFields -->


#### RemoveOrgMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveOrgMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectGrantMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectGrantMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectRoleRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| role_key | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveProjectRoleResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveSecondFactorFromLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ zitadel.policy.v1.SecondFactorType](#zitadelpolicyv1secondfactortype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveSecondFactorFromLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveUserGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveUserGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveUserRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RemoveUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResendHumanEmailVerificationRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResendHumanEmailVerificationResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResendHumanInitializationRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| email | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResendHumanInitializationResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResendHumanPhoneVerificationRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResendHumanPhoneVerificationResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetLabelPolicyToDefaultRequest


 <!-- end HasFields -->


#### ResetLabelPolicyToDefaultResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetLoginPolicyToDefaultRequest


 <!-- end HasFields -->


#### ResetLoginPolicyToDefaultResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetPasswordAgePolicyToDefaultRequest


 <!-- end HasFields -->


#### ResetPasswordAgePolicyToDefaultResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetPasswordComplexityPolicyToDefaultRequest


 <!-- end HasFields -->


#### ResetPasswordComplexityPolicyToDefaultResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ResetPasswordLockoutPolicyToDefaultRequest


 <!-- end HasFields -->


#### ResetPasswordLockoutPolicyToDefaultResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SendHumanResetPasswordNotificationRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| type | [ SendHumanResetPasswordNotificationRequest.Type](#sendhumanresetpasswordnotificationrequesttype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SendHumanResetPasswordNotificationResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetHumanInitialPasswordRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| password | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetHumanInitialPasswordResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetPrimaryOrgDomainRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| domain | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SetPrimaryOrgDomainResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UnlockUserRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UnlockUserResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateAPIAppConfigRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
| auth_method_type | [ zitadel.app.v1.APIAuthMethodType](#zitadelappv1apiauthmethodtype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateAPIAppConfigResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateAppRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
| name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateAppResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomLabelPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| primary_color | [ string](#string) | - |
| secondary_color | [ string](#string) | - |
| hide_login_name_suffix | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomLabelPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomLoginPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| allow_username_password | [ bool](#bool) | - |
| allow_register | [ bool](#bool) | - |
| allow_external_idp | [ bool](#bool) | - |
| force_mfa | [ bool](#bool) | - |
| passwordless_type | [ zitadel.policy.v1.PasswordlessType](#zitadelpolicyv1passwordlesstype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomLoginPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomPasswordAgePolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| max_age_days | [ uint32](#uint32) | - |
| expire_warn_days | [ uint32](#uint32) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomPasswordAgePolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomPasswordComplexityPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| min_length | [ uint64](#uint64) | - |
| has_uppercase | [ bool](#bool) | - |
| has_lowercase | [ bool](#bool) | - |
| has_number | [ bool](#bool) | - |
| has_symbol | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomPasswordComplexityPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomPasswordLockoutPolicyRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| max_attempts | [ uint32](#uint32) | - |
| show_lockout_failure | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateCustomPasswordLockoutPolicyResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateHumanEmailRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| email | [ string](#string) | - |
| is_email_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateHumanEmailResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateHumanPhoneRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| phone | [ string](#string) | - |
| is_phone_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateHumanPhoneResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateHumanProfileRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| first_name | [ string](#string) | - |
| last_name | [ string](#string) | - |
| nick_name | [ string](#string) | - |
| display_name | [ string](#string) | - |
| preferred_language | [ string](#string) | - |
| gender | [ zitadel.user.v1.Gender](#zitadeluserv1gender) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateHumanProfileResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateMachineRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| description | [ string](#string) | - |
| name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateMachineResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOIDCAppConfigRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| app_id | [ string](#string) | - |
| redirect_uris | [repeated string](#string) | - |
| response_types | [repeated zitadel.app.v1.OIDCResponseType](#zitadelappv1oidcresponsetype) | - |
| grant_types | [repeated zitadel.app.v1.OIDCGrantType](#zitadelappv1oidcgranttype) | - |
| app_type | [ zitadel.app.v1.OIDCAppType](#zitadelappv1oidcapptype) | - |
| auth_method_type | [ zitadel.app.v1.OIDCAuthMethodType](#zitadelappv1oidcauthmethodtype) | - |
| post_logout_redirect_uris | [repeated string](#string) | - |
| dev_mode | [ bool](#bool) | - |
| access_token_type | [ zitadel.app.v1.OIDCTokenType](#zitadelappv1oidctokentype) | - |
| access_token_role_assertion | [ bool](#bool) | - |
| id_token_role_assertion | [ bool](#bool) | - |
| id_token_userinfo_assertion | [ bool](#bool) | - |
| clock_skew | [ google.protobuf.Duration](#googleprotobufduration) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOIDCAppConfigResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOrgIDPOIDCConfigRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
| client_id | [ string](#string) | - |
| client_secret | [ string](#string) | - |
| issuer | [ string](#string) | - |
| scopes | [repeated string](#string) | - |
| display_name_mapping | [ zitadel.idp.v1.OIDCMappingField](#zitadelidpv1oidcmappingfield) | - |
| username_mapping | [ zitadel.idp.v1.OIDCMappingField](#zitadelidpv1oidcmappingfield) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOrgIDPOIDCConfigResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOrgIDPRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id | [ string](#string) | - |
| name | [ string](#string) | - |
| styling_type | [ zitadel.idp.v1.IDPStylingType](#zitadelidpv1idpstylingtype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOrgIDPResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOrgMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateOrgMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectGrantMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectGrantMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
| role_keys | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectMemberRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| user_id | [ string](#string) | - |
| roles | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectMemberResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| name | [ string](#string) | - |
| project_role_assertion | [ bool](#bool) | - |
| project_role_check | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectRoleRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
| role_key | [ string](#string) | - |
| display_name | [ string](#string) | - |
| group | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateProjectRoleResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateUserGrantRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| grant_id | [ string](#string) | - |
| role_keys | [repeated string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateUserGrantResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateUserNameRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| user_name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UpdateUserNameResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ValidateOrgDomainRequest



| Field | Type | Description |
| ----- | ---- | ----------- |
| domain | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ValidateOrgDomainResponse



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### SendHumanResetPasswordNotificationRequest.Type {#sendhumanresetpasswordnotificationrequesttype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_EMAIL | 0 | - |
| TYPE_SMS | 1 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### EmailQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| email | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### FirstNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### LastNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| last_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Member



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| roles | [repeated string](#string) | - |
| preferred_login_name | [ string](#string) | - |
| email | [ string](#string) | - |
| first_name | [ string](#string) | - |
| last_name | [ string](#string) | - |
| display_name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SearchQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.first_name_query | [ FirstNameQuery](#firstnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.last_name_query | [ LastNameQuery](#lastnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.email_query | [ EmailQuery](#emailquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.user_id_query | [ UserIDQuery](#useridquery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserIDQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->

 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### ErrorDetail



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| message | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### LocalizedMessage



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ string](#string) | - |
| localized_message | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->

 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### ListDetails



| Field | Type | Description |
| ----- | ---- | ----------- |
| total_result | [ uint64](#uint64) | - |
| processed_sequence | [ uint64](#uint64) | - |
| view_timestamp | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ListQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| offset | [ uint64](#uint64) | - |
| limit | [ uint32](#uint32) | - |
| asc | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ObjectDetails



| Field | Type | Description |
| ----- | ---- | ----------- |
| sequence | [ uint64](#uint64) | sequence represents the order of events. It's always upcounting

on read: the sequence of the last event reduced by the projection

on manipulation: the timestamp of the event(s) added by the manipulation |
| creation_date | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | creation_date is the timestamp where the first operation on the object was made

on read: the timestamp of the first event of the object

on create: the timestamp of the event(s) added by the manipulation |
| change_date | [ google.protobuf.Timestamp](#googleprotobuftimestamp) | change_date is the timestamp when the object was changed

on read: the timestamp of the last event reduced by the projection

on manipulation: the |
| resource_owner | [ string](#string) | resource_owner is the organisation an object belongs to |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### TextQueryMethod {#textquerymethod}


| Name | Number | Description |
| ---- | ------ | ----------- |
| TEXT_QUERY_METHOD_EQUALS | 0 | - |
| TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE | 1 | - |
| TEXT_QUERY_METHOD_STARTS_WITH | 2 | - |
| TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE | 3 | - |
| TEXT_QUERY_METHOD_CONTAINS | 4 | - |
| TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE | 5 | - |
| TEXT_QUERY_METHOD_ENDS_WITH | 6 | - |
| TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE | 7 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### AuthOption



| Field | Type | Description |
| ----- | ---- | ----------- |
| permission | [ string](#string) | - |
| check_field_name | [ string](#string) | - |
| feature | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->

 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### Domain



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| domain_name | [ string](#string) | - |
| is_verified | [ bool](#bool) | - |
| is_primary | [ bool](#bool) | - |
| validation_type | [ DomainValidationType](#domainvalidationtype) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DomainNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DomainSearchQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.domain_name_query | [ DomainNameQuery](#domainnamequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Org



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| state | [ OrgState](#orgstate) | - |
| name | [ string](#string) | - |
| primary_domain | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### OrgDomainQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| domain | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### OrgNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### OrgQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.name_query | [ OrgNameQuery](#orgnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.domain_query | [ OrgDomainQuery](#orgdomainquery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### DomainValidationType {#domainvalidationtype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| DOMAIN_VALIDATION_TYPE_UNSPECIFIED | 0 | - |
| DOMAIN_VALIDATION_TYPE_HTTP | 1 | - |
| DOMAIN_VALIDATION_TYPE_DNS | 2 | - |




#### OrgFieldName {#orgfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| ORG_FIELD_NAME_UNSPECIFIED | 0 | - |
| ORG_FIELD_NAME_NAME | 1 | - |




#### OrgState {#orgstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| ORG_STATE_UNSPECIFIED | 0 | - |
| ORG_STATE_ACTIVE | 1 | - |
| ORG_STATE_INACTIVE | 2 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### LabelPolicy



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| primary_color | [ string](#string) | - |
| secondary_color | [ string](#string) | - |
| is_default | [ bool](#bool) | - |
| hide_login_name_suffix | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### LoginPolicy



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| allow_username_password | [ bool](#bool) | - |
| allow_register | [ bool](#bool) | - |
| allow_external_idp | [ bool](#bool) | - |
| force_mfa | [ bool](#bool) | - |
| passwordless_type | [ PasswordlessType](#passwordlesstype) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### OrgIAMPolicy



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| user_login_must_be_domain | [ bool](#bool) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### PasswordAgePolicy



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| max_age_days | [ uint64](#uint64) | - |
| expire_warn_days | [ uint64](#uint64) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### PasswordComplexityPolicy



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| min_length | [ uint64](#uint64) | - |
| has_uppercase | [ bool](#bool) | - |
| has_lowercase | [ bool](#bool) | - |
| has_number | [ bool](#bool) | - |
| has_symbol | [ bool](#bool) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### PasswordLockoutPolicy



| Field | Type | Description |
| ----- | ---- | ----------- |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| max_attempts | [ uint64](#uint64) | - |
| show_lockout_failure | [ bool](#bool) | - |
| is_default | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### MultiFactorType {#multifactortype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| MULTI_FACTOR_TYPE_UNSPECIFIED | 0 | - |
| MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION | 1 | TODO: what does livio think after the weekend? :D |




#### PasswordlessType {#passwordlesstype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PASSWORDLESS_TYPE_NOT_ALLOWED | 0 | - |
| PASSWORDLESS_TYPE_ALLOWED | 1 | PLANNED: PASSWORDLESS_TYPE_WITH_CERT |




#### SecondFactorType {#secondfactortype}


| Name | Number | Description |
| ---- | ------ | ----------- |
| SECOND_FACTOR_TYPE_UNSPECIFIED | 0 | - |
| SECOND_FACTOR_TYPE_OTP | 1 | - |
| SECOND_FACTOR_TYPE_U2F | 2 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### GrantProjectNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GrantRoleKeyQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| role_key | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### GrantedProject



| Field | Type | Description |
| ----- | ---- | ----------- |
| grant_id | [ string](#string) | - |
| granted_org_id | [ string](#string) | - |
| granted_org_name | [ string](#string) | - |
| granted_role_keys | [repeated string](#string) | - |
| state | [ ProjectGrantState](#projectgrantstate) | - |
| project_id | [ string](#string) | - |
| project_name | [ string](#string) | - |
| project_owner_id | [ string](#string) | - |
| project_owner_name | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Project



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| name | [ string](#string) | - |
| state | [ ProjectState](#projectstate) | - |
| project_role_assertion | [ bool](#bool) | describes if roles of user should be added in token |
| project_role_check | [ bool](#bool) | ZITADEL checks if the user has at least one on this project |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ProjectGrantQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_name_query | [ GrantProjectNameQuery](#grantprojectnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.role_key_query | [ GrantRoleKeyQuery](#grantrolekeyquery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ProjectNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### ProjectQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.name_query | [ ProjectNameQuery](#projectnamequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Role



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| display_name | [ string](#string) | - |
| group | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RoleDisplayNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| display_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RoleKeyQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### RoleQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.key_query | [ RoleKeyQuery](#rolekeyquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.display_name_query | [ RoleDisplayNameQuery](#roledisplaynamequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### ProjectGrantState {#projectgrantstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PROJECT_GRANT_STATE_UNSPECIFIED | 0 | - |
| PROJECT_GRANT_STATE_ACTIVE | 1 | - |
| PROJECT_GRANT_STATE_INACTIVE | 2 | - |




#### ProjectState {#projectstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| PROJECT_STATE_UNSPECIFIED | 0 | - |
| PROJECT_STATE_ACTIVE | 1 | - |
| PROJECT_STATE_INACTIVE | 2 | - |


 <!-- end Enums -->
 <!-- end if Enums -->


 <!-- end services -->

### Messages


#### AuthFactor



| Field | Type | Description |
| ----- | ---- | ----------- |
| state | [ AuthFactorState](#authfactorstate) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.otp | [ AuthFactorOTP](#authfactorotp) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.u2f | [ AuthFactorU2F](#authfactoru2f) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### AuthFactorOTP


 <!-- end HasFields -->


#### AuthFactorU2F



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### DisplayNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| display_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Email



| Field | Type | Description |
| ----- | ---- | ----------- |
| email | [ string](#string) | - |
| is_email_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### EmailQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| email_address | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### FirstNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Human



| Field | Type | Description |
| ----- | ---- | ----------- |
| profile | [ Profile](#profile) | - |
| email | [ Email](#email) | - |
| phone | [ Phone](#phone) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### LastNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| last_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Machine



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [ string](#string) | - |
| description | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Membership



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| roles | [repeated string](#string) | - |
| display_name | [ string](#string) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.iam | [ bool](#bool) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.org_id | [ string](#string) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.project_id | [ string](#string) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.project_grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### MembershipIAMQuery
this query is always equals


| Field | Type | Description |
| ----- | ---- | ----------- |
| iam | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### MembershipOrgQuery
this query is always equals


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### MembershipProjectGrantQuery
this query is always equals


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### MembershipProjectQuery
this query is always equals


| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### MembershipQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.org_query | [ MembershipOrgQuery](#membershiporgquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_query | [ MembershipProjectQuery](#membershipprojectquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_grant_query | [ MembershipProjectGrantQuery](#membershipprojectgrantquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.iam_query | [ MembershipIAMQuery](#membershipiamquery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### NickNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| nick_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Phone



| Field | Type | Description |
| ----- | ---- | ----------- |
| phone | [ string](#string) | - |
| is_phone_verified | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Profile



| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name | [ string](#string) | - |
| last_name | [ string](#string) | - |
| nick_name | [ string](#string) | - |
| display_name | [ string](#string) | - |
| preferred_language | [ string](#string) | - |
| gender | [ Gender](#gender) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### SearchQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.user_name_query | [ UserNameQuery](#usernamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.first_name_query | [ FirstNameQuery](#firstnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.last_name_query | [ LastNameQuery](#lastnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.nick_name_query | [ NickNameQuery](#nicknamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.display_name_query | [ DisplayNameQuery](#displaynamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.email_query | [ EmailQuery](#emailquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.state_query | [ StateQuery](#statequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.type_query | [ TypeQuery](#typequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### Session



| Field | Type | Description |
| ----- | ---- | ----------- |
| session_id | [ string](#string) | - |
| agent_id | [ string](#string) | - |
| auth_state | [ SessionState](#sessionstate) | - |
| user_id | [ string](#string) | - |
| user_name | [ string](#string) | - |
| login_name | [ string](#string) | - |
| display_name | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### StateQuery
UserStateQuery is always equals


| Field | Type | Description |
| ----- | ---- | ----------- |
| state | [ UserState](#userstate) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### TypeQuery
UserTypeQuery is always equals


| Field | Type | Description |
| ----- | ---- | ----------- |
| type | [ Type](#type) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### User



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| state | [ UserState](#userstate) | - |
| user_name | [ string](#string) | - |
| login_names | [repeated string](#string) | - |
| preferred_login_name | [ string](#string) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.human | [ Human](#human) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) type.machine | [ Machine](#machine) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrant



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| details | [ zitadel.v1.ObjectDetails](#zitadelv1objectdetails) | - |
| role_keys | [repeated string](#string) | - |
| state | [ UserGrantState](#usergrantstate) | - |
| user_id | [ string](#string) | - |
| user_name | [ string](#string) | - |
| first_name | [ string](#string) | - |
| last_name | [ string](#string) | - |
| email | [ string](#string) | - |
| display_name | [ string](#string) | - |
| org_id | [ string](#string) | - |
| org_name | [ string](#string) | - |
| org_domain | [ string](#string) | - |
| project_id | [ string](#string) | - |
| project_name | [ string](#string) | - |
| project_grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantDisplayNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| display_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantEmailQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| email | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantFirstNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantLastNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| last_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantOrgDomainQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_domain | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantOrgNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| org_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantProjectGrantIDQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_grant_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantProjectIDQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantProjectNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| project_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_id_query | [ UserGrantProjectIDQuery](#usergrantprojectidquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.user_id_query | [ UserGrantUserIDQuery](#usergrantuseridquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.with_granted_query | [ UserGrantWithGrantedQuery](#usergrantwithgrantedquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.role_key_query | [ UserGrantRoleKeyQuery](#usergrantrolekeyquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_grant_id_query | [ UserGrantProjectGrantIDQuery](#usergrantprojectgrantidquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.user_name_query | [ UserGrantUserNameQuery](#usergrantusernamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.first_name_query | [ UserGrantFirstNameQuery](#usergrantfirstnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.last_name_query | [ UserGrantLastNameQuery](#usergrantlastnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.email_query | [ UserGrantEmailQuery](#usergrantemailquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.org_name_query | [ UserGrantOrgNameQuery](#usergrantorgnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.org_domain_query | [ UserGrantOrgDomainQuery](#usergrantorgdomainquery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.project_name_query | [ UserGrantProjectNameQuery](#usergrantprojectnamequery) | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.display_name_query | [ UserGrantDisplayNameQuery](#usergrantdisplaynamequery) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantRoleKeyQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| role_key | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantUserIDQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantUserNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserGrantWithGrantedQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| with_granted | [ bool](#bool) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### UserNameQuery



| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name | [ string](#string) | - |
| method | [ zitadel.v1.TextQueryMethod](#zitadelv1textquerymethod) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### WebAuthNKey



| Field | Type | Description |
| ----- | ---- | ----------- |
| public_key | [ bytes](#bytes) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### WebAuthNToken



| Field | Type | Description |
| ----- | ---- | ----------- |
| id | [ string](#string) | - |
| state | [ AuthFactorState](#authfactorstate) | - |
| name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->


#### WebAuthNVerification



| Field | Type | Description |
| ----- | ---- | ----------- |
| public_key_credential | [ bytes](#bytes) | - |
| token_name | [ string](#string) | - |
 <!-- end Fields -->
 <!-- end HasFields -->
 <!-- end messages -->


### Enums


#### AuthFactorState {#authfactorstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| AUTH_FACTOR_STATE_UNSPECIFIED | 0 | - |
| AUTH_FACTOR_STATE_NOT_READY | 1 | - |
| AUTH_FACTOR_STATE_READY | 2 | - |
| AUTH_FACTOR_STATE_REMOVED | 3 | - |




#### Gender {#gender}


| Name | Number | Description |
| ---- | ------ | ----------- |
| GENDER_UNSPECIFIED | 0 | - |
| GENDER_FEMALE | 1 | - |
| GENDER_MALE | 2 | - |
| GENDER_DIVERSE | 3 | - |




#### SessionState {#sessionstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| SESSION_STATE_UNSPECIFIED | 0 | - |
| SESSION_STATE_ACTIVE | 1 | - |
| SESSION_STATE_TERMINATED | 2 | - |




#### Type {#type}


| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNSPECIFIED | 0 | - |
| TYPE_HUMAN | 1 | - |
| TYPE_MACHINE | 2 | - |




#### UserFieldName {#userfieldname}


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_FIELD_NAME_UNSPECIFIED | 0 | - |
| USER_FIELD_NAME_USER_NAME | 1 | - |
| USER_FIELD_NAME_FIRST_NAME | 2 | - |
| USER_FIELD_NAME_LAST_NAME | 3 | - |
| USER_FIELD_NAME_NICK_NAME | 4 | - |
| USER_FIELD_NAME_DISPLAY_NAME | 5 | - |
| USER_FIELD_NAME_EMAIL | 6 | - |
| USER_FIELD_NAME_STATE | 7 | - |
| USER_FIELD_NAME_TYPE | 8 | - |




#### UserGrantState {#usergrantstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_GRANT_STATE_UNSPECIFIED | 0 | - |
| USER_GRANT_STATE_ACTIVE | 1 | - |
| USER_GRANT_STATE_INACTIVE | 2 | - |




#### UserState {#userstate}


| Name | Number | Description |
| ---- | ------ | ----------- |
| USER_STATE_UNSPECIFIED | 0 | - |
| USER_STATE_ACTIVE | 1 | - |
| USER_STATE_INACTIVE | 2 | - |
| USER_STATE_DELETED | 3 | - |
| USER_STATE_LOCKED | 4 | - |
| USER_STATE_SUSPEND | 5 | - |
| USER_STATE_INITIAL | 6 | - |


 <!-- end Enums -->
 <!-- end if Enums -->
 <!-- end Files -->

### Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <div><h4 id="double" /></div><a name="double" /> double |  | double | double | float |
| <div><h4 id="float" /></div><a name="float" /> float |  | float | float | float |
| <div><h4 id="int32" /></div><a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <div><h4 id="int64" /></div><a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <div><h4 id="uint32" /></div><a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <div><h4 id="uint64" /></div><a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <div><h4 id="sint32" /></div><a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <div><h4 id="sint64" /></div><a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <div><h4 id="fixed32" /></div><a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <div><h4 id="fixed64" /></div><a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <div><h4 id="sfixed32" /></div><a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <div><h4 id="sfixed64" /></div><a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <div><h4 id="bool" /></div><a name="bool" /> bool |  | bool | boolean | boolean |
| <div><h4 id="string" /></div><a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <div><h4 id="bytes" /></div><a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |


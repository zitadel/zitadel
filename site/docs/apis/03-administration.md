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




- [Scalar Value Types](#scalar-value-types)



### AdminService {#zitadeladminv1adminservice}


#### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)

Indicates if ZITADEL is running.
It respondes as soon as ZITADEL started




GET: /healthz


#### IsOrgUnique

> **rpc** IsOrgUnique([IsOrgUniqueRequest](#isorguniquerequest))
[IsOrgUniqueResponse](#isorguniqueresponse)

Checks whether an organisation exists by the given parameters




GET: /orgs/_is_unique


#### GetOrgByID

> **rpc** GetOrgByID([GetOrgByIDRequest](#getorgbyidrequest))
[GetOrgByIDResponse](#getorgbyidresponse)






GET: /orgs/{id}


#### ListOrgs

> **rpc** ListOrgs([ListOrgsRequest](#listorgsrequest))
[ListOrgsResponse](#listorgsresponse)

Returns all organisations matching the request
all queries need to match (ANDed)




POST: /orgs/_search


#### SetUpOrg

> **rpc** SetUpOrg([SetUpOrgRequest](#setuporgrequest))
[SetUpOrgResponse](#setuporgresponse)

Creates a new org and user 
and adds the user to the orgs members as ORG_OWNER




POST: /orgs/_setup


#### GetIDPByID

> **rpc** GetIDPByID([GetIDPByIDRequest](#getidpbyidrequest))
[GetIDPByIDResponse](#getidpbyidresponse)






GET: /idps/{id}


#### ListIDPs

> **rpc** ListIDPs([ListIDPsRequest](#listidpsrequest))
[ListIDPsResponse](#listidpsresponse)






POST: /idps/_search


#### AddOIDCIDP

> **rpc** AddOIDCIDP([AddOIDCIDPRequest](#addoidcidprequest))
[AddOIDCIDPResponse](#addoidcidpresponse)






POST: /idps/oidc


#### UpdateIDP

> **rpc** UpdateIDP([UpdateIDPRequest](#updateidprequest))
[UpdateIDPResponse](#updateidpresponse)

Updates the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.




PUT: /idps/{idp_id}


#### DeactivateIDP

> **rpc** DeactivateIDP([DeactivateIDPRequest](#deactivateidprequest))
[DeactivateIDPResponse](#deactivateidpresponse)

Sets the state of the idp to IDP_STATE_INACTIVE
the state MUST be IDP_STATE_ACTIVE for this call




POST: /idps/{idp_id}/_deactivate


#### ReactivateIDP

> **rpc** ReactivateIDP([ReactivateIDPRequest](#reactivateidprequest))
[ReactivateIDPResponse](#reactivateidpresponse)

Sets the state of the idp to IDP_STATE_ACTIVE
the state MUST be IDP_STATE_INACTIVE for this call




POST: /idps/{idp_id}/_reactivate


#### RemoveIDP

> **rpc** RemoveIDP([RemoveIDPRequest](#removeidprequest))
[RemoveIDPResponse](#removeidpresponse)

RemoveIDP deletes the IDP permanetly




DELETE: /idps/{idp_id}


#### UpdateIDPOIDCConfig

> **rpc** UpdateIDPOIDCConfig([UpdateIDPOIDCConfigRequest](#updateidpoidcconfigrequest))
[UpdateIDPOIDCConfigResponse](#updateidpoidcconfigresponse)

Updates the oidc configuration of the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.




PUT: /idps/{idp_id}/oidc_config


#### GetDefaultFeatures

> **rpc** GetDefaultFeatures([GetDefaultFeaturesRequest](#getdefaultfeaturesrequest))
[GetDefaultFeaturesResponse](#getdefaultfeaturesresponse)






GET: /features


#### SetDefaultFeatures

> **rpc** SetDefaultFeatures([SetDefaultFeaturesRequest](#setdefaultfeaturesrequest))
[SetDefaultFeaturesResponse](#setdefaultfeaturesresponse)






PUT: /features


#### GetOrgFeatures

> **rpc** GetOrgFeatures([GetOrgFeaturesRequest](#getorgfeaturesrequest))
[GetOrgFeaturesResponse](#getorgfeaturesresponse)






GET: /orgs/{org_id}/features


#### SetOrgFeatures

> **rpc** SetOrgFeatures([SetOrgFeaturesRequest](#setorgfeaturesrequest))
[SetOrgFeaturesResponse](#setorgfeaturesresponse)






PUT: /orgs/{org_id}/features


#### ResetOrgFeatures

> **rpc** ResetOrgFeatures([ResetOrgFeaturesRequest](#resetorgfeaturesrequest))
[ResetOrgFeaturesResponse](#resetorgfeaturesresponse)






DELETE: /orgs/{org_id}/features


#### GetOrgIAMPolicy

> **rpc** GetOrgIAMPolicy([GetOrgIAMPolicyRequest](#getorgiampolicyrequest))
[GetOrgIAMPolicyResponse](#getorgiampolicyresponse)

Returns the IAM policy defined by the administrators of ZITADEL




GET: /policies/orgiam


#### UpdateOrgIAMPolicy

> **rpc** UpdateOrgIAMPolicy([UpdateOrgIAMPolicyRequest](#updateorgiampolicyrequest))
[UpdateOrgIAMPolicyResponse](#updateorgiampolicyresponse)

Updates the default IAM policy.
it impacts all organisations without a customised policy




PUT: /policies/orgiam


#### GetCustomOrgIAMPolicy

> **rpc** GetCustomOrgIAMPolicy([GetCustomOrgIAMPolicyRequest](#getcustomorgiampolicyrequest))
[GetCustomOrgIAMPolicyResponse](#getcustomorgiampolicyresponse)

Returns the customised policy or the default if not customised




GET: /orgs/{org_id}/policies/orgiam


#### AddCustomOrgIAMPolicy

> **rpc** AddCustomOrgIAMPolicy([AddCustomOrgIAMPolicyRequest](#addcustomorgiampolicyrequest))
[AddCustomOrgIAMPolicyResponse](#addcustomorgiampolicyresponse)

Defines a custom ORGIAM policy as specified




POST: /orgs/{org_id}/policies/orgiam


#### UpdateCustomOrgIAMPolicy

> **rpc** UpdateCustomOrgIAMPolicy([UpdateCustomOrgIAMPolicyRequest](#updatecustomorgiampolicyrequest))
[UpdateCustomOrgIAMPolicyResponse](#updatecustomorgiampolicyresponse)

Updates a custom ORGIAM policy as specified




PUT: /orgs/{org_id}/policies/orgiam


#### ResetCustomOrgIAMPolicyToDefault

> **rpc** ResetCustomOrgIAMPolicyToDefault([ResetCustomOrgIAMPolicyToDefaultRequest](#resetcustomorgiampolicytodefaultrequest))
[ResetCustomOrgIAMPolicyToDefaultResponse](#resetcustomorgiampolicytodefaultresponse)

Resets the org iam policy of the organisation to default
ZITADEL will fallback to the default policy defined by the ZITADEL administrators




DELETE: /orgs/{org_id}/policies/orgiam


#### GetLabelPolicy

> **rpc** GetLabelPolicy([GetLabelPolicyRequest](#getlabelpolicyrequest))
[GetLabelPolicyResponse](#getlabelpolicyresponse)

Returns the label policy defined by the administrators of ZITADEL




GET: /policies/label


#### UpdateLabelPolicy

> **rpc** UpdateLabelPolicy([UpdateLabelPolicyRequest](#updatelabelpolicyrequest))
[UpdateLabelPolicyResponse](#updatelabelpolicyresponse)

Updates the default label policy of ZITADEL
it impacts all organisations without a customised policy




PUT: /policies/label


#### GetLoginPolicy

> **rpc** GetLoginPolicy([GetLoginPolicyRequest](#getloginpolicyrequest))
[GetLoginPolicyResponse](#getloginpolicyresponse)

Returns the login policy defined by the administrators of ZITADEL




GET: /policies/login


#### UpdateLoginPolicy

> **rpc** UpdateLoginPolicy([UpdateLoginPolicyRequest](#updateloginpolicyrequest))
[UpdateLoginPolicyResponse](#updateloginpolicyresponse)

Updates the default login policy of ZITADEL
it impacts all organisations without a customised policy




PUT: /policies/login


#### ListLoginPolicyIDPs

> **rpc** ListLoginPolicyIDPs([ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest))
[ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)

Returns the idps linked to the default login policy,
defined by the administrators of ZITADEL




POST: /policies/login/idps/_search


#### AddIDPToLoginPolicy

> **rpc** AddIDPToLoginPolicy([AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest))
[AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)

Adds the povided idp to the default login policy.
It impacts all organisations without a customised policy




POST: /policies/login/idps


#### RemoveIDPFromLoginPolicy

> **rpc** RemoveIDPFromLoginPolicy([RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest))
[RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)

Removes the povided idp from the default login policy.
It impacts all organisations without a customised policy




DELETE: /policies/login/idps/{idp_id}


#### ListLoginPolicySecondFactors

> **rpc** ListLoginPolicySecondFactors([ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest))
[ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)

Returns the available second factors defined by the administrators of ZITADEL




POST: /policies/login/second_factors/_search


#### AddSecondFactorToLoginPolicy

> **rpc** AddSecondFactorToLoginPolicy([AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest))
[AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)

Adds a second factor to the default login policy.
It impacts all organisations without a customised policy




POST: /policies/login/second_factors


#### RemoveSecondFactorFromLoginPolicy

> **rpc** RemoveSecondFactorFromLoginPolicy([RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest))
[RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)

Removes a second factor from the default login policy.
It impacts all organisations without a customised policy




DELETE: /policies/login/second_factors/{type}


#### ListLoginPolicyMultiFactors

> **rpc** ListLoginPolicyMultiFactors([ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest))
[ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)

Returns the available multi factors defined by the administrators of ZITADEL




POST: /policies/login/multi_factors/_search


#### AddMultiFactorToLoginPolicy

> **rpc** AddMultiFactorToLoginPolicy([AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest))
[AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)

Adds a multi factor to the default login policy.
It impacts all organisations without a customised policy




POST: /policies/login/multi_factors


#### RemoveMultiFactorFromLoginPolicy

> **rpc** RemoveMultiFactorFromLoginPolicy([RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest))
[RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)

Removes a multi factor from the default login policy.
It impacts all organisations without a customised policy




DELETE: /policies/login/multi_factors/{type}


#### GetPasswordComplexityPolicy

> **rpc** GetPasswordComplexityPolicy([GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest))
[GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)

Returns the password complexity policy defined by the administrators of ZITADEL




GET: /policies/password/complexity


#### UpdatePasswordComplexityPolicy

> **rpc** UpdatePasswordComplexityPolicy([UpdatePasswordComplexityPolicyRequest](#updatepasswordcomplexitypolicyrequest))
[UpdatePasswordComplexityPolicyResponse](#updatepasswordcomplexitypolicyresponse)

Updates the default password complexity policy of ZITADEL
it impacts all organisations without a customised policy




PUT: /policies/password/complexity


#### GetPasswordAgePolicy

> **rpc** GetPasswordAgePolicy([GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest))
[GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)

Returns the password age policy defined by the administrators of ZITADEL




GET: /policies/password/age


#### UpdatePasswordAgePolicy

> **rpc** UpdatePasswordAgePolicy([UpdatePasswordAgePolicyRequest](#updatepasswordagepolicyrequest))
[UpdatePasswordAgePolicyResponse](#updatepasswordagepolicyresponse)

Updates the default password age policy of ZITADEL
it impacts all organisations without a customised policy




PUT: /policies/password/age


#### GetPasswordLockoutPolicy

> **rpc** GetPasswordLockoutPolicy([GetPasswordLockoutPolicyRequest](#getpasswordlockoutpolicyrequest))
[GetPasswordLockoutPolicyResponse](#getpasswordlockoutpolicyresponse)

Returns the password lockout policy defined by the administrators of ZITADEL




GET: /policies/password/lockout


#### UpdatePasswordLockoutPolicy

> **rpc** UpdatePasswordLockoutPolicy([UpdatePasswordLockoutPolicyRequest](#updatepasswordlockoutpolicyrequest))
[UpdatePasswordLockoutPolicyResponse](#updatepasswordlockoutpolicyresponse)

Updates the default password lockout policy of ZITADEL
it impacts all organisations without a customised policy




PUT: /policies/password/lockout


#### ListIAMMemberRoles

> **rpc** ListIAMMemberRoles([ListIAMMemberRolesRequest](#listiammemberrolesrequest))
[ListIAMMemberRolesResponse](#listiammemberrolesresponse)

Returns the IAM roles visible for the requested user




POST: /members/roles/_search


#### ListIAMMembers

> **rpc** ListIAMMembers([ListIAMMembersRequest](#listiammembersrequest))
[ListIAMMembersResponse](#listiammembersresponse)

Returns all members matching the request
all queries need to match (ANDed)




POST: /members/_search


#### AddIAMMember

> **rpc** AddIAMMember([AddIAMMemberRequest](#addiammemberrequest))
[AddIAMMemberResponse](#addiammemberresponse)

Adds a user to the membership list of ZITADEL with the given roles
undefined roles will be dropped




POST: /members


#### UpdateIAMMember

> **rpc** UpdateIAMMember([UpdateIAMMemberRequest](#updateiammemberrequest))
[UpdateIAMMemberResponse](#updateiammemberresponse)

Sets the given roles on a member.
The member has only roles provided by this call




PUT: /members/{user_id}


#### RemoveIAMMember

> **rpc** RemoveIAMMember([RemoveIAMMemberRequest](#removeiammemberrequest))
[RemoveIAMMemberResponse](#removeiammemberresponse)

Removes the user from the membership list of ZITADEL




DELETE: /members/{user_id}


#### ListViews

> **rpc** ListViews([ListViewsRequest](#listviewsrequest))
[ListViewsResponse](#listviewsresponse)

Returns all stored read models of ZITADEL
views are used for search optimisation and optimise request latencies
they represent the delta of the event happend on the objects




POST: /views/_search


#### ClearView

> **rpc** ClearView([ClearViewRequest](#clearviewrequest))
[ClearViewResponse](#clearviewresponse)

Truncates the delta of the change stream
be carefull with this function because ZITADEL has to 
recompute the deltas after they got cleared. 
Search requests will return wrong results until all deltas are recomputed




POST: /views/{database}/{view_name}


#### ListFailedEvents

> **rpc** ListFailedEvents([ListFailedEventsRequest](#listfailedeventsrequest))
[ListFailedEventsResponse](#listfailedeventsresponse)

Returns event descriptions which cannot be processed.
It's possible that some events need some retries. 
For example if the SMTP-API wasn't able to send an email at the first time




POST: /failedevents/_search


#### RemoveFailedEvent

> **rpc** RemoveFailedEvent([RemoveFailedEventRequest](#removefailedeventrequest))
[RemoveFailedEventResponse](#removefailedeventresponse)

Deletes the event from failed events view.
the event is not removed from the change stream
This call is usefull if the system was able to process the event later. 
e.g. if the second try of sending an email was successful. the first try produced a
failed event. You can find out if it worked on the `failure_count`




DELETE: /failedevents/{database}/{view_name}/{failed_sequence}


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
if name or domain is already in use, org is not unique


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


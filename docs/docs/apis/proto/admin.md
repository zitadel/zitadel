---
title: zitadel/admin.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)


## AdminService {#zitadeladminv1adminservice}


### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)

Indicates if ZITADEL is running.
It respondes as soon as ZITADEL started



    GET: /healthz


### IsOrgUnique

> **rpc** IsOrgUnique([IsOrgUniqueRequest](#isorguniquerequest))
[IsOrgUniqueResponse](#isorguniqueresponse)

Checks whether an organisation exists by the given parameters



    GET: /orgs/_is_unique


### GetOrgByID

> **rpc** GetOrgByID([GetOrgByIDRequest](#getorgbyidrequest))
[GetOrgByIDResponse](#getorgbyidresponse)





    GET: /orgs/{id}


### ListOrgs

> **rpc** ListOrgs([ListOrgsRequest](#listorgsrequest))
[ListOrgsResponse](#listorgsresponse)

Returns all organisations matching the request
all queries need to match (ANDed)



    POST: /orgs/_search


### SetUpOrg

> **rpc** SetUpOrg([SetUpOrgRequest](#setuporgrequest))
[SetUpOrgResponse](#setuporgresponse)

Creates a new org and user 
and adds the user to the orgs members as ORG_OWNER



    POST: /orgs/_setup


### GetIDPByID

> **rpc** GetIDPByID([GetIDPByIDRequest](#getidpbyidrequest))
[GetIDPByIDResponse](#getidpbyidresponse)





    GET: /idps/{id}


### ListIDPs

> **rpc** ListIDPs([ListIDPsRequest](#listidpsrequest))
[ListIDPsResponse](#listidpsresponse)





    POST: /idps/_search


### AddOIDCIDP

> **rpc** AddOIDCIDP([AddOIDCIDPRequest](#addoidcidprequest))
[AddOIDCIDPResponse](#addoidcidpresponse)





    POST: /idps/oidc


### UpdateIDP

> **rpc** UpdateIDP([UpdateIDPRequest](#updateidprequest))
[UpdateIDPResponse](#updateidpresponse)

Updates the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.



    PUT: /idps/{idp_id}


### DeactivateIDP

> **rpc** DeactivateIDP([DeactivateIDPRequest](#deactivateidprequest))
[DeactivateIDPResponse](#deactivateidpresponse)

Sets the state of the idp to IDP_STATE_INACTIVE
the state MUST be IDP_STATE_ACTIVE for this call



    POST: /idps/{idp_id}/_deactivate


### ReactivateIDP

> **rpc** ReactivateIDP([ReactivateIDPRequest](#reactivateidprequest))
[ReactivateIDPResponse](#reactivateidpresponse)

Sets the state of the idp to IDP_STATE_ACTIVE
the state MUST be IDP_STATE_INACTIVE for this call



    POST: /idps/{idp_id}/_reactivate


### RemoveIDP

> **rpc** RemoveIDP([RemoveIDPRequest](#removeidprequest))
[RemoveIDPResponse](#removeidpresponse)

RemoveIDP deletes the IDP permanetly



    DELETE: /idps/{idp_id}


### UpdateIDPOIDCConfig

> **rpc** UpdateIDPOIDCConfig([UpdateIDPOIDCConfigRequest](#updateidpoidcconfigrequest))
[UpdateIDPOIDCConfigResponse](#updateidpoidcconfigresponse)

Updates the oidc configuration of the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.



    PUT: /idps/{idp_id}/oidc_config


### GetDefaultFeatures

> **rpc** GetDefaultFeatures([GetDefaultFeaturesRequest](#getdefaultfeaturesrequest))
[GetDefaultFeaturesResponse](#getdefaultfeaturesresponse)





    GET: /features


### SetDefaultFeatures

> **rpc** SetDefaultFeatures([SetDefaultFeaturesRequest](#setdefaultfeaturesrequest))
[SetDefaultFeaturesResponse](#setdefaultfeaturesresponse)





    PUT: /features


### GetOrgFeatures

> **rpc** GetOrgFeatures([GetOrgFeaturesRequest](#getorgfeaturesrequest))
[GetOrgFeaturesResponse](#getorgfeaturesresponse)





    GET: /orgs/{org_id}/features


### SetOrgFeatures

> **rpc** SetOrgFeatures([SetOrgFeaturesRequest](#setorgfeaturesrequest))
[SetOrgFeaturesResponse](#setorgfeaturesresponse)





    PUT: /orgs/{org_id}/features


### ResetOrgFeatures

> **rpc** ResetOrgFeatures([ResetOrgFeaturesRequest](#resetorgfeaturesrequest))
[ResetOrgFeaturesResponse](#resetorgfeaturesresponse)





    DELETE: /orgs/{org_id}/features


### GetOrgIAMPolicy

> **rpc** GetOrgIAMPolicy([GetOrgIAMPolicyRequest](#getorgiampolicyrequest))
[GetOrgIAMPolicyResponse](#getorgiampolicyresponse)

Returns the IAM policy defined by the administrators of ZITADEL



    GET: /policies/orgiam


### UpdateOrgIAMPolicy

> **rpc** UpdateOrgIAMPolicy([UpdateOrgIAMPolicyRequest](#updateorgiampolicyrequest))
[UpdateOrgIAMPolicyResponse](#updateorgiampolicyresponse)

Updates the default IAM policy.
it impacts all organisations without a customised policy



    PUT: /policies/orgiam


### GetCustomOrgIAMPolicy

> **rpc** GetCustomOrgIAMPolicy([GetCustomOrgIAMPolicyRequest](#getcustomorgiampolicyrequest))
[GetCustomOrgIAMPolicyResponse](#getcustomorgiampolicyresponse)

Returns the customised policy or the default if not customised



    GET: /orgs/{org_id}/policies/orgiam


### AddCustomOrgIAMPolicy

> **rpc** AddCustomOrgIAMPolicy([AddCustomOrgIAMPolicyRequest](#addcustomorgiampolicyrequest))
[AddCustomOrgIAMPolicyResponse](#addcustomorgiampolicyresponse)

Defines a custom ORGIAM policy as specified



    POST: /orgs/{org_id}/policies/orgiam


### UpdateCustomOrgIAMPolicy

> **rpc** UpdateCustomOrgIAMPolicy([UpdateCustomOrgIAMPolicyRequest](#updatecustomorgiampolicyrequest))
[UpdateCustomOrgIAMPolicyResponse](#updatecustomorgiampolicyresponse)

Updates a custom ORGIAM policy as specified



    PUT: /orgs/{org_id}/policies/orgiam


### ResetCustomOrgIAMPolicyToDefault

> **rpc** ResetCustomOrgIAMPolicyToDefault([ResetCustomOrgIAMPolicyToDefaultRequest](#resetcustomorgiampolicytodefaultrequest))
[ResetCustomOrgIAMPolicyToDefaultResponse](#resetcustomorgiampolicytodefaultresponse)

Resets the org iam policy of the organisation to default
ZITADEL will fallback to the default policy defined by the ZITADEL administrators



    DELETE: /orgs/{org_id}/policies/orgiam


### GetLabelPolicy

> **rpc** GetLabelPolicy([GetLabelPolicyRequest](#getlabelpolicyrequest))
[GetLabelPolicyResponse](#getlabelpolicyresponse)

Returns the label policy defined by the administrators of ZITADEL



    GET: /policies/label


### UpdateLabelPolicy

> **rpc** UpdateLabelPolicy([UpdateLabelPolicyRequest](#updatelabelpolicyrequest))
[UpdateLabelPolicyResponse](#updatelabelpolicyresponse)

Updates the default label policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/label


### GetLoginPolicy

> **rpc** GetLoginPolicy([GetLoginPolicyRequest](#getloginpolicyrequest))
[GetLoginPolicyResponse](#getloginpolicyresponse)

Returns the login policy defined by the administrators of ZITADEL



    GET: /policies/login


### UpdateLoginPolicy

> **rpc** UpdateLoginPolicy([UpdateLoginPolicyRequest](#updateloginpolicyrequest))
[UpdateLoginPolicyResponse](#updateloginpolicyresponse)

Updates the default login policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/login


### ListLoginPolicyIDPs

> **rpc** ListLoginPolicyIDPs([ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest))
[ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)

Returns the idps linked to the default login policy,
defined by the administrators of ZITADEL



    POST: /policies/login/idps/_search


### AddIDPToLoginPolicy

> **rpc** AddIDPToLoginPolicy([AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest))
[AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)

Adds the povided idp to the default login policy.
It impacts all organisations without a customised policy



    POST: /policies/login/idps


### RemoveIDPFromLoginPolicy

> **rpc** RemoveIDPFromLoginPolicy([RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest))
[RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)

Removes the povided idp from the default login policy.
It impacts all organisations without a customised policy



    DELETE: /policies/login/idps/{idp_id}


### ListLoginPolicySecondFactors

> **rpc** ListLoginPolicySecondFactors([ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest))
[ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)

Returns the available second factors defined by the administrators of ZITADEL



    POST: /policies/login/second_factors/_search


### AddSecondFactorToLoginPolicy

> **rpc** AddSecondFactorToLoginPolicy([AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest))
[AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)

Adds a second factor to the default login policy.
It impacts all organisations without a customised policy



    POST: /policies/login/second_factors


### RemoveSecondFactorFromLoginPolicy

> **rpc** RemoveSecondFactorFromLoginPolicy([RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest))
[RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)

Removes a second factor from the default login policy.
It impacts all organisations without a customised policy



    DELETE: /policies/login/second_factors/{type}


### ListLoginPolicyMultiFactors

> **rpc** ListLoginPolicyMultiFactors([ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest))
[ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)

Returns the available multi factors defined by the administrators of ZITADEL



    POST: /policies/login/multi_factors/_search


### AddMultiFactorToLoginPolicy

> **rpc** AddMultiFactorToLoginPolicy([AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest))
[AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)

Adds a multi factor to the default login policy.
It impacts all organisations without a customised policy



    POST: /policies/login/multi_factors


### RemoveMultiFactorFromLoginPolicy

> **rpc** RemoveMultiFactorFromLoginPolicy([RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest))
[RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)

Removes a multi factor from the default login policy.
It impacts all organisations without a customised policy



    DELETE: /policies/login/multi_factors/{type}


### GetPasswordComplexityPolicy

> **rpc** GetPasswordComplexityPolicy([GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest))
[GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)

Returns the password complexity policy defined by the administrators of ZITADEL



    GET: /policies/password/complexity


### UpdatePasswordComplexityPolicy

> **rpc** UpdatePasswordComplexityPolicy([UpdatePasswordComplexityPolicyRequest](#updatepasswordcomplexitypolicyrequest))
[UpdatePasswordComplexityPolicyResponse](#updatepasswordcomplexitypolicyresponse)

Updates the default password complexity policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/password/complexity


### GetPasswordAgePolicy

> **rpc** GetPasswordAgePolicy([GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest))
[GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)

Returns the password age policy defined by the administrators of ZITADEL



    GET: /policies/password/age


### UpdatePasswordAgePolicy

> **rpc** UpdatePasswordAgePolicy([UpdatePasswordAgePolicyRequest](#updatepasswordagepolicyrequest))
[UpdatePasswordAgePolicyResponse](#updatepasswordagepolicyresponse)

Updates the default password age policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/password/age


### GetPasswordLockoutPolicy

> **rpc** GetPasswordLockoutPolicy([GetPasswordLockoutPolicyRequest](#getpasswordlockoutpolicyrequest))
[GetPasswordLockoutPolicyResponse](#getpasswordlockoutpolicyresponse)

Returns the password lockout policy defined by the administrators of ZITADEL



    GET: /policies/password/lockout


### UpdatePasswordLockoutPolicy

> **rpc** UpdatePasswordLockoutPolicy([UpdatePasswordLockoutPolicyRequest](#updatepasswordlockoutpolicyrequest))
[UpdatePasswordLockoutPolicyResponse](#updatepasswordlockoutpolicyresponse)

Updates the default password lockout policy of ZITADEL
it impacts all organisations without a customised policy



    PUT: /policies/password/lockout


### ListIAMMemberRoles

> **rpc** ListIAMMemberRoles([ListIAMMemberRolesRequest](#listiammemberrolesrequest))
[ListIAMMemberRolesResponse](#listiammemberrolesresponse)

Returns the IAM roles visible for the requested user



    POST: /members/roles/_search


### ListIAMMembers

> **rpc** ListIAMMembers([ListIAMMembersRequest](#listiammembersrequest))
[ListIAMMembersResponse](#listiammembersresponse)

Returns all members matching the request
all queries need to match (ANDed)



    POST: /members/_search


### AddIAMMember

> **rpc** AddIAMMember([AddIAMMemberRequest](#addiammemberrequest))
[AddIAMMemberResponse](#addiammemberresponse)

Adds a user to the membership list of ZITADEL with the given roles
undefined roles will be dropped



    POST: /members


### UpdateIAMMember

> **rpc** UpdateIAMMember([UpdateIAMMemberRequest](#updateiammemberrequest))
[UpdateIAMMemberResponse](#updateiammemberresponse)

Sets the given roles on a member.
The member has only roles provided by this call



    PUT: /members/{user_id}


### RemoveIAMMember

> **rpc** RemoveIAMMember([RemoveIAMMemberRequest](#removeiammemberrequest))
[RemoveIAMMemberResponse](#removeiammemberresponse)

Removes the user from the membership list of ZITADEL



    DELETE: /members/{user_id}


### ListViews

> **rpc** ListViews([ListViewsRequest](#listviewsrequest))
[ListViewsResponse](#listviewsresponse)

Returns all stored read models of ZITADEL
views are used for search optimisation and optimise request latencies
they represent the delta of the event happend on the objects



    POST: /views/_search


### ClearView

> **rpc** ClearView([ClearViewRequest](#clearviewrequest))
[ClearViewResponse](#clearviewresponse)

Truncates the delta of the change stream
be carefull with this function because ZITADEL has to 
recompute the deltas after they got cleared. 
Search requests will return wrong results until all deltas are recomputed



    POST: /views/{database}/{view_name}


### ListFailedEvents

> **rpc** ListFailedEvents([ListFailedEventsRequest](#listfailedeventsrequest))
[ListFailedEventsResponse](#listfailedeventsresponse)

Returns event descriptions which cannot be processed.
It's possible that some events need some retries. 
For example if the SMTP-API wasn't able to send an email at the first time



    POST: /failedevents/_search


### RemoveFailedEvent

> **rpc** RemoveFailedEvent([RemoveFailedEventRequest](#removefailedeventrequest))
[RemoveFailedEventResponse](#removefailedeventresponse)

Deletes the event from failed events view.
the event is not removed from the change stream
This call is usefull if the system was able to process the event later. 
e.g. if the second try of sending an email was successful. the first try produced a
failed event. You can find out if it worked on the `failure_count`



    DELETE: /failedevents/{database}/{view_name}/{failed_sequence}







## Messages


### AddCustomOrgIAMPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id |  string | - |
| user_login_must_be_domain |  bool | the username has to end with the domain of it's organisation (uniqueness is organisation based) |



### AddCustomOrgIAMPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddIAMMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| roles | repeated string | if no roles provided the user won't have any rights |



### AddIAMMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddIDPToLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | Id of the predefined idp configuration |



### AddIDPToLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddMultiFactorToLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - |



### AddMultiFactorToLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### AddOIDCIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| name |  string | - |
| styling_type |  zitadel.idp.v1.IDPStylingType | - |
| client_id |  string | - |
| client_secret |  string | - |
| issuer |  string | - |
| scopes | repeated string | - |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - |



### AddOIDCIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| idp_id |  string | - |



### AddSecondFactorToLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - |



### AddSecondFactorToLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ClearViewRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| database |  string | - |
| view_name |  string | - |



### ClearViewResponse


| Field | Type | Description |
| ----- | ---- | ----------- |



### DeactivateIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |



### DeactivateIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### FailedEvent


| Field | Type | Description |
| ----- | ---- | ----------- |
| database |  string | - |
| view_name |  string | - |
| failed_sequence |  uint64 | - |
| failure_count |  uint64 | - |
| error_message |  string | - |



### GetCustomOrgIAMPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id |  string | - |



### GetCustomOrgIAMPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.OrgIAMPolicy | - |
| is_default |  bool | - |



### GetDefaultFeaturesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetDefaultFeaturesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| features |  zitadel.features.v1.Features | - |



### GetIDPByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### GetIDPByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp |  zitadel.idp.v1.IDP | - |



### GetLabelPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetLabelPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |



### GetLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.LoginPolicy | - |



### GetOrgByIDRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| id |  string | - |



### GetOrgByIDResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| org |  zitadel.org.v1.Org | - |



### GetOrgFeaturesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id |  string | - |



### GetOrgFeaturesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| features |  zitadel.features.v1.Features | - |



### GetOrgIAMPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetOrgIAMPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.OrgIAMPolicy | - |



### GetPasswordAgePolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetPasswordAgePolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordAgePolicy | - |



### GetPasswordComplexityPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetPasswordComplexityPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |



### GetPasswordLockoutPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetPasswordLockoutPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordLockoutPolicy | - |



### HealthzRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### HealthzResponse


| Field | Type | Description |
| ----- | ---- | ----------- |



### IDPQuery


| Field | Type | Description |
| ----- | ---- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_id_query |  zitadel.idp.v1.IDPIDQuery | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_name_query |  zitadel.idp.v1.IDPNameQuery | - |



### IsOrgUniqueRequest
if name or domain is already in use, org is not unique

| Field | Type | Description |
| ----- | ---- | ----------- |
| name |  string | - |
| domain |  string | - |



### IsOrgUniqueResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| is_unique |  bool | - |



### ListFailedEventsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListFailedEventsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated FailedEvent | TODO: list details |



### ListIAMMemberRolesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListIAMMemberRolesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| roles | repeated string | - |



### ListIAMMembersRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |



### ListIAMMembersResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.member.v1.Member | - |



### ListIDPsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| sorting_column |  zitadel.idp.v1.IDPFieldName | the field the result is sorted |
| queries | repeated IDPQuery | criterias the client is looking for |



### ListIDPsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| sorting_column |  zitadel.idp.v1.IDPFieldName | - |
| result | repeated zitadel.idp.v1.IDP | - |



### ListLoginPolicyIDPsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |



### ListLoginPolicyIDPsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.idp.v1.IDPLoginPolicyLink | - |



### ListLoginPolicyMultiFactorsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListLoginPolicyMultiFactorsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.policy.v1.MultiFactorType | - |



### ListLoginPolicySecondFactorsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListLoginPolicySecondFactorsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.policy.v1.SecondFactorType | - |



### ListOrgsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| sorting_column |  zitadel.org.v1.OrgFieldName | the field the result is sorted |
| queries | repeated zitadel.org.v1.OrgQuery | criterias the client is looking for |



### ListOrgsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| sorting_column |  zitadel.org.v1.OrgFieldName | - |
| result | repeated zitadel.org.v1.Org | - |



### ListViewsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListViewsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated View | TODO: list details |



### ReactivateIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |



### ReactivateIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveFailedEventRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| database |  string | - |
| view_name |  string | - |
| failed_sequence |  uint64 | - |



### RemoveFailedEventResponse


| Field | Type | Description |
| ----- | ---- | ----------- |



### RemoveIAMMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |



### RemoveIAMMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveIDPFromLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |



### RemoveIDPFromLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |



### RemoveIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveMultiFactorFromLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - |



### RemoveMultiFactorFromLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveSecondFactorFromLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - |



### RemoveSecondFactorFromLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResetCustomOrgIAMPolicyToDefaultRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id |  string | - |



### ResetCustomOrgIAMPolicyToDefaultResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResetOrgFeaturesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id |  string | - |



### ResetOrgFeaturesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### SetDefaultFeaturesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| tier_name |  string | - |
| description |  string | - |
| audit_log_retention |  google.protobuf.Duration | - |
| login_policy_username_login |  bool | - |
| login_policy_registration |  bool | - |
| login_policy_idp |  bool | - |
| login_policy_factors |  bool | - |
| login_policy_passwordless |  bool | - |
| password_complexity_policy |  bool | - |
| label_policy |  bool | - |



### SetDefaultFeaturesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### SetOrgFeaturesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id |  string | - |
| tier_name |  string | - |
| description |  string | - |
| state |  zitadel.features.v1.FeaturesState | - |
| state_description |  string | - |
| audit_log_retention |  google.protobuf.Duration | - |
| login_policy_username_login |  bool | - |
| login_policy_registration |  bool | - |
| login_policy_idp |  bool | - |
| login_policy_factors |  bool | - |
| login_policy_passwordless |  bool | - |
| password_complexity_policy |  bool | - |
| label_policy |  bool | - |



### SetOrgFeaturesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### SetUpOrgRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| org |  SetUpOrgRequest.Org | - |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) user.human |  SetUpOrgRequest.Human | oneof field for the user managing the organisation |



### SetUpOrgRequest.Human


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name |  string | - |
| profile |  SetUpOrgRequest.Human.Profile | - |
| email |  SetUpOrgRequest.Human.Email | - |
| phone |  SetUpOrgRequest.Human.Phone | - |
| password |  string | - |



### SetUpOrgRequest.Human.Email


| Field | Type | Description |
| ----- | ---- | ----------- |
| email |  string | TODO: check if no value is allowed |
| is_email_verified |  bool | - |



### SetUpOrgRequest.Human.Phone


| Field | Type | Description |
| ----- | ---- | ----------- |
| phone |  string | has to be a global number |
| is_phone_verified |  bool | - |



### SetUpOrgRequest.Human.Profile


| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name |  string | - |
| last_name |  string | - |
| nick_name |  string | - |
| display_name |  string | - |
| preferred_language |  string | - |
| gender |  zitadel.user.v1.Gender | - |



### SetUpOrgRequest.Org


| Field | Type | Description |
| ----- | ---- | ----------- |
| name |  string | - |
| domain |  string | - |



### SetUpOrgResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| org_id |  string | - |
| user_id |  string | - |



### UpdateCustomOrgIAMPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id |  string | - |
| user_login_must_be_domain |  bool | - |



### UpdateCustomOrgIAMPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateIAMMemberRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_id |  string | - |
| roles | repeated string | if no roles provided the user won't have any rights |



### UpdateIAMMemberResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateIDPOIDCConfigRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |
| issuer |  string | - |
| client_id |  string | - |
| client_secret |  string | - |
| scopes | repeated string | - |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - |



### UpdateIDPOIDCConfigResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |
| name |  string | - |
| styling_type |  zitadel.idp.v1.IDPStylingType | - |



### UpdateIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateLabelPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| primary_color |  string | - |
| secondary_color |  string | - |
| hide_login_name_suffix |  bool | - |



### UpdateLabelPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateLoginPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| allow_username_password |  bool | - |
| allow_register |  bool | - |
| allow_external_idp |  bool | - |
| force_mfa |  bool | - |
| passwordless_type |  zitadel.policy.v1.PasswordlessType | - |



### UpdateLoginPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateOrgIAMPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_login_must_be_domain |  bool | - |



### UpdateOrgIAMPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdatePasswordAgePolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| max_age_days |  uint32 | - |
| expire_warn_days |  uint32 | - |



### UpdatePasswordAgePolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdatePasswordComplexityPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| min_length |  uint32 | - |
| has_uppercase |  bool | - |
| has_lowercase |  bool | - |
| has_number |  bool | - |
| has_symbol |  bool | - |



### UpdatePasswordComplexityPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdatePasswordLockoutPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| max_attempts |  uint32 | failed attempts until a user gets locked |
| show_lockout_failure |  bool | If an error should be displayed during a lockout or not |



### UpdatePasswordLockoutPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### View


| Field | Type | Description |
| ----- | ---- | ----------- |
| database |  string | - |
| view_name |  string | - |
| processed_sequence |  uint64 | - |
| event_timestamp |  google.protobuf.Timestamp | The timestamp the event occured |
| last_successful_spooler_run |  google.protobuf.Timestamp | - |






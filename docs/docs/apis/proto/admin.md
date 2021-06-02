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




### IsOrgUnique

> **rpc** IsOrgUnique([IsOrgUniqueRequest](#isorguniquerequest))
[IsOrgUniqueResponse](#isorguniqueresponse)

Checks whether an organisation exists by the given parameters




### GetOrgByID

> **rpc** GetOrgByID([GetOrgByIDRequest](#getorgbyidrequest))
[GetOrgByIDResponse](#getorgbyidresponse)

Returns an organisation by id




### ListOrgs

> **rpc** ListOrgs([ListOrgsRequest](#listorgsrequest))
[ListOrgsResponse](#listorgsresponse)

Returns all organisations matching the request
all queries need to match (AND)




### SetUpOrg

> **rpc** SetUpOrg([SetUpOrgRequest](#setuporgrequest))
[SetUpOrgResponse](#setuporgresponse)

Creates a new org and user 
and adds the user to the orgs members as ORG_OWNER




### GetIDPByID

> **rpc** GetIDPByID([GetIDPByIDRequest](#getidpbyidrequest))
[GetIDPByIDResponse](#getidpbyidresponse)

Returns a identity provider configuration of the IAM




### ListIDPs

> **rpc** ListIDPs([ListIDPsRequest](#listidpsrequest))
[ListIDPsResponse](#listidpsresponse)

Returns all identity provider configurations of the IAM




### AddOIDCIDP

> **rpc** AddOIDCIDP([AddOIDCIDPRequest](#addoidcidprequest))
[AddOIDCIDPResponse](#addoidcidpresponse)

Adds a new oidc identity provider configuration the IAM




### UpdateIDP

> **rpc** UpdateIDP([UpdateIDPRequest](#updateidprequest))
[UpdateIDPResponse](#updateidpresponse)

Updates the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.




### DeactivateIDP

> **rpc** DeactivateIDP([DeactivateIDPRequest](#deactivateidprequest))
[DeactivateIDPResponse](#deactivateidpresponse)

Sets the state of the idp to IDP_STATE_INACTIVE
the state MUST be IDP_STATE_ACTIVE for this call




### ReactivateIDP

> **rpc** ReactivateIDP([ReactivateIDPRequest](#reactivateidprequest))
[ReactivateIDPResponse](#reactivateidpresponse)

Sets the state of the idp to IDP_STATE_ACTIVE
the state MUST be IDP_STATE_INACTIVE for this call




### RemoveIDP

> **rpc** RemoveIDP([RemoveIDPRequest](#removeidprequest))
[RemoveIDPResponse](#removeidpresponse)

RemoveIDP deletes the IDP permanetly




### UpdateIDPOIDCConfig

> **rpc** UpdateIDPOIDCConfig([UpdateIDPOIDCConfigRequest](#updateidpoidcconfigrequest))
[UpdateIDPOIDCConfigResponse](#updateidpoidcconfigresponse)

Updates the oidc configuration of the specified idp
all fields are updated. If no value is provided the field will be empty afterwards.




### GetDefaultFeatures

> **rpc** GetDefaultFeatures([GetDefaultFeaturesRequest](#getdefaultfeaturesrequest))
[GetDefaultFeaturesResponse](#getdefaultfeaturesresponse)






### SetDefaultFeatures

> **rpc** SetDefaultFeatures([SetDefaultFeaturesRequest](#setdefaultfeaturesrequest))
[SetDefaultFeaturesResponse](#setdefaultfeaturesresponse)






### GetOrgFeatures

> **rpc** GetOrgFeatures([GetOrgFeaturesRequest](#getorgfeaturesrequest))
[GetOrgFeaturesResponse](#getorgfeaturesresponse)






### SetOrgFeatures

> **rpc** SetOrgFeatures([SetOrgFeaturesRequest](#setorgfeaturesrequest))
[SetOrgFeaturesResponse](#setorgfeaturesresponse)






### ResetOrgFeatures

> **rpc** ResetOrgFeatures([ResetOrgFeaturesRequest](#resetorgfeaturesrequest))
[ResetOrgFeaturesResponse](#resetorgfeaturesresponse)






### GetOrgIAMPolicy

> **rpc** GetOrgIAMPolicy([GetOrgIAMPolicyRequest](#getorgiampolicyrequest))
[GetOrgIAMPolicyResponse](#getorgiampolicyresponse)

Returns the IAM policy defined by the administrators of ZITADEL




### UpdateOrgIAMPolicy

> **rpc** UpdateOrgIAMPolicy([UpdateOrgIAMPolicyRequest](#updateorgiampolicyrequest))
[UpdateOrgIAMPolicyResponse](#updateorgiampolicyresponse)

Updates the default IAM policy.
it impacts all organisations without a customised policy




### GetCustomOrgIAMPolicy

> **rpc** GetCustomOrgIAMPolicy([GetCustomOrgIAMPolicyRequest](#getcustomorgiampolicyrequest))
[GetCustomOrgIAMPolicyResponse](#getcustomorgiampolicyresponse)

Returns the customised policy or the default if not customised




### AddCustomOrgIAMPolicy

> **rpc** AddCustomOrgIAMPolicy([AddCustomOrgIAMPolicyRequest](#addcustomorgiampolicyrequest))
[AddCustomOrgIAMPolicyResponse](#addcustomorgiampolicyresponse)

Defines a custom ORGIAM policy as specified




### UpdateCustomOrgIAMPolicy

> **rpc** UpdateCustomOrgIAMPolicy([UpdateCustomOrgIAMPolicyRequest](#updatecustomorgiampolicyrequest))
[UpdateCustomOrgIAMPolicyResponse](#updatecustomorgiampolicyresponse)

Updates a custom ORGIAM policy as specified




### ResetCustomOrgIAMPolicyToDefault

> **rpc** ResetCustomOrgIAMPolicyToDefault([ResetCustomOrgIAMPolicyToDefaultRequest](#resetcustomorgiampolicytodefaultrequest))
[ResetCustomOrgIAMPolicyToDefaultResponse](#resetcustomorgiampolicytodefaultresponse)

Resets the org iam policy of the organisation to default
ZITADEL will fallback to the default policy defined by the ZITADEL administrators




### GetLabelPolicy

> **rpc** GetLabelPolicy([GetLabelPolicyRequest](#getlabelpolicyrequest))
[GetLabelPolicyResponse](#getlabelpolicyresponse)

Returns the label policy defined by the administrators of ZITADEL




### UpdateLabelPolicy

> **rpc** UpdateLabelPolicy([UpdateLabelPolicyRequest](#updatelabelpolicyrequest))
[UpdateLabelPolicyResponse](#updatelabelpolicyresponse)

Updates the default label policy of ZITADEL
it impacts all organisations without a customised policy




### GetLoginPolicy

> **rpc** GetLoginPolicy([GetLoginPolicyRequest](#getloginpolicyrequest))
[GetLoginPolicyResponse](#getloginpolicyresponse)

Returns the login policy defined by the administrators of ZITADEL




### UpdateLoginPolicy

> **rpc** UpdateLoginPolicy([UpdateLoginPolicyRequest](#updateloginpolicyrequest))
[UpdateLoginPolicyResponse](#updateloginpolicyresponse)

Updates the default login policy of ZITADEL
it impacts all organisations without a customised policy




### ListLoginPolicyIDPs

> **rpc** ListLoginPolicyIDPs([ListLoginPolicyIDPsRequest](#listloginpolicyidpsrequest))
[ListLoginPolicyIDPsResponse](#listloginpolicyidpsresponse)

Returns the idps linked to the default login policy,
defined by the administrators of ZITADEL




### AddIDPToLoginPolicy

> **rpc** AddIDPToLoginPolicy([AddIDPToLoginPolicyRequest](#addidptologinpolicyrequest))
[AddIDPToLoginPolicyResponse](#addidptologinpolicyresponse)

Adds the povided idp to the default login policy.
It impacts all organisations without a customised policy




### RemoveIDPFromLoginPolicy

> **rpc** RemoveIDPFromLoginPolicy([RemoveIDPFromLoginPolicyRequest](#removeidpfromloginpolicyrequest))
[RemoveIDPFromLoginPolicyResponse](#removeidpfromloginpolicyresponse)

Removes the povided idp from the default login policy.
It impacts all organisations without a customised policy




### ListLoginPolicySecondFactors

> **rpc** ListLoginPolicySecondFactors([ListLoginPolicySecondFactorsRequest](#listloginpolicysecondfactorsrequest))
[ListLoginPolicySecondFactorsResponse](#listloginpolicysecondfactorsresponse)

Returns the available second factors defined by the administrators of ZITADEL




### AddSecondFactorToLoginPolicy

> **rpc** AddSecondFactorToLoginPolicy([AddSecondFactorToLoginPolicyRequest](#addsecondfactortologinpolicyrequest))
[AddSecondFactorToLoginPolicyResponse](#addsecondfactortologinpolicyresponse)

Adds a second factor to the default login policy.
It impacts all organisations without a customised policy




### RemoveSecondFactorFromLoginPolicy

> **rpc** RemoveSecondFactorFromLoginPolicy([RemoveSecondFactorFromLoginPolicyRequest](#removesecondfactorfromloginpolicyrequest))
[RemoveSecondFactorFromLoginPolicyResponse](#removesecondfactorfromloginpolicyresponse)

Removes a second factor from the default login policy.
It impacts all organisations without a customised policy




### ListLoginPolicyMultiFactors

> **rpc** ListLoginPolicyMultiFactors([ListLoginPolicyMultiFactorsRequest](#listloginpolicymultifactorsrequest))
[ListLoginPolicyMultiFactorsResponse](#listloginpolicymultifactorsresponse)

Returns the available multi factors defined by the administrators of ZITADEL




### AddMultiFactorToLoginPolicy

> **rpc** AddMultiFactorToLoginPolicy([AddMultiFactorToLoginPolicyRequest](#addmultifactortologinpolicyrequest))
[AddMultiFactorToLoginPolicyResponse](#addmultifactortologinpolicyresponse)

Adds a multi factor to the default login policy.
It impacts all organisations without a customised policy




### RemoveMultiFactorFromLoginPolicy

> **rpc** RemoveMultiFactorFromLoginPolicy([RemoveMultiFactorFromLoginPolicyRequest](#removemultifactorfromloginpolicyrequest))
[RemoveMultiFactorFromLoginPolicyResponse](#removemultifactorfromloginpolicyresponse)

Removes a multi factor from the default login policy.
It impacts all organisations without a customised policy




### GetPasswordComplexityPolicy

> **rpc** GetPasswordComplexityPolicy([GetPasswordComplexityPolicyRequest](#getpasswordcomplexitypolicyrequest))
[GetPasswordComplexityPolicyResponse](#getpasswordcomplexitypolicyresponse)

Returns the password complexity policy defined by the administrators of ZITADEL




### UpdatePasswordComplexityPolicy

> **rpc** UpdatePasswordComplexityPolicy([UpdatePasswordComplexityPolicyRequest](#updatepasswordcomplexitypolicyrequest))
[UpdatePasswordComplexityPolicyResponse](#updatepasswordcomplexitypolicyresponse)

Updates the default password complexity policy of ZITADEL
it impacts all organisations without a customised policy




### GetPasswordAgePolicy

> **rpc** GetPasswordAgePolicy([GetPasswordAgePolicyRequest](#getpasswordagepolicyrequest))
[GetPasswordAgePolicyResponse](#getpasswordagepolicyresponse)

Returns the password age policy defined by the administrators of ZITADEL




### UpdatePasswordAgePolicy

> **rpc** UpdatePasswordAgePolicy([UpdatePasswordAgePolicyRequest](#updatepasswordagepolicyrequest))
[UpdatePasswordAgePolicyResponse](#updatepasswordagepolicyresponse)

Updates the default password age policy of ZITADEL
it impacts all organisations without a customised policy




### GetPasswordLockoutPolicy

> **rpc** GetPasswordLockoutPolicy([GetPasswordLockoutPolicyRequest](#getpasswordlockoutpolicyrequest))
[GetPasswordLockoutPolicyResponse](#getpasswordlockoutpolicyresponse)

Returns the password lockout policy defined by the administrators of ZITADEL




### UpdatePasswordLockoutPolicy

> **rpc** UpdatePasswordLockoutPolicy([UpdatePasswordLockoutPolicyRequest](#updatepasswordlockoutpolicyrequest))
[UpdatePasswordLockoutPolicyResponse](#updatepasswordlockoutpolicyresponse)

Updates the default password lockout policy of ZITADEL
it impacts all organisations without a customised policy




### GetDefaultInitMessageText

> **rpc** GetDefaultInitMessageText([GetDefaultInitMessageTextRequest](#getdefaultinitmessagetextrequest))
[GetDefaultInitMessageTextResponse](#getdefaultinitmessagetextresponse)

Returns the custom text for initial message




### SetDefaultInitMessageText

> **rpc** SetDefaultInitMessageText([SetDefaultInitMessageTextRequest](#setdefaultinitmessagetextrequest))
[SetDefaultInitMessageTextResponse](#setdefaultinitmessagetextresponse)

Sets the default custom text for initial message
it impacts all organisations without customized initial message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}




### GetDefaultPasswordResetMessageText

> **rpc** GetDefaultPasswordResetMessageText([GetDefaultPasswordResetMessageTextRequest](#getdefaultpasswordresetmessagetextrequest))
[GetDefaultPasswordResetMessageTextResponse](#getdefaultpasswordresetmessagetextresponse)

Returns the custom text for password reset message




### SetDefaultPasswordResetMessageText

> **rpc** SetDefaultPasswordResetMessageText([SetDefaultPasswordResetMessageTextRequest](#setdefaultpasswordresetmessagetextrequest))
[SetDefaultPasswordResetMessageTextResponse](#setdefaultpasswordresetmessagetextresponse)

Sets the default custom text for password reset message
it impacts all organisations without customized password reset message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}




### GetDefaultVerifyEmailMessageText

> **rpc** GetDefaultVerifyEmailMessageText([GetDefaultVerifyEmailMessageTextRequest](#getdefaultverifyemailmessagetextrequest))
[GetDefaultVerifyEmailMessageTextResponse](#getdefaultverifyemailmessagetextresponse)

Returns the custom text for verify email message




### SetDefaultVerifyEmailMessageText

> **rpc** SetDefaultVerifyEmailMessageText([SetDefaultVerifyEmailMessageTextRequest](#setdefaultverifyemailmessagetextrequest))
[SetDefaultVerifyEmailMessageTextResponse](#setdefaultverifyemailmessagetextresponse)

Sets the default custom text for verify email message
it impacts all organisations without customized verify email message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}




### GetDefaultVerifyPhoneMessageText

> **rpc** GetDefaultVerifyPhoneMessageText([GetDefaultVerifyPhoneMessageTextRequest](#getdefaultverifyphonemessagetextrequest))
[GetDefaultVerifyPhoneMessageTextResponse](#getdefaultverifyphonemessagetextresponse)

Returns the custom text for verify phone message




### SetDefaultVerifyPhoneMessageText

> **rpc** SetDefaultVerifyPhoneMessageText([SetDefaultVerifyPhoneMessageTextRequest](#setdefaultverifyphonemessagetextrequest))
[SetDefaultVerifyPhoneMessageTextResponse](#setdefaultverifyphonemessagetextresponse)

Sets the default custom text for verify phone message
it impacts all organisations without customized verify phone message text
The Following Variables can be used:
{{.Code}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}




### GetDefaultDomainClaimedMessageText

> **rpc** GetDefaultDomainClaimedMessageText([GetDefaultDomainClaimedMessageTextRequest](#getdefaultdomainclaimedmessagetextrequest))
[GetDefaultDomainClaimedMessageTextResponse](#getdefaultdomainclaimedmessagetextresponse)

Returns the custom text for domain claimed message




### SetDefaultDomainClaimedMessageText

> **rpc** SetDefaultDomainClaimedMessageText([SetDefaultDomainClaimedMessageTextRequest](#setdefaultdomainclaimedmessagetextrequest))
[SetDefaultDomainClaimedMessageTextResponse](#setdefaultdomainclaimedmessagetextresponse)

Sets the default custom text for domain claimed phone message
it impacts all organisations without customized verify phone message text
The Following Variables can be used:
{{.Domain}} {{.TempUsername}} {{.UserName}} {{.FirstName}} {{.LastName}} {{.NickName}} {{.DisplayName}} {{.LastEmail}} {{.VerifiedEmail}} {{.LastPhone}} {{.VerifiedPhone}} {{.PreferredLoginName}} {{.LoginNames}} {{.ChangeDate}}




### ListIAMMemberRoles

> **rpc** ListIAMMemberRoles([ListIAMMemberRolesRequest](#listiammemberrolesrequest))
[ListIAMMemberRolesResponse](#listiammemberrolesresponse)

Returns the IAM roles visible for the requested user




### ListIAMMembers

> **rpc** ListIAMMembers([ListIAMMembersRequest](#listiammembersrequest))
[ListIAMMembersResponse](#listiammembersresponse)

Returns all members matching the request
all queries need to match (ANDed)




### AddIAMMember

> **rpc** AddIAMMember([AddIAMMemberRequest](#addiammemberrequest))
[AddIAMMemberResponse](#addiammemberresponse)

Adds a user to the membership list of ZITADEL with the given roles
undefined roles will be dropped




### UpdateIAMMember

> **rpc** UpdateIAMMember([UpdateIAMMemberRequest](#updateiammemberrequest))
[UpdateIAMMemberResponse](#updateiammemberresponse)

Sets the given roles on a member.
The member has only roles provided by this call




### RemoveIAMMember

> **rpc** RemoveIAMMember([RemoveIAMMemberRequest](#removeiammemberrequest))
[RemoveIAMMemberResponse](#removeiammemberresponse)

Removes the user from the membership list of ZITADEL




### ListViews

> **rpc** ListViews([ListViewsRequest](#listviewsrequest))
[ListViewsResponse](#listviewsresponse)

Returns all stored read models of ZITADEL
views are used for search optimisation and optimise request latencies
they represent the delta of the event happend on the objects




### ClearView

> **rpc** ClearView([ClearViewRequest](#clearviewrequest))
[ClearViewResponse](#clearviewresponse)

Truncates the delta of the change stream
be carefull with this function because ZITADEL has to 
recompute the deltas after they got cleared. 
Search requests will return wrong results until all deltas are recomputed




### ListFailedEvents

> **rpc** ListFailedEvents([ListFailedEventsRequest](#listfailedeventsrequest))
[ListFailedEventsResponse](#listfailedeventsresponse)

Returns event descriptions which cannot be processed.
It's possible that some events need some retries. 
For example if the SMTP-API wasn't able to send an email at the first time




### RemoveFailedEvent

> **rpc** RemoveFailedEvent([RemoveFailedEventRequest](#removefailedeventrequest))
[RemoveFailedEventResponse](#removefailedeventresponse)

Deletes the event from failed events view.
the event is not removed from the change stream
This call is usefull if the system was able to process the event later. 
e.g. if the second try of sending an email was successful. the first try produced a
failed event. You can find out if it worked on the `failure_count`









## Messages


### AddCustomOrgIAMPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_login_must_be_domain |  bool | the username has to end with the domain of it's organisation (uniqueness is organisation based) |  |




### AddCustomOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddIAMMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | if no roles provided the user won't have any rights |  |




### AddIAMMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddIDPToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | Id of the predefined idp configuration | string.min_len: 1<br /> string.max_len: 200<br />  |




### AddIDPToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddMultiFactorToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### AddMultiFactorToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddOIDCIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| styling_type |  zitadel.idp.v1.IDPStylingType | - | enum.defined_only: true<br />  |
| client_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_secret |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| scopes | repeated string | - |  |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |




### AddOIDCIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| idp_id |  string | - |  |




### AddSecondFactorToLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### AddSecondFactorToLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ClearViewRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| view_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ClearViewResponse
This is an empty response




### DeactivateIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### DeactivateIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### FailedEvent



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - |  |
| view_name |  string | - |  |
| failed_sequence |  uint64 | - |  |
| failure_count |  uint64 | - |  |
| error_message |  string | - |  |




### GetCustomOrgIAMPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetCustomOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.OrgIAMPolicy | - |  |
| is_default |  bool | - |  |




### GetDefaultDomainClaimedMessageTextRequest
This is an empty request


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultDomainClaimedMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultFeaturesRequest





### GetDefaultFeaturesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| features |  zitadel.features.v1.Features | - |  |




### GetDefaultInitMessageTextRequest
This is an empty request


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultInitMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultPasswordResetMessageTextRequest
This is an empty request


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultPasswordResetMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultVerifyEmailMessageTextRequest
This is an empty request


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultVerifyEmailMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetDefaultVerifyPhoneMessageTextRequest
This is an empty request


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetDefaultVerifyPhoneMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| custom_text |  zitadel.text.v1.MessageCustomText | - |  |




### GetIDPByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetIDPByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp |  zitadel.idp.v1.IDP | - |  |




### GetLabelPolicyRequest
This is an empty request




### GetLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LabelPolicy | - |  |




### GetLoginPolicyRequest
This is an empty request




### GetLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.LoginPolicy | - |  |




### GetOrgByIDRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetOrgByIDResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org |  zitadel.org.v1.Org | - |  |




### GetOrgFeaturesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetOrgFeaturesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| features |  zitadel.features.v1.Features | - |  |




### GetOrgIAMPolicyRequest





### GetOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.OrgIAMPolicy | - |  |




### GetPasswordAgePolicyRequest
This is an empty request




### GetPasswordAgePolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordAgePolicy | - |  |




### GetPasswordComplexityPolicyRequest





### GetPasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |  |




### GetPasswordLockoutPolicyRequest
This is an empty request




### GetPasswordLockoutPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordLockoutPolicy | - |  |




### HealthzRequest
This is an empty request




### HealthzResponse
This is an empty response




### IDPQuery



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_id_query |  zitadel.idp.v1.IDPIDQuery | - |  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) query.idp_name_query |  zitadel.idp.v1.IDPNameQuery | - |  |




### IsOrgUniqueRequest
if name or domain is already in use, org is not unique


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| domain |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### IsOrgUniqueResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| is_unique |  bool | - |  |




### ListFailedEventsRequest
This is an empty request




### ListFailedEventsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated FailedEvent | TODO: list details |  |




### ListIAMMemberRolesRequest
This is an empty request




### ListIAMMemberRolesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| roles | repeated string | - |  |




### ListIAMMembersRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.member.v1.SearchQuery | criterias the client is looking for |  |




### ListIAMMembersResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.member.v1.Member | - |  |




### ListIDPsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| sorting_column |  zitadel.idp.v1.IDPFieldName | the field the result is sorted |  |
| queries | repeated IDPQuery | criterias the client is looking for |  |




### ListIDPsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.idp.v1.IDPFieldName | - |  |
| result | repeated zitadel.idp.v1.IDP | - |  |




### ListLoginPolicyIDPsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |




### ListLoginPolicyIDPsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.idp.v1.IDPLoginPolicyLink | - |  |




### ListLoginPolicyMultiFactorsRequest
This is an empty request




### ListLoginPolicyMultiFactorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.policy.v1.MultiFactorType | - |  |




### ListLoginPolicySecondFactorsRequest
This is an empty request




### ListLoginPolicySecondFactorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.policy.v1.SecondFactorType | - |  |




### ListOrgsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| sorting_column |  zitadel.org.v1.OrgFieldName | the field the result is sorted |  |
| queries | repeated zitadel.org.v1.OrgQuery | criterias the client is looking for |  |




### ListOrgsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| sorting_column |  zitadel.org.v1.OrgFieldName | - |  |
| result | repeated zitadel.org.v1.Org | - |  |




### ListViewsRequest
This is an empty request




### ListViewsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated View | TODO: list details |  |




### ReactivateIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ReactivateIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveFailedEventRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| view_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| failed_sequence |  uint64 | - |  |




### RemoveFailedEventResponse
This is an empty response




### RemoveIAMMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveIAMMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveIDPFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveIDPFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveMultiFactorFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.MultiFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### RemoveMultiFactorFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RemoveSecondFactorFromLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| type |  zitadel.policy.v1.SecondFactorType | - | enum.defined_only: true<br /> enum.not_in: [0]<br />  |




### RemoveSecondFactorFromLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetCustomOrgIAMPolicyToDefaultRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetCustomOrgIAMPolicyToDefaultResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResetOrgFeaturesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### ResetOrgFeaturesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultDomainClaimedMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultDomainClaimedMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultFeaturesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| tier_name |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 200<br />  |
| audit_log_retention |  google.protobuf.Duration | - | duration.gte.seconds: 0<br /> duration.gte.nanos: 0<br />  |
| login_policy_username_login |  bool | - |  |
| login_policy_registration |  bool | - |  |
| login_policy_idp |  bool | - |  |
| login_policy_factors |  bool | - |  |
| login_policy_passwordless |  bool | - |  |
| password_complexity_policy |  bool | - |  |
| label_policy |  bool | - |  |
| custom_domain |  bool | - |  |
| custom_text |  bool | - |  |




### SetDefaultFeaturesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultInitMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 1000<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultInitMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultPasswordResetMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultPasswordResetMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultVerifyEmailMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultVerifyEmailMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetDefaultVerifyPhoneMessageTextRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| language |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| title |  string | - | string.max_len: 200<br />  |
| pre_header |  string | - | string.max_len: 200<br />  |
| subject |  string | - | string.max_len: 200<br />  |
| greeting |  string | - | string.max_len: 200<br />  |
| text |  string | - | string.max_len: 800<br />  |
| button_text |  string | - | string.max_len: 200<br />  |
| footer_text |  string | - | string.max_len: 200<br />  |




### SetDefaultVerifyPhoneMessageTextResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetOrgFeaturesRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| tier_name |  string | - | string.max_len: 200<br />  |
| description |  string | - | string.max_len: 200<br />  |
| state |  zitadel.features.v1.FeaturesState | - |  |
| state_description |  string | - | string.max_len: 200<br />  |
| audit_log_retention |  google.protobuf.Duration | - | duration.gte.seconds: 0<br /> duration.gte.nanos: 0<br />  |
| login_policy_username_login |  bool | - |  |
| login_policy_registration |  bool | - |  |
| login_policy_idp |  bool | - |  |
| login_policy_factors |  bool | - |  |
| login_policy_passwordless |  bool | - |  |
| password_complexity_policy |  bool | - |  |
| label_policy |  bool | - |  |
| custom_domain |  bool | - |  |
| custom_text |  bool | - |  |




### SetOrgFeaturesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SetUpOrgRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org |  SetUpOrgRequest.Org | - | message.required: true<br />  |
| [**oneof**](https://developers.google.com/protocol-buffers/docs/proto3#oneof) user.human |  SetUpOrgRequest.Human | oneof field for the user managing the organisation |  |




### SetUpOrgRequest.Human



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| profile |  SetUpOrgRequest.Human.Profile | - | message.required: true<br />  |
| email |  SetUpOrgRequest.Human.Email | - | message.required: true<br />  |
| phone |  SetUpOrgRequest.Human.Phone | - |  |
| password |  string | - |  |




### SetUpOrgRequest.Human.Email



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | TODO: check if no value is allowed | string.email: true<br />  |
| is_email_verified |  bool | - |  |




### SetUpOrgRequest.Human.Phone



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| phone |  string | has to be a global number | string.min_len: 1<br /> string.max_len: 50<br /> string.prefix: +<br />  |
| is_phone_verified |  bool | - |  |




### SetUpOrgRequest.Human.Profile



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| last_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| nick_name |  string | - | string.max_len: 200<br />  |
| display_name |  string | - | string.max_len: 200<br />  |
| preferred_language |  string | - | string.max_len: 10<br />  |
| gender |  zitadel.user.v1.Gender | - |  |




### SetUpOrgRequest.Org



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| domain |  string | - | string.max_len: 200<br />  |




### SetUpOrgResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| org_id |  string | - |  |
| user_id |  string | - |  |




### UpdateCustomOrgIAMPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| user_login_must_be_domain |  bool | - |  |




### UpdateCustomOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateIAMMemberRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| roles | repeated string | if no roles provided the user won't have any rights |  |




### UpdateIAMMemberResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateIDPOIDCConfigRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| issuer |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| client_secret |  string | - | string.max_len: 200<br />  |
| scopes | repeated string | - |  |
| display_name_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |
| username_mapping |  zitadel.idp.v1.OIDCMappingField | - | enum.defined_only: true<br />  |




### UpdateIDPOIDCConfigResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateIDPRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| styling_type |  zitadel.idp.v1.IDPStylingType | - | enum.defined_only: true<br />  |




### UpdateIDPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateLabelPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| primary_color |  string | - | string.min_len: 1<br /> string.max_len: 50<br />  |
| secondary_color |  string | - | string.min_len: 1<br /> string.max_len: 50<br />  |
| hide_login_name_suffix |  bool | - |  |




### UpdateLabelPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateLoginPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| allow_username_password |  bool | - |  |
| allow_register |  bool | - |  |
| allow_external_idp |  bool | - |  |
| force_mfa |  bool | - |  |
| passwordless_type |  zitadel.policy.v1.PasswordlessType | - | enum.defined_only: true<br />  |




### UpdateLoginPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdateOrgIAMPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_login_must_be_domain |  bool | - |  |




### UpdateOrgIAMPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdatePasswordAgePolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| max_age_days |  uint32 | - |  |
| expire_warn_days |  uint32 | - |  |




### UpdatePasswordAgePolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdatePasswordComplexityPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| min_length |  uint32 | - |  |
| has_uppercase |  bool | - |  |
| has_lowercase |  bool | - |  |
| has_number |  bool | - |  |
| has_symbol |  bool | - |  |




### UpdatePasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### UpdatePasswordLockoutPolicyRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| max_attempts |  uint32 | failed attempts until a user gets locked |  |
| show_lockout_failure |  bool | If an error should be displayed during a lockout or not |  |




### UpdatePasswordLockoutPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### View



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| database |  string | - |  |
| view_name |  string | - |  |
| processed_sequence |  uint64 | - |  |
| event_timestamp |  google.protobuf.Timestamp | The timestamp the event occured |  |
| last_successful_spooler_run |  google.protobuf.Timestamp | - |  |







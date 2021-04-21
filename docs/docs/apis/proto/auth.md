---
title: zitadel/auth.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)


## AuthService {#zitadelauthv1authservice}


### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)






### GetMyUser

> **rpc** GetMyUser([GetMyUserRequest](#getmyuserrequest))
[GetMyUserResponse](#getmyuserresponse)






### ListMyUserChanges

> **rpc** ListMyUserChanges([ListMyUserChangesRequest](#listmyuserchangesrequest))
[ListMyUserChangesResponse](#listmyuserchangesresponse)






### ListMyUserSessions

> **rpc** ListMyUserSessions([ListMyUserSessionsRequest](#listmyusersessionsrequest))
[ListMyUserSessionsResponse](#listmyusersessionsresponse)






### UpdateMyUserName

> **rpc** UpdateMyUserName([UpdateMyUserNameRequest](#updatemyusernamerequest))
[UpdateMyUserNameResponse](#updatemyusernameresponse)






### GetMyPasswordComplexityPolicy

> **rpc** GetMyPasswordComplexityPolicy([GetMyPasswordComplexityPolicyRequest](#getmypasswordcomplexitypolicyrequest))
[GetMyPasswordComplexityPolicyResponse](#getmypasswordcomplexitypolicyresponse)






### UpdateMyPassword

> **rpc** UpdateMyPassword([UpdateMyPasswordRequest](#updatemypasswordrequest))
[UpdateMyPasswordResponse](#updatemypasswordresponse)






### GetMyProfile

> **rpc** GetMyProfile([GetMyProfileRequest](#getmyprofilerequest))
[GetMyProfileResponse](#getmyprofileresponse)






### UpdateMyProfile

> **rpc** UpdateMyProfile([UpdateMyProfileRequest](#updatemyprofilerequest))
[UpdateMyProfileResponse](#updatemyprofileresponse)






### GetMyEmail

> **rpc** GetMyEmail([GetMyEmailRequest](#getmyemailrequest))
[GetMyEmailResponse](#getmyemailresponse)






### SetMyEmail

> **rpc** SetMyEmail([SetMyEmailRequest](#setmyemailrequest))
[SetMyEmailResponse](#setmyemailresponse)






### VerifyMyEmail

> **rpc** VerifyMyEmail([VerifyMyEmailRequest](#verifymyemailrequest))
[VerifyMyEmailResponse](#verifymyemailresponse)






### ResendMyEmailVerification

> **rpc** ResendMyEmailVerification([ResendMyEmailVerificationRequest](#resendmyemailverificationrequest))
[ResendMyEmailVerificationResponse](#resendmyemailverificationresponse)






### GetMyPhone

> **rpc** GetMyPhone([GetMyPhoneRequest](#getmyphonerequest))
[GetMyPhoneResponse](#getmyphoneresponse)






### SetMyPhone

> **rpc** SetMyPhone([SetMyPhoneRequest](#setmyphonerequest))
[SetMyPhoneResponse](#setmyphoneresponse)






### VerifyMyPhone

> **rpc** VerifyMyPhone([VerifyMyPhoneRequest](#verifymyphonerequest))
[VerifyMyPhoneResponse](#verifymyphoneresponse)






### ResendMyPhoneVerification

> **rpc** ResendMyPhoneVerification([ResendMyPhoneVerificationRequest](#resendmyphoneverificationrequest))
[ResendMyPhoneVerificationResponse](#resendmyphoneverificationresponse)






### RemoveMyPhone

> **rpc** RemoveMyPhone([RemoveMyPhoneRequest](#removemyphonerequest))
[RemoveMyPhoneResponse](#removemyphoneresponse)






### ListMyLinkedIDPs

> **rpc** ListMyLinkedIDPs([ListMyLinkedIDPsRequest](#listmylinkedidpsrequest))
[ListMyLinkedIDPsResponse](#listmylinkedidpsresponse)






### RemoveMyLinkedIDP

> **rpc** RemoveMyLinkedIDP([RemoveMyLinkedIDPRequest](#removemylinkedidprequest))
[RemoveMyLinkedIDPResponse](#removemylinkedidpresponse)






### ListMyAuthFactors

> **rpc** ListMyAuthFactors([ListMyAuthFactorsRequest](#listmyauthfactorsrequest))
[ListMyAuthFactorsResponse](#listmyauthfactorsresponse)






### AddMyAuthFactorOTP

> **rpc** AddMyAuthFactorOTP([AddMyAuthFactorOTPRequest](#addmyauthfactorotprequest))
[AddMyAuthFactorOTPResponse](#addmyauthfactorotpresponse)






### VerifyMyAuthFactorOTP

> **rpc** VerifyMyAuthFactorOTP([VerifyMyAuthFactorOTPRequest](#verifymyauthfactorotprequest))
[VerifyMyAuthFactorOTPResponse](#verifymyauthfactorotpresponse)






### RemoveMyAuthFactorOTP

> **rpc** RemoveMyAuthFactorOTP([RemoveMyAuthFactorOTPRequest](#removemyauthfactorotprequest))
[RemoveMyAuthFactorOTPResponse](#removemyauthfactorotpresponse)






### AddMyAuthFactorU2F

> **rpc** AddMyAuthFactorU2F([AddMyAuthFactorU2FRequest](#addmyauthfactoru2frequest))
[AddMyAuthFactorU2FResponse](#addmyauthfactoru2fresponse)






### VerifyMyAuthFactorU2F

> **rpc** VerifyMyAuthFactorU2F([VerifyMyAuthFactorU2FRequest](#verifymyauthfactoru2frequest))
[VerifyMyAuthFactorU2FResponse](#verifymyauthfactoru2fresponse)






### RemoveMyAuthFactorU2F

> **rpc** RemoveMyAuthFactorU2F([RemoveMyAuthFactorU2FRequest](#removemyauthfactoru2frequest))
[RemoveMyAuthFactorU2FResponse](#removemyauthfactoru2fresponse)






### ListMyPasswordless

> **rpc** ListMyPasswordless([ListMyPasswordlessRequest](#listmypasswordlessrequest))
[ListMyPasswordlessResponse](#listmypasswordlessresponse)






### AddMyPasswordless

> **rpc** AddMyPasswordless([AddMyPasswordlessRequest](#addmypasswordlessrequest))
[AddMyPasswordlessResponse](#addmypasswordlessresponse)






### VerifyMyPasswordless

> **rpc** VerifyMyPasswordless([VerifyMyPasswordlessRequest](#verifymypasswordlessrequest))
[VerifyMyPasswordlessResponse](#verifymypasswordlessresponse)






### RemoveMyPasswordless

> **rpc** RemoveMyPasswordless([RemoveMyPasswordlessRequest](#removemypasswordlessrequest))
[RemoveMyPasswordlessResponse](#removemypasswordlessresponse)






### ListMyUserGrants

> **rpc** ListMyUserGrants([ListMyUserGrantsRequest](#listmyusergrantsrequest))
[ListMyUserGrantsResponse](#listmyusergrantsresponse)






### ListMyProjectOrgs

> **rpc** ListMyProjectOrgs([ListMyProjectOrgsRequest](#listmyprojectorgsrequest))
[ListMyProjectOrgsResponse](#listmyprojectorgsresponse)






### ListMyZitadelFeatures

> **rpc** ListMyZitadelFeatures([ListMyZitadelFeaturesRequest](#listmyzitadelfeaturesrequest))
[ListMyZitadelFeaturesResponse](#listmyzitadelfeaturesresponse)






### ListMyZitadelPermissions

> **rpc** ListMyZitadelPermissions([ListMyZitadelPermissionsRequest](#listmyzitadelpermissionsrequest))
[ListMyZitadelPermissionsResponse](#listmyzitadelpermissionsresponse)






### ListMyProjectPermissions

> **rpc** ListMyProjectPermissions([ListMyProjectPermissionsRequest](#listmyprojectpermissionsrequest))
[ListMyProjectPermissionsResponse](#listmyprojectpermissionsresponse)











## Messages


### AddMyAuthFactorOTPRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### AddMyAuthFactorOTPResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| url |  string | - |  |
| secret |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |



### AddMyAuthFactorU2FRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### AddMyAuthFactorU2FResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  zitadel.user.v1.WebAuthNKey | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |



### AddMyPasswordlessRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### AddMyPasswordlessResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  zitadel.user.v1.WebAuthNKey | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |



### GetMyEmailRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### GetMyEmailResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| email |  zitadel.user.v1.Email | - |  |



### GetMyPasswordComplexityPolicyRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### GetMyPasswordComplexityPolicyResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |  |



### GetMyPhoneRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### GetMyPhoneResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| phone |  zitadel.user.v1.Phone | - |  |



### GetMyProfileRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### GetMyProfileResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| profile |  zitadel.user.v1.Profile | - |  |



### GetMyUserRequest
GetMyUserRequest is an empty request
the request parameters are read from the token-header

| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### GetMyUserResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user |  zitadel.user.v1.User | - |  |
| last_login |  google.protobuf.Timestamp | - |  |



### HealthzRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### HealthzResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ListMyAuthFactorsRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ListMyAuthFactorsResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated zitadel.user.v1.AuthFactor | - |  |



### ListMyLinkedIDPsRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering

PLANNED: queries for idp name and login name |  |



### ListMyLinkedIDPsResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.idp.v1.IDPUserLink | - |  |



### ListMyPasswordlessRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ListMyPasswordlessResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated zitadel.user.v1.WebAuthNToken | - |  |



### ListMyProjectOrgsRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |
| queries | repeated zitadel.org.v1.OrgQuery | criterias the client is looking for |  |



### ListMyProjectOrgsResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.org.v1.Org | - |  |



### ListMyProjectPermissionsRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ListMyProjectPermissionsResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated string | - |  |



### ListMyUserChangesRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | - |  |



### ListMyUserChangesResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.change.v1.Change | - |  |



### ListMyUserGrantsRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |



### ListMyUserGrantsResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated UserGrant | - |  |



### ListMyUserSessionsRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ListMyUserSessionsResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated zitadel.user.v1.Session | - |  |



### ListMyZitadelFeaturesRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ListMyZitadelFeaturesResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated string | - |  |



### ListMyZitadelPermissionsRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ListMyZitadelPermissionsResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated string | - |  |



### RemoveMyAuthFactorOTPRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### RemoveMyAuthFactorOTPResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### RemoveMyAuthFactorU2FRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| token_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |



### RemoveMyAuthFactorU2FResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### RemoveMyLinkedIDPRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| idp_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| linked_user_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |



### RemoveMyLinkedIDPResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### RemoveMyPasswordlessRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| token_id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |



### RemoveMyPasswordlessResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### RemoveMyPhoneRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### RemoveMyPhoneResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### ResendMyEmailVerificationRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ResendMyEmailVerificationResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### ResendMyPhoneVerificationRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |



### ResendMyPhoneVerificationResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### SetMyEmailRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| email |  string | TODO: check if no value is allowed | string.email: true<br />  |



### SetMyEmailResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### SetMyPhoneRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| phone |  string | - | string.min_len: 1<br /> string.max_len: 50<br /> string.prefix: +<br />  |



### SetMyPhoneResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### UpdateMyPasswordRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| old_password |  string | - | string.min_len: 1<br /> string.max_bytes: 70<br />  |
| new_password |  string | - | string.min_len: 1<br /> string.max_bytes: 70<br />  |



### UpdateMyPasswordResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### UpdateMyProfileRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| first_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| last_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| nick_name |  string | - | string.max_len: 200<br />  |
| display_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| preferred_language |  string | - | string.max_len: 10<br />  |
| gender |  zitadel.user.v1.Gender | - |  |



### UpdateMyProfileResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### UpdateMyUserNameRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user_name |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |



### UpdateMyUserNameResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### UserGrant


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| org_id |  string | - |  |
| project_id |  string | - |  |
| user_id |  string | - |  |
| roles | repeated string | - |  |
| org_name |  string | - |  |
| grant_id |  string | - |  |



### VerifyMyAuthFactorOTPRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| code |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |



### VerifyMyAuthFactorOTPResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### VerifyMyAuthFactorU2FRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| verification |  zitadel.user.v1.WebAuthNVerification | - | message.required: true<br />  |



### VerifyMyAuthFactorU2FResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### VerifyMyEmailRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| code |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |



### VerifyMyEmailResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### VerifyMyPasswordlessRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| verification |  zitadel.user.v1.WebAuthNVerification | - | message.required: true<br />  |



### VerifyMyPasswordlessResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |



### VerifyMyPhoneRequest


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| code |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |



### VerifyMyPhoneResponse


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |






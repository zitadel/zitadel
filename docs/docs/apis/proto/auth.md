---
title: zitadel/auth.proto
---


## AuthService {#zitadelauthv1authservice}


### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)





    GET: /healthz


### GetMyUser

> **rpc** GetMyUser([GetMyUserRequest](#getmyuserrequest))
[GetMyUserResponse](#getmyuserresponse)





    GET: /users/me


### ListMyUserChanges

> **rpc** ListMyUserChanges([ListMyUserChangesRequest](#listmyuserchangesrequest))
[ListMyUserChangesResponse](#listmyuserchangesresponse)





    POST: /users/me/changes/_search


### ListMyUserSessions

> **rpc** ListMyUserSessions([ListMyUserSessionsRequest](#listmyusersessionsrequest))
[ListMyUserSessionsResponse](#listmyusersessionsresponse)





    POST: /users/me/sessions/_search


### UpdateMyUserName

> **rpc** UpdateMyUserName([UpdateMyUserNameRequest](#updatemyusernamerequest))
[UpdateMyUserNameResponse](#updatemyusernameresponse)





    PUT: /users/me/username


### GetMyPasswordComplexityPolicy

> **rpc** GetMyPasswordComplexityPolicy([GetMyPasswordComplexityPolicyRequest](#getmypasswordcomplexitypolicyrequest))
[GetMyPasswordComplexityPolicyResponse](#getmypasswordcomplexitypolicyresponse)





    GET: /policies/passwords/complexity


### UpdateMyPassword

> **rpc** UpdateMyPassword([UpdateMyPasswordRequest](#updatemypasswordrequest))
[UpdateMyPasswordResponse](#updatemypasswordresponse)





    PUT: /users/me/password


### GetMyProfile

> **rpc** GetMyProfile([GetMyProfileRequest](#getmyprofilerequest))
[GetMyProfileResponse](#getmyprofileresponse)





    GET: /users/me/profile


### UpdateMyProfile

> **rpc** UpdateMyProfile([UpdateMyProfileRequest](#updatemyprofilerequest))
[UpdateMyProfileResponse](#updatemyprofileresponse)





    PUT: /users/me/profile


### GetMyEmail

> **rpc** GetMyEmail([GetMyEmailRequest](#getmyemailrequest))
[GetMyEmailResponse](#getmyemailresponse)





    GET: /users/me/email


### SetMyEmail

> **rpc** SetMyEmail([SetMyEmailRequest](#setmyemailrequest))
[SetMyEmailResponse](#setmyemailresponse)





    PUT: /users/me/email


### VerifyMyEmail

> **rpc** VerifyMyEmail([VerifyMyEmailRequest](#verifymyemailrequest))
[VerifyMyEmailResponse](#verifymyemailresponse)





    POST: /users/me/email/_verify


### ResendMyEmailVerification

> **rpc** ResendMyEmailVerification([ResendMyEmailVerificationRequest](#resendmyemailverificationrequest))
[ResendMyEmailVerificationResponse](#resendmyemailverificationresponse)





    POST: /users/me/email/_resend_verification


### GetMyPhone

> **rpc** GetMyPhone([GetMyPhoneRequest](#getmyphonerequest))
[GetMyPhoneResponse](#getmyphoneresponse)





    GET: /users/me/phone


### SetMyPhone

> **rpc** SetMyPhone([SetMyPhoneRequest](#setmyphonerequest))
[SetMyPhoneResponse](#setmyphoneresponse)





    PUT: /users/me/phone


### VerifyMyPhone

> **rpc** VerifyMyPhone([VerifyMyPhoneRequest](#verifymyphonerequest))
[VerifyMyPhoneResponse](#verifymyphoneresponse)





    POST: /users/me/phone/_verify


### ResendMyPhoneVerification

> **rpc** ResendMyPhoneVerification([ResendMyPhoneVerificationRequest](#resendmyphoneverificationrequest))
[ResendMyPhoneVerificationResponse](#resendmyphoneverificationresponse)





    POST: /users/me/phone/_resend_verification


### RemoveMyPhone

> **rpc** RemoveMyPhone([RemoveMyPhoneRequest](#removemyphonerequest))
[RemoveMyPhoneResponse](#removemyphoneresponse)





    DELETE: /users/me/phone


### ListMyLinkedIDPs

> **rpc** ListMyLinkedIDPs([ListMyLinkedIDPsRequest](#listmylinkedidpsrequest))
[ListMyLinkedIDPsResponse](#listmylinkedidpsresponse)





    POST: /users/me/idps/_search


### RemoveMyLinkedIDP

> **rpc** RemoveMyLinkedIDP([RemoveMyLinkedIDPRequest](#removemylinkedidprequest))
[RemoveMyLinkedIDPResponse](#removemylinkedidpresponse)





    DELETE: /users/me/idps/{idp_id}/{linked_user_id}


### ListMyAuthFactors

> **rpc** ListMyAuthFactors([ListMyAuthFactorsRequest](#listmyauthfactorsrequest))
[ListMyAuthFactorsResponse](#listmyauthfactorsresponse)





    POST: /users/me/auth_factors/_search


### AddMyAuthFactorOTP

> **rpc** AddMyAuthFactorOTP([AddMyAuthFactorOTPRequest](#addmyauthfactorotprequest))
[AddMyAuthFactorOTPResponse](#addmyauthfactorotpresponse)





    POST: /users/me/auth_factors/otp


### VerifyMyAuthFactorOTP

> **rpc** VerifyMyAuthFactorOTP([VerifyMyAuthFactorOTPRequest](#verifymyauthfactorotprequest))
[VerifyMyAuthFactorOTPResponse](#verifymyauthfactorotpresponse)





    POST: /users/me/auth_factors/otp/_verify


### RemoveMyAuthFactorOTP

> **rpc** RemoveMyAuthFactorOTP([RemoveMyAuthFactorOTPRequest](#removemyauthfactorotprequest))
[RemoveMyAuthFactorOTPResponse](#removemyauthfactorotpresponse)





    DELETE: /users/me/auth_factors/otp


### AddMyAuthFactorU2F

> **rpc** AddMyAuthFactorU2F([AddMyAuthFactorU2FRequest](#addmyauthfactoru2frequest))
[AddMyAuthFactorU2FResponse](#addmyauthfactoru2fresponse)





    POST: /users/me/auth_factors/u2f


### VerifyMyAuthFactorU2F

> **rpc** VerifyMyAuthFactorU2F([VerifyMyAuthFactorU2FRequest](#verifymyauthfactoru2frequest))
[VerifyMyAuthFactorU2FResponse](#verifymyauthfactoru2fresponse)





    POST: /users/me/auth_factors/u2f/_verify


### RemoveMyAuthFactorU2F

> **rpc** RemoveMyAuthFactorU2F([RemoveMyAuthFactorU2FRequest](#removemyauthfactoru2frequest))
[RemoveMyAuthFactorU2FResponse](#removemyauthfactoru2fresponse)





    DELETE: /users/me/auth_factors/u2f/{token_id}


### ListMyPasswordless

> **rpc** ListMyPasswordless([ListMyPasswordlessRequest](#listmypasswordlessrequest))
[ListMyPasswordlessResponse](#listmypasswordlessresponse)





    POST: /users/me/passwordless/_search


### AddMyPasswordless

> **rpc** AddMyPasswordless([AddMyPasswordlessRequest](#addmypasswordlessrequest))
[AddMyPasswordlessResponse](#addmypasswordlessresponse)





    POST: /users/me/passwordless


### VerifyMyPasswordless

> **rpc** VerifyMyPasswordless([VerifyMyPasswordlessRequest](#verifymypasswordlessrequest))
[VerifyMyPasswordlessResponse](#verifymypasswordlessresponse)





    POST: /users/me/passwordless/_verify


### RemoveMyPasswordless

> **rpc** RemoveMyPasswordless([RemoveMyPasswordlessRequest](#removemypasswordlessrequest))
[RemoveMyPasswordlessResponse](#removemypasswordlessresponse)





    DELETE: /users/me/passwordless/{token_id}


### ListMyUserGrants

> **rpc** ListMyUserGrants([ListMyUserGrantsRequest](#listmyusergrantsrequest))
[ListMyUserGrantsResponse](#listmyusergrantsresponse)





    POST: /usergrants/me/_search


### ListMyProjectOrgs

> **rpc** ListMyProjectOrgs([ListMyProjectOrgsRequest](#listmyprojectorgsrequest))
[ListMyProjectOrgsResponse](#listmyprojectorgsresponse)





    POST: /global/projectorgs/_search


### ListMyZitadelFeatures

> **rpc** ListMyZitadelFeatures([ListMyZitadelFeaturesRequest](#listmyzitadelfeaturesrequest))
[ListMyZitadelFeaturesResponse](#listmyzitadelfeaturesresponse)





    POST: /features/zitadel/me/_search


### ListMyZitadelPermissions

> **rpc** ListMyZitadelPermissions([ListMyZitadelPermissionsRequest](#listmyzitadelpermissionsrequest))
[ListMyZitadelPermissionsResponse](#listmyzitadelpermissionsresponse)





    POST: /permissions/zitadel/me/_search


### ListMyProjectPermissions

> **rpc** ListMyProjectPermissions([ListMyProjectPermissionsRequest](#listmyprojectpermissionsrequest))
[ListMyProjectPermissionsResponse](#listmyprojectpermissionsresponse)





    POST: /permissions/me/_search







## Messages


### AddMyAuthFactorOTPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### AddMyAuthFactorOTPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| url |  string | - |
| secret |  string | - |
| details |  zitadel.v1.ObjectDetails | - |



### AddMyAuthFactorU2FRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### AddMyAuthFactorU2FResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| key |  zitadel.user.v1.WebAuthNKey | - |
| details |  zitadel.v1.ObjectDetails | - |



### AddMyPasswordlessRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### AddMyPasswordlessResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| key |  zitadel.user.v1.WebAuthNKey | - |
| details |  zitadel.v1.ObjectDetails | - |



### GetMyEmailRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetMyEmailResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| email |  zitadel.user.v1.Email | - |



### GetMyPasswordComplexityPolicyRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetMyPasswordComplexityPolicyResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |



### GetMyPhoneRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetMyPhoneResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| phone |  zitadel.user.v1.Phone | - |



### GetMyProfileRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### GetMyProfileResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |
| profile |  zitadel.user.v1.Profile | - |



### GetMyUserRequest
GetMyUserRequest is an empty request
the request parameters are read from the token-header

| Field | Type | Description |
| ----- | ---- | ----------- |



### GetMyUserResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| user |  zitadel.user.v1.User | - |
| last_login |  google.protobuf.Timestamp | - |



### HealthzRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### HealthzResponse


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListMyAuthFactorsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListMyAuthFactorsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated zitadel.user.v1.AuthFactor | - |



### ListMyLinkedIDPsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering

PLANNED: queries for idp name and login name |



### ListMyLinkedIDPsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.idp.v1.IDPUserLink | - |



### ListMyPasswordlessRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListMyPasswordlessResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated zitadel.user.v1.WebAuthNToken | - |



### ListMyProjectOrgsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |
| queries | repeated zitadel.org.v1.OrgQuery | criterias the client is looking for |



### ListMyProjectOrgsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.org.v1.Org | - |



### ListMyProjectPermissionsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListMyProjectPermissionsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated string | - |



### ListMyUserChangesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.change.v1.ChangeQuery | - |



### ListMyUserChangesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated zitadel.change.v1.Change | - |



### ListMyUserGrantsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |



### ListMyUserGrantsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ListDetails | - |
| result | repeated UserGrant | - |



### ListMyUserSessionsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListMyUserSessionsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated zitadel.user.v1.Session | - |



### ListMyZitadelFeaturesRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListMyZitadelFeaturesResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated string | - |



### ListMyZitadelPermissionsRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ListMyZitadelPermissionsResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| result | repeated string | - |



### RemoveMyAuthFactorOTPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### RemoveMyAuthFactorOTPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveMyAuthFactorU2FRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| token_id |  string | - |



### RemoveMyAuthFactorU2FResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveMyLinkedIDPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| idp_id |  string | - |
| linked_user_id |  string | - |



### RemoveMyLinkedIDPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveMyPasswordlessRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| token_id |  string | - |



### RemoveMyPasswordlessResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### RemoveMyPhoneRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### RemoveMyPhoneResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResendMyEmailVerificationRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ResendMyEmailVerificationResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### ResendMyPhoneVerificationRequest


| Field | Type | Description |
| ----- | ---- | ----------- |



### ResendMyPhoneVerificationResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### SetMyEmailRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| email |  string | TODO: check if no value is allowed |



### SetMyEmailResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### SetMyPhoneRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| phone |  string | - |



### SetMyPhoneResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateMyPasswordRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| old_password |  string | - |
| new_password |  string | - |



### UpdateMyPasswordResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateMyProfileRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| first_name |  string | - |
| last_name |  string | - |
| nick_name |  string | - |
| display_name |  string | - |
| preferred_language |  string | - |
| gender |  zitadel.user.v1.Gender | - |



### UpdateMyProfileResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UpdateMyUserNameRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| user_name |  string | - |



### UpdateMyUserNameResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### UserGrant


| Field | Type | Description |
| ----- | ---- | ----------- |
| org_id |  string | - |
| project_id |  string | - |
| user_id |  string | - |
| roles | repeated string | - |
| org_name |  string | - |
| grant_id |  string | - |



### VerifyMyAuthFactorOTPRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| code |  string | - |



### VerifyMyAuthFactorOTPResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### VerifyMyAuthFactorU2FRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| verification |  zitadel.user.v1.WebAuthNVerification | - |



### VerifyMyAuthFactorU2FResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### VerifyMyEmailRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| code |  string | - |



### VerifyMyEmailResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### VerifyMyPasswordlessRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| verification |  zitadel.user.v1.WebAuthNVerification | - |



### VerifyMyPasswordlessResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |



### VerifyMyPhoneRequest


| Field | Type | Description |
| ----- | ---- | ----------- |
| code |  string | - |



### VerifyMyPhoneResponse


| Field | Type | Description |
| ----- | ---- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |






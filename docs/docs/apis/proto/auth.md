---
title: zitadel/auth.proto
---
> This document reflects the state from API 1.0 (available from 20.04.2021)


## AuthService {#zitadelauthv1authservice}


### Healthz

> **rpc** Healthz([HealthzRequest](#healthzrequest))
[HealthzResponse](#healthzresponse)





    GET: /healthz


### GetSupportedLanguages

> **rpc** GetSupportedLanguages([GetSupportedLanguagesRequest](#getsupportedlanguagesrequest))
[GetSupportedLanguagesResponse](#getsupportedlanguagesresponse)

Returns the default languages



    GET: /languages


### GetMyUser

> **rpc** GetMyUser([GetMyUserRequest](#getmyuserrequest))
[GetMyUserResponse](#getmyuserresponse)

Returns my full blown user



    GET: /users/me


### ListMyUserChanges

> **rpc** ListMyUserChanges([ListMyUserChangesRequest](#listmyuserchangesrequest))
[ListMyUserChangesResponse](#listmyuserchangesresponse)

Returns the history of the authorized user (each event)



    POST: /users/me/changes/_search


### ListMyUserSessions

> **rpc** ListMyUserSessions([ListMyUserSessionsRequest](#listmyusersessionsrequest))
[ListMyUserSessionsResponse](#listmyusersessionsresponse)

Returns the user sessions of the authorized user of the current useragent



    POST: /users/me/sessions/_search


### SetMyMetadata

> **rpc** SetMyMetadata([SetMyMetadataRequest](#setmymetadatarequest))
[SetMyMetadataResponse](#setmymetadataresponse)

Sets a user meta data by key to the authorized user



    POST: /users/me/metadata/{key}


### BulkSetMyMetadata

> **rpc** BulkSetMyMetadata([BulkSetMyMetadataRequest](#bulksetmymetadatarequest))
[BulkSetMyMetadataResponse](#bulksetmymetadataresponse)

Set a list of user meta data to the authorized user



    POST: /users/me/metadata/_bulk


### ListMyMetadata

> **rpc** ListMyMetadata([ListMyMetadataRequest](#listmymetadatarequest))
[ListMyMetadataResponse](#listmymetadataresponse)

Returns the user meta data of the authorized user



    POST: /users/me/metadata/_search


### GetMyMetadata

> **rpc** GetMyMetadata([GetMyMetadataRequest](#getmymetadatarequest))
[GetMyMetadataResponse](#getmymetadataresponse)

Returns the user meta data by key of the authorized user



    GET: /users/me/metadata/{key}


### RemoveMyMetadata

> **rpc** RemoveMyMetadata([RemoveMyMetadataRequest](#removemymetadatarequest))
[RemoveMyMetadataResponse](#removemymetadataresponse)

Removes a user meta data by key to the authorized user



    DELETE: /users/me/metadata/{key}


### BulkRemoveMyMetadata

> **rpc** BulkRemoveMyMetadata([BulkRemoveMyMetadataRequest](#bulkremovemymetadatarequest))
[BulkRemoveMyMetadataResponse](#bulkremovemymetadataresponse)

Set a list of user meta data to the authorized user



    DELETE: /users/me/metadata/_bulk


### ListMyRefreshTokens

> **rpc** ListMyRefreshTokens([ListMyRefreshTokensRequest](#listmyrefreshtokensrequest))
[ListMyRefreshTokensResponse](#listmyrefreshtokensresponse)

Returns the refresh tokens of the authorized user



    POST: /users/me/tokens/refresh/_search


### RevokeMyRefreshToken

> **rpc** RevokeMyRefreshToken([RevokeMyRefreshTokenRequest](#revokemyrefreshtokenrequest))
[RevokeMyRefreshTokenResponse](#revokemyrefreshtokenresponse)

Revokes a single refresh token of the authorized user by its (token) id



    DELETE: /users/me/tokens/refresh/{id}


### RevokeAllMyRefreshTokens

> **rpc** RevokeAllMyRefreshTokens([RevokeAllMyRefreshTokensRequest](#revokeallmyrefreshtokensrequest))
[RevokeAllMyRefreshTokensResponse](#revokeallmyrefreshtokensresponse)

Revokes all refresh tokens of the authorized user



    POST: /users/me/tokens/refresh/_revoke_all


### UpdateMyUserName

> **rpc** UpdateMyUserName([UpdateMyUserNameRequest](#updatemyusernamerequest))
[UpdateMyUserNameResponse](#updatemyusernameresponse)

Change the user name of the authorize user



    PUT: /users/me/username


### GetMyPasswordComplexityPolicy

> **rpc** GetMyPasswordComplexityPolicy([GetMyPasswordComplexityPolicyRequest](#getmypasswordcomplexitypolicyrequest))
[GetMyPasswordComplexityPolicyResponse](#getmypasswordcomplexitypolicyresponse)

Returns the password complexity policy of my organisation
This policy defines how the password should look



    GET: /policies/passwords/complexity


### UpdateMyPassword

> **rpc** UpdateMyPassword([UpdateMyPasswordRequest](#updatemypasswordrequest))
[UpdateMyPasswordResponse](#updatemypasswordresponse)

Change the password of the authorized user



    PUT: /users/me/password


### GetMyProfile

> **rpc** GetMyProfile([GetMyProfileRequest](#getmyprofilerequest))
[GetMyProfileResponse](#getmyprofileresponse)

Returns the profile information of the authorized user



    GET: /users/me/profile


### UpdateMyProfile

> **rpc** UpdateMyProfile([UpdateMyProfileRequest](#updatemyprofilerequest))
[UpdateMyProfileResponse](#updatemyprofileresponse)

Changes the profile information of the authorized user



    PUT: /users/me/profile


### GetMyEmail

> **rpc** GetMyEmail([GetMyEmailRequest](#getmyemailrequest))
[GetMyEmailResponse](#getmyemailresponse)

Returns the email address of the authorized user



    GET: /users/me/email


### SetMyEmail

> **rpc** SetMyEmail([SetMyEmailRequest](#setmyemailrequest))
[SetMyEmailResponse](#setmyemailresponse)

Changes the email address of the authorized user
An email is sent to the given address, to verify it



    PUT: /users/me/email


### VerifyMyEmail

> **rpc** VerifyMyEmail([VerifyMyEmailRequest](#verifymyemailrequest))
[VerifyMyEmailResponse](#verifymyemailresponse)

Sets the email address to verified



    POST: /users/me/email/_verify


### ResendMyEmailVerification

> **rpc** ResendMyEmailVerification([ResendMyEmailVerificationRequest](#resendmyemailverificationrequest))
[ResendMyEmailVerificationResponse](#resendmyemailverificationresponse)

Sends a new email to the last given address to verify it



    POST: /users/me/email/_resend_verification


### GetMyPhone

> **rpc** GetMyPhone([GetMyPhoneRequest](#getmyphonerequest))
[GetMyPhoneResponse](#getmyphoneresponse)

Returns the phone number of the authorized user



    GET: /users/me/phone


### SetMyPhone

> **rpc** SetMyPhone([SetMyPhoneRequest](#setmyphonerequest))
[SetMyPhoneResponse](#setmyphoneresponse)

Sets the phone number of the authorized user
An sms is sent to the number with a verification code



    PUT: /users/me/phone


### VerifyMyPhone

> **rpc** VerifyMyPhone([VerifyMyPhoneRequest](#verifymyphonerequest))
[VerifyMyPhoneResponse](#verifymyphoneresponse)

Sets the phone number to verified



    POST: /users/me/phone/_verify


### ResendMyPhoneVerification

> **rpc** ResendMyPhoneVerification([ResendMyPhoneVerificationRequest](#resendmyphoneverificationrequest))
[ResendMyPhoneVerificationResponse](#resendmyphoneverificationresponse)

Resends a sms to the last given phone number, to verify it



    POST: /users/me/phone/_resend_verification


### RemoveMyPhone

> **rpc** RemoveMyPhone([RemoveMyPhoneRequest](#removemyphonerequest))
[RemoveMyPhoneResponse](#removemyphoneresponse)

Removed the phone number of the authorized user



    DELETE: /users/me/phone


### RemoveMyAvatar

> **rpc** RemoveMyAvatar([RemoveMyAvatarRequest](#removemyavatarrequest))
[RemoveMyAvatarResponse](#removemyavatarresponse)

Remove my avatar



    DELETE: /users/me/avatar


### ListMyLinkedIDPs

> **rpc** ListMyLinkedIDPs([ListMyLinkedIDPsRequest](#listmylinkedidpsrequest))
[ListMyLinkedIDPsResponse](#listmylinkedidpsresponse)

Returns a list of all linked identity providers (social logins, eg. Google, Microsoft, AD, etc.)



    POST: /users/me/idps/_search


### RemoveMyLinkedIDP

> **rpc** RemoveMyLinkedIDP([RemoveMyLinkedIDPRequest](#removemylinkedidprequest))
[RemoveMyLinkedIDPResponse](#removemylinkedidpresponse)

Removes a linked identity provider (social logins, eg. Google, Microsoft, AD, etc.)



    DELETE: /users/me/idps/{idp_id}/{linked_user_id}


### ListMyAuthFactors

> **rpc** ListMyAuthFactors([ListMyAuthFactorsRequest](#listmyauthfactorsrequest))
[ListMyAuthFactorsResponse](#listmyauthfactorsresponse)

Returns all configured authentication factors (second and multi)



    POST: /users/me/auth_factors/_search


### AddMyAuthFactorOTP

> **rpc** AddMyAuthFactorOTP([AddMyAuthFactorOTPRequest](#addmyauthfactorotprequest))
[AddMyAuthFactorOTPResponse](#addmyauthfactorotpresponse)

Adds a new OTP (One Time Password) Second Factor to the authorized user
Only one OTP can be configured per user



    POST: /users/me/auth_factors/otp


### VerifyMyAuthFactorOTP

> **rpc** VerifyMyAuthFactorOTP([VerifyMyAuthFactorOTPRequest](#verifymyauthfactorotprequest))
[VerifyMyAuthFactorOTPResponse](#verifymyauthfactorotpresponse)

Verify the last added OTP (One Time Password)



    POST: /users/me/auth_factors/otp/_verify


### RemoveMyAuthFactorOTP

> **rpc** RemoveMyAuthFactorOTP([RemoveMyAuthFactorOTPRequest](#removemyauthfactorotprequest))
[RemoveMyAuthFactorOTPResponse](#removemyauthfactorotpresponse)

Removed the configured OTP (One Time Password) Factor



    DELETE: /users/me/auth_factors/otp


### AddMyAuthFactorU2F

> **rpc** AddMyAuthFactorU2F([AddMyAuthFactorU2FRequest](#addmyauthfactoru2frequest))
[AddMyAuthFactorU2FResponse](#addmyauthfactoru2fresponse)

Adds a new U2F (Universal Second Factor) to the authorized user
Multiple U2Fs can be configured



    POST: /users/me/auth_factors/u2f


### VerifyMyAuthFactorU2F

> **rpc** VerifyMyAuthFactorU2F([VerifyMyAuthFactorU2FRequest](#verifymyauthfactoru2frequest))
[VerifyMyAuthFactorU2FResponse](#verifymyauthfactoru2fresponse)

Verifies the last added U2F (Universal Second Factor) of the authorized user



    POST: /users/me/auth_factors/u2f/_verify


### RemoveMyAuthFactorU2F

> **rpc** RemoveMyAuthFactorU2F([RemoveMyAuthFactorU2FRequest](#removemyauthfactoru2frequest))
[RemoveMyAuthFactorU2FResponse](#removemyauthfactoru2fresponse)

Removes the U2F Authentication from the authorized user



    DELETE: /users/me/auth_factors/u2f/{token_id}


### ListMyPasswordless

> **rpc** ListMyPasswordless([ListMyPasswordlessRequest](#listmypasswordlessrequest))
[ListMyPasswordlessResponse](#listmypasswordlessresponse)

Returns all configured passwordless authenticators of the authorized user



    POST: /users/me/passwordless/_search


### AddMyPasswordless

> **rpc** AddMyPasswordless([AddMyPasswordlessRequest](#addmypasswordlessrequest))
[AddMyPasswordlessResponse](#addmypasswordlessresponse)

Adds a new passwordless authenticator to the authorized user
Multiple passwordless authentications can be configured



    POST: /users/me/passwordless


### AddMyPasswordlessLink

> **rpc** AddMyPasswordlessLink([AddMyPasswordlessLinkRequest](#addmypasswordlesslinkrequest))
[AddMyPasswordlessLinkResponse](#addmypasswordlesslinkresponse)

Adds a new passwordless authenticator link to the authorized user and returns it directly
This link enables the user to register a new device if current passwordless devices are all platform authenticators
e.g. User has already registered Windows Hello and wants to register FaceID on the iPhone



    POST: /users/me/passwordless/_link


### SendMyPasswordlessLink

> **rpc** SendMyPasswordlessLink([SendMyPasswordlessLinkRequest](#sendmypasswordlesslinkrequest))
[SendMyPasswordlessLinkResponse](#sendmypasswordlesslinkresponse)

Adds a new passwordless authenticator link to the authorized user and sends it to the registered email address
This link enables the user to register a new device if current passwordless devices are all platform authenticators
e.g. User has already registered Windows Hello and wants to register FaceID on the iPhone



    POST: /users/me/passwordless/_send_link


### VerifyMyPasswordless

> **rpc** VerifyMyPasswordless([VerifyMyPasswordlessRequest](#verifymypasswordlessrequest))
[VerifyMyPasswordlessResponse](#verifymypasswordlessresponse)

Verifies the last added passwordless configuration



    POST: /users/me/passwordless/_verify


### RemoveMyPasswordless

> **rpc** RemoveMyPasswordless([RemoveMyPasswordlessRequest](#removemypasswordlessrequest))
[RemoveMyPasswordlessResponse](#removemypasswordlessresponse)

Removes the passwordless configuration from the authorized user



    DELETE: /users/me/passwordless/{token_id}


### ListMyUserGrants

> **rpc** ListMyUserGrants([ListMyUserGrantsRequest](#listmyusergrantsrequest))
[ListMyUserGrantsResponse](#listmyusergrantsresponse)

Returns all user grants (authorizations) of the authorized user



    POST: /usergrants/me/_search


### ListMyProjectOrgs

> **rpc** ListMyProjectOrgs([ListMyProjectOrgsRequest](#listmyprojectorgsrequest))
[ListMyProjectOrgsResponse](#listmyprojectorgsresponse)

Returns a list of organisations where the authorized user has a user grant (authorization) in the context of the requested project



    POST: /global/projectorgs/_search


### ListMyZitadelFeatures

> **rpc** ListMyZitadelFeatures([ListMyZitadelFeaturesRequest](#listmyzitadelfeaturesrequest))
[ListMyZitadelFeaturesResponse](#listmyzitadelfeaturesresponse)

Returns a list of features, which are allowed on these organisation based on the subscription of the organisation



    POST: /features/zitadel/me/_search


### ListMyZitadelPermissions

> **rpc** ListMyZitadelPermissions([ListMyZitadelPermissionsRequest](#listmyzitadelpermissionsrequest))
[ListMyZitadelPermissionsResponse](#listmyzitadelpermissionsresponse)

Returns the permissions the authorized user has in ZITADEL based on his manager roles (e.g ORG_OWNER)



    POST: /permissions/zitadel/me/_search


### ListMyProjectPermissions

> **rpc** ListMyProjectPermissions([ListMyProjectPermissionsRequest](#listmyprojectpermissionsrequest))
[ListMyProjectPermissionsResponse](#listmyprojectpermissionsresponse)

Returns a list of roles for the authorized user and project



    POST: /permissions/me/_search







## Messages


### AddMyAuthFactorOTPRequest
This is an empty request




### AddMyAuthFactorOTPResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| url |  string | - |  |
| secret |  string | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddMyAuthFactorU2FRequest
This is an empty request




### AddMyAuthFactorU2FResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  zitadel.user.v1.WebAuthNKey | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### AddMyPasswordlessLinkRequest
This is an empty request




### AddMyPasswordlessLinkResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| link |  string | - |  |
| expiration |  google.protobuf.Duration | - |  |




### AddMyPasswordlessRequest
This is an empty request




### AddMyPasswordlessResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  zitadel.user.v1.WebAuthNKey | - |  |
| details |  zitadel.v1.ObjectDetails | - |  |




### BulkRemoveMyMetadataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| keys | repeated string | - | repeated.items.string.min_len: 1<br /> repeated.items.string.max_len: 200<br />  |




### BulkRemoveMyMetadataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### BulkSetMyMetadataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| metadata | repeated BulkSetMyMetadataRequest.Metadata | - |  |




### BulkSetMyMetadataRequest.Metadata



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| value |  string | - | string.min_len: 1<br />  |




### BulkSetMyMetadataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### GetMyEmailRequest
This is an empty request




### GetMyEmailResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| email |  zitadel.user.v1.Email | - |  |




### GetMyMetadataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### GetMyMetadataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| metadata |  zitadel.metadata.v1.Metadata | - |  |




### GetMyPasswordComplexityPolicyRequest
This is an empty request




### GetMyPasswordComplexityPolicyResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| policy |  zitadel.policy.v1.PasswordComplexityPolicy | - |  |




### GetMyPhoneRequest
This is an empty request




### GetMyPhoneResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| phone |  zitadel.user.v1.Phone | - |  |




### GetMyProfileRequest
This is an empty request




### GetMyProfileResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |
| profile |  zitadel.user.v1.Profile | - |  |




### GetMyUserRequest
This is an empty request
the request parameters are read from the token-header




### GetMyUserResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| user |  zitadel.user.v1.User | - |  |
| last_login |  google.protobuf.Timestamp | - |  |




### GetSupportedLanguagesRequest
This is an empty request




### GetSupportedLanguagesResponse
This is an empty response


| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| languages | repeated string | - |  |




### HealthzRequest
This is an empty request




### HealthzResponse
This is an empty response




### ListMyAuthFactorsRequest
This is an empty request




### ListMyAuthFactorsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated zitadel.user.v1.AuthFactor | - |  |




### ListMyLinkedIDPsRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | list limitations and ordering |  |




### ListMyLinkedIDPsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.idp.v1.IDPUserLink | - |  |




### ListMyMetadataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| query |  zitadel.v1.ListQuery | - |  |
| queries | repeated zitadel.metadata.v1.MetadataQuery | - |  |




### ListMyMetadataResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.metadata.v1.Metadata | - |  |




### ListMyPasswordlessRequest
This is an empty request




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
This is an empty request




### ListMyProjectPermissionsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated string | - |  |




### ListMyRefreshTokensRequest
This is an empty request




### ListMyRefreshTokensResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ListDetails | - |  |
| result | repeated zitadel.user.v1.RefreshToken | - |  |




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
This is an empty request




### ListMyUserSessionsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated zitadel.user.v1.Session | - |  |




### ListMyZitadelFeaturesRequest
This is an empty request




### ListMyZitadelFeaturesResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated string | - |  |




### ListMyZitadelPermissionsRequest
This is an empty request




### ListMyZitadelPermissionsResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| result | repeated string | - |  |




### RemoveMyAuthFactorOTPRequest
This is an empty request




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




### RemoveMyAvatarRequest
This is an empty request




### RemoveMyAvatarResponse



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




### RemoveMyMetadataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RemoveMyMetadataResponse



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
This is an empty request




### RemoveMyPhoneResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResendMyEmailVerificationRequest
This is an empty request




### ResendMyEmailVerificationResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### ResendMyPhoneVerificationRequest
This is an empty request




### ResendMyPhoneVerificationResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### RevokeAllMyRefreshTokensRequest
This is an empty request




### RevokeAllMyRefreshTokensResponse
This is an empty response




### RevokeMyRefreshTokenRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| id |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |




### RevokeMyRefreshTokenResponse



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| details |  zitadel.v1.ObjectDetails | - |  |




### SendMyPasswordlessLinkRequest
This is an empty request




### SendMyPasswordlessLinkResponse



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




### SetMyMetadataRequest



| Field | Type | Description | Validation |
| ----- | ---- | ----------- | ----------- |
| key |  string | - | string.min_len: 1<br /> string.max_len: 200<br />  |
| value |  string | - | string.min_len: 1<br />  |




### SetMyMetadataResponse



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







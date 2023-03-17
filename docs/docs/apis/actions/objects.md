---
title: Objects
---

## External User

- `externalId` *string*  
  User id from the identity provider
- `externalIdpId` *string*  
  Id of the identity provider
- `human`
  - `firstName` *string*
  - `lastName` *string*
  - `nickName` *string*
  - `displayName` *string*
  - `preferredLanguage` *string*  
    In [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646) fromat
  - `email` *string*
  - `isEmailVerified` *boolean*
  - `phone` *string*
  - `isPhoneVerified` *boolean*

## metadata with value as bytes

- `key` *string*
- `value` Array of *byte*

## metadata result

- `count` *number*
- `sequence` *number*
- `timestamp` *Date*
- `metadata` Array of [*metadata*](#metadata)

## metadata

- `creationDate` *Date*
- `changeDate` *Date*
- `resourceOwner` *string*
- `sequence` *number*
- `key` *string*
- `value` `Any`

## user grant

- `projectId` *string*  
  Required. Id of the project to be granted
- `projectGrantId` *string*  
  If the grant is for a project grant
- `roles` Array of *string*  
  Containing the roles

## user

- `id` *string*
- `creationDate` *Date*
- `changeDate` *Date*
- `resourceOwner` *string*
- `sequence` *number*  
  Unsigned 64 bit integer
- `state` *number*  
  <ul><li>0: unspecified</li><li>1: active</li><li>2: inactive</li><li>3: deleted</li><li>4: locked</li><li>5: suspended</li><li>6: initial</li></ul>
- `username` *string*
- `loginNames` Array of *string*
- `preferredLoginName` *string*
- `human`  
  Set if the user is human
  - `firstName` *string*
  - `lastName` *string*
  - `nickName` *string*
  - `displayName` *string*
  - `avatarKey` *string*
  - `preferredLanguage` *string*  
    In [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646) fromat
  - `gender` *number*  
    <ul><li>0: unspecified</li><li>1: female</li><li>2: male</li><li>3: diverse</li></ul>
  - `email` *string*
  - `isEmailVerified` *boolean*
  - `phone` *string*
  - `isPhoneVerified` *boolean*
- `machine`  
  Set if the user is a machine
  - `name` *string*
  - `description` *string*

## human user

- `id` *string*
- `creationDate` *Date*
- `changeDate` *Date*
- `resourceOwner` *string*
- `sequence` *number*
- `state` *number*  
  <ul><li>0: unspecified</li><li>1: active</li><li>2: inactive</li><li>3: deleted</li><li>4: locked</li><li>5: suspended</li><li>6: initial</li></ul>
- `username` *string*
- `loginNames` Array of *string*
- `preferredLoginName` *string*
- `profile`
  - `firstName` *string*
  - `lastName` *string*
  - `nickName` *string*
  - `displayName` *string*
  - `preferredLanguage` *string*  
    In [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646) fromat
- `email`
  - `email` *string*
  - `isEmailVerified` *boolean*
- `phone`
  - `phone` *string*
  - `isPhoneVerified` *boolean*

## Auth Request

This object contains context information about the request to the [authorization endpoint](/docs/apis/openidoauth/endpoints#authorization_endpoint).

- `id` *string*
- `agentId` *string*
- `creationDate` *Date*
- `changeDate` *Date*
- `browserInfo` *browserInfo*
  - `userAgent` *string*
  - `acceptLanguage` *string*
  - `remoteIp` *string*
- `applicationId` *string*
- `callbackUri` *string*
- `transferState` *string*
- `prompt` Array of *Number*  
   <ul><li>0: not specified</li><li>1: none</li><li>2: login</li><li>3: consent</li><li>4: select_account</li><li>5: create</li></ul>
- `uiLocales` Array of *string*
- `loginHint` *string*
- `maxAuthAge` *Number*  
  Duration in nanoseconds
- `instanceId` *string*
- `request`
  - `oidc`
    - `scopes` Array of *string*
- `userId` *string*
- `userName` *string*
- `loginName` *string*
- `displayName` *string*
- `resourceOwner` *string*
- `requestedOrgId` *string*
- `requestedOrgName` *string*
- `requestedPrimaryDomain` *string*
- `requestedOrgDomain` *bool*
- `applicationResourceOwner` *string*
- `privateLabelingSetting` *Number*
  <ul><li>0: Unspecified</li><li>1: Enforce project resource owner policy</li><li>2: Allow login user resource owner policy</li></ul>
- `selectedIdpConfigId` *string*
- `linkingUsers` Array of [*ExternalUser*](#external-user)
- `passwordVerified` *bool*
- `mfasVerified` Array of *Number*  
  <ul><li>0: OTP</li><li>1: U2F</li><li>2: U2F User verification</li></ul>
- `audience` Array of *string*
- `authTime` *Date*

## HTTP Request

This object is based on the Golang struct [http.Request](https://pkg.go.dev/net/http#Request), some attributes are removed as not all provided information is usable in this context.

- `method` *string*
- `url` *string*
- `proto` *string*
- `contentLength` *number*
- `host` *string*
- `form` Map *string* of Array of *string*
- `postForm` Map *string* of Array of *string*
- `remoteAddr` *string*
- `headers` Map *string* of Array of *string*

## Claims

This object represents [the claims](../openidoauth/claims) which will be written into the oidc token.

- `sub` *string*
- `name` *string*
- `email` *string*
- `locale` *string*
- `given_name` *string*
- `family_name` *string*
- `preferred_username` *string*
- `email_verified` *bool*
- `updated_at` *Number*

Additionally there could additional fields depending on the configuration of your [project](../../guides/manage/console/projects#role-settings) and your [application](../../guides/manage/console/applications#token-settings)

## user grant list

This object represents a list of user grant stored in ZITADEL.

- `count` *Number*
- `sequence` *Number*
- `timestamp` *Date*
- `grants` Array of
  - `id` *string*
  - `projectGrantId` *string*  
    The id of the [project grant](../../concepts/usecases/saas#project-grant)
  - `state` *Number*  
    <ul><li>0: unspecified</li><li>1: active</li><li>2: inactive</li><li>3: removed</li></ul>
  - `creationDate` *Date*
  - `changeDate` *Date*
  - `sequence` *Number*
  - `userId` *string*
  - `roles` Array of *string*
  - `userResourceOwner` *string*
  - `userGrantResourceOwner` *string*
  - `userGrantResourceOwnerName` *string*
  - `projectId` *string*
  - `projectName` *string*

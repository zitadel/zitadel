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

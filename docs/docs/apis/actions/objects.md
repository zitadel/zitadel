---
title: Objects
---

## External User

- `externalId`: string
- `externalIdpId`: string
- `human`
  - `firstName`: `string`
  - `lastName`: `string`
  - `nickName`: `string`
  - `displayName`: `string`
  - `preferredLanguage`: `string`
  - `email`: `string`
  - `isEmailVerified`: `boolean`
  - `phone`: `string`
  - `isPhoneVerified`: `boolean`

## id token claims

## metadata

- key: `string`
- value: Array of `byte`

## user grant

- `projectId`: Required. Id of the project to be granted
- `projectGrantId`: If the grant is for a project grant
- `roles`: Array of `string` containing the roles

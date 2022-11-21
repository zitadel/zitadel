---
title: Login flows
---

## Complement Token

This flow is executed during the creation of tokens and token introspection.

### Pre Userinfo creation

This trigger is called before userinfo are set in the token or response.

#### Parameters of Pre Userinfo creation

`ctx`: is always `null`

`api`:

**Fields**: None

**Methods**

| name | description | parameter types | response |
|---|---|---|---|
| setClaim(key, value) | sets an additional claim in user info. The claim can be set once and the key must not be a reserved key | `string`, `any` | none |
| appendLogIntoClaims(entry) | appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an `array` | `string` | none |

#### Available modules of Pre Userinfo creation

- [zitadel/http](#zitadelhttp)
- [zitadel/metadata/user](#zitadelmetadatauser)

### Pre access token creation

This trigger is called before the claims are set in the access token and the token type is `jwt`.

#### Parameters of Pre access token creation

`ctx`: is always `null`

`api`:

**Fields**: None

**Methods**

| name | description | parameter types | response |
|---|---|---|---|
| setClaim(key, value) | sets an additional claim in access token. The claim can be set once and the key must not be a reserved key. | `string`, `any` | none |
| appendLogIntoClaims(entry) | appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an `array` | `string` | none |
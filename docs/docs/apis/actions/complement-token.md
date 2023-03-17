---
title: Complement Token Flow
---

This flow is executed during the creation of tokens and token introspection.

## Pre Userinfo creation

This trigger is called before userinfo are set in the token or response.

### Parameters of Pre Userinfo creation

- `ctx`  
  The first parameter contains the following fields:
  - `v1`
    - `claims` [*Claims*](./objects#claims)
    - `getUser()` [*User*](./objects#user)
    - `user`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
    - `grants` [*UserGrantList*](./objects#user-grant-list)
- `api`  
  The second parameter contains the following fields:
  - `v1`
    - `userinfo`  
      This function is deprecated, please use `api.v1.claims`
      - `setClaim(string, Any)`  
        Sets any value if the key is not already present. If it's already present there is a message added to `urn:zitadel:iam:action:${action.name}:log`
    - `claims`
      - `setClaim(string, Any)`  
        Sets any value if the key is not already present. If it's already present there is a message added to `urn:zitadel:iam:action:${action.name}:log`
    - `user`
      - `setMetadata(string, Any)`  
        Key of the metadata and any value

## Pre access token creation

This trigger is called before the claims are set in the access token and the token type is `jwt`.

### Parameters of Pre access token creation

- `ctx`  
  The first parameter contains the following fields:
  - `v1`
    - `claims` [*Claims*](./objects#claims)
    - `getUser()` [*User*](./objects#user)
    - `user`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
    - `grants` [*UserGrantList*](./objects#user-grant-list)
- `api`  
  The second parameter contains the following fields:
  - `v1`
    - `claims`
      - `setClaim(string, Any)`  
        Sets any value if the key is not already present. If it's already present there is a message added to `urn:zitadel:iam:action:${action.name}:log`
      - `appendLogIntoClaims(string)`  
        Appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an Array of *string*
    - `user`
      - `setMetadata(string, Any)`  
        Key of the metadata and any value

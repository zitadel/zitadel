---
title: Complement Token Flow
---

This flow is executed during the creation of tokens and token introspection.

The flow is represented by the following Ids in the API: `2`

## Pre Userinfo creation (id_token / userinfo / introspection endpoint)

This trigger is called before userinfo are set in the id_token or userinfo and introspection endpoint response.

The trigger is represented by the following Ids in the API: `4`

### Parameters of Pre Userinfo creation

- `ctx`  
  The first parameter contains the following fields:
  - `v1`
    - `claims` [*Claims*](./objects#claims)
    - `getUser()` [*User*](./objects#user)
    - `user`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
      - `grants` [*UserGrantList*](./objects#user-grant-list)
    - `org`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
- `api`  
  The second parameter contains the following fields:
  - `v1`
    - `userinfo`  
      This function is deprecated, please use `api.v1.claims`
      - `setClaim(string, Any)`  
        Sets any value if the key is not already present. If it's already present there is a message added to `urn:zitadel:iam:action:${action.name}:log` 
        Note that keys with prefix `urn:zitadel:iam` will be ignored.
    - `claims`
      - `setClaim(string, Any)`  
        Sets any value if the key is not already present. If it's already present there is a message added to `urn:zitadel:iam:action:${action.name}:log`
        Note that keys with prefix `urn:zitadel:iam` will be ignored.
    - `user`
      - `setMetadata(string, Any)`  
        Key of the metadata and any value

## Pre access token creation

This trigger is called before the claims are set in the access token and the token type is `jwt`.

The trigger is represented by the following Ids in the API: `5`

### Parameters of Pre access token creation

- `ctx`  
  The first parameter contains the following fields:
  - `v1`
    - `claims` [*Claims*](./objects#claims)
    - `getUser()` [*User*](./objects#user)
    - `user`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
      - `grants` [*UserGrantList*](./objects#user-grant-list)
    - `org`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
- `api`  
  The second parameter contains the following fields:
  - `v1`
    - `claims`
      - `setClaim(string, Any)`  
        Sets any value if the key is not already present. If it's already present there is a message added to `urn:zitadel:iam:action:${action.name}:log`
        Note that keys with prefix `urn:zitadel:iam` will be ignored.
      - `appendLogIntoClaims(string)`  
        Appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an Array of *string*
    - `user`
      - `setMetadata(string, Any)`  
        Key of the metadata and any value

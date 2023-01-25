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
    - `user`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
- `api`  
  The second parameter contains the following fields:
  - `v1`
    - `userinfo`
      - `setClaim(string, Any)`  
        Key of the claim and any value
    - `user`
      - `setMetadata(string, Any)`  
        Key of the metadata and any value

## Pre access token creation

This trigger is called before the claims are set in the access token and the token type is `jwt`.

### Parameters of Pre access token creation

- `ctx`  
  The first parameter contains the following fields:
  - `v1`
    - `user`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
- `api`  
  The second parameter contains the following fields:
  - `v1`
    - `claims`
      - `setClaim(string, Any)`  
        Sets any value if the key is not already present
      - `appendLogIntoClaims(string)`  
        Appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an Array of *string*
    - `user`
      - `setMetadata(string, Any)`  
        Key of the metadata and any value

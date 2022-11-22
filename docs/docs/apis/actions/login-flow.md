---
title: Login flows
---

## Complement Token

This flow is executed during the creation of tokens and token introspection.

### Pre Userinfo creation

This trigger is called before userinfo are set in the token or response.

#### Parameters of Pre Userinfo creation

- `ctx`: The first parameter contains the following fields:
  - `v1`
    - `user`
      - `getMetadata()` [`metadataResult`](./objects#metadata-result)
- `api`: The second parameter contains the following fields:
  - `v1`
    - `userinfo`
      - `setClaim(string, Object)`: key of the claim and an object as value
    - `user`
      - `setMetadata(string, Object)`: key of the metadata and an object as value

### Pre access token creation

This trigger is called before the claims are set in the access token and the token type is `jwt`.

#### Parameters of Pre access token creation

- `ctx`: The first parameter contains the following fields:
  - `v1`
    - `user`
      - `getMetadata()` [`metadataResult`](./objects#metadata-result)
- `api`: The second parameter contains the following fields:
  - `v1`
    - `claims`
      - `setClaim(string, Object)`: sets the value if the key is not already present
      - `appendLogIntoClaims(string)`: Appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an Array of `string`
    - `user`
      - `setMetadata(string, Object)`: key of the metadata and an object as value

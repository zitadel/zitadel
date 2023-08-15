---
title: Complement SAMLResponse
---

This flow is executed before the return of the SAMLResponse.

## Pre SAMLResponse creation

This trigger is called before attributes are set in the SAMLResponse.

### Parameters of Pre SAMLResponse creation

- `ctx`  
  The first parameter contains the following fields:
  - `v1`
    - `getUser()` [*User*](./objects#user)
    - `user`
      - `getMetadata()` [*metadataResult*](./objects#metadata-result)
      - `grants` [*UserGrantList*](./objects#user-grant-list)
- `api`  
  The second parameter contains the following fields:
  - `v1`
    - `attributes`
      - `setCustomAttribute(string, string, ...string)`  
        Sets any value as attribute in addition to the default attributes, if the key is not already present. The parameters represent the key, nameFormat and the attributeValue(s).
    - `user`
      - `setMetadata(string, Any)`  
        Key of the metadata and any value

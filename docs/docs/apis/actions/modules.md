---
title: Modules
---

ZITADEL provides the following modules.

## HTTP

This module provides functionality to call REST APIs.

### Import

```js
    let http = require('zitadel/http')
```

### `fetch()` function

This function allows to call HTTP servers. The function does NOT fulfil the [Fetch API specification](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API).

#### Parameters

- `url` *string*
- `options`  
  **Optional**, containing custom settings that you want to apply to the request.
  - `headers`  
    Overwrites the default headers. One of the following types
    - *map[string] string*  
      The value is split into separate values after each comma `,`.
    - *map[string] Array of string*  
      The value is a string array
    - default:
      - `Content-Type`: `application/json`
      - `Accept`: `application/json`
  - `method`  
    The request method. Allowed values are `GET`, `POST`, `PUT`, `DELETE`
  - `body` *Object*  
    JSON representation

#### Response

If the request was invalid, an error will be thrown, otherwise a Response object will be returned.

The object has the following fields and methods:

- `status` *number*  
  Status code of response
- `body` *string*  
  Return value
- `json()` *Object*  
  Returns the body as JSON object, or throws an error if the body is not a json object.
- `text()` *string*  
  Returns the body

## UUID

This module provides functionality to generate a UUID

### Import

```js
    let uuid = require("zitadel/uuid")
```

### `uuid.vX()` function

This function generates a UUID using [google/uuid](https://github.com/google/uuid). `vX` allows to define the UUID version:

- `uuid.v1()` *string*  
  Generates a UUID version 1, based on date-time and MAC address
- `uuid.v3(namespace, data)` *string*  
  Generates a UUID version 3, based on the provided namespace using MD5
- `uuid.v4()` *string*  
  Generates a UUID version 4, which is randomly generated
- `uuid.v5(namespace, data)` *string*  
  Generates a UUID version 5, based on the provided namespace using SHA1

#### Parameters

- `namespace` *UUID*/*string*  
  Namespace to be used in the hashing function. Either provide one of defined [namespaces](#namespaces) or a string representing a UUID.
- `data` *[]byte*/*string*  
  data to be used in the hashing function. Possible types are []byte or string.

### Namespaces

The following predefined namespaces can be used for `uuid.v3` and `uuid.v5`:

- `uuid.namespaceDNS` *UUID*  
  6ba7b810-9dad-11d1-80b4-00c04fd430c8
- `uuid.namespaceURL` *UUID*  
  6ba7b811-9dad-11d1-80b4-00c04fd430c8
- `uuid.namespaceOID` *UUID*  
  6ba7b812-9dad-11d1-80b4-00c04fd430c8
- `uuid.namespaceX500` *UUID*  
  6ba7b814-9dad-11d1-80b4-00c04fd430c8

### Example
```js
let uuid = require("zitadel/uuid")
function setUUID(ctx, api) {
  if (api.metadata === undefined) {
    return;
  }

  api.v1.user.appendMetadata('custom-id', uuid.v4());
}
```
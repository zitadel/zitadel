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

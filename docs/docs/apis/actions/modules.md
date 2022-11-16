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

- url - url as string
- options - optional, object containing custom settings that you want to apply to the request.
  - `headers`: overwrites the default headers
    - map[string] string: The value is splitted into seperate values after each comma `,`.
    - map[string] string[]: The value is a string array
    - default:
      - `Content-Type`: `application/json`
      - `Accept`: `application/json`
  - `method`: The request method. Allowed values are `GET`, `POST`, `PUT`, `DELETE`
  - `body`: JSON object

#### Response

If the request was valid, an error will be thrown, otherwise a Reponse object will be returned.

The object has the following fields and methods:

- `status`: Status code as number
- `body`: return value as string
- `json()`: returns the body as JSON object, or throws an error if the body is not a json object.
- `text()`: returns the body as string

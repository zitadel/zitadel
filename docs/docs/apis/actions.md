---
title: Actions
---

This page describes the options you have when writing ZITADEL actions scripts.

## Language

ZITADEL interpretes the scripts as JavaScript.
Make sure your scripts are ECMAScript 5.1(+) compliant.
Go to the [goja GitHub page](https://github.com/dop251/goja) for detailed reference about the underlying library features and limitations.

Actions provide a defined set of libraries. The provided libraries vary depending on trigger types.

Actions do not have access to any libraries yet.
Also, sending HTTP requests is not supported yet.
[We plan to add such features in the future](https://zitadel.com/roadmap).

## Action

The action object describes the logic defined by the user. Actions can be linked to multiple [triggers](#flows).

The name of the action must corelate with the javascript function in the code section. This function will be called by ZITADEL.

For reading and mutating state, the runtime executes the function that has the same name as the action.
The function receives the JavaScript objects `ctx` and `api`.

The object `ctx` provides readable information as object properties and by callable functions.
The object `api` provides mutable properties and state mutating functions.

The script of an action called **doSomething** should have a function called `doSomething` and look something like this:

```js
function doSomething(ctx, api){
    // read from ctx and manipulate with api
}
```

## Flows

Flows are the links between an [action](#action) and a specific point during a user interaction with ZITADEL. These specific point are called [Trigger Types](#trigger-types).

## Trigger Types

Trigger types define the point during execution of request. Each trigger defines which readable information (`ctx`) and mutual properties (`api`) are passed into the called function as well as which libraries are available for the function. Each trigger type is described in the [available flow types section](#available-flow-types)

## Available Flow Types

This section describes the flow types available.

### External Authentication

<!-- link idp and jwt -->
This flow is executed if the user logs in using an identity provider or using a jwt token.

#### Post Authentication

A user has authenticated externally. ZITADEL retrieved and mapped the external information.

##### Parameters of Post Authentication

`ctx`:

**Fields:**

| name | description | type |
|---|---|---|
| accessToken | the access token which will be returned to the user. This can be an opaque token or a JWT | string |
| idToken | the id token which will be returned to the user | string |

**Methods:**

| name | description | return value |
|---|---|---|---|
| getClaims(string) | returns the requested `claim` | `any` |
| claimsJSON() | Returns the complete payload of the `ctx.idToken` | `Object` |

`api`:

**Fields:** none

| name | description |
|---|---|
| metadata | array of metadata `{key: string, value: bytes}` |

**Methods:**

| name | description | return value |
|---|---|---|
| setFirstName(string) | sets the first name | none |
| setLastName(string) | sets the last name | none |
| setNickName(string) | sets the nick name | none |
| setDisplayName(string) | sets the display name | none |
| setPreferredLanguage(string) | sets the preferred language, the string has to be a valid language tag | none |
| setGender(int) | sets the gender. <br/><ul><li>0: unspecified</li><li>1: female</li><li>2: male</li><li>3: diverse</li></ul> | none |
| setPreferredUsername(string) | sets the preferred username | none |
| setEmail(string) | sets the email | |
| setEmailVerified(bool) | if true the email set is verified without user interaction | |
| setPhone(string) | sets the phone number | |
| setPhoneVerified(bool) | if true the phone number set is verified without user interaction | |

#### Pre Creation

A user selected **Register** on the overview page after external authentication. ZITADEL did not create the user yet.

##### Parameters of Pre Creation

`ctx`:

**Fields:**

| name | description | type |
|---|---|---|
| accessToken | the access token which will be returned to the user. This can be an opaque token or a JWT | string |
| idToken | the id token which will be returned to the user | string |

**Methods:**

| name | description | return value |
|---|---|---|---|
| getClaims(string) | returns the requested `claim` | `any` |
| claimsJSON() | Returns the complete payload of the `ctx.idToken` | `Object` |

`api`:

**Fields:** none

| name | description |
|---|---|
| metadata | array of metadata `{key: string, value: bytes}` |

**Methods:**

| name | description | return value |
|---|---|---|
| setFirstName(string) | sets the first name | none |
| setLastName(string) | sets the last name | none |
| setNickName(string) | sets the nick name | none |
| setDisplayName(string) | sets the display name | none |
| setPreferredLanguage(string) | sets the preferred language, the string has to be a valid language tag | none |
| setGender(int) | sets the gender. <br/><ul><li>0: unspecified</li><li>1: female</li><li>2: male</li><li>3: diverse</li></ul> | none |
| setUsername(string) | sets the username | none |
| setEmail(string) | sets the email | |
| setEmailVerified(bool) | if true the email set is verified without user interaction | |
| setPhone(string) | sets the phone number | |
| setPhoneVerified(bool) | if true the phone number set is verified without user interaction | |

#### Post Creation

A user selected **Register** on the overview page after external authentication. ZITADEL created the user.

##### Parameters of Post Creation

`ctx`:

**Fields:**

| name | description | type |
|---|---|---|
| accessToken | the access token which will be returned to the user. This can be an opaque token or a JWT | string |
| idToken | the id token which will be returned to the user | string |

**Methods:**

| name | description | return value |
|---|---|---|---|
| getClaims(string) | returns the requested `claim` | `any` |
| claimsJSON() | Returns the complete payload of the `ctx.idToken` | `Object` |

`api`:

**Fields:** none

| name | description |
|---|---|
| metadata | array of metadata `{key: string, value: bytes}` |
| userGrants | array of user grants `{projectID: string, projectGrantID: (optional)string, roles: [string]}` |

**Methods:**

| name | description | return value |
|---|---|---|
| setFirstName(string) | sets the first name | none |
| setLastName(string) | sets the last name | none |
| setNickName(string) | sets the nick name | none |
| setDisplayName(string) | sets the display name | none |
| setPreferredLanguage(string) | sets the preferred language, the string has to be a valid language tag | none |
| setGender(int) | sets the gender. <br/><ul><li>0: unspecified</li><li>1: female</li><li>2: male</li><li>3: diverse</li></ul> | none |
| setEmail(string) | sets the email | |
| setEmailVerified(bool) | if true the email set is verified without user interaction | |
| setPhone(string) | sets the phone number | |
| setPhoneVerified(bool) | if true the phone number set is verified without user interaction | |

### Complement Token

This flow is executed during the creation of tokens and token introspection.

#### Pre Userinfo creation

This trigger is called before userinfo are set in the token or response.

##### Parameters of Pre Userinfo creation

`ctx`: is always `null`

`api`:

**Fields**: None

**Methods**

| name | description | parameter types | response |
|---|---|---|---|
| setClaim(key, value) | sets an additional claim in user info. The claim can be set once and the key must not be a reserved key | `string`, `any` | none |
| appendLogIntoClaims(entry) | appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an `array` | `string` | none |

##### Available modules of Pre Userinfo creation

- [zitadel/http](#zitadelhttp)
- [zitadel/metadata/user](#zitadelmetadatauser)

#### Pre access token creation

This trigger is called before the claims are set in the access token and the token type is `jwt`.

##### Parameters of Pre access token creation

`ctx`: is always `null`

`api`:

**Fields**: None

**Methods**

| name | description | parameter types | response |
|---|---|---|---|
| setClaim(key, value) | sets an additional claim in access token. The claim can be set once and the key must not be a reserved key. | `string`, `any` | none |
| appendLogIntoClaims(entry) | appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an `array` | `string` | none |

##### Available modules of Pre access token creation

- [zitadel/http](#zitadelhttp)
- [zitadel/metadata/user](#zitadelmetadatauser)

## Provided Modules

This section descibes the modules which can be `require`d by actions if available.

### zitadel/http

This modules provides http functionality.

#### fetch

Fetch calls defined API's. It abstracts the golang `http.Do` function.

parameters:

- url: `string`
- options (optional): `Object`
  - headers: `map`[`string`](`string` OR `string array`)
  - method: `string`
    - one of `GET`, `POST`, `PUT`, `DELETE`
  - body: `Object`

response: `Object`

- body: `string`
- statusCode: `int`
- headers `map[string][]string`
- json(): Object representation of `body`
- text(): string representation of `body`

example:

```javascript
let http = require('zitadel/http');

function callPostman() {
  let res = http.fetch('http://postman-echo.com/get', {})
  let headersOfResponseBody = res.json().headers
}
```

### zitadel/metadata/user

This library abstracts the storage of metadata of the given user.

#### set

parameters:

- key: `string`
  - it's recommended to use urn-annotated keys
- value: `any`
  - The value will be json-marshalled before its stored

response: none

example:

```javascript
let userMD = require('zitadel/metadata/user')

const KEY = 'urn:mycorp:example'

function example() {
  let myNewMD = md.get().metadata.find(md => md.key === KEY)
}
```

#### get

parameters: none

response: `Object`

- count: `int`
  - amount of metadata found
- sequence: `int`
  - the sequence of the last event processed to calculate the response
- timestamp: `Date`
  - the timestamp of the last event processed to calculate the response
- metadata: `array of object`
  - creationDate: `Date`
  - changeDate: `Date`
  - resourceOwner: `string`
    - org id of the owner of this key
  - sequence: `int`
    - last event sequence which changed this key
  - key: `string`
  - value: `Object`
    - json representation of the value stored

example:

```javascript
let userMD = require('zitadel/metadata/user')

const KEY = 'urn:mycorp:example'

function example() {
  userMD.set(KEY, {key: 'value'})
}
```

## Further reading

- [Actions concept](../concepts/features/actions)
- [Actions guide](../guides/manage/customize/behavior)
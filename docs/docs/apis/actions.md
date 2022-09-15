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

This trigger is called after the oidc tokens are created.

<!-- 
##### Parameters of Post Authentication

`ctx`:

| field name | description | type | methods |
|---|---|---|---|
| accessToken | the access token which will be returned to the user | string |  |
| idToken | the id token which will be returned to the user | string |  |

`api`:

| field name | description | type | methods |
|---|---|---|---|
|  |  |  |  |
|  |  |  |  |

##### Functions of Post Authentication

The following functions are available if `id token claims` are available.

| name | description | parameters | return types |
|---|---|---|---|
| getClaims(claim) | returns the requested `claim` | claim string representation of the claim | Object |
| claimsJSON() | returns all claims of the `id token` | none | Object |

##### Modules of Post Authentication
-->

#### Pre Creation

This trigger is called before the oidc tokens are created.

#### Post Creation

This trigger is called after the oidc tokens are created.

### Complement Token

This flow is executed during the creation of tokens and token introspection.

#### Pre Userinfo creation

This trigger is called before userinfo are set in the token or response.

##### Parameters of Pre Userinfo creation

`ctx`: is always `null`

`api`:

**Fields**

None

**Methods**

| name | description | parameter types | response |
|---|---|---|---|
| setClaim(key, value) | sets an additional claim in user info. The claim can be set once and the key must not be a reserved key | `string`, `any` | none |
| appendLogIntoClaims(entry) | appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an `array` | `string` | none |

##### Available modules of Pre Userinfo creation

- [zitadel/http](#zitadel/http)
- [zitadel/metadata/user](#zitadel/metadata/user)

#### Pre access token creation

This trigger is called before the claims are set in the access token and the token type is `jwt`.

##### Parameters of Pre access token creation

`ctx`: is always `null`

`api`:

**Fields**

None

**Methods**

| name | description | parameter types | response |
|---|---|---|---|
| setClaim(key, value) | sets an additional claim in access token. The claim can be set once and the key must not be a reserved key. | `string`, `any` | none |
| appendLogIntoClaims(entry) | appends the entry into the claim `urn:zitadel:action:{action.name}:log` the value of the claim is an `array` | `string` | none |

##### Available modules of Pre access token creation

- [zitadel/http](#zitadel/http)
- [zitadel/metadata/user](#zitadel/metadata/user)

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

response: Object

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

parameters:

- asdf

response: Object

- asdf

example:

```javascript
let userMD = require('zitadel/metadata/user')

const KEY = 'urn:mycorp:example'

function example() {
  userMD.set(KEY, {key: 'value'})
  let myNewMD = md.get().metadata.find(md => md.key === KEY)
}
```












Each flow type supports its own set of:

- Triggers
- Readable information
- Writable information
- Libraries provided by ZITADEL

### External Authentication

The external authentication flow 

ZITADEL supports only the external authentication flow at the moment.
[More flows are coming soon](https://zitadel.com/roadmap).

### External authentication flow triggers

- Post authentication: A user has authenticated externally. ZITADEL retrieved and mapped the external information.
- Pre creation:  A user selected **Register** on the overview page after external authentication. ZITADEL did not create the user yet.
- Post creation: A user selected **Register** on the overview page after external authentication. ZITADEL created the user.

### External authentication flow context

- `ctx.accessToken string`  
  This can be an opaque token or a JWT
- `ctx.idToken string`
- `ctx.getClaim(string) any`  
  Returns the requested claim
- `ctx.claimsJSON() object`  
  Returns the complete payload of the `ctx.idToken`

### External authentication flow api

- `api.setFirstName(string)`
- `api.setLastName(string)`
- `api.setNickName(string)`
- `api.setDisplayName(string)`
- `api.setPreferredLanguage(string)`
- `api.setGender(Gender)`  
- `api.setUsername(string)`  
  This function is only available for the pre creation trigger
- `api.setPreferredUsername(string)`  
  This function is only available for the post authentication trigger
- `api.setEmail(string)`
- `api.setEmailVerified(bool)`
- `api.setPhone(string)`
- `api.setPhoneVerified(bool)`
- `api.metadata array<Metadata>`  
  Push entries.  
- `api.userGrants array<UserGrant>`  
  Push entries.  
  This field is only available for the post creation trigger


### External authentication flow types <!-- TODO: Are these types correct? -->

- `Gender` is a code number

| code | gender |
| ---- | ------ |
| 0 | unspecified |
| 1 | female |
| 2 | male |
| 3 | diverse |

- `UserGrant` is a JavaScript object

```ts
{
    ProjectID: string,
    ProjectGrantID: string,
    Roles: Array<string>,
}
```

- `Metadata` is a JavaScript object with string values.
  The string values must be Base64 encoded

## Further reading

- [Actions concept](../concepts/features/actions)
- [Actions guide](../guides/manage/customize/behavior)
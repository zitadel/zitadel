---
title: Actions
---

This page describes the options you have when writing ZITADEL actions scripts.

## Language
ZITADEL interpretes the scripts as JavaScript.
Make sure your scripts are ECMAScript 5.1(+) compliant.
Go to the [goja GitHub page](https://github.com/dop251/goja) for detailed reference about the underlying library features and limitations.

Actions do not have access to any libraries yet.
Also, sending HTTP requests is not supported yet.
[We plan to add such features in the future](https://zitadel.ch/roadmap).

## Flows

Each flow type supports its own set of:
- Triggers
- Readable information
- Writable information

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

ZITADEL supports only the external authentication flow at the moment.
[More flows are coming soon](https://zitadel.ch/roadmap).

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

- [Actions concept](../../concepts/features/actions)
- [Actions guide](../../guides/customization/behavior)

---
title: Actions
---

This page describes how you can write ZITADEL actions scripts.

## Language
ZITADEL interpretes the scripts as JavaScript.
Make sure your scripts are ECMAScript 5.1(+) compliant.
Go to the [goja GitHub page](https://github.com/dop251/goja) for detailed reference about the underlying library features and limitations.

## Flows

Each flow type supports its own set of:
- Triggers
- Readable information
- Writable information

For reading and mutating state, each action has access to the JavaScript objects `ctx` and `api`.

The object `ctx` provides readable information as object properties and by callable functions.
The object `api` provides mutable properties and state mutating functions.

### External authentication flow

Triggers:
- Post authentication: A user has authenticated externally. ZITADEL retrieved and mapped the external information.
- Pre creation:  A user selected **Register** on the overview page after external authentication. ZITADEL did not create the user yet.
- Post creation: A user selected **Register** on the overview page after external authentication. ZITADEL created the user.

Readable user state:
- `ctx.accessToken string` This can be an opaque token or a JWT
- `ctx.idToken string`
- `ctx.getClaim(string) any`: Returns the requested claim
- `ctx.claimsJSON() object`: Returns the complete payload of the `ctx.idToken`

Writable user state:
- `api.setFirstName(string)`
- `api.setLastName(string)`
- `api.setNickName(string)`
- `api.setDisplayName(string)`
- `api.setPreferredLanguage(string)`
- `api.setGender(Gender)` <!-- TODO: What type is Gender? -->
- `api.setUsername(string)` This function is only available for the pre creation trigger
- `api.setPreferredUsername(string)` This function is only available for the post authentication trigger
- `api.setEmail(string)`
- `api.setEmailVerified(bool)`
- `api.setPhone(string)`
- `api.setPhoneVerified(bool)`
- `api.metadata array<Metadata>` Push entries. <!-- TODO: What type is Metadata? -->
- `api.userGrants array<UserGrant>` Push entries. This field is only available for the post creation trigger <!-- TODO: What type is UserGrant? -->

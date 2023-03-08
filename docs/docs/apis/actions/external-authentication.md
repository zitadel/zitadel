---
title: External Authentication Flow
---

This flow is executed if the user logs in using an [identity provider](/guides/integrate/identity-providers/introduction.md) or using a [jwt token](/concepts/structure/jwt_idp).

## Post Authentication

A user has authenticated externally. ZITADEL retrieved and mapped the external information.

### Parameters of Post Authentication Action

- `ctx`  
The first parameter contains the following fields
  - `accessToken` *string*  
    The access token which will be returned to the user. This can be an opaque token or a JWT
  - `claimsJSON()` [*idTokenClaims*](../openidoauth/claims)  
    Returns all claims of the id token
  - `getClaim(key)` *Any*  
    Returns the requested [id token claim](../openidoauth/claims)
  - `idToken` *string*  
    The id token which will be returned to the user
  - `v1`
    - `externalUser()` [*externalUser*](./objects#external-user)
    - `authError` *string*  
      This is a verification errors string representation. If the verification succeeds, this is "none"
    - `authRequest` [*auth request*](/docs/apis/actions/objects#auth-request)
    - `httpRequest` [*http request*](/docs/apis/actions/objects#http-request)
    - `providerInfo` *Any*  
      Returns the response of the provider. In case the provider is a Generic OAuth Provider, the information is accessible through:
      - `rawInfo`  *Any*
- `api`  
  The second parameter contains the following fields
  - `v1`
    - `user`
      - `appendMetadata(string, Any)`  
        The first parameter represents the key and the second a value which will be stored
  - `setFirstName(string)`  
    Sets the first name
  - `setLastName(string)`  
    Sets the last name
  - `setNickName(string)`  
    Sets the nickname
  - `setDisplayName(string)`  
    Sets the display name
  - `setPreferredLanguage(string)`  
    Sets the preferred language. Please use the format defined in [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646)
  - `setPreferredUsername(string)`  
    Sets the preferred username
  - `setEmail(string)`  
    Sets the email address of the user
  - `setEmailVerified(boolean)`  
    Sets the email address verified or unverified
  - `setPhone(string)`  
    Sets the phone number of the user
  - `setPhoneVerified(boolean)`  
    Sets the phone number verified or unverified
  - `metadata`  
    Array of [*metadata*](./objects#metadata-with-value-as-bytes). This function is deprecated, please use `api.v1.user.appendMetadata`

## Pre Creation

A user selected **Register** on the overview page after external authentication. ZITADEL did not create the user yet.

### Parameters of Pre Creation

- `ctx`  
  The first parameter contains the following fields
  - `v1`
    - `user` [*human*](./objects#human-user)
    - `authRequest` [*auth request*](/docs/apis/actions/objects#auth-request)
    - `httpRequest` [*http request*](/docs/apis/actions/objects#http-request)
- `api`  
  The second parameter contains the following fields
  - `metadata`  
    Array of [*metadata*](./objects#metadata-with-value-as-bytes). This function is deprecated, please use `api.v1.user.appendMetadata`
  - `setFirstName(string)`  
    Sets the first name
  - `setLastName(string)`  
    Sets the last name
  - `setNickName(string)`  
    Sets the nick name
  - `setDisplayName(string)`  
    Sets the display name
  - `setPreferredLanguage(string)`  
    Sets the preferred language, the string has to be a valid language tag as defined in [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646)
  - `setGender(int)`  
    Sets the gender.  
    <ul><li>0: unspecified</li><li>1: female</li><li>2: male</li><li>3: diverse</li></ul>
  - `setUsername(string)`  
    Sets the username
  - `setEmail(string)`  
    Sets the email
  - `setEmailVerified(bool)`  
    If true the email set is verified without user interaction
  - `setPhone(string)`  
    Sets the phone number
  - `setPhoneVerified(bool)`  
    If true the phone number set is verified without user interaction
  - `v1`
    - `user`
      - `appendMetadata(string, Any)`  
        The first parameter represents the key and the second a value which will be stored

## Post Creation

A user selected **Register** on the overview page after external authentication and ZITADEL successfully created the user.

### Parameters of Post Creation

- `ctx`  
  The first parameter contains the following fields
  - `v1`
    - `getUser()` [*user*](./objects#user)
    - `authRequest` [*auth request*](/docs/apis/actions/objects#auth-request)
    - `httpRequest` [*http request*](/docs/apis/actions/objects#http-request)
- `api`  
  The second parameter contains the following fields
  - `userGrants` Array of [*userGrant*](./objects#user-grant)'s
  - `v1`
    - `appendUserGrant(`[`userGrant`](./objects#user-grant)`)`

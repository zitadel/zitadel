---
title: Internal Authentication Flow
---

## Post Authentication

A user has authenticated directly at ZITADEL.
ZITADEL validated the users inputs for password, one-time password, security key or passwordless factor.
Each validation step triggers the action.

### Parameters of Post Authentication Action

- `ctx`  
  The first parameter contains the following fields
    - `v1`
        - `authMethod` *string*  
          This is one of "password", "OTP", "U2F" or "passwordless"
        - `authError` *string*  
          This is a verification errors string representation. If the verification succeeds, this is "none"
        - `authRequest` [*auth request*](/docs/apis/actions/objects#auth-request)
        - `httpRequest` [*http request*](/docs/apis/actions/objects#http-request)
- `api`  
  The second parameter contains the following fields
    - `metadata`  
      Array of [*metadata*](./objects#metadata-with-value-as-bytes). This function is deprecated, please use `api.v1.user.appendMetadata`
    - `v1`
        - `user`
            - `appendMetadata(string, Any)`  
              The first parameter represents the key and the second a value which will be stored

## Pre Creation

A user registers directly at ZITADEL.
ZITADEL did not create the user yet.

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

A user registers directly at ZITADEL.  
ZITADEL successfully created the user.

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

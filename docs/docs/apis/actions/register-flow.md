---
title: Register flows
---

## External Authentication

<!-- link idp and jwt -->
This flow is executed if the user logs in using an identity provider or using a jwt token.

### Post Authentication

A user has authenticated externally. ZITADEL retrieved and mapped the external information.

#### Parameters of post authentication action

- `ctx`: The first parameter contains the following fields:
  - `accessToken`: `string` The access token which will be returned to the user. This can be an opaque token or a JWT
  - `claimsJSON()` [`idTokenClaims`](./objects#id-token-claims): returns all claims of the id token
  - `getClaim(key) Object)`: returns the requested [id token claim](./objects#id-token-claims)
  - `idToken`: `string` The id token which will be returned to the user
  - `v1`:
    - `externalUser()`: [`externalUser`](./objects#external-user)
- `api`: The second parameter contains the following fields:
  - `v1`
    - `user`
      - `appendMetadata(string, Object)`: the first parameter represents the key and the second a value which will be stored
  - `setFirstName(string)`: sets the parameter as first name
  - `setLastName(string)`: sets the parameter as last name
  - `setNickName(string)`: sets the parameter as nickname
  - `setDisplayName(string)`: sets the parameter as display name
  - `setPreferredLanguage(string)`: sets the parameter as preferred language. Please use the format defined in [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646)
  - `setPreferredUsername(string)`: sets the preferred username
  - `setEmail(string)`: sets the email address of the user
  - `setEmailVerified(boolean)`: sets the email address verified or unverified
  - `setPhone(string)`: sets the phone number of the user
  - `setPhoneVerified(boolean)`: sets the phone number verified or unverified
  - `metadata` Array of [`metadata`](./objects#metadata). This function is deprecated, please use `api.v1.user.appendMetadata`

### Pre Creation

A user selected **Register** on the overview page after external authentication. ZITADEL did not create the user yet.

#### Parameters of Pre Creation

`ctx`:

**Fields:**

| name | description | type |
|---|---|---|
| accessToken | the access token which will be returned to the user. This can be an opaque token or a JWT | `string` |
| idToken | the id token which will be returned to the user | `string` |

**Methods:**

| name | description | return value |
|---|---|---|
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
| setEmail(string) | sets the email | none |
| setEmailVerified(bool) | if true the email set is verified without user interaction | none |
| setPhone(string) | sets the phone number | none |
| setPhoneVerified(bool) | if true the phone number set is verified without user interaction | none |

### Post Creation

A user selected **Register** on the overview page after external authentication and ZITADEL successfully created the user.

#### Parameters of Post Creation

`ctx`:

**Fields:**

| name | description | type |
|---|---|---|
| accessToken | the access token which will be returned to the user. This can be an opaque token or a JWT | `string` |
| idToken | the id token which will be returned to the user | `string` |

**Methods:**

| name | description | return value |
|---|---|---|
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
| setEmail(string) | sets the email | none |
| setEmailVerified(bool) | if true the email set is verified without user interaction | none |
| setPhone(string) | sets the phone number | none |
| setPhoneVerified(bool) | if true the phone number set is verified without user interaction | none |

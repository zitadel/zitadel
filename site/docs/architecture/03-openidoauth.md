---
title: OpenID Connect 1.0 & OAuth 2.0
---

### Endpoints and Domains
This chapter documents the [OpenID Connect 1.0](https://openid.net/connect/) and [OAuth 2.0](https://oauth.net/2/) features provided by **ZITADEL**.

Under normal circumstances **ZITADEL** need four domain names to operate properly.

| Domain Name | Example               | Description                                                                                                                          |
|:------------|:----------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| issuer      | `issuer.zitadel.ch`   | Provides the [OpenID Connect 1.0 Discovery Endpoint](#OpenID_Connect_1_0_Discovery)                                                   |
| api         | `api.zitadel.ch`      | All ZITADEL API's are located under this domain see [API explanation](apis#APIs) for details                                      |
| login       | `accounts.zitadel.ch` | The accounts.* page provides server renderer pages like login and register and as well the authorization_endpoint for OpenID Connect |
| console     | `console.zitadel.ch`  | With the console.* domain we serve the assets for the management gui                                                                 |

#### OpenID Connect 1.0 Discovery

The OpenID Connect Discovery Endpoint is located within the issuer domain.
For example with [zitadel.ch](https://zitadel.ch), issuer.zitadel.ch would be the domain. This would give us [https://issuer.zitadel.ch/.well-known/openid-configuration](https://issuer.zitadel.ch/.well-known/openid-configuration).

**Link to spec.** [OpenID Connect Discovery 1.0 incorporating errata set 1](https://openid.net/specs/openid-connect-discovery-1_0.html)

#### authorization_endpoint

[https://accounts.zitadel.ch/oauth/v2/authorize](https://accounts.zitadel.ch/oauth/v2/authorize)

> The authorization_endpoint is located with the login page, due to the need of accessing the same cookie domain

Required request Parameters

| Parameter     | Description                                                                                                                                       |
|---------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| client_id     | The id of your client as shown in Console.                                                                                                        |
| redirect_uri  | Callback uri of the authorization request where the code or tokens will be sent to. Must match exactly one of the preregistered in Console.       |
| response_type | Determines whether a `code`, `id_token token` or just `id_token` will be returned. Most use cases will need `code`. See flow guide for more info. |
| scope         | `openid` is required, see [Scopes](architecture#Scopes) for more possible values. Scopes are space delimited, e.g. `openid email profile`                     |

Required parameters for PKCE (see PKCE guide for more information)

| Parameter             | Description                                           |
|-----------------------|-------------------------------------------------------|
| code_challenge        | The SHA-256 value of the generated code_verifier      | 
| code_challenge_method | Method used to generate the challenge, must be `S256` |

Optional parameters

| Parameter     | Description                                                                                                                              |
|---------------|------------------------------------------------------------------------------------------------------------------------------------------|
| id_token_hint | Valid `id_token` (of an existing session) used to identity the subject. Should be provided when using prompt `none`.                     | 
| login_hint    | A valid logon name of a user. Will be used for username inputs or preselecting a user on `select_account`                                |
| max_age       | | 
| nonce         | Random string value to associate the client session with the ID Token and for replay attacks mitigation.                                 | 
| prompt        | If the Auth Server prompts the user for (re)authentication. <br>no prompt: the user will have to choose a session if more than one session exists<br>`none`: user must be authenticated without interaction, an error is returned otherwise <br>`login`: user must reauthenticate / provide a user name <br>`select_account`: user is prompted to select one of the existing sessions or create a new one |
| state         | Opaque value used to maintain state between the request and the callback. Used for Cross-Site Request Forgery (CSRF) mitigation as well. |

Successful Code Response

| Property | Description                                                                   |
|----------|-------------------------------------------------------------------------------| 
| code     | Opaque string which will be necessary to request tokens on the token endpoint |
| state    | Unmodified `state` parameter from the request                                 |

Successful Implicit Response

| Property     | Description                                                 |
|--------------|-------------------------------------------------------------| 
| access_token | Only returned if `response_type` included `token`           |
| expires_in   | Number of second until the expiration of the `access_token` |
| id_token     | Only returned if `response_type` included `id_token`        |
| token_type   | Type of the `access_token`. Value is always `Bearer`        |

Error Response

Regardless of the authorization flow chosen, if an error occurs the following response will be returned to the redirect_uri.

> If the redirect_uri is not provided, was not registered or anything other prevents the auth server form returning the response to the client,
the error will be display directly to the user on the auth server


| Property          | Description                                                          |
|-------------------|----------------------------------------------------------------------| 
| error             | An OAuth / OIDC error_type  (//TODO: list error types)               |
| error_description | Description of the error type or additional information of the error |
| state             | Unmodified `state` parameter from the request                        |

#### token_endpoint

[https://api.zitadel.ch/oauth/v2/token](https://api.zitadel.ch/oauth/v2/token)

##### Authorization Code Grant (Code Exchange)

Required request Parameters

| Parameter     | Description                                                                                                   |
|---------------|---------------------------------------------------------------------------------------------------------------|
| code          | Code that was issued from the authorization request.                                                          |
| grant_type    | must be `authorization_code`
| redirect_uri  | Callback uri where the code was be sent to. Must match exactly the redirect_uri of the authorization request. |

Depending on your authorization method you will have to provide additional parameters or headers:

When using `client_secret_basic`

Send your `client_id` and `client_secret` as Basic Auth Header in the following manner:

```markdown
Authorization: "Basic " + base64( formUrlEncode(client_id) + ":" + formUrlEncode(client_secret) )
```

Given the client_id `78366401571920522@amce` and client_secret `veryweaksecret!`, this would result in the following `Authorization` header: 
`Basic NzgzNjY0MDE1NzE5MjA1MjIlNDBhbWNlOnZlcnl3ZWFrc2VjcmV0JTIx`

When using `client_secret_post`

Send your `client_id` and `client_secret` as parameters in the body:

| Parameter     | Description                      |
|---------------|----------------------------------|
| client_id     | client_id of the application     |
| client_secret | client_secret of the application |

When using `none` (PKCE)

Send your code_verifier for us to recompute the code_challenge of the authorization request.

| Parameter     | Description                                                  |
|---------------|--------------------------------------------------------------|
| code_verifier | code_verifier previously used to generate the code_challenge |

When using `private_key_jwt`

Send a client assertion as JWT for us to validate the signature against the registered public key.

| Parameter             | Description                                                                                                     |
|-----------------------|-----------------------------------------------------------------------------------------------------------------|
| client_assertion      | JWT built and signed according to [Using JWTs for Client Authentication](#Using JWTs for Client Authentication) |
| client_assertion_type | must be `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`                                                |

##### JWT Profile Grant

> TODO: describe or link

#### introspection_endpoint

[https://api.zitadel.ch/oauth/v2/introspection](https://api.zitadel.ch/oauth/v2/introspection)


| Parameter | Description     |
|-----------|-----------------|
| token     | An access token |

Depending on your authorization method you will have to provide additional parameters or headers:

When using `client_secret_basic`

Send your `client_id` and `client_secret` as Basic Auth Header in the following manner:

```markdown
Authorization: "Basic " + base64( formUrlEncode(client_id) + ":" + formUrlEncode(client_secret) )
```

Given the client_id `78366401571920522@amce` and client_secret `veryweaksecret!`, this would result in the following `Authorization` header:
`Basic NzgzNjY0MDE1NzE5MjA1MjIlNDBhbWNlOnZlcnl3ZWFrc2VjcmV0JTIx`

When using `private_key_jwt`

Send a client assertion as JWT for us to validate the signature against the registered public key.

| Parameter             | Description                                                                                                     |
|-----------------------|-----------------------------------------------------------------------------------------------------------------|
| client_assertion      | JWT built and signed according to [Using JWTs for Client Authentication](#Using JWTs for Client Authentication) |
| client_assertion_type | must be `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`                                                |

#### userinfo_endpoint

[https://api.zitadel.ch/oauth/v2/userinfo](https://api.zitadel.ch/oauth/v2/userinfo)

#### end_session_endpoint

[https://accounts.zitadel.ch/oauth/v2/endsession](https://accounts.zitadel.ch/oauth/v2/endsession)

> The end_session_endpoint is located with the login page, due to the need of accessing the same cookie domain

#### jwks_uri

[https://api.zitadel.ch/oauth/v2/keys](https://api.zitadel.ch/oauth/v2/keys)

> Be aware that these keys can be rotated without any prior notice. We will however make sure that a proper `kid` is set with each key!

#### OAuth 2.0 Metadata

**ZITADEL** does not yet provide a OAuth 2.0 Metadata endpoint but instead provides a [OpenID Connect Discovery Endpoint](#OpenID_Connect_1_0_Discovery).

### Scopes

ZITADEL supports the usage of scopes as way of requesting information from the IAM and also instruct ZITADEL to do certain operations.

#### Standard Scopes

| Scopes  | Example   | Description                                          |
|:--------|:----------|------------------------------------------------------|
| openid  | `openid`  | When using openid connect this is a mandatory scope  |
| profile | `profile` | Optional scope to request the profile of the subject |
| email   | `email`   | Optional scope to request the email of the subject   |
| address | `address` | Optional scope to request the address of the subject |

#### Custom Scopes

> This feature is not yet released

#### Reserved Scopes

In addition to the standard compliant scopes we utilize the following scopes.

| Scopes                                          | Example                                                                        | Description                                                                                                                                                                                                                                                                                                                                                             |
|:------------------------------------------------|:-------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| urn:zitadel:iam:org:project:role:{rolename}     | `urn:zitadel:iam:org:project:role:user`                                        | By using this scope a [client](administrate#clients) can request the claim urn:zitadel:iam:roles:rolename} to be asserted when possible. As an alternative approach you can enable all [roles](administrate#Roles) to be asserted from the [project](administrate#projects) a [client](administrate#clients) belongs to. See details [here](administrate#RBAC_Settings) |
| urn:zitadel:iam:org:domain:primary:{domainname} | `urn:zitadel:iam:org:domain:primary:acme.ch`                                   | When requesting this scope **ZITADEL** will enforce that the user is a member of the selected organization. If the organization does not exist a failure is displayed                                                                                                                                                                                                   |
| urn:zitadel:iam:role:{rolename}                 |                                                                                |                                                                                                                                                                                                                                                                                                                                                                         |
| urn:zitadel:iam:org:project:id:{projectid}:aud  | ZITADEL's Project id is `urn:zitadel:iam:org:project:id:69234237810729019:aud` | By adding this scope, the requested projectid will be added to the audience of the access and id token                                                                                                                                                                                                                                                                  |

> If access to ZITADEL's API's is needed with a service user the scope `urn:zitadel:iam:org:project:id:69234237810729019:aud` needs to be used with the JWT Profile request

### Claims

ZITADEL asserts claims on different places according to the corresponding specifications or project and clients settings.
Please check below the matrix for an overview where which scope is asserted.

> Some claims will only be returned if certain scopes were requested (e.g. `address`), custom scopes are marked with `*`.
> See [Reserved Scopes](architecture#Reserved_Scopes) for details.
> 
> Additionally, some are only returned if response_type is only `id_token` 🆔 or if configured in Console ⚙.
> 
> Scopes in Access Tokens can only be asserted if type is JWT <img src="tech/jwt.png" alt="jwt icon">

| Claims                                          | Userinfo          | Introspect           | ID Token             | Access Token                                             |
|:------------------------------------------------|:------------------|----------------------|----------------------|----------------------------------------------------------|
| acr                                             | ✅                | ✅                   | ✅                   | ❌                                                       |
| address                                         | `address`         | `address`            | `address` and 🆔 / ⚙ | ❌                                                       |
| amr                                             | ✅                | ✅                   | ❌                                                       |
| aud                                             | ❌                | ❌                   | ✅                   | <img src="tech/jwt.png" alt="jwt">                       |
| auth_time                                       | ❌                | ❌                   | ✅                   | ❌                                                       |
| azp                                             | ❌                | ❌                   | ✅                   | <img src="tech/jwt.png" alt="jwt">                       |
| email                                           | `email`           | `email`              | `email` and 🆔 / ⚙   | ❌                                                       |
| email_verified                                  | `email`           | `email`              | `email` and 🆔 / ⚙   | ❌                                                       |
| exp                                             | ❌                | ❌                   | ✅                   | <img src="tech/jwt.png" alt="jwt">                       |
| family_name                                     | `profile`         | `profile`            | `profile` and 🆔 / ⚙ | ❌                                                       |
| gender                                          | `profile`         | `profile`            | `profile` and 🆔 / ⚙ | ❌                                                       |
| given_name                                      | `profile`         | `profile`            | `profile` and 🆔 / ⚙ | ❌                                                       |
| iat                                             | ❌                | ❌                   | ✅                   | <img src="tech/jwt.png" alt="jwt">                       |
| iss                                             | ❌                | ❌                   | ✅                   | <img src="tech/jwt.png" alt="jwt">                       |
| locale                                          | `profile`         | `profile`            | `profile` and 🆔 / ⚙ | ❌                                                       |
| name                                            | `profile`         | `profile`            | `profile` and 🆔 / ⚙ | ❌                                                       |
| nonce                                           | ❌                | ❌                   | ✅                   | ❌                                                       |
| phone                                           | `phone`           | `phone`              | `phone` and 🆔 / ⚙   | ❌                                                       |
| phone_verified                                  | `phone`           | `phone`              | `phone` and 🆔 / ⚙   | ❌                                                       |
| preferred_username (username when Introspect )  | `profile`         | `profile`            | ✅                   | ❌                                                       |
| sub                                             | ✅                | ✅                   | ✅                   | <img src="tech/jwt.png" alt="jwt">                       |
| urn:zitadel:iam:org:domain:primary:{domainname} | `Primary Domain*` | `Primary Domain*`    | `Primary Domain*`    | <img src="tech/jwt.png" alt="jwt"> and `Primary Domain*` |
| urn:zitadel:iam:org:project:roles:{rolename}    | `Roles*` / ⚙      | `Roles*` / ⚙         | `Roles*` / ⚙         | <img src="tech/jwt.png" alt="jwt"> and `Roles*` / ⚙      |

#### Standard Claims

| Claims             | Example                                  | Description                                                                                   |
|:-------------------|:-----------------------------------------|-----------------------------------------------------------------------------------------------|
| acr                | TBA                                      | TBA                                                                                           |
| address            | `Teufener Strasse 19, 9000 St. Gallen`   | TBA                                                                                           |
| amr                | `pwd mfa`                                | Authentication Method References as defined in [RFC8176](https://tools.ietf.org/html/rfc8176) |
| aud                | `69234237810729019`                      | By default all client id's and the project id is included                                     |
| auth_time          | `1311280969`                             | Unix time of the authentication                                                               |
| azp                | `69234237810729234`                      | Client id of the client who requested the token                                               |
| email              | `road.runner@acme.ch`                    | Email Address of the subject                                                                  |
| email_verified     | `true`                                   | Boolean if the email was verified by ZITADEL                                                  |
| exp                | `1311281970`                             | Time the token expires as unix time                                                           |
| family_name        | `Runner`                                 | The subjects family name                                                                      |
| gender             | `other`                                  | Gender of the subject                                                                         |
| given_name         | `Road`                                   | Given name of the subject                                                                     |
| iat                | `1311280970`                             | Issued at time of the token as unix time                                                      |
| iss                | `https://issuer.zitadel.ch`              | Issuing domain of a token                                                                     |
| locale             | `en`                                     | Language from the subject                                                                     |
| name               | `Road Runner`                            | The subjects full name                                                                        |
| nonce              | `blQtVEJHNTF0WHhFQmhqZ0RqeHJsdzdkd2d...` | The nonce provided by the client                                                              |
| phone              | `+41 79 XXX XX XX`                       | Phone number provided by the user                                                             |
| preferred_username | `road.runner@acme.caos.ch`               | ZITADEL's login name of the user. Consist of `username@primarydomain`                         |
| sub                | `77776025198584418`                      | Subject ID of the user                                                                        |

#### Custom Claims

> This feature is not yet released

#### Reserved Claims

ZITADEL reserves some claims to assert certain data.

| Claims                                          | Example                                                                                              | Description                                                                                                                                                                        |
|:------------------------------------------------|:-----------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| urn:zitadel:iam:org:domain:primary:{domainname} | `{"urn:zitadel:iam:org:domain:primary": "acme.ch"}`                                                  | This claim represents the primary domain of the organization the user belongs to.                                                                                                  |
| urn:zitadel:iam:org:project:roles:{rolename}    | `{"urn:zitadel:iam:org:project:roles": [ {"user": {"id1": "acme.zitade.ch", "id2": "caos.ch"} } ] }` | When roles are asserted, ZITADEL does this by providing the `id` and `primaryDomain` below the role. This gives you the option to check in which organization a user has the role. |
| urn:zitadel:iam:roles:{rolename}                | TBA                                                                                                  | TBA                                                                                                                                                                                |

### Grant Types

For a list of supported or unsupported `Grant Types` please have a look at the table below.

| Grant Type                                            | Supported           |
|:------------------------------------------------------|:--------------------|
| Authorization Code                                    | yes                 |
| Authorization Code with PKCE                          | yes                 |
| Client Credentials                                    | yes                 |
| Device Authorization                                  | under consideration |
| Implicit                                              | yes                 |
| JSON Web Token (JWT) Profile                          | yes                 |
| Refresh Token                                         | work in progress    |
| Resource Owner Password Credentials                   | no                  |
| Security Assertion Markup Language (SAML) 2.0 Profile | no                  |
| Token Exchange                                        | work in progress    |

#### Authorization Code

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.1](https://tools.ietf.org/html/rfc6749#section-1.3.1)

#### Proof Key for Code Exchange

**Link to spec.** [Proof Key for Code Exchange by OAuth Public Clients](https://tools.ietf.org/html/rfc7636)

#### Implicit

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.2](https://tools.ietf.org/html/rfc6749#section-1.3.2)

#### Client Credentials

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.4](https://tools.ietf.org/html/rfc6749#section-1.3.4)

#### Refresh Token

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.5](https://tools.ietf.org/html/rfc6749#section-1.5)

#### JSON Web Token (JWT) Profile

**Link to spec.** [JSON Web Token (JWT) Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7523)

##### Using JWTs as Authorization Grants

Our service user work with the JWT profile to authenticate them against ZITADEL.

1. Create or use an existing service user
2. Create a new key and download it
3. Generate a JWT with the structure below and sign it with the downloaded key
4. Send the JWT Base64 encoded to ZITADEL's token endpoint
5. Use the received access token

---

Key JSON

| Key    | Example                                                             | Description                                                        |
|:-------|:--------------------------------------------------------------------|:-------------------------------------------------------------------|
| type   | `"serviceaccount"`                                                  | The type of account, right now only serviceaccount is valid        |
| keyId  | `"81693565968772648"`                                               | This is unique ID of the key                                       |
| key    | `"-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----"` | The private key generated by ZITADEL, this can not be regenerated! |
| userId | `78366401571647008`                                                 | The service users ID, this is the same as the subject from tokens  |

```JSON
{
	"type": "serviceaccount",
	"keyId": "81693565968772648",
	"key": "-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----",
	"userId": "78366401571647008"
}
```

---

JWT

| Claim | Example                       | Description                                                                                                   |
|:------|:------------------------------|:--------------------------------------------------------------------------------------------------------------|
| aud   | `"https://issuer.zitadel.ch"` | String or Array of intended audiences MUST include ZITADEL's issuing domain                                   |
| exp   | `1605183582`                  | Unix timestamp of the expiry, MUST NOT be longer than 1h                                                      |
| iat   | `1605179982`                  | Unix timestamp of the creation singing time of the JWT                                                        |
| iss   | `"77479219772321307"`         | String which represents the requesting party (owner of the key), normally the `userId` from the json key file |
| sub   | `"77479219772321307"`         | The subject ID of the service user, normally the `userId` from the json key file                              |

```JSON
{
	"iss": "77479219772321307",
	"sub": "77479219772321307",
	"aud": "https://issuer.zitadel.ch",
	"exp": 1605183582,
	"iat": 1605179982
}
```

> To identify your key, it is necessary that you provide a JWT with a `kid` header claim representing your keyId from the Key JSON:
> ```json
> {
> 	"alg": "RS256",
> 	"kid": "81693565968772648"
> }
> ```

---

Access Token Request

| Parameter    | Example                                                                     | Description                                   |
|:-------------|:----------------------------------------------------------------------------|:----------------------------------------------|
| Content-Type | `application/x-www-form-urlencoded`                                         |                                               |
| grant_type   | `urn:ietf:params:oauth:grant-type:jwt-bearer`                               | Using JWTs as Authorization Grants            |
| assertion    | `eyJhbGciOiJSUzI1Ni...`                                                     | The base64 encoded JWT created above          |
| scope        | `openid profile email urn:zitadel:iam:org:project:id:69234237810729019:aud` | Scopes you would like to request from ZITADEL |

```BASH
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data assertion=eyJhbGciOiJSUzI1Ni... \
  --data scope=openid profile email address
```

##### Using JWTs for Client Authentication

JWT can be used for Client Authentication for Code Exchange as well as Token Introspection //TODO: links

1. Create or use an existing app (OIDC or API)
2. Create a new key and download it
3. Generate a JWT with the structure below and sign it with the downloaded key
4. Use the JWT Base64 encoded in Code Exchange or Token Introspection request

---

Key JSON

| Key      | Example                                                             | Description                                                                    |
|:---------|:--------------------------------------------------------------------|:-------------------------------------------------------------------------------|
| type     | `"application"`                                                     | The type of account, right now only application is valid                       |
| keyId    | `"81693565968962154"`                                               | This is unique ID of the key                                                   |
| key      | `"-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----"` | The private key generated by ZITADEL, this can not be regenerated!             |
| clientId | `78366401571920522@acme`                                            | The client_id of the application, this is the same as the subject from tokens  |
| appId    | `78366403256846242`                                                 | The id of the application (just for completeness, not used for JWT)            |

```JSON
{
	"type": "serviceaccount",
	"keyId": "81693565968962154",
	"key": "-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----",
	"clientId": "78366401571920522@acme",
	"appId": "78366403256846242"
}
```

---

JWT

| Claim | Example                       | Description                                                                                                     |
|:------|:------------------------------|:----------------------------------------------------------------------------------------------------------------|
| aud   | `"https://issuer.zitadel.ch"` | String or Array of intended audiences MUST include ZITADEL's issuing domain                                     |
| exp   | `1605183582`                  | Unix timestamp of the expiry, MUST NOT be longer than 1h                                                        |
| iat   | `1605179982`                  | Unix timestamp of the creation singing time of the JWT                                                          |
| iss   | `"78366401571920522@acme"`    | String which represents the requesting party (owner of the key), normally the `clientID` from the json key file |
| sub   | `"78366401571920522@acme"`    | The subject ID of the application, normally the `clientID` from the json key file                               |

```JSON
{
	"iss": "78366401571920522@acme",
	"sub": "78366401571920522@acme",
	"aud": "https://issuer.zitadel.ch",
	"exp": 1605183582,
	"iat": 1605179982
}
```

> To identify your key, it is necessary that you provide a JWT with a `kid` header claim representing your keyId from the Key JSON:
> ```json
> {
> 	"alg": "RS256",
> 	"kid": "81693565968962154"
> }
> ```

---

Access Token Request

| Parameter             | Example                                                      | Description                                   |
|:----------------------|:-------------------------------------------------------------|:----------------------------------------------|
| Content-Type          | `application/x-www-form-urlencoded`                          |                                               |
| grant_type            | `authorization_code`                                         | Using JWTs as Client Authentication           |
| code                  | `DKLvnksjndjsflkdjlkfgjslow...`                              | The code you received from the Auth Endpoint  |
| client_assertion_type | `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`     | Using JWTs as Client Authentication           |
| client_assertion      | `eyJhbGciOiJSUzI1Ni...`                                      | The base64 encoded JWT created above          |

```BASH
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=authorization_code \
  --data code=DKLvnksjndjsflkdjlkfgjslow... \
  --data client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer \
  --data client_assertion=eyJhbGciOiJSUzI1Ni...
```

---

Introspection Request

| Parameter             | Example                                                      | Description                                   |
|:----------------------|:-------------------------------------------------------------|:----------------------------------------------|
| Content-Type          | `application/x-www-form-urlencoded`                          |                                               |
| client_assertion_type | `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`     | Using JWTs as Client Authentication           |
| client_assertion      | `eyJhbGciOiJSUzI1Ni...`                                      | The base64 encoded JWT created above          |
| token                 | `VjVxyCZmRmWYqd3_F5db9Pb9mHR5fqzhn...`                       | The (access) token you would like to check    |

```BASH
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer \
  --data client_assertion=eyJhbGciOiJSUzI1Ni... \
  --data token=VjVxyCZmRmWYqd3_F5db9Pb9mHR5fqzhn...
```


#### Token Exchange

**Link to spec.** [OAuth 2.0 Token Exchange](https://tools.ietf.org/html/rfc8693)

### Device Authorization

**Link to spec.** [OAuth 2.0 Device Authorization Grant](https://tools.ietf.org/html/rfc8628)

### Not Supported Grant Types

#### Resource Owner Password Credentials

> Due to growing security concerns we do not support this grant type. With OAuth 2.1 it looks like this grant will be removed.

**Link to spec.** [OThe OAuth 2.0 Authorization Framework Section 1.3.3](https://tools.ietf.org/html/rfc6749#section-1.3.3)

#### Security Assertion Markup Language (SAML) 2.0 Profile

**Link to spec.** [Security Assertion Markup Language (SAML) 2.0 Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7522)


---
title: Endpoints
---

## OpenID Connect 1.0 Discovery

The OpenID Connect Discovery Endpoint is located within the issuer domain.
For example with [zitadel.ch](https://zitadel.ch), issuer.zitadel.ch would be the domain. This would give us [https://issuer.zitadel.ch/.well-known/openid-configuration](https://issuer.zitadel.ch/.well-known/openid-configuration).

**Link to spec.** [OpenID Connect Discovery 1.0 incorporating errata set 1](https://openid.net/specs/openid-connect-discovery-1_0.html)

## authorization_endpoint

[https://accounts.zitadel.ch/oauth/v2/authorize](https://accounts.zitadel.ch/oauth/v2/authorize)

> The authorization_endpoint is located with the login page, due to the need of accessing the same cookie domain

Required request Parameters

| Parameter     | Description                                                                                                                                       |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| client_id     | The id of your client as shown in Console.                                                                                                        |
| redirect_uri  | Callback uri of the authorization request where the code or tokens will be sent to. Must match exactly one of the preregistered in Console.       |
| response_type | Determines whether a `code`, `id_token token` or just `id_token` will be returned. Most use cases will need `code`. See flow guide for more info. |
| scope         | `openid` is required, see [Scopes](architecture#Scopes) for more possible values. Scopes are space delimited, e.g. `openid email profile`         |

Required parameters for PKCE (see PKCE guide for more information)

| Parameter             | Description                                           |
| --------------------- | ----------------------------------------------------- |
| code_challenge        | The SHA-256 value of the generated code_verifier      |
| code_challenge_method | Method used to generate the challenge, must be `S256` |

Optional parameters

| Parameter     | Description                                                                                                                                                                                                                                                                                                                                                                                                       |
| ------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| id_token_hint | Valid `id_token` (of an existing session) used to identity the subject. Should be provided when using prompt `none`.                                                                                                                                                                                                                                                                                              |
| login_hint    | A valid logon name of a user. Will be used for username inputs or preselecting a user on `select_account`                                                                                                                                                                                                                                                                                                         |
| max_age       | Seconds since the last active successful authentication of the user                                                                                                                                                                                                                                                                                                                                               |
| nonce         | Random string value to associate the client session with the ID Token and for replay attacks mitigation.                                                                                                                                                                                                                                                                                                          |
| prompt        | If the Auth Server prompts the user for (re)authentication. <br />no prompt: the user will have to choose a session if more than one session exists<br />`none`: user must be authenticated without interaction, an error is returned otherwise <br />`login`: user must reauthenticate / provide a user name <br />`select_account`: user is prompted to select one of the existing sessions or create a new one |
| state         | Opaque value used to maintain state between the request and the callback. Used for Cross-Site Request Forgery (CSRF) mitigation as well.                                                                                                                                                                                                                                                                          |

Successful Code Response

| Property | Description                                                                   |
| -------- | ----------------------------------------------------------------------------- |
| code     | Opaque string which will be necessary to request tokens on the token endpoint |
| state    | Unmodified `state` parameter from the request                                 |

Successful Implicit Response

| Property     | Description                                                 |
| ------------ | ----------------------------------------------------------- |
| access_token | Only returned if `response_type` included `token`           |
| expires_in   | Number of second until the expiration of the `access_token` |
| id_token     | Only returned if `response_type` included `id_token`        |
| token_type   | Type of the `access_token`. Value is always `Bearer`        |

Error Response

Regardless of the authorization flow chosen, if an error occurs the following response will be returned to the redirect_uri.

> If the redirect_uri is not provided, was not registered or anything other prevents the auth server form returning the response to the client,
the error will be display directly to the user on the auth server

| Property          | Description                                                          |
| ----------------- | -------------------------------------------------------------------- |
| error             | An OAuth / OIDC error_type                                           |
| error_description | Description of the error type or additional information of the error |
| state             | Unmodified `state` parameter from the request                        |

## token_endpoint

[https://api.zitadel.ch/oauth/v2/token](https://api.zitadel.ch/oauth/v2/token)

### Authorization Code Grant (Code Exchange)

Required request Parameters

| Parameter    | Description                                                                                                   |
| ------------ | ------------------------------------------------------------------------------------------------------------- |
| code         | Code that was issued from the authorization request.                                                          |
| grant_type   | Must be `authorization_code`                                                                                  |
| redirect_uri | Callback uri where the code was be sent to. Must match exactly the redirect_uri of the authorization request. |

Depending on your authorization method you will have to provide additional parameters or headers:

When using `client_secret_basic`

Send your `client_id` and `client_secret` as Basic Auth Header. Check [Client Secret Basic Auth Method](architecture#Client_Secret_Basic) on how to build it correctly.

When using `client_secret_post`

Send your `client_id` and `client_secret` as parameters in the body:

| Parameter     | Description                      |
| ------------- | -------------------------------- |
| client_id     | client_id of the application     |
| client_secret | client_secret of the application |

When using `none` (PKCE)

Send your code_verifier for us to recompute the code_challenge of the authorization request.

| Parameter     | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| code_verifier | code_verifier previously used to generate the code_challenge |

When using `private_key_jwt`

Send a client assertion as JWT for us to validate the signature against the registered public key.

| Parameter             | Description                                                                                                     |
| --------------------- | --------------------------------------------------------------------------------------------------------------- |
| client_assertion      | JWT built and signed according to [Using JWTs for Client Authentication](#Using JWTs for Client Authentication) |
| client_assertion_type | Must be `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`                                                |

### JWT Profile Grant

---

Required request Parameters

| Parameter  | Description                                                                                                                   |
| ---------- | ----------------------------------------------------------------------------------------------------------------------------- |
| grant_type | Must be `urn:ietf:params:oauth:grant-type:jwt-bearer`                                                                         |
| assertion  | JWT built and signed according to [Using JWTs for Client Authentication](#Using JWTs for Client Authentication)               |
| scope      | [Scopes](architecture#Scopes) you would like to request from ZITADEL. Scopes are space delimited, e.g. `openid email profile` |

```BASH
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=authorization_code \
  --data code=DKLvnksjndjsflkdjlkfgjslow... \
  --data client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer \
  --data client_assertion=eyJhbGciOiJSUzI1Ni...
```

## introspection_endpoint

[https://api.zitadel.ch/oauth/v2/introspect](https://api.zitadel.ch/oauth/v2/introspect)

| Parameter | Description     |
| --------- | --------------- |
| token     | An access token |

Depending on your authorization method you will have to provide additional parameters or headers:

When using `client_secret_basic`

Send your `client_id` and `client_secret` as Basic Auth Header. Check [Client Secret Basic Auth Method](architecture#Client_Secret_Basic) on how to build it correctly.

---

When using `private_key_jwt`

Send a client assertion as JWT for us to validate the signature against the registered public key.

| Parameter             | Description                                                                                                 |
| --------------------- | ----------------------------------------------------------------------------------------------------------- |
| client_assertion      | JWT built and signed according to [Using JWTs for Client Authentication](architecture#JWT_with_Private_Key) |
| client_assertion_type | must be `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`                                            |

```BASH
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/introspect \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer \
  --data client_assertion=eyJhbGciOiJSUzI1Ni... \
  --data token=VjVxyCZmRmWYqd3_F5db9Pb9mHR5fqzhn...
```

## userinfo_endpoint

[https://api.zitadel.ch/oauth/v2/userinfo](https://api.zitadel.ch/oauth/v2/userinfo)

## end_session_endpoint

[https://accounts.zitadel.ch/oauth/v2/endsession](https://accounts.zitadel.ch/oauth/v2/endsession)

> The end_session_endpoint is located with the login page, due to the need of accessing the same cookie domain

## jwks_uri

[https://api.zitadel.ch/oauth/v2/keys](https://api.zitadel.ch/oauth/v2/keys)

> Be aware that these keys can be rotated without any prior notice. We will however make sure that a proper `kid` is set with each key!

## OAuth 2.0 Metadata

**ZITADEL** does not yet provide a OAuth 2.0 Metadata endpoint but instead provides a [OpenID Connect Discovery Endpoint](#OpenID_Connect_1_0_Discovery).

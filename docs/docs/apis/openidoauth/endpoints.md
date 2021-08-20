---
title: Endpoints
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

## OpenID Connect 1.0 Discovery

The OpenID Connect Discovery Endpoint is located within the issuer domain.
For example with [zitadel.ch](https://zitadel.ch), issuer.zitadel.ch would be the domain. This would give us [https://issuer.zitadel.ch/.well-known/openid-configuration](https://issuer.zitadel.ch/.well-known/openid-configuration).

**Link to spec.** [OpenID Connect Discovery 1.0 incorporating errata set 1](https://openid.net/specs/openid-connect-discovery-1_0.html)

## authorization_endpoint

[https://accounts.zitadel.ch/oauth/v2/authorize](https://accounts.zitadel.ch/oauth/v2/authorize)

:::note
The authorization_endpoint is located with the login page, due to the need of accessing the same cookie domain
:::

The authorization_endpoint is the starting point for all initial user authentications. The user agent (browser) will be redirected to this endpoint to
authenticate the user in exchange for an authorization_code (authorization code flow) or tokens (implicit flow). 

<details>
    <summary>Links to specs</summary>
    <ul>
        <li><a href="https://datatracker.ietf.org/doc/html/rfc6749#section-3.1">Section 3.1 of OAuth2.0 (RFC6749)</a></li>
        <li><a href="https://openid.net/specs/openid-connect-core-1_0.html#AuthorizationEndpoint">Section 3.1.2 of OpenID Connect Core 1.0 incorporating errata set 1</a></li>
    </ul>
</details>

### Required request parameters

| Parameter     | Description                                                                                                                                       |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| client_id     | The id of your client as shown in Console.                                                                                                        |
| redirect_uri  | Callback uri of the authorization request where the code or tokens will be sent to. Must match exactly one of the preregistered in Console.       |
| response_type | Determines whether a `code`, `id_token token` or just `id_token` will be returned. Most use cases will need `code`. See flow guide for more info. |
| scope         | `openid` is required, see [Scopes](scopes) for more possible values. Scopes are space delimited, e.g. `openid email profile`                      |

:::important
Following the [OIDC Core 1.0 specs](https://openid.net/specs/openid-connect-core-1_0.html#ScopeClaims) whenever an access_token is issued, 
the id_token will not contain any claims of the scopes `profile`, `email`, `phone` and `address`. 

Send the access_token to the [userinfo_endpoint](#userinfo_endpoint) or [introspection_endpoint](#introspection_endpoint) the retrieve these claims
or set the `id_token_userinfo_assertion` Option ("User Info inside ID Token" in Console) to true.
:::

Depending on your authorization method you will have to provide additional parameters or headers:

<Tabs
    groupId="token-auth-methods"
    defaultValue="client_secret_basic"
    values={[
        {label: 'client_secret_basic', value: 'client_secret_basic'},
        {label: 'client_secret_post', value: 'client_secret_post'},
        {label: 'none (PKCE)', value: 'none'},
        {label: 'private_key_jwt', value: 'private_key_jwt'},
    ]}
>
<TabItem value="client_secret_basic">
no additional parameters required
</TabItem>
<TabItem value="client_secret_post">
no additional parameters required
</TabItem>
<TabItem value="none">

| Parameter             | Description                                           |
| --------------------- | ----------------------------------------------------- |
| code_challenge        | The SHA-256 value of the generated `code_verifier`    |
| code_challenge_method | Method used to generate the challenge, must be `S256` |

see PKCE guide for more information

</TabItem>
<TabItem value="private_key_jwt">
no additional parameters required
</TabItem>
</Tabs>

### Additional parameters

| Parameter     | Description                                                                                                                                                                                                                                     |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id_token_hint | Valid `id_token` (of an existing session) used to identity the subject. **Should** be provided when using prompt `none`.                                                                                                                        |
| login_hint    | A valid logon name of a user. Will be used for username inputs or preselecting a user on `select_account`                                                                                                                                       |
| max_age       | Seconds since the last active successful authentication of the user                                                                                                                                                                             |
| nonce         | Random string value to associate the client session with the ID Token and for replay attacks mitigation. **Must** be provided when using **implicit flow**.                                                                                     |
| prompt        | If the Auth Server prompts the user for (re)authentication. <br />no prompt: the user will have to choose a session if more than one session exists<br />`none`: user must be authenticated without interaction, an error is returned otherwise <br />`login`: user must reauthenticate / provide a user name <br />`select_account`: user is prompted to select one of the existing sessions or create a new one <br />`create`: the registration form will be displayed to the user directly |
| state         | Opaque value used to maintain state between the request and the callback. Used for Cross-Site Request Forgery (CSRF) mitigation as well, therefore highly **recommended**.                                                                      |
| ui_locales    | Spaces delimited list of preferred locales for the login UI, e.g. `de-CH de en`. If none is provided or matches the possible locales provided by the login UI, the `accept-language` header of the browser will be taken into account.          |

### Successful Code Response

When your `response_type` was `code` and no error occurred, the following response will be returned: 

| Property | Description                                                                   |
| -------- | ----------------------------------------------------------------------------- |
| code     | Opaque string which will be necessary to request tokens on the token endpoint |
| state    | Unmodified `state` parameter from the request                                 |

### Successful Implicit Response

When your `response_type` was either `it_token` or `id_token token` and no error occurred, the following response will be returned:

| Property     | Description                                                                           |
| ------------ | ------------------------------------------------------------------------------------- |
| access_token | Only returned if `response_type` included `token`                                     |
| expires_in   | Number of second until the expiration of the `access_token`                           |
| id_token     | An `id_token` of the authorized user                                                  |
| token_type   | Type of the `access_token`. Value is always `Bearer`                                  |
| scope        | Scopes of the `access_token`. These might differ from the provided `scope` parameter. |
| state        | Unmodified `state` parameter from the request                                         |

### Error Response

Regardless of the authorization flow chosen, if an error occurs the following response will be returned to the redirect_uri.

:::note
If the redirect_uri is not provided, was not registered or anything other prevents the auth server form returning the response to the client,
the error will be display directly to the user on the auth server
:::

| Property          | Description                                                          |
| ----------------- | -------------------------------------------------------------------- |
| error             | An OAuth / OIDC [error_type](#authorize-errors)                      |
| error_description | Description of the error type or additional information of the error |
| state             | Unmodified `state` parameter from the request                        |

#### Possible errors {#authorize-errors}

| error_type                | Possible reason                                                                                                                                                              |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| invalid_request           | The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.                            |
| invalid_scope             | The requested scope is invalid. Typically the required `openid` value is missing.                                                                                            |
| unauthorized_client       | The client is not authorized to request an access_token using this method. Check in Console that the requested `response_type` is allowed in your application configuration. |
| unsupported_response_type | The authorization server does not support the requested response_type.                                                                                                       |
| server_error              | The authorization server encountered an unexpected condition that prevented it from fulfilling the request.                                                                  |

## token_endpoint

[https://api.zitadel.ch/oauth/v2/token](https://api.zitadel.ch/oauth/v2/token)

The token_endpoint will as the name suggests return various tokens (access, id and refresh) depending on the used `grant_type`. 
When using [`authorization_code`](#authorization-code-grant-code-exchange) flow call this endpoint after receiving the code from the authorization_endpoint.
When using [`refresh_token`](#authorization-code-grant-code-exchange) or [`urn:ietf:params:oauth:grant-type:jwt-bearer` (JWT Profile)](#jwt-profile-grant) you will call this endpoint directly.

### Authorization Code Grant (Code Exchange)

As mention above, when using `authorization_code` grant, this endpoint will be your second request for authorizing a user with its user agent (browser).

#### Required request Parameters

| Parameter    | Description                                                                                                   |
| ------------ | ------------------------------------------------------------------------------------------------------------- |
| code         | Code that was issued from the authorization request.                                                          |
| grant_type   | Must be `authorization_code`                                                                                  |
| redirect_uri | Callback uri where the code was be sent to. Must match exactly the redirect_uri of the authorization request. |

Depending on your authorization method you will have to provide additional parameters or headers:

<Tabs
    groupId="token-auth-methods"
    defaultValue="client_secret_basic"
    values={[
        {label: 'client_secret_basic', value: 'client_secret_basic'},
        {label: 'client_secret_post', value: 'client_secret_post'},
        {label: 'none (PKCE)', value: 'none'},
        {label: 'private_key_jwt', value: 'private_key_jwt'},
    ]}
>
<TabItem value="client_secret_basic">

Send your `client_id` and `client_secret` as Basic Auth Header. Check [Client Secret Basic Auth Method](authn-methods#client-secret-basic) on how to build it correctly.

</TabItem>
<TabItem value="client_secret_post">

Send your `client_id` and `client_secret` as parameters in the body:

| Parameter     | Description                      |
| ------------- | -------------------------------- |
| client_id     | client_id of the application     |
| client_secret | client_secret of the application |

</TabItem>
<TabItem value="none">

Send your `code_verifier` for us to recompute the `code_challenge` of the authorization request.

| Parameter     | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| code_verifier | code_verifier previously used to generate the code_challenge |

</TabItem>
<TabItem value="private_key_jwt">

Send a client assertion as JWT for us to validate the signature against the registered public key.

| Parameter             | Description                                                                                                     |
| --------------------- | --------------------------------------------------------------------------------------------------------------- |
| client_assertion      | JWT built and signed according to [Using JWTs for Client Authentication](authn-methods#jwt-with-private-key) |
| client_assertion_type | Must be `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`                                                |

</TabItem>
</Tabs>

#### Successful code response {#token-code-response}

| Property      | Description                                                                           |
| ------------- | ------------------------------------------------------------------------------------- |
| access_token  | An `access_token` as JWT or opaque token                                              |
| expires_in    | Number of second until the expiration of the `access_token`                           |
| id_token      | An `id_token` of the authorized user                                                  |
| scope         | Scopes of the `access_token`. These might differ from the provided `scope` parameter. |
| refresh_token | An opaque token. Only returned if `offline_access` scope was requested                |
| token_type    | Type of the `access_token`. Value is always `Bearer`                                  |

### JWT Profile Grant

#### Required request Parameters

| Parameter  | Description                                                                                                                   |
| ---------- | ----------------------------------------------------------------------------------------------------------------------------- |
| grant_type | Must be `urn:ietf:params:oauth:grant-type:jwt-bearer`                                                                         |
| assertion  | JWT built and signed according to [Using JWTs for Authorization Grants](grant-types#using-jwts-as-authorization-grants)               |
| scope      | [Scopes](Scopes) you would like to request from ZITADEL. Scopes are space delimited, e.g. `openid email profile`              |

```BASH
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=authorization_code \
  --data code=DKLvnksjndjsflkdjlkfgjslow... \
  --data client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer \
  --data client_assertion=eyJhbGciOiJSUzI1Ni...
```

#### Successful JWT Profile response {#token-jwt-response}

| Property      | Description                                                                           |
| ------------- | ------------------------------------------------------------------------------------- |
| access_token  | An `access_token` as JWT or opaque token                                              |
| expires_in    | Number of second until the expiration of the `access_token`                           |
| id_token      | An `id_token` of the authorized service user                                          |
| scope         | Scopes of the `access_token`. These might differ from the provided `scope` parameter. |
| token_type    | Type of the `access_token`. Value is always `Bearer`                                  |

### Refresh Token Grant

To request a new `access_token` without user interaction, you can use the `refresh_token` grant. 
See [offline_access Scope](Scopes#standard-scopes) for how to request a `refresh_token` in the authorization request.

#### Required request Parameters

| Parameter     | Description                                                                                  |
| ------------- | -------------------------------------------------------------------------------------------- |
| grant_type    | Must be `refresh_token`                                                                      |
| refresh_token | The refresh_token previously issued in the last authorization_code or refresh_token request. |
| scope         | [Scopes](Scopes) you would like to request from ZITADEL for the new access_token. Must be a subset of the scope originally requested by the corresponding auth request. When omitted, the scopes requested by the original auth request will be reused. Scopes are space delimited, e.g. `openid email profile` |

Depending on your authorization method you will have to provide additional parameters or headers:

<Tabs
    groupId="token-auth-methods"
    defaultValue="client_secret_basic"
    values={[
        {label: 'client_secret_basic', value: 'client_secret_basic'},
        {label: 'client_secret_post', value: 'client_secret_post'},
        {label: 'none (PKCE)', value: 'none'},
        {label: 'private_key_jwt', value: 'private_key_jwt'},
    ]}
>
<TabItem value="client_secret_basic">

Send your `client_id` and `client_secret` as Basic Auth Header. Check [Client Secret Basic Auth Method](authn-methods#client-secret-basic) on how to build it correctly.

</TabItem>
<TabItem value="client_secret_post">

Send your `client_id` and `client_secret` as parameters in the body:

| Parameter     | Description                      |
| ------------- | -------------------------------- |
| client_id     | client_id of the application     |
| client_secret | client_secret of the application |

</TabItem>
<TabItem value="none">

Send your `client_id` as parameter in the body. No authentication is required.

</TabItem>
<TabItem value="private_key_jwt">

Send a `client_assertion` as JWT for us to validate the signature against the registered public key.

| Parameter             | Description                                                                                                     |
| --------------------- | --------------------------------------------------------------------------------------------------------------- |
| client_assertion      | JWT built and signed according to [Using JWTs for Client Authentication](authn-methods#jwt-with-private-key) |
| client_assertion_type | Must be `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`                                                |

</TabItem>
</Tabs>

#### Successful refresh token response {#token-refresh-response}

| Property      | Description                                                                           |
| ------------- | ------------------------------------------------------------------------------------- |
| access_token  | An `access_token` as JWT or opaque token                                              |
| expires_in    | Number of second until the expiration of the `access_token`                           |
| id_token      | An `id_token` of the authorized user                                                  |
| scope         | Scopes of the `access_token`. These might differ from the provided `scope` parameter. |
| refresh_token | An new opaque refresh_token.                                                          |
| token_type    | Type of the `access_token`. Value is always `Bearer`                                  |

### Error response

> //TODO: errors

## introspection_endpoint

[https://api.zitadel.ch/oauth/v2/introspect](https://api.zitadel.ch/oauth/v2/introspect)

This endpoint enables client to validate an `acccess_token`, either opaque or JWT. Unlike client side JWT validation,
this endpoint will check if the token is not revoked (by client or logout).

| Parameter | Description     |
| --------- | --------------- |
| token     | An access token |

Depending on your authorization method you will have to provide additional parameters or headers:

<Tabs
    groupId="introspect-auth-methods"
    defaultValue="client_secret_basic"
    values={[
        {label: 'client_secret_basic', value: 'client_secret_basic'},
        {label: 'private_key_jwt', value: 'private_key_jwt'},
    ]}
>
<TabItem value="client_secret_basic">

Send your `client_id` and `client_secret` as Basic Auth Header. Check [Client Secret Basic Auth Method](authn-methods#client-secret-basic) on how to build it correctly.

</TabItem>

<TabItem value="private_key_jwt">

Send a `client_assertion` as JWT for us to validate the signature against the registered public key.

| Parameter             | Description                                                                                                 |
| --------------------- | ----------------------------------------------------------------------------------------------------------- |
| client_assertion      | JWT built and signed according to [Using JWTs for Client Authentication](authn-methods#client-secret-basic) |
| client_assertion_type | must be `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`                                            |

```BASH
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/introspect \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer \
  --data client_assertion=eyJhbGciOiJSUzI1Ni... \
  --data token=VjVxyCZmRmWYqd3_F5db9Pb9mHR5fqzhn...
```

</TabItem>
</Tabs>

### Successful introspection response {#introspect-response}

Upon successful authorization of the client a response with the boolean `active` is returned, indicating if the provided token 
is active and the requesting client is part of the token audience.

If `active` is **true**, further information will be provided:

| Property  | Description                                          |
| --------- | ---------------------------------------------------- |
| scope     | Space delimited list of scopes granted to the token. |

Additionally and depending on the granted scopes, information about the authorized user is provided. 
Check the [Claims](claims) page if a specific claims might be returned and for detailed description.

### Error response {#introspect-error-response}

If the authorization fails, an HTTP 401 with `invalid_client` will be returned.

## userinfo_endpoint

[https://api.zitadel.ch/oauth/v2/userinfo](https://api.zitadel.ch/oauth/v2/userinfo)

This endpoint will return information about the authorized user.

Send the `access_token` of the **user** (not the client) as Bearer Token in the `authorization` header:
```BASH
curl --request GET \
  --url https://api.zitadel.ch/oauth/v2/userinfo
  --header 'Authorization: Bearer dsfdsjk29fm2as...'
```

### Successful userinfo response {#userinfo-response}

If the `access_token` is valid, the information about the user depending on the granted scopes is returned.
Check the [Claims](claims) page if a specific claims might be returned and for detailed description.

### Error response {#userinfo-error-response}

If the token is invalid or expired, an HTTP 401 will be returned.

## end_session_endpoint

[https://accounts.zitadel.ch/oauth/v2/endsession](https://accounts.zitadel.ch/oauth/v2/endsession)

> The end_session_endpoint is located with the login page, due to the need of accessing the same cookie domain

## jwks_uri

[https://api.zitadel.ch/oauth/v2/keys](https://api.zitadel.ch/oauth/v2/keys)

> Be aware that these keys can be rotated without any prior notice. We will however make sure that a proper `kid` is set with each key!

## OAuth 2.0 Metadata

**ZITADEL** does not yet provide a OAuth 2.0 Metadata endpoint but instead provides a [OpenID Connect Discovery Endpoint](#OpenID_Connect_1_0_Discovery).

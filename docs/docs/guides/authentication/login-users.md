---
title: Login Users into your Application
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import CodeFlowChart from '../../imports/_code-flow-chart.md';

## Overview

This guide will show you how to use ZITADEL to login users into your application (authentication).
It will guide you step-by-step through the basics and point out on how to customize process.

## OIDC / OAuth Flow

OAuth and therefore OIDC know three different application types:
- Web (Server-side web applications such as java, .net, ...)
- Native (native, mobile or desktop applications)
- User Agent (single page applications / SPA, generally JavaScript executed in the browser)

Depending on the app type you're trying to register, there are small differences.
But regardless of the app type we recommend using Proof Key for Code Exchange (PKCE).

Please read the following guide about the [different-client-profiles](../authorization/oauth-recommended-flows#different-client-profiles) and why to use PKCE.

### Code Flow

For a basic understanding of OAuth and its flows, we'll briefly describe the most important flow: **Authorization Code**.

The following diagram demonstrates a basic authorization_code flow: 

<CodeFlowChart />

1. When an unauthenticated user visits your application,
2. you will create an authorization request to the authorization endpoint.
3. The Authorization Server (ZITADEL) will send an HTTP 302 to the user's browser, which will redirect him to the login UI.
4. The user will have to authenticate using the demanded auth mechanics.
5. Your application will be called on the registered callback path (redirect_uri) and be provided an authorization_code.
6. This authorization_code must then be sent together with you applications authentication (client_id + client_secret) to the token_endpoint
7. In exchange the Authorization Server (ZITADEL) will return an access_token and if requested a refresh_token and in the case of OIDC an id_token as well 

This flow is the same when using PKCE or JWT with Private Key for authentication.

## Create Application

- choose app type
- go into console
- follow wizard

<Tabs
    groupId="app-type"
    default="web"
    values={[
        {'label': 'Web', 'value': 'web'},
        {'label': 'Native', 'value': 'native'},
        {'label': 'SPA', 'value': 'spa'},
    ]}
>
<TabItem value="web">

<Tabs
    groupId="auth-type"
    default="pkce"
    values={[
        {'label': 'PKCE', 'value': 'pkce'},
        {'label': 'Basic Auth', 'value': 'basic'},
        {'label': 'JWT with Private Key', 'value': 'jwt'},
    ]}
>
<TabItem value="web">
</TabItem>
</Tabs>
1. pkce
2. basic
3. jwt
4. post

redirects...

</TabItem>
<TabItem value="native">

<Tabs
    groupId="auth-type"
    default="pkce"
    values={[
        {'label': 'PKCE', 'value': 'pkce'},
    ]}
>
<TabItem value="pkce">
</TabItem>
</Tabs>

redirects...

enable refresh_token

additional origins

</TabItem>
<TabItem value="spa">

<Tabs
    groupId="auth-type"
    default="pkce"
    values={[
        {'label': 'PKCE', 'value': 'pkce'},
        {'label': 'Implicit', 'value': 'implicit'},
    ]}
>
<TabItem value="pkce">
</TabItem>
</Tabs>

redirects...

</TabItem>
</Tabs>

## Auth Request

To initialize the user authentication, you will have to create an authorization request using HTTP GET in the user agent (browser) 
on /authorize with at least the following parameters:
- `client_id`: this tells the authorization server which application it is, copy from Console
- `redirect_uri`: where the authorization code is sent to after the user authentication, must be one of the registered in the previous step
- `response_type`: if you want to have a code (authorization code flow) or directly a token (implicit flow), so when ever possible use `code`
- `scope`: what scope you want to grant to the access_token / id_token, minimum is `openid`, typically you will have `openid profile email`

We recommend always using two additional parameters `state` and `nonce`. The former enables you to transfer a state through
the authentication process. The latter is used to bind the client session with the id_token and to mitigate replay attacks.

Depending on your authentication method you might need to provide additional parameters:

<Tabs
    groupId="auth-type"
    default="pkce"
    values={[
        {'label': 'PKCE', 'value': 'pkce'},
        {'label': 'Basic Auth', 'value': 'basic'},
        {'label': 'JWT with Private Key', 'value': 'jwt'},
    ]}
>
<TabItem value="pkce">

PKCE stands for Proof Key for Code Exchange. So other than "normal" code exchange, the does not authenticate using
client_id and client_secret but an additional code. You will have to generate a random string, hash it and send this hash
on the authorization_endpoint. On the token_endpoint you will then send the plain string for the authorization to compute
the hash as well and to verify it's correct. In order to do so you're required to send the following two parameters as well:
- `code_challenge`: the base64url representation of the (sha256) hash of your random string
- `code_challenge_method`: must always be `S256` standing for sha256, this is the only algorithm we support

For example for `random-string` the code_challenge would be `9az09PjcfuENS7oDK7jUd2xAWRb-B3N7Sr3kDoWECOY`

The request would finally look like (linebreaks and whitespace for display reasons):

```curl
curl --request GET \
  --url 'https://accounts.zitadel.ch/oauth/v2/authorize
    ?client_id=${client_id}
    &redirect_uri=${redirect_uri}
    &response_type=code
    &scope=openid%20email%20profile
    &code_challenge=${code_challenge}
    &code_challenge_method=S256'
```

</TabItem>
<TabItem value="basic">

You don't need any additional parameter for this request. We're identifying the app by the `client_id` parameter.

So your request might look like this (linebreaks and whitespace for display reasons):

```curl
curl --request GET \
  --url 'https://accounts.zitadel.ch/oauth/v2/authorize
    ?client_id=${client_id}
    &redirect_uri=${redirect_uri}
    &response_type=code
    &scope=openid%20email%20profile'
```

</TabItem>
<TabItem value="jwt">

You don't need any additional parameter for this request. We're identifying the app by the `client_id` parameter.

So your request might look like this (linebreaks and whitespace for display reasons):

```curl
curl --request GET \
  --url 'https://accounts.zitadel.ch/oauth/v2/authorize
    ?client_id=${client_id}
    &redirect_uri=${redirect_uri}
    &response_type=code
    &scope=openid%20email%20profile'
```

</TabItem>
</Tabs>

### Additional parameters and customization

There are additional parameters and values you can provide to satisfy your use case and to customize the user's authentication flow.
Please check the [authorization_endpoint reference](../../apis/openidoauth/endpoints.md#authorization_endpoint) in the OAuth / OIDC documentation.

## Callback

Regardless of a successful or error response from the authorization_endpoint, the authorization server will call your
callback endpoint you provided by the `redirect_uri`.

:::note
If the redirect_uri is not provided, was not registered or anything other prevents the auth server form returning the response to the client,
the error will be display directly to the user on the auth server.
:::

Upon successful authentication you'll be given a `code` and if provided the unmodified `state` parameter.
You will need this `code` in the token request.

If a parameter was missing, malformed or any other error occurred, your answer will contain an `error` stating the error type,
possibly an `error_description` providing some information about the error and its reason and the `state` parameter.
Check the [error response section](../../apis/openidoauth/endpoints#error-response) in the authorization_endpoint reference.

## Token request

Next you will have to exchange the given `code` for the tokens. For this HTTP POST request you will need to provide the following:
- code: the code that was issued from the authorization request
- grant_type: must be `authorization_code`
- redirect_uri: callback uri where the code was sent to. Must match exactly the redirect_uri of the authorization request

Depending on your authentication method you'll need additional headers and parameters:

<Tabs
    groupId="auth-type"
    defaultValue="pkce"
    values={[
        {label: 'PKCE', value: 'pkce'},
        {label: 'Basic Auth', value: 'basic'},
        {label: 'JWT with Private Key', value: 'jwt'},
    ]}
>
<TabItem value="pkce">

Send your `client_id` and the previously generated string as `code_verifier` for us to recompute the `code_challenge` of the authorization request:

```curl
curl --request POST \
--url https://api.zitadel.ch/oauth/v2/token \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data grant_type=authorization_code \
--data code=${code} \
--data redirect_uri=${redirect_uri} \
--data client_id=${client_id} \
--data code_verifier=${code_verifier}
```

</TabItem>
<TabItem value="basic">

Send your `client_id` and `client_secret` as Basic Auth Header. Note that OAuth2 requires client_id and client_secret to be form url encoded. 
So check [Client Secret Basic Auth Method](../../apis/openidoauth/authn-methods#client-secret-basic) on how to build it correctly.

```curl
curl --request POST \
--url https://api.zitadel.ch/oauth/v2/token \
--header 'Authorization: Basic ${basic}' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data grant_type=authorization_code \
--data code=${code} \
--data redirect_uri=${redirect_uri}
```

</TabItem>
<TabItem value="jwt">

Send a `client_assertion` as JWT and the `client_assertion_type` `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`
for us to validate the signature against the registered public key:

```curl
curl --request POST \
--url https://api.zitadel.ch/oauth/v2/token \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data grant_type=authorization_code \
--data code=${code} \
--data redirect_uri=${redirect_uri} \
--data client_assertion=${client_assertion} \
--data client_assertion_type=urn%3Aietf%3Aparams%3Aoauth%3Aclient-assertion-type%3Ajwt-bearer
```

</TabItem>
</Tabs>

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

scopes

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
custom scopes

## Token Request


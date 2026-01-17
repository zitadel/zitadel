---
title: Log Out Users from an Application with ZITADEL
sidebar_label: Logout
---

This guide shows you the different concepts and use cases of the logout process and how to use it in ZITADEL.

## OpenID Connect Single Logout 

### Single Sign On (SSO) vs Single Logout (SLO)

Single Sign On (SSO) allows a user to login once without the need for authentication across multiple applications.
Single Logout (SLO) is the counterpart to SSO. With SLO a user can logout and terminate sessions across many applications, without actively logging out from them.

The purpose of a logout is to terminate a user session.
Depending on how the session handling is implemented, there are different mechanisms that can be used.
There are two possibilities where sessions are stored:
- The User Agent (e.g the Browser or Mobile App) stores the session information (e.g. in a cookie)
- A Server stores the session information (e.g. in a database or api)

### OpenID Connect Logout

OpenID Connect defines three logout mechanisms to address the different architectures:
- [OpenID Connect Session Management 1.0](https://openid.net/specs/openid-connect-session-1_0.html)
- [OpenID Connect Front-Channel Logout 1.0](https://openid.net/specs/openid-connect-frontchannel-1_0.html)
- [OpenID Connect Back-Channel Logout 1.0](https://openid.net/specs/openid-connect-backchannel-1_0.html)

#### Session Management

Session Management in OpenID Connect defines a mechanism for a client (Relying Party, RP) to monitor the state of the user session from a identity provider (OP, e.g ZITADEL).
When a user logs out of the provider, the user's session is terminated and the client can in turn reflect that in its behavior.

#### RP initiated Logout

With the RP initiated flow all logout processes are triggered by a request from the client (e.g your application) through a well defined standard API by redirecting the user-agent to the [end_session_endpoint](/docs/apis/openidoauth/endpoints#end_session_endpoint).
If you have specified some post_logout_redirect_uris on your client you have to send either the id_token_hint or the client_id as param in your request.
So ZITADEL is able to read the configured redirect uris.

```
GET $CUSTOM-DOMAIN/oidc/v1/end_session
    ?id_token_hint={id_token}
    &post_logout_redirect_uri=https://rp.example.com/logged_out
    &state=random_string
```

#### Front-Channel Logout

The user agent handles the front-channel logout. 
Each client with an OpenID Session of the user that supports front-channel renders an iframe so the logout request is performed on all clients parallel.

:::note
This is not yet implemented in ZITADEL
:::

#### Back-Channel Logout

The back-channel logout is a mechanism on the server-side and the user agent does not have to do anything.
The user will logout from all clients even in the case the user agent was closed.

:::note
This is not yet implemented in ZITADEL
:::

## Scenarios

1. Logout all users from the current user-agent/browser (current implementation of ZITADEL end_session_endpoint)
2. Logout my user from the current user-agent/browser
3. Logout my user from the all devices

## Session Handling in ZITADEL

The session management in ZITADEL is done on the server side. 
As soon as a user authenticates the first time, ZITADEL generates a user-agent cookie.
All open sessions on that user-agent (browser) will be stored to the same cookie.
If you delete the cookie in your browser, we will not be able to find out which sessions belong to your user-agent.

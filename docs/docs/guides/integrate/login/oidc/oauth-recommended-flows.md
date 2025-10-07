---
title: Recommended authorization flows with OpenID Connect (OIDC) and OAuth 2.x
sidebar_label: Recommended authorization flows 
---

## Introduction

In this guide we will go over some basics on how to obtain an authorization with OpenID Connect 1.x and OAuth 2.x.

ZITADEL does not make assumptions about the application type you are about to integrate. Therefore you must qualify and define an appropriate method for your users to gain authorization to access your application (“authentication flow”). Your choice depends on the application’s ability to maintain the confidentiality of client credentials and the technological capabilities of your application. If you choose a deprecated or unfeasible flow to obtain authorization for your application, your users’ credentials may be compromised.

We invite you to further explore the different authorization flows in the OAuth 2.x standard. We assume that you have a brand-new application (ie. without legacy requirements), and you've found a reliable SDK/Library that does the heavy lifting for you.

So this module will only go over the basics and explain why we recommend the flow “Authorization Flow with PKCE” as default for most applications. We will also cover the case of machine-to-machine communication, ie. where there is no interactive login. Further we will guide you to further reading viable alternatives, if the default flow is not feasible.

## Basics of Federated Identity

Although Federated Identities are not a new concept ([RFC 6749](https://tools.ietf.org/html/rfc6749), “The OAuth 2.0 Authorization Framework” was released in 2012) it is important to highlight the difference between the traditional client-server authentication model and the concept of delegated authorization and authentication.

The aforementioned RFC provides us with some [problems and limitations](https://tools.ietf.org/html/rfc6749#section-1) of the client-server authentication, where a client requests a protected resource on the server by authenticating with the user’s credentials:

- Applications need to store users credentials (eg, password) for future use; compromise of any application results in compromise of the end-users credentials
- Servers are required to support password authentication
- Without means of limiting scope when providing the user’s credentials, the application gains overly broad access to protected resources
- Users cannot revoke access for a single application, but only for all by changing credentials

So what do we want to achieve with delegated authentication?

- Instead of implementing authentication on each server and trusting each server

  - Users only **authenticate** with a trusted server (ie. ZITADEL), that validates the user’s identity through a challenge (eg, multiple authentication factors) and issues an **ID token** (OpenID Connect 1.x)
  - Applications have means of **validating the integrity** of presented access and ID tokens

- Instead of sending around the user’s credentials
  - Clients may access protected resources with an **access token** that is only valid for specific scope and limited lifetime (OAuth 2.x)
  - Users have to **authorize** applications to access certain [**scopes**](/apis/openidoauth/scopes) (eg, email address or custom roles). Applications can request [**claims**](/apis/openidoauth/claims) (key:value pairs, eg email address) for the authorized scopes with the access token or ID token from ZITADEL
  - Access tokens are bearer tokens, meaning that possession of the token provides access to a resource. But the tokens expire frequently and the application must request a new access token via **refresh token** or the user must reauthenticate

![Overview federated identities](/img/guides/consulting_federated_identities_basics.png)

This is where the so-called “flows” come into play: There are a number of different flows on how to handle the process from authentication, over authorization, getting tokens and requesting additional information about the user.

Maybe interesting to mention is that we are mostly concerned with choosing the right OAuth 2.x flows (alas “authorization”). OpenID Connect extends the OAuth 2.x flow with useful features like endpoint discovery (where to ask), ID Token (who is the user, when and how did she authenticate), and UserInfo Endpoint (getting additional information about the user).

## Different client profiles

As mentioned in the beginning of this module, there are two main determinants for choosing the optimal authorization flow:

1. Client’s ability to maintain the confidentiality of client credentials
2. Technological limitations

OAuth 2.x defines two [client types](https://tools.ietf.org/html/rfc6749#section-2.1) based on their ability to maintain the confidentiality of their client credentials:

- Confidential: Clients capable of maintaining the confidentiality of their credentials (e.g., client implemented on a secure server with restricted access to the client credentials), or capable of secure client authentication using other means.
- Public: Clients incapable of maintaining the confidentiality of their credentials (e.g., clients executing on the device used by the resource owner, such as an installed native application or a web browser-based application), and incapable of secure client authentication via any other means.

The following table gives you a brief overview of different client profiles.

<table className="table-wrapper">
    <tbody>
	<tr>
		<th>Confidentiality of client credentials</th>
		<th>Client profile</th>
		<th>Examples of clients</th>
    </tr>
	<tr>
		<td rowSpan="2">Public</td>
		<td>User-Agent</td>
		<td>Single-page web applications (SPA), generally JavaScript executed in Browser</td>
	</tr>
	<tr>
		<td>Native</td>
		<td>Native, Mobile, or Desktop applications</td>
	</tr>
	<tr>
		<td rowSpan="2">Confidential</td>
		<td>Web</td>
		<td>Server-side web applications such as java, .net, …</td>
	</tr>
	<tr>
		<td>Machine-to-Machine</td>
		<td>APIs, generally services without direct user-interaction at the authorization server</td>
	</tr>
    </tbody>
</table>

## Our recommended authorization flows

We recommend using the flow **“Authorization Code with Proof Key of Code Exchange (PKCE)”** ([RFC7636](https://tools.ietf.org/html/rfc7636)) for **User-Agent**, **Native**, and **Web** clients.

If you don’t have any technical limitations, you should favor the flow Authorization Code with PKCE over other methods. The PKCE part makes the flow resistant against authorization code interception attack as described well in RFC7636.

_So what about APIs?_

We recommend using **“JWT bearer token with private key”** ([RFC7523](https://tools.ietf.org/html/rfc7523)) for Machine-to-Machine clients.

What this means is that you have to send an JWT token, containing the [standard claims for access tokens](/apis/openidoauth/claims) and that is signed with your private key, to the token endpoint to request the access token. We will see how this works in another module about Service Accounts.

If you don’t have any technical limitations, you should prefer this method over other methods.

A JWT with a private key can also be used with client profile web to further enhance security.

In case you need alternative flows and their advantages and drawbacks, there will be a module to outline more methods and our recommended fallback strategy per client profile that are available in ZITADEL.

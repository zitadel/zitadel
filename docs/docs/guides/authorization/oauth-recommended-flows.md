---
title: Recommended authorization flows
---

<table class="table-wrapper">
    <tr>
        <td>Description</td>
        <td>Learn about the different authentication flows and about the flow that we recommend for your application.</td>
    </tr>
    <tr>
        <td>Learning Outcomes</td>
        <td>
            In this module you will:
            <ul>
              <li>Learn the basics of federated identities</li>
              <li>Understand the basics of OAuth 2.x client profiles and their importance for authorization flows</li>
              <li>Get a recommended flow for Web, Native, User-Agent, and API</li>
            </ul>
        </td>
    </tr>
     <tr>
        <td>Prerequisites</td>
        <td>
            Basic knowledge about federated identities
        </td>
    </tr>
</table>


## Introduction

Before we set up our first application within ZITADEL, we need outline the basics of how to authorize with OpenID Connect 1.x and OAuth 2.x.

ZITADEL makes no assumptions about the application type that you are going to integrate.
Therefore you must choose an appropriate method to authorize your users (i.e. define an *authentication flow*).

What you choose depends on your application’s ability to keep client credentials confidential, and on the app's technological capabilities.
If you chose a deprecated or unfeasible flow to obtain authorization for your application, you might compromise your users' credentials.

We invite you to explore the different authorization flows in the OAuth 2.x standard.
At the start, we assume that you have a brand-new application (ie. without legacy requirements),
and that you found a reliable SDK/Library that does the heavy lifting for you.

This module covers only the basics.
It explains why we recommend the flow “Authorization Flow with PKCE” as default for most applications.
We also cover the case of machine-to-machine communication, ie. where there is no interactive login.
For cases where the default flow is not feasible.

## Basics of Federated Identity

Federated Identities are not a new concept.
[RFC 6749](https://tools.ietf.org/html/rfc6749),
“The OAuth 2.0 Authorization Framework,” was published in 2012.
Nevertheless, it is important to highlight the difference between the traditional client-server authentication model and the concept of delegated authorization and authentication.

In client-server authentication, a client requests a protected resource on the server by authenticating with the user’s credentials.
The aforementioned RFC describes some [problems and limitations](https://tools.ietf.org/html/rfc6749#section-1) of this approach: 

* Applications need to store users credentials (eg, password) for future use; compromising any application also compromises the end-user credentials.
* Servers must support password authentication.
* If an application cannot limit scope when providing the user’s credentials, the application gains overly broad access to protected resources.
* Users cannot revoke access for a single application, but only for all (by changing credentials).

So what do we want to achieve with delegated authentication?

* Instead of implementing authentication on each server and trusting each server:
  * Users *authenticate* only with a trusted server (ie. ZITADEL), which validates the user’s identity through a challenge (eg, multiple authentication factors) and issues an *ID token* (OpenID Connect 1.x).
  * Applications have a way to *validate the integrity* of the presented access and ID tokens.

* Instead of sending around the user’s credentials:
  * Clients can use an *access token* to access protected resources. This token is valid only for a specific scope and lifetime (OAuth 2.x).
  * Users have to *authorize* applications to access certain [**scopes**](https://docs.zitadel.ch/architecture#Scopes) (e.g. email address or custom roles). Applications can request [**claims**](https://docs.zitadel.ch/architecture#Claims) (key:value pairs, e.g. email address) for the authorized scopes with the access token or ID token from ZITADEL
  * Access tokens are bearer tokens.
   They provide access to whoever possesses the token. But the tokens expire frequently and the application must either request a new access token via *refresh token*, or have the user reauthenticate.

![Overview federated identities](/img/guides/consulting_federated_identities_basics.png)

This is where the so-called “flows” come into play: There are a number of different flows to handle the process of authenticating, authorizating, getting tokens, and requesting additional user information.

It is worth mentioning that we are mostly concerned with choosing the right OAuth 2.x flows (“authorization”).
OpenID Connect extends the OAuth 2.x flow with useful features like:
* Endpoint discovery&mdash;where to ask?
* ID Tokens&mdash;who are the users, when and how did they authenticate?
* UserInfo Endpoint&mdash;what additional information is there about the user?

## Different client profiles

As mentioned earlier in this module, two main considerations determine the optimal authorization flow:

1. The client’s ability to maintain client-credential confidentiality. 
2. Technological limitations

OAuth 2.x defines two [client types](https://tools.ietf.org/html/rfc6749#section-2.1) based on their ability to maintain the confidentiality of their client credentials:

* Confidential: Clients capable of maintaining the confidentiality of their credentials (e.g., the client is implemented on a secure server, with restricted access to client credentials), or capable of securing client authentication using other means.
* Public: Clients incapable of maintaining the confidentiality of their credentials (e.g., clients executing on the device used by the resource owner, such as an installed native application or a web browser-based application), and incapable of secure client authentication via any other means.

The following table gives you a brief overview of different client profiles.

<table class="table-wrapper">
	<tr>
		<th>Confidentiality of client credentials</th>
		<th>Client profile</th>
		<th>Examples of clients</th>
    </tr>
	<tr>
		<td rowspan="2">Public</td>
		<td>User-Agent</td>
		<td>Single-page web applications (SPA), generally JavaScript executed in Browser</td>
	</tr>
	<tr>
		<td>Native</td>
		<td>Native, Mobile, or Desktop applications</td>
	</tr>
	<tr>
		<td rowspan="2">Confidential</td>
		<td>Web</td>
		<td>Server-side web applications such as java, .net, etc.</td>
	</tr>
	<tr>
		<td>Machine-to-Machine</td>
		<td>APIs, generally services without direct user-interaction, on the authorization server</td>
	</tr>
</table>

## Our recommended authorization flows

For User-Agent, Native, and Web clients,
**we recommend using the flow “Authorization Code with Proof Key of Code Exchange (PKCE)”** ([RFC7636](https://tools.ietf.org/html/rfc7636)).

If you don’t have any technical limitations, favor the flow *Authorization Code with PKCE* over other methods.
The PKCE part makes the flow resistant to authorization code interception attacks, as described well in RFC7636.

*So what about APIs?*

**For APIs, we recommend using _JWT bearer token with private key_ ([RFC7523](https://tools.ietf.org/html/rfc7523)) for Machine-to-Machine clients.

What this means is that, to request the access token, you must send a JWT token to the token endpoint.
This JWT must contain the [standard claims for access tokens](https://docs.zitadel.ch/architecture#Claims), and be signed with your private key.
We will see how this works in another module about Service Accounts.

If you don’t have any technical limitations, you should prefer this method over other methods.

To further enhance security, you can also use a JWT with a private key with a client profile web.

In some cases, you might need alternative flows, with their advantages and drawbacks.
There will be a module that outlines more methods and our recommended fallback strategies per client profile available in ZITADEL.

## Knowledge Check (3)

* With federated identities the user sends credentials to the server holding the protected resource.
    - [ ] yes
    - [ ] no
* ZITADEL discovers your client profile automatically and sets the correct flow.
    - [ ] yes
    - [ ] no
* When working with APIs / machine-to-machine communication, it is recommended to exchange a JWT, singed with your private key, for an access token
    - [ ] yes
    - [ ] no

<details>
    <summary>
        Solutions
    </summary>

* With federated identities the user sends credentials to the server holding the protected resource
    - [ ] yes
    - [x] no (Users are authenticated against a centralized IDP, only access tokens are sent to the requested resources).
* ZITADEL will discover your client profile automatically and set the correct flow
    - [ ] yes
    - [x] no (ZITADEL does not make any assumptions about your application’s requirements).
* When working with APIs / machine-to-machine communication its recommended to exchange a JWT that is singed with your private key for an access token
    - [x] yes
    - [ ] no

</details>

## Summary (3)

* Federated Identities solve key problems and challenges with traditional server-client architecture.
* Use _Authorization Code with Proof Key of Code Exchange (PKCE)_ for User-Agent, Native, and Web clients.
* Use a *JWT bearer token with private key* for Machine-to-Machine clients.
* If these flows are technically not possible, ZITADEL supports fallback flows and strategies.

Where to go from here

* Applications
* Service Accounts
* Alternative authentication flows (aka. "The Zoo")

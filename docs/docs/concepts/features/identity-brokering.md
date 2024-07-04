---
title: Identity Brokering
sidebar_label: Identity Brokering
---

Link social logins and external identity providers with your identity management platform allowing users to login with their preferred identity provider.

Establish a trusted connection between your central identity provider (IdP) and third party identity providers.

By using a central identity brokering service you don't need to develop and establish a trust relationship between each application and each identity provider individually.

## What are federated identities?

Federated identity management is an arrangement built upon the trust between two or more domains. Users of these domains are allowed to access applications and services using the same identity.
This identity is known as federated identity and the pattern behind this is identity federation.

Compatibility across various IdPs is ensured by using industry standard protocols, such as: 

* OpenID Connect (OIDC): A modern and versatile protocol for secure authentication.
* SAML2: A widely adopted protocol for secure single sign-on (SSO) in enterprise environments.
* LDAP: A lightweight protocol for accessing user data directories commonly used in corporate networks.

## What is identity brokering?

A service provider that specializes in brokering access control between multiple service providers (also referred to as relying parties) is called an identity broker.
Federated identity management is an arrangement that is made between two or more such identity brokers across organizations.

For example, if Google is configured as an identity provider in your organization, the user will get the option to use his Google Account on the Login Screen of ZITADEL.
Because Google is registered as a trusted identity provider, the user will be able to login in with the Google account after the user is linked with an existing ZITADEL account (if the user is already registered) or a new one with the claims provided by Google.

![Diagram of an identity brokering scheme using a central identity provider that has a trust link to the Google IdP and Entra ID](/img/concepts/features/identity-brokering.png)

The schema is a very simplified version, but shows the essential steps for identity brokering

1. An unauthenticated user wants to use the alpha.com's application.
2. The application redirects the user to alpha.com's identity provider (IdP).
3. Based on the user's tenants configuration the IdP presents the configured identity providers, or redirects the user directly to the primary external IdP. The user authenticates with their external identity provider (eg, Entra ID).
4. After the authentication, the user is redirected back to alpha.com's identity provider. If the user doesn't exist in the IdP the user will be created just-in-time and linked to the external identity provider for future reference.
5. As with a local authentication, the IdP issues a token to the user that can be used to access the application. The IdP redirects the user, which is now authenticated, eventually to the application.

## Is single-sign-on (SSO) the same as identity brokering?

Sometimes single-sign-on (SSO) and login with third party identity providers is used interchangeably.
Typically SSO describes an authentication scheme that allows users to log in once at a central identity provider and access service providers (client applications) without to login again.

Identity brokering describes an authentication scheme where users can login with external identity providers that have a established trust with an identity provider which facilitates the authentication for the requested applications.

The connection between the two lies in how SSO can be implemented as part of an identity brokering solution.
In such cases, the identity broker uses SSO to enable seamless access across multiple systems, handling the complexities of different authentication protocols and standards behind the scenes.
This allows users to log in once and gain access to multiple systems that the broker facilitates.

## Multitenancy and identity brokering

In a multi-tenancy application, you want to be able to configure an external identity provider per tenant.
For example some organizations might use their EntraID, some other want to login with their OKTA, or Google Workspace.

Using an identity provider with strong multitenancy capabilities such as ZITADEL, you can configure a different set of external identity providers per organization.

[Domain discovery](/docs/guides/solution-scenarios/domain-discovery) ensures that users are redirected to their external identity provider based on their email-address or username.
[Managers](../structure/managers) can configure organization domains that are used for domain-based redirection to an external IdP.

![Diagram explaining domain discovery](/img/concepts/features/domain-discovery.png)

## Simplify identity brokering with ZITADEL templates

ZITADEL works with SAML, OpenID Connect, and LDAP external identity providers.

For popular IdPs such as EntraID, Okta, Google, Facebook, and GitHub, ZITADEL [offers pre-configured templates](/docs/guides/integrate/identity-providers/introduction).
These templates expedite the configuration process, allowing organizations to quickly integrate these providers with minimal effort.

ZITADEL recognizes that specific needs may extend beyond pre-built templates.
To address this, ZITADEL provides generic templates that enable connection to virtually any IdP.  This ensures maximum flexibility and future-proofs login infrastructure, accommodating future integrations with ease.

### References

* [Detailed integration guide for many identity providers](/guides/integrate/identity-providers/introduction)
* [Setup identity providers with Console](/guides/manage/console/default-settings#identity-providers)
* [Configure identity providers with the ZITADEL API](/docs/apis/resources/mgmt/identity-providers)

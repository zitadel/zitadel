---
title: Identity Brokering
---

## What are Identity Brokering and Federated Identities?

Federated identity management is an arrangement built upon the trust between two or more domains. Users of these domains are allowed to access applications and services using the same identity.
This identity is known as federated identity and the pattern behind this is identity federation.

A service provider that specializes in brokering access control between multiple service providers (also referred to as relying parties) is called an identity broker.
Federated identity management is an arrangement that is made between two or more such identity brokers across organizations.

For example, if Google is configured as an identity provider in your organization, the user will get the option to use his Google Account on the Login Screen of ZITADEL. Because Google is registered as a trusted identity provider, the user will be able to login in with the Google account after the user is linked with an existing ZITADEL account (if he is already registered) or a new one with the claims provided by Google.

![Identity Brokering](/img/guides/identity_brokering.png)

## How to use external identity providers in ZITADEL

Configure external identity providers on the instance level or just for one organization via the [Console](/guides/manage/console/instance-settings#identity-providers) or ZITADEL APIs.

The guides listed in this section will help you set up specific identity providers.
ZITADEL also provides templates to configure generic identity providers, which don't have templates.

---
title: Identity Brokering
---

## What is Identity Brokering and Federated Identities?

Federated identity management is an arrangement built upon the trust between two or more domains. Users of these domains are allowed to access applications and services using the same identity.
This identity is known as federated identity and the pattern behind this as identity federation.

A service provider that specializes in brokering access control between multiple service providers (also referred to as relying parties) is called identity broker.
Federated identity management is an arrangement that is made between two or more such identity brokers across organizations.

Example:
If Google is configured as identity provider on your organization, the user will get the option to use his Google Account on the Login Screen of ZITADEL (1).
ZITADEL will redirect the user to the login screen of Google where he as to authenticated himself (2) and is sent back after he has finished that (3).
Because Google is registered as trusted identity provider the user will be able to login in with the Google account after he linked an existing ZITADEL Account or just registered a new one with the claims provided by Google (4)(5).

![Identity Brokering](/img/guides/identity_brokering.png)

## How to use external identity providers in ZITADEL

Configure external identity providers on instance level or just for one organization via [Console](/manage/console/instance-settings#identity-providers) or APIs.
The guides in this will help you to set up specific identity providers.
ZITADEL provides also templates to configure generic identity providers, which don't have a template.
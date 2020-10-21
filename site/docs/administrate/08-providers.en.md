---
title: Identity Providers
---

### What are Identity Providers

Identity providers or in short idp's are external systems to which **ZITADEL** can create a **federation** or use their **directory service**.
Normally federation uses protocols like [OpenID Connect 1.0](https://openid.net/connect/), [OAuth 2.0](https://oauth.net/2/) and [SAML 2.0](http://docs.oasis-open.org/security/saml/Post2.0/sstc-saml-tech-overview-2.0.html).

Some examples include:

#### Social Providers

- Google Account
- Microsoft Live Account
- Apple ID
- GitHub
- GitLab
- ...

#### Enterprise Providers**

- Azure AD Tenant
- Gsuite hosted domain
- ...

### Generic

- ADFS
- ADDS
- Keycloak
- LDAP

### What is Identity Brokering

ZITADEL supports the usage as identity broker, by linking multiple external idp's into one user.
With identity brokering the client which relies on ZITADEL does not need to care about the linking of identity.

### Manage Identity Providers

> Screenshot here

### Federation Protocols

Currently supported are the following protocols.

- OpenID Connect 1.0
- OAuth 2.0

SAML 2.0 will follow later on.

### Storage Federation

> This is a work in progress.

Storage federation is a means of integrating existing identity storage like [LDAP](https://tools.ietf.org/html/rfc4511) and [ADDS](https://docs.microsoft.com/en-us/windows-server/identity/ad-ds/get-started/virtual-dc/active-directory-domain-services-overview).
With this process **ZITADEL** can authenticate users with LDAP Binding and SPNNEGO for ADDS. It is also possible to synchronize the users just-in-time or scheduled.

#### Sync Settings

Here we will document all the different sync options

- Readonly
- Writeback
- just-in-time sync
- scheduled sync

> TBD

### Audit identity provider changes

> Screenshot here

---
title: Identity Providers
---

### What are Identity providers

Identity providers or in short idp's are external systems to which ZITADEL can create a federation.
Normally this federation uses protocols like OpenID Connect 1.0 / OAuth 2.0 and SAML 2.0.

Some examples include:

**Social Providers**

- Google Account
- Microsoft Live Account
- Apple ID
- GitHub
- GitLab
- ...

**Enterprise Providers**

- Azure AD Tenant
- Gsuite hosted domain
- ...

**Generic**

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

Upcoming is SAML 2.0

> Screenshot here

### Audit identity provider changes

> Screenshot here

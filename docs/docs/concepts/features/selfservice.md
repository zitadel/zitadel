---
title: Self-Service
---

ZITADEL allows users to perform certain tasks themselves.
For these tasks we either provide an user interface, or the tasks can be initiated or completed through our APIs.

It is important to understand that, depending on your use case, there will exist different user-types that want to perform different actions:  

- `User` are the end-users of your application. Users should be able to perform tasks like register/join, update their profile, manage authenticators etc.There are certain actions that can be executed pre-login, yet others require the user to have a valid session.
- `Manager` are users with a [special manager role within ZITADEL](http://localhost:3000/docs/concepts/structure/managers) and can perform administrative actions such as system configuration or granting access rights to users.

:::info
It is important to note that a `Manager` is not simply an administrative user, but can be used to create much more advanced self-service scenarios.
For example you can create an organization and assign a user from that organization the Manager Role `ORG_OWNER`.
Given this role the user could perform actions like configuring their own SSO/Identity Provider, set security policy for their organization, or assign roles to other users.
:::

## Registration

Unauthenticated users (pre-login).

- Email

## Login

Unauthenticated users (pre-login).

### Browser

- Explain hosted login
- Explain basic flow
- Branding: Trigger based on domain/org, or primary domain scope
- Link to Login user Guide, Scopes, etc.

### Mobile Applications

- Embedded Browser
- Redirect protocol (guide?)

### SSO / External IdP

### Social Login

### APIs

### Others

- Games etc.

## Logout

Authenticated users.

- End all sessions
- SLO

## Secrets

Authenticated users.

### Password reset

### Change password

### MFA / FIDO Passkeys

- Add & remove second factor
- Add & remove passwordless authenticator

## Profile

Authenticated users.

### Update Information

### Email Verification

### Phone Verification

### Account Linking

- Add external identity providers

## Managers

[Roles](/docs/concepts/structure/managers#roles)

Can be human users or also service users (eg, to manage programmatically)
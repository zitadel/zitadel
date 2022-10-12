---
title: Users
---

## Types of users

ZITADEL supports authentication and authorization for different user types.
We can mainly differentiate between Human Users and Machine Users.
We typically call human users simply "Users" and machine users "Service Users".

### Human users

Human users typically logon with an interactive login.
This means that an application redirects a user to a website ("login page") where the user can provide her credentials.
ZITADEL handles the authentication and provides the application only with a token that verifies the authentication process.

### Service users

Service users are for machine-to-machine communication and you would typically secure backend services.
For example in ZITADEL you would require an authenticated Service User to access the Management API.
The main difference between human and machine users is the type of credentials that can be used for authentication: Human users typically logon via an login prompt, but Machine users require a non-interactive logon process.

### Managers

Any user, human or service user, can be given a [Manager](/docs/concepts/structure/managers) role.
Given a manager role, a user is not only an end-user of ZITADEL but can also manage certain aspects of ZITADEL itself.

## Uniqueness

- Can only exist in one organization
  - identified by loginname
  - user grant
  - moving identities

## Identity linking

- Identity linking:
  - add external idps to an identity
  - auto-linking

More about how to manage your users read our [users guide](../../guides/manage/console/users).
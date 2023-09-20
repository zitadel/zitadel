---
title: ZITADEL Users
sidebar_label: Users
---

## Types of users

ZITADEL supports authentication and authorization for different user types.
We can mainly differentiate between Human Users and Machine Users.
We typically call human users simply "Users" and machine users "Service Users".

### Human users

Human users typically logon with an interactive login.
This means that an application redirects a user to a website ("login page") where the user can provide the credentials.
ZITADEL handles the authentication and provides the application with a token that verifies the authentication process.

### Service users

Service users are for machine-to-machine communication and you would use those typically to access secure backend services.
For example in ZITADEL you would require an authenticated Service User to access the Management API.
The main difference between human and machine users is the type of credentials that can be used for authentication: Human users typically logon via an login prompt, but Machine users require a non-interactive logon process.

### Managers

Any user, human or service user, can be given a [Manager](/concepts/structure/managers) role.
Given a manager role, a user is not only an end-user of ZITADEL but can also manage certain aspects of ZITADEL itself.

## Constraints

Users can only exist within one [organization](/concepts/structure/organizations).
It is currently not possible to move users between organizations.

User accounts are uniquely identified by their `id` or `loginname` in combination of the `organization domain` (eg, `road.runner@acme.zitadel.local`).
You can use the same email address for different user accounts.

## Where to store users

Depending on your [scenario](/guides/solution-scenarios/introduction), you might want to store all users in one organization (CIAM / B2C) or create a new organization for each logical group of users, e.g. each business customer (B2B).
With a project grant, you can delegate the access management of an organization's project to another organization.
You can also create a user grant to allow single users to access projects from another organization.
This is also an alternative to cases where you might want to move users between organizations.

## Identity linking

When using external identity providers (ie. social login, enterprise SSO), a user account will be created in ZITADEL.
The external identity will be linked to the ZITADEL account.

You can link multiple external accounts to a ZITADEl account.
If login with "Username / Password" (ie. local account) is enabled and you have configured external IDPs, the user can decide if she wants to login with an external IDP or the local account.
When only one external identity provider is configured and login with "Username / Password" is disabled, then the user is immediately redirected to the external identity provider.

More about how to manage your users read our [users guide](../../guides/manage/console/users).
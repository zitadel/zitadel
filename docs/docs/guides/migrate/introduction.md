---
title: Migrate to ZITADEL
sidebar_label: Introduction
---

This section of guides shows you how to migrate from your current auth system to ZITADEL.
The individual guides in this section should give you an overview of things to consider before you start the migration.

When moving from a previous auth solution to ZITADEL, it is important to note that some decisions and features are unique to ZITADEL.
Without duplicating too much content here are some important features and patterns to consider in terms of solution architecture.
You can read more about the basic structure and important concepts of ZITADEL in our [concepts section](/docs/concepts/).

## Multi-tenancy architecture

Multi-tenancy in ZITADEL can be achieved through either [Instances](/docs/concepts/structure/instance) or [Organizations](/docs/concepts/structure/organizations).
Where instances represent isolated ZITADEL instances, Organizations provide a more permeable approach to multi-tenancy.

In most cases, when you want to achieve multi-tenancy, you use Organizations. Each organization can have their own set of Settings (eg, Security Policies, IDPs, Branding), Managers, and Users.
Please also consult our guide on [Solution Scenarios](/docs/guides/solution-scenarios/introduction
) for B2C and B2B for more details.

## Delegated access management

Some solutions, that offer multi-tenancy, require you to copy applications and settings to each tenant and manage changes individually.
ZITADEL works differently by using [Granted Projects](/docs/concepts/structure/granted_projects).

Projects can be granted to [Organization](/docs/concepts/structure/projects#granted-organizations) or even to individual users.
You can think of it as a logical link to a Project, which can be used by the receiving Organization or User as if it was their own project, except privileges to modify the Project itself.

Delegated access management is a great way of keeping the management overhead low and enabling [self-service](/docs/concepts/features/selfservice#managers-in-delegation) for Organizations to manage their own Settings and Authorizations.

## Actions

ZITADEL [Actions](/docs/apis/actions/introduction) is the key feature to customize and create workflows and change the default behavior of the platform.

You define custom code that should be run on a specific Trigger.
A trigger could be the creation of a new user, getting profile information about a user, or a login attempt.

With the [HTTP module](/docs/apis/actions/modules) you can even make calls to third party systems for example to receive additional user information from a backend system, or triggering a workflow in other systems (Webhook).

You can also create custom claims or manipulate tokens and information on the userinfo endpoint by using the [Complement Token Flow](/docs/apis/actions/complement-token).
This might be required, if an application expects roles/permissions in a certain format or additional attributes (eg, a backend user-id) as claims.

## Metadata

You can store arbitrary key-value pairs of data on objects such as Users or Organizations.
Metadata could link a user to a specific backend user-id or represent an "organizational unit" for your business logic.
Metadata can be access directly with the correct [scopes](/docs/apis/openidoauth/scopes#reserved-scopes) or transformed to custom claims (see above).

## Migrating resources

### Migrating users

Migrating users with minimal impact on users can be a challenging task.
We provide some more information on migrating users and secrets in [this guide](./users.md).

### Migrating clients / applications

After you have set up or imported your applications to ZITADEL, you need to update your client's configurations, such as issuer, clientID or credentials.
It is not possible to create an application with a pre-defined clientID or import existing credentials.

## Technical considerations

### Batch migration

**Batch migration** is the easiest way, if you can afford some minimal downtime to move all users and applications over to ZITADEL.
See the [User guide](./users.md) for batch migration of users.

### Just-in-time migration

In case all your applications depend on ZITADEL after the migration date, and ZITADEL is able to retrieve the required user information, including secrets, from the legacy system, then the recommended way is to let **ZITADEL orchestrate the user migration just-in-time**:

- Create a pre-authentication [Action](/docs/apis/actions/introduction) to request user data from the legacy system and create a new user in ZITADEL.
- Optionally, create a post-authentication Action to flag successfully migrated users in your legacy system

For all other cases, we recommend that the **legacy system orchestrates the migration** of users to ZITADEL for more flexibility:

- Update your legacy system to create a user in ZITADEL on their next login, if not already flagged as migrated, by using our APIs (you can set the password and a verified email)
- Redirect migrated users with a login hint in the [auth request](/docs/apis/openidoauth/authrequest.mdx) to ZITADEL to pre-select the user

In this case the migration can also be done as an import job or also allowing to create user session in both the legacy auth solution and ZITADEL in parallel with identity brokering:

- Setup ZITADEL to use your legacy system as external identity provider (note: you can also use JWT-IDP, if you only have a token).
- Configure your app to use ZITADEL, which will redirect users automatically to the external identity provider to login.
- A session will be created both on the legacy system and ZITADEL
- If a user does not exist already in ZITADEL you can auto-register new users and use an Action to pull additional information (eg, Secrets) from your legacy system. Note: ZITADEL links external identity information to users, meaning you can have users use both a password and external identity providers to login with the same user.

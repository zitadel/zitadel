---
title: Migrate to ZITADEL
sidebar_label: Introduction
---

This section of guides shows you how to migrate from your current auth system to ZITADEL.
The individual guides in this section should give you an overview of things to consider before you start the migration.

When moving from a previous auth solution to ZITADEL, it is important to note that some decisions and features are unique to ZITADEL.
Without duplicating too much content here are some important features and patterns to consider in terms of solution architecture.

## Multi-Tenancy Architecture

Multi-tenancy in ZITADEL can be achieved through either [Instances](/docs/concepts/structure/instance) or [Organizations](/docs/concepts/structure/organizations).
Where instances represent isolated ZITADEL instances, Organizations provide a more permeable approach to multi-tenancy.

In most cases, when you want to achieve multi-tenancy, you use Organizations. Each organization can have their own set of Settings (eg, Security Policies, IDPs, Branding), Managers, and Users.
Please also consult our guide on [Solution Scenarios](/docs/guides/solution-scenarios/introduction
) for B2C and B2B for more details.

## Delegated Access Management

Some solutions, that offer multi-tenancy, require you to copy applications and settings to each tenant and manage changes individually.
ZITADEL works differently by using [Granted Projects](/docs/concepts/structure/granted_projects).

Projects can be granted to [Organization](/docs/concepts/structure/projects#granted-organizations) or even to individual users.
You can think of it as a logical link to a Project, which can be used by the receiving Organization or User as if it was their own project, except privileges to modify the Project itself

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

## Migrating users

Migrating users with minimal impact on users can be a challenging task.
We provide technical considerations and solution scenarios in [this guide](./users.md).
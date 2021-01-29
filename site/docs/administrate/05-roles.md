---
title: Roles
---

### What are Roles

**ZITADEL** lets [projects](administrate#projects) define their **role based access control**.

**Roles** can be consumed by the [clients](administrate#clients) which exist within a specific [project](administrate#projects).

For more information about how **roles** can be consumed, have a look the protocol specific information.

- [OpenID Connect / OAuth](quickstarts#How_to_consume_authorizations_in_your_application_or_service)

### Manage Roles

Each **role** consist of three fields.

| Field        | Description                                                                  | Example                                          |
|:-------------|:-----------------------------------------------------------------------------|--------------------------------------------------|
| Key          | This is the **Roles** actual name which can be used to verify the users roles.                                            | User                                             |
| Display Name | A descriptive text for the purpose of the **Role**                           | User is the default role provided to each person |
| Group        | The group field allows to group certain roles who belong in the same context | User and Admin in the group **default**          |

### Granting Roles

To give someone (or somewhat) access to a [project's](administrate#projects) resources and services **ZITADEL** provides two processes: **Roles** can either be granted to [users](administrate#Users) or to [organisations](administrate#Organisations).

#### Grant Roles to Organisations

The possibility to grant **roles** to an [organisation](administrate#Organisations) is intended as "delegation" so that a [organisation](administrate#Organisations) can on their own grant access to [users](administrate#Users).

For example a **service provider** could grant the **roles** `user`, and `manager` to an [organisation](administrate#Organisations) as soon as they purchases his service. This can be automated by utilising a [service user](administrate#Manage_Service_Users) in the **service providers** business process.

> Screenshot here

#### Grant Roles to Users

By granting **roles** to [users](administrate#Users), be it [humans or machines](administrate#Human_vs_Service_Users), this [user](administrate#Users) receives the authorization to access a project's resources.

> Screenshot here

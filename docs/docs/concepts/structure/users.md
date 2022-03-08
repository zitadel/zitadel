---
title: Users
---

Users belong to one, and only one, organization in ZITADEL. It is not possible to move a user to another organization, at least for the moment.

![Overview ZITADEL Organizations](/img/concepts/objects/organization.png)

## Managing Users

Users can be created, changed, and deleted via our APIs or manually in Console.

### Self-service

Depending on the organization's policies users may self-register and thus create accounts on their own. Users can self-manage their user information or authentication methods. A self-service UI is provided in Console, yet you can also integrate our APIs with your user interface.

### Metadata

ZITADEL allows storing arbitrary key-value metadata on a user object. This can be used to store own identifiers or additional information about the user. We recommend to keep only required information and user other patterns such as distributed claims for more advanced use cases. 

## User Grant

Instead of an [Organization Grant](/docs/guides/basics/projects#exercise---grant-a-project) you can also grant roles from a project to an individual user from another Organization. This feature is called User Grant.

## Human vs. Machine (Service Users)

ZITADEL supports human an machine users. We call human users simply "Users" and machine users "Service Users".

With Service Users you would typically secure backend services. For example in ZITADEL you would require an authenticated Service User to access the Management API. The main difference between human and machine users is the type of credentials that can be used for authentication: Human users typically logon via an login prompt, but Machine users require a non-interactive logon process.

---
title: Projects
---

### What are projects

The idea of projects is to have a vessel for all components who are closely related to each other.
In ZITADEL all clients located in the same project share their roles, grants and authorizations.
From a access management perspective you manage who has what role in the project and your application consume this information.
A project belongs to exactly one organisation.

**Clients**

Clients are described here [What are clients](###what_are_clients)
Basically these are you applications who initiate the authorization flow

**Roles**

Roles (or Project Roles) is a mean of managing users access rights for a certain project. 
These roles are opaque for ZITADEL and have no weight in relation to each other. 
So if a user has two roles, admin and user in a certain project, the information will be treated additive.

**Grants**

With ZITADEL it is possible to give third parties (other organisations) the possibility to manage certain roles on their own.
To achieve this the owner of a project can grant (some could say delegate) certain roles or all roles to a organisation.
After granting that organisation it can manage on its own which user has what roles.
This feature is especially useful for service providers, because they are able to establish a great self-service culture for their business customers.

**Authorizations** 

#### Project vs. granted Project

The simple difference of a project vs a granted project is that a project belongs to your organisation and the granted project belongs to a third party who did grant you some rights to manage certain roles of their project.
To make it more easily to differentiate ZITADEL Console displays these both as separate menu in the project section.

### Setup new project

> Screenshot here

#### Define project specific roles

> Screenshot here

### Grant project to a third party

> Screenshot here

### Audit project changes

> Screenshot here

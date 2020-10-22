---
title: Projects
---

### What are projects

The idea of projects is to have a vessel for all components who are closely related to each other.
In ZITADEL all clients located in the same project share their roles, grants and authorizations.
From an access management perspective you manage who has what role in the project and your application consumes this information.
A project belongs to exactly one organisation.
The attribute project role assertion defines, if the roles should be integrated in the tokens without sending corresponding scope (urn:zitadel:iam:org:project:role:{rolename})
With the project role check you can define if a user should have a requested role to be able to logon.

**Clients**

Clients are described here [What are clients](administrate#What_are_clients)
Basically these are your applications who initiate the authorization flow.

**Roles**

[Roles (or Project Roles)](administrate#Roles) is a means of managing users access rights for a certain project.
These [roles](administrate#Roles)  are opaque for ZITADEL and have no weight in relation to each other.
So if a [user](administrate#Users) has two roles, admin and user in a certain project, the information will be treated additive.

**Grants**

With ZITADEL it is possible to give third parties (other organisations) the possibility to manage certain roles on their own.
To achieve this the owner of a project can grant (some could say delegate) certain roles or all roles to an organisation.
After granting that organisation it can manage on its own which user has what roles.
This feature is especially useful for service providers, because they are able to establish a great self-service culture for their business customers.

**Authorizations** 

#### Project vs. granted Project

The simple difference of a project vs a granted project is that a project belongs to your organisation and the granted project belongs to a third party who did grant you some rights to manage certain roles of their project.
To make it more easier to differentiate, ZITADEL Console displays these both as separate menu in the project section.

### Manage a project

#### Create a project

To create your project go to [https://console.zitadel.ch/projects](https://console.zitadel.ch/projects)

<img src="img/console_projects_empty.png" alt="Manage Projects" width="1000px" height="auto">

Create a new project with a name which explains what's the intended use of this project.

<img src="img/console_projects_my_first_project.png" alt="Manage Projects" width="1000px" height="auto">

#### RBAC Settings

- Authorisation Check option (Check if the user at least has one role granted)
- Enable Project_Role Assertion (if this is enabled assert project_roles, with the config of the corresponding client)

#### Define project specific roles

> Screenshot here

### Grant project to a third party

> Screenshot here

### Audit project changes

> Screenshot here

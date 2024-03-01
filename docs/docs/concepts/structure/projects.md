---
title: ZITADEL Projects
sidebar_label: Projects
---

# Project

import ProjectDescription from './\_project_description.mdx';

<ProjectDescription name="ProjectDescription" />

To learn how to set up a project read this console guide [here](../../guides/manage/console/projects.mdx).

## Applications

Applications define the different clients, that share the same roles.
At the moment we support OIDC and almost every OAuth2 client. We'll be expanding this with SAML shortly.
Go to [Applications](./applications) for more details.

## Granted Organizations

To enable another organization to use a project, the organization needs a grant to the project.
Only the selected roles will be available to the granted organization.

The granted organization will be able to manage the authorizations of their users for the granted project by themselves in their organization.

More about granted projects: [Granted Projects](./granted_projects)

## Roles

A role consists of different attributes. Only the key is relevant to the authorization and must therefore be unique.
The display name is only to provide a human-readable name if needed.
And the group should enable a better handling in ZITADEL console, like give a user all the roles of a specific group. (Not implemented yet)

All applications in a project share the roles. Read more about roles [here](../../guides/manage/console/roles)

## Default Project

When creating a new ZITADEL instance you will find an automatically created project on the first created organization.
This default project does represent the ZITADEL project and is used to secure the different apps shipped with ZITADEL.
This includes the ZITADEL Management Console and APIs.

We do not recommend changing any settings or using this project for anything else, as this could influence the behavior of ZITADEL.

The default name of the project is "ZITADEL", this might differ on self-hosted instances, when you change the default name.

![Default Project](/img/guides/solution-scenarios/console-default-project.png)

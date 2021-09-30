---
title: Projects
---

# Project

import ProjectDescription from './_project_description.mdx';

<ProjectDescription name="ProjectDescription" />

## Project Settings

On default the login screen will be shown in the private labeling settings of the system (e.g zitadel.ch).
With the [primary domain scope](../../../apis/openidoauth/scopes#reserves-scopes) it is possible to trigger the setting of the given organization. 
But this will also restrict, the login to user of the given organization.  

With the private labeling setting it is possible to choose which settings should trigger.

| Setting | Description |
| --- | --- |
| Unspecified | If nothing is specified the default will trigger. (System settings zitadel.ch) |
| Enforce project resource owner policy | This setting will enforce the private labeling of the organization (resource owner) of the project through the whole login process. |
| Allow Login User resource owner policy | With this setting first the private labeling of the organization (resource owner) of the project will trigger. As soon as the user and its organization (resource owner) is identified by ZITADEL, the settings will change to the organization of the user. |

## Applications

Applications define the different clients, that share the same roles. 
At the moment we support OIDC and almost every OAuth2 client. We'll be expanding this with SAML shortly.
Go to [Applications](./applications) for more details.

## Granted Organizations

To enable another organization to use a project, the organization needs a grant to the project.
Only the selected roles will be available to the granted organization.

The granted organization will be able to manage the authorizations of his users for the granted project by himself in his own organization.

More about granted projects: [Granted Projects](./granted_projects)

## Roles

A role consists of different attributes. For the authorization only the key is relevant, which must therefore be unique.
The display name is only to provide a human-readable name if needed. 
And the group should enable a better handling in ZITADEL console, like give a user all the roles of a specific group. (Not implemented yet)

All applications in a project share the roles.

### Role specific Project Settings

| Setting | Description |
| --- | --- |
| Assert roles on authentication | If this setting is enabled role information is sent from userinfo endpoint and depending on your application settings in tokens and other types.  |
| Check roles on authentication | If set, users are only allowed to authenticate if any role is assigned to their account. |
| Check for project on authentication | The user will only be able to authenticate if his organization is the owner or has a grant to the project. |

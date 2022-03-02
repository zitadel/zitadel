---
title: Projects

---

|                   |                                                                                                                                                                                                                          |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Description       | Learn the basics about applications, roles and authorizations, then learn how projects allow you to group these together.                                                                                                       |
| Learning Outcomes | In this module you will: <ul><li>Learn about projects and granted projects</li><li>Create a new project</li><li>Creating simple roles and authorizations</li><li>Create an organization grant for your project</li></ul> |
| Prerequisites     | <ul><li>ZITADEL organizations</li><li>Role Based Access Management (RBAC)</li></ul>                                                                                                                                      |

## What is a project?

import ProjectDescription from '../../concepts/structure/_project_description.mdx';

<ProjectDescription name="ProjectDescription" />

The goal of this module is to give you an overview of how to manage access rights and delegate role management to third parties.
This module tries not to dive too deeply into the details.

So, let’s create a straightforward example project.

## Exercise - Create a simple project

1. Visit <https://console.zitadel.ch/projects>, or select **Projects** within your organization.
2. Select **+ Create New Project**.

![Empty Project](/img/console_projects_empty.png)

3. Enter the name `My first project` and continue.

Let’s make this more interesting by adding some basic roles and authorizations to your project, then confirm the scope of the roles and authorizations.

4. Jump to the section **ROLES**. Create two new roles with the following values

* Key: reader
* Display Name: Reader
* Group: user

and

* Key: editor
* Display Name: Editor
* Group: user

![Add New Roles](/img/console_projects_add_new_roles.gif)

Now, you can add roles to your own user, or you can create a new user.
To create a new user:

1. Go to **Users** and select **New**.
2. Enter the required contact details and select **Create**.

![Create new user](/img/console_users_create_new_user.gif)

To grant users certain roles, you need to create authorizations.
Go back to the project, and jump to the section **AUTHORIZATIONS**.

![Verify your authorization](/img/console_projects_create_authorization.gif)

To verify the role granted to the user:

1. Select **Users** from the navigation menu.
2. Select the user `Coyote`.
3. Scroll down to the section **AUTHORIZATION**. There you should be able to verify that the user has the role `reader` for your project `My first project`.

![Organization grant](/img/console_projects_authorization_created.png)

4. Now create another project (eg. “My second project”).
5. Verify that there are no roles or authorizations on your second project.

## What is a granted project?

import GrantedProjectDescription from '../../concepts/structure/_granted_project_description.mdx';

<GrantedProjectDescription name="GrantedProjectDescription" />

## Exercise - Grant a project

1. Visit the project that you created before. In the section GRANTED ORGANIZATIONS, select 
**New**.
2. Enter the domain `acme.caos.ch`. Search the organization and continue to the next step.
3. Select some roles that you would like to grant to the organization ACME. Confirm the change.
4. You should now see ACME-CAOS in the section GRANTED ORGANIZATIONS

![Grant a project](/img/projects_create_org_grant_caos2acme.gif)

## Knowledge Check (2)

* You can set up multiple projects within an organization to manage scope.
    - [ ] yes
    - [ ] no
* Authorizations define more detailed access rights within your application.
    - [ ] yes
    - [ ] no
* Your projects, as well as projects granted to your organization, are visible within the Tab Projects of your organization.
    - [ ] yes
    - [ ] no

<details>
    <summary>
        Solutions
    </summary>

* You can setup multiple projects within an organization to manage scope
    - [x] yes
    - [ ] no
* Authorizations are define more detailed access rights within your application
    - [ ] yes
    - [x] no (Authorizations link users to certain roles)
* Your projects as well as projects granted to your organization are visible within the Tab Projects of your organization
    - [ ] yes
    - [x] no (Projects and Granted Projects are shown on different tabs)

</details>

## Summary (2)

* Use projects to manage the scope of roles, authorizations and applications
* Create and assign roles to users of your organization within your project
* Use project grants to enable other organizations to self-manage access rights (roles) to your applications

Where to go from here:

* Manage roles for your project
* Grant roles to other organizations or users
* Service Users
* Manage IAM Roles
* Setup a SaaS Application with granted projects (Learning Path)

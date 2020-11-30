---
title: Projects
---

### What are projects

The idea of projects is to have a vessel for all components who are closely related to each other.
In ZITADEL all clients located in the same project share their roles, grants and authorizations.
From an access management perspective you manage who has what role in the project and your application consumes this information. A project belongs to exactly one [organisation](administrate#Organisations).

The attribute `project role assertion` defines, if the roles should be integrated in the tokens without sending corresponding scope (`urn:zitadel:iam:org:project:role:{rolename}`)

With the project role check you can define, if a user should have a requested role to be able to logon.

**Clients**

These are your applications who initiate the authorization flow (see  [What are clients](administrate#What_are_clients)). 

**Roles**

[Roles (or Project Roles)](administrate#Roles) are a means of managing users access rights for a certain project. These [roles](administrate#Roles) are opaque for ZITADEL and have no weight in relation to each other. 

As example, if [user](administrate#Users) has two roles, `admin` and `user` in a certain project, the information will be treated additive. There is no meaning or hierarchy implied by these roles.

**Grants**

With ZITADEL it is possible to give third parties (other organisations) the possibility to manage certain roles on their own. As a service provider, you will find this feature useful, as it allows you to establish a self-service culture for your business customers.

The owner of a project can grant (some would say "delegate") certain roles or all roles to another organisation. The target organization can then indipendently manage the assignment of their users to  the role within the [granted project](administrate#Project_vs_granted_Project).

**Authorizations**

> TODO, Link to authorizations

#### Project vs. granted Project

A project belongs to your organisation. You can [grant certain roles](administrate#Grant_project_to_a_third_party) to another organisation. A granted project, on the other hand, belongs to a third party, granting you some rights to manage certain roles of their project. 

To make it more easier to differentiate ZITADEL Console displays these both as separate menu in the project section.

### Manage a project

#### Create a project

To create your project go to [https://console.zitadel.ch/projects](https://console.zitadel.ch/projects)

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_projects_empty.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_projects_empty.png" itemprop="thumbnail" alt="Manage Projects" />
        </a>
        <figcaption itemprop="caption description">Manage Projects</figcaption>
    </figure>
</div>

Create a new project with a name which explains what's the intended use of this project.

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_projects_my_first_project.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_projects_my_first_project.png" itemprop="thumbnail" alt="Manage Projects" />
        </a>
        <figcaption itemprop="caption description">Manage Projects</figcaption>
    </figure>
</div>

#### RBAC Settings

- Authorisation Check option (Check if the user at least has one role granted)
- Enable Project_Role Assertion (if this is enabled assert project_roles, with the config of the corresponding client)

#### Define project specific roles

> Screenshot here

### Grant project to a third party

> Screenshot here

### Manage Project Authorisations

> Screenshot here

### Manage Project ZITADEL Roles

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_project_manage_roles_1.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_project_manage_roles_1.png" itemprop="thumbnail" alt="Manage ZITADEL Roles 1" />
        </a>
        <figcaption itemprop="caption description">Manage ZITADEL Roles 1</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_project_manage_roles_2.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_project_manage_roles_2.png" itemprop="thumbnail" alt="Manage ZITADEL Roles 2" />
        </a>
        <figcaption itemprop="caption description">Manage ZITADEL Roles 2</figcaption>
    </figure>
</div>

### Audit project changes

> Screenshot here

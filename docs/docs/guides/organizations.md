---
title: Organizations
---


| | |
| --- | --- |
| Description | Learn how ZITADEL is structured around Organizations and how to create your organization and verify a domain to use with that new organization. |
| Learning Outcomes | In this module you will: <ul><li>Learn about organizations</li><li>Create a new organization</li><li>Verify your domain name </li></ul> |
|Prerequisites|None|

## What is an organization?

ZITADEL is organized around the idea that 
* Multiple organizations share the same system. In this case multiple organizations share the same service, zitadel.ch
* organizations can grant each other rights to self-manage certain aspects of the IAM (eg, roles for access management)
* organizations are vessels for users and projects

![Overview ZITADEL Organizations](/img/zitadel_organizations.png)

Organizations in ZITADEL are therefore comparable to tenants of a system or organizational units of a directory based system.

You can use projects within your organization to manage the security context of closely related components, such as roles, grants and authorizations for multiple clients. You can set up multiple projects within your organization. 

ZITADEL allows you to give other organizations permission to manage certain aspects of a project within your organization on their own. This means you could set up a project with roles that should exist within your service/software, but allow another organization to allocate the roles to users within their own organization. As a service provider, you will find this feature useful, as it allows you to establish a self-service culture for your business customers.

![Organization Grant](/img/zitadel_organization_grant.png)

Each organization has its own pool of usernames, which includes human and service users, for its domain (`{username}@{domainname}.{zitadeldomain}`). A username is unique within your organization. You can configure ZITADEL to use your own domain, and simplify user experience (`{loginname}@{yourdomain.tld}`).

There are several more modules in our documentation to go into more detail regarding organization management, projects, clients, and users. But first let’s create a new organization and verify your domain name.

## Exercise - Create a new organization

To register your organization and create a user for ZITADEL, visit zitadel.ch or directly https://accounts.zitadel.ch/register/org and fill in the required information.

![Register new Organization](/img/console_org_register.png)

If you already have an existing login for zitadel.ch, you need to visit the console, then click on your organization’s name in the upper left corner, and then select “New organization”.

![Select Organization](/img/console_org_select.png)

## How ZITADEL handles usernames

As we mentioned before, each organization has its own pool of usernames, which includes human and service. 

This means that, for example a user with the username road.runner, can only exist once in an organization called ACME. ZITADEL will automatically generate a "logonname" for each  consisting of `{username}@{domainname}.{zitadeldomain}`, in our example road.runner@acme.zitadel.ch.

When you verify your domain name, then ZITADEL will generate additional logonames for each user with the verified domain. If our example organization would own the domain acme.ch and verify within the organization ACME, then the resulting logonname in our example would be road.runner@acme.ch in addition to the already generated road.runner@acme.zitadel.ch. The user can now use either logonname to authenticate with your application.

## Domain verification and primary domain

Once you have successfully registered your organization, ZITADEL will automatically generate a domain name for your organization (eg, acme.zitadel.ch). Users that you create within your organization will be suffixed with this domain name.

You can improve the user experience, by suffixing users with a domain name that is in your control. For that you can prove the ownership of your domain, by DNS or HTTP challenge.

An organization can have multiple domain names, but only one domain can be primary. The primary domain defines which login name ZITADEL displays to the user, and what information gets asserted in access_tokens (`preferred_username`).

Please note that domain verification also removes the logonname from all users, who might have used this combination in the global organization (ie. users not belonging to a specific organization). Relating to our example with acme.ch: If a user ‘coyote’ exists in the global organization with the logonname coyote@acme.ch, then after verification of acme.ch, this logonname will be replaced with `coyote@{randomvalue.tld}`. ZITADEL will notify users affected by this change.

## Exercise - Verify your domain name

1. Browse to your organization
2. Click **Add Domain**
3. To start the domain verification click the domain name and a dialog will appear, where you can choose between DNS or HTTP challenge methods.
4. For example, create a TXT record with your DNS provider for the used domain and click verify. ZITADEL will then proceed an check your DNS.
5. When the verification is successful you have the option to activate the domain by clicking **Set as primary**

![Verify Domain](/img/console_verify_domain.gif)

<figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
    <a href="img/console_verify_domain.gif" itemprop="contentUrl" data-size="1920x1080">
        <img src="img/console_verify_domain.gif" itemprop="thumbnail" alt="console_verify_domain" />
    </a>
    <figcaption itemprop="caption description">Verify Domain</figcaption>
</figure>

> **_Please note:_** Do not delete the verification code, as ZITADEL will re-check the ownership of your domain from time to time

## Knowledge Check

* Users exist only within projects or clients
    - [ ] yes
    - [ ] no
* User can only login with `{username}@{domainname}.{zitadeldomain}`
    - [ ] yes
    - [ ] no
* You can delegate access management self-service to another organization
    - [ ] yes
    - [ ] no

<details>
    <summary>
        Solutions
    </summary>

* Users exist only within projects or clients
    - [ ] yes
    - [x] no (users exist within organizations)
* User can only login with `{username}@{domainname}.{zitadeldomain}`
    - [ ] yes
    - [x] no (You can validate your own domain and login with `{loginname}@{yourdomain.tld}`)
* You can delegate access management self-service to another organization
    - [x] yes
    - [ ] no
    
</details>

## Summary

* Create your organization and a new user by visiting zitadel.ch
* Organizations are the top-most vessel for your IAM objects, such as users or projects
* Verify your domain in the Console to improve user experience; remember to not delete the verification code to allow recheck of ownership
* You can delegate certain aspects of your IAM to other organizations for self-service

Where to go from here: 
* Create a project
* Setup Passwordless MFA
* Manage ZITADEL Roles
* Grant roles to other organizations or users
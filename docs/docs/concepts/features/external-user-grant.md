---
title: External user grant
---

ZITADEL's external user grant is a feature that allows you to grant access to projects within your organization to users from other organizations.
This is useful in scenarios where you want to collaborate with external users without needing them to be part of your organization.
By using external user grants, you can streamline collaboration with external users while maintaining control over access to your projects within ZITADEL.

![](/img/concepts/features/external-user-grant.png)

## Where to store users

### Consumer Identity Management (CIAM) / Business-to-Consumer (B2C)

You might typically store all users in a single ZITADEL [organization](../structure/organizations) for managing customer accounts.
We recommend creating a second organization for your own team that also contains all the projects and applications that will be [granted](../structure/granted_projects) to the first organization with the B2C customer accounts.
Instead of duplicating user accounts for your team members in the B2C organization, you can create external user grants on the B2C organization.

### Multitenancy / Business-to-Business (B2B)

ZITADEL allows you to create separate [organizations](../structure/organizations) for each of your business partner or tenant.
There might be cases were users from one organization need access to projects from another organization.
You can create an external user grant, that allows the inviting organization to manage the roles for the external user.

## Project Grants vs. User Grants

Project grants are used to delegate access management of an entire project (or specific roles of the project) to another organization.

User grants provide a more granular approach, allowing specific users from external organizations to access your projects.

## Alternative to multiple user accounts

A user account is always unique across a ZITADEL instance.
In some use cases, external user grants are a simple way to allow users access to multiple tenants.

## References

* [API reference for user grants](/docs/apis/resources/mgmt/user-grants)
* [How to manage user grants through ZITADEL's console](/docs/guides/manage/console/roles#authorizations)
* [More about multi-tenancy with ZITADEL](https://zitadel.com/blog/multi-tenancy-with-organizations)

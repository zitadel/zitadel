---
title: ZITADEL Users
sidebar_label: Users
sidebar_position: 3
---

## Types of users

ZITADEL supports authentication and authorization for different user types.
We can mainly differentiate between Human Users and Machine Users.
We typically call human users simply "Users" and machine users "Service Users".

### Human users

Human users typically logon with an interactive login.
This means that an application redirects a user to a website ("login page") where the user can provide the credentials.
ZITADEL handles the authentication and provides the application with a token that verifies the authentication process.

Read more on how to [login users with ZITADEL](/docs/guides/integrate/login/login-users).

### Service users

Service users are for machine-to-machine communication and you would use those typically to access secure backend services.
For example in ZITADEL you would require an authenticated Service User to access the Management API.
The main difference between human and machine users is the type of credentials that can be used for authentication: Human users typically logon via an login prompt, but Machine users require a non-interactive logon process.

Learn how to [use service users](/docs/guides/integrate/service-users/authenticate-service-users) with ZITADEL.

### Managers

Any user, human or service user, can be given a [Manager](/concepts/structure/managers) role.
Given a manager role, a user is not only an end-user of ZITADEL but can also manage certain aspects of ZITADEL itself.

### Federated users

Federated users are identities that are managed by a third-party identity provider.
Users can login via an external identity provider, using [identity brokering](../features/identity-brokering) ("Single-sign-on").
Federated user [accounts are linked](../features/account-linking) to internal users to be able to assign roles and keep an audit trail.

### External users

In a multi-tenancy architecture, you might use [organizations](organizations) to separate user groups.
By using [external user grants](../features/external-user-grant) an organization is able to invite users from another organization.
These invited users are called external users.

## Considerations

### Uniqueness of users

Users can only exist within one [organization](/concepts/structure/organizations).
It is currently not possible to move users between organizations.

User accounts are uniquely identified by their `id` or `loginname` in combination of the `organization domain` (eg, `road.runner@acme.zitadel.local`).
You can use the same email address for different user accounts.

### How to structure user pools

Consider this general recommendation as a starting point:

- Create one organization ("default organization") for your own company
- Configure projects and applications in the default organization
- Structure users in organizations based on common domains that are self-managed (eg, company)
- Grant your projects to the organizations, allow Managers to give granted roles to their users

You might want to adjust this general setup based on your [scenario](/guides/solution-scenarios/introduction). 

One important consideration in the setup is that you can only have a domain once for an organization. If you have multiple teams working with the same email address, you might need to add them to one single organization that has the domain verified for the teams' domain.

For a CIAM / B2C setup, you might want to store all users in one organization and allow that organization to use a specific set of social logins.
In a multitenancy / B2B scenario, you might have thousands of smaller teams.

#### Hierarchy

There is no concept of hierarchies and inheritance based on users or organizations.
This is why we recommend to structure users along the smallest unit of groups.
You can use organization metadata or your own business logic to describe a hierarchy of organizations or user groups.

## References

- [Manage users in the Console](../../guides/manage/console/users)
- [ZITADEL APIs: Users](/docs/apis/resources/mgmt/users)
- [User onboarding and registration](/docs/guides/integrate/onboarding)

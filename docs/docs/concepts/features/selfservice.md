---
title: Self Service in ZITADEL
sidebar_label: Self Service
---

ZITADEL allows users to perform many tasks themselves.
For these tasks we either provide an user interface, or the tasks can be initiated or completed through ZITADEL's APIs.

It is important to understand that, depending on your use case, there will exist different user-types that want to perform different actions:

- `Users` are the end-users of your application. Like with any CIAM solution, users should be able to perform tasks like register/join, update their profile, manage authenticators etc. There are certain actions that can be executed pre-login, yet others require the user to have a valid session.
- `Managers` are users with a [special manager role](../../guides/manage/console/managers) within ZITADEL and can perform administrative actions such as system configuration or granting access rights to users.

All self-service interfaces are available in different [languages](/guides/manage/customize/texts#internationalization).

:::info
ZITADEL covers the typical "CIAM" self-service capabilities as well as delegated access management for multi-tenancy scenarios. Please refer to the section [Managers](#managers).
:::

## Registration

:::info
You can pre-select a given organization by passing the scope `urn:zitadel:iam:org:domain:primary:{domainname}` in the authorization request.
This will force users to register only with the specified organization.
Furthermore the branding and login settings (e.g. Social Login Providers) are directly shown to the user.
:::

### Local account

Allows anonymous users and authenticated users, who are not yet in the given organization, to create an account (register) themselves.

- Mandatory profile fields
- Set password
- Accept terms of service and privacy policy
- User receives an email with an one-time code for verification purpose
- User has to enter the one-time code to finish registration
- User can re-request a new one-time code

The user is prompted on first login to setup MFA.
This step can be made mandatory if MFA is enforced in the login policy.

### Existing Identity / SSO / Social Login

Anonymous users and authenticated users, who are not yet in the given organization, to register with an external identity provider.
An external identity provider can be a Social Login Provider or a pre-configured identity provider.

- Information from the external identity provider is used to pre-fill the profile information
- User can update the profile information
- An account is created within ZITADEL and linked with the external identity provider

#### Account Linking

When you login with an external identity provider, and the user does not exist in ZITADEL, then an autoregister flow is triggered. The user is presented with two options:

- Create a new account: A new account will be created as stated above
- Autolinking: The user is prompted to login with an existing [local account](#local-account). If successful, the existing identity from the external identity provider will be linked with the local account. A user can now login with either the local account or any of the linked external accounts.

## Login

### SSO / Social Logins

Given an external identity provider is configured on the instance or on the organization, then:

- the user will be shown a button for each identity provider as alternative to login with a [local account](#local-account)
- when clicking the button the user will be redirected to the identity provider
- after successful login the user will be redirected to the application

### Machines

Machine accounts can't use an interactive login but require other means of authentication, such as privately-signed JWT or personal access tokens.
Read more about [Service Users](/guides/integrate/service-users/authenticate-service-users) and recommended [OpenID Connect Flows](/guides/integrate/login/oidc/oauth-recommended-flows#different-client-profiles).

## Logout

Users can terminate the session for all their users (logout).
A client can also implement this, by calling the [specific endpoint](/apis/openidoauth/endpoints#end_session_endpoint).

## Profile

These actions are available for authenticated users only.
ZITADEL provides a self-service UI for the user profile out-of-the box under the path _$CUSTOM-DOMAIN/ui/console/users/me_.
You can also implement your own version in your application by using our APIs.

### Change password

Users can change their passwords.
The current password must be entered first.

### MFA / FIDO Passkeys

Users can setup and delete a second factor and FIDO Passkeys (Passwordless).
Available authenticators are:

- Time-based one-time password (TOTP) (Which are Authenticator Apps like Google/Microsoft Authenticator, Authy, etc.)
- One-time password sent as E-Mail
- One-time password sent as SMS
- FIDO Universal Second Factor (U2F) (Security Keys, Device, etc.)
- FIDO2 WebAuthN (Passkeys)

### Update Information

Users can change their profile information. This includes

- UserName
- First- and Last name
- Nickname
- Display Name
- Gender
- Language
- Email address
- Phone number

### Email Verification

Users can change their email address.
The user receives a one-time password and can verify control over the given email address.

### Phone Verification

Users can change their phone number.
The user receives a one-time password and can verify control over the given phone number.

### Identity providers

Users can create a connection between a [local user account](#local-account) and an [external identity](#existing-identity--sso--social-login).
The user can login with any of the linked accounts.
[Linking of external accounts](#account-linking) is done during the login process.

## Managers

It is important to note that a `Manager` is not simply an administrative user, but can be used to create much more advanced scenarios such as delegating administration of a whole organization to a user, acting then as administrator and permission manager of that user group.

Thus we will explain service for two very common scenarios in ZITADEL:

- `Managers in isolation`: Granting administrative permissions within a single organization context.
- `Managers in delegation`: Granting administrative permissions to a user from a different organization where the organizations depend on each other

A list of [Manager Roles](../../guides/manage/console/managers#roles) is available with a description of permissions.
Managers can be assigned to both human users and service users eg, for managing certain tasks programmatically.

### Managers in isolation

An user with the Manager roles `IAM_OWNER` or `ORG_OWNER` might want to assign other users from their organization elevated permissions to handle certain aspects of the IAM tasks.
This could be permission to assign authorizations within this isolated organization (`ORG_USER_MANAGER`) or handling setup of projects and applications (`PROJECT_OWNER`).

### Managers in delegation

In a setup like described in the [B2B Scenario](/guides/solution-scenarios/b2b), there exists an organization of the project owner and a customer organization.
The project is granted to the customer organization, such that the customer can access the project and assign authorization to their users.

Given such as setup the owner might want to give one administrative user of the customer organization the role `ORG_OWNER`.
Equipped with this Manager Role, the user can perform actions like configuring their own SSO/Identity Provider, set security policy for their organization, customize branding, or assign project or Manager roles to other users.

An `ORG_OWNER` can also not only delegate Manager roles to other users [as described in the earlier section](#managers-in-isolation) but also manage all aspects of their own organization as well as authorize users to use the granted project.
With ZITADEL there is no need to replicate all settings and projects across organizations.
Instead you set-up the project in one organization, delegate it to different organizations, and then appoint users as Managers of that organization to allow for self-service in a multi-tenancy scenario.

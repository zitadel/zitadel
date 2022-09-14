---
title: Self-Service
---

ZITADEL allows users to perform certain tasks themselves.
For these tasks we either provide an user interface, or the tasks can be initiated or completed through our APIs.

It is important to understand that, depending on your use case, there will exist different user-types that want to perform different actions:  

- `User` are the end-users of your application. Users should be able to perform tasks like register/join, update their profile, manage authenticators etc.There are certain actions that can be executed pre-login, yet others require the user to have a valid session.
- `Manager` are users with a [special manager role within ZITADEL](http://localhost:3000/docs/concepts/structure/managers) and can perform administrative actions such as system configuration or granting access rights to users.

:::info
It is important to note that a `Manager` is not simply an administrative user, but can be used to create much more advanced self-service scenarios.
For example you can create an organization and assign a user from that organization the Manager Role `ORG_OWNER`.
Given this role the user could perform actions like configuring their own SSO/Identity Provider, set security policy for their organization, or assign roles to other users.
:::

All self-service interfaces are available in different [languages](http://localhost:3000/docs/guides/manage/customize/texts#internationalization).

## Registration

:::info
You can pre-select a given organization by passing the scope `urn:zitadel:iam:org:domain:primary:{domainname}` in the authorization request.
This will force users to register only with the specified organization.
Also the branding and login settings (e.g. Social Login Providers) are directly shown to the user.
:::

### Local account

Allows unauthenticated and users, who are not yet in the given organization, to create an account (register) themselves.

- Mandatory profile fields
- Set password
- Accept terms of service and privacy policy
- User receives an email with an one-time code
- User has to enter the one-time code to finish registration
- User can re-quest a new one-time code

The user is prompted on first login to setup MFA.
This step is mandatory in case MFA is enforced in the login policy.

### Existing Identity / SSO / Social Login

Allows unauthenticated and users, who are not yet in the given organization, to register with an external identity provider.
An external identity provider can be a Social Login Provider or a pre-configured identity provider.

- Information from the external identity provider is used to pre-fill the profile information
- User can update the profile information
- An account is created within ZITADEL and linked with the external identity provider

## Login

### Web, Mobile, and Single-Page Applications

[This guide](/docs/guides/integrate/login-users) explains in more detail the login-flows for different application types.
Human users are redirected to ZITADEL's login page and complete sign-in with the interactive login flow.
It is important to understand that ZITADEL provides a hosted login page and the device of the users opens this login page in a browser, even on Native/Mobile apps.

#### MFA / 2FA

Users are automatically prompted to provide a second factor, when

- Instance or organization [login policy](https://docs.zitadel.com/docs/concepts/structure/policies#login-policy)
- Requested by the client
- A multi-factor is setup for the user

When a multi-factor is required, but not set-up, then the user is requested to set-up an additional factor.

#### FIDO Passkeys

Users can select a button to initiate passwordless login or use a fall-back method (ie. login with username/password), if available.

The passwordless login flow follows the FIDO 2 / WebAuthN standard.
Briefly explained the following happens:

- User selects button
- User's device will ask the user to provide a gesture (eg, FaceID, Windows Hello, Fingerprint, PIN)
- The user is being redirected to the application

With the introduction of passkeys the gesture can be provided on ANY of the user's devices.
This is not strictly the device where the login flow is being executed (eg, push notification on a mobile device).
The user experience depends mainly on the used operating system and browser.

### SSO / Social Logins

Given an external identity provider is configured on the instance or on the organization, then: 

- the user will be shown a button for each identity provider as alternative to login with a [local account](#local-account)
- when clicking the button the user will be redirected to the identity provider
- after successful login the user will be redirected to the application

### Machines

Machine accounts can't use an interactive login but require other means of authentication, such as privately-signed JWT or personal access tokens.
Read more about [Service Users](/docs/guides/integrate/serviceusers) and recommended [OpenID Connect Flows](/docs/guides/integrate/oauth-recommended-flows#different-client-profiles).

### Other Clients

We currently do not expose the Login APIs.
Whereas you can register users via the management API, you can't login users with our APIs.
This might be important in cases where you can't use a website (eg, Games, VR, ...).

### Customization

The login page can be changed by customizing different branding aspects and you can define a custom domain for the login (eg, login.acme.com).

:::info
By default, the displayed branding is defined based on the user's domain. In case you want to show the branding of a specific organization by default, you need to either pass a primary domain scope (`urn:zitadel:iam:org:domain:primary:{domainname}`) with the authorization request, or define the behavior on your Project's settings.
:::

### Account picker

A list of accounts that were used to log-in are shown to the user.
The user can click the account in the list and does not need to type the username.

:::info
This behavior can be changed with the authorization request. Please refer to our [guide](/docs/guides/integrate/login-users).
:::

## Logout

Users can terminate all their sessions (logout).
A client will implement this by calling the [specific endpoint](http://localhost:3000/docs/apis/openidoauth/endpoints#end_session_endpoint).

## Secrets

### Password reset

Unauthenticated users can request a password reset after providing the loginname during the login flow.

- User selects reset password
- An email will be sent to the verified email address
- User has 

### Change password

### MFA / FIDO Passkeys

- Add & remove second factor
- Add & remove passwordless authenticator

## Profile

Authenticated users.

### Update Information

### Email Verification

### Phone Verification

### Account Linking

- Add external identity providers

## Managers

[Roles](/docs/concepts/structure/managers#roles)

Can be human users or also service users (eg, to manage programmatically)
---
title: ZITADEL Policies
---

Policies are configurations of all the different parts of the IAM. For all parts we have a suitable default in the IAM.
The default configuration can be overridden for each organization.

## General

You can find this settings in the menu organisation in the section polcies.
Each policy can be overriden and resetted to the default.

## Password Complexity

With the password complexity policy you can define how a password should look.

The following properties can be set: 
- Minimum Length
- Has Uppercase
- Has Lowercase
- Has Number
- Has Symbol

![Password Complexity Policy](/img/manuals/policies/console_org_pw_complexity.png)

## Login Policy

The Login Policy defines how the login process should look like and which authentication options a user has to authenticate.

| Setting | Description |
| --- | --- |
| Register allowed | Enable self register possibility in the login ui |
| Username Password allowed | Possibility to login with username and password |
| External IDP allowed | Possibility to login with an external identity (e.g Google, Microsoft, Apple, etc)|
| Force MFA | Force a user to register and use a multifactor authentication |
| Passwordless | Choose if passwordless login is allowed or not |

![Login Policy](/img/manuals/policies/console_org_login.png)

### Multifactors / Second Factors

In the multifactors section you can configure what kind of multifactors should be allowed.
Multifactors: 
- U2F (Universal Second Factor) with PIN

Secondfacrots: 
- OTP (One Time Password)
- U2F (Universal Second Factor)

![Second- and Multifactors](/img/manuals/policies/console_org_second_and_multi_factors.png)

### Identity Providers

You can configure all kinds of external identity providers for identity brokering, which support OIDC (OpenIDConnect).
Make a new identity provider configuration and enable it in the list after.

For a detailed guide about how to configure a new identity provider for identity brokering have a look at our guide:
[Identity Brokering](../guides/usage/identity-brokering)
---
title: ZITADEL Policies
---

Policies are configurations of all the different parts of the IAM. For all parts we have a suitable default in the IAM.
The default configuration can be overridden for each organization.

## General

You can find these settings in the menu organisation in the section polcies.
Each policy can be overriden and resetted to the default.

## Password Complexity

With the password complexity policy you can define the requirements for a users password.

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

Secondfactors: 
- OTP (One Time Password)
- U2F (Universal Second Factor)

![Second- and Multifactors](/img/manuals/policies/console_org_second_and_multi_factors.png)

### Identity Providers

You can configure all kinds of external identity providers for identity brokering, which support OIDC (OpenID Connect).
Create a new identity provider configuration and enable it in the list afterwards.

For a detailed guide about how to configure a new identity provider for identity brokering have a look at our guide:
[Identity Brokering](../guides/usage/identity-brokering)

## Private Labeling / Branding

With private labeling you can brand and customize your login page and emails, that it matches your CI/CD.
You can configure a light and a dark design.

Make sure you click the "Set preview as current configuration" button after you finish your configuration. Before this it will only be set as your preview configuration.

| Setting | Description |
| --- | --- |
| Logo | Upload your logo for the light and the dark design. |
| Colors | You can set four different colors to design your login page and email. (Background-, Primary-, Warn- and Font Color) |
| Font | Upload your custom font |
| Hide Loginname suffix | If enabled,  your loginname suffix (Domain) will not be shown in the login page |
| Disable Watermark | If you disable the watermark you will not see the "Powered by ZITADEL" in the login page |

![Private Labeling](/img/manuals/policies/console_org_private_labeling.png)

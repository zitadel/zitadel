---
title: Settings/Policies
---

Settings and policies are configurations of all the different parts of the Instance or an organization. For all parts we have a suitable default in the Instance.
The default configuration can be overridden for each organization, some policies are currently only available on the instance level. If thats the case it will be mentioned on the descriptions below.

You can find these settings in the instance page under settings, or on a specific organization menu organization in the section policies.
Each policy can be overridden and reset to the default.

## General

:::info
Only available on the instance settings
:::

At the moment general settings is only one configuration. This defines the default language of the whole instance.

![General Settings](/img/console_instance_policy_general.png)

## Notification

:::info
Only available on the instance settings
:::

In the notification settings you can configure your SMTP and an SMS Provider. At the moment only Twilio is available as SMS provider.

### SMTP 
On each instance we configure our default SMTP provider. To make sure, that you only send some E-Mails from domains you own. You need to add a custom domain on your instance.
Go to the ZITADEL [customer portal](https://zitadel.cloud) to configure a custom domain. 

![Notification Providers](/img/console_instance_policy_notification.png)

### SMS

No default provider is configured to send some sms to your users. If you like to validate the phone numbers of your users make sure to add your twilio configuration.

![Notification Providers](/img/console_instance_policy_notification_twilio.png)

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

### Passwordless

Passwordless authentication means that the user doesn't need to enter a password to login. In our case the user has to enter his loginname and as the next step proof the identity through a registered device or token.
There are two different types one is depending on the device (e.g. Fingerprint, Face recognition, WindowsHello) and the other is independent (eg. Yubikey, Solokey).

### Multifactor

In the multifactors section you can configure what kind of multifactors should be allowed. For passwordless to work, it's required to enable U2F (Universial Second Factor) with PIN. There is no other option at the moment.
Multifactors:
- U2F (Universal Second Factor) with PIN

Secondfactors:
- OTP (One Time Password)
- U2F (Universal Second Factor)

![Second- and Multifactors](/img/manuals/policies/console_org_second_and_multi_factors.png)

## Password Complexity

With the password complexity policy you can define the requirements for a users password.

The following properties can be set:
- Minimum Length
- Has Uppercase
- Has Lowercase
- Has Number
- Has Symbol

![Password Complexity Policy](/img/manuals/policies/console_org_pw_complexity.png)

## Lockout Policy

Define when an account should be locked.

The following settings are available:
- Maximum Password Attempts: When the user has reached the maximum password attempts the account will be locked

If an account is locked, the administrator has to unlock it in the ZITADEL console

## Identity Providers

You can configure all kinds of external identity providers for identity brokering, which support OIDC (OpenID Connect).
Create a new identity provider configuration and enable it in the list afterwards.

For a detailed guide about how to configure a new identity provider for identity brokering have a look at our guide:
[Identity Brokering](../../guides/integrate/identity-brokering)

## Domain policy

In the domain policy you have two different settings.
One is the "user_login_must_be_domain", by setting this all the users within an organisation will be suffixed with the domain of the organisation.

The second is "validate_org_domains" if this is set to true all created domains on an organisation must be verified per acme challenge.
More about how to verify a domain [here](../../guides/manage/console/organizations#domain-verification-and-primary-domain).
If it is set to false, all registered domain will automatically be created as verified and the users will be able to use the domain for login.

## Branding

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

![Private Labeling](/img/console_private_labeling.png)

## Privacy Policy and TOS

Each organization is able to configure its own privacy policy, terms of service and help.
A link to the current policies can be provided. On register each user has to accept these policies.

By clicking on an input field you can see the language attribute to integrate into a link, for the possibility to have different links for different languages.
The language of the user will be set into the url.
Example:
https://demo.com/tos-{{.Lang}}

![Privacy Policy](/img/manuals/policies/console_org_privacy.png)

## OIDC token lifetime and expiration

:::info
Only available on the instance settings
:::

Configure how long the different oidc tokens should life.
You can set the following times:
- Access Token Lifetime
- ID Token Lifetime
- Refresh Token Expiration
- Refresh Token Idle Expiration

![OIDC Token Lifetimes](/img/manuals/policies/console_policy_oidc.png)


## Secret appearance

:::info
Only available on the instance settings
:::

ZITADEL has some different codes and secrets, that can be specified.
You can configure what kind of characters should be included, how long the secret should be and the expiration.
The following secrets can be configured:
- Initialization Mail Code
- Email verification code
- Phone verification code
- Password reset code
- Passwordless initialization code
- Application secrets

![OIDC Token Lifetimes](/img/manuals/policies/console_policy_secrets.png)
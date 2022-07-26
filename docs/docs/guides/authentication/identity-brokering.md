---
title: Identity Brokering
---

<table class="table-wrapper">
    <tr>
        <td>Description</td>
        <td>Learn about identity brokering/federation and how to add an external identity provider to authenticate your users.</td>
    </tr>
    <tr>
        <td>Learning Outcomes</td>
        <td>
            In this module you will:
            <ul>
                <li>Learn about Identity Providers</li>
                <li>Add a new identity provider</li>
                <li>Example with Google Login</li>
            </ul>
        </td>
    </tr>
     <tr>
        <td>Prerequisites</td>
        <td>
            <ul>
                <li>Knowledge of <a href="/docs/guides/start/organizations">Organizations</a></li>
            </ul>
        </td>
    </tr>
</table>

## What is Identity Brokering and Federated Identities?

Federated identity management is an arrangement built upon the trust between two or more domains. Users of these domains are allowed to access applications and services using the same identity.
This identity is known as federated identity and the pattern behind this as identity federation.

A service provider that specializes in brokering access control between multiple service providers (also referred to as relying parties) is called identity broker.
Federated identity management is an arrangement that is made between two or more such identity brokers across organizations.

Example:
If Google is configured as identity provider on your organization, the user will get the option to use his Google Account on the Login Screen of ZITADEL (1).
ZITADEL will redirect the user to the login screen of Google where he as to authenticated himself (2) and is sent back after he has finished that (3).
Because Google is registered as trusted identity provider the user will be able to login in with the Google account after he linked an existing ZITADEL Account or just registered a new one with the claims provided by Google (4)(5).

![Identity Brokering](/img/guides/identity_brokering.png)

## Exercise: Register an external identity provider

In this exercise we will add a new Google identity provider to federate identities with ZITADEL.

### 1. Create new OIDC Client

1. Register an OIDC Client in your preferred provider
2. Make sure you add the ZITADEL callback redirect uris
   - {your-domain}/ui/login/register/externalidp/callback
   - {your-domain}/ui/login/login/externalidp/callback

> **Information:** Make sure the provider is OIDC 1.0 compliant with a proper Discovery Endpoint

Google Example:

1. Go to the Google Gloud Platform and choose youre project: <https://console.cloud.google.com/apis/credentials>
2. Click on "+ CREATE CREDENTIALS" and choose "OAuth client ID"
3. Choose Web application as Application type and give a name
4. Add the redirect uris
   - {your-domain}/ui/login/register/externalidp/callback
   - {your-domain}/ui/login/login/externalidp/callback
5. Save clientid and client secret

![Add new oAuth credentials in Google Console](/img/google_add_credentials.gif)

### 2. Add custom login policy

The login policy can be configured on two levels. Once as default on the instance and this can be overwritten for each organization.
This case describes how to change it on the organization.

1. Go to your organization settings by clicking on "Organization" in the menu
2. Modify your login policy in the menu "Login Behaviour and Security"

![Add custom login policy](/img/console_org_custom_login_policy.gif)

### 3. Configure new identity provider

1. Go to the settings of your instance or a specific organization (depending on where you need the identity provider)
2. Go to the identity providers section and click "New"
3. Select "OIDC Configuration" and fill out the form
   - Use the issuer, clientid and client secret provided by your provider (Google Issuer: https://accounts.google.com)
   - The scopes will be prefilled with openid, profile and email, because this information is relevant for ZITADEL
   - You can choose what fields you like to map as the display name and as username. The fields you can choose are preferred_username and email
     (Example: For Google you should choose email for both fields)
4. Save your configuration
5. You will now see the created configuration in the list. Click on the activate icon at the end of the row you can see when hovering over the row, to activate it in the login flow.

![Configure identity provider](/img/console_org_identity_provider.gif)

### 4. Send the primary domain scope on the authorization request
ZITADEL will show a set of identity providers by default. This configuration can be changed by users with the [manager role] (https://docs.zitadel.com/docs/concepts/zitadel/objects/managers) `IAM_OWNER`.

An organization's login settings will be shown 

- as soon as the user has entered the loginname and ZITADEL can identitfy to which organization he belongs; or
- by sending a primary domain scope.
To get your own configuration you will have to send the [primary domain scope](https://docs.zitadel.com/docs/apis/openidoauth/scopes#reserved-scopes) in your [authorization request](https://docs.zitadel.com/docs/guides/authentication/login-users/#auth-request) .
The primary domain scope will restrict the login to your organization, so only users of your own organization will be able to login, also your branding and policies will trigger.

:::note

You need to create your own auth request with your applications parameters. Please see the docs to construct an [Auth Request](https://docs.zitadel.com/docs/guides/authentication/login-users/#auth-request).

:::

Your user will now be able to choose Google for login instead of username/password or mfa.

## Knowledge Check

* The issuer for your identity provider is <https://issuer.zitadel.ch>
    - [ ] yes
    - [ ] no
* The identity provider has to be oAuth 2.0 compliant
    - [ ] yes
    - [ ] no

<details>
    <summary>
        Solutions
    </summary>

* The issuer for your identity provider is https://issuer.zitadel.ch
    - [ ] yes
    - [x] no (The issuer is provided by your choosen identity provider. In the case of Google it's https://accounts.google.com)
* The identity provider has to be oAuth 2.0 compliant
    - [x] yes
    - [ ] no

</details>

## Summary

* You can federate identities of all oAuth 2.0 compliant external identity providers
* Configure the provider in your custom login policy

Where to go from here:

* ZITADEL Projects
* Service users

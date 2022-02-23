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
                <li>See an example with Google Login</li>
            </ul>
        </td>
    </tr>
     <tr>
        <td>Prerequisites</td>
        <td>
            <ul>
                <li>Knowledge of <a href="/docs/guides/basics/organizations">Organizations</a></li>
            </ul>
        </td>
    </tr>
</table>

## What is Identity Brokering and Federated Identities?

_Federated identity management_ is an arrangement built upon the trust between two or more domains.
Users of these domains can access applications and services using the same identity.
This identity is known as _federated identity_.
The pattern behind it is known as _identity federation_.

An _identity broker_ is a service provider that specializes in brokering access control between multiple service providers (also referred to as relying parties).
_Federated identity management_ is an arrangement that is made between two or more such identity brokers across organizations.

For example,
if Google is configured as an identity provider on your organization, users can use their Google Account on the Login Screen of ZITADEL (1).
ZITADEL redirects users to the login screen of Google, where they authenticate themselves, (2) and are sent back (3).
Because Google is registered as a trusted identity provider, users can log in with the Google account after they link an existing ZITADEL Account, or after they register a new one with the claims provided by Google (4)(5).

![Identity Brokering](/img/guides/identity_brokering.png)

## Exercise: Register an external identity provider

In this exercise, we will add a new Google identity provider to federate identities with ZITADEL.

### 1. Create new OIDC Client

1. Register an OIDC Client in your preferred provider.
2. Make sure you add the ZITADEL callback redirect URIs.
   https://accounts.zitadel.ch/register/externalidp/callback
   https://accounts.zitadel.ch/login/externalidp/callback

> **Information:** Make sure the provider is OIDC 1.0 compliant with a proper Discovery Endpoint

Google Example:

1. Go to the Google Gloud Platform and choose youre project: <https://console.cloud.google.com/apis/credentials>
2. Select **+ CREATE CREDENTIALS**. Choose **OAuth client ID**.
3. Choose **Web application** as the Application type and give it a name.
4. Add the redirect URIs from above.
5. Save the `clientid` and client secret.

![Add new oAuth credentials in Google Console](/img/google_add_credentials.gif)

### 2. Add custom login policy on your organization

1. Go to your organization settings by selecting **Organization** in the menu, or by following this link: <https://console.zitadel.ch/org>.
2. Modify your login policy.
3. As long as you have the default policy, you can't change the policy. To set your own settings, select **create custom**.

![Add custom login policy](/img/console_org_custom_login_policy.gif)

### 3.Configure new identity provider

1. Go to the identity providers section. Select **new**.
2. Fill out the form:
   - Use the issuer, `clientid` and client secret provided by your provider
   - The scopes will be prefilled with `openid`, `profile` and `email`, because this information is relevant for ZITADEL.
   - Choose what fields you like to map as the display name and username. The fields you can choose are `preferred_username` and `email`
     (Example: For Google, you should choose `email` for both fields)
3. Save your configuration
4. Link your new configuration to your login policy.
   Search the organization category to get your own configuration.
   If you choose system you can link all predefined providers.

![Configure identity provider](/img/console_org_identity_provider.gif)

### 4.Send the primary domain scope on the authorization request
ZITADEL shows a set of identity providers by default. This configuration can be changed by users with the [manager role] (https://docs.zitadel.ch/docs/concepts/zitadel/objects/managers) `IAM_OWNER`.

An organization's login settings are shown:

- as soon as users enter the `loginname`, and ZITADEL can identify which organizations they belong to; or
- by sending a primary domain scope.
To get your own configuration, send the [primary domain scope](https://docs.zitadel.ch/docs/apis/openidoauth/scopes#reserved-scopes) in your [authorization request](https://docs.zitadel.ch/docs/guides/authentication/login-users/#auth-request) .
The primary domain scope restricts the login to your organization, so only users of your own organization can login.
This also triggers your branding and policies.

See the following link as an example. Users can register and login to the organization that verified the `@caos.ch` domain only.
```
https://accounts.zitadel.ch/oauth/v2/authorize?client_id=69234247558357051%40zitadel&scope=openid%20profile%20urn%3Azitadel%3Aiam%3Aorg%3Adomain%3Aprimary%3Acaos.ch&redirect_uri=https%3A%2F%2Fconsole.zitadel.ch%2Fauth%2Fcallback&state=testd&response_type=code&nonce=test&code_challenge=UY30LKMy4bZFwF7Oyk6BpJemzVblLRf0qmFT8rskUW0
```

:::info

Make sure to replace the domain `caos.ch` with your own domain to trigger the correct branding.

:::

:::caution

This example uses the ZITADEL Cloud Application for demonstration.
You need to create your own auth request with your applications parameters.
Please see the docs to construct an [Auth Request](https://docs.zitadel.ch/docs/guides/authentication/login-users/#auth-request).

:::

Your user now can choose Google to login instead of username/password or mfa.

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
    - [x] no (The issuer is provided by your chosen identity provider. In the case of Google, it's https://accounts.google.com)
* The identity provider has to be oAuth 2.0 compliant
    - [x] yes
    - [ ] no

</details>

## Summary

* You can federate identities of all oAuth 2.0 compliant external identity providers.
* Configure the provider in your custom login policy.

Where to go from here:

* ZITADEL Projects
* Service users

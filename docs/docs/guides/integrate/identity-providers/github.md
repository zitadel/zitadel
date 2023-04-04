---
title: Configure GitHub as Identity Provider
sidebar_label: GitHub
---

This guides shows you how to connect GitHub or GitHub Enterprise as an identity provider in ZITADEL.

:::info
In ZITADEL you can connect an Identity Provider (IdP) like GitHub to your instance and provide it as default to all organizations or you can register the IdP to a specific organization only. This can also be done through your customers in a self-service fashion.
:::

## Prerequisite

Make sure to read and follow the [General Guide](./general) on how to setup identity providers for your ZITADEL instance first, before you go through the specific guide here.

## GitHub Configuration

### Register a new application

For **GitHub** browse to the [Register a new OAuth application](https://github.com/settings/applications/new). You can find this link withing [Settings](https://github.com/settings/profile) - [Developer Settings](https://github.com/settings/apps) - - [OAuth Apps](https://github.com/settings/developers).

For **GitHub Enterprise** go to your GitHub Enterprise home page and then to Settings - Developer Settings - OAuth Apps - Register a new application/New OAuth App

Fill in the application name and homepage URL.

You have to add the authorization callback URL, where GitHub should redirect, after the user has authenticated himself.
In this example our test instance has the domain `https://acme-gzoe4x.zitadel.cloud`.
This results in the following authorization callback URL:
 `https://acme-gzoe4x.zitadel.cloud/ui/login/login/externalidp/callback`

:::info
To adapt this for you setup just replace the domain
:::

![Register an OAuth application](/img/guides/github_oauth_app_registration.png)

### Client ID and Secret

After clicking "Register application" , you will see the detail page of the application you have just created.
To be able to connect GitHub to ZITADEL you will need a client ID and a client secret. 
The client ID you can copy directly. A secret you have to generate by clicking "Generate new client secret".
Make sure to save the secret, as you will not be able to show it again.

![Client ID and Secret](/img/guides/github_oauth_client_id_secret.png)

## ZITADEL Configuration

### Create new GitHub Provider

Go to the settings of your ZITADEL instance or the organization where you like to add a new GitHub provider.
Choose the GitHub provider template. This template has everything you need preconfigured. You only have to add the client ID and secret, you have created in the step before.

You can configure the following settings if you like, a useful default will be filled if you don't change anything:

**Scopes**: The scopes define which scopes will be sent to the provider, `openid`, `profile`, and `email` are prefilled. This informations will be taken to create/update the user within ZITADEL.

**Automatic creation**: If this setting is enabled the user will be created automatically within ZITADEL, if it doesn't exist.

**Automatic update**: If this setting is enabled, the user will be updated within ZITADEL, if some user data are changed withing the provider. E.g if the lastname changes on the GitHub account, the information will be changed on the ZITADEL account on the next login. 

**Account creation allowed**: This setting determines if account creation within ZITADEL is allowed or not.

**Account linking allowed**: This setting determines if account linking is allowed. (E.g an account within ZITADEL should already be existing and the when login with GitHub an account should be linked)

:::info
Either account creation or account linking have to be enabled. Otherwise, the provider can't be used.
:::

![GitHub Provider](/img/guides/zitadel_github_create_provider.png)

### Activate IdP

Once you created the IdP you need to activate it, to make it usable for your users.

![Activate the GitHub](/img/guides/zitadel_activate_github.png)

## Test the setup

To test the setup use incognito mode and browse to your login page.
If you succeeded you should see a new button which should redirect you to your GitHub.

Per default the login of your instance will be shown, read the following section on how to trigger it for a specific organization: [Organization Scope](./general#trigger-configuration-on-the-login-for-a-specific-organization)


![GitHub Button](/img/guides/zitadel_login_github.png)

![GitHub Login](/img/guides/github_login.png)

If the user is not yet linked in ZITADEL the user will see the screen below.
Because GitHub is an OAuth provider and oAuth does not provide a standardized way to get the user data not all of the data can be taken over. First and Lastname are not filled.
The user has to enter the rest of the data himself.

![GitHub Login](/img/guides/zitadel_login_external_not_found_registration.png)

### Optional: Add ZITADEL action to autofill userdata

If you don't want the user to have to enter his first and lastname himself, you can add a ZITADEL action in which you specify how the data should be transferred.

1. Go to the settings of the organization where the users will be registered
2. Add an new action with the following body. Make sure the action name is the same as in the script itself. Make sure to change the id in the script to the id of your own identity provider configuration. 

```js reference
https://github.com/zitadel/actions/blob/main/examples/github_identity_provider
```


3. Add the action to the flow "External Authentication" on the trigger Post Authentication
  

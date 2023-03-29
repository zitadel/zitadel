---
title: Configure GitLab as Identity Provider
sidebar_label: GitLab
---

This guides shows you how to connect GitLab or GitLab SelfHosted as an identity provider in ZITADEL.

:::info
In ZITADEL you can connect an Identity Provider (IdP) like GitLab to your instance and provide it as default to all organizations or you can register the IdP to a specific organization only. This can also be done through your customers in a self-service fashion.
:::

## Prerequisite

To be able to use GitLab to authenticate your users you need to enable OpenID Connect for OAuth applications. You can find more informations in the following link: [GitLab as OpenID Connect identity provider](https://docs.gitlab.com/ee/integration/openid_connect_provider.html)

## GitLab Configuration

### Register a new application

1. Login to [gitlab.com](https://gitlab.com)
2. Select [Edit Profile](https://gitlab.com/-/profile)
3. Click on [Applications](https://gitlab.com/-/profile/applications) in the side navigation

For **GitLab Self-Hosted** go to your GitLab Selfhosted instanceand follow the same steps as for GitLab.

Fill in the application name.

You have to add the redirect URI, where GitLab should redirect, after the user has authenticated himself.
In this example our test instance has the domain `https://acme-gzoe4x.zitadel.cloud`.
This results in the following redirect URI:
 `https://acme-gzoe4x.zitadel.cloud/ui/login/login/externalidp/callback`

:::info
To adapt this for you setup just replace the domain
:::

![Register an OAuth application](/img/guides/gitlab_app_registration.png)

### Client ID and Secret

After clicking "Save application", you will see the detail page of the application you have just created.
To be able to connect GitLab to ZITADEL you will need a client ID and a client secret. 
Save the ID and the Secret, you will not be able to copy the secret again, if you lose it you have to generate a new one.

![Client ID and Secret](/img/guides/gitlab_app_id_secret.png)

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

![GitHub Provider](/img/guides/zitadel_gitlab_create_provider.png)

### Activate IdP

Once you created the IdP you need to activate it, to make it usable for your users.

![Activate the GitHub](/img/guides/zitadel_activate_gitlab.png)

## Test the setup

To test the setup use incognito mode and browse to your login page.
If you succeeded you should see a new button which should redirect you to your GitHub.

![GitHub Button](/img/guides/zitadel_login_gitlab.png)

![GitHub Login](/img/guides/gitlab_login.png)

If the user is not yet linked in ZITADEL the user will see the screen below.
Because GitHub is an OAuth provider and oAuth does not provide a standardized way to get the user data not all of the data can be taken over. First and Lastname are not filled.
The user has to enter the rest of the data himself.

![GitHub Login](/img/guides/zitadel_login_external_not_found_registration.png)

### Optional: Add ZITADEL action to autofill userdata

If you don't want the user to have to enter his first and lastname himself, you can add a ZITADEL action in which you specify how the data should be transferred.

1. Go to the settings of the organization where the users will be registered

2. Add an new action with the following body. Make sure the action name is the same as in the script itself. Make sure to change the id in the script to the id of your own identity provider configuration. 

```js reference

https://github.com/zitadel/actions/blob/main/examples/gitlab_identity_provider

```
3. Add the action to the flow "External Authentication" on the trigger Post Authentication
  

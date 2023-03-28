---
title: Configure Google as Identity Provider
sidebar_label: Google
---

This guides shows you how to connect Google as an identity provider in ZITADEL.

:::info
In ZITADEL you can connect an Identity Provider (IdP) like Google to your instance and provide it as default to all organizations or you can register the IdP to a specific organization only. This can also be done through your customers in a self-service fashion.
:::

## Prerequisite

To be able to use Google to authenticate your users you need an OIDC Client in Google.

## Google Configuration

### Register a new client

1. Go to the Google Cloud Platform and choose your project: <https://console.cloud.google.com/apis/credentials>
2. Click on "+ CREATE CREDENTIALS" and choose "OAuth client ID"
3. Choose "Web application" as application type and give a name
4. Add the redirect uris
 - {your-domain}/ui/login/login/externalidp/callback
 - Example redirect url for the domain `https://acme-gzoe4x.zitadel.cloud` would look like this:  `https://acme-gzoe4x.zitadel.cloud/ui/login/login/externalidp/callback`
5. Save the Client ID and Client secret

![Google OAuth App Registration](/img/guides/google_oauth_app_registration.png)

![Google Client ID and Secret](/img/guides/google_client_id_secret.png)

## ZITADEL Configuration

### Create new GitHub Provider

Go to the settings of your ZITADEL instance or the organization where you like to add a new Google provider.
Choose the Google provider template. This template has everything you need preconfigured. You only have to add the client ID and secret, you have created in the step before.

You can configure the following settings if you like, a useful default will be filled if you don't change anything:

**Scopes**: The scopes define which scopes will be sent to the provider, `openid`, `profile`, and `email` are prefilled. This information will be taken to create/update the user within ZITADEL.

**Automatic creation**: If this setting is enabled the user will be created automatically within ZITADEL, if it doesn't exist.

**Automatic update**: If this setting is enabled, the user will be updated within ZITADEL, if some user data are changed withing the provider. E.g if the lastname changes on the GitHub account, the information will be changed on the ZITADEL account on the next login. 

**Account creation allowed**: This setting determines if account creation within ZITADEL is allowed or not.

**Account linking allowed**: This setting determines if account linking is allowed. (E.g an account within ZITADEL should already be existing and the when login with GitHub an account should be linked)

:::info
Either account creation or account linking have to be enabled. Otherwise, the provider can't be used.
:::

![GitHub Provider](/img/guides/zitadel_google_create_provider.png)

### Activate IdP

Once you created the IdP you need to activate it, to make it usable for your users.

![Activate the GitHub](/img/guides/zitadel_activate_google.png)

## Test the setup

To test the setup use incognito mode and browse to your login page.
If you succeeded you should see a new button which should redirect you to your GitHub.

![GitHub Button](/img/guides/zitadel_login_github.png)

![GitHub Login](/img/guides/google_login.png)

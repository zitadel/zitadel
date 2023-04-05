---
title: Configure Azure AD as Identity Provider
sidebar_label: Azure AD
---

This guides shows you how to connect Azure AD as an identity provider in ZITADEL.

:::info
In ZITADEL you can connect an Identity Provider (IdP) like Azure AD to your instance and provide it as default to all organizations or you can register the IdP to a specific organization only. This can also be done through your customers in a self-service fashion.
:::

## Prerequisite

Make sure to read and follow the [General Guide](./general) on how to setup identity providers for your ZITADEL instance first, before you go through the specific guide here.

## Azure AD Configuration

You need to have access to an AzureAD Tenant. If you do not yet have one follow [this guide from Microsoft](https://docs.microsoft.com/en-us/azure/active-directory/develop/quickstart-create-new-tenant) to create one for free.

### Register a new client

1. Browse to the [App registration menus create dialog](https://portal.azure.com/#view/Microsoft_AAD_RegisteredApps/CreateApplicationBlade/quickStartType~/null/isMSAApp~/false) to create a new app.
2. Give the application a name and choose who should be able to login (Single-Tenant, Multi-Tenant, Personal Accounts, etc.) This setting will also have an impact on how to configure the provider later on in ZITADEL.
3. Choose "Web" in the redirect uri field and add the URL:
 - {your-domain}/ui/login/login/externalidp/callback
 - Example redirect url for the domain `https://acme-gzoe4x.zitadel.cloud` would look like this: `https://acme-gzoe4x.zitadel.cloud/ui/login/login/externalidp/callback`
5. Save the Application (client) ID and the Directory (tenant) ID from the detail page

![Azure App Registration](/img/guides/azure_app_registration.png)

![Azure Client ID and Tenant ID](/img/guides/azure_client_tenant_id.png)

### Add client secret

To be able to authenticate your Azure application you have to generate a new client secret.

1. Click on client credentials on the detail page of the application or use the menu "Certificates & secrets"
2. Click on "+ New client secret" and enter a description and an expiry date, add the secret afterwards
3. Copy the value of the secret. You will not be able to see the value again after some time 

![Azure Client Secret](/img/guides/azure_client_secret.png)

### Token configuration

To allow ZITADEL to get the information from the authenticating user you have to configure what kind of optional claims should be returned in the token.

1. Click on Token configuration in the side menu
2. Click on "+ Add optional claim"
3. Add email, family_name, given_name and preferred_username to the id token

![Azure Token configuration](/img/guides/azure_token_configuration.png)

### API permissions

To be able to get all the information that ZITADEL needs, you have to configure the correct permissions.

1. Go to "API permissions" in the side menu
2. Make sure the permissions include "Microsoft Graph": email, profile and User.Read
3. The "Other permissions granted" should include "Microsoft Graph: openid"

![Azure API permissions](/img/guides/azure_api_permissions.png)

## ZITADEL Configuration

### Create new Azure AD Provider

Go to the settings of your ZITADEL instance or the organization where you like to add a new **Azure AD** provider.
Choose the **Microsoft** provider template.
This template has everything you need preconfigured.
You only have to add the client ID and secret, you have created in the step before.

You can configure the following settings if you like, a useful default will be filled if you don't change anything:

**Scopes**: The scopes define which scopes will be sent to the provider, `openid`, `profile`, and `email` are prefilled.
This information will be taken to create/update the user within ZITADEL. Make sure to also add `User.Read`

**Email Verified**: Azure AD doesn't send the email verified claim in the users token, if you don't enable this setting.
The user is then created with an unverified email, which results in an email verification message.
If you want to avoid that, make sure to enable "Email verified".
In that case, the user is created with a verified email address.

**Tenant Type**: Configure the tenant type according to what you have chosen in the settings of your Azure AD application previously.
- Common: Choose common if you want all Microsoft accounts being able to login.
In this case, configure "Accounts in any organizational directory and personal Microsoft accounts" in your Azure AD App.
- Organizations: Choose organization if you have Azure AD Tenants and no personal accounts. (You have configured either "Accounts in this organization" or "Accounts in any organizational directory" on your Azure APP)
- Consumers: Choose this if you want to allow public accounts. (In your Azure AD App you have configured "Personal Microsoft accounts only")

**Tenant ID**: If you have selected either the *Organizations* or *Customers* as the *Tenant Type*, you have to enter the *Directory (Tenant) ID*, copied previously in the Azure App configuration, here.

**Automatic creation**: If this setting is enabled the user will be created automatically within ZITADEL, if it doesn't exist.

**Automatic update**: If this setting is enabled, the user will be updated within ZITADEL, if some user data are changed withing the provider. E.g if the lastname changes on the GitHub account, the information will be changed on the ZITADEL account on the next login. 

**Account creation allowed**: This setting determines if account creation within ZITADEL is allowed or not.

**Account linking allowed**: This setting determines if account linking is allowed. (E.g an account within ZITADEL should already be existing and the when login with GitHub an account should be linked)

:::info
Either account creation or account linking have to be enabled. Otherwise, the provider can't be used.
:::

![Azure Provider](/img/guides/zitadel_azure_provider.png)

### Activate IdP

Once you created the IdP you need to activate it.

![Activate Azure AD](/img/guides/zitadel_activate_azure.png)

## Test the setup

To test the setup, use incognito mode and browse to your login page.
If you succeeded you should see a new button which should redirect you to your Microsoft Login.

Per default the login of your instance will be shown, read the following section on how to trigger it for a specific organization: [Organization Scope](./general#trigger-configuration-on-the-login-for-a-specific-organization)

![Azure Button](/img/guides/zitadel_login_azure.png)

![GitHub Login](/img/guides/microsoft_login.png)

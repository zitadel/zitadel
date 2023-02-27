---
title: Configure AzureAD as Identity Provider
sidebar_label: AzureAD
---

## AzureAD Tenant as Identity Provider for ZITADEL

This guides shows you how to connect an AzureAD Tenant to ZITADEL.

:::info
In ZITADEL you can connect an Identity Provider (IdP) like an AzureAD to your instance and provide it as default to all organizations or you can register the IdP to a specific organization only. This can also be done through your customers in a self-service fashion.
:::

### Prerequisite

You need to have access to an AzureAD Tenant. If you do not yet have one follow [this guide from Microsoft](https://docs.microsoft.com/en-us/azure/active-directory/develop/quickstart-create-new-tenant) to create one for free.

### AzureAD Configuration

#### Create a new Application

Browse to the [App registration menus create dialog](https://portal.azure.com/#view/Microsoft_AAD_RegisteredApps/CreateApplicationBlade/quickStartType~/null/isMSAApp~/false) to create a new app.

![Create an Application](/img/guides/azure_app_register.png)

:::info
Make sure to select `web` as application type in the `Redirect URI (optional)` section.
You can leave the second field empty since we will change this in the next step.
:::

![Create an Application](/img/guides/azure_app.png)

#### Configure Redirect URIS

For this to work you need to whitelist the redirect URIs from your ZITADEL Instance.
In this example our test instance has the domain `test-qcon0h.zitadel.cloud`. In this case we need to whitelist these two entries:

- `https://test-qcon0h.zitadel.cloud/ui/login/register/externalidp/callback`
- `https://test-qcon0h.zitadel.cloud/ui/login/login/externalidp/callback`

:::info
To adapt this for you setup just replace the domain
:::

![Configure Redirect URIS](/img/guides/azure_app_redirects.png)

#### Create Client Secret

To allow your ZITADEL to communicate with the AzureAD you need to create a Secret

![Create Client Secret](/img/guides/azure_app_secrets.png)

:::info
Please save this for the later configuration of ZITADEL
:::

#### Configure ID Token Claims

![Configure ID Token Claims](/img/guides/azure_app_token.png)

### ZITADEL Configuration

#### Create IdP

Use the values displayed on the AzureAD Application page in your ZITADEL IdP Settings.

- You need to extract the `issuer` of your AzureAD Tenant from the OpenID configuration (`OpenID Connect metadata document`) in the `Endpoints submenu`. It should be your tenant's domain appended with `/v2.0`
- The `Client ID` of ZITADEL corresponds to the `Application (client) ID` in the Overview page
- The `Client Secret` was generated during the `Create Client Secret` step

![Azure Application](/img/guides/azure_app.png)

![Create IdP](/img/guides/azure_zitadel_settings.png)

#### Activate IdP

Once you created the IdP you need to activate it, to make it usable for your users.

![Activate the AzureAD](/img/guides/azure_zitadel_activate.png)

![Active AzureAD](/img/guides/azure_zitadel_active.png)

#### Disable 2-Factor prompt

If a user has no 2-factor configured, ZITADEL does ask on a regularly basis, if the user likes to add a new 2-factor for more security.
If you don't want your users to get this prompt when using Azure, you have to disable this feature.

1. Go to the login behaviour settings of your instance or organization, depending if you like to disable it for all or just a specific organization respectively
2. Set "Multi-factor init lifetimes" to 0

![img.png](/img/guides/login_lifetimes.png)

#### Create user with verified email

Azure AD does not send the "email verified claim" in its token.
Due to that the user will get an email verification mail to verify his email address.

To create the user with a verified email address you must add an action.

1. Go to the actions of your organization
2. Create a new action with the following code to set the email to verified automatically
3. Make sure the action name matches the function in the action itself e.g: "setEmailVerified"

```js reference
https://github.com/zitadel/actions/blob/main/examples/verify_email.js
```

![img.png](/img/guides/action_email_verify.png)

3. Add the action "email verify" to the flow "external authentication" and to the trigger "pre creation"

![img.png](/img/guides/action_pre_creation_email_verify.png)

#### Automatically redirect to Azure AD

If you like to get automatically redirected to your Azure AD login instead of showing the ZITADEL login with the Username/Password and a button "Login with AzureAD" you have to do the following steps:

1. Go to the login behaviour settings of your instance or organization
2. Disable login with username and password
3. Make sure you have only configured AzureAD as external identity provider
4. If you did all your settings on the organization level make sure to send the organization scope in your authorization request: [scope](/apis/openidoauth/scopes#reserved-scopes)

### Test the setup

To test the setup use incognito mode and browse to your login page.
If you succeeded you should see a new button which should redirect you to your AzureAD Tenant.

![AzureAD Button](/img/guides/azure_zitadel_button.png)

![AzureAD Login](/img/guides/azure_login.png)

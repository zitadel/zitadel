---
title: Brand Customization
---

ZITADEL offers various customization options for your projects brand design. The branding can be configured on two different levels.
The configuration on the instance level will set the default settings, which are triggered for all users if not overwritten on an organization specifically.
The second possibility is to configure it on each organization. This guide will describe the second possibility.
For this head over to the Branding Setting on your Organization Page.

## How it works

You are able to customize the light and a dark mode separately.
All your changes will be shown in the preview window on the right side.
As soon as you are happy with your configuration click the "Apply configuration" button.
After this your settings will trigger in your system. The login and the emails will be sent with your branding.

## Settings

![Private Labeling](/img/console_private_labeling.png)

### Logo

Upload your logo for the chosen theme, as soon as it is uploaded the preview on the right side of the screen should show it.

### Colors

In the next part you can configure your colors. 
Background colour is self-explanatory, the primary color will be used for buttons, links and some highlights. 
The warn color is used for all the error messages and warnings and the font colour for texts. 

### Font

Last step to apply to your branding is the font upload. 
The best way is to upload a ttf file after a successful upload you will see it in the font part, but not in the preview.

### Advanced Settings

In the advanced behavior you can choose if the loginname suffix (domain e.g road.runner@acme.caos.ch) should be shown in the loginname screen or not and if the “ZITADEL watermark” should be hidden.

## Trigger the private labeling for the login

If you like to trigger your settings for your applications you have different possibilities.

### 1. Primary Domain Scope

Send a [primary domain scope](../../../apis/openidoauth/scopes) with your [authorization request](../../integrate/login-users#auth-request) to trigger your organization.
The primary domain scope will restrict the login to your organization, so only users of your own organization will be able to login.

See the following link as an example. Users will be able to register and login to the organization that verified the @caos.ch domain only.

```
https://{your_domain.zitadel.cloud}/oauth/v2/authorize?client_id=69234247558357051%40zitadel&scope=openid%20profile%20urn%3Azitadel%3Aiam%3Aorg%3Adomain%3Aprimary%3Acaos.ch&redirect_uri=https%3A%2F%2Fconsole.zitadel.cloud%2Fauth%2Fcallback&state=testd&response_type=code&nonce=test&code_challenge=UY30LKMy4bZFwF7Oyk6BpJemzVblLRf0qmFT8rskUW0
```

:::info
Make sure to replace the domain `caos.ch` with your own domain to trigger the correct branding.
:::

:::caution
This example uses the ZITADEL Cloud Application for demonstration. You need to create your own auth request with your applications parameters. Please see the docs to construct an [Auth Request]../integrate/login-users#auth-request).
:::

### 2. Setting on your Project

Set the private labeling setting on your project to define which branding should trigger.

## Reset to default

If you don't like your customization anymore click the "reset policy" button.
All your settings will be removed and the default settings of the system will trigger.

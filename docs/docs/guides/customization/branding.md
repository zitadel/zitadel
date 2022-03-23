---
title: Brand Customization
---

ZITADEL lets you customize the interface for your project's brand design. 
To customize, head to your Organization page, then to **Private Labeling Policy**.

## How it works
You can customize light and dark modes separately.
All your changes are shown in the preview window on the right side.
As soon as you are happy with your configuration, select **Apply configuration**.
After this your settings will trigger in your system.
The login and the emails will be sent with your branding.

## Settings

![Private Labeling](/img/console_private_labeling.png)

### Logo
Upload your logo for the chosen theme.
As soon as it is uploaded, the preview on the right side of the screen shows it.

### Colors
In the next part you can configure your colors. 
* Background color is the background.
* The primary color is used for buttons, links and some highlights. 
* The warn color is used for all the error messages and warnings.
* The font color is for texts. 

### Font

The last step is to upload your font. 
The best way is to upload a TTF file.

After a successful upload, you will see it in the font part, but not in the preview.

### Advanced Settings

In the advanced behavior, you can choose whether to show the loginname suffix (domain e.g road.runner@acme.caos.ch) in the loginname screen.
You can also configure whether to hide the ZITADEL watermark.

## Trigger the private labeling for the login

If you want to trigger your settings for your applications, you have different possibilities.

### A. Primary Domain Scope
Send a [primary domain scope](https://docs.zitadel.ch/docs/apis/openidoauth/scopes#reserved-scopes) with your [authorization request](https://docs.zitadel.ch/docs/guides/authentication/login-users/#auth-request) to trigger your organization.
The primary domain scope will restrict the login to your organization.
Only users of your own organization will be able to login.

See the following link as an example.
Users can register and log in only to the organization that verified the @caos.ch domain.
```
https://accounts.zitadel.ch/oauth/v2/authorize?client_id=69234247558357051%40zitadel&scope=openid%20profile%20urn%3Azitadel%3Aiam%3Aorg%3Adomain%3Aprimary%3Acaos.ch&redirect_uri=https%3A%2F%2Fconsole.zitadel.ch%2Fauth%2Fcallback&state=testd&response_type=code&nonce=test&code_challenge=UY30LKMy4bZFwF7Oyk6BpJemzVblLRf0qmFT8rskUW0
```

:::info

Make sure to replace the domain `caos.ch` with your own domain to trigger the correct branding.

:::

:::caution

This example uses the ZITADEL Cloud Application for demonstration.
You need to create your own auth request with your application parameters.
Refer to the docs [to construct an Auth Request](https://docs.zitadel.ch/docs/guides/authentication/login-users/#auth-request).

:::



### B. Setting on your Project
Set the private labeling setting on your project to define which branding should trigger.

## Reset to default
If you don't like your customization anymore, select the **reset policy** button.
All your settings will be removed and the default settings of the system will trigger.

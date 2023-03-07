---
title: Troubleshoot
---

You will find some possible error messages here, what the problem is and what some possible solutions can be.

:::tip Didn't find an answer?  
Join or [Chat](https://zitadel.com/chat) or open a [Discussion](https://github.com/zitadel/zitadel/discussions).
:::


## User Agent does not correspond

This error appeared for some users as soon as they were redirected to the login page of ZITADEL.
ZITADEL uses some cookies to identify the browser/user agent of the user, so it is able to store the active user sessions. By blocking the cookies the functions of ZITADEL will be affected.

We only found this issue with iPhone users, and it was dependent on the settings of the device.

Go to the settings of the app Safari and check in the "Experimental WebKit Features" if SameSite strict enforcement (ITP) is disabled
Also check if "block all cookies" is active. If so please disable this setting.

To make sure, that your new settings will trigger, please restart your mobile phone and try it again.

**Settings > Safari > Advanced > Experimental Features > disable: „SameSite strict enforcement (ITP)“**

![Same Site Strict Enforvement](/img/manuals/errors/same-site-strict.png)

**Settings > Safari > disable: "Block All cookies"** 
![Block all cookies](/img/manuals/errors/block-cookies.png)

Do you still face this issue? Please contact us, and we will help you find out what the problem is.

## Instance not found

`ID=QUERY-n0wng Message=Instance not found`

If you're in an self-hosting scenario with a custom domain, you need to instruct ZITADEL to use the `ExternalDomain`.
You can find more instruction in our guide about [custom domains](https://zitadel.com/docs/self-hosting/manage/custom-domain).
We also provide a guide on how to [configure](https://zitadel.com/docs/self-hosting/manage/configure) ZITADEL with variables from files or environment variables.

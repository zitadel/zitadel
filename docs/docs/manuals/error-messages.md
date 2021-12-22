---
title: Error Messages
---

You will find some possible error messages here, what the problem is and what some possible solutions can be.

## User Agent does not correspond

This error appeared for some users as soon as they were redirected to the login page of ZITADEL.

We only found this issue with iPhone users and it was dependent on the settings of the device.

### Solution

Go to the settings of the app Safari and check in the "Experimental WebKit Features" if SameSite strict enforcement (ITP) is disabled
Also check if "block all cookies" is active. If so please disable this setting.

To make sure, that your new settings will trigger, please restart your mobile phone and try it again.

Do you still face this issue? Please contact us, and we will try to find out what the problem is.


![Same Site Strict Enforvement](/img/manuals/errors/same-site-strict.png)
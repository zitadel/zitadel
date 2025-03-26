---
title: ZITADEL Instances
sidebar_label: Instances
---

The ZITADEL Customer Portal is used to manage all your different ZITADEL instances.
Instances are containers for your organizations, users and projects.
A recommended setup could look like the following:
1. Instance: "Dev Environment"
2. Instance: "Test Environment"
3. Instance: "Prod Environment"

In the free subscription model you have one instance included.
To be able to add more instances please upgrade to "ZITADEL Pro".

## Overview

The overview shows all the instances that are registered for your customer.
You can directly see the custom domain and data region.
With a click on an instance you get to the detail of the chosen instance.


![Instance Overview](/img/manuals/portal/customer_portal_instance_overview.png)


## New instance

Click on the new button above the instance table to create a new instance.

1. Enter the name of your new instance
2. Add the credentials for your first administrator
   - Username (prefilled)
   - Password
3. Instance created! You can now see the details of your first instance.

:::info
Every new instance gets a generated domain of the form [instancename][randomnumber].zitadel.cloud
:::

![New Instance](/img/manuals/portal/customer_portal_new_instance.gif)

## Detail

The detail shows you general information about your instance, the region and your usage.

![New Instance](/img/manuals/portal/customer_portal_instance_detail.png)

### Upgrade to Pro

Your first instance is included in the free subscription.
As soon as you want to create your second instance or use a "pro" feature like choosing the region, you will have to upgrade to the Pro subscription.
To upgrade you must enter your billing information.

If you hit a limit from the free tier you will automatically be asked to add your credit card information and to subscribe to the pro tier.
You can also upgrade manually at any time.

1. Click the "Upgrade to PRO" button in the menu or go to the billing menu
2. If you choose the billing menu, you can now see your Free plan, click "Upgrade to Pro"
4. Add the missing data
   - Payment method: Credit Card Information
   - Customer: At least you have to fill the country
5. Save the information

![Upgrade to Pro](/img/manuals/portal/customer_portal_upgrade_tier.png)

### Add Custom Domain

We recommend register a custom domain to access your ZITADEL instance.
The primary custom domain of your ZITADEL instance will be the issuer of the instance. All other custom domains can be used to access the instance itself

1. Browse to the "Custom Domains" Tab
2. Click **Add domain**
3. Enter the domain you want and select the instance where the domain should belong to
4. In the next screen you will get all the information you will have to add to your DNS provider to verify your domain

> **_Please note:_** Do not delete the verification code, as ZITADEL Customer Portal will re-check the ownership of your domain from time to time

Be aware that it has some impacts if you change the primary domain of your instance.

1. The urls and issuer have to change in your app
2. Passkey authentication is based on the domain, if you change it, your users will not be able to login with the registered passkey authentication

![Add custom domain](/img/manuals/portal/customer_portal_add_domain.png)

#### Verify Custom Domain

As soon as you have added your custom domain you will have to verify it, by adding a CNAME record to your DNS provider.

1. Go to your DNS provider
2. Add a new CNAME record (You can find the target on the detail page of your instance)
3. After adding the CNAME you need to wait till the domain is verified (this can take some time)

You will now be able to use the added custom domain to access your ZITADEL instance

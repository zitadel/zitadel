---
title: Instances
---

The ZITADEL customer Portal is used to manage all your different ZITADEL instances.
You can also manage your subscriptions, billing, newsletters and support requests.

## Overview

The overview shows all the instances that are registered for a specific customer.
You can directly see what kind of subscription the instance has and in which data region it is stored.
With a click on a instance row you get to the detail of the chosen instance.

## New instance

Click on the new button above the instance table to create a new instance.

1. Enter the name of your new instance
2. Choose if you like to start with the free or the pay as you go tier
3. Choose your options (pay as you go)
   1. Data Region: The region where your data is stored
   2. Custom Domain: We generate a default domain ({instance-name}-{random-string}.zitadel.cloud), but you can choose you custom domain
   3. If our basic SLA and Support is not enough, you can extend it
4. Check the summary
5. Add you payment method (pay as you go)
6. Return to Customer Portal
7. Instance created!

You will get an email to initialize your first user of the instance and to access the new created ZITADEL instance.

:::info
Every new instance gets a generated domain of the form [instancename][randomnumber].zitadel.cloud
:::

![New Instance](/img/manuals/portal/customer_portal_new_instance.gif)

## Detail

The detail shows you general information about your instance, which options you have and your usage.

![New Instance](/img/manuals/portal/customer_portal_instance_detail.png)

### Upgrade Instance

A free instance can be upgraded to a "pay as you go" instance. By upgrading your authenticated request will no longer be capped and you will be able to choose more options. To upgrade you must enter your billing information.

1. Go to the detail of your instance
2. Click "Upgrade to paid tier!" in the General Information
3. Choose the options you need (can be changed later)
   1. Data Region
   2. Custom Domain
   3. Extended SLA
4. Add a payment method or choose an existing one

### Add Custom Domain

We recommend register a custom domain to access your ZITADEL instance.
The primary domain of your ZITADEL instance will be the issuer of the instance. All other domains can be used to access the instance itself

1. Browse to your instance
2. Click **Add custom domain**
3. To start the domain verification click the domain name and a dialog will appear, where you can choose between DNS or HTTP challenge methods.
4. For example, create a TXT record with your DNS provider for the used domain and click verify. ZITADEL will then proceed an check your DNS.
5. When the verification is successful you have the option to activate the domain by clicking **Set as primary**

> **_Please note:_** Do not delete the verification code, as ZITADEL Customer Portal will re-check the ownership of your domain from time to time

Be aware that it has some impacts if you change the primary domain of your instance.

1. The urls and issuer have to change in your app
2. Passwordless authentication is based on the domain, if you change it, your users will not be able to login with the registered passwordless authentication

![Add custom domain](/img/manuals/portal/portal_add_domain.png)

#### Verify Domain

If you need a custom domain for your ZITADEL instance, you need to verify the domain.

1. Go to your DNS provider
2. Add a new CNAME record (You can find the target on the detail page of your instance)
3. After adding the CNAME you need to wait till the domain is verified (this can take some time)

You will now be able to use the added custom domain to access your ZITADEL instance

### Change Options

You can change your selected options in the detail of your instance.
This can have an impact on your instance cost.

1. Go to the detail of your instance
2. Click the edit button on the Options section
3. Choose your options
   1. Extended SLA
   2. Data Region
4. Save

![Edit Options](/img/manuals/portal/portal_edit_options.png)

### Downgrade Instance

If you are in the "Pay as you go tier" with your instance, you can downgrade it to the free tier.

:::caution
Be aware that this might have an impact for your users and application.
If you have registered a custom domain, it will be deleted.
The data region will be set to "Global", if you have selected something else.
:::

1. Go to the detail of your instance
2. Click "Change to free tier" in the General Information
3. You will see an overview of what happens when downgrading, click "Downgrade anyway"
4. In the popup you need to confirm by clicking "I am sure"

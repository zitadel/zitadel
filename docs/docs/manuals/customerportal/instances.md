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
   a. Data Region: The region where your data is stored
   b. Custom Domain: We generate a default domain ({instance-name}-{random-string}.zitadel.cloud), but you can choose you custom domain
   c. If our basic SLA and Support is not enough, you can extend it
4. Check the summary
5. Add you payment method (pay as you go)
6. Return to customer portal
7. Instance created!

You will get an email to initialize your first user of the instance and to access the new created ZITADEL instance.

![New Instance](/img/manuals/portal/customer_portal_new_instance.gif)

## Detail

The detail shows you general information about your instance, which options you have and your usage.

![New Instance](/img/manuals/portal/customer_portal_instance_detail.png)

## Custom Domain

We recommend register a custom domain to access your ZITADEL instance.
The primary domain of your ZITADEL instance will be the issuer of the instance. All other domains can be used to access the instance itself

Be aware that it has some impacts if you change the primary domain of your instance.
1. The urls and issuer have to change in your app
2. Passwordless authentication is based on the domain, if you change it, your users will not be able to login with the registered passwordless authentication

### Verify Domain

If you need a custom domain for your ZITADEL instance, you need to verify the domain.

1. Go to your DNS provider
2. Add a new CNAME record (You can find the target on the detail page of your instance)
3. After adding the CNAME you need to wait till the domain is verified (this can take some time)

You will now be able to use the added custom domain to access your ZITADEL instance
---
title: Organisations
---

### What are organisations

Organisations are comparable to tenants of a system or OU's (organisational units) if we speak of a directory based system.
ZITADEL is organised around the idea that multiple organisations share the same [System](administrate#What_is_meant_by_system) and that they can grant each other rights so self manage certain things.

#### Global organisation

ZITADEL provides a global organisation for users who manage their accounts on their own. Think of this like the difference between a "Microsoft Live Login" vs. "AzureAD User"
or if you think of Google "Gmail" vs "Gsuite".

### Create an organisation without existing login

ZITADEL allows you to create a new organisation without a pre-existing user. For [ZITADEL.ch](https://zitadel.ch) you can create a org by visiting the [Register organisation](https://accounts.zitadel.ch/register/org)

> Screenshot here

For dedicated ZITADEL instances this URL might be different, but in most cases should be something like https://accounts.YOURDOMAIN.TLD/register/org

### Create an organisation with existing login

You can simply create a new organisation by visiting the [ZITADEL Console](https://console.zitadel.ch) and clicking "new organisation" in the upper left corner.

> Screenshot here

For dedicated ZITADEL instances this URL might be different, but in most cases should be something like `https://console.YOURDOMAIN.TLD`

### Verify a domain name

Once you created your organisation you will receive a generated domain name from ZITADEL for your organisation. For example if you call your organisation `ACME` you will receive `acme.zitadel.ch` as name. Furthermore the users you create will be suffixed with this domain name. To improve the user experience you can verify a domain name which you control. If you control acme.ch you can verify the ownership by DNS or HTTP challenge.
After the domain is verified your users can use both domain names to log-in. The user "coyote" can now use "coyote@acme.zitadel.ch" and "coyote@acme.ch".
An organisation can have multiple domain names, but only one of it can be primary. The primary domain defines which login name ZITADEL displays to the user, and also what information gets asserted in access_tokens (preferred_username).

Browse to your [organisation](administrate#Organisations) by visiting [https://console.zitadel.ch/org](https://console.zitadel.ch/org).

Add the domain to your [organisation](administrate#Organisations) by clicking the button **Add Domain**.
<img src="img/console_org_domain_default.png" alt="Organisation Overview" width="1000px" height="auto">

Input the domain in the input field and click **Add**
<img src="img/console_org_domain_add.png" alt="Organisation Add Domain" width="1000px" height="auto">

<img src="img/console_org_domain_added.png" alt="Organisation Domain Added" width="1000px" height="auto">

To start the domain verification click the domain name and a dialog will appear, where you can choose between DNS or HTTP challenge methods.
<img src="img/console_org_domain_verify.png" alt="Organisation Domain Verify" width="1000px" height="auto">

For example, create a TXT record with your DNS provider for the used domain an click verify. **ZITADEL** will then proceed an check your DNS.
<img src="img/console_org_domain_verify_dns.png" alt="Organisation Domain Verify DNS" width="1000px" height="auto">

> Do not delete the verification code **ZITADEL** will recheck the ownership from time to time

When the verification is successful you have the option to activate the domain by clicking **Set as primary**
<img src="img/console_org_domain_verified.png" alt="Organisation Domain Verified" width="1000px" height="auto">

> This changes the **preferred loginnames** of your [users](administrate#Users) as indicated [here](administrate#How_ZITADEL_handles_usernames).

Gratulations your are done! You can check this by visiting [https://console.zitadel.ch/users/me](https://console.zitadel.ch/users/me)
<img src="img/console_user_personal_info.png" alt="User Personal Information" width="1000px" height="auto">

<div class="my-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_verified.png" itemprop="contentUrl">
            <img src="img/console_org_domain_verified.png" itemprop="thumbnail" alt="Image description" />
        </a>
        <figcaption itemprop="caption description">Image caption</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_personal_info.png" itemprop="contentUrl">
            <img src="img/console_user_personal_info.png" itemprop="thumbnail" alt="Image description" />
        </a>
        <figcaption itemprop="caption description">Image caption</figcaption>
    </figure>
</div>


> This only works when the [user](administrate#Users) is member of this [organisation](administrate#Organisations)

### Audit organisation changes

All changes to the organisation are displayed on the organisation menu within [ZITADEL Console](https://console.zitadel.ch/org) organisation menu. Located on the right hand side under "activity".

> Screenshot here

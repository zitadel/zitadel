---
title: Organisations
---

### What are organisations

Organisations are comparable to tenants of a system or OU's (organisational units) if we speak of a directory based system.
ZITADEL is organised around the idea that multiple organisations share the same [System](administrate#What_is_meant_by_system) and that they can grant each other rights so self manage certain things.

#### Global organisation

ZITADEL provides a global organisation for users who manage their accounts on their own. Think of this like the difference between a "Microsoft Live Login" vs. "AzureAD User"
or if you think of Google "Gmail" vs "Gsuite".

[//]: # (@fforootd I don't understand the point of this chapter, and the comparison doesn't make much sense to me.)

### Create an organisation without existing login

ZITADEL allows you to create a new organisation without a pre-existing user. For [ZITADEL.ch](https://zitadel.ch) you can create a org by visiting the [Register organisation](https://accounts.zitadel.ch/register/org)

> Screenshot here

<details>
    <summary>
        Dedicated Instance
    </summary>
For dedicated ZITADEL instances this URL might be different, but in most cases should be something like https://accounts.YOURDOMAIN.TLD/register/org
</details>

### Create an organisation with existing login

You can simply create a new organisation by visiting the [ZITADEL Console](https://console.zitadel.ch) and clicking "new organisation" in the upper left corner.

> Screenshot here

<details>
    <summary>
        Dedicated Instance
    </summary>
For dedicated ZITADEL instances this URL might be different, but in most cases should be something like `https://console.YOURDOMAIN.TLD`
</details>

### Verify a domain name

Once you created your organisation you will receive a generated domain name from ZITADEL for your organisation. For example if you call your organisation `ACME` you will receive `acme.zitadel.ch` as name. Furthermore the users you create will be suffixed with this domain name. To improve the user experience you can verify a domain name which you control. If you control acme.ch you can verify the ownership by DNS or HTTP challenge.
After the domain is verified your users can use both domain names to log-in. The user "coyote" can now use "coyote@acme.zitadel.ch" and "coyote@acme.ch".
An organisation can have multiple domain names, but only one of it can be primary. The primary domain defines which login name ZITADEL displays to the user, and also what information gets asserted in access_tokens (preferred_username).

Browse to your [organisation](administrate#Organisations) by visiting [https://console.zitadel.ch/org](https://console.zitadel.ch/org).

Add the domain to your [organisation](administrate#Organisations) by clicking the button **Add Domain**.
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_default.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_default.png" itemprop="thumbnail" alt="Organisation Overview" />
        </a>
        <figcaption itemprop="caption description">Organisation Overview</figcaption>
    </figure>
</div>

Input the domain in the input field and click **Add**
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_add.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_add.png" itemprop="thumbnail" alt="Organisation Add Domain" />
        </a>
        <figcaption itemprop="caption description">Organisation Add Domain</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_added.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_added.png" itemprop="thumbnail" alt="Organisation Domain Added" />
        </a>
        <figcaption itemprop="caption description">Organisation Domain Added</figcaption>
    </figure>
</div>
To start the domain verification click the domain name and a dialog will appear, where you can choose between DNS or HTTP challenge methods.
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_verify.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_verify.png" itemprop="thumbnail" alt="Organisation Domain Verify" />
        </a>
        <figcaption itemprop="caption description">Organisation Domain Verify</figcaption>
    </figure>
</div>
For example, create a TXT record with your DNS provider for the used domain and click verify. **ZITADEL** will then proceed an check your DNS.
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_verify_dns.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_verify_dns.png" itemprop="thumbnail" alt="Organisation Domain Verify DNS" />
        </a>
        <figcaption itemprop="caption description">Organisation Domain Verify DNS</figcaption>
    </figure>
</div>

> Do not delete the verification code **ZITADEL** will recheck the ownership from time to time

When the verification is successful you have the option to activate the domain by clicking **Set as primary**

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_verified.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_verified.png" itemprop="thumbnail" alt="Organization Domain Verified" />
        </a>
        <figcaption itemprop="caption description">Organisation verified</figcaption>
    </figure>
</div>

> This changes the **preferred loginnames** of your [users](administrate#Users) as indicated [here](administrate#How_ZITADEL_handles_usernames).

Congratulations your are done! You can check this by visiting [https://console.zitadel.ch/users/me](https://console.zitadel.ch/users/me)
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_personal_info.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_personal_info.png" itemprop="thumbnail" alt="User Personal Information" />
        </a>
        <figcaption itemprop="caption description">User Personal Information</figcaption>
    </figure>
</div>

> This only works when the [user](administrate#Users) is member of this [organisation](administrate#Organisations)

### Manage Organisation ZITADEL Roles

[//]: # (@fforootd screenshot sais "projects" - bit confusing)
[//]: # (@fforootd Is a "Manager" another "Organisation owner", or is this a special role?)

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_manage_roles_1.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_manage_roles_1.png" itemprop="thumbnail" alt="Manage ZITADEL Roles 1" />
        </a>
        <figcaption itemprop="caption description">Manage ZITADEL Roles 1</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_manage_roles_2.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_manage_roles_2.png" itemprop="thumbnail" alt="Manage ZITADEL Roles 2" />
        </a>
        <figcaption itemprop="caption description">Manage ZITADEL Roles 2</figcaption>
    </figure>
</div>

### Audit organisation changes

All changes to the organisation are displayed on the organisation menu within [ZITADEL Console](https://console.zitadel.ch/org) organisation menu. Located on the right hand side under "activity".

> Screenshot here

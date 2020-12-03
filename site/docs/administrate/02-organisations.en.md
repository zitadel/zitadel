---
title: Organizations
---

### What are organizations

Organizations are comparable to tenants of a system or OU's (organizational units) if we speak of a directory based system.
ZITADEL is organized around the idea that 
* multiple organizations share the same [System](administrate#What_is_meant_by_system) 
* these organizations can grant each other rights to self-manage certain things (eg, delegating roles)
* organizations are a vessels for [users](administrate#What_are_users) and [projects](administrate#What_are_projects)

#### Global organization

The global organization holds users that are not assigned to any other organization in the [System](administrate#What_is_meant_by_system). Thus ZITADEL provides a global organization for users who manage their accounts on their own.

<details open>
    <summary>
        Example
    </summary>
Let's look at our example company `acme.ch`: Suppose ACME sells online-tickets for concert venues. ACME created an organization `iam` to manage their own enterprise users (employees) and projects to manage the provided services. They also created an organization `b2b-partner-1`, allowing the partner self-manage their access. A partner could be a concert venue, that can administrate the backend of the service (e.g. posting new concerts, setting up billing, ...), and you want to allow them to self-manage access of users (e.g. employees of the venue) to their backend. Lastly, the organization `global` holds all the b2c customers of `acme.ch` that registered to the service to buy concert tickets.
</details>

### Create an organization without existing login

ZITADEL allows you to create a new organization without a pre-existing user. For [ZITADEL.ch](https://zitadel.ch) you can create a org by visiting the [Register organization](https://accounts.zitadel.ch/register/org)

> Screenshot here

<details>
    <summary>
        Dedicated Instance
    </summary>
For dedicated ZITADEL instances this URL might be different, but in most cases should be something like https://accounts.YOURDOMAIN.TLD/register/org
</details>

### Create an organization with existing login

You can simply create a new organization by visiting the [ZITADEL Console](https://console.zitadel.ch) and clicking "new organization" in the upper left corner.

> Screenshot here

<details>
    <summary>
        Dedicated Instance
    </summary>
For dedicated ZITADEL instances this URL might be different, but in most cases should be something like `https://console.YOURDOMAIN.TLD`
</details>

### Verify a domain name

Once you created your organization you will receive a generated domain name from ZITADEL for your organization. For example if you call your organization `ACME` you will receive `acme.zitadel.ch` as name. Furthermore the users you create will be suffixed with this domain name. To improve the user experience you can verify a domain name which you control. If you control acme.ch you can verify the ownership by DNS or HTTP challenge.
After the domain is verified your users can use both domain names to log-in. The user "coyote" can now use "coyote@acme.zitadel.ch" and "coyote@acme.ch".
An organization can have multiple domain names, but only one of it can be primary. The primary domain defines which login name ZITADEL displays to the user, and also what information gets asserted in access_tokens (preferred_username).

Browse to your [organization](administrate#Organizations) by visiting [https://console.zitadel.ch/org](https://console.zitadel.ch/org).

Add the domain to your [organization](administrate#Organizations) by clicking the button **Add Domain**.
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_default.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_default.png" itemprop="thumbnail" alt="Organization Overview" />
        </a>
        <figcaption itemprop="caption description">Organization Overview</figcaption>
    </figure>
</div>

Input the domain in the input field and click **Add**
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_add.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_add.png" itemprop="thumbnail" alt="Organization Add Domain" />
        </a>
        <figcaption itemprop="caption description">Organization Add Domain</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_added.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_added.png" itemprop="thumbnail" alt="Organization Domain Added" />
        </a>
        <figcaption itemprop="caption description">Organization Domain Added</figcaption>
    </figure>
</div>
To start the domain verification click the domain name and a dialog will appear, where you can choose between DNS or HTTP challenge methods.
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_verify.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_verify.png" itemprop="thumbnail" alt="Organization Domain Verify" />
        </a>
        <figcaption itemprop="caption description">Organization Domain Verify</figcaption>
    </figure>
</div>
For example, create a TXT record with your DNS provider for the used domain and click verify. **ZITADEL** will then proceed an check your DNS.
<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_verify_dns.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_verify_dns.png" itemprop="thumbnail" alt="Organization Domain Verify DNS" />
        </a>
        <figcaption itemprop="caption description">Organization Domain Verify DNS</figcaption>
    </figure>
</div>

> Do not delete the verification code **ZITADEL** will recheck the ownership from time to time

When the verification is successful you have the option to activate the domain by clicking **Set as primary**

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_org_domain_verified.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_org_domain_verified.png" itemprop="thumbnail" alt="Organization Domain Verified" />
        </a>
        <figcaption itemprop="caption description">Organization verified</figcaption>
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

> This only works when the [user](administrate#Users) is member of this [organization](administrate#Organizations)

### Manage Organization ZITADEL Roles

You can assign users [management roles](https://docs.zitadel.ch/administrate#ZITADEL_s_management_Roles) to your new organization.

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

### Audit organization changes

All changes to the organization are displayed on the organization menu within [ZITADEL Console](https://console.zitadel.ch/org) organization menu. Located on the right hand side under "activity".

> Screenshot here

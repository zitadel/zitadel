---
title: Users
---

### What are users

In **ZITADEL** there are different [users](administrate#Users). Some belong to dedicated [organisations](administrate#Organisations) other belong to the [global organisation](administrate#Global_organisation). Some of them are human [users](administrate#Users) others are machines.
Nonetheless we treat them all the same in regard to [roles](administrate#Roles) management and audit trail.

#### Human vs. Service Users

The major difference between human vs. machine [users](administrate#Users) is the type of credentials that can be used: With machine [users](administrate#Users) there is only a non-interactive logon process possible. As such we utilize “JWT as Authorization Grant”.

> TODO Link to “JWT as Authorization Grant” explanation.

[//]: # (maybe add own section with concepts / knowledge. Same would go for concepts such as "authorization server"... see okta docs)

### How ZITADEL handles usernames

[//]: # (this part is convoluted an needs some more structure)

**ZITADEL** is built around the concept of [organisations](administrate#Organisations). Each [organisation](administrate#Organisations) has its own pool of usernames which includes human and service [users](administrate#Users).

For example a [user](administrate#Users) with the username `road.runner` can only exist once in the [organisation](administrate#Organisations) `ACME`. **ZITADEL** will automatically generate a "logonname" for each [user](administrate#Users) consisting of `{username}@{domainname}.{zitadeldomain}`. Without [verifying the domain name](administrate#Verify_a_domain_name) this would result in the logonname `road.runner@acme.zitadel.ch`. 

> If you use a dedicated instance **ZITADEL** replace `zitadel.ch` with your domain name.

If someone [verifies a domain name](administrate#Verify_a_domain_name) within the organisation, **ZITADEL** will generate additional logonames for each [user](administrate#Users) with the verified domain. For example if the domain is `acme.ch` the resulting logonname would be `road.runner@acme.ch` in addition to the already generated `road.runner@acme.zitadel.ch`.

> Domain verification also removes the logonname from all [users](administrate#Users), who might have used this combination in the [global organisation](administrate#Global_organisation).
> Relating to example with `acme.ch` if a user in the [global organisation](administrate#Global_organisation), let's call him `coyote`, used `coyote@acme.ch` this logonname will be replaced with `coyote@randomvalue.tld`
> **ZITADEL** notifies the user about this change

### Manage Users

#### Search Users

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_list_search.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_list_search.png" itemprop="thumbnail" alt="User list Search" />
        </a>
        <figcaption itemprop="caption description">User list Search</figcaption>
    </figure>
</div>

Image 1: User List Search

#### Create Users

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_list.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_list.png" itemprop="thumbnail" alt="User list" />
        </a>
        <figcaption itemprop="caption description">User list</figcaption>
    </figure>
</div>

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_create_form.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_create_form.png" itemprop="thumbnail" alt="User Create Form" />
        </a>
        <figcaption itemprop="caption description">User Create Form</figcaption>
    </figure>
</div>

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_create_done.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_create_done.png" itemprop="thumbnail" alt="User Create Done" />
        </a>
        <figcaption itemprop="caption description">User Create Done</figcaption>
    </figure>
</div>

#### Set Password

> Screenshot here

### Manage Service Users

> Screenshot here

### Manage User Authorisations

> Screenshot here

### Manage User ZITADEL Roles

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_manage_roles_1.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_manage_roles_1.png" itemprop="thumbnail" alt="Manage ZITADEL Roles 1" />
        </a>
        <figcaption itemprop="caption description">Manage ZITADEL Roles 1</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_manage_roles_2.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_manage_roles_2.png" itemprop="thumbnail" alt="Manage ZITADEL Roles 2" />
        </a>
        <figcaption itemprop="caption description">Manage ZITADEL Roles 2</figcaption>
    </figure>
</div>

### Audit user changes

> Screenshot here

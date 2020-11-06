---
title: Console
---

### What is Console

Console is the ZITADEL Graphical User Interface. 

ZITADEL Console can be reached at [console.zitadel.ch](https://console.zitadel.ch/).
 For dedicated ZITADEL instances this URL might be different, but in most cases should be something like `https://console.YOURDOMAIN.TLD`

#### ZITADEL Users

**Users** can manage some information on their own.

- profile information
- credentials
- external logins

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_entry.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_entry.png" itemprop="thumbnail" alt="User Entry" />
        </a>
        <figcaption itemprop="caption description">User Entry</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_user_personal_information.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_user_personal_information.png" itemprop="thumbnail" alt="User Personal Information" />
        </a>
        <figcaption itemprop="caption description">User Personal Information</figcaption>
    </figure>
</div>

#### ZITADEL Organisation Owners

Users who manage organisations (**organisation owners**) do this also with Console.

- Organisation settings (policies, domains, idps)
- Manage users
- Manage projects, clients and roles
- Give access to users

#### ZITADEL Administrators

For the **IAM Administrators** there is also a section in Console solely intended to manage the system.

- Check failed events
- Reset read models
- Manage system settings and policies

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_admin_entry.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_admin_entry.png" itemprop="thumbnail" alt="Adminstrator Entry" />
        </a>
        <figcaption itemprop="caption description">Adminstrator Entry</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_admin_system.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_admin_system.png" itemprop="thumbnail" alt="System Administration" />
        </a>
        <figcaption itemprop="caption description">System Administration</figcaption>
    </figure>
</div>

> ZITADEL does display a banner to warn the administrator that his account has elevated privileges!

### Technologies

Console is built with Angular and interfaces with ZITADEL by utilizing the GRPC APIs over GRPC-web.

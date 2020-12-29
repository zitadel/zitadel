---
title: System Administration
---

### What is meant by system

System describes the root of ZITADEL and includes all other elements like organisations, users and so on. Most of the time this part is managed by an user with the role IAM_OWNER.

### Default Policies

When ZITADEL is set up for the first time we establish certain default policies for the whole system.

> TODO Document default policy settings

### Manage Read Models

Read Models are a way to normalize data out of the event stream for certain aspects. For example there is a model which consists of logonname and the password hash so that the login process can query that data.

All read models are eventually consistent by nature and sometimes an administrator would like to verify they are still up-to-date.
In the ZITADEL Console is a section called administration available where the admin can check all read models and their current state.
There is even a possibility to regenerate a read model.

> When a read model is regenerated it might take up some time to be fully operational again
> Depending on the model which is regenerated this might have a operational impact for the end-users

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_iam_admin_views.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_iam_admin_views.png" itemprop="thumbnail" alt="IAM View Management" />
        </a>
        <figcaption itemprop="caption description">IAM View Management</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/console_iam_admin_failed.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/console_iam_admin_failed.png" itemprop="thumbnail" alt="IAM Failed Events" />
        </a>
        <figcaption itemprop="caption description">IAM Failed Events</figcaption>
    </figure>
</div>

> Additional infos to the architecture of ZITADEL is located in [Architecture Docs](documentation#Architecture)

### Secret Handling

ZITADEL stores secrets always encrypted or hashed in it's storage.
Whenever feasible we try to utilize public / private key mechanics to handle secrets.

**Encryption**
We use `AES256` as default mechanic for storing secrets.

**Password Hashing**
By default `bcrypt` is used with a salt of `14`.

> This mechanic is used for user passwords and client secrets

**Signing Keys**
These keys are randomly generated within ZITADEL and are rotated on a regular basis (e.g all 6h).

> Signing keys are stored with AES256 encryption

**TLS**
Under normal operations ZITADEL's API nodes are located behind a reverse proxy. So the TLS Key handling are out of context in this regard.
However ZITADEL can use TLS keys at runtime level.

> TODO Document TLS config

### IAM Configuration

> TODO Document ZITADEL config

### Audit system changes

> Screenshot here

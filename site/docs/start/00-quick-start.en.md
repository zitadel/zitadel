---
title: Quick-Start
description: A quick-start reference for the impatient reader.
---

> All documentations are under active work and subject to change soon!

### Try ZITADEL

You can either use [ZITADEL.ch](https://zitadel.ch) or deploy a dedicated **ZITADEL** instance.

### Use ZITADEL.ch

To register your free [organisation](administrate#Organisations), visit this link [register organisation](https://accounts.zitadel.ch/register/org).
After accepting the TOS and filling out all the required fields you will receive a email with further instructions.

<img id="org-register-img" style="cursor: pointer; border-radius: 8px;" src="img/accounts_org_register.png" alt="Organisation Register" width="800px" height="auto">
<script>
    var openPhotoSwipe = function() {
        console.log('show image');
        var pswpElement = document.querySelectorAll('.pswp')[0];
        var options = {
            history: false,
            focus: false,
            showAnimationDuration: 0,
            hideAnimationDuration: 0
        };
        var items = [
            {
                src: 'img/accounts_org_register.png',
                w: 1024,
                h: 768,
                msrc: 'path/to/small-image.jpg',
                title: 'Image Caption'
            },
        ];
        var gallery = new PhotoSwipe( pswpElement, PhotoSwipeUI_Default, items, options);
        gallery.init();
    }
    document.getElementById('org-register-img').onclick = openPhotoSwipe;
</script>

#### Verify your domain name (optional)

When you verify your domain you get the benefit that your [organisations](administrate#Organisations) [users](administrate#Users) can use this domain as **preferred loginname**. You find a more detailed explanation [How ZITADEL handles usernames](administrate#How_ZITADEL_handles_usernames).

The verification process is documented [here](administrate#Verify_a_domain_name)

#### Add Users to your organisation

To add new user just follow [this guide](administrate#Create_Users)

#### Setup an application

First [create a project](administrate#Create_a_project)

Then create within this [project](administrate#Projects) a [new client](administrate#Create_a_client)

The wizard should provide some guidance what client is the proper for you. If you are still unsure consult our [Integration Guide](integrate#Overview)

### Use ORBOS to install ZITADEL

> This will be added later on

### Install ZITADEL with static manifest

> This will be added later on

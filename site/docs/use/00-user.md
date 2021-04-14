---
title: User Manual
---

<i class="las la-book-reader" style="font-size: 100px; height: 100px; color:#6c8eef"></i>

> All documentations are under active work and subject to change soon!

### Organization and user registration

Zitadel allows users to register a organization and/or user with just a few steps.

A. Register an organization

 1. Create an organization
 2. Verify your email
 3. Login to Zitadel and manage the organization

B. Create User
 1. An administrator can create and manage users within console.

C. Enable selfregistration for User

 1. Create an organization as above
 2. Create custom policy
 3. Enable the "Register allowed" flag in the Login Policy
 4. Connect your application and add the applications [scope](https://docs.zitadel.ch/architecture/#Custom_Scopes) to the redirect URL.

This will enable the register option in the login dialog and will register the user within your organization if he does not already have an account.

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/register.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/register.gif" itemprop="thumbnail" alt="Register organization" />
        </a>
        <figcaption itemprop="caption description">Register Organization</figcaption>
    </figure>
      <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/create-user.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/create-user.gif" itemprop="thumbnail" alt="Create User" />
        </a>
        <figcaption itemprop="caption description">Create User</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/enable-selfregister.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/enable-selfregister.gif" itemprop="thumbnail" alt="enable selfregister" />
        </a>
        <figcaption itemprop="caption description">Enable Selfregister</figcaption>
    </figure>
      <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/email-verify.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/email-verify.gif" itemprop="thumbnail" alt="Verify EMail" />
        </a>
        <figcaption itemprop="caption description">Verify EMail</figcaption>
    </figure>
</div>


### Change EMail

To change your email address visit your Personal Information page and amend the email field.



<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/change-email.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/change-email.gif" itemprop="thumbnail" alt="Change EMail" />
        </a>
        <figcaption itemprop="caption description">Change EMail</figcaption>
    </figure>
</div>


### Verify Phone

tbd

### Change Password

To change your password you can hit the link right at the overview page. Alternatively  you can set it in the "Personal Information" page.



<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/change-password.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/change-password.gif" itemprop="thumbnail" alt="Change Password" />
        </a>
        <figcaption itemprop="caption description">Change Password</figcaption>
    </figure>
</div>

### Manage Multi Factor

To enable multifactor authentication visit the "Personal Information" page of your account and scroll to the "multifactor authentication". 
You can either:

1. Configure OTP (One Time Password)

An OTP application creates a dynamic Token that changes periodically and needs to be added in addition to the password. Install an aproppriate OTP application of your choice and register Zitadel. The most convenient way is to scan the QR code with your the application on your mobile device. 

> **Information:** Some example Authenticator Apps for OTP are: Google Authenticator, Microsoft Authenticator, Authy. You can choose the one you like the most.

2. Add U2F (Universal Second Factor)

Unuversal Second Factor basically is a piece of hardware such as an USB key that gets linked to your Identity and authorizes as second factor when a button on the device is pressed.

> **Information:**  some example Keys are [Solokeys](https://solokeys.com) or [Yubikey](https://www.yubico.com/) You can choose the one you like the most.






<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/enable-mfa-handling.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/enable-mfa-handling.gif" itemprop="thumbnail" alt="Enable Multi Factor" />
        </a>
        <figcaption itemprop="caption description">Enable Multi Factor</figcaption>
    </figure>
   <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/login-mfa.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/login-mfa.gif" itemprop="thumbnail" alt="Login Multi Factor" />
        </a>
        <figcaption itemprop="caption description">Login Multi Factor</figcaption>
    </figure>
</div>


### Identity Linking

To link an external Identity Provider with a Zitadel Account you have to:

1. choose your IDP
2. Login to your IDP

you can then either

1. link the Identity to an existing ZITADEL useraccount
2. auto register a new ZITADEL useraccount 


<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/linking-accounts.gif" itemprop="contentUrl" data-size="1100x906">
            <img src="img/linking-accounts.gif" itemprop="thumbnail" alt="Linking Accounts" />
        </a>
        <figcaption itemprop="caption description">Linking Accounts</figcaption>
    </figure>
</div>


#### Auto Register

see Identity Linking above


#### Manage Account Linking

You can manage the linked external IDP Providers within the "Personal Information" Page.


<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/manage-external-idp.png" itemprop="contentUrl" data-size="1710x747">
            <img src="img/manage-external-idp.png" itemprop="thumbnail" alt="manage external IDP" />
        </a>
        <figcaption itemprop="caption description">Linking Accounts</figcaption>
    </figure>
</div>





### Login User

<div class="zitadel-gallery" itemscope itemtype="http://schema.org/ImageGallery">
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/accounts_page.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/accounts_page.png" itemprop="thumbnail" alt="Login Username" />
        </a>
        <figcaption itemprop="caption description">Login Username</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/accounts_password.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/accounts_password.png" itemprop="thumbnail" alt="Login Password" />
        </a>
        <figcaption itemprop="caption description">Login Password</figcaption>
    </figure>
    <figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject">
        <a href="img/accounts_otp_verify.png" itemprop="contentUrl" data-size="1920x1080">
            <img src="img/accounts_otp_verify.png" itemprop="thumbnail" alt="Login OTP" />
        </a>
        <figcaption itemprop="caption description">Login OTP</figcaption>
    </figure>
</div>

---
title: User Manual
---


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

Register Organization
![Register Organization](/img/register.gif)


Create User
![Create User](/img/create-user.gif)


Enable Selfregister
![Enable Selfregister](/img/enable-selfregister.gif)


Verify EMail
![Verify EMail](/img/email-verify.gif)


### Change EMail

To change your email address visit your Personal Information page and amend the email field.


Change EMail
![Change EMail](/img/change-email.gif)


### Verify Phone

tbd

### Change Password

To change your password you can hit the link right at the overview page. Alternatively  you can set it in the "Personal Information" page.


Change Password
![Change Password](/img/change-password.gif)


### Manage Multi Factor

To enable multifactor authentication visit the "Personal Information" page of your account and scroll to the "multifactor authentication". 
You can either:

1. Configure OTP (One Time Password)

An OTP application creates a dynamic Token that changes periodically and needs to be added in addition to the password. Install an aproppriate OTP application of your choice and register Zitadel. The most convenient way is to scan the QR code with your the application on your mobile device. 

> **Information:** Some example Authenticator Apps for OTP are: Google Authenticator, Microsoft Authenticator, Authy. You can choose the one you like the most.

2. Add U2F (Universal Second Factor)

Unuversal Second Factor basically is a piece of hardware such as an USB key that gets linked to your Identity and authorizes as second factor when a button on the device is pressed.

> **Information:**  some example Keys are [Solokeys](https://solokeys.com) or [Yubikey](https://www.yubico.com/) You can choose the one you like the most.



Enable Multi Factor
![Enable Multi Factor](/img/enable-mfa-handling.gif)


Login Multi Factor
![Login Multi Factor](/img/login-mfa.gif)


### Identity Linking

To link an external Identity Provider with a Zitadel Account you have to:

1. choose your IDP
2. Login to your IDP

you can then either

1. link the Identity to an existing ZITADEL useraccount
2. auto register a new ZITADEL useraccount 


Linking Accounts
![Linking Accounts](/img/linking-accounts.gif)


####  Self Register

When self registration is enabled, users can register themselfes in the organanization without any administrative effort.



Self Register
![Self Register](/img/self-register.gif)


#### Manage Account Linking

You can manage the linked external IDP Providers within the "Personal Information" Page.


Manage External IDP
![Manage External IDP](/img/manage-external-idp.png)



### Login User


Login Username
![Login Username](/img/accounts_page.png)


Login Password
![Login Password](/img/accounts_password.png)


Login OTP
![Login OTP](/img/accounts_multifactor.png)

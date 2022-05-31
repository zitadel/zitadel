---
title: User Profile
---

To get to your user profile you have to login to your ZITADEL Console {your-domain}-{randomstring}.zitadel.cloud or {your-custom-domain}.
If you have no special permissions in the ZITADEL Console, you will get directly to your profile page. 
Otherwise click on your user avatar in the top right of the console. A menu will open, with the "Edit Account" button you will be redirected to your profile page.

## Loginname

You are able to login with some different login names. The login name consists of the username and the organization suffix. The organization suffix are the registered domains on your organization.
![Loginname](/img/manuals/console_profile_loginname.png)

## General

In the general section you can find your profile data and contact information.
In the profile data you can change the following data:
- Avatar
- Username
- Firstname
- Lastname
- Nickname
- Display Name
- Gender
- Language

In the contact information you can change your password, email and phone number. The Email and Phone number need to be verified.

![Profile](/img/manuals/console_profile.png)
### Change Password

Change your password by entering your old, new and new confirmation password.

![Change Password](/img/change_password.gif)

### Change Email

Click on the edit button next to the email to change your email address.
You will now get an email to verify that this is your account. This can take a moment. 
Click on the button in the mail to verify the address. If you now reload your profile page the email address should be shown as verified.

If you wait to long to verify the email, your code will probably be expired. 
The get a new verification mail click on "resend code" next to the "not verified" label.

The email doesn't need to be unique within the whole system.

### Change Phone number

The phone number is not mandatory withing ZITADEL. If you like to add it, you have to verify it. 

1. Click "edit button" and add your number
2. Get an SMS with a verification code to the added number
3. Click "Verify" below the added number
4. A popup with an Input field for your code will be shown
5. Enter the code a click "OK"

Your phone number should now be verified.

## Identity Providers

The identity provider section shows you, if you have linked an account from another system. (e.g. Google Account, Github, Azure AD, etc)
If you have some linked accounts, in this section you can remove them, if you don't need them anymore.

## Passwordless

ZITADEL provides some different authentication methods, passwordless is one of them.
Passwordless has two different types, system based or system independent.

If you use system based methods make sure to register all the different devices you need to login. (e.g. Notebook, Mobile Phone, etc)
 
Examples for passwordless authentication methods are: Fingerprint, Windows Hello, Face Recognition, etc.
For device independent authentication you can use some hardware tokens. e.g. Yubikey, Solokey, etc.

There are different options how to add a passwordless autehntication.
1. Add directly on the current device
2. Send a registration link to your email. You can open this email and use the link on any device you like to register
3. Generate a qr code with a registration link and scann the QR Code with the device where you like to register 

Make sure to add at least to different devices or a device independent method

![Add Passwordless fingerprint](/img/manuals/console_profile_passwordless.gif)

## Multifactor Authentication

Multifactor authentication means that after entering the password, you need some kind of second authentication.
At the moment ZITADEL provides Webauthn and OTP.
Webauthn uses your device to authenticate e.g Fingerprint, Face Recognition, Windows Hello.
OTP means One time password, to use this method you need to install some kind of Authenticator App like Google Authenticator, Authy, Microsoft Authenticator.

## Fingerprint, Security Keys, Face ID, etc.

Use a method that is provided by your device to authenticate yourself.

1. Click the button "Add Factor" in the multifactor authentication section of your profile
2. Choose Fingerprint, Security Keys, Face ID and others
3. Enter a name which identifies your authentication (e.g iPhone Road.Runner, Mac Book 1, Yubikey), The name is used for nothing just for yourself to recognize what you have registered.
4. Your device will show you a popup to choose what method you like to register
5. Choose the method ond follow the instructions (e.g. Scan your finger, Enter Pin, etc.)

![Add MFA Fingerprint](/img/manuals/console_profile_mfa_webauthn.gif)

### One time Password (OTP)

For One time password (OTP) you will need an Authenticator app of your choice that provides an authentication code.

1. Download an Authenticator App of your choice (e.g. Authy, Google Authenticator, Microsoft Authenticator, etc.)
2. Click the button "Add Factor" in the multifactor authentication section of your profile
3. Choose OTP (One-Time-Password)
4. Scan the QR Code with your app
5. Enter the code you get in the app in the Code input field

You will now be able to use otp as a second factor during the login process



## Authorization

In the authorization section you can see all the permissions and roles you have to some different applications.

## Memberships

Membership is the role model ZITADEL provides for itself. If you have any permissions to manage something within ZITADEL you will have a membership.
This memeberships are hierarchical and have the following layers:
- System
- Organization
- Project
- Granted Project

To read more about the different roles withing ZITADEL click [here](../concepts/structure/managers.md).

## Metadata

Sometimes it is needed to store some more data on a user. This data can be stored in the metadata.

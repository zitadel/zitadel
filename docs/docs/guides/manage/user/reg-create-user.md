---
title: Register and Create User
---

The ZITADEL API has different possibilities to create users.
This can be used, if you are building your own registration page.

Use the following API call to create your users:
[Create User (Human)](/apis/mgmt/management-service-import-human-user.api.mdx)

## With Username and Password

If you are collecting all the user information and a password, you can directly create the user with the password.
With the `password_change_required` flag you can choose if the user has to change the password on the first login or not.
This might make sense if an administrator created the user.

## With passwordless

You can directly ask for a link to create the passwordless registration for the user. 
Fill the user data and set the attribute `request_passwordless_registration` to true.
You will get a link for the registration and an expiration time in the response.

If you add `requestPlatformType` as query parameter to the link you can define what type the platform should be.
- **platform**: Device itself e.g. FaceID, Fingerprint etc.
- **crossPlatform** A hardware token e.g SoloKey
- **unspecified** The user is free to choose

If nothing is requested, the type will not be restricted and all possibilities of the device will be taken into account.

### Add passwordless to existing user

If you already have a user in ZITADEL, it is possible to add passwordless later.

[Add Passwordless Registration ](/apis/mgmt)

Send the user_id in the request and you will get a link and an expiration as response.
You can then customize the link the same as described above in the creation process.

The second possibility is to send the link directly to the user per email.
Use the following request in that case:

[Send Passwordless Registration ](/apis/mgmt)

## Verified Email Address

When creating users you can define wether you want to verify their email address before activating the user.

In case you want to create the user with an verified email address set the `isEmailVerified` property to `true`.
No initialization email will be sent to the user.

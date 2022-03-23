---
title: Factors
---

## Manage Multi Factor

To enable multifactor authentication visit the "Personal Information" page of your account and scroll to the "multifactor authentication".

:::caution
In order to avoid being locked out if a factor does not work, we recommend registering several options
:::

### Configure OTP (One Time Password)

An OTP application creates a dynamic Token that changes periodically and needs to be added in addition to the password. 
1. Install an appropriate OTP application of your choice
2. Click Add AuthFactor
3. Choose OTP Option
4. Scan the QR Code with you chosen authenticator app
5. Enter the code from your app in the ZITADEL Console

:::info
Some example Authenticator Apps for OTP are: Google Authenticator, Microsoft Authenticator, Authy. You can choose the one you like the most.
:::

![Add One Time Password](/img/manuals/console_add_otp.gif)

### Configure U2F (Universal Second Factor)

U2F is dependent on the device and browser you are currently working.
In general there might be the following possibilities:
- FingerScan
- FaceRecognition (e.g. FaceID)
- Hardware Tokens (e.g. YubiKey, Solokeys)

Hardware Tokens are basically a piece of hardware such as a USB key that gets linked to your identity and authorizes as second factor when a button on the device is pressed.

:::info
Some example Keys are [Solokeys](https://solokeys.com) or [Yubikey](https://www.yubico.com/) You can choose the one you like the most.
:::

![Add Universal Second Factor](/img/manuals/console_add_u2f.gif)




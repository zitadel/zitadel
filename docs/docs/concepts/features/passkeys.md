---
title: "Passkeys in ZITADEL: Passwordless phishing-resistant authentication"
sidebar_label: Passkeys
---

ZITADEL's passkeys feature enables passwordless authentication, offering a **smoother and more secure** login experience for your users. This document explains the essential details for developers.

### What are Passkeys?

Imagine signing in without passwords! Passkeys, replacing traditional passwords, leverage **public-key cryptography** similar to FIDO2 and WebAuthn. Users rely on their devices' **biometrics or PINs** for authentication, eliminating password burdens.

### Benefits for Developers

* **Enhanced Security:** Phishing-resistant passkeys minimize credential theft risks.
* **Streamlined User Experience:** Faster, easier logins free users from managing passwords.
* **Platform Agnostic:** Works across devices and platforms supporting passkeys.
* **Modern Standard:** Complies with the FIDO2 and WebAuthn standards.

### Features

* **Seamless Registration:** Create unique passkeys for users on various devices. Optionally pair them with specific users and choose cross-platform or platform-specific options.
* **User Control:** Users manage their passkeys directly through ZITADEL's self-service portal, allowing registration, viewing, and deletion.
* **Intuitive Login:** Users initiate passwordless login by selecting the passkey option and verifying themselves with the device's biometrics (fingerprint, face ID, etc.).
* **Robust Fallback:** Traditional password login remains available for users without passkeys.

### Developer Resources

* **Documentation:** Passkeys Guide: [https://zitadel.com/docs/guides/integrate/login-ui/passkey](/docs/guides/integrate/login-ui/passkey)
* **Create Passkey Registration Link API:** [https://zitadel.com/docs/guides/manage/user/reg-create-user](/docs/guides/manage/user/reg-create-user)

### Notes

* Passkey support is still evolving in browsers and platforms. Check compatibility for your target audience.
* ZITADEL actively develops its passkey features. Stay updated with documentation and releases.
* Passkeys are bound to your domain, thus we recommend configuring a [custom domain](/docs/concepts/features/custom-domain.md) before setting up passkeys.

Don't hesitate to ask if you have further questions about integrating passkeys in your ZITADEL application!
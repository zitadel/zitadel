---
title: Secrets
---

In this chapter you can find information of how ZITADEL processes and stores secrects. 

## Secrets Principles

ZITADEL uses the following principles when handling Secrets across their lifecycle:

- Automate rotation
- Limit lifetime
- Show only once
- Prefer public / private key concepts (FIDO2, U2F, JWT Profile, ...)
- Irreversible hash secrets / passwords
- Encrypt secrets storage
- When sent across unsecure channels (eMail, SMS, ...)
  - Forced changed through receiver
  - Verify that the secret can only be used once

## Secrets Storage

By default ZITADEL stores secrets from its users, clients as well as its generated secrets like signing keys in it database.
To protect the secrets against extraction from database as well as database dumps they are encrypted with AES256.

:::info
The key used to encrypt and decrypt the secrets in the ZITADEL database is called `masterkey` and needs to be exactly 32 bytes long.
:::

:::info
The only secrets stored outside of the Secrets Storage are the masterkey, the TLS Keys, the initial Admin User (including the password)
:::

## Secrets stored in the Secrets Storage

### Passwords / Client Secrets

ZITADEL does handle many different passwords and secrets. These include:

- Validation Secrets
- User Passwords
- Client Secrets

:::info
ZITADEL uses `bcrypt` by default to store all Passwords and Client Secrets in an non reversible way to further reduce the risk of a Secrets Storage breach.
:::

### Public Keys

ZITADEL does handle many different public keys. These include:

- FIDO2
- U2F
- JWT Profile
- Signing Keys

### Private Keys

The only private keys currently stored by ZITADEL are the signing keys used to sign Tokens.
Signing Keys are rotated in a default schedule of 6 hours by creating new key pairs each turn.

:::info
By default ZITADEL uses `RSA256`.
:::

### Validation Secrets

Validation Secrets are used for different purposes, these include:

- Verifying contact information like eMail, Phonenumbers
- Verifying proof of ownership over domain names (DNS)
- Resting accounts of users (password, MFA reset, ...)

:::info
All validation secrets are protected against reuse and are treated as passwords.
:::

### Unencrypted Secrets

Some secrets cannot be hashed because they need to be used in their raw form. These include:

- Client Secrets of Identity Providers (IdPs)
- Personal Access Tokens

## Secrets stored outside the Secrets Storage

### Masterkey

Since the Masterkey is used as means of protecting the Secrets Storage it cannot be stored in the storage.
You find [here the many ways how ZITADEL can consume the Masterkey](../../guides/manage/self-hosted/configure).

### TLS Material

ZITADEL does support end to end TLS as such it can consume TLS Key Material.
Please check our [TLS Modes documentation](../../guides/manage/self-hosted/tls_modes) for more details.

### Admin User

The initial Admin User of ZITADEL can be configured through [ZITADELs config options](../../guides/manage/self-hosted/configure).

:::info
To prevent elevated breaches ZITADEL forces the Admin Users password to be changed during the first login.
:::


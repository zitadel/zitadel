---
title: OpenID Connect and Oauth2 web keys
sidebar_label: Web keys
---

Web Keys in ZITADEL are used to sign and verify JSON Web Tokens (JWT).
ID tokens are created, signed and returned by ZITADEL when a OpenID connect (OIDC) or Oauth2
authorization flow completes and a user is authenticated.
Optionally zitadel can return JWTs for access tokens if the OIDC Application if configured for it.

## Introduction

ZITADEL uses asymmetric cryptography to sign and validate JWTs.
Keys are generated in pairs giving a private and public key.
Private keys are used to sign tokens.
Public keys are used to verify tokens.
OIDC clients need the public key to verify ID tokens.
Oauth2 API apps might need the public key if they want to client-side verification of a
JWT access tokens, instead of introspection.
ZITADEL uses the public key verification when API calls are made, the user info or introspection
endpoints are called with a JWT access token.

### JSON Web Key

ZITADEL implement the [RFC7517 - JSON Web Key (JWK)](https://www.rfc-editor.org/rfc/rfc7517) format for storage and distribution of public keys.
Web keys in ZITADEL support a number of [JSON Web Algorithms (JWA)](https://www.rfc-editor.org/rfc/rfc7518) for digital signatures:

| Identifier | Description                     |
| ---------- | ------------------------------- |
| RS256      | RSASSA-PKCS1-v1_5 using SHA-256 |
| RS384      | RSASSA-PKCS1-v1_5 using SHA-384 |
| RS512      | RSASSA-PKCS1-v1_5 using SHA-512 |
| ES256      | ECDSA using P-256 and SHA-256   |
| ES384      | ECDSA using P-384 and SHA-384   |
| ES512      | ECDSA using P-521 and SHA-512   |
| EdDSA      | EdDSA signature algorithms[^1]  |

[^1]: EdDSA refers to both Ed25519 and Ed448 curves. ZITADEL only supports Ed25519 with a SHA-512 hashing algorithm. EdDSA is for JSON Object Signing is defined in [RFC8037](https://www.rfc-editor.org/rfc/rfc8037).


## Web Key management

ZITADEL provides a resource API for management of web keys. The API allows the creation, activation, deletion and listing of web keys.
All public keys that are stored for an instance are served on the [JWKS endpoint](#json-web-key).
Applications need the the public keys for token verification and not all applications are capable of on-demand key fetching when receiving a token with an unknown key ID (`kid` header claim).
Instead, those application may do a refresh or only load keys at startup.

Using the web keys API, keys can be created and activated for signing later. This allows the keys to be distributed to the instance's apps and caches.
Once a key is deactivated, its public key will remain available for token verification on the web key is deleted.
Delayed deletion makes sure tokens that were signed before the key got deactivated remain valid.

### Creation

ss

## JSON web key set


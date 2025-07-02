---
title: OpenID Connect and Oauth2 web keys
sidebar_label: Web keys [Beta]
---

Web Keys in ZITADEL are used to sign and verify JSON Web Tokens (JWT).
ID tokens are created, signed and returned by ZITADEL when a OpenID connect (OIDC) or OAuth2
authorization flow completes and a user is authenticated.
Optionally, ZITADEL can return JWTs for access tokens if the OIDC Application is configured for it.

## Introduction

ZITADEL uses asymmetric cryptography to sign and validate JWTs.
Keys are generated in pairs resulting in a private and public key.
Private keys are used to sign tokens.
Public keys are used to verify tokens.
OIDC clients need the public key to verify ID tokens.
OAuth2 API apps might need the public key if they want to client-side verification of a
JWT access tokens, instead of [introspection](/docs/apis/openidoauth/endpoints#introspection_endpoint).
ZITADEL uses public key verification when API calls are made or when the userInfo or introspection
endpoints are called with a JWT access token.

### JSON Web Key

ZITADEL implements the [RFC7517 - JSON Web Key (JWK)](https://www.rfc-editor.org/rfc/rfc7517) format for storage and distribution of public keys.
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

### Client Algorithm Support

Before customizing the algorithm the instance admin **MUST** make sure the complete app and API ecosystem
supports the chosen algorithm.

When all OIDC applications of an instance use opaque access tokens, and they call APIs which only use
introspection for token validation, only the OIDC applications will need to support the chosen algorithm.
If JWT access tokens are used and APIs do public key verification, those APIs need to support the chosen algorithm as well.

RS256 is widely considered the default algorithm and must be supported by all OIDC/Oauth2 providers, relying parties and resource servers.
This is also the default ZITADEL uses when [creating web keys](#creation).
It might be reasonable to assume RS384 and RS512 are just as supported, because those are just variations on RSA based keys.
The ES256, ES384 and ES512 might have reasonable support as well,
ECDSA is part of the same [JSON Web Algorithms (JWA)](https://www.rfc-editor.org/rfc/rfc7518) as RSA.

EdDSA usage is defined in the supplemental [RFC8037](https://www.rfc-editor.org/rfc/rfc8037),
and therefore may be less supported than the others.
Also, the `at_hash` claim in the ID token is a hashed string of the access token.
The hasher is usually defined by the keys `alg` header. For example:

- RS256 defines an RSA key and a SHA256 hasher.
- ES512 defines an elliptic curve key with the P-512 and SHA512 hasher.

Unfortunately, there is no published standard for the `at_hash` hasher used for EdDSA.
In fact, EdDSA may use different curves and internally uses different hashers:

- ed25519 uses SHA512;
- ed448 uses SHAKE256;

This resulted in a [proposal](https://bitbucket.org/openid/connect/issues/1125/_hash-algorithm-for-eddsa-id-tokens) at
the Open ID workgroup to follow suit and use the same hashing algorithms for the `at_hash` claim.
This means both signers and verifiers can't know the hasher by the `alg` value alone and need to inspect `crv` value as well.
Since the decision in the proposal isn't published yet,
there is a big change some OIDC client libraries don't have proper support for EdDSA / ed25519.

The ZITADEL back-end is written in Go. The Go developers have denied ed448 curve implementations to be included.
Therefore ZITADEL only uses ed25519 with a SHA512.
The same counts for [zitadel/oidc](https://github.com/zitadel/oidc) Go library.

## Web Key management

ZITADEL provides a resource based [web keys API](/docs/apis/resources/webkey_service_v2).
The API allows the creation, activation, deletion and listing of web keys.
All public keys that are stored for an instance are served on the [JWKS endpoint](#json-web-key-set).
Applications need public keys for token verification and not all applications are capable of on-demand
key fetching when receiving a token with an unknown key ID (`kid` header claim).
Instead, those application may do a time-based refresh or only load keys at startup.

Using the web keys API, keys can be created first and activated for signing later.
This allows the keys to be distributed to the instance's apps and caches.
Once a key is deactivated, its public key will remain available for token verification until the web key is deleted.
Delayed deletion makes sure tokens that were signed before the key got deactivated remain valid.

When the `web_key` [feature](/docs/apis/resources/feature_service_v2/feature-service-set-instance-features) is enabled the first time,
two web key pairs are created with one activated.

### Creation

The web key [create](/docs/apis/resources/webkey_service_v3/zitadel-web-keys-create-web-key) endpoint generates a new web key pair,
using the passed generator configuration from the request. This config is a one-of field of:

- RSA
- ECDSA
- ED25519

When the request does not contain any specific configuration,
[RSA](#rsa) is used as default with the default options as described below:

```bash
curl -L 'https://$CUSTOM-DOMAIN/v2beta/web_keys' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
-d '{}'
```

#### RSA

The RSA generator config takes two enum values.

- The `bits` fields determines the size of the RSA key:
  - `RSA_BITS_2048` (**default**)
  - `RSA_BITS_3072`
  - `RSA_BITS_4096`
- The `hasher` field sets the hash mode and
  determines the `alg` header value of the web key:
  - `RSA_HASHER_SHA256` results in the RS256 algorithm header. (**default**)
  - `RSA_HASHER_SHA384` results in the RS384 algorithm header.
  - `RSA_HASHER_SHA512` results in the RS512 algorithm header.

For example, to create a RSA web key with the size of 3072 bits and the SHA512 algorithm (RS512):

```bash
curl -L 'https://$CUSTOM-DOMAIN/v2beta/web_keys' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
-d '{
  "rsa": {
    "bits": "RSA_BITS_3072",
    "hasher": "RSA_HASHER_SHA512"
  }
}'

```

#### ECDSA

The ECDSA generator config takes a single `curve` enum value which determines both the key's curve parameters and hashing algorithm:

- `ECDSA_CURVE_P256` uses the NIST P-256 curve and sets the ES256 algorithm header.
- `ECDSA_CURVE_P384` uses the NIST P-384 curve and sets the ES384 algorithm header.
- `ECDSA_CURVE_P512` uses the NIST P-512 curve and sets the ES512 algorithm header.

For example, to create a ECDSA web key with a P-256 curve and the SHA256 algorithm:

```bash
curl -L 'https://$CUSTOM-DOMAIN/v2beta/web_keys' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
-d '{
  "ecdsa": {
    "curve": "ECDSA_CURVE_P256"
  }
}'
```

#### ED25519

ED25519 is an EdDSA curve and currently the only EdDSA curve supported by ZITADEL.[^2]
No config is needed for ed25519 as its specification already includes the curve parameters.
ed25519 always uses the SHA512 hasher.

Note that the `alg` header for ed25519 is `EdDSA` and refers to both ed25519 and ed448 curves.
Both curves specify different hashers.
Clients which support both curves must inspect `crv` header value to assert the difference.

For example, to create a ed25519 web key:

```bash
curl -L 'https://$CUSTOM-DOMAIN/v2beta/web_keys' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
-d '{
  "ed25519": {}
}'
```

[^2]: The ZITADEL back-end is written in Go. The Go developers have denied ed448 curve implementations to be included.
  Therefore ZITADEL won't support this either.

### Activation

When a generated web key is [activated](/docs/apis/resources/webkey_service_v3/zitadel-web-keys-activate-web-key),
its private key will be used to sign new tokens.
There can be only one active key on an instance.
Activating a key implies deactivation of the previously active key.

Public keys on the [JWKS](#json-web-key-set) endpoint may be [cached](#caching).
Therefore it is advised to delay activation after generating a key,
at least for the duration of the max-age setting plus any time it might take for client applications to refresh.

### Deletion

Non-active keys may be [deleted](/docs/apis/resources/webkey_service_v3/zitadel-web-keys-delete-web-key).
Deletion also means tokens signed with this key become invalid.
Active keys can't be deleted.
As each public key is available on the [JWKS](#json-web-key-set) endpoint,
it is important to cleanup old web keys that are no longer needed.
Otherwise the endpoint's response size will only grow over time, which might lead to performance issues.

Once a key was activated and deactivated (by activation of the next key) deletion should wait:

- Until access and ID tokens are expired. See [OIDC token lifetimes](/docs/guides/manage/console/default-settings#oidc-token-lifetimes-and-expiration).
- ID tokens may be used as `id_token_hint` in authentication and end-session requests. The hint typically doesn't expire, but becomes invalid once the key is deleted.
  It might be desired to keep keys around long enough to minimalize user impact.

### Rotation example

This section gives an example on a key rotation strategy.
This strategy aims to fulfill the following requirements:

1. Web keys are rotated monthly.
2. Applications have enough time to see the next activated web key on the [JWKS](#json-web-key-set) endpoint.
3. Web keys are kept long enough to cover the access and ID token validity of 24 hours.
4. Web keys are kept long enough to to allow usage of the `id_token_hint` for at least 3 months.
  Users that haven't logged in / refreshed tokens with the client app for that period,
  will need to re-enter their username.

When the feature flag was enabled the first time, the instance got two keys with the first one activated. When this feature becomes general available, instance creation will setup the first two keys in the same way. So the initial state always looks like this:

| id  | created    | changed    | state           |
| --- | ---------- | ---------- | --------------- |
| 1   | 2025-01-01 | 2025-01-01 | `STATE_ACTIVE`  |
| 2   | 2025-01-01 | 2025-01-01 | `STATE_INITIAL` |

For the sake of this example we will use simplified IDs and restrict timestamps to dates.

After one month, on 2025-02-01, we wish to activate the next available key and create a new key to be available for activation next month. This fulfills requirements 1 and 2.

```bash
curl -L -X POST 'https://$CUSTOM-DOMAIN/v2beta/web_keys/2/_activate' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>'

curl -L 'https://$CUSTOM-DOMAIN/v2beta/web_keys' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
-d '{}'
```

Key ID 2 became active, Key ID 1 became inactive and a new key with ID 3 was created:

| id  | created    | changed    | state            |
| --- | ---------- | ---------- | ---------------- |
| 1   | 2025-01-01 | 2025-02-01 | `STATE_INACTIVE` |
| 2   | 2025-01-01 | 2025-02-01 | `STATE_ACTIVE`   |
| 3   | 2025-02-01 | 2025-02-01 | `STATE_INITIAL`  |

No keys are deleted yet.
We continue like this monthly.
At one point (on 2025-05-01) we will have a web key with `STATE_INACTIVE` with a changed date of 3 months ago:

| id  | created    | changed    | state            |
| --- | ---------- | ---------- | ---------------- |
| 1   | 2025-01-01 | 2025-02-01 | `STATE_INACTIVE` |
| 2   | 2025-01-01 | 2025-03-01 | `STATE_INACTIVE` |
| 3   | 2025-02-01 | 2025-04-01 | `STATE_INACTIVE` |
| 4   | 2025-03-01 | 2025-05-01 | `STATE_INACTIVE` |
| 5   | 2025-04-01 | 2025-05-01 | `STATE_ACTIVE`   |
| 6   | 2025-05-01 | 2025-05-01 | `STATE_INITIAL`  |

In addition to the activate and create calls we made on this iteration,
we can now safely delete the oldest key, as both requirement 3 and 4 are now fulfilled:

```bash
curl -L -X DELETE 'https://$CUSTOM-DOMAIN/v2beta/web_keys/1' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>'
```

The final state:

| id  | created    | changed    | state            |
| --- | ---------- | ---------- | ---------------- |
| 2   | 2025-01-01 | 2025-03-01 | `STATE_INACTIVE` |
| 3   | 2025-02-01 | 2025-04-01 | `STATE_INACTIVE` |
| 4   | 2025-03-01 | 2025-05-01 | `STATE_INACTIVE` |
| 5   | 2025-04-01 | 2025-05-01 | `STATE_ACTIVE`   |
| 6   | 2025-05-01 | 2025-05-01 | `STATE_INITIAL`  |

Next month, Key ID 6 will be activated, an new key added and Key ID 2 can be deleted.

## JSON web key set

The JSON web key set (JWKS) endpoint serves all available public keys for the instance on
`{your_domain}/oauth/v2/keys`. This includes activated, newly non-activated and deactivated web keys. The response format is defined in [RFC7517, section 5: JWK Set Format](https://www.rfc-editor.org/rfc/rfc7517#section-5).

And looks like:

```json
{
  "keys": [
    {
      "use": "sig",
      "kty": "RSA",
      "kid": "280543383892525058",
      "alg": "RS384",
      "n": "0pVcbjTEr-awBmvztGLbBJB_-_YwjCKKXURJRpoXrChlaqtAvbkxby7mu9wSKAibxnvaobfuxnQydlB4CoKObUr00ARVBNeP5HLzeQUEx3CZh3s1LsjiuYov_yyvK9D12WH1LikP4ZPS68j-DVoEOEcFAE6cNikXTeDyCKa-ixROALieRXUQXTlvVyA_s0FhevmH0-M6rEN4YcfQuIZACEv2nQ4AJo0sNnugwrrqNn595ONKMSh2XTVngxxAD3TGHXg9bELB-WmgnZamVbO-ObpDBp5Ov73HL60_UoBTzBDECM6ovl52fHusLFw6Vkdt9_W3QhuRFljNqTPnna6rB-bLptQltBpnSBV3TxmklBcQ1EO3qeGvgOJsmDwSRlr28Du_1pyFMFANnG174eX5XrYASqTgJ1Wq7AfMBmv7YwGU7PbMce1V_CAV9u_hNkMJf0xQ4AIqrQ98f9hC5VCdCoKSOH1-1d8icEu7UmDyJohWqvY7xGOM_0Abx8ekMRT2O9PulmQ22me_GI5zXh7iv9yaoNq8EUNP5bdtr-ZG4PG8mqpLDSLpCpobYRK5AynyJkf-7_6neSy-ihu604ADKsNzB-uO58V8MPFdSPncyuUeTPX4dAVajbFyMtoAjtI1k_HYMU8nojRUrLSCJae9b0KtcPm9s7dCIL1Zpa4B-YM",
      "e": "AQAB"
    },
    {
      "use": "sig",
      "kty": "OKP",
      "kid": "280998627474669570",
      "crv": "Ed25519",
      "alg": "EdDSA",
      "x": "B51hFhRUHMHpqO1f-OThtnk3PfnRFaPFJWCLXSM_kuI"
    },
    {
      "use": "sig",
      "kty": "EC",
      "kid": "282465789963927554",
      "crv": "P-256",
      "alg": "ES256",
      "x": "X5s3tNoIXd5odp_-IwQq5oaAgMSoAxj0hwQ1DgHihmI",
      "y": "JqmTlRjoOv5bY5E9tAZXHaUHUamAAAFshO8zLhEZ9ZM"
    }
  ]
}
```

After the `web_key` feature is enabled, the response may still contain legacy keys, in order not to invalidate older sessions.
The legacy keys will disappear once they expire.

### Caching

As web keys can be created and distributed ahead of time, it is safe for JWKS responses to be cached at intermediate proxies.
Once the `web_key` feature is enabled, ZITADEL will send a `Cache-Control` header which allows caching.

By default and in ZITADEL Cloud we allow 5 minutes of caching:

```
Cache-Control: max-age=300, must-revalidate
```

Self-hosters can modify this setting through the `ZITADEL_OIDC_JWKSCACHECONTROLMAXAGE`
environment variable or in the configuration yaml:

```yaml
OIDC:
  JWKSCacheControlMaxAge: 5m
```

Setting the value to `0` will result in a `no-store` value in the `Cache-Control` header.

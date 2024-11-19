---
title: Migrate Users
sidebar_label: Users
---

Migrating users from an existing system, while minimizing impact on said users, can be a challenging task.

## Individual Users

Creating individual users can be done with this endpoint: [ImportHumanUser](/docs/apis/resources/mgmt/management-service-import-human-user).
Please also consult our [guide](/docs/guides/manage/user/reg-create-user) on how to create users.

```json
{
  "userName": "test9@test9",
  "profile": {
    "firstName": "Road",
    "lastName": "Runner",
    "displayName": "Road Runner",
    "preferredLanguage": "en"
  },
  "email": {
    "email": "test@test.com",
    "isEmailVerified": false
  },
  "hashedPassword": {
    "value": "$2a$14$aPbwhMVJSVrRRW2NoM/5.esSJO6o/EIGzGxWiM5SAEZlGqCsr9DAK",
    "algorithm": "bcrypt"
  },
  "passwordChangeRequired": false,
  "otpCode": "testotp",
  "requestPasswordlessRegistration": false,
  "idps": [
    {
      "configId": "124425861423228496",
      "externalUserId": "roadrunner@mailonline.com",
      "displayName": "name"
    }
  ]
}
```

## Bulk import

For bulk import use the [import endpoint](/docs/apis/resources/admin/admin-service-import-data) on the admin API:

```json
{
  "timeout": "10m",
  "data_orgs": {
    "orgs": [
      {
        "orgId": "104133391254874632",
        "org": {
          "name": "ACME"
        },
        "humanUsers": [
          {
            "userId": "104133391271651848",
            "user": {
              "userName": "test9@test9",
              "profile": {
                "firstName": "Road",
                "lastName": "Runner",
                "displayName": "Road Runner",
                "preferredLanguage": "de"
              },
              "email": {
                "email": "test@acme.tld",
                "isEmailVerified": true
              },
              "hashedPassword": {
                "value": "$2a$14$aPbwhMVJSVrRRW2NoM/5.esSJO6o/EIGzGxWiM5SAEZlGqCsr9DAK",
                "algorithm": "bcrypt"
              }
            }
          },
          {
            "userId": "120080115081209416",
            "user": {
              "userName": "testuser",
              "profile": {
                "firstName": "Test",
                "lastName": "User",
                "displayName": "Test User",
                "preferredLanguage": "und"
              },
              "email": {
                "email": "fabienne@caos.ch",
                "isEmailVerified": true
              },
              "hashedPassword": {
                "value": "$2a$14$785Fcdbpo9rn5L7E21nIAOJvGCPgWFrZhIAIfDonYXzWuZIKRAQkO",
                "algorithm": "bcrypt"
              }
            }
          },
          {
            "userId": "145195347319252359",
            "user": {
              "userName": "wile@test9",
              "profile": {
                "firstName": "Wile E.",
                "lastName": "Coyote",
                "displayName": "Wile E. Coyote",
                "preferredLanguage": "en"
              },
              "email": {
                "email": "wile.e@acme.tld"
              }
            }
          }
        ]
      }
    ]
  }
}
```

:::info
We will improve the bulk import interface for users in the future.
You can show your interest or join the discussion on [this issue](https://github.com/zitadel/zitadel/issues/5524).
:::

## Migrate secrets

Besides user data you need to migrate secrets, such as password hashes, OTP seeds, and public keys for passkeys (FIDO2).
The snippets in the sections below are parts from the bulk import endpoint, to clarify how the different objects can be imported.

### Passwords

ZITADEL stores passwords only as irreversible hashes, never in clear text.
Existing password hashes can be imported if they use a supported [hash algorithm](/docs/concepts/architecture/secrets#hashed-secrets).

Import password hashes using the import API (snippet from [bulk-import](#bulk-import)):
```json
{
  "userName": "test9@test9",
    ...,
    "hashedPassword": {
        "value": "$2a$14$aPbwhMVJSVrRRW2NoM/5.esSJO6o/EIGzGxWiM5SAEZlGqCsr9DAK",
        "algorithm": "bcrypt"
    },
    "passwordChangeRequired": false,
    ...,
}
```

Upon initial login, ZITADEL validates the imported password using the appropriate verifier.

:::info Verifiers
In ZITADEL, a password verifier checks the validity of a password hash created with an algorithm different from the currently configured one.
It acts as a translator, allowing ZITADEL to understand and validate hashes made with older algorithms like MD5 even when the system has transitioned to newer ones like Argon2.  
This is crucial during migrations or when importing user data.
Essentially, a verifier ensures ZITADEL can work with passwords hashed using various algorithms, maintaining security while transitioning to stronger hashing methods.
:::

Regardless of the `passwordChangeRequired` setting, the password is rehashed using the configured hasher algorithm and stored.
This ensures consistency and allows for automatic updates even when hasher configurations are changed, such as increasing salt cost for bcrypt.

To configure the default hasher for new user passwords, set the `Algorithm` of the `PasswordHasher` in the [runtime configuration file](/docs/self-hosting/manage/configure#runtime-configuration-file)
or by the environment variable `ZITADEL_SYSTEMDEFAULTS_PASSWORDHASHER_HASHER_ALGORITHM`, for example:

```
ZITADEL_SYSTEMDEFAULTS_PASSWORDHASHER_HASHER_ALGORITHM='pbkdf2'
```

Hasher configuration updates will automatically rehash existing passwords when they are validated or changed.


In case the hashes can't be transferred directly, you always have the option to create a user in ZITADEL without password and prompt users to create a new password.

If your legacy system receives the passwords in clear text (eg, login form) you could also directly create users via ZITADEL API.
We will explain this pattern in more detail in this guide.

### One-time passwords (OTP)

You can pass the OTP secret when creating users:

_snippet from [bulk-import](#bulk-import) example:_
```json
{
  "userName": "test9@test9",
    ...,
    "otpCode": "testotp",
    ...,
}
```

### Passkeys

When creating new users, you can trigger a workflow that prompts the users to setup a passkey authenticator.

_snippet from [bulk-import](#bulk-import) example:_
```json
{
  "userName": "test9@test9",
    ...,
    "requestPasswordlessRegistration": false,
    ...,
}
```

For passkeys to work on the new system you need to make sure that the new auth server has the same domain as the legacy auth server.

:::info
Currently it is not possible to migrate passkeys directly from another system.
:::

## Users linked to an external IDP

A users `sub` is bound to the external [IDP's Client ID](/docs/guides/manage/console/default-settings#identity-providers).
This means that the IDP Client ID configured in ZITADEL must be the same ID as in the legacy system.

Users should be imported with their `externalUserId`.

_snippet from [bulk-import](#bulk-import) example:_
```json
{
  "userName": "test9@test9",
    ...,
    "idps": [
        {
        "configId": "124425861423228496",
        "externalUserId": "roadrunner@mailonline.com",
        "displayName": "name"
        }
    ...,
}
```

You can use an Action with [post-creation flow](/docs/apis/actions/external-authentication#post-creation) to pull information such as roles from the old system and apply them to the user in ZITADEL.

## Metadata

You can store arbitrary key-value information on a user (or Organization) in ZITADEL.
Use metadata to store additional attributes of the users, such as organizational unit, backend-id, etc.

:::info
Metadata must be added to users after the users were created. Currently metadata can't be added during user creation.  
[API reference: User Metadata](/docs/apis/resources/mgmt/user-metadata)
:::

Request metadata from the userinfo endpoint by passing the required [reserved scope](/docs/apis/openidoauth/scopes#reserved-scopes) in your auth request.
With the [complement token flow](/docs/apis/actions/complement-token), you can also transform metadata (or roles) to custom claims.

## Authorizations / Roles

You can assign roles from owned or granted projects to a user.

:::info
Authorizations must be added to users after the users were created. Currently metadata can't be added during user creation.  
[API reference: User Authorization / Grants](/docs/apis/resources/auth/user-authorizations-grants)
:::
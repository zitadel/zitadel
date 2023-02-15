---
title: Users
---

Migrating users from an existing system, while minimizing impact on said users, can be a challenging task.
This guide gives you an overview of technical considerations, explains the most common patterns for migrating users and gives you some implementation details.

We will also offer more detailed guides how to migrate users from a specific auth provider to ZITADEL.

## Technical Considerations

### Evaluating migration patterns

There will be multiple ways for migrating users from your existing auth system ("legacy") to ZITADEL.
Which migration pattern to use depends on your requirements.

This section should help you to get an overview of the different migration patterns and help you design an ideal solution for your use case. Your solution might require adjustments from the presented baseline patterns.

#### Batch vs. Just-in-time Migration

```mermaid
%%{init: {'theme':'dark'}}%%
flowchart LR
    start([Start]) --> downtime{Zero downtime?}
    downtime -- No --> batch[[Batch migration]]
    downtime -- Yes --> clients{Can apps</br>switch</br>at day0?}
    subgraph jit [Just-in-time Migration]
        clients -- Yes --> user_api
        user_api{User API</br>available</br>on legacy?} -- Yes --> jit_zitadel[[ZITADEL</br>orchestrates migration]]
        user_api -- No --> jit_legacy
        clients -- No --> jit_legacy[[Legacy system</br>orchestrates migration]]
    end
```

[Batch migration](#batch-migration) is the easiest way, if you can afford some minimal downtime to move all users and applications over to ZITADEL.

In case all your applications depend on ZITADEL after the migration date, and ZITADEL is able to retrieve the required user information, including secrets, from the legacy system, then the recommended way is to let [ZITADEL orchestrate the user migration](#just-in-time-zitadel).

For all other cases, the legacy system needs to orchestrate the migration of users to ZITADEL for most flexibility.

#### Legacy System orchestrates migration

```mermaid
%%{init: {'theme':'dark'}}%%
flowchart LR
    jit_legacy[[Legacy system</br>orchestrates migration]] --> parallel{Requires session</br>on both, legacy</br>and ZITADEL?}
    parallel -- No --> import([Import user or</br> password reset])
    parallel -- Yes --> brokering([Identity Brokering + Action])
```

### Migrating Secrets

- Hashes
- Passkeys
- OTP
- For simplicity only once mentioned: If the hash algorithm is not available or the password can't be migrated, then use password reset flow; or
- Can also capture credentials behind login and provision user to ZITADEL

```mermaid
%%{init: {'theme':'dark'}}%%
flowchart LR
    batch[[Batch migration]] --> hash{Hash algo.</br>supported?}
    hash -- No --> reset([Use password reset])
    hash -- Yes --> import([Import user])
```

### Users linked to an external IDP

TODO: https://github.com/zitadel/zitadel/issues/5176

### JWT IDP

TODO: 

## Migration Patterns

### Batch migration

TODO: Chart
TODO: Example API Call - import user
TODO: Example API Call - password reset

### Just-in-time: ZITADEL

TODO: Chart
TODO: Example Action (HTTP, Metadata)?

### Just-in-time: Legacy

#### Provision users from legacy to ZITADEL

```mermaid
%%{init: {'theme':'dark'}}%%
sequenceDiagram
    actor U as User
    participant L as Legacy
    participant Z as ZITADEL
    
    U ->> L: Enter Username
    L -->> L: Flagged as migrated?
    alt migrated
        L ->> Z: redirect with login hint
    else not migrated
        L ->> Z: create user
        Z ->> L: response
        L -->> L: flag as migrated
        L ->> Z: redirect with login hint
    end

```

#### Identity Brokering and Action (parallel sessions)

TODO: Normal SSO with sessions in Legacy and ZITADEL (user already migrated)

TODO:
- IDP can be OIDC compliant / LDAP; or
- JWT IDP

```mermaid
%%{init: {'theme':'dark'}}%%
sequenceDiagram
    actor U as User
    participant A as App
    participant L as Legacy
    participant Z as ZITADEL
    
    U ->> A: start login
    A ->> L: redirect to IDP
    L ->> L: check if user is migrated

    L ->> Z: send auth request
    Z ->> L: redirect to IDP
    U ->> L: login with credentials

    L -->> L: create session (auto)
    L ->> Z: redirect
    opt create user
        Z -->> Z: Auto-register user
        Z -->> L: (Action) Get password / user info
        L -->> Z: (Action) Update user
        Z -->> L: (Action) Flag user as created
    end
    Z -->> Z: create session (auto)
    Z ->> A: redirect to callback url
```
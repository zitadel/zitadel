---
title: Policies
---

### What are policies

Policies are a means of enforcing certain behaviour of ZITADEL.
ZITADEL defines a default policy on the system level. However an organisation owner can change these aspects within his own organisation.

Below is a list of available policies

### Password complexity

This policy enforces passwords of users within the org. to be compliant.

- min length
- has number
- has symbol
- has lower case
- has upper case

> Screenshot here

### IAM Access Preference

This policy enforces, when set to true, that usernames are suffixed with the organisations domain.
Under normal operation this policy is only false on the `global` org. so that users can choose their email as their username.
Only available for the `IAM Administrator`

> Screenshot here

### Login Options

With this policy it is possible to define what options a user sees in the login process.

- Username Password allowed
- Self Register allowed
- External IDP allowed
- List of allowed external IDPs

> Screenshot here

### Audit policy changes

> Screenshot here

### Upcoming Policies

- Password age
- Password failure count

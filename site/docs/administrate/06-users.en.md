---
title: Users
---

### What are users

In ZITADEL there are different users. Some belong to dedicated organisations other belong to the global org. Some of them are human users others are machines.
Nonetheless we treat them all the same in regard to roles management and audit trail.

#### Human vs. Service Users

The major difference between humane vs. machine users is the type of credentials who can be used.
With machine users there is only a non interactive login process possible. As such we utilize “JWT as Authorization Grant”.

> TODO Link to “JWT as Authorization Grant” explanation.

### How ZITADEL handles usernames

ZITADEL is built around the concept of organisations. Each organisation has it's own pool of usernames which include human and service users.
For example a user with the username `alice` can only exist once the org. `ACME`. ZITADEL will automatically generate a "logonname" for each user consisting of `{username}@{domainname}.{zitadeldomain}`. Without verifying the domain name this would result in the logonname `alice@acme.zitadel.ch`. If you use a dedicated ZITADEL replace `zitadel.ch` with your domain name.

If someone verifies a domain name within the org. ZITADEL will generate additional logonames for each user with that domain. For example if the domain is `acme.ch` the resulting logonname would be `alice@acme.ch` and as well the generated one `alice@acme.zitadel.ch`.

> Domain verification also removes the logonname from all users who might have used this combination in the global org.
> Relating to example with `acme.ch` if a user in the global org, let's call him `bob` used `bob@acme.ch` this logonname will be replaced with `bob@randomvalue.tld`
> ZITADEL notifies the user about this change

### Manage Users

#### Create User

> Screenshot here

#### Set Password

> Screenshot here

### Manage Service Users

> Screenshot here

### Authorizations

> Screenshot here

### Audit user changes

> Screenshot here

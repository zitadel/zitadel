---
title: Users
---

### What are users

In **ZITADEL** there are different [users](administrate#Users). Some belong to dedicated [organisations](administrate#Organisations) other belong to the global [organisations](administrate#Organisations). Some of them are human [users](administrate#Users) others are machines.
Nonetheless we treat them all the same in regard to [roles](administrate#Roles) management and audit trail.

#### Human vs. Service Users

The major difference between humane vs. machine [users](administrate#Users) is the type of credentials who can be used.
With machine [users](administrate#Users) there is only a non interactive login process possible. As such we utilize “JWT as Authorization Grant”.

> TODO Link to “JWT as Authorization Grant” explanation.

### How ZITADEL handles usernames

**ZITADEL** is built around the concept of [organisations](administrate#Organisations). Each [organisation](administrate#Organisations) has it's own pool of usernames which include human and service [users](administrate#Users).
For example a [user](administrate#Users) with the username `road.runner` can only exist once the [organisation](administrate#Organisations) `ACME`. **ZITADEL** will automatically generate a "logonname" for each [user](administrate#Users) consisting of `{username}@{domainname}.{zitadeldomain}`. Without verifying the domain name this would result in the logonname `road.runner@acme.zitadel.ch`. If you use a dedicated **ZITADEL** replace `zitadel.ch` with your domain name.

If someone verifies a domain name within the organisation **ZITADEL** will generate additional logonames for each [user](administrate#Users) with that domain. For example if the domain is `acme.ch` the resulting logonname would be `road.runner@acme.ch` and as well the generated one `road.runner@acme.zitadel.ch`.

> Domain verification also removes the logonname from all [users](administrate#Users who might have used this combination in the global [organisation](administrate#Organisations).
> Relating to example with `acme.ch` if a user in the global [organisation](administrate#Organisations), let's call him `coyote` used `coyote@acme.ch` this logonname will be replaced with `coyote@randomvalue.tld`
> **ZITADEL** notifies the user about this change

### Manage Users

#### Search User

<img src="img/console_user_list_search.png" alt="User list Search" width="1000px" height="auto">

Image 1: User List Search

#### Create User

<img src="img/console_user_list.png" alt="User list" width="1000px" height="auto">

Image 2: User List

<img src="img/console_user_create_form.png" alt="User Create Form" width="1000px" height="auto">

Image 3: User Create Form

<img src="img/console_user_create_done.png" alt="User Create Done" width="1000px" height="auto">

Image 4: User Create Done

#### Set Password

> Screenshot here

### Manage Service Users

> Screenshot here

### Authorizations

> Screenshot here

### Audit user changes

> Screenshot here

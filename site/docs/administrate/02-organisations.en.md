---
title: Organisations
---

### What are organisations

Organisations are comparable to tenants of a system or OU's (organisational units) if we speak of a directory based system.
ZITADEL is organised around the idea that multiple organisations share the same [System](#What_is_meant_by_system) and that they can grant each other rights so self manage certain things.

### Create an organisation without existing login

ZITADEL allows you to create a new organisation without a preexisting user. For [ZITADEL.ch](https://zitadel.ch) you can create a org by visiting the [Register organisation](https://accounts.zitadel.ch/register/org)

> Screenshot here

For dedicated ZITADEL instances this url might be different, but in most cases should be something like https://accounts.YOURDOMAIN.TLD/register/org

### Create an organisation with existing login

You can simply create a new organisation by visiting the [ZITADEL Console](https://console.zitadel.ch) and clicking "new organisation" in the upper left corner.

> Screenshot here

For dedicated ZITADEL instances this url might be different, but in most cases should be something like https://console.YOURDOMAIN.TLD

### Verify a domain name

Once you created your organisation you will receive a generated domain name from ZITADEL for your organisation. For example if you call your organisation "ACME" you will receive "acme.zitadel.ch" as name. Furthermore the users you create will be suffixed with this domain name. To improve the user experience you can verify a domain name which you control. If you control acme.ch you can verify the ownership by DNS or HTTP challenge.
After the domain is verified your users can use both domain names to log-in. The user "coyote" can now use "coyote@acme.zitadel.ch" and "coyote@acme.ch".
An organisation can have multiple domain names, but only one of it can be primary. The primary domain defines which loginname ZITADEL displays to the user, and also what information gets asserted in access_tokens (preferred_username).

> Screenshot here

### Audit organisation changes

All changes to the organisation are displayed on the organisation menu within [ZITADEL Console](https://console.zitadel.ch/org) organisation menu. Located on the right hand side under "activity"

> Screenshot here

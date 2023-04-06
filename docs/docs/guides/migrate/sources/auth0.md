---
title: Migrate from Auth0
sidebar_label: From Auth0
---

Migrating users from Auth0 to ZITADEL requires the following steps:

- Request and download hashed passwords
- Export all user data
- Import users and password hashes to ZITADEL

## Export hashed passwords

Auth0 does not export hashed passwords as part of the bulk user export.
You must create a support ticket to download password hashes and password-related information.
Please also refer to the Auth0 guide on how to [Export Data](https://auth0.com/docs/troubleshoot/customer-support/manage-subscriptions/export-data#user-passwords).

:::info
You can also import users into ZITADEL with an verified email but without the passwords.
Users will be prompted to create a new password after they login for the first time after migration.
:::

1. Go to https://support.auth0.com/tickets and click on **Open Ticket**
2. Issue Type: **I have a question regarding my Auth0 account**
3. What can we help you with?: **I would like to obtain an export of my tenant password hases**
4. Fill out the form: Request password hashes as bcrypt
5. **Submit ticket**

You will receive a JSON file including the password hashes.
See the this [community post](https://community.auth0.com/t/password-hashes-export-data-format/58730) for more information about the contents and format.

## Export all user data

Create a [bulk user export](https://auth0.com/docs/manage-users/user-migration/bulk-user-exports) from the Auth0 Management API.
You will receive a newline-delimited JSON with the requested user data.

## Import users and password hashes to ZITADEL

You will need to merge the received password hashes with the user bulk export.

After you successfully merged the datasets, you can follow the instructions described in the [Migrate Users](../users) guide to import users to ZITADEL.

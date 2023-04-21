---
title: Migrate from Auth0
sidebar_label: From Auth0
---

Migrating users from Auth0 to ZITADEL requires the following steps:

- Request and download hashed passwords
- Export all user data
- Run migration tool to merge Auth0 users and passwords
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
3. What can we help you with?: **I would like to obtain an export of my tenant password hashes**
4. Fill out the form: Request password hashes as bcrypt
5. **Submit ticket**

You will receive a JSON file including the password hashes.
See this [community post](https://community.auth0.com/t/password-hashes-export-data-format/58730) for more information about the contents and format.

## Export all user data

Create a [bulk user export](https://auth0.com/docs/manage-users/user-migration/bulk-user-exports) from the Auth0 Management API.
You will receive a newline-delimited JSON with the requested user data.

This is an example request, we have included the user id, the email and the name of the user. Make sure to export the users in a json format.

```bash
curl --request POST \
  --url $AUTH0_DOMAIN/api/v2/jobs/users-exports \
  --header 'authorization: Bearer $TOKEN' \
  --header 'content-type: application/json' \
  --data '{
	"connection_id": "$CONNECTION_ID",
	"format": "json", 
	"fields": [
		{"name": "user_id"},
		{"name": "email"},
		{"name": "name"},
	]
}'
```

## Run Migration Tool

We have developed a tool that combines your exported user data with their corresponding passwords to generate the import request body for ZITADEL.

1. Download the latest release of [github.com/zitadel/zitadel-tools](https://github.com/zitadel/zitadel-tools/releases)
2. Execute the binary with the following flags:
 ```bash
 ./zitadel-tools migrate auth0 --org=<organisation id> --users=./users.json --passwords=./passwords.json --output=./importBody.json
 ```
 Use the Organization ID from your ZITADEL instance where you like to add the users.
3. You will now get a new file importBody.json which contains the body for the request to the import of ZITADEL

## Import users and password hashes to ZITADEL

Copy the content from the importBody.json file created in the last step.
You can now follow the instructions described in the [Migrate Users](../users) guide to import users to ZITADEL.

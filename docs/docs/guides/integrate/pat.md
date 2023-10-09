---
title: ZITADEL's Personal Access Tokens(PATs)
sidebar_label: Personal Access Tokens(PATs)
---


A Personal Access Token (PAT) is a ready to use token which can be used as _Authorization_ header.
At the moment ZITADEL only allows PATs for machine accounts (service users).

It is an alternative to the JWT profile authentication where the service user has a key to authenticate. Read more about that [here](serviceusers)

## Create a Service User with a PAT


1. Navigate to Service Users
2. Click on **New**
3. Enter a user name and a display name
4. Click on the Personal Access Token menu point in the detail of your user
5. Click on **New**
6. You can either set an expiration date or leave it empty if you don't want it to expire
7. Copy the token from the dialog (You will not see this again)

![Create new service user](/img/guides/console-service-user-pat.gif)

## Grant role for ZITADEL

To be able to access the ZITADEL APIs your service user needs permissions to ZITADEL.

1. Go to the detail page of your organization
2. Click in the top right corner the "+" button
3. Search for your service user
4. Give the user the role you need, for the example we choose Org Owner (More about [ZITADEL Permissions](../manage/console/managers))

![Add org owner to service user](/img/guides/console-service-user-org-owner.gif)


## Call ZITADEL API with PAT

Because the PAT is a ready to use Token, you can add it as Authorization Header and send it in your requests to the ZITADEL API.
In this example we read the organization of the service user.

```bash
curl --request GET \
  --url $CUSTOM-DOMAIN/management/v1/orgs/me \
  --header 'Authorization: Bearer {PAT}' 
```
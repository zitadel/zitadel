---
title: Client Credentials with Service Users
---

This is a guide on how to use Client Credentials with service users in ZITADEL. You can read more about users [here](/concepts/structure/users.md).

In ZITADEL, the Client Credentials grant can be used for this non-interactive authentication as alternative to the [JWT profile authentication](serviceusers).

## Create a Service User with a Secret

1. Navigate to Service Users
2. Click on **New**
3. Enter a username and a display name
4. Click on **Create**
5. Open **Actions** in the top right corner and click on **Generate Client Secret**
6. Copy the **ClientID** and **ClientSecret** from the dialog

:::note
Be sure to copy in particular the ClientSecret. You won't be able to retrieve it again.
If you lose it, you will have to generate a new one.
:::

![Create new service user](/img/console_serviceusers_secret.gif)

## Grant role for ZITADEL

To be able to access the ZITADEL APIs your service user needs permissions to ZITADEL.

1. Go to the detail page of your organization
2. Click in the top right corner the "+" button
3. Search for your service user
4. Give the user the role you need, for the example we choose Org Owner (More about [ZITADEL Permissions](/guides/manage/console/managers))

![Add org owner to service user](/img/guides/console-service-user-org-owner.gif)

## Authenticating a service user

In this step we will authenticate a service user and receive an access_token to use against the ZITADEL API.

You will need to craft a POST request to ZITADEL's token endpoint:

```bash
curl --request POST \
  --url https://{your_domain}.zitadel.cloud/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header 'Authorization: Basic ${BASIC_AUTH}' \
  --data grant_type=client_credentials \
  --data scope='openid profile email urn:zitadel:iam:org:project:id:zitadel:aud'
```

* `grant_type` should be set to `client_credentials`
* `scope` should contain any [Scopes](/apis/openidoauth/scopes) you want to include, but must include `openid`. For this example, please include `profile`, `email`
  and `urn:zitadel:iam:org:project:id:zitadel:aud`. The latter provides access to the ZITADEL API.

You should receive a successful response with `access_token`,  `token_type` and time to expiry in seconds as `expires_in`.

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{
  "access_token": "MtjHodGy4zxKylDOhg6kW90WeEQs2q...",
  "token_type": "Bearer",
  "expires_in": 43199
}
```

## Call ZITADEL API with Token

Because the received Token includes the `urn:zitadel:iam:org:project:id:zitadel:aud` scope, we can send it in your requests to the ZITADEL API as Authorization Header.
In this example we read the organization of the service user.

```bash
curl --request GET \
  --url {your-domain}/management/v1/orgs/me \
  --header 'Authorization: Bearer ${TOKEN}' 
```

## Summary

* With service users you can secure machine-to-machine communication
* Client Credentials provide an alternative way to JWT Profile for service user authentication
* After successful authorization you can use an access token like for human users

Where to go from here:

* Management API
* Securing backend API

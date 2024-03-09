---
title: Configure client credential authentication for service users
sidebar_label: Client credential authentication
sidebar_position: 3
---

This guide demonstrates how developers can leverage Client Credential authentication to secure communication between [service users](/concepts/structure/users) and client applications within ZITADEL.

In ZITADEL, the Client Credentials grant can be used for this [non-interactive authentication](authenticate-service-users) as alternative to the [JWT profile authentication](serviceusers).

## Create a Service User with a client secret

1. Navigate to Service Users
2. Click on **New**
3. Enter a username and a display name
4. Click on **Create**
5. Open **Actions** in the top right corner and click on **Generate Client Secret**
6. Copy the **ClientID** and **ClientSecret** from the dialog

:::note
Make sure to copy in particular the ClientSecret. You won't be able to retrieve it again.
If you lose it, you will have to generate a new one.
:::

![Create new service user](/img/console_serviceusers_secret.gif)

## Grant a manager role to the service user

## Authenticating a service user

In this step we will authenticate a service user and receive an access_token to use against the ZITADEL API.

You will need to craft a POST request to ZITADEL's token endpoint:

```bash
curl --request POST \
  --url https://$CUSTOM-DOMAIN/oauth/v2/token \
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


## Accessing ZITADEL's Management API with client credentials




## Call ZITADEL API with Token

Because the received Token includes the `urn:zitadel:iam:org:project:id:zitadel:aud` scope, we can send it in your requests to the ZITADEL API as Authorization Header.
In this example we read the organization of the service user.

```bash
curl --request GET \
  --url $CUSTOM-DOMAIN/management/v1/orgs/me \
  --header 'Authorization: Bearer ${TOKEN}' 
```

## Summary

* With service users you can secure machine-to-machine communication
* Client Credentials provide an alternative way to JWT Profile for service user authentication
* After successful authorization you can use an access token like for human users

## Notes

* Read about the [different methods to authenticate service users](./authenticate-service-users)
* [Service User API reference](/docs/category/apis/resources/mgmt/user-machine)
* [OIDC client secret basic](/docs/apis/openidoauth/authn-methods#client-secret-basic) authentication method reference
* [Access ZITADEL APIs](../zitadel-apis/)
* Validate access tokens with [token introspection with basic auth](../token-introspection/basic-auth)

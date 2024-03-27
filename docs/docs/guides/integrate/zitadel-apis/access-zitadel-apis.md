---
title: Access ZITADEL APIs
---

This guide explains what ZITADEL APIs are and how to access ZITADEL APIs using a service user to manage all types of resources and settings.

## What are the ZITADEL APIs

ZITADEL exposes a variety of APIs that allow you to interact with its functionalities programmatically.
These APIs are offered through different protocols including gRPC and REST.
Additionally, ZITADEL provides [SDKs for popular languages](/docs/sdk-examples/introduction) and frameworks to simplify integration.

Here's a breakdown of some key points about ZITADEL APIs:

* **Auth API:** Used by authenticated users for tasks related to their accounts.
* **Management API:** Used by organization managers for administrative tasks.
* **Admin API:** Used for administrative functions on the ZITADEL instance itself (may require separate user setup).  
* **System API:** For ZITADEL self-hosted deployments only, providing superordinate control (requires specific configuration).

:::note Resource-based APIs
ZITADEL is transitioning from a use-case based API structure to a resource-based one, aiming to simplify API usage.
:::

For further details and in-depth exploration, you can refer to the [ZITADEL API documentation](/docs/apis/introduction).

## How to access ZITADEL APIs

Accessing ZITADEL APIs, except for the Auth API and the System API, requires these basic steps:

1. **Create a service user**: A service user is a special type of account used to grant programmatic access to ZITADEL's functionalities. Unlike regular users who log in with a username and password, [service users rely on a more secure mechanism involving digital keys and tokens](../service-users/authenticate-service-users).
2. **Give permission to access ZITADEL APIs**: Assign a Manager role to the service user, giving it permission to make changes to certain resources in ZITADEL.
3. **Authenticate the service user**: Like human users, service users must authenticate and request an OAuth token with the scope `urn:zitadel:iam:org:project:id:zitadel:aud` to access ZITADEL APIs. [Service users can be authenticated](../service-users/authenticate-service-users) using private key JWT, client credentials, or personal access tokens.
4. **Access ZITADEL APIs with the token**: The OAuth token must be included in the Authorization Header of calls to ZITADEL APIs.

### Accessing Auth API

The Auth API can be used for all operations on the requesting user, meaning the user id in the sub claim of the used token.
Using this API doesn't require a service user to be authenticated.
Instead, you call the Auth API with the token of the user.

[Reference documentation for authentication API](/docs/apis/introduction#authentication)

### Accessing System API

With the System API developers can manage different ZITADEL instances.
The System API can't be accessed by service users and requires a special configuration and authentication that can be found in our [guide to access ZITADEL's System API](./access-zitadel-system-api).

[Reference documentation for system API](/docs/apis/introduction#system)

## 1. Create a service user

First, you need to create a new service user through the console or ZITADEL APIs.

Via Console:

1. In an organization, navigate to Users > Service Users
2. Click on **New**
3. Enter a username and a display name
4. Click on **Create**

Via APIs:

* [Create User (Machine)](/docs/apis/resources/mgmt/management-service-add-machine-user)

## 2. Grant a Manager role to the service user

ZITADEL Managers are Users who have permission to manage ZITADEL itself.
There are some different levels for managers.

- **IAM Managers**: This is the highest level. Users with IAM Manager roles are able to manage the whole instance.
- **Org Managers**: Managers in the Organization Level are able to manage everything within the granted Organization.
- **Project Managers**: At this level, the user is able to manage a project.
- **Project Grant Manager**: The project grant manager is for projects, which are granted of another organization.

On each level, we have some different Roles. Here you can find more about the different roles: [ZITADEL Manager Roles](/guides/manage/console/managers#roles)

To be able to access the ZITADEL APIs your service user needs permissions to ZITADEL.

1. Go to the detail page of your organization
2. Click in the top right corner the "+" button
3. Search for your service user
4. Give the user the role you need, for the example we choose Org Owner (More about [ZITADEL Permissions](/guides/manage/console/managers))

![Add Org Manager](/img/console_org_manager_add.gif)

## 3. Authenticate service user and request token

Service users can be authenticated using private key JWT, client credentials, or personal access tokens.
The [service user authentication](../service-users/authenticate-service-users) can be used to make machine-to-machine requests to any Resource Server (eg, a backend service / API) by requesting a token from the Authorization Server (ZITADEL) and sending the short-lived token (access token) in the Header of requests.

This guide covers a specific case of service user authentication when requesting access to the [ZITADEL APIs](/docs/apis/introduction).
While PAT can be used directly to access the ZITADEL APIS, the more secure authentication methods private key JWT and client credentials must include the [reserved scope](/docs/apis/openidoauth/scopes) `urn:zitadel:iam:org:project:id:zitadel:aud` when requesting an access from the token endpoint.
This scope will add the ZITADEL APIs to the audience of the access token.
ZITADEL APIs will check if they are in the audience of the access token, and reject the token in case they are not in the audience.

The following sections will explain the more specific authentication to access the ZITADEL APIs.

### Authenticate with private key JWT

Follow the steps in this guide to [generate an key file](../service-users/private-key-jwt#2-generate-a-private-key-file) and [create a JWT and sign with private key](../service-users/private-key-jwt#3-create-a-jwt-and-sign-with-private-key).

With the encoded JWT (assertion) from the prior step, you will need to craft a POST request to ZITADEL's token endpoint.

**To access the ZITADEL APIs you need the ZITADEL Project ID in the audience of your token.**
This is possible by sending a [reserved scope](/apis/openidoauth/scopes) for the audience.
Use the scope `urn:zitadel:iam:org:project:id:zitadel:aud` to include the ZITADEL project id in your audience

A sample request will look like this

```bash {5}
curl --request POST \
  --url $CUSTOM-DOMAIN/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data scope='openid profile urn:zitadel:iam:org:project:id:zitadel:aud' \
  --data assertion=eyJ0eXAiOiJKV1QiL...
```

where 

- `grant_type` must be set to `urn:ietf:params:oauth:grant-type:jwt-bearer`
- `scope` should contain any [Scopes](/apis/openidoauth/scopes) you want to include, but must include `openid` and `urn:zitadel:iam:org:project:id:zitadel:aud` to acces the ZITADEL APIs. For this example include `profile`.
- `assertion` is the encoded value of the JWT that was signed with your private key from the prior step

You should receive a successful response with `access_token`, `token_type` and time to expiry in seconds as `expires_in`.

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{
  "access_token": "MtjHodGy4zxKylDOhg6kW90WeEQs2q...",
  "token_type": "Bearer",
  "expires_in": 43199
}
```

Use the access_token in the Authorization header to make requests to the ZITADEL APIs.
In the following example, we read the organization of the service user.

```bash
curl --request GET \
  --url $CUSTOM-DOMAIN/management/v1/orgs/me \
  --header 'Authorization: Bearer ${TOKEN}' 
```

### Client credentials

Get the client id and client secret by

1. navigating to the service user, then
2. Open **Actions** in the top right corner and click on **Generate Client Secret**
3. Copy the **ClientID** and **ClientSecret** from the dialog

![Create new service user](/img/console_serviceusers_secret.gif)

With the ClientId and ClientSecret from the prior step, you will need to craft a POST request to ZITADEL's token endpoint.

#### Audience scope

**To access the ZITADEL APIs you need the ZITADEL Project ID in the audience of your token.**
This is possible by sending a [reserved scope](/apis/openidoauth/scopes) for the audience.
Use the scope `urn:zitadel:iam:org:project:id:zitadel:aud` to include the ZITADEL project id in your audience

In this step we will authenticate a service user and receive an access_token to use against the ZITADEL API.

#### Basic authentication

When using `client_secret_basic` on the token or introspection endpoints, provide an `Authorization` header with a Basic auth value in the following form:

```markdown
Authorization: "Basic " + base64( formUrlEncode(client_id) + ":" + formUrlEncode(client_secret) )
```

For an example see the [client secret basic authentication method reference](/docs/apis/openidoauth/authn-methods#client-secret-basic).
We recommend using an OpenID / OAuth library that handles the encoding for you.

#### Post request

You will need to craft a POST request to ZITADEL's token endpoint:

```bash {6}
curl --request POST \
  --url https://$CUSTOM-DOMAIN/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --header 'Authorization: Basic ${BASIC_AUTH}' \
  --data grant_type=client_credentials \
  --data scope='openid profile urn:zitadel:iam:org:project:id:zitadel:aud'
```

where

* `grant_type` should be set to `client_credentials`
* `scope` should contain any [Scopes](/apis/openidoauth/scopes) you want to include, but must include `openid`. For this example, please include `profile`
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

Because the received token includes the `urn:zitadel:iam:org:project:id:zitadel:aud` scope, we can send it in your requests to the ZITADEL API as the Authorization Header.
In this example we read the organization of the service user.

```bash
curl --request GET \
  --url $CUSTOM-DOMAIN/management/v1/orgs/me \
  --header 'Authorization: Bearer ${TOKEN}' 
```

### Personal access token (PAT)

A Personal Access Token (PAT) is a ready-to-use token which can be used as _Authorization_ header.

Because the PAT is a ready-to-use token, you can add it as Authorization Header and send it in your requests to the ZITADEL API.
In this example, we read the organization of the service user.

```bash
curl --request GET \
  --url $CUSTOM-DOMAIN/management/v1/orgs/me \
  --header 'Authorization: Bearer {PAT}' 
```

## Notes

- [Example application in Go](./example-zitadel-api-with-go) to access ZITADEL APIs
- [Example application in .NET](./example-zitadel-api-with-dot-net) to access ZITADEL APIs
- Learn how to use the [Event API](./event-api) to retrieve your audit trail
- Read about the [different methods to authenticate service users](../service-users/authenticate-service-users)

---
title: Access ZITADEL APIs
---

This guide explains what ZITADEL APIs are and how to access ZITADEL APIs using a service user to manage all types of resources and settings.

## What are the ZITADEL APIs

ZITADEL exposes a variety of APIs that allow you to interact with its functionalities programmatically.
These APIs are offered through different protocols including gRPC and REST.
Additionally, ZITADEL provides [SDKs for popular languages](../../../sdk-examples/) and frameworks to simplify integration.

Here's a breakdown of some key points about ZITADEL APIs:

* **Auth API:** Used by authenticated users for tasks related to their accounts.
* **Management API:** Used by organization managers for administrative tasks.
* **Admin API:** Used for administrative functions on the ZITADEL instance itself (may require separate user setup).  
* **System API:** For ZITADEL self-hosted deployments only, providing superordinate control (requires specific configuration).

:::note Migration
ZITADEL is transitioning from a use-case based API structure to a resource-based one, aiming to simplify API usage.
:::

For further details and in-depth exploration, you can refer to the [Zitadel API documentation](/docs/apis/introduction).

## How to access ZITADEL APIs

Accessing ZITADEL APIs, except for the Auth API and the System API, requires these basic steps:

1. **Create a service user**: A service user is a special type of account used to grant programmatic access to ZITADEL's functionalities. Unlike regular users who log in with a username and password, [service users rely on a more secure mechanism involving digital keys and tokens](../service-users/authenticate-service-users).
2. **Give permission to access ZITADEL APIs**: Assign a Manager role to the service  user, giving it permission to make changes to certain resources in ZITADEL.
3. **Authenticate the service user**: Like human users, service users must authenticate and request a OAuth token with the scope `urn:zitadel:iam:org:project:id:zitadel:aud` to access ZITADEL APIs.
4. **Access ZITADEL APIs with the token**: The OAuth token must be included in the Authorization Header of calls to ZITADEL APIs.

### Auth API

The Auth API can be used for all operations on the requesting user, meaning the user id in the sub claim of the used token.
Using this API doesn't require a service user to be authenticated.
Instead you call the Auth API with the token of the user.

### System API

With the System API developers can manage different ZITADEL instances.
The System API can't be accessed with service users and requires a special configuration and authentication that can be found in our [guide to access ZITADEL's System API](./access-zitadel-system-api).

## ZITADEL Managers

ZITADEL Managers are Users who have permission to manage ZITADEL itself. There are some different levels for managers.

- **IAM Managers**: This is the highest level. Users with IAM Manager roles are able to manage the whole instance.
- **Org Managers**: Managers in the Organization Level are able to manage everything within the granted Organization.
- **Project Mangers**: In this level the user is able to manage a project.
- **Project Grant Manager**: The project grant manager is for projects, which are granted of another organization.

On each level we have some different Roles. Here you can find more about the different roles: [ZITADEL Manager Roles](/guides/manage/console/managers#roles)

## Add ORG_OWNER to Service User

Make sure you have a Service User with a Key. (For more detailed informations about creating a service user go to [Service User](serviceusers))

1. Navigate to Organization Detail
2. Click the **+** button in the right part of console, in the managers part of details
3. Search the user and select it
4. Choose the role ORG_OWNER

![Add Org Manager](/img/console_org_manager_add.gif)

## Authenticating a service user

In ZITADEL we use the `urn:ietf:params:oauth:grant-type:jwt-bearer` (**“JWT bearer token with private key”**, [RFC7523](https://tools.ietf.org/html/rfc7523)) authorization grant for this non-interactive authentication.
This is already described in the [Service User](./serviceusers), so make sure you follow this guide.

### Request an OAuth token, with audience for ZITADEL

With the encoded JWT from the prior step, you will need to craft a POST request to ZITADEL's token endpoint:

To access the ZITADEL APIs you need the ZITADEL Project ID in the audience of your token.
This is possible by sending a custom scope for the audience. More about [Custom Scopes](/apis/openidoauth/scopes)

Use the scope `urn:zitadel:iam:org:project:id:zitadel:aud` to include the ZITADEL project id in your audience

```bash
curl --request POST \
  --url $CUSTOM-DOMAIN/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data scope='openid profile email urn:zitadel:iam:org:project:id:zitadel:aud' \
  --data assertion=eyJ0eXAiOiJKV1QiL...
```

- `grant_type` must be set to `urn:ietf:params:oauth:grant-type:jwt-bearer`
- `scope` should contain any [Scopes](/apis/openidoauth/scopes) you want to include, but must include `openid`. For this example, please include `profile` and `email`
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

With this token you are allowed to access the [ZITADEL APIs](/apis/introduction) .

## Summary

- Grant a user for ZITADEL
- Because there is no interactive logon, you need to use a JWT signed with your private key to authorize the user
- With a custom scope (`urn:zitadel:iam:org:project:id:zitadel:aud`) you can access ZITADEL APIs

Where to go from here:

- [ZITADEL API Documentation](/apis/introduction)

## Notes

- Example Go
- Example DotNet


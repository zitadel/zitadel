---
title: Access ZITADEL APIs
---

:::note
This guide focuses on the Admin, Auth and Management APIs. To access the ZITADEL System API, please checkout [this guide](./access-zitadel-system-api).
:::

## ZITADEL Managers

ZITADEL Managers are Users who have permission to manage ZITADEL itself. There are some different levels for managers.

- **IAM Managers**: This is the highest level. Users with IAM Manager roles are able to manage the whole instance.
- **Org Managers**: Managers in the Organization Level are able to manage everything within the granted Organization.
- **Project Mangers**: In this level the user is able to manage a project.
- **Project Grant Manager**: The project grant manager is for projects, which are granted of another organization.

On each level we have some different Roles. Here you can find more about the different roles: [ZITADEL Manager Roles](/guides/manage/console/managers#roles)

## Add ORG_OWNER to Service User

Make sure you have a Service User with a Key. (For more detailed informations about creating a service user go to [Service User](serviceusers.md))

1. Navigate to Organization Detail
2. Click the **+** button in the right part of console, in the managers part of details
3. Search the user and select it
4. Choose the role ORG_OWNER

![Add Org Manager](/img/console_org_manager_add.gif)

## Authenticating a service user

In ZITADEL we use the `urn:ietf:params:oauth:grant-type:jwt-bearer` (**“JWT bearer token with private key”**, [RFC7523](https://tools.ietf.org/html/rfc7523)) authorization grant for this non-interactive authentication.
This is already described in the [Service User](serviceusers.md), so make sure you follow this guide.

### Request an OAuth token, with audience for ZITADEL

With the encoded JWT from the prior step, you will need to craft a POST request to ZITADEL's token endpoint:

To access the ZITADEL APIs you need the ZITADEL Project ID in the audience of your token.
This is possible by sending a custom scope for the audience. More about [Custom Scopes](/apis/openidoauth/scopes)

Use the scope `urn:zitadel:iam:org:project:id:zitadel:aud` to include the ZITADEL project id in your audience

```bash
curl --request POST \
  --url {your_domain}/oauth/v2/token \
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

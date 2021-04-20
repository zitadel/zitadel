---
title: Access ZITADEL APIs
---

<table class="table-wrapper">
    <tr>
        <td>Description</td>
        <td>Learn how to authorize Service Users to access ZITADEL APIs.</td>
    </tr>
    <tr>
        <td>Learning Outcomes</td>
        <td>
            In this module you will:
            <ul>
                <li>Grant a Service User for ZITADEL</li>
                <li>Authorize a Service User with JWT signed with your private key</li>
            </ul>
        </td>
    </tr>
     <tr>
        <td>Prerequisites</td>
        <td>
            <ul>
                <li>Knowledge of <a href="/docs/guides/oauth-recommended-flows">Recommended Authorization Flows</a></li>
                <li>Knowledge of <a href="/docs/guides/serviceusers">Service Users</a></li>
            </ul>
        </td>
    </tr>
</table>

## ZITADEL Managers

ZITADEL Managers are Users who have permission to manage ZITADEL itself. There are some different levels for managers. 

- **IAM Managers**: This is the highest level. Users with IAM Manager roles are able to manage the whole IAM. 
- **Org Managers**: Managers in the Organisation Level are able to manage everything within the granted Organisation.
- **Project Mangers**: In this level the user is able to manage a project.
- **Project Grant Manager**: The project grant manager is for projects, which are granted of another organisation.

On each level we have some different Roles. Here you can find more about the different roles: [ZITADEL Manager Roles](../manuals/admin-managers)


## Exercise: Add ORG_OWNER to Service User

Make sure you have a Service User with a Key. (For more detailed informations about creating a service user go to [Service User](serviceusers))

1. Navigate to Organisation Detail
2. Click the **+** button in the right part of console, in the managers part of details
3. Search the user and select it
4. Choose the role ORG_OWNER

![Add Org Manager](/img/console_org_manager_add.gif)

## Authenticating a service user

In ZITADEL we use the `private_jwt` (**“JWT bearer token with private key”**, [RFC7523](https://tools.ietf.org/html/rfc7523)) authorization grant for this non-interactive authentication.
This is already described in the [Service User](serviceusers), so make sure you follow this guide.

### Request an OAuth token, with audience for ZITADEL

With the encoded JWT from the prior step, you will need to craft a POST request to ZITADEL's token endpoint:

To access the ZITADEL APIs you need the ZITADEL Project ID in the audience of your token.
This is possible by sending a custom scope for the audience. More about [Custom Scopes](../apis/openidoauth/scopes)

Use the scope `urn:zitadel:iam:org:project:id:{projectid}:aud` to include the project id in your audience

> The scope for zitadel.ch is: `urn:zitadel:iam:org:project:id:69234237810729019:aud`

```bash
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer \
  --data scope='openid profile email urn:zitadel:iam:org:project:id:69234237810729019:aud' \
  --data assertion=eyJ0eXAiOiJKV1QiL...
```

* `grant_type` must be set to `urn:ietf:params:oauth:client-assertion-type:jwt-bearer`
* `scope` should contain any [Scopes](../apis/openidoauth/scopes) you want to include, but must include `openid`. For this example, please include `profile` and `email`
* `assertion` is the encoded value of the JWT that was signed with your private key from the prior step

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

With this token you are allowed to access the [ZITADEL APIs](../apis/introduction) .
## Knowledge Check


* Managers are used for all users to get authorizations for all projects?
    - [ ] yes
    - [ ] no
* You need the ZITADEL Project Id in your audience and can access this with a custom scope?
    - [ ] yes
    - [ ] no


<details>
    <summary>
        Solutions
    </summary>


* Managers are used for all users to get authorizations for all projects?
    - [ ] yes
    - [x] no (Managers are only used to grant users for ZITADEL)
* You need the ZITADEL Project Id in your audience and can access this with a custom scope?
    - [x] yes
    - [ ] no

</details>

## Summary

* Grant a user for ZITADEL
* Because there is no interactive logon, you need to use a JWT signed with your private key to authorize the user
* With a custom scope (`urn:zitadel:iam:org:project:id:{projectid}:aud`) you can access ZITADEL APIs


Where to go from here:

* [ZITADEL API Documentation](../apis/introduction)

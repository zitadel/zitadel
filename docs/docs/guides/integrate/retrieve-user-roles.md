---
title: Retrieve user roles
---

This guide explains all the possible ways of retrieving user roles across different organizations and projects using ZITADEL's APIs. 

## What are roles/authorizations/grants in ZITADEL? 
User roles, user grants, or authorizations refer to the roles that are assigned to a user. These terms are used interchangeably to mean the roles assigned to the user, e.g., the ZITADEL Console refers to the pairing of roles and users as authorizations, whereas the APIs refer to them as grants. This guide will use the term roles for application-specific roles (e.g., `admin`, `accountant`, `employee`, `hr`, etc.) and ZITADEL-specific manager roles (e.g., `IAM_OWNER`, `ORG_OWNER`, `PROJECT_OWNER`, etc.). 

Roles are critical to managing permissions in a single-tenant or multi-tenant application. It can, however, be tricky to retrieve them, especially when spanning multiple organizations and projects. 

## Assign roles and memberships

Human or service users can be assigned roles. You can do this via the ZITADEL Console or the ZITADEL APIs. As mentioned earlier, there are two types of roles in ZITADEL. You can have your own application-specific roles, and alternatively, ZITADEL also has manager roles, such as `ORG_OWNER` and `IAM_OWNER`. 

Follow the links below to assign roles to your users. 

- [Add application roles via the ZITADEL Console](/docs/guides/manage/console/roles)
- [Add manager roles via the ZITADEL Console](/docs/guides/manage/console/managers)
- [Add application roles via the ZITADEL Management API](/docs/category/apis/resources/mgmt/project-roles)
- [Add manager roles to users via the ZITADEL Management API](/category/apis/resources/mgmt/members)

## Retrieve roles

Roles can be requested via our auth and management APIs, from userinfo endpoint or ID token. Currently, manager roles cannot be directly included in the token. You will need to use the ZITADEL APIs to retrieve them.

### Generate a token

You must first of all generate a token for the user. If it’s a human user, he would be using a front-end application and logging in via the browser or device. An access token will be returned after they log in successfully. A machine user will use a script or other program to generate a token using the JWT profile or client credentials grant types. 

How to generate a token: 

- [Generate tokens for human users](/docs/guides/integrate/login-users)
- [Generate tokens for service users](/docs/guides/integrate/serviceusers)

In order to access role information via the token you must include the right audience and the necessary role claims in the scope and/or select the required role settings in the ZITADEL console before requesting the token. 

### Determine the audience

An important concept in OpenID Connect (OIDC) is the 'audience' (`aud`) claim, which is part of the token payload. The `aud` claim identifies who or what this token is intended for. If the recipient (e.g., a resource server) does not identify itself with a value in the `aud` claim when this claim is present, then the token must be rejected. 

The audience is essential in multi-tier systems, where you may authenticate with one client (in one project) but access resources from another client (in a different project) or when you are accessing ZITADEL’s management APIs. Without the correct audience in your token, you will run into errors, such as the ‘Invalid audience’ error in ZITADEL. 

You can determine the audience in two ways:

**1. Use the explicit scope for ZITADEL to access only ZITADEL APIs:**

If your application needs to access ZITADEL's APIs (for example, to pull a list of all users), follow this steps:

- Add the scope `urn:zitadel:iam:org:project:id:zitadel:aud` to the authentication request when authenticating the user. Now, the application can make calls to ZITADEL's API without running into an ‘Invalid audience’ error.

**2. Include the project id of the ZITADEL project in the scope for accessing the ZITADEL APIs or anything else:**

Let's assume you have a frontend application and a backend application under different projects. Here's how to add the correct audience:

- Authenticate the end-users to an application in your front-end project.
- In the authentication request sent from the front-end application, add the scope `urn:zitadel:iam:org:project:id:{projectId}:aud`, replacing `{projectId}` with the project ID of your backend application.
- Now, the front end can send requests to the backend. The backend will validate the token with ZITADEL's introspection endpoint and will not return an ‘Invalid audience’ error.

And you can also use the same to access the ZITADEL APIs. 

### Role settings in the ZITADEL Console 

If you need user roles returned from the userinfo endpoint, you must select the **’Assert Roles on Authentication’** checkbox in your project under general settings. 

<img
    src="/docs/img/guides/integrate/retrieve-user-roles-1.png"
    width="75%"
    alt="Assert Roles on Authentication"
/>

If you need them included in your ID Token, select **’User Roles Inside ID Token’** in application settings. This has to be set in your applications as this is dependent on your application type. Navigate to your application and select this setting.

<img
    src="/docs/img/guides/integrate/retrieve-user-roles-2.png"
    width="75%"
    alt="Assert Roles on Authentication"
/>


Alternatively, you can include the claims `urn:iam:org:project:roles` or/and `urn:zitadel:iam:org:projects:roles` in your scope to achieve the same as above. 

### Retrieve roles from the userinfo endpoint

The user info endpoint is  **ZITADEL_DOMAIN/oidc/v1/userinfo**.

This endpoint will return information about the authenticated user.
Send the access token of the user as `Bearer Token` in the `Authorization` header:

**cURL Request:**
```bash
curl --request GET \
 --url $ZITADEL_DOMAIN/oidc/v1/userinfo
 --header 'Authorization: Bearer <TOKEN>'
```

If the access token is valid, the information about the user (depending on the granted scopes) is returned. Check the [Claims page](/docs/apis/openidoauth/claims) for more details. 

**Sample responses:**

**1. Scope used:** `openid email profile urn:zitadel:iam:org:project:id:zitadel:aud` 

**Sample response**: 

```bash
{
  "email": "david.wallace@dundermifflin.com",
  "email_verified": true,
  "family_name": "Wallace",
  "gender": "male",
  "given_name": "David",
  "locale": "en",
  "name": "David Wallace",
  "nickname": "David",
  "preferred_username": "david.wallace",
  "sub": "223427827918176513",
  "updated_at": 1689669364,
  "urn:zitadel:iam:org:project:223281986649719041:roles": {
    "cfo": {
      "223281939119866113": "corporate.user-authorizations-io8epz.zitadel.cloud"
    },
    "corporate member": {
      "223279178798072065": "org-a.user-authorizations-io8epz.zitadel.cloud",
      "223279223391912193": "org-b.user-authorizations-io8epz.zitadel.cloud"
    }
  },
  "urn:zitadel:iam:org:project:roles": {
    "cfo": {
      "223281939119866113": "corporate.user-authorizations-io8epz.zitadel.cloud"
    },
    "corporate member": {
      "223279178798072065": "org-a.user-authorizations-io8epz.zitadel.cloud",
      "223279223391912193": "org-b.user-authorizations-io8epz.zitadel.cloud"
    }
  }
}
```

This request can be tested out in the following way:
1. Select the **‘Assert Roles on Authentication’** checkbox.
2. Do not include the roles claims in the scope.
3. When you run the command, you will see that the roles were returned.
4. If you unselect the **‘Assert Roles on Authentication’** checkbox, you will not see the roles.


**2. Scope used:** `openid email profile urn:zitadel:iam:org:project:id:{projectId}:aud urn:iam:org:project:roles urn:zitadel:iam:org:projects:roles`

:::note
In order to stay up-to-date with the latest ZITADEL standards, we recommend that you use the roles from the identifier `urn:zitadel:iam:org:project:{projectId}:roles` rather than `urn:zitadel:iam:org:project:roles`. While both identifiers are maintained for backwards compatibility, the format which includes the specific ID represents our more recent model.
:::

**Sample response:** 

```bash
{
  "email": "david.wallace@dundermifflin.com",
  "email_verified": true,
  "family_name": "Wallace",
  "gender": "male",
  "given_name": "David",
  "locale": "en",
  "name": "David Wallace",
  "nickname": "David",
  "preferred_username": "david.wallace",
  "sub": "223427827918176513",
  "updated_at": 1689669364,
  "urn:zitadel:iam:org:project:223281986649719041:roles": {
    "cfo": {
      "223281939119866113": "corporate.user-authorizations-io8epz.zitadel.cloud"
    },
    "corporate member": {
      "223279178798072065": "org-a.user-authorizations-io8epz.zitadel.cloud",
      "223279223391912193": "org-b.user-authorizations-io8epz.zitadel.cloud"
    }
  },
  "urn:zitadel:iam:org:project:roles": {
    "cfo": {
      "223281939119866113": "corporate.user-authorizations-io8epz.zitadel.cloud"
    },
    "corporate member": {
      "223279178798072065": "org-a.user-authorizations-io8epz.zitadel.cloud",
      "223279223391912193": "org-b.user-authorizations-io8epz.zitadel.cloud"
    }
  }
}
```

This request can be tested out in the following way:
1. Do not select the **‘Assert Roles on Authentication’** checkbox 
2. Include the role claims in the scope as given.
3. When you run the command, you will see the roles in the response.
4. If you remove the role claims in the scope and run the command, you will not receive the roles.

### Retrieve roles using the auth API

Now we will use the auth API to retrieve roles from a logged in user using the user’s token
The base URL is: **https://$ZITADEL_DOMAIN/auth/v1**

Let’s start with a user who has multiple roles in different organizations in a multi-tenanted set up. You can use the logged in user’s token or the machine user’s token to retrieve the authorizations using the [APIs listed under user authorizations/grants in the auth API](/docs/category/apis/resources/auth/user-authorizations-grants). 

**Scope used:** `openid urn:zitadel:iam:org:project:id:zitadel:aud`


#### **1. [List my project roles](/docs/apis/resources/auth/auth-service-list-my-project-permissions)**

Returns a list of roles for the authenticated user and for the requesting project (based on the token).

**URL: https://$ZITADEL_DOMAIN/auth/v1/permissions/me/_search**

**cURL request:** 
```bash
curl -L -X POST 'https://$ZITADEL_DOMAIN/auth/v1/permissions/me/_search' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>'
```

**Sample response:** 

```bash
{
  "result": [
    "cfo"
  ]
}
```

#### **2.[List my ZITADEL permissions](/docs/apis/resources/auth/auth-service-list-my-zitadel-permissions)​**

Returns a list of permissions the authenticated user has in ZITADEL based on the manager roles the user has. (e.g: `ORG_OWNER` = `org.read`, `org.write`, ...).

This request can be used if you are building a management UI. For instance, if the UI is managing users, you can show the management functionality based on the permissions the user has. Here’s an example: if the user has `user.read` and `user.write` permission you can show the edit buttons, if the user only has `user.read` permission, you can hide the edit buttons.

**URL: https://ZITADEL_DOMAIN/auth/v1/permissions/zitadel/me/_search**

**cURL Request:** 

```bash
curl -L -X POST 'https://$ZITADEL_DOMAIN/auth/v1/permissions/zitadel/me/_search' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>'
```

**Sample result:** 

```bash
{
  "result": [
    "org.read",
    "user.read",
    "user.global.read",
    "user.write",
    "user.delete",
    "user.grant.read",
    "user.grant.write",
    "user.grant.delete",
    "user.membership.read",
    "policy.read",
    "project.read",
    "project.role.read",
    "org.member.read",
    "org.idp.read",
    "org.action.read",
    "org.flow.read",
    "project.member.read",
    "project.app.read",
    "project.grant.read",
    "project.grant.member.read",
    "project.grant.user.grant.read",
    "project.read:self",
    "project.create"
  ]
}
```

#### **[3. List my authorizations/grants​](/docs/apis/resources/auth/auth-service-list-my-user-grants)**

Returns a list of user grants the authenticated user has. User grants consist of an organization, a project and roles.

**URL: https://$ZITADEL_DOMAIN/auth/v1/usergrants/me/_search**

**cURL request:**

```bash
curl -L -X POST 'https://$ZITADEL_DOMAIN/auth/v1/usergrants/me/_search' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data-raw '{
  "query": {
    "offset": "0",
    "limit": 100,
    "asc": true
  }
}'
```

**Sample result:**
```bash
{
  "details": {
    "totalResult": "3",
    "processedSequence": "339",
    "viewTimestamp": "2023-07-19T09:20:34.371331Z"
  },
  "result": [
    {
      "orgId": "223279178798072065",
      "projectId": "223281986649719041",
      "userId": "223427827918176513",
      "roles": [
        "corporate member"
      ],
      "orgName": "Org A",
      "grantId": "223428842084106497",
      "details": {
        "sequence": "296",
        "creationDate": "2023-07-18T08:46:07.692435Z",
        "changeDate": "2023-07-18T08:46:07.692435Z",
        "resourceOwner": "223279178798072065"
      },
      "orgDomain": "org-a.user-authorizations-io8epz.zitadel.cloud",
      "projectName": "HR",
      "projectGrantId": "223282340514758913",
      "roleKeys": [
        "corporate member"
      ],
      "userType": "TYPE_HUMAN"
    },
    {
      "orgId": "223279223391912193",
      "projectId": "223281986649719041",
      "userId": "223427827918176513",
      "roles": [
        "corporate member"
      ],
      "orgName": "Org B",
      "grantId": "223428980244480257",
      "details": {
        "sequence": "298",
        "creationDate": "2023-07-18T08:47:30.015324Z",
        "changeDate": "2023-07-18T08:47:30.015324Z",
        "resourceOwner": "223279223391912193"
      },
      "orgDomain": "org-b.user-authorizations-io8epz.zitadel.cloud",
      "projectName": "HR",
      "projectGrantId": "223282930787549441",
      "roleKeys": [
        "corporate member"
      ],
      "userType": "TYPE_HUMAN"
    },
    {
      "orgId": "223281939119866113",
      "projectId": "223281986649719041",
      "userId": "223427827918176513",
      "roles": [
        "cfo"
      ],
      "orgName": "Corporate",
      "grantId": "223428420858544385",
      "details": {
        "sequence": "293",
        "creationDate": "2023-07-18T08:41:56.649257Z",
        "changeDate": "2023-07-18T08:44:33.094117Z",
        "resourceOwner": "223281939119866113"
      },
      "orgDomain": "corporate.user-authorizations-io8epz.zitadel.cloud",
      "projectName": "HR",
      "roleKeys": [
        "cfo"
      ],
      "userType": "TYPE_HUMAN"
    }
  ]
}
```

### Retrieve roles using the management API
Now we will use the management API to retrieve user roles under an admin user. 

The base URL is: **https://$ZITADEL_DOMAIN/management/v1**

In [APIs listed under user grants in the management API](/docs/category/apis/resources/mgmt/user-grants), you will see that you can use the management API to retrieve and modify user grants. The two API paths that we are interested in to fetch user roles are given below.

**Scope used:** `openid urn:zitadel:iam:org:project:id:zitadel:aud`

#### **1. [Search user grants](/docs/apis/resources/mgmt/management-service-list-user-grants)​**

Returns a list of user roles that match the search queries. A user with manager permissions will call this API and will also have to reside in the same organization as the user. 

**URL: https://$ZITADEL_DOMAIN/management/v1/users/grants/_search**

**cURL request:** 

```bash
curl -L -X POST 'https://$ZITADEL_DOMAIN/management/v1/users/grants/_search' \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>' \
--data '{
  "query": {
    "offset": "0",
    "limit": 100,
    "asc": true
  },
  "queries": [
    {
      "user_id_query": {
        "user_id": "223427827918176513"
      }
    }
  ]
}'
```

**Sample result:**

```bash
{
  "details": {
    "totalResult": "1",
    "processedSequence": "342",
    "viewTimestamp": "2023-07-19T11:24:58.769023Z"
  },
  "result": [
    {
      "id": "223428420858544385",
      "details": {
        "sequence": "293",
        "creationDate": "2023-07-18T08:41:56.649257Z",
        "changeDate": "2023-07-18T08:44:33.094117Z",
        "resourceOwner": "223281939119866113"
      },
      "roleKeys": [
        "cfo"
      ],
      "state": "USER_GRANT_STATE_ACTIVE",
      "userId": "223427827918176513",
      "userName": "david.wallace",
      "firstName": "David",
      "lastName": "Wallace",
      "email": "david.wallace@dundermifflin.com",
      "displayName": "David Wallace",
      "orgId": "223281939119866113",
      "orgName": "Corporate",
      "orgDomain": "corporate.user-authorizations-io8epz.zitadel.cloud",
      "projectId": "223281986649719041",
      "projectName": "HR",
      "preferredLoginName": "david.wallace",
      "userType": "TYPE_HUMAN"
    }
  ]
}
```

#### **2. [User grant by ID](/docs/apis/resources/mgmt/management-service-get-user-grant-by-id)​** 

Returns a user grant per ID. A user grant is a role a user has for a specific project and organization.

**URL: https://$ZITADEL_DOMAIN//management/v1/users/:userId/grants/:grantId**

**cURL request:**

```bash 
curl -L -X GET 'https://$ZITADEL_DOMAIN/management/v1/users/:userId/grants/:grantId' \
-H 'Accept: application/json' \
-H 'Authorization: Bearer <TOKEN>'
```


## Summary

The process of retrieving a user's roles involves understanding the audience scope, getting a token, and accessing the correct API endpoints based on your requirement. Following these steps will help efficiently manage roles in single and multi-tenant applications.


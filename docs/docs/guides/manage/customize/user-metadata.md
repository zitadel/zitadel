---
title: User Metadata
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

This guide shows you how to request metadata from a user.
ZITADEL offers multiple methods to retrieve metadata.
Pick the one that works best for your solution.

## Use cases for metadata

Typical use cases for user metadata include:

- Link the user to an internal identifier (eg, userId, contract number, etc.)
- Save custom user data when registering a user
- Route upstream traffic based on user attributes

## Before you start

Before you start you need to add some metadata to an existing user.
You can do so by using [Console](../console/users) or [setting user metadata](/docs/apis/resources/mgmt/management-service-set-user-metadata) through the management API.

Most of the methods below require you to login with the correct user while setting some scopes.
Make sure you pick the right user when logging into the test application.
Use the [OIDC authentication request playground](/docs/apis/openidoauth/authrequest) or the configuration of an [example client](/docs/examples/introduction) to set the required scopes and receive a valid access token.

:::info Getting a token
In case you want to test out different settings configure an application with code flow (PKCE).
Grab the code from the url parameter after a successful login and exchange the code for tokens by calling the [token endpoint](docs/apis/openidoauth/endpoints#token_endpoint).
:::

## Use tokens to get user metadata

Use one of these methods to get the metadata for the currently logged in user.

In case you want to manage metadata for other users than the currently logged in user, then you must use the [management service](#manage-user-metadata-through-the-management-api).

### Request metadata from userinfo endpoint

--> requires metadata scope in auth request! `urn:zitadel:iam:user:metadata`
--> OIDC
--> error?

With the access token we can make a request to the userinfo endpoint to get the user's metadata. This method is the preferred method to retrieve a user's information in combination with opaque tokens, to insure that the token is valid.

```bash
curl --request GET \
  --url "https://$ZITADEL_DOMAIN/oidc/v1/userinfo" \
  --header "Authorization: Bearer $ACCESS_TOKEN"
```

The response will look something like this

```json
{
    "email":"road.runner@zitadel.com",
    "email_verified":true,
    "family_name":"Runner",
    "given_name":"Road",
    "locale":"en",
    "name":"Road Runner",
    "preferred_username":"road.runner@...asd.zitadel.cloud",
    "sub":"166.....729",
    "updated_at":1655467738,
    //highlight-start
    "urn:zitadel:iam:user:metadata":{
        "ContractNumber":"MTIzNA",
        }
    //highlight-end
    }
```

```json
{"email":"mpa+admin.alice@zitadel.com","email_verified":true,"family_name":"Admin","given_name":"Alice","locale":null,"name":"Alice Admin","preferred_username":"admin.alice@demo-customer.b2b-demo-rbxajm.zitadel.cloud","sub":"170848145649959169","updated_at":1658329554,"urn:zitadel:iam:org:project:170086774599581953:roles":{"reader":{"170086363054473473":"demo-customer.b2b-demo-rbxajm.zitadel.cloud","190957560872829185":"demo-customer2.b2b-demo-rbxajm.zitadel.cloud"},"support:read":{"170086363054473473":"demo-customer.b2b-demo-rbxajm.zitadel.cloud"}},"urn:zitadel:iam:org:project:roles":{"reader":{"170086363054473473":"demo-customer.b2b-demo-rbxajm.zitadel.cloud","190957560872829185":"demo-customer2.b2b-demo-rbxajm.zitadel.cloud"},"support:read":{"170086363054473473":"demo-customer.b2b-demo-rbxajm.zitadel.cloud"}},"urn:zitadel:iam:user:metadata":{"P-ID":"MTIzNTE5ODM"}}
```

You can grab the metadata from the reserved claim `"urn:zitadel:iam:user:metadata"` as key-value pairs. Note that the values are base64 encoded. So the value `MTIzNA` decodes to `1234`.

### Send metadata inside the ID token

Check "User Info inside ID Token" in the configuration of your application.

![](/img/console_projects_application_token_settings.png)

Now request a new token from ZITADEL.

The result will give you something like:

```json
{
    "access_token":"jZuRixKQTVecEjKqw...kc3G4",
    "token_type":"Bearer",
    "expires_in":43199,
    "id_token":"ey...Ww"
}
```

```json
{
  "amr": [
    "password",
    "pwd",
    "mfa",
    "otp"
  ],
  "at_hash": "lGIblkTr8faHz2zd0oTddA",
  "aud": [
    "170086824411201793@portal",
    "209806276543185153@portal",
    "170086774599581953"
  ],
  "auth_time": 1687418556,
  "azp": "170086824411201793@portal",
  "c_hash": "dA3wre4ytCJCn11f7cIm0A",
  "client_id": "170086824411201793@portal",
  "email": "mpa+admin.alice@zitadel.com",
  "email_verified": true,
  "exp": 1687422272,
  "family_name": "Admin",
  "given_name": "Alice",
  "iat": 1687418672,
  "iss": "https://b2b-demo-rbxajm.zitadel.cloud",
  "locale": null,
  "name": "Alice Admin",
  "preferred_username": "admin.alice@demo-customer.b2b-demo-rbxajm.zitadel.cloud",
  "sub": "170848145649959169",
  "updated_at": 1658329554,
  "urn:zitadel:iam:org:project:170086774599581953:roles": {
    "reader": {
      "170086363054473473": "demo-customer.b2b-demo-rbxajm.zitadel.cloud",
      "190957560872829185": "demo-customer2.b2b-demo-rbxajm.zitadel.cloud"
    },
    "support:read": {
      "170086363054473473": "demo-customer.b2b-demo-rbxajm.zitadel.cloud"
    }
  },
  "urn:zitadel:iam:org:project:roles": {
    "reader": {
      "170086363054473473": "demo-customer.b2b-demo-rbxajm.zitadel.cloud",
      "190957560872829185": "demo-customer2.b2b-demo-rbxajm.zitadel.cloud"
    },
    "support:read": {
      "170086363054473473": "demo-customer.b2b-demo-rbxajm.zitadel.cloud"
    }
  },
  "urn:zitadel:iam:user:metadata": {
    "P-ID": "MTIzNTE5ODM"
  }
}
```

```bash
jq -R 'split(".") | .[1] | @base64d | fromjson' <<< $ID_TOKEN
```

### Request metadata from authentication API

--> omit queries to find all
--> `urn:zitadel:iam:org:project:id:zitadel:aud` else invalid audience (APP-Zxfako)

https://zitadel.com/docs/apis/resources/auth/auth-service-list-my-metadata

curl -L -X POST "https://$ZITADEL_DOMAIN/auth/v1/users/me/metadata/_search" \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H "Authorization: Bearer $ACCESS_TOKEN" \
--data-raw '{
  "query": {
    "offset": "0",
    "limit": 100,
    "asc": true
  },
  "queries": [
    {
      "keyQuery": {
        "key": "P-ID",
        "method": "TEXT_QUERY_METHOD_EQUALS"
      }
    }
  ]
}'

```json
{
    "details":{
        "totalResult":"1",
        "processedSequence":"2935",
        "viewTimestamp":"2023-06-21T16:01:52.829838Z"
    },
    "result":[
        {
            "details":{
                "sequence":"409",
                "creationDate":"2022-08-04T09:09:06.259324Z",
                "changeDate":"2022-08-04T09:09:06.259324Z",
                "resourceOwner":"170086363054473473"
                },
            "key":"P-ID",
            "value":"MTIzNTE5ODM="
        }
    ]
}
```

Grab the id_token and inspect the contents of the token at [jwt.io](https://jwt.io/). You should get the same info in the ID token as when requested from the user endpoint.

## Manage user metadata through the management API

:::warning
:::

http://localhost:3000/docs/apis/resources/mgmt/management-service-list-user-metadata
http://localhost:3000/docs/apis/resources/mgmt/management-service-get-user-metadata
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
Use the [OIDC authentication request playground](https://zitadel.com/playgrounds/oidc) or the configuration of an [example client](/docs/sdk-examples/introduction) to set the required scopes and receive a valid access token.

:::info Getting a token
In case you want to test out different settings configure an application with code flow (PKCE).
Grab the code from the url parameter after a successful login and exchange the code for tokens by calling the [token endpoint](/docs/apis/openidoauth/endpoints#token_endpoint).
You will find more information in our guides on how to [authenticate users](/docs/guides/integrate/login/oidc/login-users).
:::

## Use tokens to get user metadata

Use one of these methods to get the metadata for the currently logged in user.

In case you want to manage metadata for other users than the currently logged in user, then you must use the [management service](#manage-user-metadata-through-the-management-api).

### Request metadata from userinfo endpoint

With the access token we can make a request to the [userinfo endpoint](/docs/apis/openidoauth/endpoints#introspection_endpoint) to get the user's metadata.
This method is the preferred method to retrieve a user's information in combination with opaque tokens, to insure that the token is valid.

You must pass the [reserved scope](/docs/apis/openidoauth/scopes#reserved-scopes) `urn:zitadel:iam:user:metadata` in your authentication request.
If you don't include this scope the response will contain user data, but not the metadata object.

Request the user information by calling the [userinfo endpoint](/docs/apis/openidoauth/endpoints#introspection_endpoint):

```bash
curl --request GET \
  --url "https://$CUSTOM-DOMAIN/oidc/v1/userinfo" \
  --header "Authorization: Bearer $ACCESS_TOKEN"
```

Replace `$ACCESS_TOKEN` with your user's access token.

The response will look something like this

```json
{
  "email": "road.runner@zitadel.com",
  "email_verified": true,
  "family_name": "Runner",
  "given_name": "Road",
  "locale": "en",
  "name": "Road Runner",
  "preferred_username": "road.runner@...asd.zitadel.cloud",
  "sub": "166.....729",
  "updated_at": 1655467738,
  //highlight-start
  "urn:zitadel:iam:user:metadata": {
    "ContractNumber": "MTIzNA"
  }
  //highlight-end
}
```

You can grab the metadata from the reserved claim `"urn:zitadel:iam:user:metadata"` as key-value pairs.
Note that the values are base64 encoded.
So the value `MTIzNA` decodes to `1234`.

### Send metadata inside the ID token

You might want to include metadata directly into the ID Token.
For that you need to enable "User Info inside ID Token" in your application's settings.

![](/img/console_projects_application_token_settings.png)

Now request a new token from ZITADEL by logging in with the user that has metadata attached.
Make sure you log into the correct client/application where you enabled the settings.

The result will give you something like:

```json
{
  "access_token": "jZuRixKQTVecEjKqw...kc3G4",
  "token_type": "Bearer",
  "expires_in": 43199,
  "id_token": "ey...Ww"
}
```

When you decode the value of `id_token`, then the response will include the metadata claim:

```json
{
  "amr": ["password", "pwd", "mfa", "otp"],
  "at_hash": "lGIblkTr8faHz2zd0oTddA",
  "aud": [
    "170086824411201793@portal",
    "209806276543185153@portal",
    "170086774599581953"
  ],
  "auth_time": 1687418556,
  "azp": "170086824411201793@portal",
  "c_hash": "dA3wre4ytCJCn11f7cIm0A",
  "client_id": "1700...1793@portal",
  "email": "road.runner@zitadel.com",
  "email_verified": true,
  "exp": 1687422272,
  "family_name": "Runner",
  "given_name": "Road",
  "iat": 1687418672,
  "iss": "https://...-abcd.zitadel.cloud",
  "locale": null,
  "name": "Road Runner",
  "preferred_username": "road.runner@...-abcd.zitadel.cloud",
  "sub": "170848145649959169",
  "updated_at": 1658329554,
  //highlight-start
  "urn:zitadel:iam:user:metadata": {
    "ContractNumber": "MTIzNA"
  }
  //highlight-end
}
```

Note that the values are base64 encoded.
So the value `MTIzNA` decodes to `1234`.

:::info decoding the jwt token
Use a website like [jwt.io](https://jwt.io/) to decode the token.  
With jq installed you can also use `jq -R 'split(".") | .[1] | @base64d | fromjson' <<< $ID_TOKEN`
:::

### Request metadata from authentication API

You can use the authentication service to request and search for the user's metadata.

The introspection endpoint and the token endpoint in the examples above do not require a special scope to access.
Yet when accessing the authentication service, you need to pass the [reserved scope](/docs/apis/openidoauth/scopes#reserved-scopes) `urn:zitadel:iam:org:project:id:zitadel:aud` along with the authentication request.
This scope allows the user to access ZITADEL's APIs, specifically the authentication API that we need for this method.
Use the [OIDC authentication request playground](https://zitadel.com/playgrounds/oidc) or the configuration of an [example client](/docs/sdk-examples/introduction) to set the required scopes and receive a valid access token.

:::note Invalid audience
If you get the error "invalid audience (APP-Zxfako)", then you need to add the reserved scope `urn:zitadel:iam:org:project:id:zitadel:aud` to your authentication request.
:::

You can request the user's metadata with the [List My Metadata](/docs/apis/resources/auth/auth-service-list-my-metadata) method:

```bash
curl -L -X POST "https://$CUSTOM-DOMAIN/auth/v1/users/me/metadata/_search" \
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
        "key": "$METADATA_KEY",
        "method": "TEXT_QUERY_METHOD_EQUALS"
      }
    }
  ]
}'
```

Replace `$ACCESS_TOKEN` with your user's access token.  
Replace `$CUSTOM-DOMAIN` with your ZITADEL instance's url.  
Replace `$METADATA_KEY` with they key you want to search for (f.e. "ContractNumber")

:::info Get all metadata
You can omit the queries array to retrieve all metadata key-value pairs.
:::

An example response for your search looks like this:

```json
{
  "details": {
    "totalResult": "1",
    "processedSequence": "2935",
    "viewTimestamp": "2023-06-21T16:01:52.829838Z"
  },
  "result": [
    {
      "details": {
        "sequence": "409",
        "creationDate": "2022-08-04T09:09:06.259324Z",
        "changeDate": "2022-08-04T09:09:06.259324Z",
        "resourceOwner": "170086363054473473"
      },
      "key": "ContractNumber",
      "value": "MTIzNA"
    }
  ]
}
```

## Register user with custom metadata

When you build your own registration UI you have the possibility to have custom fields and add them to the metadata of your user.
Learn everything about how to build your own registration UI [here](/docs/guides/integrate/onboarding/end-users#build-your-own-registration-form).

## Manage user metadata through the management API

The previous methods allowed you to retrieve metadata only for the `sub` in the access token.
In case you want to get the metadata for another user, you need to use the management service.
The user that calls the management service must have [manager permissions](/docs/guides/manage/console/managers).
A user can be either a human user or a service user.

You can get [metadata of a user filtered by your query](/docs/apis/resources/mgmt/management-service-list-user-metadata) or [get a metadata object from a user by a specific key](/docs/apis/resources/mgmt/management-service-get-user-metadata).
The management service allows you to set and delete metadata, see the [API documentation for users](/docs/apis/resources/mgmt/users).

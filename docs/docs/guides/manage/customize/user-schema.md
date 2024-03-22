---
title: User Schema
---

ZITADEL allows you to define schemas for users, based on the [JSON Schema Standard](https://json-schema.org/). 
This gives you the possibility to define your own data models for your users, validate them based on these definitions
and making sure who has access or manipulate information of the user.

By defining multiple schemas, you can even differentiate between different personas of your organization or application
and restrictions, resp. requirements for them to authenticate. 

For example, you could have separate schemas for your employees and your customers. While you might want to make sure that
certain data like given name and family name are required for employees, they might be optional for the latter.
Or you might want to disable username password authentication for your admins and only allow phishing resistant methods like passkeys,
but still allow it for your customers.

:::info
Please be aware that User Schema is in a [preview stage](/support/software-release-cycles-support#preview) not feature complete
and therefore not generally available.

Do not use it for production yet. To test it out, you need to enable the `UserSchema` [feature flag](/apis/resources/feature_service_v2/feature-service).
:::

## Create your first schema

Let's create the first very simple schema `user`, which defines a `givenName` and `familyName` for a user and allows them to
authenticate with username and password.

We can do so by calling the [create user schema endpoint](/docs/apis/resources/user_schema_service_v3/user-schema-service-create-user-schema)
with the following data. Make sure to provide an access_token with an IAM_OWNER role.

```bash
curl -X POST "https://$CUSTOM-DOMAIN/v3alpha/user_schemas" \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H "Authorization: Bearer $ACCESS_TOKEN" \
--data-raw '{
  "type": "user",
  "schema": {
    "$schema": "urn:zitadel:schema:v1",
    "type": "object",
    "properties": {
      "givenName": {
        "type": "string"
      },
      "familyName": {
        "type": "string"
      }
    }
  },
  "possibleAuthenticators": [
    "AUTHENTICATOR_TYPE_USERNAME",
    "AUTHENTICATOR_TYPE_PASSWORD"
  ]
}'
```

This will return something similar to:
```json
{
  "id": "257199613398745508",
  "details": {
    "sequence": "2",
    "change_date": "2024-03-07T08:08:35.963956Z",
    "resource_owner": "253750309325636004"
  }
}
```

So you successfully create a schema and could use that to manage your users based on that.
But let's first checkout some possibilities ZITADEL offers.

## Assign Permissions

In the first step we've created a very simple `user` schema with only `givenName` and `familyName`.
This allows any user with the permission to edit the user's data to change these values.
Let's now update the schema and add some more properties and restrict who's able to retrieve and change data.

By setting `urn:zitadel:schema:permission` to fields, we can define the permissions for that field of different user roles / context.

For example by adding it to the `givenName` and `familyName` we can keep the state from before, where any `owner` (e.g. ORG_OWNER)
as well as the user themselves (`self`) are allowed to read (`r`) and write (`w`) the data.

Let's now assume our service provides some profile information of the user on a dedicated page.
Since we do not want the user to be able to change that value, we set the permission of `self` to `r`, meaning they will be able
to see the `profileUri` value, but cannot update it.

Maybe we also have some `customerId`, which the user should not even know about. We therefore can simply omit the `self` permission
and only set `owner` to `rw`, so admins are able to read and change the id if needed.

Finally, we call the [update user schema endpoint](/docs/apis/resources/user_schema_service_v3/user-schema-service-update-user-schema)
with the following data. Be sure to provide the id of the previously created schema.

```bash
curl -X PUT "https://$CUSTOM-DOMAIN/v3alpha/user_schemas/$SCHEMA_ID" \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H "Authorization: Bearer $ACCESS_TOKEN" \
--data-raw '{
  "schema": {
    "$schema": "urn:zitadel:schema:v1",
    "type": "object",
    "properties": {
      "givenName": {
        "type": "string",
        "urn:zitadel:schema:permission": {
          "owner": "rw",
          "self": "rw"
        }
      },
      "familyName": {
        "type": "string",
        "urn:zitadel:schema:permission": {
          "owner": "rw",
          "self": "rw"
        }
      },
      "profileUri": {
        "type": "string",
        "format": "uri",
        "urn:zitadel:schema:permission": {
          "owner": "rw",
          "self": "r"
        }
      },
      "customerId": {
        "type": "string",
        "urn:zitadel:schema:permission": {
          "owner": "rw"
        }
      }
    }
  }
}'
```

## Retrieve the Existing Schemas

To check the state of existing schemas you can simply [list them](/apis/resources/user_schema_service_v3/user-schema-service-list-user-schemas).
In this case we will query for the one with state `active`. Check out the api documentation for detailed information on possible filters.
The API also allows to retrieve a single [schema by ID](/apis/resources/user_schema_service_v3/user-schema-service-get-user-schema-by-id).

```bash
curl -X POST "https://$CUSTOM-DOMAIN/v3alpha/user_schemas/search" \
-H 'Content-Type: application/json' \
-H 'Accept: application/json' \
-H "Authorization: Bearer $ACCESS_TOKEN" \
--data-raw '{
  "query": {
    "offset": "0",
    "limit": 100,
    "asc": true
  },
  "sortingColumn": "FIELD_NAME_TYPE",
  "queries": [
    {
      "stateQuery": {
        "state": "STATE_ACTIVE"
      }
    }
  ]
}'
```

If you've followed this guide, it should list you a singe schema:

```json
{
  "details": {
    "totalResult": "1",
    "timestamp": "2024-03-21T16:35:19.685700Z"
  },
  "result": [
    {
      "id": "259279890237358500",
      "details": {
        "sequence": "2",
        "changeDate": "2024-03-21T16:35:19.685700Z",
        "resourceOwner": "224313188550750765"
      },
      "type": "user",
      "state": "STATE_ACTIVE",
      "revision": 2,
      "schema": {
        "$schema": "urn:zitadel:schema:v1",
        "properties": {
          "customerId": {
            "type": "string",
            "urn:zitadel:schema:permission": {
              "owner": "rw"
            }
          },
          "familyName": {
            "type": "string",
            "urn:zitadel:schema:permission": {
              "owner": "rw",
              "self": "rw"
            }
          },
          "givenName": {
            "type": "string",
            "urn:zitadel:schema:permission": {
              "owner": "rw",
              "self": "rw"
            }
          },
          "profileUri": {
            "format": "uri",
            "type": "string",
            "urn:zitadel:schema:permission": {
              "owner": "rw",
              "self": "r"
            }
          }
        },
        "type": "object"
      },
      "possibleAuthenticators": [
        "AUTHENTICATOR_TYPE_USERNAME",
        "AUTHENTICATOR_TYPE_PASSWORD"
      ]
    }
  ]
}
```

### Revision

Note the `revision` property, which is currently `2`. Each update to the `schema`-property will increase
it by `1`. The revision will later be reflected on the managed users to state based on which revision of the schema
they were last updated on.
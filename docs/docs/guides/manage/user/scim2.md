---
title: SCIM v2.0
---

:::info
The SCIM v2 interface of Zitadel is currently in a [preview stage](/support/software-release-cycles-support#preview).
It is not yet feature-complete, may contain bugs, and is not generally available.

Do not use it for production yet.

As long as the feature is in a preview state, it will be available for free, it will be put behind a commercial license once it is fully available.
:::

The Zitadel [SCIM v2](https://scim.cloud/) service provider interface enables seamless integration of identity and
access management (IAM) systems with Zitadel,
following the System for Cross-domain Identity Management (SCIM) v2.0 specification.
This interface allows standardized management of IAM resources, making it easier to automate user provisioning and
deprovisioning.

## Supported endpoints

The Zitadel SCIM v2.0 service provider implementation supports the following endpoints.
The base URL for the SCIM endpoint in Zitadel is: `https://${ZITADEL_DOMAIN}/scim/v2/{orgId}`.

| Endpoint                                                                | Remarks                                                    |
|-------------------------------------------------------------------------|------------------------------------------------------------|
| `GET /scim/v2/{orgId}/ServiceProviderConfig`                            | Retrieve the configuration of the Zitadel service provider |
| `GET /scim/v2/{orgId}/Schemas`                                          | Retrieve all supported schemas                             |
| `GET /scim/v2/{orgId}/Schemas/{id}`                                     | Retrieve a known supported schema                          |
| `GET /scim/v2/{orgId}/ResourceTypes`                                    | Retrieve all supported resource types                      |
| `GET /scim/v2/{orgId}/ResourceTypes/{name}`                             | Retrieve a known supported resource type                   |
| `GET /scim/v2/{orgId}/Users/{id}`                                       | Retrieve a known user                                      |
| `GET /scim/v2/{orgId}/Users`<br />`POST /scim/v2/{orgId}/Users/.search` | Query users (including filtering, sorting, paging)         |
| `POST /scim/v2/{orgId}/Users`                                           | Create a user                                              |
| `PUT /scim/v2/{orgId}/Users/{id}`                                       | Replace a user                                             |
| `PATCH /scim/v2/{orgId}/Users/{id}`                                     | Modify a user                                              |
| `DELETE /scim/v2/{orgId}/Users/{id}`                                    | Delete a user                                              |
| `POST /scim/v2/{orgId}/Bulk`                                            | Apply multiple operations in a single request              |

## Authentication

The SCIM interface adheres to Zitadel's standard API authentication methods.
For detailed instructions on authenticating with the SCIM interface, refer to the [Authenticate Service Users Guide](/guides/integrate/service-users/authenticate-service-users).

## Query

The list users endpoint supports sorting and filtering for both `GET /scim/v2/{orgId}/Users` and `POST /scim/v2/{orgId}/Users/.search` requests.
By default, the response includes up to 100 users, with a maximum allowable value for `count` set to 100.

### Sort

The following attributes are supported in the `SortBy` attribute.

- `meta.created`
- `meta.lastModified`
- `id`
- `username`
- `name.familyName`
- `name.givenName`
- `emails` and `emails.value`

### Filter

The following filter attributes and operators are supported:

| Attribute                    | Supported operators          |
|------------------------------|------------------------------|
| `meta.created`               | `EQ`, `GT`, `GE`, `LT`, `LE` |
| `meta.lastModified`          | `EQ`, `GT`, `GE`, `LT`, `LE` |
| `id`                         | `EQ`, `NE`, `CO`, `SW`, `EW` |
| `externalId`                 | `EQ`, `NE`                   |
| `username`                   | `EQ`, `NE`, `CO`, `SW`, `EW` |
| `name.familyName`            | `EQ`, `NE`, `CO`, `SW`, `EW` |
| `name.givenName`             | `EQ`, `NE`, `CO`, `SW`, `EW` |
| `emails`<br />`emails.value` | `EQ`, `NE`, `CO`, `SW`, `EW` |
| `active`                     | `EQ`, `NE`                   |

Filters can have a maximum length of 1000 characters.

## Examples

Here are practical examples demonstrating how to interact with the SCIM API,
providing clear guidance on common use cases such as creating a user.

<details>
<summary>Create a minimal user</summary>

```bash
curl -X POST "https://${DOMAIN}/scim/v2/${ORG_ID}/Users" \
  -H 'Content-Type: application/scim+json' \
  -H 'Accept: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-raw '
    {
      "schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
      "userName": "john.doe",
      "name": {
        "familyName": "Doe",
        "givenName": "John"
      },
      "password": "Password1!",
      "emails": [
        {
          "value": "john.doe@example.com",
          "primary": true
        }
      ]
    }
  '
```

</details>
<details>
<summary>List users created after a given date sorted by the creation date</summary>

```bash
curl -G "http://${DOMAIN}/scim/v2/${ORG_ID}/Users" \
  -H 'Accept: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-urlencode "sortBy=meta.created" \
  --data-urlencode "sortOrder=descending" \
  --data-urlencode "filter=meta.created gt \"2025-01-24T09:22:35.695245Z\""
```

</details>
<details>
<summary>Set a user inactive</summary>

```bash
curl -X PATCH "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/${USER_ID}" \
  -H 'Content-Type: application/scim+json' \
  -H 'Accept: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-raw '
    {
      "schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
      "Operations": [
        {
          "op": "replace",
          "path": "active",
          "value": false
        }
      ]
    }
  '
```

</details>
<details>
<summary>Set the password of a user</summary>

```bash
curl -X PATCH "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/${USER_ID}" \
  -H 'Content-Type: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-raw '
    {
      "schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
      "Operations": [
        {
          "op": "replace",
          "path": "password",
          "value": "Password2!"
        }
      ]
    }
  '
```

</details>
<details>
<summary>Delete a user</summary>

```bash
curl -X DELETE "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/${USER_ID}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
```

</details>

## Provisioning domain

A provisioning domain refers to an administrative domain that exists outside the domain of a service provider due to
legal or technical reasons.
For more details, refer to the [definitions](https://datatracker.ietf.org/doc/html/rfc7643#section-1.2)
of [RFC7643](https://datatracker.ietf.org/doc/html/rfc7643).

The `externalId` of a user is scoped to the provisioning domain.
To set a provisioning domain for a machine user,
add a metadata entry with the key `urn:zitadel:scim:provisioningDomain` and assign its value to the corresponding
provisioning domain.

When a machine user has a `urn:zitadel:scim:provisioningDomain` metadata set,
the `externalId` of all users provisioned or queried by that machine user is stored in the users' metadata.
The key format is `urn:zitadel:scim:{provisioningDomain}:externalId`,
where `{provisioningDomain}` is replaced with the machine user's provisioning domain.
If the machine user does not have a provisioning domain set,
a simplified metadata key `urn:zitadel:scim:externalId` is used to store and retrieve the `externalId` of users.

## Mapping

The table below outlines how supported SCIM attributes in the user schema map to corresponding Zitadel user attributes.
Some attributes are directly mapped to Zitadel user attributes, while others are stored in the user's metadata.
For more information about user metadata, see [here](../customize/user-metadata).

| SCIM                   | Zitadel                                                                                                   | Remarks                                                                                                                                                                                                                                        |
|------------------------|-----------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `id`                   | `userId`                                                                                                  |                                                                                                                                                                                                                                                |
| `username`             | `username`                                                                                                |                                                                                                                                                                                                                                                |
| `name.formatted`       | `profile.displayName`                                                                                     | The SCIM attribute `displayName` takes precedence over `name.formatted`                                                                                                                                                                        |
| `name.familyName`      | `profile.familyName`                                                                                      |                                                                                                                                                                                                                                                |
| `name.middleName`      | `metadata[urn:zitadel:scim:name.middleName]`                                                              |                                                                                                                                                                                                                                                |
| `name.honorificPrefix` | `metadata[urn:zitadel:scim:name.honorificPrefix]`                                                         |                                                                                                                                                                                                                                                |
| `name.honorificSuffix` | `metadata[urn:zitadel:scim:name.honorificSuffix]`                                                         |                                                                                                                                                                                                                                                |
| `displayName`          | `profile.displayName`                                                                                     | The SCIM attribute `displayName` takes precedence over `name.formatted`                                                                                                                                                                        |
| `nickName`             | `profile.nickName`                                                                                        |                                                                                                                                                                                                                                                |
| `profileUrl`           | `metadata[urn:zitadel:scim:profileUrl]`                                                                   |                                                                                                                                                                                                                                                |
| `title`                | `metadata[urn:zitadel:scim:title]`                                                                        |                                                                                                                                                                                                                                                |
| `preferredLanguage`    | `profile.preferredLanguage`                                                                               |                                                                                                                                                                                                                                                |
| `locale`               | `metadata[urn:zitadel:scim:locale]`                                                                       |                                                                                                                                                                                                                                                |
| `timezone`             | `metadata[urn:zitadel:scim:timezone]`                                                                     |                                                                                                                                                                                                                                                |
| `active`               | `state`                                                                                                   | `Initial` and `Active` are mapped to `active = true`, all other states are mapped to `active = false`.<br />The `active` value can only be updated if the user is in the state `Active` or `Inactive`.                                         |
| `password`             | `password`                                                                                                |                                                                                                                                                                                                                                                |
| `emails`               | `email`                                                                                                   | Only the `primary` email is stored in Zitadel, if there is no `primary` email, the first one is stored. By default emails from SCIM are considered verified, this can be adjusted in the [configuration](#configuration).                      |
| `phoneNumbers`         | `phone`                                                                                                   | Only the `primary` phone number is stored in Zitadel, if there is no `primary` phone number, the first one is stored. By default phone numbers from SCIM are considered verified, this can be adjusted in the [configuration](#configuration). |
| `ims`                  | `metadata[urn:zitadel:scim:ims]`                                                                          | Serialized as JSON.                                                                                                                                                                                                                            |
| `photos`               | `metadata[urn:zitadel:scim:photos]`                                                                       | Serialized as JSON.                                                                                                                                                                                                                            |
| `addresses`            | `metadata[urn:zitadel:scim:addresses]`                                                                    | Serialized as JSON.                                                                                                                                                                                                                            |
| `entitlements`         | `metadata[urn:zitadel:scim:entitlements]`                                                                 | Serialized as JSON.                                                                                                                                                                                                                            |
| `roles`                | `metadata[urn:zitadel:scim:roles]`                                                                        | Serialized as JSON.                                                                                                                                                                                                                            |
| `externalId`           | `metadata[urn:zitadel:scim:externalId]`<br />`metadata[urn:zitadel:scim:{provisioningDomain}:externalId]` | See [provisioning domain](#provisioning-domain).                                                                                                                                                                                               |

## Error handling

The SCIM interface uses standard HTTP status codes and error messages to indicate the success or failure of API requests
following the error handling guidelines of [RFC7644](https://datatracker.ietf.org/doc/html/rfc7644#section-3.12).

In addition to the default SCIM error schema (`urn:ietf:params:scim:api:messages:2.0:Error`),
Zitadel extends the error response with a custom schema, `urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail`.
This schema includes additional attributes, such as the untranslated error message and an error id,
which aids pinpointing the source of the error in the system.

## Configuration

This section provides details on the runtime configuration of the SCIM interface of Zitadel.

By default, Zitadel's SCIM interface assumes that email addresses and phone numbers are verified.  
The bulk endpoint supports up to 100 operations per request, with a maximum request body size of 1 MB.  
This behavior can be adjusted through the Zitadel runtime configuration settings:

```yaml
SCIM:
  EmailVerified: true
  PhoneVerified: true
  MaxRequestBodySize: 1_000_000
  Bulk:
    MaxOperationsCount: 100
 ```

## Limitations

This section outlines the known limitations of the Zitadel SCIM implementation,
including unsupported features, partial compliance with the SCIM specification,
and any potential edge cases to consider during integration.

### Discovery

The discovery endpoints `GET /ServiceProviderConfig`, `GET /ResourceTypes` and `GET /Schemas` are not yet supported.

### Supported schemas

Only the users schema `urn:ietf:params:scim:schemas:core:2.0:User` is supported.

### Required attributes

The following SCIM user attributes are required, in addition to those required by the SCIM standard:

* `name.familyName`
* `name.givenName`
* `emails`: at least one email is required

### Duplicated attribute mapping

The SCIM user attributes `name.formatted` and `displayName` are both mapped to the `profile.displayName` attribute in
Zitadel.
When a user is provisioned with different values for these attributes, `displayName` takes precedence.
Only the value of `displayName` is stored and returned in subsequent queries.

## Resources

- **[SCIM](https://scim.cloud/)**: The Webpage of SCIM.
- **[RFC7643](https://tools.ietf.org/html/rfc7643) Core Schema**:
  The Core Schema provides a platform-neutral schema and extension model for representing users and groups.
- **[RFC7644](https://tools.ietf.org/html/rfc7644) Protocol**:
  The SCIM Protocol is an application-level, REST protocol for provisioning and managing identity data on the web.
- **[RFC7642](https://tools.ietf.org/html/rfc7642) Definitions, Overview, Concepts, and Requirements**:
  This document lists the user scenarios and use cases of System for Cross-domain Identity Management (SCIM).

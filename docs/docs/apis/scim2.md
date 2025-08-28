---
title: SCIM v2.0 (Preview)
---

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
Make sure to replace any placeholder values (`${}`) with the actual values from your environment.

<details>
<summary>`POST /Users`: Create a minimal user</summary>

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

**Response (`201 Created`)**
```json
{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:User"
  ],
  "meta": {
    "resourceType": "User",
    "created": "2025-01-27T15:30:27.651321Z",
    "lastModified": "2025-01-27T15:30:27.651321Z",
    "version": "2",
    "location": "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/304499468865155777"
  },
  "id": "304499468865155777",
  "userName": "john.doe",
  "name": {
    "familyName": "Doe",
    "givenName": "John"
  },
  "preferredLanguage": "en",
  "emails": [
    {
      "value": "john.doe@example.com",
      "primary": true
    }
  ]
}
```

</details>
<details>
<summary>`POST /Users`: Create a full user</summary>

```bash
curl -X POST "https://${DOMAIN}/scim/v2/${ORG_ID}/Users" \
  -H 'Content-Type: application/scim+json' \
  -H 'Accept: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-raw '
    {
      "schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
      "externalId": "8d4b51c0-51bd-4386-ae17-79ce5fd36517",
      "userName": "john.doe@example.com",
      "name": {
        "formatted": "Mr. John J Doe, III",
        "familyName": "Doe",
        "givenName": "John",
        "middleName": "Jim",
        "honorificPrefix": "Mr.",
        "honorificSuffix": "III"
      },
      "displayName": "John Doe",
      "nickName": "Johnny",
      "profileUrl": "https://login.example.com/john.doe",
      "emails": [
        {
          "value": "john.doe@example.com",
          "type": "work",
          "primary": true
        }
      ],
      "addresses": [
        {
          "type": "work",
          "streetAddress": "100 Universal City Plaza",
          "locality": "Hollywood",
          "region": "CA",
          "postalCode": "91608",
          "country": "USA",
          "formatted": "100 Universal City Plaza\nHollywood, CA 91608 USA",
          "primary": true
        }
      ],
      "phoneNumbers": [
        {
          "value": "+1 555-555-5555",
          "type": "work",
          "primary": true
        }
      ],
      "ims": [
        {
          "value": "@j.doe",
          "type": "X"
        }
      ],
      "photos": [
        {
          "value": "https://photos.example.com/profilephoto/john.doe/F",
          "type": "photo"
        }
      ],
      "roles": [
        {
          "value": "user-admin",
          "display": "User administrator"
        }
      ],
      "entitlements": [
        {
          "value": "read-passports",
          "display": "Read Passports"
        }
      ],
      "userType": "Employee",
      "title": "Tour Guide",
      "preferredLanguage": "en-US",
      "locale": "en-US",
      "timezone": "America/Los_Angeles",
      "active": true,
      "password": "Password1!"
    }'
```

**Response (`201 Created`)**
```json
{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:User"
  ],
  "meta": {
    "resourceType": "User",
    "created": "2025-01-27T15:31:47.84572Z",
    "lastModified": "2025-01-27T15:31:47.84572Z",
    "version": "16",
    "location": "https://localhost:8080/scim/v2/303879575732073153/Users/304499603368096449"
  },
  "id": "304499603368096449",
  "externalId": "8d4b51c0-51bd-4386-ae17-79ce5fd36517",
  "userName": "john.doe@example.com",
  "name": {
    "formatted": "John Doe",
    "familyName": "Doe",
    "givenName": "John",
    "middleName": "Jim",
    "honorificPrefix": "Mr.",
    "honorificSuffix": "III"
  },
  "displayName": "John Doe",
  "nickName": "Johnny",
  "profileUrl": "https://login.example.com/john.doe",
  "title": "Tour Guide",
  "preferredLanguage": "en-US",
  "locale": "en-US",
  "timezone": "America/Los_Angeles",
  "active": true,
  "emails": [
    {
      "value": "john.doe@example.com",
      "primary": true
    }
  ],
  "phoneNumbers": [
    {
      "value": "+15555555555",
      "primary": true
    }
  ],
  "ims": [
    {
      "value": "@j.doe",
      "type": "X"
    }
  ],
  "addresses": [
    {
      "type": "work",
      "streetAddress": "100 Universal City Plaza",
      "locality": "Hollywood",
      "region": "CA",
      "postalCode": "91608",
      "country": "USA",
      "formatted": "100 Universal City Plaza\nHollywood, CA 91608 USA",
      "primary": true
    }
  ],
  "photos": [
    {
      "value": "https://photos.example.com/profilephoto/john.doe/F",
      "type": "photo"
    }
  ],
  "entitlements": [
    {
      "value": "read-passports",
      "display": "Read Passports"
    }
  ],
  "roles": [
    {
      "value": "user-admin",
      "display": "User administrator"
    }
  ]
}
```

</details>
<details>
<summary>`GET /Users/{id}`: Retrive a known user</summary>

```bash
curl -G "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/${USER_ID}" \
  -H 'Accept: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
```


**Response (`200 OK`)**
```json
{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:User"
  ],
  "meta": {
    "resourceType": "User",
    "created": "2025-01-27T15:31:47.84572Z",
    "lastModified": "2025-01-27T15:31:47.84572Z",
    "version": "16",
    "location": "https://localhost:8080/scim/v2/303879575732073153/Users/304499603368096449"
  },
  "id": "304499603368096449",
  "externalId": "8d4b51c0-51bd-4386-ae17-79ce5fd36517",
  "userName": "john.doe@example.com",
  "name": {
    "formatted": "John Doe",
    "familyName": "Doe",
    "givenName": "John",
    "middleName": "Jim",
    "honorificPrefix": "Mr.",
    "honorificSuffix": "III"
  },
  "displayName": "John Doe",
  "nickName": "Johnny",
  "profileUrl": "https://login.example.com/john.doe",
  "title": "Tour Guide",
  "preferredLanguage": "en-US",
  "locale": "en-US",
  "timezone": "America/Los_Angeles",
  "active": true,
  "emails": [
    {
      "value": "john.doe@example.com",
      "primary": true
    }
  ],
  "phoneNumbers": [
    {
      "value": "+15555555555",
      "primary": true
    }
  ],
  "ims": [
    {
      "value": "@j.doe",
      "type": "X"
    }
  ],
  "addresses": [
    {
      "type": "work",
      "streetAddress": "100 Universal City Plaza",
      "locality": "Hollywood",
      "region": "CA",
      "postalCode": "91608",
      "country": "USA",
      "formatted": "100 Universal City Plaza\nHollywood, CA 91608 USA",
      "primary": true
    }
  ],
  "photos": [
    {
      "value": "https://photos.example.com/profilephoto/john.doe/F",
      "type": "photo"
    }
  ],
  "entitlements": [
    {
      "value": "read-passports",
      "display": "Read Passports"
    }
  ],
  "roles": [
    {
      "value": "user-admin",
      "display": "User administrator"
    }
  ]
}
```

</details>
<details>
<summary>`GET /Users`: List users created after a given date sorted by the creation date</summary>

```bash
curl -G "https://${DOMAIN}/scim/v2/${ORG_ID}/Users" \
  -H 'Accept: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-urlencode "sortBy=meta.created" \
  --data-urlencode "sortOrder=descending" \
  --data-urlencode "filter=meta.created gt \"2025-01-24T09:22:35.695245Z\""
```

**Response (`200 OK`)**
```json
{
  "schemas": ["urn:ietf:params:scim:api:messages:2.0:ListResponse"],
  "itemsPerPage": 100,
  "totalResults": 1,
  "startIndex": 1,
  "Resources": [
    {
      "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:User"
      ],
      "meta": {
        "resourceType": "User",
        "created": "2025-01-27T15:31:47.84572Z",
        "lastModified": "2025-01-27T15:31:47.84572Z",
        "version": "3",
        "location": "https://localhost:8080/scim/v2/303879575732073153/Users/304499603368096449"
      },
      "id": "304499603368096449",
      "externalId": "8d4b51c0-51bd-4386-ae17-79ce5fd36517",
      "userName": "john.doe@example.com",
      "name": {
        "formatted": "John Doe",
        "familyName": "Doe",
        "givenName": "John",
        "middleName": "Jim",
        "honorificPrefix": "Mr.",
        "honorificSuffix": "III"
      },
      "displayName": "John Doe",
      "nickName": "Johnny",
      "profileUrl": "https://login.example.com/john.doe",
      "title": "Tour Guide",
      "preferredLanguage": "und",
      "locale": "en-US",
      "timezone": "America/Los_Angeles",
      "active": true,
      "emails": [
        {
          "value": "john.doe@example.com",
          "primary": true
        }
      ],
      "phoneNumbers": [
        {
          "value": "+15555555555",
          "primary": true
        }
      ],
      "ims": [
        {
          "value": "@j.doe",
          "type": "X"
        }
      ],
      "addresses": [
        {
          "type": "work",
          "streetAddress": "100 Universal City Plaza",
          "locality": "Hollywood",
          "region": "CA",
          "postalCode": "91608",
          "country": "USA",
          "formatted": "100 Universal City Plaza\nHollywood, CA 91608 USA",
          "primary": true
        }
      ],
      "photos": [
        {
          "value": "https://photos.example.com/profilephoto/john.doe/F",
          "type": "photo"
        }
      ],
      "entitlements": [
        {
          "value": "read-passports",
          "display": "Read Passports"
        }
      ],
      "roles": [
        {
          "value": "user-admin",
          "display": "User administrator"
        }
      ]
    }
  ]
}
```

</details>
<details>
<summary>`PATCH /Users/{id}`: Set a user inactive</summary>

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

**Response**: `204 No Content`

</details>
<details>
<summary>`PATCH /Users/{id}`: Set the password of a user</summary>

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

**Response**: `204 No Content`

</details>
<details>
<summary>`PUT /Users/{id}`: Replace a full user</summary>

```bash
curl -X PUT "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/${USER_ID}" \
  -H 'Content-Type: application/scim+json' \
  -H 'Accept: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-raw '
    {
      "schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
      "externalId": "8d4b51c0-51bd-4386-ae17-79ce5fd36517",
      "userName": "john.doe@example.com",
      "name": {
        "formatted": "Mr. John J Doe, III",
        "familyName": "Doe",
        "givenName": "John",
        "middleName": "Jim",
        "honorificPrefix": "Mr.",
        "honorificSuffix": "III"
      },
      "displayName": "John Doe",
      "nickName": "Johnny",
      "profileUrl": "https://login.example.com/john.doe",
      "emails": [
        {
          "value": "john.doe@example.com",
          "type": "work",
          "primary": true
        }
      ],
      "addresses": [
        {
          "type": "work",
          "streetAddress": "100 Universal City Plaza",
          "locality": "Hollywood",
          "region": "CA",
          "postalCode": "91608",
          "country": "USA",
          "formatted": "100 Universal City Plaza\nHollywood, CA 91608 USA",
          "primary": true
        }
      ],
      "phoneNumbers": [
        {
          "value": "+1 555-555-5555",
          "type": "work",
          "primary": true
        }
      ],
      "ims": [
        {
          "value": "@j.doe",
          "type": "X"
        }
      ],
      "photos": [
        {
          "value": "https://photos.example.com/profilephoto/john.doe/F",
          "type": "photo"
        }
      ],
      "roles": [
        {
          "value": "user-admin",
          "display": "User administrator"
        }
      ],
      "entitlements": [
        {
          "value": "read-passports",
          "display": "Read Passports"
        }
      ],
      "userType": "Employee",
      "title": "Tour Guide",
      "preferredLanguage": "en-US",
      "locale": "en-US",
      "timezone": "America/Los_Angeles",
      "active": true,
      "password": "Password1!"
    }'
```

**Response (`200 OK`)**
```json
{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:User"
  ],
  "meta": {
    "resourceType": "User",
    "created": "2025-01-27T15:31:47.84572Z",
    "lastModified": "2025-01-27T15:31:47.84572Z",
    "version": "16",
    "location": "https://localhost:8080/scim/v2/303879575732073153/Users/304499603368096449"
  },
  "id": "304499603368096449",
  "externalId": "8d4b51c0-51bd-4386-ae17-79ce5fd36517",
  "userName": "john.doe@example.com",
  "name": {
    "formatted": "John Doe",
    "familyName": "Doe",
    "givenName": "John",
    "middleName": "Jim",
    "honorificPrefix": "Mr.",
    "honorificSuffix": "III"
  },
  "displayName": "John Doe",
  "nickName": "Johnny",
  "profileUrl": "https://login.example.com/john.doe",
  "title": "Tour Guide",
  "preferredLanguage": "en-US",
  "locale": "en-US",
  "timezone": "America/Los_Angeles",
  "active": true,
  "emails": [
    {
      "value": "john.doe@example.com",
      "primary": true
    }
  ],
  "phoneNumbers": [
    {
      "value": "+15555555555",
      "primary": true
    }
  ],
  "ims": [
    {
      "value": "@j.doe",
      "type": "X"
    }
  ],
  "addresses": [
    {
      "type": "work",
      "streetAddress": "100 Universal City Plaza",
      "locality": "Hollywood",
      "region": "CA",
      "postalCode": "91608",
      "country": "USA",
      "formatted": "100 Universal City Plaza\nHollywood, CA 91608 USA",
      "primary": true
    }
  ],
  "photos": [
    {
      "value": "https://photos.example.com/profilephoto/john.doe/F",
      "type": "photo"
    }
  ],
  "entitlements": [
    {
      "value": "read-passports",
      "display": "Read Passports"
    }
  ],
  "roles": [
    {
      "value": "user-admin",
      "display": "User administrator"
    }
  ]
}
```

</details>
<details>
<summary>`DELETE /Users/{id}`: Delete a user</summary>

```bash
curl -X DELETE "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/${USER_ID}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
```

**Response**: `204 No Content`

</details>
<details>
<summary>`POST /Bulk`: Update the password of one user and delete another one</summary>

```bash
curl -X POST "https://${DOMAIN}/scim/v2/${ORG_ID}/Bulk" \
  -H 'Content-Type: application/scim+json' \
  -H 'Accept: application/scim+json' \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  --data-raw '
    {
      "schemas": ["urn:ietf:params:scim:api:messages:2.0:BulkRequest"],
      "Operations": [
        {
          "method": "PATCH",
          "path": "/Users/${USER_ID}",
          "data": {
            "schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
            "Operations": [
              {
                "op": "replace",
                "path": "password",
                "value": "Password2!"
              }
            ]
          }
        },
        {
          "method": "DELETE",
          "path": "/Users/${USER_ID2}"
        }
      ]
    }'
```

**Response**: `200 OK`

```json
{
  "schemas": ["urn:ietf:params:scim:api:messages:2.0:BulkResponse"],
  "Operations": [
    {
      "method": "PATCH",
      "location": "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/${USER_ID}",
      "status": "204"
    },
    {
      "method": "DELETE",
      "location": "https://${DOMAIN}/scim/v2/${ORG_ID}/Users/${USER_ID2}",
      "status": "204"
    }
  ]
}
```

</details>

<details>
<summary>`GET /ServiceProviderConfig`: Get service provider configuration</summary>

```bash
curl -G "https://${DOMAIN}/scim/v2/${ORG_ID}/ServiceProviderConfig" \
  -H 'Accept: application/scim+json'
```

**Response**: `200 OK`

```json
{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"
  ],
  "meta": {
    "resourceType": "ServiceProviderConfig",
    "location": "https://${DOMAIN}/scim/v2/${ORG_ID}/ServiceProviderConfig"
  },
  "documentationUri": "https://zitadel.com/docs/guides/manage/user/scim2",
  "patch": {
    "supported": true
  },
  "bulk": {
    "supported": true,
    "maxOperations": 100,
    "maxPayloadSize": 1000000
  },
  "filter": {
    "supported": true,
    "maxResults": 100
  },
  "changePassword": {
    "supported": true
  },
  "sort": {
    "supported": true
  },
  "etag": {
    "supported": false
  },
  "authenticationSchemes": [
    {
      "name": "Zitadel authentication token",
      "description": "Authentication scheme using the OAuth Bearer Token Standard",
      "specUri": "https://www.rfc-editor.org/info/rfc6750",
      "documentationUri": "https://zitadel.com/docs/guides/integrate/service-users/authenticate-service-users",
      "type": "oauthbearertoken",
      "primary": false
    }
  ]
}
```

</details>

## Error handling

The SCIM interface uses standard HTTP status codes and error messages to indicate the success or failure of API requests
following the error handling guidelines of [RFC7644](https://datatracker.ietf.org/doc/html/rfc7644#section-3.12).

In addition to the default SCIM error schema (`urn:ietf:params:scim:api:messages:2.0:Error`),
Zitadel extends the error response with a custom schema, `urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail`.
This schema includes additional attributes, such as the untranslated error message and an error id,
which aids pinpointing the source of the error in the system.

## Resources

- **[Zitadel SCIM Documentation](/guides/manage/user/scim2)**: Documentation of Zitadel's SCIM implementation, including configuration details and known limitations.
- **[SCIM](https://scim.cloud/)**: The Webpage of SCIM.
- **[RFC7643](https://tools.ietf.org/html/rfc7643) Core Schema**:
  The Core Schema provides a platform-neutral schema and extension model for representing users and groups.
- **[RFC7644](https://tools.ietf.org/html/rfc7644) Protocol**:
  The SCIM Protocol is an application-level, REST protocol for provisioning and managing identity data on the web.
- **[RFC7642](https://tools.ietf.org/html/rfc7642) Definitions, Overview, Concepts, and Requirements**:
  This document lists the user scenarios and use cases of System for Cross-domain Identity Management (SCIM).

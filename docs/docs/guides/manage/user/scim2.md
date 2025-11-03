---
title: SCIM v2.0 (Preview)
---

The Zitadel [SCIM v2](https://scim.cloud/) service provider interface enables seamless integration of identity and
access management (IAM) systems with Zitadel,
following the System for Cross-domain Identity Management (SCIM) v2.0 specification.
This interface allows standardized management of IAM resources, making it easier to automate user provisioning and
deprovisioning.

## API

To learn more about Zitadel's SCIM API, see the API documentation [here](/apis/scim2).

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

- **[Zitadel SCIM API Documentation](/apis/scim2)**: Documentation of Zitadel's SCIM API implementation.
- **[SCIM](https://scim.cloud/)**: The Webpage of SCIM.
- **[RFC7643](https://tools.ietf.org/html/rfc7643) Core Schema**:
  The Core Schema provides a platform-neutral schema and extension model for representing users and groups.
- **[RFC7644](https://tools.ietf.org/html/rfc7644) Protocol**:
  The SCIM Protocol is an application-level, REST protocol for provisioning and managing identity data on the web.
- **[RFC7642](https://tools.ietf.org/html/rfc7642) Definitions, Overview, Concepts, and Requirements**:
  This document lists the user scenarios and use cases of System for Cross-domain Identity Management (SCIM).

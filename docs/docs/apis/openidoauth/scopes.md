---
title: Scopes
---

ZITADEL supports the usage of scopes as way of requesting information from the IAM and also instruct ZITADEL to do certain operations.

## Standard Scopes

| Scopes         | Description                                                                    |
| :------------- | ------------------------------------------------------------------------------ |
| openid         | When using openid connect this is a mandatory scope                            |
| profile        | Optional scope to request the profile of the subject                           |
| email          | Optional scope to request the email of the subject                             |
| address        | Optional scope to request the address of the subject                           |
| offline_access | Optional scope to request a refresh_token (only possible when using code flow) |

## Custom Scopes

> This feature is not yet released

## Reserved Scopes

In addition to the standard compliant scopes we utilize the following scopes.

| Scopes                                            | Example                                                | Description                                                                                                                                                                                                                                                           |
| :------------------------------------------------ | :----------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `urn:zitadel:iam:org:project:role:{rolekey}`      | `urn:zitadel:iam:org:project:role:user`                | By using this scope a client can request the claim urn:zitadel:iam:roles to be asserted when possible. As an alternative approach you can enable all roles to be asserted from the [project](/guides/manage/console/roles#authorizations) a client belongs to.   |
| `urn:zitadel:iam:org:id:{id}`                     | `urn:zitadel:iam:org:id:178204173316174381`            | When requesting this scope **ZITADEL** will enforce that the user is a member of the selected organization. If the organization does not exist a failure is displayed. It will assert the `urn:zitadel:iam:user:resourceowner` claims.                                |
| `urn:zitadel:iam:org:domain:primary:{domainname}` | `urn:zitadel:iam:org:domain:primary:acme.ch`           | When requesting this scope **ZITADEL** will enforce that the user is a member of the selected organization and the username is suffixed by the provided domain. If the organization does not exist a failure is displayed                                             |
| `urn:zitadel:iam:role:{rolename}`                 |                                                        |                                                                                                                                                                                                                                                                       |
| `urn:zitadel:iam:org:project:id:{projectid}:aud`  | `urn:zitadel:iam:org:project:id:69234237810729019:aud` | By adding this scope, the requested projectid will be added to the audience of the access token                                                                                                                                                                       |
| `urn:zitadel:iam:org:project:id:zitadel:aud`      | `urn:zitadel:iam:org:project:id:zitadel:aud`           | By adding this scope, the ZITADEL project ID will be added to the audience of the access token                                                                                                                                                                        |
| `urn:zitadel:iam:user:metadata`                   | `urn:zitadel:iam:user:metadata`                        | By adding this scope, the metadata of the user will be included in the token. The values are base64 encoded.                                                                                                                                                          |
| `urn:zitadel:iam:user:resourceowner`              | `urn:zitadel:iam:user:resourceowner`                   | By adding this scope, the resourceowner (id, name, primary_domain) of the user will be included in the token.                                                                                                                                                         |
| `urn:zitadel:iam:org:idp:id:{idp_id}`             | `urn:zitadel:iam:org:idp:id:76625965177954913`         | By adding this scope the user will directly be redirected to the identity provider to authenticate. Make sure you also send the primary domain scope if a custom login policy is configured. Otherwise the system will not be able to identify the identity provider. |

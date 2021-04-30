---
title: Claims
---

ZITADEL asserts claims on different places according to the corresponding specifications or project and clients settings.
Please check below the matrix for an overview where which scope is asserted.

| Claims                                          | Userinfo       | Introspection  | ID Token                                    | Access Token                         |
|:------------------------------------------------|:---------------|----------------|---------------------------------------------|--------------------------------------|
| acr                                             | No             | No             | Yes                                         | No                                   |
| address                                         | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| amr                                             | No             | No             | Yes                                         | No                                   |
| aud                                             | No             | No             | Yes                                         | When JWT                             |
| auth_time                                       | No             | No             | Yes                                         | No                                   |
| azp                                             | No             | No             | Yes                                         | When JWT                             |
| email                                           | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| email_verified                                  | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| exp                                             | No             | No             | Yes                                         | When JWT                             |
| family_name                                     | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| gender                                          | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| given_name                                      | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| iat                                             | No             | No             | Yes                                         | When JWT                             |
| iss                                             | No             | No             | Yes                                         | When JWT                             |
| locale                                          | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| name                                            | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| nonce                                           | No             | No             | Yes                                         | No                                   |
| phone                                           | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| phone_verified                                  | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| preferred_username (username when Introspect )  | When requested | When requested | Yes                                         | No                                   |
| sub                                             | Yes            | Yes            | Yes                                         | When JWT                             |
| urn:zitadel:iam:org:domain:primary:{domainname} | When requested | When requested | When requested                              | When JWT and requested               |
| urn:zitadel:iam:org:project:roles:{rolename}    | When requested | When requested | When requested or configured                | When JWT and requested or configured |

## Standard Claims

| Claims             | Example                                  | Description                                                                                   |
|:-------------------|:-----------------------------------------|-----------------------------------------------------------------------------------------------|
| acr                | TBA                                      | TBA                                                                                           |
| address            | `Teufener Strasse 19, 9000 St. Gallen`   | TBA                                                                                           |
| amr                | `pwd mfa`                                | Authentication Method References as defined in [RFC8176](https://tools.ietf.org/html/rfc8176) |
| aud                | `69234237810729019`                      | By default all client id's and the project id is included                                     |
| auth_time          | `1311280969`                             | Unix time of the authentication                                                               |
| azp                | `69234237810729234`                      | Client id of the client who requested the token                                               |
| email              | `road.runner@acme.ch`                    | Email Address of the subject                                                                  |
| email_verified     | `true`                                   | Boolean if the email was verified by ZITADEL                                                  |
| exp                | `1311281970`                             | Time the token expires as unix time                                                           |
| family_name        | `Runner`                                 | The subjects family name                                                                      |
| gender             | `other`                                  | Gender of the subject                                                                         |
| given_name         | `Road`                                   | Given name of the subject                                                                     |
| iat                | `1311280970`                             | Issued at time of the token as unix time                                                      |
| iss                | `https://issuer.zitadel.ch`              | Issuing domain of a token                                                                     |
| locale             | `en`                                     | Language from the subject                                                                     |
| name               | `Road Runner`                            | The subjects full name                                                                        |
| nonce              | `blQtVEJHNTF0WHhFQmhqZ0RqeHJsdzdkd2d...` | The nonce provided by the client                                                              |
| phone              | `+41 79 XXX XX XX`                       | Phone number provided by the user                                                             |
| preferred_username | `road.runner@acme.caos.ch`               | ZITADEL's login name of the user. Consist of `username@primarydomain`                         |
| sub                | `77776025198584418`                      | Subject ID of the user                                                                        |

## Custom Claims

> This feature is not yet released

## Reserved Claims

ZITADEL reserves some claims to assert certain data.

| Claims                                          | Example                                                                                              | Description                                                                                                                                                                        |
|:------------------------------------------------|:-----------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| urn:zitadel:iam:org:domain:primary:{domainname} | `{"urn:zitadel:iam:org:domain:primary": "acme.ch"}`                                                  | This claim represents the primary domain of the organization the user belongs to.                                                                                                  |
| urn:zitadel:iam:org:project:roles:{rolename}    | `{"urn:zitadel:iam:org:project:roles": [ {"user": {"id1": "acme.zitade.ch", "id2": "caos.ch"} } ] }` | When roles are asserted, ZITADEL does this by providing the `id` and `primaryDomain` below the role. This gives you the option to check in which organization a user has the role. |
| urn:zitadel:iam:roles:{rolename}                | TBA                                                                                                  | TBA                                                                                                                                                                                |

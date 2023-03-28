---
title: Claims
---

ZITADEL asserts claims on different places according to the corresponding specifications or project and clients settings.
Please check below the matrix for an overview where which scope is asserted.

| Claims                                            | Userinfo       | Introspection  | ID Token                                    | Access Token                         |
|:--------------------------------------------------|:---------------|----------------|---------------------------------------------|--------------------------------------|
| acr                                               | No             | No             | Yes                                         | No                                   |
| address                                           | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| amr                                               | No             | No             | Yes                                         | No                                   |
| aud                                               | No             | Yes            | Yes                                         | When JWT                             |
| auth_time                                         | No             | No             | Yes                                         | No                                   |
| azp (client_id when Introspect)                   | No             | Yes            | Yes                                         | When JWT                             |
| email                                             | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| email_verified                                    | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| exp                                               | No             | Yes            | Yes                                         | When JWT                             |
| family_name                                       | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| gender                                            | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| given_name                                        | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| iat                                               | No             | Yes            | Yes                                         | When JWT                             |
| iss                                               | No             | Yes            | Yes                                         | When JWT                             |
| jti                                               | No             | Yes            | No                                          | When JWT                             |
| locale                                            | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| name                                              | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| nbf                                               | No             | Yes            | Yes                                         | When JWT                             |
| nonce                                             | No             | No             | Yes                                         | No                                   |
| phone                                             | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| phone_verified                                    | When requested | When requested | When requested amd response_type `id_token` | No                                   |
| preferred_username (username when Introspect)     | When requested | When requested | Yes                                         | No                                   |
| sub                                               | Yes            | Yes            | Yes                                         | When JWT                             |
| urn:zitadel:iam:org:domain:primary:{domainname}   | When requested | When requested | When requested                              | When JWT and requested               |
| urn:zitadel:iam:org:project:roles                 | When requested | When requested | When requested or configured                | When JWT and requested or configured |
| urn:zitadel:iam:user:metadata                     | When requested | When requested | When requested                              | When JWT and requested               |
| urn:zitadel:iam:user:resourceowner:id             | When requested | When requested | When requested                              | When JWT and requested               |
| urn:zitadel:iam:user:resourceowner:name           | When requested | When requested | When requested                              | When JWT and requested               |
| urn:zitadel:iam:user:resourceowner:primary_domain | When requested | When requested | When requested                              | When JWT and requested               |

## Standard Claims

| Claims             | Example                                  | Description                                                                                                                                            |
|:-------------------|:-----------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| acr                | TBA                                      | TBA                                                                                                                                                    |
| address            | `Teufener Strasse 19, 9000 St. Gallen`   | TBA                                                                                                                                                    |
| amr                | `pwd mfa`                                | Authentication Method References as defined in [RFC8176](https://tools.ietf.org/html/rfc8176) <br/> `password` value is deprecated, please check `pwd` |
| aud                | `69234237810729019`                      | The audience of the token, by default all client id's and the project id are included                                                                  |
| auth_time          | `1311280969`                             | Unix time of the authentication                                                                                                                        |
| azp                | `69234237810729234`                      | Client id of the client who requested the token                                                                                                        |
| email              | `road.runner@acme.ch`                    | Email Address of the subject                                                                                                                           |
| email_verified     | `true`                                   | Boolean if the email was verified by ZITADEL                                                                                                           |
| exp                | `1311281970`                             | Time the token expires (as unix time)                                                                                                                  |
| family_name        | `Runner`                                 | The subjects family name                                                                                                                               |
| gender             | `other`                                  | Gender of the subject                                                                                                                                  |
| given_name         | `Road`                                   | Given name of the subject                                                                                                                              |
| iat                | `1311280970`                             | Time of the token was issued at (as unix time)                                                                                                         |
| iss                | `{your_domain}`                          | Issuing domain of a token                                                                                                                              |
| jti                | `69234237813329048`                      | Unique id of the token                                                                                                                                 |
| locale             | `en`                                     | Language from the subject                                                                                                                              |
| name               | `Road Runner`                            | The subjects full name                                                                                                                                 |
| nbf                | `1311280970`                             | Time the token must not be used before (as unix time)                                                                                                  |
| nonce              | `blQtVEJHNTF0WHhFQmhqZ0RqeHJsdzdkd2d...` | The nonce provided by the client                                                                                                                       |
| phone              | `+41 79 XXX XX XX`                       | Phone number provided by the user                                                                                                                      |
| phone_verified     | `true`                                   | Boolean if the phone was verified by ZITADEL                                                                                                           |
| preferred_username | `road.runner@acme.caos.ch`               | ZITADEL's login name of the user. Consist of `username@primarydomain`                                                                                  |
| sub                | `77776025198584418`                      | Subject ID of the user                                                                                                                                 |

## Custom Claims

You can add custom claims using the [complement token flow](/docs/apis/actions/complement-token) of the [actions feature](/docs/apis/actions/introduction).

## Reserved Claims

ZITADEL reserves some claims to assert certain data. Please check out the [reserved scopes](scopes#reserved-scopes).

| Claims                                            | Example                                                                                              | Description                                                                                                                                                                        |
|:--------------------------------------------------|:-----------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| urn:zitadel:iam:action:{actionname}:log           | `{"urn:zitadel:iam:action:appendCustomClaims:log": ["test log", "another test log"]}`                | This claim is set during Actions as a log, e.g. if two custom claims with the same keys are set.                                                                                   |
| urn:zitadel:iam:org:domain:primary:{domainname}   | `{"urn:zitadel:iam:org:domain:primary": "acme.ch"}`                                                  | This claim represents the primary domain of the organization the user belongs to.                                                                                                  |
| urn:zitadel:iam:org:project:roles                 | `{"urn:zitadel:iam:org:project:roles": [ {"user": {"id1": "acme.zitade.ch", "id2": "caos.ch"} } ] }` | When roles are asserted, ZITADEL does this by providing the `id` and `primaryDomain` below the role. This gives you the option to check in which organization a user has the role. |
| urn:zitadel:iam:roles:{rolename}                  | TBA                                                                                                  | TBA                                                                                                                                                                                |
| urn:zitadel:iam:user:metadata                     | `{"urn:zitadel:iam:user:metadata": [ {"key": "VmFsdWU=" } ] }`                                       | The metadata claim will include all metadata of a user. The values are base64 encoded.                                                                                             |
| urn:zitadel:iam:user:resourceowner:id             | `{"urn:zitadel:iam:user:resourceowner:id": "orgid"}`                                                 | This claim represents the id of the resource owner organisation of the user.                                                                                                       |
| urn:zitadel:iam:user:resourceowner:name           | `{"urn:zitadel:iam:user:resourceowner:name": "ACME"}`                                                | This claim represents the name of the resource owner organisation of the user.                                                                                                     |
| urn:zitadel:iam:user:resourceowner:primary_domain | `{"urn:zitadel:iam:user:resourceowner:primary_domain": "acme.ch"}`                                   | This claim represents the primary domain of the resource owner organisation of the user.                                                                                           |

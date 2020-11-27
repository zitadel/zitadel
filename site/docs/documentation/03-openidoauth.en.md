---
title: OpenID Connect 1.0 & OAuth 2.0
---

### Endpoints and Domains

This chapter documents the [OpenID Connect 1.0](https://openid.net/connect/) and [OAuth 2.0](https://oauth.net/2/) features provided by **ZITADEL**.

Under normal circumstances **ZITADEL** need four domain names to operate properly.

| Domain Name | Example               | Description                                                                                                                          |
|:------------|:----------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| issuer      | `issuer.zitadel.ch`   | Provides the [OpenID Connect 1.0 Discovery Endpoint](#openid-connect-10-discovery)                                                   |
| api         | `api.zitadel.ch`      | All ZITADEL API's are located under this domain see [API explanation](develop#APIs) for details                                      |
| login       | `accounts.zitadel.ch` | The accounts.* page provides server renderer pages like login and register and as well the authorization_endpoint for OpenID Connect |
| console     | `console.zitadel.ch`  | With the console.* domain we serve the assets for the management gui                                                                 |

#### OpenID Connect 1.0 Discovery

The OpenID Connect Discovery Endpoint is located within the issuer domain.
For example with [zitadel.ch](zitadel.ch) this would be the domain [issuer.zitadel.ch](issuer.zitadel.ch). This would give us [https://issuer.zitadel.ch/.well-known/openid-configuration](https://issuer.zitadel.ch/.well-known/openid-configuration).

**Link to spec.** [OpenID Connect Discovery 1.0 incorporating errata set 1](https://openid.net/specs/openid-connect-discovery-1_0.html)

#### authorization_endpoint

[https://accounts.zitadel.ch/oauth/v2/authorize](https://accounts.zitadel.ch/oauth/v2/authorize)

> The authorization_endpoint is located with the login page, due to the need of accessing the same cookie domain

#### token_endpoint

[https://api.zitadel.ch/oauth/v2/token](https://api.zitadel.ch/oauth/v2/token)

#### userinfo_endpoint

[https://api.zitadel.ch/oauth/v2/userinfo](https://api.zitadel.ch/oauth/v2/userinfo)

#### end_session_endpoint

[https://accounts.zitadel.ch/oauth/v2/endsession](https://accounts.zitadel.ch/oauth/v2/endsession)

> The end_session_endpoint is located with the login page, due to the need of accessing the same cookie domain

#### jwks_uri

[https://api.zitadel.ch/oauth/v2/keys](https://api.zitadel.ch/oauth/v2/keys)

> Be aware that these keys can be rotated without any prior notice. We will however make sure that a proper `kid` is set with each key!

#### OAuth 2.0 Metadata

**ZITADEL** does not yet provide a OAuth 2.0 Metadata endpoint but instead provides a [OpenID Connect Discovery Endpoint](#openid-connect-10-discovery).

### Scopes

ZITADEL supports the usage of scopes as way of requesting information from the IAM and also instruct ZITADEL to do certain operations.

#### Standard Scopes

| Scopes  | Example   | Description                                          |
|:--------|:----------|------------------------------------------------------|
| openid  | `openid`  | When using openid connect this is a mandatory scope  |
| profile | `profile` | Optional scope to request the profile of the subject |
| email   | `email`   | Optional scope to request the email of the subject   |
| address | `address` | Optional scope to request the address of the subject |

#### Custom Scopes

> This feature is not yet released

#### Reserved Scopes

In addition to the standard compliant scopes we utilize the following scopes.

| Scopes   | Example    | Description    |
|:------------------------------------------------|:---------------|--------------------------------------------|
| urn:zitadel:iam:org:project:role:{rolename}     | `urn:zitadel:iam:org:project:role:user` | By using this scope a [client](administrate#clients) can request the claim urn:zitadel:iam:roles:rolename} to be asserted when possible. As an alternative approach you can enable all [roles](administrate#Roles) to be asserted from the [project](administrate#projects) a [client](administrate#clients) belongs to. See details [here](administrate#RBAC_Settings) |
| urn:zitadel:iam:org:domain:primary:{domainname} | `urn:zitadel:iam:org:domain:primary:acme.ch` |When requesting this scope **ZITADEL** will enforce that the user is a member of the selected organization. If the organization does not exist a failure is displayed |
| urn:zitadel:iam:role:{rolename}                 | | |
| urn:zitadel:iam:org:project:id:{projectid}:aud  | ZITADEL's Project id is `urn:zitadel:iam:org:project:id:69234237810729019:aud` | By adding this scope, the requested projectid will be added to the audience of the access and id token |

> If access to ZITADEL's API's is needed with a service user the scope `urn:zitadel:iam:org:project:id:69234237810729019:aud` needs to be used with the JWT Profile request

### Claims

ZITADEL asserts claims on different places according to the corresponding specifications or project and clients settings.
Please check below the matrix for an overview where which scope is asserted.

| Claims                                          | Userinfo           | ID Token                               | Access Token                             |
|:------------------------------------------------|:-------------------|----------------------------------------|------------------------------------------|
| acr                                             | Yes                | Yes                                    | No                                       |
| address                                         | Yes when requested | Yes only when response type `id_token` | No                                       |
| amr                                             | Yes                | Yes                                    | No                                       |
| aud                                             | No                 | Yes                                    | Yes when JWT                             |
| auth_time                                       | Yes                | Yes                                    | No                                       |
| azp                                             | No                 | Yes                                    | Yes when JWT                             |
| email                                           | Yes when requested | Yes only when response type `id_token` | No                                       |
| email_verified                                  | Yes when requested | Yes only when response type `id_token` | No                                       |
| exp                                             | No                 | Yes                                    | Yes when JWT                             |
| family_name                                     | Yes when requested | Yes when requested                     | No                                       |
| gender                                          | Yes when requested | Yes when requested                     | No                                       |
| given_name                                      | Yes when requested | Yes when requested                     | No                                       |
| iat                                             | No                 | Yes                                    | Yes when JWT                             |
| iss                                             | No                 | Yes                                    | Yes when JWT                             |
| locale                                          | Yes when requested | Yes when requested                     | No                                       |
| name                                            | Yes when requested | Yes when requested                     | No                                       |
| nonce                                           | No                 | Yes                                    | No                                       |
| phone                                           | Yes when requested | Yes only when response type `id_token` | No                                       |
| preferred_username                              | Yes when requested | Yes                                    | No                                       |
| sub                                             | Yes                | Yes                                    | Yes when JWT                             |
| urn:zitadel:iam:org:domain:primary:{domainname} | Yes when requested | Yes when requested                     | Yes when JWT and requested               |
| urn:zitadel:iam:org:project:roles:{rolename}    | Yes when requested | Yes when requested or configured       | Yes when JWT and requested or configured |

#### Standard Claims

| Claims             | Example                                  | Description                                                                                   |
|:-------------------|:-----------------------------------------|-----------------------------------------------------------------------------------------------|
| acr                | TBA                                      | TBA                                                                                           |
| address            | `Teufener Strasse 19, 9000 St. Gallen`   |                                                                                               |
| amr                | `"amr": "pwd mfa"`                       | Authentication Method References as defined in [RFC8176](https://tools.ietf.org/html/rfc8176) |
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
| iss                | `"iss": "https://issuer.zitadel.ch"`     | Issuing domain of a token                                                                     |
| locale             | `en`                                     | Language from the subject                                                                     |
| name               | `Road Runner`                            | The subjects full name                                                                        |
| nonce              | `blQtVEJHNTF0WHhFQmhqZ0RqeHJsdzdkd2d...` | The nonce provided by the client                                                              |
| phone              | `+41 71 XXX XX XX`                       |                                                                                               |
| preferred_username | `road.runner@acme.caos.ch`               |                                                                                               |
| sub                | `77776025198584418`                      | Subject ID of the user                                                                        |

#### Custom Claims

> This feature is not yet released

#### Reserved Claims

ZITADEL reserves some claims to assert certain data.

| Claims                                          | Example                                                                                              | Description                                                                                                                                                                        |
|:------------------------------------------------|:-----------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| urn:zitadel:iam:org:domain:primary:{domainname} | `{"urn:zitadel:iam:org:domain:primary": "acme.ch"}`                                                  | This claim represents the primary domain of the organization the user belongs to.                                                                                                  |
| urn:zitadel:iam:org:project:roles:{rolename}    | `{"urn:zitadel:iam:org:project:roles": [ {"user": {"id1": "acme.zitade.ch", "id2": "caos.ch"} } ] }` | When roles are asserted, ZITADEL does this by providing the `id` and `primaryDomain` below the role. This gives you the option to check in which organization a user has the role. |
| urn:zitadel:iam:roles:{rolename}                | TBA                                                                                                  | TBA                                                                                                                                                                                |

### Grant Types

For a list of supported or unsupported `Grant Types` please have a look at the table below.

| Grant Type                                            | Supported           |
|:------------------------------------------------------|:--------------------|
| Authorization Code                                    | yes                 |
| Authorization Code with PKCE                          | yes                 |
| Client Credentials                                    | yes                 |
| Device Authorization                                  | under consideration |
| Implicit                                              | yes                 |
| JSON Web Token (JWT) Profile                          | partially           |
| Refresh Token                                         | work in progress    |
| Resource Owner Password Credentials                   | no                  |
| Security Assertion Markup Language (SAML) 2.0 Profile | no                  |
| Token Exchange                                        | work in progress    |

#### Authorization Code

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.1](https://tools.ietf.org/html/rfc6749#section-1.3.1)

#### Proof Key for Code Exchange

**Link to spec.** [Proof Key for Code Exchange by OAuth Public Clients](https://tools.ietf.org/html/rfc7636)

#### Implicit

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.2](https://tools.ietf.org/html/rfc6749#section-1.3.2)

#### Client Credentials

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.4](https://tools.ietf.org/html/rfc6749#section-1.3.4)

#### Refresh Token

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.5](https://tools.ietf.org/html/rfc6749#section-1.5)

#### JSON Web Token (JWT) Profile

**Link to spec.** [JSON Web Token (JWT) Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7523)

##### Using JWTs as Authorization Grants

Our service user work with the JWT profile to authenticate them against ZITADEL.

1. Create or use an existing service user
2. Create a new key and download it
3. Generate a JWT with the structure below and sing it with the downloaded key
4. Send the JWT Base64 encoded to ZITADEL's token endpoint
5. Use the received access token

---

Key JSON

| Key    | Example                                                           | Description                                                        |
|:-------|:------------------------------------------------------------------|:-------------------------------------------------------------------|
| type   | `"serviceaccount"`                                                  | The type of account, right now only serviceaccount is valid        |
| keyId  | `"81693565968772648"`                                               | This is unique ID of the key                                       |
| key    | `"-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----"` | The private key generated by ZITADEL, this can not be regenerated! |
| userId | `78366401571647008`                                                 | The service users ID, this is the same as the subject from tokens  |

```JSON
{
	"type": "serviceaccount",
	"keyId": "81693565968772648",
	"key": "-----BEGIN RSA PRIVATE KEY-----...-----END RSA PRIVATE KEY-----",
	"userId": "78366401571647008"
}
```

---

JWT

| Claim | Example                       | Description                                                                      |
|:------|:------------------------------|:---------------------------------------------------------------------------------|
| aud   | `"https://issuer.zitadel.ch"` | String or Array of intended audiences MUST include ZITADEL's issuing domain      |
| exp   | `1605183582`                  | Unix timestamp of the expiry, MUST NOT be longer than 1h                         |
| iat   | `1605179982`                  | Unix timestamp of the creation singing time of the JWT                           |
| iss   | `"http://localhost:50003"`    | String which represents the requesting party                                     |
| sub   | `"77479219772321307"`         | The subject ID of the service user, normally the `userId` from the json key file |

```JSON
{
	"iss": "http://localhost:50003",
	"sub": "77479219772321307",
	"aud": "https://issuer.zitadel.ch",
	"exp": 1605183582,
	"iat": 1605179982
}
```

---

Access Token Request

| Parameter    | Example                                                                   | Description                                   |
|:-------------|:--------------------------------------------------------------------------|:----------------------------------------------|
| Content-Type | `application/x-www-form-urlencoded`                                         |                                               |
| grant_type   | `urn:ietf:params:oauth:grant-type:jwt-bearer`                               | Using JWTs as Authorization Grants            |
| assertion    | `eyJhbGciOiJSUzI1Ni...`                                                     | The base64 encoded JWT created above          |
| scope        | `openid profile email urn:zitadel:iam:org:project:id:69234237810729019:aud` | Scopes you would like to request from ZITADEL |

```BASH
curl --request POST \
  --url https://api.zitadel.ch/oauth/v2/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
  --data assertion=eyJhbGciOiJSUzI1Ni...
  --data scope=openid profile email address
```

##### Using JWTs for Client Authentication

> Not yet supported

#### Token Exchange

**Link to spec.** [OAuth 2.0 Token Exchange](https://tools.ietf.org/html/rfc8693)

### Device Authorization

**Link to spec.** [OAuth 2.0 Device Authorization Grant](https://tools.ietf.org/html/rfc8628)

### Not Supported Grant Types

#### Resource Owner Password Credentials

> Due to growing security concerns we do not support this grant type. With OAuth 2.1 it looks like this grant will be removed.

**Link to spec.** [OThe OAuth 2.0 Authorization Framework Section 1.3.3](https://tools.ietf.org/html/rfc6749#section-1.3.3)

#### Security Assertion Markup Language (SAML) 2.0 Profile

**Link to spec.** [Security Assertion Markup Language (SAML) 2.0 Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7522)

---
title: OpenID Connect 1.0 & OAuth 2.0
---

### Endpoints

This chapter documents the OpenID Connect 1.0 and OAuth 2.0 features provided by **ZITADEL**.

#### OpenID Connect 1.0 Discovery

The OpenID Connect Discovery Endpoint is located with the issuer domain.
For example with [zitadel.ch](zitadel.ch) this would be the domain [issuer.zitadel.ch](issuer.zitadel.ch). This would give us [https://issuer.zitadel.ch/.well-known/openid-configuration](https://issuer.zitadel.ch/.well-known/openid-configuration).

**Link to spec.** [OpenID Connect Discovery 1.0 incorporating errata set 1](https://openid.net/specs/openid-connect-discovery-1_0.html)

#### authorization_endpoint

[https://accounts.zitadel.ch/oauth/v2/authorize](https://accounts.zitadel.ch/oauth/v2/authorize)

#### token_endpoint

#### userinfo_endpoint

#### end_session_endpoint

#### jwks_uri

#### OAuth 2.0 Metadata

ZITADEL does not provide a OAuth 2.0 Metadata endpoint but instead provides a OpenID Connect Discovery Endpoint.
TODO: Insert Link.

### Grant Types

| Grant Type                                            | Supported           |
|:------------------------------------------------------|:--------------------|
| Authorization Code                                    | yes                 |
| Implicit                                              | yes                 |
| Resource Owner Password Credentials                   | no                  |
| Client Credentials                                    | yes                 |
| Device Authorization                                  | under consideration |
| Refresh Token                                         | work in progress    |
| JSON Web Token (JWT) Profile                          | yes                 |
| Security Assertion Markup Language (SAML) 2.0 Profile | no                  |
| Token Exchange                                        | work in progress    |

#### Authorization Code

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.1](https://tools.ietf.org/html/rfc6749#section-1.3.1)  

#### Implicit

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.2](https://tools.ietf.org/html/rfc6749#section-1.3.2)

#### Client Credentials

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.3.4](https://tools.ietf.org/html/rfc6749#section-1.3.4)

#### Refresh Token

**Link to spec.** [The OAuth 2.0 Authorization Framework Section 1.5](https://tools.ietf.org/html/rfc6749#section-1.5)

#### JSON Web Token (JWT) Profile

**Link to spec.** [JSON Web Token (JWT) Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7523)

#### Token Exchange

**Link to spec.** [OAuth 2.0 Token Exchange](https://tools.ietf.org/html/rfc8693)

### Device Authorization

**Link to spec.** [OAuth 2.0 Device Authorization Grant](https://tools.ietf.org/html/rfc8628)

### Not Supported Grant Types

#### Resource Owner Password Credentials

> Due to growing security concern we do not support this grant type. With OAuth 2.1 it looks like this grant will be removed.

**Link to spec.** [OThe OAuth 2.0 Authorization Framework Section 1.3.3](https://tools.ietf.org/html/rfc6749#section-1.3.3)

#### Security Assertion Markup Language (SAML) 2.0 Profile

**Link to spec.** [Security Assertion Markup Language (SAML) 2.0 Profile for OAuth 2.0 Client Authentication and Authorization Grants](https://tools.ietf.org/html/rfc7522)
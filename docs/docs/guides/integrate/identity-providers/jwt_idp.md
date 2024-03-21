---
title: JWT IdP
sidebar_label: JWT IdP
---


## Configure a JWT as an Identity Provider in ZITADEL

JSON Web Token Identity Provider (JWT IdP) gives you the possibility to use an (existing) JWT as federated identity.
Imagine you have a Web Application Firewall (WAF) which handles your session for an existing application.
You're now creating a new application which uses ZITADEL as its IDP / Authentication Server.
The new app might even be opened from within the existing application and you want to reuse the session information for the new application.
This is where JWT IdP comes into place.

All you need to provide is an endpoint where ZITADEL can receive a JWT and some information for its signature verification.

## Authentication using JWT IdP

The authentication process then might look like the following:

![JWT IdP Architecture](/img/guides/jwt_idp.png)

1. The user is logged into the existing application and the WAF holds the session information. It might even send a JWT to the application.
   The new application is opened by clicking on a link in the existing application.
2. The application bootstraps and since it cannot find a session, it will create an OIDC Authorization Request to ZITADEL.
   In this request it provides a scope to directly request the JWT IdP.
3. ZITADEL will do so and redirect to the preconfigured JWT Endpoint. While the endpoint is behind the WAF, ZITADEL is able to receive a JWT from a defined HTTP header.
   It will then validate its signature, which might require to call the configured Keys Endpoint.
   If the signature is valid and token is not expired, ZITADEL will then use the token and the enclosed `sub` claim as if it was an id_token returned from an OIDC IDP:
   It will try to match a user's external identity and if not possible, create a new user with the information of the provided token.
   Prerequisite for this is that the IDP setting `autoregister` is set to `true`.
4. ZITADEL will then redirect to its main instance and the login flow will proceed.
5. The user will be redirected to the Callback Endpoint of the new Application, where the application will exchange the code for tokens.
   The user is finally logged in the new application, without any direct interaction.

### Terms and example values

To further explain and illustrate how a JWT IdP works, we will assume the following:

- the **Existing Application** is deployed under `apps.test.com/existing/`
- the **New Application** is deployed under `new.test.com`
- the **Login UI of ZITADEL** is deployed under `accounts.test.com`

The **JWT IdP Configuration** might then be:
  - **JWT Endpoint** (Endpoint where ZITADEL will redirect to):<br/>`https://apps.test.com/existing/auth-new`
  - **Issuer** (of the JWT):<br/>`https://issuer.test.internal`
  - **Keys Endpoint** (where keys of the JWT Signature can be gathered):<br/>`https://issuer.test.internal/keys`
  - **Header Name** (of the JWT, Authorization if omitted):<br/>`x-custom-tkn`

Therefore, if the user is redirected from ZITADEL to the JWT Endpoint on the WAF (`https://apps.test.com/existing/auth-new`), 
the session cookies previously issued by the WAF, will be sent along by the browser due to the path being on the same domain as the exiting application.
The WAF will reuse the session and send the JWT in the HTTP header `x-custom-tkn` to its upstream, the ZITADEL JWT Endpoint (`https://accounts.test.com/ui/login/login/jwt/authorize`).

For the signature validation, ZITADEL must be able to connect to Keys Endpoint (`https://issuer.test.internal/keys`) 
and it will check if the token was signed (claim `iss`) by the defined Issuer (`https://issuer.test.internal`).


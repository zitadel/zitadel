---
title:  Proxy / WAF
description: ...
---

### Proxy Protocol and Flow recommendation

### Ambassador Example

According to [https://www.getambassador.io/docs/latest/](https://www.getambassador.io/docs/latest/) Ambassador is a:

>The Ambassador Edge Stack is a comprehensive, self-service edge stack and API Gateway for Kubernetes built on Envoy Proxy. The shift to Kubernetes and microservices has profound consequences for the capabilities you need at the edge, as well as how you manage the edge. The Ambassador Edge Stack has been engineered with this world in mind.

#### Configure ZITADEL for Ambassador

You can use **ZITADEL** for Authentication and Authorization with **Ambassador**.

> The redirect URI is `https://{AMBASSADOR_URL}/.ambassador/oauth2/redirection-endpoint`

#### Use Ambassador to Authenticate with ZITADEL

With this you can use Ambassador to initiate the Authorization Code Flow.

```yaml
apiVersion: getambassador.io/v2
kind: Filter
metadata:
  name: zitadel-filter
  namespace: default
spec:
  OAuth2:
    authorizationURL: https://accounts.zitadel.ch/oauth/v2/authorize
    clientID: {ZITADEL_GENERATED_CLIENT_ID}
    secret: {ZITADEL_GENERATED_CLIENT_SECRET}
    protectedOrigins:
    - origin: https://{PROTECTED_URL}
```

```yaml
apiVersion: getambassador.io/v2
kind: FilterPolicy
metadata:
  name: zitadel-policy
  namespace: default
spec:
  rules:
    - host: "*"
      path: /backend/get-quote/
      filters:
        - name: zitadel-filter
```

#### Use Ambassador to check JWT Bearer Tokens

If you would like **Ambassador** to verify a JWT token from the autorization header you can do so by configuring **ZITADEL's** endpoints.

> Make sure that in your client settings of **ZITADEL** the "AuthToken Options" is **JWT** by default **ZITADEL** will use opaque tokens!

```yaml
apiVersion: getambassador.io/v2
kind: Filter
metadata:
  name: zitadel-filter
  namespace: default
spec:
  JWT:
    jwksURI:            "https://api.zitadel.ch/oauth/v2/keys"
    validAlgorithms:
    - "RS256"
    issuer:             "https://issuer.zitadel.ch"
    requireIssuer:      true
```

```yaml
apiVersion: getambassador.io/v2
kind: FilterPolicy
metadata:
  name: zitadel-policy
  namespace: default
spec:
  rules:
    - host: "*"
      path: /backend/get-quote/
      filters:
        - name: zitadel-filter
```

> Additional Infos can be found with [Ambassadors Documentation](https://www.getambassador.io/docs/latest/howtos/oauth-oidc-auth/)

### NGINX Example

> TODO

### OAuth2 Proxy Example

[OAuth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy) is a project which allows services to delegate the authentication flow to a IDP, for example **ZITADEL**

> Right now the OAuth 2.0 proxy is not [spec. compliant](https://openid.net/specs/openid-connect-core-1_0.html#ScopeClaims) because **ZITADEL** does not assert information into the ID Token when an access Token is delivered, we will change this however with [ISSUE #940](https://github.com/caos/zitadel/issues/940)

```toml
provider = "oidc"
provider_display_name "ZITADEL"
redirect_url = "http://127.0.0.1:4180/oauth2/callback"
oidc_issuer_url = "https://issuer.zitadel.ch"
upstreams = [
    "https://example.corp.com"
]
email_domains = [
    "*"
]
client_id = "{ZITADEL_GENERATED_CLIENT_ID}"
client_secret = "{ZITADEL_GENERATED_CLIENT_SECRET}"
pass_access_token = true
cookie_secret = "{SUPPLY_SOME_SECRET_HERE}"
skip_provider_button = true
```

### Cloudflare Access Example

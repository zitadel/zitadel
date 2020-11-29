---
title:  Proxy / WAF
description: ...
---

### Proxy Protocol and Flow recommendation

### Ambassador Example

According to [https://www.getambassador.io/docs/latest/](https://www.getambassador.io/docs/latest/) Ambassador is a:

>The Ambassador Edge Stack is a comprehensive, self-service edge stack and API Gateway for Kubernetes built on Envoy Proxy. The shift to Kubernetes and microservices has profound consequences for the capabilities you need at the edge, as well as how you manage the edge. The Ambassador Edge Stack has been engineered with this world in mind.

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

### OAuth2 Proxy Example

[OAuth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy) is a project which allows services to delegate the authentication flow to a IDP, for example **ZITADEL**

#### OAuth2 Proxy Authentication Example

```toml
provider = "oidc"
user_id_claim = "sub" #uses the subject as ID instead of the email
provider_display_name = "ZITADEL"
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
cookie_secure = false #localdev only false
http_address = "127.0.0.1:4180" #localdev only
```

> This was tested with version `oauth2-proxy v6.1.1 (built with go1.14.2)`

#### OAuth2 Proxy Authorization Example

> Not yet supported but with the work of [https://github.com/oauth2-proxy/oauth2-proxy/pull/797](https://github.com/oauth2-proxy/oauth2-proxy/pull/797) it should be possible in the future

### Cloudflare Access Example

> TODO

### NGINX Example

> TODO

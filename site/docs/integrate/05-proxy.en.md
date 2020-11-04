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

#### Ambassador Filter Authentication

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

#### Ambassador Filter Authorisation

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


#### Ambassador FilterPolicy

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

### OAuth2 Proxy Example

---
Some text

```toml
provider = "ZITADEL"
redirect_url = "https://example.corp.com/oauth2/callback"
oidc_issuer_url = "https://issuer.zitadel.ch/.well-known/openid-configuration"
upstreams = [
    "https://example.corp.com"
]
email_domains = [
    "corp.com"
]
client_id = "XXXXX"
client_secret = "YYYYY"
pass_access_token = true
cookie_secret = "ZZZZZ"
skip_provider_button = true
```

### Cloudflare Access Example
